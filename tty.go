//go:build !windows
// +build !windows

package wizard

import (
	"os"
	"os/exec"
	"syscall"
	"time"
)

// DrainStdin discards bytes pending in stdin (e.g. cursor-position responses
// left by interactive rendering in raw mode). Safe to call from non-TTY contexts.
func DrainStdin() {
	fd := int(os.Stdin.Fd())
	if err := syscall.SetNonblock(fd, true); err != nil {
		return
	}
	defer func() {
		_ = syscall.SetNonblock(fd, false)
	}()

	buf := make([]byte, 256)
	deadline := time.Now().Add(80 * time.Millisecond)
	for time.Now().Before(deadline) {
		n, err := syscall.Read(fd, buf)
		if n > 0 {
			deadline = time.Now().Add(80 * time.Millisecond)
			continue
		}
		if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		break
	}
}

// RestoreTTY best-effort restores terminal state. Call from signal handlers
// or os.Exit paths that bypass deferred term.Restore calls.
func RestoreTTY() {
	fd := int(os.Stdin.Fd())
	_ = syscall.SetNonblock(fd, false)

	cmd := exec.Command("stty", "sane")
	cmd.Stdin = os.Stdin
	_ = cmd.Run()
}
