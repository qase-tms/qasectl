package create

import (
	"fmt"
	"github.com/qase-tms/qasectl/cmd/flags"
	"github.com/qase-tms/qasectl/internal/client"
	"github.com/qase-tms/qasectl/internal/service/milestone"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"path"
	"time"
)

const (
	titleFlag       = "title"
	descriptionFlag = "description"
	statusFlag      = "status"
	dueDateFlag     = "due-date"
	outputFlag      = "output"
)

// Command returns a new cobra command for create milestones
func Command() *cobra.Command {
	var (
		title       string
		description string
		status      string
		dueDate     string
		output      string
	)

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a new milestone",
		Example: "qli testops milestone create --title 'New milestone' --description 'This is a milestone' --status active --due-date 2024-01-30 --project 'PRJ' --token 'TOKEN'",
		RunE: func(cmd *cobra.Command, args []string) error {
			if status != "" && status != "active" && status != "completed" {
				return fmt.Errorf("invalid status value: %s. allowed values: active, completed", status)
			}

			t := int64(0)
			if dueDate != "" {
				d, err := time.Parse("2006-01-02", dueDate)
				if err != nil {
					return fmt.Errorf("failed to parse due date: %w", err)
				}
				t = d.Unix()
			}

			token := viper.GetString(flags.TokenFlag)
			project := viper.GetString(flags.ProjectFlag)

			c := client.NewClientV1(token)
			s := milestone.NewService(c)

			e, err := s.CreateMilestone(cmd.Context(), project, title, description, status, t)
			if err != nil {
				return fmt.Errorf("failed to create milestone: %w", err)
			}

			if output == "" {
				dir, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current directory: %w", err)
				}
				output = path.Join(dir, "qase.env")
			}

			err = os.WriteFile(output, []byte(fmt.Sprintf("QASE_MILESTONE=%d", e.ID)), 0644)
			if err != nil {
				return fmt.Errorf("failed to write milestone ID to file: %w", err)
			}

			slog.Info(fmt.Sprintf("Milestone created with ID: %d", e.ID))

			return nil
		},
	}

	cmd.Flags().StringVar(&title, titleFlag, "", "title of the milestone")
	err := cmd.MarkFlagRequired(titleFlag)
	if err != nil {
		fmt.Println(err)
	}
	cmd.Flags().StringVarP(&description, descriptionFlag, "d", "", "description of the milestone")
	cmd.Flags().StringVarP(&status, statusFlag, "s", "", "status of the milestone. Allowed values: active, completed")
	cmd.Flags().StringVar(&dueDate, dueDateFlag, "", "due date of the milestone. Format: YYYY-MM-DD")
	cmd.Flags().StringVarP(&output, outputFlag, "o", "", "output path for the milestone ID")

	return cmd
}
