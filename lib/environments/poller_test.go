package environments

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Unit tests for the small helpers in poller.go. These don't need a server.
// ---------------------------------------------------------------------------

func TestBackoff(t *testing.T) {
	tests := []struct {
		description string
		failures    int
		want        time.Duration
	}{
		{"first failure backs off 2s (2^1) so a transient blip recovers quickly", 1, 2 * time.Second},
		{"second consecutive failure doubles to 4s following the exponential curve", 2, 4 * time.Second},
		{"fifth failure reaches 32s while still under the cap", 5, 32 * time.Second},
		{"sixth failure would compute 64s but is clamped to the 60s cap", 6, 60 * time.Second},
		{"very large failure counts stay clamped at the cap rather than overflowing", 30, 60 * time.Second},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			require.Equal(t, tc.want, backoff(tc.failures))
		})
	}
}

func TestJitter(t *testing.T) {
	low, high := 1*time.Second, 3*time.Second
	for range 100 {
		got := jitter(low, high)
		require.GreaterOrEqual(t, got, low, "jitter must never undershoot lower bound")
		require.Less(t, got, high, "jitter must stay strictly below upper bound")
	}
}

func TestIsStatus(t *testing.T) {
	tests := []struct {
		description string
		err         error
		code        int
		want        bool
	}{
		{"matching status code on a wrapped api error is detected so 409 on Stop can be suppressed",
			&anthropic.Error{StatusCode: 409}, 409, true},
		{"different status code returns false so only the intended code is matched",
			&anthropic.Error{StatusCode: 500}, 409, false},
		{"non-api error type is never matched even if the message looks similar",
			errors.New("409 conflict"), 409, false},
		{"nil error is treated as no-match rather than panicking",
			nil, 409, false},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			require.Equal(t, tc.want, isStatus(tc.err, tc.code))
		})
	}
}

func TestIsFatal4xx(t *testing.T) {
	tests := []struct {
		description string
		err         error
		want        bool
	}{
		{"400 is fatal because the request body cannot succeed on retry",
			&anthropic.Error{StatusCode: 400}, true},
		{"408 is excluded because timeouts deserve backoff, not teardown",
			&anthropic.Error{StatusCode: 408}, false},
		{"429 is excluded because rate-limits deserve backoff, not teardown",
			&anthropic.Error{StatusCode: 429}, false},
		{"500 is excluded because server-side errors retry, not abort",
			&anthropic.Error{StatusCode: 500}, false},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			require.Equal(t, tc.want, isFatal4xx(tc.err))
		})
	}
}

// ---------------------------------------------------------------------------
// Integration tests against an httptest server scripting poll/ack/stop.
// ---------------------------------------------------------------------------

type recordedCall struct {
	method string
	path   string
	auth   string
	// apiKey captures the X-Api-Key header value on the recorded request.
	// Runner helpers must NOT leak the parent client's default API key
	// alongside their bearer credential; tests assert this is empty on
	// every helper-issued call.
	apiKey string
	body   string
	header http.Header
}

// fakeWorkServer is a minimal httptest server that records every request
// and dispatches to per-endpoint handler functions provided by the test.
// Tests assign Handle* fields before constructing the helper under test;
// calls to any unassigned endpoint fail the test. Shared by both
// poller_test.go and dispatcher_test.go.
type fakeWorkServer struct {
	t      *testing.T
	server *httptest.Server

	mu    sync.Mutex
	calls []recordedCall

	// Environment work endpoints (poller).
	HandlePoll http.HandlerFunc
	HandleAck  http.HandlerFunc
	HandleStop http.HandlerFunc

	// Environment work endpoints (dispatcher).
	HandleHeartbeat http.HandlerFunc

	// Session event endpoints (dispatcher).
	HandleStream http.HandlerFunc // GET  /v1/sessions/{id}/events/stream
	HandleSend   http.HandlerFunc // POST /v1/sessions/{id}/events
	HandleList   http.HandlerFunc // GET  /v1/sessions/{id}/events

	// Session lookup endpoint (skill setup).
	HandleSessionGet http.HandlerFunc // GET /v1/sessions/{id}
}

