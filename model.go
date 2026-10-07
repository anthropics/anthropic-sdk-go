package anthropic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
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

// ModelService contains methods and other services that help with interacting with
// the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewModelService] method instead.
type ModelService struct {
	Options []option.RequestOption
}

// NewModelService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewModelService(opts ...option.RequestOption) (r ModelService) {
	r = ModelService{}
	r.Options = opts
	return
}

// Get a specific model.
//
// The Models API response can be used to determine information about a specific
// model or resolve a model alias to a model ID.
func (r *ModelService) Get(ctx context.Context, modelID string, query ModelGetParams, opts ...option.RequestOption) (res *ModelInfo, err error) {
	if len(query.Betas) > 0 {
		headerValues := make([]string, len(query.Betas))
		for i, v := range query.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	if !param.IsOmitted(query.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", query.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	if modelID == "" {
		err = errors.New("missing required model_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/models/%s", url.PathEscape(modelID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List available models.
//
// The Models API response can be used to determine which models are available for
// use in the API. More recently released models are listed first.
func (r *ModelService) List(ctx context.Context, params ModelListParams, opts ...option.RequestOption) (res *pagination.Page[ModelInfo], err error) {
	var raw *http.Response
	if len(params.Betas) > 0 {
		headerValues := make([]string, len(params.Betas))
		for i, v := range params.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	if !param.IsOmitted(params.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", params.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/models"
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
func (r *ModelService) ListAutoPaging(ctx context.Context, params ModelListParams, opts ...option.RequestOption) *pagination.PageAutoPager[ModelInfo] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

// Indicates whether a capability is supported.
type CapabilitySupport struct {
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
func (r CapabilitySupport) RawJSON() string { return r.JSON.raw }
func (r *CapabilitySupport) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Context management capability details.
type ContextManagementCapability struct {
	// Whether the clear_thinking_20251015 strategy is supported.
	ClearThinking20251015 CapabilitySupport `json:"clear_thinking_20251015" api:"required"`
	// Whether the clear_tool_uses_20250919 strategy is supported.
	ClearToolUses20250919 CapabilitySupport `json:"clear_tool_uses_20250919" api:"required"`
	// Whether the compact_20260112 strategy is supported.
	Compact20260112 CapabilitySupport `json:"compact_20260112" api:"required"`
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
func (r ContextManagementCapability) RawJSON() string { return r.JSON.raw }
func (r *ContextManagementCapability) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Effort (reasoning_effort) capability details.
type EffortCapability struct {
	// Whether the model supports high effort level.
	High CapabilitySupport `json:"high" api:"required"`
	// Whether the model supports low effort level.
	Low CapabilitySupport `json:"low" api:"required"`
	// Whether the model supports max effort level.
	Max CapabilitySupport `json:"max" api:"required"`
	// Whether the model supports medium effort level.
	Medium CapabilitySupport `json:"medium" api:"required"`
	// Whether this capability is supported by the model.
	Supported bool `json:"supported" api:"required"`
	// Whether the model supports xhigh effort level.
	Xhigh CapabilitySupport `json:"xhigh" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		High        respjson.Field
		Low         respjson.Field
		Max         respjson.Field
		Medium      respjson.Field
		Supported   respjson.Field
		Xhigh       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EffortCapability) RawJSON() string { return r.JSON.raw }
func (r *EffortCapability) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Model capability information.
type ModelCapabilities struct {
	// Whether the model supports the Batch API.
	Batch CapabilitySupport `json:"batch" api:"required"`
	// Whether the model supports citation generation.
	Citations CapabilitySupport `json:"citations" api:"required"`
	// Whether code that the model runs in the code execution tool can call the
	// request's other tools, as in programmatic tool calling and dynamic filtering for
	// web search and web fetch. Support for the code execution tool itself is in
	// `server_tools.code_execution`.
	CodeExecution CapabilitySupport `json:"code_execution" api:"required"`
	// Context management support and available strategies.
	ContextManagement ContextManagementCapability `json:"context_management" api:"required"`
	// Effort (reasoning_effort) support and available levels.
	Effort EffortCapability `json:"effort" api:"required"`
	// Whether the model accepts image content blocks.
	ImageInput CapabilitySupport `json:"image_input" api:"required"`
	// Whether the model accepts PDF content blocks.
	PDFInput CapabilitySupport `json:"pdf_input" api:"required"`
	// Whether this model supports the web search and code execution server tools.
	// `supported` is true when the model supports at least one of the tools. A
	// supported tool can still be rejected for your organization, for example when an
	// admin has turned web search off.
	ServerTools ServerToolsCapability `json:"server_tools" api:"required"`
	// Whether the model supports structured output / JSON mode / strict tool schemas.
	StructuredOutputs CapabilitySupport `json:"structured_outputs" api:"required"`
	// Thinking capability and supported type configurations.
	Thinking ThinkingCapability `json:"thinking" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Batch             respjson.Field
		Citations         respjson.Field
		CodeExecution     respjson.Field
		ContextManagement respjson.Field
		Effort            respjson.Field
		ImageInput        respjson.Field
		PDFInput          respjson.Field
		ServerTools       respjson.Field
		StructuredOutputs respjson.Field
		Thinking          respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ModelCapabilities) RawJSON() string { return r.JSON.raw }
func (r *ModelCapabilities) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ModelInfo struct {
	// Unique model identifier.
	ID string `json:"id" api:"required"`
	// Object mapping capability names to their support details. Keys are always
	// present for all known capabilities.
	Capabilities ModelCapabilities `json:"capabilities" api:"required"`
	// RFC 3339 datetime string representing the time at which the model was released.
	// May be set to an epoch value if the release date is unknown.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// RFC 3339 datetime string representing the time of the model's most recent
	// deprecation. Populated for `deprecated` and `retired` models; `null` while the
	// model is `active`.
	DeprecatedAt time.Time `json:"deprecated_at" api:"required" format:"date-time"`
	// A human-readable name for the model.
	DisplayName string `json:"display_name" api:"required"`
	// The model's current lifecycle stage.
	//
	// - `active`: The model is available for use, open to new adopters, and not
	//   scheduled for retirement.
	// - `deprecated`: The model remains callable for organizations with existing
	//   access, but is headed for retirement and closed to new adopters.
	// - `retired`: The model is no longer available for use; inference requests naming
	//   it fail. It remains in the catalogue as the historical record of its
	//   retirement.
	//
	// Any of "active", "deprecated", "retired".
	Lifecycle ModelInfoLifecycle `json:"lifecycle" api:"required"`
	// The model line this model belongs to, such as `opus` for both Claude Opus 4.5
	// and Claude Opus 4.6. More lines may be added. `null` when the model belongs to
	// no line; do not infer a line from the `id`.
	//
	// Any of "haiku", "sonnet", "opus", "fable", "mythos".
	Line ModelLine `json:"line" api:"required"`
	// Maximum input context window size in tokens for this model.
	MaxInputTokens int64 `json:"max_input_tokens" api:"required"`
	// Maximum value for the `max_tokens` parameter when using this model.
	MaxTokens int64 `json:"max_tokens" api:"required"`
	// RFC 3339 datetime string representing the model's currently scheduled retirement
	// date. The schedule can be revised until retirement occurs; `null` while the
	// model is `active` or while no retirement is scheduled. A past date on a
	// `deprecated` model means retirement is overdue, not that it has occurred:
	// `lifecycle` is the retirement signal.
	RetiresAt time.Time `json:"retires_at" api:"required" format:"date-time"`
	// Object type.
	//
	// For Models, this is always `"model"`.
	Type constant.Model `json:"type" default:"model"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		Capabilities   respjson.Field
		CreatedAt      respjson.Field
		DeprecatedAt   respjson.Field
		DisplayName    respjson.Field
		Lifecycle      respjson.Field
		Line           respjson.Field
		MaxInputTokens respjson.Field
		MaxTokens      respjson.Field
		RetiresAt      respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ModelInfo) RawJSON() string { return r.JSON.raw }
func (r *ModelInfo) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The model's current lifecycle stage.
//
//   - `active`: The model is available for use, open to new adopters, and not
//     scheduled for retirement.
//   - `deprecated`: The model remains callable for organizations with existing
//     access, but is headed for retirement and closed to new adopters.
//   - `retired`: The model is no longer available for use; inference requests naming
//     it fail. It remains in the catalogue as the historical record of its
//     retirement.
type ModelInfoLifecycle string

const (
	ModelInfoLifecycleActive     ModelInfoLifecycle = "active"
	ModelInfoLifecycleDeprecated ModelInfoLifecycle = "deprecated"
	ModelInfoLifecycleRetired    ModelInfoLifecycle = "retired"
)

// A Claude model line, such as `opus` or `sonnet`. More lines may be added as new
// values.
type ModelLine string

const (
	ModelLineHaiku  ModelLine = "haiku"
	ModelLineSonnet ModelLine = "sonnet"
	ModelLineOpus   ModelLine = "opus"
	ModelLineFable  ModelLine = "fable"
	ModelLineMythos ModelLine = "mythos"
)

// Web search and code execution tool support, with one entry per tool.
type ServerToolsCapability struct {
	// Whether the model supports the code execution tool: true when the model supports
	// at least one version of the tool, not necessarily every version.
	CodeExecution CapabilitySupport `json:"code_execution" api:"required"`
	// Whether this capability is supported by the model.
	Supported bool `json:"supported" api:"required"`
	// Whether the model supports the web search tool: true when the model supports at
	// least one version of the tool, not necessarily every version.
	WebSearch CapabilitySupport `json:"web_search" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CodeExecution respjson.Field
		Supported     respjson.Field
		WebSearch     respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ServerToolsCapability) RawJSON() string { return r.JSON.raw }
func (r *ServerToolsCapability) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Thinking capability details.
type ThinkingCapability struct {
	// Whether this capability is supported by the model.
	Supported bool `json:"supported" api:"required"`
	// Supported thinking type configurations.
	Types ThinkingTypes `json:"types" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Supported   respjson.Field
		Types       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ThinkingCapability) RawJSON() string { return r.JSON.raw }
func (r *ThinkingCapability) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Which `thinking.type` values the model accepts on requests. Read each key on its
// own: for example, `enabled` can be false while `disabled` is true.
type ThinkingTypes struct {
	// Whether the model accepts thinking with type 'adaptive' (the model decides
	// whether and how much to think).
	Adaptive CapabilitySupport `json:"adaptive" api:"required"`
	// Whether the model accepts thinking with type 'disabled' (thinking turned off).
	// False exactly when a request that sends it gets a 400 from this model. True on a
	// model that does not support thinking.
	Disabled CapabilitySupport `json:"disabled" api:"required"`
	// Whether the model accepts thinking with type 'enabled' (extended thinking with a
	// caller-set `budget_tokens`).
	Enabled CapabilitySupport `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Adaptive    respjson.Field
		Disabled    respjson.Field
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ThinkingTypes) RawJSON() string { return r.JSON.raw }
func (r *ThinkingTypes) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ModelGetParams struct {
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type ModelListParams struct {
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
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// Filter the list to models in any of the given lifecycle stages (`active`,
	// `deprecated`, or `retired`). Up to 3 values. When omitted, the list contains the
	// `active` and `deprecated` models; `retired` models appear only when `retired` is
	// requested explicitly.
	//
	// Any of "active", "deprecated", "retired".
	Lifecycle []string `query:"lifecycle,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [ModelListParams]'s query parameters as `url.Values`.
func (r ModelListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
