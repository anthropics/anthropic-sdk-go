package anthropic

import (
	"encoding/json"
	"fmt"
	"testing"
)

func sseEvent(t *testing.T, raw string) BetaManagedAgentsStreamSessionEventsUnion {
	t.Helper()
	var ev BetaManagedAgentsStreamSessionEventsUnion
	if err := json.Unmarshal([]byte(raw), &ev); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return ev
}

func eventStart(t *testing.T, eventID string) BetaManagedAgentsStreamSessionEventsUnion {
	return sseEvent(t, fmt.Sprintf(`{"type":"event_start","event":{"type":"agent.message","id":%q}}`, eventID))
}

func eventDelta(t *testing.T, eventID, text string, index int) BetaManagedAgentsStreamSessionEventsUnion {
	return sseEvent(t, fmt.Sprintf(`{"type":"event_delta","event_id":%q,"delta":{"type":"content_delta","index":%d,"content":{"type":"text","text":%q}}}`, eventID, index, text))
}

func feed(acc *BetaManagedAgentsEventAccumulator, evs ...BetaManagedAgentsStreamSessionEventsUnion) {
	for _, ev := range evs {
		acc.Accumulate(ev)
	}
}

func TestBetaManagedAgentsEventAccumulator_StartOpensEmptyPreview(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc, eventStart(t, "evt_1"))
	msg, ok := acc.AgentMessages["evt_1"]
	if !ok {
		t.Fatal("expected preview")
	}
	if msg.ID != "evt_1" || msg.Type != "agent.message" || len(msg.Content) != 0 {
		t.Fatalf("unexpected seed: %+v", msg)
	}
}

func TestBetaManagedAgentsEventAccumulator_StartIgnoresNonAgentMessage(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc, sseEvent(t, `{"type":"event_start","event":{"type":"agent.thinking","id":"evt_1"}}`))
	if _, ok := acc.AgentMessages["evt_1"]; ok {
		t.Fatal("expected no preview for agent.thinking start")
	}
}

func TestBetaManagedAgentsEventAccumulator_DeltaForIgnoredPreviewIsNoOp(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc,
		sseEvent(t, `{"type":"event_start","event":{"type":"agent.thinking","id":"evt_1"}}`),
		eventDelta(t, "evt_1", "x", 0),
	)
	if _, ok := acc.AgentMessages["evt_1"]; ok {
		t.Fatal("expected no preview for ignored type")
	}
}

func TestBetaManagedAgentsEventAccumulator_DeltaAppendsAndInserts(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc,
		eventStart(t, "evt_1"),
		eventDelta(t, "evt_1", "Hel", 0),
		eventDelta(t, "evt_1", "lo", 0),
		eventDelta(t, "evt_1", "World", 1),
	)
	msg := acc.AgentMessages["evt_1"]
	if len(msg.Content) != 2 || msg.Content[0].Text != "Hello" || msg.Content[1].Text != "World" {
		t.Fatalf("unexpected content: %+v", msg.Content)
	}
	if text := acc.AgentMessageText("evt_1"); text != "HelloWorld" {
		t.Fatalf("unexpected joined text: %q", text)
	}
}

func TestBetaManagedAgentsEventAccumulator_DeltaDefaultsIndexZero(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc,
		eventStart(t, "evt_1"),
		sseEvent(t, `{"type":"event_delta","event_id":"evt_1","delta":{"type":"content_delta","content":{"type":"text","text":"a"}}}`),
		sseEvent(t, `{"type":"event_delta","event_id":"evt_1","delta":{"type":"content_delta","content":{"type":"text","text":"b"}}}`),
	)
	if text := acc.AgentMessageText("evt_1"); text != "ab" {
		t.Fatalf("unexpected joined text: %q", text)
	}
}

func TestBetaManagedAgentsEventAccumulator_DeltaBeforeStartIsNoOp(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc, eventDelta(t, "evt_1", "x", 0))
	if _, ok := acc.AgentMessages["evt_1"]; ok {
		t.Fatal("expected no preview without event_start")
	}
}

func TestBetaManagedAgentsEventAccumulator_OutOfRangeIndexIsNoOp(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc,
		eventStart(t, "evt_1"),
		eventDelta(t, "evt_1", "x", 2),
		eventDelta(t, "evt_1", "y", -1),
	)
	if len(acc.AgentMessages["evt_1"].Content) != 0 {
		t.Fatalf("expected out-of-range deltas to be dropped, got %+v", acc.AgentMessages["evt_1"].Content)
	}
}

