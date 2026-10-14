package anthropic_test

import (
	"context"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

var traceContextHeaders = []string{"traceparent", "tracestate", "baggage"}

// tracedContext returns the context of a caller whose span has a trace state and
// whose context has baggage, so a request that propagates it has all of traceContextHeaders.
func tracedContext(t *testing.T) context.Context {
	t.Helper()
	member, err := baggage.NewMember("tenant", "acme")
	if err != nil {
		t.Fatal(err)
	}
	bag, err := baggage.New(member)
	if err != nil {
		t.Fatal(err)
	}
	state, err := trace.ParseTraceState("vendor=value")
	if err != nil {
		t.Fatal(err)
	}
	return trace.ContextWithRemoteSpanContext(baggage.ContextWithBaggage(context.Background(), bag), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{1},
		SpanID:     trace.SpanID{2},
		TraceFlags: trace.FlagsSampled,
		TraceState: state,
		Remote:     true,
	}))
}

// sentHeaders returns a transport that answers with the message, and the headers of the last request it got.
func sentHeaders(contentType, body string) (*closureTransport, *http.Header) {
	sent := &http.Header{}
	rt := respond(200, contentType, body)
	return &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		*sent = req.Header.Clone()
		return rt.fn(req)
	}}, sent
}

func wantTraceContextHeaders(t *testing.T, sent http.Header, want bool) {
	t.Helper()
	for _, header := range traceContextHeaders {
		if got := sent.Values(header); (len(got) > 0) != want {
			t.Errorf("%s = %q, want sent: %t", header, got, want)
		}
	}
}

func TestOpenTelemetryPropagation(t *testing.T) {
	setPropagator(t, propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	for name, tc := range map[string]struct {
		env     string
		opts    []option.RequestOption
		headers bool
		noSpan  bool
	}{
		"default":                {headers: true},
		"on":                     {opts: []option.RequestOption{option.WithOpenTelemetryPropagation(true)}, headers: true},
		"off":                    {opts: []option.RequestOption{option.WithOpenTelemetryPropagation(false)}},
		"env":                    {env: "false"},
		"env in another case":    {env: "FALSE"},
		"env with another value": {env: "0", headers: true},
		"env with on":            {env: "false", opts: []option.RequestOption{option.WithOpenTelemetryPropagation(true)}, headers: true},
		"env true with off":      {env: "true", opts: []option.RequestOption{option.WithOpenTelemetryPropagation(false)}},
		"on after off":           {opts: []option.RequestOption{option.WithOpenTelemetryPropagation(false), option.WithOpenTelemetryPropagation(true)}, headers: true},
		"off with traces on":     {opts: []option.RequestOption{option.WithOpenTelemetryPropagation(false), option.WithOpenTelemetryTraces(true)}},
		"on without OpenTelemetry": {
			opts:   []option.RequestOption{option.WithOpenTelemetryPropagation(true), option.WithoutOpenTelemetry()},
			noSpan: true,
		},
		"on with traces off": {
			opts:   []option.RequestOption{option.WithOpenTelemetryPropagation(true), option.WithOpenTelemetryTraces(false)},
			noSpan: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Setenv("ANTHROPIC_OPEN_TELEMETRY_PROPAGATION", tc.env)
			rt, sent := sentHeaders("application/json", messageJSON)
			client, recorder := newTracedClient(t, rt, tc.opts...)

			if _, err := client.Messages.New(tracedContext(t), messageParams()); err != nil {
				t.Fatal(err)
			}
			wantTraceContextHeaders(t, *sent, tc.headers)
			if tc.noSpan {
				if got := len(recorder.Started()); got != 0 {
					t.Errorf("got %d spans, want none", got)
				}
				return
			}
			span := onlySpan(t, recorder)
			if span.Name() != "anthropic.messages.create" || span.Status().Code != codes.Ok {
				t.Errorf("span %q with status %v, want anthropic.messages.create with OK", span.Name(), span.Status().Code)
			}
			if got, want := span.Parent().SpanID(), (trace.SpanID{2}); got != want {
				t.Errorf("the span has the parent %s, want the caller's %s", got, want)
			}
			wantAttrs(t, span, map[string]any{"gen_ai.request.model": "claude-requested", "http.response.status_code": int64(200)})
		})
	}
}

