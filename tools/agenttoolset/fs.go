package agenttoolset

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	anthropic "github.com/anthropics/anthropic-sdk-go/v2"
)

// defaultMaxFileBytes is the read/edit size cap used when
// AgentToolContext.MaxFileBytes is unset (zero).
const defaultMaxFileBytes = 256 * 1024

// defaultMaxImageBase64Bytes bounds the base64-encoded form; defaultMaxPDFBytes
// bounds raw PDF bytes before base64 encoding.
const (
	defaultMaxImageBase64Bytes = 5 * 1024 * 1024
	defaultMaxPDFBytes         = 20 * 1024 * 1024
	pdfMediaType               = "application/pdf"
)

// binaryMediaTypes maps lowercase extensions to the media types read returns as
// content blocks; a test pins the values to the generated source-param types.
var binaryMediaTypes = map[string]string{
	".png":  string(anthropic.BetaBase64ImageSourceMediaTypeImagePNG),
	".jpg":  string(anthropic.BetaBase64ImageSourceMediaTypeImageJPEG),
	".jpeg": string(anthropic.BetaBase64ImageSourceMediaTypeImageJPEG),
	".gif":  string(anthropic.BetaBase64ImageSourceMediaTypeImageGIF),
	".webp": string(anthropic.BetaBase64ImageSourceMediaTypeImageWebP),
	".pdf":  pdfMediaType,
}

// readChunkBytes is how much readRangeStreaming reads from the file at a time.
const readChunkBytes = 64 * 1024

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
// env.Workdir or an allowed root.
func BetaReadTool(env *AgentToolContext) anthropic.BetaTool {
	return &funcTool{
		name:        "read",
		description: "Read a UTF-8 text file rooted at the workdir.",
		schema: objectSchema(map[string]any{
			"file_path": prop("string", "Path of the file to read, rooted at the workdir (absolute paths inside it are allowed)."),
			"view_range": map[string]any{
				"type":        "array",
				"items":       map[string]any{"type": "integer"},
				"description": "[start_line, end_line] 1-indexed inclusive",
			},
		}, "file_path"),
		env:       env,
		runBlocks: execRead,
	}
}

// BetaWriteTool returns an anthropic.BetaTool that writes file contents under
// env.Workdir or an allowed root, creating parent directories as needed.
func BetaWriteTool(env *AgentToolContext) anthropic.BetaTool {
	return &funcTool{
		name:        "write",
		description: "Write a UTF-8 text file rooted at the workdir, creating parent directories as needed.",
		schema: objectSchema(map[string]any{
			"file_path": prop("string", "Path of the file to write, rooted at the workdir (absolute paths inside it are allowed)."),
			"content":   prop("string", "Full file contents to write."),
		}, "file_path", "content"),
		env: env,
		run: execWrite,
	}
}

// BetaEditTool returns an anthropic.BetaTool that performs unique-match string
// replacement in a file under env.Workdir or an allowed root.
func BetaEditTool(env *AgentToolContext) anthropic.BetaTool {
	return &funcTool{
		name:        "edit",
		description: "Replace a unique occurrence of old_string with new_string in a file (set replace_all to replace every occurrence).",
		schema: objectSchema(map[string]any{
			"file_path":   prop("string", "Path of the file to edit, rooted at the workdir (absolute paths inside it are allowed)."),
			"old_string":  prop("string", "Substring to find and replace."),
			"new_string":  prop("string", "Replacement text."),
			"replace_all": prop("boolean", "Replace every occurrence instead of requiring a unique match."),
		}, "file_path", "old_string", "new_string"),
		env: env,
		run: execEdit,
	}
}

func toolErrorf(format string, a ...any) error {
	return &ToolError{Content: fmt.Sprintf(format, a...)}
}

func notUTF8ReadError(filePath string) error {
	return toolErrorf("read: %s: not valid UTF-8 text (binary files are only supported for image/PDF extensions)", filePath)
}

