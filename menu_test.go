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
		t.Fatalf("expected BackLabel to be recognized")
	}
	if !IsBackChoice(BackMenuLabel(16)) {
		t.Fatalf("expected BackMenuLabel to be recognized")
	}
	if !IsExitChoice("Exit") {
		t.Fatalf("expected plain Exit to be recognized")
	}
	if !IsExitChoice(ExitLabel()) {
		t.Fatalf("expected ExitLabel to be recognized")
	}
	if !IsExitChoice(ExitMenuLabel(16)) {
		t.Fatalf("expected ExitMenuLabel to be recognized")
	}
}

func TestSelectHint(t *testing.T) {
	got := SelectHint()
	want := "[Use arrows to move, type to filter, Esc=Back, Ctrl+C=Exit]"
	if got != want {
		t.Fatalf("unexpected select hint: got %q want %q", got, want)
	}
}

func TestFormatActionLabel(t *testing.T) {
	got := FormatActionLabel("Create new cluster spec", 9)
	want := "Create    new cluster spec"
	if got != want {
		t.Fatalf("unexpected action label: got %q want %q", got, want)
	}
	if FormatActionLabel("Back", 9) != "Back" {
		t.Fatalf("Back should stay plain")
	}
	if FormatActionLabel("Exit", 9) != "Exit" {
		t.Fatalf("Exit should stay plain")
	}
}

func TestActionVerb(t *testing.T) {
	if ActionVerb("Create    new cluster spec") != "Create" {
		t.Fatalf("unexpected verb")
	}
	if ActionVerb(Colorize("Delete   config", "\033[31m")) != "Delete" {
		t.Fatalf("unexpected verb with ansi")
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
