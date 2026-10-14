package anthropic

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"

	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/internal/tracing"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// sessionsCreateSpan returns the option that makes the OpenTelemetry span of a
// call to [BetaSessionService.New] record what the API's span of the same name
// does: the Agent and the Environment of params, and the name and the model of
// the Agent of the Session that New decodes into *res. *res stays nil when the
// application had the body decoded into a value of its own.
func sessionsCreateSpan(params BetaSessionNewParams, res **BetaManagedAgentsSession) option.RequestOption {
	return requestconfig.WithSpanRecorder("anthropic.sessions.create", requestconfig.SpanRecorder{
		StartAttributes: func(http.Header) []attribute.KeyValue {
			agentID := params.Agent.OfString.Value
			if id := params.Agent.GetID(); id != nil {
				agentID = *id
			}
			var attrs []attribute.KeyValue
			if agentID != "" {
				attrs = append(attrs, attribute.String("gen_ai.agent.id", agentID))
			}
			if params.EnvironmentID != "" {
				attrs = append(attrs, attribute.String("anthropic.session.environment.id", params.EnvironmentID))
			}
			return attrs
		},
		RecordResult: func(call *tracing.APICall, err error) {
			if err != nil || *res == nil {
				return
			}
			agent := (*res).Agent
			// The model is the one that the Session resolved, which params may not name.
			call.Span.SetAttributes(attribute.String("gen_ai.agent.name", agent.Name), attribute.String("gen_ai.request.model", agent.Model.ID))
		},
	})
}
