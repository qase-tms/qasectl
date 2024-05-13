package result

import "context"
import (
	models "github.com/qase-tms/qasectl/internal/models/result"
)

type client interface {
	UploadData(ctx context.Context, project string, runID int64, results []models.Result) error
}

type Parser interface {
	Parse() ([]models.Result, error)
}

// Service is a service for importing data
type Service struct {
	client client
	parser Parser
}

// NewService creates a new service
func NewService(client client, parser Parser) *Service {
	return &Service{client: client, parser: parser}
}

// Import imports the data
func (s *Service) Import(ctx context.Context, project string, runID int64) error {
	results, err := s.parser.Parse()
	if err != nil {
		return err
	}

	return s.client.UploadData(ctx, project, runID, results)
}
