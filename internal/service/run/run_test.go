package run

import (
	"context"
	"errors"
	"github.com/magiconair/properties/assert"
	"github.com/qase-tms/qasectl/internal/models/run"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestService_CompleteRun(t *testing.T) {

	type args struct {
		projectCode string
		runId       int64
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		err        error
		errMessage string
	}{
		{
			name: "success",
			args: args{
				projectCode: "test",
				runId:       1,
			},
			wantErr:    false,
			err:        nil,
			errMessage: "",
		},
		{
			name: "error",
			args: args{
				projectCode: "test",
				runId:       1,
			},
			wantErr:    true,
			err:        errors.New("error"),
			errMessage: "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)

			f.client.EXPECT().CompleteRun(gomock.Any(), tt.args.projectCode, tt.args.runId).Return(tt.err)

			s := NewService(f.client)

			if err := s.CompleteRun(context.Background(), tt.args.projectCode, tt.args.runId); err != nil {
				if !tt.wantErr {
					t.Errorf("CompleteRun() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				assert.Equal(t, tt.errMessage, err.Error())
			}
		})
	}
}

func TestService_CreateRun(t *testing.T) {
	type args struct {
		pc   string
		t    string
		d    string
		e    string
		m    string
		plan string
		args baseArgs
	}
	type mArgs struct {
		models []run.Milestone
		args   baseArgs
	}
	type eArgs struct {
		models []run.Environment
		args   baseArgs
	}
	type pArgs struct {
		models []run.Plan
		args   baseArgs
	}
	tests := []struct {
		name       string
		args       args
		mArgs      mArgs
		eArgs      eArgs
		pArgs      pArgs
		want       int64
		wantErr    bool
		errMessage string
	}{
		{
			name: "success",
			args: args{
				pc:   "test",
				t:    "test",
				d:    "test",
				e:    "test",
				m:    "test",
				plan: "test",
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			mArgs: mArgs{
				models: []run.Milestone{
					{
						ID:    1,
						Title: "test",
					},
				},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			eArgs: eArgs{
				models: []run.Environment{
					{
						ID:   1,
						Slug: "test",
					}},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			pArgs: pArgs{
				models: []run.Plan{
					{
						ID:    1,
						Title: "test",
					},
				},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			want:       1,
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "environment not found",
			args: args{
				pc:   "test",
				t:    "test",
				d:    "test",
				e:    "test",
				m:    "test",
				plan: "test",
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			mArgs: mArgs{
				models: []run.Milestone{
					{
						ID:    1,
						Title: "test",
					},
				},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			eArgs: eArgs{
				models: []run.Environment{
					{
						ID:   0,
						Slug: "test1",
					}},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			pArgs: pArgs{
				models: []run.Plan{
					{
						ID:    1,
						Title: "test",
					},
				},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			want:       1,
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "failed to get environments",
			args: args{
				pc:   "test",
				t:    "test",
				d:    "test",
				e:    "test",
				m:    "test",
				plan: "test",
				args: baseArgs{
					err:    nil,
					isUsed: false,
				},
			},
			mArgs: mArgs{
				models: []run.Milestone{
					{
						ID:    1,
						Title: "test",
					},
				},
				args: baseArgs{
					err:    nil,
					isUsed: false,
				},
			},
			eArgs: eArgs{
				models: []run.Environment{
					{
						ID:   0,
						Slug: "test",
					}},
				args: baseArgs{
					err:    errors.New("error"),
					isUsed: true,
				},
			},
			pArgs: pArgs{
				models: []run.Plan{
					{
						ID:    1,
						Title: "test",
					},
				},
				args: baseArgs{
					err:    nil,
					isUsed: false,
				},
			},
			want:       0,
			wantErr:    true,
			errMessage: "failed to get environments: error",
		},
		{
			name: "milestone not found",
			args: args{
				pc:   "test",
				t:    "test",
				d:    "test",
				e:    "test",
				m:    "test",
				plan: "test",
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			mArgs: mArgs{
				models: []run.Milestone{
					{
						ID:    0,
						Title: "test",
					},
				},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			eArgs: eArgs{
				models: []run.Environment{
					{
						ID:   1,
						Slug: "test",
					}},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			pArgs: pArgs{
				models: []run.Plan{
					{
						ID:    1,
						Title: "test",
					},
				},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			want:       1,
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "failed to get milestones",
			args: args{
				pc:   "test",
				t:    "test",
				d:    "test",
				e:    "test",
				m:    "test",
				plan: "test",
				args: baseArgs{
					err:    nil,
					isUsed: false,
				},
			},
			mArgs: mArgs{
				models: []run.Milestone{
					{
						ID:    1,
						Title: "test",
					},
				},
				args: baseArgs{
					err:    errors.New("error"),
					isUsed: true,
				},
			},
			eArgs: eArgs{
				models: []run.Environment{
					{
						ID:   0,
						Slug: "test",
					}},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			pArgs: pArgs{
				models: []run.Plan{
					{
						ID:    1,
						Title: "test",
					},
				},
				args: baseArgs{
					err:    nil,
					isUsed: false,
				},
			},
			want:       0,
			wantErr:    true,
			errMessage: "failed to get milestones: error",
		},
		{
			name: "plan not found",
			args: args{
				pc:   "test",
				t:    "test",
				d:    "test",
				e:    "test",
				m:    "test",
				plan: "test",
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			mArgs: mArgs{
				models: []run.Milestone{
					{
						ID:    0,
						Title: "test",
					},
				},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			eArgs: eArgs{
				models: []run.Environment{
					{
						ID:   1,
						Slug: "test",
					}},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			pArgs: pArgs{
				models: []run.Plan{
					{
						ID:    0,
						Title: "test1",
					},
				},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			want:       1,
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "failed to get plans",
			args: args{
				pc:   "test",
				t:    "test",
				d:    "test",
				e:    "test",
				m:    "test",
				plan: "test",
				args: baseArgs{
					err:    nil,
					isUsed: false,
				},
			},
			mArgs: mArgs{
				models: []run.Milestone{
					{
						ID:    1,
						Title: "test",
					},
				},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			eArgs: eArgs{
				models: []run.Environment{
					{
						ID:   0,
						Slug: "test",
					}},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			pArgs: pArgs{
				models: []run.Plan{
					{
						ID:    1,
						Title: "test",
					},
				},
				args: baseArgs{
					err:    errors.New("error"),
					isUsed: true,
				},
			},
			want:       0,
			wantErr:    true,
			errMessage: "failed to get plans: error",
		},
		{
			name: "failed to create run",
			args: args{
				pc:   "test",
				t:    "test",
				d:    "test",
				e:    "test",
				m:    "test",
				plan: "test",
				args: baseArgs{
					err:    errors.New("error"),
					isUsed: true,
				},
			},
			mArgs: mArgs{
				models: []run.Milestone{
					{
						ID:    1,
						Title: "test",
					},
				},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			eArgs: eArgs{
				models: []run.Environment{
					{
						ID:   0,
						Slug: "test",
					}},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			pArgs: pArgs{
				models: []run.Plan{
					{
						ID:    1,
						Title: "test",
					},
				},
				args: baseArgs{
					err:    nil,
					isUsed: true,
				},
			},
			want:       0,
			wantErr:    true,
			errMessage: "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)

			if tt.eArgs.args.isUsed {
				f.client.EXPECT().GetEnvironments(gomock.Any(), tt.args.pc).Return(tt.eArgs.models, tt.eArgs.args.err)
			}
			if tt.mArgs.args.isUsed {
				f.client.EXPECT().GetMilestones(gomock.Any(), tt.args.pc, tt.args.m).Return(tt.mArgs.models, tt.mArgs.args.err)
			}
			if tt.pArgs.args.isUsed {
				f.client.EXPECT().GetPlans(gomock.Any(), tt.args.pc).Return(tt.pArgs.models, tt.pArgs.args.err)
			}

			if tt.args.args.isUsed {
				f.client.EXPECT().CreateRun(gomock.Any(),
					tt.args.pc,
					tt.args.t,
					tt.args.d,
					tt.eArgs.models[0].ID,
					tt.mArgs.models[0].ID,
					tt.pArgs.models[0].ID,
				).
					Return(tt.want, tt.args.args.err)
			}

			s := NewService(f.client)

			got, err := s.CreateRun(context.Background(), tt.args.pc, tt.args.t, tt.args.d, tt.args.e, tt.args.m, tt.args.plan)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("CreateRun() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				assert.Equal(t, tt.errMessage, err.Error())
			}
			if got != tt.want {
				t.Errorf("CreateRun() got = %v, want %v", got, tt.want)
			}
		})
	}
}