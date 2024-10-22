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

// BetaMessageService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaMessageService] method instead.
type BetaMessageService struct {
	Options []option.RequestOption
	Batches *BetaMessageBatchService
}

// NewBetaMessageService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaMessageService(opts ...option.RequestOption) (r *BetaMessageService) {
	r = &BetaMessageService{}
	r.Options = opts
	r.Batches = NewBetaMessageBatchService(opts...)
	return
}

// Send a structured list of input messages with text and/or image content, and the
// model will generate the next message in the conversation.
//
// The Messages API can be used for either single queries or stateless multi-turn
// conversations.
//
// Note: If you choose to set a timeout for this request, we recommend 10 minutes.
func (r *BetaMessageService) New(ctx context.Context, params BetaMessageNewParams, opts ...option.RequestOption) (res *BetaMessage, err error) {
	opts = append(r.Options[:], opts...)
	path := "v1/messages?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return
}

// Send a structured list of input messages with text and/or image content, and the
// model will generate the next message in the conversation.
//
// The Messages API can be used for either single queries or stateless multi-turn
// conversations.
//
// Note: If you choose to set a timeout for this request, we recommend 10 minutes.
func (r *BetaMessageService) NewStreaming(ctx context.Context, params BetaMessageNewParams, opts ...option.RequestOption) (stream *ssestream.Stream[BetaRawMessageStreamEvent]) {
	var (
		raw *http.Response
		err error
	)
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithJSONSet("stream", true)}, opts...)
	path := "v1/messages?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &raw, opts...)
	return ssestream.NewStream[BetaRawMessageStreamEvent](ssestream.NewDecoder(raw), err)
}

type BetaCacheControlEphemeralParam struct {
	Type param.Field[BetaCacheControlEphemeralType] `json:"type,required"`
}

func (r BetaCacheControlEphemeralParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaCacheControlEphemeralType string

const (
	BetaCacheControlEphemeralTypeEphemeral BetaCacheControlEphemeralType = "ephemeral"
)

func (r BetaCacheControlEphemeralType) IsKnown() bool {
	switch r {
	case BetaCacheControlEphemeralTypeEphemeral:
		return true
	}
	return false
}

type BetaContentBlock struct {
	Type BetaContentBlockType `json:"type,required"`
	Text string               `json:"text"`
	ID   string               `json:"id"`
	Name string               `json:"name"`
	// This field can have the runtime type of [interface{}].
	Input interface{}          `json:"input,required"`
	JSON  betaContentBlockJSON `json:"-"`
	union BetaContentBlockUnion
}

// betaContentBlockJSON contains the JSON metadata for the struct
// [BetaContentBlock]
type betaContentBlockJSON struct {
	Type        apijson.Field
	Text        apijson.Field
	ID          apijson.Field
	Name        apijson.Field
	Input       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r betaContentBlockJSON) RawJSON() string {
	return r.raw
}

func (r *BetaContentBlock) UnmarshalJSON(data []byte) (err error) {
	*r = BetaContentBlock{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [BetaContentBlockUnion] interface which you can cast to the
// specific types for more type safety.
//
// Possible runtime types of the union are [BetaTextBlock], [BetaToolUseBlock].
func (r BetaContentBlock) AsUnion() BetaContentBlockUnion {
	return r.union
}

// Union satisfied by [BetaTextBlock] or [BetaToolUseBlock].
type BetaContentBlockUnion interface {
	implementsBetaContentBlock()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*BetaContentBlockUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaTextBlock{}),
			DiscriminatorValue: "text",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaToolUseBlock{}),
			DiscriminatorValue: "tool_use",
		},
	)
}

type BetaContentBlockType string

const (
	BetaContentBlockTypeText    BetaContentBlockType = "text"
	BetaContentBlockTypeToolUse BetaContentBlockType = "tool_use"
)

func (r BetaContentBlockType) IsKnown() bool {
	switch r {
	case BetaContentBlockTypeText, BetaContentBlockTypeToolUse:
		return true
	}
	return false
}

type BetaContentBlockParam struct {
	CacheControl param.Field[BetaCacheControlEphemeralParam] `json:"cache_control"`
	Type         param.Field[BetaContentBlockParamType]      `json:"type,required"`
	Text         param.Field[string]                         `json:"text"`
	Source       param.Field[interface{}]                    `json:"source,required"`
	ID           param.Field[string]                         `json:"id"`
	Name         param.Field[string]                         `json:"name"`
	Input        param.Field[interface{}]                    `json:"input,required"`
	ToolUseID    param.Field[string]                         `json:"tool_use_id"`
	IsError      param.Field[bool]                           `json:"is_error"`
	Content      param.Field[interface{}]                    `json:"content,required"`
}

func (r BetaContentBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaContentBlockParam) implementsBetaContentBlockParamUnion() {}

// Satisfied by [BetaTextBlockParam], [BetaImageBlockParam],
// [BetaToolUseBlockParam], [BetaToolResultBlockParam], [BetaContentBlockParam].
type BetaContentBlockParamUnion interface {
	implementsBetaContentBlockParamUnion()
}

type BetaContentBlockParamType string

const (
	BetaContentBlockParamTypeText       BetaContentBlockParamType = "text"
	BetaContentBlockParamTypeImage      BetaContentBlockParamType = "image"
	BetaContentBlockParamTypeToolUse    BetaContentBlockParamType = "tool_use"
	BetaContentBlockParamTypeToolResult BetaContentBlockParamType = "tool_result"
)

func (r BetaContentBlockParamType) IsKnown() bool {
	switch r {
	case BetaContentBlockParamTypeText, BetaContentBlockParamTypeImage, BetaContentBlockParamTypeToolUse, BetaContentBlockParamTypeToolResult:
		return true
	}
	return false
}

