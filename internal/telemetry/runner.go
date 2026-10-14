package telemetry

import (
	"context"
	"encoding/json"
	"errors"
	"unicode/utf8"

	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/anthropics/anthropic-sdk-go/internal/tracing"
)

// StartToolRunnerSpan starts the anthropic.messages.tool_runner span of one loop
// of a Messages tool runner whose requests use opts and ask for model, as a
// child of the span in ctx. hasOutputFormat is whether the runner's parameters
// set an output format: output_config.format, or the deprecated output_format.
// It returns nil when opts turn the spans off. The runner sets the span's status
// and ends it.
func StartToolRunnerSpan(ctx context.Context, opts tracing.Options, model string, hasOutputFormat bool) (span trace.Span) {
	// A panic of the application's tracer leaves the loop without a span.
	defer func() { tracing.HandleRecovered(recover()) }()

	tracer := opts.Tracer()
	if tracer == nil {
		return nil
	}
	attrs := append(CommonAttributes(opts), GenAIOperationNameInvokeAgent, GenAIRequestModelKey.String(model))
	if hasOutputFormat {
		attrs = append(attrs, GenAIOutputTypeJSON)
	}
	_, span = tracer.Start(ctx, "anthropic.messages.tool_runner", trace.WithSpanKind(trace.SpanKindInternal), trace.WithAttributes(attrs...))
	return span
}

// MadeByMCPHelper marks a tool that the SDK's MCP helper made: the helper's tool
// type embeds it, and so implements [MCPHelperTool].
type MadeByMCPHelper struct{}

func (MadeByMCPHelper) madeByMCPHelper() {}

// MCPHelperTool is a tool that the SDK's MCP helper made. The tool runner runs it
// in an anthropic.mcp_tool_use span.
type MCPHelperTool interface{ madeByMCPHelper() }

// StartToolUseSpan starts the span of one run of a tool whose name is name, for
// the tool call id, as a child of the span in ctx, and returns ctx with the
// span. The span is nil, and ctx unchanged, when opts turn the spans off. The
// span's name is spanName, which says where the tool comes from
// (anthropic.custom_tool_use, anthropic.tool_use or anthropic.mcp_tool_use),
// then a space and name, or nothing when name is empty. description returns the
// description in the tool's definition that the runner has: the span has it as
// gen_ai.tool.description, without content mode too, unless it is empty or not
// valid UTF-8. It is the application's code, so it is called only when the spans
// are on. toolsetName is the toolset_name of the tool call, or empty when it has
// none. resourceIDs are the IDs of what the tool call belongs to, which only
// the runner knows. [EndToolUseSpan] ends it.
func StartToolUseSpan(ctx context.Context, opts tracing.Options, spanName, name, toolsetName string, description func() string, id string, resourceIDs ...attribute.KeyValue) (toolCtx context.Context, span trace.Span) {
	// A panic of the application's tracer, or of description, leaves the run
	// without a span, in ctx.
	toolCtx = ctx
	defer func() { tracing.HandleRecovered(recover()) }()

	tracer := opts.Tracer()
	if tracer == nil {
		return ctx, nil
	}
	if name != "" {
		spanName += " " + name
	}
	attrs := append(CommonAttributes(opts),
		GenAIOperationNameExecuteTool,
		GenAIToolNameKey.String(name),
		GenAIToolCallIDKey.String(id),
		GenAIToolTypeFunction,
	)
	if toolsetName != "" {
		attrs = append(attrs, AnthropicToolsetNameKey.String(toolsetName))
	}
	attrs = append(attrs, resourceIDs...)
	// A string that is not valid UTF-8 can fail the whole export.
	if description := description(); description != "" && utf8.ValidString(description) {
		attrs = append(attrs, GenAIToolDescriptionKey.String(description))
	}
	return tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindInternal), trace.WithAttributes(attrs...))
}

// EndToolUseSpan ends the span of a tool run that returned err in ctx, now that
// the runner has made the result of it. A run that failed fails the span with
// tool_error, unless the application canceled ctx while the tool ran or err is
// that cancellation: then the run was stopped, not failed, and the status stays
// unset, as it does on the tool runner span and the request spans. ctx is the application's
// context, not one that carries the runner's own per-tool timeout, so a tool
// that runs out of time still fails its span.
//
// In content mode the span gets input, the tool's input, and, unless the run was
// canceled, result: the content that the runner sends back to the model, a
// string or a list of content blocks, which is encoded here. Together they stay
// within the size limit of opts ([tracing.Options.MaxContentBytes]): an input
// that is over the limit on its own is left off the span, and blocks are
// dropped from the end of the result. anthropic.content.truncated says when
// either happened. A result that JSON cannot hold is left off the span. It does
// nothing when span is nil.
func EndToolUseSpan(ctx context.Context, span trace.Span, opts tracing.Options, input json.RawMessage, result any, err error) {
	if span == nil {
		return
	}
	// A panic of the application's span stays out of the runner, and the span still ends.
	defer func() { tracing.HandleRecovered(recover()) }()
	defer span.End()

	canceled := errors.Is(err, context.Canceled) || errors.Is(ctx.Err(), context.Canceled)
	if opts.Content() && span.IsRecording() {
		setToolUseContent(span, opts.MaxContentBytes(), input, result, canceled)
	}
	switch {
	case err == nil:
		span.SetStatus(codes.Ok, "")
	case canceled:
		// Stopped by the application: no status.
	default:
		tracing.SetErrorStatus(span, err, "tool_error")
	}
}

// setToolUseContent sets on span, the span of a tool run in content mode, the
// attributes that [toolUseContentAttributes] returns.
func setToolUseContent(span trace.Span, limit int, input json.RawMessage, result any, canceled bool) {
	// The runner runs a tool in a goroutine of its own, where a failure to record
	// telemetry would end the process. The span still gets its status and ends.
	defer func() { tracing.HandleRecovered(recover()) }()

	span.SetAttributes(toolUseContentAttributes(limit, input, result, canceled)...)
}

// toolUseContentAttributes returns the content attributes of a tool run: the
// tool's input and, unless the run was canceled, result, within limit bytes, as
// [EndToolUseSpan] describes.
func toolUseContentAttributes(limit int, input json.RawMessage, result any, canceled bool) []attribute.KeyValue {
	var attrs []attribute.KeyValue
	// The result has what the arguments leave of the limit, which is all of it
	// when they are left out.
	room := limit
	arguments := withoutHTMLEscapes(string(input))
	truncated := len(arguments) > limit
	if !truncated {
		attrs = append(attrs, GenAIToolCallArgumentsKey.String(arguments))
		room -= len(arguments)
	}
	if resultJSON, err := json.Marshal(result); err == nil && !canceled {
		content, dropped := truncateToolResultContent(toolResultContent(gjson.ParseBytes(resultJSON)), room)
		attrs = append(attrs, AnthropicToolResultContentKey.String(content))
		truncated = truncated || dropped
	}
	if truncated {
		attrs = append(attrs, AnthropicContentTruncatedKey.Bool(true))
	}
	return attrs
}
