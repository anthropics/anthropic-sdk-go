// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"context"
	"fmt"
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
// Learn more about the Messages API in our [user guide](/en/docs/initial-setup)
//
// Note: If you choose to set a timeout for this request, we recommend 10 minutes.
func (r *BetaMessageService) New(ctx context.Context, params BetaMessageNewParams, opts ...option.RequestOption) (res *BetaMessage, err error) {
	for _, v := range params.Betas.Value {
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
func (r *BetaMessageService) NewStreaming(ctx context.Context, params BetaMessageNewParams, opts ...option.RequestOption) (stream *ssestream.Stream[BetaRawMessageStreamEvent]) {
	var (
		raw *http.Response
		err error
	)
	for _, v := range params.Betas.Value {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%s", v)))
	}
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithJSONSet("stream", true)}, opts...)
	path := "v1/messages?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &raw, opts...)
	return ssestream.NewStream[BetaRawMessageStreamEvent](ssestream.NewDecoder(raw), err)
}

// Count the number of tokens in a Message.
//
// The Token Count API can be used to count the number of tokens in a Message,
// including tools, images, and documents, without creating it.
//
// Learn more about token counting in our
// [user guide](/en/docs/build-with-claude/token-counting)
func (r *BetaMessageService) CountTokens(ctx context.Context, params BetaMessageCountTokensParams, opts ...option.RequestOption) (res *BetaMessageTokensCount, err error) {
	for _, v := range params.Betas.Value {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%s", v)))
	}
	opts = append(r.Options[:], opts...)
	path := "v1/messages/count_tokens?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return
}

type BetaBase64ImageSourceParam struct {
	Data      param.Field[string]                         `json:"data,required" format:"byte"`
	MediaType param.Field[BetaBase64ImageSourceMediaType] `json:"media_type,required"`
	Type      param.Field[BetaBase64ImageSourceType]      `json:"type,required"`
}

func (r BetaBase64ImageSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaBase64ImageSourceParam) implementsBetaImageBlockParamSourceUnion() {}

type BetaBase64ImageSourceMediaType string

const (
	BetaBase64ImageSourceMediaTypeImageJPEG BetaBase64ImageSourceMediaType = "image/jpeg"
	BetaBase64ImageSourceMediaTypeImagePNG  BetaBase64ImageSourceMediaType = "image/png"
	BetaBase64ImageSourceMediaTypeImageGIF  BetaBase64ImageSourceMediaType = "image/gif"
	BetaBase64ImageSourceMediaTypeImageWebP BetaBase64ImageSourceMediaType = "image/webp"
)

func (r BetaBase64ImageSourceMediaType) IsKnown() bool {
	switch r {
	case BetaBase64ImageSourceMediaTypeImageJPEG, BetaBase64ImageSourceMediaTypeImagePNG, BetaBase64ImageSourceMediaTypeImageGIF, BetaBase64ImageSourceMediaTypeImageWebP:
		return true
	}
	return false
}

type BetaBase64ImageSourceType string

const (
	BetaBase64ImageSourceTypeBase64 BetaBase64ImageSourceType = "base64"
)

func (r BetaBase64ImageSourceType) IsKnown() bool {
	switch r {
	case BetaBase64ImageSourceTypeBase64:
		return true
	}
	return false
}

