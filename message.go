// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/resp"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
	"github.com/tidwall/gjson"
)

// MessageService contains methods and other services that help with interacting
// with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewMessageService] method instead.
type MessageService struct {
	Options []option.RequestOption
	Batches MessageBatchService
}

// NewMessageService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewMessageService(opts ...option.RequestOption) (r MessageService) {
	r = MessageService{}
	r.Options = opts
	r.Batches = NewMessageBatchService(opts...)
	return
}

// Send a structured list of input messages with text and/or image content, and the
// model will generate the next message in the conversation.
//
// The Messages API can be used for either single queries or stateless multi-turn
// conversations.
//
// Learn more about the Messages API in our [user guide](/en/docs/initial-setup)
//
// Note: If you choose to set a timeout for this request, we recommend 10 minutes.
func (r *MessageService) New(ctx context.Context, body MessageNewParams, opts ...option.RequestOption) (res *Message, err error) {
	opts = append(r.Options[:], opts...)
	path := "v1/messages"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Send a structured list of input messages with text and/or image content, and the
// model will generate the next message in the conversation.
//
// The Messages API can be used for either single queries or stateless multi-turn
// conversations.
//
// Learn more about the Messages API in our [user guide](/en/docs/initial-setup)
//
// Note: If you choose to set a timeout for this request, we recommend 10 minutes.
func (r *MessageService) NewStreaming(ctx context.Context, body MessageNewParams, opts ...option.RequestOption) (stream *ssestream.Stream[MessageStreamEventUnion]) {
	var (
		raw *http.Response
		err error
	)
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithJSONSet("stream", true)}, opts...)
	path := "v1/messages"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &raw, opts...)
	return ssestream.NewStream[MessageStreamEventUnion](ssestream.NewDecoder(raw), err)
}

