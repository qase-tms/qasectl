package cmd

import (
	"fmt"
	"github.com/qase-tms/qasectl/internal"
	"github.com/spf13/cobra"
)

func init() {
	runCmd.Flags().String(titleFlag, "", "")
	requireFlags(runCmd, titleFlag)

	runCmd.Flags().String(descriptionFlag, "", "")

	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:  "run",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subcmd := args[0]

		switch subcmd {
		case "create":
			title := cmd.Flag(titleFlag).Value.String()
			description := cmd.Flag(descriptionFlag).Value.String()

			api, err := internal.NewApiFromConfig()
			if err != nil {
				return err
			}

			return api.CreateRun(title, description)
		default:
			return fmt.Errorf("unknown subcmd %q", subcmd)
		}
	},
}