// readBinaryBlock applies the media size cap (tightened by an explicit positive
// MaxFileBytes) before reading path into a base64 image or document block.
func readBinaryBlock(path, filePath string, size int64, mediaType string, env *AgentToolContext) ([]anthropic.BetaToolResultBlockParamContentUnion, error) {
	var (
		limit  int64
		capped bool
	)
	if mediaType == pdfMediaType {
		limit, capped = resolveMaxBytes(env.MaxPDFBytes, defaultMaxPDFBytes)
	} else {
		if b64Cap, ok := resolveMaxBytes(env.MaxImageBase64Bytes, defaultMaxImageBase64Bytes); ok {
			limit, capped = (b64Cap/4)*3, true
		}
	}
	if env.MaxFileBytes > 0 && (!capped || env.MaxFileBytes < limit) {
		limit, capped = env.MaxFileBytes, true
	}
	if capped && size > limit {
		return nil, toolErrorf("read: %s is %d bytes, exceeds %d-byte limit for image/PDF files.", filePath, size, limit)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, toolErrorf("read %s: %s", filePath, fsErrorMessage(err))
	}
	encoded := base64.StdEncoding.EncodeToString(raw)
	if mediaType == pdfMediaType {
		return []anthropic.BetaToolResultBlockParamContentUnion{{
			OfDocument: &anthropic.BetaRequestDocumentBlockParam{
				Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
					OfBase64: &anthropic.BetaBase64PDFSourceParam{Data: encoded},
				},
			},
		}}, nil
	}
	return []anthropic.BetaToolResultBlockParamContentUnion{{
		OfImage: &anthropic.BetaImageBlockParam{
			Source: anthropic.BetaImageBlockParamSourceUnion{
				OfBase64: &anthropic.BetaBase64ImageSourceParam{
					Data:      encoded,
					MediaType: anthropic.BetaBase64ImageSourceMediaType(mediaType),
				},
			},
		},
	}}, nil
}

func execRead(_ context.Context, raw json.RawMessage, env *AgentToolContext) ([]anthropic.BetaToolResultBlockParamContentUnion, error) {
	var in anthropic.BetaManagedAgentsAgentToolset20260401ReadInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return nil, toolErrorf("invalid read input: %v", err)
	}
	if in.FilePath == "" {
		return nil, toolErrorf("read: file_path is required")
	}
	path, err := resolvePath(env, in.FilePath)
	if err != nil {
		return nil, toolErrorf("read: %v", err)
	}
	if len(in.ViewRange) != 0 && len(in.ViewRange) != 2 {
		return nil, toolErrorf("read: view_range must be [start_line, end_line]")
	}
	// Stat before any open: the size cap stops a multi-GB file from OOM'ing
	// the runner, and the mode check rejects FIFOs/devices/dirs before
	// open() can block on them.
	info, err := os.Stat(path)
	if err != nil {
		return nil, toolErrorf("read %s: %s", in.FilePath, fsErrorMessage(err))
	}
	if !info.Mode().IsRegular() {
		return nil, toolErrorf("read: %s is not a regular file", in.FilePath)
	}
	if mediaType, ok := binaryMediaTypes[strings.ToLower(filepath.Ext(path))]; ok {
		if len(in.ViewRange) != 0 {
			return nil, toolErrorf("read: view_range is not supported for image/PDF files")
		}
		return readBinaryBlock(path, in.FilePath, info.Size(), mediaType, env)
	}
	if limit, capped := resolveMaxBytes(env.MaxFileBytes, defaultMaxFileBytes); capped && info.Size() > limit {
		if len(in.ViewRange) == 0 {
			return nil, toolErrorf("read: %s is %d bytes, exceeds %d-byte limit. Use the view_range parameter to read specific line ranges, e.g. view_range: [1, 500].",
				in.FilePath, info.Size(), limit)
		}
		out, err := readRangeStreaming(path, in.FilePath, in.ViewRange[0], in.ViewRange[1], limit)
		if err != nil {
			return nil, err
		}
		return textResult(out), nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, toolErrorf("read %s: %s", in.FilePath, fsErrorMessage(err))
	}
	if !utf8.Valid(data) {
		return nil, notUTF8ReadError(in.FilePath)
	}
	if len(in.ViewRange) == 0 {
		return textResult(string(data)), nil
	}
	lines := strings.Split(string(data), "\n")
	start := max(0, int(in.ViewRange[0])-1)
	if start >= len(lines) {
		return textResult(""), nil
	}
	end := len(lines)
	if endLine := int(in.ViewRange[1]); endLine > 0 && endLine < end {
		end = endLine
	}
	// An inverted range selects nothing; without this lines[start:end] panics.
	if end < start {
		return textResult(""), nil
	}
	return textResult(strings.Join(lines[start:end], "\n")), nil
}

