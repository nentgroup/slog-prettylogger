//go:build ignoretests

package main

import (
	"log/slog"
	"os"

	prettylogger "github.com/nentgroup/slog-prettylogger"
)

func main() {
	handler := prettylogger.NewHandler(os.Stdout, prettylogger.HandlerOptions{
		SlogOpts: slog.HandlerOptions{
			AddSource: true, // Include caller information
			Level:     slog.LevelDebug,
		},
		LevelColors: map[slog.Level]string{
			slog.LevelDebug: "\033[36m",   // cyan
			slog.LevelInfo:  "\033[32m",   // green
			slog.LevelWarn:  "\033[33m",   // yellow
			slog.LevelError: "\033[1;35m", // bold magenta
		},
	})
	logger := slog.New(handler)

	logger.Debug("debug with cyan")
	logger.Error("error with bold magenta")
}
