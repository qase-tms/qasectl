package cmd

import (
	"github.com/qase-tms/qasectl/internal"
	"github.com/spf13/cobra"
	"strconv"
)

func init() {
	setCmd.Flags().String(projectCodeFlag, "", "")
	setCmd.Flags().Int(runIdFlag, 0, "")
	rootCmd.AddCommand(setCmd)
}

var setCmd = &cobra.Command{
	Use: "set",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectCode := cmd.Flag(projectCodeFlag)
		runId := cmd.Flag(runIdFlag)

		return internal.UpdateConfig(func(cfg internal.Config) internal.Config {
			if projectCode != nil && projectCode.Changed {
				cfg.ProjectCode = projectCode.Value.String()
			}
			if runId != nil && runId.Changed {
				runId, err := strconv.Atoi(runId.Value.String())
				// looks like cobra should pass only correct int values
				// but just in case here is an err check
				if err == nil {
					cfg.RunId = runId
				}
			}

			return cfg
		})
	},
}
