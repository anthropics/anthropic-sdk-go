package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/tracing"
)

// StartEnvironmentWorkSpan starts the anthropic.environment.work span of one work
// item of a self-hosted environment, whose ID is workID and which runs the
// Session sessionID, as a child of the span in ctx, and returns ctx with the
// span. The span is nil, and ctx unchanged, when opts turn the spans off. The
// worker sets the span's status and ends it.
func StartEnvironmentWorkSpan(ctx context.Context, opts tracing.Options, environmentID, workID, sessionID string) (workCtx context.Context, span trace.Span) {
	// A panic of the application's tracer leaves the item without a span, in ctx.
	workCtx = ctx
	defer func() { tracing.HandleRecovered(recover()) }()

	tracer := opts.Tracer()
	if tracer == nil {
		return ctx, nil
	}
	// An item is a whole Session, not one conversation of the application's.
	opts.ConversationID = ""
	attrs := append(CommonAttributes(opts),
		attribute.String("anthropic.environment.id", environmentID),
		attribute.String("anthropic.environment.work.id", workID),
		attribute.String("anthropic.session.id", sessionID),
	)
	return tracer.Start(ctx, "anthropic.environment.work", trace.WithSpanKind(trace.SpanKindInternal), trace.WithAttributes(attrs...))
}
