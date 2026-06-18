// Package environments provides helpers for running self-hosted environment
// workers: the control-plane [WorkPoller] (claim work items and hand them
// back) and the [EnvironmentWorker] composition (poll + skills + run the
// session tools while heartbeating + force-stop + loop). They are designed
// to feel at home alongside the SDK's other iterators (ssestream.Stream[T],
// BetaToolRunner): pull-style Next/Current/Err/Close plus a Go 1.23
// range-over-func All().
//
// The per-session tool-execution loop itself lives next to the Messages tool
// runner as [github.com/anthropics/anthropic-sdk-go.SessionToolRunner]
// (client.Beta.Sessions.Events.NewToolRunner), and the agent_toolset_20260401
// tool implementations live in
// [github.com/anthropics/anthropic-sdk-go/tools/agenttoolset].
package environments

import (
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/stainlessheader"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

// bearerReqOpts returns the auth-only per-request options every
// runner-helper API call must carry:
//
//   - WithHeaderDel("X-Api-Key") clears the parent client's default API key
//     so it doesn't ride on the wire alongside the bearer credential.
//     option.WithAuthToken alone only *adds* an Authorization header — the
//     parent's WithAPIKey middleware still sets X-Api-Key, so both end up on
//     every request and the server's route-by-route precedence is what would
//     otherwise decide which one auth gets.
//   - WithAuthToken sets the environment key as the Authorization bearer.
//
// Use this for pass-through call sites where the downstream helper stamps
// its own x-stainless-helper header (the SessionToolRunner's stream/list/send
// path). For call sites that own their helper tag (the WorkPoller and the
// EnvironmentWorker's heartbeat/force-stop), use [helperReqOpts] instead.
//
// Returns an error on an empty environmentKey: option.WithAuthToken("")
// produces a malformed "Authorization: Bearer " header that 401s with a
// non-obvious server-side error. All callers already validate
// EnvironmentKey at the worker/poller option boundary; this guard is a
// belt-and-braces for a future caller that reaches the helpers through a
// path that misses that validation.
func bearerReqOpts(environmentKey string) ([]option.RequestOption, error) {
	if environmentKey == "" {
		return nil, errors.New("environments: environment key is required")
	}
	return []option.RequestOption{
		option.WithHeaderDel("X-Api-Key"),
		option.WithAuthToken(environmentKey),
	}, nil
}

// helperReqOpts wraps [bearerReqOpts] with an x-stainless-helper telemetry
// tag so the control plane can attribute requests to the specific helper
// that issued them. Use this for the WorkPoller and the EnvironmentWorker's
// heartbeat/force-stop; the SessionToolRunner stamps its own tag internally
// and so should pass [bearerReqOpts] through unchanged.
//
// Surfaces [bearerReqOpts]'s empty-key error to the caller.
func helperReqOpts(environmentKey string, helper stainlessheader.Value) ([]option.RequestOption, error) {
	opts, err := bearerReqOpts(environmentKey)
	if err != nil {
		return nil, err
	}
	return append(opts, stainlessheader.With(helper)), nil
}

// ManagedAgentsBeta is the anthropic-beta value gating self-hosted
// environments access. The work resource and the sessions.events.*
// resource both auto-inject this header on every call, so the runner
// does not need to set it explicitly; the constant is exported so
// consumers issuing custom requests against the same APIs can reference
// it.
//
// As of the managed-agents-2026-04-01 rollup, this beta replaces the
// older `environments-2026-03-01` header and shares the value with the
// managed-agents beta used elsewhere in the SDK.
//
// TODO(codegen): the generator emits "managed-agents-2026-04-01" as a
// string literal in every relevant method body rather than as an
// exported named constant. Source from a generated constant when
// one exists.
const ManagedAgentsBeta = "managed-agents-2026-04-01"

const (
	// pollBlockMillis is the long-poll block_ms we request from the API.
	// The server caps this at 999 — we rely on client-side jitter between
	// empty polls to spread reconnections across workers.
	pollBlockMillis = 999

	pollBackoffCap = 60 * time.Second

	// stopTimeout bounds the deferred Stop call so a slow server cannot
	// stall the consumer's next iteration or Close indefinitely.
	stopTimeout = 10 * time.Second
)

// WorkPollerOptions configures a WorkPoller.
type WorkPollerOptions struct {
	// EnvironmentID is the self-hosted environment to claim work from.
	// Required.
	EnvironmentID string

	// EnvironmentKey is the bearer token used to authorize every work.*
	// call the poller issues — Poll, Ack and Stop. Required.
	EnvironmentKey string

	// WorkerID is a stable identifier reported back to the server for
	// observability and Redis consumer-group routing. Defaults to
	// "<os.Hostname()>-<random hex>".
	WorkerID string

	// RequestOptions are applied to every work.* call the poller issues —
	// Poll, Ack and Stop — on top of the environment-key auth and
	// x-stainless-helper telemetry the poller adds itself. Use it for a
	// proxy/custom header or a base-URL override. These options are applied
	// first, so the poller's own environment-key auth and helper header take
	// precedence and cannot be clobbered by a caller option.
	RequestOptions []option.RequestOption

	// Drain controls what the poller does when the work queue is empty.
	// When false (default) the poller long-polls forever, sleeping with
	// jitter between empty polls — the long-running-runner shape. When true
	// the poller returns as soon as a poll comes back empty: Next returns
	// false and Err stays nil (normal termination, same as ctx
	// cancellation), so a webhook-driven dispatcher can drain the queue and
	// respond instead of blocking. Pair with BlockMs set to param.Null for a
	// single non-blocking pass.
	Drain bool

	// BlockMs is the long-poll block_ms forwarded to work.poll — how long
	// the server holds an empty poll open. Three states:
	//   - omitted (zero value): the poller uses the default of 999ms (the
	//     server caps block_ms at 999).
	//   - param.NewOpt(n): request block_ms=n.
	//   - param.Null[int64](): omit block_ms entirely for a non-blocking
	//     poll (the server rejects an explicit 0). Drain callers usually
	//     want this so the final empty poll returns immediately.
	BlockMs param.Opt[int64]

	// ReclaimOlderThanMs is forwarded to work.poll: reclaim un-ack'd work
	// older than this many milliseconds. Omitted (zero value) by default so
	// no reclaim is requested. Useful in drain mode so a dead runner's work
	// re-surfaces on the next poll. Set with param.NewOpt(n).
	ReclaimOlderThanMs param.Opt[int64]

	// Logger receives non-fatal warnings (poll backoff, decode/ack skip,
	// stop failure). Defaults to slog.Default().
	Logger *slog.Logger
}

// WorkPoller long-polls an environment's work queue, claims items, and posts
// ack — yielding each claimed [anthropic.BetaSelfHostedWork] to the consumer
// via Next/Current. After the consumer finishes with a yielded item (the next
// Next call or Close), the poller posts Stop for that work item. Poll, Ack and
// Stop are all authorized with the environment key.
//
// A WorkPoller is NOT safe for concurrent use. All methods must be called
// from a single goroutine. Always pair construction with a Close, ideally
// via defer, so the final yielded work item is reliably stopped.
//
// Typical usage:
//
//	poller := environments.NewWorkPoller(ctx, client, environments.WorkPollerOptions{
//	    EnvironmentID:  envID,
//	    EnvironmentKey: environmentKey,
//	})
//	defer poller.Close()
//	for poller.Next() {
//	    work := poller.Current()
//	    // handle work
//	}
//	if err := poller.Err(); err != nil {
//	    log.Fatal(err)
//	}
type WorkPoller struct {
	ctx    context.Context
	client anthropic.Client
	opts   WorkPollerOptions
	log    *slog.Logger

	current     *anthropic.BetaSelfHostedWork
	err         error  // first construction-time or per-call error, or nil
	pendingStop func() // runs on next Next or Close; nil between yields
	failures    int    // consecutive poll failures, for poll backoff
	discards    int    // consecutive unprocessable items, for discard backoff
	closed      bool
}

// NewWorkPoller returns a WorkPoller bound to ctx and client. WorkerID
// defaults to "<os.Hostname()>-<random hex>"; Logger defaults to
// slog.Default().
//
// Required option fields are validated here: if EnvironmentID or
// EnvironmentKey is empty, Next returns false on the first call and Err
// returns the corresponding error. The returned pointer is never nil so
// `defer poller.Close()` is always safe.
func NewWorkPoller(ctx context.Context, client anthropic.Client, opts WorkPollerOptions) *WorkPoller {
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}
	if opts.WorkerID == "" {
		opts.WorkerID = defaultWorkerID(log)
	}
	log = log.With(
		slog.String("component", "work-poller"),
		slog.String("environment_id", opts.EnvironmentID),
	)
	p := &WorkPoller{
		ctx:    ctx,
		client: client,
		opts:   opts,
		log:    log,
	}
	switch {
	case opts.EnvironmentID == "":
		p.err = errors.New("environments: WorkPollerOptions.EnvironmentID is required")
	case opts.EnvironmentKey == "":
		p.err = errors.New("environments: WorkPollerOptions.EnvironmentKey is required")
	}
	return p
}

