package result

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"

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
	CreateRun(ctx context.Context, p, t string, d, e string, m, plan int64, tags []string, isCloud bool, browser string, startTime *int64) (int64, error)
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
		// Find the earliest StartTime from all results
		var startTime *int64
		if minStartTime := s.findMinStartTime(results); minStartTime != nil {
			// Subtract 10 seconds (10000 milliseconds) from the earliest start time
			runStartTime := int64(*minStartTime) - 10000
			startTime = &runStartTime
			logger.Debug("calculated run start time", "startTime", runStartTime, "minResultStartTime", *minStartTime)
		}

		ID, err := s.rs.CreateRun(ctx, p.Project, p.Title, p.Description, "", 0, 0, []string{}, false, "", startTime)
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

	if p.SkipParams {
		for i := range results {
			results[i].Params = nil
		}
	}

	// Filter attachments if extensions are specified
	if p.AttachmentExtensions != "" {
		results = s.filterAttachments(results, p.AttachmentExtensions)
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

// filterAttachments filters attachments based on file extensions
func (s *Service) filterAttachments(results []models.Result, extensions string) []models.Result {
	const op = "result.filterAttachments"
	logger := slog.With("op", op)

	// If no extensions specified, return results as-is
	if extensions == "" {
		return results
	}

	// Parse extensions
	extList := strings.Split(extensions, ",")
	for i, ext := range extList {
		extList[i] = strings.TrimSpace(ext)
		if !strings.HasPrefix(extList[i], ".") {
			extList[i] = "." + extList[i]
		}
	}

	logger.Debug("filtering attachments", "extensions", extList)

	for i := range results {
		// Filter main result attachments
		results[i].Attachments = s.filterAttachmentList(results[i].Attachments, extList)

		// Filter step attachments recursively
		results[i].Steps = s.filterStepAttachments(results[i].Steps, extList)
	}

	return results
}

// filterAttachmentList filters a list of attachments based on extensions
func (s *Service) filterAttachmentList(attachments []models.Attachment, extensions []string) []models.Attachment {
	filtered := make([]models.Attachment, 0, len(attachments))

	for _, attachment := range attachments {
		if s.shouldIncludeAttachment(attachment.Name, extensions) {
			filtered = append(filtered, attachment)
		}
	}

	return filtered
}

// filterStepAttachments recursively filters attachments in steps
func (s *Service) filterStepAttachments(steps []models.Step, extensions []string) []models.Step {
	for i := range steps {
		// Filter step execution attachments
		steps[i].Execution.Attachments = s.filterAttachmentList(steps[i].Execution.Attachments, extensions)

		// Recursively filter nested steps
		if len(steps[i].Steps) > 0 {
			steps[i].Steps = s.filterStepAttachments(steps[i].Steps, extensions)
		}
	}

	return steps
}

// shouldIncludeAttachment checks if an attachment should be included based on its name and allowed extensions
func (s *Service) shouldIncludeAttachment(filename string, extensions []string) bool {
	if len(extensions) == 0 {
		return true
	}

	filename = strings.ToLower(filename)
	for _, ext := range extensions {
		if strings.HasSuffix(filename, strings.ToLower(ext)) {
			return true
		}
	}

	return false
}

// findMinStartTime finds the minimum StartTime from all results
func (s *Service) findMinStartTime(results []models.Result) *float64 {
	var minStartTime *float64

	for _, result := range results {
		if result.Execution.StartTime != nil {
			if minStartTime == nil || *result.Execution.StartTime < *minStartTime {
				minStartTime = result.Execution.StartTime
			}
		}
	}

	return minStartTime
}
