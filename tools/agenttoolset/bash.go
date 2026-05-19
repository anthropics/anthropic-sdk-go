package agenttoolset

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/creack/pty"
)

const (
	bashOutputLimit = 100 * 1024
	bashDefaultTO   = 120 * time.Second
)

// credentialEnvPrefixes are the environment-variable namespaces stripped from
// the bash tool's spawned shell. The runner is typically started with
// ANTHROPIC_* credentials in its environment (API key, auth token, the
// environment key the SDK and the environment worker read); none of those
// belong in a model-controlled shell, so the whole namespace is dropped.
// PATH, HOME, locale, etc. are kept so the shell still behaves normally.
var credentialEnvPrefixes = []string{"ANTHROPIC_"}

// scrubbedEnviron returns os.Environ() with every credential-bearing variable
// (see credentialEnvPrefixes) removed, so the bash tool cannot leak the
// runner's credentials into a command's environment.
func scrubbedEnviron() []string {
	src := os.Environ()
	out := make([]string, 0, len(src))
	for _, kv := range src {
		name := kv
		if i := strings.IndexByte(kv, '='); i >= 0 {
			name = kv[:i]
		}
		scrub := false
		for _, prefix := range credentialEnvPrefixes {
			if strings.HasPrefix(name, prefix) {
				scrub = true
				break
			}
		}
		if !scrub {
			out = append(out, kv)
		}
	}
	return out
}

var ansi = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)

// BashSession is a persistent /bin/bash process attached to a PTY. State
// (cwd, env, background jobs) survives across Exec calls.
type BashSession struct {
	mu        sync.Mutex
	pty       *os.File
	cmd       *exec.Cmd
	buf       bytes.Buffer
	truncated bool
	closed    bool
	done      chan struct{}
	// notify carries a wakeup whenever the drain goroutine appends output, so
	// exec waits on an event instead of polling the buffer on a timer. It is
	// buffered (cap 1) and drained on every select iteration, so a signal is
	// never missed and a burst of writes collapses harmlessly into one wakeup.
	notify chan struct{}
}

var _ io.Closer = (*BashSession)(nil)

// NewBashSession starts bash in dir and begins draining the PTY into an
// internal buffer.
//
// env selects the base environment for the spawned shell. When nil the shell
// inherits the runner's environment minus its ANTHROPIC_* credentials (see
// scrubbedEnviron), so a model-issued command cannot read the API key / auth
// token / environment key out of its own environment. When non-nil it FULLY
// REPLACES that default — the mapping is used verbatim and is NOT merged with
// the scrubbed process environment. PS1/PS2/TERM are always overlaid on the
// chosen base so output stays clean and parseable regardless of the base.
func NewBashSession(dir string, env map[string]string) (*BashSession, error) {
	cmd := exec.Command("/bin/bash", "--noprofile", "--norc")
	cmd.Dir = dir
	// When env is nil, spawn the shell with the runner's environment minus
	// its ANTHROPIC_* credentials. When env is non-nil it FULLY REPLACES
	// that default — the mapping is used verbatim, nothing is merged in.
	base := scrubbedEnviron()
	if env != nil {
		base = make([]string, 0, len(env))
		for k, v := range env {
			base = append(base, k+"="+v)
		}
	}
	cmd.Env = append(base, "PS1=", "PS2=", "TERM=dumb")

	p, err := pty.Start(cmd)
	if err != nil {
		return nil, fmt.Errorf("start bash pty: %w", err)
	}

	s := &BashSession{pty: p, cmd: cmd, done: make(chan struct{}), notify: make(chan struct{}, 1)}
	go s.drain()
	// Disable echo and job-control noise so output is just command results.
	// This call must keep the PTY as its stdin (no </dev/null redirect):
	// stty needs a terminal on fd 0 to set its modes.
	_, _, _ = s.exec(context.Background(), "stty -echo 2>/dev/null; set +m", 5*time.Second, false)
	return s, nil
}

