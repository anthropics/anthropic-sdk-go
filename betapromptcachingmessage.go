// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"context"
	"net/http"
	"reflect"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/param"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
	"github.com/tidwall/gjson"
)

// BetaPromptCachingMessageService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaPromptCachingMessageService] method instead.
type BetaPromptCachingMessageService struct {
	Options []option.RequestOption
}

// NewBetaPromptCachingMessageService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaPromptCachingMessageService(opts ...option.RequestOption) (r *BetaPromptCachingMessageService) {
	r = &BetaPromptCachingMessageService{}
	r.Options = opts
	return
}

// Create a Message.
//
// Send a structured list of input messages with text and/or image content, and the
// model will generate the next message in the conversation.
//
// The Messages API can be used for either single queries or stateless multi-turn
// conversations.
//
// Note: If you choose to set a timeout for this request, we recommend 10 minutes.
func (r *BetaPromptCachingMessageService) New(ctx context.Context, body BetaPromptCachingMessageNewParams, opts ...option.RequestOption) (res *PromptCachingBetaMessage, err error) {
	opts = append(r.Options[:], opts...)
	path := "v1/messages?beta=prompt_caching"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Create a Message.
//
// Send a structured list of input messages with text and/or image content, and the
// model will generate the next message in the conversation.
//
// The Messages API can be used for either single queries or stateless multi-turn
// conversations.
//
// Note: If you choose to set a timeout for this request, we recommend 10 minutes.
func (r *BetaPromptCachingMessageService) NewStreaming(ctx context.Context, body BetaPromptCachingMessageNewParams, opts ...option.RequestOption) (stream *ssestream.Stream[RawPromptCachingBetaMessageStreamEvent]) {
	var (
		raw *http.Response
		err error
	)
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithJSONSet("stream", true)}, opts...)
	path := "v1/messages?beta=prompt_caching"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &raw, opts...)
	return ssestream.NewStream[RawPromptCachingBetaMessageStreamEvent](ssestream.NewDecoder(raw), err)
}

type PromptCachingBetaCacheControlEphemeralParam struct {
	Type param.Field[PromptCachingBetaCacheControlEphemeralType] `json:"type,required"`
}

