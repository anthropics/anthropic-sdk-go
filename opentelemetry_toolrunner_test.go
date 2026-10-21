package anthropic_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/mcp"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// weatherTool is a tool that returns its result, or fails with err.
type weatherTool struct {
	mu    sync.Mutex
	err   error
	calls int
}

func (w *weatherTool) Name() string        { return "get_weather" }
func (w *weatherTool) Description() string { return "Gets the weather" }
func (w *weatherTool) InputSchema() anthropic.BetaToolInputSchemaParam {
	return anthropic.BetaToolInputSchemaParam{Properties: map[string]any{}}
}

func (w *weatherTool) Execute(ctx context.Context, _ json.RawMessage) ([]anthropic.BetaToolResultBlockParamContentUnion, error) {
	w.mu.Lock()
	w.calls++
	w.mu.Unlock()
	if w.err != nil {
		return nil, w.err
	}
	return []anthropic.BetaToolResultBlockParamContentUnion{{OfText: &anthropic.BetaTextBlockParam{Text: "sunny"}}}, nil
}

const toolUseJSON = `{"id":"msg_1","type":"message","role":"assistant","model":"claude-served","stop_reason":"tool_use",
	"content":[{"type":"tool_use","id":"toolu_1","name":"get_weather","input":{}},{"type":"tool_use","id":"toolu_2","name":"unknown","input":{}}],
	"usage":{"input_tokens":1,"output_tokens":1}}`

// toolRunnerTransport answers the first request with tool calls, and the next with a final message.
func toolRunnerTransport() *closureTransport {
	calls := 0
	return &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		calls++
		body := messageJSON
		if calls == 1 {
			body = toolUseJSON
		}
		return respond(200, "application/json", body).fn(req)
	}}
}

// toolUseStreamBody is the stream of a Message that calls get_weather.
const toolUseStreamBody = "event: message_start\n" +
	`data: {"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","model":"claude-served","content":[],"usage":{"input_tokens":1,"output_tokens":1}}}` + "\n\n" +
	"event: content_block_start\n" + `data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}` + "\n\n" +
	"event: content_block_delta\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Let me "}}` + "\n\n" +
	"event: content_block_delta\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"check."}}` + "\n\n" +
	"event: content_block_stop\n" + `data: {"type":"content_block_stop","index":0}` + "\n\n" +
	"event: content_block_start\n" +
	`data: {"type":"content_block_start","index":1,"content_block":{"type":"tool_use","id":"toolu_1","name":"get_weather","input":{}}}` + "\n\n" +
	"event: content_block_delta\n" +
	`data: {"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"{\"city\": "}}` + "\n\n" +
	"event: content_block_delta\n" +
	`data: {"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"\"Paris\"}"}}` + "\n\n" +
	"event: content_block_stop\n" + `data: {"type":"content_block_stop","index":1}` + "\n\n" +
	"event: message_delta\n" +
	`data: {"type":"message_delta","delta":{"stop_reason":"tool_use"},"usage":{"output_tokens":5}}` + "\n\n" +
	"event: message_stop\n" + `data: {"type":"message_stop"}` + "\n\n"

// toolRunnerStreamingTransport streams a tool call for the first request, and a final message for the next.
func toolRunnerStreamingTransport() *closureTransport {
	calls := 0
	return &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		calls++
		body := streamBody
		if calls == 1 {
			body = toolUseStreamBody
		}
		return respond(200, "text/event-stream", body).fn(req)
	}}
}

func toolRunnerParams() anthropic.BetaToolRunnerParams {
	return anthropic.BetaToolRunnerParams{BetaMessageNewParams: anthropic.BetaMessageNewParams{
		Model:     "claude-requested",
		MaxTokens: 100,
		Messages:  []anthropic.BetaMessageParam{anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("weather?"))},
	}}
}

func TestOpenTelemetryToolRunner(t *testing.T) {
	tool := &weatherTool{err: errors.New("no forecast")}
	client, recorder := newTracedClient(t, toolRunnerTransport())

	ctx, app := recorder.Tracer("app").Start(context.Background(), "handle_order")
	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{tool}, toolRunnerParams(),
		option.WithOpenTelemetryConversationID("conv_1"))
	if _, err := runner.RunToCompletion(ctx); err != nil {
		t.Fatal(err)
	}
	app.End()

	var names []string
	byName := map[string][]sdktrace.ReadOnlySpan{}
	for _, s := range recorder.Spans() {
		names = append(names, s.Name())
		byName[s.Name()] = append(byName[s.Name()], s)
	}
	toolRunnerSpan := byName["anthropic.messages.tool_runner"]
	uses := byName["anthropic.custom_tool_use get_weather"]
	creates := byName["anthropic.messages.create"]
	if len(toolRunnerSpan) != 1 || len(uses) != 1 || len(creates) != 2 {
		t.Fatalf("got spans %q, want one tool runner span, one tool use, and two creates", names)
	}
	if toolRunnerSpan[0].SpanKind() != trace.SpanKindInternal || toolRunnerSpan[0].Status().Code != codes.Ok {
		t.Errorf("tool runner span of kind %v with status %v, want INTERNAL with OK", toolRunnerSpan[0].SpanKind(), toolRunnerSpan[0].Status().Code)
	}
	if toolRunnerSpan[0].Parent().SpanID() != app.SpanContext().SpanID() {
		t.Error("the tool runner span is not a child of the application's span")
	}
	for _, child := range append(uses, creates...) {
		if child.Parent().SpanID() != toolRunnerSpan[0].SpanContext().SpanID() {
			t.Errorf("span %q is not a child of the tool runner span", child.Name())
		}
		wantAttrs(t, child, map[string]any{"gen_ai.conversation.id": "conv_1"})
	}
	for i, reason := range []string{"tool_use", "end_turn"} {
		wantAttrs(t, creates[i], map[string]any{
			"anthropic.message.stop_reason":  reason,
			"gen_ai.response.finish_reasons": nil,
			"http.request.method":            "POST",
			"http.route":                     "/v1/messages",
			"url.full":                       "http://localhost:4010/v1/messages",
		})
	}
	// Only a request has a method, a route and a URL: the tool runner's span and
	// the tool's send none.
	for _, helper := range append(toolRunnerSpan, uses...) {
		wantAttrs(t, helper, absent(map[string]any{}, append([]string{"http.request.method", "http.route", "http.request.resend_count", "server.address"}, urlAttrs...)...))
	}
	wantAttrs(t, toolRunnerSpan[0], map[string]any{
		"gen_ai.operation.name":  "invoke_agent",
		"gen_ai.request.model":   "claude-requested",
		"gen_ai.provider.name":   "anthropic",
		"gen_ai.conversation.id": "conv_1",
		"gen_ai.output.type":     nil,
		"anthropic.request.id":   nil,
		// The stop reason is each request's, not the loop's.
		"anthropic.message.stop_reason":  nil,
		"gen_ai.response.finish_reasons": nil,
		"error.type":                     nil,
	})
	if uses[0].SpanKind() != trace.SpanKindInternal || uses[0].Status().Code != codes.Error {
		t.Errorf("tool span of kind %v with status %v, want INTERNAL with ERROR", uses[0].SpanKind(), uses[0].Status().Code)
	}
	wantAttrs(t, uses[0], map[string]any{
		"gen_ai.operation.name": "execute_tool",
		"gen_ai.tool.name":      "get_weather",
		"gen_ai.tool.call.id":   "toolu_1",
		"gen_ai.tool.type":      "function",
		"gen_ai.provider.name":  "anthropic",
		"error.type":            "tool_error",
		// The description that the application gave is metadata: the client records no content.
		"gen_ai.tool.description":    "Gets the weather",
		"gen_ai.tool.call.arguments": nil,
		// Only the Sessions tool runner has a Session.
		"anthropic.session.id":        nil,
		"anthropic.session.thread.id": nil,
	})
}

// undescribedTool is a weatherTool whose definition has the description
// description: none, unless the test sets one.
type undescribedTool struct {
	weatherTool
	description string
}

func (u *undescribedTool) Description() string { return u.description }

