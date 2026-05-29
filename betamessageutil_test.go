package anthropic_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
)

// TestBetaAccumulateRecoversFromInvalidToolUseInput mirrors
// TestAccumulateRecoversFromInvalidToolUseInput for the Beta API surface —
// see that test's comment for the bug we're guarding against
// (https://github.com/anthropics/anthropic-sdk-go/issues/292).
func TestBetaAccumulateRecoversFromInvalidToolUseInput(t *testing.T) {
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

			message := anthropic.BetaMessage{}
			for _, eventStr := range events {
				event := anthropic.BetaRawMessageStreamEventUnion{}
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

func TestBetaAccumulateLeavesValidInputUntouched(t *testing.T) {
	events := []string{
		`{"type": "message_start", "message": {"id": "msg_x", "type": "message", "role": "assistant", "content": [], "model": "test", "usage": {"input_tokens": 0, "output_tokens": 0}}}`,
		`{"type": "content_block_start", "index": 0, "content_block": {"type": "tool_use", "id": "toolu_id", "name": "tool_name", "input": {}}}`,
		`{"type": "content_block_delta", "index": 0, "delta": {"type": "input_json_delta", "partial_json": "{\"argument\":"}}`,
		`{"type": "content_block_delta", "index": 0, "delta": {"type": "input_json_delta", "partial_json": " \"value\"}"}}`,
		`{"type": "content_block_stop", "index": 0}`,
		`{"type": "message_stop"}`,
	}
	message := anthropic.BetaMessage{}
	for _, eventStr := range events {
		event := anthropic.BetaRawMessageStreamEventUnion{}
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
