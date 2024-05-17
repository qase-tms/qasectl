package client

import (
	"context"
	"fmt"
	apiV1Client "github.com/qase-tms/qase-go/qase-api-client"
	models "github.com/qase-tms/qasectl/internal/models/result"
	"log/slog"
	"os"
)

// ClientV1 is a client for Qase API v1
type ClientV1 struct {
	// token is a token for Qase API
	token string
}

// NewClientV1 creates a new client for Qase API v1
func NewClientV1(token string) *ClientV1 {
	return &ClientV1{
		token: token,
	}
}

// CreateRun creates a new run
func (c *ClientV1) CreateRun(ctx context.Context, projectCode, title string, description *string) (int64, error) {
	const op = "client.clientv1.createrun"
	logger := slog.With("op", op)

	logger.Debug("creating run", "projectCode", projectCode, "title", title, "description", description)

	ctx, client := c.getApiV1Client(ctx)

	resp, r, err := client.RunsAPI.
		CreateRun(ctx, projectCode).
		RunCreate(apiV1Client.RunCreate{
			Title:       title,
			Description: description,
		}).
		Execute()

	if err != nil {
		logger.Debug("failed to create run", "response", r)
		return 0, fmt.Errorf("failed to create run: %w", err)
	}

	logger.Info("created run", "runID", resp.Result.GetId(), "title", title, "description", description)

	return resp.Result.GetId(), nil
}

// CompleteRun completes a run
func (c *ClientV1) CompleteRun(ctx context.Context, projectCode string, runId int64) error {
	const op = "client.clientv1.completerun"
	logger := slog.With("op", op)

	ctx, client := c.getApiV1Client(ctx)

	logger.Debug("completing run", "projectCode", projectCode, "runId", runId)

	_, r, err := client.RunsAPI.
		CompleteRun(ctx, projectCode, int32(runId)).
		Execute()

	if err != nil {
		logger.Debug("failed to complete run", "response", r)
		return fmt.Errorf("failed to complete run: %w", err)
	}

	logger.Info("completed run", "runId", runId)

	return nil
}

// UploadData uploads results to Qase
func (c *ClientV1) UploadData(ctx context.Context, project string, runID int64, results []models.Result) error {
	const op = "client.clientv1.uploaddata"
	logger := slog.With("op", op)

	logger.Debug("uploading data", "project", project, "runID", runID, "results", results)

	ctx, client := c.getApiV1Client(ctx)

	resultModels := make([]apiV1Client.ResultCreate, 0, len(results))
	for _, result := range results {
		resultModels = append(resultModels, c.convertResultToApiModel(ctx, project, result))
	}

	logger.Debug("converted results", "resultModels", resultModels)

	bulkModel := apiV1Client.NewResultcreateBulk(resultModels)

	resp, r, err := client.ResultsAPI.
		CreateResultBulk(ctx, project, int32(runID)).
		ResultcreateBulk(*bulkModel).
		Execute()

	if err != nil {
		logger.Debug("failed to upload data", "model", resp, "response", r)
		return fmt.Errorf("failed to upload data: %w", err)
	}

	return nil
}

// uploadAttachment uploads attachments to Qase
func (c *ClientV1) uploadAttachment(ctx context.Context, projectCode string, file []*os.File) (string, error) {
	const op = "client.clientv1.uploadattachment"
	logger := slog.With("op", op)

	logger.Debug("uploading attachment", "projectCode", projectCode, "file", file[0].Name())

	ctx, client := c.getApiV1Client(ctx)

	resp, r, err := client.AttachmentsAPI.
		UploadAttachment(ctx, projectCode).
		File(file).
		Execute()
	if err != nil {
		logger.Debug("failed to upload attachment", "response", r)
		return "", fmt.Errorf("failed to upload attachment: %w", err)
	}

	return *resp.Result[0].Hash, nil
}

// getApiV1Client returns a context and a client for Qase API v1
func (c *ClientV1) getApiV1Client(ctx context.Context) (context.Context, *apiV1Client.APIClient) {
	ctx = context.WithValue(ctx, apiV1Client.ContextAPIKeys,
		map[string]apiV1Client.APIKey{
			"TokenAuth": {Key: c.token},
		})

	cfg := apiV1Client.NewConfiguration()
	client := apiV1Client.NewAPIClient(cfg)

	return ctx, client
}
