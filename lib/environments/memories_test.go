package environments

// The memory-store sync contract.
//
// An agent keeps a folder of notes on disk. The same folder lives on the
// server. Each sync makes the two agree: each side gets the other's
// changes, and when both changed the same file the server wins.
//
// Each test tells one story with two actors — `local` is the agent's folder
// on disk, `server` is the remote copy. Paths are the same string on both
// sides. server.Received holds only what the sync sent; arranging remote
// state never touches it.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/filestore"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stretchr/testify/require"
)

func shaHex(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}

// Requests recorded on memoryServer.Received, in domain terms: paths and
// content. `was` is the pre-image content the request's precondition guards
// on — the builder computes the sha.

func created(path, content string) string { return "create " + path + " " + content }

func updated(path, content, was string) string {
	return "update " + path + " " + content + " if=" + shaHex(was)
}

func deletedReq(path, was string) string { return "delete " + path + " if=" + shaHex(was) }

// memoryServer is the test's view of one remote memory store named "notes".
//
// Files is its current state (path → content); Received is every write the
// code under test sent it. Arranging via Write / Delete changes Files without
// touching Received — only the sync's own requests are recorded. It serves
// both the sessions lookup (recording Retrieves) and the memory_stores
// endpoints, so one server backs the pure SessionMemoryStores tests and the
// worker e2e tests alike.
type memoryServer struct {
	t *testing.T

	mu        sync.Mutex
	Files     map[string]string
	Received  []string
	Retrieves []string
	// ContentFetches is every single-memory content fetch (the sync's
	// content pass), by path, in arrival order — a test asserting "no
	// content moved" checks this stays empty.
	ContentFetches []string

	access   string // "" (read_write) or "read_only"
	OnUpdate func() // race hook: runs before an update is processed
	onList   func() // race hook: runs after a listing is served; set via SetOnList
	// onFetch is a race hook that runs as each content fetch arrives, with
	// the server mutex held (Files is safe to touch); a non-zero return
	// fails the fetch with that status. Set via SetOnFetch.
	onFetch func(path string) int
	// uploadHook runs as each create or update arrives, with the memory's
	// path, before the request is recorded and without the server mutex —
	// block in it to hold an upload open. A request whose client gave up
	// while held is dropped unrecorded. Set via SetUploadHook.
	uploadHook func(ctx context.Context, path string)
	// FailUpdates answers every update with a 500 — recorded, so the attempt
	// is visible, but never applied.
	FailUpdates bool
	// FailUploads answers every create or update with this status, before it
	// is recorded or applied. Zero means none.
	FailUploads int
	// MaxContentBytes answers any create/update whose content exceeds it
	// with a 400 — recorded, never applied. Zero means no cap.
	MaxContentBytes int
	// Broken attaches a second store, "memstore_broken", listed before the
	// healthy one, whose listing yields one memory and then explodes on the
	// next page — so a half-download is observable.
	Broken bool
	// MountPath is emitted on the "notes" store resource (null when empty).
	MountPath string
	// Noise adds a non-store resource to the session and a memory_prefix
	// rollup to every listing.
	Noise bool
	// NoStores drops the "notes" store from the session, leaving one with no
	// memory attached at all.
	NoStores bool
}

func newMemoryServer(t *testing.T, initial map[string]string, access string) *memoryServer {
	files := map[string]string{}
	for k, v := range initial {
		files[k] = v
	}
	return &memoryServer{t: t, Files: files, access: access}
}

func (m *memoryServer) Write(path, content string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Files[path] = content
}

func (m *memoryServer) Delete(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Files, path)
}

func (m *memoryServer) SetOnList(fn func()) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onList = fn
}

func (m *memoryServer) SetOnFetch(fn func(path string) int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onFetch = fn
}

func (m *memoryServer) SetUploadHook(fn func(ctx context.Context, path string)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.uploadHook = fn
}

func (m *memoryServer) ReceivedSnapshot() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]string(nil), m.Received...)
}

// ContentFetchesSnapshot returns the content fetches so far, sorted — they
// fan out, so arrival order is not deterministic.
func (m *memoryServer) ContentFetchesSnapshot() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return slices.Sorted(slices.Values(m.ContentFetches))
}

func (m *memoryServer) ClearContentFetches() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ContentFetches = nil
}

// memID round-trips a path through the URL-safe id the fake hands out.
func memID(path string) string { return "mem_" + hex.EncodeToString([]byte(path)) }

func memIDPath(id string) string {
	raw, _ := hex.DecodeString(strings.TrimPrefix(id, "mem_"))
	return string(raw)
}

// itemJSON renders one memory the way the API does: the basic view carries
// the sha but no content, the full view carries both.
func (m *memoryServer) itemJSON(path string, full bool) map[string]any {
	content := m.Files[path]
	item := map[string]any{
		"type": "memory", "id": memID(path), "path": "/" + path,
		"content_sha256":     shaHex(content),
		"content_size_bytes": len(content),
		"created_at":         "2026-05-11T12:00:00Z", "updated_at": "2026-05-11T12:00:00Z",
		"memory_store_id": "memstore_notes", "memory_version_id": "memver_1",
	}
	if full {
		item["content"] = content
	}
	return item
}

func isFullView(r *http.Request) bool { return r.URL.Query().Get("view") == "full" }

func (m *memoryServer) sessionJSON() string {
	var access any
	if m.access != "" {
		access = m.access
	}
	resources := []map[string]any{}
	if m.Broken {
		resources = append(resources, map[string]any{
			"type": "memory_store", "memory_store_id": "memstore_broken",
			"mount_path": nil, "name": "broken", "access": nil,
		})
	}
	if m.Noise {
		resources = append(resources, map[string]any{"type": "file", "file_id": "file_1"})
	}
	var mount any
	if m.MountPath != "" {
		mount = m.MountPath
	}
	if !m.NoStores {
		resources = append(resources, map[string]any{
			"type": "memory_store", "memory_store_id": "memstore_notes",
			"mount_path": mount, "name": "notes", "access": access,
		})
	}
	body, _ := json.Marshal(map[string]any{
		"id":        "s1",
		"agent":     map[string]any{"skills": []any{}},
		"resources": resources,
	})
	return string(body)
}

