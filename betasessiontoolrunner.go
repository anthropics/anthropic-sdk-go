package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"log/slog"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"time"

	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"golang.org/x/sync/errgroup"
)

// sessionRunnerHelperHeader and its value tag the requests the
// SessionToolRunner issues (stream / list / send) so the control plane can
// attribute traffic to the helper rather than to bare client calls.
const (
	sessionRunnerHelperHeader = "x-stainless-helper"
	sessionRunnerHelperValue  = "session-tool-runner"
)

// DefaultMaxIdle is used for [SessionToolRunnerOptions.MaxIdle] when it is nil:
// once the session goes idle with stop_reason "end_turn", the runner keeps
// running for this long before stopping (any new event resets the countdown).
const DefaultMaxIdle = 60 * time.Second

const (
	sessionRunnerStreamBackoffStart = 500 * time.Millisecond
	sessionRunnerStreamBackoffCap   = 10 * time.Second
	// sessionRunnerStreamHealthyAfter is how long a stream must stay open
	// before the reconnect backoff resets.
	sessionRunnerStreamHealthyAfter = 5 * time.Second
	sessionRunnerSendRetries        = 3
	sessionRunnerResultsBuffer      = 32 // enough to absorb a brief consumer pause without back-pressuring dispatch
	sessionRunnerToolUseQueueBuffer = 32
)

// Defaults for the tunable per-runner timeouts. Each is overridable via the
// matching [SessionToolRunnerOptions] field (ToolTimeout / SendTimeout /
// DrainTimeout); a zero option falls back to the default here, so the
// out-of-the-box behavior is identical when nothing is set.
const (
	defaultSessionRunnerToolTimeout  = 120 * time.Second
	defaultSessionRunnerSendTimeout  = 15 * time.Second
	defaultSessionRunnerDrainTimeout = 30 * time.Second
)

// Sentinel errors returned via (*SessionToolRunner).Err so consumers can
// distinguish normal end-of-session from idle bailout.
var (
	// ErrSessionTerminated means the server emitted session.status_terminated
	// or session.deleted. Iteration ends normally; the consumer's outer loop
	// can move on.
	ErrSessionTerminated = errors.New("session terminated")

	// ErrIdleTimeout means the session went idle with stop_reason "end_turn"
	// and stayed quiet for MaxIdle. The consumer chose this timeout, so treat
	// it as expected.
	ErrIdleTimeout = errors.New("session idle after end_turn")
)

// SessionToolRunnerOptions configures a [SessionToolRunner]. The
// managed-agents session id is not part of this struct — it is a positional
// argument to NewToolRunner, matching list/send/stream on the events resource.
type SessionToolRunnerOptions struct {
	// Tools is the registry of locally-executable tools, in the same
	// [BetaTool] shape (*BetaMessageService).NewToolRunner accepts. The runner
	// looks up each agent.tool_use and agent.custom_tool_use event's Name
	// against this slice and routes to the matching tool. Use
	// agenttoolset.BetaAgentToolset20260401(env) (from
	// github.com/anthropics/anthropic-sdk-go/tools/agenttoolset) for the
	// standard agent_toolset_20260401 set; filter or extend the slice to
	// customise. Tool lifetime — including Close on tools that implement
	// io.Closer — is the caller's responsibility: the runner never closes
	// tools, so the same slice can be reused across multiple sessions.
	Tools []BetaTool

	// MaxIdle is how long the runner keeps running after the session goes idle
	// with stop_reason "end_turn" before it stops; any new event resets the
	// countdown and it re-arms on the next "end_turn" idle. nil uses
	// [DefaultMaxIdle] (60s). A non-nil value <= 0 disables it — the runner
	// then only stops on session termination, the consumer breaking out, or
	// ctx cancellation.
	MaxIdle *time.Duration

	// Logger receives non-fatal warnings (transient stream errors, send
	// retries). Defaults to slog.Default().
	Logger *slog.Logger

	// RequestOptions are applied to every request the runner issues — the
	// event stream, the reconcile list, and the tool-result send. The runner
	// additionally tags each request with an x-stainless-helper header of
	// its own (appended last, so it always wins on a collision with a
	// caller-supplied tag).
	//
	// Self-hosted-environment callers must pass the environment key here.
	// option.WithAuthToken alone only ADDS an Authorization header — the
	// parent client's WithAPIKey middleware still emits X-Api-Key on every
	// request, so both creds would land on the wire and the server rejects
	// the dual auth on the events stream. Pair the bearer with an explicit
	// X-Api-Key delete:
	//
	//	RequestOptions: []option.RequestOption{
	//		option.WithHeaderDel("X-Api-Key"),
	//		option.WithAuthToken(environmentKey),
	//	}
	RequestOptions []option.RequestOption

	// ToolTimeout bounds a single tool's Execute call; a tool that runs longer
	// is cancelled and reported as an error result. Zero uses the default of
	// 120s.
	ToolTimeout time.Duration

	// SendTimeout bounds a single attempt to post a tool-result event back to
	// the session (the runner still retries transient failures up to
	// sessionRunnerSendRetries times). Zero uses the default of 15s.
	SendTimeout time.Duration

	// DrainTimeout bounds how long Close waits for in-flight tool executions
	// to finish before it logs a warning and returns. Zero uses the default of
	// 30s.
	DrainTimeout time.Duration
}

