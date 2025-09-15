package field

import (
	"github.com/qase-tms/qasectl/cmd/testops/field/custom"
	"github.com/spf13/cobra"
)

// Command returns a new cobra command for fields
func Command() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "field",
		Short: "Manage fields",
	}

	cmd.AddCommand(custom.Command())

	return cmd
}
