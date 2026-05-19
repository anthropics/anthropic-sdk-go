package environments

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
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

	// EnvironmentKey is the environment's key. It authorizes both the
	// work-poll calls (Poll/Ack/Stop) and the per-session calls (the lease
	// heartbeat, the force-stop, and the SessionToolRunner's event
	// stream/list/send). Required by [EnvironmentWorker.Run]; for
	// [EnvironmentWorker.HandleItem] it is one of the resolution sources for
	// the per-item environment key (see [HandleItemOptions]).
	EnvironmentKey string

	// WorkerID is a stable identifier reported back to the server for
	// observability. Defaults to "<os.Hostname()>-<random hex>" — forwarded
	// to [WorkPoller].
	WorkerID string

	// Workdir is the base directory for the per-session
	// [agenttoolset.AgentToolContext]. When empty it defaults to the process
	// working directory captured when [NewEnvironmentWorker] is called.
	Workdir string

	// UnrestrictedPaths is forwarded to the per-session
	// [agenttoolset.AgentToolContext].
	UnrestrictedPaths bool

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
// ([agenttoolset.AgentToolContext.SetupSkills]), then runs a SessionToolRunner for the
// session WHILE heartbeating the work-item lease in parallel; on exit it
// force-stops the work item and loops to the next one. The lease heartbeat
// reports state "stopping"/"stopped" or a lost lease back into the run by
// cancelling the session runner.
//
// [EnvironmentWorker.Run] drives the poll loop. [EnvironmentWorker.HandleItem]
// exposes the per-item flow for callers that already hold a claimed work item
// (for example a `worker poll --on-work` hook) and only need the
// run/heartbeat/force-stop machinery.
type EnvironmentWorker struct {
	client anthropic.Client
	opts   EnvironmentWorkerOptions
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

// Run polls the environment and services each claimed session until ctx is
// cancelled. A cancelled ctx (deadline or otherwise) is normal termination and
// yields a nil error; a non-retryable poll error is returned.
//
// EnvironmentID and EnvironmentKey are required to poll; Run returns an
// error immediately if either is unset.
func (w *EnvironmentWorker) Run(ctx context.Context) error {
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
}

// HandleItem services a single already-claimed session work item — the per-item
// flow [EnvironmentWorker.Run] runs for each claimed item: it builds the
// per-session [agenttoolset.AgentToolContext] (workdir/UnrestrictedPaths from the
// worker's options), downloads the session agent's skills
// ([agenttoolset.AgentToolContext.SetupSkills]), then runs a SessionToolRunner for the
// session WHILE heartbeating the work-item lease in parallel; on exit — success
// or error — it force-stops the work item. Use it from a `worker poll
// --on-work` hook (or any caller that has already claimed a work item itself).
//
// Each empty field of opts is read from the matching ANTHROPIC_* environment
// variable (see [HandleItemOptions]); inside a `worker poll --on-work` child
// process every value is already exported, so HandleItem(ctx, HandleItemOptions{})
// just works. After the env-var fallback, WorkID/EnvironmentID/SessionID must
// all be non-empty or HandleItem returns an error naming the missing one.
// EnvironmentKey resolves in order — the opts field, the [EnvironmentWorker]'s
// own EnvironmentKey option, then ANTHROPIC_ENVIRONMENT_KEY — and must also
// resolve to a non-empty value.
//
// It returns the SessionToolRunner's terminal error unless that error is a
// benign session termination or idle timeout, in which case it returns nil.
func (w *EnvironmentWorker) HandleItem(ctx context.Context, opts HandleItemOptions) error {
	workID := cmp.Or(opts.WorkID, os.Getenv("ANTHROPIC_WORK_ID"))
	environmentID := cmp.Or(opts.EnvironmentID, os.Getenv("ANTHROPIC_ENVIRONMENT_ID"))
	sessionID := cmp.Or(opts.SessionID, os.Getenv("ANTHROPIC_SESSION_ID"))
	environmentKey := cmp.Or(opts.EnvironmentKey, w.opts.EnvironmentKey, os.Getenv("ANTHROPIC_ENVIRONMENT_KEY"))

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

	// The per-item code only reads work.ID / work.EnvironmentID /
	// work.Data.Type / work.Data.ID, so minimal struct literals are enough.
	work := &anthropic.BetaSelfHostedWork{
		ID:            workID,
		EnvironmentID: environmentID,
		Data: anthropic.BetaSessionWorkData{
			ID: sessionID,
		},
	}
	return w.handleItem(ctx, work, environmentKey)
}

// handleItem is the per-item flow shared by [EnvironmentWorker.Run]'s poll loop
// and [EnvironmentWorker.HandleItem]: build the per-session [agenttoolset.AgentToolContext],
// download the session agent's skills, run a SessionToolRunner WHILE
// heartbeating the work-item lease in parallel, and force-stop the work item on
// exit. environmentKey authorizes the per-session heartbeat/stop calls and the
// SessionToolRunner's stream/list/send; work must be a `session` work item.
func (w *EnvironmentWorker) handleItem(ctx context.Context, work *anthropic.BetaSelfHostedWork, environmentKey string) error {
	log := w.opts.Logger
	if log == nil {
		log = slog.Default()
	}

	sessionID := work.Data.ID
	log = log.With(slog.String("work_id", work.ID), slog.String("session_id", sessionID))

	// The environment key authorizes the per-session calls: the lease
	// heartbeat and the force-stop here, the skill download below, plus the
	// SessionToolRunner's stream/list/send. The x-stainless-helper header
	// attributes the heartbeat/force-stop/skill traffic to this helper.
	// helperReqOpts also clears the parent client's default X-Api-Key so it
	// doesn't ride alongside the bearer credential. Caller-supplied
	// RequestOptions are applied first so a proxy/custom header reaches every
	// per-session call while the worker's own X-Api-Key delete,
	// environment-key auth and helper header (appended last) still win.
	helperOpts, err := helperReqOpts(environmentKey, helperHeaderEnvironmentWorker)
	if err != nil {
		// Run and HandleItem validate environmentKey at their entry points
		// (see worker.go:130 / 209-218), so an empty key here means a future
		// code path bypassed that validation; surface it rather than fire
		// requests with the parent client's credentials.
		return err
	}
	hbStopOpts := make([]option.RequestOption, 0, len(w.opts.RequestOptions)+len(helperOpts))
	hbStopOpts = append(hbStopOpts, w.opts.RequestOptions...)
	hbStopOpts = append(hbStopOpts, helperOpts...)

	// Per-session context: cancelled when the outer ctx is cancelled (it is a
	// child), when the session runner finishes, or when the lease heartbeat
	// says to stop. Constructed BEFORE skill setup so the heartbeat goroutine
	// below can use it.
	sessCtx, sessCancel := context.WithCancel(ctx)
	defer sessCancel()

	// Start the lease heartbeat BEFORE skill setup. The poller already acked
	// this work item when it yielded — every second between the ack and the
	// first heartbeat is a window during which the control plane sees no
	// liveness signal and may reclaim the lease. SetupSkills below can be
	// slow (it issues a session lookup plus a per-skill download/extract that
	// can dwarf the lease TTL on a slow network or a large bundle), so
	// starting the heartbeat afterwards was a race that let a second worker
	// pick up the same session.
	hbDone := make(chan struct{})
	go func() {
		defer close(hbDone)
		runHeartbeat(sessCtx, sessCancel, w.client, work, hbStopOpts, log)
	}()

	env := &agenttoolset.AgentToolContext{
		Workdir:           w.opts.Workdir,
		UnrestrictedPaths: w.opts.UnrestrictedPaths,
	}
	// The session lookup and skill download are environment-scoped, so they
	// need the environment key like the heartbeat/stop and the runner do —
	// without it they fall back to the client's default credentials and fail.
	// Use sessCtx so a heartbeat-driven lease loss (the heartbeat goroutine
	// cancels sessCtx on a permanent failure / stopping state / 412 reclaim)
	// also aborts the skill download instead of letting it run to completion
	// on a session we no longer own.
	if err := env.SetupSkills(sessCtx, w.client, sessionID, hbStopOpts...); err != nil {
		log.Warn("skill setup failed", slog.Any("error", err))
	}
	// Clean up the skills this work item downloaded so one session's skills
	// don't leak into the next item served by the same worker.
	defer func() {
		if err := env.Cleanup(); err != nil {
			log.Warn("skill cleanup failed", slog.Any("error", err))
		}
	}()

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

	// Authorize the runner's stream/list/send calls with the environment key
	// (via bearerReqOpts — the runner stamps its own x-stainless-helper
	// header "session-tool-runner" internally *after* these options, so we
	// just need the auth bits here; any helper tag we passed in would be
	// overwritten). bearerReqOpts also clears the parent client's default
	// X-Api-Key so it doesn't ride alongside the bearer credential.
	// Caller-supplied RequestOptions go first so a proxy/custom header
	// reaches stream/list/send too; the X-Api-Key delete + environment-key
	// auth (appended last) still win.
	runnerBearerOpts, err := bearerReqOpts(environmentKey)
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
	}
	var runErr error
	if err := runner.Err(); err != nil &&
		!errors.Is(err, anthropic.ErrSessionTerminated) &&
		!errors.Is(err, anthropic.ErrIdleTimeout) {
		log.Warn("session tool runner exited with error", slog.Any("error", err))
		runErr = err
	}
	_ = runner.Close()

	sessCancel()
	<-hbDone

	// Force-stop the work item on a fresh context so a cancelled outer ctx
	// doesn't skip the stop.
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

	return runErr
}

