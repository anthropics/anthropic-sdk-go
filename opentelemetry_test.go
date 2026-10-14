package anthropic_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/embedded"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

const messageJSON = `{
	"id": "msg_1", "type": "message", "role": "assistant", "model": "claude-served",
	"content": [{"type": "text", "text": "hi"}],
	"stop_reason": "end_turn",
	"usage": {
		"input_tokens": 10, "output_tokens": 5,
		"cache_creation_input_tokens": 3, "cache_read_input_tokens": 2,
		"cache_creation": {"ephemeral_5m_input_tokens": 1, "ephemeral_1h_input_tokens": 2},
		"service_tier": "standard"
	}
}`

const streamBody = "event: message_start\n" +
	`data: {"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","model":"claude-served","content":[],"stop_reason":null,"usage":{"input_tokens":10,"output_tokens":1,"cache_read_input_tokens":2,"cache_creation_input_tokens":3,"cache_creation":{"ephemeral_5m_input_tokens":3,"ephemeral_1h_input_tokens":0},"service_tier":"standard"}}}` + "\n\n" +
	"event: content_block_start\n" +
	`data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}` + "\n\n" +
	"event: ping\n" +
	`data: {"type":"ping"}` + "\n\n" +
	"event: content_block_delta\n" +
	`data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"hi"}}` + "\n\n" +
	"event: content_block_stop\n" +
	`data: {"type":"content_block_stop","index":0}` + "\n\n" +
	"event: message_delta\n" +
	`data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":7,"output_tokens_details":{"thinking_tokens":4}}}` + "\n\n" +
	"event: message_stop\n" +
	`data: {"type":"message_stop"}` + "\n\n"

// respond returns a transport that answers every request with status, headers, and body.
func respond(status int, contentType, body string, header ...string) *closureTransport {
	return &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		h := http.Header{"Content-Type": {contentType}}
		for i := 0; i+1 < len(header); i += 2 {
			h.Set(header[i], header[i+1])
		}
		return &http.Response{StatusCode: status, Header: h, Body: io.NopCloser(strings.NewReader(body)), Request: req}, nil
	}}
}

// spanRecorder is a tracer provider that records its spans in memory.
type spanRecorder struct {
	*sdktrace.TracerProvider
	recorded *tracetest.SpanRecorder
	sampler  *attributeSampler
}

func newSpanRecorder() *spanRecorder {
	recorded := tracetest.NewSpanRecorder()
	sampler := &attributeSampler{started: map[string][]map[string]any{}}
	return &spanRecorder{
		TracerProvider: sdktrace.NewTracerProvider(sdktrace.WithSampler(sampler), sdktrace.WithSpanProcessor(recorded)),
		recorded:       recorded,
		sampler:        sampler,
	}
}

// Spans returns the spans that have ended, in the order they ended.
func (r *spanRecorder) Spans() []sdktrace.ReadOnlySpan { return r.recorded.Ended() }

// Started returns every span that has started, in the order they started.
func (r *spanRecorder) Started() []sdktrace.ReadWriteSpan { return r.recorded.Started() }

func attr(span sdktrace.ReadOnlySpan, key string) (attribute.Value, bool) {
	for _, kv := range span.Attributes() {
		if string(kv.Key) == key {
			return kv.Value, true
		}
	}
	return attribute.Value{}, false
}

func ended(span sdktrace.ReadOnlySpan) bool { return !span.EndTime().IsZero() }

func newTracedClient(t *testing.T, rt http.RoundTripper, opts ...option.RequestOption) (anthropic.Client, *spanRecorder) {
	t.Helper()
	recorder := newSpanRecorder()
	opts = append([]option.RequestOption{
		option.WithAPIKey("key"),
		option.WithBaseURL("http://localhost:4010"),
		option.WithHTTPClient(&http.Client{Transport: rt}),
		option.WithMaxRetries(0),
		option.WithTracerProvider(recorder),
	}, opts...)
	return anthropic.NewClient(opts...), recorder
}

func messageParams() anthropic.MessageNewParams {
	return anthropic.MessageNewParams{
		Model:       "claude-requested",
		MaxTokens:   100,
		Temperature: anthropic.Float(0.5),
		Messages:    []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hello"))},
	}
}

func onlySpan(t *testing.T, recorder *spanRecorder) sdktrace.ReadOnlySpan {
	t.Helper()
	spans := recorder.Started()
	if len(spans) != 1 {
		t.Fatalf("got %d spans, want 1", len(spans))
	}
	if !ended(spans[0]) {
		t.Fatalf("span %q has not ended", spans[0].Name())
	}
	return spans[0]
}

func wantAttrs(t *testing.T, span sdktrace.ReadOnlySpan, want map[string]any) {
	t.Helper()
	for key, value := range want {
		got, ok := attr(span, key)
		if value == nil {
			if ok {
				t.Errorf("%s: attribute %s = %v, want absent", span.Name(), key, got.AsInterface())
			}
			continue
		}
		if !ok {
			t.Errorf("%s: attribute %s is absent, want %v", span.Name(), key, value)
		} else if !reflect.DeepEqual(got.AsInterface(), value) {
			t.Errorf("%s: attribute %s = %#v, want %#v", span.Name(), key, got.AsInterface(), value)
		}
	}
}

func setPropagator(t *testing.T, p propagation.TextMapPropagator) {
	old := otel.GetTextMapPropagator()
	otel.SetTextMapPropagator(p)
	t.Cleanup(func() { otel.SetTextMapPropagator(old) })
}

// wantJSON checks that the attribute key of span holds the JSON value want.
func wantJSON(t *testing.T, span sdktrace.ReadOnlySpan, key, want string) {
	t.Helper()
	got, ok := attr(span, key)
	if !ok {
		t.Errorf("%s: attribute %s is absent", span.Name(), key)
		return
	}
	var gotValue, wantValue any
	if err := json.Unmarshal([]byte(got.AsString()), &gotValue); err != nil {
		t.Fatalf("%s: attribute %s is not JSON: %v", span.Name(), key, err)
	}
	if err := json.Unmarshal([]byte(want), &wantValue); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(gotValue, wantValue) {
		t.Errorf("%s: attribute %s = %s, want %s", span.Name(), key, got.AsString(), want)
	}
}

func TestOpenTelemetryMessageSpan(t *testing.T) {
	setPropagator(t, propagation.TraceContext{})
	var traceparent string
	rt := respond(200, "application/json", messageJSON,
		"request-id", "req_1", "anthropic-organization-id", "org-uuid", "anthropic-workspace-id", "wrkspc_1")
	client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		traceparent = req.Header.Get("traceparent")
		return rt.fn(req)
	}})

	_, err := client.Messages.New(context.Background(), messageParams(),
		option.WithHeader("anthropic-beta", "beta-a, beta-b"),
		option.WithOpenTelemetryConversationID("conv_1"))
	if err != nil {
		t.Fatal(err)
	}

	span := onlySpan(t, recorder)
	if span.Name() != "anthropic.messages.create" || span.SpanKind() != trace.SpanKindClient {
		t.Errorf("span %q of kind %v, want anthropic.messages.create of kind CLIENT", span.Name(), span.SpanKind())
	}
	if span.InstrumentationScope().Name != "com.anthropic.sdk.go" || span.InstrumentationScope().Version != internal.PackageVersion {
		t.Errorf("scope %s %s, want com.anthropic.sdk.go %s", span.InstrumentationScope().Name, span.InstrumentationScope().Version, internal.PackageVersion)
	}
	if span.Status().Code != codes.Ok {
		t.Errorf("status %v, want OK", span.Status().Code)
	}
	if want := span.SpanContext().TraceID().String(); !strings.Contains(traceparent, want) {
		t.Errorf("traceparent %q does not carry the trace %s", traceparent, want)
	}
	wantAttrs(t, span, map[string]any{
		"gen_ai.provider.name":                                     "anthropic",
		"gen_ai.operation.name":                                    "chat",
		"gen_ai.conversation.id":                                   "conv_1",
		"server.address":                                           "localhost",
		"server.port":                                              int64(4010),
		"http.response.status_code":                                int64(200),
		"anthropic.request.id":                                     "req_1",
		"anthropic.organization.id":                                "org-uuid",
		"anthropic.workspace.id":                                   "wrkspc_1",
		"gen_ai.request.model":                                     "claude-requested",
		"gen_ai.request.max_tokens":                                int64(100),
		"gen_ai.request.temperature":                               0.5,
		"gen_ai.request.top_p":                                     nil,
		"gen_ai.request.stream":                                    nil,
		"gen_ai.output.type":                                       nil,
		"anthropic.beta":                                           []string{"beta-a", "beta-b"},
		"gen_ai.response.id":                                       "msg_1",
		"gen_ai.response.model":                                    "claude-served",
		"anthropic.message.stop_reason":                            "end_turn",
		"gen_ai.response.finish_reasons":                           nil,
		"gen_ai.usage.input_tokens":                                int64(15),
		"gen_ai.usage.output_tokens":                               int64(5),
		"gen_ai.usage.cache_write.input_tokens":                    int64(3),
		"gen_ai.usage.cache_read.input_tokens":                     int64(2),
		"anthropic.usage.cache_creation.ephemeral_5m_input_tokens": int64(1),
		"anthropic.usage.cache_creation.ephemeral_1h_input_tokens": int64(2),
		"anthropic.response.service_tier":                          "standard",
		"gen_ai.usage.reasoning.output_tokens":                     nil,
		"gen_ai.response.time_to_first_chunk":                      nil,
		"error.type":                                               nil,
	})
}

