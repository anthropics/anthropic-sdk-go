package anthropic_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/otel/codes"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

func TestOpenTelemetryStreamSpan(t *testing.T) {
	client, recorder := newTracedClient(t, respond(200, "text/event-stream", streamBody))

	stream := client.Messages.NewStreaming(context.Background(), messageParams())
	if len(recorder.Spans()) != 0 {
		t.Fatal("the span ended before the stream was read")
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
	if _, ok := attr(span, "gen_ai.response.time_to_first_chunk"); !ok {
		t.Error("time_to_first_chunk is absent")
	}
	wantAttrs(t, span, map[string]any{
		"gen_ai.request.stream":                                    true,
		"gen_ai.response.id":                                       "msg_1",
		"gen_ai.response.model":                                    "claude-served",
		"anthropic.message.stop_reason":                            "end_turn",
		"gen_ai.usage.input_tokens":                                int64(15),
		"gen_ai.usage.output_tokens":                               int64(7),
		"gen_ai.usage.reasoning.output_tokens":                     int64(4),
		"anthropic.usage.cache_creation.ephemeral_5m_input_tokens": int64(3),
		"anthropic.response.service_tier":                          "standard",
	})
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}
	if span.Status().Code != codes.Ok {
		t.Errorf("closing a finished stream changed the status to %v", span.Status().Code)
	}
}

func TestOpenTelemetryStreamClosedEarly(t *testing.T) {
	client, recorder := newTracedClient(t, respond(200, "text/event-stream", streamBody))

	stream := client.Messages.NewStreaming(context.Background(), messageParams())
	if !stream.Next() {
		t.Fatal(stream.Err())
	}
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}

	span := onlySpan(t, recorder)
	if span.Status().Code != codes.Unset {
		t.Errorf("status %v, want unset", span.Status().Code)
	}
	wantAttrs(t, span, map[string]any{
		"gen_ai.response.id":              "msg_1",
		"gen_ai.usage.input_tokens":       int64(15),
		"anthropic.response.service_tier": "standard",
		"gen_ai.usage.output_tokens":      nil,
		"anthropic.message.stop_reason":   nil,
		"gen_ai.response.finish_reasons":  nil,
		"error.type":                      nil,
	})
}

// messageStreams reads the Message stream of each Messages API: it passes seen,
// when it is not nil, the type of each event, until the stream ends or seen
// returns false, and then closes the stream.
var messageStreams = map[string]func(client anthropic.Client, seen func(eventType string) bool) error{
	"Messages.NewStreaming": func(client anthropic.Client, seen func(string) bool) error {
		stream := client.Messages.NewStreaming(context.Background(), messageParams())
		for stream.Next() {
			if seen != nil && !seen(stream.Current().Type) {
				break
			}
		}
		if err := stream.Err(); err != nil {
			return err
		}
		return stream.Close()
	},
	"Beta.Messages.NewStreaming": func(client anthropic.Client, seen func(string) bool) error {
		stream := client.Beta.Messages.NewStreaming(context.Background(), anthropic.BetaMessageNewParams{Model: "claude-requested", MaxTokens: 100})
		for stream.Next() {
			if seen != nil && !seen(stream.Current().Type) {
				break
			}
		}
		if err := stream.Err(); err != nil {
			return err
		}
		return stream.Close()
	},
}

func TestOpenTelemetryStreamClosedAfterStop(t *testing.T) {
	// The application stops at message_stop, before the server ends the body: it cancels the call.
	for name, read := range messageStreams {
		t.Run(name, func(t *testing.T) {
			client, recorder := newTracedClient(t, respond(200, "text/event-stream", streamBody))
			if err := read(client, func(eventType string) bool { return eventType != "message_stop" }); err != nil {
				t.Fatal(err)
			}
			span := onlySpan(t, recorder)
			if span.Status().Code != codes.Unset {
				t.Errorf("status %v, want unset", span.Status().Code)
			}
			// The whole Message was read, so the span has all of it.
			wantAttrs(t, span, map[string]any{
				"gen_ai.response.id":            "msg_1",
				"anthropic.message.stop_reason": "end_turn",
				"gen_ai.usage.output_tokens":    int64(7),
				"error.type":                    nil,
			})
		})
	}
}

