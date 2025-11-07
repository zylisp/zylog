/*
Package logger performs basic setup of the logrus library with custom formatting.

# Overview

Zylog logger's primary features include:

  - Exceedingly simple setup
  - Colored output (enabled/disabled with a boolean)
  - Logging level (lower-case string)
  - Output (only stdout and stderr currently supported)
  - ReportCaller (enabled/disabled with a boolean; prints package, function
    and line number)
  - Custom format (similar to the Clojure twig library and the LFE logjam
    libraries)

Setup is done with the zylog logger, after which logrus may be used as designed
by its author.

Installation

	$ go get github.com/geomyidia/zylog/logger

Additionally, there is a demo you may install and run:

	$ go get github.com/geomyidia/zylog/cmd/zylog-demo

# Configuration

To configure the logger, simply pass an options struct reference to
SetupLogging. For example,

package main

	import (
		logger "github.com/geomyidia/zylog/logger"
		log "github.com/sirupsen/logrus"
	)

	func main () {
		log.SetupLogging(&log.ZyLogOptions{
			Colored:      true,
			Level:        "info",
			Output:       "stdout",
			ReportCaller: false,
		})
		// More app code
		log.Info("App started up!")
	}
*/
package logger

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/zylog/errors"
)

// Format represents a timestamp format type
type Format int

// Timestamp format constants
const (
	Unset Format = iota
	RFC3339
	Simple
	Time
)

// TextFormatter formats logs into text.
type TextFormatter struct {
	// Force disabling colors.
	DisableColors bool
	// TimestampFormat specifies the format for timestamps.
	TimestampFormat Format
}

// ZyLogOptions are used by the zylog logger to set up logrus.
type ZyLogOptions struct {
	Colored         bool
	Level           string
	Output          string // stdout, stderr, or filesystem
	ReportCaller    bool
	TimestampFormat Format // RFC3339, Simple (YYYYMMDD.HHmmSS), or Time (HH:mm:SS); defaults to Simple
}

// ToTimeFormat converts a Format to its corresponding time format string.
func (f Format) ToTimeFormat() string {
	switch f {
	case RFC3339:
		return time.RFC3339
	case Simple:
		return "20060102.150405"
	case Time:
		return "15:04:05"
	case Unset:
		return "20060102.150405" // Default to Simple
	default:
		return "20060102.150405" // Default to Simple
	}
}

// SetupLogging performs the setup of the zylog logger.
func SetupLogging(opts *ZyLogOptions) {
	level, err := log.ParseLevel(opts.Level)
	if err != nil {
		panic(errors.ErrLogLevel)
	}
	log.SetLevel(level)
	switch opts.Output {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	case "filesystem":
		panic(errors.ErrNotImplemented("filesystem log output"))
	default:
		panic(errors.ErrLogOutput(opts.Output))
	}
	disableColors := !opts.Colored
	color.NoColor = disableColors
	timestampFormat := opts.TimestampFormat
	if timestampFormat == Unset {
		// Default to Simple if not set
		timestampFormat = Simple
	}
	log.SetFormatter(&TextFormatter{
		DisableColors:   disableColors,
		TimestampFormat: timestampFormat,
	})
	log.SetReportCaller(opts.ReportCaller)
	log.Info("Logging initialized.")
}

// Format provides the custom formatting of the zylog logger.
//
// In particular, logs output in the following form:
//
//	YYYY-mm-DDTHH:MM:SS-TZ:00 LEVEL ▶ logged message ...
//
// If the ReportCaller option is set to true, the log output will have the
// following form:
//
//	YYYY-mm-DDTHH:MM:SS-TZ:00 LEVEL [pkghost/auth/proj/file.Func:LINENUM] ▶ logged message ...
//
// Any structured data passed as logrus fields will be appended to the above
// line forms.
func (f *TextFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	time := color.HiBlackString(entry.Time.Format(f.TimestampFormat.ToTimeFormat()))
	level := ColorLevel(strings.ToUpper(entry.Level.String()))

	fmt.Fprintf(b, "%s %s", time, level)
	if entry.Logger.ReportCaller {
		fmt.Fprintf(b, " [%s:%s]",
			color.HiYellowString(entry.Caller.Function),
			color.YellowString(strconv.Itoa(entry.Caller.Line)))
	}
	if entry.Message != "" {
		b.WriteString(color.CyanString(" ▶ "))
		b.WriteString(color.GreenString(entry.Message))
	}

	if len(entry.Data) > 0 {
		b.WriteString(" || ")
	}
	for key, value := range entry.Data {
		fmt.Fprintf(b, "%s={%s}, ", key, value)
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

// ColorLevel determines the color of the log level based upon the string
// value of the log level.
func ColorLevel(level string) string {
	switch level {
	case "TRACE":
		level = color.HiMagentaString(level)
	case "DEBUG":
		level = color.HiCyanString(level)
	case "INFO":
		level = color.HiGreenString(level)
	case "WARNING":
		level = color.HiYellowString(level)
	case "ERROR":
		level = color.RedString(level)
	case "FATAL":
		level = color.HiRedString(level)
	case "PANIC":
		level = color.HiWhiteString(level)
	}
	return level
}
