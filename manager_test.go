package wizard

import "testing"

func TestManageCreate(t *testing.T) {
	sel, err := Manage(ManageOptions[string]{
		Items:       nil,
		ItemLabel:   func(item string) string { return item },
		MainMessage: "Action:",
		ItemMessage: "Item:",
		Select: func(options []string, defaultOption, message string) string {
			if message == "Action:" {
				return "Create new"
			}
			return ""
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sel.Action != ManageCreate {
		t.Fatalf("expected create action, got %q", sel.Action)
	}
	if sel.Item != nil {
		t.Fatal("expected nil item for create action")
	}
}

func TestManageEdit(t *testing.T) {
	sel, err := Manage(ManageOptions[string]{
		Items:         []string{"vmware", "prod"},
		ItemLabel:     func(item string) string { return item },
		MainMessage:   "Action:",
		ItemMessage:   "Item:",
		DefaultLabel:  "prod",
		EditLabel:     "Edit existing",
		CreateLabel:   "Create new",
		DeleteLabel:   "Delete",
		SetDefaultLbl: "Set default",
		CancelLabel:   "Cancel",
		Select: func(options []string, defaultOption, message string) string {
			if message == "Action:" {
				return "Edit existing"
			}
			if message == "Item:" {
				return "prod"
			}
			return "Cancel"
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sel.Action != ManageEdit {
		t.Fatalf("expected edit action, got %q", sel.Action)
	}
	if sel.Item == nil || *sel.Item != "prod" {
		t.Fatalf("expected selected item prod, got %#v", sel.Item)
	}
}

func TestManageSetDefault(t *testing.T) {
	sel, err := Manage(ManageOptions[string]{
		Items:       []string{"a", "b"},
		ItemLabel:   func(item string) string { return item },
		MainMessage: "Action:",
		ItemMessage: "Item:",
		Select: func(options []string, defaultOption, message string) string {
			if message == "Action:" {
				return "Set default"
			}
			return "b"
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sel.Action != ManageSetDefault {
		t.Fatalf("expected set default action, got %q", sel.Action)
	}
	if sel.Item == nil || *sel.Item != "b" {
		t.Fatalf("expected selected item b, got %#v", sel.Item)
	}
}

func TestManageBackFromItemSelection(t *testing.T) {
	sel, err := Manage(ManageOptions[string]{
		Items:       []string{"a", "b"},
		ItemLabel:   func(item string) string { return item },
		MainMessage: "Action:",
		ItemMessage: "Item:",
		Select: func(options []string, defaultOption, message string) string {
			if message == "Action:" {
				return "Edit existing"
			}
			return "Back"
		},
		CancelLabel: "Back",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sel.Action != ManageCancel {
		t.Fatalf("expected cancel action when Back is selected, got %q", sel.Action)
	}
}
