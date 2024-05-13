package create

import (
	"fmt"
	"github.com/qase-tms/qasectl/cmd/flags"
	"github.com/qase-tms/qasectl/internal/client"
	"github.com/qase-tms/qasectl/internal/service/run"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	titleFlag       = "title"
	descriptionFlag = "description"
)

// Command returns a new cobra command for create runs
func Command() *cobra.Command {
	var (
		title       string
		description string
	)

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a new test run",
		Example: "qli run create --title 'My test run' --description 'This is a test run' --project 'PRJ' --token 'TOKEN'",
		RunE: func(cmd *cobra.Command, args []string) error {
			token := viper.GetString(flags.TokenFlag)
			project := viper.GetString(flags.ProjectFlag)

			c := client.NewClientV1(token)
			s := run.NewService(c)

			id, err := s.CreateRun(cmd.Context(), project, title, &description)
			if err != nil {
				return err
			}

			fmt.Println("Run created with ID:", id)

			return nil
		},
	}

	cmd.Flags().StringVarP(&title, titleFlag, "", "", "title of the test run")
	err := cmd.MarkFlagRequired(titleFlag)
	if err != nil {
		fmt.Println(err)
	}
	cmd.Flags().StringVarP(&description, descriptionFlag, "d", "", "description of the test run")

	return cmd
}