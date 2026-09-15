package environments

// Tests for continuing a lease another process started: HandleItem sends the
// caller's lease token as its first expected_last_heartbeat, then echoes each
// response's last_heartbeat as usual, and the three outcomes the hand-off
// contract defines for that first heartbeat end the item the right way.

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// stoppedHeartbeat is the 200 the control plane answers once a lease lapsed
// without anyone re-claiming the item.
const stoppedHeartbeat = `{"last_heartbeat":"2026-05-11T12:00:01Z","lease_extended":false,"state":"stopped","ttl_seconds":30,"type":"work_heartbeat"}`

// handleItemEndedByHeartbeat services work_1 through HandleItem with opts,
// answering each heartbeat with answerHeartbeat once the event stream is open
// (see holdHeartbeatsUntilServing), and returns the server and the logs.
func handleItemEndedByHeartbeat(t *testing.T, opts HandleItemOptions, answerHeartbeat func(w http.ResponseWriter, beat int)) (*fakeWorkServer, string) {
	t.Helper()
	server := newFakeWorkServer(t)
	holdHeartbeatsUntilServing(t, server, answerHeartbeat)

	var logs lockedBuffer
	worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		EnvironmentKey: "env_key",
		Workdir:        t.TempDir(),
		Logger:         slog.New(slog.NewTextHandler(&logs, nil)),
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts.WorkID, opts.EnvironmentID, opts.SessionID = "work_1", "env_1", "sesn_test"
	require.NoError(t, worker.HandleItem(ctx, opts))
	return server, logs.String()
}

// The first heartbeat carries the caller's lease token when one is given and
// the NO_HEARTBEAT sentinel otherwise; every later heartbeat echoes the
// last_heartbeat the previous response returned.
func TestEnvironmentWorker_HandleItemFirstHeartbeatExpectation(t *testing.T) {
	tests := []struct {
		description string
		option      string
		envVar      string
		want        string
	}{
		{"starts the lease when no token is given", "", "", noHeartbeatSentinel},
		{"continues the lease from the option", "hb-poller", "", "hb-poller"},
		{"continues the lease from ANTHROPIC_WORK_LAST_HEARTBEAT", "", "hb-env", "hb-env"},
		{"the option wins over the environment variable", "hb-poller", "hb-env", "hb-poller"},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			t.Setenv("ANTHROPIC_WORK_LAST_HEARTBEAT", tc.envVar)
			server, _ := handleItemEndedByHeartbeat(t, HandleItemOptions{ExpectedLastHeartbeat: tc.option}, func(w http.ResponseWriter, beat int) {
				w.Header().Set("Content-Type", "application/json")
				if beat == 1 {
					_, _ = w.Write([]byte(shortTTLHeartbeat))
					return
				}
				_, _ = w.Write([]byte(stoppedHeartbeat))
			})

			beats := callsEndingIn(server.Calls(), "/heartbeat")
			require.Len(t, beats, 2)
			require.Contains(t, beats[0].query, "expected_last_heartbeat="+tc.want)
			require.Contains(t, beats[1].query, "expected_last_heartbeat="+url.QueryEscape("2026-05-11T12:00:00Z"))
		})
	}
}

// A token the server no longer holds for the item — the hand-off lapsed and
// the item was re-queued or claimed by another worker — releases the item on
// the first heartbeat without stopping it.
func TestEnvironmentWorker_HandleItemContinuedLeaseLostReleasesItemWithoutStoppingIt(t *testing.T) {
	server, logs := handleItemEndedByHeartbeat(t, HandleItemOptions{ExpectedLastHeartbeat: "hb-stale"}, func(w http.ResponseWriter, _ int) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusPreconditionFailed)
		_, _ = w.Write([]byte(leaseLostBody))
	})

	calls := server.Calls()
	beats := callsEndingIn(calls, "/heartbeat")
	require.Len(t, beats, 1)
	require.Contains(t, beats[0].query, "expected_last_heartbeat=hb-stale")
	require.Empty(t, callsEndingIn(calls, "/stop"), "an item someone else holds must not be stopped")
	require.Contains(t, logs, `msg="lease lost; released without stopping it" work_id=work_1 session_id=sesn_test`)
}

// A lease that lapsed with nobody re-claiming the item answers the first
// heartbeat with 200, lease_extended false and state stopped: the item ends
// after that one heartbeat and is force-stopped on the way out.
func TestEnvironmentWorker_HandleItemLapsedLeaseExits(t *testing.T) {
	server, logs := handleItemEndedByHeartbeat(t, HandleItemOptions{ExpectedLastHeartbeat: "hb-lapsed"}, func(w http.ResponseWriter, _ int) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(stoppedHeartbeat))
	})

	calls := server.Calls()
	require.Len(t, callsEndingIn(calls, "/heartbeat"), 1)
	require.Len(t, callsEndingIn(calls, "/stop"), 1)
	require.Contains(t, logs, `msg="heartbeat reports shutdown" work_id=work_1 session_id=sesn_test state=stopped`)
	require.NotContains(t, logs, "without stopping it")
}
