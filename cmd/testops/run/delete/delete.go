package delete

import (
	"fmt"
	"github.com/qase-tms/qasectl/cmd/flags"
	"github.com/qase-tms/qasectl/internal/client"
	"github.com/qase-tms/qasectl/internal/service/run"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"time"
)

const (
	idsFlag       = "ids"
	allFlag       = "all"
	startTimeFlag = "start"
	endTimeFlag   = "end"
)

// Command returns a new cobra command for delete runs
func Command() *cobra.Command {
	var (
		ids       []int64
		all       bool
		startTime string
		endTime   string
	)

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete test runs",
		Example: "qli testops run delete --ids 1,2,3 --start 2024-01-02 --end 2024-12-31 --project 'PRJ' --token 'TOKEN'",
		RunE: func(cmd *cobra.Command, args []string) error {
			token := viper.GetString(flags.TokenFlag)
			project := viper.GetString(flags.ProjectFlag)

			var start, end int64 = 0, 0

			if startTime != "" {
				t, err := time.Parse(time.DateOnly, startTime)
				if err != nil {
					return fmt.Errorf("failed to parse start time: %w", err)
				}
				start = t.Unix()
			}

			if endTime != "" {
				t, err := time.Parse(time.DateOnly, endTime)
				if err != nil {
					return fmt.Errorf("failed to parse end time: %w", err)
				}
				end = t.Unix()
			}

			c := client.NewClientV1(token)
			s := run.NewService(c)

			err := s.DeleteRun(cmd.Context(), project, ids, all, start, end)
			if err != nil {
				return fmt.Errorf("failed to delete test runs: %w", err)
			}

			slog.Info("Test runs deleted")

			return nil
		},
	}

	cmd.Flags().Int64SliceVar(&ids, idsFlag, []int64{}, "IDs of test runs to delete. format: --ids 1,2,3")
	cmd.Flags().BoolVar(&all, allFlag, false, "delete all test runs in the project")
	cmd.MarkFlagsOneRequired(idsFlag, allFlag)
	cmd.MarkFlagsMutuallyExclusive(idsFlag, allFlag)

	cmd.Flags().StringVarP(&startTime, startTimeFlag, "s", "", "start date of the test runs. Format: YYYY-MM-DD")
	cmd.Flags().StringVarP(&endTime, endTimeFlag, "e", "", "end date of the test runs. Format: YYYY-MM-DD")

	return cmd
}
