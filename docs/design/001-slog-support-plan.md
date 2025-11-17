# Claude Code Implementation Prompt: Add slog Support to zylog

## Context

You are working on the `zylog` library, a Go logging wrapper that currently supports only `logrus`. The goal is to add support for Go's standard library `slog` package while maintaining backward compatibility with existing `logrus` functionality.

## Current Architecture Analysis

### Existing Structure
```
zylog/
├── cmd/zylog-demo/main.go          # Demo application
├── errors/errors.go                 # Custom error types
├── formatter/formatter.go           # logrus formatter implementation
├── level/level.go                   # Log level constants
├── logger/logger.go                 # Package documentation
├── options/options.go               # Configuration options
├── main.go                          # Package declaration
├── version.go                       # Version information
└── zylog.go                         # Main setup logic
```

### Key Components

1. **Options System** (`options/options.go`):
   - `ZyLog` struct contains: `Colored`, `Level`, `Output`, `ReportCaller`, `TimestampFormat`, `PadLevel`, `PadAmount`, `PadSide`
   - Provides `Default()`, `NoLevelPadding()`, and `NoCaller()` convenience functions

2. **Formatter** (`formatter/formatter.go`):
   - `LogLine` struct implements `logrus.Formatter` interface
   - Provides colored output with custom format: `TIMESTAMP LEVEL [CALLER] ▶ MESSAGE`
   - `ColorLevel()` function applies colors based on log level
   - Supports three timestamp formats: RFC3339, SimpleTimestamp (YYYYMMDD.HHmmSS), TimeOnly (HH:mm:SS)
   - Supports level padding (left or right alignment)

3. **Main Setup** (`zylog.go`):
   - `SetupLogging()` configures logrus with provided options
   - Handles output destination (stdout, stderr, filesystem-not-implemented)
   - Sets color mode, timestamp format, and caller reporting

## Implementation Requirements

### Phase 1: Create slog Handler

**File**: `handler/slog_handler.go`

Create a new `slog.Handler` implementation that mimics the exact formatting style of the existing `formatter.LogLine`:

```go
package handler

import (
    "context"
    "fmt"
    "io"
    "log/slog"
    "runtime"
    "strconv"
    "strings"
    "time"

    "github.com/fatih/color"
    "github.com/zylisp/zylog/formatter"
    "github.com/zylisp/zylog/level"
    "github.com/zylisp/zylog/options"
)

type ZylogHandler struct {
    opts   *options.ZyLog
    writer io.Writer
    attrs  []slog.Attr
    groups []string
}

func NewZylogHandler(writer io.Writer, opts *options.ZyLog) *ZylogHandler {
    return &ZylogHandler{
        opts:   opts,
        writer: writer,
        attrs:  make([]slog.Attr, 0),
        groups: make([]string, 0),
    }
}

func (h *ZylogHandler) Enabled(ctx context.Context, level slog.Level) bool {
    // Map slog levels to zylog levels and check against configured level
    // Implement level comparison logic
}

func (h *ZylogHandler) Handle(ctx context.Context, r slog.Record) error {
    // Format the log record using the same style as formatter.LogLine
    // 1. Format timestamp using h.opts.TimestampFormat
    // 2. Format level using ColorLevel() equivalent for slog levels
    // 3. If ReportCaller is true, extract and format caller info
    // 4. Format message with cyan arrow
    // 5. Append any attributes from r.Attrs() and h.attrs
    // 6. Write to h.writer
}

func (h *ZylogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    // Return a new handler with additional attributes
    // Clone the handler and append attrs to h.attrs
}

func (h *ZylogHandler) WithGroup(name string) slog.Handler {
    // Return a new handler with a group name
    // Clone the handler and append name to h.groups
}
```

**Critical Implementation Details**:

1. **Level Mapping**: Create a function to map `slog.Level` to zylog level strings:
   ```go
   func slogLevelToString(level slog.Level) string {
       switch level {
       case slog.LevelDebug:
           return level.Debug
       case slog.LevelInfo:
           return level.Info
       case slog.LevelWarn:
           return level.Warn
       case slog.LevelError:
           return level.Error
       default:
           if level < slog.LevelDebug {
               return level.Trace
           }
           return level.Fatal
       }
   }
   ```

