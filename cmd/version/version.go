package version

import (
	"fmt"
	"github.com/qase-tms/qasectl/internal"
	"github.com/spf13/cobra"
)

func VersionCmd() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:                "version",
		Short:              "Print the version number of this application",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			if internal.Version == "" {
				fmt.Println("Qase CLI: version not set during build.")
			} else {
				fmt.Println("Qase CLI " + internal.Version)
			}
		},
	}

	return versionCmd
}