func TestOpenTelemetryToolWithoutADescriptionHasNoDescriptionAttribute(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	// A description that is not valid UTF-8 can fail the whole export, so the span
	// leaves it out as it does an empty one, and keeps its other attributes.
	for name, description := range map[string]string{"none": "", "not valid UTF-8": "Gets the \xff\xfe weather"} {
		t.Run(name, func(t *testing.T) {
			client, recorder := newTracedClient(t, toolRunnerTransport(), content)
			runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&undescribedTool{description: description}}, toolRunnerParams())
			if _, err := runner.RunToCompletion(context.Background()); err != nil {
				t.Fatal(err)
			}
			found := false
			for _, s := range recorder.Spans() {
				if s.Name() == "anthropic.custom_tool_use get_weather" {
					found = true
					wantAttrs(t, s, map[string]any{"gen_ai.tool.name": "get_weather", "gen_ai.tool.call.id": "toolu_1", "gen_ai.tool.description": nil})
					if _, ok := attr(s, "anthropic.tool_result.content"); !ok {
						t.Error("attribute anthropic.tool_result.content is absent, want only the description left out")
					}
				}
			}
			if !found {
				t.Fatal("got no anthropic.custom_tool_use span")
			}
		})
	}
}

func TestOpenTelemetryStreamingToolRunnerRecordsTheToolDescription(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	client, recorder := newTracedClient(t, toolRunnerStreamingTransport(), content)

	runner := client.Beta.Messages.NewToolRunnerStreaming([]anthropic.BetaTool{&weatherTool{}}, toolRunnerParams())
	for events, err := range runner.AllStreaming(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		for _, err := range events {
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	found := false
	for _, s := range recorder.Spans() {
		if s.Name() == "anthropic.custom_tool_use get_weather" {
			found = true
			wantAttrs(t, s, map[string]any{"gen_ai.tool.description": "Gets the weather"})
		}
	}
	if !found {
		t.Fatal("got no anthropic.custom_tool_use span")
	}
}

func TestOpenTelemetryToolRunnerOutputType(t *testing.T) {
	format := anthropic.BetaJSONOutputFormatParam{Schema: map[string]any{"type": "object", "properties": map[string]any{}}}
	for name, set := range map[string]func(params *anthropic.BetaToolRunnerParams){
		"output_config.format":     func(params *anthropic.BetaToolRunnerParams) { params.OutputConfig.Format = format },
		"deprecated output_format": func(params *anthropic.BetaToolRunnerParams) { params.OutputFormat = format },
	} {
		t.Run(name, func(t *testing.T) {
			client, recorder := newTracedClient(t, toolRunnerTransport())
			params := toolRunnerParams()
			set(&params)
			runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&weatherTool{}}, params)
			if _, err := runner.RunToCompletion(context.Background()); err != nil {
				t.Fatal(err)
			}
			// The loop has it from the runner's parameters, and each request from its own.
			for _, s := range recorder.Spans() {
				want := any("json")
				if strings.HasPrefix(s.Name(), "anthropic.custom_tool_use") {
					want = nil
				}
				wantAttrs(t, s, map[string]any{"gen_ai.output.type": want})
			}
		})
	}
}

// tracingTool starts a span of its own in the context that the runner runs it in.
type tracingTool struct {
	weatherTool
	tracer trace.Tracer
}

func (w *tracingTool) Execute(ctx context.Context, input json.RawMessage) ([]anthropic.BetaToolResultBlockParamContentUnion, error) {
	_, span := w.tracer.Start(ctx, "look_up_forecast")
	defer span.End()
	return w.weatherTool.Execute(ctx, input)
}

func TestOpenTelemetryToolRunsInItsSpan(t *testing.T) {
	client, recorder := newTracedClient(t, toolRunnerTransport())
	tracer := recorder.Tracer("app")

	ctx, app := tracer.Start(context.Background(), "handle_order")
	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&tracingTool{tracer: tracer}}, toolRunnerParams())
	if _, err := runner.RunToCompletion(ctx); err != nil {
		t.Fatal(err)
	}
	// The runner's spans were never current in the caller's context.
	if got := trace.SpanContextFromContext(ctx).SpanID(); got != app.SpanContext().SpanID() {
		t.Errorf("the span in the caller's context is %s after the run, want the caller's %s", got, app.SpanContext().SpanID())
	}
	_, after := tracer.Start(ctx, "after_the_run")
	after.End()
	app.End()

	byName := map[string]sdktrace.ReadOnlySpan{}
	for _, s := range recorder.Spans() {
		byName[s.Name()] = s
	}
	use, inTool := byName["anthropic.custom_tool_use get_weather"], byName["look_up_forecast"]
	if use == nil || inTool == nil {
		t.Fatalf("got spans %v, want the tool's span and the one it started", byName)
	}
	if inTool.Parent().SpanID() != use.SpanContext().SpanID() {
		t.Error("the span that the tool started is not a child of the tool's span")
	}
	// The application set no conversation ID on the runner.
	wantAttrs(t, use, map[string]any{"gen_ai.conversation.id": nil})
	wantAttrs(t, byName["anthropic.messages.create"], map[string]any{"gen_ai.conversation.id": nil})
	if byName["after_the_run"].Parent().SpanID() != app.SpanContext().SpanID() {
		t.Error("a span started after the run is not a child of the caller's span")
	}
}

func TestOpenTelemetryToolErrorTextStaysOffSpans(t *testing.T) {
	const text = "forecast-service-password-hunter2"
	// Where the text of the tool's error may be, by span name: in content mode it
	// is content of the tool's result, and of the next request, which sends it on.
	for name, test := range map[string]struct {
		opts    []option.RequestOption
		allowed map[string]string
	}{
		"metadata only": {},
		"content": {
			opts: []option.RequestOption{option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)},
			allowed: map[string]string{
				"anthropic.custom_tool_use get_weather": "anthropic.tool_result.content",
				"anthropic.messages.create":             "gen_ai.input.messages",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			client, recorder := newTracedClient(t, toolRunnerTransport(), test.opts...)
			runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&weatherTool{err: errors.New(text)}}, toolRunnerParams())
			if _, err := runner.RunToCompletion(context.Background()); err != nil {
				t.Fatal(err)
			}
			spans := recorder.Spans()
			if len(spans) != 4 {
				t.Fatalf("got %d spans, want the tool runner span, the tool use, and two creates", len(spans))
			}
			found := map[string]bool{}
			for _, span := range spans {
				if strings.Contains(span.Status().Description, text) {
					t.Errorf("%s: the status description has the text of the tool's error", span.Name())
				}
				if len(span.Events()) != 0 {
					t.Errorf("%s: got %d span events, want none", span.Name(), len(span.Events()))
				}
				for _, kv := range span.Attributes() {
					if !strings.Contains(kv.Value.Emit(), text) {
						continue
					}
					if string(kv.Key) != test.allowed[span.Name()] {
						t.Errorf("%s: attribute %s has the text of the tool's error", span.Name(), kv.Key)
					}
					found[span.Name()] = true
				}
			}
			if len(found) != len(test.allowed) {
				t.Errorf("the text of the tool's error is on spans %v, want it in %v", found, test.allowed)
			}
		})
	}
}

// mcpWeatherTool is a get_weather tool that the MCP helper made, over an
// in-memory MCP server.
func mcpWeatherTool(t *testing.T) anthropic.BetaTool {
	t.Helper()
	ctx := context.Background()
	server := mcpsdk.NewServer(&mcpsdk.Implementation{Name: "weather", Version: "1"}, nil)
	definition := &mcpsdk.Tool{Name: "get_weather", Description: "The weather, as the server describes it", InputSchema: map[string]any{"type": "object"}}
	mcpsdk.AddTool(server, definition, func(context.Context, *mcpsdk.CallToolRequest, map[string]any) (*mcpsdk.CallToolResult, any, error) {
		return &mcpsdk.CallToolResult{Content: []mcpsdk.Content{&mcpsdk.TextContent{Text: "sunny"}}}, nil, nil
	})
	serverTransport, clientTransport := mcpsdk.NewInMemoryTransports()
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = serverSession.Close() })
	session, err := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "test", Version: "1"}, nil).Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = session.Close() })
	tool, err := mcp.NewBetaTool(definition, session)
	if err != nil {
		t.Fatal(err)
	}
	return tool
}