func (r PromptCachingBetaCacheControlEphemeralParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type PromptCachingBetaCacheControlEphemeralType string

const (
	PromptCachingBetaCacheControlEphemeralTypeEphemeral PromptCachingBetaCacheControlEphemeralType = "ephemeral"
)

func (r PromptCachingBetaCacheControlEphemeralType) IsKnown() bool {
	switch r {
	case PromptCachingBetaCacheControlEphemeralTypeEphemeral:
		return true
	}
	return false
}

type PromptCachingBetaImageBlockParam struct {
	Source       param.Field[PromptCachingBetaImageBlockParamSource]      `json:"source,required"`
	Type         param.Field[PromptCachingBetaImageBlockParamType]        `json:"type,required"`
	CacheControl param.Field[PromptCachingBetaCacheControlEphemeralParam] `json:"cache_control"`
}

func (r PromptCachingBetaImageBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r PromptCachingBetaImageBlockParam) implementsPromptCachingBetaMessageParamContentUnion() {}

func (r PromptCachingBetaImageBlockParam) implementsPromptCachingBetaToolResultBlockParamContentArrayUnionItem() {
}

type PromptCachingBetaImageBlockParamSource struct {
	Data      param.Field[string]                                          `json:"data,required" format:"byte"`
	MediaType param.Field[PromptCachingBetaImageBlockParamSourceMediaType] `json:"media_type,required"`
	Type      param.Field[PromptCachingBetaImageBlockParamSourceType]      `json:"type,required"`
}

func (r PromptCachingBetaImageBlockParamSource) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type PromptCachingBetaImageBlockParamSourceMediaType string

const (
	PromptCachingBetaImageBlockParamSourceMediaTypeImageJPEG PromptCachingBetaImageBlockParamSourceMediaType = "image/jpeg"
	PromptCachingBetaImageBlockParamSourceMediaTypeImagePNG  PromptCachingBetaImageBlockParamSourceMediaType = "image/png"
	PromptCachingBetaImageBlockParamSourceMediaTypeImageGIF  PromptCachingBetaImageBlockParamSourceMediaType = "image/gif"
	PromptCachingBetaImageBlockParamSourceMediaTypeImageWebP PromptCachingBetaImageBlockParamSourceMediaType = "image/webp"
)

func (r PromptCachingBetaImageBlockParamSourceMediaType) IsKnown() bool {
	switch r {
	case PromptCachingBetaImageBlockParamSourceMediaTypeImageJPEG, PromptCachingBetaImageBlockParamSourceMediaTypeImagePNG, PromptCachingBetaImageBlockParamSourceMediaTypeImageGIF, PromptCachingBetaImageBlockParamSourceMediaTypeImageWebP:
		return true
	}
	return false
}

type PromptCachingBetaImageBlockParamSourceType string

const (
	PromptCachingBetaImageBlockParamSourceTypeBase64 PromptCachingBetaImageBlockParamSourceType = "base64"
)

func (r PromptCachingBetaImageBlockParamSourceType) IsKnown() bool {
	switch r {
	case PromptCachingBetaImageBlockParamSourceTypeBase64:
		return true
	}
	return false
}

type PromptCachingBetaImageBlockParamType string

const (
	PromptCachingBetaImageBlockParamTypeImage PromptCachingBetaImageBlockParamType = "image"
)

func (r PromptCachingBetaImageBlockParamType) IsKnown() bool {
	switch r {
	case PromptCachingBetaImageBlockParamTypeImage:
		return true
	}
	return false
}

type PromptCachingBetaMessage struct {
	// Unique object identifier.
	//
	// The format and length of IDs may change over time.
	ID string `json:"id,required"`
	// Content generated by the model.
	//
	// This is an array of content blocks, each of which has a `type` that determines
	// its shape.
	//
	// Example:
	//
	// ```json
	// [{ "type": "text", "text": "Hi, I'm Claude." }]
	// ```
	//
	// If the request input `messages` ended with an `assistant` turn, then the
	// response `content` will continue directly from that last turn. You can use this
	// to constrain the model's output.
	//
	// For example, if the input `messages` were:
	//
	// ```json
	// [
	//
	//	{
	//	  "role": "user",
	//	  "content": "What's the Greek name for Sun? (A) Sol (B) Helios (C) Sun"
	//	},
	//	{ "role": "assistant", "content": "The best answer is (" }
	//
	// ]
	// ```
	//
	// Then the response `content` might be:
	//
	// ```json
	// [{ "type": "text", "text": "B)" }]
	// ```
	Content []ContentBlock `json:"content,required"`
	// The model that will complete your prompt.\n\nSee
	// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	Model Model `json:"model,required"`
	// Conversational role of the generated message.
	//
	// This will always be `"assistant"`.
	Role PromptCachingBetaMessageRole `json:"role,required"`
	// The reason that we stopped.
	//
	// This may be one the following values:
	//
	// - `"end_turn"`: the model reached a natural stopping point
	// - `"max_tokens"`: we exceeded the requested `max_tokens` or the model's maximum
	// - `"stop_sequence"`: one of your provided custom `stop_sequences` was generated
	// - `"tool_use"`: the model invoked one or more tools
	//
	// In non-streaming mode this value is always non-null. In streaming mode, it is
	// null in the `message_start` event and non-null otherwise.
	StopReason PromptCachingBetaMessageStopReason `json:"stop_reason,required,nullable"`
	// Which custom stop sequence was generated, if any.
	//
	// This value will be a non-null string if one of your custom stop sequences was
	// generated.
	StopSequence string `json:"stop_sequence,required,nullable"`
	// Object type.
	//
	// For Messages, this is always `"message"`.
	Type PromptCachingBetaMessageType `json:"type,required"`
	// Billing and rate-limit usage.
	//
	// Anthropic's API bills and rate-limits by token counts, as tokens represent the
	// underlying cost to our systems.
	//
	// Under the hood, the API transforms requests into a format suitable for the
	// model. The model's output then goes through a parsing stage before becoming an
	// API response. As a result, the token counts in `usage` will not match one-to-one
	// with the exact visible content of an API request or response.
	//
	// For example, `output_tokens` will be non-zero, even for an empty string response
	// from Claude.
	Usage PromptCachingBetaUsage       `json:"usage,required"`
	JSON  promptCachingBetaMessageJSON `json:"-"`
}

// promptCachingBetaMessageJSON contains the JSON metadata for the struct
// [PromptCachingBetaMessage]
type promptCachingBetaMessageJSON struct {
	ID           apijson.Field
	Content      apijson.Field
	Model        apijson.Field
	Role         apijson.Field
	StopReason   apijson.Field
	StopSequence apijson.Field
	Type         apijson.Field
	Usage        apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *PromptCachingBetaMessage) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r promptCachingBetaMessageJSON) RawJSON() string {
	return r.raw
}

// Conversational role of the generated message.
//
// This will always be `"assistant"`.
type PromptCachingBetaMessageRole string

const (
	PromptCachingBetaMessageRoleAssistant PromptCachingBetaMessageRole = "assistant"
)

func (r PromptCachingBetaMessageRole) IsKnown() bool {
	switch r {
	case PromptCachingBetaMessageRoleAssistant:
		return true
	}
	return false
}

