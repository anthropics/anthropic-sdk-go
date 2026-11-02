package option

import (
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"go.opentelemetry.io/otel/trace"
)

// The OpenTelemetry options configure the span of each API call. They are
// client options: pass them to the client. Like every RequestOption, one that
// is passed to a single method call applies to that call only.

// WithTracerProvider returns a RequestOption that sets the OpenTelemetry tracer
// provider that creates the span of each API call, in place of the provider
// registered with the OpenTelemetry API. A nil provider selects the registered one.
//
// The client emits spans without this option, through the registered provider.
func WithTracerProvider(provider trace.TracerProvider) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.OpenTelemetry.TracerProvider = provider
		return nil
	})
}

// WithOpenTelemetryPropagation returns a RequestOption that turns the
// trace-context headers of the requests of API calls on or off, and leaves the
// spans as they are. They are on by default: turn them off when the
// instrumentation of your HTTP client already sends the trace context, so that a
// request does not carry it twice.
//
// Without this option, the environment variable
// ANTHROPIC_OPEN_TELEMETRY_PROPAGATION turns them off when its value is "false".
func WithOpenTelemetryPropagation(enabled bool) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.OpenTelemetry.Propagation = &enabled
		return nil
	})
}

// WithOpenTelemetryTraces returns a RequestOption that turns the OpenTelemetry
// spans of API calls, and the trace-context headers of their requests, on or off.
// They are on by default.
//
// Without this option, or [WithoutOpenTelemetry], the environment variable
// ANTHROPIC_OPEN_TELEMETRY turns them off when its value is "false".
func WithOpenTelemetryTraces(enabled bool) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.OpenTelemetry.Traces.Enabled = &enabled
		return nil
	})
}

// WithoutOpenTelemetry returns a RequestOption that turns off the OpenTelemetry
// spans of API calls. It is the short form of WithOpenTelemetryTraces(false).
func WithoutOpenTelemetry() RequestOption {
	return WithOpenTelemetryTraces(false)
}

// OpenTelemetryContentMode is whether the OpenTelemetry spans carry the content
// of requests and responses.
type OpenTelemetryContentMode string

const (
	// OpenTelemetryContentModeMetadataOnly leaves the content out of the spans.
	OpenTelemetryContentModeMetadataOnly OpenTelemetryContentMode = "metadata_only"
	// OpenTelemetryContentModeContent puts the content in the spans.
	OpenTelemetryContentModeContent OpenTelemetryContentMode = "content"
)

// WithOpenTelemetryTracesContentMode returns a RequestOption that sets whether
// the OpenTelemetry spans of API calls carry the content of requests and
// responses. The default is [OpenTelemetryContentModeMetadataOnly].
//
// Without this option, the environment variable
// ANTHROPIC_OPEN_TELEMETRY_TRACES_CONTENT_MODE selects
// [OpenTelemetryContentModeContent] when its value is "content".
func WithOpenTelemetryTracesContentMode(mode OpenTelemetryContentMode) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.OpenTelemetry.Traces.ContentMode = string(mode)
		return nil
	})
}

// WithOpenTelemetryTracesMaxContentBytes returns a RequestOption that sets how
// many bytes of content the OpenTelemetry span of an API call carries at most, in
// [OpenTelemetryContentModeContent]. Content over the maximum is cut from the
// span, which then has the attribute anthropic.content.truncated. The default is
// 750,000 bytes, and so is a maximum below 1.
//
// Without this option, the environment variable
// ANTHROPIC_OPEN_TELEMETRY_TRACES_MAX_CONTENT_BYTES sets the maximum when its
// value is a whole number of at least 1.
func WithOpenTelemetryTracesMaxContentBytes(maxBytes int64) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.OpenTelemetry.Traces.MaxContentBytes = &maxBytes
		return nil
	})
}