func TestOpenTelemetryCompactionTokens(t *testing.T) {
	ctx := context.Background()
	const compaction1 = `{"type":"compaction","input_tokens":100,"cache_read_input_tokens":20,"cache_creation_input_tokens":30,"output_tokens":40}`
	const compaction2 = `{"type":"compaction","input_tokens":1,"cache_read_input_tokens":2,"cache_creation_input_tokens":3,"output_tokens":4}`
	const message = `{"type":"message","input_tokens":10,"cache_read_input_tokens":2,"cache_creation_input_tokens":3,"output_tokens":5}`
	const fallback = `{"type":"fallback_message","model":"claude-first","input_tokens":9,"cache_read_input_tokens":0,"cache_creation_input_tokens":0,"output_tokens":1}`
	// The usage of a Message leaves out its compactions, which usage.iterations
	// has apart: the span has their sum, beside the Message's own tokens.
	cases := map[string]struct {
		usage                                            string
		input, output, compactionInput, compactionOutput any
	}{
		"two compactions and a message": {
			`{"input_tokens":10,"cache_read_input_tokens":2,"cache_creation_input_tokens":3,"output_tokens":5,"iterations":[` + compaction1 + "," + message + "," + compaction2 + `]}`,
			int64(15), int64(5), int64(156), int64(44),
		},
		"a compaction only": {
			`{"input_tokens":0,"output_tokens":0,"iterations":[` + compaction2 + `]}`,
			int64(0), int64(0), int64(6), int64(4),
		},
		"no compaction": {
			`{"input_tokens":10,"output_tokens":5,"iterations":[` + fallback + "," + message + `]}`,
			int64(10), int64(5), nil, nil,
		},
		"no iteration":  {`{"input_tokens":10,"output_tokens":5,"iterations":[]}`, int64(10), int64(5), nil, nil},
		"no iterations": {`{"input_tokens":10,"output_tokens":5,"iterations":null}`, int64(10), int64(5), nil, nil},
		"no such field": {`{"input_tokens":10,"output_tokens":5}`, int64(10), int64(5), nil, nil},
	}
	creates := map[string]func(client anthropic.Client) error{
		"Messages.New": func(client anthropic.Client) error {
			_, err := client.Messages.New(ctx, messageParams())
			return err
		},
		"Beta.Messages.New": func(client anthropic.Client) error {
			_, err := client.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{Model: "claude-requested", MaxTokens: 100})
			return err
		},
	}
	modes := []option.OpenTelemetryContentMode{option.OpenTelemetryContentModeMetadataOnly, option.OpenTelemetryContentModeContent}
	for name, test := range cases {
		for method, create := range creates {
			for _, mode := range modes {
				t.Run(name+" "+method+" "+string(mode), func(t *testing.T) {
					body := `{"id":"msg_1","type":"message","role":"assistant","model":"claude-served","content":[],"stop_reason":"end_turn","usage":` + test.usage + `}`
					client, recorder := newTracedClient(t, respond(200, "application/json", body), option.WithOpenTelemetryTracesContentMode(mode))
					if err := create(client); err != nil {
						t.Fatal(err)
					}
					wantAttrs(t, onlySpan(t, recorder), map[string]any{
						"gen_ai.usage.input_tokens":                test.input,
						"gen_ai.usage.output_tokens":               test.output,
						"anthropic.usage.compaction.input_tokens":  test.compactionInput,
						"anthropic.usage.compaction.output_tokens": test.compactionOutput,
					})
				})
			}
		}
	}
}

// TestOpenTelemetryResponseDecodedElsewhere covers the methods that record
// attributes of their response. The application can have the body decoded into
// a value of its own, and the method then has no response to read.
func TestOpenTelemetryResponseDecodedElsewhere(t *testing.T) {
	ctx := context.Background()
	calls := map[string]struct {
		body string
		call func(client anthropic.Client, opts ...option.RequestOption) (decoded bool, err error)
	}{
		"messages create": {messageJSON, func(client anthropic.Client, opts ...option.RequestOption) (bool, error) {
			res, err := client.Messages.New(ctx, messageParams(), opts...)
			return res != nil, err
		}},
		"beta messages create": {messageJSON, func(client anthropic.Client, opts ...option.RequestOption) (bool, error) {
			res, err := client.Beta.Messages.New(ctx, betaContentParams(), opts...)
			return res != nil, err
		}},
		"count tokens": {`{"input_tokens":7}`, func(client anthropic.Client, opts ...option.RequestOption) (bool, error) {
			res, err := client.Messages.CountTokens(ctx, anthropic.MessageCountTokensParams{Model: "claude-requested"}, opts...)
			return res != nil, err
		}},
		"beta count tokens": {`{"input_tokens":7}`, func(client anthropic.Client, opts ...option.RequestOption) (bool, error) {
			res, err := client.Beta.Messages.CountTokens(ctx, anthropic.BetaMessageCountTokensParams{Model: "claude-requested"}, opts...)
			return res != nil, err
		}},
		"batches create": {`{"id":"msgbatch_new","type":"message_batch"}`, func(client anthropic.Client, opts ...option.RequestOption) (bool, error) {
			res, err := client.Messages.Batches.New(ctx, anthropic.MessageBatchNewParams{}, opts...)
			return res != nil, err
		}},
		"beta batches create": {`{"id":"msgbatch_new","type":"message_batch"}`, func(client anthropic.Client, opts ...option.RequestOption) (bool, error) {
			res, err := client.Beta.Messages.Batches.New(ctx, anthropic.BetaMessageBatchNewParams{}, opts...)
			return res != nil, err
		}},
	}
	for name, test := range calls {
		for _, traced := range []bool{true, false} {
			t.Run(name+map[bool]string{true: " traced", false: " untraced"}[traced], func(t *testing.T) {
				var opts []option.RequestOption
				if !traced {
					opts = append(opts, option.WithoutOpenTelemetry())
				}
				client, recorder := newTracedClient(t, respond(200, "application/json", test.body), opts...)

				// Decoded into a value of the application's.
				var into map[string]any
				decoded, err := test.call(client, option.WithResponseBodyInto(&into))
				if err != nil || decoded || len(into) == 0 {
					t.Errorf("got a response %t, error %v, and %v in the application's value; want it decoded there only", decoded, err, into)
				}
				// Handed to the application unread.
				var raw *http.Response
				decoded, err = test.call(client, option.WithResponseBodyInto(&raw))
				if err != nil || decoded || raw == nil {
					t.Fatalf("got a response %t, error %v, and the raw response %v; want the raw response only", decoded, err, raw)
				}
				if body, err := io.ReadAll(raw.Body); err != nil || string(body) != test.body {
					t.Errorf("read %q, error %v; want the body", body, err)
				}
				raw.Body.Close()

				spans := recorder.Spans()
				if !traced {
					if len(spans) != 0 {
						t.Errorf("got %d spans, want none", len(spans))
					}
					return
				}
				if len(spans) != 2 {
					t.Fatalf("got %d ended spans, want one for each call", len(spans))
				}
				for _, span := range spans {
					if span.Status().Code != codes.Ok {
						t.Errorf("%s: status %v, want OK", span.Name(), span.Status().Code)
					}
					// The method never saw the response, so it records none of it.
					wantAttrs(t, span, map[string]any{
						"gen_ai.response.id":                          nil,
						"anthropic.message.count_tokens.input_tokens": nil,
						"anthropic.message.batch.id":                  nil,
						"gen_ai.provider.name":                        "anthropic",
					})
				}
			})
		}
	}
}

func TestOpenTelemetryStreamResponseDecodedElsewhere(t *testing.T) {
	client, recorder := newTracedClient(t, respond(200, "text/event-stream", streamBody))
	var raw *http.Response
	// The stream has no decoder then, and is not read. The application reads the response.
	_ = client.Messages.NewStreaming(context.Background(), messageParams(), option.WithResponseBodyInto(&raw))
	if raw == nil {
		t.Fatal("the application did not get the raw response")
	}
	if _, err := io.ReadAll(raw.Body); err != nil {
		t.Fatal(err)
	}
	raw.Body.Close()
	span := onlySpan(t, recorder)
	if span.Status().Code != codes.Ok {
		t.Errorf("status %v, want OK", span.Status().Code)
	}
}

