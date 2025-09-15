package testops

import (
	"fmt"

	"github.com/qase-tms/qasectl/cmd/flags"
	"github.com/qase-tms/qasectl/cmd/testops/env"
	"github.com/qase-tms/qasectl/cmd/testops/field"
	"github.com/qase-tms/qasectl/cmd/testops/filter"
	"github.com/qase-tms/qasectl/cmd/testops/milestone"
	"github.com/qase-tms/qasectl/cmd/testops/result"
	"github.com/qase-tms/qasectl/cmd/testops/run"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	tokenFlag   = "token"
	projectFlag = "project"
)

// Command returns a new cobra command for testops
func Command() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "testops",
		Short: "Manage test operations",
	}

	cmd.PersistentFlags().StringP(tokenFlag, "t", "", "token for Qase API access")
	err := viper.BindPFlag(flags.TokenFlag, cmd.PersistentFlags().Lookup(tokenFlag))
	if err != nil {
		fmt.Println(err)
	}
	err = cmd.MarkPersistentFlagRequired(tokenFlag)
	if err != nil {
		fmt.Println(err)
	}

	cmd.PersistentFlags().StringP(projectFlag, "p", "", "project code for Qase API access")
	err = viper.BindPFlag(flags.ProjectFlag, cmd.PersistentFlags().Lookup(projectFlag))
	if err != nil {
		fmt.Println(err)
	}
	err = cmd.MarkPersistentFlagRequired(projectFlag)
	if err != nil {
		fmt.Println(err)
	}

	cmd.AddCommand(run.Command())
	cmd.AddCommand(result.Command())
	cmd.AddCommand(env.Command())
	cmd.AddCommand(milestone.Command())
	cmd.AddCommand(filter.Command())
	cmd.AddCommand(field.Command())

	return cmd
}
