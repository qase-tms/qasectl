package maestro

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	models "github.com/qase-tms/qasectl/internal/models/result"
)

const commandsFile = "commands.json"

type Parser struct {
	path  string
	steps string
}

func NewParser(path, steps string) *Parser {
	return &Parser{path: path, steps: steps}
}

func (p *Parser) Parse() ([]models.Result, error) {
	dirs, err := findFlowDirs(p.path)
	if err != nil {
		return nil, err
	}
	if len(dirs) == 0 {
		return nil, fmt.Errorf("no %s found in %s: run Maestro with --debug-output to produce it", commandsFile, p.path)
	}

	var flows []flow
	for _, dir := range dirs {
		f, err := loadFlow(dir)
		if err != nil {
			slog.Warn("skipping flow", "dir", dir, "error", err)
			continue
		}
		flows = append(flows, f)
	}
	if err := checkDuplicates(flows); err != nil {
		return nil, err
	}

	devices := deviceNames(p.path)
	sharded := hasSeveralShards(flows)
	var results []models.Result
	var missingDevice bool
	for _, f := range flows {
		if _, ok := devices[f.UDID]; !ok {
			missingDevice = true
		}
		results = append(results, p.convert(f, devices, sharded))
	}
	if missingDevice {
		slog.Info("Maestro JUnit report not found; device names are not set. Run Maestro with --format junit --output <run-dir>/report.xml")
	}
	return results, nil
}

func findFlowDirs(root string) ([]string, error) {
	var dirs []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && d.Name() == commandsFile {
			dirs = append(dirs, filepath.Dir(path))
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan %s: %w", root, err)
	}
	sort.Strings(dirs)
	return dirs, nil
}

func loadFlow(dir string) (flow, error) {
	data, err := os.ReadFile(filepath.Join(dir, commandsFile))
	if err != nil {
		return flow{}, err
	}
	var cmds []command
	if err := json.Unmarshal(data, &cmds); err != nil {
		return flow{}, fmt.Errorf("invalid %s: %w", commandsFile, err)
	}
	if len(cmds) == 0 {
		return flow{}, errors.New("no commands")
	}

	f := flow{Dir: dir, Commands: cmds}
	var hasConfig, seenEnv bool
	for _, c := range cmds {
		if c.Depth != 0 {
			continue
		}
		switch c.Type {
		case "applyConfigurationCommand":
			if hasConfig {
				continue
			}
			var a struct {
				Config struct {
					Name       string            `json:"name"`
					Tags       []string          `json:"tags"`
					Properties map[string]string `json:"properties"`
				} `json:"config"`
			}
			if err := json.Unmarshal(c.Args, &a); err != nil {
				return flow{}, fmt.Errorf("invalid applyConfigurationCommand: %w", err)
			}
			hasConfig = true
			f.Name = a.Config.Name
			f.Tags = a.Config.Tags
			f.Properties = a.Config.Properties
		case "defineVariablesCommand":
			if seenEnv {
				continue
			}
			var a struct {
				Env map[string]string `json:"env"`
			}
			if err := json.Unmarshal(c.Args, &a); err != nil {
				return flow{}, fmt.Errorf("invalid defineVariablesCommand: %w", err)
			}
			seenEnv = true
			f.FileName = a.Env["MAESTRO_FILENAME"]
			f.UDID = a.Env["MAESTRO_DEVICE_UDID"]
			f.ShardID = a.Env["MAESTRO_SHARD_ID"]
		}
	}
	if !hasConfig {
		return flow{}, errors.New("no applyConfigurationCommand")
	}
	if f.Name == "" && f.FileName == "" {
		return flow{}, errors.New("flow has neither config name nor MAESTRO_FILENAME")
	}
	return f, nil
}

func (p *Parser) convert(f flow, devices map[string]string, sharded bool) models.Result {
	title := f.Name
	if title == "" {
		title = f.FileName
	}

	status, interrupted := flowStatus(f.Commands)
	start := float64(f.Commands[0].Timestamp)
	var end int64
	for _, c := range f.Commands {
		if e := c.Timestamp + c.Duration; e > end {
			end = e
		}
	}
	endF := float64(end)
	dur := endF - start

	ids := parseQaseIDs(f.Properties[qaseIDKey])
	suites := parseSuites(f.Properties[qaseSuiteKey])
	var params map[string]string
	if name, ok := devices[f.UDID]; ok {
		params = map[string]string{"device": name}
	}
	sigParts := append(append([]string{}, suites...), f.FileName, title)
	sig := models.GenerateSignature(ids, sigParts, params)
	res := models.Result{
		Title:     title,
		Signature: &sig,
		Fields:    fieldsFor(f),
		Params:    params,
		Execution: models.Execution{
			Status:    status,
			StartTime: &start,
			EndTime:   &endF,
			Duration:  &dur,
		},
	}
	if sharded && f.ShardID != "" {
		thread := "shard-" + f.ShardID
		res.Execution.Thread = &thread
	}
	if len(ids) > 0 {
		res.TestOpsIDs = &ids
	}
	for _, s := range suites {
		res.Relations.Suite.Data = append(res.Relations.Suite.Data, models.SuiteData{Title: s})
	}
	res.Steps = buildSteps(f.Dir, f.Commands, strings.EqualFold(p.steps, "all"))
	if fail := deepestFailure(f.Commands); fail != nil && fail.Error != nil {
		msg, trace := fail.Error.Message, fail.Error.DebugMessage
		res.Message = &msg
		if trace != "" {
			res.Execution.StackTrace = &trace
		}
	} else if interrupted {
		msg := "Flow did not finish"
		res.Message = &msg
	}
	if status == "failed" {
		res.Attachments = failureLogs(f.Dir)
	}
	if status == "passed" && hadFailedAttempt(f.Commands) {
		res.Fields["is_flaky"] = "true"
	}
	return res
}

func flowStatus(cmds []command) (status string, interrupted bool) {
	status = "passed"
	for _, c := range cmds {
		switch {
		case c.Status == "RUNNING" || c.Status == "PENDING":
			interrupted = true
			status = "failed"
		case c.Depth == 0 && c.Status == "FAILED":
			status = "failed"
		}
	}
	return status, interrupted
}
