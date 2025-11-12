package allure

import (
	"reflect"
	"testing"

	models "github.com/qase-tms/qasectl/internal/models/result"
)

func TestParser_extractTestOpsID(t *testing.T) {
	parser := &Parser{}

	tests := []struct {
		name     string
		input    string
		expected *int64
	}{
		{name: "empty string", input: "", expected: nil},
		{name: "no separator", input: "test", expected: nil},
		{name: "invalid ID", input: "test-case-abc", expected: nil},
		{name: "valid case", input: "test-case-12345", expected: func() *int64 { i := int64(12345); return &i }()},
		{name: "empty ID", input: "test-case-", expected: nil},
		{name: "ID at the start", input: "12345-case", expected: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.extractTestOpsID(tt.input)

			if (result == nil && tt.expected != nil) || (result != nil && tt.expected == nil) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			} else if result != nil && *result != *tt.expected {
				t.Errorf("expected %v, got %v", *tt.expected, *result)
			}
		})
	}
}

func TestParser_convertTest_LayerField(t *testing.T) {
	parser := &Parser{}

	tests := []struct {
		name           string
		labelName      string
		labelValue     string
		expectedField  string
		expectedValue  string
	}{
		{name: "valid layer unknown", labelName: "layer", labelValue: "unknown", expectedField: "layer", expectedValue: "unknown"},
		{name: "valid layer e2e", labelName: "layer", labelValue: "e2e", expectedField: "layer", expectedValue: "e2e"},
		{name: "valid layer api", labelName: "layer", labelValue: "api", expectedField: "layer", expectedValue: "api"},
		{name: "valid layer unit", labelName: "layer", labelValue: "unit", expectedField: "layer", expectedValue: "unit"},
		{name: "invalid layer custom", labelName: "layer", labelValue: "custom", expectedField: "custom layer", expectedValue: "custom"},
		{name: "invalid layer integration", labelName: "layer", labelValue: "integration", expectedField: "custom layer", expectedValue: "integration"},
		{name: "invalid layer empty", labelName: "layer", labelValue: "", expectedField: "custom layer", expectedValue: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := Test{
				Name:   "Test",
				Start:  1000,
				Stop:   2000,
				Status: "passed",
				Labels: []Label{
					{Name: tt.labelName, Value: tt.labelValue},
				},
			}

			result := parser.convertTest(test)

			if result.Fields[tt.expectedField] != tt.expectedValue {
				t.Errorf("expected field %s to have value %s, got %s", tt.expectedField, tt.expectedValue, result.Fields[tt.expectedField])
			}

			// Check that layer field doesn't exist if it was renamed
			if tt.expectedField == "custom layer" {
				if _, exists := result.Fields["layer"]; exists {
					t.Errorf("expected layer field to be renamed to 'custom layer', but 'layer' field still exists with value %s", result.Fields["layer"])
				}
			}
		})
	}
}

func TestParser_convertTest_SuiteField(t *testing.T) {
	parser := &Parser{}

	test := Test{
		Name:   "Test",
		Start:  1000,
		Stop:   2000,
		Status: "passed",
		Labels: []Label{
			{Name: "suite", Value: "com.example.test"},
		},
	}

	result := parser.convertTest(test)

	expectedSuites := []models.SuiteData{
		{Title: "com"},
		{Title: "example"},
		{Title: "test"},
	}

	if !reflect.DeepEqual(result.Relations.Suite.Data, expectedSuites) {
		t.Errorf("expected suites %v, got %v", expectedSuites, result.Relations.Suite.Data)
	}
}

func TestParser_convertTest_PackageField(t *testing.T) {
	parser := &Parser{}

	test := Test{
		Name:   "Test",
		Start:  1000,
		Stop:   2000,
		Status: "passed",
		Labels: []Label{
			{Name: "package", Value: "com.example.test"},
		},
	}

	result := parser.convertTest(test)

	expectedSuites := []models.SuiteData{
		{Title: "com"},
		{Title: "example"},
		{Title: "test"},
	}

	if !reflect.DeepEqual(result.Relations.Suite.Data, expectedSuites) {
		t.Errorf("expected suites %v, got %v", expectedSuites, result.Relations.Suite.Data)
	}
}

func TestParser_convertTest_SuiteOverridesPackage(t *testing.T) {
	parser := &Parser{}

	test := Test{
		Name:   "Test",
		Start:  1000,
		Stop:   2000,
		Status: "passed",
		Labels: []Label{
			{Name: "package", Value: "com.example.old"},
			{Name: "suite", Value: "com.example.new"},
		},
	}

	result := parser.convertTest(test)

	expectedSuites := []models.SuiteData{
		{Title: "com"},
		{Title: "example"},
		{Title: "new"},
	}

	if !reflect.DeepEqual(result.Relations.Suite.Data, expectedSuites) {
		t.Errorf("expected suites %v, got %v", expectedSuites, result.Relations.Suite.Data)
	}
}
