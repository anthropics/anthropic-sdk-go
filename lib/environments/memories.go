package environments

// Session-level memory-store download and sync.
//
// A session may have several memory stores attached. SessionMemoryStores
// resolves where each store's folder goes on disk, opens a confined
// filestore.FileStore there, and reconciles each folder with its remote
// store — the merge rules live on Sync.

import (
	"cmp"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"maps"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/filestore"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"golang.org/x/sync/errgroup"
)

// MemoryDeleteMode is whether a locally deleted file may delete its
// memory on the server. [MemorySyncDeletionsLogOnly] runs the same checks and
// only logs what it would delete.
type MemoryDeleteMode int

const (
	MemorySyncDeletionsEnabled MemoryDeleteMode = iota
	MemorySyncDeletionsLogOnly
	MemorySyncDeletionsDisabled
)

// DefaultMemorySyncInterval is how often the worker syncs the session's
// memory stores back while the session runs. Checked after each dispatched
// tool call.
const DefaultMemorySyncInterval = 15 * time.Second

// MinMemorySyncInterval is the shortest sync interval accepted; a shorter
// positive interval is rejected with [ErrMemorySyncIntervalTooShort].
const MinMemorySyncInterval = 5 * time.Second

// ErrMemorySyncIntervalTooShort reports a positive sync interval below
// [MinMemorySyncInterval]. Zero (the default) and negative (disabled)
// intervals are never rejected.
var ErrMemorySyncIntervalTooShort = fmt.Errorf(
	"memory sync interval must be at least %s (or zero for the default, negative to disable)",
	MinMemorySyncInterval,
)

// MemoryFlushTimeout bounds each of the two network passes
// [SessionMemoryStores.Cleanup] runs — the Finish, then the
// [SessionMemoryStores.FlushWrites] — so a slow server cannot stall
// teardown. A pass that hits the bound logs a warning.
var MemoryFlushTimeout = 30 * time.Second

// MarkerPath is the reserved file at each store root whose content gates
// whether a sync trusts the folder. Never itself synced.
const MarkerPath = ".anthropic-memory-store"

const markerVersion = 1

func markerSHA(memoryStoreID string) string {
	sum := sha256.Sum256(fmt.Appendf(nil, "version %d\n%s", markerVersion, memoryStoreID))
	return hex.EncodeToString(sum[:])
}

// DeleteCorroborationWindow is how long a file must stay missing locally
// before its server delete is sent.
const DeleteCorroborationWindow = 30 * time.Second

// One sync sends at most a quarter of the store's files as deletes,
// clamped to [floor, ceiling] so small stores don't crawl.
const (
	deleteCapFloor   = 8
	deleteCapCeiling = 50
)

// Page sizes for memory listings — the API's maximum per view: basic pages
// carry up to 100 items, full pages are capped by the server at 20.
const (
	listPageSize     = 100
	fullListPageSize = 20
)

// fetchConcurrency is how many single-memory content fetches may be in
// flight at once during one store's pull pass. A sync rarely pulls more
// than a handful of memories, so a higher cap buys nothing in the common
// case.
const fetchConcurrency = 16

// uploadConcurrency is how many uploads one store's flush keeps in flight.
// At ~0.3s per upload, 32 clears the server's 2000-memories-per-store cap
// inside [MemoryFlushTimeout].
const uploadConcurrency = 32

// pendingPull is a memory the sha pass decided to write to disk; the
// content pass fetches and writes it.
type pendingPull struct {
	rel    string
	listed anthropic.BetaManagedAgentsMemory
}

// pendingUpload is a file the flush decided to push; uploadAll sends it.
type pendingUpload struct {
	rel      string
	localSHA string
	existing *anthropic.BetaManagedAgentsMemory
}

type deletePass struct {
	mode MemoryDeleteMode
	cap  int
	// True on the final sync, where waiting would drop the delete.
	waiveWindow bool
	attempted   int
	capped      int
	suppressed  int
}

func (d *deletePass) takeSlot() bool {
	if d.attempted >= d.cap {
		d.capped++
		return false
	}
	d.attempted++
	return true
}

// markerScan is one directory scan with the marker checked and removed —
// same scan, so a folder deleted mid-sync cannot pass the check and then
// read as empty.
type markerScan struct {
	files map[string]string
	// The only state in which the scan may drive uploads or deletes.
	markerOK       bool
	distrustReason string
}

// SessionMemoryError reports a memory store that could not be materialised on
// disk. Only [SessionMemoryStores.Download] returns it.
//
// Err is the cause — the refusal of a directory that already existed, a failed
// listing, an unusable root — and [errors.As]/[errors.Is] reach it through
// Unwrap.
type SessionMemoryError struct {
	// MemoryStoreID is the store that could not be materialised.
	MemoryStoreID string
	// Err is why it could not be. Never nil.
	Err error
}

func (e *SessionMemoryError) Error() string {
	return fmt.Sprintf("memory store %s: %v", e.MemoryStoreID, e.Err)
}

func (e *SessionMemoryError) Unwrap() error { return e.Err }

// ErrSessionMemoryNoToken reports a work item that carried no sessions token
// for a session with memory stores attached. The item fails, as a hosted
// sandbox does; the poller must issue a per-item secret carrying a sessions
// token, or the operator must disable memory sync with a negative interval.
var ErrSessionMemoryNoToken = errors.New(
	"the work item carried no sessions token, so the session's memory stores cannot be mounted",
)

