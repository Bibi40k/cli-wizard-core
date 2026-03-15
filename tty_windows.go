//go:build windows
// +build windows

package wizard

// DrainStdin is a no-op on Windows.
func DrainStdin() {}

// RestoreTTY is a no-op on Windows.
func RestoreTTY() {}