// DispatchedToolCall describes a single tool-use event that the runner has
// executed locally and (attempted to) post a result for. The consumer observes
// one of these per iteration of (*SessionToolRunner).Next.
//
// The runner dispatches both kinds of tool-call event the session API emits:
// agent.tool_use (the builtin agent_toolset_20260401 tools, answered with
// user.tool_result) and agent.custom_tool_use (user-defined function tools,
// answered with user.custom_tool_result). The Custom field reports which kind
// this dispatch came from and selects which embedded-event pair is populated —
// ToolUse/Result for a builtin call, CustomToolUse/CustomResult for a custom
// one. The flat top-level convenience fields (ToolUseID / Name / IsError /
// Posted) are populated for both kinds; read whichever is more convenient,
// they describe the same dispatch. The raw tool input lives on the triggering
// event (ToolUse.Input / CustomToolUse.Input) and the tool's output blocks on
// the posted result (Result.Content / CustomResult.Content).
//
// IsError is true if the tool reported failure or exceeded its per-call
// timeout. Posted is orthogonal to IsError: it reports whether a result event
// was successfully sent back to the session. Posted=false means the agent will
// not see a result from this runner, regardless of IsError — either because
// the send failed or because the tool name is not one this runner owns and the
// runner deliberately posted nothing and left the id pending for its owner.
type DispatchedToolCall struct {
	// Custom reports whether this dispatch was triggered by an
	// agent.custom_tool_use event (a user-defined function tool) rather than an
	// agent.tool_use event (a builtin agent_toolset_20260401 tool). It selects
	// which of the ToolUse/CustomToolUse and Result/CustomResult pairs below is
	// populated; the other pair is left zero.
	Custom bool

	// ToolUse is the full agent.tool_use event that triggered this dispatch,
	// exactly as it arrived on the stream (or via reconcile). Populated only
	// when Custom is false.
	ToolUse BetaManagedAgentsAgentToolUseEvent

	// CustomToolUse is the full agent.custom_tool_use event that triggered this
	// dispatch, exactly as it arrived on the stream (or via reconcile).
	// Populated only when Custom is true.
	CustomToolUse BetaManagedAgentsAgentCustomToolUseEvent

	// Result is the user.tool_result event the runner built and posted back to
	// the session for a builtin agent.tool_use. Populated only when Custom is
	// false. Read Result.Content for the tool's output blocks; the flat
	// IsError mirrors Result.IsError and Posted reports whether the send
	// succeeded.
	Result BetaManagedAgentsUserToolResultEventParams

	// CustomResult is the user.custom_tool_result event the runner built and
	// posted back to the session for an agent.custom_tool_use. Populated only
	// when Custom is true. Read CustomResult.Content for the tool's output
	// blocks; the flat IsError mirrors CustomResult.IsError and Posted
	// reports whether the send succeeded.
	CustomResult BetaManagedAgentsUserCustomToolResultEventParams

	// ToolUseID is the id of the tool-use event that triggered this dispatch
	// (flat convenience copy of ToolUse.ID / CustomToolUse.ID). It is the id
	// the posted result event references — tool_use_id for a builtin call,
	// custom_tool_use_id for a custom one. Use it to correlate logs with the
	// session history.
	ToolUseID string

	// Name is the tool name from the tool-use event (flat convenience copy of
	// ToolUse.Name / CustomToolUse.Name).
	Name string

	// IsError captures whether the tool's outcome should be surfaced to the
	// agent as an error.
	IsError bool

	// Posted reports whether a result event for this call reached the session.
	// False on a permanent 4xx or exhausted retries, and also false — with no
	// result event ever built — for a tool name this runner does not own when
	// it deliberately posts nothing and leaves the id pending for its owner.
	Posted bool
}

