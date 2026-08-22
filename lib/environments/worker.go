package environments

import (
	"cmp"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/sendwindow"
	"github.com/anthropics/anthropic-sdk-go/internal/stainlessheader"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/tools/agenttoolset"
)

const (
	heartbeatDefault    = 30 * time.Second
	heartbeatFloor      = 1 * time.Second
	noHeartbeatSentinel = "NO_HEARTBEAT"
)

// EnvironmentWorkerOptions configures an [EnvironmentWorker].
type EnvironmentWorkerOptions struct {
	// EnvironmentID is the self-hosted environment to poll for work. Required
	// by [EnvironmentWorker.Run]; unused by [EnvironmentWorker.HandleItem] (which
	// takes the work item's environment via [HandleItemOptions]).
	EnvironmentID string

	// EnvironmentKey is the worker's standing credential: polling always
	// uses it, and per-session calls fall back to it when a work item's
	// secret doesn't yield a sessions token. Required by
	// [EnvironmentWorker.Run]; for [EnvironmentWorker.HandleItem] see
	// [HandleItemOptions].
	EnvironmentKey string

	// WorkerID is a stable identifier reported back to the server for
	// observability. Defaults to "<os.Hostname()>-<random hex>" — forwarded
	// to [WorkPoller].
	WorkerID string

	// Workdir is the base directory for the per-session
	// [agenttoolset.AgentToolContext]. When empty it defaults to the process
	// working directory captured when [NewEnvironmentWorker] is called.
	Workdir string

	// Deprecated: no longer supported and slated for removal. The file tools
	// are always confined to Workdir plus the session's memory-store folders,
	// so a true value is not reinterpreted: [EnvironmentWorker.Run] and
	// [EnvironmentWorker.HandleItem] return
	// [agenttoolset.ErrUnrestrictedPathsUnsupported] instead.
	UnrestrictedPaths bool

	// MaxFileBytes is forwarded to the per-session
	// [agenttoolset.AgentToolContext], capping the size of files the read and
	// edit tools load into memory. Zero uses the built-in 256 KiB default; a
	// positive value sets a custom cap; a negative value disables the cap (use
	// only when the sandbox can absorb arbitrarily large files).
	MaxFileBytes int64

	// Tools, if non-nil, is exposed to every claimed session as-is. Ignored
	// when ToolsFunc is set. When both Tools and ToolsFunc are nil the worker
	// uses agenttoolset.BetaAgentToolset20260401(env) — the standard
	// agent_toolset_20260401 set bound to the per-session AgentToolContext. Tool
	// lifetime is the caller's responsibility; the worker never closes tools it
	// was given via Tools.
	Tools []anthropic.BetaTool

	// ToolsFunc, if non-nil, is invoked once per claimed session with that
	// session's [agenttoolset.AgentToolContext] — use it to bind
	// agenttoolset.BetaAgentToolset20260401 (or any tool that needs the workdir) to the
	// right session. The worker calls [agenttoolset.CloseAll] on the result
	// after the session finishes.
	ToolsFunc func(env *agenttoolset.AgentToolContext) []anthropic.BetaTool

	// MaxIdle is forwarded to the per-session
	// [github.com/anthropics/anthropic-sdk-go.SessionToolRunner].
	MaxIdle *time.Duration

	// MemorySyncInterval is how often the session's attached memory stores
	// are synced back while it runs — checked after each dispatched tool
	// call, plus one final sync when the session ends cleanly. A session
	// that ends on an error or cancel instead gets a push-only
	// [SessionMemoryStores.FlushWrites], bounded by [MemoryFlushTimeout] like
	// the final sync, with a warning logged when either bound cuts work off;
	// it can still lose edits. Zero uses [DefaultMemorySyncInterval]; a
	// positive value below [MinMemorySyncInterval] makes Run and HandleItem
	// fail up front with [ErrMemorySyncIntervalTooShort]; a negative value
	// disables memory download and sync entirely, silently.
	// Memory stores are only touched for work items whose secret carries a
	// sessions token, because the memory endpoints reject the environment
	// key. Without a token, a work item whose session has stores attached
	// fails — its memory cannot be mounted — unless the interval is
	// negative, in which case the session quietly runs without memory:
	// disabling sync is the operator's explicit choice.
	//
	// A store the worker cannot materialise on disk fails the work item too:
	// see [SessionMemoryStores.Download].
	MemorySyncInterval time.Duration

	// MemorySyncDeletions gates whether local deletions may delete server
	// memories; see [MemoryDeleteMode]. Uploads and pulls are unaffected.
	// The zero value is [MemorySyncDeletionsEnabled].
	MemorySyncDeletions MemoryDeleteMode

	// RequestOptions are applied to every request the worker issues, on top
	// of the environment-key auth and x-stainless-helper telemetry it adds
	// itself: the [WorkPoller]'s Poll/Ack/Stop, the lease heartbeat and
	// force-stop, the per-session skill download, and the SessionToolRunner's
	// event stream/list/send. Use it for a proxy/custom header or a base-URL
	// override that must reach the whole self-hosted runner. These options are
	// applied first, so the worker's own environment-key auth and helper
	// header take precedence and cannot be clobbered by a caller option.
	RequestOptions []option.RequestOption

	// Logger receives non-fatal warnings. Defaults to slog.Default().
	Logger *slog.Logger
}

