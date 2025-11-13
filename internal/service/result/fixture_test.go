package result

import (
	"strings"
	"testing"

	models "github.com/qase-tms/qasectl/internal/models/result"
	"github.com/qase-tms/qasectl/internal/service/result/mocks"
	"go.uber.org/mock/gomock"
)

type fixture struct {
	client *mocks.Mockclient
	rs     *mocks.MockrunService
	parser *mocks.MockParser
}

func newFixture(t *testing.T) *fixture {
	ctr := gomock.NewController(t)

	return &fixture{
		client: mocks.NewMockclient(ctr),
		rs:     mocks.NewMockrunService(ctr),
		parser: mocks.NewMockParser(ctr),
	}
}

func prepareModels() []models.Result {
	return []models.Result{
		{
			ID:        nil,
			Title:     "Test 1",
			Signature: nil,
			TestOpsID: nil,
			Execution: models.Execution{
				StartTime:  nil,
				EndTime:    nil,
				Status:     "passed",
				Duration:   nil,
				StackTrace: nil,
				Thread:     nil,
			},
			Fields:      make(map[string]string),
			Attachments: []models.Attachment{},
			Steps:       []models.Step{},
			StepType:    "text",
			Params: map[string]string{
				"browser": "chrome",
				"version": "1.0.0",
			},
			Relations: models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{},
				},
			},
			Muted:   false,
			Message: nil,
		},
		{
			ID:        nil,
			Title:     "Test 2",
			Signature: nil,
			TestOpsID: nil,
			Execution: models.Execution{
				StartTime:  nil,
				EndTime:    nil,
				Status:     "failed",
				Duration:   nil,
				StackTrace: nil,
				Thread:     nil,
			},
			Fields:      make(map[string]string),
			Attachments: []models.Attachment{},
			Steps:       []models.Step{},
			StepType:    "text",
			Params: map[string]string{
				"browser": "firefox",
				"version": "1.0.0",
			},
			Relations: models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{},
				},
			},
			Muted:   false,
			Message: nil,
		},
	}
}

func prepareModelsWithEmptyParams() []models.Result {
	return []models.Result{
		{
			ID:        nil,
			Title:     "Test with Empty Params",
			Signature: nil,
			TestOpsID: nil,
			Execution: models.Execution{
				StartTime:  nil,
				EndTime:    nil,
				Status:     "passed",
				Duration:   nil,
				StackTrace: nil,
				Thread:     nil,
			},
			Fields:      make(map[string]string),
			Attachments: []models.Attachment{},
			Steps:       []models.Step{},
			StepType:    "text",
			Params:      nil, // Empty params for SkipParams testing
			Relations: models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{},
				},
			},
			Muted:   false,
			Message: nil,
		},
	}
}

func prepareModelsWithAttachments() []models.Result {
	return []models.Result{
		{
			ID:        nil,
			Title:     "Test with Attachments",
			Signature: nil,
			TestOpsID: nil,
			Execution: models.Execution{
				StartTime:  nil,
				EndTime:    nil,
				Status:     "passed",
				Duration:   nil,
				StackTrace: nil,
				Thread:     nil,
			},
			Fields: make(map[string]string),
			Attachments: []models.Attachment{
				{Name: "screenshot.png"},
				{Name: "photo.jpg"},
				{Name: "document.pdf"},
				{Name: "log.txt"},
			},
			Steps:    []models.Step{},
			StepType: "text",
			Params: map[string]string{
				"browser": "chrome",
				"version": "1.0.0",
			},
			Relations: models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{},
				},
			},
			Muted:   false,
			Message: nil,
		},
	}
}

func prepareModelsWithStepAttachments() []models.Result {
	return []models.Result{
		{
			ID:        nil,
			Title:     "Test with Step Attachments",
			Signature: nil,
			TestOpsID: nil,
			Execution: models.Execution{
				StartTime:  nil,
				EndTime:    nil,
				Status:     "passed",
				Duration:   nil,
				StackTrace: nil,
				Thread:     nil,
			},
			Fields:      make(map[string]string),
			Attachments: []models.Attachment{},
			Steps: []models.Step{
				{
					Data: models.Data{
						Action: "Step 1",
					},
					Execution: models.StepExecution{
						Attachments: []models.Attachment{
							{Name: "step1.png"},
							{Name: "step1.pdf"},
							{Name: "step1.jpg"},
						},
						Duration: nil,
						Status:   "passed",
					},
					Steps: []models.Step{},
				},
				{
					Data: models.Data{
						Action: "Step 2",
					},
					Execution: models.StepExecution{
						Attachments: []models.Attachment{
							{Name: "step2.pdf"},
							{Name: "step2.txt"},
						},
						Duration: nil,
						Status:   "passed",
					},
					Steps: []models.Step{},
				},
			},
			StepType: "text",
			Params: map[string]string{
				"browser": "chrome",
				"version": "1.0.0",
			},
			Relations: models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{},
				},
			},
			Muted:   false,
			Message: nil,
		},
	}
}

