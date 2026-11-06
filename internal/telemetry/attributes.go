package telemetry

import (
	"go.opentelemetry.io/otel/attribute"
)

// The values of gen_ai.provider.name.
const (
	ProviderBedrock = "aws.bedrock"
	ProviderVertex  = "gcp.vertex_ai"
)

// Attribute keys and values of the OpenTelemetry generative-AI semantic
// conventions, named as the semconv packages of go.opentelemetry.io/otel name
// them. They are copied rather than imported: the gen_ai conventions are still
// in development there, and the SDK does not pin a version of the conventions.
const (
	GenAIConversationIDKey           = attribute.Key("gen_ai.conversation.id")
	GenAIInputMessagesKey            = attribute.Key("gen_ai.input.messages")
	GenAIOperationNameKey            = attribute.Key("gen_ai.operation.name")
	GenAIOutputMessagesKey           = attribute.Key("gen_ai.output.messages")
	GenAIOutputTypeKey               = attribute.Key("gen_ai.output.type")
	GenAIProviderNameKey             = attribute.Key("gen_ai.provider.name")
	GenAIRequestModelKey             = attribute.Key("gen_ai.request.model")
	GenAIRequestStopSequencesKey     = attribute.Key("gen_ai.request.stop_sequences")
	GenAIResponseTimeToFirstChunkKey = attribute.Key("gen_ai.response.time_to_first_chunk")
	GenAISystemInstructionsKey       = attribute.Key("gen_ai.system_instructions")
	GenAIToolCallArgumentsKey        = attribute.Key("gen_ai.tool.call.arguments")
	GenAIToolCallIDKey               = attribute.Key("gen_ai.tool.call.id")
	GenAIToolDescriptionKey          = attribute.Key("gen_ai.tool.description")
	GenAIToolNameKey                 = attribute.Key("gen_ai.tool.name")
	GenAIToolTypeKey                 = attribute.Key("gen_ai.tool.type")
	GenAIUsageInputTokensKey         = attribute.Key("gen_ai.usage.input_tokens")
)

var (
	GenAIOperationNameExecuteTool = GenAIOperationNameKey.String("execute_tool")
	GenAIOperationNameInvokeAgent = GenAIOperationNameKey.String("invoke_agent")
	GenAIOutputTypeJSON           = GenAIOutputTypeKey.String("json")
	GenAIToolTypeFunction         = GenAIToolTypeKey.String("function")
)

// The attribute keys this SDK defines, named in the same style.
const (
	AnthropicContentTruncatedKey            = attribute.Key("anthropic.content.truncated")
	AnthropicMessageBatchRequestCountKey    = attribute.Key("anthropic.message.batch.request_count")
	AnthropicSessionThreadIDKey             = attribute.Key("anthropic.session.thread.id")
	AnthropicToolResultContentKey           = attribute.Key("anthropic.tool_result.content")
	AnthropicToolsetNameKey                 = attribute.Key("anthropic.toolset.name")
	AnthropicUsageCompactionInputTokensKey  = attribute.Key("anthropic.usage.compaction.input_tokens")
	AnthropicUsageCompactionOutputTokensKey = attribute.Key("anthropic.usage.compaction.output_tokens")
)