type BetaBase64PDFBlockParam struct {
	Source       param.Field[BetaBase64PDFBlockSourceUnionParam] `json:"source,required"`
	Type         param.Field[BetaBase64PDFBlockType]             `json:"type,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam]     `json:"cache_control"`
	Citations    param.Field[BetaCitationsConfigParam]           `json:"citations"`
	Context      param.Field[string]                             `json:"context"`
	Title        param.Field[string]                             `json:"title"`
}

func (r BetaBase64PDFBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaBase64PDFBlockParam) implementsBetaContentBlockParamUnion() {}

type BetaBase64PDFBlockSourceParam struct {
	Type      param.Field[BetaBase64PDFBlockSourceType]      `json:"type,required"`
	Content   param.Field[interface{}]                       `json:"content"`
	Data      param.Field[string]                            `json:"data" format:"byte"`
	MediaType param.Field[BetaBase64PDFBlockSourceMediaType] `json:"media_type"`
	URL       param.Field[string]                            `json:"url"`
}

func (r BetaBase64PDFBlockSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaBase64PDFBlockSourceParam) implementsBetaBase64PDFBlockSourceUnionParam() {}

// Satisfied by [BetaBase64PDFSourceParam], [BetaPlainTextSourceParam],
// [BetaContentBlockSourceParam], [BetaURLPDFSourceParam],
// [BetaBase64PDFBlockSourceParam].
type BetaBase64PDFBlockSourceUnionParam interface {
	implementsBetaBase64PDFBlockSourceUnionParam()
}

type BetaBase64PDFBlockSourceType string

const (
	BetaBase64PDFBlockSourceTypeBase64  BetaBase64PDFBlockSourceType = "base64"
	BetaBase64PDFBlockSourceTypeText    BetaBase64PDFBlockSourceType = "text"
	BetaBase64PDFBlockSourceTypeContent BetaBase64PDFBlockSourceType = "content"
	BetaBase64PDFBlockSourceTypeURL     BetaBase64PDFBlockSourceType = "url"
)

func (r BetaBase64PDFBlockSourceType) IsKnown() bool {
	switch r {
	case BetaBase64PDFBlockSourceTypeBase64, BetaBase64PDFBlockSourceTypeText, BetaBase64PDFBlockSourceTypeContent, BetaBase64PDFBlockSourceTypeURL:
		return true
	}
	return false
}

type BetaBase64PDFBlockSourceMediaType string

const (
	BetaBase64PDFBlockSourceMediaTypeApplicationPDF BetaBase64PDFBlockSourceMediaType = "application/pdf"
	BetaBase64PDFBlockSourceMediaTypeTextPlain      BetaBase64PDFBlockSourceMediaType = "text/plain"
)

func (r BetaBase64PDFBlockSourceMediaType) IsKnown() bool {
	switch r {
	case BetaBase64PDFBlockSourceMediaTypeApplicationPDF, BetaBase64PDFBlockSourceMediaTypeTextPlain:
		return true
	}
	return false
}

type BetaBase64PDFBlockType string

const (
	BetaBase64PDFBlockTypeDocument BetaBase64PDFBlockType = "document"
)

func (r BetaBase64PDFBlockType) IsKnown() bool {
	switch r {
	case BetaBase64PDFBlockTypeDocument:
		return true
	}
	return false
}

type BetaBase64PDFSourceParam struct {
	Data      param.Field[string]                       `json:"data,required" format:"byte"`
	MediaType param.Field[BetaBase64PDFSourceMediaType] `json:"media_type,required"`
	Type      param.Field[BetaBase64PDFSourceType]      `json:"type,required"`
}

func (r BetaBase64PDFSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaBase64PDFSourceParam) implementsBetaBase64PDFBlockSourceUnionParam() {}

type BetaBase64PDFSourceMediaType string

const (
	BetaBase64PDFSourceMediaTypeApplicationPDF BetaBase64PDFSourceMediaType = "application/pdf"
)

func (r BetaBase64PDFSourceMediaType) IsKnown() bool {
	switch r {
	case BetaBase64PDFSourceMediaTypeApplicationPDF:
		return true
	}
	return false
}

type BetaBase64PDFSourceType string

const (
	BetaBase64PDFSourceTypeBase64 BetaBase64PDFSourceType = "base64"
)

func (r BetaBase64PDFSourceType) IsKnown() bool {
	switch r {
	case BetaBase64PDFSourceTypeBase64:
		return true
	}
	return false
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

type BetaCitationCharLocation struct {
	CitedText      string                       `json:"cited_text,required"`
	DocumentIndex  int64                        `json:"document_index,required"`
	DocumentTitle  string                       `json:"document_title,required,nullable"`
	EndCharIndex   int64                        `json:"end_char_index,required"`
	StartCharIndex int64                        `json:"start_char_index,required"`
	Type           BetaCitationCharLocationType `json:"type,required"`
	JSON           betaCitationCharLocationJSON `json:"-"`
}

// betaCitationCharLocationJSON contains the JSON metadata for the struct
// [BetaCitationCharLocation]
type betaCitationCharLocationJSON struct {
	CitedText      apijson.Field
	DocumentIndex  apijson.Field
	DocumentTitle  apijson.Field
	EndCharIndex   apijson.Field
	StartCharIndex apijson.Field
	Type           apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *BetaCitationCharLocation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaCitationCharLocationJSON) RawJSON() string {
	return r.raw
}

func (r BetaCitationCharLocation) implementsBetaCitationsDeltaCitation() {}

func (r BetaCitationCharLocation) implementsBetaTextCitation() {}

type BetaCitationCharLocationType string

const (
	BetaCitationCharLocationTypeCharLocation BetaCitationCharLocationType = "char_location"
)

func (r BetaCitationCharLocationType) IsKnown() bool {
	switch r {
	case BetaCitationCharLocationTypeCharLocation:
		return true
	}
	return false
}

type BetaCitationCharLocationParam struct {
	CitedText      param.Field[string]                            `json:"cited_text,required"`
	DocumentIndex  param.Field[int64]                             `json:"document_index,required"`
	DocumentTitle  param.Field[string]                            `json:"document_title,required"`
	EndCharIndex   param.Field[int64]                             `json:"end_char_index,required"`
	StartCharIndex param.Field[int64]                             `json:"start_char_index,required"`
	Type           param.Field[BetaCitationCharLocationParamType] `json:"type,required"`
}

func (r BetaCitationCharLocationParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaCitationCharLocationParam) implementsBetaTextCitationParamUnion() {}

type BetaCitationCharLocationParamType string

const (
	BetaCitationCharLocationParamTypeCharLocation BetaCitationCharLocationParamType = "char_location"
)

func (r BetaCitationCharLocationParamType) IsKnown() bool {
	switch r {
	case BetaCitationCharLocationParamTypeCharLocation:
		return true
	}
	return false
}

type BetaCitationContentBlockLocation struct {
	CitedText       string                               `json:"cited_text,required"`
	DocumentIndex   int64                                `json:"document_index,required"`
	DocumentTitle   string                               `json:"document_title,required,nullable"`
	EndBlockIndex   int64                                `json:"end_block_index,required"`
	StartBlockIndex int64                                `json:"start_block_index,required"`
	Type            BetaCitationContentBlockLocationType `json:"type,required"`
	JSON            betaCitationContentBlockLocationJSON `json:"-"`
}

// betaCitationContentBlockLocationJSON contains the JSON metadata for the struct
// [BetaCitationContentBlockLocation]
type betaCitationContentBlockLocationJSON struct {
	CitedText       apijson.Field
	DocumentIndex   apijson.Field
	DocumentTitle   apijson.Field
	EndBlockIndex   apijson.Field
	StartBlockIndex apijson.Field
	Type            apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *BetaCitationContentBlockLocation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaCitationContentBlockLocationJSON) RawJSON() string {
	return r.raw
}

func (r BetaCitationContentBlockLocation) implementsBetaCitationsDeltaCitation() {}

func (r BetaCitationContentBlockLocation) implementsBetaTextCitation() {}

type BetaCitationContentBlockLocationType string

const (
	BetaCitationContentBlockLocationTypeContentBlockLocation BetaCitationContentBlockLocationType = "content_block_location"
)

func (r BetaCitationContentBlockLocationType) IsKnown() bool {
	switch r {
	case BetaCitationContentBlockLocationTypeContentBlockLocation:
		return true
	}
	return false
}

type BetaCitationContentBlockLocationParam struct {
	CitedText       param.Field[string]                                    `json:"cited_text,required"`
	DocumentIndex   param.Field[int64]                                     `json:"document_index,required"`
	DocumentTitle   param.Field[string]                                    `json:"document_title,required"`
	EndBlockIndex   param.Field[int64]                                     `json:"end_block_index,required"`
	StartBlockIndex param.Field[int64]                                     `json:"start_block_index,required"`
	Type            param.Field[BetaCitationContentBlockLocationParamType] `json:"type,required"`
}

func (r BetaCitationContentBlockLocationParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaCitationContentBlockLocationParam) implementsBetaTextCitationParamUnion() {}

type BetaCitationContentBlockLocationParamType string

const (
	BetaCitationContentBlockLocationParamTypeContentBlockLocation BetaCitationContentBlockLocationParamType = "content_block_location"
)

func (r BetaCitationContentBlockLocationParamType) IsKnown() bool {
	switch r {
	case BetaCitationContentBlockLocationParamTypeContentBlockLocation:
		return true
	}
	return false
}

type BetaCitationPageLocation struct {
	CitedText       string                       `json:"cited_text,required"`
	DocumentIndex   int64                        `json:"document_index,required"`
	DocumentTitle   string                       `json:"document_title,required,nullable"`
	EndPageNumber   int64                        `json:"end_page_number,required"`
	StartPageNumber int64                        `json:"start_page_number,required"`
	Type            BetaCitationPageLocationType `json:"type,required"`
	JSON            betaCitationPageLocationJSON `json:"-"`
}

// betaCitationPageLocationJSON contains the JSON metadata for the struct
// [BetaCitationPageLocation]
type betaCitationPageLocationJSON struct {
	CitedText       apijson.Field
	DocumentIndex   apijson.Field
	DocumentTitle   apijson.Field
	EndPageNumber   apijson.Field
	StartPageNumber apijson.Field
	Type            apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *BetaCitationPageLocation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaCitationPageLocationJSON) RawJSON() string {
	return r.raw
}

func (r BetaCitationPageLocation) implementsBetaCitationsDeltaCitation() {}

func (r BetaCitationPageLocation) implementsBetaTextCitation() {}

type BetaCitationPageLocationType string

const (
	BetaCitationPageLocationTypePageLocation BetaCitationPageLocationType = "page_location"
)

func (r BetaCitationPageLocationType) IsKnown() bool {
	switch r {
	case BetaCitationPageLocationTypePageLocation:
		return true
	}
	return false
}

type BetaCitationPageLocationParam struct {
	CitedText       param.Field[string]                            `json:"cited_text,required"`
	DocumentIndex   param.Field[int64]                             `json:"document_index,required"`
	DocumentTitle   param.Field[string]                            `json:"document_title,required"`
	EndPageNumber   param.Field[int64]                             `json:"end_page_number,required"`
	StartPageNumber param.Field[int64]                             `json:"start_page_number,required"`
	Type            param.Field[BetaCitationPageLocationParamType] `json:"type,required"`
}

func (r BetaCitationPageLocationParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaCitationPageLocationParam) implementsBetaTextCitationParamUnion() {}

type BetaCitationPageLocationParamType string

const (
	BetaCitationPageLocationParamTypePageLocation BetaCitationPageLocationParamType = "page_location"
)

func (r BetaCitationPageLocationParamType) IsKnown() bool {
	switch r {
	case BetaCitationPageLocationParamTypePageLocation:
		return true
	}
	return false
}

type BetaCitationsConfigParam struct {
	Enabled param.Field[bool] `json:"enabled"`
}

func (r BetaCitationsConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaCitationsDelta struct {
	Citation BetaCitationsDeltaCitation `json:"citation,required"`
	Type     BetaCitationsDeltaType     `json:"type,required"`
	JSON     betaCitationsDeltaJSON     `json:"-"`
}

// betaCitationsDeltaJSON contains the JSON metadata for the struct
// [BetaCitationsDelta]
type betaCitationsDeltaJSON struct {
	Citation    apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaCitationsDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaCitationsDeltaJSON) RawJSON() string {
	return r.raw
}

func (r BetaCitationsDelta) implementsBetaRawContentBlockDeltaEventDelta() {}

type BetaCitationsDeltaCitation struct {
	CitedText       string                         `json:"cited_text,required"`
	DocumentIndex   int64                          `json:"document_index,required"`
	DocumentTitle   string                         `json:"document_title,required,nullable"`
	Type            BetaCitationsDeltaCitationType `json:"type,required"`
	EndBlockIndex   int64                          `json:"end_block_index"`
	EndCharIndex    int64                          `json:"end_char_index"`
	EndPageNumber   int64                          `json:"end_page_number"`
	StartBlockIndex int64                          `json:"start_block_index"`
	StartCharIndex  int64                          `json:"start_char_index"`
	StartPageNumber int64                          `json:"start_page_number"`
	JSON            betaCitationsDeltaCitationJSON `json:"-"`
	union           BetaCitationsDeltaCitationUnion
}

// betaCitationsDeltaCitationJSON contains the JSON metadata for the struct
// [BetaCitationsDeltaCitation]
type betaCitationsDeltaCitationJSON struct {
	CitedText       apijson.Field
	DocumentIndex   apijson.Field
	DocumentTitle   apijson.Field
	Type            apijson.Field
	EndBlockIndex   apijson.Field
	EndCharIndex    apijson.Field
	EndPageNumber   apijson.Field
	StartBlockIndex apijson.Field
	StartCharIndex  apijson.Field
	StartPageNumber apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r betaCitationsDeltaCitationJSON) RawJSON() string {
	return r.raw
}

func (r *BetaCitationsDeltaCitation) UnmarshalJSON(data []byte) (err error) {
	*r = BetaCitationsDeltaCitation{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [BetaCitationsDeltaCitationUnion] interface which you can cast
// to the specific types for more type safety.
//
// Possible runtime types of the union are [BetaCitationCharLocation],
// [BetaCitationPageLocation], [BetaCitationContentBlockLocation].
func (r BetaCitationsDeltaCitation) AsUnion() BetaCitationsDeltaCitationUnion {
	return r.union
}

// Union satisfied by [BetaCitationCharLocation], [BetaCitationPageLocation] or
// [BetaCitationContentBlockLocation].
type BetaCitationsDeltaCitationUnion interface {
	implementsBetaCitationsDeltaCitation()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*BetaCitationsDeltaCitationUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaCitationCharLocation{}),
			DiscriminatorValue: "char_location",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaCitationPageLocation{}),
			DiscriminatorValue: "page_location",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaCitationContentBlockLocation{}),
			DiscriminatorValue: "content_block_location",
		},
	)
}

type BetaCitationsDeltaCitationType string

const (
	BetaCitationsDeltaCitationTypeCharLocation         BetaCitationsDeltaCitationType = "char_location"
	BetaCitationsDeltaCitationTypePageLocation         BetaCitationsDeltaCitationType = "page_location"
	BetaCitationsDeltaCitationTypeContentBlockLocation BetaCitationsDeltaCitationType = "content_block_location"
)

func (r BetaCitationsDeltaCitationType) IsKnown() bool {
	switch r {
	case BetaCitationsDeltaCitationTypeCharLocation, BetaCitationsDeltaCitationTypePageLocation, BetaCitationsDeltaCitationTypeContentBlockLocation:
		return true
	}
	return false
}

type BetaCitationsDeltaType string

const (
	BetaCitationsDeltaTypeCitationsDelta BetaCitationsDeltaType = "citations_delta"
)

func (r BetaCitationsDeltaType) IsKnown() bool {
	switch r {
	case BetaCitationsDeltaTypeCitationsDelta:
		return true
	}
	return false
}

type BetaContentBlock struct {
	Type BetaContentBlockType `json:"type,required"`
	ID   string               `json:"id"`
	// This field can have the runtime type of [[]BetaTextCitation].
	Citations interface{} `json:"citations"`
	Data      string      `json:"data"`
	// This field can have the runtime type of [interface{}].
	Input     interface{}          `json:"input"`
	Name      string               `json:"name"`
	Signature string               `json:"signature"`
	Text      string               `json:"text"`
	Thinking  string               `json:"thinking"`
	JSON      betaContentBlockJSON `json:"-"`
	union     BetaContentBlockUnion
}

// betaContentBlockJSON contains the JSON metadata for the struct
// [BetaContentBlock]
type betaContentBlockJSON struct {
	Type        apijson.Field
	ID          apijson.Field
	Citations   apijson.Field
	Data        apijson.Field
	Input       apijson.Field
	Name        apijson.Field
	Signature   apijson.Field
	Text        apijson.Field
	Thinking    apijson.Field
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
// Possible runtime types of the union are [BetaTextBlock], [BetaToolUseBlock],
// [BetaThinkingBlock], [BetaRedactedThinkingBlock].
func (r BetaContentBlock) AsUnion() BetaContentBlockUnion {
	return r.union
}

// Union satisfied by [BetaTextBlock], [BetaToolUseBlock], [BetaThinkingBlock] or
// [BetaRedactedThinkingBlock].
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
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaThinkingBlock{}),
			DiscriminatorValue: "thinking",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaRedactedThinkingBlock{}),
			DiscriminatorValue: "redacted_thinking",
		},
	)
}

type BetaContentBlockType string

const (
	BetaContentBlockTypeText             BetaContentBlockType = "text"
	BetaContentBlockTypeToolUse          BetaContentBlockType = "tool_use"
	BetaContentBlockTypeThinking         BetaContentBlockType = "thinking"
	BetaContentBlockTypeRedactedThinking BetaContentBlockType = "redacted_thinking"
)

func (r BetaContentBlockType) IsKnown() bool {
	switch r {
	case BetaContentBlockTypeText, BetaContentBlockTypeToolUse, BetaContentBlockTypeThinking, BetaContentBlockTypeRedactedThinking:
		return true
	}
	return false
}

type BetaContentBlockParam struct {
	Type         param.Field[BetaContentBlockParamType]      `json:"type,required"`
	ID           param.Field[string]                         `json:"id"`
	CacheControl param.Field[BetaCacheControlEphemeralParam] `json:"cache_control"`
	Citations    param.Field[interface{}]                    `json:"citations"`
	Content      param.Field[interface{}]                    `json:"content"`
	Context      param.Field[string]                         `json:"context"`
	Data         param.Field[string]                         `json:"data"`
	Input        param.Field[interface{}]                    `json:"input"`
	IsError      param.Field[bool]                           `json:"is_error"`
	Name         param.Field[string]                         `json:"name"`
	Signature    param.Field[string]                         `json:"signature"`
	Source       param.Field[interface{}]                    `json:"source"`
	Text         param.Field[string]                         `json:"text"`
	Thinking     param.Field[string]                         `json:"thinking"`
	Title        param.Field[string]                         `json:"title"`
	ToolUseID    param.Field[string]                         `json:"tool_use_id"`
}

func (r BetaContentBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaContentBlockParam) implementsBetaContentBlockParamUnion() {}

// Satisfied by [BetaTextBlockParam], [BetaImageBlockParam],
// [BetaToolUseBlockParam], [BetaToolResultBlockParam], [BetaBase64PDFBlockParam],
// [BetaThinkingBlockParam], [BetaRedactedThinkingBlockParam],
// [BetaContentBlockParam].
type BetaContentBlockParamUnion interface {
	implementsBetaContentBlockParamUnion()
}

type BetaContentBlockParamType string

const (
	BetaContentBlockParamTypeText             BetaContentBlockParamType = "text"
	BetaContentBlockParamTypeImage            BetaContentBlockParamType = "image"
	BetaContentBlockParamTypeToolUse          BetaContentBlockParamType = "tool_use"
	BetaContentBlockParamTypeToolResult       BetaContentBlockParamType = "tool_result"
	BetaContentBlockParamTypeDocument         BetaContentBlockParamType = "document"
	BetaContentBlockParamTypeThinking         BetaContentBlockParamType = "thinking"
	BetaContentBlockParamTypeRedactedThinking BetaContentBlockParamType = "redacted_thinking"
)

func (r BetaContentBlockParamType) IsKnown() bool {
	switch r {
	case BetaContentBlockParamTypeText, BetaContentBlockParamTypeImage, BetaContentBlockParamTypeToolUse, BetaContentBlockParamTypeToolResult, BetaContentBlockParamTypeDocument, BetaContentBlockParamTypeThinking, BetaContentBlockParamTypeRedactedThinking:
		return true
	}
	return false
}

type BetaContentBlockSourceParam struct {
	Content param.Field[BetaContentBlockSourceContentUnionParam] `json:"content,required"`
	Type    param.Field[BetaContentBlockSourceType]              `json:"type,required"`
}

func (r BetaContentBlockSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaContentBlockSourceParam) implementsBetaBase64PDFBlockSourceUnionParam() {}

// Satisfied by [shared.UnionString],
// [BetaContentBlockSourceContentBetaContentBlockSourceContentParam].
type BetaContentBlockSourceContentUnionParam interface {
	ImplementsBetaContentBlockSourceContentUnionParam()
}

type BetaContentBlockSourceContentBetaContentBlockSourceContentParam []BetaContentBlockSourceContentUnionParam

func (r BetaContentBlockSourceContentBetaContentBlockSourceContentParam) ImplementsBetaContentBlockSourceContentUnionParam() {
}

type BetaContentBlockSourceType string

const (
	BetaContentBlockSourceTypeContent BetaContentBlockSourceType = "content"
)

func (r BetaContentBlockSourceType) IsKnown() bool {
	switch r {
	case BetaContentBlockSourceTypeContent:
		return true
	}
	return false
}

type BetaImageBlockParam struct {
	Source       param.Field[BetaImageBlockParamSourceUnion] `json:"source,required"`
	Type         param.Field[BetaImageBlockParamType]        `json:"type,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam] `json:"cache_control"`
}

