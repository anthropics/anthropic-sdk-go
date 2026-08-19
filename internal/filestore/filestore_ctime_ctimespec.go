//go:build darwin || freebsd || netbsd

package filestore

import (
	"os"
	"syscall"
)

// ctimeNs is fi's inode change time in nanoseconds — see hashEntry.
// These GOOS values spell the syscall.Stat_t field Ctimespec.
func ctimeNs(fi os.FileInfo) int64 {
	return fi.Sys().(*syscall.Stat_t).Ctimespec.Nano()
}