type BetaImageBlockParam struct {
	Source       param.Field[BetaImageBlockParamSource]      `json:"source,required"`
	Type         param.Field[BetaImageBlockParamType]        `json:"type,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam] `json:"cache_control"`
}

func (r BetaImageBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaImageBlockParam) implementsBetaContentBlockParamUnion() {}

func (r BetaImageBlockParam) implementsBetaToolResultBlockParamContentUnion() {}

type BetaImageBlockParamSource struct {
	Data      param.Field[string]                             `json:"data,required" format:"byte"`
	MediaType param.Field[BetaImageBlockParamSourceMediaType] `json:"media_type,required"`
	Type      param.Field[BetaImageBlockParamSourceType]      `json:"type,required"`
}

func (r BetaImageBlockParamSource) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaImageBlockParamSourceMediaType string

const (
	BetaImageBlockParamSourceMediaTypeImageJPEG BetaImageBlockParamSourceMediaType = "image/jpeg"
	BetaImageBlockParamSourceMediaTypeImagePNG  BetaImageBlockParamSourceMediaType = "image/png"
	BetaImageBlockParamSourceMediaTypeImageGIF  BetaImageBlockParamSourceMediaType = "image/gif"
	BetaImageBlockParamSourceMediaTypeImageWebP BetaImageBlockParamSourceMediaType = "image/webp"
)

func (r BetaImageBlockParamSourceMediaType) IsKnown() bool {
	switch r {
	case BetaImageBlockParamSourceMediaTypeImageJPEG, BetaImageBlockParamSourceMediaTypeImagePNG, BetaImageBlockParamSourceMediaTypeImageGIF, BetaImageBlockParamSourceMediaTypeImageWebP:
		return true
	}
	return false
}

type BetaImageBlockParamSourceType string

const (
	BetaImageBlockParamSourceTypeBase64 BetaImageBlockParamSourceType = "base64"
)

func (r BetaImageBlockParamSourceType) IsKnown() bool {
	switch r {
	case BetaImageBlockParamSourceTypeBase64:
		return true
	}
	return false
}

type BetaImageBlockParamType string

const (
	BetaImageBlockParamTypeImage BetaImageBlockParamType = "image"
)

func (r BetaImageBlockParamType) IsKnown() bool {
	switch r {
	case BetaImageBlockParamTypeImage:
		return true
	}
	return false
}

type BetaInputJSONDelta struct {
	PartialJSON string                 `json:"partial_json,required"`
	Type        BetaInputJSONDeltaType `json:"type,required"`
	JSON        betaInputJSONDeltaJSON `json:"-"`
}

// betaInputJSONDeltaJSON contains the JSON metadata for the struct
// [BetaInputJSONDelta]
type betaInputJSONDeltaJSON struct {
	PartialJSON apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaInputJSONDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaInputJSONDeltaJSON) RawJSON() string {
	return r.raw
}

func (r BetaInputJSONDelta) implementsBetaRawContentBlockDeltaEventDelta() {}

type BetaInputJSONDeltaType string

const (
	BetaInputJSONDeltaTypeInputJSONDelta BetaInputJSONDeltaType = "input_json_delta"
)

func (r BetaInputJSONDeltaType) IsKnown() bool {
	switch r {
	case BetaInputJSONDeltaTypeInputJSONDelta:
		return true
	}
	return false
}

type BetaMessage struct {
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
	Content []BetaContentBlock `json:"content,required"`
	// The model that will complete your prompt.\n\nSee
	// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	Model Model `json:"model,required"`
	// Conversational role of the generated message.
	//
	// This will always be `"assistant"`.
	Role BetaMessageRole `json:"role,required"`
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
	StopReason BetaMessageStopReason `json:"stop_reason,required,nullable"`
	// Which custom stop sequence was generated, if any.
	//
	// This value will be a non-null string if one of your custom stop sequences was
	// generated.
	StopSequence string `json:"stop_sequence,required,nullable"`
	// Object type.
	//
	// For Messages, this is always `"message"`.
	Type BetaMessageType `json:"type,required"`
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
	Usage BetaUsage       `json:"usage,required"`
	JSON  betaMessageJSON `json:"-"`
}

// betaMessageJSON contains the JSON metadata for the struct [BetaMessage]
type betaMessageJSON struct {
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

func (r *BetaMessage) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaMessageJSON) RawJSON() string {
	return r.raw
}

// Conversational role of the generated message.
//
// This will always be `"assistant"`.
type BetaMessageRole string

const (
	BetaMessageRoleAssistant BetaMessageRole = "assistant"
)

func (r BetaMessageRole) IsKnown() bool {
	switch r {
	case BetaMessageRoleAssistant:
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
type BetaMessageStopReason string

const (
	BetaMessageStopReasonEndTurn      BetaMessageStopReason = "end_turn"
	BetaMessageStopReasonMaxTokens    BetaMessageStopReason = "max_tokens"
	BetaMessageStopReasonStopSequence BetaMessageStopReason = "stop_sequence"
	BetaMessageStopReasonToolUse      BetaMessageStopReason = "tool_use"
)

func (r BetaMessageStopReason) IsKnown() bool {
	switch r {
	case BetaMessageStopReasonEndTurn, BetaMessageStopReasonMaxTokens, BetaMessageStopReasonStopSequence, BetaMessageStopReasonToolUse:
		return true
	}
	return false
}

// Object type.
//
// For Messages, this is always `"message"`.
type BetaMessageType string

const (
	BetaMessageTypeMessage BetaMessageType = "message"
)

func (r BetaMessageType) IsKnown() bool {
	switch r {
	case BetaMessageTypeMessage:
		return true
	}
	return false
}

type BetaMessageDeltaUsage struct {
	// The cumulative number of output tokens which were used.
	OutputTokens int64                     `json:"output_tokens,required"`
	JSON         betaMessageDeltaUsageJSON `json:"-"`
}

// betaMessageDeltaUsageJSON contains the JSON metadata for the struct
// [BetaMessageDeltaUsage]
type betaMessageDeltaUsageJSON struct {
	OutputTokens apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *BetaMessageDeltaUsage) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaMessageDeltaUsageJSON) RawJSON() string {
	return r.raw
}

type BetaMessageParam struct {
	Content param.Field[[]BetaContentBlockParamUnion] `json:"content,required"`
	Role    param.Field[BetaMessageParamRole]         `json:"role,required"`
}

func (r BetaMessageParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaMessageParamRole string

const (
	BetaMessageParamRoleUser      BetaMessageParamRole = "user"
	BetaMessageParamRoleAssistant BetaMessageParamRole = "assistant"
)

func (r BetaMessageParamRole) IsKnown() bool {
	switch r {
	case BetaMessageParamRoleUser, BetaMessageParamRoleAssistant:
		return true
	}
	return false
}

type BetaMetadataParam struct {
	// An external identifier for the user who is associated with the request.
	//
	// This should be a uuid, hash value, or other opaque identifier. Anthropic may use
	// this id to help detect abuse. Do not include any identifying information such as
	// name, email address, or phone number.
	UserID param.Field[string] `json:"user_id"`
}

func (r BetaMetadataParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaRawContentBlockDeltaEvent struct {
	Delta BetaRawContentBlockDeltaEventDelta `json:"delta,required"`
	Index int64                              `json:"index,required"`
	Type  BetaRawContentBlockDeltaEventType  `json:"type,required"`
	JSON  betaRawContentBlockDeltaEventJSON  `json:"-"`
}

// betaRawContentBlockDeltaEventJSON contains the JSON metadata for the struct
// [BetaRawContentBlockDeltaEvent]
type betaRawContentBlockDeltaEventJSON struct {
	Delta       apijson.Field
	Index       apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaRawContentBlockDeltaEvent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaRawContentBlockDeltaEventJSON) RawJSON() string {
	return r.raw
}

func (r BetaRawContentBlockDeltaEvent) implementsBetaRawMessageStreamEvent() {}

type BetaRawContentBlockDeltaEventDelta struct {
	Type        BetaRawContentBlockDeltaEventDeltaType `json:"type,required"`
	Text        string                                 `json:"text"`
	PartialJSON string                                 `json:"partial_json"`
	JSON        betaRawContentBlockDeltaEventDeltaJSON `json:"-"`
	union       BetaRawContentBlockDeltaEventDeltaUnion
}

// betaRawContentBlockDeltaEventDeltaJSON contains the JSON metadata for the struct
// [BetaRawContentBlockDeltaEventDelta]
type betaRawContentBlockDeltaEventDeltaJSON struct {
	Type        apijson.Field
	Text        apijson.Field
	PartialJSON apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r betaRawContentBlockDeltaEventDeltaJSON) RawJSON() string {
	return r.raw
}

func (r *BetaRawContentBlockDeltaEventDelta) UnmarshalJSON(data []byte) (err error) {
	*r = BetaRawContentBlockDeltaEventDelta{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [BetaRawContentBlockDeltaEventDeltaUnion] interface which you
// can cast to the specific types for more type safety.
//
// Possible runtime types of the union are [BetaTextDelta], [BetaInputJSONDelta].
func (r BetaRawContentBlockDeltaEventDelta) AsUnion() BetaRawContentBlockDeltaEventDeltaUnion {
	return r.union
}

// Union satisfied by [BetaTextDelta] or [BetaInputJSONDelta].
type BetaRawContentBlockDeltaEventDeltaUnion interface {
	implementsBetaRawContentBlockDeltaEventDelta()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*BetaRawContentBlockDeltaEventDeltaUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaTextDelta{}),
			DiscriminatorValue: "text_delta",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaInputJSONDelta{}),
			DiscriminatorValue: "input_json_delta",
		},
	)
}

type BetaRawContentBlockDeltaEventDeltaType string

const (
	BetaRawContentBlockDeltaEventDeltaTypeTextDelta      BetaRawContentBlockDeltaEventDeltaType = "text_delta"
	BetaRawContentBlockDeltaEventDeltaTypeInputJSONDelta BetaRawContentBlockDeltaEventDeltaType = "input_json_delta"
)

func (r BetaRawContentBlockDeltaEventDeltaType) IsKnown() bool {
	switch r {
	case BetaRawContentBlockDeltaEventDeltaTypeTextDelta, BetaRawContentBlockDeltaEventDeltaTypeInputJSONDelta:
		return true
	}
	return false
}

type BetaRawContentBlockDeltaEventType string

const (
	BetaRawContentBlockDeltaEventTypeContentBlockDelta BetaRawContentBlockDeltaEventType = "content_block_delta"
)

func (r BetaRawContentBlockDeltaEventType) IsKnown() bool {
	switch r {
	case BetaRawContentBlockDeltaEventTypeContentBlockDelta:
		return true
	}
	return false
}

type BetaRawContentBlockStartEvent struct {
	ContentBlock BetaRawContentBlockStartEventContentBlock `json:"content_block,required"`
	Index        int64                                     `json:"index,required"`
	Type         BetaRawContentBlockStartEventType         `json:"type,required"`
	JSON         betaRawContentBlockStartEventJSON         `json:"-"`
}

// betaRawContentBlockStartEventJSON contains the JSON metadata for the struct
// [BetaRawContentBlockStartEvent]
type betaRawContentBlockStartEventJSON struct {
	ContentBlock apijson.Field
	Index        apijson.Field
	Type         apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *BetaRawContentBlockStartEvent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaRawContentBlockStartEventJSON) RawJSON() string {
	return r.raw
}

func (r BetaRawContentBlockStartEvent) implementsBetaRawMessageStreamEvent() {}

type BetaRawContentBlockStartEventContentBlock struct {
	Type BetaRawContentBlockStartEventContentBlockType `json:"type,required"`
	Text string                                        `json:"text"`
	ID   string                                        `json:"id"`
	Name string                                        `json:"name"`
	// This field can have the runtime type of [interface{}].
	Input interface{}                                   `json:"input,required"`
	JSON  betaRawContentBlockStartEventContentBlockJSON `json:"-"`
	union BetaRawContentBlockStartEventContentBlockUnion
}

// betaRawContentBlockStartEventContentBlockJSON contains the JSON metadata for the
// struct [BetaRawContentBlockStartEventContentBlock]
type betaRawContentBlockStartEventContentBlockJSON struct {
	Type        apijson.Field
	Text        apijson.Field
	ID          apijson.Field
	Name        apijson.Field
	Input       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r betaRawContentBlockStartEventContentBlockJSON) RawJSON() string {
	return r.raw
}

func (r *BetaRawContentBlockStartEventContentBlock) UnmarshalJSON(data []byte) (err error) {
	*r = BetaRawContentBlockStartEventContentBlock{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [BetaRawContentBlockStartEventContentBlockUnion] interface
// which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are [BetaTextBlock], [BetaToolUseBlock].
func (r BetaRawContentBlockStartEventContentBlock) AsUnion() BetaRawContentBlockStartEventContentBlockUnion {
	return r.union
}

// Union satisfied by [BetaTextBlock] or [BetaToolUseBlock].
type BetaRawContentBlockStartEventContentBlockUnion interface {
	implementsBetaRawContentBlockStartEventContentBlock()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*BetaRawContentBlockStartEventContentBlockUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaTextBlock{}),
			DiscriminatorValue: "text",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaToolUseBlock{}),
			DiscriminatorValue: "tool_use",
		},
	)
}

type BetaRawContentBlockStartEventContentBlockType string

const (
	BetaRawContentBlockStartEventContentBlockTypeText    BetaRawContentBlockStartEventContentBlockType = "text"
	BetaRawContentBlockStartEventContentBlockTypeToolUse BetaRawContentBlockStartEventContentBlockType = "tool_use"
)

func (r BetaRawContentBlockStartEventContentBlockType) IsKnown() bool {
	switch r {
	case BetaRawContentBlockStartEventContentBlockTypeText, BetaRawContentBlockStartEventContentBlockTypeToolUse:
		return true
	}
	return false
}

type BetaRawContentBlockStartEventType string

const (
	BetaRawContentBlockStartEventTypeContentBlockStart BetaRawContentBlockStartEventType = "content_block_start"
)

func (r BetaRawContentBlockStartEventType) IsKnown() bool {
	switch r {
	case BetaRawContentBlockStartEventTypeContentBlockStart:
		return true
	}
	return false
}

type BetaRawContentBlockStopEvent struct {
	Index int64                            `json:"index,required"`
	Type  BetaRawContentBlockStopEventType `json:"type,required"`
	JSON  betaRawContentBlockStopEventJSON `json:"-"`
}

// betaRawContentBlockStopEventJSON contains the JSON metadata for the struct
// [BetaRawContentBlockStopEvent]
type betaRawContentBlockStopEventJSON struct {
	Index       apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaRawContentBlockStopEvent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaRawContentBlockStopEventJSON) RawJSON() string {
	return r.raw
}

func (r BetaRawContentBlockStopEvent) implementsBetaRawMessageStreamEvent() {}

type BetaRawContentBlockStopEventType string

const (
	BetaRawContentBlockStopEventTypeContentBlockStop BetaRawContentBlockStopEventType = "content_block_stop"
)

func (r BetaRawContentBlockStopEventType) IsKnown() bool {
	switch r {
	case BetaRawContentBlockStopEventTypeContentBlockStop:
		return true
	}
	return false
}

type BetaRawMessageDeltaEvent struct {
	Delta BetaRawMessageDeltaEventDelta `json:"delta,required"`
	Type  BetaRawMessageDeltaEventType  `json:"type,required"`
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
	Usage BetaMessageDeltaUsage        `json:"usage,required"`
	JSON  betaRawMessageDeltaEventJSON `json:"-"`
}

// betaRawMessageDeltaEventJSON contains the JSON metadata for the struct
// [BetaRawMessageDeltaEvent]
type betaRawMessageDeltaEventJSON struct {
	Delta       apijson.Field
	Type        apijson.Field
	Usage       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaRawMessageDeltaEvent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaRawMessageDeltaEventJSON) RawJSON() string {
	return r.raw
}

func (r BetaRawMessageDeltaEvent) implementsBetaRawMessageStreamEvent() {}

type BetaRawMessageDeltaEventDelta struct {
	StopReason   BetaRawMessageDeltaEventDeltaStopReason `json:"stop_reason,required,nullable"`
	StopSequence string                                  `json:"stop_sequence,required,nullable"`
	JSON         betaRawMessageDeltaEventDeltaJSON       `json:"-"`
}

// betaRawMessageDeltaEventDeltaJSON contains the JSON metadata for the struct
// [BetaRawMessageDeltaEventDelta]
type betaRawMessageDeltaEventDeltaJSON struct {
	StopReason   apijson.Field
	StopSequence apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *BetaRawMessageDeltaEventDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaRawMessageDeltaEventDeltaJSON) RawJSON() string {
	return r.raw
}

type BetaRawMessageDeltaEventDeltaStopReason string

const (
	BetaRawMessageDeltaEventDeltaStopReasonEndTurn      BetaRawMessageDeltaEventDeltaStopReason = "end_turn"
	BetaRawMessageDeltaEventDeltaStopReasonMaxTokens    BetaRawMessageDeltaEventDeltaStopReason = "max_tokens"
	BetaRawMessageDeltaEventDeltaStopReasonStopSequence BetaRawMessageDeltaEventDeltaStopReason = "stop_sequence"
	BetaRawMessageDeltaEventDeltaStopReasonToolUse      BetaRawMessageDeltaEventDeltaStopReason = "tool_use"
)

func (r BetaRawMessageDeltaEventDeltaStopReason) IsKnown() bool {
	switch r {
	case BetaRawMessageDeltaEventDeltaStopReasonEndTurn, BetaRawMessageDeltaEventDeltaStopReasonMaxTokens, BetaRawMessageDeltaEventDeltaStopReasonStopSequence, BetaRawMessageDeltaEventDeltaStopReasonToolUse:
		return true
	}
	return false
}

type BetaRawMessageDeltaEventType string

const (
	BetaRawMessageDeltaEventTypeMessageDelta BetaRawMessageDeltaEventType = "message_delta"
)

func (r BetaRawMessageDeltaEventType) IsKnown() bool {
	switch r {
	case BetaRawMessageDeltaEventTypeMessageDelta:
		return true
	}
	return false
}

type BetaRawMessageStartEvent struct {
	Message BetaMessage                  `json:"message,required"`
	Type    BetaRawMessageStartEventType `json:"type,required"`
	JSON    betaRawMessageStartEventJSON `json:"-"`
}

// betaRawMessageStartEventJSON contains the JSON metadata for the struct
// [BetaRawMessageStartEvent]
type betaRawMessageStartEventJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaRawMessageStartEvent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaRawMessageStartEventJSON) RawJSON() string {
	return r.raw
}

func (r BetaRawMessageStartEvent) implementsBetaRawMessageStreamEvent() {}

type BetaRawMessageStartEventType string

const (
	BetaRawMessageStartEventTypeMessageStart BetaRawMessageStartEventType = "message_start"
)

func (r BetaRawMessageStartEventType) IsKnown() bool {
	switch r {
	case BetaRawMessageStartEventTypeMessageStart:
		return true
	}
	return false
}

type BetaRawMessageStopEvent struct {
	Type BetaRawMessageStopEventType `json:"type,required"`
	JSON betaRawMessageStopEventJSON `json:"-"`
}

// betaRawMessageStopEventJSON contains the JSON metadata for the struct
// [BetaRawMessageStopEvent]
type betaRawMessageStopEventJSON struct {
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaRawMessageStopEvent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaRawMessageStopEventJSON) RawJSON() string {
	return r.raw
}

func (r BetaRawMessageStopEvent) implementsBetaRawMessageStreamEvent() {}

type BetaRawMessageStopEventType string

const (
	BetaRawMessageStopEventTypeMessageStop BetaRawMessageStopEventType = "message_stop"
)

func (r BetaRawMessageStopEventType) IsKnown() bool {
	switch r {
	case BetaRawMessageStopEventTypeMessageStop:
		return true
	}
	return false
}

type BetaRawMessageStreamEvent struct {
	Type    BetaRawMessageStreamEventType `json:"type,required"`
	Message BetaMessage                   `json:"message"`
	// This field can have the runtime type of [BetaRawMessageDeltaEventDelta],
	// [BetaRawContentBlockDeltaEventDelta].
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
	Usage BetaMessageDeltaUsage `json:"usage"`
	Index int64                 `json:"index"`
	// This field can have the runtime type of
	// [BetaRawContentBlockStartEventContentBlock].
	ContentBlock interface{}                   `json:"content_block,required"`
	JSON         betaRawMessageStreamEventJSON `json:"-"`
	union        BetaRawMessageStreamEventUnion
}

// betaRawMessageStreamEventJSON contains the JSON metadata for the struct
// [BetaRawMessageStreamEvent]
type betaRawMessageStreamEventJSON struct {
	Type         apijson.Field
	Message      apijson.Field
	Delta        apijson.Field
	Usage        apijson.Field
	Index        apijson.Field
	ContentBlock apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r betaRawMessageStreamEventJSON) RawJSON() string {
	return r.raw
}

func (r *BetaRawMessageStreamEvent) UnmarshalJSON(data []byte) (err error) {
	*r = BetaRawMessageStreamEvent{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [BetaRawMessageStreamEventUnion] interface which you can cast
// to the specific types for more type safety.
//
// Possible runtime types of the union are [BetaRawMessageStartEvent],
// [BetaRawMessageDeltaEvent], [BetaRawMessageStopEvent],
// [BetaRawContentBlockStartEvent], [BetaRawContentBlockDeltaEvent],
// [BetaRawContentBlockStopEvent].
func (r BetaRawMessageStreamEvent) AsUnion() BetaRawMessageStreamEventUnion {
	return r.union
}

// Union satisfied by [BetaRawMessageStartEvent], [BetaRawMessageDeltaEvent],
// [BetaRawMessageStopEvent], [BetaRawContentBlockStartEvent],
// [BetaRawContentBlockDeltaEvent] or [BetaRawContentBlockStopEvent].
type BetaRawMessageStreamEventUnion interface {
	implementsBetaRawMessageStreamEvent()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*BetaRawMessageStreamEventUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaRawMessageStartEvent{}),
			DiscriminatorValue: "message_start",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaRawMessageDeltaEvent{}),
			DiscriminatorValue: "message_delta",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaRawMessageStopEvent{}),
			DiscriminatorValue: "message_stop",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaRawContentBlockStartEvent{}),
			DiscriminatorValue: "content_block_start",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaRawContentBlockDeltaEvent{}),
			DiscriminatorValue: "content_block_delta",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaRawContentBlockStopEvent{}),
			DiscriminatorValue: "content_block_stop",
		},
	)
}

type BetaRawMessageStreamEventType string

const (
	BetaRawMessageStreamEventTypeMessageStart      BetaRawMessageStreamEventType = "message_start"
	BetaRawMessageStreamEventTypeMessageDelta      BetaRawMessageStreamEventType = "message_delta"
	BetaRawMessageStreamEventTypeMessageStop       BetaRawMessageStreamEventType = "message_stop"
	BetaRawMessageStreamEventTypeContentBlockStart BetaRawMessageStreamEventType = "content_block_start"
	BetaRawMessageStreamEventTypeContentBlockDelta BetaRawMessageStreamEventType = "content_block_delta"
	BetaRawMessageStreamEventTypeContentBlockStop  BetaRawMessageStreamEventType = "content_block_stop"
)

func (r BetaRawMessageStreamEventType) IsKnown() bool {
	switch r {
	case BetaRawMessageStreamEventTypeMessageStart, BetaRawMessageStreamEventTypeMessageDelta, BetaRawMessageStreamEventTypeMessageStop, BetaRawMessageStreamEventTypeContentBlockStart, BetaRawMessageStreamEventTypeContentBlockDelta, BetaRawMessageStreamEventTypeContentBlockStop:
		return true
	}
	return false
}

type BetaTextBlock struct {
	Text string            `json:"text,required"`
	Type BetaTextBlockType `json:"type,required"`
	JSON betaTextBlockJSON `json:"-"`
}

// betaTextBlockJSON contains the JSON metadata for the struct [BetaTextBlock]
type betaTextBlockJSON struct {
	Text        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaTextBlock) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaTextBlockJSON) RawJSON() string {
	return r.raw
}

func (r BetaTextBlock) implementsBetaContentBlock() {}

func (r BetaTextBlock) implementsBetaRawContentBlockStartEventContentBlock() {}

type BetaTextBlockType string

const (
	BetaTextBlockTypeText BetaTextBlockType = "text"
)

func (r BetaTextBlockType) IsKnown() bool {
	switch r {
	case BetaTextBlockTypeText:
		return true
	}
	return false
}

type BetaTextBlockParam struct {
	Text         param.Field[string]                         `json:"text,required"`
	Type         param.Field[BetaTextBlockParamType]         `json:"type,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam] `json:"cache_control"`
}