func TestBetaManagedAgentsEventAccumulator_BufferedEventReplacesPreview(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc,
		eventStart(t, "evt_1"),
		eventDelta(t, "evt_1", "partial", 0),
		sseEvent(t, `{"type":"agent.message","id":"evt_1","processed_at":"2024-01-01T00:00:00Z","content":[{"type":"text","text":"complete"}]}`),
	)
	if text := acc.AgentMessageText("evt_1"); text != "complete" {
		t.Fatalf("expected final event to replace preview, got %q", text)
	}
}

func TestBetaManagedAgentsEventAccumulator_StragglerDeltaAfterBufferedEventIsDropped(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc,
		eventStart(t, "evt_1"),
		eventDelta(t, "evt_1", "partial", 0),
		sseEvent(t, `{"type":"agent.message","id":"evt_1","processed_at":"2024-01-01T00:00:00Z","content":[{"type":"text","text":"complete"}]}`),
		eventDelta(t, "evt_1", "straggler", 0),
	)
	if text := acc.AgentMessageText("evt_1"); text != "complete" {
		t.Fatalf("expected straggler delta after the canonical event to be dropped, got %q", text)
	}
}

func TestBetaManagedAgentsEventAccumulator_ModelRequestEndClearsPreviews(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc,
		eventStart(t, "evt_1"),
		eventDelta(t, "evt_1", "x", 0),
		sseEvent(t, `{"type":"span.model_request_end","id":"sevt_2","model_request_start_id":"sevt_1","is_error":true,"processed_at":"2024-01-01T00:00:00Z"}`),
	)
	if len(acc.AgentMessages) != 0 {
		t.Fatal("expected previews to be cleared by span.model_request_end")
	}
}

func TestBetaManagedAgentsEventAccumulator_ModelRequestEndKeepsCanonicalMessages(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc,
		eventStart(t, "evt_1"),
		eventDelta(t, "evt_1", "partial", 0),
		sseEvent(t, `{"type":"agent.message","id":"evt_1","processed_at":"2024-01-01T00:00:00Z","content":[{"type":"text","text":"complete"}]}`),
		eventStart(t, "evt_2"),
		eventDelta(t, "evt_2", "open", 0),
		sseEvent(t, `{"type":"span.model_request_end","id":"sevt_2","model_request_start_id":"sevt_1","is_error":true,"processed_at":"2024-01-01T00:00:00Z"}`),
	)
	if text := acc.AgentMessageText("evt_1"); text != "complete" {
		t.Fatalf("expected canonical message to survive span.model_request_end, got %q", text)
	}
	if _, ok := acc.AgentMessages["evt_2"]; ok {
		t.Fatal("expected open preview to be discarded by span.model_request_end")
	}
}

func TestBetaManagedAgentsEventAccumulator_OtherEventsAreNoOps(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc, sseEvent(t, `{"type":"session.status_running","id":"sevt_1","processed_at":"2024-01-01T00:00:00Z"}`))
	if _, ok := acc.AgentMessages["sevt_1"]; ok {
		t.Fatal("expected no preview")
	}
}

func TestBetaManagedAgentsEventAccumulator_MultiplePreviews(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc,
		eventStart(t, "evt_a"),
		eventDelta(t, "evt_a", "alpha", 0),
		eventStart(t, "evt_b"),
		eventDelta(t, "evt_b", "beta", 0),
		sseEvent(t, `{"type":"event_start","event":{"type":"agent.thinking","id":"evt_c"}}`),
	)
	if len(acc.AgentMessages) != 2 {
		t.Fatalf("expected 2 previews, got %d", len(acc.AgentMessages))
	}
	if acc.AgentMessageText("evt_a") != "alpha" || acc.AgentMessageText("evt_b") != "beta" {
		t.Fatalf("unexpected: a=%q b=%q", acc.AgentMessageText("evt_a"), acc.AgentMessageText("evt_b"))
	}
}

func TestBetaManagedAgentsEventAccumulator_ZeroValue(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	if acc.AgentMessageText("evt_1") != "" {
		t.Fatal("expected empty text from zero value")
	}
	if len(acc.AgentMessages) != 0 {
		t.Fatal("expected empty map from zero value")
	}
}
