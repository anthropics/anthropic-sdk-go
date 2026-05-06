// Package mcp provides helpers for using Model Context Protocol tools and
// resources with the Anthropic SDK.
//
// This package is imported separately from the core SDK so that consumers who
// don't use MCP don't pull in the MCP SDK and its transitive dependencies.
package mcp

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// UnsupportedValueError is returned when an MCP value cannot be converted to
// a format supported by the Claude API.
type UnsupportedValueError struct {
	Message string
}

func (e *UnsupportedValueError) Error() string { return e.Message }

// -----------------------------------------------------------------------------
// Tool conversion
// -----------------------------------------------------------------------------

// NewBetaTool creates an [anthropic.BetaTool] from an MCP tool and session
// for use with [anthropic.BetaMessageService.NewToolRunner].
//
// The session must outlive any [anthropic.BetaToolRunner] that uses the
// returned tool, since execution dispatches to the session.
func NewBetaTool(tool *mcpsdk.Tool, session *mcpsdk.ClientSession) (anthropic.BetaTool, error) {
	var inputSchema anthropic.BetaToolInputSchemaParam
	if tool.InputSchema != nil {
		// Round-trip via JSON: the MCP and Anthropic schema types use different
		// representations internally, but their wire formats are compatible.
		b, err := json.Marshal(tool.InputSchema)
		if err == nil {
			err = json.Unmarshal(b, &inputSchema)
		}
		if err != nil {
			return nil, err
		}
	}
	return &betaTool{tool: tool, session: session, schema: inputSchema}, nil
}

