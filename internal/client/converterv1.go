package client

import (
	"context"
	"fmt"
	apiV1Client "github.com/qase-tms/qase-go/qase-api-client"
	models "github.com/qase-tms/qasectl/internal/models/result"
	"os"
	"path/filepath"
)

func (c *ClientV1) convertResultToApiModel(ctx context.Context, projectCode string, result models.Result) apiV1Client.ResultCreate {
	defect := false
	model := apiV1Client.ResultCreate{
		Status:      result.Execution.Status,
		Comment:     *apiV1Client.NewNullableString(result.Message),
		Defect:      *apiV1Client.NewNullableBool(&defect),
		Stacktrace:  *apiV1Client.NewNullableString(result.Execution.StackTrace),
		Param:       result.Params,
		Steps:       c.convertStepToApiModel(ctx, projectCode, result.Steps),
		Attachments: c.convertAttachments(ctx, projectCode, result.Attachments),
	}

	if result.Execution.StartTime != nil {
		startTime := int32(*result.Execution.StartTime)
		model.StartTime = *apiV1Client.NewNullableInt32(&startTime)
	}

	if result.Execution.EndTime != nil {
		endTime := int64(*result.Execution.EndTime)
		model.Time = *apiV1Client.NewNullableInt64(&endTime)
	}

	if result.Execution.Duration != nil {
		duration := int64(*result.Execution.Duration)
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
	} else {
		model.CaseId = result.TestOpsID
	}

	return model
}

func (c *ClientV1) convertAttachments(ctx context.Context, projectCode string, attachments []models.Attachment) []string {
	results := make([]string, 0, len(attachments))

	for _, attachment := range attachments {
		rmAttach := false
		if attachment.FilePath == nil {
			path, err := os.Getwd()
			if err != nil {
				fmt.Println("cannot get executable path", err)
			}

			fp := filepath.Join(path, attachment.Name)
			err = os.WriteFile(fp, *attachment.Content, 0644)
			if err != nil {
				fmt.Println("cannot write file", "error", err)
			}

			attachment.FilePath = &fp

			rmAttach = true
		}

		file, err := os.Open(*attachment.FilePath)
		if err != nil {
			fmt.Println("failed to open file: %w", err)
			if rmAttach {
				removeFile(*attachment.FilePath)
			}
			continue
		}

		hash, err := c.uploadAttachment(ctx, projectCode, []*os.File{file})
		if err != nil {
			fmt.Println("failed to upload attachment: %w", err)
			if rmAttach {
				removeFile(*attachment.FilePath)
			}
			continue
		}

		results = append(results, hash)

		if rmAttach {
			removeFile(*attachment.FilePath)
		}
	}

	return results
}

func removeFile(path string) {
	err := os.Remove(path)
	if err != nil {
		fmt.Println("cannot remove file", "error", err)
	}
}

func (c *ClientV1) convertStepToApiModel(ctx context.Context, projectCode string, steps []models.Step) []apiV1Client.TestStepResultCreate {
	stepModels := make([]apiV1Client.TestStepResultCreate, 0, len(steps))

	for i := range steps {
		stepModel := apiV1Client.TestStepResultCreate{
			Status:      steps[i].Execution.Status,
			Comment:     *apiV1Client.NewNullableString(&steps[i].Execution.Comment),
			Attachments: c.convertAttachments(ctx, projectCode, steps[i].Execution.Attachments),
			Steps:       c.convertStepToMaps(ctx, projectCode, steps[i].Steps),
			Action:      &steps[i].Data.Action,
		}

		stepModels = append(stepModels, stepModel)
	}

	return stepModels
}

func (c *ClientV1) convertStepToMaps(ctx context.Context, projectCode string, steps []models.Step) []map[string]interface{} {
	stepModels := make([]map[string]interface{}, 0, len(steps))

	for i := range steps {

		stepModel := map[string]interface{}{
			"status":      steps[i].Execution.Status,
			"comment":     *apiV1Client.NewNullableString(&steps[i].Execution.Comment),
			"attachments": c.convertAttachments(ctx, projectCode, steps[i].Execution.Attachments),
			"steps":       c.convertStepToMaps(ctx, projectCode, steps[i].Steps),
			"action":      steps[i].Data.Action,
		}

		stepModels = append(stepModels, stepModel)
	}

	return stepModels
}
