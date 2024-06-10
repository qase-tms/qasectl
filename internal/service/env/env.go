package env

import (
	"context"
	"fmt"
	"github.com/qase-tms/qasectl/internal/models/run"
	"log/slog"
)

// client is a client for env
//
//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type client interface {
	CreateEnvironment(ctx context.Context, pc, n, d, s, h string) (run.Environment, error)
	GetEnvironments(ctx context.Context, projectCode string) ([]run.Environment, error)
}

// Service is a Service for env
type Service struct {
	client client
}

// NewService creates a new env Service
func NewService(client client) *Service {
	return &Service{client: client}
}

// CreateEnvironment creates a new environment
func (srv *Service) CreateEnvironment(ctx context.Context, pc, n, d, s, h string) (run.Environment, error) {
	envs, err := srv.client.GetEnvironments(ctx, pc)
	if err != nil {
		return run.Environment{}, fmt.Errorf("failed to get environments: %w", err)
	}

	for _, env := range envs {
		if env.Slug == s {
			slog.Info("environment already exists")
			return env, nil
		}
	}

	return srv.client.CreateEnvironment(ctx, pc, n, d, s, h)
}
