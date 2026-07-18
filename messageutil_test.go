package anthropic_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

func unmarshalContentBlockParam(t *testing.T, jsonData string) anthropic.ContentBlockParamUnion {
	var block anthropic.ContentBlockUnion
	err := json.Unmarshal([]byte(jsonData), &block)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	result := block.ToParam()
	return result
}

func TestContentBlockUnionToParam(t *testing.T) {
	t.Run("TextBlock with text only", func(t *testing.T) {
		result := unmarshalContentBlockParam(t, `{"type":"text","text":"Hello, world!"}`)
		if result.OfText == nil {
			t.Fatal("Expected OfText to be non-nil")
		}
		if result.OfText.Text != "Hello, world!" {
			t.Errorf("Expected text 'Hello, world!', got '%s'", result.OfText.Text)
		}
		if result.OfText.Type != constant.Text("text") {
			t.Errorf("Expected type 'text', got '%s'", result.OfText.Type)
		}
	})

	t.Run("WebSearchToolResultBlock with search results", func(t *testing.T) {
		result := unmarshalContentBlockParam(t, `{"type":"web_search_tool_result","tool_use_id":"test123","content":[{"type":"web_search_result","title":"Test Web Title","url":"https://test.com","encrypted_content":"abc123","page_age":"1 day ago"}]}`)
		var block anthropic.ContentBlockUnion
		if err := json.Unmarshal([]byte(`{"type":"web_search_tool_result","tool_use_id":"test123","content":[{"type":"web_search_result","title":"Test Web Title","url":"https://test.com","encrypted_content":"abc123","page_age":"1 day ago"}]}`), &block); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if len(block.Content.OfWebSearchResultBlockArray) != 1 {
			t.Errorf("Expected Content.OfWebSearchResultBlockArray to have 1 result, got %d", len(block.Content.OfWebSearchResultBlockArray))
		}
		if len(block.Content.OfWebSearchResultBlockArray) > 0 && block.Content.OfWebSearchResultBlockArray[0].Title != "Test Web Title" {
			t.Errorf("Expected title '', got '%s'", block.Content.OfWebSearchResultBlockArray[0].Title)
		}

		if result.OfWebSearchToolResult == nil {
			t.Fatal("Expected OfWebSearchToolResult to be non-nil")
		}
		if len(result.OfWebSearchToolResult.Content.OfWebSearchToolResultBlockItem) != 1 {
			t.Errorf("Expected 1 search result in param, got %d", len(result.OfWebSearchToolResult.Content.OfWebSearchToolResultBlockItem))
		}
	})
}

// TestJSONOutputFormatParamTypeAlwaysMarshaled is a regression test for issue
// https://github.com/anthropics/anthropic-sdk-go/issues/328 — the structured-
// output `type` discriminator was being elided from the request body, which
// caused the API to hang. The field is documented as elidable (its zero value
// should marshal as "json_schema") via the `default:"json_schema"` struct tag,
// so request bodies must always carry it.
func TestJSONOutputFormatParamTypeAlwaysMarshaled(t *testing.T) {
	wantType := []byte(`"type":"json_schema"`)

	t.Run("zero Type emits default", func(t *testing.T) {
		body, err := json.Marshal(anthropic.JSONOutputFormatParam{
			Schema: map[string]any{"type": "object"},
		})
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if !bytes.Contains(body, wantType) {
			t.Errorf("missing %s in %s", wantType, body)
		}
	})

	t.Run("explicit Type round-trips", func(t *testing.T) {
		body, err := json.Marshal(anthropic.JSONOutputFormatParam{
			Schema: map[string]any{"type": "object"},
			Type:   constant.JSONSchema("json_schema"),
		})
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if !bytes.Contains(body, wantType) {
			t.Errorf("missing %s in %s", wantType, body)
		}
	})

	t.Run("nested in MessageNewParams", func(t *testing.T) {
		params := anthropic.MessageNewParams{
			MaxTokens: 100,
			Model:     anthropic.ModelClaudeSonnet4_5,
			Messages: []anthropic.MessageParam{
				anthropic.NewUserMessage(anthropic.NewTextBlock("hi")),
			},
			OutputConfig: anthropic.OutputConfigParam{
				Format: anthropic.JSONOutputFormatParam{
					Schema: map[string]any{"type": "object"},
				},
			},
		}
		body, err := json.Marshal(params)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if !bytes.Contains(body, wantType) {
			t.Errorf("missing %s in %s", wantType, body)
		}
		// Sanity check: the `type` we found is the one inside output_config.format,
		// not a stray copy from a nested schema.
		if !bytes.Contains(body, []byte(`"format":{"schema":`)) {
			t.Errorf("unexpected body shape: %s", body)
		}
	})
}
