package environments

// Tests for the per-item work secret: the sessions token extracted from a
// claimed item's secret payload is preferred over the environment key for
// that item's heartbeat / force-stop / skill-download / session calls, while
// polling stays on the environment key. The payload and the token must never
// reach the logs — asserted on happy and error paths alike.

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// encodeSecret builds a work-item secret the way the control plane does:
// URL-safe base64 of a JSON payload, without padding.
func encodeSecret(t *testing.T, payload map[string]any) string {
	t.Helper()
	body, err := json.Marshal(payload)
	require.NoError(t, err)
	return base64.RawURLEncoding.EncodeToString(body)
}

func TestSessionsTokenFromSecret(t *testing.T) {
	body, err := json.Marshal(map[string]any{
		"sessions_token": "sessions-token-abc",
		"auth":           []map[string]string{{"type": "anthropic_oauth", "token": "oauth"}},
	})
	require.NoError(t, err)
	unpadded := base64.RawURLEncoding.EncodeToString(body)
	padded := base64.URLEncoding.EncodeToString(body)

	// The payload decodes from URL-safe base64 JSON, with or without padding.
	require.Equal(t, "sessions-token-abc", sessionsTokenFromSecret(unpadded))
	require.Equal(t, "sessions-token-abc", sessionsTokenFromSecret(padded))
	// Anything malformed or token-less resolves to "".
	require.Empty(t, sessionsTokenFromSecret(""))
	require.Empty(t, sessionsTokenFromSecret("!!! not base64 !!!"))
	// Valid base64 but not JSON / not an object / missing or empty token.
	require.Empty(t, sessionsTokenFromSecret(base64.RawURLEncoding.EncodeToString([]byte("plain text"))))
	require.Empty(t, sessionsTokenFromSecret(encodeSecret(t, map[string]any{"session_ingress_token": "ingress-token-1"})))
	require.Empty(t, sessionsTokenFromSecret(encodeSecret(t, map[string]any{"sessions_token": ""})))
}

