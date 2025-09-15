package delete

import (
	"fmt"
	"log/slog"

	"github.com/qase-tms/qasectl/cmd/flags"
	"github.com/qase-tms/qasectl/internal/client"
	"github.com/qase-tms/qasectl/internal/service/fields"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	fieldIDFlag = "id"
	allFlag     = "all"
)

// Command returns a new cobra command for delete custom fields
func Command() *cobra.Command {
	var (
		fieldID int32
		all     bool
	)

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete custom fields",
		Example: "qasectl testops field custom delete --id 1 --token 'TOKEN'",
		RunE: func(cmd *cobra.Command, args []string) error {
			token := viper.GetString(flags.TokenFlag)

			c := client.NewClientV1(token)
			s := fields.NewService(c)

			var params fields.RemoveCustomFieldsParams
			if fieldID != 0 {
				params.FieldID = &fieldID
			}
			params.All = all

			err := s.RemoveCustomFields(cmd.Context(), params)
			if err != nil {
				return fmt.Errorf("failed to delete custom fields: %w", err)
			}

			slog.Info("Custom fields deleted")

			return nil
		},
	}

	cmd.Flags().Int32Var(&fieldID, fieldIDFlag, 0, "ID of custom field to delete")
	cmd.Flags().BoolVar(&all, allFlag, false, "delete all custom fields in the project")
	cmd.MarkFlagsOneRequired(fieldIDFlag, allFlag)
	cmd.MarkFlagsMutuallyExclusive(fieldIDFlag, allFlag)

	return cmd
}
