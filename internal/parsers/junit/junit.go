package junit

import (
	"encoding/xml"
	"fmt"
	models "github.com/qase-tms/qasectl/internal/models/result"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Parser is a parser for Junit XML files
type Parser struct {
	path string
}

// NewParser creates a new Parser
func NewParser(path string) *Parser {
	return &Parser{
		path: path,
	}
}

// Parse parses the Junit XML file and returns the results
func (p *Parser) Parse() ([]models.Result, error) {
	var results []models.Result

	fileInfo, err := os.Stat(p.path)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		err := filepath.Walk(p.path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				result, err := p.parseFile(path)
				if err != nil {
					return err
				}
				results = append(results, result...)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		result, err := p.parseFile(p.path)
		if err != nil {
			return nil, err
		}
		results = append(results, result...)
	}

	return results, nil
}

// parseFile parses a single Junit XML file
func (p *Parser) parseFile(path string) ([]models.Result, error) {
	xmlFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := xmlFile.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	byteValue, _ := io.ReadAll(xmlFile)

	var testSuites TestSuites
	err = xml.Unmarshal(byteValue, &testSuites)
	if err != nil {
		return nil, err
	}

	return convertTestSuites(testSuites), nil
}

// convertTestSuites converts a TestSuites to Results
func convertTestSuites(testSuites TestSuites) []models.Result {
	results := make([]models.Result, 0)

	for _, testSuite := range testSuites.TestSuites {
		for _, testCase := range testSuite.TestCases {
			relation := models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{
						{
							Title: testSuites.Name,
						},
						{
							Title: testSuite.Name,
						},
					},
				},
			}

			signature := fmt.Sprintf("%s::%s::%s::%s", testSuites.Name, testSuite.Name, testCase.ClassName, testCase.Name)

			status := "passed"
			var stackTrace *string
			var message *string

			if testCase.Failure != nil {
				status = "failed"
				stackTrace = &testCase.Failure.Body
				message = &testCase.Failure.Message
			}

			if testCase.Error != nil {
				status = "invalid"
				stackTrace = &testCase.Error.Body
				message = &testCase.Error.Message
			}

			if testCase.Skipped != nil {
				status = "skipped"
				message = &testCase.Skipped.Message
			}

			fields := make(map[string]string)

			for k := range testCase.Properties.Property {
				fields[testCase.Properties.Property[k].Name] = testCase.Properties.Property[k].Value
			}

			result := models.Result{
				Title:     testCase.Name,
				Signature: &signature,
				Relations: relation,
				Execution: models.Execution{
					Duration:   &testCase.Time,
					Status:     status,
					StackTrace: stackTrace,
				},
				Attachments: make([]models.Attachment, 0),
				Steps:       make([]models.Step, 0),
				StepType:    "text",
				Params:      make(map[string]string),
				Muted:       false,
				Fields:      fields,
				Message:     message,
			}
			results = append(results, result)
		}
	}

	return results
}
