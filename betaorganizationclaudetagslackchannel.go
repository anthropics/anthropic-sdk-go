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

// BetaOrganizationClaudeTagSlackChannelService contains methods and other services
// that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationClaudeTagSlackChannelService] method instead.
type BetaOrganizationClaudeTagSlackChannelService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationClaudeTagSlackChannelService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationClaudeTagSlackChannelService(opts ...option.RequestOption) (r BetaOrganizationClaudeTagSlackChannelService) {
	r = BetaOrganizationClaudeTagSlackChannelService{}
	r.Options = opts
	return
}

// Set a Slack channel up for Claude Tag and return it with its identifier. If the
// channel is already set up, whether in claude.ai, by Claude's first use in it or
// by an earlier request, the existing channel is returned, so the request can be
// repeated. `settings`, when sent, is applied as an update applies it. Before a
// channel is set up, Claude's Slack app confirms with Slack that the channel is in
// the workspace, is not archived and is visible to the app; a channel that fails
// the confirmation returns 404. Setting the channel up and applying `settings` are
// two steps: when the second is refused, the response is the refusal, the channel
// stays set up with every setting unset, and the error's
// `details.slack_channel_id` holds the channel's identifier.
func (r *BetaOrganizationClaudeTagSlackChannelService) New(ctx context.Context, params BetaOrganizationClaudeTagSlackChannelNewParams, opts ...option.RequestOption) (res *BetaClaudeTagSlackChannel, err error) {
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
	path := "v1/organizations/claude_tag/slack_channels"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_tag.slack_channels.create", "/v1/organizations/claude_tag/slack_channels"), requestconfig.WithSpanCreatedID("anthropic.organization.claude_tag.slack_channel.id", &res, func(created *BetaClaudeTagSlackChannel) string { return created.ID }))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get one Slack channel and its settings. By default `settings` holds what the
// channel itself stores: a setting that is null there is not set on the channel
// and is inherited. With `include_inherited=true`, `settings` holds what applies
// in the channel once the settings levels above it are taken into account, and the
// request fails with 503 rather than resolve against a settings level that could
// not be read. Either way `sources` names the settings level each value comes
// from, so a caller tells the channel's own value from an inherited one by its
// source, not by whether it is null.
func (r *BetaOrganizationClaudeTagSlackChannelService) Get(ctx context.Context, id string, params BetaOrganizationClaudeTagSlackChannelGetParams, opts ...option.RequestOption) (res *BetaClaudeTagSlackChannel, err error) {
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
	path := fmt.Sprintf("v1/organizations/claude_tag/slack_channels/%s", url.PathEscape(id))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_tag.slack_channels.retrieve", "/v1/organizations/claude_tag/slack_channels/{id}", "include_inherited"), requestconfig.WithSpanAttribute("anthropic.organization.claude_tag.slack_channel.id", id))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Change the settings a Slack channel stores. The update is partial: a setting
// left out is unchanged, and a setting sent as null is cleared, so the channel
// inherits it. If any value is refused, no setting in the request is stored. A
// group direct message takes no settings and returns 400.
func (r *BetaOrganizationClaudeTagSlackChannelService) Update(ctx context.Context, id string, params BetaOrganizationClaudeTagSlackChannelUpdateParams, opts ...option.RequestOption) (res *BetaClaudeTagSlackChannel, err error) {
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
	path := fmt.Sprintf("v1/organizations/claude_tag/slack_channels/%s", url.PathEscape(id))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_tag.slack_channels.update", "/v1/organizations/claude_tag/slack_channels/{id}"), requestconfig.WithSpanAttribute("anthropic.organization.claude_tag.slack_channel.id", id))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List the set-up Slack channels of one Slack workspace, newest first by
// `created_at`. A channel that is not set up is not listed: Slack holds the full
// list of a workspace's channels. Pass `slack_channel_id` as well to look one
// channel up by its Slack identifiers: `data` then holds that channel, or is empty
// when it is not set up, and `limit` and `page` are ignored. A workspace that the
// organization owns only through a Slack Enterprise Grid returns 404 until it is
// set up. With `include_inherited=true` each channel's `settings` holds what
// applies in it, and the request fails with 503 rather than resolve against a
// settings level that could not be read. A `page` cursor resumes only a request
// with the same `include_inherited`.
func (r *BetaOrganizationClaudeTagSlackChannelService) List(ctx context.Context, params BetaOrganizationClaudeTagSlackChannelListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaClaudeTagSlackChannel], err error) {
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
	path := "v1/organizations/claude_tag/slack_channels"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_tag.slack_channels.list", "/v1/organizations/claude_tag/slack_channels", "include_inherited", "limit", "page", "slack_channel_id", "slack_workspace_id"))
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

// List the set-up Slack channels of one Slack workspace, newest first by
// `created_at`. A channel that is not set up is not listed: Slack holds the full
// list of a workspace's channels. Pass `slack_channel_id` as well to look one
// channel up by its Slack identifiers: `data` then holds that channel, or is empty
// when it is not set up, and `limit` and `page` are ignored. A workspace that the
// organization owns only through a Slack Enterprise Grid returns 404 until it is
// set up. With `include_inherited=true` each channel's `settings` holds what
// applies in it, and the request fails with 503 rather than resolve against a
// settings level that could not be read. A `page` cursor resumes only a request
// with the same `include_inherited`.
func (r *BetaOrganizationClaudeTagSlackChannelService) ListAutoPaging(ctx context.Context, params BetaOrganizationClaudeTagSlackChannelListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaClaudeTagSlackChannel] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Remove a Slack channel from Claude Tag. This deletes the channel's stored
// settings and its Bundle Targets at once, and then, in the background, the
// channel's sessions and their transcripts, its memory, its routines and the
// artifacts published from it. None of it can be restored. The channel in Slack
// and its messages are unchanged. The request is refused with 400 and
// `error.details.error_code` `slack_channel_enabled` while `enabled` resolves to
// `true` for the channel, as a read with `include_inherited=true` returns it: set
// `enabled` to `false` on the channel first. This guards against removing a
// channel in use by mistake; it is not a way to keep Claude out, because the
// removal deletes the channel's own `enabled: false` with its other settings,
// after which the channel inherits again and Claude's next use of it can set it up
// anew with a new identifier. The identifier returns 404 as soon as the request
// returns, so a repeated removal returns 404. A group direct message returns 400,
// and so does a channel that holds configuration this API does not manage.
func (r *BetaOrganizationClaudeTagSlackChannelService) Delete(ctx context.Context, id string, body BetaOrganizationClaudeTagSlackChannelDeleteParams, opts ...option.RequestOption) (res *BetaClaudeTagDeletedSlackChannel, err error) {
	if len(body.Betas) > 0 {
		headerValues := make([]string, len(body.Betas))
		for i, v := range body.Betas {
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
	path := fmt.Sprintf("v1/organizations/claude_tag/slack_channels/%s", url.PathEscape(id))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_tag.slack_channels.delete", "/v1/organizations/claude_tag/slack_channels/{id}"), requestconfig.WithSpanAttribute("anthropic.organization.claude_tag.slack_channel.id", id))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Confirms that a Slack channel was removed from Claude Tag.
type BetaClaudeTagDeletedSlackChannel struct {
	// The removed channel's identifier.
	ID   string                       `json:"id" api:"required"`
	Type constant.SlackChannelDeleted `json:"type" default:"slack_channel_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagDeletedSlackChannel) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagDeletedSlackChannel) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Slack channel set up for Claude Tag, with the settings it stores or, with
// `include_inherited`, the settings that apply in it.
type BetaClaudeTagSlackChannel struct {
	// The channel's identifier. Starts with `slch_`. It does not change for as long as
	// the channel is set up, with one exception: a channel that is not configured gets
	// a new one when claude.ai replaces its record, which happens only while the
	// organization's read-only mode is off. A channel that is removed and set up again
	// also has a new one. The old one keeps returning 404; look the channel up by its
	// Slack identifiers for the new one.
	ID string `json:"id" api:"required"`
	// Time the channel was set up, by whichever route. It does not change for as long
	// as `id` does not.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// What the channel itself stores or, with `include_inherited`, what applies in it.
	// `sources` says which.
	Settings BetaClaudeTagSettings `json:"settings" api:"required"`
	// Slack's own identifiers for the channel.
	Slack BetaClaudeTagSlackChannelIdentifiers `json:"slack" api:"required"`
	// For each setting, the settings level its value in `settings` comes from.
	Sources BetaClaudeTagSettingSources `json:"sources" api:"required"`
	Type    constant.SlackChannel       `json:"type" default:"slack_channel"`
	// Time Claude Tag's record of the channel last changed. It moves on every change
	// to a setting, and can also move when something no response holds changes, so a
	// later value says to read the channel again, not that a setting changed.
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
func (r BetaClaudeTagSlackChannel) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagSlackChannel) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationClaudeTagSlackChannelNewParams struct {
	// Slack's own identifiers for the channel to set up. The organization must own the
	// workspace.
	Slack BetaClaudeTagSlackChannelIdentifiersParam `json:"slack,omitzero" api:"required"`
	// Settings to store on the channel, applied as Update Slack Channel applies them.
	Settings BetaClaudeTagSettingsCreateParam `json:"settings,omitzero"`
	// This endpoint is in beta: requests must send
	// `experimental-claude-tag-2026-10-21` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationClaudeTagSlackChannelNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationClaudeTagSlackChannelNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationClaudeTagSlackChannelNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationClaudeTagSlackChannelGetParams struct {
	// When `true`, `settings` holds what applies in the channel once the settings
	// levels above it are taken into account, and every member of `sources` names the
	// level a value comes from. When `false` or left out, `settings` holds what the
	// channel itself stores.
	IncludeInherited param.Opt[bool] `query:"include_inherited,omitzero" json:"-"`
	// This endpoint is in beta: requests must send
	// `experimental-claude-tag-2026-10-21` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationClaudeTagSlackChannelGetParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationClaudeTagSlackChannelGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationClaudeTagSlackChannelUpdateParams struct {
	// The settings to change.
	Settings BetaClaudeTagSettingsUpdateParam `json:"settings,omitzero" api:"required"`
	// This endpoint is in beta: requests must send
	// `experimental-claude-tag-2026-10-21` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationClaudeTagSlackChannelUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationClaudeTagSlackChannelUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationClaudeTagSlackChannelUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationClaudeTagSlackChannelListParams struct {
	// The Slack workspace whose channels to list: Slack's identifier for it, as in the
	// channel object's `slack.workspace_id`, not the workspace object's `id`.
	SlackWorkspaceID string `query:"slack_workspace_id" api:"required" json:"-"`
	// When `true`, `settings` holds what applies in each channel once the settings
	// levels above it are taken into account, and every member of `sources` names the
	// level a value comes from. When `false` or left out, `settings` holds what each
	// channel itself stores.
	IncludeInherited param.Opt[bool] `query:"include_inherited,omitzero" json:"-"`
	// Number of items to return per page, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Slack's identifier for one channel of that workspace, as in the channel object's
	// `slack.channel_id`. It is not the channel object's `id`. With it `data` holds
	// that channel alone, or nothing when the channel is not set up.
	SlackChannelID param.Opt[string] `query:"slack_channel_id,omitzero" json:"-"`
	// This endpoint is in beta: requests must send
	// `experimental-claude-tag-2026-10-21` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationClaudeTagSlackChannelListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationClaudeTagSlackChannelListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationClaudeTagSlackChannelDeleteParams struct {
	// This endpoint is in beta: requests must send
	// `experimental-claude-tag-2026-10-21` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
