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

// OrganizationFederationIssuerService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationFederationIssuerService] method instead.
type OrganizationFederationIssuerService struct {
	Options []option.RequestOption
}

// NewOrganizationFederationIssuerService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewOrganizationFederationIssuerService(opts ...option.RequestOption) (r OrganizationFederationIssuerService) {
	r = OrganizationFederationIssuerService{}
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
func (r *OrganizationFederationIssuerService) New(ctx context.Context, body OrganizationFederationIssuerNewParams, opts ...option.RequestOption) (res *FederationIssuer, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/organizations/federation_issuers"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// Retrieve a federation issuer by its ID (`fdis_...`).
func (r *OrganizationFederationIssuerService) Get(ctx context.Context, federationIssuerID string, opts ...option.RequestOption) (res *FederationIssuer, err error) {
	opts = slices.Concat(r.Options, opts)
	if federationIssuerID == "" {
		err = errors.New("missing required federation_issuer_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/federation_issuers/%s", url.PathEscape(federationIssuerID))
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
func (r *OrganizationFederationIssuerService) Update(ctx context.Context, federationIssuerID string, body OrganizationFederationIssuerUpdateParams, opts ...option.RequestOption) (res *FederationIssuer, err error) {
	opts = slices.Concat(r.Options, opts)
	if federationIssuerID == "" {
		err = errors.New("missing required federation_issuer_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/federation_issuers/%s", url.PathEscape(federationIssuerID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
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
func (r *OrganizationFederationIssuerService) List(ctx context.Context, query OrganizationFederationIssuerListParams, opts ...option.RequestOption) (res *pagination.PageCursor[FederationIssuer], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/federation_issuers"
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

// **Requires an OAuth access token with the `org:admin` scope**, from
// `ant auth login --scope org:admin` or a workload identity federation rule; Admin
// API keys are not accepted. See
// [Manage WIF with the Admin API](/docs/en/manage-claude/wif-admin-api).
//
// List federation issuers in your organization.
//
// Archived issuers are excluded unless `include_archived=true`.
func (r *OrganizationFederationIssuerService) ListAutoPaging(ctx context.Context, query OrganizationFederationIssuerListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[FederationIssuer] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
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
func (r *OrganizationFederationIssuerService) Archive(ctx context.Context, federationIssuerID string, opts ...option.RequestOption) (res *FederationIssuer, err error) {
	opts = slices.Concat(r.Options, opts)
	if federationIssuerID == "" {
		err = errors.New("missing required federation_issuer_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/federation_issuers/%s/archive", url.PathEscape(federationIssuerID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Registered external OIDC identity provider.
//
// Records an external IdP the organization trusts for the RFC 7523 jwt-bearer
// grant. The `issuer_url` must match the JWT `iss` claim exactly.
type FederationIssuer struct {
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
	JWKS FederationIssuerJWKSUnion `json:"jwks" api:"required"`
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
	// Live state of Anthropic's JWKS polling for this issuer. Populated on both
	// single-issuer retrieval and list responses, including archived issuers.
	// Typically null for inline-key issuers (no polling), or when poll status is
	// temporarily unavailable or polling has not started yet.
	PollStatus FederationIssuerPollStatus `json:"poll_status" api:"required"`
	Type       constant.FederationIssuer  `json:"type" default:"federation_issuer"`
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
func (r FederationIssuer) RawJSON() string { return r.JSON.raw }
func (r *FederationIssuer) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// FederationIssuerJWKSUnion contains all possible properties and values from
// [JWKSDiscovery], [JWKSExplicitURL], [JWKSInline].
//
// Use the [FederationIssuerJWKSUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type FederationIssuerJWKSUnion struct {
	// Any of "discovery", "explicit_url", "inline".
	Type      string `json:"type"`
	CACertPEM string `json:"ca_cert_pem"`
	// This field is from variant [JWKSDiscovery].
	DiscoveryBase string `json:"discovery_base"`
	// This field is from variant [JWKSExplicitURL].
	URL string `json:"url"`
	// This field is from variant [JWKSInline].
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

// anyFederationIssuerJWKS is implemented by each variant of
// [FederationIssuerJWKSUnion] to add type safety for the return type of
// [FederationIssuerJWKSUnion.AsAny]
type anyFederationIssuerJWKS interface {
	implFederationIssuerJWKSUnion()
}

func (JWKSDiscovery) implFederationIssuerJWKSUnion()   {}
func (JWKSExplicitURL) implFederationIssuerJWKSUnion() {}
func (JWKSInline) implFederationIssuerJWKSUnion()      {}

// Use the following switch statement to find the correct variant
//
//	switch variant := FederationIssuerJWKSUnion.AsAny().(type) {
//	case anthropic.JWKSDiscovery:
//	case anthropic.JWKSExplicitURL:
//	case anthropic.JWKSInline:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u FederationIssuerJWKSUnion) AsAny() anyFederationIssuerJWKS {
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

func (u FederationIssuerJWKSUnion) AsDiscovery() (v JWKSDiscovery) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FederationIssuerJWKSUnion) AsExplicitURL() (v JWKSExplicitURL) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u FederationIssuerJWKSUnion) AsInline() (v JWKSInline) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u FederationIssuerJWKSUnion) RawJSON() string { return u.JSON.raw }

func (r *FederationIssuerJWKSUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Status of automatic JWKS polling for a federation issuer.
//
// Anthropic periodically fetches the issuer's signing keys in the background.
// These fields summarize the most recent fetches so the health of the JWKS
// endpoint can be monitored.
type FederationIssuerPollStatus struct {
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
func (r FederationIssuerPollStatus) RawJSON() string { return r.JSON.raw }
func (r *FederationIssuerPollStatus) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// JWKS via the issuer's OIDC discovery document.
type JWKSDiscovery struct {
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
func (r JWKSDiscovery) RawJSON() string { return r.JSON.raw }
func (r *JWKSDiscovery) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this JWKSDiscovery to a JWKSDiscoveryParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// JWKSDiscoveryParam.Overrides()
func (r JWKSDiscovery) ToParam() JWKSDiscoveryParam {
	return param.Override[JWKSDiscoveryParam](json.RawMessage(r.RawJSON()))
}

// JWKS via the issuer's OIDC discovery document.
//
// The property Type is required.
type JWKSDiscoveryParam struct {
	// Optional custom CA (PEM) for TLS verification of the JWKS fetch.
	CACertPEM param.Opt[string] `json:"ca_cert_pem,omitzero"`
	// Set when the discovery URL differs from `issuer_url`.
	DiscoveryBase param.Opt[string] `json:"discovery_base,omitzero"`
	// This field can be elided, and will marshal its zero value as "discovery".
	Type constant.Discovery `json:"type" default:"discovery"`
	paramObj
}

func (r JWKSDiscoveryParam) MarshalJSON() (data []byte, err error) {
	type shadow JWKSDiscoveryParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *JWKSDiscoveryParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// JWKS fetched from a fixed endpoint.
type JWKSExplicitURL struct {
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
func (r JWKSExplicitURL) RawJSON() string { return r.JSON.raw }
func (r *JWKSExplicitURL) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this JWKSExplicitURL to a JWKSExplicitURLParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// JWKSExplicitURLParam.Overrides()
func (r JWKSExplicitURL) ToParam() JWKSExplicitURLParam {
	return param.Override[JWKSExplicitURLParam](json.RawMessage(r.RawJSON()))
}

// JWKS fetched from a fixed endpoint.
//
// The properties Type, URL are required.
type JWKSExplicitURLParam struct {
	// JWKS endpoint.
	URL string `json:"url" api:"required"`
	// Optional custom CA (PEM) for TLS verification of the JWKS fetch.
	CACertPEM param.Opt[string] `json:"ca_cert_pem,omitzero"`
	// This field can be elided, and will marshal its zero value as "explicit_url".
	Type constant.ExplicitURL `json:"type" default:"explicit_url"`
	paramObj
}

func (r JWKSExplicitURLParam) MarshalJSON() (data []byte, err error) {
	type shadow JWKSExplicitURLParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *JWKSExplicitURLParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// JWKS supplied directly; no network fetch.
type JWKSInline struct {
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
func (r JWKSInline) RawJSON() string { return r.JSON.raw }
func (r *JWKSInline) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this JWKSInline to a JWKSInlineParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// JWKSInlineParam.Overrides()
func (r JWKSInline) ToParam() JWKSInlineParam {
	return param.Override[JWKSInlineParam](json.RawMessage(r.RawJSON()))
}

// JWKS supplied directly; no network fetch.
//
// The properties Keys, Type are required.
type JWKSInlineParam struct {
	// Inline JWK objects.
	Keys []map[string]any `json:"keys,omitzero" api:"required"`
	// This field can be elided, and will marshal its zero value as "inline".
	Type constant.Inline `json:"type" default:"inline"`
	paramObj
}

func (r JWKSInlineParam) MarshalJSON() (data []byte, err error) {
	type shadow JWKSInlineParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *JWKSInlineParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationFederationIssuerNewParams struct {
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
	JWKS OrganizationFederationIssuerNewParamsJWKSUnion `json:"jwks,omitzero"`
	paramObj
}

func (r OrganizationFederationIssuerNewParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationFederationIssuerNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationFederationIssuerNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type OrganizationFederationIssuerNewParamsJWKSUnion struct {
	OfDiscovery   *JWKSDiscoveryParam   `json:",omitzero,inline"`
	OfExplicitURL *JWKSExplicitURLParam `json:",omitzero,inline"`
	OfInline      *JWKSInlineParam      `json:",omitzero,inline"`
	paramUnion
}

func (u OrganizationFederationIssuerNewParamsJWKSUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDiscovery, u.OfExplicitURL, u.OfInline)
}
func (u *OrganizationFederationIssuerNewParamsJWKSUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *OrganizationFederationIssuerNewParamsJWKSUnion) asAny() any {
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
func (u OrganizationFederationIssuerNewParamsJWKSUnion) GetDiscoveryBase() *string {
	if vt := u.OfDiscovery; vt != nil && vt.DiscoveryBase.Valid() {
		return &vt.DiscoveryBase.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationFederationIssuerNewParamsJWKSUnion) GetURL() *string {
	if vt := u.OfExplicitURL; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationFederationIssuerNewParamsJWKSUnion) GetKeys() []map[string]any {
	if vt := u.OfInline; vt != nil {
		return vt.Keys
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationFederationIssuerNewParamsJWKSUnion) GetType() *string {
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
func (u OrganizationFederationIssuerNewParamsJWKSUnion) GetCACertPEM() *string {
	if vt := u.OfDiscovery; vt != nil && vt.CACertPEM.Valid() {
		return &vt.CACertPEM.Value
	} else if vt := u.OfExplicitURL; vt != nil && vt.CACertPEM.Valid() {
		return &vt.CACertPEM.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[OrganizationFederationIssuerNewParamsJWKSUnion](
		"type",
		apijson.Discriminator[JWKSDiscoveryParam]("discovery"),
		apijson.Discriminator[JWKSExplicitURLParam]("explicit_url"),
		apijson.Discriminator[JWKSInlineParam]("inline"),
	)
}

type OrganizationFederationIssuerUpdateParams struct {
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
	JWKS OrganizationFederationIssuerUpdateParamsJWKSUnion `json:"jwks,omitzero"`
	paramObj
}

func (r OrganizationFederationIssuerUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationFederationIssuerUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationFederationIssuerUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type OrganizationFederationIssuerUpdateParamsJWKSUnion struct {
	OfDiscovery   *JWKSDiscoveryParam   `json:",omitzero,inline"`
	OfExplicitURL *JWKSExplicitURLParam `json:",omitzero,inline"`
	OfInline      *JWKSInlineParam      `json:",omitzero,inline"`
	paramUnion
}

func (u OrganizationFederationIssuerUpdateParamsJWKSUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfDiscovery, u.OfExplicitURL, u.OfInline)
}
func (u *OrganizationFederationIssuerUpdateParamsJWKSUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *OrganizationFederationIssuerUpdateParamsJWKSUnion) asAny() any {
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
func (u OrganizationFederationIssuerUpdateParamsJWKSUnion) GetDiscoveryBase() *string {
	if vt := u.OfDiscovery; vt != nil && vt.DiscoveryBase.Valid() {
		return &vt.DiscoveryBase.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationFederationIssuerUpdateParamsJWKSUnion) GetURL() *string {
	if vt := u.OfExplicitURL; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationFederationIssuerUpdateParamsJWKSUnion) GetKeys() []map[string]any {
	if vt := u.OfInline; vt != nil {
		return vt.Keys
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OrganizationFederationIssuerUpdateParamsJWKSUnion) GetType() *string {
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
func (u OrganizationFederationIssuerUpdateParamsJWKSUnion) GetCACertPEM() *string {
	if vt := u.OfDiscovery; vt != nil && vt.CACertPEM.Valid() {
		return &vt.CACertPEM.Value
	} else if vt := u.OfExplicitURL; vt != nil && vt.CACertPEM.Valid() {
		return &vt.CACertPEM.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[OrganizationFederationIssuerUpdateParamsJWKSUnion](
		"type",
		apijson.Discriminator[JWKSDiscoveryParam]("discovery"),
		apijson.Discriminator[JWKSExplicitURLParam]("explicit_url"),
		apijson.Discriminator[JWKSInlineParam]("inline"),
	)
}

type OrganizationFederationIssuerListParams struct {
	// Opaque cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Include archived resources. Defaults to false.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Number of results per page.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrganizationFederationIssuerListParams]'s query parameters
// as `url.Values`.
func (r OrganizationFederationIssuerListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
