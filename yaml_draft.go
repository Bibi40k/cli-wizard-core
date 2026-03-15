package wizard

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// DraftTargetToken converts a target file path to a safe filename token.
// Uses path→"__" conversion to preserve directory context.
// Example: "configs/prod/node.yaml" → "configs__prod__node.yaml"
func DraftTargetToken(targetPath string) string {
	p := filepath.Clean(strings.TrimSpace(targetPath))
	p = strings.TrimPrefix(p, "./")
	p = strings.ReplaceAll(p, "\\", "/")
	return strings.ReplaceAll(p, "/", "__")
}

// WriteDraft writes plaintext YAML to tmp/<token>.draft.<timestamp>.yaml.
// Creates the tmp/ directory if needed. Returns the draft file path.
func WriteDraft(targetPath string, plaintext []byte) (string, error) {
	if err := os.MkdirAll("tmp", 0700); err != nil {
		return "", err
	}
	base := DraftTargetToken(targetPath)
	ts := time.Now().Format("20060102-150405")
	draftPath := filepath.Join("tmp", fmt.Sprintf("%s.draft.%s.yaml", base, ts))
	if err := os.WriteFile(draftPath, plaintext, 0600); err != nil {
		return "", err
	}
	return draftPath, nil
}

// LoadDraftYAML reads a draft file and unmarshals it into out.
// Returns (true, nil) on success, (false, nil) if draftPath is empty,
// and (false, err) on read or parse failure.
func LoadDraftYAML(draftPath string, out any) (bool, error) {
	draftPath = strings.TrimSpace(draftPath)
	if draftPath == "" {
		return false, nil
	}
	data, err := os.ReadFile(draftPath)
	if err != nil {
		return false, err
	}
	if err := yaml.Unmarshal(data, out); err != nil {
		return false, fmt.Errorf("parse draft: %w", err)
	}
	return true, nil
}

// CleanupDrafts removes all draft files for the given target path.
func CleanupDrafts(targetPath string) error {
	base := DraftTargetToken(targetPath)
	pattern := filepath.Join("tmp", fmt.Sprintf("%s.draft.*.yaml", base))
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, p := range matches {
		if rmErr := os.Remove(p); rmErr != nil && !os.IsNotExist(rmErr) {
			return rmErr
		}
	}
	return nil
}

// LatestDraftForTarget returns the path of the most-recently-modified draft
// for the given target, or empty string if none exist.
func LatestDraftForTarget(targetPath string) string {
	base := DraftTargetToken(targetPath)
	pattern := filepath.Join("tmp", fmt.Sprintf("%s.draft.*.yaml", base))
	matches, _ := filepath.Glob(pattern)
	if len(matches) == 0 {
		return ""
	}
	type info struct {
		path string
		mod  int64
	}
	var files []info
	for _, p := range matches {
		st, err := os.Stat(p)
		if err != nil {
			continue
		}
		files = append(files, info{path: p, mod: st.ModTime().UnixNano()})
	}
	if len(files) == 0 {
		return ""
	}
	sort.Slice(files, func(i, j int) bool { return files[i].mod > files[j].mod })
	return files[0].path
}
