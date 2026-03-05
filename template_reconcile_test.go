package wizard

import "testing"

func TestReconcileWithTemplate_AddMissingAndKeepExisting(t *testing.T) {
	current := map[string]any{
		"vm": map[string]any{
			"name": "node-01",
		},
	}
	template := map[string]any{
		"vm": map[string]any{
			"name":       "",
			"profile":    "talos",
			"ip_address": "",
		},
	}

	out, report, err := ReconcileWithTemplate(current, template, ReconcileOptions{})
	if err != nil {
		t.Fatalf("ReconcileWithTemplate error: %v", err)
	}
	vm, ok := out["vm"].(map[string]any)
	if !ok {
		t.Fatalf("vm should be object")
	}
	if got := vm["name"]; got != "node-01" {
		t.Fatalf("expected existing name to stay, got %v", got)
	}
	if got := vm["profile"]; got != "talos" {
		t.Fatalf("expected missing profile from template, got %v", got)
	}
	if len(report.Added) != 2 {
		t.Fatalf("expected 2 added paths, got %d: %v", len(report.Added), report.Added)
	}
}

func TestReconcileWithTemplate_DropUnknown(t *testing.T) {
	current := map[string]any{
		"vm": map[string]any{
			"name":     "node-01",
			"username": "",
		},
	}
	template := map[string]any{
		"vm": map[string]any{
			"name": "node-01",
		},
	}
	out, report, err := ReconcileWithTemplate(current, template, ReconcileOptions{DropUnknown: true})
	if err != nil {
		t.Fatalf("ReconcileWithTemplate error: %v", err)
	}
	vm := out["vm"].(map[string]any)
	if _, ok := vm["username"]; ok {
		t.Fatalf("expected username to be removed")
	}
	if len(report.Removed) != 1 || report.Removed[0] != "vm.username" {
		t.Fatalf("unexpected removed paths: %v", report.Removed)
	}
}

func TestReconcileWithTemplate_RequiredPaths(t *testing.T) {
	current := map[string]any{
		"vm": map[string]any{
			"name": "node-01",
		},
	}
	template := map[string]any{
		"vm": map[string]any{
			"name": "node-01",
		},
	}
	_, report, err := ReconcileWithTemplate(current, template, ReconcileOptions{
		RequiredPaths: []string{"vm.name", "vm.ip_address"},
	})
	if err != nil {
		t.Fatalf("ReconcileWithTemplate error: %v", err)
	}
	if len(report.MissingRequired) != 1 || report.MissingRequired[0] != "vm.ip_address" {
		t.Fatalf("unexpected missing required: %v", report.MissingRequired)
	}
}
