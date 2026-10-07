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

// BetaOrganizationRBACGroupService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationRBACGroupService] method instead.
type BetaOrganizationRBACGroupService struct {
	Options []option.RequestOption
	Members BetaOrganizationRBACGroupMemberService
}

// NewBetaOrganizationRBACGroupService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationRBACGroupService(opts ...option.RequestOption) (r BetaOrganizationRBACGroupService) {
	r = BetaOrganizationRBACGroupService{}
	r.Options = opts
	r.Members = NewBetaOrganizationRBACGroupMemberService(opts...)
	return
}

// Create an RBAC Group in the Claude Enterprise tenant. Groups created via the API
// have source type `"direct"`.
//
// The RBAC Groups API is available to Claude Enterprise organizations only.
func (r *BetaOrganizationRBACGroupService) New(ctx context.Context, body BetaOrganizationRBACGroupNewParams, opts ...option.RequestOption) (res *BetaRBACGroup, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/rbac_groups?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve an RBAC Group by ID.
//
// The RBAC Groups API is available to Claude Enterprise organizations only.
func (r *BetaOrganizationRBACGroupService) Get(ctx context.Context, rbacGroupID string, opts ...option.RequestOption) (res *BetaRBACGroup, err error) {
	opts = slices.Concat(r.Options, opts)
	if rbacGroupID == "" {
		err = errors.New("missing required rbac_group_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/rbac_groups/%s?beta=true", url.PathEscape(rbacGroupID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update an RBAC Group's name. Groups provisioned by an identity provider (source
// type `"scim"`) cannot be modified via the API while an organization in the
// tenant uses SCIM provisioning.
//
// The RBAC Groups API is available to Claude Enterprise organizations only.
func (r *BetaOrganizationRBACGroupService) Update(ctx context.Context, rbacGroupID string, body BetaOrganizationRBACGroupUpdateParams, opts ...option.RequestOption) (res *BetaRBACGroup, err error) {
	opts = slices.Concat(r.Options, opts)
	if rbacGroupID == "" {
		err = errors.New("missing required rbac_group_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/rbac_groups/%s?beta=true", url.PathEscape(rbacGroupID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// List RBAC Groups in the Claude Enterprise tenant.
//
// The RBAC Groups API is available to Claude Enterprise organizations only.
func (r *BetaOrganizationRBACGroupService) List(ctx context.Context, query BetaOrganizationRBACGroupListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaRBACGroup], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/rbac_groups?beta=true"
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

// List RBAC Groups in the Claude Enterprise tenant.
//
// The RBAC Groups API is available to Claude Enterprise organizations only.
func (r *BetaOrganizationRBACGroupService) ListAutoPaging(ctx context.Context, query BetaOrganizationRBACGroupListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaRBACGroup] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

// Delete an RBAC Group. Groups provisioned by an identity provider (source type
// `"scim"`) cannot be deleted via the API while an organization in the tenant uses
// SCIM provisioning.
//
// The RBAC Groups API is available to Claude Enterprise organizations only.
func (r *BetaOrganizationRBACGroupService) Delete(ctx context.Context, rbacGroupID string, opts ...option.RequestOption) (res *BetaOrganizationRBACGroupDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if rbacGroupID == "" {
		err = errors.New("missing required rbac_group_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/rbac_groups/%s?beta=true", url.PathEscape(rbacGroupID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type BetaRBACGroup struct {
	// ID of the RBAC Group.
	ID string `json:"id" api:"required"`
	// RFC 3339 timestamp of when the RBAC Group was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Name of the RBAC Group. Not uniqueness-enforced.
	Name string `json:"name" api:"required"`
	// RBAC Role IDs attached to this RBAC Group. Role attachment is managed in the
	// admin settings and is read-only on this API. `null` means role data was
	// temporarily unavailable — retry to distinguish from an empty list.
	RoleIDs []string `json:"role_ids" api:"required"`
	// How the RBAC Group was created: `"direct"` for groups created directly (for
	// example, in the organization's admin settings), `"scim"` for groups provisioned
	// by the identity provider.
	//
	// Any of "direct", "scim".
	SourceType BetaRBACGroupSourceType `json:"source_type" api:"required"`
	// Object type.
	//
	// For RBAC Groups, this is always `"rbac_group"`.
	Type constant.RBACGroup `json:"type" default:"rbac_group"`
	// RFC 3339 timestamp of when the RBAC Group was last updated.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		RoleIDs     respjson.Field
		SourceType  respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaRBACGroup) RawJSON() string { return r.JSON.raw }
func (r *BetaRBACGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// How the RBAC Group was created: `"direct"` for groups created directly (for
// example, in the organization's admin settings), `"scim"` for groups provisioned
// by the identity provider.
type BetaRBACGroupSourceType string

const (
	BetaRBACGroupSourceTypeDirect BetaRBACGroupSourceType = "direct"
	BetaRBACGroupSourceTypeSCIM   BetaRBACGroupSourceType = "scim"
)

type BetaOrganizationRBACGroupDeleteResponse struct {
	// ID of the RBAC Group.
	ID string `json:"id" api:"required"`
	// Deleted object type.
	//
	// For RBAC Groups, this is always `"rbac_group_deleted"`.
	Type constant.RBACGroupDeleted `json:"type" default:"rbac_group_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationRBACGroupDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationRBACGroupDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationRBACGroupNewParams struct {
	// Name of the RBAC Group. Not uniqueness-enforced.
	Name string `json:"name" api:"required"`
	paramObj
}

func (r BetaOrganizationRBACGroupNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationRBACGroupNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationRBACGroupNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationRBACGroupUpdateParams struct {
	// Name of the RBAC Group. Not uniqueness-enforced.
	Name param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r BetaOrganizationRBACGroupUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationRBACGroupUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationRBACGroupUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationRBACGroupListParams struct {
	// Optionally set to the `next_page` token from the previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `1000`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationRBACGroupListParams]'s query parameters as
// `url.Values`.
func (r BetaOrganizationRBACGroupListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
