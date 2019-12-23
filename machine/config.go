package machine

import (
	"io"

	"go.uber.org/zap"
)

type Config struct {
	inReader  io.Reader
	logger    *zap.Logger
	outWriter io.Writer
}

type Option func(c *Config)

func InReader(r io.Reader) Option {
	return func(c *Config) {
		c.inReader = r
	}
}

func Logger(l *zap.Logger) Option {
	return func(c *Config) {
		c.logger.Sync()
		c.logger = l
	}
}

func OutWriter(w io.Writer) Option {
	return func(c *Config) {
		c.outWriter = w
	}
}
