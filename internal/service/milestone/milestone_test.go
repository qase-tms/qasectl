package milestone

import (
	"context"
	"errors"
	"github.com/magiconair/properties/assert"
	"github.com/qase-tms/qasectl/internal/models/run"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
)

func TestService_CreateMilestone(t *testing.T) {
	type args struct {
		projectCode string
		n           string
		d           string
		s           string
		t           int64
	}
	tests := []struct {
		name       string
		args       args
		want       run.Milestone
		envs       []run.Milestone
		wantErr    bool
		errGet     error
		createUse  bool
		errCreate  error
		errMessage string
	}{
		{
			name: "success create milestone",
			args: args{
				projectCode: "projectCode",
				n:           "name",
				d:           "description",
				s:           "status",
				t:           1,
			},
			want: run.Milestone{
				ID:    1,
				Title: "name",
			},
			envs:       []run.Milestone{},
			wantErr:    false,
			errCreate:  nil,
			createUse:  true,
			errGet:     nil,
			errMessage: "",
		},
		{
			name: "success get milestone",
			args: args{
				projectCode: "projectCode",
				n:           "name",
				d:           "description",
				s:           "status",
				t:           1,
			},
			want: run.Milestone{
				ID:    1,
				Title: "name",
			},
			envs: []run.Milestone{
				{
					ID:    1,
					Title: "name",
				},
				{
					ID:    2,
					Title: "name2",
				},
			},
			wantErr:    false,
			errCreate:  nil,
			createUse:  false,
			errGet:     nil,
			errMessage: "",
		},
		{
			name: "failed get milestone",
			args: args{
				projectCode: "projectCode",
				n:           "name",
				d:           "description",
				s:           "status",
				t:           1,
			},
			want:       run.Milestone{},
			envs:       []run.Milestone{},
			wantErr:    true,
			errCreate:  nil,
			createUse:  false,
			errGet:     errors.New("error"),
			errMessage: "failed to get milestones: error",
		},
		{
			name: "failed create milestone",
			args: args{
				projectCode: "projectCode",
				n:           "name",
				d:           "description",
				s:           "status",
				t:           1,
			},
			want:       run.Milestone{},
			envs:       []run.Milestone{},
			wantErr:    true,
			errCreate:  errors.New("error"),
			createUse:  true,
			errGet:     nil,
			errMessage: "error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)

			f.client.EXPECT().GetMilestones(gomock.Any(), tt.args.projectCode, tt.args.n).Return(tt.envs, tt.errGet)
			if tt.createUse {
				f.client.EXPECT().CreateMilestone(gomock.Any(), tt.args.projectCode, tt.args.n, tt.args.d, tt.args.s, tt.args.t).Return(tt.want, tt.errCreate)
			}

			srv := NewService(f.client)
			got, err := srv.CreateMilestone(context.Background(), tt.args.projectCode, tt.args.n, tt.args.d, tt.args.s, tt.args.t)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("CreateMilestone() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				assert.Equal(t, tt.errMessage, err.Error())
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateMilestone() got = %v, want %v", got, tt.want)
			}
		})
	}
}
