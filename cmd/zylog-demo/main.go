/*
Package main offers a demo utility for the zylog logger wrapper.

Log entries with both caller (package, function, and line number) as well as
without caller information are demonstrated.
*/
package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/zylog"
	"github.com/geomyidia/zylog/options"
)

// SetupLogger ...
func SetupLogger() {
	zylog.SetupLogging(options.Default())
}

// SetupLoggerNoCaller ...
func SetupLoggerNoCaller() {
	zylog.SetupLogging(options.NoCaller())
}

func main() {
	zylog.PrintVersions()
	SetupLogger()
	log.Trace("This is trace")
	log.Debug("This is debug")
	log.Info("This is info")
	log.Warn("This is warn")
	log.Error("This is error")
	log.Info("Fatal and Panic are also supported; " +
		"Fatal will os.Exit, and Panic will log, then panic().")
	log.Info("When not testing, you'll want to turn off caller reporting:")
	SetupLoggerNoCaller()
	log.Trace("This is trace")
	log.Debug("This is debug")
	log.Info("This is info")
	log.Warn("This is warn")
	log.Error("This is error")
}
