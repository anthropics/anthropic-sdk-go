package environments

import (
	"context"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentWorker_CancelledCtxReturnsNil(t *testing.T) {
	server := newFakeWorkServer(t)
	worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "envkey",
		Logger:         silentLogger,
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	require.NoError(t, worker.Run(ctx), "a cancelled ctx is normal termination, not an error")
}

func TestEnvironmentWorker_DefaultWorkdirSnapshotsCwd(t *testing.T) {
	server := newFakeWorkServer(t)

	wd, err := os.Getwd()
	require.NoError(t, err)

	// An empty Workdir is snapshotted to the cwd at construction, not left as
	// a lazily-resolved ".". This mirrors the TS/Python helpers and keeps the
	// session workdir stable across a later os.Chdir.
	worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "envkey",
		Logger:         silentLogger,
	})
	require.Equal(t, wd, worker.opts.Workdir)
	require.NotEqual(t, ".", worker.opts.Workdir)

	// An explicit Workdir is honored as-is.
	explicit := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "envkey",
		Workdir:        "/some/explicit/dir",
		Logger:         silentLogger,
	})
	require.Equal(t, "/some/explicit/dir", explicit.opts.Workdir)
}

func TestEnvironmentWorker_RequiresEnvironmentID(t *testing.T) {
	server := newFakeWorkServer(t)
	worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		EnvironmentKey: "envkey",
		Logger:         silentLogger,
	})
	// The underlying WorkPoller surfaces the missing-id error on the first poll.
	require.Error(t, worker.Run(context.Background()))
}

// customHeaderName/customHeaderValue stand in for a caller-supplied
// proxy/routing header that must reach every request the self-hosted runner
// issues when threaded through EnvironmentWorkerOptions.RequestOptions.
const (
	customHeaderName  = "X-Custom-Proxy"
	customHeaderValue = "proxy-token-xyz"
)

// assertCustomHeaderEverywhere asserts that every recorded call carries the
// caller's custom header, authenticates with the environment key, AND does
// NOT leak the parent client's X-Api-Key. The third assertion is what closes
// the regression window the original #854 review missed: a future change
// that drops WithHeaderDel("X-Api-Key") from the worker side of the auth
// recipe (heartbeat, force-stop, runner stream/list/send, skill-setup
// session lookup) flips this assertion immediately, just as the poller-side
// assertion at poller_test.go:307-315 does for poll/ack/stop.
func assertCustomHeaderEverywhere(t *testing.T, calls []recordedCall) {
	t.Helper()
	require.NotEmpty(t, calls, "expected the worker to issue at least one request")
	for _, c := range calls {
		t.Run(c.method+" "+c.path, func(t *testing.T) {
			require.Equal(t, customHeaderValue, c.header.Get(customHeaderName),
				"caller-supplied custom header must reach %s %s", c.method, c.path)
			require.Equal(t, "Bearer env_key", c.auth,
				"environment-key auth must still win on %s %s", c.method, c.path)
			require.Empty(t, c.apiKey,
				"%s %s must not leak the parent client's X-Api-Key alongside the bearer", c.method, c.path)
		})
	}
}

// TestWorkPoller_ThreadsCustomRequestOptions covers the poller leg: a custom
// header supplied via WorkPollerOptions.RequestOptions must appear on
// Poll/Ack/Stop.
func TestWorkPoller_ThreadsCustomRequestOptions(t *testing.T) {
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
		w.WriteHeader(http.StatusNoContent)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := NewWorkPoller(ctx, server.Client(), WorkPollerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "env_key",
		WorkerID:       "test-worker",
		RequestOptions: []option.RequestOption{option.WithHeader(customHeaderName, customHeaderValue)},
		Logger:         silentLogger,
	})

	require.True(t, p.Next(), "first Next() must yield the scripted work")
	require.NoError(t, p.Close())

	calls := server.Calls()
	require.Len(t, calls, 3, "expected poll + ack + stop")
	assertCustomHeaderEverywhere(t, calls)
}