// Next advances the poller. Returns true if a work item is now available
// via Current. Returns false when the bound context is cancelled, Close
// has been called, or a non-retryable error occurred (check Err).
//
// Each Next call first runs the deferred Stop for the previously yielded
// work item (if any), then long-polls until a claimable item arrives and
// acks it. An ack failure causes the item to be force-stopped and skipped,
// and polling continues.
func (p *WorkPoller) Next() bool {
	p.runPendingStop()
	if p.err != nil || p.closed {
		return false
	}

	// The environment key authorizes every work.* call the poller issues:
	// Poll, Ack and Stop. helperReqOpts also clears the parent client's
	// default X-Api-Key so it doesn't ride alongside the bearer credential.
	// Caller-supplied RequestOptions are applied first so a proxy/custom
	// header reaches Poll/Ack/Stop while the poller's own X-Api-Key delete,
	// environment-key auth and helper header (appended last) still win.
	helperOpts, err := helperReqOpts(p.opts.EnvironmentKey, stainlessheader.EnvironmentsWorkPoller)
	if err != nil {
		// NewWorkPoller already validates EnvironmentKey at construction
		// (see line 240-245), so reaching here means the validation was
		// bypassed by a future code path; surface the failure on Err()
		// and end the iterator rather than firing requests with the
		// parent client's credentials.
		p.err = err
		return false
	}
	reqOpts := make([]option.RequestOption, 0, len(p.opts.RequestOptions)+len(helperOpts))
	reqOpts = append(reqOpts, p.opts.RequestOptions...)
	reqOpts = append(reqOpts, helperOpts...)

	// block_ms: an explicit value is forwarded as-is, an explicit null omits
	// the param (non-blocking poll), and the omitted/zero case falls back to
	// the default 999 so the historical behaviour is preserved.
	pollParams := anthropic.BetaEnvironmentWorkPollParams{
		AnthropicWorkerID: param.NewOpt(p.opts.WorkerID),
	}
	switch {
	case p.opts.BlockMs.Valid():
		pollParams.BlockMs = p.opts.BlockMs
	case param.IsNull(p.opts.BlockMs):
		// non-blocking poll: leave BlockMs omitted from the request
	default:
		pollParams.BlockMs = param.NewOpt(int64(pollBlockMillis))
	}
	if p.opts.ReclaimOlderThanMs.Valid() {
		pollParams.ReclaimOlderThanMs = p.opts.ReclaimOlderThanMs
	}

	for {
		if p.ctx.Err() != nil {
			// ctx cancellation is normal termination, not an error — same
			// convention as io.EOF for Go iterators.
			return false
		}

		work, err := p.client.Beta.Environments.Work.Poll(
			p.ctx,
			p.opts.EnvironmentID,
			pollParams,
			reqOpts...,
		)
		if err != nil {
			if p.ctx.Err() != nil {
				return false
			}
			// A bad environment key or a missing environment is a 4xx that
			// will never succeed on retry — surface it instead of backing
			// off forever, matching the heartbeat/stream sibling loops.
			if isFatal4xx(err) {
				p.log.ErrorContext(p.ctx, "poll hit permanent 4xx; stopping",
					slog.Any("error", err))
				p.err = fmt.Errorf("environments: poll failed: %w", err)
				return false
			}
			p.failures++
			// Jittered backoff so a fleet of workers recovering from the
			// same blip don't synchronise their retries.
			base := backoff(p.failures)
			d := jitter(base/2, base)
			p.log.WarnContext(p.ctx, "poll failed, backing off",
				slog.Any("error", err), slog.Duration("sleep", d))
			sleep(p.ctx, d)
			continue
		}
		p.failures = 0

		if work == nil || work.ID == "" {
			if p.opts.Drain {
				// Queue empty and Drain set: end iteration normally (Err
				// stays nil) rather than sleeping and re-polling, the same
				// way ctx cancellation ends it.
				p.log.InfoContext(p.ctx, "work queue drained")
				return false
			}
			sleep(p.ctx, jitter(time.Second, 3*time.Second))
			continue
		}

		log := p.log.With(slog.String("work_id", work.ID),
			slog.String("work_type", string(work.Data.Type)))

		if _, err := p.client.Beta.Environments.Work.Ack(
			p.ctx,
			work.ID,
			anthropic.BetaEnvironmentWorkAckParams{
				EnvironmentID: p.opts.EnvironmentID,
			},
			reqOpts...,
		); err != nil {
			// The claim could not be confirmed — discard the item rather than
			// leaving it dangling, and back off before polling again.
			log.ErrorContext(p.ctx, "ack failed, discarding work item",
				slog.Any("error", err))
			p.discardUnprocessable(work, reqOpts, log)
			continue
		}
		log.InfoContext(p.ctx, "claimed work")

		p.discards = 0
		p.current = work
		p.pendingStop = p.makeStopClosure(work.ID, work.EnvironmentID, reqOpts, log)
		return true
	}
}

