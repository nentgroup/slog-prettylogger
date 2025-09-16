<img src="./.github/assets/gopher.png" align="right" height="96" width="96"/>

# Slog Pretty Logger

[![Tests](https://github.com/nentgroup/slog-prettylogger/actions/workflows/test.yml/badge.svg)](https://github.com/nentgroup/slog-prettylogger/actions/workflows/test.yml)
[![Go Coverage](https://github.com/nentgroup/slog-prettylogger/wiki/coverage.svg)](https://raw.githack.com/wiki/nentgroup/slog-prettylogger/coverage.html)
[![Go Reference](https://pkg.go.dev/badge/github.com/nentgroup/slog-prettylogger.svg)](https://pkg.go.dev/github.com/nentgroup/slog-prettylogger)
[![Go Report Card](https://goreportcard.com/badge/github.com/nentgroup/slog-prettylogger)](https://goreportcard.com/report/github.com/nentgroup/slog-prettylogger)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/rs/zerolog/master/LICENSE)

A colorful, readable logging handler for Go's standard `slog` package. Makes logs easier to read during development with color-coded levels and clean formatting. Inspired by zerolog's console output.

## Features

- 🎨 **Color-coded log levels**: DEBUG, INFO, WARN, ERROR are visually distinguishable
- 🔍 **Emphasized messages**: Important parts stand out with bold formatting
- 🧩 **Smart field formatting**: Automatically formats fields based on type
- 📊 **JSON object support**: Properly formats nested JSON objects

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
	slog.Info("Application started", "version", "1.0.0")
	slog.Debug("Debug information", "cache_hits", 42)
	slog.Warn("Resource usage high", "cpu", 85.5, "memory", "3.2GB")
	slog.Error("Failed to connect to database",
		"error", "connection refused",
		"retry_count", 3,
		"db_host", "localhost:5432")

	// With structured data
	logger.Info("User logged in",
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

The `HandlerOptions` struct allows customization of the logger behavior:

```go
type HandlerOptions struct {
	SlogOpts slog.HandlerOptions // Standard slog handler options
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgements

Developed and maintained by [Viaplay Group](https://github.com/nentgroup).
