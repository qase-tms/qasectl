package qase

import (
	"encoding/json"
	"fmt"
	models "github.com/qase-tms/qasectl/internal/models/result"
	"log/slog"
	"os"
	"path"
	"path/filepath"
)

// Parser is a parser for Qase files
type Parser struct {
	path string
}

// NewParser creates a new Parser
func NewParser(path string) *Parser {
	return &Parser{
		path: path,
	}
}

// Parse parses the Qase file and returns the results
func (p *Parser) Parse() ([]models.Result, error) {
	const op = "qase.Parser.Parse"
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

// parseFile parses a single Qase file
func (p *Parser) parseFile(path string) (models.Result, error) {
	byteValue, err := os.ReadFile(path)
	if err != nil {
		return models.Result{}, fmt.Errorf("failed to read file: %w", err)
	}

	var result models.Result
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		return result, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	result.Attachments = p.convertAttachments(result.Attachments)

	result.Steps = p.convertStepAttachments(result.Steps)

	result = p.calculateDuration(result)

	return result, nil
}

func (p *Parser) convertStepAttachments(steps []models.Step) []models.Step {
	for i := range steps {
		steps[i].Execution.Attachments = p.convertAttachments(steps[i].Execution.Attachments)

		if steps[i].Steps != nil {
			steps[i].Steps = p.convertStepAttachments(steps[i].Steps)
		}
	}

	return steps
}

func (p *Parser) convertAttachments(attachments []models.Attachment) []models.Attachment {
	for i := range attachments {
		if attachments[i].FilePath == nil {
			continue
		}

		_, err := os.Stat(*attachments[i].FilePath)
		if err == nil {
			continue
		}

		if os.IsNotExist(err) {
			id := filepath.Base(*attachments[i].FilePath)
			dir := filepath.Dir(p.path)
			*attachments[i].FilePath = path.Join(dir, "attachments", id)
		}
	}

	return attachments
}

func (p *Parser) calculateDuration(result models.Result) models.Result {
	if result.StartTime != nil && result.EndTime != nil {
		duration := *result.Execution.EndTime - *result.Execution.StartTime
		result.Execution.Duration = &duration
	}

	result.Execution.StartTime = nil
	result.Execution.EndTime = nil

	return result
}
