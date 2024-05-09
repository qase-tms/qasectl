package complete

import (
	"fmt"
	"github.com/qase-tms/qasectl/cmd/flags"
	"github.com/qase-tms/qasectl/internal/client"
	"github.com/qase-tms/qasectl/internal/service/run"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	idFlag = "id"
)

// Command returns a new cobra command for complete runs
func Command() *cobra.Command {
	var (
		runID int64
	)

	cmd := &cobra.Command{
		Use:     "complete",
		Short:   "Complete a test run",
		Example: "qli run complete --id 123 --project 'PRJ' --token 'TOKEN'",
		RunE: func(cmd *cobra.Command, args []string) error {
			token := viper.GetString(flags.TokenFlag)
			project := viper.GetString(flags.ProjectFlag)

			c := client.NewClientV1(token)
			s := run.NewService(c)

			err := s.CompleteRun(cmd.Context(), project, runID)
			if err != nil {
				if err.Error() == "404 Not Found" {
					fmt.Printf("Run %v not found", runID)
					return nil
				}
				return err
			}

			fmt.Printf("Run %v completed", runID)
			return err
		},
	}

	cmd.Flags().Int64Var(&runID, idFlag, 0, "ID of the test run")
	err := cmd.MarkFlagRequired(idFlag)
	if err != nil {
		fmt.Println(err)
	}

	return cmd
}