func (r BetaTextBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaTextBlockParam) implementsBetaContentBlockParamUnion() {}

func (r BetaTextBlockParam) implementsBetaToolResultBlockParamContentUnion() {}

type BetaTextBlockParamType string

const (
	BetaTextBlockParamTypeText BetaTextBlockParamType = "text"
)

func (r BetaTextBlockParamType) IsKnown() bool {
	switch r {
	case BetaTextBlockParamTypeText:
		return true
	}
	return false
}

type BetaTextDelta struct {
	Text string            `json:"text,required"`
	Type BetaTextDeltaType `json:"type,required"`
	JSON betaTextDeltaJSON `json:"-"`
}

// betaTextDeltaJSON contains the JSON metadata for the struct [BetaTextDelta]
type betaTextDeltaJSON struct {
	Text        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaTextDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaTextDeltaJSON) RawJSON() string {
	return r.raw
}

func (r BetaTextDelta) implementsBetaRawContentBlockDeltaEventDelta() {}

type BetaTextDeltaType string

const (
	BetaTextDeltaTypeTextDelta BetaTextDeltaType = "text_delta"
)

func (r BetaTextDeltaType) IsKnown() bool {
	switch r {
	case BetaTextDeltaTypeTextDelta:
		return true
	}
	return false
}

