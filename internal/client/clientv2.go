package client

import (
	"context"
	apiV2Client "github.com/qase-tms/qase-go/qase-api-v2-client"
	models "github.com/qase-tms/qasectl/internal/models/result"
	"log/slog"
)

// ClientV2 is a client for Qase API v2
type ClientV2 struct {
	// token is a token for Qase API
	token    string
	clientV1 *ClientV1
}

// NewClientV2 creates a new client for Qase API v1
func NewClientV2(token string, clientV1 *ClientV1) *ClientV2 {
	return &ClientV2{
		token:    token,
		clientV1: clientV1,
	}
}

// UploadData uploads results to Qase
func (c *ClientV2) UploadData(ctx context.Context, project string, runID int64, results []models.Result) error {
	const op = "client.clientv2.uploaddata"
	logger := slog.With("op", op)

	logger.Debug("uploading data", "project", project, "runID", runID, "results", results)

	ctx, client := c.getApiV2Client(ctx)

	resultModels := make([]apiV2Client.ResultCreate, 0, len(results))
	for _, result := range results {
		resultModels = append(resultModels, c.convertResultToApiModel(ctx, project, result))
	}

	logger.Debug("converted results", "resultModels", resultModels)

	bulkModel := apiV2Client.NewCreateResultsRequestV2()
	bulkModel.SetResults(resultModels)

	r, err := client.ResultsAPI.
		CreateResultsV2(ctx, project, runID).
		CreateResultsRequestV2(*bulkModel).
		Execute()

	if err != nil {
		return NewQaseApiError(err.Error(), r.Body)
	}
	return nil
}

// getApiV2Client returns a context and a client for Qase API v2
func (c *ClientV2) getApiV2Client(ctx context.Context) (context.Context, *apiV2Client.APIClient) {
	ctx = context.WithValue(ctx, apiV2Client.ContextAPIKeys,
		map[string]apiV2Client.APIKey{
			"TokenAuth": {Key: c.token},
		})

	cfg := apiV2Client.NewConfiguration()
	client := apiV2Client.NewAPIClient(cfg)

	return ctx, client
}