// Current returns the most recent claimed work item. Only valid after Next
// returned true and before the next Next call.
func (p *WorkPoller) Current() *anthropic.BetaSelfHostedWork {
	return p.current
}

// Err returns the first non-retryable error, or nil if iteration ended
// normally (consumer break, ctx cancellation, or Close). Construction-time
// option errors are surfaced here on the first Next call.
func (p *WorkPoller) Err() error {
	return p.err
}

// Close runs the deferred Stop for the last yielded work item (if any)
// and marks the poller closed. Safe to call multiple times; subsequent
// Next calls return false. Always returns nil — the signature satisfies
// io.Closer so callers can `defer poller.Close()` uniformly.
func (p *WorkPoller) Close() error {
	if p.closed {
		return nil
	}
	p.closed = true
	p.runPendingStop()
	return nil
}

// All returns a Go 1.23 range-over-func iterator yielding each claimed
// [anthropic.BetaSelfHostedWork]. On the final iteration (or on early break)
// err carries the value of Err.
//
//	for work, err := range poller.All() {
//	    if err != nil { return err }
//	    // ...
//	}
func (p *WorkPoller) All() iter.Seq2[*anthropic.BetaSelfHostedWork, error] {
	return func(yield func(*anthropic.BetaSelfHostedWork, error) bool) {
		for p.Next() {
			if !yield(p.Current(), nil) {
				return
			}
		}
		if err := p.Err(); err != nil {
			yield(nil, err)
		}
	}
}