2. **Caller Information**: Extract caller from `slog.Record`:
   ```go
   if h.opts.ReportCaller && r.PC != 0 {
       fs := runtime.CallersFrames([]uintptr{r.PC})
       f, _ := fs.Next()
       // Format: [function:line]
   }
   ```

3. **Timestamp Formatting**: Reuse the existing `TSFormat.ToTimeFormat()` method:
   ```go
   timestamp := r.Time.Format(h.opts.TimestampFormat.ToTimeFormat())
   ```

4. **Color Handling**: Respect the `h.opts.Colored` flag and use `color.NoColor`

5. **Attribute Formatting**: Format attributes similarly to logrus fields:
   ```go
   // Append " || " before attributes if present
   // Format each attribute as "key={value}, "
   ```

### Phase 2: Refactor Shared Formatting Logic

**File**: `formatter/common.go`

Extract common formatting functions that can be used by both logrus and slog handlers:

```go
package formatter

import (
    "fmt"
    "github.com/fatih/color"
    "github.com/zylisp/zylog/level"
    "strings"
)

// FormatTimestamp formats a time string with color
func FormatTimestamp(timestamp string) string {
    return color.HiBlackString(timestamp)
}

// FormatMessage formats a log message with color
func FormatMessage(message string) string {
    return color.GreenString(message)
}

// FormatArrow returns the colored arrow separator
func FormatArrow() string {
    return color.CyanString(" ▶ ")
}

// FormatCaller formats caller information
func FormatCaller(function string, line int) string {
    return fmt.Sprintf(" [%s:%s]",
        color.HiYellowString(function),
        color.YellowString(fmt.Sprintf("%d", line)))
}

// FormatAttributeSeparator returns the attribute separator
func FormatAttributeSeparator() string {
    return " || "
}

// Note: ColorLevel already exists in formatter.go and can be reused
```

**File**: Update `formatter/formatter.go` to use these common functions in `LogLine.Format()`

### Phase 3: Add slog Setup Function

**File**: `zylog.go`

Add a new function to setup slog with zylog formatting:

```go
import (
    "log/slog"
    "os"
    
    "github.com/fatih/color"
    log "github.com/sirupsen/logrus"
    
    "github.com/zylisp/zylog/errors"
    "github.com/zylisp/zylog/formatter"
    "github.com/zylisp/zylog/handler"
    "github.com/zylisp/zylog/options"
)

// SetupSlog configures and returns a new slog.Logger with zylog formatting.
// Returns a configured *slog.Logger instance that can be used directly or set as the default logger.
func SetupSlog(opts *options.ZyLog) *slog.Logger {
    // 1. Determine output writer based on opts.Output
    var writer io.Writer
    switch opts.Output {
    case StdOut:
        writer = os.Stdout
    case StdErr:
        writer = os.Stderr
    case FileSystem:
        panic(errors.ErrNotImplemented("filesystem log output"))
    default:
        panic(errors.ErrUnsupLogOutput(opts.Output))
    }
    
    // 2. Configure color mode
    disableColors := !opts.Colored
    color.NoColor = disableColors
    
    // 3. Set default timestamp format if unset
    timestampFormat := opts.TimestampFormat
    if timestampFormat == formatter.TSUnset {
        timestampFormat = formatter.SimpleTimestamp
    }
    
    // 4. Update opts with resolved timestamp format
    opts.TimestampFormat = timestampFormat
    
    // 5. Create handler
    h := handler.NewZylogHandler(writer, opts)
    
    // 6. Parse and map log level
    level := parseLogLevel(opts.Level)
    handlerOpts := &slog.HandlerOptions{
        Level:     level,
        AddSource: opts.ReportCaller,
    }
    
    // 7. Create and return logger
    logger := slog.New(h)
    logger.Info("Slog logging initialized.")
    return logger
}

// parseLogLevel converts string level to slog.Level
func parseLogLevel(levelStr string) slog.Level {
    switch strings.ToLower(levelStr) {
    case "trace":
        return slog.LevelDebug - 1 // Use lower than debug for trace
    case "debug":
        return slog.LevelDebug
    case "info":
        return slog.LevelInfo
    case "warn", "warning":
        return slog.LevelWarn
    case "error":
        return slog.LevelError
    case "fatal", "panic":
        return slog.LevelError + 1 // Use higher than error
    default:
        return slog.LevelInfo
    }
}

// Keep existing SetupLogging function for backward compatibility
```