func TestOpenTelemetryRequestAttributes(t *testing.T) {
	want := map[string]any{
		"gen_ai.request.model":             "claude-requested",
		"gen_ai.request.max_tokens":        int64(100),
		"gen_ai.request.temperature":       0.5,
		"gen_ai.request.top_k":             int64(5),
		"gen_ai.request.top_p":             0.9,
		"gen_ai.request.reasoning.level":   "high",
		"anthropic.request.service_tier":   "auto",
		"anthropic.thinking.type":          "enabled",
		"anthropic.thinking.budget_tokens": int64(1024),
		"gen_ai.request.stream":            nil,
		"gen_ai.output.type":               nil,
	}
	for name, send := range map[string]func(client anthropic.Client) error{
		"messages": func(client anthropic.Client) error {
			params := messageParams()
			params.TopK = anthropic.Int(5)
			params.TopP = anthropic.Float(0.9)
			params.ServiceTier = anthropic.MessageNewParamsServiceTierAuto
			params.OutputConfig = anthropic.OutputConfigParam{Effort: anthropic.OutputConfigEffortHigh}
			// The type is left to its default, which the request still sends.
			params.Thinking = anthropic.ThinkingConfigParamUnion{OfEnabled: &anthropic.ThinkingConfigEnabledParam{BudgetTokens: 1024}}
			_, err := client.Messages.New(context.Background(), params)
			return err
		},
		"beta messages": func(client anthropic.Client) error {
			_, err := client.Beta.Messages.New(context.Background(), anthropic.BetaMessageNewParams{
				Model:        "claude-requested",
				MaxTokens:    100,
				Temperature:  anthropic.Float(0.5),
				TopK:         anthropic.Int(5),
				TopP:         anthropic.Float(0.9),
				ServiceTier:  anthropic.BetaMessageNewParamsServiceTierAuto,
				OutputConfig: anthropic.BetaOutputConfigParam{Effort: anthropic.BetaOutputConfigEffortHigh},
				Thinking:     anthropic.BetaThinkingConfigParamUnion{OfEnabled: &anthropic.BetaThinkingConfigEnabledParam{BudgetTokens: 1024}},
				Messages:     []anthropic.BetaMessageParam{anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("hello"))},
			})
			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON))
			if err := send(client); err != nil {
				t.Fatal(err)
			}
			wantAttrs(t, onlySpan(t, recorder), want)
		})
	}

	t.Run("adaptive thinking", func(t *testing.T) {
		client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON))
		params := messageParams()
		params.Thinking = anthropic.ThinkingConfigParamUnion{OfAdaptive: &anthropic.ThinkingConfigAdaptiveParam{}}
		if _, err := client.Messages.New(context.Background(), params); err != nil {
			t.Fatal(err)
		}
		wantAttrs(t, onlySpan(t, recorder), map[string]any{
			"anthropic.thinking.type":          "adaptive",
			"anthropic.thinking.budget_tokens": nil,
			"gen_ai.request.top_k":             nil,
			"anthropic.request.service_tier":   nil,
			"gen_ai.request.reasoning.level":   nil,
		})
	})

	// The span reads the type that the request sends, whichever variant has it.
	for name, thinking := range map[string]anthropic.ThinkingConfigParamUnion{
		"between_tools": {OfBetweenTools: &anthropic.ThinkingConfigBetweenToolsParam{}},
		"later_variant": param.Override[anthropic.ThinkingConfigParamUnion](map[string]any{"type": "later_variant"}),
	} {
		t.Run(name+" thinking", func(t *testing.T) {
			client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON))
			params := messageParams()
			params.Thinking = thinking
			if _, err := client.Messages.New(context.Background(), params); err != nil {
				t.Fatal(err)
			}
			wantAttrs(t, onlySpan(t, recorder), map[string]any{
				"anthropic.thinking.type":          name,
				"anthropic.thinking.budget_tokens": nil,
			})
		})
	}
	t.Run("beta between_tools thinking", func(t *testing.T) {
		client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON))
		_, err := client.Beta.Messages.New(context.Background(), anthropic.BetaMessageNewParams{
			Model:     "claude-requested",
			MaxTokens: 100,
			Thinking:  anthropic.BetaThinkingConfigParamUnion{OfBetweenTools: &anthropic.BetaThinkingConfigBetweenToolsParam{}},
			Messages:  []anthropic.BetaMessageParam{anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("hello"))},
		})
		if err != nil {
			t.Fatal(err)
		}
		wantAttrs(t, onlySpan(t, recorder), map[string]any{"anthropic.thinking.type": "between_tools"})
	})
}

func TestOpenTelemetryOutputType(t *testing.T) {
	format := map[string]any{"type": "object", "properties": map[string]any{}}
	for name, send := range map[string]func(client anthropic.Client) error{
		"output_config.format": func(client anthropic.Client) error {
			params := messageParams()
			params.OutputConfig = anthropic.OutputConfigParam{Format: anthropic.JSONOutputFormatParam{Schema: format}}
			_, err := client.Messages.New(context.Background(), params)
			return err
		},
		"deprecated output_format": func(client anthropic.Client) error {
			_, err := client.Beta.Messages.New(context.Background(), anthropic.BetaMessageNewParams{
				Model:        "claude-requested",
				MaxTokens:    100,
				Messages:     []anthropic.BetaMessageParam{anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("hello"))},
				OutputFormat: anthropic.BetaJSONOutputFormatParam{Schema: format},
			})
			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			// Metadata, so it is set whatever the content mode is.
			client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON))
			if err := send(client); err != nil {
				t.Fatal(err)
			}
			wantAttrs(t, onlySpan(t, recorder), map[string]any{"gen_ai.output.type": "json"})
		})
	}
}

func TestOpenTelemetryErrorResponse(t *testing.T) {
	attempts := 0
	rt := &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		attempts++
		status := 529
		if attempts == 2 {
			status = 429
		}
		return respond(status, "application/json",
			`{"type":"error","error":{"type":"rate_limit_error","message":"slow down"}}`,
			"request-id", "req_final", "retry-after-ms", "1").fn(req)
	}}
	client, recorder := newTracedClient(t, rt, option.WithMaxRetries(1))

	_, err := client.Messages.New(context.Background(), messageParams())
	if err == nil {
		t.Fatal("the call did not fail")
	}
	if attempts != 2 {
		t.Fatalf("made %d attempts, want 2", attempts)
	}

	span := onlySpan(t, recorder)
	if span.Status().Code != codes.Error {
		t.Errorf("status %v, want ERROR", span.Status().Code)
	}
	wantAttrs(t, span, map[string]any{
		"error.type":                "rate_limit_error",
		"http.response.status_code": int64(429),
		"anthropic.request.id":      "req_final",
		"gen_ai.request.model":      "claude-requested",
		"url.path":                  "/v1/messages",
		"gen_ai.response.id":        nil,
		"gen_ai.usage.input_tokens": nil,
		// An error response has no stop reason.
		"anthropic.message.stop_reason": nil,
	})
}

func TestOpenTelemetryErrorWithoutEnvelope(t *testing.T) {
	client, recorder := newTracedClient(t, respond(502, "text/html", "<html>bad gateway</html>"))

	if _, err := client.Models.Get(context.Background(), "claude", anthropic.ModelGetParams{}); err == nil {
		t.Fatal("the call did not fail")
	}
	wantAttrs(t, onlySpan(t, recorder), map[string]any{"error.type": "502"})
}

func TestOpenTelemetryConnectionError(t *testing.T) {
	client, recorder := newTracedClient(t, &closureTransport{fn: func(*http.Request) (*http.Response, error) {
		return nil, errors.New("connection refused")
	}})

	if _, err := client.Models.Get(context.Background(), "claude", anthropic.ModelGetParams{}); err == nil {
		t.Fatal("the call did not fail")
	}
	span := onlySpan(t, recorder)
	if span.Status().Code != codes.Error {
		t.Errorf("status %v, want ERROR", span.Status().Code)
	}
	wantAttrs(t, span, map[string]any{
		"error.type":                "connection_error",
		"http.response.status_code": nil,
		"anthropic.request.id":      nil,
		"server.address":            "localhost",
		"server.port":               int64(4010),
		// The request was sent, so the span has where it went.
		"url.path": "/v1/models/claude",
	})
}

func TestOpenTelemetryConnectionErrorAfterAnErrorResponse(t *testing.T) {
	attempts := 0
	rt := &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		attempts++
		if attempts == 1 {
			return respond(429, "application/json",
				`{"type":"error","error":{"type":"rate_limit_error","message":"slow down"}}`,
				"request-id", "req_first", "retry-after-ms", "1").fn(req)
		}
		return nil, errors.New("connection reset")
	}}
	client, recorder := newTracedClient(t, rt, option.WithMaxRetries(1))

	if _, err := client.Models.Get(context.Background(), "claude", anthropic.ModelGetParams{}); err == nil {
		t.Fatal("the call did not fail")
	}
	// The final attempt received no response, so the span has nothing of the first attempt's.
	wantAttrs(t, onlySpan(t, recorder), map[string]any{
		"error.type":                "connection_error",
		"http.response.status_code": nil,
		"anthropic.request.id":      nil,
	})
}

func TestOpenTelemetryCanceledCall(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		cancel()
		return nil, req.Context().Err()
	}})

	if _, err := client.Messages.New(ctx, messageParams()); !errors.Is(err, context.Canceled) {
		t.Fatalf("got error %v, want context.Canceled", err)
	}
	span := onlySpan(t, recorder)
	if span.Status().Code != codes.Unset {
		t.Errorf("status %v, want unset", span.Status().Code)
	}
	wantAttrs(t, span, map[string]any{"error.type": nil})
}

