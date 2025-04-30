package result

import (
	models "github.com/qase-tms/qasectl/internal/models/result"
	"github.com/qase-tms/qasectl/internal/service/result/mocks"
	"go.uber.org/mock/gomock"
	"testing"
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
			Params:      make(map[string]string),
			Relations:   models.Relation{},
			Muted:       false,
			Message:     nil,
		},
		{
			ID:        nil,
			Title:     "Test 2",
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
			Relations:   models.Relation{},
			Muted:       false,
			Message:     nil,
		},
	}
}
