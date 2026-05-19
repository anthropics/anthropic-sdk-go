package agenttoolset

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	anthropic "github.com/anthropics/anthropic-sdk-go"
)

const (
	globResultLimit  = 200
	grepOutputLimit  = 100 * 1024
	grepMaxFileBytes = 8 * 1024 * 1024
	walkMaxEntries   = 50_000
)

const truncationNotice = "[output truncated]"

// BetaGlobTool returns an anthropic.BetaTool that globs under env.Workdir.
func BetaGlobTool(env *AgentToolContext) anthropic.BetaTool {
	return &funcTool{
		name:        "glob",
		description: "List paths matching a glob pattern (e.g. **/*.go), newest first.",
		schema: objectSchema(map[string]any{
			"pattern": prop("string", "Glob pattern, e.g. **/*.go (** matches any depth)."),
			"path":    prop("string", "Directory to search in. Defaults to the workdir."),
		}, "pattern"),
		env: env,
		run: execGlob,
	}
}

// BetaGrepTool returns an anthropic.BetaTool that searches file contents under
// env.Workdir.
func BetaGrepTool(env *AgentToolContext) anthropic.BetaTool {
	return &funcTool{
		name:        "grep",
		description: "Search file contents for a regex. Uses ripgrep if available, otherwise a built-in walker.",
		schema: objectSchema(map[string]any{
			"pattern": prop("string", "Regular expression to search for."),
			"path":    prop("string", "Directory to search in. Defaults to the workdir."),
		}, "pattern"),
		env: env,
		run: execGrep,
	}
}

func execGlob(_ context.Context, raw json.RawMessage, env *AgentToolContext) (string, bool) {
	var in anthropic.BetaManagedAgentsAgentToolset20260401GlobInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return errorf("invalid glob input: %v", err)
	}
	if in.Pattern == "" {
		return errorf("glob: pattern is required")
	}

	root := env.Workdir
	pattern := in.Pattern
	if filepath.IsAbs(pattern) {
		if !env.UnrestrictedPaths {
			return errorf("glob: absolute pattern not permitted")
		}
		root = "/"
		pattern = strings.TrimPrefix(pattern, "/")
	} else if in.Path != "" {
		p, err := resolvePath(env, in.Path)
		if err != nil {
			return errorf("glob: %v", err)
		}
		root = p
	}

	// Reject a ".." segment in the pattern itself. The WalkDir below is rooted
	// at the (confined) root and matches against paths relative to it, so a
	// "../.." pattern matches nothing today — but rejecting it outright keeps
	// the confinement explicit and consistent with the other SDKs' glob tools,
	// which feed the pattern to a filesystem globber where ".." would escape
	// the workdir.
	if !env.UnrestrictedPaths && hasParentDirSegment(pattern) {
		return errorf("glob: pattern %q must not contain a %q segment", pattern, "..")
	}

	type entry struct {
		path  string
		mtime int64
	}
	// Walk the tree ourselves (stdlib only — no third-party glob dependency)
	// and match each entry against the pattern. filepath.WalkDir never
	// follows symlinks, so the walk cannot escape root. We stop after
	// walkMaxEntries visited so a pattern over an enormous tree can't stall
	// the runner.
	var entries []entry
	visited := 0
	errStop := errors.New("stop")
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if p == root {
			return nil
		}
		visited++
		if visited > walkMaxEntries {
			return errStop
		}
		rel, relErr := filepath.Rel(root, p)
		if relErr != nil {
			return nil
		}
		if !globMatch(pattern, filepath.ToSlash(rel)) {
			return nil
		}
		var mt int64
		if info, statErr := d.Info(); statErr == nil {
			mt = info.ModTime().UnixNano()
		}
		entries = append(entries, entry{path: p, mtime: mt})
		return nil
	})
	if err != nil && !errors.Is(err, errStop) {
		return errorf("glob: %v", err)
	}
	if len(entries) == 0 {
		return "no matches", false
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].mtime > entries[j].mtime })
	if len(entries) > globResultLimit {
		entries = entries[:globResultLimit]
	}
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.path
	}
	return strings.Join(out, "\n"), false
}

// hasParentDirSegment reports whether pattern contains a ".." path segment,
// which could escape the workdir when a pattern is matched against the
// filesystem.
func hasParentDirSegment(pattern string) bool {
	for _, seg := range strings.Split(filepath.ToSlash(pattern), "/") {
		if seg == ".." {
			return true
		}
	}
	return false
}

// globMatch reports whether rel — a slash-separated path relative to the search
// root — matches pattern. A "**" segment matches any number of path segments
// (including none); every other segment is matched against a single path
// segment with [filepath.Match] semantics (* ? [..]). This is the small subset
// of doublestar behaviour the glob tool actually documents, implemented with
// the standard library so the SDK carries no third-party glob dependency.
func globMatch(pattern, rel string) bool {
	return matchGlobSegments(strings.Split(pattern, "/"), strings.Split(rel, "/"))
}

