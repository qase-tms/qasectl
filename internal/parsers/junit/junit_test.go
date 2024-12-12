package junit

import (
	models "github.com/qase-tms/qasectl/internal/models/result"
	"reflect"
	"testing"
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
