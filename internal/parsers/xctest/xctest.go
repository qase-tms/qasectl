package xctest

import (
	"encoding/json"
	"fmt"
	models "github.com/qase-tms/qasectl/internal/models/result"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Parser is a parser for XCTest files
type Parser struct {
	path  string
	level stepLevel
}

// NewParser creates a new Parser
func NewParser(path, level string) (*Parser, error) {
	if !strings.HasSuffix(path, ".xcresult") {
		return nil, fmt.Errorf("unsupported format: %s", path)
	}

	return &Parser{
		path:  path,
		level: parseStepLevel(level),
	}, nil
}

var (
	layoutTime        = "2006-01-02T15:04:05.000-0700"
	caseId            *int64
	failures          map[string]FailureSummary
	processedFailures []string
)

const (
	internalStep   = "com.apple.dt.xctest.activity-type.internal"
	deleteStep     = "com.apple.dt.xctest.activity-type.deletedAttachment"
	attachmentStep = "com.apple.dt.xctest.activity-type.attachmentContainer"
	qaseConfig     = "Qase config"
	heicExt        = ".heic"
)

// Parse parses the XCTest file and returns the results
func (p *Parser) Parse() ([]models.Result, error) {
	const op = "xctest.Parser.Parse"
	logger := slog.With("op", op)

	logger.Debug("parsing XCTest file", "path", p.path)

	testPlanSumIDs, err := p.getTestPlanSumIDs()
	if err != nil {
		return nil, fmt.Errorf("failed to get test plan summaries: %w", err)
	}

	logger.Debug("got test plan summary IDs", "testPlanSumIDs", testPlanSumIDs)

	testPlanSums, err := p.getTestPlanSums(testPlanSumIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get test plan summaries: %w", err)
	}

	logger.Debug("got test plan summaries", "testPlanSums", testPlanSums)

	tests := p.getXCTests(testPlanSums)

	logger.Debug("got tests", "xctest", tests)

	results := p.getResults(tests)

	logger.Debug("got results", "results", results)

	return results, nil
}

func (p *Parser) readJson(id *string) ([]byte, error) {
	args := []string{"xcresulttool", "get", "--path", p.path, "--format", "json"}
	if id != nil {
		args = append(args, "--id", *id)
	}

	executeCommand := func(args []string) ([]byte, error) {
		out, err := exec.Command("xcrun", args...).Output()
		if err != nil {
			return nil, fmt.Errorf("failed to get XCResult: %w", err)
		}
		return out, nil
	}

	out, err := executeCommand(args)
	if err != nil {
		if strings.Contains(err.Error(), "exit status 64") {
			args = append(args, "--legacy")
			out, err = executeCommand(args)
			if err != nil {
				return nil, fmt.Errorf("failed to get XCResult with --legacy: %w", err)
			}
			return out, nil
		}
		return nil, err
	}

	return out, nil
}

func (p *Parser) readAttachment(id string) ([]byte, error) {
	args := []string{"xcresulttool", "get", "--path", p.path, "--format", "raw", "--id", id}

	executeCommand := func(args []string) ([]byte, error) {
		out, err := exec.Command("xcrun", args...).Output()
		if err != nil {
			return nil, fmt.Errorf("failed to get XCResult attachment: %w", err)
		}
		return out, nil
	}

	out, err := executeCommand(args)
	if err != nil {
		if strings.Contains(err.Error(), "exit status 64") {
			args = append(args, "--legacy")
			out, err = executeCommand(args)
			if err != nil {
				return nil, fmt.Errorf("failed to get XCResult attachment with --legacy: %w", err)
			}
		} else {
			return nil, err
		}
	}

	return out, nil
}

func (p *Parser) getTestPlanSumIDs() (map[string]TestMeta, error) {
	out, err := p.readJson(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get test plan summaries: %w", err)
	}

	var structure Structure
	err = json.Unmarshal(out, &structure)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal test plan summaries: %w", err)
	}

	testPlanSumIDs := make(map[string]TestMeta)

	for _, action := range structure.Actions.Values {
		if action.ActionResult.TestsRef == nil {
			continue
		}

		m := TestMeta{}

		if action.StartedTime != nil {
			st := float64(parseTime(action.StartedTime.Value).Unix())
			m.StartTime = &st
		}
		if action.EndedTime != nil {
			et := float64(parseTime(action.StartedTime.Value).Unix())
			m.EndTime = &et
		}
		if action.RunDestination != nil {
			m.Device = action.RunDestination.DisplayName.Value
		}

		testPlanSumIDs[action.ActionResult.TestsRef.ID.Value] = m
	}

	return testPlanSumIDs, nil
}

