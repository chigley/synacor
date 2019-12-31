package cmd

import (
	"fmt"
	"strings"

	"github.com/chigley/synacor/adventure"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search <challenge.bin>",
	Short: "Search a text adventure for the ruins",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path, err := adventure.FindRuins(prg, logger)
		if err != nil {
			logger.Fatal("finding ruins", zap.Error(err))
		}
		for _, n := range path {
			fmt.Printf("%s (%s)\n", n.ExitToHere, strings.Join(n.Inv, ", "))
		}
	},
}
