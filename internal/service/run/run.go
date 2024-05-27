package run

import (
	"context"
	"fmt"
	"github.com/qase-tms/qasectl/internal/models/run"
)

// client is a client for run
type client interface {
	GetEnvironments(ctx context.Context, projectCode string) ([]run.Environment, error)
	CreateRun(ctx context.Context, projectCode, title string, description string, envID int64) (int64, error)
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
func (s *Service) CreateRun(ctx context.Context, p, t string, d, e string) (int64, error) {
	var envID int64 = 0
	if e != "" {
		es, err := s.client.GetEnvironments(ctx, p)
		if err != nil {
			return 0, fmt.Errorf("failed to get environments: %w", err)
		}
		for _, env := range es {
			if env.Slug == e {
				envID = env.ID
			}
		}
	}

	return s.client.CreateRun(ctx, p, t, d, envID)
}

// CompleteRun completes a run
func (s *Service) CompleteRun(ctx context.Context, projectCode string, runId int64) error {
	return s.client.CompleteRun(ctx, projectCode, runId)
}