func TestOpenTelemetryToolTheMCPHelperMadeRunsInAnMCPToolUseSpan(t *testing.T) {
	client, recorder := newTracedClient(t, toolRunnerTransport())

	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{mcpWeatherTool(t)}, toolRunnerParams())
	if _, err := runner.RunToCompletion(context.Background()); err != nil {
		t.Fatal(err)
	}

	var names []string
	var use sdktrace.ReadOnlySpan
	for _, s := range recorder.Spans() {
		names = append(names, s.Name())
		if s.Name() == "anthropic.mcp_tool_use get_weather" {
			use = s
		}
		if strings.HasPrefix(s.Name(), "anthropic.custom_tool_use") {
			t.Errorf("span %q, want no custom tool span", s.Name())
		}
	}
	if use == nil {
		t.Fatalf("got spans %q, want an anthropic.mcp_tool_use span", names)
	}
	if use.SpanKind() != trace.SpanKindInternal || use.Status().Code != codes.Ok {
		t.Errorf("tool span of kind %v with status %v, want INTERNAL with OK", use.SpanKind(), use.Status().Code)
	}
	wantAttrs(t, use, map[string]any{
		"gen_ai.operation.name": "execute_tool",
		"gen_ai.tool.name":      "get_weather",
		"gen_ai.tool.call.id":   "toolu_1",
		"gen_ai.tool.type":      "function",
		"gen_ai.provider.name":  "anthropic",
		// The MCP server's description, which the helper's tool definition has.
		"gen_ai.tool.description": "The weather, as the server describes it",
	})
}

func TestOpenTelemetryToolRunnerFailedCall(t *testing.T) {
	client, recorder := newTracedClient(t, respond(400, "application/json",
		`{"type":"error","error":{"type":"invalid_request_error","message":"bad"}}`))

	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&weatherTool{}}, toolRunnerParams())
	if _, err := runner.NextMessage(context.Background()); err == nil {
		t.Fatal("the runner did not fail")
	}
	for _, s := range recorder.Spans() {
		if s.Status().Code != codes.Error {
			t.Errorf("span %q has status %v, want ERROR", s.Name(), s.Status().Code)
		}
		wantAttrs(t, s, map[string]any{"error.type": "invalid_request_error"})
	}
	if got := len(recorder.Spans()); got != 2 {
		t.Errorf("got %d spans, want the tool runner span and the call", got)
	}
}

func TestOpenTelemetryToolRunnerStreamingFailedCall(t *testing.T) {
	client, recorder := newTracedClient(t, respond(400, "application/json",
		`{"type":"error","error":{"type":"invalid_request_error","message":"bad"}}`))

	// The request is sent when the loop opens the stream, and its error arrives as the application reads it.
	runner := client.Beta.Messages.NewToolRunnerStreaming([]anthropic.BetaTool{&weatherTool{}}, toolRunnerParams())
	var failed error
	for events := range runner.AllStreaming(context.Background()) {
		for _, err := range events {
			if err != nil {
				failed = err
			}
		}
		if failed != nil {
			break
		}
	}
	if failed == nil {
		t.Fatal("the runner did not fail")
	}
	var toolRunnerSpan sdktrace.ReadOnlySpan
	for _, s := range recorder.Started() {
		if s.Name() == "anthropic.messages.tool_runner" {
			toolRunnerSpan = s
		}
	}
	if toolRunnerSpan == nil || !ended(toolRunnerSpan) {
		t.Fatal("the tool runner span did not end")
	}
	if toolRunnerSpan.Status().Code != codes.Error {
		t.Errorf("status %v, want ERROR", toolRunnerSpan.Status().Code)
	}
	wantAttrs(t, toolRunnerSpan, map[string]any{"error.type": "invalid_request_error"})
}

// canceledTool blocks until the application cancels its context, then fails
// with an error of its own rather than the context's.
type canceledTool struct {
	weatherTool
	started chan struct{}
}

func (c *canceledTool) Execute(ctx context.Context, _ json.RawMessage) ([]anthropic.BetaToolResultBlockParamContentUnion, error) {
	close(c.started)
	<-ctx.Done()
	return nil, errors.New("gave up")
}

func TestOpenTelemetryToolRunCanceledByTheApplicationLeavesItsSpanAndTheLoopSpanUnset(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	client, recorder := newTracedClient(t, toolRunnerTransport(), content)
	tool := &canceledTool{started: make(chan struct{})}
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-tool.started
		cancel()
	}()

	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{tool}, toolRunnerParams())
	if _, err := runner.RunToCompletion(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want the cancellation", err)
	}
	// A cancellation is the application stopping the run, which is no failure of the tool or the loop.
	for _, name := range []string{"anthropic.custom_tool_use get_weather", "anthropic.messages.tool_runner"} {
		var span sdktrace.ReadOnlySpan
		for _, s := range recorder.Started() {
			if s.Name() == name {
				span = s
			}
		}
		if span == nil || !ended(span) {
			t.Fatalf("the %s span did not end", name)
		}
		if span.Status().Code != codes.Unset {
			t.Errorf("%s: status %v, want unset", name, span.Status().Code)
		}
		// The runner made no result of the canceled run.
		wantAttrs(t, span, map[string]any{"error.type": nil, "anthropic.tool_result.content": nil, "gen_ai.tool.call.result": nil})
	}
}

func TestOpenTelemetryToolRunnerStopped(t *testing.T) {
	for name, run := range map[string]func(context.Context, anthropic.Client){
		"break": func(ctx context.Context, client anthropic.Client) {
			runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&weatherTool{}}, toolRunnerParams())
			for range runner.All(ctx) {
				break
			}
		},
		"break streaming": func(ctx context.Context, client anthropic.Client) {
			runner := client.Beta.Messages.NewToolRunnerStreaming([]anthropic.BetaTool{&weatherTool{}}, toolRunnerParams())
			for events := range runner.AllStreaming(ctx) {
				for range events {
					break
				}
				break
			}
		},
		"break streaming after a turn": func(ctx context.Context, client anthropic.Client) {
			runner := client.Beta.Messages.NewToolRunnerStreaming([]anthropic.BetaTool{&weatherTool{}}, toolRunnerParams())
			for events := range runner.AllStreaming(ctx) {
				for range events {
				}
				break
			}
		},
		"cancel": func(ctx context.Context, client anthropic.Client) {
			ctx, cancel := context.WithCancel(ctx)
			cancel()
			runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&weatherTool{}}, toolRunnerParams())
			_, _ = runner.RunToCompletion(ctx)
		},
	} {
		t.Run(name, func(t *testing.T) {
			rt := toolRunnerTransport()
			if strings.Contains(name, "streaming") {
				rt = respond(200, "text/event-stream", streamBody)
			}
			client, recorder := newTracedClient(t, rt)
			run(context.Background(), client)

			var toolRunnerSpan sdktrace.ReadOnlySpan
			for _, s := range recorder.Started() {
				if s.Name() == "anthropic.messages.tool_runner" {
					toolRunnerSpan = s
				}
			}
			if toolRunnerSpan == nil || !ended(toolRunnerSpan) {
				t.Fatal("the tool runner span did not end")
			}
			if toolRunnerSpan.Status().Code != codes.Unset {
				t.Errorf("status %v, want unset", toolRunnerSpan.Status().Code)
			}
			wantAttrs(t, toolRunnerSpan, map[string]any{"error.type": nil})
		})
	}
}

func TestOpenTelemetryToolContent(t *testing.T) {
	// The result is the content blocks of the tool_result block, for an error result too.
	for name, test := range map[string]struct {
		tool anthropic.BetaTool
		want string
	}{
		"result": {&weatherTool{}, `[{"type":"text","text":"sunny"}]`},
		"error":  {&weatherTool{err: errors.New("no forecast")}, `[{"type":"text","text":"Error: no forecast"}]`},
		// Text and search results keep their fields but the encrypted index of the
		// citation, and any other block holds its type only.
		"blocks": {&blocksTool{}, `[{"type":"text","text":"sunny","citations":[{"type":"web_search_result_location","cited_text":"sun","url":"https://example.com"}],
			"cache_control":{"type":"ephemeral"}},
			{"type":"search_result","source":"https://example.com","title":"Forecast","content":[{"type":"text","text":"rain"}]},{"type":"image"}]`},
		"no content": {&blocksTool{empty: true}, `[]`},
	} {
		t.Run(name, func(t *testing.T) {
			content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
			client, recorder := newTracedClient(t, toolRunnerTransport(), content)

			runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{test.tool}, toolRunnerParams())
			if _, err := runner.RunToCompletion(context.Background()); err != nil {
				t.Fatal(err)
			}
			for _, s := range recorder.Spans() {
				if s.Name() == "anthropic.custom_tool_use get_weather" {
					wantJSON(t, s, "gen_ai.tool.call.arguments", `{}`)
					wantJSON(t, s, "anthropic.tool_result.content", test.want)
					wantAttrs(t, s, map[string]any{"gen_ai.tool.call.result": nil, "anthropic.content.truncated": nil})
					return
				}
			}
			t.Fatal("the tool's span did not end")
		})
	}
}

