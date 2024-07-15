package run

import (
	"context"
)

// client is a client for run
//
//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type client interface {
	CreateRun(ctx context.Context, projectCode, title string, description, envSlug string, mileID, planID int64) (int64, error)
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
	return s.client.CreateRun(ctx, pc, t, d, e, m, plan)
}

// CompleteRun completes a run
func (s *Service) CompleteRun(ctx context.Context, projectCode string, runId int64) error {
	return s.client.CompleteRun(ctx, projectCode, runId)
}