func TestOpenTelemetryParentAndPropagation(t *testing.T) {
	setPropagator(t, propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	var header http.Header
	rt := respond(200, "application/json", `{"id":"claude","type":"model"}`)
	client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		header = req.Header.Clone()
		return rt.fn(req)
	}})

	ctx, parent := recorder.Tracer("app").Start(context.Background(), "handle_order")
	_, err := client.Models.Get(ctx, "claude", anthropic.ModelGetParams{},
		option.WithHeader("traceparent", "00-11111111111111111111111111111111-2222222222222222-01"),
		option.WithHeader("tracestate", "vendor=1"),
		option.WithHeader("baggage", "user=1"))
	if err != nil {
		t.Fatal(err)
	}
	parent.End()

	span := recorder.Started()[1]
	if span.Name() != "anthropic.models.retrieve" {
		t.Fatalf("span %q, want anthropic.models.retrieve", span.Name())
	}
	if span.Parent().SpanID() != parent.SpanContext().SpanID() {
		t.Error("the span is not a child of the current span")
	}
	want := "00-" + span.SpanContext().TraceID().String() + "-" + span.SpanContext().SpanID().String() + "-01"
	if got := header.Values("traceparent"); len(got) != 1 || got[0] != want {
		t.Errorf("traceparent %q, want %q", got, want)
	}
	// The propagators write no tracestate and no baggage, so the application's values stay.
	if got := header.Get("tracestate"); got != "vendor=1" {
		t.Errorf("tracestate %q, want vendor=1", got)
	}
	if got := header.Get("baggage"); got != "user=1" {
		t.Errorf("baggage %q, want user=1", got)
	}
}

func TestOpenTelemetryDisabled(t *testing.T) {
	setPropagator(t, propagation.TraceContext{})
	for name, tc := range map[string]struct {
		env   string
		opts  []option.RequestOption
		spans int
	}{
		"option":                     {opts: []option.RequestOption{option.WithoutOpenTelemetry()}},
		"traces off":                 {opts: []option.RequestOption{option.WithOpenTelemetryTraces(false)}},
		"env":                        {env: "FALSE"},
		"env with another value":     {env: "0", spans: 1},
		"env with a tracer provider": {env: "false", opts: []option.RequestOption{option.WithTracerProvider(nil)}},
		"env with traces on":         {env: "false", opts: []option.RequestOption{option.WithOpenTelemetryTraces(true)}, spans: 1},
		"traces on after off":        {opts: []option.RequestOption{option.WithoutOpenTelemetry(), option.WithOpenTelemetryTraces(true)}, spans: 1},
		"provider after off":         {opts: []option.RequestOption{option.WithoutOpenTelemetry(), option.WithTracerProvider(nil)}},
	} {
		t.Run(name, func(t *testing.T) {
			t.Setenv("ANTHROPIC_OPEN_TELEMETRY", tc.env)
			recorder := newSpanRecorder()
			otel.SetTracerProvider(recorder)
			t.Cleanup(func() { otel.SetTracerProvider(trace.NewNoopTracerProvider()) })

			var traceparent string
			rt := respond(200, "application/json", `{"id":"claude","type":"model"}`)
			client := anthropic.NewClient(append([]option.RequestOption{
				option.WithAPIKey("key"),
				option.WithBaseURL("http://localhost:4010"),
				option.WithHTTPClient(&http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
					traceparent = req.Header.Get("traceparent")
					return rt.fn(req)
				}}}),
			}, tc.opts...)...)

			if _, err := client.Models.Get(context.Background(), "claude", anthropic.ModelGetParams{}); err != nil {
				t.Fatal(err)
			}
			if got := len(recorder.Spans()); got != tc.spans {
				t.Errorf("got %d spans, want %d", got, tc.spans)
			}
			if (traceparent != "") != (tc.spans > 0) {
				t.Errorf("traceparent %q with %d spans", traceparent, tc.spans)
			}
		})
	}
}

// Each OpenTelemetry option sets one setting, and the options of a request apply
// after the client's. So a request that sets one setting keeps the others as the
// client set them, and a setting that both set is the request's.
func TestOpenTelemetryRequestOptionKeepsTheOtherSettingsOfTheClient(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	for name, test := range map[string]struct {
		client, request []option.RequestOption
		// requestProvider gives the request a tracer provider of its own.
		requestProvider bool
		content         bool
		conversation    any
	}{
		"conversation ID": {
			client:       []option.RequestOption{content},
			request:      []option.RequestOption{option.WithOpenTelemetryConversationID("conv_request")},
			content:      true,
			conversation: "conv_request",
		},
		"conversation ID over the client's": {
			client:       []option.RequestOption{content, option.WithOpenTelemetryConversationID("conv_client")},
			request:      []option.RequestOption{option.WithOpenTelemetryConversationID("conv_request")},
			content:      true,
			conversation: "conv_request",
		},
		"content mode": {
			client:       []option.RequestOption{option.WithOpenTelemetryConversationID("conv_client")},
			request:      []option.RequestOption{content},
			content:      true,
			conversation: "conv_client",
		},
		"content mode over the client's": {
			client:       []option.RequestOption{content, option.WithOpenTelemetryConversationID("conv_client")},
			request:      []option.RequestOption{option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeMetadataOnly)},
			conversation: "conv_client",
		},
		"traces on": {
			client:       []option.RequestOption{content, option.WithOpenTelemetryConversationID("conv_client")},
			request:      []option.RequestOption{option.WithOpenTelemetryTraces(true)},
			content:      true,
			conversation: "conv_client",
		},
		"tracer provider": {
			client:          []option.RequestOption{content, option.WithOpenTelemetryConversationID("conv_client")},
			requestProvider: true,
			content:         true,
			conversation:    "conv_client",
		},
		"none": {
			client:       []option.RequestOption{content, option.WithOpenTelemetryConversationID("conv_client")},
			content:      true,
			conversation: "conv_client",
		},
	} {
		t.Run(name, func(t *testing.T) {
			registered := newSpanRecorder()
			otel.SetTracerProvider(registered)
			t.Cleanup(func() { otel.SetTracerProvider(trace.NewNoopTracerProvider()) })

			client, ofClient := newTracedClient(t, respond(200, "application/json", contentResponseJSON), test.client...)
			request, recorder, unused := test.request, ofClient, newSpanRecorder()
			if test.requestProvider {
				recorder, unused = unused, ofClient
				request = append(request, option.WithTracerProvider(recorder))
			}
			if _, err := client.Messages.New(context.Background(), contentParams(), request...); err != nil {
				t.Fatal(err)
			}

			span := onlySpan(t, recorder)
			if got := len(unused.Started()) + len(registered.Started()); got != 0 {
				t.Errorf("%d spans went to a tracer provider that the call does not have", got)
			}
			wantAttrs(t, span, map[string]any{
				"gen_ai.conversation.id": test.conversation,
				"gen_ai.response.id":     "msg_1",
			})
			if test.content {
				wantJSON(t, span, "gen_ai.input.messages", contentInputJSON)
				wantJSON(t, span, "gen_ai.output.messages", contentOutputJSON)
			} else {
				wantAttrs(t, span, map[string]any{"gen_ai.input.messages": nil, "gen_ai.output.messages": nil})
			}
		})
	}
}

// A request that turns the spans off does so for itself only: it leaves the
// client's settings as they are for the next call.
func TestOpenTelemetryTracesOffForOneRequest(t *testing.T) {
	setPropagator(t, propagation.TraceContext{})
	for name, off := range map[string]option.RequestOption{
		"traces off":            option.WithOpenTelemetryTraces(false),
		"without OpenTelemetry": option.WithoutOpenTelemetry(),
	} {
		t.Run(name, func(t *testing.T) {
			registered := newSpanRecorder()
			otel.SetTracerProvider(registered)
			t.Cleanup(func() { otel.SetTracerProvider(trace.NewNoopTracerProvider()) })

			var traceparent string
			rt := respond(200, "application/json", `{"id":"claude","type":"model"}`)
			client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				traceparent = req.Header.Get("traceparent")
				return rt.fn(req)
			}}, option.WithOpenTelemetryTraces(true), option.WithOpenTelemetryConversationID("conv_client"))

			if _, err := client.Models.Get(context.Background(), "claude", anthropic.ModelGetParams{}, off); err != nil {
				t.Fatal(err)
			}
			if got := len(recorder.Started()) + len(registered.Started()); got != 0 || traceparent != "" {
				t.Fatalf("got %d spans and the traceparent %q from a request that turns the spans off", got, traceparent)
			}

			if _, err := client.Models.Get(context.Background(), "claude", anthropic.ModelGetParams{}); err != nil {
				t.Fatal(err)
			}
			span := onlySpan(t, recorder)
			if span.Name() != "anthropic.models.retrieve" {
				t.Errorf("span %q, want anthropic.models.retrieve", span.Name())
			}
			wantAttrs(t, span, map[string]any{"gen_ai.conversation.id": "conv_client"})
			if !strings.Contains(traceparent, span.SpanContext().SpanID().String()) {
				t.Errorf("traceparent %q does not carry the span %s", traceparent, span.SpanContext().SpanID())
			}
			if got := len(registered.Started()); got != 0 {
				t.Errorf("%d spans went to the registered tracer provider, want the client's", got)
			}
		})
	}
}

