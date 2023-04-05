package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/qase-tms/qasectl/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(authCmd)
}

var authCmd = &cobra.Command{
	Use: "auth",
	RunE: func(cmd *cobra.Command, args []string) error {
		var token string
		prompt := &survey.Password{
			Message: "Please type your token:",
		}
		survey.AskOne(prompt, &token)

		return internal.UpdateToken(token)
	},
}
