//go:build !windows

package filestore

// Tests for FileStore — the confined-folder filesystem layer.

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

func sha(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func mustOpen(t *testing.T, root string, opts ...Option) *FileStore {
	t.Helper()
	store, err := Open(root, opts...)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	return store
}

func createdStore(t *testing.T, root string, opts ...Option) *FileStore {
	t.Helper()
	store := mustOpen(t, root, opts...)
	if err := store.CreateRoot(); err != nil {
		t.Fatalf("CreateRoot: %v", err)
	}
	return store
}

// newStore opens a FileStore on a fresh root, creates the root, and returns
// the store with the root path.
func newStore(t *testing.T) (*FileStore, string) {
	t.Helper()
	root := filepath.Join(t.TempDir(), "store")
	return createdStore(t, root), root
}

func put(t *testing.T, store *FileStore, rel, data string) {
	t.Helper()
	if err := store.Put(rel, []byte(data)); err != nil {
		t.Fatalf("Put(%q): %v", rel, err)
	}
}

func requireRefused(t *testing.T, err error) {
	t.Helper()
	var fse *FileStoreError
	if !errors.As(err, &fse) {
		t.Fatalf("expected *FileStoreError, got %T: %v", err, err)
	}
}

// requireNotCreated fails the test unless err matches fs.ErrNotExist and dir
// is still absent: the write failed instead of creating dir.
func requireNotCreated(t *testing.T, err error, dir string) {
	t.Helper()
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("err = %v; want one matching fs.ErrNotExist", err)
	}
	if _, statErr := os.Lstat(dir); !errors.Is(statErr, fs.ErrNotExist) {
		t.Fatalf("%s exists after the failed write (lstat: %v)", dir, statErr)
	}
}

func TestFileStorePutWritesTextAndCreatesOwnerOnlyParents(t *testing.T) {
	store, root := newStore(t)
	put(t, store, "projects/foo/notes.md", "héllo")

	dest := filepath.Join(root, "projects", "foo", "notes.md")
	got, err := os.ReadFile(dest)
	if err != nil || string(got) != "héllo" {
		t.Fatalf("read back %q, %v", got, err)
	}
	if st, _ := os.Stat(filepath.Dir(dest)); st.Mode().Perm() != 0o700 {
		t.Fatalf("parent mode = %v, want 0o700", st.Mode().Perm())
	}
	if st, _ := os.Stat(dest); st.Mode().Perm() != 0o600 {
		t.Fatalf("file mode = %v, want 0o600", st.Mode().Perm())
	}
}

func TestFileStorePutExecutableSetsOwnerExec(t *testing.T) {
	store, root := newStore(t)
	if err := store.Put("run.sh", []byte("#!/bin/sh\necho hi\n"), Executable()); err != nil {
		t.Fatalf("Put: %v", err)
	}
	if st, _ := os.Stat(filepath.Join(root, "run.sh")); st.Mode().Perm() != 0o700 {
		t.Fatalf("mode = %v, want 0o700", st.Mode().Perm())
	}
}

func TestFileStorePutIsAtomicAndLeavesNoTempOnSuccess(t *testing.T) {
	store, root := newStore(t)
	put(t, store, "a.md", "first")
	put(t, store, "a.md", "second")

	if got, _ := os.ReadFile(filepath.Join(root, "a.md")); string(got) != "second" {
		t.Fatalf("got %q, want second", got)
	}
	if entries, _ := os.ReadDir(root); len(entries) != 1 || entries[0].Name() != "a.md" {
		t.Fatalf("dir entries = %v", entries)
	}
}

func TestFileStoreEveryVerbRefusesAnEscapingPath(t *testing.T) {
	store, root := newStore(t)
	secret := filepath.Join(filepath.Dir(root), "secret.md")
	if err := os.WriteFile(secret, []byte("top secret"), 0o644); err != nil {
		t.Fatal(err)
	}
	put(t, store, "ok.md", "ok")

	for _, bad := range []string{"../secret.md", "a/../../secret.md"} {
		requireRefused(t, store.Put(bad, []byte("nope")))
		_, err := store.Get(bad)
		requireRefused(t, err)
		_, err = store.Ls(bad)
		requireRefused(t, err)
		_, err = store.HashTree(bad)
		requireRefused(t, err)
		_, err = store.HashFile(bad)
		requireRefused(t, err)
		_, err = store.FindSymlinks(bad)
		requireRefused(t, err)
		requireRefused(t, store.Move("ok.md", bad))
		requireRefused(t, store.Move(bad, "landed.md"))
		requireRefused(t, store.Remove(bad))
	}

	if got, _ := os.ReadFile(secret); string(got) != "top secret" {
		t.Fatalf("secret was modified: %q", got)
	}
	// Nothing materialised outside the root either — not even a temp file.
	entries, _ := os.ReadDir(filepath.Dir(root))
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name()
	}
	if !slices.Equal(names, []string{"secret.md", "store"}) {
		t.Fatalf("outside the root: %v", names)
	}
}

