package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/paramutil"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// BetaOrganizationClaudeCodeManagedSettingConfigurationService contains methods
// and other services that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationClaudeCodeManagedSettingConfigurationService] method
// instead.
type BetaOrganizationClaudeCodeManagedSettingConfigurationService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationClaudeCodeManagedSettingConfigurationService generates a new
// service that applies the given options to each request. These options are
// applied after the parent client's options (if there is one), and before any
// request-specific options.
func NewBetaOrganizationClaudeCodeManagedSettingConfigurationService(opts ...option.RequestOption) (r BetaOrganizationClaudeCodeManagedSettingConfigurationService) {
	r = BetaOrganizationClaudeCodeManagedSettingConfigurationService{}
	r.Options = opts
	return
}

// Retrieve the organization's Claude Code managed settings mode and RBAC Group
// order.
//
// **Accepted credentials:** an organization-level Admin API key with the
// `read:claude_code_managed_settings` scope.
//
// Every request must include the beta header
// `anthropic-beta: ce-claude-code-managed-settings-2026-10-06`. A request without
// it, or from an organization not enabled for this beta, returns `404`. The Claude
// Code Managed Settings API is in beta and is available to Claude Enterprise
// organizations only.
func (r *BetaOrganizationClaudeCodeManagedSettingConfigurationService) Get(ctx context.Context, query BetaOrganizationClaudeCodeManagedSettingConfigurationGetParams, opts ...option.RequestOption) (res *BetaClaudeCodeManagedSettingsConfiguration, err error) {
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
	path := "v1/organizations/claude_code/managed_settings/configuration"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_code.managed_settings.configuration.retrieve", "/v1/organizations/claude_code/managed_settings/configuration"))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Change the organization's Claude Code managed settings mode, the RBAC Group
// order, or both, in one atomic write.
//
// To change the order, send `rbac_group_order` with the `rbac_group_order_version`
// from the read the change is based on. If the order has changed since that read,
// the request returns `409`: read it again and retry. To switch modes only if the
// order has not changed, send the current `rbac_group_order` and
// `rbac_group_order_version` with the new `mode`.
//
// A switch to `rbac_groups` mode with no RBAC Group settings and no default
// settings copies the organization settings into a new default, in the same write,
// so members keep the settings they had. The copy does not follow later changes to
// the organization settings. If the organization settings cannot be copied, the
// request returns `400` and changes nothing. If they cannot be read, are over 1
// MiB or are not valid, setting the default settings first avoids the copy.
//
// After a switch to `rbac_groups` mode, a member in no listed RBAC Group receives
// the default settings, or no settings if no default is set, not the organization
// settings.
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
//   - [ErrorCodeBetaRequestBodyTimeout] (HTTP 400)
//   - [ErrorCodeBetaSettingsDocumentLimitExceeded] (HTTP 400)
//   - [ErrorCodeBetaSettingsHistoryLimitExceeded] (HTTP 400)
//   - [ErrorCodeSandboxEgressRestricted] (HTTP 403)
func (r *BetaOrganizationClaudeCodeManagedSettingConfigurationService) Update(ctx context.Context, params BetaOrganizationClaudeCodeManagedSettingConfigurationUpdateParams, opts ...option.RequestOption) (res *BetaClaudeCodeManagedSettingsConfiguration, err error) {
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
	path := "v1/organizations/claude_code/managed_settings/configuration"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_code.managed_settings.configuration.update", "/v1/organizations/claude_code/managed_settings/configuration"))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type BetaClaudeCodeManagedSettingsConfiguration struct {
	// How the organization serves managed settings. In `organization` mode every
	// member receives the organization settings, which are set in claude.ai. In
	// `rbac_groups` mode each member receives the settings of the first RBAC Group in
	// `rbac_group_order` that they belong to. A member in none of them receives the
	// default settings, or no settings if no default is set.
	Mode BetaClaudeCodeManagedSettingsModeUnion `json:"mode" api:"required"`
	// When the mode last changed, in RFC 3339 format, or `null` if it never has.
	ModeUpdatedAt time.Time `json:"mode_updated_at" api:"required" format:"date-time"`
	// The ID of every RBAC Group with settings, deleted groups included, highest
	// precedence first.
	RBACGroupOrder []string `json:"rbac_group_order" api:"required"`
	// An opaque version of `rbac_group_order`, which changes whenever the order does.
	// Send it back with a new order.
	RBACGroupOrderVersion string `json:"rbac_group_order_version" api:"required"`
	// Always `claude_code_managed_settings_configuration`.
	Type constant.ClaudeCodeManagedSettingsConfiguration `json:"type" default:"claude_code_managed_settings_configuration"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Mode                  respjson.Field
		ModeUpdatedAt         respjson.Field
		RBACGroupOrder        respjson.Field
		RBACGroupOrderVersion respjson.Field
		Type                  respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeCodeManagedSettingsConfiguration) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeCodeManagedSettingsConfiguration) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaClaudeCodeManagedSettingsModeUnion contains all possible properties and
// values from [BetaClaudeCodeManagedSettingsOrganizationMode],
// [BetaClaudeCodeManagedSettingsRBACGroupsMode].
//
// Use the [BetaClaudeCodeManagedSettingsModeUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaClaudeCodeManagedSettingsModeUnion struct {
	// Any of "organization", "rbac_groups".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaClaudeCodeManagedSettingsMode is implemented by each variant of
// [BetaClaudeCodeManagedSettingsModeUnion] to add type safety for the return type
// of [BetaClaudeCodeManagedSettingsModeUnion.AsAny]
type anyBetaClaudeCodeManagedSettingsMode interface {
	implBetaClaudeCodeManagedSettingsModeUnion()
}

func (BetaClaudeCodeManagedSettingsOrganizationMode) implBetaClaudeCodeManagedSettingsModeUnion() {}
func (BetaClaudeCodeManagedSettingsRBACGroupsMode) implBetaClaudeCodeManagedSettingsModeUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaClaudeCodeManagedSettingsModeUnion.AsAny().(type) {
//	case anthropic.BetaClaudeCodeManagedSettingsOrganizationMode:
//	case anthropic.BetaClaudeCodeManagedSettingsRBACGroupsMode:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaClaudeCodeManagedSettingsModeUnion) AsAny() anyBetaClaudeCodeManagedSettingsMode {
	switch u.Type {
	case "organization":
		return u.AsOrganization()
	case "rbac_groups":
		return u.AsRBACGroups()
	}
	return nil
}

func (u BetaClaudeCodeManagedSettingsModeUnion) AsOrganization() (v BetaClaudeCodeManagedSettingsOrganizationMode) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaClaudeCodeManagedSettingsModeUnion) AsRBACGroups() (v BetaClaudeCodeManagedSettingsRBACGroupsMode) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaClaudeCodeManagedSettingsModeUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaClaudeCodeManagedSettingsModeUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaClaudeCodeManagedSettingsModeParamUnion struct {
	OfOrganization *BetaClaudeCodeManagedSettingsOrganizationModeParam `json:",omitzero,inline"`
	OfRBACGroups   *BetaClaudeCodeManagedSettingsRBACGroupsModeParam   `json:",omitzero,inline"`
	paramUnion
}

func (u BetaClaudeCodeManagedSettingsModeParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfOrganization, u.OfRBACGroups)
}
func (u *BetaClaudeCodeManagedSettingsModeParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaClaudeCodeManagedSettingsModeParamUnion) asAny() any {
	if !param.IsOmitted(u.OfOrganization) {
		return u.OfOrganization
	} else if !param.IsOmitted(u.OfRBACGroups) {
		return u.OfRBACGroups
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaClaudeCodeManagedSettingsModeParamUnion) GetType() *string {
	if vt := u.OfOrganization; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	} else if vt := u.OfRBACGroups; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaClaudeCodeManagedSettingsModeParamUnion](
		"type",
		apijson.Discriminator[BetaClaudeCodeManagedSettingsOrganizationModeParam]("organization"),
		apijson.Discriminator[BetaClaudeCodeManagedSettingsRBACGroupsModeParam]("rbac_groups"),
	)
}

type BetaClaudeCodeManagedSettingsOrganizationMode struct {
	// Always `organization`.
	Type constant.Organization `json:"type" default:"organization"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeCodeManagedSettingsOrganizationMode) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeCodeManagedSettingsOrganizationMode) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewBetaClaudeCodeManagedSettingsOrganizationModeParam() BetaClaudeCodeManagedSettingsOrganizationModeParam {
	return BetaClaudeCodeManagedSettingsOrganizationModeParam{
		Type: "organization",
	}
}