// TestEnvironmentWorker_ThreadsCustomRequestOptions covers the per-session
// leg: a custom header supplied via EnvironmentWorkerOptions.RequestOptions
// must reach the skill-setup session lookup, the lease heartbeat, the
// SessionToolRunner's stream/list, and the force-stop on exit.
func TestEnvironmentWorker_ThreadsCustomRequestOptions(t *testing.T) {
	server := newFakeWorkServer(t)

	// Gate the stream's terminate event on the heartbeat having fired, so the
	// heartbeat is deterministically observed before the session ends and the
	// per-session context is cancelled — no timing-based flakiness.
	heartbeatSeen := make(chan struct{})
	var heartbeatOnce sync.Once

	server.HandleSessionGet = func(w http.ResponseWriter, _ *http.Request) {
		// No skills -> SetupSkills does only the session lookup and returns.
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"agent":{"skills":[]}}`))
	}
	server.HandleHeartbeat = func(w http.ResponseWriter, _ *http.Request) {
		heartbeatOnce.Do(func() { close(heartbeatSeen) })
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"last_heartbeat":"2026-05-11T12:00:00Z","lease_extended":true,"state":"active","ttl_seconds":30,"type":"work_heartbeat"}`))
	}
	server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[],"first_id":null,"has_more":false,"last_id":null}`))
	}
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		require.True(t, ok, "stream response writer must support flushing")
		w.WriteHeader(http.StatusOK)
		flusher.Flush()
		// Wait for at least one heartbeat before terminating the session.
		select {
		case <-heartbeatSeen:
		case <-r.Context().Done():
			return
		case <-time.After(5 * time.Second):
			t.Error("heartbeat was never observed before stream terminate")
		}
		_, _ = w.Write([]byte("event: session.status_terminated\n" +
			`data: {"type":"session.status_terminated","id":"evt_term","processed_at":"2026-05-11T12:00:00Z"}` +
			"\n\n"))
		flusher.Flush()
	}
	server.HandleSend = func(w http.ResponseWriter, _ *http.Request) {
		// No tool_use fires, so Send should never be called.
		t.Error("Send must not be called when no tool_use event arrives")
		w.WriteHeader(http.StatusOK)
	}
	server.HandleStop = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "env_key",
		Workdir:        t.TempDir(),
		RequestOptions: []option.RequestOption{option.WithHeader(customHeaderName, customHeaderValue)},
		Logger:         silentLogger,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	require.NoError(t, worker.HandleItem(ctx, HandleItemOptions{
		WorkID:        "work_1",
		EnvironmentID: "env_1",
		SessionID:     "sesn_test",
	}), "a terminated session is benign and HandleItem should return nil")

	calls := server.Calls()
	assertCustomHeaderEverywhere(t, calls)

	// Sanity-check that the calls we expect to carry the header actually
	// happened — otherwise the assertion above could pass vacuously.
	want := map[string]bool{
		"sessions/sesn_test":               false, // skill-setup session lookup
		"work/work_1/heartbeat":            false, // lease heartbeat
		"sessions/sesn_test/events/stream": false, // runner stream
		"work/work_1/stop":                 false, // force-stop on exit
	}
	for _, c := range calls {
		for suffix := range want {
			if strings.HasSuffix(c.path, suffix) {
				want[suffix] = true
			}
		}
	}
	for suffix, seen := range want {
		require.True(t, seen, "expected a request to %q but none was recorded", suffix)
	}
}

// TestEnvironmentWorker_HeartbeatStartsBeforeSkillSetup pins the invariant
// that the lease heartbeat is running before SetupSkills's session
// lookup completes. The poller acks the item when it
// yields, so any gap between ack and the first heartbeat is a window where
// the control plane can reclaim the lease — and SetupSkills (a session
// lookup plus a per-skill download/extract) can take longer than the lease
// TTL on a slow network or a large bundle.
//
// The test enforces this by making HandleSessionGet block until at least
// one heartbeat has fired. If the heartbeat goroutine is started AFTER
// SetupSkills returns (the pre-fix shape), no heartbeat ever fires and
// HandleSessionGet blocks until the outer ctx times out — HandleItem then
// returns the wrapped ctx error rather than nil. With the fix, the
// heartbeat goroutine is started before SetupSkills, the channel closes
// while SetupSkills's session lookup is in flight, and the session
// completes normally.
func TestEnvironmentWorker_HeartbeatStartsBeforeSkillSetup(t *testing.T) {
	server := newFakeWorkServer(t)

	heartbeatSeen := make(chan struct{})
	var heartbeatOnce sync.Once

	server.HandleHeartbeat = func(w http.ResponseWriter, _ *http.Request) {
		heartbeatOnce.Do(func() { close(heartbeatSeen) })
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"last_heartbeat":"2026-05-11T12:00:00Z","lease_extended":true,"state":"active","ttl_seconds":30,"type":"work_heartbeat"}`))
	}
	server.HandleSessionGet = func(w http.ResponseWriter, r *http.Request) {
		// Block the session-get until a heartbeat has been observed. This is
		// what proves the heartbeat starts BEFORE skill setup finishes — if
		// it didn't, this handler would block forever (the heartbeat
		// goroutine never starts) and the outer context would time out.
		select {
		case <-heartbeatSeen:
		case <-r.Context().Done():
			return
		case <-time.After(5 * time.Second):
			t.Error("heartbeat did not fire while SetupSkills was in flight")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"agent":{"skills":[]}}`))
	}
	server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[],"first_id":null,"has_more":false,"last_id":null}`))
	}
	server.HandleStream = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		require.True(t, ok, "stream response writer must support flushing")
		w.WriteHeader(http.StatusOK)
		flusher.Flush()
		_, _ = w.Write([]byte("event: session.status_terminated\n" +
			`data: {"type":"session.status_terminated","id":"evt_term","processed_at":"2026-05-11T12:00:00Z"}` +
			"\n\n"))
		flusher.Flush()
	}
	server.HandleStop = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "env_key",
		Workdir:        t.TempDir(),
		Logger:         silentLogger,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	require.NoError(t, worker.HandleItem(ctx, HandleItemOptions{
		WorkID:        "work_1",
		EnvironmentID: "env_1",
		SessionID:     "sesn_test",
	}), "HandleItem must complete normally; if the heartbeat never starts, HandleSessionGet blocks until ctx timeout and HandleItem returns a wrapped ctx error")
}

