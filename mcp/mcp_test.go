package mcp_test

import (
	"errors"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/mcp"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestToBlock_Text(t *testing.T) {
	block, err := mcp.ToBlock(&mcpsdk.TextContent{Text: "hello world"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.OfText == nil || block.OfText.Text != "hello world" {
		t.Fatalf("expected text block with %q", "hello world")
	}
}

func TestToBlock_Image(t *testing.T) {
	for _, mimeType := range []string{"image/jpeg", "image/png", "image/gif", "image/webp"} {
		t.Run(mimeType, func(t *testing.T) {
			block, err := mcp.ToBlock(&mcpsdk.ImageContent{
				Data:     []byte("raw bytes"),
				MIMEType: mimeType,
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if block.OfImage == nil || block.OfImage.Source.OfBase64 == nil {
				t.Fatal("expected base64 image block")
			}
			if string(block.OfImage.Source.OfBase64.MediaType) != mimeType {
				t.Fatalf("expected media_type %q, got %q", mimeType, block.OfImage.Source.OfBase64.MediaType)
			}
		})
	}
}

func TestToBlock_ImageUnsupportedMimeType(t *testing.T) {
	_, err := mcp.ToBlock(&mcpsdk.ImageContent{MIMEType: "image/bmp"})
	var mcpErr *mcp.UnsupportedValueError
	if !errors.As(err, &mcpErr) {
		t.Fatalf("expected UnsupportedValueError, got %T: %v", err, err)
	}
}

func TestToBlock_EmbeddedResource_Text(t *testing.T) {
	block, err := mcp.ToBlock(&mcpsdk.EmbeddedResource{
		Resource: &mcpsdk.ResourceContents{URI: "file:///hello.txt", MIMEType: "text/plain", Text: "hello"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.OfDocument == nil || block.OfDocument.Source.OfText == nil {
		t.Fatal("expected text document block")
	}
	if block.OfDocument.Source.OfText.Data != "hello" {
		t.Fatalf("expected %q, got %q", "hello", block.OfDocument.Source.OfText.Data)
	}
}

func TestToBlock_EmbeddedResource_Image(t *testing.T) {
	block, err := mcp.ToBlock(&mcpsdk.EmbeddedResource{
		Resource: &mcpsdk.ResourceContents{URI: "file:///img.png", MIMEType: "image/png", Blob: []byte("binary")},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.OfImage == nil {
		t.Fatal("expected image block")
	}
}

func TestToBlock_EmbeddedResource_PDF(t *testing.T) {
	block, err := mcp.ToBlock(&mcpsdk.EmbeddedResource{
		Resource: &mcpsdk.ResourceContents{URI: "file:///doc.pdf", MIMEType: "application/pdf", Blob: []byte("%PDF")},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.OfDocument == nil || block.OfDocument.Source.OfBase64 == nil {
		t.Fatal("expected base64 document block")
	}
}

func TestToBlock_Audio_Unsupported(t *testing.T) {
	_, err := mcp.ToBlock(&mcpsdk.AudioContent{})
	var mcpErr *mcp.UnsupportedValueError
	if !errors.As(err, &mcpErr) {
		t.Fatalf("expected UnsupportedValueError, got %T", err)
	}
}

func TestToBlock_ResourceLink_Unsupported(t *testing.T) {
	_, err := mcp.ToBlock(&mcpsdk.ResourceLink{})
	var mcpErr *mcp.UnsupportedValueError
	if !errors.As(err, &mcpErr) {
		t.Fatalf("expected UnsupportedValueError, got %T", err)
	}
}

func TestToBlock_NilResource(t *testing.T) {
	_, err := mcp.ToBlock(&mcpsdk.EmbeddedResource{Resource: nil})
	if err == nil {
		t.Fatal("expected error for nil resource")
	}
}

func TestToMessage_User(t *testing.T) {
	msg, err := mcp.ToMessage(&mcpsdk.PromptMessage{
		Role:    "user",
		Content: &mcpsdk.TextContent{Text: "hello"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Role != "user" {
		t.Fatalf("expected role %q, got %q", "user", msg.Role)
	}
	if len(msg.Content) != 1 || msg.Content[0].OfText == nil {
		t.Fatal("expected single text content block")
	}
}

func TestToMessage_Assistant(t *testing.T) {
	msg, err := mcp.ToMessage(&mcpsdk.PromptMessage{
		Role:    "assistant",
		Content: &mcpsdk.TextContent{Text: "reply"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Role != "assistant" {
		t.Fatalf("expected role %q, got %q", "assistant", msg.Role)
	}
}

func TestToMessage_UnsupportedContent(t *testing.T) {
	_, err := mcp.ToMessage(&mcpsdk.PromptMessage{
		Role:    "user",
		Content: &mcpsdk.AudioContent{},
	})
	if err == nil {
		t.Fatal("expected error for unsupported content")
	}
}

func TestResourceToBlock_Text(t *testing.T) {
	block, err := mcp.ResourceToBlock(&mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{
			{URI: "file:///hello.txt", MIMEType: "text/plain", Text: "hello world"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.OfDocument == nil || block.OfDocument.Source.OfText == nil {
		t.Fatal("expected text document block")
	}
	if block.OfDocument.Source.OfText.Data != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", block.OfDocument.Source.OfText.Data)
	}
}

func TestResourceToBlock_NoMimeType(t *testing.T) {
	block, err := mcp.ResourceToBlock(&mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{{URI: "file:///data", Text: "plain data"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.OfDocument == nil || block.OfDocument.Source.OfText == nil {
		t.Fatal("expected text document block")
	}
}

func TestResourceToBlock_SkipsUnsupportedSelectsSupported(t *testing.T) {
	block, err := mcp.ResourceToBlock(&mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{
			{URI: "file:///audio.mp3", MIMEType: "audio/mpeg"},
			{URI: "file:///notes.txt", MIMEType: "text/plain", Text: "notes"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.OfDocument == nil || block.OfDocument.Source.OfText == nil || block.OfDocument.Source.OfText.Data != "notes" {
		t.Fatal("expected text from second item")
	}
}

func TestResourceToBlock_EmptyContents(t *testing.T) {
	_, err := mcp.ResourceToBlock(&mcpsdk.ReadResourceResult{})
	var mcpErr *mcp.UnsupportedValueError
	if !errors.As(err, &mcpErr) {
		t.Fatalf("expected UnsupportedValueError, got %T", err)
	}
}

func TestResourceToBlock_NoSupportedMimeType(t *testing.T) {
	_, err := mcp.ResourceToBlock(&mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{{URI: "file:///video.mp4", MIMEType: "video/mp4"}},
	})
	var mcpErr *mcp.UnsupportedValueError
	if !errors.As(err, &mcpErr) {
		t.Fatalf("expected UnsupportedValueError, got %T", err)
	}
}

func TestResourceToBlock_ImageWithTextData_Error(t *testing.T) {
	_, err := mcp.ResourceToBlock(&mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{
			{URI: "file:///img.png", MIMEType: "image/png", Text: "not binary"},
		},
	})
	if err == nil {
		t.Fatal("expected error for image resource with text data")
	}
}

func TestResourceToBlock_BlobIsRawBytes(t *testing.T) {
	// Blob is []byte (already decoded) — no base64 required from the caller.
	block, err := mcp.ResourceToBlock(&mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{
			{URI: "file:///file.txt", MIMEType: "text/plain", Blob: []byte("decoded text")},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.OfDocument == nil || block.OfDocument.Source.OfText == nil {
		t.Fatal("expected text document")
	}
	if block.OfDocument.Source.OfText.Data != "decoded text" {
		t.Fatalf("expected %q, got %q", "decoded text", block.OfDocument.Source.OfText.Data)
	}
}

func TestResourceToFile_Text(t *testing.T) {
	r, err := mcp.ResourceToFile(&mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{
			{URI: "file:///hello.txt", MIMEType: "text/plain", Text: "hello"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil reader")
	}
}

func TestResourceToFile_EmptyContents(t *testing.T) {
	_, err := mcp.ResourceToFile(&mcpsdk.ReadResourceResult{})
	var mcpErr *mcp.UnsupportedValueError
	if !errors.As(err, &mcpErr) {
		t.Fatalf("expected UnsupportedValueError, got %T", err)
	}
}
