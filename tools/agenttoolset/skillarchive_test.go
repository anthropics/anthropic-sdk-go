package agenttoolset

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func zipBytes(t *testing.T, entries map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, body := range entries {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatalf("zip create %q: %v", name, err)
		}
		if _, err := w.Write([]byte(body)); err != nil {
			t.Fatalf("zip write %q: %v", name, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("zip close: %v", err)
	}
	return buf.Bytes()
}

func tarGzBytes(t *testing.T, entries map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	for name, body := range entries {
		if err := tw.WriteHeader(&tar.Header{Name: name, Mode: 0o644, Size: int64(len(body)), Typeflag: tar.TypeReg}); err != nil {
			t.Fatalf("tar header %q: %v", name, err)
		}
		if _, err := tw.Write([]byte(body)); err != nil {
			t.Fatalf("tar write %q: %v", name, err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("tar close: %v", err)
	}
	if err := gw.Close(); err != nil {
		t.Fatalf("gzip close: %v", err)
	}
	return buf.Bytes()
}

// extractArchiveBytes writes archive bytes to a temp file and extracts it,
// exercising the on-disk extraction path the real skill download now uses.
func extractArchiveBytes(t *testing.T, data []byte, dest string) error {
	t.Helper()
	path := filepath.Join(t.TempDir(), "archive")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write archive: %v", err)
	}
	return extractSkillArchive(path, dest)
}

func TestExtractSkillArchive_Zip(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "s")
	if err := extractArchiveBytes(t, zipBytes(t, map[string]string{"SKILL.md": "hi", "d/x.txt": "x"}), dest); err != nil {
		t.Fatalf("extract: %v", err)
	}
	if got, _ := os.ReadFile(filepath.Join(dest, "SKILL.md")); string(got) != "hi" {
		t.Fatalf("SKILL.md = %q", got)
	}
	if got, _ := os.ReadFile(filepath.Join(dest, "d", "x.txt")); string(got) != "x" {
		t.Fatalf("d/x.txt = %q", got)
	}
}

func TestExtractSkillArchive_TarGz(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "s")
	if err := extractArchiveBytes(t, tarGzBytes(t, map[string]string{"SKILL.md": "tarred", "sub/y": "Y"}), dest); err != nil {
		t.Fatalf("extract: %v", err)
	}
	if got, _ := os.ReadFile(filepath.Join(dest, "SKILL.md")); string(got) != "tarred" {
		t.Fatalf("SKILL.md = %q", got)
	}
	if got, _ := os.ReadFile(filepath.Join(dest, "sub", "y")); string(got) != "Y" {
		t.Fatalf("sub/y = %q", got)
	}
}

func TestArchiveTopDir(t *testing.T) {
	cases := []struct {
		desc  string
		names []string
		want  string
	}{
		{"single wrapped skill", []string{"pdf/SKILL.md", "pdf/scripts/x.py"}, "pdf"},
		{"single file under wrapper", []string{"pdf/SKILL.md"}, "pdf"},
		{"dot-prefixed (tar -C dir .)", []string{"./", "./pdf/", "./pdf/SKILL.md", "./pdf/scripts/x"}, "pdf"},
		{"flat archive", []string{"SKILL.md", "scripts/x.py"}, ""},
		{"multiple roots", []string{"a/x", "b/y"}, ""},
		{"only bare top dir", []string{"pdf/"}, ""},
		{"empty", nil, ""},
	}
	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			if got := archiveTopDir(c.names); got != c.want {
				t.Errorf("archiveTopDir(%v) = %q, want %q", c.names, got, c.want)
			}
		})
	}
}

func TestStripTopDir(t *testing.T) {
	cases := []struct{ name, top, want string }{
		{"pdf/SKILL.md", "pdf", "SKILL.md"},
		{"pdf/scripts/x.py", "pdf", "scripts/x.py"},
		{"pdf", "pdf", ""},                    // bare top-dir entry
		{"./pdf/SKILL.md", "pdf", "SKILL.md"}, // dot-prefixed (tar -C dir .)
		{"./pdf/", "pdf", ""},
		{"SKILL.md", "", "SKILL.md"}, // no wrapper -> unchanged
		{"other/x", "pdf", "other/x"},
	}
	for _, c := range cases {
		if got := stripTopDir(c.name, c.top); got != c.want {
			t.Errorf("stripTopDir(%q,%q) = %q, want %q", c.name, c.top, got, c.want)
		}
	}
}