// EnvironmentWorker is the self-hosted environment runner, composed from the
// control-plane [WorkPoller] and the per-session
// [github.com/anthropics/anthropic-sdk-go.SessionToolRunner].
//
// For each claimed `session` work item it builds the per-session
// [agenttoolset.AgentToolContext], downloads the session agent's skills
// ([agenttoolset.AgentToolContext.SetupSkillsFromSession]), then runs a SessionToolRunner for the
// session WHILE heartbeating the work-item lease in parallel; on exit it
// force-stops the work item (unless the lease was lost, in which case the item
// is left to whoever holds it now) and loops to the next one. The lease
// heartbeat reports state "stopping"/"stopped" or a lost lease back into the
// run by cancelling the session runner.
//
// [EnvironmentWorker.Run] drives the poll loop. [EnvironmentWorker.HandleItem]
// exposes the per-item flow for callers that already hold a claimed work item
// (for example a `worker poll --on-work` hook) and only need the
// run/heartbeat/force-stop machinery.
type EnvironmentWorker struct {
	client anthropic.Client
	opts   EnvironmentWorkerOptions

	// memoryClock overrides the SessionMemoryStores time source; tests use
	// it to drive the sync cadence deterministically. Nil means time.Now.
	memoryClock func() time.Time
}

// NewEnvironmentWorker returns an [EnvironmentWorker] bound to client. Call
// [EnvironmentWorker.Run] to start polling.
func NewEnvironmentWorker(client anthropic.Client, opts EnvironmentWorkerOptions) *EnvironmentWorker {
	if opts.Workdir == "" {
		// Snapshot the cwd at construction so a later os.Chdir cannot move the
		// per-session workdir out from under the worker. Mirrors the TS/Python
		// helpers, which capture process.cwd()/os.getcwd() the same way. If the
		// lookup fails, fall back to "." (resolved at use time).
		if wd, err := os.Getwd(); err == nil {
			opts.Workdir = wd
		} else {
			opts.Workdir = "."
		}
	}
	return &EnvironmentWorker{client: client, opts: opts}
}

// checkOptions rejects option values that can never work, before any request
// is sent: the deprecated UnrestrictedPaths, and a MemorySyncInterval below
// the floor.
func (w *EnvironmentWorker) checkOptions() error {
	if w.opts.UnrestrictedPaths {
		return agenttoolset.ErrUnrestrictedPathsUnsupported
	}
	return checkMemorySyncInterval(w.opts.MemorySyncInterval)
}

// Run polls the environment and services each claimed session until ctx is
// cancelled. A cancelled ctx (deadline or otherwise) is normal termination and
// yields a nil error; a non-retryable poll error is returned.
//
// EnvironmentID and EnvironmentKey are required to poll; Run returns an
// error immediately if either is unset, and likewise if the deprecated
// UnrestrictedPaths option is set or MemorySyncInterval is below
// [MinMemorySyncInterval].
func (w *EnvironmentWorker) Run(ctx context.Context) error {
	if err := w.checkOptions(); err != nil {
		return fmt.Errorf("EnvironmentWorker.Run: %w", err)
	}
	if w.opts.EnvironmentID == "" || w.opts.EnvironmentKey == "" {
		return errors.New("EnvironmentWorker.Run: EnvironmentID and EnvironmentKey are required to poll for work")
	}

	log := w.opts.Logger
	if log == nil {
		log = slog.Default()
	}

	poller := NewWorkPoller(ctx, w.client, WorkPollerOptions{
		EnvironmentID:  w.opts.EnvironmentID,
		EnvironmentKey: w.opts.EnvironmentKey,
		WorkerID:       w.opts.WorkerID,
		RequestOptions: w.opts.RequestOptions,
		AutoStop:       param.NewOpt(false),
		Logger:         log,
	})
	defer poller.Close()

	for poller.Next() {
		work := poller.Current()
		if work == nil {
			continue
		}
		// handleItem logs its own per-item failures; the poll loop keeps going.
		_ = w.handleItem(ctx, work, w.opts.EnvironmentKey)
	}
	if err := poller.Err(); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return nil
}