func (r BetaImageBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaImageBlockParam) implementsBetaContentBlockParamUnion() {}

func (r BetaImageBlockParam) implementsBetaContentBlockSourceContentUnionParam() {}

func (r BetaImageBlockParam) implementsBetaToolResultBlockParamContentUnion() {}

type BetaImageBlockParamSource struct {
	Type      param.Field[BetaImageBlockParamSourceType]      `json:"type,required"`
	Data      param.Field[string]                             `json:"data" format:"byte"`
	MediaType param.Field[BetaImageBlockParamSourceMediaType] `json:"media_type"`
	URL       param.Field[string]                             `json:"url"`
}

func (r BetaImageBlockParamSource) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaImageBlockParamSource) implementsBetaImageBlockParamSourceUnion() {}

// Satisfied by [BetaBase64ImageSourceParam], [BetaURLImageSourceParam],
// [BetaImageBlockParamSource].
type BetaImageBlockParamSourceUnion interface {
	implementsBetaImageBlockParamSourceUnion()
}

type BetaImageBlockParamSourceType string

const (
	BetaImageBlockParamSourceTypeBase64 BetaImageBlockParamSourceType = "base64"
	BetaImageBlockParamSourceTypeURL    BetaImageBlockParamSourceType = "url"
)

func (r BetaImageBlockParamSourceType) IsKnown() bool {
	switch r {
	case BetaImageBlockParamSourceTypeBase64, BetaImageBlockParamSourceTypeURL:
		return true
	}
	return false
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
	//
	// Total input tokens in a request is the summation of `input_tokens`,
	// `cache_creation_input_tokens`, and `cache_read_input_tokens`.
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

type BetaMessageTokensCount struct {
	// The total number of tokens across the provided list of messages, system prompt,
	// and tools.
	InputTokens int64                      `json:"input_tokens,required"`
	JSON        betaMessageTokensCountJSON `json:"-"`
}

// betaMessageTokensCountJSON contains the JSON metadata for the struct
// [BetaMessageTokensCount]
type betaMessageTokensCountJSON struct {
	InputTokens apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaMessageTokensCount) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaMessageTokensCountJSON) RawJSON() string {
	return r.raw
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

type BetaPlainTextSourceParam struct {
	Data      param.Field[string]                       `json:"data,required"`
	MediaType param.Field[BetaPlainTextSourceMediaType] `json:"media_type,required"`
	Type      param.Field[BetaPlainTextSourceType]      `json:"type,required"`
}

func (r BetaPlainTextSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaPlainTextSourceParam) implementsBetaBase64PDFBlockSourceUnionParam() {}

type BetaPlainTextSourceMediaType string

const (
	BetaPlainTextSourceMediaTypeTextPlain BetaPlainTextSourceMediaType = "text/plain"
)

func (r BetaPlainTextSourceMediaType) IsKnown() bool {
	switch r {
	case BetaPlainTextSourceMediaTypeTextPlain:
		return true
	}
	return false
}

type BetaPlainTextSourceType string

const (
	BetaPlainTextSourceTypeText BetaPlainTextSourceType = "text"
)

func (r BetaPlainTextSourceType) IsKnown() bool {
	switch r {
	case BetaPlainTextSourceTypeText:
		return true
	}
	return false
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
	Type BetaRawContentBlockDeltaEventDeltaType `json:"type,required"`
	// This field can have the runtime type of [BetaCitationsDeltaCitation].
	Citation    interface{}                            `json:"citation"`
	PartialJSON string                                 `json:"partial_json"`
	Signature   string                                 `json:"signature"`
	Text        string                                 `json:"text"`
	Thinking    string                                 `json:"thinking"`
	JSON        betaRawContentBlockDeltaEventDeltaJSON `json:"-"`
	union       BetaRawContentBlockDeltaEventDeltaUnion
}

// betaRawContentBlockDeltaEventDeltaJSON contains the JSON metadata for the struct
// [BetaRawContentBlockDeltaEventDelta]
type betaRawContentBlockDeltaEventDeltaJSON struct {
	Type        apijson.Field
	Citation    apijson.Field
	PartialJSON apijson.Field
	Signature   apijson.Field
	Text        apijson.Field
	Thinking    apijson.Field
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
// Possible runtime types of the union are [BetaTextDelta], [BetaInputJSONDelta],
// [BetaCitationsDelta], [BetaThinkingDelta], [BetaSignatureDelta].
func (r BetaRawContentBlockDeltaEventDelta) AsUnion() BetaRawContentBlockDeltaEventDeltaUnion {
	return r.union
}

// Union satisfied by [BetaTextDelta], [BetaInputJSONDelta], [BetaCitationsDelta],
// [BetaThinkingDelta] or [BetaSignatureDelta].
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
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaCitationsDelta{}),
			DiscriminatorValue: "citations_delta",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaThinkingDelta{}),
			DiscriminatorValue: "thinking_delta",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaSignatureDelta{}),
			DiscriminatorValue: "signature_delta",
		},
	)
}

