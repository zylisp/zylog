# zylog

[![Build Status][gh-actions-badge]][gh-actions]
[![Tags][github-tags-badge]][github-tags]

*A flexible logging wrapper supporting both slog and logrus with beautiful, consistent formatting*

## Features

- 🎨 Beautifully coloured output with customizable colours
- 🔄 Support for both **slog** (Go standard library) and **logrus**
- 📍 Optional caller information (package, function, line number)
- ⏰ Multiple timestamp formats (RFC3339, Standard, Simple, Time-only)
- 📏 Configurable level padding for perfect alignment
- 🎯 Unified API for both logging backends
- ⚙️ Simple, minimal configuration

## Installation

```bash
go get github.com/zylisp/zylog
```

## Example Use

Run the comprehensive demo:

```bash
make demo
```

Or build and run manually:

```bash
make build
./bin/zylog-demo
```

At which point you should see something like the following:

![screenshot](assets/images/screenshot.png)

See the demo code in  `cmd/zylog-demo/main.go` for examples of both logging backends.

## Quick Start

### Unified Setup (Recommended)

Use the unified `SetupLogging()` function to configure either logger:

```go
package main

import (
    "github.com/zylisp/zylog"
    "github.com/zylisp/zylog/formatter"
    "github.com/zylisp/zylog/options"
)

func main() {
    // Configure with slog (default)
    logger, err := zylog.SetupLogging(&options.ZyLog{
        Coloured:         true,
        Level:           "info",
        Output:          "stdout",
        ReportCaller:    true,
        TimestampFormat: formatter.SimpleTimestamp,
        PadLevel:        true,
        PadAmount:       7,
        PadSide:         "left",
        Logger:          options.Slog, // or options.LogRUs
    })
    if err != nil {
        panic(err)
    }

    // For slog, use the returned logger
    if logger != nil {
        logger.Info("Application started")
        logger.Info("User logged in", "user_id", 12345, "ip", "192.168.1.1")
    }
}
```

### Using with slog

```go
import (
    "log/slog"
    "github.com/zylisp/zylog"
    "github.com/zylisp/zylog/options"
)

func main() {
    // SetupLogging returns the configured slog logger and sets it as default
    logger, _ := zylog.SetupLogging(options.Default())

    // Use the returned logger directly
    logger.Info("Application started")

    // Or use the default slog logger
    slog.Info("This also works")

    // Structured logging
    logger.Info("User action",
        slog.String("user", "alice"),
        slog.Int("request_id", 12345))
}
```

### Using with logrus

```go
import (
    log "github.com/sirupsen/logrus"
    "github.com/zylisp/zylog"
    "github.com/zylisp/zylog/options"
)

func main() {
    // For logrus, SetupLogging returns nil (uses global instance)
    opts := options.Default()
    opts.Logger = options.LogRUs
    _, _ = zylog.SetupLogging(opts)

    // Use global logrus instance
    log.Info("Application started")
    log.WithFields(log.Fields{
        "user": "alice",
        "request_id": 12345,
    }).Info("User action")
}
```

## Configuration Options

```go
type ZyLog struct {
    Coloured         bool                // Enable coloured output
    Level           string              // Log level: "trace", "debug", "info", "warn", "error"
    Output          string              // Output destination: "stdout", "stderr"
    ReportCaller    bool                // Include caller information
    TimestampFormat formatter.TSFormat  // Timestamp format
    PadLevel        bool                // Pad level strings for alignment
    PadAmount       int                 // Number of characters to pad to
    PadSide         string              // "left" or "right"
    MsgSeparator    string              // Separator between message and attributes
    Logger          options.Logger      // LogRUs or Slog
}
```

### Preset Configurations

```go
// Default configuration (slog with all features enabled)
opts := options.Default()

// Disable caller reporting
opts := options.NoCaller()
```

## Colour Customization

Zylog allows you to customize the foreground and background colours of every formatted element. By default, zylog uses sensible colour defaults, but you can override any colour you want.

### Simple Example - Changing a Few Colours

```go
import (
    "github.com/fatih/color"
    "github.com/zylisp/zylog"
    "github.com/zylisp/zylog/colours"
    "github.com/zylisp/zylog/options"
)

func main() {
    opts := options.Default()

    // Customize just the colours you want to change
    opts.Colours.LevelError = &colours.Colour{
        Fg: color.FgHiRed,
        Bg: color.BgYellow,  // Add yellow background to errors
    }
    opts.Colours.Message = &colours.Colour{
        Fg: color.FgHiWhite,
        Bg: color.Reset,  // No background
    }

    logger, _ := zylog.SetupLogging(opts)
    logger.Error("This error has a yellow background!")
}
```

### Disabling Colour for Specific Elements

To disable colour for a specific element while keeping others coloured, set both Fg and Bg to `color.Reset`:

