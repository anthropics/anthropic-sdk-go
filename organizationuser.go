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

// OrganizationUserService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationUserService] method instead.
type OrganizationUserService struct {
	Options []option.RequestOption
}

// NewOrganizationUserService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewOrganizationUserService(opts ...option.RequestOption) (r OrganizationUserService) {
	r = OrganizationUserService{}
	r.Options = opts
	return
}

// Retrieve a member of the organization by user ID.
func (r *OrganizationUserService) Get(ctx context.Context, userID string, opts ...option.RequestOption) (res *OrganizationUser, err error) {
	opts = slices.Concat(r.Options, opts)
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/users/%s", url.PathEscape(userID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a member's organization role.
func (r *OrganizationUserService) Update(ctx context.Context, userID string, body OrganizationUserUpdateParams, opts ...option.RequestOption) (res *OrganizationUser, err error) {
	opts = slices.Concat(r.Options, opts)
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/users/%s", url.PathEscape(userID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// List the organization's members.
func (r *OrganizationUserService) List(ctx context.Context, query OrganizationUserListParams, opts ...option.RequestOption) (res *pagination.Page[OrganizationUser], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/users"
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

// List the organization's members.
func (r *OrganizationUserService) ListAutoPaging(ctx context.Context, query OrganizationUserListParams, opts ...option.RequestOption) *pagination.PageAutoPager[OrganizationUser] {
	return pagination.NewPageAutoPager(r.List(ctx, query, opts...))
}

// Remove a member from the organization.
func (r *OrganizationUserService) Remove(ctx context.Context, userID string, opts ...option.RequestOption) (res *OrganizationUserRemoveResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/users/%s", url.PathEscape(userID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type OrganizationUser struct {
	// ID of the User.
	ID string `json:"id" api:"required"`
	// RFC 3339 datetime string indicating when the User joined the Organization.
	AddedAt time.Time `json:"added_at" api:"required" format:"date-time"`
	// Email of the User.
	Email string `json:"email" api:"required"`
	// Name of the User.
	Name string `json:"name" api:"required"`
	// Organization role of the User.
	//
	// Any of "admin", "billing", "claude_code_user", "developer", "managed",
	// "membership_admin", "owner", "primary_owner", "user".
	Role OrganizationRole `json:"role" api:"required"`
	// Object type.
	//
	// For Users, this is always `"user"`.
	Type constant.User `json:"type" default:"user"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		AddedAt     respjson.Field
		Email       respjson.Field
		Name        respjson.Field
		Role        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationUser) RawJSON() string { return r.JSON.raw }
func (r *OrganizationUser) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationUserRemoveResponse struct {
	// ID of the User.
	ID string `json:"id" api:"required"`
	// Deleted object type.
	//
	// For Users, this is always `"user_deleted"`.
	Type constant.UserDeleted `json:"type" default:"user_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationUserRemoveResponse) RawJSON() string { return r.JSON.raw }
func (r *OrganizationUserRemoveResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationUserUpdateParams struct {
	// New role for the User.
	//
	// The accepted values depend on the organization type. Console and API
	// organizations accept `user`, `developer`, `billing`, and `claude_code_user`;
	// `admin` cannot be assigned through the API. Claude Enterprise organizations
	// accept `user` and `managed`.
	//
	// Any of "billing", "claude_code_user", "developer", "managed", "user".
	Role OrganizationUserUpdateParamsRole `json:"role,omitzero" api:"required"`
	paramObj
}

func (r OrganizationUserUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationUserUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationUserUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// New role for the User.
//
// The accepted values depend on the organization type. Console and API
// organizations accept `user`, `developer`, `billing`, and `claude_code_user`;
// `admin` cannot be assigned through the API. Claude Enterprise organizations
// accept `user` and `managed`.
type OrganizationUserUpdateParamsRole string

const (
	OrganizationUserUpdateParamsRoleBilling        OrganizationUserUpdateParamsRole = "billing"
	OrganizationUserUpdateParamsRoleClaudeCodeUser OrganizationUserUpdateParamsRole = "claude_code_user"
	OrganizationUserUpdateParamsRoleDeveloper      OrganizationUserUpdateParamsRole = "developer"
	OrganizationUserUpdateParamsRoleManaged        OrganizationUserUpdateParamsRole = "managed"
	OrganizationUserUpdateParamsRoleUser           OrganizationUserUpdateParamsRole = "user"
)

type OrganizationUserListParams struct {
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately after this object.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately before this object.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// Filter by user email.
	Email param.Opt[string] `query:"email,omitzero" format:"email" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `1000`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Filter to items whose `role` equals one of the supplied values. Repeatable;
	// values are OR'ed together.
	//
	// Accepted values depend on the organization type: Console and API organizations
	// accept `user`, `developer`, `billing`, `admin`, and `claude_code_user`; Claude
	// Enterprise organizations accept `user`, `owner`, `primary_owner`,
	// `membership_admin`, and `managed`.
	Roles []string `query:"roles,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrganizationUserListParams]'s query parameters as
// `url.Values`.
func (r OrganizationUserListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
