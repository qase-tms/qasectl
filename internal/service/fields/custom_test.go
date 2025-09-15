package fields

import (
	"context"
	"errors"
	"testing"

	"github.com/qase-tms/qasectl/internal/models/fields/custom"
	"github.com/qase-tms/qasectl/internal/service/fields/mocks"
	"go.uber.org/mock/gomock"
)

func TestService_RemoveCustomFields(t *testing.T) {
	type args struct {
		params RemoveCustomFieldsParams
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		errMessage    string
		mockSetup     func(*mocks.Mockclient)
		expectedCalls int
	}{
		{
			name: "success remove custom field by ID",
			args: args{
				params: RemoveCustomFieldsParams{
					FieldID: func() *int32 { id := int32(123); return &id }(),
					All:     false,
				},
			},
			wantErr:    false,
			errMessage: "",
			mockSetup: func(m *mocks.Mockclient) {
				m.EXPECT().RemoveCustomFieldByID(gomock.Any(), int32(123)).Return(nil).Times(1)
			},
			expectedCalls: 1,
		},
		{
			name: "success remove all custom fields",
			args: args{
				params: RemoveCustomFieldsParams{
					FieldID: nil,
					All:     true,
				},
			},
			wantErr:    false,
			errMessage: "",
			mockSetup: func(m *mocks.Mockclient) {
				fields := []custom.CustomField{
					{ID: 1, Title: "Field 1"},
					{ID: 2, Title: "Field 2"},
				}
				m.EXPECT().GetCustomFields(gomock.Any()).Return(fields, nil).Times(1)
				m.EXPECT().RemoveCustomFieldByID(gomock.Any(), int32(1)).Return(nil).Times(1)
				m.EXPECT().RemoveCustomFieldByID(gomock.Any(), int32(2)).Return(nil).Times(1)
			},
			expectedCalls: 3,
		},
		{
			name: "success remove all custom fields with empty list",
			args: args{
				params: RemoveCustomFieldsParams{
					FieldID: nil,
					All:     true,
				},
			},
			wantErr:    false,
			errMessage: "",
			mockSetup: func(m *mocks.Mockclient) {
				fields := []custom.CustomField{}
				m.EXPECT().GetCustomFields(gomock.Any()).Return(fields, nil).Times(1)
			},
			expectedCalls: 1,
		},
		{
			name: "error when fieldID and all are both not set",
			args: args{
				params: RemoveCustomFieldsParams{
					FieldID: nil,
					All:     false,
				},
			},
			wantErr:    true,
			errMessage: "fieldID or all is required",
			mockSetup: func(m *mocks.Mockclient) {
				// No mock expectations for this case
			},
			expectedCalls: 0,
		},
		{
			name: "error when RemoveCustomFieldByID fails",
			args: args{
				params: RemoveCustomFieldsParams{
					FieldID: func() *int32 { id := int32(123); return &id }(),
					All:     false,
				},
			},
			wantErr:    true,
			errMessage: "failed to remove custom field",
			mockSetup: func(m *mocks.Mockclient) {
				m.EXPECT().RemoveCustomFieldByID(gomock.Any(), int32(123)).Return(errors.New("failed to remove custom field")).Times(1)
			},
			expectedCalls: 1,
		},
		{
			name: "error when GetCustomFields fails",
			args: args{
				params: RemoveCustomFieldsParams{
					FieldID: nil,
					All:     true,
				},
			},
			wantErr:    true,
			errMessage: "failed to get custom fields: failed to fetch fields",
			mockSetup: func(m *mocks.Mockclient) {
				m.EXPECT().GetCustomFields(gomock.Any()).Return(nil, errors.New("failed to fetch fields")).Times(1)
			},
			expectedCalls: 1,
		},
		{
			name: "error when RemoveCustomFieldByID fails during all fields removal",
			args: args{
				params: RemoveCustomFieldsParams{
					FieldID: nil,
					All:     true,
				},
			},
			wantErr:    true,
			errMessage: "failed to remove custom field 1 with title Field 1: failed to remove custom field",
			mockSetup: func(m *mocks.Mockclient) {
				fields := []custom.CustomField{
					{ID: 1, Title: "Field 1"},
				}
				m.EXPECT().GetCustomFields(gomock.Any()).Return(fields, nil).Times(1)
				m.EXPECT().RemoveCustomFieldByID(gomock.Any(), int32(1)).Return(errors.New("failed to remove custom field")).Times(1)
			},
			expectedCalls: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFixture(t)
			tt.mockSetup(f.client)

			srv := NewService(f.client)
			err := srv.RemoveCustomFields(context.Background(), tt.args.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("RemoveCustomFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMessage {
				t.Errorf("RemoveCustomFields() error = %v, wantErr %v", err, tt.errMessage)
			}
		})
	}
}
