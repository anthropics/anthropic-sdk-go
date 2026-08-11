package agenttoolset

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

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
			got, isErr := execRead(context.Background(), mustJSON(t, tc.input), env)
			require.Equal(t, tc.wantErr, isErr, "is_error mismatch; output=%q", got)
			if !tc.wantErr {
				require.Equal(t, tc.want, got)
			}
		})
	}
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

func TestExecReadRejectsOversizedFile(t *testing.T) {
	work := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(work, "big.txt"), make([]byte, defaultMaxFileBytes+1), 0o644))
	out, isErr := runTool(t, BetaReadTool(&AgentToolContext{Workdir: work}), mustJSON(t, map[string]any{"file_path": "big.txt"}))
	require.True(t, isErr)
	require.Contains(t, out, "exceeds")
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
