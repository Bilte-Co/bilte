// heavily inspired by: https://github.com/google/exposure-notifications-server/blob/main/pkg/logging/logger.go

// Sets up / configures global logging for the application.
package logging

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
)

// contextKey is a private string type to prevent collisions in the context map.
type contextKey string

// loggerKey points to the value in the context where the logger is stored.
const loggerKey = contextKey("logger")

var (
	// defaultLogger is the default logger. It is initialized once per package
	// include upon calling DefaultLogger.
	defaultLogger     *slog.Logger
	defaultLoggerOnce sync.Once
)

func NewLogger(level string, development bool) *slog.Logger {
	w := os.Stderr
	options := &tint.Options{
		TimeFormat: time.RFC3339,
	}

	switch level {
	case "debug":
		options.Level = slog.LevelDebug
	case "info":
		options.Level = slog.LevelInfo
	case "warn":
		options.Level = slog.LevelWarn
	case "error":
		options.Level = slog.LevelError
	default:
		options.Level = slog.LevelInfo
	}

	logger := slog.New(tint.NewHandler(w, options))

	return logger
}

// NewLoggerFromEnv creates a new logger from the environment. It consumes
// LOG_LEVEL for determining the level and LOG_MODE for determining the output
// parameters.
func NewLoggerFromEnv() *slog.Logger {
	_ = godotenv.Load()

	level := os.Getenv("LOG_LEVEL")
	development := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV"))) == "development"

	return NewLogger(level, development)
}

// DefaultLogger returns the default logger for the package.
func DefaultLogger() *slog.Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLoggerFromEnv()
	})
	return defaultLogger
}

// WithLogger creates a new context with the provided logger attached.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns the logger stored in the context. If no such logger
// exists, a default logger is returned.
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return DefaultLogger()
}
