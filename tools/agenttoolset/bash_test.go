//go:build !windows

package agenttoolset

import (
	"context"
	"io"
	"strings"
	"testing"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/stretchr/testify/require"
)

func newBashTool(t *testing.T, dir string) anthropic.BetaTool {
	t.Helper()
	tool := BetaBashTool(&AgentToolContext{Workdir: dir})
	t.Cleanup(func() {
		if c, ok := tool.(io.Closer); ok {
			_ = c.Close()
		}
	})
	return tool
}

func TestBashSessionPersistence(t *testing.T) {
	tool := newBashTool(t, t.TempDir())

	tests := []struct {
		description string
		command     string
		want        string
	}{
		{
			description: "first command exports a variable into the persistent shell environment",
			command:     "export FOO=bar; echo set",
			want:        "set",
		},
		{
			description: "second command in the same session can read the variable exported by the first",
			command:     "echo $FOO",
			want:        "bar",
		},
		{
			description: "changing directory persists so a later pwd reflects the new location",
			command:     "cd /tmp && pwd",
			want:        "/tmp",
		},
		{
			description: "subsequent pwd without cd still reports /tmp confirming cwd survived across calls",
			command:     "pwd",
			want:        "/tmp",
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			out, isErr := runTool(t, tool, mustJSON(t, map[string]any{"command": tc.command}))
			require.False(t, isErr, "unexpected error output: %q", out)
			require.Equal(t, tc.want, strings.TrimSpace(out))
		})
	}
}

func TestBashRestart(t *testing.T) {
	tool := newBashTool(t, t.TempDir())

	_, isErr := runTool(t, tool, mustJSON(t, map[string]any{"command": "export GONE=1"}))
	require.False(t, isErr)

	out, isErr := runTool(t, tool, mustJSON(t, map[string]any{"restart": true, "command": "echo [$GONE]"}))
	require.False(t, isErr, "output=%q", out)
	require.Equal(t, "[]", strings.TrimSpace(out), "variable should be unset after restart")
}

func TestBashNonZeroExit(t *testing.T) {
	tool := newBashTool(t, t.TempDir())

	// Use a subshell so the non-zero status is observed without terminating
	// the persistent session itself.
	out, isErr := runTool(t, tool, mustJSON(t, map[string]any{"command": "(exit 7)"}))
	require.True(t, isErr, "non-zero exit must set is_error so the model knows the command failed")
	require.Contains(t, out, "exit code: 7")
}

func TestBashTimeout(t *testing.T) {
	work := t.TempDir()
	sess, err := NewBashSession(work, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = sess.Close() })

	start := time.Now()
	out, code, err := sess.Exec(context.Background(), "sleep 5", 200*time.Millisecond)
	require.ErrorIs(t, err, ErrTimedOut, "Exec returns ErrTimedOut so the caller knows to restart the session")
	require.Equal(t, -1, code, "timed-out command reports sentinel exit code")
	require.Contains(t, out, "[timed out]")
	require.Less(t, time.Since(start), 2*time.Second, "Exec must return promptly after the timeout rather than blocking on sleep")
}

// TestBashToolRestartsAfterCtxCancel pins the guarantee that the
// persistent shell is restarted when the per-tool ctx fires before
// bashDefaultTO does — Exec returns ctx.Err() rather than ErrTimedOut,
// and the still-running command's stdout plus the per-call exit-code
// sentinel would otherwise land in the next Exec's buffer (the next
// call then sees stale output or, worse, false-matches the prior
// sentinel and reports the wrong exit code).
//
// This test uses a tight ctx (200 ms) with a long bash-side timeout
// (30 s, much greater than bashDefaultTO would matter) so ctx fires
// first; the wrapped command writes "before" then sleeps. With the fix,
// the second call sees only "after". Without the fix, the second call's
// output either contains "before" or the prior sentinel.
func TestBashToolRestartsAfterCtxCancel(t *testing.T) {
	tool := newBashTool(t, t.TempDir())

	// First call: ctx fires before the bash-side timeout would. The shell
	// is left mid-command on the other side of the PTY with `sleep 5`
	// still running and "before\n" buffered (plus the per-call sentinel
	// that will land once `sleep` returns).
	tightCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	_, err := tool.Execute(tightCtx, mustJSON(t, map[string]any{
		"command":    "echo before; sleep 5",
		"timeout_ms": 30000,
	}))
	require.Error(t, err, "ctx cancel must surface as a tool error")

	// Second call must run in a fresh session — neither the previous
	// command's stdout nor its sentinel should be visible.
	out2, isErr2 := runTool(t, tool, mustJSON(t, map[string]any{"command": "echo after"}))
	require.False(t, isErr2)
	require.Contains(t, out2, "after")
	require.NotContains(t, out2, "before", "second call must not see leftover output from the ctx-cancelled command")
	require.NotContains(t, out2, "__ANT_CMD_", "second call must not see the ctx-cancelled command's sentinel")
}