// storeFiles is the slice of [filestore.FileStore] the sync uses; an
// interface so tests can inject disk failures.
type storeFiles interface {
	Put(rel string, data []byte, opts ...filestore.PutOption) error
	Get(rel string) ([]byte, error)
	HashTree(under string) (map[string]string, error)
	HashFile(rel string) (string, error)
	Remove(rel string) error
	CreateRoot() error
	Root() filestore.Root
	Dispose() error
}

func openFileStore(root string, opts ...filestore.Option) (storeFiles, error) {
	return filestore.Open(root, opts...)
}

// attachedStore is one attached store: its FileStore on disk plus the sync baseline.
type attachedStore struct {
	memoryStoreID string
	files         storeFiles
	readOnly      bool
	// mu guards baseline and refusedSHAs while a pass fans out its fetches
	// or uploads; the merge decisions before the fan-out run serially and
	// take no lock.
	mu sync.Mutex
	// baseline is {rel path → content sha} as of the last download or
	// successful sync.
	baseline map[string]string
	// refusedSHAs is {rel path → sha} the server refused; retried only after
	// the file changes.
	refusedSHAs map[string]string
	// When each path was first seen missing locally.
	pendingDeletes map[string]time.Time
}

// SessionMemoryStoresOptions configures a [SessionMemoryStores].
type SessionMemoryStoresOptions struct {
	// Workdir anchors the {workdir}/memory/<name> fallback for stores that
	// carry no mount path. Required.
	Workdir string
	// SyncInterval is the minimum gap between SyncIfDue syncs. Zero uses
	// [DefaultMemorySyncInterval]; a positive value below
	// [MinMemorySyncInterval] is rejected.
	SyncInterval time.Duration
	// SyncDeletions is whether a file deleted locally may delete its memory
	// on the server; the zero value is [MemorySyncDeletionsEnabled].
	SyncDeletions MemoryDeleteMode
	// RequestOptions are applied to every request — the memory endpoints
	// reject the environment key, so the worker passes its per-item
	// token-scoped auth options here.
	RequestOptions []option.RequestOption
	// Logger receives non-fatal warnings. Defaults to slog.Default().
	Logger *slog.Logger
}

// SessionMemoryStores is the memory stores attached to one session,
// materialised on disk.
//
// [SessionMemoryStores.Download] opens a confined store at each attached
// store's directory (its mount path, or a workdir fallback — see Download),
// pulls its memories, and records each one's content_sha256 as the sync
// baseline. Each sync ([SessionMemoryStores.SyncIfDue] on the worker's
// cadence, [SessionMemoryStores.Finish] once at the end) reconciles disk
// against server, per store and per path:
//
//   - a memory changed only remotely is written to disk;
//   - a file changed only locally is uploaded — an update with a
//     content_sha256 precondition, or a create for a new file;
//   - a file changed on both sides logs a warning and takes the server
//     version;
//   - a file the server refuses (too large, invalid content) is skipped —
//     warned once and retried only after the file changes; other files keep
//     syncing;
//   - a file deleted locally is deleted on the server, guarded by an
//     expected_content_sha256 precondition so a concurrent server-side
//     edit wins and restores the file. The delete waits
//     [DeleteCorroborationWindow] for a later sync to confirm the file is
//     still gone (Finish skips the wait), and each sync sends a bounded
//     number. [MemoryDeleteMode] gates whether deletes are sent at all;
//   - a memory deleted on the server is deleted on disk — unless the
//     local file holds an un-pushed edit, in which case a writable
//     store re-creates the memory from it and a read-only store keeps
//     the file unsynced;
//   - a store attached read-only pulls but never pushes.
//
// A download pulls the whole store, so it lists with content included. The
// recurring syncs instead run two phases: a content-free listing (paths and
// shas) drives the merge decisions, then only the memories actually being
// written to disk are fetched, a bounded number at a time. A sync that finds
// nothing changed moves no content at all.
//
// A file whose write to disk failed is never in the baseline, so its absence
// reads as a failed download — it is pulled again, never deleted. A write
// never re-creates a store folder that vanished mid-sync: it fails, and the
// next sync's scan finds whatever is at the path by then — nothing
// (re-downloaded) or someone else's files (left alone) — under the rules
// below.
//
// Download writes a [MarkerPath] file into each folder, and sync only acts
// on a folder whose marker matches the store. An empty or missing folder is
// rebuilt from the server. A folder with files but no matching marker is
// left untouched: nothing is uploaded from it, and nothing is deleted on
// the server because of it.
//
// The syncs never fail: one bad store or file is logged and the rest
// continue. Instances are not safe for concurrent use. Lifecycle:
// Download, SyncIfDue per tool call, Cleanup at the end.
type SessionMemoryStores struct {
	client        anthropic.Client
	workdir       string
	syncInterval  time.Duration
	syncDeletions MemoryDeleteMode
	reqOpts       []option.RequestOption
	log           *slog.Logger
	lastSync      time.Time
	finished      bool
	stores        []*attachedStore

	// now is the time source SyncIfDue measures the gap with; setClock
	// swaps it (the worker's test hook and this package's tests).
	now func() time.Time
	// openStore opens the confined folder for one store's root; tests swap
	// it to inject disk failures.
	openStore func(root string, opts ...filestore.Option) (storeFiles, error)
}

// checkMemorySyncInterval rejects a positive interval below the floor.
func checkMemorySyncInterval(d time.Duration) error {
	if d > 0 && d < MinMemorySyncInterval {
		return fmt.Errorf("%w: got %s", ErrMemorySyncIntervalTooShort, d)
	}
	return nil
}

