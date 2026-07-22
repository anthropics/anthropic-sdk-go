package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

var sessionRunnerSilentLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

// stubBetaTool is a BetaTool with controllable Execute behavior, used to
// exercise both happy and failure paths in the runner.
type stubBetaTool struct {
	name string
	run  func(ctx context.Context, input json.RawMessage) (string, bool)
	runs atomic.Int32
}

func (s *stubBetaTool) Name() string        { return s.name }
func (s *stubBetaTool) Description() string { return s.name }
func (s *stubBetaTool) InputSchema() BetaToolInputSchemaParam {
	return BetaToolInputSchemaParam{Properties: map[string]any{}}
}
func (s *stubBetaTool) Execute(ctx context.Context, input json.RawMessage) ([]BetaToolResultBlockParamContentUnion, error) {
	s.runs.Add(1)
	content, isErr := "ok from "+s.name, false
	if s.run != nil {
		content, isErr = s.run(ctx, input)
	}
	if isErr {
		return nil, fmt.Errorf("%s", content)
	}
	return []BetaToolResultBlockParamContentUnion{{OfText: &BetaTextBlockParam{Text: content}}}, nil
}

// dispatchedResultText joins the text blocks of the result event the runner
// built for call — CustomResult when call.Custom, otherwise Result. Replaces
// the removed flat DispatchedToolCall.Content convenience field.
func dispatchedResultText(call DispatchedToolCall) string {
	var out string
	add := func(s string) {
		if out != "" {
			out += "\n"
		}
		out += s
	}
	if call.Custom {
		for _, b := range call.CustomResult.Content {
			if b.OfText != nil {
				add(b.OfText.Text)
			}
		}
		return out
	}
	for _, b := range call.Result.Content {
		if b.OfText != nil {
			add(b.OfText.Text)
		}
	}
	return out
}

func sseLine(eventType string, payload any) string {
	body, _ := json.Marshal(payload)
	return fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, body)
}

func toolUseEvt(id, name string, input map[string]any) map[string]any {
	return map[string]any{
		"type":         "agent.tool_use",
		"id":           id,
		"name":         name,
		"input":        input,
		"processed_at": "2026-05-11T12:00:00Z",
	}
}

func customToolUseEvt(id, name string, input map[string]any) map[string]any {
	return map[string]any{
		"type":         "agent.custom_tool_use",
		"id":           id,
		"name":         name,
		"input":        input,
		"processed_at": "2026-05-11T12:00:00Z",
	}
}

func plainEvt(id, eventType string) map[string]any {
	return map[string]any{"type": eventType, "id": id, "processed_at": "2026-05-11T12:00:00Z"}
}

// askToolUseEvt is a toolUseEvt with evaluated_permission set (e.g. "ask"
// for an always_ask tool, "deny" for a server-blocked one).
func askToolUseEvt(id, name string, input map[string]any, permission string) map[string]any {
	ev := toolUseEvt(id, name, input)
	ev["evaluated_permission"] = permission
	return ev
}

// toolConfirmationEvt is a user.tool_confirmation resolving a held tool_use:
// "allow" releases it to run, "deny" drops it.
func toolConfirmationEvt(toolUseID, result string) map[string]any {
	return map[string]any{
		"type":         "user.tool_confirmation",
		"id":           "evt_conf_" + toolUseID,
		"tool_use_id":  toolUseID,
		"result":       result,
		"processed_at": "2026-05-11T12:00:00Z",
	}
}

func idleEndTurnEvt(id string) map[string]any {
	return map[string]any{
		"type":         "session.status_idle",
		"id":           id,
		"stop_reason":  map[string]any{"type": "end_turn"},
		"processed_at": "2026-05-11T12:00:00Z",
	}
}

// idleRequiresActionEvt is a session.status_idle with stop_reason
// "requires_action" — the server parks here waiting on the listed events
// (tool confirmations, tool results).
func idleRequiresActionEvt(id string, eventIDs ...string) map[string]any {
	return map[string]any{
		"type":         "session.status_idle",
		"id":           id,
		"stop_reason":  map[string]any{"type": "requires_action", "event_ids": eventIDs},
		"processed_at": "2026-05-11T12:00:00Z",
	}
}

func interruptEvt(id string, processed bool, sessionThreadID string) map[string]any {
	event := map[string]any{
		"type":         "user.interrupt",
		"id":           id,
		"processed_at": nil,
	}
	if processed {
		event["processed_at"] = "2026-05-11T12:00:01Z"
	}
	if sessionThreadID != "" {
		event["session_thread_id"] = sessionThreadID
	}
	return event
}

func emptyEventList() string {
	body, _ := json.Marshal(map[string]any{"data": []any{}, "first_id": nil, "has_more": false, "last_id": nil})
	return string(body)
}

func sendOK() string {
	body, _ := json.Marshal(map[string]any{"type": "send_session_events"})
	return string(body)
}

// sessionEventsServer is a minimal httptest server scripting the
// /v1/sessions/{id}/events{,/stream} endpoints the SessionToolRunner uses.
type sessionEventsServer struct {
	t      *testing.T
	server *httptest.Server

	HandleStream http.HandlerFunc // GET  /v1/sessions/{id}/events/stream
	HandleSend   http.HandlerFunc // POST /v1/sessions/{id}/events
	HandleList   http.HandlerFunc // GET  /v1/sessions/{id}/events
}

func newSessionEventsServer(t *testing.T) *sessionEventsServer {
	s := &sessionEventsServer{t: t}
	s.server = httptest.NewServer(http.HandlerFunc(s.serveHTTP))
	t.Cleanup(s.server.Close)
	return s
}

func (s *sessionEventsServer) serveHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasSuffix(r.URL.Path, "/events/stream"):
		require.NotNil(s.t, s.HandleStream, "unscripted StreamEvents call")
		s.HandleStream(w, r)
	case strings.HasSuffix(r.URL.Path, "/events"):
		switch r.Method {
		case http.MethodPost:
			require.NotNil(s.t, s.HandleSend, "unscripted events.Send call")
			s.HandleSend(w, r)
		case http.MethodGet:
			if s.HandleList != nil {
				s.HandleList(w, r)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(emptyEventList()))
		default:
			http.Error(w, "unexpected", http.StatusMethodNotAllowed)
		}
	default:
		s.t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		http.Error(w, "unexpected", http.StatusNotImplemented)
	}
}

func (s *sessionEventsServer) Client() Client {
	return NewClient(
		option.WithBaseURL(s.server.URL),
		option.WithAPIKey("test-key"),
		option.WithMaxRetries(0),
	)
}

// streamWriter writes scripted SSE events and optionally holds the connection
// open until the request context is cancelled.
func streamWriter(w http.ResponseWriter, r *http.Request, events []string, holdOpen bool) {
	w.Header().Set("Content-Type", "text/event-stream")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "no flusher", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	flusher.Flush()
	for _, ev := range events {
		if _, err := w.Write([]byte(ev)); err != nil {
			return
		}
		flusher.Flush()
	}
	if holdOpen {
		<-r.Context().Done()
	}
}

func newShortIdleRunner(t *testing.T, ctx context.Context, client Client, tools []BetaTool, maxIdle time.Duration) *SessionToolRunner {
	r := client.Beta.Sessions.Events.NewToolRunner(ctx, "sesn_test", SessionToolRunnerOptions{
		Tools:   tools,
		MaxIdle: &maxIdle,
		Logger:  sessionRunnerSilentLogger,
	})
	t.Cleanup(func() { _ = r.Close() })
	return r
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestSessionToolRunner_RequiresSessionID(t *testing.T) {
	client := NewClient(option.WithAPIKey("k"))
	r := client.Beta.Sessions.Events.NewToolRunner(context.Background(), "", SessionToolRunnerOptions{})
	require.False(t, r.Next())
	require.Error(t, r.Err())
	require.NoError(t, r.Close())
}

func TestSessionToolRunner_YieldsAndPostsResult(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, []string{
			sseLine("agent.tool_use", toolUseEvt("evt_1", "echo", map[string]any{"text": "hi"})),
		}, true)
	}
	var sent atomic.Int32
	server.HandleSend = func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		require.Contains(t, string(body), "evt_1")
		require.Contains(t, string(body), "user.tool_result")
		sent.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sendOK()))
	}

	echo := &stubBetaTool{name: "echo", run: func(_ context.Context, input json.RawMessage) (string, bool) {
		return "echoed: " + string(input), false
	}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 500*time.Millisecond)

	require.True(t, r.Next(), "first Next must yield the dispatched tool call")
	call := r.Current()
	require.Equal(t, "evt_1", call.ToolUseID)
	require.Equal(t, "echo", call.Name)
	require.False(t, call.IsError)
	require.True(t, call.Posted)
	require.Contains(t, dispatchedResultText(call), "echoed")
	require.Equal(t, int32(1), echo.runs.Load())
	require.Equal(t, int32(1), sent.Load())
	require.NoError(t, r.Close())
}

