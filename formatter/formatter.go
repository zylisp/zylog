// Package formatter provides custom log formatting functionality for zylog.
package formatter

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"

	"github.com/zylisp/zylog/level"
)

// TSFormat represents a timestamp format type.
type TSFormat int

// Timestamp format constants
const (
	TSUnset TSFormat = iota
	RFC3339
	StandardTimestamp
	SimpleTimestamp
	TimeOnly

	TSSimple   = "20060102.150405"
	TSStandard = "2006/01/02 15:04:05"
	TSTimeOnly = "15:04:05"
)

// ToTimeFormat converts a Format to its corresponding time format string.
func (f TSFormat) ToTimeFormat() string {
	switch f {
	case RFC3339:
		return time.RFC3339
	case StandardTimestamp:
		return TSStandard
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
	// PadLevel specifies whether to pad level strings for alignment.
	PadLevel bool
	// PadAmount specifies the total width to pad level strings to.
	PadAmount int
	// PadSide specifies which side to add padding on ("left" or "right").
	PadSide string
	// AttrSeparator specifies the separator between message and attributes.
	MsgSeparator string
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

	timestamp := FormatTimestamp(entry.Time.Format(f.TimestampFormat.ToTimeFormat()))
	level := ColorLevel(strings.ToUpper(entry.Level.String()), f.PadLevel, f.PadAmount, f.PadSide)

	fmt.Fprintf(b, "%s %s", timestamp, level)
	if entry.Logger.ReportCaller {
		b.WriteString(FormatCaller(entry.Caller.Function, entry.Caller.Line))
	}
	if entry.Message != "" {
		b.WriteString(FormatArrow())
		b.WriteString(FormatMessage(entry.Message))
	}

	if len(entry.Data) > 0 {
		b.WriteString(f.MsgSeparator)
		first := true
		for key, value := range entry.Data {
			if !first {
				b.WriteString(", ")
			}
			fmt.Fprintf(b, "%s={%s}", FormatAttrKey(key), FormatAttrValue(fmt.Sprintf("%v", value)))
			first = false
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

// ColorLevel determines the color of the log level based upon the string
// value of the log level. If padLevel is true, the level string will be
// padded to padAmount characters, aligned according to padSide.
func ColorLevel(lvl string, padLevel bool, padAmount int, padSide string) string {
	// Apply padding before colorizing
	if padLevel && padAmount > 0 {
		if padSide == "left" {
			// Right-align: add spaces on the left
			lvl = fmt.Sprintf("%*s", padAmount, lvl)
		} else {
			// Left-align (default): add spaces on the right
			lvl = fmt.Sprintf("%-*s", padAmount, lvl)
		}
	}

	// Now colorize the padded string
	switch strings.TrimSpace(lvl) {
	case level.Trace:
		lvl = color.HiMagentaString(lvl)
	case level.Debug:
		lvl = color.HiCyanString(lvl)
	case level.Info:
		lvl = color.HiGreenString(lvl)
	case level.Warn, level.Warning:
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