```go
opts := options.Default()
opts.Colours.Timestamp = &colours.Colour{
    Fg: color.Reset,
    Bg: color.Reset,
}
// Timestamp will now be uncoloured, but everything else remains coloured
```

### Complete Colour Configuration Reference

The `Colours` struct provides fine-grained control over every coloured element:

```go
type Colours struct {
    // Timestamp colours (default: HiBlack/grey)
    Timestamp *Colour

    // Log level colours
    LevelTrace   *Colour  // default: HiMagenta
    LevelDebug   *Colour  // default: HiCyan
    LevelInfo    *Colour  // default: HiGreen
    LevelWarn    *Colour  // default: HiYellow
    LevelWarning *Colour  // default: HiYellow
    LevelError   *Colour  // default: Red
    LevelFatal   *Colour  // default: HiRed
    LevelPanic   *Colour  // default: HiWhite

    // Message text colour (default: Green)
    Message *Colour

    // Arrow separator " ▶ " (default: Cyan)
    Arrow *Colour

    // Caller information colours
    CallerFunction *Colour  // default: HiYellow
    CallerLine     *Colour  // default: Yellow

    // Structured logging attribute colours
    AttrKey   *Colour  // default: Yellow
    AttrValue *Colour  // default: HiYellow
}

type Colour struct {
    Fg color.Attribute  // Foreground colour from github.com/fatih/color
    Bg color.Attribute  // Background colour from github.com/fatih/color
}
```

Available colour attributes from `github.com/fatih/color`:

**Foreground colours:**

- `color.FgBlack`, `color.FgRed`, `color.FgGreen`, `color.FgYellow`
- `color.FgBlue`, `color.FgMagenta`, `color.FgCyan`, `color.FgWhite`
- `color.FgHiBlack`, `color.FgHiRed`, `color.FgHiGreen`, `color.FgHiYellow`
- `color.FgHiBlue`, `color.FgHiMagenta`, `color.FgHiCyan`, `color.FgHiWhite`

**Background colours:**

- `color.BgBlack`, `color.BgRed`, `color.BgGreen`, `color.BgYellow`
- `color.BgBlue`, `color.BgMagenta`, `color.BgCyan`, `color.BgWhite`
- `color.BgHiBlack`, `color.BgHiRed`, `color.BgHiGreen`, `color.BgHiYellow`
- `color.BgHiBlue`, `color.BgHiMagenta`, `color.BgHiCyan`, `color.BgHiWhite`

**Special:**

- `color.Reset` - No colour (use for both Fg and Bg to disable colouring for an element)

### Global Colour Disable

The existing `Coloured: false` option continues to work and will disable ALL colours regardless of individual colour settings:

```go
opts := options.Default()
opts.Coloured = false  // Disables all colours globally
```

## Timestamp Formats

Zylog supports multiple timestamp formats:

- `formatter.RFC3339` - Full RFC3339 format (e.g., `2025-11-07T14:30:45-08:00`)
- `formatter.StandardTimestamp` - Standard format (e.g., `2006/01/02 15:04:05`)
- `formatter.SimpleTimestamp` - Compact format: YYYYMMDD.HHmmSS (e.g., `20251107.143045`) **[Default]**
- `formatter.TimeOnly` - Time only: HH:mm:SS (e.g., `14:30:45`)

```go
opts.TimestampFormat = formatter.RFC3339         // 2024-11-07T14:30:52-05:00
opts.TimestampFormat = formatter.StandardTimestamp // 2006/01/02 15:04:05
opts.TimestampFormat = formatter.SimpleTimestamp  // 20241107.143052
opts.TimestampFormat = formatter.TimeOnly         // 14:30:52
```

## Output Format

```
20241107.143052    INFO [main.main:42] ▶ Application started
20241107.143052    WARN [auth.Login:127] ▶ Failed login attempt: user={alice}, ip={192.168.1.1}
```

## Logger Comparison

### slog (Recommended for new projects)

- ✅ Standard library (no external dependencies)
- ✅ Better performance
- ✅ Native structured logging
- ✅ Logger instances for better control
- ✅ Modern Go idioms

### logrus (For existing projects)

- ✅ Global singleton pattern
- ✅ Simpler for quick setup
- ✅ Wide ecosystem support
- ✅ Familiar to many Go developers

## Background

The formatting style is inspired by the [Twig Clojure](https://github.com/clojusc/twig) and [Logjam LFE](https://github.com/lfex/logjam) libraries.

## License

© 2019,2025 ZYLISP Project
© 2019-2025, geomyidia Project

Apache License, Version 2.0

[//]: ---Named-Links---

[gh-actions-badge]: https://github.com/zylisp/zylog/actions/workflows/cicd.yml/badge.svg
[gh-actions]: https://github.com/zylisp/zylog/actions
[github-tags]: https://github.com/zylisp/zylog/tags
[github-tags-badge]: https://img.shields.io/github/tag/zylisp/zylog.svg
