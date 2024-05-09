package cmd

import (
	"fmt"
	"github.com/qase-tms/qasectl/cmd/flags"
	"github.com/qase-tms/qasectl/cmd/run"
	"github.com/qase-tms/qasectl/cmd/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	tokenFlag   = "token"
	projectFlag = "project"
)

var rootCmd = &cobra.Command{
	Use:   "qli",
	Short: "CLI tool for Qase TestOps",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	viper.SetEnvPrefix("QASE_TESTOPS")
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().StringP(tokenFlag, "t", "", "token for Qase API access")
	err := viper.BindPFlag(flags.TokenFlag, rootCmd.PersistentFlags().Lookup(tokenFlag))
	if err != nil {
		fmt.Println(err)
	}
	err = rootCmd.MarkPersistentFlagRequired(tokenFlag)
	if err != nil {
		fmt.Println(err)
	}

	rootCmd.PersistentFlags().StringP(projectFlag, "p", "", "project code for Qase API access")
	err = viper.BindPFlag(flags.ProjectFlag, rootCmd.PersistentFlags().Lookup(projectFlag))
	if err != nil {
		fmt.Println(err)
	}
	err = rootCmd.MarkPersistentFlagRequired(projectFlag)
	if err != nil {
		fmt.Println(err)
	}

	rootCmd.AddCommand(run.Command())
	rootCmd.AddCommand(version.VersionCmd())
}
