package create

import (
	"fmt"
	"github.com/qase-tms/qasectl/cmd/flags"
	"github.com/qase-tms/qasectl/internal/client"
	"github.com/qase-tms/qasectl/internal/service/run"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"path"
)

const (
	titleFlag       = "title"
	descriptionFlag = "description"
	environmentFlag = "environment"
	milestoneFlag   = "milestone"
	planFlag        = "plan"
	outputFlag      = "output"
)

// Command returns a new cobra command for create runs
func Command() *cobra.Command {
	var (
		title       string
		description string
		environment string
		milestone   string
		plan        string
		output      string
	)

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a new test run",
		Example: "qli run create --title 'My test run' --description 'This is a test run' --environment local --project 'PRJ' --token 'TOKEN'",
		RunE: func(cmd *cobra.Command, args []string) error {
			token := viper.GetString(flags.TokenFlag)
			project := viper.GetString(flags.ProjectFlag)

			c := client.NewClientV1(token)
			s := run.NewService(c)

			id, err := s.CreateRun(cmd.Context(), project, title, description, environment, milestone, plan)
			if err != nil {
				return fmt.Errorf("failed to create run: %w", err)
			}

			if output == "" {
				dir, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current directory: %w", err)
				}
				output = path.Join(dir, "qase.env")
			}

			err = os.WriteFile(output, []byte(fmt.Sprintf("QASE_TESTOPS_RUN_ID=%d", id)), 0644)
			if err != nil {
				return fmt.Errorf("failed to write run ID to file: %w", err)
			}

			slog.Info(fmt.Sprintf("Run created with ID: %d", id))

			return nil
		},
	}

	cmd.Flags().StringVarP(&title, titleFlag, "", "", "title of the test run")
	err := cmd.MarkFlagRequired(titleFlag)
	if err != nil {
		fmt.Println(err)
	}
	cmd.Flags().StringVarP(&description, descriptionFlag, "d", "", "description of the test run")
	cmd.Flags().StringVarP(&environment, environmentFlag, "e", "", "environment of the test run")
	cmd.Flags().StringVarP(&milestone, milestoneFlag, "m", "", "milestone of the test run")
	cmd.Flags().StringVar(&plan, planFlag, "", "plan of the test run")
	cmd.Flags().StringVarP(&output, outputFlag, "o", "", "output path for the test run ID")

	return cmd
}