// The reason that we stopped.
//
// This may be one the following values:
//
// - `"end_turn"`: the model reached a natural stopping point
// - `"max_tokens"`: we exceeded the requested `max_tokens` or the model's maximum
// - `"stop_sequence"`: one of your provided custom `stop_sequences` was generated
// - `"tool_use"`: the model invoked one or more tools
//
// In non-streaming mode this value is always non-null. In streaming mode, it is
// null in the `message_start` event and non-null otherwise.
type PromptCachingBetaMessageStopReason string

const (
	PromptCachingBetaMessageStopReasonEndTurn      PromptCachingBetaMessageStopReason = "end_turn"
	PromptCachingBetaMessageStopReasonMaxTokens    PromptCachingBetaMessageStopReason = "max_tokens"
	PromptCachingBetaMessageStopReasonStopSequence PromptCachingBetaMessageStopReason = "stop_sequence"
	PromptCachingBetaMessageStopReasonToolUse      PromptCachingBetaMessageStopReason = "tool_use"
)

func (r PromptCachingBetaMessageStopReason) IsKnown() bool {
	switch r {
	case PromptCachingBetaMessageStopReasonEndTurn, PromptCachingBetaMessageStopReasonMaxTokens, PromptCachingBetaMessageStopReasonStopSequence, PromptCachingBetaMessageStopReasonToolUse:
		return true
	}
	return false
}

// Object type.
//
// For Messages, this is always `"message"`.
type PromptCachingBetaMessageType string

const (
	PromptCachingBetaMessageTypeMessage PromptCachingBetaMessageType = "message"
)

func (r PromptCachingBetaMessageType) IsKnown() bool {
	switch r {
	case PromptCachingBetaMessageTypeMessage:
		return true
	}
	return false
}

type PromptCachingBetaMessageParam struct {
	Content param.Field[[]PromptCachingBetaMessageParamContentUnion] `json:"content,required"`
	Role    param.Field[PromptCachingBetaMessageParamRole]           `json:"role,required"`
}

func (r PromptCachingBetaMessageParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type PromptCachingBetaMessageParamContent struct {
	CacheControl param.Field[PromptCachingBetaCacheControlEphemeralParam] `json:"cache_control"`
	Type         param.Field[PromptCachingBetaMessageParamContentType]    `json:"type,required"`
	Text         param.Field[string]                                      `json:"text"`
	Source       param.Field[interface{}]                                 `json:"source,required"`
	ID           param.Field[string]                                      `json:"id"`
	Name         param.Field[string]                                      `json:"name"`
	Input        param.Field[interface{}]                                 `json:"input,required"`
	ToolUseID    param.Field[string]                                      `json:"tool_use_id"`
	IsError      param.Field[bool]                                        `json:"is_error"`
	Content      param.Field[interface{}]                                 `json:"content,required"`
}

func (r PromptCachingBetaMessageParamContent) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r PromptCachingBetaMessageParamContent) implementsPromptCachingBetaMessageParamContentUnion() {}

// Satisfied by [PromptCachingBetaTextBlockParam],
// [PromptCachingBetaImageBlockParam], [PromptCachingBetaToolUseBlockParam],
// [PromptCachingBetaToolResultBlockParam], [PromptCachingBetaMessageParamContent].
type PromptCachingBetaMessageParamContentUnion interface {
	implementsPromptCachingBetaMessageParamContentUnion()
}

type PromptCachingBetaMessageParamContentType string

const (
	PromptCachingBetaMessageParamContentTypeText       PromptCachingBetaMessageParamContentType = "text"
	PromptCachingBetaMessageParamContentTypeImage      PromptCachingBetaMessageParamContentType = "image"
	PromptCachingBetaMessageParamContentTypeToolUse    PromptCachingBetaMessageParamContentType = "tool_use"
	PromptCachingBetaMessageParamContentTypeToolResult PromptCachingBetaMessageParamContentType = "tool_result"
)

func (r PromptCachingBetaMessageParamContentType) IsKnown() bool {
	switch r {
	case PromptCachingBetaMessageParamContentTypeText, PromptCachingBetaMessageParamContentTypeImage, PromptCachingBetaMessageParamContentTypeToolUse, PromptCachingBetaMessageParamContentTypeToolResult:
		return true
	}
	return false
}

type PromptCachingBetaMessageParamRole string

const (
	PromptCachingBetaMessageParamRoleUser      PromptCachingBetaMessageParamRole = "user"
	PromptCachingBetaMessageParamRoleAssistant PromptCachingBetaMessageParamRole = "assistant"
)

func (r PromptCachingBetaMessageParamRole) IsKnown() bool {
	switch r {
	case PromptCachingBetaMessageParamRoleUser, PromptCachingBetaMessageParamRoleAssistant:
		return true
	}
	return false
}

