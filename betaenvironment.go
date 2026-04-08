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

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// BetaEnvironmentService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaEnvironmentService] method instead.
type BetaEnvironmentService struct {
	Options []option.RequestOption
}

// NewBetaEnvironmentService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaEnvironmentService(opts ...option.RequestOption) (r BetaEnvironmentService) {
	r = BetaEnvironmentService{}
	r.Options = opts
	return
}

// Create a new environment with the specified configuration.
func (r *BetaEnvironmentService) New(ctx context.Context, params BetaEnvironmentNewParams, opts ...option.RequestOption) (res *BetaEnvironment, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	path := "v1/environments?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Retrieve a specific environment by ID.
func (r *BetaEnvironmentService) Get(ctx context.Context, environmentID string, query BetaEnvironmentGetParams, opts ...option.RequestOption) (res *BetaEnvironment, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if environmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/environments/%s?beta=true", environmentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update an existing environment's configuration.
func (r *BetaEnvironmentService) Update(ctx context.Context, environmentID string, params BetaEnvironmentUpdateParams, opts ...option.RequestOption) (res *BetaEnvironment, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if environmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/environments/%s?beta=true", environmentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List environments with pagination support.
func (r *BetaEnvironmentService) List(ctx context.Context, params BetaEnvironmentListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaEnvironment], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/environments?beta=true"
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

// List environments with pagination support.
func (r *BetaEnvironmentService) ListAutoPaging(ctx context.Context, params BetaEnvironmentListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaEnvironment] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Delete an environment by ID. Returns a confirmation of the deletion.
func (r *BetaEnvironmentService) Delete(ctx context.Context, environmentID string, body BetaEnvironmentDeleteParams, opts ...option.RequestOption) (res *BetaEnvironmentDeleteResponse, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if environmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/environments/%s?beta=true", environmentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Archive an environment by ID. Archived environments cannot be used to create new
// sessions.
func (r *BetaEnvironmentService) Archive(ctx context.Context, environmentID string, body BetaEnvironmentArchiveParams, opts ...option.RequestOption) (res *BetaEnvironment, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if environmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/environments/%s/archive?beta=true", environmentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// `cloud` environment configuration.
type BetaCloudConfig struct {
	// Network configuration policy.
	Networking BetaCloudConfigNetworkingUnion `json:"networking" api:"required"`
	// Package manager configuration.
	Packages BetaPackages `json:"packages" api:"required"`
	// Environment type
	Type constant.Cloud `json:"type" default:"cloud"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Networking  respjson.Field
		Packages    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaCloudConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaCloudConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaCloudConfigNetworkingUnion contains all possible properties and values from
// [BetaUnrestrictedNetwork], [BetaLimitedNetwork].
//
// Use the [BetaCloudConfigNetworkingUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaCloudConfigNetworkingUnion struct {
	// Any of "unrestricted", "limited".
	Type string `json:"type"`
	// This field is from variant [BetaLimitedNetwork].
	AllowMCPServers bool `json:"allow_mcp_servers"`
	// This field is from variant [BetaLimitedNetwork].
	AllowPackageManagers bool `json:"allow_package_managers"`
	// This field is from variant [BetaLimitedNetwork].
	AllowedHosts []string `json:"allowed_hosts"`
	JSON         struct {
		Type                 respjson.Field
		AllowMCPServers      respjson.Field
		AllowPackageManagers respjson.Field
		AllowedHosts         respjson.Field
		raw                  string
	} `json:"-"`
}

// anyBetaCloudConfigNetworking is implemented by each variant of
// [BetaCloudConfigNetworkingUnion] to add type safety for the return type of
// [BetaCloudConfigNetworkingUnion.AsAny]
type anyBetaCloudConfigNetworking interface {
	implBetaCloudConfigNetworkingUnion()
}

func (BetaUnrestrictedNetwork) implBetaCloudConfigNetworkingUnion() {}
func (BetaLimitedNetwork) implBetaCloudConfigNetworkingUnion()      {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaCloudConfigNetworkingUnion.AsAny().(type) {
//	case anthropic.BetaUnrestrictedNetwork:
//	case anthropic.BetaLimitedNetwork:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaCloudConfigNetworkingUnion) AsAny() anyBetaCloudConfigNetworking {
	switch u.Type {
	case "unrestricted":
		return u.AsUnrestricted()
	case "limited":
		return u.AsLimited()
	}
	return nil
}

func (u BetaCloudConfigNetworkingUnion) AsUnrestricted() (v BetaUnrestrictedNetwork) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaCloudConfigNetworkingUnion) AsLimited() (v BetaLimitedNetwork) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaCloudConfigNetworkingUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaCloudConfigNetworkingUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Request params for `cloud` environment configuration.
//
// Fields default to null; on update, omitted fields preserve the existing value.
//
// The property Type is required.
type BetaCloudConfigParams struct {
	// Network configuration policy. Omit on update to preserve the existing value.
	Networking BetaCloudConfigParamsNetworkingUnion `json:"networking,omitzero"`
	// Specify packages (and optionally their versions) available in this environment.
	//
	// When versioning, use the version semantics relevant for the package manager,
	// e.g. for `pip` use `package==1.0.0`. You are responsible for validating the
	// package and version exist. Unversioned installs the latest.
	Packages BetaPackagesParams `json:"packages,omitzero"`
	// Environment type
	//
	// This field can be elided, and will marshal its zero value as "cloud".
	Type constant.Cloud `json:"type" default:"cloud"`
	paramObj
}

func (r BetaCloudConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaCloudConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaCloudConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaCloudConfigParamsNetworkingUnion struct {
	OfUnrestricted *BetaUnrestrictedNetworkParam `json:",omitzero,inline"`
	OfLimited      *BetaLimitedNetworkParams     `json:",omitzero,inline"`
	paramUnion
}

func (u BetaCloudConfigParamsNetworkingUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfUnrestricted, u.OfLimited)
}
func (u *BetaCloudConfigParamsNetworkingUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaCloudConfigParamsNetworkingUnion) asAny() any {
	if !param.IsOmitted(u.OfUnrestricted) {
		return u.OfUnrestricted
	} else if !param.IsOmitted(u.OfLimited) {
		return u.OfLimited
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaCloudConfigParamsNetworkingUnion) GetAllowMCPServers() *bool {
	if vt := u.OfLimited; vt != nil && vt.AllowMCPServers.Valid() {
		return &vt.AllowMCPServers.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaCloudConfigParamsNetworkingUnion) GetAllowPackageManagers() *bool {
	if vt := u.OfLimited; vt != nil && vt.AllowPackageManagers.Valid() {
		return &vt.AllowPackageManagers.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaCloudConfigParamsNetworkingUnion) GetAllowedHosts() []string {
	if vt := u.OfLimited; vt != nil {
		return vt.AllowedHosts
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaCloudConfigParamsNetworkingUnion) GetType() *string {
	if vt := u.OfUnrestricted; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfLimited; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaCloudConfigParamsNetworkingUnion](
		"type",
		apijson.Discriminator[BetaUnrestrictedNetworkParam]("unrestricted"),
		apijson.Discriminator[BetaLimitedNetworkParams]("limited"),
	)
}

// Unified Environment resource for both cloud and BYOC environments.
type BetaEnvironment struct {
	// Environment identifier (e.g., 'env\_...')
	ID string `json:"id" api:"required"`
	// RFC 3339 timestamp when environment was archived, or null if not archived
	ArchivedAt string `json:"archived_at" api:"required"`
	// `cloud` environment configuration.
	Config BetaCloudConfig `json:"config" api:"required"`
	// RFC 3339 timestamp when environment was created
	CreatedAt string `json:"created_at" api:"required"`
	// User-provided description for the environment
	Description string `json:"description" api:"required"`
	// User-provided metadata key-value pairs
	Metadata map[string]string `json:"metadata" api:"required"`
	// Human-readable name for the environment
	Name string `json:"name" api:"required"`
	// The type of object (always 'environment')
	Type constant.Environment `json:"type" default:"environment"`
	// RFC 3339 timestamp when environment was last updated
	UpdatedAt string `json:"updated_at" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ArchivedAt  respjson.Field
		Config      respjson.Field
		CreatedAt   respjson.Field
		Description respjson.Field
		Metadata    respjson.Field
		Name        respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaEnvironment) RawJSON() string { return r.JSON.raw }
func (r *BetaEnvironment) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Response after deleting an environment.
type BetaEnvironmentDeleteResponse struct {
	// Environment identifier
	ID string `json:"id" api:"required"`
	// The type of response
	Type constant.EnvironmentDeleted `json:"type" default:"environment_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaEnvironmentDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaEnvironmentDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Limited network access.
type BetaLimitedNetwork struct {
	// Permits outbound access to MCP server endpoints configured on the agent, beyond
	// those listed in the `allowed_hosts` array.
	AllowMCPServers bool `json:"allow_mcp_servers" api:"required"`
	// Permits outbound access to public package registries (PyPI, npm, etc.) beyond
	// those listed in the `allowed_hosts` array.
	AllowPackageManagers bool `json:"allow_package_managers" api:"required"`
	// Specifies domains the container can reach.
	AllowedHosts []string `json:"allowed_hosts" api:"required"`
	// Network policy type
	Type constant.Limited `json:"type" default:"limited"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AllowMCPServers      respjson.Field
		AllowPackageManagers respjson.Field
		AllowedHosts         respjson.Field
		Type                 respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaLimitedNetwork) RawJSON() string { return r.JSON.raw }
func (r *BetaLimitedNetwork) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Limited network request params.
//
// Fields default to null; on update, omitted fields preserve the existing value.
//
// The property Type is required.
type BetaLimitedNetworkParams struct {
	// Permits outbound access to MCP server endpoints configured on the agent, beyond
	// those listed in the `allowed_hosts` array. Defaults to `false`.
	AllowMCPServers param.Opt[bool] `json:"allow_mcp_servers,omitzero"`
	// Permits outbound access to public package registries (PyPI, npm, etc.) beyond
	// those listed in the `allowed_hosts` array. Defaults to `false`.
	AllowPackageManagers param.Opt[bool] `json:"allow_package_managers,omitzero"`
	// Specifies domains the container can reach.
	AllowedHosts []string `json:"allowed_hosts,omitzero"`
	// Network policy type
	//
	// This field can be elided, and will marshal its zero value as "limited".
	Type constant.Limited `json:"type" default:"limited"`
	paramObj
}

func (r BetaLimitedNetworkParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaLimitedNetworkParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaLimitedNetworkParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Packages (and their versions) available in this environment.
type BetaPackages struct {
	// Ubuntu/Debian packages to install
	Apt []string `json:"apt" api:"required"`
	// Rust packages to install
	Cargo []string `json:"cargo" api:"required"`
	// Ruby packages to install
	Gem []string `json:"gem" api:"required"`
	// Go packages to install
	Go []string `json:"go" api:"required"`
	// Node.js packages to install
	Npm []string `json:"npm" api:"required"`
	// Python packages to install
	Pip []string `json:"pip" api:"required"`
	// Package configuration type
	//
	// Any of "packages".
	Type BetaPackagesType `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Apt         respjson.Field
		Cargo       respjson.Field
		Gem         respjson.Field
		Go          respjson.Field
		Npm         respjson.Field
		Pip         respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaPackages) RawJSON() string { return r.JSON.raw }
func (r *BetaPackages) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Package configuration type
type BetaPackagesType string

const (
	BetaPackagesTypePackages BetaPackagesType = "packages"
)

// Specify packages (and optionally their versions) available in this environment.
//
// When versioning, use the version semantics relevant for the package manager,
// e.g. for `pip` use `package==1.0.0`. You are responsible for validating the
// package and version exist. Unversioned installs the latest.
type BetaPackagesParams struct {
	// Ubuntu/Debian packages to install
	Apt []string `json:"apt,omitzero"`
	// Rust packages to install
	Cargo []string `json:"cargo,omitzero"`
	// Ruby packages to install
	Gem []string `json:"gem,omitzero"`
	// Go packages to install
	Go []string `json:"go,omitzero"`
	// Node.js packages to install
	Npm []string `json:"npm,omitzero"`
	// Python packages to install
	Pip []string `json:"pip,omitzero"`
	// Package configuration type
	//
	// Any of "packages".
	Type BetaPackagesParamsType `json:"type,omitzero"`
	paramObj
}

func (r BetaPackagesParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaPackagesParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaPackagesParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Package configuration type
type BetaPackagesParamsType string

const (
	BetaPackagesParamsTypePackages BetaPackagesParamsType = "packages"
)

// Unrestricted network access.
type BetaUnrestrictedNetwork struct {
	// Network policy type
	Type constant.Unrestricted `json:"type" default:"unrestricted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaUnrestrictedNetwork) RawJSON() string { return r.JSON.raw }
func (r *BetaUnrestrictedNetwork) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaUnrestrictedNetwork to a BetaUnrestrictedNetworkParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaUnrestrictedNetworkParam.Overrides()
func (r BetaUnrestrictedNetwork) ToParam() BetaUnrestrictedNetworkParam {
	return param.Override[BetaUnrestrictedNetworkParam](json.RawMessage(r.RawJSON()))
}

func NewBetaUnrestrictedNetworkParam() BetaUnrestrictedNetworkParam {
	return BetaUnrestrictedNetworkParam{
		Type: "unrestricted",
	}
}

// Unrestricted network access.
//
// This struct has a constant value, construct it with
// [NewBetaUnrestrictedNetworkParam].
type BetaUnrestrictedNetworkParam struct {
	// Network policy type
	Type constant.Unrestricted `json:"type" default:"unrestricted"`
	paramObj
}

func (r BetaUnrestrictedNetworkParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaUnrestrictedNetworkParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaUnrestrictedNetworkParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaEnvironmentNewParams struct {
	// Human-readable name for the environment
	Name string `json:"name" api:"required"`
	// Optional description of the environment
	Description param.Opt[string] `json:"description,omitzero"`
	// Request params for `cloud` environment configuration.
	//
	// Fields default to null; on update, omitted fields preserve the existing value.
	Config BetaCloudConfigParams `json:"config,omitzero"`
	// User-provided metadata key-value pairs
	Metadata map[string]string `json:"metadata,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaEnvironmentNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaEnvironmentNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaEnvironmentNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaEnvironmentGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaEnvironmentUpdateParams struct {
	// Updated description of the environment
	Description param.Opt[string] `json:"description,omitzero"`
	// Updated name for the environment
	Name param.Opt[string] `json:"name,omitzero"`
	// Request params for `cloud` environment configuration.
	//
	// Fields default to null; on update, omitted fields preserve the existing value.
	Config BetaCloudConfigParams `json:"config,omitzero"`
	// User-provided metadata key-value pairs. Set a value to null or empty string to
	// delete the key.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaEnvironmentUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaEnvironmentUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaEnvironmentUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaEnvironmentListParams struct {
	// Opaque cursor from previous response for pagination. Pass the `next_page` value
	// from the previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Include archived environments in the response
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum number of environments to return
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaEnvironmentListParams]'s query parameters as
// `url.Values`.
func (r BetaEnvironmentListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaEnvironmentDeleteParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaEnvironmentArchiveParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
