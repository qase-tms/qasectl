package junit

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	models "github.com/qase-tms/qasectl/internal/models/result"
)

func TestParseSteps(t *testing.T) {
	tests := []struct {
		name       string
		properties Properties
		expected   []models.Step
	}{
		{
			name:       "Empty properties",
			properties: Properties{},
			expected:   []models.Step{},
		},
		{
			name: "Single step",
			properties: Properties{
				Property: []Property{
					{Name: "step[passed]", Value: "A"},
				},
			},
			expected: []models.Step{
				{
					Data: models.Data{Action: "A"},
					Execution: models.StepExecution{
						Status: "passed",
					},
					Steps: nil,
				},
			},
		},
		{
			name: "Multiple steps",
			properties: Properties{
				Property: []Property{
					{Name: "step[passed]", Value: "A"},
					{Name: "step[failed]", Value: "B"},
				},
			},
			expected: []models.Step{
				{
					Data: models.Data{Action: "A"},
					Execution: models.StepExecution{
						Status: "passed",
					},
				},
				{
					Data: models.Data{Action: "B"},
					Execution: models.StepExecution{
						Status: "failed",
					},
				},
			},
		},
		{
			name: "Nested steps",
			properties: Properties{
				Property: []Property{
					{Name: "step[passed]", Value: "C/D"},
					{Name: "step[passed]", Value: "C"},
				},
			},
			expected: []models.Step{
				{
					Data: models.Data{Action: "C"},
					Execution: models.StepExecution{
						Status: "passed",
					},
					Steps: []models.Step{
						{
							Data: models.Data{Action: "D"},
							Execution: models.StepExecution{
								Status: "passed",
							},
						},
					},
				},
			},
		},
		{
			name: "Duplicate steps",
			properties: Properties{
				Property: []Property{
					{Name: "step[passed]", Value: "A"},
					{Name: "step[passed]", Value: "A"},
				},
			},
			expected: []models.Step{
				{
					Data: models.Data{Action: "A"},
					Execution: models.StepExecution{
						Status: "passed",
					},
				},
				{
					Data: models.Data{Action: "A"},
					Execution: models.StepExecution{
						Status: "passed",
					},
				},
			},
		},
		{
			name: "Invalid property name",
			properties: Properties{
				Property: []Property{
					{Name: "invalid[step]", Value: "A"},
				},
			},
			expected: []models.Step{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseSteps(tt.properties)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %+v, want %+v", result, tt.expected)
			}
		})
	}
}

