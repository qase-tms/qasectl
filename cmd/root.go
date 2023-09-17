package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:     "qasectl",
	Example: "qasectl run create",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	rootCmd.SetUsageTemplate(rootUsageTemplate())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func rootUsageTemplate() string {
	return `USAGE{{if .HasAvailableSubCommands}}
  {{.CommandPath}} <command> <subcommand> [flags]{{end}}{{if .HasAvailableSubCommands}}

CORE COMMANDS{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

FLAGS
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasExample}}

EXAMPLES
{{.Example}}{{end}}

LEARN MORE
  Use 'qasectl <command> <subcommand> --help' for more information about a command.
  Read the manual at https://developers.qase.io/cli
`
}
