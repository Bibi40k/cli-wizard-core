package wizard

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

// Selector provides an interactive arrow-key selection widget for terminal UIs.
type Selector struct {
	// MaxVisible limits items visible at once. 0 means show all.
	MaxVisible  int
	interrupted bool
}

// NewSelector creates a Selector with sensible defaults (MaxVisible = 10).
func NewSelector() *Selector {
	return &Selector{MaxVisible: 10}
}

// Select displays an interactive menu with arrow-key navigation, type-to-filter,
// and optional scrolling. Returns the selected item.
// On ESC or Ctrl+C, returns the preferred cancel option (Back > Exit > Cancel)
// if present, or empty string. Call WasInterrupted() to check for Ctrl+C.
//
// Signature is compatible with ManageOptions.Select.
func (s *Selector) Select(items []string, defaultItem, message string) string {
	s.interrupted = false
	if len(items) == 0 {
		return defaultItem
	}
	if defaultItem == "" {
		defaultItem = items[0]
	}

	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return defaultItem
	}

	state, err := term.MakeRaw(fd)
	if err != nil {
		return defaultItem
	}
	defer func() { _ = term.Restore(fd, state) }()

	query := ""
	cursor := defaultItemIndex(items, defaultItem)
	cancel := preferredCancel(items)
	maxVis := s.MaxVisible
	offset := 0
	renderedLines := 0
	isDestructive := strings.HasPrefix(strings.ToLower(strings.TrimSpace(message)), "delete ")

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
			fmt.Fprintf(os.Stdout, "\r\033[%dA\033[J", renderedLines)
		} else {
			fmt.Fprint(os.Stdout, "\r\033[J")
		}

		filtered := filterItems(items, query)
		clampView(len(filtered))

		vis := len(filtered)
		if maxVis > 0 && vis > maxVis {
			vis = maxVis
		}

		// Header
		hdrColor := "\033[1;37m"
		if isDestructive {
			hdrColor = "\033[1;31m"
		}
		fmt.Fprintf(os.Stdout, "\033[32m?\033[0m %s%s\033[0m  \033[36m%s\033[0m",
			hdrColor, strings.TrimSpace(message), SelectHint())
		if query != "" {
			fmt.Fprintf(os.Stdout, "  (filter: %s)", query)
		}
		fmt.Fprint(os.Stdout, "\r\n")
		lines := 1

		// Items
		end := offset + vis
		if end > len(filtered) {
			end = len(filtered)
		}
		for i := offset; i < end; i++ {
			idx := filtered[i]
			item := items[idx]
			clean := NormalizeChoice(item)
			if IsBackChoice(clean) || IsExitChoice(clean) || IsCancelChoice(clean) {
				item = clean
			}

			if i == cursor {
				color := "\033[36m"
				if IsBackChoice(clean) || IsExitChoice(clean) || IsCancelChoice(clean) {
					color = "\033[33m"
				} else if isDestructive || strings.HasPrefix(strings.ToLower(clean), "delete ") {
					color = "\033[31m"
				}
				fmt.Fprintf(os.Stdout, "  %s❯ %s\033[0m\r\n", color, item)
			} else {
				fmt.Fprintf(os.Stdout, "    %s\r\n", item)
			}
			lines++
		}

		if len(filtered) == 0 {
			fmt.Fprint(os.Stdout, "    (no matches)\r\n")
			lines++
		}

		// Footer
		if maxVis > 0 && len(filtered) > maxVis {
			fmt.Fprintf(os.Stdout, "  \033[2m%d/%d\033[0m\r\n", cursor+1, len(filtered))
		} else {
			fmt.Fprint(os.Stdout, "\r\n")
		}
		lines++

		renderedLines = lines
	}

	showResult := func(result string) {
		if renderedLines > 0 {
			fmt.Fprintf(os.Stdout, "\r\033[%dA\033[J", renderedLines)
		} else {
			fmt.Fprint(os.Stdout, "\r\033[J")
		}
		if result != "" {
			rc := "\033[36m"
			clean := NormalizeChoice(result)
			if strings.HasPrefix(strings.ToLower(clean), "delete ") {
				rc = "\033[31m"
			} else if IsBackChoice(clean) || IsExitChoice(clean) || IsCancelChoice(clean) {
				rc = "\033[33m"
			}
			fmt.Fprintf(os.Stdout, "\033[32m?\033[0m %s %s%s\033[0m\r\n",
				strings.TrimSpace(message), rc, result)
		}
	}

	render()

	buf := make([]byte, 8)
	for {
		n, rerr := os.Stdin.Read(buf)
		if rerr != nil || n == 0 {
			s.interrupted = true
			showResult(cancel)
			return cancel
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
			showResult(cancel)
			if cancel != "" {
				return cancel
			}
			return ""

		case n == 1 && (buf[0] == '\r' || buf[0] == '\n'): // Enter
			filtered := filterItems(items, query)
			if len(filtered) == 0 {
				continue
			}
			result := items[filtered[cursor]]
			showResult(result)
			return result

		case n == 1 && buf[0] == 27: // Bare ESC
			showResult(cancel)
			if cancel != "" {
				return cancel
			}
			return defaultItem

		case n == 1 && (buf[0] == 127 || buf[0] == 8): // Backspace
			if len(query) > 0 {
				_, size := utf8.DecodeLastRuneInString(query)
				if size > 0 && size <= len(query) {
					query = query[:len(query)-size]
				}
				cursor = 0
				render()
			}

		case n == 1 && buf[0] >= 32 && buf[0] <= 126: // Printable
			query += string(buf[0])
			cursor = 0
			render()
		}
	}
}

// WasInterrupted returns true if the last Select ended with Ctrl+C.
func (s *Selector) WasInterrupted() bool {
	return s.interrupted
}

// filterItems returns indices of items matching the query (case-insensitive substring).
func filterItems(items []string, query string) []int {
	q := strings.TrimSpace(strings.ToLower(query))
	out := make([]int, 0, len(items))
	for i, it := range items {
		if q == "" || strings.Contains(strings.ToLower(NormalizeChoice(it)), q) {
			out = append(out, i)
		}
	}
	return out
}

// preferredCancel finds the best cancel option: Back > Exit > Cancel.
func preferredCancel(items []string) string {
	for _, it := range items {
		if IsBackChoice(it) {
			return it
		}
	}
	for _, it := range items {
		if IsExitChoice(it) {
			return it
		}
	}
	for _, it := range items {
		if IsCancelChoice(it) {
			return it
		}
	}
	return ""
}

// defaultItemIndex returns the index of defaultItem in items, or 0.
func defaultItemIndex(items []string, defaultItem string) int {
	for i, it := range items {
		if it == defaultItem {
			return i
		}
	}
	return 0
}
