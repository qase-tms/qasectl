package env

import (
	"context"
	"errors"
	"github.com/magiconair/properties/assert"
	"github.com/qase-tms/qasectl/internal/models/run"
	"go.uber.org/mock/gomock"
	"reflect"
	"testing"
)

func TestService_CreateEnvironment(t *testing.T) {
	type args struct {
		pc string
		n  string
		d  string
		s  string
		h  string
	}
	tests := []struct {
		name       string
		args       args
		want       run.Environment
		envs       []run.Environment
		wantErr    bool
		errGet     error
		createUse  bool
		errCreate  error
		errMessage string
	}{
		{
			name: "success create environment",
			args: args{
				pc: "projectCode",
				n:  "name",
				d:  "description",
				s:  "slug",
				h:  "host",
			},
			want: run.Environment{
				ID:    1,
				Title: "name",
				Slug:  "slug",
			},
			envs:       []run.Environment{},
			wantErr:    false,
			errCreate:  nil,
			createUse:  true,
			errGet:     nil,
			errMessage: "",
		},
		{
			name: "success get environment",
			args: args{
				pc: "projectCode",
				n:  "name",
				d:  "description",
				s:  "slug",
				h:  "host",
			},
			want: run.Environment{
				ID:    1,
				Title: "name",
				Slug:  "slug",
			},
			envs: []run.Environment{
				{
					ID:    1,
					Title: "name",
					Slug:  "slug",
				},
				{
					ID:    2,
					Title: "name2",
					Slug:  "slug2",
				},
			},
			wantErr:    false,
			errCreate:  nil,
			createUse:  false,
			errGet:     nil,
			errMessage: "",
		},
		{
			name: "failed get environment",
			args: args{
				pc: "projectCode",
				n:  "name",
				d:  "description",
				s:  "slug",
				h:  "host",
			},
			want:       run.Environment{},
			envs:       []run.Environment{},
			wantErr:    true,
			errCreate:  nil,
			createUse:  false,
			errGet:     errors.New("error"),
			errMessage: "failed to get environments: error",
		},
		{
			name: "failed create environment",
			args: args{
				pc: "projectCode",
				n:  "name",
				d:  "description",
				s:  "slug",
				h:  "host",
			},
			want:       run.Environment{},
			envs:       []run.Environment{},
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

			f.client.EXPECT().GetEnvironments(gomock.Any(), tt.args.pc).Return(tt.envs, tt.errGet)
			if tt.createUse {
				f.client.EXPECT().CreateEnvironment(gomock.Any(), tt.args.pc, tt.args.n, tt.args.d, tt.args.s, tt.args.h).Return(tt.want, tt.errCreate)
			}

			srv := NewService(f.client)
			got, err := srv.CreateEnvironment(context.Background(), tt.args.pc, tt.args.n, tt.args.d, tt.args.s, tt.args.h)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("CreateEnvironment() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				assert.Equal(t, tt.errMessage, err.Error())
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateEnvironment() got = %v, want %v", got, tt.want)
			}
		})
	}
}
