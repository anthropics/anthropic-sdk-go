package telemetry

import (
	"go.opentelemetry.io/otel/attribute"
)

// The values of gen_ai.provider.name.
const (
	ProviderAnthropic = "anthropic"
	ProviderBedrock   = "aws.bedrock"
	ProviderVertex    = "gcp.vertex_ai"
)

// Attribute keys and values of the OpenTelemetry generative-AI semantic
// conventions, named as the semconv packages of go.opentelemetry.io/otel name
// them. They are copied rather than imported: the gen_ai conventions are still
// in development there, and the SDK does not pin a version of the conventions.
const (
	GenAIConversationIDKey             = attribute.Key("gen_ai.conversation.id")
	GenAIInputMessagesKey              = attribute.Key("gen_ai.input.messages")
	GenAIOperationNameKey              = attribute.Key("gen_ai.operation.name")
	GenAIOutputMessagesKey             = attribute.Key("gen_ai.output.messages")
	GenAIOutputTypeKey                 = attribute.Key("gen_ai.output.type")
	GenAIProviderNameKey               = attribute.Key("gen_ai.provider.name")
	GenAIRequestMaxTokensKey           = attribute.Key("gen_ai.request.max_tokens")
	GenAIRequestModelKey               = attribute.Key("gen_ai.request.model")
	GenAIRequestReasoningLevelKey      = attribute.Key("gen_ai.request.reasoning.level")
	GenAIRequestStopSequencesKey       = attribute.Key("gen_ai.request.stop_sequences")
	GenAIRequestStreamKey              = attribute.Key("gen_ai.request.stream")
	GenAIRequestTemperatureKey         = attribute.Key("gen_ai.request.temperature")
	GenAIRequestTopKKey                = attribute.Key("gen_ai.request.top_k")
	GenAIRequestTopPKey                = attribute.Key("gen_ai.request.top_p")
	GenAIResponseIDKey                 = attribute.Key("gen_ai.response.id")
	GenAIResponseModelKey              = attribute.Key("gen_ai.response.model")
	GenAIResponseTimeToFirstChunkKey   = attribute.Key("gen_ai.response.time_to_first_chunk")
	GenAISystemInstructionsKey         = attribute.Key("gen_ai.system_instructions")
	GenAIToolCallArgumentsKey          = attribute.Key("gen_ai.tool.call.arguments")
	GenAIToolCallIDKey                 = attribute.Key("gen_ai.tool.call.id")
	GenAIToolDescriptionKey            = attribute.Key("gen_ai.tool.description")
	GenAIToolNameKey                   = attribute.Key("gen_ai.tool.name")
	GenAIToolTypeKey                   = attribute.Key("gen_ai.tool.type")
	GenAIUsageCacheReadInputTokensKey  = attribute.Key("gen_ai.usage.cache_read.input_tokens")
	GenAIUsageCacheWriteInputTokensKey = attribute.Key("gen_ai.usage.cache_write.input_tokens")
	GenAIUsageInputTokensKey           = attribute.Key("gen_ai.usage.input_tokens")
	GenAIUsageOutputTokensKey          = attribute.Key("gen_ai.usage.output_tokens")
	GenAIUsageReasoningOutputTokensKey = attribute.Key("gen_ai.usage.reasoning.output_tokens")
)

var (
	GenAIOperationNameChat        = GenAIOperationNameKey.String("chat")
	GenAIOperationNameExecuteTool = GenAIOperationNameKey.String("execute_tool")
	GenAIOperationNameInvokeAgent = GenAIOperationNameKey.String("invoke_agent")
	GenAIOutputTypeJSON           = GenAIOutputTypeKey.String("json")
	GenAIToolTypeFunction         = GenAIToolTypeKey.String("function")
)

// The attribute keys this SDK defines, named in the same style.
const (
	AnthropicRequestIDKey                                = attribute.Key("anthropic.request.id")
	AnthropicOrganizationIDKey                           = attribute.Key("anthropic.organization.id")
	AnthropicWorkspaceIDKey                              = attribute.Key("anthropic.workspace.id")
	AnthropicBetaKey                                     = attribute.Key("anthropic.beta")
	AnthropicContentTruncatedKey                         = attribute.Key("anthropic.content.truncated")
	AnthropicMessageBatchRequestCountKey                 = attribute.Key("anthropic.message.batch.request_count")
	AnthropicMessageStopReasonKey                        = attribute.Key("anthropic.message.stop_reason")
	AnthropicMessageCountTokensInputTokensKey            = attribute.Key("anthropic.message.count_tokens.input_tokens")
	AnthropicRequestServiceTierKey                       = attribute.Key("anthropic.request.service_tier")
	AnthropicResponseServiceTierKey                      = attribute.Key("anthropic.response.service_tier")
	AnthropicSessionThreadIDKey                          = attribute.Key("anthropic.session.thread.id")
	AnthropicThinkingBudgetTokensKey                     = attribute.Key("anthropic.thinking.budget_tokens")
	AnthropicThinkingTypeKey                             = attribute.Key("anthropic.thinking.type")
	AnthropicToolResultContentKey                        = attribute.Key("anthropic.tool_result.content")
	AnthropicToolsetNameKey                              = attribute.Key("anthropic.toolset.name")
	AnthropicUsageCacheCreationEphemeral1hInputTokensKey = attribute.Key("anthropic.usage.cache_creation.ephemeral_1h_input_tokens")
	AnthropicUsageCacheCreationEphemeral5mInputTokensKey = attribute.Key("anthropic.usage.cache_creation.ephemeral_5m_input_tokens")
	AnthropicUsageCompactionInputTokensKey               = attribute.Key("anthropic.usage.compaction.input_tokens")
	AnthropicUsageCompactionOutputTokensKey              = attribute.Key("anthropic.usage.compaction.output_tokens")
)
