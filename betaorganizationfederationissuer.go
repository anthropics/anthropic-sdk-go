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

// BetaOrganizationFederationIssuerService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationFederationIssuerService] method instead.
type BetaOrganizationFederationIssuerService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationFederationIssuerService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationFederationIssuerService(opts ...option.RequestOption) (r BetaOrganizationFederationIssuerService) {
	r = BetaOrganizationFederationIssuerService{}
	r.Options = opts
	return
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Register an OIDC issuer that Anthropic will trust for workload identity
// federation in your organization.
//
// The `jwks` field controls how the issuer's signing keys are obtained and takes
// one of three shapes selected by `type`: `discovery` (resolve keys through OIDC
// discovery), `explicit_url` (fetch keys from a fixed JWKS URL), or `inline`
// (provide a static key set). When `jwks.type` is `discovery` and no
// `discovery_base` is set, the issuer URL must be publicly reachable over HTTPS so
// Anthropic can fetch the discovery document; for `explicit_url` and `inline`
// modes the issuer URL is only matched as the JWT's `iss` claim and is not
// fetched.
func (r *BetaOrganizationFederationIssuerService) New(ctx context.Context, params BetaOrganizationFederationIssuerNewParams, opts ...option.RequestOption) (res *BetaFederationIssuer, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/federation_issuers?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Retrieve a federation issuer by its ID (`fdis_...`).
func (r *BetaOrganizationFederationIssuerService) Get(ctx context.Context, federationIssuerID string, query BetaOrganizationFederationIssuerGetParams, opts ...option.RequestOption) (res *BetaFederationIssuer, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	if federationIssuerID == "" {
		err = errors.New("missing required federation_issuer_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/federation_issuers/%s?beta=true", federationIssuerID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Partially update a federation issuer.
//
// Setting `jwks` replaces the full JWKS shape at once. Archived issuers cannot be
// updated; this returns 400. Create a new issuer instead.
//
// Updating an issuer that backs a rule with a scope outside `workspace:developer`
// or `workspace:inference` requires a Console session.
func (r *BetaOrganizationFederationIssuerService) Update(ctx context.Context, federationIssuerID string, params BetaOrganizationFederationIssuerUpdateParams, opts ...option.RequestOption) (res *BetaFederationIssuer, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	if federationIssuerID == "" {
		err = errors.New("missing required federation_issuer_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/federation_issuers/%s?beta=true", federationIssuerID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// List federation issuers in your organization.
//
// Archived issuers are excluded unless `include_archived=true`.
func (r *BetaOrganizationFederationIssuerService) List(ctx context.Context, params BetaOrganizationFederationIssuerListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaFederationIssuer], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/federation_issuers?beta=true"
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

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// List federation issuers in your organization.
//
// Archived issuers are excluded unless `include_archived=true`.
func (r *BetaOrganizationFederationIssuerService) ListAutoPaging(ctx context.Context, params BetaOrganizationFederationIssuerListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaFederationIssuer] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Archive a federation issuer.
//
// Idempotent; re-archiving returns the issuer with its original `archived_at`.
// Rejected with 400 if any live (non-archived) federation rule still references
// the issuer; archive those rules first (a rule's issuer cannot be changed), or
// recreate them against another issuer.
func (r *BetaOrganizationFederationIssuerService) Archive(ctx context.Context, federationIssuerID string, body BetaOrganizationFederationIssuerArchiveParams, opts ...option.RequestOption) (res *BetaFederationIssuer, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	if federationIssuerID == "" {
		err = errors.New("missing required federation_issuer_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/federation_issuers/%s/archive?beta=true", federationIssuerID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Registered external OIDC identity provider.
//
// Records an external IdP the organization trusts for the RFC 7523 jwt-bearer
// grant. The `issuer_url` must match the JWT `iss` claim exactly.
type BetaFederationIssuer struct {
	// Tagged ID of the federation issuer.
	ID string `json:"id" api:"required"`
	// If set, all rules referencing this issuer reject token exchange.
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// Tagged ID (`user_`/`svac_`) of the actor that archived this issuer.
	ArchivedByActorID string `json:"archived_by_actor_id" api:"required"`
	// Whether the jwt-bearer exchange enforces JTI single-use (replay protection) for
	// tokens from this issuer. Applies only to assertions carrying a `jti` claim;
	// tokens without one are accepted without single-use enforcement.
	CheckJTI bool `json:"check_jti" api:"required"`
	// When this issuer was created.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Tagged ID (`user_`/`svac_`) of the actor that created this issuer.
	CreatedByActorID string `json:"created_by_actor_id" api:"required"`
	// The `iss` claim value. Incoming JWTs must match exactly.
	IssuerURL string `json:"issuer_url" api:"required"`
	// How signing keys are obtained for signature verification.
	JWKS BetaFederationIssuerJWKSUnion `json:"jwks" api:"required"`
	// If set, Anthropic's JWKS poller has paused polling for this issuer after
	// repeated fetch failures. Re-enable by sending `jwks_polling_disabled: false` via
	// the issuer update endpoint (POST) once the upstream JWKS endpoint is fixed. An
	// OAuth caller cannot send this when the issuer backs a rule with any scope other
	// than `workspace:developer` or `workspace:inference`; use a Console session.
	JWKSPollingDisabledAt time.Time `json:"jwks_polling_disabled_at" api:"required" format:"date-time"`
	// Maximum allowed iat→exp spread for assertions from this issuer (1-176400
	// seconds, i.e. up to 49h). Assertions must carry both `iat` and `exp`; a missing
	// `iat` is rejected.
	MaxJWTLifetimeSeconds int64 `json:"max_jwt_lifetime_seconds" api:"required"`
	// Admin-chosen slug identifier.
	Name string `json:"name" api:"required"`
	// Status of automatic JWKS polling for a federation issuer.
	//
	// Anthropic periodically fetches the issuer's signing keys in the background.
	// These fields summarize the most recent fetches so the health of the JWKS
	// endpoint can be monitored.
	PollStatus BetaFederationIssuerPollStatus `json:"poll_status" api:"required"`
	Type       constant.FederationIssuer      `json:"type" default:"federation_issuer"`
	// When this issuer was last updated.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// Tagged ID (`user_`/`svac_`) of the actor that last updated this issuer.
	UpdatedByActorID string `json:"updated_by_actor_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                    respjson.Field
		ArchivedAt            respjson.Field
		ArchivedByActorID     respjson.Field
		CheckJTI              respjson.Field
		CreatedAt             respjson.Field
		CreatedByActorID      respjson.Field
		IssuerURL             respjson.Field
		JWKS                  respjson.Field
		JWKSPollingDisabledAt respjson.Field
		MaxJWTLifetimeSeconds respjson.Field
		Name                  respjson.Field
		PollStatus            respjson.Field
		Type                  respjson.Field
		UpdatedAt             respjson.Field
		UpdatedByActorID      respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaFederationIssuer) RawJSON() string { return r.JSON.raw }
func (r *BetaFederationIssuer) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaFederationIssuerJWKSUnion contains all possible properties and values from
// [BetaJWKSDiscovery], [BetaJWKSExplicitURL], [BetaJWKSInline].
//
// Use the [BetaFederationIssuerJWKSUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaFederationIssuerJWKSUnion struct {
	// Any of "discovery", "explicit_url", "inline".
	Type      string `json:"type"`
	CACertPEM string `json:"ca_cert_pem"`
	// This field is from variant [BetaJWKSDiscovery].
	DiscoveryBase string `json:"discovery_base"`
	// This field is from variant [BetaJWKSExplicitURL].
	URL string `json:"url"`
	// This field is from variant [BetaJWKSInline].
	Keys []map[string]any `json:"keys"`
	JSON struct {
		Type          respjson.Field
		CACertPEM     respjson.Field
		DiscoveryBase respjson.Field
		URL           respjson.Field
		Keys          respjson.Field
		raw           string
	} `json:"-"`
}

// anyBetaFederationIssuerJWKS is implemented by each variant of
// [BetaFederationIssuerJWKSUnion] to add type safety for the return type of
// [BetaFederationIssuerJWKSUnion.AsAny]
type anyBetaFederationIssuerJWKS interface {
	implBetaFederationIssuerJWKSUnion()
}

func (BetaJWKSDiscovery) implBetaFederationIssuerJWKSUnion()   {}
func (BetaJWKSExplicitURL) implBetaFederationIssuerJWKSUnion() {}
func (BetaJWKSInline) implBetaFederationIssuerJWKSUnion()      {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaFederationIssuerJWKSUnion.AsAny().(type) {
//	case anthropic.BetaJWKSDiscovery:
//	case anthropic.BetaJWKSExplicitURL:
//	case anthropic.BetaJWKSInline:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaFederationIssuerJWKSUnion) AsAny() anyBetaFederationIssuerJWKS {
	switch u.Type {
	case "discovery":
		return u.AsDiscovery()
	case "explicit_url":
		return u.AsExplicitURL()
	case "inline":
		return u.AsInline()
	}
	return nil
}

func (u BetaFederationIssuerJWKSUnion) AsDiscovery() (v BetaJWKSDiscovery) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaFederationIssuerJWKSUnion) AsExplicitURL() (v BetaJWKSExplicitURL) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaFederationIssuerJWKSUnion) AsInline() (v BetaJWKSInline) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaFederationIssuerJWKSUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaFederationIssuerJWKSUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Status of automatic JWKS polling for a federation issuer.
//
// Anthropic periodically fetches the issuer's signing keys in the background.
// These fields summarize the most recent fetches so the health of the JWKS
// endpoint can be monitored.
type BetaFederationIssuerPollStatus struct {
	// Consecutive fetch failures since the last success.
	ConsecutiveFailures int64 `json:"consecutive_failures" api:"required"`
	// When the last successful fetch completed.
	LastFetchedAt time.Time `json:"last_fetched_at" api:"required" format:"date-time"`
	// When the next fetch is scheduled. Null if paused.
	NextPollAt time.Time `json:"next_poll_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ConsecutiveFailures respjson.Field
		LastFetchedAt       respjson.Field
		NextPollAt          respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaFederationIssuerPollStatus) RawJSON() string { return r.JSON.raw }
func (r *BetaFederationIssuerPollStatus) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// JWKS via the issuer's OIDC discovery document.
type BetaJWKSDiscovery struct {
	Type constant.Discovery `json:"type" default:"discovery"`
	// Optional custom CA (PEM) for TLS verification of the JWKS fetch.
	CACertPEM string `json:"ca_cert_pem" api:"nullable"`
	// Set when the discovery URL differs from `issuer_url`.
	DiscoveryBase string `json:"discovery_base" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type          respjson.Field
		CACertPEM     respjson.Field
		DiscoveryBase respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaJWKSDiscovery) RawJSON() string { return r.JSON.raw }
func (r *BetaJWKSDiscovery) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaJWKSDiscovery to a BetaJWKSDiscoveryParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaJWKSDiscoveryParam.Overrides()
func (r BetaJWKSDiscovery) ToParam() BetaJWKSDiscoveryParam {
	return param.Override[BetaJWKSDiscoveryParam](json.RawMessage(r.RawJSON()))
}

// JWKS via the issuer's OIDC discovery document.
//
// The property Type is required.
type BetaJWKSDiscoveryParam struct {
	// Optional custom CA (PEM) for TLS verification of the JWKS fetch.
	CACertPEM param.Opt[string] `json:"ca_cert_pem,omitzero"`
	// Set when the discovery URL differs from `issuer_url`.
	DiscoveryBase param.Opt[string] `json:"discovery_base,omitzero"`
	// This field can be elided, and will marshal its zero value as "discovery".
	Type constant.Discovery `json:"type" default:"discovery"`
	paramObj
}

func (r BetaJWKSDiscoveryParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaJWKSDiscoveryParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaJWKSDiscoveryParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// JWKS fetched from a fixed endpoint.
type BetaJWKSExplicitURL struct {
	Type constant.ExplicitURL `json:"type" default:"explicit_url"`
	// JWKS endpoint.
	URL string `json:"url" api:"required"`
	// Optional custom CA (PEM) for TLS verification of the JWKS fetch.
	CACertPEM string `json:"ca_cert_pem" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		URL         respjson.Field
		CACertPEM   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaJWKSExplicitURL) RawJSON() string { return r.JSON.raw }
func (r *BetaJWKSExplicitURL) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaJWKSExplicitURL to a BetaJWKSExplicitURLParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaJWKSExplicitURLParam.Overrides()
func (r BetaJWKSExplicitURL) ToParam() BetaJWKSExplicitURLParam {
	return param.Override[BetaJWKSExplicitURLParam](json.RawMessage(r.RawJSON()))
}

// JWKS fetched from a fixed endpoint.
//
// The properties Type, URL are required.
type BetaJWKSExplicitURLParam struct {
	// JWKS endpoint.
	URL string `json:"url" api:"required"`
	// Optional custom CA (PEM) for TLS verification of the JWKS fetch.
	CACertPEM param.Opt[string] `json:"ca_cert_pem,omitzero"`
	// This field can be elided, and will marshal its zero value as "explicit_url".
	Type constant.ExplicitURL `json:"type" default:"explicit_url"`
	paramObj
}

func (r BetaJWKSExplicitURLParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaJWKSExplicitURLParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaJWKSExplicitURLParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// JWKS supplied directly; no network fetch.
type BetaJWKSInline struct {
	// Inline JWK objects.
	Keys []map[string]any `json:"keys" api:"required"`
	Type constant.Inline  `json:"type" default:"inline"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Keys        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaJWKSInline) RawJSON() string { return r.JSON.raw }
func (r *BetaJWKSInline) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaJWKSInline to a BetaJWKSInlineParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaJWKSInlineParam.Overrides()
func (r BetaJWKSInline) ToParam() BetaJWKSInlineParam {
	return param.Override[BetaJWKSInlineParam](json.RawMessage(r.RawJSON()))
}

// JWKS supplied directly; no network fetch.
//
// The properties Keys, Type are required.
type BetaJWKSInlineParam struct {
	// Inline JWK objects.
	Keys []map[string]any `json:"keys,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "inline".
	Type constant.Inline `json:"type" default:"inline"`
	paramObj
}

func (r BetaJWKSInlineParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaJWKSInlineParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaJWKSInlineParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationFederationIssuerNewParams struct {
	// The `iss` claim value to match against.
	IssuerURL string `json:"issuer_url" api:"required"`
	// Slug identifier (lowercase, digits, hyphens). Unique within the organization; a
	// duplicate name returns 409.
	Name string `json:"name" api:"required"`
	// Whether the jwt-bearer exchange enforces JTI single-use (replay protection) for
	// tokens from this issuer. Defaults to true. Applies only to assertions carrying a
	// `jti` claim; tokens without one are accepted without single-use enforcement.
	CheckJTI param.Opt[bool] `json:"check_jti,omitzero"`
	// Maximum allowed iat→exp spread for assertions from this issuer (1-176400
	// seconds, i.e. up to 49h). Defaults to 3600 (1h). Assertions must carry both
	// `iat` and `exp`; a missing `iat` is rejected.
	MaxJWTLifetimeSeconds param.Opt[int64] `json:"max_jwt_lifetime_seconds,omitzero"`
	// How signing keys are obtained. Defaults to OIDC discovery.
	JWKS BetaOrganizationFederationIssuerNewParamsJWKSUnion `json:"jwks,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationFederationIssuerNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationFederationIssuerNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationFederationIssuerNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaOrganizationFederationIssuerNewParamsJWKSUnion struct {
	OfDiscovery   *BetaJWKSDiscoveryParam   `json:",omitzero,inline"`
	OfExplicitURL *BetaJWKSExplicitURLParam `json:",omitzero,inline"`
	OfInline      *BetaJWKSInlineParam      `json:",omitzero,inline"`
	paramUnion
}

func (u BetaOrganizationFederationIssuerNewParamsJWKSUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDiscovery, u.OfExplicitURL, u.OfInline)
}
func (u *BetaOrganizationFederationIssuerNewParamsJWKSUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaOrganizationFederationIssuerNewParamsJWKSUnion) asAny() any {
	if !param.IsOmitted(u.OfDiscovery) {
		return u.OfDiscovery
	} else if !param.IsOmitted(u.OfExplicitURL) {
		return u.OfExplicitURL
	} else if !param.IsOmitted(u.OfInline) {
		return u.OfInline
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationFederationIssuerNewParamsJWKSUnion) GetDiscoveryBase() *string {
	if vt := u.OfDiscovery; vt != nil && vt.DiscoveryBase.Valid() {
		return &vt.DiscoveryBase.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationFederationIssuerNewParamsJWKSUnion) GetURL() *string {
	if vt := u.OfExplicitURL; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationFederationIssuerNewParamsJWKSUnion) GetKeys() []map[string]any {
	if vt := u.OfInline; vt != nil {
		return vt.Keys
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationFederationIssuerNewParamsJWKSUnion) GetType() *string {
	if vt := u.OfDiscovery; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfExplicitURL; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfInline; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationFederationIssuerNewParamsJWKSUnion) GetCACertPEM() *string {
	if vt := u.OfDiscovery; vt != nil && vt.CACertPEM.Valid() {
		return &vt.CACertPEM.Value
	} else if vt := u.OfExplicitURL; vt != nil && vt.CACertPEM.Valid() {
		return &vt.CACertPEM.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaOrganizationFederationIssuerNewParamsJWKSUnion](
		"type",
		apijson.Discriminator[BetaJWKSDiscoveryParam]("discovery"),
		apijson.Discriminator[BetaJWKSExplicitURLParam]("explicit_url"),
		apijson.Discriminator[BetaJWKSInlineParam]("inline"),
	)
}

type BetaOrganizationFederationIssuerGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaOrganizationFederationIssuerUpdateParams struct {
	// Whether the jwt-bearer exchange enforces JTI single-use (replay protection) for
	// tokens from this issuer. Applies only to assertions carrying a `jti` claim;
	// tokens without one are accepted without single-use enforcement.
	CheckJTI param.Opt[bool] `json:"check_jti,omitzero"`
	// Replaces the `iss` claim value to match against. For discovery-mode issuers
	// without a `discovery_base`, this is also the URL Anthropic fetches the OIDC
	// discovery document and signing keys from, so changing it repoints the JWKS
	// source. Changing the issuer URL to a well-known shared platform is rejected
	// while any live rule under this issuer would not constrain tenant identity.
	IssuerURL param.Opt[string] `json:"issuer_url,omitzero"`
	// Only `false` is accepted, to re-enable polling after the system pauses it.
	// Polling is paused automatically; sending `true` is rejected.
	JWKSPollingDisabled param.Opt[bool] `json:"jwks_polling_disabled,omitzero"`
	// Maximum allowed iat→exp spread for assertions from this issuer (1-176400
	// seconds, i.e. up to 49h). Assertions must carry both `iat` and `exp`; a missing
	// `iat` is rejected.
	MaxJWTLifetimeSeconds param.Opt[int64] `json:"max_jwt_lifetime_seconds,omitzero"`
	// Replaces the slug identifier (lowercase, digits, hyphens). Unique within the
	// organization; a duplicate name returns 409.
	Name param.Opt[string] `json:"name,omitzero"`
	// Replaces the entire JWKS configuration.
	JWKS BetaOrganizationFederationIssuerUpdateParamsJWKSUnion `json:"jwks,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaOrganizationFederationIssuerUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationFederationIssuerUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationFederationIssuerUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaOrganizationFederationIssuerUpdateParamsJWKSUnion struct {
	OfDiscovery   *BetaJWKSDiscoveryParam   `json:",omitzero,inline"`
	OfExplicitURL *BetaJWKSExplicitURLParam `json:",omitzero,inline"`
	OfInline      *BetaJWKSInlineParam      `json:",omitzero,inline"`
	paramUnion
}

func (u BetaOrganizationFederationIssuerUpdateParamsJWKSUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDiscovery, u.OfExplicitURL, u.OfInline)
}
func (u *BetaOrganizationFederationIssuerUpdateParamsJWKSUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaOrganizationFederationIssuerUpdateParamsJWKSUnion) asAny() any {
	if !param.IsOmitted(u.OfDiscovery) {
		return u.OfDiscovery
	} else if !param.IsOmitted(u.OfExplicitURL) {
		return u.OfExplicitURL
	} else if !param.IsOmitted(u.OfInline) {
		return u.OfInline
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationFederationIssuerUpdateParamsJWKSUnion) GetDiscoveryBase() *string {
	if vt := u.OfDiscovery; vt != nil && vt.DiscoveryBase.Valid() {
		return &vt.DiscoveryBase.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationFederationIssuerUpdateParamsJWKSUnion) GetURL() *string {
	if vt := u.OfExplicitURL; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationFederationIssuerUpdateParamsJWKSUnion) GetKeys() []map[string]any {
	if vt := u.OfInline; vt != nil {
		return vt.Keys
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationFederationIssuerUpdateParamsJWKSUnion) GetType() *string {
	if vt := u.OfDiscovery; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfExplicitURL; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfInline; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOrganizationFederationIssuerUpdateParamsJWKSUnion) GetCACertPEM() *string {
	if vt := u.OfDiscovery; vt != nil && vt.CACertPEM.Valid() {
		return &vt.CACertPEM.Value
	} else if vt := u.OfExplicitURL; vt != nil && vt.CACertPEM.Valid() {
		return &vt.CACertPEM.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaOrganizationFederationIssuerUpdateParamsJWKSUnion](
		"type",
		apijson.Discriminator[BetaJWKSDiscoveryParam]("discovery"),
		apijson.Discriminator[BetaJWKSExplicitURLParam]("explicit_url"),
		apijson.Discriminator[BetaJWKSInlineParam]("inline"),
	)
}

type BetaOrganizationFederationIssuerListParams struct {
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Include archived resources. Defaults to false.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Number of results per page.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationFederationIssuerListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationFederationIssuerListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationFederationIssuerArchiveParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
