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

// BetaMessageService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaMessageService] method instead.
type BetaMessageService struct {
	Options []option.RequestOption
	Batches BetaMessageBatchService
}

// NewBetaMessageService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaMessageService(opts ...option.RequestOption) (r BetaMessageService) {
	r = BetaMessageService{}
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
// Learn more about the Messages API in our [user guide](/en/docs/initial-setup)
//
// Note: If you choose to set a timeout for this request, we recommend 10 minutes.
func (r *BetaMessageService) New(ctx context.Context, params BetaMessageNewParams, opts ...option.RequestOption) (res *BetaMessage, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%s", v)))
	}
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
// Learn more about the Messages API in our [user guide](/en/docs/initial-setup)
//
// Note: If you choose to set a timeout for this request, we recommend 10 minutes.
func (r *BetaMessageService) NewStreaming(ctx context.Context, params BetaMessageNewParams, opts ...option.RequestOption) (stream *ssestream.Stream[BetaRawMessageStreamEventUnion]) {
	var (
		raw *http.Response
		err error
	)
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%s", v)))
	}
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithJSONSet("stream", true)}, opts...)
	path := "v1/messages?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &raw, opts...)
	return ssestream.NewStream[BetaRawMessageStreamEventUnion](ssestream.NewDecoder(raw), err)
}

// Count the number of tokens in a Message.
//
// The Token Count API can be used to count the number of tokens in a Message,
// including tools, images, and documents, without creating it.
//
// Learn more about token counting in our
// [user guide](/en/docs/build-with-claude/token-counting)
func (r *BetaMessageService) CountTokens(ctx context.Context, params BetaMessageCountTokensParams, opts ...option.RequestOption) (res *BetaMessageTokensCount, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%s", v)))
	}
	opts = append(r.Options[:], opts...)
	path := "v1/messages/count_tokens?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return
}