### Phase 4: Update Demo Application

**File**: `cmd/zylog-demo/main.go`

Add slog demonstrations alongside existing logrus examples:

```go
package main

import (
    "log/slog"
    
    log "github.com/sirupsen/logrus"
    
    "github.com/zylisp/zylog"
    "github.com/zylisp/zylog/formatter"
    "github.com/zylisp/zylog/options"
)

// ... keep existing SetupLogger functions ...

func main() {
    zylog.PrintVersions()
    
    // === LOGRUS DEMOS ===
    log.Info("========================================")
    log.Info("=== LOGRUS DEMONSTRATIONS ===")
    log.Info("========================================")
    
    SetupLogger()
    log.Info(" *** Logging with default options ***")
    log.Trace("This is trace")
    log.Debug("This is debug")
    log.Info("This is info")
    log.Warn("This is warn")
    log.Error("This is error")
    log.Info("Fatal and Panic are also supported; " +
        "Fatal will os.Exit, and Panic will log, then panic().")
    
    SetupLoggerNoCaller()
    log.Info(" *** Logging with calling function disabled ***")
    log.Trace("This is trace")
    log.Debug("This is debug")
    log.Info("This is info")
    log.Warn("This is warn")
    log.Error("This is error")
    
    SetupLoggerNoPad()
    log.Info(" *** Logging with level padding disabled and timestamps in RFC3339 format ***")
    log.Trace("This is trace")
    log.Debug("This is debug")
    log.Info("This is info")
    log.Warn("This is warn")
    log.Error("This is error")
    
    SetupLoggerNoColour()
    log.Info(" *** Logging with level padding disabled, timestamps in time-only format, and no colour ***")
    log.Trace("This is trace")
    log.Debug("This is debug")
    log.Info("This is info")
    log.Warn("This is warn")
    log.Error("This is error")
    
    // === SLOG DEMOS ===
    log.Info("========================================")
    log.Info("=== SLOG DEMONSTRATIONS ===")
    log.Info("========================================")
    
    // Demo 1: Default options
    logger := zylog.SetupSlog(options.Default())
    logger.Info(" *** Slog with default options ***")
    logger.Debug("This is debug")
    logger.Info("This is info")
    logger.Warn("This is warn")
    logger.Error("This is error")
    
    // Demo 2: No caller
    logger2 := zylog.SetupSlog(options.NoCaller())
    logger2.Info(" *** Slog with calling function disabled ***")
    logger2.Debug("This is debug")
    logger2.Info("This is info")
    logger2.Warn("This is warn")
    logger2.Error("This is error")
    
    // Demo 3: No padding, RFC3339
    opts := options.NoLevelPadding()
    opts.TimestampFormat = formatter.RFC3339
    logger3 := zylog.SetupSlog(opts)
    logger3.Info(" *** Slog with level padding disabled and timestamps in RFC3339 format ***")
    logger3.Debug("This is debug")
    logger3.Info("This is info")
    logger3.Warn("This is warn")
    logger3.Error("This is error")
    
    // Demo 4: No color, time-only
    opts4 := options.NoLevelPadding()
    opts4.TimestampFormat = formatter.TimeOnly
    opts4.Colored = false
    logger4 := zylog.SetupSlog(opts4)
    logger4.Info(" *** Slog with level padding disabled, timestamps in time-only format, and no colour ***")
    logger4.Debug("This is debug")
    logger4.Info("This is info")
    logger4.Warn("This is warn")
    logger4.Error("This is error")
    
    // Demo 5: Structured logging with attributes
    logger.Info("Structured logging example",
        "user", "alice",
        "request_id", "12345",
        "duration_ms", 42)
    
    // Demo 6: With context attributes
    userLogger := logger.With("user", "bob", "role", "admin")
    userLogger.Info("User performed action", "action", "delete")
    userLogger.Warn("User attempted restricted operation", "action", "access_logs")
}
```

### Phase 5: Update Documentation