// SessionToolRunner attaches to a managed-agents session, executes incoming
// agent.tool_use and agent.custom_tool_use events via a local tool registry,
// and posts the matching user.tool_result / user.custom_tool_result events
// back. agent.mcp_tool_use events are left alone — MCP tools are server-side.
// It is the sessions-side counterpart to (*BetaMessageService).NewToolRunner:
// it does ONLY the tool-execution loop — attach to the event stream, reconcile
// via the events list endpoint, dispatch the registered tools, post results,
// and the idle-after-end_turn timeout.
// Lease heartbeating, work claiming, and skill download are not its concern —
// see [github.com/anthropics/anthropic-sdk-go/lib/environments.EnvironmentWorker]
// for the full self-hosted runner composition.
//
// A SessionToolRunner is NOT safe for concurrent use. All methods must be
// called from a single goroutine. Tool BetaTool.Execute handlers ARE called
// from a background goroutine, but only one at a time per runner (serial
// execution per the agent.tool_use contract).
//
// Tool lifetime is the caller's responsibility; the runner never closes tools,
// so the same Tools slice can be reused across multiple sessions.
//
// Typical usage:
//
//	r := client.Beta.Sessions.Events.NewToolRunner(ctx, sessionID, anthropic.SessionToolRunnerOptions{
//	    Tools: agenttoolset.BetaAgentToolset20260401(&agenttoolset.AgentToolContext{Workdir: "."}),
//	})
//	defer r.Close()
//	for r.Next() {
//	    call := r.Current()
//	    // observe call.Name, call.IsError, call.Posted
//	}
//	if err := r.Err(); err != nil && !errors.Is(err, anthropic.ErrSessionTerminated) {
//	    log.Print(err)
//	}
type SessionToolRunner struct {
	eventService *BetaSessionEventService
	sessionID    string
	opts         SessionToolRunnerOptions
	log          *slog.Logger
	byName       map[string]BetaTool

	// reqOpts are applied to every request the runner issues: the caller's
	// SessionToolRunnerOptions.RequestOptions plus the runner's own
	// x-stainless-helper telemetry header. Immutable after construction.
	reqOpts []option.RequestOption

	// ctx is the runner's internal context — a child of the caller's ctx —
	// used by all background loops. cancel terminates them.
	ctx    context.Context
	cancel context.CancelFunc

	// started/closed are plain bools; Next/Close/All are documented as
	// single-goroutine.
	started bool
	closed  bool

	// Set on start(); read-only after.
	group   *errgroup.Group
	results chan DispatchedToolCall

	// Tracks in-flight tool executions for the drain phase.
	inFlight sync.WaitGroup

	// Dedup sets — touched by the stream goroutine, the reconcile pass before
	// each stream reconnect, and the dispatch goroutine. `seen` prevents
	// re-dispatching the same tool_use across reconcile+stream overlaps;
	// `answered` prevents re-executing a tool whose result the server already
	// has. One mutex covers both — they're always touched together at low
	// frequency.
	mu       sync.Mutex
	seen     map[string]struct{}
	answered map[string]struct{}

	// endTurnAtNano is the time.Now().UnixNano() of the most recent
	// session.status_idle with stop_reason "end_turn" for which no newer event
	// has since arrived; 0 whenever the session is not in that state. Set by
	// the stream goroutine (and the reconcile pass), read by the idle watchdog
	// — hence atomic.
	endTurnAtNano atomic.Int64

	// idleSignal wakes the idle watchdog whenever the stream goroutine (or the
	// reconcile pass) updates endTurnAtNano, so the watchdog re-arms its timer
	// off the relevant event instead of polling endTurnAtNano on a ticker. It
	// is buffered (cap 1) and written with a non-blocking send: a single
	// pending nudge is enough, since the watchdog reads the authoritative
	// endTurnAtNano stamp once it runs. Created in start().
	idleSignal chan struct{}

	// constructErr is set only at construction and is immutable thereafter;
	// safe to read without a lock. Used to bail out of Next/start when required
	// options are missing.
	constructErr error

	// terminalErr is the first sentinel/wrapped error any loop produced.
	// Surfaced via Err() after the results channel drains.
	terminalErrMu sync.Mutex
	terminalErr   error

	// Yielded value most recently read from results.
	current DispatchedToolCall
}

// NewToolRunner returns a [SessionToolRunner] bound to ctx and this event
// service for the given managed-agents session. It is the sessions-side
// counterpart to (*BetaMessageService).NewToolRunner. sessionID is a leading
// positional argument, matching list/send/stream on the events resource. The
// first call to Next launches background goroutines (stream, dispatch, idle
// watchdog). Logger defaults to slog.Default().
//
// sessionID is validated here: if it is empty, Next returns false on the first
// call and Err returns the corresponding error. The returned pointer is never
// nil so `defer r.Close()` is always safe.
func (r *BetaSessionEventService) NewToolRunner(ctx context.Context, sessionID string, opts SessionToolRunnerOptions) *SessionToolRunner {
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}
	log = log.With(
		slog.String("component", "session-tool-runner"),
		slog.String("session_id", sessionID),
	)
	byName := make(map[string]BetaTool, len(opts.Tools))
	for _, t := range opts.Tools {
		byName[t.Name()] = t
	}
	internalCtx, cancel := context.WithCancel(ctx)
	reqOpts := make([]option.RequestOption, 0, len(opts.RequestOptions)+1)
	reqOpts = append(reqOpts, opts.RequestOptions...)
	reqOpts = append(reqOpts, option.WithHeader(sessionRunnerHelperHeader, sessionRunnerHelperValue))
	rn := &SessionToolRunner{
		eventService: r,
		sessionID:    sessionID,
		opts:         opts,
		log:          log,
		byName:       byName,
		reqOpts:      reqOpts,
		ctx:          internalCtx,
		cancel:       cancel,
		seen:         map[string]struct{}{},
		answered:     map[string]struct{}{},
	}
	if sessionID == "" {
		rn.constructErr = errors.New("anthropic: NewToolRunner requires a non-empty session id")
	}
	return rn
}

// NewSessionToolRunner is the package-level equivalent of
// (*BetaSessionEventService).NewToolRunner — useful when you have a [Client]
// value. Prefer client.Beta.Sessions.Events.NewToolRunner in new code.
func NewSessionToolRunner(ctx context.Context, client Client, sessionID string, opts SessionToolRunnerOptions) *SessionToolRunner {
	return client.Beta.Sessions.Events.NewToolRunner(ctx, sessionID, opts)
}

// Next advances the runner. Returns true if a new DispatchedToolCall is
// available via Current. Returns false when the bound context is cancelled,
// Close has been called, the session terminated, or MaxIdle elapsed (check Err
// to distinguish).
func (r *SessionToolRunner) Next() bool {
	if r.constructErr != nil {
		return false
	}
	r.start()
	select {
	case <-r.ctx.Done():
		return false
	case call, ok := <-r.results:
		if !ok {
			return false
		}
		r.current = call
		return true
	}
}