func TestSessionToolRunner_DispatchesCustomTool(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, []string{
			sseLine("agent.custom_tool_use", customToolUseEvt("cevt_1", "lookup", map[string]any{"q": "x"})),
		}, true)
	}
	var sent atomic.Int32
	server.HandleSend = func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		require.Contains(t, string(body), "user.custom_tool_result",
			"a custom tool call must be answered with user.custom_tool_result")
		require.Contains(t, string(body), "custom_tool_use_id")
		require.Contains(t, string(body), "cevt_1")
		require.NotContains(t, string(body), "user.tool_result")
		sent.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sendOK()))
	}

	lookup := &stubBetaTool{name: "lookup", run: func(_ context.Context, input json.RawMessage) (string, bool) {
		return "looked up: " + string(input), false
	}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{lookup}, 500*time.Millisecond)

	require.True(t, r.Next(), "first Next must yield the dispatched custom tool call")
	call := r.Current()
	require.True(t, call.Custom, "dispatch must be flagged as a custom tool call")
	require.Equal(t, "cevt_1", call.ToolUseID)
	require.Equal(t, "cevt_1", call.CustomToolUse.ID, "the embedded agent.custom_tool_use event must be populated")
	require.Equal(t, "lookup", call.Name)
	require.Equal(t, "cevt_1", call.CustomResult.CustomToolUseID, "the posted user.custom_tool_result must reference the event id")
	require.False(t, call.IsError)
	require.True(t, call.Posted)
	require.Contains(t, dispatchedResultText(call), "looked up")
	require.Equal(t, int32(1), lookup.runs.Load())
	require.Equal(t, int32(1), sent.Load())
	require.NoError(t, r.Close())
}

func TestSessionToolRunner_DispatchesBuiltinAndCustom(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, []string{
			sseLine("agent.tool_use", toolUseEvt("evt_b", "echo", map[string]any{})),
			sseLine("agent.custom_tool_use", customToolUseEvt("evt_c", "echo", map[string]any{})),
		}, true)
	}
	var builtinResults, customResults atomic.Int32
	server.HandleSend = func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		switch {
		case strings.Contains(string(body), "user.custom_tool_result"):
			require.Contains(t, string(body), "evt_c")
			customResults.Add(1)
		case strings.Contains(string(body), "user.tool_result"):
			require.Contains(t, string(body), "evt_b")
			builtinResults.Add(1)
		default:
			t.Errorf("send body carries no tool result: %s", body)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sendOK()))
	}

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 500*time.Millisecond)

	seen := map[string]DispatchedToolCall{}
	for range 2 {
		require.True(t, r.Next())
		call := r.Current()
		seen[call.ToolUseID] = call
	}
	require.NoError(t, r.Close())

	builtin, ok := seen["evt_b"]
	require.True(t, ok, "the builtin agent.tool_use was not dispatched")
	require.False(t, builtin.Custom)
	require.Equal(t, "evt_b", builtin.Result.ToolUseID)
	require.True(t, builtin.Posted)

	custom, ok := seen["evt_c"]
	require.True(t, ok, "the agent.custom_tool_use was not dispatched")
	require.True(t, custom.Custom)
	require.Equal(t, "evt_c", custom.CustomResult.CustomToolUseID)
	require.True(t, custom.Posted)

	require.Equal(t, int32(1), builtinResults.Load())
	require.Equal(t, int32(1), customResults.Load())
	require.Equal(t, int32(2), echo.runs.Load())
}

func TestSessionToolRunner_PermanentResultRejectionTerminates(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, []string{
			sseLine("agent.tool_use", toolUseEvt("evt_1", "echo", map[string]any{})),
		}, true)
	}

	var sends atomic.Int32
	server.HandleSend = func(w http.ResponseWriter, _ *http.Request) {
		sends.Add(1)
		http.Error(w, "bad", http.StatusBadRequest)
	}

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 0)

	require.True(t, r.Next(), "the rejected dispatch remains observable")
	require.False(t, r.Current().Posted)
	require.False(t, r.Next(), "a permanent result rejection must stop the runner")
	require.ErrorContains(t, r.Err(), "tool result send")
	require.Equal(t, int32(1), echo.runs.Load())
	require.Equal(t, int32(1), sends.Load())
	require.NoError(t, r.Close())
}

func TestSessionToolRunner_DrainsRejectedCallAfterTerminalCancel(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, []string{
			sseLine("agent.tool_use", toolUseEvt("evt_ok", "echo", nil)),
			sseLine("agent.tool_use", toolUseEvt("evt_rejected", "echo", nil)),
		}, true)
	}
	var sends atomic.Int32
	server.HandleSend = func(w http.ResponseWriter, _ *http.Request) {
		if sends.Add(1) == 1 {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(sendOK()))
			return
		}
		http.Error(w, "bad", http.StatusBadRequest)
	}

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 0)

	require.True(t, r.Next())
	require.Equal(t, "evt_ok", r.Current().ToolUseID)
	require.True(t, r.Current().Posted)
	select {
	case <-r.ctx.Done():
	case <-time.After(time.Second):
		t.Fatal("runner did not terminate after permanent result rejection")
	}
	require.True(t, r.Next(), "a buffered rejected call must be drained before cancellation ends iteration")
	require.Equal(t, "evt_rejected", r.Current().ToolUseID)
	require.False(t, r.Current().Posted)
	require.False(t, r.Next())
	require.ErrorContains(t, r.Err(), "tool result send")
	require.Equal(t, int32(2), echo.runs.Load())
	require.Equal(t, int32(2), sends.Load())
	require.NoError(t, r.Close())
}

func TestSessionToolRunner_SurfaceCallPrefersBufferAfterCancellation(t *testing.T) {
	r := &SessionToolRunner{results: make(chan DispatchedToolCall, 1)}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	r.surfaceCall(ctx, DispatchedToolCall{ToolUseID: "evt_done"})

	select {
	case call := <-r.results:
		require.Equal(t, "evt_done", call.ToolUseID)
	default:
		t.Fatal("completed call was dropped despite available result buffer")
	}
}

func TestSessionToolRunner_ReconcileRetriesTransientPostFailure(t *testing.T) {
	server := newSessionEventsServer(t)
	var streamConns atomic.Int32
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		if streamConns.Add(1) == 1 {
			streamWriter(w, r, []string{
				sseLine("agent.tool_use", toolUseEvt("evt_retry", "echo", nil)),
			}, false)
			return
		}
		streamWriter(w, r, nil, true)
	}
	var lists atomic.Int32
	server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if lists.Add(1) == 1 {
			_, _ = w.Write([]byte(emptyEventList()))
			return
		}
		body, _ := json.Marshal(map[string]any{
			"data":     []any{toolUseEvt("evt_retry", "echo", nil)},
			"first_id": "evt_retry", "has_more": false, "last_id": "evt_retry",
		})
		_, _ = w.Write(body)
	}
	var sends atomic.Int32
	server.HandleSend = func(w http.ResponseWriter, _ *http.Request) {
		if sends.Add(1) <= sessionRunnerSendRetries {
			http.Error(w, "retry", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sendOK()))
	}

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 0)

	require.True(t, r.Next())
	require.False(t, r.Current().Posted)
	require.True(t, r.Next())
	require.True(t, r.Current().Posted)
	require.Equal(t, int32(2), echo.runs.Load())
	require.Equal(t, int32(4), sends.Load())
	require.NoError(t, r.Close())
}

