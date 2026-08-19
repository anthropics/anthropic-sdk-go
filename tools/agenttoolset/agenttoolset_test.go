package agenttoolset

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/stretchr/testify/require"
)

func TestResolvePath(t *testing.T) {
	work := t.TempDir()
	// Canonicalise the workdir the way resolvePath does so comparisons hold on
	// platforms where the temp dir lives behind a symlink (e.g. /var on macOS).
	absWork := canonicalRoot(work)

	tests := []struct {
		description string
		env         *AgentToolContext
		input       string
		want        string
		wantErr     string
	}{
		{
			description: "plain relative file resolves under the workdir",
			env:         &AgentToolContext{Workdir: work},
			input:       "a.txt",
			want:        filepath.Join(absWork, "a.txt"),
		},
		{
			description: "nested relative path resolves with intermediate directories joined",
			env:         &AgentToolContext{Workdir: work},
			input:       filepath.Join("sub", "b.txt"),
			want:        filepath.Join(absWork, "sub", "b.txt"),
		},
		{
			description: "dot-dot that climbs out of the workdir is rejected to keep tools jailed by default",
			env:         &AgentToolContext{Workdir: work},
			input:       filepath.Join("..", "escape.txt"),
			wantErr:     `path "../escape.txt" is outside the session's working directory`,
		},
		{
			description: "dot-dot that stays inside the workdir after normalisation is permitted",
			env:         &AgentToolContext{Workdir: work},
			input:       filepath.Join("sub", "..", "c.txt"),
			want:        filepath.Join(absWork, "c.txt"),
		},
		{
			description: "absolute path outside workdir is rejected",
			env:         &AgentToolContext{Workdir: work},
			input:       "/etc/passwd",
			wantErr:     `path "/etc/passwd" is outside the session's working directory`,
		},
		{
			description: "absolute path inside workdir is permitted",
			env:         &AgentToolContext{Workdir: work},
			input:       filepath.Join(absWork, "a.txt"),
			want:        filepath.Join(absWork, "a.txt"),
		},
		{
			description: "absolute path naming the workdir itself is permitted",
			env:         &AgentToolContext{Workdir: work},
			input:       absWork,
			want:        absWork,
		},
		{
			description: "absolute sibling that string-prefixes the workdir is rejected (segment-aware contain)",
			env:         &AgentToolContext{Workdir: work},
			input:       filepath.Join(absWork+"extra", "secret.txt"),
			wantErr:     "is outside the session's working directory",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			got, err := resolvePath(tc.env, tc.input)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

// TestResolvePathAllowedRoots covers the second tier of permitted roots: a
// path under an AllowedRoots entry resolves, everything outside Workdir and
// every entry is refused, and the refusal only mentions "other permitted
// directories" when a non-empty entry exists.
func TestResolvePathAllowedRoots(t *testing.T) {
	tmp := canonicalRoot(t.TempDir())
	work := filepath.Join(tmp, "work")
	mount := filepath.Join(tmp, "mount")
	for _, d := range []string{work, mount, mount + "2", filepath.Join(tmp, "cwd")} {
		require.NoError(t, os.Mkdir(d, 0o755))
	}
	require.NoError(t, os.WriteFile(filepath.Join(mount, "one.md"), []byte("1"), 0o644))

	const (
		short = "is outside the session's working directory"
		long  = "is outside the session's working directory and its other permitted directories"
	)
	tests := []struct {
		description  string
		allowedRoots []string
		input        string
		want         string
		wantErr      string
	}{
		{
			description:  "absolute path under an allowed root resolves to its canonical form",
			allowedRoots: []string{mount},
			input:        filepath.Join(mount, "note.md"),
			want:         filepath.Join(mount, "note.md"),
		},
		{
			description:  "relative path still resolves under the workdir when allowed roots are set",
			allowedRoots: []string{mount},
			input:        "inside.txt",
			want:         filepath.Join(work, "inside.txt"),
		},
		{
			description:  "absolute path outside the workdir and every allowed root is refused with the long message",
			allowedRoots: []string{mount},
			input:        filepath.Join(tmp, "elsewhere", "x"),
			wantErr:      long,
		},
		{
			description:  "relative escape outside every root is refused with the long message",
			allowedRoots: []string{mount},
			input:        filepath.Join("..", "elsewhere", "x"),
			wantErr:      long,
		},
		{
			description:  "sibling that string-prefixes an allowed root is refused (segment-aware for allowed roots too)",
			allowedRoots: []string{mount},
			input:        filepath.Join(mount+"2", "x"),
			wantErr:      long,
		},
		{
			description:  "an allowed root that is a regular file permits exactly that file",
			allowedRoots: []string{filepath.Join(mount, "one.md")},
			input:        filepath.Join(mount, "one.md"),
			want:         filepath.Join(mount, "one.md"),
		},
		{
			description:  "an allowed root that is a regular file does not permit its siblings",
			allowedRoots: []string{filepath.Join(mount, "one.md")},
			input:        filepath.Join(mount, "two.md"),
			wantErr:      long,
		},
		{
			description:  "the filesystem root as an allowed root permits everything",
			allowedRoots: []string{"/"},
			input:        filepath.Join(tmp, "elsewhere", "x"),
			want:         filepath.Join(tmp, "elsewhere", "x"),
		},
		{
			description:  "an empty AllowedRoots list refuses with the short message",
			allowedRoots: []string{},
			input:        "/etc/passwd",
			wantErr:      `path "/etc/passwd" ` + short,
		},
		{
			description:  "an empty entry grants nothing and does not count as another permitted directory",
			allowedRoots: []string{""},
			input:        filepath.Join(tmp, "cwd", "x"),
			wantErr:      `" ` + short,
		},
	}

	// An empty entry must not silently mean the process working directory.
	t.Chdir(filepath.Join(tmp, "cwd"))
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			env := &AgentToolContext{Workdir: work, AllowedRoots: tc.allowedRoots}
			got, err := resolvePath(env, tc.input)
			if tc.wantErr != "" {
				require.ErrorContains(t, err, tc.wantErr)
				if tc.wantErr != long {
					require.NotContains(t, err.Error(), "other permitted directories")
				}
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestResolvePathRejectsUnrestrictedPaths(t *testing.T) {
	env := &AgentToolContext{Workdir: t.TempDir(), UnrestrictedPaths: true}
	_, err := resolvePath(env, "a.txt")
	require.ErrorIs(t, err, ErrUnrestrictedPathsUnsupported)
}

// TestResolvePathSymlinkOutOfAllowedRootIsRejected: symlinks are canonicalised
// before the containment check, so a link inside an allowed root cannot reach
// outside it.
func TestResolvePathSymlinkOutOfAllowedRootIsRejected(t *testing.T) {
	work := t.TempDir()
	mount := t.TempDir()
	require.NoError(t, os.Symlink("/etc/passwd", filepath.Join(mount, "leak")))
	env := &AgentToolContext{Workdir: work, AllowedRoots: []string{mount}}
	_, err := resolvePath(env, filepath.Join(mount, "leak"))
	require.ErrorContains(t, err, "is outside the session's working directory and its other permitted directories")
}

// TestASymlinkedReadOnlyRootStillRefusesWrites: AllowedRoots entries resolve
// at check time, so ReadOnlyRoots must too — the same symlinked path must not
// grant access on one side while its write protection misses on the other.
func TestASymlinkedReadOnlyRootStillRefusesWrites(t *testing.T) {
	tmp := canonicalRoot(t.TempDir())
	work := filepath.Join(tmp, "work")
	realStore := filepath.Join(tmp, "real-store")
	link := filepath.Join(tmp, "link-store")
	require.NoError(t, os.Mkdir(work, 0o755))
	require.NoError(t, os.Mkdir(realStore, 0o755))
	require.NoError(t, os.Symlink(realStore, link))

	for _, readOnly := range []string{link, realStore} {
		env := &AgentToolContext{Workdir: work, AllowedRoots: []string{link}, ReadOnlyRoots: []string{readOnly}}
		got, err := resolvePath(env, filepath.Join(link, "note.md"))
		require.NoError(t, err)
		require.Equal(t, filepath.Join(realStore, "note.md"), got)

		_, err = resolveWritablePath(env, filepath.Join(link, "note.md"))
		require.EqualError(t, err, filepath.Join(link, "note.md")+" is inside read-only directory "+readOnly)
	}
}

// Every tool built from a context still carrying the deprecated flag refuses
// to run, so the misconfiguration surfaces on the first call rather than as a
// silently different confinement policy.
func TestToolsRefuseToRunWithUnrestrictedPaths(t *testing.T) {
	env := &AgentToolContext{Workdir: t.TempDir(), UnrestrictedPaths: true}
	tools := BetaAgentToolset20260401(env)
	t.Cleanup(func() { CloseAll(tools) })
	require.Len(t, tools, 6)
	for _, tool := range tools {
		t.Run(tool.Name(), func(t *testing.T) {
			_, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
			require.ErrorIs(t, err, ErrUnrestrictedPathsUnsupported)
			var toolErr *ToolError
			require.False(t, errors.As(err, &toolErr), "a configuration error, not a tool refusal shown to the model as its own mistake")
		})
	}
}

// TestResolvePathConfinesSymlinks verifies the real (non-lexical) confinement:
// a symlink that lives inside the workdir but points outside it is resolved
// before the workdir check, so the operation is rejected — even when the link
// target does not exist (dangling).
func TestResolvePathConfinesSymlinks(t *testing.T) {
	work := t.TempDir()
	outside := t.TempDir()
	env := &AgentToolContext{Workdir: work}

	// Existing target outside the workdir.
	require.NoError(t, os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("x"), 0o644))
	require.NoError(t, os.Symlink(filepath.Join(outside, "secret.txt"), filepath.Join(work, "live")))
	_, err := resolvePath(env, "live")
	require.Error(t, err, "a symlink inside the workdir that points outside it must be rejected")

	// Dangling target outside the workdir.
	require.NoError(t, os.Symlink(filepath.Join(outside, "nope.txt"), filepath.Join(work, "dangling")))
	_, err = resolvePath(env, "dangling")
	require.Error(t, err, "a dangling symlink inside the workdir that points outside it must be rejected")

	// A symlink whose resolved target stays inside the workdir is fine.
	require.NoError(t, os.Mkdir(filepath.Join(work, "sub"), 0o755))
	require.NoError(t, os.Symlink(filepath.Join(work, "sub"), filepath.Join(work, "inside")))
	got, err := resolvePath(env, filepath.Join("inside", "f.txt"))
	require.NoError(t, err)
	require.Equal(t, filepath.Join(canonicalRoot(work), "sub", "f.txt"), got)

	// A dangling symlink whose target is inside the workdir resolves to that
	// target even when the target's parent directory does not exist yet.
	require.NoError(t, os.Symlink(filepath.Join("newdir", "f.txt"), filepath.Join(work, "d")))
	got, err = resolvePath(env, "d")
	require.NoError(t, err)
	require.Equal(t, filepath.Join(canonicalRoot(work), "newdir", "f.txt"), got)
}

// symlinkLoopFixture lays out, under a fresh workdir: loop_a <-> loop_b,
// self -> self, evil_link -> <outside>/secret.txt, and L -> loop_a/../evil_link.
// It returns the workdir and the outside directory.
func symlinkLoopFixture(t *testing.T) (work, outside string) {
	t.Helper()
	work = t.TempDir()
	outside = t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("SECRET"), 0o644))
	for name, target := range map[string]string{
		"loop_a":    "loop_b",
		"loop_b":    "loop_a",
		"self":      "self",
		"evil_link": filepath.Join(outside, "secret.txt"),
		"L":         "loop_a/../evil_link",
	} {
		require.NoError(t, os.Symlink(target, filepath.Join(work, name)))
	}
	return work, outside
}

func TestResolvePathRejectsSymlinkLoops(t *testing.T) {
	work, outside := symlinkLoopFixture(t)
	env := &AgentToolContext{Workdir: work}
	absWork := canonicalRoot(work)

	tests := []struct {
		description string
		input       string
		wantErr     string
	}{
		{"two-link cycle", "loop_a", `path "loop_a": too many levels of symbolic links`},
		{"child under a two-link cycle", "loop_a/child.txt", `path "loop_a/child.txt": too many levels of symbolic links`},
		{"self-referencing link", "self", `path "self": too many levels of symbolic links`},
		{"child under a self-referencing link", "self/x", `path "self/x": too many levels of symbolic links`},
		{"dot-dot past a cycle is collapsed lexically and lands on the escaping link", "loop_a/../evil_link", `path "loop_a/../evil_link" is outside the session's working directory`},
		{"dot-dot past a self link is collapsed lexically and lands on the escaping link", "self/../evil_link", `path "self/../evil_link" is outside the session's working directory`},
		{"dangling link pointing outside the workdir", "dangle_out", `path "dangle_out" is outside the session's working directory`},
	}
	require.NoError(t, os.Symlink(filepath.Join(outside, "nope"), filepath.Join(work, "dangle_out")))

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			errc := make(chan error, 1)
			go func() {
				_, err := resolvePath(env, tc.input)
				errc <- err
			}()
			select {
			case err := <-errc:
				require.EqualError(t, err, tc.wantErr)
				require.NotContains(t, err.Error(), absWork)
			case <-time.After(2 * time.Second):
				t.Fatal("resolvePath did not return within 2s")
			}
		})
	}

	t.Run("link whose target text crosses a cycle is denied one way or the other", func(t *testing.T) {
		_, err := resolvePath(env, "L")
		require.Error(t, err)
		require.Regexp(t, `too many levels of symbolic links|is outside the session's working directory`, err.Error())
	})

	t.Run("deep missing path under an outside-pointing link is still an escape", func(t *testing.T) {
		require.NoError(t, os.Symlink(outside, filepath.Join(work, "outlink")))
		parts := []string{"outlink"}
		for i := range 300 {
			parts = append(parts, fmt.Sprintf("d%d", i))
		}
		parts = append(parts, "f.txt")
		_, err := resolvePath(env, filepath.Join(parts...))
		require.Error(t, err)
		require.Contains(t, err.Error(), "is outside the session's working directory")
	})
}