func TestFileStoreGetDoesNotFollowASymlinkLeaf(t *testing.T) {
	store, root := newStore(t)
	put(t, store, "real.md", "original")
	if err := os.Symlink(filepath.Join(root, "real.md"), filepath.Join(root, "link.md")); err != nil {
		t.Fatal(err)
	}

	_, err := store.Get("link.md")
	requireRefused(t, err)
	// Put replaces the symlink atomically rather than following it.
	put(t, store, "link.md", "new")
	if got, _ := os.ReadFile(filepath.Join(root, "real.md")); string(got) != "original" {
		t.Fatalf("real.md was overwritten: %q", got)
	}
	if st, _ := os.Lstat(filepath.Join(root, "link.md")); st.Mode()&fs.ModeSymlink != 0 {
		t.Fatal("link.md is still a symlink")
	}
	if got, _ := os.ReadFile(filepath.Join(root, "link.md")); string(got) != "new" {
		t.Fatalf("link.md = %q", got)
	}
}

func TestFileStoreGetReturnsNilWhenAbsent(t *testing.T) {
	store, _ := newStore(t)
	got, err := store.Get("absent.md")
	if got != nil || err != nil {
		t.Fatalf("absent: got %v, %v", got, err)
	}
	put(t, store, "present.md", "here")
	got, err = store.Get("present.md")
	if err != nil || string(got) != "here" {
		t.Fatalf("present: got %q, %v", got, err)
	}
}

func TestFileStoreLsListsRecursivelyAndNeverDescendsSymlinks(t *testing.T) {
	store, root := newStore(t)
	put(t, store, "a.md", "alpha")
	put(t, store, "sub/b.md", "bravo")
	outside := filepath.Join(filepath.Dir(root), "outside")
	if err := os.Mkdir(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outside, "s.md"), []byte("secret"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "lnkdir")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, "s.md"), filepath.Join(root, "lnk.md")); err != nil {
		t.Fatal(err)
	}

	want := map[string]bool{"a.md": true, "sub/b.md": true}
	if got, err := store.Ls("/"); err != nil || !maps.Equal(got, want) {
		t.Fatalf("Ls(/) = %v, %v; want %v", got, err, want)
	}
	if got, err := store.Ls("sub"); err != nil || !maps.Equal(got, map[string]bool{"sub/b.md": true}) {
		t.Fatalf("Ls(sub) = %v, %v", got, err)
	}
	if got, err := store.Ls("absent"); err != nil || len(got) != 0 {
		t.Fatalf("Ls(absent) = %v, %v; want empty map", got, err)
	}
	if got, err := store.Ls("a.md/under"); err != nil || len(got) != 0 {
		t.Fatalf("Ls(a.md/under) = %v, %v; want empty map", got, err)
	}
	_, err := store.Ls("a.md")
	requireRefused(t, err)
}

