package formatter

import (
	"fmt"

	"github.com/fatih/color"
)

// FormatTimestamp formats a time string with light grey color.
func FormatTimestamp(timestamp string) string {
	return color.HiBlackString(timestamp)
}

// FormatMessage formats a log message with dark green color.
func FormatMessage(message string) string {
	return color.GreenString(message)
}

// FormatArrow returns the colored arrow separator.
func FormatArrow() string {
	return color.CyanString(" ▶ ")
}

// FormatCaller formats caller information with yellow colors.
func FormatCaller(function string, line int) string {
	return fmt.Sprintf(" [%s:%s]",
		color.HiYellowString(function),
		color.YellowString(fmt.Sprintf("%d", line)))
}

// FormatAttrKey formats an attribute key with yellow color.
func FormatAttrKey(key string) string {
	return color.YellowString(key)
}

// FormatAttrValue formats an attribute value with bright yellow color.
func FormatAttrValue(value string) string {
	return color.HiYellowString(value)
}
