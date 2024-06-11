package run

import (
	"context"
	"fmt"
	"github.com/qase-tms/qasectl/internal/models/run"
)

// client is a client for run
//
//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type client interface {
	GetEnvironments(ctx context.Context, projectCode string) ([]run.Environment, error)
	CreateRun(ctx context.Context, projectCode, title, description string, envID, mileID, planID int64) (int64, error)
	CompleteRun(ctx context.Context, projectCode string, runId int64) error
}

// Service is a Service for run
type Service struct {
	client client
}

// NewService creates a new run Service
func NewService(client client) *Service {
	return &Service{client: client}
}

// CreateRun creates a new run
func (s *Service) CreateRun(ctx context.Context, pc, t, d, e string, m, plan int64) (int64, error) {
	var envID int64 = 0
	if e != "" {
		es, err := s.client.GetEnvironments(ctx, pc)
		if err != nil {
			return 0, fmt.Errorf("failed to get environments: %w", err)
		}
		for _, env := range es {
			if env.Slug == e {
				envID = env.ID
			}
		}
	}

	return s.client.CreateRun(ctx, pc, t, d, envID, m, plan)
}

// CompleteRun completes a run
func (s *Service) CompleteRun(ctx context.Context, projectCode string, runId int64) error {
	return s.client.CompleteRun(ctx, projectCode, runId)
}
