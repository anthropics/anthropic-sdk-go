package agenttoolset

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// This file holds the skill-archive extraction logic, kept separate from the
// skill download/setup flow in skills.go and from the agent_toolset_20260401
// tool implementations.

// Bounds on archive extraction to guard against decompression bombs. These are
// vars rather than consts only so tests can lower them and exercise the limits
// without building multi-gigabyte archives; the public surface is unchanged.
var (
	skillArchiveMaxMembers       = 10_000
	skillArchiveMaxBytes   int64 = 1 << 30 // 1 GiB
)

// extractSkillArchive extracts the skill download at archivePath (a zip or
// gzip/bzip2/plain tar archive) into dest, refusing any member that would
// escape dest (zip-slip / tar-slip): skill archives come from the API, but
// skills can be third-party. The archive is read straight from disk — never
// buffered whole into memory — so a large skill bundle cannot OOM the runner.
func extractSkillArchive(archivePath, dest string) (retErr error) {
	root, err := filepath.Abs(dest)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return err
	}
	// extractSkillArchive creates dest, so it owns cleaning it up on failure:
	// a half-extracted skill (over the member/byte cap, a corrupt member,
	// zip-slip rejection, disk full, ...) is worse than none. Best effort —
	// dest is either a complete extraction or absent.
	defer func() {
		if retErr != nil {
			_ = os.RemoveAll(root)
		}
	}()

	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Sniff the leading bytes for the zip magic ("PK\x03\x04") without reading
	// the whole archive, then rewind for the real extractor.
	var magic [4]byte
	n, _ := io.ReadFull(f, magic[:])
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return err
	}
	if n >= 4 && magic[0] == 'P' && magic[1] == 'K' && magic[2] == 0x03 && magic[3] == 0x04 {
		// zip needs random access over the file; hand it the path so it seeks
		// the file on disk rather than us buffering it into memory.
		return extractZip(archivePath, root)
	}
	return extractTar(f, root)
}

// archiveTopDir returns the single top-level directory shared by every entry
// in names, or "" if the entries don't all live under one common directory.
//
// Skill bundles are packaged wrapped in one directory named after the skill
// (e.g. pdf/SKILL.md, pdf/scripts/...). The extractor strips that wrapper so
// contents land directly in the skill's destination dir instead of a redundant
// nested <skill>/<skill>/ level. A flat or multi-root archive yields "" and is
// extracted unchanged.
// cleanSegments splits a slash path into its components, dropping empty and
// "." segments so a "./pdf/x"-style name (e.g. from `tar -C dir .`) is treated
// the same as "pdf/x".
func cleanSegments(name string) []string {
	out := make([]string, 0, 4)
	for _, p := range strings.Split(filepath.ToSlash(name), "/") {
		if p != "" && p != "." {
			out = append(out, p)
		}
	}
	return out
}

func archiveTopDir(names []string) string {
	top := ""
	seen := false
	nested := false
	for _, n := range names {
		parts := cleanSegments(n)
		if len(parts) == 0 {
			continue
		}
		if !seen {
			top, seen = parts[0], true
		} else if parts[0] != top {
			return ""
		}
		if len(parts) > 1 {
			nested = true
		}
	}
	if seen && nested {
		return top
	}
	return ""
}

// stripTopDir drops the leading top component from name. Returns "" for the
// bare top-dir entry itself (nothing to write).
func stripTopDir(name, top string) string {
	if top == "" {
		return name
	}
	parts := cleanSegments(name)
	if len(parts) == 0 || parts[0] != top {
		return name
	}
	return strings.Join(parts[1:], "/") // "" for the bare top-dir entry itself
}

// tarDecompressor returns a reader over f's (optionally gzip/bzip2-compressed)
// tar bytes. f must be positioned at the start; the caller closes the returned
// closer when done with this pass.
func tarDecompressor(f *os.File) (io.Reader, func() error, error) {
	var magic [3]byte
	n, _ := io.ReadFull(f, magic[:])
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, nil, err
	}
	switch {
	case n >= 2 && magic[0] == 0x1f && magic[1] == 0x8b: // gzip
		gz, err := gzip.NewReader(f)
		if err != nil {
			return nil, nil, err
		}
		return gz, gz.Close, nil
	case n >= 3 && magic[0] == 'B' && magic[1] == 'Z' && magic[2] == 'h': // bzip2
		return bzip2.NewReader(f), func() error { return nil }, nil
	}
	return f, func() error { return nil }, nil
}

// safeJoin joins name under root, anchoring at "/" first so a ".." component
// is absorbed at the root rather than escaping it. Returns an error if the
// result still somehow lands outside root.
func safeJoin(root, name string) (string, error) {
	clean := filepath.Clean("/" + filepath.ToSlash(name))
	target := filepath.Join(root, clean)
	if target != root && !strings.HasPrefix(target, root+string(os.PathSeparator)) {
		return "", fmt.Errorf("refusing archive member %q", name)
	}
	return target, nil
}

