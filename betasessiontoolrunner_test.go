package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
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

func idleEndTurnEvt(id string) map[string]any {
	return map[string]any{
		"type":         "session.status_idle",
		"id":           id,
		"stop_reason":  map[string]any{"type": "end_turn"},
		"processed_at": "2026-05-11T12:00:00Z",
	}
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

// TestSessionToolRunner_ReconcileRetriesFailedPost covers the regression where
// a tool_use whose result post failed was marked answered anyway and so never
// retried. The post must only count as answered once it actually succeeds, so
// the reconcile after a reconnect re-dispatches the still-unanswered call.
func TestSessionToolRunner_ReconcileRetriesFailedPost(t *testing.T) {
	server := newSessionEventsServer(t)

	var streamConns atomic.Int32
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		if streamConns.Add(1) == 1 {
			// First connection: deliver the tool_use, then close so the runner
			// reconnects and reconciles.
			streamWriter(w, r, []string{
				sseLine("agent.tool_use", toolUseEvt("evt_1", "echo", map[string]any{})),
			}, false)
			return
		}
		streamWriter(w, r, nil, true) // later connections: hold open, no events
	}

	var lists atomic.Int32
	server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// The reconcile after the reconnect must still see the unanswered
		// tool_use in history so it can re-dispatch it.
		if lists.Add(1) >= 2 {
			body, _ := json.Marshal(map[string]any{
				"data":     []any{toolUseEvt("evt_1", "echo", map[string]any{})},
				"first_id": "evt_1", "has_more": false, "last_id": "evt_1",
			})
			_, _ = w.Write(body)
			return
		}
		_, _ = w.Write([]byte(emptyEventList()))
	}

	var sends atomic.Int32
	server.HandleSend = func(w http.ResponseWriter, _ *http.Request) {
		if sends.Add(1) == 1 {
			// First post fails permanently — the call must NOT be marked
			// answered, so the next reconcile retries it.
			http.Error(w, "bad", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(sendOK()))
	}

	echo := &stubBetaTool{name: "echo"}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), []BetaTool{echo}, 0)

	require.True(t, r.Next(), "first dispatch should be yielded")
	first := r.Current()
	require.Equal(t, "evt_1", first.ToolUseID)
	require.False(t, first.Posted, "the first result post fails (permanent 4xx)")

	require.True(t, r.Next(), "reconcile after reconnect must re-dispatch the unanswered tool_use")
	second := r.Current()
	require.Equal(t, "evt_1", second.ToolUseID)
	require.True(t, second.Posted, "the retried result post succeeds")

	require.GreaterOrEqual(t, echo.runs.Load(), int32(2), "the tool is re-executed on the retry")
	require.NoError(t, r.Close())
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

// TestSessionToolRunner_SkippedUnownedToolDoesNotTripIdle pins that a skipped
// (unanswered) unowned tool_use stays out of the end-turn accounting:
// reconcile sees history ending on an end_turn idle but with the unowned
// tool_use still unanswered, so it must NOT arm the idle countdown — the
// runner has not actually handled that call, its owner still has to.
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

	// Bound the run: with a short MaxIdle, a wrongly-armed idle countdown
	// would end the runner ~150ms in; a correct runner instead blocks until
	// this ctx expires ~2s in. The gap makes the two outcomes unambiguous.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	r := newShortIdleRunner(t, ctx, server.Client(), nil, 150*time.Millisecond)

	require.True(t, r.Next(), "reconcile must still surface the unowned call")
	call := r.Current()
	require.Equal(t, "evt_pending", call.ToolUseID)
	require.False(t, call.Posted)
	require.False(t, call.IsError)
	require.Empty(t, call.Result.ToolUseID, "no result was built for the skipped call")

	// The skipped tool_use is unanswered, so reconcile must keep it OUT of
	// the end-turn accounting even though history ends on an end_turn idle.
	// Single-goroutine: drain to termination and assert it was ctx expiry,
	// not an idle timeout.
	start := time.Now()
	for r.Next() {
		t.Fatalf("unexpected extra yield: %+v", r.Current())
	}
	elapsed := time.Since(start)

	require.NotErrorIs(t, r.Err(), ErrIdleTimeout,
		"runner falsely went idle on a skipped unowned tool_use")
	require.GreaterOrEqual(t, elapsed, time.Second,
		"runner ended after only %s — it idle-timed-out instead of staying up for the pending unowned call", elapsed)
	require.NoError(t, r.Close())
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
