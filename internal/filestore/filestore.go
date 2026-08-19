// Package filestore confines all reads and writes to one folder — a relative
// path cannot escape it.
//
// Beta scope: symlinks are refused or skipped wherever the store meets them,
// but there is no hardening against a process racing the store's own
// syscalls; fsync durability, non-unix platforms (Open refuses where
// O_NOFOLLOW is missing), and read-size caps are out of scope.
package filestore

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode/utf8"
)

// Folders the store creates hold downloaded user/model content — keep them
// owner-only rather than inheriting the process umask.
const (
	ownerOnlyDirMode  os.FileMode = 0o700
	ownerOnlyFileMode os.FileMode = 0o600
	ownerOnlyExecMode os.FileMode = 0o700
)

// Refusal reasons carried by [FileStoreError].
const (
	ReasonEscapesRoot           = "escapes the store root"
	ReasonIsASymlink            = "is a symlink"
	ReasonNotAFile              = "is not a regular file"
	ReasonNotADirectory         = "is not a directory"
	ReasonNotUTF8               = "is not valid utf-8"
	ReasonMoveDestinationExists = "already exists"
)

// FileStoreError marks a refused operation — input the store will not act on.
// OS-level errors propagate as their usual *PathError values.
type FileStoreError struct {
	Reason  string
	RelPath string
}

func (e *FileStoreError) Error() string { return fmt.Sprintf("path %q %s", e.RelPath, e.Reason) }

// Root is a store's resolved root path and what [FileStore.Dispose] will do to
// it. RemovedOnDispose is true when Open found no root and this store created
// it, so Dispose removes it. A root that was already there is someone else's —
// a pre-seeded mount, a caller's workdir — and is kept.
type Root struct {
	Path             string
	RemovedOnDispose bool
}

// IsPathLegal reports whether path is usable verbatim as a store location:
// absolute (in POSIX terms — mount paths arrive in wire form) with no ".."
// elements.
func IsPathLegal(path string) bool {
	return strings.HasPrefix(path, "/") && !slices.Contains(strings.Split(path, "/"), "..")
}

// Option configures a store at [Open].
type Option func(*config)

type config struct {
	utf8 bool
}

// WithUTF8 restricts the store to valid UTF-8: a Put of binary bytes and a
// Get of a binary file are refused with a *FileStoreError, so a caller that
// decodes what Get returns can never hit a decode error.
func WithUTF8() Option {
	return func(c *config) { c.utf8 = true }
}

// PutOption configures a single Put call.
type PutOption func(*putConfig)

type putConfig struct {
	executable bool
}

// Executable marks the written file executable (owner-only exec bit).
func Executable() PutOption { return func(c *putConfig) { c.executable = true } }

// hashed is one hashed file version: the stat identity it had, and its sha.
type hashed struct {
	mtimeNs int64
	// userspace cannot set ctime, so writers that preserve mtimes
	// (rsync -t, cp -p) still miss the cache.
	ctimeNs int64
	size    int64
	sha     string
}

// FileStore is one confined folder of regular files.
//
// Every rel path is relative to the root (a leading "/" also means the root)
// and refused with a *FileStoreError when it escapes. The store holds regular
// files only: symlinks are refused on read and skipped by listings —
// [FileStore.FindSymlinks] reports them. A rel path resolving to the root
// itself is banned by this interface: Put and Get refuse it, Move and Remove
// do nothing. A store opened with [WithUTF8] refuses binary content the same
// way — on a Put of such bytes and on a Get of such a file. Only
// [FileStore.CreateRoot] makes the root: writes create directories below it,
// never the root itself, so a root removed while the store is open stays
// removed and the write fails with an error matching [fs.ErrNotExist].
type FileStore struct {
	root             string
	removedOnDispose bool
	utf8Only         bool
	// hashMu guards hashes: HashTree and HashFile share the cache.
	hashMu sync.Mutex
	hashes map[string]hashed
}

