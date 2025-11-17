# Implementation Prompt: Configurable Colors for zylog

## Overview

Add user-configurable foreground and background colors for all formatted output elements in zylog. Currently, all colors are hardcoded using `github.com/fatih/color` definitions. This implementation will introduce a new color configuration system that allows users to customize every colored element while maintaining backward compatibility with existing behavior.

## Current State

- Version: `0.2.0` (will be incremented to `0.2.1`)
- Colors are hardcoded in:
  - `formatter/common.go`: timestamp, message, arrow, caller, attribute keys/values
  - `formatter/formatter.go`: log levels (ColorLevel function)
- Uses `github.com/fatih/color` package for all color output
- Global color disable via `opts.Colored bool` field

## Goals

1. Make ALL colored output elements user-configurable
2. Support both foreground and background colors for each element
3. Maintain backward compatibility - nil Colors pointer uses current defaults
4. Allow per-element color disable (using fatih/color's no-color capability)
5. Respect existing global `Colored bool` disable flag
6. Provide clear documentation and examples

## Implementation Steps

### Step 1: Create New Color Types (`colors/colors.go`)

Create a new package `colors` with the following structure:

```go
// Package colors provides color configuration types for zylog.
package colors

import "github.com/fatih/color"

// Color represents foreground and background color attributes for a single element.
// Use color.Attribute constants from github.com/fatih/color (e.g., color.FgRed, color.BgBlue).
// Set to color.Reset (or 0) to disable color for that component.
type Color struct {
	Fg color.Attribute // Foreground color
	Bg color.Attribute // Background color
}

// Colours holds color configuration for all formatted output elements in zylog.
// A nil Colours pointer will use default colors (current hardcoded behavior).
type Colours struct {
	// Timestamp colors (currently light grey/HiBlack)
	Timestamp *Color
	
	// Log level colors
	LevelTrace   *Color // Currently HiMagenta
	LevelDebug   *Color // Currently HiCyan
	LevelInfo    *Color // Currently HiGreen
	LevelWarn    *Color // Currently HiYellow
	LevelWarning *Color // Currently HiYellow (alias for Warn)
	LevelError   *Color // Currently Red
	LevelFatal   *Color // Currently HiRed
	LevelPanic   *Color // Currently HiWhite
	
	// Message text color (currently green)
	Message *Color
	
	// Arrow separator (currently cyan " ▶ ")
	Arrow *Color
	
	// Caller information colors
	CallerFunction *Color // Currently HiYellow
	CallerLine     *Color // Currently Yellow
	
	// Structured logging attribute colors
	AttrKey   *Color // Currently Yellow
	AttrValue *Color // Currently HiYellow
}

// Default returns the default color configuration matching current hardcoded behavior.
func Default() *Colours {
	return &Colours{
		Timestamp:      &Color{Fg: color.FgHiBlack, Bg: color.Reset},
		LevelTrace:     &Color{Fg: color.FgHiMagenta, Bg: color.Reset},
		LevelDebug:     &Color{Fg: color.FgHiCyan, Bg: color.Reset},
		LevelInfo:      &Color{Fg: color.FgHiGreen, Bg: color.Reset},
		LevelWarn:      &Color{Fg: color.FgHiYellow, Bg: color.Reset},
		LevelWarning:   &Color{Fg: color.FgHiYellow, Bg: color.Reset},
		LevelError:     &Color{Fg: color.FgRed, Bg: color.Reset},
		LevelFatal:     &Color{Fg: color.FgHiRed, Bg: color.Reset},
		LevelPanic:     &Color{Fg: color.FgHiWhite, Bg: color.Reset},
		Message:        &Color{Fg: color.FgGreen, Bg: color.Reset},
		Arrow:          &Color{Fg: color.FgCyan, Bg: color.Reset},
		CallerFunction: &Color{Fg: color.FgHiYellow, Bg: color.Reset},
		CallerLine:     &Color{Fg: color.FgYellow, Bg: color.Reset},
		AttrKey:        &Color{Fg: color.FgYellow, Bg: color.Reset},
		AttrValue:      &Color{Fg: color.FgHiYellow, Bg: color.Reset},
	}
}

// ApplyColor applies the color to a string. If color is nil, returns string unchanged.
// If both Fg and Bg are color.Reset (0), returns string unchanged (no color).
func (c *Color) ApplyColor(s string) string {
	if c == nil {
		return s
	}
	if c.Fg == color.Reset && c.Bg == color.Reset {
		return s
	}
	
	// Build attribute list
	attrs := []color.Attribute{}
	if c.Fg != color.Reset {
		attrs = append(attrs, c.Fg)
	}
	if c.Bg != color.Reset {
		attrs = append(attrs, c.Bg)
	}
	
	if len(attrs) == 0 {
		return s
	}
	
	return color.New(attrs...).Sprint(s)
}
```

### Step 2: Update `options/options.go`

Add the Colours field to the ZyLog struct:

```go
import (
	"fmt"

	"github.com/zylisp/zylog/colors"
	"github.com/zylisp/zylog/formatter"
)

// ZyLog are used by the zylog logger to set up logrus.
type ZyLog struct {
	Colored         bool
	Level           string
	Output          string // stdout, stderr, or filesystem
	ReportCaller    bool
	TimestampFormat formatter.TSFormat // RFC3339, Simple (YYYYMMDD.HHmmSS), or Time (HH:mm:SS)
	PadLevel        bool               // Whether to pad level strings for alignment
	PadAmount       int                // Number of characters to pad level strings to
	PadSide         string             // "left" or "right"; which side to pad level strings on
	MsgSeparator    string             // Separator between message and attributes
	Logger          Logger             // Logger type: Logrus or Slog
	Colours         *colors.Colours    // Color configuration (nil uses defaults)
}
```

Update the `defaultOpts` variable:

```go
var (
	defaultOpts = &ZyLog{
		Colored:         true,
		Level:           "trace",
		Output:          "stdout",
		ReportCaller:    true,
		TimestampFormat: formatter.SimpleTimestamp,
		PadLevel:        false,
		PadAmount:       5,
		PadSide:         "left",
		MsgSeparator:    ": ",
		Logger:          Slog,
		Colours:         nil, // nil means use defaults
	}
)
```

### Step 3: Update `formatter/common.go`

Replace all hardcoded color functions with configurable versions:

```go
package formatter

import (
	"fmt"

	"github.com/zylisp/zylog/colors"
)

// FormatTimestamp formats a time string with the configured color.
func FormatTimestamp(timestamp string, colours *colors.Colours) string {
	if colours == nil {
		colours = colors.Default()
	}
	if colours.Timestamp == nil {
		return timestamp
	}
	return colours.Timestamp.ApplyColor(timestamp)
}

// FormatMessage formats a log message with the configured color.
func FormatMessage(message string, colours *colors.Colours) string {
	if colours == nil {
		colours = colors.Default()
	}
	if colours.Message == nil {
		return message
	}
	return colours.Message.ApplyColor(message)
}

// FormatArrow returns the colored arrow separator.
func FormatArrow(colours *colors.Colours) string {
	if colours == nil {
		colours = colors.Default()
	}
	arrow := " ▶ "
	if colours.Arrow == nil {
		return arrow
	}
	return colours.Arrow.ApplyColor(arrow)
}

// FormatCaller formats caller information with the configured colors.
func FormatCaller(function string, line int, colours *colors.Colours) string {
	if colours == nil {
		colours = colors.Default()
	}
	
	functionStr := function
	if colours.CallerFunction != nil {
		functionStr = colours.CallerFunction.ApplyColor(function)
	}
	
	lineStr := fmt.Sprintf("%d", line)
	if colours.CallerLine != nil {
		lineStr = colours.CallerLine.ApplyColor(lineStr)
	}
	
	return fmt.Sprintf(" [%s:%s]", functionStr, lineStr)
}

// FormatAttrKey formats an attribute key with the configured color.
func FormatAttrKey(key string, colours *colors.Colours) string {
	if colours == nil {
		colours = colors.Default()
	}
	if colours.AttrKey == nil {
		return key
	}
	return colours.AttrKey.ApplyColor(key)
}

// FormatAttrValue formats an attribute value with the configured color.
func FormatAttrValue(value string, colours *colors.Colours) string {
	if colours == nil {
		colours = colors.Default()
	}
	if colours.AttrValue == nil {
		return value
	}
	return colours.AttrValue.ApplyColor(value)
}
```

### Step 4: Update `formatter/formatter.go`

Update the LogLine struct to include colors:

```go
// LogLine formats logs into a complete line.
type LogLine struct {
	// Force disabling colors.
	DisableColors bool
	// TimestampFormat specifies the format for timestamps.
	TimestampFormat TSFormat
	// PadLevel specifies whether to pad level strings for alignment.
	PadLevel bool
	// PadAmount specifies the total width to pad level strings to.
	PadAmount int
	// PadSide specifies which side to add padding on ("left" or "right").
	PadSide string
	// AttrSeparator specifies the separator between message and attributes.
	MsgSeparator string
	// Colours specifies the color configuration (nil uses defaults).
	Colours *colors.Colours
}
```

Update the Format method to pass colours through:

```go
func (f *LogLine) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := FormatTimestamp(entry.Time.Format(f.TimestampFormat.ToTimeFormat()), f.Colours)
	level := ColorLevel(strings.ToUpper(entry.Level.String()), f.PadLevel, f.PadAmount, f.PadSide, f.Colours)

	fmt.Fprintf(b, "%s %s", timestamp, level)
	if entry.Logger.ReportCaller {
		b.WriteString(FormatCaller(entry.Caller.Function, entry.Caller.Line, f.Colours))
	}
	if entry.Message != "" {
		b.WriteString(FormatArrow(f.Colours))
		b.WriteString(FormatMessage(entry.Message, f.Colours))
	}

	if len(entry.Data) > 0 {
		b.WriteString(f.MsgSeparator)
		first := true
		for key, value := range entry.Data {
			if !first {
				b.WriteString(", ")
			}
			fmt.Fprintf(b, "%s={%s}", 
				FormatAttrKey(key, f.Colours), 
				FormatAttrValue(fmt.Sprintf("%v", value), f.Colours))
			first = false
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}
```

Update the ColorLevel function signature and implementation:

```go
// ColorLevel determines the color of the log level based upon the string
// value of the log level. If padLevel is true, the level string will be
// padded to padAmount characters, aligned according to padSide.
func ColorLevel(lvl string, padLevel bool, padAmount int, padSide string, colours *colors.Colours) string {
	if colours == nil {
		colours = colors.Default()
	}
	
	// Apply padding before colorizing
	if padLevel && padAmount > 0 {
		if padSide == "left" {
			// Right-align: add spaces on the left
			lvl = fmt.Sprintf("%*s", padAmount, lvl)
		} else {
			// Left-align (default): add spaces on the right
			lvl = fmt.Sprintf("%-*s", padAmount, lvl)
		}
	}

	// Now colorize the padded string based on level
	trimmedLvl := strings.TrimSpace(lvl)
	var colorConfig *colors.Color
	
	switch trimmedLvl {
	case level.Trace:
		colorConfig = colours.LevelTrace
	case level.Debug:
		colorConfig = colours.LevelDebug
	case level.Info:
		colorConfig = colours.LevelInfo
	case level.Warn, level.Warning:
		colorConfig = colours.LevelWarn
	case level.Error:
		colorConfig = colours.LevelError
	case level.Fatal:
		colorConfig = colours.LevelFatal
	case level.Panic:
		colorConfig = colours.LevelPanic
	default:
		return lvl
	}
	
	if colorConfig == nil {
		return lvl
	}
	
	return colorConfig.ApplyColor(lvl)
}
```

### Step 5: Update `logger/logrus.go`

Pass the Colours configuration to the formatter:

```go
func SetupLogRUs(opts *options.ZyLog) {
	level, err := log.ParseLevel(opts.Level)
	if err != nil {
		panic(errors.ErrLogLevel)
	}
	log.SetLevel(level)
	switch opts.Output {
	case StdOut:
		log.SetOutput(os.Stdout)
	case StdErr:
		log.SetOutput(os.Stderr)
	case FileSystem:
		panic(errors.ErrNotImplemented("filesystem log output"))
	default:
		panic(errors.ErrUnsupLogOutput(opts.Output))
	}
	disableColors := !opts.Colored
	color.NoColor = disableColors
	timestampFormat := opts.TimestampFormat
	if timestampFormat == formatter.TSUnset {
		// Default to Simple if not set
		timestampFormat = formatter.SimpleTimestamp
	}
	log.SetFormatter(&formatter.LogLine{
		DisableColors:   disableColors,
		TimestampFormat: timestampFormat,
		PadLevel:        opts.PadLevel,
		PadAmount:       opts.PadAmount,
		PadSide:         opts.PadSide,
		MsgSeparator:    opts.MsgSeparator,
		Colours:         opts.Colours, // Add this line
	})
	log.SetReportCaller(opts.ReportCaller)
	log.Info("Logging initialized.")
}
```

### Step 6: Update `logger/slog.go`

Update the SLogHandler struct:

```go
// SLogHandler implements slog.Handler with zylog formatting.
type SLogHandler struct {
	opts   *options.ZyLog
	writer io.Writer
	attrs  []slog.Attr
	groups []string
}
```

Update the Handle method to pass colours through to all format functions:

```go
func (h *SLogHandler) Handle(_ context.Context, r slog.Record) error {
	// Build the log line using the same format as formatter.LogLine
	var buf strings.Builder

	// Get colours (nil means use defaults)
	colours := h.opts.Colours

	// 1. Format timestamp
	timestampStr := r.Time.Format(h.opts.TimestampFormat.ToTimeFormat())
	buf.WriteString(formatter.FormatTimestamp(timestampStr, colours))
	buf.WriteString(" ")

	// 2. Format level
	levelStr := slogLevelToString(r.Level)
	levelFormatted := formatter.ColorLevel(levelStr, h.opts.PadLevel, h.opts.PadAmount, h.opts.PadSide, colours)
	buf.WriteString(levelFormatted)

	// 3. Format caller if enabled
	if h.opts.ReportCaller && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		buf.WriteString(formatter.FormatCaller(f.Function, f.Line, colours))
	}

	// 4. Format message
	if r.Message != "" {
		buf.WriteString(formatter.FormatArrow(colours))
		buf.WriteString(formatter.FormatMessage(r.Message, colours))
	}

	// 5. Format attributes
	hasAttrs := len(h.attrs) > 0 || r.NumAttrs() > 0
	if hasAttrs {
		buf.WriteString(h.opts.MsgSeparator)
		first := true

		// Add handler-level attributes first
		for _, attr := range h.attrs {
			if !first {
				buf.WriteString(", ")
			}
			h.appendAttr(&buf, attr, colours)
			first = false
		}

		// Add record-level attributes
		r.Attrs(func(a slog.Attr) bool {
			if !first {
				buf.WriteString(", ")
			}
			h.appendAttr(&buf, a, colours)
			first = false
			return true
		})
	}

	// 6. Add newline
	buf.WriteString("\n")

	// Write to output
	_, err := h.writer.Write([]byte(buf.String()))
	return err
}
```

Update the appendAttr method signature:

```go
// appendAttr appends a single attribute to the buffer in zylog format.
func (h *SLogHandler) appendAttr(buf *strings.Builder, attr slog.Attr, colours *colors.Colours) {
	// Handle groups
	prefix := ""
	if len(h.groups) > 0 {
		prefix = strings.Join(h.groups, ".") + "."
	}

	key := prefix + attr.Key
	value := attr.Value.String()

	fmt.Fprintf(buf, "%s={%s}", 
		formatter.FormatAttrKey(key, colours), 
		formatter.FormatAttrValue(value, colours))
}
```

### Step 7: Update `README.md`

Add a new section after "Configuration Options" showing color customization:

```markdown
## Color Customization

Zylog allows you to customize the foreground and background colors of every formatted element. If you don't specify colors (by leaving `Colours` as `nil`), zylog uses sensible defaults that match the current behavior.

### Simple Example - Changing a Few Colors

```go
import (
    "github.com/fatih/color"
    "github.com/zylisp/zylog"
    "github.com/zylisp/zylog/colors"
    "github.com/zylisp/zylog/options"
)

func main() {
    opts := options.Default()
    
    // Start with default colors
    opts.Colours = colors.Default()
    
    // Customize just the colors you want to change
    opts.Colours.LevelError = &colors.Color{
        Fg: color.FgHiRed,
        Bg: color.BgYellow,  // Add yellow background to errors
    }
    opts.Colours.Message = &colors.Color{
        Fg: color.FgHiWhite,
        Bg: color.Reset,  // No background
    }
    
    logger, _ := zylog.SetupLogging(opts)
    logger.Error("This error has a yellow background!")
}
```

### Disabling Color for Specific Elements

To disable color for a specific element while keeping others colored, set both Fg and Bg to `color.Reset`:

```go
opts.Colours = colors.Default()
opts.Colours.Timestamp = &colors.Color{
    Fg: color.Reset,
    Bg: color.Reset,
}
// Timestamp will now be uncolored, but everything else remains colored
```

### Complete Color Configuration Reference

The `Colours` struct provides fine-grained control over every colored element:

```go
type Colours struct {
    // Timestamp colors (default: HiBlack/grey)
    Timestamp *Color
    
    // Log level colors
    LevelTrace   *Color  // default: HiMagenta
    LevelDebug   *Color  // default: HiCyan
    LevelInfo    *Color  // default: HiGreen
    LevelWarn    *Color  // default: HiYellow
    LevelWarning *Color  // default: HiYellow
    LevelError   *Color  // default: Red
    LevelFatal   *Color  // default: HiRed
    LevelPanic   *Color  // default: HiWhite
    
    // Message text color (default: Green)
    Message *Color
    
    // Arrow separator " ▶ " (default: Cyan)
    Arrow *Color
    
    // Caller information colors
    CallerFunction *Color  // default: HiYellow
    CallerLine     *Color  // default: Yellow
    
    // Structured logging attribute colors
    AttrKey   *Color  // default: Yellow
    AttrValue *Color  // default: HiYellow
}

type Color struct {
    Fg color.Attribute  // Foreground color
    Bg color.Attribute  // Background color
}
```

Available color attributes from `github.com/fatih/color`:

**Foreground colors:**
- `color.FgBlack`, `color.FgRed`, `color.FgGreen`, `color.FgYellow`
- `color.FgBlue`, `color.FgMagenta`, `color.FgCyan`, `color.FgWhite`
- `color.FgHiBlack`, `color.FgHiRed`, `color.FgHiGreen`, `color.FgHiYellow`
- `color.FgHiBlue`, `color.FgHiMagenta`, `color.FgHiCyan`, `color.FgHiWhite`

**Background colors:**
- `color.BgBlack`, `color.BgRed`, `color.BgGreen`, `color.BgYellow`
- `color.BgBlue`, `color.BgMagenta`, `color.BgCyan`, `color.BgWhite`
- `color.BgHiBlack`, `color.BgHiRed`, `color.BgHiGreen`, `color.BgHiYellow`
- `color.BgHiBlue`, `color.BgHiMagenta`, `color.BgHiCyan`, `color.BgHiWhite`

**Special:**
- `color.Reset` - No color (use for both Fg and Bg to disable coloring for an element)

### Global Color Disable

The existing `Colored: false` option continues to work and will disable ALL colors regardless of individual color settings:

```go
opts := options.Default()
opts.Colored = false  // Disables all colors globally
opts.Colours = colors.Default()  // This will be ignored due to Colored: false
```
```

### Step 8: Update Version in Makefile

Change the VERSION line from:

```makefile
VERSION = 0.2.0
```

to:

```makefile
VERSION = 0.2.1
```

## Testing Checklist

After implementation, verify:

1. **Default behavior unchanged**: Running existing code without Colours specified should produce identical output
2. **Color customization works**: Setting custom colors for individual elements should apply correctly
3. **Nil handling**: `Colours: nil` should use default colors
4. **Per-element disable**: Setting `Fg: color.Reset, Bg: color.Reset` should disable color for that element only
5. **Global disable respected**: `Colored: false` should override all color settings
6. **Both backends**: Test with both slog and logrus
7. **Background colors**: Verify background colors are applied when set
8. **All elements**: Test customization of all color fields (timestamp, levels, message, arrow, caller, attributes)

## Implementation Order

Follow this order to minimize errors:

1. Create `colors/colors.go` with Color and Colours types
2. Add Colours field to options.ZyLog
3. Update formatter/common.go functions
4. Update formatter/formatter.go (LogLine struct and ColorLevel function)
5. Update logger/logrus.go to pass Colours through
6. Update logger/slog.go to pass Colours through
7. Update README.md with documentation
8. Update Makefile version
9. Run `make demo` to verify everything works

## Important Notes

- The global `opts.Colored bool` flag takes precedence - when false, it disables all colors via `color.NoColor = true`
- All format functions should check for nil colours and use `colors.Default()` as fallback
- Individual Color pointers can be nil within a Colours struct - treat nil as "use the element without color changes"
- The `color.Reset` constant (value 0) is used to explicitly disable color for Fg or Bg
- Background colors default to `color.Reset` (no background) in the Default() function
- Maintain backward compatibility: existing code without Colours should work identically

## Go Module Updates

No new dependencies are required - `github.com/fatih/color` is already in use.

## Files to Create/Modify

**Create:**
- `colors/colors.go`

**Modify:**
- `options/options.go`
- `formatter/common.go`
- `formatter/formatter.go`
- `logger/logrus.go`
- `logger/slog.go`
- `README.md`
- `Makefile`

## Validation

After implementation, run:
```bash
make demo
```

The output should be identical to the current output (since we're using default colors when Colours is nil). Then create a simple test program that sets custom colors to verify the new functionality works.