// Skill bundles are wrapped in a directory named after the skill; extraction
// must strip it so files land at dest/SKILL.md, not the doubled
// dest/<skill>/SKILL.md.
func TestExtractSkillArchive_StripsWrapperDir(t *testing.T) {
	for _, tc := range []struct {
		kind string
		make func(*testing.T, map[string]string) []byte
	}{
		{"zip", zipBytes},
		{"targz", tarGzBytes},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			dest := filepath.Join(t.TempDir(), "skills", "pdf")
			// Dot-prefixed names mirror a real `tar -C dir .` / zip of "."
			// bundle; the wrapper must still be detected and stripped.
			data := tc.make(t, map[string]string{
				"./pdf/SKILL.md":       "# PDF",
				"./pdf/scripts/run.py": "print(1)",
			})
			if err := extractArchiveBytes(t, data, dest); err != nil {
				t.Fatalf("extract: %v", err)
			}
			if got, _ := os.ReadFile(filepath.Join(dest, "SKILL.md")); string(got) != "# PDF" {
				t.Fatalf("SKILL.md = %q, want %q", got, "# PDF")
			}
			if got, _ := os.ReadFile(filepath.Join(dest, "scripts", "run.py")); string(got) != "print(1)" {
				t.Fatalf("scripts/run.py = %q", got)
			}
			if _, err := os.Stat(filepath.Join(dest, "pdf")); err == nil {
				t.Fatal("wrapper dir was not stripped (skills/pdf/pdf/ doubling)")
			}
		})
	}
}

// setSkillArchiveLimits temporarily lowers the extraction bounds so tests can
// exercise the limits without building 10K-member / 1 GiB archives. Restored
// via t.Cleanup; tests that use this must not run in parallel.
func setSkillArchiveLimits(t *testing.T, members int, byteLimit int64) {
	t.Helper()
	oldM, oldB := skillArchiveMaxMembers, skillArchiveMaxBytes
	skillArchiveMaxMembers, skillArchiveMaxBytes = members, byteLimit
	t.Cleanup(func() { skillArchiveMaxMembers, skillArchiveMaxBytes = oldM, oldB })
}

// requireDestAbsent asserts the extraction destination does not exist — a
// failed extraction must remove the partially written directory rather than
// leave a half-extracted skill on disk.
func requireDestAbsent(t *testing.T, dest string) {
	t.Helper()
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		t.Fatalf("dest %q should have been removed after extraction failure, stat err = %v", dest, err)
	}
}

func TestExtractSkillArchive_FailureRemovesDest(t *testing.T) {
	cases := []struct {
		desc string
		data func(t *testing.T) []byte
	}{
		{
			// Member-count cap fires before any file is written (zip pass 1).
			desc: "over member limit",
			data: func(t *testing.T) []byte {
				return zipBytes(t, map[string]string{"a": "1", "b": "2", "c": "3"})
			},
		},
		{
			// Byte budget fires mid-extraction, after some files have already
			// been written into dest.
			desc: "over byte budget mid-extraction",
			data: func(t *testing.T) []byte {
				return zipBytes(t, map[string]string{"a": "abcdef", "b": "ghijkl"})
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			setSkillArchiveLimits(t, 2, 8)
			dest := filepath.Join(t.TempDir(), "s")
			if err := extractArchiveBytes(t, tc.data(t), dest); err == nil {
				t.Fatal("expected extraction to fail, got nil")
			}
			requireDestAbsent(t, dest)
		})
	}
}

// TestExtractTar_Pass1MemberBound proves the member-count cap fires during
// pass 1 (top-dir detection) rather than only in pass 2. The archive ends in a
// garbage block instead of the usual end-of-archive trailer: a bounded pass 1
// stops at member N+1 with the "exceeds N members" error, while an unbounded
// pass 1 would read through to the garbage and surface a tar parse error.
func TestExtractTar_Pass1MemberBound(t *testing.T) {
	setSkillArchiveLimits(t, 3, skillArchiveMaxBytes)

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for i := 0; i <= skillArchiveMaxMembers; i++ {
		hdr := &tar.Header{Name: fmt.Sprintf("f%d", i), Mode: 0o644, Size: 0, Typeflag: tar.TypeReg}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatalf("tar header: %v", err)
		}
	}
	if err := tw.Flush(); err != nil {
		t.Fatalf("tar flush: %v", err)
	}
	buf.Write(bytes.Repeat([]byte{0xFF}, 512))

	dest := filepath.Join(t.TempDir(), "s")
	err := extractArchiveBytes(t, buf.Bytes(), dest)
	if err == nil {
		t.Fatal("expected member-count error, got nil")
	}
	if !strings.Contains(err.Error(), "members") {
		t.Fatalf("expected pass-1 member-count error, got: %v", err)
	}
	requireDestAbsent(t, dest)
}

