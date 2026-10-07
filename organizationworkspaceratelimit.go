package anthropic

import (
	"context"
	"encoding/json"
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

// OrganizationWorkspaceRateLimitService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationWorkspaceRateLimitService] method instead.
type OrganizationWorkspaceRateLimitService struct {
	Options []option.RequestOption
}

// NewOrganizationWorkspaceRateLimitService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewOrganizationWorkspaceRateLimitService(opts ...option.RequestOption) (r OrganizationWorkspaceRateLimitService) {
	r = OrganizationWorkspaceRateLimitService{}
	r.Options = opts
	return
}

// List a workspace's rate limits.
//
// By default, returns only the groups and limiter types that have a
// workspace-level override. With `include_inherited=true`, returns every group
// with organization-level limits the workspace can see, listing for each the
// values it inherits from the organization as well as its own overrides. Each
// value's `source` says which it is.
//
// When `limit` is omitted, every matching entry is returned in a single page; when
// `limit` truncates the result, follow `next_page` to fetch the remaining entries.
func (r *OrganizationWorkspaceRateLimitService) List(ctx context.Context, workspaceID string, query OrganizationWorkspaceRateLimitListParams, opts ...option.RequestOption) (res *pagination.PageCursor[WorkspaceRateLimit], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s/rate_limits", url.PathEscape(workspaceID))
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

// List a workspace's rate limits.
//
// By default, returns only the groups and limiter types that have a
// workspace-level override. With `include_inherited=true`, returns every group
// with organization-level limits the workspace can see, listing for each the
// values it inherits from the organization as well as its own overrides. Each
// value's `source` says which it is.
//
// When `limit` is omitted, every matching entry is returned in a single page; when
// `limit` truncates the result, follow `next_page` to fetch the remaining entries.
func (r *OrganizationWorkspaceRateLimitService) ListAutoPaging(ctx context.Context, workspaceID string, query OrganizationWorkspaceRateLimitListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[WorkspaceRateLimit] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, workspaceID, query, opts...))
}

