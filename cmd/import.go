package cmd

import (
	"github.com/qase-tms/qasectl/internal"
	"github.com/spf13/cobra"
)

func init() {
	importCmd.Flags().String(pathFlag, "", "")
	requireFlags(importCmd, pathFlag)

	importCmd.Flags().String(runIdFlag, "", "")
	requireFlags(importCmd, runIdFlag)

	importCmd.Flags().String(projectCodeFlag, "", "")
	requireFlags(importCmd, projectCodeFlag)

	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use: "import",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cmd.Flag(pathFlag).Value.String()

		return internal.Import(path)
	},
}
