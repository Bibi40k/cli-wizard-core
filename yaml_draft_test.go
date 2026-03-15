package wizard

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestDraftTargetToken(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"configs/prod/node.yaml", "configs__prod__node.yaml"},
		{"./configs/node.yaml", "configs__node.yaml"},
		{"node.yaml", "node.yaml"},
		{"a/b/c.yaml", "a__b__c.yaml"},
	}
	for _, tc := range tests {
		got := DraftTargetToken(tc.input)
		if got != tc.want {
			t.Errorf("DraftTargetToken(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestLoadDraftYAMLEmptyPath(t *testing.T) {
	var out any
	ok, err := LoadDraftYAML("", &out)
	if ok || err != nil {
		t.Errorf("LoadDraftYAML(\"\") = (%v, %v), want (false, nil)", ok, err)
	}
}

func TestLoadDraftYAMLMissingFile(t *testing.T) {
	var out any
	ok, err := LoadDraftYAML("/nonexistent/path.yaml", &out)
	if ok {
		t.Error("LoadDraftYAML on missing file should return ok=false")
	}
	if err == nil {
		t.Error("LoadDraftYAML on missing file should return an error")
	}
}

func TestWriteAndLoadDraftYAML(t *testing.T) {
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	type myStruct struct {
		Name  string `yaml:"name"`
		Value int    `yaml:"value"`
	}

	original := myStruct{Name: "test", Value: 42}
	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	draftPath, err := WriteDraft("configs/prod/node.yaml", data)
	if err != nil {
		t.Fatalf("WriteDraft: %v", err)
	}
	if draftPath == "" {
		t.Fatal("WriteDraft returned empty path")
	}
	if _, err := os.Stat(draftPath); err != nil {
		t.Fatalf("draft file not created: %v", err)
	}

	var loaded myStruct
	ok, err := LoadDraftYAML(draftPath, &loaded)
	if !ok || err != nil {
		t.Fatalf("LoadDraftYAML = (%v, %v), want (true, nil)", ok, err)
	}
	if loaded.Name != original.Name || loaded.Value != original.Value {
		t.Errorf("loaded = %+v, want %+v", loaded, original)
	}
}

func TestLatestDraftForTarget(t *testing.T) {
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	// No drafts yet.
	if got := LatestDraftForTarget("configs/node.yaml"); got != "" {
		t.Errorf("LatestDraftForTarget with no drafts = %q, want \"\"", got)
	}

	// Create two drafts with different modification times.
	if err := os.MkdirAll("tmp", 0700); err != nil {
		t.Fatalf("mkdir tmp: %v", err)
	}
	token := DraftTargetToken("configs/node.yaml")
	older := filepath.Join("tmp", token+".draft.20240101-100000.yaml")
	newer := filepath.Join("tmp", token+".draft.20240102-100000.yaml")
	_ = os.WriteFile(older, []byte("a: 1"), 0600)
	_ = os.WriteFile(newer, []byte("a: 2"), 0600)
	// Set explicit mtime so sort is deterministic.
	base := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	_ = os.Chtimes(older, base, base)
	_ = os.Chtimes(newer, base.Add(time.Hour), base.Add(time.Hour))

	got := LatestDraftForTarget("configs/node.yaml")
	if got != newer {
		t.Errorf("LatestDraftForTarget = %q, want %q", got, newer)
	}
}

func TestCleanupDrafts(t *testing.T) {
	tmp := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmp); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	defer func() { _ = os.Chdir(origDir) }()

	if err := os.MkdirAll("tmp", 0700); err != nil {
		t.Fatalf("mkdir tmp: %v", err)
	}

	token := DraftTargetToken("configs/node.yaml")
	draft1 := filepath.Join("tmp", token+".draft.20240101-100000.yaml")
	draft2 := filepath.Join("tmp", token+".draft.20240101-120000.yaml")
	other := filepath.Join("tmp", "other.draft.20240101-100000.yaml")
	_ = os.WriteFile(draft1, []byte("x"), 0600)
	_ = os.WriteFile(draft2, []byte("x"), 0600)
	_ = os.WriteFile(other, []byte("x"), 0600)

	if err := CleanupDrafts("configs/node.yaml"); err != nil {
		t.Fatalf("CleanupDrafts: %v", err)
	}

	if _, err := os.Stat(draft1); !os.IsNotExist(err) {
		t.Error("draft1 should have been removed")
	}
	if _, err := os.Stat(draft2); !os.IsNotExist(err) {
		t.Error("draft2 should have been removed")
	}
	if _, err := os.Stat(other); err != nil {
		t.Error("other draft should NOT have been removed")
	}
}
