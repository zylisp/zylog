package zylog

import (
	"os"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/zylog/errors"
	"github.com/geomyidia/zylog/formatter"
	"github.com/geomyidia/zylog/options"
)

// Output destination constants
const (
	StdOut     = "stdout"
	StdErr     = "stderr"
	FileSystem = "filesystem"
)

// SetupLogging performs the setup of the zylog logger.
func SetupLogging(opts *options.ZyLog) {
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
	})
	log.SetReportCaller(opts.ReportCaller)
	log.Info("Logging initialized.")
}
