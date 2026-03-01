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
