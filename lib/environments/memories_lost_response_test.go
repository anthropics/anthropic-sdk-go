package environments

// An upload whose response never arrives.
//
// The server can save a create or update and the response still be lost — a
// timeout, a dropped connection. The sync then cannot tell whether the bytes
// landed. If they did, the server's copy is the sync's own upload, not
// another writer's edit, and a newer local edit must go out over it.

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stretchr/testify/require"
)

func TestLostResponse_AnEditMadeAfterALostUploadResponseStillReachesTheServer(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"note.md": "v0"}, "", nil, nil)
	ctx := context.Background()

	writeLocal(t, local, "note.md", "v1")
	writeLocal(t, local, "new.md", "v1")
	server.SetLoseUploadResponses(true)
	runSync(ctx, stores)
	server.SetLoseUploadResponses(false)
	// The server holds both writes; the sync never heard back.
	require.Equal(t, map[string]string{"note.md": "v1", "new.md": "v1"}, server.FilesSnapshot())

	writeLocal(t, local, "note.md", "v2")
	writeLocal(t, local, "new.md", "v2")
	runSync(ctx, stores)

	// The server's v1 was this sync's own upload, so v2 goes out over it.
	require.Equal(t, map[string]string{"note.md": "v2", "new.md": "v2"}, server.FilesSnapshot())
	require.Equal(t, "v2", readLocal(t, local, "note.md"))
	require.Equal(t, "v2", readLocal(t, local, "new.md"))
	require.Equal(t, []string{
		created("new.md", "v1"),
		updated("note.md", "v1", "v0"),
		updated("new.md", "v2", "v1"),
		updated("note.md", "v2", "v1"),
	}, server.ReceivedSnapshot())

	// Settled: nothing more to send.
	runSync(ctx, stores)
	require.Len(t, server.ReceivedSnapshot(), 4)
}

func TestLostResponse_ARetryWhoseResponseIsLostIsRememberedToo(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"note.md": "v0"}, "", nil, nil)
	ctx := context.Background()

	writeLocal(t, local, "note.md", "v1")
	server.FailUploads = 500
	runSync(ctx, stores)
	server.FailUploads = 0
	// The retry sends the same bytes; this time they land and the response is lost.
	server.SetLoseUploadResponses(true)
	runSync(ctx, stores)
	server.SetLoseUploadResponses(false)
	require.Equal(t, "v1", server.FilesSnapshot()["note.md"])

	writeLocal(t, local, "note.md", "v2")
	runSync(ctx, stores)

	require.Equal(t, "v2", server.FilesSnapshot()["note.md"])
	require.Equal(t, "v2", readLocal(t, local, "note.md"))
}

func TestLostResponse_AClientRetryRefusedBecauseTheFirstAttemptLandedIsRemembered(t *testing.T) {
	// The client's own retry re-sends the update with the old precondition;
	// the first attempt already moved the server, so the retry gets a 409.
	server := newMemoryServer(t, map[string]string{"note.md": "v0"}, "")
	httpSrv := httptest.NewServer(http.HandlerFunc(server.serveHTTP))
	t.Cleanup(httpSrv.Close)
	client := anthropic.NewClient(option.WithBaseURL(httpSrv.URL), option.WithAPIKey("k"), option.WithMaxRetries(1))
	workdir := t.TempDir()
	stores := mustStores(t, client, SessionMemoryStoresOptions{Workdir: workdir, Logger: silentLogger})
	ctx := context.Background()
	require.NoError(t, stores.Download(ctx, fetchSession(t, client)))
	local := filepath.Join(workdir, "memory", "notes")

	writeLocal(t, local, "note.md", "v1")
	server.SetLoseUploadResponses(true)
	runSync(ctx, stores)
	server.SetLoseUploadResponses(false)
	require.Equal(t, "v1", server.FilesSnapshot()["note.md"])
	sent := updated("note.md", "v1", "v0")
	require.Equal(t, []string{sent, sent}, server.ReceivedSnapshot())

	writeLocal(t, local, "note.md", "v2")
	runSync(ctx, stores)

	require.Equal(t, "v2", server.FilesSnapshot()["note.md"])
	require.Equal(t, "v2", readLocal(t, local, "note.md"))
}

func TestLostResponse_ADeletionMadeAfterALostUploadResponseStillReachesTheServer(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"note.md": "v0", "keep.md": "v0"}, "", nil, nil)
	ctx := context.Background()

	writeLocal(t, local, "note.md", "v1")
	writeLocal(t, local, "new.md", "v1")
	server.SetLoseUploadResponses(true)
	runSync(ctx, stores)
	server.SetLoseUploadResponses(false)
	require.NoError(t, os.Remove(filepath.Join(local, "note.md")))
	require.NoError(t, os.Remove(filepath.Join(local, "new.md")))

	stores.Finish(ctx)

	// The server's copies are this sync's own uploads, so the deletions go
	// out, guarded by the content that was sent.
	require.Equal(t, map[string]string{"keep.md": "v0"}, server.FilesSnapshot())
	received := server.ReceivedSnapshot()
	require.Equal(t, []string{deletedReq("new.md", "v1"), deletedReq("note.md", "v1")}, received[len(received)-2:])
}

