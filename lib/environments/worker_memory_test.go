package environments

// Memory-store sync driven end to end through the real EnvironmentWorker.
//
// Real: the worker, SessionMemoryStores, the filestore, and the filesystem.
// Faked: the network (the memoryServer behind the fake work server) and the
// clock. So the chain worker → gate on the sessions token → download →
// SyncIfDue per tool call → Finish → Dispose runs for real against a fake
// control plane; assertions read the server and the disk.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/tools/agenttoolset"
	"github.com/stretchr/testify/require"
)

// memoryTestTool is a minimal anthropic.BetaTool whose Execute delegates to
// the test — the agent's stand-in for editing its notes folder.
type memoryTestTool struct {
	exec    func(input json.RawMessage)
	onClose func()
}

func (m *memoryTestTool) Name() string        { return "note_edit" }
func (m *memoryTestTool) Description() string { return "edit the note" }
func (m *memoryTestTool) InputSchema() anthropic.BetaToolInputSchemaParam {
	return anthropic.BetaToolInputSchemaParam{Properties: map[string]any{}}
}
func (m *memoryTestTool) Execute(_ context.Context, input json.RawMessage) ([]anthropic.BetaToolResultBlockParamContentUnion, error) {
	m.exec(input)
	return []anthropic.BetaToolResultBlockParamContentUnion{{OfText: &anthropic.BetaTextBlockParam{Text: "ok"}}}, nil
}
func (m *memoryTestTool) Close() error {
	if m.onClose != nil {
		m.onClose()
	}
	return nil
}

func toolUseSSE(i int) string {
	return fmt.Sprintf("event: agent.tool_use\n"+
		`data: {"type":"agent.tool_use","id":"evt_%d","name":"note_edit","input":{"i":%d},"processed_at":"2026-05-11T12:00:00Z"}`+
		"\n\n", i, i)
}

const terminatedSSE = "event: session.status_terminated\n" +
	`data: {"type":"session.status_terminated","id":"evt_term","processed_at":"2026-05-11T12:00:00Z"}` +
	"\n\n"

