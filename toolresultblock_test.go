package anthropic_test

import (
	"encoding/json"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/tidwall/gjson"
)

func TestToolResultBlockParam_MarshalJSON_StringContent(t *testing.T) {
	// When ContentString is set and Content is empty, content should serialize as a string.
	block := anthropic.ToolResultBlockParam{
		ToolUseID:     "call_123",
		ContentString: "15 degrees",
		IsError:       param.Opt[bool]{Value: false},
	}

	data, err := json.Marshal(block)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	result := gjson.ParseBytes(data)

	if result.Get("tool_use_id").String() != "call_123" {
		t.Errorf("Expected tool_use_id 'call_123', got '%s'", result.Get("tool_use_id").String())
	}
	if result.Get("type").String() != "tool_result" {
		t.Errorf("Expected type 'tool_result', got '%s'", result.Get("type").String())
	}
	// Content should be a string, not an array
	contentResult := result.Get("content")
	if contentResult.Type != gjson.String {
		t.Fatalf("Expected content to be a string, got type %v: %s", contentResult.Type, contentResult.Raw)
	}
	if contentResult.String() != "15 degrees" {
		t.Errorf("Expected content '15 degrees', got '%s'", contentResult.String())
	}
}

func TestToolResultBlockParam_MarshalJSON_ArrayContent(t *testing.T) {
	// When Content array is set, content should serialize as an array.
	block := anthropic.ToolResultBlockParam{
		ToolUseID: "call_456",
		Content: []anthropic.ToolResultBlockParamContentUnion{
			{OfText: &anthropic.TextBlockParam{Text: "some result"}},
		},
		IsError: param.Opt[bool]{Value: false},
	}

	data, err := json.Marshal(block)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	result := gjson.ParseBytes(data)
	contentResult := result.Get("content")
	if !contentResult.IsArray() {
		t.Fatalf("Expected content to be an array, got: %s", contentResult.Raw)
	}
	items := contentResult.Array()
	if len(items) != 1 {
		t.Fatalf("Expected 1 content item, got %d", len(items))
	}
	if items[0].Get("text").String() != "some result" {
		t.Errorf("Expected text 'some result', got '%s'", items[0].Get("text").String())
	}
}

func TestToolResultBlockParam_UnmarshalJSON_StringContent(t *testing.T) {
	// The API docs say content can be a string. When deserializing, a string content
	// should be normalized into a Content array with a single TextBlock.
	jsonStr := `{
		"tool_use_id": "call_789",
		"type": "tool_result",
		"content": "15 degrees",
		"is_error": false
	}`

	var block anthropic.ToolResultBlockParam
	err := json.Unmarshal([]byte(jsonStr), &block)
	if err != nil {
		t.Fatalf("Failed to unmarshal string content: %v", err)
	}

	if block.ToolUseID != "call_789" {
		t.Errorf("Expected tool_use_id 'call_789', got '%s'", block.ToolUseID)
	}
	if len(block.Content) != 1 {
		t.Fatalf("Expected 1 content item after string normalization, got %d", len(block.Content))
	}
	if block.Content[0].OfText == nil {
		t.Fatal("Expected OfText to be non-nil")
	}
	if block.Content[0].OfText.Text != "15 degrees" {
		t.Errorf("Expected text '15 degrees', got '%s'", block.Content[0].OfText.Text)
	}
}

func TestToolResultBlockParam_UnmarshalJSON_ArrayContent(t *testing.T) {
	// Standard array content format should still work.
	jsonStr := `{
		"tool_use_id": "call_abc",
		"type": "tool_result",
		"content": [{"type": "text", "text": "hello world"}],
		"is_error": false
	}`

	var block anthropic.ToolResultBlockParam
	err := json.Unmarshal([]byte(jsonStr), &block)
	if err != nil {
		t.Fatalf("Failed to unmarshal array content: %v", err)
	}

	if block.ToolUseID != "call_abc" {
		t.Errorf("Expected tool_use_id 'call_abc', got '%s'", block.ToolUseID)
	}
	if len(block.Content) != 1 {
		t.Fatalf("Expected 1 content item, got %d", len(block.Content))
	}
	if block.Content[0].OfText == nil {
		t.Fatal("Expected OfText to be non-nil")
	}
	if block.Content[0].OfText.Text != "hello world" {
		t.Errorf("Expected text 'hello world', got '%s'", block.Content[0].OfText.Text)
	}
}