// TestEnvironmentWorker_HeartbeatPerRequestTimeoutAndStalenessCeiling pins
// two related guarantees about runHeartbeat:
//
//  1. No per-request timeout — Heartbeat inherits the SDK's ~10-minute
//     default, so one slow/hung call outlives the entire lease window.
//  2. No staleness ceiling — a run of transient errors retries forever while
//     the lease silently expires server-side; the worker keeps executing
//     tools against a session another worker may also have claimed.
//
// The fixes interact: each Heartbeat is sent with WithRequestTimeout(interval)
// (so hung calls fail fast), and the transient-error branch bails when
// time.Since(lastSuccess) > ttl (so a run of fast-failing or fast-cancelled
// beats also terminates the heartbeat goroutine within a bounded window).
//
// This test scripts a server that returns one successful beat with a short
// TTL (2 s — shrinks interval to 1 s and ttl to 2 s) and then hangs every
// subsequent heartbeat. With the fixes, each hung beat is cancelled at the
// 1 s request timeout, and the staleness check fires once we're past the
// 2 s TTL since the last success — so the heartbeat goroutine cancels
// sessCtx, the runner exits with ctx.Canceled, and HandleItem returns the
// wrapped ctx error well within the 15 s outer ctx.
//
// Without the per-request-timeout fix, beat 2 would block on the hung
// handler for the full SDK default (10 min) — only one heartbeat call would
// land, and HandleItem would only return when the outer ctx times out.
// Without the staleness-ceiling fix (even with the request timeout) the
// goroutine would log transient failures forever; sessCtx would never be
// cancelled by the heartbeat.
func TestEnvironmentWorker_HeartbeatPerRequestTimeoutAndStalenessCeiling(t *testing.T) {
	server := newFakeWorkServer(t)

	var heartbeatCalls atomic.Int32
	server.HandleHeartbeat = func(w http.ResponseWriter, r *http.Request) {
		n := heartbeatCalls.Add(1)
		if n == 1 {
			// First beat succeeds with a tight TTL so the test runs fast:
			// ttl becomes 2 s and interval shrinks to 1 s.
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"last_heartbeat":"2026-05-11T12:00:00Z","lease_extended":true,"state":"active","ttl_seconds":2,"type":"work_heartbeat"}`))
			return
		}
		// Every subsequent beat hangs. With the per-request-timeout fix the
		// client cancels at interval (1 s); the handler's request context
		// fires and the handler returns. Without the fix, this would block
		// for the SDK default (~10 min).
		select {
		case <-r.Context().Done():
			return
		case <-time.After(20 * time.Second):
			t.Error("hung heartbeat was never cancelled — per-request timeout did not fire")
		}
	}

	server.HandleSessionGet = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"agent":{"skills":[]}}`))
	}
	server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[],"first_id":null,"has_more":false,"last_id":null}`))
	}
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		// Stream stays open until the per-session ctx is cancelled — which
		// is what the staleness ceiling triggers.
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		require.True(t, ok)
		w.WriteHeader(http.StatusOK)
		flusher.Flush()
		<-r.Context().Done()
	}
	server.HandleStop = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "env_key",
		Workdir:        t.TempDir(),
		Logger:         silentLogger,
	})

	// Outer ctx is generous so a regression (hung beat blocking on the SDK
	// default 10-min timeout) doesn't masquerade as "test finished within
	// the deadline" — we explicitly cap the expected runtime below.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	start := time.Now()
	// HandleItem returns nil here — the runner's Err() is documented to
	// return nil when iteration ended via consumer ctx cancellation (which
	// is what the staleness ceiling triggers via sessCtx). The two
	// behavioural assertions below are what proves the fixes work.
	_ = worker.HandleItem(ctx, HandleItemOptions{
		WorkID:        "work_1",
		EnvironmentID: "env_1",
		SessionID:     "sesn_test",
	})
	elapsed := time.Since(start)

	// With the fix, the staleness ceiling fires once the 2 s ttl has
	// elapsed since the first (and only) successful beat, terminating
	// HandleItem in roughly that window. 8 s is comfortably below the 15 s
	// outer ctx, so a regression where only the outer ctx terminates the
	// run fails this check.
	require.Less(t, elapsed, 8*time.Second, "HandleItem ran too long — staleness ceiling didn't terminate the heartbeat (got %s)", elapsed)

	// At least 2 heartbeat calls must have landed: the first success and at
	// least one timeout-cancelled retry. Without the per-request-timeout
	// fix, the second call would still be blocked on the SDK's ~10 min
	// default when HandleItem returns (via the outer ctx), so the count
	// would be 1.
	require.GreaterOrEqual(t, heartbeatCalls.Load(), int32(2), "expected at least 2 heartbeats; got %d — per-request timeout likely missing", heartbeatCalls.Load())
}
