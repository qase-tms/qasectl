package run

import (
	"github.com/qase-tms/qasectl/cmd/testops/run/complete"
	"github.com/qase-tms/qasectl/cmd/testops/run/create"
	"github.com/qase-tms/qasectl/cmd/testops/run/delete"
	"github.com/spf13/cobra"
)

// Command returns a new cobra command for runs
func Command() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Manage test runs",
	}

	cmd.AddCommand(create.Command())
	cmd.AddCommand(complete.Command())
	cmd.AddCommand(delete.Command())

	return cmd
}