// NewSessionMemoryStores returns a [SessionMemoryStores] bound to client.
// Call [SessionMemoryStores.Download] to materialise the stores. It fails
// only on an invalid SyncInterval ([ErrMemorySyncIntervalTooShort]).
func NewSessionMemoryStores(client anthropic.Client, opts SessionMemoryStoresOptions) (*SessionMemoryStores, error) {
	if err := checkMemorySyncInterval(opts.SyncInterval); err != nil {
		return nil, err
	}
	if opts.SyncInterval == 0 {
		opts.SyncInterval = DefaultMemorySyncInterval
	}
	if opts.Logger == nil {
		opts.Logger = slog.Default()
	}
	return &SessionMemoryStores{
		client:        client,
		workdir:       opts.Workdir,
		syncInterval:  opts.SyncInterval,
		syncDeletions: opts.SyncDeletions,
		reqOpts:       opts.RequestOptions,
		log:           opts.Logger,
		now:           time.Now,
		lastSync:      time.Now(),
		openStore:     openFileStore,
	}, nil
}

// setClock replaces the time source and re-seeds the last-sync mark so the
// cadence math runs entirely on the new clock.
func (s *SessionMemoryStores) setClock(now func() time.Time) {
	s.now = now
	s.lastSync = now()
}

// Roots returns every attached store's folder. The worker lists these in the
// file tools' AllowedRoots so a store mounted outside the working directory
// stays reachable.
func (s *SessionMemoryStores) Roots() []string {
	var roots []string
	for _, st := range s.stores {
		roots = append(roots, st.files.Root().Path)
	}
	return roots
}

// ReadOnlyRoots returns the root directories of stores attached read-only.
// The file tools consult this to refuse writes into read-only stores.
func (s *SessionMemoryStores) ReadOnlyRoots() []string {
	var roots []string
	for _, st := range s.stores {
		if st.readOnly {
			roots = append(roots, st.files.Root().Path)
		}
	}
	return roots
}

// Download downloads every attached store's memories to disk.
//
// session arrives already fetched — one snapshot shared with the skills
// download, so the two cannot disagree about the resources.
//
// Fails on the first store it cannot materialise, returning a
// [SessionMemoryError] — a session missing a folder its system prompt names
// runs with amnesia and syncs nothing back. Stores that did land stay
// tracked, so the caller's Dispose still removes the directories this created.
func (s *SessionMemoryStores) Download(ctx context.Context, session *anthropic.BetaManagedAgentsSession) error {
	for _, res := range session.Resources {
		if res.Type != "memory_store" {
			continue
		}
		resource := res.AsMemoryStore()
		root, err := s.storeRoot(resource)
		if err != nil {
			return &SessionMemoryError{MemoryStoreID: resource.MemoryStoreID, Err: err}
		}
		if err := s.downloadStore(ctx, resource, root); err != nil {
			return &SessionMemoryError{MemoryStoreID: resource.MemoryStoreID, Err: err}
		}
	}
	s.lastSync = s.now()
	return nil
}

// storeRoot is where one store's files land on disk.
//
// The store's files land at its MountPath — the very location the agent's
// system prompt tells it to read. A MountPath we cannot use verbatim is refused
// rather than quietly relocated: the agent would read an empty folder at the
// path it was told about, and write notes somewhere the next session looks for
// nothing.
func (s *SessionMemoryStores) storeRoot(resource anthropic.BetaManagedAgentsMemoryStoreResource) (string, error) {
	if root := resource.MountPath; root != "" {
		if !filestore.IsPathLegal(root) {
			return "", fmt.Errorf("mount_path is not a clean absolute path: %q", root)
		}
		return root, nil
	}
	// No mount path at all: nothing points the agent anywhere, so the workdir
	// is as good a home as any — {workdir}/memory/<name>, or the store id when
	// there is no name.
	return filepath.Join(s.workdir, "memory", cmp.Or(resource.Name, resource.MemoryStoreID)), nil
}

func (s *SessionMemoryStores) downloadStore(ctx context.Context, resource anthropic.BetaManagedAgentsMemoryStoreResource, root string) error {
	// WithUTF8: a binary file is refused at Put/Get, not mid-sync.
	files, err := s.openStore(root, filestore.WithUTF8())
	if err != nil {
		return err
	}
	// A root Open did not create is a dead run's leftovers; the first sync
	// would upload them into the customer's store. Left exactly as found.
	if !files.Root().RemovedOnDispose {
		// The configured path, not Root().Path — log the operator's own
		// string, not our normalization of it.
		return fmt.Errorf("something already exists at the memory store's path: %s; it must not exist when the session starts", root)
	}
	if err := files.CreateRoot(); err != nil {
		// An unmountable root fails the item, not just one file.
		_ = files.Dispose() // removes the root iff the failed mkdir got that far
		return fmt.Errorf("cannot create the memory store's folder: %s: %w; the worker host must make this mount path writable", root, err)
	}
	store := &attachedStore{
		memoryStoreID:  resource.MemoryStoreID,
		files:          files,
		readOnly:       resource.Access == anthropic.BetaManagedAgentsMemoryStoreResourceAccessReadOnly,
		baseline:       map[string]string{},
		refusedSHAs:    map[string]string{},
		pendingDeletes: map[string]time.Time{},
	}
	if err := s.stampAndPull(ctx, store); err != nil {
		// A half-downloaded folder self-destructs; Dispose removes only what
		// this Open created — the guard above already established it did.
		_ = files.Dispose()
		return err
	}
	s.log.Info("downloaded memories",
		slog.Int("count", len(store.baseline)),
		slog.String("memory_store_id", store.memoryStoreID),
		slog.String("dir", store.files.Root().Path))
	s.stores = append(s.stores, store)
	return nil
}

