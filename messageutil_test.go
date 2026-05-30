package anthropic_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

func requireJSONEqual(t *testing.T, want string, got any) {
	t.Helper()

	gotJSON, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	var wantValue any
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}
	var gotValue any
	if err := json.Unmarshal(gotJSON, &gotValue); err != nil {
		t.Fatalf("Failed to unmarshal actual JSON: %v", err)
	}
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Fatalf("JSON mismatch:\n got: %s\nwant: %s", gotJSON, want)
	}
}

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

func TestWebFetchToolResultBlockToParam(t *testing.T) {
	const jsonData = `{"type":"web_fetch_tool_result","tool_use_id":"fetch123","caller":{"type":"direct"},"content":{"type":"web_fetch_result","url":"https://example.com","retrieved_at":"2026-05-30T00:00:00Z","content":{"type":"document","title":"Example","citations":{"enabled":true},"source":{"type":"text","data":"hello"}}}}`

	t.Run("stable", func(t *testing.T) {
		var block anthropic.ContentBlockUnion
		if err := json.Unmarshal([]byte(jsonData), &block); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}
		requireJSONEqual(t, jsonData, block.ToParam())
	})

	t.Run("beta", func(t *testing.T) {
		var block anthropic.BetaContentBlockUnion
		if err := json.Unmarshal([]byte(jsonData), &block); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}
		requireJSONEqual(t, jsonData, block.ToParam())
	})
}

func TestWebSearchToolResultBlockToParam(t *testing.T) {
	const jsonData = `{"type":"web_search_tool_result","tool_use_id":"search123","caller":{"type":"direct"},"content":[{"type":"web_search_result","title":"Example","url":"https://example.com","encrypted_content":"abc123"}]}`

	t.Run("stable", func(t *testing.T) {
		var block anthropic.ContentBlockUnion
		if err := json.Unmarshal([]byte(jsonData), &block); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}
		requireJSONEqual(t, jsonData, block.ToParam())
	})

	t.Run("beta", func(t *testing.T) {
		var block anthropic.BetaContentBlockUnion
		if err := json.Unmarshal([]byte(jsonData), &block); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}
		requireJSONEqual(t, jsonData, block.ToParam())
	})
}
