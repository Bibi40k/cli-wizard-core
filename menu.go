package wizard

import "fmt"

// FormatMenuLabel renders a two-column menu label: [tag] + aligned text.
// width controls the fixed width for the [tag] column; values <= 0 default to 12.
func FormatMenuLabel(tag, text string, width int) string {
	if width <= 0 {
		width = 12
	}
	return fmt.Sprintf("%-*s %s", width, "["+tag+"]", text)
}