func (s *BashSession) drain() {
	tmp := make([]byte, 4096)
	for {
		n, err := s.pty.Read(tmp)
		if n > 0 {
			s.mu.Lock()
			s.buf.Write(tmp[:n])
			// Cap the buffer during accumulation so a command that streams
			// unboundedly can't OOM the runner. Keep the tail so the sentinel
			// stays detectable.
			if over := s.buf.Len() - bashOutputLimit; over > 0 {
				s.buf.Next(over)
				s.truncated = true
			}
			s.mu.Unlock()
			// Wake any waiting exec; non-blocking so drain never stalls.
			select {
			case s.notify <- struct{}{}:
			default:
			}
		}
		if err != nil {
			close(s.done)
			return
		}
	}
}

// Exec runs cmd in the persistent shell and returns combined output and the
// command's exit code embedded in the output tail.
func (s *BashSession) Exec(ctx context.Context, cmd string, timeout time.Duration) (string, int, error) {
	if timeout <= 0 {
		timeout = bashDefaultTO
	}
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return "", -1, errBashClosed
	}
	s.buf.Reset()
	s.truncated = false
	s.mu.Unlock()
	return s.exec(ctx, cmd, timeout, true)
}

// exec writes cmd wrapped with a completion sentinel and waits for it. When
// redirectStdin is true the command group's stdin is /dev/null so a
// stdin-reading command (`cat`, `read`) gets EOF instead of blocking on the
// PTY until the timeout; the internal `stty -echo` setup call passes false
// because stty needs the PTY on fd 0.
func (s *BashSession) exec(ctx context.Context, cmd string, timeout time.Duration, redirectStdin bool) (string, int, error) {
	// The trailing printf surfaces the exit code on its own line so we can
	// strip it from the user-visible output and report is_error accurately.
	// Per-call random nonce so a command that prints a fixed marker can't
	// spoof the exit-code framing. The sentinel is split across two adjacent
	// quoted strings so PTY echo of the input cannot contain the full
	// sentinel and false-trigger the reader.
	sentinel := newSentinel()
	half := len(sentinel) / 2
	redir := ""
	if redirectStdin {
		redir = " </dev/null"
	}
	wrapped := fmt.Sprintf("{ %s\n}%s 2>&1; printf '\\n%s''%s%%d\\n' $?\n", cmd, redir, sentinel[:half], sentinel[half:])
	if _, err := io.WriteString(s.pty, wrapped); err != nil {
		return "", -1, fmt.Errorf("write to pty: %w", err)
	}

	// Event-driven wait: the drain goroutine signals s.notify whenever it
	// appends output, and the timer bounds the call. No busy polling.
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	timedOut := false
	for {
		select {
		case <-ctx.Done():
			return "", -1, ctx.Err()
		case <-s.done:
			return "", -1, fmt.Errorf("bash session terminated")
		case <-s.notify:
		case <-timer.C:
			timedOut = true
		}
		s.mu.Lock()
		out := append([]byte(nil), s.buf.Bytes()...)
		truncated := s.truncated
		s.mu.Unlock()
		if i := bytes.Index(out, []byte(sentinel)); i >= 0 {
			body := out[:i]
			tail := out[i+len(sentinel):]
			var code int
			// A parse miss means the framing is corrupt; report -1 instead of
			// letting code default to 0.
			if n, err := fmt.Fscan(bytes.NewReader(tail), &code); err != nil || n != 1 {
				slog.Default().Warn("bash: failed to parse exit code from sentinel tail; treating as failure",
					slog.Any("error", err), slog.Int("scanned", n))
				code = -1
			}
			return cleanOutput(body, truncated), code, nil
		}
		if timedOut {
			return cleanOutput(out, truncated) + "\n[timed out]", -1, ErrTimedOut
		}
	}
}

// newSentinel returns a per-call completion marker carrying a random nonce so
// a command can't predict and emit it to spoof the exit-code framing.
func newSentinel() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "__ANT_CMD_" + hex.EncodeToString(b) + "_DONE__"
}

// ErrTimedOut signals the wrapped command did not complete before the
// per-call deadline; the caller must restart the session because the
// running group's sentinel will land in a future buffer otherwise.
var ErrTimedOut = errors.New("bash: command timed out")

// errBashClosed is returned when Exec is called on a closed session.
var errBashClosed = errors.New("session closed")

