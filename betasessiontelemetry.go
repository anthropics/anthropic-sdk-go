package anthropic

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// sessionsCreateSpan returns the option that makes the OpenTelemetry span of a
// call to [BetaSessionService.New] start with the ID of the Agent of params, as the
// API's span of the same name has it. The Agent is an ID, or an object that has one.
func sessionsCreateSpan(params BetaSessionNewParams) option.RequestOption {
	return requestconfig.WithSpanRecorder("anthropic.sessions.create", requestconfig.SpanRecorder{
		StartAttributes: func(http.Header) []attribute.KeyValue {
			agentID := params.Agent.OfString.Value
			if id := params.Agent.GetID(); id != nil {
				agentID = *id
			}
			if agentID == "" {
				return nil
			}
			return []attribute.KeyValue{attribute.String("gen_ai.agent.id", agentID)}
		},
	})
}
