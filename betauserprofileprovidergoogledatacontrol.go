package anthropic

import (
	"context"
	"encoding/json"
	"errors"
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

// BetaUserProfileProviderGoogleDataControlService contains methods and other
// services that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaUserProfileProviderGoogleDataControlService] method instead.
type BetaUserProfileProviderGoogleDataControlService struct {
	Options []option.RequestOption
}

// NewBetaUserProfileProviderGoogleDataControlService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaUserProfileProviderGoogleDataControlService(opts ...option.RequestOption) (r BetaUserProfileProviderGoogleDataControlService) {
	r = BetaUserProfileProviderGoogleDataControlService{}
	r.Options = opts
	return
}

// Get the profile's encryption keys and storage buckets for Claude on Google
// Cloud's Agent Platform, with each region's status.
func (r *BetaUserProfileProviderGoogleDataControlService) Get(ctx context.Context, userProfileID string, query BetaUserProfileProviderGoogleDataControlGetParams, opts ...option.RequestOption) (res *BetaGoogleManagedDataControls, err error) {
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
	if !param.IsOmitted(query.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", query.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "user-profiles-2026-09-04,google-managed-data-controls-2026-09-01")}, opts...)
	if userProfileID == "" {
		err = errors.New("missing required user_profile_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/user_profiles/%s/providers/google/data_controls", requestconfig.EscapePathSegment(userProfileID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.user_profiles.providers.google.data_controls.retrieve", "/v1/user_profiles/{user_profile_id}/providers/google/data_controls"), requestconfig.WithSpanAttribute("anthropic.user_profile.id", userProfileID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Turn a control on or off, or change the key or bucket set for a region. Controls
// and regions that you leave out stay as they are.
func (r *BetaUserProfileProviderGoogleDataControlService) Update(ctx context.Context, userProfileID string, params BetaUserProfileProviderGoogleDataControlUpdateParams, opts ...option.RequestOption) (res *BetaGoogleManagedDataControls, err error) {
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
	if !param.IsOmitted(params.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", params.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "user-profiles-2026-09-04,google-managed-data-controls-2026-09-01")}, opts...)
	if userProfileID == "" {
		err = errors.New("missing required user_profile_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/user_profiles/%s/providers/google/data_controls", requestconfig.EscapePathSegment(userProfileID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.user_profiles.providers.google.data_controls.update", "/v1/user_profiles/{user_profile_id}/providers/google/data_controls"), requestconfig.WithSpanAttribute("anthropic.user_profile.id", userProfileID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// A Cloud Storage bucket, named by the project that holds it and by its own name.
type BetaGoogleManagedBucketReference struct {
	// The bucket's Cloud Storage name.
	BucketName string `json:"bucket_name" api:"required"`
	// The number of the Google Cloud project that holds the bucket, in decimal.
	ProjectNumber string `json:"project_number" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		BucketName    respjson.Field
		ProjectNumber respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedBucketReference) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedBucketReference) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaGoogleManagedBucketReference to a
// BetaGoogleManagedBucketReferenceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaGoogleManagedBucketReferenceParam.Overrides()
func (r BetaGoogleManagedBucketReference) ToParam() BetaGoogleManagedBucketReferenceParam {
	return param.Override[BetaGoogleManagedBucketReferenceParam](json.RawMessage(r.RawJSON()))
}

// A Cloud Storage bucket, named by the project that holds it and by its own name.
//
// The properties BucketName, ProjectNumber are required.
type BetaGoogleManagedBucketReferenceParam struct {
	// The bucket's Cloud Storage name.
	BucketName string `json:"bucket_name" api:"required"`
	// The number of the Google Cloud project that holds the bucket, in decimal.
	ProjectNumber string `json:"project_number" api:"required"`
	paramObj
}

func (r BetaGoogleManagedBucketReferenceParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaGoogleManagedBucketReferenceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaGoogleManagedBucketReferenceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaGoogleManagedCMEKControlUnion contains all possible properties and values
// from [BetaGoogleManagedCMEKDisabled], [BetaGoogleManagedCMEKEnabled].
//
// Use the [BetaGoogleManagedCMEKControlUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaGoogleManagedCMEKControlUnion struct {
	RetiredKeys []BetaGoogleManagedRetiredKey `json:"retired_keys"`
	// Any of "disabled", "enabled".
	Type string `json:"type"`
	// This field is from variant [BetaGoogleManagedCMEKEnabled].
	ExternalKeyIDs map[string]string `json:"external_key_ids"`
	// This field is from variant [BetaGoogleManagedCMEKEnabled].
	Status map[string]BetaGoogleManagedCMEKRegionStatusUnion `json:"status"`
	JSON   struct {
		RetiredKeys    respjson.Field
		Type           respjson.Field
		ExternalKeyIDs respjson.Field
		Status         respjson.Field
		raw            string
	} `json:"-"`
}

// anyBetaGoogleManagedCMEKControl is implemented by each variant of
// [BetaGoogleManagedCMEKControlUnion] to add type safety for the return type of
// [BetaGoogleManagedCMEKControlUnion.AsAny]
type anyBetaGoogleManagedCMEKControl interface {
	implBetaGoogleManagedCMEKControlUnion()
}

func (BetaGoogleManagedCMEKDisabled) implBetaGoogleManagedCMEKControlUnion() {}
func (BetaGoogleManagedCMEKEnabled) implBetaGoogleManagedCMEKControlUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaGoogleManagedCMEKControlUnion.AsAny().(type) {
//	case anthropic.BetaGoogleManagedCMEKDisabled:
//	case anthropic.BetaGoogleManagedCMEKEnabled:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaGoogleManagedCMEKControlUnion) AsAny() anyBetaGoogleManagedCMEKControl {
	switch u.Type {
	case "disabled":
		return u.AsDisabled()
	case "enabled":
		return u.AsEnabled()
	}
	return nil
}

func (u BetaGoogleManagedCMEKControlUnion) AsDisabled() (v BetaGoogleManagedCMEKDisabled) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedCMEKControlUnion) AsEnabled() (v BetaGoogleManagedCMEKEnabled) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaGoogleManagedCMEKControlUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaGoogleManagedCMEKControlUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaGoogleManagedCMEKControlParamsUnion struct {
	OfDisabled *BetaGoogleManagedCMEKDisabledParams `json:",omitzero,inline"`
	OfEnabled  *BetaGoogleManagedCMEKEnabledParams  `json:",omitzero,inline"`
	paramUnion
}

func (u BetaGoogleManagedCMEKControlParamsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDisabled, u.OfEnabled)
}
func (u *BetaGoogleManagedCMEKControlParamsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaGoogleManagedCMEKControlParamsUnion) asAny() any {
	if !param.IsOmitted(u.OfDisabled) {
		return u.OfDisabled
	} else if !param.IsOmitted(u.OfEnabled) {
		return u.OfEnabled
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaGoogleManagedCMEKControlParamsUnion) GetExternalKeyIDs() map[string]string {
	if vt := u.OfEnabled; vt != nil {
		return vt.ExternalKeyIDs
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaGoogleManagedCMEKControlParamsUnion) GetType() *string {
	if vt := u.OfDisabled; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	} else if vt := u.OfEnabled; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaGoogleManagedCMEKControlParamsUnion](
		"type",
		apijson.Discriminator[BetaGoogleManagedCMEKDisabledParams]("disabled"),
		apijson.Discriminator[BetaGoogleManagedCMEKEnabledParams]("enabled"),
	)
}

// The control is off and has no keys.
type BetaGoogleManagedCMEKDisabled struct {
	// Keys that a region has retired, earliest `expires_at` first.
	RetiredKeys []BetaGoogleManagedRetiredKey `json:"retired_keys" api:"required"`
	Type        constant.Disabled             `json:"type" default:"disabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		RetiredKeys respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedCMEKDisabled) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedCMEKDisabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewBetaGoogleManagedCMEKDisabledParams() BetaGoogleManagedCMEKDisabledParams {
	return BetaGoogleManagedCMEKDisabledParams{
		Type: "disabled",
	}
}

// Turns the control off and removes every region's key.
//
// This struct has a constant value, construct it with
// [NewBetaGoogleManagedCMEKDisabledParams].
type BetaGoogleManagedCMEKDisabledParams struct {
	Type constant.Disabled `json:"type" default:"disabled"`
	paramObj
}

func (r BetaGoogleManagedCMEKDisabledParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaGoogleManagedCMEKDisabledParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaGoogleManagedCMEKDisabledParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The control is on and has a key for at least one region.
type BetaGoogleManagedCMEKEnabled struct {
	// The ID of the External Key configured for each region (`ekey_...`), keyed by
	// region code. A region with no key is left out.
	ExternalKeyIDs map[string]string `json:"external_key_ids" api:"required"`
	// Keys that a region has retired, earliest `expires_at` first.
	RetiredKeys []BetaGoogleManagedRetiredKey `json:"retired_keys" api:"required"`
	// Each region's status for the control, keyed by region code. A region that your
	// organization is not offered may have no entry, and Anthropic may add regions.
	Status map[string]BetaGoogleManagedCMEKRegionStatusUnion `json:"status" api:"required"`
	Type   constant.Enabled                                  `json:"type" default:"enabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ExternalKeyIDs respjson.Field
		RetiredKeys    respjson.Field
		Status         respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedCMEKEnabled) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedCMEKEnabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Turns the control on, or changes its keys while it is on.
//
// The property Type is required.
type BetaGoogleManagedCMEKEnabledParams struct {
	// The External Key to use in each region, keyed by region code, such as `us` or
	// `eu`. An External Key ID (`ekey_...`) sets that region's key, `null` removes it,
	// and a region you leave out keeps its key. While the control is on, at least one
	// region must have a key.
	ExternalKeyIDs map[string]string `json:"external_key_ids,omitzero"`
	// This field can be elided, and will marshal its zero value as "enabled".
	Type constant.Enabled `json:"type" default:"enabled"`
	paramObj
}

func (r BetaGoogleManagedCMEKEnabledParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaGoogleManagedCMEKEnabledParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaGoogleManagedCMEKEnabledParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The configured key passed validation, and the region applies it.
type BetaGoogleManagedCMEKRegionActive struct {
	// The ID of the External Key that the region applies now (`ekey_...`). While a
	// replacement key is `pending` or `failed`, this is the previous key. `null` when
	// the region applies no key.
	EffectiveExternalKeyID string `json:"effective_external_key_id" api:"required"`
	// Why Google Cloud KMS refused the region's latest use of the key. `null` when the
	// region has no refusal to report. A `null` does not mean that the key works.
	Error string `json:"error" api:"required"`
	// When the region last got a successful answer from Google Cloud KMS for the key,
	// in RFC 3339 format. `null` when it has none to report.
	LastCheckedAt time.Time       `json:"last_checked_at" api:"required" format:"date-time"`
	Type          constant.Active `json:"type" default:"active"`
	// When the status last became `active`, in RFC 3339 format.
	ValidatedAt time.Time `json:"validated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EffectiveExternalKeyID respjson.Field
		Error                  respjson.Field
		LastCheckedAt          respjson.Field
		Type                   respjson.Field
		ValidatedAt            respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedCMEKRegionActive) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedCMEKRegionActive) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The region stopped retrying validation before the configured key passed. To
// restart validation, send the region a key again. It can be the same key.
type BetaGoogleManagedCMEKRegionFailed struct {
	// The ID of the External Key that the region applies now (`ekey_...`). While a
	// replacement key is `pending` or `failed`, this is the previous key. `null` when
	// the region applies no key.
	EffectiveExternalKeyID string `json:"effective_external_key_id" api:"required"`
	// Why the region cannot use the configured key.
	Error string          `json:"error" api:"required"`
	Type  constant.Failed `json:"type" default:"failed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EffectiveExternalKeyID respjson.Field
		Error                  respjson.Field
		Type                   respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedCMEKRegionFailed) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedCMEKRegionFailed) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// No key is set for this region.
type BetaGoogleManagedCMEKRegionNotConfigured struct {
	// The ID of the External Key that the region applies now (`ekey_...`). While a
	// replacement key is `pending` or `failed`, this is the previous key. `null` when
	// the region applies no key.
	EffectiveExternalKeyID string                 `json:"effective_external_key_id" api:"required"`
	Type                   constant.NotConfigured `json:"type" default:"not_configured"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EffectiveExternalKeyID respjson.Field
		Type                   respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedCMEKRegionNotConfigured) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedCMEKRegionNotConfigured) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The region is validating the configured key, or the key has passed validation
// and the region has not started to apply it.
type BetaGoogleManagedCMEKRegionPending struct {
	// The ID of the External Key that the region applies now (`ekey_...`). While a
	// replacement key is `pending` or `failed`, this is the previous key. `null` when
	// the region applies no key.
	EffectiveExternalKeyID string `json:"effective_external_key_id" api:"required"`
	// Why the latest validation attempt failed. `null` when no attempt has failed.
	Error string           `json:"error" api:"required"`
	Type  constant.Pending `json:"type" default:"pending"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EffectiveExternalKeyID respjson.Field
		Error                  respjson.Field
		Type                   respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedCMEKRegionPending) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedCMEKRegionPending) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaGoogleManagedCMEKRegionStatusUnion contains all possible properties and
// values from [BetaGoogleManagedCMEKRegionNotConfigured],
// [BetaGoogleManagedCMEKRegionPending], [BetaGoogleManagedCMEKRegionActive],
// [BetaGoogleManagedCMEKRegionFailed], [BetaGoogleManagedCMEKRegionSuspended].
//
// Use the [BetaGoogleManagedCMEKRegionStatusUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaGoogleManagedCMEKRegionStatusUnion struct {
	EffectiveExternalKeyID string `json:"effective_external_key_id"`
	// Any of "not_configured", "pending", "active", "failed", "suspended".
	Type  string `json:"type"`
	Error string `json:"error"`
	// This field is from variant [BetaGoogleManagedCMEKRegionActive].
	LastCheckedAt time.Time `json:"last_checked_at"`
	// This field is from variant [BetaGoogleManagedCMEKRegionActive].
	ValidatedAt time.Time `json:"validated_at"`
	JSON        struct {
		EffectiveExternalKeyID respjson.Field
		Type                   respjson.Field
		Error                  respjson.Field
		LastCheckedAt          respjson.Field
		ValidatedAt            respjson.Field
		raw                    string
	} `json:"-"`
}

// anyBetaGoogleManagedCMEKRegionStatus is implemented by each variant of
// [BetaGoogleManagedCMEKRegionStatusUnion] to add type safety for the return type
// of [BetaGoogleManagedCMEKRegionStatusUnion.AsAny]
type anyBetaGoogleManagedCMEKRegionStatus interface {
	implBetaGoogleManagedCMEKRegionStatusUnion()
}

func (BetaGoogleManagedCMEKRegionNotConfigured) implBetaGoogleManagedCMEKRegionStatusUnion() {}
func (BetaGoogleManagedCMEKRegionPending) implBetaGoogleManagedCMEKRegionStatusUnion()       {}
func (BetaGoogleManagedCMEKRegionActive) implBetaGoogleManagedCMEKRegionStatusUnion()        {}
func (BetaGoogleManagedCMEKRegionFailed) implBetaGoogleManagedCMEKRegionStatusUnion()        {}
func (BetaGoogleManagedCMEKRegionSuspended) implBetaGoogleManagedCMEKRegionStatusUnion()     {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaGoogleManagedCMEKRegionStatusUnion.AsAny().(type) {
//	case anthropic.BetaGoogleManagedCMEKRegionNotConfigured:
//	case anthropic.BetaGoogleManagedCMEKRegionPending:
//	case anthropic.BetaGoogleManagedCMEKRegionActive:
//	case anthropic.BetaGoogleManagedCMEKRegionFailed:
//	case anthropic.BetaGoogleManagedCMEKRegionSuspended:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaGoogleManagedCMEKRegionStatusUnion) AsAny() anyBetaGoogleManagedCMEKRegionStatus {
	switch u.Type {
	case "not_configured":
		return u.AsNotConfigured()
	case "pending":
		return u.AsPending()
	case "active":
		return u.AsActive()
	case "failed":
		return u.AsFailed()
	case "suspended":
		return u.AsSuspended()
	}
	return nil
}

func (u BetaGoogleManagedCMEKRegionStatusUnion) AsNotConfigured() (v BetaGoogleManagedCMEKRegionNotConfigured) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedCMEKRegionStatusUnion) AsPending() (v BetaGoogleManagedCMEKRegionPending) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedCMEKRegionStatusUnion) AsActive() (v BetaGoogleManagedCMEKRegionActive) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedCMEKRegionStatusUnion) AsFailed() (v BetaGoogleManagedCMEKRegionFailed) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedCMEKRegionStatusUnion) AsSuspended() (v BetaGoogleManagedCMEKRegionSuspended) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaGoogleManagedCMEKRegionStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaGoogleManagedCMEKRegionStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Anthropic has withdrawn or paused its approval of customer-managed keys for this
// profile, so the region does not apply the control.
type BetaGoogleManagedCMEKRegionSuspended struct {
	// The ID of the External Key that the region applies now (`ekey_...`). While a
	// replacement key is `pending` or `failed`, this is the previous key. `null` when
	// the region applies no key.
	EffectiveExternalKeyID string             `json:"effective_external_key_id" api:"required"`
	Type                   constant.Suspended `json:"type" default:"suspended"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EffectiveExternalKeyID respjson.Field
		Type                   respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedCMEKRegionSuspended) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedCMEKRegionSuspended) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A user profile's encryption key and storage bucket settings for Claude on Google
// Cloud's Agent Platform. Every user profile has one of these objects. You can
// read and update it, but not create or delete it. Keys and buckets are set for
// each region. A request's region depends on the Google endpoint it is sent to:
// `us` is the US multi-region endpoint and the regional endpoint of every location
// whose name starts with `us-`, `eu` is the EU multi-region endpoint and the
// regional endpoint of every location whose name starts with `europe-`, `sg` is
// the `asia-southeast1` regional endpoint, and `global` is the global endpoint,
// wherever Google serves the request. Anthropic may add regions.
type BetaGoogleManagedDataControls struct {
	// The label value that every bucket set up for this profile must have. It never
	// changes.
	BucketLabelValue string `json:"bucket_label_value" api:"required"`
	// The customer-managed encryption key control: the key in the end customer's
	// Google Cloud KMS that is set for each region.
	CMEK BetaGoogleManagedCMEKControlUnion `json:"cmek" api:"required"`
	// The customer-owned storage control: the end customer's Cloud Storage bucket that
	// is set for each region.
	Storage BetaGoogleManagedStorageControlUnion `json:"storage" api:"required"`
	// Object type. Always `google_managed_data_controls`.
	Type constant.GoogleManagedDataControls `json:"type" default:"google_managed_data_controls"`
	// When an update last changed a control, or a region's status last changed, in RFC
	// 3339 format. A region that starts or stops reading `suspended`, a change to
	// `last_checked_at` alone, and the removal of a retired entry do not count.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// The ID of the user profile that these controls belong to (`uprof_...`).
	UserProfileID string `json:"user_profile_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		BucketLabelValue respjson.Field
		CMEK             respjson.Field
		Storage          respjson.Field
		Type             respjson.Field
		UpdatedAt        respjson.Field
		UserProfileID    respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedDataControls) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedDataControls) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A bucket that a region has retired. Until `expires_at`, keep the bucket and
// Anthropic's access to it as they were when the bucket passed verification.
type BetaGoogleManagedRetiredBucket struct {
	// The bucket's Cloud Storage name.
	BucketName string `json:"bucket_name" api:"required"`
	// When this entry is removed, in RFC 3339 format.
	ExpiresAt time.Time `json:"expires_at" api:"required" format:"date-time"`
	// The number of the Google Cloud project that holds the bucket, in decimal.
	ProjectNumber string `json:"project_number" api:"required"`
	// The code of the region that retired the bucket.
	Region string `json:"region" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		BucketName    respjson.Field
		ExpiresAt     respjson.Field
		ProjectNumber respjson.Field
		Region        respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedRetiredBucket) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedRetiredBucket) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A key that a region has retired. Keep the key usable until `expires_at`.
type BetaGoogleManagedRetiredKey struct {
	// When this entry is removed, in RFC 3339 format.
	ExpiresAt time.Time `json:"expires_at" api:"required" format:"date-time"`
	// The ID of the retired External Key (`ekey_...`).
	ExternalKeyID string `json:"external_key_id" api:"required"`
	// The code of the region that retired the key.
	Region string `json:"region" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ExpiresAt     respjson.Field
		ExternalKeyID respjson.Field
		Region        respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedRetiredKey) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedRetiredKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The bucket configured for a region.
type BetaGoogleManagedStorageBucket struct {
	// The bucket's Cloud Storage name.
	BucketName string `json:"bucket_name" api:"required"`
	// The resource name of the Cloud KMS key that the bucket encrypts objects under by
	// default. `null` until the region has confirmed the bucket, and when the bucket
	// has no default key.
	KMSKeyName string `json:"kms_key_name" api:"required"`
	// The bucket's location, as Cloud Storage reports it. `null` until the region has
	// confirmed the bucket.
	Location string `json:"location" api:"required"`
	// The ID of the project that holds the bucket. `null` until the region has
	// confirmed the bucket.
	ProjectID string `json:"project_id" api:"required"`
	// The number of the Google Cloud project that holds the bucket, in decimal.
	ProjectNumber string `json:"project_number" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		BucketName    respjson.Field
		KMSKeyName    respjson.Field
		Location      respjson.Field
		ProjectID     respjson.Field
		ProjectNumber respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedStorageBucket) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedStorageBucket) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaGoogleManagedStorageControlUnion contains all possible properties and values
// from [BetaGoogleManagedStorageDisabled], [BetaGoogleManagedStorageEnabled].
//
// Use the [BetaGoogleManagedStorageControlUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaGoogleManagedStorageControlUnion struct {
	RetiredBuckets []BetaGoogleManagedRetiredBucket `json:"retired_buckets"`
	// Any of "disabled", "enabled".
	Type string `json:"type"`
	// This field is from variant [BetaGoogleManagedStorageEnabled].
	Buckets map[string]BetaGoogleManagedStorageBucket `json:"buckets"`
	// This field is from variant [BetaGoogleManagedStorageEnabled].
	Status map[string]BetaGoogleManagedStorageRegionStatusUnion `json:"status"`
	JSON   struct {
		RetiredBuckets respjson.Field
		Type           respjson.Field
		Buckets        respjson.Field
		Status         respjson.Field
		raw            string
	} `json:"-"`
}

// anyBetaGoogleManagedStorageControl is implemented by each variant of
// [BetaGoogleManagedStorageControlUnion] to add type safety for the return type of
// [BetaGoogleManagedStorageControlUnion.AsAny]
type anyBetaGoogleManagedStorageControl interface {
	implBetaGoogleManagedStorageControlUnion()
}

func (BetaGoogleManagedStorageDisabled) implBetaGoogleManagedStorageControlUnion() {}
func (BetaGoogleManagedStorageEnabled) implBetaGoogleManagedStorageControlUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaGoogleManagedStorageControlUnion.AsAny().(type) {
//	case anthropic.BetaGoogleManagedStorageDisabled:
//	case anthropic.BetaGoogleManagedStorageEnabled:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaGoogleManagedStorageControlUnion) AsAny() anyBetaGoogleManagedStorageControl {
	switch u.Type {
	case "disabled":
		return u.AsDisabled()
	case "enabled":
		return u.AsEnabled()
	}
	return nil
}

func (u BetaGoogleManagedStorageControlUnion) AsDisabled() (v BetaGoogleManagedStorageDisabled) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedStorageControlUnion) AsEnabled() (v BetaGoogleManagedStorageEnabled) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaGoogleManagedStorageControlUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaGoogleManagedStorageControlUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaGoogleManagedStorageControlParamsUnion struct {
	OfDisabled *BetaGoogleManagedStorageDisabledParams `json:",omitzero,inline"`
	OfEnabled  *BetaGoogleManagedStorageEnabledParams  `json:",omitzero,inline"`
	paramUnion
}

func (u BetaGoogleManagedStorageControlParamsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDisabled, u.OfEnabled)
}
func (u *BetaGoogleManagedStorageControlParamsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaGoogleManagedStorageControlParamsUnion) asAny() any {
	if !param.IsOmitted(u.OfDisabled) {
		return u.OfDisabled
	} else if !param.IsOmitted(u.OfEnabled) {
		return u.OfEnabled
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaGoogleManagedStorageControlParamsUnion) GetBuckets() map[string]BetaGoogleManagedBucketReferenceParam {
	if vt := u.OfEnabled; vt != nil {
		return vt.Buckets
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaGoogleManagedStorageControlParamsUnion) GetType() *string {
	if vt := u.OfDisabled; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	} else if vt := u.OfEnabled; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaGoogleManagedStorageControlParamsUnion](
		"type",
		apijson.Discriminator[BetaGoogleManagedStorageDisabledParams]("disabled"),
		apijson.Discriminator[BetaGoogleManagedStorageEnabledParams]("enabled"),
	)
}

// The control is off and has no buckets.
type BetaGoogleManagedStorageDisabled struct {
	// Buckets that a region has retired, earliest `expires_at` first.
	RetiredBuckets []BetaGoogleManagedRetiredBucket `json:"retired_buckets" api:"required"`
	Type           constant.Disabled                `json:"type" default:"disabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		RetiredBuckets respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedStorageDisabled) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedStorageDisabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewBetaGoogleManagedStorageDisabledParams() BetaGoogleManagedStorageDisabledParams {
	return BetaGoogleManagedStorageDisabledParams{
		Type: "disabled",
	}
}

// Turns the control off and removes every region's bucket.
//
// This struct has a constant value, construct it with
// [NewBetaGoogleManagedStorageDisabledParams].
type BetaGoogleManagedStorageDisabledParams struct {
	Type constant.Disabled `json:"type" default:"disabled"`
	paramObj
}

func (r BetaGoogleManagedStorageDisabledParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaGoogleManagedStorageDisabledParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaGoogleManagedStorageDisabledParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The control is on and has a bucket for at least one region.
type BetaGoogleManagedStorageEnabled struct {
	// The bucket configured for each region, keyed by region code. A region with no
	// bucket is left out.
	Buckets map[string]BetaGoogleManagedStorageBucket `json:"buckets" api:"required"`
	// Buckets that a region has retired, earliest `expires_at` first.
	RetiredBuckets []BetaGoogleManagedRetiredBucket `json:"retired_buckets" api:"required"`
	// Each region's status for the control, keyed by region code. A region that your
	// organization is not offered may have no entry, and Anthropic may add regions.
	Status map[string]BetaGoogleManagedStorageRegionStatusUnion `json:"status" api:"required"`
	Type   constant.Enabled                                     `json:"type" default:"enabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Buckets        respjson.Field
		RetiredBuckets respjson.Field
		Status         respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedStorageEnabled) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedStorageEnabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Turns the control on, or changes its buckets while it is on.
//
// The property Type is required.
type BetaGoogleManagedStorageEnabledParams struct {
	// The bucket to use in each region, keyed by region code, such as `us` or `eu`. An
	// object with `project_number` and `bucket_name` sets that region's bucket, `null`
	// removes it, and a region you leave out keeps its bucket. A bucket can serve only
	// one region. While the control is on, at least one region must have a bucket.
	Buckets map[string]BetaGoogleManagedBucketReferenceParam `json:"buckets,omitzero"`
	// This field can be elided, and will marshal its zero value as "enabled".
	Type constant.Enabled `json:"type" default:"enabled"`
	paramObj
}

func (r BetaGoogleManagedStorageEnabledParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaGoogleManagedStorageEnabledParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaGoogleManagedStorageEnabledParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// One check that a verification attempt failed.
type BetaGoogleManagedStorageFailedCheck struct {
	// Which check failed, such as `label_mismatch` or `grant_mismatch`. Anthropic may
	// add values.
	//
	// Any of "label_mismatch", "location_mismatch", "encryption", "versioning",
	// "bucket_settings", "lifecycle_rules", "deletion_policy", "grant_mismatch",
	// "kms_key_refused", "probe_failed", "project_not_verified".
	Check BetaGoogleManagedStorageFailedCheckCheck `json:"check" api:"required"`
	// What failed and, where the end customer can fix it, what to change.
	Message string `json:"message" api:"required"`
	// On a check about one of the accounts that Anthropic uses to reach the bucket,
	// that account's kind of access, such as `ingest`. `null` on any other check.
	// Anthropic may add values.
	//
	// Any of "ingest", "read", "maintenance", "media_erasure", "verify".
	Role BetaGoogleManagedStorageFailedCheckRole `json:"role" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Check       respjson.Field
		Message     respjson.Field
		Role        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedStorageFailedCheck) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedStorageFailedCheck) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Which check failed, such as `label_mismatch` or `grant_mismatch`. Anthropic may
// add values.
type BetaGoogleManagedStorageFailedCheckCheck string

const (
	BetaGoogleManagedStorageFailedCheckCheckLabelMismatch      BetaGoogleManagedStorageFailedCheckCheck = "label_mismatch"
	BetaGoogleManagedStorageFailedCheckCheckLocationMismatch   BetaGoogleManagedStorageFailedCheckCheck = "location_mismatch"
	BetaGoogleManagedStorageFailedCheckCheckEncryption         BetaGoogleManagedStorageFailedCheckCheck = "encryption"
	BetaGoogleManagedStorageFailedCheckCheckVersioning         BetaGoogleManagedStorageFailedCheckCheck = "versioning"
	BetaGoogleManagedStorageFailedCheckCheckBucketSettings     BetaGoogleManagedStorageFailedCheckCheck = "bucket_settings"
	BetaGoogleManagedStorageFailedCheckCheckLifecycleRules     BetaGoogleManagedStorageFailedCheckCheck = "lifecycle_rules"
	BetaGoogleManagedStorageFailedCheckCheckDeletionPolicy     BetaGoogleManagedStorageFailedCheckCheck = "deletion_policy"
	BetaGoogleManagedStorageFailedCheckCheckGrantMismatch      BetaGoogleManagedStorageFailedCheckCheck = "grant_mismatch"
	BetaGoogleManagedStorageFailedCheckCheckKMSKeyRefused      BetaGoogleManagedStorageFailedCheckCheck = "kms_key_refused"
	BetaGoogleManagedStorageFailedCheckCheckProbeFailed        BetaGoogleManagedStorageFailedCheckCheck = "probe_failed"
	BetaGoogleManagedStorageFailedCheckCheckProjectNotVerified BetaGoogleManagedStorageFailedCheckCheck = "project_not_verified"
)

// On a check about one of the accounts that Anthropic uses to reach the bucket,
// that account's kind of access, such as `ingest`. `null` on any other check.
// Anthropic may add values.
type BetaGoogleManagedStorageFailedCheckRole string

const (
	BetaGoogleManagedStorageFailedCheckRoleIngest       BetaGoogleManagedStorageFailedCheckRole = "ingest"
	BetaGoogleManagedStorageFailedCheckRoleRead         BetaGoogleManagedStorageFailedCheckRole = "read"
	BetaGoogleManagedStorageFailedCheckRoleMaintenance  BetaGoogleManagedStorageFailedCheckRole = "maintenance"
	BetaGoogleManagedStorageFailedCheckRoleMediaErasure BetaGoogleManagedStorageFailedCheckRole = "media_erasure"
	BetaGoogleManagedStorageFailedCheckRoleVerify       BetaGoogleManagedStorageFailedCheckRole = "verify"
)

// The configured bucket passed verification, and the region applies it.
type BetaGoogleManagedStorageRegionActive struct {
	// The bucket that the region applies now. While a replacement bucket is `pending`
	// or `failed`, this is the previous bucket. `null` when the region applies no
	// bucket.
	EffectiveBucket BetaGoogleManagedBucketReference `json:"effective_bucket" api:"required"`
	// The checks that the latest reported daily verification failed. Empty when it
	// passed.
	FailedChecks []BetaGoogleManagedStorageFailedCheck `json:"failed_checks" api:"required"`
	Type         constant.Active                       `json:"type" default:"active"`
	// When the status last became `active`, in RFC 3339 format.
	ValidatedAt time.Time `json:"validated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EffectiveBucket respjson.Field
		FailedChecks    respjson.Field
		Type            respjson.Field
		ValidatedAt     respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedStorageRegionActive) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedStorageRegionActive) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The region stopped retrying verification before the configured bucket passed. To
// restart verification, send the region a bucket again. It can be the same bucket.
type BetaGoogleManagedStorageRegionFailed struct {
	// The bucket that the region applies now. While a replacement bucket is `pending`
	// or `failed`, this is the previous bucket. `null` when the region applies no
	// bucket.
	EffectiveBucket BetaGoogleManagedBucketReference `json:"effective_bucket" api:"required"`
	// The checks that the last reported verification attempt failed.
	FailedChecks []BetaGoogleManagedStorageFailedCheck `json:"failed_checks" api:"required"`
	Type         constant.Failed                       `json:"type" default:"failed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EffectiveBucket respjson.Field
		FailedChecks    respjson.Field
		Type            respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedStorageRegionFailed) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedStorageRegionFailed) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// No bucket is set for this region.
type BetaGoogleManagedStorageRegionNotConfigured struct {
	// The bucket that the region applies now. While a replacement bucket is `pending`
	// or `failed`, this is the previous bucket. `null` when the region applies no
	// bucket.
	EffectiveBucket BetaGoogleManagedBucketReference `json:"effective_bucket" api:"required"`
	Type            constant.NotConfigured           `json:"type" default:"not_configured"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EffectiveBucket respjson.Field
		Type            respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedStorageRegionNotConfigured) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedStorageRegionNotConfigured) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The region is verifying the configured bucket.
type BetaGoogleManagedStorageRegionPending struct {
	// The bucket that the region applies now. While a replacement bucket is `pending`
	// or `failed`, this is the previous bucket. `null` when the region applies no
	// bucket.
	EffectiveBucket BetaGoogleManagedBucketReference `json:"effective_bucket" api:"required"`
	// The checks that the latest reported verification attempt failed. Empty when none
	// has failed.
	FailedChecks []BetaGoogleManagedStorageFailedCheck `json:"failed_checks" api:"required"`
	Type         constant.Pending                      `json:"type" default:"pending"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EffectiveBucket respjson.Field
		FailedChecks    respjson.Field
		Type            respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedStorageRegionPending) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedStorageRegionPending) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaGoogleManagedStorageRegionStatusUnion contains all possible properties and
// values from [BetaGoogleManagedStorageRegionNotConfigured],
// [BetaGoogleManagedStorageRegionPending], [BetaGoogleManagedStorageRegionActive],
// [BetaGoogleManagedStorageRegionFailed],
// [BetaGoogleManagedStorageRegionSuspended].
//
// Use the [BetaGoogleManagedStorageRegionStatusUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaGoogleManagedStorageRegionStatusUnion struct {
	// This field is from variant [BetaGoogleManagedStorageRegionNotConfigured].
	EffectiveBucket BetaGoogleManagedBucketReference `json:"effective_bucket"`
	// Any of "not_configured", "pending", "active", "failed", "suspended".
	Type         string                                `json:"type"`
	FailedChecks []BetaGoogleManagedStorageFailedCheck `json:"failed_checks"`
	// This field is from variant [BetaGoogleManagedStorageRegionActive].
	ValidatedAt time.Time `json:"validated_at"`
	JSON        struct {
		EffectiveBucket respjson.Field
		Type            respjson.Field
		FailedChecks    respjson.Field
		ValidatedAt     respjson.Field
		raw             string
	} `json:"-"`
}

// anyBetaGoogleManagedStorageRegionStatus is implemented by each variant of
// [BetaGoogleManagedStorageRegionStatusUnion] to add type safety for the return
// type of [BetaGoogleManagedStorageRegionStatusUnion.AsAny]
type anyBetaGoogleManagedStorageRegionStatus interface {
	implBetaGoogleManagedStorageRegionStatusUnion()
}

func (BetaGoogleManagedStorageRegionNotConfigured) implBetaGoogleManagedStorageRegionStatusUnion() {}
func (BetaGoogleManagedStorageRegionPending) implBetaGoogleManagedStorageRegionStatusUnion()       {}
func (BetaGoogleManagedStorageRegionActive) implBetaGoogleManagedStorageRegionStatusUnion()        {}
func (BetaGoogleManagedStorageRegionFailed) implBetaGoogleManagedStorageRegionStatusUnion()        {}
func (BetaGoogleManagedStorageRegionSuspended) implBetaGoogleManagedStorageRegionStatusUnion()     {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaGoogleManagedStorageRegionStatusUnion.AsAny().(type) {
//	case anthropic.BetaGoogleManagedStorageRegionNotConfigured:
//	case anthropic.BetaGoogleManagedStorageRegionPending:
//	case anthropic.BetaGoogleManagedStorageRegionActive:
//	case anthropic.BetaGoogleManagedStorageRegionFailed:
//	case anthropic.BetaGoogleManagedStorageRegionSuspended:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaGoogleManagedStorageRegionStatusUnion) AsAny() anyBetaGoogleManagedStorageRegionStatus {
	switch u.Type {
	case "not_configured":
		return u.AsNotConfigured()
	case "pending":
		return u.AsPending()
	case "active":
		return u.AsActive()
	case "failed":
		return u.AsFailed()
	case "suspended":
		return u.AsSuspended()
	}
	return nil
}

func (u BetaGoogleManagedStorageRegionStatusUnion) AsNotConfigured() (v BetaGoogleManagedStorageRegionNotConfigured) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedStorageRegionStatusUnion) AsPending() (v BetaGoogleManagedStorageRegionPending) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedStorageRegionStatusUnion) AsActive() (v BetaGoogleManagedStorageRegionActive) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedStorageRegionStatusUnion) AsFailed() (v BetaGoogleManagedStorageRegionFailed) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedStorageRegionStatusUnion) AsSuspended() (v BetaGoogleManagedStorageRegionSuspended) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaGoogleManagedStorageRegionStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaGoogleManagedStorageRegionStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Anthropic has withdrawn or paused its approval of customer-owned storage for
// this profile, so the region does not apply the control.
type BetaGoogleManagedStorageRegionSuspended struct {
	// The bucket that the region applies now. While a replacement bucket is `pending`
	// or `failed`, this is the previous bucket. `null` when the region applies no
	// bucket.
	EffectiveBucket BetaGoogleManagedBucketReference `json:"effective_bucket" api:"required"`
	Type            constant.Suspended               `json:"type" default:"suspended"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EffectiveBucket respjson.Field
		Type            respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedStorageRegionSuspended) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedStorageRegionSuspended) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaUserProfileProviderGoogleDataControlGetParams struct {
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// This endpoint is in beta: requests must send
	// `google-managed-data-controls-2026-09-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaUserProfileProviderGoogleDataControlUpdateParams struct {
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// The customer-managed encryption key control. Leave it out to keep the control as
	// it is. Send `cmek`, `storage`, or both.
	CMEK BetaGoogleManagedCMEKControlParamsUnion `json:"cmek,omitzero"`
	// The customer-owned storage control. Leave it out to keep the control as it is.
	// Send `cmek`, `storage`, or both.
	Storage BetaGoogleManagedStorageControlParamsUnion `json:"storage,omitzero"`
	// This endpoint is in beta: requests must send
	// `google-managed-data-controls-2026-09-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaUserProfileProviderGoogleDataControlUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaUserProfileProviderGoogleDataControlUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaUserProfileProviderGoogleDataControlUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