// textsTool returns one text block for each of texts.
type textsTool struct {
	weatherTool
	texts []string
}

func (b *textsTool) Execute(context.Context, json.RawMessage) ([]anthropic.BetaToolResultBlockParamContentUnion, error) {
	var blocks []anthropic.BetaToolResultBlockParamContentUnion
	for _, text := range b.texts {
		blocks = append(blocks, anthropic.BetaToolResultBlockParamContentUnion{OfText: &anthropic.BetaTextBlockParam{Text: text}})
	}
	return blocks, nil
}

// toolsetToolUseJSON is a Message that calls a member of the computer toolset
// and a tool that is in no toolset.
const toolsetToolUseJSON = `{"id":"msg_1","type":"message","role":"assistant","model":"claude-served","stop_reason":"tool_use",
	"content":[{"type":"tool_use","id":"toolu_1","name":"screenshot","toolset_name":"computer","input":{}},
		{"type":"tool_use","id":"toolu_2","name":"get_weather","toolset_name":null,"input":{}}],
	"usage":{"input_tokens":1,"output_tokens":1}}`

type screenshotTool struct{ weatherTool }

func (*screenshotTool) Name() string { return "screenshot" }

func TestOpenTelemetryToolsetName(t *testing.T) {
	calls := 0
	client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		calls++
		body := messageJSON
		if calls == 1 {
			body = toolsetToolUseJSON
		}
		return respond(200, "application/json", body).fn(req)
	}})

	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&screenshotTool{}, &weatherTool{}}, toolRunnerParams())
	if _, err := runner.RunToCompletion(context.Background()); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"anthropic.custom_tool_use screenshot":  "computer",
		"anthropic.custom_tool_use get_weather": nil,
	}
	for _, s := range recorder.Spans() {
		if toolset, ok := want[s.Name()]; ok {
			wantAttrs(t, s, map[string]any{"anthropic.toolset.name": toolset})
			delete(want, s.Name())
		}
	}
	if len(want) != 0 {
		t.Errorf("the spans of %v did not end", want)
	}
}

func TestOpenTelemetryToolContentIsNotHTMLEscaped(t *testing.T) {
	calls := 0
	rt := &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		calls++
		body := messageJSON
		if calls == 1 {
			body = markupResponseJSON
		}
		return respond(200, "application/json", body).fn(req)
	}}
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	client, recorder := newTracedClient(t, rt, content)

	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&textsTool{texts: []string{markupText}}}, toolRunnerParams())
	if _, err := runner.RunToCompletion(context.Background()); err != nil {
		t.Fatal(err)
	}
	var tool, second sdktrace.ReadOnlySpan
	for _, s := range recorder.Spans() {
		switch s.Name() {
		case "anthropic.custom_tool_use get_weather":
			tool = s
		case "anthropic.messages.create":
			second = s
		}
	}
	if tool == nil || second == nil {
		t.Fatal("the span of the tool or of the second request did not end")
	}
	wantMarkupAsTyped(t, tool, "gen_ai.tool.call.arguments", 1)
	wantMarkupAsTyped(t, tool, "anthropic.tool_result.content", 1)
	// The second request sends back the text, the tool call and the tool's result.
	wantMarkupAsTyped(t, second, "gen_ai.input.messages", 3)
}

func TestOpenTelemetryToolContentSizeLimit(t *testing.T) {
	big := strings.Repeat("x", 600<<10)
	const first = `[{"type":"text","text":"first"}]`
	// Blocks go from the end of the result, whole, until the span's content fits.
	for name, test := range map[string]struct {
		texts []string
		// maxBytes is the maximum that the client's options set, or 0 when they set none.
		maxBytes  int64
		want      string
		truncated any
	}{
		"last block":  {texts: []string{"first", big, big}, want: `[{"type":"text","text":"first"},{"type":"text","text":"` + big + `"}]`, truncated: true},
		"every block": {texts: []string{big + big, "last"}, want: `[]`, truncated: true},
		// The maximum counts the tool's arguments, {}, too.
		"last block over the maximum of the options": {texts: []string{"first", "second"}, maxBytes: int64(len(first) + 2), want: first, truncated: true},
		"one byte over the maximum of the options":   {texts: []string{"first", "second"}, maxBytes: int64(len(first) + 1), want: `[]`, truncated: true},
		"a maximum of the options that all blocks fit": {
			texts: []string{"first", big, big}, maxBytes: 4 << 20,
			want: `[{"type":"text","text":"first"},{"type":"text","text":"` + big + `"},{"type":"text","text":"` + big + `"}]`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			opts := []option.RequestOption{option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)}
			if test.maxBytes != 0 {
				opts = append(opts, option.WithOpenTelemetryTracesMaxContentBytes(test.maxBytes))
			}
			client, recorder := newTracedClient(t, toolRunnerTransport(), opts...)

			runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&textsTool{texts: test.texts}}, toolRunnerParams())
			if _, err := runner.RunToCompletion(context.Background()); err != nil {
				t.Fatal(err)
			}
			for _, s := range recorder.Spans() {
				if s.Name() == "anthropic.custom_tool_use get_weather" {
					wantJSON(t, s, "gen_ai.tool.call.arguments", `{}`)
					wantJSON(t, s, "anthropic.tool_result.content", test.want)
					wantAttrs(t, s, map[string]any{"anthropic.content.truncated": test.truncated})
					return
				}
			}
			t.Fatal("the tool's span did not end")
		})
	}
}

func TestOpenTelemetryToolArgumentsThatAloneAreOverTheSizeLimitAreLeftOut(t *testing.T) {
	const arguments = `{"city":"Paris, the capital of France"}`
	const ok = `[{"type":"text","text":"ok"}]`
	toolUse := strings.Replace(toolUseJSON, `"name":"get_weather","input":{}`, `"name":"get_weather","input":`+arguments, 1)
	for name, test := range map[string]struct {
		maxBytes int64
		// arguments is the attribute's value, or "" for an absent attribute.
		arguments, result string
	}{
		// The result has what the arguments leave of the maximum, which is nothing.
		"arguments that fill the maximum stay": {maxBytes: int64(len(arguments)), arguments: arguments, result: `[]`},
		// It has all of it when they are left out.
		"one byte over the maximum, they are left out": {maxBytes: int64(len(arguments) - 1), result: ok},
	} {
		t.Run(name, func(t *testing.T) {
			calls := 0
			rt := &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
				calls++
				body := messageJSON
				if calls == 1 {
					body = toolUse
				}
				return respond(200, "application/json", body).fn(req)
			}}
			client, recorder := newTracedClient(t, rt,
				option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent),
				option.WithOpenTelemetryTracesMaxContentBytes(test.maxBytes),
			)

			runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&textsTool{texts: []string{"ok"}}}, toolRunnerParams())
			if _, err := runner.RunToCompletion(context.Background()); err != nil {
				t.Fatal(err)
			}
			for _, s := range recorder.Spans() {
				if s.Name() == "anthropic.custom_tool_use get_weather" {
					if test.arguments == "" {
						wantAttrs(t, s, map[string]any{"gen_ai.tool.call.arguments": nil})
					} else {
						wantJSON(t, s, "gen_ai.tool.call.arguments", test.arguments)
					}
					wantJSON(t, s, "anthropic.tool_result.content", test.result)
					wantAttrs(t, s, map[string]any{"anthropic.content.truncated": true})
					return
				}
			}
			t.Fatal("the tool's span did not end")
		})
	}
}

// blocksTool returns a text block with a citation and a cache breakpoint, a
// search result and an image, or no content when empty.
type blocksTool struct {
	weatherTool
	empty bool
}

