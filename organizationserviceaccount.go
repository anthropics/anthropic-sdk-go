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

// OrganizationServiceAccountService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationServiceAccountService] method instead.
type OrganizationServiceAccountService struct {
	Options    []option.RequestOption
	Workspaces OrganizationServiceAccountWorkspaceService
}

// NewOrganizationServiceAccountService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewOrganizationServiceAccountService(opts ...option.RequestOption) (r OrganizationServiceAccountService) {
	r = OrganizationServiceAccountService{}
	r.Options = opts
	r.Workspaces = NewOrganizationServiceAccountWorkspaceService(opts...)
	return
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Create a service account.
//
// A service account is a named workload identity that federation rules target.
// `organization_role` is `developer` (default) or `admin`; a rule may only be
// created or retargeted to grant `org:admin` scope when the target's
// `organization_role` is `admin`. Creating an `admin`-role service account
// requires an interactive credential (a user OAuth token or a Console session) — a
// workload may only create `developer`-role service accounts.
func (r *OrganizationServiceAccountService) New(ctx context.Context, body OrganizationServiceAccountNewParams, opts ...option.RequestOption) (res *ServiceAccount, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/service_accounts"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Retrieve a service account by its ID (`svac_...`).
func (r *OrganizationServiceAccountService) Get(ctx context.Context, serviceAccountID string, opts ...option.RequestOption) (res *ServiceAccount, err error) {
	opts = slices.Concat(r.Options, opts)
	if serviceAccountID == "" {
		err = errors.New("missing required service_account_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/service_accounts/%s", serviceAccountID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Update a service account.
//
// Only `description` and `organization_role` are mutable; `name` cannot be
// changed. Archived service accounts cannot be updated; this returns 400. Setting
// `organization_role` to `admin` (even when unchanged) requires an interactive
// credential (a user OAuth token or a Console session).
func (r *OrganizationServiceAccountService) Update(ctx context.Context, serviceAccountID string, body OrganizationServiceAccountUpdateParams, opts ...option.RequestOption) (res *ServiceAccount, err error) {
	opts = slices.Concat(r.Options, opts)
	if serviceAccountID == "" {
		err = errors.New("missing required service_account_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/service_accounts/%s", serviceAccountID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// List service accounts in the caller's organization.
//
// Results are ordered by creation time, newest first. Use `limit` and the
// `next_page` cursor to paginate; set `include_archived=true` to include archived
// service accounts.
func (r *OrganizationServiceAccountService) List(ctx context.Context, query OrganizationServiceAccountListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ServiceAccount], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/service_accounts"
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
// List service accounts in the caller's organization.
//
// Results are ordered by creation time, newest first. Use `limit` and the
// `next_page` cursor to paginate; set `include_archived=true` to include archived
// service accounts.
func (r *OrganizationServiceAccountService) ListAutoPaging(ctx context.Context, query OrganizationServiceAccountListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ServiceAccount] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Archive a service account.
//
// Idempotent; re-archiving returns the service account with its original
// `archived_at`. Rejected with 400 if any live (non-archived) federation rule
// still targets this service account, same as issuer archival; archive those rules
// first or change their target to another service account.
func (r *OrganizationServiceAccountService) Archive(ctx context.Context, serviceAccountID string, opts ...option.RequestOption) (res *ServiceAccount, err error) {
	opts = slices.Concat(r.Options, opts)
	if serviceAccountID == "" {
		err = errors.New("missing required service_account_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/service_accounts/%s/archive", serviceAccountID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Named non-human identity within the caller's organization.
//
// A service account is a pure identity: name + org. Authorization lives on
// whatever references it (federation rules).
type ServiceAccount struct {
	// Tagged ID of the service account.
	ID string `json:"id" api:"required"`
	// If set, this service account is archived.
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// Tagged ID (`user_`/`svac_`) of the actor that archived this service account.
	ArchivedByActorID string `json:"archived_by_actor_id" api:"required"`
	// When this service account was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Tagged ID (`user_`/`svac_`) of the actor that created this service account.
	CreatedByActorID string `json:"created_by_actor_id" api:"required"`
	// Optional free-text description.
	Description string `json:"description" api:"required"`
	// Admin-chosen slug identifier.
	Name string `json:"name" api:"required"`
	// Org-level role. A federation rule may only be created or retargeted to grant
	// `org:admin` scope when this is `admin`. A rule granting `org:admin` whose target
	// is later demoted to `developer` is rejected at token exchange. Rules granting
	// `org:admin` are managed in the Console.
	//
	// Any of "admin", "developer".
	OrganizationRole ServiceAccountOrganizationRole `json:"organization_role" api:"required"`
	Type             constant.ServiceAccount        `json:"type" default:"service_account"`
	// When this service account was last updated.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// Tagged ID (`user_`/`svac_`) of the actor that last updated this service account.
	UpdatedByActorID string `json:"updated_by_actor_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		ArchivedAt        respjson.Field
		ArchivedByActorID respjson.Field
		CreatedAt         respjson.Field
		CreatedByActorID  respjson.Field
		Description       respjson.Field
		Name              respjson.Field
		OrganizationRole  respjson.Field
		Type              respjson.Field
		UpdatedAt         respjson.Field
		UpdatedByActorID  respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ServiceAccount) RawJSON() string { return r.JSON.raw }
func (r *ServiceAccount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Org-level role. A federation rule may only be created or retargeted to grant
// `org:admin` scope when this is `admin`. A rule granting `org:admin` whose target
// is later demoted to `developer` is rejected at token exchange. Rules granting
// `org:admin` are managed in the Console.
type ServiceAccountOrganizationRole string

const (
	ServiceAccountOrganizationRoleAdmin     ServiceAccountOrganizationRole = "admin"
	ServiceAccountOrganizationRoleDeveloper ServiceAccountOrganizationRole = "developer"
)

type ServiceAccountWorkspaceMember struct {
	// Tagged ID (`user_...`/`svac_...`) of the actor who created this membership.
	CreatedByActorID string `json:"created_by_actor_id" api:"required"`
	// True when this is the implicit default-workspace membership every service
	// account has when no explicit membership exists. Implicit memberships have role
	// `workspace_user` and cannot be removed.
	Implicit bool `json:"implicit" api:"required"`
	// Tagged service account ID (`svac_...`).
	ServiceAccountID string                                 `json:"service_account_id" api:"required"`
	Type             constant.ServiceAccountWorkspaceMember `json:"type" default:"service_account_workspace_member"`
	// Tagged workspace ID (`wrkspc_...`).
	WorkspaceID string `json:"workspace_id" api:"required"`
	// Role of the service account in this workspace. Service accounts cannot hold the
	// `workspace_billing` role.
	//
	// Any of "workspace_admin", "workspace_billing", "workspace_developer",
	// "workspace_restricted_developer", "workspace_user".
	WorkspaceRole WorkspaceRole `json:"workspace_role" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CreatedByActorID respjson.Field
		Implicit         respjson.Field
		ServiceAccountID respjson.Field
		Type             respjson.Field
		WorkspaceID      respjson.Field
		WorkspaceRole    respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ServiceAccountWorkspaceMember) RawJSON() string { return r.JSON.raw }
func (r *ServiceAccountWorkspaceMember) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationServiceAccountNewParams struct {
	// Slug identifier (lowercase, digits, hyphens). Unique within the organization; a
	// duplicate name returns 409.
	Name string `json:"name" api:"required"`
	// Optional free-text description.
	Description param.Opt[string] `json:"description,omitzero"`
	// Org-level role. Defaults to `developer`.
	//
	// Any of "admin", "developer".
	OrganizationRole OrganizationServiceAccountNewParamsOrganizationRole `json:"organization_role,omitzero"`
	paramObj
}

func (r OrganizationServiceAccountNewParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationServiceAccountNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationServiceAccountNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Org-level role. Defaults to `developer`.
type OrganizationServiceAccountNewParamsOrganizationRole string

const (
	OrganizationServiceAccountNewParamsOrganizationRoleAdmin     OrganizationServiceAccountNewParamsOrganizationRole = "admin"
	OrganizationServiceAccountNewParamsOrganizationRoleDeveloper OrganizationServiceAccountNewParamsOrganizationRole = "developer"
)

type OrganizationServiceAccountUpdateParams struct {
	// Replaces the description. Omit to leave unchanged; send `null` to clear (the
	// field is stored as an empty string).
	Description param.Opt[string] `json:"description,omitzero"`
	// Replaces the org-level role. Omit or send `null` to leave unchanged.
	//
	// Any of "admin", "developer".
	OrganizationRole OrganizationServiceAccountUpdateParamsOrganizationRole `json:"organization_role,omitzero"`
	paramObj
}

func (r OrganizationServiceAccountUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationServiceAccountUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationServiceAccountUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Replaces the org-level role. Omit or send `null` to leave unchanged.
type OrganizationServiceAccountUpdateParamsOrganizationRole string

const (
	OrganizationServiceAccountUpdateParamsOrganizationRoleAdmin     OrganizationServiceAccountUpdateParamsOrganizationRole = "admin"
	OrganizationServiceAccountUpdateParamsOrganizationRoleDeveloper OrganizationServiceAccountUpdateParamsOrganizationRole = "developer"
)

type OrganizationServiceAccountListParams struct {
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Include archived resources. Defaults to false.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Number of results per page.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrganizationServiceAccountListParams]'s query parameters as
// `url.Values`.
func (r OrganizationServiceAccountListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
