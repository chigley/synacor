package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

var (
	prg    *os.File
	logger *zap.Logger

	verbose bool
)

var rootCmd = &cobra.Command{
	Use: "synacor",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		logger, err = initLogger(verbose)
		if err != nil {
			return err
		}

		prg, err = os.Open(args[0])
		return err
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		prg.Close()
		logger.Sync()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func initLogger(verbose bool) (*zap.Logger, error) {
	if verbose {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}
