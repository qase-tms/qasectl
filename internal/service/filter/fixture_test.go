package filter

import (
	"testing"

	"github.com/qase-tms/qasectl/internal/service/filter/mocks"
	"go.uber.org/mock/gomock"
)

type fixture struct {
	client *mocks.Mockclient
}

func newFixture(t *testing.T) *fixture {
	ctr := gomock.NewController(t)

	return &fixture{
		client: mocks.NewMockclient(ctr),
	}
}
