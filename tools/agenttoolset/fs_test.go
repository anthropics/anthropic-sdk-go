package agenttoolset

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/stretchr/testify/require"
)

func mustJSON(t *testing.T, v any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return b
}

func TestExecRead(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(work, "a.txt"), []byte("line1\nline2\nline3"), 0o644))
	env := &AgentToolContext{Workdir: work}

	tests := []struct {
		description string
		input       map[string]any
		want        string
		wantErr     bool
	}{
		{
			description: "reading an existing file returns its full contents with no error flag",
			input:       map[string]any{"file_path": "a.txt"},
			want:        "line1\nline2\nline3",
		},
		{
			description: "view_range slices the file by 1-indexed inclusive line numbers like the 1P read tool",
			input:       map[string]any{"file_path": "a.txt", "view_range": []int{2, 2}},
			want:        "line2",
		},
		{
			description: "reading a missing file surfaces the os error and sets is_error so the model can retry",
			input:       map[string]any{"file_path": "missing.txt"},
			wantErr:     true,
		},
		{
			description: "empty file_path is rejected with a clear validation message before touching the filesystem",
			input:       map[string]any{},
			wantErr:     true,
		},
		{
			description: "file_path that escapes the workdir is rejected by the resolvePath jail",
			input:       map[string]any{"file_path": "../outside.txt"},
			wantErr:     true,
		},
		{
			description: "view_range with the wrong arity is rejected so the model gets a clear error",
			input:       map[string]any{"file_path": "a.txt", "view_range": []int{2}},
			wantErr:     true,
		},
		{
			description: "inverted view_range selects nothing and returns empty content rather than an error or a slice-bounds panic",
			input:       map[string]any{"file_path": "a.txt", "view_range": []int{3, 1}},
			want:        "",
		},
		{
			description: "empty view_range reads the whole file like an omitted one",
			input:       map[string]any{"file_path": "a.txt", "view_range": []int{}},
			want:        "line1\nline2\nline3",
		},
		{
			description: "end_line of 0 reads through to the end of the file",
			input:       map[string]any{"file_path": "a.txt", "view_range": []int{2, 0}},
			want:        "line2\nline3",
		},
		{
			description: "start_line past the end of the file returns empty content",
			input:       map[string]any{"file_path": "a.txt", "view_range": []int{10, 12}},
			want:        "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			got, isErr := runTool(t, BetaReadTool(env), mustJSON(t, tc.input))
			require.Equal(t, tc.wantErr, isErr, "is_error mismatch; output=%q", got)
			if !tc.wantErr {
				require.Equal(t, tc.want, got)
			}
		})
	}
}

// The read tool returns file bytes verbatim: CRLF and lone-CR line endings are
// never normalized, and view_range splits on "\n" only, so a trailing "\r"
// stays attached to its line.
func TestExecReadPreservesCRLF(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(work, "crlf.txt"), []byte("line1\r\nline2\r\nline3\r\n"), 0o644))
	tool := BetaReadTool(&AgentToolContext{Workdir: work})

	out, isErr := runTool(t, tool, mustJSON(t, map[string]any{"file_path": "crlf.txt"}))
	require.False(t, isErr, "output=%q", out)
	require.Equal(t, "line1\r\nline2\r\nline3\r\n", out)

	out, isErr = runTool(t, tool, mustJSON(t, map[string]any{"file_path": "crlf.txt", "view_range": []int{1, 2}}))
	require.False(t, isErr, "output=%q", out)
	require.Equal(t, "line1\r\nline2\r", out)
}

func TestExecReadPreservesLoneCR(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(work, "cr.txt"), []byte("a\rb\rc\r"), 0o644))
	tool := BetaReadTool(&AgentToolContext{Workdir: work})

	out, isErr := runTool(t, tool, mustJSON(t, map[string]any{"file_path": "cr.txt"}))
	require.False(t, isErr, "output=%q", out)
	require.Equal(t, "a\rb\rc\r", out)

	// A lone "\r" is not a line break, so the file is a single line and line 2
	// is past the end.
	out, isErr = runTool(t, tool, mustJSON(t, map[string]any{"file_path": "cr.txt", "view_range": []int{2, 2}}))
	require.False(t, isErr, "output=%q", out)
	require.Equal(t, "", out)
}

