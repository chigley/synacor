package machine

import (
	"io"

	"go.uber.org/zap"
)

type Config struct {
	logger    *zap.Logger
	outWriter io.Writer
}

type Option func(c *Config)

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
