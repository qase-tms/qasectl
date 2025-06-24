package filter

import (
	"context"
	"errors"
	"testing"

	"github.com/qase-tms/qasectl/internal/models/plan"
	"go.uber.org/mock/gomock"
)

func TestService_GetFilteredResults(t *testing.T) {
	type args struct {
		project   string
		planID    int64
		framework string
	}
	type planArgs struct {
		plan plan.PlanDetailed
		err  error
	}
	tests := []struct {
		name       string
		args       args
		planArgs   planArgs
		want       string
		wantErr    bool
		errMessage string
	}{
		{
			name: "success playwright with single case",
			args: args{
				project:   "test",
				planID:    1,
				framework: "playwright",
			},
			planArgs: planArgs{
				plan: plan.PlanDetailed{
					ID:    1,
					Title: "Test Plan",
					Cases: []int64{123},
				},
				err: nil,
			},
			want:       "(Qase ID: 123)",
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "success playwright with multiple cases",
			args: args{
				project:   "test",
				planID:    1,
				framework: "playwright",
			},
			planArgs: planArgs{
				plan: plan.PlanDetailed{
					ID:    1,
					Title: "Test Plan",
					Cases: []int64{123, 456, 789},
				},
				err: nil,
			},
			want:       "(Qase ID: 123|456|789)",
			wantErr:    false,
			errMessage: "",
		},
		{
			name: "success playwright with empty cases",
			args: args{
				project:   "test",
				planID:    1,
				framework: "playwright",
			},
			planArgs: planArgs{
				plan: plan.PlanDetailed{
					ID:    1,
					Title: "Test Plan",
					Cases: []int64{},
				},
				err: nil,
			},
			want:       "",
			wantErr:    true,
			errMessage: "no cases found in plan",
		},
		{
			name: "unsupported framework",
			args: args{
				project:   "test",
				planID:    1,
				framework: "unsupported",
			},
			planArgs: planArgs{
				plan: plan.PlanDetailed{
					ID:    1,
					Title: "Test Plan",
					Cases: []int64{123},
				},
				err: nil,
			},
			want:       "",
			wantErr:    true,
			errMessage: "unsupported framework: unsupported",
		},
		{
			name: "failed to get plan",
			args: args{
				project:   "test",
				planID:    1,
				framework: "playwright",
			},
			planArgs: planArgs{
				plan: plan.PlanDetailed{},
				err:  errors.New("failed to get plan"),
			},
			want:       "",
			wantErr:    true,
			errMessage: "failed to get plan",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)

			f.client.EXPECT().GetPlan(gomock.Any(), tt.args.project, tt.args.planID).Return(tt.planArgs.plan, tt.planArgs.err)

			s := NewService(f.client)

			got, err := s.GetFilteredResults(context.Background(), tt.args.project, tt.args.planID, tt.args.framework)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetFilteredResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMessage {
				t.Errorf("Service.GetFilteredResults() error = %v, wantErr %v", err, tt.errMessage)
				return
			}
			if got != tt.want {
				t.Errorf("Service.GetFilteredResults() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrepareForPlaywright(t *testing.T) {
	tests := []struct {
		name string
		IDs  []int64
		want string
	}{
		{
			name: "single ID",
			IDs:  []int64{123},
			want: "(Qase ID: 123)",
		},
		{
			name: "multiple IDs",
			IDs:  []int64{123, 456, 789},
			want: "(Qase ID: 123|456|789)",
		},
		{
			name: "empty IDs",
			IDs:  []int64{},
			want: "(Qase ID: )",
		},
		{
			name: "large numbers",
			IDs:  []int64{999999, 1000000},
			want: "(Qase ID: 999999|1000000)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prepareForPlaywright(tt.IDs); got != tt.want {
				t.Errorf("prepareForPlaywright() = %v, want %v", got, tt.want)
			}
		})
	}
}
