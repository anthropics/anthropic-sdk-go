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

// BetaOrganizationProviderGoogleProjectDataControlService contains methods and
// other services that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationProviderGoogleProjectDataControlService] method instead.
type BetaOrganizationProviderGoogleProjectDataControlService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationProviderGoogleProjectDataControlService generates a new
// service that applies the given options to each request. These options are
// applied after the parent client's options (if there is one), and before any
// request-specific options.
func NewBetaOrganizationProviderGoogleProjectDataControlService(opts ...option.RequestOption) (r BetaOrganizationProviderGoogleProjectDataControlService) {
	r = BetaOrganizationProviderGoogleProjectDataControlService{}
	r.Options = opts
	return
}

// Get a Google Cloud project's encryption keys and storage buckets for Claude on
// Google Cloud's Agent Platform, with each region's status.
//
// Your organization must have proved that it owns the project.
func (r *BetaOrganizationProviderGoogleProjectDataControlService) Get(ctx context.Context, projectNumber string, query BetaOrganizationProviderGoogleProjectDataControlGetParams, opts ...option.RequestOption) (res *BetaGoogleManagedProjectDataControls, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "google-managed-data-controls-2026-09-01")}, opts...)
	if projectNumber == "" {
		err = errors.New("missing required project_number parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/providers/google/projects/%s/data_controls", requestconfig.EscapePathSegment(projectNumber))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.providers.google.projects.data_controls.retrieve", "/v1/organizations/providers/google/projects/{project_number}/data_controls"), requestconfig.WithSpanAttribute("anthropic.organization.provider.google.project.id", projectNumber))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Turn a control on or off, or change the key or bucket set for a region. Controls
// and regions that you leave out stay as they are.
//
// Your organization must have proved that it owns the project. Only one
// organization at a time can have controls turned on for a project.
//
// **Requires a credential that identifies a user or a service account**: an OAuth
// access token with the `org:admin` scope, or an organization-scoped API key that
// belongs to a user or service account whose organization role can update these
// controls. The `admin`, `owner` and `primary_owner` roles can. An Admin API key
// gets a `403` error.
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeBetaCMEKNotApproved] (HTTP 400)
//   - [ErrorCodeBetaExternalKeyNotFound] (HTTP 400)
//   - [ErrorCodeBetaExternalKeyProviderMismatch] (HTTP 400)
//   - [ErrorCodeBetaGCPProjectNotVerified] (HTTP 400)
//   - [ErrorCodeKeyProfileMismatch] (HTTP 400)
//   - [ErrorCodeBetaOrganizationNotEligible] (HTTP 400)
//   - [ErrorCodeBetaRegionNotSupported] (HTTP 400)
//   - [ErrorCodeBetaStorageNotApproved] (HTTP 400)
//   - [ErrorCodeSandboxEgressRestricted] (HTTP 403)
//   - [ErrorCodeBetaGCPProjectControlledByAnotherOrganization] (HTTP 409)
//   - [ErrorCodeBetaStorageInUse] (HTTP 409)
//   - [ErrorCodeBetaAuthorizationCheckUnavailable] (HTTP 503)
func (r *BetaOrganizationProviderGoogleProjectDataControlService) Update(ctx context.Context, projectNumber string, params BetaOrganizationProviderGoogleProjectDataControlUpdateParams, opts ...option.RequestOption) (res *BetaGoogleManagedProjectDataControls, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "google-managed-data-controls-2026-09-01")}, opts...)
	if projectNumber == "" {
		err = errors.New("missing required project_number parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/providers/google/projects/%s/data_controls", requestconfig.EscapePathSegment(projectNumber))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.providers.google.projects.data_controls.update", "/v1/organizations/providers/google/projects/{project_number}/data_controls"), requestconfig.WithSpanAttribute("anthropic.organization.provider.google.project.id", projectNumber))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// A Cloud Storage bucket, named by the project that holds it and by its own name.
//
// The properties BucketName, ProjectNumber are required.
type BetaGoogleManagedProjectBucketReferenceParam struct {
	// The bucket's Cloud Storage name.
	BucketName string `json:"bucket_name" api:"required"`
	// The number of the Google Cloud project that holds the bucket, in decimal.
	ProjectNumber string `json:"project_number" api:"required"`
	paramObj
}

func (r BetaGoogleManagedProjectBucketReferenceParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaGoogleManagedProjectBucketReferenceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaGoogleManagedProjectBucketReferenceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaGoogleManagedProjectCMEKControlUnion contains all possible properties and
// values from [BetaGoogleManagedProjectCMEKDisabled],
// [BetaGoogleManagedProjectCMEKEnabled].
//
// Use the [BetaGoogleManagedProjectCMEKControlUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaGoogleManagedProjectCMEKControlUnion struct {
	RetiredKeys []BetaGoogleManagedProjectRetiredKey `json:"retired_keys"`
	// Any of "disabled", "enabled".
	Type string `json:"type"`
	// This field is from variant [BetaGoogleManagedProjectCMEKEnabled].
	ExternalKeyIDs map[string]string `json:"external_key_ids"`
	// This field is from variant [BetaGoogleManagedProjectCMEKEnabled].
	Status map[string]BetaGoogleManagedProjectCMEKRegionStatusUnion `json:"status"`
	JSON   struct {
		RetiredKeys    respjson.Field
		Type           respjson.Field
		ExternalKeyIDs respjson.Field
		Status         respjson.Field
		raw            string
	} `json:"-"`
}

// anyBetaGoogleManagedProjectCMEKControl is implemented by each variant of
// [BetaGoogleManagedProjectCMEKControlUnion] to add type safety for the return
// type of [BetaGoogleManagedProjectCMEKControlUnion.AsAny]
type anyBetaGoogleManagedProjectCMEKControl interface {
	implBetaGoogleManagedProjectCMEKControlUnion()
}

func (BetaGoogleManagedProjectCMEKDisabled) implBetaGoogleManagedProjectCMEKControlUnion() {}
func (BetaGoogleManagedProjectCMEKEnabled) implBetaGoogleManagedProjectCMEKControlUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaGoogleManagedProjectCMEKControlUnion.AsAny().(type) {
//	case anthropic.BetaGoogleManagedProjectCMEKDisabled:
//	case anthropic.BetaGoogleManagedProjectCMEKEnabled:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaGoogleManagedProjectCMEKControlUnion) AsAny() anyBetaGoogleManagedProjectCMEKControl {
	switch u.Type {
	case "disabled":
		return u.AsDisabled()
	case "enabled":
		return u.AsEnabled()
	}
	return nil
}

func (u BetaGoogleManagedProjectCMEKControlUnion) AsDisabled() (v BetaGoogleManagedProjectCMEKDisabled) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedProjectCMEKControlUnion) AsEnabled() (v BetaGoogleManagedProjectCMEKEnabled) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaGoogleManagedProjectCMEKControlUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaGoogleManagedProjectCMEKControlUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaGoogleManagedProjectCMEKControlParamsUnion struct {
	OfDisabled *BetaGoogleManagedProjectCMEKDisabledParams `json:",omitzero,inline"`
	OfEnabled  *BetaGoogleManagedProjectCMEKEnabledParams  `json:",omitzero,inline"`
	paramUnion
}

