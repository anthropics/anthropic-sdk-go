package anthropic_test

import (
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

func marshalContentBlockParamObject(t *testing.T, jsonData string) map[string]any {
	t.Helper()

	paramBlock := unmarshalContentBlockParam(t, jsonData)
	data, err := json.Marshal(paramBlock)
	if err != nil {
		t.Fatalf("Failed to marshal param JSON: %v", err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to decode marshaled param JSON: %v\nJSON: %s", err, data)
	}
	return result
}

func marshalBetaContentBlockParamObject(t *testing.T, jsonData string) map[string]any {
	t.Helper()

	var block anthropic.BetaContentBlockUnion
	if err := json.Unmarshal([]byte(jsonData), &block); err != nil {
		t.Fatalf("Failed to unmarshal beta JSON: %v", err)
	}
	data, err := json.Marshal(block.ToParam())
	if err != nil {
		t.Fatalf("Failed to marshal beta param JSON: %v", err)
	}
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to decode marshaled beta param JSON: %v\nJSON: %s", err, data)
	}
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

	t.Run("WebSearchToolResultBlock marshals search results under list and preserves caller", func(t *testing.T) {
		result := marshalContentBlockParamObject(t, `{
			"type":"web_search_tool_result",
			"tool_use_id":"toolu_search_123",
			"caller":{"type":"code_execution_20250825","tool_id":"srvtool_123"},
			"content":[{"type":"web_search_result","title":"Example Article","url":"https://example.com/article","encrypted_content":"abc123","page_age":"1 day ago"}]
		}`)

		content, ok := result["content"].(map[string]any)
		if !ok {
			t.Fatalf("content = %#v, want object with list", result["content"])
		}
		list, ok := content["list"].([]any)
		if !ok || len(list) != 1 {
			t.Fatalf("content.list = %#v, want one result", content["list"])
		}
		caller, ok := result["caller"].(map[string]any)
		if !ok {
			t.Fatalf("caller = %#v, want object", result["caller"])
		}
		if caller["type"] != "code_execution_20250825" || caller["tool_id"] != "srvtool_123" {
			t.Fatalf("caller = %#v, want code execution caller", caller)
		}
	})

	t.Run("WebFetchToolResultBlock marshals content and preserves caller", func(t *testing.T) {
		result := marshalContentBlockParamObject(t, `{
			"type":"web_fetch_tool_result",
			"tool_use_id":"toolu_fetch_456",
			"caller":{"type":"code_execution_20250825","tool_id":"srvtool_456"},
			"content":{
				"type":"web_fetch_result",
				"url":"https://example.com/article",
				"retrieved_at":"2026-01-01T00:00:00Z",
				"content":{
					"type":"document",
					"title":"Example Article",
					"citations":{"enabled":false},
					"source":{"type":"text","media_type":"text/plain","data":"Fetched body"}
				}
			}
		}`)

		content, ok := result["content"].(map[string]any)
		if !ok {
			t.Fatalf("content = %#v, want object", result["content"])
		}
		if content["type"] != "web_fetch_result" || content["url"] != "https://example.com/article" {
			t.Fatalf("content = %#v, want web fetch result", content)
		}
		caller, ok := result["caller"].(map[string]any)
		if !ok {
			t.Fatalf("caller = %#v, want object", result["caller"])
		}
		if caller["type"] != "code_execution_20250825" || caller["tool_id"] != "srvtool_456" {
			t.Fatalf("caller = %#v, want code execution caller", caller)
		}
	})
}

func TestBetaContentBlockUnionToParamWebToolResults(t *testing.T) {
	t.Run("BetaWebSearchToolResultBlock marshals search results under list and preserves caller", func(t *testing.T) {
		result := marshalBetaContentBlockParamObject(t, `{
			"type":"web_search_tool_result",
			"tool_use_id":"toolu_search_beta",
			"caller":{"type":"code_execution_20250825","tool_id":"srvtool_beta_search"},
			"content":[{"type":"web_search_result","title":"Example Article","url":"https://example.com/article","encrypted_content":"abc123","page_age":"1 day ago"}]
		}`)

		content, ok := result["content"].(map[string]any)
		if !ok {
			t.Fatalf("content = %#v, want object with list", result["content"])
		}
		list, ok := content["list"].([]any)
		if !ok || len(list) != 1 {
			t.Fatalf("content.list = %#v, want one result", content["list"])
		}
		caller, ok := result["caller"].(map[string]any)
		if !ok {
			t.Fatalf("caller = %#v, want object", result["caller"])
		}
		if caller["type"] != "code_execution_20250825" || caller["tool_id"] != "srvtool_beta_search" {
			t.Fatalf("caller = %#v, want code execution caller", caller)
		}
	})

	t.Run("BetaWebFetchToolResultBlock marshals content and preserves caller", func(t *testing.T) {
		result := marshalBetaContentBlockParamObject(t, `{
			"type":"web_fetch_tool_result",
			"tool_use_id":"toolu_fetch_beta",
			"caller":{"type":"code_execution_20250825","tool_id":"srvtool_beta_fetch"},
			"content":{
				"type":"web_fetch_result",
				"url":"https://example.com/article",
				"retrieved_at":"2026-01-01T00:00:00Z",
				"content":{
					"type":"document",
					"title":"Example Article",
					"citations":{"enabled":false},
					"source":{"type":"text","media_type":"text/plain","data":"Fetched body"}
				}
			}
		}`)

		content, ok := result["content"].(map[string]any)
		if !ok {
			t.Fatalf("content = %#v, want object", result["content"])
		}
		if content["type"] != "web_fetch_result" || content["url"] != "https://example.com/article" {
			t.Fatalf("content = %#v, want web fetch result", content)
		}
		caller, ok := result["caller"].(map[string]any)
		if !ok {
			t.Fatalf("caller = %#v, want object", result["caller"])
		}
		if caller["type"] != "code_execution_20250825" || caller["tool_id"] != "srvtool_beta_fetch" {
			t.Fatalf("caller = %#v, want code execution caller", caller)
		}
	})
}
