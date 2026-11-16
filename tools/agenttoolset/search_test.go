package agenttoolset

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
		wantErr     string
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
			wantErr:     "pattern is required",
		},
		{
			description: "absolute pattern is always rejected",
			pattern:     "/etc/*",
			wantErr:     "pass a relative pattern",
		},
		{
			description: "pattern containing a .. segment is rejected so it cannot escape the workdir",
			pattern:     "../*.go",
			wantErr:     "must not contain",
		},
		{
			description: "doublestar pattern with an embedded .. segment is also rejected",
			pattern:     "**/../*.go",
			wantErr:     "must not contain",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out, isErr := execGlob(context.Background(), mustJSON(t, map[string]any{"pattern": tc.pattern}), env)
			require.Equal(t, tc.wantErr != "", isErr, "output=%q", out)
			if tc.wantErr != "" {
				require.Contains(t, out, tc.wantErr)
				return
			}
			if tc.assert == nil {
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

	t.Run("a symlinked workdir is searched through its canonical path by default", func(t *testing.T) {
		link := filepath.Join(t.TempDir(), "link")
		require.NoError(t, os.Symlink(work, link))
		out, isErr := execGlob(context.Background(), mustJSON(t, map[string]any{"pattern": "**/*.go"}), &AgentToolContext{Workdir: link})
		require.False(t, isErr, "output=%q", out)
		lines := strings.Split(strings.TrimSpace(out), "\n")
		require.Len(t, lines, 3)
		for _, l := range lines {
			require.True(t, strings.HasPrefix(l, canonicalRoot(work)+string(filepath.Separator)), "result %q is not under the canonical workdir", l)
		}
	})
}

// TestGlobAndGrepReachAnAllowedRoot: a memory folder mounted outside the
// workdir is searchable when listed in AllowedRoots, and a directory that is
// in neither is refused by both tools.
func TestGlobAndGrepReachAnAllowedRoot(t *testing.T) {
	tmp := t.TempDir()
	work := filepath.Join(tmp, "work")
	mount := filepath.Join(tmp, "mnt", "notes")
	outside := filepath.Join(tmp, "outside")
	for _, d := range []string{work, mount, outside} {
		require.NoError(t, os.MkdirAll(d, 0o755))
	}
	require.NoError(t, os.WriteFile(filepath.Join(mount, "note.md"), []byte("remembered fact"), 0o644))
	env := &AgentToolContext{Workdir: work, AllowedRoots: []string{mount}}
	const refused = "is outside the session's working directory and its other permitted directories"

	out, isErr := execGlob(context.Background(), mustJSON(t, map[string]any{"pattern": "*.md", "path": mount}), env)
	require.False(t, isErr, "output=%q", out)
	require.Equal(t, filepath.Join(canonicalRoot(mount), "note.md"), out)

	for _, tc := range []struct {
		description string
		hidePath    bool
	}{
		{description: "grep via ripgrep when it is on PATH"},
		{description: "grep via the built-in walker", hidePath: true},
	} {
		t.Run(tc.description, func(t *testing.T) {
			if tc.hidePath {
				t.Setenv("PATH", "")
			}
			out, isErr := execGrep(context.Background(), mustJSON(t, map[string]any{"pattern": "remembered", "path": mount}), env)
			require.False(t, isErr, "output=%q", out)
			require.Contains(t, out, "note.md")
		})
	}

	out, isErr = execGlob(context.Background(), mustJSON(t, map[string]any{"pattern": "*", "path": outside}), env)
	require.True(t, isErr)
	require.Contains(t, out, refused)
	out, isErr = execGrep(context.Background(), mustJSON(t, map[string]any{"pattern": "x", "path": outside}), env)
	require.True(t, isErr)
	require.Contains(t, out, refused)
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

func TestExecGlobToleratesSymlinkLoopsInWorkdir(t *testing.T) {
	work, _ := symlinkLoopFixture(t)
	require.NoError(t, os.WriteFile(filepath.Join(work, "ok.txt"), []byte("x"), 0o644))

	out, isErr := execGlob(context.Background(), mustJSON(t, map[string]any{"pattern": "*"}), &AgentToolContext{Workdir: work})
	require.False(t, isErr, "output=%q", out)
	require.Contains(t, strings.Split(out, "\n"), filepath.Join(work, "ok.txt"))
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

func TestExecGrepIgnoresRipgrepConfig(t *testing.T) {
	if _, err := exec.LookPath("rg"); err != nil {
		t.Skip("ripgrep is not on PATH")
	}
	work := t.TempDir()
	outside := t.TempDir()
	config := filepath.Join(outside, "ripgreprc")
	require.NoError(t, os.WriteFile(config, []byte("--follow\n"), 0o644))
	t.Setenv("RIPGREP_CONFIG_PATH", config)
	require.NoError(t, os.WriteFile(filepath.Join(work, "notes.txt"), []byte("needle"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("needle"), 0o600))
	require.NoError(t, os.Symlink(filepath.Join(outside, "secret.txt"), filepath.Join(work, "leak")))

	out, isErr := execGrep(context.Background(), mustJSON(t, map[string]any{"pattern": "needle"}), &AgentToolContext{Workdir: work})
	require.False(t, isErr, "output=%q", out)
	require.Contains(t, out, "notes.txt")
	require.NotContains(t, out, "leak")
}

func TestExecGrepPassesNoConfigToRipgrep(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the stand-in rg is a shell script")
	}
	bin := t.TempDir()
	args := filepath.Join(bin, "args")
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" >\"$RG_ARGS\"\nexit 1\n"
	require.NoError(t, os.WriteFile(filepath.Join(bin, "rg"), []byte(script), 0o755))
	t.Setenv("PATH", bin)
	t.Setenv("RG_ARGS", args)

	out, isErr := execGrep(context.Background(), mustJSON(t, map[string]any{"pattern": "needle"}), &AgentToolContext{Workdir: t.TempDir()})
	require.False(t, isErr, "output=%q", out)
	require.Equal(t, "no matches", out)

	got, err := os.ReadFile(args)
	require.NoError(t, err)
	require.Contains(t, strings.Split(string(got), "\n"), "--no-config")
}

func TestGrepRunsOnlyARipgrepOutsideThePermittedRoots(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the stand-in rg is a shell script")
	}
	host := t.TempDir()
	at := func(rel string) string { return filepath.Join(host, rel) }
	for _, dir := range []string{"work/bin", "WORK/bin", "allowed/bin", "caf\u00e9", "cafe\u0301/bin", "system", "linked-file", "linked-batch"} {
		require.NoError(t, os.MkdirAll(at(dir), 0o755))
	}
	for _, file := range []string{"work/bin/rg", "WORK/bin/rg", "allowed/bin/rg", "cafe\u0301/bin/rg", "system/rg", "system/rg.CMD"} {
		require.NoError(t, os.WriteFile(at(file), []byte("#!/bin/sh\necho \"$0\"\n"), 0o755))
	}
	require.NoError(t, os.Symlink(at("work/bin"), at("linked-dir")))
	require.NoError(t, os.Symlink(at("system"), at("linked-system")))
	require.NoError(t, os.Symlink(at("system/rg.CMD"), at("linked-batch/rg")))
	require.NoError(t, os.Symlink(at("work/bin/rg"), at("linked-file/rg")))
	require.NoError(t, os.Symlink("loop", at("loop")))
	require.NoError(t, os.WriteFile(at("work/a.txt"), []byte("alpha\nbeta\n"), 0o644))
	t.Chdir(at("system"))

	ripgrep := filepath.Join(canonicalRoot(at("system")), "rg") + "\n"
	builtIn := filepath.Join(canonicalRoot(at("work")), "a.txt") + ":2:beta"

	for _, tc := range []struct {
		description  string
		path         []string
		allowedRoots []string
		want         string
	}{
		{description: "rg outside the workdir and the allowed roots", path: []string{at("system")}, want: ripgrep},
		{description: "rg in the current directory with an empty allowed root", path: []string{at("system")}, allowedRoots: []string{""}, want: ripgrep},
		{description: "rg outside behind a symlink is started by its real path", path: []string{at("linked-system")}, want: ripgrep},
		{description: "rg in the workdir", path: []string{at("work/bin")}, want: builtIn},
		{description: "rg in the workdir spelled in another letter case", path: []string{at("WORK/bin")}, want: builtIn},
		{description: "rg in an allowed root", path: []string{at("allowed/bin")}, want: builtIn},
		{description: "rg in an allowed root spelled in another Unicode form", path: []string{at("cafe\u0301/bin")}, allowedRoots: []string{at("caf\u00e9")}, want: builtIn},
		{description: "rg outside after that one", path: []string{at("cafe\u0301/bin"), at("system")}, allowedRoots: []string{at("caf\u00e9")}, want: ripgrep},
		{description: "PATH entry that is a symlink into the workdir", path: []string{at("linked-dir")}, want: builtIn},
		{description: "rg that is a symlink into the workdir", path: []string{at("linked-file")}, want: builtIn},
		{description: "rg that is a symlink to a batch file", path: []string{at("linked-batch")}, want: builtIn},
		{description: "relative PATH entry", path: []string{"../system"}, want: builtIn},
		{
			description: "rg outside after the ones that are passed over",
			path:        []string{at("work/bin"), at("allowed/bin"), at("linked-dir"), at("linked-file"), at("linked-batch"), "../system", at("system")},
			want:        ripgrep,
		},
		{description: "allowed root that cannot be resolved", path: []string{at("system")}, allowedRoots: []string{at("loop")}, want: builtIn},
	} {
		t.Run(tc.description, func(t *testing.T) {
			t.Setenv("PATH", strings.Join(tc.path, string(filepath.ListSeparator)))
			env := &AgentToolContext{Workdir: at("work"), AllowedRoots: append([]string{at("allowed")}, tc.allowedRoots...)}
			out, isErr := runTool(t, BetaGrepTool(env), mustJSON(t, map[string]any{"pattern": "beta"}))
			require.False(t, isErr, "output=%q", out)
			require.Equal(t, tc.want, out)
		})
	}
}

func TestGrepStartsRipgrepUnderTheNameRg(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the stand-in rg is a link to /bin/sh")
	}
	bin := t.TempDir()
	require.NoError(t, os.Symlink("/bin/sh", filepath.Join(bin, "rg")))
	t.Setenv("PATH", bin)

	out, isErr := runTool(t, BetaGrepTool(&AgentToolContext{Workdir: t.TempDir()}), mustJSON(t, map[string]any{"pattern": "needle"}))
	require.True(t, isErr, "output=%q", out)
	require.Regexp(t, `^grep: rg failed: exit status \d+: rg: `, out)
}

func TestGrepKeepsCredentialsFromRipgrep(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the stand-in rg is a shell script")
	}
	bin := t.TempDir()
	script := "#!/bin/sh\necho \"key=[$ANTHROPIC_API_KEY] mixed=[$Anthropic_Fake] kept=[$KEPT]\"\n"
	require.NoError(t, os.WriteFile(filepath.Join(bin, "rg"), []byte(script), 0o755))
	t.Setenv("PATH", bin)
	t.Setenv("ANTHROPIC_API_KEY", "sk-ant-should-not-leak")
	t.Setenv("Anthropic_Fake", "should-not-leak")
	t.Setenv("KEPT", "kept")

	out, isErr := runTool(t, BetaGrepTool(&AgentToolContext{Workdir: t.TempDir()}), mustJSON(t, map[string]any{"pattern": "needle"}))
	require.False(t, isErr, "output=%q", out)
	require.Equal(t, "key=[] mixed=[] kept=[kept]\n", out)
}

func TestGrepBuiltInSearchSniffsAFileBeforeReadingAllOfIt(t *testing.T) {
	f, err := os.Open("/proc/self/pagemap")
	if err != nil {
		t.Skip("needs a readable /proc/self/pagemap")
	}
	require.NoError(t, f.Close())
	t.Setenv("PATH", "")
	out, isErr := runTool(t, BetaGrepTool(&AgentToolContext{Workdir: "/"}), mustJSON(t, map[string]any{"pattern": "needle", "path": "/proc/self/pagemap"}))
	require.False(t, isErr, "output=%q", out)
	require.Equal(t, "no matches", out)
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