func TestFileStoreFindSymlinksReportsEverySymlink(t *testing.T) {
	store, root := newStore(t)
	put(t, store, "a.md", "alpha")
	put(t, store, "sub/b.md", "bravo")
	if got, err := store.FindSymlinks("/"); err != nil || len(got) != 0 {
		t.Fatalf("FindSymlinks on clean store = %v, %v; want empty map", got, err)
	}

	outside := filepath.Join(filepath.Dir(root), "outside")
	if err := os.Mkdir(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outside, "s.md"), []byte("secret"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "lnkdir")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(outside, "s.md"), filepath.Join(root, "sub", "lnk.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "gone.md"), filepath.Join(root, "dangling.md")); err != nil {
		t.Fatal(err)
	}

	want := map[string]bool{"lnkdir": true, "sub/lnk.md": true, "dangling.md": true}
	if got, err := store.FindSymlinks("/"); err != nil || !maps.Equal(got, want) {
		t.Fatalf("FindSymlinks(/) = %v, %v; want %v", got, err, want)
	}
	if got, err := store.FindSymlinks("sub"); err != nil || !maps.Equal(got, map[string]bool{"sub/lnk.md": true}) {
		t.Fatalf("FindSymlinks(sub) = %v, %v", got, err)
	}
	if got, err := store.FindSymlinks("absent"); err != nil || len(got) != 0 {
		t.Fatalf("FindSymlinks(absent) = %v, %v; want empty map", got, err)
	}
	// The base itself being a symlink is reported, not descended.
	if got, err := store.FindSymlinks("lnkdir"); err != nil || !maps.Equal(got, map[string]bool{"lnkdir": true}) {
		t.Fatalf("FindSymlinks(lnkdir) = %v, %v", got, err)
	}
	// A regular-file base is refused like any non-directory.
	_, err := store.FindSymlinks("a.md")
	requireRefused(t, err)
}

func TestFileStoreHashTreeHashesRegularFilesAndSkipsSymlinks(t *testing.T) {
	store, root := newStore(t)
	put(t, store, "a.md", "alpha")
	put(t, store, "sub/b.md", "bravo")
	outside := filepath.Join(filepath.Dir(root), "secret.txt")
	_ = os.WriteFile(outside, []byte("secret"), 0o644)
	_ = os.Symlink(outside, filepath.Join(root, "link.md"))

	want := map[string]string{"a.md": sha("alpha"), "sub/b.md": sha("bravo")}
	if got, err := store.HashTree("/"); err != nil || !maps.Equal(got, want) {
		t.Fatalf("HashTree(/) = %v, %v; want %v", got, err, want)
	}
	wantSub := map[string]string{"sub/b.md": sha("bravo")}
	if got, err := store.HashTree("sub/"); err != nil || !maps.Equal(got, wantSub) {
		t.Fatalf("HashTree(sub/) = %v, %v; want %v", got, err, wantSub)
	}
}

func TestFileStoreListingsErrorOnAnUnreadableSubdirectory(t *testing.T) {
	// An unreadable subdirectory fails the listing rather than dropping out
	// of it — a caller diffing listings would otherwise read its files as gone.
	if os.Geteuid() == 0 {
		t.Skip("permission checks do not bind root")
	}
	store, root := newStore(t)
	put(t, store, "a.md", "alpha")
	put(t, store, "locked/b.md", "bravo")
	locked := filepath.Join(root, "locked")
	if err := os.Chmod(locked, 0o000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(locked, 0o700)

	if _, err := store.HashTree("/"); !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("HashTree = %v; want a permission error", err)
	}
	if _, err := store.Ls("/"); !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("Ls = %v; want a permission error", err)
	}
	if _, err := store.FindSymlinks("/"); !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("FindSymlinks = %v; want a permission error", err)
	}

	if err := os.Chmod(locked, 0o700); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"a.md": sha("alpha"), "locked/b.md": sha("bravo")}
	if got, err := store.HashTree("/"); err != nil || !maps.Equal(got, want) {
		t.Fatalf("HashTree(/) after chmod = %v, %v; want %v", got, err, want)
	}
}

func TestFileStoreHashTreeRefusesASymlinkedPrefix(t *testing.T) {
	// A prefix that exists but is a symlink is refused rather than listed
	// (or silently empty); a plain file prefix is refused too.
	store, root := newStore(t)
	put(t, store, "real/a.md", "alpha")
	if err := os.Symlink(filepath.Join(root, "real"), filepath.Join(root, "alias")); err != nil {
		t.Fatal(err)
	}

	_, err := store.HashTree("alias")
	requireRefused(t, err)
	_, err = store.HashTree("real/a.md")
	requireRefused(t, err)
	if got, err := store.HashTree("absent"); err != nil || len(got) != 0 {
		t.Fatalf("HashTree(absent) = %v, %v; want empty map", got, err)
	}
}