type BetaRawContentBlockDeltaEventDeltaType string

const (
	BetaRawContentBlockDeltaEventDeltaTypeTextDelta      BetaRawContentBlockDeltaEventDeltaType = "text_delta"
	BetaRawContentBlockDeltaEventDeltaTypeInputJSONDelta BetaRawContentBlockDeltaEventDeltaType = "input_json_delta"
	BetaRawContentBlockDeltaEventDeltaTypeCitationsDelta BetaRawContentBlockDeltaEventDeltaType = "citations_delta"
	BetaRawContentBlockDeltaEventDeltaTypeThinkingDelta  BetaRawContentBlockDeltaEventDeltaType = "thinking_delta"
	BetaRawContentBlockDeltaEventDeltaTypeSignatureDelta BetaRawContentBlockDeltaEventDeltaType = "signature_delta"
)

func (r BetaRawContentBlockDeltaEventDeltaType) IsKnown() bool {
	switch r {
	case BetaRawContentBlockDeltaEventDeltaTypeTextDelta, BetaRawContentBlockDeltaEventDeltaTypeInputJSONDelta, BetaRawContentBlockDeltaEventDeltaTypeCitationsDelta, BetaRawContentBlockDeltaEventDeltaTypeThinkingDelta, BetaRawContentBlockDeltaEventDeltaTypeSignatureDelta:
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
	ID   string                                        `json:"id"`
	// This field can have the runtime type of [[]BetaTextCitation].
	Citations interface{} `json:"citations"`
	Data      string      `json:"data"`
	// This field can have the runtime type of [interface{}].
	Input     interface{}                                   `json:"input"`
	Name      string                                        `json:"name"`
	Signature string                                        `json:"signature"`
	Text      string                                        `json:"text"`
	Thinking  string                                        `json:"thinking"`
	JSON      betaRawContentBlockStartEventContentBlockJSON `json:"-"`
	union     BetaRawContentBlockStartEventContentBlockUnion
}

// betaRawContentBlockStartEventContentBlockJSON contains the JSON metadata for the
// struct [BetaRawContentBlockStartEventContentBlock]
type betaRawContentBlockStartEventContentBlockJSON struct {
	Type        apijson.Field
	ID          apijson.Field
	Citations   apijson.Field
	Data        apijson.Field
	Input       apijson.Field
	Name        apijson.Field
	Signature   apijson.Field
	Text        apijson.Field
	Thinking    apijson.Field
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
// Possible runtime types of the union are [BetaTextBlock], [BetaToolUseBlock],
// [BetaThinkingBlock], [BetaRedactedThinkingBlock].
func (r BetaRawContentBlockStartEventContentBlock) AsUnion() BetaRawContentBlockStartEventContentBlockUnion {
	return r.union
}

// Union satisfied by [BetaTextBlock], [BetaToolUseBlock], [BetaThinkingBlock] or
// [BetaRedactedThinkingBlock].
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
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaThinkingBlock{}),
			DiscriminatorValue: "thinking",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaRedactedThinkingBlock{}),
			DiscriminatorValue: "redacted_thinking",
		},
	)
}

type BetaRawContentBlockStartEventContentBlockType string

const (
	BetaRawContentBlockStartEventContentBlockTypeText             BetaRawContentBlockStartEventContentBlockType = "text"
	BetaRawContentBlockStartEventContentBlockTypeToolUse          BetaRawContentBlockStartEventContentBlockType = "tool_use"
	BetaRawContentBlockStartEventContentBlockTypeThinking         BetaRawContentBlockStartEventContentBlockType = "thinking"
	BetaRawContentBlockStartEventContentBlockTypeRedactedThinking BetaRawContentBlockStartEventContentBlockType = "redacted_thinking"
)