func (p *WorkPoller) runPendingStop() {
	if p.pendingStop == nil {
		return
	}
	stop := p.pendingStop
	p.pendingStop = nil
	stop()
}

func (p *WorkPoller) makeStopClosure(workID, envID string, reqOpts []option.RequestOption, log *slog.Logger) func() {
	return func() {
		// Fresh context so a cancelled caller ctx doesn't skip the Stop.
		stopCtx, cancel := context.WithTimeout(context.Background(), stopTimeout)
		defer cancel()
		stopErr := stopWork(stopCtx, p.client, workID,
			anthropic.BetaEnvironmentWorkStopParams{EnvironmentID: envID},
			reqOpts...)
		if stopErr != nil && !isStatus(stopErr, 409) {
			log.WarnContext(stopCtx, "stop failed", slog.Any("error", stopErr))
		}
	}
}

// stopWork posts work.Stop and discards the response body.
//
// TODO: remove this helper once either (a) the OpenAPI spec for
// work.Stop is updated to declare a 204 No Content response (currently
// it declares a BetaSelfHostedWork JSON body that the server never
// sends), or (b) the Go SDK's internal/requestconfig short-circuits 204
// responses the way the TypeScript SDK does in src/internal/parse.ts.
// Today the server returns 204 with no body / no Content-Type, and the
// strict Go decoder errors with
//
//	"expected destination type of 'string' or '[]byte' for responses
//	 with content-type '' that is not 'application/json'"
//
// for what is actually a successful call. We work around it by
// rebinding the response destination to **http.Response via
// WithResponseBodyInto, which trips the bypass in internal/requestconfig
// (the executor hands the body off without parsing). TS and Python
// never hit this because their decoders handle 204/empty bodies
// gracefully; only Go is strict here.
func stopWork(ctx context.Context, client anthropic.Client, workID string, params anthropic.BetaEnvironmentWorkStopParams, opts ...option.RequestOption) error {
	var raw *http.Response
	opts = append(opts, option.WithResponseBodyInto(&raw))
	_, err := client.Beta.Environments.Work.Stop(ctx, workID, params, opts...)
	if raw != nil && raw.Body != nil {
		_ = raw.Body.Close()
	}
	return err
}

