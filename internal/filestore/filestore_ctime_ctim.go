//go:build linux || openbsd || dragonfly

package filestore

import (
	"os"
	"syscall"
)

// ctimeNs is fi's inode change time in nanoseconds — see hashEntry.
// These GOOS values spell the syscall.Stat_t field Ctim.
func ctimeNs(fi os.FileInfo) int64 {
	return fi.Sys().(*syscall.Stat_t).Ctim.Nano()
}