func TestExecWrite(t *testing.T) {
	work := t.TempDir()
	env := &AgentToolContext{Workdir: work}

	tests := []struct {
		description string
		input       map[string]any
		wantOnDisk  string
		wantErr     bool
	}{
		{
			description: "writing a new file creates it with the given content",
			input:       map[string]any{"file_path": "new.txt", "content": "hello"},
			wantOnDisk:  "hello",
		},
		{
			description: "writing into a nested directory creates the parent directories automatically",
			input:       map[string]any{"file_path": "deep/nested/f.txt", "content": "x"},
			wantOnDisk:  "x",
		},
		{
			description: "writing outside the workdir via dot-dot is rejected before any IO happens",
			input:       map[string]any{"file_path": "../evil.txt", "content": "x"},
			wantErr:     true,
		},
		{
			description: "missing file_path is rejected with a validation error",
			input:       map[string]any{"content": "x"},
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			_, isErr := execWrite(context.Background(), mustJSON(t, tc.input), env)
			require.Equal(t, tc.wantErr, isErr)
			if tc.wantErr {
				return
			}
			data, err := os.ReadFile(filepath.Join(work, tc.input["file_path"].(string)))
			require.NoError(t, err)
			require.Equal(t, tc.wantOnDisk, string(data))
		})
	}

	t.Run("writing the same path again overwrites rather than appending", func(t *testing.T) {
		_, isErr := execWrite(context.Background(), mustJSON(t, map[string]any{"file_path": "ow.txt", "content": "first"}), env)
		require.False(t, isErr)
		_, isErr = execWrite(context.Background(), mustJSON(t, map[string]any{"file_path": "ow.txt", "content": "second"}), env)
		require.False(t, isErr)
		data, _ := os.ReadFile(filepath.Join(work, "ow.txt"))
		require.Equal(t, "second", string(data))
	})
}

func TestExecReadRejectsSymlinkLoop(t *testing.T) {
	work, _ := symlinkLoopFixture(t)
	tool := BetaReadTool(&AgentToolContext{Workdir: work})

	out, isErr := runTool(t, tool, mustJSON(t, map[string]any{"file_path": "loop_a"}))
	require.True(t, isErr)
	require.Equal(t, `read: path "loop_a": too many levels of symbolic links`, out)

	for _, p := range []string{"loop_a/../evil_link", "L"} {
		out, isErr := runTool(t, tool, mustJSON(t, map[string]any{"file_path": p}))
		require.True(t, isErr, "file_path=%q output=%q", p, out)
		require.NotContains(t, out, "SECRET", "file_path=%q", p)
	}
}

func TestExecWriteRejectsSymlinkLoop(t *testing.T) {
	work, _ := symlinkLoopFixture(t)
	before, err := os.ReadDir(work)
	require.NoError(t, err)

	out, isErr := runTool(t, BetaWriteTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{
		"file_path": "loop_a/child.txt", "content": "x",
	}))
	require.True(t, isErr)
	require.Equal(t, `write: path "loop_a/child.txt": too many levels of symbolic links`, out)

	after, err := os.ReadDir(work)
	require.NoError(t, err)
	names := func(es []os.DirEntry) []string {
		out := make([]string, len(es))
		for i, e := range es {
			out[i] = e.Name()
		}
		return out
	}
	require.Equal(t, names(before), names(after), "a rejected write must not create anything in the workdir")
}

func TestExecWriteThroughDanglingInsideSymlink(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.Symlink(filepath.Join("newdir", "f.txt"), filepath.Join(work, "d")))

	out, isErr := runTool(t, BetaWriteTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{
		"file_path": "d", "content": "via-link",
	}))
	require.False(t, isErr, "output=%q", out)
	data, err := os.ReadFile(filepath.Join(work, "newdir", "f.txt"))
	require.NoError(t, err)
	require.Equal(t, "via-link", string(data))
}

func TestExecReadFollowsSymlinkChainInsideWorkdir(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(work, "real.txt"), []byte("target"), 0o644))
	require.NoError(t, os.Symlink("real.txt", filepath.Join(work, "c2")))
	require.NoError(t, os.Symlink("c2", filepath.Join(work, "c1")))
	require.NoError(t, os.Symlink("c1", filepath.Join(work, "c0")))

	out, isErr := runTool(t, BetaReadTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{"file_path": "c0"}))
	require.False(t, isErr, "output=%q", out)
	require.Equal(t, "target", out)
}

// Only symlink reads count against the resolver's hop cap, so a new file
// under more missing directories than the cap allows still resolves.
func TestExecWriteDeepMissingPathIsNotASymlinkLoop(t *testing.T) {
	work := t.TempDir()
	parts := make([]string, 0, 51)
	for i := range 50 {
		parts = append(parts, fmt.Sprintf("d%d", i))
	}
	parts = append(parts, "f.txt")
	rel := filepath.Join(parts...)

	out, isErr := runTool(t, BetaWriteTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{
		"file_path": rel, "content": "deep",
	}))
	require.False(t, isErr, "output=%q", out)
	data, err := os.ReadFile(filepath.Join(work, rel))
	require.NoError(t, err)
	require.Equal(t, "deep", string(data))
}