// HandleItemOptions selects the already-claimed work item that
// [EnvironmentWorker.HandleItem] services. Each empty field falls back to the
// matching ANTHROPIC_* environment variable — the same variables the
// `ant worker poll --on-work` hook exports into the child process:
//
//	WorkID          ← ANTHROPIC_WORK_ID
//	EnvironmentID   ← ANTHROPIC_ENVIRONMENT_ID
//	SessionID       ← ANTHROPIC_SESSION_ID
//	EnvironmentKey  ← ANTHROPIC_ENVIRONMENT_KEY
//
// WorkID, EnvironmentID and SessionID are required (after the env-var
// fallback). EnvironmentKey resolves in order: this field, then the
// [EnvironmentWorker]'s own EnvironmentKey option, then
// ANTHROPIC_ENVIRONMENT_KEY — and is also required.
type HandleItemOptions struct {
	// WorkID identifies the already-claimed work item; falls back to
	// ANTHROPIC_WORK_ID when empty.
	WorkID string
	// EnvironmentID is the self-hosted environment the work item belongs to;
	// falls back to ANTHROPIC_ENVIRONMENT_ID when empty.
	EnvironmentID string
	// SessionID is the managed-agents session the work item refers to; falls
	// back to ANTHROPIC_SESSION_ID when empty.
	SessionID string
	// EnvironmentKey authorizes the per-session calls; falls back to the
	// [EnvironmentWorker]'s own EnvironmentKey, then ANTHROPIC_ENVIRONMENT_KEY.
	EnvironmentKey string
	// WorkSecret is the work item's per-item secret payload from the poll
	// response; falls back to ANTHROPIC_WORK_SECRET (the variable the
	// `ant worker poll --on-work` command sets alongside the others). Unlike
	// the fields above it is optional — when present, the sessions token
	// extracted from it is preferred over EnvironmentKey as the Bearer
	// credential for this item's heartbeat / force-stop / skill-download /
	// session calls; when absent (or undecodable) those calls use
	// EnvironmentKey.
	WorkSecret string
}

