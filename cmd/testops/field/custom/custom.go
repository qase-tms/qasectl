package custom

import (
	"github.com/qase-tms/qasectl/cmd/testops/field/custom/delete"
	"github.com/spf13/cobra"
)

// Command returns a new cobra command for custom fields
func Command() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "custom",
		Short: "Manage custom fields",
	}

	cmd.AddCommand(delete.Command())

	return cmd
}
