// Package agenttoolset provides Node-equivalent local executors for the
// `agent_toolset_20260401` tool set — `bash`, `read`, `write`, `edit`, `glob`,
// `grep` — plus the workdir/skills [AgentToolContext] and the skill-download
// helpers ([AgentToolContext.SetupSkillsFromSession], and the deprecated
// [AgentToolContext.SetupSkills] for a caller holding only a session id).
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
//     [BetaGlobTool], [BetaGrepTool]) confine to Workdir plus any
//     AllowedRoots. resolvePath canonicalizes the target — resolving every
//     symlink, including the leaf, even a dangling one — before the
//     containment check and returns that canonical path for the operation, so
//     a symlink inside a permitted root that points outside it neither passes
//     the check nor gets followed afterwards. This is a real boundary, not a
//     lexical hint (modulo the residual TOCTOU noted on resolvePath).
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
// Workdir and reject anything outside Workdir and AllowedRoots (symlinks
// resolved); [BetaBashTool] runs an unrestricted /bin/bash regardless.
type AgentToolContext struct {
	// Workdir is the base directory for resolving relative tool paths.
	Workdir string
	// Deprecated: no longer supported and slated for removal. The file tools are
	// always confined to Workdir plus AllowedRoots, so a true value is not
	// reinterpreted: every call to a tool built from this context fails with
	// ErrUnrestrictedPathsUnsupported instead. List extra directories the file
	// tools may reach in AllowedRoots.
	UnrestrictedPaths bool
	// AllowedRoots are directories outside Workdir the file tools may also
	// reach; the environment worker sets this to the session's memory-store
	// folders. Entries are resolved (symlinks followed) on every call, so an
	// entry reaches whatever it points at when the tool runs; empty entries are
	// ignored and relative entries resolve against the process working
	// directory, like Workdir. Does not constrain [BetaBashTool].
	AllowedRoots []string

	// MaxFileBytes caps the size of a file the read and edit tools will load
	// into memory: it applies to a whole-file read and to edit, which both read
	// the whole file. A read with view_range on a larger file streams it and
	// applies the cap to the selected lines instead. Zero (the default) uses the
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

	// ReadOnlyRoots are directories the write and edit tools refuse to
	// modify, resolved on every call like AllowedRoots. The environment worker
	// sets this to the roots of read-only memory stores so the agent sees the
	// error at write time instead of the change silently never syncing; the
	// mechanism is generic to any directory. Does not constrain [BetaBashTool].
	ReadOnlyRoots []string

	// downloadedSkillDirs are the per-skill directories
	// [AgentToolContext.SetupSkillsFromSession] created under {Workdir}/skills;
	// [AgentToolContext.Cleanup] removes exactly these and nothing else.
	downloadedSkillDirs []string
}

// BetaAgentToolset20260401 returns the six built-in agent_toolset_20260401
// implementations bound to env, in the anthropic.BetaTool shape the SDK's tool
// runners accept. The slice is owned by the caller; filter or append before
// passing to a runner.
//
// [BetaBashTool] keeps a persistent shell open until its Close is called;
// [CloseAll] releases every tool in the slice that implements io.Closer.
// Every tool call fails with [ErrUnrestrictedPathsUnsupported] while the
// deprecated env.UnrestrictedPaths is set.
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

// ErrUnrestrictedPathsUnsupported is what every tool call returns, resolvePath
// returns, and the environment worker's Run/HandleItem wrap, when the
// deprecated UnrestrictedPaths option is set.
var ErrUnrestrictedPathsUnsupported = errors.New(
	"agenttoolset: UnrestrictedPaths is no longer supported; remove it and list any extra directories the file tools may reach in AllowedRoots")

func rejectUnrestrictedPaths(env *AgentToolContext) error {
	if env.UnrestrictedPaths {
		return ErrUnrestrictedPathsUnsupported
	}
	return nil
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

// ToolError is a tool's error result — a refused path, a missing argument,
// a failed command — carrying the text the model sees with is_error set.
type ToolError struct {
	// Content is the error-result text shown to the model.
	Content string
}

func (e *ToolError) Error() string { return e.Content }

// funcTool adapts a plain function into an anthropic.BetaTool. Used for
// stateless tools that only need the shared AgentToolContext. Soft tool failures (the
// (string, true) return) become a *[ToolError] so the surrounding tool
// runner surfaces them to the model as an error result.
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
	if err := rejectUnrestrictedPaths(t.env); err != nil {
		return nil, err
	}
	content, isErr := t.run(ctx, input, t.env)
	if isErr {
		return nil, &ToolError{Content: content}
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
// ENOENT from another SDK). It never returns Go's "op /abs/path: ..." text,
// so the host's absolute path does not reach the model.
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
	case errors.Is(err, errSymlinkLoop), errors.Is(err, syscall.ELOOP):
		return errSymlinkLoop.Error()
	}
	var errno syscall.Errno
	if errors.As(err, &errno) {
		return errno.Error()
	}
	return "i/o error"
}

