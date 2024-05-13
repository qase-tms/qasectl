package client

import (
	"context"
	apiV1Client "github.com/qase-tms/qase-go/qase-api-client"
	models "github.com/qase-tms/qasectl/internal/models/result"
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
	ctx, client := c.getApiV1Client(ctx)

	resp, _, err := client.RunsAPI.
		CreateRun(ctx, projectCode).
		RunCreate(apiV1Client.RunCreate{
			Title:       title,
			Description: description,
		}).
		Execute()

	if err != nil {
		return 0, err
	}

	return resp.Result.GetId(), nil
}

// CompleteRun completes a run
func (c *ClientV1) CompleteRun(ctx context.Context, projectCode string, runId int64) error {
	ctx, client := c.getApiV1Client(ctx)

	_, _, err := client.RunsAPI.
		CompleteRun(ctx, projectCode, int32(runId)).
		Execute()

	return err

}

// UploadData uploads results to Qase
func (c *ClientV1) UploadData(ctx context.Context, project string, runID int64, results []models.Result) error {
	ctx, client := c.getApiV1Client(ctx)

	resultModels := make([]apiV1Client.ResultCreate, 0, len(results))
	for _, result := range results {
		resultModels = append(resultModels, c.convertResultToApiModel(ctx, project, result))
	}

	bulkModel := apiV1Client.NewResultcreateBulk(resultModels)

	_, _, err := client.ResultsAPI.
		CreateResultBulk(ctx, project, int32(runID)).
		ResultcreateBulk(*bulkModel).
		Execute()

	return err
}

// uploadAttachment uploads attachments to Qase
func (c *ClientV1) uploadAttachment(ctx context.Context, projectCode string, file []*os.File) (string, error) {
	ctx, client := c.getApiV1Client(ctx)

	resp, _, err := client.AttachmentsAPI.
		UploadAttachment(ctx, projectCode).
		File(file).
		Execute()
	if err != nil {
		return "", err
	}

	return *resp.Result[0].Hash, err
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
