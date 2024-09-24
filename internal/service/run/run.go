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
	CreateRun(ctx context.Context, projectCode, title string, description, envSlug string, mileID, planID int64) (int64, error)
	CompleteRun(ctx context.Context, projectCode string, runId int64) error
	GetTestRuns(ctx context.Context, projectCode string, start, end int64) ([]run.Run, error)
	DeleteTestRun(ctx context.Context, projectCode string, id int64) error
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
	return s.client.CreateRun(ctx, pc, t, d, e, m, plan)
}

// CompleteRun completes a run
func (s *Service) CompleteRun(ctx context.Context, projectCode string, runId int64) error {
	return s.client.CompleteRun(ctx, projectCode, runId)
}

func (s *Service) DeleteRun(ctx context.Context, projectCode string, ids []int64, all bool, start, end int64) error {
	if len(ids) == 0 && !all {
		return fmt.Errorf("no ids provided")
	}

	foundIDs, err := s.client.GetTestRuns(ctx, projectCode, start, end)
	if err != nil {
		return fmt.Errorf("failed to get test runs: %w", err)
	}

	delIDs := make([]int64, 0)

	if len(ids) > 0 {
		for _, v := range foundIDs {
			for _, id := range ids {
				if v.ID == id {
					delIDs = append(delIDs, id)
				}
			}
		}
	}

	if all {
		for _, v := range foundIDs {
			delIDs = append(delIDs, v.ID)
		}
	}

	for _, id := range delIDs {
		err := s.client.DeleteTestRun(ctx, projectCode, id)
		if err != nil {
			return fmt.Errorf("failed to delete run with id %d: %w", id, err)
		}
	}

	return nil
}
