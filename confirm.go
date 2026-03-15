package wizard

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// Confirm displays an interactive Y/N prompt.
// Returns true for yes, false for no.
// On Ctrl+C, returns false and sets the interrupted flag (check WasInterrupted()).
// On ESC, returns false (not interrupted).
// On non-TTY stdin, returns defaultYes without prompting.
func (s *Selector) Confirm(message string, defaultYes bool) bool {
	s.interrupted = false

	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		return defaultYes
	}

	hint := "[y/N]"
	if defaultYes {
		hint = "[Y/n]"
	}

	state, err := term.MakeRaw(fd)
	if err != nil {
		return defaultYes
	}
	defer func() { _ = term.Restore(fd, state) }()

	_, _ = fmt.Fprintf(os.Stdout, "\033[32m?\033[0m \033[1;37m%s\033[0m %s: ", strings.TrimSpace(message), hint)

	buf := make([]byte, 4)
	for {
		n, rerr := os.Stdin.Read(buf)
		if rerr != nil || n == 0 {
			s.interrupted = true
			_, _ = fmt.Fprint(os.Stdout, "\r\n")
			return false
		}

		switch {
		case n == 1 && (buf[0] == 'y' || buf[0] == 'Y'):
			_, _ = fmt.Fprintf(os.Stdout, "\033[36mYes\033[0m\r\n")
			return true

		case n == 1 && (buf[0] == 'n' || buf[0] == 'N'):
			_, _ = fmt.Fprintf(os.Stdout, "\033[33mNo\033[0m\r\n")
			return false

		case n == 1 && (buf[0] == '\r' || buf[0] == '\n'):
			if defaultYes {
				_, _ = fmt.Fprintf(os.Stdout, "\033[36mYes\033[0m\r\n")
				return true
			}
			_, _ = fmt.Fprintf(os.Stdout, "\033[33mNo\033[0m\r\n")
			return false

		case n == 1 && buf[0] == 3: // Ctrl+C
			s.interrupted = true
			_, _ = fmt.Fprintf(os.Stdout, "\033[33mNo\033[0m\r\n")
			return false

		case n == 1 && buf[0] == 27: // ESC
			_, _ = fmt.Fprintf(os.Stdout, "\033[33mNo\033[0m\r\n")
			return false
		}
		// Ignore other bytes (e.g. escape sequences).
	}
}
