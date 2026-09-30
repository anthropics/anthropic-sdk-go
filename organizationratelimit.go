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

// OrganizationRateLimitService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationRateLimitService] method instead.
type OrganizationRateLimitService struct {
	Options []option.RequestOption
}

// NewOrganizationRateLimitService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewOrganizationRateLimitService(opts ...option.RequestOption) (r OrganizationRateLimitService) {
	r = OrganizationRateLimitService{}
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
func (r *OrganizationRateLimitService) List(ctx context.Context, query OrganizationRateLimitListParams, opts ...option.RequestOption) (res *pagination.PageCursor[OrganizationRateLimit], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/rate_limits"
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
func (r *OrganizationRateLimitService) ListAutoPaging(ctx context.Context, query OrganizationRateLimitListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[OrganizationRateLimit] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

type OrganizationRateLimit struct {
	// Identifier of this rate-limit entry. It is stable within the organization and
	// differs between organizations; the group's own identifier is `group.id`.
	ID string `json:"id" api:"required"`
	// The rate-limit group this entry's limits apply to. Its `type` equals
	// `group_type`.
	Group OrganizationRateLimitGroupUnion `json:"group" api:"required"`
	// The limiter values that apply to this group.
	Limits []OrganizationRateLimitValue `json:"limits" api:"required"`
	// Model names this entry's limits apply to, including aliases. `null` when
	// `group_type` is not `"model_group"`.
	Models []string `json:"models" api:"required"`
	// Object type. Always `rate_limit` for organization rate-limit entries.
	Type constant.RateLimit `json:"type" default:"rate_limit"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Group       respjson.Field
		Limits      respjson.Field
		Models      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationRateLimit) RawJSON() string { return r.JSON.raw }
func (r *OrganizationRateLimit) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// OrganizationRateLimitGroupUnion contains all possible properties and values from
// [OrganizationRateLimitModelGroup], [OrganizationRateLimitBatchGroup],
// [OrganizationRateLimitTokenCountGroup], [OrganizationRateLimitFilesGroup],
// [OrganizationRateLimitSkillsGroup], [OrganizationRateLimitWebSearchGroup].
//
// Use the [OrganizationRateLimitGroupUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type OrganizationRateLimitGroupUnion struct {
	ID string `json:"id"`
	// This field is from variant [OrganizationRateLimitModelGroup].
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

// anyOrganizationRateLimitGroup is implemented by each variant of
// [OrganizationRateLimitGroupUnion] to add type safety for the return type of
// [OrganizationRateLimitGroupUnion.AsAny]
type anyOrganizationRateLimitGroup interface {
	implOrganizationRateLimitGroupUnion()
}

func (OrganizationRateLimitModelGroup) implOrganizationRateLimitGroupUnion()      {}
func (OrganizationRateLimitBatchGroup) implOrganizationRateLimitGroupUnion()      {}
func (OrganizationRateLimitTokenCountGroup) implOrganizationRateLimitGroupUnion() {}
func (OrganizationRateLimitFilesGroup) implOrganizationRateLimitGroupUnion()      {}
func (OrganizationRateLimitSkillsGroup) implOrganizationRateLimitGroupUnion()     {}
func (OrganizationRateLimitWebSearchGroup) implOrganizationRateLimitGroupUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := OrganizationRateLimitGroupUnion.AsAny().(type) {
//	case anthropic.OrganizationRateLimitModelGroup:
//	case anthropic.OrganizationRateLimitBatchGroup:
//	case anthropic.OrganizationRateLimitTokenCountGroup:
//	case anthropic.OrganizationRateLimitFilesGroup:
//	case anthropic.OrganizationRateLimitSkillsGroup:
//	case anthropic.OrganizationRateLimitWebSearchGroup:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u OrganizationRateLimitGroupUnion) AsAny() anyOrganizationRateLimitGroup {
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

func (u OrganizationRateLimitGroupUnion) AsModelGroup() (v OrganizationRateLimitModelGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u OrganizationRateLimitGroupUnion) AsBatch() (v OrganizationRateLimitBatchGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u OrganizationRateLimitGroupUnion) AsTokenCount() (v OrganizationRateLimitTokenCountGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u OrganizationRateLimitGroupUnion) AsFiles() (v OrganizationRateLimitFilesGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u OrganizationRateLimitGroupUnion) AsSkills() (v OrganizationRateLimitSkillsGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u OrganizationRateLimitGroupUnion) AsWebSearch() (v OrganizationRateLimitWebSearchGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u OrganizationRateLimitGroupUnion) RawJSON() string { return u.JSON.raw }

func (r *OrganizationRateLimitGroupUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationRateLimitBatchGroup struct {
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
func (r OrganizationRateLimitBatchGroup) RawJSON() string { return r.JSON.raw }
func (r *OrganizationRateLimitBatchGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationRateLimitFilesGroup struct {
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
func (r OrganizationRateLimitFilesGroup) RawJSON() string { return r.JSON.raw }
func (r *OrganizationRateLimitFilesGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationRateLimitModelGroup struct {
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
func (r OrganizationRateLimitModelGroup) RawJSON() string { return r.JSON.raw }
func (r *OrganizationRateLimitModelGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationRateLimitSkillsGroup struct {
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
func (r OrganizationRateLimitSkillsGroup) RawJSON() string { return r.JSON.raw }
func (r *OrganizationRateLimitSkillsGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationRateLimitTokenCountGroup struct {
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
func (r OrganizationRateLimitTokenCountGroup) RawJSON() string { return r.JSON.raw }
func (r *OrganizationRateLimitTokenCountGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationRateLimitValue struct {
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
func (r OrganizationRateLimitValue) RawJSON() string { return r.JSON.raw }
func (r *OrganizationRateLimitValue) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationRateLimitWebSearchGroup struct {
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
func (r OrganizationRateLimitWebSearchGroup) RawJSON() string { return r.JSON.raw }
func (r *OrganizationRateLimitWebSearchGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationRateLimitListParams struct {
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
	GroupType OrganizationRateLimitListParamsGroupType `query:"group_type,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrganizationRateLimitListParams]'s query parameters as
// `url.Values`.
func (r OrganizationRateLimitListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Filter by group type.
type OrganizationRateLimitListParamsGroupType string

const (
	OrganizationRateLimitListParamsGroupTypeBatch      OrganizationRateLimitListParamsGroupType = "batch"
	OrganizationRateLimitListParamsGroupTypeFiles      OrganizationRateLimitListParamsGroupType = "files"
	OrganizationRateLimitListParamsGroupTypeModelGroup OrganizationRateLimitListParamsGroupType = "model_group"
	OrganizationRateLimitListParamsGroupTypeSkills     OrganizationRateLimitListParamsGroupType = "skills"
	OrganizationRateLimitListParamsGroupTypeTokenCount OrganizationRateLimitListParamsGroupType = "token_count"
	OrganizationRateLimitListParamsGroupTypeWebSearch  OrganizationRateLimitListParamsGroupType = "web_search"
)
