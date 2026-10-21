package environments

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// startAttributes is a sampler that samples every span, and keeps the
// attributes that each span started with, by span name.
type startAttributes struct {
	mu      sync.Mutex
	started map[string]map[string]any
}

func (s *startAttributes) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	attrs := map[string]any{}
	for _, kv := range p.Attributes {
		attrs[string(kv.Key)] = kv.Value.AsInterface()
	}
	s.mu.Lock()
	s.started[p.Name] = attrs
	s.mu.Unlock()
	return sdktrace.AlwaysSample().ShouldSample(p)
}

func (s *startAttributes) Description() string { return "startAttributes" }

// newWorkSpanRecorder returns the request option that records the SDK's spans,
// the recorder, and the sampler that saw them start.
func newWorkSpanRecorder() (option.RequestOption, *tracetest.SpanRecorder, *startAttributes) {
	recorded := tracetest.NewSpanRecorder()
	sampler := &startAttributes{started: map[string]map[string]any{}}
	provider := sdktrace.NewTracerProvider(sdktrace.WithSampler(sampler), sdktrace.WithSpanProcessor(recorded))
	return option.WithTracerProvider(provider), recorded, sampler
}

func spansNamed(recorded *tracetest.SpanRecorder, name string) []sdktrace.ReadOnlySpan {
	var spans []sdktrace.ReadOnlySpan
	for _, s := range recorded.Ended() {
		if s.Name() == name {
			spans = append(spans, s)
		}
	}
	return spans
}

func spanAttributes(span sdktrace.ReadOnlySpan) map[string]any {
	attrs := map[string]any{}
	for _, kv := range span.Attributes() {
		attrs[string(kv.Key)] = kv.Value.AsInterface()
	}
	return attrs
}

// pollOnce makes server hand out work on the first poll, and run afterFirst
// before it answers every later one with no work.
func pollOnce(server *fakeWorkServer, work string, afterFirst func()) {
	var mu sync.Mutex
	polls := 0
	server.HandlePoll = func(w http.ResponseWriter, _ *http.Request) {
		mu.Lock()
		polls++
		first := polls == 1
		mu.Unlock()
		if !first {
			afterFirst()
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(work))
	}
}