func TestOpenTelemetryStreamEndedBeforeStop(t *testing.T) {
	// The server ends the body before message_stop: the span fails, and has the
	// response attributes of the events that arrived.
	beforeDelta, _, _ := strings.Cut(streamBody, "event: message_delta\n")
	beforeStop, _, _ := strings.Cut(streamBody, "event: message_stop\n")
	tests := map[string]struct {
		body string
		want map[string]any
	}{
		"before message_delta": {beforeDelta, map[string]any{
			"anthropic.message.stop_reason":        nil,
			"gen_ai.usage.output_tokens":           nil,
			"gen_ai.usage.reasoning.output_tokens": nil,
		}},
		"after message_delta": {beforeStop, map[string]any{
			"anthropic.message.stop_reason":        "end_turn",
			"gen_ai.usage.output_tokens":           int64(7),
			"gen_ai.usage.reasoning.output_tokens": int64(4),
		}},
	}
	for name, read := range messageStreams {
		for when, test := range tests {
			t.Run(name+" "+when, func(t *testing.T) {
				client, recorder := newTracedClient(t, respond(200, "text/event-stream", test.body))
				if err := read(client, nil); err != nil {
					t.Fatal(err)
				}

				span := onlySpan(t, recorder)
				if span.Status().Code != codes.Error {
					t.Errorf("status %v, want ERROR", span.Status().Code)
				}
				wantAttrs(t, span, test.want)
				wantAttrs(t, span, map[string]any{
					"gen_ai.response.id":              "msg_1",
					"gen_ai.response.model":           "claude-served",
					"gen_ai.usage.input_tokens":       int64(15),
					"anthropic.response.service_tier": "standard",
					"gen_ai.response.finish_reasons":  nil,
					"error.type":                      "api_error",
				})
			})
		}
	}
}

func TestOpenTelemetryStreamCompactionTokens(t *testing.T) {
	// The message_delta event has the usage of each step that made the Message.
	body := strings.Replace(streamBody, `"usage":{"output_tokens":7,`, `"usage":{"output_tokens":7,"iterations":[`+
		`{"type":"compaction","input_tokens":100,"cache_read_input_tokens":20,"cache_creation_input_tokens":30,"output_tokens":40},`+
		`{"type":"message","input_tokens":10,"cache_read_input_tokens":2,"cache_creation_input_tokens":3,"output_tokens":7}],`, 1)
	beforeStop, _, _ := strings.Cut(body, "event: message_stop\n")
	beforeDelta, _, _ := strings.Cut(body, "event: message_delta\n")
	streams := map[string]func(client anthropic.Client) error{
		"Messages.NewStreaming": func(client anthropic.Client) error {
			stream := client.Messages.NewStreaming(context.Background(), messageParams())
			for stream.Next() {
			}
			return stream.Err()
		},
		"Beta.Messages.NewStreaming": func(client anthropic.Client) error {
			stream := client.Beta.Messages.NewStreaming(context.Background(), anthropic.BetaMessageNewParams{Model: "claude-requested", MaxTokens: 100})
			for stream.Next() {
			}
			return stream.Err()
		},
	}
	modes := []option.OpenTelemetryContentMode{option.OpenTelemetryContentModeMetadataOnly, option.OpenTelemetryContentModeContent}
	for name, read := range streams {
		for _, mode := range modes {
			t.Run(name+" "+string(mode), func(t *testing.T) {
				client, recorder := newTracedClient(t, respond(200, "text/event-stream", body), option.WithOpenTelemetryTracesContentMode(mode))
				if err := read(client); err != nil {
					t.Fatal(err)
				}
				wantAttrs(t, onlySpan(t, recorder), map[string]any{
					"gen_ai.usage.input_tokens":                int64(15),
					"gen_ai.usage.output_tokens":               int64(7),
					"anthropic.usage.compaction.input_tokens":  int64(150),
					"anthropic.usage.compaction.output_tokens": int64(40),
				})
			})
			// A stream that ends after message_delta keeps what that event brought.
			t.Run(name+" "+string(mode)+" ended before message_stop", func(t *testing.T) {
				client, recorder := newTracedClient(t, respond(200, "text/event-stream", beforeStop), option.WithOpenTelemetryTracesContentMode(mode))
				if err := read(client); err != nil {
					t.Fatal(err)
				}
				wantAttrs(t, onlySpan(t, recorder), map[string]any{
					"error.type":                               "api_error",
					"gen_ai.usage.input_tokens":                int64(15),
					"gen_ai.usage.output_tokens":               int64(7),
					"anthropic.usage.compaction.input_tokens":  int64(150),
					"anthropic.usage.compaction.output_tokens": int64(40),
				})
			})
			// As the output tokens, they are absent before message_delta.
			t.Run(name+" "+string(mode)+" ended before message_delta", func(t *testing.T) {
				client, recorder := newTracedClient(t, respond(200, "text/event-stream", beforeDelta), option.WithOpenTelemetryTracesContentMode(mode))
				if err := read(client); err != nil {
					t.Fatal(err)
				}
				wantAttrs(t, onlySpan(t, recorder), map[string]any{
					"error.type":                               "api_error",
					"gen_ai.usage.input_tokens":                int64(15),
					"gen_ai.usage.output_tokens":               nil,
					"anthropic.usage.compaction.input_tokens":  nil,
					"anthropic.usage.compaction.output_tokens": nil,
				})
			})
		}
	}
}

