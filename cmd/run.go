package cmd

import (
	"github.com/chigley/synacor/machine"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run <challenge.bin>",
	Short: "Run a program interactively",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := machine.New(prg, machine.Logger(logger))
		if err != nil {
			logger.Fatal("creating machine", zap.Error(err))
		}
		if err := m.Run(); err != nil {
			logger.Fatal("running machine", zap.Error(err))
		}
	},
}