func TestExecReadUnderUnreadableDirReportsPermissionDenied(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("directory permission bits are not enforced for this user")
	}
	work := t.TempDir()
	noperm := filepath.Join(work, "noperm")
	require.NoError(t, os.Mkdir(noperm, 0o000))
	t.Cleanup(func() { _ = os.Chmod(noperm, 0o755) })

	out, isErr := runTool(t, BetaReadTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{"file_path": "noperm/x"}))
	require.True(t, isErr)
	require.Equal(t, `read: path "noperm/x": permission denied`, out)
}

func TestExecEdit(t *testing.T) {
	tests := []struct {
		description string
		initial     string
		input       map[string]any
		want        string
		wantErr     bool
	}{
		{
			description: "unique old_string is replaced once leaving the rest of the file untouched",
			initial:     "alpha\nbeta\ngamma",
			input:       map[string]any{"file_path": "f.txt", "old_string": "beta", "new_string": "BETA"},
			want:        "alpha\nBETA\ngamma",
		},
		{
			description: "ambiguous old_string without replace_all errors so the model must disambiguate",
			initial:     "x\nx\n",
			input:       map[string]any{"file_path": "f.txt", "old_string": "x", "new_string": "y"},
			wantErr:     true,
		},
		{
			description: "ambiguous old_string with replace_all true rewrites every occurrence",
			initial:     "x\nx\n",
			input:       map[string]any{"file_path": "f.txt", "old_string": "x", "new_string": "y", "replace_all": true},
			want:        "y\ny\n",
		},
		{
			description: "old_string not present in the file errors rather than silently writing an unchanged file",
			initial:     "abc",
			input:       map[string]any{"file_path": "f.txt", "old_string": "zzz", "new_string": "y"},
			wantErr:     true,
		},
		{
			description: "empty old_string is rejected as a validation error before reading the file",
			initial:     "abc",
			input:       map[string]any{"file_path": "f.txt", "old_string": "", "new_string": "y"},
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			work := t.TempDir()
			env := &AgentToolContext{Workdir: work}
			require.NoError(t, os.WriteFile(filepath.Join(work, "f.txt"), []byte(tc.initial), 0o644))

			_, isErr := execEdit(context.Background(), mustJSON(t, tc.input), env)
			require.Equal(t, tc.wantErr, isErr)
			if tc.wantErr {
				return
			}
			data, err := os.ReadFile(filepath.Join(work, "f.txt"))
			require.NoError(t, err)
			require.Equal(t, tc.want, string(data))
		})
	}
}

func TestExecEditRejectsOversizedFile(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(work, "big.txt"), make([]byte, defaultMaxFileBytes+1), 0o644))
	out, isErr := runTool(t, BetaEditTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{
		"file_path": "big.txt", "old_string": "a", "new_string": "b",
	}))
	require.True(t, isErr)
	require.Contains(t, out, "exceeds")
	require.NotContains(t, out, "bash", "the toolset may not include bash, so the error must not point at it")
}

func TestExecEditRejectsDirectory(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(work, "sub"), 0o755))
	out, isErr := runTool(t, BetaEditTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{
		"file_path": "sub", "old_string": "a", "new_string": "b",
	}))
	require.True(t, isErr)
	require.Contains(t, out, "not a regular file")
}

func TestExecEditAllowsNormalFile(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(work, "f.txt"), []byte("alpha\nbeta\ngamma"), 0o644))
	out, isErr := runTool(t, BetaEditTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{
		"file_path": "f.txt", "old_string": "beta", "new_string": "BETA",
	}))
	require.False(t, isErr, "output=%q", out)
	data, err := os.ReadFile(filepath.Join(work, "f.txt"))
	require.NoError(t, err)
	require.Equal(t, "alpha\nBETA\ngamma", string(data))
}

func TestExecEditPreservesCRLF(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(work, "crlf.txt"), []byte("line1\r\nline2\r\nline3\r\n"), 0o644))
	out, isErr := runTool(t, BetaEditTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{
		"file_path": "crlf.txt", "old_string": "line2", "new_string": "LINE2",
	}))
	require.False(t, isErr, "output=%q", out)
	data, err := os.ReadFile(filepath.Join(work, "crlf.txt"))
	require.NoError(t, err)
	require.Equal(t, []byte("line1\r\nLINE2\r\nline3\r\n"), data)
}

