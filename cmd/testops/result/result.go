package result

import (
	"github.com/qase-tms/qasectl/cmd/testops/result/upload"
	"github.com/spf13/cobra"
)

// Command returns a new cobra command for results
func Command() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "result",
		Short: "Manage test results",
	}

	cmd.AddCommand(upload.Command())

	return cmd
}
