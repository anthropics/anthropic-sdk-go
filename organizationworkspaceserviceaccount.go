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

// OrganizationWorkspaceServiceAccountService contains methods and other services
// that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationWorkspaceServiceAccountService] method instead.
type OrganizationWorkspaceServiceAccountService struct {
	Options []option.RequestOption
}

// NewOrganizationWorkspaceServiceAccountService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewOrganizationWorkspaceServiceAccountService(opts ...option.RequestOption) (r OrganizationWorkspaceServiceAccountService) {
	r = OrganizationWorkspaceServiceAccountService{}
	r.Options = opts
	return
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Retrieve a service account's membership in a workspace.
//
// Returns the membership record, including the service account's `workspace_role`
// in this workspace. Archived workspaces return 400. For the default workspace,
// returns the implicit (`implicit: true`) membership when no explicit membership
// exists; an explicitly added membership is returned with its assigned role. An
// archived service account returns 404.
func (r *OrganizationWorkspaceServiceAccountService) Get(ctx context.Context, serviceAccountID string, query OrganizationWorkspaceServiceAccountGetParams, opts ...option.RequestOption) (res *ServiceAccountWorkspaceMember, err error) {
	opts = slices.Concat(r.Options, opts)
	if query.WorkspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	if serviceAccountID == "" {
		err = errors.New("missing required service_account_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s/service_accounts/%s", url.PathEscape(query.WorkspaceID), url.PathEscape(serviceAccountID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Change a service account's role in a workspace.
//
// The new `workspace_role` replaces the current one. Only explicit memberships can
// be updated; to set a role on the implicit default-workspace membership, add the
// service account explicitly with
// `POST /workspaces/{workspace_id}/service_accounts`. Archived workspaces
// return 400. Archived service accounts cannot be updated and are rejected.
func (r *OrganizationWorkspaceServiceAccountService) Update(ctx context.Context, serviceAccountID string, params OrganizationWorkspaceServiceAccountUpdateParams, opts ...option.RequestOption) (res *ServiceAccountWorkspaceMember, err error) {
	opts = slices.Concat(r.Options, opts)
	if params.WorkspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	if serviceAccountID == "" {
		err = errors.New("missing required service_account_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s/service_accounts/%s", url.PathEscape(params.WorkspaceID), url.PathEscape(serviceAccountID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// List the service accounts that are members of a workspace.
//
// Each entry includes the service account's `workspace_role`. Use `limit` and the
// `next_page` cursor to paginate. Archived workspaces return 400; use
// `GET /service_accounts/{id}/workspaces` to audit memberships of an archived
// workspace. The implicit default-workspace membership is not included in this
// list. Memberships of archived service accounts are omitted from the results.
func (r *OrganizationWorkspaceServiceAccountService) List(ctx context.Context, workspaceID string, query OrganizationWorkspaceServiceAccountListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ServiceAccountWorkspaceMember], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s/service_accounts", url.PathEscape(workspaceID))
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

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// List the service accounts that are members of a workspace.
//
// Each entry includes the service account's `workspace_role`. Use `limit` and the
// `next_page` cursor to paginate. Archived workspaces return 400; use
// `GET /service_accounts/{id}/workspaces` to audit memberships of an archived
// workspace. The implicit default-workspace membership is not included in this
// list. Memberships of archived service accounts are omitted from the results.
func (r *OrganizationWorkspaceServiceAccountService) ListAutoPaging(ctx context.Context, workspaceID string, query OrganizationWorkspaceServiceAccountListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ServiceAccountWorkspaceMember] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, workspaceID, query, opts...))
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Add a service account to a workspace with the given `workspace_role`.
//
// The role determines what the service account can do in the workspace and which
// workspace-scoped permissions it can be granted when authenticating through
// federation. Every service account is already an implicit `workspace_user` member
// of the default workspace; adding it explicitly assigns a chosen role. If the
// service account is already an explicit member of the workspace, its
// `workspace_role` is replaced with the value supplied here. Archived workspaces
// return 400. Archived service accounts cannot be added and are rejected.
func (r *OrganizationWorkspaceServiceAccountService) Add(ctx context.Context, workspaceID string, body OrganizationWorkspaceServiceAccountAddParams, opts ...option.RequestOption) (res *ServiceAccountWorkspaceMember, err error) {
	opts = slices.Concat(r.Options, opts)
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s/service_accounts", url.PathEscape(workspaceID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Remove a service account from a workspace.
//
// Removal is idempotent (returns 200 even if the membership was already removed).
// A DELETE against the implicit default-workspace membership returns 200 but is a
// no-op and the membership persists; deleting an explicit default-workspace row
// reverts to the implicit `workspace_user` membership. Archived workspaces
// return 400.
func (r *OrganizationWorkspaceServiceAccountService) Remove(ctx context.Context, serviceAccountID string, body OrganizationWorkspaceServiceAccountRemoveParams, opts ...option.RequestOption) (res *OrganizationWorkspaceServiceAccountRemoveResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if body.WorkspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	if serviceAccountID == "" {
		err = errors.New("missing required service_account_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s/service_accounts/%s", url.PathEscape(body.WorkspaceID), url.PathEscape(serviceAccountID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type OrganizationWorkspaceServiceAccountRemoveResponse struct {
	// Tagged service account ID (`svac_...`) named in the delete request. Removal is
	// idempotent; see the endpoint description for the implicit-membership no-op.
	ServiceAccountID string                                        `json:"service_account_id" api:"required"`
	Type             constant.ServiceAccountWorkspaceMemberDeleted `json:"type" default:"service_account_workspace_member_deleted"`
	// Tagged workspace ID (`wrkspc_...`) named in the delete request.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ServiceAccountID respjson.Field
		Type             respjson.Field
		WorkspaceID      respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationWorkspaceServiceAccountRemoveResponse) RawJSON() string { return r.JSON.raw }
func (r *OrganizationWorkspaceServiceAccountRemoveResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationWorkspaceServiceAccountGetParams struct {
	// ID of the workspace.
	WorkspaceID string `path:"workspace_id" api:"required" json:"-"`
	paramObj
}

type OrganizationWorkspaceServiceAccountUpdateParams struct {
	// ID of the workspace.
	WorkspaceID string `path:"workspace_id" api:"required" json:"-"`
	// New role for the service account in this workspace.
	//
	// Any of "workspace_admin", "workspace_developer",
	// "workspace_restricted_developer", "workspace_user".
	WorkspaceRole NoBillingWorkspaceRole `json:"workspace_role,omitzero" api:"required"`
	paramObj
}

func (r OrganizationWorkspaceServiceAccountUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationWorkspaceServiceAccountUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationWorkspaceServiceAccountUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationWorkspaceServiceAccountListParams struct {
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of results per page.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrganizationWorkspaceServiceAccountListParams]'s query
// parameters as `url.Values`.
func (r OrganizationWorkspaceServiceAccountListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type OrganizationWorkspaceServiceAccountAddParams struct {
	// Tagged service account ID to add.
	ServiceAccountID string `json:"service_account_id" api:"required"`
	// Role to assign to the service account in this workspace.
	//
	// Any of "workspace_admin", "workspace_developer",
	// "workspace_restricted_developer", "workspace_user".
	WorkspaceRole NoBillingWorkspaceRole `json:"workspace_role,omitzero" api:"required"`
	paramObj
}

func (r OrganizationWorkspaceServiceAccountAddParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationWorkspaceServiceAccountAddParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationWorkspaceServiceAccountAddParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationWorkspaceServiceAccountRemoveParams struct {
	// ID of the workspace.
	WorkspaceID string `path:"workspace_id" api:"required" json:"-"`
	paramObj
}