// Current returns the most recent DispatchedToolCall. Only valid after Next
// returned true.
func (r *SessionToolRunner) Current() DispatchedToolCall {
	return r.current
}

// Err returns the first non-recoverable error, or nil if iteration ended via
// the consumer's ctx cancellation or Close. Distinguish causes with errors.Is
// against [ErrSessionTerminated] and [ErrIdleTimeout]. Construction-time option
// errors are surfaced here as well.
func (r *SessionToolRunner) Err() error {
	if r.constructErr != nil {
		return r.constructErr
	}
	r.terminalErrMu.Lock()
	defer r.terminalErrMu.Unlock()
	return r.terminalErr
}

// Close stops all background goroutines and waits for the cleanup goroutine
// (which drains in-flight tools) to finish. Safe to call multiple times. Always
// returns nil — the signature satisfies io.Closer so callers can `defer
// r.Close()` uniformly. The runner never closes the tools it was given.
func (r *SessionToolRunner) Close() error {
	if r.closed {
		return nil
	}
	r.closed = true
	// Mark started so a late Next becomes a no-op.
	r.started = true
	r.cancel()
	if r.results == nil {
		return nil
	}
	for range r.results { //nolint:revive // drain channel; body intentionally empty
	}
	return nil
}

// All returns a Go 1.23 range-over-func iterator yielding DispatchedToolCall
// values. On the final iteration err carries Err.
func (r *SessionToolRunner) All() iter.Seq2[DispatchedToolCall, error] {
	return func(yield func(DispatchedToolCall, error) bool) {
		for r.Next() {
			if !yield(r.Current(), nil) {
				return
			}
		}
		if err := r.Err(); err != nil {
			yield(DispatchedToolCall{}, err)
		}
	}
}

// start lazily launches the background goroutines on the first Next call.
// Idempotent via the started flag.
func (r *SessionToolRunner) start() {
	if r.started {
		return
	}
	r.started = true
	if r.constructErr != nil {
		return
	}
	r.results = make(chan DispatchedToolCall, sessionRunnerResultsBuffer)
	r.idleSignal = make(chan struct{}, 1)
	toolUseQ := make(chan pendingToolUse, sessionRunnerToolUseQueueBuffer)
	g, gctx := errgroup.WithContext(r.ctx)
	r.group = g

	g.Go(func() error { return r.streamLoop(gctx, toolUseQ) })
	g.Go(func() error { return r.dispatchLoop(gctx, toolUseQ) })
	g.Go(func() error { return r.idleWatchdog(gctx) })

	go r.coordinate(g)
}

func (r *SessionToolRunner) coordinate(g *errgroup.Group) {
	groupErr := g.Wait()
	r.setTerminalErr(groupErr)
	r.drainInFlight()
	r.cancel()
	close(r.results)
}

// drainInFlight waits up to the configured drain timeout (opts.DrainTimeout, or
// defaultSessionRunnerDrainTimeout) for any in-flight tool executions to
// complete, then logs a warning and returns. On timeout the waiter goroutine
// below stays blocked on r.inFlight.Wait() — and pins the runner state it closes
// over — until each tool actually returns; that is bounded by the per-tool
// timeout (opts.ToolTimeout, or defaultSessionRunnerToolTimeout), not unbounded,
// but it does mean a slow tool can outlive Close.
func (r *SessionToolRunner) drainInFlight() {
	drainTimeout := r.drainTimeout()
	done := make(chan struct{})
	go func() {
		r.inFlight.Wait()
		close(done)
	}()
	t := time.NewTimer(drainTimeout)
	defer t.Stop()
	select {
	case <-done:
	case <-t.C:
		r.log.Warn("drain timeout exceeded; in-flight tools may still be running",
			slog.Duration("drain_timeout", drainTimeout))
	}
}

func (r *SessionToolRunner) setTerminalErr(err error) {
	if err == nil {
		return
	}
	r.terminalErrMu.Lock()
	defer r.terminalErrMu.Unlock()
	if r.terminalErr == nil {
		r.terminalErr = err
	}
}

// pendingToolUse is a tool-call event the stream loop (or the reconcile pass)
// has picked up and queued for dispatch. Exactly one of toolUse / customToolUse
// is populated, selected by custom: an agent.tool_use is answered with a
// user.tool_result, an agent.custom_tool_use with a user.custom_tool_result.
type pendingToolUse struct {
	custom        bool
	toolUse       BetaManagedAgentsAgentToolUseEvent
	customToolUse BetaManagedAgentsAgentCustomToolUseEvent
}

// id returns the tool-use event id — the key for the seen/answered dedup sets
// and the id the posted result event references.
func (p pendingToolUse) id() string {
	if p.custom {
		return p.customToolUse.ID
	}
	return p.toolUse.ID
}

// name returns the tool name to look up in the local registry.
func (p pendingToolUse) name() string {
	if p.custom {
		return p.customToolUse.Name
	}
	return p.toolUse.Name
}

// input returns the raw tool input the agent supplied.
func (p pendingToolUse) input() map[string]any {
	if p.custom {
		return p.customToolUse.Input
	}
	return p.toolUse.Input
}

// signalIdle nudges the idle watchdog to re-evaluate endTurnAtNano. It is a
// non-blocking send on a cap-1 channel — only the stream goroutine calls it, so
// a dropped nudge just means a wake is already pending.
func (r *SessionToolRunner) signalIdle() {
	select {
	case r.idleSignal <- struct{}{}:
	default:
	}
}

