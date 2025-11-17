package logger

import (
	"os"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"

	"github.com/zylisp/zylog/errors"
	"github.com/zylisp/zylog/formatter"
	"github.com/zylisp/zylog/options"
)

// Output destination constants
const (
	StdOut     = "stdout"
	StdErr     = "stderr"
	FileSystem = "filesystem"
)

// SetupLogRUs performs the setup of the logrus logger with zylog formatting.
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
		Colours:         opts.Colours,
	})
	log.SetReportCaller(opts.ReportCaller)
	log.Info("Logging initialized.")
}
