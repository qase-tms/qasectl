package cmd

import "github.com/spf13/cobra"

const (
	titleFlag       = "title"
	descriptionFlag = "description"
	pathFlag        = "path"
	runIdFlag       = "runId"
	projectCodeFlag = "projectCode"
)

func requireFlags(cmd *cobra.Command, flags ...string) {
	for _, flag := range flags {
		err := cmd.MarkFlagRequired(flag)
		if err != nil {
			panic(err) // shouldn't happen
		}
	}
}
