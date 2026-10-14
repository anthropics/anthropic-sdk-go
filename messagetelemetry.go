package anthropic

import (
	"errors"
	"net/http"

	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel/attribute"

	shimjson "github.com/anthropics/anthropic-sdk-go/internal/encoding/json"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/internal/telemetry"
	"github.com/anthropics/anthropic-sdk-go/internal/tracing"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

// messagesCreateSpan returns the option that makes the OpenTelemetry span of a
// call to [MessageService.New] or [MessageService.NewStreaming] record params
// and the Message: the one that New decodes into *res, or, when res is nil, the
// one that the application reads from the stream.
func messagesCreateSpan(params MessageNewParams, res **Message) option.RequestOption {
	var request *telemetry.MessagesCreateRequest
	return requestconfig.WithSpanRecorder("anthropic.messages.create", requestconfig.SpanRecorder{
		StartAttributes: func(header http.Header) []attribute.KeyValue {
			request = messagesCreateRequest(params, header, res == nil)
			return telemetry.MessagesCreateStartAttributes(request)
		},
		RecordRequest: func(call *tracing.APICall) {
			if !call.Content() {
				return
			}
			// A system prompt that the parameters leave out is not sent.
			if len(params.System) > 0 {
				request.System = marshalParams(params.System)
			}
			request.Messages = marshalParams(params.Messages)
			request.StopSequences = params.StopSequences
			telemetry.SetMessagesCreateRequestContent(call, request)
		},
		RecordResult: func(call *tracing.APICall, err error) {
			if res != nil {
				setMessagesCreateResponseAttributes(call, request, *res, err)
			}
		},
		StreamObserver: func(call *tracing.APICall) tracing.StreamObserver {
			return telemetry.MessageStreamObserver(call, request, &streamedMessage[MessageStreamEventUnion, *Message]{message: new(Message), stop: MessageStreamEventUnion{Type: "message_stop"}})
		},
	})
}

// betaMessagesCreateSpan is [messagesCreateSpan] for [BetaMessageService.New]
// and [BetaMessageService.NewStreaming].
func betaMessagesCreateSpan(params BetaMessageNewParams, res **BetaMessage) option.RequestOption {
	var request *telemetry.MessagesCreateRequest
	return requestconfig.WithSpanRecorder("anthropic.messages.create", requestconfig.SpanRecorder{
		StartAttributes: func(header http.Header) []attribute.KeyValue {
			request = betaMessagesCreateRequest(params, header, res == nil)
			return telemetry.MessagesCreateStartAttributes(request)
		},
		RecordRequest: func(call *tracing.APICall) {
			if !call.Content() {
				return
			}
			// A system prompt that the parameters leave out is not sent.
			if len(params.System) > 0 {
				request.System = marshalParams(params.System)
			}
			request.Messages = marshalParams(params.Messages)
			request.StopSequences = params.StopSequences
			telemetry.SetMessagesCreateRequestContent(call, request)
		},
		RecordResult: func(call *tracing.APICall, err error) {
			if res != nil {
				setMessagesCreateResponseAttributes(call, request, *res, err)
			}
		},
		StreamObserver: func(call *tracing.APICall) tracing.StreamObserver {
			return telemetry.MessageStreamObserver(call, request, &streamedMessage[BetaRawMessageStreamEventUnion, *BetaMessage]{message: new(BetaMessage), stop: BetaRawMessageStreamEventUnion{Type: "message_stop"}})
		},
	})
}

// countTokensSpan returns the option that makes the OpenTelemetry span of a call
// to a CountTokens method record model and the count that the method decodes
// into *res, which inputTokens reads. *res stays nil when the application had
// the body decoded into a value of its own.
func countTokensSpan[T any](model string, res **T, inputTokens func(*T) int64) option.RequestOption {
	return requestconfig.WithSpanRecorder("anthropic.messages.count_tokens", requestconfig.SpanRecorder{
		StartAttributes: func(http.Header) []attribute.KeyValue {
			return []attribute.KeyValue{telemetry.GenAIRequestModelKey.String(model)}
		},
		RecordResult: func(call *tracing.APICall, err error) {
			if err == nil && *res != nil {
				call.Span.SetAttributes(telemetry.AnthropicMessageCountTokensInputTokensKey.Int64(inputTokens(*res)))
			}
		},
	})
}

// messageBatchesCreateSpan returns the option that makes the OpenTelemetry span
// of a call to [MessageBatchService.New] or [BetaMessageBatchService.New] start
// with requests, the number of requests of the batch.
func messageBatchesCreateSpan(requests int) option.RequestOption {
	return requestconfig.WithSpanRecorder("anthropic.messages.batches.create", requestconfig.SpanRecorder{
		StartAttributes: func(http.Header) []attribute.KeyValue {
			return []attribute.KeyValue{telemetry.AnthropicMessageBatchRequestCountKey.Int(requests)}
		},
	})
}

// messagesCreateRequest returns what the OpenTelemetry span of a call to
// [MessageService.New] or [MessageService.NewStreaming] starts with of params.
// header is the header of the request, and stream is whether the method streams
// the Message.
func messagesCreateRequest(params MessageNewParams, header http.Header, stream bool) *telemetry.MessagesCreateRequest {
	return &telemetry.MessagesCreateRequest{
		Model:                string(params.Model),
		MaxTokens:            params.MaxTokens,
		Temperature:          params.Temperature,
		TopK:                 params.TopK,
		TopP:                 params.TopP,
		ServiceTier:          string(params.ServiceTier),
		Effort:               string(params.OutputConfig.Effort),
		ThinkingType:         thinkingType(params.Thinking),
		ThinkingBudgetTokens: params.Thinking.GetBudgetTokens(),
		HasOutputFormat:      !param.IsOmitted(params.OutputConfig.Format),
		Stream:               stream,
		Betas:                telemetry.BetaNames(header),
	}
}

// betaMessagesCreateRequest is [messagesCreateRequest] for
// [BetaMessageService.New] and [BetaMessageService.NewStreaming].
func betaMessagesCreateRequest(params BetaMessageNewParams, header http.Header, stream bool) *telemetry.MessagesCreateRequest {
	return &telemetry.MessagesCreateRequest{
		Model:                string(params.Model),
		MaxTokens:            params.MaxTokens,
		Temperature:          params.Temperature,
		TopK:                 params.TopK,
		TopP:                 params.TopP,
		ServiceTier:          string(params.ServiceTier),
		Effort:               string(params.OutputConfig.Effort),
		ThinkingType:         thinkingType(params.Thinking),
		ThinkingBudgetTokens: params.Thinking.GetBudgetTokens(),
		HasOutputFormat:      !param.IsOmitted(params.OutputConfig.Format) || !param.IsOmitted(params.OutputFormat),
		Stream:               stream,
		Betas:                telemetry.BetaNames(header),
	}
}

// thinkingType returns thinking.type as the request sends it, read from the JSON
// of thinking, the thinking configuration of the parameters. The type is sent
// even when the application leaves it to its default, and a configuration that
// this version of the SDK has no variant for has one too.
func thinkingType(thinking any) string {
	if param.IsOmitted(thinking) {
		return ""
	}
	return gjson.GetBytes(marshalParams(thinking), "type").Str
}

// marshalParams returns the JSON of v, a value of the request parameters, as the
// SDK encodes the request itself: the JSON that the parameter types make is
// taken as it is. encoding/json would validate and compact all of it again.
func marshalParams(v any) []byte {
	data, _ := shimjson.Marshal(v, shimjson.WithSkipCompaction(true))
	return data
}

// setMessagesCreateResponseAttributes sets on the span of call the attributes
// of res, the Message that [MessageService.New] or [BetaMessageService.New]
// returned with err. request is what the span recorded of the call's parameters.
// res is nil when the application had the body decoded into a value of its own
// ([option.WithResponseBodyInto]).
func setMessagesCreateResponseAttributes[M interface {
	Message | BetaMessage
	RawJSON() string
}](call *tracing.APICall, request *telemetry.MessagesCreateRequest, res *M, err error) {
	if err == nil && res != nil {
		telemetry.SetMessagesCreateResponseAttributes(call, request, (*res).RawJSON())
	}
}

// streamedMessage is the [telemetry.StreamedMessage] of a Message of type M
// whose stream has events of type E: the OpenTelemetry span of the streaming
// method reads the streamed Message from what Accumulate builds.
type streamedMessage[E any, M interface {
	Accumulate(E) error
	RawJSON() string
}] struct {
	message M
	// stop is the message_stop event of the stream.
	stop E
}

func (s *streamedMessage[E, M]) AccumulateEvent(event any) error {
	decoded, ok := event.(*E)
	if !ok {
		return errors.New("the stream decodes its events into a type of its own")
	}
	// Accumulate takes the event by value and keeps nothing that the next decode writes to.
	return s.message.Accumulate(*decoded)
}

func (s *streamedMessage[E, M]) RawJSON() string { return s.message.RawJSON() }

// SnapshotJSON has Accumulate finish the Message as message_stop does, which
// writes every block that the events built into its JSON.
func (s *streamedMessage[E, M]) SnapshotJSON() string {
	if s.message.Accumulate(s.stop) != nil {
		return ""
	}
	return s.message.RawJSON()
}