func TestSessionToolRunner_GlobalInterruptCancelsHeldHistory(t *testing.T) {
	server := newSessionEventsServer(t)
	var streamConns atomic.Int32
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		if streamConns.Add(1) == 1 {
			streamWriter(w, r, []string{
				sseLine("agent.tool_use", askToolUseEvt("evt_held", "echo", nil, "ask")),
			}, false)
			return
		}
		streamWriter(w, r, nil, true)
	}

	var lists atomic.Int32
	server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if lists.Add(1) == 1 {
			_, _ = w.Write([]byte(emptyEventList()))
			return
		}
		events := []any{interruptEvt("evt_interrupt", true, ""), idleEndTurnEvt("evt_idle")}
		body, _ := json.Marshal(map[string]any{
			"data": events, "first_id": "evt_interrupt", "has_more": false, "last_id": "evt_idle",
		})
		_, _ = w.Write(body)
	}
	server.HandleSend = func(http.ResponseWriter, *http.Request) {
		t.Fatal("an interrupted tool call must not post a result")
	}

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 50*time.Millisecond)

	require.False(t, r.Next())
	require.ErrorIs(t, r.Err(), ErrIdleTimeout)
	require.Equal(t, int32(0), echo.runs.Load())
	require.False(t, r.isAnswered("evt_held"))
	require.True(t, r.isSettled("evt_held"))
	require.NoError(t, r.Close())
}

func TestSessionToolRunner_GlobalInterruptCancelsHeldLiveCall(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, []string{
			sseLine("agent.tool_use", askToolUseEvt("evt_held", "echo", nil, "ask")),
			sseLine("user.interrupt", interruptEvt("evt_interrupt", true, "")),
			sseLine("session.status_idle", idleEndTurnEvt("evt_idle")),
		}, true)
	}
	server.HandleSend = func(http.ResponseWriter, *http.Request) {
		t.Fatal("an interrupted tool call must not post a result")
	}

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 50*time.Millisecond)

	require.False(t, r.Next())
	require.ErrorIs(t, r.Err(), ErrIdleTimeout)
	require.Equal(t, int32(0), echo.runs.Load())
	require.False(t, r.isAnswered("evt_held"))
	require.True(t, r.isSettled("evt_held"))
	require.NoError(t, r.Close())
}

func TestSessionToolRunner_GlobalInterruptUsesProcessedTimeCutoff(t *testing.T) {
	tests := []struct {
		name      string
		toolUseID string
		toolUse   map[string]any
	}{
		{
			name:      "builtin",
			toolUseID: "evt_stale",
			toolUse:   toolUseEvt("evt_stale", "echo", nil),
		},
		{
			name:      "custom",
			toolUseID: "cevt_stale",
			toolUse:   customToolUseEvt("cevt_stale", "echo", nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newSessionEventsServer(t)
			server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
				streamWriter(w, r, nil, true)
			}
			server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
				// History is ordered by created_at. A queued interrupt can therefore
				// appear before a tool call that it later invalidated at processing time.
				events := []any{interruptEvt("evt_interrupt", true, ""), tt.toolUse, idleEndTurnEvt("evt_idle")}
				body, _ := json.Marshal(map[string]any{
					"data": events, "first_id": "evt_interrupt", "has_more": false, "last_id": "evt_idle",
				})
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(body)
			}
			server.HandleSend = func(http.ResponseWriter, *http.Request) {
				t.Fatal("a tool call before the processed interrupt cutoff must not post a result")
			}

			echo := &stubBetaTool{name: "echo"}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 50*time.Millisecond)

			require.False(t, r.Next())
			require.ErrorIs(t, r.Err(), ErrIdleTimeout)
			require.Equal(t, int32(0), echo.runs.Load())
			require.True(t, r.isCanceled(tt.toolUseID))
			require.NoError(t, r.Close())
		})
	}
}

func TestSessionToolRunner_GlobalInterruptCancelsInFlightCall(t *testing.T) {
	server := newSessionEventsServer(t)
	started := make(chan struct{})
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher := w.(http.Flusher)
		w.WriteHeader(http.StatusOK)
		flusher.Flush()
		_, _ = w.Write([]byte(sseLine("agent.tool_use", toolUseEvt("evt_running", "echo", nil))))
		flusher.Flush()
		select {
		case <-started:
		case <-r.Context().Done():
			return
		}
		_, _ = w.Write([]byte(sseLine("user.interrupt", interruptEvt("evt_interrupt", true, ""))))
		flusher.Flush()
		_, _ = w.Write([]byte(sseLine("session.status_idle", idleEndTurnEvt("evt_idle"))))
		flusher.Flush()
		<-r.Context().Done()
	}
	server.HandleSend = func(http.ResponseWriter, *http.Request) {
		t.Fatal("an interrupted in-flight tool call must not post a result")
	}

	echo := &stubBetaTool{name: "echo", run: func(ctx context.Context, _ json.RawMessage) (string, bool) {
		close(started)
		<-ctx.Done()
		return ctx.Err().Error(), true
	}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 50*time.Millisecond)

	require.True(t, r.Next(), "the interrupted dispatch remains observable")
	require.Equal(t, "evt_running", r.Current().ToolUseID)
	require.False(t, r.Current().Posted)
	require.False(t, r.Next())
	require.ErrorIs(t, r.Err(), ErrIdleTimeout)
	require.Equal(t, int32(1), echo.runs.Load())
	require.True(t, r.isCanceled("evt_running"))
	require.NoError(t, r.Close())
}

func TestSessionToolRunner_PermanentRejectionRefreshesSettlementHistory(t *testing.T) {
	server := newSessionEventsServer(t)
	sendStarted := make(chan struct{})
	releaseIdle := make(chan struct{})
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher := w.(http.Flusher)
		w.WriteHeader(http.StatusOK)
		flusher.Flush()
		_, _ = w.Write([]byte(sseLine("agent.tool_use", toolUseEvt("evt_race", "echo", nil))))
		_, _ = w.Write([]byte(sseLine("agent.tool_use", toolUseEvt("evt_queued_answered", "echo", nil))))
		_, _ = w.Write([]byte(sseLine("agent.tool_use", toolUseEvt("evt_queued_interrupted", "echo", nil))))
		flusher.Flush()
		select {
		case <-sendStarted:
		case <-r.Context().Done():
			return
		}
		select {
		case <-releaseIdle:
		case <-r.Context().Done():
			return
		}
		_, _ = w.Write([]byte(sseLine("session.status_idle", idleEndTurnEvt("evt_idle"))))
		flusher.Flush()
		<-r.Context().Done()
	}

	var lists atomic.Int32
	server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if lists.Add(1) == 1 {
			_, _ = w.Write([]byte(emptyEventList()))
			return
		}
		result := map[string]any{
			"type":         "user.tool_result",
			"id":           "evt_result",
			"tool_use_id":  "evt_race",
			"content":      []any{map[string]any{"type": "text", "text": "answered elsewhere"}},
			"is_error":     false,
			"processed_at": "2026-05-11T12:00:01Z",
		}
		queuedResult := map[string]any{
			"type":         "user.tool_result",
			"id":           "evt_queued_result",
			"tool_use_id":  "evt_queued_answered",
			"content":      []any{map[string]any{"type": "text", "text": "answered elsewhere"}},
			"is_error":     false,
			"processed_at": "2026-05-11T12:00:01Z",
		}
		interrupt := interruptEvt("evt_interrupt", true, "")
		body, _ := json.Marshal(map[string]any{
			"data": []any{result, queuedResult, interrupt}, "first_id": "evt_result", "has_more": false, "last_id": "evt_interrupt",
		})
		_, _ = w.Write(body)
	}
	var sends atomic.Int32
	server.HandleSend = func(w http.ResponseWriter, _ *http.Request) {
		if sends.Add(1) == 1 {
			close(sendStarted)
		}
		http.Error(w, "interrupted", http.StatusBadRequest)
	}

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	runner := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 50*time.Millisecond)

	require.True(t, runner.Next(), "the interrupted dispatch remains observable")
	require.False(t, runner.Current().Posted)
	close(releaseIdle)
	require.False(t, runner.Next())
	require.ErrorIs(t, runner.Err(), ErrIdleTimeout)
	require.Equal(t, int32(1), echo.runs.Load())
	require.Equal(t, int32(1), sends.Load())
	require.Equal(t, int32(2), lists.Load())
	require.True(t, runner.isAnswered("evt_race"))
	require.False(t, runner.isCanceled("evt_race"))
	require.True(t, runner.isAnswered("evt_queued_answered"))
	require.False(t, runner.isCanceled("evt_queued_answered"))
	require.True(t, runner.isCanceled("evt_queued_interrupted"))
	require.NoError(t, runner.Close())
}