// The traces switch is one of the settings that a request keeps. A client that
// turns the spans on where the environment turns them off keeps them on for a
// request that sets another setting.
func TestOpenTelemetryRequestOptionKeepsTheTracesSwitchOfTheClient(t *testing.T) {
	setPropagator(t, propagation.TraceContext{})
	t.Setenv("ANTHROPIC_OPEN_TELEMETRY", "false")
	for name, request := range map[string]option.RequestOption{
		"conversation ID": option.WithOpenTelemetryConversationID("conv_request"),
		"content mode":    option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent),
	} {
		t.Run(name, func(t *testing.T) {
			var traceparent string
			rt := respond(200, "application/json", `{"id":"claude","type":"model"}`)
			client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				traceparent = req.Header.Get("traceparent")
				return rt.fn(req)
			}}, option.WithOpenTelemetryTraces(true))

			if _, err := client.Models.Get(context.Background(), "claude", anthropic.ModelGetParams{}, request); err != nil {
				t.Fatal(err)
			}
			span := onlySpan(t, recorder)
			if !strings.Contains(traceparent, span.SpanContext().SpanID().String()) {
				t.Errorf("traceparent %q does not carry the span %s", traceparent, span.SpanContext().SpanID())
			}
		})
	}
}

// encodings counts how many times it is encoded as JSON.
type encodings struct{ count *int }

func (e encodings) MarshalJSON() ([]byte, error) {
	*e.count++
	return []byte("true"), nil
}

// With the spans off, a call does none of the work of a span: the parameters
// are encoded for the request body, and not again for the span's attributes.
func TestOpenTelemetryDisabledEncodesTheParametersOnce(t *testing.T) {
	var count int
	thinking := anthropic.ThinkingConfigEnabledParam{BudgetTokens: 1024}
	thinking.SetExtraFields(map[string]any{"counted": encodings{&count}})
	params := messageParams()
	params.Thinking = anthropic.ThinkingConfigParamUnion{OfEnabled: &thinking}
	if _, err := json.Marshal(params); err != nil {
		t.Fatal(err)
	}
	body := count

	for name, tc := range map[string]struct {
		opts  []option.RequestOption
		spans int
	}{
		// The span reads thinking.type from the encoded configuration.
		"spans on":  {spans: 1},
		"spans off": {opts: []option.RequestOption{option.WithoutOpenTelemetry()}},
	} {
		t.Run(name, func(t *testing.T) {
			client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON), tc.opts...)
			count = 0
			if _, err := client.Messages.New(context.Background(), params); err != nil {
				t.Fatal(err)
			}
			if got := len(recorder.Started()); got != tc.spans {
				t.Fatalf("got %d spans, want %d", got, tc.spans)
			}
			if (count > body) != (tc.spans > 0) || count < body {
				t.Errorf("the thinking configuration was encoded %d times with %d spans, and %d times for the body alone", count, tc.spans, body)
			}
		})
	}
}

// A request that a middleware leaves without a URL fails in the HTTP client, as
// it does without a span: recording the request on the span must not panic.
func TestOpenTelemetryMiddlewareRequestWithoutURL(t *testing.T) {
	var handled []error
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) { handled = append(handled, err) }))
	t.Cleanup(func() { otel.SetErrorHandler(otel.ErrorHandlerFunc(func(error) {})) })
	withoutURL := func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		req.URL = nil
		return next(req)
	}
	client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON), option.WithMiddleware(withoutURL))

	_, err := client.Messages.New(context.Background(), messageParams())
	if err == nil || !strings.Contains(err.Error(), "nil Request.URL") {
		t.Fatalf("the call returned %v, want the error of the HTTP client", err)
	}
	if span := onlySpan(t, recorder); span.Status().Code != codes.Error {
		t.Errorf("the span has status %v, want ERROR", span.Status().Code)
	}
	if len(handled) != 1 {
		t.Errorf("the error handler got %v, want the failure to record the request", handled)
	}
}

func TestOpenTelemetryWithoutTracerProvider(t *testing.T) {
	setPropagator(t, propagation.TraceContext{})
	var header http.Header
	body := &trackingBody{Reader: strings.NewReader("file contents")}
	client := anthropic.NewClient(
		option.WithAPIKey("key"),
		option.WithBaseURL("http://localhost:4010"),
		option.WithHTTPClient(&http.Client{Transport: &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
			header = req.Header.Clone()
			return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/octet-stream"}}, Body: body}, nil
		}}}),
	)

	// The application reads the body of a download itself, which is the body that a span would follow.
	res, err := client.Beta.Files.Download(context.Background(), "file_1", anthropic.BetaFileDownloadParams{})
	if err != nil {
		t.Fatal(err)
	}
	if got := header.Get("traceparent"); got != "" {
		t.Errorf("traceparent %q without a tracer provider", got)
	}
	if res.Body != io.ReadCloser(body) {
		t.Error("the response body was wrapped without a tracer provider")
	}
}

type trackingBody struct{ io.Reader }

func (*trackingBody) Close() error { return nil }

func TestOpenTelemetryRawBody(t *testing.T) {
	client, recorder := newTracedClient(t, respond(200, "application/octet-stream", "file contents"))

	res, err := client.Beta.Files.Download(context.Background(), "file_1", anthropic.BetaFileDownloadParams{})
	if err != nil {
		t.Fatal(err)
	}
	if len(recorder.Spans()) != 0 {
		t.Fatal("the span ended before the body was read")
	}
	if _, err := io.ReadAll(res.Body); err != nil {
		t.Fatal(err)
	}
	span := onlySpan(t, recorder)
	if span.Name() != "anthropic.files.download" || span.Status().Code != codes.Ok {
		t.Errorf("span %q with status %v, want anthropic.files.download with OK", span.Name(), span.Status().Code)
	}
	_ = res.Body.Close()
}

func TestOpenTelemetryCountTokens(t *testing.T) {
	client, recorder := newTracedClient(t, respond(200, "application/json", `{"input_tokens":42}`))

	_, err := client.Messages.CountTokens(context.Background(), anthropic.MessageCountTokensParams{
		Model:    "claude-requested",
		Messages: []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hi"))},
	})
	if err != nil {
		t.Fatal(err)
	}
	span := onlySpan(t, recorder)
	if span.Name() != "anthropic.messages.count_tokens" {
		t.Errorf("span %q", span.Name())
	}
	wantAttrs(t, span, map[string]any{
		"gen_ai.request.model":                        "claude-requested",
		"anthropic.message.count_tokens.input_tokens": int64(42),
		"gen_ai.operation.name":                       nil,
		"gen_ai.usage.input_tokens":                   nil,
	})
}

func TestOpenTelemetryBatches(t *testing.T) {
	results := `{"custom_id":"a","result":{"type":"canceled"}}` + "\n" + `{"custom_id":"b","result":{"type":"canceled"}}` + "\n"
	rt := &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		switch {
		case strings.HasSuffix(req.URL.Path, "/results"):
			return respond(200, "application/binary", results).fn(req)
		case req.Method == http.MethodPost:
			return respond(200, "application/json", `{"id":"msgbatch_new","type":"message_batch"}`).fn(req)
		default:
			return respond(200, "application/json", `{"id":"msgbatch_1","type":"message_batch","results_url":"http://localhost:4010/v1/messages/batches/msgbatch_1/results"}`).fn(req)
		}
	}}
	client, recorder := newTracedClient(t, rt)
	ctx := context.Background()

	if _, err := client.Messages.Batches.New(ctx, anthropic.MessageBatchNewParams{}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Messages.Batches.Get(ctx, "msgbatch_1", anthropic.MessageBatchGetParams{}); err != nil {
		t.Fatal(err)
	}
	stream := client.Messages.Batches.ResultsStreaming(ctx, "msgbatch_1", anthropic.MessageBatchResultsParams{})
	n := 0
	for stream.Next() {
		n++
	}
	if err := stream.Err(); err != nil || n != 2 {
		t.Fatalf("read %d results, error %v", n, err)
	}

	// The batch is the one the method was called with, whatever the response names.
	if _, err := client.Beta.Messages.Batches.Cancel(ctx, "msgbatch_1", anthropic.BetaMessageBatchCancelParams{}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Beta.Messages.Batches.Delete(ctx, "msgbatch_1", anthropic.BetaMessageBatchDeleteParams{}); err != nil {
		t.Fatal(err)
	}

	spans := recorder.Spans()
	if len(spans) != 5 {
		t.Fatalf("got %d spans, want 5", len(spans))
	}
	for i, want := range []struct{ name, id string }{
		{"anthropic.messages.batches.create", "msgbatch_new"},
		{"anthropic.messages.batches.retrieve", "msgbatch_1"},
		{"anthropic.messages.batches.results", "msgbatch_1"},
		{"anthropic.messages.batches.cancel", "msgbatch_1"},
		{"anthropic.messages.batches.delete", "msgbatch_1"},
	} {
		if spans[i].Name() != want.name || spans[i].Status().Code != codes.Ok {
			t.Errorf("span %d: %q with status %v, want %q with OK", i, spans[i].Name(), spans[i].Status().Code, want.name)
		}
		wantAttrs(t, spans[i], map[string]any{"anthropic.message.batch.id": want.id})
	}
}

func TestOpenTelemetryPages(t *testing.T) {
	pages := []string{
		`{"data":[{"id":"a","type":"model"}],"has_more":true,"first_id":"a","last_id":"a"}`,
		`{"data":[{"id":"b","type":"model"}],"has_more":false,"first_id":"b","last_id":"b"}`,
	}
	calls := 0
	client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		calls++
		return respond(200, "application/json", pages[calls-1]).fn(req)
	}})

	iter := client.Models.ListAutoPaging(context.Background(), anthropic.ModelListParams{})
	for iter.Next() {
	}
	if err := iter.Err(); err != nil {
		t.Fatal(err)
	}
	spans := recorder.Spans()
	if len(spans) != 2 || spans[0].Name() != "anthropic.models.list" || spans[1].Name() != "anthropic.models.list" {
		t.Fatalf("got %d spans, want two anthropic.models.list spans", len(spans))
	}
}

