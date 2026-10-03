package anthropic

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"

	"fmt"
	"github.com/anthropics/anthropic-sdk-go/option"
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

func TestBetaManagedAgentsEventAccumulator_DeltaSeedsTheUnionVariant(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc,
		eventStart(t, "evt_1"),
		eventDelta(t, "evt_1", "Hel", 0),
	)
	block := acc.AgentMessages["evt_1"].Content[0]
	if _, ok := block.AsAny().(BetaManagedAgentsTextBlock); !ok {
		t.Fatalf("expected preview block to report the text variant, got %T", block.AsAny())
	}
}

func TestBetaManagedAgentsEventAccumulator_ContentlessDeltaKeepsIndexAlignment(t *testing.T) {
	// Deltas address entries by index, so a fragment carrying no content still
	// has to occupy its slot or every later fragment falls out of range.
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc,
		eventStart(t, "evt_1"),
		sseEvent(t, `{"type":"event_delta","event_id":"evt_1","delta":{"type":"content_delta","index":0}}`),
		eventDelta(t, "evt_1", "World", 1),
	)
	msg := acc.AgentMessages["evt_1"]
	if len(msg.Content) != 2 || msg.Content[1].Text != "World" {
		t.Fatalf("unexpected content: %+v", msg.Content)
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

func TestBetaManagedAgentsEventAccumulator_UnrecognizedEventTypeIsNoOp(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc,
		eventStart(t, "evt_1"),
		eventDelta(t, "evt_1", "Hi", 0),
		sseEvent(t, `{"type":"session.some_future_event","id":"evt_1","result":{"type":"object"}}`),
	)
	if text := acc.AgentMessageText("evt_1"); text != "Hi" {
		t.Fatalf("expected preview to be untouched by an unrecognized event type, got %q", text)
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

func TestBetaManagedAgentsEventAccumulator_ReplayedStartKeepsCanonicalMessage(t *testing.T) {
	for _, withPreview := range []bool{false, true} {
		t.Run(fmt.Sprintf("initial_preview=%t", withPreview), func(t *testing.T) {
			var acc BetaManagedAgentsEventAccumulator
			if withPreview {
				feed(&acc, eventStart(t, "evt_final"), eventDelta(t, "evt_final", "unfinished", 0))
			}
			raw := `{"type":"agent.message","id":"evt_final","processed_at":"2026-01-01T00:00:00Z","content":[{"type":"text","text":"complete"}],"future_field":{"n":9007199254740993}}`
			final := sseEvent(t, raw)
			feed(&acc, final, eventStart(t, "evt_final"))
			got := acc.AgentMessages["evt_final"]
			if got.RawJSON() != raw || !got.JSON.ProcessedAt.Valid() || acc.AgentMessageText("evt_final") != "complete" {
				t.Fatalf("replayed start overwrote final state: %+v", got)
			}
			feed(&acc, eventDelta(t, "evt_final", "stale", 0), eventDelta(t, "evt_final", "extra", 1),
				sseEvent(t, `{"type":"span.model_request_end","id":"end","processed_at":"2026-01-01T00:00:01Z"}`))
			if got := acc.AgentMessages["evt_final"]; got.RawJSON() != raw || acc.AgentMessageText("evt_final") != "complete" {
				t.Fatalf("final message was changed or removed after restart: %+v", got)
			}
			if final.RawJSON() != raw {
				t.Fatal("input event was changed")
			}
		})
	}
}

func TestBetaManagedAgentsEventAccumulator_RestartControlsAndFinalReplacement(t *testing.T) {
	var acc BetaManagedAgentsEventAccumulator
	feed(&acc, eventStart(t, "open"), eventDelta(t, "open", "old preview", 0), eventStart(t, "open"), eventDelta(t, "open", "fresh preview", 0))
	if acc.AgentMessageText("open") != "fresh preview" {
		t.Fatal("an unfinished preview can no longer restart")
	}
	feed(&acc, sseEvent(t, `{"type":"agent.message","id":"final","processed_at":"2026-01-01T00:00:00Z","content":[{"type":"text","text":"first"}]}`),
		sseEvent(t, `{"type":"agent.message","id":"final","processed_at":"2026-01-01T00:00:01Z","content":[{"type":"text","text":"replacement"}]}`))
	if acc.AgentMessageText("final") != "replacement" {
		t.Fatal("new canonical event did not replace its predecessor")
	}
	feed(&acc, eventStart(t, "independent"), eventDelta(t, "independent", "another preview", 0))
	if acc.AgentMessageText("independent") != "another preview" {
		t.Fatal("other event IDs stopped accumulating")
	}
	var nilAccumulator *BetaManagedAgentsEventAccumulator
	nilAccumulator.Accumulate(eventStart(t, "ignored"))
}

func TestBetaManagedAgentsEventAccumulator_HTTPReplayRetainsFinalTranscript(t *testing.T) {
	rawFinal := `{"type":"agent.message","id":"final","processed_at":"2026-01-01T00:00:00Z","content":[{"type":"text","text":"canonical output"}]}`
	events := []BetaManagedAgentsStreamSessionEventsUnion{
		eventStart(t, "final"), eventDelta(t, "final", "partial", 0), sseEvent(t, rawFinal),
		eventStart(t, "final"), eventDelta(t, "final", "stale", 0),
		sseEvent(t, `{"type":"span.model_request_end","id":"end","processed_at":"2026-01-01T00:00:01Z"}`),
	}
	var wire strings.Builder
	for _, event := range events {
		fmt.Fprintf(&wire, "event: %s\ndata: %s\n\n", event.Type, event.RawJSON())
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/sessions/session_test/events/stream" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		text := wire.String()
		for len(text) > 0 {
			n := min(7, len(text))
			if _, err := fmt.Fprint(w, text[:n]); err != nil {
				return
			}
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
			text = text[n:]
		}
	}))
	defer server.Close()
	client := NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test-key"), option.WithMaxRetries(0))
	stream := client.Beta.Sessions.Events.StreamEvents(context.Background(), "session_test", BetaSessionEventStreamParams{})
	defer stream.Close()
	var acc BetaManagedAgentsEventAccumulator
	count := 0
	for stream.Next() {
		acc.Accumulate(stream.Current())
		count++
	}
	if err := stream.Err(); err != nil {
		t.Fatal(err)
	}
	if count != len(events) {
		t.Fatalf("decoded %d events; want %d", count, len(events))
	}
	if len(acc.AgentMessages) != 1 || acc.AgentMessageText("final") != "canonical output" || acc.AgentMessages["final"].RawJSON() != rawFinal {
		t.Fatalf("final transcript was lost after decoded replay: %+v", acc.AgentMessages)
	}
}