func (r BetaRawContentBlockStartEventContentBlockType) IsKnown() bool {
	switch r {
	case BetaRawContentBlockStartEventContentBlockTypeText, BetaRawContentBlockStartEventContentBlockTypeToolUse, BetaRawContentBlockStartEventContentBlockTypeThinking, BetaRawContentBlockStartEventContentBlockTypeRedactedThinking:
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
	//
	// Total input tokens in a request is the summation of `input_tokens`,
	// `cache_creation_input_tokens`, and `cache_read_input_tokens`.
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
	Type BetaRawMessageStreamEventType `json:"type,required"`
	// This field can have the runtime type of
	// [BetaRawContentBlockStartEventContentBlock].
	ContentBlock interface{} `json:"content_block"`
	// This field can have the runtime type of [BetaRawMessageDeltaEventDelta],
	// [BetaRawContentBlockDeltaEventDelta].
	Delta   interface{} `json:"delta"`
	Index   int64       `json:"index"`
	Message BetaMessage `json:"message"`
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
	Usage BetaMessageDeltaUsage         `json:"usage"`
	JSON  betaRawMessageStreamEventJSON `json:"-"`
	union BetaRawMessageStreamEventUnion
}

// betaRawMessageStreamEventJSON contains the JSON metadata for the struct
// [BetaRawMessageStreamEvent]
type betaRawMessageStreamEventJSON struct {
	Type         apijson.Field
	ContentBlock apijson.Field
	Delta        apijson.Field
	Index        apijson.Field
	Message      apijson.Field
	Usage        apijson.Field
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

type BetaRedactedThinkingBlock struct {
	Data string                        `json:"data,required"`
	Type BetaRedactedThinkingBlockType `json:"type,required"`
	JSON betaRedactedThinkingBlockJSON `json:"-"`
}

// betaRedactedThinkingBlockJSON contains the JSON metadata for the struct
// [BetaRedactedThinkingBlock]
type betaRedactedThinkingBlockJSON struct {
	Data        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaRedactedThinkingBlock) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaRedactedThinkingBlockJSON) RawJSON() string {
	return r.raw
}

func (r BetaRedactedThinkingBlock) implementsBetaContentBlock() {}

func (r BetaRedactedThinkingBlock) implementsBetaRawContentBlockStartEventContentBlock() {}

type BetaRedactedThinkingBlockType string

const (
	BetaRedactedThinkingBlockTypeRedactedThinking BetaRedactedThinkingBlockType = "redacted_thinking"
)

func (r BetaRedactedThinkingBlockType) IsKnown() bool {
	switch r {
	case BetaRedactedThinkingBlockTypeRedactedThinking:
		return true
	}
	return false
}

type BetaRedactedThinkingBlockParam struct {
	Data param.Field[string]                             `json:"data,required"`
	Type param.Field[BetaRedactedThinkingBlockParamType] `json:"type,required"`
}

func (r BetaRedactedThinkingBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaRedactedThinkingBlockParam) implementsBetaContentBlockParamUnion() {}

type BetaRedactedThinkingBlockParamType string

const (
	BetaRedactedThinkingBlockParamTypeRedactedThinking BetaRedactedThinkingBlockParamType = "redacted_thinking"
)

func (r BetaRedactedThinkingBlockParamType) IsKnown() bool {
	switch r {
	case BetaRedactedThinkingBlockParamTypeRedactedThinking:
		return true
	}
	return false
}

type BetaSignatureDelta struct {
	Signature string                 `json:"signature,required"`
	Type      BetaSignatureDeltaType `json:"type,required"`
	JSON      betaSignatureDeltaJSON `json:"-"`
}

// betaSignatureDeltaJSON contains the JSON metadata for the struct
// [BetaSignatureDelta]
type betaSignatureDeltaJSON struct {
	Signature   apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaSignatureDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaSignatureDeltaJSON) RawJSON() string {
	return r.raw
}

func (r BetaSignatureDelta) implementsBetaRawContentBlockDeltaEventDelta() {}

type BetaSignatureDeltaType string

const (
	BetaSignatureDeltaTypeSignatureDelta BetaSignatureDeltaType = "signature_delta"
)

func (r BetaSignatureDeltaType) IsKnown() bool {
	switch r {
	case BetaSignatureDeltaTypeSignatureDelta:
		return true
	}
	return false
}

type BetaTextBlock struct {
	// Citations supporting the text block.
	//
	// The type of citation returned will depend on the type of document being cited.
	// Citing a PDF results in `page_location`, plain text results in `char_location`,
	// and content document results in `content_block_location`.
	Citations []BetaTextCitation `json:"citations,required,nullable"`
	Text      string             `json:"text,required"`
	Type      BetaTextBlockType  `json:"type,required"`
	JSON      betaTextBlockJSON  `json:"-"`
}

// betaTextBlockJSON contains the JSON metadata for the struct [BetaTextBlock]
type betaTextBlockJSON struct {
	Citations   apijson.Field
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
	Citations    param.Field[[]BetaTextCitationParamUnion]   `json:"citations"`
}

func (r BetaTextBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaTextBlockParam) implementsBetaContentBlockParamUnion() {}

func (r BetaTextBlockParam) implementsBetaContentBlockSourceContentUnionParam() {}

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

type BetaTextCitation struct {
	CitedText       string               `json:"cited_text,required"`
	DocumentIndex   int64                `json:"document_index,required"`
	DocumentTitle   string               `json:"document_title,required,nullable"`
	Type            BetaTextCitationType `json:"type,required"`
	EndBlockIndex   int64                `json:"end_block_index"`
	EndCharIndex    int64                `json:"end_char_index"`
	EndPageNumber   int64                `json:"end_page_number"`
	StartBlockIndex int64                `json:"start_block_index"`
	StartCharIndex  int64                `json:"start_char_index"`
	StartPageNumber int64                `json:"start_page_number"`
	JSON            betaTextCitationJSON `json:"-"`
	union           BetaTextCitationUnion
}

// betaTextCitationJSON contains the JSON metadata for the struct
// [BetaTextCitation]
type betaTextCitationJSON struct {
	CitedText       apijson.Field
	DocumentIndex   apijson.Field
	DocumentTitle   apijson.Field
	Type            apijson.Field
	EndBlockIndex   apijson.Field
	EndCharIndex    apijson.Field
	EndPageNumber   apijson.Field
	StartBlockIndex apijson.Field
	StartCharIndex  apijson.Field
	StartPageNumber apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r betaTextCitationJSON) RawJSON() string {
	return r.raw
}

func (r *BetaTextCitation) UnmarshalJSON(data []byte) (err error) {
	*r = BetaTextCitation{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [BetaTextCitationUnion] interface which you can cast to the
// specific types for more type safety.
//
// Possible runtime types of the union are [BetaCitationCharLocation],
// [BetaCitationPageLocation], [BetaCitationContentBlockLocation].
func (r BetaTextCitation) AsUnion() BetaTextCitationUnion {
	return r.union
}

// Union satisfied by [BetaCitationCharLocation], [BetaCitationPageLocation] or
// [BetaCitationContentBlockLocation].
type BetaTextCitationUnion interface {
	implementsBetaTextCitation()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*BetaTextCitationUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaCitationCharLocation{}),
			DiscriminatorValue: "char_location",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaCitationPageLocation{}),
			DiscriminatorValue: "page_location",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaCitationContentBlockLocation{}),
			DiscriminatorValue: "content_block_location",
		},
	)
}

type BetaTextCitationType string

const (
	BetaTextCitationTypeCharLocation         BetaTextCitationType = "char_location"
	BetaTextCitationTypePageLocation         BetaTextCitationType = "page_location"
	BetaTextCitationTypeContentBlockLocation BetaTextCitationType = "content_block_location"
)

func (r BetaTextCitationType) IsKnown() bool {
	switch r {
	case BetaTextCitationTypeCharLocation, BetaTextCitationTypePageLocation, BetaTextCitationTypeContentBlockLocation:
		return true
	}
	return false
}

type BetaTextCitationParam struct {
	CitedText       param.Field[string]                    `json:"cited_text,required"`
	DocumentIndex   param.Field[int64]                     `json:"document_index,required"`
	DocumentTitle   param.Field[string]                    `json:"document_title,required"`
	Type            param.Field[BetaTextCitationParamType] `json:"type,required"`
	EndBlockIndex   param.Field[int64]                     `json:"end_block_index"`
	EndCharIndex    param.Field[int64]                     `json:"end_char_index"`
	EndPageNumber   param.Field[int64]                     `json:"end_page_number"`
	StartBlockIndex param.Field[int64]                     `json:"start_block_index"`
	StartCharIndex  param.Field[int64]                     `json:"start_char_index"`
	StartPageNumber param.Field[int64]                     `json:"start_page_number"`
}

