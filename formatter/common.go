package formatter

import (
	"fmt"

	"github.com/zylisp/zylog/colours"
)

// FormatTimestamp formats a time string with the configured colour.
func FormatTimestamp(timestamp string, cols *colours.Colours) string {
	return cols.Timestamp.ApplyColour(timestamp)
}

// FormatMessage formats a log message with the configured colour.
func FormatMessage(message string, cols *colours.Colours) string {
	return cols.Message.ApplyColour(message)
}

// FormatArrow returns the coloured arrow separator.
func FormatArrow(cols *colours.Colours) string {
	return cols.Arrow.ApplyColour(" ▶ ")
}

// FormatCaller formats caller information with the configured colours.
func FormatCaller(function string, line int, cols *colours.Colours) string {
	functionStr := cols.CallerFunction.ApplyColour(function)
	lineStr := cols.CallerLine.ApplyColour(fmt.Sprintf("%d", line))
	return fmt.Sprintf(" [%s:%s]", functionStr, lineStr)
}

// FormatAttrKey formats an attribute key with the configured colour.
func FormatAttrKey(key string, cols *colours.Colours) string {
	return cols.AttrKey.ApplyColour(key)
}

// FormatAttrValue formats an attribute value with the configured colour.
func FormatAttrValue(value string, cols *colours.Colours) string {
	return cols.AttrValue.ApplyColour(value)
}