type PromptCachingBetaTextBlockParam struct {
	Text         param.Field[string]                                      `json:"text,required"`
	Type         param.Field[PromptCachingBetaTextBlockParamType]         `json:"type,required"`
	CacheControl param.Field[PromptCachingBetaCacheControlEphemeralParam] `json:"cache_control"`
}

func (r PromptCachingBetaTextBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r PromptCachingBetaTextBlockParam) implementsPromptCachingBetaMessageParamContentUnion() {}

func (r PromptCachingBetaTextBlockParam) implementsPromptCachingBetaToolResultBlockParamContentArrayUnionItem() {
}

type PromptCachingBetaTextBlockParamType string

const (
	PromptCachingBetaTextBlockParamTypeText PromptCachingBetaTextBlockParamType = "text"
)

func (r PromptCachingBetaTextBlockParamType) IsKnown() bool {
	switch r {
	case PromptCachingBetaTextBlockParamTypeText:
		return true
	}
	return false
}

type PromptCachingBetaToolParam struct {
	// [JSON schema](https://json-schema.org/) for this tool's input.
	//
	// This defines the shape of the `input` that your tool accepts and that the model
	// will produce.
	InputSchema  param.Field[PromptCachingBetaToolInputSchemaParam]       `json:"input_schema,required"`
	Name         param.Field[string]                                      `json:"name,required"`
	CacheControl param.Field[PromptCachingBetaCacheControlEphemeralParam] `json:"cache_control"`
	// Description of what this tool does.
	//
	// Tool descriptions should be as detailed as possible. The more information that
	// the model has about what the tool is and how to use it, the better it will
	// perform. You can use natural language descriptions to reinforce important
	// aspects of the tool input JSON schema.
	Description param.Field[string] `json:"description"`
}