func TestExecReadRejectsOversizedFile(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(work, "big.txt"), make([]byte, defaultMaxFileBytes+1), 0o644))
	out, isErr := runTool(t, BetaReadTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{"file_path": "big.txt"}))
	require.True(t, isErr)
	require.Contains(t, out, "exceeds")
	require.Contains(t, out, "view_range", "the error should steer the model to a ranged read")
	require.NotContains(t, out, "bash", "the toolset may not include bash, so the error must not point at it")
}

func TestExecReadRejectsDirectory(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(work, "sub"), 0o755))
	out, isErr := runTool(t, BetaReadTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{"file_path": "sub"}))
	require.True(t, isErr)
	require.Contains(t, out, "not a regular file")
}

// markedFile writes a file whose contents are the unique marker "OLD" followed
// by pad zero bytes, so the edit tool has a unique old_string to replace while
// the file's total size is controllable.
func markedFile(t *testing.T, dir string, pad int) {
	t.Helper()
	content := append([]byte("OLD"), make([]byte, pad)...)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "f.txt"), content, 0o644))
}

func TestExecEditCustomMaxBytesRejectsBelowCap(t *testing.T) {
	work := t.TempDir()
	markedFile(t, work, 2000) // ~2 KiB file
	out, isErr := runTool(t, BetaEditTool(&AgentToolContext{Workdir: work, MaxFileBytes: 1000}), mustJSON(t, map[string]any{
		"file_path": "f.txt", "old_string": "OLD", "new_string": "NEW",
	}))
	require.True(t, isErr)
	require.Contains(t, out, "exceeds")
	require.NotContains(t, out, "bash")
}

func TestExecEditCustomMaxBytesAllowsAboveDefault(t *testing.T) {
	work := t.TempDir()
	markedFile(t, work, defaultMaxFileBytes) // just over the built-in default
	out, isErr := runTool(t, BetaEditTool(&AgentToolContext{Workdir: work, MaxFileBytes: defaultMaxFileBytes * 2}), mustJSON(t, map[string]any{
		"file_path": "f.txt", "old_string": "OLD", "new_string": "NEW",
	}))
	require.False(t, isErr, "output=%q", out)
	data, err := os.ReadFile(filepath.Join(work, "f.txt"))
	require.NoError(t, err)
	require.True(t, len(data) > defaultMaxFileBytes && string(data[:3]) == "NEW")
}

func TestExecEditUncappedAllowsOversized(t *testing.T) {
	work := t.TempDir()
	markedFile(t, work, defaultMaxFileBytes) // would exceed the default cap
	out, isErr := runTool(t, BetaEditTool(&AgentToolContext{Workdir: work, MaxFileBytes: -1}), mustJSON(t, map[string]any{
		"file_path": "f.txt", "old_string": "OLD", "new_string": "NEW",
	}))
	require.False(t, isErr, "output=%q", out)
}

func TestExecEditRejectsDirectoryEvenWhenUncapped(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(work, "sub"), 0o755))
	out, isErr := runTool(t, BetaEditTool(&AgentToolContext{Workdir: work, MaxFileBytes: -1}), mustJSON(t, map[string]any{
		"file_path": "sub", "old_string": "a", "new_string": "b",
	}))
	require.True(t, isErr)
	require.Contains(t, out, "not a regular file")
}

func TestExecReadCustomMaxBytesRejectsBelowCap(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(work, "f.txt"), make([]byte, 2000), 0o644))
	out, isErr := runTool(t, BetaReadTool(&AgentToolContext{Workdir: work, MaxFileBytes: 1000}), mustJSON(t, map[string]any{"file_path": "f.txt"}))
	require.True(t, isErr)
	require.Contains(t, out, "exceeds")
	require.Contains(t, out, "view_range")
	require.NotContains(t, out, "bash")
}

