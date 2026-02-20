package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config holds configuration for the logger
type Config struct {
	EnableConsole bool
	ConsoleJSON   bool
	Verbose       bool
}

// Init initializes the global logger
func Init(cfg Config) {
	var output io.Writer = os.Stdout

	if cfg.EnableConsole {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	zerolog.TimeFieldFormat = time.RFC3339
	l := zerolog.New(output).With().Timestamp().Logger()

	if cfg.Verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Logger = l
}

// Log returns the global logger
func Log() *zerolog.Logger {
	return &log.Logger
}

// FromContext returns the logger associated with the context.
// If no logger is associated, it returns the global logger.
func FromContext(ctx context.Context) *zerolog.Logger {
	logger := zerolog.Ctx(ctx)
	if logger.GetLevel() == zerolog.Disabled {
		return &log.Logger
	}
	return logger
}

// WithContext returns a context with the logger associated.
func WithContext(ctx context.Context, logger *zerolog.Logger) context.Context {
	return logger.WithContext(ctx)
}
