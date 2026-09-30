package anthropic

import (
	"context"
	"encoding/json"
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

// OrganizationFederationRuleService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationFederationRuleService] method instead.
type OrganizationFederationRuleService struct {
	Options    []option.RequestOption
	Workspaces OrganizationFederationRuleWorkspaceService
}

// NewOrganizationFederationRuleService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewOrganizationFederationRuleService(opts ...option.RequestOption) (r OrganizationFederationRuleService) {
	r = OrganizationFederationRuleService{}
	r.Options = opts
	r.Workspaces = NewOrganizationFederationRuleWorkspaceService(opts...)
	return
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Create a federation rule owned by your organization.
//
// The referenced issuer and the target service account must already exist in the
// same organization; invalid references are rejected with a 400 error. The
// workspace reference is validated. Membership is not checked at rule creation:
// token exchange resolves a single enabled workspace per call and is rejected
// unless the target service account is a member of that workspace (it is
// implicitly a member of the default workspace). Rules on well-known shared
// issuers (GitHub Actions, GitLab, Buildkite, Terraform Cloud, Google) must
// constrain tenant identity via an identity-bearing claim, a tenant-pinning
// subject prefix (such as `repo:YOUR_ORG/...`), or a CEL condition referencing one
// of those identity claims (e.g. `claims.repository_owner`). OAuth callers may
// only manage rules whose `oauth_scope` is `workspace:developer` or
// `workspace:inference`; other scopes require a Console session.
func (r *OrganizationFederationRuleService) New(ctx context.Context, body OrganizationFederationRuleNewParams, opts ...option.RequestOption) (res *FederationRule, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/federation_rules"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Retrieve a federation rule by its ID (`fdrl_...`).
func (r *OrganizationFederationRuleService) Get(ctx context.Context, federationRuleID string, opts ...option.RequestOption) (res *FederationRule, err error) {
	opts = slices.Concat(r.Options, opts)
	if federationRuleID == "" {
		err = errors.New("missing required federation_rule_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/federation_rules/%s", federationRuleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Partially update a federation rule.
//
// `issuer_id` is immutable. `match` and `target` are replaced as whole objects
// when set. Referenced service accounts and workspaces must exist in your
// organization; invalid references are rejected with a 400 error. Archived rules
// cannot be updated; this returns 400. Create a new rule instead. Rules on
// well-known shared issuers (GitHub Actions, GitLab, Buildkite, Terraform Cloud,
// Google) must constrain tenant identity via an identity-bearing claim, a
// tenant-pinning subject prefix (such as `repo:YOUR_ORG/...`), or a CEL condition
// referencing one of those identity claims (e.g. `claims.repository_owner`). On
// these issuers the requirement is re-checked on every update; if an existing
// rule's stored match does not yet constrain tenant identity, any update (even a
// rename or description change) must also supply a conforming `match` in the same
// request. OAuth callers may only manage rules whose `oauth_scope` is
// `workspace:developer` or `workspace:inference`; other scopes require a Console
// session.
func (r *OrganizationFederationRuleService) Update(ctx context.Context, federationRuleID string, body OrganizationFederationRuleUpdateParams, opts ...option.RequestOption) (res *FederationRule, err error) {
	opts = slices.Concat(r.Options, opts)
	if federationRuleID == "" {
		err = errors.New("missing required federation_rule_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/federation_rules/%s", federationRuleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// List federation rules in your organization.
//
// Optionally filter by issuer with `issuer_id`. Archived rules are excluded unless
// `include_archived=true`.
func (r *OrganizationFederationRuleService) List(ctx context.Context, query OrganizationFederationRuleListParams, opts ...option.RequestOption) (res *pagination.PageCursor[FederationRule], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/federation_rules"
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
// List federation rules in your organization.
//
// Optionally filter by issuer with `issuer_id`. Archived rules are excluded unless
// `include_archived=true`.
func (r *OrganizationFederationRuleService) ListAutoPaging(ctx context.Context, query OrganizationFederationRuleListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[FederationRule] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Archive a federation rule.
//
// Token exchange through this rule stops immediately. Idempotent; re-archiving
// returns the rule with its original `archived_at`. Archiving clears the rule's
// workspace targeting (`workspace_id` and `workspace_ids` are emptied). Tokens
// already minted before archive remain valid until they expire. OAuth callers may
// only manage rules whose `oauth_scope` is `workspace:developer` or
// `workspace:inference`; other scopes require a Console session.
func (r *OrganizationFederationRuleService) Archive(ctx context.Context, federationRuleID string, opts ...option.RequestOption) (res *FederationRule, err error) {
	opts = slices.Concat(r.Options, opts)
	if federationRuleID == "" {
		err = errors.New("missing required federation_rule_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/federation_rules/%s/archive", federationRuleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Authorization rule binding an external OIDC identity to Anthropic.
//
// Evaluates the match conditions and mints an OAuth access token for the resolved
// target, scoped to a single workspace where the rule is enabled (chosen by the
// caller at exchange time when the rule is enabled for more than one). For rules
// enabled via `workspace_ids` or `applies_to_all_workspaces`, the target service
// account must be a member of that workspace (it is implicitly a member of the
// default workspace); rules carrying only the legacy `workspace_id` binding do not
// enforce this.
type FederationRule struct {
	// Tagged ID of the federation rule.
	ID string `json:"id" api:"required"`
	// When true, this rule is enabled for every workspace in the org (including ones
	// created after the rule). `workspace_ids` is ignored at exchange time.
	AppliesToAllWorkspaces bool `json:"applies_to_all_workspaces" api:"required"`
	// If set, this rule is archived and rejects token exchange.
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// Tagged ID (`user_`/`svac_`) of the actor that archived this rule.
	ArchivedByActorID string `json:"archived_by_actor_id" api:"required"`
	// When this rule was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Tagged ID (`user_`/`svac_`) of the actor that created this rule.
	CreatedByActorID string `json:"created_by_actor_id" api:"required"`
	// Optional free-text description.
	Description string `json:"description" api:"required"`
	// Tagged ID of the issuer whose tokens this rule accepts.
	IssuerID string `json:"issuer_id" api:"required"`
	// Issuer's display name at read time.
	IssuerName string `json:"issuer_name" api:"required"`
	// Conditions the verified JWT must satisfy for this rule to apply. All populated
	// matcher fields must pass.
	Match FederationRuleMatch `json:"match" api:"required"`
	// Admin-chosen slug identifier.
	Name string `json:"name" api:"required"`
	// Space-separated OAuth scopes granted on the minted token.
	OAuthScope string `json:"oauth_scope" api:"required"`
	// Identity that tokens minted via this rule act as. Currently always a
	// `service_account` target.
	Target ServiceAccountTarget `json:"target" api:"required"`
	// Lifetime in seconds of access tokens minted via this rule. Minted tokens are
	// capped at `max(60, min(this value, 2 × remaining assertion validity))` seconds.
	TokenLifetimeSeconds int64                   `json:"token_lifetime_seconds" api:"required"`
	Type                 constant.FederationRule `json:"type" default:"federation_rule"`
	// When this rule was last updated.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// Tagged ID (`user_`/`svac_`) of the actor that last updated this rule.
	UpdatedByActorID string `json:"updated_by_actor_id" api:"required"`
	// Tagged IDs of the workspaces this rule is enabled for. May be empty for older
	// rules that only carry the legacy `workspace_id` binding. Ignored at exchange
	// time when `applies_to_all_workspaces` is true (the list may still be non-empty).
	WorkspaceIDs []string `json:"workspace_ids" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                     respjson.Field
		AppliesToAllWorkspaces respjson.Field
		ArchivedAt             respjson.Field
		ArchivedByActorID      respjson.Field
		CreatedAt              respjson.Field
		CreatedByActorID       respjson.Field
		Description            respjson.Field
		IssuerID               respjson.Field
		IssuerName             respjson.Field
		Match                  respjson.Field
		Name                   respjson.Field
		OAuthScope             respjson.Field
		Target                 respjson.Field
		TokenLifetimeSeconds   respjson.Field
		Type                   respjson.Field
		UpdatedAt              respjson.Field
		UpdatedByActorID       respjson.Field
		WorkspaceIDs           respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FederationRule) RawJSON() string { return r.JSON.raw }
func (r *FederationRule) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Does the incoming JWT qualify?
//
// All populated fields must pass; omitted fields are skipped. At least one of
// `subject_prefix` (other than a wildcard-only value like `*`), `claims`, or
// `condition` is required; `audience` alone is not sufficient.
type FederationRuleMatch struct {
	// Exact match against the `aud` claim (any element if array). When omitted, the
	// JWT's `aud` must still equal Anthropic's expected audience for the issuer;
	// setting this field overrides that default.
	Audience string `json:"audience" api:"nullable"`
	// Exact-match `{claim: value}` pairs against top-level claims. Only string-valued
	// claims can be matched; use `condition` for non-string claims.
	Claims map[string]string `json:"claims" api:"nullable"`
	// CEL expression over claims for logic the structural fields can't express. Must
	// evaluate to a boolean and may reference only the `claims` variable; a
	// constant-true expression (such as `true`) is rejected with 400.
	Condition string `json:"condition" api:"nullable"`
	// Match the verified JWT `sub` claim. Exact match unless the value ends with `*`,
	// in which case it is a prefix match. Example:
	// `repo:my-org/my-repo:ref:refs/heads/main`.
	SubjectPrefix string `json:"subject_prefix" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Audience      respjson.Field
		Claims        respjson.Field
		Condition     respjson.Field
		SubjectPrefix respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FederationRuleMatch) RawJSON() string { return r.JSON.raw }
func (r *FederationRuleMatch) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this FederationRuleMatch to a FederationRuleMatchParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// FederationRuleMatchParam.Overrides()
func (r FederationRuleMatch) ToParam() FederationRuleMatchParam {
	return param.Override[FederationRuleMatchParam](json.RawMessage(r.RawJSON()))
}

// Does the incoming JWT qualify?
//
// All populated fields must pass; omitted fields are skipped. At least one of
// `subject_prefix` (other than a wildcard-only value like `*`), `claims`, or
// `condition` is required; `audience` alone is not sufficient.
type FederationRuleMatchParam struct {
	// Exact match against the `aud` claim (any element if array). When omitted, the
	// JWT's `aud` must still equal Anthropic's expected audience for the issuer;
	// setting this field overrides that default.
	Audience param.Opt[string] `json:"audience,omitzero"`
	// CEL expression over claims for logic the structural fields can't express. Must
	// evaluate to a boolean and may reference only the `claims` variable; a
	// constant-true expression (such as `true`) is rejected with 400.
	Condition param.Opt[string] `json:"condition,omitzero"`
	// Match the verified JWT `sub` claim. Exact match unless the value ends with `*`,
	// in which case it is a prefix match. Example:
	// `repo:my-org/my-repo:ref:refs/heads/main`.
	SubjectPrefix param.Opt[string] `json:"subject_prefix,omitzero"`
	// Exact-match `{claim: value}` pairs against top-level claims. Only string-valued
	// claims can be matched; use `condition` for non-string claims.
	Claims map[string]string `json:"claims,omitzero"`
	paramObj
}

func (r FederationRuleMatchParam) MarshalJSON() (data []byte, err error) {
	type shadow FederationRuleMatchParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *FederationRuleMatchParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type FederationRuleWorkspace struct {
	// When this workspace was enabled for the rule.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Tagged ID (`user_...` or `svac_...`) of the actor that enabled this workspace
	// for the rule, if known.
	CreatedByActorID string `json:"created_by_actor_id" api:"required"`
	// Tagged ID of the federation rule.
	FederationRuleID string                           `json:"federation_rule_id" api:"required"`
	Type             constant.FederationRuleWorkspace `json:"type" default:"federation_rule_workspace"`
	// Tagged ID of the workspace this rule is enabled for.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// Workspace display name. Populated when listing; null in the enable response.
	WorkspaceName string `json:"workspace_name" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CreatedAt        respjson.Field
		CreatedByActorID respjson.Field
		FederationRuleID respjson.Field
		Type             respjson.Field
		WorkspaceID      respjson.Field
		WorkspaceName    respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r FederationRuleWorkspace) RawJSON() string { return r.JSON.raw }
func (r *FederationRuleWorkspace) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Bind to a fixed service account by ID.
type ServiceAccountTarget struct {
	// Tagged ID of the service account to mint tokens for.
	ServiceAccountID string                  `json:"service_account_id" api:"required"`
	Type             constant.ServiceAccount `json:"type" default:"service_account"`
	// Service account's display name at read time. Ignored on writes.
	ServiceAccountName string `json:"service_account_name" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ServiceAccountID   respjson.Field
		Type               respjson.Field
		ServiceAccountName respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ServiceAccountTarget) RawJSON() string { return r.JSON.raw }
func (r *ServiceAccountTarget) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ServiceAccountTarget to a ServiceAccountTargetParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ServiceAccountTargetParam.Overrides()
func (r ServiceAccountTarget) ToParam() ServiceAccountTargetParam {
	return param.Override[ServiceAccountTargetParam](json.RawMessage(r.RawJSON()))
}

// Bind to a fixed service account by ID.
//
// The properties ServiceAccountID, Type are required.
type ServiceAccountTargetParam struct {
	// Tagged ID of the service account to mint tokens for.
	ServiceAccountID string `json:"service_account_id" api:"required"`
	// Service account's display name at read time. Ignored on writes.
	ServiceAccountName param.Opt[string] `json:"service_account_name,omitzero"`
	// This field can be elided, and will marshal its zero value as "service_account".
	Type constant.ServiceAccount `json:"type" default:"service_account"`
	paramObj
}

func (r ServiceAccountTargetParam) MarshalJSON() (data []byte, err error) {
	type shadow ServiceAccountTargetParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ServiceAccountTargetParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationFederationRuleNewParams struct {
	// Tagged ID of the federation issuer.
	IssuerID string `json:"issuer_id" api:"required"`
	// Conditions the verified JWT must satisfy for this rule to apply. At least one of
	// `subject_prefix` (other than a wildcard-only value like `*`), `claims`, or
	// `condition` is required; `audience` alone is not sufficient.
	Match FederationRuleMatchParam `json:"match,omitzero" api:"required"`
	// Slug identifier (lowercase, digits, hyphens). Unique within the organization; a
	// duplicate name returns 409.
	Name string `json:"name" api:"required"`
	// Space-separated OAuth scopes. OAuth callers may only set `workspace:developer`
	// or `workspace:inference`; other scopes (such as `org:admin`) require a Console
	// session.
	OAuthScope string `json:"oauth_scope" api:"required"`
	// Identity that tokens minted via this rule act as. Currently always a
	// `service_account` target.
	Target ServiceAccountTargetParam `json:"target,omitzero" api:"required"`
	// Optional free-text description.
	Description param.Opt[string] `json:"description,omitzero"`
	// Tagged ID of the workspace to enable this rule for. Required unless
	// `applies_to_all_workspaces` is true. Additional workspaces can be added via the
	// `/federation_rules/{federation_rule_id}/workspaces` sub-resource.
	WorkspaceID param.Opt[string] `json:"workspace_id,omitzero"`
	// When true, enable this rule for every workspace in the org (including workspaces
	// created later).
	AppliesToAllWorkspaces param.Opt[bool] `json:"applies_to_all_workspaces,omitzero"`
	// Lifetime in seconds for access tokens minted via this rule (60-86400). Defaults
	// to 3600 (1h). Minted tokens are capped at
	// `max(60, min(this value, 2 × remaining assertion validity))` seconds.
	TokenLifetimeSeconds param.Opt[int64] `json:"token_lifetime_seconds,omitzero"`
	paramObj
}

func (r OrganizationFederationRuleNewParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationFederationRuleNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationFederationRuleNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationFederationRuleUpdateParams struct {
	// When true, enables this rule for every workspace in the org (including
	// workspaces created later). Setting `false` is rejected with 400 if no workspace
	// would remain enabled; a rule with only a legacy `workspace_id` binding continues
	// to mint.
	AppliesToAllWorkspaces param.Opt[bool] `json:"applies_to_all_workspaces,omitzero"`
	// Replaces the description. Omit to leave unchanged; send `null` to clear (the
	// field is stored as an empty string).
	Description param.Opt[string] `json:"description,omitzero"`
	// Replaces the slug identifier (lowercase, digits, hyphens). Unique within the
	// organization; a duplicate name returns 409.
	Name param.Opt[string] `json:"name,omitzero"`
	// Replaces the space-separated OAuth scopes granted on minted tokens. OAuth
	// callers may only set `workspace:developer` or `workspace:inference`; other
	// scopes (such as `org:admin`) require a Console session.
	OAuthScope param.Opt[string] `json:"oauth_scope,omitzero"`
	// Replaces the lifetime in seconds for access tokens minted via this rule
	// (60-86400). Minted tokens are capped at
	// `max(60, min(this value, 2 × remaining assertion validity))` seconds.
	TokenLifetimeSeconds param.Opt[int64] `json:"token_lifetime_seconds,omitzero"`
	// Replaces the existing single workspace enablement (the previous one is removed).
	// Rejected with 400 if the rule is enabled for more than one workspace; use the
	// `/federation_rules/{federation_rule_id}/workspaces` sub-resource instead.
	WorkspaceID param.Opt[string] `json:"workspace_id,omitzero"`
	// Replaces the entire match object. All populated matcher fields must pass.
	Match FederationRuleMatchParam `json:"match,omitzero"`
	// Replaces the entire target object. Currently always a `service_account` target.
	Target ServiceAccountTargetParam `json:"target,omitzero"`
	paramObj
}

func (r OrganizationFederationRuleUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationFederationRuleUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationFederationRuleUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationFederationRuleListParams struct {
	// Filter to rules referencing this federation issuer.
	IssuerID param.Opt[string] `query:"issuer_id,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Include archived resources. Defaults to false.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Number of results per page.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrganizationFederationRuleListParams]'s query parameters as
// `url.Values`.
func (r OrganizationFederationRuleListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
