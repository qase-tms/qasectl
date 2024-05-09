package cmd

import (
	"github.com/qase-tms/qasectl/internal"
	"github.com/spf13/cobra"
)

const pathFlag = "path"

//func init() {
//	importCmd.Flags().String(pathFlag, "", "")
//	importCmd.MarkFlagRequired(pathFlag)
//
//	rootCmd.AddCommand(importCmd)
//}

var importCmd = &cobra.Command{
	Use: "import",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cmd.Flag(pathFlag).Value.String()

		return internal.Import(path)
	},
}
