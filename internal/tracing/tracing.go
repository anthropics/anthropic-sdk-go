package tracing

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal"
	"github.com/anthropics/anthropic-sdk-go/internal/apierror"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationScope = "com.anthropic.sdk.go"

// Attribute keys of the OpenTelemetry semantic conventions, named as the semconv
// packages of go.opentelemetry.io/otel name them. They are copied rather than
// imported so that the SDK does not pin a version of the conventions.
const (
	ErrorTypeKey              = attribute.Key("error.type")
	HTTPRequestMethodKey      = attribute.Key("http.request.method")
	HTTPRequestResendCountKey = attribute.Key("http.request.resend_count")
	HTTPResponseStatusCodeKey = attribute.Key("http.response.status_code")
	HTTPRouteKey              = attribute.Key("http.route")
	ServerAddressKey          = attribute.Key("server.address")
	ServerPortKey             = attribute.Key("server.port")
	URLFullKey                = attribute.Key("url.full")
	URLPathKey                = attribute.Key("url.path")
	URLQueryKey               = attribute.Key("url.query")
	URLSchemeKey              = attribute.Key("url.scheme")
)

// DefaultMaxContentBytes is how many bytes the content attributes of one span
// hold at most, apart from the emptied lists and messages that stay when their
// content is removed, when no option and no environment variable sets a
// maximum. Sizes are string lengths, which in Go are the bytes of the
// serialized values.
const DefaultMaxContentBytes = 750_000

// Options are the OpenTelemetry options of a client or a request.
type Options struct {
	// Traces are the options of the spans.
	Traces TracesOptions
	// Propagation turns the trace-context headers of the requests on or off.
	// When it is nil, the environment variable
	// ANTHROPIC_OPEN_TELEMETRY_PROPAGATION decides.
	Propagation *bool
	// TracerProvider creates the spans, in place of the provider registered
	// with the OpenTelemetry API, when it is not nil.
	TracerProvider trace.TracerProvider
	// ProviderName is the gen_ai.provider.name of the spans when a third-party
	// platform serves the API. The platform's client options set it.
	ProviderName string
	// ConversationID is the gen_ai.conversation.id of the spans, which links the
	// calls of one conversation. The SDK does not send it to the API.
	ConversationID string
	// Platform is the last part of the name of the instrumentation scope of the
	// spans when the client of a third-party platform makes the calls: the name of
	// the client's package, such as bedrock. The platform's client options set it.
	Platform string
}

// TracesOptions are the options of the OpenTelemetry spans.
type TracesOptions struct {
	// Enabled turns the spans and the trace-context headers on or off. When it is
	// nil, the environment variable ANTHROPIC_OPEN_TELEMETRY decides.
	Enabled *bool
	// ContentMode is "content" when the spans carry the content of requests
	// and responses, and "metadata_only" when they don't. When it is empty, the
	// environment variable ANTHROPIC_OPEN_TELEMETRY_TRACES_CONTENT_MODE decides.
	ContentMode string
	// MaxContentBytes is how many bytes the content attributes of one span hold
	// at most. When it is nil, the environment variable
	// ANTHROPIC_OPEN_TELEMETRY_TRACES_MAX_CONTENT_BYTES decides.
	MaxContentBytes *int64
}

// Content reports whether the spans carry the content of requests and responses.
func (o Options) Content() bool {
	mode := o.Traces.ContentMode
	if mode == "" {
		mode = os.Getenv("ANTHROPIC_OPEN_TELEMETRY_TRACES_CONTENT_MODE")
	}
	return strings.EqualFold(mode, "content")
}