func (r BetaTextCitationParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaTextCitationParam) implementsBetaTextCitationParamUnion() {}

// Satisfied by [BetaCitationCharLocationParam], [BetaCitationPageLocationParam],
// [BetaCitationContentBlockLocationParam], [BetaTextCitationParam].
type BetaTextCitationParamUnion interface {
	implementsBetaTextCitationParamUnion()
}

type BetaTextCitationParamType string

const (
	BetaTextCitationParamTypeCharLocation         BetaTextCitationParamType = "char_location"
	BetaTextCitationParamTypePageLocation         BetaTextCitationParamType = "page_location"
	BetaTextCitationParamTypeContentBlockLocation BetaTextCitationParamType = "content_block_location"
)

func (r BetaTextCitationParamType) IsKnown() bool {
	switch r {
	case BetaTextCitationParamTypeCharLocation, BetaTextCitationParamTypePageLocation, BetaTextCitationParamTypeContentBlockLocation:
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

type BetaThinkingBlock struct {
	Signature string                `json:"signature,required"`
	Thinking  string                `json:"thinking,required"`
	Type      BetaThinkingBlockType `json:"type,required"`
	JSON      betaThinkingBlockJSON `json:"-"`
}

// betaThinkingBlockJSON contains the JSON metadata for the struct
// [BetaThinkingBlock]
type betaThinkingBlockJSON struct {
	Signature   apijson.Field
	Thinking    apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaThinkingBlock) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaThinkingBlockJSON) RawJSON() string {
	return r.raw
}

func (r BetaThinkingBlock) implementsBetaContentBlock() {}

func (r BetaThinkingBlock) implementsBetaRawContentBlockStartEventContentBlock() {}

type BetaThinkingBlockType string

const (
	BetaThinkingBlockTypeThinking BetaThinkingBlockType = "thinking"
)

func (r BetaThinkingBlockType) IsKnown() bool {
	switch r {
	case BetaThinkingBlockTypeThinking:
		return true
	}
	return false
}

type BetaThinkingBlockParam struct {
	Signature param.Field[string]                     `json:"signature,required"`
	Thinking  param.Field[string]                     `json:"thinking,required"`
	Type      param.Field[BetaThinkingBlockParamType] `json:"type,required"`
}

func (r BetaThinkingBlockParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaThinkingBlockParam) implementsBetaContentBlockParamUnion() {}

type BetaThinkingBlockParamType string

const (
	BetaThinkingBlockParamTypeThinking BetaThinkingBlockParamType = "thinking"
)

func (r BetaThinkingBlockParamType) IsKnown() bool {
	switch r {
	case BetaThinkingBlockParamTypeThinking:
		return true
	}
	return false
}

type BetaThinkingConfigDisabledParam struct {
	Type param.Field[BetaThinkingConfigDisabledType] `json:"type,required"`
}

func (r BetaThinkingConfigDisabledParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaThinkingConfigDisabledParam) implementsBetaThinkingConfigParamUnion() {}

type BetaThinkingConfigDisabledType string

const (
	BetaThinkingConfigDisabledTypeDisabled BetaThinkingConfigDisabledType = "disabled"
)

func (r BetaThinkingConfigDisabledType) IsKnown() bool {
	switch r {
	case BetaThinkingConfigDisabledTypeDisabled:
		return true
	}
	return false
}

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
	BudgetTokens param.Field[int64]                         `json:"budget_tokens,required"`
	Type         param.Field[BetaThinkingConfigEnabledType] `json:"type,required"`
}

func (r BetaThinkingConfigEnabledParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaThinkingConfigEnabledParam) implementsBetaThinkingConfigParamUnion() {}

type BetaThinkingConfigEnabledType string

const (
	BetaThinkingConfigEnabledTypeEnabled BetaThinkingConfigEnabledType = "enabled"
)

func (r BetaThinkingConfigEnabledType) IsKnown() bool {
	switch r {
	case BetaThinkingConfigEnabledTypeEnabled:
		return true
	}
	return false
}

// Configuration for enabling Claude's extended thinking.
//
// When enabled, responses include `thinking` content blocks showing Claude's
// thinking process before the final answer. Requires a minimum budget of 1,024
// tokens and counts towards your `max_tokens` limit.
//
// See
// [extended thinking](https://docs.anthropic.com/en/docs/build-with-claude/extended-thinking)
// for details.
type BetaThinkingConfigParam struct {
	Type param.Field[BetaThinkingConfigParamType] `json:"type,required"`
	// Determines how many tokens Claude can use for its internal reasoning process.
	// Larger budgets can enable more thorough analysis for complex problems, improving
	// response quality.
	//
	// Must be ≥1024 and less than `max_tokens`.
	//
	// See
	// [extended thinking](https://docs.anthropic.com/en/docs/build-with-claude/extended-thinking)
	// for details.
	BudgetTokens param.Field[int64] `json:"budget_tokens"`
}

func (r BetaThinkingConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaThinkingConfigParam) implementsBetaThinkingConfigParamUnion() {}

// Configuration for enabling Claude's extended thinking.
//
// When enabled, responses include `thinking` content blocks showing Claude's
// thinking process before the final answer. Requires a minimum budget of 1,024
// tokens and counts towards your `max_tokens` limit.
//
// See
// [extended thinking](https://docs.anthropic.com/en/docs/build-with-claude/extended-thinking)
// for details.
//
// Satisfied by [BetaThinkingConfigEnabledParam],
// [BetaThinkingConfigDisabledParam], [BetaThinkingConfigParam].
type BetaThinkingConfigParamUnion interface {
	implementsBetaThinkingConfigParamUnion()
}

type BetaThinkingConfigParamType string

const (
	BetaThinkingConfigParamTypeEnabled  BetaThinkingConfigParamType = "enabled"
	BetaThinkingConfigParamTypeDisabled BetaThinkingConfigParamType = "disabled"
)

func (r BetaThinkingConfigParamType) IsKnown() bool {
	switch r {
	case BetaThinkingConfigParamTypeEnabled, BetaThinkingConfigParamTypeDisabled:
		return true
	}
	return false
}

type BetaThinkingDelta struct {
	Thinking string                `json:"thinking,required"`
	Type     BetaThinkingDeltaType `json:"type,required"`
	JSON     betaThinkingDeltaJSON `json:"-"`
}

// betaThinkingDeltaJSON contains the JSON metadata for the struct
// [BetaThinkingDelta]
type betaThinkingDeltaJSON struct {
	Thinking    apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaThinkingDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaThinkingDeltaJSON) RawJSON() string {
	return r.raw
}

func (r BetaThinkingDelta) implementsBetaRawContentBlockDeltaEventDelta() {}

type BetaThinkingDeltaType string

const (
	BetaThinkingDeltaTypeThinkingDelta BetaThinkingDeltaType = "thinking_delta"
)

func (r BetaThinkingDeltaType) IsKnown() bool {
	switch r {
	case BetaThinkingDeltaTypeThinkingDelta:
		return true
	}
	return false
}

type BetaToolParam struct {
	// [JSON schema](https://json-schema.org/draft/2020-12) for this tool's input.
	//
	// This defines the shape of the `input` that your tool accepts and that the model
	// will produce.
	InputSchema param.Field[BetaToolInputSchemaParam] `json:"input_schema,required"`
	// Name of the tool.
	//
	// This is how the tool will be called by the model and in tool_use blocks.
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

func (r BetaToolParam) implementsBetaMessageCountTokensParamsToolUnion() {}

// [JSON schema](https://json-schema.org/draft/2020-12) for this tool's input.
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
	// Name of the tool.
	//
	// This is how the tool will be called by the model and in tool_use blocks.
	Name         param.Field[BetaToolBash20241022Name]       `json:"name,required"`
	Type         param.Field[BetaToolBash20241022Type]       `json:"type,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam] `json:"cache_control"`
}

func (r BetaToolBash20241022Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolBash20241022Param) implementsBetaToolUnionUnionParam() {}

func (r BetaToolBash20241022Param) implementsBetaMessageCountTokensParamsToolUnion() {}

// Name of the tool.
//
// This is how the tool will be called by the model and in tool_use blocks.
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

type BetaToolBash20250124Param struct {
	// Name of the tool.
	//
	// This is how the tool will be called by the model and in tool_use blocks.
	Name         param.Field[BetaToolBash20250124Name]       `json:"name,required"`
	Type         param.Field[BetaToolBash20250124Type]       `json:"type,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam] `json:"cache_control"`
}

func (r BetaToolBash20250124Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolBash20250124Param) implementsBetaToolUnionUnionParam() {}