func TestBetaAgentToolset(t *testing.T) {
	ts := BetaAgentToolset20260401(&AgentToolContext{Workdir: t.TempDir()})
	defer CloseAll(ts)
	got := map[string]bool{}
	for _, tool := range ts {
		got[tool.Name()] = true
		require.NotEmpty(t, tool.Description(), "tool %q must have a description", tool.Name())
	}
	for _, name := range []string{"bash", "read", "write", "edit", "glob", "grep"} {
		require.True(t, got[name], "agent_toolset_20260401 tool %q must be returned by BetaAgentToolset20260401", name)
	}
}

// closerTool is a BetaTool whose Close behaviour is controlled by onClose. Used
// to verify CloseAll's per-tool isolation.
type closerTool struct {
	name    string
	onClose func() error
	closed  bool
}

func (c *closerTool) Name() string        { return c.name }
func (c *closerTool) Description() string { return c.name }
func (c *closerTool) InputSchema() anthropic.BetaToolInputSchemaParam {
	return anthropic.BetaToolInputSchemaParam{Properties: map[string]any{}}
}
func (c *closerTool) Execute(context.Context, json.RawMessage) ([]anthropic.BetaToolResultBlockParamContentUnion, error) {
	return nil, nil
}
func (c *closerTool) Close() error {
	c.closed = true
	if c.onClose != nil {
		return c.onClose()
	}
	return nil
}

