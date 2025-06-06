package result

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"

	models "github.com/qase-tms/qasectl/internal/models/result"
	"golang.org/x/sync/errgroup"
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
	CreateRun(ctx context.Context, p, t string, d, e string, m, plan int64, tags []string) (int64, error)
	CompleteRun(ctx context.Context, projectCode string, runId int64) error
}

const (
	// MaxWorkerCount is the maximum number of workers
	MaxWorkerCount = 5
)

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
		ID, err := s.rs.CreateRun(ctx, p.Project, p.Title, p.Description, "", 0, 0, []string{})
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

	if len(p.Statuses) > 0 {
		for i := range results {
			if status, ok := p.Statuses[results[i].Execution.Status]; ok {
				results[i].Execution.Status = status
			}
		}
	}

	err = s.uploadResults(ctx, p.Project, p.Batch, runID, results)
	if err != nil {
		return fmt.Errorf("failed to upload results: %w", err)
	}

	if isTestRunCreated {
		err := s.rs.CompleteRun(ctx, p.Project, runID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) uploadResults(ctx context.Context, project string, batchSize, runID int64, results []models.Result) error {
	batchCount := (int64(len(results)) + batchSize - 1) / batchSize
	batches := make([][]models.Result, 0, batchCount)

	for i := int64(0); i < int64(len(results)); i += batchSize {
		end := i + batchSize
		if end > int64(len(results)) {
			end = int64(len(results))
		}
		batches = append(batches, results[i:end])
	}

	g, ctx := errgroup.WithContext(ctx)

	batchCh := make(chan []models.Result)

	workerCount := runtime.NumCPU()
	if workerCount > MaxWorkerCount {
		workerCount = MaxWorkerCount
	}

	for i := 0; i < workerCount; i++ {
		g.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case batch, ok := <-batchCh:
					if !ok {
						return nil
					}

					if err := s.client.UploadData(ctx, project, runID, batch); err != nil {
						return err
					}
				}
			}
		})
	}

	g.Go(func() error {
		defer close(batchCh)
		for _, batch := range batches {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case batchCh <- batch:
			}
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
