package agenttoolset

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
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
			description: "inverted view_range is rejected instead of panicking the runner with a slice-bounds error",
			input:       map[string]any{"file_path": "a.txt", "view_range": []int{3, 1}},
			wantErr:     true,
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
	require.NoError(t, os.WriteFile(filepath.Join(work, "big.txt"), make([]byte, editMaxBytes+1), 0o644))
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
	require.NoError(t, os.WriteFile(filepath.Join(work, "big.txt"), make([]byte, readMaxBytes+1), 0o644))
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
