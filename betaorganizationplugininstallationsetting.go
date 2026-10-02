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

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// BetaOrganizationPluginInstallationSettingService contains methods and other
// services that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationPluginInstallationSettingService] method instead.
type BetaOrganizationPluginInstallationSettingService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationPluginInstallationSettingService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationPluginInstallationSettingService(opts ...option.RequestOption) (r BetaOrganizationPluginInstallationSettingService) {
	r = BetaOrganizationPluginInstallationSettingService{}
	r.Options = opts
	return
}

// List an organization-owned Plugin's installation settings, which say which
// members it is for, most recently created first.
//
// The list holds the Plugin's own organization-wide setting (absent while the
// Plugin inherits its marketplace's default) and each RBAC Group's own setting. A
// member-owned Plugin has shares instead, so this path returns 404 for one.
//
// **Accepted credentials:** an Admin API key with the `read:plugins` or
// `read:org_audit` scope, or a Compliance Access Key with the
// `read:compliance_org_data` scope.
//
// Every request must include the beta header
// `anthropic-beta: ce-plugins-2026-09-01`. A request without it returns `404`,
// exactly as if the endpoint did not exist. The Plugins API is in beta and is
// available to Claude Enterprise organizations only. It is not available to Claude
// Platform (Claude Console) organizations, or to organizations with HIPAA
// readiness enabled.
func (r *BetaOrganizationPluginInstallationSettingService) List(ctx context.Context, pluginID string, params BetaOrganizationPluginInstallationSettingListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaPluginInstallationSetting], err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01"), option.WithResponseInto(&raw)}, opts...)
	if pluginID == "" {
		err = errors.New("missing required plugin_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/plugins/%s/installation_settings?beta=true", pluginID)
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

// List an organization-owned Plugin's installation settings, which say which
// members it is for, most recently created first.
//
// The list holds the Plugin's own organization-wide setting (absent while the
// Plugin inherits its marketplace's default) and each RBAC Group's own setting. A
// member-owned Plugin has shares instead, so this path returns 404 for one.
//
// **Accepted credentials:** an Admin API key with the `read:plugins` or
// `read:org_audit` scope, or a Compliance Access Key with the
// `read:compliance_org_data` scope.
//
// Every request must include the beta header
// `anthropic-beta: ce-plugins-2026-09-01`. A request without it returns `404`,
// exactly as if the endpoint did not exist. The Plugins API is in beta and is
// available to Claude Enterprise organizations only. It is not available to Claude
// Platform (Claude Console) organizations, or to organizations with HIPAA
// readiness enabled.
func (r *BetaOrganizationPluginInstallationSettingService) ListAutoPaging(ctx context.Context, pluginID string, params BetaOrganizationPluginInstallationSettingListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaPluginInstallationSetting] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, pluginID, params, opts...))
}

