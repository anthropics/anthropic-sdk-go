// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

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

// BetaOrganizationServiceAccountService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationServiceAccountService] method instead.
type BetaOrganizationServiceAccountService struct {
	Options    []option.RequestOption
	Workspaces BetaOrganizationServiceAccountWorkspaceService
}

// NewBetaOrganizationServiceAccountService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationServiceAccountService(opts ...option.RequestOption) (r BetaOrganizationServiceAccountService) {
	r = BetaOrganizationServiceAccountService{}
	r.Options = opts
	r.Workspaces = NewBetaOrganizationServiceAccountWorkspaceService(opts...)
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
func (r *BetaOrganizationServiceAccountService) New(ctx context.Context, params BetaOrganizationServiceAccountNewParams, opts ...option.RequestOption) (res *BetaServiceAccount, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/service_accounts?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Retrieve a service account by its ID (`svac_...`).
func (r *BetaOrganizationServiceAccountService) Get(ctx context.Context, serviceAccountID string, query BetaOrganizationServiceAccountGetParams, opts ...option.RequestOption) (res *BetaServiceAccount, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	if serviceAccountID == "" {
		err = errors.New("missing required service_account_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/service_accounts/%s?beta=true", serviceAccountID)
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
func (r *BetaOrganizationServiceAccountService) Update(ctx context.Context, serviceAccountID string, params BetaOrganizationServiceAccountUpdateParams, opts ...option.RequestOption) (res *BetaServiceAccount, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	if serviceAccountID == "" {
		err = errors.New("missing required service_account_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/service_accounts/%s?beta=true", serviceAccountID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
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
func (r *BetaOrganizationServiceAccountService) List(ctx context.Context, params BetaOrganizationServiceAccountListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaServiceAccount], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/service_accounts?beta=true"
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
// List service accounts in the caller's organization.
//
// Results are ordered by creation time, newest first. Use `limit` and the
// `next_page` cursor to paginate; set `include_archived=true` to include archived
// service accounts.
func (r *BetaOrganizationServiceAccountService) ListAutoPaging(ctx context.Context, params BetaOrganizationServiceAccountListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaServiceAccount] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
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
func (r *BetaOrganizationServiceAccountService) Archive(ctx context.Context, serviceAccountID string, body BetaOrganizationServiceAccountArchiveParams, opts ...option.RequestOption) (res *BetaServiceAccount, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	if serviceAccountID == "" {
		err = errors.New("missing required service_account_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/service_accounts/%s/archive?beta=true", serviceAccountID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Named non-human identity within the caller's organization.
//
// A service account is a pure identity: name + org. Authorization lives on
// whatever references it (federation rules).
type BetaServiceAccount struct {
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
	OrganizationRole BetaServiceAccountOrganizationRole `json:"organization_role" api:"required"`
	Type             constant.ServiceAccount            `json:"type" default:"service_account"`
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
func (r BetaServiceAccount) RawJSON() string { return r.JSON.raw }
func (r *BetaServiceAccount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Org-level role. A federation rule may only be created or retargeted to grant
// `org:admin` scope when this is `admin`. A rule granting `org:admin` whose target
// is later demoted to `developer` is rejected at token exchange. Rules granting
// `org:admin` are managed in the Console.
type BetaServiceAccountOrganizationRole string

const (
	BetaServiceAccountOrganizationRoleAdmin     BetaServiceAccountOrganizationRole = "admin"
	BetaServiceAccountOrganizationRoleDeveloper BetaServiceAccountOrganizationRole = "developer"
)

type BetaServiceAccountWorkspaceMember struct {
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
	WorkspaceRole BetaWorkspaceRole `json:"workspace_role" api:"required"`
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
func (r BetaServiceAccountWorkspaceMember) RawJSON() string { return r.JSON.raw }
func (r *BetaServiceAccountWorkspaceMember) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationServiceAccountNewParams struct {
	// Slug identifier (lowercase, digits, hyphens). Unique within the organization; a
	// duplicate name returns 409.
	Name string `json:"name" api:"required"`
	// Optional free-text description.
	Description param.Opt[string] `json:"description,omitzero"`
	// Org-level role. Defaults to `developer`.
	//
	// Any of "admin", "developer".
	OrganizationRole BetaOrganizationServiceAccountNewParamsOrganizationRole `json:"organization_role,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationServiceAccountNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationServiceAccountNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationServiceAccountNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Org-level role. Defaults to `developer`.
type BetaOrganizationServiceAccountNewParamsOrganizationRole string

const (
	BetaOrganizationServiceAccountNewParamsOrganizationRoleAdmin     BetaOrganizationServiceAccountNewParamsOrganizationRole = "admin"
	BetaOrganizationServiceAccountNewParamsOrganizationRoleDeveloper BetaOrganizationServiceAccountNewParamsOrganizationRole = "developer"
)

type BetaOrganizationServiceAccountGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationServiceAccountUpdateParams struct {
	// Replaces the description. Omit to leave unchanged; send `null` to clear (the
	// field is stored as an empty string).
	Description param.Opt[string] `json:"description,omitzero"`
	// Replaces the org-level role. Omit or send `null` to leave unchanged.
	//
	// Any of "admin", "developer".
	OrganizationRole BetaOrganizationServiceAccountUpdateParamsOrganizationRole `json:"organization_role,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationServiceAccountUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationServiceAccountUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationServiceAccountUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Replaces the org-level role. Omit or send `null` to leave unchanged.
type BetaOrganizationServiceAccountUpdateParamsOrganizationRole string

const (
	BetaOrganizationServiceAccountUpdateParamsOrganizationRoleAdmin     BetaOrganizationServiceAccountUpdateParamsOrganizationRole = "admin"
	BetaOrganizationServiceAccountUpdateParamsOrganizationRoleDeveloper BetaOrganizationServiceAccountUpdateParamsOrganizationRole = "developer"
)

type BetaOrganizationServiceAccountListParams struct {
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Include archived resources. Defaults to false.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Number of results per page.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationServiceAccountListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationServiceAccountListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationServiceAccountArchiveParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
