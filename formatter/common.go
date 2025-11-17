package formatter

import (
	"fmt"

	"github.com/zylisp/zylog/colors"
)

// FormatTimestamp formats a time string with the configured color.
func FormatTimestamp(timestamp string, colours *colors.Colours) string {
	return colours.Timestamp.ApplyColor(timestamp)
}

// FormatMessage formats a log message with the configured color.
func FormatMessage(message string, colours *colors.Colours) string {
	return colours.Message.ApplyColor(message)
}

// FormatArrow returns the colored arrow separator.
func FormatArrow(colours *colors.Colours) string {
	return colours.Arrow.ApplyColor(" ▶ ")
}

// FormatCaller formats caller information with the configured colors.
func FormatCaller(function string, line int, colours *colors.Colours) string {
	functionStr := colours.CallerFunction.ApplyColor(function)
	lineStr := colours.CallerLine.ApplyColor(fmt.Sprintf("%d", line))
	return fmt.Sprintf(" [%s:%s]", functionStr, lineStr)
}

// FormatAttrKey formats an attribute key with the configured color.
func FormatAttrKey(key string, colours *colors.Colours) string {
	return colours.AttrKey.ApplyColor(key)
}

// FormatAttrValue formats an attribute value with the configured color.
func FormatAttrValue(value string, colours *colors.Colours) string {
	return colours.AttrValue.ApplyColor(value)
}
