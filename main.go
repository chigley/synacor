package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/chigley/synacor/machine"
	"go.uber.org/zap"
)

func main() {
	verbose := flag.Bool("verbose", false, "enable verbose output")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-verbose] challenge.bin\n", os.Args[0])
		os.Exit(1)
	}

	logger, err := logger(*verbose)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	path := args[0]
	prg, err := os.Open(path)
	if err != nil {
		logger.Fatal("opening program", zap.Error(err))
	}
	defer prg.Close()

	machine, err := machine.New(prg, machine.Logger(logger))
	if err != nil {
		logger.Fatal("creating machine", zap.Error(err))
	}

	if err := machine.Run(); err != nil {
		logger.Fatal("running machine", zap.Error(err))
	}
}

func logger(verbose bool) (*zap.Logger, error) {
	if verbose {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}
