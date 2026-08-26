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

// BetaOrganizationFederationRuleWorkspaceService contains methods and other
// services that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationFederationRuleWorkspaceService] method instead.
type BetaOrganizationFederationRuleWorkspaceService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationFederationRuleWorkspaceService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationFederationRuleWorkspaceService(opts ...option.RequestOption) (r BetaOrganizationFederationRuleWorkspaceService) {
	r = BetaOrganizationFederationRuleWorkspaceService{}
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
func (r *BetaOrganizationFederationRuleWorkspaceService) List(ctx context.Context, federationRuleID string, params BetaOrganizationFederationRuleWorkspaceListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaFederationRuleWorkspace], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if federationRuleID == "" {
		err = errors.New("missing required federation_rule_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/federation_rules/%s/workspaces?beta=true", federationRuleID)
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
// List workspaces where this federation rule is enabled.
//
// Returns all workspace enablements in a single response; the `limit` and `page`
// parameters are accepted but have no effect, and `next_page` is always `null`.
// Returns explicit per-workspace enablements only; for rules with
// `applies_to_all_workspaces` or a legacy single `workspace_id`, check those
// fields on the rule itself.
func (r *BetaOrganizationFederationRuleWorkspaceService) ListAutoPaging(ctx context.Context, federationRuleID string, params BetaOrganizationFederationRuleWorkspaceListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaFederationRuleWorkspace] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, federationRuleID, params, opts...))
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
func (r *BetaOrganizationFederationRuleWorkspaceService) Add(ctx context.Context, federationRuleID string, params BetaOrganizationFederationRuleWorkspaceAddParams, opts ...option.RequestOption) (res *BetaFederationRuleWorkspace, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	if federationRuleID == "" {
		err = errors.New("missing required federation_rule_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/federation_rules/%s/workspaces?beta=true", federationRuleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
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
func (r *BetaOrganizationFederationRuleWorkspaceService) Remove(ctx context.Context, workspaceID string, params BetaOrganizationFederationRuleWorkspaceRemoveParams, opts ...option.RequestOption) (res *BetaOrganizationFederationRuleWorkspaceRemoveResponse, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	if params.FederationRuleID == "" {
		err = errors.New("missing required federation_rule_id parameter")
		return nil, err
	}
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/federation_rules/%s/workspaces/%s?beta=true", params.FederationRuleID, workspaceID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type BetaOrganizationFederationRuleWorkspaceRemoveResponse struct {
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
func (r BetaOrganizationFederationRuleWorkspaceRemoveResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationFederationRuleWorkspaceRemoveResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationFederationRuleWorkspaceListParams struct {
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of results per page.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationFederationRuleWorkspaceListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationFederationRuleWorkspaceListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationFederationRuleWorkspaceAddParams struct {
	// Tagged ID of the workspace to enable this rule for.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationFederationRuleWorkspaceAddParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationFederationRuleWorkspaceAddParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationFederationRuleWorkspaceAddParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationFederationRuleWorkspaceRemoveParams struct {
	// ID of the federation rule.
	FederationRuleID string `path:"federation_rule_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
