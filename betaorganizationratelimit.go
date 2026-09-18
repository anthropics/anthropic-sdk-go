package anthropic

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"slices"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// BetaOrganizationRateLimitService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationRateLimitService] method instead.
type BetaOrganizationRateLimitService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationRateLimitService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationRateLimitService(opts ...option.RequestOption) (r BetaOrganizationRateLimitService) {
	r = BetaOrganizationRateLimitService{}
	r.Options = opts
	return
}

// List Messages API rate limits for your organization.
//
// Each entry corresponds to one rate-limit group (either a model family or an
// API-surface category such as the Files API or Message Batches) and contains the
// set of limiter values that apply to it.
//
// When `limit` is omitted, every matching entry is returned in a single page; when
// `limit` truncates the result, follow `next_page` to fetch the remaining entries.
func (r *BetaOrganizationRateLimitService) List(ctx context.Context, query BetaOrganizationRateLimitListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaOrganizationRateLimit], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/rate_limits?beta=true"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
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

// List Messages API rate limits for your organization.
//
// Each entry corresponds to one rate-limit group (either a model family or an
// API-surface category such as the Files API or Message Batches) and contains the
// set of limiter values that apply to it.
//
// When `limit` is omitted, every matching entry is returned in a single page; when
// `limit` truncates the result, follow `next_page` to fetch the remaining entries.
func (r *BetaOrganizationRateLimitService) ListAutoPaging(ctx context.Context, query BetaOrganizationRateLimitListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaOrganizationRateLimit] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