// Open resolves root; creates nothing — only [FileStore.CreateRoot] makes the
// folder. On platforms without O_NOFOLLOW it returns an error (a deployment
// condition, not a *FileStoreError input refusal).
func Open(root string, opts ...Option) (*FileStore, error) {
	if err := platformCheck(); err != nil {
		return nil, err
	}
	var cfg config
	for _, o := range opts {
		o(&cfg)
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	_, err = os.Lstat(abs)
	// Only "no entry" means we created it. Any other failure (EACCES, ELOOP)
	// must not mark a real directory ours to delete on dispose.
	removedOnDispose := errors.Is(err, fs.ErrNotExist)
	if err != nil && !removedOnDispose {
		return nil, err
	}
	return &FileStore{
		root:             abs,
		removedOnDispose: removedOnDispose,
		utf8Only:         cfg.utf8,
		hashes:           map[string]hashed{},
	}, nil
}

// CreateRoot creates the root directory and any missing ancestors; already
// existing is fine.
func (s *FileStore) CreateRoot() error {
	return makeDirAndAncestors(s.root)
}

// Root returns the resolved root and what Dispose will do to it.
func (s *FileStore) Root() Root {
	return Root{Path: s.root, RemovedOnDispose: s.removedOnDispose}
}

// Dispose removes the root iff this Open created it; pre-existing roots are kept.
func (s *FileStore) Dispose() error {
	if !s.removedOnDispose {
		return nil
	}
	return os.RemoveAll(s.root)
}

// Put writes data atomically to the file at rel. Missing directories below
// the root are created; a missing root is not — the write fails with an
// error matching [fs.ErrNotExist].
func (s *FileStore) Put(rel string, data []byte, opts ...PutOption) error {
	var cfg putConfig
	for _, o := range opts {
		o(&cfg)
	}
	// "dir/." names a directory just like a trailing "/".
	tail := strings.ReplaceAll(rel, "\\", "/")
	if tail == "" || tail == "." || strings.HasSuffix(tail, "/") || strings.HasSuffix(tail, "/.") {
		return &FileStoreError{Reason: ReasonNotAFile, RelPath: rel}
	}
	dest, err := s.resolveUnderRoot(rel)
	if err != nil {
		return err
	}
	if err := s.requireUTF8(rel, data); err != nil {
		return err
	}
	if err := makeDirsBelowRoot(s.root, filepath.Dir(dest)); err != nil {
		return err
	}
	return replaceViaTemp(dest, data, cfg.executable)
}

// Get returns the file's bytes; (nil, nil) when absent.
func (s *FileStore) Get(rel string) ([]byte, error) {
	dest, err := s.resolveUnderRoot(rel)
	if err != nil {
		return nil, err
	}
	f, err := openRegularFile(rel, dest)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if err := s.requireUTF8(rel, data); err != nil {
		return nil, err
	}
	return data, nil
}

// Ls returns the relative path of every file under the directory under.
func (s *FileStore) Ls(under string) (map[string]bool, error) {
	base, err := s.resolveUnderRoot(under)
	if err != nil {
		return nil, err
	}
	files, err := filenamesInDir(s.root, under, base)
	if err != nil {
		return nil, err
	}
	out := make(map[string]bool)
	for _, f := range files {
		out[f.rel] = true
	}
	return out, nil
}

// FindSymlinks returns every symlink under under — listings skip them and
// reads refuse them.
func (s *FileStore) FindSymlinks(under string) (map[string]bool, error) {
	base, err := s.resolveUnderRoot(under)
	if err != nil {
		return nil, err
	}
	return symlinksInDir(s.root, under, base)
}

// HashTree returns {rel path → sha256-hex} of every file under the directory
// under. Unchanged files — same size, mtime, and ctime since the last call —
// reuse their recorded hash instead of being re-read.
func (s *FileStore) HashTree(under string) (map[string]string, error) {
	base, err := s.resolveUnderRoot(under)
	if err != nil {
		return nil, err
	}
	walkStartNs := nowNs()
	files, err := filenamesInDir(s.root, under, base)
	if err != nil {
		return nil, err
	}
	s.hashMu.Lock()
	defer s.hashMu.Unlock()
	out := make(map[string]string)
	for _, f := range files {
		sha, err := s.hashViaCache(f.rel, f.path, walkStartNs)
		if err != nil {
			return nil, err
		}
		if sha != "" {
			out[f.rel] = sha
		}
	}
	return out, nil
}

// HashFile returns one file's sha256-hex; ("", nil) when absent. Shares
// [FileStore.HashTree]'s cache.
func (s *FileStore) HashFile(rel string) (string, error) {
	dest, err := s.resolveUnderRoot(rel)
	if err != nil {
		return "", err
	}
	st, err := os.Lstat(dest)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	if st.Mode()&fs.ModeSymlink != 0 {
		return "", &FileStoreError{Reason: ReasonIsASymlink, RelPath: rel}
	}
	if !st.Mode().IsRegular() {
		return "", &FileStoreError{Reason: ReasonNotAFile, RelPath: rel}
	}
	relKey, err := filepath.Rel(s.root, dest)
	if err != nil {
		return "", err
	}
	s.hashMu.Lock()
	defer s.hashMu.Unlock()
	return s.hashViaCache(filepath.ToSlash(relKey), dest, nowNs())
}

// Move renames src to dst; an existing dst is refused. The banned store root
// as either end does nothing.
func (s *FileStore) Move(src, dst string) error {
	from, err := s.resolveUnderRoot(src)
	if err != nil {
		return err
	}
	to, err := s.resolveUnderRoot(dst)
	if err != nil {
		return err
	}
	if from == s.root || to == s.root {
		return nil
	}
	// Stat, not Lstat: a dangling symlink at dst reads absent and is replaced
	// by the rename.
	if _, err := os.Stat(to); err == nil {
		return &FileStoreError{Reason: ReasonMoveDestinationExists, RelPath: dst}
	}
	if err := makeDirsBelowRoot(s.root, filepath.Dir(to)); err != nil {
		return err
	}
	return os.Rename(from, to)
}

// Remove deletes a file or subtree; absent — and the banned store root — do nothing.
func (s *FileStore) Remove(rel string) error {
	dest, err := s.resolveUnderRoot(rel)
	if err != nil {
		return err
	}
	if dest == s.root {
		return nil
	}
	// os.RemoveAll unlinks a symlink rather than following it, so Remove
	// cannot reach through one.
	return os.RemoveAll(dest)
}

func (s *FileStore) resolveUnderRoot(rel string) (string, error) {
	norm := strings.ReplaceAll(rel, "\\", "/")
	if strings.HasPrefix(norm, "/") {
		norm = strings.TrimLeft(norm, "/")
	} else if filepath.IsAbs(rel) {
		return "", &FileStoreError{Reason: ReasonEscapesRoot, RelPath: rel}
	}
	var parts []string
	for _, p := range strings.Split(norm, "/") {
		switch p {
		case "", ".":
		case "..":
			return "", &FileStoreError{Reason: ReasonEscapesRoot, RelPath: rel}
		default:
			parts = append(parts, p)
		}
	}
	return filepath.Join(append([]string{s.root}, parts...)...), nil
}

func (s *FileStore) requireUTF8(rel string, data []byte) error {
	if s.utf8Only && !utf8.Valid(data) {
		return &FileStoreError{Reason: ReasonNotUTF8, RelPath: rel}
	}
	return nil
}

// hashViaCache hashes one regular file through the shared cache; the caller
// holds hashMu. ("", nil) means the file left this snapshot — vanished or
// swapped for a non-regular entry since the walk.
func (s *FileStore) hashViaCache(rel, path string, walkStartNs int64) (string, error) {
	st, err := os.Lstat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", nil // vanished since the walk: not in this snapshot
		}
		return "", err
	}
	if !st.Mode().IsRegular() {
		return "", nil
	}
	var sha string
	if cached, ok := s.hashes[rel]; ok && unchangedSinceHashed(cached, st) {
		sha = cached.sha
	} else if sha, err = hashFile(path); err != nil {
		var fse *FileStoreError
		if errors.Is(err, fs.ErrNotExist) || errors.As(err, &fse) {
			return "", nil
		}
		// FreeBSD reports EMLINK rather than ELOOP for O_NOFOLLOW.
		if errors.Is(err, syscall.ELOOP) || errors.Is(err, syscall.EMLINK) {
			return "", nil
		}
		return "", err
	}
	if oldEnoughToCache(st, walkStartNs) {
		s.hashes[rel] = hashed{mtimeNs: st.ModTime().UnixNano(), ctimeNs: ctimeNs(st), size: st.Size(), sha: sha}
	}
	return sha, nil
}

