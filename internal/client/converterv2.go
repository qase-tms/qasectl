package client

import (
	"context"
	apiV2Client "github.com/qase-tms/qase-go/qase-api-v2-client"
	models "github.com/qase-tms/qasectl/internal/models/result"
)

func (c *ClientV2) convertResultToApiModel(ctx context.Context, projectCode string, result models.Result) apiV2Client.ResultCreate {
	model := apiV2Client.ResultCreate{
		Title:       result.Title,
		Signature:   result.Signature,
		Attachments: c.clientV1.convertAttachments(ctx, projectCode, result.Attachments),
		Steps:       c.convertSteps(ctx, projectCode, result.Steps),
		Params:      &result.Params,
		ParamGroups: result.ParamGroups,
	}

	if result.TestOpsID != nil {
		model.SetTestopsId(*result.TestOpsID)
	}

	if result.Message != nil {
		model.SetMessage(*result.Message)
	}

	model.SetDefect(false)
	model.SetExecution(c.createExecution(result.Execution))
	model.SetStepsType(c.getStepsType("classic"))
	model.SetRelations(c.createRelations(result.Relations.Suite.Data))
	model.SetFields(c.convertFields(result.Fields))

	return model
}

func (c *ClientV2) createExecution(execution models.Execution) apiV2Client.ResultExecution {
	exec := apiV2Client.NewResultExecution(execution.Status)

	if execution.StartTime != nil {
		exec.SetStartTime(*execution.StartTime / 1000.0)
	}

	if execution.EndTime != nil {
		exec.SetEndTime(*execution.EndTime / 1000.0)
	}

	if execution.Duration != nil {
		exec.SetDuration(int64(*execution.Duration))
	}

	if execution.StackTrace != nil {
		exec.SetStacktrace(*execution.StackTrace)
	}

	if execution.Thread != nil {
		exec.SetThread(*execution.Thread)
	}

	return *exec
}

func (c *ClientV2) getStepsType(stepType string) apiV2Client.ResultStepsType {
	st, _ := apiV2Client.NewResultStepsTypeFromValue(stepType)
	return *st
}

func (c *ClientV2) createRelations(suiteData []models.SuiteData) apiV2Client.ResultRelations {
	relationItems := make([]apiV2Client.RelationSuiteItem, 0, len(suiteData))
	for _, data := range suiteData {
		if data.Title == "" {
			continue
		}
		relationItems = append(relationItems, *apiV2Client.NewRelationSuiteItem(data.Title))
	}

	suite := apiV2Client.NewRelationSuiteWithDefaults()
	suite.SetData(relationItems)

	relations := apiV2Client.NewResultRelations()
	relations.SetSuite(*suite)
	return *relations
}

func (c *ClientV2) convertFields(fields map[string]string) apiV2Client.ResultCreateFields {
	rcf := apiV2Client.NewResultCreateFields()
	rcf.AdditionalProperties = make(map[string]interface{})

	for k, v := range fields {
		switch k {
		case "author":
			rcf.SetAuthor(v)
		case "description":
			rcf.SetDescription(v)
		case "preconditions":
			rcf.SetPreconditions(v)
		case "postconditions":
			rcf.SetPostconditions(v)
		case "layer":
			rcf.SetLayer(v)
		case "severity":
			rcf.SetSeverity(v)
		case "priority":
			rcf.SetPriority(v)
		case "behavior":
			rcf.SetBehavior(v)
		case "type":
			rcf.SetType(v)
		case "muted":
			rcf.SetMuted(v)
		case "isFlaky":
			rcf.SetIsFlaky(v)
		default:
			rcf.AdditionalProperties[k] = v
		}
	}

	return *rcf
}

func (c *ClientV2) convertSteps(ctx context.Context, projectCode string, steps []models.Step) []apiV2Client.ResultStep {
	stepModels := make([]apiV2Client.ResultStep, len(steps))
	for i := range steps {
		stepModels[i] = c.createStepModel(ctx, projectCode, steps[i])
	}
	return stepModels
}

func (c *ClientV2) createStepModel(ctx context.Context, projectCode string, step models.Step) apiV2Client.ResultStep {
	m := apiV2Client.NewResultStep()
	d := apiV2Client.NewResultStepData(step.Data.Action)
	m.SetData(*d)
	m.SetExecution(c.createStepExecution(ctx, projectCode, step.Execution))
	m.SetSteps(c.convertStepMaps(ctx, projectCode, step.Steps))

	return *m
}

func (c *ClientV2) createStepExecution(ctx context.Context, projectCode string, execution models.StepExecution) apiV2Client.ResultStepExecution {
	status, _ := apiV2Client.NewResultStepStatusFromValue(execution.Status)
	exec := apiV2Client.NewResultStepExecution(*status)

	exec.Attachments = c.clientV1.convertAttachments(ctx, projectCode, execution.Attachments)

	if execution.Duration != nil {
		exec.SetDuration(int64(*execution.Duration))
	}

	exec.SetComment(execution.Comment)
	return *exec
}

func (c *ClientV2) convertStepMaps(ctx context.Context, projectCode string, steps []models.Step) []map[string]interface{} {
	stepMaps := make([]map[string]interface{}, len(steps))
	for i := range steps {
		stepMaps[i] = map[string]interface{}{
			"data":      apiV2Client.NewResultStepData(steps[i].Data.Action),
			"execution": c.createStepExecution(ctx, projectCode, steps[i].Execution),
			"steps":     c.convertStepMaps(ctx, projectCode, steps[i].Steps),
		}
	}
	return stepMaps
}
