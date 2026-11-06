package requestconfig

import (
	"context"
	"net/http"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/tracing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace/noop"
)

// The request for a later page of a list is a clone: its span has the name, the
// route, the query parameters, the attributes and the options of the first
// request, and not the ID of an object that the first request decoded.
func TestCloneOpenTelemetry(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "https://example.com/v1/things", nil)
	require.NoError(t, err)
	cfg := &RequestConfig{Request: req}
	var res *struct{ ID string }
	traces := false
	cfg.OpenTelemetry.Traces.Enabled = &traces
	require.NoError(t, cfg.Apply(
		WithSpan("anthropic.things.list", "/v1/things", "limit", "name", "page"),
		WithSpanContentQuery("name"),
		WithSpanRecorder("anthropic.things.list", SpanRecorder{RecordResult: func(*tracing.APICall, error) {}}),
		WithSpanAttribute("anthropic.thing.id", "thing_1"),
		WithSpanCreatedID("anthropic.thing.id", &res, func(r *struct{ ID string }) string { return r.ID }),
	))

	clone := cfg.Clone(context.Background()).OpenTelemetry
	assert.Equal(t, "anthropic.things.list", clone.SpanName)
	assert.Equal(t, "/v1/things", clone.SpanRoute)
	assert.Equal(t, []string{"limit", "name", "page"}, clone.SpanQueryParams)
	assert.Equal(t, []string{"name"}, clone.SpanContentQueryParams)
	assert.Equal(t, []attribute.KeyValue{attribute.String("anthropic.thing.id", "thing_1")}, clone.SpanAttributes)
	assert.NotNil(t, clone.spanRecorder.RecordResult)
	assert.Same(t, &traces, clone.Traces.Enabled)
	assert.Nil(t, clone.spanCreatedID)
	assert.NotNil(t, cfg.OpenTelemetry.spanCreatedID)
}

// keptSpan is a span that records, and keeps the attributes that it is given.
type keptSpan struct {
	noop.Span
	attrs []attribute.KeyValue
}

func (*keptSpan) IsRecording() bool { return true }

func (span *keptSpan) SetAttributes(attrs ...attribute.KeyValue) {
	span.attrs = append(span.attrs, attrs...)
}

// The span of a later page of a list reads the body of the response, as that of
// the first page does: the page decodes it into a value that the SDK set. A value
// that an option of the application set keeps the body off the span.
func TestNextPageSpanReadsTheResponseBody(t *testing.T) {
	var first, next, own struct{ ID string }
	cfg, err := NewRequestConfig(context.Background(), http.MethodPost, "v1/messages", nil, &first)
	require.NoError(t, err)
	attributesOf := func(cfg *RequestConfig) []attribute.KeyValue {
		span := &keptSpan{}
		body := &bufferedBody{contents: []byte(`{"id":"msg_1"}`)}
		cfg.setSpanResponseAttributes(&tracing.APICall{Name: "anthropic.messages.create", Span: span, Response: &http.Response{Body: body}})
		return span.attrs
	}
	want := []attribute.KeyValue{attribute.String("gen_ai.response.id", "msg_1")}
	assert.Equal(t, want, attributesOf(cfg))

	page := cfg.Clone(context.Background())
	page.SetResponseBodyInto(&next)
	assert.Equal(t, want, attributesOf(page))

	page.ResponseBodyInto = &own
	assert.Empty(t, attributesOf(page))
}