type BetaToolParam struct {
	// [JSON schema](https://json-schema.org/) for this tool's input.
	//
	// This defines the shape of the `input` that your tool accepts and that the model
	// will produce.
	InputSchema  param.Field[BetaToolInputSchemaParam]       `json:"input_schema,required"`
	Name         param.Field[string]                         `json:"name,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam] `json:"cache_control"`
	// Description of what this tool does.
	//
	// Tool descriptions should be as detailed as possible. The more information that
	// the model has about what the tool is and how to use it, the better it will
	// perform. You can use natural language descriptions to reinforce important
	// aspects of the tool input JSON schema.
	Description param.Field[string]       `json:"description"`
	Type        param.Field[BetaToolType] `json:"type"`
}

func (r BetaToolParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolParam) implementsBetaToolUnionUnionParam() {}

// [JSON schema](https://json-schema.org/) for this tool's input.
//
// This defines the shape of the `input` that your tool accepts and that the model
// will produce.
type BetaToolInputSchemaParam struct {
	Type        param.Field[BetaToolInputSchemaType] `json:"type,required"`
	Properties  param.Field[interface{}]             `json:"properties"`
	ExtraFields map[string]interface{}               `json:"-,extras"`
}

func (r BetaToolInputSchemaParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaToolInputSchemaType string

const (
	BetaToolInputSchemaTypeObject BetaToolInputSchemaType = "object"
)

func (r BetaToolInputSchemaType) IsKnown() bool {
	switch r {
	case BetaToolInputSchemaTypeObject:
		return true
	}
	return false
}

type BetaToolType string

const (
	BetaToolTypeCustom BetaToolType = "custom"
)

func (r BetaToolType) IsKnown() bool {
	switch r {
	case BetaToolTypeCustom:
		return true
	}
	return false
}

type BetaToolBash20241022Param struct {
	Name         param.Field[BetaToolBash20241022Name]       `json:"name,required"`
	Type         param.Field[BetaToolBash20241022Type]       `json:"type,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam] `json:"cache_control"`
}