// resolvePath resolves p against env.Workdir and rejects results outside the
// permitted roots: Workdir plus each non-empty entry of env.AllowedRoots.
// Absolute and relative inputs go through the same canonicalise-then-contain
// check — an absolute path that lands inside a permitted root is accepted,
// only paths that resolve outside all of them are rejected. Every symlink in p
// (including the leaf, even a dangling one) is resolved before the check, the
// roots are resolved the same way on every call, and the resolved path is what
// the tool then operates on, so a symlink inside a permitted root that points
// outside it can neither pass the check nor be followed afterwards. ".." is
// collapsed lexically before any symlink is followed, and a path whose
// symlinks cannot be resolved (a loop, an unreadable component) is rejected
// rather than used as given. See the package-level trust model.
//
// Residual TOCTOU: a component could still be swapped for a symlink between this
// call and the eventual filesystem operation. Closing that fully needs
// per-component O_NOFOLLOW/openat, which Go does not expose ergonomically; the
// same residual exposure exists in the SDK's other file-tool helpers and is why
// a sandbox is still recommended for the toolset as a whole.
func resolvePath(env *AgentToolContext, p string) (string, error) {
	if err := rejectUnrestrictedPaths(env); err != nil {
		return "", err
	}
	root := canonicalRoot(env.Workdir)
	abs := filepath.Clean(p)
	if !filepath.IsAbs(p) {
		abs = filepath.Join(root, p)
	}
	real, err := canonicalize(abs)
	if err != nil {
		return "", fmt.Errorf("path %q: %s", p, fsErrorMessage(err))
	}
	if within(real, root) {
		return real, nil
	}
	hasOtherRoots := false
	for _, r := range env.AllowedRoots {
		if r == "" {
			continue // names no directory; filepath.Abs would turn it into the cwd
		}
		hasOtherRoots = true
		if within(real, canonicalRoot(r)) {
			return real, nil
		}
	}
	if hasOtherRoots {
		return "", fmt.Errorf("path %q is outside the session's working directory and its other permitted directories", p)
	}
	return "", fmt.Errorf("path %q is outside the session's working directory", p)
}

// resolveWritablePath is resolvePath for the write and edit tools: the
// resolved target must also lie outside every env.ReadOnlyRoots entry.
func resolveWritablePath(env *AgentToolContext, p string) (string, error) {
	real, err := resolvePath(env, p)
	if err != nil {
		return "", err
	}
	if ro := readOnlyRootFor(env, real); ro != "" {
		return "", fmt.Errorf("%s is inside read-only directory %s", p, ro)
	}
	return real, nil
}

// readOnlyRootFor returns the entry of env.ReadOnlyRoots that target falls
// under, or "". Roots resolve here on every call, symmetric with AllowedRoots
// in resolvePath, so a symlinked store path cannot grant access on one side
// while its write protection misses on the other. target is canonicalized
// again so the comparison also holds for a caller that did not obtain it from
// resolvePath, e.g. under a symlinked prefix such as /tmp on macOS.
func readOnlyRootFor(env *AgentToolContext, target string) string {
	if len(env.ReadOnlyRoots) == 0 {
		return ""
	}
	tgt := filepath.Clean(target)
	if real, err := canonicalize(tgt); err == nil {
		tgt = real
	}
	for _, root := range env.ReadOnlyRoots {
		if within(tgt, canonicalRoot(root)) {
			return root
		}
	}
	return ""
}

// within reports whether path equals root or lies beneath it, comparing whole
// segments so /work does not contain /workextra. Both must be clean absolute
// paths; root may be the filesystem root.
func within(path, root string) bool {
	sep := string(filepath.Separator)
	return path == root || strings.HasPrefix(path, strings.TrimSuffix(root, sep)+sep)
}

// canonicalRoot resolves a configured root (Workdir, an AllowedRoots or
// ReadOnlyRoots entry) the way resolvePath resolves targets, so both sides of
// a containment check have had the same symlinks followed. A root whose
// symlinks cannot be resolved is compared as its absolute path, which no
// canonicalized target lies under, so it grants nothing.
func canonicalRoot(r string) string {
	abs := absOrSelf(r)
	if real, err := canonicalize(abs); err == nil {
		return real
	}
	return abs
}

func absOrSelf(p string) string {
	if abs, err := filepath.Abs(p); err == nil {
		return abs
	}
	return filepath.Clean(p)
}

var errSymlinkLoop = errors.New("too many levels of symbolic links")

const maxSymlinkHops = 40

// canonicalize returns abs with every symlink resolved, or an error — never
// abs itself or a partly resolved path. It EvalSymlinks the longest existing
// ancestor and re-appends the missing remainder so paths being created
// (write/edit) still resolve, but a component that is itself a symlink is read
// and followed rather than re-appended, so a dangling link pointing outside the
// workdir cannot slip through. Only symlink reads count against
// maxSymlinkHops; exceeding it, or an ELOOP from the OS, is a loop.
func canonicalize(abs string) (string, error) {
	var tail []string
	prefix := filepath.Clean(abs)
	hops := 0
	for {
		real, evalErr := filepath.EvalSymlinks(prefix)
		if evalErr == nil {
			parts := make([]string, 0, len(tail)+1)
			parts = append(parts, real)
			for i := len(tail) - 1; i >= 0; i-- {
				parts = append(parts, tail[i])
			}
			return filepath.Join(parts...), nil
		}
		fi, lerr := os.Lstat(prefix)
		switch {
		case lerr == nil && fi.Mode()&os.ModeSymlink != 0:
			hops++
			if hops > maxSymlinkHops {
				return "", errSymlinkLoop
			}
			dest, err := os.Readlink(prefix)
			if err != nil {
				return "", err
			}
			if !filepath.IsAbs(dest) {
				dest = filepath.Join(filepath.Dir(prefix), dest)
			}
			prefix = filepath.Clean(dest)
		case lerr == nil:
			return "", evalErr
		case errors.Is(lerr, fs.ErrNotExist), errors.Is(lerr, syscall.ENOTDIR):
			parent := filepath.Dir(prefix)
			if parent == prefix {
				return "", lerr
			}
			tail = append(tail, filepath.Base(prefix))
			prefix = parent
		default:
			return "", lerr
		}
	}
}