func statusJSON(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"type":"error","error":{"type":"api_error","message":"status %d"}}`, code)
}

// uploadTarget reports the memory path a create or update request is for;
// false for every other request. The body is drained and rewound first: a
// create names its path there, and net/http only notices a client that
// gave up once the body has been read to the end.
func uploadTarget(r *http.Request) (string, bool) {
	if r.Method != http.MethodPost {
		return "", false
	}
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(body))
	rest := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/v1/memory_stores/"), "/", 3)
	switch {
	case len(rest) == 2 && rest[1] == "memories":
		var in struct{ Path string }
		_ = json.Unmarshal(body, &in)
		return strings.TrimLeft(in.Path, "/"), true
	case len(rest) == 3 && rest[1] == "memories":
		id, _ := url.PathUnescape(rest[2])
		return memIDPath(strings.SplitN(id, "?", 2)[0]), true
	default:
		return "", false
	}
}

// handleMemories serves the /v1/memory_stores/... endpoints against the
// in-memory store.
func (m *memoryServer) handleMemories(w http.ResponseWriter, r *http.Request) {
	if rel, ok := uploadTarget(r); ok {
		m.mu.Lock()
		hook := m.uploadHook
		m.mu.Unlock()
		if hook != nil {
			hook(r.Context(), rel)
			if r.Context().Err() != nil {
				return
			}
		}
		if m.FailUploads != 0 {
			statusJSON(w, m.FailUploads)
			return
		}
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	rest := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/v1/memory_stores/"), "/", 3)
	switch {
	case rest[0] == "memstore_broken":
		// One memory is listed, then the pager explodes on the next page.
		if r.URL.Query().Get("page") == "" {
			half := map[string]any{
				"type": "memory", "id": "mem_half", "path": "/half.md",
				"content_sha256":     shaHex("partial"),
				"content_size_bytes": len("partial"),
				"created_at":         "2026-05-11T12:00:00Z", "updated_at": "2026-05-11T12:00:00Z",
				"memory_store_id": "memstore_broken", "memory_version_id": "memver_1",
			}
			if isFullView(r) {
				half["content"] = "partial"
			}
			data := []any{}
			if m.Noise {
				data = append(data, map[string]any{"type": "memory_prefix", "path": "/projects/"})
			}
			data = append(data, half)
			_ = json.NewEncoder(w).Encode(map[string]any{"data": data, "next_page": "boom"})
			return
		}
		statusJSON(w, 500)
	case len(rest) == 2 && rest[1] == "memories" && r.Method == http.MethodGet:
		items := make([]map[string]any, 0, len(m.Files)+1)
		if m.Noise {
			// Depth-limited listings roll directories up as prefixes; only
			// "memory" items carry content.
			items = append(items, map[string]any{"type": "memory_prefix", "path": "/projects/"})
		}
		full := isFullView(r)
		for path := range m.Files {
			items = append(items, m.itemJSON(path, full))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": items, "next_page": nil})
		if m.onList != nil {
			m.onList()
		}
	case len(rest) == 2 && rest[1] == "memories" && r.Method == http.MethodPost:
		var in struct{ Path, Content string }
		body, _ := io.ReadAll(r.Body)
		require.NoError(m.t, json.Unmarshal(body, &in))
		bare := strings.TrimLeft(in.Path, "/")
		m.Received = append(m.Received, created(bare, in.Content))
		if m.MaxContentBytes > 0 && len(in.Content) > m.MaxContentBytes {
			statusJSON(w, 400)
			return
		}
		m.Files[bare] = in.Content
		_ = json.NewEncoder(w).Encode(m.itemJSON(bare, true))
	case len(rest) == 3 && rest[1] == "memories":
		id, _ := url.PathUnescape(rest[2])
		id = strings.SplitN(id, "?", 2)[0]
		path := memIDPath(id)
		switch r.Method {
		case http.MethodGet: // single-memory content fetch
			m.ContentFetches = append(m.ContentFetches, path)
			if m.onFetch != nil {
				if code := m.onFetch(path); code != 0 {
					statusJSON(w, code)
					return
				}
			}
			if _, ok := m.Files[path]; !ok {
				statusJSON(w, 404)
				return
			}
			_ = json.NewEncoder(w).Encode(m.itemJSON(path, isFullView(r)))
		case http.MethodPost: // update
			if m.OnUpdate != nil {
				m.OnUpdate()
			}
			var in struct {
				Content      string
				Precondition struct {
					ContentSha256 string `json:"content_sha256"`
				}
			}
			body, _ := io.ReadAll(r.Body)
			require.NoError(m.t, json.Unmarshal(body, &in))
			m.Received = append(m.Received, "update "+path+" "+in.Content+" if="+in.Precondition.ContentSha256)
			if m.MaxContentBytes > 0 && len(in.Content) > m.MaxContentBytes {
				statusJSON(w, 400)
				return
			}
			if m.FailUpdates {
				statusJSON(w, 500)
				return
			}
			current, ok := m.Files[path]
			if !ok {
				statusJSON(w, 404)
				return
			}
			if shaHex(current) != in.Precondition.ContentSha256 {
				statusJSON(w, 409)
				return
			}
			m.Files[path] = in.Content
			_ = json.NewEncoder(w).Encode(m.itemJSON(path, true))
		case http.MethodDelete:
			expected := r.URL.Query().Get("expected_content_sha256")
			m.Received = append(m.Received, "delete "+path+" if="+expected)
			current, ok := m.Files[path]
			if !ok {
				statusJSON(w, 404)
				return
			}
			if expected != "" && shaHex(current) != expected {
				statusJSON(w, 409)
				return
			}
			delete(m.Files, path)
			_ = json.NewEncoder(w).Encode(map[string]any{"id": id, "type": "memory_deleted"})
		default:
			statusJSON(w, 405)
		}
	default:
		m.t.Errorf("unexpected memory_stores request: %s %s", r.Method, r.URL.Path)
		statusJSON(w, 404)
	}
}

func (m *memoryServer) serveHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/v1/memory_stores/"):
		m.handleMemories(w, r)
	case strings.HasPrefix(r.URL.Path, "/v1/sessions/"):
		m.mu.Lock()
		m.Retrieves = append(m.Retrieves, strings.TrimPrefix(r.URL.Path, "/v1/sessions/"))
		m.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(m.sessionJSON()))
	default:
		m.t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		statusJSON(w, 404)
	}
}

// fetchSession retrieves the fake server's session. Download takes that
// snapshot and never looks the session up itself, so every test fetches it the
// way the worker does — once, up front.
func fetchSession(t *testing.T, client anthropic.Client) *anthropic.BetaManagedAgentsSession {
	t.Helper()
	session, err := client.Beta.Sessions.Get(context.Background(), "s1", anthropic.BetaSessionGetParams{})
	require.NoError(t, err)
	return session
}

// mustStores is NewSessionMemoryStores for options a test knows are valid.
func mustStores(t *testing.T, client anthropic.Client, opts SessionMemoryStoresOptions) *SessionMemoryStores {
	t.Helper()
	stores, err := NewSessionMemoryStores(client, opts)
	require.NoError(t, err)
	return stores
}

// faults selects which disk verbs fail, and for which paths, and can act on
// the disk just before a write. The zero value does nothing; a test sets a
// field to break the verb it is exercising and clears it again to watch the
// retry. The disk events the baseline rules are written against — a write
// that never lands, a removal that cannot happen, a folder wiped between two
// writes — have no reliable fixture, so they are injected here.
type faults struct {
	put      func(rel string) bool
	remove   func(rel string) bool
	hashFile func(rel string) bool
	// onPut runs as each Put arrives, before it reaches the disk — a race
	// hook to land a change on disk between two of a pass's writes.
	onPut func(rel string)
}

// always breaks a verb for every path.
func always(string) bool { return true }

// open is the [SessionMemoryStores.openStore] seam: every store opens through
// a faultyStore reading these faults live.
func (f *faults) open(root string, opts ...filestore.Option) (storeFiles, error) {
	files, err := filestore.Open(root, opts...)
	if err != nil {
		return nil, err
	}
	return &faultyStore{storeFiles: files, faults: f}, nil
}

type faultyStore struct {
	storeFiles
	faults *faults
}

func (f *faultyStore) Put(rel string, data []byte, opts ...filestore.PutOption) error {
	if f.faults.onPut != nil {
		f.faults.onPut(rel)
	}
	if f.faults.put != nil && f.faults.put(rel) {
		return fmt.Errorf("injected write failure: %s", rel)
	}
	return f.storeFiles.Put(rel, data, opts...)
}

func (f *faultyStore) Remove(rel string) error {
	if f.faults.remove != nil && f.faults.remove(rel) {
		return fmt.Errorf("injected removal failure: %s", rel)
	}
	return f.storeFiles.Remove(rel)
}

func (f *faultyStore) HashFile(rel string) (string, error) {
	if f.faults.hashFile != nil && f.faults.hashFile(rel) {
		return "", fmt.Errorf("injected hash failure: %s", rel)
	}
	return f.storeFiles.HashFile(rel)
}

// captureLog is a debug-level text logger writing into buf, for tests that
// assert on log output.
func captureLog(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

// downloadedStores returns a SessionMemoryStores with the fake server's one
// store already downloaded to disk, plus the local dir and the server handle.
// disk, when non-nil, opens every store through a faultyStore.
func downloadedStores(t *testing.T, initial map[string]string, access string, logBuf *bytes.Buffer, disk *faults) (string, *memoryServer, *SessionMemoryStores) {
	t.Helper()
	server := newMemoryServer(t, initial, access)
	httpSrv := httptest.NewServer(http.HandlerFunc(server.serveHTTP))
	t.Cleanup(httpSrv.Close)
	client := anthropic.NewClient(option.WithBaseURL(httpSrv.URL), option.WithAPIKey("k"), option.WithMaxRetries(0))

	logger := silentLogger
	if logBuf != nil {
		logger = captureLog(logBuf)
	}
	workdir := t.TempDir()
	stores := mustStores(t, client, SessionMemoryStoresOptions{Workdir: workdir, Logger: logger})
	if disk != nil {
		stores.openStore = disk.open
	}
	require.NoError(t, stores.Download(context.Background(), fetchSession(t, client)))
	return filepath.Join(workdir, "memory", "notes"), server, stores
}

// testClock is a settable stand-in for the store's time source.
type testClock struct{ now time.Time }

func (c *testClock) Now() time.Time { return c.now }

// withClock swaps in a fake clock at time zero and returns a handle to it.
func withClock(stores *SessionMemoryStores) *testClock {
	c := &testClock{now: time.Unix(0, 0)}
	stores.setClock(c.Now)
	return c
}

// runSync drives one reconcile pass without clock bookkeeping.
//
// The public cadence path (SyncIfDue plus the clock) is exercised by the
// worker tests; unit tests call the pass directly.
func runSync(ctx context.Context, s *SessionMemoryStores) {
	s.syncAll(ctx, false)
}

func writeLocal(t *testing.T, dir, name, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
}

func readLocal(t *testing.T, dir, name string) string {
	t.Helper()
	got, err := os.ReadFile(filepath.Join(dir, name))
	require.NoError(t, err)
	return string(got)
}

func TestMemorySync_EditsFlowBothWaysAndTheServerWinsConflicts(t *testing.T) {
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t,
		map[string]string{"push.md": "v1", "pull.md": "v1", "both.md": "v1", "gone.md": "v1"}, "", &logBuf, nil)
	ctx := context.Background()

	writeLocal(t, local, "push.md", "v2 from agent")
	writeLocal(t, local, "new.md", "new from agent")
	server.Write("pull.md", "v2 from server")
	writeLocal(t, local, "both.md", "v2 from agent")
	server.Write("both.md", "v2 from server")
	server.Delete("gone.md")

	runSync(ctx, stores)

	// Each side has the other's change; the server won the conflict; the
	// server's deletion took the local copy with it.
	require.Equal(t, "v2 from agent", server.Files["push.md"])
	require.Equal(t, "new from agent", server.Files["new.md"])
	require.Equal(t, "v2 from server", readLocal(t, local, "pull.md"))
	require.Equal(t, "v2 from server", readLocal(t, local, "both.md"))
	require.NoFileExists(t, filepath.Join(local, "gone.md"))
	require.NotContains(t, server.Files, "gone.md")
	require.Contains(t, logBuf.String(), "changed both locally and remotely")

	// The push was guarded: overwrite only if the server still had v1.
	sent := []string{created("new.md", "new from agent"), updated("push.md", "v2 from agent", "v1")}
	require.ElementsMatch(t, sent, server.ReceivedSnapshot())

	// Settled: a second sync sends nothing more.
	runSync(ctx, stores)
	require.Len(t, server.ReceivedSnapshot(), len(sent))
}

func TestMemorySync_LocalDeletionsReachTheServerUnlessItChanged(t *testing.T) {
	// One memory's download write fails, so it is never on disk to begin with
	// — that absence is not a deletion.
	disk := &faults{put: func(rel string) bool { return rel == "never.md" }}
	// keep.md survives on disk so the pass is not an all-files-gone wipe.
	local, server, stores := downloadedStores(t,
		map[string]string{"drop.md": "v1", "raced.md": "v1", "never.md": "v1", "keep.md": "v1"}, "", nil, disk)
	disk.put = nil
	clock := withClock(stores)
	ctx := context.Background()
	require.NoFileExists(t, filepath.Join(local, "never.md"))

	require.NoError(t, os.Remove(filepath.Join(local, "drop.md")))
	require.NoError(t, os.Remove(filepath.Join(local, "raced.md")))
	server.Write("raced.md", "v2 from server")

	runSync(ctx, stores)
	// The first pass only records the deletion; a corroborating pass
	// after the window sends it.
	require.Contains(t, server.Files, "drop.md")
	clock.now = clock.now.Add(DeleteCorroborationWindow + time.Second)
	runSync(ctx, stores)

	// Only the un-raced deletion went out, guarded by the last-synced content.
	sent := []string{deletedReq("drop.md", "v1")}
	require.Equal(t, sent, server.ReceivedSnapshot())
	require.NotContains(t, server.Files, "drop.md")
	// The raced deletion lost — the server's edit came back to disk instead.
	require.Equal(t, "v2 from server", readLocal(t, local, "raced.md"))
	// The never-downloaded file was pulled, not deleted.
	require.Equal(t, "v1", readLocal(t, local, "never.md"))
	require.Equal(t, "v1", server.Files["never.md"])

	runSync(ctx, stores)
	require.Equal(t, sent, server.ReceivedSnapshot())
}

func TestMemorySync_ReadOnlyStoresAreNeverWritten(t *testing.T) {
	local, server, stores := downloadedStores(t,
		map[string]string{"edit.md": "v1", "drop.md": "v1", "pull.md": "v1", "vanish.md": "v1"}, "read_only", nil, nil)
	ctx := context.Background()

	writeLocal(t, local, "edit.md", "v2 from agent")
	require.NoError(t, os.Remove(filepath.Join(local, "drop.md")))
	writeLocal(t, local, "new.md", "new from agent")
	server.Write("pull.md", "v2 from server")
	server.Delete("vanish.md")

	runSync(ctx, stores)

	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"edit.md": "v1", "drop.md": "v1", "pull.md": "v2 from server"}, server.Files)
	// Pulls still happen...
	require.Equal(t, "v2 from server", readLocal(t, local, "pull.md"))
	// ...including a removal, which is a pull and not a push: the server no
	// longer holds vanish.md, so neither does the folder.
	require.NoFileExists(t, filepath.Join(local, "vanish.md"))
	// And the store's root is reported for the write/edit tool refusal, as
	// well as for the file tools' reach.
	require.Equal(t, []string{local}, stores.ReadOnlyRoots())
	require.Equal(t, stores.ReadOnlyRoots(), stores.Roots())
}

func TestMemorySync_PathsThatEscapeTheStoreAreRefused(t *testing.T) {
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t,
		map[string]string{"../evil.md": "boom", "ok.md": "fine"}, "", &logBuf, nil)

	runSync(context.Background(), stores)

	require.Equal(t, "fine", readLocal(t, local, "ok.md"))
	// Nothing named evil.md landed anywhere under the workdir.
	workdir := filepath.Dir(filepath.Dir(local))
	require.NoError(t, filepath.Walk(workdir, func(path string, _ os.FileInfo, err error) error {
		require.NotEqual(t, "evil.md", filepath.Base(path))
		return err
	}))
	require.Empty(t, server.ReceivedSnapshot())
	require.Contains(t, logBuf.String(), "escapes")
}

func TestMemorySync_LifecycleCadenceAndDispose(t *testing.T) {
	// The worker's call sequence: download → SyncIfDue per call → Finish →
	// Dispose. A fake clock makes the cadence observable.
	server := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	httpSrv := httptest.NewServer(http.HandlerFunc(server.serveHTTP))
	t.Cleanup(httpSrv.Close)
	client := anthropic.NewClient(option.WithBaseURL(httpSrv.URL), option.WithAPIKey("k"), option.WithMaxRetries(0))

	now := time.Unix(0, 0)
	workdir := t.TempDir()
	stores := mustStores(t, client, SessionMemoryStoresOptions{
		Workdir:      workdir,
		SyncInterval: 300 * time.Second,
		Logger:       silentLogger,
	})
	stores.setClock(func() time.Time { return now })
	ctx := context.Background()
	require.NoError(t, stores.Download(ctx, fetchSession(t, client)))
	local := filepath.Join(workdir, "memory", "notes")

	writeLocal(t, local, "note.md", "edit 1")
	stores.SyncIfDue(ctx)
	require.Empty(t, server.ReceivedSnapshot())

	// One long gap fires exactly once, and not again until another interval.
	now = time.Unix(2000, 0)
	stores.SyncIfDue(ctx)
	stores.SyncIfDue(ctx)
	require.Equal(t, []string{updated("note.md", "edit 1", "v1")}, server.ReceivedSnapshot())
	now = time.Unix(2299, 0)
	stores.SyncIfDue(ctx)
	require.Len(t, server.ReceivedSnapshot(), 1)

	// The last sync is Finish: unconditional, once.
	writeLocal(t, local, "note.md", "edit 2")
	stores.Finish(ctx)
	received := server.ReceivedSnapshot()
	require.Equal(t, updated("note.md", "edit 2", "edit 1"), received[len(received)-1])

	// Dispose removes the store dir the download created.
	require.DirExists(t, local)
	stores.Dispose()
	require.NoDirExists(t, local)
}

func TestMemorySync_TheSyncIntervalHasADefaultAndAFloor(t *testing.T) {
	client := anthropic.NewClient(option.WithAPIKey("k"))
	build := func(d time.Duration) (*SessionMemoryStores, error) {
		return NewSessionMemoryStores(client, SessionMemoryStoresOptions{
			Workdir:      t.TempDir(),
			SyncInterval: d,
			Logger:       silentLogger,
		})
	}
	for _, tc := range []struct{ in, want time.Duration }{
		{0, DefaultMemorySyncInterval},
		{MinMemorySyncInterval, MinMemorySyncInterval},
		{time.Minute, time.Minute},
		{-1, -1}, // negative passes through; the worker gates on it
	} {
		stores, err := build(tc.in)
		require.NoError(t, err, tc.in)
		require.Equal(t, tc.want, stores.syncInterval, tc.in)
	}
	_, err := build(time.Second)
	require.ErrorIs(t, err, ErrMemorySyncIntervalTooShort)
}

func TestMemorySync_ABinaryLocalFileNeverBlocksTheOtherFiles(t *testing.T) {
	// The store is utf-8-restricted, so a binary file the agent drops in the
	// folder is refused at read — warned once and skipped until it changes —
	// while every other file keeps syncing.
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"note.md": "v1"}, "", &logBuf, nil)
	ctx := context.Background()

	require.NoError(t, os.WriteFile(filepath.Join(local, "blob.bin"), []byte{0xff, 0xfe, 0x00, ' ', 'x'}, 0o600))
	writeLocal(t, local, "note.md", "v2")

	runSync(ctx, stores)
	runSync(ctx, stores)

	require.Equal(t, []string{updated("note.md", "v2", "v1")}, server.ReceivedSnapshot())
	require.NotContains(t, server.Files, "blob.bin")
	require.Equal(t, 1, strings.Count(logBuf.String(), "stays un-synced until its content changes"))
	require.Contains(t, logBuf.String(), "path=blob.bin")
}

func TestMemorySync_AnOversizedFileIsSkippedUntilItShrinks(t *testing.T) {
	// A file the server refuses as too large is warned once and skipped —
	// not retried while unchanged, never blocking the other files — and
	// syncs as soon as an edit brings it under the cap.
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"note.md": "v1"}, "", &logBuf, nil)
	ctx := context.Background()
	server.MaxContentBytes = 64

	writeLocal(t, local, "big.md", strings.Repeat("x", 100))
	writeLocal(t, local, "note.md", "v2")

	runSync(ctx, stores)
	runSync(ctx, stores)

	// The small edit landed; the big file was attempted exactly once.
	received := server.ReceivedSnapshot()
	require.Contains(t, received, updated("note.md", "v2", "v1"))
	creates := 0
	for _, r := range received {
		if strings.HasPrefix(r, "create ") {
			creates++
		}
	}
	require.Equal(t, 1, creates)
	require.NotContains(t, server.Files, "big.md")
	require.Equal(t, 1, strings.Count(logBuf.String(), "stays un-synced until its content changes"))
	require.Contains(t, logBuf.String(), "path=big.md")

	writeLocal(t, local, "big.md", "small now")
	runSync(ctx, stores)
	require.Equal(t, "small now", server.Files["big.md"])
}

func TestMemorySync_ARefusalThatSaysNothingAboutTheContentIsRetriedNextSync(t *testing.T) {
	// Anything but a content refusal is retried by the next sync.
	for _, status := range []int{401, 403, 404} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			var logBuf bytes.Buffer
			local, server, stores := downloadedStores(t, map[string]string{"note.md": "v1"}, "", &logBuf, nil)
			ctx := context.Background()
			writeLocal(t, local, "note.md", "v2")
			writeLocal(t, local, "new.md", "fresh")

			server.FailUploads = status
			runSync(ctx, stores)
			server.FailUploads = 0
			runSync(ctx, stores)

			require.Equal(t, map[string]string{"note.md": "v2", "new.md": "fresh"}, server.Files)
			require.NotContains(t, logBuf.String(), "stays un-synced")
		})
	}
}

func TestMemorySync_AnUploadThatLosesTheRacePullsTheWinnerNextSync(t *testing.T) {
	// The server changes between this sync's list and its update: the
	// content_sha256 precondition rejects the push (409), and the next sync
	// applies the conflict rule — the winning remote edit is pulled, the
	// local edit is never re-pushed.
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"note.md": "v1"}, "", &logBuf, nil)
	ctx := context.Background()

	writeLocal(t, local, "note.md", "local edit")
	server.OnUpdate = func() {
		// Race: the remote changes after this sync's list, before its
		// update. (Called with the server mutex held.)
		server.Files["note.md"] = "raced remote edit"
	}
	runSync(ctx, stores)
	server.OnUpdate = nil

	// The precondition rejected the push (one guarded attempt, no retry) and
	// the server kept its raced edit.
	require.Equal(t, []string{updated("note.md", "local edit", "v1")}, server.ReceivedSnapshot())
	require.Equal(t, "raced remote edit", server.Files["note.md"])
	require.Contains(t, logBuf.String(), "changed both locally and remotely")
	// The warning says what actually happened: the push was dropped and the
	// local file is the stale one. Nothing "kept the remote version" here —
	// disk still holds the losing edit until the next sync overwrites it.
	require.Contains(t, logBuf.String(), "the upload was refused and the local edit loses")
	require.Equal(t, "local edit", readLocal(t, local, "note.md"))

	// The next sync pulls the winner instead of re-pushing.
	runSync(ctx, stores)
	require.Equal(t, "raced remote edit", readLocal(t, local, "note.md"))
	require.Len(t, server.ReceivedSnapshot(), 1)
}

// storesAt wires a SessionMemoryStores to server over workdir without
// downloading — for the tests that assert on Download's own refusals.
func storesAt(t *testing.T, server *memoryServer, workdir string) *SessionMemoryStores {
	t.Helper()
	httpSrv := httptest.NewServer(http.HandlerFunc(server.serveHTTP))
	t.Cleanup(httpSrv.Close)
	client := anthropic.NewClient(option.WithBaseURL(httpSrv.URL), option.WithAPIKey("k"), option.WithMaxRetries(0))
	return mustStores(t, client, SessionMemoryStoresOptions{Workdir: workdir, Logger: silentLogger})
}

func TestMemorySync_ABrokenStoreFailsTheWholeDownload(t *testing.T) {
	// The first attached store's listing explodes mid-download: its half-written
	// folder self-destructs and the whole download fails. A session served
	// without a folder its system prompt names would run with amnesia, so the
	// worker turns this into a failed work item rather than a warning.
	server := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	server.Broken = true
	workdir := t.TempDir()
	stores := storesAt(t, server, workdir)

	err := stores.Download(context.Background(), fetchSession(t, stores.client))
	var memErr *SessionMemoryError
	require.ErrorAs(t, err, &memErr)
	require.Equal(t, "memstore_broken", memErr.MemoryStoreID)

	// The broken store got one file down before exploding; nothing survives.
	// The healthy store behind it was never reached, and nothing is tracked.
	require.NoDirExists(t, filepath.Join(workdir, "memory", "broken"))
	require.NoDirExists(t, filepath.Join(workdir, "memory", "notes"))
	require.Empty(t, stores.ReadOnlyRoots())
	require.Empty(t, server.ReceivedSnapshot())
}

func TestMemorySync_ALegalMountPathIsHonored(t *testing.T) {
	// mount_path is the location the agent's system prompt names. A clean
	// absolute one is honored verbatim, and Dispose removes it when the
	// download created it.
	tmp := t.TempDir()
	mount := filepath.Join(tmp, "mnt", "notes")
	server := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	server.MountPath = mount
	stores := storesAt(t, server, filepath.Join(tmp, "wd"))

	require.NoError(t, stores.Download(context.Background(), fetchSession(t, stores.client)))
	require.Equal(t, "v1", readLocal(t, mount, "note.md"))
	require.Equal(t, []string{mount}, stores.Roots())
	stores.Dispose()
	require.NoDirExists(t, mount)
}

func TestMemorySync_AnUncleanMountPathFailsTheDownload(t *testing.T) {
	// A mount_path we cannot use verbatim fails the session rather than quietly
	// relocating the store: the agent would read an empty folder at the path it
	// was told about, and write notes somewhere the next session never looks.
	for _, tc := range []struct {
		name  string
		mount string
		// "" means the resource carries no mount_path at all — that one is
		// allowed, and falls back to the workdir, because nothing pointed the
		// agent anywhere.
		ok bool
	}{
		{"escapes with ..", "/mnt/../../etc/cron.d", false},
		{"relative, with a separator", "relative/notes", false},
		{"relative, a bare name", "notes", false},
		{"no mount path at all", "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tmp := t.TempDir()
			workdir := filepath.Join(tmp, "wd")
			server := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
			server.MountPath = tc.mount
			stores := storesAt(t, server, workdir)

			err := stores.Download(context.Background(), fetchSession(t, stores.client))
			if tc.ok {
				require.NoError(t, err)
				require.Equal(t, "v1", readLocal(t, filepath.Join(workdir, "memory", "notes"), "note.md"))
				return
			}

			var memErr *SessionMemoryError
			require.ErrorAs(t, err, &memErr)
			require.Equal(t, "memstore_notes", memErr.MemoryStoreID)
			require.ErrorContains(t, err, "not a clean absolute path")

			// Nothing was written anywhere, nothing was sent, nothing is tracked.
			require.NoDirExists(t, workdir)
			require.Empty(t, server.ReceivedSnapshot())
			require.Empty(t, stores.ReadOnlyRoots())
		})
	}
}

func TestMemorySync_AnUnwritableMountPathFailsTheDownload(t *testing.T) {
	// A mount the host cannot provide is an environment failure: the download
	// fails loudly (failing the work item) instead of degrading to a session
	// with silently-empty memory.
	if os.Geteuid() == 0 {
		t.Skip("directory modes do not bind root")
	}
	tmp := t.TempDir()
	parent := filepath.Join(tmp, "mnt")
	require.NoError(t, os.Mkdir(parent, 0o700))
	require.NoError(t, os.Chmod(parent, 0o500)) // the worker may not create the store's folder here
	t.Cleanup(func() { _ = os.Chmod(parent, 0o700) })

	server := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	server.MountPath = filepath.Join(parent, "notes")
	stores := storesAt(t, server, filepath.Join(tmp, "wd"))

	err := stores.Download(context.Background(), fetchSession(t, stores.client))
	var memErr *SessionMemoryError
	require.ErrorAs(t, err, &memErr)
	require.Equal(t, "memstore_notes", memErr.MemoryStoreID)
	require.ErrorContains(t, err, "writable")
	require.Empty(t, server.ReceivedSnapshot())
}

func TestMemorySync_AStoreDirectoryThatAlreadyExistsIsRefused(t *testing.T) {
	// Nothing but a dead run leaves a folder at the store's path. Syncing into
	// it would fold that run's debris into the customer's store, so refuse —
	// and leave the folder exactly as it was found.
	workdir := t.TempDir()
	stale := filepath.Join(workdir, "memory", "notes")
	require.NoError(t, os.MkdirAll(stale, 0o755))
	writeLocal(t, stale, "debris.md", "left by another session")

	server := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	stores := storesAt(t, server, workdir)

	err := stores.Download(context.Background(), fetchSession(t, stores.client))
	var memErr *SessionMemoryError
	require.ErrorAs(t, err, &memErr)
	require.Equal(t, "memstore_notes", memErr.MemoryStoreID)
	require.ErrorContains(t, err, "already exists")

	// Refused, not adopted and not deleted — not even by Dispose.
	require.Equal(t, "left by another session", readLocal(t, stale, "debris.md"))
	require.NoFileExists(t, filepath.Join(stale, "note.md"))
	require.Empty(t, server.ReceivedSnapshot())
	require.Empty(t, stores.ReadOnlyRoots())
	stores.Dispose()
	require.DirExists(t, stale)

	// A sibling directory in the same parent is none of our business.
	other := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(other, "memory", "unrelated"), 0o755))
	server2 := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	stores2 := storesAt(t, server2, other)

	require.NoError(t, stores2.Download(context.Background(), fetchSession(t, stores2.client)))
	require.Equal(t, "v1", readLocal(t, filepath.Join(other, "memory", "notes"), "note.md"))
	require.DirExists(t, filepath.Join(other, "memory", "unrelated"))
}

func TestMemorySync_PrefixItemsAndForeignResourcesLeaveNoTrace(t *testing.T) {
	// Only memory_store resources are downloaded and only "memory" items
	// carry content — a "file" resource and memory_prefix rollups pass
	// through without touching the disk or the server.
	server := newMemoryServer(t, map[string]string{"ok.md": "fine"}, "")
	server.Noise = true
	workdir := t.TempDir()
	stores := storesAt(t, server, workdir)
	ctx := context.Background()
	require.NoError(t, stores.Download(ctx, fetchSession(t, stores.client)))
	local := filepath.Join(workdir, "memory", "notes")

	require.Equal(t, "fine", readLocal(t, local, "ok.md"))
	entries, err := os.ReadDir(filepath.Join(workdir, "memory"))
	require.NoError(t, err)
	require.Len(t, entries, 1)
	require.Equal(t, "notes", entries[0].Name())
	require.NoDirExists(t, filepath.Join(local, "projects"))

	runSync(ctx, stores)
	require.Empty(t, server.ReceivedSnapshot())
}

// ---- the sync matrix -------------------------------------------------------

// Matrix cells: keep means the side did nothing, gone means it deleted the
// path. Any other value is the content that side holds now. The sentinels
// cannot collide with real content.
const (
	keep = "\x00keep"
	gone = "\x00gone"
)

func TestMemorySync_TheSyncMatrix(t *testing.T) {
	// Every combination of (local moved?, remote moved?) for a path the store
	// downloaded at "v1". local and remote say what happened since; disk and
	// store are the state the two sides must agree on afterwards, and sent is
	// exactly what went over the wire.
	for _, tc := range []struct {
		name          string
		local, remote string
		disk, store   string
		sent          []string
	}{
		{"untouched", keep, keep, "v1", "v1", nil},
		{"local edit", "v2", keep, "v2", "v2", []string{updated("f.md", "v2", "v1")}},
		{"remote edit", keep, "v2", "v2", "v2", nil},
		{"both edit, the server wins", "v2 local", "v2 remote", "v2 remote", "v2 remote", nil},
		{"local delete", gone, keep, gone, gone, []string{deletedReq("f.md", "v1")}},
		{"remote delete", keep, gone, gone, gone, nil},
		{"both delete", gone, gone, gone, gone, nil},
		{"local delete, remote edit", gone, "v2", "v2", "v2", nil},
		{"local edit, remote delete", "v2", gone, "v2", "v2", []string{created("f.md", "v2")}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			local, server, stores := downloadedStores(t, map[string]string{"f.md": "v1"}, "", nil, nil)
			clock := withClock(stores)
			ctx := context.Background()

			switch tc.local {
			case keep:
			case gone:
				require.NoError(t, os.Remove(filepath.Join(local, "f.md")))
			default:
				writeLocal(t, local, "f.md", tc.local)
			}
			switch tc.remote {
			case keep:
			case gone:
				server.Delete("f.md")
			default:
				server.Write("f.md", tc.remote)
			}

			settled := func() {
				t.Helper()
				if tc.disk == gone {
					require.NoFileExists(t, filepath.Join(local, "f.md"))
				} else {
					require.Equal(t, tc.disk, readLocal(t, local, "f.md"))
				}
				if tc.store == gone {
					require.NotContains(t, server.Files, "f.md")
				} else {
					require.Equal(t, tc.store, server.Files["f.md"])
				}
				require.Equal(t, tc.sent, server.ReceivedSnapshot())
			}

			// Two passes: a local deletion is only sent after its
			// corroboration window.
			runSync(ctx, stores)
			clock.now = clock.now.Add(DeleteCorroborationWindow + time.Second)
			runSync(ctx, stores)
			settled()

			// Settled: a second pass sends nothing more and changes nothing.
			runSync(ctx, stores)
			settled()
		})
	}
}

// ---- the failure column ----------------------------------------------------
//
// The same paths, one operation broken. A failed operation must leave disk and
// server in a state the next sync can still reconcile — never one where the
// retry does the opposite of what was asked.

func TestMemorySync_AFailedUploadKeepsTheLocalEditAndRetries(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"f.md": "v1"}, "", nil, nil)
	ctx := context.Background()
	writeLocal(t, local, "f.md", "v2 from agent")

	server.FailUpdates = true
	runSync(ctx, stores)

	// The agent's edit survives; the server is untouched.
	require.Equal(t, "v2 from agent", readLocal(t, local, "f.md"))
	require.Equal(t, "v1", server.Files["f.md"])

	// The baseline still says v1, so the next (working) sync retries the push
	// with the same precondition.
	server.FailUpdates = false
	runSync(ctx, stores)
	require.Equal(t, "v2 from agent", server.Files["f.md"])
	push := updated("f.md", "v2 from agent", "v1")
	require.Equal(t, []string{push, push}, server.ReceivedSnapshot())
}

func TestMemorySync_AFailedPullLeavesTheOldFileAndRetries(t *testing.T) {
	disk := &faults{}
	local, server, stores := downloadedStores(t, map[string]string{"f.md": "v1"}, "", nil, disk)
	ctx := context.Background()
	server.Write("f.md", "v2 from server")

	disk.put = always
	runSync(ctx, stores)
	disk.put = nil

	// The stale file is still there and the baseline never advanced...
	require.Equal(t, "v1", readLocal(t, local, "f.md"))
	require.Empty(t, server.ReceivedSnapshot())

	// ...so the next sync pulls it, rather than mistaking v1 for a local edit
	// and pushing it back over the server's v2.
	runSync(ctx, stores)
	require.Equal(t, "v2 from server", readLocal(t, local, "f.md"))
	require.Empty(t, server.ReceivedSnapshot())
}

func TestMemorySync_AFailedLocalRemovalNeverReUploadsTheDeletedMemory(t *testing.T) {
	// The server deleted the memory; removing our copy fails. The baseline must
	// stay put — dropping it would read the leftover file as a new local file
	// and create the memory the server just deleted.
	var logBuf bytes.Buffer
	disk := &faults{}
	local, server, stores := downloadedStores(t, map[string]string{"f.md": "v1"}, "", &logBuf, disk)
	ctx := context.Background()
	server.Delete("f.md")

	disk.remove = always
	runSync(ctx, stores)
	disk.remove = nil

	require.Equal(t, "v1", readLocal(t, local, "f.md"))
	require.Empty(t, server.ReceivedSnapshot())
	require.NotContains(t, server.Files, "f.md")
	require.Contains(t, logBuf.String(), "failed to remove memory deleted remotely")

	// The retry deletes it locally; it is never pushed back.
	runSync(ctx, stores)
	require.NoFileExists(t, filepath.Join(local, "f.md"))
	require.Empty(t, server.ReceivedSnapshot())
	require.NotContains(t, server.Files, "f.md")
}

// twoStoreSession builds a session attaching two memory stores at the given
// mount paths, the way the control plane would send it.
func twoStoreSession(t *testing.T, aRoot, bRoot string) *anthropic.BetaManagedAgentsSession {
	t.Helper()
	body := fmt.Sprintf(`{"resources":[
		{"type":"memory_store","memory_store_id":"store_a","mount_path":%q,"name":"a"},
		{"type":"memory_store","memory_store_id":"store_b","mount_path":%q,"name":"b"}
	]}`, aRoot, bRoot)
	var session anthropic.BetaManagedAgentsSession
	require.NoError(t, json.Unmarshal([]byte(body), &session))
	return &session
}

func TestMemorySync_TwoStoresDownloadToTheirOwnFolders(t *testing.T) {
	wd := t.TempDir()
	server := newMemoryServer(t, map[string]string{"note.md": "v1"}, "")
	httpSrv := httptest.NewServer(http.HandlerFunc(server.serveHTTP))
	t.Cleanup(httpSrv.Close)
	client := anthropic.NewClient(option.WithBaseURL(httpSrv.URL), option.WithAPIKey("k"), option.WithMaxRetries(0))
	stores := mustStores(t, client, SessionMemoryStoresOptions{Workdir: wd, Logger: silentLogger})

	a, b := filepath.Join(wd, "m", "a"), filepath.Join(wd, "m", "b")
	require.NoError(t, stores.Download(context.Background(), twoStoreSession(t, a, b)))
	require.Equal(t, "v1", readLocal(t, a, "note.md"))
	require.Equal(t, "v1", readLocal(t, b, "note.md"))
}

func TestMemorySync_StoresSyncInParallelAndOneFailureSparesTheRest(t *testing.T) {
	// Store A's listing waits for store B's listing to start — only a
	// parallel sync satisfies that — then fails. A serial sync times the
	// wait out instead (a different error message), and store B's edit must
	// land either way.
	wd := t.TempDir()
	slow := newMemoryServer(t, map[string]string{"a.md": "v1"}, "")
	fast := newMemoryServer(t, map[string]string{"b.md": "v1"}, "")
	var syncing atomic.Bool
	fastStarted := make(chan struct{})
	var once sync.Once
	fail := func(w http.ResponseWriter, msg string) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, `{"error":{"type":"api_error","message":%q}}`, msg)
	}
	httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isList := r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/memories")
		if strings.Contains(r.URL.Path, "store_a") {
			if syncing.Load() && isList {
				select {
				case <-fastStarted:
					fail(w, "slow store exploded")
				case <-time.After(10 * time.Second):
					fail(w, "fast store never started")
				}
				return
			}
			slow.serveHTTP(w, r)
			return
		}
		if syncing.Load() && isList {
			once.Do(func() { close(fastStarted) })
		}
		fast.serveHTTP(w, r)
	}))
	t.Cleanup(httpSrv.Close)
	client := anthropic.NewClient(option.WithBaseURL(httpSrv.URL), option.WithAPIKey("k"), option.WithMaxRetries(0))

	var logBuf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuf, nil))
	stores := mustStores(t, client, SessionMemoryStoresOptions{Workdir: wd, Logger: logger})
	aRoot, bRoot := filepath.Join(wd, "m", "a"), filepath.Join(wd, "m", "b")
	ctx := context.Background()
	require.NoError(t, stores.Download(ctx, twoStoreSession(t, aRoot, bRoot)))

	syncing.Store(true)
	writeLocal(t, bRoot, "b.md", "v2 from agent")
	runSync(ctx, stores)

	// Store A failed *after* seeing store B's listing start — proof the two
	// ran concurrently — and store B still pushed.
	require.Contains(t, logBuf.String(), "slow store exploded")
	require.Equal(t, []string{updated("b.md", "v2 from agent", "v1")}, fast.ReceivedSnapshot())
	// The failed store's baseline was kept, so nothing is re-uploaded once it heals.
	require.Equal(t, "v1", readLocal(t, aRoot, "a.md"))
}

// ---- folder-trust gating ---------------------------------------------------
//
// A missing, emptied, or swapped store folder is never trusted as a basis
// for deletions or uploads.
//
// Download stamps MarkerPath (holding the store's id) into the folder. A
// sync pass that finds the folder destroyed re-downloads it and pushes
// nothing; a sync that finds files without a matching marker leaves the
// folder as found and syncs nothing at all.

const (
	folderTrustStoreID = "memstore_notes"
	folderTrustMarker  = "version 1\n" + folderTrustStoreID
)

func TestFolderTrust_MarkerBytesArePinned(t *testing.T) {
	// A format change breaks CI here rather than as a marker mismatch at runtime.
	// printf 'version 1\nmemstore_test' | sha256sum
	require.Equal(t, "b8c0fbafe3e9ea979f0fbe5d984ba8834e7a2771d859462b3f7aa556fdb55a67",
		markerSHA("memstore_test"))
}

func TestFolderTrust_ADeletedFolderIssuesNoDeletesAndIsRebuilt(t *testing.T) {
	// rm -rf of the whole folder between syncs: the next sync must not read
	// the emptiness as "the agent deleted everything" and wipe the server.
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1", "b.md": "v1"}, "", &logBuf, nil)
	clock := withClock(stores)
	ctx := context.Background()
	require.NoError(t, os.RemoveAll(local))

	runSync(ctx, stores)

	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"a.md": "v1", "b.md": "v1"}, server.Files)
	require.Equal(t, "v1", readLocal(t, local, "a.md"))
	require.Equal(t, "v1", readLocal(t, local, "b.md"))
	require.Equal(t, folderTrustMarker, readLocal(t, local, MarkerPath))
	require.Contains(t, logBuf.String(), "the folder or its marker is gone; re-downloading")

	// The rebuilt folder syncs normally: an edit and a deletion both
	// propagate (the deletion after its waiting window).
	writeLocal(t, local, "a.md", "v2")
	require.NoError(t, os.Remove(filepath.Join(local, "b.md")))
	runSync(ctx, stores)
	clock.now = clock.now.Add(DeleteCorroborationWindow + time.Second)
	runSync(ctx, stores)
	require.Equal(t, []string{updated("a.md", "v2", "v1"), deletedReq("b.md", "v1")}, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"a.md": "v2"}, server.Files)
}

func TestFolderTrust_AnEmptiedFolderIssuesNoDeletesAndIsRepopulated(t *testing.T) {
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1", "b.md": "v1"}, "", &logBuf, nil)
	clock := withClock(stores)
	ctx := context.Background()
	entries, _ := os.ReadDir(local)
	for _, e := range entries {
		require.NoError(t, os.Remove(filepath.Join(local, e.Name())))
	}

	runSync(ctx, stores)

	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"a.md": "v1", "b.md": "v1"}, server.Files)
	require.Equal(t, "v1", readLocal(t, local, "a.md"))
	require.Equal(t, folderTrustMarker, readLocal(t, local, MarkerPath))
	require.Contains(t, logBuf.String(), "the folder or its marker is gone; re-downloading")

	require.NoError(t, os.Remove(filepath.Join(local, "a.md")))
	runSync(ctx, stores)
	clock.now = clock.now.Add(DeleteCorroborationWindow + time.Second)
	runSync(ctx, stores)
	require.Equal(t, []string{deletedReq("a.md", "v1")}, server.ReceivedSnapshot())
}

func TestFolderTrust_AnEmptiedFolderWithTheMarkerLeftBehindIssuesNoDeletes(t *testing.T) {
	// rm <mount>/*: shell globs skip dotfiles, so the marker survives the
	// wipe — an intact marker alone must not clear a folder whose every
	// memory file vanished at once.
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1", "b.md": "v1"}, "", &logBuf, nil)
	clock := withClock(stores)
	ctx := context.Background()
	entries, _ := os.ReadDir(local)
	for _, e := range entries {
		if e.Name() != MarkerPath {
			require.NoError(t, os.Remove(filepath.Join(local, e.Name())))
		}
	}

	runSync(ctx, stores)

	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"a.md": "v1", "b.md": "v1"}, server.Files)
	require.Equal(t, "v1", readLocal(t, local, "a.md"))
	require.Contains(t, logBuf.String(), "every memory file is gone at once; re-downloading")

	// With files present again, an ordinary single deletion still
	// propagates (after its waiting window).
	require.NoError(t, os.Remove(filepath.Join(local, "a.md")))
	runSync(ctx, stores)
	clock.now = clock.now.Add(DeleteCorroborationWindow + time.Second)
	runSync(ctx, stores)
	require.Equal(t, []string{deletedReq("a.md", "v1")}, server.ReceivedSnapshot())
}

func TestFolderTrust_ASwappedFolderNeitherDeletesNorUploadsEver(t *testing.T) {
	// A folder carrying another store's marker is foreign: the pass must not
	// adopt it — no deletes, no uploads, no restamp — on this pass or any
	// later one, or the foreign files leak into the customer's store one
	// pass later.
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1"}, "", &logBuf, nil)
	ctx := context.Background()
	require.NoError(t, os.RemoveAll(local))
	require.NoError(t, os.MkdirAll(local, 0o755))
	writeLocal(t, local, MarkerPath, "memstore_other")
	writeLocal(t, local, "foreign.md", "someone else's notes")

	runSync(ctx, stores)
	runSync(ctx, stores)

	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"a.md": "v1"}, server.Files)
	require.Equal(t, "memstore_other", readLocal(t, local, MarkerPath))
	require.Equal(t, "someone else's notes", readLocal(t, local, "foreign.md"))
	require.NoFileExists(t, filepath.Join(local, "a.md"))
	require.Contains(t, logBuf.String(), "the marker file does not match this store; leaving the memory store folder as found")
}

func TestFolderTrust_FilesWithoutAMarkerAreNeitherPushedNorPulledOver(t *testing.T) {
	// Only the marker was deleted (a dotfile cleanup): the remaining files
	// may hold un-pushed edits, so the pass must neither upload them, nor
	// delete by them, nor destroy them by re-downloading over the folder.
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1", "b.md": "v1"}, "", &logBuf, nil)
	ctx := context.Background()
	require.NoError(t, os.Remove(filepath.Join(local, MarkerPath)))
	writeLocal(t, local, "a.md", "un-pushed edit")
	require.NoError(t, os.Remove(filepath.Join(local, "b.md")))

	runSync(ctx, stores)

	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"a.md": "v1", "b.md": "v1"}, server.Files)
	require.Equal(t, "un-pushed edit", readLocal(t, local, "a.md"))
	require.NoFileExists(t, filepath.Join(local, MarkerPath))
	require.Contains(t, logBuf.String(), "the marker file is gone; leaving the memory store folder as found")
}

func TestFolderTrust_ASingleLocalDeletionWithTheMarkerIntactStillPropagates(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1", "b.md": "v1"}, "", nil, nil)
	clock := withClock(stores)
	ctx := context.Background()
	require.NoError(t, os.Remove(filepath.Join(local, "a.md")))

	runSync(ctx, stores)
	clock.now = clock.now.Add(DeleteCorroborationWindow + time.Second)
	runSync(ctx, stores)

	require.Equal(t, []string{deletedReq("a.md", "v1")}, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"b.md": "v1"}, server.Files)
	require.Equal(t, folderTrustMarker, readLocal(t, local, MarkerPath))
}

func TestFolderTrust_TheMarkerNeverSyncsInEitherDirection(t *testing.T) {
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1"}, "", &logBuf, nil)
	ctx := context.Background()

	// The marker sits in the folder yet is never uploaded or remote-deleted.
	runSync(ctx, stores)
	require.Empty(t, server.ReceivedSnapshot())

	// A marker path the server lists is refused: skipped with a warning,
	// never written over the real marker on disk.
	server.Write(MarkerPath, "not a store id")
	runSync(ctx, stores)
	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, folderTrustMarker, readLocal(t, local, MarkerPath))
	require.Contains(t, logBuf.String(), "reserved marker path")
}

func TestFolderTrust_DisposeToleratesAnAlreadyDeletedFolder(t *testing.T) {
	local, _, stores := downloadedStores(t, map[string]string{"a.md": "v1"}, "", nil, nil)
	require.NoError(t, os.RemoveAll(local))

	stores.Dispose()

	require.NoDirExists(t, local)
}

func TestFolderTrust_DisposeKeepsAFolderThatFailsTheMarkerCheck(t *testing.T) {
	// Sync left the folder as found after a failed marker check; teardown
	// must not then delete the files sync refused to touch.
	var logBuf bytes.Buffer
	local, _, stores := downloadedStores(t, map[string]string{"a.md": "v1"}, "", &logBuf, nil)
	require.NoError(t, os.Remove(filepath.Join(local, MarkerPath)))
	writeLocal(t, local, "a.md", "un-pushed edit")

	stores.Dispose()

	require.Equal(t, "un-pushed edit", readLocal(t, local, "a.md"))
	require.Contains(t, logBuf.String(), "leaving the memory store folder on disk")
}

func TestFolderTrust_DisposeNeverPanicsWhenAStoreRootIsReplacedByAFile(t *testing.T) {
	// The worker calls Dispose from its teardown: a store whose folder
	// something replaced with a regular file must log and be skipped, not
	// panic out of the teardown.
	local, _, stores := downloadedStores(t, map[string]string{"a.md": "v1"}, "", nil, nil)
	require.NoError(t, os.RemoveAll(local))
	require.NoError(t, os.WriteFile(local, []byte("a file where the folder was"), 0o600))

	stores.Dispose()

	got, err := os.ReadFile(local)
	require.NoError(t, err)
	require.Equal(t, "a file where the folder was", string(got))
}

func TestFolderTrust_AMarkerFromAnOlderFormatNoLongerClearsTheFolder(t *testing.T) {
	// The marker carries a format version so an SDK upgrade can retire every
	// folder written before it: a marker in the old format (the bare store
	// id) must fail the check, and the folder is left as found.
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1", "b.md": "v1"}, "", &logBuf, nil)
	ctx := context.Background()
	writeLocal(t, local, MarkerPath, folderTrustStoreID)
	writeLocal(t, local, "a.md", "edit under the old marker")
	require.NoError(t, os.Remove(filepath.Join(local, "b.md")))

	runSync(ctx, stores)

	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"a.md": "v1", "b.md": "v1"}, server.Files)
	require.Equal(t, "edit under the old marker", readLocal(t, local, "a.md"))
	require.Contains(t, logBuf.String(), "the marker file does not match this store")
}

func TestFolderTrust_AFolderWipedWhileAPullIsInFlightIsRebuiltNotDistrusted(t *testing.T) {
	// rm -rf lands between a sync's scan and its content write: the write
	// must not bring the folder back without its marker, or every later sync
	// reads the folder as foreign and the session's memory writes are lost.
	// A flat path fails at the temp-file open, a nested one at the mkdir.
	for _, tc := range []struct{ name, rel string }{{"flat", "z.md"}, {"nested", "sub/z.md"}} {
		t.Run(tc.name, func(t *testing.T) {
			var logBuf bytes.Buffer
			local, server, stores := downloadedStores(t,
				map[string]string{"x.md": "v1", "y.md": "v1", tc.rel: "v1"}, "", &logBuf, nil)
			ctx := context.Background()
			server.Write(tc.rel, "v2")
			server.SetOnFetch(func(string) int {
				_ = os.RemoveAll(local)
				return 0
			})

			runSync(ctx, stores)

			require.NoDirExists(t, local)
			require.Empty(t, server.ReceivedSnapshot())
			require.Contains(t, logBuf.String(), "failed to write memory")

			server.SetOnFetch(nil)
			runSync(ctx, stores)

			require.Equal(t, folderTrustMarker, readLocal(t, local, MarkerPath))
			require.Equal(t, "v1", readLocal(t, local, "x.md"))
			require.Equal(t, "v2", readLocal(t, local, tc.rel))
			require.Contains(t, logBuf.String(), "re-downloading the memory store folder")
			require.NotContains(t, logBuf.String(), "leaving the memory store folder as found")

			stores.Finish(ctx)
			stores.FlushWrites(ctx)
			stores.Dispose()
			require.NoDirExists(t, local)
			require.Empty(t, server.ReceivedSnapshot())
			require.Equal(t, map[string]string{"x.md": "v1", "y.md": "v1", tc.rel: "v2"}, server.Files)
		})
	}
}

func TestFolderTrust_AFolderWipedWhileALocallyDeletedMemoryIsBeingRestoredIsRebuilt(t *testing.T) {
	// A file deleted locally but changed on the server is restored by the
	// next sync; an rm -rf that lands while that restore is downloading must
	// leave no folder behind, the sync after it must rebuild the whole store,
	// and the earlier absence must never become a server delete.
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t,
		map[string]string{"x.md": "v1", "y.md": "v1", "z.md": "v1"}, "", &logBuf, nil)
	clock := withClock(stores)
	ctx := context.Background()
	require.NoError(t, os.Remove(filepath.Join(local, "z.md")))
	server.Write("z.md", "v2")
	server.SetOnFetch(func(string) int {
		_ = os.RemoveAll(local)
		return 0
	})

	runSync(ctx, stores)
	require.NoDirExists(t, local)

	server.SetOnFetch(nil)
	runSync(ctx, stores)
	require.Equal(t, folderTrustMarker, readLocal(t, local, MarkerPath))
	require.Equal(t, "v1", readLocal(t, local, "x.md"))
	require.Equal(t, "v1", readLocal(t, local, "y.md"))
	require.Equal(t, "v2", readLocal(t, local, "z.md"))

	clock.now = clock.now.Add(DeleteCorroborationWindow + time.Second)
	runSync(ctx, stores)
	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"x.md": "v1", "y.md": "v1", "z.md": "v2"}, server.Files)
	require.NotContains(t, logBuf.String(), "leaving the memory store folder as found")
}

func TestFolderTrust_AFolderWipedAgainMidRecoveryIsRebuiltByTheNextSync(t *testing.T) {
	// A re-download writes the marker first and the memories after it; a
	// second rm -rf between those writes must not leave memories on disk
	// without a marker — the next sync re-downloads once more instead.
	var logBuf bytes.Buffer
	disk := &faults{}
	local, server, stores := downloadedStores(t,
		map[string]string{"a.md": "v1", "b.md": "v1", "c.md": "v1"}, "", &logBuf, disk)
	ctx := context.Background()
	require.NoError(t, os.RemoveAll(local))

	memoryPuts := 0
	disk.onPut = func(rel string) {
		if rel == MarkerPath {
			return
		}
		memoryPuts++
		if memoryPuts == 2 { // the marker and one memory are back on disk
			_ = os.RemoveAll(local)
		}
	}
	runSync(ctx, stores)
	disk.onPut = nil

	require.NoDirExists(t, local)
	require.Empty(t, server.ReceivedSnapshot())

	runSync(ctx, stores)

	require.Equal(t, folderTrustMarker, readLocal(t, local, MarkerPath))
	for _, name := range []string{"a.md", "b.md", "c.md"} {
		require.Equal(t, "v1", readLocal(t, local, name))
	}
	require.Equal(t, 2, strings.Count(logBuf.String(), "re-downloading the memory store folder"))
	require.NotContains(t, logBuf.String(), "leaving the memory store folder as found")
	require.Empty(t, server.ReceivedSnapshot())

	stores.Dispose()
	require.NoDirExists(t, local)
}

// ---- server-delete discipline ----------------------------------------------
//
// A locally deleted file is deleted on the server only after a second sync
// confirms it, at a bounded rate per sync, and only when server deletes
// are enabled at all.
//
// The first sync that sees the file missing records it; a later sync —
// after DeleteCorroborationWindow (skipped on the session's final Finish, the
// session's last sync) and a re-check that the file is still gone and the
// marker still intact — sends the DELETE. Each sync sends a bounded
// number of deletes; SyncDeletions = MemorySyncDeletionsLogOnly only logs them,
// and MemorySyncDeletionsDisabled sends none, ever.

var deleteElapsed = DeleteCorroborationWindow + time.Second

// downloadedMode is downloadedStores with a SyncDeletions mode and a fake
// clock at time zero.
func downloadedMode(t *testing.T, initial map[string]string, mode MemoryDeleteMode, logBuf *bytes.Buffer) (string, *memoryServer, *SessionMemoryStores, *testClock) {
	t.Helper()
	local, server, stores := downloadedStores(t, initial, "", logBuf, nil)
	stores.syncDeletions = mode
	clock := withClock(stores)
	return local, server, stores, clock
}

func TestDeleteDiscipline_ALocalDeletionIsConfirmedBeforeItReachesTheServer(t *testing.T) {
	local, server, stores, clock := downloadedMode(t,
		map[string]string{"a.md": "v1", "keep.md": "v1"}, MemorySyncDeletionsEnabled, nil)
	ctx := context.Background()
	require.NoError(t, os.Remove(filepath.Join(local, "a.md")))

	// The first pass only records the absence.
	runSync(ctx, stores)
	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"a.md": "v1", "keep.md": "v1"}, server.Files)

	// Still inside the window: nothing goes out.
	clock.now = clock.now.Add(DeleteCorroborationWindow - time.Second)
	runSync(ctx, stores)
	require.Empty(t, server.ReceivedSnapshot())

	// Window elapsed and the file is still gone: the guarded delete is sent.
	clock.now = time.Unix(0, 0).Add(deleteElapsed)
	runSync(ctx, stores)
	require.Equal(t, []string{deletedReq("a.md", "v1")}, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"keep.md": "v1"}, server.Files)
}

func TestDeleteDiscipline_AFileThatReappearsIsNotDeletedAndLeavesPending(t *testing.T) {
	local, server, stores, clock := downloadedMode(t, map[string]string{"a.md": "v1"}, MemorySyncDeletionsEnabled, nil)
	ctx := context.Background()
	require.NoError(t, os.Remove(filepath.Join(local, "a.md")))
	runSync(ctx, stores)

	writeLocal(t, local, "a.md", "v1")
	clock.now = time.Unix(0, 0).Add(deleteElapsed)
	runSync(ctx, stores)
	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"a.md": "v1"}, server.Files)

	// The reappearance cleared the pending entry: deleting again starts a
	// fresh window instead of firing on the stale observation.
	require.NoError(t, os.Remove(filepath.Join(local, "a.md")))
	runSync(ctx, stores)
	require.Empty(t, server.ReceivedSnapshot())
	clock.now = time.Unix(0, 0).Add(2 * deleteElapsed)
	runSync(ctx, stores)
	require.Equal(t, []string{deletedReq("a.md", "v1")}, server.ReceivedSnapshot())
}

func TestDeleteDiscipline_DisabledModeNeverDeletesButStillSyncs(t *testing.T) {
	var logBuf bytes.Buffer
	local, server, stores, clock := downloadedMode(t,
		map[string]string{"gone.md": "v1", "push.md": "v1", "pull.md": "v1"}, MemorySyncDeletionsDisabled, &logBuf)
	ctx := context.Background()
	require.NoError(t, os.Remove(filepath.Join(local, "gone.md")))
	writeLocal(t, local, "push.md", "v2 from agent")
	server.Write("pull.md", "v2 from server")

	runSync(ctx, stores)
	clock.now = time.Unix(0, 0).Add(deleteElapsed)
	runSync(ctx, stores)
	clock.now = time.Unix(0, 0).Add(2 * deleteElapsed)
	runSync(ctx, stores)

	// Uploads and pulls still flow; the deletion never does.
	require.Equal(t, []string{updated("push.md", "v2 from agent", "v1")}, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"gone.md": "v1", "push.md": "v2 from agent", "pull.md": "v2 from server"}, server.Files)
	require.Equal(t, "v2 from server", readLocal(t, local, "pull.md"))
	require.Contains(t, logBuf.String(), "server deletes are disabled")
}

func TestDeleteDiscipline_ThePerPassCapSpreadsDeletesOverPasses(t *testing.T) {
	var logBuf bytes.Buffer
	files := map[string]string{}
	for i := 0; i < 14; i++ {
		files[fmt.Sprintf("f%02d.md", i)] = "v1"
	}
	local, server, stores, clock := downloadedMode(t, files, MemorySyncDeletionsEnabled, &logBuf)
	ctx := context.Background()
	var doomed []string
	for k := range files {
		doomed = append(doomed, k)
	}
	slices.Sort(doomed)
	doomed = doomed[:12]
	for _, name := range doomed {
		require.NoError(t, os.Remove(filepath.Join(local, name)))
	}

	runSync(ctx, stores)
	require.Empty(t, server.ReceivedSnapshot())

	// cap = max(8, min(50, 14/4)) = 8: the first confirming sync stops there.
	clock.now = time.Unix(0, 0).Add(deleteElapsed)
	runSync(ctx, stores)
	var want []string
	for _, name := range doomed[:8] {
		want = append(want, deletedReq(name, "v1"))
	}
	require.Equal(t, want, server.ReceivedSnapshot())
	require.Contains(t, logBuf.String(), "delete cap reached: sent 8 deletes, held 4 for later syncs")

	// The held-back deletions confirm on the next pass.
	runSync(ctx, stores)
	for _, name := range doomed[8:] {
		want = append(want, deletedReq(name, "v1"))
	}
	require.Equal(t, want, server.ReceivedSnapshot())
	remaining := map[string]string{}
	for i := 12; i < 14; i++ {
		remaining[fmt.Sprintf("f%02d.md", i)] = "v1"
	}
	require.Equal(t, remaining, server.Files)
}

func TestDeleteDiscipline_TheFinalSyncWaivesTheWindowButStillReChecks(t *testing.T) {
	// A deletion made in a session shorter than the waiting window still
	// goes through: the final sync skips the wait, and the re-check just
	// before the delete stands in for the second observation.
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1", "keep.md": "v1"}, "", nil, nil)
	ctx := context.Background()
	require.NoError(t, os.Remove(filepath.Join(local, "a.md")))

	stores.Finish(ctx)

	require.Equal(t, []string{deletedReq("a.md", "v1")}, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"keep.md": "v1"}, server.Files)
}

func TestDeleteDiscipline_AFolderDestroyedMidSyncSendsNoDeletes(t *testing.T) {
	// The directory scan saw an intact folder, then the folder vanished
	// while the server listing was in flight: the re-check just before
	// the delete must notice the marker is gone and hold every delete.
	local, server, stores, clock := downloadedMode(t, map[string]string{"a.md": "v1", "keep.md": "v1"}, MemorySyncDeletionsEnabled, nil)
	ctx := context.Background()
	require.NoError(t, os.Remove(filepath.Join(local, "a.md")))
	runSync(ctx, stores)
	clock.now = time.Unix(0, 0).Add(deleteElapsed)

	server.SetOnList(func() { _ = os.RemoveAll(local) })
	runSync(ctx, stores)
	server.SetOnList(nil)

	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"a.md": "v1", "keep.md": "v1"}, server.Files)
}

func TestDeleteDiscipline_AnUnreadableSubdirectoryFailsThePassInsteadOfReadingAsDeletions(t *testing.T) {
	// Files under a directory the scan cannot enter are not treated as deleted.
	if os.Geteuid() == 0 {
		t.Skip("permission checks do not bind root")
	}
	var logBuf bytes.Buffer
	local, server, stores, clock := downloadedMode(t,
		map[string]string{"sub/a.md": "v1", "sub/b.md": "v1", "top.md": "v1"}, MemorySyncDeletionsEnabled, &logBuf)
	ctx := context.Background()
	writeLocal(t, local, "top.md", "v2")
	sub := filepath.Join(local, "sub")
	require.NoError(t, os.Chmod(sub, 0o000))
	t.Cleanup(func() { _ = os.Chmod(sub, 0o700) })

	runSync(ctx, stores)
	clock.now = time.Unix(0, 0).Add(deleteElapsed)
	runSync(ctx, stores)
	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, 2, strings.Count(logBuf.String(), "memory sync failed"))

	// Readable again: the edit goes out and nothing is deleted.
	require.NoError(t, os.Chmod(sub, 0o700))
	clock.now = time.Unix(0, 0).Add(2 * deleteElapsed)
	runSync(ctx, stores)
	require.Equal(t, []string{updated("top.md", "v2", "v1")}, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"sub/a.md": "v1", "sub/b.md": "v1", "top.md": "v2"}, server.Files)
}

func TestDeleteDiscipline_ARecoveryReDownloadClearsPendingDeletes(t *testing.T) {
	// The re-download rebuilt the disk, so an absence observed before it
	// says nothing about the rebuilt folder.
	local, server, stores, clock := downloadedMode(t, map[string]string{"a.md": "v1"}, MemorySyncDeletionsEnabled, nil)
	ctx := context.Background()
	require.NoError(t, os.Remove(filepath.Join(local, "a.md")))
	runSync(ctx, stores)

	clock.now = time.Unix(0, 0).Add(deleteElapsed)
	require.NoError(t, os.RemoveAll(local))
	runSync(ctx, stores)
	require.Equal(t, "v1", readLocal(t, local, "a.md"))

	// A fresh deletion starts a fresh window — the pre-recovery
	// observation must not corroborate it.
	require.NoError(t, os.Remove(filepath.Join(local, "a.md")))
	runSync(ctx, stores)
	require.Empty(t, server.ReceivedSnapshot())
	clock.now = time.Unix(0, 0).Add(2 * deleteElapsed)
	runSync(ctx, stores)
	require.Equal(t, []string{deletedReq("a.md", "v1")}, server.ReceivedSnapshot())
}

func TestDeleteDiscipline_LogOnlyModeLogsTheDeleteButNeverSendsIt(t *testing.T) {
	// The dry-run setting: the full pipeline runs — window, re-checks,
	// cap — but the confirmed delete is logged instead of sent, on every
	// sync while the file stays missing.
	var logBuf bytes.Buffer
	local, server, stores, clock := downloadedMode(t,
		map[string]string{"a.md": "v1", "keep.md": "v1"}, MemorySyncDeletionsLogOnly, &logBuf)
	ctx := context.Background()
	require.NoError(t, os.Remove(filepath.Join(local, "a.md")))

	runSync(ctx, stores)
	require.NotContains(t, logBuf.String(), "log-only") // still inside the waiting window

	clock.now = time.Unix(0, 0).Add(deleteElapsed)
	runSync(ctx, stores)
	runSync(ctx, stores)

	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"a.md": "v1", "keep.md": "v1"}, server.Files)
	require.Equal(t, 2, strings.Count(logBuf.String(), "log-only: sync would delete this memory on the server"))
	require.Contains(t, logBuf.String(), "path=a.md")
}

func TestDeleteDiscipline_FinishRunsOnceAndASecondCallPanics(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1"}, "", nil, nil)
	ctx := context.Background()
	writeLocal(t, local, "a.md", "v2")

	stores.Finish(ctx)
	require.Equal(t, map[string]string{"a.md": "v2"}, server.Files)

	require.PanicsWithValue(t,
		"Finish() was already called: it is the session's last sync and runs once",
		func() { stores.Finish(ctx) })
}

// ---- never destroy the only copy -------------------------------------------
//
// A memory deleted remotely while its file holds an un-pushed edit keeps
// the file — a writable store re-creates the memory from it, a read-only
// store keeps it on disk unsynced. Only an unedited file follows the
// remote delete off the disk.

func TestKeepOnlyCopy_AnEditedFileSurvivesARemoteDeleteAndIsReCreated(t *testing.T) {
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t,
		map[string]string{"edited.md": "v1", "plain.md": "v1"}, "", &logBuf, nil)
	ctx := context.Background()
	writeLocal(t, local, "edited.md", "v2 from agent")
	server.Delete("edited.md")
	server.Delete("plain.md")

	runSync(ctx, stores)

	// The edit's only copy survives and the memory is re-created from it;
	// the unedited file follows the remote delete off the disk.
	require.Equal(t, "v2 from agent", readLocal(t, local, "edited.md"))
	require.Equal(t, map[string]string{"edited.md": "v2 from agent"}, server.Files)
	require.Equal(t, []string{created("edited.md", "v2 from agent")}, server.ReceivedSnapshot())
	require.NoFileExists(t, filepath.Join(local, "plain.md"))
	require.Contains(t, logBuf.String(), "re-creating it from the file")

	// Settled: nothing more goes out.
	runSync(ctx, stores)
	require.Equal(t, []string{created("edited.md", "v2 from agent")}, server.ReceivedSnapshot())

	// The upload's sha became the baseline: the next edit is a guarded update.
	writeLocal(t, local, "edited.md", "v3")
	runSync(ctx, stores)
	got := server.ReceivedSnapshot()
	require.Equal(t, updated("edited.md", "v3", "v2 from agent"), got[len(got)-1])
}

func TestKeepOnlyCopy_AReadOnlyStoreKeepsTheEditedFileWithoutPushing(t *testing.T) {
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t,
		map[string]string{"edited.md": "v1", "plain.md": "v1"}, "read_only", &logBuf, nil)
	ctx := context.Background()
	writeLocal(t, local, "edited.md", "v2 from agent")
	server.Delete("edited.md")
	server.Delete("plain.md")

	runSync(ctx, stores)

	require.Equal(t, "v2 from agent", readLocal(t, local, "edited.md"))
	require.NoFileExists(t, filepath.Join(local, "plain.md"))
	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, 1, strings.Count(logBuf.String(), "cannot push"))

	// Later passes see a local-only file in a read-only store: a quiet no-op.
	runSync(ctx, stores)
	runSync(ctx, stores)
	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, "v2 from agent", readLocal(t, local, "edited.md"))
	require.Equal(t, 1, strings.Count(logBuf.String(), "cannot push"))
}

func TestKeepOnlyCopy_ARefusedEditSurvivesARemoteDeleteWithoutARetry(t *testing.T) {
	// The server already refused this content, so the remote delete
	// neither removes the file nor triggers a re-create the server would
	// refuse again.
	local, server, stores := downloadedStores(t, map[string]string{"big.md": "v1"}, "", nil, nil)
	ctx := context.Background()
	server.MaxContentBytes = 64

	writeLocal(t, local, "big.md", strings.Repeat("x", 100))
	runSync(ctx, stores)
	require.Len(t, server.ReceivedSnapshot(), 1) // the one refused update

	server.Delete("big.md")
	runSync(ctx, stores)
	runSync(ctx, stores)

	require.Equal(t, strings.Repeat("x", 100), readLocal(t, local, "big.md"))
	require.Len(t, server.ReceivedSnapshot(), 1)
	require.NotContains(t, server.Files, "big.md")
}

func TestKeepOnlyCopy_AMemoryDeletedBetweenTheListingAndTheUpdateIsReCreated(t *testing.T) {
	// An update that meets a 404 re-creates the memory in the same pass.
	local, server, stores := downloadedStores(t, map[string]string{"note.md": "v1"}, "", nil, nil)
	ctx := context.Background()
	writeLocal(t, local, "note.md", "v2 from agent")

	server.SetUploadHook(func(_ context.Context, path string) {
		server.Delete(path)
		server.SetUploadHook(nil)
	})
	runSync(ctx, stores)

	require.Equal(t, map[string]string{"note.md": "v2 from agent"}, server.Files)
	require.Equal(t, []string{updated("note.md", "v2 from agent", "v1"), created("note.md", "v2 from agent")}, server.ReceivedSnapshot())

	// The create's sha is the new baseline.
	writeLocal(t, local, "note.md", "v3")
	runSync(ctx, stores)
	received := server.ReceivedSnapshot()
	require.Equal(t, updated("note.md", "v3", "v2 from agent"), received[len(received)-1])
}

func TestKeepOnlyCopy_AnEditWrittenMidSyncStillSurvivesARemoteDelete(t *testing.T) {
	// The directory scan saw the file unedited; an edit is written while
	// the server listing is in flight. The removal must re-read the file
	// instead of destroying the edit based on the stale scan.
	local, server, stores := downloadedStores(t, map[string]string{"z.md": "v1", "keep.md": "v1"}, "", nil, nil)
	ctx := context.Background()
	server.Delete("z.md")

	server.SetOnList(func() { require.NoError(t, os.WriteFile(filepath.Join(local, "z.md"), []byte("late edit"), 0o600)) })
	runSync(ctx, stores)
	server.SetOnList(nil)

	require.Equal(t, "late edit", readLocal(t, local, "z.md"))
	require.Equal(t, "late edit", server.Files["z.md"])
}

func TestKeepOnlyCopy_AReadOnlyPullOverALocalEditIsWarned(t *testing.T) {
	// Read-only stores always take the server version, but overwriting a
	// local edit must leave a trace instead of happening silently.
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"pull.md": "v1"}, "read_only", &logBuf, nil)
	ctx := context.Background()
	writeLocal(t, local, "pull.md", "local edit")
	server.Write("pull.md", "v2 from server")

	runSync(ctx, stores)

	require.Equal(t, "v2 from server", readLocal(t, local, "pull.md"))
	require.Empty(t, server.ReceivedSnapshot())
	require.Contains(t, logBuf.String(), "changed both locally and remotely")
}

func TestKeepOnlyCopy_AFailedReReadDoesNotResurrectAServerDeletedMemory(t *testing.T) {
	// The re-read before removal can fail (a transient I/O error). That
	// must mean "retry next sync", not "file gone": dropping the entry
	// would make the next sync upload the file as new, re-creating the
	// memory the server just deleted.
	disk := &faults{}
	local, server, stores := downloadedStores(t, map[string]string{"z.md": "v1", "keep.md": "v1"}, "", nil, disk)
	ctx := context.Background()
	server.Delete("z.md")

	disk.hashFile = func(rel string) bool { return rel == "z.md" }
	runSync(ctx, stores)
	disk.hashFile = nil

	// The error held everything: file still on disk, nothing pushed.
	require.Equal(t, "v1", readLocal(t, local, "z.md"))
	require.Empty(t, server.ReceivedSnapshot())

	// With the error gone, the next sync completes the server delete.
	runSync(ctx, stores)
	require.NoFileExists(t, filepath.Join(local, "z.md"))
	require.NotContains(t, server.Files, "z.md")
	require.Empty(t, server.ReceivedSnapshot())
}

// ---- FlushWrites — the push-only shutdown pass -----------------------------
//
// A session that ends on an error or cancel gets no full reconcile, but
// the writes already on disk are the only copy of the agent's edits.
// FlushWrites pushes those and does nothing else: no remote deletes, no
// local removals, no pulls.

func TestFinalFlush_PushesNewAndChangedFilesAndNothingElse(t *testing.T) {
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{
		"a-pull.md":               "v1",
		"b-push.md":               "v1",
		"d-local-gone.md":         "v1",
		"e-remote-gone.md":        "v1",
		"f-remote-gone-edited.md": "v1",
		"g-both.md":               "v1",
	}, "", &logBuf, nil)
	ctx := context.Background()

	writeLocal(t, local, "b-push.md", "v2 from agent")
	writeLocal(t, local, "c-new.md", "new from agent")
	require.NoError(t, os.Remove(filepath.Join(local, "d-local-gone.md")))
	writeLocal(t, local, "f-remote-gone-edited.md", "v2 from agent")
	writeLocal(t, local, "g-both.md", "v2 from agent")
	server.Write("a-pull.md", "v2 from server")
	server.Delete("e-remote-gone.md")
	server.Delete("f-remote-gone-edited.md")
	server.Write("g-both.md", "v2 from server")

	stores.FlushWrites(ctx)

	// Only additions went out: the update, the new file, and the
	// re-created edited file whose remote copy was deleted.
	sent := []string{
		updated("b-push.md", "v2 from agent", "v1"),
		created("c-new.md", "new from agent"),
		created("f-remote-gone-edited.md", "v2 from agent"),
	}
	require.ElementsMatch(t, sent, server.ReceivedSnapshot())
	// A locally missing file was not deleted remotely.
	require.Equal(t, "v1", server.Files["d-local-gone.md"])
	// A remotely deleted, locally unedited memory was not re-uploaded —
	// and its file was not removed from disk either.
	require.NotContains(t, server.Files, "e-remote-gone.md")
	require.Equal(t, "v1", readLocal(t, local, "e-remote-gone.md"))
	// Remote changes were not pulled.
	require.Equal(t, "v1", readLocal(t, local, "a-pull.md"))
	// The conflict was left to the server, and logged.
	require.Equal(t, "v2 from server", server.Files["g-both.md"])
	require.Equal(t, "v2 from agent", readLocal(t, local, "g-both.md"))
	require.Contains(t, logBuf.String(), "changed both locally and remotely")

	// Settled: a second flush pushes nothing more.
	stores.FlushWrites(ctx)
	require.ElementsMatch(t, sent, server.ReceivedSnapshot())
}

func TestFinalFlush_UploadsSeveralFilesAtOnceUpToTheBound(t *testing.T) {
	// The bound's worth of uploads are in flight together, and never one
	// more — a serial flush would time out waiting for the second to start.
	local, server, stores := downloadedStores(t, map[string]string{}, "", nil, nil)
	want := map[string]string{}
	for i := range uploadConcurrency + 3 {
		name, content := fmt.Sprintf("f%02d.md", i), fmt.Sprintf("file %d", i)
		writeLocal(t, local, name, content)
		want[name] = content
	}

	var mu sync.Mutex
	inFlight, peak := 0, 0
	release := make(chan struct{})
	var once sync.Once
	server.SetUploadHook(func(ctx context.Context, _ string) {
		mu.Lock()
		inFlight++
		peak = max(peak, inFlight)
		full := inFlight == uploadConcurrency
		mu.Unlock()
		if full {
			once.Do(func() {
				// Hold the full set briefly so a flush that ignored the bound
				// has time to start one more before any of these finish.
				time.Sleep(100 * time.Millisecond)
				close(release)
			})
		}
		select {
		case <-release:
		case <-ctx.Done():
		case <-time.After(10 * time.Second):
			t.Error("the uploads never filled the bound: they ran serially")
			once.Do(func() { close(release) })
		}
		mu.Lock()
		inFlight--
		mu.Unlock()
	})

	stores.FlushWrites(context.Background())

	mu.Lock()
	got := peak
	mu.Unlock()
	require.Equal(t, uploadConcurrency, got)
	require.Equal(t, want, server.Files)
}

func TestFinalFlush_AFlushCutOffPartWayLogsHowManyFilesItHadNotUploaded(t *testing.T) {
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"z.md": "z0"}, "", &logBuf, nil)
	names := []string{"a.md", "b.md", "c.md", "y.md", "z.md"}
	for _, name := range names {
		writeLocal(t, local, name, "content of "+name)
	}

	// Two uploads — one create, one update — never complete; the flush's
	// own deadline cuts them off.
	server.SetUploadHook(func(ctx context.Context, path string) {
		if path == "y.md" || path == "z.md" {
			<-ctx.Done()
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	stores.FlushWrites(ctx)

	require.Equal(t, 1, strings.Count(logBuf.String(), "memory flush cut off part-way"))
	require.Contains(t, logBuf.String(),
		`2 of 5 changed files had not finished uploading" memory_store_id=memstore_notes`)
	require.Equal(t, map[string]string{
		"a.md": "content of a.md", "b.md": "content of b.md", "c.md": "content of c.md", "z.md": "z0",
	}, server.Files)
	require.NotContains(t, logBuf.String(), "failed to upload memory")
	require.NotContains(t, logBuf.String(), "memory flush failed")

	// The three that made it are not sent again; the two that did not are.
	// A flush that finishes in time adds no cut-off line.
	server.SetUploadHook(nil)
	stores.FlushWrites(context.Background())
	want := map[string]string{}
	for _, name := range names {
		want[name] = "content of " + name
	}
	require.Equal(t, want, server.Files)
	require.ElementsMatch(t, []string{
		created("a.md", "content of a.md"),
		created("b.md", "content of b.md"),
		created("c.md", "content of c.md"),
		created("y.md", "content of y.md"),
		updated("z.md", "content of z.md", "z0"),
	}, server.ReceivedSnapshot())
	require.Equal(t, 1, strings.Count(logBuf.String(), "cut off part-way"))
}

