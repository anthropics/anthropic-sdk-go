package requestconfig

import (
	"net/http"
	"slices"

	"github.com/anthropics/anthropic-sdk-go/internal/telemetry"
	"github.com/anthropics/anthropic-sdk-go/internal/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// OpenTelemetryConfig is what a RequestConfig holds for the OpenTelemetry span
// of its call: the options of the application, and what the SDK method that
// sends the request says of the span.
type OpenTelemetryConfig struct {
	tracing.Options
	// SpanName names the OpenTelemetry span of the call. The SDK method that sends the
	// request sets it.
	SpanName string
	// SpanRoute is the http.route of the span: the path of the endpoint as the API
	// description gives it, such as /v1/files/{file_id}. The SDK method sets it.
	SpanRoute string
	// SpanQueryParams are the names of the query parameters that the endpoint
	// defines. The span's url.query has no other parameter. The SDK method sets them.
	SpanQueryParams []string
	// SpanContentQueryParams are the ones of SpanQueryParams whose value is content
	// on this endpoint only. A url.query that has one is content. The SDK method sets them.
	SpanContentQueryParams []string
	// SpanAttributes are the resource IDs on the span that the SDK method has from
	// its own arguments: the IDs in the request's path. The request for a later
	// page of a list has them too.
	SpanAttributes []attribute.KeyValue
	// spanCreatedID returns the attribute for the ID of the object that the
	// method created, once the response is decoded, and false when there is none.
	// It is nil unless the method creates one object.
	spanCreatedID func() (attribute.KeyValue, bool)
	// spanRecorder is what the SDK method records on the span of its own. It is
	// zero for a method that has nothing of its own to record.
	spanRecorder SpanRecorder
}

// clone returns what a clone of the RequestConfig has of o: all of it but
// spanCreatedID, which reads the object that this call decodes.
func (o OpenTelemetryConfig) clone() OpenTelemetryConfig {
	o.spanCreatedID = nil
	return o
}

// SpanRecorder is what an SDK method that has attributes of its own records on
// the OpenTelemetry span of its call. Execute starts and ends the span, as it
// does for every method. A field that is nil records nothing, and neither does
// one that panics: the call goes on without what it records.
type SpanRecorder struct {
	// StartAttributes returns what the method knows of the request, whose header
	// is header, before the span starts. The span starts with it, so that a
	// sampler sees it.
	StartAttributes func(header http.Header) []attribute.KeyValue
	// RecordRequest records, once the span has started, what of the request
	// depends on whether the span records.
	RecordRequest func(call *tracing.APICall)
	// RecordResult records the result of the call, which failed with err or whose
	// response is decoded, before the span ends.
	RecordResult func(call *tracing.APICall, err error)
	// StreamObserver returns what the span follows of the events of a response
	// that the application decodes as a stream.
	StreamObserver func(call *tracing.APICall) tracing.StreamObserver
}

// WithSpanRecorder returns a RequestOption that names the OpenTelemetry span of
// the call name, and makes the span record what recorder does.
func WithSpanRecorder(name string, recorder SpanRecorder) RequestOption {
	return RequestOptionFunc(func(r *RequestConfig) error {
		r.OpenTelemetry.SpanName, r.OpenTelemetry.spanRecorder = name, recorder
		return nil
	})
}

// WithSpan returns a RequestOption that describes the OpenTelemetry span of the
// call: the span's name, the path of the endpoint as the API description gives
// it, and the names of the query parameters that the endpoint defines.
func WithSpan(name, route string, queryParams ...string) RequestOption {
	return RequestOptionFunc(func(r *RequestConfig) error {
		r.OpenTelemetry.SpanName, r.OpenTelemetry.SpanRoute, r.OpenTelemetry.SpanQueryParams = name, route, queryParams
		return nil
	})
}

// WithSpanContentQuery returns a RequestOption that lists the query parameters
// whose value is content on this endpoint only. The method passes it next to [WithSpan].
func WithSpanContentQuery(names ...string) RequestOption {
	return RequestOptionFunc(func(r *RequestConfig) error {
		r.OpenTelemetry.SpanContentQueryParams = names
		return nil
	})
}

// WithSpanAttribute returns a RequestOption that sets the attribute key of the
// call's OpenTelemetry span to value: a resource ID that the SDK method has from
// its own arguments. The method passes it next to [WithSpan], with key as a
// string literal.
func WithSpanAttribute(key, value string) RequestOption {
	return RequestOptionFunc(func(r *RequestConfig) error {
		r.OpenTelemetry.SpanAttributes = append(slices.Clip(r.OpenTelemetry.SpanAttributes), attribute.String(key, value))
		return nil
	})
}

// WithSpanCreatedID returns a RequestOption for a method that creates one object
// and decodes it into *res: the call's OpenTelemetry span gets the attribute key
// with the object's ID, which id reads from the decoded object. The attribute is
// left out when the call fails, when the application had the response decoded
// into a value of its own, which leaves *res nil, and when the object has no ID.
func WithSpanCreatedID[T any](key string, res **T, id func(*T) string) RequestOption {
	return RequestOptionFunc(func(r *RequestConfig) error {
		r.OpenTelemetry.spanCreatedID = func() (attribute.KeyValue, bool) {
			if *res == nil {
				return attribute.KeyValue{}, false
			}
			value := id(*res)
			return attribute.String(key, value), value != ""
		}
		return nil
	})
}

// sendInSpan sends the request, as send does, inside the OpenTelemetry span of
// the call, which it starts and ends. It only sends the request when the spans
// are off.
func (cfg *RequestConfig) sendInSpan() (err error) {
	ctx := cfg.Request.Context()
	recorder := cfg.OpenTelemetry.spanRecorder
	call := tracing.StartAPICallSpan(ctx, cfg.OpenTelemetry.Options, cfg.OpenTelemetry.SpanName, cfg.spanStartAttributes)
	if call == nil {
		return cfg.send(ctx, nil)
	}
	if recorder.RecordRequest != nil {
		tracing.Guard(func() { recorder.RecordRequest(call) })
	}
	// The attempts run inside the span, so that the application's HTTP
	// instrumentation parents their spans to it.
	err = cfg.send(trace.ContextWithSpan(ctx, call.Span), call)
	tracing.Guard(func() { telemetry.SetAPICallAttributes(call) })
	if err == nil && cfg.OpenTelemetry.spanCreatedID != nil {
		tracing.Guard(func() {
			if id, ok := cfg.OpenTelemetry.spanCreatedID(); ok {
				call.Span.SetAttributes(id)
			}
		})
	}
	if recorder.RecordResult != nil {
		tracing.Guard(func() { recorder.RecordResult(call, err) })
	}
	// The span of a body that the application reads ends when that read does.
	if _, ok := cfg.ResponseBodyInto.(**http.Response); !ok || err != nil {
		tracing.EndAPICallSpan(call, err)
		return err
	}
	var observer tracing.StreamObserver
	if recorder.StreamObserver != nil {
		tracing.Guard(func() { observer = recorder.StreamObserver(call) })
	}
	tracing.EndAPICallSpanWhenBodyEnds(call, observer)
	return err
}

// spanStartAttributes returns the attributes that the span of the call starts
// with. They are made only once the spans are known to be on.
func (cfg *RequestConfig) spanStartAttributes() []attribute.KeyValue {
	attrs := slices.Concat(telemetry.APICallAttributes(cfg.OpenTelemetry.Options, cfg.OpenTelemetry.SpanAttributes), tracing.ServerAttributes(cfg.BaseURL), tracing.RequestAttributes(cfg.Request.Method, cfg.OpenTelemetry.SpanRoute), cfg.OpenTelemetry.SpanAttributes)
	if start := cfg.OpenTelemetry.spanRecorder.StartAttributes; start != nil {
		tracing.Guard(func() { attrs = append(attrs, start(cfg.Request.Header)...) })
	}
	return attrs
}

// tracedDoer returns what sends one attempt through next, the HTTP client, and
// records on call the request that the client gets: the middlewares have run by
// then, so its URL is the one that the SDK sends the request to, and a redirect
// that the client follows never shows.
func tracedDoer(call *tracing.APICall, queryParams, endpointContentQueryParams []string, next middlewareNext) middlewareNext {
	return func(req *http.Request) (*http.Response, error) {
		tracing.RecordSentRequest(call, req.URL, queryParams, endpointContentQueryParams)
		return next(req)
	}
}

// tracedHandler returns what sends one attempt through next with the
// trace-context headers of the span of call, and keeps its response on call.
func tracedHandler(call *tracing.APICall, next middlewareNext) middlewareNext {
	return func(req *http.Request) (*http.Response, error) {
		tracing.InjectTraceHeaders(call, req)
		res, err := next(req)
		call.Response = res
		return res, err
	}
}