// The worker's poll loop emits one span for the work item, under the span of its
// own context, and everything it does for the item is under that span: the API
// calls and the tool's span. The polls are the item's siblings.
func TestEnvironmentWorker_EmitsASpanForEachWorkItem(t *testing.T) {
	server := newFakeWorkServer(t)
	scriptWorkServer(t, server, workJSON("work_1", "env_1", "session"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// The span of the first item has ended by the time the worker polls again.
	pollOnce(server, workJSON("work_1", "env_1", "session"), cancel)
	sent := make(chan struct{})
	var sentOnce sync.Once
	server.HandleSend = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[]}`))
		sentOnce.Do(func() { close(sent) })
	}
	server.HandleStream = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("event: agent.tool_use\n" +
			`data: {"type":"agent.tool_use","id":"evt_1","name":"echo","input":{},"processed_at":"2026-05-11T12:00:00Z"}` + "\n\n"))
		w.(http.Flusher).Flush()
		select {
		case <-sent:
		case <-r.Context().Done():
			return
		}
		_, _ = w.Write([]byte("event: session.status_terminated\n" +
			`data: {"type":"session.status_terminated","id":"evt_term","processed_at":"2026-05-11T12:00:00Z"}` + "\n\n"))
		w.(http.Flusher).Flush()
	}
	traced, recorded, sampler := newWorkSpanRecorder()
	worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "env_key",
		Workdir:        t.TempDir(),
		Tools:          []anthropic.BetaTool{echoTool{}},
		RequestOptions: []option.RequestOption{traced},
		Logger:         silentLogger,
	})

	ctx, app := sdktrace.NewTracerProvider().Tracer("app").Start(ctx, "serve")
	require.NoError(t, worker.Run(ctx))
	app.End()

	works := spansNamed(recorded, "anthropic.environment.work")
	require.Len(t, works, 1)
	work := works[0]
	require.Equal(t, app.SpanContext().SpanID(), work.Parent().SpanID(), "the item's span is a child of the span in the worker's context")
	require.Equal(t, trace.SpanKindInternal, work.SpanKind())
	require.Equal(t, codes.Ok, work.Status().Code)
	want := map[string]any{
		"gen_ai.provider.name":          "anthropic",
		"anthropic.environment.id":      "env_1",
		"anthropic.environment.work.id": "work_1",
		"anthropic.session.id":          "sesn_test",
	}
	require.Equal(t, want, spanAttributes(work))
	require.Equal(t, want, sampler.started["anthropic.environment.work"], "the span starts with its attributes")
	require.Equal(t, "com.anthropic.sdk.go", work.InstrumentationScope().Name)

	for _, name := range []string{
		"anthropic.environments.work.ack",
		"anthropic.environments.work.heartbeat",
		"anthropic.environments.work.stop",
		"anthropic.sessions.retrieve",
		"anthropic.tool_use echo",
	} {
		children := spansNamed(recorded, name)
		require.NotEmpty(t, children, name)
		for _, child := range children {
			require.Equal(t, work.SpanContext().SpanID(), child.Parent().SpanID(), "%s is a child of the item's span", name)
		}
	}
	for _, s := range recorded.Ended() {
		require.Equal(t, app.SpanContext().TraceID(), s.SpanContext().TraceID(), s.Name())
		if s.Name() == "anthropic.environments.work.poll" {
			require.Equal(t, app.SpanContext().SpanID(), s.Parent().SpanID(), "a poll is a sibling of the item's span")
			continue
		}
		if s == work {
			continue
		}
		require.NotEqual(t, app.SpanContext().SpanID(), s.Parent().SpanID(), "%s is under the item's span", s.Name())
		require.False(t, s.EndTime().After(work.EndTime()), "%s ends before the item's span", s.Name())
	}
}

func TestEnvironmentWorker_WorkSpanStatus(t *testing.T) {
	for name, test := range map[string]struct {
		script        func(server *fakeWorkServer)
		wantStatus    codes.Code
		wantErrorType any
	}{
		// The session lookup fails the item: the span has the error.type of that call's span.
		"API error": {
			script: func(server *fakeWorkServer) {
				server.HandleSessionGet = func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusNotFound)
					_, _ = w.Write([]byte(`{"type":"error","error":{"type":"not_found_error","message":"no such session"}}`))
				}
			},
			wantStatus:    codes.Error,
			wantErrorType: "not_found_error",
		},
		// The item is another worker's now: it was neither finished with nor failed.
		"lost lease": {
			script: func(server *fakeWorkServer) {
				server.HandleHeartbeat = func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusPreconditionFailed)
					_, _ = w.Write([]byte(`{"type":"error","error":{"type":"invalid_request_error","message":"lease lost"}}`))
				}
				server.HandleStop = func(http.ResponseWriter, *http.Request) {
					t.Error("an item whose lease was lost must not be stopped")
				}
			},
			wantStatus: codes.Unset,
		},
		// A heartbeat that the API rejects for good ends the lease too, and the item is still stopped.
		"rejected heartbeat": {
			script: func(server *fakeWorkServer) {
				server.HandleHeartbeat = func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusForbidden)
					_, _ = w.Write([]byte(`{"type":"error","error":{"type":"permission_error","message":"not allowed"}}`))
				}
			},
			wantStatus: codes.Unset,
		},
	} {
		t.Run(name, func(t *testing.T) {
			server := newFakeWorkServer(t)
			scriptWorkServer(t, server, workJSON("work_1", "env_1", "session"))
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			pollOnce(server, workJSON("work_1", "env_1", "session"), cancel)
			test.script(server)
			traced, recorded, _ := newWorkSpanRecorder()
			worker := NewEnvironmentWorker(server.Client(), EnvironmentWorkerOptions{
				EnvironmentID:  "env_1",
				EnvironmentKey: "env_key",
				Workdir:        t.TempDir(),
				RequestOptions: []option.RequestOption{traced},
				Logger:         silentLogger,
			})
			require.NoError(t, worker.Run(ctx))

			works := spansNamed(recorded, "anthropic.environment.work")
			require.Len(t, works, 1)
			require.Equal(t, test.wantStatus, works[0].Status().Code)
			require.Equal(t, test.wantErrorType, spanAttributes(works[0])["error.type"])
		})
	}
}

// For an application that drives the poller itself, the item's span ends when it asks for the next
// item or closes the poller.
func TestWorkPoller_WorkSpanStatus(t *testing.T) {
	for name, test := range map[string]struct {
		finish     func(t *testing.T, p *WorkPoller, cancel context.CancelFunc)
		wantStatus codes.Code
	}{
		"closed": {
			finish:     func(*testing.T, *WorkPoller, context.CancelFunc) {},
			wantStatus: codes.Ok,
		},
		"next item asked for": {
			finish:     func(t *testing.T, p *WorkPoller, _ context.CancelFunc) { require.False(t, p.Next()) },
			wantStatus: codes.Ok,
		},
		"poller stopped": {
			finish:     func(_ *testing.T, _ *WorkPoller, cancel context.CancelFunc) { cancel() },
			wantStatus: codes.Unset,
		},
	} {
		t.Run(name, func(t *testing.T) {
			server := newFakeWorkServer(t)
			scriptWorkServer(t, server, workJSON("work_1", "env_1", "session"))
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			// The second poll finds the queue empty, and stops the poller.
			pollOnce(server, workJSON("work_1", "env_1", "session"), cancel)
			traced, recorded, _ := newWorkSpanRecorder()
			p := NewWorkPoller(ctx, server.Client(), WorkPollerOptions{
				EnvironmentID:  "env_1",
				EnvironmentKey: "env_key",
				RequestOptions: []option.RequestOption{traced},
				Logger:         silentLogger,
			})
			require.True(t, p.Next())
			require.Empty(t, spansNamed(recorded, "anthropic.environment.work"), "the span is open while the consumer handles the item")
			test.finish(t, p, cancel)
			require.NoError(t, p.Close())

			works := spansNamed(recorded, "anthropic.environment.work")
			require.Len(t, works, 1)
			require.False(t, works[0].Parent().IsValid(), "without a span in the poller's context the item's span is a root span")
			require.Equal(t, test.wantStatus, works[0].Status().Code)
			require.Nil(t, spanAttributes(works[0])["error.type"])
			// The poller stops the item once the consumer is done with it, in the item's trace.
			stops := spansNamed(recorded, "anthropic.environments.work.stop")
			require.Len(t, stops, 1)
			require.Equal(t, works[0].SpanContext().SpanID(), stops[0].Parent().SpanID())
			require.Equal(t, codes.Ok, stops[0].Status().Code, "the API answered the stop, with nothing to read")
		})
	}
}

// An item whose claim the poller cannot confirm fails with that call's error.
func TestWorkPoller_WorkSpanOfAnItemThatCannotBeAcknowledged(t *testing.T) {
	server := newFakeWorkServer(t)
	scriptWorkServer(t, server, workJSON("work_1", "env_1", "session"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// The poller discards the item, backs off, and polls again.
	pollOnce(server, workJSON("work_1", "env_1", "session"), cancel)
	server.HandleAck = func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"type":"error","error":{"type":"invalid_request_error","message":"already claimed"}}`))
	}
	traced, recorded, _ := newWorkSpanRecorder()
	p := NewWorkPoller(ctx, server.Client(), WorkPollerOptions{
		EnvironmentID:  "env_1",
		EnvironmentKey: "env_key",
		RequestOptions: []option.RequestOption{traced},
		Logger:         silentLogger,
	})
	require.False(t, p.Next())
	require.NoError(t, p.Close())

	works := spansNamed(recorded, "anthropic.environment.work")
	require.Len(t, works, 1)
	acks := spansNamed(recorded, "anthropic.environments.work.ack")
	require.Len(t, acks, 1)
	require.Equal(t, works[0].SpanContext().SpanID(), acks[0].Parent().SpanID())
	require.Equal(t, codes.Error, works[0].Status().Code)
	require.Equal(t, spanAttributes(acks[0])["error.type"], spanAttributes(works[0])["error.type"])
	require.Equal(t, "invalid_request_error", spanAttributes(works[0])["error.type"])
}

