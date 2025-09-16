package prettylogger

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func newTestLogger() (*slog.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	handler := NewHandler(&buf, HandlerOptions{
		NoColor: true, // Disable ANSI codes for testing
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	})
	logger := slog.New(handler)
	return logger, &buf
}

func TestBasicLogging(t *testing.T) {
	logger, buf := newTestLogger()
	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected output to contain 'test message', got: %s", output)
	}
	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected output to contain 'INFO', got: %s", output)
	}
}

func TestLogLevels(t *testing.T) {
	tests := []struct {
		name    string
		logFunc func(*slog.Logger, string, ...any)
		level   string
		message string
	}{
		{"debug", (*slog.Logger).Debug, "DEBUG", "debug message"},
		{"info", (*slog.Logger).Info, "INFO", "info message"},
		{"warn", (*slog.Logger).Warn, "WARN", "warning message"},
		{"error", (*slog.Logger).Error, "ERROR", "error message"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, buf := newTestLogger()
			tt.logFunc(logger, tt.message)

			output := buf.String()
			if !strings.Contains(output, tt.message) {
				t.Errorf("Expected output to contain '%s', got: %s", tt.message, output)
			}
			if !strings.Contains(output, tt.level) {
				t.Errorf("Expected output to contain '%s', got: %s", tt.level, output)
			}
		})
	}
}

func TestWithAttrs(t *testing.T) {
	logger, buf := newTestLogger()
	logger = logger.With("string", "value", "number", 42)
	logger.Info("message with attrs")

	output := buf.String()
	if !strings.Contains(output, "string=value") || !strings.Contains(output, "number=42") {
		t.Errorf("Expected output to contain attributes, got: %s", output)
	}
}

func TestWithError(t *testing.T) {
	logger, buf := newTestLogger()
	err := errors.New("test error")
	logger.Error("error occurred", "error", err)

	output := buf.String()
	if !strings.Contains(output, "test error") {
		t.Errorf("Expected output to contain 'test error', got: %s", output)
	}
}

func TestWithNestedData(t *testing.T) {
	logger, buf := newTestLogger()
	metadata := map[string]interface{}{
		"user": map[string]interface{}{
			"id":   123,
			"name": "Test User",
		},
	}
	logger.Info("user action", "metadata", metadata)

	output := buf.String()
	expected := "metadata={user={id=123 name=Test User}}"
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain nested pretty format '%s', got: %s", expected, output)
	}
}

func TestSourceLocation(t *testing.T) {
	var buf bytes.Buffer
	handler := NewHandler(&buf, HandlerOptions{
		NoColor: true,
		SlogOpts: slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		},
	})
	logger := slog.New(handler)
	logger.Info("with source")

	output := buf.String()
	if !strings.Contains(output, "logger_test.go") {
		t.Errorf("Expected output to contain source file info, got: %s", output)
	}
}

func TestHandleRecordDirectly(t *testing.T) {
	var buf bytes.Buffer
	handler := NewHandler(&buf, HandlerOptions{NoColor: true})

	record := slog.Record{
		Time:    time.Date(2025, 9, 16, 12, 0, 0, 0, time.UTC),
		Message: "direct record",
		Level:   slog.LevelInfo,
	}
	record.AddAttrs(slog.String("key", "value"))

	err := handler.Handle(context.Background(), record)
	if err != nil {
		t.Errorf("Handler.Handle returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "direct record") {
		t.Errorf("Expected output to contain 'direct record', got: %s", output)
	}
	if !strings.Contains(output, "key=value") {
		t.Errorf("Expected output to contain 'key=value', got: %s", output)
	}
}
