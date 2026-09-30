package anthropic

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// OrganizationComplianceSettingService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationComplianceSettingService] method instead.
type OrganizationComplianceSettingService struct {
	Options []option.RequestOption
}

// NewOrganizationComplianceSettingService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewOrganizationComplianceSettingService(opts ...option.RequestOption) (r OrganizationComplianceSettingService) {
	r = OrganizationComplianceSettingService{}
	r.Options = opts
	return
}

// Retrieve your organization's Compliance Settings.
//
// Compliance Settings is a singleton resource: there is exactly one per
// organization, addressed without an identifier. The `state` field reflects
// whether the Compliance API is enabled. An organization with a parent
// organization reads the state inherited from the parent's configuration.
func (r *OrganizationComplianceSettingService) Get(ctx context.Context, opts ...option.RequestOption) (res *OrganizationComplianceSettings, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/compliance_settings"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update your organization's Compliance Settings.
//
// Setting `state` to `enabled` turns on the Compliance API and begins capturing
// organization activity events. Setting it to `disabled` turns both off. `state`
// reflects whether the Compliance API is enabled.
//
// A request that sets `state` to its current value succeeds and leaves the
// resource unchanged. A `disabled` request stays in effect until a later `enabled`
// request or the organization's next provisioning action that enables Access
// Transparency: enabling Access Transparency also enables the Compliance API,
// which serves its activity events, so such provisioning (including re-runs)
// re-enables the Compliance API even after a `disabled` request. Automated
// provisioning never disables compliance settings.
func (r *OrganizationComplianceSettingService) Update(ctx context.Context, body OrganizationComplianceSettingUpdateParams, opts ...option.RequestOption) (res *OrganizationComplianceSettings, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/compliance_settings"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// ComplianceSettingsStateUnion contains all possible properties and values from
// [ComplianceSettingsStateEnabled], [ComplianceSettingsStateDisabled].
//
// Use the [ComplianceSettingsStateUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ComplianceSettingsStateUnion struct {
	// Any of "enabled", "disabled".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyComplianceSettingsState is implemented by each variant of
// [ComplianceSettingsStateUnion] to add type safety for the return type of
// [ComplianceSettingsStateUnion.AsAny]
type anyComplianceSettingsState interface {
	implComplianceSettingsStateUnion()
}

func (ComplianceSettingsStateEnabled) implComplianceSettingsStateUnion()  {}
func (ComplianceSettingsStateDisabled) implComplianceSettingsStateUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ComplianceSettingsStateUnion.AsAny().(type) {
//	case anthropic.ComplianceSettingsStateEnabled:
//	case anthropic.ComplianceSettingsStateDisabled:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ComplianceSettingsStateUnion) AsAny() anyComplianceSettingsState {
	switch u.Type {
	case "enabled":
		return u.AsEnabled()
	case "disabled":
		return u.AsDisabled()
	}
	return nil
}

func (u ComplianceSettingsStateUnion) AsEnabled() (v ComplianceSettingsStateEnabled) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ComplianceSettingsStateUnion) AsDisabled() (v ComplianceSettingsStateDisabled) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ComplianceSettingsStateUnion) RawJSON() string { return u.JSON.raw }

func (r *ComplianceSettingsStateUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ComplianceSettingsStateDisabled struct {
	Type constant.Disabled `json:"type" default:"disabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ComplianceSettingsStateDisabled) RawJSON() string { return r.JSON.raw }
func (r *ComplianceSettingsStateDisabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewComplianceSettingsStateDisabledParam() ComplianceSettingsStateDisabledParam {
	return ComplianceSettingsStateDisabledParam{
		Type: "disabled",
	}
}

// This struct has a constant value, construct it with
// [NewComplianceSettingsStateDisabledParam].
type ComplianceSettingsStateDisabledParam struct {
	Type constant.Disabled `json:"type" default:"disabled"`
	paramObj
}

func (r ComplianceSettingsStateDisabledParam) MarshalJSON() (data []byte, err error) {
	type shadow ComplianceSettingsStateDisabledParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ComplianceSettingsStateDisabledParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ComplianceSettingsStateEnabled struct {
	Type constant.Enabled `json:"type" default:"enabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ComplianceSettingsStateEnabled) RawJSON() string { return r.JSON.raw }
func (r *ComplianceSettingsStateEnabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewComplianceSettingsStateEnabledParam() ComplianceSettingsStateEnabledParam {
	return ComplianceSettingsStateEnabledParam{
		Type: "enabled",
	}
}

// This struct has a constant value, construct it with
// [NewComplianceSettingsStateEnabledParam].
type ComplianceSettingsStateEnabledParam struct {
	Type constant.Enabled `json:"type" default:"enabled"`
	paramObj
}

func (r ComplianceSettingsStateEnabledParam) MarshalJSON() (data []byte, err error) {
	type shadow ComplianceSettingsStateEnabledParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ComplianceSettingsStateEnabledParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ComplianceSettingsStateParamUnion struct {
	OfEnabled  *ComplianceSettingsStateEnabledParam  `json:",omitzero,inline"`
	OfDisabled *ComplianceSettingsStateDisabledParam `json:",omitzero,inline"`
	paramUnion
}

func (u ComplianceSettingsStateParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfEnabled, u.OfDisabled)
}
func (u *ComplianceSettingsStateParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ComplianceSettingsStateParamUnion) asAny() any {
	if !param.IsOmitted(u.OfEnabled) {
		return u.OfEnabled
	} else if !param.IsOmitted(u.OfDisabled) {
		return u.OfDisabled
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ComplianceSettingsStateParamUnion) GetType() *string {
	if vt := u.OfEnabled; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDisabled; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ComplianceSettingsStateParamUnion](
		"type",
		apijson.Discriminator[ComplianceSettingsStateEnabledParam]("enabled"),
		apijson.Discriminator[ComplianceSettingsStateDisabledParam]("disabled"),
	)
}

type OrganizationComplianceSettings struct {
	// Whether the Compliance API is enabled for this organization.
	State ComplianceSettingsStateUnion `json:"state" api:"required"`
	Type  constant.ComplianceSettings  `json:"type" default:"compliance_settings"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		State       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationComplianceSettings) RawJSON() string { return r.JSON.raw }
func (r *OrganizationComplianceSettings) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationComplianceSettingUpdateParams struct {
	// Desired state. Accepts the string shorthand "enabled" or "disabled" in place of
	// the object form; the response always returns the canonical object form.
	State ComplianceSettingsStateParamUnion `json:"state,omitzero" api:"required"`
	paramObj
}

func (r OrganizationComplianceSettingUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationComplianceSettingUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationComplianceSettingUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