func TestLostResponse_AFolderWipedAfterALostCreateResponseIsRebuiltNotReadAsDeletions(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"note.md": "v0"}, "", nil, nil)
	ctx := context.Background()

	writeLocal(t, local, "new.md", "v1")
	server.SetLoseUploadResponses(true)
	runSync(ctx, stores)
	server.SetLoseUploadResponses(false)
	// Both files vanish at once; the marker stays.
	require.NoError(t, os.Remove(filepath.Join(local, "note.md")))
	require.NoError(t, os.Remove(filepath.Join(local, "new.md")))

	// The final sync waives the delete window, so a misread wipe would delete
	// at once.
	stores.Finish(ctx)

	// Two known files gone together is a wipe: re-downloaded, nothing deleted.
	require.Equal(t, map[string]string{"note.md": "v0", "new.md": "v1"}, server.FilesSnapshot())
	require.Equal(t, "v0", readLocal(t, local, "note.md"))
	require.Equal(t, "v1", readLocal(t, local, "new.md"))
	require.Equal(t, []string{created("new.md", "v1")}, server.ReceivedSnapshot())
}

func TestLostResponse_TheShutdownFlushPushesARevertMadeAfterALostUploadResponse(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"note.md": "v0"}, "", nil, nil)
	ctx := context.Background()

	writeLocal(t, local, "note.md", "v1")
	server.SetLoseUploadResponses(true)
	runSync(ctx, stores)
	server.SetLoseUploadResponses(false)
	// Back to the last content the sync knows it saved — but the server moved on.
	writeLocal(t, local, "note.md", "v0")

	stores.FlushWrites(ctx)

	require.Equal(t, "v0", server.FilesSnapshot()["note.md"])
	received := server.ReceivedSnapshot()
	require.Equal(t, updated("note.md", "v0", "v1"), received[len(received)-1])
}

func TestLostResponse_TheShutdownFlushSendsNothingForARevertWhoseUploadNeverLanded(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"note.md": "v0", "keep.md": "v0"}, "", nil, nil)
	ctx := context.Background()

	writeLocal(t, local, "note.md", "v1")
	server.FailUploads = 500
	runSync(ctx, stores)
	server.FailUploads = 0
	writeLocal(t, local, "note.md", "v0")
	// Another writer deletes the memory; the file holds no edit to save.
	server.Delete("note.md")

	stores.FlushWrites(ctx)

	require.Equal(t, map[string]string{"keep.md": "v0"}, server.FilesSnapshot())
	require.Empty(t, server.ReceivedSnapshot())
}

func TestLostResponse_TheShutdownFlushPushesAnEditMadeAfterALostUploadResponse(t *testing.T) {
	local, server, stores := downloadedStores(t, map[string]string{"note.md": "v0"}, "", nil, nil)
	ctx := context.Background()

	writeLocal(t, local, "note.md", "v1")
	server.SetLoseUploadResponses(true)
	runSync(ctx, stores)
	server.SetLoseUploadResponses(false)
	writeLocal(t, local, "note.md", "v2")

	stores.FlushWrites(ctx)

	require.Equal(t, "v2", server.FilesSnapshot()["note.md"])
	received := server.ReceivedSnapshot()
	require.Equal(t, updated("note.md", "v2", "v1"), received[len(received)-1])
}

func TestLostResponse_AnotherWritersEditAfterALostUploadResponseStillWins(t *testing.T) {
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"note.md": "v0"}, "", &logBuf, nil)
	ctx := context.Background()

	writeLocal(t, local, "note.md", "v1")
	server.SetLoseUploadResponses(true)
	runSync(ctx, stores)
	server.SetLoseUploadResponses(false)
	server.Write("note.md", "theirs")
	writeLocal(t, local, "note.md", "v2")

	runSync(ctx, stores)

	// The server's copy is not what this sync sent, so the usual conflict rule holds.
	require.Equal(t, "theirs", server.FilesSnapshot()["note.md"])
	require.Equal(t, "theirs", readLocal(t, local, "note.md"))
	require.Contains(t, logBuf.String(), "changed both locally and remotely")
}

func TestLostResponse_ASendIsForgottenOnceALaterSyncHasSeenTheServer(t *testing.T) {
	// A remembered send must not outlive the next completed sync: kept
	// longer, it would mistake another writer's identical bytes for this
	// sync's upload and overwrite them.
	var logBuf bytes.Buffer
	local, server, stores := downloadedStores(t, map[string]string{"note.md": "v0"}, "", &logBuf, nil)
	ctx := context.Background()

	writeLocal(t, local, "note.md", "v1")
	server.FailUploads = 500
	runSync(ctx, stores)
	server.FailUploads = 0
	require.Equal(t, "v0", server.FilesSnapshot()["note.md"])

	// The agent reverts, so the next sync has nothing to send — its listing
	// shows the v1 upload never landed.
	writeLocal(t, local, "note.md", "v0")
	runSync(ctx, stores)

	// Another writer now saves the very bytes that failed upload carried.
	server.Write("note.md", "v1")
	writeLocal(t, local, "note.md", "v2")
	runSync(ctx, stores)

	require.Equal(t, "v1", server.FilesSnapshot()["note.md"])
	require.Equal(t, "v1", readLocal(t, local, "note.md"))
	require.Contains(t, logBuf.String(), "changed both locally and remotely")
	require.Empty(t, server.ReceivedSnapshot())
}
