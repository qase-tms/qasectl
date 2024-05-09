package run

import "context"

// client is a client for run
type client interface {
	CreateRun(ctx context.Context, projectCode, title string, description *string) (int64, error)
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
func (s *Service) CreateRun(ctx context.Context, projectCode, title string, description *string) (int64, error) {
	return s.client.CreateRun(ctx, projectCode, title, description)
}

// CompleteRun completes a run
func (s *Service) CompleteRun(ctx context.Context, projectCode string, runId int64) error {
	return s.client.CompleteRun(ctx, projectCode, runId)
}
