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
	absWork := realpathOrSelf(absOrSelf(work))

	// A sibling directory whose name shares a prefix with the workdir, used to
	// verify the prefix check compares whole path segments and not raw string
	// prefixes.
	sibling := absWork + "extra"

	tests := []struct {
		description string
		env         *AgentToolContext
		input       string
		want        string
		wantErr     bool
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
			wantErr:     true,
		},
		{
			description: "dot-dot that stays inside the workdir after normalisation is permitted",
			env:         &AgentToolContext{Workdir: work},
			input:       filepath.Join("sub", "..", "c.txt"),
			want:        filepath.Join(absWork, "c.txt"),
		},
		{
			description: "absolute path outside workdir is rejected when UnrestrictedPaths is false",
			env:         &AgentToolContext{Workdir: work},
			input:       "/etc/passwd",
			wantErr:     true,
		},
		{
			description: "absolute path inside workdir is permitted when UnrestrictedPaths is false",
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
			wantErr:     true,
		},
		{
			description: "absolute path is returned cleaned when UnrestrictedPaths is true",
			env:         &AgentToolContext{Workdir: work, UnrestrictedPaths: true},
			input:       "/etc/passwd",
			want:        "/etc/passwd",
		},
		{
			description: "dot-dot escape is permitted when UnrestrictedPaths is true since the operator opted out of the jail",
			env:         &AgentToolContext{Workdir: work, UnrestrictedPaths: true},
			input:       filepath.Join("..", "escape.txt"),
			want:        filepath.Join(absWork, "..", "escape.txt"),
		},
		{
			description: "sibling directory that string-prefixes the workdir is accepted as an absolute path under UnrestrictedPaths",
			env:         &AgentToolContext{Workdir: work, UnrestrictedPaths: true},
			input:       sibling,
			want:        sibling,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			got, err := resolvePath(tc.env, tc.input)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
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
	require.Equal(t, filepath.Join(realpathOrSelf(absOrSelf(work)), "sub", "f.txt"), got)

	// A dangling symlink whose target is inside the workdir resolves to that
	// target even when the target's parent directory does not exist yet.
	require.NoError(t, os.Symlink(filepath.Join("newdir", "f.txt"), filepath.Join(work, "d")))
	got, err = resolvePath(env, "d")
	require.NoError(t, err)
	require.Equal(t, filepath.Join(realpathOrSelf(absOrSelf(work)), "newdir", "f.txt"), got)
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
	absWork := realpathOrSelf(absOrSelf(work))

	tests := []struct {
		description string
		input       string
		wantErr     string
	}{
		{"two-link cycle", "loop_a", `path "loop_a": too many levels of symbolic links`},
		{"child under a two-link cycle", "loop_a/child.txt", `path "loop_a/child.txt": too many levels of symbolic links`},
		{"self-referencing link", "self", `path "self": too many levels of symbolic links`},
		{"child under a self-referencing link", "self/x", `path "self/x": too many levels of symbolic links`},
		{"dot-dot past a cycle is collapsed lexically and lands on the escaping link", "loop_a/../evil_link", `path "loop_a/../evil_link" escapes workdir`},
		{"dot-dot past a self link is collapsed lexically and lands on the escaping link", "self/../evil_link", `path "self/../evil_link" escapes workdir`},
		{"dangling link pointing outside the workdir", "dangle_out", `path "dangle_out" escapes workdir`},
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
		require.Regexp(t, `too many levels of symbolic links|escapes workdir`, err.Error())
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
		require.Contains(t, err.Error(), "escapes workdir")
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
			wantContent: `read: path "../outside.txt" escapes workdir`,
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