func newFakeWorkServer(t *testing.T) *fakeWorkServer {
	f := &fakeWorkServer{t: t}
	f.server = httptest.NewServer(http.HandlerFunc(f.serveHTTP))
	t.Cleanup(func() { f.server.Close() })
	return f
}

func (f *fakeWorkServer) serveHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	// Rewind so the handler can re-read the body. Httptest handlers commonly
	// inspect the payload to assert what the SDK sent.
	r.Body = io.NopCloser(bytes.NewReader(body))
	f.mu.Lock()
	f.calls = append(f.calls, recordedCall{
		method: r.Method,
		path:   r.URL.Path,
		auth:   r.Header.Get("Authorization"),
		apiKey: r.Header.Get("X-Api-Key"),
		body:   string(body),
		header: r.Header.Clone(),
	})
	f.mu.Unlock()

	switch {
	case strings.HasSuffix(r.URL.Path, "/work/poll"):
		require.NotNil(f.t, f.HandlePoll, "unscripted Poll call: %s", r.URL.Path)
		f.HandlePoll(w, r)
	case strings.HasSuffix(r.URL.Path, "/ack"):
		require.NotNil(f.t, f.HandleAck, "unscripted Ack call: %s", r.URL.Path)
		f.HandleAck(w, r)
	case strings.HasSuffix(r.URL.Path, "/stop"):
		require.NotNil(f.t, f.HandleStop, "unscripted Stop call: %s", r.URL.Path)
		f.HandleStop(w, r)
	case strings.HasSuffix(r.URL.Path, "/heartbeat"):
		require.NotNil(f.t, f.HandleHeartbeat, "unscripted Heartbeat call: %s", r.URL.Path)
		f.HandleHeartbeat(w, r)
	case strings.HasSuffix(r.URL.Path, "/events/stream"):
		require.NotNil(f.t, f.HandleStream, "unscripted StreamEvents call: %s", r.URL.Path)
		f.HandleStream(w, r)
	case strings.HasSuffix(r.URL.Path, "/events"):
		switch r.Method {
		case http.MethodPost:
			require.NotNil(f.t, f.HandleSend, "unscripted events.Send call: %s", r.URL.Path)
			f.HandleSend(w, r)
		case http.MethodGet:
			require.NotNil(f.t, f.HandleList, "unscripted events.List call: %s", r.URL.Path)
			f.HandleList(w, r)
		default:
			f.t.Errorf("unexpected method on events endpoint: %s", r.Method)
			http.Error(w, "unexpected", http.StatusMethodNotAllowed)
		}
	case strings.Contains(r.URL.Path, "/sessions/"):
		// Bare GET /v1/sessions/{id} — the SetupSkills session lookup. The
		// /events{,/stream} cases above already handled the event endpoints.
		require.NotNil(f.t, f.HandleSessionGet, "unscripted Sessions.Get call: %s", r.URL.Path)
		f.HandleSessionGet(w, r)
	default:
		f.t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		http.Error(w, "unexpected", http.StatusNotImplemented)
	}
}

func (f *fakeWorkServer) Calls() []recordedCall {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.calls)
}

func (f *fakeWorkServer) Client() anthropic.Client {
	return anthropic.NewClient(
		option.WithBaseURL(f.server.URL),
		option.WithAPIKey("client-default-key"),
		option.WithMaxRetries(0),
	)
}

