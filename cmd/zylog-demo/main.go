/*
Package main offers a demo utility for the zylog logger wrapper.

Log entries with both caller (package, function, and line number) as well as
without caller information are demonstrated.
*/
package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/zylisp/zylog"
	"github.com/zylisp/zylog/formatter"
	"github.com/zylisp/zylog/options"
)

// SetupLogger ...
func SetupLogger() {
	zylog.SetupLogging(options.Default())
}

// SetupLoggerNoCaller ...
func SetupLoggerNoCaller() {
	zylog.SetupLogging(options.NoCaller())
}

// SetupLoggerNoPad ...
func SetupLoggerNoPad() {
	opts := options.NoLevelPadding()
	opts.TimestampFormat = formatter.RFC3339
	zylog.SetupLogging(opts)
}

// SetupLoggerNoColour ...
func SetupLoggerNoColour() {
	opts := options.NoLevelPadding()
	opts.TimestampFormat = formatter.TimeOnly
	opts.Colored = false
	zylog.SetupLogging(opts)
}

func main() {
	zylog.PrintVersions()
	SetupLogger()
	log.Info(" *** Logging with default options ***")
	log.Trace("This is trace")
	log.Debug("This is debug")
	log.Info("This is info")
	log.Warn("This is warn")
	log.Error("This is error")
	log.Info("Fatal and Panic are also supported; " +
		"Fatal will os.Exit, and Panic will log, then panic().")
	log.Info("When not testing, you'll want to turn off caller reporting:")
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
}
