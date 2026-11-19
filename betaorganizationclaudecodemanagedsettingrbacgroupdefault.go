package anthropic

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// BetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultService contains methods
// and other services that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultService] method
// instead.
type BetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultService generates a
// new service that applies the given options to each request. These options are
// applied after the parent client's options (if there is one), and before any
// request-specific options.
func NewBetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultService(opts ...option.RequestOption) (r BetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultService) {
	r = BetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultService{}
	r.Options = opts
	return
}

// Retrieve the default Claude Code managed settings, which members who belong to
// no RBAC Group in `rbac_group_order` receive in `rbac_groups` mode.
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
func (r *BetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultService) Get(ctx context.Context, query BetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultGetParams, opts ...option.RequestOption) (res *BetaClaudeCodeRBACGroupDefaultManagedSettings, err error) {
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
	path := "v1/organizations/claude_code/managed_settings/rbac_group_default"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_code.managed_settings.rbac_group_default.retrieve", "/v1/organizations/claude_code/managed_settings/rbac_group_default"))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Set or replace the default Claude Code managed settings, or clear them with
// `null`.
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
//   - [ErrorCodeBetaRequestBodyTimeout] (HTTP 400)
//   - [ErrorCodeBetaSettingsDocumentLimitExceeded] (HTTP 400)
//   - [ErrorCodeBetaSettingsHistoryLimitExceeded] (HTTP 400)
//   - [ErrorCodeSandboxEgressRestricted] (HTTP 403)
func (r *BetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultService) Update(ctx context.Context, params BetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultUpdateParams, opts ...option.RequestOption) (res *BetaClaudeCodeRBACGroupDefaultManagedSettings, err error) {
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
	path := "v1/organizations/claude_code/managed_settings/rbac_group_default"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.claude_code.managed_settings.rbac_group_default.update", "/v1/organizations/claude_code/managed_settings/rbac_group_default"))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type BetaClaudeCodeRBACGroupDefaultManagedSettings struct {
	// The default settings document, or `null` if none are set or the stored document
	// is damaged and can no longer be read. To tell these apart, `updated_at` is
	// `null` only when none are set. Replacing a damaged document fixes this.
	Settings map[string]any `json:"settings" api:"required"`
	// Always `claude_code_rbac_group_default_managed_settings`.
	Type constant.ClaudeCodeRBACGroupDefaultManagedSettings `json:"type" default:"claude_code_rbac_group_default_managed_settings"`
	// When the default settings last changed, in RFC 3339 format, or `null` if none
	// are set.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Settings    respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaClaudeCodeRBACGroupDefaultManagedSettings) RawJSON() string { return r.JSON.raw }
func (r *BetaClaudeCodeRBACGroupDefaultManagedSettings) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultGetParams struct {
	// This endpoint is in beta: requests must send
	// `ce-claude-code-managed-settings-2026-10-06` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultUpdateParams struct {
	// The whole settings document, a JSON object of Claude Code managed settings. It
	// replaces any stored document. At most 1 MiB as the server stores it, which can
	// be larger than the request body.
	//
	// To clear the default settings, set `Settings` to
	// `param.NullMap[map[string]any]()`. If `Settings` is nil, the API returns a
	// `400`. An empty map stores an empty document instead of clearing the default
	// settings.
	Settings map[string]any `json:"settings,omitzero" api:"required"`
	// This endpoint is in beta: requests must send
	// `ce-claude-code-managed-settings-2026-10-06` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationClaudeCodeManagedSettingRBACGroupDefaultUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
