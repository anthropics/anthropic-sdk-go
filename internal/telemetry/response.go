package telemetry

import (
	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel/attribute"

	"github.com/anthropics/anthropic-sdk-go/internal/tracing"
)

// SetMessagesCreateResponseAttributes sets on the span of call the attributes
// of the Message that a call to the Messages create method returned, whose JSON
// is messageJSON, and in content mode its content. request is what
// the span recorded of the call's parameters.
func SetMessagesCreateResponseAttributes(call *tracing.APICall, request *MessagesCreateRequest, messageJSON string) {
	if !call.Span.IsRecording() {
		return
	}
	// A failure to record telemetry must not fail the API call.
	defer func() { tracing.HandleRecovered(recover()) }()

	var message responseMessage
	message.fromJSON(gjson.Parse(messageJSON))
	attrs := message.attributes(true)
	if call.Content() && message.content.Exists() {
		attrs = append(attrs, outputContentAttributes(call, request, contentParts(message.content))...)
	}
	call.Span.SetAttributes(attrs...)
}

// outputContentAttributes returns gen_ai.output.messages for the Message whose
// content is output. The request's content was given the whole size limit of
// call's options when the call started, so when the output does not fit beside
// what the span has of it, the request's content attributes are among the
// returned ones too: all of the request's parts are made to fit with the output,
// not only those that the span has.
func outputContentAttributes(call *tracing.APICall, request *MessagesCreateRequest, output []part) []attribute.KeyValue {
	limit := call.Options.MaxContentBytes()
	messages := []messageContent{newMessageContent("assistant", output)}
	value := messagesJSON(messages, [][]string{messages[0].parts})
	if request.contentSize+len(value) <= limit {
		return []attribute.KeyValue{GenAIOutputMessagesKey.String(value)}
	}
	attrs, _ := request.content.attributes(messages, true, limit)
	return attrs
}

// usageFields are the fields of a Message's usage that the span reports.
var usageFields = []string{
	"input_tokens",
	"output_tokens",
	"cache_creation_input_tokens",
	"cache_read_input_tokens",
	"cache_creation.ephemeral_5m_input_tokens",
	"cache_creation.ephemeral_1h_input_tokens",
	"service_tier",
	"output_tokens_details.thinking_tokens",
}

// responseMessage is what the span records of a Message: the response body, or
// the Message that the SDK's accumulator built from a stream.
type responseMessage struct {
	id, model, stopReason string
	usage                 map[string]gjson.Result
	// iterations is usage.iterations: the usage of each step that made the Message.
	iterations gjson.Result
	content    gjson.Result
}

// fromJSON fills m from message, the JSON of a Message.
func (m *responseMessage) fromJSON(message gjson.Result) {
	m.id = message.Get("id").String()
	m.model = message.Get("model").String()
	if reason := message.Get("stop_reason"); reason.Type == gjson.String {
		m.stopReason = reason.String()
	}
	m.content = message.Get("content")
	usage := message.Get("usage")
	m.iterations = usage.Get("iterations")
	for _, field := range usageFields {
		if v := usage.Get(field); v.Exists() && v.Type != gjson.Null {
			if m.usage == nil {
				m.usage = map[string]gjson.Result{}
			}
			m.usage[field] = v
		}
	}
}

// attributes returns the attributes that come from the Message. The ones whose
// value the message_delta event of a stream ends are among them only when
// complete is true: the Message is a response body, or its stream is past that event.
func (m *responseMessage) attributes(complete bool) []attribute.KeyValue {
	var attrs []attribute.KeyValue
	if m.id != "" {
		attrs = append(attrs, GenAIResponseIDKey.String(m.id))
	}
	if m.model != "" {
		attrs = append(attrs, GenAIResponseModelKey.String(m.model))
	}
	if input, ok := m.usage["input_tokens"]; ok {
		// The conventions define the input tokens as the sum of these three.
		sum := input.Int() + m.usage["cache_read_input_tokens"].Int() + m.usage["cache_creation_input_tokens"].Int()
		attrs = append(attrs, GenAIUsageInputTokensKey.Int64(sum))
	}
	for field, key := range map[string]attribute.Key{
		"cache_creation_input_tokens":              GenAIUsageCacheWriteInputTokensKey,
		"cache_creation.ephemeral_5m_input_tokens": AnthropicUsageCacheCreationEphemeral5mInputTokensKey,
		"cache_creation.ephemeral_1h_input_tokens": AnthropicUsageCacheCreationEphemeral1hInputTokensKey,
		"cache_read_input_tokens":                  GenAIUsageCacheReadInputTokensKey,
	} {
		if v, ok := m.usage[field]; ok {
			attrs = append(attrs, key.Int64(v.Int()))
		}
	}
	if tier, ok := m.usage["service_tier"]; ok {
		attrs = append(attrs, AnthropicResponseServiceTierKey.String(tier.String()))
	}
	if !complete {
		return attrs
	}
	if v, ok := m.usage["output_tokens"]; ok {
		attrs = append(attrs, GenAIUsageOutputTokensKey.Int64(v.Int()))
	}
	if v, ok := m.usage["output_tokens_details.thinking_tokens"]; ok {
		attrs = append(attrs, GenAIUsageReasoningOutputTokensKey.Int64(v.Int()))
	}
	if m.stopReason != "" {
		attrs = append(attrs, AnthropicMessageStopReasonKey.String(m.stopReason))
	}
	compactions := m.iterations.Get(`#(type=="compaction")#`).Array()
	if len(compactions) == 0 {
		return attrs
	}
	// The Message's own usage leaves out the tokens of its compactions.
	var input, output int64
	for _, usage := range compactions {
		input += usage.Get("input_tokens").Int() + usage.Get("cache_read_input_tokens").Int() + usage.Get("cache_creation_input_tokens").Int()
		output += usage.Get("output_tokens").Int()
	}
	return append(attrs, AnthropicUsageCompactionInputTokensKey.Int64(input), AnthropicUsageCompactionOutputTokensKey.Int64(output))
}