// workJSON returns a minimally-valid BetaSelfHostedWork JSON body. All
// api:"required" fields are present so the SDK's strict unmarshal accepts
// it. dataType is "session" or "health_check". The "secret" field is still
// emitted (codegen marks it required) but the helper no longer reads it.
func workJSON(workID, envID, dataType string) string {
	ts := "2026-05-11T12:00:00Z"
	data := map[string]any{"type": dataType, "id": "sesn_test"}
	if dataType == "session" {
		data["mode"] = "private"
	}
	body, _ := json.Marshal(map[string]any{
		"id":                  workID,
		"acknowledged_at":     ts,
		"actor":               map[string]any{"type": "api_key", "api_key_id": "apk_test"},
		"created_at":          ts,
		"data":                data,
		"environment_id":      envID,
		"latest_heartbeat_at": ts,
		"metadata":            map[string]string{},
		"secret":              "unused-by-helper",
		"started_at":          ts,
		"state":               "queued",
		"stop_requested_at":   "",
		"stopped_at":          "",
		"type":                "work",
	})
	return string(body)
}

var silentLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestWorkPoller_YieldsAndPostsStopOnClose(t *testing.T) {
	server := newFakeWorkServer(t)
	work := workJSON("work_1", "env_1", "session")

	pollCount := 0
	server.HandlePoll = func(w http.ResponseWriter, _ *http.Request) {
		pollCount++
		w.Header().Set("Content-Type", "application/json")
		if pollCount == 1 {
			_, _ = w.Write([]byte(work))
		} else {
			// Empty poll — null body — return 204 with nothing.
			w.WriteHeader(http.StatusNoContent)
		}
	}
	server.HandleAck = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(work))
	}
	server.HandleStop = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(work))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := NewWorkPoller(ctx, server.Client(), WorkPollerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "env_key",
		WorkerID:       "test-worker",
		Logger:         silentLogger,
	})

	require.True(t, p.Next(), "first Next() must yield the scripted work")
	require.Equal(t, "work_1", p.Current().ID)

	require.NoError(t, p.Close())

	calls := server.Calls()
	require.Len(t, calls, 3, "expected poll + ack + stop")
	require.Contains(t, calls[0].path, "/work/poll")
	require.Contains(t, calls[1].path, "/ack")
	require.Contains(t, calls[2].path, "/stop")

	require.Equal(t, "Bearer env_key", calls[0].auth,
		"Poll must use the environment key")
	require.Equal(t, "Bearer env_key", calls[1].auth,
		"Ack must use the environment key")
	require.Equal(t, "Bearer env_key", calls[2].auth,
		"Stop must use the environment key")

	// option.WithAuthToken alone only *adds* an Authorization header; the
	// parent client's WithAPIKey middleware still sets X-Api-Key. The
	// poller's helperReqOpts must explicitly delete X-Api-Key per-request so
	// both headers don't ride on the wire together.
	require.Empty(t, calls[0].apiKey,
		"Poll must not leak the parent client's X-Api-Key")
	require.Empty(t, calls[1].apiKey,
		"Ack must not leak the parent client's X-Api-Key")
	require.Empty(t, calls[2].apiKey,
		"Stop must not leak the parent client's X-Api-Key")
}

func TestWorkPoller_StopRunsBeforeNextPoll(t *testing.T) {
	server := newFakeWorkServer(t)
	work1 := workJSON("work_1", "env_1", "session")
	work2 := workJSON("work_2", "env_1", "session")

	pollCount := 0
	server.HandlePoll = func(w http.ResponseWriter, _ *http.Request) {
		pollCount++
		w.Header().Set("Content-Type", "application/json")
		switch pollCount {
		case 1:
			_, _ = w.Write([]byte(work1))
		case 2:
			_, _ = w.Write([]byte(work2))
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	}
	server.HandleAck = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(work1))
	}
	server.HandleStop = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(work1))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := NewWorkPoller(ctx, server.Client(), WorkPollerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "env_key",
		Logger:         silentLogger,
	})
	t.Cleanup(func() { _ = p.Close() })

	require.True(t, p.Next())
	require.Equal(t, "work_1", p.Current().ID)

	require.True(t, p.Next())
	require.Equal(t, "work_2", p.Current().ID)

	calls := server.Calls()
	require.GreaterOrEqual(t, len(calls), 5)

	// Expected order: poll → ack(1) → stop(1) → poll → ack(2). Stop for
	// work_1 must appear BEFORE the second poll, even though work_1 is
	// already a different item from the second Next().
	stop1Idx := -1
	poll2Idx := -1
	for i, c := range calls {
		if stop1Idx == -1 && strings.Contains(c.path, "work_1/stop") {
			stop1Idx = i
		}
		if stop1Idx != -1 && poll2Idx == -1 && strings.Contains(c.path, "/work/poll") {
			poll2Idx = i
		}
	}
	require.NotEqual(t, -1, stop1Idx, "stop for work_1 was never posted")
	require.NotEqual(t, -1, poll2Idx, "second poll was never posted")
	require.Less(t, stop1Idx, poll2Idx,
		"deferred stop for work_1 must run before the second poll")
}

