// Package prettylogger provides a colorized, human-readable logging handler for slog.
package prettylogger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"sort"
	"strings"
	"time"
)

// ANSI color codes
const (
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Magenta   = "\033[35m"
	Cyan      = "\033[36m"
	White     = "\033[37m"
	BoldRed   = "\033[1;31m"
	BoldWhite = "\033[1;37m"
)

const defaultTimeFormat = time.Kitchen

// HandlerOptions configures the behavior of a pretty logger Handler.
type HandlerOptions struct {
	SlogOpts    slog.HandlerOptions
	TimeFormat  string
	NoColor     bool
	LevelColors map[slog.Level]string
}

// Handler implements slog.Handler with pretty, colorized output formatting.
type Handler struct {
	slog.Handler
	opts  HandlerOptions
	l     *log.Logger
	attrs []slog.Attr
}

// NewHandler creates a new Handler that writes colorized, human-readable logs to out.
// It accepts HandlerOptions for configuration.
func NewHandler(out io.Writer, opts HandlerOptions) *Handler {
	if opts.TimeFormat == "" {
		opts.TimeFormat = defaultTimeFormat
	}
	if opts.LevelColors == nil {
		opts.LevelColors = map[slog.Level]string{
			slog.LevelDebug: Magenta,
			slog.LevelInfo:  Blue,
			slog.LevelWarn:  Yellow,
			slog.LevelError: Red,
		}
	}
	return &Handler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		opts:    opts,
		l:       log.New(out, "", 0),
	}
}

// WithAttrs returns a new Handler with merged attributes for pretty printing.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	newAttrs = append(newAttrs, h.attrs...)
	newAttrs = append(newAttrs, attrs...)
	return &Handler{
		Handler: h.Handler,
		l:       h.l,
		opts:    h.opts,
		attrs:   newAttrs,
	}
}

// WithGroup returns a new Handler (groups are ignored for pretty printing).
func (h *Handler) WithGroup(_ string) slog.Handler {
	return &Handler{
		Handler: h.Handler,
		l:       h.l,
		opts:    h.opts,
		attrs:   h.attrs,
	}
}

// Handle processes a log record, formatting it with colors and proper spacing.
func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String()
	if !h.opts.NoColor {
		if color, ok := h.opts.LevelColors[r.Level]; ok {
			level = colorize(color, level)
		}
	}

	fields := map[string]any{}
	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	timeStr := r.Time.Format(h.opts.TimeFormat)
	var parts []string
	parts = append(parts, timeStr, level)

	if h.opts.SlogOpts.AddSource {
		src := fmt.Sprintf("%s:%d", r.Source().File, r.Source().Line)
		parts = append(parts, colorize(BoldWhite, src))
	}
	parts = append(parts, ">", colorize(BoldWhite, r.Message))

	// Stable key order
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := fields[k]
		parts = append(parts, h.formatField(k, v))
	}

	h.l.Println(strings.Join(parts, " "))
	return nil
}

func (h *Handler) formatField(k string, v any) string {
	// Format key
	key := k + "="
	if !h.opts.NoColor {
		key = colorize(Cyan, key)
	}

	// Format value
	val := formatValue(v)
	if !h.opts.NoColor {
		if k == "error" {
			val = colorize(BoldRed, val)
		} else {
			val = colorize(White, val)
		}
	}

	return key + val
}

// formatValue formats a value for pretty printing.
func formatValue(v interface{}) string {
	switch vv := v.(type) {
	case string:
		return vv
	case int, int32, int64, float32, float64:
		return fmt.Sprintf("%v", vv)
	case error:
		return vv.Error()
	case map[string]interface{}:
		// Sort keys for consistent order
		keys := make([]string, 0, len(vv))
		for k := range vv {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		parts := make([]string, 0, len(vv))
		for _, k := range keys {
			parts = append(parts, k+"="+formatValue(vv[k]))
		}
		return "{" + strings.Join(parts, " ") + "}"
	case []slog.Attr:
		parts := make([]string, 0, len(vv))
		for _, attr := range vv {
			parts = append(parts, attr.Key+"="+formatValue(attr.Value.Any()))
		}
		return "{" + strings.Join(parts, " ") + "}"
	case []any:
		parts := make([]string, 0, len(vv))
		for _, item := range vv {
			parts = append(parts, formatValue(item))
		}
		return "{" + strings.Join(parts, " ") + "}"
	default:
		return fmt.Sprintf("%v", vv)
	}
}

func colorize(color, s string) string {
	return color + s + Reset
}