type WorkspaceRateLimit struct {
	// The rate-limit group this entry's limits apply to. Its `type` equals
	// `group_type`.
	Group WorkspaceRateLimitGroupUnion `json:"group" api:"required"`
	// The workspace's limiter values for this group. By default only the limiter types
	// with a workspace-level override are listed. With `include_inherited` set to
	// `true`, the limiter types the workspace inherits from the organization are
	// listed too, each marked by `source`.
	Limits []WorkspaceRateLimitValue `json:"limits" api:"required"`
	// Model names this entry's limits apply to, including aliases. `null` when
	// `group_type` is not `"model_group"`.
	Models []string `json:"models" api:"required"`
	// The `id` of the organization's RateLimit entry this entry applies to.
	RateLimitID string `json:"rate_limit_id" api:"required"`
	// Object type. Always `workspace_rate_limit` for workspace rate-limit entries.
	Type constant.WorkspaceRateLimit `json:"type" default:"workspace_rate_limit"`
	// ID of the Workspace this entry applies to.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Group       respjson.Field
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
func (r WorkspaceRateLimit) RawJSON() string { return r.JSON.raw }
func (r *WorkspaceRateLimit) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// WorkspaceRateLimitGroupUnion contains all possible properties and values from
// [OrganizationRateLimitModelGroup], [OrganizationRateLimitBatchGroup],
// [OrganizationRateLimitTokenCountGroup], [OrganizationRateLimitFilesGroup],
// [OrganizationRateLimitSkillsGroup], [OrganizationRateLimitWebSearchGroup].
//
// Use the [WorkspaceRateLimitGroupUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type WorkspaceRateLimitGroupUnion struct {
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

// anyWorkspaceRateLimitGroup is implemented by each variant of
// [WorkspaceRateLimitGroupUnion] to add type safety for the return type of
// [WorkspaceRateLimitGroupUnion.AsAny]
type anyWorkspaceRateLimitGroup interface {
	implWorkspaceRateLimitGroupUnion()
}

func (OrganizationRateLimitModelGroup) implWorkspaceRateLimitGroupUnion()      {}
func (OrganizationRateLimitBatchGroup) implWorkspaceRateLimitGroupUnion()      {}
func (OrganizationRateLimitTokenCountGroup) implWorkspaceRateLimitGroupUnion() {}
func (OrganizationRateLimitFilesGroup) implWorkspaceRateLimitGroupUnion()      {}
func (OrganizationRateLimitSkillsGroup) implWorkspaceRateLimitGroupUnion()     {}
func (OrganizationRateLimitWebSearchGroup) implWorkspaceRateLimitGroupUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := WorkspaceRateLimitGroupUnion.AsAny().(type) {
//	case anthropic.OrganizationRateLimitModelGroup:
//	case anthropic.OrganizationRateLimitBatchGroup:
//	case anthropic.OrganizationRateLimitTokenCountGroup:
//	case anthropic.OrganizationRateLimitFilesGroup:
//	case anthropic.OrganizationRateLimitSkillsGroup:
//	case anthropic.OrganizationRateLimitWebSearchGroup:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u WorkspaceRateLimitGroupUnion) AsAny() anyWorkspaceRateLimitGroup {
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

func (u WorkspaceRateLimitGroupUnion) AsModelGroup() (v OrganizationRateLimitModelGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u WorkspaceRateLimitGroupUnion) AsBatch() (v OrganizationRateLimitBatchGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u WorkspaceRateLimitGroupUnion) AsTokenCount() (v OrganizationRateLimitTokenCountGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u WorkspaceRateLimitGroupUnion) AsFiles() (v OrganizationRateLimitFilesGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u WorkspaceRateLimitGroupUnion) AsSkills() (v OrganizationRateLimitSkillsGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u WorkspaceRateLimitGroupUnion) AsWebSearch() (v OrganizationRateLimitWebSearchGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u WorkspaceRateLimitGroupUnion) RawJSON() string { return u.JSON.raw }

func (r *WorkspaceRateLimitGroupUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkspaceRateLimitOrganizationSource struct {
	// Always `organization`: no workspace-level override is stored, so the
	// organization's value applies.
	Type constant.Organization `json:"type" default:"organization"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkspaceRateLimitOrganizationSource) RawJSON() string { return r.JSON.raw }
func (r *WorkspaceRateLimitOrganizationSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkspaceRateLimitValue struct {
	// The organization-level value for the same limiter type, for reference. `null`
	// when the organization has no limit configured for this limiter type.
	OrgLimit int64 `json:"org_limit" api:"required"`
	// Where `value` comes from. `organization` values are listed only when
	// `include_inherited` is `true`, and then `value` equals `org_limit`.
	Source WorkspaceRateLimitValueSourceUnion `json:"source" api:"required"`
	// The limiter type (for example, `requests_per_minute` or
	// `input_tokens_per_minute`).
	Type string `json:"type" api:"required"`
	// The workspace's value for this limiter type: the workspace-level override when
	// `source.type` is `workspace`, otherwise the organization's value.
	Value int64 `json:"value" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		OrgLimit    respjson.Field
		Source      respjson.Field
		Type        respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkspaceRateLimitValue) RawJSON() string { return r.JSON.raw }
func (r *WorkspaceRateLimitValue) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// WorkspaceRateLimitValueSourceUnion contains all possible properties and values
// from [WorkspaceRateLimitWorkspaceSource],
// [WorkspaceRateLimitOrganizationSource].
//
// Use the [WorkspaceRateLimitValueSourceUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type WorkspaceRateLimitValueSourceUnion struct {
	// Any of "workspace", "organization".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyWorkspaceRateLimitValueSource is implemented by each variant of
// [WorkspaceRateLimitValueSourceUnion] to add type safety for the return type of
// [WorkspaceRateLimitValueSourceUnion.AsAny]
type anyWorkspaceRateLimitValueSource interface {
	implWorkspaceRateLimitValueSourceUnion()
}

func (WorkspaceRateLimitWorkspaceSource) implWorkspaceRateLimitValueSourceUnion()    {}
func (WorkspaceRateLimitOrganizationSource) implWorkspaceRateLimitValueSourceUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := WorkspaceRateLimitValueSourceUnion.AsAny().(type) {
//	case anthropic.WorkspaceRateLimitWorkspaceSource:
//	case anthropic.WorkspaceRateLimitOrganizationSource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u WorkspaceRateLimitValueSourceUnion) AsAny() anyWorkspaceRateLimitValueSource {
	switch u.Type {
	case "workspace":
		return u.AsWorkspace()
	case "organization":
		return u.AsOrganization()
	}
	return nil
}

func (u WorkspaceRateLimitValueSourceUnion) AsWorkspace() (v WorkspaceRateLimitWorkspaceSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u WorkspaceRateLimitValueSourceUnion) AsOrganization() (v WorkspaceRateLimitOrganizationSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u WorkspaceRateLimitValueSourceUnion) RawJSON() string { return u.JSON.raw }

func (r *WorkspaceRateLimitValueSourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkspaceRateLimitWorkspaceSource struct {
	// Always `workspace`: a workspace-level override is stored.
	Type constant.Workspace `json:"type" default:"workspace"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r WorkspaceRateLimitWorkspaceSource) RawJSON() string { return r.JSON.raw }
func (r *WorkspaceRateLimitWorkspaceSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationWorkspaceRateLimitListParams struct {
	// Maximum number of items to return per page. Ranges from `1` to `1000`.
	//
	// When omitted, every remaining entry is returned in a single page and `next_page`
	// is `null`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Also list the limiter values the workspace inherits from the organization,
	// including groups with no workspace-level override.
	IncludeInherited param.Opt[bool] `query:"include_inherited,omitzero" json:"-"`
	// Filter by group type.
	//
	// Any of "batch", "files", "model_group", "skills", "token_count", "web_search".
	GroupType OrganizationWorkspaceRateLimitListParamsGroupType `query:"group_type,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrganizationWorkspaceRateLimitListParams]'s query
// parameters as `url.Values`.
func (r OrganizationWorkspaceRateLimitListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Filter by group type.
type OrganizationWorkspaceRateLimitListParamsGroupType string

const (
	OrganizationWorkspaceRateLimitListParamsGroupTypeBatch      OrganizationWorkspaceRateLimitListParamsGroupType = "batch"
	OrganizationWorkspaceRateLimitListParamsGroupTypeFiles      OrganizationWorkspaceRateLimitListParamsGroupType = "files"
	OrganizationWorkspaceRateLimitListParamsGroupTypeModelGroup OrganizationWorkspaceRateLimitListParamsGroupType = "model_group"
	OrganizationWorkspaceRateLimitListParamsGroupTypeSkills     OrganizationWorkspaceRateLimitListParamsGroupType = "skills"
	OrganizationWorkspaceRateLimitListParamsGroupTypeTokenCount OrganizationWorkspaceRateLimitListParamsGroupType = "token_count"
	OrganizationWorkspaceRateLimitListParamsGroupTypeWebSearch  OrganizationWorkspaceRateLimitListParamsGroupType = "web_search"
)