func TestNewToolResultBlock_SerializesAsString(t *testing.T) {
	// NewToolResultBlock should produce content as a string in JSON.
	block := anthropic.NewToolResultBlock("call_123", "sunny 72F", false)

	data, err := json.Marshal(block)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	result := gjson.ParseBytes(data)
	contentResult := result.Get("content")
	if contentResult.Type != gjson.String {
		t.Fatalf("Expected content to be a string, got type %v: %s", contentResult.Type, contentResult.Raw)
	}
	if contentResult.String() != "sunny 72F" {
		t.Errorf("Expected content 'sunny 72F', got '%s'", contentResult.String())
	}
}

func TestNewToolResultBlockFromArray_SerializesAsArray(t *testing.T) {
	// NewToolResultBlockFromArray should produce content as an array in JSON.
	content := []anthropic.ToolResultBlockParamContentUnion{
		{OfText: &anthropic.TextBlockParam{Text: "result text"}},
	}
	block := anthropic.NewToolResultBlockFromArray("call_456", content, false)

	data, err := json.Marshal(block)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	result := gjson.ParseBytes(data)
	contentResult := result.Get("content")
	if !contentResult.IsArray() {
		t.Fatalf("Expected content to be an array, got: %s", contentResult.Raw)
	}
}

func TestToolResultBlockParam_RoundTrip_StringContent(t *testing.T) {
	// Marshal with string content, then unmarshal should normalize to Content array.
	original := anthropic.ToolResultBlockParam{
		ToolUseID:     "call_rt",
		ContentString: "round trip test",
		IsError:       param.Opt[bool]{Value: true},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Verify JSON has string content
	result := gjson.ParseBytes(data)
	if result.Get("content").Type != gjson.String {
		t.Fatalf("Expected marshaled content to be a string, got: %s", result.Get("content").Raw)
	}

	var restored anthropic.ToolResultBlockParam
	err = json.Unmarshal(data, &restored)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if restored.ToolUseID != original.ToolUseID {
		t.Errorf("ToolUseID mismatch: got '%s', want '%s'", restored.ToolUseID, original.ToolUseID)
	}
	// After unmarshaling string content, it should be in Content array as a TextBlock.
	if len(restored.Content) != 1 || restored.Content[0].OfText == nil {
		t.Fatal("Expected Content to have 1 TextBlock after round-trip")
	}
	if restored.Content[0].OfText.Text != "round trip test" {
		t.Errorf("Content text mismatch: got '%s', want 'round trip test'", restored.Content[0].OfText.Text)
	}
}

func TestBetaToolResultBlockParam_MarshalJSON_StringContent(t *testing.T) {
	block := anthropic.BetaToolResultBlockParam{
		ToolUseID:     "call_beta_123",
		ContentString: "beta string content",
		IsError:       param.Opt[bool]{Value: false},
	}

	data, err := json.Marshal(block)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	result := gjson.ParseBytes(data)
	contentResult := result.Get("content")
	if contentResult.Type != gjson.String {
		t.Fatalf("Expected content to be a string, got type %v: %s", contentResult.Type, contentResult.Raw)
	}
	if contentResult.String() != "beta string content" {
		t.Errorf("Expected content 'beta string content', got '%s'", contentResult.String())
	}
}

func TestBetaToolResultBlockParam_UnmarshalJSON_StringContent(t *testing.T) {
	jsonStr := `{
		"tool_use_id": "call_beta_789",
		"type": "tool_result",
		"content": "beta 15 degrees",
		"is_error": false
	}`

	var block anthropic.BetaToolResultBlockParam
	err := json.Unmarshal([]byte(jsonStr), &block)
	if err != nil {
		t.Fatalf("Failed to unmarshal string content: %v", err)
	}

	if len(block.Content) != 1 || block.Content[0].OfText == nil {
		t.Fatal("Expected Content to have 1 TextBlock after string normalization")
	}
	if block.Content[0].OfText.Text != "beta 15 degrees" {
		t.Errorf("Expected text 'beta 15 degrees', got '%s'", block.Content[0].OfText.Text)
	}
}