func (r BetaToolBash20241022Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolBash20241022Param) implementsBetaToolUnionUnionParam() {}

type BetaToolBash20241022Name string

const (
	BetaToolBash20241022NameBash BetaToolBash20241022Name = "bash"
)

func (r BetaToolBash20241022Name) IsKnown() bool {
	switch r {
	case BetaToolBash20241022NameBash:
		return true
	}
	return false
}

type BetaToolBash20241022Type string

const (
	BetaToolBash20241022TypeBash20241022 BetaToolBash20241022Type = "bash_20241022"
)

func (r BetaToolBash20241022Type) IsKnown() bool {
	switch r {
	case BetaToolBash20241022TypeBash20241022:
		return true
	}
	return false
}

// How the model should use the provided tools. The model can use a specific tool,
// any available tool, or decide by itself.
type BetaToolChoiceParam struct {
	Type param.Field[BetaToolChoiceType] `json:"type,required"`
	// Whether to disable parallel tool use.
	//
	// Defaults to `false`. If set to `true`, the model will output at most one tool
	// use.
	DisableParallelToolUse param.Field[bool] `json:"disable_parallel_tool_use"`
	// The name of the tool to use.
	Name param.Field[string] `json:"name"`
}

func (r BetaToolChoiceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolChoiceParam) implementsBetaToolChoiceUnionParam() {}

