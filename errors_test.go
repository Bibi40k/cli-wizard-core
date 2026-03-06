package wizard

import (
	"errors"
	"strings"
	"testing"
)

func TestNewUserError(t *testing.T) {
	err := NewUserError("invalid input", "provide a value")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got := err.Error(); got != "invalid input" {
		t.Fatalf("unexpected error message: %q", got)
	}
	if got := ErrorHint(err); got != "provide a value" {
		t.Fatalf("unexpected hint: %q", got)
	}
}

func TestWithHint(t *testing.T) {
	base := errors.New("boom")
	err := WithHint(base, "retry with --dry-run")
	if err == nil {
		t.Fatal("expected wrapped error, got nil")
	}
	if got := err.Error(); got != "boom" {
		t.Fatalf("unexpected wrapped error message: %q", got)
	}
	if got := ErrorHint(err); got != "retry with --dry-run" {
		t.Fatalf("unexpected wrapped hint: %q", got)
	}
}

func TestErrorHint_NoHint(t *testing.T) {
	if got := ErrorHint(errors.New("plain")); got != "" {
		t.Fatalf("expected empty hint, got %q", got)
	}
}

func TestFormatCLIError_WithHint(t *testing.T) {
	out := FormatCLIError(NewUserError("doctor requires --spec", "use --spec examples/spec.sample.json"))
	if !strings.Contains(out, "error:") {
		t.Fatalf("expected error label, got %q", out)
	}
	if !strings.Contains(out, "hint:") {
		t.Fatalf("expected hint label, got %q", out)
	}
	if !strings.Contains(out, "doctor requires --spec") {
		t.Fatalf("expected error content, got %q", out)
	}
}

func TestFormatCLIError_WithoutHint(t *testing.T) {
	out := FormatCLIError(errors.New("plain failure"))
	if !strings.Contains(out, "plain failure") {
		t.Fatalf("expected error content, got %q", out)
	}
	if strings.Contains(out, "hint:") {
		t.Fatalf("did not expect hint label, got %q", out)
	}
}

func TestIsInterrupted(t *testing.T) {
	if !IsInterrupted(ErrInterrupted) {
		t.Fatal("expected sentinel to be recognized")
	}
	wrapped := WithHint(ErrInterrupted, "stop")
	if !IsInterrupted(wrapped) {
		t.Fatal("expected wrapped interrupted sentinel to be recognized")
	}
	if IsInterrupted(errors.New("other")) {
		t.Fatal("unexpected interrupted match")
	}
}
