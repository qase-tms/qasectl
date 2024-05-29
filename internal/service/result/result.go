package result

import (
	"context"
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
	CreateRun(ctx context.Context, p, t string, d, e, m, plan string) (int64, error)
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
func (s *Service) Upload(ctx context.Context, p UploadParams) {
	const op = "result.parser.import"
	logger := slog.With("op", op)

	results, err := s.parser.Parse()
	if err != nil {
		logger.Error("failed to parse results", "error", err)
		return
	}

	logger.Info("number of results found", "count", len(results))

	if len(results) == 0 {
		logger.Info("no results to upload")
		return
	}
	runID := p.RunID
	isTestRunCreated := false
	if runID == 0 {
		ID, err := s.rs.CreateRun(ctx, p.Project, p.Title, p.Description, "", "", "")
		if err != nil {
			logger.Error("failed to create run", "error", err)
			return
		}
		runID = ID
		isTestRunCreated = true
	}

	if int64(len(results)) < p.Batch {
		err := s.client.UploadData(ctx, p.Project, runID, results)
		if err != nil {
			logger.Error("failed to upload results", "error", err)
		}
	} else {
		for i := int64(0); i < int64(len(results)); i += p.Batch {
			err := s.client.UploadData(ctx, p.Project, runID, results[i:i+p.Batch])
			if err != nil {
				logger.Error("failed to upload results", "error", err)
			}
		}
	}

	if isTestRunCreated {
		err := s.rs.CompleteRun(ctx, p.Project, runID)
		if err != nil {
			logger.Error("failed to complete run", "error", err)
		}
	}
}
