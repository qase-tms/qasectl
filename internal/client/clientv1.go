package client

import (
	"context"
	"fmt"
	apiV1Client "github.com/qase-tms/qase-go/qase-api-client"
	models "github.com/qase-tms/qasectl/internal/models/result"
	"github.com/qase-tms/qasectl/internal/models/run"
	"log/slog"
	"os"
)

const (
	defaultLimit = 50
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

// CreateEnvironment creates a new environment
func (c *ClientV1) CreateEnvironment(ctx context.Context, pc, n, d, s, h string) (run.Environment, error) {
	const op = "client.clientv1.createenvironment"
	logger := slog.With("op", op)

	logger.Debug("creating environment", "projectCode", pc, "name", n, "description", d, "slug", s, "host", h)

	ctx, client := c.getApiV1Client(ctx)

	m := apiV1Client.EnvironmentCreate{
		Title: n,
		Slug:  s,
	}

	if d != "" {
		m.SetDescription(d)
	}

	if h != "" {
		m.SetHost(h)
	}

	resp, r, err := client.EnvironmentsAPI.
		CreateEnvironment(ctx, pc).
		EnvironmentCreate(m).
		Execute()

	if err != nil {
		logger.Debug("failed to create environment", "response", r)
		return run.Environment{}, fmt.Errorf("failed to create environment: %w", err)
	}

	env := run.Environment{
		Title: n,
		Slug:  s,
		ID:    resp.Result.GetId(),
	}

	logger.Info("created environment", "environment", env)

	return env, nil
}

// GetEnvironments returns environments
func (c *ClientV1) GetEnvironments(ctx context.Context, projectCode string) ([]run.Environment, error) {
	const op = "client.clientv1.getenvironments"
	logger := slog.With("op", op)

	logger.Debug("getting environments", "projectCode", projectCode)

	ctx, client := c.getApiV1Client(ctx)

	environments := make([]run.Environment, 0)

	var offset int32 = 0
	for {
		resp, r, err := client.EnvironmentsAPI.
			GetEnvironments(ctx, projectCode).
			Limit(defaultLimit).
			Offset(offset).
			Execute()

		if err != nil {
			logger.Debug("failed to get environments", "response", r)
			return nil, fmt.Errorf("failed to get environments: %w", err)
		}

		for _, env := range resp.Result.Entities {
			environments = append(environments, run.Environment{
				Title: env.GetTitle(),
				ID:    env.GetId(),
				Slug:  env.GetSlug(),
			})
		}

		if resp.Result.GetTotal() <= offset {
			break
		} else {
			offset += defaultLimit
		}
	}

	logger.Debug("got environments", "environments", environments)

	return environments, nil
}

// GetMilestones returns milestones
func (c *ClientV1) GetMilestones(ctx context.Context, projectCode, milestoneName string) ([]run.Milestone, error) {
	const op = "client.clientv1.getmilestones"
	logger := slog.With("op", op)

	logger.Debug("getting milestones", "projectCode", projectCode, "milestoneName", milestoneName)

	ctx, client := c.getApiV1Client(ctx)

	resp, r, err := client.MilestonesAPI.
		GetMilestones(ctx, projectCode).
		Search(milestoneName).
		Execute()

	if err != nil {
		logger.Debug("failed to get milestones", "response", r)
		return nil, fmt.Errorf("failed to get milestones: %w", err)
	}

	milestones := make([]run.Milestone, 0, len(resp.Result.Entities))
	for _, milestone := range resp.Result.Entities {
		milestones = append(milestones, run.Milestone{
			Title: milestone.GetTitle(),
			ID:    milestone.GetId(),
		})
	}

	logger.Debug("got milestones", "milestones", milestones)

	return milestones, nil
}

// GetPlans returns plans
func (c *ClientV1) GetPlans(ctx context.Context, projectCode string) ([]run.Plan, error) {
	const op = "client.clientv1.getplans"
	logger := slog.With("op", op)

	logger.Debug("getting plans", "projectCode", projectCode)

	ctx, client := c.getApiV1Client(ctx)

	plans := make([]run.Plan, 0)

	var offset int32 = 0
	for {
		resp, r, err := client.PlansAPI.
			GetPlans(ctx, projectCode).
			Limit(defaultLimit).
			Offset(offset).
			Execute()

		if err != nil {
			logger.Debug("failed to get plans", "response", r)
			return nil, fmt.Errorf("failed to get plans: %w", err)
		}

		for _, plan := range resp.Result.Entities {
			plans = append(plans, run.Plan{
				Title: plan.GetTitle(),
				ID:    plan.GetId(),
			})
		}

		if resp.Result.GetTotal() <= offset {
			break
		} else {
			offset += defaultLimit
		}
	}

	logger.Debug("got plans", "plans", plans)

	return plans, nil
}

// CreateRun creates a new run
func (c *ClientV1) CreateRun(ctx context.Context, projectCode, title string, description string, envID, mileID, planID int64) (int64, error) {
	const op = "client.clientv1.createrun"
	logger := slog.With("op", op)

	ctx, client := c.getApiV1Client(ctx)

	m := apiV1Client.RunCreate{
		Title: title,
	}

	if description != "" {
		m.SetDescription(description)
	}

	if envID != 0 {
		m.SetEnvironmentId(envID)
	}

	if mileID != 0 {
		m.SetMilestoneId(mileID)
	}

	if planID != 0 {
		m.SetPlanId(planID)
	}

	logger.Debug("creating run", "projectCode", projectCode, "model", m)

	resp, r, err := client.RunsAPI.
		CreateRun(ctx, projectCode).
		RunCreate(m).
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