// scriptWorkServer wires the handlers a full worker pass needs: one poll
// yielding work, ack, heartbeat, the skill-setup session lookup, an event
// stream that terminates the session once a heartbeat landed, and stop.
func scriptWorkServer(t *testing.T, server *fakeWorkServer, work string) {
	t.Helper()
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

	heartbeatSeen := make(chan struct{})
	var heartbeatOnce sync.Once
	server.HandleHeartbeat = func(w http.ResponseWriter, _ *http.Request) {
		heartbeatOnce.Do(func() { close(heartbeatSeen) })
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"last_heartbeat":"2026-05-11T12:00:00Z","lease_extended":true,"state":"active","ttl_seconds":30,"type":"work_heartbeat"}`))
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
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		require.True(t, ok, "stream response writer must support flushing")
		w.WriteHeader(http.StatusOK)
		flusher.Flush()
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
	server.HandleStop = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}

// bearersByLeg partitions the recorded calls into the polling leg (poll/ack),
// the per-item leg (heartbeat, session lookup, event stream/list/send), and
// the worker's exit-path stop, and returns the set of Authorization values
// seen on each.
func bearersByLeg(calls []recordedCall) (poll, item, stop map[string]bool) {
	poll, item, stop = map[string]bool{}, map[string]bool{}, map[string]bool{}
	for _, c := range calls {
		switch {
		case strings.HasSuffix(c.path, "/work/poll"), strings.HasSuffix(c.path, "/ack"):
			poll[c.auth] = true
		case strings.HasSuffix(c.path, "/stop"):
			stop[c.auth] = true
		case strings.HasSuffix(c.path, "/heartbeat"), strings.Contains(c.path, "/sessions/"):
			item[c.auth] = true
		}
	}
	return poll, item, stop
}

func TestEnvironmentWorker_PrefersWorkItemSecret(t *testing.T) {
	// A claimed item carrying a per-item secret payload authenticates that
	// item's heartbeat / force-stop / skill-download / session-runner calls
	// with the sessions token extracted from it; polling stays on the
	// environment key. Neither the payload nor the token reaches the logs.
	secret := encodeSecret(t, map[string]any{
		"sessions_token":        "sessions-token-item-1",
		"session_ingress_token": "ingress-token-1",
	})
	server := newFakeWorkServer(t)
	scriptWorkServer(t, server, workJSONWithSecret("work_1", "env_1", "session", secret))

	var logBuf bytes.Buffer

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// End the run after the worker's own force-stop: the second poll returns
	// 204, and cancelling there lets Run return nil promptly.
	stopSeen := make(chan struct{})
	var stopOnce sync.Once
	prevStop := server.HandleStop
	server.HandleStop = func(w http.ResponseWriter, r *http.Request) {
		prevStop(w, r)
		stopOnce.Do(func() { close(stopSeen) })
	}
	go func() {
		select {
		case <-stopSeen:
		case <-ctx.Done():
		}
		cancel()
	}()

	worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "env_key",
		Workdir:        t.TempDir(),
		Logger:         slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelDebug})),
	})
	require.NoError(t, worker.Run(ctx))

	poll, item, stop := bearersByLeg(server.Calls())
	require.Equal(t, map[string]bool{"Bearer env_key": true}, poll,
		"polling must stay on the environment key")
	require.Equal(t, map[string]bool{"Bearer sessions-token-item-1": true}, item,
		"per-item calls must switch to the sessions token")
	require.Equal(t, map[string]bool{"Bearer sessions-token-item-1": true}, stop,
		"the worker force-stop must carry the sessions token and be the item's only stop")

	// The credentials are opaque — neither the payload nor the extracted
	// token may appear in log output.
	logs := logBuf.String()
	require.NotContains(t, logs, secret)
	require.NotContains(t, logs, "sessions-token-item-1")
}

func TestEnvironmentWorker_FallsBackWhenSecretUndecodable(t *testing.T) {
	// A secret payload that doesn't decode falls back to the environment
	// key, with a warning that doesn't include the payload itself.
	server := newFakeWorkServer(t)
	scriptWorkServer(t, server, workJSON("work_1", "env_1", "session"))

	var logBuf bytes.Buffer
	worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		EnvironmentKey: "env_key",
		Workdir:        t.TempDir(),
		Logger:         slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelDebug})),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	require.NoError(t, worker.HandleItem(ctx, HandleItemOptions{
		WorkID:        "work_1",
		EnvironmentID: "env_1",
		SessionID:     "sesn_test",
		WorkSecret:    "not-a-valid-payload",
	}))

	_, item, stop := bearersByLeg(server.Calls())
	require.Equal(t, map[string]bool{"Bearer env_key": true}, item,
		"an undecodable secret must fall back to the environment key")
	require.Equal(t, map[string]bool{"Bearer env_key": true}, stop,
		"the force-stop must fall back to the environment key too")

	logs := logBuf.String()
	require.Contains(t, logs, "no sessions token could be extracted")
	require.NotContains(t, logs, "not-a-valid-payload")
}

func TestEnvironmentWorker_SecretNotLoggedOnErrorPaths(t *testing.T) {
	// Error-path logging (heartbeat shutdown, force-stop failure) must not
	// leak the per-item secret payload or its sessions token.
	secret := encodeSecret(t, map[string]any{"sessions_token": "sessions-token-err"})
	server := newFakeWorkServer(t)
	scriptWorkServer(t, server, workJSONWithSecret("work_1", "env_1", "session", secret))
	// The heartbeat reports "stopping" so the shutdown branch logs, and the
	// force-stop fails so its error branch logs too.
	server.HandleHeartbeat = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"last_heartbeat":"2026-05-11T12:00:00Z","lease_extended":true,"state":"stopping","ttl_seconds":30,"type":"work_heartbeat"}`))
	}
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		require.True(t, ok, "stream response writer must support flushing")
		w.WriteHeader(http.StatusOK)
		flusher.Flush()
		<-r.Context().Done() // the heartbeat's shutdown cancels the session
	}
	server.HandleStop = func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"error":{"type":"api_error","message":"force-stop exploded"}}`, http.StatusInternalServerError)
	}

	var logBuf bytes.Buffer
	worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		EnvironmentKey: "env_key",
		Workdir:        t.TempDir(),
		Logger:         slog.New(slog.NewTextHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelDebug})),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = worker.HandleItem(ctx, HandleItemOptions{
		WorkID:        "work_1",
		EnvironmentID: "env_1",
		SessionID:     "sesn_test",
		WorkSecret:    secret,
	})

	logs := logBuf.String()
	// The error branches actually ran...
	require.Contains(t, logs, "force-stop on exit failed")
	// ...and none of them leaked the payload or the token inside it.
	require.NotContains(t, logs, secret)
	require.NotContains(t, logs, "sessions-token-err")
}

func TestHandleItem_UsesWorkSecretArgument(t *testing.T) {
	// An explicit WorkSecret supplies the per-item Bearer credential (its
	// sessions token); EnvironmentKey is still required but only a fallback.
	secret := encodeSecret(t, map[string]any{"sessions_token": "sessions-token-arg"})
	server := newFakeWorkServer(t)
	scriptWorkServer(t, server, workJSON("work_1", "env_1", "session"))

	worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		Workdir: t.TempDir(),
		Logger:  silentLogger,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	require.NoError(t, worker.HandleItem(ctx, HandleItemOptions{
		WorkID:         "work_1",
		EnvironmentID:  "env_1",
		SessionID:      "sesn_test",
		EnvironmentKey: "env_key",
		WorkSecret:     secret,
	}))

	_, item, stop := bearersByLeg(server.Calls())
	require.Equal(t, map[string]bool{"Bearer sessions-token-arg": true}, item)
	require.Equal(t, map[string]bool{"Bearer sessions-token-arg": true}, stop)
}

func TestHandleItem_FallsBackToWorkSecretEnvVar(t *testing.T) {
	// WorkSecret falls back to ANTHROPIC_WORK_SECRET (the env var the
	// `worker poll --on-work` command sets alongside the others); when
	// neither is present the environment key is used — the existing
	// HandleItem tests cover that path.
	secret := encodeSecret(t, map[string]any{"sessions_token": "sessions-token-env"})
	server := newFakeWorkServer(t)
	scriptWorkServer(t, server, workJSON("w_env", "e_env", "session"))

	t.Setenv("ANTHROPIC_WORK_ID", "w_env")
	t.Setenv("ANTHROPIC_ENVIRONMENT_ID", "e_env")
	t.Setenv("ANTHROPIC_SESSION_ID", "sesn_test")
	t.Setenv("ANTHROPIC_ENVIRONMENT_KEY", "key_env")
	t.Setenv("ANTHROPIC_WORK_SECRET", secret)

	worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		Workdir: t.TempDir(),
		Logger:  silentLogger,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	require.NoError(t, worker.HandleItem(ctx, HandleItemOptions{}))

	_, item, stop := bearersByLeg(server.Calls())
	require.Equal(t, map[string]bool{"Bearer sessions-token-env": true}, item)
	require.Equal(t, map[string]bool{"Bearer sessions-token-env": true}, stop)
}
