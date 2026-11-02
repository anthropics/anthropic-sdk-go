package anthropic_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// threadFields records the type of the thread field of each request that next
// answers, or "" for a request without one.
type threadFields struct {
	next  http.RoundTripper
	types []string
}

func (f *threadFields) RoundTrip(req *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewReader(body))
	f.types = append(f.types, gjson.GetBytes(body, "thread.type").String())
	return f.next.RoundTrip(req)
}

func toolRunnerThreadParams() anthropic.BetaMessageThreadCreateParams {
	params := toolRunnerParams()
	return anthropic.BetaMessageThreadCreateParams{
		Model:     params.Model,
		MaxTokens: params.MaxTokens,
		Messages:  params.Messages,
	}
}

// runThreadToolRunner runs a tool runner with tool to completion. The runner comes
// from thread when there is one, and from client with params otherwise.
func runThreadToolRunner(t *testing.T, ctx context.Context, streaming bool, client anthropic.Client, thread *anthropic.BetaMessageThread, params anthropic.BetaToolRunnerParams, tool anthropic.BetaTool) {
	t.Helper()
	tools := []anthropic.BetaTool{tool}
	if !streaming {
		runner := client.Beta.Messages.NewToolRunner(tools, params)
		if thread != nil {
			runner = thread.NewToolRunner(tools)
		}
		if _, err := runner.RunToCompletion(ctx); err != nil {
			t.Fatal(err)
		}
		return
	}
	runner := client.Beta.Messages.NewToolRunnerStreaming(tools, params)
	if thread != nil {
		runner = thread.NewToolRunnerStreaming(tools)
	}
	for events, err := range runner.AllStreaming(ctx) {
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

func toolRunnerTransportFor(streaming bool) *closureTransport {
	if streaming {
		return toolRunnerStreamingTransport()
	}
	return toolRunnerTransport()
}

// wantToolRunnerSpans checks that recorder has the spans of one finished loop of
// two requests and one tool call, with the requests and the tool call inside the
// tool runner span. It returns the tool runner span.
func wantToolRunnerSpans(t *testing.T, recorder *spanRecorder) sdktrace.ReadOnlySpan {
	t.Helper()
	var names []string
	byName := map[string][]sdktrace.ReadOnlySpan{}
	for _, s := range recorder.Spans() {
		names = append(names, s.Name())
		byName[s.Name()] = append(byName[s.Name()], s)
	}
	loops := byName["anthropic.messages.tool_runner"]
	uses := byName["anthropic.custom_tool_use get_weather"]
	creates := byName["anthropic.messages.create"]
	if len(names) != 4 || len(loops) != 1 || len(uses) != 1 || len(creates) != 2 {
		t.Fatalf("got spans %q, want one tool runner span, one tool use, and two creates", names)
	}
	if loops[0].Status().Code != codes.Ok {
		t.Errorf("tool runner span status %v, want OK", loops[0].Status().Code)
	}
	for _, child := range append(uses, creates...) {
		if child.Parent().SpanID() != loops[0].SpanContext().SpanID() {
			t.Errorf("span %q is not a child of the tool runner span", child.Name())
		}
	}
	return loops[0]
}

// The requests of a tool runner are inside its span whichever way it sends them.
func TestOpenTelemetryToolRunnerRequestsAreInsideItsSpanWithAndWithoutAMessageThread(t *testing.T) {
	unthreaded := toolRunnerParams()
	unthreaded.DisableThread = true

	for _, tt := range []struct {
		name       string
		params     anthropic.BetaToolRunnerParams
		overThread bool
		want       []string
	}{
		{name: "on a thread of its own", params: toolRunnerParams(), want: []string{"create", "continue"}},
		{name: "without a thread", params: unthreaded, want: []string{"", ""}},
		{name: "over a thread", overThread: true, want: []string{"create", "continue"}},
	} {
		for _, streaming := range []bool{false, true} {
			t.Run(tt.name+map[bool]string{false: ", messages", true: ", streaming"}[streaming], func(t *testing.T) {
				sent := &threadFields{next: toolRunnerTransportFor(streaming)}
				client, recorder := newTracedClient(t, sent)
				var thread *anthropic.BetaMessageThread
				if tt.overThread {
					thread = client.Beta.Messages.NewThread(toolRunnerThreadParams())
				}

				ctx, app := recorder.Tracer("app").Start(context.Background(), "handle_order")
				runThreadToolRunner(t, ctx, streaming, client, thread, tt.params, &weatherTool{})

				if len(sent.types) != 2 || sent.types[0] != tt.want[0] || sent.types[1] != tt.want[1] {
					t.Fatalf("the requests had the thread types %q, want %q", sent.types, tt.want)
				}
				if loop := wantToolRunnerSpans(t, recorder); loop.Parent().SpanID() != app.SpanContext().SpanID() {
					t.Error("the tool runner span is not a child of the application's span")
				}
			})
		}
	}
}

// The thread's options apply to the requests of a runner over it, so they
// configure the runner's own spans too.
func TestOpenTelemetryToolRunnerOverAMessageThreadTakesTheOptionsOfTheThread(t *testing.T) {
	for _, streaming := range []bool{false, true} {
		name := map[bool]string{false: "messages", true: "streaming"}[streaming]

		t.Run("tracer provider, "+name, func(t *testing.T) {
			client, ofClient := newTracedClient(t, toolRunnerTransportFor(streaming))
			ofThread := newSpanRecorder()
			thread := client.Beta.Messages.NewThread(toolRunnerThreadParams(), option.WithTracerProvider(ofThread))

			runThreadToolRunner(t, context.Background(), streaming, client, thread, anthropic.BetaToolRunnerParams{}, &weatherTool{})

			wantToolRunnerSpans(t, ofThread)
			if spans := ofClient.Started(); len(spans) != 0 {
				t.Errorf("the client's tracer provider started %d spans, want none", len(spans))
			}
		})

		t.Run("spans off, "+name, func(t *testing.T) {
			client, recorder := newTracedClient(t, toolRunnerTransportFor(streaming))
			thread := client.Beta.Messages.NewThread(toolRunnerThreadParams(), option.WithoutOpenTelemetry())
			tool := &weatherTool{}

			runThreadToolRunner(t, context.Background(), streaming, client, thread, anthropic.BetaToolRunnerParams{}, tool)

			if tool.calls != 1 {
				t.Fatalf("the tool ran %d times, want once", tool.calls)
			}
			if spans := recorder.Started(); len(spans) != 0 {
				t.Errorf("got %d spans, want none", len(spans))
			}
		})
	}
}

func TestOpenTelemetryToolRunnerOverAMessageThreadDroppedBetweenStepsEndsTheLoopSpan(t *testing.T) {
	ctx := context.Background()

	t.Run("messages", func(t *testing.T) {
		client, recorder := newTracedClient(t, toolRunnerTransport())
		thread := client.Beta.Messages.NewThread(toolRunnerThreadParams())
		func() {
			runner := thread.NewToolRunner([]anthropic.BetaTool{&weatherTool{}})
			if message, err := runner.NextMessage(ctx); err != nil || message == nil {
				t.Fatalf("got message %v, error %v; want the turn that calls the tool", message, err)
			}
		}()

		waitForToolRunnerSpan(t, recorder)
		wantStoppedToolRunnerSpan(t, recorder)
		// The application still holds the thread, which must not keep the runner alive.
		if thread.LastMessage() == nil {
			t.Error("the thread has no reply, want the turn that calls the tool")
		}
	})

	t.Run("streaming", func(t *testing.T) {
		client, recorder := newTracedClient(t, toolRunnerStreamingTransport())
		thread := client.Beta.Messages.NewThread(toolRunnerThreadParams())
		func() {
			runner := thread.NewToolRunnerStreaming([]anthropic.BetaTool{&weatherTool{}})
			for _, err := range runner.NextStreaming(ctx) {
				if err != nil {
					t.Fatal(err)
				}
			}
		}()

		waitForToolRunnerSpan(t, recorder)
		wantStoppedToolRunnerSpan(t, recorder)
		if thread.LastMessage() == nil {
			t.Error("the thread has no reply, want the turn that calls the tool")
		}
	})
}
