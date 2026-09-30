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
			fmt.Println(versionOutput(internal.ResolveVersion()))
		},
	}

	return versionCmd
}

// versionOutput renders the line printed by the version command.
func versionOutput(version string) string {
	if version == internal.UnknownVersion {
		return "Qase CLI: version not set during build."
	}

	return "Qase CLI " + version
}
