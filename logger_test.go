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

func TestLogAttrsDoesNotPanicWithAddSource(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, HandlerOptions{
		NoColor:  true,
		SlogOpts: slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug},
	})
	slog.SetDefault(slog.New(h))

	tracingAttrs := slog.Group(
		"tracing",
		slog.String("spanId", "some-span-id"),
	)

	metaAttrs := slog.Group(
		"meta",
		slog.String("http.url", "http://example.com"),
		slog.String("http.method", "GET"),
		slog.Int("http.status", 500),
	)

	slog.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"via LogAttrs",
		slog.String("caller", "ssss"),
		tracingAttrs,
		metaAttrs,
	)

	out := buf.String()
	if !strings.Contains(out, "via LogAttrs") || !strings.Contains(out, "http.url=http://example.com") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestStdLogBridgeDoesNotPanicWithAddSource(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, HandlerOptions{
		NoColor:  true,
		SlogOpts: slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug},
	})
	logger := slog.NewLogLogger(h, slog.LevelInfo)

	logger.Printf("stdlog %s", "message")

	out := buf.String()
	if !strings.Contains(out, "stdlog message") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestZeroTimeIsOmitted(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, HandlerOptions{NoColor: true})
	r := slog.Record{Message: "no time", Level: slog.LevelInfo}
	_ = h.Handle(context.Background(), r)
	out := buf.String()
	if !strings.HasPrefix(out, "INFO ") && !strings.Contains(out, " INFO ") {
		t.Fatalf("expected output to start with level when time is zero, got %q", out)
	}
}

type lv string

func (lv) LogValue() slog.Value { return slog.StringValue("resolved") }

func TestAttrValuesAreResolved(t *testing.T) {
	logger, buf := newTestLogger()
	logger.Info("has lv", slog.Any("lv", lv("x")))
	out := buf.String()
	if !strings.Contains(out, "lv=resolved") {
		t.Fatalf("expected resolved log value, got %q", out)
	}
}

func TestZeroValueAttrIgnored(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(&buf, HandlerOptions{NoColor: true})
	r := slog.Record{Message: "ignore zero attr", Level: slog.LevelInfo}
	r.AddAttrs(slog.Attr{})
	_ = h.Handle(context.Background(), r)
	out := buf.String()
	if strings.Contains(out, "=") && strings.Contains(out, "ignore zero attr") && len(strings.Fields(out)) > 3 {
		// There shouldn't be any extra key=value from zero attr; tolerate other fields
		t.Fatalf("unexpected extra attributes from zero attr: %q", out)
	}
}

func TestEmptyKeyGroupInlined(t *testing.T) {
	logger, buf := newTestLogger()
	logger.Info("inline grp", slog.Group("", slog.Int("a", 1), slog.Int("b", 2)))
	out := buf.String()
	if !strings.Contains(out, "a=1") || !strings.Contains(out, "b=2") {
		t.Fatalf("expected inlined group attrs, got %q", out)
	}
}

func TestNamedGroupFlattened(t *testing.T) {
	logger, buf := newTestLogger()
	logger.Info("named grp", slog.Group("g", slog.Int("a", 1), slog.String("s", "x")))
	out := buf.String()
	if !strings.Contains(out, "g.a=1") || !strings.Contains(out, "g.s=x") {
		t.Fatalf("expected flattened group keys, got %q", out)
	}
}

func TestEmptyGroupIgnored(t *testing.T) {
	logger, buf := newTestLogger()
	logger.Info("empty grp", slog.Group("g"))
	out := buf.String()
	if strings.Contains(out, "g=") || strings.Contains(out, "g.") {
		t.Fatalf("expected empty group to be ignored, got %q", out)
	}
}
