package anthropic_test

import (
	"encoding/json"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
)

func TestBetaManagedAgentsEventParamsOfUserMessageSetsType(t *testing.T) {
	union := anthropic.BetaManagedAgentsEventParamsOfUserMessage(
		[]anthropic.BetaManagedAgentsUserMessageEventParamsContentUnion{{
			OfText: &anthropic.BetaManagedAgentsTextBlockParam{
				Text: "hello",
				Type: anthropic.BetaManagedAgentsTextBlockTypeText,
			},
		}},
	)

	data, err := json.Marshal(union)
	if err != nil {
		t.Fatalf("failed to marshal: %s", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal: %s", err)
	}

	typ, ok := raw["type"]
	if !ok {
		t.Fatal("expected 'type' field to be present in marshaled JSON")
	}
	if typ != "user.message" {
		t.Fatalf("expected type 'user.message', got %q", typ)
	}
}

func TestBetaManagedAgentsEventParamsOfUserCustomToolResultSetsType(t *testing.T) {
	union := anthropic.BetaManagedAgentsEventParamsOfUserCustomToolResult("tool_use_123")

	data, err := json.Marshal(union)
	if err != nil {
		t.Fatalf("failed to marshal: %s", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal: %s", err)
	}

	typ, ok := raw["type"]
	if !ok {
		t.Fatal("expected 'type' field to be present in marshaled JSON")
	}
	if typ != "user.custom_tool_result" {
		t.Fatalf("expected type 'user.custom_tool_result', got %q", typ)
	}
}
