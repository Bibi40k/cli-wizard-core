package wizard

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

// MultiSelect displays an interactive checkbox list.
// Arrow keys navigate, Space toggles selection, Enter confirms.
// Ctrl+C/ESC returns defaults and sets interrupted (Ctrl+C only).
// Items are returned in original (not selection) order.
// On non-TTY stdin, returns defaults without prompting.
func (s *Selector) MultiSelect(items []string, defaults []string, message string) []string {
	s.interrupted = false

	if len(items) == 0 {
		return defaults
	}

	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return defaults
	}

	// Build initial selection state from defaults.
	selected := make([]bool, len(items))
	defaultSet := make(map[string]bool, len(defaults))
	for _, d := range defaults {
		defaultSet[d] = true
	}
	for i, it := range items {
		if defaultSet[it] {
			selected[i] = true
		}
	}

	state, err := term.MakeRaw(fd)
	if err != nil {
		return defaults
	}
	defer func() { _ = term.Restore(fd, state) }()

	cursor := 0
	query := ""
	maxVis := s.MaxVisible
	offset := 0
	renderedLines := 0

	clampView := func(count int) {
		if count == 0 {
			cursor = 0
			offset = 0
			return
		}
		if cursor >= count {
			cursor = count - 1
		}
		if cursor < 0 {
			cursor = 0
		}
		vis := count
		if maxVis > 0 && vis > maxVis {
			vis = maxVis
		}
		if cursor < offset {
			offset = cursor
		} else if cursor >= offset+vis {
			offset = cursor - vis + 1
		}
	}

	render := func() {
		if renderedLines > 0 {
			_, _ = fmt.Fprintf(os.Stdout, "\r\033[%dA\033[J", renderedLines)
		} else {
			_, _ = fmt.Fprint(os.Stdout, "\r\033[J")
		}

		filtered := filterItems(items, query)
		clampView(len(filtered))

		vis := len(filtered)
		if maxVis > 0 && vis > maxVis {
			vis = maxVis
		}

		// Header
		_, _ = fmt.Fprintf(os.Stdout, "\033[32m?\033[0m \033[1;37m%s\033[0m  \033[36m%s\033[0m",
			strings.TrimSpace(message), MultiSelectHint())
		if query != "" {
			_, _ = fmt.Fprintf(os.Stdout, "  (filter: %s)", query)
		}
		_, _ = fmt.Fprint(os.Stdout, "\r\n")
		lines := 1

		end := offset + vis
		if end > len(filtered) {
			end = len(filtered)
		}
		for i := offset; i < end; i++ {
			idx := filtered[i]
			item := items[idx]
			check := "[ ]"
			checkColor := ""
			if selected[idx] {
				check = "[x]"
				checkColor = "\033[32m"
			}
			if i == cursor {
				_, _ = fmt.Fprintf(os.Stdout, "  \033[36m❯ %s%s\033[0m %s\r\n", checkColor, check, item)
			} else {
				_, _ = fmt.Fprintf(os.Stdout, "    %s%s\033[0m %s\r\n", checkColor, check, item)
			}
			lines++
		}

		if len(filtered) == 0 {
			_, _ = fmt.Fprint(os.Stdout, "    (no matches)\r\n")
			lines++
		}

		// Footer
		if maxVis > 0 && len(filtered) > maxVis {
			_, _ = fmt.Fprintf(os.Stdout, "  \033[2m%d/%d\033[0m\r\n", cursor+1, len(filtered))
		} else {
			_, _ = fmt.Fprint(os.Stdout, "\r\n")
		}
		lines++

		renderedLines = lines
	}

	collectResult := func() []string {
		var result []string
		for i, it := range items {
			if selected[i] {
				result = append(result, it)
			}
		}
		return result
	}

	showResult := func(result []string) {
		if renderedLines > 0 {
			_, _ = fmt.Fprintf(os.Stdout, "\r\033[%dA\033[J", renderedLines)
		} else {
			_, _ = fmt.Fprint(os.Stdout, "\r\033[J")
		}
		summary := strings.Join(result, ", ")
		if summary == "" {
			summary = "(none)"
		}
		_, _ = fmt.Fprintf(os.Stdout, "\033[32m?\033[0m \033[1;37m%s\033[0m \033[36m%s\033[0m\r\n",
			strings.TrimSpace(message), summary)
	}

	render()

	buf := make([]byte, 8)
	for {
		n, rerr := os.Stdin.Read(buf)
		if rerr != nil || n == 0 {
			s.interrupted = true
			showResult(defaults)
			return defaults
		}

		switch {
		case n >= 3 && buf[0] == 27 && buf[1] == '[':
			filtered := filterItems(items, query)
			switch buf[2] {
			case 'A': // Up
				if len(filtered) > 0 {
					cursor--
					if cursor < 0 {
						cursor = len(filtered) - 1
					}
				}
				render()
			case 'B': // Down
				if len(filtered) > 0 {
					cursor++
					if cursor >= len(filtered) {
						cursor = 0
					}
				}
				render()
			}

		case n == 1 && buf[0] == 3: // Ctrl+C
			s.interrupted = true
			showResult(defaults)
			return defaults

		case n == 1 && buf[0] == 27: // ESC
			showResult(defaults)
			return defaults

		case n == 1 && buf[0] == ' ': // Space — toggle selection
			filtered := filterItems(items, query)
			if len(filtered) > 0 {
				idx := filtered[cursor]
				selected[idx] = !selected[idx]
			}
			render()

		case n == 1 && (buf[0] == '\r' || buf[0] == '\n'): // Enter — confirm
			result := collectResult()
			showResult(result)
			return result

		case n == 1 && (buf[0] == 127 || buf[0] == 8): // Backspace
			if len(query) > 0 {
				_, size := utf8.DecodeLastRuneInString(query)
				if size > 0 && size <= len(query) {
					query = query[:len(query)-size]
				}
				cursor = 0
				render()
			}

		case n == 1 && buf[0] >= 33 && buf[0] <= 126: // Printable (excl. space, handled above)
			query += string(buf[0])
			cursor = 0
			render()
		}
	}
}

// MultiSelectHint returns the standard keyboard hint for multi-select prompts.
func MultiSelectHint() string {
	return "[arrows move, space toggle, enter confirm, Esc=Back, Ctrl+C=Exit]"
}