func (r PromptCachingBetaToolParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// [JSON schema](https://json-schema.org/) for this tool's input.
//
// This defines the shape of the `input` that your tool accepts and that the model
// will produce.
type PromptCachingBetaToolInputSchemaParam struct {
	Type        param.Field[PromptCachingBetaToolInputSchemaType] `json:"type,required"`
	Properties  param.Field[interface{}]                          `json:"properties"`
	ExtraFields map[string]interface{}                            `json:"-,extras"`
}

func (r PromptCachingBetaToolInputSchemaParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type PromptCachingBetaToolInputSchemaType string

const (
	PromptCachingBetaToolInputSchemaTypeObject PromptCachingBetaToolInputSchemaType = "object"
)

func (r PromptCachingBetaToolInputSchemaType) IsKnown() bool {
	switch r {
	case PromptCachingBetaToolInputSchemaTypeObject:
		return true
	}
	return false
}

type PromptCachingBetaToolResultBlockParam struct {
	ToolUseID    param.Field[string]                                            `json:"tool_use_id,required"`
	Type         param.Field[PromptCachingBetaToolResultBlockParamType]         `json:"type,required"`
	CacheControl param.Field[PromptCachingBetaCacheControlEphemeralParam]       `json:"cache_control"`
	Content      param.Field[PromptCachingBetaToolResultBlockParamContentUnion] `json:"content"`
	IsError      param.Field[bool]                                              `json:"is_error"`
}

func (r PromptCachingBetaToolResultBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r PromptCachingBetaToolResultBlockParam) implementsPromptCachingBetaMessageParamContentUnion() {
}

type PromptCachingBetaToolResultBlockParamType string

const (
	PromptCachingBetaToolResultBlockParamTypeToolResult PromptCachingBetaToolResultBlockParamType = "tool_result"
)

func (r PromptCachingBetaToolResultBlockParamType) IsKnown() bool {
	switch r {
	case PromptCachingBetaToolResultBlockParamTypeToolResult:
		return true
	}
	return false
}

// Satisfied by [shared.UnionString],
// [PromptCachingBetaToolResultBlockParamContentArray].
type PromptCachingBetaToolResultBlockParamContentUnion interface {
	ImplementsPromptCachingBetaToolResultBlockParamContentUnion()
}

type PromptCachingBetaToolResultBlockParamContentArray []PromptCachingBetaToolResultBlockParamContentArrayUnionItem

func (r PromptCachingBetaToolResultBlockParamContentArray) ImplementsPromptCachingBetaToolResultBlockParamContentUnion() {
}

type PromptCachingBetaToolResultBlockParamContentArrayItem struct {
	CacheControl param.Field[PromptCachingBetaCacheControlEphemeralParam]           `json:"cache_control"`
	Type         param.Field[PromptCachingBetaToolResultBlockParamContentArrayType] `json:"type,required"`
	Text         param.Field[string]                                                `json:"text"`
	Source       param.Field[interface{}]                                           `json:"source,required"`
}

func (r PromptCachingBetaToolResultBlockParamContentArrayItem) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r PromptCachingBetaToolResultBlockParamContentArrayItem) implementsPromptCachingBetaToolResultBlockParamContentArrayUnionItem() {
}

// Satisfied by [PromptCachingBetaTextBlockParam],
// [PromptCachingBetaImageBlockParam],
// [PromptCachingBetaToolResultBlockParamContentArrayItem].
type PromptCachingBetaToolResultBlockParamContentArrayUnionItem interface {
	implementsPromptCachingBetaToolResultBlockParamContentArrayUnionItem()
}

type PromptCachingBetaToolResultBlockParamContentArrayType string

const (
	PromptCachingBetaToolResultBlockParamContentArrayTypeText  PromptCachingBetaToolResultBlockParamContentArrayType = "text"
	PromptCachingBetaToolResultBlockParamContentArrayTypeImage PromptCachingBetaToolResultBlockParamContentArrayType = "image"
)

func (r PromptCachingBetaToolResultBlockParamContentArrayType) IsKnown() bool {
	switch r {
	case PromptCachingBetaToolResultBlockParamContentArrayTypeText, PromptCachingBetaToolResultBlockParamContentArrayTypeImage:
		return true
	}
	return false
}

type PromptCachingBetaToolUseBlockParam struct {
	ID           param.Field[string]                                      `json:"id,required"`
	Input        param.Field[interface{}]                                 `json:"input,required"`
	Name         param.Field[string]                                      `json:"name,required"`
	Type         param.Field[PromptCachingBetaToolUseBlockParamType]      `json:"type,required"`
	CacheControl param.Field[PromptCachingBetaCacheControlEphemeralParam] `json:"cache_control"`
}

func (r PromptCachingBetaToolUseBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r PromptCachingBetaToolUseBlockParam) implementsPromptCachingBetaMessageParamContentUnion() {}

type PromptCachingBetaToolUseBlockParamType string

const (
	PromptCachingBetaToolUseBlockParamTypeToolUse PromptCachingBetaToolUseBlockParamType = "tool_use"
)

func (r PromptCachingBetaToolUseBlockParamType) IsKnown() bool {
	switch r {
	case PromptCachingBetaToolUseBlockParamTypeToolUse:
		return true
	}
	return false
}

type PromptCachingBetaUsage struct {
	// The number of input tokens used to create the cache entry.
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens,required,nullable"`
	// The number of input tokens read from the cache.
	CacheReadInputTokens int64 `json:"cache_read_input_tokens,required,nullable"`
	// The number of input tokens which were used.
	InputTokens int64 `json:"input_tokens,required"`
	// The number of output tokens which were used.
	OutputTokens int64                      `json:"output_tokens,required"`
	JSON         promptCachingBetaUsageJSON `json:"-"`
}

// promptCachingBetaUsageJSON contains the JSON metadata for the struct
// [PromptCachingBetaUsage]
type promptCachingBetaUsageJSON struct {
	CacheCreationInputTokens apijson.Field
	CacheReadInputTokens     apijson.Field
	InputTokens              apijson.Field
	OutputTokens             apijson.Field
	raw                      string
	ExtraFields              map[string]apijson.Field
}

func (r *PromptCachingBetaUsage) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r promptCachingBetaUsageJSON) RawJSON() string {
	return r.raw
}

type RawPromptCachingBetaMessageStartEvent struct {
	Message PromptCachingBetaMessage                  `json:"message,required"`
	Type    RawPromptCachingBetaMessageStartEventType `json:"type,required"`
	JSON    rawPromptCachingBetaMessageStartEventJSON `json:"-"`
}

// rawPromptCachingBetaMessageStartEventJSON contains the JSON metadata for the
// struct [RawPromptCachingBetaMessageStartEvent]
type rawPromptCachingBetaMessageStartEventJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RawPromptCachingBetaMessageStartEvent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r rawPromptCachingBetaMessageStartEventJSON) RawJSON() string {
	return r.raw
}

func (r RawPromptCachingBetaMessageStartEvent) implementsRawPromptCachingBetaMessageStreamEvent() {}

type RawPromptCachingBetaMessageStartEventType string

const (
	RawPromptCachingBetaMessageStartEventTypeMessageStart RawPromptCachingBetaMessageStartEventType = "message_start"
)

func (r RawPromptCachingBetaMessageStartEventType) IsKnown() bool {
	switch r {
	case RawPromptCachingBetaMessageStartEventTypeMessageStart:
		return true
	}
	return false
}