// matchGlobSegments matches pat against name, treating "**" as zero or more
// path segments; memoized to avoid exponential backtracking.
func matchGlobSegments(pat, name []string) bool {
	return matchGlobSegmentsMemo(pat, name, make(map[[2]int]bool))
}

func matchGlobSegmentsMemo(pat, name []string, memo map[[2]int]bool) bool {
	key := [2]int{len(pat), len(name)}
	if v, ok := memo[key]; ok {
		return v
	}
	res := func() bool {
		for len(pat) > 0 {
			if pat[0] == "**" {
				// Collapse consecutive "**" segments — they are equivalent.
				for len(pat) > 0 && pat[0] == "**" {
					pat = pat[1:]
				}
				if len(pat) == 0 {
					return true // trailing "**" matches whatever remains
				}
				// Try matching the rest of the pattern at every suffix of name.
				for i := 0; i <= len(name); i++ {
					if matchGlobSegmentsMemo(pat, name[i:], memo) {
						return true
					}
				}
				return false
			}
			if len(name) == 0 {
				return false
			}
			ok, err := filepath.Match(pat[0], name[0])
			if err != nil || !ok {
				return false
			}
			pat, name = pat[1:], name[1:]
		}
		return len(name) == 0
	}()
	memo[key] = res
	return res
}

func execGrep(ctx context.Context, raw json.RawMessage, env *AgentToolContext) (string, bool) {
	var in anthropic.BetaManagedAgentsAgentToolset20260401GrepInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return errorf("invalid grep input: %v", err)
	}
	if in.Pattern == "" {
		return errorf("grep: pattern is required")
	}
	searchPath := env.Workdir
	if in.Path != "" {
		p, err := resolvePath(env, in.Path)
		if err != nil {
			return errorf("grep: %v", err)
		}
		searchPath = p
	}

	if rg, err := exec.LookPath("rg"); err == nil {
		return runRipgrep(ctx, rg, in.Pattern, searchPath)
	}
	return runWalkGrep(in.Pattern, searchPath)
}

func runRipgrep(ctx context.Context, rg, pattern, path string) (string, bool) {
	// --max-filesize bounds per-file size; cappedBuffer bounds total stdout.
	cmd := exec.CommandContext(ctx, rg, "-n", "--no-heading",
		"--max-filesize", fmt.Sprintf("%d", grepMaxFileBytes),
		"-e", pattern, "--", path)
	stdout := newCappedBuffer(grepOutputLimit)
	var stderr bytes.Buffer
	cmd.Stdout = stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return "no matches", false
	}
	if err != nil {
		return errorf("grep: rg failed: %v: %s", err, stderr.String())
	}
	out := stdout.String()
	if stdout.truncated {
		out += "\n" + truncationNotice
	}
	return out, false
}

// cappedBuffer is an io.Writer that accumulates at most cap bytes and
// discards the rest.
type cappedBuffer struct {
	buf       bytes.Buffer
	cap       int
	truncated bool
}

func newCappedBuffer(capacity int) *cappedBuffer { return &cappedBuffer{cap: capacity} }

func (w *cappedBuffer) Write(p []byte) (int, error) {
	remaining := w.cap - w.buf.Len()
	if remaining <= 0 {
		w.truncated = true
		return len(p), nil // pretend we consumed it so the producer keeps going
	}
	if len(p) > remaining {
		w.buf.Write(p[:remaining])
		w.truncated = true
		return len(p), nil
	}
	return w.buf.Write(p)
}

func (w *cappedBuffer) String() string { return w.buf.String() }

func runWalkGrep(pattern, root string) (string, bool) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return errorf("grep: invalid regex: %v", err)
	}
	var hits []string
	budget := grepOutputLimit
	push := func(line string) bool {
		budget -= len(line) + 1
		if budget < 0 {
			hits = append(hits, truncationNotice)
			return false
		}
		hits = append(hits, line)
		return true
	}
	seen := 0
	errStop := errors.New("stop")
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		// Skip symlinks entirely: a symlink inside the workdir pointing at,
		// say, /etc/shadow must not be read (or descended into) by grep.
		// WalkDir never follows symlinks for traversal, but it still hands
		// the link entry to the callback as a non-dir, so os.ReadFile below
		// would follow it. Only descend real dirs, only read real files.
		if d.Type()&fs.ModeSymlink != 0 {
			return nil
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.Type().IsRegular() {
			return nil // skip FIFOs, devices, sockets
		}
		seen++
		if seen > walkMaxEntries {
			return errStop
		}
		// Stat before reading so a multi-GB file can't OOM the runner. Go's
		// regexp is RE2 (linear-time) so there is no ReDoS line-length concern.
		if info, err := d.Info(); err != nil || info.Size() > grepMaxFileBytes {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		sniff := min(len(data), 512)
		if bytes.IndexByte(data[:sniff], 0) >= 0 {
			return nil
		}
		for i, line := range strings.Split(string(data), "\n") {
			if re.MatchString(line) && !push(fmt.Sprintf("%s:%d:%s", path, i+1, line)) {
				return errStop
			}
		}
		return nil
	})
	if len(hits) == 0 {
		return "no matches", false
	}
	return strings.Join(hits, "\n"), false
}