// MaxContentBytes returns how many bytes the content attributes of one span hold
// at most: the option when it is set, and otherwise the environment variable.
// A maximum below 1, and a variable that is not a whole number, select
// [DefaultMaxContentBytes].
func (o Options) MaxContentBytes() int {
	limit := int64(DefaultMaxContentBytes)
	if o.Traces.MaxContentBytes != nil {
		if *o.Traces.MaxContentBytes >= 1 {
			limit = *o.Traces.MaxContentBytes
		}
	} else if fromEnv, err := strconv.ParseInt(os.Getenv("ANTHROPIC_OPEN_TELEMETRY_TRACES_MAX_CONTENT_BYTES"), 10, 64); err == nil && fromEnv >= 1 {
		limit = fromEnv
	}
	return int(min(limit, math.MaxInt))
}

// Propagates reports whether the requests carry the trace context of their span:
// the option when it is set, and otherwise not when the environment variable is "false".
func (o Options) Propagates() bool {
	if o.Propagation != nil {
		return *o.Propagation
	}
	return !strings.EqualFold(os.Getenv("ANTHROPIC_OPEN_TELEMETRY_PROPAGATION"), "false")
}

// Tracer returns the tracer for these options, or nil when the spans are disabled.
func (o Options) Tracer() trace.Tracer {
	if o.Traces.Enabled != nil && !*o.Traces.Enabled {
		return nil
	}
	if o.Traces.Enabled == nil && strings.EqualFold(os.Getenv("ANTHROPIC_OPEN_TELEMETRY"), "false") {
		return nil
	}
	provider := o.TracerProvider
	if provider == nil {
		provider = otel.GetTracerProvider()
	}
	scope := instrumentationScope
	if o.Platform != "" {
		scope += "." + o.Platform
	}
	return provider.Tracer(scope, trace.WithInstrumentationVersion(internal.PackageVersion))
}

// APICall is one invocation of an SDK method that sends a request to the API,
// once it has a span: every attempt that retries make, and reading the response
// body when the method hands it to the application. The span is a CLIENT span
// named after the method. The SDK emits no span per attempt; the application's
// HTTP instrumentation, when it has one, parents its spans to this one.
//
// requestconfig starts the span before the first attempt ([StartAPICallSpan])
// and sends its trace context with each attempt ([InjectTraceHeaders]). The
// span ends with the call's result: at once when the call failed or the
// response body was decoded ([EndAPICallSpan]), or, when the application reads
// the body, once it has ([EndAPICallSpanWhenBodyEnds]).
type APICall struct {
	Span    trace.Span
	Options Options
	// Start is when the call started. It is kept beside the span, which cannot
	// report its own start time, to measure gen_ai.response.time_to_first_chunk from.
	Start time.Time
	// Response is the response of the latest attempt, or nil when that attempt received none.
	Response *http.Response
	// requests is how many requests the SDK has sent for the call ([RecordSentRequest]).
	// A middleware can send one while the application reads the response body.
	requests atomic.Int32
}

// Content reports whether the span of call carries the content of the request
// and the response: it records, and the call's options are in content mode.
func (call *APICall) Content() bool {
	return call.Span.IsRecording() && call.Options.Content()
}

// HandleRecovered reports recovered, what recover returned while the SDK recorded
// telemetry, to the error handler that the application registered with
// OpenTelemetry: the failure must not reach the API call, and must not be
// invisible either. It does nothing when nothing was recovered.
func HandleRecovered(recovered any) {
	if recovered == nil {
		return
	}
	err, ok := recovered.(error)
	if !ok {
		err = fmt.Errorf("%v", recovered)
	}
	otel.Handle(fmt.Errorf("anthropic: recording telemetry failed: %w", err))
}

// Guard runs record, with which an SDK method records telemetry of its own, and
// hands a panic in it to [HandleRecovered], so that the API call goes on as if
// record had returned.
func Guard(record func()) {
	defer func() { HandleRecovered(recover()) }()

	record()
}

