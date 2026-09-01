// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

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

// BetaOrganizationExternalKeyService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationExternalKeyService] method instead.
type BetaOrganizationExternalKeyService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationExternalKeyService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationExternalKeyService(opts ...option.RequestOption) (r BetaOrganizationExternalKeyService) {
	r = BetaOrganizationExternalKeyService{}
	r.Options = opts
	return
}

// Create an external key config owned by the caller's organization.
func (r *BetaOrganizationExternalKeyService) New(ctx context.Context, body BetaOrganizationExternalKeyNewParams, opts ...option.RequestOption) (res *BetaExternalKey, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/external_keys?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve a single external key config in the caller's organization by ID.
func (r *BetaOrganizationExternalKeyService) Get(ctx context.Context, externalKeyID string, opts ...option.RequestOption) (res *BetaExternalKey, err error) {
	opts = slices.Concat(r.Options, opts)
	if externalKeyID == "" {
		err = errors.New("missing required external_key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/external_keys/%s?beta=true", externalKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Partially update an external key config. Omitted fields are left unchanged.
//
// `display_name` is always editable. `geo` and `provider_config` cannot be changed
// once any workspace references this config, because previously encrypted data
// requires the original key identity to decrypt.
func (r *BetaOrganizationExternalKeyService) Update(ctx context.Context, externalKeyID string, body BetaOrganizationExternalKeyUpdateParams, opts ...option.RequestOption) (res *BetaExternalKey, err error) {
	opts = slices.Concat(r.Options, opts)
	if externalKeyID == "" {
		err = errors.New("missing required external_key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/external_keys/%s?beta=true", externalKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// List external key configs in the caller's organization.
//
// Results are ordered by creation time (newest first). Use the `next_page` cursor
// from the response to fetch subsequent pages.
func (r *BetaOrganizationExternalKeyService) List(ctx context.Context, query BetaOrganizationExternalKeyListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaExternalKey], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/external_keys?beta=true"
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
func (r *BetaOrganizationExternalKeyService) ListAutoPaging(ctx context.Context, query BetaOrganizationExternalKeyListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaExternalKey] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

// Delete an external key config.
//
// The request is rejected if any workspace still references this config.
func (r *BetaOrganizationExternalKeyService) Delete(ctx context.Context, externalKeyID string, opts ...option.RequestOption) (res *BetaOrganizationExternalKeyDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if externalKeyID == "" {
		err = errors.New("missing required external_key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/external_keys/%s?beta=true", externalKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Validate an external key config against the customer's KMS.
//
// Anthropic performs an encrypt/decrypt roundtrip against the configured KMS key
// and waits up to 30 seconds for the result. The response status is `success` if
// the roundtrip succeeded, or `failure` with an error message if it failed or
// timed out.
func (r *BetaOrganizationExternalKeyService) Validate(ctx context.Context, externalKeyID string, opts ...option.RequestOption) (res *BetaOrganizationExternalKeyValidateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if externalKeyID == "" {
		err = errors.New("missing required external_key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/external_keys/%s/validate?beta=true", externalKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type BetaAWSExternalKeyConfig struct {
	// Full ARN of the AWS KMS key. On Claude Platform on AWS the key must be a
	// single-Region key in your organization's own AWS account; cross-account keys,
	// multi-Region keys, and alias ARNs are rejected.
	KMSARN string       `json:"kms_arn" api:"required"`
	Type   constant.AWS `json:"type" default:"aws"`
	// AWS region. Derived from `kms_arn` if omitted.
	Region string `json:"region" api:"nullable"`
	// IAM role ARN. Deprecated — Anthropic reaches the KMS key through its own
	// intermediate role (or, on Claude Platform on AWS, with credentials AWS issues
	// for the Workspace); this field is ignored.
	//
	// Deprecated: deprecated
	RoleARN string `json:"role_arn" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		KMSARN      respjson.Field
		Type        respjson.Field
		Region      respjson.Field
		RoleARN     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAWSExternalKeyConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaAWSExternalKeyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaAWSExternalKeyConfig to a
// BetaAWSExternalKeyConfigParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaAWSExternalKeyConfigParam.Overrides()
func (r BetaAWSExternalKeyConfig) ToParam() BetaAWSExternalKeyConfigParam {
	return param.Override[BetaAWSExternalKeyConfigParam](json.RawMessage(r.RawJSON()))
}

// The properties KMSARN, Type are required.
type BetaAWSExternalKeyConfigParam struct {
	// Full ARN of the AWS KMS key. On Claude Platform on AWS the key must be a
	// single-Region key in your organization's own AWS account; cross-account keys,
	// multi-Region keys, and alias ARNs are rejected.
	KMSARN string `json:"kms_arn" api:"required"`
	// AWS region. Derived from `kms_arn` if omitted.
	Region param.Opt[string] `json:"region,omitzero"`
	// IAM role ARN. Deprecated — Anthropic reaches the KMS key through its own
	// intermediate role (or, on Claude Platform on AWS, with credentials AWS issues
	// for the Workspace); this field is ignored.
	//
	// Deprecated: deprecated
	RoleARN param.Opt[string] `json:"role_arn,omitzero"`
	// This field can be elided, and will marshal its zero value as "aws".
	Type constant.AWS `json:"type" default:"aws"`
	paramObj
}

func (r BetaAWSExternalKeyConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaAWSExternalKeyConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAWSExternalKeyConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAzureExternalKeyConfig struct {
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
func (r BetaAzureExternalKeyConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaAzureExternalKeyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Azure Key Vault provider configuration.
//
// The properties KeyName, TenantID, Type, VaultURI are required.
type BetaAzureExternalKeyConfigParam struct {
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

func (r BetaAzureExternalKeyConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaAzureExternalKeyConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAzureExternalKeyConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// CMEK external key config belonging to the caller's organization.
//
// Configs are organization-scoped. Workspaces attach to a config; once any
// workspace references it, the provider fields become effectively immutable
// (existing encrypted data needs the config for decrypt).
type BetaExternalKey struct {
	// Identifier of the external key config. A tagged ID prefixed `ekey_`, or — for
	// organizations on the Claude Platform on AWS — the AWS KMS key ARN.
	ID string `json:"id" api:"required"`
	// Whether any workspace uses this config to encrypt its data — counting live and
	// archived workspaces (an archived workspace's data remains encrypted under the
	// config), excluding deleted ones. Only an attached config is used by the
	// encryption path; an `unattached` config is inert and can be deleted.
	Attachment BetaExternalKeyAttachmentUnion `json:"attachment" api:"required"`
	CreatedAt  time.Time                      `json:"created_at" api:"required" format:"date-time"`
	// Human-friendly display name. Null if none was set.
	DisplayName string `json:"display_name" api:"required"`
	// Data residency geo. Selects which regional validator handles this key's
	// encrypt/decrypt roundtrips.
	Geo string `json:"geo" api:"required"`
	// KMS provider identity and auth coordinates.
	ProviderConfig BetaExternalKeyProviderConfigUnion `json:"provider_config" api:"required"`
	Type           constant.ExternalKey               `json:"type" default:"external_key"`
	UpdatedAt      time.Time                          `json:"updated_at" api:"required" format:"date-time"`
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
func (r BetaExternalKey) RawJSON() string { return r.JSON.raw }
func (r *BetaExternalKey) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaExternalKeyAttachmentUnion contains all possible properties and values from
// [BetaExternalKeyAttachedAttachment], [BetaExternalKeyUnattachedAttachment].
//
// Use the [BetaExternalKeyAttachmentUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaExternalKeyAttachmentUnion struct {
	// Any of "attached", "unattached".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaExternalKeyAttachment is implemented by each variant of
// [BetaExternalKeyAttachmentUnion] to add type safety for the return type of
// [BetaExternalKeyAttachmentUnion.AsAny]
type anyBetaExternalKeyAttachment interface {
	implBetaExternalKeyAttachmentUnion()
}

func (BetaExternalKeyAttachedAttachment) implBetaExternalKeyAttachmentUnion()   {}
func (BetaExternalKeyUnattachedAttachment) implBetaExternalKeyAttachmentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaExternalKeyAttachmentUnion.AsAny().(type) {
//	case anthropic.BetaExternalKeyAttachedAttachment:
//	case anthropic.BetaExternalKeyUnattachedAttachment:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaExternalKeyAttachmentUnion) AsAny() anyBetaExternalKeyAttachment {
	switch u.Type {
	case "attached":
		return u.AsAttached()
	case "unattached":
		return u.AsUnattached()
	}
	return nil
}

func (u BetaExternalKeyAttachmentUnion) AsAttached() (v BetaExternalKeyAttachedAttachment) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaExternalKeyAttachmentUnion) AsUnattached() (v BetaExternalKeyUnattachedAttachment) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaExternalKeyAttachmentUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaExternalKeyAttachmentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaExternalKeyProviderConfigUnion contains all possible properties and values
// from [BetaAWSExternalKeyConfig], [BetaGCPExternalKeyConfig],
// [BetaAzureExternalKeyConfig].
//
// Use the [BetaExternalKeyProviderConfigUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaExternalKeyProviderConfigUnion struct {
	// This field is from variant [BetaAWSExternalKeyConfig].
	KMSARN string `json:"kms_arn"`
	// Any of "aws", "gcp", "azure".
	Type string `json:"type"`
	// This field is from variant [BetaAWSExternalKeyConfig].
	Region string `json:"region"`
	// This field is from variant [BetaAWSExternalKeyConfig].
	RoleARN string `json:"role_arn"`
	KeyName string `json:"key_name"`
	// This field is from variant [BetaAzureExternalKeyConfig].
	TenantID string `json:"tenant_id"`
	// This field is from variant [BetaAzureExternalKeyConfig].
	VaultURI string `json:"vault_uri"`
	// This field is from variant [BetaAzureExternalKeyConfig].
	ClientID string `json:"client_id"`
	JSON     struct {
		KMSARN   respjson.Field
		Type     respjson.Field
		Region   respjson.Field
		RoleARN  respjson.Field
		KeyName  respjson.Field
		TenantID respjson.Field
		VaultURI respjson.Field
		ClientID respjson.Field
		raw      string
	} `json:"-"`
}

// anyBetaExternalKeyProviderConfig is implemented by each variant of
// [BetaExternalKeyProviderConfigUnion] to add type safety for the return type of
// [BetaExternalKeyProviderConfigUnion.AsAny]
type anyBetaExternalKeyProviderConfig interface {
	implBetaExternalKeyProviderConfigUnion()
}

func (BetaAWSExternalKeyConfig) implBetaExternalKeyProviderConfigUnion()   {}
func (BetaGCPExternalKeyConfig) implBetaExternalKeyProviderConfigUnion()   {}
func (BetaAzureExternalKeyConfig) implBetaExternalKeyProviderConfigUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaExternalKeyProviderConfigUnion.AsAny().(type) {
//	case anthropic.BetaAWSExternalKeyConfig:
//	case anthropic.BetaGCPExternalKeyConfig:
//	case anthropic.BetaAzureExternalKeyConfig:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaExternalKeyProviderConfigUnion) AsAny() anyBetaExternalKeyProviderConfig {
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

func (u BetaExternalKeyProviderConfigUnion) AsAWS() (v BetaAWSExternalKeyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaExternalKeyProviderConfigUnion) AsGCP() (v BetaGCPExternalKeyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaExternalKeyProviderConfigUnion) AsAzure() (v BetaAzureExternalKeyConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaExternalKeyProviderConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaExternalKeyProviderConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaExternalKeyAttachedAttachment struct {
	Type constant.Attached `json:"type" default:"attached"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaExternalKeyAttachedAttachment) RawJSON() string { return r.JSON.raw }
func (r *BetaExternalKeyAttachedAttachment) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaExternalKeyUnattachedAttachment struct {
	Type constant.Unattached `json:"type" default:"unattached"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaExternalKeyUnattachedAttachment) RawJSON() string { return r.JSON.raw }
func (r *BetaExternalKeyUnattachedAttachment) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaGCPExternalKeyConfig struct {
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
func (r BetaGCPExternalKeyConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaGCPExternalKeyConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaGCPExternalKeyConfig to a
// BetaGCPExternalKeyConfigParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaGCPExternalKeyConfigParam.Overrides()
func (r BetaGCPExternalKeyConfig) ToParam() BetaGCPExternalKeyConfigParam {
	return param.Override[BetaGCPExternalKeyConfigParam](json.RawMessage(r.RawJSON()))
}

// The properties KeyName, Type are required.
type BetaGCPExternalKeyConfigParam struct {
	// Full resource name of the Cloud KMS key.
	KeyName string `json:"key_name" api:"required"`
	// This field can be elided, and will marshal its zero value as "gcp".
	Type constant.GCP `json:"type" default:"gcp"`
	paramObj
}

func (r BetaGCPExternalKeyConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaGCPExternalKeyConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaGCPExternalKeyConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationExternalKeyDeleteResponse struct {
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
func (r BetaOrganizationExternalKeyDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationExternalKeyDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Result of a validation roundtrip against the customer's KMS.
//
// HTTP 200 for both outcomes — the operation completed; `status` says whether the
// key works.
type BetaOrganizationExternalKeyValidateResponse struct {
	// Error message when status is `failure`. Null otherwise.
	Error string `json:"error" api:"required"`
	// `success` — encrypt/decrypt roundtrip succeeded. `failure` — the roundtrip
	// failed or timed out; see `error`.
	//
	// Any of "failure", "success".
	Status BetaOrganizationExternalKeyValidateResponseStatus `json:"status" api:"required"`
	Type   constant.ExternalKeyValidation                    `json:"type" default:"external_key_validation"`
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
func (r BetaOrganizationExternalKeyValidateResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationExternalKeyValidateResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// `success` — encrypt/decrypt roundtrip succeeded. `failure` — the roundtrip
// failed or timed out; see `error`.
type BetaOrganizationExternalKeyValidateResponseStatus string

const (
	BetaOrganizationExternalKeyValidateResponseStatusFailure BetaOrganizationExternalKeyValidateResponseStatus = "failure"
	BetaOrganizationExternalKeyValidateResponseStatusSuccess BetaOrganizationExternalKeyValidateResponseStatus = "success"
)

type BetaOrganizationExternalKeyNewParams struct {
	// KMS provider identity and auth coordinates.
	ProviderConfig BetaOrganizationExternalKeyNewParamsProviderConfigUnion `json:"provider_config,omitzero" api:"required"`
	// Human-friendly display name.
	DisplayName param.Opt[string] `json:"display_name,omitzero"`
	// Data residency geo. Only `us` is supported.
	//
	// Any of "us".
	Geo BetaOrganizationExternalKeyNewParamsGeo `json:"geo,omitzero"`
	paramObj
}

func (r BetaOrganizationExternalKeyNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationExternalKeyNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationExternalKeyNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaOrganizationExternalKeyNewParamsProviderConfigUnion struct {
	OfAWS   *BetaAWSExternalKeyConfigParam   `json:",omitzero,inline"`
	OfGCP   *BetaGCPExternalKeyConfigParam   `json:",omitzero,inline"`
	OfAzure *BetaAzureExternalKeyConfigParam `json:",omitzero,inline"`
	paramUnion
}

func (u BetaOrganizationExternalKeyNewParamsProviderConfigUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAWS, u.OfGCP, u.OfAzure)
}
func (u *BetaOrganizationExternalKeyNewParamsProviderConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaOrganizationExternalKeyNewParamsProviderConfigUnion) asAny() any {
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
func (u BetaOrganizationExternalKeyNewParamsProviderConfigUnion) GetKMSARN() *string {
	if vt := u.OfAWS; vt != nil {
		return &vt.KMSARN
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalKeyNewParamsProviderConfigUnion) GetRegion() *string {
	if vt := u.OfAWS; vt != nil && vt.Region.Valid() {
		return &vt.Region.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalKeyNewParamsProviderConfigUnion) GetRoleARN() *string {
	if vt := u.OfAWS; vt != nil && vt.RoleARN.Valid() {
		return &vt.RoleARN.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalKeyNewParamsProviderConfigUnion) GetTenantID() *string {
	if vt := u.OfAzure; vt != nil {
		return &vt.TenantID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalKeyNewParamsProviderConfigUnion) GetVaultURI() *string {
	if vt := u.OfAzure; vt != nil {
		return &vt.VaultURI
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalKeyNewParamsProviderConfigUnion) GetClientID() *string {
	if vt := u.OfAzure; vt != nil && vt.ClientID.Valid() {
		return &vt.ClientID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalKeyNewParamsProviderConfigUnion) GetType() *string {
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
func (u BetaOrganizationExternalKeyNewParamsProviderConfigUnion) GetKeyName() *string {
	if vt := u.OfGCP; vt != nil {
		return (*string)(&vt.KeyName)
	} else if vt := u.OfAzure; vt != nil {
		return (*string)(&vt.KeyName)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaOrganizationExternalKeyNewParamsProviderConfigUnion](
		"type",
		apijson.Discriminator[BetaAWSExternalKeyConfigParam]("aws"),
		apijson.Discriminator[BetaGCPExternalKeyConfigParam]("gcp"),
		apijson.Discriminator[BetaAzureExternalKeyConfigParam]("azure"),
	)
}

// Data residency geo. Only `us` is supported.
type BetaOrganizationExternalKeyNewParamsGeo string

const (
	BetaOrganizationExternalKeyNewParamsGeoUs BetaOrganizationExternalKeyNewParamsGeo = "us"
)

type BetaOrganizationExternalKeyUpdateParams struct {
	// Human-friendly display name.
	DisplayName param.Opt[string] `json:"display_name,omitzero"`
	// Data residency geo. Only `us` is supported.
	//
	// Any of "us".
	Geo BetaOrganizationExternalKeyUpdateParamsGeo `json:"geo,omitzero"`
	// KMS provider identity and auth coordinates.
	ProviderConfig BetaOrganizationExternalKeyUpdateParamsProviderConfigUnion `json:"provider_config,omitzero"`
	paramObj
}

func (r BetaOrganizationExternalKeyUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationExternalKeyUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationExternalKeyUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Data residency geo. Only `us` is supported.
type BetaOrganizationExternalKeyUpdateParamsGeo string

const (
	BetaOrganizationExternalKeyUpdateParamsGeoUs BetaOrganizationExternalKeyUpdateParamsGeo = "us"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaOrganizationExternalKeyUpdateParamsProviderConfigUnion struct {
	OfAWS   *BetaAWSExternalKeyConfigParam   `json:",omitzero,inline"`
	OfGCP   *BetaGCPExternalKeyConfigParam   `json:",omitzero,inline"`
	OfAzure *BetaAzureExternalKeyConfigParam `json:",omitzero,inline"`
	paramUnion
}

func (u BetaOrganizationExternalKeyUpdateParamsProviderConfigUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAWS, u.OfGCP, u.OfAzure)
}
func (u *BetaOrganizationExternalKeyUpdateParamsProviderConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaOrganizationExternalKeyUpdateParamsProviderConfigUnion) asAny() any {
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
func (u BetaOrganizationExternalKeyUpdateParamsProviderConfigUnion) GetKMSARN() *string {
	if vt := u.OfAWS; vt != nil {
		return &vt.KMSARN
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalKeyUpdateParamsProviderConfigUnion) GetRegion() *string {
	if vt := u.OfAWS; vt != nil && vt.Region.Valid() {
		return &vt.Region.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalKeyUpdateParamsProviderConfigUnion) GetRoleARN() *string {
	if vt := u.OfAWS; vt != nil && vt.RoleARN.Valid() {
		return &vt.RoleARN.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalKeyUpdateParamsProviderConfigUnion) GetTenantID() *string {
	if vt := u.OfAzure; vt != nil {
		return &vt.TenantID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalKeyUpdateParamsProviderConfigUnion) GetVaultURI() *string {
	if vt := u.OfAzure; vt != nil {
		return &vt.VaultURI
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalKeyUpdateParamsProviderConfigUnion) GetClientID() *string {
	if vt := u.OfAzure; vt != nil && vt.ClientID.Valid() {
		return &vt.ClientID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationExternalKeyUpdateParamsProviderConfigUnion) GetType() *string {
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
func (u BetaOrganizationExternalKeyUpdateParamsProviderConfigUnion) GetKeyName() *string {
	if vt := u.OfGCP; vt != nil {
		return (*string)(&vt.KeyName)
	} else if vt := u.OfAzure; vt != nil {
		return (*string)(&vt.KeyName)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaOrganizationExternalKeyUpdateParamsProviderConfigUnion](
		"type",
		apijson.Discriminator[BetaAWSExternalKeyConfigParam]("aws"),
		apijson.Discriminator[BetaGCPExternalKeyConfigParam]("gcp"),
		apijson.Discriminator[BetaAzureExternalKeyConfigParam]("azure"),
	)
}

type BetaOrganizationExternalKeyListParams struct {
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Number of results per page.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationExternalKeyListParams]'s query parameters
// as `url.Values`.
func (r BetaOrganizationExternalKeyListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