// TestExtractTar_Pass1ByteBound proves the declared-size byte budget fires
// during pass 1. The budget is computed from hdr.Size alone, so the test only
// needs entries whose declared sizes sum past the (lowered) limit. As above,
// the trailing garbage block proves pass 1 stopped at the cap rather than
// reading through the whole archive.
func TestExtractTar_Pass1ByteBound(t *testing.T) {
	setSkillArchiveLimits(t, skillArchiveMaxMembers, 8)

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	body := []byte("xxxxx") // 5 bytes; 3 entries declare 15 > 8
	for i := 0; i < 3; i++ {
		hdr := &tar.Header{Name: fmt.Sprintf("f%d", i), Mode: 0o644, Size: int64(len(body)), Typeflag: tar.TypeReg}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatalf("tar header: %v", err)
		}
		if _, err := tw.Write(body); err != nil {
			t.Fatalf("tar write: %v", err)
		}
	}
	if err := tw.Flush(); err != nil {
		t.Fatalf("tar flush: %v", err)
	}
	buf.Write(bytes.Repeat([]byte{0xFF}, 512))

	dest := filepath.Join(t.TempDir(), "s")
	err := extractArchiveBytes(t, buf.Bytes(), dest)
	if err == nil {
		t.Fatal("expected byte-budget error, got nil")
	}
	if !strings.Contains(err.Error(), "bytes") {
		t.Fatalf("expected pass-1 byte-budget error, got: %v", err)
	}
	requireDestAbsent(t, dest)
}

type tarMember struct {
	hdr  tar.Header
	body string
}

// rawTarBytes writes members in order as an uncompressed tar stream. A member
// whose Typeflag is tar.TypeRegA is written as TypeReg and then patched back
// to the legacy '\x00' flag on the wire, since tar.Writer refuses to emit it.
func rawTarBytes(t *testing.T, members ...tarMember) []byte {
	t.Helper()
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	var legacyOffsets []int
	for _, m := range members {
		hdr := m.hdr
		hdr.Size = int64(len(m.body))
		if hdr.Typeflag == tar.TypeRegA {
			hdr.Typeflag = tar.TypeReg
			require.NoError(t, tw.Flush())
			legacyOffsets = append(legacyOffsets, buf.Len())
		}
		require.NoError(t, tw.WriteHeader(&hdr), "tar header %q", hdr.Name)
		if len(m.body) > 0 {
			_, err := tw.Write([]byte(m.body))
			require.NoError(t, err, "tar write %q", hdr.Name)
		}
	}
	require.NoError(t, tw.Close())
	out := buf.Bytes()
	for _, off := range legacyOffsets {
		block := out[off : off+512]
		block[156] = tar.TypeRegA
		copy(block[148:156], "        ")
		var sum int
		for _, b := range block {
			sum += int(b)
		}
		copy(block[148:156], fmt.Sprintf("%06o\x00 ", sum))
	}
	return out
}

type zipMember struct {
	hdr  zip.FileHeader
	body string
}

// Unix st_mode type bits as they appear in a zip entry's external attributes.
const (
	zipUnixHost = 3 << 8
	sIFDIR      = 0o040000
	sIFIFO      = 0o010000
	sIFLNK      = 0o120000
)

func rawZipBytes(t *testing.T, members ...zipMember) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for _, m := range members {
		hdr := m.hdr
		w, err := zw.CreateHeader(&hdr)
		require.NoError(t, err, "zip header %q", hdr.Name)
		_, err = w.Write([]byte(m.body))
		require.NoError(t, err, "zip write %q", hdr.Name)
	}
	require.NoError(t, zw.Close())
	return buf.Bytes()
}

// requireOnlyPlainEntries asserts that every entry under dir, checked with
// lstat, is a regular file or a directory.
func requireOnlyPlainEntries(t *testing.T, dir string) {
	t.Helper()
	require.NoError(t, filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		require.NoError(t, err)
		require.True(t, d.Type().IsRegular() || d.IsDir(), "%s is neither a regular file nor a directory (%v)", p, d.Type())
		return nil
	}))
}

