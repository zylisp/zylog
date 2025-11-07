package zylog

import (
	"log/slog"

	"github.com/zylisp/zylog/errors"
	"github.com/zylisp/zylog/logger"
	"github.com/zylisp/zylog/options"
)

// SetupLogging configures the selected logger (slog or logrus) with zylog formatting.
// For slog, it returns the configured logger instance and sets it as default.
// For logrus, it configures the global logrus instance and returns nil.
func SetupLogging(opts *options.ZyLog) (*slog.Logger, error) {
	switch opts.Logger {
	case options.LogRUs:
		logger.SetupLogRUs(opts)
		return nil, nil
	case options.Slog:
		l := logger.SetupSlog(opts)
		slog.SetDefault(l)
		return l, nil
	default:
		return nil, errors.ErrUnsupLogger(opts.Logger)
	}
}
