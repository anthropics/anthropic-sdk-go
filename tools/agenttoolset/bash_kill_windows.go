//go:build windows

package agenttoolset

import "os"

// killProcessGroup falls back to a single-process kill on Windows.
func killProcessGroup(p *os.Process) error {
	if p == nil {
		return nil
	}
	return p.Kill()
}
