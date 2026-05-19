package agenttoolset

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExecGlob(t *testing.T) {
	work := t.TempDir()
	env := &AgentToolContext{Workdir: work}

	write := func(rel string, mtime time.Time) {
		full := filepath.Join(work, rel)
		require.NoError(t, os.MkdirAll(filepath.Dir(full), 0o755))
		require.NoError(t, os.WriteFile(full, []byte("x"), 0o644))
		require.NoError(t, os.Chtimes(full, mtime, mtime))
	}
	now := time.Now()
	write("a.go", now.Add(-2*time.Hour))
	write("b.go", now)
	write("sub/c.go", now.Add(-1*time.Hour))
	write("d.txt", now)

	tests := []struct {
		description string
		pattern     string
		assert      func(t *testing.T, lines []string)
		wantErr     bool
	}{
		{
			description: "doublestar pattern matches across directories and excludes non-matching extensions",
			pattern:     "**/*.go",
			assert: func(t *testing.T, lines []string) {
				require.Len(t, lines, 3)
				for _, l := range lines {
					require.True(t, strings.HasSuffix(l, ".go"), "non-go result leaked: %q", l)
				}
			},
		},
		{
			description: "results are ordered newest-mtime first so the model sees the most recently touched files",
			pattern:     "**/*.go",
			assert: func(t *testing.T, lines []string) {
				require.True(t, strings.HasSuffix(lines[0], "b.go"), "expected b.go first, got %q", lines[0])
				require.True(t, strings.HasSuffix(lines[len(lines)-1], "a.go"), "expected a.go last, got %q", lines[len(lines)-1])
			},
		},
		{
			description: "empty pattern is rejected before walking the filesystem",
			pattern:     "",
			wantErr:     true,
		},
		{
			description: "absolute pattern is rejected when UnrestrictedPaths is false",
			pattern:     "/etc/*",
			wantErr:     true,
		},
		{
			description: "pattern containing a .. segment is rejected so it cannot escape the workdir",
			pattern:     "../*.go",
			wantErr:     true,
		},
		{
			description: "doublestar pattern with an embedded .. segment is also rejected",
			pattern:     "**/../*.go",
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out, isErr := execGlob(context.Background(), mustJSON(t, map[string]any{"pattern": tc.pattern}), env)
			require.Equal(t, tc.wantErr, isErr, "output=%q", out)
			if tc.wantErr || tc.assert == nil {
				return
			}
			tc.assert(t, strings.Split(strings.TrimSpace(out), "\n"))
		})
	}

	t.Run("pattern with no matches returns the sentinel string rather than empty output", func(t *testing.T) {
		out, isErr := execGlob(context.Background(), mustJSON(t, map[string]any{"pattern": "*.nomatch"}), env)
		require.False(t, isErr)
		require.Equal(t, "no matches", out)
	})

	t.Run("path argument scopes the search to a subdirectory of the workdir", func(t *testing.T) {
		out, isErr := execGlob(context.Background(), mustJSON(t, map[string]any{"pattern": "*.go", "path": "sub"}), env)
		require.False(t, isErr)
		require.Contains(t, out, "c.go")
		require.NotContains(t, out, "a.go")
	})
}

func TestExecGrep(t *testing.T) {
	work := t.TempDir()
	env := &AgentToolContext{Workdir: work}
	require.NoError(t, os.WriteFile(filepath.Join(work, "a.txt"), []byte("hello world\nfoo bar\nhello again"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(work, "b.txt"), []byte("nothing here"), 0o644))

	tests := []struct {
		description string
		pattern     string
		hidePath    bool
		assert      func(t *testing.T, out string)
		wantErr     bool
	}{
		{
			description: "literal pattern matching multiple lines reports each hit with file:line: prefix",
			pattern:     "hello",
			assert: func(t *testing.T, out string) {
				require.Contains(t, out, "a.txt:1:")
				require.Contains(t, out, "a.txt:3:")
				require.NotContains(t, out, "b.txt")
			},
		},
		{
			description: "no-match pattern returns the sentinel rather than an error so the model can continue",
			pattern:     "zzz_absent",
			assert: func(t *testing.T, out string) {
				require.Equal(t, "no matches", out)
			},
		},
		{
			description: "empty pattern is rejected as a validation error",
			pattern:     "",
			wantErr:     true,
		},
		{
			description: "fallback walker handles the same query when ripgrep is absent from PATH",
			pattern:     "hello",
			hidePath:    true,
			assert: func(t *testing.T, out string) {
				require.Contains(t, out, "a.txt:1:")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			if tc.hidePath {
				t.Setenv("PATH", "")
			}
			out, isErr := execGrep(context.Background(), mustJSON(t, map[string]any{"pattern": tc.pattern}), env)
			require.Equal(t, tc.wantErr, isErr, "output=%q", out)
			if tc.assert != nil {
				tc.assert(t, out)
			}
		})
	}
}

// TestExecGrepSkipsSymlinks verifies the built-in walker never reads through a
// symlink: a link inside the workdir pointing at a file outside it (e.g. a
// stand-in for /etc/shadow) must not have its target's contents surface in
// grep output. PATH is cleared so the walker is exercised rather than ripgrep.
func TestExecGrepSkipsSymlinks(t *testing.T) {
	t.Setenv("PATH", "")
	work := t.TempDir()
	outside := t.TempDir()
	env := &AgentToolContext{Workdir: work}

	require.NoError(t, os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("TOPSECRET_NEEDLE"), 0o600))
	require.NoError(t, os.Symlink(filepath.Join(outside, "secret.txt"), filepath.Join(work, "leak")))
	require.NoError(t, os.WriteFile(filepath.Join(work, "real.txt"), []byte("ordinary line"), 0o644))

	out, isErr := execGrep(context.Background(), mustJSON(t, map[string]any{"pattern": "TOPSECRET_NEEDLE"}), env)
	require.False(t, isErr, "output=%q", out)
	require.Equal(t, "no matches", out,
		"grep must not follow a symlink out of the workdir; %q leaked", out)
}

// TestGlobMatchDeepDoublestarBacktracking checks that a pattern alternating
// "**" and a literal segment matches a deeply nested path without exponential
// backtracking. Without memoization, ~30 alternating "**"/"x" segments against
// a 30-deep "x/x/.../x" path is ~2^30 recursive calls; with memoization it is
// O(n^2). Bound the wall time so a regression fails fast rather than hanging
// the suite.
func TestGlobMatchDeepDoublestarBacktracking(t *testing.T) {
	const depth = 30
	patSegs := make([]string, 0, depth*2)
	nameSegs := make([]string, 0, depth)
	for i := 0; i < depth; i++ {
		patSegs = append(patSegs, "**", "x")
		nameSegs = append(nameSegs, "x")
	}
	pat := strings.Join(patSegs, "/")
	name := strings.Join(nameSegs, "/")

	start := time.Now()
	require.True(t, globMatch(pat, name), "pattern %q should match %q", pat, name)
	// A non-match (one extra trailing literal that the path lacks) is the
	// worst case for a backtracking matcher: it explores every failed branch.
	require.False(t, globMatch(pat+"/y", name))
	elapsed := time.Since(start)
	require.Less(t, elapsed, 5*time.Second,
		"globMatch took %v for a %d-deep ** pattern; backtracking is unbounded", elapsed, depth)
}