func (b *blocksTool) Execute(context.Context, json.RawMessage) ([]anthropic.BetaToolResultBlockParamContentUnion, error) {
	if b.empty {
		return nil, nil
	}
	return []anthropic.BetaToolResultBlockParamContentUnion{
		{OfText: &anthropic.BetaTextBlockParam{
			Text:         "sunny",
			Citations:    []anthropic.BetaTextCitationParamUnion{{OfWebSearchResultLocation: &anthropic.BetaCitationWebSearchResultLocationParam{CitedText: "sun", URL: "https://example.com", EncryptedIndex: "opaque"}}},
			CacheControl: anthropic.NewBetaCacheControlEphemeralParam(),
		}},
		{OfSearchResult: &anthropic.BetaSearchResultBlockParam{Source: "https://example.com", Title: "Forecast", Content: []anthropic.BetaTextBlockParam{{Text: "rain"}}}},
		{OfImage: &anthropic.BetaImageBlockParam{Source: anthropic.BetaImageBlockParamSourceUnion{
			OfBase64: &anthropic.BetaBase64ImageSourceParam{Data: "AAAA", MediaType: anthropic.BetaBase64ImageSourceMediaTypeImagePNG},
		}}},
	}, nil
}

// unencodableTool returns content that JSON cannot encode.
type unencodableTool struct{ weatherTool }

func (*unencodableTool) Execute(context.Context, json.RawMessage) ([]anthropic.BetaToolResultBlockParamContentUnion, error) {
	block := &anthropic.BetaTextBlockParam{Text: "sunny"}
	block.SetExtraFields(map[string]any{"channel": make(chan int)})
	return []anthropic.BetaToolResultBlockParamContentUnion{{OfText: block}}, nil
}

func TestOpenTelemetryToolContentThatCannotBeEncoded(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	client, recorder := newTracedClient(t, toolRunnerTransport(), content)

	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&unencodableTool{}}, toolRunnerParams())
	// The run goes on to send the result, which fails to encode as the span's copy did.
	_, _ = runner.RunToCompletion(context.Background())
	for _, s := range recorder.Spans() {
		if s.Name() == "anthropic.custom_tool_use get_weather" {
			if s.Status().Code != codes.Ok {
				t.Errorf("status %v, want OK", s.Status().Code)
			}
			wantJSON(t, s, "gen_ai.tool.call.arguments", `{}`)
			wantAttrs(t, s, map[string]any{"anthropic.tool_result.content": nil})
			return
		}
	}
	t.Fatal("the tool's span did not end")
}

// startedToolRunnerSpans returns the tool runner spans,
// anthropic.messages.tool_runner, that recorder has seen start, ended or not.
func startedToolRunnerSpans(recorder *spanRecorder) []sdktrace.ReadOnlySpan {
	var toolRunnerSpans []sdktrace.ReadOnlySpan
	for _, s := range recorder.Started() {
		if s.Name() == "anthropic.messages.tool_runner" {
			toolRunnerSpans = append(toolRunnerSpans, s)
		}
	}
	return toolRunnerSpans
}

// wantStoppedToolRunnerSpan fails the test unless recorder has one tool runner span, which ended
// as a run that the application stopped: with no status and no error.type, and
// with the attributes that it started with.
func wantStoppedToolRunnerSpan(t *testing.T, recorder *spanRecorder) sdktrace.ReadOnlySpan {
	t.Helper()
	toolRunnerSpans := startedToolRunnerSpans(recorder)
	if len(toolRunnerSpans) != 1 {
		t.Fatalf("got %d tool runner spans, want 1", len(toolRunnerSpans))
	}
	toolRunnerSpan := toolRunnerSpans[0]
	if !ended(toolRunnerSpan) {
		t.Fatal("the tool runner span has not ended")
	}
	if toolRunnerSpan.Status().Code != codes.Unset {
		t.Errorf("tool runner span status %v, want none", toolRunnerSpan.Status().Code)
	}
	wantAttrs(t, toolRunnerSpan, map[string]any{
		"gen_ai.operation.name": "invoke_agent",
		"gen_ai.request.model":  "claude-requested",
		"gen_ai.provider.name":  "anthropic",
		"error.type":            nil,
	})
	return toolRunnerSpan
}

// waitForToolRunnerSpan collects garbage until recorder has the tool runner span as ended, or fails the test.
func waitForToolRunnerSpan(t *testing.T, recorder *spanRecorder) {
	t.Helper()
	for range 100 {
		runtime.GC()
		for _, s := range recorder.Spans() {
			if s.Name() == "anthropic.messages.tool_runner" {
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("the span of a dropped runner is still open after 100 collections")
}

func TestOpenTelemetryToolRunnerDroppedBetweenStepsEndsTheLoopSpan(t *testing.T) {
	ctx := context.Background()

	t.Run("messages", func(t *testing.T) {
		client, recorder := newTracedClient(t, toolRunnerTransport())
		func() {
			runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&weatherTool{}}, toolRunnerParams())
			if message, err := runner.NextMessage(ctx); err != nil || message == nil {
				t.Fatalf("got message %v, error %v; want the turn that calls the tool", message, err)
			}
		}()

		waitForToolRunnerSpan(t, recorder)
		wantStoppedToolRunnerSpan(t, recorder)
	})

	t.Run("streaming", func(t *testing.T) {
		client, recorder := newTracedClient(t, toolRunnerStreamingTransport())
		func() {
			runner := client.Beta.Messages.NewToolRunnerStreaming([]anthropic.BetaTool{&weatherTool{}}, toolRunnerParams())
			for _, err := range runner.NextStreaming(ctx) {
				if err != nil {
					t.Fatal(err)
				}
			}
		}()

		waitForToolRunnerSpan(t, recorder)
		wantStoppedToolRunnerSpan(t, recorder)
	})
}

func TestOpenTelemetryToolRunnerInUseKeepsTheLoopSpanOpen(t *testing.T) {
	ctx := context.Background()
	client, recorder := newTracedClient(t, toolRunnerTransport())
	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&weatherTool{}}, toolRunnerParams())
	if message, err := runner.NextMessage(ctx); err != nil || message == nil {
		t.Fatalf("got message %v, error %v; want the turn that calls the tool", message, err)
	}

	for range 5 {
		runtime.GC()
		time.Sleep(10 * time.Millisecond)
	}
	if toolRunnerSpans := startedToolRunnerSpans(recorder); len(toolRunnerSpans) != 1 || ended(toolRunnerSpans[0]) {
		t.Fatal("want the tool runner span open while the application holds the runner")
	}

	if _, err := runner.RunToCompletion(ctx); err != nil {
		t.Fatal(err)
	}
	toolRunnerSpans := startedToolRunnerSpans(recorder)
	if len(toolRunnerSpans) != 1 || !ended(toolRunnerSpans[0]) || toolRunnerSpans[0].Status().Code != codes.Ok {
		t.Fatalf("got %d tool runner spans, want the one of the completed loop, with status OK", len(toolRunnerSpans))
	}
}

// panicsAtTheEndOfAToolRunnerSpan is a span processor with a bug.
type panicsAtTheEndOfAToolRunnerSpan struct{}

func (panicsAtTheEndOfAToolRunnerSpan) OnStart(context.Context, sdktrace.ReadWriteSpan) {}
func (panicsAtTheEndOfAToolRunnerSpan) Shutdown(context.Context) error                  { return nil }
func (panicsAtTheEndOfAToolRunnerSpan) ForceFlush(context.Context) error                { return nil }
func (panicsAtTheEndOfAToolRunnerSpan) OnEnd(s sdktrace.ReadOnlySpan) {
	if s.Name() == "anthropic.messages.tool_runner" {
		panic("exporter bug")
	}
}

// The span of a dropped runner ends on a goroutine that the application cannot recover on.
func TestOpenTelemetryToolRunnerDroppedSurvivesASpanProcessorThatPanics(t *testing.T) {
	handled := make(chan error, 1)
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) { handled <- err }))
	t.Cleanup(func() { otel.SetErrorHandler(otel.ErrorHandlerFunc(func(error) {})) })
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(panicsAtTheEndOfAToolRunnerSpan{}))
	client, _ := newTracedClient(t, toolRunnerTransport(), option.WithTracerProvider(provider))
	func() {
		runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&weatherTool{}}, toolRunnerParams())
		if message, err := runner.NextMessage(context.Background()); err != nil || message == nil {
			t.Fatalf("got message %v, error %v; want the turn that calls the tool", message, err)
		}
	}()

	for range 100 {
		runtime.GC()
		select {
		case err := <-handled:
			if !strings.Contains(err.Error(), "exporter bug") {
				t.Fatalf("the error handler got %v, want the panic of the span processor", err)
			}
			return
		case <-time.After(10 * time.Millisecond):
		}
	}
	t.Fatal("the error handler got nothing after 100 collections")
}

