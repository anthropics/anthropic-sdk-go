package agenttoolset

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	anthropic "github.com/anthropics/anthropic-sdk-go"
)

// defaultMaxFileBytes is the read/edit size cap used when
// AgentToolContext.MaxFileBytes is unset (zero).
const defaultMaxFileBytes = 256 * 1024

// resolveMaxBytes turns a configured cap into an effective size limit. Zero
// selects def (the built-in default); a negative value disables the size check
// entirely (capped == false). It governs only the size guard — callers still
// reject non-regular files, since the FIFO/device hang hazard is unrelated to
// memory headroom.
func resolveMaxBytes(configured, def int64) (limit int64, capped bool) {
	switch {
	case configured < 0:
		return 0, false
	case configured == 0:
		return def, true
	default:
		return configured, true
	}
}

// BetaReadTool returns an anthropic.BetaTool that reads file contents under
// env.Workdir.
func BetaReadTool(env *AgentToolContext) anthropic.BetaTool {
	return &funcTool{
		name:        "read",
		description: "Read a UTF-8 text file under the workdir.",
		schema: objectSchema(map[string]any{
			"file_path": prop("string", "Path of the file to read, relative to the workdir or absolute under it."),
			"view_range": map[string]any{
				"type":        "array",
				"items":       map[string]any{"type": "integer"},
				"description": "[start_line, end_line] 1-indexed inclusive",
			},
		}, "file_path"),
		env: env,
		run: execRead,
	}
}

// BetaWriteTool returns an anthropic.BetaTool that writes file contents under
// env.Workdir, creating parent directories as needed.
func BetaWriteTool(env *AgentToolContext) anthropic.BetaTool {
	return &funcTool{
		name:        "write",
		description: "Write a UTF-8 text file under the workdir, creating parent directories as needed.",
		schema: objectSchema(map[string]any{
			"file_path": prop("string", "Path of the file to write, relative to the workdir or absolute under it."),
			"content":   prop("string", "Full file contents to write."),
		}, "file_path", "content"),
		env: env,
		run: execWrite,
	}
}

// BetaEditTool returns an anthropic.BetaTool that performs unique-match string
// replacement in a file under env.Workdir.
func BetaEditTool(env *AgentToolContext) anthropic.BetaTool {
	return &funcTool{
		name:        "edit",
		description: "Replace a unique occurrence of old_string with new_string in a file (set replace_all to replace every occurrence).",
		schema: objectSchema(map[string]any{
			"file_path":   prop("string", "Path of the file to edit, relative to the workdir or absolute under it."),
			"old_string":  prop("string", "Substring to find and replace."),
			"new_string":  prop("string", "Replacement text."),
			"replace_all": prop("boolean", "Replace every occurrence instead of requiring a unique match."),
		}, "file_path", "old_string", "new_string"),
		env: env,
		run: execEdit,
	}
}

func execRead(_ context.Context, raw json.RawMessage, env *AgentToolContext) (string, bool) {
	var in anthropic.BetaManagedAgentsAgentToolset20260401ReadInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return errorf("invalid read input: %v", err)
	}
	if in.FilePath == "" {
		return errorf("read: file_path is required")
	}
	path, err := resolvePath(env, in.FilePath)
	if err != nil {
		return errorf("read: %v", err)
	}
	// Stat before any open: the size cap stops a multi-GB file from OOM'ing
	// the runner, and the mode check rejects FIFOs/devices/dirs before
	// open() can block on them.
	info, err := os.Stat(path)
	if err != nil {
		return errorf("read %s: %s", in.FilePath, fsErrorMessage(err))
	}
	if !info.Mode().IsRegular() {
		return errorf("read: %s is not a regular file", in.FilePath)
	}
	if limit, capped := resolveMaxBytes(env.MaxFileBytes, defaultMaxFileBytes); capped && info.Size() > limit {
		return errorf("read: %s is %d bytes, exceeds %d-byte limit. Use bash (head/tail/sed) to read a slice.",
			in.FilePath, info.Size(), limit)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return errorf("read %s: %s", in.FilePath, fsErrorMessage(err))
	}
	if len(in.ViewRange) == 0 {
		return string(data), false
	}
	if len(in.ViewRange) != 2 {
		return errorf("read: view_range must be [start_line, end_line]")
	}
	lines := strings.Split(string(data), "\n")
	start := max(0, int(in.ViewRange[0])-1)
	if start >= len(lines) {
		return "", false
	}
	end := len(lines)
	if endLine := int(in.ViewRange[1]); endLine > 0 && endLine < end {
		end = endLine
	}
	// Guard against an inverted range (end_line < start_line): slicing
	// lines[start:end] with end < start panics, and a confused or
	// prompt-injected model can crash the whole runner with one tool call.
	if end < start {
		return errorf("read: view_range end line %d is before start line %d", in.ViewRange[1], in.ViewRange[0])
	}
	return strings.Join(lines[start:end], "\n"), false
}

