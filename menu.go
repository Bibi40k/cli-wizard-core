package wizard

import (
	"fmt"
	"regexp"
	"strings"
)

var ansiEscapeRE = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)

// FormatMenuLabel renders a two-column menu label: [tag] + aligned text.
// width controls the fixed width for the [tag] column; values <= 0 default to 12.
func FormatMenuLabel(tag, text string, width int) string {
	if width <= 0 {
		width = 12
	}
	left := ""
	if tag != "" {
		left = "[" + tag + "]"
	}
	return fmt.Sprintf("%-*s %s", width, left, text)
}

// Colorize wraps a string with an ANSI color code and resets formatting.
// Pass empty color to keep the text unchanged.
func Colorize(text, color string) string {
	if color == "" {
		return text
	}
	return color + text + "\033[0m"
}

// BackLabel returns a consistently colored "Back" label.
func BackLabel() string {
	return Colorize("Back", "\033[33m")
}

// BackMenuLabel renders a colored Back label aligned like menu entries.
func BackMenuLabel(width int) string {
	return Colorize(FormatMenuLabel("", "Back", width), "\033[33m")
}

// NormalizeChoice strips ANSI colors and trims whitespace for robust comparisons.
func NormalizeChoice(value string) string {
	return strings.TrimSpace(ansiEscapeRE.ReplaceAllString(value, ""))
}

// IsBackChoice checks whether a selected value maps to Back.
func IsBackChoice(value string) bool {
	return strings.EqualFold(NormalizeChoice(value), "Back")
}
