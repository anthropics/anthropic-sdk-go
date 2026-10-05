package ssestream

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type chunkedSSEReader struct {
	data     []byte
	limit    int
	terminal error
	reads    int
	closes   int
}

func (r *chunkedSSEReader) Read(p []byte) (int, error) {
	r.reads++
	if len(r.data) == 0 {
		return 0, r.terminal
	}
	size := min(len(p), len(r.data), r.limit)
	copy(p, r.data[:size])
	r.data = r.data[size:]
	return size, nil
}
func (r *chunkedSSEReader) Close() error { r.closes++; return nil }

func TestSSELineEndingVariants(t *testing.T) {
	for _, ending := range []string{"\n", "\r\n", "\r"} {
		for _, limit := range []int{1, 2, 3, 17, 4096} {
			t.Run(strings.ReplaceAll(strings.ReplaceAll(ending, "\r", "CR"), "\n", "LF")+"/chunk-"+strconv.Itoa(limit), func(t *testing.T) {
				input := strings.Join([]string{
					": keepalive", "event: message_start", `data: {"value":`, `data: "café"}`, "",
					"event: ping", "data: {}", "", "event: message_stop", `data: {"value":"done"}`, "", "",
				}, ending)
				body := &chunkedSSEReader{data: []byte(input), limit: limit, terminal: io.EOF}
				stream := NewStream[map[string]string](NewDecoder(&http.Response{Body: body, Header: http.Header{}}), nil)
				var values []string
				for stream.Next() {
					values = append(values, stream.Current()["value"])
				}
				if stream.Err() != nil {
					t.Fatal(stream.Err())
				}
				if !reflect.DeepEqual(values, []string{"café", "done"}) {
					t.Fatalf("events=%q", values)
				}
				if err := stream.Close(); err != nil || body.closes != 1 {
					t.Fatalf("Close: %v count=%d", err, body.closes)
				}
			})
		}
	}
}

func TestSSEMixedLineEndingsAndPartialFinalEvent(t *testing.T) {
	input := "event: message_start\r\ndata: {\"value\":\"first\"}\n\r" +
		"event: message_stop\ndata: {\"value\":\"second\"}\r\n\n" +
		"event: message_stop\rdata: {\"value\":\"incomplete\"}"
	for _, limit := range []int{1, 2, 7, len(input)} {
		body := &chunkedSSEReader{data: []byte(input), limit: limit, terminal: io.EOF}
		stream := NewStream[map[string]string](NewDecoder(&http.Response{Body: body}), nil)
		var values []string
		for stream.Next() {
			values = append(values, stream.Current()["value"])
		}
		if stream.Err() != nil {
			t.Fatal(stream.Err())
		}
		if !reflect.DeepEqual(values, []string{"first", "second"}) {
			t.Fatalf("chunk=%d events=%q", limit, values)
		}
		stream.Close()
	}
}

func TestSSECompleteCREventDoesNotWaitForAnotherRead(t *testing.T) {
	terminal := errors.New("no subsequent network data")
	body := &chunkedSSEReader{data: []byte("event: message\rdata: {\"value\":\"ready\"}\r\r"), limit: 4096, terminal: terminal}
	stream := NewStream[map[string]string](NewDecoder(&http.Response{Body: body}), nil)
	defer stream.Close()
	if !stream.Next() || stream.Current()["value"] != "ready" {
		t.Fatalf("complete event not delivered: %v", stream.Err())
	}
	if body.reads != 1 {
		t.Fatalf("read past completed event: %d reads", body.reads)
	}
	if stream.Next() || !errors.Is(stream.Err(), terminal) {
		t.Fatalf("transport error lost: %v", stream.Err())
	}
}