func (s *SessionMemoryStores) scanMarker(store *attachedStore) (markerScan, error) {
	local, err := store.files.HashTree("/")
	if err != nil {
		return markerScan{}, err
	}
	marker, hasMarker := local[MarkerPath]
	delete(local, MarkerPath)
	if hasMarker && marker == markerSHA(store.memoryStoreID) {
		return markerScan{files: local, markerOK: true}, nil
	}
	reason := "the marker file is gone"
	if hasMarker {
		reason = "the marker file does not match this store"
	}
	return markerScan{files: local, markerOK: false, distrustReason: reason}, nil
}

// rebuild rebuilds a destroyed folder from the server; sends no deletes, no uploads.
func (s *SessionMemoryStores) rebuild(ctx context.Context, store *attachedStore, reason string) error {
	s.log.Warn(reason+"; re-downloading the memory store folder instead of syncing",
		slog.String("root", store.files.Root().Path),
		slog.String("memory_store_id", store.memoryStoreID))
	if err := store.files.CreateRoot(); err != nil {
		return err
	}
	return s.stampAndPull(ctx, store)
}

// stampAndPull writes the trust marker, then pulls every remote memory;
// pushes nothing. The baseline is rebuilt from only the writes that
// succeed, so it never holds a file that isn't on disk. Every memory is
// needed here, so the listing carries the content — pages cost far fewer
// round-trips than a request per memory.
func (s *SessionMemoryStores) stampAndPull(ctx context.Context, store *attachedStore) error {
	store.baseline = map[string]string{}
	clear(store.pendingDeletes) // earlier absence observations no longer describe this disk
	if err := store.files.Put(MarkerPath,
		fmt.Appendf(nil, "version %d\n%s", markerVersion, store.memoryStoreID)); err != nil {
		return err
	}
	for item, err := range s.listMemories(ctx, store.memoryStoreID, anthropic.BetaManagedAgentsMemoryViewFull) {
		if err != nil {
			return err
		}
		rel := strings.TrimLeft(item.Path, "/")
		if s.write(store, rel, item.Content) {
			store.baseline[rel] = item.ContentSha256
		}
	}
	return nil
}

// listMemories iterates the store's memories at the largest page size the
// view allows — basic (shas, no content) at [listPageSize], full at
// [fullListPageSize]. Prefix rollups and the marker path are skipped.
func (s *SessionMemoryStores) listMemories(ctx context.Context, memoryStoreID string, view anthropic.BetaManagedAgentsMemoryView) iter.Seq2[anthropic.BetaManagedAgentsMemory, error] {
	limit := int64(listPageSize)
	if view == anthropic.BetaManagedAgentsMemoryViewFull {
		limit = fullListPageSize
	}
	return func(yield func(anthropic.BetaManagedAgentsMemory, error) bool) {
		pager := s.client.Beta.MemoryStores.Memories.ListAutoPaging(ctx, memoryStoreID,
			anthropic.BetaMemoryStoreMemoryListParams{View: view, Limit: param.NewOpt(limit)}, s.reqOpts...)
		for pager.Next() {
			item := pager.Current()
			if item.Type != "memory" {
				continue
			}
			if strings.TrimLeft(item.Path, "/") == MarkerPath {
				s.log.Warn("the server listed the reserved marker path; skipping",
					slog.String("path", item.Path), slog.String("memory_store_id", memoryStoreID))
				continue
			}
			if !yield(item.AsMemory(), nil) {
				return
			}
		}
		if err := pager.Err(); err != nil {
			yield(anthropic.BetaManagedAgentsMemory{}, err)
		}
	}
}

// Cleanup is the teardown, run once the session's tools are closed: the
// last full sync when the session ended cleanly, a flush of unsynced
// writes always, then the folders are removed. Each network step gets its
// own [MemoryFlushTimeout]-bounded context — the session context may
// already be cancelled — and a step that hits the bound logs a warning.
func (s *SessionMemoryStores) Cleanup(cleanEnd bool) {
	if cleanEnd {
		ctx, cancel := context.WithTimeout(context.Background(), MemoryFlushTimeout)
		s.Finish(ctx)
		if ctx.Err() != nil {
			s.log.Warn(fmt.Sprintf("final memory sync cut off after %s; the flush that follows still uploads changed files",
				MemoryFlushTimeout))
		}
		cancel()
	}
	// A fresh bound of its own even after Finish: Finish may have hit its
	// bound or swallowed a failure, and after a clean Finish the flush
	// finds nothing dirty and makes no network calls.
	ctx, cancel := context.WithTimeout(context.Background(), MemoryFlushTimeout)
	s.FlushWrites(ctx)
	if ctx.Err() != nil {
		s.log.Warn(fmt.Sprintf("memory flush cut off after %s; changed files it had not uploaded yet are not saved",
			MemoryFlushTimeout))
	}
	cancel()
	s.Dispose()
}

// Finish is the session's last sync — call once, at a clean end. Server
// deletes skip the [DeleteCorroborationWindow] wait. Panics if called twice.
func (s *SessionMemoryStores) Finish(ctx context.Context) {
	if s.finished {
		panic("Finish() was already called: it is the session's last sync and runs once")
	}
	s.finished = true
	s.syncAll(ctx, true)
}

func (s *SessionMemoryStores) syncAll(ctx context.Context, final bool) {
	var wg sync.WaitGroup
	for _, store := range s.stores {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.syncStore(ctx, store, final); err != nil {
				s.log.Warn("memory sync failed",
					slog.String("memory_store_id", store.memoryStoreID), slog.Any("error", err))
			}
		}()
	}
	wg.Wait()
	s.lastSync = s.now()
}