func TestWorkPoller_AckFailureSkips(t *testing.T) {
	server := newFakeWorkServer(t)
	first := workJSON("work_1", "env_1", "session")
	second := workJSON("work_2", "env_1", "session")

	pollCount := 0
	server.HandlePoll = func(w http.ResponseWriter, _ *http.Request) {
		pollCount++
		w.Header().Set("Content-Type", "application/json")
		switch pollCount {
		case 1:
			_, _ = w.Write([]byte(first))
		case 2:
			_, _ = w.Write([]byte(second))
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	}
	ackCount := 0
	server.HandleAck = func(w http.ResponseWriter, r *http.Request) {
		ackCount++
		if strings.Contains(r.URL.Path, "work_1") {
			http.Error(w, `{"type":"error","error":{"type":"server_error","message":"sim ack failure"}}`,
				http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(second))
	}
	server.HandleStop = func(w http.ResponseWriter, _ *http.Request) {
		// Stop is now posted both for the item whose ack failed (a best-effort
		// force-stop so it doesn't dangle) and, on Close, for the item that
		// successfully yielded.
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(second))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := NewWorkPoller(ctx, server.Client(), WorkPollerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "env_key",
		Logger:         silentLogger,
	})
	t.Cleanup(func() { _ = p.Close() })

	require.True(t, p.Next())
	require.Equal(t, "work_2", p.Current().ID,
		"work whose ack 5xx'd should be skipped; iterator yields the next item")
	require.Equal(t, 2, ackCount, "both acks were attempted; only the second succeeded")

	// The unprocessable item (ack 5xx'd) must be force-stopped rather than left
	// dangling, before the poller moves on to the next item.
	var work1Stop *recordedCall
	for _, c := range server.Calls() {
		if strings.Contains(c.path, "work_1/stop") {
			c := c
			work1Stop = &c
			break
		}
	}
	require.NotNil(t, work1Stop, "ack-failed work_1 must be force-stopped, not silently skipped")
	require.Contains(t, work1Stop.body, `"force":true`, "discard stop must force")
}

func TestWorkPoller_CtxCancelDuringPollExitsCleanly(t *testing.T) {
	server := newFakeWorkServer(t)
	pollUnblocked := make(chan struct{})
	server.HandlePoll = func(w http.ResponseWriter, r *http.Request) {
		<-pollUnblocked
		select {
		case <-r.Context().Done():
			// Server saw the cancellation — return 499-ish; client side
			// will surface ctx.Err.
			http.Error(w, "cancelled", 499)
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	}
	server.HandleAck = func(w http.ResponseWriter, _ *http.Request) {
		t.Fatal("ack must not be called when ctx is cancelled before any work is claimed")
	}
	server.HandleStop = func(w http.ResponseWriter, _ *http.Request) {
		t.Fatal("stop must not be called when nothing was claimed")
	}

	ctx, cancel := context.WithCancel(context.Background())

	p := NewWorkPoller(ctx, server.Client(), WorkPollerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "env_key",
		Logger:         silentLogger,
	})

	done := make(chan bool, 1)
	go func() {
		done <- p.Next()
	}()

	// Give the poll request time to be in flight, then cancel.
	time.Sleep(50 * time.Millisecond)
	cancel()
	close(pollUnblocked)

	select {
	case got := <-done:
		require.False(t, got, "Next() must return false after ctx cancellation")
		require.NoError(t, p.Err(), "ctx cancellation is normal termination, not an error")
	case <-time.After(2 * time.Second):
		t.Fatal("Next() did not return after ctx cancellation")
	}

	require.NoError(t, p.Close())
}

func TestWorkPoller_CloseIsIdempotent(t *testing.T) {
	server := newFakeWorkServer(t)
	server.HandlePoll = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
	server.HandleAck = func(w http.ResponseWriter, _ *http.Request) {
		t.Fatal("ack should not be called")
	}
	server.HandleStop = func(w http.ResponseWriter, _ *http.Request) {
		t.Fatal("stop should not be called when nothing was claimed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	p := NewWorkPoller(ctx, server.Client(), WorkPollerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "env_key",
		Logger:         silentLogger,
	})

	require.NoError(t, p.Close())
	require.NoError(t, p.Close(), "Close must be safe to call multiple times")

	require.False(t, p.Next(), "Next after Close must return false immediately")
}

func TestWorkPoller_All_RangeOverFunc(t *testing.T) {
	server := newFakeWorkServer(t)
	work := workJSON("work_1", "env_1", "session")

	pollCount := 0
	server.HandlePoll = func(w http.ResponseWriter, _ *http.Request) {
		pollCount++
		w.Header().Set("Content-Type", "application/json")
		if pollCount == 1 {
			_, _ = w.Write([]byte(work))
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
	server.HandleAck = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(work))
	}
	server.HandleStop = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(work))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := NewWorkPoller(ctx, server.Client(), WorkPollerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "env_key",
		Logger:         silentLogger,
	})
	t.Cleanup(func() { _ = p.Close() })

	var seen []string
	for item, err := range p.All() {
		require.NoError(t, err)
		seen = append(seen, item.ID)
		if len(seen) == 1 {
			break
		}
	}
	require.Equal(t, []string{"work_1"}, seen)
}

