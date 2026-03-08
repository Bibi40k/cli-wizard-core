package wizard

import (
	"testing"
)

func TestFilterItems(t *testing.T) {
	items := []string{"Apple", "Banana", "Cherry", "Avocado"}

	tests := []struct {
		query string
		want  []int
	}{
		{"", []int{0, 1, 2, 3}},
		{"a", []int{0, 1, 3}},
		{"an", []int{1}},
		{"z", []int{}},
		{"CHERRY", []int{2}},
	}
	for _, tc := range tests {
		got := filterItems(items, tc.query)
		if len(got) != len(tc.want) {
			t.Errorf("filterItems(%q) = %v, want %v", tc.query, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("filterItems(%q)[%d] = %d, want %d", tc.query, i, got[i], tc.want[i])
			}
		}
	}
}

func TestPreferredCancel(t *testing.T) {
	tests := []struct {
		items []string
		want  string
	}{
		{[]string{"Create", "Edit", "Back"}, "Back"},
		{[]string{"Create", "Exit"}, "Exit"},
		{[]string{"Create", "Cancel"}, "Cancel"},
		{[]string{"Create", "Exit", "Back"}, "Back"},
		{[]string{"Create", "Edit"}, ""},
	}
	for _, tc := range tests {
		got := preferredCancel(tc.items)
		if got != tc.want {
			t.Errorf("preferredCancel(%v) = %q, want %q", tc.items, got, tc.want)
		}
	}
}

func TestDefaultItemIndex(t *testing.T) {
	items := []string{"A", "B", "C"}

	tests := []struct {
		def  string
		want int
	}{
		{"A", 0},
		{"B", 1},
		{"C", 2},
		{"D", 0},
		{"", 0},
	}
	for _, tc := range tests {
		got := defaultItemIndex(items, tc.def)
		if got != tc.want {
			t.Errorf("defaultItemIndex(%v, %q) = %d, want %d", items, tc.def, got, tc.want)
		}
	}
}

func TestIsCancelChoice(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		{"Cancel", true},
		{"cancel", true},
		{"CANCEL", true},
		{"  Cancel  ", true},
		{"Back", false},
		{"Exit", false},
		{"Something", false},
	}
	for _, tc := range tests {
		if got := IsCancelChoice(tc.value); got != tc.want {
			t.Errorf("IsCancelChoice(%q) = %v, want %v", tc.value, got, tc.want)
		}
	}
}

func TestNewSelector(t *testing.T) {
	s := NewSelector()
	if s.MaxVisible != 10 {
		t.Errorf("NewSelector().MaxVisible = %d, want 10", s.MaxVisible)
	}
	if s.WasInterrupted() {
		t.Error("NewSelector().WasInterrupted() should be false")
	}
}

func TestSelectEmptyItems(t *testing.T) {
	s := NewSelector()
	got := s.Select(nil, "fallback", "Pick:")
	if got != "fallback" {
		t.Errorf("Select(nil) = %q, want %q", got, "fallback")
	}
}
