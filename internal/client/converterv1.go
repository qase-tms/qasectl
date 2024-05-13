package client

import (
	"context"
	"fmt"
	apiV1Client "github.com/qase-tms/qase-go/qase-api-client"
	models "github.com/qase-tms/qasectl/internal/models/result"
	"os"
)

func (c *ClientV1) convertResultToApiModel(ctx context.Context, projectCode string, result models.Result) apiV1Client.ResultCreate {
	defect := false
	model := apiV1Client.ResultCreate{
		CaseId:      result.TestOpsID,
		Status:      result.Execution.Status,
		Comment:     *apiV1Client.NewNullableString(result.Message),
		Defect:      *apiV1Client.NewNullableBool(&defect),
		Stacktrace:  *apiV1Client.NewNullableString(result.Execution.StackTrace),
		Param:       result.Params,
		Steps:       c.convertStepToApiModel(ctx, projectCode, result.Steps),
		Attachments: c.convertAttachments(ctx, projectCode, result.Attachments),
	}
	if result.Execution.StartTime != nil {
		startTime := int32(result.Execution.StartTime.Unix())
		model.StartTime = *apiV1Client.NewNullableInt32(&startTime)
	}

	if result.Execution.EndTime != nil {
		endTime := result.Execution.EndTime.Unix()
		model.Time = *apiV1Client.NewNullableInt64(&endTime)
	}

	if result.Execution.Duration != 0 {
		duration := int64(result.Execution.Duration.Seconds())
		model.TimeMs = *apiV1Client.NewNullableInt64(&duration)
	}

	if result.TestOpsID == nil {
		caseModel := apiV1Client.ResultCreateCase{
			Title: &result.Title,
		}

		if v, ok := result.Fields["description"]; ok {
			caseModel.Description = *apiV1Client.NewNullableString(&v)
		}

		if v, ok := result.Fields["severity"]; ok {
			caseModel.Severity = &v
		}

		if v, ok := result.Fields["priority"]; ok {
			caseModel.Priority = &v
		}

		if v, ok := result.Fields["layer"]; ok {
			caseModel.Layer = &v
		}

		if v, ok := result.Fields["preconditions"]; ok {
			caseModel.Preconditions = *apiV1Client.NewNullableString(&v)
		}

		if v, ok := result.Fields["postconditions"]; ok {
			caseModel.Postconditions = *apiV1Client.NewNullableString(&v)
		}

		if result.Relations.Suite.Data != nil {
			var suite string

			for _, s := range result.Relations.Suite.Data {
				suite += s.Title + "\t"
			}

			caseModel.SuiteTitle = *apiV1Client.NewNullableString(&suite)
		}

		model.Case = &caseModel
	}

	return model
}

func (c *ClientV1) convertAttachments(ctx context.Context, projectCode string, attachments []models.Attachment) []string {
	results := make([]string, 0, len(attachments))

	for _, attachment := range attachments {
		// TODO: convert attachment from content to file
		file, err := os.Open(*attachment.FilePath)
		if err != nil {
			fmt.Println("failed to open file: %w", err)
			continue
		}

		hash, err := c.uploadAttachment(ctx, projectCode, []*os.File{file})
		if err != nil {
			fmt.Println("failed to upload attachment: %w", err)
			continue
		}

		results = append(results, hash)
	}

	return results
}

func (c *ClientV1) convertStepToApiModel(ctx context.Context, projectCode string, steps []models.Step) []apiV1Client.TestStepResultCreate {
	stepModels := make([]apiV1Client.TestStepResultCreate, 0, len(steps))

	for _, step := range steps {
		stepModel := apiV1Client.TestStepResultCreate{
			Status:      step.Execution.Status,
			Comment:     *apiV1Client.NewNullableString(&step.Execution.Comment),
			Attachments: c.convertAttachments(ctx, projectCode, step.Execution.Attachments),
			Steps:       c.convertStepToMaps(ctx, projectCode, step.Steps),
			Action:      &step.Data.Action,
		}

		stepModels = append(stepModels, stepModel)
	}

	return stepModels
}

func (c *ClientV1) convertStepToMaps(ctx context.Context, projectCode string, steps []models.Step) []map[string]interface{} {
	stepModels := make([]map[string]interface{}, 0, len(steps))

	for _, step := range steps {

		stepModel := map[string]interface{}{
			"Status":      step.Execution.Status,
			"Comment":     *apiV1Client.NewNullableString(&step.Execution.Comment),
			"Attachments": c.convertAttachments(ctx, projectCode, step.Execution.Attachments),
			"Steps":       c.convertStepToMaps(ctx, projectCode, step.Steps),
			"Action":      step.Data.Action,
		}

		stepModels = append(stepModels, stepModel)
	}

	return stepModels
}
