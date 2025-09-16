//go:build ignoretests

package main

import (
	"log/slog"
	"os"

	prettylogger "github.com/nentgroup/slog-prettylogger"
)

func main() {
	// Create a new pretty logger
	logger := slog.New(prettylogger.NewHandler(os.Stdout, prettylogger.HandlerOptions{
		SlogOpts: slog.HandlerOptions{
			AddSource: true, // Include caller information
			Level:     slog.LevelDebug,
		},
	}))

	// Set as default logger
	slog.SetDefault(logger)

	// Demonstrate different log levels
	slog.Debug("this is a debug message", "detail", "some debug info")
	slog.Info("application started", "version", "1.0.0", "env", "development")
	slog.Warn("resource usage high", "cpu", 85.5, "memory", "3.2GB")
	slog.Error("failed to connect to database",
		"error", "connection refused",
		"retry_count", 3,
		"db_host", "localhost:5432")

	// Log with nested structured data
	logger.Info("User logged in",
		"user_id", 123,
		"session", map[string]interface{}{
			"id": "abc-123-xyz",
			"ip": "192.168.1.1",
			"user_agent": map[string]interface{}{
				"browser":  "Chrome",
				"version":  "98.0.4758.102",
				"platform": "macOS",
			},
		})
}
