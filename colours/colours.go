// Package colours provides colour configuration types for zylog.
package colours

import "github.com/fatih/color"

// Colour represents foreground and background colour attributes for a single element.
// Use color.Attribute constants from github.com/fatih/color (e.g., color.FgRed, color.BgBlue).
// Set to color.Reset (or 0) to disable colour for that component.
type Colour struct {
	Fg color.Attribute // Foreground colour
	Bg color.Attribute // Background colour
}

// Colours holds colour configuration for all formatted output elements in zylog.
// A nil Colours pointer will use default colours (current hardcoded behavior).
type Colours struct {
	// Timestamp colours (currently light grey/HiBlack)
	Timestamp *Colour

	// Log level colours
	LevelTrace   *Colour // Currently HiMagenta
	LevelDebug   *Colour // Currently HiCyan
	LevelInfo    *Colour // Currently HiGreen
	LevelWarn    *Colour // Currently HiYellow
	LevelWarning *Colour // Currently HiYellow (alias for Warn)
	LevelError   *Colour // Currently Red
	LevelFatal   *Colour // Currently HiRed
	LevelPanic   *Colour // Currently HiWhite

	// Message text colour (currently green)
	Message *Colour

	// Arrow separator (currently cyan " ▶ ")
	Arrow *Colour

	// Caller information colours
	CallerFunction *Colour // Currently HiYellow
	CallerLine     *Colour // Currently Yellow

	// Structured logging attribute colours
	AttrKey   *Colour // Currently Yellow
	AttrValue *Colour // Currently HiYellow
}

// Default returns the default colour configuration matching current hardcoded behavior.
func Default() *Colours {
	return &Colours{
		Timestamp:      &Colour{Fg: color.FgHiBlack, Bg: color.Reset},
		LevelTrace:     &Colour{Fg: color.FgHiMagenta, Bg: color.Reset},
		LevelDebug:     &Colour{Fg: color.FgHiCyan, Bg: color.Reset},
		LevelInfo:      &Colour{Fg: color.FgHiGreen, Bg: color.Reset},
		LevelWarn:      &Colour{Fg: color.FgHiYellow, Bg: color.Reset},
		LevelWarning:   &Colour{Fg: color.FgHiYellow, Bg: color.Reset},
		LevelError:     &Colour{Fg: color.FgRed, Bg: color.Reset},
		LevelFatal:     &Colour{Fg: color.FgHiRed, Bg: color.Reset},
		LevelPanic:     &Colour{Fg: color.FgHiWhite, Bg: color.Reset},
		Message:        &Colour{Fg: color.FgGreen, Bg: color.Reset},
		Arrow:          &Colour{Fg: color.FgCyan, Bg: color.Reset},
		CallerFunction: &Colour{Fg: color.FgHiYellow, Bg: color.Reset},
		CallerLine:     &Colour{Fg: color.FgYellow, Bg: color.Reset},
		AttrKey:        &Colour{Fg: color.FgYellow, Bg: color.Reset},
		AttrValue:      &Colour{Fg: color.FgGreen, Bg: color.Reset},
	}
}

// ApplyColour applies the colour to a string. If colour is nil, returns string unchanged.
// If both Fg and Bg are color.Reset (0), returns string unchanged (no colour).
func (c *Colour) ApplyColour(s string) string {
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
