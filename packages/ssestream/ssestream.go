package ssestream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"
	"strings"

	"github.com/anthropics/anthropic-sdk-go/internal/apierror"
)

type Decoder interface {
	Event() Event
	Next() bool
	Close() error
	Err() error
}

func NewDecoder(res *http.Response) Decoder {
	if res == nil || res.Body == nil {
		return nil
	}

	body := res.Body
	// The span of a traced API call follows the decoded events, below, instead of the body.
	traced, isTraced := body.(tracedBody)
	if isTraced {
		body = traced.Unwrap()
	}

	var decoder Decoder
	contentType := res.Header.Get("content-type")
	if t, ok := decoderTypes[contentType]; ok {
		decoder = t(body)
	} else {
		scn := bufio.NewScanner(body)
		scn.Buffer(nil, bufio.MaxScanTokenSize<<9)
		decoder = &eventStreamDecoder{rc: body, scn: scn}
	}
	if isTraced {
		decoder = traced.TraceEvents(decoder)
	}

	// richDecoder needs the http request to provide helpful errors
	// Aside from the error case there should be no difference
	// from the underlying decoder.
	if res.Request != nil {
		return richErrorDecoder{Decoder: decoder, resp: res}
	}

	return decoder
}

// tracedBody is the response body of an API call that has an OpenTelemetry
// span: a *ResponseBody of package internal/tracing, which this package cannot import.
type tracedBody interface {
	// Unwrap returns the body to decode. Reading it leaves the span open.
	Unwrap() io.ReadCloser
	// TraceEvents returns the decoder to use in place of decoder, so that the
	// span follows the decoded events.
	TraceEvents(decoder Decoder) Decoder
}

// tracedDecoder is the decoder that TraceEvents of a [tracedBody] returns. A
// [Stream] shows it what each event was decoded into, so that the span does not
// decode the event again.
type tracedDecoder interface {
	// TraceDecoded sees a pointer to the value that the current event was decoded
	// into. The next event overwrites the value.
	TraceDecoded(event any)
}

func asTracedDecoder(decoder Decoder) tracedDecoder {
	if rich, ok := decoder.(richErrorDecoder); ok {
		decoder = rich.Decoder
	}
	traced, _ := decoder.(tracedDecoder)
	return traced
}

// richErrorDecoder wraps a Decoder and carries the original [*http.Response]
// so that it can construct rich API errors from SSE error events.
//
// A richErrorDecoder is used by [Stream] to construct errors whenever possible.
// [Stream] needs to rely on the underlying decoder, because the [Stream] has
// no access to the original http request.
//
// There should be no other differences from the underlying decoder.
type richErrorDecoder struct {
	Decoder
	resp *http.Response
}

func (d *richErrorDecoder) newAPIError(errorJSON []byte) error {
	aerr := &apierror.Error{}
	if d.resp != nil {
		aerr.Request = d.resp.Request
		aerr.Response = d.resp
		aerr.StatusCode = d.resp.StatusCode
		aerr.RequestID = d.resp.Header.Get("request-id")
		aerr.WorkspaceID = d.resp.Header.Get("anthropic-workspace-id")
	}
	if aerr.UnmarshalJSON(errorJSON) != nil {
		return fmt.Errorf("received error while streaming: %s", string(errorJSON))
	}
	return aerr
}

var decoderTypes = map[string](func(io.ReadCloser) Decoder){}

func RegisterDecoder(contentType string, decoder func(io.ReadCloser) Decoder) {
	decoderTypes[strings.ToLower(contentType)] = decoder
}

type Event struct {
	Type string
	Data []byte
}

// A base implementation of a Decoder for text/event-stream.
type eventStreamDecoder struct {
	evt Event
	rc  io.ReadCloser
	scn *bufio.Scanner
	err error
}

func (s *eventStreamDecoder) Next() bool {
	if s.err != nil {
		return false
	}

	event := ""
	data := bytes.NewBuffer(nil)

	for s.scn.Scan() {
		txt := s.scn.Bytes()

		// Dispatch event on an empty line
		if len(txt) == 0 {
			if event == "" && data.Len() == 0 {
				continue
			}
			s.evt = Event{
				Type: event,
				Data: data.Bytes(),
			}
			return true
		}

		// Split a string like "event: bar" into name="event" and value=" bar".
		name, value, _ := bytes.Cut(txt, []byte(":"))

		// Consume an optional space after the colon if it exists.
		if len(value) > 0 && value[0] == ' ' {
			value = value[1:]
		}

		switch string(name) {
		case "":
			// An empty line in the for ": something" is a comment and should be ignored.
			continue
		case "event":
			event = string(value)
		case "data":
			_, s.err = data.Write(value)
			if s.err != nil {
				break
			}
			_, s.err = data.WriteRune('\n')
			if s.err != nil {
				break
			}
		}
	}

	if s.scn.Err() != nil {
		s.err = s.scn.Err()
	}

	return false
}

func (s *eventStreamDecoder) Event() Event {
	return s.evt
}

func (s *eventStreamDecoder) Close() error {
	return s.rc.Close()
}

func (s *eventStreamDecoder) Err() error {
	return s.err
}

func (e *eventStreamDecoder) All() iter.Seq[Event] {
	return func(yield func(Event) bool) {
		for e.Next() {
			if !yield(e.Event()) {
				break
			}
		}
	}
}

type Stream[T any] struct {
	decoder Decoder
	cur     T
	err     error
	traced  tracedDecoder
}

func NewStream[T any](decoder Decoder, err error) *Stream[T] {
	return &Stream[T]{
		decoder: decoder,
		err:     err,
		traced:  asTracedDecoder(decoder),
	}
}

// Next returns false if the stream has ended or an error occurred.
// Call Stream.Current() to get the current value.
// Call Stream.Err() to get the error.
//
//		for stream.Next() {
//			data := stream.Current()
//		}
//
//	 	if stream.Err() != nil {
//			...
//	 	}
func (s *Stream[T]) Next() bool {
	if s.err != nil {
		return false
	}

	for s.decoder.Next() {
		switch s.decoder.Event().Type {
		case "ping":
			continue
		case "":
			continue
		case "error":
			data := s.decoder.Event().Data
			if ed, ok := s.decoder.(richErrorDecoder); ok {
				s.err = ed.newAPIError(data)
			} else {
				s.err = fmt.Errorf("received error while streaming: %s", string(data))
			}
			return false
		}

		if len(bytes.TrimSuffix(s.decoder.Event().Data, []byte("\n"))) == 0 {
			continue
		}

		prev := s.cur
		var zero T
		s.cur = zero // decode in place; a fresh local would escape and allocate per event
		s.err = json.Unmarshal(s.decoder.Event().Data, &s.cur)
		if s.err != nil {
			s.cur = prev
			return false
		}
		if s.traced != nil {
			s.traced.TraceDecoded(&s.cur)
		}
		return true
	}

	// decoder.Next() may be false because of an error
	s.err = s.decoder.Err()

	return false
}

func (s *Stream[T]) Current() T {
	return s.cur
}

func (s *Stream[T]) Err() error {
	return s.err
}

func (s *Stream[T]) Close() error {
	if s.decoder == nil {
		// already closed
		return nil
	}
	return s.decoder.Close()
}

func (s *Stream[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for s.Next() {
			if !yield(s.Current()) {
				break
			}
		}
	}
}
