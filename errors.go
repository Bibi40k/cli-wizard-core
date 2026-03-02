package wizard

import (
	"errors"
	"strings"
)

const (
	ANSIReset  = "\033[0m"
	ANSIRed    = "\033[31m"
	ANSIYellow = "\033[33m"
	ANSICyan   = "\033[36m"
)

type hintProvider interface {
	Hint() string
}

// UserError is a user-facing error with an optional actionable hint.
type UserError struct {
	Message string
	HintMsg string
}

func (e *UserError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

func (e *UserError) Hint() string {
	if e == nil {
		return ""
	}
	return e.HintMsg
}

// NewUserError creates a user-facing error message with an optional hint.
func NewUserError(message, hint string) error {
	return &UserError{
		Message: strings.TrimSpace(message),
		HintMsg: strings.TrimSpace(hint),
	}
}

type hintedError struct {
	cause error
	hint  string
}

func (e *hintedError) Error() string {
	if e == nil || e.cause == nil {
		return ""
	}
	return e.cause.Error()
}

func (e *hintedError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

func (e *hintedError) Hint() string {
	if e == nil {
		return ""
	}
	return e.hint
}

// WithHint wraps an existing error with a user-facing hint.
func WithHint(err error, hint string) error {
	if err == nil {
		return nil
	}
	return &hintedError{cause: err, hint: strings.TrimSpace(hint)}
}

// ErrorHint extracts a hint from an error if available.
func ErrorHint(err error) string {
	if err == nil {
		return ""
	}
	var hp hintProvider
	if errors.As(err, &hp) {
		return strings.TrimSpace(hp.Hint())
	}
	return ""
}

// FormatCLIError renders a colored CLI error and optional hint.
func FormatCLIError(err error) string {
	if err == nil {
		return ""
	}
	msg := Colorize("error:", ANSIRed) + " " + err.Error()
	if hint := ErrorHint(err); hint != "" {
		msg += "\n" + Colorize("hint:", ANSIYellow) + " " + Colorize(hint, ANSICyan)
	}
	return msg
}
