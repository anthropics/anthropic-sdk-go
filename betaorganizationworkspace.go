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
	"github.com/anthropics/anthropic-sdk-go/v2/internal/paramutil"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// BetaOrganizationWorkspaceService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationWorkspaceService] method instead.
type BetaOrganizationWorkspaceService struct {
	Options         []option.RequestOption
	RateLimits      BetaOrganizationWorkspaceRateLimitService
	Members         BetaOrganizationWorkspaceMemberService
	ServiceAccounts BetaOrganizationWorkspaceServiceAccountService
}

// NewBetaOrganizationWorkspaceService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationWorkspaceService(opts ...option.RequestOption) (r BetaOrganizationWorkspaceService) {
	r = BetaOrganizationWorkspaceService{}
	r.Options = opts
	r.RateLimits = NewBetaOrganizationWorkspaceRateLimitService(opts...)
	r.Members = NewBetaOrganizationWorkspaceMemberService(opts...)
	r.ServiceAccounts = NewBetaOrganizationWorkspaceServiceAccountService(opts...)
	return
}

// Create Workspace
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeBetaApprovalRequired] (HTTP 400)
//   - [ErrorCodeExternalKeyValidationFailed] (HTTP 400)
//   - [ErrorCodeExternalKeyValidationTimedOut] (HTTP 400)
//   - [ErrorCodeBetaExternalStorageNotAvailable] (HTTP 400)
//   - [ErrorCodeBetaGrantRequired] (HTTP 400)
//   - [ErrorCodeBetaHumanReviewNotAvailable] (HTTP 400)
//   - [ErrorCodeBetaHumanReviewOnWithExternalStorage] (HTTP 400)
//   - [ErrorCodeBetaInvalidStorageARN] (HTTP 400)
//   - [ErrorCodeKeyProfileMismatch] (HTTP 400)
//   - [ErrorCodeBetaOrganizationNotEligible] (HTTP 400)
//   - [ErrorCodeBetaStorageNotAttachable] (HTTP 400)
//   - [ErrorCodeStorageProfileMismatch] (HTTP 400)
//   - [ErrorCodeBetaStorageVerificationFailed] (HTTP 400)
//   - [ErrorCodeBetaTranscriptsStorageRequired] (HTTP 400)
//   - [ErrorCodeUserProfileNotFound] (HTTP 400)
//   - [ErrorCodeWorkspacePolicyViolation] (HTTP 400)
//   - [ErrorCodeSandboxEgressRestricted] (HTTP 403)
//   - [ErrorCodeExternalKeyChanged] (HTTP 409)
//   - [ErrorCodeBetaAuthorizationCheckUnavailable] (HTTP 503)
//   - [ErrorCodeWorkspaceCreatedIncomplete] (HTTP 503)
func (r *BetaOrganizationWorkspaceService) New(ctx context.Context, params BetaOrganizationWorkspaceNewParams, opts ...option.RequestOption) (res *BetaWorkspace, err error) {
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
	path := "v1/organizations/workspaces"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithConservativeRetries())
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.workspaces.create", "/v1/organizations/workspaces"), requestconfig.WithSpanCreatedID("anthropic.organization.workspace.id", &res, func(created *BetaWorkspace) string { return created.ID }))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get Workspace
func (r *BetaOrganizationWorkspaceService) Get(ctx context.Context, workspaceID string, query BetaOrganizationWorkspaceGetParams, opts ...option.RequestOption) (res *BetaWorkspace, err error) {
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
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s", requestconfig.EscapePathSegment(workspaceID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.workspaces.retrieve", "/v1/organizations/workspaces/{workspace_id}"), requestconfig.WithSpanAttribute("anthropic.organization.workspace.id", workspaceID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Workspace
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeBetaApprovalRequired] (HTTP 400)
//   - [ErrorCodeExternalKeyValidationFailed] (HTTP 400)
//   - [ErrorCodeExternalKeyValidationTimedOut] (HTTP 400)
//   - [ErrorCodeBetaExternalStorageNotAvailable] (HTTP 400)
//   - [ErrorCodeBetaGrantRequired] (HTTP 400)
//   - [ErrorCodeBetaHumanReviewNotAvailable] (HTTP 400)
//   - [ErrorCodeBetaHumanReviewOnWithExternalStorage] (HTTP 400)
//   - [ErrorCodeBetaInvalidStorageARN] (HTTP 400)
//   - [ErrorCodeKeyProfileMismatch] (HTTP 400)
//   - [ErrorCodeBetaOrganizationNotEligible] (HTTP 400)
//   - [ErrorCodeBetaStorageAttachLimitReached] (HTTP 400)
//   - [ErrorCodeBetaStorageNotAttachable] (HTTP 400)
//   - [ErrorCodeStorageProfileMismatch] (HTTP 400)
//   - [ErrorCodeBetaStorageProfileStillRequired] (HTTP 400)
//   - [ErrorCodeBetaStorageVerificationFailed] (HTTP 400)
//   - [ErrorCodeBetaTranscriptsStorageRequired] (HTTP 400)
//   - [ErrorCodeUserProfileNotFound] (HTTP 400)
//   - [ErrorCodeWorkspacePolicyViolation] (HTTP 400)
//   - [ErrorCodeSandboxEgressRestricted] (HTTP 403)
//   - [ErrorCodeExternalKeyChanged] (HTTP 409)
//   - [ErrorCodeBetaStorageChangedConcurrently] (HTTP 409)
//   - [ErrorCodeBetaWorkspaceChangedConcurrently] (HTTP 409)
//   - [ErrorCodeBetaAuthorizationCheckUnavailable] (HTTP 503)
//   - [ErrorCodeChangeOutcomeUnknown] (HTTP 503)
func (r *BetaOrganizationWorkspaceService) Update(ctx context.Context, workspaceID string, params BetaOrganizationWorkspaceUpdateParams, opts ...option.RequestOption) (res *BetaWorkspace, err error) {
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
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s", requestconfig.EscapePathSegment(workspaceID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.workspaces.update", "/v1/organizations/workspaces/{workspace_id}"), requestconfig.WithSpanAttribute("anthropic.organization.workspace.id", workspaceID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List Workspaces
func (r *BetaOrganizationWorkspaceService) List(ctx context.Context, params BetaOrganizationWorkspaceListParams, opts ...option.RequestOption) (res *pagination.Page[BetaWorkspace], err error) {
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
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/workspaces"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.workspaces.list", "/v1/organizations/workspaces", "after_id", "before_id", "ids", "include_archived", "include_default", "limit"))
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

