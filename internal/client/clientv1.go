package client

import (
	"context"
	apiV1Client "github.com/qase-tms/qase-go/qase-api-client"
)

// ClientV1 is a client for Qase API v1
type ClientV1 struct {
	// token is a token for Qase API
	token string
	// projectCode is a code of the project
	projectCode string
}

// NewClientV1 creates a new client for Qase API v1
func NewClientV1(token, projectCode string) *ClientV1 {
	return &ClientV1{
		token:       token,
		projectCode: projectCode,
	}
}

// CreateRun creates a new run
func (c *ClientV1) CreateRun(ctx context.Context, title string, description *string) (int64, error) {
	ctx, client := c.getApiV1Client(ctx)

	resp, _, err := client.RunsAPI.
		CreateRun(ctx, c.projectCode).
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
func (c *ClientV1) CompleteRun(ctx context.Context, runId int64) error {
	ctx, client := c.getApiV1Client(ctx)

	_, _, err := client.RunsAPI.
		CompleteRun(ctx, c.projectCode, int32(runId)).
		Execute()

	return err

}

// getApiV1Client returns a context and a client for Qase API v1
func (c *ClientV1) getApiV1Client(ctx context.Context) (context.Context, *apiV1Client.APIClient) {
	ctx = context.WithValue(ctx, apiV1Client.ContextAPIKeys,
		map[string]apiV1Client.APIKey{
			"Token": {Key: c.token},
		})

	cfg := apiV1Client.NewConfiguration()
	client := apiV1Client.NewAPIClient(cfg)

	return ctx, client
}