func TestExecReadViewRangeStreamsFileOverCap(t *testing.T) {
	work := t.TempDir()
	// 18 bytes, over the 16-byte cap below.
	require.NoError(t, os.WriteFile(filepath.Join(work, "a.txt"), []byte("line1\nline2\nline3\n"), 0o644))
	// A single 100-byte line followed by a short one.
	require.NoError(t, os.WriteFile(filepath.Join(work, "long.txt"), []byte(strings.Repeat("x", 100)+"\nok\n"), 0o644))
	env := &AgentToolContext{Workdir: work, MaxFileBytes: 16}

	tests := []struct {
		description  string
		input        map[string]any
		want         string
		wantErr      bool
		errContains  string
		errOmitsBash bool
	}{
		{
			description: "a single line of a file over the cap is streamed and returned",
			input:       map[string]any{"file_path": "a.txt", "view_range": []int{2, 2}},
			want:        "line2",
		},
		{
			description: "end_line of 0 streams through to the end of the file and keeps the trailing newline like the whole-file path",
			input:       map[string]any{"file_path": "a.txt", "view_range": []int{2, 0}},
			want:        "line2\nline3\n",
		},
		{
			description:  "a view_range whose selected lines exceed the cap is rejected with advice to narrow it, not to use bash",
			input:        map[string]any{"file_path": "a.txt", "view_range": []int{1, 3}},
			wantErr:      true,
			errContains:  "Narrow the view_range",
			errOmitsBash: true,
		},
		{
			description:  "a single line that alone exceeds the cap is rejected because read cannot return part of a line",
			input:        map[string]any{"file_path": "long.txt", "view_range": []int{1, 1}},
			wantErr:      true,
			errContains:  "cannot return part of a line",
			errOmitsBash: true,
		},
		{
			description: "a short line after an over-cap line is still reachable",
			input:       map[string]any{"file_path": "long.txt", "view_range": []int{2, 2}},
			want:        "ok",
		},
		{
			description: "an inverted view_range on a file over the cap selects nothing rather than erroring",
			input:       map[string]any{"file_path": "a.txt", "view_range": []int{3, 1}},
			want:        "",
		},
		{
			description: "start_line past the end of a file over the cap returns empty content",
			input:       map[string]any{"file_path": "a.txt", "view_range": []int{10, 12}},
			want:        "",
		},
		{
			description: "view_range with the wrong arity is rejected before the size check",
			input:       map[string]any{"file_path": "a.txt", "view_range": []int{2}},
			wantErr:     true,
			errContains: "view_range must be [start_line, end_line]",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			got, isErr := runTool(t, BetaReadTool(env), mustJSON(t, tc.input))
			require.Equal(t, tc.wantErr, isErr, "is_error mismatch; output=%q", got)
			if !tc.wantErr {
				require.Equal(t, tc.want, got)
				return
			}
			require.Contains(t, got, tc.errContains)
			if tc.errOmitsBash {
				require.NotContains(t, got, "bash")
			}
		})
	}
}

func TestExecReadViewRangeStreamsAcrossChunkBoundaries(t *testing.T) {
	work := t.TempDir()
	// ~500 KB, over the default cap and several read chunks long, so chunk
	// boundaries land mid-line.
	var sb strings.Builder
	rows := make([]string, 0, 5000)
	for i := 1; i <= 5000; i++ {
		row := fmt.Sprintf("row%04d", i) + strings.Repeat(".", 90)
		rows = append(rows, row)
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	require.Greater(t, sb.Len(), defaultMaxFileBytes)
	require.Greater(t, sb.Len(), 2*readChunkBytes)
	require.NoError(t, os.WriteFile(filepath.Join(work, "big.txt"), []byte(sb.String()), 0o644))

	got, isErr := runTool(t, BetaReadTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{"file_path": "big.txt", "view_range": []int{1000, 1002}}))
	require.False(t, isErr, "output=%q", got)
	require.Equal(t, strings.Join(rows[999:1002], "\n"), got)
}

func TestExecReadUncappedAllowsOversized(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(work, "big.txt"), make([]byte, defaultMaxFileBytes+1), 0o644))
	_, isErr := runTool(t, BetaReadTool(&AgentToolContext{Workdir: work, MaxFileBytes: -1}), mustJSON(t, map[string]any{"file_path": "big.txt"}))
	require.False(t, isErr)
}

func TestExecReadRejectsDirectoryEvenWhenUncapped(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(work, "sub"), 0o755))
	out, isErr := runTool(t, BetaReadTool(&AgentToolContext{Workdir: work, MaxFileBytes: -1}), mustJSON(t, map[string]any{"file_path": "sub"}))
	require.True(t, isErr)
	require.Contains(t, out, "not a regular file")
}