// SyncIfDue runs one sync when the interval has elapsed since the last.
// Never returns an error — a store that fails is logged and the rest
// continue.
func (s *SessionMemoryStores) SyncIfDue(ctx context.Context) {
	if s.now().Sub(s.lastSync) < s.syncInterval {
		return
	}
	s.syncAll(ctx, false)
}

// FlushWrites uploads new and changed files and does nothing else — no
// server deletes, no local removals, no pulls.
//
// The shutdown rescue for a session that ended on an error or cancel:
// best-effort, bounded by the caller, so an errored session can still
// lose edits. Each store uploads up to [uploadConcurrency] files at a
// time; a store cut off part-way logs how many changed files it had not
// uploaded. Read-only stores, refused files, files the server already
// holds, and folders that fail the marker check are skipped. Never
// returns an error.
func (s *SessionMemoryStores) FlushWrites(ctx context.Context) {
	var wg sync.WaitGroup
	for _, store := range s.stores {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.flushStore(ctx, store); err != nil {
				s.log.Warn("memory flush failed",
					slog.String("memory_store_id", store.memoryStoreID), slog.Any("error", err))
			}
		}()
	}
	wg.Wait()
}

func (s *SessionMemoryStores) flushStore(ctx context.Context, store *attachedStore) error {
	if store.readOnly {
		return nil
	}
	scan, err := s.scanMarker(store)
	if err != nil {
		return err
	}
	if !scan.markerOK {
		s.log.Warn(scan.distrustReason+"; not uploading anything from the memory store folder",
			slog.String("root", store.files.Root().Path),
			slog.String("memory_store_id", store.memoryStoreID))
		return nil
	}
	dirty := map[string]string{}
	for rel, sha := range scan.files {
		if sha != store.baseline[rel] && store.refusedSHAs[rel] != sha {
			dirty[rel] = sha
		}
	}
	if len(dirty) == 0 {
		return nil
	}
	unfinished := len(dirty)
	defer func() {
		if ctx.Err() != nil && unfinished > 0 {
			s.log.Warn(fmt.Sprintf("memory flush cut off part-way; %d of %d changed files had not finished uploading",
				unfinished, len(dirty)),
				slog.String("memory_store_id", store.memoryStoreID))
		}
	}()
	remote := map[string]anthropic.BetaManagedAgentsMemory{}
	for item, err := range s.listMemories(ctx, store.memoryStoreID, anthropic.BetaManagedAgentsMemoryViewBasic) {
		if err != nil {
			if ctx.Err() != nil {
				return nil // reported by the cut-off line, not as a failure
			}
			return err
		}
		remote[strings.TrimLeft(item.Path, "/")] = item
	}
	var uploads []pendingUpload
	for _, rel := range slices.Sorted(maps.Keys(dirty)) {
		localSHA := dirty[rel]
		baseSHA, hasBase := store.baseline[rel]
		var existing *anthropic.BetaManagedAgentsMemory
		if item, ok := remote[rel]; ok {
			existing = &item
		}
		if existing != nil && existing.ContentSha256 == localSHA {
			store.baseline[rel] = localSHA
			continue
		}
		if existing != nil && (!hasBase || existing.ContentSha256 != baseSHA) {
			s.log.Warn("memory changed both locally and remotely; the flush leaves the remote version",
				slog.String("path", rel), slog.String("memory_store_id", store.memoryStoreID))
			continue
		}
		uploads = append(uploads, pendingUpload{rel: rel, localSHA: localSHA, existing: existing})
	}
	unfinished = s.uploadAll(ctx, store, uploads)
	return nil
}

// Dispose removes every store directory that Download created. A folder
// that holds files but fails the marker check is kept — sync promised to
// leave it as found.
func (s *SessionMemoryStores) Dispose() {
	for _, store := range s.stores {
		root := store.files.Root()
		scan, err := s.scanMarker(store)
		if err == nil && !scan.markerOK && len(scan.files) > 0 {
			s.log.Warn(scan.distrustReason+"; leaving the memory store folder on disk",
				slog.String("root", root.Path),
				slog.String("memory_store_id", store.memoryStoreID))
			continue
		}
		if err == nil {
			err = store.files.Dispose()
		}
		if err != nil {
			s.log.Warn("failed to remove the memory store folder",
				slog.String("root", root.Path),
				slog.String("memory_store_id", store.memoryStoreID), slog.Any("error", err))
			continue
		}
		if root.RemovedOnDispose {
			s.log.Info("removed memory store dir",
				slog.String("dir", root.Path), slog.String("memory_store_id", store.memoryStoreID))
		}
	}
}

