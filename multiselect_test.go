package wizard

import (
	"reflect"
	"testing"
)

func TestMultiSelectNonTTY(t *testing.T) {
	// In test environment stdin is a pipe (non-TTY), so MultiSelect returns defaults.
	s := NewSelector()
	items := []string{"alpha", "beta", "gamma"}
	defaults := []string{"alpha", "gamma"}

	got := s.MultiSelect(items, defaults, "Pick items:")
	if !reflect.DeepEqual(got, defaults) {
		t.Errorf("MultiSelect on non-TTY = %v, want %v", got, defaults)
	}
	if s.WasInterrupted() {
		t.Error("WasInterrupted() should be false after non-TTY MultiSelect")
	}
}

func TestMultiSelectNonTTYEmptyDefaults(t *testing.T) {
	s := NewSelector()
	got := s.MultiSelect([]string{"a", "b"}, nil, "Pick:")
	if got != nil {
		t.Errorf("MultiSelect with nil defaults on non-TTY = %v, want nil", got)
	}
}

func TestMultiSelectEmptyItems(t *testing.T) {
	s := NewSelector()
	defaults := []string{"x"}
	got := s.MultiSelect(nil, defaults, "Pick:")
	if !reflect.DeepEqual(got, defaults) {
		t.Errorf("MultiSelect(nil items) = %v, want %v", got, defaults)
	}
}

func TestMultiSelectHint(t *testing.T) {
	h := MultiSelectHint()
	if h == "" {
		t.Error("MultiSelectHint() should not be empty")
	}
}