// brokenWorkSpans is a tracer provider of the application's that panics for the
// span of a work item, and for no other span: when the span starts, or when it
// is given its status, as panicsIn says.
type brokenWorkSpans struct {
	trace.TracerProvider
	panicsIn string
}

func (p brokenWorkSpans) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return brokenWorkTracer{p.TracerProvider.Tracer(name, opts...), p.panicsIn}
}

type brokenWorkTracer struct {
	trace.Tracer
	panicsIn string
}

func (tr brokenWorkTracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if name != "anthropic.environment.work" {
		return tr.Tracer.Start(ctx, name, opts...)
	}
	if tr.panicsIn == "Tracer.Start" {
		panic("Tracer.Start is broken")
	}
	ctx, span := tr.Tracer.Start(ctx, name, opts...)
	span = brokenWorkSpan{span}
	return trace.ContextWithSpan(ctx, span), span
}

type brokenWorkSpan struct{ trace.Span }

func (brokenWorkSpan) SetStatus(codes.Code, string) { panic("Span.SetStatus is broken") }

// A panic of the application's tracer or span stays out of the poller: it hands
// the item out and stops it as it does without the span, and a span that started ends.
func TestWorkPoller_WorkSpanThatPanics(t *testing.T) {
	for panicsIn, wantSpans := range map[string]int{"Tracer.Start": 0, "Span.SetStatus": 1} {
		t.Run(panicsIn, func(t *testing.T) {
			var mu sync.Mutex
			var handled []error
			otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
				mu.Lock()
				defer mu.Unlock()
				handled = append(handled, err)
			}))
			t.Cleanup(func() { otel.SetErrorHandler(otel.ErrorHandlerFunc(func(error) {})) })
			server := newFakeWorkServer(t)
			scriptWorkServer(t, server, workJSON("work_1", "env_1", "session"))
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			pollOnce(server, workJSON("work_1", "env_1", "session"), cancel)
			recorded := tracetest.NewSpanRecorder()
			provider := brokenWorkSpans{sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorded)), panicsIn}
			p := NewWorkPoller(ctx, server.Client(), WorkPollerOptions{
				EnvironmentID:  "env_1",
				EnvironmentKey: "env_key",
				RequestOptions: []option.RequestOption{option.WithTracerProvider(provider)},
				Logger:         silentLogger,
			})

			require.True(t, p.Next())
			require.Equal(t, "work_1", p.Current().ID)
			require.False(t, p.Next())
			require.NoError(t, p.Close())

			require.Len(t, spansNamed(recorded, "anthropic.environment.work"), wantSpans)
			require.Len(t, spansNamed(recorded, "anthropic.environments.work.stop"), 1)
			mu.Lock()
			defer mu.Unlock()
			require.Len(t, handled, 1)
			require.ErrorContains(t, handled[0], panicsIn+" is broken")
		})
	}
}