func TestOpenTelemetryConcurrentCallsShareOptions(t *testing.T) {
	setPropagator(t, propagation.TraceContext{})
	var mu sync.Mutex
	traceparents := map[string]string{}
	rt := &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		mu.Lock()
		traceparents[req.URL.Path] = req.Header.Get("traceparent")
		mu.Unlock()
		return respond(200, "application/json", `{"id":"claude","type":"model"}`).fn(req)
	}}
	client, recorder := newTracedClient(t, rt)

	// One option slice, shared by both calls.
	shared := []option.RequestOption{option.WithHeader("x-shared", "1")}
	var wg sync.WaitGroup
	for _, id := range []string{"model-a", "model-b"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := client.Models.Get(context.Background(), id, anthropic.ModelGetParams{}, shared...); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()

	spans := recorder.Spans()
	if len(spans) != 2 {
		t.Fatalf("got %d spans, want 2", len(spans))
	}
	seen := map[string]bool{}
	for _, span := range spans {
		seen[span.SpanContext().SpanID().String()] = true
	}
	for path, traceparent := range traceparents {
		parts := strings.Split(traceparent, "-")
		if len(parts) != 4 || !seen[parts[2]] {
			t.Errorf("%s: traceparent %q names no span of its own", path, traceparent)
		}
		delete(seen, parts[2])
	}
	if len(seen) != 0 {
		t.Errorf("the two calls did not send two traceparents: %v", traceparents)
	}
}

// instrumentedTransport stands in for an HTTP client instrumentation: it reports
// each attempt as a span, a child of the span that is current in the request's context.
type instrumentedTransport struct{ next *closureTransport }

func (t *instrumentedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	parent := trace.SpanFromContext(req.Context())
	_, span := parent.TracerProvider().Tracer("http").Start(req.Context(), req.Method, trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()
	return t.next.fn(req)
}

func TestOpenTelemetryAttemptSpans(t *testing.T) {
	var retryCounts []string
	rt := &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		retryCounts = append(retryCounts, req.Header.Get("X-Stainless-Retry-Count"))
		if len(retryCounts) == 1 {
			return respond(529, "application/json", `{"type":"error","error":{"type":"overloaded_error","message":"busy"}}`, "retry-after-ms", "1").fn(req)
		}
		return respond(200, "application/json", messageJSON, "request-id", "req_2").fn(req)
	}}
	client, recorder := newTracedClient(t, &instrumentedTransport{next: rt}, option.WithMaxRetries(2))

	if _, err := client.Messages.New(context.Background(), messageParams()); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(retryCounts, []string{"0", "1"}) {
		t.Fatalf("retry counts %q, want 0 and 1", retryCounts)
	}

	var call sdktrace.ReadOnlySpan
	var attempts []sdktrace.ReadOnlySpan
	for _, s := range recorder.Spans() {
		if s.Name() == "anthropic.messages.create" {
			call = s
		} else {
			attempts = append(attempts, s)
		}
	}
	if call == nil || len(attempts) != 2 {
		t.Fatalf("got %d spans, want one call span and two attempt spans", len(recorder.Spans()))
	}
	for _, attempt := range attempts {
		if attempt.Parent().SpanID() != call.SpanContext().SpanID() {
			t.Errorf("attempt span %q is not a child of the call span", attempt.Name())
		}
	}
	if call.Status().Code != codes.Ok {
		t.Errorf("status %v, want OK", call.Status().Code)
	}
	wantAttrs(t, call, map[string]any{
		"http.response.status_code": int64(200),
		"anthropic.request.id":      "req_2",
		"gen_ai.response.id":        "msg_1",
		"error.type":                nil,
	})
}

// batchResultsClient returns a traced client whose batch results download is body.
func batchResultsClient(t *testing.T, body io.Reader) (anthropic.Client, *spanRecorder) {
	return newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/binary"}}, Body: io.NopCloser(body), Request: req}, nil
	}})
}

func TestOpenTelemetryBatchResultsSpanCoversTheDownload(t *testing.T) {
	results := `{"custom_id":"a","result":{"type":"canceled"}}` + "\n" + `{"custom_id":"b","result":{"type":"canceled"}}` + "\n"
	client, recorder := batchResultsClient(t, strings.NewReader(results))

	stream := client.Messages.Batches.ResultsStreaming(context.Background(), "msgbatch_1", anthropic.MessageBatchResultsParams{})
	if !stream.Next() {
		t.Fatal(stream.Err())
	}
	if len(recorder.Spans()) != 0 {
		t.Fatal("the span ended before the results were read")
	}
	for stream.Next() {
	}
	if err := stream.Err(); err != nil {
		t.Fatal(err)
	}
	span := onlySpan(t, recorder)
	if span.Status().Code != codes.Ok {
		t.Errorf("status %v, want OK", span.Status().Code)
	}
	wantAttrs(t, span, map[string]any{"anthropic.message.batch.id": "msgbatch_1", "http.response.status_code": int64(200)})
	_ = stream.Close()
}

// failingReader delivers data and then fails as a dropped connection does.
type failingReader struct{ data string }

func (r *failingReader) Read(p []byte) (int, error) {
	if r.data == "" {
		return 0, io.ErrUnexpectedEOF
	}
	n := copy(p, r.data)
	r.data = r.data[n:]
	return n, nil
}

func TestOpenTelemetryBatchResultsFailureWhileReading(t *testing.T) {
	client, recorder := batchResultsClient(t, &failingReader{data: `{"custom_id":"a","result":{"type":"canceled"}}` + "\n" + `{"custom_id":`})

	stream := client.Messages.Batches.ResultsStreaming(context.Background(), "msgbatch_1", anthropic.MessageBatchResultsParams{})
	for stream.Next() {
	}
	if stream.Err() == nil {
		t.Fatal("the download did not fail")
	}
	span := onlySpan(t, recorder)
	if span.Status().Code != codes.Error {
		t.Errorf("status %v, want ERROR", span.Status().Code)
	}
	wantAttrs(t, span, map[string]any{"error.type": "connection_error"})
}

func TestOpenTelemetryBatchResultsClosedEarly(t *testing.T) {
	results := `{"custom_id":"a","result":{"type":"canceled"}}` + "\n" + `{"custom_id":"b","result":{"type":"canceled"}}` + "\n"
	for _, reads := range []int{0, 1} {
		client, recorder := batchResultsClient(t, strings.NewReader(results))
		stream := client.Messages.Batches.ResultsStreaming(context.Background(), "msgbatch_1", anthropic.MessageBatchResultsParams{})
		for range reads {
			if !stream.Next() {
				t.Fatal(stream.Err())
			}
		}
		if err := stream.Close(); err != nil {
			t.Fatal(err)
		}
		_ = stream.Close()
		span := onlySpan(t, recorder)
		if span.Status().Code != codes.Unset {
			t.Errorf("closed after %d reads: status %v, want unset", reads, span.Status().Code)
		}
		wantAttrs(t, span, map[string]any{"error.type": nil})
	}
}

func TestOpenTelemetryRawBodyClosedUnread(t *testing.T) {
	client, recorder := newTracedClient(t, respond(200, "application/json", `{"id":"claude","type":"model"}`))

	var raw *http.Response
	if _, err := client.Models.Get(context.Background(), "claude", anthropic.ModelGetParams{}, option.WithResponseBodyInto(&raw)); err != nil {
		t.Fatal(err)
	}
	if len(recorder.Spans()) != 0 {
		t.Fatal("the span ended before the application was done with the body")
	}
	// The body is handed over unread, so closing it is stopping before the end.
	if err := raw.Body.Close(); err != nil {
		t.Fatal(err)
	}
	span := onlySpan(t, recorder)
	if span.Status().Code != codes.Unset {
		t.Errorf("status %v, want unset", span.Status().Code)
	}
	wantAttrs(t, span, map[string]any{"http.response.status_code": int64(200)})
}

// A response without a body, as the HTTP client hands it over for a 204, has
// nothing for the application to read: the call is over when it returns.
func TestOpenTelemetryRawResponseWithoutABody(t *testing.T) {
	client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 204, Header: http.Header{}, Body: http.NoBody, Request: req}, nil
	}})

	var raw *http.Response
	_, err := client.Beta.Environments.Work.Stop(context.Background(), "work_1", anthropic.BetaEnvironmentWorkStopParams{EnvironmentID: "env_1"},
		option.WithResponseBodyInto(&raw))
	if err != nil {
		t.Fatal(err)
	}
	span := onlySpan(t, recorder)
	if span.Name() != "anthropic.environments.work.stop" || span.Status().Code != codes.Ok {
		t.Errorf("span %q with status %v, want anthropic.environments.work.stop with OK", span.Name(), span.Status().Code)
	}
	wantAttrs(t, span, map[string]any{"http.response.status_code": int64(204)})
}

// attributeSampler samples as a tracer provider does by default, and keeps the
// attributes that each span started with, by span name, in the order the spans started.
type attributeSampler struct {
	mu      sync.Mutex
	started map[string][]map[string]any
}

