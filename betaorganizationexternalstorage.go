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

// BetaOrganizationExternalStorageService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationExternalStorageService] method instead.
type BetaOrganizationExternalStorageService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationExternalStorageService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationExternalStorageService(opts ...option.RequestOption) (r BetaOrganizationExternalStorageService) {
	r = BetaOrganizationExternalStorageService{}
	r.Options = opts
	return
}

// Register a bucket in the organization's own cloud account. It is saved as
// `pending_configuration` and verified when activated.
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeBetaGrantRequired] (HTTP 400)
//   - [ErrorCodeBetaOrganizationNotEligible] (HTTP 400)
//   - [ErrorCodeBetaProviderNotAllowed] (HTTP 400)
//   - [ErrorCodeBetaUnsupportedRegion] (HTTP 400)
//   - [ErrorCodeUserProfileNotFound] (HTTP 400)
//   - [ErrorCodeSandboxEgressRestricted] (HTTP 403)
//   - [ErrorCodeBetaStorageAlreadyRegistered] (HTTP 409)
//   - [ErrorCodeBetaAuthorizationCheckUnavailable] (HTTP 503)
func (r *BetaOrganizationExternalStorageService) New(ctx context.Context, params BetaOrganizationExternalStorageNewParams, opts ...option.RequestOption) (res *BetaExternalStorage, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "efs-2026-10-01")}, opts...)
	path := "v1/organizations/external_storage"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.external_storage.create", "/v1/organizations/external_storage"), requestconfig.WithSpanCreatedID("anthropic.organization.external_storage.id", &res, func(created *BetaExternalStorage) string { return created.ID }))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get one External Storage in the organization.
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeBetaInvalidStorageARN] (HTTP 400)
func (r *BetaOrganizationExternalStorageService) Get(ctx context.Context, externalStorageID string, query BetaOrganizationExternalStorageGetParams, opts ...option.RequestOption) (res *BetaExternalStorage, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "efs-2026-10-01")}, opts...)
	if externalStorageID == "" {
		err = errors.New("missing required external_storage_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/external_storage/%s", requestconfig.EscapePathSegment(externalStorageID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.external_storage.retrieve", "/v1/organizations/external_storage/{external_storage_id}"), requestconfig.WithSpanAttribute("anthropic.organization.external_storage.id", externalStorageID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Rename External Storage or clear its name, in any status.
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeBetaInvalidStorageARN] (HTTP 400)
//   - [ErrorCodeSandboxEgressRestricted] (HTTP 403)
func (r *BetaOrganizationExternalStorageService) Update(ctx context.Context, externalStorageID string, params BetaOrganizationExternalStorageUpdateParams, opts ...option.RequestOption) (res *BetaExternalStorage, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "efs-2026-10-01")}, opts...)
	if externalStorageID == "" {
		err = errors.New("missing required external_storage_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/external_storage/%s", requestconfig.EscapePathSegment(externalStorageID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.external_storage.update", "/v1/organizations/external_storage/{external_storage_id}"), requestconfig.WithSpanAttribute("anthropic.organization.external_storage.id", externalStorageID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List the organization's External Storage, newest first.
func (r *BetaOrganizationExternalStorageService) List(ctx context.Context, params BetaOrganizationExternalStorageListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaExternalStorage], err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "efs-2026-10-01"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/external_storage"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.external_storage.list", "/v1/organizations/external_storage", "limit", "page", "user_profile_ids"))
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

// List the organization's External Storage, newest first.
func (r *BetaOrganizationExternalStorageService) ListAutoPaging(ctx context.Context, params BetaOrganizationExternalStorageListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaExternalStorage] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Delete External Storage. Storage that is pending configuration, available, or
// active but never attached is deleted at once, unless an Enterprise Frontier
// Safeguards workspace policy lists it. Other storage gets a 400 that says what to
// do next. On `draining` storage that 400 can mean this delete started a removal
// that can't be undone. To check, get the storage: if `status.removal_started` is
// true, delete it again later to finish.
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeBetaExternalStorageInUse] (HTTP 400)
//   - [ErrorCodeBetaExternalStorageInWorkspacePolicy] (HTTP 400)
//   - [ErrorCodeBetaInvalidStorageARN] (HTTP 400)
//   - [ErrorCodeSandboxEgressRestricted] (HTTP 403)
func (r *BetaOrganizationExternalStorageService) Delete(ctx context.Context, externalStorageID string, body BetaOrganizationExternalStorageDeleteParams, opts ...option.RequestOption) (res *BetaDeletedExternalStorage, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "efs-2026-10-01")}, opts...)
	if externalStorageID == "" {
		err = errors.New("missing required external_storage_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/external_storage/%s", requestconfig.EscapePathSegment(externalStorageID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.external_storage.delete", "/v1/organizations/external_storage/{external_storage_id}"), requestconfig.WithSpanAttribute("anthropic.organization.external_storage.id", externalStorageID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Verify pending External Storage and make it available. Storage that is already
// available, active or draining comes back unchanged, without any check.
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeBetaGrantRequired] (HTTP 400)
//   - [ErrorCodeBetaInvalidStorageARN] (HTTP 400)
//   - [ErrorCodeBetaOrganizationNotEligible] (HTTP 400)
//   - [ErrorCodeBetaProviderNotAllowed] (HTTP 400)
//   - [ErrorCodeBetaStorageVerificationFailed] (HTTP 400)
//   - [ErrorCodeUserProfileNotFound] (HTTP 400)
//   - [ErrorCodeSandboxEgressRestricted] (HTTP 403)
//   - [ErrorCodeBetaAuthorizationCheckUnavailable] (HTTP 503)
func (r *BetaOrganizationExternalStorageService) Activate(ctx context.Context, externalStorageID string, body BetaOrganizationExternalStorageActivateParams, opts ...option.RequestOption) (res *BetaExternalStorage, err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "efs-2026-10-01")}, opts...)
	if externalStorageID == "" {
		err = errors.New("missing required external_storage_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/external_storage/%s/activate", requestconfig.EscapePathSegment(externalStorageID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.external_storage.activate", "/v1/organizations/external_storage/{external_storage_id}/activate"), requestconfig.WithSpanAttribute("anthropic.organization.external_storage.id", externalStorageID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// At least one Workspace has the storage attached.
type BetaActiveExternalStorageStatus struct {
	// Whether the bucket passed its latest recorded check. It can stay `true` after
	// the bucket stops working, so don't use it for monitoring.
	Healthy bool `json:"healthy" api:"required"`
	// When the bucket's latest recorded check ran, in RFC 3339.
	LastCheckedAt time.Time       `json:"last_checked_at" api:"required" format:"date-time"`
	Type          constant.Active `json:"type" default:"active"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Healthy       respjson.Field
		LastCheckedAt respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaActiveExternalStorageStatus) RawJSON() string { return r.JSON.raw }
func (r *BetaActiveExternalStorageStatus) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The storage is verified, no Workspace has it attached, and Anthropic doesn't use
// the bucket.
type BetaAvailableExternalStorageStatus struct {
	// Whether the bucket passed its latest recorded check. It can stay `true` after
	// the bucket stops working, so don't use it for monitoring.
	Healthy bool `json:"healthy" api:"required"`
	// When the bucket's latest recorded check ran, in RFC 3339.
	LastCheckedAt time.Time          `json:"last_checked_at" api:"required" format:"date-time"`
	Type          constant.Available `json:"type" default:"available"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Healthy       respjson.Field
		LastCheckedAt respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAvailableExternalStorageStatus) RawJSON() string { return r.JSON.raw }
func (r *BetaAvailableExternalStorageStatus) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An Amazon S3 bucket.
type BetaAWSExternalStorageConfig struct {
	// 12-digit ID of the AWS account that owns the bucket.
	AccountID string `json:"account_id" api:"required"`
	// The bucket's name, without `s3://`.
	BucketName string `json:"bucket_name" api:"required"`
	// The bucket's AWS region.
	Region string       `json:"region" api:"required"`
	Type   constant.AWS `json:"type" default:"aws"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AccountID   respjson.Field
		BucketName  respjson.Field
		Region      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAWSExternalStorageConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaAWSExternalStorageConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An Amazon S3 bucket.
//
// The properties AccountID, BucketName, Region, Type are required.
type BetaAWSExternalStorageConfigParam struct {
	// 12-digit ID of the AWS account that owns the bucket.
	AccountID string `json:"account_id" api:"required"`
	// The bucket's name, without `s3://`: 3 to 63 lower-case letters, digits and
	// hyphens, starting and ending with a letter or digit.
	BucketName string `json:"bucket_name" api:"required"`
	// The bucket's AWS region. Activating the storage fails if Anthropic can't use
	// buckets in that region.
	Region string `json:"region" api:"required"`
	// This field can be elided, and will marshal its zero value as "aws".
	Type constant.AWS `json:"type" default:"aws"`
	paramObj
}

func (r BetaAWSExternalStorageConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaAWSExternalStorageConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAWSExternalStorageConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaDeletedExternalStorage struct {
	// The External Storage ID, which starts with `estrg_`. Through Claude Platform on
	// AWS, storage whose `provider_config.type` is `aws` has its bucket ARN as its ID
	// instead: `arn:<partition>:s3:::<bucket_name>`, where `<partition>` is the
	// partition of the bucket's region.
	ID   string                          `json:"id" api:"required"`
	Type constant.ExternalStorageDeleted `json:"type" default:"external_storage_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDeletedExternalStorage) RawJSON() string { return r.JSON.raw }
func (r *BetaDeletedExternalStorage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// No Workspace has the storage attached. A Workspace can attach it again, unless
// `removal_started` is true.
type BetaDrainingExternalStorageStatus struct {
	// Whether the bucket passed its latest recorded check. It can stay `true` after
	// the bucket stops working, so don't use it for monitoring.
	Healthy bool `json:"healthy" api:"required"`
	// When the bucket's latest recorded check ran, in RFC 3339.
	LastCheckedAt time.Time `json:"last_checked_at" api:"required" format:"date-time"`
	// True once a delete has started removing this storage, which can't be undone. It
	// then can't be attached, and another delete finishes the removal.
	RemovalStarted bool              `json:"removal_started" api:"required"`
	Type           constant.Draining `json:"type" default:"draining"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Healthy        respjson.Field
		LastCheckedAt  respjson.Field
		RemovalStarted respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDrainingExternalStorageStatus) RawJSON() string { return r.JSON.raw }
func (r *BetaDrainingExternalStorageStatus) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A storage bucket in the organization's own cloud account. Anthropic writes
// inference data to it from each Workspace it is attached to.
type BetaExternalStorage struct {
	// The External Storage ID, which starts with `estrg_`. Through Claude Platform on
	// AWS, storage whose `provider_config.type` is `aws` has its bucket ARN as its ID
	// instead: `arn:<partition>:s3:::<bucket_name>`, where `<partition>` is the
	// partition of the bucket's region.
	ID        string    `json:"id" api:"required"`
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Human-friendly display name. Null if none was set.
	DisplayName    string                                 `json:"display_name" api:"required"`
	ProviderConfig BetaExternalStorageProviderConfigUnion `json:"provider_config" api:"required"`
	// The version of Anthropic's setup template whose IAM role names this storage was
	// found to use at activation. Null if no version was recorded, as for storage not
	// yet activated or not on S3.
	SetupTemplateVersion string                         `json:"setup_template_version" api:"required"`
	Status               BetaExternalStorageStatusUnion `json:"status" api:"required"`
	Type                 constant.ExternalStorage       `json:"type" default:"external_storage"`
	UpdatedAt            time.Time                      `json:"updated_at" api:"required" format:"date-time"`
	// The user profile the storage is bound to, or null.
	UserProfileID string `json:"user_profile_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                   respjson.Field
		CreatedAt            respjson.Field
		DisplayName          respjson.Field
		ProviderConfig       respjson.Field
		SetupTemplateVersion respjson.Field
		Status               respjson.Field
		Type                 respjson.Field
		UpdatedAt            respjson.Field
		UserProfileID        respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaExternalStorage) RawJSON() string { return r.JSON.raw }
func (r *BetaExternalStorage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaExternalStorageProviderConfigUnion contains all possible properties and
// values from [BetaAWSExternalStorageConfig], [BetaGCPExternalStorageConfig].
//
// Use the [BetaExternalStorageProviderConfigUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaExternalStorageProviderConfigUnion struct {
	// This field is from variant [BetaAWSExternalStorageConfig].
	AccountID  string `json:"account_id"`
	BucketName string `json:"bucket_name"`
	Region     string `json:"region"`
	// Any of "aws", "gcp".
	Type string `json:"type"`
	// This field is from variant [BetaGCPExternalStorageConfig].
	ProjectNumber string `json:"project_number"`
	JSON          struct {
		AccountID     respjson.Field
		BucketName    respjson.Field
		Region        respjson.Field
		Type          respjson.Field
		ProjectNumber respjson.Field
		raw           string
	} `json:"-"`
}

// anyBetaExternalStorageProviderConfig is implemented by each variant of
// [BetaExternalStorageProviderConfigUnion] to add type safety for the return type
// of [BetaExternalStorageProviderConfigUnion.AsAny]
type anyBetaExternalStorageProviderConfig interface {
	implBetaExternalStorageProviderConfigUnion()
}

func (BetaAWSExternalStorageConfig) implBetaExternalStorageProviderConfigUnion() {}
func (BetaGCPExternalStorageConfig) implBetaExternalStorageProviderConfigUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaExternalStorageProviderConfigUnion.AsAny().(type) {
//	case anthropic.BetaAWSExternalStorageConfig:
//	case anthropic.BetaGCPExternalStorageConfig:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaExternalStorageProviderConfigUnion) AsAny() anyBetaExternalStorageProviderConfig {
	switch u.Type {
	case "aws":
		return u.AsAWS()
	case "gcp":
		return u.AsGCP()
	}
	return nil
}

func (u BetaExternalStorageProviderConfigUnion) AsAWS() (v BetaAWSExternalStorageConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaExternalStorageProviderConfigUnion) AsGCP() (v BetaGCPExternalStorageConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaExternalStorageProviderConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaExternalStorageProviderConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaExternalStorageStatusUnion contains all possible properties and values from
// [BetaPendingConfigurationExternalStorageStatus],
// [BetaAvailableExternalStorageStatus], [BetaActiveExternalStorageStatus],
// [BetaDrainingExternalStorageStatus].
//
// Use the [BetaExternalStorageStatusUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaExternalStorageStatusUnion struct {
	// Any of "pending_configuration", "available", "active", "draining".
	Type          string    `json:"type"`
	Healthy       bool      `json:"healthy"`
	LastCheckedAt time.Time `json:"last_checked_at"`
	// This field is from variant [BetaDrainingExternalStorageStatus].
	RemovalStarted bool `json:"removal_started"`
	JSON           struct {
		Type           respjson.Field
		Healthy        respjson.Field
		LastCheckedAt  respjson.Field
		RemovalStarted respjson.Field
		raw            string
	} `json:"-"`
}

// anyBetaExternalStorageStatus is implemented by each variant of
// [BetaExternalStorageStatusUnion] to add type safety for the return type of
// [BetaExternalStorageStatusUnion.AsAny]
type anyBetaExternalStorageStatus interface {
	implBetaExternalStorageStatusUnion()
}

func (BetaPendingConfigurationExternalStorageStatus) implBetaExternalStorageStatusUnion() {}
func (BetaAvailableExternalStorageStatus) implBetaExternalStorageStatusUnion()            {}
func (BetaActiveExternalStorageStatus) implBetaExternalStorageStatusUnion()               {}
func (BetaDrainingExternalStorageStatus) implBetaExternalStorageStatusUnion()             {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaExternalStorageStatusUnion.AsAny().(type) {
//	case anthropic.BetaPendingConfigurationExternalStorageStatus:
//	case anthropic.BetaAvailableExternalStorageStatus:
//	case anthropic.BetaActiveExternalStorageStatus:
//	case anthropic.BetaDrainingExternalStorageStatus:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaExternalStorageStatusUnion) AsAny() anyBetaExternalStorageStatus {
	switch u.Type {
	case "pending_configuration":
		return u.AsPendingConfiguration()
	case "available":
		return u.AsAvailable()
	case "active":
		return u.AsActive()
	case "draining":
		return u.AsDraining()
	}
	return nil
}

func (u BetaExternalStorageStatusUnion) AsPendingConfiguration() (v BetaPendingConfigurationExternalStorageStatus) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaExternalStorageStatusUnion) AsAvailable() (v BetaAvailableExternalStorageStatus) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaExternalStorageStatusUnion) AsActive() (v BetaActiveExternalStorageStatus) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaExternalStorageStatusUnion) AsDraining() (v BetaDrainingExternalStorageStatus) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaExternalStorageStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaExternalStorageStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Google Cloud Storage bucket.
type BetaGCPExternalStorageConfig struct {
	// The bucket's name, without `gs://`.
	BucketName string `json:"bucket_name" api:"required"`
	// Number of the Google Cloud project that holds the bucket.
	ProjectNumber string `json:"project_number" api:"required"`
	// The bucket's Google Cloud region.
	Region string       `json:"region" api:"required"`
	Type   constant.GCP `json:"type" default:"gcp"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		BucketName    respjson.Field
		ProjectNumber respjson.Field
		Region        respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaGCPExternalStorageConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaGCPExternalStorageConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Google Cloud Storage bucket.
//
// The properties BucketName, ProjectNumber, Region, Type are required.
type BetaGCPExternalStorageConfigParam struct {
	// The bucket's name, without `gs://`: 3 to 63 lower-case letters, digits and
	// hyphens, starting and ending with a letter or digit. It can't start with `goog`
	// or contain `google`.
	BucketName string `json:"bucket_name" api:"required"`
	// Project number, not project ID, of the Google Cloud project that holds the
	// bucket.
	ProjectNumber string `json:"project_number" api:"required"`
	// The bucket's Google Cloud region, such as `us-east5`: one of the single regions
	// that the setup guide lists. Dual-regions and multi-regions aren't supported.
	Region string `json:"region" api:"required"`
	// This field can be elided, and will marshal its zero value as "gcp".
	Type constant.GCP `json:"type" default:"gcp"`
	paramObj
}

func (r BetaGCPExternalStorageConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaGCPExternalStorageConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaGCPExternalStorageConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The storage isn't verified yet.
type BetaPendingConfigurationExternalStorageStatus struct {
	Type constant.PendingConfiguration `json:"type" default:"pending_configuration"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPendingConfigurationExternalStorageStatus) RawJSON() string { return r.JSON.raw }
func (r *BetaPendingConfigurationExternalStorageStatus) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationExternalStorageNewParams struct {
	ProviderConfig BetaOrganizationExternalStorageNewParamsProviderConfigUnion `json:"provider_config,omitzero" api:"required"`
	// Human-friendly display name.
	DisplayName param.Opt[string] `json:"display_name,omitzero"`
	// The user profile to bind the storage to. It can't be changed later. The storage
	// can be attached only to Workspaces with this profile, and the bucket's setup
	// must name the same ID.
	UserProfileID param.Opt[string] `json:"user_profile_id,omitzero"`
	// This endpoint is in beta: requests must send `efs-2026-10-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationExternalStorageNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationExternalStorageNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationExternalStorageNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaOrganizationExternalStorageNewParamsProviderConfigUnion struct {
	OfAWS *BetaAWSExternalStorageConfigParam `json:",omitzero,inline"`
	OfGCP *BetaGCPExternalStorageConfigParam `json:",omitzero,inline"`
	paramUnion
}

func (u BetaOrganizationExternalStorageNewParamsProviderConfigUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAWS, u.OfGCP)
}
func (u *BetaOrganizationExternalStorageNewParamsProviderConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaOrganizationExternalStorageNewParamsProviderConfigUnion) asAny() any {
	if !param.IsOmitted(u.OfAWS) {
		return u.OfAWS
	} else if !param.IsOmitted(u.OfGCP) {
		return u.OfGCP
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalStorageNewParamsProviderConfigUnion) GetAccountID() *string {
	if vt := u.OfAWS; vt != nil {
		return &vt.AccountID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalStorageNewParamsProviderConfigUnion) GetProjectNumber() *string {
	if vt := u.OfGCP; vt != nil {
		return &vt.ProjectNumber
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalStorageNewParamsProviderConfigUnion) GetBucketName() *string {
	if vt := u.OfAWS; vt != nil {
		return (*string)(&vt.BucketName)
	} else if vt := u.OfGCP; vt != nil {
		return (*string)(&vt.BucketName)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalStorageNewParamsProviderConfigUnion) GetRegion() *string {
	if vt := u.OfAWS; vt != nil {
		return (*string)(&vt.Region)
	} else if vt := u.OfGCP; vt != nil {
		return (*string)(&vt.Region)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalStorageNewParamsProviderConfigUnion) GetType() *string {
	if vt := u.OfAWS; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	} else if vt := u.OfGCP; vt != nil {
		return paramutil.AddrOrDefault(vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaOrganizationExternalStorageNewParamsProviderConfigUnion](
		"type",
		apijson.Discriminator[BetaAWSExternalStorageConfigParam]("aws"),
		apijson.Discriminator[BetaGCPExternalStorageConfigParam]("gcp"),
	)
}

type BetaOrganizationExternalStorageGetParams struct {
	// This endpoint is in beta: requests must send `efs-2026-10-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationExternalStorageUpdateParams struct {
	// Human-friendly display name; null clears it.
	DisplayName param.Opt[string] `json:"display_name,omitzero" api:"required"`
	// This endpoint is in beta: requests must send `efs-2026-10-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationExternalStorageUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationExternalStorageUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationExternalStorageUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationExternalStorageListParams struct {
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of results per page.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Returns only External Storage bound to one of these user profile IDs
	// (`uprof_...`). An empty entry matches storage bound to no user profile. At most
	// 10 entries, counting repeats.
	UserProfileIDs []string `query:"user_profile_ids,omitzero" json:"-"`
	// This endpoint is in beta: requests must send `efs-2026-10-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationExternalStorageListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationExternalStorageListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationExternalStorageDeleteParams struct {
	// This endpoint is in beta: requests must send `efs-2026-10-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationExternalStorageActivateParams struct {
	// This endpoint is in beta: requests must send `efs-2026-10-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