// How the model should use the provided tools. The model can use a specific tool,
// any available tool, or decide by itself.
//
// Satisfied by [BetaToolChoiceAutoParam], [BetaToolChoiceAnyParam],
// [BetaToolChoiceToolParam], [BetaToolChoiceParam].
type BetaToolChoiceUnionParam interface {
	implementsBetaToolChoiceUnionParam()
}

type BetaToolChoiceType string

const (
	BetaToolChoiceTypeAuto BetaToolChoiceType = "auto"
	BetaToolChoiceTypeAny  BetaToolChoiceType = "any"
	BetaToolChoiceTypeTool BetaToolChoiceType = "tool"
)

func (r BetaToolChoiceType) IsKnown() bool {
	switch r {
	case BetaToolChoiceTypeAuto, BetaToolChoiceTypeAny, BetaToolChoiceTypeTool:
		return true
	}
	return false
}

// The model will use any available tools.
type BetaToolChoiceAnyParam struct {
	Type param.Field[BetaToolChoiceAnyType] `json:"type,required"`
	// Whether to disable parallel tool use.
	//
	// Defaults to `false`. If set to `true`, the model will output exactly one tool
	// use.
	DisableParallelToolUse param.Field[bool] `json:"disable_parallel_tool_use"`
}

func (r BetaToolChoiceAnyParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolChoiceAnyParam) implementsBetaToolChoiceUnionParam() {}

type BetaToolChoiceAnyType string