func TestCloseAllIsolatesPanicsAndErrors(t *testing.T) {
	good := &closerTool{name: "good"}
	panicker := &closerTool{name: "panicker", onClose: func() error {
		panic("simulated tool close panic")
	}}
	errored := &closerTool{name: "errored", onClose: func() error {
		return errors.New("simulated tool close error")
	}}

	// Order chosen so the panicking tool sits between two well-behaved ones:
	// this verifies CloseAll doesn't short-circuit on either side.
	CloseAll([]anthropic.BetaTool{good, panicker, errored})

	require.True(t, good.closed,
		"tool before the panicker must still see Close")
	require.True(t, errored.closed,
		"tool after the panicker must still see Close — per-tool recover keeps the loop running")
}

// A tool-input refusal is matchable by name, not only by its wording.
func TestToolRefusalIsAToolError(t *testing.T) {
	env := &AgentToolContext{Workdir: t.TempDir()}
	tools := []anthropic.BetaTool{BetaReadTool(env), BetaBashTool(env)}
	defer CloseAll(tools)

	tests := []struct {
		description string
		tool        anthropic.BetaTool
		input       map[string]any
		wantContent string
	}{
		{
			description: "a file path that escapes the workdir is refused as a *ToolError",
			tool:        tools[0],
			input:       map[string]any{"file_path": "../outside.txt"},
			wantContent: `read: path "../outside.txt" is outside the session's working directory`,
		},
		{
			description: "a missing required argument is refused as a *ToolError",
			tool:        tools[1],
			input:       map[string]any{},
			wantContent: "bash: command is required",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			_, err := tc.tool.Execute(context.Background(), mustJSON(t, tc.input))
			var toolErr *ToolError
			require.ErrorAs(t, err, &toolErr)
			require.Equal(t, tc.wantContent, toolErr.Content)
			require.EqualError(t, err, tc.wantContent, "the typed error keeps the exact message")
		})
	}
}
