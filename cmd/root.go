package cmd

import (
	"github.com/qase-tms/qasectl/cmd/testops"
	"github.com/qase-tms/qasectl/cmd/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log/slog"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "qasectl",
	Short: "CLI tool for Qase TestOps",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(setLogger)
	viper.SetEnvPrefix("QASE_TESTOPS")
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	err := viper.BindPFlag("DEBUG", rootCmd.PersistentFlags().Lookup("verbose"))
	if err != nil {
		slog.Error("failed to bind flag", "flag", "verbose", "error", err)
	}

	rootCmd.AddCommand(testops.Command())
	rootCmd.AddCommand(version.VersionCmd())
}

func setLogger() {
	debug := viper.GetBool("Debug")

	ll := slog.LevelInfo
	if debug {
		ll = slog.LevelDebug
	}
	opts := &slog.HandlerOptions{
		Level: ll,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)

	slog.SetDefault(slog.New(handler))
}