// blocksAtTheEndOfAToolRunnerSpan is a span processor that is slow: OnEnd returns once release is closed.
type blocksAtTheEndOfAToolRunnerSpan struct{ entered, release chan struct{} }

func (blocksAtTheEndOfAToolRunnerSpan) OnStart(context.Context, sdktrace.ReadWriteSpan) {}
func (blocksAtTheEndOfAToolRunnerSpan) Shutdown(context.Context) error                  { return nil }
func (blocksAtTheEndOfAToolRunnerSpan) ForceFlush(context.Context) error                { return nil }
func (p blocksAtTheEndOfAToolRunnerSpan) OnEnd(s sdktrace.ReadOnlySpan) {
	if s.Name() == "anthropic.messages.tool_runner" {
		close(p.entered)
		<-p.release
	}
}

// collectUntil collects garbage until done is closed, and reports whether it was within 100 collections.
func collectUntil(done <-chan struct{}) bool {
	for range 100 {
		runtime.GC()
		select {
		case <-done:
			return true
		case <-time.After(10 * time.Millisecond):
		}
	}
	return false
}

// The cleanups of a program share goroutines of the runtime, so the runner's must not wait for a span processor.
func TestOpenTelemetryToolRunnerCleanupDoesNotWaitForASpanProcessor(t *testing.T) {
	for name, run := range map[string]func(client anthropic.Client){
		"dropped between steps": func(client anthropic.Client) {
			runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&weatherTool{}}, toolRunnerParams())
			if message, err := runner.NextMessage(context.Background()); err != nil || message == nil {
				t.Errorf("got message %v, error %v; want the turn that calls the tool", message, err)
			}
		},
		// The runner is unreachable while the loop that the application left ends the span.
		"stopped": func(client anthropic.Client) {
			for range client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&weatherTool{}}, toolRunnerParams()).All(context.Background()) {
				break
			}
		},
	} {
		t.Run(name, func(t *testing.T) {
			processor := blocksAtTheEndOfAToolRunnerSpan{entered: make(chan struct{}), release: make(chan struct{})}
			defer close(processor.release)
			client, _ := newTracedClient(t, toolRunnerTransport(), option.WithTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(processor))))
			go run(client)
			if !collectUntil(processor.entered) {
				t.Fatal("the tool runner span has not ended after 100 collections")
			}
			// The runner's cleanup starts.
			for range 5 {
				runtime.GC()
				time.Sleep(10 * time.Millisecond)
			}

			other := make(chan struct{})
			runtime.AddCleanup(new([64]byte), func(ran chan struct{}) { close(ran) }, other)
			if !collectUntil(other) {
				t.Fatal("another cleanup has not run after 100 collections")
			}
		})
	}
}

// A runner that only its iterator refers to is in use for as long as the loop over that iterator runs.
func TestOpenTelemetryToolRunnerHeldByItsIteratorAloneKeepsTheLoopSpanOpen(t *testing.T) {
	client, recorder := newTracedClient(t, toolRunnerTransport())
	messages := 0
	for _, err := range client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{&weatherTool{}}, toolRunnerParams()).All(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		messages++
		for range 5 {
			runtime.GC()
			time.Sleep(10 * time.Millisecond)
		}
		if toolRunnerSpans := startedToolRunnerSpans(recorder); len(toolRunnerSpans) != 1 || ended(toolRunnerSpans[0]) {
			t.Fatal("want the tool runner span open in the body of the loop")
		}
	}
	toolRunnerSpans := startedToolRunnerSpans(recorder)
	if messages != 2 || len(toolRunnerSpans) != 1 || toolRunnerSpans[0].Status().Code != codes.Ok {
		t.Fatalf("got %d messages and %d tool runner spans, want 2 and the one of the completed loop, with status OK", messages, len(toolRunnerSpans))
	}
}

func TestOpenTelemetryStreamingToolRunnerRecordsOutputContentAndBuildsItsMessagesItself(t *testing.T) {
	// The conversation that the runner builds from its streams, for each way that
	// the spans of its requests treat them.
	conversations := map[string]string{}
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	panicking := &panickingProvider{TracerProvider: sdktrace.NewTracerProvider(), key: "gen_ai.output.messages"}
	for name, opts := range map[string][]option.RequestOption{
		"content":       {content},
		"metadata only": nil,
		"no spans":      {option.WithoutOpenTelemetry()},
		// Spans that record nothing: the runner builds its Messages all the same.
		"content, sampled out":       {content, option.WithTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.NeverSample())))},
		"content, no-op provider":    {content, option.WithTracerProvider(noop.NewTracerProvider())},
		"content, recording panics":  {content, option.WithTracerProvider(panicking)},
		"metadata, recording panics": {option.WithTracerProvider(&panickingProvider{TracerProvider: sdktrace.NewTracerProvider(), key: "gen_ai.response.id"})},
	} {
		t.Run(name, func(t *testing.T) {
			if name == "content, sampled out" {
				// The requests carry the trace context of spans that do not record.
				setPropagator(t, propagation.TraceContext{})
			}
			client, recorder := newTracedClient(t, toolRunnerStreamingTransport(), opts...)
			runner := client.Beta.Messages.NewToolRunnerStreaming([]anthropic.BetaTool{&weatherTool{}}, toolRunnerParams())
			for events, err := range runner.AllStreaming(context.Background()) {
				if err != nil {
					t.Fatal(err)
				}
				for _, err := range events {
					if err != nil {
						t.Fatal(err)
					}
				}
			}
			if !runner.IsCompleted() || runner.IterationCount() != 2 {
				t.Fatalf("the runner made %d requests and completed: %t, want 2 and true", runner.IterationCount(), runner.IsCompleted())
			}
			conversation, err := json.Marshal(runner.Messages())
			if err != nil {
				t.Fatal(err)
			}
			conversations[name] = string(conversation)

			var outputs []any
			for _, s := range recorder.Spans() {
				if s.Name() != "anthropic.messages.create" {
					continue
				}
				if value, ok := attr(s, "gen_ai.output.messages"); ok {
					outputs = append(outputs, value.AsString())
				} else {
					outputs = append(outputs, nil)
				}
			}
			switch name {
			case "content":
				// Each request's span has the Message that its stream built.
				if len(outputs) != 2 || outputs[0] == nil || outputs[1] == nil {
					t.Fatalf("got output content %v, want it on both request spans", outputs)
				}
				for i, want := range []string{
					`[{"role":"assistant","parts":[{"type":"text","content":"Let me check."},
						{"type":"tool_call","id":"toolu_1","name":"get_weather","arguments":{"city":"Paris"}}]}]`,
					`[{"role":"assistant","parts":[{"type":"text","content":"hi"}]}]`,
				} {
					var got, wanted any
					if err := json.Unmarshal([]byte(outputs[i].(string)), &got); err != nil {
						t.Fatal(err)
					}
					if err := json.Unmarshal([]byte(want), &wanted); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(got, wanted) {
						t.Errorf("request %d: gen_ai.output.messages = %s, want %s", i, outputs[i], want)
					}
				}
			case "metadata only":
				if len(outputs) != 2 || outputs[0] != nil || outputs[1] != nil {
					t.Errorf("got output content %v, want two request spans without any", outputs)
				}
			default:
				if len(recorder.Started()) != 0 {
					t.Errorf("got %d spans, want none", len(recorder.Started()))
				}
			}
		})
	}
	// Recording the content of each of the two streams failed, and the run went on.
	if got := panicking.panics.Load(); got != 2 {
		t.Errorf("recording the output content panicked %d times, want once for each stream", got)
	}
	// The runner builds the Message of each stream itself, from the events that it
	// reads, and the span only reads it: the conversation never depends on the spans.
	want := conversations["no spans"]
	if !strings.Contains(want, `"Let me check."`) || !strings.Contains(want, `"city":"Paris"`) || !strings.Contains(want, `"tool_result"`) {
		t.Fatalf("the conversation is %s, want the streamed turn and the tool's result in it", want)
	}
	for name, got := range conversations {
		if got != want {
			t.Errorf("%s: the conversation is %s, want %s as without spans", name, got, want)
		}
	}
}

