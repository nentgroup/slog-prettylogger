//go:build ignoretests

package main

import (
	"log/slog"
	"os"

	prettylogger "github.com/nentgroup/slog-prettylogger"
)

func main() {
	handler := prettylogger.NewHandler(os.Stdout, prettylogger.HandlerOptions{})
	logger := slog.New(handler).With("app", "myservice", "env", "dev")

	logger.Info("service started")
	logger.Warn("resource limit approaching", "cpu", 92)
}
