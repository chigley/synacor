package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/chigley/synacor/machine"
	"go.uber.org/zap"
)

func main() {
	verbose := flag.Bool("verbose", false, "enable verbose output")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		usage()
	}
	cmd, bin := args[0], args[1]

	logger, err := logger(*verbose)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	prg, err := os.Open(bin)
	if err != nil {
		logger.Fatal("opening program", zap.Error(err))
	}
	defer prg.Close()

	switch cmd {
	case "run":
		if err := run(prg, logger); err != nil {
			logger.Fatal("run", zap.Error(err))
		}
	default:
		usage()
	}
}

func run(prg io.Reader, logger *zap.Logger) error {
	machine, err := machine.New(prg, machine.Logger(logger))
	if err != nil {
		return fmt.Errorf("creating machine: %w", err)
	}
	return fmt.Errorf("running machine: %w", machine.Run())
}

func logger(verbose bool) (*zap.Logger, error) {
	if verbose {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-verbose] run challenge.bin\n", os.Args[0])
	os.Exit(1)
}