// panickingProvider is a tracer provider whose spans panic when they are given
// the attribute key, as a broken exporter or processor might.
type panickingProvider struct {
	trace.TracerProvider
	key    string
	panics atomic.Int32
}

func (p *panickingProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
	return &panickingTracer{Tracer: p.TracerProvider.Tracer(name, opts...), provider: p}
}

type panickingTracer struct {
	trace.Tracer
	provider *panickingProvider
}

func (t *panickingTracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	ctx, span := t.Tracer.Start(ctx, name, opts...)
	span = &panickingSpan{Span: span, provider: t.provider}
	return trace.ContextWithSpan(ctx, span), span
}

type panickingSpan struct {
	trace.Span
	provider *panickingProvider
}

func (s *panickingSpan) SetAttributes(attrs ...attribute.KeyValue) {
	for _, kv := range attrs {
		if string(kv.Key) == s.provider.key {
			s.provider.panics.Add(1)
			panic("the span cannot record " + s.provider.key)
		}
	}
	s.Span.SetAttributes(attrs...)
}

func TestOpenTelemetryToolRunnerStepsRunInsideTheToolRunnerSpanWhateverContextTheyAreCalledWith(t *testing.T) {
	for _, streaming := range []bool{false, true} {
		t.Run(map[bool]string{false: "messages", true: "streaming"}[streaming], func(t *testing.T) {
			recorder := newSpanRecorder()
			tracer := recorder.Tracer("app")
			// Other instrumentation, which starts a span from the context that the
			// runner's request has: a middleware that the application passed.
			instrumentation := option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
				_, span := tracer.Start(req.Context(), "http_instrumentation")
				defer span.End()
				return next(req)
			})
			rt := toolRunnerTransport()
			if streaming {
				rt = toolRunnerStreamingTransport()
			}
			client, _ := newTracedClient(t, rt, option.WithTracerProvider(recorder), instrumentation)
			tools := []anthropic.BetaTool{&tracingTool{tracer: tracer}}

			// step runs one step of a runner of either kind with ctx, to the end of the step.
			var step func(ctx context.Context)
			if streaming {
				runner := client.Beta.Messages.NewToolRunnerStreaming(tools, toolRunnerParams())
				step = func(ctx context.Context) {
					for _, err := range runner.NextStreaming(ctx) {
						if err != nil {
							t.Fatal(err)
						}
					}
				}
			} else {
				runner := client.Beta.Messages.NewToolRunner(tools, toolRunnerParams())
				step = func(ctx context.Context) {
					if _, err := runner.NextMessage(ctx); err != nil {
						t.Fatal(err)
					}
				}
			}

			first, firstRequest := tracer.Start(context.Background(), "first_request")
			step(first)
			// Between steps the application's context is as it was: it has the
			// application's span, and a span started from it is that span's child.
			if got := trace.SpanFromContext(first).SpanContext().SpanID(); got != firstRequest.SpanContext().SpanID() {
				t.Errorf("the span in the application's context is %s after a step, want its own %s", got, firstRequest.SpanContext().SpanID())
			}
			_, between := tracer.Start(first, "between_steps")
			between.End()
			// The next steps come from an unrelated context, in another trace.
			second, secondRequest := tracer.Start(context.Background(), "second_request")
			step(second)
			step(second)
			firstRequest.End()
			secondRequest.End()

			byName := map[string][]sdktrace.ReadOnlySpan{}
			for _, s := range recorder.Spans() {
				byName[s.Name()] = append(byName[s.Name()], s)
			}
			toolRunnerSpans, creates := byName["anthropic.messages.tool_runner"], byName["anthropic.messages.create"]
			uses, inTool, inRequest := byName["anthropic.custom_tool_use get_weather"], byName["look_up_forecast"], byName["http_instrumentation"]
			if len(toolRunnerSpans) != 1 || len(creates) != 2 || len(uses) != 1 || len(inTool) != 1 || len(inRequest) != 2 {
				t.Fatalf("got %d tool runner spans, %d requests, %d tool runs, %d spans from the tool and %d from the middleware; want 1, 2, 1, 1 and 2",
					len(toolRunnerSpans), len(creates), len(uses), len(inTool), len(inRequest))
			}
			toolRunnerSpan := toolRunnerSpans[0]
			if toolRunnerSpan.Parent().SpanID() != firstRequest.SpanContext().SpanID() {
				t.Error("the tool runner span is not a child of the span of the first step's context")
			}
			if byName["between_steps"][0].Parent().SpanID() != firstRequest.SpanContext().SpanID() {
				t.Error("a span that the application started between steps is not a child of its own span")
			}
			// Everything that ran inside a step is under the tool runner span, the
			// steps that were called with the unrelated context too.
			for _, s := range append(append([]sdktrace.ReadOnlySpan{}, creates...), uses...) {
				if s.Parent().SpanID() != toolRunnerSpan.SpanContext().SpanID() {
					t.Errorf("span %q is not a child of the tool runner span", s.Name())
				}
			}
			requests := map[trace.SpanID]bool{creates[0].SpanContext().SpanID(): true, creates[1].SpanContext().SpanID(): true}
			for i, s := range inRequest {
				if !requests[s.Parent().SpanID()] {
					t.Errorf("the middleware's span %d is not a child of a request span", i)
				}
			}
			if inTool[0].Parent().SpanID() != uses[0].SpanContext().SpanID() {
				t.Error("the span that the tool started is not a child of the tool's span")
			}
			for _, s := range recorder.Spans() {
				inSecond := s.SpanContext().TraceID() == secondRequest.SpanContext().TraceID()
				if s.Name() == "second_request" != inSecond {
					t.Errorf("span %q is in the trace of the unrelated context: %t", s.Name(), inSecond)
				}
			}
		})
	}
}

