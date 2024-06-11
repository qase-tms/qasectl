package result

import (
	"context"
	"errors"
	models "github.com/qase-tms/qasectl/internal/models/result"
	"go.uber.org/mock/gomock"
	"testing"
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
			wantErr:    false,
			errMessage: "",
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
			wantErr:    false,
			errMessage: "",
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
			wantErr:    false,
			errMessage: "",
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
				},
				err:    errors.New("failed upload data"),
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
			name: "failed upload data",
			args: args{
				p: UploadParams{
					Project:     "project",
					Title:       "title",
					Description: "description",
					RunID:       0,
					Batch:       20,
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
				isUsed: true,
			},
			wantErr:    false,
			errMessage: "",
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
				f.rs.EXPECT().CreateRun(gomock.Any(), tt.args.p.Project, tt.args.p.Title, tt.args.p.Description, "", int64(0), int64(0)).Return(tt.rArgs.model, tt.rArgs.err)
			}

			if tt.cArgs.isUsed {
				f.rs.EXPECT().CompleteRun(gomock.Any(), tt.args.p.Project, tt.rArgs.model).Return(tt.cArgs.err)
			}

			if tt.args.isUsed {
				if tt.args.count != 1 {
					f.client.EXPECT().UploadData(gomock.Any(), tt.args.p.Project, tt.args.runID, gomock.Any()).Return(tt.args.err).Times(tt.args.count)
				} else {
					f.client.EXPECT().UploadData(gomock.Any(), tt.args.p.Project, tt.args.runID, tt.pArgs.models).Return(tt.args.err).Times(tt.args.count)
				}
			}
			s := NewService(f.client, f.parser, f.rs)

			s.Upload(context.Background(), tt.args.p)
		})
	}
}