func TestFileStoreHashFileHashesOneFileAndSharesTheCache(t *testing.T) {
	store, root := newStore(t)
	realMargin := timestampTrustMarginNs
	timestampTrustMarginNs = -1_000_000_000_000_000
	t.Cleanup(func() { timestampTrustMarginNs = realMargin })
	hashed := countHashes(t)

	put(t, store, "a.md", "alpha")
	if got, err := store.HashFile("a.md"); err != nil || got != sha("alpha") {
		t.Fatalf("HashFile(a.md) = %q, %v", got, err)
	}
	if got, err := store.HashFile("absent.md"); err != nil || got != "" {
		t.Fatalf("HashFile(absent.md) = %q, %v; want empty", got, err)
	}
	// HashTree reuses the entry HashFile just cached.
	want := map[string]string{"a.md": sha("alpha")}
	if got, err := store.HashTree("/"); err != nil || !maps.Equal(got, want) {
		t.Fatalf("HashTree = %v, %v; want %v", got, err, want)
	}
	if !slices.Equal(*hashed, []string{"a.md"}) {
		t.Fatalf("hashed %v, want one hash shared across both verbs", *hashed)
	}

	if err := os.Symlink(filepath.Join(root, "a.md"), filepath.Join(root, "lnk.md")); err != nil {
		t.Fatal(err)
	}
	_, err := store.HashFile("lnk.md")
	requireRefused(t, err)
	var fse *FileStoreError
	if !errors.As(err, &fse) || fse.Reason != ReasonIsASymlink {
		t.Fatalf("HashFile(lnk.md) = %v; want an is-a-symlink refusal", err)
	}
}

func TestFileStoreMoveRenamesWithinRootAndRefusesExistingDest(t *testing.T) {
	store, _ := newStore(t)
	put(t, store, "a.md", "alpha")
	if err := store.Move("a.md", "nested/b.md"); err != nil {
		t.Fatalf("Move: %v", err)
	}
	if got, _ := store.Get("nested/b.md"); string(got) != "alpha" {
		t.Fatalf("Get nested/b.md = %q", got)
	}
	if got, _ := store.Get("a.md"); got != nil {
		t.Fatal("a.md still present")
	}
	put(t, store, "c.md", "charlie")
	requireRefused(t, store.Move("nested/b.md", "c.md"))
}

func TestFileStoreRootPathIsBannedFromTheInterface(t *testing.T) {
	// A rel resolving to the root: Put/Get refuse, Move/Remove do nothing.
	store, _ := newStore(t)
	put(t, store, "x.md", "x")

	if err := store.Remove("/"); err != nil {
		t.Fatalf("Remove(/): %v", err)
	}
	if got, err := store.Get("x.md"); err != nil || string(got) != "x" {
		t.Fatalf("x.md after Remove(/) = %q, %v", got, err)
	}
	if err := store.Move("/", "elsewhere"); err != nil {
		t.Fatalf("Move(/, elsewhere): %v", err)
	}
	if err := store.Move("x.md", "/"); err != nil {
		t.Fatalf("Move(x.md, /): %v", err)
	}
	if got, err := store.Ls("/"); err != nil || !maps.Equal(got, map[string]bool{"x.md": true}) {
		t.Fatalf("Ls = %v, %v; want only x.md", got, err)
	}

	requireRefused(t, store.Put("/", []byte("data")))
	_, err := store.Get("/")
	requireRefused(t, err)
	var fse *FileStoreError
	if !errors.As(err, &fse) || fse.Reason != ReasonNotAFile {
		t.Fatalf("Get(/) = %v; want a not-a-regular-file refusal", err)
	}
}

func TestFileStoreRemoveDeletesFileAndSubtree(t *testing.T) {
	store, root := newStore(t)
	put(t, store, "a.md", "a")
	put(t, store, "sub/b.md", "b")

	if err := store.Remove("a.md"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "a.md")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatal("a.md still exists")
	}
	if err := store.Remove("sub"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "sub")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatal("sub/ still exists")
	}
	if err := store.Remove("absent.md"); err != nil {
		t.Fatalf("Remove(absent.md): %v", err)
	}
}

func TestFileStoreRemoveUnlinksADanglingSymlink(t *testing.T) {
	store, root := newStore(t)
	if err := os.Symlink(filepath.Join(root, "nowhere"), filepath.Join(root, "dangle")); err != nil {
		t.Fatal(err)
	}

	if err := store.Remove("dangle"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(root, "dangle")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("dangling symlink still present: %v", err)
	}
}

