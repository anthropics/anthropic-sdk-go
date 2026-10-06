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

// OrganizationWorkspaceService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationWorkspaceService] method instead.
type OrganizationWorkspaceService struct {
	Options         []option.RequestOption
	RateLimits      OrganizationWorkspaceRateLimitService
	Members         OrganizationWorkspaceMemberService
	ServiceAccounts OrganizationWorkspaceServiceAccountService
}

// NewOrganizationWorkspaceService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewOrganizationWorkspaceService(opts ...option.RequestOption) (r OrganizationWorkspaceService) {
	r = OrganizationWorkspaceService{}
	r.Options = opts
	r.RateLimits = NewOrganizationWorkspaceRateLimitService(opts...)
	r.Members = NewOrganizationWorkspaceMemberService(opts...)
	r.ServiceAccounts = NewOrganizationWorkspaceServiceAccountService(opts...)
	return
}

// Create Workspace
func (r *OrganizationWorkspaceService) New(ctx context.Context, body OrganizationWorkspaceNewParams, opts ...option.RequestOption) (res *Workspace, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/workspaces"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get Workspace
func (r *OrganizationWorkspaceService) Get(ctx context.Context, workspaceID string, opts ...option.RequestOption) (res *Workspace, err error) {
	opts = slices.Concat(r.Options, opts)
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s", url.PathEscape(workspaceID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Workspace
func (r *OrganizationWorkspaceService) Update(ctx context.Context, workspaceID string, body OrganizationWorkspaceUpdateParams, opts ...option.RequestOption) (res *Workspace, err error) {
	opts = slices.Concat(r.Options, opts)
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s", url.PathEscape(workspaceID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// List Workspaces
func (r *OrganizationWorkspaceService) List(ctx context.Context, query OrganizationWorkspaceListParams, opts ...option.RequestOption) (res *pagination.Page[Workspace], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/workspaces"
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

// List Workspaces
func (r *OrganizationWorkspaceService) ListAutoPaging(ctx context.Context, query OrganizationWorkspaceListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Workspace] {
	return pagination.NewPageAutoPager(r.List(ctx, query, opts...))
}

// Archive Workspace
func (r *OrganizationWorkspaceService) Archive(ctx context.Context, workspaceID string, opts ...option.RequestOption) (res *Workspace, err error) {
	opts = slices.Concat(r.Options, opts)
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/workspaces/%s/archive", url.PathEscape(workspaceID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type AllowedInferenceGeo string

const (
	AllowedInferenceGeoGlobal AllowedInferenceGeo = "global"
	AllowedInferenceGeoUs     AllowedInferenceGeo = "us"
)

type DataResidency struct {
	// Permitted inference geo values. 'unrestricted' means all geos are allowed.
	AllowedInferenceGeos DataResidencyAllowedInferenceGeosUnion `json:"allowed_inference_geos" api:"required"`
	// Default inference geo applied when requests omit the parameter.
	//
	// Any of "global", "us".
	DefaultInferenceGeo DataResidencyDefaultInferenceGeo `json:"default_inference_geo" api:"required"`
	// Geographic region for workspace data storage. Immutable after creation.
	//
	// Any of "us".
	WorkspaceGeo DataResidencyWorkspaceGeo `json:"workspace_geo" api:"required"`
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
func (r DataResidency) RawJSON() string { return r.JSON.raw }
func (r *DataResidency) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DataResidencyAllowedInferenceGeosUnion contains all possible properties and
// values from [[]AllowedInferenceGeo], [constant.Unrestricted].
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfGeos OfUnrestricted]
type DataResidencyAllowedInferenceGeosUnion struct {
	// This field will be present if the value is a [[]AllowedInferenceGeo] instead of
	// an object.
	OfGeos []AllowedInferenceGeo `json:",inline"`
	// This field will be present if the value is a [constant.Unrestricted] instead of
	// an object.
	OfUnrestricted constant.Unrestricted `json:",inline"`
	JSON           struct {
		OfGeos         respjson.Field
		OfUnrestricted respjson.Field
		raw            string
	} `json:"-"`
}

func (u DataResidencyAllowedInferenceGeosUnion) AsGeos() (v []AllowedInferenceGeo) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DataResidencyAllowedInferenceGeosUnion) AsUnrestricted() (v constant.Unrestricted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DataResidencyAllowedInferenceGeosUnion) RawJSON() string { return u.JSON.raw }

func (r *DataResidencyAllowedInferenceGeosUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Default inference geo applied when requests omit the parameter.
type DataResidencyDefaultInferenceGeo string

const (
	DataResidencyDefaultInferenceGeoGlobal DataResidencyDefaultInferenceGeo = "global"
	DataResidencyDefaultInferenceGeoUs     DataResidencyDefaultInferenceGeo = "us"
)

// Geographic region for workspace data storage. Immutable after creation.
type DataResidencyWorkspaceGeo string

const (
	DataResidencyWorkspaceGeoUs DataResidencyWorkspaceGeo = "us"
)

type DataResidencyCreateConfigParam struct {
	// Permitted inference geo values. Defaults to 'unrestricted' if omitted, which
	// allows all geos. Use the string 'unrestricted' to allow all geos, or a list of
	// specific geos.
	AllowedInferenceGeos DataResidencyCreateConfigAllowedInferenceGeosUnionParam `json:"allowed_inference_geos,omitzero"`
	// Default inference geo applied when requests omit the parameter. Defaults to
	// 'global' if omitted. Must be a member of `allowed_inference_geos` unless
	// `allowed_inference_geos` is `"unrestricted"`.
	//
	// Any of "global", "us".
	DefaultInferenceGeo DataResidencyCreateConfigDefaultInferenceGeo `json:"default_inference_geo,omitzero"`
	// Geographic region for workspace data storage. Immutable after creation. Defaults
	// to 'us' if omitted.
	//
	// Any of "us".
	WorkspaceGeo DataResidencyCreateConfigWorkspaceGeo `json:"workspace_geo,omitzero"`
	paramObj
}

func (r DataResidencyCreateConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow DataResidencyCreateConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DataResidencyCreateConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DataResidencyCreateConfigAllowedInferenceGeosUnionParam struct {
	OfGeos []AllowedInferenceGeo `json:",omitzero,inline"`
	// Construct this variant with constant.ValueOf[constant.Unrestricted]()
	OfUnrestricted constant.Unrestricted `json:",omitzero,inline"`
	paramUnion
}

func (u DataResidencyCreateConfigAllowedInferenceGeosUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfGeos, u.OfUnrestricted)
}
func (u *DataResidencyCreateConfigAllowedInferenceGeosUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DataResidencyCreateConfigAllowedInferenceGeosUnionParam) asAny() any {
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
type DataResidencyCreateConfigDefaultInferenceGeo string

const (
	DataResidencyCreateConfigDefaultInferenceGeoGlobal DataResidencyCreateConfigDefaultInferenceGeo = "global"
	DataResidencyCreateConfigDefaultInferenceGeoUs     DataResidencyCreateConfigDefaultInferenceGeo = "us"
)

// Geographic region for workspace data storage. Immutable after creation. Defaults
// to 'us' if omitted.
type DataResidencyCreateConfigWorkspaceGeo string

const (
	DataResidencyCreateConfigWorkspaceGeoUs DataResidencyCreateConfigWorkspaceGeo = "us"
)

type DataResidencyUpdateConfigParam struct {
	// Permitted inference geo values. Use 'unrestricted' to allow all geos, or a list
	// of specific geos.
	AllowedInferenceGeos DataResidencyUpdateConfigAllowedInferenceGeosUnionParam `json:"allowed_inference_geos,omitzero"`
	// Default inference geo applied when requests omit the parameter. Must be a member
	// of `allowed_inference_geos` unless `allowed_inference_geos` is `"unrestricted"`.
	//
	// Any of "global", "us".
	DefaultInferenceGeo DataResidencyUpdateConfigDefaultInferenceGeo `json:"default_inference_geo,omitzero"`
	paramObj
}

func (r DataResidencyUpdateConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow DataResidencyUpdateConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DataResidencyUpdateConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DataResidencyUpdateConfigAllowedInferenceGeosUnionParam struct {
	OfGeos []AllowedInferenceGeo `json:",omitzero,inline"`
	// Construct this variant with constant.ValueOf[constant.Unrestricted]()
	OfUnrestricted constant.Unrestricted `json:",omitzero,inline"`
	paramUnion
}

func (u DataResidencyUpdateConfigAllowedInferenceGeosUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfGeos, u.OfUnrestricted)
}
func (u *DataResidencyUpdateConfigAllowedInferenceGeosUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DataResidencyUpdateConfigAllowedInferenceGeosUnionParam) asAny() any {
	if !param.IsOmitted(u.OfGeos) {
		return &u.OfGeos
	} else if !param.IsOmitted(u.OfUnrestricted) {
		return &u.OfUnrestricted
	}
	return nil
}

// Default inference geo applied when requests omit the parameter. Must be a member
// of `allowed_inference_geos` unless `allowed_inference_geos` is `"unrestricted"`.
type DataResidencyUpdateConfigDefaultInferenceGeo string

const (
	DataResidencyUpdateConfigDefaultInferenceGeoGlobal DataResidencyUpdateConfigDefaultInferenceGeo = "global"
	DataResidencyUpdateConfigDefaultInferenceGeoUs     DataResidencyUpdateConfigDefaultInferenceGeo = "us"
)

type NoBillingWorkspaceRole string

const (
	NoBillingWorkspaceRoleWorkspaceAdmin               NoBillingWorkspaceRole = "workspace_admin"
	NoBillingWorkspaceRoleWorkspaceDeveloper           NoBillingWorkspaceRole = "workspace_developer"
	NoBillingWorkspaceRoleWorkspaceRestrictedDeveloper NoBillingWorkspaceRole = "workspace_restricted_developer"
	NoBillingWorkspaceRoleWorkspaceUser                NoBillingWorkspaceRole = "workspace_user"
)

type Workspace struct {
	// ID of the Workspace.
	ID string `json:"id" api:"required"`
	// RFC 3339 datetime string indicating when the Workspace was archived, or `null`
	// if the Workspace is not archived.
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// Identifier for this Workspace's encryption compartment. When you configure a
	// customer-managed encryption key (CMEK) on AWS, reference this value in your KMS
	// key-policy condition so the key is scoped to this compartment. On GCP and Azure,
	// Anthropic enforces the compartment binding automatically; you do not need to
	// reference this value in your key configuration. See the CMEK integration guide
	// for the required key configuration; unless your organization is on Claude
	// Platform on AWS, it includes a separate value used during key validation. On
	// Claude Platform on AWS there is no separate validation value: the key is
	// validated against this Workspace's own value when it is attached, so if your key
	// policy uses the compartment condition, add this value to it before attaching the
	// key.
	CompartmentID string `json:"compartment_id" api:"required"`
	// RFC 3339 datetime string indicating when the Workspace was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Data residency configuration.
	DataResidency DataResidency `json:"data_residency" api:"required"`
	// Hex color code representing the Workspace in the Anthropic Console.
	DisplayColor string `json:"display_color" api:"required"`
	// ID of the customer-managed encryption key (CMEK) configuration to use for this
	// Workspace. Setting this field requires CMEK to be enabled for your organization.
	// When set, data stored for this Workspace is encrypted with the referenced key.
	// Create key configurations with the External Keys API. On Claude Platform on AWS
	// the value is the AWS KMS key ARN, and the key must be a single-Region key in the
	// same AWS account and Region as the Workspace. On that platform the key is
	// validated against this Workspace when it is attached, so a key-policy problem is
	// reported as an error on this request. This field is write-once: once a key is
	// attached to a Workspace it cannot be detached or replaced. To rotate key
	// material, rotate the underlying key on your cloud KMS; the `external_key_id`
	// stays the same.
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
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		ArchivedAt    respjson.Field
		CompartmentID respjson.Field
		CreatedAt     respjson.Field
		DataResidency respjson.Field
		DisplayColor  respjson.Field
		ExternalKeyID respjson.Field
		Name          respjson.Field
		Tags          respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Workspace) RawJSON() string { return r.JSON.raw }
func (r *Workspace) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkspaceMember struct {
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
	WorkspaceRole WorkspaceRole `json:"workspace_role" api:"required"`
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
func (r WorkspaceMember) RawJSON() string { return r.JSON.raw }
func (r *WorkspaceMember) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type WorkspaceRole string

const (
	WorkspaceRoleWorkspaceAdmin               WorkspaceRole = "workspace_admin"
	WorkspaceRoleWorkspaceBilling             WorkspaceRole = "workspace_billing"
	WorkspaceRoleWorkspaceDeveloper           WorkspaceRole = "workspace_developer"
	WorkspaceRoleWorkspaceRestrictedDeveloper WorkspaceRole = "workspace_restricted_developer"
	WorkspaceRoleWorkspaceUser                WorkspaceRole = "workspace_user"
)

type OrganizationWorkspaceNewParams struct {
	// Name of the Workspace.
	Name string `json:"name" api:"required"`
	// Hex color code representing the Workspace in the Anthropic Console.
	DisplayColor param.Opt[string] `json:"display_color,omitzero"`
	// ID of the customer-managed encryption key (CMEK) configuration to use for this
	// Workspace. Setting this field requires CMEK to be enabled for your organization.
	// When set, data stored for this Workspace is encrypted with the referenced key.
	// Create key configurations with the External Keys API. On Claude Platform on AWS
	// the value is the AWS KMS key ARN, and the key must be a single-Region key in the
	// same AWS account and Region as the Workspace. On that platform the key is
	// validated against this Workspace when it is attached, so a key-policy problem is
	// reported as an error on this request. This field is write-once: once a key is
	// attached to a Workspace it cannot be detached or replaced. To rotate key
	// material, rotate the underlying key on your cloud KMS; the `external_key_id`
	// stays the same.
	ExternalKeyID param.Opt[string] `json:"external_key_id,omitzero"`
	// User-defined tags as string key-value pairs. Keys may not begin with
	// `anthropic`.
	Tags map[string]string `json:"tags,omitzero"`
	// Data residency configuration for the workspace. If omitted, defaults to
	// `workspace_geo: "us"`, `allowed_inference_geos: "unrestricted"`, and
	// `default_inference_geo: "global"`.
	DataResidency DataResidencyCreateConfigParam `json:"data_residency,omitzero"`
	paramObj
}

func (r OrganizationWorkspaceNewParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationWorkspaceNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationWorkspaceNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationWorkspaceUpdateParams struct {
	// Hex color code representing the Workspace in the Anthropic Console.
	DisplayColor param.Opt[string] `json:"display_color,omitzero"`
	// ID of the customer-managed encryption key (CMEK) configuration to use for this
	// Workspace. Setting this field requires CMEK to be enabled for your organization.
	// When set, data stored for this Workspace is encrypted with the referenced key.
	// Create key configurations with the External Keys API. On Claude Platform on AWS
	// the value is the AWS KMS key ARN, and the key must be a single-Region key in the
	// same AWS account and Region as the Workspace. On that platform the key is
	// validated against this Workspace when it is attached, so a key-policy problem is
	// reported as an error on this request. This field is write-once: once a key is
	// attached to a Workspace it cannot be detached or replaced. To rotate key
	// material, rotate the underlying key on your cloud KMS; the `external_key_id`
	// stays the same.
	ExternalKeyID param.Opt[string] `json:"external_key_id,omitzero"`
	// Name of the Workspace.
	Name param.Opt[string] `json:"name,omitzero"`
	// User-defined tags as string key-value pairs. Keys may not begin with
	// `anthropic`.
	Tags map[string]string `json:"tags,omitzero"`
	// Data residency configuration for the workspace.
	DataResidency DataResidencyUpdateConfigParam `json:"data_residency,omitzero"`
	paramObj
}

func (r OrganizationWorkspaceUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationWorkspaceUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationWorkspaceUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationWorkspaceListParams struct {
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
	paramObj
}

// URLQuery serializes [OrganizationWorkspaceListParams]'s query parameters as
// `url.Values`.
func (r OrganizationWorkspaceListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