func (p *Parser) getTestPlanSums(IDs map[string]TestMeta) (map[TestMeta][]ActionTestPlanRunSummaries, error) {
	testPlanSums := make(map[TestMeta][]ActionTestPlanRunSummaries)

	for k, v := range IDs {
		out, err := p.readJson(&k)
		if err != nil {
			return nil, fmt.Errorf("failed to get test plan summaries: %w", err)
		}
		var sum ActionTestPlanRunSummaries
		err = json.Unmarshal(out, &sum)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal test plan summaries: %w", err)
		}

		testPlanSums[v] = append(testPlanSums[v], sum)
	}

	return testPlanSums, nil
}

func (p *Parser) getActionTestSummary(ID string) (ActionTestSummary, error) {
	out, err := p.readJson(&ID)
	if err != nil {
		return ActionTestSummary{}, fmt.Errorf("failed to get action test summary: %w", err)
	}

	var summary ActionTestSummary
	err = json.Unmarshal(out, &summary)
	if err != nil {
		return ActionTestSummary{}, fmt.Errorf("failed to unmarshal action test summary: %w", err)
	}

	return summary, nil
}

func (p *Parser) getXCTests(testPlanSums map[TestMeta][]ActionTestPlanRunSummaries) []XCTest {
	tests := make([]XCTest, 0)

	for k, v := range testPlanSums {
		for _, testPlanSum := range v {
			for _, testPlanRunSum := range testPlanSum.Summaries.Values {
				k.Configuration = testPlanRunSum.Name.Value
				for _, testSum := range testPlanRunSum.TestableSummaries.Values {
					for _, test := range testSum.Tests.Values {
						if test.Subtests != nil {
							tt := p.getXCTestsFromSubtest(*test.Subtests, make([]string, 0))
							for _, t := range tt {
								t.Metadata = k
								tests = append(tests, t)
							}
						}
					}
				}
			}
		}
	}

	return tests
}

func (p *Parser) getXCTestsFromSubtest(s Subtests, suites []string) []XCTest {
	const op = "xctest.Parser.getXCTestsFromSubtest"
	logger := slog.With("op", op)

	var tests []XCTest

	for _, v := range s.Values {
		if v.SummaryRef != nil {
			act, err := p.getActionTestSummary(v.SummaryRef.ID.Value)
			if err != nil {
				logger.Error("failed to get action test summary", "error", err, "ID", v.SummaryRef.ID.Value)
				continue
			}

			d, err := strconv.ParseFloat(v.Duration.Value, 64)
			if err != nil {
				logger.Error("failed to parse duration", "error", err)
				d = 0
			}

			test := XCTest{
				Name:      v.Name.Value,
				Action:    act,
				Suites:    make([]string, 0),
				Signature: v.IdentifierURL.Value,
				Duration:  d,
			}

			if len(suites) > 0 {
				test.Suites = append(test.Suites, suites...)
			}

			tests = append(tests, test)

			continue
		}

		if v.Subtests != nil {
			suites = append(suites, v.Name.Value)

			tt := p.getXCTestsFromSubtest(*v.Subtests, suites)
			tests = append(tests, tt...)
		}
	}

	return tests
}

func (p *Parser) getResults(tests []XCTest) []models.Result {
	const op = "xctest.Parser.getResults"
	logger := slog.With("op", op)

	var results []models.Result

	for _, t := range tests {
		caseId = nil
		failures = make(map[string]FailureSummary)
		processedFailures = make([]string, 0)

		logger := logger.With("testTitle", t.Name)
		logger.Debug("processing test", "test", t)

		suites := make([]models.SuiteData, 0)
		suites = append(suites, models.SuiteData{Title: t.Metadata.Suite})
		for _, s := range t.Suites {
			suites = append(suites, models.SuiteData{
				Title: s,
			})
		}

		result := models.Result{
			Title:     t.Name,
			Signature: &t.Signature,
			Execution: models.Execution{
				Status:    getStatus(t.Action.TestStatus.Value),
				Duration:  &t.Duration,
				StartTime: t.Metadata.StartTime,
				EndTime:   t.Metadata.EndTime,
			},
			Fields:      map[string]string{},
			Attachments: make([]models.Attachment, 0),
			Params: map[string]string{
				"Device":        t.Metadata.Device,
				"Configuration": t.Metadata.Configuration,
			},
			Relations: models.Relation{
				Suite: models.Suite{
					Data: suites,
				},
			},
			StepType: "text",
			Muted:    false,
		}

		if t.Action.FailureSummaries != nil {
			var message string

			for _, f := range t.Action.FailureSummaries.Values {
				failures[f.UUID.Value] = f
				message += f.Message.Value + "\n"

			}

			result.Execution.StackTrace = &message
		}

		steps, _, a := p.getSteps(t.Action.ActivitySummaries, 0)
		result.Attachments = append(result.Attachments, a...)

		result.Steps = steps
		if caseId != nil {
			result.TestOpsID = caseId
		}

		for k, v := range failures {
			if isFailureProcessed(k) {
				continue
			}

			if v.Attachments != nil {
				for _, att := range v.Attachments.Values {
					a, err := p.getAttachment(att)
					if err != nil {
						logger.Error("failed to get attachments", "error", err)
						continue
					}

					result.Attachments = append(result.Attachments, a)
				}
			}
		}

		results = append(results, result)
	}

	return results
}