func TestFileStoreOpenErrorsOnAnUnreadableMountPath(t *testing.T) {
	// A mount path the process cannot stat errors at Open — it is never read
	// as absent, which would mark a real directory ours to remove on Dispose.
	if os.Geteuid() == 0 {
		t.Skip("permission checks do not bind root")
	}
	locked := filepath.Join(t.TempDir(), "locked")
	if err := os.MkdirAll(filepath.Join(locked, "store"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(locked, 0o000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(locked, 0o700)

	_, err := Open(filepath.Join(locked, "store"))
	if !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("Open = %v; want a permission error", err)
	}
	var fse *FileStoreError
	if errors.As(err, &fse) {
		t.Fatalf("Open = %v; a deployment failure must not be a *FileStoreError", err)
	}
}

func TestFileStoreDisposeRemovesRootOnlyIfCreated(t *testing.T) {
	store, fresh := newStore(t)
	put(t, store, "x.md", "x")
	if r := store.Root(); !r.RemovedOnDispose || r.Path != fresh {
		t.Fatalf("Root = %+v", r)
	}
	if err := store.Dispose(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(fresh); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("root %s not removed", fresh)
	}

	old := filepath.Join(t.TempDir(), "old")
	if err := os.Mkdir(old, 0o755); err != nil {
		t.Fatal(err)
	}
	store = mustOpen(t, old)
	put(t, store, "y.md", "y")
	if r := store.Root(); r.RemovedOnDispose {
		t.Fatal("Root marked a pre-existing root for removal")
	}
	if err := store.Dispose(); err != nil {
		t.Fatal(err)
	}
	if got, err := os.ReadFile(filepath.Join(old, "y.md")); err != nil || string(got) != "y" {
		t.Fatalf("y.md = %q, %v", got, err)
	}
}

func TestFileStorePutWritesArchiveBytesVerbatim(t *testing.T) {
	// Archive-shaped bytes are just bytes — never extracted.
	store, _ := newStore(t)
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, _ := zw.Create("pdf/SKILL.md")
	_, _ = w.Write([]byte("# pdf"))
	_ = zw.Close()
	blob := buf.Bytes()

	if err := store.Put("bundle.zip", blob); err != nil {
		t.Fatalf("Put: %v", err)
	}
	got, err := store.Get("bundle.zip")
	if err != nil || !bytes.Equal(got, blob) {
		t.Fatalf("bundle.zip is not byte-identical to the input: %v", err)
	}
	want := map[string]string{"bundle.zip": sha(string(blob))}
	if tree, _ := store.HashTree("/"); !maps.Equal(tree, want) {
		t.Fatalf("HashTree = %v, want %v", tree, want)
	}
}

func TestFileStorePutRefusesDirectoryShapedPaths(t *testing.T) {
	store, _ := newStore(t)
	for _, dir := range []string{"bundles/v1/", "dir/.", ".", "/", ""} {
		requireRefused(t, store.Put(dir, []byte("data")))
	}
	if got, err := store.HashTree("/"); err != nil || len(got) != 0 {
		t.Fatalf("HashTree = %v, %v; want nothing written", got, err)
	}
}

func TestFileStoreUTF8RestrictionRefusesBinaryContent(t *testing.T) {
	root := filepath.Join(t.TempDir(), "store")
	store := createdStore(t, root, WithUTF8())
	binary := []byte{0xff, 0xfe, 0x00, ' ', 'n', 'o', 't', ' ', 'u', 't', 'f', '8'}

	// A Put of invalid bytes is refused and writes nothing.
	requireRefused(t, store.Put("blob.bin", binary))
	if _, err := os.Stat(filepath.Join(root, "blob.bin")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatal("blob.bin materialised despite the refusal")
	}

	// Conforming content is unaffected.
	put(t, store, "note.md", "héllo")
	if got, err := store.Get("note.md"); err != nil || string(got) != "héllo" {
		t.Fatalf("Get note.md = %q, %v", got, err)
	}

	// A file smuggled in behind the store's back (written directly to disk)
	// is refused at Get.
	if err := os.WriteFile(filepath.Join(root, "smuggled.bin"), binary, 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := store.Get("smuggled.bin")
	requireRefused(t, err)

	// An unrestricted store keeps returning raw bytes.
	loose := createdStore(t, filepath.Join(t.TempDir(), "store2"))
	if err := loose.Put("blob.bin", binary); err != nil {
		t.Fatalf("Put on unrestricted store: %v", err)
	}
	if got, err := loose.Get("blob.bin"); err != nil || !bytes.Equal(got, binary) {
		t.Fatalf("Get blob.bin = %q, %v", got, err)
	}
}

func TestFileStoreIsPathLegal(t *testing.T) {
	for path, want := range map[string]bool{
		"/mnt/memory/notes":     true,
		"mnt/relative":          false,
		"/mnt/../../etc/cron.d": false,
		"":                      false,
	} {
		if got := IsPathLegal(path); got != want {
			t.Errorf("IsPathLegal(%q) = %v, want %v", path, got, want)
		}
	}
}

// countHashes replaces hashFile with a counting wrapper for the duration of
// the test and returns the slice of hashed base names.
func countHashes(t *testing.T) *[]string {
	t.Helper()
	var hashed []string
	realHash := hashFile
	hashFile = func(path string) (string, error) {
		hashed = append(hashed, filepath.Base(path))
		return realHash(path)
	}
	t.Cleanup(func() { hashFile = realHash })
	return &hashed
}

func TestFileStoreHashTreeRehashesOnlyWhatChanged(t *testing.T) {
	// A file whose stat identity (mtime, ctime, size) is unchanged since the
	// last walk is served from the cache; a changed file is re-read. The
	// trust margin is disabled here so fresh files cache immediately; a
	// writer that restores mtime and size (rsync -t style) is still
	// re-hashed, because userspace cannot restore ctime.
	store, root := newStore(t)
	realMargin := timestampTrustMarginNs
	timestampTrustMarginNs = -1_000_000_000_000_000
	t.Cleanup(func() { timestampTrustMarginNs = realMargin })
	hashed := countHashes(t)

	put(t, store, "a.md", "alpha")
	put(t, store, "b.md", "bravo")
	want := map[string]string{"a.md": sha("alpha"), "b.md": sha("bravo")}
	if got, err := store.HashTree("/"); err != nil || !maps.Equal(got, want) {
		t.Fatalf("HashTree = %v, %v; want %v", got, err, want)
	}
	slices.Sort(*hashed)
	if !slices.Equal(*hashed, []string{"a.md", "b.md"}) {
		t.Fatalf("first walk hashed %v, want both files", *hashed)
	}

	*hashed = nil
	if got, err := store.HashTree("/"); err != nil || !maps.Equal(got, want) {
		t.Fatalf("HashTree = %v, %v; want %v", got, err, want)
	}
	if len(*hashed) != 0 {
		t.Fatalf("unchanged walk hashed %v, want nothing", *hashed)
	}

	// An mtime-and-size-preserving rewrite: same length, mtime restored. The
	// sleep guarantees the rewrite's ctime lands on a later coarse-clock tick
	// than the original's — ctime is the one stamp the rewrite cannot restore.
	time.Sleep(50 * time.Millisecond)
	b := filepath.Join(root, "b.md")
	st, err := os.Stat(b)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(b, []byte("BRAVO"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(b, st.ModTime(), st.ModTime()); err != nil {
		t.Fatal(err)
	}
	got, err := store.HashTree("/")
	if err != nil || got["b.md"] != sha("BRAVO") {
		t.Fatalf("HashTree after rewrite = %v, %v; want b.md re-read", got, err)
	}
	if !slices.Equal(*hashed, []string{"b.md"}) {
		t.Fatalf("rewrite walk hashed %v, want only b.md", *hashed)
	}
}

func TestFileStoreHashTreeDistrustsFreshStamps(t *testing.T) {
	// Filesystems stamp times with coarse clocks, so a same-tick rewrite can
	// reuse the exact stat identity of the version just hashed. A file whose
	// stamps are within the trust margin of the walk is never cached —
	// pinned by freezing the walk clock at the file's own stamp time.
	store, root := newStore(t)
	hashed := countHashes(t)

	put(t, store, "racy.md", "AAAAAA")
	racy := filepath.Join(root, "racy.md")
	st, err := os.Lstat(racy)
	if err != nil {
		t.Fatal(err)
	}
	stamp := ctimeNs(st)
	realNow := nowNs
	nowNs = func() int64 { return stamp }
	t.Cleanup(func() { nowNs = realNow })

	if got, err := store.HashTree("/"); err != nil || got["racy.md"] != sha("AAAAAA") {
		t.Fatalf("HashTree = %v, %v", got, err)
	}
	if got, err := store.HashTree("/"); err != nil || got["racy.md"] != sha("AAAAAA") {
		t.Fatalf("HashTree = %v, %v", got, err)
	}
	// Both walks hashed: nothing this fresh is ever cached.
	if !slices.Equal(*hashed, []string{"racy.md", "racy.md"}) {
		t.Fatalf("hashed = %v, want both walks to hash", *hashed)
	}

	// What the margin defends: a same-tick rewrite that preserves the whole
	// stat identity must be seen by the next walk.
	if err := os.WriteFile(racy, []byte("BBBBBB"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got, err := store.HashTree("/"); err != nil || got["racy.md"] != sha("BBBBBB") {
		t.Fatalf("HashTree after rewrite = %v, %v; want the new bytes", got, err)
	}

	// The margin's width is part of the contract: stamps 1.5s old are still
	// distrusted (coarse filesystems stamp in whole seconds); 10s-old stamps
	// are safely past it.
	*hashed = nil
	st, err = os.Lstat(racy)
	if err != nil {
		t.Fatal(err)
	}
	now := ctimeNs(st)
	nowNs = func() int64 { return now + 1_500_000_000 }
	if _, err := store.HashTree("/"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.HashTree("/"); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(*hashed, []string{"racy.md", "racy.md"}) {
		t.Fatalf("hashed inside the margin = %v, want both walks to hash", *hashed)
	}
	nowNs = func() int64 { return now + 10_000_000_000 }
	if _, err := store.HashTree("/"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.HashTree("/"); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(*hashed, []string{"racy.md", "racy.md", "racy.md"}) {
		t.Fatalf("hashed past the margin = %v, want exactly one more hash", *hashed)
	}
}

func TestCreateRootMakesTheFolderOwnerOnly(t *testing.T) {
	root := filepath.Join(t.TempDir(), "store")
	store := mustOpen(t, root)
	if err := store.CreateRoot(); err != nil {
		t.Fatalf("CreateRoot: %v", err)
	}
	st, err := os.Stat(root)
	if err != nil {
		t.Fatal(err)
	}
	if got := st.Mode().Perm(); got != 0o700 {
		t.Fatalf("root mode = %v, want 0700", got)
	}
	if err := store.CreateRoot(); err != nil { // already existing is fine
		t.Fatalf("second CreateRoot: %v", err)
	}
	if !store.Root().RemovedOnDispose {
		t.Fatal("RemovedOnDispose = false, want true: Open found no root, so this store created it")
	}
}

func TestFileStorePutNeverCreatesAMissingRoot(t *testing.T) {
	// Only CreateRoot makes the folder: a Put into a store that was never
	// created, or whose folder was removed, fails and re-creates nothing —
	// not the root and nothing above it.
	mnt := filepath.Join(t.TempDir(), "mnt")
	root := filepath.Join(mnt, "store")
	store := mustOpen(t, root)

	requireNotCreated(t, store.Put("a.md", []byte("a")), mnt)
	requireNotCreated(t, store.Put("sub/a.md", []byte("a")), mnt)

	if err := store.CreateRoot(); err != nil {
		t.Fatalf("CreateRoot: %v", err)
	}
	put(t, store, "a.md", "a")
	if err := os.RemoveAll(root); err != nil {
		t.Fatal(err)
	}
	requireNotCreated(t, store.Put("b.md", []byte("b")), root)
	requireNotCreated(t, store.Put("sub/b.md", []byte("b")), root)

	if err := os.RemoveAll(mnt); err != nil {
		t.Fatal(err)
	}
	requireNotCreated(t, store.Put("z.md", []byte("z")), mnt)

	if err := store.CreateRoot(); err != nil {
		t.Fatalf("CreateRoot: %v", err)
	}
	put(t, store, "z.md", "z")
	if got, err := store.Get("z.md"); err != nil || string(got) != "z" {
		t.Fatalf("Get z.md = %q, %v", got, err)
	}
}

func TestFileStoreMoveNeverCreatesAMissingRoot(t *testing.T) {
	store, root := newStore(t)
	put(t, store, "a.md", "a")
	if err := os.RemoveAll(root); err != nil {
		t.Fatal(err)
	}

	requireNotCreated(t, store.Move("a.md", "sub/b.md"), root)
}
