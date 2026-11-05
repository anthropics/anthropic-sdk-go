package anthropic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// OrganizationRBACGroupMemberService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationRBACGroupMemberService] method instead.
type OrganizationRBACGroupMemberService struct {
	Options []option.RequestOption
}

// NewOrganizationRBACGroupMemberService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewOrganizationRBACGroupMemberService(opts ...option.RequestOption) (r OrganizationRBACGroupMemberService) {
	r = OrganizationRBACGroupMemberService{}
	r.Options = opts
	return
}

// List members of an RBAC Group.
//
// The RBAC Groups API is available to Claude Enterprise organizations only.
func (r *OrganizationRBACGroupMemberService) List(ctx context.Context, rbacGroupID string, query OrganizationRBACGroupMemberListParams, opts ...option.RequestOption) (res *pagination.PageCursor[RBACGroupMember], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if rbacGroupID == "" {
		err = errors.New("missing required rbac_group_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/rbac_groups/%s/members", requestconfig.EscapePathSegment(rbacGroupID))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.rbac_groups.members.list", "/v1/organizations/rbac_groups/{rbac_group_id}/members", "limit", "page"), requestconfig.WithSpanAttribute("anthropic.organization.rbac_group.id", rbacGroupID))
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
func (r *OrganizationRBACGroupMemberService) ListAutoPaging(ctx context.Context, rbacGroupID string, query OrganizationRBACGroupMemberListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[RBACGroupMember] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, rbacGroupID, query, opts...))
}

// Add a User to an RBAC Group. Membership of groups provisioned by an identity
// provider (source type `"scim"`) cannot be modified via the API while an
// organization in the tenant uses SCIM provisioning.
//
// The RBAC Groups API is available to Claude Enterprise organizations only.
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeRBACGroupGrantsOwnerAccess] (HTTP 400)
//   - [ErrorCodeRBACGroupOwnerAccessUnconfirmed] (HTTP 400)
func (r *OrganizationRBACGroupMemberService) Add(ctx context.Context, rbacGroupID string, body OrganizationRBACGroupMemberAddParams, opts ...option.RequestOption) (res *RBACGroupMember, err error) {
	opts = slices.Concat(r.Options, opts)
	if rbacGroupID == "" {
		err = errors.New("missing required rbac_group_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/rbac_groups/%s/members", requestconfig.EscapePathSegment(rbacGroupID))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.rbac_groups.members.add", "/v1/organizations/rbac_groups/{rbac_group_id}/members"), requestconfig.WithSpanAttribute("anthropic.organization.rbac_group.id", rbacGroupID), requestconfig.WithSpanCreatedID("anthropic.organization.rbac_group.member.id", &res, func(created *RBACGroupMember) string { return created.UserID }))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Remove a User from an RBAC Group. Membership of groups provisioned by an
// identity provider (source type `"scim"`) cannot be modified via the API while an
// organization in the tenant uses SCIM provisioning.
//
// The RBAC Groups API is available to Claude Enterprise organizations only.
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeRBACGroupGrantsOwnerAccess] (HTTP 400)
//   - [ErrorCodeRBACGroupOwnerAccessUnconfirmed] (HTTP 400)
func (r *OrganizationRBACGroupMemberService) Remove(ctx context.Context, userID string, body OrganizationRBACGroupMemberRemoveParams, opts ...option.RequestOption) (res *DeletedRBACGroupMember, err error) {
	opts = slices.Concat(r.Options, opts)
	if body.RBACGroupID == "" {
		err = errors.New("missing required rbac_group_id parameter")
		return nil, err
	}
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/rbac_groups/%s/members/%s", requestconfig.EscapePathSegment(body.RBACGroupID), requestconfig.EscapePathSegment(userID))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.rbac_groups.members.remove", "/v1/organizations/rbac_groups/{rbac_group_id}/members/{user_id}"), requestconfig.WithSpanAttribute("anthropic.organization.rbac_group.id", body.RBACGroupID), requestconfig.WithSpanAttribute("anthropic.organization.rbac_group.member.id", userID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type DeletedRBACGroupMember struct {
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
func (r DeletedRBACGroupMember) RawJSON() string { return r.JSON.raw }
func (r *DeletedRBACGroupMember) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type RBACGroupMember struct {
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
func (r RBACGroupMember) RawJSON() string { return r.JSON.raw }
func (r *RBACGroupMember) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationRBACGroupMemberListParams struct {
	// Optionally set to the `next_page` token from the previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `1000`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrganizationRBACGroupMemberListParams]'s query parameters
// as `url.Values`.
func (r OrganizationRBACGroupMemberListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type OrganizationRBACGroupMemberAddParams struct {
	// ID of the User.
	UserID string `json:"user_id" api:"required"`
	paramObj
}

func (r OrganizationRBACGroupMemberAddParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationRBACGroupMemberAddParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationRBACGroupMemberAddParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationRBACGroupMemberRemoveParams struct {
	// ID of the RBAC Group.
	RBACGroupID string `path:"rbac_group_id" api:"required" json:"-"`
	paramObj
}
