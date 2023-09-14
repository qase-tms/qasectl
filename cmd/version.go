package cmd

import (
	"fmt"

	"github.com/qase-tms/qasectl/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of this application",
	Long:  `All software has versions. This is this application's version.`,
	Run: func(cmd *cobra.Command, args []string) {
		if internal.Version == "" {
			fmt.Println("Qase CLI: version not set during build.")
		} else {
			fmt.Println("Qase CLI v" + internal.Version)
		}
	},
}
