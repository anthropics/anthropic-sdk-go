//go:build !windows

package filestore

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestFileStoreCreatesEveryParentOwnerOnlyUnderPermissiveUmask(t *testing.T) {
	// Every directory the store creates — including intermediate parents,
	// not just the leaf — is 0o700, even under a permissive umask.
	old := syscall.Umask(0)
	defer syscall.Umask(old)

	tmp := t.TempDir()
	root := filepath.Join(tmp, "nested", "store")
	store := createdStore(t, root)
	put(t, store, "projects/foo/notes.md", "hi")
	for _, d := range []string{
		filepath.Join(tmp, "nested"),
		root,
		filepath.Join(root, "projects"),
		filepath.Join(root, "projects", "foo"),
	} {
		st, err := os.Stat(d)
		if err != nil {
			t.Fatalf("stat %s: %v", d, err)
		}
		if st.Mode().Perm() != 0o700 {
			t.Fatalf("%s has mode %v, want 0o700", d, st.Mode().Perm())
		}
	}
}

func TestFileStoreGetRefusesANonRegularFile(t *testing.T) {
	// Get must refuse — not hang on — a FIFO, and refuse a symlink leaf.
	store, root := newStore(t)
	put(t, store, "real.md", "data")
	if err := syscall.Mkfifo(filepath.Join(root, "pipe"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "real.md"), filepath.Join(root, "link.md")); err != nil {
		t.Fatal(err)
	}

	_, err := store.Get("pipe")
	requireRefused(t, err)
	if !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("error %q does not say not a regular file", err)
	}
	_, err = store.Get("link.md")
	requireRefused(t, err)
	if !strings.Contains(err.Error(), "is a symlink") {
		t.Fatalf("error %q does not say is a symlink", err)
	}

	// The open itself must refuse the fifo (O_NONBLOCK + fstat) — a swap
	// after a caller's lstat would otherwise block the read forever.
	_, err = openRegularFile("pipe", filepath.Join(root, "pipe"))
	requireRefused(t, err)
}