// prepareModelsWithStartTime creates models with different StartTime values for sorting tests
func prepareModelsWithStartTime() []models.Result {
	startTime1 := float64(1000)
	endTime1 := float64(2000)
	startTime2 := float64(500)
	endTime2 := float64(1500)
	startTime3 := float64(2000)
	endTime3 := float64(3000)

	return []models.Result{
		{
			ID:        nil,
			Title:     "Test 3",
			Signature: nil,
			TestOpsID: nil,
			Execution: models.Execution{
				StartTime:  &startTime3,
				EndTime:    &endTime3,
				Status:     "passed",
				Duration:   nil,
				StackTrace: nil,
				Thread:     nil,
			},
			Fields:      make(map[string]string),
			Attachments: []models.Attachment{},
			Steps:       []models.Step{},
			StepType:    "text",
			Params:      make(map[string]string),
			Relations: models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{},
				},
			},
			Muted:   false,
			Message: nil,
		},
		{
			ID:        nil,
			Title:     "Test 1",
			Signature: nil,
			TestOpsID: nil,
			Execution: models.Execution{
				StartTime:  &startTime1,
				EndTime:    &endTime1,
				Status:     "passed",
				Duration:   nil,
				StackTrace: nil,
				Thread:     nil,
			},
			Fields:      make(map[string]string),
			Attachments: []models.Attachment{},
			Steps:       []models.Step{},
			StepType:    "text",
			Params:      make(map[string]string),
			Relations: models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{},
				},
			},
			Muted:   false,
			Message: nil,
		},
		{
			ID:        nil,
			Title:     "Test 2",
			Signature: nil,
			TestOpsID: nil,
			Execution: models.Execution{
				StartTime:  &startTime2,
				EndTime:    &endTime2,
				Status:     "failed",
				Duration:   nil,
				StackTrace: nil,
				Thread:     nil,
			},
			Fields:      make(map[string]string),
			Attachments: []models.Attachment{},
			Steps:       []models.Step{},
			StepType:    "text",
			Params:      make(map[string]string),
			Relations: models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{},
				},
			},
			Muted:   false,
			Message: nil,
		},
	}
}

// prepareModelsWithNilStartTime creates models with some nil StartTime values for nil handling tests
func prepareModelsWithNilStartTime() []models.Result {
	startTime1 := float64(1000)
	endTime1 := float64(2000)
	startTime2 := float64(500)
	endTime2 := float64(1500)

	return []models.Result{
		{
			ID:        nil,
			Title:     "Test with nil StartTime",
			Signature: nil,
			TestOpsID: nil,
			Execution: models.Execution{
				StartTime:  nil,
				EndTime:    nil,
				Status:     "passed",
				Duration:   nil,
				StackTrace: nil,
				Thread:     nil,
			},
			Fields:      make(map[string]string),
			Attachments: []models.Attachment{},
			Steps:       []models.Step{},
			StepType:    "text",
			Params:      make(map[string]string),
			Relations: models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{},
				},
			},
			Muted:   false,
			Message: nil,
		},
		{
			ID:        nil,
			Title:     "Test 1",
			Signature: nil,
			TestOpsID: nil,
			Execution: models.Execution{
				StartTime:  &startTime1,
				EndTime:    &endTime1,
				Status:     "passed",
				Duration:   nil,
				StackTrace: nil,
				Thread:     nil,
			},
			Fields:      make(map[string]string),
			Attachments: []models.Attachment{},
			Steps:       []models.Step{},
			StepType:    "text",
			Params:      make(map[string]string),
			Relations: models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{},
				},
			},
			Muted:   false,
			Message: nil,
		},
		{
			ID:        nil,
			Title:     "Test 2",
			Signature: nil,
			TestOpsID: nil,
			Execution: models.Execution{
				StartTime:  &startTime2,
				EndTime:    &endTime2,
				Status:     "failed",
				Duration:   nil,
				StackTrace: nil,
				Thread:     nil,
			},
			Fields:      make(map[string]string),
			Attachments: []models.Attachment{},
			Steps:       []models.Step{},
			StepType:    "text",
			Params:      make(map[string]string),
			Relations: models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{},
				},
			},
			Muted:   false,
			Message: nil,
		},
	}
}