func (s *attributeSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	attrs := map[string]any{}
	for _, kv := range p.Attributes {
		attrs[string(kv.Key)] = kv.Value.AsInterface()
	}
	s.mu.Lock()
	s.started[p.Name] = append(s.started[p.Name], attrs)
	s.mu.Unlock()
	return sdktrace.ParentBased(sdktrace.AlwaysSample()).ShouldSample(p)
}

// last returns the attributes that the latest span named name started with.
func (s *attributeSampler) last(name string) map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.started[name][len(s.started[name])-1]
}

func (s *attributeSampler) Description() string { return "attributeSampler" }

func TestOpenTelemetrySamplerSeesTheAttributesThatASpanStartsWith(t *testing.T) {
	pages := []string{
		`{"data":[{"id":"a","type":"model"}],"has_more":true,"first_id":"a","last_id":"a"}`,
		`{"data":[{"id":"b","type":"model"}],"has_more":false,"first_id":"b","last_id":"b"}`,
	}
	listed := 0
	client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		body := messageJSON
		if req.URL.Path == "/v1/models" {
			body = pages[listed]
			listed++
		}
		return respond(200, "application/json", body).fn(req)
	}}, option.WithOpenTelemetryConversationID("conv_1"))

	ctx := context.Background()
	if _, err := client.Messages.New(ctx, messageParams()); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Messages.CountTokens(ctx, anthropic.MessageCountTokensParams{
		Model:    "claude-requested",
		Messages: []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hi"))},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Messages.Batches.Get(ctx, "msgbatch_1", anthropic.MessageBatchGetParams{}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Models.Get(ctx, "claude-requested", anthropic.ModelGetParams{}); err != nil {
		t.Fatal(err)
	}
	iter := client.Models.ListAutoPaging(ctx, anthropic.ModelListParams{})
	for iter.Next() {
	}
	if err := iter.Err(); err != nil {
		t.Fatal(err)
	}

	// Every span starts with these, with the method and the route of its endpoint,
	// and with what its method knows of the request.
	with := func(method, route string, attrs map[string]any) map[string]any {
		attrs["gen_ai.provider.name"] = "anthropic"
		attrs["gen_ai.conversation.id"] = "conv_1"
		attrs["server.address"] = "localhost"
		attrs["server.port"] = int64(4010)
		attrs["http.request.method"] = method
		attrs["http.route"] = route
		return attrs
	}
	want := map[string][]map[string]any{
		"anthropic.messages.create": {with("POST", "/v1/messages", map[string]any{
			"gen_ai.operation.name":      "chat",
			"gen_ai.request.model":       "claude-requested",
			"gen_ai.request.max_tokens":  int64(100),
			"gen_ai.request.temperature": 0.5,
		})},
		"anthropic.messages.count_tokens":     {with("POST", "/v1/messages/count_tokens", map[string]any{"gen_ai.request.model": "claude-requested"})},
		"anthropic.messages.batches.retrieve": {with("GET", "/v1/messages/batches/{message_batch_id}", map[string]any{"anthropic.message.batch.id": "msgbatch_1"})},
		"anthropic.models.retrieve":           {with("GET", "/v1/models/{model_id}", map[string]any{"anthropic.model.id": "claude-requested"})},
		"anthropic.models.list":               {with("GET", "/v1/models", map[string]any{}), with("GET", "/v1/models", map[string]any{})},
	}
	if got := recorder.sampler.started; !reflect.DeepEqual(got, want) {
		t.Errorf("the spans started with %v, want %v", got, want)
	}
}

func TestOpenTelemetrySamplerSeesTheRequestParametersOfMessages(t *testing.T) {
	sampler := &attributeSampler{started: map[string][]map[string]any{}}
	newClient := func(contentType, body string) anthropic.Client {
		return anthropic.NewClient(
			option.WithAPIKey("key"),
			option.WithBaseURL("http://localhost:4010"),
			option.WithHTTPClient(&http.Client{Transport: respond(200, contentType, body)}),
			option.WithMaxRetries(0),
			option.WithTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSampler(sampler))),
		)
	}
	ctx := context.Background()
	opts := []option.RequestOption{
		option.WithHeader("anthropic-beta", "beta-a, beta-b"),
		option.WithOpenTelemetryConversationID("conv_1"),
	}
	params := messageParams()
	params.TopK = anthropic.Int(5)
	params.TopP = anthropic.Float(0.9)
	params.ServiceTier = anthropic.MessageNewParamsServiceTierAuto
	params.Thinking = anthropic.ThinkingConfigParamUnion{OfEnabled: &anthropic.ThinkingConfigEnabledParam{BudgetTokens: 2048}}
	params.OutputConfig = anthropic.OutputConfigParam{
		Effort: anthropic.OutputConfigEffortHigh,
		Format: anthropic.JSONOutputFormatParam{Schema: map[string]any{"type": "object"}},
	}
	want := map[string]any{
		"gen_ai.operation.name":            "chat",
		"gen_ai.provider.name":             "anthropic",
		"gen_ai.conversation.id":           "conv_1",
		"server.address":                   "localhost",
		"server.port":                      int64(4010),
		"http.request.method":              "POST",
		"http.route":                       "/v1/messages",
		"gen_ai.request.model":             "claude-requested",
		"gen_ai.request.max_tokens":        int64(100),
		"gen_ai.request.temperature":       0.5,
		"gen_ai.request.top_k":             int64(5),
		"gen_ai.request.top_p":             0.9,
		"gen_ai.request.reasoning.level":   "high",
		"gen_ai.output.type":               "json",
		"anthropic.request.service_tier":   "auto",
		"anthropic.thinking.type":          "enabled",
		"anthropic.thinking.budget_tokens": int64(2048),
		"anthropic.beta":                   []string{"beta-a", "beta-b"},
	}

	client := newClient("application/json", messageJSON)
	if _, err := client.Messages.New(ctx, params, opts...); err != nil {
		t.Fatal(err)
	}
	if got := sampler.last("anthropic.messages.create"); !reflect.DeepEqual(got, want) {
		t.Errorf("Messages.New: the span started with %v, want %v", got, want)
	}

	streaming := newClient("text/event-stream", streamBody)
	stream := streaming.Messages.NewStreaming(ctx, params, opts...)
	for stream.Next() {
	}
	if err := stream.Err(); err != nil {
		t.Fatal(err)
	}
	want["gen_ai.request.stream"] = true
	if got := sampler.last("anthropic.messages.create"); !reflect.DeepEqual(got, want) {
		t.Errorf("Messages.NewStreaming: the span started with %v, want %v", got, want)
	}

	betaParams := anthropic.BetaMessageNewParams{
		Model:     "claude-requested",
		MaxTokens: 100,
		Messages:  []anthropic.BetaMessageParam{anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("hello"))},
		Betas:     []anthropic.AnthropicBeta{"beta-c"},
	}
	if _, err := client.Beta.Messages.New(ctx, betaParams, opts...); err != nil {
		t.Fatal(err)
	}
	// The header and the parameters both name betas, and the span starts with all of them,
	// in the order the request sends them: the parameters' betas come first.
	wantBeta := map[string]any{
		"gen_ai.operation.name":     "chat",
		"gen_ai.provider.name":      "anthropic",
		"gen_ai.conversation.id":    "conv_1",
		"server.address":            "localhost",
		"server.port":               int64(4010),
		"http.request.method":       "POST",
		"http.route":                "/v1/messages",
		"gen_ai.request.model":      "claude-requested",
		"gen_ai.request.max_tokens": int64(100),
		"anthropic.beta":            []string{"beta-c", "beta-a", "beta-b"},
	}
	if got := sampler.last("anthropic.messages.create"); !reflect.DeepEqual(got, wantBeta) {
		t.Errorf("Beta.Messages.New: the span started with %v, want %v", got, wantBeta)
	}
}

func TestOpenTelemetrySamplerSeesTheBatchRequestCount(t *testing.T) {
	ctx := context.Background()
	limit := anthropic.Int(1)
	const requestCountKey = "anthropic.message.batch.request_count"
	// Only the creation of a batch has the number of its requests.
	cases := []struct {
		name string
		want any
		call func(c anthropic.Client) error
	}{
		{"anthropic.messages.create", nil, func(c anthropic.Client) error {
			_, err := c.Messages.New(ctx, messageParams())
			return err
		}},
		{"anthropic.messages.batches.create", int64(2), func(c anthropic.Client) error {
			_, err := c.Messages.Batches.New(ctx, anthropic.MessageBatchNewParams{Requests: []anthropic.MessageBatchNewParamsRequest{{CustomID: "a"}, {CustomID: "b"}}})
			return err
		}},
		{"anthropic.messages.batches.create", int64(1), func(c anthropic.Client) error {
			_, err := c.Beta.Messages.Batches.New(ctx, anthropic.BetaMessageBatchNewParams{Requests: []anthropic.BetaMessageBatchNewParamsRequest{{CustomID: "a"}}})
			return err
		}},
		{"anthropic.messages.batches.retrieve", nil, func(c anthropic.Client) error {
			_, err := c.Messages.Batches.Get(ctx, "msgbatch_1", anthropic.MessageBatchGetParams{})
			return err
		}},
		{"anthropic.messages.batches.list", nil, func(c anthropic.Client) error {
			return drain(c.Messages.Batches.ListAutoPaging(ctx, anthropic.MessageBatchListParams{Limit: limit}))
		}},
	}
	for _, test := range cases {
		client, recorder := newTracedClient(t, resourceTransport(messageJSON))
		if err := test.call(client); err != nil {
			t.Fatalf("%s: %v", test.name, err)
		}
		for _, started := range recorder.sampler.started[test.name] {
			if started[requestCountKey] != test.want {
				t.Errorf("%s: the span started with %s = %v, want %v", test.name, requestCountKey, started[requestCountKey], test.want)
			}
		}
		if len(recorder.sampler.started[test.name]) == 0 {
			t.Errorf("%s: no such span started", test.name)
		}
	}
}

