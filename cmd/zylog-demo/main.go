/*
Package main offers a demo utility for the zylog logger wrapper.

Log entries with both caller (package, function, and line number) as well as
without caller information are demonstrated.
*/
package main

import (
	"log/slog"

	"github.com/sirupsen/logrus"

	"github.com/zylisp/zylog"
	"github.com/zylisp/zylog/formatter"
	"github.com/zylisp/zylog/options"
)

// SetupLogger ...
func SetupLogger(ol options.Logger) *slog.Logger {
	opts := options.Default()
	opts.Logger = ol
	l, _ := zylog.SetupLogging(opts)
	return l
}

// SetupLoggerNoCaller ...
func SetupLoggerNoCaller(ol options.Logger) *slog.Logger {
	opts := options.NoCaller()
	opts.Logger = ol
	l, _ := zylog.SetupLogging(opts)
	return l
}

// SetupLoggerWithPad ...
func SetupLoggerWithPad(ol options.Logger) *slog.Logger {
	opts := options.WithLevelPadding()
	if ol == options.LogRUs {
		opts.PadAmount = 7
	}
	opts.Logger = ol
	opts.TimestampFormat = formatter.RFC3339
	l, _ := zylog.SetupLogging(opts)
	return l
}

// SetupLoggerNoColour ...
func SetupLoggerNoColour(ol options.Logger) *slog.Logger {
	opts := options.Default()
	if ol == options.LogRUs {
		opts.PadAmount = 7
	}
	opts.Logger = ol
	opts.TimestampFormat = formatter.TimeOnly
	opts.Coloured = false
	l, _ := zylog.SetupLogging(opts)
	return l
}

func main() {
	zylog.PrintVersions()

	slog.Info("=======================================")
	slog.Info("=======   SLOG DEMONSTRATIONS   =======")
	slog.Info("=======================================")

	// Demo 1: Default options
	log := SetupLogger(options.Slog)
	log.Info(" *** Slog with default options ***")
	log.Debug("This is debug")
	log.Info("This is info")
	log.Warn("This is warn")
	log.Error("This is error")

	// Demo 2: No caller
	log = SetupLoggerNoCaller(options.Slog)
	log.Info(" *** Slog with calling function disabled ***")
	log.Debug("This is debug")
	log.Info("This is info")
	log.Warn("This is warn")
	log.Error("This is error")

	// Demo 3: With padding, RFC3339
	log = SetupLoggerWithPad(options.Slog)
	log.Info(" *** Slog with level padding enabled and timestamps in RFC3339 format ***")
	log.Debug("This is debug")
	log.Info("This is info")
	log.Warn("This is warn")
	log.Error("This is error")
	log.Info("Structured logging example",
		slog.String("user", "alice"),
		slog.String("request_id", "12345"),
		slog.Int("duration_ms", 42))
	userLogger := log.With(slog.String("user", "bob"), slog.String("role", "admin"))
	userLogger.Info("User performed action", slog.String("action", "delete"))
	userLogger.Warn("User attempted restricted operation", slog.String("action", "access_logs"))

	// Demo 4: No colour, time-only
	log = SetupLoggerNoColour(options.Slog)
	log.Info(" *** Slog with timestamps in time-only format, and no colour ***")
	log.Debug("This is debug")
	log.Info("This is info")
	log.Warn("This is warn")
	log.Error("This is error")

	// === LOGRUS DEMOS ===
	logrus.Info("=======================================")
	logrus.Info("======   LOGRUS DEMONSTRATIONS   ======")
	logrus.Info("=======================================")

	_ = SetupLogger(options.LogRUs)
	logrus.Info(" *** Logging with default options ***")
	logrus.Trace("This is trace")
	logrus.Debug("This is debug")
	logrus.Info("This is info")
	logrus.Warn("This is warn")
	logrus.Error("This is error")
	logrus.Info("Fatal and Panic are also supported; " +
		"Fatal will os.Exit, and Panic will log, then panic().")

	_ = SetupLoggerNoCaller(options.LogRUs)
	logrus.Info(" *** Logging with calling function disabled ***")
	logrus.Trace("This is trace")
	logrus.Debug("This is debug")
	logrus.Info("This is info")
	logrus.Warn("This is warn")
	logrus.Error("This is error")

	_ = SetupLoggerWithPad(options.LogRUs)
	logrus.Info(" *** Logging with level padding rnabled and timestamps in RFC3339 format ***")
	logrus.Trace("This is trace")
	logrus.Debug("This is debug")
	logrus.Info("This is info")
	logrus.Warn("This is warn")
	logrus.Error("This is error")

	_ = SetupLoggerNoColour(options.LogRUs)
	logrus.Info(" *** Logging with timestamps in time-only format, and no colour ***")
	logrus.Trace("This is trace")
	logrus.Debug("This is debug")
	logrus.Info("This is info")
	logrus.Warn("This is warn")
	logrus.Error("This is error")
}
