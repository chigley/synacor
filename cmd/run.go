package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/chigley/synacor/machine"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&saveFile, "save-file", "s", "", "save final machine state to disk")
}

var saveFile string

var runCmd = &cobra.Command{
	Use:   "run <challenge.bin>",
	Short: "Run a program interactively",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m, err := machine.New(prg, machine.Logger(logger))
		if err != nil {
			logger.Fatal("creating machine", zap.Error(err))
		}
		if err := m.Run(); err != nil && !errors.Is(err, machine.ErrNeedInput) {
			logger.Fatal("running machine", zap.Error(err))
		}

		if saveFile != "" {
			if err := save(m, saveFile); err != nil {
				logger.Fatal("saving machine", zap.Error(err))
			}
		}
	},
}

func save(m *machine.Machine, path string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating save file: %w", err)
	}
	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
	}()
	return m.Encode(f)
}