// Count the number of tokens in a Message.
//
// The Token Count API can be used to count the number of tokens in a Message,
// including tools, images, and documents, without creating it.
//
// Learn more about token counting in our
// [user guide](/en/docs/build-with-claude/token-counting)
func (r *MessageService) CountTokens(ctx context.Context, body MessageCountTokensParams, opts ...option.RequestOption) (res *MessageTokensCount, err error) {
	opts = append(r.Options[:], opts...)
	path := "v1/messages/count_tokens"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// The properties Data, MediaType, Type are required.
type Base64ImageSourceParam struct {
	Data string `json:"data,required" format:"byte"`
	// Any of "image/jpeg", "image/png", "image/gif", "image/webp".
	MediaType Base64ImageSourceMediaType `json:"media_type,omitzero,required"`
	// This field can be elided, and will marshal its zero value as "base64".
	Type constant.Base64 `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f Base64ImageSourceParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r Base64ImageSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow Base64ImageSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type Base64ImageSourceMediaType string

const (
	Base64ImageSourceMediaTypeImageJPEG Base64ImageSourceMediaType = "image/jpeg"
	Base64ImageSourceMediaTypeImagePNG  Base64ImageSourceMediaType = "image/png"
	Base64ImageSourceMediaTypeImageGIF  Base64ImageSourceMediaType = "image/gif"
	Base64ImageSourceMediaTypeImageWebP Base64ImageSourceMediaType = "image/webp"
)

// The properties Data, MediaType, Type are required.
type Base64PDFSourceParam struct {
	Data string `json:"data,required" format:"byte"`
	// This field can be elided, and will marshal its zero value as "application/pdf".
	MediaType constant.ApplicationPDF `json:"media_type,required"`
	// This field can be elided, and will marshal its zero value as "base64".
	Type constant.Base64 `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f Base64PDFSourceParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r Base64PDFSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow Base64PDFSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The property Type is required.
type CacheControlEphemeralParam struct {
	// This field can be elided, and will marshal its zero value as "ephemeral".
	Type constant.Ephemeral `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CacheControlEphemeralParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r CacheControlEphemeralParam) MarshalJSON() (data []byte, err error) {
	type shadow CacheControlEphemeralParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type CitationCharLocation struct {
	CitedText      string                `json:"cited_text,required"`
	DocumentIndex  int64                 `json:"document_index,required"`
	DocumentTitle  string                `json:"document_title,required"`
	EndCharIndex   int64                 `json:"end_char_index,required"`
	StartCharIndex int64                 `json:"start_char_index,required"`
	Type           constant.CharLocation `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		CitedText      resp.Field
		DocumentIndex  resp.Field
		DocumentTitle  resp.Field
		EndCharIndex   resp.Field
		StartCharIndex resp.Field
		Type           resp.Field
		ExtraFields    map[string]resp.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CitationCharLocation) RawJSON() string { return r.JSON.raw }
func (r *CitationCharLocation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties CitedText, DocumentIndex, DocumentTitle, EndCharIndex,
// StartCharIndex, Type are required.
type CitationCharLocationParam struct {
	DocumentTitle  param.Opt[string] `json:"document_title,omitzero,required"`
	CitedText      string            `json:"cited_text,required"`
	DocumentIndex  int64             `json:"document_index,required"`
	EndCharIndex   int64             `json:"end_char_index,required"`
	StartCharIndex int64             `json:"start_char_index,required"`
	// This field can be elided, and will marshal its zero value as "char_location".
	Type constant.CharLocation `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CitationCharLocationParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r CitationCharLocationParam) MarshalJSON() (data []byte, err error) {
	type shadow CitationCharLocationParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type CitationContentBlockLocation struct {
	CitedText       string                        `json:"cited_text,required"`
	DocumentIndex   int64                         `json:"document_index,required"`
	DocumentTitle   string                        `json:"document_title,required"`
	EndBlockIndex   int64                         `json:"end_block_index,required"`
	StartBlockIndex int64                         `json:"start_block_index,required"`
	Type            constant.ContentBlockLocation `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		CitedText       resp.Field
		DocumentIndex   resp.Field
		DocumentTitle   resp.Field
		EndBlockIndex   resp.Field
		StartBlockIndex resp.Field
		Type            resp.Field
		ExtraFields     map[string]resp.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CitationContentBlockLocation) RawJSON() string { return r.JSON.raw }
func (r *CitationContentBlockLocation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties CitedText, DocumentIndex, DocumentTitle, EndBlockIndex,
// StartBlockIndex, Type are required.
type CitationContentBlockLocationParam struct {
	DocumentTitle   param.Opt[string] `json:"document_title,omitzero,required"`
	CitedText       string            `json:"cited_text,required"`
	DocumentIndex   int64             `json:"document_index,required"`
	EndBlockIndex   int64             `json:"end_block_index,required"`
	StartBlockIndex int64             `json:"start_block_index,required"`
	// This field can be elided, and will marshal its zero value as
	// "content_block_location".
	Type constant.ContentBlockLocation `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CitationContentBlockLocationParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r CitationContentBlockLocationParam) MarshalJSON() (data []byte, err error) {
	type shadow CitationContentBlockLocationParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type CitationPageLocation struct {
	CitedText       string                `json:"cited_text,required"`
	DocumentIndex   int64                 `json:"document_index,required"`
	DocumentTitle   string                `json:"document_title,required"`
	EndPageNumber   int64                 `json:"end_page_number,required"`
	StartPageNumber int64                 `json:"start_page_number,required"`
	Type            constant.PageLocation `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		CitedText       resp.Field
		DocumentIndex   resp.Field
		DocumentTitle   resp.Field
		EndPageNumber   resp.Field
		StartPageNumber resp.Field
		Type            resp.Field
		ExtraFields     map[string]resp.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CitationPageLocation) RawJSON() string { return r.JSON.raw }
func (r *CitationPageLocation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties CitedText, DocumentIndex, DocumentTitle, EndPageNumber,
// StartPageNumber, Type are required.
type CitationPageLocationParam struct {
	DocumentTitle   param.Opt[string] `json:"document_title,omitzero,required"`
	CitedText       string            `json:"cited_text,required"`
	DocumentIndex   int64             `json:"document_index,required"`
	EndPageNumber   int64             `json:"end_page_number,required"`
	StartPageNumber int64             `json:"start_page_number,required"`
	// This field can be elided, and will marshal its zero value as "page_location".
	Type constant.PageLocation `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CitationPageLocationParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r CitationPageLocationParam) MarshalJSON() (data []byte, err error) {
	type shadow CitationPageLocationParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type CitationsConfigParam struct {
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f CitationsConfigParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r CitationsConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow CitationsConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type CitationsDelta struct {
	Citation CitationsDeltaCitationUnion `json:"citation,required"`
	Type     constant.CitationsDelta     `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Citation    resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CitationsDelta) RawJSON() string { return r.JSON.raw }
func (r *CitationsDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// CitationsDeltaCitationUnion contains all possible properties and values from
// [CitationCharLocation], [CitationPageLocation], [CitationContentBlockLocation].
//
// Use the [CitationsDeltaCitationUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type CitationsDeltaCitationUnion struct {
	CitedText     string `json:"cited_text"`
	DocumentIndex int64  `json:"document_index"`
	DocumentTitle string `json:"document_title"`
	// This field is from variant [CitationCharLocation].
	EndCharIndex int64 `json:"end_char_index"`
	// This field is from variant [CitationCharLocation].
	StartCharIndex int64 `json:"start_char_index"`
	// Any of "char_location", "page_location", "content_block_location".
	Type string `json:"type"`
	// This field is from variant [CitationPageLocation].
	EndPageNumber int64 `json:"end_page_number"`
	// This field is from variant [CitationPageLocation].
	StartPageNumber int64 `json:"start_page_number"`
	// This field is from variant [CitationContentBlockLocation].
	EndBlockIndex int64 `json:"end_block_index"`
	// This field is from variant [CitationContentBlockLocation].
	StartBlockIndex int64 `json:"start_block_index"`
	JSON            struct {
		CitedText       resp.Field
		DocumentIndex   resp.Field
		DocumentTitle   resp.Field
		EndCharIndex    resp.Field
		StartCharIndex  resp.Field
		Type            resp.Field
		EndPageNumber   resp.Field
		StartPageNumber resp.Field
		EndBlockIndex   resp.Field
		StartBlockIndex resp.Field
		raw             string
	} `json:"-"`
}

// Use the following switch statement to find the correct variant
//
//	switch variant := CitationsDeltaCitationUnion.AsAny().(type) {
//	case CitationCharLocation:
//	case CitationPageLocation:
//	case CitationContentBlockLocation:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u CitationsDeltaCitationUnion) AsAny() any {
	switch u.Type {
	case "char_location":
		return u.AsResponseCharLocationCitation()
	case "page_location":
		return u.AsResponsePageLocationCitation()
	case "content_block_location":
		return u.AsResponseContentBlockLocationCitation()
	}
	return nil
}

func (u CitationsDeltaCitationUnion) AsResponseCharLocationCitation() (v CitationCharLocation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CitationsDeltaCitationUnion) AsResponsePageLocationCitation() (v CitationPageLocation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CitationsDeltaCitationUnion) AsResponseContentBlockLocationCitation() (v CitationContentBlockLocation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u CitationsDeltaCitationUnion) RawJSON() string { return u.JSON.raw }

func (r *CitationsDeltaCitationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ContentBlockUnion contains all possible properties and values from [TextBlock],
// [ToolUseBlock], [ThinkingBlock], [RedactedThinkingBlock].
//
// Use the [ContentBlockUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ContentBlockUnion struct {
	// This field is from variant [TextBlock].
	Citations []TextCitationUnion `json:"citations"`
	// This field is from variant [TextBlock].
	Text string `json:"text"`
	// Any of "text", "tool_use", "thinking", "redacted_thinking".
	Type string `json:"type"`
	// This field is from variant [ToolUseBlock].
	ID string `json:"id"`
	// This field is from variant [ToolUseBlock].
	Input json.RawMessage `json:"input"`
	// This field is from variant [ToolUseBlock].
	Name string `json:"name"`
	// This field is from variant [ThinkingBlock].
	Signature string `json:"signature"`
	// This field is from variant [ThinkingBlock].
	Thinking string `json:"thinking"`
	// This field is from variant [RedactedThinkingBlock].
	Data string `json:"data"`
	JSON struct {
		Citations resp.Field
		Text      resp.Field
		Type      resp.Field
		ID        resp.Field
		Input     resp.Field
		Name      resp.Field
		Signature resp.Field
		Thinking  resp.Field
		Data      resp.Field
		raw       string
	} `json:"-"`
}

func (r ContentBlockUnion) ToParam() ContentBlockParamUnion {
	switch variant := r.AsAny().(type) {
	case TextBlock:
		p := variant.ToParam()
		return ContentBlockParamUnion{OfRequestTextBlock: &p}
	case ToolUseBlock:
		p := variant.ToParam()
		return ContentBlockParamUnion{OfRequestToolUseBlock: &p}
	case ThinkingBlock:
		p := variant.ToParam()
		return ContentBlockParamUnion{OfRequestThinkingBlock: &p}
	case RedactedThinkingBlock:
		p := variant.ToParam()
		return ContentBlockParamUnion{OfRequestRedactedThinkingBlock: &p}
	}
	return ContentBlockParamUnion{}
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ContentBlockUnion.AsAny().(type) {
//	case TextBlock:
//	case ToolUseBlock:
//	case ThinkingBlock:
//	case RedactedThinkingBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ContentBlockUnion) AsAny() any {
	switch u.Type {
	case "text":
		return u.AsResponseTextBlock()
	case "tool_use":
		return u.AsResponseToolUseBlock()
	case "thinking":
		return u.AsResponseThinkingBlock()
	case "redacted_thinking":
		return u.AsResponseRedactedThinkingBlock()
	}
	return nil
}

func (u ContentBlockUnion) AsResponseTextBlock() (v TextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ContentBlockUnion) AsResponseToolUseBlock() (v ToolUseBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ContentBlockUnion) AsResponseThinkingBlock() (v ThinkingBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ContentBlockUnion) AsResponseRedactedThinkingBlock() (v RedactedThinkingBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ContentBlockUnion) RawJSON() string { return u.JSON.raw }

func (r *ContentBlockUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func ContentBlockParamOfRequestTextBlock(text string) ContentBlockParamUnion {
	var variant TextBlockParam
	variant.Text = text
	return ContentBlockParamUnion{OfRequestTextBlock: &variant}
}

func ContentBlockParamOfRequestImageBlock[T Base64ImageSourceParam | URLImageSourceParam](source T) ContentBlockParamUnion {
	var variant ImageBlockParam
	switch v := any(source).(type) {
	case Base64ImageSourceParam:
		variant.Source.OfBase64ImageSource = &v
	case URLImageSourceParam:
		variant.Source.OfURLImageSource = &v
	}
	return ContentBlockParamUnion{OfRequestImageBlock: &variant}
}

func ContentBlockParamOfRequestToolUseBlock(id string, input interface{}, name string) ContentBlockParamUnion {
	var variant ToolUseBlockParam
	variant.ID = id
	variant.Input = input
	variant.Name = name
	return ContentBlockParamUnion{OfRequestToolUseBlock: &variant}
}

func ContentBlockParamOfRequestToolResultBlock(toolUseID string) ContentBlockParamUnion {
	var variant ToolResultBlockParam
	variant.ToolUseID = toolUseID
	return ContentBlockParamUnion{OfRequestToolResultBlock: &variant}
}

func ContentBlockParamOfRequestDocumentBlock[
	T Base64PDFSourceParam | PlainTextSourceParam | ContentBlockSourceParam | URLPDFSourceParam,
](source T) ContentBlockParamUnion {
	var variant DocumentBlockParam
	switch v := any(source).(type) {
	case Base64PDFSourceParam:
		variant.Source.OfBase64PDFSource = &v
	case PlainTextSourceParam:
		variant.Source.OfPlainTextSource = &v
	case ContentBlockSourceParam:
		variant.Source.OfContentBlockSource = &v
	case URLPDFSourceParam:
		variant.Source.OfUrlpdfSource = &v
	}
	return ContentBlockParamUnion{OfRequestDocumentBlock: &variant}
}

func ContentBlockParamOfRequestThinkingBlock(signature string, thinking string) ContentBlockParamUnion {
	var variant ThinkingBlockParam
	variant.Signature = signature
	variant.Thinking = thinking
	return ContentBlockParamUnion{OfRequestThinkingBlock: &variant}
}

func ContentBlockParamOfRequestRedactedThinkingBlock(data string) ContentBlockParamUnion {
	var variant RedactedThinkingBlockParam
	variant.Data = data
	return ContentBlockParamUnion{OfRequestRedactedThinkingBlock: &variant}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ContentBlockParamUnion struct {
	OfRequestTextBlock             *TextBlockParam             `json:",omitzero,inline"`
	OfRequestImageBlock            *ImageBlockParam            `json:",omitzero,inline"`
	OfRequestToolUseBlock          *ToolUseBlockParam          `json:",omitzero,inline"`
	OfRequestToolResultBlock       *ToolResultBlockParam       `json:",omitzero,inline"`
	OfRequestDocumentBlock         *DocumentBlockParam         `json:",omitzero,inline"`
	OfRequestThinkingBlock         *ThinkingBlockParam         `json:",omitzero,inline"`
	OfRequestRedactedThinkingBlock *RedactedThinkingBlockParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u ContentBlockParamUnion) IsPresent() bool { return !param.IsOmitted(u) && !u.IsNull() }
func (u ContentBlockParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[ContentBlockParamUnion](u.OfRequestTextBlock,
		u.OfRequestImageBlock,
		u.OfRequestToolUseBlock,
		u.OfRequestToolResultBlock,
		u.OfRequestDocumentBlock,
		u.OfRequestThinkingBlock,
		u.OfRequestRedactedThinkingBlock)
}

func (u *ContentBlockParamUnion) asAny() any {
	if !param.IsOmitted(u.OfRequestTextBlock) {
		return u.OfRequestTextBlock
	} else if !param.IsOmitted(u.OfRequestImageBlock) {
		return u.OfRequestImageBlock
	} else if !param.IsOmitted(u.OfRequestToolUseBlock) {
		return u.OfRequestToolUseBlock
	} else if !param.IsOmitted(u.OfRequestToolResultBlock) {
		return u.OfRequestToolResultBlock
	} else if !param.IsOmitted(u.OfRequestDocumentBlock) {
		return u.OfRequestDocumentBlock
	} else if !param.IsOmitted(u.OfRequestThinkingBlock) {
		return u.OfRequestThinkingBlock
	} else if !param.IsOmitted(u.OfRequestRedactedThinkingBlock) {
		return u.OfRequestRedactedThinkingBlock
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetText() *string {
	if vt := u.OfRequestTextBlock; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetID() *string {
	if vt := u.OfRequestToolUseBlock; vt != nil {
		return &vt.ID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetInput() *interface{} {
	if vt := u.OfRequestToolUseBlock; vt != nil {
		return &vt.Input
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetName() *string {
	if vt := u.OfRequestToolUseBlock; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetToolUseID() *string {
	if vt := u.OfRequestToolResultBlock; vt != nil {
		return &vt.ToolUseID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetContent() *[]ToolResultBlockParamContentUnion {
	if vt := u.OfRequestToolResultBlock; vt != nil {
		return &vt.Content
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetIsError() *bool {
	if vt := u.OfRequestToolResultBlock; vt != nil && vt.IsError.IsPresent() {
		return &vt.IsError.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetContext() *string {
	if vt := u.OfRequestDocumentBlock; vt != nil && vt.Context.IsPresent() {
		return &vt.Context.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetTitle() *string {
	if vt := u.OfRequestDocumentBlock; vt != nil && vt.Title.IsPresent() {
		return &vt.Title.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetSignature() *string {
	if vt := u.OfRequestThinkingBlock; vt != nil {
		return &vt.Signature
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetThinking() *string {
	if vt := u.OfRequestThinkingBlock; vt != nil {
		return &vt.Thinking
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetData() *string {
	if vt := u.OfRequestRedactedThinkingBlock; vt != nil {
		return &vt.Data
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ContentBlockParamUnion) GetType() *string {
	if vt := u.OfRequestTextBlock; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestImageBlock; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestToolUseBlock; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestToolResultBlock; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestDocumentBlock; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestThinkingBlock; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestRedactedThinkingBlock; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's CacheControl property, if present.
func (u ContentBlockParamUnion) GetCacheControl() *CacheControlEphemeralParam {
	if vt := u.OfRequestTextBlock; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfRequestImageBlock; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfRequestToolUseBlock; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfRequestToolResultBlock; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfRequestDocumentBlock; vt != nil {
		return &vt.CacheControl
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ContentBlockParamUnion) GetCitations() (res contentBlockParamUnionCitations) {
	if vt := u.OfRequestTextBlock; vt != nil {
		res.ofTextBlockCitations = &vt.Citations
	} else if vt := u.OfRequestDocumentBlock; vt != nil {
		res.ofCitationsConfig = &vt.Citations
	}
	return
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type contentBlockParamUnionCitations struct {
	ofTextBlockCitations *[]TextCitationParamUnion
	ofCitationsConfig    *CitationsConfigParam
}

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]anthropic.TextCitationParamUnion:
//	case *anthropic.CitationsConfigParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u contentBlockParamUnionCitations) AsAny() any {
	if !param.IsOmitted(u.ofTextBlockCitations) {
		return u.ofTextBlockCitations
	} else if !param.IsOmitted(u.ofCitationsConfig) {
		return u.ofCitationsConfig
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionCitations) GetEnabled() *bool {
	if vt := u.ofCitationsConfig; vt != nil && vt.Enabled.IsPresent() {
		return &vt.Enabled.Value
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ContentBlockParamUnion) GetSource() (res contentBlockParamUnionSource) {
	if vt := u.OfRequestImageBlock; vt != nil {
		res.ofImageBlockSource = &vt.Source
	} else if vt := u.OfRequestDocumentBlock; vt != nil {
		res.ofDocumentBlockSource = &vt.Source
	}
	return
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type contentBlockParamUnionSource struct {
	ofImageBlockSource    *ImageBlockParamSourceUnion
	ofDocumentBlockSource *DocumentBlockParamSourceUnion
}

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *anthropic.Base64ImageSourceParam:
//	case *anthropic.URLImageSourceParam:
//	case *anthropic.Base64PDFSourceParam:
//	case *anthropic.PlainTextSourceParam:
//	case *anthropic.ContentBlockSourceParam:
//	case *anthropic.URLPDFSourceParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u contentBlockParamUnionSource) AsAny() any {
	if !param.IsOmitted(u.ofImageBlockSource) {
		return u.ofImageBlockSource.asAny()
	} else if !param.IsOmitted(u.ofDocumentBlockSource) {
		return u.ofDocumentBlockSource.asAny()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionSource) GetContent() *ContentBlockSourceContentUnionParam {
	if u.ofDocumentBlockSource != nil {
		return u.ofDocumentBlockSource.GetContent()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionSource) GetData() *string {
	if u.ofImageBlockSource != nil {
		return u.ofImageBlockSource.GetData()
	} else if u.ofDocumentBlockSource != nil {
		return u.ofDocumentBlockSource.GetData()
	} else if u.ofDocumentBlockSource != nil {
		return u.ofDocumentBlockSource.GetData()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionSource) GetMediaType() *string {
	if u.ofImageBlockSource != nil {
		return u.ofImageBlockSource.GetMediaType()
	} else if u.ofDocumentBlockSource != nil {
		return u.ofDocumentBlockSource.GetMediaType()
	} else if u.ofDocumentBlockSource != nil {
		return u.ofDocumentBlockSource.GetMediaType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionSource) GetType() *string {
	if u.ofImageBlockSource != nil {
		return u.ofImageBlockSource.GetType()
	} else if u.ofImageBlockSource != nil {
		return u.ofImageBlockSource.GetType()
	} else if u.ofDocumentBlockSource != nil {
		return u.ofDocumentBlockSource.GetType()
	} else if u.ofDocumentBlockSource != nil {
		return u.ofDocumentBlockSource.GetType()
	} else if u.ofDocumentBlockSource != nil {
		return u.ofDocumentBlockSource.GetType()
	} else if u.ofDocumentBlockSource != nil {
		return u.ofDocumentBlockSource.GetType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u contentBlockParamUnionSource) GetURL() *string {
	if u.ofImageBlockSource != nil {
		return u.ofImageBlockSource.GetURL()
	} else if u.ofDocumentBlockSource != nil {
		return u.ofDocumentBlockSource.GetURL()
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ContentBlockParamUnion](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(TextBlockParam{}),
			DiscriminatorValue: "text",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ImageBlockParam{}),
			DiscriminatorValue: "image",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ToolUseBlockParam{}),
			DiscriminatorValue: "tool_use",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ToolResultBlockParam{}),
			DiscriminatorValue: "tool_result",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(DocumentBlockParam{}),
			DiscriminatorValue: "document",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ThinkingBlockParam{}),
			DiscriminatorValue: "thinking",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(RedactedThinkingBlockParam{}),
			DiscriminatorValue: "redacted_thinking",
		},
	)
}

// The properties Content, Type are required.
type ContentBlockSourceParam struct {
	Content ContentBlockSourceContentUnionParam `json:"content,omitzero,required"`
	// This field can be elided, and will marshal its zero value as "content".
	Type constant.Content `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f ContentBlockSourceParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r ContentBlockSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow ContentBlockSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ContentBlockSourceContentUnionParam struct {
	OfString                    param.Opt[string]                     `json:",omitzero,inline"`
	OfContentBlockSourceContent []ContentBlockSourceContentUnionParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u ContentBlockSourceContentUnionParam) IsPresent() bool {
	return !param.IsOmitted(u) && !u.IsNull()
}
func (u ContentBlockSourceContentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[ContentBlockSourceContentUnionParam](u.OfString, u.OfContentBlockSourceContent)
}

func (u *ContentBlockSourceContentUnionParam) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfContentBlockSourceContent) {
		return &u.OfContentBlockSourceContent
	}
	return nil
}

// The properties Source, Type are required.
type DocumentBlockParam struct {
	Source       DocumentBlockParamSourceUnion `json:"source,omitzero,required"`
	Context      param.Opt[string]             `json:"context,omitzero"`
	Title        param.Opt[string]             `json:"title,omitzero"`
	CacheControl CacheControlEphemeralParam    `json:"cache_control,omitzero"`
	Citations    CitationsConfigParam          `json:"citations,omitzero"`
	// This field can be elided, and will marshal its zero value as "document".
	Type constant.Document `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f DocumentBlockParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r DocumentBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow DocumentBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DocumentBlockParamSourceUnion struct {
	OfBase64PDFSource    *Base64PDFSourceParam    `json:",omitzero,inline"`
	OfPlainTextSource    *PlainTextSourceParam    `json:",omitzero,inline"`
	OfContentBlockSource *ContentBlockSourceParam `json:",omitzero,inline"`
	OfUrlpdfSource       *URLPDFSourceParam       `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u DocumentBlockParamSourceUnion) IsPresent() bool { return !param.IsOmitted(u) && !u.IsNull() }
func (u DocumentBlockParamSourceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[DocumentBlockParamSourceUnion](u.OfBase64PDFSource, u.OfPlainTextSource, u.OfContentBlockSource, u.OfUrlpdfSource)
}

func (u *DocumentBlockParamSourceUnion) asAny() any {
	if !param.IsOmitted(u.OfBase64PDFSource) {
		return u.OfBase64PDFSource
	} else if !param.IsOmitted(u.OfPlainTextSource) {
		return u.OfPlainTextSource
	} else if !param.IsOmitted(u.OfContentBlockSource) {
		return u.OfContentBlockSource
	} else if !param.IsOmitted(u.OfUrlpdfSource) {
		return u.OfUrlpdfSource
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DocumentBlockParamSourceUnion) GetContent() *ContentBlockSourceContentUnionParam {
	if vt := u.OfContentBlockSource; vt != nil {
		return &vt.Content
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DocumentBlockParamSourceUnion) GetURL() *string {
	if vt := u.OfUrlpdfSource; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DocumentBlockParamSourceUnion) GetData() *string {
	if vt := u.OfBase64PDFSource; vt != nil {
		return (*string)(&vt.Data)
	} else if vt := u.OfPlainTextSource; vt != nil {
		return (*string)(&vt.Data)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DocumentBlockParamSourceUnion) GetMediaType() *string {
	if vt := u.OfBase64PDFSource; vt != nil {
		return (*string)(&vt.MediaType)
	} else if vt := u.OfPlainTextSource; vt != nil {
		return (*string)(&vt.MediaType)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DocumentBlockParamSourceUnion) GetType() *string {
	if vt := u.OfBase64PDFSource; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfPlainTextSource; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfContentBlockSource; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfUrlpdfSource; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DocumentBlockParamSourceUnion](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(Base64PDFSourceParam{}),
			DiscriminatorValue: "base64",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(PlainTextSourceParam{}),
			DiscriminatorValue: "text",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ContentBlockSourceParam{}),
			DiscriminatorValue: "content",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(URLPDFSourceParam{}),
			DiscriminatorValue: "url",
		},
	)
}

// The properties Source, Type are required.
type ImageBlockParam struct {
	Source       ImageBlockParamSourceUnion `json:"source,omitzero,required"`
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as "image".
	Type constant.Image `json:"type,required"`
	paramObj
}

func NewImageBlockBase64(mediaType string, encodedData string) ContentBlockParamUnion {
	return ContentBlockParamUnion{
		OfRequestImageBlock: &ImageBlockParam{
			Source: ImageBlockParamSourceUnion{
				OfBase64ImageSource: &Base64ImageSourceParam{
					Data:      encodedData,
					MediaType: Base64ImageSourceMediaType(mediaType),
				},
			},
		},
	}
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f ImageBlockParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r ImageBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ImageBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ImageBlockParamSourceUnion struct {
	OfBase64ImageSource *Base64ImageSourceParam `json:",omitzero,inline"`
	OfURLImageSource    *URLImageSourceParam    `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u ImageBlockParamSourceUnion) IsPresent() bool { return !param.IsOmitted(u) && !u.IsNull() }
func (u ImageBlockParamSourceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[ImageBlockParamSourceUnion](u.OfBase64ImageSource, u.OfURLImageSource)
}

func (u *ImageBlockParamSourceUnion) asAny() any {
	if !param.IsOmitted(u.OfBase64ImageSource) {
		return u.OfBase64ImageSource
	} else if !param.IsOmitted(u.OfURLImageSource) {
		return u.OfURLImageSource
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ImageBlockParamSourceUnion) GetData() *string {
	if vt := u.OfBase64ImageSource; vt != nil {
		return &vt.Data
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ImageBlockParamSourceUnion) GetMediaType() *string {
	if vt := u.OfBase64ImageSource; vt != nil {
		return (*string)(&vt.MediaType)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ImageBlockParamSourceUnion) GetURL() *string {
	if vt := u.OfURLImageSource; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ImageBlockParamSourceUnion) GetType() *string {
	if vt := u.OfBase64ImageSource; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfURLImageSource; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ImageBlockParamSourceUnion](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(Base64ImageSourceParam{}),
			DiscriminatorValue: "base64",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(URLImageSourceParam{}),
			DiscriminatorValue: "url",
		},
	)
}

type InputJSONDelta struct {
	PartialJSON string                  `json:"partial_json,required"`
	Type        constant.InputJSONDelta `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		PartialJSON resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r InputJSONDelta) RawJSON() string { return r.JSON.raw }
func (r *InputJSONDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Message struct {
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
	Content []ContentBlockUnion `json:"content,required"`
	// The model that will complete your prompt.\n\nSee
	// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	Model Model `json:"model,required"`
	// Conversational role of the generated message.
	//
	// This will always be `"assistant"`.
	Role constant.Assistant `json:"role,required"`
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
	//
	// Any of "end_turn", "max_tokens", "stop_sequence", "tool_use".
	StopReason MessageStopReason `json:"stop_reason,required"`
	// Which custom stop sequence was generated, if any.
	//
	// This value will be a non-null string if one of your custom stop sequences was
	// generated.
	StopSequence string `json:"stop_sequence,required"`
	// Object type.
	//
	// For Messages, this is always `"message"`.
	Type constant.Message `json:"type,required"`
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
	//
	// Total input tokens in a request is the summation of `input_tokens`,
	// `cache_creation_input_tokens`, and `cache_read_input_tokens`.
	Usage Usage `json:"usage,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID           resp.Field
		Content      resp.Field
		Model        resp.Field
		Role         resp.Field
		StopReason   resp.Field
		StopSequence resp.Field
		Type         resp.Field
		Usage        resp.Field
		ExtraFields  map[string]resp.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Message) RawJSON() string { return r.JSON.raw }
func (r *Message) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (r Message) ToParam() MessageParam {
	var p MessageParam
	p.Role = MessageParamRole(r.Role)
	p.Content = make([]ContentBlockParamUnion, len(r.Content))
	for i, c := range r.Content {
		p.Content[i] = c.ToParam()
	}
	return p
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
type MessageStopReason string

const (
	MessageStopReasonEndTurn      MessageStopReason = "end_turn"
	MessageStopReasonMaxTokens    MessageStopReason = "max_tokens"
	MessageStopReasonStopSequence MessageStopReason = "stop_sequence"
	MessageStopReasonToolUse      MessageStopReason = "tool_use"
)

func MessageCountTokensToolParamOfTool(inputSchema ToolInputSchemaParam, name string) MessageCountTokensToolUnionParam {
	var variant ToolParam
	variant.InputSchema = inputSchema
	variant.Name = name
	return MessageCountTokensToolUnionParam{OfTool: &variant}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type MessageCountTokensToolUnionParam struct {
	OfTool               *ToolParam                   `json:",omitzero,inline"`
	OfBashTool20250124   *ToolBash20250124Param       `json:",omitzero,inline"`
	OfTextEditor20250124 *ToolTextEditor20250124Param `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u MessageCountTokensToolUnionParam) IsPresent() bool { return !param.IsOmitted(u) && !u.IsNull() }
func (u MessageCountTokensToolUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[MessageCountTokensToolUnionParam](u.OfTool, u.OfBashTool20250124, u.OfTextEditor20250124)
}

func (u *MessageCountTokensToolUnionParam) asAny() any {
	if !param.IsOmitted(u.OfTool) {
		return u.OfTool
	} else if !param.IsOmitted(u.OfBashTool20250124) {
		return u.OfBashTool20250124
	} else if !param.IsOmitted(u.OfTextEditor20250124) {
		return u.OfTextEditor20250124
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u MessageCountTokensToolUnionParam) GetInputSchema() *ToolInputSchemaParam {
	if vt := u.OfTool; vt != nil {
		return &vt.InputSchema
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u MessageCountTokensToolUnionParam) GetDescription() *string {
	if vt := u.OfTool; vt != nil && vt.Description.IsPresent() {
		return &vt.Description.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u MessageCountTokensToolUnionParam) GetName() *string {
	if vt := u.OfTool; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfBashTool20250124; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfTextEditor20250124; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u MessageCountTokensToolUnionParam) GetType() *string {
	if vt := u.OfBashTool20250124; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfTextEditor20250124; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's CacheControl property, if present.
func (u MessageCountTokensToolUnionParam) GetCacheControl() *CacheControlEphemeralParam {
	if vt := u.OfTool; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfBashTool20250124; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfTextEditor20250124; vt != nil {
		return &vt.CacheControl
	}
	return nil
}

type MessageDeltaUsage struct {
	// The cumulative number of output tokens which were used.
	OutputTokens int64 `json:"output_tokens,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		OutputTokens resp.Field
		ExtraFields  map[string]resp.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r MessageDeltaUsage) RawJSON() string { return r.JSON.raw }
func (r *MessageDeltaUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Content, Role are required.
type MessageParam struct {
	Content []ContentBlockParamUnion `json:"content,omitzero,required"`
	// Any of "user", "assistant".
	Role MessageParamRole `json:"role,omitzero,required"`
	paramObj
}

func NewUserMessage(blocks ...ContentBlockParamUnion) MessageParam {
	return MessageParam{
		Role:    MessageParamRoleUser,
		Content: blocks,
	}
}

func NewAssistantMessage(blocks ...ContentBlockParamUnion) MessageParam {
	return MessageParam{
		Role:    MessageParamRoleAssistant,
		Content: blocks,
	}
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f MessageParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r MessageParam) MarshalJSON() (data []byte, err error) {
	type shadow MessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type MessageParamRole string

const (
	MessageParamRoleUser      MessageParamRole = "user"
	MessageParamRoleAssistant MessageParamRole = "assistant"
)

type MessageTokensCount struct {
	// The total number of tokens across the provided list of messages, system prompt,
	// and tools.
	InputTokens int64 `json:"input_tokens,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		InputTokens resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r MessageTokensCount) RawJSON() string { return r.JSON.raw }
func (r *MessageTokensCount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type MetadataParam struct {
	// An external identifier for the user who is associated with the request.
	//
	// This should be a uuid, hash value, or other opaque identifier. Anthropic may use
	// this id to help detect abuse. Do not include any identifying information such as
	// name, email address, or phone number.
	UserID param.Opt[string] `json:"user_id,omitzero"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f MetadataParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r MetadataParam) MarshalJSON() (data []byte, err error) {
	type shadow MetadataParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The model that will complete your prompt.\n\nSee
// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
// details and options.
type Model = string

const (
	ModelClaude3_7SonnetLatest      Model = "claude-3-7-sonnet-latest"
	ModelClaude3_7Sonnet20250219    Model = "claude-3-7-sonnet-20250219"
	ModelClaude3_5HaikuLatest       Model = "claude-3-5-haiku-latest"
	ModelClaude3_5Haiku20241022     Model = "claude-3-5-haiku-20241022"
	ModelClaude3_5SonnetLatest      Model = "claude-3-5-sonnet-latest"
	ModelClaude3_5Sonnet20241022    Model = "claude-3-5-sonnet-20241022"
	ModelClaude_3_5_Sonnet_20240620 Model = "claude-3-5-sonnet-20240620"
	ModelClaude3OpusLatest          Model = "claude-3-opus-latest"
	ModelClaude_3_Opus_20240229     Model = "claude-3-opus-20240229"
	// Deprecated: Will reach end-of-life on July 21st, 2025. Please migrate to a newer
	// model. Visit https://docs.anthropic.com/en/docs/resources/model-deprecations for
	// more information.
	ModelClaude_3_Sonnet_20240229 Model = "claude-3-sonnet-20240229"
	ModelClaude_3_Haiku_20240307  Model = "claude-3-haiku-20240307"
	// Deprecated: Will reach end-of-life on July 21st, 2025. Please migrate to a newer
	// model. Visit https://docs.anthropic.com/en/docs/resources/model-deprecations for
	// more information.
	ModelClaude_2_1 Model = "claude-2.1"
	// Deprecated: Will reach end-of-life on July 21st, 2025. Please migrate to a newer
	// model. Visit https://docs.anthropic.com/en/docs/resources/model-deprecations for
	// more information.
	ModelClaude_2_0 Model = "claude-2.0"
)

// The properties Data, MediaType, Type are required.
type PlainTextSourceParam struct {
	Data string `json:"data,required"`
	// This field can be elided, and will marshal its zero value as "text/plain".
	MediaType constant.TextPlain `json:"media_type,required"`
	// This field can be elided, and will marshal its zero value as "text".
	Type constant.Text `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f PlainTextSourceParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r PlainTextSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow PlainTextSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type ContentBlockDeltaEvent struct {
	Delta ContentBlockDeltaEventDeltaUnion `json:"delta,required"`
	Index int64                            `json:"index,required"`
	Type  constant.ContentBlockDelta       `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Delta       resp.Field
		Index       resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ContentBlockDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *ContentBlockDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ContentBlockDeltaEventDeltaUnion contains all possible properties and values
// from [TextDelta], [InputJSONDelta], [CitationsDelta], [ThinkingDelta],
// [SignatureDelta].
//
// Use the [ContentBlockDeltaEventDeltaUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ContentBlockDeltaEventDeltaUnion struct {
	// This field is from variant [TextDelta].
	Text string `json:"text"`
	// Any of "text_delta", "input_json_delta", "citations_delta", "thinking_delta",
	// "signature_delta".
	Type string `json:"type"`
	// This field is from variant [InputJSONDelta].
	PartialJSON string `json:"partial_json"`
	// This field is from variant [CitationsDelta].
	Citation CitationsDeltaCitationUnion `json:"citation"`
	// This field is from variant [ThinkingDelta].
	Thinking string `json:"thinking"`
	// This field is from variant [SignatureDelta].
	Signature string `json:"signature"`
	JSON      struct {
		Text        resp.Field
		Type        resp.Field
		PartialJSON resp.Field
		Citation    resp.Field
		Thinking    resp.Field
		Signature   resp.Field
		raw         string
	} `json:"-"`
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ContentBlockDeltaEventDeltaUnion.AsAny().(type) {
//	case TextDelta:
//	case InputJSONDelta:
//	case CitationsDelta:
//	case ThinkingDelta:
//	case SignatureDelta:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ContentBlockDeltaEventDeltaUnion) AsAny() any {
	switch u.Type {
	case "text_delta":
		return u.AsTextContentBlockDelta()
	case "input_json_delta":
		return u.AsInputJSONContentBlockDelta()
	case "citations_delta":
		return u.AsCitationsDelta()
	case "thinking_delta":
		return u.AsThinkingContentBlockDelta()
	case "signature_delta":
		return u.AsSignatureContentBlockDelta()
	}
	return nil
}

func (u ContentBlockDeltaEventDeltaUnion) AsTextContentBlockDelta() (v TextDelta) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ContentBlockDeltaEventDeltaUnion) AsInputJSONContentBlockDelta() (v InputJSONDelta) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ContentBlockDeltaEventDeltaUnion) AsCitationsDelta() (v CitationsDelta) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ContentBlockDeltaEventDeltaUnion) AsThinkingContentBlockDelta() (v ThinkingDelta) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ContentBlockDeltaEventDeltaUnion) AsSignatureContentBlockDelta() (v SignatureDelta) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ContentBlockDeltaEventDeltaUnion) RawJSON() string { return u.JSON.raw }

func (r *ContentBlockDeltaEventDeltaUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ContentBlockStartEvent struct {
	ContentBlock ContentBlockStartEventContentBlockUnion `json:"content_block,required"`
	Index        int64                                   `json:"index,required"`
	Type         constant.ContentBlockStart              `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ContentBlock resp.Field
		Index        resp.Field
		Type         resp.Field
		ExtraFields  map[string]resp.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ContentBlockStartEvent) RawJSON() string { return r.JSON.raw }
func (r *ContentBlockStartEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ContentBlockStartEventContentBlockUnion contains all possible properties and
// values from [TextBlock], [ToolUseBlock], [ThinkingBlock],
// [RedactedThinkingBlock].
//
// Use the [ContentBlockStartEventContentBlockUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ContentBlockStartEventContentBlockUnion struct {
	// This field is from variant [TextBlock].
	Citations []TextCitationUnion `json:"citations"`
	// This field is from variant [TextBlock].
	Text string `json:"text"`
	// Any of "text", "tool_use", "thinking", "redacted_thinking".
	Type string `json:"type"`
	// This field is from variant [ToolUseBlock].
	ID string `json:"id"`
	// This field is from variant [ToolUseBlock].
	Input interface{} `json:"input"`
	// This field is from variant [ToolUseBlock].
	Name string `json:"name"`
	// This field is from variant [ThinkingBlock].
	Signature string `json:"signature"`
	// This field is from variant [ThinkingBlock].
	Thinking string `json:"thinking"`
	// This field is from variant [RedactedThinkingBlock].
	Data string `json:"data"`
	JSON struct {
		Citations resp.Field
		Text      resp.Field
		Type      resp.Field
		ID        resp.Field
		Input     resp.Field
		Name      resp.Field
		Signature resp.Field
		Thinking  resp.Field
		Data      resp.Field
		raw       string
	} `json:"-"`
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ContentBlockStartEventContentBlockUnion.AsAny().(type) {
//	case TextBlock:
//	case ToolUseBlock:
//	case ThinkingBlock:
//	case RedactedThinkingBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ContentBlockStartEventContentBlockUnion) AsAny() any {
	switch u.Type {
	case "text":
		return u.AsResponseTextBlock()
	case "tool_use":
		return u.AsResponseToolUseBlock()
	case "thinking":
		return u.AsResponseThinkingBlock()
	case "redacted_thinking":
		return u.AsResponseRedactedThinkingBlock()
	}
	return nil
}

func (u ContentBlockStartEventContentBlockUnion) AsResponseTextBlock() (v TextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ContentBlockStartEventContentBlockUnion) AsResponseToolUseBlock() (v ToolUseBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ContentBlockStartEventContentBlockUnion) AsResponseThinkingBlock() (v ThinkingBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ContentBlockStartEventContentBlockUnion) AsResponseRedactedThinkingBlock() (v RedactedThinkingBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ContentBlockStartEventContentBlockUnion) RawJSON() string { return u.JSON.raw }

func (r *ContentBlockStartEventContentBlockUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ContentBlockStopEvent struct {
	Index int64                     `json:"index,required"`
	Type  constant.ContentBlockStop `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Index       resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ContentBlockStopEvent) RawJSON() string { return r.JSON.raw }
func (r *ContentBlockStopEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type MessageDeltaEvent struct {
	Delta MessageDeltaEventDelta `json:"delta,required"`
	Type  constant.MessageDelta  `json:"type,required"`
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
	//
	// Total input tokens in a request is the summation of `input_tokens`,
	// `cache_creation_input_tokens`, and `cache_read_input_tokens`.
	Usage MessageDeltaUsage `json:"usage,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Delta       resp.Field
		Type        resp.Field
		Usage       resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r MessageDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *MessageDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type MessageDeltaEventDelta struct {
	// Any of "end_turn", "max_tokens", "stop_sequence", "tool_use".
	StopReason   string `json:"stop_reason,required"`
	StopSequence string `json:"stop_sequence,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		StopReason   resp.Field
		StopSequence resp.Field
		ExtraFields  map[string]resp.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r MessageDeltaEventDelta) RawJSON() string { return r.JSON.raw }
func (r *MessageDeltaEventDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type MessageStartEvent struct {
	Message Message               `json:"message,required"`
	Type    constant.MessageStart `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Message     resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r MessageStartEvent) RawJSON() string { return r.JSON.raw }
func (r *MessageStartEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type MessageStopEvent struct {
	Type constant.MessageStop `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r MessageStopEvent) RawJSON() string { return r.JSON.raw }
func (r *MessageStopEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// MessageStreamEventUnion contains all possible properties and values from
// [MessageStartEvent], [MessageDeltaEvent], [MessageStopEvent],
// [ContentBlockStartEvent], [ContentBlockDeltaEvent], [ContentBlockStopEvent].
//
// Use the [MessageStreamEventUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type MessageStreamEventUnion struct {
	// This field is from variant [MessageStartEvent].
	Message Message `json:"message"`
	// Any of "message_start", "message_delta", "message_stop", "content_block_start",
	// "content_block_delta", "content_block_stop".
	Type string `json:"type"`
	// This field is a union of [MessageDeltaEventDelta],
	// [ContentBlockDeltaEventDeltaUnion]
	Delta MessageStreamEventUnionDelta `json:"delta"`
	// This field is from variant [MessageDeltaEvent].
	Usage MessageDeltaUsage `json:"usage"`
	// This field is from variant [ContentBlockStartEvent].
	ContentBlock ContentBlockStartEventContentBlockUnion `json:"content_block"`
	Index        int64                                   `json:"index"`
	JSON         struct {
		Message      resp.Field
		Type         resp.Field
		Delta        resp.Field
		Usage        resp.Field
		ContentBlock resp.Field
		Index        resp.Field
		raw          string
	} `json:"-"`
}

// Use the following switch statement to find the correct variant
//
//	switch variant := MessageStreamEventUnion.AsAny().(type) {
//	case MessageStartEvent:
//	case MessageDeltaEvent:
//	case MessageStopEvent:
//	case ContentBlockStartEvent:
//	case ContentBlockDeltaEvent:
//	case ContentBlockStopEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u MessageStreamEventUnion) AsAny() any {
	switch u.Type {
	case "message_start":
		return u.AsMessageStartEvent()
	case "message_delta":
		return u.AsMessageDeltaEvent()
	case "message_stop":
		return u.AsMessageStopEvent()
	case "content_block_start":
		return u.AsContentBlockStartEvent()
	case "content_block_delta":
		return u.AsContentBlockDeltaEvent()
	case "content_block_stop":
		return u.AsContentBlockStopEvent()
	}
	return nil
}

func (u MessageStreamEventUnion) AsMessageStartEvent() (v MessageStartEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageStreamEventUnion) AsMessageDeltaEvent() (v MessageDeltaEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageStreamEventUnion) AsMessageStopEvent() (v MessageStopEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageStreamEventUnion) AsContentBlockStartEvent() (v ContentBlockStartEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageStreamEventUnion) AsContentBlockDeltaEvent() (v ContentBlockDeltaEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u MessageStreamEventUnion) AsContentBlockStopEvent() (v ContentBlockStopEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u MessageStreamEventUnion) RawJSON() string { return u.JSON.raw }

func (r *MessageStreamEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// MessageStreamEventUnionDelta is an implicit subunion of
// [MessageStreamEventUnion]. MessageStreamEventUnionDelta provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [MessageStreamEventUnion].
type MessageStreamEventUnionDelta struct {
	// This field is from variant [MessageDeltaEventDelta].
	StopReason string `json:"stop_reason"`
	// This field is from variant [MessageDeltaEventDelta].
	StopSequence string `json:"stop_sequence"`
	// This field is from variant [ContentBlockDeltaEventDeltaUnion].
	Text string `json:"text"`
	Type string `json:"type"`
	// This field is from variant [ContentBlockDeltaEventDeltaUnion].
	PartialJSON string `json:"partial_json"`
	// This field is from variant [ContentBlockDeltaEventDeltaUnion].
	Citation CitationsDeltaCitationUnion `json:"citation"`
	// This field is from variant [ContentBlockDeltaEventDeltaUnion].
	Thinking string `json:"thinking"`
	// This field is from variant [ContentBlockDeltaEventDeltaUnion].
	Signature string `json:"signature"`
	JSON      struct {
		StopReason   resp.Field
		StopSequence resp.Field
		Text         resp.Field
		Type         resp.Field
		PartialJSON  resp.Field
		Citation     resp.Field
		Thinking     resp.Field
		Signature    resp.Field
		raw          string
	} `json:"-"`
}

func (r *MessageStreamEventUnionDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Accumulate builds up the Message incrementally from a MessageStreamEvent. The Message then can be used as
// any other Message, except with the caveat that the Message.JSON field which normally can be used to inspect
// the JSON sent over the network may not be populated fully.
//
//	message := anthropic.Message{}
//	for stream.Next() {
//		event := stream.Current()
//		message.Accumulate(event)
//	}
func (acc *Message) Accumulate(event MessageStreamEventUnion) error {
	if acc == nil {
		return fmt.Errorf("accumulate: cannot accumlate into nil Message")
	}

	switch event := event.AsAny().(type) {
	case MessageStartEvent:
		*acc = event.Message
	case MessageDeltaEvent:
		acc.StopReason = MessageStopReason(event.Delta.StopReason)
		acc.StopSequence = event.Delta.StopSequence
		acc.Usage.OutputTokens = event.Usage.OutputTokens

		// acc.JSON.StopReason = event.Delta.JSON.StopReason
		// acc.JSON.StopSequence = event.Delta.JSON.StopSequence
		// acc.Usage.JSON.OutputTokens = event.Usage.JSON.OutputTokens
	case MessageStopEvent:
		accJson, err := json.Marshal(acc)
		if err != nil {
			return fmt.Errorf("error converting content block to JSON: %w", err)
		}
		acc.JSON.raw = string(accJson)
	case ContentBlockStartEvent:
		acc.Content = append(acc.Content, ContentBlockUnion{})
		err := acc.Content[len(acc.Content)-1].UnmarshalJSON([]byte(event.ContentBlock.RawJSON()))
		if err != nil {
			return err
		}
	case ContentBlockDeltaEvent:
		if len(acc.Content) == 0 {
			return fmt.Errorf("received event of type %s but there was no content block", event.Type)
		}
		cb := &acc.Content[len(acc.Content)-1]
		switch delta := event.Delta.AsAny().(type) {
		case TextDelta:
			cb.Text += delta.Text
		case InputJSONDelta:
			if string(cb.Input) == "{}" {
				cb.Input = json.RawMessage{}
			}
			cb.Input = append(cb.Input, []byte(delta.PartialJSON)...)
		case ThinkingDelta:
			cb.Thinking += delta.Thinking
		case SignatureDelta:
			cb.Signature += delta.Signature
		}
	case ContentBlockStopEvent:
		if len(acc.Content) == 0 {
			return fmt.Errorf("received event of type %s but there was no content block", event.Type)
		}
		contentBlock := &acc.Content[len(acc.Content)-1]
		cbJson, err := json.Marshal(contentBlock)
		if err != nil {
			return fmt.Errorf("error converting content block to JSON: %w", err)
		}
		contentBlock.JSON.raw = string(cbJson)
	}

	return nil
}

type RedactedThinkingBlock struct {
	Data string                    `json:"data,required"`
	Type constant.RedactedThinking `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Data        resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r RedactedThinkingBlock) RawJSON() string { return r.JSON.raw }
func (r *RedactedThinkingBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (r RedactedThinkingBlock) ToParam() RedactedThinkingBlockParam {
	var p RedactedThinkingBlockParam
	p.Type = r.Type
	p.Data = r.Data
	return p
}

// The properties Data, Type are required.
type RedactedThinkingBlockParam struct {
	Data string `json:"data,required"`
	// This field can be elided, and will marshal its zero value as
	// "redacted_thinking".
	Type constant.RedactedThinking `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f RedactedThinkingBlockParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r RedactedThinkingBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow RedactedThinkingBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type SignatureDelta struct {
	Signature string                  `json:"signature,required"`
	Type      constant.SignatureDelta `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Signature   resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SignatureDelta) RawJSON() string { return r.JSON.raw }
func (r *SignatureDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type TextBlock struct {
	// Citations supporting the text block.
	//
	// The type of citation returned will depend on the type of document being cited.
	// Citing a PDF results in `page_location`, plain text results in `char_location`,
	// and content document results in `content_block_location`.
	Citations []TextCitationUnion `json:"citations,required"`
	Text      string              `json:"text,required"`
	Type      constant.Text       `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Citations   resp.Field
		Text        resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TextBlock) RawJSON() string { return r.JSON.raw }
func (r *TextBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (r TextBlock) ToParam() TextBlockParam {
	var p TextBlockParam
	p.Type = r.Type
	p.Text = r.Text
	p.Citations = make([]TextCitationParamUnion, len(r.Citations))
	for i, citation := range r.Citations {
		switch citationVariant := citation.AsAny().(type) {
		case CitationCharLocation:
			var citationParam CitationCharLocationParam
			citationParam.Type = citationVariant.Type
			citationParam.DocumentTitle = toParam(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
			citationParam.CitedText = citationVariant.CitedText
			citationParam.DocumentIndex = citationVariant.DocumentIndex
			citationParam.EndCharIndex = citationVariant.EndCharIndex
			citationParam.StartCharIndex = citationVariant.StartCharIndex
			p.Citations[i] = TextCitationParamUnion{OfRequestCharLocationCitation: &citationParam}
		case CitationPageLocation:
			var citationParam CitationPageLocationParam
			citationParam.Type = citationVariant.Type
			citationParam.DocumentTitle = toParam(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
			citationParam.DocumentIndex = citationVariant.DocumentIndex
			citationParam.EndPageNumber = citationVariant.EndPageNumber
			citationParam.StartPageNumber = citationVariant.StartPageNumber
			p.Citations[i] = TextCitationParamUnion{OfRequestPageLocationCitation: &citationParam}
		case CitationContentBlockLocation:
			var citationParam CitationContentBlockLocationParam
			citationParam.Type = citationVariant.Type
			citationParam.DocumentTitle = toParam(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
			citationParam.CitedText = citationVariant.CitedText
			citationParam.DocumentIndex = citationVariant.DocumentIndex
			citationParam.EndBlockIndex = citationVariant.EndBlockIndex
			citationParam.StartBlockIndex = citationVariant.StartBlockIndex
			p.Citations[i] = TextCitationParamUnion{OfRequestContentBlockLocationCitation: &citationParam}
		}
	}
	return p
}

// The properties Text, Type are required.
type TextBlockParam struct {
	Text         string                     `json:"text,required"`
	Citations    []TextCitationParamUnion   `json:"citations,omitzero"`
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as "text".
	Type constant.Text `json:"type,required"`
	paramObj
}

func NewTextBlock(text string) ContentBlockParamUnion {
	return ContentBlockParamUnion{
		OfRequestTextBlock: &TextBlockParam{
			Text: text,
		},
	}
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f TextBlockParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r TextBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow TextBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// TextCitationUnion contains all possible properties and values from
// [CitationCharLocation], [CitationPageLocation], [CitationContentBlockLocation].
//
// Use the [TextCitationUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type TextCitationUnion struct {
	CitedText     string `json:"cited_text"`
	DocumentIndex int64  `json:"document_index"`
	DocumentTitle string `json:"document_title"`
	// This field is from variant [CitationCharLocation].
	EndCharIndex int64 `json:"end_char_index"`
	// This field is from variant [CitationCharLocation].
	StartCharIndex int64 `json:"start_char_index"`
	// Any of "char_location", "page_location", "content_block_location".
	Type string `json:"type"`
	// This field is from variant [CitationPageLocation].
	EndPageNumber int64 `json:"end_page_number"`
	// This field is from variant [CitationPageLocation].
	StartPageNumber int64 `json:"start_page_number"`
	// This field is from variant [CitationContentBlockLocation].
	EndBlockIndex int64 `json:"end_block_index"`
	// This field is from variant [CitationContentBlockLocation].
	StartBlockIndex int64 `json:"start_block_index"`
	JSON            struct {
		CitedText       resp.Field
		DocumentIndex   resp.Field
		DocumentTitle   resp.Field
		EndCharIndex    resp.Field
		StartCharIndex  resp.Field
		Type            resp.Field
		EndPageNumber   resp.Field
		StartPageNumber resp.Field
		EndBlockIndex   resp.Field
		StartBlockIndex resp.Field
		raw             string
	} `json:"-"`
}

// Use the following switch statement to find the correct variant
//
//	switch variant := TextCitationUnion.AsAny().(type) {
//	case CitationCharLocation:
//	case CitationPageLocation:
//	case CitationContentBlockLocation:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u TextCitationUnion) AsAny() any {
	switch u.Type {
	case "char_location":
		return u.AsResponseCharLocationCitation()
	case "page_location":
		return u.AsResponsePageLocationCitation()
	case "content_block_location":
		return u.AsResponseContentBlockLocationCitation()
	}
	return nil
}

func (u TextCitationUnion) AsResponseCharLocationCitation() (v CitationCharLocation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u TextCitationUnion) AsResponsePageLocationCitation() (v CitationPageLocation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u TextCitationUnion) AsResponseContentBlockLocationCitation() (v CitationContentBlockLocation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u TextCitationUnion) RawJSON() string { return u.JSON.raw }

func (r *TextCitationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type TextCitationParamUnion struct {
	OfRequestCharLocationCitation         *CitationCharLocationParam         `json:",omitzero,inline"`
	OfRequestPageLocationCitation         *CitationPageLocationParam         `json:",omitzero,inline"`
	OfRequestContentBlockLocationCitation *CitationContentBlockLocationParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u TextCitationParamUnion) IsPresent() bool { return !param.IsOmitted(u) && !u.IsNull() }
func (u TextCitationParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[TextCitationParamUnion](u.OfRequestCharLocationCitation, u.OfRequestPageLocationCitation, u.OfRequestContentBlockLocationCitation)
}

func (u *TextCitationParamUnion) asAny() any {
	if !param.IsOmitted(u.OfRequestCharLocationCitation) {
		return u.OfRequestCharLocationCitation
	} else if !param.IsOmitted(u.OfRequestPageLocationCitation) {
		return u.OfRequestPageLocationCitation
	} else if !param.IsOmitted(u.OfRequestContentBlockLocationCitation) {
		return u.OfRequestContentBlockLocationCitation
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetEndCharIndex() *int64 {
	if vt := u.OfRequestCharLocationCitation; vt != nil {
		return &vt.EndCharIndex
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetStartCharIndex() *int64 {
	if vt := u.OfRequestCharLocationCitation; vt != nil {
		return &vt.StartCharIndex
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetEndPageNumber() *int64 {
	if vt := u.OfRequestPageLocationCitation; vt != nil {
		return &vt.EndPageNumber
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetStartPageNumber() *int64 {
	if vt := u.OfRequestPageLocationCitation; vt != nil {
		return &vt.StartPageNumber
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetEndBlockIndex() *int64 {
	if vt := u.OfRequestContentBlockLocationCitation; vt != nil {
		return &vt.EndBlockIndex
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetStartBlockIndex() *int64 {
	if vt := u.OfRequestContentBlockLocationCitation; vt != nil {
		return &vt.StartBlockIndex
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetCitedText() *string {
	if vt := u.OfRequestCharLocationCitation; vt != nil {
		return (*string)(&vt.CitedText)
	} else if vt := u.OfRequestPageLocationCitation; vt != nil {
		return (*string)(&vt.CitedText)
	} else if vt := u.OfRequestContentBlockLocationCitation; vt != nil {
		return (*string)(&vt.CitedText)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetDocumentIndex() *int64 {
	if vt := u.OfRequestCharLocationCitation; vt != nil {
		return (*int64)(&vt.DocumentIndex)
	} else if vt := u.OfRequestPageLocationCitation; vt != nil {
		return (*int64)(&vt.DocumentIndex)
	} else if vt := u.OfRequestContentBlockLocationCitation; vt != nil {
		return (*int64)(&vt.DocumentIndex)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetDocumentTitle() *string {
	if vt := u.OfRequestCharLocationCitation; vt != nil && vt.DocumentTitle.IsPresent() {
		return &vt.DocumentTitle.Value
	} else if vt := u.OfRequestPageLocationCitation; vt != nil && vt.DocumentTitle.IsPresent() {
		return &vt.DocumentTitle.Value
	} else if vt := u.OfRequestContentBlockLocationCitation; vt != nil && vt.DocumentTitle.IsPresent() {
		return &vt.DocumentTitle.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u TextCitationParamUnion) GetType() *string {
	if vt := u.OfRequestCharLocationCitation; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestPageLocationCitation; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestContentBlockLocationCitation; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[TextCitationParamUnion](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CitationCharLocationParam{}),
			DiscriminatorValue: "char_location",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CitationPageLocationParam{}),
			DiscriminatorValue: "page_location",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CitationContentBlockLocationParam{}),
			DiscriminatorValue: "content_block_location",
		},
	)
}

type TextDelta struct {
	Text string             `json:"text,required"`
	Type constant.TextDelta `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Text        resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r TextDelta) RawJSON() string { return r.JSON.raw }
func (r *TextDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ThinkingBlock struct {
	Signature string            `json:"signature,required"`
	Thinking  string            `json:"thinking,required"`
	Type      constant.Thinking `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Signature   resp.Field
		Thinking    resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ThinkingBlock) RawJSON() string { return r.JSON.raw }
func (r *ThinkingBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (r ThinkingBlock) ToParam() ThinkingBlockParam {
	var p ThinkingBlockParam
	p.Type = r.Type
	p.Signature = r.Signature
	p.Thinking = r.Thinking
	return p
}

// The properties Signature, Thinking, Type are required.
type ThinkingBlockParam struct {
	Signature string `json:"signature,required"`
	Thinking  string `json:"thinking,required"`
	// This field can be elided, and will marshal its zero value as "thinking".
	Type constant.Thinking `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f ThinkingBlockParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r ThinkingBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ThinkingBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The property Type is required.
type ThinkingConfigDisabledParam struct {
	// This field can be elided, and will marshal its zero value as "disabled".
	Type constant.Disabled `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f ThinkingConfigDisabledParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r ThinkingConfigDisabledParam) MarshalJSON() (data []byte, err error) {
	type shadow ThinkingConfigDisabledParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties BudgetTokens, Type are required.
type ThinkingConfigEnabledParam struct {
	// Determines how many tokens Claude can use for its internal reasoning process.
	// Larger budgets can enable more thorough analysis for complex problems, improving
	// response quality.
	//
	// Must be ≥1024 and less than `max_tokens`.
	//
	// See
	// [extended thinking](https://docs.anthropic.com/en/docs/build-with-claude/extended-thinking)
	// for details.
	BudgetTokens int64 `json:"budget_tokens,required"`
	// This field can be elided, and will marshal its zero value as "enabled".
	Type constant.Enabled `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f ThinkingConfigEnabledParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r ThinkingConfigEnabledParam) MarshalJSON() (data []byte, err error) {
	type shadow ThinkingConfigEnabledParam
	return param.MarshalObject(r, (*shadow)(&r))
}

func ThinkingConfigParamOfThinkingConfigEnabled(budgetTokens int64) ThinkingConfigParamUnion {
	var variant ThinkingConfigEnabledParam
	variant.BudgetTokens = budgetTokens
	return ThinkingConfigParamUnion{OfThinkingConfigEnabled: &variant}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ThinkingConfigParamUnion struct {
	OfThinkingConfigEnabled  *ThinkingConfigEnabledParam  `json:",omitzero,inline"`
	OfThinkingConfigDisabled *ThinkingConfigDisabledParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u ThinkingConfigParamUnion) IsPresent() bool { return !param.IsOmitted(u) && !u.IsNull() }
func (u ThinkingConfigParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[ThinkingConfigParamUnion](u.OfThinkingConfigEnabled, u.OfThinkingConfigDisabled)
}

func (u *ThinkingConfigParamUnion) asAny() any {
	if !param.IsOmitted(u.OfThinkingConfigEnabled) {
		return u.OfThinkingConfigEnabled
	} else if !param.IsOmitted(u.OfThinkingConfigDisabled) {
		return u.OfThinkingConfigDisabled
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ThinkingConfigParamUnion) GetBudgetTokens() *int64 {
	if vt := u.OfThinkingConfigEnabled; vt != nil {
		return &vt.BudgetTokens
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ThinkingConfigParamUnion) GetType() *string {
	if vt := u.OfThinkingConfigEnabled; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfThinkingConfigDisabled; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ThinkingConfigParamUnion](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ThinkingConfigEnabledParam{}),
			DiscriminatorValue: "enabled",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ThinkingConfigDisabledParam{}),
			DiscriminatorValue: "disabled",
		},
	)
}

type ThinkingDelta struct {
	Thinking string                 `json:"thinking,required"`
	Type     constant.ThinkingDelta `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		Thinking    resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ThinkingDelta) RawJSON() string { return r.JSON.raw }
func (r *ThinkingDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties InputSchema, Name are required.
type ToolParam struct {
	// [JSON schema](https://json-schema.org/draft/2020-12) for this tool's input.
	//
	// This defines the shape of the `input` that your tool accepts and that the model
	// will produce.
	InputSchema ToolInputSchemaParam `json:"input_schema,omitzero,required"`
	// Name of the tool.
	//
	// This is how the tool will be called by the model and in tool_use blocks.
	Name string `json:"name,required"`
	// Description of what this tool does.
	//
	// Tool descriptions should be as detailed as possible. The more information that
	// the model has about what the tool is and how to use it, the better it will
	// perform. You can use natural language descriptions to reinforce important
	// aspects of the tool input JSON schema.
	Description  param.Opt[string]          `json:"description,omitzero"`
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f ToolParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r ToolParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// [JSON schema](https://json-schema.org/draft/2020-12) for this tool's input.
//
// This defines the shape of the `input` that your tool accepts and that the model
// will produce.
//
// The property Type is required.
type ToolInputSchemaParam struct {
	Properties interface{} `json:"properties,omitzero"`
	// This field can be elided, and will marshal its zero value as "object".
	Type        constant.Object        `json:"type,required"`
	ExtraFields map[string]interface{} `json:"-,extras"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f ToolInputSchemaParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r ToolInputSchemaParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolInputSchemaParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties Name, Type are required.
type ToolBash20250124Param struct {
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// Name of the tool.
	//
	// This is how the tool will be called by the model and in tool_use blocks.
	//
	// This field can be elided, and will marshal its zero value as "bash".
	Name constant.Bash `json:"name,required"`
	// This field can be elided, and will marshal its zero value as "bash_20250124".
	Type constant.Bash20250124 `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f ToolBash20250124Param) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r ToolBash20250124Param) MarshalJSON() (data []byte, err error) {
	type shadow ToolBash20250124Param
	return param.MarshalObject(r, (*shadow)(&r))
}

func ToolChoiceParamOfToolChoiceTool(name string) ToolChoiceUnionParam {
	var variant ToolChoiceToolParam
	variant.Name = name
	return ToolChoiceUnionParam{OfToolChoiceTool: &variant}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ToolChoiceUnionParam struct {
	OfToolChoiceAuto *ToolChoiceAutoParam `json:",omitzero,inline"`
	OfToolChoiceAny  *ToolChoiceAnyParam  `json:",omitzero,inline"`
	OfToolChoiceTool *ToolChoiceToolParam `json:",omitzero,inline"`
	OfToolChoiceNone *ToolChoiceNoneParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u ToolChoiceUnionParam) IsPresent() bool { return !param.IsOmitted(u) && !u.IsNull() }
func (u ToolChoiceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[ToolChoiceUnionParam](u.OfToolChoiceAuto, u.OfToolChoiceAny, u.OfToolChoiceTool, u.OfToolChoiceNone)
}

func (u *ToolChoiceUnionParam) asAny() any {
	if !param.IsOmitted(u.OfToolChoiceAuto) {
		return u.OfToolChoiceAuto
	} else if !param.IsOmitted(u.OfToolChoiceAny) {
		return u.OfToolChoiceAny
	} else if !param.IsOmitted(u.OfToolChoiceTool) {
		return u.OfToolChoiceTool
	} else if !param.IsOmitted(u.OfToolChoiceNone) {
		return u.OfToolChoiceNone
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolChoiceUnionParam) GetName() *string {
	if vt := u.OfToolChoiceTool; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolChoiceUnionParam) GetType() *string {
	if vt := u.OfToolChoiceAuto; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfToolChoiceAny; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfToolChoiceTool; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfToolChoiceNone; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolChoiceUnionParam) GetDisableParallelToolUse() *bool {
	if vt := u.OfToolChoiceAuto; vt != nil && vt.DisableParallelToolUse.IsPresent() {
		return &vt.DisableParallelToolUse.Value
	} else if vt := u.OfToolChoiceAny; vt != nil && vt.DisableParallelToolUse.IsPresent() {
		return &vt.DisableParallelToolUse.Value
	} else if vt := u.OfToolChoiceTool; vt != nil && vt.DisableParallelToolUse.IsPresent() {
		return &vt.DisableParallelToolUse.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ToolChoiceUnionParam](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ToolChoiceAutoParam{}),
			DiscriminatorValue: "auto",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ToolChoiceAnyParam{}),
			DiscriminatorValue: "any",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ToolChoiceToolParam{}),
			DiscriminatorValue: "tool",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ToolChoiceNoneParam{}),
			DiscriminatorValue: "none",
		},
	)
}

// The model will use any available tools.
//
// The property Type is required.
type ToolChoiceAnyParam struct {
	// Whether to disable parallel tool use.
	//
	// Defaults to `false`. If set to `true`, the model will output exactly one tool
	// use.
	DisableParallelToolUse param.Opt[bool] `json:"disable_parallel_tool_use,omitzero"`
	// This field can be elided, and will marshal its zero value as "any".
	Type constant.Any `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f ToolChoiceAnyParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r ToolChoiceAnyParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolChoiceAnyParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The model will automatically decide whether to use tools.
//
// The property Type is required.
type ToolChoiceAutoParam struct {
	// Whether to disable parallel tool use.
	//
	// Defaults to `false`. If set to `true`, the model will output at most one tool
	// use.
	DisableParallelToolUse param.Opt[bool] `json:"disable_parallel_tool_use,omitzero"`
	// This field can be elided, and will marshal its zero value as "auto".
	Type constant.Auto `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f ToolChoiceAutoParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r ToolChoiceAutoParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolChoiceAutoParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The model will not be allowed to use tools.
//
// The property Type is required.
type ToolChoiceNoneParam struct {
	// This field can be elided, and will marshal its zero value as "none".
	Type constant.None `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f ToolChoiceNoneParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r ToolChoiceNoneParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolChoiceNoneParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The model will use the specified tool with `tool_choice.name`.
//
// The properties Name, Type are required.
type ToolChoiceToolParam struct {
	// The name of the tool to use.
	Name string `json:"name,required"`
	// Whether to disable parallel tool use.
	//
	// Defaults to `false`. If set to `true`, the model will output exactly one tool
	// use.
	DisableParallelToolUse param.Opt[bool] `json:"disable_parallel_tool_use,omitzero"`
	// This field can be elided, and will marshal its zero value as "tool".
	Type constant.Tool `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f ToolChoiceToolParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r ToolChoiceToolParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolChoiceToolParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties ToolUseID, Type are required.
type ToolResultBlockParam struct {
	ToolUseID    string                             `json:"tool_use_id,required"`
	IsError      param.Opt[bool]                    `json:"is_error,omitzero"`
	CacheControl CacheControlEphemeralParam         `json:"cache_control,omitzero"`
	Content      []ToolResultBlockParamContentUnion `json:"content,omitzero"`
	// This field can be elided, and will marshal its zero value as "tool_result".
	Type constant.ToolResult `json:"type,required"`
	paramObj
}

func NewToolResultBlock(toolUseID string, content string, isError bool) ContentBlockParamUnion {
	blockParam := ToolResultBlockParam{
		ToolUseID: toolUseID,
		Content: []ToolResultBlockParamContentUnion{{OfRequestTextBlock: &TextBlockParam{
			Text: content,
		}}},
		IsError: Bool(isError),
	}
	return ContentBlockParamUnion{OfRequestToolResultBlock: &blockParam}
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f ToolResultBlockParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r ToolResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ToolResultBlockParamContentUnion struct {
	OfRequestTextBlock  *TextBlockParam  `json:",omitzero,inline"`
	OfRequestImageBlock *ImageBlockParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u ToolResultBlockParamContentUnion) IsPresent() bool { return !param.IsOmitted(u) && !u.IsNull() }
func (u ToolResultBlockParamContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[ToolResultBlockParamContentUnion](u.OfRequestTextBlock, u.OfRequestImageBlock)
}

func (u *ToolResultBlockParamContentUnion) asAny() any {
	if !param.IsOmitted(u.OfRequestTextBlock) {
		return u.OfRequestTextBlock
	} else if !param.IsOmitted(u.OfRequestImageBlock) {
		return u.OfRequestImageBlock
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolResultBlockParamContentUnion) GetText() *string {
	if vt := u.OfRequestTextBlock; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolResultBlockParamContentUnion) GetCitations() []TextCitationParamUnion {
	if vt := u.OfRequestTextBlock; vt != nil {
		return vt.Citations
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolResultBlockParamContentUnion) GetSource() *ImageBlockParamSourceUnion {
	if vt := u.OfRequestImageBlock; vt != nil {
		return &vt.Source
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolResultBlockParamContentUnion) GetType() *string {
	if vt := u.OfRequestTextBlock; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestImageBlock; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's CacheControl property, if present.
func (u ToolResultBlockParamContentUnion) GetCacheControl() *CacheControlEphemeralParam {
	if vt := u.OfRequestTextBlock; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfRequestImageBlock; vt != nil {
		return &vt.CacheControl
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ToolResultBlockParamContentUnion](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(TextBlockParam{}),
			DiscriminatorValue: "text",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(ImageBlockParam{}),
			DiscriminatorValue: "image",
		},
	)
}

// The properties Name, Type are required.
type ToolTextEditor20250124Param struct {
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// Name of the tool.
	//
	// This is how the tool will be called by the model and in tool_use blocks.
	//
	// This field can be elided, and will marshal its zero value as
	// "str_replace_editor".
	Name constant.StrReplaceEditor `json:"name,required"`
	// This field can be elided, and will marshal its zero value as
	// "text_editor_20250124".
	Type constant.TextEditor20250124 `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f ToolTextEditor20250124Param) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r ToolTextEditor20250124Param) MarshalJSON() (data []byte, err error) {
	type shadow ToolTextEditor20250124Param
	return param.MarshalObject(r, (*shadow)(&r))
}

func ToolUnionParamOfTool(inputSchema ToolInputSchemaParam, name string) ToolUnionParam {
	var variant ToolParam
	variant.InputSchema = inputSchema
	variant.Name = name
	return ToolUnionParam{OfTool: &variant}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ToolUnionParam struct {
	OfTool               *ToolParam                   `json:",omitzero,inline"`
	OfBashTool20250124   *ToolBash20250124Param       `json:",omitzero,inline"`
	OfTextEditor20250124 *ToolTextEditor20250124Param `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u ToolUnionParam) IsPresent() bool { return !param.IsOmitted(u) && !u.IsNull() }
func (u ToolUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[ToolUnionParam](u.OfTool, u.OfBashTool20250124, u.OfTextEditor20250124)
}

func (u *ToolUnionParam) asAny() any {
	if !param.IsOmitted(u.OfTool) {
		return u.OfTool
	} else if !param.IsOmitted(u.OfBashTool20250124) {
		return u.OfBashTool20250124
	} else if !param.IsOmitted(u.OfTextEditor20250124) {
		return u.OfTextEditor20250124
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolUnionParam) GetInputSchema() *ToolInputSchemaParam {
	if vt := u.OfTool; vt != nil {
		return &vt.InputSchema
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolUnionParam) GetDescription() *string {
	if vt := u.OfTool; vt != nil && vt.Description.IsPresent() {
		return &vt.Description.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolUnionParam) GetName() *string {
	if vt := u.OfTool; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfBashTool20250124; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfTextEditor20250124; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ToolUnionParam) GetType() *string {
	if vt := u.OfBashTool20250124; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfTextEditor20250124; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's CacheControl property, if present.
func (u ToolUnionParam) GetCacheControl() *CacheControlEphemeralParam {
	if vt := u.OfTool; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfBashTool20250124; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfTextEditor20250124; vt != nil {
		return &vt.CacheControl
	}
	return nil
}

type ToolUseBlock struct {
	ID    string           `json:"id,required"`
	Input json.RawMessage  `json:"input,required"`
	Name  string           `json:"name,required"`
	Type  constant.ToolUse `json:"type,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		ID          resp.Field
		Input       resp.Field
		Name        resp.Field
		Type        resp.Field
		ExtraFields map[string]resp.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ToolUseBlock) RawJSON() string { return r.JSON.raw }
func (r *ToolUseBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (r ToolUseBlock) ToParam() ToolUseBlockParam {
	var toolUse ToolUseBlockParam
	toolUse.Type = r.Type
	toolUse.ID = r.ID
	toolUse.Input = r.Input
	toolUse.Name = r.Name
	return toolUse
}

// The properties ID, Input, Name, Type are required.
type ToolUseBlockParam struct {
	ID           string                     `json:"id,required"`
	Input        interface{}                `json:"input,omitzero,required"`
	Name         string                     `json:"name,required"`
	CacheControl CacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as "tool_use".
	Type constant.ToolUse `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f ToolUseBlockParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r ToolUseBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ToolUseBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties Type, URL are required.
type URLImageSourceParam struct {
	URL string `json:"url,required"`
	// This field can be elided, and will marshal its zero value as "url".
	Type constant.URL `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f URLImageSourceParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r URLImageSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow URLImageSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties Type, URL are required.
type URLPDFSourceParam struct {
	URL string `json:"url,required"`
	// This field can be elided, and will marshal its zero value as "url".
	Type constant.URL `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f URLPDFSourceParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r URLPDFSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow URLPDFSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type Usage struct {
	// The number of input tokens used to create the cache entry.
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens,required"`
	// The number of input tokens read from the cache.
	CacheReadInputTokens int64 `json:"cache_read_input_tokens,required"`
	// The number of input tokens which were used.
	InputTokens int64 `json:"input_tokens,required"`
	// The number of output tokens which were used.
	OutputTokens int64 `json:"output_tokens,required"`
	// Metadata for the response, check the presence of optional fields with the
	// [resp.Field.IsPresent] method.
	JSON struct {
		CacheCreationInputTokens resp.Field
		CacheReadInputTokens     resp.Field
		InputTokens              resp.Field
		OutputTokens             resp.Field
		ExtraFields              map[string]resp.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Usage) RawJSON() string { return r.JSON.raw }
func (r *Usage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type MessageNewParams struct {
	// The maximum number of tokens to generate before stopping.
	//
	// Note that our models may stop _before_ reaching this maximum. This parameter
	// only specifies the absolute maximum number of tokens to generate.
	//
	// Different models have different maximum values for this parameter. See
	// [models](https://docs.anthropic.com/en/docs/models-overview) for details.
	MaxTokens int64 `json:"max_tokens,required"`
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
	Messages []MessageParam `json:"messages,omitzero,required"`
	// The model that will complete your prompt.\n\nSee
	// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	Model Model `json:"model,omitzero,required"`
	// Amount of randomness injected into the response.
	//
	// Defaults to `1.0`. Ranges from `0.0` to `1.0`. Use `temperature` closer to `0.0`
	// for analytical / multiple choice, and closer to `1.0` for creative and
	// generative tasks.
	//
	// Note that even with `temperature` of `0.0`, the results will not be fully
	// deterministic.
	Temperature param.Opt[float64] `json:"temperature,omitzero"`
	// Only sample from the top K options for each subsequent token.
	//
	// Used to remove "long tail" low probability responses.
	// [Learn more technical details here](https://towardsdatascience.com/how-to-sample-from-language-models-682bceb97277).
	//
	// Recommended for advanced use cases only. You usually only need to use
	// `temperature`.
	TopK param.Opt[int64] `json:"top_k,omitzero"`
	// Use nucleus sampling.
	//
	// In nucleus sampling, we compute the cumulative distribution over all the options
	// for each subsequent token in decreasing probability order and cut it off once it
	// reaches a particular probability specified by `top_p`. You should either alter
	// `temperature` or `top_p`, but not both.
	//
	// Recommended for advanced use cases only. You usually only need to use
	// `temperature`.
	TopP param.Opt[float64] `json:"top_p,omitzero"`
	// An object describing metadata about the request.
	Metadata MetadataParam `json:"metadata,omitzero"`
	// Custom text sequences that will cause the model to stop generating.
	//
	// Our models will normally stop when they have naturally completed their turn,
	// which will result in a response `stop_reason` of `"end_turn"`.
	//
	// If you want the model to stop generating when it encounters custom strings of
	// text, you can use the `stop_sequences` parameter. If the model encounters one of
	// the custom sequences, the response `stop_reason` value will be `"stop_sequence"`
	// and the response `stop_sequence` value will contain the matched stop sequence.
	StopSequences []string `json:"stop_sequences,omitzero"`
	// System prompt.
	//
	// A system prompt is a way of providing context and instructions to Claude, such
	// as specifying a particular goal or role. See our
	// [guide to system prompts](https://docs.anthropic.com/en/docs/system-prompts).
	System []TextBlockParam `json:"system,omitzero"`
	// Configuration for enabling Claude's extended thinking.
	//
	// When enabled, responses include `thinking` content blocks showing Claude's
	// thinking process before the final answer. Requires a minimum budget of 1,024
	// tokens and counts towards your `max_tokens` limit.
	//
	// See
	// [extended thinking](https://docs.anthropic.com/en/docs/build-with-claude/extended-thinking)
	// for details.
	Thinking ThinkingConfigParamUnion `json:"thinking,omitzero"`
	// How the model should use the provided tools. The model can use a specific tool,
	// any available tool, decide by itself, or not use tools at all.
	ToolChoice ToolChoiceUnionParam `json:"tool_choice,omitzero"`
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
	//   - `input_schema`: [JSON schema](https://json-schema.org/draft/2020-12) for the
	//     tool `input` shape that the model will produce in `tool_use` output content
	//     blocks.
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
	Tools []ToolUnionParam `json:"tools,omitzero"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f MessageNewParams) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }

func (r MessageNewParams) MarshalJSON() (data []byte, err error) {
	type shadow MessageNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}

type MessageCountTokensParams struct {
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
	Messages []MessageParam `json:"messages,omitzero,required"`
	// The model that will complete your prompt.\n\nSee
	// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	Model Model `json:"model,omitzero,required"`
	// System prompt.
	//
	// A system prompt is a way of providing context and instructions to Claude, such
	// as specifying a particular goal or role. See our
	// [guide to system prompts](https://docs.anthropic.com/en/docs/system-prompts).
	System MessageCountTokensParamsSystemUnion `json:"system,omitzero"`
	// Configuration for enabling Claude's extended thinking.
	//
	// When enabled, responses include `thinking` content blocks showing Claude's
	// thinking process before the final answer. Requires a minimum budget of 1,024
	// tokens and counts towards your `max_tokens` limit.
	//
	// See
	// [extended thinking](https://docs.anthropic.com/en/docs/build-with-claude/extended-thinking)
	// for details.
	Thinking ThinkingConfigParamUnion `json:"thinking,omitzero"`
	// How the model should use the provided tools. The model can use a specific tool,
	// any available tool, decide by itself, or not use tools at all.
	ToolChoice ToolChoiceUnionParam `json:"tool_choice,omitzero"`
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
	//   - `input_schema`: [JSON schema](https://json-schema.org/draft/2020-12) for the
	//     tool `input` shape that the model will produce in `tool_use` output content
	//     blocks.
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
	Tools []MessageCountTokensToolUnionParam `json:"tools,omitzero"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f MessageCountTokensParams) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }

func (r MessageCountTokensParams) MarshalJSON() (data []byte, err error) {
	type shadow MessageCountTokensParams
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type MessageCountTokensParamsSystemUnion struct {
	OfString                         param.Opt[string] `json:",omitzero,inline"`
	OfMessageCountTokenssSystemArray []TextBlockParam  `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u MessageCountTokensParamsSystemUnion) IsPresent() bool {
	return !param.IsOmitted(u) && !u.IsNull()
}
func (u MessageCountTokensParamsSystemUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[MessageCountTokensParamsSystemUnion](u.OfString, u.OfMessageCountTokenssSystemArray)
}

func (u *MessageCountTokensParamsSystemUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfMessageCountTokenssSystemArray) {
		return &u.OfMessageCountTokenssSystemArray
	}
	return nil
}