// StartAPICallSpan starts the span of an API call under name, the name of the
// SDK method, as a child of the span in ctx, with the attributes that attrs
// returns, which the method has from its own arguments. It returns nil when
// opts turn the spans off, when name is empty because no SDK method made the
// call, or when the span neither records nor has a trace context that goes on,
// in the headers or, with the propagation off, in the context of the attempts:
// such a span has no effect, so the call runs as if tracing were off. attrs is
// not called in the first two cases, so that a call with the spans off costs
// no more than finding that out.
func StartAPICallSpan(ctx context.Context, opts Options, name string, attrs func() []attribute.KeyValue) *APICall {
	// A panic of the application's tracer, or of attrs, leaves the call without a span.
	defer func() { HandleRecovered(recover()) }()

	tracer := opts.Tracer()
	if tracer == nil || name == "" {
		return nil
	}
	start := time.Now()
	spanCtx, span := tracer.Start(ctx, name, trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs()...))
	if !span.IsRecording() {
		// With the propagation off, a valid context still reaches the spans that
		// the application's HTTP instrumentation parents to this one.
		goesOn := span.SpanContext().IsValid()
		if opts.Propagates() {
			goesOn = len(inject(spanCtx)) > 0
		}
		if !goesOn {
			return nil
		}
	}
	return &APICall{Span: span, Options: opts, Start: start}
}

// ServerAttributes returns the server.address and server.port that the span of
// an API call starts with: those of u, the client's base URL.
func ServerAttributes(u *url.URL) []attribute.KeyValue {
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		port = 443
		if u.Scheme == "http" {
			port = 80
		}
	}
	return []attribute.KeyValue{
		ServerAddressKey.String(u.Hostname()),
		ServerPortKey.Int(port),
	}
}

// RequestAttributes returns the http.request.method and http.route that the
// span of an API call starts with: method is the HTTP method of the request, and
// route the path of the endpoint as the API description gives it, such as
// /v1/files/{file_id}. A call that no SDK method describes has no route.
func RequestAttributes(method, route string) []attribute.KeyValue {
	attrs := []attribute.KeyValue{HTTPRequestMethodKey.String(method)}
	if route != "" {
		attrs = append(attrs, HTTPRouteKey.String(route))
	}
	return attrs
}

// URLAttributes returns the attributes of u, the URL that the SDK sends a request
// to: url.scheme, url.path, url.query and url.full. queryParams are the names of
// the query parameters that the endpoint defines: url.query has only those
// ([definedQuery]), and so has url.full, which leaves out the user information too.
// A query with a parameter whose value is content ([contentQueryParams], and
// endpointContentQueryParams, the ones that are content on this endpoint only)
// makes both of them content: they are among the attributes only when content is true.
func URLAttributes(u *url.URL, queryParams, endpointContentQueryParams []string, content bool) []attribute.KeyValue {
	path, query := u.EscapedPath(), definedQuery(u.RawQuery, queryParams)
	full := u.Scheme + "://" + u.Host + path
	attrs := []attribute.KeyValue{URLSchemeKey.String(u.Scheme), URLPathKey.String(path)}
	if !content && definedQuery(query, slices.Concat(contentQueryParams, endpointContentQueryParams)) != "" {
		return attrs
	}
	if query != "" {
		attrs = append(attrs, URLQueryKey.String(query))
		full += "?" + query
	}
	return append(attrs, URLFullKey.String(full))
}

// RecordSentRequest records that the SDK sends a request for call to u, as the
// HTTP client gets it. requestconfig calls it for each request. The first one
// gives the span its [URLAttributes], which no later one, a retry, changes.
func RecordSentRequest(call *APICall, u *url.URL, queryParams, endpointContentQueryParams []string) {
	// A middleware made the request: one without a URL fails in the HTTP client, not here.
	defer func() { HandleRecovered(recover()) }()

	if call.requests.Add(1) == 1 && call.Span.IsRecording() {
		call.Span.SetAttributes(URLAttributes(u, queryParams, endpointContentQueryParams, call.Options.Content())...)
	}
}