const sessionEventsBody = "event: session.status_running\n" +
	`data: {"type":"session.status_running","id":"evt_1"}` + "\n\n" +
	"event: session.status_idle\n" +
	`data: {"type":"session.status_idle","id":"evt_2"}` + "\n\n"

func TestOpenTelemetryEventStreamClosed(t *testing.T) {
	client, recorder := newTracedClient(t, respond(200, "text/event-stream", sessionEventsBody))

	stream := client.Beta.Sessions.Events.StreamEvents(context.Background(), "sess_1", anthropic.BetaSessionEventStreamParams{})
	if !stream.Next() {
		t.Fatal(stream.Err())
	}
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}

	// An event stream has no last event, so closing it cancels the call.
	span := onlySpan(t, recorder)
	if span.Name() != "anthropic.sessions.events.stream" || span.Status().Code != codes.Unset {
		t.Errorf("span %q with status %v, want anthropic.sessions.events.stream with status unset", span.Name(), span.Status().Code)
	}
	wantAttrs(t, span, map[string]any{"error.type": nil})
}

func TestOpenTelemetryEventStreamEndedByServer(t *testing.T) {
	client, recorder := newTracedClient(t, respond(200, "text/event-stream", sessionEventsBody))

	stream := client.Beta.Sessions.Events.StreamEvents(context.Background(), "sess_1", anthropic.BetaSessionEventStreamParams{})
	for stream.Next() {
	}
	if err := stream.Err(); err != nil {
		t.Fatal(err)
	}

	// An event stream has no last event, so it is never read to the end: the server ending it sets no status.
	span := onlySpan(t, recorder)
	if span.Name() != "anthropic.sessions.events.stream" || span.Status().Code != codes.Unset {
		t.Errorf("span %q with status %v, want anthropic.sessions.events.stream with status unset", span.Name(), span.Status().Code)
	}
	if span.EndTime().IsZero() {
		t.Error("the span has not ended")
	}
	wantAttrs(t, span, map[string]any{
		"error.type":                nil,
		"http.response.status_code": int64(200),
		"gen_ai.provider.name":      "anthropic",
	})
}

func TestOpenTelemetryStreamErrorEvent(t *testing.T) {
	body := "event: error\n" + `data: {"type":"error","error":{"type":"overloaded_error","message":"Overloaded"}}` + "\n\n"
	client, recorder := newTracedClient(t, respond(200, "text/event-stream", body))

	stream := client.Messages.NewStreaming(context.Background(), messageParams())
	for stream.Next() {
	}
	if stream.Err() == nil {
		t.Fatal("the stream did not fail")
	}

	span := onlySpan(t, recorder)
	if span.Status().Code != codes.Error {
		t.Errorf("status %v, want ERROR", span.Status().Code)
	}
	wantAttrs(t, span, map[string]any{"error.type": "overloaded_error"})
}