// writeEmptyPoll writes the response the Go SDK decodes as an empty queue:
// 200 with a JSON `null` body, leaving the decoded *BetaSelfHostedWork nil.
// (A bare 204 has no Content-Type and trips the strict decoder instead.)
func writeEmptyPoll(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte("null"))
}

// TestWorkPoller_DrainReturnsOnEmptyQueue asserts that with Drain set the
// poller ends iteration as soon as a poll comes back empty, instead of
// sleeping and re-polling forever. Err must stay nil (normal termination).
func TestWorkPoller_DrainReturnsOnEmptyQueue(t *testing.T) {
	server := newFakeWorkServer(t)

	pollCount := 0
	server.HandlePoll = func(w http.ResponseWriter, _ *http.Request) {
		pollCount++
		writeEmptyPoll(w)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := NewWorkPoller(ctx, server.Client(), WorkPollerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "env_key",
		Drain:          true,
		Logger:         silentLogger,
	})
	t.Cleanup(func() { _ = p.Close() })

	start := time.Now()
	require.False(t, p.Next(), "Drain must end iteration on the first empty poll")
	require.NoError(t, p.Err(), "draining the queue is normal termination, not an error")
	require.Less(t, time.Since(start), 2*time.Second,
		"Drain must return promptly rather than sleeping between empty polls")
	require.Equal(t, 1, pollCount, "Drain must not re-poll after an empty queue")
}

