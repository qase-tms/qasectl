package env

import (
	"github.com/qase-tms/qasectl/cmd/testops/env/create"
	"github.com/spf13/cobra"
)

// Command returns a new cobra command for envs
func Command() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "env",
		Short: "Manage environments",
	}

	cmd.AddCommand(create.Command())

	return cmd
}
