// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"context"
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
	// Stable identifier for this rate-limit group within the organization.
	ID string `json:"id" api:"required"`
	// The kind of rate-limit group this entry represents. `model_group` entries apply
	// to a family of models (listed in `models`); other values apply to an API-surface
	// category and have `models` set to `null`.
	//
	// Any of "batch", "files", "model_group", "skills", "token_count", "web_search".
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

// The kind of rate-limit group this entry represents. `model_group` entries apply
// to a family of models (listed in `models`); other values apply to an API-surface
// category and have `models` set to `null`.
type BetaOrganizationRateLimitGroupType string

const (
	BetaOrganizationRateLimitGroupTypeBatch      BetaOrganizationRateLimitGroupType = "batch"
	BetaOrganizationRateLimitGroupTypeFiles      BetaOrganizationRateLimitGroupType = "files"
	BetaOrganizationRateLimitGroupTypeModelGroup BetaOrganizationRateLimitGroupType = "model_group"
	BetaOrganizationRateLimitGroupTypeSkills     BetaOrganizationRateLimitGroupType = "skills"
	BetaOrganizationRateLimitGroupTypeTokenCount BetaOrganizationRateLimitGroupType = "token_count"
	BetaOrganizationRateLimitGroupTypeWebSearch  BetaOrganizationRateLimitGroupType = "web_search"
)

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