func TestSessionToolRunner_ExternalResultReleasesConfirmationHold(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		result := map[string]any{
			"type":         "user.tool_result",
			"id":           "evt_result",
			"tool_use_id":  "evt_held",
			"content":      []any{map[string]any{"type": "text", "text": "answered elsewhere"}},
			"is_error":     false,
			"processed_at": "2026-05-11T12:00:01Z",
		}
		streamWriter(w, r, []string{
			sseLine("agent.tool_use", askToolUseEvt("evt_held", "echo", nil, "ask")),
			sseLine("user.tool_result", result),
			sseLine("session.status_idle", idleEndTurnEvt("evt_idle")),
		}, true)
	}
	server.HandleSend = func(http.ResponseWriter, *http.Request) {
		t.Fatal("a tool call answered elsewhere must not post another result")
	}

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 50*time.Millisecond)

	require.False(t, r.Next())
	require.ErrorIs(t, r.Err(), ErrIdleTimeout)
	require.Equal(t, int32(0), echo.runs.Load())
	require.True(t, r.isAnswered("evt_held"))
	require.NoError(t, r.Close())
}

func TestSessionToolRunner_HistoryInterruptIsNotReappliedFromStream(t *testing.T) {
	server := newSessionEventsServer(t)
	interrupt := interruptEvt("evt_interrupt", true, "")
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, []string{sseLine("user.interrupt", interrupt)}, true)
	}
	server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
		tool := toolUseEvt("evt_after", "echo", nil)
		tool["processed_at"] = "2026-05-11T12:00:02Z"
		body, _ := json.Marshal(map[string]any{
			"data": []any{interrupt, tool}, "first_id": "evt_interrupt", "has_more": false, "last_id": "evt_after",
		})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}
	var sends atomic.Int32
	server.HandleSend = func(w http.ResponseWriter, _ *http.Request) {
		sends.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sendOK()))
	}

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 0)

	require.True(t, r.Next())
	require.Equal(t, "evt_after", r.Current().ToolUseID)
	require.True(t, r.Current().Posted)
	require.Equal(t, int32(1), echo.runs.Load())
	require.Equal(t, int32(1), sends.Load())
	require.NoError(t, r.Close())
}

func TestSessionToolRunner_DoesNotCancelForTargetedOrQueuedInterrupt(t *testing.T) {
	tests := map[string]map[string]any{
		"targeted": interruptEvt("evt_interrupt", true, "sthr_other"),
		"queued":   interruptEvt("evt_interrupt", false, ""),
	}
	for name, interrupt := range tests {
		t.Run(name, func(t *testing.T) {
			server := newSessionEventsServer(t)
			server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
				streamWriter(w, r, nil, true)
			}
			server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				events := []any{toolUseEvt("evt_live", "echo", nil), interrupt, idleEndTurnEvt("evt_idle")}
				body, _ := json.Marshal(map[string]any{
					"data": events, "first_id": "evt_live", "has_more": false, "last_id": "evt_idle",
				})
				_, _ = w.Write(body)
			}
			var sends atomic.Int32
			server.HandleSend = func(w http.ResponseWriter, _ *http.Request) {
				sends.Add(1)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(sendOK()))
			}

			echo := &stubBetaTool{name: "echo"}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 0)

			require.True(t, r.Next())
			require.True(t, r.Current().Posted)
			require.Equal(t, int32(1), echo.runs.Load())
			require.Equal(t, int32(1), sends.Load())
			require.NoError(t, r.Close())
		})
	}
}

// TestSessionToolRunner_SkipsUnownedToolByDefault pins the default
// split-client behavior: a tool-use event whose Name is not in the runner's
// registry belongs to the other client servicing the session (e.g. the
// customer's app backend handling custom tools). The runner must post NO
// result for it, claim nothing, and leave the tool_use_id pending — while
// still yielding the DispatchedToolCall so the caller can observe the unowned
// dispatch (Posted=false, IsError=false, no result event populated). It must
// not panic on the registry miss.
func TestSessionToolRunner_SkipsUnownedToolByDefault(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, []string{
			sseLine("agent.tool_use", toolUseEvt("evt_99", "not_ours", map[string]any{})),
			sseLine("agent.custom_tool_use", customToolUseEvt("cevt_99", "app_backend_tool", map[string]any{})),
		}, true)
	}
	server.HandleSend = func(http.ResponseWriter, *http.Request) {
		t.Fatal("runner must not post any result for a tool it does not own")
	}

	// The runner owns "echo" but neither streamed event is for it.
	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 500*time.Millisecond)

	seen := map[string]DispatchedToolCall{}
	for range 2 {
		require.True(t, r.Next(), "the unowned call must still be surfaced to the caller")
		call := r.Current()
		seen[call.ToolUseID] = call
	}

	builtin, ok := seen["evt_99"]
	require.True(t, ok, "the unowned agent.tool_use must still be yielded")
	require.Equal(t, "not_ours", builtin.Name)
	require.False(t, builtin.Custom)
	require.False(t, builtin.IsError, "a skipped call is not an error")
	require.False(t, builtin.Posted, "nothing was sent for an unowned tool")
	require.Equal(t, "evt_99", builtin.ToolUse.ID, "the triggering event is still surfaced")
	require.Empty(t, builtin.Result.ToolUseID, "no user.tool_result was ever built")
	require.Empty(t, dispatchedResultText(builtin))

	custom, ok := seen["cevt_99"]
	require.True(t, ok, "the unowned agent.custom_tool_use must still be yielded")
	require.Equal(t, "app_backend_tool", custom.Name)
	require.True(t, custom.Custom)
	require.False(t, custom.IsError, "a skipped call is not an error")
	require.False(t, custom.Posted, "nothing was sent for an unowned custom tool")
	require.Equal(t, "cevt_99", custom.CustomToolUse.ID, "the triggering event is still surfaced")
	require.Empty(t, custom.CustomResult.CustomToolUseID, "no user.custom_tool_result was ever built")
	require.Empty(t, dispatchedResultText(custom))

	require.Equal(t, int32(0), echo.runs.Load(), "no registered tool should have run")
	require.NoError(t, r.Close())
	require.NoError(t, r.Err(), "skipping an unowned tool is not a terminal error")
}

// A skipped (unanswered) unowned tool_use stays OUT of the end-turn
// accounting: reconcile sees history ending on an end_turn idle but with the
// unowned tool_use still unanswered, so it must NOT arm the countdown — the
// runner has not handled that call, its owner still has to.
func TestSessionToolRunner_SkippedUnownedToolDoesNotTripIdle(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, nil, true) // no live events; reconcile drives the test
	}
	server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{
			"data": []any{
				toolUseEvt("evt_pending", "not_ours", map[string]any{}),
				idleEndTurnEvt("evt_idle"),
			},
			"first_id": "evt_pending", "has_more": false, "last_id": "evt_idle",
		})
		_, _ = w.Write(body)
	}
	server.HandleSend = func(http.ResponseWriter, *http.Request) {
		t.Fatal("runner must not post any result for a tool it does not own")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), nil, 100*time.Millisecond)

	require.True(t, r.Next(), "reconcile must still surface the unowned call")
	call := r.Current()
	require.Equal(t, "evt_pending", call.ToolUseID)
	require.False(t, call.Posted)
	require.False(t, call.IsError)
	require.Empty(t, call.Result.ToolUseID, "no result was built for the skipped call")

	for r.Next() {
		t.Fatalf("unexpected extra yield: %+v", r.Current())
	}
	require.NotErrorIs(t, r.Err(), ErrIdleTimeout,
		"runner idled out with an unowned tool_use still unanswered — reconcile must not arm over outstanding work")
}

func TestSessionToolRunner_SessionTerminatedEndsIteration(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, []string{
			sseLine("session.status_terminated", plainEvt("evt_term", "session.status_terminated")),
		}, false)
	}
	server.HandleSend = func(http.ResponseWriter, *http.Request) {
		t.Fatal("send must not be called when no tool_use ever fires")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), nil, 0)
	for r.Next() {
		t.Fatalf("unexpected yield: %+v", r.Current())
	}
	require.ErrorIs(t, r.Err(), ErrSessionTerminated)
}

