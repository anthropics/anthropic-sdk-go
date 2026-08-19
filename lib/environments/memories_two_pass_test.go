package environments

// The sync's two passes: a content-free listing, then targeted fetches.
//
// A sync lists shas only (basic view) and fetches a memory's body just
// before writing it to disk; only the download takes full pages, since it
// needs every memory anyway. These tests pin that fetch discipline: a sync
// that writes nothing fetches nothing, a changed path fetches exactly
// itself, fetches fan out rather than queue, and a fetch that fails — or
// races a server-side delete — leaves the old state for the next sync.

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stretchr/testify/require"
)

func TestTwoPass_ASyncThatPullsNothingFetchesNoContent(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1", "b.md": "v1", "c.md": "v1"}, "", nil, nil)
	ctx := context.Background()
	// The download took full listing pages: no single-memory fetches at all.
	require.Empty(t, server.ContentFetchesSnapshot())

	writeLocal(t, local, "a.md", "v2 from agent") // a push needs no remote content
	runSync(ctx, stores)

	require.Empty(t, server.ContentFetchesSnapshot())
	require.Equal(t, "v2 from agent", server.Files["a.md"])
}

func TestTwoPass_OnlyTheChangedMemoryIsFetched(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1", "b.md": "v1"}, "", nil, nil)
	ctx := context.Background()
	server.ClearContentFetches()

	server.Write("a.md", "v2 from server")
	runSync(ctx, stores)

	require.Equal(t, []string{"a.md"}, server.ContentFetchesSnapshot())
	require.Equal(t, "v2 from server", readLocal(t, local, "a.md"))
	require.Equal(t, "v1", readLocal(t, local, "b.md"))
}

func TestTwoPass_AFileAlreadyHoldingTheRemoteContentIsAdoptedWithoutAFetch(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1"}, "", nil, nil)
	ctx := context.Background()
	server.ClearContentFetches()

	// Both sides landed on the same bytes independently.
	writeLocal(t, local, "a.md", "v2 on both")
	server.Write("a.md", "v2 on both")
	runSync(ctx, stores)

	require.Empty(t, server.ContentFetchesSnapshot())
	require.Empty(t, server.ReceivedSnapshot())

	// The baseline advanced to the adopted sha: the next edit pushes guarded by it.
	writeLocal(t, local, "a.md", "v3 from agent")
	runSync(ctx, stores)
	require.Equal(t, []string{updated("a.md", "v3 from agent", "v2 on both")}, server.ReceivedSnapshot())
}

func TestTwoPass_ContentFetchesFanOutWithinAStore(t *testing.T) {
	// The first memory's fetch completes only after the second's starts —
	// only parallel fetches satisfy that; serial ones would time out.
	server := newMemoryServer(t, map[string]string{"a.md": "v1", "b.md": "v1"}, "")
	secondStarted := make(chan struct{})
	httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isFetch := r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/memories/")
		if isFetch {
			if strings.HasSuffix(r.URL.Path, "/"+memID("a.md")) {
				select {
				case <-secondStarted:
				case <-time.After(10 * time.Second):
					t.Error("the second fetch never started: fetches ran serially")
				}
			} else {
				close(secondStarted)
			}
		}
		server.serveHTTP(w, r)
	}))
	t.Cleanup(httpSrv.Close)
	client := anthropic.NewClient(option.WithBaseURL(httpSrv.URL), option.WithAPIKey("k"), option.WithMaxRetries(0))
	workdir := t.TempDir()
	stores := mustStores(t, client, SessionMemoryStoresOptions{Workdir: workdir, Logger: silentLogger})
	ctx := context.Background()
	require.NoError(t, stores.Download(ctx, fetchSession(t, client)))
	local := filepath.Join(workdir, "memory", "notes")

	server.Write("a.md", "v2 from server")
	server.Write("b.md", "v2 from server")
	runSync(ctx, stores)

	require.Equal(t, "v2 from server", readLocal(t, local, "a.md"))
	require.Equal(t, "v2 from server", readLocal(t, local, "b.md"))
}

func TestTwoPass_AFailedFetchKeepsTheOldStateAndTheNextSyncRetries(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1"}, "", nil, nil)
	ctx := context.Background()
	server.Write("a.md", "v2 from server")

	server.SetOnFetch(func(string) int { return 500 })
	runSync(ctx, stores)
	require.Equal(t, "v1", readLocal(t, local, "a.md"))

	server.SetOnFetch(nil)
	runSync(ctx, stores)
	require.Equal(t, "v2 from server", readLocal(t, local, "a.md"))
	// The stale baseline drove no uploads or deletes along the way.
	require.Empty(t, server.ReceivedSnapshot())
}

func TestTwoPass_AMemoryDeletedBetweenListingAndFetchWaitsForTheNextSync(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"a.md": "v1", "keep.md": "v1"}, "", nil, nil)
	ctx := context.Background()
	server.Write("a.md", "v2 from server")

	// Called with the server mutex held, so Files is safe to touch.
	server.SetOnFetch(func(path string) int {
		delete(server.Files, path)
		return 0
	})
	runSync(ctx, stores)
	// The fetch 404ed: disk and baseline keep the old version for now.
	require.Equal(t, "v1", readLocal(t, local, "a.md"))

	server.SetOnFetch(nil)
	runSync(ctx, stores)
	// The next sync saw the server-side delete and took the local copy with it.
	require.NoFileExists(t, filepath.Join(local, "a.md"))
	require.Empty(t, server.ReceivedSnapshot())
}

func TestTwoPass_ListingsAskForTheLargestPageTheViewAllows(t *testing.T) {
	// The download lists full pages at the server's cap for that view; the
	// sync lists basic pages at the API maximum.
	server := newMemoryServer(t, map[string]string{"a.md": "v1"}, "")
	var listQueries []string
	httpSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/memories") {
			q := r.URL.Query()
			listQueries = append(listQueries, "view="+q.Get("view")+" limit="+q.Get("limit"))
		}
		server.serveHTTP(w, r)
	}))
	t.Cleanup(httpSrv.Close)
	client := anthropic.NewClient(option.WithBaseURL(httpSrv.URL), option.WithAPIKey("k"), option.WithMaxRetries(0))
	stores := mustStores(t, client, SessionMemoryStoresOptions{Workdir: t.TempDir(), Logger: silentLogger})
	ctx := context.Background()
	require.NoError(t, stores.Download(ctx, fetchSession(t, client)))
	runSync(ctx, stores)

	require.Equal(t, []string{"view=full limit=20", "view=basic limit=100"}, listQueries)
}