func TestOpenTelemetryStreamClosedBeforeReading(t *testing.T) {
	client, recorder := newTracedClient(t, respond(200, "text/event-stream", streamBody))

	stream := client.Messages.NewStreaming(context.Background(), messageParams())
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}
	_ = stream.Close()
	span := onlySpan(t, recorder)
	if span.Status().Code != codes.Unset {
		t.Errorf("status %v, want unset", span.Status().Code)
	}
	wantAttrs(t, span, map[string]any{"gen_ai.request.model": "claude-requested", "error.type": nil, "gen_ai.response.id": nil})
}

// blockingBody is a response body whose reads wait for the request's context to end.
type blockingBody struct{ ctx context.Context }

func (b blockingBody) Read([]byte) (int, error) {
	<-b.ctx.Done()
	return 0, b.ctx.Err()
}

func (blockingBody) Close() error { return nil }

func TestOpenTelemetryStreamCanceledBeforeReading(t *testing.T) {
	client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"text/event-stream"}}, Body: blockingBody{req.Context()}, Request: req}, nil
	}})

	ctx, cancel := context.WithCancel(context.Background())
	stream := client.Messages.NewStreaming(ctx, messageParams())
	cancel()
	if stream.Next() {
		t.Fatal("a canceled stream delivered an event")
	}
	if !errors.Is(stream.Err(), context.Canceled) {
		t.Fatalf("error %v, want context.Canceled", stream.Err())
	}
	span := onlySpan(t, recorder)
	if span.Status().Code != codes.Unset {
		t.Errorf("status %v, want unset", span.Status().Code)
	}
	wantAttrs(t, span, map[string]any{"error.type": nil})
	_ = stream.Close()
	if len(recorder.Spans()) != 1 {
		t.Errorf("got %d ended spans after closing, want 1", len(recorder.Spans()))
	}
}

func TestOpenTelemetryStreamRetriedMeasuresTimeToFirstChunkFromTheStart(t *testing.T) {
	attempts := 0
	client, recorder := newTracedClient(t, &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		attempts++
		if attempts == 1 {
			return respond(529, "application/json", `{"type":"error","error":{"type":"overloaded_error","message":"Overloaded"}}`).fn(req)
		}
		return respond(200, "text/event-stream", streamBody).fn(req)
	}}, option.WithMaxRetries(1))

	start := time.Now()
	stream := client.Messages.NewStreaming(context.Background(), messageParams())
	for stream.Next() {
	}
	if err := stream.Err(); err != nil || attempts != 2 {
		t.Fatalf("error %v after %d attempts", err, attempts)
	}
	elapsed := time.Since(start)
	span := onlySpan(t, recorder)
	ttfc, ok := attr(span, "gen_ai.response.time_to_first_chunk")
	if !ok {
		t.Fatal("time_to_first_chunk is absent")
	}
	// The retry waits at least 0.375s: the first chunk is timed from before the first attempt.
	if got := ttfc.AsFloat64(); got < 0.3 || got > elapsed.Seconds() {
		t.Errorf("time_to_first_chunk = %v, want between the retry delay and the call's %v", got, elapsed)
	}
}

// respondInParts returns a transport that answers every request with an event
// stream that serves parts one at a time, waiting for pause before each part after the first.
func respondInParts(pause time.Duration, parts ...string) *closureTransport {
	return &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		body, server := io.Pipe()
		go func() {
			for i, part := range parts {
				if i > 0 {
					time.Sleep(pause)
				}
				if _, err := io.WriteString(server, part); err != nil {
					return
				}
			}
			server.Close()
		}()
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"text/event-stream"}}, Body: body, Request: req}, nil
	}}
}