// TestSessionToolRunner_ReconcileSurfacesSessionTerminatedFromHistory pins
// the guarantee that reconcile() shuts down streamLoop when the listed
// history contains a session.status_terminated or session.deleted event.
// The live stream only replays current events, not history, so a session
// that terminated before the runner attached is only visible to reconcile;
// without surfacing it, streamLoop would reconnect forever against a dead
// session after the (eventless) live stream disconnected.
//
// The test scripts an event-list response containing a single
// session.status_terminated entry and a stream handler that simply holds
// the connection open (no events ever fire). reconcile() returns
// ErrSessionTerminated, streamLoop propagates it, and the runner exits
// with ErrSessionTerminated well within the outer ctx.
func TestSessionToolRunner_ReconcileSurfacesSessionTerminatedFromHistory(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		// Stream is eager-opened by streamLoop BEFORE reconcile, so this
		// handler always runs. With the fix, reconcile then returns the
		// sentinel and the stream is closed; without the fix, the runner
		// sits here forever waiting for events that never come.
		streamWriter(w, r, nil, true)
	}
	server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{
			"data":     []any{plainEvt("evt_term", "session.status_terminated")},
			"first_id": "evt_term",
			"has_more": false,
			"last_id":  "evt_term",
		})
		_, _ = w.Write(body)
	}
	server.HandleSend = func(http.ResponseWriter, *http.Request) {
		t.Fatal("send must not be called when no tool_use ever fires")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), nil, 0)

	start := time.Now()
	for r.Next() {
		t.Fatalf("unexpected yield: %+v", r.Current())
	}
	elapsed := time.Since(start)

	require.ErrorIs(t, r.Err(), ErrSessionTerminated, "reconcile must surface ErrSessionTerminated when the listed history contains session.status_terminated")
	require.Less(t, elapsed, 3*time.Second, "runner ran too long — reconcile didn't shut down streamLoop (got %s)", elapsed)
}

func TestSessionToolRunner_IdleTimeoutEndsIteration(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, []string{sseLine("session.status_idle", idleEndTurnEvt("evt_idle"))}, true)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), nil, 100*time.Millisecond)
	for r.Next() {
		t.Fatalf("unexpected yield: %+v", r.Current())
	}
	require.ErrorIs(t, r.Err(), ErrIdleTimeout)
}

func TestSessionToolRunner_SendRetriesOnTransientError(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, []string{
			sseLine("agent.tool_use", toolUseEvt("evt_1", "echo", map[string]any{})),
		}, true)
	}
	var attempts atomic.Int32
	server.HandleSend = func(w http.ResponseWriter, _ *http.Request) {
		if attempts.Add(1) == 1 {
			http.Error(w, "boom", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sendOK()))
	}

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	// Shrink the per-attempt send timeout via the options surface so the test
	// stays quick without mutating package-level state.
	maxIdle := 500 * time.Millisecond
	client := server.Client()
	r := client.Beta.Sessions.Events.NewToolRunner(ctx, "sesn_test", SessionToolRunnerOptions{
		Tools:       []BetaTool{echo},
		MaxIdle:     &maxIdle,
		Logger:      sessionRunnerSilentLogger,
		SendTimeout: 2 * time.Second,
	})
	t.Cleanup(func() { _ = r.Close() })

	require.True(t, r.Next())
	call := r.Current()
	require.True(t, call.Posted, "send should succeed after a retry")
	require.GreaterOrEqual(t, attempts.Load(), int32(2))
}

// TestToToolResultContent_BlockPassthrough verifies that block kinds shared by
// both unions round-trip natively, and block kinds the destination union
// cannot represent fall back to a text block holding the raw JSON.
func TestToToolResultContent_BlockPassthrough(t *testing.T) {
	src := []BetaToolResultBlockParamContentUnion{
		{OfText: &BetaTextBlockParam{Text: "hello"}},
		{OfImage: &BetaImageBlockParam{
			Source: BetaImageBlockParamSourceUnion{
				OfBase64: &BetaBase64ImageSourceParam{
					Data:      "QUJD",
					MediaType: BetaBase64ImageSourceMediaTypeImagePNG,
				},
			},
		}},
		{OfDocument: &BetaRequestDocumentBlockParam{
			Source: BetaRequestDocumentBlockSourceUnionParam{
				OfBase64: &BetaBase64PDFSourceParam{Data: "UERG"},
			},
		}},
		{OfDocument: &BetaRequestDocumentBlockParam{
			Source: BetaRequestDocumentBlockSourceUnionParam{
				OfText: &BetaPlainTextSourceParam{Data: "plain text doc"},
			},
		}},
		{OfDocument: &BetaRequestDocumentBlockParam{
			Source: BetaRequestDocumentBlockSourceUnionParam{
				OfContent: &BetaContentBlockSourceParam{
					Content: BetaContentBlockSourceContentUnionParam{
						OfString: param.NewOpt("nested content"),
					},
				},
			},
		}},
		{OfSearchResult: &BetaSearchResultBlockParam{
			Source:  "https://example.com",
			Title:   "result",
			Content: []BetaTextBlockParam{{Text: "snippet"}},
		}},
		{OfToolReference: &BetaToolReferenceBlockParam{ToolName: "weather"}},
	}

	got := toToolResultContent(src)
	require.Len(t, got, 7)

	require.NotNil(t, got[0].OfText)
	require.Equal(t, "hello", got[0].OfText.Text, "text block passes through unchanged")

	require.NotNil(t, got[1].OfImage, "base64 image block passes through natively")
	require.Nil(t, got[1].OfText)
	require.NotNil(t, got[1].OfImage.Source.OfBase64)
	require.Equal(t, "QUJD", got[1].OfImage.Source.OfBase64.Data)

	require.NotNil(t, got[2].OfDocument, "base64 document block passes through natively")
	require.Nil(t, got[2].OfText)
	require.NotNil(t, got[2].OfDocument.Source.OfBase64)
	require.Equal(t, "UERG", got[2].OfDocument.Source.OfBase64.Data)

	require.NotNil(t, got[3].OfDocument, "plain-text document block passes through natively")
	require.Nil(t, got[3].OfText)
	require.NotNil(t, got[3].OfDocument.Source.OfText)
	require.Equal(t, "plain text doc", got[3].OfDocument.Source.OfText.Data)

	require.NotNil(t, got[4].OfText, "content-source document block falls back to stringified JSON")
	require.Nil(t, got[4].OfDocument)
	require.Contains(t, got[4].OfText.Text, `"type":"document"`)
	require.Contains(t, got[4].OfText.Text, "nested content")

	require.NotNil(t, got[5].OfSearchResult, "search_result block passes through natively")
	require.Nil(t, got[5].OfText)
	require.Equal(t, "https://example.com", got[5].OfSearchResult.Source)
	require.Equal(t, "result", got[5].OfSearchResult.Title)

	require.NotNil(t, got[6].OfText, "tool_reference block falls back to stringified JSON")
	require.Contains(t, got[6].OfText.Text, "tool_reference")
	require.Contains(t, got[6].OfText.Text, "weather")

	// The custom-tool variant must produce equivalent content.
	custom := toCustomToolResultContent(src)
	require.Len(t, custom, 7)
	for i := range got {
		gotJSON, err := json.Marshal(got[i])
		require.NoError(t, err)
		customJSON, err := json.Marshal(custom[i])
		require.NoError(t, err)
		require.JSONEq(t, string(gotJSON), string(customJSON))
	}
	require.NotNil(t, custom[1].OfImage)
	require.NotNil(t, custom[2].OfDocument)
	require.NotNil(t, custom[5].OfSearchResult)
}

// ===== tool calls that need user approval =====
//
// The server marks each agent.tool_use with evaluated_permission. Only "allow"
// (or no mark) runs on arrival; "ask" waits for the user's
// user.tool_confirmation; "deny" and unrecognized values never run.