func TestBashToolRestartsAfterTimeout(t *testing.T) {
	tool := newBashTool(t, t.TempDir())

	// First call hangs past its timeout; the tool must restart so the still-running
	// `sleep` cannot leak its sentinel into the next call's buffer.
	out, isErr := runTool(t, tool, mustJSON(t, map[string]any{
		"command":    "echo before; sleep 5",
		"timeout_ms": 200,
	}))
	require.True(t, isErr)
	require.Contains(t, out, "session restarted after timeout")

	// Second call must execute in a fresh session and see only its own output.
	out2, isErr2 := runTool(t, tool, mustJSON(t, map[string]any{"command": "echo after"}))
	require.False(t, isErr2)
	require.Contains(t, out2, "after")
	require.NotContains(t, out2, "before", "second call must not see leftover output from the timed-out group")
	require.NotContains(t, out2, "__ANT_CMD_", "second call must not see the timed-out group's sentinel")
}

func TestBashMissingCommand(t *testing.T) {
	tool := newBashTool(t, t.TempDir())
	_, isErr := runTool(t, tool, mustJSON(t, map[string]any{}))
	require.True(t, isErr, "empty command without restart must be rejected")
}

// TestBashScrubsAnthropicCredentials verifies the spawned shell does not
// inherit the runner's ANTHROPIC_* credentials (API key, auth/session tokens),
// while ordinary variables like PATH still pass through.
func TestBashScrubsAnthropicCredentials(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "sk-ant-should-not-leak")
	t.Setenv("ANTHROPIC_AUTH_TOKEN", "tok-should-not-leak")
	tool := newBashTool(t, t.TempDir())

	out, isErr := runTool(t, tool, mustJSON(t, map[string]any{
		"command": `echo "key=[$ANTHROPIC_API_KEY] token=[$ANTHROPIC_AUTH_TOKEN] path=[${PATH:+set}]"`,
	}))
	require.False(t, isErr, "output=%q", out)
	require.Contains(t, out, "key=[]", "ANTHROPIC_API_KEY must not be visible to the spawned shell")
	require.Contains(t, out, "token=[]", "ANTHROPIC_AUTH_TOKEN must not be visible to the spawned shell")
	require.Contains(t, out, "path=[set]", "PATH must survive the credential scrub")
}

// TestBashEnvFullyReplaces verifies that a non-nil AgentToolContext.Env FULLY
// REPLACES the default scrubbed process environment rather than being merged
// into it: only the provided keys are visible, and an inherited variable like
// PATH is gone.
func TestBashEnvFullyReplaces(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "sk-ant-should-not-leak")
	t.Setenv("BASH_TEST_INHERITED", "from-process")

	tool := BetaBashTool(&AgentToolContext{
		Workdir: t.TempDir(),
		Env:     map[string]string{"CUSTOM_ONLY": "yes"},
	})
	t.Cleanup(func() {
		if c, ok := tool.(io.Closer); ok {
			_ = c.Close()
		}
	})

	// BASH_TEST_INHERITED is a non-ANTHROPIC process var: it would survive the
	// credential scrub and so would be visible if Env were merged into the
	// default instead of replacing it. Its absence proves full replacement.
	// (PATH is not a reliable probe — bash synthesizes a compiled-in default
	// PATH even when started with an empty environment.)
	out, isErr := runTool(t, tool, mustJSON(t, map[string]any{
		"command": `echo "custom=[$CUSTOM_ONLY] inherited=[$BASH_TEST_INHERITED] key=[$ANTHROPIC_API_KEY]"`,
	}))
	require.False(t, isErr, "output=%q", out)
	require.Contains(t, out, "custom=[yes]", "the provided Env mapping must be used verbatim")
	require.Contains(t, out, "inherited=[]",
		"a non-nil Env FULLY REPLACES the default; inherited process vars must not leak in")
	require.Contains(t, out, "key=[]", "ANTHROPIC_* credentials must never be visible")
}

func TestBashSentinelNotSpoofable(t *testing.T) {
	sess, err := NewBashSession(t.TempDir(), nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = sess.Close() })

	// A command that prints a hardcoded sentinel-like marker must not truncate
	// its own output or spoof the exit code: the per-call random nonce makes
	// the real marker unguessable.
	out, code, err := sess.Exec(context.Background(), "printf '__ANT_CMD_DONE__7\\nafter\\n'; (exit 3)", 0)
	require.NoError(t, err)
	require.Contains(t, out, "__ANT_CMD_DONE__7")
	require.Contains(t, out, "after")
	require.Equal(t, 3, code)
}

func TestBashStdinRedirect(t *testing.T) {
	sess, err := NewBashSession(t.TempDir(), nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = sess.Close() })

	// `cat` with an open stdin would block until the timeout; the </dev/null
	// redirect on the wrapper lets it exit cleanly with EOF.
	out, code, err := sess.Exec(context.Background(), "cat; echo done", 2*time.Second)
	require.NoError(t, err)
	require.Equal(t, "done", strings.TrimSpace(out))
	require.Equal(t, 0, code)
}

func TestBashOutputBufferBounded(t *testing.T) {
	sess, err := NewBashSession(t.TempDir(), nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = sess.Close() })

	// A command that streams more than the cap must be truncated rather than
	// buffered unbounded; the tail (and the sentinel) must survive.
	out, code, err := sess.Exec(context.Background(), "head -c 300000 /dev/zero | tr '\\0' a; echo END", 0)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(out, "[output truncated]\n"), "truncated output is annotated")
	require.True(t, strings.HasSuffix(strings.TrimSpace(out), "END"), "the tail of the output is kept")
	require.LessOrEqual(t, len(out), bashOutputLimit+len("[output truncated]\n"))
	require.Equal(t, 0, code)
}