// HandleItem services a single already-claimed session work item — the per-item
// flow [EnvironmentWorker.Run] runs for each claimed item: it builds the
// per-session [agenttoolset.AgentToolContext] (Workdir/MaxFileBytes from the
// worker's options), downloads the session agent's skills
// ([agenttoolset.AgentToolContext.SetupSkillsFromSession]), then runs a SessionToolRunner for the
// session WHILE heartbeating the work-item lease in parallel; on exit — success
// or error — it force-stops the work item. The one exception is a lost lease:
// the item then belongs to the queue or another worker and is left alone. Use
// it from a `worker poll --on-work` hook (or any caller that has already
// claimed a work item itself).
//
// Each empty field of opts is read from the matching ANTHROPIC_* environment
// variable (see [HandleItemOptions]); inside a `worker poll --on-work` child
// process every value is already exported, so HandleItem(ctx, HandleItemOptions{})
// just works. After the env-var fallback, WorkID/EnvironmentID/SessionID must
// all be non-empty or HandleItem returns an error naming the missing one.
// EnvironmentKey resolves in order — the opts field, the [EnvironmentWorker]'s
// own EnvironmentKey option, then ANTHROPIC_ENVIRONMENT_KEY — and must also
// resolve to a non-empty value. A worker built with the deprecated
// UnrestrictedPaths option, or a MemorySyncInterval below
// [MinMemorySyncInterval], returns an error before any of this.
//
// It returns the SessionToolRunner's terminal error unless that error is a
// benign session termination or idle timeout, in which case it returns nil.
// A failed session lookup fails the item (the force-stop still runs), and so
// do memory stores that cannot be mounted: the work item carried no sessions
// token for a session that has stores attached, or a store failed to
// download (a [SessionMemoryError]). A negative MemorySyncInterval is the
// exception — the item runs without its memory.
func (w *EnvironmentWorker) HandleItem(ctx context.Context, opts HandleItemOptions) error {
	if err := w.checkOptions(); err != nil {
		return fmt.Errorf("EnvironmentWorker.HandleItem: %w", err)
	}
	workID := cmp.Or(opts.WorkID, os.Getenv("ANTHROPIC_WORK_ID"))
	environmentID := cmp.Or(opts.EnvironmentID, os.Getenv("ANTHROPIC_ENVIRONMENT_ID"))
	sessionID := cmp.Or(opts.SessionID, os.Getenv("ANTHROPIC_SESSION_ID"))
	environmentKey := cmp.Or(opts.EnvironmentKey, w.opts.EnvironmentKey, os.Getenv("ANTHROPIC_ENVIRONMENT_KEY"))
	// The per-item secret is optional: the field, then ANTHROPIC_WORK_SECRET,
	// then empty (use the environment key).
	workSecret := cmp.Or(opts.WorkSecret, os.Getenv("ANTHROPIC_WORK_SECRET"))

	for _, req := range []struct{ name, val, env string }{
		{"work_id", workID, "ANTHROPIC_WORK_ID"},
		{"environment_id", environmentID, "ANTHROPIC_ENVIRONMENT_ID"},
		{"session_id", sessionID, "ANTHROPIC_SESSION_ID"},
		{"environment_key", environmentKey, "ANTHROPIC_ENVIRONMENT_KEY"},
	} {
		if req.val == "" {
			return fmt.Errorf("EnvironmentWorker.HandleItem: %s is required — pass it in HandleItemOptions or set %s", req.name, req.env)
		}
	}

	// The per-item code only reads work.ID / work.EnvironmentID / work.Secret /
	// work.Data.Type / work.Data.ID, so a minimally populated item is enough.
	work := &anthropic.BetaSelfHostedWork{
		ID:            workID,
		EnvironmentID: environmentID,
		Secret:        workSecret,
		Data: anthropic.BetaSessionWorkData{
			ID:   sessionID,
			Type: "session",
		},
	}
	return w.handleItem(ctx, work, environmentKey)
}

