// Package options defines configuration options for zylog.
package options

import "github.com/zylisp/zylog/formatter"

var (
	defaultOpts = &ZyLog{
		Colored:      true,
		Level:        "trace",
		Output:       "stdout",
		ReportCaller: true,
		PadLevel:     true,
		PadSide:      "left",
	}
)

// ZyLog are used by the zylog logger to set up logrus.
type ZyLog struct {
	Colored         bool
	Level           string
	Output          string // stdout, stderr, or filesystem
	ReportCaller    bool
	TimestampFormat formatter.TSFormat // RFC3339, Simple (YYYYMMDD.HHmmSS), or Time (HH:mm:SS); defaults to Simple
	PadLevel        bool               // Whether to pad level strings for alignment
	PadSide         string             // "left" or "right"; which side to pad level strings on
}

// Default returns the default ZyLog configuration options.
func Default() *ZyLog {
	return defaultOpts
}

// NoCaller returns ZyLog configuration options with ReportCaller disabled.
func NoCaller() *ZyLog {
	noCaller := defaultOpts
	noCaller.ReportCaller = false
	return noCaller
}