type BetaOrganizationRateLimit struct {
	// Identifier of this rate-limit entry. It is stable within the organization and
	// differs between organizations; the group's own identifier is `group.id`.
	ID string `json:"id" api:"required"`
	// The rate-limit group this entry's limits apply to. Its `type` equals
	// `group_type`.
	Group BetaOrganizationRateLimitGroupUnion `json:"group" api:"required"`
	// Deprecated: use `group.type` instead. The kind of rate-limit group this entry
	// represents. `model_group` entries apply to a family of models (listed in
	// `models`); other values apply to an API-surface category and have `models` set
	// to `null`. Always equal to `group.type`.
	//
	// Any of "batch", "files", "model_group", "skills", "token_count", "web_search".
	//
	// Deprecated: Use `group.type` instead. `group_type` is still returned and always
	// equals `group.type`.
	GroupType BetaOrganizationRateLimitGroupType `json:"group_type" api:"required"`
	// The limiter values that apply to this group.
	Limits []BetaOrganizationRateLimitValue `json:"limits" api:"required"`
	// Model names this entry's limits apply to, including aliases. `null` when
	// `group_type` is not `"model_group"`.
	Models []string `json:"models" api:"required"`
	// Object type. Always `rate_limit` for organization rate-limit entries.
	Type constant.RateLimit `json:"type" default:"rate_limit"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Group       respjson.Field
		GroupType   respjson.Field
		Limits      respjson.Field
		Models      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationRateLimit) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationRateLimit) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaOrganizationRateLimitGroupUnion contains all possible properties and values
// from [BetaOrganizationRateLimitModelGroup],
// [BetaOrganizationRateLimitBatchGroup],
// [BetaOrganizationRateLimitTokenCountGroup],
// [BetaOrganizationRateLimitFilesGroup], [BetaOrganizationRateLimitSkillsGroup],
// [BetaOrganizationRateLimitWebSearchGroup].
//
// Use the [BetaOrganizationRateLimitGroupUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaOrganizationRateLimitGroupUnion struct {
	ID string `json:"id"`
	// This field is from variant [BetaOrganizationRateLimitModelGroup].
	DisplayName string `json:"display_name"`
	// Any of "model_group", "batch", "token_count", "files", "skills", "web_search".
	Type string `json:"type"`
	JSON struct {
		ID          respjson.Field
		DisplayName respjson.Field
		Type        respjson.Field
		raw         string
	} `json:"-"`
}

// anyBetaOrganizationRateLimitGroup is implemented by each variant of
// [BetaOrganizationRateLimitGroupUnion] to add type safety for the return type of
// [BetaOrganizationRateLimitGroupUnion.AsAny]
type anyBetaOrganizationRateLimitGroup interface {
	implBetaOrganizationRateLimitGroupUnion()
}

func (BetaOrganizationRateLimitModelGroup) implBetaOrganizationRateLimitGroupUnion()      {}
func (BetaOrganizationRateLimitBatchGroup) implBetaOrganizationRateLimitGroupUnion()      {}
func (BetaOrganizationRateLimitTokenCountGroup) implBetaOrganizationRateLimitGroupUnion() {}
func (BetaOrganizationRateLimitFilesGroup) implBetaOrganizationRateLimitGroupUnion()      {}
func (BetaOrganizationRateLimitSkillsGroup) implBetaOrganizationRateLimitGroupUnion()     {}
func (BetaOrganizationRateLimitWebSearchGroup) implBetaOrganizationRateLimitGroupUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaOrganizationRateLimitGroupUnion.AsAny().(type) {
//	case anthropic.BetaOrganizationRateLimitModelGroup:
//	case anthropic.BetaOrganizationRateLimitBatchGroup:
//	case anthropic.BetaOrganizationRateLimitTokenCountGroup:
//	case anthropic.BetaOrganizationRateLimitFilesGroup:
//	case anthropic.BetaOrganizationRateLimitSkillsGroup:
//	case anthropic.BetaOrganizationRateLimitWebSearchGroup:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaOrganizationRateLimitGroupUnion) AsAny() anyBetaOrganizationRateLimitGroup {
	switch u.Type {
	case "model_group":
		return u.AsModelGroup()
	case "batch":
		return u.AsBatch()
	case "token_count":
		return u.AsTokenCount()
	case "files":
		return u.AsFiles()
	case "skills":
		return u.AsSkills()
	case "web_search":
		return u.AsWebSearch()
	}
	return nil
}

func (u BetaOrganizationRateLimitGroupUnion) AsModelGroup() (v BetaOrganizationRateLimitModelGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaOrganizationRateLimitGroupUnion) AsBatch() (v BetaOrganizationRateLimitBatchGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaOrganizationRateLimitGroupUnion) AsTokenCount() (v BetaOrganizationRateLimitTokenCountGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaOrganizationRateLimitGroupUnion) AsFiles() (v BetaOrganizationRateLimitFilesGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaOrganizationRateLimitGroupUnion) AsSkills() (v BetaOrganizationRateLimitSkillsGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaOrganizationRateLimitGroupUnion) AsWebSearch() (v BetaOrganizationRateLimitWebSearchGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaOrganizationRateLimitGroupUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaOrganizationRateLimitGroupUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Deprecated: use `group.type` instead. The kind of rate-limit group this entry
// represents. `model_group` entries apply to a family of models (listed in
// `models`); other values apply to an API-surface category and have `models` set
// to `null`. Always equal to `group.type`.
type BetaOrganizationRateLimitGroupType string

const (
	BetaOrganizationRateLimitGroupTypeBatch      BetaOrganizationRateLimitGroupType = "batch"
	BetaOrganizationRateLimitGroupTypeFiles      BetaOrganizationRateLimitGroupType = "files"
	BetaOrganizationRateLimitGroupTypeModelGroup BetaOrganizationRateLimitGroupType = "model_group"
	BetaOrganizationRateLimitGroupTypeSkills     BetaOrganizationRateLimitGroupType = "skills"
	BetaOrganizationRateLimitGroupTypeTokenCount BetaOrganizationRateLimitGroupType = "token_count"
	BetaOrganizationRateLimitGroupTypeWebSearch  BetaOrganizationRateLimitGroupType = "web_search"
)

type BetaOrganizationRateLimitBatchGroup struct {
	// Opaque identifier of the rate-limit group (for example,
	// `rlg_01VPTCmyiu5ZLsWkcxYG2pY8`). It is the same in every organization and never
	// changes, unlike the entry's own identifier, which differs per organization.
	ID string `json:"id" api:"required"`
	// Always `batch`: the Message Batches API.
	Type constant.Batch `json:"type" default:"batch"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationRateLimitBatchGroup) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationRateLimitBatchGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationRateLimitFilesGroup struct {
	// Opaque identifier of the rate-limit group (for example,
	// `rlg_01VPTCmyiu5ZLsWkcxYG2pY8`). It is the same in every organization and never
	// changes, unlike the entry's own identifier, which differs per organization.
	ID string `json:"id" api:"required"`
	// Always `files`: the Files API.
	Type constant.Files `json:"type" default:"files"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationRateLimitFilesGroup) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationRateLimitFilesGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationRateLimitModelGroup struct {
	// Opaque identifier of the rate-limit group (for example,
	// `rlg_01VPTCmyiu5ZLsWkcxYG2pY8`). It is the same in every organization and never
	// changes, unlike the entry's own identifier, which differs per organization.
	ID string `json:"id" api:"required"`
	// Human-readable name of the model group (for example, `Claude Sonnet 4.x`). For
	// display only; it may change.
	DisplayName string `json:"display_name" api:"required"`
	// Always `model_group`: a family of models.
	Type constant.ModelGroup `json:"type" default:"model_group"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		DisplayName respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationRateLimitModelGroup) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationRateLimitModelGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationRateLimitSkillsGroup struct {
	// Opaque identifier of the rate-limit group (for example,
	// `rlg_01VPTCmyiu5ZLsWkcxYG2pY8`). It is the same in every organization and never
	// changes, unlike the entry's own identifier, which differs per organization.
	ID string `json:"id" api:"required"`
	// Always `skills`: the Skills API.
	Type constant.Skills `json:"type" default:"skills"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationRateLimitSkillsGroup) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationRateLimitSkillsGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationRateLimitTokenCountGroup struct {
	// Opaque identifier of the rate-limit group (for example,
	// `rlg_01VPTCmyiu5ZLsWkcxYG2pY8`). It is the same in every organization and never
	// changes, unlike the entry's own identifier, which differs per organization.
	ID string `json:"id" api:"required"`
	// Always `token_count`: the Token Count API.
	Type constant.TokenCount `json:"type" default:"token_count"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationRateLimitTokenCountGroup) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationRateLimitTokenCountGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationRateLimitValue struct {
	// The limiter type (for example, `requests_per_minute` or
	// `input_tokens_per_minute`).
	Type string `json:"type" api:"required"`
	// The configured limit value for this limiter type.
	Value int64 `json:"value" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationRateLimitValue) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationRateLimitValue) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationRateLimitWebSearchGroup struct {
	// Opaque identifier of the rate-limit group (for example,
	// `rlg_01VPTCmyiu5ZLsWkcxYG2pY8`). It is the same in every organization and never
	// changes, unlike the entry's own identifier, which differs per organization.
	ID string `json:"id" api:"required"`
	// Always `web_search`: the Messages API web search tool.
	Type constant.WebSearch `json:"type" default:"web_search"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationRateLimitWebSearchGroup) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationRateLimitWebSearchGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationRateLimitListParams struct {
	// Maximum number of items to return per page. Ranges from `1` to `1000`.
	//
	// When omitted, every remaining entry is returned in a single page and `next_page`
	// is `null`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Filter to the single entry containing this model. Accepts full model names and
	// aliases. Returns 404 if the model is not found or has no rate limits for this
	// organization.
	Model param.Opt[string] `query:"model,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Filter by group type.
	//
	// Any of "batch", "files", "model_group", "skills", "token_count", "web_search".
	GroupType BetaOrganizationRateLimitListParamsGroupType `query:"group_type,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationRateLimitListParams]'s query parameters as
// `url.Values`.
func (r BetaOrganizationRateLimitListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Filter by group type.
type BetaOrganizationRateLimitListParamsGroupType string

const (
	BetaOrganizationRateLimitListParamsGroupTypeBatch      BetaOrganizationRateLimitListParamsGroupType = "batch"
	BetaOrganizationRateLimitListParamsGroupTypeFiles      BetaOrganizationRateLimitListParamsGroupType = "files"
	BetaOrganizationRateLimitListParamsGroupTypeModelGroup BetaOrganizationRateLimitListParamsGroupType = "model_group"
	BetaOrganizationRateLimitListParamsGroupTypeSkills     BetaOrganizationRateLimitListParamsGroupType = "skills"
	BetaOrganizationRateLimitListParamsGroupTypeTokenCount BetaOrganizationRateLimitListParamsGroupType = "token_count"
	BetaOrganizationRateLimitListParamsGroupTypeWebSearch  BetaOrganizationRateLimitListParamsGroupType = "web_search"
)