// TestFileToolsConfinementMatrix drives the full path policy through the real
// tools: the workdir works, a memory folder mounted outside it works, a
// read-only memory folder reads but refuses writes, and everywhere else is
// refused.
func TestFileToolsConfinementMatrix(t *testing.T) {
	tmp := t.TempDir()
	work := filepath.Join(tmp, "work")
	notes := filepath.Join(tmp, "mnt", "notes")
	facts := filepath.Join(tmp, "mnt", "facts")
	outside := filepath.Join(tmp, "outside")
	for _, d := range []string{work, notes, facts, outside} {
		require.NoError(t, os.MkdirAll(d, 0o755))
	}
	require.NoError(t, os.WriteFile(filepath.Join(facts, "facts.md"), []byte("immutable"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("secret"), 0o644))

	env := &AgentToolContext{Workdir: work, AllowedRoots: []string{notes, facts}, ReadOnlyRoots: []string{facts}}
	read, write, edit := BetaReadTool(env), BetaWriteTool(env), BetaEditTool(env)
	const refused = "is outside the session's working directory and its other permitted directories"

	msg, isErr := runTool(t, write, mustJSON(t, map[string]any{"file_path": "w.txt", "content": "in workdir"}))
	require.False(t, isErr, msg)
	got, err := os.ReadFile(filepath.Join(work, "w.txt"))
	require.NoError(t, err)
	require.Equal(t, "in workdir", string(got))

	note := filepath.Join(notes, "note.md")
	msg, isErr = runTool(t, write, mustJSON(t, map[string]any{"file_path": note, "content": "remembered"}))
	require.False(t, isErr, msg)
	msg, isErr = runTool(t, read, mustJSON(t, map[string]any{"file_path": note}))
	require.False(t, isErr, msg)
	require.Equal(t, "remembered", msg)
	msg, isErr = runTool(t, edit, mustJSON(t, map[string]any{"file_path": note, "old_string": "remembered", "new_string": "edited"}))
	require.False(t, isErr, msg)
	got, err = os.ReadFile(note)
	require.NoError(t, err)
	require.Equal(t, "edited", string(got))

	msg, isErr = runTool(t, read, mustJSON(t, map[string]any{"file_path": filepath.Join(facts, "facts.md")}))
	require.False(t, isErr, msg)
	require.Equal(t, "immutable", msg)
	msg, isErr = runTool(t, write, mustJSON(t, map[string]any{"file_path": filepath.Join(facts, "facts.md"), "content": "x"}))
	require.True(t, isErr)
	require.Contains(t, msg, "read-only")
	got, err = os.ReadFile(filepath.Join(facts, "facts.md"))
	require.NoError(t, err)
	require.Equal(t, "immutable", string(got))

	msg, isErr = runTool(t, read, mustJSON(t, map[string]any{"file_path": filepath.Join(outside, "secret.txt")}))
	require.True(t, isErr)
	require.Contains(t, msg, refused)
	msg, isErr = runTool(t, write, mustJSON(t, map[string]any{"file_path": filepath.Join(outside, "new.txt"), "content": "x"}))
	require.True(t, isErr)
	require.Contains(t, msg, refused)
	require.NoFileExists(t, filepath.Join(outside, "new.txt"))
}

func TestWriteAndEditToolsRefuseReadOnlyRoots(t *testing.T) {
	work := t.TempDir()
	ro := filepath.Join(work, "memory", "notes")
	require.NoError(t, os.MkdirAll(ro, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(ro, "facts.md"), []byte("immutable"), 0o644))

	env := &AgentToolContext{Workdir: work, ReadOnlyRoots: []string{ro}}
	write := BetaWriteTool(env)
	edit := BetaEditTool(env)

	// Overwriting an existing file, creating a new one, and editing in place
	// are all refused with an error naming the read-only directory.
	for _, input := range []map[string]any{
		{"file_path": "memory/notes/facts.md", "content": "overwritten"},
		{"file_path": "memory/notes/new.md", "content": "new"},
	} {
		msg, isErr := runTool(t, write, mustJSON(t, input))
		require.True(t, isErr, "write into a read-only root must fail: %v", input)
		require.Contains(t, msg, "read-only")
	}
	msg, isErr := runTool(t, edit, mustJSON(t, map[string]any{
		"file_path": "memory/notes/facts.md", "old_string": "immutable", "new_string": "x",
	}))
	require.True(t, isErr)
	require.Contains(t, msg, "read-only")

	got, err := os.ReadFile(filepath.Join(ro, "facts.md"))
	require.NoError(t, err)
	require.Equal(t, "immutable", string(got))
	require.NoFileExists(t, filepath.Join(ro, "new.md"))

	// Paths outside the read-only root are unaffected.
	msg, isErr = runTool(t, write, mustJSON(t, map[string]any{"file_path": "scratch.txt", "content": "ok"}))
	require.False(t, isErr, msg)
	got, err = os.ReadFile(filepath.Join(work, "scratch.txt"))
	require.NoError(t, err)
	require.Equal(t, "ok", string(got))
}

// TestBinaryMediaTypesTrackGeneratedAPITypes ensures binaryMediaTypes covers
// every media type generated on BetaBase64ImageSourceMediaType and
// BetaBase64PDFSourceParam.
func TestBinaryMediaTypesTrackGeneratedAPITypes(t *testing.T) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filepath.Join("..", "..", "betamessage.go"), nil, 0)
	require.NoError(t, err)

	want := map[string]bool{
		string(anthropic.BetaBase64PDFSourceParam{}.MediaType.Default()): true,
	}
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			ident, ok := vs.Type.(*ast.Ident)
			if !ok || ident.Name != "BetaBase64ImageSourceMediaType" {
				continue
			}
			for _, v := range vs.Values {
				lit, ok := v.(*ast.BasicLit)
				require.True(t, ok && lit.Kind == token.STRING)
				unquoted, err := strconv.Unquote(lit.Value)
				require.NoError(t, err)
				want[unquoted] = true
			}
		}
	}

	got := map[string]bool{}
	for _, mediaType := range binaryMediaTypes {
		got[mediaType] = true
	}
	require.Equal(t, want, got)
}

