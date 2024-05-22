package allure

import (
	"encoding/json"
	"fmt"
	models "github.com/qase-tms/qasectl/internal/models/result"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Parser is a parser for Allure files
type Parser struct {
	path string
}

// NewParser creates a new Parser
func NewParser(path string) *Parser {
	return &Parser{
		path: path,
	}
}

// Parse parses the Allure file and returns the results
func (p *Parser) Parse() ([]models.Result, error) {
	const op = "allure.Parser.Parse"
	logger := slog.With("path", p.path, "op", op)

	var files []string
	fileInfo, err := os.Stat(p.path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	if fileInfo.IsDir() {
		err := filepath.Walk(p.path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("failed to walk path: %w", err)
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to walk path: %w", err)
		}
	} else {
		files = append(files, p.path)
	}

	if len(files) == 0 {
		logger.Info("no files found")
		return nil, nil
	}

	results := make([]models.Result, 0, len(files))

	for _, file := range files {
		if filepath.Ext(file) != ".json" {
			logger.Debug("skipping file. Only support json format", "file", file)
			continue
		}

		logger.Debug("parsing file", "file", file)

		result, err := p.parseFile(file)
		if err != nil {
			logger.Error("failed to parse file", "file", file, "error", err)
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

func (p *Parser) parseFile(file string) (models.Result, error) {
	byteValue, err := os.ReadFile(file)
	if err != nil {
		return models.Result{}, fmt.Errorf("failed to read file: %w", err)
	}

	var test Test
	err = json.Unmarshal(byteValue, &test)
	if err != nil {
		return models.Result{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return p.convertTest(test), nil
}

func (p *Parser) convertTest(test Test) models.Result {
	d := test.Stop - test.Start
	result := models.Result{
		Title:       test.Name,
		Steps:       make([]models.Step, 0, len(test.Steps)),
		Attachments: p.convertAttachments(test.Attachments),
		StepType:    "text",
		Duration:    &d,
		Execution: models.Execution{
			Duration:   &d,
			Status:     test.Status,
			StackTrace: test.StatusDetails.Trace,
		},
		Message: test.StatusDetails.Message,
	}

	for _, step := range test.Steps {
		result.Steps = append(result.Steps, p.convertStep(step))
	}

	for _, v := range test.Labels {
		if v.Value == "thread" {
			result.Execution.Thread = &v.Value
		}

		if v.Name != "package" {
			continue
		}

		suites := strings.Split(v.Value, ".")
		data := make([]models.SuiteData, 0, len(suites))

		for i := range suites {
			data = append(data, models.SuiteData{Title: suites[i]})
		}

		result.Relations = models.Relation{
			Suite: models.Suite{
				Data: data,
			},
		}
	}

	return result
}

func (p *Parser) convertStep(step TestStep) models.Step {
	d := step.Stop - step.Start
	result := models.Step{
		Data: models.Data{
			Action: step.Name,
		},
		Execution: models.StepExecution{
			Attachments: p.convertAttachments(step.Attachments),
			Duration:    &d,
			Status:      step.Status,
		},
		Steps: make([]models.Step, 0, len(step.Steps)),
	}

	for _, s := range step.Steps {
		result.Steps = append(result.Steps, p.convertStep(s))
	}

	return result
}

func (p *Parser) convertAttachments(attachments []Attachment) []models.Attachment {
	result := make([]models.Attachment, 0, len(attachments))

	for _, attachment := range attachments {
		d := filepath.Dir(p.path)
		p := path.Join(d, attachment.Source)
		result = append(result, models.Attachment{
			Name:     attachment.Name,
			FilePath: &p,
		})
	}

	return result
}