func (s *SessionMemoryStores) syncStore(ctx context.Context, store *attachedStore, final bool) error {
	scan, err := s.scanMarker(store)
	if err != nil {
		return err
	}
	local := scan.files
	if !scan.markerOK {
		if len(local) > 0 {
			s.log.Warn(scan.distrustReason+"; leaving the memory store folder as found and not syncing. "+
				"Delete it to recover — the server copy is intact.",
				slog.String("root", store.files.Root().Path),
				slog.String("memory_store_id", store.memoryStoreID))
			return nil
		}
		return s.rebuild(ctx, store, "the folder or its marker is gone")
	}
	// A lone file vanishing is an ordinary deletion; two or more at once
	// with nothing left is a wiped folder.
	if len(local) == 0 && len(store.baseline) > 1 {
		return s.rebuild(ctx, store, "every memory file is gone at once")
	}

	remote := map[string]anthropic.BetaManagedAgentsMemory{}
	for item, err := range s.listMemories(ctx, store.memoryStoreID, anthropic.BetaManagedAgentsMemoryViewBasic) {
		if err != nil {
			return err
		}
		remote[strings.TrimLeft(item.Path, "/")] = item
	}

	deletes := &deletePass{
		mode:        s.syncDeletions,
		cap:         max(deleteCapFloor, min(deleteCapCeiling, len(store.baseline)/4)),
		waiveWindow: final,
	}
	var pulls []pendingPull
	baseline := map[string]string{}
	for _, rel := range sortedKeyUnion(remote, local, store.baseline) {
		var remoteItem *anthropic.BetaManagedAgentsMemory
		if item, ok := remote[rel]; ok {
			remoteItem = &item
		}
		localSHA, localPresent := local[rel]
		baseSHA, hasBase := store.baseline[rel]
		// Locally deleted, remotely unchanged: a candidate server delete.
		if !localPresent && hasBase && remoteItem != nil && remoteItem.ContentSha256 == baseSHA && !store.readOnly {
			if sha, ok := s.corroboratedDelete(ctx, store, rel, remoteItem, baseSHA, deletes); ok {
				baseline[rel] = sha
			}
			continue
		}
		if sha, ok := s.syncPath(ctx, store, rel, remoteItem, localSHA, localPresent, &pulls); ok {
			baseline[rel] = sha
		}
	}
	store.baseline = baseline
	// The content pass: everything above moved only shas.
	s.pullAll(ctx, store, pulls)
	if deletes.suppressed > 0 {
		s.log.Debug("server deletes are disabled; locally deleted memories stay on the server",
			slog.Int("count", deletes.suppressed), slog.String("memory_store_id", store.memoryStoreID))
	}
	if deletes.capped > 0 {
		verb := "sent"
		if deletes.mode == MemorySyncDeletionsLogOnly {
			verb = "would send"
		}
		s.log.Warn(fmt.Sprintf("delete cap reached: %s %d deletes, held %d for later syncs",
			verb, deletes.attempted, deletes.capped),
			slog.String("memory_store_id", store.memoryStoreID))
	}
	return nil
}

// sortedKeyUnion returns the sorted union of the three maps' keys.
// Deterministic order keeps each store's logs and server traffic stable
// across runs (stores sync in parallel, so cross-store order is not).
func sortedKeyUnion(remote map[string]anthropic.BetaManagedAgentsMemory, local, base map[string]string) []string {
	keys := slices.Collect(maps.Keys(remote))
	keys = slices.AppendSeq(keys, maps.Keys(local))
	keys = slices.AppendSeq(keys, maps.Keys(base))
	slices.Sort(keys)
	return slices.Compact(keys)
}

// syncPath reconciles a single path and returns its entry in the next
// sync's baseline. The baseline advances only on a successful disk write
// or server call. A path that needs the remote content is appended to
// pulls and keeps its old baseline entry until the content pass writes it.
func (s *SessionMemoryStores) syncPath(ctx context.Context, store *attachedStore, rel string, remote *anthropic.BetaManagedAgentsMemory, localSHA string, localPresent bool, pulls *[]pendingPull) (string, bool) {
	baseSHA, hasBase := store.baseline[rel]
	if localPresent {
		delete(store.pendingDeletes, rel)
	}

	if remote == nil {
		if !localPresent {
			delete(store.pendingDeletes, rel)
			return "", false
		}
		if hasBase {
			if localSHA == baseSHA {
				fresh, onDisk := s.removeLocal(store, rel, baseSHA)
				if !onDisk {
					return "", false
				}
				if fresh == baseSHA {
					return baseSHA, true
				}
				localSHA = fresh
			}
			// The local edit is the only copy: fall through as a
			// local-only file (uploaded on a writable store, kept on a
			// read-only one).
			if store.readOnly {
				s.log.Warn("memory deleted remotely but edited locally; keeping the file, "+
					"which a read-only store cannot push",
					slog.String("path", rel), slog.String("memory_store_id", store.memoryStoreID))
			} else if store.refusedSHAs[rel] != localSHA {
				s.log.Info("memory deleted remotely but edited locally; re-creating it from the file",
					slog.String("path", rel), slog.String("memory_store_id", store.memoryStoreID))
			}
		}
		if store.readOnly {
			return "", false
		}
		if store.refusedSHAs[rel] == localSHA {
			return "", false // refused before, unchanged since
		}
		return s.upload(ctx, store, rel, localSHA, nil)
	}

	remoteSHA := remote.ContentSha256
	remoteChanged := !hasBase || remoteSHA != baseSHA
	locallyEdited := localPresent && (!hasBase || localSHA != baseSHA) && localSHA != remoteSHA
	// Read-only stores never push, so their local edits don't count.
	localChanged := !store.readOnly && locallyEdited

	if !localPresent && hasBase {
		// Only successful writes enter the baseline, so this file was
		// verifiably on disk and is now gone: a real local deletion.
		if remoteChanged {
			s.log.Warn("memory deleted locally but changed remotely; restoring the remote version",
				slog.String("path", rel), slog.String("memory_store_id", store.memoryStoreID))
			delete(store.pendingDeletes, rel)
			*pulls = append(*pulls, pendingPull{rel: rel, listed: *remote})
		}
		return baseSHA, true
	}

	switch {
	case remoteChanged:
		if localPresent && localSHA == remoteSHA {
			// The file already holds the remote bytes — adopt without a fetch.
			return remoteSHA, true
		}
		// Read-only stores too: overwriting a local edit deserves a trace.
		if locallyEdited {
			s.log.Warn("memory changed both locally and remotely; keeping the remote version",
				slog.String("path", rel), slog.String("memory_store_id", store.memoryStoreID))
		}
		*pulls = append(*pulls, pendingPull{rel: rel, listed: *remote})
		return baseSHA, hasBase
	case localChanged:
		if store.refusedSHAs[rel] == localSHA {
			return remoteSHA, true // refused before, unchanged since
		}
		if sha, ok := s.upload(ctx, store, rel, localSHA, remote); ok {
			return sha, true
		}
		return remoteSHA, true
	default:
		return remoteSHA, true
	}
}

