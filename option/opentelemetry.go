package option

import (
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
)

// WithOpenTelemetryConversationID returns a RequestOption that links the calls it
// applies to into one conversation, with an ID that the application chooses.
//
// The SDK does not send the ID to the API. It sets the gen_ai.conversation.id
// attribute on the OpenTelemetry span of each call. A tool runner passes the
// option to each request it sends, so every span of its loop has the same ID.
//
// Without the option, the span of a call to a Session thread has the thread's ID
// as its conversation ID, as the API's spans of the thread have, and the span of
// any other call has none.
func WithOpenTelemetryConversationID(id string) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.OpenTelemetry.ConversationID = id
		return nil
	})
}
