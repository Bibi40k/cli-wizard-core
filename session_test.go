package wizard

import "testing"

func TestSessionLifecycle(t *testing.T) {
	t.Helper()
	loaded := false
	started := false
	stopped := false
	finalized := false

	s := NewSession(
		"target.yaml",
		"draft.yaml",
		&struct{}{},
		func() bool { return false },
		func(draftPath string, state any) (bool, error) {
			if draftPath != "draft.yaml" {
				t.Fatalf("draftPath = %q", draftPath)
			}
			loaded = true
			return true, nil
		},
		func(targetPath, draftPath string, state any, isEmpty func() bool) func() {
			started = true
			if targetPath != "target.yaml" || draftPath != "draft.yaml" {
				t.Fatalf("unexpected paths: %q %q", targetPath, draftPath)
			}
			return func() { stopped = true }
		},
		func(targetPath string) error {
			if targetPath != "target.yaml" {
				t.Fatalf("targetPath = %q", targetPath)
			}
			finalized = true
			return nil
		},
	)

	ok, err := s.LoadDraft()
	if err != nil || !ok || !loaded {
		t.Fatalf("LoadDraft failed: ok=%v err=%v loaded=%v", ok, err, loaded)
	}
	s.Start()
	s.Stop()
	if err := s.Finalize(); err != nil {
		t.Fatalf("Finalize: %v", err)
	}
	if !started || !stopped || !finalized {
		t.Fatalf("lifecycle flags: started=%v stopped=%v finalized=%v", started, stopped, finalized)
	}
}