// streamLoop tails the live SSE stream with reconnect backoff. On each
// (re)connect it opens the stream FIRST and only then reconciles full history,
// so a tool_use the server emits during the reconcile window lands on the
// stream rather than in the gap between the two; the per-event seen/answered
// sets dedup events both passes observe. Both agent.tool_use and
// agent.custom_tool_use are dispatched; agent.mcp_tool_use is ignored (MCP
// tools run server-side). Closes toolUseQ on exit so dispatchLoop can return.
func (r *SessionToolRunner) streamLoop(ctx context.Context, out chan<- pendingToolUse) error {
	defer close(out)

	backoff := sessionRunnerStreamBackoffStart
	for {
		if ctx.Err() != nil {
			return nil
		}

		// Open the live stream before reconciling history: StreamEvents makes
		// its HTTP request eagerly, so once it returns the connection is
		// established and no event emitted during reconcile can slip through
		// the gap. seen/answered dedup whatever both passes pick up.
		stream := r.eventService.StreamEvents(
			ctx,
			r.sessionID,
			BetaSessionEventStreamParams{},
			r.reqOpts...,
		)
		if err := r.reconcile(ctx, out); err != nil {
			_ = stream.Close()
			if errors.Is(err, ErrSessionTerminated) {
				return err
			}
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		if ctx.Err() != nil {
			_ = stream.Close()
			return nil
		}

		// Reset backoff only once the connection has been healthy.
		connectedAt := time.Now()
		var sentinel error
		for stream.Next() {
			if time.Since(connectedAt) > sessionRunnerStreamHealthyAfter {
				backoff = sessionRunnerStreamBackoffStart
			}
			ev := stream.Current()
			if ev.Type == string(BetaManagedAgentsSessionStatusIdleEventTypeSessionStatusIdle) &&
				ev.StopReason.Type == string(BetaManagedAgentsSessionEndTurnTypeEndTurn) {
				r.endTurnAtNano.Store(time.Now().UnixNano())
			} else {
				r.endTurnAtNano.Store(0)
			}
			r.signalIdle()
			switch ev.Type {
			case string(BetaManagedAgentsAgentToolUseEventTypeAgentToolUse):
				if r.markSeen(ev.ID) {
					select {
					case <-ctx.Done():
						sentinel = ctx.Err()
					case out <- pendingToolUse{toolUse: ev.AsAgentToolUse()}:
					}
				}
			case string(BetaManagedAgentsAgentCustomToolUseEventTypeAgentCustomToolUse):
				if r.markSeen(ev.ID) {
					select {
					case <-ctx.Done():
						sentinel = ctx.Err()
					case out <- pendingToolUse{custom: true, customToolUse: ev.AsAgentCustomToolUse()}:
					}
				}
			case string(BetaManagedAgentsUserToolResultEventTypeUserToolResult):
				r.markAnswered(ev.ToolUseID)
			case string(BetaManagedAgentsUserCustomToolResultEventTypeUserCustomToolResult):
				r.markAnswered(ev.CustomToolUseID)
			case string(BetaManagedAgentsSessionStatusTerminatedEventTypeSessionStatusTerminated),
				string(BetaManagedAgentsSessionDeletedEventTypeSessionDeleted):
				r.log.Info("session terminated", slog.String("type", ev.Type))
				sentinel = ErrSessionTerminated
			}
			if sentinel != nil {
				break
			}
		}
		_ = stream.Close()
		if sentinel != nil {
			return sentinel
		}
		if err := stream.Err(); err != nil {
			if ctx.Err() != nil {
				return nil
			}
			if isFatal4xxStatus(err) {
				r.log.Error("permanent stream failure", slog.Any("error", err))
				return fmt.Errorf("stream: %w", err)
			}
			r.log.Warn("stream disconnected; reconnecting",
				slog.Any("error", err), slog.Duration("backoff", backoff))
		}
		if ctx.Err() != nil {
			return nil
		}
		// Jittered backoff so a fleet of runners reconnecting after the same
		// blip don't synchronise into a thundering herd.
		sleepCtx(ctx, jitterDuration(backoff))
		backoff = min(backoff*2, sessionRunnerStreamBackoffCap)
	}
}

// jitterDuration returns a random duration in [d/2, d) so concurrent runners
// spread their retries instead of synchronising. d <= 1 is returned unchanged.
func jitterDuration(d time.Duration) time.Duration {
	if d <= 1 {
		return d
	}
	half := d / 2
	return half + time.Duration(rand.Int64N(int64(d-half)))
}

// reconcile lists full event history and emits any tool_use the runner has not
// already posted a result for. It pairs agent.tool_use with user.tool_result
// and agent.custom_tool_use with user.custom_tool_result when deciding what is
// still unanswered.
//
// Every tool-use event is collected — and marked seen so the concurrently
// running stream loop does not also enqueue it — but the enqueue decision is
// gated on `answered`, NOT `seen`: a tool_use whose result post previously
// FAILED was never marked answered, so it is re-dispatched on the next
// reconcile instead of being silently dropped.
//
// Returns ErrSessionTerminated if the listed history contains a
// session.status_terminated or session.deleted event — the live stream
// will never replay a terminate that fired before we attached, so without
// this check streamLoop would reconnect forever against a dead session.
func (r *SessionToolRunner) reconcile(ctx context.Context, out chan<- pendingToolUse) error {
	var pending []pendingToolUse
	lastWasEndTurn := false
	pager := r.eventService.ListAutoPaging(ctx, r.sessionID,
		BetaSessionEventListParams{
			Limit: param.NewOpt(int64(1000)),
			// Reconcile assumes oldest-to-newest event order.
			Order: BetaSessionEventListParamsOrderAsc,
		},
		r.reqOpts...)
	for pager.Next() {
		ev := pager.Current()
		switch ev.Type {
		case string(BetaManagedAgentsAgentToolUseEventTypeAgentToolUse):
			r.markSeen(ev.ID)
			pending = append(pending, pendingToolUse{toolUse: ev.AsAgentToolUse()})
		case string(BetaManagedAgentsAgentCustomToolUseEventTypeAgentCustomToolUse):
			r.markSeen(ev.ID)
			pending = append(pending, pendingToolUse{custom: true, customToolUse: ev.AsAgentCustomToolUse()})
		case string(BetaManagedAgentsUserToolResultEventTypeUserToolResult):
			r.markAnswered(ev.ToolUseID)
		case string(BetaManagedAgentsUserCustomToolResultEventTypeUserCustomToolResult):
			r.markAnswered(ev.CustomToolUseID)
		case string(BetaManagedAgentsSessionStatusTerminatedEventTypeSessionStatusTerminated),
			string(BetaManagedAgentsSessionDeletedEventTypeSessionDeleted):
			// The session is already over. Any pending tool_use we
			// collected above can't be answered (the server won't accept
			// results against a terminated session) and the live stream
			// will never replay this event, so streamLoop must shut down.
			r.log.Info("reconcile: session already terminated", slog.String("type", ev.Type))
			return ErrSessionTerminated
		}
		lastWasEndTurn = ev.Type == string(BetaManagedAgentsSessionStatusIdleEventTypeSessionStatusIdle) &&
			ev.StopReason.Type == string(BetaManagedAgentsSessionEndTurnTypeEndTurn)
	}
	if err := pager.Err(); err != nil {
		r.log.Warn("reconcile list failed", slog.Any("error", err))
		return nil
	}
	var unanswered []pendingToolUse
	for _, p := range pending {
		if !r.isAnswered(p.id()) {
			unanswered = append(unanswered, p)
		}
	}
	if lastWasEndTurn && len(unanswered) == 0 {
		r.endTurnAtNano.Store(time.Now().UnixNano())
	} else {
		r.endTurnAtNano.Store(0)
	}
	r.signalIdle()
	for _, p := range unanswered {
		select {
		case <-ctx.Done():
			return nil
		case out <- p:
		}
	}
	return nil
}

// dispatchLoop reads tool-use events (agent.tool_use and agent.custom_tool_use)
// serially and executes each. The session contract guarantees one outstanding
// tool-use per session at a time, so serial execution is correct. Pushes the
// resulting DispatchedToolCall to r.results.
func (r *SessionToolRunner) dispatchLoop(ctx context.Context, in <-chan pendingToolUse) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case p, ok := <-in:
			if !ok {
				return nil
			}
			if r.isAnswered(p.id()) {
				continue
			}
			call := r.execute(ctx, p)
			r.results <- call
		}
	}
}

