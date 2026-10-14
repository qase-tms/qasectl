package result

import (
	"testing"

	models "github.com/qase-tms/qasectl/internal/models/result"
)

// resultWithSteps returns a result whose step tree holds n nodes: one root with n-1 children.
func resultWithSteps(title string, n int) models.Result {
	r := models.Result{Title: title}
	if n == 0 {
		return r
	}
	children := make([]models.Step, n-1)
	r.Steps = []models.Step{{Steps: children}}
	return r
}

func batchTitles(batches [][]models.Result) [][]string {
	out := make([][]string, 0, len(batches))
	for _, b := range batches {
		titles := make([]string, 0, len(b))
		for _, r := range b {
			titles = append(titles, r.Title)
		}
		out = append(out, titles)
	}
	return out
}

func TestSplitBatches(t *testing.T) {
	tests := []struct {
		name      string
		steps     []int
		batchSize int64
		maxSteps  int
		want      [][]string
	}{
		{name: "by result count", steps: []int{1, 1, 1, 1, 1}, batchSize: 2, maxSteps: 100, want: [][]string{{"r0", "r1"}, {"r2", "r3"}, {"r4"}}},
		{name: "by step count", steps: []int{40, 40, 40, 20}, batchSize: 200, maxSteps: 100, want: [][]string{{"r0", "r1"}, {"r2", "r3"}}},
		{name: "exactly at the step limit", steps: []int{50, 50, 1}, batchSize: 200, maxSteps: 100, want: [][]string{{"r0", "r1"}, {"r2"}}},
		{name: "oversized result goes alone", steps: []int{10, 150, 10}, batchSize: 200, maxSteps: 100, want: [][]string{{"r0"}, {"r1"}, {"r2"}}},
		{name: "results without steps", steps: []int{0, 0, 0}, batchSize: 200, maxSteps: 100, want: [][]string{{"r0", "r1", "r2"}}},
		{name: "empty", steps: nil, batchSize: 200, maxSteps: 100, want: [][]string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := make([]models.Result, 0, len(tt.steps))
			for i, n := range tt.steps {
				results = append(results, resultWithSteps("r"+string(rune('0'+i)), n))
			}

			got := batchTitles(splitBatches(results, tt.batchSize, tt.maxSteps))
			if len(got) != len(tt.want) {
				t.Fatalf("splitBatches() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if len(got[i]) != len(tt.want[i]) {
					t.Fatalf("splitBatches() = %v, want %v", got, tt.want)
				}
				for j := range got[i] {
					if got[i][j] != tt.want[i][j] {
						t.Fatalf("splitBatches() = %v, want %v", got, tt.want)
					}
				}
			}
		})
	}
}

func TestCountSteps_CountsNestedSteps(t *testing.T) {
	steps := []models.Step{
		{Steps: []models.Step{{}, {Steps: []models.Step{{}}}}},
		{},
	}
	if got := countSteps(steps); got != 5 {
		t.Errorf("countSteps() = %d, want 5", got)
	}
}
