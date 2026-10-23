package anthropic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
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

// BetaOrganizationClaudeTagSlackWorkspaceService contains methods and other
// services that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationClaudeTagSlackWorkspaceService] method instead.
type BetaOrganizationClaudeTagSlackWorkspaceService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationClaudeTagSlackWorkspaceService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationClaudeTagSlackWorkspaceService(opts ...option.RequestOption) (r BetaOrganizationClaudeTagSlackWorkspaceService) {
	r = BetaOrganizationClaudeTagSlackWorkspaceService{}
	r.Options = opts
	return
}

// Get one Slack workspace and its settings. By default `settings` holds what the
// workspace itself stores: a setting that is null there is not set on the
// workspace and is inherited. With `include_inherited=true`, `settings` holds what
// applies in the workspace once the organization's Slack defaults are taken into
// account, and the request fails with 503 rather than resolve against a settings
// level that could not be read. Either way `sources` names the settings level each
// value comes from.
func (r *BetaOrganizationClaudeTagSlackWorkspaceService) Get(ctx context.Context, id string, params BetaOrganizationClaudeTagSlackWorkspaceGetParams, opts ...option.RequestOption) (res *BetaClaudeTagSlackWorkspace, err error) {
	if len(params.Betas) > 0 {
		headerValues := make([]string, len(params.Betas))
		for i, v := range params.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "experimental-claude-tag-2026-10-21")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/claude_tag/slack_workspaces/%s", url.PathEscape(id))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_tag.slack_workspaces.retrieve", "/v1/organizations/claude_tag/slack_workspaces/{id}", "include_inherited"), requestconfig.WithSpanAttribute("anthropic.organization.claude_tag.slack_workspace.id", id))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Change the settings a Slack workspace stores. The update is partial: a setting
// left out is unchanged, and a setting sent as null is cleared, so the workspace
// inherits it. If any value is refused, no setting in the request is stored.
func (r *BetaOrganizationClaudeTagSlackWorkspaceService) Update(ctx context.Context, id string, params BetaOrganizationClaudeTagSlackWorkspaceUpdateParams, opts ...option.RequestOption) (res *BetaClaudeTagSlackWorkspace, err error) {
	if len(params.Betas) > 0 {
		headerValues := make([]string, len(params.Betas))
		for i, v := range params.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "experimental-claude-tag-2026-10-21")}, opts...)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/claude_tag/slack_workspaces/%s", url.PathEscape(id))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_tag.slack_workspaces.update", "/v1/organizations/claude_tag/slack_workspaces/{id}"), requestconfig.WithSpanAttribute("anthropic.organization.claude_tag.slack_workspace.id", id))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List the set-up Slack workspaces the organization owns, newest first by
// `created_at`. Pass `slack_workspace_id` to look one workspace up by its Slack
// identifier: `data` then holds that workspace, or is empty when it is not set up.
// A workspace of a Slack Enterprise Grid the organization has connected returns
// 404 until it is set up. With `include_inherited=true` each workspace's
// `settings` holds what applies in it, and the request fails with 503 rather than
// resolve against a settings level that could not be read. A `page` cursor resumes
// only a request with the same `include_inherited`.
func (r *BetaOrganizationClaudeTagSlackWorkspaceService) List(ctx context.Context, params BetaOrganizationClaudeTagSlackWorkspaceListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaClaudeTagSlackWorkspace], err error) {
	var raw *http.Response
	if len(params.Betas) > 0 {
		headerValues := make([]string, len(params.Betas))
		for i, v := range params.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "experimental-claude-tag-2026-10-21"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/claude_tag/slack_workspaces"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_tag.slack_workspaces.list", "/v1/organizations/claude_tag/slack_workspaces", "include_inherited", "limit", "page", "slack_workspace_id"))
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

