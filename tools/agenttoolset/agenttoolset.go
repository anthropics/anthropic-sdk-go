// Package agenttoolset provides Node-equivalent local executors for the
// `agent_toolset_20260401` tool set — `bash`, `read`, `write`, `edit`, `glob`,
// `grep` — plus the workdir/skills [AgentToolContext] and the skill-download helper
// ([AgentToolContext.SetupSkills]).
//
// This mirrors the SDK's other first-class tool modules: it is the explicit
// entry point for these implementations. Importing it pulls in os/exec, a PTY
// dependency, etc., so it is kept separate from the rest of the SDK — depending
// on it is opt-in.
//
// The result of [BetaAgentToolset20260401] is a plain []anthropic.BetaTool; hand it to
// any tool runner — client.Beta.Messages.NewToolRunner(tools, …) for the
// Messages API, or client.Beta.Sessions.Events.NewToolRunner(…) for a
// managed-agents session:
//
//	import "github.com/anthropics/anthropic-sdk-go/tools/agenttoolset"
//
//	env := &agenttoolset.AgentToolContext{Workdir: "/work"}
//	tools := agenttoolset.BetaAgentToolset20260401(env)
//
// Trust model — two tiers:
//
//   - The file tools ([BetaReadTool], [BetaWriteTool], [BetaEditTool],
//     [BetaGlobTool], [BetaGrepTool]) confine to Workdir unless
//     UnrestrictedPaths is set. resolvePath canonicalizes the target —
//     resolving every symlink, including the leaf, even a dangling one — before
//     the workdir check and returns that canonical path for the operation, so a
//     symlink inside the workdir that points outside it neither passes the check
//     nor gets followed afterwards. This is a real boundary, not a lexical hint
//     (modulo the residual TOCTOU noted on resolvePath).
//   - [BetaBashTool] runs an unrestricted /bin/bash and cannot be confined. Run
//     it — and, for defense in depth, the whole toolset — inside a sandbox the
//     host controls (e.g. a self-hosted environment runner).
package agenttoolset

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	anthropic "github.com/anthropics/anthropic-sdk-go"
)

// AgentToolContext carries per-session configuration the agent_toolset_20260401
// executors need.
//
// See the package-level trust model: the file tools resolve paths against
// Workdir and reject escapes (symlinks resolved) unless UnrestrictedPaths is
// set; [BetaBashTool] runs an unrestricted /bin/bash regardless.
type AgentToolContext struct {
	// Workdir is the base directory for resolving relative tool paths.
	Workdir string
	// UnrestrictedPaths controls whether the file tools accept paths that
	// resolve outside Workdir. When false (default) they are rejected.
	// Does not constrain [BetaBashTool].
	UnrestrictedPaths bool

	// MaxFileBytes caps the size of a file the read and edit tools will load
	// into memory (both read the whole file). Zero (the default) uses the
	// built-in 256 KiB cap; a positive value sets a custom cap; a negative
	// value disables the size cap entirely. Disabling it reintroduces the OOM
	// risk on a model-controlled path, so set it negative only when the sandbox
	// can absorb arbitrarily large files. The non-regular-file (FIFO/device)
	// guard always applies regardless of this value.
	MaxFileBytes int64

	// Env sets the bash subprocess environment. When non-nil it fully replaces
	// the inherited environment with exactly these entries; when nil the
	// subprocess inherits the runner's environment with ANTHROPIC_* credentials
	// removed.
	//
	// Env is used verbatim and is not filtered, so a tool whose script calls the
	// Anthropic API can deliberately pass routing or config vars such as
	// ANTHROPIC_BASE_URL. For that same reason, do not build Env by copying
	// os.Environ() and adding a few entries — that leaks the runner's own
	// ANTHROPIC_* credentials into a model-driven shell. Populate Env with only
	// the variables the tools need.
	Env map[string]string
}

// BetaAgentToolset20260401 returns the six built-in agent_toolset_20260401
// implementations bound to env, in the anthropic.BetaTool shape the SDK's tool
// runners accept. The slice is owned by the caller; filter or append before
// passing to a runner.
//
// [BetaBashTool] keeps a persistent shell open until its Close is called;
// [CloseAll] releases every tool in the slice that implements io.Closer.
func BetaAgentToolset20260401(env *AgentToolContext) []anthropic.BetaTool {
	return []anthropic.BetaTool{
		BetaBashTool(env),
		BetaReadTool(env),
		BetaWriteTool(env),
		BetaEditTool(env),
		BetaGlobTool(env),
		BetaGrepTool(env),
	}
}

// CloseAll releases resources held by any tools that implement io.Closer. Each
// Close runs under its own recover so one panicking tool cannot skip cleanup
// for the rest of the slice. Errors and panics are swallowed — callers wanting
// visibility should close tools themselves and inspect the return values.
func CloseAll(ts []anthropic.BetaTool) {
	for _, t := range ts {
		func(t anthropic.BetaTool) {
			defer func() { _ = recover() }()
			if c, ok := t.(io.Closer); ok {
				_ = c.Close()
			}
		}(t)
	}
}

