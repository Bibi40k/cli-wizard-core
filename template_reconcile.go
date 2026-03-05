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
}

// ReconcileReport summarizes changes and validation findings.
type ReconcileReport struct {
	Added           []string
	Removed         []string
	MissingRequired []string
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
	sort.Strings(report.Added)
	sort.Strings(report.Removed)
	sort.Strings(report.MissingRequired)
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
