package wizard

import (
	"testing"
)

func TestConfirmNonTTY(t *testing.T) {
	// In test environment stdin is a pipe (non-TTY), so Confirm returns defaultYes.
	s := NewSelector()

	if got := s.Confirm("Continue?", true); got != true {
		t.Errorf("Confirm(defaultYes=true) on non-TTY = %v, want true", got)
	}
	if s.WasInterrupted() {
		t.Error("WasInterrupted() should be false after non-TTY Confirm")
	}

	if got := s.Confirm("Continue?", false); got != false {
		t.Errorf("Confirm(defaultYes=false) on non-TTY = %v, want false", got)
	}
	if s.WasInterrupted() {
		t.Error("WasInterrupted() should be false after non-TTY Confirm")
	}
}

func TestConfirmInitialState(t *testing.T) {
	s := NewSelector()
	if s.WasInterrupted() {
		t.Error("fresh Selector should not be interrupted")
	}
}
