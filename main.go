package main

import (
	"flag"
	"log"
	"os"

	"github.com/chigley/synacor/machine"
	"go.uber.org/zap"
)

func main() {
	verbose := flag.Bool("verbose", false, "enable verbose output")
	flag.Parse()

	logger, err := logger(*verbose)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	machine, err := machine.New(os.Stdin, machine.Logger(logger))
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
