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
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal/stainlessheader"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"golang.org/x/sync/errgroup"
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
	// countdown and it re-arms on the next "end_turn" idle. The countdown is
	// deferred while a confirmation-gated call is held or still dispatching,
	// and starts fresh once the last such call resolves. nil uses
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

	// Confirmation is the verdict that gated this call, if any. "allow"
	// means the call required user confirmation and was approved before it
	// ran; "deny" means the user denied it or the server evaluated its
	// permission to "deny" — the tool was never executed and nothing was
	// posted (Posted=false, IsError=false, no Result). Empty means the call
	// needed no confirmation.
	Confirmation string

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
// A call the server gated with evaluated_permission "ask" (the always_ask
// policy — or any unrecognized value, failing closed) is held until its
// user.tool_confirmation arrives: only an explicit "allow" runs it. A "deny"
// verdict — or a call the server already evaluated to "deny" — is never
// executed and posts nothing, but IS yielded (Confirmation="deny",
// Posted=false, no Result) so the caller can observe it. While a call is held
// the server parks the session at stop_reason "requires_action" and the
// runner waits until the verdict arrives, the session terminates, or ctx is
// cancelled — wrap in [environments.EnvironmentWorker] for lease-driven
// detach, or pass a context.WithTimeout for a wall-clock bound.
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

	// eventMu serializes idle/settlement transitions from the stream/reconcile
	// goroutine with the dispatch-side history refresh after a rejected result.
	eventMu sync.Mutex

	// Tool-call state shared by the stream/reconcile and dispatch goroutines.
	// `seen` stores each call's processed timestamp for stream/history dedup and
	// interrupt cutoff checks; `answered` tracks results accepted by the server;
	// `canceled` tracks calls invalidated by a processed session-wide interrupt.
	// The active-call fields let such an interrupt cancel local execution too.
	mu                      sync.Mutex
	seen                    map[string]time.Time
	answered                map[string]struct{}
	canceled                map[string]struct{}
	latestGlobalInterruptAt time.Time
	activeCallID            string
	activeCallCancel        context.CancelFunc
	// Confirmation gating (always_ask tools). `confirmationVerdicts` records
	// every user.tool_confirmation verdict by tool_use_id (persistent for the
	// run, like seen/answered). `awaitingConfirmation` holds the ask-gated tool
	// calls whose verdict has not yet arrived. Both are only touched from the
	// stream goroutine (streamLoop / reconcile), so no lock is needed for them.
	confirmationVerdicts map[string]string
	awaitingConfirmation map[string]pendingToolUse

	// idle owns the single stop-countdown, armed off the server's
	// session.status_idle stop_reason, including its deferral while
	// confirmation-gated tool work is outstanding.
	idle *idleClock

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
// first call to Next launches background goroutines (stream, dispatch).
// Logger defaults to slog.Default().
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
	reqOpts = append(reqOpts, stainlessheader.With(stainlessheader.SessionToolRunner))
	rn := &SessionToolRunner{
		eventService: r,
		sessionID:    sessionID,
		opts:         opts,
		log:          log,
		byName:       byName,
		reqOpts:      reqOpts,
		ctx:          internalCtx,
		cancel:       cancel,
		seen:         map[string]time.Time{},
		answered:     map[string]struct{}{},
		canceled:     map[string]struct{}{},
	}
	rn.confirmationVerdicts = map[string]string{}
	rn.awaitingConfirmation = map[string]pendingToolUse{}
	rn.idle = newIdleClock(rn.maxIdle())
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
	// Prefer already-buffered calls over cancellation so a terminal error does
	// not hide the rejected dispatch that produced it.
	select {
	case call, ok := <-r.results:
		if !ok {
			return false
		}
		r.current = call
		return true
	default:
	}
	select {
	case call, ok := <-r.results:
		if !ok {
			return false
		}
		r.current = call
		return true
	case <-r.ctx.Done():
		// coordinate may have canceled the context after buffering the final
		// call. Check once more before ending iteration.
		select {
		case call, ok := <-r.results:
			if !ok {
				return false
			}
			r.current = call
			return true
		default:
			return false
		}
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
// confirmation is the verdict that released it onto the queue — "allow" for an
// ask-gated call the user approved, "" for a call that needed no confirmation.
// Denied calls never reach the queue.
type pendingToolUse struct {
	custom        bool
	toolUse       BetaManagedAgentsAgentToolUseEvent
	customToolUse BetaManagedAgentsAgentCustomToolUseEvent
	confirmation  string
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

// idleClock is the idle stop-countdown state. idleWatchdog reads it and
// stops the runner once armedAt + maxIdle has elapsed.
//
// It also owns the deferral of that countdown: gated tool work registered via
// block — a call held for user confirmation, or a user-approved call still
// dispatching — keeps an arm pending until unblock retires the last blocker.
// The server can report an end_turn idle while a call is held (a legacy
// v1beta idle is up-converted with stop_reason "end_turn" unconditionally), so
// the countdown must never run over gated work: stopping then would drop the
// held call when its verdict later arrives, or cut the runner off before a
// released call's result can drive the next turn.
//
// Every transition runs under mu, so a disarm racing the pending-arm apply
// cannot be overwritten and resurrect a cancelled countdown.
type idleClock struct {
	maxIdle time.Duration

	wake chan struct{} // cap-1; nudges idleWatchdog to re-read the state

	mu      sync.Mutex
	armedAt time.Time // zero when disarmed
	// armPending is set when an arm found blockers outstanding; the unblock
	// retiring the last blocker applies it. Cleared by any disarm. While it is
	// set — and, more generally, while blockers is non-empty — armedAt is zero.
	armPending bool
	blockers   map[string]struct{}
}

func newIdleClock(maxIdle time.Duration) *idleClock {
	return &idleClock{maxIdle: maxIdle, wake: make(chan struct{}, 1), blockers: map[string]struct{}{}}
}

// noteEvent arms on session.status_idle with stop_reason "end_turn";
// disarms on anything else.
//
// user.tool_confirmation is neutral: it signals neither agent activity nor an
// idle, and its effect on the clock flows through block/unblock instead —
// disarming here would discard the pending arm the verdict is about to settle.
func (c *idleClock) noteEvent(eventType, stopReason string) {
	if eventType == "user.tool_confirmation" {
		return
	}
	if eventType == string(BetaManagedAgentsSessionStatusIdleEventTypeSessionStatusIdle) &&
		stopReason == string(BetaManagedAgentsSessionEndTurnTypeEndTurn) {
		c.arm()
	} else {
		c.disarm()
	}
}

// arm (re)starts the countdown from now and wakes the watchdog — or, while
// blockers are outstanding, holds the arm pending instead.
func (c *idleClock) arm() {
	if c.maxIdle <= 0 {
		return
	}
	c.mu.Lock()
	if len(c.blockers) > 0 {
		// Any running countdown is stale: the arm defers because gated work is
		// outstanding, so the session is not idly waiting to stop.
		c.armPending = true
		was := !c.armedAt.IsZero()
		c.armedAt = time.Time{}
		c.mu.Unlock()
		if was {
			c.signal()
		}
		return
	}
	c.armPending = false
	c.armedAt = time.Now()
	c.mu.Unlock()
	c.signal()
}

// disarm cancels the countdown and any pending arm. Blockers persist — they
// track real outstanding work, retired only by unblock.
func (c *idleClock) disarm() {
	c.mu.Lock()
	c.armPending = false
	was := !c.armedAt.IsZero()
	c.armedAt = time.Time{}
	c.mu.Unlock()
	if was {
		c.signal()
	}
}

// block registers gated work that must finish before an idle countdown starts.
func (c *idleClock) block(id string) {
	c.mu.Lock()
	c.blockers[id] = struct{}{}
	// Defensive: a blocker taken while the countdown runs converts it into a
	// pending arm, so a stale end_turn cannot stop the runner mid-gate.
	was := !c.armedAt.IsZero()
	if was {
		c.armPending = true
		c.armedAt = time.Time{}
	}
	c.mu.Unlock()
	if was {
		c.signal()
	}
}

// unblock retires gated work (a no-op for ids never blocked) and applies a
// pending arm once the last blocker retires — the countdown then runs a full
// fresh window from now. The apply happens inside the critical section:
// deciding under the lock but stamping after releasing it would let a
// concurrent disarm (a new stream event — the session is not idle) land in the
// gap and be overwritten, resurrecting a cancelled countdown.
func (c *idleClock) unblock(id string) {
	c.mu.Lock()
	delete(c.blockers, id)
	fire := len(c.blockers) == 0 && c.armPending
	if fire {
		c.armPending = false
		c.armedAt = time.Now()
	}
	c.mu.Unlock()
	if fire {
		c.signal()
	}
}

func (c *idleClock) signal() {
	select {
	case c.wake <- struct{}{}:
	default:
	}
}

// snapshot returns the current arm timestamp for the watchdog.
func (c *idleClock) snapshot() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.armedAt
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
			r.eventMu.Lock()
			r.idle.noteEvent(ev.Type, ev.StopReason.Type)
			switch ev.Type {
			case "user.interrupt":
				r.processGlobalInterrupt(ev.ProcessedAt, ev.SessionThreadID)
			case string(BetaManagedAgentsUserToolResultEventTypeUserToolResult):
				r.noteAnswered(ev.ToolUseID)
			case string(BetaManagedAgentsUserCustomToolResultEventTypeUserCustomToolResult):
				r.noteAnswered(ev.CustomToolUseID)
			}
			r.eventMu.Unlock()
			switch ev.Type {
			case string(BetaManagedAgentsAgentToolUseEventTypeAgentToolUse):
				if r.markSeen(ev.ID, ev.ProcessedAt) {
					r.routeToolEvent(ctx, out, pendingToolUse{toolUse: ev.AsAgentToolUse()})
				}
			case string(BetaManagedAgentsAgentCustomToolUseEventTypeAgentCustomToolUse):
				if r.markSeen(ev.ID, ev.ProcessedAt) {
					r.routeToolEvent(ctx, out, pendingToolUse{custom: true, customToolUse: ev.AsAgentCustomToolUse()})
				}
			case "user.tool_confirmation":
				r.noteConfirmation(ctx, out, ev.ToolUseID, string(ev.Result))
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
// gated on whether the call is settled, NOT merely seen. Calls canceled by a
// processed session-wide interrupt are settled locally because the server no
// longer accepts their results; transient send failures remain eligible for a
// later reconcile.
//
// Returns ErrSessionTerminated if the listed history contains a
// session.status_terminated or session.deleted event — the live stream
// will never replay a terminate that fired before we attached, so without
// this check streamLoop would reconnect forever against a dead session.
func (r *SessionToolRunner) reconcile(ctx context.Context, out chan<- pendingToolUse) error {
	// Reconciliation replaces the runner's view of session idleness. Clear any
	// prior pending arm before interrupts/results retire confirmation blockers.
	r.eventMu.Lock()
	r.idle.disarm()
	r.eventMu.Unlock()
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
			r.markSeen(ev.ID, ev.ProcessedAt)
			pending = append(pending, pendingToolUse{toolUse: ev.AsAgentToolUse()})
		case string(BetaManagedAgentsAgentCustomToolUseEventTypeAgentCustomToolUse):
			r.markSeen(ev.ID, ev.ProcessedAt)
			pending = append(pending, pendingToolUse{custom: true, customToolUse: ev.AsAgentCustomToolUse()})
		case "user.tool_confirmation":
			// Record the verdict only, before the routing pass below, so a
			// tool call whose confirmation appears later in the same history
			// is routed with its verdict already known. Skip settled calls so
			// we don't re-record their verdict on every reconcile.
			if !r.isSettled(ev.ToolUseID) {
				r.confirmationVerdicts[ev.ToolUseID] = string(ev.Result)
			}
		case "user.interrupt":
			// Event history is ordered by created_at, which can differ from
			// processed_at for queued user events. Keep a processed-time cutoff
			// instead of assuming every later history entry survived the interrupt.
			r.eventMu.Lock()
			r.processGlobalInterrupt(ev.ProcessedAt, ev.SessionThreadID)
			r.eventMu.Unlock()
		case string(BetaManagedAgentsUserToolResultEventTypeUserToolResult):
			r.eventMu.Lock()
			r.noteAnswered(ev.ToolUseID)
			r.eventMu.Unlock()
		case string(BetaManagedAgentsUserCustomToolResultEventTypeUserCustomToolResult):
			r.eventMu.Lock()
			r.noteAnswered(ev.CustomToolUseID)
			r.eventMu.Unlock()
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
		if !r.isSettled(p.id()) {
			unanswered = append(unanswered, p)
		}
	}
	for _, p := range unanswered {
		r.routeToolEvent(ctx, out, p)
	}
	// A held call's verdict is normally applied by the routing pass above; if
	// its tool_use fell outside the listed window the pass never saw it, so
	// apply the recorded verdict to the held copy here.
	for id, held := range r.awaitingConfirmation {
		if r.isSettled(id) {
			r.cancelHeldCall(id)
			continue
		}
		if v, ok := r.confirmationVerdicts[id]; ok {
			r.applyVerdict(ctx, out, held, v)
		}
	}
	// Routing resolved denied calls in place (marking them answered) and held
	// ask-gated calls for their verdict. If the last event in history is an
	// end_turn idle and no tool work is outstanding, the session is done — arm
	// the stop-countdown so the runner stops even if that end_turn arrived
	// during a disconnect. A still-held call doesn't count as outstanding: the
	// clock blocks on it, so the arm is held pending until its verdict lands.
	if lastWasEndTurn && len(r.outstanding(unanswered)) == 0 {
		r.eventMu.Lock()
		r.idle.arm()
		r.eventMu.Unlock()
	}
	return nil
}

// outstanding returns the reconciled tool calls that still owe the session a
// result: neither answered (a denial resolves a call in place) nor held
// awaiting a user confirmation.
func (r *SessionToolRunner) outstanding(unanswered []pendingToolUse) []pendingToolUse {
	var out []pendingToolUse
	for _, p := range unanswered {
		id := p.id()
		if _, held := r.awaitingConfirmation[id]; held {
			continue
		}
		if r.isSettled(id) {
			continue
		}
		out = append(out, p)
	}
	return out
}

// dispatchLoop reads tool-use events (agent.tool_use and agent.custom_tool_use)
// serially and executes each. The session contract guarantees one outstanding
// tool-use per session at a time, so serial execution is correct. Pushes the
// resulting DispatchedToolCall to r.results.
//
// Each call is unblocked once it is fully disposed of — executed, or moot
// because it was answered elsewhere — retiring the idle-clock blocker
// applyVerdict took on a user-approved call. unblock is a no-op for a call
// that was never gated.
func (r *SessionToolRunner) dispatchLoop(ctx context.Context, in <-chan pendingToolUse) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case p, ok := <-in:
			if !ok {
				return nil
			}
			// Retire the blocker only once the call is fully disposed of —
			// executed and surfaced, or moot because it was answered
			// elsewhere. Surfacing can block on a full results buffer, and
			// the countdown must not start while the call is still queued.
			if err := func() error {
				defer r.idle.unblock(p.id())
				if r.isSettled(p.id()) {
					return nil
				}
				call, surface, err := r.execute(ctx, p)
				if surface {
					r.surfaceCall(ctx, call)
				}
				return err
			}(); err != nil {
				return err
			}
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

// idleWatchdog returns ErrIdleTimeout once the session has been idle past
// the bound matching its stop_reason. It is event-driven: the stream
// goroutine nudges idle.wake whenever it updates the clock, and the
// watchdog re-arms a single timer off that authoritative stamp — no
// polling ticker.
func (r *SessionToolRunner) idleWatchdog(ctx context.Context) error {
	timer := time.NewTimer(time.Hour)
	if !timer.Stop() {
		<-timer.C
	}
	defer timer.Stop()

	rearm := func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		if at := r.idle.snapshot(); !at.IsZero() {
			remaining := r.idle.maxIdle - time.Since(at)
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
		case <-r.idle.wake:
			rearm()
		case <-timer.C:
			at := r.idle.snapshot()
			if !at.IsZero() && time.Since(at) >= r.idle.maxIdle {
				r.log.Info("session idle after end_turn; stopping",
					slog.Duration("max_idle", r.idle.maxIdle))
				return ErrIdleTimeout
			}
			rearm()
		}
	}
}

// execute looks up the tool by name, runs it under a per-tool timeout, posts
// the matching result event (user.tool_result for an agent.tool_use,
// user.custom_tool_result for an agent.custom_tool_use), and returns the
// DispatchedToolCall to be yielded.
func (r *SessionToolRunner) execute(ctx context.Context, p pendingToolUse) (DispatchedToolCall, bool, error) {
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
		Confirmation:  p.confirmation,
	}
	callCtx, cancelCall := context.WithCancel(ctx)
	if !r.startCall(id, cancelCall) {
		cancelCall()
		return call, false, nil
	}
	defer func() {
		r.finishCall(id)
		cancelCall()
	}()

	rawInput, err := json.Marshal(p.input())
	if err != nil {
		log.Warn("re-encoding tool input failed", slog.Any("error", err))
		call.IsError = true
		call, err = r.postCall(callCtx, call, textOnlyResult(fmt.Sprintf("tool input could not be re-encoded: %v", err)))
		return call, true, err
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
		return call, true, nil
	} else {
		// Derive the per-tool timeout from the runner ctx (not
		// context.WithoutCancel) so cancelling the runner also aborts an
		// in-flight tool instead of leaving it to run out the full timeout
		// while teardown blocks on the drain.
		toolTimeout := r.toolTimeout()
		toolCtx, cancel := context.WithTimeout(callCtx, toolTimeout)
		out, runErr := tool.Execute(toolCtx, rawInput)
		toolCtxErr := toolCtx.Err()
		cancel()
		if r.isCanceled(id) {
			return call, true, nil
		}
		switch {
		case errors.Is(toolCtxErr, context.DeadlineExceeded):
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
	call, err = r.postCall(callCtx, call, blocks)
	return call, true, err
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
// the result post actually succeeds. Transient send exhaustion leaves the call
// pending for reconcile; a permanent rejection is returned as a terminal
// runner error.
func (r *SessionToolRunner) postCall(ctx context.Context, call DispatchedToolCall, blocks []BetaToolResultBlockParamContentUnion) (DispatchedToolCall, error) {
	if r.isSettled(call.ToolUseID) {
		return call, nil
	}
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
	var err error
	call.Posted, err = r.sendResult(ctx, event, call.ToolUseID)
	if err != nil {
		return call, err
	}
	if call.Posted {
		r.markAnswered(call.ToolUseID)
	}
	return call, nil
}

// sendResult posts a tool-result event (user.tool_result or
// user.custom_tool_result). Returns true on success. Retries up to
// sessionRunnerSendRetries times on transient failures, backing off between
// attempts but NOT after the final one; a permanent 4xx short-circuits and
// returns a terminal error so the worker releases the session instead of
// redispatching an unrecoverable call forever.
func (r *SessionToolRunner) sendResult(ctx context.Context, event BetaManagedAgentsEventParamsUnion, toolUseID string) (bool, error) {
	params := BetaSessionEventSendParams{
		Events: []BetaManagedAgentsEventParamsUnion{event},
	}
	for attempt := range sessionRunnerSendRetries {
		if r.isSettled(toolUseID) {
			return false, nil
		}
		sendCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), r.sendTimeout())
		_, err := r.eventService.Send(sendCtx, r.sessionID, params, r.reqOpts...)
		cancel()
		if err == nil {
			return true, nil
		}
		if r.isSettled(toolUseID) {
			r.log.Debug("tool result send abandoned: call settled elsewhere",
				slog.String("tool_use_id", toolUseID))
			return false, nil
		}
		if ctx.Err() != nil {
			r.log.Debug("tool result send abandoned: ctx cancelled",
				slog.String("tool_use_id", toolUseID))
			return false, nil
		}
		if isFatal4xxStatus(err) {
			if r.refreshRejectedCall(ctx, toolUseID) {
				r.log.Debug("tool result rejection matched settled history",
					slog.String("tool_use_id", toolUseID))
				return false, nil
			}
			if ctx.Err() != nil {
				return false, nil
			}
			r.log.Error("tool result send hit permanent 4xx; not retrying",
				slog.String("tool_use_id", toolUseID), slog.Any("error", err))
			return false, fmt.Errorf("tool result send: %w", err)
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
	return false, nil
}

// refreshRejectedCall resolves the race where the server processes an
// interrupt or another runner's result before the local SSE stream observes it.
// It is only called after a permanent result rejection, so the extra list
// request stays off the normal dispatch path.
func (r *SessionToolRunner) refreshRejectedCall(ctx context.Context, toolUseID string) bool {
	r.mu.Lock()
	_, seen := r.seen[toolUseID]
	r.mu.Unlock()
	if !seen {
		return false
	}

	listCtx, cancel := context.WithTimeout(ctx, r.sendTimeout())
	defer cancel()
	pager := r.eventService.ListAutoPaging(listCtx, r.sessionID,
		BetaSessionEventListParams{
			Limit: param.NewOpt(int64(1000)),
			Order: BetaSessionEventListParamsOrderAsc,
			Types: []string{
				"user.interrupt",
				string(BetaManagedAgentsUserToolResultEventTypeUserToolResult),
				string(BetaManagedAgentsUserCustomToolResultEventTypeUserCustomToolResult),
			},
		},
		r.reqOpts...)
	var latestGlobalInterruptAt time.Time
	answeredIDs := map[string]struct{}{}
	for pager.Next() {
		ev := pager.Current()
		switch ev.Type {
		case "user.interrupt":
			if ev.SessionThreadID == "" && ev.ProcessedAt.After(latestGlobalInterruptAt) {
				latestGlobalInterruptAt = ev.ProcessedAt
			}
		case string(BetaManagedAgentsUserToolResultEventTypeUserToolResult):
			answeredIDs[ev.ToolUseID] = struct{}{}
		case string(BetaManagedAgentsUserCustomToolResultEventTypeUserCustomToolResult):
			answeredIDs[ev.CustomToolUseID] = struct{}{}
		}
	}

	r.eventMu.Lock()
	answeredAdvanced := false
	for id := range answeredIDs {
		if r.markAnswered(id) {
			answeredAdvanced = true
		}
	}
	var canceledIDs []string
	var cancelActive context.CancelFunc
	interruptAdvanced := false
	if !latestGlobalInterruptAt.IsZero() {
		var ok bool
		canceledIDs, cancelActive, ok, interruptAdvanced = r.recordGlobalInterrupt(latestGlobalInterruptAt, "")
		if !ok {
			interruptAdvanced = false
		}
	}
	if answeredAdvanced || interruptAdvanced {
		// These non-idle events have not reached the local stream yet. Clear
		// stale idle state before retiring blockers; a later idle event arms a
		// fresh countdown. eventMu prevents a stream idle from racing past this.
		r.idle.disarm()
	}
	if cancelActive != nil {
		cancelActive()
	}
	// awaitingConfirmation belongs to the stream goroutine. Retire blockers
	// here; the next stream/reconcile pass removes the corresponding map entries.
	for id := range answeredIDs {
		r.idle.unblock(id)
	}
	for _, id := range canceledIDs {
		r.idle.unblock(id)
	}
	r.eventMu.Unlock()

	if err := pager.Err(); err != nil {
		r.log.Warn("tool result rejection history check failed",
			slog.String("tool_use_id", toolUseID), slog.Any("error", err))
	}
	return r.isSettled(toolUseID)
}

func (r *SessionToolRunner) markSeen(id string, processedAt time.Time) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.seen[id]; ok {
		return false
	}
	r.seen[id] = processedAt
	if !r.latestGlobalInterruptAt.IsZero() && !processedAt.After(r.latestGlobalInterruptAt) {
		r.canceled[id] = struct{}{}
	}
	return true
}

func (r *SessionToolRunner) markAnswered(id string) bool {
	r.mu.Lock()
	_, existed := r.answered[id]
	r.answered[id] = struct{}{}
	r.mu.Unlock()
	return !existed
}

// noteAnswered records a result observed from the server and retires any local
// execution or confirmation hold for the same call. Only the stream goroutine
// calls this method, so awaitingConfirmation needs no additional lock.
func (r *SessionToolRunner) noteAnswered(id string) {
	r.mu.Lock()
	r.answered[id] = struct{}{}
	var cancelActive context.CancelFunc
	if r.activeCallID == id {
		cancelActive = r.activeCallCancel
	}
	r.mu.Unlock()

	if cancelActive != nil {
		cancelActive()
	}
	r.cancelHeldCall(id)
}

func (r *SessionToolRunner) cancelHeldCall(id string) {
	if _, held := r.awaitingConfirmation[id]; held {
		delete(r.awaitingConfirmation, id)
		delete(r.confirmationVerdicts, id)
		r.idle.unblock(id)
	}
}

// processGlobalInterrupt advances the session-wide processed-time cutoff. Calls
// at or before the cutoff are no longer accepted by the server, regardless of
// where they appear in created_at-ordered history.
func (r *SessionToolRunner) processGlobalInterrupt(processedAt time.Time, sessionThreadID string) {
	canceledIDs, cancelActive, ok, _ := r.recordGlobalInterrupt(processedAt, sessionThreadID)
	if !ok {
		return
	}

	if cancelActive != nil {
		cancelActive()
	}
	for _, id := range canceledIDs {
		r.cancelHeldCall(id)
	}
}

// recordGlobalInterrupt applies the synchronized portion of a global interrupt
// transition. It can run from the stream loop or the dispatch-side rejection
// refresh; stream-owned confirmation maps are cleaned separately.
func (r *SessionToolRunner) recordGlobalInterrupt(processedAt time.Time, sessionThreadID string) ([]string, context.CancelFunc, bool, bool) {
	if processedAt.IsZero() || sessionThreadID != "" {
		return nil, nil, false, false
	}

	r.mu.Lock()
	advanced := processedAt.After(r.latestGlobalInterruptAt)
	if advanced {
		r.latestGlobalInterruptAt = processedAt
	}
	var canceledIDs []string
	for id, callProcessedAt := range r.seen {
		if callProcessedAt.After(r.latestGlobalInterruptAt) {
			continue
		}
		if _, answered := r.answered[id]; answered {
			continue
		}
		r.canceled[id] = struct{}{}
		canceledIDs = append(canceledIDs, id)
	}
	var cancelActive context.CancelFunc
	if _, canceled := r.canceled[r.activeCallID]; canceled {
		cancelActive = r.activeCallCancel
	}
	r.mu.Unlock()
	return canceledIDs, cancelActive, true, advanced
}

func (r *SessionToolRunner) startCall(id string, cancel context.CancelFunc) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, answered := r.answered[id]; answered {
		return false
	}
	if _, canceled := r.canceled[id]; canceled {
		return false
	}
	r.activeCallID = id
	r.activeCallCancel = cancel
	return true
}

func (r *SessionToolRunner) finishCall(id string) {
	r.mu.Lock()
	if r.activeCallID == id {
		r.activeCallID = ""
		r.activeCallCancel = nil
	}
	r.mu.Unlock()
}

func (r *SessionToolRunner) isAnswered(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.answered[id]
	return ok
}

func (r *SessionToolRunner) isSettled(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, answered := r.answered[id]
	_, canceled := r.canceled[id]
	return answered || canceled
}

func (r *SessionToolRunner) isCanceled(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, canceled := r.canceled[id]
	return canceled
}

// routeToolEvent enqueues p for dispatch, honoring its evaluated permission.
// A call the server gated with "ask" (or any unrecognized permission — fail
// closed) is held until its user.tool_confirmation arrives; a call already
// evaluated to "deny" is resolved as denied regardless of any recorded verdict.
func (r *SessionToolRunner) routeToolEvent(ctx context.Context, out chan<- pendingToolUse, p pendingToolUse) {
	id := p.id()
	if r.isSettled(id) {
		return
	}
	// Only builtin agent.tool_use carries evaluated_permission today; if the
	// field ever lands on agent.custom_tool_use the gate must keep failing
	// closed rather than dispatch a gated call by event type.
	perm := ""
	if !p.custom {
		perm = string(p.toolUse.EvaluatedPermission)
	}
	// Comma-ok, not the zero value: a verdict this SDK doesn't recognise —
	// including an empty one — must reach applyVerdict and fail closed, not
	// read as "no verdict yet" and dispatch.
	verdict, hasVerdict := r.confirmationVerdicts[id]
	if perm == "deny" {
		// A server-side deny overrides any (stray) recorded verdict.
		verdict, hasVerdict = "deny", true
	}
	if !hasVerdict {
		if perm == "" || perm == "allow" {
			select {
			case <-ctx.Done():
			case out <- p:
			}
		} else if _, held := r.awaitingConfirmation[id]; !held {
			// "ask" — or a permission this SDK doesn't recognise — waits for
			// the user's verdict. (Already held: a reconcile after a reconnect
			// re-routes the call; keep the existing hold.)
			r.log.Info("tool call awaiting confirmation; holding",
				slog.String("tool", p.name()), slog.String("tool_use_id", id))
			r.awaitingConfirmation[id] = p
			r.idle.block(id)
			// A dispatch-side history refresh can apply an interrupt between
			// the initial settled check and this hold registration.
			if r.isSettled(id) {
				r.cancelHeldCall(id)
			}
		}
		return
	}
	r.applyVerdict(ctx, out, p, verdict)
}

// noteConfirmation records an allow/deny verdict and releases the held call it
// gates, if any. A verdict for a call this runner has not seen yet (or one it
// never gates, e.g. an agent.mcp_tool_use) is kept in confirmationVerdicts so
// a later route of that call resolves instantly.
func (r *SessionToolRunner) noteConfirmation(ctx context.Context, out chan<- pendingToolUse, toolUseID, result string) {
	if r.isSettled(toolUseID) {
		return
	}
	r.confirmationVerdicts[toolUseID] = result
	if held, ok := r.awaitingConfirmation[toolUseID]; ok {
		r.applyVerdict(ctx, out, held, result)
	}
}

// applyVerdict dispatches or resolves a gated call according to its verdict.
// Only an explicit "allow" releases the call to execute; anything else fails
// closed as a denial. A denied call is marked answered and yielded directly
// (nothing ran, nothing posted) — it never reaches the dispatch queue.
//
// The idle-clock blocker accounting lives here: a denial retires the held
// call's blocker, while an allow keeps one on the call — taking it now if the
// verdict was already known when the call was routed, so no hold was ever
// taken — until the dispatch loop has finished with it. The countdown must not
// run over gated work that is still in flight.
func (r *SessionToolRunner) applyVerdict(ctx context.Context, out chan<- pendingToolUse, p pendingToolUse, verdict string) {
	id := p.id()
	_, wasHeld := r.awaitingConfirmation[id]
	delete(r.awaitingConfirmation, id)
	if verdict == "allow" {
		r.log.Info("tool call confirmed", slog.String("tool", p.name()), slog.String("tool_use_id", id))
		if !wasHeld {
			r.idle.block(id)
		}
		p.confirmation = "allow"
		select {
		case <-ctx.Done():
		case out <- p:
		}
		return
	}
	if wasHeld {
		r.idle.unblock(id)
	}
	r.markAnswered(id)
	r.log.Info("tool call denied; not executing", slog.String("tool", p.name()), slog.String("tool_use_id", id))
	r.surfaceCall(ctx, DispatchedToolCall{
		Custom:        p.custom,
		ToolUse:       p.toolUse,
		CustomToolUse: p.customToolUse,
		ToolUseID:     id,
		Name:          p.name(),
		IsError:       false,
		Posted:        false,
		Confirmation:  "deny",
	})
}

// surfaceCall yields call to the consumer, tolerating a consumer that left
// early (ctx cancelled). The underlying work already happened; only the
// observability event is lost.
func (r *SessionToolRunner) surfaceCall(ctx context.Context, call DispatchedToolCall) {
	// Prefer an immediately available buffer slot over cancellation. A sibling
	// loop can terminate the errgroup just as this completed call is surfaced.
	select {
	case r.results <- call:
		return
	default:
	}
	select {
	case r.results <- call:
	case <-ctx.Done():
		// Cancellation and buffer availability can become ready together.
		// Preserve the call when possible without blocking shutdown.
		select {
		case r.results <- call:
		default:
		}
	}
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