func extractZip(archivePath, root string) error {
	zr, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer zr.Close()
	// zr.File is fully materialized at OpenReader; bound it before building the
	// names list so pass 1 (top-dir detection) is not unbounded on a zip bomb.
	if len(zr.File) > skillArchiveMaxMembers {
		return fmt.Errorf("skill archive exceeds %d members", skillArchiveMaxMembers)
	}
	names := make([]string, 0, len(zr.File))
	for _, f := range zr.File {
		names = append(names, f.Name)
	}
	top := archiveTopDir(names)
	members := 0
	var remaining int64 = skillArchiveMaxBytes
	for _, f := range zr.File {
		nm := stripTopDir(f.Name, top)
		if nm == "" {
			continue
		}
		members++
		if members > skillArchiveMaxMembers {
			return fmt.Errorf("skill archive exceeds %d members", skillArchiveMaxMembers)
		}
		target, err := safeJoin(root, nm)
		if err != nil {
			return err
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		written, err := copyZipEntry(f, target, remaining)
		if err != nil {
			return err
		}
		remaining -= written
	}
	return nil
}

// copyZipEntry decompresses f to target, capped at limit bytes; returns bytes
// written.
func copyZipEntry(f *zip.File, target string, limit int64) (int64, error) {
	rc, err := f.Open()
	if err != nil {
		return 0, err
	}
	defer rc.Close()
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, archiveFileMode(f.Mode()))
	if err != nil {
		return 0, err
	}
	defer out.Close()
	return copyBounded(out, rc, limit)
}

// extractTar streams a (optionally gzip/bzip2-compressed) tar archive from f
// into root. f must be positioned at the start of the archive; the tar reader
// consumes it as a stream, so the archive is never buffered whole in memory.
func extractTar(f *os.File, root string) error {
	skip := func(flag byte) bool {
		switch flag {
		case tar.TypeSymlink, tar.TypeLink, tar.TypeChar, tar.TypeBlock, tar.TypeFifo:
			return true // skip symlinks / hardlinks / devices
		}
		return false
	}

	// Pass 1: read headers only to detect the skill bundle's wrapper directory.
	// The archive is on disk, so a second pass just rewinds and re-reads it —
	// nothing is buffered whole in memory.
	r, closeR, err := tarDecompressor(f)
	if err != nil {
		return err
	}
	// Bound pass 1 too: tr.Next() reads through every preceding entry's body to
	// advance, so an unbounded pass 1 would burn CPU/wall time on a tar bomb
	// before pass 2's caps ever run. The byte budget is computed from hdr.Size
	// (the declared size) so we don't have to read bodies to enforce it.
	var names []string
	tr := tar.NewReader(r)
	members := 0
	var declared int64
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			_ = closeR()
			return err
		}
		if skip(hdr.Typeflag) {
			continue
		}
		members++
		if members > skillArchiveMaxMembers {
			_ = closeR()
			return fmt.Errorf("skill archive exceeds %d members", skillArchiveMaxMembers)
		}
		declared += hdr.Size
		if declared > skillArchiveMaxBytes {
			_ = closeR()
			return fmt.Errorf("skill archive exceeds %d bytes decompressed", skillArchiveMaxBytes)
		}
		names = append(names, hdr.Name)
	}
	if err := closeR(); err != nil {
		return err
	}
	top := archiveTopDir(names)

	// Pass 2: extract, stripping the wrapper directory.
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return err
	}
	r, closeR, err = tarDecompressor(f)
	if err != nil {
		return err
	}
	defer closeR()
	tr = tar.NewReader(r)
	members = 0
	var remaining int64 = skillArchiveMaxBytes
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if skip(hdr.Typeflag) {
			continue
		}
		nm := stripTopDir(hdr.Name, top)
		if nm == "" {
			continue
		}
		members++
		if members > skillArchiveMaxMembers {
			return fmt.Errorf("skill archive exceeds %d members", skillArchiveMaxMembers)
		}
		target, err := safeJoin(root, nm)
		if err != nil {
			return err
		}
		if hdr.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		written, err := copyTarEntry(tr, target, hdr.FileInfo().Mode(), remaining)
		if err != nil {
			return err
		}
		remaining -= written
	}
	return nil
}

// copyTarEntry decompresses tr's current entry to target, capped at limit
// bytes; returns bytes written.
func copyTarEntry(tr *tar.Reader, target string, entryMode os.FileMode, limit int64) (int64, error) {
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, archiveFileMode(entryMode))
	if err != nil {
		return 0, err
	}
	defer out.Close()
	return copyBounded(out, tr, limit)
}

// archiveFileMode reduces an archive entry's mode to 0o755 if executable,
// 0o644 otherwise.
func archiveFileMode(entryMode os.FileMode) os.FileMode {
	if entryMode&0o111 != 0 {
		return 0o755
	}
	return 0o644
}

// copyBounded copies up to limit bytes from src to dst, returning an error if
// src exceeds limit.
func copyBounded(dst io.Writer, src io.Reader, limit int64) (int64, error) {
	if limit < 0 {
		limit = 0
	}
	n, err := io.Copy(dst, io.LimitReader(src, limit+1))
	if err != nil {
		return n, err
	}
	if n > limit {
		return n, fmt.Errorf("skill archive exceeds %d bytes decompressed", skillArchiveMaxBytes)
	}
	return n, nil
}