// readRangeStreaming returns lines [startLine, endLine] of the file at path,
// capping the selected bytes at limit.
func readRangeStreaming(path, filePath string, startLine, endLine, limit int64) (string, error) {
	lines := newLineRangeCollector(filePath, startLine, endLine, limit)
	if lines.rangeIsEmpty() {
		return "", nil
	}
	f, err := os.Open(path)
	if err != nil {
		return "", toolErrorf("read %s: %s", filePath, fsErrorMessage(err))
	}
	defer f.Close()
	// Raw byte chunks rather than bufio.Scanner or ReadString: a single huge
	// line must never be buffered whole, so memory stays bounded by limit plus
	// one chunk.
	chunk := make([]byte, readChunkBytes)
	for {
		n, readErr := f.Read(chunk)
		if err := lines.collectFrom(chunk[:n]); err != nil {
			return "", err
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return "", toolErrorf("read %s: %s", filePath, fsErrorMessage(readErr))
		}
		if lines.rangeIsCollected() {
			break
		}
	}
	if !utf8.Valid(lines.collected) {
		return "", notUTF8ReadError(filePath)
	}
	return lines.text(), nil
}

// lineRangeCollector collects the bytes of lines [startLine, endLine] from
// consecutive file chunks, capped at limit.
type lineRangeCollector struct {
	filePath           string
	startLine, endLine int64
	start, end         int
	limit              int64
	line               int
	collected          []byte
}

func newLineRangeCollector(filePath string, startLine, endLine, limit int64) *lineRangeCollector {
	end := math.MaxInt
	if endLine > 0 {
		end = int(endLine)
	}
	return &lineRangeCollector{
		filePath:  filePath,
		startLine: startLine,
		endLine:   endLine,
		start:     max(0, int(startLine)-1),
		end:       end,
		limit:     limit,
	}
}

func (c *lineRangeCollector) rangeIsEmpty() bool {
	return c.end <= c.start
}

func (c *lineRangeCollector) rangeIsCollected() bool {
	return c.line >= c.end
}

func (c *lineRangeCollector) collectFrom(chunk []byte) error {
	lineStart := 0
	for lineStart < len(chunk) && !c.rangeIsCollected() {
		newline := bytes.IndexByte(chunk[lineStart:], '\n')
		lineEnd := len(chunk)
		if newline >= 0 {
			newline += lineStart
			lineEnd = newline
		}
		if c.line >= c.start {
			if err := c.collect(chunk[lineStart:lineEnd], newline >= 0); err != nil {
				return err
			}
		}
		if newline < 0 {
			break
		}
		c.line++
		lineStart = newline + 1
	}
	return nil
}

func (c *lineRangeCollector) collect(lineBytes []byte, newlineTerminated bool) error {
	c.collected = append(c.collected, lineBytes...)
	if newlineTerminated && c.line+1 < c.end {
		c.collected = append(c.collected, '\n')
	}
	if int64(len(c.collected)) > c.limit {
		return c.overLimitError()
	}
	return nil
}

func (c *lineRangeCollector) overLimitError() error {
	if c.end-c.start == 1 {
		return &ToolError{Content: fmt.Sprintf("read: line %d of %s alone exceeds %d-byte limit. The read tool cannot return part of a line, so view_range cannot narrow this further.",
			c.start+1, c.filePath, c.limit)}
	}
	return &ToolError{Content: fmt.Sprintf("read: view_range [%d, %d] of %s exceeds %d-byte limit. Narrow the view_range to read a smaller portion.",
		c.startLine, c.endLine, c.filePath, c.limit)}
}

func (c *lineRangeCollector) text() string {
	return string(c.collected)
}

func execWrite(_ context.Context, raw json.RawMessage, env *AgentToolContext) (string, bool) {
	var in anthropic.BetaManagedAgentsAgentToolset20260401WriteInput
	if err := json.Unmarshal(raw, &in); err != nil {
		return errorf("invalid write input: %v", err)
	}
	if in.FilePath == "" {
		return errorf("write: file_path is required")
	}
	path, err := resolveWritablePath(env, in.FilePath)
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
	path, err := resolveWritablePath(env, in.FilePath)
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
		return errorf("edit: %s is %d bytes, exceeds %d-byte limit. The edit tool loads the whole file and cannot modify a file this large.",
			in.FilePath, info.Size(), limit)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return errorf("edit %s: %s", in.FilePath, fsErrorMessage(err))
	}
	if !utf8.Valid(data) {
		return errorf("edit: %s: not valid UTF-8 text (cannot edit binary files)", in.FilePath)
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
