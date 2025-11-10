// Package logger provides logging implementations for both slog and logrus with zylog formatting.
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/fatih/color"

	"github.com/zylisp/zylog/errors"
	"github.com/zylisp/zylog/formatter"
	"github.com/zylisp/zylog/level"
	"github.com/zylisp/zylog/options"
)

// SLogHandler implements slog.Handler with zylog formatting.
type SLogHandler struct {
	opts   *options.ZyLog
	writer io.Writer
	attrs  []slog.Attr
	groups []string
}

// SetupSlog configures and returns a new slog.Logger with zylog formatting.
// Returns a configured *slog.Logger instance that can be used directly or set as the default logger.
func SetupSlog(opts *options.ZyLog) *slog.Logger {
	if opts == nil {
		opts = options.Default()
	}

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
	opts.TimestampFormat = timestampFormat

	// 4. Create handler
	h := NewSLogHandler(writer, opts)

	// 5. Create and return logger
	logger := slog.New(h)
	slog.SetDefault(logger)
	logger.Info("Slog logging initialized.")
	return logger
}

// NewSLogHandler creates a new SLogHandler with the given writer and options.
func NewSLogHandler(writer io.Writer, opts *options.ZyLog) *SLogHandler {
	if opts == nil {
		opts = options.Default()
	}
	return &SLogHandler{
		opts:   opts,
		writer: writer,
		attrs:  make([]slog.Attr, 0),
		groups: make([]string, 0),
	}
}

// Enabled reports whether the handler handles records at the given level.
func (h *SLogHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	minLevel := parseSlogLevel(h.opts.Level)
	return lvl >= minLevel
}

// Handle handles the Record.
func (h *SLogHandler) Handle(_ context.Context, r slog.Record) error {
	// Build the log line using the same format as formatter.LogLine
	var buf strings.Builder

	// 1. Format timestamp
	timestampStr := r.Time.Format(h.opts.TimestampFormat.ToTimeFormat())
	buf.WriteString(formatter.FormatTimestamp(timestampStr))
	buf.WriteString(" ")

	// 2. Format level
	levelStr := slogLevelToString(r.Level)
	levelFormatted := formatter.ColorLevel(levelStr, h.opts.PadLevel, h.opts.PadAmount, h.opts.PadSide)
	buf.WriteString(levelFormatted)

	// 3. Format caller if enabled
	if h.opts.ReportCaller && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		buf.WriteString(formatter.FormatCaller(f.Function, f.Line))
	}

	// 4. Format message
	if r.Message != "" {
		buf.WriteString(formatter.FormatArrow())
		buf.WriteString(formatter.FormatMessage(r.Message))
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
			h.appendAttr(&buf, attr)
			first = false
		}

		// Add record-level attributes
		r.Attrs(func(a slog.Attr) bool {
			if !first {
				buf.WriteString(", ")
			}
			h.appendAttr(&buf, a)
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

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
func (h *SLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Create a new handler with cloned attributes
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &SLogHandler{
		opts:   h.opts,
		writer: h.writer,
		attrs:  newAttrs,
		groups: h.groups, // TODO: implement group support if needed
	}
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
func (h *SLogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = name

	return &SLogHandler{
		opts:   h.opts,
		writer: h.writer,
		attrs:  h.attrs,
		groups: newGroups,
	}
}

// appendAttr appends a single attribute to the buffer in zylog format.
func (h *SLogHandler) appendAttr(buf *strings.Builder, attr slog.Attr) {
	// Handle groups
	prefix := ""
	if len(h.groups) > 0 {
		prefix = strings.Join(h.groups, ".") + "."
	}

	key := prefix + attr.Key
	value := attr.Value.String()

	fmt.Fprintf(buf, "%s={%s}", formatter.FormatAttrKey(key), formatter.FormatAttrValue(value))
}

// slogLevelToString converts a slog.Level to a zylog level string.
func slogLevelToString(lvl slog.Level) string {
	switch {
	case lvl < slog.LevelDebug:
		return level.Trace
	case lvl < slog.LevelInfo:
		return level.Debug
	case lvl < slog.LevelWarn:
		return level.Info
	case lvl < slog.LevelError:
		return level.Warn
	case lvl < slog.LevelError+4:
		return level.Error
	case lvl < slog.LevelError+8:
		return level.Fatal
	default:
		return level.Panic
	}
}

// parseSlogLevel converts a string level to slog.Level.
func parseSlogLevel(levelStr string) slog.Level {
	switch strings.ToLower(levelStr) {
	case "trace":
		return slog.LevelDebug - 1
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "fatal":
		return slog.LevelError + 4
	case "panic":
		return slog.LevelError + 8
	default:
		return slog.LevelInfo
	}
}
