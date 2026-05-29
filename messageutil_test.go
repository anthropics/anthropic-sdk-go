package anthropic_test

import (
	"encoding/json"
	"fmt"
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

// TestAccumulateRecoversFromInvalidToolUseInput pins the fix for
// https://github.com/anthropics/anthropic-sdk-go/issues/292. Before the fix,
// Accumulate's json.Marshal at content_block_stop / message_stop would fail
// with "unexpected end of JSON input" whenever an input_json_delta sequence
// left cb.Input as either empty-non-nil, truncated, or with an unclosed
// string — typical when max_tokens cuts off mid tool_use. The fix replaces
// such Input with []byte("{}") so the stream survives and the caller's tool
// dispatcher gets a structurally valid tool_use to handle.
func TestAccumulateRecoversFromInvalidToolUseInput(t *testing.T) {
	for name, partials := range map[string][]string{
		"empty non-nil":         {""},
		"truncated":             {`{"argument":`},
		"unclosed string":       {`{"x": "abc`},
		"multi-delta truncated": {`{"args`, `ument":`},
	} {
		t.Run(name, func(t *testing.T) {
			events := []string{
				`{"type": "message_start", "message": {"id": "msg_x", "type": "message", "role": "assistant", "content": [], "model": "test", "usage": {"input_tokens": 0, "output_tokens": 0}}}`,
				`{"type": "content_block_start", "index": 0, "content_block": {"type": "tool_use", "id": "toolu_id", "name": "tool_name", "input": {}}}`,
			}
			for _, partial := range partials {
				events = append(events, fmt.Sprintf(
					`{"type": "content_block_delta", "index": 0, "delta": {"type": "input_json_delta", "partial_json": %q}}`,
					partial,
				))
			}
			events = append(events,
				`{"type": "content_block_stop", "index": 0}`,
				`{"type": "message_stop"}`,
			)

			message := anthropic.Message{}
			for _, eventStr := range events {
				event := anthropic.MessageStreamEventUnion{}
				if err := (&event).UnmarshalJSON([]byte(eventStr)); err != nil {
					t.Fatalf("unmarshal %s: %v", eventStr, err)
				}
				if err := (&message).Accumulate(event); err != nil {
					t.Fatalf("Accumulate must not error on the malformed-input case; got %v", err)
				}
			}
			if _, err := json.Marshal(message); err != nil {
				t.Fatalf("json.Marshal must succeed after sanitisation; got %v", err)
			}
			if len(message.Content) != 1 {
				t.Fatalf("expected one content block; got %d", len(message.Content))
			}
			if got := string(message.Content[0].Input); got != "{}" {
				t.Errorf("Input should be sanitised to {}; got %q", got)
			}
		})
	}
}

// TestAccumulateLeavesValidInputUntouched is the negative case for the same
// fix: a tool_use whose accumulated Input is valid JSON must not be rewritten.
func TestAccumulateLeavesValidInputUntouched(t *testing.T) {
	events := []string{
		`{"type": "message_start", "message": {"id": "msg_x", "type": "message", "role": "assistant", "content": [], "model": "test", "usage": {"input_tokens": 0, "output_tokens": 0}}}`,
		`{"type": "content_block_start", "index": 0, "content_block": {"type": "tool_use", "id": "toolu_id", "name": "tool_name", "input": {}}}`,
		`{"type": "content_block_delta", "index": 0, "delta": {"type": "input_json_delta", "partial_json": "{\"argument\":"}}`,
		`{"type": "content_block_delta", "index": 0, "delta": {"type": "input_json_delta", "partial_json": " \"value\"}"}}`,
		`{"type": "content_block_stop", "index": 0}`,
		`{"type": "message_stop"}`,
	}
	message := anthropic.Message{}
	for _, eventStr := range events {
		event := anthropic.MessageStreamEventUnion{}
		if err := (&event).UnmarshalJSON([]byte(eventStr)); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if err := (&message).Accumulate(event); err != nil {
			t.Fatal(err)
		}
	}
	if got := string(message.Content[0].Input); got != `{"argument": "value"}` {
		t.Errorf("valid Input must not be touched; got %q", got)
	}
}
