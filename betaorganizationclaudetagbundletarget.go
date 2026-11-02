package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/apiquery"
	shimjson "github.com/anthropics/anthropic-sdk-go/v2/internal/encoding/json"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/paramutil"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// BetaOrganizationClaudeTagBundleTargetService contains methods and other services
// that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationClaudeTagBundleTargetService] method instead.
type BetaOrganizationClaudeTagBundleTargetService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationClaudeTagBundleTargetService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationClaudeTagBundleTargetService(opts ...option.RequestOption) (r BetaOrganizationClaudeTagBundleTargetService) {
	r = BetaOrganizationClaudeTagBundleTargetService{}
	r.Options = opts
	return
}

// Attach a Bundle to a Slack channel or a Slack workspace. A Bundle attached to a
// workspace can be used in the workspace and in every channel in it. Attaching a
// Bundle to a channel or workspace it is already attached to returns the existing
// Bundle Target. The channel or workspace must be set up: attaching never sets a
// channel up. A group direct message takes no Bundle and returns 400.
func (r *BetaOrganizationClaudeTagBundleTargetService) New(ctx context.Context, bundleID string, params BetaOrganizationClaudeTagBundleTargetNewParams, opts ...option.RequestOption) (res *BetaClaudeTagBundleTargetUnion, err error) {
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
	if bundleID == "" {
		err = errors.New("missing required bundle_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/claude_tag/bundles/%s/targets", url.PathEscape(bundleID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_tag.bundles.targets.create", "/v1/organizations/claude_tag/bundles/{bundle_id}/targets"), requestconfig.WithSpanAttribute("anthropic.organization.claude_tag.bundle.id", bundleID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get one Bundle Target. Returns 404 when `target_id` is the identifier of another
// Bundle's Bundle Target, or the Bundle Target was detached.
func (r *BetaOrganizationClaudeTagBundleTargetService) Get(ctx context.Context, targetID string, params BetaOrganizationClaudeTagBundleTargetGetParams, opts ...option.RequestOption) (res *BetaClaudeTagBundleTargetUnion, err error) {
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
	if params.BundleID == "" {
		err = errors.New("missing required bundle_id parameter")
		return nil, err
	}
	if targetID == "" {
		err = errors.New("missing required target_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/claude_tag/bundles/%s/targets/%s", url.PathEscape(params.BundleID), url.PathEscape(targetID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_tag.bundles.targets.retrieve", "/v1/organizations/claude_tag/bundles/{bundle_id}/targets/{target_id}"), requestconfig.WithSpanAttribute("anthropic.organization.claude_tag.bundle.id", params.BundleID), requestconfig.WithSpanAttribute("anthropic.organization.claude_tag.bundle.target.id", targetID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List the Slack channels and Slack workspaces a Bundle is attached to, newest
// first by `created_at`. Only channels and workspaces the organization owns are
// listed; a Bundle that is attached to anything else has no Bundle Target for it.
func (r *BetaOrganizationClaudeTagBundleTargetService) List(ctx context.Context, bundleID string, params BetaOrganizationClaudeTagBundleTargetListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaClaudeTagBundleTargetUnion], err error) {
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
	if bundleID == "" {
		err = errors.New("missing required bundle_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/claude_tag/bundles/%s/targets", url.PathEscape(bundleID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_tag.bundles.targets.list", "/v1/organizations/claude_tag/bundles/{bundle_id}/targets", "limit", "page"), requestconfig.WithSpanAttribute("anthropic.organization.claude_tag.bundle.id", bundleID))
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

// List the Slack channels and Slack workspaces a Bundle is attached to, newest
// first by `created_at`. Only channels and workspaces the organization owns are
// listed; a Bundle that is attached to anything else has no Bundle Target for it.
func (r *BetaOrganizationClaudeTagBundleTargetService) ListAutoPaging(ctx context.Context, bundleID string, params BetaOrganizationClaudeTagBundleTargetListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaClaudeTagBundleTargetUnion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, bundleID, params, opts...))
}

// Detach a Bundle from a Slack channel or a Slack workspace. A repeated detach
// returns 404.
func (r *BetaOrganizationClaudeTagBundleTargetService) Delete(ctx context.Context, targetID string, params BetaOrganizationClaudeTagBundleTargetDeleteParams, opts ...option.RequestOption) (res *BetaClaudeTagDeletedBundleTarget, err error) {
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
	if params.BundleID == "" {
		err = errors.New("missing required bundle_id parameter")
		return nil, err
	}
	if targetID == "" {
		err = errors.New("missing required target_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/claude_tag/bundles/%s/targets/%s", url.PathEscape(params.BundleID), url.PathEscape(targetID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_tag.bundles.targets.delete", "/v1/organizations/claude_tag/bundles/{bundle_id}/targets/{target_id}"), requestconfig.WithSpanAttribute("anthropic.organization.claude_tag.bundle.id", params.BundleID), requestconfig.WithSpanAttribute("anthropic.organization.claude_tag.bundle.target.id", targetID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// BetaClaudeTagBundleTargetUnion contains all possible properties and values from
// [BetaClaudeTagSlackChannelBundleTarget],
// [BetaClaudeTagSlackWorkspaceBundleTarget].
//
// Use the [BetaClaudeTagBundleTargetUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaClaudeTagBundleTargetUnion struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	// This field is from variant [BetaClaudeTagSlackChannelBundleTarget].
	SlackChannelID string `json:"slack_channel_id"`
	// Any of "slack_channel", "slack_workspace".
	Type string `json:"type"`
	// This field is from variant [BetaClaudeTagSlackWorkspaceBundleTarget].
	SlackWorkspaceID string `json:"slack_workspace_id"`
	JSON             struct {
		ID               respjson.Field
		CreatedAt        respjson.Field
		SlackChannelID   respjson.Field
		Type             respjson.Field
		SlackWorkspaceID respjson.Field
		raw              string
	} `json:"-"`
}

// anyBetaClaudeTagBundleTarget is implemented by each variant of
// [BetaClaudeTagBundleTargetUnion] to add type safety for the return type of
// [BetaClaudeTagBundleTargetUnion.AsAny]
type anyBetaClaudeTagBundleTarget interface {
	implBetaClaudeTagBundleTargetUnion()
}

func (BetaClaudeTagSlackChannelBundleTarget) implBetaClaudeTagBundleTargetUnion()   {}
func (BetaClaudeTagSlackWorkspaceBundleTarget) implBetaClaudeTagBundleTargetUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaClaudeTagBundleTargetUnion.AsAny().(type) {
//	case anthropic.BetaClaudeTagSlackChannelBundleTarget:
//	case anthropic.BetaClaudeTagSlackWorkspaceBundleTarget:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaClaudeTagBundleTargetUnion) AsAny() anyBetaClaudeTagBundleTarget {
	switch u.Type {
	case "slack_channel":
		return u.AsSlackChannel()
	case "slack_workspace":
		return u.AsSlackWorkspace()
	}
	return nil
}

func (u BetaClaudeTagBundleTargetUnion) AsSlackChannel() (v BetaClaudeTagSlackChannelBundleTarget) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaClaudeTagBundleTargetUnion) AsSlackWorkspace() (v BetaClaudeTagSlackWorkspaceBundleTarget) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaClaudeTagBundleTargetUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaClaudeTagBundleTargetUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func BetaClaudeTagBundleTargetCreateParamOfSlackChannel(slackChannelID string) BetaClaudeTagBundleTargetCreateUnionParam {
	var slackChannel BetaClaudeTagSlackChannelBundleTargetCreateParam
	slackChannel.SlackChannelID = slackChannelID
	return BetaClaudeTagBundleTargetCreateUnionParam{OfSlackChannel: &slackChannel}
}

func BetaClaudeTagBundleTargetCreateParamOfSlackWorkspace(slackWorkspaceID string) BetaClaudeTagBundleTargetCreateUnionParam {
	var slackWorkspace BetaClaudeTagSlackWorkspaceBundleTargetCreateParam
	slackWorkspace.SlackWorkspaceID = slackWorkspaceID
	return BetaClaudeTagBundleTargetCreateUnionParam{OfSlackWorkspace: &slackWorkspace}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaClaudeTagBundleTargetCreateUnionParam struct {
	OfSlackChannel   *BetaClaudeTagSlackChannelBundleTargetCreateParam   `json:",omitzero,inline"`
	OfSlackWorkspace *BetaClaudeTagSlackWorkspaceBundleTargetCreateParam `json:",omitzero,inline"`
	paramUnion
}

func (u BetaClaudeTagBundleTargetCreateUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfSlackChannel, u.OfSlackWorkspace)
}
func (u *BetaClaudeTagBundleTargetCreateUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaClaudeTagBundleTargetCreateUnionParam) asAny() any {
	if !param.IsOmitted(u.OfSlackChannel) {
		return u.OfSlackChannel
	} else if !param.IsOmitted(u.OfSlackWorkspace) {
		return u.OfSlackWorkspace
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaClaudeTagBundleTargetCreateUnionParam) GetSlackChannelID() *string {
	if vt := u.OfSlackChannel; vt != nil {
		return &vt.SlackChannelID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaClaudeTagBundleTargetCreateUnionParam) GetSlackWorkspaceID() *string {
	if vt := u.OfSlackWorkspace; vt != nil {
		return &vt.SlackWorkspaceID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaClaudeTagBundleTargetCreateUnionParam) GetType() *string {
	if vt := u.OfSlackChannel; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	} else if vt := u.OfSlackWorkspace; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaClaudeTagBundleTargetCreateUnionParam](
		"type",
		apijson.Discriminator[BetaClaudeTagSlackChannelBundleTargetCreateParam]("slack_channel"),
		apijson.Discriminator[BetaClaudeTagSlackWorkspaceBundleTargetCreateParam]("slack_workspace"),
	)
}

// Confirms that a Bundle was detached.
type BetaClaudeTagDeletedBundleTarget struct {
	// The identifier the Bundle Target had.
	ID   string                       `json:"id" api:"required"`
	Type constant.BundleTargetDeleted `json:"type" default:"bundle_target_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagDeletedBundleTarget) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagDeletedBundleTarget) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// One Bundle attached to one Slack channel. The Bundle can be used in that channel
// only.
type BetaClaudeTagSlackChannelBundleTarget struct {
	// The Bundle Target's identifier. Starts with `btgt_`. No two Bundle Targets have
	// the same one, and the same Bundle and target always give the same one, also
	// after a detach and a new attach. It is opaque.
	ID string `json:"id" api:"required"`
	// Time the Bundle was attached.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// The identifier of the Slack channel the Bundle is attached to, starting with
	// `slch_`. It is the channel object's `id`, not Slack's `channel_id`.
	SlackChannelID string                `json:"slack_channel_id" api:"required"`
	Type           constant.SlackChannel `json:"type" default:"slack_channel"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		CreatedAt      respjson.Field
		SlackChannelID respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagSlackChannelBundleTarget) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagSlackChannelBundleTarget) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Attach the Bundle to a Slack channel.
//
// The properties SlackChannelID, Type are required.
type BetaClaudeTagSlackChannelBundleTargetCreateParam struct {
	// The identifier of the Slack channel to attach the Bundle to, starting with
	// `slch_`. It is the channel object's `id`, not Slack's `channel_id`.
	SlackChannelID string `json:"slack_channel_id" api:"required"`
	// This field can be elided, and will marshal its zero value as "slack_channel".
	Type constant.SlackChannel `json:"type" default:"slack_channel"`
	paramObj
}

func (r BetaClaudeTagSlackChannelBundleTargetCreateParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaClaudeTagSlackChannelBundleTargetCreateParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaClaudeTagSlackChannelBundleTargetCreateParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// One Bundle attached to one Slack workspace. The Bundle can be used in the
// workspace and in every channel in it.
type BetaClaudeTagSlackWorkspaceBundleTarget struct {
	// The Bundle Target's identifier. Starts with `btgt_`. No two Bundle Targets have
	// the same one, and the same Bundle and target always give the same one, also
	// after a detach and a new attach. It is opaque.
	ID string `json:"id" api:"required"`
	// Time the Bundle was attached.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// The identifier of the Slack workspace the Bundle is attached to, starting with
	// `slws_`.
	SlackWorkspaceID string                  `json:"slack_workspace_id" api:"required"`
	Type             constant.SlackWorkspace `json:"type" default:"slack_workspace"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		CreatedAt        respjson.Field
		SlackWorkspaceID respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeTagSlackWorkspaceBundleTarget) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeTagSlackWorkspaceBundleTarget) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Attach the Bundle to a Slack workspace.
//
// The properties SlackWorkspaceID, Type are required.
type BetaClaudeTagSlackWorkspaceBundleTargetCreateParam struct {
	// The identifier of the Slack workspace to attach the Bundle to, starting with
	// `slws_`.
	SlackWorkspaceID string `json:"slack_workspace_id" api:"required"`
	// This field can be elided, and will marshal its zero value as "slack_workspace".
	Type constant.SlackWorkspace `json:"type" default:"slack_workspace"`
	paramObj
}

func (r BetaClaudeTagSlackWorkspaceBundleTargetCreateParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaClaudeTagSlackWorkspaceBundleTargetCreateParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaClaudeTagSlackWorkspaceBundleTargetCreateParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationClaudeTagBundleTargetNewParams struct {
	Target BetaClaudeTagBundleTargetCreateUnionParam
	// This endpoint is in beta: requests must send
	// `experimental-claude-tag-2026-10-21` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationClaudeTagBundleTargetNewParams) MarshalJSON() (data []byte, err error) {
	return shimjson.Marshal(r.Target)
}
func (r *BetaOrganizationClaudeTagBundleTargetNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationClaudeTagBundleTargetGetParams struct {
	// The Bundle's identifier, starting with `capp_`.
	BundleID string `path:"bundle_id" api:"required" json:"-"`
	// This endpoint is in beta: requests must send
	// `experimental-claude-tag-2026-10-21` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationClaudeTagBundleTargetListParams struct {
	// Number of items to return per page, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// This endpoint is in beta: requests must send
	// `experimental-claude-tag-2026-10-21` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationClaudeTagBundleTargetListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationClaudeTagBundleTargetListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationClaudeTagBundleTargetDeleteParams struct {
	// The Bundle's identifier, starting with `capp_`.
	BundleID string `path:"bundle_id" api:"required" json:"-"`
	// This endpoint is in beta: requests must send
	// `experimental-claude-tag-2026-10-21` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
