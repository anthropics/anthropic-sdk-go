package telemetry

import (
	"sync"
	"time"

	"github.com/tidwall/gjson"

	"github.com/anthropics/anthropic-sdk-go/internal/tracing"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
)

// MessageStreamObserver returns what the span of call, whose response streams a
// Message to the application, follows of the events: it records the Message as
// the application decodes them, and message_stop is the stream's last event. A
// stream that ends before it leaves the span what the events built of the
// Message's content. request is what the span recorded of the call's
// parameters, and message is an empty Message of the method's API, which the
// events build.
func MessageStreamObserver(call *tracing.APICall, request *MessagesCreateRequest, message StreamedMessage) tracing.StreamObserver {
	recorder := &messageStreamRecorder{call: call, request: request, message: message}
	return tracing.StreamObserver{OnEvent: recorder.arrived, OnDecoded: recorder.record, LastEvent: "message_stop", OnIncomplete: recorder.cutShort}
}

// StreamedMessage is a Message that the SDK's accumulator builds from the
// events of its stream: a Message of the Messages API or of its beta, behind its
// Accumulate method. The streaming method passes it in, because this package
// cannot import the package that imports it.
type StreamedMessage interface {
	// AccumulateEvent applies one event of the stream to the Message: event points
	// to the value that the application's stream decoded it into, so that it is
	// not decoded again. The stream overwrites the value with its next event.
	AccumulateEvent(event any) error
	// RawJSON returns the JSON of the Message as the events have built it so far,
	// as the API would have sent it unstreamed. It is whole once message_stop was
	// applied.
	RawJSON() string
	// SnapshotJSON returns the JSON of the Message of a stream that ended before
	// message_stop, with the content that the events built: a block that was still
	// open has what its deltas brought, and a tool call whose input is not whole
	// yet has an empty one.
	SnapshotJSON() string
}

// messageStreamRecorder records a streamed Message on the span of its call as
// its events arrive, so that nothing is left to record when the span ends,
// whichever way it ends.
type messageStreamRecorder struct {
	call *tracing.APICall
	// request is what the span recorded of the call's parameters.
	request *MessagesCreateRequest
	// timed is set once the time to the first chunk is on the span.
	timed bool
	// mu guards the fields after it: the application may close its stream while
	// it is decoding an event of it.
	mu sync.Mutex
	// message is the Message that the SDK's accumulator builds from the events.
	message StreamedMessage
	// started is set once message_start began the Message.
	started bool
	// noContent is set once the Message's content failed to build: it would be
	// wrong, so the stream is recorded without gen_ai.output.messages.
	noContent bool
}

// arrived sets the time to the first chunk on the span when the stream's first
// content_block_delta event arrives, whatever the type of its delta and whether
// or not the application's stream decodes it. A stream without such an event
// has no time to the first chunk.
func (d *messageStreamRecorder) arrived(event ssestream.Event) {
	defer d.recovered()

	if event.Type == "content_block_delta" && !d.timed {
		d.timed = true
		d.call.Span.SetAttributes(GenAIResponseTimeToFirstChunkKey.Float64(time.Since(d.call.Start).Seconds()))
	}
}

// record applies decoded, which the application's stream decoded event into, to
// the Message that the SDK's accumulator builds, and sets on the span what the
// event adds to it: the response attributes after message_start, with those
// that message_delta ends, such as the stop reason and the output tokens, after
// message_delta, and the content at message_stop in content mode. Outside
// content mode the blocks of the content are started and never filled. The size
// of the content is not limited: the application asked for it, and max_tokens
// bounds it.
func (d *messageStreamRecorder) record(event ssestream.Event, decoded any) {
	d.mu.Lock()
	defer d.mu.Unlock()
	defer d.recovered()

	call, message := d.call, d.message
	content := call.Content() && !d.noContent
	switch event.Type {
	case "message_start", "message_delta":
		if message.AccumulateEvent(decoded) != nil {
			d.noContent = true
			return
		}
		d.started = true
		var response responseMessage
		response.fromJSON(gjson.Parse(message.RawJSON()))
		// message_start has output tokens too, of a Message that has yet to stream.
		call.Span.SetAttributes(response.attributes(event.Type == "message_delta")...)
	case "content_block_start":
		// Applied without content mode too: a fallback block names the model that
		// served the response, and the accumulator takes the Message's model from it.
		if !d.noContent && message.AccumulateEvent(decoded) != nil {
			d.noContent = true
		}
	case "content_block_delta", "content_block_stop":
		if content && message.AccumulateEvent(decoded) != nil {
			d.noContent = true
		}
	case "message_stop":
		if content && message.AccumulateEvent(decoded) != nil {
			content = false
		}
		var response responseMessage
		response.fromJSON(gjson.Parse(message.RawJSON()))
		attrs := response.attributes(true)
		if content && response.content.Exists() {
			attrs = append(attrs, outputContentAttributes(call, d.request, contentParts(response.content))...)
		}
		call.Span.SetAttributes(attrs...)
	}
}

// cutShort sets gen_ai.output.messages on the span of a stream that ends before
// message_stop, in content mode, from what the events built of the Message. The
// span keeps the response attributes of the events that arrived. The content of
// a stream without message_start, and of one with an event that the accumulator
// rejected, stays off the span.
func (d *messageStreamRecorder) cutShort() {
	d.mu.Lock()
	defer d.mu.Unlock()
	defer d.recovered()

	if !d.started || d.noContent || !d.call.Content() {
		return
	}
	if content := gjson.Get(d.message.SnapshotJSON(), "content"); content.Exists() {
		d.call.Span.SetAttributes(outputContentAttributes(d.call, d.request, contentParts(content))...)
	}
}

// recovered is deferred by what records an event. The application is decoding
// its stream: a failure to record must not reach it.
func (d *messageStreamRecorder) recovered() {
	if recovered := recover(); recovered != nil {
		d.noContent = true
		tracing.HandleRecovered(recovered)
	}
}
