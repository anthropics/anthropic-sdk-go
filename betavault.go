// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"context"
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

// BetaVaultService contains methods and other services that help with interacting
// with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaVaultService] method instead.
type BetaVaultService struct {
	Options     []option.RequestOption
	Credentials BetaVaultCredentialService
}

// NewBetaVaultService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaVaultService(opts ...option.RequestOption) (r BetaVaultService) {
	r = BetaVaultService{}
	r.Options = opts
	r.Credentials = NewBetaVaultCredentialService(opts...)
	return
}

// Create Vault
func (r *BetaVaultService) New(ctx context.Context, params BetaVaultNewParams, opts ...option.RequestOption) (res *BetaManagedAgentsVault, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	path := "v1/vaults?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get Vault
func (r *BetaVaultService) Get(ctx context.Context, vaultID string, query BetaVaultGetParams, opts ...option.RequestOption) (res *BetaManagedAgentsVault, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/vaults/%s?beta=true", vaultID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Vault
func (r *BetaVaultService) Update(ctx context.Context, vaultID string, params BetaVaultUpdateParams, opts ...option.RequestOption) (res *BetaManagedAgentsVault, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/vaults/%s?beta=true", vaultID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List Vaults
func (r *BetaVaultService) List(ctx context.Context, params BetaVaultListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaManagedAgentsVault], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/vaults?beta=true"
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

// List Vaults
func (r *BetaVaultService) ListAutoPaging(ctx context.Context, params BetaVaultListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsVault] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Delete Vault
func (r *BetaVaultService) Delete(ctx context.Context, vaultID string, body BetaVaultDeleteParams, opts ...option.RequestOption) (res *BetaManagedAgentsDeletedVault, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/vaults/%s?beta=true", vaultID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Archive Vault
func (r *BetaVaultService) Archive(ctx context.Context, vaultID string, body BetaVaultArchiveParams, opts ...option.RequestOption) (res *BetaManagedAgentsVault, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/vaults/%s/archive?beta=true", vaultID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Confirmation of a deleted vault.
type BetaManagedAgentsDeletedVault struct {
	// Unique identifier of the deleted vault.
	ID string `json:"id" api:"required"`
	// Any of "vault_deleted".
	Type BetaManagedAgentsDeletedVaultType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsDeletedVault) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsDeletedVault) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsDeletedVaultType string

const (
	BetaManagedAgentsDeletedVaultTypeVaultDeleted BetaManagedAgentsDeletedVaultType = "vault_deleted"
)

// A vault that stores credentials for use by agents during sessions.
type BetaManagedAgentsVault struct {
	// Unique identifier for the vault.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Human-readable name for the vault.
	DisplayName string `json:"display_name" api:"required"`
	// Arbitrary key-value metadata attached to the vault.
	Metadata map[string]string `json:"metadata" api:"required"`
	// Any of "vault".
	Type BetaManagedAgentsVaultType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ArchivedAt  respjson.Field
		CreatedAt   respjson.Field
		DisplayName respjson.Field
		Metadata    respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsVault) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsVault) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsVaultType string

const (
	BetaManagedAgentsVaultTypeVault BetaManagedAgentsVaultType = "vault"
)

type BetaVaultNewParams struct {
	// Human-readable name for the vault. 1-255 characters.
	DisplayName string `json:"display_name" api:"required"`
	// Arbitrary key-value metadata to attach to the vault. Maximum 16 pairs, keys up
	// to 64 chars, values up to 512 chars.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaVaultNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaVaultNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaVaultNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaVaultGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaVaultUpdateParams struct {
	// Updated human-readable name for the vault. 1-255 characters.
	DisplayName param.Opt[string] `json:"display_name,omitzero"`
	// Metadata patch. Set a key to a string to upsert it, or to null to delete it.
	// Omitted keys are preserved.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaVaultUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaVaultUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaVaultUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaVaultListParams struct {
	// Whether to include archived vaults in the results.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum number of vaults to return per page. Defaults to 20, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination token from a previous `list_vaults` response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaVaultListParams]'s query parameters as `url.Values`.
func (r BetaVaultListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaVaultDeleteParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaVaultArchiveParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