// List Workspaces
func (r *BetaOrganizationWorkspaceService) ListAutoPaging(ctx context.Context, params BetaOrganizationWorkspaceListParams, opts ...option.RequestOption) *pagination.PageAutoPager[BetaWorkspace] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

// Archive Workspace
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeSandboxEgressRestricted] (HTTP 403)
//   - [ErrorCodeStorageAttachedDuringArchive] (HTTP 409)
func (r *BetaOrganizationWorkspaceService) Archive(ctx context.Context, workspaceID string, body BetaOrganizationWorkspaceArchiveParams, opts ...option.RequestOption) (res *BetaWorkspace, err error) {
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
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s/archive", requestconfig.EscapePathSegment(workspaceID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.workspaces.archive", "/v1/organizations/workspaces/{workspace_id}/archive"), requestconfig.WithSpanAttribute("anthropic.organization.workspace.id", workspaceID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type BetaAllowedInferenceGeo string

const (
	BetaAllowedInferenceGeoGlobal BetaAllowedInferenceGeo = "global"
	BetaAllowedInferenceGeoUs     BetaAllowedInferenceGeo = "us"
)

// The Workspace's human review setting.
type BetaAnthropicHumanReview struct {
	// Any of "disabled", "enabled".
	Type BetaAnthropicHumanReviewType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnthropicHumanReview) RawJSON() string { return r.JSON.raw }
func (r *BetaAnthropicHumanReview) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAnthropicHumanReviewType string

const (
	BetaAnthropicHumanReviewTypeDisabled BetaAnthropicHumanReviewType = "disabled"
	BetaAnthropicHumanReviewTypeEnabled  BetaAnthropicHumanReviewType = "enabled"
)

// Whether Anthropic's human review applies to the Workspace's inference data.
//
// The property Type is required.
type BetaAnthropicHumanReviewParam struct {
	// Any of "disabled", "enabled".
	Type BetaAnthropicHumanReviewParamType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaAnthropicHumanReviewParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaAnthropicHumanReviewParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAnthropicHumanReviewParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAnthropicHumanReviewParamType string

const (
	BetaAnthropicHumanReviewParamTypeDisabled BetaAnthropicHumanReviewParamType = "disabled"
	BetaAnthropicHumanReviewParamTypeEnabled  BetaAnthropicHumanReviewParamType = "enabled"
)

type BetaDataResidency struct {
	// Permitted inference geo values. 'unrestricted' means all geos are allowed.
	AllowedInferenceGeos BetaDataResidencyAllowedInferenceGeosUnion `json:"allowed_inference_geos" api:"required"`
	// Default inference geo applied when requests omit the parameter.
	//
	// Any of "global", "us".
	DefaultInferenceGeo BetaDataResidencyDefaultInferenceGeo `json:"default_inference_geo" api:"required"`
	// Geographic region for workspace data storage. Immutable after creation.
	//
	// Any of "us".
	WorkspaceGeo BetaDataResidencyWorkspaceGeo `json:"workspace_geo" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AllowedInferenceGeos respjson.Field
		DefaultInferenceGeo  respjson.Field
		WorkspaceGeo         respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDataResidency) RawJSON() string { return r.JSON.raw }
func (r *BetaDataResidency) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaDataResidencyAllowedInferenceGeosUnion contains all possible properties and
// values from [[]BetaAllowedInferenceGeo], [constant.Unrestricted].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfGeos OfUnrestricted]
type BetaDataResidencyAllowedInferenceGeosUnion struct {
	// This field will be present if the value is a [[]BetaAllowedInferenceGeo] instead
	// of an object.
	OfGeos []BetaAllowedInferenceGeo `json:",inline"`
	// This field will be present if the value is a [constant.Unrestricted] instead of
	// an object.
	OfUnrestricted constant.Unrestricted `json:",inline"`
	JSON           struct {
		OfGeos         respjson.Field
		OfUnrestricted respjson.Field
		raw            string
	} `json:"-"`
}

func (u BetaDataResidencyAllowedInferenceGeosUnion) AsGeos() (v []BetaAllowedInferenceGeo) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaDataResidencyAllowedInferenceGeosUnion) AsUnrestricted() (v constant.Unrestricted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaDataResidencyAllowedInferenceGeosUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaDataResidencyAllowedInferenceGeosUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Default inference geo applied when requests omit the parameter.
type BetaDataResidencyDefaultInferenceGeo string

const (
	BetaDataResidencyDefaultInferenceGeoGlobal BetaDataResidencyDefaultInferenceGeo = "global"
	BetaDataResidencyDefaultInferenceGeoUs     BetaDataResidencyDefaultInferenceGeo = "us"
)

// Geographic region for workspace data storage. Immutable after creation.
type BetaDataResidencyWorkspaceGeo string

const (
	BetaDataResidencyWorkspaceGeoUs BetaDataResidencyWorkspaceGeo = "us"
)

type BetaDataResidencyCreateConfigParam struct {
	// Permitted inference geo values. Defaults to 'unrestricted' if omitted, which
	// allows all geos. Use the string 'unrestricted' to allow all geos, or a list of
	// specific geos.
	AllowedInferenceGeos BetaDataResidencyCreateConfigAllowedInferenceGeosUnionParam `json:"allowed_inference_geos,omitzero"`
	// Default inference geo applied when requests omit the parameter. Defaults to
	// 'global' if omitted. Must be a member of `allowed_inference_geos` unless
	// `allowed_inference_geos` is `"unrestricted"`.
	//
	// Any of "global", "us".
	DefaultInferenceGeo BetaDataResidencyCreateConfigDefaultInferenceGeo `json:"default_inference_geo,omitzero"`
	// Geographic region for workspace data storage. Immutable after creation. Defaults
	// to 'us' if omitted.
	//
	// Any of "us".
	WorkspaceGeo BetaDataResidencyCreateConfigWorkspaceGeo `json:"workspace_geo,omitzero"`
	paramObj
}

func (r BetaDataResidencyCreateConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaDataResidencyCreateConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaDataResidencyCreateConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaDataResidencyCreateConfigAllowedInferenceGeosUnionParam struct {
	OfGeos []BetaAllowedInferenceGeo `json:",omitzero,inline"`
	// Construct this variant with constant.ValueOf[constant.Unrestricted]()
	OfUnrestricted constant.Unrestricted `json:",omitzero,inline"`
	paramUnion
}

func (u BetaDataResidencyCreateConfigAllowedInferenceGeosUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfGeos, u.OfUnrestricted)
}
func (u *BetaDataResidencyCreateConfigAllowedInferenceGeosUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaDataResidencyCreateConfigAllowedInferenceGeosUnionParam) asAny() any {
	if !param.IsOmitted(u.OfGeos) {
		return &u.OfGeos
	} else if !param.IsOmitted(u.OfUnrestricted) {
		return &u.OfUnrestricted
	}
	return nil
}

// Default inference geo applied when requests omit the parameter. Defaults to
// 'global' if omitted. Must be a member of `allowed_inference_geos` unless
// `allowed_inference_geos` is `"unrestricted"`.
type BetaDataResidencyCreateConfigDefaultInferenceGeo string

const (
	BetaDataResidencyCreateConfigDefaultInferenceGeoGlobal BetaDataResidencyCreateConfigDefaultInferenceGeo = "global"
	BetaDataResidencyCreateConfigDefaultInferenceGeoUs     BetaDataResidencyCreateConfigDefaultInferenceGeo = "us"
)

// Geographic region for workspace data storage. Immutable after creation. Defaults
// to 'us' if omitted.
type BetaDataResidencyCreateConfigWorkspaceGeo string

const (
	BetaDataResidencyCreateConfigWorkspaceGeoUs BetaDataResidencyCreateConfigWorkspaceGeo = "us"
)

type BetaDataResidencyUpdateConfigParam struct {
	// Permitted inference geo values. Use 'unrestricted' to allow all geos, or a list
	// of specific geos.
	AllowedInferenceGeos BetaDataResidencyUpdateConfigAllowedInferenceGeosUnionParam `json:"allowed_inference_geos,omitzero"`
	// Default inference geo applied when requests omit the parameter. Must be a member
	// of `allowed_inference_geos` unless `allowed_inference_geos` is `"unrestricted"`.
	//
	// Any of "global", "us".
	DefaultInferenceGeo BetaDataResidencyUpdateConfigDefaultInferenceGeo `json:"default_inference_geo,omitzero"`
	paramObj
}

func (r BetaDataResidencyUpdateConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaDataResidencyUpdateConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaDataResidencyUpdateConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaDataResidencyUpdateConfigAllowedInferenceGeosUnionParam struct {
	OfGeos []BetaAllowedInferenceGeo `json:",omitzero,inline"`
	// Construct this variant with constant.ValueOf[constant.Unrestricted]()
	OfUnrestricted constant.Unrestricted `json:",omitzero,inline"`
	paramUnion
}

func (u BetaDataResidencyUpdateConfigAllowedInferenceGeosUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfGeos, u.OfUnrestricted)
}
func (u *BetaDataResidencyUpdateConfigAllowedInferenceGeosUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaDataResidencyUpdateConfigAllowedInferenceGeosUnionParam) asAny() any {
	if !param.IsOmitted(u.OfGeos) {
		return &u.OfGeos
	} else if !param.IsOmitted(u.OfUnrestricted) {
		return &u.OfUnrestricted
	}
	return nil
}

// Default inference geo applied when requests omit the parameter. Must be a member
// of `allowed_inference_geos` unless `allowed_inference_geos` is `"unrestricted"`.
type BetaDataResidencyUpdateConfigDefaultInferenceGeo string

const (
	BetaDataResidencyUpdateConfigDefaultInferenceGeoGlobal BetaDataResidencyUpdateConfigDefaultInferenceGeo = "global"
	BetaDataResidencyUpdateConfigDefaultInferenceGeoUs     BetaDataResidencyUpdateConfigDefaultInferenceGeo = "us"
)

// The organization's External Storage that `external_storage_id` names.
type BetaExternalStorageMember struct {
	// The External Storage ID, which starts with `estrg_`. Through Claude Platform on
	// AWS, storage whose `provider_config.type` is `aws` has its bucket ARN as its ID
	// instead: `arn:<partition>:s3:::<bucket_name>`, where `<partition>` is the
	// partition of the bucket's region.
	ExternalStorageID string            `json:"external_storage_id" api:"required"`
	Type              constant.External `json:"type" default:"external"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ExternalStorageID respjson.Field
		Type              respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaExternalStorageMember) RawJSON() string { return r.JSON.raw }
func (r *BetaExternalStorageMember) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The organization's External Storage that `external_storage_id` names.
//
// The properties ExternalStorageID, Type are required.
type BetaExternalStorageMemberParam struct {
	// The External Storage ID, which starts with `estrg_`. Through Claude Platform on
	// AWS, storage whose `provider_config.type` is `aws` has its bucket ARN as its ID
	// instead: `arn:<partition>:s3:::<bucket_name>`, where `<partition>` is the
	// partition of the bucket's region.
	ExternalStorageID string `json:"external_storage_id" api:"required"`
	// This field can be elided, and will marshal its zero value as "external".
	Type constant.External `json:"type" default:"external"`
	paramObj
}

func (r BetaExternalStorageMemberParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaExternalStorageMemberParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaExternalStorageMemberParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaInferenceDataPrivacy struct {
	// The Workspace's human review setting. Human review is off only while this is
	// `disabled` and Anthropic's approval stands.
	HumanReview BetaAnthropicHumanReview `json:"human_review"`
	// Where the data Anthropic's safeguards systems retain for the Workspace goes.
	// `external` sends it to the External Storage that `external_storage_id` names.
	// `internal` means no External Storage: it goes to Anthropic's storage or, on a
	// cloud provider's platform, to that provider's storage. A recent change may not
	// have taken effect yet.
	SafeguardsStorage BetaStorageMemberUnion `json:"safeguards_storage"`
	// Where copies of the Workspace's transcripts and its cases (Anthropic's safety
	// findings about those transcripts) go. `external` sends them to the External
	// Storage that `external_storage_id` names. `internal` means no External Storage,
	// so no copies are made. A recent change may not have taken effect yet.
	TranscriptsStorage BetaStorageMemberUnion `json:"transcripts_storage"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		HumanReview        respjson.Field
		SafeguardsStorage  respjson.Field
		TranscriptsStorage respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaInferenceDataPrivacy) RawJSON() string { return r.JSON.raw }
func (r *BetaInferenceDataPrivacy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Where the Workspace's retained inference data and transcript copies go, and
// whether Anthropic's human review applies. A member cannot be null; omit it
// instead.
type BetaInferenceDataPrivacyCreateParam struct {
	// Whether Anthropic's human review applies to the Workspace. `disabled` needs
	// Anthropic's approval. When omitted, `enabled`. Cannot be null.
	HumanReview BetaAnthropicHumanReviewParam `json:"human_review,omitzero"`
	// Where the data Anthropic's safeguards systems retain for the Workspace goes.
	// `external` sends it, once the attach takes effect, to the External Storage that
	// `external_storage_id` names. `internal` means no External Storage: it goes to
	// Anthropic's storage or, on a cloud provider's platform, to that provider's
	// storage. When omitted, `internal`. Cannot be null.
	SafeguardsStorage BetaStorageMemberParamUnion `json:"safeguards_storage,omitzero"`
	// Where copies of the Workspace's transcripts and its cases (Anthropic's safety
	// findings about those transcripts) go. `external` sends them, once the attach
	// takes effect, to the External Storage that `external_storage_id` names.
	// `internal` means no External Storage, so no copies are made. When omitted,
	// `internal`. Cannot be null.
	TranscriptsStorage BetaStorageMemberParamUnion `json:"transcripts_storage,omitzero"`
	paramObj
}

func (r BetaInferenceDataPrivacyCreateParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaInferenceDataPrivacyCreateParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaInferenceDataPrivacyCreateParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Changes to where the Workspace's retained inference data and transcript copies
// go, and to whether Anthropic's human review applies. An omitted member is
// unchanged; a member cannot be null.
type BetaInferenceDataPrivacyUpdateParam struct {
	// Whether Anthropic's human review applies to the Workspace from now on.
	// `disabled` needs Anthropic's approval. When omitted, unchanged. Cannot be null.
	HumanReview BetaAnthropicHumanReviewParam `json:"human_review,omitzero"`
	// Where the data Anthropic's safeguards systems retain for the Workspace goes once
	// the change takes effect. `external` sends it to the External Storage that
	// `external_storage_id` names. `internal` means no External Storage: it goes to
	// Anthropic's storage or, on a cloud provider's platform, to that provider's
	// storage. When omitted, unchanged. Cannot be null.
	SafeguardsStorage BetaStorageMemberParamUnion `json:"safeguards_storage,omitzero"`
	// Where copies of the Workspace's transcripts and its cases (Anthropic's safety
	// findings about those transcripts) go once the change takes effect. `external`
	// sends them to the External Storage that `external_storage_id` names. `internal`
	// means no External Storage, so no more copies are made. When omitted, unchanged.
	// Cannot be null.
	TranscriptsStorage BetaStorageMemberParamUnion `json:"transcripts_storage,omitzero"`
	paramObj
}

func (r BetaInferenceDataPrivacyUpdateParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaInferenceDataPrivacyUpdateParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaInferenceDataPrivacyUpdateParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// No External Storage for this purpose.
type BetaInternalStorage struct {
	Type constant.Internal `json:"type" default:"internal"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaInternalStorage) RawJSON() string { return r.JSON.raw }
func (r *BetaInternalStorage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewBetaInternalStorageParam() BetaInternalStorageParam {
	return BetaInternalStorageParam{
		Type: "internal",
	}
}

// No External Storage for this purpose.
//
// This struct has a constant value, construct it with
// [NewBetaInternalStorageParam].
type BetaInternalStorageParam struct {
	Type constant.Internal `json:"type" default:"internal"`
	paramObj
}

func (r BetaInternalStorageParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaInternalStorageParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaInternalStorageParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaNoBillingWorkspaceRole string

const (
	BetaNoBillingWorkspaceRoleWorkspaceAdmin               BetaNoBillingWorkspaceRole = "workspace_admin"
	BetaNoBillingWorkspaceRoleWorkspaceDeveloper           BetaNoBillingWorkspaceRole = "workspace_developer"
	BetaNoBillingWorkspaceRoleWorkspaceRestrictedDeveloper BetaNoBillingWorkspaceRole = "workspace_restricted_developer"
	BetaNoBillingWorkspaceRoleWorkspaceUser                BetaNoBillingWorkspaceRole = "workspace_user"
)

// BetaStorageMemberUnion contains all possible properties and values from
// [BetaInternalStorage], [BetaExternalStorageMember].
//
// Use the [BetaStorageMemberUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaStorageMemberUnion struct {
	// Any of "internal", "external".
	Type string `json:"type"`
	// This field is from variant [BetaExternalStorageMember].
	ExternalStorageID string `json:"external_storage_id"`
	JSON              struct {
		Type              respjson.Field
		ExternalStorageID respjson.Field
		raw               string
	} `json:"-"`
}

// anyBetaStorageMember is implemented by each variant of [BetaStorageMemberUnion]
// to add type safety for the return type of [BetaStorageMemberUnion.AsAny]
type anyBetaStorageMember interface {
	implBetaStorageMemberUnion()
}

func (BetaInternalStorage) implBetaStorageMemberUnion()       {}
func (BetaExternalStorageMember) implBetaStorageMemberUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaStorageMemberUnion.AsAny().(type) {
//	case anthropic.BetaInternalStorage:
//	case anthropic.BetaExternalStorageMember:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaStorageMemberUnion) AsAny() anyBetaStorageMember {
	switch u.Type {
	case "internal":
		return u.AsInternal()
	case "external":
		return u.AsExternal()
	}
	return nil
}

func (u BetaStorageMemberUnion) AsInternal() (v BetaInternalStorage) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaStorageMemberUnion) AsExternal() (v BetaExternalStorageMember) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaStorageMemberUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaStorageMemberUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func BetaStorageMemberParamOfExternal(externalStorageID string) BetaStorageMemberParamUnion {
	var external BetaExternalStorageMemberParam
	external.ExternalStorageID = externalStorageID
	return BetaStorageMemberParamUnion{OfExternal: &external}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaStorageMemberParamUnion struct {
	OfInternal *BetaInternalStorageParam       `json:",omitzero,inline"`
	OfExternal *BetaExternalStorageMemberParam `json:",omitzero,inline"`
	paramUnion
}

func (u BetaStorageMemberParamUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfInternal, u.OfExternal)
}
func (u *BetaStorageMemberParamUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaStorageMemberParamUnion) asAny() any {
	if !param.IsOmitted(u.OfInternal) {
		return u.OfInternal
	} else if !param.IsOmitted(u.OfExternal) {
		return u.OfExternal
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaStorageMemberParamUnion) GetExternalStorageID() *string {
	if vt := u.OfExternal; vt != nil {
		return &vt.ExternalStorageID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaStorageMemberParamUnion) GetType() *string {
	if vt := u.OfInternal; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	} else if vt := u.OfExternal; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaStorageMemberParamUnion](
		"type",
		apijson.Discriminator[BetaInternalStorageParam]("internal"),
		apijson.Discriminator[BetaExternalStorageMemberParam]("external"),
	)
}

type BetaWorkspace struct {
	// ID of the Workspace.
	ID string `json:"id" api:"required"`
	// RFC 3339 datetime string indicating when the Workspace was archived, or `null`
	// if the Workspace is not archived.
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// Identifier for this Workspace's encryption compartment. When you configure an
	// AWS KMS customer-managed encryption key (CMEK) that isn't bound to a user
	// profile, reference this value in your KMS key-policy condition so the key is
	// scoped to this compartment. On GCP and Azure, Anthropic enforces the compartment
	// binding automatically; you do not need to reference this value in your key
	// configuration. See the CMEK integration guide for the required key
	// configuration; unless your organization is on Claude Platform on AWS, it
	// includes a separate value used during key validation. On Claude Platform on AWS
	// there is no separate validation value: your organization's own key is validated
	// against this Workspace's own value when it is attached, so if its key policy
	// uses the compartment condition, add this value to it before attaching the key.
	CompartmentID string `json:"compartment_id" api:"required"`
	// RFC 3339 datetime string indicating when the Workspace was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Data residency configuration.
	DataResidency BetaDataResidency `json:"data_residency" api:"required"`
	// Hex color code representing the Workspace in the Anthropic Console.
	DisplayColor string `json:"display_color" api:"required"`
	// ID of the customer-managed encryption key (CMEK) configuration to use for this
	// Workspace. Setting this field requires CMEK to be enabled for your organization.
	// When set, data stored for this Workspace is encrypted with the referenced key.
	// Create key configurations with the External Keys API. On Claude Platform on AWS
	// the value for an AWS KMS key is the key's ARN, and your organization's own key
	// must be a single-Region key in the same Region as the Workspace. If that key is
	// outside the Workspace's AWS account, it can be attached only by updating an
	// active Workspace. On that platform the key is validated against this Workspace
	// when it is attached, so a key-policy problem is reported as an error on this
	// request. This field is write-once: once a key is attached to a Workspace it
	// cannot be detached or replaced. To rotate key material, rotate the underlying
	// key on your cloud KMS; the `external_key_id` stays the same.
	ExternalKeyID string `json:"external_key_id" api:"required"`
	// Name of the Workspace.
	Name string `json:"name" api:"required"`
	// User-defined tags as string key-value pairs. Keys may not begin with
	// `anthropic`.
	Tags map[string]string `json:"tags" api:"required"`
	// Object type.
	//
	// For Workspaces, this is always `"workspace"`.
	Type constant.Workspace `json:"type" default:"workspace"`
	// ID of the user profile attached to the Workspace, or null if it has none. Every
	// inference request through a Workspace that has a user profile must name that
	// profile.
	UserProfileID string `json:"user_profile_id" api:"required"`
	// Where the Workspace's retained inference data and transcript copies go, and
	// whether Anthropic's human review applies. Returned only to requests that send
	// the `efs-2026-10-01` `anthropic-beta` value, and each member only to
	// organizations enabled for it.
	InferenceDataPrivacy BetaInferenceDataPrivacy `json:"inference_data_privacy"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                   respjson.Field
		ArchivedAt           respjson.Field
		CompartmentID        respjson.Field
		CreatedAt            respjson.Field
		DataResidency        respjson.Field
		DisplayColor         respjson.Field
		ExternalKeyID        respjson.Field
		Name                 respjson.Field
		Tags                 respjson.Field
		Type                 respjson.Field
		UserProfileID        respjson.Field
		InferenceDataPrivacy respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWorkspace) RawJSON() string { return r.JSON.raw }
func (r *BetaWorkspace) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWorkspaceMember struct {
	// Object type.
	//
	// For Workspace Members, this is always `"workspace_member"`.
	Type constant.WorkspaceMember `json:"type" default:"workspace_member"`
	// ID of the User.
	UserID string `json:"user_id" api:"required"`
	// ID of the Workspace.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// Role of the Workspace Member.
	//
	// Any of "workspace_admin", "workspace_billing", "workspace_developer",
	// "workspace_restricted_developer", "workspace_user".
	WorkspaceRole BetaWorkspaceRole `json:"workspace_role" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type          respjson.Field
		UserID        respjson.Field
		WorkspaceID   respjson.Field
		WorkspaceRole respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaWorkspaceMember) RawJSON() string { return r.JSON.raw }
func (r *BetaWorkspaceMember) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaWorkspaceRole string

const (
	BetaWorkspaceRoleWorkspaceAdmin               BetaWorkspaceRole = "workspace_admin"
	BetaWorkspaceRoleWorkspaceBilling             BetaWorkspaceRole = "workspace_billing"
	BetaWorkspaceRoleWorkspaceDeveloper           BetaWorkspaceRole = "workspace_developer"
	BetaWorkspaceRoleWorkspaceRestrictedDeveloper BetaWorkspaceRole = "workspace_restricted_developer"
	BetaWorkspaceRoleWorkspaceUser                BetaWorkspaceRole = "workspace_user"
)

type BetaOrganizationWorkspaceNewParams struct {
	// Name of the Workspace.
	Name string `json:"name" api:"required"`
	// Hex color code representing the Workspace in the Anthropic Console.
	DisplayColor param.Opt[string] `json:"display_color,omitzero"`
	// ID of the customer-managed encryption key (CMEK) configuration to use for this
	// Workspace. Setting this field requires CMEK to be enabled for your organization.
	// When set, data stored for this Workspace is encrypted with the referenced key.
	// Create key configurations with the External Keys API. On Claude Platform on AWS
	// the value for an AWS KMS key is the key's ARN, and your organization's own key
	// must be a single-Region key in the same Region as the Workspace. If that key is
	// outside the Workspace's AWS account, it can be attached only by updating an
	// active Workspace. On that platform the key is validated against this Workspace
	// when it is attached, so a key-policy problem is reported as an error on this
	// request. This field is write-once: once a key is attached to a Workspace it
	// cannot be detached or replaced. To rotate key material, rotate the underlying
	// key on your cloud KMS; the `external_key_id` stays the same.
	ExternalKeyID param.Opt[string] `json:"external_key_id,omitzero"`
	// ID of a user profile to attach to the Workspace. Every inference request through
	// the Workspace must then name this user profile; the API rejects a request that
	// names none or a different one. The user profile can't be replaced or removed
	// later; only deleting the user profile clears it.
	UserProfileID param.Opt[string] `json:"user_profile_id,omitzero"`
	// User-defined tags as string key-value pairs. Keys may not begin with
	// `anthropic`.
	Tags map[string]string `json:"tags,omitzero"`
	// Data residency configuration for the workspace. If omitted, defaults to
	// `workspace_geo: "us"`, `allowed_inference_geos: "unrestricted"`, and
	// `default_inference_geo: "global"`.
	DataResidency BetaDataResidencyCreateConfigParam `json:"data_residency,omitzero"`
	// Where the Workspace's retained inference data and transcript copies go, and
	// whether Anthropic's human review applies. Requires the `efs-2026-10-01`
	// `anthropic-beta` value. When omitted, `safeguards_storage` and
	// `transcripts_storage` are `internal` and `human_review` is `enabled`. Cannot be
	// null. An Admin API key can't set this field, because it doesn't belong to a user
	// or a service account. Use an API key that belongs to one, or an OAuth token.
	InferenceDataPrivacy BetaInferenceDataPrivacyCreateParam `json:"inference_data_privacy,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationWorkspaceNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationWorkspaceNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationWorkspaceNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationWorkspaceGetParams struct {
	// Optional. Send `efs-2026-10-01` in this header to use this endpoint's beta
	// features: additional parameters, accepted values and response fields.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationWorkspaceUpdateParams struct {
	// Hex color code representing the Workspace in the Anthropic Console.
	DisplayColor param.Opt[string] `json:"display_color,omitzero"`
	// ID of the customer-managed encryption key (CMEK) configuration to use for this
	// Workspace. Setting this field requires CMEK to be enabled for your organization.
	// When set, data stored for this Workspace is encrypted with the referenced key.
	// Create key configurations with the External Keys API. On Claude Platform on AWS
	// the value for an AWS KMS key is the key's ARN, and your organization's own key
	// must be a single-Region key in the same Region as the Workspace. If that key is
	// outside the Workspace's AWS account, it can be attached only by updating an
	// active Workspace. On that platform the key is validated against this Workspace
	// when it is attached, so a key-policy problem is reported as an error on this
	// request. This field is write-once: once a key is attached to a Workspace it
	// cannot be detached or replaced. To rotate key material, rotate the underlying
	// key on your cloud KMS; the `external_key_id` stays the same.
	ExternalKeyID param.Opt[string] `json:"external_key_id,omitzero"`
	// Name of the Workspace.
	Name param.Opt[string] `json:"name,omitzero"`
	// ID of a user profile to attach to the Workspace. Every inference request through
	// the Workspace must then name this user profile; the API rejects a request that
	// names none or a different one. Once set, it can't be replaced or removed
	// (sending the same ID again changes nothing); only deleting the user profile
	// clears it. It can't be set on a Workspace whose External Key isn't bound to this
	// user profile. The default Workspace can't have a user profile.
	UserProfileID param.Opt[string] `json:"user_profile_id,omitzero"`
	// User-defined tags as string key-value pairs. Keys may not begin with
	// `anthropic`.
	Tags map[string]string `json:"tags,omitzero"`
	// Data residency configuration for the workspace.
	DataResidency BetaDataResidencyUpdateConfigParam `json:"data_residency,omitzero"`
	// Changes to where the Workspace's retained inference data and transcript copies
	// go, and to whether Anthropic's human review applies. Requires the
	// `efs-2026-10-01` `anthropic-beta` value. An omitted member is unchanged. Cannot
	// be null. An Admin API key can't set this field, because it doesn't belong to a
	// user or a service account. Use an API key that belongs to one, or an OAuth
	// token.
	InferenceDataPrivacy BetaInferenceDataPrivacyUpdateParam `json:"inference_data_privacy,omitzero"`
	// Optional. Send `efs-2026-10-01` in this header to use this endpoint's beta
	// features: additional parameters, accepted values and response fields.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationWorkspaceUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationWorkspaceUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationWorkspaceUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationWorkspaceListParams struct {
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately after this object.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately before this object.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// Whether to include Workspaces that have been archived in the response
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Whether to include the organization's default Workspace in the response
	IncludeDefault param.Opt[bool] `query:"include_default,omitzero" json:"-"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `1000`.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Only the Workspaces with these IDs, at most 100. `include_archived` and
	// `include_default` still apply, and IDs of Workspaces the caller can't list are
	// left out.
	IDs []string `query:"ids,omitzero" json:"-"`
	// Optional. Send `efs-2026-10-01` in this header to use this endpoint's beta
	// features: additional parameters, accepted values and response fields.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationWorkspaceListParams]'s query parameters as
// `url.Values`.
func (r BetaOrganizationWorkspaceListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationWorkspaceArchiveParams struct {
	// Optional. Send `efs-2026-10-01` in this header to use this endpoint's beta
	// features: additional parameters, accepted values and response fields.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