// Remove an organization-owned Plugin's own installation setting for the whole
// organization or for one RBAC Group.
//
// Removing the `organization` target returns the Plugin to its marketplace's
// default installation setting and leaves the groups' settings in place. Removing
// a group's setting makes the group's members fall back to the Plugin's
// organization-wide setting or to the settings of their other groups.
//
// A target that holds no setting of its own returns 404 (a Plugin that already
// inherits its marketplace's default holds no `organization` setting), and so does
// a member-owned Plugin.
//
// A removal counts as one of the Plugin's installation-setting writes: send all of
// those writes one at a time. If several arrive for the same Plugin at the same
// time, the server handles them one after another and can answer some of them with
// `503` and `x-should-retry: true` instead of applying them; wait a second or two
// and send the removal again. A `404` on the repeat means the setting is already
// gone.
//
// **Accepted credentials:** an Admin API key with the `write:plugins` scope.
//
// Every request must include the beta header
// `anthropic-beta: ce-plugins-2026-09-01`. A request without it returns `404`,
// exactly as if the endpoint did not exist. The Plugins API is in beta and is
// available to Claude Enterprise organizations only. It is not available to Claude
// Platform (Claude Console) organizations, or to organizations with HIPAA
// readiness enabled.
func (r *BetaOrganizationPluginInstallationSettingService) Remove(ctx context.Context, target string, params BetaOrganizationPluginInstallationSettingRemoveParams, opts ...option.RequestOption) (res *BetaDeletedPluginInstallationSetting, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01")}, opts...)
	if params.PluginID == "" {
		err = errors.New("missing required plugin_id parameter")
		return nil, err
	}
	if target == "" {
		err = errors.New("missing required target parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/plugins/%s/installation_settings/%s?beta=true", params.PluginID, target)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Set or change an organization-owned Plugin's installation setting for the whole
// organization or for one RBAC Group.
//
// Writing the value a target already holds of its own changes nothing.
//
// A member-owned Plugin has shares instead of installation settings, so this path
// returns 404 for one.
//
// Send a Plugin's installation-setting writes one at a time. If several writes for
// the same Plugin arrive at the same time, the server handles them one after
// another and can answer some of them with `503` instead of applying them. That
// `503` carries `x-should-retry: true`, and the write is safe to repeat: wait a
// second or two, then send it again.
//
// **Accepted credentials:** an Admin API key with the `write:plugins` scope.
//
// Every request must include the beta header
// `anthropic-beta: ce-plugins-2026-09-01`. A request without it returns `404`,
// exactly as if the endpoint did not exist. The Plugins API is in beta and is
// available to Claude Enterprise organizations only. It is not available to Claude
// Platform (Claude Console) organizations, or to organizations with HIPAA
// readiness enabled.
func (r *BetaOrganizationPluginInstallationSettingService) Set(ctx context.Context, target string, params BetaOrganizationPluginInstallationSettingSetParams, opts ...option.RequestOption) (res *BetaPluginInstallationSetting, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "ce-plugins-2026-09-01")}, opts...)
	if params.PluginID == "" {
		err = errors.New("missing required plugin_id parameter")
		return nil, err
	}
	if target == "" {
		err = errors.New("missing required target parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/plugins/%s/installation_settings/%s?beta=true", params.PluginID, target)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Confirmation that one target's installation setting was removed, naming the
// Plugin and the target in place of an ID.
type BetaDeletedPluginInstallationSetting struct {
	// The Plugin's ID.
	PluginID string `json:"plugin_id" api:"required"`
	// Whose setting was removed.
	Target BetaDeletedPluginInstallationSettingTargetUnion `json:"target" api:"required"`
	// Always `plugin_installation_setting_deleted`.
	Type constant.PluginInstallationSettingDeleted `json:"type" default:"plugin_installation_setting_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		PluginID    respjson.Field
		Target      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDeletedPluginInstallationSetting) RawJSON() string { return r.JSON.raw }
func (r *BetaDeletedPluginInstallationSetting) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaDeletedPluginInstallationSettingTargetUnion contains all possible properties
// and values from [BetaPluginTargetOrganization], [BetaPluginTargetRBACGroup],
// [BetaPluginTargetOrganizationMember].
//
// Use the [BetaDeletedPluginInstallationSettingTargetUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaDeletedPluginInstallationSettingTargetUnion struct {
	// Any of "organization", "rbac_group", "organization_member".
	Type string `json:"type"`
	// This field is from variant [BetaPluginTargetRBACGroup].
	RBACGroupID string `json:"rbac_group_id"`
	// This field is from variant [BetaPluginTargetOrganizationMember].
	UserID string `json:"user_id"`
	JSON   struct {
		Type        respjson.Field
		RBACGroupID respjson.Field
		UserID      respjson.Field
		raw         string
	} `json:"-"`
}

// anyBetaDeletedPluginInstallationSettingTarget is implemented by each variant of
// [BetaDeletedPluginInstallationSettingTargetUnion] to add type safety for the
// return type of [BetaDeletedPluginInstallationSettingTargetUnion.AsAny]
type anyBetaDeletedPluginInstallationSettingTarget interface {
	implBetaDeletedPluginInstallationSettingTargetUnion()
}

func (BetaPluginTargetOrganization) implBetaDeletedPluginInstallationSettingTargetUnion()       {}
func (BetaPluginTargetRBACGroup) implBetaDeletedPluginInstallationSettingTargetUnion()          {}
func (BetaPluginTargetOrganizationMember) implBetaDeletedPluginInstallationSettingTargetUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaDeletedPluginInstallationSettingTargetUnion.AsAny().(type) {
//	case anthropic.BetaPluginTargetOrganization:
//	case anthropic.BetaPluginTargetRBACGroup:
//	case anthropic.BetaPluginTargetOrganizationMember:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaDeletedPluginInstallationSettingTargetUnion) AsAny() anyBetaDeletedPluginInstallationSettingTarget {
	switch u.Type {
	case "organization":
		return u.AsOrganization()
	case "rbac_group":
		return u.AsRBACGroup()
	case "organization_member":
		return u.AsOrganizationMember()
	}
	return nil
}

func (u BetaDeletedPluginInstallationSettingTargetUnion) AsOrganization() (v BetaPluginTargetOrganization) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaDeletedPluginInstallationSettingTargetUnion) AsRBACGroup() (v BetaPluginTargetRBACGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaDeletedPluginInstallationSettingTargetUnion) AsOrganizationMember() (v BetaPluginTargetOrganizationMember) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaDeletedPluginInstallationSettingTargetUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaDeletedPluginInstallationSettingTargetUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The installation setting an organization-owned Plugin holds for one target. It
// has no ID of its own: it is addressed by the Plugin's ID and the target.
type BetaPluginInstallationSetting struct {
	// When the target was first given a setting for this Plugin.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// The setting the target holds for this Plugin. One of `required`, `auto_install`,
	// `available`, `not_available`; a value this API does not yet name is returned as
	// stored.
	//
	// Any of "auto_install", "available", "not_available", "required".
	InstallationPreference BetaPluginInstallationSettingInstallationPreference `json:"installation_preference" api:"required"`
	// The Plugin's ID.
	PluginID string `json:"plugin_id" api:"required"`
	// Whose setting this is: `organization` (the Plugin's own organization-wide
	// setting) or `rbac_group` (one RBAC Group's own setting); `organization_member`
	// does not occur here.
	Target BetaPluginInstallationSettingTargetUnion `json:"target" api:"required"`
	// Always `plugin_installation_setting`.
	Type constant.PluginInstallationSetting `json:"type" default:"plugin_installation_setting"`
	// When its setting last changed.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CreatedAt              respjson.Field
		InstallationPreference respjson.Field
		PluginID               respjson.Field
		Target                 respjson.Field
		Type                   respjson.Field
		UpdatedAt              respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPluginInstallationSetting) RawJSON() string { return r.JSON.raw }
func (r *BetaPluginInstallationSetting) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The setting the target holds for this Plugin. One of `required`, `auto_install`,
// `available`, `not_available`; a value this API does not yet name is returned as
// stored.
type BetaPluginInstallationSettingInstallationPreference string

const (
	BetaPluginInstallationSettingInstallationPreferenceAutoInstall  BetaPluginInstallationSettingInstallationPreference = "auto_install"
	BetaPluginInstallationSettingInstallationPreferenceAvailable    BetaPluginInstallationSettingInstallationPreference = "available"
	BetaPluginInstallationSettingInstallationPreferenceNotAvailable BetaPluginInstallationSettingInstallationPreference = "not_available"
	BetaPluginInstallationSettingInstallationPreferenceRequired     BetaPluginInstallationSettingInstallationPreference = "required"
)

// BetaPluginInstallationSettingTargetUnion contains all possible properties and
// values from [BetaPluginTargetOrganization], [BetaPluginTargetRBACGroup],
// [BetaPluginTargetOrganizationMember].
//
// Use the [BetaPluginInstallationSettingTargetUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaPluginInstallationSettingTargetUnion struct {
	// Any of "organization", "rbac_group", "organization_member".
	Type string `json:"type"`
	// This field is from variant [BetaPluginTargetRBACGroup].
	RBACGroupID string `json:"rbac_group_id"`
	// This field is from variant [BetaPluginTargetOrganizationMember].
	UserID string `json:"user_id"`
	JSON   struct {
		Type        respjson.Field
		RBACGroupID respjson.Field
		UserID      respjson.Field
		raw         string
	} `json:"-"`
}

// anyBetaPluginInstallationSettingTarget is implemented by each variant of
// [BetaPluginInstallationSettingTargetUnion] to add type safety for the return
// type of [BetaPluginInstallationSettingTargetUnion.AsAny]
type anyBetaPluginInstallationSettingTarget interface {
	implBetaPluginInstallationSettingTargetUnion()
}

func (BetaPluginTargetOrganization) implBetaPluginInstallationSettingTargetUnion()       {}
func (BetaPluginTargetRBACGroup) implBetaPluginInstallationSettingTargetUnion()          {}
func (BetaPluginTargetOrganizationMember) implBetaPluginInstallationSettingTargetUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaPluginInstallationSettingTargetUnion.AsAny().(type) {
//	case anthropic.BetaPluginTargetOrganization:
//	case anthropic.BetaPluginTargetRBACGroup:
//	case anthropic.BetaPluginTargetOrganizationMember:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaPluginInstallationSettingTargetUnion) AsAny() anyBetaPluginInstallationSettingTarget {
	switch u.Type {
	case "organization":
		return u.AsOrganization()
	case "rbac_group":
		return u.AsRBACGroup()
	case "organization_member":
		return u.AsOrganizationMember()
	}
	return nil
}

func (u BetaPluginInstallationSettingTargetUnion) AsOrganization() (v BetaPluginTargetOrganization) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaPluginInstallationSettingTargetUnion) AsRBACGroup() (v BetaPluginTargetRBACGroup) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaPluginInstallationSettingTargetUnion) AsOrganizationMember() (v BetaPluginTargetOrganizationMember) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaPluginInstallationSettingTargetUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaPluginInstallationSettingTargetUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationPluginInstallationSettingListParams struct {
	// For a `read:org_audit` or `read:compliance_org_data` key created for all of a
	// parent organization's linked organizations: a child organization of that parent
	// to read instead of the organization the key was created in, given as the
	// organization's UUID or its `org_`-prefixed ID. A value that is neither returns a
	// 400; an organization that is not a child of the key's parent, or where the
	// Plugins API is not available, returns a 404. Any other key may pass only its own
	// organization's ID here; another organization returns a 404.
	OrganizationID param.Opt[string] `query:"organization_id,omitzero" json:"-"`
	// Optionally set to the `next_page` token from the previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `100`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Only settings for this kind of target: `organization` (the organization-wide
	// setting) or `rbac_group` (an RBAC Group's).
	//
	// Any of "organization", "rbac_group".
	TargetType BetaOrganizationPluginInstallationSettingListParamsTargetType `query:"target_type,omitzero" json:"-"`
	// This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationPluginInstallationSettingListParams]'s
// query parameters as `url.Values`.
func (r BetaOrganizationPluginInstallationSettingListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Only settings for this kind of target: `organization` (the organization-wide
// setting) or `rbac_group` (an RBAC Group's).
type BetaOrganizationPluginInstallationSettingListParamsTargetType string

const (
	BetaOrganizationPluginInstallationSettingListParamsTargetTypeOrganization BetaOrganizationPluginInstallationSettingListParamsTargetType = "organization"
	BetaOrganizationPluginInstallationSettingListParamsTargetTypeRBACGroup    BetaOrganizationPluginInstallationSettingListParamsTargetType = "rbac_group"
)

type BetaOrganizationPluginInstallationSettingRemoveParams struct {
	// ID of the Plugin (prefixed `plugin_`).
	PluginID string `path:"plugin_id" api:"required" json:"-"`
	// This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationPluginInstallationSettingSetParams struct {
	// ID of the Plugin (prefixed `plugin_`).
	PluginID string `path:"plugin_id" api:"required" json:"-"`
	// The installation setting the target is to hold for this Plugin: one of
	// `required`, `auto_install`, `available`, `not_available`.
	//
	// Any of "auto_install", "available", "not_available", "required".
	InstallationPreference BetaOrganizationPluginInstallationSettingSetParamsInstallationPreference `json:"installation_preference,omitzero" api:"required"`
	// This endpoint is in beta: requests must send `ce-plugins-2026-09-01` in this
	// header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationPluginInstallationSettingSetParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationPluginInstallationSettingSetParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationPluginInstallationSettingSetParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The installation setting the target is to hold for this Plugin: one of
// `required`, `auto_install`, `available`, `not_available`.
type BetaOrganizationPluginInstallationSettingSetParamsInstallationPreference string

const (
	BetaOrganizationPluginInstallationSettingSetParamsInstallationPreferenceAutoInstall  BetaOrganizationPluginInstallationSettingSetParamsInstallationPreference = "auto_install"
	BetaOrganizationPluginInstallationSettingSetParamsInstallationPreferenceAvailable    BetaOrganizationPluginInstallationSettingSetParamsInstallationPreference = "available"
	BetaOrganizationPluginInstallationSettingSetParamsInstallationPreferenceNotAvailable BetaOrganizationPluginInstallationSettingSetParamsInstallationPreference = "not_available"
	BetaOrganizationPluginInstallationSettingSetParamsInstallationPreferenceRequired     BetaOrganizationPluginInstallationSettingSetParamsInstallationPreference = "required"
)