// startOnlyServer is a tracer provider that fails the test when a span is given
// server.address, server.port, http.request.method or http.route once it has started.
type startOnlyServer struct {
	trace.TracerProvider
	t *testing.T
}

func (p startOnlyServer) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return startOnlyServerTracer{p.TracerProvider.Tracer(name, opts...), p.t}
}

type startOnlyServerTracer struct {
	trace.Tracer
	t *testing.T
}

func (tr startOnlyServerTracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	ctx, span := tr.Tracer.Start(ctx, name, opts...)
	span = startOnlyServerSpan{span, tr.t}
	return trace.ContextWithSpan(ctx, span), span
}

type startOnlyServerSpan struct {
	trace.Span
	t *testing.T
}

func (s startOnlyServerSpan) SetAttributes(attrs ...attribute.KeyValue) {
	for _, kv := range attrs {
		if strings.HasPrefix(string(kv.Key), "server.") || kv.Key == "http.request.method" || kv.Key == "http.route" {
			s.t.Errorf("%s was set after the span started", kv.Key)
		}
	}
	s.Span.SetAttributes(attrs...)
}

func TestOpenTelemetryServerAttributesAreSetOnlyWhenTheSpanStarts(t *testing.T) {
	attempts := 0
	recorder := newSpanRecorder()
	client, _ := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		attempts++
		if attempts == 1 {
			return respond(429, "application/json",
				`{"type":"error","error":{"type":"rate_limit_error","message":"slow down"}}`, "retry-after-ms", "1").fn(req)
		}
		return respond(200, "application/json", messageJSON).fn(req)
	}}, option.WithMaxRetries(1), option.WithTracerProvider(startOnlyServer{recorder, t}))

	if _, err := client.Messages.New(context.Background(), messageParams()); err != nil {
		t.Fatal(err)
	}
	if attempts != 2 {
		t.Fatalf("got %d attempts, want 2", attempts)
	}
	wantAttrs(t, onlySpan(t, recorder), map[string]any{
		"server.address":      "localhost",
		"server.port":         int64(4010),
		"http.request.method": "POST",
		"http.route":          "/v1/messages",
	})
}

func TestOpenTelemetryAttemptTimeoutAfterTheResponseHasNoStatusCode(t *testing.T) {
	// The response arrives as the time of the attempt runs out, so the attempt
	// fails with the timeout and has no response.
	rt := &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		<-req.Context().Done()
		return respond(200, "application/json", messageJSON).fn(req)
	}}
	client, recorder := newTracedClient(t, rt)
	_, err := client.Models.Get(context.Background(), "claude-requested", anthropic.ModelGetParams{}, option.WithRequestTimeout(10*time.Millisecond))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want the timeout", err)
	}
	wantAttrs(t, onlySpan(t, recorder), map[string]any{"error.type": "client_timeout_error", "http.response.status_code": nil})
}

// brokenProvider is a tracer provider of the application's that panics in the
// method that panicsIn names, and counts the spans that it starts and ends.
type brokenProvider struct {
	embedded.TracerProvider
	panicsIn       string
	started, ended atomic.Int32
}

func (p *brokenProvider) fail(method string) {
	if p.panicsIn == method {
		panic(method + " is broken")
	}
}

func (p *brokenProvider) Tracer(string, ...trace.TracerOption) trace.Tracer {
	p.fail("TracerProvider.Tracer")
	return &brokenTracer{provider: p}
}

type brokenTracer struct {
	embedded.Tracer
	provider *brokenProvider
}

func (t *brokenTracer) Start(ctx context.Context, _ string, _ ...trace.SpanStartOption) (context.Context, trace.Span) {
	t.provider.fail("Tracer.Start")
	span := &brokenSpan{provider: t.provider, spanContext: trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{1},
		SpanID:     trace.SpanID{byte(t.provider.started.Add(1))},
		TraceFlags: trace.FlagsSampled,
	})}
	return trace.ContextWithSpan(ctx, span), span
}

type brokenSpan struct {
	noop.Span
	provider    *brokenProvider
	spanContext trace.SpanContext
}

func (s *brokenSpan) IsRecording() bool                   { return true }
func (s *brokenSpan) SpanContext() trace.SpanContext      { return s.spanContext }
func (s *brokenSpan) SetAttributes(...attribute.KeyValue) { s.provider.fail("Span.SetAttributes") }
func (s *brokenSpan) SetStatus(codes.Code, string)        { s.provider.fail("Span.SetStatus") }
func (s *brokenSpan) End(...trace.SpanEndOption) {
	s.provider.ended.Add(1)
	s.provider.fail("Span.End")
}

// brokenPropagator is a propagator of the application's that panics when it injects.
type brokenPropagator struct{ propagation.TraceContext }

func (brokenPropagator) Inject(context.Context, propagation.TextMapCarrier) {
	panic("TextMapPropagator.Inject is broken")
}

// A panic in the application's tracer provider, tracer, span or propagator is a
// failure to record telemetry: the call returns what it returns without it, the
// error handler of OpenTelemetry gets the panic, and every span that started ends.
func TestOpenTelemetryPanicsStayOutOfTheCall(t *testing.T) {
	calls := map[string]struct {
		transport func() *closureTransport
		spans     int32
		run       func(t *testing.T, client anthropic.Client)
	}{
		"create": {func() *closureTransport { return respond(200, "application/json", messageJSON) }, 1, func(t *testing.T, client anthropic.Client) {
			message, err := client.Messages.New(context.Background(), messageParams())
			if err != nil || message.ID != "msg_1" {
				t.Errorf("got the Message %q and the error %v, want msg_1", message.ID, err)
			}
		}},
		"stream": {func() *closureTransport { return respond(200, "text/event-stream", streamBody) }, 1, func(t *testing.T, client anthropic.Client) {
			stream := client.Messages.NewStreaming(context.Background(), messageParams())
			var message anthropic.Message
			for stream.Next() {
				if err := message.Accumulate(stream.Current()); err != nil {
					t.Fatal(err)
				}
			}
			if err := stream.Err(); err != nil || message.ID != "msg_1" || message.StopReason != anthropic.StopReasonEndTurn {
				t.Errorf("got the Message %q that stopped with %q and the error %v, want msg_1 and end_turn", message.ID, message.StopReason, err)
			}
		}},
		// The method records its result on the span itself.
		"count tokens": {func() *closureTransport { return respond(200, "application/json", `{"input_tokens":42}`) }, 1, func(t *testing.T, client anthropic.Client) {
			count, err := client.Messages.CountTokens(context.Background(), anthropic.MessageCountTokensParams{
				Model:    "claude-requested",
				Messages: []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hi"))},
			})
			if err != nil || count.InputTokens != 42 {
				t.Errorf("got %v and the error %v, want 42 input tokens", count, err)
			}
		}},
		// The loop, its two requests and its tool have a span each.
		"tool runner": {toolRunnerTransport, 4, func(t *testing.T, client anthropic.Client) {
			tool := &weatherTool{}
			runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{tool}, toolRunnerParams())
			message, err := runner.RunToCompletion(context.Background())
			if err != nil || message == nil || message.StopReason != anthropic.BetaStopReasonEndTurn || tool.calls != 1 {
				t.Errorf("got the Message %v and the error %v after %d runs of the tool, want the final Message after one", message, err, tool.calls)
			}
		}},
	}
	for _, panicsIn := range []string{"TracerProvider.Tracer", "Tracer.Start", "Span.SetAttributes", "Span.SetStatus", "Span.End", "TextMapPropagator.Inject"} {
		for name, call := range calls {
			t.Run(panicsIn+" in "+name, func(t *testing.T) {
				var mu sync.Mutex
				var handled []error
				otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
					mu.Lock()
					defer mu.Unlock()
					handled = append(handled, err)
				}))
				t.Cleanup(func() { otel.SetErrorHandler(otel.ErrorHandlerFunc(func(error) {})) })
				if panicsIn == "TextMapPropagator.Inject" {
					setPropagator(t, brokenPropagator{})
				} else {
					setPropagator(t, propagation.TraceContext{})
				}
				provider := &brokenProvider{panicsIn: panicsIn}
				client, _ := newTracedClient(t, call.transport(), option.WithTracerProvider(provider))

				call.run(t, client)

				mu.Lock()
				defer mu.Unlock()
				if len(handled) == 0 {
					t.Error("the error handler got nothing, want the panic")
				}
				for _, err := range handled {
					if !strings.Contains(err.Error(), panicsIn+" is broken") {
						t.Errorf("the error handler got %v, want the panic of %s", err, panicsIn)
					}
				}
				want := call.spans
				if panicsIn == "TracerProvider.Tracer" || panicsIn == "Tracer.Start" {
					want = 0
				}
				if started, ended := provider.started.Load(), provider.ended.Load(); started != want || ended != want {
					t.Errorf("%d spans started and %d ended, want %d of each", started, ended, want)
				}
			})
		}
	}
}
