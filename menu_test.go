package wizard

import "testing"

func TestFormatMenuLabel_DefaultWidth(t *testing.T) {
	got := FormatMenuLabel("vm", "Edit vm.node.sops.yaml", 0)
	want := "[vm]         Edit vm.node.sops.yaml"
	if got != want {
		t.Fatalf("unexpected label: got %q want %q", got, want)
	}
}

func TestFormatMenuLabel_CustomWidth(t *testing.T) {
	got := FormatMenuLabel("schematic", "Manage talos.schematics.sops.yaml", 14)
	want := "[schematic]    Manage talos.schematics.sops.yaml"
	if got != want {
		t.Fatalf("unexpected label: got %q want %q", got, want)
	}
}

func TestFormatMenuLabel_EmptyTag(t *testing.T) {
	got := FormatMenuLabel("", "Exit", 14)
	want := "               Exit"
	if got != want {
		t.Fatalf("unexpected label: got %q want %q", got, want)
	}
}

func TestBackHelpers(t *testing.T) {
	if !IsBackChoice("Back") {
		t.Fatalf("expected plain Back to be recognized")
	}
	if !IsBackChoice(BackLabel()) {
		t.Fatalf("expected colored BackLabel to be recognized")
	}
	if !IsBackChoice(BackMenuLabel(16)) {
		t.Fatalf("expected colored BackMenuLabel to be recognized")
	}
}

func TestColorize_EmptyColor(t *testing.T) {
	got := Colorize("abc", "")
	if got != "abc" {
		t.Fatalf("unexpected colorized text: %q", got)
	}
}

func TestColorize_WithColor(t *testing.T) {
	got := Colorize("abc", "\033[31m")
	want := "\033[31mabc\033[0m"
	if got != want {
		t.Fatalf("unexpected colorized text: got %q want %q", got, want)
	}
}