func TestFinalFlush_AdoptsAFileTheServerAlreadyHolds(t *testing.T) {
	// A final sync cut off after uploading a file leaves the baseline stale;
	// the flush sees the server already holds those bytes and records that
	// instead of reporting a conflict.
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"n.md": "v1"}, "", &logBuf, nil)
	ctx := context.Background()
	writeLocal(t, local, "n.md", "v2")
	server.Write("n.md", "v2")

	stores.FlushWrites(ctx)

	require.Empty(t, server.ReceivedSnapshot())
	require.NotContains(t, logBuf.String(), "changed both locally and remotely")

	// The adopted sha became the baseline: the next edit is a guarded update.
	writeLocal(t, local, "n.md", "v3")
	stores.FlushWrites(ctx)
	require.Equal(t, []string{updated("n.md", "v3", "v2")}, server.ReceivedSnapshot())
}

func TestFinalFlush_ARefusedFileDoesNotHoldBackTheOthers(t *testing.T) {
	// A refusal and a success are book-kept from two uploads of one flush:
	// the refused file is tried once and not retried, the other lands once.
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{}, "", &logBuf, nil)
	ctx := context.Background()
	server.MaxContentBytes = 64
	writeLocal(t, local, "big.md", strings.Repeat("x", 100))
	writeLocal(t, local, "ok.md", "fits")

	stores.FlushWrites(ctx)
	stores.FlushWrites(ctx)

	require.Equal(t, "fits", server.Files["ok.md"])
	require.NotContains(t, server.Files, "big.md")
	require.ElementsMatch(t,
		[]string{created("big.md", strings.Repeat("x", 100)), created("ok.md", "fits")},
		server.ReceivedSnapshot())
	require.Equal(t, 1, strings.Count(logBuf.String(), "stays un-synced until its content changes"))
}

