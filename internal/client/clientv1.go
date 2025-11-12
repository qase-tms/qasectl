package client

import (
	"context"
	"log/slog"
	"os"
	"time"

	apiV1Client "github.com/qase-tms/qase-go/qase-api-client"
	"github.com/qase-tms/qasectl/internal/models/fields/custom"
	"github.com/qase-tms/qasectl/internal/models/plan"
	models "github.com/qase-tms/qasectl/internal/models/result"
	"github.com/qase-tms/qasectl/internal/models/run"
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

// CreateMilestone creates a new milestone
func (c *ClientV1) CreateMilestone(ctx context.Context, projectCode, n, d, s string, t int64) (run.Milestone, error) {
	const op = "client.clientv1.createmilestone"
	logger := slog.With("op", op)

	logger.Debug("creating milestone", "projectCode", projectCode, "name", n, "description", d, "status", s, "dueDate", t)

	ctx, client := c.getApiV1Client(ctx)

	m := apiV1Client.MilestoneCreate{
		Title: n,
	}

	if d != "" {
		m.SetDescription(d)
	}

	if s != "" {
		m.SetStatus(s)
	}

	if t != 0 {
		m.SetDueDate(t)
	}

	resp, r, err := client.MilestonesAPI.
		CreateMilestone(ctx, projectCode).
		MilestoneCreate(m).
		Execute()

	if err != nil {
		return run.Milestone{}, NewQaseApiError(err.Error(), extractBody(r))
	}

	milestone := run.Milestone{
		Title: n,
		ID:    resp.Result.GetId(),
	}

	logger.Info("created milestone", "milestone", milestone)

	return milestone, nil
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
		return run.Environment{}, NewQaseApiError(err.Error(), extractBody(r))
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
			return nil, NewQaseApiError(err.Error(), extractBody(r))
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
		return nil, NewQaseApiError(err.Error(), extractBody(r))
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
			return nil, NewQaseApiError(err.Error(), extractBody(r))
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

func (c *ClientV1) GetPlan(ctx context.Context, projectCode string, planID int64) (plan.PlanDetailed, error) {
	const op = "client.clientv1.getplan"
	logger := slog.With("op", op)

	logger.Debug("getting plan", "projectCode", projectCode, "planID", planID)

	ctx, client := c.getApiV1Client(ctx)

	resp, r, err := client.PlansAPI.
		GetPlan(ctx, projectCode, int32(planID)).
		Execute()

	if err != nil {
		return plan.PlanDetailed{}, NewQaseApiError(err.Error(), extractBody(r))
	}

	cases := make([]int64, 0, len(resp.Result.GetCases()))
	for _, c := range resp.Result.GetCases() {
		cases = append(cases, c.GetCaseId())
	}

	return plan.PlanDetailed{
		ID:    resp.Result.GetId(),
		Title: resp.Result.GetTitle(),
		Cases: cases,
	}, nil
}

// CreateRun creates a new run
func (c *ClientV1) CreateRun(ctx context.Context, projectCode, title string, description, envSlug string, mileID, planID int64, tags []string, isCloud bool, browser string, startTime *int64) (int64, error) {
	const op = "client.clientv1.createrun"
	logger := slog.With("op", op)

	ctx, client := c.getApiV1Client(ctx)

	m := apiV1Client.RunCreate{
		Title: title,
	}

	m.SetIsAutotest(true)

	if description != "" {
		m.SetDescription(description)
	}

	if envSlug != "" {
		m.SetEnvironmentSlug(envSlug)
	}

	if mileID != 0 {
		m.SetMilestoneId(mileID)
	}

	if planID != 0 {
		m.SetPlanId(planID)
	}

	if len(tags) > 0 {
		m.SetTags(tags)
	}

	if isCloud {
		m.SetIsCloud(true)
	}

	if browser != "" {
		cloudConfig := apiV1Client.RunCreateCloudRunConfig{
			Browser: &browser,
		}

		m.SetCloudRunConfig(cloudConfig)
	}

	if startTime != nil {
		// Convert milliseconds to time.Time and format as "YYYY-MM-DD HH:MM:SS" in UTC
		startTimeSeconds := *startTime / 1000
		t := time.Unix(startTimeSeconds, 0).UTC()
		startTimeStr := t.Format("2006-01-02 15:04:05")
		m.StartTime = &startTimeStr
		logger.Debug("setting run start time", "startTime", *startTime, "startTimeSeconds", startTimeSeconds, "formatted", startTimeStr)
	}

	logger.Debug("creating run", "projectCode", projectCode, "model", m)

	resp, r, err := client.RunsAPI.
		CreateRun(ctx, projectCode).
		RunCreate(m).
		Execute()

	if err != nil {
		return 0, NewQaseApiError(err.Error(), extractBody(r))
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
		return NewQaseApiError(err.Error(), extractBody(r))
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

	_, r, err := client.ResultsAPI.
		CreateResultBulk(ctx, project, int32(runID)).
		ResultcreateBulk(*bulkModel).
		Execute()

	if err != nil {
		return NewQaseApiError(err.Error(), extractBody(r))
	}

	return nil
}

// GetTestRuns returns test runs
func (c *ClientV1) GetTestRuns(ctx context.Context, projectCode string, start, end int64) ([]run.Run, error) {
	const op = "client.clientv1.gettestruns"
	logger := slog.With("op", op)

	logger.Debug("getting test runs", "projectCode", projectCode)

	ctx, client := c.getApiV1Client(ctx)

	testRuns := make([]run.Run, 0)

	var offset int32 = 0
	for {
		req := client.RunsAPI.
			GetRuns(ctx, projectCode).
			Limit(defaultLimit).
			Offset(offset)

		if start != 0 {
			req = req.FromStartTime(start)
		}

		if end != 0 {
			req = req.ToStartTime(end)
		}

		resp, r, err := req.Execute()

		if err != nil {
			return nil, NewQaseApiError(err.Error(), extractBody(r))
		}

		for _, testRun := range resp.Result.Entities {
			testRuns = append(testRuns, run.Run{
				ID: testRun.GetId(),
			})
		}

		if resp.Result.GetFiltered() <= offset {
			break
		} else {
			offset += defaultLimit
		}
	}

	logger.Debug("got test runs", "testRuns", testRuns)

	return testRuns, nil
}

// DeleteTestRun deletes test run
func (c *ClientV1) DeleteTestRun(ctx context.Context, projectCode string, id int64) error {
	const op = "client.clientv1.deletetestrun"
	logger := slog.With("op", op)

	logger.Debug("deleting test run", "projectCode", projectCode, "id", id)

	ctx, client := c.getApiV1Client(ctx)

	_, r, err := client.RunsAPI.
		DeleteRun(ctx, projectCode, int32(id)).
		Execute()

	if err != nil {
		return NewQaseApiError(err.Error(), extractBody(r))
	}

	return nil
}

// GetCustomFields returns custom fields
func (c *ClientV1) GetCustomFields(ctx context.Context) ([]custom.CustomField, error) {
	const op = "client.clientv1.getcustomfields"
	logger := slog.With("op", op)

	logger.Debug("getting custom fields")

	ctx, client := c.getApiV1Client(ctx)

	customFields := make([]custom.CustomField, 0)

	var offset int32 = 0
	for {
		resp, r, err := client.CustomFieldsAPI.
			GetCustomFields(ctx).
			Limit(defaultLimit).
			Offset(offset).
			Execute()

		if err != nil {
			return nil, NewQaseApiError(err.Error(), extractBody(r))
		}

		for _, field := range resp.Result.Entities {
			customFields = append(customFields, custom.CustomField{
				ID:    field.GetId(),
				Title: field.GetTitle(),
			})
		}

		if resp.Result.GetTotal() <= offset {
			break
		} else {
			offset += defaultLimit
		}
	}

	logger.Debug("got custom fields", "customFields", customFields)

	return customFields, nil
}

// RemoveCustomFieldByID removes a custom field by ID
func (c *ClientV1) RemoveCustomFieldByID(ctx context.Context, fieldID int32) error {
	const op = "client.clientv1.removecustomfield"
	logger := slog.With("op", op)

	logger.Debug("removing custom field", "fieldID", fieldID)

	ctx, client := c.getApiV1Client(ctx)

	_, r, err := client.CustomFieldsAPI.
		DeleteCustomField(ctx, fieldID).
		Execute()

	if err != nil {
		return NewQaseApiError(err.Error(), extractBody(r))
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
		return "", NewQaseApiError(err.Error(), extractBody(r))
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