func TestOpenTelemetryPropagationOffForAStream(t *testing.T) {
	setPropagator(t, propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	rt, sent := sentHeaders("text/event-stream", streamBody)
	client, recorder := newTracedClient(t, rt, option.WithOpenTelemetryPropagation(false))

	stream := client.Messages.NewStreaming(tracedContext(t), messageParams())
	for stream.Next() {
	}
	if err := stream.Err(); err != nil {
		t.Fatal(err)
	}
	wantTraceContextHeaders(t, *sent, false)
	span := onlySpan(t, recorder)
	if span.Status().Code != codes.Ok {
		t.Errorf("status %v, want OK", span.Status().Code)
	}
	wantAttrs(t, span, map[string]any{"gen_ai.response.id": "msg_1"})
}

func TestOpenTelemetryPropagationOffForOneRequest(t *testing.T) {
	setPropagator(t, propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	rt, sent := sentHeaders("application/json", messageJSON)
	client, recorder := newTracedClient(t, rt)

	if _, err := client.Messages.New(tracedContext(t), messageParams(), option.WithOpenTelemetryPropagation(false)); err != nil {
		t.Fatal(err)
	}
	wantTraceContextHeaders(t, *sent, false)

	if _, err := client.Messages.New(tracedContext(t), messageParams()); err != nil {
		t.Fatal(err)
	}
	wantTraceContextHeaders(t, *sent, true)
	if got := len(recorder.Spans()); got != 2 {
		t.Errorf("got %d ended spans, want one for each call", got)
	}
}

func TestOpenTelemetryPropagationOffKeepsTheHeadersOfTheApplication(t *testing.T) {
	setPropagator(t, propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	rt, sent := sentHeaders("application/json", messageJSON)
	client, recorder := newTracedClient(t, rt, option.WithOpenTelemetryPropagation(false))

	const traceparent = "00-11111111111111111111111111111111-2222222222222222-01"
	_, err := client.Messages.New(tracedContext(t), messageParams(),
		option.WithHeader("traceparent", traceparent),
		option.WithHeader("baggage", "user=1"))
	if err != nil {
		t.Fatal(err)
	}
	onlySpan(t, recorder)
	if got := sent.Values("traceparent"); len(got) != 1 || got[0] != traceparent {
		t.Errorf("traceparent %q, want the application's %q", got, traceparent)
	}
	if got := sent.Values("baggage"); len(got) != 1 || got[0] != "user=1" {
		t.Errorf("baggage %q, want the application's user=1", got)
	}
	if got := sent.Values("tracestate"); len(got) != 0 {
		t.Errorf("tracestate %q, want none", got)
	}
}

// Without a tracer provider there is no span to record, and the request still
// carries the baggage of the application unless the propagation is off.
func TestOpenTelemetryPropagationWithoutTracerProvider(t *testing.T) {
	setPropagator(t, propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	for name, tc := range map[string]struct {
		opts    []option.RequestOption
		baggage string
	}{
		"default": {baggage: "tenant=acme"},
		"off":     {opts: []option.RequestOption{option.WithOpenTelemetryPropagation(false)}},
	} {
		t.Run(name, func(t *testing.T) {
			rt, sent := sentHeaders("application/json", messageJSON)
			client := anthropic.NewClient(append([]option.RequestOption{
				option.WithAPIKey("key"),
				option.WithBaseURL("http://localhost:4010"),
				option.WithHTTPClient(&http.Client{Transport: rt}),
			}, tc.opts...)...)

			ctx := baggage.ContextWithBaggage(context.Background(), baggage.FromContext(tracedContext(t)))
			if _, err := client.Messages.New(ctx, messageParams()); err != nil {
				t.Fatal(err)
			}
			if got := sent.Get("baggage"); got != tc.baggage {
				t.Errorf("baggage %q, want %q", got, tc.baggage)
			}
			if got := sent.Get("traceparent"); got != "" {
				t.Errorf("traceparent %q without a tracer provider", got)
			}
		})
	}
}

// countingPropagator propagates the W3C trace context, and counts how often it injects one.
type countingPropagator struct {
	propagation.TraceContext
	injects atomic.Int32
}

func (p *countingPropagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	p.injects.Add(1)
	p.TraceContext.Inject(ctx, carrier)
}

// The span of a call that the sampler drops still has a trace context. With the
// propagation off it stays in the context of the attempts, where the application's
// HTTP instrumentation parents its spans to it, though no header carries it.
func TestOpenTelemetryPropagationOffKeepsTheContextOfADroppedSpan(t *testing.T) {
	sampledOut := func() trace.TracerProvider {
		return sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.NeverSample()))
	}
	for name, tc := range map[string]struct {
		provider    trace.TracerProvider
		propagation bool
		valid       bool
	}{
		"sampled out, propagation off":    {provider: sampledOut(), valid: true},
		"no-op provider, propagation off": {provider: noop.NewTracerProvider()},
		"sampled out, propagation on":     {provider: sampledOut(), propagation: true, valid: true},
	} {
		t.Run(name, func(t *testing.T) {
			propagator := &countingPropagator{}
			setPropagator(t, propagator)
			var inAttempt trace.SpanContext
			var traceparent string
			rt := &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				inAttempt = trace.SpanContextFromContext(req.Context())
				traceparent = req.Header.Get("traceparent")
				return respond(200, "application/json", messageJSON).fn(req)
			}}
			// The provider replaces the recorder of newTracedClient.
			client, _ := newTracedClient(t, rt, option.WithTracerProvider(tc.provider), option.WithOpenTelemetryPropagation(tc.propagation))

			if _, err := client.Messages.New(context.Background(), messageParams()); err != nil {
				t.Fatal(err)
			}
			if inAttempt.IsValid() != tc.valid || inAttempt.IsSampled() {
				t.Errorf("the attempt ran with the span context %+v, want valid: %t, and not sampled", inAttempt, tc.valid)
			}
			if tc.propagation {
				if !strings.HasPrefix(traceparent, "00-"+inAttempt.TraceID().String()+"-"+inAttempt.SpanID().String()) || !strings.HasSuffix(traceparent, "-00") {
					t.Errorf("traceparent %q, want that of the span %+v, not sampled", traceparent, inAttempt)
				}
				return
			}
			if traceparent != "" {
				t.Errorf("traceparent %q, want none", traceparent)
			}
			if got := propagator.injects.Load(); got != 0 {
				t.Errorf("the propagator injected %d times, want none", got)
			}
		})
	}
}