// maxIdle returns the configured idle-grace duration: DefaultMaxIdle when
// opts.MaxIdle is nil, otherwise the pointed-to value (<= 0 means disabled).
func (r *SessionToolRunner) maxIdle() time.Duration {
	if r.opts.MaxIdle == nil {
		return DefaultMaxIdle
	}
	return *r.opts.MaxIdle
}

// toolTimeout returns the per-tool Execute deadline: opts.ToolTimeout when set
// to a positive value, otherwise defaultSessionRunnerToolTimeout.
func (r *SessionToolRunner) toolTimeout() time.Duration {
	if r.opts.ToolTimeout > 0 {
		return r.opts.ToolTimeout
	}
	return defaultSessionRunnerToolTimeout
}

// sendTimeout returns the per-attempt tool-result send deadline:
// opts.SendTimeout when set to a positive value, otherwise
// defaultSessionRunnerSendTimeout.
func (r *SessionToolRunner) sendTimeout() time.Duration {
	if r.opts.SendTimeout > 0 {
		return r.opts.SendTimeout
	}
	return defaultSessionRunnerSendTimeout
}

// drainTimeout returns how long Close waits for in-flight tools: opts.DrainTimeout
// when set to a positive value, otherwise defaultSessionRunnerDrainTimeout.
func (r *SessionToolRunner) drainTimeout() time.Duration {
	if r.opts.DrainTimeout > 0 {
		return r.opts.DrainTimeout
	}
	return defaultSessionRunnerDrainTimeout
}

