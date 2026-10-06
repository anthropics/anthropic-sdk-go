package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
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

// OrganizationExternalKeyService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationExternalKeyService] method instead.
type OrganizationExternalKeyService struct {
	Options []option.RequestOption
}

// NewOrganizationExternalKeyService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewOrganizationExternalKeyService(opts ...option.RequestOption) (r OrganizationExternalKeyService) {
	r = OrganizationExternalKeyService{}
	r.Options = opts
	return
}

// Create an external key config owned by the caller's organization.
func (r *OrganizationExternalKeyService) New(ctx context.Context, body OrganizationExternalKeyNewParams, opts ...option.RequestOption) (res *ExternalKey, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/external_keys"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve a single external key config in the caller's organization by ID.
func (r *OrganizationExternalKeyService) Get(ctx context.Context, externalKeyID string, opts ...option.RequestOption) (res *ExternalKey, err error) {
	opts = slices.Concat(r.Options, opts)
	if externalKeyID == "" {
		err = errors.New("missing required external_key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/external_keys/%s", url.PathEscape(externalKeyID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Partially update an external key config. Omitted fields are left unchanged.
//
// `display_name` is always editable. `geo` and `provider_config` cannot be changed
// once any workspace references this config, because previously encrypted data
// requires the original key identity to decrypt.
func (r *OrganizationExternalKeyService) Update(ctx context.Context, externalKeyID string, body OrganizationExternalKeyUpdateParams, opts ...option.RequestOption) (res *ExternalKey, err error) {
	opts = slices.Concat(r.Options, opts)
	if externalKeyID == "" {
		err = errors.New("missing required external_key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/external_keys/%s", url.PathEscape(externalKeyID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// List external key configs in the caller's organization.
//
// Results are ordered by creation time (newest first). Use the `next_page` cursor
// from the response to fetch subsequent pages.
func (r *OrganizationExternalKeyService) List(ctx context.Context, query OrganizationExternalKeyListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ExternalKey], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/external_keys"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
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

// List external key configs in the caller's organization.
//
// Results are ordered by creation time (newest first). Use the `next_page` cursor
// from the response to fetch subsequent pages.
func (r *OrganizationExternalKeyService) ListAutoPaging(ctx context.Context, query OrganizationExternalKeyListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ExternalKey] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

// Delete an external key config.
//
// The request is rejected if any workspace still references this config.
func (r *OrganizationExternalKeyService) Delete(ctx context.Context, externalKeyID string, opts ...option.RequestOption) (res *OrganizationExternalKeyDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if externalKeyID == "" {
		err = errors.New("missing required external_key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/external_keys/%s", url.PathEscape(externalKeyID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Validate an external key config against the customer's KMS.
//
// Anthropic performs an encrypt/decrypt roundtrip against the configured KMS key
// and waits up to 30 seconds for the result. The response status is `success` if
// the roundtrip succeeded, or `failure` with an error message if it failed or
// timed out.
func (r *OrganizationExternalKeyService) Validate(ctx context.Context, externalKeyID string, opts ...option.RequestOption) (res *OrganizationExternalKeyValidateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if externalKeyID == "" {
		err = errors.New("missing required external_key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/external_keys/%s/validate", url.PathEscape(externalKeyID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type AWSExternalKeyConfig struct {
	// Full ARN of the AWS KMS key. On Claude Platform on AWS the key must be a
	// single-Region key in your organization's own AWS account; cross-account keys,
	// multi-Region keys, and alias ARNs are rejected.
	KMSARN string       `json:"kms_arn" api:"required"`
	Type   constant.AWS `json:"type" default:"aws"`
	// AWS region. Derived from `kms_arn` if omitted.
	Region string `json:"region" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		KMSARN      respjson.Field
		Type        respjson.Field
		Region      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AWSExternalKeyConfig) RawJSON() string { return r.JSON.raw }
func (r *AWSExternalKeyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this AWSExternalKeyConfig to a AWSExternalKeyConfigParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// AWSExternalKeyConfigParam.Overrides()
func (r AWSExternalKeyConfig) ToParam() AWSExternalKeyConfigParam {
	return param.Override[AWSExternalKeyConfigParam](json.RawMessage(r.RawJSON()))
}

// The properties KMSARN, Type are required.
type AWSExternalKeyConfigParam struct {
	// Full ARN of the AWS KMS key. On Claude Platform on AWS the key must be a
	// single-Region key in your organization's own AWS account; cross-account keys,
	// multi-Region keys, and alias ARNs are rejected.
	KMSARN string `json:"kms_arn" api:"required"`
	// AWS region. Derived from `kms_arn` if omitted.
	Region param.Opt[string] `json:"region,omitzero"`
	// This field can be elided, and will marshal its zero value as "aws".
	Type constant.AWS `json:"type" default:"aws"`
	paramObj
}

func (r AWSExternalKeyConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow AWSExternalKeyConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AWSExternalKeyConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type AzureExternalKeyConfig struct {
	// Name of the key within the vault.
	KeyName string `json:"key_name" api:"required"`
	// Azure AD tenant ID.
	TenantID string         `json:"tenant_id" api:"required"`
	Type     constant.Azure `json:"type" default:"azure"`
	// Key Vault data-plane URI — `https://{vault-name}.vault.azure.net` or
	// `https://{hsm-name}.managedhsm.azure.net`.
	VaultURI string `json:"vault_uri" api:"required"`
	// Azure AD application (client) ID. Omit to use Anthropic's multitenant app.
	// Provide only if using a single-tenant app registration in the customer's
	// directory.
	ClientID string `json:"client_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		KeyName     respjson.Field
		TenantID    respjson.Field
		Type        respjson.Field
		VaultURI    respjson.Field
		ClientID    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AzureExternalKeyConfig) RawJSON() string { return r.JSON.raw }
func (r *AzureExternalKeyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Azure Key Vault provider configuration.
//
// The properties KeyName, TenantID, Type, VaultURI are required.
type AzureExternalKeyConfigParam struct {
	// Name of the key within the vault.
	KeyName string `json:"key_name" api:"required"`
	// Azure AD tenant ID.
	TenantID string `json:"tenant_id" api:"required"`
	// Key Vault data-plane URI — `https://{vault-name}.vault.azure.net` or
	// `https://{hsm-name}.managedhsm.azure.net`.
	VaultURI string `json:"vault_uri" api:"required"`
	// Azure AD application (client) ID. Omit to use Anthropic's multitenant app.
	// Provide only if using a single-tenant app registration in the customer's
	// directory.
	ClientID param.Opt[string] `json:"client_id,omitzero"`
	// This field can be elided, and will marshal its zero value as "azure".
	Type constant.Azure `json:"type" default:"azure"`
	paramObj
}

func (r AzureExternalKeyConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow AzureExternalKeyConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AzureExternalKeyConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// CMEK external key config belonging to the caller's organization.
//
// Configs are organization-scoped. Workspaces attach to a config; once any
// workspace references it, the provider fields become effectively immutable
// (existing encrypted data needs the config for decrypt).
type ExternalKey struct {
	// Identifier of the external key config. A tagged ID prefixed `ekey_`, or — for
	// organizations on the Claude Platform on AWS — the AWS KMS key ARN.
	ID string `json:"id" api:"required"`
	// Whether any workspace uses this config to encrypt its data — counting live and
	// archived workspaces (an archived workspace's data remains encrypted under the
	// config), excluding deleted ones. Only an attached config is used by the
	// encryption path; an `unattached` config is inert and can be deleted.
	Attachment ExternalKeyAttachmentUnion `json:"attachment" api:"required"`
	CreatedAt  time.Time                  `json:"created_at" api:"required" format:"date-time"`
	// Human-friendly display name. Null if none was set.
	DisplayName string `json:"display_name" api:"required"`
	// Data residency geo. Selects which regional validator handles this key's
	// encrypt/decrypt roundtrips.
	Geo string `json:"geo" api:"required"`
	// KMS provider identity and auth coordinates.
	ProviderConfig ExternalKeyProviderConfigUnion `json:"provider_config" api:"required"`
	Type           constant.ExternalKey           `json:"type" default:"external_key"`
	UpdatedAt      time.Time                      `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		Attachment     respjson.Field
		CreatedAt      respjson.Field
		DisplayName    respjson.Field
		Geo            respjson.Field
		ProviderConfig respjson.Field
		Type           respjson.Field
		UpdatedAt      respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ExternalKey) RawJSON() string { return r.JSON.raw }
func (r *ExternalKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ExternalKeyAttachmentUnion contains all possible properties and values from
// [ExternalKeyAttachedAttachment], [ExternalKeyUnattachedAttachment].
//
// Use the [ExternalKeyAttachmentUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ExternalKeyAttachmentUnion struct {
	// Any of "attached", "unattached".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyExternalKeyAttachment is implemented by each variant of
// [ExternalKeyAttachmentUnion] to add type safety for the return type of
// [ExternalKeyAttachmentUnion.AsAny]
type anyExternalKeyAttachment interface {
	implExternalKeyAttachmentUnion()
}

func (ExternalKeyAttachedAttachment) implExternalKeyAttachmentUnion()   {}
func (ExternalKeyUnattachedAttachment) implExternalKeyAttachmentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ExternalKeyAttachmentUnion.AsAny().(type) {
//	case anthropic.ExternalKeyAttachedAttachment:
//	case anthropic.ExternalKeyUnattachedAttachment:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ExternalKeyAttachmentUnion) AsAny() anyExternalKeyAttachment {
	switch u.Type {
	case "attached":
		return u.AsAttached()
	case "unattached":
		return u.AsUnattached()
	}
	return nil
}

func (u ExternalKeyAttachmentUnion) AsAttached() (v ExternalKeyAttachedAttachment) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ExternalKeyAttachmentUnion) AsUnattached() (v ExternalKeyUnattachedAttachment) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ExternalKeyAttachmentUnion) RawJSON() string { return u.JSON.raw }

func (r *ExternalKeyAttachmentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ExternalKeyProviderConfigUnion contains all possible properties and values from
// [AWSExternalKeyConfig], [GCPExternalKeyConfig], [AzureExternalKeyConfig].
//
// Use the [ExternalKeyProviderConfigUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ExternalKeyProviderConfigUnion struct {
	// This field is from variant [AWSExternalKeyConfig].
	KMSARN string `json:"kms_arn"`
	// Any of "aws", "gcp", "azure".
	Type string `json:"type"`
	// This field is from variant [AWSExternalKeyConfig].
	Region  string `json:"region"`
	KeyName string `json:"key_name"`
	// This field is from variant [AzureExternalKeyConfig].
	TenantID string `json:"tenant_id"`
	// This field is from variant [AzureExternalKeyConfig].
	VaultURI string `json:"vault_uri"`
	// This field is from variant [AzureExternalKeyConfig].
	ClientID string `json:"client_id"`
	JSON     struct {
		KMSARN   respjson.Field
		Type     respjson.Field
		Region   respjson.Field
		KeyName  respjson.Field
		TenantID respjson.Field
		VaultURI respjson.Field
		ClientID respjson.Field
		raw      string
	} `json:"-"`
}

// anyExternalKeyProviderConfig is implemented by each variant of
// [ExternalKeyProviderConfigUnion] to add type safety for the return type of
// [ExternalKeyProviderConfigUnion.AsAny]
type anyExternalKeyProviderConfig interface {
	implExternalKeyProviderConfigUnion()
}

func (AWSExternalKeyConfig) implExternalKeyProviderConfigUnion()   {}
func (GCPExternalKeyConfig) implExternalKeyProviderConfigUnion()   {}
func (AzureExternalKeyConfig) implExternalKeyProviderConfigUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ExternalKeyProviderConfigUnion.AsAny().(type) {
//	case anthropic.AWSExternalKeyConfig:
//	case anthropic.GCPExternalKeyConfig:
//	case anthropic.AzureExternalKeyConfig:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ExternalKeyProviderConfigUnion) AsAny() anyExternalKeyProviderConfig {
	switch u.Type {
	case "aws":
		return u.AsAWS()
	case "gcp":
		return u.AsGCP()
	case "azure":
		return u.AsAzure()
	}
	return nil
}

func (u ExternalKeyProviderConfigUnion) AsAWS() (v AWSExternalKeyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ExternalKeyProviderConfigUnion) AsGCP() (v GCPExternalKeyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ExternalKeyProviderConfigUnion) AsAzure() (v AzureExternalKeyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ExternalKeyProviderConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *ExternalKeyProviderConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ExternalKeyAttachedAttachment struct {
	Type constant.Attached `json:"type" default:"attached"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ExternalKeyAttachedAttachment) RawJSON() string { return r.JSON.raw }
func (r *ExternalKeyAttachedAttachment) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ExternalKeyUnattachedAttachment struct {
	Type constant.Unattached `json:"type" default:"unattached"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ExternalKeyUnattachedAttachment) RawJSON() string { return r.JSON.raw }
func (r *ExternalKeyUnattachedAttachment) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type GCPExternalKeyConfig struct {
	// Full resource name of the Cloud KMS key.
	KeyName string       `json:"key_name" api:"required"`
	Type    constant.GCP `json:"type" default:"gcp"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		KeyName     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r GCPExternalKeyConfig) RawJSON() string { return r.JSON.raw }
func (r *GCPExternalKeyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this GCPExternalKeyConfig to a GCPExternalKeyConfigParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// GCPExternalKeyConfigParam.Overrides()
func (r GCPExternalKeyConfig) ToParam() GCPExternalKeyConfigParam {
	return param.Override[GCPExternalKeyConfigParam](json.RawMessage(r.RawJSON()))
}

// The properties KeyName, Type are required.
type GCPExternalKeyConfigParam struct {
	// Full resource name of the Cloud KMS key.
	KeyName string `json:"key_name" api:"required"`
	// This field can be elided, and will marshal its zero value as "gcp".
	Type constant.GCP `json:"type" default:"gcp"`
	paramObj
}

func (r GCPExternalKeyConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow GCPExternalKeyConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *GCPExternalKeyConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationExternalKeyDeleteResponse struct {
	// ID of the deleted External Key.
	ID   string                      `json:"id" api:"required"`
	Type constant.ExternalKeyDeleted `json:"type" default:"external_key_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationExternalKeyDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *OrganizationExternalKeyDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Result of a validation roundtrip against the customer's KMS.
//
// HTTP 200 for both outcomes — the operation completed; `status` says whether the
// key works.
type OrganizationExternalKeyValidateResponse struct {
	// Error message when status is `failure`. Null otherwise.
	Error string `json:"error" api:"required"`
	// `success` — encrypt/decrypt roundtrip succeeded. `failure` — the roundtrip
	// failed or timed out; see `error`.
	//
	// Any of "failure", "success".
	Status OrganizationExternalKeyValidateResponseStatus `json:"status" api:"required"`
	Type   constant.ExternalKeyValidation                `json:"type" default:"external_key_validation"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Error       respjson.Field
		Status      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OrganizationExternalKeyValidateResponse) RawJSON() string { return r.JSON.raw }
func (r *OrganizationExternalKeyValidateResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// `success` — encrypt/decrypt roundtrip succeeded. `failure` — the roundtrip
// failed or timed out; see `error`.
type OrganizationExternalKeyValidateResponseStatus string

const (
	OrganizationExternalKeyValidateResponseStatusFailure OrganizationExternalKeyValidateResponseStatus = "failure"
	OrganizationExternalKeyValidateResponseStatusSuccess OrganizationExternalKeyValidateResponseStatus = "success"
)

type OrganizationExternalKeyNewParams struct {
	// KMS provider identity and auth coordinates.
	ProviderConfig OrganizationExternalKeyNewParamsProviderConfigUnion `json:"provider_config,omitzero" api:"required"`
	// Human-friendly display name.
	DisplayName param.Opt[string] `json:"display_name,omitzero"`
	// Data residency geo. Only `us` is supported.
	//
	// Any of "us".
	Geo OrganizationExternalKeyNewParamsGeo `json:"geo,omitzero"`
	paramObj
}

func (r OrganizationExternalKeyNewParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationExternalKeyNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationExternalKeyNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type OrganizationExternalKeyNewParamsProviderConfigUnion struct {
	OfAWS   *AWSExternalKeyConfigParam   `json:",omitzero,inline"`
	OfGCP   *GCPExternalKeyConfigParam   `json:",omitzero,inline"`
	OfAzure *AzureExternalKeyConfigParam `json:",omitzero,inline"`
	paramUnion
}

func (u OrganizationExternalKeyNewParamsProviderConfigUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAWS, u.OfGCP, u.OfAzure)
}
func (u *OrganizationExternalKeyNewParamsProviderConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *OrganizationExternalKeyNewParamsProviderConfigUnion) asAny() any {
	if !param.IsOmitted(u.OfAWS) {
		return u.OfAWS
	} else if !param.IsOmitted(u.OfGCP) {
		return u.OfGCP
	} else if !param.IsOmitted(u.OfAzure) {
		return u.OfAzure
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationExternalKeyNewParamsProviderConfigUnion) GetKMSARN() *string {
	if vt := u.OfAWS; vt != nil {
		return &vt.KMSARN
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationExternalKeyNewParamsProviderConfigUnion) GetRegion() *string {
	if vt := u.OfAWS; vt != nil && vt.Region.Valid() {
		return &vt.Region.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationExternalKeyNewParamsProviderConfigUnion) GetTenantID() *string {
	if vt := u.OfAzure; vt != nil {
		return &vt.TenantID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationExternalKeyNewParamsProviderConfigUnion) GetVaultURI() *string {
	if vt := u.OfAzure; vt != nil {
		return &vt.VaultURI
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationExternalKeyNewParamsProviderConfigUnion) GetClientID() *string {
	if vt := u.OfAzure; vt != nil && vt.ClientID.Valid() {
		return &vt.ClientID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationExternalKeyNewParamsProviderConfigUnion) GetType() *string {
	if vt := u.OfAWS; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfGCP; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAzure; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationExternalKeyNewParamsProviderConfigUnion) GetKeyName() *string {
	if vt := u.OfGCP; vt != nil {
		return (*string)(&vt.KeyName)
	} else if vt := u.OfAzure; vt != nil {
		return (*string)(&vt.KeyName)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[OrganizationExternalKeyNewParamsProviderConfigUnion](
		"type",
		apijson.Discriminator[AWSExternalKeyConfigParam]("aws"),
		apijson.Discriminator[GCPExternalKeyConfigParam]("gcp"),
		apijson.Discriminator[AzureExternalKeyConfigParam]("azure"),
	)
}

// Data residency geo. Only `us` is supported.
type OrganizationExternalKeyNewParamsGeo string

const (
	OrganizationExternalKeyNewParamsGeoUs OrganizationExternalKeyNewParamsGeo = "us"
)

type OrganizationExternalKeyUpdateParams struct {
	// Human-friendly display name.
	DisplayName param.Opt[string] `json:"display_name,omitzero"`
	// Data residency geo. Only `us` is supported.
	//
	// Any of "us".
	Geo OrganizationExternalKeyUpdateParamsGeo `json:"geo,omitzero"`
	// KMS provider identity and auth coordinates.
	ProviderConfig OrganizationExternalKeyUpdateParamsProviderConfigUnion `json:"provider_config,omitzero"`
	paramObj
}

func (r OrganizationExternalKeyUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationExternalKeyUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationExternalKeyUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Data residency geo. Only `us` is supported.
type OrganizationExternalKeyUpdateParamsGeo string

const (
	OrganizationExternalKeyUpdateParamsGeoUs OrganizationExternalKeyUpdateParamsGeo = "us"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type OrganizationExternalKeyUpdateParamsProviderConfigUnion struct {
	OfAWS   *AWSExternalKeyConfigParam   `json:",omitzero,inline"`
	OfGCP   *GCPExternalKeyConfigParam   `json:",omitzero,inline"`
	OfAzure *AzureExternalKeyConfigParam `json:",omitzero,inline"`
	paramUnion
}

func (u OrganizationExternalKeyUpdateParamsProviderConfigUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAWS, u.OfGCP, u.OfAzure)
}
func (u *OrganizationExternalKeyUpdateParamsProviderConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *OrganizationExternalKeyUpdateParamsProviderConfigUnion) asAny() any {
	if !param.IsOmitted(u.OfAWS) {
		return u.OfAWS
	} else if !param.IsOmitted(u.OfGCP) {
		return u.OfGCP
	} else if !param.IsOmitted(u.OfAzure) {
		return u.OfAzure
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationExternalKeyUpdateParamsProviderConfigUnion) GetKMSARN() *string {
	if vt := u.OfAWS; vt != nil {
		return &vt.KMSARN
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationExternalKeyUpdateParamsProviderConfigUnion) GetRegion() *string {
	if vt := u.OfAWS; vt != nil && vt.Region.Valid() {
		return &vt.Region.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationExternalKeyUpdateParamsProviderConfigUnion) GetTenantID() *string {
	if vt := u.OfAzure; vt != nil {
		return &vt.TenantID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationExternalKeyUpdateParamsProviderConfigUnion) GetVaultURI() *string {
	if vt := u.OfAzure; vt != nil {
		return &vt.VaultURI
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationExternalKeyUpdateParamsProviderConfigUnion) GetClientID() *string {
	if vt := u.OfAzure; vt != nil && vt.ClientID.Valid() {
		return &vt.ClientID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationExternalKeyUpdateParamsProviderConfigUnion) GetType() *string {
	if vt := u.OfAWS; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfGCP; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAzure; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationExternalKeyUpdateParamsProviderConfigUnion) GetKeyName() *string {
	if vt := u.OfGCP; vt != nil {
		return (*string)(&vt.KeyName)
	} else if vt := u.OfAzure; vt != nil {
		return (*string)(&vt.KeyName)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[OrganizationExternalKeyUpdateParamsProviderConfigUnion](
		"type",
		apijson.Discriminator[AWSExternalKeyConfigParam]("aws"),
		apijson.Discriminator[GCPExternalKeyConfigParam]("gcp"),
		apijson.Discriminator[AzureExternalKeyConfigParam]("azure"),
	)
}

type OrganizationExternalKeyListParams struct {
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of results per page.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrganizationExternalKeyListParams]'s query parameters as
// `url.Values`.
func (r OrganizationExternalKeyListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
