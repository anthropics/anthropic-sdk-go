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

// OrganizationFederationRuleWorkspaceService contains methods and other services
// that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationFederationRuleWorkspaceService] method instead.
type OrganizationFederationRuleWorkspaceService struct {
	Options []option.RequestOption
}

// NewOrganizationFederationRuleWorkspaceService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewOrganizationFederationRuleWorkspaceService(opts ...option.RequestOption) (r OrganizationFederationRuleWorkspaceService) {
	r = OrganizationFederationRuleWorkspaceService{}
	r.Options = opts
	return
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// List workspaces where this federation rule is enabled.
//
// Returns all workspace enablements in a single response; the `limit` and `page`
// parameters are accepted but have no effect, and `next_page` is always `null`.
// Returns explicit per-workspace enablements only; for rules with
// `applies_to_all_workspaces` or a legacy single `workspace_id`, check those
// fields on the rule itself.
func (r *OrganizationFederationRuleWorkspaceService) List(ctx context.Context, federationRuleID string, query OrganizationFederationRuleWorkspaceListParams, opts ...option.RequestOption) (res *pagination.PageCursor[FederationRuleWorkspace], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if federationRuleID == "" {
		err = errors.New("missing required federation_rule_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/federation_rules/%s/workspaces", url.PathEscape(federationRuleID))
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
// List workspaces where this federation rule is enabled.
//
// Returns all workspace enablements in a single response; the `limit` and `page`
// parameters are accepted but have no effect, and `next_page` is always `null`.
// Returns explicit per-workspace enablements only; for rules with
// `applies_to_all_workspaces` or a legacy single `workspace_id`, check those
// fields on the rule itself.
func (r *OrganizationFederationRuleWorkspaceService) ListAutoPaging(ctx context.Context, federationRuleID string, query OrganizationFederationRuleWorkspaceListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[FederationRuleWorkspace] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, federationRuleID, query, opts...))
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Enable a federation rule for a workspace.
//
// Idempotent; re-enabling returns the existing enablement. The rule and workspace
// must both belong to your organization. Membership of the rule's target service
// account in this workspace is not checked at enablement: token exchange into this
// workspace is rejected unless the target is a member (it is implicitly a member
// of the default workspace). Archived rules are rejected with 400. OAuth callers
// may only manage rules whose `oauth_scope` is `workspace:developer` or
// `workspace:inference`; other scopes require a Console session.
func (r *OrganizationFederationRuleWorkspaceService) Add(ctx context.Context, federationRuleID string, body OrganizationFederationRuleWorkspaceAddParams, opts ...option.RequestOption) (res *FederationRuleWorkspace, err error) {
	opts = slices.Concat(r.Options, opts)
	if federationRuleID == "" {
		err = errors.New("missing required federation_rule_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/federation_rules/%s/workspaces", url.PathEscape(federationRuleID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Disable a federation rule for a workspace.
//
// Idempotent; succeeds even if the enablement was already removed. OAuth callers
// may only manage rules whose `oauth_scope` is `workspace:developer` or
// `workspace:inference`; other scopes require a Console session.
func (r *OrganizationFederationRuleWorkspaceService) Remove(ctx context.Context, workspaceID string, body OrganizationFederationRuleWorkspaceRemoveParams, opts ...option.RequestOption) (res *OrganizationFederationRuleWorkspaceRemoveResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if body.FederationRuleID == "" {
		err = errors.New("missing required federation_rule_id parameter")
		return nil, err
	}
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/federation_rules/%s/workspaces/%s", url.PathEscape(body.FederationRuleID), url.PathEscape(workspaceID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type OrganizationFederationRuleWorkspaceRemoveResponse struct {
	// Tagged ID of the federation rule.
	FederationRuleID string                                  `json:"federation_rule_id" api:"required"`
	Type             constant.FederationRuleWorkspaceDeleted `json:"type" default:"federation_rule_workspace_deleted"`
	// Tagged ID of the workspace named in the delete request. Removal is idempotent.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FederationRuleID respjson.Field
		Type             respjson.Field
		WorkspaceID      respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationFederationRuleWorkspaceRemoveResponse) RawJSON() string { return r.JSON.raw }
func (r *OrganizationFederationRuleWorkspaceRemoveResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationFederationRuleWorkspaceListParams struct {
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of results per page.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrganizationFederationRuleWorkspaceListParams]'s query
// parameters as `url.Values`.
func (r OrganizationFederationRuleWorkspaceListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type OrganizationFederationRuleWorkspaceAddParams struct {
	// Tagged ID of the workspace to enable this rule for.
	WorkspaceID string `json:"workspace_id" api:"required"`
	paramObj
}

func (r OrganizationFederationRuleWorkspaceAddParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationFederationRuleWorkspaceAddParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationFederationRuleWorkspaceAddParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationFederationRuleWorkspaceRemoveParams struct {
	// ID of the federation rule.
	FederationRuleID string `path:"federation_rule_id" api:"required" json:"-"`
	paramObj
}
