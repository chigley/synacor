package cmd

import (
	"errors"

	"github.com/chigley/synacor/machine"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "save final machine state to disk")
}

var overwrite bool

var loadCmd = &cobra.Command{
	Use:   "load <save-file>",
	Short: "Run a program interactively, loading a save file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := machine.Load(prg, machine.Logger(logger))
		if err != nil {
			logger.Fatal("loading machine", zap.Error(err))
		}
		if err := m.Run(); err != nil && !errors.Is(err, machine.ErrNeedInput) {
			logger.Fatal("running machine", zap.Error(err))
		}

		prg.Close()
		if overwrite {
			if err := save(m, args[0]); err != nil {
				logger.Fatal("saving machine", zap.Error(err))
			}
		}
	},
}
