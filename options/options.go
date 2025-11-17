// Package options defines configuration options for zylog.
package options

import (
	"fmt"

	"github.com/zylisp/zylog/colours"
	"github.com/zylisp/zylog/formatter"
)

// Logger represents the type of logger backend to use.
type Logger int

// Logger type constants
const (
	LogRUs Logger = iota // LogRUs uses the logrus logging library
	Slog                 // Slog uses Go's standard library slog
)

var (
	defaultOpts = &ZyLog{
		Coloured:        true,
		Level:           "trace",
		Output:          "stdout",
		ReportCaller:    true,
		TimestampFormat: formatter.SimpleTimestamp,
		PadLevel:        false,
		PadAmount:       5,
		PadSide:         "left",
		MsgSeparator:    ": ",
		Logger:          Slog,
		Colours:         colours.Default(),
	}
)

func (l Logger) String() string {
	switch l {
	case LogRUs:
		return "logrus"
	case Slog:
		return "slog"
	default:
		return fmt.Sprintf("unknown logger (iota '%d')", l)
	}
}

// ZyLog are used by the zylog logger to set up logrus.
type ZyLog struct {
	Coloured        bool
	Level           string
	Output          string // stdout, stderr, or filesystem
	ReportCaller    bool
	TimestampFormat formatter.TSFormat // RFC3339, Simple (YYYYMMDD.HHmmSS), or Time (HH:mm:SS)
	PadLevel        bool               // Whether to pad level strings for alignment
	PadAmount       int                // Number of characters to pad level strings to
	PadSide         string             // "left" or "right"; which side to pad level strings on
	MsgSeparator    string             // Separator between message and attributes
	Logger          Logger             // Logger type: Logrus or Slog
	Colours         *colours.Colours   // Colour configuration (nil uses defaults)
}

// Default returns the default ZyLog configuration options.
func Default() *ZyLog {
	opts := *defaultOpts
	return &opts
}

// WithLevelPadding returns ZyLog configuration options with PadLevel disabled.
func WithLevelPadding() *ZyLog {
	opts := *defaultOpts
	opts.PadLevel = true
	return &opts
}

// NoCaller returns ZyLog configuration options with ReportCaller disabled.
func NoCaller() *ZyLog {
	opts := *defaultOpts
	opts.ReportCaller = false
	return &opts
}