// idleWatchdog returns ErrIdleTimeout once the session has been idle with
// stop_reason "end_turn" for MaxIdle with no new events. It is event-driven:
// the stream goroutine nudges idleSignal whenever it updates endTurnAtNano, and
// the watchdog re-arms a single timer off that authoritative stamp — no polling
// ticker. MaxIdle <= 0 disables it (waits only for ctx).
func (r *SessionToolRunner) idleWatchdog(ctx context.Context) error {
	maxIdle := r.maxIdle()
	if maxIdle <= 0 {
		<-ctx.Done()
		return nil
	}
	// Start with a stopped, drained timer; it is armed only once an
	// end_turn-idle stamp is observed.
	timer := time.NewTimer(maxIdle)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()

	// rearm stops/drains the timer and, if the session is currently idle on
	// end_turn, resets it to fire when MaxIdle elapses from that stamp.
	rearm := func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		if at := r.endTurnAtNano.Load(); at != 0 {
			remaining := maxIdle - time.Since(time.Unix(0, at))
			if remaining < time.Millisecond {
				remaining = time.Millisecond
			}
			timer.Reset(remaining)
		}
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-r.idleSignal:
			rearm()
		case <-timer.C:
			at := r.endTurnAtNano.Load()
			if at != 0 && time.Since(time.Unix(0, at)) >= maxIdle {
				r.log.Info("session idle after end_turn; stopping", slog.Duration("max_idle", maxIdle))
				return ErrIdleTimeout
			}
			// A newer event moved (or cleared) the stamp between the timer
			// firing and this read — re-arm from the current state instead of
			// timing out spuriously.
			rearm()
		}
	}
}

// execute looks up the tool by name, runs it under a per-tool timeout, posts
// the matching result event (user.tool_result for an agent.tool_use,
// user.custom_tool_result for an agent.custom_tool_use), and returns the
// DispatchedToolCall to be yielded.
func (r *SessionToolRunner) execute(ctx context.Context, p pendingToolUse) DispatchedToolCall {
	id, name := p.id(), p.name()
	log := r.log.With(
		slog.String("tool", name),
		slog.String("tool_use_id", id),
		slog.Bool("custom", p.custom),
	)
	log.Info("executing tool")

	r.inFlight.Add(1)
	defer r.inFlight.Done()

	call := DispatchedToolCall{
		Custom:        p.custom,
		ToolUse:       p.toolUse,
		CustomToolUse: p.customToolUse,
		ToolUseID:     id,
		Name:          name,
	}

	rawInput, err := json.Marshal(p.input())
	if err != nil {
		log.Warn("re-encoding tool input failed", slog.Any("error", err))
		call.IsError = true
		return r.postCall(ctx, call, textOnlyResult(fmt.Sprintf("tool input could not be re-encoded: %v", err)))
	}

	var blocks []BetaToolResultBlockParamContentUnion
	tool, ok := r.byName[name]
	if !ok {
		// Skip (split-client partial fulfilment): a name this runner is not
		// registered for belongs to the other client servicing this session
		// (typically the customer's app backend handling custom tools). Post
		// NO result, do not mark it answered, and leave the tool_use_id
		// pending for its owner — claiming it would corrupt the conversation.
		// Still yield the call so the caller can observe the unowned
		// dispatch; nothing was sent, so Posted and IsError stay false and no
		// result event is populated. The id remains unanswered, so reconcile
		// keeps it out of the idle/end-turn accounting and re-surfaces it
		// after a reconnect until its owner answers it.
		log.Info("tool not owned by this runner; leaving the tool_use_id pending for its owner")
		return call
	} else {
		// Derive the per-tool timeout from the runner ctx (not
		// context.WithoutCancel) so cancelling the runner also aborts an
		// in-flight tool instead of leaving it to run out the full timeout
		// while teardown blocks on the drain.
		toolTimeout := r.toolTimeout()
		toolCtx, cancel := context.WithTimeout(ctx, toolTimeout)
		out, runErr := tool.Execute(toolCtx, rawInput)
		cancel()
		switch {
		case errors.Is(toolCtx.Err(), context.DeadlineExceeded):
			call.IsError = true
			blocks = textOnlyResult(fmt.Sprintf("tool %q timed out after %s", name, toolTimeout))
		case runErr != nil:
			call.IsError = true
			blocks = textOnlyResult(runErr.Error())
		default:
			blocks = out
			call.IsError = false
		}
	}
	if len(blocks) == 0 {
		blocks = textOnlyResult("(no output)")
	}
	return r.postCall(ctx, call, blocks)
}

// textOnlyResult wraps s as a single text tool-result block.
func textOnlyResult(s string) []BetaToolResultBlockParamContentUnion {
	return []BetaToolResultBlockParamContentUnion{{OfText: &BetaTextBlockParam{Text: s}}}
}

// toToolResultContent converts tool result blocks for a user.tool_result event
// by JSON round-tripping each block into the event content union, falling back
// to a text block holding the raw JSON when the round-trip is incomplete.
func toToolResultContent(blocks []BetaToolResultBlockParamContentUnion) []BetaManagedAgentsUserToolResultEventParamsContentUnion {
	out := make([]BetaManagedAgentsUserToolResultEventParamsContentUnion, 0, len(blocks))
	for _, b := range blocks {
		raw, err := json.Marshal(b)
		if err != nil || len(raw) == 0 || string(raw) == "null" {
			continue
		}
		var dst BetaManagedAgentsUserToolResultEventParamsContentUnion
		if json.Unmarshal(raw, &dst) == nil && roundTripComplete(&dst) {
			out = append(out, dst)
			continue
		}
		out = append(out, BetaManagedAgentsUserToolResultEventParamsContentUnion{
			OfText: &BetaManagedAgentsTextBlockParam{Text: string(raw)},
		})
	}
	return out
}

