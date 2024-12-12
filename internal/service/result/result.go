package result

import (
	"context"
	"fmt"
	models "github.com/qase-tms/qasectl/internal/models/result"
	"log/slog"
)

//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type client interface {
	UploadData(ctx context.Context, project string, runID int64, results []models.Result) error
}

//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type Parser interface {
	Parse() ([]models.Result, error)
}

//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type runService interface {
	CreateRun(ctx context.Context, p, t string, d, e string, m, plan int64) (int64, error)
	CompleteRun(ctx context.Context, projectCode string, runId int64) error
}

// Service is a service for importing data
type Service struct {
	client client
	parser Parser
	rs     runService
}

// NewService creates a new service
func NewService(client client, parser Parser, rs runService) *Service {
	return &Service{client: client, parser: parser, rs: rs}
}

// Upload imports the data
func (s *Service) Upload(ctx context.Context, p UploadParams) error {
	const op = "result.parser.import"
	logger := slog.With("op", op)

	results, err := s.parser.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse results: %w", err)
	}

	logger.Info("number of results found", "count", len(results))

	if len(results) == 0 {
		return fmt.Errorf("no results to upload")
	}

	runID := p.RunID
	isTestRunCreated := false
	if runID == 0 {
		ID, err := s.rs.CreateRun(ctx, p.Project, p.Title, p.Description, "", 0, 0)
		if err != nil {
			return err
		}
		runID = ID
		isTestRunCreated = true
	}

	if p.Suite != "" {
		s := []models.SuiteData{
			{Title: p.Suite,
				PublicID: nil,
			},
		}

		for i := range results {
			results[i].Relations.Suite.Data = append(s, results[i].Relations.Suite.Data...)
		}
	}

	if int64(len(results)) < p.Batch {
		err := s.client.UploadData(ctx, p.Project, runID, results)
		if err != nil {
			return fmt.Errorf("failed to upload results: %w", err)
		}
	} else {
		for i := int64(0); i < int64(len(results)); i += p.Batch {
			end := i + p.Batch

			if end > int64(len(results)) {
				end = int64(len(results))
			}

			err := s.client.UploadData(ctx, p.Project, runID, results[i:end])
			if err != nil {
				return fmt.Errorf("failed to upload results: %w", err)
			}
		}
	}

	if isTestRunCreated {
		err := s.rs.CompleteRun(ctx, p.Project, runID)
		if err != nil {
			return err
		}
	}

	return nil
}