// TestWorkPoller_DrainYieldsThenReturns asserts Drain still yields available
// work and only returns once the queue empties — and that the last item is
// still stopped on the way out (the poller's auto-stop model is unchanged).
func TestWorkPoller_DrainYieldsThenReturns(t *testing.T) {
	server := newFakeWorkServer(t)
	work := workJSON("work_1", "env_1", "session")

	pollCount := 0
	server.HandlePoll = func(w http.ResponseWriter, _ *http.Request) {
		pollCount++
		if pollCount == 1 {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(work))
			return
		}
		writeEmptyPoll(w)
	}
	server.HandleAck = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(work))
	}
	stopped := 0
	server.HandleStop = func(w http.ResponseWriter, _ *http.Request) {
		stopped++
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(work))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := NewWorkPoller(ctx, server.Client(), WorkPollerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "env_key",
		Drain:          true,
		Logger:         silentLogger,
	})
	t.Cleanup(func() { _ = p.Close() })

	require.True(t, p.Next(), "first Next must yield the scripted work")
	require.Equal(t, "work_1", p.Current().ID)
	require.False(t, p.Next(), "second Next must drain-return on the empty queue")
	require.NoError(t, p.Err())
	require.Equal(t, 1, stopped, "the yielded item must still be stopped on drain exit")
}

// TestWorkPoller_BlockMsWiring asserts how WorkPollerOptions.BlockMs and
// ReclaimOlderThanMs map onto the poll query: omitted -> default 999,
// param.Null -> omitted entirely (non-blocking), an explicit value passes
// through, and ReclaimOlderThanMs only appears when set.
func TestWorkPoller_BlockMsWiring(t *testing.T) {
	cases := []struct {
		name        string
		opts        WorkPollerOptions
		wantBlock   string // expected block_ms query value, "" means absent
		wantReclaim string // expected reclaim_older_than_ms, "" means absent
	}{
		{
			name:      "default omits to 999",
			opts:      WorkPollerOptions{},
			wantBlock: "999",
		},
		{
			name:      "explicit value passes through",
			opts:      WorkPollerOptions{BlockMs: param.NewOpt(int64(250))},
			wantBlock: "250",
		},
		{
			name:      "null omits block_ms for a non-blocking poll",
			opts:      WorkPollerOptions{BlockMs: param.Null[int64]()},
			wantBlock: "",
		},
		{
			name:        "reclaim_older_than_ms threads through when set",
			opts:        WorkPollerOptions{ReclaimOlderThanMs: param.NewOpt(int64(5000))},
			wantBlock:   "999",
			wantReclaim: "5000",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := newFakeWorkServer(t)
			var gotBlock, gotReclaim string
			var gotBlockSet, gotReclaimSet bool
			server.HandlePoll = func(w http.ResponseWriter, r *http.Request) {
				q := r.URL.Query()
				gotBlock, gotBlockSet = q.Get("block_ms"), q.Has("block_ms")
				gotReclaim, gotReclaimSet = q.Get("reclaim_older_than_ms"), q.Has("reclaim_older_than_ms")
				writeEmptyPoll(w)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			opts := tc.opts
			opts.EnvironmentID = "env_1"
			opts.EnvironmentKey = "env_key"
			opts.Drain = true // end after the single empty poll
			opts.Logger = silentLogger
			p := NewWorkPoller(ctx, server.Client(), opts)
			t.Cleanup(func() { _ = p.Close() })

			require.False(t, p.Next())
			require.NoError(t, p.Err())

			if tc.wantBlock == "" {
				require.False(t, gotBlockSet, "block_ms must be absent for a non-blocking poll")
			} else {
				require.True(t, gotBlockSet, "block_ms must be present")
				require.Equal(t, tc.wantBlock, gotBlock)
			}
			if tc.wantReclaim == "" {
				require.False(t, gotReclaimSet, "reclaim_older_than_ms must be absent unless set")
			} else {
				require.True(t, gotReclaimSet, "reclaim_older_than_ms must be present")
				require.Equal(t, tc.wantReclaim, gotReclaim)
			}
		})
	}
}