**File**: `README.md` (create if doesn't exist)

Add comprehensive documentation:

```markdown
# zylog - Styled Logging for Go

A Go logging library that provides beautiful, consistent formatting for both `logrus` and Go's standard library `slog`.

## Features

- 🎨 Colored output with customizable colors
- 📍 Optional caller information (package, function, line number)
- ⏰ Multiple timestamp formats (RFC3339, Simple, Time-only)
- 📏 Configurable level padding for alignment
- 🔄 Support for both logrus and slog
- 🚀 Simple, minimal configuration

## Installation

```bash
go get github.com/zylisp/zylog
```

## Quick Start

### Using with logrus

```go
package main

import (
    log "github.com/sirupsen/logrus"
    "github.com/zylisp/zylog"
    "github.com/zylisp/zylog/options"
)

func main() {
    zylog.SetupLogging(options.Default())
    log.Info("Application started")
}
```

### Using with slog

```go
package main

import (
    "github.com/zylisp/zylog"
    "github.com/zylisp/zylog/options"
)

func main() {
    logger := zylog.SetupSlog(options.Default())
    logger.Info("Application started")
    
    // Structured logging
    logger.Info("User logged in",
        "user_id", 12345,
        "ip", "192.168.1.1")
}
```

## Configuration Options

```go
type ZyLog struct {
    Colored         bool              // Enable colored output
    Level           string            // Log level: "trace", "debug", "info", "warn", "error"
    Output          string            // Output destination: "stdout", "stderr"
    ReportCaller    bool              // Include caller information
    TimestampFormat formatter.TSFormat // Timestamp format
    PadLevel        bool              // Pad level strings for alignment
    PadAmount       int               // Number of characters to pad to
    PadSide         string            // "left" or "right"
}
```

### Preset Configurations

```go
// Default configuration with all features enabled
opts := options.Default()

// Disable caller reporting
opts := options.NoCaller()

// Disable level padding
opts := options.NoLevelPadding()
```

## Output Format

```
20241107.143052 INFO    [main.main:42] ▶ Application started
20241107.143052 WARN    [auth.Login:127] ▶ Failed login attempt || user={alice}, ip={192.168.1.1}
```

## Timestamp Formats

```go
opts.TimestampFormat = formatter.RFC3339         // 2024-11-07T14:30:52-05:00
opts.TimestampFormat = formatter.SimpleTimestamp // 20241107.143052
opts.TimestampFormat = formatter.TimeOnly        // 14:30:52
```

## Comparison: logrus vs slog

### logrus
- Global logger with package-level functions
- Simpler for quick setup
- Fields added via `WithFields()`
- Best for applications already using logrus

### slog
- Logger instances
- Better performance
- Native structured logging
- Recommended for new projects
- Standard library (no external dependencies)

## Examples

See `cmd/zylog-demo` for comprehensive examples.

## License

[Your License]
```

**File**: Update `logger/logger.go` package documentation to mention slog support

### Phase 6: Add Tests

**File**: `handler/slog_handler_test.go`

```go
package handler

import (
    "bytes"
    "context"
    "log/slog"
    "strings"
    "testing"
    
    "github.com/zylisp/zylog/formatter"
    "github.com/zylisp/zylog/options"
)

func TestZylogHandler_Basic(t *testing.T) {
    buf := &bytes.Buffer{}
    opts := &options.ZyLog{
        Colored:         false,
        Level:           "debug",
        Output:          "stdout",
        ReportCaller:    false,
        TimestampFormat: formatter.SimpleTimestamp,
        PadLevel:        true,
        PadAmount:       7,
        PadSide:         "left",
    }
    
    h := NewZylogHandler(buf, opts)
    logger := slog.New(h)
    
    logger.Info("test message")
    
    output := buf.String()
    if !strings.Contains(output, "INFO") {
        t.Errorf("Expected INFO in output, got: %s", output)
    }
    if !strings.Contains(output, "test message") {
        t.Errorf("Expected message in output, got: %s", output)
    }
}

func TestZylogHandler_WithAttrs(t *testing.T) {
    buf := &bytes.Buffer{}
    opts := &options.ZyLog{
        Colored:         false,
        Level:           "debug",
        Output:          "stdout",
        ReportCaller:    false,
        TimestampFormat: formatter.SimpleTimestamp,
        PadLevel:        false,
        PadAmount:       0,
        PadSide:         "left",
    }
    
    h := NewZylogHandler(buf, opts)
    logger := slog.New(h)
    
    logger.Info("test message", "key", "value", "number", 42)
    
    output := buf.String()
    if !strings.Contains(output, "key={value}") {
        t.Errorf("Expected key=value in output, got: %s", output)
    }
    if !strings.Contains(output, "number={42}") {
        t.Errorf("Expected number=42 in output, got: %s", output)
    }
}

func TestZylogHandler_Caller(t *testing.T) {
    buf := &bytes.Buffer{}
    opts := &options.ZyLog{
        Colored:         false,
        Level:           "debug",
        Output:          "stdout",
        ReportCaller:    true,
        TimestampFormat: formatter.SimpleTimestamp,
        PadLevel:        false,
        PadAmount:       0,
        PadSide:         "left",
    }
    
    h := NewZylogHandler(buf, opts)
    logger := slog.New(h)
    
    logger.Info("test message")
    
    output := buf.String()
    if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
        t.Errorf("Expected caller info in brackets, got: %s", output)
    }
}

func TestZylogHandler_Levels(t *testing.T) {
    tests := []struct {
        name     string
        logLevel string
        logFunc  func(*slog.Logger, string)
        want     string
    }{
        {"debug", "debug", (*slog.Logger).Debug, "DEBUG"},
        {"info", "info", (*slog.Logger).Info, "INFO"},
        {"warn", "warn", (*slog.Logger).Warn, "WARN"},
        {"error", "error", (*slog.Logger).Error, "ERROR"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            buf := &bytes.Buffer{}
            opts := &options.ZyLog{
                Colored:         false,
                Level:           tt.logLevel,
                Output:          "stdout",
                ReportCaller:    false,
                TimestampFormat: formatter.SimpleTimestamp,
                PadLevel:        false,
                PadAmount:       0,
                PadSide:         "left",
            }
            
            h := NewZylogHandler(buf, opts)
            logger := slog.New(h)
            
            tt.logFunc(logger, "test")
            
            output := buf.String()
            if !strings.Contains(output, tt.want) {
                t.Errorf("Expected %s in output, got: %s", tt.want, output)
            }
        })
    }
}
```

**File**: `zylog_test.go`

```go
package zylog

import (
    "bytes"
    "log/slog"
    "os"
    "strings"
    "testing"
    
    "github.com/zylisp/zylog/formatter"
    "github.com/zylisp/zylog/options"
)

func TestSetupSlog_Basic(t *testing.T) {
    opts := &options.ZyLog{
        Colored:         false,
        Level:           "info",
        Output:          "stdout",
        ReportCaller:    false,
        TimestampFormat: formatter.SimpleTimestamp,
        PadLevel:        false,
        PadAmount:       0,
        PadSide:         "left",
    }
    
    logger := SetupSlog(opts)
    if logger == nil {
        t.Fatal("Expected logger, got nil")
    }
}

func TestSetupSlog_InvalidOutput(t *testing.T) {
    opts := &options.ZyLog{
        Colored:         false,
        Level:           "info",
        Output:          "invalid",
        ReportCaller:    false,
        TimestampFormat: formatter.SimpleTimestamp,
        PadLevel:        false,
        PadAmount:       0,
        PadSide:         "left",
    }
    
    defer func() {
        if r := recover(); r == nil {
            t.Error("Expected panic for invalid output")
        }
    }()
    
    SetupSlog(opts)
}

func TestParseLogLevel(t *testing.T) {
    tests := []struct {
        input string
        want  slog.Level
    }{
        {"trace", slog.LevelDebug - 1},
        {"debug", slog.LevelDebug},
        {"info", slog.LevelInfo},
        {"warn", slog.LevelWarn},
        {"error", slog.LevelError},
        {"invalid", slog.LevelInfo}, // defaults to info
    }
    
    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            got := parseLogLevel(tt.input)
            if got != tt.want {
                t.Errorf("parseLogLevel(%q) = %v, want %v", tt.input, got, tt.want)
            }
        })
    }
}
```

## Implementation Checklist

- [ ] Create `handler/slog_handler.go` with complete `ZylogHandler` implementation
- [ ] Implement `Enabled()`, `Handle()`, `WithAttrs()`, and `WithGroup()` methods
- [ ] Add level mapping function `slogLevelToString()`
- [ ] Implement caller extraction from `slog.Record`
- [ ] Create `formatter/common.go` with shared formatting functions
- [ ] Refactor `formatter/formatter.go` to use common functions
- [ ] Add `SetupSlog()` function to `zylog.go`
- [ ] Add `parseLogLevel()` helper function
- [ ] Update `cmd/zylog-demo/main.go` with slog demonstrations
- [ ] Create comprehensive `README.md`
- [ ] Update `logger/logger.go` documentation
- [ ] Create `handler/slog_handler_test.go` with unit tests
- [ ] Create `zylog_test.go` with integration tests
- [ ] Ensure all tests pass: `go test ./...`
- [ ] Run demo to verify output: `go run cmd/zylog-demo/main.go`
- [ ] Verify output format matches between logrus and slog versions
- [ ] Check that colors work correctly (run demo in terminal)
- [ ] Verify caller information displays correctly
- [ ] Test all timestamp formats
- [ ] Test level padding options
- [ ] Verify structured logging (attributes) works correctly

## Success Criteria

1. **Backward Compatibility**: All existing logrus functionality continues to work exactly as before
2. **Format Consistency**: slog output matches logrus output format exactly (colors, spacing, arrows, etc.)
3. **Feature Parity**: All zylog options work identically for both logrus and slog
4. **Tests Pass**: All unit and integration tests pass
5. **Demo Works**: The demo application successfully demonstrates both logrus and slog with identical output styles
6. **Documentation**: Clear documentation explaining how to use both backends

## Technical Notes

### Key Differences to Handle

1. **Caller Reporting**:
   - logrus: Uses `entry.Caller.Function` and `entry.Caller.Line`
   - slog: Uses `runtime.CallersFrames([]uintptr{r.PC})`

2. **Attributes vs Fields**:
   - logrus: `entry.Data` is a `map[string]interface{}`
   - slog: Attributes accessed via `r.Attrs(func(a slog.Attr) bool)`

3. **Level Filtering**:
   - logrus: Built into logger with `log.SetLevel()`
   - slog: Handler's `Enabled()` method must implement filtering

4. **Handler Immutability**:
   - slog handlers should be immutable
   - `WithAttrs()` and `WithGroup()` must return new handlers, not modify existing ones

### Performance Considerations

- Avoid unnecessary allocations in `Handle()` method
- Reuse buffers where possible
- Consider using sync.Pool for buffer allocation if performance becomes critical

### Color Handling

- Respect `color.NoColor` global setting
- Ensure colors are disabled when `opts.Colored == false`
- Colors should work identically between logrus and slog implementations

## Questions to Resolve

1. Should `SetupSlog()` return a logger or also set a default logger globally?
   - Recommendation: Return a logger (current approach) since slog encourages explicit logger passing
   
2. How to handle Fatal and Panic levels in slog?
   - Recommendation: Map to `slog.LevelError + 1` but don't call `os.Exit()` or `panic()` - let applications handle this

3. Should there be a `SetupSlogAsDefault()` function?
   - Recommendation: Yes, add this for convenience:
     ```go
     func SetupSlogAsDefault(opts *options.ZyLog) {
         logger := SetupSlog(opts)
         slog.SetDefault(logger)
     }
     ```

## Additional Enhancements (Optional)

1. **Benchmark Tests**: Add benchmarks comparing logrus vs slog performance
2. **Context Support**: Demonstrate extracting values from `context.Context` in slog
3. **Custom Levels**: Show how to use custom slog levels beyond the standard ones
4. **Handler Wrapping**: Provide example of wrapping zylog handler with other handlers
5. **Migration Guide**: Document how to migrate from logrus to slog within the same codebase

## Validation Steps

1. Build the project: `go build ./...`
2. Run all tests: `go test -v ./...`
3. Run the demo: `go run cmd/zylog-demo/main.go`
4. Visually compare logrus and slog output sections
5. Test with colors: `TERM=xterm-256color go run cmd/zylog-demo/main.go`
6. Test without colors: `NO_COLOR=1 go run cmd/zylog-demo/main.go`
7. Verify structured logging output includes attributes
8. Check that caller info matches between implementations
9. Ensure timestamp formats are identical
10. Verify level padding works consistently

## End Result

After implementation, users should be able to:

```go
// Use logrus (existing)
zylog.SetupLogging(options.Default())
log.Info("Using logrus")

// Use slog (new)
logger := zylog.SetupSlog(options.Default())
logger.Info("Using slog")

// Both produce identical visual output format
```

The library should feel like a unified logging solution that happens to support two backends, not two separate implementations awkwardly joined together.
