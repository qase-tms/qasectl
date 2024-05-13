package upload

import (
	"fmt"
	"github.com/qase-tms/qasectl/cmd/flags"
	"github.com/qase-tms/qasectl/internal/client"
	"github.com/qase-tms/qasectl/internal/parsers/junit"
	"github.com/qase-tms/qasectl/internal/service/result"
	"github.com/qase-tms/qasectl/internal/service/run"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	pathFlag        = "path"
	formatFlag      = "format"
	runIDFlag       = "id"
	titleFlag       = "title"
	descriptionFlag = "description"
)

// Command returns a new cobra command for upload
func Command() *cobra.Command {
	var (
		path        string
		format      string
		runID       int64
		title       string
		description string
	)

	cmd := &cobra.Command{
		Use:     "upload",
		Short:   "Upload test results",
		Example: "qli result upload --path 'path' --format 'junit' --id 123 --project 'PRJ' --token 'TOKEN'",
		RunE: func(cmd *cobra.Command, args []string) error {
			token := viper.GetString(flags.TokenFlag)
			project := viper.GetString(flags.ProjectFlag)

			var p result.Parser
			switch format {
			case "junit":
				p = junit.NewParser(path)
			case "qase":
				fmt.Println("Uploading Qase results")
			case "allure":
				fmt.Println("Uploading Allure results")
			case "xctest":
				fmt.Println("Uploading XTest results")
			default:
				return fmt.Errorf("unknown format: %s. allowed formats: junit, qase, allure, xctest", format)
			}

			c := client.NewClientV1(token)

			if runID == 0 {
				rs := run.NewService(c)

				id, err := rs.CreateRun(cmd.Context(), project, title, &description)
				if err != nil {
					return err
				}

				runID = id
			}

			s := result.NewService(c, p)

			err := s.Import(cmd.Context(), project, runID)
			if err != nil {
				return err
			}

			fmt.Println("Results uploaded successfully")
			return nil
		},
	}

	cmd.Flags().StringVar(&path, pathFlag, "", "path to the results file")
	err := cmd.MarkFlagRequired(pathFlag)
	if err != nil {
		fmt.Println(err)
	}

	cmd.Flags().StringVar(&format, formatFlag, "", "format of the results file: junit, qase, allure, xctest")
	err = cmd.MarkFlagRequired(formatFlag)
	if err != nil {
		fmt.Println(err)
	}

	cmd.Flags().Int64Var(&runID, runIDFlag, 0, "ID of the test run")
	cmd.Flags().StringVar(&title, titleFlag, "", "Title of the test run")
	cmd.Flags().StringVarP(&description, descriptionFlag, "d", "", "Description of the test run")
	cmd.MarkFlagsOneRequired(runIDFlag, titleFlag)
	cmd.MarkFlagsMutuallyExclusive(runIDFlag, titleFlag)

	return cmd
}