func TestExecReadImageReturnsImageContentBlock(t *testing.T) {
	cases := []struct {
		name      string
		mediaType anthropic.BetaBase64ImageSourceMediaType
	}{
		{"photo.png", anthropic.BetaBase64ImageSourceMediaTypeImagePNG},
		{"slide.JPG", anthropic.BetaBase64ImageSourceMediaTypeImageJPEG},
		{"slide.jpeg", anthropic.BetaBase64ImageSourceMediaTypeImageJPEG},
		{"anim.gif", anthropic.BetaBase64ImageSourceMediaTypeImageGIF},
		{"hero.webp", anthropic.BetaBase64ImageSourceMediaTypeImageWebP},
	}
	raw := []byte("\x89PNG\r\n\x1a\n\xff\xd8")
	wantB64 := base64.StdEncoding.EncodeToString(raw)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			work := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(work, tc.name), raw, 0o644))
			tool := BetaReadTool(&AgentToolContext{Workdir: work})

			blocks, err := tool.Execute(context.Background(), mustJSON(t, map[string]any{"file_path": tc.name}))
			require.NoError(t, err)
			require.Len(t, blocks, 1)
			require.NotNil(t, blocks[0].OfImage)
			require.NotNil(t, blocks[0].OfImage.Source.OfBase64)
			require.Equal(t, tc.mediaType, blocks[0].OfImage.Source.OfBase64.MediaType)
			require.Equal(t, wantB64, blocks[0].OfImage.Source.OfBase64.Data)
			wire, err := json.Marshal(blocks[0])
			require.NoError(t, err)
			require.JSONEq(t, fmt.Sprintf(`{"type":"image","source":{"type":"base64","media_type":%q,"data":%q}}`, tc.mediaType, wantB64), string(wire))
		})
	}
}

func TestExecReadPDFReturnsDocumentContentBlock(t *testing.T) {
	raw := []byte("%PDF-1.4\n%\xe2\xe3\xcf\xd3\n")
	wantB64 := base64.StdEncoding.EncodeToString(raw)

	for _, name := range []string{"report.pdf", "REPORT.PDF"} {
		t.Run(name, func(t *testing.T) {
			work := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(work, name), raw, 0o644))
			tool := BetaReadTool(&AgentToolContext{Workdir: work})

			blocks, err := tool.Execute(context.Background(), mustJSON(t, map[string]any{"file_path": name}))
			require.NoError(t, err)
			require.Len(t, blocks, 1)
			require.NotNil(t, blocks[0].OfDocument)
			require.NotNil(t, blocks[0].OfDocument.Source.OfBase64)
			require.Equal(t, wantB64, blocks[0].OfDocument.Source.OfBase64.Data)
			wire, err := json.Marshal(blocks[0])
			require.NoError(t, err)
			require.JSONEq(t, fmt.Sprintf(`{"type":"document","source":{"type":"base64","media_type":"application/pdf","data":%q}}`, wantB64), string(wire))
		})
	}
}

func TestExecReadImageExceedingTextCapSucceeds(t *testing.T) {
	work := t.TempDir()
	raw := append([]byte("\x89PNG\r\n\x1a\n"), make([]byte, defaultMaxFileBytes+1024)...)
	require.NoError(t, os.WriteFile(filepath.Join(work, "big.png"), raw, 0o644))
	tool := BetaReadTool(&AgentToolContext{Workdir: work})

	blocks, err := tool.Execute(context.Background(), mustJSON(t, map[string]any{"file_path": "big.png"}))
	require.NoError(t, err)
	require.Len(t, blocks, 1)
	require.NotNil(t, blocks[0].OfImage)
	require.Equal(t, base64.StdEncoding.EncodeToString(raw), blocks[0].OfImage.Source.OfBase64.Data)
}