// The properties Data, MediaType, Type are required.
type BetaBase64ImageSourceParam struct {
	Data string `json:"data,required" format:"byte"`
	// Any of "image/jpeg", "image/png", "image/gif", "image/webp".
	MediaType BetaBase64ImageSourceMediaType `json:"media_type,omitzero,required"`
	// This field can be elided, and will marshal its zero value as "base64".
	Type constant.Base64 `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaBase64ImageSourceParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaBase64ImageSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaBase64ImageSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaBase64ImageSourceMediaType string

const (
	BetaBase64ImageSourceMediaTypeImageJPEG BetaBase64ImageSourceMediaType = "image/jpeg"
	BetaBase64ImageSourceMediaTypeImagePNG  BetaBase64ImageSourceMediaType = "image/png"
	BetaBase64ImageSourceMediaTypeImageGIF  BetaBase64ImageSourceMediaType = "image/gif"
	BetaBase64ImageSourceMediaTypeImageWebP BetaBase64ImageSourceMediaType = "image/webp"
)

// The properties Source, Type are required.
type BetaBase64PDFBlockParam struct {
	Source       BetaBase64PDFBlockSourceUnionParam `json:"source,omitzero,required"`
	Context      param.Opt[string]                  `json:"context,omitzero"`
	Title        param.Opt[string]                  `json:"title,omitzero"`
	CacheControl BetaCacheControlEphemeralParam     `json:"cache_control,omitzero"`
	Citations    BetaCitationsConfigParam           `json:"citations,omitzero"`
	// This field can be elided, and will marshal its zero value as "document".
	Type constant.Document `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaBase64PDFBlockParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaBase64PDFBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaBase64PDFBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaBase64PDFBlockSourceUnionParam struct {
	OfBase64PDFSource    *BetaBase64PDFSourceParam    `json:",omitzero,inline"`
	OfPlainTextSource    *BetaPlainTextSourceParam    `json:",omitzero,inline"`
	OfContentBlockSource *BetaContentBlockSourceParam `json:",omitzero,inline"`
	OfUrlpdfSource       *BetaURLPDFSourceParam       `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u BetaBase64PDFBlockSourceUnionParam) IsPresent() bool {
	return !param.IsOmitted(u) && !u.IsNull()
}
func (u BetaBase64PDFBlockSourceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[BetaBase64PDFBlockSourceUnionParam](u.OfBase64PDFSource, u.OfPlainTextSource, u.OfContentBlockSource, u.OfUrlpdfSource)
}

func (u *BetaBase64PDFBlockSourceUnionParam) asAny() any {
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
func (u BetaBase64PDFBlockSourceUnionParam) GetContent() *BetaContentBlockSourceContentUnionParam {
	if vt := u.OfContentBlockSource; vt != nil {
		return &vt.Content
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaBase64PDFBlockSourceUnionParam) GetURL() *string {
	if vt := u.OfUrlpdfSource; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaBase64PDFBlockSourceUnionParam) GetData() *string {
	if vt := u.OfBase64PDFSource; vt != nil {
		return (*string)(&vt.Data)
	} else if vt := u.OfPlainTextSource; vt != nil {
		return (*string)(&vt.Data)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaBase64PDFBlockSourceUnionParam) GetMediaType() *string {
	if vt := u.OfBase64PDFSource; vt != nil {
		return (*string)(&vt.MediaType)
	} else if vt := u.OfPlainTextSource; vt != nil {
		return (*string)(&vt.MediaType)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaBase64PDFBlockSourceUnionParam) GetType() *string {
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
	apijson.RegisterUnion[BetaBase64PDFBlockSourceUnionParam](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaBase64PDFSourceParam{}),
			DiscriminatorValue: "base64",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaPlainTextSourceParam{}),
			DiscriminatorValue: "text",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaContentBlockSourceParam{}),
			DiscriminatorValue: "content",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaURLPDFSourceParam{}),
			DiscriminatorValue: "url",
		},
	)
}

// The properties Data, MediaType, Type are required.
type BetaBase64PDFSourceParam struct {
	Data string `json:"data,required" format:"byte"`
	// This field can be elided, and will marshal its zero value as "application/pdf".
	MediaType constant.ApplicationPDF `json:"media_type,required"`
	// This field can be elided, and will marshal its zero value as "base64".
	Type constant.Base64 `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaBase64PDFSourceParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaBase64PDFSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaBase64PDFSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The property Type is required.
type BetaCacheControlEphemeralParam struct {
	// This field can be elided, and will marshal its zero value as "ephemeral".
	Type constant.Ephemeral `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaCacheControlEphemeralParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaCacheControlEphemeralParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaCacheControlEphemeralParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaCitationCharLocation struct {
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
func (r BetaCitationCharLocation) RawJSON() string { return r.JSON.raw }
func (r *BetaCitationCharLocation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties CitedText, DocumentIndex, DocumentTitle, EndCharIndex,
// StartCharIndex, Type are required.
type BetaCitationCharLocationParam struct {
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
func (f BetaCitationCharLocationParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaCitationCharLocationParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaCitationCharLocationParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaCitationContentBlockLocation struct {
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
func (r BetaCitationContentBlockLocation) RawJSON() string { return r.JSON.raw }
func (r *BetaCitationContentBlockLocation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties CitedText, DocumentIndex, DocumentTitle, EndBlockIndex,
// StartBlockIndex, Type are required.
type BetaCitationContentBlockLocationParam struct {
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
func (f BetaCitationContentBlockLocationParam) IsPresent() bool {
	return !param.IsOmitted(f) && !f.IsNull()
}
func (r BetaCitationContentBlockLocationParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaCitationContentBlockLocationParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaCitationPageLocation struct {
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
func (r BetaCitationPageLocation) RawJSON() string { return r.JSON.raw }
func (r *BetaCitationPageLocation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties CitedText, DocumentIndex, DocumentTitle, EndPageNumber,
// StartPageNumber, Type are required.
type BetaCitationPageLocationParam struct {
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
func (f BetaCitationPageLocationParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaCitationPageLocationParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaCitationPageLocationParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaCitationsConfigParam struct {
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaCitationsConfigParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaCitationsConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaCitationsConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaCitationsDelta struct {
	Citation BetaCitationsDeltaCitationUnion `json:"citation,required"`
	Type     constant.CitationsDelta         `json:"type,required"`
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
func (r BetaCitationsDelta) RawJSON() string { return r.JSON.raw }
func (r *BetaCitationsDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaCitationsDeltaCitationUnion contains all possible properties and values from
// [BetaCitationCharLocation], [BetaCitationPageLocation],
// [BetaCitationContentBlockLocation].
//
// Use the [BetaCitationsDeltaCitationUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaCitationsDeltaCitationUnion struct {
	CitedText     string `json:"cited_text"`
	DocumentIndex int64  `json:"document_index"`
	DocumentTitle string `json:"document_title"`
	// This field is from variant [BetaCitationCharLocation].
	EndCharIndex int64 `json:"end_char_index"`
	// This field is from variant [BetaCitationCharLocation].
	StartCharIndex int64 `json:"start_char_index"`
	// Any of "char_location", "page_location", "content_block_location".
	Type string `json:"type"`
	// This field is from variant [BetaCitationPageLocation].
	EndPageNumber int64 `json:"end_page_number"`
	// This field is from variant [BetaCitationPageLocation].
	StartPageNumber int64 `json:"start_page_number"`
	// This field is from variant [BetaCitationContentBlockLocation].
	EndBlockIndex int64 `json:"end_block_index"`
	// This field is from variant [BetaCitationContentBlockLocation].
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
//	switch variant := BetaCitationsDeltaCitationUnion.AsAny().(type) {
//	case BetaCitationCharLocation:
//	case BetaCitationPageLocation:
//	case BetaCitationContentBlockLocation:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaCitationsDeltaCitationUnion) AsAny() any {
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

func (u BetaCitationsDeltaCitationUnion) AsResponseCharLocationCitation() (v BetaCitationCharLocation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaCitationsDeltaCitationUnion) AsResponsePageLocationCitation() (v BetaCitationPageLocation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaCitationsDeltaCitationUnion) AsResponseContentBlockLocationCitation() (v BetaCitationContentBlockLocation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaCitationsDeltaCitationUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaCitationsDeltaCitationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaContentBlockUnion contains all possible properties and values from
// [BetaTextBlock], [BetaToolUseBlock], [BetaThinkingBlock],
// [BetaRedactedThinkingBlock].
//
// Use the [BetaContentBlockUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaContentBlockUnion struct {
	// This field is from variant [BetaTextBlock].
	Citations []BetaTextCitationUnion `json:"citations"`
	// This field is from variant [BetaTextBlock].
	Text string `json:"text"`
	// Any of "text", "tool_use", "thinking", "redacted_thinking".
	Type string `json:"type"`
	// This field is from variant [BetaToolUseBlock].
	ID string `json:"id"`
	// This field is from variant [BetaToolUseBlock].
	Input interface{} `json:"input"`
	// This field is from variant [BetaToolUseBlock].
	Name string `json:"name"`
	// This field is from variant [BetaThinkingBlock].
	Signature string `json:"signature"`
	// This field is from variant [BetaThinkingBlock].
	Thinking string `json:"thinking"`
	// This field is from variant [BetaRedactedThinkingBlock].
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

func (r BetaContentBlockUnion) ToParam() BetaContentBlockParamUnion {
	switch variant := r.AsAny().(type) {
	case BetaTextBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfRequestTextBlock: &p}
	case BetaToolUseBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfRequestToolUseBlock: &p}
	case BetaThinkingBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfRequestThinkingBlock: &p}
	case BetaRedactedThinkingBlock:
		p := variant.ToParam()
		return BetaContentBlockParamUnion{OfRequestRedactedThinkingBlock: &p}
	}
	return BetaContentBlockParamUnion{}
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaContentBlockUnion.AsAny().(type) {
//	case BetaTextBlock:
//	case BetaToolUseBlock:
//	case BetaThinkingBlock:
//	case BetaRedactedThinkingBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaContentBlockUnion) AsAny() any {
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

func (u BetaContentBlockUnion) AsResponseTextBlock() (v BetaTextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaContentBlockUnion) AsResponseToolUseBlock() (v BetaToolUseBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaContentBlockUnion) AsResponseThinkingBlock() (v BetaThinkingBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaContentBlockUnion) AsResponseRedactedThinkingBlock() (v BetaRedactedThinkingBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaContentBlockUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaContentBlockUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func BetaContentBlockParamOfRequestTextBlock(text string) BetaContentBlockParamUnion {
	var variant BetaTextBlockParam
	variant.Text = text
	return BetaContentBlockParamUnion{OfRequestTextBlock: &variant}
}

func BetaContentBlockParamOfRequestImageBlock[
	T BetaBase64ImageSourceParam | BetaURLImageSourceParam,
](source T) BetaContentBlockParamUnion {
	var variant BetaImageBlockParam
	switch v := any(source).(type) {
	case BetaBase64ImageSourceParam:
		variant.Source.OfBase64ImageSource = &v
	case BetaURLImageSourceParam:
		variant.Source.OfURLImageSource = &v
	}
	return BetaContentBlockParamUnion{OfRequestImageBlock: &variant}
}

func BetaContentBlockParamOfRequestToolUseBlock(id string, input interface{}, name string) BetaContentBlockParamUnion {
	var variant BetaToolUseBlockParam
	variant.ID = id
	variant.Input = input
	variant.Name = name
	return BetaContentBlockParamUnion{OfRequestToolUseBlock: &variant}
}

func BetaContentBlockParamOfRequestToolResultBlock(toolUseID string) BetaContentBlockParamUnion {
	var variant BetaToolResultBlockParam
	variant.ToolUseID = toolUseID
	return BetaContentBlockParamUnion{OfRequestToolResultBlock: &variant}
}

func BetaContentBlockParamOfRequestDocumentBlock[
	T BetaBase64PDFSourceParam | BetaPlainTextSourceParam | BetaContentBlockSourceParam | BetaURLPDFSourceParam,
](source T) BetaContentBlockParamUnion {
	var variant BetaBase64PDFBlockParam
	switch v := any(source).(type) {
	case BetaBase64PDFSourceParam:
		variant.Source.OfBase64PDFSource = &v
	case BetaPlainTextSourceParam:
		variant.Source.OfPlainTextSource = &v
	case BetaContentBlockSourceParam:
		variant.Source.OfContentBlockSource = &v
	case BetaURLPDFSourceParam:
		variant.Source.OfUrlpdfSource = &v
	}
	return BetaContentBlockParamUnion{OfRequestDocumentBlock: &variant}
}

func BetaContentBlockParamOfRequestThinkingBlock(signature string, thinking string) BetaContentBlockParamUnion {
	var variant BetaThinkingBlockParam
	variant.Signature = signature
	variant.Thinking = thinking
	return BetaContentBlockParamUnion{OfRequestThinkingBlock: &variant}
}

func BetaContentBlockParamOfRequestRedactedThinkingBlock(data string) BetaContentBlockParamUnion {
	var variant BetaRedactedThinkingBlockParam
	variant.Data = data
	return BetaContentBlockParamUnion{OfRequestRedactedThinkingBlock: &variant}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaContentBlockParamUnion struct {
	OfRequestTextBlock             *BetaTextBlockParam             `json:",omitzero,inline"`
	OfRequestImageBlock            *BetaImageBlockParam            `json:",omitzero,inline"`
	OfRequestToolUseBlock          *BetaToolUseBlockParam          `json:",omitzero,inline"`
	OfRequestToolResultBlock       *BetaToolResultBlockParam       `json:",omitzero,inline"`
	OfRequestDocumentBlock         *BetaBase64PDFBlockParam        `json:",omitzero,inline"`
	OfRequestThinkingBlock         *BetaThinkingBlockParam         `json:",omitzero,inline"`
	OfRequestRedactedThinkingBlock *BetaRedactedThinkingBlockParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u BetaContentBlockParamUnion) IsPresent() bool { return !param.IsOmitted(u) && !u.IsNull() }
func (u BetaContentBlockParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[BetaContentBlockParamUnion](u.OfRequestTextBlock,
		u.OfRequestImageBlock,
		u.OfRequestToolUseBlock,
		u.OfRequestToolResultBlock,
		u.OfRequestDocumentBlock,
		u.OfRequestThinkingBlock,
		u.OfRequestRedactedThinkingBlock)
}

func (u *BetaContentBlockParamUnion) asAny() any {
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
func (u BetaContentBlockParamUnion) GetText() *string {
	if vt := u.OfRequestTextBlock; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaContentBlockParamUnion) GetID() *string {
	if vt := u.OfRequestToolUseBlock; vt != nil {
		return &vt.ID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaContentBlockParamUnion) GetInput() *interface{} {
	if vt := u.OfRequestToolUseBlock; vt != nil {
		return &vt.Input
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaContentBlockParamUnion) GetName() *string {
	if vt := u.OfRequestToolUseBlock; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaContentBlockParamUnion) GetToolUseID() *string {
	if vt := u.OfRequestToolResultBlock; vt != nil {
		return &vt.ToolUseID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaContentBlockParamUnion) GetContent() *[]BetaToolResultBlockParamContentUnion {
	if vt := u.OfRequestToolResultBlock; vt != nil {
		return &vt.Content
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaContentBlockParamUnion) GetIsError() *bool {
	if vt := u.OfRequestToolResultBlock; vt != nil && vt.IsError.IsPresent() {
		return &vt.IsError.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaContentBlockParamUnion) GetContext() *string {
	if vt := u.OfRequestDocumentBlock; vt != nil && vt.Context.IsPresent() {
		return &vt.Context.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaContentBlockParamUnion) GetTitle() *string {
	if vt := u.OfRequestDocumentBlock; vt != nil && vt.Title.IsPresent() {
		return &vt.Title.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaContentBlockParamUnion) GetSignature() *string {
	if vt := u.OfRequestThinkingBlock; vt != nil {
		return &vt.Signature
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaContentBlockParamUnion) GetThinking() *string {
	if vt := u.OfRequestThinkingBlock; vt != nil {
		return &vt.Thinking
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaContentBlockParamUnion) GetData() *string {
	if vt := u.OfRequestRedactedThinkingBlock; vt != nil {
		return &vt.Data
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaContentBlockParamUnion) GetType() *string {
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
func (u BetaContentBlockParamUnion) GetCacheControl() *BetaCacheControlEphemeralParam {
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
func (u BetaContentBlockParamUnion) GetCitations() (res betaContentBlockParamUnionCitations) {
	if vt := u.OfRequestTextBlock; vt != nil {
		res.ofBetaTextBlockCitations = &vt.Citations
	} else if vt := u.OfRequestDocumentBlock; vt != nil {
		res.ofBetaCitationsConfig = &vt.Citations
	}
	return
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type betaContentBlockParamUnionCitations struct {
	ofBetaTextBlockCitations *[]BetaTextCitationParamUnion
	ofBetaCitationsConfig    *BetaCitationsConfigParam
}

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]anthropic.BetaTextCitationParamUnion:
//	case *anthropic.BetaCitationsConfigParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u betaContentBlockParamUnionCitations) AsAny() any {
	if !param.IsOmitted(u.ofBetaTextBlockCitations) {
		return u.ofBetaTextBlockCitations
	} else if !param.IsOmitted(u.ofBetaCitationsConfig) {
		return u.ofBetaCitationsConfig
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u betaContentBlockParamUnionCitations) GetEnabled() *bool {
	if vt := u.ofBetaCitationsConfig; vt != nil && vt.Enabled.IsPresent() {
		return &vt.Enabled.Value
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u BetaContentBlockParamUnion) GetSource() (res betaContentBlockParamUnionSource) {
	if vt := u.OfRequestImageBlock; vt != nil {
		res.ofBetaImageBlockSource = &vt.Source
	} else if vt := u.OfRequestDocumentBlock; vt != nil {
		res.ofBetaBase64PDFBlockSourceUnion = &vt.Source
	}
	return
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type betaContentBlockParamUnionSource struct {
	ofBetaImageBlockSource          *BetaImageBlockParamSourceUnion
	ofBetaBase64PDFBlockSourceUnion *BetaBase64PDFBlockSourceUnionParam
}

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *anthropic.BetaBase64ImageSourceParam:
//	case *anthropic.BetaURLImageSourceParam:
//	case *anthropic.BetaBase64PDFSourceParam:
//	case *anthropic.BetaPlainTextSourceParam:
//	case *anthropic.BetaContentBlockSourceParam:
//	case *anthropic.BetaURLPDFSourceParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u betaContentBlockParamUnionSource) AsAny() any {
	if !param.IsOmitted(u.ofBetaImageBlockSource) {
		return u.ofBetaImageBlockSource.asAny()
	} else if !param.IsOmitted(u.ofBetaBase64PDFBlockSourceUnion) {
		return u.ofBetaBase64PDFBlockSourceUnion.asAny()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u betaContentBlockParamUnionSource) GetContent() *BetaContentBlockSourceContentUnionParam {
	if u.ofBetaBase64PDFBlockSourceUnion != nil {
		return u.ofBetaBase64PDFBlockSourceUnion.GetContent()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u betaContentBlockParamUnionSource) GetData() *string {
	if u.ofBetaImageBlockSource != nil {
		return u.ofBetaImageBlockSource.GetData()
	} else if u.ofBetaBase64PDFBlockSourceUnion != nil {
		return u.ofBetaBase64PDFBlockSourceUnion.GetData()
	} else if u.ofBetaBase64PDFBlockSourceUnion != nil {
		return u.ofBetaBase64PDFBlockSourceUnion.GetData()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u betaContentBlockParamUnionSource) GetMediaType() *string {
	if u.ofBetaImageBlockSource != nil {
		return u.ofBetaImageBlockSource.GetMediaType()
	} else if u.ofBetaBase64PDFBlockSourceUnion != nil {
		return u.ofBetaBase64PDFBlockSourceUnion.GetMediaType()
	} else if u.ofBetaBase64PDFBlockSourceUnion != nil {
		return u.ofBetaBase64PDFBlockSourceUnion.GetMediaType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u betaContentBlockParamUnionSource) GetType() *string {
	if u.ofBetaImageBlockSource != nil {
		return u.ofBetaImageBlockSource.GetType()
	} else if u.ofBetaImageBlockSource != nil {
		return u.ofBetaImageBlockSource.GetType()
	} else if u.ofBetaBase64PDFBlockSourceUnion != nil {
		return u.ofBetaBase64PDFBlockSourceUnion.GetType()
	} else if u.ofBetaBase64PDFBlockSourceUnion != nil {
		return u.ofBetaBase64PDFBlockSourceUnion.GetType()
	} else if u.ofBetaBase64PDFBlockSourceUnion != nil {
		return u.ofBetaBase64PDFBlockSourceUnion.GetType()
	} else if u.ofBetaBase64PDFBlockSourceUnion != nil {
		return u.ofBetaBase64PDFBlockSourceUnion.GetType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u betaContentBlockParamUnionSource) GetURL() *string {
	if u.ofBetaImageBlockSource != nil {
		return u.ofBetaImageBlockSource.GetURL()
	} else if u.ofBetaBase64PDFBlockSourceUnion != nil {
		return u.ofBetaBase64PDFBlockSourceUnion.GetURL()
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaContentBlockParamUnion](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaTextBlockParam{}),
			DiscriminatorValue: "text",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaImageBlockParam{}),
			DiscriminatorValue: "image",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaToolUseBlockParam{}),
			DiscriminatorValue: "tool_use",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaToolResultBlockParam{}),
			DiscriminatorValue: "tool_result",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaBase64PDFBlockParam{}),
			DiscriminatorValue: "document",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaThinkingBlockParam{}),
			DiscriminatorValue: "thinking",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaRedactedThinkingBlockParam{}),
			DiscriminatorValue: "redacted_thinking",
		},
	)
}

// The properties Content, Type are required.
type BetaContentBlockSourceParam struct {
	Content BetaContentBlockSourceContentUnionParam `json:"content,omitzero,required"`
	// This field can be elided, and will marshal its zero value as "content".
	Type constant.Content `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaContentBlockSourceParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaContentBlockSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaContentBlockSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaContentBlockSourceContentUnionParam struct {
	OfString                        param.Opt[string]                         `json:",omitzero,inline"`
	OfBetaContentBlockSourceContent []BetaContentBlockSourceContentUnionParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u BetaContentBlockSourceContentUnionParam) IsPresent() bool {
	return !param.IsOmitted(u) && !u.IsNull()
}
func (u BetaContentBlockSourceContentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[BetaContentBlockSourceContentUnionParam](u.OfString, u.OfBetaContentBlockSourceContent)
}

func (u *BetaContentBlockSourceContentUnionParam) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfBetaContentBlockSourceContent) {
		return &u.OfBetaContentBlockSourceContent
	}
	return nil
}

// The properties Source, Type are required.
type BetaImageBlockParam struct {
	Source       BetaImageBlockParamSourceUnion `json:"source,omitzero,required"`
	CacheControl BetaCacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as "image".
	Type constant.Image `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaImageBlockParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaImageBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaImageBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaImageBlockParamSourceUnion struct {
	OfBase64ImageSource *BetaBase64ImageSourceParam `json:",omitzero,inline"`
	OfURLImageSource    *BetaURLImageSourceParam    `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u BetaImageBlockParamSourceUnion) IsPresent() bool { return !param.IsOmitted(u) && !u.IsNull() }
func (u BetaImageBlockParamSourceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[BetaImageBlockParamSourceUnion](u.OfBase64ImageSource, u.OfURLImageSource)
}

func (u *BetaImageBlockParamSourceUnion) asAny() any {
	if !param.IsOmitted(u.OfBase64ImageSource) {
		return u.OfBase64ImageSource
	} else if !param.IsOmitted(u.OfURLImageSource) {
		return u.OfURLImageSource
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaImageBlockParamSourceUnion) GetData() *string {
	if vt := u.OfBase64ImageSource; vt != nil {
		return &vt.Data
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaImageBlockParamSourceUnion) GetMediaType() *string {
	if vt := u.OfBase64ImageSource; vt != nil {
		return (*string)(&vt.MediaType)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaImageBlockParamSourceUnion) GetURL() *string {
	if vt := u.OfURLImageSource; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaImageBlockParamSourceUnion) GetType() *string {
	if vt := u.OfBase64ImageSource; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfURLImageSource; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaImageBlockParamSourceUnion](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaBase64ImageSourceParam{}),
			DiscriminatorValue: "base64",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaURLImageSourceParam{}),
			DiscriminatorValue: "url",
		},
	)
}

