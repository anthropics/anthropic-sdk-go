// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// BetaModelService contains methods and other services that help with interacting
// with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaModelService] method instead.
type BetaModelService struct {
	Options []option.RequestOption
}

// NewBetaModelService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaModelService(opts ...option.RequestOption) (r BetaModelService) {
	r = BetaModelService{}
	r.Options = opts
	return
}

// Get a specific model.
//
// The Models API response can be used to determine information about a specific
// model or resolve a model alias to a model ID.
func (r *BetaModelService) Get(ctx context.Context, modelID string, query BetaModelGetParams, opts ...option.RequestOption) (res *BetaModelInfo, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	if modelID == "" {
		err = errors.New("missing required model_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/models/%s?beta=true", modelID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List available models.
//
// The Models API response can be used to determine which models are available for
// use in the API. More recently released models are listed first.
func (r *BetaModelService) List(ctx context.Context, params BetaModelListParams, opts ...option.RequestOption) (res *pagination.Page[BetaModelInfo], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/models?beta=true"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, params, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// List available models.
//
// The Models API response can be used to determine which models are available for
// use in the API. More recently released models are listed first.
func (r *BetaModelService) ListAutoPaging(ctx context.Context, params BetaModelListParams, opts ...option.RequestOption) *pagination.PageAutoPager[BetaModelInfo] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

// Indicates whether a capability is supported.
type BetaCapabilitySupport struct {
	// Whether this capability is supported by the model.
	Supported bool `json:"supported" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Supported   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaCapabilitySupport) RawJSON() string { return r.JSON.raw }
func (r *BetaCapabilitySupport) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Context management capability details.
type BetaContextManagementCapability struct {
	// Indicates whether a capability is supported.
	ClearThinking20251015 BetaCapabilitySupport `json:"clear_thinking_20251015" api:"required"`
	// Indicates whether a capability is supported.
	ClearToolUses20250919 BetaCapabilitySupport `json:"clear_tool_uses_20250919" api:"required"`
	// Indicates whether a capability is supported.
	Compact20260112 BetaCapabilitySupport `json:"compact_20260112" api:"required"`
	// Whether this capability is supported by the model.
	Supported bool `json:"supported" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ClearThinking20251015 respjson.Field
		ClearToolUses20250919 respjson.Field
		Compact20260112       respjson.Field
		Supported             respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaContextManagementCapability) RawJSON() string { return r.JSON.raw }
func (r *BetaContextManagementCapability) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Effort (reasoning_effort) capability details.
type BetaEffortCapability struct {
	// Whether the model supports high effort level.
	High BetaCapabilitySupport `json:"high" api:"required"`
	// Whether the model supports low effort level.
	Low BetaCapabilitySupport `json:"low" api:"required"`
	// Whether the model supports max effort level.
	Max BetaCapabilitySupport `json:"max" api:"required"`
	// Whether the model supports medium effort level.
	Medium BetaCapabilitySupport `json:"medium" api:"required"`
	// Whether this capability is supported by the model.
	Supported bool `json:"supported" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		High        respjson.Field
		Low         respjson.Field
		Max         respjson.Field
		Medium      respjson.Field
		Supported   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaEffortCapability) RawJSON() string { return r.JSON.raw }
func (r *BetaEffortCapability) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Model capability information.
type BetaModelCapabilities struct {
	// Whether the model supports the Batch API.
	Batch BetaCapabilitySupport `json:"batch" api:"required"`
	// Whether the model supports citation generation.
	Citations BetaCapabilitySupport `json:"citations" api:"required"`
	// Whether the model supports code execution tools.
	CodeExecution BetaCapabilitySupport `json:"code_execution" api:"required"`
	// Context management support and available strategies.
	ContextManagement BetaContextManagementCapability `json:"context_management" api:"required"`
	// Effort (reasoning_effort) support and available levels.
	Effort BetaEffortCapability `json:"effort" api:"required"`
	// Whether the model accepts image content blocks.
	ImageInput BetaCapabilitySupport `json:"image_input" api:"required"`
	// Whether the model accepts PDF content blocks.
	PDFInput BetaCapabilitySupport `json:"pdf_input" api:"required"`
	// Whether the model supports structured output / JSON mode / strict tool schemas.
	StructuredOutputs BetaCapabilitySupport `json:"structured_outputs" api:"required"`
	// Thinking capability and supported type configurations.
	Thinking BetaThinkingCapability `json:"thinking" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Batch             respjson.Field
		Citations         respjson.Field
		CodeExecution     respjson.Field
		ContextManagement respjson.Field
		Effort            respjson.Field
		ImageInput        respjson.Field
		PDFInput          respjson.Field
		StructuredOutputs respjson.Field
		Thinking          respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaModelCapabilities) RawJSON() string { return r.JSON.raw }
func (r *BetaModelCapabilities) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaModelInfo struct {
	// Unique model identifier.
	ID string `json:"id" api:"required"`
	// Model capability information.
	Capabilities BetaModelCapabilities `json:"capabilities" api:"required"`
	// RFC 3339 datetime string representing the time at which the model was released.
	// May be set to an epoch value if the release date is unknown.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// A human-readable name for the model.
	DisplayName string `json:"display_name" api:"required"`
	// Maximum input context window size in tokens for this model.
	MaxInputTokens int64 `json:"max_input_tokens" api:"required"`
	// Maximum value for the `max_tokens` parameter when using this model.
	MaxTokens int64 `json:"max_tokens" api:"required"`
	// Object type.
	//
	// For Models, this is always `"model"`.
	Type constant.Model `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		Capabilities   respjson.Field
		CreatedAt      respjson.Field
		DisplayName    respjson.Field
		MaxInputTokens respjson.Field
		MaxTokens      respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaModelInfo) RawJSON() string { return r.JSON.raw }
func (r *BetaModelInfo) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Thinking capability details.
type BetaThinkingCapability struct {
	// Whether this capability is supported by the model.
	Supported bool `json:"supported" api:"required"`
	// Supported thinking type configurations.
	Types BetaThinkingTypes `json:"types" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Supported   respjson.Field
		Types       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaThinkingCapability) RawJSON() string { return r.JSON.raw }
func (r *BetaThinkingCapability) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Supported thinking type configurations.
type BetaThinkingTypes struct {
	// Whether the model supports thinking with type 'adaptive' (auto).
	Adaptive BetaCapabilitySupport `json:"adaptive" api:"required"`
	// Whether the model supports thinking with type 'enabled'.
	Enabled BetaCapabilitySupport `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Adaptive    respjson.Field
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaThinkingTypes) RawJSON() string { return r.JSON.raw }
func (r *BetaThinkingTypes) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaModelGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaModelListParams struct {
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately after this object.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately before this object.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `1000`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaModelListParams]'s query parameters as `url.Values`.
func (r BetaModelListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