type RawPromptCachingBetaMessageStreamEvent struct {
	Type    RawPromptCachingBetaMessageStreamEventType `json:"type,required"`
	Message PromptCachingBetaMessage                   `json:"message"`
	// This field can have the runtime type of [MessageDeltaEventDelta],
	// [ContentBlockDeltaEventDelta].
	Delta interface{} `json:"delta,required"`
	// Billing and rate-limit usage.
	//
	// Anthropic's API bills and rate-limits by token counts, as tokens represent the
	// underlying cost to our systems.
	//
	// Under the hood, the API transforms requests into a format suitable for the
	// model. The model's output then goes through a parsing stage before becoming an
	// API response. As a result, the token counts in `usage` will not match one-to-one
	// with the exact visible content of an API request or response.
	//
	// For example, `output_tokens` will be non-zero, even for an empty string response
	// from Claude.
	Usage MessageDeltaUsage `json:"usage"`
	Index int64             `json:"index"`
	// This field can have the runtime type of [ContentBlockStartEventContentBlock].
	ContentBlock interface{}                                `json:"content_block,required"`
	JSON         rawPromptCachingBetaMessageStreamEventJSON `json:"-"`
	union        RawPromptCachingBetaMessageStreamEventUnion
}

// rawPromptCachingBetaMessageStreamEventJSON contains the JSON metadata for the
// struct [RawPromptCachingBetaMessageStreamEvent]
type rawPromptCachingBetaMessageStreamEventJSON struct {
	Type         apijson.Field
	Message      apijson.Field
	Delta        apijson.Field
	Usage        apijson.Field
	Index        apijson.Field
	ContentBlock apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r rawPromptCachingBetaMessageStreamEventJSON) RawJSON() string {
	return r.raw
}

func (r *RawPromptCachingBetaMessageStreamEvent) UnmarshalJSON(data []byte) (err error) {
	*r = RawPromptCachingBetaMessageStreamEvent{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [RawPromptCachingBetaMessageStreamEventUnion] interface which
// you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are [RawPromptCachingBetaMessageStartEvent],
// [MessageDeltaEvent], [MessageStopEvent], [ContentBlockStartEvent],
// [ContentBlockDeltaEvent], [ContentBlockStopEvent].
func (r RawPromptCachingBetaMessageStreamEvent) AsUnion() RawPromptCachingBetaMessageStreamEventUnion {
	return r.union
}

// Union satisfied by [RawPromptCachingBetaMessageStartEvent], [MessageDeltaEvent],
// [MessageStopEvent], [ContentBlockStartEvent], [ContentBlockDeltaEvent] or
// [ContentBlockStopEvent].
type RawPromptCachingBetaMessageStreamEventUnion interface {
	implementsRawPromptCachingBetaMessageStreamEvent()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*RawPromptCachingBetaMessageStreamEventUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(RawPromptCachingBetaMessageStartEvent{}),
			DiscriminatorValue: "message_start",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(MessageDeltaEvent{}),
			DiscriminatorValue: "message_delta",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(MessageStopEvent{}),
			DiscriminatorValue: "message_stop",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ContentBlockStartEvent{}),
			DiscriminatorValue: "content_block_start",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ContentBlockDeltaEvent{}),
			DiscriminatorValue: "content_block_delta",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ContentBlockStopEvent{}),
			DiscriminatorValue: "content_block_stop",
		},
	)
}

type RawPromptCachingBetaMessageStreamEventType string

const (
	RawPromptCachingBetaMessageStreamEventTypeMessageStart      RawPromptCachingBetaMessageStreamEventType = "message_start"
	RawPromptCachingBetaMessageStreamEventTypeMessageDelta      RawPromptCachingBetaMessageStreamEventType = "message_delta"
	RawPromptCachingBetaMessageStreamEventTypeMessageStop       RawPromptCachingBetaMessageStreamEventType = "message_stop"
	RawPromptCachingBetaMessageStreamEventTypeContentBlockStart RawPromptCachingBetaMessageStreamEventType = "content_block_start"
	RawPromptCachingBetaMessageStreamEventTypeContentBlockDelta RawPromptCachingBetaMessageStreamEventType = "content_block_delta"
	RawPromptCachingBetaMessageStreamEventTypeContentBlockStop  RawPromptCachingBetaMessageStreamEventType = "content_block_stop"
)

func (r RawPromptCachingBetaMessageStreamEventType) IsKnown() bool {
	switch r {
	case RawPromptCachingBetaMessageStreamEventTypeMessageStart, RawPromptCachingBetaMessageStreamEventTypeMessageDelta, RawPromptCachingBetaMessageStreamEventTypeMessageStop, RawPromptCachingBetaMessageStreamEventTypeContentBlockStart, RawPromptCachingBetaMessageStreamEventTypeContentBlockDelta, RawPromptCachingBetaMessageStreamEventTypeContentBlockStop:
		return true
	}
	return false
}