func TestSessionToolRunner_OnlyToolsTheServerOrUserAllowedEverRun(t *testing.T) {
	server := newSessionEventsServer(t)
	// From the history endpoint (read once when the runner connects):
	server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
		body, _ := json.Marshal(map[string]any{
			"data": []any{
				askToolUseEvt("ask_answered_in_history", "echo", nil, "ask"),
				toolConfirmationEvt("ask_answered_in_history", "allow"),
			},
			"has_more": false, "first_id": "ask_answered_in_history", "last_id": "evt_conf_ask_answered_in_history",
		})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}
	// From the live event stream (held open — the dispatcher is pipelined, so a
	// trailing terminated would race the last dispatch):
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, []string{
			sseLine("agent.tool_use", askToolUseEvt("marked_allow", "echo", nil, "allow")),
			sseLine("agent.tool_use", askToolUseEvt("marked_deny", "echo", nil, "deny")),
			sseLine("agent.tool_use", askToolUseEvt("ask_then_allowed", "echo", nil, "ask")),
			sseLine("agent.tool_use", askToolUseEvt("ask_then_denied", "echo", nil, "ask")),
			sseLine("agent.tool_use", askToolUseEvt("ask_never_answered", "echo", nil, "ask")),
			sseLine("agent.tool_use", askToolUseEvt("ask_verdict_unrecognized", "echo", nil, "ask")),
			sseLine("agent.tool_use", askToolUseEvt("unrecognized_mark", "echo", nil, "something_new")),
			sseLine("user.tool_confirmation", toolConfirmationEvt("ask_then_allowed", "allow")),
			sseLine("user.tool_confirmation", toolConfirmationEvt("ask_then_denied", "deny")),
			sseLine("user.tool_confirmation", toolConfirmationEvt("ask_verdict_unrecognized", "escalate")),
		}, true)
	}
	var posted []string
	var mu sync.Mutex
	server.HandleSend = func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload struct {
			Events []struct {
				ToolUseID string `json:"tool_use_id"`
			} `json:"events"`
		}
		_ = json.Unmarshal(body, &payload)
		mu.Lock()
		for _, e := range payload.Events {
			posted = append(posted, e.ToolUseID)
		}
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sendOK()))
	}

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 0)

	// Allowed calls execute and post; denied calls (server "deny", user "deny",
	// unrecognised verdict) are yielded with Confirmation="deny" but never run.
	// Calls whose confirmation never arrives stay held and are not yielded.
	got := map[string]DispatchedToolCall{}
	for range 6 {
		require.True(t, r.Next(), "expected six resolved calls; got %v then stopped: %v", got, r.Err())
		c := r.Current()
		got[c.ToolUseID] = c
	}
	require.NoError(t, r.Close())

	allowed := []string{"ask_answered_in_history", "ask_then_allowed", "marked_allow"}
	denied := []string{"ask_then_denied", "ask_verdict_unrecognized", "marked_deny"}
	var ran, blocked []string
	for id, c := range got {
		if c.Confirmation == "deny" {
			blocked = append(blocked, id)
			require.False(t, c.Posted, "%s: a denied call posts nothing", id)
		} else {
			ran = append(ran, id)
			require.True(t, c.Posted, "%s: an allowed call posts its result", id)
		}
	}
	sort.Strings(ran)
	sort.Strings(blocked)
	sort.Strings(posted)
	require.Equal(t, allowed, ran, "only explicitly allowed calls execute")
	require.Equal(t, denied, blocked, "denied calls are yielded with Confirmation=\"deny\"")
	require.Equal(t, allowed, posted)
	require.Equal(t, int32(3), echo.runs.Load())
	require.NotContains(t, got, "ask_never_answered", "a call whose confirmation never arrives stays held")
	require.NotContains(t, got, "unrecognized_mark", "an unrecognised permission with no verdict stays held")
}

// A denied call is still an outcome the consumer needs to observe: the
// agent tried to invoke a tool and was blocked. Dropping it silently makes
// "the agent called nothing" indistinguishable from "the agent called five
// tools and the user denied every one", which breaks audit trails and any
// UI that surfaces per-call outcomes. So a denied call — whether the server
// evaluated permission to "deny" or the user's confirmation verdict was a
// deny — must be yielded with DispatchedToolCall.Confirmation == "deny"
// (Posted=false, IsError=false, no result), and an ask-then-allow call must
// carry Confirmation == "allow" so the consumer can tell it was gated.
func TestSessionToolRunner_DeniedCallsAreYieldedWithConfirmation(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
		body, _ := json.Marshal(map[string]any{
			"data": []any{
				askToolUseEvt("srv_deny", "echo", nil, "deny"),
				askToolUseEvt("usr_deny", "echo", nil, "ask"),
				askToolUseEvt("usr_allow", "echo", nil, "ask"),
				idleRequiresActionEvt("evt_ra", "usr_deny", "usr_allow"),
			},
			"has_more": false, "first_id": "srv_deny", "last_id": "evt_ra",
		})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, []string{
			sseLine("user.tool_confirmation", toolConfirmationEvt("usr_deny", "deny")),
			sseLine("user.tool_confirmation", toolConfirmationEvt("usr_allow", "allow")),
			sseLine("session.status_idle", idleEndTurnEvt("evt_et")),
		}, true)
	}
	var sends atomic.Int32
	server.HandleSend = func(w http.ResponseWriter, _ *http.Request) {
		sends.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sendOK()))
	}

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 200*time.Millisecond)

	got := map[string]DispatchedToolCall{}
	for r.Next() {
		c := r.Current()
		got[c.ToolUseID] = c
	}
	require.ErrorIs(t, r.Err(), ErrIdleTimeout)

	require.Contains(t, got, "srv_deny", "server-denied call was not yielded")
	require.Contains(t, got, "usr_deny", "user-denied call was not yielded")
	require.Contains(t, got, "usr_allow", "user-allowed call was not yielded")

	for _, id := range []string{"srv_deny", "usr_deny"} {
		c := got[id]
		require.Equal(t, "deny", c.Confirmation, "%s: Confirmation must be \"deny\"", id)
		require.False(t, c.Posted, "%s: a denied call posts no result", id)
		require.False(t, c.IsError, "%s: a denial is not a tool error", id)
		require.Empty(t, c.Result.ToolUseID, "%s: no result was built for a denied call", id)
	}
	require.Equal(t, "allow", got["usr_allow"].Confirmation, "an ask-then-allow call must carry Confirmation == \"allow\"")
	require.True(t, got["usr_allow"].Posted)

	require.Equal(t, int32(1), echo.runs.Load(), "only the allowed call executes")
	require.Equal(t, int32(1), sends.Load(), "only the allowed call posts a result")
}

// A verdict this SDK cannot read must fail closed, and a server-side deny
// outranks a stray allow. Each verdict here lands BEFORE its tool_use, so the
// gate resolves it out of the recorded verdicts rather than through the
// held-call path — and that read must tell "no verdict recorded" apart from
// "a verdict recorded that we can't interpret", which a bare map lookup on
// Go's zero value does not.
func TestSessionToolRunner_GateFailsClosedOnStrayVerdicts(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, []string{
			// Verdicts first: each is recorded before its call is routed.
			sseLine("user.tool_confirmation", toolConfirmationEvt("ungated_stray_deny", "deny")),
			sseLine("user.tool_confirmation", toolConfirmationEvt("ungated_stray_empty", "")),
			sseLine("user.tool_confirmation", toolConfirmationEvt("pre_denied", "allow")),

			sseLine("agent.tool_use", toolUseEvt("ungated_stray_deny", "echo", nil)),
			sseLine("agent.tool_use", toolUseEvt("ungated_stray_empty", "echo", nil)),
			sseLine("agent.tool_use", askToolUseEvt("pre_denied", "echo", nil, "deny")),
			sseLine("agent.tool_use", askToolUseEvt("marked_allow", "echo", nil, "allow")),
			sseLine("session.status_idle", idleEndTurnEvt("evt_idle")),
		}, true)
	}
	server.HandleSend = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sendOK()))
	}

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 200*time.Millisecond)

	got := map[string]DispatchedToolCall{}
	for r.Next() {
		c := r.Current()
		got[c.ToolUseID] = c
	}
	require.ErrorIs(t, r.Err(), ErrIdleTimeout, "nothing is held, so the end_turn countdown must run")

	require.Equal(t, int32(1), echo.runs.Load(), "only marked_allow may execute")
	require.True(t, got["marked_allow"].Posted)
	require.Empty(t, got["marked_allow"].Confirmation, "an allow-marked call needed no confirmation")

	for _, id := range []string{"ungated_stray_deny", "ungated_stray_empty", "pre_denied"} {
		c, ok := got[id]
		require.True(t, ok, "%s: a resolved-denied call must still be yielded", id)
		require.Equal(t, "deny", c.Confirmation, "%s: must fail closed as a denial", id)
		require.False(t, c.Posted, "%s: a denied call posts no result", id)
	}
}

