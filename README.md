<img src="./.github/assets/gopher.png" align="right" height="96" width="96"/>

# Slog Pretty Logger

[![Tests](https://github.com/nentgroup/slog-prettylogger/actions/workflows/test.yml/badge.svg)](https://github.com/nentgroup/slog-prettylogger/actions/workflows/test.yml)
[![Go Coverage](https://github.com/nentgroup/slog-prettylogger/wiki/coverage.svg)](https://raw.githack.com/wiki/nentgroup/slog-prettylogger/coverage.html)
[![Go Reference](https://pkg.go.dev/badge/github.com/nentgroup/slog-prettylogger.svg)](https://pkg.go.dev/github.com/nentgroup/slog-prettylogger)
[![Go Report Card](https://goreportcard.com/badge/github.com/nentgroup/slog-prettylogger)](https://goreportcard.com/report/github.com/nentgroup/slog-prettylogger)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/rs/zerolog/master/LICENSE)

A colorful, readable logging handler for Go's standard `slog` package. Makes logs easier to read during development with color-coded levels and clean formatting. Inspired by zerolog's console output.

## Features

- 🎨 **Color-coded levels**: DEBUG, INFO, WARN, ERROR
- 🔍 **Bold messages**: Highlights log messages and errors
- 🧩 **Field formatting**: Formats fields based on type
- ⏱ **Custom time**: Override the default `time.Kitchen` format
- 🚫 **No-color mode**: Disable ANSI colors
- 🖌️ **Custom colors**: Set your own colors per log level


## Installation

```bash
go get github.com/nentgroup/slog-prettylogger
```

## Usage

```go
package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/nentgroup/slog-prettylogger"
)

func main() {
	// Create a new pretty logger
	logger := slog.New(prettylogger.NewHandler(os.Stdout, prettylogger.HandlerOptions{
		SlogOpts: slog.HandlerOptions{
			AddSource: true, // Include caller information
			Level:     slog.LevelDebug,
		},
		TimeFormat: time.TimeOnly, // Customize time format, default is time.Kitchen
	}))

	// Set as default logger
	slog.SetDefault(logger)

	// Basic logging
	slog.Info("application started", "version", "1.0.0")
	slog.Debug("debug information", "cache_hits", 42)
	slog.Warn("resource usage high", "cpu", 85.5, "memory", "3.2GB")
	slog.Error("failed to connect to database",
		"error", "connection refused",
		"retry_count", 3,
		"db_host", "localhost:5432")

	// With structured data
	logger.Info("user logged in",
		"user_id", 123,
		"metadata", map[string]interface{}{
			"browser":  "Chrome",
			"version":  "98.0.4758.102",
			"platform": "macOS",
		})
}
```

## Output Example

When using pretty logger, your console output will look similar to:

![Pretty Logger Output Example](./.github/assets/example.png)

## Configuration Options

Use the `HandlerOptions` struct to customize the logger:

```go
type HandlerOptions struct {
    SlogOpts    slog.HandlerOptions        // Standard slog handler options (level, AddSource, etc.)
    TimeFormat  string                     // Optional: custom time format (default is time.Kitchen)
    NoColor     bool                       // Optional: disable ANSI colors
    LevelColors map[slog.Level]string      // Optional: override colors per log level (DEBUG, INFO, WARN, ERROR)
}
```

### Example with custom colors

```go
logger := slog.New(prettylogger.NewHandler(os.Stdout, prettylogger.HandlerOptions{
	SlogOpts: slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	},
	TimeFormat:  time.RFC3339,
	LevelColors: map[slog.Level]string{slog.LevelError: prettylogger.BoldRed}, // or use any ANSI color code
}))
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgements

Developed and maintained by [Viaplay Group](https://github.com/nentgroup).
