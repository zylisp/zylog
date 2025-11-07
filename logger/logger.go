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

	$ go get github.com/zylisp/zylog/logger

Additionally, there is a demo you may install and run:

	$ go get github.com/zylisp/zylog/cmd/zylog-demo

# Configuration

To configure the logger, simply pass an options struct reference to
SetupLogging. For example,

package main

	import (
		logger "github.com/zylisp/zylog/logger"
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