func TestParser_parseFile(t *testing.T) {
	tests := []struct {
		name     string
		xmlData  string
		expected []models.Result
		wantErr  bool
	}{
		{
			name: "Parse TestSuites format",
			xmlData: `<?xml version="1.0" encoding="UTF-8"?>
<testsuites name="Test Suite">
  <testsuite name="TestSuite1" tests="2" failures="1" errors="0" skipped="0" time="0.1">
    <testcase name="Test1" classname="TestClass" time="0.05">
      <failure message="Test failed">Failure details</failure>
    </testcase>
    <testcase name="Test2" classname="TestClass" time="0.05">
    </testcase>
  </testsuite>
</testsuites>`,
			expected: []models.Result{
				{
					Title: "Test1",
					Relations: models.Relation{
						Suite: models.Suite{
							Data: []models.SuiteData{
								{Title: "Test Suite"},
								{Title: "TestSuite1"},
							},
						},
					},
					Execution: models.Execution{
						Status:     "failed",
						Duration:   float64Ptr(50.0),
						StackTrace: stringPtr("Failure details"),
					},
					Message: stringPtr("Test failed"),
				},
				{
					Title: "Test2",
					Relations: models.Relation{
						Suite: models.Suite{
							Data: []models.SuiteData{
								{Title: "Test Suite"},
								{Title: "TestSuite1"},
							},
						},
					},
					Execution: models.Execution{
						Status:   "passed",
						Duration: float64Ptr(50.0),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Parse TestSuite format",
			xmlData: `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="TestSuite1" tests="2" failures="1" errors="0" skipped="0" time="0.1">
  <testcase name="Test1" classname="TestClass" time="0.05">
    <failure message="Test failed">Failure details</failure>
  </testcase>
  <testcase name="Test2" classname="TestClass" time="0.05">
  </testcase>
</testsuite>`,
			expected: []models.Result{
				{
					Title: "Test1",
					Relations: models.Relation{
						Suite: models.Suite{
							Data: []models.SuiteData{
								{Title: "TestSuite1"},
							},
						},
					},
					Execution: models.Execution{
						Status:     "failed",
						Duration:   float64Ptr(50.0),
						StackTrace: stringPtr("Failure details"),
					},
					Message: stringPtr("Test failed"),
				},
				{
					Title: "Test2",
					Relations: models.Relation{
						Suite: models.Suite{
							Data: []models.SuiteData{
								{Title: "TestSuite1"},
							},
						},
					},
					Execution: models.Execution{
						Status:   "passed",
						Duration: float64Ptr(50.0),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Parse invalid XML",
			xmlData: `<?xml version="1.0" encoding="UTF-8"?>
<invalid>This is not a valid JUnit XML</invalid>`,
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test.xml")
			err := os.WriteFile(tmpFile, []byte(tt.xmlData), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			parser := NewParser(tmpFile)
			result, err := parser.parseFile(tmpFile)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseFile() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("parseFile() unexpected error: %v", err)
				return
			}

			// Compare results (simplified comparison for test purposes)
			if len(result) != len(tt.expected) {
				t.Errorf("parseFile() got %d results, want %d", len(result), len(tt.expected))
				return
			}

			for i, res := range result {
				expected := tt.expected[i]
				if res.Title != expected.Title {
					t.Errorf("Result[%d].Title = %v, want %v", i, res.Title, expected.Title)
				}
				if res.Execution.Status != expected.Execution.Status {
					t.Errorf("Result[%d].Execution.Status = %v, want %v", i, res.Execution.Status, expected.Execution.Status)
				}
				if expected.Execution.Duration != nil {
					if res.Execution.Duration == nil {
						t.Errorf("Result[%d].Execution.Duration = nil, want %v", i, *expected.Execution.Duration)
					} else if *res.Execution.Duration != *expected.Execution.Duration {
						t.Errorf("Result[%d].Execution.Duration = %v, want %v", i, *res.Execution.Duration, *expected.Execution.Duration)
					}
				}
			}
		})
	}
}

// Helper functions for creating pointers
func float64Ptr(v float64) *float64 {
	return &v
}

func stringPtr(v string) *string {
	return &v
}

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected int
		wantErr  bool
	}{
		{
			name: "Parse directory with multiple files",
			files: map[string]string{
				"test1.xml": `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="TestSuite1" tests="1" failures="0" errors="0" skipped="0" time="0.1">
  <testcase name="Test1" classname="TestClass" time="0.05">
  </testcase>
</testsuite>`,
				"test2.xml": `<?xml version="1.0" encoding="UTF-8"?>
<testsuites name="Test Suite">
  <testsuite name="TestSuite2" tests="1" failures="0" errors="0" skipped="0" time="0.1">
    <testcase name="Test2" classname="TestClass" time="0.05">
    </testcase>
  </testsuite>
</testsuites>`,
			},
			expected: 2,
			wantErr:  false,
		},
		{
			name: "Parse single file",
			files: map[string]string{
				"test.xml": `<?xml version="1.0" encoding="UTF-8"?>
<testsuite name="TestSuite1" tests="1" failures="0" errors="0" skipped="0" time="0.1">
  <testcase name="Test1" classname="TestClass" time="0.05">
  </testcase>
</testsuite>`,
			},
			expected: 1,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tmpDir := t.TempDir()

			// Create test files
			for filename, content := range tt.files {
				filePath := filepath.Join(tmpDir, filename)
				err := os.WriteFile(filePath, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file %s: %v", filename, err)
				}
			}

			parser := NewParser(tmpDir)
			result, err := parser.Parse()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
				return
			}

			if len(result) != tt.expected {
				t.Errorf("Parse() got %d results, want %d", len(result), tt.expected)
			}
		})
	}
}
