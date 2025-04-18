package logging_test

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/bilte-co/bilte/internal/logging"

	"github.com/rs/zerolog"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		level       string
		development bool
		expected    zerolog.Level
	}{
		{"debug", false, zerolog.DebugLevel},
		{"info", false, zerolog.InfoLevel},
		{"warn", false, zerolog.WarnLevel},
		{"error", false, zerolog.ErrorLevel},
		{"invalid", false, zerolog.InfoLevel}, // Default case
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			logger := logging.NewLogger(tt.level, tt.development)
			if logger.GetLevel() != tt.expected {
				t.Errorf("Expected level %v, got %v", tt.expected, logger.GetLevel())
			}
		})
	}
}

func TestNewLoggerFromEnv(t *testing.T) {
	// Backup original env variables
	origLogLevel, hasLogLevel := os.LookupEnv("LOG_LEVEL")
	origAppEnv, hasAppEnv := os.LookupEnv("APP_ENV")

	// Restore original environment variables after test
	defer func() {
		if hasLogLevel {
			t.Setenv("LOG_LEVEL", origLogLevel)
		} else {
			os.Unsetenv("LOG_LEVEL")
		}

		if hasAppEnv {
			t.Setenv("APP_ENV", origAppEnv)
		} else {
			os.Unsetenv("APP_ENV")
		}
	}()

	// Set test env variables
	t.Setenv("LOG_LEVEL", "warn")
	t.Setenv("APP_ENV", "development")

	logger := logging.NewLoggerFromEnv()
	if logger.GetLevel() != zerolog.WarnLevel {
		t.Errorf("Expected log level %v, got %v", zerolog.WarnLevel, logger.GetLevel())
	}

	// Test with production
	t.Setenv("APP_ENV", "production")
	logger = logging.NewLoggerFromEnv()
	if logger.GetLevel() != zerolog.WarnLevel {
		t.Errorf("Expected log level %v, got %v", zerolog.WarnLevel, logger.GetLevel())
	}
}

func TestDefaultLogger_Singleton(t *testing.T) {
	var wg sync.WaitGroup
	numRoutines := 10
	loggers := make([]*zerolog.Logger, numRoutines)

	// Ensure DefaultLogger() is only initialized once
	wg.Add(numRoutines)
	for i := range numRoutines {
		go func(i int) {
			defer wg.Done()
			loggers[i] = logging.DefaultLogger()
		}(i)
	}
	wg.Wait()

	for i := 1; i < numRoutines; i++ {
		if loggers[i] != loggers[0] {
			t.Fatalf("Expected DefaultLogger to be singleton, but found different instances")
		}
	}
}

func TestWithLoggerAndFromContext(t *testing.T) {
	ctx := context.Background()
	logger := logging.NewLogger("info", false)

	// Attach logger to context
	ctx = logging.WithLogger(ctx, logger)
	retrievedLogger := logging.FromContext(ctx)

	// Validate that the retrieved logger is the same as what we set
	if retrievedLogger != logger {
		t.Fatalf("FromContext did not return the same logger that was set using WithLogger")
	}

	// Test default behavior
	newCtx := context.Background()
	defaultLogger := logging.FromContext(newCtx)

	if defaultLogger != logging.DefaultLogger() {
		t.Fatalf("FromContext should return DefaultLogger when no logger is set in the context")
	}
}