// roundTripComplete reports whether dst decoded into a variant whose required
// fields survived the JSON round-trip.
func roundTripComplete(dst *BetaManagedAgentsUserToolResultEventParamsContentUnion) bool {
	switch {
	case dst.OfText != nil:
		return true
	case dst.OfImage != nil:
		return dst.OfImage.Source.asAny() != nil
	case dst.OfDocument != nil:
		return dst.OfDocument.Source.asAny() != nil
	case dst.OfSearchResult != nil:
		return dst.OfSearchResult.Source != "" && dst.OfSearchResult.Title != ""
	default:
		return false
	}
}

// toCustomToolResultContent converts tool result blocks for a
// user.custom_tool_result event using the same round-trip as toToolResultContent.
func toCustomToolResultContent(blocks []BetaToolResultBlockParamContentUnion) []BetaManagedAgentsUserCustomToolResultEventParamsContentUnion {
	src := toToolResultContent(blocks)
	out := make([]BetaManagedAgentsUserCustomToolResultEventParamsContentUnion, 0, len(src))
	for _, b := range src {
		out = append(out, BetaManagedAgentsUserCustomToolResultEventParamsContentUnion{
			OfText:         b.OfText,
			OfImage:        b.OfImage,
			OfDocument:     b.OfDocument,
			OfSearchResult: b.OfSearchResult,
		})
	}
	return out
}

// postCall builds the matching result event for call — user.custom_tool_result
// when call.Custom, otherwise user.tool_result — from the tool's result blocks,
// sends it, and records the event on call.Result / call.CustomResult and the
// send outcome on call.Posted. The tool-use id is marked answered ONLY when
// the result post actually succeeds: a failed post leaves the call unanswered
// so the next reconcile re-dispatches it instead of silently dropping it.
func (r *SessionToolRunner) postCall(ctx context.Context, call DispatchedToolCall, blocks []BetaToolResultBlockParamContentUnion) DispatchedToolCall {
	var event BetaManagedAgentsEventParamsUnion
	if call.Custom {
		call.CustomResult = BetaManagedAgentsUserCustomToolResultEventParams{
			CustomToolUseID: call.ToolUseID,
			Type:            BetaManagedAgentsUserCustomToolResultEventParamsTypeUserCustomToolResult,
			IsError:         param.NewOpt(call.IsError),
			Content:         toCustomToolResultContent(blocks),
		}
		event = BetaManagedAgentsEventParamsUnion{OfUserCustomToolResult: &call.CustomResult}
	} else {
		call.Result = BetaManagedAgentsUserToolResultEventParams{
			ToolUseID: call.ToolUseID,
			Type:      BetaManagedAgentsUserToolResultEventParamsTypeUserToolResult,
			IsError:   param.NewOpt(call.IsError),
			Content:   toToolResultContent(blocks),
		}
		event = BetaManagedAgentsEventParamsUnion{OfUserToolResult: &call.Result}
	}
	call.Posted = r.sendResult(ctx, event, call.ToolUseID)
	if call.Posted {
		r.markAnswered(call.ToolUseID)
	}
	return call
}

// sendResult posts a tool-result event (user.tool_result or
// user.custom_tool_result). Returns true on success. Retries up to
// sessionRunnerSendRetries times on transient failures, backing off between
// attempts but NOT after the final one; a permanent 4xx short-circuits and
// returns false.
func (r *SessionToolRunner) sendResult(ctx context.Context, event BetaManagedAgentsEventParamsUnion, toolUseID string) bool {
	params := BetaSessionEventSendParams{
		Events: []BetaManagedAgentsEventParamsUnion{event},
	}
	for attempt := range sessionRunnerSendRetries {
		sendCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), r.sendTimeout())
		_, err := r.eventService.Send(sendCtx, r.sessionID, params, r.reqOpts...)
		cancel()
		if err == nil {
			return true
		}
		if isFatal4xxStatus(err) {
			r.log.Error("tool result send hit permanent 4xx; not retrying",
				slog.String("tool_use_id", toolUseID), slog.Any("error", err))
			return false
		}
		if ctx.Err() != nil {
			r.log.Debug("tool result send abandoned: ctx cancelled",
				slog.String("tool_use_id", toolUseID))
			return false
		}
		r.log.Warn("tool result send failed; retrying",
			slog.String("tool_use_id", toolUseID),
			slog.Int("attempt", attempt+1),
			slog.Any("error", err))
		// Don't back off after the final attempt — there is no retry left to
		// wait for, so sleeping would only delay returning the failure.
		if attempt < sessionRunnerSendRetries-1 {
			sleepCtx(ctx, time.Duration(attempt+1)*time.Second)
		}
	}
	return false
}

func (r *SessionToolRunner) markSeen(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.seen[id]; ok {
		return false
	}
	r.seen[id] = struct{}{}
	return true
}

func (r *SessionToolRunner) markAnswered(id string) {
	r.mu.Lock()
	r.answered[id] = struct{}{}
	r.mu.Unlock()
}

func (r *SessionToolRunner) isAnswered(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.answered[id]
	return ok
}

// sleepCtx sleeps for d or until ctx is done, whichever comes first.
func sleepCtx(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}

// isFatal4xxStatus reports whether err is a client error that will not succeed
// on retry. 408 (timeout) and 429 (rate-limited) are excluded so callers can
// back off rather than tear down.
func isFatal4xxStatus(err error) bool {
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		return false
	}
	c := apiErr.StatusCode
	return c >= 400 && c < 500 && c != 408 && c != 429
}
