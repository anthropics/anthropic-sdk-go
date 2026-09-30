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

// OrganizationWorkspaceMemberService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationWorkspaceMemberService] method instead.
type OrganizationWorkspaceMemberService struct {
	Options []option.RequestOption
}

// NewOrganizationWorkspaceMemberService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewOrganizationWorkspaceMemberService(opts ...option.RequestOption) (r OrganizationWorkspaceMemberService) {
	r = OrganizationWorkspaceMemberService{}
	r.Options = opts
	return
}

// Get Workspace Member
func (r *OrganizationWorkspaceMemberService) Get(ctx context.Context, userID string, query OrganizationWorkspaceMemberGetParams, opts ...option.RequestOption) (res *WorkspaceMember, err error) {
	opts = slices.Concat(r.Options, opts)
	if query.WorkspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s/members/%s", query.WorkspaceID, userID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Workspace Member
func (r *OrganizationWorkspaceMemberService) Update(ctx context.Context, userID string, params OrganizationWorkspaceMemberUpdateParams, opts ...option.RequestOption) (res *WorkspaceMember, err error) {
	opts = slices.Concat(r.Options, opts)
	if params.WorkspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s/members/%s", params.WorkspaceID, userID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List Workspace Members
func (r *OrganizationWorkspaceMemberService) List(ctx context.Context, workspaceID string, query OrganizationWorkspaceMemberListParams, opts ...option.RequestOption) (res *pagination.Page[WorkspaceMember], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s/members", workspaceID)
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

// List Workspace Members
func (r *OrganizationWorkspaceMemberService) ListAutoPaging(ctx context.Context, workspaceID string, query OrganizationWorkspaceMemberListParams, opts ...option.RequestOption) *pagination.PageAutoPager[WorkspaceMember] {
	return pagination.NewPageAutoPager(r.List(ctx, workspaceID, query, opts...))
}

// Create Workspace Member
func (r *OrganizationWorkspaceMemberService) Add(ctx context.Context, workspaceID string, body OrganizationWorkspaceMemberAddParams, opts ...option.RequestOption) (res *WorkspaceMember, err error) {
	opts = slices.Concat(r.Options, opts)
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s/members", workspaceID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Delete Workspace Member
func (r *OrganizationWorkspaceMemberService) Remove(ctx context.Context, userID string, body OrganizationWorkspaceMemberRemoveParams, opts ...option.RequestOption) (res *OrganizationWorkspaceMemberRemoveResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if body.WorkspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	if userID == "" {
		err = errors.New("missing required user_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s/members/%s", body.WorkspaceID, userID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type OrganizationWorkspaceMemberRemoveResponse struct {
	// Deleted object type.
	//
	// For Workspace Members, this is always `"workspace_member_deleted"`.
	Type constant.WorkspaceMemberDeleted `json:"type" default:"workspace_member_deleted"`
	// ID of the User.
	UserID string `json:"user_id" api:"required"`
	// ID of the Workspace.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		UserID      respjson.Field
		WorkspaceID respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationWorkspaceMemberRemoveResponse) RawJSON() string { return r.JSON.raw }
func (r *OrganizationWorkspaceMemberRemoveResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationWorkspaceMemberGetParams struct {
	// ID of the Workspace.
	WorkspaceID string `path:"workspace_id" api:"required" json:"-"`
	paramObj
}

type OrganizationWorkspaceMemberUpdateParams struct {
	// ID of the Workspace.
	WorkspaceID string `path:"workspace_id" api:"required" json:"-"`
	// New workspace role for the User.
	//
	// Any of "workspace_admin", "workspace_billing", "workspace_developer",
	// "workspace_restricted_developer", "workspace_user".
	WorkspaceRole WorkspaceRole `json:"workspace_role,omitzero" api:"required"`
	paramObj
}

func (r OrganizationWorkspaceMemberUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationWorkspaceMemberUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationWorkspaceMemberUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationWorkspaceMemberListParams struct {
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
	paramObj
}

// URLQuery serializes [OrganizationWorkspaceMemberListParams]'s query parameters
// as `url.Values`.
func (r OrganizationWorkspaceMemberListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type OrganizationWorkspaceMemberAddParams struct {
	// ID of the User.
	UserID string `json:"user_id" api:"required"`
	// Role of the new Workspace Member. Cannot be `workspace_billing`.
	//
	// Any of "workspace_admin", "workspace_developer",
	// "workspace_restricted_developer", "workspace_user".
	WorkspaceRole NoBillingWorkspaceRole `json:"workspace_role,omitzero" api:"required"`
	paramObj
}

func (r OrganizationWorkspaceMemberAddParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationWorkspaceMemberAddParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationWorkspaceMemberAddParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationWorkspaceMemberRemoveParams struct {
	// ID of the Workspace.
	WorkspaceID string `path:"workspace_id" api:"required" json:"-"`
	paramObj
}