func TestExtractSkillArchive_SkipsSpecialMembers(t *testing.T) {
	tests := []struct {
		description string
		data        func(t *testing.T) []byte
		wantFiles   map[string]string
		wantDirs    []string
		wantModes   map[string]os.FileMode
		wantAbsent  []string
	}{
		{
			description: "tar symlink, hardlink, FIFO and device members are skipped and the rest extracted with the wrapper stripped",
			data: func(t *testing.T) []byte {
				return rawTarBytes(t,
					tarMember{hdr: tar.Header{Name: "pdf/SKILL.md", Typeflag: tar.TypeReg}, body: "# PDF"},
					tarMember{hdr: tar.Header{Name: "pdf/scripts/run.sh", Typeflag: tar.TypeReg, Mode: 0o755}, body: "#!/bin/sh"},
					tarMember{hdr: tar.Header{Name: "pdf/old", Typeflag: tar.TypeRegA}, body: "legacy"},
					tarMember{hdr: tar.Header{Name: "pdf/abs", Typeflag: tar.TypeSymlink, Linkname: "/etc/passwd"}},
					tarMember{hdr: tar.Header{Name: "pdf/rel", Typeflag: tar.TypeSymlink, Linkname: "SKILL.md"}},
					tarMember{hdr: tar.Header{Name: "pdf/hl", Typeflag: tar.TypeLink, Linkname: "pdf/SKILL.md"}},
					tarMember{hdr: tar.Header{Name: "pdf/fifo", Typeflag: tar.TypeFifo}},
					tarMember{hdr: tar.Header{Name: "pdf/dev", Typeflag: tar.TypeChar, Devmajor: 1, Devminor: 3}},
				)
			},
			wantFiles:  map[string]string{"SKILL.md": "# PDF", "scripts/run.sh": "#!/bin/sh", "old": "legacy"},
			wantModes:  map[string]os.FileMode{"scripts/run.sh": 0o755},
			wantAbsent: []string{"abs", "rel", "hl", "fifo", "dev", "pdf"},
		},
		{
			description: "a top-level tar symlink does not defeat wrapper stripping",
			data: func(t *testing.T) []byte {
				return rawTarBytes(t,
					tarMember{hdr: tar.Header{Name: "link", Typeflag: tar.TypeSymlink, Linkname: "/tmp"}},
					tarMember{hdr: tar.Header{Name: "pdf/SKILL.md", Typeflag: tar.TypeReg}, body: "# PDF"},
				)
			},
			wantFiles:  map[string]string{"SKILL.md": "# PDF"},
			wantAbsent: []string{"link", "pdf"},
		},
		{
			description: "a git-archive style PAX global header is not written out as a file and does not defeat wrapper stripping",
			data: func(t *testing.T) []byte {
				return rawTarBytes(t,
					tarMember{hdr: tar.Header{Name: "pax_global_header", Typeflag: tar.TypeXGlobalHeader, PAXRecords: map[string]string{"comment": "0123abcd"}, Format: tar.FormatPAX}},
					tarMember{hdr: tar.Header{Name: "pdf/SKILL.md", Typeflag: tar.TypeReg}, body: "# PDF"},
				)
			},
			wantFiles:  map[string]string{"SKILL.md": "# PDF"},
			wantAbsent: []string{"pax_global_header", "pdf"},
		},
		{
			description: "a tar symlink named ../x is skipped and nothing named x appears beside dest",
			data: func(t *testing.T) []byte {
				return rawTarBytes(t,
					tarMember{hdr: tar.Header{Name: "../x", Typeflag: tar.TypeSymlink, Linkname: "/etc"}},
					tarMember{hdr: tar.Header{Name: "pdf/SKILL.md", Typeflag: tar.TypeReg}, body: "# PDF"},
				)
			},
			wantFiles:  map[string]string{"SKILL.md": "# PDF"},
			wantAbsent: []string{"x", "../x", "pdf"},
		},
		{
			description: "zip Unix-host symlink and FIFO entries are skipped, not written out as regular files",
			data: func(t *testing.T) []byte {
				return rawZipBytes(t,
					zipMember{hdr: zip.FileHeader{Name: "pdf/SKILL.md"}, body: "# PDF"},
					zipMember{hdr: zip.FileHeader{Name: "pdf/lnk", CreatorVersion: zipUnixHost, ExternalAttrs: (sIFLNK | 0o777) << 16}, body: "/etc/passwd"},
					zipMember{hdr: zip.FileHeader{Name: "pdf/fifo", CreatorVersion: zipUnixHost, ExternalAttrs: (sIFIFO | 0o644) << 16}},
				)
			},
			wantFiles:  map[string]string{"SKILL.md": "# PDF"},
			wantAbsent: []string{"lnk", "fifo", "pdf"},
		},
		{
			description: "zip symlink type bits on a non-Unix-host entry are not a mode, so the entry is a regular file",
			data: func(t *testing.T) []byte {
				return rawZipBytes(t,
					zipMember{hdr: zip.FileHeader{Name: "pdf/SKILL.md"}, body: "# PDF"},
					zipMember{hdr: zip.FileHeader{Name: "pdf/lnk", CreatorVersion: 0, ExternalAttrs: (sIFLNK | 0o777) << 16}, body: "/etc/passwd"},
				)
			},
			wantFiles:  map[string]string{"SKILL.md": "# PDF", "lnk": "/etc/passwd"},
			wantAbsent: []string{"pdf"},
		},
		{
			description: "zip Unix-host entries with zero type bits are files and S_IFDIR without a trailing slash is a directory",
			data: func(t *testing.T) []byte {
				return rawZipBytes(t,
					zipMember{hdr: zip.FileHeader{Name: "pdf/SKILL.md", CreatorVersion: zipUnixHost, ExternalAttrs: 0o600 << 16}, body: "# PDF"},
					zipMember{hdr: zip.FileHeader{Name: "pdf/run.sh", CreatorVersion: zipUnixHost, ExternalAttrs: 0o755 << 16}, body: "#!/bin/sh"},
					zipMember{hdr: zip.FileHeader{Name: "pdf/sub", CreatorVersion: zipUnixHost, ExternalAttrs: (sIFDIR | 0o755) << 16}},
				)
			},
			wantFiles:  map[string]string{"SKILL.md": "# PDF", "run.sh": "#!/bin/sh"},
			wantModes:  map[string]os.FileMode{"run.sh": 0o755},
			wantDirs:   []string{"sub"},
			wantAbsent: []string{"pdf"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			parent := t.TempDir()
			dest := filepath.Join(parent, "s")
			require.NoError(t, extractArchiveBytes(t, tc.data(t), dest))

			for rel, want := range tc.wantFiles {
				fi, err := os.Lstat(filepath.Join(dest, rel))
				require.NoError(t, err)
				require.True(t, fi.Mode().IsRegular(), "%s must be a regular file, got %v", rel, fi.Mode())
				got, err := os.ReadFile(filepath.Join(dest, rel))
				require.NoError(t, err)
				require.Equal(t, want, string(got), rel)
			}
			if runtime.GOOS != "windows" {
				for rel, want := range tc.wantModes {
					fi, err := os.Lstat(filepath.Join(dest, rel))
					require.NoError(t, err)
					require.Equal(t, want, fi.Mode().Perm(), rel)
				}
			}
			for _, rel := range tc.wantDirs {
				fi, err := os.Lstat(filepath.Join(dest, rel))
				require.NoError(t, err)
				require.True(t, fi.IsDir(), "%s must be a directory, got %v", rel, fi.Mode())
			}
			for _, rel := range tc.wantAbsent {
				_, err := os.Lstat(filepath.Join(dest, rel))
				require.ErrorIs(t, err, os.ErrNotExist, "%s must not exist", rel)
			}
			requireOnlyPlainEntries(t, dest)
			siblings, err := os.ReadDir(parent)
			require.NoError(t, err)
			require.Len(t, siblings, 1, "nothing may be created outside dest")
		})
	}
}

func TestExtractSkillArchive_ConfinesTraversal(t *testing.T) {
	parent := t.TempDir()
	dest := filepath.Join(parent, "nested", "s")
	if err := extractArchiveBytes(t, zipBytes(t, map[string]string{"../evil.txt": "pwn", "ok.txt": "fine"}), dest); err != nil {
		t.Fatalf("extract: %v", err)
	}
	if _, err := os.Stat(filepath.Join(parent, "nested", "evil.txt")); err == nil {
		t.Fatal("zip-slip member escaped the destination directory")
	}
	if got, _ := os.ReadFile(filepath.Join(dest, "evil.txt")); string(got) != "pwn" {
		t.Fatalf("traversal member not confined into dest: %q", got)
	}
	if got, _ := os.ReadFile(filepath.Join(dest, "ok.txt")); string(got) != "fine" {
		t.Fatalf("ok.txt = %q", got)
	}
}
