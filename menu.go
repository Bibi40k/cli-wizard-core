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
	return "Back"
}

// BackMenuLabel renders a colored Back label aligned like menu entries.
func BackMenuLabel(width int) string {
	return FormatMenuLabel("", "Back", width)
}

// ExitLabel returns a normalized "Exit" label.
func ExitLabel() string {
	return "Exit"
}

// ExitMenuLabel renders an Exit label aligned like menu entries.
func ExitMenuLabel(width int) string {
	return FormatMenuLabel("", "Exit", width)
}

// SelectHint returns the standard selector help hint used across wizard UIs.
func SelectHint() string {
	return "[Use arrows to move, type to filter, Esc=Back, Ctrl+C=Exit]"
}

// FormatActionLabel aligns action menus as two columns: verb + description.
// "Back"/"Exit" remain plain.
func FormatActionLabel(text string, verbWidth int) string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return trimmed
	}
	if IsBackChoice(trimmed) || IsExitChoice(trimmed) {
		return NormalizeChoice(trimmed)
	}
	if verbWidth <= 0 {
		verbWidth = 9
	}
	fields := strings.Fields(trimmed)
	if len(fields) < 2 {
		return trimmed
	}
	verb := fields[0]
	rest := strings.TrimSpace(strings.TrimPrefix(trimmed, verb))
	return fmt.Sprintf("%-*s %s", verbWidth, verb, rest)
}

// ActionVerb returns the first word from an action label (after ANSI normalization).
func ActionVerb(text string) string {
	trimmed := NormalizeChoice(text)
	fields := strings.Fields(trimmed)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

// NormalizeChoice strips ANSI colors and trims whitespace for robust comparisons.
func NormalizeChoice(value string) string {
	return strings.TrimSpace(ansiEscapeRE.ReplaceAllString(value, ""))
}

// IsBackChoice checks whether a selected value maps to Back.
func IsBackChoice(value string) bool {
	return strings.EqualFold(NormalizeChoice(value), "Back")
}

// IsExitChoice checks whether a selected value maps to Exit.
func IsExitChoice(value string) bool {
	return strings.EqualFold(NormalizeChoice(value), "Exit")
}
