package filter

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/qase-tms/qasectl/cmd/flags"
	"github.com/qase-tms/qasectl/internal/client"
	"github.com/qase-tms/qasectl/internal/service/filter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	frameworkFlag = "framework"
	planIDFlag    = "planID"
	outputFlag    = "output"
)

// Command returns a new cobra command for upload
func Command() *cobra.Command {
	var (
		framework string
		planID    int64
		output    string
	)

	cmd := &cobra.Command{
		Use:     "filter",
		Short:   "Get filtered results for the given plan ID and framework",
		Example: "qasectl testops filter --framework 'playwright' --planID 123 --project 'PRJ' --token 'TOKEN'",
		RunE: func(cmd *cobra.Command, args []string) error {
			const op = "filter"
			logger := slog.With("op", op)

			token := viper.GetString(flags.TokenFlag)
			project := viper.GetString(flags.ProjectFlag)

			cv1 := client.NewClientV1(token)
			s := filter.NewService(cv1)

			filteredResults, err := s.GetFilteredResults(cmd.Context(), project, planID, framework)
			if err != nil {
				return err
			}

			if output == "" {
				dir, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current directory: %w", err)
				}
				output = path.Join(dir, "qase.env")
			}

			err = os.WriteFile(output, []byte(fmt.Sprintf("QASE_FILTERED_RESULTS=%s", filteredResults)), 0644)
			if err != nil {
				return fmt.Errorf("failed to write filtered results to file: %w", err)
			}

			logger.Info(fmt.Sprintf("Filtered results saved to %s", output))

			return nil
		},
	}

	cmd.Flags().StringVarP(&framework, frameworkFlag, "f", "", "framework of the results file: playwright")
	err := cmd.MarkFlagRequired(frameworkFlag)
	if err != nil {
		slog.Error("Error while marking flag as required", "error", err)
	}

	cmd.Flags().Int64Var(&planID, planIDFlag, 0, "ID of the test plan")
	err = cmd.MarkFlagRequired(planIDFlag)
	if err != nil {
		slog.Error("Error while marking flag as required", "error", err)
	}

	return cmd
}
