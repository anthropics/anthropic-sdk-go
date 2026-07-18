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

	// Regression test for https://github.com/anthropics/anthropic-sdk-go/issues/317.
	// The API requires `error_code` on a tool_search_tool_result_error block, but
	// the param field is tagged `omitzero`, so a zero-value enum (e.g. an enum value
	// the SDK does not yet recognize) would be elided and the API would 400 on the
	// next request.
	t.Run("ToolSearchToolResultBlock error preserves known error_code", func(t *testing.T) {
		result := unmarshalContentBlockParam(t, `{"type":"tool_search_tool_result","tool_use_id":"tsu_1","content":{"type":"tool_search_tool_result_error","error_code":"unavailable"}}`)
		if result.OfToolSearchToolResult == nil {
			t.Fatal("Expected OfToolSearchToolResult to be non-nil")
		}
		if result.OfToolSearchToolResult.Content.OfRequestToolSearchToolResultError == nil {
			t.Fatal("Expected OfRequestToolSearchToolResultError to be non-nil")
		}

		data, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}
		var decoded struct {
			Content struct {
				ErrorCode *string `json:"error_code"`
				Type      string  `json:"type"`
			} `json:"content"`
		}
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal round-trip: %v", err)
		}
		if decoded.Content.Type != "tool_search_tool_result_error" {
			t.Errorf("Expected content.type 'tool_search_tool_result_error', got %q", decoded.Content.Type)
		}
		if decoded.Content.ErrorCode == nil {
			t.Fatalf("Expected content.error_code to be present, got missing field; raw: %s", string(data))
		}
		if *decoded.Content.ErrorCode != "unavailable" {
			t.Errorf("Expected content.error_code 'unavailable', got %q", *decoded.Content.ErrorCode)
		}
	})

	t.Run("ToolSearchToolResultBlock error preserves empty error_code", func(t *testing.T) {
		// Regression: when the response's error_code resolves to the zero value
		// of ToolSearchToolResultErrorCode (an empty string), the param field
		// is tagged `omitzero` and the field gets silently dropped. The API
		// then rejects the next request with
		// "tool_search_tool_result.content.RequestToolSearchToolResultError.error_code: Field required".
		// ToParam() must still emit the field with the original value.
		result := unmarshalContentBlockParam(t, `{"type":"tool_search_tool_result","tool_use_id":"tsu_1","content":{"type":"tool_search_tool_result_error","error_code":""}}`)
		if result.OfToolSearchToolResult == nil {
			t.Fatal("Expected OfToolSearchToolResult to be non-nil")
		}
		if result.OfToolSearchToolResult.Content.OfRequestToolSearchToolResultError == nil {
			t.Fatal("Expected OfRequestToolSearchToolResultError to be non-nil")
		}

		data, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}
		var decoded struct {
			Content struct {
				ErrorCode *string `json:"error_code"`
				Type      string  `json:"type"`
			} `json:"content"`
		}
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Failed to unmarshal round-trip: %v", err)
		}
		if decoded.Content.Type != "tool_search_tool_result_error" {
			t.Errorf("Expected content.type 'tool_search_tool_result_error', got %q", decoded.Content.Type)
		}
		if decoded.Content.ErrorCode == nil {
			t.Fatalf("Expected content.error_code to be present even when empty, got missing field; raw: %s", string(data))
		}
		if *decoded.Content.ErrorCode != "" {
			t.Errorf("Expected content.error_code '', got %q", *decoded.Content.ErrorCode)
		}
	})
}