// handleItem is the per-item flow shared by [EnvironmentWorker.Run]'s poll loop
// and [EnvironmentWorker.HandleItem]: build the per-session [agenttoolset.AgentToolContext],
// download the session agent's skills, run a SessionToolRunner WHILE
// heartbeating the work-item lease in parallel, and force-stop the work item on
// exit unless the lease was lost. A non-session work item is only force-stopped.
//
// The per-session calls — heartbeat, force-stop, skill download, and the
// SessionToolRunner's stream/list/send — are authorized by the sessions token
// extracted from work.Secret (the item's per-item secret payload) when one is
// present, and by environmentKey otherwise; a payload that yields no token
// logs a warning and falls back to environmentKey unchanged.
func (w *EnvironmentWorker) handleItem(ctx context.Context, work *anthropic.BetaSelfHostedWork, environmentKey string) error {
	log := w.opts.Logger
	if log == nil {
		log = slog.Default()
	}

	sessionID := work.Data.ID
	log = log.With(slog.String("work_id", work.ID), slog.String("session_id", sessionID))

	// The per-item credential: the sessions token carried inside the work
	// item's secret payload when the server issued one, otherwise the
	// environment key. Never log the payload or the token.
	sessionsToken := sessionsTokenFromSecret(work.Secret)
	if work.Secret != "" && sessionsToken == "" {
		log.Warn("work item carried a secret payload but no sessions token could be extracted; falling back to the environment key")
	}
	itemCredential := cmp.Or(sessionsToken, environmentKey)

	// The per-item credential authorizes the per-session calls: the lease
	// heartbeat and the force-stop here, the skill download below, plus the
	// SessionToolRunner's stream/list/send. The x-stainless-helper header
	// attributes the heartbeat/force-stop/skill traffic to this helper.
	// helperReqOpts also clears the parent client's default X-Api-Key so it
	// doesn't ride alongside the bearer credential. Caller-supplied
	// RequestOptions are applied first so a proxy/custom header reaches every
	// per-session call while the worker's own X-Api-Key delete, per-item
	// bearer auth and helper header (appended last) still win.
	helperOpts, err := helperReqOpts(itemCredential, stainlessheader.EnvironmentsWorker)
	if err != nil {
		// Run and HandleItem validate environmentKey at their entry points
		// (itemCredential can only be empty if it was), so an empty
		// credential here means a future code path bypassed that
		// validation; surface it rather than fire requests with the parent
		// client's credentials.
		return err
	}
	hbStopOpts := make([]option.RequestOption, 0, len(w.opts.RequestOptions)+len(helperOpts))
	hbStopOpts = append(hbStopOpts, w.opts.RequestOptions...)
	hbStopOpts = append(hbStopOpts, helperOpts...)

	// The heartbeat records each server-reported lease TTL here; the runner
	// picks it up from sessCtx as its tool-result send retry window.
	var leaseTTL sendwindow.Window

	// Per-session context: cancelled when the outer ctx is cancelled (it is a
	// child), when the session runner finishes, or when the lease heartbeat
	// says to stop. Constructed BEFORE skill setup so the heartbeat goroutine
	// below can use it.
	sessCtx, sessCancel := context.WithCancel(sendwindow.NewContext(ctx, &leaseTTL))
	defer sessCancel()

	// Start the lease heartbeat BEFORE skill setup. The poller already acked
	// this work item when it yielded — every second between the ack and the
	// first heartbeat is a window during which the control plane sees no
	// liveness signal and may reclaim the lease. The session lookup and skill
	// setup below can be slow (a per-skill download/extract can dwarf the
	// lease TTL on a slow network or a large bundle), so starting the
	// heartbeat afterwards was a race that let a second worker pick up the
	// same session.
	hbDone := make(chan struct{})
	// Written by the heartbeat goroutine before hbDone closes; read only
	// after receiving from hbDone.
	var leaseEnd leaseEndReason
	go func() {
		defer close(hbDone)
		leaseEnd = runHeartbeat(sessCtx, w.client, work, hbStopOpts, log, &leaseTTL)
		sessCancel()
	}()

	// Deferred first so it runs last on every exit path: end the heartbeat,
	// then force-stop the item on a fresh context (a cancelled ctx must not
	// skip it; a 409 means it already stopped) unless the lease was lost.
	defer func() {
		sessCancel()
		<-hbDone
		if leaseEnd.lost() {
			log.Info("lease lost; released without stopping it")
			return
		}

		stopCtx, cancel := context.WithTimeout(context.Background(), stopTimeout)
		defer cancel()
		if err := stopWork(stopCtx, w.client, work.ID,
			anthropic.BetaEnvironmentWorkStopParams{
				EnvironmentID: work.EnvironmentID,
				BetaSelfHostedWorkStopRequest: anthropic.BetaSelfHostedWorkStopRequestParam{
					Force: param.NewOpt(true),
				},
			},
			hbStopOpts...,
		); err != nil && !isStatus(err, 409) {
			log.Warn("force-stop on exit failed", slog.Any("error", err))
		}
	}()

	// The queue also sends "healthcheck" items, which the generated type does not model.
	if work.Data.Type != "session" {
		log.Debug("skipping non-session work item", slog.String("type", string(work.Data.Type)))
		return nil
	}

	env := &agenttoolset.AgentToolContext{
		Workdir:      w.opts.Workdir,
		MaxFileBytes: w.opts.MaxFileBytes,
	}
	// The session lookup, skill download and memory calls are
	// environment-scoped, so they need the per-item credential like the
	// heartbeat/stop and the runner do — without it they fall back to the
	// client's default credentials and fail.
	// Use sessCtx so a heartbeat-driven lease loss (the heartbeat goroutine
	// cancels sessCtx on a permanent failure / stopping state / 412 reclaim)
	// also aborts the downloads instead of letting them run to completion on
	// a session we no longer own.
	//
	// One session fetch, shared by the skills download and the memory-store
	// download — two fetches could disagree about the attached resources.
	session, err := w.client.Beta.Sessions.Get(sessCtx, sessionID, anthropic.BetaSessionGetParams{}, hbStopOpts...)
	if err != nil {
		// A session that can't be fetched fails the whole item (the deferred
		// force-stop still runs), as does a memory store that can't be
		// materialised below. Per-skill failures stay tolerant: a missing
		// skill degrades the session, a missing memory folder corrupts it.
		err = fmt.Errorf("retrieve session %s: %w", sessionID, err)
		log.Warn("session lookup failed", slog.Any("error", err))
		return err
	}
	if err := env.SetupSkillsFromSession(sessCtx, w.client, session, hbStopOpts...); err != nil {
		log.Warn("skill setup failed", slog.Any("error", err))
	}
	// Clean up the skills this work item downloaded so one session's skills
	// don't leak into the next item served by the same worker.
	defer func() {
		if err := env.Cleanup(); err != nil {
			log.Warn("skill cleanup failed", slog.Any("error", err))
		}
	}()

	// Memory stores: the memory_stores endpoints accept the per-item
	// sessions token but reject the environment key, so download and sync
	// only run when the item carried a usable secret (and the interval isn't
	// negative). The memory calls ride the same token-scoped request options
	// as the heartbeat and skill download.
	var stores *SessionMemoryStores
	var cleanEnd bool
	if sessionsToken != "" && w.opts.MemorySyncInterval >= 0 {
		var err error
		stores, err = NewSessionMemoryStores(w.client, SessionMemoryStoresOptions{
			Workdir:        w.opts.Workdir,
			SyncInterval:   w.opts.MemorySyncInterval,
			SyncDeletions:  w.opts.MemorySyncDeletions,
			RequestOptions: hbStopOpts,
			Logger:         log,
		})
		if err != nil {
			// Unreachable after checkOptions; kept so a bad interval can never
			// run a session without its memory.
			return err
		}
		if w.memoryClock != nil {
			stores.setClock(w.memoryClock)
		}
		// Registered before the download, so a store that landed before a
		// later one failed still has its directory removed. Runs after the
		// agenttoolset.CloseAll defer — bash killed, nothing still writing.
		defer func() { stores.Cleanup(cleanEnd) }()
		if err := stores.Download(sessCtx, session); err != nil {
			// A store that cannot be materialised fails the whole item (the
			// deferred force-stop still runs): a session whose system prompt
			// names a memory folder that isn't there would run with amnesia
			// and sync nothing back, which is worse than not running at all.
			log.Error("memory store download failed", slog.Any("error", err))
			return err
		}
		// A store mounted outside the workdir must stay reachable by the file
		// tools; read-only stores still refuse writes via ReadOnlyRoots.
		env.AllowedRoots = stores.Roots()
		env.ReadOnlyRoots = stores.ReadOnlyRoots()
	} else {
		// The gate above cannot tell a session that simply has no memory from
		// one whose memory we cannot mount — only the fetched session can. A
		// negative interval is a deliberate opt-out and stays quiet; a missing
		// sessions token on a session that *does* have stores fails the whole
		// item (the deferred force-stop still runs) — a hosted sandbox refuses
		// to start in this state, and running here without memories would
		// silently diverge from it.
		if w.opts.MemorySyncInterval >= 0 && hasMemoryStore(session) {
			err := fmt.Errorf("session %s: %w", sessionID, ErrSessionMemoryNoToken)
			log.Error("memory stores cannot be mounted", slog.Any("error", err))
			return err
		}
		log.Debug("memory stores disabled for this item")
	}

	var (
		tools      []anthropic.BetaTool
		closeTools bool
	)
	switch {
	case w.opts.ToolsFunc != nil:
		tools = w.opts.ToolsFunc(env)
		closeTools = true
	case w.opts.Tools != nil:
		tools = w.opts.Tools
	default:
		tools = agenttoolset.BetaAgentToolset20260401(env)
		closeTools = true
	}
	if closeTools {
		defer agenttoolset.CloseAll(tools)
	}

	// Authorize the runner's stream/list/send calls with the per-item
	// credential (via bearerReqOpts — the runner stamps its own x-stainless-helper
	// header "session-tool-runner" internally *after* these options, so we
	// just need the auth bits here; any helper tag we passed in would be
	// overwritten). bearerReqOpts also clears the parent client's default
	// X-Api-Key so it doesn't ride alongside the bearer credential.
	// Caller-supplied RequestOptions go first so a proxy/custom header
	// reaches stream/list/send too; the X-Api-Key delete + per-item bearer
	// auth (appended last) still win.
	runnerBearerOpts, err := bearerReqOpts(itemCredential)
	if err != nil {
		// Same validation invariant as the hbStopOpts construction above —
		// surface rather than send with the parent client's credentials.
		return err
	}
	runnerReqOpts := make([]option.RequestOption, 0, len(w.opts.RequestOptions)+len(runnerBearerOpts))
	runnerReqOpts = append(runnerReqOpts, w.opts.RequestOptions...)
	runnerReqOpts = append(runnerReqOpts, runnerBearerOpts...)
	runner := w.client.Beta.Sessions.Events.NewToolRunner(sessCtx, sessionID, anthropic.SessionToolRunnerOptions{
		Tools:          tools,
		MaxIdle:        w.opts.MaxIdle,
		Logger:         log,
		RequestOptions: runnerReqOpts,
	})
	for runner.Next() {
		call := runner.Current()
		log.Info("dispatched tool",
			slog.String("tool", call.Name),
			slog.Bool("is_error", call.IsError),
			slog.Bool("posted", call.Posted))
		if stores != nil {
			stores.SyncIfDue(sessCtx)
		}
	}
	var runErr error
	if err := runner.Err(); err != nil &&
		!errors.Is(err, anthropic.ErrSessionTerminated) &&
		!errors.Is(err, anthropic.ErrIdleTimeout) {
		log.Warn("session tool runner exited with error", slog.Any("error", err))
		runErr = err
	}
	cleanEnd = runErr == nil && sessCtx.Err() == nil
	// Close is documented to always return nil (it satisfies io.Closer so
	// callers can defer it uniformly), but don't bank on that: a future
	// non-nil error is surfaced rather than dropped.
	if err := runner.Close(); err != nil {
		log.Warn("session tool runner close failed", slog.Any("error", err))
		if runErr == nil {
			runErr = err
		}
	}

	return runErr
}

