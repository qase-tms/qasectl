package version

import (
	"fmt"
	"github.com/qase-tms/qasectl/internal"
	"github.com/spf13/cobra"
	"runtime/debug"
)

func VersionCmd() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:                "version",
		Short:              "Print the version number of this application",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			if internal.Version == "" {
				info, ok := debug.ReadBuildInfo()
				if ok && info.Main.Version != "(devel)" {
					fmt.Println("Qase CLI " + info.Main.Version)
					return
				}
				fmt.Println("Qase CLI: version not set during build.")
			} else {
				fmt.Println("Qase CLI " + internal.Version)
			}
		},
	}

	return versionCmd
}