// runHeartbeat keeps the work-item lease alive while a session is being served.
// It calls cancel when the control plane reports the work is stopping/stopped,
// when the lease is no longer extended, or on a permanent heartbeat failure.
//
// Each Heartbeat request is sent with a per-request timeout derived from the
// last server-reported TTL (or heartbeatDefault for the first beat); a hung
// request can't outlive the lease window. A run of transient failures is
// bounded by a staleness ceiling: if more than the last known TTL has elapsed
// since the most recent successful beat, the lease is presumed expired
// server-side and the session is cancelled rather than executing tools
// against a session another worker may also have claimed.
func runHeartbeat(ctx context.Context, cancel context.CancelFunc, client anthropic.Client, work *anthropic.BetaSelfHostedWork, reqOpts []option.RequestOption, log *slog.Logger) {
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

	// beat returns false when the worker should stop heartbeating (and the
	// caller should cancel the session).
	beat := func() bool {
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
			if ctx.Err() != nil {
				return false
			}
			// 412 means our expected_last_heartbeat no longer matches the
			// server's: another worker reclaimed the lease (or won the
			// first-heartbeat race against our "NO_HEARTBEAT" claim). Call it
			// out distinctly from an auth/spec 4xx — the operational response
			// is "someone else owns this work now", not "the client is broken".
			if isStatus(err, 412) {
				log.Warn("heartbeat precondition failed; lease reclaimed by another worker",
					slog.Any("error", err))
				return false
			}
			if isFatal4xx(err) {
				log.Error("permanent heartbeat failure", slog.Any("error", err))
				return false
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
				return false
			}
			log.Warn("transient heartbeat failure", slog.Any("error", err))
			return true
		}
		last = resp.LastHeartbeat
		lastSuccess = time.Now()
		if resp.TTLSeconds > 0 {
			ttl = max(time.Duration(resp.TTLSeconds)*time.Second, heartbeatFloor)
			interval = clampDur(ttl/2, heartbeatFloor, heartbeatDefault)
		}
		switch resp.State {
		case anthropic.BetaSelfHostedWorkHeartbeatResponseStateStopping,
			anthropic.BetaSelfHostedWorkHeartbeatResponseStateStopped:
			log.Info("heartbeat reports shutdown", slog.String("state", string(resp.State)))
			return false
		}
		if !resp.LeaseExtended {
			log.Warn("lease not extended; shutting down")
			return false
		}
		return true
	}

	if !beat() {
		cancel()
		return
	}
	for {
		t := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
		}
		if !beat() {
			cancel()
			return
		}
	}
}

func clampDur(v, lo, hi time.Duration) time.Duration {
	return max(lo, min(hi, v))
}