// leaseEndReason is why heartbeating of a work item ended. The zero value,
// leaseHeld, means it has not.
type leaseEndReason int

const (
	leaseHeld leaseEndReason = iota
	// ctx ended first: the run finished or was cancelled with the lease held.
	leaseRunnerDone
	// State stopping/stopped, or lease_extended false.
	leaseControlPlaneStop
	// A heartbeat met 412: the item was re-queued or another worker holds it.
	leaseLost
	// A heartbeat met another 4xx that retrying will not fix.
	leaseHeartbeatRejected
	// No heartbeat succeeded for a full lease TTL.
	leaseAssumedLost
)

// lost reports whether the item now belongs to the queue or another worker,
// so this one must not stop it.
func (r leaseEndReason) lost() bool {
	return r == leaseLost || r == leaseAssumedLost
}

// runHeartbeat keeps the work-item lease alive while a session is being served
// and returns why heartbeating ended: ctx ended first (the run finished or was
// cancelled), the control plane reported the work stopping/stopped or no longer
// extended the lease, a heartbeat was rejected (a 412 means the lease already
// belongs to someone else), or transient failures ran long enough that the
// lease must be assumed lost. The caller cancels the session once it returns,
// so two runners don't end up serving the same work.
//
// Each Heartbeat request is sent with a per-request timeout derived from the
// last server-reported TTL (or heartbeatDefault for the first beat); a hung
// request can't outlive the lease window. A run of transient failures is
// bounded by a staleness ceiling: if more than the last known TTL has elapsed
// since the most recent successful beat, the lease is presumed expired
// server-side and the session is cancelled rather than executing tools
// against a session another worker may also have claimed. leaseTTL is updated
// with the server-reported TTL after every successful beat.
func runHeartbeat(ctx context.Context, client anthropic.Client, work *anthropic.BetaSelfHostedWork, reqOpts []option.RequestOption, log *slog.Logger, leaseTTL *sendwindow.Window) leaseEndReason {
	interval := heartbeatDefault
	// ttl tracks the last server-reported TTL. It bounds the staleness
	// ceiling: a run of transient errors lasting longer than this means the
	// server has almost certainly let the lease expire, so we must stop
	// heartbeating and cancel the session before posting more tool results
	// against a reclaimed lease. Initialized to heartbeatDefault so a
	// permanently-failing first beat is bounded too.
	ttl := heartbeatDefault
	last := noHeartbeatSentinel
	// lastSuccess seeds the staleness clock from goroutine start so the
	// first-beat-never-succeeds case is bounded by ttl rather than retrying
	// forever.
	lastSuccess := time.Now()

	// beat sends one heartbeat and returns leaseHeld to keep beating, or why
	// the lease ended.
	beat := func() leaseEndReason {
		// Per-request timeout: cap at the current wait-between-beats interval
		// (which tracks ttl/2). A single Heartbeat call must never outlive
		// the lease window — without this, the call inherits the SDK default
		// (~10 minutes) and one hung request can let the lease expire while
		// we sit on the connection. The full slice expression caps capacity
		// so append cannot mutate the caller's reqOpts backing array.
		beatOpts := append(reqOpts[:len(reqOpts):len(reqOpts)], option.WithRequestTimeout(interval))
		resp, err := client.Beta.Environments.Work.Heartbeat(
			ctx,
			work.ID,
			anthropic.BetaEnvironmentWorkHeartbeatParams{
				EnvironmentID:         work.EnvironmentID,
				ExpectedLastHeartbeat: param.NewOpt(last),
			},
			beatOpts...,
		)
		if err != nil {
			// A 412 means the server re-queued the item or another worker holds
			// the lease, so this worker must neither keep serving it nor stop
			// it; checked before ctx so a 412 racing the run's finish counts.
			if isStatus(err, 412) {
				server := serverLeaseState(err)
				log.Error("lease lost: heartbeat precondition failed",
					slog.Any("server_state", server["state"]),
					slog.Any("server_ttl_seconds", server["ttl_seconds"]),
					slog.Any("server_last_heartbeat", server["last_heartbeat"]))
				return leaseLost
			}
			if ctx.Err() != nil {
				return leaseRunnerDone
			}
			if isFatal4xx(err) {
				log.Error("permanent heartbeat failure", slog.Any("error", err))
				return leaseHeartbeatRejected
			}
			// Bound the transient-retry window. The control plane lets the
			// lease lapse silently after ttl elapsed without a successful
			// heartbeat — without this check we'd keep executing tools
			// against a reclaimed session.
			if stale := time.Since(lastSuccess); stale > ttl {
				log.Error("heartbeat staleness ceiling exceeded; lease presumed expired",
					slog.Duration("since_last_success", stale),
					slog.Duration("ttl", ttl),
					slog.Any("error", err))
				return leaseAssumedLost
			}
			log.Warn("transient heartbeat failure", slog.Any("error", err))
			return leaseHeld
		}
		last = resp.LastHeartbeat
		lastSuccess = time.Now()
		if resp.TTLSeconds > 0 {
			ttl = max(time.Duration(resp.TTLSeconds)*time.Second, heartbeatFloor)
			interval = clampDur(ttl/2, heartbeatFloor, heartbeatDefault)
			leaseTTL.Set(ttl)
		}
		switch resp.State {
		case anthropic.BetaSelfHostedWorkHeartbeatResponseStateStopping,
			anthropic.BetaSelfHostedWorkHeartbeatResponseStateStopped:
			log.Info("heartbeat reports shutdown", slog.String("state", string(resp.State)))
			return leaseControlPlaneStop
		}
		if !resp.LeaseExtended {
			log.Warn("lease not extended; shutting down")
			return leaseControlPlaneStop
		}
		return leaseHeld
	}

	if end := beat(); end != leaseHeld {
		return end
	}
	for {
		t := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			t.Stop()
			return leaseRunnerDone
		case <-t.C:
		}
		if end := beat(); end != leaseHeld {
			return end
		}
	}
}