func (r BetaToolBash20250124Param) implementsBetaMessageCountTokensParamsToolUnion() {}

// Name of the tool.
//
// This is how the tool will be called by the model and in tool_use blocks.
type BetaToolBash20250124Name string

const (
	BetaToolBash20250124NameBash BetaToolBash20250124Name = "bash"
)

func (r BetaToolBash20250124Name) IsKnown() bool {
	switch r {
	case BetaToolBash20250124NameBash:
		return true
	}
	return false
}

type BetaToolBash20250124Type string

const (
	BetaToolBash20250124TypeBash20250124 BetaToolBash20250124Type = "bash_20250124"
)

func (r BetaToolBash20250124Type) IsKnown() bool {
	switch r {
	case BetaToolBash20250124TypeBash20250124:
		return true
	}
	return false
}

// How the model should use the provided tools. The model can use a specific tool,
// any available tool, decide by itself, or not use tools at all.
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
// any available tool, decide by itself, or not use tools at all.
//
// Satisfied by [BetaToolChoiceAutoParam], [BetaToolChoiceAnyParam],
// [BetaToolChoiceToolParam], [BetaToolChoiceNoneParam], [BetaToolChoiceParam].
type BetaToolChoiceUnionParam interface {
	implementsBetaToolChoiceUnionParam()
}

type BetaToolChoiceType string

const (
	BetaToolChoiceTypeAuto BetaToolChoiceType = "auto"
	BetaToolChoiceTypeAny  BetaToolChoiceType = "any"
	BetaToolChoiceTypeTool BetaToolChoiceType = "tool"
	BetaToolChoiceTypeNone BetaToolChoiceType = "none"
)