const (
	BetaToolChoiceAnyTypeAny BetaToolChoiceAnyType = "any"
)

func (r BetaToolChoiceAnyType) IsKnown() bool {
	switch r {
	case BetaToolChoiceAnyTypeAny:
		return true
	}
	return false
}

// The model will automatically decide whether to use tools.
type BetaToolChoiceAutoParam struct {
	Type param.Field[BetaToolChoiceAutoType] `json:"type,required"`
	// Whether to disable parallel tool use.
	//
	// Defaults to `false`. If set to `true`, the model will output at most one tool
	// use.
	DisableParallelToolUse param.Field[bool] `json:"disable_parallel_tool_use"`
}

func (r BetaToolChoiceAutoParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolChoiceAutoParam) implementsBetaToolChoiceUnionParam() {}

type BetaToolChoiceAutoType string

const (
	BetaToolChoiceAutoTypeAuto BetaToolChoiceAutoType = "auto"
)

func (r BetaToolChoiceAutoType) IsKnown() bool {
	switch r {
	case BetaToolChoiceAutoTypeAuto:
		return true
	}
	return false
}

// The model will use the specified tool with `tool_choice.name`.
type BetaToolChoiceToolParam struct {
	// The name of the tool to use.
	Name param.Field[string]                 `json:"name,required"`
	Type param.Field[BetaToolChoiceToolType] `json:"type,required"`
	// Whether to disable parallel tool use.
	//
	// Defaults to `false`. If set to `true`, the model will output exactly one tool
	// use.
	DisableParallelToolUse param.Field[bool] `json:"disable_parallel_tool_use"`
}

func (r BetaToolChoiceToolParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolChoiceToolParam) implementsBetaToolChoiceUnionParam() {}

type BetaToolChoiceToolType string

const (
	BetaToolChoiceToolTypeTool BetaToolChoiceToolType = "tool"
)

func (r BetaToolChoiceToolType) IsKnown() bool {
	switch r {
	case BetaToolChoiceToolTypeTool:
		return true
	}
	return false
}

type BetaToolComputerUse20241022Param struct {
	DisplayHeightPx param.Field[int64]                           `json:"display_height_px,required"`
	DisplayWidthPx  param.Field[int64]                           `json:"display_width_px,required"`
	Name            param.Field[BetaToolComputerUse20241022Name] `json:"name,required"`
	Type            param.Field[BetaToolComputerUse20241022Type] `json:"type,required"`
	CacheControl    param.Field[BetaCacheControlEphemeralParam]  `json:"cache_control"`
	DisplayNumber   param.Field[int64]                           `json:"display_number"`
}

func (r BetaToolComputerUse20241022Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolComputerUse20241022Param) implementsBetaToolUnionUnionParam() {}

type BetaToolComputerUse20241022Name string

const (
	BetaToolComputerUse20241022NameComputer BetaToolComputerUse20241022Name = "computer"
)

func (r BetaToolComputerUse20241022Name) IsKnown() bool {
	switch r {
	case BetaToolComputerUse20241022NameComputer:
		return true
	}
	return false
}

type BetaToolComputerUse20241022Type string

const (
	BetaToolComputerUse20241022TypeComputer20241022 BetaToolComputerUse20241022Type = "computer_20241022"
)

func (r BetaToolComputerUse20241022Type) IsKnown() bool {
	switch r {
	case BetaToolComputerUse20241022TypeComputer20241022:
		return true
	}
	return false
}

type BetaToolResultBlockParam struct {
	ToolUseID    param.Field[string]                                 `json:"tool_use_id,required"`
	Type         param.Field[BetaToolResultBlockParamType]           `json:"type,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam]         `json:"cache_control"`
	Content      param.Field[[]BetaToolResultBlockParamContentUnion] `json:"content"`
	IsError      param.Field[bool]                                   `json:"is_error"`
}

func (r BetaToolResultBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolResultBlockParam) implementsBetaContentBlockParamUnion() {}

type BetaToolResultBlockParamType string

const (
	BetaToolResultBlockParamTypeToolResult BetaToolResultBlockParamType = "tool_result"
)

func (r BetaToolResultBlockParamType) IsKnown() bool {
	switch r {
	case BetaToolResultBlockParamTypeToolResult:
		return true
	}
	return false
}

type BetaToolResultBlockParamContent struct {
	CacheControl param.Field[BetaCacheControlEphemeralParam]      `json:"cache_control"`
	Type         param.Field[BetaToolResultBlockParamContentType] `json:"type,required"`
	Text         param.Field[string]                              `json:"text"`
	Source       param.Field[interface{}]                         `json:"source,required"`
}

func (r BetaToolResultBlockParamContent) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolResultBlockParamContent) implementsBetaToolResultBlockParamContentUnion() {}

// Satisfied by [BetaTextBlockParam], [BetaImageBlockParam],
// [BetaToolResultBlockParamContent].
type BetaToolResultBlockParamContentUnion interface {
	implementsBetaToolResultBlockParamContentUnion()
}

type BetaToolResultBlockParamContentType string

const (
	BetaToolResultBlockParamContentTypeText  BetaToolResultBlockParamContentType = "text"
	BetaToolResultBlockParamContentTypeImage BetaToolResultBlockParamContentType = "image"
)

func (r BetaToolResultBlockParamContentType) IsKnown() bool {
	switch r {
	case BetaToolResultBlockParamContentTypeText, BetaToolResultBlockParamContentTypeImage:
		return true
	}
	return false
}

type BetaToolTextEditor20241022Param struct {
	Name         param.Field[BetaToolTextEditor20241022Name] `json:"name,required"`
	Type         param.Field[BetaToolTextEditor20241022Type] `json:"type,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam] `json:"cache_control"`
}

func (r BetaToolTextEditor20241022Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolTextEditor20241022Param) implementsBetaToolUnionUnionParam() {}

type BetaToolTextEditor20241022Name string

const (
	BetaToolTextEditor20241022NameStrReplaceEditor BetaToolTextEditor20241022Name = "str_replace_editor"
)

func (r BetaToolTextEditor20241022Name) IsKnown() bool {
	switch r {
	case BetaToolTextEditor20241022NameStrReplaceEditor:
		return true
	}
	return false
}