func isFailureProcessed(f string) bool {
	for _, p := range processedFailures {
		if f == p {
			return true
		}
	}

	return false
}

func (p *Parser) getAttachment(a Attachment) (models.Attachment, error) {
	a.Filename.Value = strings.Replace(a.Filename.Value, heicExt, ".jpeg", 1)

	att, err := p.readAttachment(a.PayloadRef.ID.Value)
	if err != nil {
		return models.Attachment{}, fmt.Errorf("failed to get attachment: %w", err)
	}

	return models.Attachment{
		Name:    a.Filename.Value,
		Content: &att,
	}, nil
}

func (p *Parser) getSteps(as ActivitySummaries, level int) ([]models.Step, bool, []models.Attachment) {
	const op = "xctest.Parser.getSteps"
	logger := slog.With("op", op)

	steps := make([]models.Step, 0, len(as.Values))
	attachments := make([]models.Attachment, 0)
	isFailedStep := false

	for _, v := range as.Values {
		stepAttachments := make([]models.Attachment, 0)
		stepComment := ""
		isStepFailed := false
		childSteps := make([]models.Step, 0)
		if v.ActivityType.Value == attachmentStep && v.Attachments != nil &&
			len(v.Attachments.Values) == 1 && v.Attachments.Values[0].Name.Value == qaseConfig {
			att, err := p.getAttachment(v.Attachments.Values[0])
			if err != nil {
				logger.Error("failed to get attachments", "error", err)
			}

			var ID QaseId
			err = json.Unmarshal(*att.Content, &ID)
			if err != nil {
				logger.Error("failed to unmarshal Qase ID", "error", err)
			}

			if caseId == nil {
				caseId = &ID.ID
			}

			continue
		}

		if v.FailureSummaryIDs != nil {
			for _, f := range v.FailureSummaryIDs.Values {
				if fm, ok := failures[f.Value]; ok {
					processedFailures = append(processedFailures, f.Value)
					isStepFailed = true
					stepComment = fmt.Sprintf("Line: %s\n%s", fm.LineNumber.Value, fm.Message.Value)
					if fm.Attachments != nil {
						for _, a := range fm.Attachments.Values {
							att, err := p.getAttachment(a)
							if err != nil {
								logger.Error("failed to get attachments", "error", err)
								continue
							}

							stepAttachments = append(stepAttachments, att)
						}
					}
				}
			}

			if v.Attachments != nil {
				for _, a := range v.Attachments.Values {
					att, err := p.getAttachment(a)
					if err != nil {
						logger.Error("failed to get attachments", "error", err)
						continue
					}

					stepAttachments = append(stepAttachments, att)
				}
			}
		}

		if v.Subactivities != nil {
			cs, f, a := p.getSteps(*v.Subactivities, level+1)
			if f {
				isStepFailed = true
			}
			childSteps = append(childSteps, cs...)

			if len(cs) == 0 {
				stepAttachments = append(stepAttachments, a...)
			}
		}

		isActivityTypeInternalOrDelete := v.ActivityType.Value == internalStep || v.ActivityType.Value == deleteStep
		isLevelNotAll := p.level != All
		isLevelZeroOrAboveFirst := level == 0 || (level > 0 && p.level != FirstLevel)

		if isActivityTypeInternalOrDelete && isLevelNotAll && isLevelZeroOrAboveFirst {
			attachments = append(attachments, stepAttachments...)

			if isStepFailed {
				isFailedStep = true
			}

			continue
		}

		logger := logger.With("stepAction", v.Title.Value)
		logger.Debug("processing step", "step", v)

		step := models.Step{
			Data: models.Data{
				Action: v.Title.Value,
			},
			Execution: models.StepExecution{
				Status:      "passed",
				Attachments: stepAttachments,
				Comment:     stepComment,
			},
			Steps: childSteps,
		}

		if v.Start.Value != "" {
			st := float64(parseTime(v.Start.Value).Unix())
			step.Execution.StartTime = &st
		}

		if v.Finish.Value != "" {
			et := float64(parseTime(v.Finish.Value).Unix())
			step.Execution.EndTime = &et
		}

		if isStepFailed {
			step.Execution.Status = "failed"
			isFailedStep = true
		}

		steps = append(steps, step)
	}

	return steps, isFailedStep, attachments
}

func getStatus(s string) string {
	switch s {
	case "Success":
		return "passed"
	case "Failure":
		return "failed"
	case "Error":
		return "invalid"
	case "Skipped":
		return "skipped"
	default:
		return "passed"
	}
}

func parseTime(s string) time.Time {
	t, err := time.Parse(layoutTime, s)
	if err != nil {
		slog.Error("failed to parse time", "error", err, "time", s)
		return time.Time{}
	}

	return t
}
