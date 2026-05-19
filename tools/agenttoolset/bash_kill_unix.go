//go:build !windows

package agenttoolset

import (
	"os"
	"syscall"
)

// killProcessGroup terminates p's whole process group (pty.Start sets Setsid,
// so PGID == PID), falling back to p.Kill() if the group kill fails.
func killProcessGroup(p *os.Process) error {
	if p == nil {
		return nil
	}
	if err := syscall.Kill(-p.Pid, syscall.SIGKILL); err == nil {
		return nil
	}
	return p.Kill()
}
