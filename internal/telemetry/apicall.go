// Package telemetry sets the attributes of the Claude API on the OpenTelemetry
// spans of API calls, and emits the spans of the tool runners.
package telemetry

import (
	"go.opentelemetry.io/otel/attribute"

	"github.com/anthropics/anthropic-sdk-go/internal/tracing"
)

// SetAPICallAttributes sets on the span of call the anthropic.* attributes that
// every API call has from the headers of its response. requestconfig calls it
// once the call has a result. The span started with [APICallAttributes], and the
// SDK methods set the attributes of their own requests and responses.
func SetAPICallAttributes(call *tracing.APICall) {
	res := call.Response
	if res == nil || !call.Span.IsRecording() {
		return
	}
	var attrs []attribute.KeyValue
	for header, key := range responseHeaderAttributes {
		if v := res.Header.Get(header); v != "" {
			attrs = append(attrs, key.String(v))
		}
	}
	call.Span.SetAttributes(attrs...)
}

// CommonAttributes returns the attributes that every span of opts has. A span
// starts with them, so that a sampler sees them.
func CommonAttributes(opts tracing.Options) []attribute.KeyValue {
	provider := opts.ProviderName
	if provider == "" {
		provider = ProviderAnthropic
	}
	attrs := []attribute.KeyValue{GenAIProviderNameKey.String(provider)}
	if opts.ConversationID != "" {
		attrs = append(attrs, GenAIConversationIDKey.String(opts.ConversationID))
	}
	return attrs
}

// APICallAttributes returns the attributes of opts that the span of an API call
// starts with: [CommonAttributes], and the conversation ID of a call to the path
// of a Session thread. pathIDs are the IDs in the request's path. The API links
// its own spans of a thread by the thread's ID, so the span has that ID as
// gen_ai.conversation.id too, unless the application set a conversation ID.
func APICallAttributes(opts tracing.Options, pathIDs []attribute.KeyValue) []attribute.KeyValue {
	attrs := CommonAttributes(opts)
	if opts.ConversationID != "" {
		return attrs
	}
	for _, id := range pathIDs {
		if id.Key == AnthropicSessionThreadIDKey {
			attrs = append(attrs, GenAIConversationIDKey.String(id.Value.AsString()))
		}
	}
	return attrs
}

var responseHeaderAttributes = map[string]attribute.Key{
	"request-id":                AnthropicRequestIDKey,
	"anthropic-organization-id": AnthropicOrganizationIDKey,
	"anthropic-workspace-id":    AnthropicWorkspaceIDKey,
}