// This struct has a constant value, construct it with
// [NewBetaClaudeCodeManagedSettingsOrganizationModeParam].
type BetaClaudeCodeManagedSettingsOrganizationModeParam struct {
	Type constant.Organization `json:"type" default:"organization"`
	paramObj
}

func (r BetaClaudeCodeManagedSettingsOrganizationModeParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaClaudeCodeManagedSettingsOrganizationModeParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaClaudeCodeManagedSettingsOrganizationModeParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaClaudeCodeManagedSettingsRBACGroupsMode struct {
	// Always `rbac_groups`.
	Type constant.RBACGroups `json:"type" default:"rbac_groups"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeCodeManagedSettingsRBACGroupsMode) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeCodeManagedSettingsRBACGroupsMode) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewBetaClaudeCodeManagedSettingsRBACGroupsModeParam() BetaClaudeCodeManagedSettingsRBACGroupsModeParam {
	return BetaClaudeCodeManagedSettingsRBACGroupsModeParam{
		Type: "rbac_groups",
	}
}

// This struct has a constant value, construct it with
// [NewBetaClaudeCodeManagedSettingsRBACGroupsModeParam].
type BetaClaudeCodeManagedSettingsRBACGroupsModeParam struct {
	Type constant.RBACGroups `json:"type" default:"rbac_groups"`
	paramObj
}

func (r BetaClaudeCodeManagedSettingsRBACGroupsModeParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaClaudeCodeManagedSettingsRBACGroupsModeParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaClaudeCodeManagedSettingsRBACGroupsModeParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationClaudeCodeManagedSettingConfigurationGetParams struct {
	// This endpoint is in beta: requests must send
	// `ce-claude-code-managed-settings-2026-10-06` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationClaudeCodeManagedSettingConfigurationUpdateParams struct {
	// The `rbac_group_order_version` from the read this change is based on. Requires
	// `rbac_group_order`.
	RBACGroupOrderVersion param.Opt[string] `json:"rbac_group_order_version,omitzero"`
	// The new mode, as an object such as `{"type": "rbac_groups"}`.
	Mode BetaClaudeCodeManagedSettingsModeParamUnion `json:"mode,omitzero"`
	// Every RBAC Group that has settings, each listed once, in the new order, highest
	// precedence first. Requires `rbac_group_order_version`.
	RBACGroupOrder []string `json:"rbac_group_order,omitzero"`
	// This endpoint is in beta: requests must send
	// `ce-claude-code-managed-settings-2026-10-06` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationClaudeCodeManagedSettingConfigurationUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationClaudeCodeManagedSettingConfigurationUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationClaudeCodeManagedSettingConfigurationUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
