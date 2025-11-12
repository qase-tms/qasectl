package run

import (
	"context"
	"errors"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/qase-tms/qasectl/internal/models/run"
	"go.uber.org/mock/gomock"
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
		pc      string
		t       string
		d       string
		e       string
		m       int64
		plan    int64
		tags    []string
		isCloud bool
		browser string
		args    baseArgs
	}
	tests := []struct {
		name       string
		args       args
		want       int64
		wantErr    bool
		errMessage string
	}{
		{
			name: "success",
			args: args{
				pc:      "test",
				t:       "test",
				d:       "test",
				e:       "test",
				m:       0,
				plan:    0,
				tags:    []string{},
				isCloud: false,
				browser: "",
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
			name: "success with tags",
			args: args{
				pc:      "test",
				t:       "test",
				d:       "test",
				e:       "test",
				m:       0,
				plan:    0,
				tags:    []string{"tag1", "tag2"},
				isCloud: false,
				browser: "",
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
			name: "success with cloud and browser",
			args: args{
				pc:      "test",
				t:       "test",
				d:       "test",
				e:       "test",
				m:       0,
				plan:    0,
				tags:    []string{},
				isCloud: true,
				browser: "chromium",
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
			name: "failed to create run",
			args: args{
				pc:      "test",
				t:       "test",
				d:       "test",
				e:       "test",
				m:       0,
				plan:    0,
				tags:    []string{},
				isCloud: false,
				browser: "",
				args: baseArgs{
					err:    errors.New("error"),
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

			if tt.args.args.isUsed {
				f.client.EXPECT().CreateRun(
					gomock.Any(),
					tt.args.pc,
					tt.args.t,
					tt.args.d,
					tt.args.e,
					tt.args.m,
					tt.args.plan,
					tt.args.tags,
					tt.args.isCloud,
					tt.args.browser,
					gomock.Any(), // startTime
				).
					Return(tt.want, tt.args.args.err)
			}

			s := NewService(f.client)

			got, err := s.CreateRun(
				context.Background(),
				tt.args.pc,
				tt.args.t,
				tt.args.d,
				tt.args.e,
				tt.args.m,
				tt.args.plan,
				tt.args.tags,
				tt.args.isCloud,
				tt.args.browser,
				nil, // startTime
			)
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

func TestService_DeleteRun(t *testing.T) {
	type args struct {
		projectCode string
		ids         []int64
		all         bool
		start       int64
		end         int64
	}
	type argsGetTr struct {
		models []run.Run
		isUsed bool
		err    error
	}
	type argsDelTr struct {
		ids    []int64
		isUsed bool
		err    error
	}
	tests := []struct {
		name       string
		args       args
		argsGetTr  argsGetTr
		argsDelTr  argsDelTr
		wantErr    bool
		errMessage string
	}{
		{
			name: "success with ids",
			args: args{
				projectCode: "test",
				ids:         []int64{1},
				all:         false,
				start:       0,
				end:         0,
			},
			argsGetTr: argsGetTr{
				models: []run.Run{{
					ID: 1,
				}},
				isUsed: true,
				err:    nil,
			},
			argsDelTr: argsDelTr{
				ids:    []int64{1},
				isUsed: true,
				err:    nil,
			},
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "success with all",
			args: args{
				projectCode: "test",
				ids:         []int64{},
				all:         true,
				start:       0,
				end:         0,
			},
			argsGetTr: argsGetTr{
				models: []run.Run{{
					ID: 1,
				}},
				isUsed: true,
				err:    nil,
			},
			argsDelTr: argsDelTr{
				ids:    []int64{1},
				isUsed: true,
				err:    nil,
			},
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "incorrect args",
			args: args{
				projectCode: "test",
				ids:         []int64{},
				all:         false,
				start:       0,
				end:         0,
			},
			argsGetTr: argsGetTr{
				models: []run.Run{{}},
				isUsed: false,
				err:    nil,
			},
			argsDelTr: argsDelTr{
				ids:    []int64{1},
				isUsed: false,
				err:    nil,
			},
			wantErr:    true,
			errMessage: "no ids provided",
		},
		{
			name: "error get test runs",
			args: args{
				projectCode: "test",
				ids:         []int64{1},
				all:         false,
				start:       0,
				end:         0,
			},
			argsGetTr: argsGetTr{
				models: []run.Run{{}},
				isUsed: true,
				err:    errors.New("error"),
			},
			argsDelTr: argsDelTr{
				ids:    []int64{},
				isUsed: false,
				err:    nil,
			},
			wantErr:    true,
			errMessage: "failed to get test runs: error",
		},
		{
			name: "error delete test run",
			args: args{
				projectCode: "test",
				ids:         []int64{1},
				all:         false,
				start:       0,
				end:         0,
			},
			argsGetTr: argsGetTr{
				models: []run.Run{{
					ID: 1,
				}},
				isUsed: true,
				err:    nil,
			},
			argsDelTr: argsDelTr{
				ids:    []int64{1},
				isUsed: true,
				err:    errors.New("error"),
			},
			wantErr:    true,
			errMessage: "failed to delete run with id 1: error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)

			if tt.argsGetTr.isUsed {
				f.client.EXPECT().GetTestRuns(gomock.Any(), tt.args.projectCode, tt.args.start, tt.args.end).Return(tt.argsGetTr.models, tt.argsGetTr.err)
			}
			if tt.argsDelTr.isUsed {
				f.client.EXPECT().DeleteTestRun(gomock.Any(), tt.args.projectCode, tt.argsDelTr.ids[0]).Return(tt.argsDelTr.err).Times(len(tt.argsDelTr.ids))
			}

			s := NewService(f.client)

			if err := s.DeleteRun(context.Background(), tt.args.projectCode, tt.args.ids, tt.args.all, tt.args.start, tt.args.end); err != nil {
				if !tt.wantErr {
					t.Errorf("DeleteRun() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				assert.Equal(t, tt.errMessage, err.Error())
			}
		})
	}
}
