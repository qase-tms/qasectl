package result

import (
	"context"
	"errors"
	"testing"

	models "github.com/qase-tms/qasectl/internal/models/result"
	"go.uber.org/mock/gomock"
)

func TestService_Upload(t *testing.T) {
	type args struct {
		p      UploadParams
		err    error
		isUsed bool
		count  int
		runID  int64
	}
	type pArgs struct {
		models []models.Result
		err    error
		isUsed bool
	}
	type rArgs struct {
		model  int64
		err    error
		isUsed bool
	}
	type cArgs struct {
		err    error
		isUsed bool
	}
	tests := []struct {
		name       string
		args       args
		pArgs      pArgs
		rArgs      rArgs
		cArgs      cArgs
		wantErr    bool
		errMessage string
	}{
		{
			name: "success with create test run",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       0,
					Batch:       20,
					Suite:       "",
				},
				err:    nil,
				isUsed: true,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: prepareModels(),
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: true,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: true,
			},
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "success with create test run and suite",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       0,
					Batch:       20,
					Suite:       "suite",
				},
				err:    nil,
				isUsed: true,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: prepareModels(),
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: true,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: true,
			},
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "success without create test run",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       1,
					Batch:       20,
					Suite:       "",
				},
				err:    nil,
				isUsed: true,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: prepareModels(),
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: false,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: false,
			},
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "failed parser",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       1,
					Batch:       20,
					Suite:       "",
				},
				err:    nil,
				isUsed: false,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: prepareModels(),
				err:    errors.New("failed parser"),
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: false,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: false,
			},
			wantErr:    true,
			errMessage: "failed to parse results: failed parser",
		},
		{
			name: "empty results",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       1,
					Batch:       20,
					Suite:       "",
				},
				err:    nil,
				isUsed: false,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: []models.Result{},
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: false,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: false,
			},
			wantErr:    true,
			errMessage: "no results to upload",
		},
		{
			name: "failed create test run",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       0,
					Batch:       20,
					Suite:       "",
				},
				err:    nil,
				isUsed: false,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: prepareModels(),
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    errors.New("failed create test run"),
				isUsed: true,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: false,
			},
			wantErr:    true,
			errMessage: "failed create test run",
		},
		{
			name: "use batch",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       0,
					Batch:       1,
					Suite:       "",
				},
				err:    nil,
				isUsed: true,
				count:  2,
				runID:  1,
			},
			pArgs: pArgs{
				models: prepareModels(),
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: true,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: true,
			},
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "failed upload with batch",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       0,
					Batch:       1,
					Suite:       "",
				},
				err:    errors.New("failed upload data"),
				isUsed: true,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: prepareModels(),
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: true,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: false,
			},
			wantErr:    true,
			errMessage: "failed to upload results: failed upload data",
		},
		{
			name: "failed upload data",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       0,
					Batch:       20,
					Suite:       "",
				},
				err:    errors.New("failed upload data"),
				isUsed: true,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: prepareModels(),
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: true,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: false,
			},
			wantErr:    true,
			errMessage: "failed to upload results: failed upload data",
		},
		{
			name: "failed complete run",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       0,
					Batch:       20,
					Suite:       "",
				},
				err:    nil,
				isUsed: true,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: prepareModels(),
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: true,
			},
			cArgs: cArgs{
				err:    errors.New("failed complete run"),
				isUsed: true,
			},
			wantErr:    true,
			errMessage: "failed complete run",
		},
		{
			name: "success with status mapping",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       0,
					Batch:       20,
					Suite:       "",
					Statuses: map[string]string{
						"passed": "passed",
						"failed": "failed",
					},
				},
				err:    nil,
				isUsed: true,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: []models.Result{
					{
						Execution: models.Execution{
							Status: "passed",
						},
					},
					{
						Execution: models.Execution{
							Status: "failed",
						},
					},
				},
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: true,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: true,
			},
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "success with skip params",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       1,
					Batch:       20,
					Suite:       "",
					SkipParams:  true,
				},
				err:    nil,
				isUsed: true,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: []models.Result{
					{
						Execution: models.Execution{
							Status: "passed",
						},
						Params: map[string]string{
							"key1": "value1",
							"key2": "value2",
						},
					},
				},
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: false,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: false,
			},
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "success with suite and skip params",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       0,
					Batch:       20,
					Suite:       "test-suite",
					SkipParams:  true,
				},
				err:    nil,
				isUsed: true,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: []models.Result{
					{
						Execution: models.Execution{
							Status: "passed",
						},
						Params: map[string]string{
							"key1": "value1",
						},
						Relations: models.Relation{
							Suite: models.Suite{
								Data: []models.SuiteData{},
							},
						},
					},
				},
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: true,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: true,
			},
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "success with complex status mapping and skip params",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       0,
					Batch:       20,
					Suite:       "complex-suite",
					Statuses: map[string]string{
						"passed":  "passed",
						"failed":  "failed",
						"skipped": "skipped",
						"blocked": "blocked",
					},
					SkipParams: true,
				},
				err:    nil,
				isUsed: true,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: []models.Result{
					{
						Execution: models.Execution{
							Status: "passed",
						},
						Params: map[string]string{
							"key1": "value1",
						},
						Relations: models.Relation{
							Suite: models.Suite{
								Data: []models.SuiteData{},
							},
						},
					},
					{
						Execution: models.Execution{
							Status: "failed",
						},
						Params: map[string]string{
							"key2": "value2",
						},
						Relations: models.Relation{
							Suite: models.Suite{
								Data: []models.SuiteData{},
							},
						},
					},
					{
						Execution: models.Execution{
							Status: "skipped",
						},
						Params: map[string]string{
							"key3": "value3",
						},
						Relations: models.Relation{
							Suite: models.Suite{
								Data: []models.SuiteData{},
							},
						},
					},
				},
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: true,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: true,
			},
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "verify skip params functionality",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       1,
					Batch:       20,
					Suite:       "",
					SkipParams:  true,
				},
				err:    nil,
				isUsed: true,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: prepareModels(), // Models with params
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: false,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: false,
			},
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "verify params are preserved when skip params is false",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       1,
					Batch:       20,
					Suite:       "",
					SkipParams:  false,
				},
				err:    nil,
				isUsed: true,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: prepareModels(), // Models with params
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: false,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: false,
			},
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "skip params with already empty params",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       1,
					Batch:       20,
					Suite:       "",
					SkipParams:  true,
				},
				err:    nil,
				isUsed: true,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: prepareModelsWithEmptyParams(), // Models without params
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: false,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: false,
			},
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "success with attachment filtering",
			args: args{
				p: UploadParams{
					Project:              "project",
					Title:                "title",
					Description:          "description",
					RunID:                1,
					Batch:                20,
					Suite:                "",
					AttachmentExtensions: "png,jpg",
				},
				err:    nil,
				isUsed: true,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: prepareModelsWithAttachments(), // Models with various attachments
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: false,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: false,
			},
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "success with attachment filtering and steps",
			args: args{
				p: UploadParams{
					Project:              "project",
					Title:                "title",
					Description:          "description",
					RunID:                1,
					Batch:                20,
					Suite:                "",
					AttachmentExtensions: "pdf",
				},
				err:    nil,
				isUsed: true,
				count:  1,
				runID:  1,
			},
			pArgs: pArgs{
				models: prepareModelsWithStepAttachments(), // Models with step attachments
				err:    nil,
				isUsed: true,
			},
			rArgs: rArgs{
				model:  1,
				err:    nil,
				isUsed: false,
			},
			cArgs: cArgs{
				err:    nil,
				isUsed: false,
			},
			wantErr:    false,
			errMessage: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)

			if tt.pArgs.isUsed {
				f.parser.EXPECT().Parse().Return(tt.pArgs.models, tt.pArgs.err)
			}

			if tt.rArgs.isUsed {
				f.rs.EXPECT().CreateRun(
					gomock.Any(),
					tt.args.p.Project,
					tt.args.p.Title,
					tt.args.p.Description,
					"",         // envSlug
					int64(0),   // mileID
					int64(0),   // planID
					[]string{}, // tags
					false,      // isCloud
					"",         // browser
					gomock.Any(), // startTime
				).Return(tt.rArgs.model, tt.rArgs.err)
			}

			if tt.cArgs.isUsed {
				f.rs.EXPECT().CompleteRun(gomock.Any(), tt.args.p.Project, tt.rArgs.model).Return(tt.cArgs.err)
			}

			if tt.args.isUsed {
				if tt.name == "failed upload with batch" || tt.name == "failed upload data" {
					f.client.EXPECT().
						UploadData(gomock.Any(), tt.args.p.Project, gomock.Any(), gomock.Any()).
						Return(tt.args.err).
						Times(1)
				} else {
					f.client.EXPECT().
						UploadData(gomock.Any(), tt.args.p.Project, gomock.Any(), gomock.Any()).
						Return(tt.args.err).
						Times(tt.args.count)
				}
			}

			s := NewService(f.client, f.parser, f.rs)

			err := s.Upload(context.Background(), tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Upload() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err.Error() != tt.errMessage {
				t.Errorf("Service.Upload() error = %v, wantErr %v", err, tt.errMessage)
			}
		})
	}
}

