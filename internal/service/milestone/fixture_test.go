package milestone

import (
	"github.com/qase-tms/qasectl/internal/service/milestone/mocks"
	"go.uber.org/mock/gomock"
	"testing"
)

type fixture struct {
	client *mocks.Mockclient
}

func newFixture(t *testing.T) *fixture {
	ctrl := gomock.NewController(t)

	return &fixture{
		client: mocks.NewMockclient(ctrl),
	}
}