// contentQueryParams are the query parameters that an endpoint defines whose
// value is content: email and path_prefix filter like a search, and
// expected_content_sha256 comes from stored content.
var contentQueryParams = []string{"email", "expected_content_sha256", "path_prefix"}

// definedQuery returns what a span records of rawQuery, a query string without
// its "?": the parameters whose name is in defined, the names that the endpoint
// defines, with their bytes and their order as they are in rawQuery. A name counts
// as defined when, percent-decoded, it is in defined as it is, as created_at[gte]
// is for some endpoints, or when it ends in "[]", as an array is sent, and the
// name before that is, as ids is for ids[]. Any other parameter is left out: it
// is not part of the API description, so its value may be anything.
func definedQuery(rawQuery string, defined []string) string {
	var kept []string
	for _, pair := range strings.Split(rawQuery, "&") {
		name, _, _ := strings.Cut(pair, "=")
		name, err := url.QueryUnescape(name)
		if err != nil {
			continue
		}
		if slices.Contains(defined, name) || slices.Contains(defined, strings.TrimSuffix(name, "[]")) {
			kept = append(kept, pair)
		}
	}
	return strings.Join(kept, "&")
}

// InjectTraceHeaders writes the trace context of the span of call into the
// headers of req, an attempt's request, with the propagator that the application
// registered, in place of any such headers already there. requestconfig calls it
// for each attempt, so a retry carries the context afresh. It leaves req as it
// is when the call's options turn the propagation off.
func InjectTraceHeaders(call *APICall, req *http.Request) {
	defer func() { HandleRecovered(recover()) }()

	if !call.Options.Propagates() {
		return
	}
	otel.GetTextMapPropagator().Inject(trace.ContextWithSpan(req.Context(), call.Span), propagation.HeaderCarrier(req.Header))
}

// EndAPICallSpan ends the span of call, now that the call has a result: it failed
// with err, or call.Response arrived and its body was decoded. The span gets
// http.response.status_code when a response arrived, http.request.resend_count
// when the SDK sent the request again, and its status: ERROR with
// error.type when the call failed, none when the application canceled the
// context, and OK otherwise. The span of a body that the application has yet to
// read ends through [EndAPICallSpanWhenBodyEnds] instead.
//
// It does nothing when the span is not recording.
func EndAPICallSpan(call *APICall, err error) {
	defer func() { HandleRecovered(recover()) }()

	span := call.Span
	if !span.IsRecording() {
		return
	}
	defer span.End()
	setResponseStatusCode(call)
	setResendCount(call)
	if err != nil {
		SetErrorStatus(span, err, "")
	} else {
		span.SetStatus(codes.Ok, "")
	}
}

// EndAPICallSpanWhenBodyEnds makes the span of call, whose response arrived and
// whose body the application reads, itself or through a stream decoder, end when
// that read ends, by replacing call.Response.Body with a [*ResponseBody]: with
// status OK at the end of the body, with the error when a read fails, and with no
// status when the application closes the body before its end or cancels the
// context. A stream of server-sent events needs its last event too for status OK,
// and observer sees each event as the application decodes it. The span gets
// http.response.status_code now, and http.request.resend_count now and again
// when it ends. Without a body to read, the span ends now, with status OK.
//
// It does nothing when the span is not recording.
func EndAPICallSpanWhenBodyEnds(call *APICall, observer StreamObserver) {
	defer func() { HandleRecovered(recover()) }()

	span, res := call.Span, call.Response
	if !span.IsRecording() {
		return
	}
	if res == nil || res.Body == nil || res.Body == http.NoBody {
		EndAPICallSpan(call, nil)
		return
	}
	// The body has the span before anything else can fail, so that the span ends.
	res.Body = &ResponseBody{ReadCloser: res.Body, Call: call, observer: observer}
	setResponseStatusCode(call)
	setResendCount(call)
}