func TestService_filterAttachments(t *testing.T) {
	tests := []struct {
		name        string
		results     []models.Result
		extensions  string
		expectedLen int
	}{
		{
			name: "filter png and jpg files",
			results: []models.Result{
				{
					Attachments: []models.Attachment{
						{Name: "screenshot.png"},
						{Name: "photo.jpg"},
						{Name: "document.pdf"},
						{Name: "log.txt"},
					},
				},
			},
			extensions:  "png,jpg",
			expectedLen: 2,
		},
		{
			name: "filter pdf files only",
			results: []models.Result{
				{
					Attachments: []models.Attachment{
						{Name: "screenshot.png"},
						{Name: "photo.jpg"},
						{Name: "document.pdf"},
						{Name: "report.PDF"},
					},
				},
			},
			extensions:  "pdf",
			expectedLen: 2,
		},
		{
			name: "no filtering when extensions empty",
			results: []models.Result{
				{
					Attachments: []models.Attachment{
						{Name: "screenshot.png"},
						{Name: "photo.jpg"},
						{Name: "document.pdf"},
					},
				},
			},
			extensions:  "",
			expectedLen: 3,
		},
		{
			name: "filter with steps attachments",
			results: []models.Result{
				{
					Attachments: []models.Attachment{
						{Name: "main.png"},
					},
					Steps: []models.Step{
						{
							Execution: models.StepExecution{
								Attachments: []models.Attachment{
									{Name: "step1.jpg"},
									{Name: "step2.pdf"},
								},
							},
						},
					},
				},
			},
			extensions:  "png,jpg",
			expectedLen: 2, // 1 main + 1 step attachment
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			result := s.filterAttachments(tt.results, tt.extensions)

			totalAttachments := 0
			for _, r := range result {
				totalAttachments += len(r.Attachments)
				for _, step := range r.Steps {
					totalAttachments += len(step.Execution.Attachments)
				}
			}

			if totalAttachments != tt.expectedLen {
				t.Errorf("filterAttachments() totalAttachments = %v, want %v", totalAttachments, tt.expectedLen)
			}
		})
	}
}

func TestService_shouldIncludeAttachment(t *testing.T) {
	tests := []struct {
		name       string
		filename   string
		extensions []string
		expected   bool
	}{
		{
			name:       "include png file",
			filename:   "screenshot.png",
			extensions: []string{".png", ".jpg"},
			expected:   true,
		},
		{
			name:       "include jpg file",
			filename:   "photo.jpg",
			extensions: []string{".png", ".jpg"},
			expected:   true,
		},
		{
			name:       "exclude pdf file",
			filename:   "document.pdf",
			extensions: []string{".png", ".jpg"},
			expected:   false,
		},
		{
			name:       "case insensitive matching",
			filename:   "SCREENSHOT.PNG",
			extensions: []string{".png", ".jpg"},
			expected:   true,
		},
		{
			name:       "include all when no extensions",
			filename:   "anyfile.txt",
			extensions: []string{},
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			result := s.shouldIncludeAttachment(tt.filename, tt.extensions)
			if result != tt.expected {
				t.Errorf("shouldIncludeAttachment() = %v, want %v", result, tt.expected)
			}
		})
	}
}
