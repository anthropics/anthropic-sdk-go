package telemetry

import (
	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel/attribute"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/tracing"
)

// SetMessagesCreateResponseAttributes sets on the span of call what package
// tracing does not read from the Message that a call to the Messages create
// method returned, whose JSON is messageJSON: the attributes that are sums of
// its token counts, and in content mode its content. request is what the span
// keeps of the content of the call's request.
func SetMessagesCreateResponseAttributes(call *tracing.APICall, request *MessagesCreateRequest, messageJSON string) {
	if !call.Span.IsRecording() {
		return
	}
	// A failure to record telemetry must not fail the API call.
	defer func() { tracing.HandleRecovered(recover()) }()

	message := gjson.Parse(messageJSON)
	attrs := tokenSumAttributes(message, true)
	if content := message.Get("content"); call.Content() && content.Exists() {
		attrs = append(attrs, outputContentAttributes(call, request, contentParts(content))...)
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

// tokenSumAttributes returns the attributes of message, the JSON of a Message,
// that are sums of its token counts. Those of its compactions are among them only
// when complete is true: it is a response body, or a stream past message_delta.
func tokenSumAttributes(message gjson.Result, complete bool) []attribute.KeyValue {
	var attrs []attribute.KeyValue
	usage := message.Get("usage")
	// The conventions define the input tokens as the sum of these three.
	inputTokens := func(usage gjson.Result) int64 {
		return usage.Get("input_tokens").Int() + usage.Get("cache_read_input_tokens").Int() + usage.Get("cache_creation_input_tokens").Int()
	}
	if input := usage.Get("input_tokens"); input.Exists() && input.Type != gjson.Null {
		attrs = append(attrs, GenAIUsageInputTokensKey.Int64(inputTokens(usage)))
	}
	// usage.iterations has the usage of each step that made the Message.
	compactions := usage.Get(`iterations.#(type=="compaction")#`).Array()
	if !complete || len(compactions) == 0 {
		return attrs
	}
	// The Message's own usage leaves out the tokens of its compactions.
	var input, output int64
	for _, usage := range compactions {
		input += inputTokens(usage)
		output += usage.Get("output_tokens").Int()
	}
	return append(attrs, AnthropicUsageCompactionInputTokensKey.Int64(input), AnthropicUsageCompactionOutputTokensKey.Int64(output))
}
