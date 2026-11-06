// Package telemetry sets the attributes of the Claude API on the OpenTelemetry
// spans of API calls, and emits the spans of the tool runners.
package telemetry

import (
	"slices"

	"go.opentelemetry.io/otel/attribute"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/tracing"
)

// CommonAttributes returns the attributes that every span of opts has: those of
// [tracing.ClientAttributes], with the provider of a third-party platform in
// place of theirs, and the conversation ID. A span starts with them, so that a
// sampler sees them.
func CommonAttributes(opts tracing.Options) []attribute.KeyValue {
	attrs := tracing.ClientAttributes()
	if opts.ProviderName != "" {
		attrs = slices.DeleteFunc(attrs, func(kv attribute.KeyValue) bool { return kv.Key == GenAIProviderNameKey })
		attrs = append(attrs, GenAIProviderNameKey.String(opts.ProviderName))
	}
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