func makeDirAndAncestors(path string) error {
	var missing []string
	current := path
	for {
		if _, err := os.Stat(current); err == nil {
			break
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		missing = append(missing, current)
		current = parent
	}
	for i := len(missing) - 1; i >= 0; i-- {
		if err := os.Mkdir(missing[i], ownerOnlyDirMode); err != nil && !errors.Is(err, fs.ErrExist) {
			return err
		}
	}
	return nil
}

// makeDirsBelowRoot creates each missing directory from root down to dir,
// never root itself: only CreateRoot makes the root, so a write racing an
// rm -rf of the folder fails with ENOENT instead of re-creating it.
func makeDirsBelowRoot(root, dir string) error {
	rel, err := filepath.Rel(root, dir)
	if err != nil {
		return err
	}
	if rel == "." {
		return nil
	}
	current := root
	for _, part := range strings.Split(rel, string(filepath.Separator)) {
		current = filepath.Join(current, part)
		if err := os.Mkdir(current, ownerOnlyDirMode); err != nil && !errors.Is(err, fs.ErrExist) {
			return err
		}
	}
	return nil
}

// hashFile is a variable so tests can count invocations.
var hashFile = func(path string) (string, error) {
	f, err := openRegularFile(filepath.Base(path), path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// walkedFile is one regular file a walk found.
type walkedFile struct {
	rel  string
	path string
}

// filenamesInDir returns every regular file under base as (rel, path); an
// absent base is empty, a present non-directory is refused (under is only
// for refusal messages).
func filenamesInDir(root, under, base string) ([]walkedFile, error) {
	st, err := os.Lstat(base)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) || errors.Is(err, syscall.ENOTDIR) {
			return nil, nil
		}
		return nil, err
	}
	if !st.IsDir() {
		return nil, &FileStoreError{Reason: ReasonNotADirectory, RelPath: under}
	}
	var out []walkedFile
	err = filepath.WalkDir(base, func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if errors.Is(walkErr, fs.ErrNotExist) {
				return nil // a listing is a snapshot, not a lock
			}
			return walkErr
		}
		// Non-regular entries are skipped; WalkDir never descends into a symlink.
		if !d.Type().IsRegular() {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		out = append(out, walkedFile{rel: filepath.ToSlash(rel), path: p})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func symlinksInDir(root, under, base string) (map[string]bool, error) {
	out := make(map[string]bool)
	st, err := os.Lstat(base)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) || errors.Is(err, syscall.ENOTDIR) {
			return out, nil
		}
		return nil, err
	}
	if st.Mode()&fs.ModeSymlink != 0 {
		rel, err := filepath.Rel(root, base)
		if err != nil {
			return nil, err
		}
		out[filepath.ToSlash(rel)] = true
		return out, nil
	}
	if !st.IsDir() {
		return nil, &FileStoreError{Reason: ReasonNotADirectory, RelPath: under}
	}
	err = filepath.WalkDir(base, func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if errors.Is(walkErr, fs.ErrNotExist) {
				return nil
			}
			return walkErr
		}
		// Symlinks to directories are visited as non-dir entries, never descended.
		if d.Type()&fs.ModeSymlink == 0 {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		out[filepath.ToSlash(rel)] = true
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func unchangedSinceHashed(cached hashed, st os.FileInfo) bool {
	return st.ModTime().UnixNano() == cached.mtimeNs && ctimeNs(st) == cached.ctimeNs && st.Size() == cached.size
}

// timestampTrustMarginNs: filesystems stamp times with coarse clocks, so a
// rewrite shortly after a hashed write can reuse the exact stamps. Files
// younger than the margin are simply re-hashed next walk.
var timestampTrustMarginNs int64 = 2_000_000_000

// nowNs is the walk clock; tests freeze it here.
var nowNs = func() int64 { return time.Now().UnixNano() }

func oldEnoughToCache(st os.FileInfo, walkStartNs int64) bool {
	return max(st.ModTime().UnixNano(), ctimeNs(st)) < walkStartNs-timestampTrustMarginNs
}