// removeLocal removes the file for a memory the server no longer has, if
// it still holds expectSHA. Returns the file's sha after the call and
// whether it is still on disk: an edited file is left in place, and an
// I/O error is treated as still present so the next sync retries.
func (s *SessionMemoryStores) removeLocal(store *attachedStore, rel, expectSHA string) (freshSHA string, onDisk bool) {
	// The scan may be stale — re-read so a mid-sync edit isn't destroyed.
	fresh, err := store.files.HashFile(rel)
	if err != nil {
		// Dropping the entry would re-create the memory the server just deleted.
		return expectSHA, true
	}
	if fresh == "" {
		return "", false
	}
	if fresh != expectSHA {
		return fresh, true
	}
	if err := store.files.Remove(rel); err != nil {
		s.log.Warn("failed to remove memory deleted remotely",
			slog.String("path", rel), slog.String("memory_store_id", store.memoryStoreID), slog.Any("error", err))
		return expectSHA, true
	}
	return "", false
}

// write puts a memory's content on disk; false (and a warning) on failure.
// A ".." component in the wire path reaches here as a *FileStoreError — that
// is the escape guard.
func (s *SessionMemoryStores) write(store *attachedStore, rel, content string) bool {
	if err := store.files.Put(rel, []byte(content)); err != nil {
		s.log.Warn("failed to write memory",
			slog.String("path", rel), slog.String("memory_store_id", store.memoryStoreID), slog.Any("error", err))
		return false
	}
	return true
}

// pullAll fetches and writes the given memories, [fetchConcurrency] at a
// time — the sync's content pass. The listing carried no content, so each
// memory is fetched individually and written as it arrives. On success the
// path's baseline advances; on a failed fetch or write the old entry stays
// and the next sync retries. A 404 means the memory was deleted after the
// listing — the next sync reconciles it.
func (s *SessionMemoryStores) pullAll(ctx context.Context, store *attachedStore, pulls []pendingPull) {
	if len(pulls) == 0 {
		return
	}
	var g errgroup.Group
	g.SetLimit(fetchConcurrency)
	for _, p := range pulls {
		g.Go(func() error {
			// The write stays inside the limited goroutine so a slow disk
			// cannot let fetched bodies pile up beyond the concurrency bound.
			item, err := s.client.Beta.MemoryStores.Memories.Get(ctx, p.listed.ID,
				anthropic.BetaMemoryStoreMemoryGetParams{
					MemoryStoreID: store.memoryStoreID,
					View:          anthropic.BetaManagedAgentsMemoryViewFull,
				}, s.reqOpts...)
			if err != nil {
				if !isStatus(err, 404) {
					s.log.Warn("failed to fetch memory content",
						slog.String("path", p.rel), slog.String("memory_store_id", store.memoryStoreID), slog.Any("error", err))
				}
				return nil
			}
			if s.write(store, p.rel, item.Content) {
				store.mu.Lock()
				store.baseline[p.rel] = item.ContentSha256
				store.mu.Unlock()
			}
			return nil
		})
	}
	_ = g.Wait()
}

// uploadAll pushes the given files, [uploadConcurrency] at a time — the
// flush's send pass — and returns how many had not finished when ctx was
// done. A finished upload advances the path's baseline; one that failed on
// its own is settled too (upload logged it); none starts once ctx is done.
func (s *SessionMemoryStores) uploadAll(ctx context.Context, store *attachedStore, uploads []pendingUpload) int {
	unfinished := len(uploads)
	var g errgroup.Group
	g.SetLimit(uploadConcurrency)
	for _, u := range uploads {
		g.Go(func() error {
			if ctx.Err() != nil {
				return nil
			}
			sha, ok := s.upload(ctx, store, u.rel, u.localSHA, u.existing)
			if !ok && ctx.Err() != nil {
				return nil
			}
			store.mu.Lock()
			defer store.mu.Unlock()
			unfinished--
			if ok {
				store.baseline[u.rel] = sha
			}
			return nil
		})
	}
	_ = g.Wait()
	return unfinished
}

