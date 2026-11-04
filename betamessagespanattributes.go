package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/v2/internal/paramutil"
	"go.opentelemetry.io/otel/attribute"
)

// spanAttributes returns the attributes that the OpenTelemetry span of a call with
// r starts with: one for each property that the span records and r sets.
func (r *BetaMessageNewParams) spanAttributes() (attrs []attribute.KeyValue) {
	if r.Model != "" {
		attrs = append(attrs, attribute.String("gen_ai.request.model", string(r.Model)))
	}
	attrs = append(attrs, attribute.Int64("gen_ai.request.max_tokens", r.MaxTokens))
	if r.OutputConfig.Effort != "" {
		attrs = append(attrs, attribute.String("gen_ai.request.reasoning.level", string(r.OutputConfig.Effort)))
	}
	if r.ServiceTier != "" {
		attrs = append(attrs, attribute.String("anthropic.request.service_tier", string(r.ServiceTier)))
	}
	if vt := r.Thinking.OfEnabled; vt != nil {
		attrs = append(attrs, attribute.String("anthropic.thinking.type", *paramutil.AddrOrDefault(vt.Type)))
	}
	if vt := r.Thinking.OfDisabled; vt != nil {
		attrs = append(attrs, attribute.String("anthropic.thinking.type", *paramutil.AddrOrDefault(vt.Type)))
	}
	if vt := r.Thinking.OfBetweenTools; vt != nil {
		attrs = append(attrs, attribute.String("anthropic.thinking.type", *paramutil.AddrOrDefault(vt.Type)))
	}
	if vt := r.Thinking.OfAdaptive; vt != nil {
		attrs = append(attrs, attribute.String("anthropic.thinking.type", *paramutil.AddrOrDefault(vt.Type)))
	}
	if vt := r.Thinking.OfEnabled; vt != nil {
		attrs = append(attrs, attribute.Int64("anthropic.thinking.budget_tokens", vt.BudgetTokens))
	}
	return attrs
}

// spanAttributes returns the attributes that the OpenTelemetry span of a call with
// r starts with: one for each property that the span records and r sets.
func (r *BetaMessageCountTokensParams) spanAttributes() (attrs []attribute.KeyValue) {
	if r.Model != "" {
		attrs = append(attrs, attribute.String("gen_ai.request.model", string(r.Model)))
	}
	return attrs
}