// funcTool adapts a plain function into an anthropic.BetaTool. Used for
// stateless tools that only need the shared AgentToolContext. Soft tool failures (the
// (string, true) return) become a Go error so the surrounding tool runner
// surfaces them to the model as an error result.
type funcTool struct {
	name        string
	description string
	schema      anthropic.BetaToolInputSchemaParam
	env         *AgentToolContext
	run         func(ctx context.Context, input json.RawMessage, env *AgentToolContext) (string, bool)
}

func (t *funcTool) Name() string                                    { return t.name }
func (t *funcTool) Description() string                             { return t.description }
func (t *funcTool) InputSchema() anthropic.BetaToolInputSchemaParam { return t.schema }
func (t *funcTool) Execute(ctx context.Context, input json.RawMessage) ([]anthropic.BetaToolResultBlockParamContentUnion, error) {
	content, isErr := t.run(ctx, input, t.env)
	if isErr {
		return nil, errors.New(content)
	}
	return textResult(content), nil
}

// textResult wraps a plain string as a single text tool-result block.
func textResult(s string) []anthropic.BetaToolResultBlockParamContentUnion {
	return []anthropic.BetaToolResultBlockParamContentUnion{{OfText: &anthropic.BetaTextBlockParam{Text: s}}}
}

func objectSchema(properties map[string]any, required ...string) anthropic.BetaToolInputSchemaParam {
	return anthropic.BetaToolInputSchemaParam{Properties: properties, Required: required}
}

func prop(typ, description string) map[string]any {
	return map[string]any{"type": typ, "description": description}
}

func errorf(format string, a ...any) (string, bool) {
	return fmt.Sprintf(format, a...), true
}

// fsErrorMessage maps a filesystem error to a consistent, language-independent
// message so the model sees the same wording regardless of the host runtime's
// raw errno text (e.g. Go's "open x: no such file or directory" vs a bare
// ENOENT from another SDK). Unrecognised errors fall through to their text.
func fsErrorMessage(err error) string {
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return "no such file or directory"
	case errors.Is(err, fs.ErrPermission):
		return "permission denied"
	case errors.Is(err, syscall.ENOTDIR):
		return "not a directory"
	case errors.Is(err, syscall.EISDIR):
		return "is a directory"
	default:
		return err.Error()
	}
}

// resolvePath resolves p against env.Workdir. Absolute and relative inputs go
// through the same canonicalise-then-contain check — an absolute path that
// lands inside the workdir is permitted, only paths that resolve outside are
// rejected. Every symlink in p (including the leaf, even a dangling one) is
// resolved before the workdir check, and the resolved path is what the tool
// then operates on, so a symlink inside the workdir that points outside it can
// neither pass the check nor be followed afterwards. See the package-level
// trust model.
//
// Residual TOCTOU: a component could still be swapped for a symlink between this
// call and the eventual filesystem operation. Closing that fully needs
// per-component O_NOFOLLOW/openat, which Go does not expose ergonomically; the
// same residual exposure exists in the SDK's other file-tool helpers and is why
// a sandbox is still recommended for the toolset as a whole.
func resolvePath(env *AgentToolContext, p string) (string, error) {
	if env.UnrestrictedPaths && filepath.IsAbs(p) {
		return filepath.Clean(p), nil
	}
	root := realpathOrSelf(absOrSelf(env.Workdir))
	abs := filepath.Clean(p)
	if !filepath.IsAbs(p) {
		abs = filepath.Join(root, p)
	}
	if env.UnrestrictedPaths {
		return abs, nil
	}
	real := canonicalize(abs)
	if real != root && !strings.HasPrefix(real, root+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q escapes workdir", p)
	}
	return real, nil
}

func absOrSelf(p string) string {
	if abs, err := filepath.Abs(p); err == nil {
		return abs
	}
	return filepath.Clean(p)
}

func realpathOrSelf(p string) string {
	if real, err := filepath.EvalSymlinks(p); err == nil {
		return real
	}
	return p
}

// canonicalize fully resolves abs: EvalSymlinks the longest existing ancestor
// and re-append the rest, but never re-append a component that is itself a
// symlink — read the link and continue from its target instead. This handles
// paths being created (write/edit) without letting a symlink leaf (e.g. a
// dangling one pointing outside the workdir) slip through unresolved. A symlink
// loop falls back to returning abs unchanged after a bounded number of hops.
func canonicalize(abs string) string {
	var tail []string
	prefix := filepath.Clean(abs)
	for hops := 0; hops < 255; hops++ {
		if real, err := filepath.EvalSymlinks(prefix); err == nil {
			if len(tail) == 0 {
				return real
			}
			parts := make([]string, 0, len(tail)+1)
			parts = append(parts, real)
			for i := len(tail) - 1; i >= 0; i-- {
				parts = append(parts, tail[i])
			}
			return filepath.Join(parts...)
		}
		isLink := false
		if fi, err := os.Lstat(prefix); err == nil {
			isLink = fi.Mode()&os.ModeSymlink != 0
		}
		if isLink {
			dest, err := os.Readlink(prefix)
			if err != nil {
				return abs
			}
			if !filepath.IsAbs(dest) {
				dest = filepath.Join(filepath.Dir(prefix), dest)
			}
			prefix = filepath.Clean(dest)
			continue
		}
		parent := filepath.Dir(prefix)
		if parent == prefix {
			return abs // walked past the filesystem root without a hit
		}
		tail = append(tail, filepath.Base(prefix))
		prefix = parent
	}
	return abs
}
