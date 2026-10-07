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

// BetaOrganizationRBACGroupMemberService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationRBACGroupMemberService] method instead.
type BetaOrganizationRBACGroupMemberService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationRBACGroupMemberService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationRBACGroupMemberService(opts ...option.RequestOption) (r BetaOrganizationRBACGroupMemberService) {
	r = BetaOrganizationRBACGroupMemberService{}
	r.Options = opts
	return
}

// List members of an RBAC Group.
//
// The RBAC Groups API is available to Claude Enterprise organizations only.
func (r *BetaOrganizationRBACGroupMemberService) List(ctx context.Context, rbacGroupID string, query BetaOrganizationRBACGroupMemberListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaRBACGroupMember], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if rbacGroupID == "" {
		err = errors.New("missing required rbac_group_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/rbac_groups/%s/members?beta=true", url.PathEscape(rbacGroupID))
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

// List members of an RBAC Group.
//
// The RBAC Groups API is available to Claude Enterprise organizations only.
func (r *BetaOrganizationRBACGroupMemberService) ListAutoPaging(ctx context.Context, rbacGroupID string, query BetaOrganizationRBACGroupMemberListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaRBACGroupMember] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, rbacGroupID, query, opts...))
}

// Add a User to an RBAC Group. Membership of groups provisioned by an identity
// provider (source type `"scim"`) cannot be modified via the API while an
// organization in the tenant uses SCIM provisioning.
//
// The RBAC Groups API is available to Claude Enterprise organizations only.
func (r *BetaOrganizationRBACGroupMemberService) Add(ctx context.Context, rbacGroupID string, body BetaOrganizationRBACGroupMemberAddParams, opts ...option.RequestOption) (res *BetaRBACGroupMember, err error) {
	opts = slices.Concat(r.Options, opts)
	if rbacGroupID == "" {
		err = errors.New("missing required rbac_group_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/rbac_groups/%s/members?beta=true", url.PathEscape(rbacGroupID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Remove a User from an RBAC Group. Membership of groups provisioned by an
// identity provider (source type `"scim"`) cannot be modified via the API while an
// organization in the tenant uses SCIM provisioning.
//
// The RBAC Groups API is available to Claude Enterprise organizations only.
func (r *BetaOrganizationRBACGroupMemberService) Remove(ctx context.Context, userID string, body BetaOrganizationRBACGroupMemberRemoveParams, opts ...option.RequestOption) (res *BetaOrganizationRBACGroupMemberRemoveResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if body.RBACGroupID == "" {
		err = errors.New("missing required rbac_group_id parameter")
		return nil, err
	}
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/rbac_groups/%s/members/%s?beta=true", url.PathEscape(body.RBACGroupID), url.PathEscape(userID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type BetaRBACGroupMember struct {
	// RFC 3339 timestamp of when the User was added to the RBAC Group.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Email of the User.
	Email string `json:"email" api:"required"`
	// ID of the RBAC Group.
	RBACGroupID string `json:"rbac_group_id" api:"required"`
	// Object type.
	//
	// For RBAC Group Members, this is always `"rbac_group_member"`.
	Type constant.RBACGroupMember `json:"type" default:"rbac_group_member"`
	// ID of the User.
	UserID string `json:"user_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CreatedAt   respjson.Field
		Email       respjson.Field
		RBACGroupID respjson.Field
		Type        respjson.Field
		UserID      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaRBACGroupMember) RawJSON() string { return r.JSON.raw }
func (r *BetaRBACGroupMember) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationRBACGroupMemberRemoveResponse struct {
	// ID of the RBAC Group.
	RBACGroupID string `json:"rbac_group_id" api:"required"`
	// Deleted object type. For RBAC Group Members, this is always
	// `"rbac_group_member_deleted"`.
	Type constant.RBACGroupMemberDeleted `json:"type" default:"rbac_group_member_deleted"`
	// ID of the User.
	UserID string `json:"user_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		RBACGroupID respjson.Field
		Type        respjson.Field
		UserID      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationRBACGroupMemberRemoveResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationRBACGroupMemberRemoveResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationRBACGroupMemberListParams struct {
	// Optionally set to the `next_page` token from the previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `1000`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationRBACGroupMemberListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationRBACGroupMemberListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationRBACGroupMemberAddParams struct {
	// ID of the User.
	UserID string `json:"user_id" api:"required"`
	paramObj
}

func (r BetaOrganizationRBACGroupMemberAddParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationRBACGroupMemberAddParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationRBACGroupMemberAddParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationRBACGroupMemberRemoveParams struct {
	// ID of the RBAC Group.
	RBACGroupID string `path:"rbac_group_id" api:"required" json:"-"`
	paramObj
}
