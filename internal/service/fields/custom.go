package fields

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/qase-tms/qasectl/internal/models/fields/custom"
)

//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type client interface {
	GetCustomFields(ctx context.Context) ([]custom.CustomField, error)
	RemoveCustomFieldByID(ctx context.Context, fieldID int32) error
}

type Service struct {
	client client
}

func NewService(client client) *Service {
	return &Service{client: client}
}

func (s *Service) RemoveCustomFields(ctx context.Context, params RemoveCustomFieldsParams) error {
	const op = "fields.custom.removecustomfields"
	logger := slog.With("op", op)

	logger.Debug("removing custom fields", "params", params)

	if params.FieldID == nil && !params.All {
		return fmt.Errorf("fieldID or all is required")
	}

	if params.FieldID != nil {
		return s.client.RemoveCustomFieldByID(ctx, *params.FieldID)
	}

	fields, err := s.client.GetCustomFields(ctx)
	if err != nil {
		return fmt.Errorf("failed to get custom fields: %w", err)
	}

	for _, field := range fields {
		err := s.client.RemoveCustomFieldByID(ctx, int32(field.ID))
		if err != nil {
			return fmt.Errorf("failed to remove custom field %d with title %s: %w", field.ID, field.Title, err)
		}
	}

	return nil
}
