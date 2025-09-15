package fields

import (
	"testing"

	"github.com/qase-tms/qasectl/internal/service/fields/mocks"
	"go.uber.org/mock/gomock"
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
