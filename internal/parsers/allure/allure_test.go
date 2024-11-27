package allure

import (
	"testing"
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
