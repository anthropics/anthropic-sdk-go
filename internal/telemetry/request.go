package telemetry

import (
	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel/attribute"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/tracing"
)

// MessagesCreateRequest is what the span of a call to the Messages create method
// records of the call's parameters, beside the attributes that their
// spanAttributes method returns. The method fills it from its typed parameters.
// A field that the parameters leave out is its zero value, and records nothing.
type MessagesCreateRequest struct {
	// HasOutputFormat is whether the parameters set an output format:
	// output_config.format, or the deprecated output_format.
	HasOutputFormat bool
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

// MessagesCreateStartAttributes returns what the span of a call to the Messages
// create method starts with of request, so that a sampler sees it: the output
// type, when the parameters set an output format.
func MessagesCreateStartAttributes(request *MessagesCreateRequest) []attribute.KeyValue {
	if request.HasOutputFormat {
		return []attribute.KeyValue{GenAIOutputTypeJSON}
	}
	return nil
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