// A permission value this SDK does not recognize must hold the call like
// "ask" — never dispatch it unconfirmed — which also defers the idle
// countdown for as long as it stays held.
func TestSessionToolRunner_UnrecognizedPermissionHoldsTheCall(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, []string{
			sseLine("agent.tool_use", askToolUseEvt("tu", "echo", nil, "something_new")),
			sseLine("session.status_idle", idleEndTurnEvt("evt_idle")),
		}, true)
	}
	server.HandleSend = func(http.ResponseWriter, *http.Request) { t.Error("a held call posted a result") }

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 100*time.Millisecond)

	for r.Next() {
		t.Fatalf("a held call must not be yielded: %+v", r.Current())
	}
	require.NotErrorIs(t, r.Err(), ErrIdleTimeout,
		"the held call must defer the countdown; the runner ends on ctx, not MaxIdle")
	require.Zero(t, echo.runs.Load())
}

// A call held from the live stream whose allow verdict only ever appears in a
// reconcile history — the tool_use itself scrolled out of the listed window —
// must still be dispatched, exactly once.
func TestSessionToolRunner_ReconcileVerdictResolvesCallMissingFromHistory(t *testing.T) {
	server := newSessionEventsServer(t)

	var streamConns atomic.Int32
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		if streamConns.Add(1) == 1 {
			// First connection: the gated call arrives live (and is held), then
			// the stream closes so the runner reconnects and reconciles.
			streamWriter(w, r, []string{
				sseLine("agent.tool_use", askToolUseEvt("tu", "echo", nil, "ask")),
			}, false)
			return
		}
		streamWriter(w, r, nil, true) // later connections: hold open, no events
	}

	var lists atomic.Int32
	server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// The reconcile after the reconnect sees only the verdict: the
		// tool_use itself is outside the listed window.
		if lists.Add(1) >= 2 {
			body, _ := json.Marshal(map[string]any{
				"data":     []any{toolConfirmationEvt("tu", "allow")},
				"has_more": false, "first_id": "evt_conf_tu", "last_id": "evt_conf_tu",
			})
			_, _ = w.Write(body)
			return
		}
		_, _ = w.Write([]byte(emptyEventList()))
	}

	server.HandleSend = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sendOK()))
	}

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 0)

	require.True(t, r.Next(), "the approved held call must be dispatched: %v", r.Err())
	call := r.Current()
	require.Equal(t, "tu", call.ToolUseID)
	require.Equal(t, "allow", call.Confirmation)
	require.True(t, call.Posted)
	require.NoError(t, r.Close())
	require.Equal(t, int32(1), echo.runs.Load(), "the approved call must run exactly once")
}

// Deny flavor of the window-eviction case: the verdict must resolve the held
// copy (nothing runs, nothing posts) and clear the hold so the end_turn idle
// in the same history can stop the runner — instead of the occupied gate
// deferring idle-out forever.
func TestSessionToolRunner_ReconcileDenyForCallMissingFromHistoryLetsRunnerStop(t *testing.T) {
	server := newSessionEventsServer(t)

	var streamConns atomic.Int32
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		if streamConns.Add(1) == 1 {
			streamWriter(w, r, []string{
				sseLine("agent.tool_use", askToolUseEvt("tu", "echo", nil, "ask")),
			}, false)
			return
		}
		streamWriter(w, r, nil, true)
	}

	var lists atomic.Int32
	server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if lists.Add(1) >= 2 {
			body, _ := json.Marshal(map[string]any{
				"data":     []any{toolConfirmationEvt("tu", "deny"), idleEndTurnEvt("evt_idle")},
				"has_more": false, "first_id": "evt_conf_tu", "last_id": "evt_idle",
			})
			_, _ = w.Write(body)
			return
		}
		_, _ = w.Write([]byte(emptyEventList()))
	}

	server.HandleSend = func(http.ResponseWriter, *http.Request) { t.Error("a denied call posted a result") }

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 100*time.Millisecond)

	var yielded []DispatchedToolCall
	for r.Next() {
		yielded = append(yielded, r.Current())
	}
	require.ErrorIs(t, r.Err(), ErrIdleTimeout, "must stop on its own once the deny resolves the held call")
	require.Len(t, yielded, 1, "the denied held call is yielded exactly once")
	require.Equal(t, "tu", yielded[0].ToolUseID)
	require.Equal(t, "deny", yielded[0].Confirmation)
	require.Equal(t, int32(0), echo.runs.Load(), "the denied call must never run")
}

// ===== the idle countdown vs. confirmation-gated calls =====
//
// The server up-converts a legacy v1beta idle to stop_reason "end_turn"
// unconditionally, so a client CAN see an end_turn while a call is held for
// confirmation. The countdown must never run over gated work: stopping then
// drops the held call when its verdict later arrives, or cuts the runner off
// before a released call's result can drive the next turn.

// A call held awaiting its user.tool_confirmation defers the countdown even
// though the live stream reported an end_turn idle right behind it.
func TestSessionToolRunner_HeldCallBlocksIdleTimeout(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		// An ask-gated call, then an end_turn idle, then silence — no verdict.
		streamWriter(w, r, []string{
			sseLine("agent.tool_use", askToolUseEvt("evt_ask", "echo", map[string]any{}, "ask")),
			sseLine("session.status_idle", idleEndTurnEvt("evt_idle")),
		}, true)
	}
	server.HandleSend = func(http.ResponseWriter, *http.Request) {
		t.Error("nothing may be posted while the call is held")
	}

	echo := &stubBetaTool{name: "echo"}
	// Short MaxIdle: a wrongly-armed countdown would end the runner ~150ms in;
	// a correct runner instead blocks until this ctx expires ~2s in.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 150*time.Millisecond)

	start := time.Now()
	for r.Next() {
		t.Fatalf("unexpected yield while holding a call: %+v", r.Current())
	}
	elapsed := time.Since(start)

	require.NotErrorIs(t, r.Err(), ErrIdleTimeout,
		"runner idle-timed-out while a call was held awaiting confirmation")
	require.GreaterOrEqual(t, elapsed, time.Second,
		"runner ended after only %s — it idle-timed-out instead of holding the call", elapsed)
	require.Equal(t, int32(0), echo.runs.Load())
	require.NoError(t, r.Close())
}

// Reconcile flavor: the held call and its end_turn both come from history, so
// the arm the reconcile pass owes must be held pending until the verdict lands.
func TestSessionToolRunner_OpenApprovalKeepsRunnerAliveThenAnswerLetsItStop(t *testing.T) {
	// The user's answer closes the approval — a deny, or an unrecognized
	// verdict resolved as one (fail closed without wedging the runner open).
	for _, verdict := range []string{"deny", "not_a_verdict"} {
		t.Run(verdict, func(t *testing.T) {
			server := newSessionEventsServer(t)
			const wait = 600 * time.Millisecond
			// History shows a held call and the turn ending.
			server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
				body, _ := json.Marshal(map[string]any{
					"data":     []any{askToolUseEvt("tu", "echo", nil, "ask"), idleEndTurnEvt("evt_idle")},
					"has_more": false, "first_id": "tu", "last_id": "evt_idle",
				})
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(body)
			}
			// The live stream stays silent while the answer is open, then
			// delivers the user's verdict.
			server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
				f := w.(http.Flusher)
				w.Header().Set("Content-Type", "text/event-stream")
				w.WriteHeader(http.StatusOK)
				f.Flush()
				time.Sleep(wait)
				_, _ = io.WriteString(w, sseLine("user.tool_confirmation", toolConfirmationEvt("tu", verdict)))
				f.Flush()
				<-r.Context().Done()
			}
			server.HandleSend = func(http.ResponseWriter, *http.Request) { t.Error("a denied call posted a result") }

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			r := newShortIdleRunner(t, ctx, server.Client(), nil, 100*time.Millisecond)

			start := time.Now()
			var yielded []DispatchedToolCall
			for r.Next() {
				yielded = append(yielded, r.Current())
			}
			require.ErrorIs(t, r.Err(), ErrIdleTimeout, "must stop on its own via MaxIdle once the answer closes the approval")
			require.GreaterOrEqual(t, time.Since(start), wait, "stopped while an approval was still open")
			require.Len(t, yielded, 1, "the denied held call is yielded exactly once")
			require.Equal(t, "deny", yielded[0].Confirmation)
		})
	}
}

