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
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// BetaOrganizationClaudeCodeManagedSettingRBACGroupService contains methods and
// other services that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationClaudeCodeManagedSettingRBACGroupService] method
// instead.
type BetaOrganizationClaudeCodeManagedSettingRBACGroupService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationClaudeCodeManagedSettingRBACGroupService generates a new
// service that applies the given options to each request. These options are
// applied after the parent client's options (if there is one), and before any
// request-specific options.
func NewBetaOrganizationClaudeCodeManagedSettingRBACGroupService(opts ...option.RequestOption) (r BetaOrganizationClaudeCodeManagedSettingRBACGroupService) {
	r = BetaOrganizationClaudeCodeManagedSettingRBACGroupService{}
	r.Options = opts
	return
}

// Retrieve an RBAC Group's Claude Code managed settings.
//
// **Accepted credentials:** an organization-level Admin API key with the
// `read:claude_code_managed_settings` scope.
//
// Every request must include the beta header
// `anthropic-beta: ce-claude-code-managed-settings-2026-10-06`. A request without
// it, or from an organization not enabled for this beta, returns `404`. The Claude
// Code Managed Settings API is in beta and is available to Claude Enterprise
// organizations only.
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeBetaCMEKFeatureUnsupported] (HTTP 400)
//   - [ErrorCodeCMEKKeyDisabled] (HTTP 400)
//   - [ErrorCodeBetaRBACGroupManagedSettingsNotFound] (HTTP 404)
func (r *BetaOrganizationClaudeCodeManagedSettingRBACGroupService) Get(ctx context.Context, rbacGroupID string, query BetaOrganizationClaudeCodeManagedSettingRBACGroupGetParams, opts ...option.RequestOption) (res *BetaClaudeCodeRBACGroupManagedSettings, err error) {
	if len(query.Betas) > 0 {
		headerValues := make([]string, len(query.Betas))
		for i, v := range query.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-claude-code-managed-settings-2026-10-06")}, opts...)
	if rbacGroupID == "" {
		err = errors.New("missing required rbac_group_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/claude_code/managed_settings/rbac_groups/%s", requestconfig.EscapePathSegment(rbacGroupID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_code.managed_settings.rbac_groups.retrieve", "/v1/organizations/claude_code/managed_settings/rbac_groups/{rbac_group_id}"), requestconfig.WithSpanAttribute("anthropic.organization.claude_code.managed_setting.rbac_group.id", rbacGroupID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List the RBAC Groups that have Claude Code managed settings, highest precedence
// first.
//
// Each entry gives the group's status and timestamps, not its settings. The list
// includes deleted groups that still have settings, and leaves out the default
// settings. If the order changes between pages, the request returns `409`: list
// again from the first page.
//
// **Accepted credentials:** an organization-level Admin API key with the
// `read:claude_code_managed_settings` scope.
//
// Every request must include the beta header
// `anthropic-beta: ce-claude-code-managed-settings-2026-10-06`. A request without
// it, or from an organization not enabled for this beta, returns `404`. The Claude
// Code Managed Settings API is in beta and is available to Claude Enterprise
// organizations only.
func (r *BetaOrganizationClaudeCodeManagedSettingRBACGroupService) List(ctx context.Context, params BetaOrganizationClaudeCodeManagedSettingRBACGroupListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaClaudeCodeRBACGroupManagedSettingsSummary], err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-claude-code-managed-settings-2026-10-06"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/claude_code/managed_settings/rbac_groups"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_code.managed_settings.rbac_groups.list", "/v1/organizations/claude_code/managed_settings/rbac_groups", "limit", "page"))
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

// List the RBAC Groups that have Claude Code managed settings, highest precedence
// first.
//
// Each entry gives the group's status and timestamps, not its settings. The list
// includes deleted groups that still have settings, and leaves out the default
// settings. If the order changes between pages, the request returns `409`: list
// again from the first page.
//
// **Accepted credentials:** an organization-level Admin API key with the
// `read:claude_code_managed_settings` scope.
//
// Every request must include the beta header
// `anthropic-beta: ce-claude-code-managed-settings-2026-10-06`. A request without
// it, or from an organization not enabled for this beta, returns `404`. The Claude
// Code Managed Settings API is in beta and is available to Claude Enterprise
// organizations only.
func (r *BetaOrganizationClaudeCodeManagedSettingRBACGroupService) ListAutoPaging(ctx context.Context, params BetaOrganizationClaudeCodeManagedSettingRBACGroupListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaClaudeCodeRBACGroupManagedSettingsSummary] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Delete an RBAC Group's Claude Code managed settings and remove the group from
// `rbac_group_order`.
//
// A group with no settings returns `404` with the error code
// `rbac_group_managed_settings_not_found`. That code on a retry of a delete that
// returned `503` means the settings are gone.
//
// **Accepted credentials:** an organization-level Admin API key with the
// `write:claude_code_managed_settings` scope.
//
// Every request must include the beta header
// `anthropic-beta: ce-claude-code-managed-settings-2026-10-06`. A request without
// it, or from an organization not enabled for this beta, returns `404`. The Claude
// Code Managed Settings API is in beta and is available to Claude Enterprise
// organizations only.
//
// A write that returns `503` may have applied. Retrying it is safe.
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeCMEKKeyDisabled] (HTTP 400)
//   - [ErrorCodeBetaSettingsHistoryLimitExceeded] (HTTP 400)
//   - [ErrorCodeSandboxEgressRestricted] (HTTP 403)
//   - [ErrorCodeBetaRBACGroupManagedSettingsNotFound] (HTTP 404)
func (r *BetaOrganizationClaudeCodeManagedSettingRBACGroupService) Delete(ctx context.Context, rbacGroupID string, body BetaOrganizationClaudeCodeManagedSettingRBACGroupDeleteParams, opts ...option.RequestOption) (res *BetaDeletedClaudeCodeRBACGroupManagedSettings, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-claude-code-managed-settings-2026-10-06")}, opts...)
	if rbacGroupID == "" {
		err = errors.New("missing required rbac_group_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/claude_code/managed_settings/rbac_groups/%s", requestconfig.EscapePathSegment(rbacGroupID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_code.managed_settings.rbac_groups.delete", "/v1/organizations/claude_code/managed_settings/rbac_groups/{rbac_group_id}"), requestconfig.WithSpanAttribute("anthropic.organization.claude_code.managed_setting.rbac_group.id", rbacGroupID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Create or replace an RBAC Group's Claude Code managed settings.
//
// New settings join the end of `rbac_group_order`, and existing settings keep
// their place. Until the group is moved up, its members who also belong to an
// earlier group keep that group's settings. In `organization` mode a write changes
// no member's settings until the organization switches to `rbac_groups` mode. Once
// a group has settings, that switch no longer copies the organization settings
// into default settings, so members in no RBAC Group in `rbac_group_order` receive
// no settings unless you set the default settings before the switch.
//
// Get the group's ID from the RBAC Groups API. One key can hold `read:rbac_groups`
// and `write:claude_code_managed_settings`.
//
// **Accepted credentials:** an organization-level Admin API key with the
// `write:claude_code_managed_settings` scope.
//
// Every request must include the beta header
// `anthropic-beta: ce-claude-code-managed-settings-2026-10-06`. A request without
// it, or from an organization not enabled for this beta, returns `404`. The Claude
// Code Managed Settings API is in beta and is available to Claude Enterprise
// organizations only.
//
// A write that returns `503` may have applied. Retrying it is safe.
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeBetaCMEKFeatureUnsupported] (HTTP 400)
//   - [ErrorCodeCMEKKeyDisabled] (HTTP 400)
//   - [ErrorCodeBetaRBACGroupOrderLimitExceeded] (HTTP 400)
//   - [ErrorCodeBetaRequestBodyTimeout] (HTTP 400)
//   - [ErrorCodeBetaSettingsDocumentLimitExceeded] (HTTP 400)
//   - [ErrorCodeBetaSettingsHistoryLimitExceeded] (HTTP 400)
//   - [ErrorCodeSandboxEgressRestricted] (HTTP 403)
//   - [ErrorCodeBetaRBACGroupNotFound] (HTTP 404)
func (r *BetaOrganizationClaudeCodeManagedSettingRBACGroupService) Set(ctx context.Context, rbacGroupID string, params BetaOrganizationClaudeCodeManagedSettingRBACGroupSetParams, opts ...option.RequestOption) (res *BetaClaudeCodeRBACGroupManagedSettings, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-claude-code-managed-settings-2026-10-06")}, opts...)
	if rbacGroupID == "" {
		err = errors.New("missing required rbac_group_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/claude_code/managed_settings/rbac_groups/%s", requestconfig.EscapePathSegment(rbacGroupID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_code.managed_settings.rbac_groups.set", "/v1/organizations/claude_code/managed_settings/rbac_groups/{rbac_group_id}"), requestconfig.WithSpanAttribute("anthropic.organization.claude_code.managed_setting.rbac_group.id", rbacGroupID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type BetaClaudeCodeRBACGroupManagedSettings struct {
	// When the settings document this group uses was first stored, in RFC 3339 format.
	// This can be earlier than when the group got its settings, and it changes if the
	// group gets a different document.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// The RBAC Group's ID, prefixed `rbac_group_`.
	RBACGroupID string `json:"rbac_group_id" api:"required"`
	// Whether the group is still one of the organization's RBAC Groups. Its `type` is
	// `active` if it is, `deleted` if it is not, and `unknown` if that cannot be
	// checked right now. A deleted group's settings stay in `rbac_group_order` but
	// match no member.
	RBACGroupStatus BetaClaudeCodeRBACGroupStatusUnion `json:"rbac_group_status" api:"required"`
	// The group's settings document, or `null` if the stored document is damaged and
	// can no longer be read. Replacing it fixes this.
	Settings map[string]any `json:"settings" api:"required"`
	// Always `claude_code_rbac_group_managed_settings`.
	Type constant.ClaudeCodeRBACGroupManagedSettings `json:"type" default:"claude_code_rbac_group_managed_settings"`
	// When this group's settings or its position in the order last changed, in RFC
	// 3339 format. Its position changes whenever the number of groups ahead of it
	// changes.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CreatedAt       respjson.Field
		RBACGroupID     respjson.Field
		RBACGroupStatus respjson.Field
		Settings        respjson.Field
		Type            respjson.Field
		UpdatedAt       respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeCodeRBACGroupManagedSettings) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeCodeRBACGroupManagedSettings) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaClaudeCodeRBACGroupManagedSettingsSummary struct {
	// When the settings document this group uses was first stored, in RFC 3339 format.
	// This can be earlier than when the group got its settings, and it changes if the
	// group gets a different document.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// The RBAC Group's ID, prefixed `rbac_group_`.
	RBACGroupID string `json:"rbac_group_id" api:"required"`
	// Whether the group is still one of the organization's RBAC Groups. Its `type` is
	// `active` if it is, `deleted` if it is not, and `unknown` if that cannot be
	// checked right now. A deleted group's settings stay in `rbac_group_order` but
	// match no member.
	RBACGroupStatus BetaClaudeCodeRBACGroupStatusUnion `json:"rbac_group_status" api:"required"`
	// Always `claude_code_rbac_group_managed_settings_summary`.
	Type constant.ClaudeCodeRBACGroupManagedSettingsSummary `json:"type" default:"claude_code_rbac_group_managed_settings_summary"`
	// When this group's settings or its position in the order last changed, in RFC
	// 3339 format. Its position changes whenever the number of groups ahead of it
	// changes.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CreatedAt       respjson.Field
		RBACGroupID     respjson.Field
		RBACGroupStatus respjson.Field
		Type            respjson.Field
		UpdatedAt       respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeCodeRBACGroupManagedSettingsSummary) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeCodeRBACGroupManagedSettingsSummary) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaClaudeCodeRBACGroupStatusUnion contains all possible properties and values
// from [BetaClaudeCodeRBACGroupStatusActive],
// [BetaClaudeCodeRBACGroupStatusDeleted], [BetaClaudeCodeRBACGroupStatusUnknown].
//
// Use the [BetaClaudeCodeRBACGroupStatusUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaClaudeCodeRBACGroupStatusUnion struct {
	// Any of "active", "deleted", "unknown".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaClaudeCodeRBACGroupStatus is implemented by each variant of
// [BetaClaudeCodeRBACGroupStatusUnion] to add type safety for the return type of
// [BetaClaudeCodeRBACGroupStatusUnion.AsAny]
type anyBetaClaudeCodeRBACGroupStatus interface {
	implBetaClaudeCodeRBACGroupStatusUnion()
}

func (BetaClaudeCodeRBACGroupStatusActive) implBetaClaudeCodeRBACGroupStatusUnion()  {}
func (BetaClaudeCodeRBACGroupStatusDeleted) implBetaClaudeCodeRBACGroupStatusUnion() {}
func (BetaClaudeCodeRBACGroupStatusUnknown) implBetaClaudeCodeRBACGroupStatusUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaClaudeCodeRBACGroupStatusUnion.AsAny().(type) {
//	case anthropic.BetaClaudeCodeRBACGroupStatusActive:
//	case anthropic.BetaClaudeCodeRBACGroupStatusDeleted:
//	case anthropic.BetaClaudeCodeRBACGroupStatusUnknown:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaClaudeCodeRBACGroupStatusUnion) AsAny() anyBetaClaudeCodeRBACGroupStatus {
	switch u.Type {
	case "active":
		return u.AsActive()
	case "deleted":
		return u.AsDeleted()
	case "unknown":
		return u.AsUnknown()
	}
	return nil
}

func (u BetaClaudeCodeRBACGroupStatusUnion) AsActive() (v BetaClaudeCodeRBACGroupStatusActive) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaClaudeCodeRBACGroupStatusUnion) AsDeleted() (v BetaClaudeCodeRBACGroupStatusDeleted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaClaudeCodeRBACGroupStatusUnion) AsUnknown() (v BetaClaudeCodeRBACGroupStatusUnknown) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaClaudeCodeRBACGroupStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaClaudeCodeRBACGroupStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaClaudeCodeRBACGroupStatusActive struct {
	// Always `active`.
	Type constant.Active `json:"type" default:"active"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeCodeRBACGroupStatusActive) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeCodeRBACGroupStatusActive) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaClaudeCodeRBACGroupStatusDeleted struct {
	// Always `deleted`.
	Type constant.Deleted `json:"type" default:"deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeCodeRBACGroupStatusDeleted) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeCodeRBACGroupStatusDeleted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaClaudeCodeRBACGroupStatusUnknown struct {
	// Always `unknown`.
	Type constant.Unknown `json:"type" default:"unknown"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeCodeRBACGroupStatusUnknown) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeCodeRBACGroupStatusUnknown) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaDeletedClaudeCodeRBACGroupManagedSettings struct {
	// The RBAC Group's ID, prefixed `rbac_group_`.
	RBACGroupID string `json:"rbac_group_id" api:"required"`
	// Always `claude_code_rbac_group_managed_settings_deleted`.
	Type constant.ClaudeCodeRBACGroupManagedSettingsDeleted `json:"type" default:"claude_code_rbac_group_managed_settings_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		RBACGroupID respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDeletedClaudeCodeRBACGroupManagedSettings) RawJSON() string { return r.JSON.raw }
func (r *BetaDeletedClaudeCodeRBACGroupManagedSettings) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationClaudeCodeManagedSettingRBACGroupGetParams struct {
	// This endpoint is in beta: requests must send
	// `ce-claude-code-managed-settings-2026-10-06` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationClaudeCodeManagedSettingRBACGroupListParams struct {
	// Optionally set to the `next_page` token from the previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `100`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// This endpoint is in beta: requests must send
	// `ce-claude-code-managed-settings-2026-10-06` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes
// [BetaOrganizationClaudeCodeManagedSettingRBACGroupListParams]'s query parameters
// as `url.Values`.
func (r BetaOrganizationClaudeCodeManagedSettingRBACGroupListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationClaudeCodeManagedSettingRBACGroupDeleteParams struct {
	// This endpoint is in beta: requests must send
	// `ce-claude-code-managed-settings-2026-10-06` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationClaudeCodeManagedSettingRBACGroupSetParams struct {
	// The whole settings document, a JSON object of Claude Code managed settings. It
	// replaces any stored document. At most 1 MiB as the server stores it, which can
	// be larger than the request body.
	Settings map[string]any `json:"settings,omitzero" api:"required"`
	// This endpoint is in beta: requests must send
	// `ce-claude-code-managed-settings-2026-10-06` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationClaudeCodeManagedSettingRBACGroupSetParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationClaudeCodeManagedSettingRBACGroupSetParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationClaudeCodeManagedSettingRBACGroupSetParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