// List the set-up Slack workspaces the organization owns, newest first by
// `created_at`. Pass `slack_workspace_id` to look one workspace up by its Slack
// identifier: `data` then holds that workspace, or is empty when it is not set up.
// A workspace of a Slack Enterprise Grid the organization has connected returns
// 404 until it is set up. With `include_inherited=true` each workspace's
// `settings` holds what applies in it, and the request fails with 503 rather than
// resolve against a settings level that could not be read. A `page` cursor resumes
// only a request with the same `include_inherited`.
func (r *BetaOrganizationClaudeTagSlackWorkspaceService) ListAutoPaging(ctx context.Context, params BetaOrganizationClaudeTagSlackWorkspaceListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaClaudeTagSlackWorkspace] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// A Slack workspace the organization has connected, itself or through its Slack
// Enterprise Grid, with the settings it stores or, with `include_inherited`, the
// settings that apply in it.
type BetaClaudeTagSlackWorkspace struct {
	// The workspace's identifier. Starts with `slws_`. It does not change for as long
	// as the workspace is set up. A workspace that is removed and set up again has a
	// new one, and the old one keeps returning 404; look the workspace up by its Slack
	// identifier for the new one.
	ID string `json:"id" api:"required"`
	// Time the workspace was set up, by whichever route. It does not change for as
	// long as `id` does not.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// What the workspace itself stores or, with `include_inherited`, what applies in
	// it. `sources` says which.
	Settings BetaClaudeTagSettings `json:"settings" api:"required"`
	// Slack's own identifier for the workspace.
	Slack BetaClaudeTagSlackWorkspaceIdentifiers `json:"slack" api:"required"`
	// For each setting, the settings level its value in `settings` comes from.
	Sources BetaClaudeTagSettingSources `json:"sources" api:"required"`
	Type    constant.SlackWorkspace     `json:"type" default:"slack_workspace"`
	// Time Claude Tag's record of the workspace last changed. It moves on every change
	// to a setting, and can also move when something no response holds changes, so a
	// later value says to read the workspace again, not that a setting changed.
	// Attaching or detaching a Bundle does not by itself move it.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Settings    respjson.Field
		Slack       respjson.Field
		Sources     respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagSlackWorkspace) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagSlackWorkspace) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationClaudeTagSlackWorkspaceGetParams struct {
	// When `true`, `settings` holds what applies in the workspace once the
	// organization's Slack defaults are taken into account, and every member of
	// `sources` names the level a value comes from. When `false` or left out,
	// `settings` holds what the workspace itself stores.
	IncludeInherited param.Opt[bool] `query:"include_inherited,omitzero" json:"-"`
	// This endpoint is in beta: requests must send
	// `experimental-claude-tag-2026-10-21` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationClaudeTagSlackWorkspaceGetParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationClaudeTagSlackWorkspaceGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationClaudeTagSlackWorkspaceUpdateParams struct {
	// The settings to change.
	Settings BetaClaudeTagSettingsUpdateParam `json:"settings,omitzero" api:"required"`
	// This endpoint is in beta: requests must send
	// `experimental-claude-tag-2026-10-21` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationClaudeTagSlackWorkspaceUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationClaudeTagSlackWorkspaceUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationClaudeTagSlackWorkspaceUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationClaudeTagSlackWorkspaceListParams struct {
	// When `true`, `settings` holds what applies in each workspace once the
	// organization's Slack defaults are taken into account, and every member of
	// `sources` names the level a value comes from. When `false` or left out,
	// `settings` holds what each workspace itself stores.
	IncludeInherited param.Opt[bool] `query:"include_inherited,omitzero" json:"-"`
	// Number of items to return per page, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Slack's identifier for one workspace, as in the workspace object's
	// `slack.workspace_id`, not its `id`. With it `data` holds that workspace alone,
	// or nothing when it is not set up.
	SlackWorkspaceID param.Opt[string] `query:"slack_workspace_id,omitzero" json:"-"`
	// This endpoint is in beta: requests must send
	// `experimental-claude-tag-2026-10-21` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationClaudeTagSlackWorkspaceListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationClaudeTagSlackWorkspaceListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
