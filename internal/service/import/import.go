package _import

import "context"
import (
	models "github.com/qase-tms/qasectl/internal/models/import"
)

type client interface {
	UploadData(ctx context.Context, project string, runID int64, results []models.Result) error
}

// Service is a service for importing data
type Service struct {
	client *client
}

// NewService creates a new service
func NewService(client *client) *Service {
	return &Service{client: client}
}
