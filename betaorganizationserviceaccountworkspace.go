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

// BetaOrganizationServiceAccountWorkspaceService contains methods and other
// services that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationServiceAccountWorkspaceService] method instead.
type BetaOrganizationServiceAccountWorkspaceService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationServiceAccountWorkspaceService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationServiceAccountWorkspaceService(opts ...option.RequestOption) (r BetaOrganizationServiceAccountWorkspaceService) {
	r = BetaOrganizationServiceAccountWorkspaceService{}
	r.Options = opts
	return
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// List the workspaces a service account is a member of.
//
// Each entry includes the service account's `workspace_role` in that workspace.
// Use `limit` and the `next_page` cursor to paginate. When the service account has
// no explicit default-workspace membership, the implicit (`implicit: true`)
// membership is returned as the first entry on the first page; with `limit=1` the
// first page may return up to 2 entries (the implicit entry plus one explicit
// membership) so a pagination cursor can be derived. Memberships are returned only
// while the service account is active. Without a `page` cursor, an archived
// service account returns an empty list. A `page` cursor that does not match an
// active membership returns a 400 invalid-request error. A cursor stops matching
// when the membership is removed, the workspace is deleted, or the service account
// is archived. Restart pagination from the first page to recover.
func (r *BetaOrganizationServiceAccountWorkspaceService) List(ctx context.Context, serviceAccountID string, params BetaOrganizationServiceAccountWorkspaceListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaServiceAccountWorkspaceMember], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if serviceAccountID == "" {
		err = errors.New("missing required service_account_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/service_accounts/%s/workspaces?beta=true", serviceAccountID)
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, params, &res, opts...)
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
// List the workspaces a service account is a member of.
//
// Each entry includes the service account's `workspace_role` in that workspace.
// Use `limit` and the `next_page` cursor to paginate. When the service account has
// no explicit default-workspace membership, the implicit (`implicit: true`)
// membership is returned as the first entry on the first page; with `limit=1` the
// first page may return up to 2 entries (the implicit entry plus one explicit
// membership) so a pagination cursor can be derived. Memberships are returned only
// while the service account is active. Without a `page` cursor, an archived
// service account returns an empty list. A `page` cursor that does not match an
// active membership returns a 400 invalid-request error. A cursor stops matching
// when the membership is removed, the workspace is deleted, or the service account
// is archived. Restart pagination from the first page to recover.
func (r *BetaOrganizationServiceAccountWorkspaceService) ListAutoPaging(ctx context.Context, serviceAccountID string, params BetaOrganizationServiceAccountWorkspaceListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaServiceAccountWorkspaceMember] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, serviceAccountID, params, opts...))
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Add a service account to a workspace with the given `workspace_role`.
//
// Mirror of `POST /workspaces/{workspace_id}/service_accounts`, addressed from the
// service-account side; both create the same membership. If the service account is
// already an explicit member of the workspace, its `workspace_role` is replaced
// with the value supplied here. Archived workspaces return 400. Archived service
// accounts cannot be added and are rejected.
func (r *BetaOrganizationServiceAccountWorkspaceService) Add(ctx context.Context, serviceAccountID string, params BetaOrganizationServiceAccountWorkspaceAddParams, opts ...option.RequestOption) (res *BetaServiceAccountWorkspaceMember, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	if serviceAccountID == "" {
		err = errors.New("missing required service_account_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/service_accounts/%s/workspaces?beta=true", serviceAccountID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Remove a service account from a workspace.
//
// Mirror of
// `DELETE /workspaces/{workspace_id}/service_accounts/{service_account_id}`,
// addressed from the service-account side. Removal is idempotent (returns 200 even
// if the membership was already removed). A DELETE against the implicit
// default-workspace membership returns 200 but is a no-op and the membership
// persists; deleting an explicit default-workspace row reverts to the implicit
// `workspace_user` membership. Archived workspaces return 400.
func (r *BetaOrganizationServiceAccountWorkspaceService) Remove(ctx context.Context, workspaceID string, params BetaOrganizationServiceAccountWorkspaceRemoveParams, opts ...option.RequestOption) (res *BetaOrganizationServiceAccountWorkspaceRemoveResponse, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	if params.ServiceAccountID == "" {
		err = errors.New("missing required service_account_id parameter")
		return nil, err
	}
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/service_accounts/%s/workspaces/%s?beta=true", params.ServiceAccountID, workspaceID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type BetaOrganizationServiceAccountWorkspaceRemoveResponse struct {
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
func (r BetaOrganizationServiceAccountWorkspaceRemoveResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationServiceAccountWorkspaceRemoveResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationServiceAccountWorkspaceListParams struct {
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of results per page.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationServiceAccountWorkspaceListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationServiceAccountWorkspaceListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationServiceAccountWorkspaceAddParams struct {
	// Tagged workspace ID to add the service account to.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// Role to assign to the service account in this workspace.
	//
	// Any of "workspace_admin", "workspace_developer",
	// "workspace_restricted_developer", "workspace_user".
	WorkspaceRole BetaNoBillingWorkspaceRole `json:"workspace_role,omitzero" api:"required"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationServiceAccountWorkspaceAddParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationServiceAccountWorkspaceAddParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationServiceAccountWorkspaceAddParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationServiceAccountWorkspaceRemoveParams struct {
	// ID of the service account.
	ServiceAccountID string `path:"service_account_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