func (u BetaGoogleManagedProjectCMEKControlParamsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDisabled, u.OfEnabled)
}
func (u *BetaGoogleManagedProjectCMEKControlParamsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaGoogleManagedProjectCMEKControlParamsUnion) asAny() any {
	if !param.IsOmitted(u.OfDisabled) {
		return u.OfDisabled
	} else if !param.IsOmitted(u.OfEnabled) {
		return u.OfEnabled
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaGoogleManagedProjectCMEKControlParamsUnion) GetExternalKeyIDs() map[string]string {
	if vt := u.OfEnabled; vt != nil {
		return vt.ExternalKeyIDs
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaGoogleManagedProjectCMEKControlParamsUnion) GetType() *string {
	if vt := u.OfDisabled; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	} else if vt := u.OfEnabled; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaGoogleManagedProjectCMEKControlParamsUnion](
		"type",
		apijson.Discriminator[BetaGoogleManagedProjectCMEKDisabledParams]("disabled"),
		apijson.Discriminator[BetaGoogleManagedProjectCMEKEnabledParams]("enabled"),
	)
}

// The control is off and has no keys.
type BetaGoogleManagedProjectCMEKDisabled struct {
	// Keys that a region has retired, earliest `expires_at` first.
	RetiredKeys []BetaGoogleManagedProjectRetiredKey `json:"retired_keys" api:"required"`
	Type        constant.Disabled                    `json:"type" default:"disabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		RetiredKeys respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedProjectCMEKDisabled) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectCMEKDisabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewBetaGoogleManagedProjectCMEKDisabledParams() BetaGoogleManagedProjectCMEKDisabledParams {
	return BetaGoogleManagedProjectCMEKDisabledParams{
		Type: "disabled",
	}
}

// Turns the control off and removes every region's key.
//
// This struct has a constant value, construct it with
// [NewBetaGoogleManagedProjectCMEKDisabledParams].
type BetaGoogleManagedProjectCMEKDisabledParams struct {
	Type constant.Disabled `json:"type" default:"disabled"`
	paramObj
}

func (r BetaGoogleManagedProjectCMEKDisabledParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaGoogleManagedProjectCMEKDisabledParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaGoogleManagedProjectCMEKDisabledParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The control is on and has a key for at least one region.
type BetaGoogleManagedProjectCMEKEnabled struct {
	// The ID of the External Key configured for each region (`ekey_...`), keyed by
	// region code. A region with no key is left out.
	ExternalKeyIDs map[string]string `json:"external_key_ids" api:"required"`
	// Keys that a region has retired, earliest `expires_at` first.
	RetiredKeys []BetaGoogleManagedProjectRetiredKey `json:"retired_keys" api:"required"`
	// Each region's status for the control, keyed by region code. A region that your
	// organization is not offered may have no entry, and Anthropic may add regions.
	Status map[string]BetaGoogleManagedProjectCMEKRegionStatusUnion `json:"status" api:"required"`
	Type   constant.Enabled                                         `json:"type" default:"enabled"`
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
func (r BetaGoogleManagedProjectCMEKEnabled) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectCMEKEnabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Turns the control on, or changes its keys while it is on.
//
// The property Type is required.
type BetaGoogleManagedProjectCMEKEnabledParams struct {
	// The External Key to use in each region, keyed by region code, such as `us` or
	// `eu`. An External Key ID (`ekey_...`) sets that region's key, `null` removes it,
	// and a region you leave out keeps its key. While the control is on, at least one
	// region must have a key.
	ExternalKeyIDs map[string]string `json:"external_key_ids,omitzero"`
	// This field can be elided, and will marshal its zero value as "enabled".
	Type constant.Enabled `json:"type" default:"enabled"`
	paramObj
}

func (r BetaGoogleManagedProjectCMEKEnabledParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaGoogleManagedProjectCMEKEnabledParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaGoogleManagedProjectCMEKEnabledParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The configured key passed validation, and the region applies it.
type BetaGoogleManagedProjectCMEKRegionActive struct {
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
func (r BetaGoogleManagedProjectCMEKRegionActive) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectCMEKRegionActive) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The region stopped retrying validation before the configured key passed. To
// restart validation, send the region a key again. It can be the same key.
type BetaGoogleManagedProjectCMEKRegionFailed struct {
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
func (r BetaGoogleManagedProjectCMEKRegionFailed) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectCMEKRegionFailed) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// No key is set for this region.
type BetaGoogleManagedProjectCMEKRegionNotConfigured struct {
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
func (r BetaGoogleManagedProjectCMEKRegionNotConfigured) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectCMEKRegionNotConfigured) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The region is validating the configured key, or the key has passed validation
// and the region has not started to apply it.
type BetaGoogleManagedProjectCMEKRegionPending struct {
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
func (r BetaGoogleManagedProjectCMEKRegionPending) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectCMEKRegionPending) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaGoogleManagedProjectCMEKRegionStatusUnion contains all possible properties
// and values from [BetaGoogleManagedProjectCMEKRegionNotConfigured],
// [BetaGoogleManagedProjectCMEKRegionPending],
// [BetaGoogleManagedProjectCMEKRegionActive],
// [BetaGoogleManagedProjectCMEKRegionFailed],
// [BetaGoogleManagedProjectCMEKRegionSuspended].
//
// Use the [BetaGoogleManagedProjectCMEKRegionStatusUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaGoogleManagedProjectCMEKRegionStatusUnion struct {
	EffectiveExternalKeyID string `json:"effective_external_key_id"`
	// Any of "not_configured", "pending", "active", "failed", "suspended".
	Type  string `json:"type"`
	Error string `json:"error"`
	// This field is from variant [BetaGoogleManagedProjectCMEKRegionActive].
	LastCheckedAt time.Time `json:"last_checked_at"`
	// This field is from variant [BetaGoogleManagedProjectCMEKRegionActive].
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

// anyBetaGoogleManagedProjectCMEKRegionStatus is implemented by each variant of
// [BetaGoogleManagedProjectCMEKRegionStatusUnion] to add type safety for the
// return type of [BetaGoogleManagedProjectCMEKRegionStatusUnion.AsAny]
type anyBetaGoogleManagedProjectCMEKRegionStatus interface {
	implBetaGoogleManagedProjectCMEKRegionStatusUnion()
}

func (BetaGoogleManagedProjectCMEKRegionNotConfigured) implBetaGoogleManagedProjectCMEKRegionStatusUnion() {
}
func (BetaGoogleManagedProjectCMEKRegionPending) implBetaGoogleManagedProjectCMEKRegionStatusUnion() {
}
func (BetaGoogleManagedProjectCMEKRegionActive) implBetaGoogleManagedProjectCMEKRegionStatusUnion() {}
func (BetaGoogleManagedProjectCMEKRegionFailed) implBetaGoogleManagedProjectCMEKRegionStatusUnion() {}
func (BetaGoogleManagedProjectCMEKRegionSuspended) implBetaGoogleManagedProjectCMEKRegionStatusUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaGoogleManagedProjectCMEKRegionStatusUnion.AsAny().(type) {
//	case anthropic.BetaGoogleManagedProjectCMEKRegionNotConfigured:
//	case anthropic.BetaGoogleManagedProjectCMEKRegionPending:
//	case anthropic.BetaGoogleManagedProjectCMEKRegionActive:
//	case anthropic.BetaGoogleManagedProjectCMEKRegionFailed:
//	case anthropic.BetaGoogleManagedProjectCMEKRegionSuspended:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaGoogleManagedProjectCMEKRegionStatusUnion) AsAny() anyBetaGoogleManagedProjectCMEKRegionStatus {
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

func (u BetaGoogleManagedProjectCMEKRegionStatusUnion) AsNotConfigured() (v BetaGoogleManagedProjectCMEKRegionNotConfigured) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedProjectCMEKRegionStatusUnion) AsPending() (v BetaGoogleManagedProjectCMEKRegionPending) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedProjectCMEKRegionStatusUnion) AsActive() (v BetaGoogleManagedProjectCMEKRegionActive) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedProjectCMEKRegionStatusUnion) AsFailed() (v BetaGoogleManagedProjectCMEKRegionFailed) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedProjectCMEKRegionStatusUnion) AsSuspended() (v BetaGoogleManagedProjectCMEKRegionSuspended) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaGoogleManagedProjectCMEKRegionStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaGoogleManagedProjectCMEKRegionStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Anthropic has withdrawn or paused its approval of customer-managed keys for this
// project, so the region does not apply the control.
type BetaGoogleManagedProjectCMEKRegionSuspended struct {
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
func (r BetaGoogleManagedProjectCMEKRegionSuspended) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectCMEKRegionSuspended) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Google Cloud project's encryption key and storage bucket settings for Claude
// on Google Cloud's Agent Platform. Your organization has one of these objects for
// each Google Cloud project that it has proved it owns. You can read and update
// it, but not create or delete it. Keys and buckets are set for each region. A
// request's region depends on the Google endpoint it is sent to: `us` is the US
// multi-region endpoint and the regional endpoint of every location whose name
// starts with `us-`, `eu` is the EU multi-region endpoint and the regional
// endpoint of every location whose name starts with `europe-`, `sg` is the
// `asia-southeast1` regional endpoint, and `global` is the global endpoint,
// wherever Google serves the request. Anthropic may add regions.
type BetaGoogleManagedProjectDataControls struct {
	// The customer-managed encryption key control: the key in your Google Cloud KMS
	// that is set for each region. It covers the project's requests that name no user
	// profile, or a user profile that has no key of its own for the project.
	CMEK BetaGoogleManagedProjectCMEKControlUnion `json:"cmek" api:"required"`
	// The number of the Google Cloud project that these controls belong to, in
	// decimal.
	ProjectNumber string `json:"project_number" api:"required"`
	// The customer-owned storage control: your Cloud Storage bucket that is set for
	// each region. It covers the project's requests that name no user profile, or a
	// user profile that has no bucket of its own for the project.
	Storage BetaGoogleManagedProjectStorageControlUnion `json:"storage" api:"required"`
	// Object type. Always `google_managed_project_data_controls`.
	Type constant.GoogleManagedProjectDataControls `json:"type" default:"google_managed_project_data_controls"`
	// When an update last changed a control, or a region's status last changed, in RFC
	// 3339 format. A region that starts or stops reading `suspended`, a change to
	// `last_checked_at` alone, and the removal of a retired entry do not count.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CMEK          respjson.Field
		ProjectNumber respjson.Field
		Storage       respjson.Field
		Type          respjson.Field
		UpdatedAt     respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedProjectDataControls) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectDataControls) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Cloud Storage bucket, named by the project that holds it and by its own name.
type BetaGoogleManagedProjectEffectiveBucket struct {
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
func (r BetaGoogleManagedProjectEffectiveBucket) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectEffectiveBucket) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A bucket that a region has retired. Until `expires_at`, keep the bucket and
// Anthropic's access to it as they were when the bucket passed verification.
type BetaGoogleManagedProjectRetiredBucket struct {
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
func (r BetaGoogleManagedProjectRetiredBucket) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectRetiredBucket) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A key that a region has retired. Keep the key usable until `expires_at`.
type BetaGoogleManagedProjectRetiredKey struct {
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
func (r BetaGoogleManagedProjectRetiredKey) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectRetiredKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The bucket configured for a region.
type BetaGoogleManagedProjectStorageBucket struct {
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
func (r BetaGoogleManagedProjectStorageBucket) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectStorageBucket) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaGoogleManagedProjectStorageControlUnion contains all possible properties and
// values from [BetaGoogleManagedProjectStorageDisabled],
// [BetaGoogleManagedProjectStorageEnabled].
//
// Use the [BetaGoogleManagedProjectStorageControlUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaGoogleManagedProjectStorageControlUnion struct {
	RetiredBuckets []BetaGoogleManagedProjectRetiredBucket `json:"retired_buckets"`
	// Any of "disabled", "enabled".
	Type string `json:"type"`
	// This field is from variant [BetaGoogleManagedProjectStorageEnabled].
	Buckets map[string]BetaGoogleManagedProjectStorageBucket `json:"buckets"`
	// This field is from variant [BetaGoogleManagedProjectStorageEnabled].
	Status map[string]BetaGoogleManagedProjectStorageRegionStatusUnion `json:"status"`
	JSON   struct {
		RetiredBuckets respjson.Field
		Type           respjson.Field
		Buckets        respjson.Field
		Status         respjson.Field
		raw            string
	} `json:"-"`
}

// anyBetaGoogleManagedProjectStorageControl is implemented by each variant of
// [BetaGoogleManagedProjectStorageControlUnion] to add type safety for the return
// type of [BetaGoogleManagedProjectStorageControlUnion.AsAny]
type anyBetaGoogleManagedProjectStorageControl interface {
	implBetaGoogleManagedProjectStorageControlUnion()
}

func (BetaGoogleManagedProjectStorageDisabled) implBetaGoogleManagedProjectStorageControlUnion() {}
func (BetaGoogleManagedProjectStorageEnabled) implBetaGoogleManagedProjectStorageControlUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaGoogleManagedProjectStorageControlUnion.AsAny().(type) {
//	case anthropic.BetaGoogleManagedProjectStorageDisabled:
//	case anthropic.BetaGoogleManagedProjectStorageEnabled:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaGoogleManagedProjectStorageControlUnion) AsAny() anyBetaGoogleManagedProjectStorageControl {
	switch u.Type {
	case "disabled":
		return u.AsDisabled()
	case "enabled":
		return u.AsEnabled()
	}
	return nil
}

func (u BetaGoogleManagedProjectStorageControlUnion) AsDisabled() (v BetaGoogleManagedProjectStorageDisabled) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedProjectStorageControlUnion) AsEnabled() (v BetaGoogleManagedProjectStorageEnabled) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaGoogleManagedProjectStorageControlUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaGoogleManagedProjectStorageControlUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaGoogleManagedProjectStorageControlParamsUnion struct {
	OfDisabled *BetaGoogleManagedProjectStorageDisabledParams `json:",omitzero,inline"`
	OfEnabled  *BetaGoogleManagedProjectStorageEnabledParams  `json:",omitzero,inline"`
	paramUnion
}

func (u BetaGoogleManagedProjectStorageControlParamsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDisabled, u.OfEnabled)
}
func (u *BetaGoogleManagedProjectStorageControlParamsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaGoogleManagedProjectStorageControlParamsUnion) asAny() any {
	if !param.IsOmitted(u.OfDisabled) {
		return u.OfDisabled
	} else if !param.IsOmitted(u.OfEnabled) {
		return u.OfEnabled
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaGoogleManagedProjectStorageControlParamsUnion) GetBuckets() map[string]BetaGoogleManagedProjectBucketReferenceParam {
	if vt := u.OfEnabled; vt != nil {
		return vt.Buckets
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaGoogleManagedProjectStorageControlParamsUnion) GetType() *string {
	if vt := u.OfDisabled; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	} else if vt := u.OfEnabled; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaGoogleManagedProjectStorageControlParamsUnion](
		"type",
		apijson.Discriminator[BetaGoogleManagedProjectStorageDisabledParams]("disabled"),
		apijson.Discriminator[BetaGoogleManagedProjectStorageEnabledParams]("enabled"),
	)
}

// The control is off and has no buckets.
type BetaGoogleManagedProjectStorageDisabled struct {
	// Buckets that a region has retired, earliest `expires_at` first.
	RetiredBuckets []BetaGoogleManagedProjectRetiredBucket `json:"retired_buckets" api:"required"`
	Type           constant.Disabled                       `json:"type" default:"disabled"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		RetiredBuckets respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedProjectStorageDisabled) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectStorageDisabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewBetaGoogleManagedProjectStorageDisabledParams() BetaGoogleManagedProjectStorageDisabledParams {
	return BetaGoogleManagedProjectStorageDisabledParams{
		Type: "disabled",
	}
}

