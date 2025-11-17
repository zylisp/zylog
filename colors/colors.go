// Package colors provides color configuration types for zylog.
package colors

import "github.com/fatih/color"

// Color represents foreground and background color attributes for a single element.
// Use color.Attribute constants from github.com/fatih/color (e.g., color.FgRed, color.BgBlue).
// Set to color.Reset (or 0) to disable color for that component.
type Color struct {
	Fg color.Attribute // Foreground color
	Bg color.Attribute // Background color
}

// Colours holds color configuration for all formatted output elements in zylog.
// A nil Colours pointer will use default colors (current hardcoded behavior).
type Colours struct {
	// Timestamp colors (currently light grey/HiBlack)
	Timestamp *Color

	// Log level colors
	LevelTrace   *Color // Currently HiMagenta
	LevelDebug   *Color // Currently HiCyan
	LevelInfo    *Color // Currently HiGreen
	LevelWarn    *Color // Currently HiYellow
	LevelWarning *Color // Currently HiYellow (alias for Warn)
	LevelError   *Color // Currently Red
	LevelFatal   *Color // Currently HiRed
	LevelPanic   *Color // Currently HiWhite

	// Message text color (currently green)
	Message *Color

	// Arrow separator (currently cyan " ▶ ")
	Arrow *Color

	// Caller information colors
	CallerFunction *Color // Currently HiYellow
	CallerLine     *Color // Currently Yellow

	// Structured logging attribute colors
	AttrKey   *Color // Currently Yellow
	AttrValue *Color // Currently HiYellow
}

// Default returns the default color configuration matching current hardcoded behavior.
func Default() *Colours {
	return &Colours{
		Timestamp:      &Color{Fg: color.FgHiBlack, Bg: color.Reset},
		LevelTrace:     &Color{Fg: color.FgHiMagenta, Bg: color.Reset},
		LevelDebug:     &Color{Fg: color.FgHiCyan, Bg: color.Reset},
		LevelInfo:      &Color{Fg: color.FgHiGreen, Bg: color.Reset},
		LevelWarn:      &Color{Fg: color.FgHiYellow, Bg: color.Reset},
		LevelWarning:   &Color{Fg: color.FgHiYellow, Bg: color.Reset},
		LevelError:     &Color{Fg: color.FgRed, Bg: color.Reset},
		LevelFatal:     &Color{Fg: color.FgHiRed, Bg: color.Reset},
		LevelPanic:     &Color{Fg: color.FgHiWhite, Bg: color.Reset},
		Message:        &Color{Fg: color.FgGreen, Bg: color.Reset},
		Arrow:          &Color{Fg: color.FgCyan, Bg: color.Reset},
		CallerFunction: &Color{Fg: color.FgHiYellow, Bg: color.Reset},
		CallerLine:     &Color{Fg: color.FgYellow, Bg: color.Reset},
		AttrKey:        &Color{Fg: color.FgYellow, Bg: color.Reset},
		AttrValue:      &Color{Fg: color.FgHiYellow, Bg: color.Reset},
	}
}

// ApplyColor applies the color to a string. If color is nil, returns string unchanged.
// If both Fg and Bg are color.Reset (0), returns string unchanged (no color).
func (c *Color) ApplyColor(s string) string {
	if c == nil {
		return s
	}
	if c.Fg == color.Reset && c.Bg == color.Reset {
		return s
	}

	// Build attribute list
	attrs := []color.Attribute{}
	if c.Fg != color.Reset {
		attrs = append(attrs, c.Fg)
	}
	if c.Bg != color.Reset {
		attrs = append(attrs, c.Bg)
	}

	if len(attrs) == 0 {
		return s
	}

	return color.New(attrs...).Sprint(s)
}
