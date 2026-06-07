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
}

func TestContentBlockUnionToParamTextEditorCodeExecutionResult(t *testing.T) {
	const wire = `{
		"type": "text_editor_code_execution_tool_result",
		"tool_use_id": "srvtoolu_1",
		"content": {
			"type": "text_editor_code_execution_view_result",
			"content": "line1\nline2\n",
			"file_type": "text",
			"num_lines": 2,
			"start_line": 1,
			"total_lines": 2
		}
	}`

	var block anthropic.ContentBlockUnion
	if err := json.Unmarshal([]byte(wire), &block); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	assertTextEditorCodeExecutionResultContentObject(t, block.ToParam())
}

func TestBetaContentBlockUnionToParamTextEditorCodeExecutionResult(t *testing.T) {
	const wire = `{
		"type": "text_editor_code_execution_tool_result",
		"tool_use_id": "srvtoolu_1",
		"content": {
			"type": "text_editor_code_execution_view_result",
			"content": "line1\nline2\n",
			"file_type": "text",
			"num_lines": 2,
			"start_line": 1,
			"total_lines": 2
		}
	}`

	var block anthropic.BetaContentBlockUnion
	if err := json.Unmarshal([]byte(wire), &block); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	assertTextEditorCodeExecutionResultContentObject(t, block.ToParam())
}

func assertTextEditorCodeExecutionResultContentObject(t *testing.T, paramBlock any) {
	t.Helper()

	out, err := json.Marshal(paramBlock)
	if err != nil {
		t.Fatalf("Failed to marshal param block: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("Failed to unmarshal param JSON: %v", err)
	}

	content, ok := got["content"].(map[string]any)
	if !ok {
		t.Fatalf("Expected content to marshal as an object, got %T: %v", got["content"], got["content"])
	}
	if got := content["type"]; got != "text_editor_code_execution_view_result" {
		t.Errorf("Expected content.type to round-trip, got %v", got)
	}
	if got := content["content"]; got != "line1\nline2\n" {
		t.Errorf("Expected content.content to round-trip, got %v", got)
	}
}
