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
)

// BetaVaultCredentialService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaVaultCredentialService] method instead.
type BetaVaultCredentialService struct {
	Options []option.RequestOption
}

// NewBetaVaultCredentialService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaVaultCredentialService(opts ...option.RequestOption) (r BetaVaultCredentialService) {
	r = BetaVaultCredentialService{}
	r.Options = opts
	return
}

// Create Credential
func (r *BetaVaultCredentialService) New(ctx context.Context, vaultID string, params BetaVaultCredentialNewParams, opts ...option.RequestOption) (res *BetaManagedAgentsCredential, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/vaults/%s/credentials?beta=true", vaultID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get Credential
func (r *BetaVaultCredentialService) Get(ctx context.Context, credentialID string, params BetaVaultCredentialGetParams, opts ...option.RequestOption) (res *BetaManagedAgentsCredential, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if params.VaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	if credentialID == "" {
		err = errors.New("missing required credential_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/vaults/%s/credentials/%s?beta=true", params.VaultID, credentialID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Credential
func (r *BetaVaultCredentialService) Update(ctx context.Context, credentialID string, params BetaVaultCredentialUpdateParams, opts ...option.RequestOption) (res *BetaManagedAgentsCredential, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if params.VaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	if credentialID == "" {
		err = errors.New("missing required credential_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/vaults/%s/credentials/%s?beta=true", params.VaultID, credentialID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List Credentials
func (r *BetaVaultCredentialService) List(ctx context.Context, vaultID string, params BetaVaultCredentialListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaManagedAgentsCredential], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01"), option.WithResponseInto(&raw)}, opts...)
	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/vaults/%s/credentials?beta=true", vaultID)
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

// List Credentials
func (r *BetaVaultCredentialService) ListAutoPaging(ctx context.Context, vaultID string, params BetaVaultCredentialListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsCredential] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, vaultID, params, opts...))
}

// Delete Credential
func (r *BetaVaultCredentialService) Delete(ctx context.Context, credentialID string, params BetaVaultCredentialDeleteParams, opts ...option.RequestOption) (res *BetaManagedAgentsDeletedCredential, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if params.VaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	if credentialID == "" {
		err = errors.New("missing required credential_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/vaults/%s/credentials/%s?beta=true", params.VaultID, credentialID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Archive Credential
func (r *BetaVaultCredentialService) Archive(ctx context.Context, credentialID string, params BetaVaultCredentialArchiveParams, opts ...option.RequestOption) (res *BetaManagedAgentsCredential, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if params.VaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	if credentialID == "" {
		err = errors.New("missing required credential_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/vaults/%s/credentials/%s/archive?beta=true", params.VaultID, credentialID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// A credential stored in a vault. Sensitive fields are never returned in
// responses.
type BetaManagedAgentsCredential struct {
	// Unique identifier for the credential.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// Authentication details for a credential.
	Auth BetaManagedAgentsCredentialAuthUnion `json:"auth" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Arbitrary key-value metadata attached to the credential.
	Metadata map[string]string `json:"metadata" api:"required"`
	// Any of "vault_credential".
	Type BetaManagedAgentsCredentialType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// Identifier of the vault this credential belongs to.
	VaultID string `json:"vault_id" api:"required"`
	// Human-readable name for the credential.
	DisplayName string `json:"display_name" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ArchivedAt  respjson.Field
		Auth        respjson.Field
		CreatedAt   respjson.Field
		Metadata    respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		VaultID     respjson.Field
		DisplayName respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsCredential) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsCredential) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsCredentialAuthUnion contains all possible properties and values
// from [BetaManagedAgentsMCPOAuthAuthResponse],
// [BetaManagedAgentsStaticBearerAuthResponse].
//
// Use the [BetaManagedAgentsCredentialAuthUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsCredentialAuthUnion struct {
	MCPServerURL string `json:"mcp_server_url"`
	// Any of "mcp_oauth", "static_bearer".
	Type string `json:"type"`
	// This field is from variant [BetaManagedAgentsMCPOAuthAuthResponse].
	ExpiresAt time.Time `json:"expires_at"`
	// This field is from variant [BetaManagedAgentsMCPOAuthAuthResponse].
	Refresh BetaManagedAgentsMCPOAuthRefreshResponse `json:"refresh"`
	JSON    struct {
		MCPServerURL respjson.Field
		Type         respjson.Field
		ExpiresAt    respjson.Field
		Refresh      respjson.Field
		raw          string
	} `json:"-"`
}

// anyBetaManagedAgentsCredentialAuth is implemented by each variant of
// [BetaManagedAgentsCredentialAuthUnion] to add type safety for the return type of
// [BetaManagedAgentsCredentialAuthUnion.AsAny]
type anyBetaManagedAgentsCredentialAuth interface {
	implBetaManagedAgentsCredentialAuthUnion()
}

func (BetaManagedAgentsMCPOAuthAuthResponse) implBetaManagedAgentsCredentialAuthUnion()     {}
func (BetaManagedAgentsStaticBearerAuthResponse) implBetaManagedAgentsCredentialAuthUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsCredentialAuthUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsMCPOAuthAuthResponse:
//	case anthropic.BetaManagedAgentsStaticBearerAuthResponse:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsCredentialAuthUnion) AsAny() anyBetaManagedAgentsCredentialAuth {
	switch u.Type {
	case "mcp_oauth":
		return u.AsMCPOAuth()
	case "static_bearer":
		return u.AsStaticBearer()
	}
	return nil
}

func (u BetaManagedAgentsCredentialAuthUnion) AsMCPOAuth() (v BetaManagedAgentsMCPOAuthAuthResponse) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsCredentialAuthUnion) AsStaticBearer() (v BetaManagedAgentsStaticBearerAuthResponse) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsCredentialAuthUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsCredentialAuthUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsCredentialType string

const (
	BetaManagedAgentsCredentialTypeVaultCredential BetaManagedAgentsCredentialType = "vault_credential"
)

// Confirmation of a deleted credential.
type BetaManagedAgentsDeletedCredential struct {
	// Unique identifier of the deleted credential.
	ID string `json:"id" api:"required"`
	// Any of "vault_credential_deleted".
	Type BetaManagedAgentsDeletedCredentialType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsDeletedCredential) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsDeletedCredential) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsDeletedCredentialType string

const (
	BetaManagedAgentsDeletedCredentialTypeVaultCredentialDeleted BetaManagedAgentsDeletedCredentialType = "vault_credential_deleted"
)

// OAuth credential details for an MCP server.
type BetaManagedAgentsMCPOAuthAuthResponse struct {
	// URL of the MCP server this credential authenticates against.
	MCPServerURL string `json:"mcp_server_url" api:"required"`
	// Any of "mcp_oauth".
	Type BetaManagedAgentsMCPOAuthAuthResponseType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	ExpiresAt time.Time `json:"expires_at" api:"nullable" format:"date-time"`
	// OAuth refresh token configuration returned in credential responses.
	Refresh BetaManagedAgentsMCPOAuthRefreshResponse `json:"refresh" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MCPServerURL respjson.Field
		Type         respjson.Field
		ExpiresAt    respjson.Field
		Refresh      respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMCPOAuthAuthResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsMCPOAuthAuthResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMCPOAuthAuthResponseType string

const (
	BetaManagedAgentsMCPOAuthAuthResponseTypeMCPOAuth BetaManagedAgentsMCPOAuthAuthResponseType = "mcp_oauth"
)

// Parameters for creating an MCP OAuth credential.
//
// The properties AccessToken, MCPServerURL, Type are required.
type BetaManagedAgentsMCPOAuthCreateParams struct {
	// OAuth access token.
	AccessToken string `json:"access_token" api:"required"`
	// URL of the MCP server this credential authenticates against.
	MCPServerURL string `json:"mcp_server_url" api:"required"`
	// Any of "mcp_oauth".
	Type BetaManagedAgentsMCPOAuthCreateParamsType `json:"type,omitzero" api:"required"`
	// A timestamp in RFC 3339 format
	ExpiresAt param.Opt[time.Time] `json:"expires_at,omitzero" format:"date-time"`
	// OAuth refresh token parameters for creating a credential with refresh support.
	Refresh BetaManagedAgentsMCPOAuthRefreshParams `json:"refresh,omitzero"`
	paramObj
}

func (r BetaManagedAgentsMCPOAuthCreateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsMCPOAuthCreateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsMCPOAuthCreateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMCPOAuthCreateParamsType string

const (
	BetaManagedAgentsMCPOAuthCreateParamsTypeMCPOAuth BetaManagedAgentsMCPOAuthCreateParamsType = "mcp_oauth"
)

// OAuth refresh token parameters for creating a credential with refresh support.
//
// The properties ClientID, RefreshToken, TokenEndpoint, TokenEndpointAuth are
// required.
type BetaManagedAgentsMCPOAuthRefreshParams struct {
	// OAuth client ID.
	ClientID string `json:"client_id" api:"required"`
	// OAuth refresh token.
	RefreshToken string `json:"refresh_token" api:"required"`
	// Token endpoint URL used to refresh the access token.
	TokenEndpoint string `json:"token_endpoint" api:"required"`
	// Token endpoint requires no client authentication.
	TokenEndpointAuth BetaManagedAgentsMCPOAuthRefreshParamsTokenEndpointAuthUnion `json:"token_endpoint_auth,omitzero" api:"required"`
	// OAuth resource indicator.
	Resource param.Opt[string] `json:"resource,omitzero"`
	// OAuth scope for the refresh request.
	Scope param.Opt[string] `json:"scope,omitzero"`
	paramObj
}

func (r BetaManagedAgentsMCPOAuthRefreshParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsMCPOAuthRefreshParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsMCPOAuthRefreshParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsMCPOAuthRefreshParamsTokenEndpointAuthUnion struct {
	OfNone              *BetaManagedAgentsTokenEndpointAuthNoneParam  `json:",omitzero,inline"`
	OfClientSecretBasic *BetaManagedAgentsTokenEndpointAuthBasicParam `json:",omitzero,inline"`
	OfClientSecretPost  *BetaManagedAgentsTokenEndpointAuthPostParam  `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsMCPOAuthRefreshParamsTokenEndpointAuthUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfNone, u.OfClientSecretBasic, u.OfClientSecretPost)
}
func (u *BetaManagedAgentsMCPOAuthRefreshParamsTokenEndpointAuthUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsMCPOAuthRefreshParamsTokenEndpointAuthUnion) asAny() any {
	if !param.IsOmitted(u.OfNone) {
		return u.OfNone
	} else if !param.IsOmitted(u.OfClientSecretBasic) {
		return u.OfClientSecretBasic
	} else if !param.IsOmitted(u.OfClientSecretPost) {
		return u.OfClientSecretPost
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsMCPOAuthRefreshParamsTokenEndpointAuthUnion) GetType() *string {
	if vt := u.OfNone; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfClientSecretBasic; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfClientSecretPost; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsMCPOAuthRefreshParamsTokenEndpointAuthUnion) GetClientSecret() *string {
	if vt := u.OfClientSecretBasic; vt != nil {
		return (*string)(&vt.ClientSecret)
	} else if vt := u.OfClientSecretPost; vt != nil {
		return (*string)(&vt.ClientSecret)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsMCPOAuthRefreshParamsTokenEndpointAuthUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsTokenEndpointAuthNoneParam]("none"),
		apijson.Discriminator[BetaManagedAgentsTokenEndpointAuthBasicParam]("client_secret_basic"),
		apijson.Discriminator[BetaManagedAgentsTokenEndpointAuthPostParam]("client_secret_post"),
	)
}

// OAuth refresh token configuration returned in credential responses.
type BetaManagedAgentsMCPOAuthRefreshResponse struct {
	// OAuth client ID.
	ClientID string `json:"client_id" api:"required"`
	// Token endpoint URL used to refresh the access token.
	TokenEndpoint string `json:"token_endpoint" api:"required"`
	// Token endpoint requires no client authentication.
	TokenEndpointAuth BetaManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion `json:"token_endpoint_auth" api:"required"`
	// OAuth resource indicator.
	Resource string `json:"resource" api:"nullable"`
	// OAuth scope for the refresh request.
	Scope string `json:"scope" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ClientID          respjson.Field
		TokenEndpoint     respjson.Field
		TokenEndpointAuth respjson.Field
		Resource          respjson.Field
		Scope             respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMCPOAuthRefreshResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsMCPOAuthRefreshResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion contains all
// possible properties and values from
// [BetaManagedAgentsTokenEndpointAuthNoneResponse],
// [BetaManagedAgentsTokenEndpointAuthBasicResponse],
// [BetaManagedAgentsTokenEndpointAuthPostResponse].
//
// Use the [BetaManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion struct {
	// Any of "none", "client_secret_basic", "client_secret_post".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuth is implemented by
// each variant of [BetaManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion]
// to add type safety for the return type of
// [BetaManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion.AsAny]
type anyBetaManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuth interface {
	implBetaManagedAgentsMcpoAuthRefreshResponseTokenEndpointAuthUnion()
}

func (BetaManagedAgentsTokenEndpointAuthNoneResponse) implBetaManagedAgentsMcpoAuthRefreshResponseTokenEndpointAuthUnion() {
}
func (BetaManagedAgentsTokenEndpointAuthBasicResponse) implBetaManagedAgentsMcpoAuthRefreshResponseTokenEndpointAuthUnion() {
}
func (BetaManagedAgentsTokenEndpointAuthPostResponse) implBetaManagedAgentsMcpoAuthRefreshResponseTokenEndpointAuthUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsTokenEndpointAuthNoneResponse:
//	case anthropic.BetaManagedAgentsTokenEndpointAuthBasicResponse:
//	case anthropic.BetaManagedAgentsTokenEndpointAuthPostResponse:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion) AsAny() anyBetaManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuth {
	switch u.Type {
	case "none":
		return u.AsNone()
	case "client_secret_basic":
		return u.AsClientSecretBasic()
	case "client_secret_post":
		return u.AsClientSecretPost()
	}
	return nil
}

func (u BetaManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion) AsNone() (v BetaManagedAgentsTokenEndpointAuthNoneResponse) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion) AsClientSecretBasic() (v BetaManagedAgentsTokenEndpointAuthBasicResponse) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion) AsClientSecretPost() (v BetaManagedAgentsTokenEndpointAuthPostResponse) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *BetaManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Parameters for updating OAuth refresh token configuration.
type BetaManagedAgentsMCPOAuthRefreshUpdateParams struct {
	// Updated OAuth refresh token.
	RefreshToken param.Opt[string] `json:"refresh_token,omitzero"`
	// Updated OAuth scope for the refresh request.
	Scope param.Opt[string] `json:"scope,omitzero"`
	// Updated HTTP Basic authentication parameters for the token endpoint.
	TokenEndpointAuth BetaManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion `json:"token_endpoint_auth,omitzero"`
	paramObj
}

func (r BetaManagedAgentsMCPOAuthRefreshUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsMCPOAuthRefreshUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsMCPOAuthRefreshUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion struct {
	OfClientSecretBasic *BetaManagedAgentsTokenEndpointAuthBasicUpdateParam `json:",omitzero,inline"`
	OfClientSecretPost  *BetaManagedAgentsTokenEndpointAuthPostUpdateParam  `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfClientSecretBasic, u.OfClientSecretPost)
}
func (u *BetaManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion) asAny() any {
	if !param.IsOmitted(u.OfClientSecretBasic) {
		return u.OfClientSecretBasic
	} else if !param.IsOmitted(u.OfClientSecretPost) {
		return u.OfClientSecretPost
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion) GetType() *string {
	if vt := u.OfClientSecretBasic; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfClientSecretPost; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion) GetClientSecret() *string {
	if vt := u.OfClientSecretBasic; vt != nil && vt.ClientSecret.Valid() {
		return &vt.ClientSecret.Value
	} else if vt := u.OfClientSecretPost; vt != nil && vt.ClientSecret.Valid() {
		return &vt.ClientSecret.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsTokenEndpointAuthBasicUpdateParam]("client_secret_basic"),
		apijson.Discriminator[BetaManagedAgentsTokenEndpointAuthPostUpdateParam]("client_secret_post"),
	)
}

// Parameters for updating an MCP OAuth credential. The `mcp_server_url` is
// immutable.
//
// The property Type is required.
type BetaManagedAgentsMCPOAuthUpdateParams struct {
	// Any of "mcp_oauth".
	Type BetaManagedAgentsMCPOAuthUpdateParamsType `json:"type,omitzero" api:"required"`
	// Updated OAuth access token.
	AccessToken param.Opt[string] `json:"access_token,omitzero"`
	// A timestamp in RFC 3339 format
	ExpiresAt param.Opt[time.Time] `json:"expires_at,omitzero" format:"date-time"`
	// Parameters for updating OAuth refresh token configuration.
	Refresh BetaManagedAgentsMCPOAuthRefreshUpdateParams `json:"refresh,omitzero"`
	paramObj
}

func (r BetaManagedAgentsMCPOAuthUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsMCPOAuthUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsMCPOAuthUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMCPOAuthUpdateParamsType string

const (
	BetaManagedAgentsMCPOAuthUpdateParamsTypeMCPOAuth BetaManagedAgentsMCPOAuthUpdateParamsType = "mcp_oauth"
)

// Static bearer token credential details for an MCP server.
type BetaManagedAgentsStaticBearerAuthResponse struct {
	// URL of the MCP server this credential authenticates against.
	MCPServerURL string `json:"mcp_server_url" api:"required"`
	// Any of "static_bearer".
	Type BetaManagedAgentsStaticBearerAuthResponseType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MCPServerURL respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsStaticBearerAuthResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsStaticBearerAuthResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsStaticBearerAuthResponseType string

const (
	BetaManagedAgentsStaticBearerAuthResponseTypeStaticBearer BetaManagedAgentsStaticBearerAuthResponseType = "static_bearer"
)

// Parameters for creating a static bearer token credential.
//
// The properties Token, MCPServerURL, Type are required.
type BetaManagedAgentsStaticBearerCreateParams struct {
	// Static bearer token value.
	Token string `json:"token" api:"required"`
	// URL of the MCP server this credential authenticates against.
	MCPServerURL string `json:"mcp_server_url" api:"required"`
	// Any of "static_bearer".
	Type BetaManagedAgentsStaticBearerCreateParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsStaticBearerCreateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsStaticBearerCreateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsStaticBearerCreateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsStaticBearerCreateParamsType string

const (
	BetaManagedAgentsStaticBearerCreateParamsTypeStaticBearer BetaManagedAgentsStaticBearerCreateParamsType = "static_bearer"
)

// Parameters for updating a static bearer token credential. The `mcp_server_url`
// is immutable.
//
// The property Type is required.
type BetaManagedAgentsStaticBearerUpdateParams struct {
	// Any of "static_bearer".
	Type BetaManagedAgentsStaticBearerUpdateParamsType `json:"type,omitzero" api:"required"`
	// Updated static bearer token value.
	Token param.Opt[string] `json:"token,omitzero"`
	paramObj
}

func (r BetaManagedAgentsStaticBearerUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsStaticBearerUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsStaticBearerUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsStaticBearerUpdateParamsType string

const (
	BetaManagedAgentsStaticBearerUpdateParamsTypeStaticBearer BetaManagedAgentsStaticBearerUpdateParamsType = "static_bearer"
)

// Token endpoint uses HTTP Basic authentication with client credentials.
//
// The properties ClientSecret, Type are required.
type BetaManagedAgentsTokenEndpointAuthBasicParam struct {
	// OAuth client secret.
	ClientSecret string `json:"client_secret" api:"required"`
	// Any of "client_secret_basic".
	Type BetaManagedAgentsTokenEndpointAuthBasicParamType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsTokenEndpointAuthBasicParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsTokenEndpointAuthBasicParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsTokenEndpointAuthBasicParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsTokenEndpointAuthBasicParamType string

const (
	BetaManagedAgentsTokenEndpointAuthBasicParamTypeClientSecretBasic BetaManagedAgentsTokenEndpointAuthBasicParamType = "client_secret_basic"
)

// Token endpoint uses HTTP Basic authentication with client credentials.
type BetaManagedAgentsTokenEndpointAuthBasicResponse struct {
	// Any of "client_secret_basic".
	Type BetaManagedAgentsTokenEndpointAuthBasicResponseType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsTokenEndpointAuthBasicResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsTokenEndpointAuthBasicResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsTokenEndpointAuthBasicResponseType string

const (
	BetaManagedAgentsTokenEndpointAuthBasicResponseTypeClientSecretBasic BetaManagedAgentsTokenEndpointAuthBasicResponseType = "client_secret_basic"
)

// Updated HTTP Basic authentication parameters for the token endpoint.
//
// The property Type is required.
type BetaManagedAgentsTokenEndpointAuthBasicUpdateParam struct {
	// Any of "client_secret_basic".
	Type BetaManagedAgentsTokenEndpointAuthBasicUpdateParamType `json:"type,omitzero" api:"required"`
	// Updated OAuth client secret.
	ClientSecret param.Opt[string] `json:"client_secret,omitzero"`
	paramObj
}

func (r BetaManagedAgentsTokenEndpointAuthBasicUpdateParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsTokenEndpointAuthBasicUpdateParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsTokenEndpointAuthBasicUpdateParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsTokenEndpointAuthBasicUpdateParamType string

const (
	BetaManagedAgentsTokenEndpointAuthBasicUpdateParamTypeClientSecretBasic BetaManagedAgentsTokenEndpointAuthBasicUpdateParamType = "client_secret_basic"
)

// Token endpoint requires no client authentication.
//
// The property Type is required.
type BetaManagedAgentsTokenEndpointAuthNoneParam struct {
	// Any of "none".
	Type BetaManagedAgentsTokenEndpointAuthNoneParamType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsTokenEndpointAuthNoneParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsTokenEndpointAuthNoneParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsTokenEndpointAuthNoneParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsTokenEndpointAuthNoneParamType string

const (
	BetaManagedAgentsTokenEndpointAuthNoneParamTypeNone BetaManagedAgentsTokenEndpointAuthNoneParamType = "none"
)

// Token endpoint requires no client authentication.
type BetaManagedAgentsTokenEndpointAuthNoneResponse struct {
	// Any of "none".
	Type BetaManagedAgentsTokenEndpointAuthNoneResponseType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsTokenEndpointAuthNoneResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsTokenEndpointAuthNoneResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsTokenEndpointAuthNoneResponseType string

const (
	BetaManagedAgentsTokenEndpointAuthNoneResponseTypeNone BetaManagedAgentsTokenEndpointAuthNoneResponseType = "none"
)

// Token endpoint uses POST body authentication with client credentials.
//
// The properties ClientSecret, Type are required.
type BetaManagedAgentsTokenEndpointAuthPostParam struct {
	// OAuth client secret.
	ClientSecret string `json:"client_secret" api:"required"`
	// Any of "client_secret_post".
	Type BetaManagedAgentsTokenEndpointAuthPostParamType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsTokenEndpointAuthPostParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsTokenEndpointAuthPostParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsTokenEndpointAuthPostParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsTokenEndpointAuthPostParamType string

const (
	BetaManagedAgentsTokenEndpointAuthPostParamTypeClientSecretPost BetaManagedAgentsTokenEndpointAuthPostParamType = "client_secret_post"
)

// Token endpoint uses POST body authentication with client credentials.
type BetaManagedAgentsTokenEndpointAuthPostResponse struct {
	// Any of "client_secret_post".
	Type BetaManagedAgentsTokenEndpointAuthPostResponseType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsTokenEndpointAuthPostResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsTokenEndpointAuthPostResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsTokenEndpointAuthPostResponseType string

const (
	BetaManagedAgentsTokenEndpointAuthPostResponseTypeClientSecretPost BetaManagedAgentsTokenEndpointAuthPostResponseType = "client_secret_post"
)

// Updated POST body authentication parameters for the token endpoint.
//
// The property Type is required.
type BetaManagedAgentsTokenEndpointAuthPostUpdateParam struct {
	// Any of "client_secret_post".
	Type BetaManagedAgentsTokenEndpointAuthPostUpdateParamType `json:"type,omitzero" api:"required"`
	// Updated OAuth client secret.
	ClientSecret param.Opt[string] `json:"client_secret,omitzero"`
	paramObj
}

func (r BetaManagedAgentsTokenEndpointAuthPostUpdateParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsTokenEndpointAuthPostUpdateParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsTokenEndpointAuthPostUpdateParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsTokenEndpointAuthPostUpdateParamType string

const (
	BetaManagedAgentsTokenEndpointAuthPostUpdateParamTypeClientSecretPost BetaManagedAgentsTokenEndpointAuthPostUpdateParamType = "client_secret_post"
)

type BetaVaultCredentialNewParams struct {
	// Authentication details for creating a credential.
	Auth BetaVaultCredentialNewParamsAuthUnion `json:"auth,omitzero" api:"required"`
	// Human-readable name for the credential. Up to 255 characters.
	DisplayName param.Opt[string] `json:"display_name,omitzero"`
	// Arbitrary key-value metadata to attach to the credential. Maximum 16 pairs, keys
	// up to 64 chars, values up to 512 chars.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaVaultCredentialNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaVaultCredentialNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaVaultCredentialNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaVaultCredentialNewParamsAuthUnion struct {
	OfMCPOAuth     *BetaManagedAgentsMCPOAuthCreateParams     `json:",omitzero,inline"`
	OfStaticBearer *BetaManagedAgentsStaticBearerCreateParams `json:",omitzero,inline"`
	paramUnion
}

func (u BetaVaultCredentialNewParamsAuthUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfMCPOAuth, u.OfStaticBearer)
}
func (u *BetaVaultCredentialNewParamsAuthUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaVaultCredentialNewParamsAuthUnion) asAny() any {
	if !param.IsOmitted(u.OfMCPOAuth) {
		return u.OfMCPOAuth
	} else if !param.IsOmitted(u.OfStaticBearer) {
		return u.OfStaticBearer
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaVaultCredentialNewParamsAuthUnion) GetAccessToken() *string {
	if vt := u.OfMCPOAuth; vt != nil {
		return &vt.AccessToken
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaVaultCredentialNewParamsAuthUnion) GetExpiresAt() *time.Time {
	if vt := u.OfMCPOAuth; vt != nil && vt.ExpiresAt.Valid() {
		return &vt.ExpiresAt.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaVaultCredentialNewParamsAuthUnion) GetRefresh() *BetaManagedAgentsMCPOAuthRefreshParams {
	if vt := u.OfMCPOAuth; vt != nil {
		return &vt.Refresh
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaVaultCredentialNewParamsAuthUnion) GetToken() *string {
	if vt := u.OfStaticBearer; vt != nil {
		return &vt.Token
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaVaultCredentialNewParamsAuthUnion) GetMCPServerURL() *string {
	if vt := u.OfMCPOAuth; vt != nil {
		return (*string)(&vt.MCPServerURL)
	} else if vt := u.OfStaticBearer; vt != nil {
		return (*string)(&vt.MCPServerURL)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaVaultCredentialNewParamsAuthUnion) GetType() *string {
	if vt := u.OfMCPOAuth; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfStaticBearer; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaVaultCredentialNewParamsAuthUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsMCPOAuthCreateParams]("mcp_oauth"),
		apijson.Discriminator[BetaManagedAgentsStaticBearerCreateParams]("static_bearer"),
	)
}

type BetaVaultCredentialGetParams struct {
	VaultID string `path:"vault_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaVaultCredentialUpdateParams struct {
	VaultID string `path:"vault_id" api:"required" json:"-"`
	// Updated human-readable name for the credential. 1-255 characters.
	DisplayName param.Opt[string] `json:"display_name,omitzero"`
	// Metadata patch. Set a key to a string to upsert it, or to null to delete it.
	// Omitted keys are preserved.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Updated authentication details for a credential.
	Auth BetaVaultCredentialUpdateParamsAuthUnion `json:"auth,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaVaultCredentialUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaVaultCredentialUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaVaultCredentialUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaVaultCredentialUpdateParamsAuthUnion struct {
	OfMCPOAuth     *BetaManagedAgentsMCPOAuthUpdateParams     `json:",omitzero,inline"`
	OfStaticBearer *BetaManagedAgentsStaticBearerUpdateParams `json:",omitzero,inline"`
	paramUnion
}

func (u BetaVaultCredentialUpdateParamsAuthUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfMCPOAuth, u.OfStaticBearer)
}
func (u *BetaVaultCredentialUpdateParamsAuthUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaVaultCredentialUpdateParamsAuthUnion) asAny() any {
	if !param.IsOmitted(u.OfMCPOAuth) {
		return u.OfMCPOAuth
	} else if !param.IsOmitted(u.OfStaticBearer) {
		return u.OfStaticBearer
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaVaultCredentialUpdateParamsAuthUnion) GetAccessToken() *string {
	if vt := u.OfMCPOAuth; vt != nil && vt.AccessToken.Valid() {
		return &vt.AccessToken.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaVaultCredentialUpdateParamsAuthUnion) GetExpiresAt() *time.Time {
	if vt := u.OfMCPOAuth; vt != nil && vt.ExpiresAt.Valid() {
		return &vt.ExpiresAt.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaVaultCredentialUpdateParamsAuthUnion) GetRefresh() *BetaManagedAgentsMCPOAuthRefreshUpdateParams {
	if vt := u.OfMCPOAuth; vt != nil {
		return &vt.Refresh
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaVaultCredentialUpdateParamsAuthUnion) GetToken() *string {
	if vt := u.OfStaticBearer; vt != nil && vt.Token.Valid() {
		return &vt.Token.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaVaultCredentialUpdateParamsAuthUnion) GetType() *string {
	if vt := u.OfMCPOAuth; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfStaticBearer; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaVaultCredentialUpdateParamsAuthUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsMCPOAuthUpdateParams]("mcp_oauth"),
		apijson.Discriminator[BetaManagedAgentsStaticBearerUpdateParams]("static_bearer"),
	)
}

type BetaVaultCredentialListParams struct {
	// Whether to include archived credentials in the results.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum number of credentials to return per page. Defaults to 20, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination token from a previous `list_credentials` response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaVaultCredentialListParams]'s query parameters as
// `url.Values`.
func (r BetaVaultCredentialListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaVaultCredentialDeleteParams struct {
	VaultID string `path:"vault_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaVaultCredentialArchiveParams struct {
	VaultID string `path:"vault_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