// runToolRunner runs a tool runner of either kind, which client makes with
// opts, to its end.
func runToolRunner(t *testing.T, client anthropic.Client, streaming bool, opts ...option.RequestOption) {
	t.Helper()
	tools := []anthropic.BetaTool{&weatherTool{}}
	if !streaming {
		if _, err := client.Beta.Messages.NewToolRunner(tools, toolRunnerParams(), opts...).RunToCompletion(context.Background()); err != nil {
			t.Fatal(err)
		}
		return
	}
	for events, err := range client.Beta.Messages.NewToolRunnerStreaming(tools, toolRunnerParams(), opts...).AllStreaming(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		for _, err := range events {
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

// Each OpenTelemetry option sets one setting, and the options of a tool runner
// apply after the client's. So a runner that sets one setting keeps the others
// as the client set them, on the tool runner span, the spans of its tools and
// the spans of its requests, and a setting that both set is the runner's.
func TestOpenTelemetryToolRunnerOptionKeepsTheOtherSettingsOfTheClient(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	for name, test := range map[string]struct {
		client, runner []option.RequestOption
		// runnerProvider gives the runner a tracer provider of its own.
		runnerProvider bool
		content        bool
		conversation   any
	}{
		"conversation ID": {
			client:       []option.RequestOption{content},
			runner:       []option.RequestOption{option.WithOpenTelemetryConversationID("conv_run")},
			content:      true,
			conversation: "conv_run",
		},
		"conversation ID over the client's": {
			client:       []option.RequestOption{content, option.WithOpenTelemetryConversationID("conv_client")},
			runner:       []option.RequestOption{option.WithOpenTelemetryConversationID("conv_run")},
			content:      true,
			conversation: "conv_run",
		},
		"content mode": {
			client:       []option.RequestOption{option.WithOpenTelemetryConversationID("conv_client")},
			runner:       []option.RequestOption{content},
			content:      true,
			conversation: "conv_client",
		},
		"content mode over the client's": {
			client:       []option.RequestOption{content, option.WithOpenTelemetryConversationID("conv_client")},
			runner:       []option.RequestOption{option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeMetadataOnly)},
			conversation: "conv_client",
		},
		"tracer provider": {
			client:         []option.RequestOption{content, option.WithOpenTelemetryConversationID("conv_client")},
			runnerProvider: true,
			content:        true,
			conversation:   "conv_client",
		},
		"none": {
			client:       []option.RequestOption{content, option.WithOpenTelemetryConversationID("conv_client")},
			content:      true,
			conversation: "conv_client",
		},
	} {
		for _, streaming := range []bool{false, true} {
			t.Run(name+map[bool]string{false: "", true: " streaming"}[streaming], func(t *testing.T) {
				registered := newSpanRecorder()
				otel.SetTracerProvider(registered)
				t.Cleanup(func() { otel.SetTracerProvider(noop.NewTracerProvider()) })

				rt := toolRunnerTransport()
				if streaming {
					rt = toolRunnerStreamingTransport()
				}
				client, ofClient := newTracedClient(t, rt, test.client...)
				opts, recorder, unused := test.runner, ofClient, newSpanRecorder()
				if test.runnerProvider {
					recorder, unused = unused, ofClient
					opts = append(opts, option.WithTracerProvider(recorder))
				}
				runToolRunner(t, client, streaming, opts...)

				if got := len(unused.Started()) + len(registered.Started()); got != 0 {
					t.Errorf("%d spans went to a tracer provider that the runner does not have", got)
				}
				byName := map[string][]sdktrace.ReadOnlySpan{}
				for _, s := range recorder.Spans() {
					byName[s.Name()] = append(byName[s.Name()], s)
					wantAttrs(t, s, map[string]any{"gen_ai.conversation.id": test.conversation})
				}
				toolRunnerSpans, uses, creates := byName["anthropic.messages.tool_runner"], byName["anthropic.custom_tool_use get_weather"], byName["anthropic.messages.create"]
				if len(toolRunnerSpans) != 1 || len(uses) != 1 || len(creates) != 2 || len(recorder.Started()) != 4 {
					t.Fatalf("got %d spans, want one tool runner span, one tool use, and two creates", len(recorder.Started()))
				}
				for _, s := range append(uses, creates...) {
					if s.Parent().SpanID() != toolRunnerSpans[0].SpanContext().SpanID() {
						t.Errorf("span %q is not a child of the tool runner span", s.Name())
					}
				}
				for key, spans := range map[string][]sdktrace.ReadOnlySpan{"gen_ai.tool.call.arguments": uses, "gen_ai.input.messages": creates} {
					for _, s := range spans {
						if _, ok := attr(s, key); ok != test.content {
							t.Errorf("%s: the attribute %s is there: %v, want %v", s.Name(), key, ok, test.content)
						}
					}
				}
			})
		}
	}
}

// A tool runner that turns the spans off does so for its own loop only: it
// leaves the client's settings as they are for the next runner.
func TestOpenTelemetryTracesOffForOneToolRunner(t *testing.T) {
	setPropagator(t, propagation.TraceContext{})
	for name, off := range map[string]option.RequestOption{
		"traces off":            option.WithOpenTelemetryTraces(false),
		"without OpenTelemetry": option.WithoutOpenTelemetry(),
	} {
		for _, streaming := range []bool{false, true} {
			t.Run(name+map[bool]string{false: "", true: " streaming"}[streaming], func(t *testing.T) {
				registered := newSpanRecorder()
				otel.SetTracerProvider(registered)
				t.Cleanup(func() { otel.SetTracerProvider(noop.NewTracerProvider()) })

				// Each runner makes two requests, so each gets a transport of its own.
				transports := []*closureTransport{toolRunnerTransport(), toolRunnerTransport()}
				if streaming {
					transports = []*closureTransport{toolRunnerStreamingTransport(), toolRunnerStreamingTransport()}
				}
				requests, traceparents := 0, 0
				client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
					requests++
					if req.Header.Get("traceparent") != "" {
						traceparents++
					}
					return transports[(requests-1)/2].fn(req)
				}}, option.WithOpenTelemetryTraces(true), option.WithOpenTelemetryConversationID("conv_client"))

				runToolRunner(t, client, streaming, off)
				if got := len(recorder.Started()) + len(registered.Started()); got != 0 || requests != 2 || traceparents != 0 {
					t.Fatalf("got %d spans and %d traceparents in %d requests from a runner that turns the spans off", got, traceparents, requests)
				}

				runToolRunner(t, client, streaming)
				if len(startedToolRunnerSpans(recorder)) != 1 || len(recorder.Started()) != 4 || len(recorder.Spans()) != 4 {
					t.Fatalf("got %d spans, %d of them ended, want the 4 of one runner", len(recorder.Started()), len(recorder.Spans()))
				}
				for _, s := range recorder.Spans() {
					wantAttrs(t, s, map[string]any{"gen_ai.conversation.id": "conv_client"})
				}
				if requests != 4 || traceparents != 2 {
					t.Errorf("%d of %d requests have a traceparent, want the 2 of the runner that has spans", traceparents, requests)
				}
				if got := len(registered.Started()); got != 0 {
					t.Errorf("%d spans went to the registered tracer provider, want the client's", got)
				}
			})
		}
	}
}

// describedTool is a weatherTool that counts the calls of its Description.
type describedTool struct {
	weatherTool
	descriptions atomic.Int32
}

func (d *describedTool) Description() string {
	d.descriptions.Add(1)
	return "Gets the weather"
}

// With the spans off, a step of the runner does none of the work of a span: the
// options of the requests are resolved once for the loop, not at each step, and
// the tool's Description, which is the application's code, is not called.
func TestOpenTelemetryToolRunnerWithTheSpansOffDoesNoSpanWork(t *testing.T) {
	for name, opts := range map[string][]option.RequestOption{
		"spans on":  nil,
		"spans off": {option.WithoutOpenTelemetry()},
	} {
		for _, streaming := range []bool{false, true} {
			t.Run(name+map[bool]string{false: ", messages", true: ", streaming"}[streaming], func(t *testing.T) {
				applied := 0
				counted := requestconfig.RequestOptionFunc(func(*requestconfig.RequestConfig) error {
					applied++
					return nil
				})
				tool := &describedTool{}
				tools := []anthropic.BetaTool{tool}
				var recorder *spanRecorder
				var defined int32
				if streaming {
					client, r := newTracedClient(t, toolRunnerStreamingTransport(), opts...)
					runner := client.Beta.Messages.NewToolRunnerStreaming(tools, toolRunnerParams(), counted)
					defined, recorder = tool.descriptions.Load(), r
					for events, err := range runner.AllStreaming(context.Background()) {
						if err != nil {
							t.Fatal(err)
						}
						for _, err := range events {
							if err != nil {
								t.Fatal(err)
							}
						}
					}
				} else {
					client, r := newTracedClient(t, toolRunnerTransport(), opts...)
					runner := client.Beta.Messages.NewToolRunner(tools, toolRunnerParams(), counted)
					defined, recorder = tool.descriptions.Load(), r
					if _, err := runner.RunToCompletion(context.Background()); err != nil {
						t.Fatal(err)
					}
				}
				if tool.calls != 1 {
					t.Fatalf("the tool ran %d times, want once", tool.calls)
				}

				// The runner's two requests, and once for the spans of the loop.
				if applied != 3 {
					t.Errorf("the runner's options were applied %d times, want 3", applied)
				}
				spans, described := len(recorder.Started()), tool.descriptions.Load()-defined
				if name == "spans off" && (spans != 0 || described != 0) {
					t.Errorf("got %d spans and %d calls of the tool's Description for them, want none", spans, described)
				}
				if name == "spans on" && (spans != 4 || described != 1) {
					t.Errorf("got %d spans and %d calls of the tool's Description for them, want 4 and 1", spans, described)
				}
			})
		}
	}
}

// A tool runs in a goroutine of the runner's, where a panic ends the process: a
// span that fails to record the tool's content must not, and the run goes on.
func TestOpenTelemetryToolRunnerSurvivesAToolSpanThatCannotRecordContent(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	provider := &panickingProvider{TracerProvider: sdktrace.NewTracerProvider(), key: "gen_ai.tool.call.arguments"}
	tool := &weatherTool{}
	client, _ := newTracedClient(t, toolRunnerTransport(), content, option.WithTracerProvider(provider))

	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{tool}, toolRunnerParams())
	message, err := runner.RunToCompletion(context.Background())
	if err != nil || message == nil || message.StopReason != anthropic.BetaStopReasonEndTurn {
		t.Fatalf("got message %v, error %v; want the final message", message, err)
	}
	if got := provider.panics.Load(); got != 1 || tool.calls != 1 {
		t.Errorf("recording the tool's content panicked %d times for %d runs of the tool, want once for one", got, tool.calls)
	}
	conversation, err := json.Marshal(runner.Messages())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(conversation), `"tool_result"`) || !strings.Contains(string(conversation), "sunny") {
		t.Errorf("the conversation is %s, want the tool's result in it", conversation)
	}
}
