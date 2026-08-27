// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"context"
	"errors"
	"fmt"
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

// BetaOrganizationWorkspaceRateLimitService contains methods and other services
// that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationWorkspaceRateLimitService] method instead.
type BetaOrganizationWorkspaceRateLimitService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationWorkspaceRateLimitService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationWorkspaceRateLimitService(opts ...option.RequestOption) (r BetaOrganizationWorkspaceRateLimitService) {
	r = BetaOrganizationWorkspaceRateLimitService{}
	r.Options = opts
	return
}

// List rate-limit overrides configured for a workspace.
//
// Returns only the groups and limiter types that have a workspace-level override.
// Groups without overrides inherit the organization limits and are not listed; use
// `GET /v1/organizations/rate_limits` to see those.
//
// When `limit` is omitted, every matching entry is returned in a single page; when
// `limit` truncates the result, follow `next_page` to fetch the remaining entries.
func (r *BetaOrganizationWorkspaceRateLimitService) List(ctx context.Context, workspaceID string, query BetaOrganizationWorkspaceRateLimitListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaWorkspaceRateLimit], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s/rate_limits?beta=true", workspaceID)
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

// List rate-limit overrides configured for a workspace.
//
// Returns only the groups and limiter types that have a workspace-level override.
// Groups without overrides inherit the organization limits and are not listed; use
// `GET /v1/organizations/rate_limits` to see those.
//
// When `limit` is omitted, every matching entry is returned in a single page; when
// `limit` truncates the result, follow `next_page` to fetch the remaining entries.
func (r *BetaOrganizationWorkspaceRateLimitService) ListAutoPaging(ctx context.Context, workspaceID string, query BetaOrganizationWorkspaceRateLimitListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaWorkspaceRateLimit] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, workspaceID, query, opts...))
}

type BetaWorkspaceRateLimit struct {
	// The kind of rate-limit group this entry represents. `model_group` entries apply
	// to a family of models (listed in `models`); other values apply to an API-surface
	// category and have `models` set to `null`.
	//
	// Any of "batch", "files", "model_group", "skills", "token_count", "web_search".
	GroupType BetaWorkspaceRateLimitGroupType `json:"group_type" api:"required"`
	// The limiter values overridden for this group in this workspace. Limiter types
	// without a workspace override are omitted and inherit the organization value.
	Limits []BetaWorkspaceRateLimitValue `json:"limits" api:"required"`
	// Model names this entry's limits apply to, including aliases. `null` when
	// `group_type` is not `"model_group"`.
	Models []string `json:"models" api:"required"`
	// The `id` of the RateLimit group this override applies to.
	RateLimitID string `json:"rate_limit_id" api:"required"`
	// Object type. Always `workspace_rate_limit` for workspace rate-limit entries.
	Type constant.WorkspaceRateLimit `json:"type" default:"workspace_rate_limit"`
	// ID of the Workspace this override applies to.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		GroupType   respjson.Field
		Limits      respjson.Field
		Models      respjson.Field
		RateLimitID respjson.Field
		Type        respjson.Field
		WorkspaceID respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWorkspaceRateLimit) RawJSON() string { return r.JSON.raw }
func (r *BetaWorkspaceRateLimit) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The kind of rate-limit group this entry represents. `model_group` entries apply
// to a family of models (listed in `models`); other values apply to an API-surface
// category and have `models` set to `null`.
type BetaWorkspaceRateLimitGroupType string

const (
	BetaWorkspaceRateLimitGroupTypeBatch      BetaWorkspaceRateLimitGroupType = "batch"
	BetaWorkspaceRateLimitGroupTypeFiles      BetaWorkspaceRateLimitGroupType = "files"
	BetaWorkspaceRateLimitGroupTypeModelGroup BetaWorkspaceRateLimitGroupType = "model_group"
	BetaWorkspaceRateLimitGroupTypeSkills     BetaWorkspaceRateLimitGroupType = "skills"
	BetaWorkspaceRateLimitGroupTypeTokenCount BetaWorkspaceRateLimitGroupType = "token_count"
	BetaWorkspaceRateLimitGroupTypeWebSearch  BetaWorkspaceRateLimitGroupType = "web_search"
)

type BetaWorkspaceRateLimitValue struct {
	// The organization-level value for the same limiter type, for reference. `null`
	// when the organization has no limit configured for this limiter type.
	OrgLimit int64 `json:"org_limit" api:"required"`
	// The limiter type (for example, `requests_per_minute` or
	// `input_tokens_per_minute`).
	Type string `json:"type" api:"required"`
	// The workspace-level override value for this limiter type.
	Value int64 `json:"value" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		OrgLimit    respjson.Field
		Type        respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWorkspaceRateLimitValue) RawJSON() string { return r.JSON.raw }
func (r *BetaWorkspaceRateLimitValue) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationWorkspaceRateLimitListParams struct {
	// Maximum number of items to return per page. Ranges from `1` to `1000`.
	//
	// When omitted, every remaining entry is returned in a single page and `next_page`
	// is `null`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Filter by group type.
	//
	// Any of "batch", "files", "model_group", "skills", "token_count", "web_search".
	GroupType BetaOrganizationWorkspaceRateLimitListParamsGroupType `query:"group_type,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationWorkspaceRateLimitListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationWorkspaceRateLimitListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Filter by group type.
type BetaOrganizationWorkspaceRateLimitListParamsGroupType string

const (
	BetaOrganizationWorkspaceRateLimitListParamsGroupTypeBatch      BetaOrganizationWorkspaceRateLimitListParamsGroupType = "batch"
	BetaOrganizationWorkspaceRateLimitListParamsGroupTypeFiles      BetaOrganizationWorkspaceRateLimitListParamsGroupType = "files"
	BetaOrganizationWorkspaceRateLimitListParamsGroupTypeModelGroup BetaOrganizationWorkspaceRateLimitListParamsGroupType = "model_group"
	BetaOrganizationWorkspaceRateLimitListParamsGroupTypeSkills     BetaOrganizationWorkspaceRateLimitListParamsGroupType = "skills"
	BetaOrganizationWorkspaceRateLimitListParamsGroupTypeTokenCount BetaOrganizationWorkspaceRateLimitListParamsGroupType = "token_count"
	BetaOrganizationWorkspaceRateLimitListParamsGroupTypeWebSearch  BetaOrganizationWorkspaceRateLimitListParamsGroupType = "web_search"
)
