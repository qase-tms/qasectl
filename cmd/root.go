package cmd

import (
	"github.com/qase-tms/qasectl/cmd/testops"
	"github.com/qase-tms/qasectl/cmd/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	rootCmd.AddCommand(testops.Command())
	rootCmd.AddCommand(version.VersionCmd())
}