func (r BetaToolChoiceType) IsKnown() bool {
	switch r {
	case BetaToolChoiceTypeAuto, BetaToolChoiceTypeAny, BetaToolChoiceTypeTool, BetaToolChoiceTypeNone:
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

// The model will not be allowed to use tools.
type BetaToolChoiceNoneParam struct {
	Type param.Field[BetaToolChoiceNoneType] `json:"type,required"`
}

func (r BetaToolChoiceNoneParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolChoiceNoneParam) implementsBetaToolChoiceUnionParam() {}

type BetaToolChoiceNoneType string

const (
	BetaToolChoiceNoneTypeNone BetaToolChoiceNoneType = "none"
)

func (r BetaToolChoiceNoneType) IsKnown() bool {
	switch r {
	case BetaToolChoiceNoneTypeNone:
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
	// The height of the display in pixels.
	DisplayHeightPx param.Field[int64] `json:"display_height_px,required"`
	// The width of the display in pixels.
	DisplayWidthPx param.Field[int64] `json:"display_width_px,required"`
	// Name of the tool.
	//
	// This is how the tool will be called by the model and in tool_use blocks.
	Name         param.Field[BetaToolComputerUse20241022Name] `json:"name,required"`
	Type         param.Field[BetaToolComputerUse20241022Type] `json:"type,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam]  `json:"cache_control"`
	// The X11 display number (e.g. 0, 1) for the display.
	DisplayNumber param.Field[int64] `json:"display_number"`
}

func (r BetaToolComputerUse20241022Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolComputerUse20241022Param) implementsBetaToolUnionUnionParam() {}

func (r BetaToolComputerUse20241022Param) implementsBetaMessageCountTokensParamsToolUnion() {}

// Name of the tool.
//
// This is how the tool will be called by the model and in tool_use blocks.
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

type BetaToolComputerUse20250124Param struct {
	// The height of the display in pixels.
	DisplayHeightPx param.Field[int64] `json:"display_height_px,required"`
	// The width of the display in pixels.
	DisplayWidthPx param.Field[int64] `json:"display_width_px,required"`
	// Name of the tool.
	//
	// This is how the tool will be called by the model and in tool_use blocks.
	Name         param.Field[BetaToolComputerUse20250124Name] `json:"name,required"`
	Type         param.Field[BetaToolComputerUse20250124Type] `json:"type,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam]  `json:"cache_control"`
	// The X11 display number (e.g. 0, 1) for the display.
	DisplayNumber param.Field[int64] `json:"display_number"`
}

func (r BetaToolComputerUse20250124Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolComputerUse20250124Param) implementsBetaToolUnionUnionParam() {}

func (r BetaToolComputerUse20250124Param) implementsBetaMessageCountTokensParamsToolUnion() {}

// Name of the tool.
//
// This is how the tool will be called by the model and in tool_use blocks.
type BetaToolComputerUse20250124Name string

const (
	BetaToolComputerUse20250124NameComputer BetaToolComputerUse20250124Name = "computer"
)

func (r BetaToolComputerUse20250124Name) IsKnown() bool {
	switch r {
	case BetaToolComputerUse20250124NameComputer:
		return true
	}
	return false
}

type BetaToolComputerUse20250124Type string

const (
	BetaToolComputerUse20250124TypeComputer20250124 BetaToolComputerUse20250124Type = "computer_20250124"
)

func (r BetaToolComputerUse20250124Type) IsKnown() bool {
	switch r {
	case BetaToolComputerUse20250124TypeComputer20250124:
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
	Type         param.Field[BetaToolResultBlockParamContentType] `json:"type,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam]      `json:"cache_control"`
	Citations    param.Field[interface{}]                         `json:"citations"`
	Source       param.Field[interface{}]                         `json:"source"`
	Text         param.Field[string]                              `json:"text"`
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
	// Name of the tool.
	//
	// This is how the tool will be called by the model and in tool_use blocks.
	Name         param.Field[BetaToolTextEditor20241022Name] `json:"name,required"`
	Type         param.Field[BetaToolTextEditor20241022Type] `json:"type,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam] `json:"cache_control"`
}

func (r BetaToolTextEditor20241022Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolTextEditor20241022Param) implementsBetaToolUnionUnionParam() {}

func (r BetaToolTextEditor20241022Param) implementsBetaMessageCountTokensParamsToolUnion() {}

// Name of the tool.
//
// This is how the tool will be called by the model and in tool_use blocks.
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

type BetaToolTextEditor20250124Param struct {
	// Name of the tool.
	//
	// This is how the tool will be called by the model and in tool_use blocks.
	Name         param.Field[BetaToolTextEditor20250124Name] `json:"name,required"`
	Type         param.Field[BetaToolTextEditor20250124Type] `json:"type,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam] `json:"cache_control"`
}

func (r BetaToolTextEditor20250124Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolTextEditor20250124Param) implementsBetaToolUnionUnionParam() {}

func (r BetaToolTextEditor20250124Param) implementsBetaMessageCountTokensParamsToolUnion() {}

// Name of the tool.
//
// This is how the tool will be called by the model and in tool_use blocks.
type BetaToolTextEditor20250124Name string

const (
	BetaToolTextEditor20250124NameStrReplaceEditor BetaToolTextEditor20250124Name = "str_replace_editor"
)

func (r BetaToolTextEditor20250124Name) IsKnown() bool {
	switch r {
	case BetaToolTextEditor20250124NameStrReplaceEditor:
		return true
	}
	return false
}

type BetaToolTextEditor20250124Type string

const (
	BetaToolTextEditor20250124TypeTextEditor20250124 BetaToolTextEditor20250124Type = "text_editor_20250124"
)

func (r BetaToolTextEditor20250124Type) IsKnown() bool {
	switch r {
	case BetaToolTextEditor20250124TypeTextEditor20250124:
		return true
	}
	return false
}

type BetaToolUnionParam struct {
	// Name of the tool.
	//
	// This is how the tool will be called by the model and in tool_use blocks.
	Name         param.Field[string]                         `json:"name,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam] `json:"cache_control"`
	// Description of what this tool does.
	//
	// Tool descriptions should be as detailed as possible. The more information that
	// the model has about what the tool is and how to use it, the better it will
	// perform. You can use natural language descriptions to reinforce important
	// aspects of the tool input JSON schema.
	Description param.Field[string] `json:"description"`
	// The height of the display in pixels.
	DisplayHeightPx param.Field[int64] `json:"display_height_px"`
	// The X11 display number (e.g. 0, 1) for the display.
	DisplayNumber param.Field[int64] `json:"display_number"`
	// The width of the display in pixels.
	DisplayWidthPx param.Field[int64]             `json:"display_width_px"`
	InputSchema    param.Field[interface{}]       `json:"input_schema"`
	Type           param.Field[BetaToolUnionType] `json:"type"`
}

func (r BetaToolUnionParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaToolUnionParam) implementsBetaToolUnionUnionParam() {}

// Satisfied by [BetaToolParam], [BetaToolComputerUse20241022Param],
// [BetaToolBash20241022Param], [BetaToolTextEditor20241022Param],
// [BetaToolComputerUse20250124Param], [BetaToolBash20250124Param],
// [BetaToolTextEditor20250124Param], [BetaToolUnionParam].
type BetaToolUnionUnionParam interface {
	implementsBetaToolUnionUnionParam()
}

type BetaToolUnionType string

const (
	BetaToolUnionTypeCustom             BetaToolUnionType = "custom"
	BetaToolUnionTypeComputer20241022   BetaToolUnionType = "computer_20241022"
	BetaToolUnionTypeBash20241022       BetaToolUnionType = "bash_20241022"
	BetaToolUnionTypeTextEditor20241022 BetaToolUnionType = "text_editor_20241022"
	BetaToolUnionTypeComputer20250124   BetaToolUnionType = "computer_20250124"
	BetaToolUnionTypeBash20250124       BetaToolUnionType = "bash_20250124"
	BetaToolUnionTypeTextEditor20250124 BetaToolUnionType = "text_editor_20250124"
)

func (r BetaToolUnionType) IsKnown() bool {
	switch r {
	case BetaToolUnionTypeCustom, BetaToolUnionTypeComputer20241022, BetaToolUnionTypeBash20241022, BetaToolUnionTypeTextEditor20241022, BetaToolUnionTypeComputer20250124, BetaToolUnionTypeBash20250124, BetaToolUnionTypeTextEditor20250124:
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

type BetaURLImageSourceParam struct {
	Type param.Field[BetaURLImageSourceType] `json:"type,required"`
	URL  param.Field[string]                 `json:"url,required"`
}

func (r BetaURLImageSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaURLImageSourceParam) implementsBetaImageBlockParamSourceUnion() {}

type BetaURLImageSourceType string

const (
	BetaURLImageSourceTypeURL BetaURLImageSourceType = "url"
)

func (r BetaURLImageSourceType) IsKnown() bool {
	switch r {
	case BetaURLImageSourceTypeURL:
		return true
	}
	return false
}

type BetaURLPDFSourceParam struct {
	Type param.Field[BetaURLPDFSourceType] `json:"type,required"`
	URL  param.Field[string]               `json:"url,required"`
}

func (r BetaURLPDFSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaURLPDFSourceParam) implementsBetaBase64PDFBlockSourceUnionParam() {}

type BetaURLPDFSourceType string

const (
	BetaURLPDFSourceTypeURL BetaURLPDFSourceType = "url"
)

func (r BetaURLPDFSourceType) IsKnown() bool {
	switch r {
	case BetaURLPDFSourceTypeURL:
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
	// Configuration for enabling Claude's extended thinking.
	//
	// When enabled, responses include `thinking` content blocks showing Claude's
	// thinking process before the final answer. Requires a minimum budget of 1,024
	// tokens and counts towards your `max_tokens` limit.
	//
	// See
	// [extended thinking](https://docs.anthropic.com/en/docs/build-with-claude/extended-thinking)
	// for details.
	Thinking param.Field[BetaThinkingConfigParamUnion] `json:"thinking"`
	// How the model should use the provided tools. The model can use a specific tool,
	// any available tool, decide by itself, or not use tools at all.
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
	Messages param.Field[[]BetaMessageParam] `json:"messages,required"`
	// The model that will complete your prompt.\n\nSee
	// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	Model param.Field[Model] `json:"model,required"`
	// System prompt.
	//
	// A system prompt is a way of providing context and instructions to Claude, such
	// as specifying a particular goal or role. See our
	// [guide to system prompts](https://docs.anthropic.com/en/docs/system-prompts).
	System param.Field[BetaMessageCountTokensParamsSystemUnion] `json:"system"`
	// Configuration for enabling Claude's extended thinking.
	//
	// When enabled, responses include `thinking` content blocks showing Claude's
	// thinking process before the final answer. Requires a minimum budget of 1,024
	// tokens and counts towards your `max_tokens` limit.
	//
	// See
	// [extended thinking](https://docs.anthropic.com/en/docs/build-with-claude/extended-thinking)
	// for details.
	Thinking param.Field[BetaThinkingConfigParamUnion] `json:"thinking"`
	// How the model should use the provided tools. The model can use a specific tool,
	// any available tool, decide by itself, or not use tools at all.
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
	Tools param.Field[[]BetaMessageCountTokensParamsToolUnion] `json:"tools"`
	// Optional header to specify the beta version(s) you want to use.
	Betas param.Field[[]AnthropicBeta] `header:"anthropic-beta"`
}

func (r BetaMessageCountTokensParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// System prompt.
//
// A system prompt is a way of providing context and instructions to Claude, such
// as specifying a particular goal or role. See our
// [guide to system prompts](https://docs.anthropic.com/en/docs/system-prompts).
//
// Satisfied by [shared.UnionString], [BetaMessageCountTokensParamsSystemArray].
type BetaMessageCountTokensParamsSystemUnion interface {
	ImplementsBetaMessageCountTokensParamsSystemUnion()
}

type BetaMessageCountTokensParamsSystemArray []BetaTextBlockParam

func (r BetaMessageCountTokensParamsSystemArray) ImplementsBetaMessageCountTokensParamsSystemUnion() {
}

type BetaMessageCountTokensParamsTool struct {
	// Name of the tool.
	//
	// This is how the tool will be called by the model and in tool_use blocks.
	Name         param.Field[string]                         `json:"name,required"`
	CacheControl param.Field[BetaCacheControlEphemeralParam] `json:"cache_control"`
	// Description of what this tool does.
	//
	// Tool descriptions should be as detailed as possible. The more information that
	// the model has about what the tool is and how to use it, the better it will
	// perform. You can use natural language descriptions to reinforce important
	// aspects of the tool input JSON schema.
	Description param.Field[string] `json:"description"`
	// The height of the display in pixels.
	DisplayHeightPx param.Field[int64] `json:"display_height_px"`
	// The X11 display number (e.g. 0, 1) for the display.
	DisplayNumber param.Field[int64] `json:"display_number"`
	// The width of the display in pixels.
	DisplayWidthPx param.Field[int64]                                 `json:"display_width_px"`
	InputSchema    param.Field[interface{}]                           `json:"input_schema"`
	Type           param.Field[BetaMessageCountTokensParamsToolsType] `json:"type"`
}

func (r BetaMessageCountTokensParamsTool) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r BetaMessageCountTokensParamsTool) implementsBetaMessageCountTokensParamsToolUnion() {}

// Satisfied by [BetaToolParam], [BetaToolComputerUse20241022Param],
// [BetaToolBash20241022Param], [BetaToolTextEditor20241022Param],
// [BetaToolComputerUse20250124Param], [BetaToolBash20250124Param],
// [BetaToolTextEditor20250124Param], [BetaMessageCountTokensParamsTool].
type BetaMessageCountTokensParamsToolUnion interface {
	implementsBetaMessageCountTokensParamsToolUnion()
}

type BetaMessageCountTokensParamsToolsType string

const (
	BetaMessageCountTokensParamsToolsTypeCustom             BetaMessageCountTokensParamsToolsType = "custom"
	BetaMessageCountTokensParamsToolsTypeComputer20241022   BetaMessageCountTokensParamsToolsType = "computer_20241022"
	BetaMessageCountTokensParamsToolsTypeBash20241022       BetaMessageCountTokensParamsToolsType = "bash_20241022"
	BetaMessageCountTokensParamsToolsTypeTextEditor20241022 BetaMessageCountTokensParamsToolsType = "text_editor_20241022"
	BetaMessageCountTokensParamsToolsTypeComputer20250124   BetaMessageCountTokensParamsToolsType = "computer_20250124"
	BetaMessageCountTokensParamsToolsTypeBash20250124       BetaMessageCountTokensParamsToolsType = "bash_20250124"
	BetaMessageCountTokensParamsToolsTypeTextEditor20250124 BetaMessageCountTokensParamsToolsType = "text_editor_20250124"
)

func (r BetaMessageCountTokensParamsToolsType) IsKnown() bool {
	switch r {
	case BetaMessageCountTokensParamsToolsTypeCustom, BetaMessageCountTokensParamsToolsTypeComputer20241022, BetaMessageCountTokensParamsToolsTypeBash20241022, BetaMessageCountTokensParamsToolsTypeTextEditor20241022, BetaMessageCountTokensParamsToolsTypeComputer20250124, BetaMessageCountTokensParamsToolsTypeBash20250124, BetaMessageCountTokensParamsToolsTypeTextEditor20250124:
		return true
	}
	return false
}