// memoryWorkServer scripts a fakeWorkServer around a memoryServer: session
// lookups return the store-carrying session (and are counted), memory calls
// hit the in-memory store, and the stream dispatches toolCalls note_edit
// calls before ending with terminate (or holding open when terminate is
// false, for the unclean-end test).
//
// The SessionToolRunner executes tool calls eagerly as events arrive, ahead
// of the consumer loop that runs SyncIfDue — so to observe one sync per
// call the stream must not emit the next tool_use until the effects of the
// previous one are visible. gate(i) is polled after tool_use i's result was
// posted; the next event is only sent once it returns true (nil gates on the
// post alone).
func memoryWorkServer(t *testing.T, mem *memoryServer, toolCalls int, terminate bool, gate func(i int) bool) *fakeWorkServer {
	t.Helper()
	server := newFakeWorkServer(t)
	server.HandleSessionGet = func(w http.ResponseWriter, r *http.Request) {
		mem.serveHTTP(w, r) // records the retrieve, returns the session JSON
	}
	server.HandleMemories = mem.handleMemories
	server.HandleHeartbeat = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"last_heartbeat":"2026-05-11T12:00:00Z","lease_extended":true,"state":"active","ttl_seconds":30,"type":"work_heartbeat"}`))
	}
	server.HandleList = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[],"first_id":null,"has_more":false,"last_id":null}`))
	}
	var sends atomic.Int32
	server.HandleSend = func(w http.ResponseWriter, _ *http.Request) {
		sends.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"type":"send_session_events"}`))
	}
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		flusher, ok := w.(http.Flusher)
		require.True(t, ok)
		w.WriteHeader(http.StatusOK)
		flusher.Flush()
		for i := 0; i < toolCalls; i++ {
			_, _ = w.Write([]byte(toolUseSSE(i)))
			flusher.Flush()
			// Wait for tool_use i's result post, then for gate(i), before
			// emitting the next event.
			for int(sends.Load()) <= i || (gate != nil && !gate(i)) {
				select {
				case <-r.Context().Done():
					return
				case <-time.After(2 * time.Millisecond):
				}
			}
		}
		if terminate {
			_, _ = w.Write([]byte(terminatedSSE))
			flusher.Flush()
			return
		}
		<-r.Context().Done()
	}
	server.HandleStop = func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
	return server
}

// memoryWorker builds an EnvironmentWorker whose one tool writes
// "edit {i}" into the store dir (when writeEdits is set — the gate tests run
// tool calls without touching the never-downloaded folder) and then advances
// the fake clock by advanceBy — so a nonzero advance makes the very next
// SyncIfDue due. opts.Logger is honored when the test wants to read the log.
func memoryWorker(t *testing.T, server *fakeWorkServer, workdir string, writeEdits bool, advanceBy time.Duration, opts EnvironmentWorkerOptions) *EnvironmentWorker {
	t.Helper()
	// The tool advances the clock on the runner's dispatch goroutine while
	// the worker loop reads it, so the fake clock takes a lock.
	var mu sync.Mutex
	now := time.Unix(0, 0)
	local := filepath.Join(workdir, "memory", "notes")
	tool := &memoryTestTool{exec: func(input json.RawMessage) {
		var in struct{ I int }
		require.NoError(t, json.Unmarshal(input, &in))
		if writeEdits {
			require.NoError(t, os.WriteFile(filepath.Join(local, "note.md"), fmt.Appendf(nil, "edit %d", in.I), 0o600))
		}
		mu.Lock()
		now = now.Add(advanceBy)
		mu.Unlock()
	}}
	opts.Workdir = workdir
	if opts.Logger == nil {
		opts.Logger = silentLogger
	}
	opts.ToolsFunc = func(*agenttoolset.AgentToolContext) []anthropic.BetaTool { return []anthropic.BetaTool{tool} }
	worker := NewEnvironmentWorker(server.Client(), opts)
	worker.memoryClock = func() time.Time {
		mu.Lock()
		defer mu.Unlock()
		return now
	}
	return worker
}

// runMemoryItem runs one work item for session s1 under a 15s bound and
// returns what HandleItem returned.
func runMemoryItem(worker *EnvironmentWorker, secret string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return worker.HandleItem(ctx, HandleItemOptions{
		WorkID:         "work_1",
		EnvironmentID:  "env_1",
		SessionID:      "s1",
		EnvironmentKey: "env_key",
		WorkSecret:     secret,
	})
}

func handleMemoryItem(t *testing.T, worker *EnvironmentWorker, secret string) {
	t.Helper()
	require.NoError(t, runMemoryItem(worker, secret))
}

// reportStoppingOnceWritten makes the heartbeat report "stopping" once file
// holds want — the tool call has run — so the session is cancelled
// mid-stream and ends uncleanly.
func reportStoppingOnceWritten(server *fakeWorkServer, file, want string) {
	toolRan := make(chan struct{})
	server.HandleHeartbeat = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		select {
		case <-toolRan:
			_, _ = w.Write([]byte(`{"last_heartbeat":"2026-05-11T12:00:00Z","lease_extended":true,"state":"stopping","ttl_seconds":2,"type":"work_heartbeat"}`))
		default:
			_, _ = w.Write([]byte(`{"last_heartbeat":"2026-05-11T12:00:00Z","lease_extended":true,"state":"active","ttl_seconds":2,"type":"work_heartbeat"}`))
		}
	}
	go func() {
		for {
			if got, err := os.ReadFile(file); err == nil && string(got) == want {
				close(toolRan)
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
}

// shortenFlushTimeout sets [MemoryFlushTimeout] to d for the rest of the test.
func shortenFlushTimeout(t *testing.T, d time.Duration) {
	t.Helper()
	old := MemoryFlushTimeout
	MemoryFlushTimeout = d
	t.Cleanup(func() { MemoryFlushTimeout = old })
}

func TestEnvironmentWorker_MemorySyncsOnCadenceAndDisposes(t *testing.T) {
	mem := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	server := memoryWorkServer(t, mem, 3, true, func(i int) bool {
		// Before dispatching call i+1, the cadence sync for call i (which
		// advanced the clock) must have pushed its edit; call 2 does not
		// advance the clock, so it gates on the result post alone.
		if i >= 2 {
			return true
		}
		return slices.Contains(mem.ReceivedSnapshot(), updated("note.md", fmt.Sprintf("edit %d", i), preimage(i)))
	})
	workdir := t.TempDir()
	local := filepath.Join(workdir, "memory", "notes")

	// Calls 0 and 1 advance the clock past the interval, so their SyncIfDue
	// polls are due; call 2 does not, so its edit is only pushed by the
	// unconditional final sync at the clean stream end.
	worker := memoryWorker(t, server, workdir, true, 400*time.Second, EnvironmentWorkerOptions{
		MemorySyncInterval: 300 * time.Second,
	})
	handleMemoryItem(t, worker, encodeSecret(t, map[string]any{"sessions_token": "tok"}))

	// One session fetch, shared by the skills and memory downloads.
	require.Equal(t, []string{"s1"}, mem.Retrieves)
	// Two cadence syncs pushed edits 0/1; the final sync pushed edit 2.
	require.Equal(t, []string{
		updated("note.md", "edit 0", "v1"),
		updated("note.md", "edit 1", "edit 0"),
		updated("note.md", "edit 2", "edit 1"),
	}, mem.ReceivedSnapshot())
	// Dispose removed the store dir the download created.
	require.NoDirExists(t, local)
}

// The file tools must reach the session's memory folders wherever they are
// mounted, so the worker lists every store's folder in AllowedRoots — the
// default location under the workdir and a mount_path outside it alike.
func TestEnvironmentWorker_GrantsFileToolsAccessToStoreFolders(t *testing.T) {
	outsideMount := filepath.Join(t.TempDir(), "mnt", "notes")
	tests := []struct {
		description string
		mountPath   string
		wantRoot    func(workdir string) string
	}{
		{
			description: "a store with no mount_path lands under the workdir and is listed",
			wantRoot:    func(workdir string) string { return filepath.Join(workdir, "memory", "notes") },
		},
		{
			description: "a store mounted outside the workdir is listed at its mount_path",
			mountPath:   outsideMount,
			wantRoot:    func(string) string { return outsideMount },
		},
	}
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			mem := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
			mem.MountPath = tc.mountPath
			server := memoryWorkServer(t, mem, 0, true, nil)
			workdir := t.TempDir()
			worker := memoryWorker(t, server, workdir, false, 0, EnvironmentWorkerOptions{})
			var got *agenttoolset.AgentToolContext
			toolsFunc := worker.opts.ToolsFunc
			worker.opts.ToolsFunc = func(env *agenttoolset.AgentToolContext) []anthropic.BetaTool {
				got = env
				return toolsFunc(env)
			}

			handleMemoryItem(t, worker, encodeSecret(t, map[string]any{"sessions_token": "tok"}))

			require.NotNil(t, got)
			require.Equal(t, []string{tc.wantRoot(workdir)}, got.AllowedRoots)
			require.Empty(t, got.ReadOnlyRoots)
		})
	}
}

func TestEnvironmentWorker_MissingSessionsTokenFailsTheItem(t *testing.T) {
	// The session has memory stores attached but the work item carried no
	// sessions token, so its memories cannot be mounted. The item fails
	// before a single tool call is served — a hosted sandbox refuses to start
	// in this state, and running here without memories would silently diverge
	// from it.
	mem := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	server := memoryWorkServer(t, mem, 2, true, nil)
	// Any memory_stores call would fail the test as unscripted.
	server.HandleMemories = nil
	workdir := t.TempDir()

	worker := memoryWorker(t, server, workdir, false, 400*time.Second, EnvironmentWorkerOptions{})
	err := runMemoryItem(worker, "")

	require.ErrorIs(t, err, ErrSessionMemoryNoToken)
	require.ErrorContains(t, err, "carried no sessions token")

	// No download, no sync, no store dir, and the session was never served...
	require.Empty(t, mem.ReceivedSnapshot())
	require.NoDirExists(t, filepath.Join(workdir, "memory"))
	calls := server.Calls()
	require.False(t, slices.ContainsFunc(calls, func(c recordedCall) bool {
		return strings.HasSuffix(c.path, "/events/stream")
	}), "the session ran")
	// ...but the deferred force-stop still ran.
	require.True(t, slices.ContainsFunc(calls, func(c recordedCall) bool {
		return strings.HasSuffix(c.path, "/stop")
	}), "force-stop was not sent")
}

func TestEnvironmentWorker_MemorySkippedWhenIntervalNegative(t *testing.T) {
	mem := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	server := memoryWorkServer(t, mem, 2, true, nil)
	server.HandleMemories = nil
	workdir := t.TempDir()

	var logBuf bytes.Buffer
	worker := memoryWorker(t, server, workdir, false, 400*time.Second, EnvironmentWorkerOptions{
		MemorySyncInterval: -1,
		Logger:             captureLog(&logBuf),
	})
	handleMemoryItem(t, worker, encodeSecret(t, map[string]any{"sessions_token": "tok"}))

	// Even with a sessions token, a negative interval turns memory off
	// entirely.
	require.Empty(t, mem.ReceivedSnapshot())
	require.NoDirExists(t, filepath.Join(workdir, "memory"))
	// Turning memory off is a deliberate opt-out, so it says nothing loud —
	// even though this session does have a store attached.
	require.NotContains(t, logBuf.String(), "level=ERROR")
	require.NotContains(t, logBuf.String(), "has memory stores attached")
}

func TestEnvironmentWorker_MemoryDownloadFailureFailsTheItem(t *testing.T) {
	// A folder already at the store's path is debris from a run that died. The
	// download refuses it and the work item fails, rather than serving a
	// session whose system prompt names a memory folder holding someone else's
	// files.
	mem := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	server := memoryWorkServer(t, mem, 0, true, nil)
	// The refusal precedes any listing, so a memory_stores call would be a bug.
	server.HandleMemories = nil
	workdir := t.TempDir()
	debris := filepath.Join(workdir, "memory", "notes")
	require.NoError(t, os.MkdirAll(debris, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(debris, "left-behind.md"), []byte("from a dead run"), 0o600))

	worker := memoryWorker(t, server, workdir, false, 0, EnvironmentWorkerOptions{})
	err := runMemoryItem(worker, encodeSecret(t, map[string]any{"sessions_token": "tok"}))

	var memErr *SessionMemoryError
	require.ErrorAs(t, err, &memErr)
	require.ErrorContains(t, err, "already exists")

	// Nothing was pushed, the debris is exactly as it was found...
	require.Empty(t, mem.ReceivedSnapshot())
	require.Equal(t, "from a dead run", readLocal(t, debris, "left-behind.md"))
	// ...and the deferred force-stop still ran.
	require.True(t, slices.ContainsFunc(server.Calls(), func(c recordedCall) bool {
		return strings.HasSuffix(c.path, "/stop")
	}), "force-stop was not sent")
}

func TestEnvironmentWorker_MemoryFlushesWritesOnUncleanEnd(t *testing.T) {
	mem := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	// The stream never terminates; instead the heartbeat reports "stopping"
	// once the tool call has run, cancelling the session mid-run.
	server := memoryWorkServer(t, mem, 1, false, nil)
	workdir := t.TempDir()
	local := filepath.Join(workdir, "memory", "notes")
	reportStoppingOnceWritten(server, filepath.Join(local, "note.md"), "edit 0")

	// The tool edits the note WITHOUT advancing the clock, so the per-call
	// SyncIfDue is never due — the unclean end skips the final sync, but
	// the shutdown flush pushes the pending edit before Dispose removes
	// the folder.
	worker := memoryWorker(t, server, workdir, true, 0, EnvironmentWorkerOptions{
		MemorySyncInterval: 300 * time.Second,
	})
	_ = runMemoryItem(worker, encodeSecret(t, map[string]any{"sessions_token": "tok"}))

	// The cancel skipped the final sync, but the shutdown flush pushed
	// the pending edit before Dispose removed the folder.
	require.Equal(t, []string{updated("note.md", "edit 0", "v1")}, mem.ReceivedSnapshot())
	require.NoDirExists(t, local)
}

func TestEnvironmentWorker_MemoryCleanEndSyncsOnceAndTheFlushAddsNothing(t *testing.T) {
	// A clean end runs the final sync exactly once, and the teardown
	// flush that follows finds nothing dirty and sends nothing more.
	mem := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	server := memoryWorkServer(t, mem, 1, true, nil)
	workdir := t.TempDir()
	local := filepath.Join(workdir, "memory", "notes")

	var logBuf bytes.Buffer
	worker := memoryWorker(t, server, workdir, true, 0, EnvironmentWorkerOptions{
		MemorySyncInterval: 300 * time.Second,
		Logger:             captureLog(&logBuf),
	})
	handleMemoryItem(t, worker, encodeSecret(t, map[string]any{"sessions_token": "tok"}))

	// The clean-end final sync pushed the edit — exactly once. The
	// teardown flush runs after it but finds nothing dirty, so nothing
	// more is sent.
	require.Equal(t, []string{updated("note.md", "edit 0", "v1")}, mem.ReceivedSnapshot())
	// Neither teardown pass hit its bound.
	require.NotContains(t, logBuf.String(), "cut off")
	require.NoDirExists(t, local)
}

func TestEnvironmentWorker_MemoryFlushIsBoundedAndTeardownStillRuns(t *testing.T) {
	// A flush that hangs is cut off at MemoryFlushTimeout and says so;
	// Dispose still runs.
	mem := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	server := memoryWorkServer(t, mem, 1, false, nil)
	workdir := t.TempDir()
	local := filepath.Join(workdir, "memory", "notes")
	reportStoppingOnceWritten(server, filepath.Join(local, "note.md"), "edit 0")
	shortenFlushTimeout(t, 50*time.Millisecond)

	var flushStarted atomic.Bool
	// Every sync listing hangs; the download lists with the full view and is
	// unaffected. No cadence sync is due (the clock never advances) and the
	// "stopping" heartbeat makes the end unclean, so the flush is the only
	// sync that runs into it.
	realList := server.HandleMemories
	server.HandleMemories = func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Query().Get("view") == "basic" {
			flushStarted.Store(true)
			<-r.Context().Done()
			return
		}
		realList(w, r)
	}

	var logBuf bytes.Buffer
	worker := memoryWorker(t, server, workdir, true, 0, EnvironmentWorkerOptions{
		MemorySyncInterval: 300 * time.Second,
		Logger:             captureLog(&logBuf),
	})
	start := time.Now()
	_ = runMemoryItem(worker, encodeSecret(t, map[string]any{"sessions_token": "tok"}))

	require.True(t, flushStarted.Load())
	// Without the timeout the hung flush would hold teardown until the
	// test's own 15s bound.
	require.Less(t, time.Since(start), 10*time.Second)
	// The pass says it was cut off, on a line carrying the item's ids; the
	// listing was what hung, so every changed file counts as not uploaded.
	require.Contains(t, logBuf.String(),
		`msg="memory flush cut off after 50ms; changed files it had not uploaded yet are not saved" work_id=work_1 session_id=s1`)
	require.Contains(t, logBuf.String(), "memory flush cut off part-way; 1 of 1 changed files had not finished uploading")
	require.NotContains(t, logBuf.String(), "memory flush failed")
	require.NotContains(t, logBuf.String(), "final memory sync cut off")
	require.NoDirExists(t, local)
}

func TestEnvironmentWorker_MemoryFinalSyncIsBoundedAndTheFlushStillRuns(t *testing.T) {
	// A final sync that hangs is cut off at MemoryFlushTimeout and says so;
	// the flush after it still uploads the last edit.
	mem := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	server := memoryWorkServer(t, mem, 1, true, nil)
	workdir := t.TempDir()
	local := filepath.Join(workdir, "memory", "notes")
	// The same bound covers the real flush that follows, which must finish
	// inside it.
	shortenFlushTimeout(t, 500*time.Millisecond)

	// Only the final sync's listing hangs: the download lists with the full
	// view, no cadence sync is due, and the flush's listing comes second.
	var finalSyncListed atomic.Bool
	realList := server.HandleMemories
	server.HandleMemories = func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Query().Get("view") == "basic" && finalSyncListed.CompareAndSwap(false, true) {
			<-r.Context().Done()
			return
		}
		realList(w, r)
	}

	var logBuf bytes.Buffer
	worker := memoryWorker(t, server, workdir, true, 0, EnvironmentWorkerOptions{
		MemorySyncInterval: 300 * time.Second,
		Logger:             captureLog(&logBuf),
	})
	start := time.Now()
	handleMemoryItem(t, worker, encodeSecret(t, map[string]any{"sessions_token": "tok"}))

	require.Less(t, time.Since(start), 10*time.Second)
	require.Contains(t, logBuf.String(), "final memory sync cut off after 500ms; the flush that follows still uploads changed files")
	// The flush finished in time, so only the final sync reports a cut-off.
	require.NotContains(t, logBuf.String(), "memory flush cut off")
	require.Equal(t, []string{updated("note.md", "edit 0", "v1")}, mem.ReceivedSnapshot())
	require.Equal(t, "edit 0", mem.Files["note.md"])
	require.NoDirExists(t, local)
}

func TestEnvironmentWorker_MemoryLogsWhatACutOffFlushLeftBehind(t *testing.T) {
	// One upload never completes: the flush is cut off, the store logs that
	// one of its two changed files did not make it, and the other one did.
	mem := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	server := memoryWorkServer(t, mem, 1, false, nil)
	workdir := t.TempDir()
	local := filepath.Join(workdir, "memory", "notes")
	shortenFlushTimeout(t, 500*time.Millisecond)

	tool := &memoryTestTool{exec: func(json.RawMessage) {
		require.NoError(t, os.WriteFile(filepath.Join(local, "kept.md"), []byte("saved"), 0o600))
		require.NoError(t, os.WriteFile(filepath.Join(local, "stuck.md"), []byte("only copy"), 0o600))
	}}
	reportStoppingOnceWritten(server, filepath.Join(local, "stuck.md"), "only copy")
	mem.SetUploadHook(func(ctx context.Context, path string) {
		if path == "stuck.md" {
			<-ctx.Done()
		}
	})

	var logBuf bytes.Buffer
	worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		Workdir:            workdir,
		MemorySyncInterval: 300 * time.Second,
		Logger:             captureLog(&logBuf),
		ToolsFunc:          func(*agenttoolset.AgentToolContext) []anthropic.BetaTool { return []anthropic.BetaTool{tool} },
	})
	start := time.Now()
	_ = runMemoryItem(worker, encodeSecret(t, map[string]any{"sessions_token": "tok"}))

	require.Less(t, time.Since(start), 10*time.Second)
	require.Equal(t, "saved", mem.Files["kept.md"])
	require.NotContains(t, mem.Files, "stuck.md")
	// The store's line carries the item's ids and then the store id.
	require.Contains(t, logBuf.String(),
		`msg="memory flush cut off part-way; 1 of 2 changed files had not finished uploading" work_id=work_1 session_id=s1 memory_store_id=memstore_notes`)
	require.Contains(t, logBuf.String(), "memory flush cut off after 500ms; changed files it had not uploaded yet are not saved")
	require.NotContains(t, logBuf.String(), "failed to upload memory")
	require.NoDirExists(t, local)
}

// preimage is the note.md content the cadence sync for edit i guards on.
func preimage(i int) string {
	if i == 0 {
		return "v1"
	}
	return fmt.Sprintf("edit %d", i-1)
}

func TestEnvironmentWorker_NoMemoryStoreAttachedStaysQuiet(t *testing.T) {
	// Same missing sessions token as above, but this session has no memory to
	// run without. Nothing is misconfigured, so nothing is said loudly.
	mem := newMemoryServer(t, map[string]string{}, "")
	mem.NoStores = true
	server := memoryWorkServer(t, mem, 2, true, nil)
	server.HandleMemories = nil
	workdir := t.TempDir()

	var logBuf bytes.Buffer
	worker := memoryWorker(t, server, workdir, false, 0, EnvironmentWorkerOptions{
		Logger: captureLog(&logBuf),
	})
	handleMemoryItem(t, worker, "")

	require.NoDirExists(t, filepath.Join(workdir, "memory"))
	require.NotContains(t, logBuf.String(), "level=ERROR")
	require.Contains(t, logBuf.String(), "memory stores disabled for this item")
}

func TestEnvironmentWorker_MemoryCanTurnSyncDeletionsOff(t *testing.T) {
	// MemorySyncDeletions = MemorySyncDeletionsDisabled reaches the stores: a
	// local deletion never propagates however many syncs run, while
	// uploads still do.
	mem := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	server := memoryWorkServer(t, mem, 3, true, nil)
	workdir := t.TempDir()
	local := filepath.Join(workdir, "memory", "notes")

	tool := &memoryTestTool{}
	var mu sync.Mutex
	now := time.Unix(0, 0)
	tool.exec = func(input json.RawMessage) {
		var in struct{ I int }
		require.NoError(t, json.Unmarshal(input, &in))
		if in.I == 0 {
			require.NoError(t, os.Remove(filepath.Join(local, "note.md")))
			require.NoError(t, os.WriteFile(filepath.Join(local, "new.md"), []byte("survives"), 0o600))
		}
		// Every gap dwarfs both the sync interval and the corroboration window.
		mu.Lock()
		now = now.Add(400 * time.Second)
		mu.Unlock()
	}
	worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		Workdir:             workdir,
		MemorySyncInterval:  300 * time.Second,
		MemorySyncDeletions: MemorySyncDeletionsDisabled,
		Logger:              silentLogger,
		ToolsFunc:           func(*agenttoolset.AgentToolContext) []anthropic.BetaTool { return []anthropic.BetaTool{tool} },
	})
	worker.memoryClock = func() time.Time { mu.Lock(); defer mu.Unlock(); return now }

	handleMemoryItem(t, worker, encodeSecret(t, map[string]any{"sessions_token": "tok"}))

	require.Equal(t, "v1", mem.Files["note.md"])
	require.Equal(t, []string{created("new.md", "survives")}, mem.ReceivedSnapshot())
}

func TestEnvironmentWorker_TheLastSyncRunsAfterToolsAreClosed(t *testing.T) {
	// Finish skips the delete waiting window, so it must only run once
	// nothing can still be rewriting files — after CloseAll kills the
	// session's subprocesses. The test tool's Close writes a memory file:
	// were CloseAll to run after the syncs, the sync would not see it and
	// this file would never reach the server.
	mem := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	server := memoryWorkServer(t, mem, 1, true, nil)
	workdir := t.TempDir()
	local := filepath.Join(workdir, "memory", "notes")

	tool := &memoryTestTool{
		exec: func(input json.RawMessage) {
			require.NoError(t, os.WriteFile(filepath.Join(local, "note.md"), []byte("last edit"), 0o600))
		},
		onClose: func() {
			require.NoError(t, os.WriteFile(filepath.Join(local, "on-close.md"), []byte("written on close"), 0o600))
		},
	}
	worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		Workdir:            workdir,
		MemorySyncInterval: 300 * time.Second,
		Logger:             silentLogger,
		ToolsFunc:          func(*agenttoolset.AgentToolContext) []anthropic.BetaTool { return []anthropic.BetaTool{tool} },
	})

	handleMemoryItem(t, worker, encodeSecret(t, map[string]any{"sessions_token": "tok"}))

	// The tool was closed before Finish/Flush scanned the folder: the
	// file Close wrote is on the server.
	require.Equal(t, "written on close", mem.Files["on-close.md"])
	require.Equal(t, "last edit", mem.Files["note.md"])
	require.NoDirExists(t, local)
}