func TestExecReadImageRejectsOverImageCap(t *testing.T) {
	work := t.TempDir()
	rawCap := (defaultMaxImageBase64Bytes / 4) * 3
	require.NoError(t, os.WriteFile(filepath.Join(work, "huge.png"), make([]byte, rawCap+1), 0o644))
	out, isErr := runTool(t, BetaReadTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{"file_path": "huge.png"}))
	require.True(t, isErr)
	require.Contains(t, out, "exceeds")
	require.Contains(t, out, "limit for image/PDF files")
}

func TestExecReadImageAndPDFCustomMaxBytes(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(work, "small.png"), make([]byte, 300), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(work, "over.png"), make([]byte, 301), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(work, "small.pdf"), make([]byte, 500), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(work, "over.pdf"), make([]byte, 501), 0o644))

	// 400 base64 bytes -> 300 raw bytes.
	imgEnv := &AgentToolContext{Workdir: work, MaxImageBase64Bytes: 400}
	_, err := BetaReadTool(imgEnv).Execute(context.Background(), mustJSON(t, map[string]any{"file_path": "small.png"}))
	require.NoError(t, err)
	out, isErr := runTool(t, BetaReadTool(imgEnv), mustJSON(t, map[string]any{"file_path": "over.png"}))
	require.True(t, isErr)
	require.Contains(t, out, "exceeds 300-byte limit")

	pdfEnv := &AgentToolContext{Workdir: work, MaxPDFBytes: 500}
	_, err = BetaReadTool(pdfEnv).Execute(context.Background(), mustJSON(t, map[string]any{"file_path": "small.pdf"}))
	require.NoError(t, err)
	out, isErr = runTool(t, BetaReadTool(pdfEnv), mustJSON(t, map[string]any{"file_path": "over.pdf"}))
	require.True(t, isErr)
	require.Contains(t, out, "exceeds 500-byte limit")

	// Uncapped media limits allow files above the default caps.
	rawCap := (defaultMaxImageBase64Bytes / 4) * 3
	require.NoError(t, os.WriteFile(filepath.Join(work, "huge.png"), make([]byte, rawCap+1), 0o644))
	_, err = BetaReadTool(&AgentToolContext{Workdir: work, MaxImageBase64Bytes: -1}).Execute(
		context.Background(), mustJSON(t, map[string]any{"file_path": "huge.png"}))
	require.NoError(t, err)

	// An explicit positive MaxFileBytes tightens the media cap.
	tightEnv := &AgentToolContext{Workdir: work, MaxFileBytes: 200, MaxImageBase64Bytes: -1}
	out, isErr = runTool(t, BetaReadTool(tightEnv), mustJSON(t, map[string]any{"file_path": "small.png"}))
	require.True(t, isErr)
	require.Contains(t, out, "exceeds 200-byte limit")
}

func TestExecReadImageRejectsViewRange(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(work, "pic.png"), []byte("\x89PNG\r\n\x1a\n"), 0o644))
	out, isErr := runTool(t, BetaReadTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{
		"file_path": "pic.png", "view_range": []int{1, 10},
	}))
	require.True(t, isErr)
	require.Contains(t, out, "view_range is not supported for image/PDF files")
}

func TestExecReadNonImageBinaryRaisesToolError(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(work, "blob.bin"), []byte("\xff\xd8\xff\xe0\x00\x10\nvalid\n"), 0o644))

	out, isErr := runTool(t, BetaReadTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{"file_path": "blob.bin"}))
	require.True(t, isErr)
	require.Equal(t, "read: blob.bin: not valid UTF-8 text (binary files are only supported for image/PDF extensions)", out)

	// Streaming ranged read over MaxFileBytes also rejects invalid UTF-8 bytes.
	out, isErr = runTool(t, BetaReadTool(&AgentToolContext{Workdir: work, MaxFileBytes: 8}), mustJSON(t, map[string]any{
		"file_path": "blob.bin", "view_range": []int{1, 1},
	}))
	require.True(t, isErr)
	require.Equal(t, "read: blob.bin: not valid UTF-8 text (binary files are only supported for image/PDF extensions)", out)
}

func TestExecEditBinaryFileRaisesToolError(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(work, "pic.png"), []byte("\x89PNG\r\n\x1a\n\xff\xfe"), 0o644))
	out, isErr := runTool(t, BetaEditTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{
		"file_path": "pic.png", "old_string": "PNG", "new_string": "JPG",
	}))
	require.True(t, isErr)
	require.Equal(t, "edit: pic.png: not valid UTF-8 text (cannot edit binary files)", out)
}