func execWrite(_ context.Context, raw json.RawMessage, env *AgentToolContext) (string, bool) {
	var in anthropic.BetaManagedAgentsAgentToolset20260401WriteInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return errorf("invalid write input: %v", err)
	}
	if in.FilePath == "" {
		return errorf("write: file_path is required")
	}
	path, err := resolvePath(env, in.FilePath)
	if err != nil {
		return errorf("write: %v", err)
	}
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return errorf("write %s: mkdir: %s", in.FilePath, fsErrorMessage(err))
		}
	}
	if err := atomicWriteFile(path, []byte(in.Content), 0o644); err != nil {
		return errorf("write %s: %s", in.FilePath, fsErrorMessage(err))
	}
	return fmt.Sprintf("wrote %d bytes to %s", len(in.Content), in.FilePath), false
}

func execEdit(_ context.Context, raw json.RawMessage, env *AgentToolContext) (string, bool) {
	var in anthropic.BetaManagedAgentsAgentToolset20260401EditInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return errorf("invalid edit input: %v", err)
	}
	if in.FilePath == "" {
		return errorf("edit: file_path is required")
	}
	if in.OldString == "" {
		return errorf("edit: old_string is required")
	}
	path, err := resolvePath(env, in.FilePath)
	if err != nil {
		return errorf("edit: %v", err)
	}
	// Stat before any open: the size cap stops a multi-GB file from OOM'ing
	// the runner, and the mode check rejects FIFOs/devices/dirs before
	// open() can block on them.
	info, err := os.Stat(path)
	if err != nil {
		return errorf("edit %s: %s", in.FilePath, fsErrorMessage(err))
	}
	if !info.Mode().IsRegular() {
		return errorf("edit: %s is not a regular file", in.FilePath)
	}
	if limit, capped := resolveMaxBytes(env.MaxFileBytes, defaultMaxFileBytes); capped && info.Size() > limit {
		return errorf("edit: %s is %d bytes, exceeds %d-byte limit. Use bash (sed/awk) to modify a large file.",
			in.FilePath, info.Size(), limit)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return errorf("edit %s: %s", in.FilePath, fsErrorMessage(err))
	}
	content := string(data)

	count := strings.Count(content, in.OldString)
	if count == 0 {
		return errorf("edit: old_string not found in %s", in.FilePath)
	}
	var updated string
	if in.ReplaceAll {
		updated = strings.ReplaceAll(content, in.OldString, in.NewString)
	} else {
		if count > 1 {
			return errorf("edit: old_string appears %d times in %s (must be unique)", count, in.FilePath)
		}
		updated = strings.Replace(content, in.OldString, in.NewString, 1)
	}
	if err := atomicWriteFile(path, []byte(updated), 0o644); err != nil {
		return errorf("edit %s: %s", in.FilePath, fsErrorMessage(err))
	}
	return fmt.Sprintf("edited %s (%d replacement(s))", in.FilePath, count), false
}

// atomicWriteFile writes data to a temp file in the destination directory and
// renames it over path, so a concurrent reader never observes a half-written
// file and a failed write leaves the original intact. The write/edit file tools
// go through this; rename is atomic only within a single filesystem, which
// holds here because the temp file is created alongside the destination.
func atomicWriteFile(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".agenttoolset-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		// Best-effort cleanup if we bail before the rename succeeds.
		_ = os.Remove(tmpName)
	}()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Chmod(perm); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}