// StreamObserver is what the span of a call follows of the server-sent events
// that the application decodes from the call's response.
type StreamObserver struct {
	// OnEvent sees each event before the event can end the span. Nil sees nothing.
	OnEvent func(ssestream.Event)
	// OnDecoded sees, after OnEvent, each event that the application's stream
	// decodes: event.Data is decoded into *decoded, which holds the value until
	// the stream decodes its next event. Nil sees nothing.
	OnDecoded func(event ssestream.Event, decoded any)
	// LastEvent is the type of the stream's last event, and empty for a stream
	// that has none. A stream that the server ends after its last event was read
	// to the end, and one that it ends before was cut short, which fails the span
	// with api_error. A stream that the application closes before the server ends
	// it was canceled, even after its last event. A stream without a last event
	// is never read to the end: when the server ends it, the span ends with no
	// status.
	LastEvent string
	// OnIncomplete runs before the span ends, when the stream ends before its last
	// event, whichever way it ends. Nil does nothing.
	OnIncomplete func()
}

// setResponseStatusCode sets http.response.status_code on the span of call from
// its response, when one arrived.
func setResponseStatusCode(call *APICall) {
	if call.Response != nil {
		call.Span.SetAttributes(HTTPResponseStatusCodeKey.Int(call.Response.StatusCode))
	}
}

// setResendCount sets http.request.resend_count on the span of call, the number
// of requests that the SDK sent after the first one, when it sent any.
func setResendCount(call *APICall) {
	if resent := int(call.requests.Load()) - 1; resent > 0 {
		call.Span.SetAttributes(HTTPRequestResendCountKey.Int(resent))
	}
}

// ResponseBody is a response body that the application reads. It ends the span
// of the API call that returned it, as [EndAPICallSpanWhenBodyEnds] describes.
type ResponseBody struct {
	io.ReadCloser
	Call *APICall
	// observer counts once the span follows the decoded events ([ResponseBody.TraceEvents]).
	observer StreamObserver
	// sawLast is whether the stream's last event was decoded.
	sawLast atomic.Bool
	// events is whether the span follows decoded server-sent events rather than the body.
	events bool
	once   sync.Once
}

// Read reads from the body. The end of the body ends the span with status OK,
// and a failed read ends it with the error.
func (b *ResponseBody) Read(p []byte) (int, error) {
	n, err := b.ReadCloser.Read(p)
	if err == io.EOF {
		b.end(nil)
	} else if err != nil {
		b.end(err)
	}
	return n, err
}

// Close closes the body. When no read has ended the span, the application is
// closing the body before its end, which ends the span with no status.
func (b *ResponseBody) Close() error {
	// The span ends first, so that a read that the close fails does not fail it.
	b.end(context.Canceled)
	return b.ReadCloser.Close()
}

// Unwrap returns the body that this one wraps. Reading it leaves the span to
// [ResponseBody.TraceEvents].
func (b *ResponseBody) Unwrap() io.ReadCloser { return b.ReadCloser }

// TraceEvents returns decoder, which decodes the body that this one wraps, so
// that the span follows the decoded events instead: it ends when they do.
func (b *ResponseBody) TraceEvents(decoder ssestream.Decoder) ssestream.Decoder {
	b.events = true
	return &eventDecoder{Decoder: decoder, body: b}
}

