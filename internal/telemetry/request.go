package telemetry

import (
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel/attribute"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/tracing"
)

// MessagesCreateRequest is what the span of a call to the Messages create method
// records of the call's parameters. The method fills it from its typed
// parameters. A field that the parameters leave out is its zero value, and
// records nothing.
type MessagesCreateRequest struct {
	Model       string
	MaxTokens   int64
	ServiceTier string
	// Effort is output_config.effort.
	Effort string
	// ThinkingType is thinking.type, and ThinkingBudgetTokens is
	// thinking.budget_tokens, or nil when the thinking configuration has none.
	ThinkingType         string
	ThinkingBudgetTokens *int64
	// HasOutputFormat is whether the parameters set an output format:
	// output_config.format, or the deprecated output_format.
	HasOutputFormat bool
	// Stream is whether the method streams the Message.
	Stream bool
	// Betas are the names in the request's anthropic-beta headers.
	Betas []string
	// System and Messages are the JSON of the parameters of those names, which the
	// span's content attributes are made of, and StopSequences is the parameter of
	// that name, which is one of them. The method sets them in content mode only
	// ([tracing.APICall.Content]).
	System, Messages []byte
	StopSequences    []string
	// content is System and Messages as [SetMessagesCreateRequestContent]
	// serialized them, and contentSize is the length of the content attributes
	// that it set, when the output was not known.
	content     inputContent
	contentSize int
}

// MessagesCreateStartAttributes returns the attributes that the span of a call to
// the Messages create method starts with, so that a sampler sees them: the
// operation, and what request has of the call's parameters and header.
func MessagesCreateStartAttributes(request *MessagesCreateRequest) []attribute.KeyValue {
	attrs := []attribute.KeyValue{GenAIOperationNameChat}
	if request.Model != "" {
		attrs = append(attrs, GenAIRequestModelKey.String(request.Model))
	}
	attrs = append(attrs, GenAIRequestMaxTokensKey.Int64(request.MaxTokens))
	if request.ServiceTier != "" {
		attrs = append(attrs, AnthropicRequestServiceTierKey.String(request.ServiceTier))
	}
	if request.Effort != "" {
		attrs = append(attrs, GenAIRequestReasoningLevelKey.String(request.Effort))
	}
	if request.ThinkingType != "" {
		attrs = append(attrs, AnthropicThinkingTypeKey.String(request.ThinkingType))
	}
	if request.ThinkingBudgetTokens != nil {
		attrs = append(attrs, AnthropicThinkingBudgetTokensKey.Int64(*request.ThinkingBudgetTokens))
	}
	if request.HasOutputFormat {
		attrs = append(attrs, GenAIOutputTypeJSON)
	}
	if request.Stream {
		attrs = append(attrs, GenAIRequestStreamKey.Bool(true))
	}
	if len(request.Betas) > 0 {
		attrs = append(attrs, AnthropicBetaKey.StringSlice(request.Betas))
	}
	return attrs
}

// SetMessagesCreateRequestContent sets on the span of call, which is in content
// mode, the stop sequences, system instructions and messages of request. They
// cannot be among the attributes that the span starts with: serializing them is
// worth it only once the span records. The response may have them set again:
// they share one size limit, that of the call's options
// ([tracing.Options.MaxContentBytes]), with the output ([inputContent.attributes]).
func SetMessagesCreateRequestContent(call *tracing.APICall, request *MessagesCreateRequest) {
	// A failure to record telemetry must not fail the API call.
	defer func() { tracing.HandleRecovered(recover()) }()

	request.content = newInputContent(gjson.ParseBytes(request.System), gjson.ParseBytes(request.Messages), request.StopSequences)
	attrs, size := request.content.attributes(nil, false, call.Options.MaxContentBytes())
	request.contentSize = size
	call.Span.SetAttributes(attrs...)
}

// BetaNames returns the names in the anthropic-beta headers of header.
func BetaNames(header http.Header) []string {
	var names []string
	for _, value := range header.Values("anthropic-beta") {
		for _, name := range strings.Split(value, ",") {
			if name = strings.TrimSpace(name); name != "" {
				names = append(names, name)
			}
		}
	}
	return names
}
