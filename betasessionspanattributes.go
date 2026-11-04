package anthropic

import (
	"go.opentelemetry.io/otel/attribute"
)

// spanAttributes returns the attributes that the OpenTelemetry span of a call with
// r starts with: one for each property that the span records and r sets.
func (r *BetaSessionNewParams) spanAttributes() (attrs []attribute.KeyValue) {
	if r.EnvironmentID != "" {
		attrs = append(attrs, attribute.String("anthropic.session.environment.id", r.EnvironmentID))
	}
	return attrs
}