type BetaPromptCachingMessageNewParams struct {
	// The maximum number of tokens to generate before stopping.
	//
	// Note that our models may stop _before_ reaching this maximum. This parameter
	// only specifies the absolute maximum number of tokens to generate.
	//
	// Different models have different maximum values for this parameter. See
	// [models](https://docs.anthropic.com/en/docs/models-overview) for details.
	MaxTokens param.Field[int64] `json:"max_tokens,required"`
	// Input messages.
	//
	// Our models are trained to operate on alternating `user` and `assistant`
	// conversational turns. When creating a new `Message`, you specify the prior
	// conversational turns with the `messages` parameter, and the model then generates
	// the next `Message` in the conversation.
	//
	// Each input message must be an object with a `role` and `content`. You can
	// specify a single `user`-role message, or you can include multiple `user` and
	// `assistant` messages. The first message must always use the `user` role.
	//
	// If the final message uses the `assistant` role, the response content will
	// continue immediately from the content in that message. This can be used to
	// constrain part of the model's response.
	//
	// Example with a single `user` message:
	//
	// ```json
	// [{ "role": "user", "content": "Hello, Claude" }]
	// ```
	//
	// Example with multiple conversational turns:
	//
	// ```json
	// [
	//
	//	{ "role": "user", "content": "Hello there." },
	//	{ "role": "assistant", "content": "Hi, I'm Claude. How can I help you?" },
	//	{ "role": "user", "content": "Can you explain LLMs in plain English?" }
	//
	// ]
	// ```
	//
	// Example with a partially-filled response from Claude:
	//
	// ```json
	// [
	//
	//	{
	//	  "role": "user",
	//	  "content": "What's the Greek name for Sun? (A) Sol (B) Helios (C) Sun"
	//	},
	//	{ "role": "assistant", "content": "The best answer is (" }
	//
	// ]
	// ```
	//
	// Each input message `content` may be either a single `string` or an array of
	// content blocks, where each block has a specific `type`. Using a `string` for
	// `content` is shorthand for an array of one content block of type `"text"`. The
	// following input messages are equivalent:
	//
	// ```json
	// { "role": "user", "content": "Hello, Claude" }
	// ```
	//
	// ```json
	// { "role": "user", "content": [{ "type": "text", "text": "Hello, Claude" }] }
	// ```
	//
	// Starting with Claude 3 models, you can also send image content blocks:
	//
	// ```json
	//
	//	{
	//	  "role": "user",
	//	  "content": [
	//	    {
	//	      "type": "image",
	//	      "source": {
	//	        "type": "base64",
	//	        "media_type": "image/jpeg",
	//	        "data": "/9j/4AAQSkZJRg..."
	//	      }
	//	    },
	//	    { "type": "text", "text": "What is in this image?" }
	//	  ]
	//	}
	//
	// ```
	//
	// We currently support the `base64` source type for images, and the `image/jpeg`,
	// `image/png`, `image/gif`, and `image/webp` media types.
	//
	// See [examples](https://docs.anthropic.com/en/api/messages-examples#vision) for
	// more input examples.
	//
	// Note that if you want to include a
	// [system prompt](https://docs.anthropic.com/en/docs/system-prompts), you can use
	// the top-level `system` parameter — there is no `"system"` role for input
	// messages in the Messages API.
	Messages param.Field[[]PromptCachingBetaMessageParam] `json:"messages,required"`
	// The model that will complete your prompt.\n\nSee
	// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	Model param.Field[Model] `json:"model,required"`
	// An object describing metadata about the request.
	Metadata param.Field[BetaPromptCachingMessageNewParamsMetadata] `json:"metadata"`
	// Custom text sequences that will cause the model to stop generating.
	//
	// Our models will normally stop when they have naturally completed their turn,
	// which will result in a response `stop_reason` of `"end_turn"`.
	//
	// If you want the model to stop generating when it encounters custom strings of
	// text, you can use the `stop_sequences` parameter. If the model encounters one of
	// the custom sequences, the response `stop_reason` value will be `"stop_sequence"`
	// and the response `stop_sequence` value will contain the matched stop sequence.
	StopSequences param.Field[[]string] `json:"stop_sequences"`
	// System prompt.
	//
	// A system prompt is a way of providing context and instructions to Claude, such
	// as specifying a particular goal or role. See our
	// [guide to system prompts](https://docs.anthropic.com/en/docs/system-prompts).
	System param.Field[BetaPromptCachingMessageNewParamsSystemUnion] `json:"system"`
	// Amount of randomness injected into the response.
	//
	// Defaults to `1.0`. Ranges from `0.0` to `1.0`. Use `temperature` closer to `0.0`
	// for analytical / multiple choice, and closer to `1.0` for creative and
	// generative tasks.
	//
	// Note that even with `temperature` of `0.0`, the results will not be fully
	// deterministic.
	Temperature param.Field[float64] `json:"temperature"`
	// How the model should use the provided tools. The model can use a specific tool,
	// any available tool, or decide by itself.
	ToolChoice param.Field[ToolChoiceUnionParam] `json:"tool_choice"`
	// Definitions of tools that the model may use.
	//
	// If you include `tools` in your API request, the model may return `tool_use`
	// content blocks that represent the model's use of those tools. You can then run
	// those tools using the tool input generated by the model and then optionally
	// return results back to the model using `tool_result` content blocks.
	//
	// Each tool definition includes:
	//
	//   - `name`: Name of the tool.
	//   - `description`: Optional, but strongly-recommended description of the tool.
	//   - `input_schema`: [JSON schema](https://json-schema.org/) for the tool `input`
	//     shape that the model will produce in `tool_use` output content blocks.
	//
	// For example, if you defined `tools` as:
	//
	// ```json
	// [
	//
	//	{
	//	  "name": "get_stock_price",
	//	  "description": "Get the current stock price for a given ticker symbol.",
	//	  "input_schema": {
	//	    "type": "object",
	//	    "properties": {
	//	      "ticker": {
	//	        "type": "string",
	//	        "description": "The stock ticker symbol, e.g. AAPL for Apple Inc."
	//	      }
	//	    },
	//	    "required": ["ticker"]
	//	  }
	//	}
	//
	// ]
	// ```
	//
	// And then asked the model "What's the S&P 500 at today?", the model might produce
	// `tool_use` content blocks in the response like this:
	//
	// ```json
	// [
	//
	//	{
	//	  "type": "tool_use",
	//	  "id": "toolu_01D7FLrfh4GYq7yT1ULFeyMV",
	//	  "name": "get_stock_price",
	//	  "input": { "ticker": "^GSPC" }
	//	}
	//
	// ]
	// ```
	//
	// You might then run your `get_stock_price` tool with `{"ticker": "^GSPC"}` as an
	// input, and return the following back to the model in a subsequent `user`
	// message:
	//
	// ```json
	// [
	//
	//	{
	//	  "type": "tool_result",
	//	  "tool_use_id": "toolu_01D7FLrfh4GYq7yT1ULFeyMV",
	//	  "content": "259.75 USD"
	//	}
	//
	// ]
	// ```
	//
	// Tools can be used for workflows that include running client-side tools and
	// functions, or more generally whenever you want the model to produce a particular
	// JSON structure of output.
	//
	// See our [guide](https://docs.anthropic.com/en/docs/tool-use) for more details.
	Tools param.Field[[]PromptCachingBetaToolParam] `json:"tools"`
	// Only sample from the top K options for each subsequent token.
	//
	// Used to remove "long tail" low probability responses.
	// [Learn more technical details here](https://towardsdatascience.com/how-to-sample-from-language-models-682bceb97277).
	//
	// Recommended for advanced use cases only. You usually only need to use
	// `temperature`.
	TopK param.Field[int64] `json:"top_k"`
	// Use nucleus sampling.
	//
	// In nucleus sampling, we compute the cumulative distribution over all the options
	// for each subsequent token in decreasing probability order and cut it off once it
	// reaches a particular probability specified by `top_p`. You should either alter
	// `temperature` or `top_p`, but not both.
	//
	// Recommended for advanced use cases only. You usually only need to use
	// `temperature`.
	TopP param.Field[float64] `json:"top_p"`
}

func (r BetaPromptCachingMessageNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// An object describing metadata about the request.
type BetaPromptCachingMessageNewParamsMetadata struct {
	// An external identifier for the user who is associated with the request.
	//
	// This should be a uuid, hash value, or other opaque identifier. Anthropic may use
	// this id to help detect abuse. Do not include any identifying information such as
	// name, email address, or phone number.
	UserID param.Field[string] `json:"user_id"`
}

func (r BetaPromptCachingMessageNewParamsMetadata) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// System prompt.
//
// A system prompt is a way of providing context and instructions to Claude, such
// as specifying a particular goal or role. See our
// [guide to system prompts](https://docs.anthropic.com/en/docs/system-prompts).
//
// Satisfied by [shared.UnionString],
// [BetaPromptCachingMessageNewParamsSystemArray].
type BetaPromptCachingMessageNewParamsSystemUnion interface {
	ImplementsBetaPromptCachingMessageNewParamsSystemUnion()
}

type BetaPromptCachingMessageNewParamsSystemArray []PromptCachingBetaTextBlockParam

func (r BetaPromptCachingMessageNewParamsSystemArray) ImplementsBetaPromptCachingMessageNewParamsSystemUnion() {
}
