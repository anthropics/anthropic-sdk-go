// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

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

// BetaOrganizationComplianceSettingService contains methods and other services
// that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationComplianceSettingService] method instead.
type BetaOrganizationComplianceSettingService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationComplianceSettingService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationComplianceSettingService(opts ...option.RequestOption) (r BetaOrganizationComplianceSettingService) {
	r = BetaOrganizationComplianceSettingService{}
	r.Options = opts
	return
}

// Retrieve your organization's Compliance Settings.
//
// Compliance Settings is a singleton resource: there is exactly one per
// organization, addressed without an identifier. The `state` field reflects
// whether the Compliance API is enabled. An organization with a parent
// organization reads the state inherited from the parent's configuration.
func (r *BetaOrganizationComplianceSettingService) Get(ctx context.Context, opts ...option.RequestOption) (res *BetaComplianceSettings, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/compliance_settings?beta=true"
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
func (r *BetaOrganizationComplianceSettingService) Update(ctx context.Context, body BetaOrganizationComplianceSettingUpdateParams, opts ...option.RequestOption) (res *BetaComplianceSettings, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/compliance_settings?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type BetaComplianceSettings struct {
	// Whether the Compliance API is enabled for this organization.
	State BetaComplianceSettingsStateUnion `json:"state" api:"required"`
	Type  constant.ComplianceSettings      `json:"type" default:"compliance_settings"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		State       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceSettings) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceSettings) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaComplianceSettingsStateUnion contains all possible properties and values
// from [BetaComplianceSettingsStateEnabled],
// [BetaComplianceSettingsStateDisabled].
//
// Use the [BetaComplianceSettingsStateUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaComplianceSettingsStateUnion struct {
	// Any of "enabled", "disabled".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaComplianceSettingsState is implemented by each variant of
// [BetaComplianceSettingsStateUnion] to add type safety for the return type of
// [BetaComplianceSettingsStateUnion.AsAny]
type anyBetaComplianceSettingsState interface {
	implBetaComplianceSettingsStateUnion()
}

func (BetaComplianceSettingsStateEnabled) implBetaComplianceSettingsStateUnion()  {}
func (BetaComplianceSettingsStateDisabled) implBetaComplianceSettingsStateUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaComplianceSettingsStateUnion.AsAny().(type) {
//	case anthropic.BetaComplianceSettingsStateEnabled:
//	case anthropic.BetaComplianceSettingsStateDisabled:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaComplianceSettingsStateUnion) AsAny() anyBetaComplianceSettingsState {
	switch u.Type {
	case "enabled":
		return u.AsEnabled()
	case "disabled":
		return u.AsDisabled()
	}
	return nil
}

func (u BetaComplianceSettingsStateUnion) AsEnabled() (v BetaComplianceSettingsStateEnabled) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaComplianceSettingsStateUnion) AsDisabled() (v BetaComplianceSettingsStateDisabled) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaComplianceSettingsStateUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaComplianceSettingsStateUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaComplianceSettingsStateDisabled struct {
	Type constant.Disabled `json:"type" default:"disabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceSettingsStateDisabled) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceSettingsStateDisabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewBetaComplianceSettingsStateDisabledParam() BetaComplianceSettingsStateDisabledParam {
	return BetaComplianceSettingsStateDisabledParam{
		Type: "disabled",
	}
}

// This struct has a constant value, construct it with
// [NewBetaComplianceSettingsStateDisabledParam].
type BetaComplianceSettingsStateDisabledParam struct {
	Type constant.Disabled `json:"type" default:"disabled"`
	paramObj
}

func (r BetaComplianceSettingsStateDisabledParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaComplianceSettingsStateDisabledParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaComplianceSettingsStateDisabledParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaComplianceSettingsStateEnabled struct {
	Type constant.Enabled `json:"type" default:"enabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceSettingsStateEnabled) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceSettingsStateEnabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewBetaComplianceSettingsStateEnabledParam() BetaComplianceSettingsStateEnabledParam {
	return BetaComplianceSettingsStateEnabledParam{
		Type: "enabled",
	}
}

// This struct has a constant value, construct it with
// [NewBetaComplianceSettingsStateEnabledParam].
type BetaComplianceSettingsStateEnabledParam struct {
	Type constant.Enabled `json:"type" default:"enabled"`
	paramObj
}

func (r BetaComplianceSettingsStateEnabledParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaComplianceSettingsStateEnabledParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaComplianceSettingsStateEnabledParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationComplianceSettingUpdateParams struct {
	// Desired state. Accepts the string shorthand "enabled" or "disabled" in place of
	// the object form; the response always returns the canonical object form.
	State BetaOrganizationComplianceSettingUpdateParamsStateUnion `json:"state,omitzero" api:"required"`
	paramObj
}

func (r BetaOrganizationComplianceSettingUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationComplianceSettingUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationComplianceSettingUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaOrganizationComplianceSettingUpdateParamsStateUnion struct {
	OfEnabled  *BetaComplianceSettingsStateEnabledParam  `json:",omitzero,inline"`
	OfDisabled *BetaComplianceSettingsStateDisabledParam `json:",omitzero,inline"`
	paramUnion
}

func (u BetaOrganizationComplianceSettingUpdateParamsStateUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfEnabled, u.OfDisabled)
}
func (u *BetaOrganizationComplianceSettingUpdateParamsStateUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaOrganizationComplianceSettingUpdateParamsStateUnion) asAny() any {
	if !param.IsOmitted(u.OfEnabled) {
		return u.OfEnabled
	} else if !param.IsOmitted(u.OfDisabled) {
		return u.OfDisabled
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationComplianceSettingUpdateParamsStateUnion) GetType() *string {
	if vt := u.OfEnabled; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDisabled; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaOrganizationComplianceSettingUpdateParamsStateUnion](
		"type",
		apijson.Discriminator[BetaComplianceSettingsStateEnabledParam]("enabled"),
		apijson.Discriminator[BetaComplianceSettingsStateDisabledParam]("disabled"),
	)
}
