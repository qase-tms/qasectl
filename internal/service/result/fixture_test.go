package result

import (
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
