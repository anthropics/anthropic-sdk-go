//go:build !windows

package filestore

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"syscall"
)

// platformCheck returns an error when the platform cannot uphold the store's
// guarantees; unix has O_NOFOLLOW, so it never fails here.
func platformCheck() error { return nil }

// replaceViaTemp writes data atomically: a random ".fs-<hex>.tmp" sibling is
// created with O_EXCL|O_NOFOLLOW, then renamed onto dest — a reader never
// sees a half-written file, and a symlink at the leaf is replaced, not
// followed. Best-effort temp cleanup on failure; never masks the original
// error.
func replaceViaTemp(dest string, data []byte, isExecutable bool) error {
	mode := ownerOnlyFileMode
	if isExecutable {
		mode = ownerOnlyExecMode
	}
	var rnd [8]byte
	_, _ = rand.Read(rnd[:])
	tmp := filepath.Join(filepath.Dir(dest), ".fs-"+hex.EncodeToString(rnd[:])+".tmp")
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_EXCL|syscall.O_NOFOLLOW, mode)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, dest); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

// openRegularFile opens dest for reading without following a symlink leaf;
// only a regular file comes back. O_NONBLOCK so a FIFO fails the stat check
// below instead of blocking the open forever.
func openRegularFile(rel, dest string) (*os.File, error) {
	f, err := os.OpenFile(dest, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		// FreeBSD reports EMLINK rather than ELOOP for O_NOFOLLOW.
		if errors.Is(err, syscall.ELOOP) || errors.Is(err, syscall.EMLINK) {
			return nil, &FileStoreError{Reason: ReasonIsASymlink, RelPath: rel}
		}
		return nil, err
	}
	st, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	if !st.Mode().IsRegular() {
		_ = f.Close()
		return nil, &FileStoreError{Reason: ReasonNotAFile, RelPath: rel}
	}
	return f, nil
}
