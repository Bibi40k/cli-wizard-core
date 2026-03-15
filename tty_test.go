package wizard

import (
	"testing"
)

func TestDrainStdinNoPanic(t *testing.T) {
	// DrainStdin should not panic even when stdin is not a TTY.
	DrainStdin()
}

func TestRestoreTTYNoPanic(t *testing.T) {
	// RestoreTTY should not panic even when stdin is not a TTY.
	RestoreTTY()
}