// The mirror image of the two tests above: once the deny retires the last
// blocker, the end_turn the runner saw mid-hold applies and it stops on its
// own — a held call defers the countdown, it must not cancel it outright.
func TestSessionToolRunner_DenyAfterEndTurnResumesIdle(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		// Hold the call, go idle on end_turn, then deny it — and stay connected.
		streamWriter(w, r, []string{
			sseLine("agent.tool_use", askToolUseEvt("evt_ask", "echo", map[string]any{}, "ask")),
			sseLine("session.status_idle", idleEndTurnEvt("evt_idle")),
			sseLine("user.tool_confirmation", toolConfirmationEvt("evt_ask", "deny")),
		}, true)
	}
	server.HandleSend = func(http.ResponseWriter, *http.Request) { t.Error("a denied tool must post nothing") }

	echo := &stubBetaTool{name: "echo"}
	// MaxIdle well under the ctx bound: a correct runner times out shortly after
	// the deny; a runner that mis-counts the released call as outstanding hangs
	// until ctx.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 200*time.Millisecond)

	start := time.Now()
	var yielded []DispatchedToolCall
	for r.Next() {
		yielded = append(yielded, r.Current())
	}
	elapsed := time.Since(start)

	require.ErrorIs(t, r.Err(), ErrIdleTimeout,
		"the idle countdown must resume after a deny releases the last held call")
	require.Less(t, elapsed, 4*time.Second,
		"runner took %s — the released call was wrongly counted as outstanding and the countdown never resumed", elapsed)
	require.Len(t, yielded, 1, "the denied call is yielded exactly once")
	require.Equal(t, "deny", yielded[0].Confirmation)
	require.Equal(t, int32(0), echo.runs.Load())
	require.NoError(t, r.Close())
}

// An end_turn armed the countdown, then the stream dropped. The reconciled
// history ends with an ask-gated call (held) and its end_turn, so the arm
// defers — but the pre-disconnect countdown is now stale evidence and must be
// cancelled, or the runner stops MaxIdle after the *old* end_turn while the
// confirmation is still pending on a human.
func TestSessionToolRunner_DeferredArmCancelsStalePreDisconnectCountdown(t *testing.T) {
	server := newSessionEventsServer(t)
	const maxIdle = 800 * time.Millisecond
	const wait = 2 * time.Second // > stale-stamp expiry; the deny arrives only after it

	var streamConns atomic.Int32
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		if streamConns.Add(1) == 1 {
			// The agent finishes a turn (end_turn arms the countdown)… then
			// the now-quiet SSE connection is dropped — exactly what load
			// balancers do to idle streams. The runner reconnects.
			streamWriter(w, r, []string{
				sseLine("session.status_idle", idleEndTurnEvt("evt_idle_1")),
			}, false)
			return
		}
		// Second connection: the approver has stepped away — silence past
		// the stale stamp's expiry — then they come back and deny, letting
		// the runner stop on a fresh window.
		f := w.(http.Flusher)
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		f.Flush()
		time.Sleep(wait)
		_, _ = io.WriteString(w, sseLine("user.tool_confirmation", toolConfirmationEvt("tu", "deny")))
		f.Flush()
		<-r.Context().Done()
	}

	var lists atomic.Int32
	server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// While the runner was reconnecting, the user sent a follow-up from
		// another client (the product UI); the agent hit a tool the policy
		// gates on approval, and its turn ended parked on that ask. The
		// reconnect's reconcile is how the runner learns all of this.
		if lists.Add(1) >= 2 {
			body, _ := json.Marshal(map[string]any{
				"data":     []any{askToolUseEvt("tu", "echo", nil, "ask"), idleEndTurnEvt("evt_idle_2")},
				"has_more": false, "first_id": "tu", "last_id": "evt_idle_2",
			})
			_, _ = w.Write(body)
			return
		}
		_, _ = w.Write([]byte(emptyEventList()))
	}
	server.HandleSend = func(http.ResponseWriter, *http.Request) { t.Error("a held call posted a result") }

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), nil, maxIdle)

	start := time.Now()
	var yielded []DispatchedToolCall
	for r.Next() {
		yielded = append(yielded, r.Current())
	}
	require.ErrorIs(t, r.Err(), ErrIdleTimeout, "must stop on its own once the deny resolves the held call")
	require.GreaterOrEqual(t, time.Since(start), wait,
		"stopped off the stale pre-disconnect stamp while a confirmation was pending")
	require.Len(t, yielded, 1, "the denied held call is yielded exactly once")
	require.Equal(t, "deny", yielded[0].Confirmation)
}

// The turn ended while the call was held, so an idle countdown is owed — but
// it must not start until the approved call has fully dispatched, or a tool
// slower than MaxIdle has the runner stop mid-flight and orphan the turn its
// result starts.
func TestSessionToolRunner_ApprovedCallStillExecutingDefersIdleStop(t *testing.T) {
	server := newSessionEventsServer(t)
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		streamWriter(w, r, []string{
			sseLine("agent.tool_use", askToolUseEvt("tu", "slow", nil, "ask")),
			sseLine("session.status_idle", idleEndTurnEvt("evt_idle")),
			sseLine("user.tool_confirmation", toolConfirmationEvt("tu", "allow")),
		}, true)
	}
	var postedAt atomic.Int64
	server.HandleSend = func(w http.ResponseWriter, _ *http.Request) {
		postedAt.Store(time.Now().UnixNano())
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sendOK()))
	}

	const maxIdle = 400 * time.Millisecond
	slow := &stubBetaTool{name: "slow", run: func(_ context.Context, _ json.RawMessage) (string, bool) {
		time.Sleep(600 * time.Millisecond)
		return "ran", false
	}}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{slow}, maxIdle)

	require.True(t, r.Next(), "the approved call must complete and be yielded: %v", r.Err())
	require.True(t, r.Current().Posted)
	require.Equal(t, "allow", r.Current().Confirmation)
	for r.Next() {
		t.Fatalf("yielded: %+v", r.Current())
	}
	stoppedAt := time.Now()
	require.ErrorIs(t, r.Err(), ErrIdleTimeout)
	// A runner that armed the countdown at the verdict has it expire during
	// the 600ms execute and stops as soon as the call is done; deferring
	// grants a full fresh window after the dispatch instead.
	require.Greater(t, stoppedAt.Sub(time.Unix(0, postedAt.Load())), maxIdle*6/10,
		"idle countdown ran over the executing approved call")
}

// A disarm racing the pending-arm apply must win: before the apply it clears
// the pending arm, after it clears the stamp the apply wrote — a cancelled
// countdown must never resurrect.
func TestIdleClock_DisarmWinsAgainstPendingArmApply(t *testing.T) {
	for i := range 10000 {
		c := newIdleClock(time.Hour)
		// A call is held for approval when the turn ends: the countdown is
		// owed but deferred (pending, unstamped).
		c.block("tu")
		c.arm()
		startGate := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(2)
		// The approved call finishes dispatching — the dispatch goroutine
		// retires the blocker, which applies the owed arm…
		go func() { defer wg.Done(); <-startGate; c.unblock("tu") }()
		// …while the stream goroutine processes a new event (e.g. the echoed
		// user.tool_result of that very call): the session is not idle, so
		// it disarms.
		go func() { defer wg.Done(); <-startGate; c.disarm() }()
		close(startGate)
		wg.Wait()
		c.mu.Lock()
		armed := !c.armedAt.IsZero() || c.armPending
		c.mu.Unlock()
		require.False(t, armed,
			"iteration %d: countdown left armed after a disarm — the pending-arm apply raced past it", i)
	}
}
