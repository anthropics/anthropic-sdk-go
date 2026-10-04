package ssestream_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
)

type fragmentedBody struct {
	io.ReadCloser
	limit int
}

func (b *fragmentedBody) Read(p []byte) (int, error) {
	if len(p) > b.limit {
		p = p[:b.limit]
	}
	return b.ReadCloser.Read(p)
}

func TestDecoderLeadingBOMAcrossReadBoundaries(t *testing.T) {
	const data = `{"type":"message_stop","text":"preserved \uFEFF marker"}`
	for _, ending := range []string{"\n", "\r\n"} {
		for _, limit := range []int{1, 2, 3, 64} {
			for _, first := range []string{"event: message_stop", "data: " + data} {
				t.Run(fmt.Sprintf("ending=%q/limit=%d/first=%s", ending, limit, first[:5]), func(t *testing.T) {
					second := "data: " + data
					if strings.HasPrefix(first, "data:") {
						second = "event: message_stop"
					}
					wire := "\ufeff" + first + ending + second + ending + ending
					response := &http.Response{Body: &fragmentedBody{ReadCloser: io.NopCloser(strings.NewReader(wire)), limit: limit}}
					decoder := ssestream.NewDecoder(response)
					defer decoder.Close()
					if !decoder.Next() {
						t.Fatalf("missing first event: %v", decoder.Err())
					}
					event := decoder.Event()
					if event.Type != "message_stop" || strings.TrimSuffix(string(event.Data), "\n") != data {
						t.Fatalf("incorrect first event: %#v", event)
					}
					if decoder.Next() || decoder.Err() != nil {
						t.Fatalf("unexpected trailing event or error: %v", decoder.Err())
					}
				})
			}
		}
	}
}

func TestDecoderOnlyStripsOneBOMAtStartOfStream(t *testing.T) {
	for _, wire := range []string{
		"\ufeff\ufeffevent: message_stop\ndata: {}\n\n",
		"\nevent: ping\ndata: {}\n\n\ufeffevent: message_stop\ndata: {}\n\n",
	} {
		decoder := ssestream.NewDecoder(&http.Response{Body: io.NopCloser(strings.NewReader(wire))})
		for decoder.Next() {
			if decoder.Event().Type == "message_stop" {
				t.Fatal("stripped a BOM that was not the first stream character")
			}
		}
		if decoder.Err() != nil {
			t.Fatal(decoder.Err())
		}
		if err := decoder.Close(); err != nil {
			t.Fatal(err)
		}
	}
	decoder := ssestream.NewDecoder(&http.Response{Body: io.NopCloser(strings.NewReader("\ufeff: comment\nevent: message_stop\ndata: {}\n\n"))})
	defer decoder.Close()
	if !decoder.Next() || decoder.Event().Type != "message_stop" {
		t.Fatal("leading BOM comment hid the following event")
	}
}

func TestMessageClientsKeepFirstBOMPrefixedEventOrError(t *testing.T) {
	for _, beta := range []bool{false, true} {
		for _, errorEvent := range []bool{false, true} {
			t.Run(fmt.Sprintf("beta=%t/error=%t", beta, errorEvent), func(t *testing.T) {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method != http.MethodPost || r.URL.Path != "/v1/messages" {
						t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
					}
					w.Header().Set("Content-Type", "text/event-stream")
					w.Header().Set("request-id", "req_bom")
					w.WriteHeader(http.StatusOK)
					frames := []string{
						"event: message_start\ndata: {\"type\":\"message_start\",\"message\":{\"id\":\"msg_bom\",\"role\":\"assistant\",\"model\":\"test\",\"type\":\"message\",\"content\":[],\"stop_reason\":null,\"stop_sequence\":null,\"usage\":{\"input_tokens\":1,\"output_tokens\":0}}}\n\n",
						"event: message_stop\ndata: {\"type\":\"message_stop\"}\n\n",
					}
					if errorEvent {
						frames = []string{"event: error\ndata: {\"type\":\"error\",\"error\":{\"type\":\"api_error\",\"message\":\"first failure\"}}\n\n"}
					}
					for _, b := range []byte("\ufeff" + strings.Join(frames, "")) {
						_, _ = w.Write([]byte{b})
						w.(http.Flusher).Flush()
					}
				}))
				defer server.Close()
				client := anthropic.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test-key"), option.WithMaxRetries(0))
				var kinds []string
				var streamError error
				if beta {
					stream := client.Beta.Messages.NewStreaming(context.Background(), anthropic.BetaMessageNewParams{Model: "test", MaxTokens: 4, Messages: []anthropic.BetaMessageParam{anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("hello"))}})
					defer stream.Close()
					for stream.Next() {
						kinds = append(kinds, stream.Current().Type)
					}
					streamError = stream.Err()
				} else {
					stream := client.Messages.NewStreaming(context.Background(), anthropic.MessageNewParams{Model: "test", MaxTokens: 4, Messages: []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hello"))}})
					defer stream.Close()
					for stream.Next() {
						kinds = append(kinds, stream.Current().Type)
					}
					streamError = stream.Err()
				}
				if errorEvent {
					var apiError *anthropic.Error
					if !errors.As(streamError, &apiError) {
						t.Fatalf("first error event was lost: %v", streamError)
					}
					if apiError.RequestID != "req_bom" || !strings.Contains(apiError.Error(), "first failure") {
						t.Fatalf("error metadata lost: %+v", apiError)
					}
					if len(kinds) != 0 {
						t.Fatalf("error stream yielded regular events: %v", kinds)
					}
				} else {
					if streamError != nil {
						t.Fatal(streamError)
					}
					if strings.Join(kinds, ",") != "message_start,message_stop" {
						t.Fatalf("first event lost: %v", kinds)
					}
				}
			})
		}
	}
}

func TestDecoderUnmarkedEventsAndInteriorBOMAreUnchanged(t *testing.T) {
	data := "{\"type\":\"message_stop\",\"text\":\"interior \ufeff mark\"}"
	wire := "event: message_stop\ndata: " + data + "\n\n"
	for _, prefix := range []string{"", "\ufeff"} {
		decoder := ssestream.NewDecoder(&http.Response{Body: io.NopCloser(strings.NewReader(prefix + wire + wire))})
		for index := 0; index < 2; index++ {
			if !decoder.Next() {
				t.Fatalf("missing event %d: %v", index, decoder.Err())
			}
			event := decoder.Event()
			if event.Type != "message_stop" || string(event.Data) != data+"\n" {
				t.Fatalf("payload changed: %#v", event)
			}
		}
		if decoder.Next() || decoder.Err() != nil {
			t.Fatal("unexpected trailing result")
		}
		if err := decoder.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