// upload pushes one local file. On failure the path stays out of the
// baseline, so the next pass retries.
//
// A refusal the server would repeat (400/413, the utf-8 gate) enters
// refusedSHAs: warned once, retried only after the file changes.
func (s *SessionMemoryStores) upload(ctx context.Context, store *attachedStore, rel, localSHA string, existing *anthropic.BetaManagedAgentsMemory) (string, bool) {
	data, err := store.files.Get(rel)
	if err != nil {
		var fse *filestore.FileStoreError
		if errors.As(err, &fse) && localSHA != "" {
			// The utf-8 gate: the same bytes would be refused again.
			store.mu.Lock()
			store.refusedSHAs[rel] = localSHA
			store.mu.Unlock()
			s.log.Warn("the server rejected this memory file, so it stays un-synced until its content changes",
				slog.String("path", rel), slog.String("memory_store_id", store.memoryStoreID), slog.Any("rejection", err))
			return "", false
		}
		s.log.Warn("failed to upload memory",
			slog.String("path", rel), slog.String("memory_store_id", store.memoryStoreID), slog.Any("error", err))
		return "", false
	}
	if data == nil {
		// The file vanished between HashTree and Get; the next sync sees the
		// deletion.
		return "", false
	}
	var item *anthropic.BetaManagedAgentsMemory
	if existing == nil {
		item, err = s.client.Beta.MemoryStores.Memories.New(ctx, store.memoryStoreID,
			anthropic.BetaMemoryStoreMemoryNewParams{
				Path:    "/" + rel,
				Content: param.NewOpt(string(data)),
			}, s.reqOpts...)
	} else {
		item, err = s.client.Beta.MemoryStores.Memories.Update(ctx, existing.ID,
			anthropic.BetaMemoryStoreMemoryUpdateParams{
				MemoryStoreID: store.memoryStoreID,
				Content:       param.NewOpt(string(data)),
				Precondition: anthropic.BetaManagedAgentsPreconditionParam{
					Type:          anthropic.BetaManagedAgentsPreconditionTypeContentSha256,
					ContentSha256: param.NewOpt(existing.ContentSha256),
				},
			}, s.reqOpts...)
	}
	if err != nil {
		switch {
		case ctx.Err() != nil:
			// The pass was cancelled or hit its bound — not this file's
			// failure, so no per-file line and no refusal recorded.
		case existing != nil && isStatus(err, 404):
			// Deleted remotely since the listing, so this file is now the only copy.
			return s.upload(ctx, store, rel, localSHA, nil)
		case existing != nil && isStatus(err, 409):
			// The precondition lost a race: the remote moved under us, so the
			// push is dropped. The local file is now stale — the next sync
			// sees remoteChanged and pulls the winner over it.
			s.log.Warn("memory changed both locally and remotely; the upload was refused and "+
				"the local edit loses",
				slog.String("path", rel), slog.String("memory_store_id", store.memoryStoreID))
		case (isStatus(err, 400) || isStatus(err, 413)) && localSHA != "":
			// The same bytes would be refused again.
			store.mu.Lock()
			store.refusedSHAs[rel] = localSHA
			store.mu.Unlock()
			s.log.Warn("the server rejected this memory file, so it stays un-synced until its content changes",
				slog.String("path", rel), slog.String("memory_store_id", store.memoryStoreID), slog.Any("rejection", err))
		default:
			s.log.Warn("failed to upload memory",
				slog.String("path", rel), slog.String("memory_store_id", store.memoryStoreID), slog.Any("error", err))
		}
		return "", false
	}
	store.mu.Lock()
	delete(store.refusedSHAs, rel)
	store.mu.Unlock()
	return item.ContentSha256, true
}

// corroboratedDelete records when the file was first seen absent and
// sends the server delete only after [DeleteCorroborationWindow], a fresh
// re-check that the file is still gone and the marker still intact,
// and the pass's cap.
func (s *SessionMemoryStores) corroboratedDelete(ctx context.Context, store *attachedStore, rel string, remote *anthropic.BetaManagedAgentsMemory, baseSHA string, deletes *deletePass) (string, bool) {
	if deletes.mode == MemorySyncDeletionsDisabled {
		deletes.suppressed++
		return baseSHA, true
	}
	if _, ok := store.pendingDeletes[rel]; !ok {
		store.pendingDeletes[rel] = s.now()
	}
	firstAbsent := store.pendingDeletes[rel]
	if !deletes.waiveWindow && s.now().Sub(firstAbsent) < DeleteCorroborationWindow {
		return baseSHA, true
	}
	// The folder may have been destroyed since this sync's own scan.
	markerNow, err := store.files.HashFile(MarkerPath)
	if err != nil || markerNow != markerSHA(store.memoryStoreID) {
		return baseSHA, true
	}
	fileNow, err := store.files.HashFile(rel)
	if err != nil {
		return baseSHA, true
	}
	if fileNow != "" {
		delete(store.pendingDeletes, rel)
		return baseSHA, true
	}
	if !deletes.takeSlot() {
		return baseSHA, true
	}
	if deletes.mode == MemorySyncDeletionsLogOnly {
		s.log.Info("log-only: sync would delete this memory on the server",
			slog.String("path", rel), slog.String("memory_store_id", store.memoryStoreID))
		return baseSHA, true
	}
	sha, ok := s.deleteRemote(ctx, store, rel, remote, baseSHA)
	if !ok {
		delete(store.pendingDeletes, rel)
	}
	return sha, ok
}

// deleteRemote propagates a local deletion, guarded by the last-synced
// content hash. A concurrent remote edit wins: the baseline entry is kept so
// the next sync restores the remote version instead of re-deleting.
func (s *SessionMemoryStores) deleteRemote(ctx context.Context, store *attachedStore, rel string, remote *anthropic.BetaManagedAgentsMemory, baseSHA string) (string, bool) {
	_, err := s.client.Beta.MemoryStores.Memories.Delete(ctx, remote.ID,
		anthropic.BetaMemoryStoreMemoryDeleteParams{
			MemoryStoreID:         store.memoryStoreID,
			ExpectedContentSha256: param.NewOpt(baseSHA),
		}, s.reqOpts...)
	if err != nil {
		switch {
		case isStatus(err, 404):
			return "", false // already gone remotely too
		case isStatus(err, 409), isStatus(err, 412):
			s.log.Warn("memory deleted locally but changed remotely; keeping the remote version",
				slog.String("path", rel), slog.String("memory_store_id", store.memoryStoreID))
		default:
			s.log.Warn("failed to delete memory",
				slog.String("path", rel), slog.String("memory_store_id", store.memoryStoreID), slog.Any("error", err))
		}
		return baseSHA, true
	}
	s.log.Info("propagated local deletion",
		slog.String("path", rel), slog.String("memory_store_id", store.memoryStoreID))
	return "", false
}
