// Package formatter provides custom log formatting functionality for zylog.
package formatter

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/zylog/level"
)

// TSFormat represents a timestamp format type.
type TSFormat int

// Timestamp format constants
const (
	TSUnset TSFormat = iota
	RFC3339
	SimpleTimestamp
	TimeOnly

	TSSimple   = "20060102.150405"
	TSTimeOnly = "15:04:05"
)

// ToTimeFormat converts a Format to its corresponding time format string.
func (f TSFormat) ToTimeFormat() string {
	switch f {
	case RFC3339:
		return time.RFC3339
	case SimpleTimestamp:
		return TSSimple
	case TimeOnly:
		return TSTimeOnly
	default:
		return TSSimple // Default to Simple
	}
}

// LogLine formats logs into a complete line.
type LogLine struct {
	// Force disabling colors.
	DisableColors bool
	// TimestampFormat specifies the format for timestamps.
	TimestampFormat TSFormat
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
func (f *LogLine) Format(entry *log.Entry) ([]byte, error) {
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
func ColorLevel(lvl string) string {
	switch lvl {
	case level.Trace:
		lvl = color.HiMagentaString(lvl)
	case level.Debug:
		lvl = color.HiCyanString(lvl)
	case level.Info:
		lvl = color.HiGreenString(lvl)
	case level.Warn:
		lvl = color.HiYellowString(lvl)
	case level.Error:
		lvl = color.RedString(lvl)
	case level.Fatal:
		lvl = color.HiRedString(lvl)
	case level.Panic:
		lvl = color.HiWhiteString(lvl)
	}
	return lvl
}
