package milestone

import (
	"context"
	"fmt"
	"github.com/qase-tms/qasectl/internal/models/run"
)

// client is a client for env
//
//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type client interface {
	CreateMilestone(ctx context.Context, projectCode, n, d, s string, t int64) (run.Milestone, error)
	GetMilestones(ctx context.Context, projectCode, milestoneName string) ([]run.Milestone, error)
}

// Service is a service for milestones
type Service struct {
	client client
}

// NewService creates a new service for milestones
func NewService(c client) *Service {
	return &Service{
		client: c,
	}
}

// CreateMilestone creates a new milestone
func (srv *Service) CreateMilestone(ctx context.Context, projectCode, n, d, s string, t int64) (run.Milestone, error) {
	mss, err := srv.client.GetMilestones(ctx, projectCode, n)
	if err != nil {
		return run.Milestone{}, fmt.Errorf("failed to get milestones: %w", err)
	}

	if len(mss) > 0 {
		return mss[0], nil
	}

	return srv.client.CreateMilestone(ctx, projectCode, n, d, s, t)
}