// Turns the control off and removes every region's bucket.
//
// This struct has a constant value, construct it with
// [NewBetaGoogleManagedProjectStorageDisabledParams].
type BetaGoogleManagedProjectStorageDisabledParams struct {
	Type constant.Disabled `json:"type" default:"disabled"`
	paramObj
}

func (r BetaGoogleManagedProjectStorageDisabledParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaGoogleManagedProjectStorageDisabledParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaGoogleManagedProjectStorageDisabledParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The control is on and has a bucket for at least one region.
type BetaGoogleManagedProjectStorageEnabled struct {
	// The bucket configured for each region, keyed by region code. A region with no
	// bucket is left out.
	Buckets map[string]BetaGoogleManagedProjectStorageBucket `json:"buckets" api:"required"`
	// Buckets that a region has retired, earliest `expires_at` first.
	RetiredBuckets []BetaGoogleManagedProjectRetiredBucket `json:"retired_buckets" api:"required"`
	// Each region's status for the control, keyed by region code. A region that your
	// organization is not offered may have no entry, and Anthropic may add regions.
	Status map[string]BetaGoogleManagedProjectStorageRegionStatusUnion `json:"status" api:"required"`
	Type   constant.Enabled                                            `json:"type" default:"enabled"`
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
func (r BetaGoogleManagedProjectStorageEnabled) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectStorageEnabled) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Turns the control on, or changes its buckets while it is on.
//
// The property Type is required.
type BetaGoogleManagedProjectStorageEnabledParams struct {
	// The bucket to use in each region, keyed by region code, such as `us` or `eu`. An
	// object with `project_number` and `bucket_name` sets that region's bucket, `null`
	// removes it, and a region you leave out keeps its bucket. A bucket can serve only
	// one region. While the control is on, at least one region must have a bucket.
	Buckets map[string]BetaGoogleManagedProjectBucketReferenceParam `json:"buckets,omitzero"`
	// This field can be elided, and will marshal its zero value as "enabled".
	Type constant.Enabled `json:"type" default:"enabled"`
	paramObj
}

func (r BetaGoogleManagedProjectStorageEnabledParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaGoogleManagedProjectStorageEnabledParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaGoogleManagedProjectStorageEnabledParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// One check that a verification attempt failed.
type BetaGoogleManagedProjectStorageFailedCheck struct {
	// Which check failed, such as `label_mismatch` or `grant_mismatch`. Anthropic may
	// add values.
	//
	// Any of "bucket_settings", "deletion_policy", "encryption", "grant_mismatch",
	// "kms_key_refused", "label_mismatch", "lifecycle_rules", "location_mismatch",
	// "probe_failed", "project_not_verified", "versioning".
	Check BetaGoogleManagedProjectStorageFailedCheckCheck `json:"check" api:"required"`
	// What failed and, where you can fix it, what to change.
	Message string `json:"message" api:"required"`
	// On a check about one of the accounts that Anthropic uses to reach the bucket,
	// that account's kind of access, such as `ingest`. `null` on any other check.
	// Anthropic may add values.
	//
	// Any of "ingest", "maintenance", "media_erasure", "read", "verify".
	Role BetaGoogleManagedProjectStorageFailedCheckRole `json:"role" api:"required"`
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
func (r BetaGoogleManagedProjectStorageFailedCheck) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectStorageFailedCheck) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Which check failed, such as `label_mismatch` or `grant_mismatch`. Anthropic may
// add values.
type BetaGoogleManagedProjectStorageFailedCheckCheck string

const (
	BetaGoogleManagedProjectStorageFailedCheckCheckBucketSettings     BetaGoogleManagedProjectStorageFailedCheckCheck = "bucket_settings"
	BetaGoogleManagedProjectStorageFailedCheckCheckDeletionPolicy     BetaGoogleManagedProjectStorageFailedCheckCheck = "deletion_policy"
	BetaGoogleManagedProjectStorageFailedCheckCheckEncryption         BetaGoogleManagedProjectStorageFailedCheckCheck = "encryption"
	BetaGoogleManagedProjectStorageFailedCheckCheckGrantMismatch      BetaGoogleManagedProjectStorageFailedCheckCheck = "grant_mismatch"
	BetaGoogleManagedProjectStorageFailedCheckCheckKMSKeyRefused      BetaGoogleManagedProjectStorageFailedCheckCheck = "kms_key_refused"
	BetaGoogleManagedProjectStorageFailedCheckCheckLabelMismatch      BetaGoogleManagedProjectStorageFailedCheckCheck = "label_mismatch"
	BetaGoogleManagedProjectStorageFailedCheckCheckLifecycleRules     BetaGoogleManagedProjectStorageFailedCheckCheck = "lifecycle_rules"
	BetaGoogleManagedProjectStorageFailedCheckCheckLocationMismatch   BetaGoogleManagedProjectStorageFailedCheckCheck = "location_mismatch"
	BetaGoogleManagedProjectStorageFailedCheckCheckProbeFailed        BetaGoogleManagedProjectStorageFailedCheckCheck = "probe_failed"
	BetaGoogleManagedProjectStorageFailedCheckCheckProjectNotVerified BetaGoogleManagedProjectStorageFailedCheckCheck = "project_not_verified"
	BetaGoogleManagedProjectStorageFailedCheckCheckVersioning         BetaGoogleManagedProjectStorageFailedCheckCheck = "versioning"
)

// On a check about one of the accounts that Anthropic uses to reach the bucket,
// that account's kind of access, such as `ingest`. `null` on any other check.
// Anthropic may add values.
type BetaGoogleManagedProjectStorageFailedCheckRole string

const (
	BetaGoogleManagedProjectStorageFailedCheckRoleIngest       BetaGoogleManagedProjectStorageFailedCheckRole = "ingest"
	BetaGoogleManagedProjectStorageFailedCheckRoleMaintenance  BetaGoogleManagedProjectStorageFailedCheckRole = "maintenance"
	BetaGoogleManagedProjectStorageFailedCheckRoleMediaErasure BetaGoogleManagedProjectStorageFailedCheckRole = "media_erasure"
	BetaGoogleManagedProjectStorageFailedCheckRoleRead         BetaGoogleManagedProjectStorageFailedCheckRole = "read"
	BetaGoogleManagedProjectStorageFailedCheckRoleVerify       BetaGoogleManagedProjectStorageFailedCheckRole = "verify"
)

// The configured bucket passed verification, and the region applies it.
type BetaGoogleManagedProjectStorageRegionActive struct {
	// The bucket that the region applies now. While a replacement bucket is `pending`
	// or `failed`, this is the previous bucket. `null` when the region applies no
	// bucket.
	EffectiveBucket BetaGoogleManagedProjectEffectiveBucket `json:"effective_bucket" api:"required"`
	// The checks that the latest reported daily verification failed. Empty when it
	// passed.
	FailedChecks []BetaGoogleManagedProjectStorageFailedCheck `json:"failed_checks" api:"required"`
	Type         constant.Active                              `json:"type" default:"active"`
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
func (r BetaGoogleManagedProjectStorageRegionActive) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectStorageRegionActive) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The region stopped retrying verification before the configured bucket passed. To
// restart verification, send the region a bucket again. It can be the same bucket.
type BetaGoogleManagedProjectStorageRegionFailed struct {
	// The bucket that the region applies now. While a replacement bucket is `pending`
	// or `failed`, this is the previous bucket. `null` when the region applies no
	// bucket.
	EffectiveBucket BetaGoogleManagedProjectEffectiveBucket `json:"effective_bucket" api:"required"`
	// The checks that the last reported verification attempt failed.
	FailedChecks []BetaGoogleManagedProjectStorageFailedCheck `json:"failed_checks" api:"required"`
	Type         constant.Failed                              `json:"type" default:"failed"`
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
func (r BetaGoogleManagedProjectStorageRegionFailed) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectStorageRegionFailed) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// No bucket is set for this region.
type BetaGoogleManagedProjectStorageRegionNotConfigured struct {
	// The bucket that the region applies now. While a replacement bucket is `pending`
	// or `failed`, this is the previous bucket. `null` when the region applies no
	// bucket.
	EffectiveBucket BetaGoogleManagedProjectEffectiveBucket `json:"effective_bucket" api:"required"`
	Type            constant.NotConfigured                  `json:"type" default:"not_configured"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EffectiveBucket respjson.Field
		Type            respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedProjectStorageRegionNotConfigured) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectStorageRegionNotConfigured) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The region is verifying the configured bucket.
type BetaGoogleManagedProjectStorageRegionPending struct {
	// The bucket that the region applies now. While a replacement bucket is `pending`
	// or `failed`, this is the previous bucket. `null` when the region applies no
	// bucket.
	EffectiveBucket BetaGoogleManagedProjectEffectiveBucket `json:"effective_bucket" api:"required"`
	// The checks that the latest reported verification attempt failed. Empty when none
	// has failed.
	FailedChecks []BetaGoogleManagedProjectStorageFailedCheck `json:"failed_checks" api:"required"`
	Type         constant.Pending                             `json:"type" default:"pending"`
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
func (r BetaGoogleManagedProjectStorageRegionPending) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectStorageRegionPending) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaGoogleManagedProjectStorageRegionStatusUnion contains all possible
// properties and values from [BetaGoogleManagedProjectStorageRegionNotConfigured],
// [BetaGoogleManagedProjectStorageRegionPending],
// [BetaGoogleManagedProjectStorageRegionActive],
// [BetaGoogleManagedProjectStorageRegionFailed],
// [BetaGoogleManagedProjectStorageRegionSuspended].
//
// Use the [BetaGoogleManagedProjectStorageRegionStatusUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaGoogleManagedProjectStorageRegionStatusUnion struct {
	// This field is from variant [BetaGoogleManagedProjectStorageRegionNotConfigured].
	EffectiveBucket BetaGoogleManagedProjectEffectiveBucket `json:"effective_bucket"`
	// Any of "not_configured", "pending", "active", "failed", "suspended".
	Type         string                                       `json:"type"`
	FailedChecks []BetaGoogleManagedProjectStorageFailedCheck `json:"failed_checks"`
	// This field is from variant [BetaGoogleManagedProjectStorageRegionActive].
	ValidatedAt time.Time `json:"validated_at"`
	JSON        struct {
		EffectiveBucket respjson.Field
		Type            respjson.Field
		FailedChecks    respjson.Field
		ValidatedAt     respjson.Field
		raw             string
	} `json:"-"`
}

// anyBetaGoogleManagedProjectStorageRegionStatus is implemented by each variant of
// [BetaGoogleManagedProjectStorageRegionStatusUnion] to add type safety for the
// return type of [BetaGoogleManagedProjectStorageRegionStatusUnion.AsAny]
type anyBetaGoogleManagedProjectStorageRegionStatus interface {
	implBetaGoogleManagedProjectStorageRegionStatusUnion()
}

func (BetaGoogleManagedProjectStorageRegionNotConfigured) implBetaGoogleManagedProjectStorageRegionStatusUnion() {
}
func (BetaGoogleManagedProjectStorageRegionPending) implBetaGoogleManagedProjectStorageRegionStatusUnion() {
}
func (BetaGoogleManagedProjectStorageRegionActive) implBetaGoogleManagedProjectStorageRegionStatusUnion() {
}
func (BetaGoogleManagedProjectStorageRegionFailed) implBetaGoogleManagedProjectStorageRegionStatusUnion() {
}
func (BetaGoogleManagedProjectStorageRegionSuspended) implBetaGoogleManagedProjectStorageRegionStatusUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaGoogleManagedProjectStorageRegionStatusUnion.AsAny().(type) {
//	case anthropic.BetaGoogleManagedProjectStorageRegionNotConfigured:
//	case anthropic.BetaGoogleManagedProjectStorageRegionPending:
//	case anthropic.BetaGoogleManagedProjectStorageRegionActive:
//	case anthropic.BetaGoogleManagedProjectStorageRegionFailed:
//	case anthropic.BetaGoogleManagedProjectStorageRegionSuspended:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaGoogleManagedProjectStorageRegionStatusUnion) AsAny() anyBetaGoogleManagedProjectStorageRegionStatus {
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

func (u BetaGoogleManagedProjectStorageRegionStatusUnion) AsNotConfigured() (v BetaGoogleManagedProjectStorageRegionNotConfigured) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedProjectStorageRegionStatusUnion) AsPending() (v BetaGoogleManagedProjectStorageRegionPending) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedProjectStorageRegionStatusUnion) AsActive() (v BetaGoogleManagedProjectStorageRegionActive) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedProjectStorageRegionStatusUnion) AsFailed() (v BetaGoogleManagedProjectStorageRegionFailed) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaGoogleManagedProjectStorageRegionStatusUnion) AsSuspended() (v BetaGoogleManagedProjectStorageRegionSuspended) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaGoogleManagedProjectStorageRegionStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaGoogleManagedProjectStorageRegionStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Anthropic has withdrawn or paused its approval of customer-owned storage for
// this project, so the region does not apply the control.
type BetaGoogleManagedProjectStorageRegionSuspended struct {
	// The bucket that the region applies now. While a replacement bucket is `pending`
	// or `failed`, this is the previous bucket. `null` when the region applies no
	// bucket.
	EffectiveBucket BetaGoogleManagedProjectEffectiveBucket `json:"effective_bucket" api:"required"`
	Type            constant.Suspended                      `json:"type" default:"suspended"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EffectiveBucket respjson.Field
		Type            respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGoogleManagedProjectStorageRegionSuspended) RawJSON() string { return r.JSON.raw }
func (r *BetaGoogleManagedProjectStorageRegionSuspended) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationProviderGoogleProjectDataControlGetParams struct {
	// This endpoint is in beta: requests must send
	// `google-managed-data-controls-2026-09-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationProviderGoogleProjectDataControlUpdateParams struct {
	// The customer-managed encryption key control. Leave it out to keep the control as
	// it is. Send `cmek`, `storage`, or both.
	CMEK BetaGoogleManagedProjectCMEKControlParamsUnion `json:"cmek,omitzero"`
	// The customer-owned storage control. Leave it out to keep the control as it is.
	// Send `cmek`, `storage`, or both.
	Storage BetaGoogleManagedProjectStorageControlParamsUnion `json:"storage,omitzero"`
	// This endpoint is in beta: requests must send
	// `google-managed-data-controls-2026-09-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationProviderGoogleProjectDataControlUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationProviderGoogleProjectDataControlUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationProviderGoogleProjectDataControlUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