func clampDur(v, lo, hi time.Duration) time.Duration {
	return max(lo, min(hi, v))
}

// serverLeaseState is the server's view of the lease carried in a rejected
// heartbeat's response body, or nil when the body carries none.
func serverLeaseState(err error) map[string]any {
	var apiErr *anthropic.Error
	if !errors.As(err, &apiErr) {
		return nil
	}
	var body struct {
		Error struct {
			Details struct {
				CurrentState map[string]any `json:"current_state"`
			} `json:"details"`
		} `json:"error"`
	}
	_ = json.Unmarshal([]byte(apiErr.RawJSON()), &body)
	return body.Error.Details.CurrentState
}

// hasMemoryStore reports whether the session has at least one memory store
// attached.
func hasMemoryStore(session *anthropic.BetaManagedAgentsSession) bool {
	for _, res := range session.Resources {
		if res.Type == "memory_store" {
			return true
		}
	}
	return false
}

// sessionsTokenFromSecret extracts the per-item sessions token from a work
// item's secret payload.
//
// The secret the poll response populates is not itself a credential: it is a
// URL-safe base64 JSON payload bundling the per-item material — the
// sessions_token (the bearer for this item's work lifecycle and session-level
// calls) plus ingress / source tokens this worker does not consume. Returns
// the sessions token, or "" (meaning: fall back to the environment key) when
// the payload is missing, doesn't decode, or carries no token. Never log the
// payload or anything extracted from it.
func sessionsTokenFromSecret(secret string) string {
	if secret == "" {
		return ""
	}
	// The payload may arrive with or without base64 padding; RawURLEncoding
	// after stripping any padding accepts both forms.
	raw, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(secret, "="))
	if err != nil {
		return ""
	}
	var parsed struct {
		SessionsToken string `json:"sessions_token"`
	}
	// A payload that is not a JSON object (or has no usable token) resolves
	// to "" — the caller falls back to the environment key.
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return ""
	}
	return parsed.SessionsToken
}