func TestFinalFlush_NeverPushesAReadOnlyStore(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"edit.md": "v1"}, "read_only", nil, nil)
	ctx := context.Background()

	writeLocal(t, local, "edit.md", "v2 from agent")
	writeLocal(t, local, "new.md", "new from agent")

	stores.FlushWrites(ctx)

	require.Empty(t, server.ReceivedSnapshot())
	require.Equal(t, map[string]string{"edit.md": "v1"}, server.Files)
}

func TestFinalFlush_SkipsFilesTheServerAlreadyRefused(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"note.md": "v1"}, "", nil, nil)
	ctx := context.Background()
	server.MaxContentBytes = 64

	writeLocal(t, local, "big.md", strings.Repeat("x", 100))
	runSync(ctx, stores)
	creates := 0
	for _, r := range server.ReceivedSnapshot() {
		if strings.HasPrefix(r, "create ") {
			creates++
		}
	}
	require.Equal(t, 1, creates)

	stores.FlushWrites(ctx)

	// The refused file was not retried; nothing else needed pushing.
	creates = 0
	for _, r := range server.ReceivedSnapshot() {
		if strings.HasPrefix(r, "create ") {
			creates++
		}
	}
	require.Equal(t, 1, creates)
	require.NotContains(t, server.Files, "big.md")
}

func TestFinalFlush_PushesNothingWhenTheMarkerCheckFails(t *testing.T) {
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1"}, "", &logBuf, nil)
	ctx := context.Background()

	require.NoError(t, os.Remove(filepath.Join(local, MarkerPath)))
	writeLocal(t, local, "foreign.md", "someone else's notes")

	stores.FlushWrites(ctx)

	require.Empty(t, server.ReceivedSnapshot())
	require.NotContains(t, server.Files, "foreign.md")
	// Unlike sync's recovery, the flush does not re-download the folder.
	require.NoFileExists(t, filepath.Join(local, MarkerPath))
	require.Contains(t, logBuf.String(), "not uploading anything from the memory store folder")
}

func TestFinalFlush_NeverFails(t *testing.T) {
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1"}, "", &logBuf, nil)
	ctx := context.Background()
	writeLocal(t, local, "a.md", "v2 from agent")

	server.SetOnList(func() { panic(http.ErrAbortHandler) }) // the listing explodes
	stores.FlushWrites(ctx)
	server.SetOnList(nil)

	require.Empty(t, server.ReceivedSnapshot())
	require.Contains(t, logBuf.String(), "memory flush failed")
}