type BetaToolTextEditor20241022Type string

const (
	BetaToolTextEditor20241022TypeTextEditor20241022 BetaToolTextEditor20241022Type = "text_editor_20241022"
)

func (r BetaToolTextEditor20241022Type) IsKnown() bool {
	switch r {
	case BetaToolTextEditor20241022TypeTextEditor20241022:
		return true
	}
	return false
}

type BetaToolUnionParam struct {
	Type param.Field[BetaToolUnionType] `json:"type"`
	// Description of what this tool does.
	//
	// Tool descriptions should be as detailed as possible. The more information that
	// the model has about what the tool is and how to use it, the better it will
	// perform. You can use natural language descriptions to reinforce important
	// aspects of the tool input JSON schema.
	Description     param.Field[string]                         `json:"description"`
	Name            param.Field[string]                         `json:"name,required"`
	InputSchema     param.Field[interface{}]                    `json:"input_schema,required"`
	CacheControl    param.Field[BetaCacheControlEphemeralParam] `json:"cache_control"`
	DisplayHeightPx param.Field[int64]                          `json:"display_height_px"`
	DisplayWidthPx  param.Field[int64]                          `json:"display_width_px"`
	DisplayNumber   param.Field[int64]                          `json:"display_number"`
}

func (r BetaToolUnionParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolUnionParam) implementsBetaToolUnionUnionParam() {}

// Satisfied by [BetaToolParam], [BetaToolComputerUse20241022Param],
// [BetaToolBash20241022Param], [BetaToolTextEditor20241022Param],
// [BetaToolUnionParam].
type BetaToolUnionUnionParam interface {
	implementsBetaToolUnionUnionParam()
}

type BetaToolUnionType string

const (
	BetaToolUnionTypeCustom             BetaToolUnionType = "custom"
	BetaToolUnionTypeComputer20241022   BetaToolUnionType = "computer_20241022"
	BetaToolUnionTypeBash20241022       BetaToolUnionType = "bash_20241022"
	BetaToolUnionTypeTextEditor20241022 BetaToolUnionType = "text_editor_20241022"
)

func (r BetaToolUnionType) IsKnown() bool {
	switch r {
	case BetaToolUnionTypeCustom, BetaToolUnionTypeComputer20241022, BetaToolUnionTypeBash20241022, BetaToolUnionTypeTextEditor20241022:
		return true
	}
	return false
}

type BetaToolUseBlock struct {
	ID    string               `json:"id,required"`
	Input interface{}          `json:"input,required"`
	Name  string               `json:"name,required"`
	Type  BetaToolUseBlockType `json:"type,required"`
	JSON  betaToolUseBlockJSON `json:"-"`
}

// betaToolUseBlockJSON contains the JSON metadata for the struct
// [BetaToolUseBlock]
type betaToolUseBlockJSON struct {
	ID          apijson.Field
	Input       apijson.Field
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaToolUseBlock) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaToolUseBlockJSON) RawJSON() string {
	return r.raw
}

func (r BetaToolUseBlock) implementsBetaContentBlock() {}

func (r BetaToolUseBlock) implementsBetaRawContentBlockStartEventContentBlock() {}

type BetaToolUseBlockType string

const (
	BetaToolUseBlockTypeToolUse BetaToolUseBlockType = "tool_use"
)

func (r BetaToolUseBlockType) IsKnown() bool {
	switch r {
	case BetaToolUseBlockTypeToolUse:
		return true
	}
	return false
}

type BetaToolUseBlockParam struct {
	ID           param.Field[string]                         `json:"id,required"`
	Input        param.Field[interface{}]                    `json:"input,required"`
	Name         param.Field[string]                         `json:"name,required"`
	Type         param.Field[BetaToolUseBlockParamType]      `json:"type,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam] `json:"cache_control"`
}

func (r BetaToolUseBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolUseBlockParam) implementsBetaContentBlockParamUnion() {}

type BetaToolUseBlockParamType string

const (
	BetaToolUseBlockParamTypeToolUse BetaToolUseBlockParamType = "tool_use"
)

func (r BetaToolUseBlockParamType) IsKnown() bool {
	switch r {
	case BetaToolUseBlockParamTypeToolUse:
		return true
	}
	return false
}

type BetaUsage struct {
	// The number of input tokens used to create the cache entry.
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens,required,nullable"`
	// The number of input tokens read from the cache.
	CacheReadInputTokens int64 `json:"cache_read_input_tokens,required,nullable"`
	// The number of input tokens which were used.
	InputTokens int64 `json:"input_tokens,required"`
	// The number of output tokens which were used.
	OutputTokens int64         `json:"output_tokens,required"`
	JSON         betaUsageJSON `json:"-"`
}

// betaUsageJSON contains the JSON metadata for the struct [BetaUsage]
type betaUsageJSON struct {
	CacheCreationInputTokens apijson.Field
	CacheReadInputTokens     apijson.Field
	InputTokens              apijson.Field
	OutputTokens             apijson.Field
	raw                      string
	ExtraFields              map[string]apijson.Field
}

func (r *BetaUsage) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaUsageJSON) RawJSON() string {
	return r.raw
}

type BetaMessageNewParams struct {
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
	// the next `Message` in the conversation. Consecutive `user` or `assistant` turns
	// in your request will be combined into a single turn.
	//
	// Each input message must be an object with a `role` and `content`. You can
	// specify a single `user`-role message, or you can include multiple `user` and
	// `assistant` messages.
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
	Messages param.Field[[]BetaMessageParam] `json:"messages,required"`
	// The model that will complete your prompt.\n\nSee
	// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	Model param.Field[Model] `json:"model,required"`
	// An object describing metadata about the request.
	Metadata param.Field[BetaMetadataParam] `json:"metadata"`
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
	System param.Field[[]BetaTextBlockParam] `json:"system"`
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
	ToolChoice param.Field[BetaToolChoiceUnionParam] `json:"tool_choice"`
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
	Tools param.Field[[]BetaToolUnionUnionParam] `json:"tools"`
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
	// Optional header to specify the beta version(s) you want to use.
	Betas param.Field[[]AnthropicBeta] `header:"anthropic-beta"`
}

func (r BetaMessageNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