type BetaInputJSONDelta struct {
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
func (r BetaInputJSONDelta) RawJSON() string { return r.JSON.raw }
func (r *BetaInputJSONDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
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
	Content []BetaContentBlockUnion `json:"content,required"`
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
	StopReason BetaMessageStopReason `json:"stop_reason,required"`
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
	Usage BetaUsage `json:"usage,required"`
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
func (r BetaMessage) RawJSON() string { return r.JSON.raw }
func (r *BetaMessage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (r BetaMessage) ToParam() BetaMessageParam {
	var p BetaMessageParam
	p.Role = BetaMessageParamRole(r.Role)
	p.Content = make([]BetaContentBlockParamUnion, len(r.Content))
	for i, c := range r.Content {
		contentParams := c.ToParam()
		p.Content[i] = contentParams
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
type BetaMessageStopReason string

const (
	BetaMessageStopReasonEndTurn      BetaMessageStopReason = "end_turn"
	BetaMessageStopReasonMaxTokens    BetaMessageStopReason = "max_tokens"
	BetaMessageStopReasonStopSequence BetaMessageStopReason = "stop_sequence"
	BetaMessageStopReasonToolUse      BetaMessageStopReason = "tool_use"
)

type BetaMessageDeltaUsage struct {
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
func (r BetaMessageDeltaUsage) RawJSON() string { return r.JSON.raw }
func (r *BetaMessageDeltaUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties Content, Role are required.
type BetaMessageParam struct {
	Content []BetaContentBlockParamUnion `json:"content,omitzero,required"`
	// Any of "user", "assistant".
	Role BetaMessageParamRole `json:"role,omitzero,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaMessageParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaMessageParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaMessageParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaMessageParamRole string

const (
	BetaMessageParamRoleUser      BetaMessageParamRole = "user"
	BetaMessageParamRoleAssistant BetaMessageParamRole = "assistant"
)

type BetaMessageTokensCount struct {
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
func (r BetaMessageTokensCount) RawJSON() string { return r.JSON.raw }
func (r *BetaMessageTokensCount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaMetadataParam struct {
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
func (f BetaMetadataParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaMetadataParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaMetadataParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties Data, MediaType, Type are required.
type BetaPlainTextSourceParam struct {
	Data string `json:"data,required"`
	// This field can be elided, and will marshal its zero value as "text/plain".
	MediaType constant.TextPlain `json:"media_type,required"`
	// This field can be elided, and will marshal its zero value as "text".
	Type constant.Text `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaPlainTextSourceParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaPlainTextSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaPlainTextSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaRawContentBlockDeltaEvent struct {
	Delta BetaRawContentBlockDeltaEventDeltaUnion `json:"delta,required"`
	Index int64                                   `json:"index,required"`
	Type  constant.ContentBlockDelta              `json:"type,required"`
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
func (r BetaRawContentBlockDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaRawContentBlockDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaRawContentBlockDeltaEventDeltaUnion contains all possible properties and
// values from [BetaTextDelta], [BetaInputJSONDelta], [BetaCitationsDelta],
// [BetaThinkingDelta], [BetaSignatureDelta].
//
// Use the [BetaRawContentBlockDeltaEventDeltaUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaRawContentBlockDeltaEventDeltaUnion struct {
	// This field is from variant [BetaTextDelta].
	Text string `json:"text"`
	// Any of "text_delta", "input_json_delta", "citations_delta", "thinking_delta",
	// "signature_delta".
	Type string `json:"type"`
	// This field is from variant [BetaInputJSONDelta].
	PartialJSON string `json:"partial_json"`
	// This field is from variant [BetaCitationsDelta].
	Citation BetaCitationsDeltaCitationUnion `json:"citation"`
	// This field is from variant [BetaThinkingDelta].
	Thinking string `json:"thinking"`
	// This field is from variant [BetaSignatureDelta].
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
//	switch variant := BetaRawContentBlockDeltaEventDeltaUnion.AsAny().(type) {
//	case BetaTextDelta:
//	case BetaInputJSONDelta:
//	case BetaCitationsDelta:
//	case BetaThinkingDelta:
//	case BetaSignatureDelta:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaRawContentBlockDeltaEventDeltaUnion) AsAny() any {
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

func (u BetaRawContentBlockDeltaEventDeltaUnion) AsTextContentBlockDelta() (v BetaTextDelta) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRawContentBlockDeltaEventDeltaUnion) AsInputJSONContentBlockDelta() (v BetaInputJSONDelta) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRawContentBlockDeltaEventDeltaUnion) AsCitationsDelta() (v BetaCitationsDelta) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRawContentBlockDeltaEventDeltaUnion) AsThinkingContentBlockDelta() (v BetaThinkingDelta) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRawContentBlockDeltaEventDeltaUnion) AsSignatureContentBlockDelta() (v BetaSignatureDelta) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaRawContentBlockDeltaEventDeltaUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaRawContentBlockDeltaEventDeltaUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaRawContentBlockStartEvent struct {
	ContentBlock BetaRawContentBlockStartEventContentBlockUnion `json:"content_block,required"`
	Index        int64                                          `json:"index,required"`
	Type         constant.ContentBlockStart                     `json:"type,required"`
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
func (r BetaRawContentBlockStartEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaRawContentBlockStartEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaRawContentBlockStartEventContentBlockUnion contains all possible properties
// and values from [BetaTextBlock], [BetaToolUseBlock], [BetaThinkingBlock],
// [BetaRedactedThinkingBlock].
//
// Use the [BetaRawContentBlockStartEventContentBlockUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaRawContentBlockStartEventContentBlockUnion struct {
	// This field is from variant [BetaTextBlock].
	Citations []BetaTextCitationUnion `json:"citations"`
	// This field is from variant [BetaTextBlock].
	Text string `json:"text"`
	// Any of "text", "tool_use", "thinking", "redacted_thinking".
	Type string `json:"type"`
	// This field is from variant [BetaToolUseBlock].
	ID string `json:"id"`
	// This field is from variant [BetaToolUseBlock].
	Input interface{} `json:"input"`
	// This field is from variant [BetaToolUseBlock].
	Name string `json:"name"`
	// This field is from variant [BetaThinkingBlock].
	Signature string `json:"signature"`
	// This field is from variant [BetaThinkingBlock].
	Thinking string `json:"thinking"`
	// This field is from variant [BetaRedactedThinkingBlock].
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
//	switch variant := BetaRawContentBlockStartEventContentBlockUnion.AsAny().(type) {
//	case BetaTextBlock:
//	case BetaToolUseBlock:
//	case BetaThinkingBlock:
//	case BetaRedactedThinkingBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaRawContentBlockStartEventContentBlockUnion) AsAny() any {
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

func (u BetaRawContentBlockStartEventContentBlockUnion) AsResponseTextBlock() (v BetaTextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRawContentBlockStartEventContentBlockUnion) AsResponseToolUseBlock() (v BetaToolUseBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRawContentBlockStartEventContentBlockUnion) AsResponseThinkingBlock() (v BetaThinkingBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRawContentBlockStartEventContentBlockUnion) AsResponseRedactedThinkingBlock() (v BetaRedactedThinkingBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaRawContentBlockStartEventContentBlockUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaRawContentBlockStartEventContentBlockUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaRawContentBlockStopEvent struct {
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
func (r BetaRawContentBlockStopEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaRawContentBlockStopEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaRawMessageDeltaEvent struct {
	Delta BetaRawMessageDeltaEventDelta `json:"delta,required"`
	Type  constant.MessageDelta         `json:"type,required"`
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
	Usage BetaMessageDeltaUsage `json:"usage,required"`
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
func (r BetaRawMessageDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaRawMessageDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaRawMessageDeltaEventDelta struct {
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
func (r BetaRawMessageDeltaEventDelta) RawJSON() string { return r.JSON.raw }
func (r *BetaRawMessageDeltaEventDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaRawMessageStartEvent struct {
	Message BetaMessage           `json:"message,required"`
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
func (r BetaRawMessageStartEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaRawMessageStartEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaRawMessageStopEvent struct {
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
func (r BetaRawMessageStopEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaRawMessageStopEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaRawMessageStreamEventUnion contains all possible properties and values from
// [BetaRawMessageStartEvent], [BetaRawMessageDeltaEvent],
// [BetaRawMessageStopEvent], [BetaRawContentBlockStartEvent],
// [BetaRawContentBlockDeltaEvent], [BetaRawContentBlockStopEvent].
//
// Use the [BetaRawMessageStreamEventUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaRawMessageStreamEventUnion struct {
	// This field is from variant [BetaRawMessageStartEvent].
	Message BetaMessage `json:"message"`
	// Any of "message_start", "message_delta", "message_stop", "content_block_start",
	// "content_block_delta", "content_block_stop".
	Type string `json:"type"`
	// This field is a union of [BetaRawMessageDeltaEventDelta],
	// [BetaRawContentBlockDeltaEventDeltaUnion]
	Delta BetaRawMessageStreamEventUnionDelta `json:"delta"`
	// This field is from variant [BetaRawMessageDeltaEvent].
	Usage BetaMessageDeltaUsage `json:"usage"`
	// This field is from variant [BetaRawContentBlockStartEvent].
	ContentBlock BetaRawContentBlockStartEventContentBlockUnion `json:"content_block"`
	Index        int64                                          `json:"index"`
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
//	switch variant := BetaRawMessageStreamEventUnion.AsAny().(type) {
//	case BetaRawMessageStartEvent:
//	case BetaRawMessageDeltaEvent:
//	case BetaRawMessageStopEvent:
//	case BetaRawContentBlockStartEvent:
//	case BetaRawContentBlockDeltaEvent:
//	case BetaRawContentBlockStopEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaRawMessageStreamEventUnion) AsAny() any {
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

func (u BetaRawMessageStreamEventUnion) AsMessageStartEvent() (v BetaRawMessageStartEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRawMessageStreamEventUnion) AsMessageDeltaEvent() (v BetaRawMessageDeltaEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRawMessageStreamEventUnion) AsMessageStopEvent() (v BetaRawMessageStopEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRawMessageStreamEventUnion) AsContentBlockStartEvent() (v BetaRawContentBlockStartEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRawMessageStreamEventUnion) AsContentBlockDeltaEvent() (v BetaRawContentBlockDeltaEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaRawMessageStreamEventUnion) AsContentBlockStopEvent() (v BetaRawContentBlockStopEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaRawMessageStreamEventUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaRawMessageStreamEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaRawMessageStreamEventUnionDelta is an implicit subunion of
// [BetaRawMessageStreamEventUnion]. BetaRawMessageStreamEventUnionDelta provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaRawMessageStreamEventUnion].
type BetaRawMessageStreamEventUnionDelta struct {
	// This field is from variant [BetaRawMessageDeltaEventDelta].
	StopReason string `json:"stop_reason"`
	// This field is from variant [BetaRawMessageDeltaEventDelta].
	StopSequence string `json:"stop_sequence"`
	// This field is from variant [BetaRawContentBlockDeltaEventDeltaUnion].
	Text string `json:"text"`
	Type string `json:"type"`
	// This field is from variant [BetaRawContentBlockDeltaEventDeltaUnion].
	PartialJSON string `json:"partial_json"`
	// This field is from variant [BetaRawContentBlockDeltaEventDeltaUnion].
	Citation BetaCitationsDeltaCitationUnion `json:"citation"`
	// This field is from variant [BetaRawContentBlockDeltaEventDeltaUnion].
	Thinking string `json:"thinking"`
	// This field is from variant [BetaRawContentBlockDeltaEventDeltaUnion].
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

func (r *BetaRawMessageStreamEventUnionDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaRedactedThinkingBlock struct {
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
func (r BetaRedactedThinkingBlock) RawJSON() string { return r.JSON.raw }
func (r *BetaRedactedThinkingBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (r BetaRedactedThinkingBlock) ToParam() BetaRedactedThinkingBlockParam {
	var p BetaRedactedThinkingBlockParam
	p.Type = r.Type
	p.Data = r.Data
	return p
}

// The properties Data, Type are required.
type BetaRedactedThinkingBlockParam struct {
	Data string `json:"data,required"`
	// This field can be elided, and will marshal its zero value as
	// "redacted_thinking".
	Type constant.RedactedThinking `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaRedactedThinkingBlockParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaRedactedThinkingBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaRedactedThinkingBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaSignatureDelta struct {
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
func (r BetaSignatureDelta) RawJSON() string { return r.JSON.raw }
func (r *BetaSignatureDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaTextBlock struct {
	// Citations supporting the text block.
	//
	// The type of citation returned will depend on the type of document being cited.
	// Citing a PDF results in `page_location`, plain text results in `char_location`,
	// and content document results in `content_block_location`.
	Citations []BetaTextCitationUnion `json:"citations,required"`
	Text      string                  `json:"text,required"`
	Type      constant.Text           `json:"type,required"`
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
func (r BetaTextBlock) RawJSON() string { return r.JSON.raw }
func (r *BetaTextBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (r BetaTextBlock) ToParam() BetaTextBlockParam {
	var p BetaTextBlockParam
	p.Type = r.Type
	p.Text = r.Text
	p.Citations = make([]BetaTextCitationParamUnion, len(r.Citations))
	for i, citation := range r.Citations {
		switch citationVariant := citation.AsAny().(type) {
		case BetaCitationCharLocation:
			var citationParam BetaCitationCharLocationParam
			citationParam.Type = citationVariant.Type
			citationParam.DocumentTitle = toParam(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
			citationParam.CitedText = citationVariant.CitedText
			citationParam.DocumentIndex = citationVariant.DocumentIndex
			citationParam.EndCharIndex = citationVariant.EndCharIndex
			citationParam.StartCharIndex = citationVariant.StartCharIndex
			p.Citations[i] = BetaTextCitationParamUnion{OfRequestCharLocationCitation: &citationParam}
		case BetaCitationPageLocation:
			var citationParam BetaCitationPageLocationParam
			citationParam.Type = citationVariant.Type
			citationParam.DocumentTitle = toParam(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
			citationParam.DocumentIndex = citationVariant.DocumentIndex
			citationParam.EndPageNumber = citationVariant.EndPageNumber
			citationParam.StartPageNumber = citationVariant.StartPageNumber
			p.Citations[i] = BetaTextCitationParamUnion{OfRequestPageLocationCitation: &citationParam}
		case BetaCitationContentBlockLocation:
			var citationParam BetaCitationContentBlockLocationParam
			citationParam.Type = citationVariant.Type
			citationParam.DocumentTitle = toParam(citationVariant.DocumentTitle, citationVariant.JSON.DocumentTitle)
			citationParam.CitedText = citationVariant.CitedText
			citationParam.DocumentIndex = citationVariant.DocumentIndex
			citationParam.EndBlockIndex = citationVariant.EndBlockIndex
			citationParam.StartBlockIndex = citationVariant.StartBlockIndex
			p.Citations[i] = BetaTextCitationParamUnion{OfRequestContentBlockLocationCitation: &citationParam}
		}
	}
	return p
}

// The properties Text, Type are required.
type BetaTextBlockParam struct {
	Text         string                         `json:"text,required"`
	Citations    []BetaTextCitationParamUnion   `json:"citations,omitzero"`
	CacheControl BetaCacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as "text".
	Type constant.Text `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaTextBlockParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaTextBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaTextBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// BetaTextCitationUnion contains all possible properties and values from
// [BetaCitationCharLocation], [BetaCitationPageLocation],
// [BetaCitationContentBlockLocation].
//
// Use the [BetaTextCitationUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaTextCitationUnion struct {
	CitedText     string `json:"cited_text"`
	DocumentIndex int64  `json:"document_index"`
	DocumentTitle string `json:"document_title"`
	// This field is from variant [BetaCitationCharLocation].
	EndCharIndex int64 `json:"end_char_index"`
	// This field is from variant [BetaCitationCharLocation].
	StartCharIndex int64 `json:"start_char_index"`
	// Any of "char_location", "page_location", "content_block_location".
	Type string `json:"type"`
	// This field is from variant [BetaCitationPageLocation].
	EndPageNumber int64 `json:"end_page_number"`
	// This field is from variant [BetaCitationPageLocation].
	StartPageNumber int64 `json:"start_page_number"`
	// This field is from variant [BetaCitationContentBlockLocation].
	EndBlockIndex int64 `json:"end_block_index"`
	// This field is from variant [BetaCitationContentBlockLocation].
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
//	switch variant := BetaTextCitationUnion.AsAny().(type) {
//	case BetaCitationCharLocation:
//	case BetaCitationPageLocation:
//	case BetaCitationContentBlockLocation:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaTextCitationUnion) AsAny() any {
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

func (u BetaTextCitationUnion) AsResponseCharLocationCitation() (v BetaCitationCharLocation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaTextCitationUnion) AsResponsePageLocationCitation() (v BetaCitationPageLocation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaTextCitationUnion) AsResponseContentBlockLocationCitation() (v BetaCitationContentBlockLocation) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaTextCitationUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaTextCitationUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaTextCitationParamUnion struct {
	OfRequestCharLocationCitation         *BetaCitationCharLocationParam         `json:",omitzero,inline"`
	OfRequestPageLocationCitation         *BetaCitationPageLocationParam         `json:",omitzero,inline"`
	OfRequestContentBlockLocationCitation *BetaCitationContentBlockLocationParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u BetaTextCitationParamUnion) IsPresent() bool { return !param.IsOmitted(u) && !u.IsNull() }
func (u BetaTextCitationParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[BetaTextCitationParamUnion](u.OfRequestCharLocationCitation, u.OfRequestPageLocationCitation, u.OfRequestContentBlockLocationCitation)
}

func (u *BetaTextCitationParamUnion) asAny() any {
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
func (u BetaTextCitationParamUnion) GetEndCharIndex() *int64 {
	if vt := u.OfRequestCharLocationCitation; vt != nil {
		return &vt.EndCharIndex
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaTextCitationParamUnion) GetStartCharIndex() *int64 {
	if vt := u.OfRequestCharLocationCitation; vt != nil {
		return &vt.StartCharIndex
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaTextCitationParamUnion) GetEndPageNumber() *int64 {
	if vt := u.OfRequestPageLocationCitation; vt != nil {
		return &vt.EndPageNumber
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaTextCitationParamUnion) GetStartPageNumber() *int64 {
	if vt := u.OfRequestPageLocationCitation; vt != nil {
		return &vt.StartPageNumber
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaTextCitationParamUnion) GetEndBlockIndex() *int64 {
	if vt := u.OfRequestContentBlockLocationCitation; vt != nil {
		return &vt.EndBlockIndex
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaTextCitationParamUnion) GetStartBlockIndex() *int64 {
	if vt := u.OfRequestContentBlockLocationCitation; vt != nil {
		return &vt.StartBlockIndex
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaTextCitationParamUnion) GetCitedText() *string {
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
func (u BetaTextCitationParamUnion) GetDocumentIndex() *int64 {
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
func (u BetaTextCitationParamUnion) GetDocumentTitle() *string {
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
func (u BetaTextCitationParamUnion) GetType() *string {
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
	apijson.RegisterUnion[BetaTextCitationParamUnion](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaCitationCharLocationParam{}),
			DiscriminatorValue: "char_location",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaCitationPageLocationParam{}),
			DiscriminatorValue: "page_location",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaCitationContentBlockLocationParam{}),
			DiscriminatorValue: "content_block_location",
		},
	)
}

type BetaTextDelta struct {
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
func (r BetaTextDelta) RawJSON() string { return r.JSON.raw }
func (r *BetaTextDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaThinkingBlock struct {
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
func (r BetaThinkingBlock) RawJSON() string { return r.JSON.raw }
func (r *BetaThinkingBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (r BetaThinkingBlock) ToParam() BetaThinkingBlockParam {
	var p BetaThinkingBlockParam
	p.Type = r.Type
	p.Signature = r.Signature
	p.Thinking = r.Thinking
	return p
}

// The properties Signature, Thinking, Type are required.
type BetaThinkingBlockParam struct {
	Signature string `json:"signature,required"`
	Thinking  string `json:"thinking,required"`
	// This field can be elided, and will marshal its zero value as "thinking".
	Type constant.Thinking `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaThinkingBlockParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaThinkingBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaThinkingBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The property Type is required.
type BetaThinkingConfigDisabledParam struct {
	// This field can be elided, and will marshal its zero value as "disabled".
	Type constant.Disabled `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaThinkingConfigDisabledParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaThinkingConfigDisabledParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaThinkingConfigDisabledParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties BudgetTokens, Type are required.
type BetaThinkingConfigEnabledParam struct {
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
func (f BetaThinkingConfigEnabledParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaThinkingConfigEnabledParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaThinkingConfigEnabledParam
	return param.MarshalObject(r, (*shadow)(&r))
}

func BetaThinkingConfigParamOfThinkingConfigEnabled(budgetTokens int64) BetaThinkingConfigParamUnion {
	var variant BetaThinkingConfigEnabledParam
	variant.BudgetTokens = budgetTokens
	return BetaThinkingConfigParamUnion{OfThinkingConfigEnabled: &variant}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaThinkingConfigParamUnion struct {
	OfThinkingConfigEnabled  *BetaThinkingConfigEnabledParam  `json:",omitzero,inline"`
	OfThinkingConfigDisabled *BetaThinkingConfigDisabledParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u BetaThinkingConfigParamUnion) IsPresent() bool { return !param.IsOmitted(u) && !u.IsNull() }
func (u BetaThinkingConfigParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[BetaThinkingConfigParamUnion](u.OfThinkingConfigEnabled, u.OfThinkingConfigDisabled)
}

func (u *BetaThinkingConfigParamUnion) asAny() any {
	if !param.IsOmitted(u.OfThinkingConfigEnabled) {
		return u.OfThinkingConfigEnabled
	} else if !param.IsOmitted(u.OfThinkingConfigDisabled) {
		return u.OfThinkingConfigDisabled
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaThinkingConfigParamUnion) GetBudgetTokens() *int64 {
	if vt := u.OfThinkingConfigEnabled; vt != nil {
		return &vt.BudgetTokens
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaThinkingConfigParamUnion) GetType() *string {
	if vt := u.OfThinkingConfigEnabled; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfThinkingConfigDisabled; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaThinkingConfigParamUnion](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaThinkingConfigEnabledParam{}),
			DiscriminatorValue: "enabled",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaThinkingConfigDisabledParam{}),
			DiscriminatorValue: "disabled",
		},
	)
}

type BetaThinkingDelta struct {
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
func (r BetaThinkingDelta) RawJSON() string { return r.JSON.raw }
func (r *BetaThinkingDelta) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The properties InputSchema, Name are required.
type BetaToolParam struct {
	// [JSON schema](https://json-schema.org/draft/2020-12) for this tool's input.
	//
	// This defines the shape of the `input` that your tool accepts and that the model
	// will produce.
	InputSchema BetaToolInputSchemaParam `json:"input_schema,omitzero,required"`
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
	Description param.Opt[string] `json:"description,omitzero"`
	// Any of "custom".
	Type         BetaToolType                   `json:"type,omitzero"`
	CacheControl BetaCacheControlEphemeralParam `json:"cache_control,omitzero"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaToolParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaToolParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaToolParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// [JSON schema](https://json-schema.org/draft/2020-12) for this tool's input.
//
// This defines the shape of the `input` that your tool accepts and that the model
// will produce.
//
// The property Type is required.
type BetaToolInputSchemaParam struct {
	Properties interface{} `json:"properties,omitzero"`
	// This field can be elided, and will marshal its zero value as "object".
	Type        constant.Object        `json:"type,required"`
	ExtraFields map[string]interface{} `json:"-,extras"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaToolInputSchemaParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaToolInputSchemaParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaToolInputSchemaParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaToolType string

const (
	BetaToolTypeCustom BetaToolType = "custom"
)

// The properties Name, Type are required.
type BetaToolBash20241022Param struct {
	CacheControl BetaCacheControlEphemeralParam `json:"cache_control,omitzero"`
	// Name of the tool.
	//
	// This is how the tool will be called by the model and in tool_use blocks.
	//
	// This field can be elided, and will marshal its zero value as "bash".
	Name constant.Bash `json:"name,required"`
	// This field can be elided, and will marshal its zero value as "bash_20241022".
	Type constant.Bash20241022 `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaToolBash20241022Param) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaToolBash20241022Param) MarshalJSON() (data []byte, err error) {
	type shadow BetaToolBash20241022Param
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties Name, Type are required.
type BetaToolBash20250124Param struct {
	CacheControl BetaCacheControlEphemeralParam `json:"cache_control,omitzero"`
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
func (f BetaToolBash20250124Param) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaToolBash20250124Param) MarshalJSON() (data []byte, err error) {
	type shadow BetaToolBash20250124Param
	return param.MarshalObject(r, (*shadow)(&r))
}

func BetaToolChoiceParamOfToolChoiceTool(name string) BetaToolChoiceUnionParam {
	var variant BetaToolChoiceToolParam
	variant.Name = name
	return BetaToolChoiceUnionParam{OfToolChoiceTool: &variant}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaToolChoiceUnionParam struct {
	OfToolChoiceAuto *BetaToolChoiceAutoParam `json:",omitzero,inline"`
	OfToolChoiceAny  *BetaToolChoiceAnyParam  `json:",omitzero,inline"`
	OfToolChoiceTool *BetaToolChoiceToolParam `json:",omitzero,inline"`
	OfToolChoiceNone *BetaToolChoiceNoneParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u BetaToolChoiceUnionParam) IsPresent() bool { return !param.IsOmitted(u) && !u.IsNull() }
func (u BetaToolChoiceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[BetaToolChoiceUnionParam](u.OfToolChoiceAuto, u.OfToolChoiceAny, u.OfToolChoiceTool, u.OfToolChoiceNone)
}

func (u *BetaToolChoiceUnionParam) asAny() any {
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
func (u BetaToolChoiceUnionParam) GetName() *string {
	if vt := u.OfToolChoiceTool; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaToolChoiceUnionParam) GetType() *string {
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
func (u BetaToolChoiceUnionParam) GetDisableParallelToolUse() *bool {
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
	apijson.RegisterUnion[BetaToolChoiceUnionParam](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaToolChoiceAutoParam{}),
			DiscriminatorValue: "auto",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaToolChoiceAnyParam{}),
			DiscriminatorValue: "any",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaToolChoiceToolParam{}),
			DiscriminatorValue: "tool",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaToolChoiceNoneParam{}),
			DiscriminatorValue: "none",
		},
	)
}

// The model will use any available tools.
//
// The property Type is required.
type BetaToolChoiceAnyParam struct {
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
func (f BetaToolChoiceAnyParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaToolChoiceAnyParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaToolChoiceAnyParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The model will automatically decide whether to use tools.
//
// The property Type is required.
type BetaToolChoiceAutoParam struct {
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
func (f BetaToolChoiceAutoParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaToolChoiceAutoParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaToolChoiceAutoParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The model will not be allowed to use tools.
//
// The property Type is required.
type BetaToolChoiceNoneParam struct {
	// This field can be elided, and will marshal its zero value as "none".
	Type constant.None `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaToolChoiceNoneParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaToolChoiceNoneParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaToolChoiceNoneParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The model will use the specified tool with `tool_choice.name`.
//
// The properties Name, Type are required.
type BetaToolChoiceToolParam struct {
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
func (f BetaToolChoiceToolParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaToolChoiceToolParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaToolChoiceToolParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties DisplayHeightPx, DisplayWidthPx, Name, Type are required.
type BetaToolComputerUse20241022Param struct {
	// The height of the display in pixels.
	DisplayHeightPx int64 `json:"display_height_px,required"`
	// The width of the display in pixels.
	DisplayWidthPx int64 `json:"display_width_px,required"`
	// The X11 display number (e.g. 0, 1) for the display.
	DisplayNumber param.Opt[int64]               `json:"display_number,omitzero"`
	CacheControl  BetaCacheControlEphemeralParam `json:"cache_control,omitzero"`
	// Name of the tool.
	//
	// This is how the tool will be called by the model and in tool_use blocks.
	//
	// This field can be elided, and will marshal its zero value as "computer".
	Name constant.Computer `json:"name,required"`
	// This field can be elided, and will marshal its zero value as
	// "computer_20241022".
	Type constant.Computer20241022 `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaToolComputerUse20241022Param) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaToolComputerUse20241022Param) MarshalJSON() (data []byte, err error) {
	type shadow BetaToolComputerUse20241022Param
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties DisplayHeightPx, DisplayWidthPx, Name, Type are required.
type BetaToolComputerUse20250124Param struct {
	// The height of the display in pixels.
	DisplayHeightPx int64 `json:"display_height_px,required"`
	// The width of the display in pixels.
	DisplayWidthPx int64 `json:"display_width_px,required"`
	// The X11 display number (e.g. 0, 1) for the display.
	DisplayNumber param.Opt[int64]               `json:"display_number,omitzero"`
	CacheControl  BetaCacheControlEphemeralParam `json:"cache_control,omitzero"`
	// Name of the tool.
	//
	// This is how the tool will be called by the model and in tool_use blocks.
	//
	// This field can be elided, and will marshal its zero value as "computer".
	Name constant.Computer `json:"name,required"`
	// This field can be elided, and will marshal its zero value as
	// "computer_20250124".
	Type constant.Computer20250124 `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaToolComputerUse20250124Param) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaToolComputerUse20250124Param) MarshalJSON() (data []byte, err error) {
	type shadow BetaToolComputerUse20250124Param
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties ToolUseID, Type are required.
type BetaToolResultBlockParam struct {
	ToolUseID    string                                 `json:"tool_use_id,required"`
	IsError      param.Opt[bool]                        `json:"is_error,omitzero"`
	CacheControl BetaCacheControlEphemeralParam         `json:"cache_control,omitzero"`
	Content      []BetaToolResultBlockParamContentUnion `json:"content,omitzero"`
	// This field can be elided, and will marshal its zero value as "tool_result".
	Type constant.ToolResult `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaToolResultBlockParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaToolResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaToolResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaToolResultBlockParamContentUnion struct {
	OfRequestTextBlock  *BetaTextBlockParam  `json:",omitzero,inline"`
	OfRequestImageBlock *BetaImageBlockParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u BetaToolResultBlockParamContentUnion) IsPresent() bool {
	return !param.IsOmitted(u) && !u.IsNull()
}
func (u BetaToolResultBlockParamContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[BetaToolResultBlockParamContentUnion](u.OfRequestTextBlock, u.OfRequestImageBlock)
}

func (u *BetaToolResultBlockParamContentUnion) asAny() any {
	if !param.IsOmitted(u.OfRequestTextBlock) {
		return u.OfRequestTextBlock
	} else if !param.IsOmitted(u.OfRequestImageBlock) {
		return u.OfRequestImageBlock
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaToolResultBlockParamContentUnion) GetText() *string {
	if vt := u.OfRequestTextBlock; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaToolResultBlockParamContentUnion) GetCitations() []BetaTextCitationParamUnion {
	if vt := u.OfRequestTextBlock; vt != nil {
		return vt.Citations
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaToolResultBlockParamContentUnion) GetSource() *BetaImageBlockParamSourceUnion {
	if vt := u.OfRequestImageBlock; vt != nil {
		return &vt.Source
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaToolResultBlockParamContentUnion) GetType() *string {
	if vt := u.OfRequestTextBlock; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRequestImageBlock; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's CacheControl property, if present.
func (u BetaToolResultBlockParamContentUnion) GetCacheControl() *BetaCacheControlEphemeralParam {
	if vt := u.OfRequestTextBlock; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfRequestImageBlock; vt != nil {
		return &vt.CacheControl
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaToolResultBlockParamContentUnion](
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaTextBlockParam{}),
			DiscriminatorValue: "text",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaImageBlockParam{}),
			DiscriminatorValue: "image",
		},
	)
}

// The properties Name, Type are required.
type BetaToolTextEditor20241022Param struct {
	CacheControl BetaCacheControlEphemeralParam `json:"cache_control,omitzero"`
	// Name of the tool.
	//
	// This is how the tool will be called by the model and in tool_use blocks.
	//
	// This field can be elided, and will marshal its zero value as
	// "str_replace_editor".
	Name constant.StrReplaceEditor `json:"name,required"`
	// This field can be elided, and will marshal its zero value as
	// "text_editor_20241022".
	Type constant.TextEditor20241022 `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaToolTextEditor20241022Param) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaToolTextEditor20241022Param) MarshalJSON() (data []byte, err error) {
	type shadow BetaToolTextEditor20241022Param
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties Name, Type are required.
type BetaToolTextEditor20250124Param struct {
	CacheControl BetaCacheControlEphemeralParam `json:"cache_control,omitzero"`
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
func (f BetaToolTextEditor20250124Param) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaToolTextEditor20250124Param) MarshalJSON() (data []byte, err error) {
	type shadow BetaToolTextEditor20250124Param
	return param.MarshalObject(r, (*shadow)(&r))
}

func BetaToolUnionParamOfTool(inputSchema BetaToolInputSchemaParam, name string) BetaToolUnionParam {
	var variant BetaToolParam
	variant.InputSchema = inputSchema
	variant.Name = name
	return BetaToolUnionParam{OfTool: &variant}
}

func BetaToolUnionParamOfComputerUseTool20241022(displayHeightPx int64, displayWidthPx int64) BetaToolUnionParam {
	var variant BetaToolComputerUse20241022Param
	variant.DisplayHeightPx = displayHeightPx
	variant.DisplayWidthPx = displayWidthPx
	return BetaToolUnionParam{OfComputerUseTool20241022: &variant}
}

func BetaToolUnionParamOfComputerUseTool20250124(displayHeightPx int64, displayWidthPx int64) BetaToolUnionParam {
	var variant BetaToolComputerUse20250124Param
	variant.DisplayHeightPx = displayHeightPx
	variant.DisplayWidthPx = displayWidthPx
	return BetaToolUnionParam{OfComputerUseTool20250124: &variant}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaToolUnionParam struct {
	OfTool                    *BetaToolParam                    `json:",omitzero,inline"`
	OfComputerUseTool20241022 *BetaToolComputerUse20241022Param `json:",omitzero,inline"`
	OfBashTool20241022        *BetaToolBash20241022Param        `json:",omitzero,inline"`
	OfTextEditor20241022      *BetaToolTextEditor20241022Param  `json:",omitzero,inline"`
	OfComputerUseTool20250124 *BetaToolComputerUse20250124Param `json:",omitzero,inline"`
	OfBashTool20250124        *BetaToolBash20250124Param        `json:",omitzero,inline"`
	OfTextEditor20250124      *BetaToolTextEditor20250124Param  `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u BetaToolUnionParam) IsPresent() bool { return !param.IsOmitted(u) && !u.IsNull() }
func (u BetaToolUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[BetaToolUnionParam](u.OfTool,
		u.OfComputerUseTool20241022,
		u.OfBashTool20241022,
		u.OfTextEditor20241022,
		u.OfComputerUseTool20250124,
		u.OfBashTool20250124,
		u.OfTextEditor20250124)
}

func (u *BetaToolUnionParam) asAny() any {
	if !param.IsOmitted(u.OfTool) {
		return u.OfTool
	} else if !param.IsOmitted(u.OfComputerUseTool20241022) {
		return u.OfComputerUseTool20241022
	} else if !param.IsOmitted(u.OfBashTool20241022) {
		return u.OfBashTool20241022
	} else if !param.IsOmitted(u.OfTextEditor20241022) {
		return u.OfTextEditor20241022
	} else if !param.IsOmitted(u.OfComputerUseTool20250124) {
		return u.OfComputerUseTool20250124
	} else if !param.IsOmitted(u.OfBashTool20250124) {
		return u.OfBashTool20250124
	} else if !param.IsOmitted(u.OfTextEditor20250124) {
		return u.OfTextEditor20250124
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaToolUnionParam) GetInputSchema() *BetaToolInputSchemaParam {
	if vt := u.OfTool; vt != nil {
		return &vt.InputSchema
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaToolUnionParam) GetDescription() *string {
	if vt := u.OfTool; vt != nil && vt.Description.IsPresent() {
		return &vt.Description.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaToolUnionParam) GetName() *string {
	if vt := u.OfTool; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfComputerUseTool20241022; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfBashTool20241022; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfTextEditor20241022; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfComputerUseTool20250124; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfBashTool20250124; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfTextEditor20250124; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaToolUnionParam) GetType() *string {
	if vt := u.OfTool; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfComputerUseTool20241022; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfBashTool20241022; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfTextEditor20241022; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfComputerUseTool20250124; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfBashTool20250124; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfTextEditor20250124; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaToolUnionParam) GetDisplayHeightPx() *int64 {
	if vt := u.OfComputerUseTool20241022; vt != nil {
		return (*int64)(&vt.DisplayHeightPx)
	} else if vt := u.OfComputerUseTool20250124; vt != nil {
		return (*int64)(&vt.DisplayHeightPx)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaToolUnionParam) GetDisplayWidthPx() *int64 {
	if vt := u.OfComputerUseTool20241022; vt != nil {
		return (*int64)(&vt.DisplayWidthPx)
	} else if vt := u.OfComputerUseTool20250124; vt != nil {
		return (*int64)(&vt.DisplayWidthPx)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaToolUnionParam) GetDisplayNumber() *int64 {
	if vt := u.OfComputerUseTool20241022; vt != nil && vt.DisplayNumber.IsPresent() {
		return &vt.DisplayNumber.Value
	} else if vt := u.OfComputerUseTool20250124; vt != nil && vt.DisplayNumber.IsPresent() {
		return &vt.DisplayNumber.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's CacheControl property, if present.
func (u BetaToolUnionParam) GetCacheControl() *BetaCacheControlEphemeralParam {
	if vt := u.OfTool; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfComputerUseTool20241022; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfBashTool20241022; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfTextEditor20241022; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfComputerUseTool20250124; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfBashTool20250124; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfTextEditor20250124; vt != nil {
		return &vt.CacheControl
	}
	return nil
}

type BetaToolUseBlock struct {
	ID    string           `json:"id,required"`
	Input interface{}      `json:"input,required"`
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
func (r BetaToolUseBlock) RawJSON() string { return r.JSON.raw }
func (r *BetaToolUseBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func (r BetaToolUseBlock) ToParam() BetaToolUseBlockParam {
	var p BetaToolUseBlockParam
	p.Type = r.Type
	p.ID = r.ID
	p.Input = r.Input
	p.Name = r.Name
	return p
}

// The properties ID, Input, Name, Type are required.
type BetaToolUseBlockParam struct {
	ID           string                         `json:"id,required"`
	Input        interface{}                    `json:"input,omitzero,required"`
	Name         string                         `json:"name,required"`
	CacheControl BetaCacheControlEphemeralParam `json:"cache_control,omitzero"`
	// This field can be elided, and will marshal its zero value as "tool_use".
	Type constant.ToolUse `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaToolUseBlockParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaToolUseBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaToolUseBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties Type, URL are required.
type BetaURLImageSourceParam struct {
	URL string `json:"url,required"`
	// This field can be elided, and will marshal its zero value as "url".
	Type constant.URL `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaURLImageSourceParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaURLImageSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaURLImageSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// The properties Type, URL are required.
type BetaURLPDFSourceParam struct {
	URL string `json:"url,required"`
	// This field can be elided, and will marshal its zero value as "url".
	Type constant.URL `json:"type,required"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaURLPDFSourceParam) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }
func (r BetaURLPDFSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaURLPDFSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaUsage struct {
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
func (r BetaUsage) RawJSON() string { return r.JSON.raw }
func (r *BetaUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaMessageNewParams struct {
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
	Messages []BetaMessageParam `json:"messages,omitzero,required"`
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
	Metadata BetaMetadataParam `json:"metadata,omitzero"`
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
	System []BetaTextBlockParam `json:"system,omitzero"`
	// Configuration for enabling Claude's extended thinking.
	//
	// When enabled, responses include `thinking` content blocks showing Claude's
	// thinking process before the final answer. Requires a minimum budget of 1,024
	// tokens and counts towards your `max_tokens` limit.
	//
	// See
	// [extended thinking](https://docs.anthropic.com/en/docs/build-with-claude/extended-thinking)
	// for details.
	Thinking BetaThinkingConfigParamUnion `json:"thinking,omitzero"`
	// How the model should use the provided tools. The model can use a specific tool,
	// any available tool, decide by itself, or not use tools at all.
	ToolChoice BetaToolChoiceUnionParam `json:"tool_choice,omitzero"`
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
	Tools []BetaToolUnionParam `json:"tools,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaMessageNewParams) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }

func (r BetaMessageNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaMessageNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}

type BetaMessageCountTokensParams struct {
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
	Messages []BetaMessageParam `json:"messages,omitzero,required"`
	// The model that will complete your prompt.\n\nSee
	// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	Model Model `json:"model,omitzero,required"`
	// System prompt.
	//
	// A system prompt is a way of providing context and instructions to Claude, such
	// as specifying a particular goal or role. See our
	// [guide to system prompts](https://docs.anthropic.com/en/docs/system-prompts).
	System BetaMessageCountTokensParamsSystemUnion `json:"system,omitzero"`
	// Configuration for enabling Claude's extended thinking.
	//
	// When enabled, responses include `thinking` content blocks showing Claude's
	// thinking process before the final answer. Requires a minimum budget of 1,024
	// tokens and counts towards your `max_tokens` limit.
	//
	// See
	// [extended thinking](https://docs.anthropic.com/en/docs/build-with-claude/extended-thinking)
	// for details.
	Thinking BetaThinkingConfigParamUnion `json:"thinking,omitzero"`
	// How the model should use the provided tools. The model can use a specific tool,
	// any available tool, decide by itself, or not use tools at all.
	ToolChoice BetaToolChoiceUnionParam `json:"tool_choice,omitzero"`
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
	Tools []BetaMessageCountTokensParamsToolUnion `json:"tools,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (f BetaMessageCountTokensParams) IsPresent() bool { return !param.IsOmitted(f) && !f.IsNull() }

func (r BetaMessageCountTokensParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaMessageCountTokensParams
	return param.MarshalObject(r, (*shadow)(&r))
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaMessageCountTokensParamsSystemUnion struct {
	OfString                             param.Opt[string]    `json:",omitzero,inline"`
	OfBetaMessageCountTokenssSystemArray []BetaTextBlockParam `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u BetaMessageCountTokensParamsSystemUnion) IsPresent() bool {
	return !param.IsOmitted(u) && !u.IsNull()
}
func (u BetaMessageCountTokensParamsSystemUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[BetaMessageCountTokensParamsSystemUnion](u.OfString, u.OfBetaMessageCountTokenssSystemArray)
}

func (u *BetaMessageCountTokensParamsSystemUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfBetaMessageCountTokenssSystemArray) {
		return &u.OfBetaMessageCountTokenssSystemArray
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaMessageCountTokensParamsToolUnion struct {
	OfTool                    *BetaToolParam                    `json:",omitzero,inline"`
	OfComputerUseTool20241022 *BetaToolComputerUse20241022Param `json:",omitzero,inline"`
	OfBashTool20241022        *BetaToolBash20241022Param        `json:",omitzero,inline"`
	OfTextEditor20241022      *BetaToolTextEditor20241022Param  `json:",omitzero,inline"`
	OfComputerUseTool20250124 *BetaToolComputerUse20250124Param `json:",omitzero,inline"`
	OfBashTool20250124        *BetaToolBash20250124Param        `json:",omitzero,inline"`
	OfTextEditor20250124      *BetaToolTextEditor20250124Param  `json:",omitzero,inline"`
	paramUnion
}

// IsPresent returns true if the field's value is not omitted and not the JSON
// "null". To check if this field is omitted, use [param.IsOmitted].
func (u BetaMessageCountTokensParamsToolUnion) IsPresent() bool {
	return !param.IsOmitted(u) && !u.IsNull()
}
func (u BetaMessageCountTokensParamsToolUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion[BetaMessageCountTokensParamsToolUnion](u.OfTool,
		u.OfComputerUseTool20241022,
		u.OfBashTool20241022,
		u.OfTextEditor20241022,
		u.OfComputerUseTool20250124,
		u.OfBashTool20250124,
		u.OfTextEditor20250124)
}

func (u *BetaMessageCountTokensParamsToolUnion) asAny() any {
	if !param.IsOmitted(u.OfTool) {
		return u.OfTool
	} else if !param.IsOmitted(u.OfComputerUseTool20241022) {
		return u.OfComputerUseTool20241022
	} else if !param.IsOmitted(u.OfBashTool20241022) {
		return u.OfBashTool20241022
	} else if !param.IsOmitted(u.OfTextEditor20241022) {
		return u.OfTextEditor20241022
	} else if !param.IsOmitted(u.OfComputerUseTool20250124) {
		return u.OfComputerUseTool20250124
	} else if !param.IsOmitted(u.OfBashTool20250124) {
		return u.OfBashTool20250124
	} else if !param.IsOmitted(u.OfTextEditor20250124) {
		return u.OfTextEditor20250124
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaMessageCountTokensParamsToolUnion) GetInputSchema() *BetaToolInputSchemaParam {
	if vt := u.OfTool; vt != nil {
		return &vt.InputSchema
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaMessageCountTokensParamsToolUnion) GetDescription() *string {
	if vt := u.OfTool; vt != nil && vt.Description.IsPresent() {
		return &vt.Description.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaMessageCountTokensParamsToolUnion) GetName() *string {
	if vt := u.OfTool; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfComputerUseTool20241022; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfBashTool20241022; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfTextEditor20241022; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfComputerUseTool20250124; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfBashTool20250124; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfTextEditor20250124; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaMessageCountTokensParamsToolUnion) GetType() *string {
	if vt := u.OfTool; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfComputerUseTool20241022; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfBashTool20241022; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfTextEditor20241022; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfComputerUseTool20250124; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfBashTool20250124; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfTextEditor20250124; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaMessageCountTokensParamsToolUnion) GetDisplayHeightPx() *int64 {
	if vt := u.OfComputerUseTool20241022; vt != nil {
		return (*int64)(&vt.DisplayHeightPx)
	} else if vt := u.OfComputerUseTool20250124; vt != nil {
		return (*int64)(&vt.DisplayHeightPx)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaMessageCountTokensParamsToolUnion) GetDisplayWidthPx() *int64 {
	if vt := u.OfComputerUseTool20241022; vt != nil {
		return (*int64)(&vt.DisplayWidthPx)
	} else if vt := u.OfComputerUseTool20250124; vt != nil {
		return (*int64)(&vt.DisplayWidthPx)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaMessageCountTokensParamsToolUnion) GetDisplayNumber() *int64 {
	if vt := u.OfComputerUseTool20241022; vt != nil && vt.DisplayNumber.IsPresent() {
		return &vt.DisplayNumber.Value
	} else if vt := u.OfComputerUseTool20250124; vt != nil && vt.DisplayNumber.IsPresent() {
		return &vt.DisplayNumber.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's CacheControl property, if present.
func (u BetaMessageCountTokensParamsToolUnion) GetCacheControl() *BetaCacheControlEphemeralParam {
	if vt := u.OfTool; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfComputerUseTool20241022; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfBashTool20241022; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfTextEditor20241022; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfComputerUseTool20250124; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfBashTool20250124; vt != nil {
		return &vt.CacheControl
	} else if vt := u.OfTextEditor20250124; vt != nil {
		return &vt.CacheControl
	}
	return nil
}
