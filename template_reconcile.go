package wizard

import (
	"fmt"
	"sort"
	"strings"
)

// ReconcileOptions controls template sync behavior.
type ReconcileOptions struct {
	// DropUnknown removes keys not present in template.
	DropUnknown bool
	// RequiredPaths are dot-notation keys that must exist after reconcile.
	RequiredPaths []string
	// CheckPlaceholders, when true, checks required paths for empty/CHANGE_ME values.
	CheckPlaceholders bool
}

// ReconcileReport summarizes changes and validation findings.
type ReconcileReport struct {
	Added             []string
	Removed           []string
	MissingRequired   []string
	PlaceholderValues []string
}

// IsPlaceholderValue returns true if v is empty, nil, or contains "CHANGE_ME".
func IsPlaceholderValue(v any) bool {
	s, ok := v.(string)
	if !ok {
		return v == nil
	}
	s = strings.TrimSpace(s)
	return s == "" || strings.Contains(strings.ToUpper(s), "CHANGE_ME")
}

// ReconcileWithTemplate syncs a config object against a template object.
// It keeps existing values for keys defined in template, fills missing keys
// from template defaults, and optionally removes unknown keys.
func ReconcileWithTemplate(current, template map[string]any, opts ReconcileOptions) (map[string]any, ReconcileReport, error) {
	if template == nil {
		return nil, ReconcileReport{}, fmt.Errorf("template is nil")
	}
	if current == nil {
		current = map[string]any{}
	}
	out := cloneMap(current)
	report := ReconcileReport{
		Added:           make([]string, 0),
		Removed:         make([]string, 0),
		MissingRequired: make([]string, 0),
	}
	reconcileObject("", out, template, opts.DropUnknown, &report)
	report.MissingRequired = append(report.MissingRequired, findMissingRequired(out, opts.RequiredPaths)...)
	if opts.CheckPlaceholders {
		report.PlaceholderValues = append(report.PlaceholderValues, findPlaceholderValues(out, opts.RequiredPaths, report.MissingRequired)...)
	}
	sort.Strings(report.Added)
	sort.Strings(report.Removed)
	sort.Strings(report.MissingRequired)
	sort.Strings(report.PlaceholderValues)
	return out, report, nil
}

func reconcileObject(prefix string, out, tpl map[string]any, dropUnknown bool, report *ReconcileReport) {
	for k, tv := range tpl {
		path := joinPath(prefix, k)
		ov, ok := out[k]
		if !ok {
			out[k] = cloneAny(tv)
			report.Added = append(report.Added, path)
			continue
		}
		tm, tok := tv.(map[string]any)
		om, ook := ov.(map[string]any)
		if tok && ook {
			reconcileObject(path, om, tm, dropUnknown, report)
			out[k] = om
		}
	}
	if !dropUnknown {
		return
	}
	for k := range out {
		if _, ok := tpl[k]; !ok {
			path := joinPath(prefix, k)
			delete(out, k)
			report.Removed = append(report.Removed, path)
		}
	}
}

func findMissingRequired(data map[string]any, required []string) []string {
	missing := make([]string, 0)
	for _, p := range required {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if !hasPath(data, p) {
			missing = append(missing, p)
		}
	}
	return missing
}

// findPlaceholderValues returns required paths whose values are placeholders,
// excluding paths already in missing (which don't exist at all).
func findPlaceholderValues(data map[string]any, required []string, missing []string) []string {
	missingSet := make(map[string]struct{}, len(missing))
	for _, p := range missing {
		missingSet[p] = struct{}{}
	}
	result := make([]string, 0)
	for _, p := range required {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if _, isMissing := missingSet[p]; isMissing {
			continue
		}
		v := getPath(data, p)
		if IsPlaceholderValue(v) {
			result = append(result, p)
		}
	}
	return result
}

func getPath(data map[string]any, path string) any {
	parts := strings.Split(path, ".")
	var cur any = data
	for _, part := range parts {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = m[part]
	}
	return cur
}

func hasPath(data map[string]any, path string) bool {
	parts := strings.Split(path, ".")
	var cur any = data
	for _, part := range parts {
		m, ok := cur.(map[string]any)
		if !ok {
			return false
		}
		v, exists := m[part]
		if !exists {
			return false
		}
		cur = v
	}
	return true
}

func joinPath(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

func cloneMap(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = cloneAny(v)
	}
	return out
}

func cloneAny(v any) any {
	switch t := v.(type) {
	case map[string]any:
		return cloneMap(t)
	case []any:
		out := make([]any, 0, len(t))
		for _, it := range t {
			out = append(out, cloneAny(it))
		}
		return out
	default:
		return t
	}
}
