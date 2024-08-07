package run

import (
	"context"
	"errors"
	"github.com/magiconair/properties/assert"
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
		m    int64
		plan int64
		args baseArgs
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
				pc:   "test",
				t:    "test",
				d:    "test",
				e:    "test",
				m:    0,
				plan: 0,
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
				pc:   "test",
				t:    "test",
				d:    "test",
				e:    "test",
				m:    0,
				plan: 0,
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
				f.client.EXPECT().CreateRun(gomock.Any(),
					tt.args.pc,
					tt.args.t,
					tt.args.d,
					tt.args.e,
					tt.args.m,
					tt.args.plan,
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