func TestOpenTelemetryStreamTimeToFirstChunkIsTheFirstContentDelta(t *testing.T) {
	const pause = 50 * time.Millisecond
	const delta = "event: content_block_delta\n"
	// message_start, content_block_start and ping arrive at once, then a delta, then another with the rest.
	beforeDeltas, rest, _ := strings.Cut(streamBody, delta)
	end := strings.Index(rest, "\n\n") + 2
	firstDelta, secondDeltaAndRest := delta+rest[:end], delta+rest

	for name, read := range messageStreams {
		t.Run(name, func(t *testing.T) {
			client, recorder := newTracedClient(t, respondInParts(pause, beforeDeltas, firstDelta, secondDeltaAndRest))

			start := time.Now()
			var untilFirstDelta time.Duration
			deltas := 0
			err := read(client, func(eventType string) bool {
				if eventType == "content_block_delta" {
					if deltas++; deltas == 1 {
						untilFirstDelta = time.Since(start)
					}
				}
				return true
			})
			if err != nil || deltas != 2 {
				t.Fatalf("error %v after %d deltas, want 2", err, deltas)
			}
			ttfc, ok := attr(onlySpan(t, recorder), "gen_ai.response.time_to_first_chunk")
			if !ok {
				t.Fatal("time_to_first_chunk is absent")
			}
			// The events before the first delta do not stop the clock, and the second delta does not move it.
			if got := ttfc.AsFloat64(); got < pause.Seconds()/2 || got > untilFirstDelta.Seconds() {
				t.Errorf("time_to_first_chunk = %v, want between the pause of %v before the first delta and the %v until the application had it", got, pause, untilFirstDelta)
			}
		})
	}
}

func TestOpenTelemetryStreamWithoutContentDeltaHasNoTimeToFirstChunk(t *testing.T) {
	// A block can arrive whole in its content_block_start event.
	beforeDelta, rest, _ := strings.Cut(streamBody, "event: content_block_delta\n")
	_, afterDelta, _ := strings.Cut(rest, "\n\n")
	client, recorder := newTracedClient(t, respond(200, "text/event-stream", beforeDelta+afterDelta))

	stream := client.Messages.NewStreaming(context.Background(), messageParams())
	for stream.Next() {
		if stream.Current().Type == "content_block_delta" {
			t.Fatal("the stream has a content_block_delta event")
		}
	}
	if err := stream.Err(); err != nil {
		t.Fatal(err)
	}
	span := onlySpan(t, recorder)
	if span.Status().Code != codes.Ok {
		t.Errorf("status %v, want OK", span.Status().Code)
	}
	wantAttrs(t, span, map[string]any{
		"gen_ai.response.time_to_first_chunk": nil,
		"gen_ai.response.id":                  "msg_1",
		"gen_ai.usage.output_tokens":          int64(7),
	})
}

func TestOpenTelemetryStreamTimeToFirstChunkCountsEveryDeltaType(t *testing.T) {
	text := `{"type":"text","text":""}`
	thinking := `{"type":"thinking","thinking":"","signature":""}`
	for name, block := range map[string]struct{ start, delta string }{
		"text_delta":       {text, `{"type":"text_delta","text":"hi"}`},
		"input_json_delta": {`{"type":"tool_use","id":"toolu_1","name":"get_weather","input":{}}`, `{"type":"input_json_delta","partial_json":"{}"}`},
		"thinking_delta":   {thinking, `{"type":"thinking_delta","thinking":"Hm."}`},
		"signature_delta":  {thinking, `{"type":"signature_delta","signature":"sig"}`},
		"citations_delta": {`{"type":"text","text":"","citations":[]}`, `{"type":"citations_delta","citation":{"type":"char_location",` +
			`"cited_text":"hi","document_index":0,"document_title":null,"start_char_index":0,"end_char_index":2,"file_id":null}}`},
		// A type that this version of the SDK does not know.
		"unknown": {text, `{"type":"later_delta"}`},
	} {
		body := strings.Replace(streamBody, `"content_block":`+text, `"content_block":`+block.start, 1)
		body = strings.Replace(body, `"delta":{"type":"text_delta","text":"hi"}`, `"delta":`+block.delta, 1)
		if !strings.Contains(body, block.start) || !strings.Contains(body, block.delta) {
			t.Fatalf("%s: the stream does not have the block and its delta", name)
		}
		for api, read := range messageStreams {
			t.Run(api+" "+name, func(t *testing.T) {
				client, recorder := newTracedClient(t, respond(200, "text/event-stream", body))
				if err := read(client, nil); err != nil {
					t.Fatal(err)
				}
				if _, ok := attr(onlySpan(t, recorder), "gen_ai.response.time_to_first_chunk"); !ok {
					t.Error("time_to_first_chunk is absent")
				}
			})
		}
	}
}
