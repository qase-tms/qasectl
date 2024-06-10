package create

import (
	"fmt"
	"github.com/qase-tms/qasectl/cmd/flags"
	"github.com/qase-tms/qasectl/internal/client"
	"github.com/qase-tms/qasectl/internal/service/env"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"path"
	"strings"
)

const (
	titleFlag       = "title"
	descriptionFlag = "description"
	slugFlag        = "slug"
	hostFlag        = "host"
	outputFlag      = "output"
)

// Command returns a new cobra command for create environments
func Command() *cobra.Command {
	var (
		title       string
		description string
		slug        string
		host        string
		output      string
	)

	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create a new environment",
		Example: "qli testops env create --title 'New environment' --slug local --description 'This is a environment' --host app.server.com --project 'PRJ' --token 'TOKEN'",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.Contains(slug, " ") {
				return fmt.Errorf("slug can't contain spaces")
			}

			token := viper.GetString(flags.TokenFlag)
			project := viper.GetString(flags.ProjectFlag)

			c := client.NewClientV1(token)
			s := env.NewService(c)

			e, err := s.CreateEnvironment(cmd.Context(), project, title, description, slug, host)
			if err != nil {
				return fmt.Errorf("failed to create environment: %w", err)
			}

			if output == "" {
				dir, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current directory: %w", err)
				}
				output = path.Join(dir, "qase.env")
			}

			err = os.WriteFile(output, []byte(fmt.Sprintf("QASE_ENVIRONMENT=%d", e.ID)), 0644)
			if err != nil {
				return fmt.Errorf("failed to write environament ID to file: %w", err)
			}

			slog.Info(fmt.Sprintf("Environment created with ID: %d", e.ID))

			return nil
		},
	}

	cmd.Flags().StringVar(&title, titleFlag, "", "title of the environment")
	err := cmd.MarkFlagRequired(titleFlag)
	if err != nil {
		fmt.Println(err)
	}
	cmd.Flags().StringVarP(&description, descriptionFlag, "d", "", "description of the environment")
	cmd.Flags().StringVarP(&slug, slugFlag, "s", "", "slug of the environment, (string without spaces)")
	err = cmd.MarkFlagRequired(slugFlag)
	if err != nil {
		fmt.Println(err)
	}
	cmd.Flags().StringVar(&host, hostFlag, "", "host of the environment")
	cmd.Flags().StringVarP(&output, outputFlag, "o", "", "output path for the environment ID")

	return cmd
}