// prepareModelsWithLongTitles creates models with titles longer than 255 characters for truncation tests
func prepareModelsWithLongTitles() []models.Result {
	// Create a title longer than 255 characters
	longTitle := strings.Repeat("A", 300) // 300 characters
	exactTitle := strings.Repeat("B", 255) // Exactly 255 characters
	shortTitle := strings.Repeat("C", 100)  // Less than 255 characters

	return []models.Result{
		{
			ID:        nil,
			Title:     longTitle,
			Signature: nil,
			TestOpsID: nil,
			Execution: models.Execution{
				StartTime:  nil,
				EndTime:    nil,
				Status:     "passed",
				Duration:   nil,
				StackTrace: nil,
				Thread:     nil,
			},
			Fields:      make(map[string]string),
			Attachments: []models.Attachment{},
			Steps:       []models.Step{},
			StepType:    "text",
			Params:      make(map[string]string),
			Relations: models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{},
				},
			},
			Muted:   false,
			Message: nil,
		},
		{
			ID:        nil,
			Title:     exactTitle,
			Signature: nil,
			TestOpsID: nil,
			Execution: models.Execution{
				StartTime:  nil,
				EndTime:    nil,
				Status:     "passed",
				Duration:   nil,
				StackTrace: nil,
				Thread:     nil,
			},
			Fields:      make(map[string]string),
			Attachments: []models.Attachment{},
			Steps:       []models.Step{},
			StepType:    "text",
			Params:      make(map[string]string),
			Relations: models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{},
				},
			},
			Muted:   false,
			Message: nil,
		},
		{
			ID:        nil,
			Title:     shortTitle,
			Signature: nil,
			TestOpsID: nil,
			Execution: models.Execution{
				StartTime:  nil,
				EndTime:    nil,
				Status:     "failed",
				Duration:   nil,
				StackTrace: nil,
				Thread:     nil,
			},
			Fields:      make(map[string]string),
			Attachments: []models.Attachment{},
			Steps:       []models.Step{},
			StepType:    "text",
			Params:      make(map[string]string),
			Relations: models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{},
				},
			},
			Muted:   false,
			Message: nil,
		},
	}
}

// prepareModelsWithUTF8Titles creates models with UTF-8 multi-byte characters (emojis) for truncation tests
func prepareModelsWithUTF8Titles() []models.Result {
	// Create a title with emojis (each emoji is multiple bytes in UTF-8)
	// 100 emojis = 100 runes, but more than 100 bytes
	emojiTitle := strings.Repeat("🚀", 300) // 300 runes, but 1200 bytes (each emoji is 4 bytes)
	// Mix of ASCII and emojis
	mixedTitle := strings.Repeat("A🚀", 130) // 260 runes (130 * 2), but more than 260 bytes

	return []models.Result{
		{
			ID:        nil,
			Title:     emojiTitle,
			Signature: nil,
			TestOpsID: nil,
			Execution: models.Execution{
				StartTime:  nil,
				EndTime:    nil,
				Status:     "passed",
				Duration:   nil,
				StackTrace: nil,
				Thread:     nil,
			},
			Fields:      make(map[string]string),
			Attachments: []models.Attachment{},
			Steps:       []models.Step{},
			StepType:    "text",
			Params:      make(map[string]string),
			Relations: models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{},
				},
			},
			Muted:   false,
			Message: nil,
		},
		{
			ID:        nil,
			Title:     mixedTitle,
			Signature: nil,
			TestOpsID: nil,
			Execution: models.Execution{
				StartTime:  nil,
				EndTime:    nil,
				Status:     "passed",
				Duration:   nil,
				StackTrace: nil,
				Thread:     nil,
			},
			Fields:      make(map[string]string),
			Attachments: []models.Attachment{},
			Steps:       []models.Step{},
			StepType:    "text",
			Params:      make(map[string]string),
			Relations: models.Relation{
				Suite: models.Suite{
					Data: []models.SuiteData{},
				},
			},
			Muted:   false,
			Message: nil,
		},
	}
}