// discardUnprocessable best-effort force-stops a claimed work item the poller
// cannot hand to a consumer — its ack failed — so the item does not dangle
// and get redelivered, then backs off before the next poll so a persistently
// bad item cannot hot-loop the poller. reqOpts carries the environment key.
func (p *WorkPoller) discardUnprocessable(work *anthropic.BetaSelfHostedWork, reqOpts []option.RequestOption, log *slog.Logger) {
	stopCtx, cancel := context.WithTimeout(context.Background(), stopTimeout)
	defer cancel()
	if err := stopWork(stopCtx, p.client, work.ID,
		anthropic.BetaEnvironmentWorkStopParams{
			EnvironmentID: work.EnvironmentID,
			BetaSelfHostedWorkStopRequest: anthropic.BetaSelfHostedWorkStopRequestParam{
				Force: param.NewOpt(true),
			},
		},
		reqOpts...,
	); err != nil && !isStatus(err, 409) {
		log.WarnContext(p.ctx, "force-stop of unprocessable work item failed",
			slog.Any("error", err))
	}

	p.discards++
	base := backoff(p.discards)
	d := jitter(base/2, base)
	log.WarnContext(p.ctx, "backing off after unprocessable work item",
		slog.Duration("sleep", d))
	sleep(p.ctx, d)
}

// Compile-time assertion that *WorkPoller satisfies io.Closer.
var _ io.Closer = (*WorkPoller)(nil)

func backoff(n int) time.Duration {
	if n <= 0 {
		return time.Second
	}
	if n > 6 {
		return pollBackoffCap
	}
	if d := time.Duration(1<<n) * time.Second; d <= pollBackoffCap {
		return d
	}
	return pollBackoffCap
}

func jitter(low, high time.Duration) time.Duration {
	if high <= low {
		return low
	}
	return low + time.Duration(rand.Int64N(int64(high-low)))
}

func sleep(ctx context.Context, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}

// defaultWorkerID returns a process-unique worker id. The control plane uses
// it for observability and Redis consumer-group routing and documents that it
// must be unique, so even when os.Hostname succeeds we append a random suffix
// rather than reusing a bare (and possibly shared) hostname; if Hostname fails
// the random suffix alone still keeps the id unique.
func defaultWorkerID(log *slog.Logger) string {
	var b [8]byte
	if _, err := crand.Read(b[:]); err != nil {
		log.Debug("crypto/rand failed generating worker id suffix", slog.Any("error", err))
	}
	suffix := hex.EncodeToString(b[:])
	host, err := os.Hostname()
	if err != nil || host == "" {
		if err != nil {
			log.Debug("os.Hostname failed; using random-only worker id", slog.Any("error", err))
		}
		return "anthropic-sdk-go-runner-" + suffix
	}
	return host + "-" + suffix
}

func isStatus(err error, code int) bool {
	var apiErr *anthropic.Error
	return errors.As(err, &apiErr) && apiErr.StatusCode == code
}

// isFatal4xx reports whether err is a client error that will not succeed
// on retry. 408 (timeout) and 429 (rate-limited) are excluded so callers
// can back off rather than tear down.
func isFatal4xx(err error) bool {
	var apiErr *anthropic.Error
	if !errors.As(err, &apiErr) {
		return false
	}
	c := apiErr.StatusCode
	return c >= 400 && c < 500 && c != 408 && c != 429
}
