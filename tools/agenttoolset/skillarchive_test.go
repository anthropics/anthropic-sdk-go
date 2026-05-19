package agenttoolset

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