// end ends the span, the first time that it is called. err is nil at the end of
// the body, [context.Canceled] when the application stopped reading, and
// otherwise the error that failed the read.
func (b *ResponseBody) end(err error) {
	b.once.Do(func() {
		defer func() { HandleRecovered(recover()) }()

		span := b.Call.Span
		defer span.End()
		// A last event shows among decoded events only, not in a body read as bytes.
		hasLast := b.events && b.observer.LastEvent != ""
		complete := hasLast && b.sawLast.Load()
		if hasLast && !complete && b.observer.OnIncomplete != nil {
			Guard(b.observer.OnIncomplete)
		}
		switch {
		case errors.Is(err, context.Canceled):
			// Stopped by the application before the server ended the body: a
			// cancellation, even after the last event.
		case err != nil:
			SetErrorStatus(span, err, "")
		case hasLast && !complete:
			// The application sees such a stream end without an error, so only
			// the span fails. The server cut the stream short, so the failure is the API's.
			SetErrorStatus(span, errIncomplete, "api_error")
		case b.events && !hasLast:
			// Events without a last event are never read to the end: the server ending them sets no status.
		default:
			span.SetStatus(codes.Ok, "")
		}
		// A middleware may have sent a request since the response arrived.
		setResendCount(b.Call)
	})
}

// eventDecoder is the decoder of a [ResponseBody] whose span follows the events.
type eventDecoder struct {
	ssestream.Decoder
	body *ResponseBody
}

// Next decodes the next event. The end of the events, a decoding error, or an
// error event ends the span.
func (d *eventDecoder) Next() bool {
	if !d.Decoder.Next() {
		d.body.end(d.Decoder.Err())
		return false
	}
	event := d.Decoder.Event()
	if d.body.observer.OnEvent != nil {
		Guard(func() { d.body.observer.OnEvent(event) })
	}
	if event.Type == d.body.observer.LastEvent {
		d.body.sawLast.Store(true)
	}
	if event.Type == "error" {
		// The stream fails with this event.
		d.body.end(&EventError{Data: event.Data})
	}
	return true
}

// TraceDecoded passes the observer decoded, a pointer to the value that the
// application's stream decoded the current event into.
func (d *eventDecoder) TraceDecoded(decoded any) {
	if d.body.observer.OnDecoded != nil {
		Guard(func() { d.body.observer.OnDecoded(d.Decoder.Event(), decoded) })
	}
}

// Close closes the decoder. When the events have not ended the span, the
// application is closing the stream early, which ends the span with no status.
func (d *eventDecoder) Close() error {
	d.body.end(context.Canceled)
	return d.Decoder.Close()
}

// SetErrorStatus gives span the status of an operation that failed with err:
// ERROR, with error.type set to typ, or else derived from err. A canceled
// context is the application's doing, not a failure, and sets no status.
func SetErrorStatus(span trace.Span, err error, typ string) {
	if errors.Is(err, context.Canceled) {
		return
	}
	if typ == "" {
		typ = errorType(err)
	}
	span.SetStatus(codes.Error, "")
	span.SetAttributes(ErrorTypeKey.String(typ))
}

// errorType returns the error.type of an operation that failed with err. An
// error response, or the error event of a stream, has the type of its error
// envelope when the API sent one, and an error response has its status code otherwise.
func errorType(err error) string {
	var apiErr *apierror.Error
	if errors.As(err, &apiErr) {
		if errType := apiErr.Type(); errType != "" {
			return string(errType)
		}
		if apiErr.StatusCode >= 400 {
			return strconv.Itoa(apiErr.StatusCode)
		}
	}
	var eventErr *EventError
	if errors.As(err, &eventErr) {
		if errType := gjson.GetBytes(eventErr.Data, "error.type").String(); errType != "" {
			return errType
		}
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return "client_timeout_error"
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		if netErr.Timeout() {
			return "client_timeout_error"
		}
		return "connection_error"
	}
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return "connection_error"
	}
	return "client_error"
}

// EventError is the error of a stream that delivered an error event.
type EventError struct {
	// Data is the data of the event.
	Data []byte
}

// Error implements error.
func (e *EventError) Error() string { return "stream error event" }

// errIncomplete is the error of a stream that ended before its last event.
var errIncomplete = errors.New("stream ended before its last event")

// inject returns the headers that carry the trace context of ctx, from the
// propagator that the application registered.
func inject(ctx context.Context) propagation.HeaderCarrier {
	carrier := propagation.HeaderCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return carrier
}