// NewBetaTools converts a slice of MCP tools into [anthropic.BetaTool] values
// for the tool runner. The session must outlive any runner that uses these
// tools.
//
//	result, _ := session.ListTools(ctx, nil)
//	tools, err := mcp.NewBetaTools(result.Tools, session)
//	runner := client.Beta.Messages.NewToolRunner(tools, params)
func NewBetaTools(tools []*mcpsdk.Tool, session *mcpsdk.ClientSession) ([]anthropic.BetaTool, error) {
	out := make([]anthropic.BetaTool, 0, len(tools))
	for _, t := range tools {
		item, err := NewBetaTool(t, session)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

type betaTool struct {
	tool    *mcpsdk.Tool
	session *mcpsdk.ClientSession
	schema  anthropic.BetaToolInputSchemaParam
}

func (t *betaTool) Name() string                                    { return t.tool.Name }
func (t *betaTool) Description() string                             { return t.tool.Description }
func (t *betaTool) InputSchema() anthropic.BetaToolInputSchemaParam { return t.schema }

func (t *betaTool) Execute(ctx context.Context, input json.RawMessage) ([]anthropic.BetaToolResultBlockParamContentUnion, error) {
	var args map[string]any
	if err := json.Unmarshal(input, &args); err != nil {
		return nil, fmt.Errorf("mcp tool %s: failed to unmarshal input: %w", t.tool.Name, err)
	}

	result, err := t.session.CallTool(ctx, &mcpsdk.CallToolParams{
		Name:      t.tool.Name,
		Arguments: args,
	})
	if err != nil {
		return nil, fmt.Errorf("mcp tool %s: %w", t.tool.Name, err)
	}

	if result.IsError {
		var parts []string
		for _, c := range result.Content {
			if tc, ok := c.(*mcpsdk.TextContent); ok && tc.Text != "" {
				parts = append(parts, tc.Text)
			}
		}
		msg := strings.Join(parts, "\n")
		if msg == "" {
			msg = "tool returned an error"
		}
		return nil, errors.New(msg)
	}

	// Per the MCP spec, when both Content and StructuredContent are present,
	// Content is a text mirror of StructuredContent and either is sufficient.
	// When Content is empty but StructuredContent is set, encode the structured
	// data as a text block — matching the TS and Python implementations.
	if len(result.Content) == 0 {
		if result.StructuredContent != nil {
			b, marshalErr := json.Marshal(result.StructuredContent)
			if marshalErr != nil {
				return nil, fmt.Errorf("mcp tool %s: failed to marshal structured content: %w", t.tool.Name, marshalErr)
			}
			return []anthropic.BetaToolResultBlockParamContentUnion{
				{OfText: &anthropic.BetaTextBlockParam{Text: string(b)}},
			}, nil
		}
		return nil, nil
	}

	blocks := make([]anthropic.BetaToolResultBlockParamContentUnion, 0, len(result.Content))
	for _, c := range result.Content {
		block, convErr := ToBlock(c)
		if convErr != nil {
			return nil, convErr
		}
		blocks = append(blocks, block)
	}
	return blocks, nil
}

// -----------------------------------------------------------------------------
// Content conversion
// -----------------------------------------------------------------------------

// ToBlock converts a single MCP content item into an Anthropic content block.
// Supported: TextContent, ImageContent (jpeg/png/gif/webp), EmbeddedResource.
// Returns [*UnsupportedValueError] for AudioContent and ResourceLink.
func ToBlock(content mcpsdk.Content) (anthropic.BetaToolResultBlockParamContentUnion, error) {
	switch v := content.(type) {
	case *mcpsdk.TextContent:
		return anthropic.BetaToolResultBlockParamContentUnion{OfText: &anthropic.BetaTextBlockParam{Text: v.Text}}, nil

	case *mcpsdk.ImageContent:
		if !isSupportedImageMimeType(v.MIMEType) {
			return anthropic.BetaToolResultBlockParamContentUnion{}, &UnsupportedValueError{
				fmt.Sprintf("unsupported image MIME type: %s", v.MIMEType),
			}
		}
		return anthropic.BetaToolResultBlockParamContentUnion{OfImage: &anthropic.BetaImageBlockParam{
			Source: anthropic.BetaImageBlockParamSourceUnion{
				OfBase64: &anthropic.BetaBase64ImageSourceParam{
					Data:      base64.StdEncoding.EncodeToString(v.Data),
					MediaType: anthropic.BetaBase64ImageSourceMediaType(v.MIMEType),
				},
			},
		}}, nil

	case *mcpsdk.EmbeddedResource:
		if v.Resource == nil {
			return anthropic.BetaToolResultBlockParamContentUnion{}, &UnsupportedValueError{"embedded resource has nil contents"}
		}
		return convertResourceContents(v.Resource)

	case *mcpsdk.AudioContent, *mcpsdk.ResourceLink:
		return anthropic.BetaToolResultBlockParamContentUnion{}, &UnsupportedValueError{
			fmt.Sprintf("unsupported MCP content type: %T", content),
		}

	default:
		return anthropic.BetaToolResultBlockParamContentUnion{}, &UnsupportedValueError{
			fmt.Sprintf("unknown MCP content type: %T", content),
		}
	}
}

// ToMessage converts an MCP prompt message into a [anthropic.BetaMessageParam].
func ToMessage(msg *mcpsdk.PromptMessage) (anthropic.BetaMessageParam, error) {
	block, err := ToBlock(msg.Content)
	if err != nil {
		return anthropic.BetaMessageParam{}, err
	}
	return anthropic.BetaMessageParam{
		Role: anthropic.BetaMessageParamRole(msg.Role),
		Content: []anthropic.BetaContentBlockParamUnion{{
			OfText:     block.OfText,
			OfImage:    block.OfImage,
			OfDocument: block.OfDocument,
		}},
	}, nil
}

// -----------------------------------------------------------------------------
// Resource conversion
// -----------------------------------------------------------------------------

// ResourceToBlock converts MCP resource read results into an Anthropic content
// block. Returns the first item in Contents with a MIME type supported by the
// Claude API.
func ResourceToBlock(result *mcpsdk.ReadResourceResult) (anthropic.BetaToolResultBlockParamContentUnion, error) {
	if len(result.Contents) == 0 {
		return anthropic.BetaToolResultBlockParamContentUnion{}, &UnsupportedValueError{
			"resource contents array must contain at least one item",
		}
	}
	for _, c := range result.Contents {
		if isSupportedResourceMimeType(c.MIMEType) {
			return convertResourceContents(c)
		}
	}
	var mimeTypes []string
	for _, c := range result.Contents {
		if c.MIMEType != "" {
			mimeTypes = append(mimeTypes, c.MIMEType)
		}
	}
	return anthropic.BetaToolResultBlockParamContentUnion{}, &UnsupportedValueError{
		fmt.Sprintf("no supported MIME type found in resource contents. Available: %s", strings.Join(mimeTypes, ", ")),
	}
}

// ResourceToFile converts MCP resource contents into an [io.Reader] suitable
// for [anthropic.BetaFileUploadParams].File. Uses the first item in Contents
// regardless of MIME type — any resource can be uploaded as a file.
func ResourceToFile(result *mcpsdk.ReadResourceResult) (io.Reader, error) {
	if len(result.Contents) == 0 {
		return nil, &UnsupportedValueError{"resource contents array must contain at least one item"}
	}
	res := result.Contents[0]
	var data []byte
	if len(res.Blob) > 0 {
		data = res.Blob
	} else {
		data = []byte(res.Text)
	}
	return anthropic.File(bytes.NewReader(data), extractFilename(res.URI), res.MIMEType), nil
}

// -----------------------------------------------------------------------------
// Internal helpers
// -----------------------------------------------------------------------------

func convertResourceContents(res *mcpsdk.ResourceContents) (anthropic.BetaToolResultBlockParamContentUnion, error) {
	mimeType := res.MIMEType

	if isSupportedImageMimeType(mimeType) {
		if len(res.Blob) == 0 {
			return anthropic.BetaToolResultBlockParamContentUnion{}, &UnsupportedValueError{
				fmt.Sprintf("image resource must have blob data, not text. URI: %s", res.URI),
			}
		}
		return anthropic.BetaToolResultBlockParamContentUnion{OfImage: &anthropic.BetaImageBlockParam{
			Source: anthropic.BetaImageBlockParamSourceUnion{
				OfBase64: &anthropic.BetaBase64ImageSourceParam{
					Data:      base64.StdEncoding.EncodeToString(res.Blob),
					MediaType: anthropic.BetaBase64ImageSourceMediaType(mimeType),
				},
			},
		}}, nil
	}

	if mimeType == "application/pdf" {
		if len(res.Blob) == 0 {
			return anthropic.BetaToolResultBlockParamContentUnion{}, &UnsupportedValueError{
				fmt.Sprintf("PDF resource must have blob data, not text. URI: %s", res.URI),
			}
		}
		return anthropic.BetaToolResultBlockParamContentUnion{OfDocument: &anthropic.BetaRequestDocumentBlockParam{
			Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
				OfBase64: &anthropic.BetaBase64PDFSourceParam{Data: base64.StdEncoding.EncodeToString(res.Blob)},
			},
		}}, nil
	}

	if mimeType == "" || strings.HasPrefix(mimeType, "text/") {
		var text string
		if len(res.Blob) > 0 {
			text = string(res.Blob)
		} else {
			text = res.Text
		}
		return anthropic.BetaToolResultBlockParamContentUnion{OfDocument: &anthropic.BetaRequestDocumentBlockParam{
			Source: anthropic.BetaRequestDocumentBlockSourceUnionParam{
				OfText: &anthropic.BetaPlainTextSourceParam{Data: text},
			},
		}}, nil
	}

	return anthropic.BetaToolResultBlockParamContentUnion{}, &UnsupportedValueError{
		fmt.Sprintf("unsupported MIME type %q for resource: %s", mimeType, res.URI),
	}
}

var supportedImageMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func isSupportedImageMimeType(mimeType string) bool {
	return supportedImageMimeTypes[mimeType]
}

func isSupportedResourceMimeType(mimeType string) bool {
	return mimeType == "" ||
		strings.HasPrefix(mimeType, "text/") ||
		mimeType == "application/pdf" ||
		isSupportedImageMimeType(mimeType)
}

func extractFilename(uri string) string {
	if uri == "" {
		return "file"
	}
	u, err := url.Parse(uri)
	if err == nil && u.Path != "" {
		if idx := strings.LastIndex(u.Path, "/"); idx >= 0 {
			if name := u.Path[idx+1:]; name != "" {
				return name
			}
		}
	}
	return "file"
}
