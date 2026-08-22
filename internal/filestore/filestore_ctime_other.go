//go:build !linux && !darwin && !freebsd && !netbsd && !openbsd && !dragonfly

package filestore

import "os"

// ctimeNs is 0 here: identity degrades to (mtime, size), so an
// mtime-preserving rewrite (rsync -t, cp -p) is served the stale hash.
func ctimeNs(os.FileInfo) int64 { return 0 }