// Close terminates the bash process group and the PTY. Safe to call multiple
// times.
func (s *BashSession) Close() error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil
	}
	s.closed = true
	s.mu.Unlock()

	var firstErr error
	if s.pty != nil {
		if err := s.pty.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	if s.cmd != nil && s.cmd.Process != nil {
		if err := killProcessGroup(s.cmd.Process); err != nil && firstErr == nil {
			firstErr = err
		}
		// Reap the process; ExitError after SIGKILL is expected.
		var exitErr *exec.ExitError
		if err := s.cmd.Wait(); err != nil && !errors.As(err, &exitErr) && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func cleanOutput(b []byte, truncated bool) string {
	out := ansi.ReplaceAll(b, nil)
	out = bytes.ReplaceAll(out, []byte("\r\n"), []byte("\n"))
	out = bytes.ReplaceAll(out, []byte("\r"), nil)
	out = bytes.TrimLeft(out, "\n")
	if truncated {
		return "[output truncated]\n" + string(out)
	}
	return string(out)
}

// BetaBashTool returns an anthropic.BetaTool backed by a persistent bash
// session rooted at env.Workdir. The session is created lazily on first use and
// persists across calls; the returned tool implements io.Closer.
//
// bash is the one explicitly-unrestricted tool in the set — it runs /bin/bash
// directly and ignores AgentToolContext.UnrestrictedPaths. Run it inside a sandbox you
// control.
func BetaBashTool(env *AgentToolContext) anthropic.BetaTool {
	return &bashTool{env: env}
}

type bashTool struct {
	env  *AgentToolContext
	mu   sync.Mutex
	sess *BashSession
}

func (t *bashTool) Name() string { return "bash" }

func (t *bashTool) Description() string {
	return "Run a bash command in a persistent shell. State (cwd, env vars) persists across calls."
}

func (t *bashTool) InputSchema() anthropic.BetaToolInputSchemaParam {
	return objectSchema(map[string]any{
		"command":    prop("string", "The command to run"),
		"restart":    prop("boolean", "Restart the persistent shell before running"),
		"timeout_ms": prop("integer", "Per-call timeout in milliseconds"),
	})
}

func (t *bashTool) Execute(ctx context.Context, raw json.RawMessage) ([]anthropic.BetaToolResultBlockParamContentUnion, error) {
	content, isErr := t.run(ctx, raw)
	if isErr {
		return nil, errors.New(content)
	}
	return textResult(content), nil
}

func (t *bashTool) session() (*BashSession, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.sess != nil {
		return t.sess, nil
	}
	s, err := NewBashSession(t.env.Workdir, t.env.Env)
	if err != nil {
		return nil, err
	}
	t.sess = s
	return s, nil
}

func (t *bashTool) restart() error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.sess == nil {
		return nil
	}
	err := t.sess.Close()
	t.sess = nil
	return err
}

func (t *bashTool) Close() error {
	return t.restart()
}

func (t *bashTool) run(ctx context.Context, raw json.RawMessage) (string, bool) {
	var in anthropic.BetaManagedAgentsAgentToolset20260401BashInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return errorf("invalid bash input: %v", err)
	}
	if in.Restart {
		_ = t.restart()
		if in.Command == "" {
			return "bash session restarted", false
		}
	}
	if in.Command == "" {
		return errorf("bash: command is required")
	}
	sess, err := t.session()
	if err != nil {
		return errorf("bash: %v", err)
	}
	to := time.Duration(in.TimeoutMs) * time.Millisecond
	out, code, err := sess.Exec(ctx, in.Command, to)
	if err != nil {
		// Any non-nil error from Exec leaves the persistent shell in an
		// unknown state: a context-cancel returns ctx.Err() with the
		// command still running on the other side of the PTY, a pty write
		// error means the shell may be unusable, a "session terminated"
		// means the underlying process is gone. If we kept the session,
		// the still-running command's stdout plus the per-call exit-code
		// sentinel would land in the next Exec's buffer and contaminate
		// it. Restart unconditionally and keep the special-case timeout
		// message so tool consumers see the same outcome as before.
		_ = t.restart()
		if errors.Is(err, ErrTimedOut) {
			return out + "\nsession restarted after timeout", true
		}
		return errorf("bash: %v", err)
	}
	if code != 0 {
		return fmt.Sprintf("%s\nexit code: %d", out, code), true
	}
	return out, false
}
