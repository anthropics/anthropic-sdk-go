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
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// BetaTunnelCertificateService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaTunnelCertificateService] method instead.
type BetaTunnelCertificateService struct {
	Options []option.RequestOption
}

// NewBetaTunnelCertificateService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaTunnelCertificateService(opts ...option.RequestOption) (r BetaTunnelCertificateService) {
	r = BetaTunnelCertificateService{}
	r.Options = opts
	return
}

// The Tunnels API is in research preview. It requires the
// `anthropic-beta: mcp-tunnels-2026-06-22` header and may change without a
// deprecation period. It supersedes the Admin API endpoints at
// `/v1/organizations/tunnels`, which remain available during a migration window.
//
// Registers a public CA certificate on a tunnel. Anthropic verifies the gateway's
// server certificate against this CA when it terminates the inner TLS session. A
// tunnel holds at most two non-archived certificates.
func (r *BetaTunnelCertificateService) New(ctx context.Context, tunnelID string, params BetaTunnelCertificateNewParams, opts ...option.RequestOption) (res *BetaTunnelCertificate, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "mcp-tunnels-2026-06-22")}, opts...)
	if tunnelID == "" {
		err = errors.New("missing required tunnel_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/tunnels/%s/certificates?beta=true", tunnelID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// The Tunnels API is in research preview. It requires the
// `anthropic-beta: mcp-tunnels-2026-06-22` header and may change without a
// deprecation period. It supersedes the Admin API endpoints at
// `/v1/organizations/tunnels`, which remain available during a migration window.
//
// Fetches a tunnel certificate by ID.
func (r *BetaTunnelCertificateService) Get(ctx context.Context, certificateID string, params BetaTunnelCertificateGetParams, opts ...option.RequestOption) (res *BetaTunnelCertificate, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "mcp-tunnels-2026-06-22")}, opts...)
	if params.TunnelID == "" {
		err = errors.New("missing required tunnel_id parameter")
		return nil, err
	}
	if certificateID == "" {
		err = errors.New("missing required certificate_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/tunnels/%s/certificates/%s?beta=true", params.TunnelID, certificateID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// The Tunnels API is in research preview. It requires the
// `anthropic-beta: mcp-tunnels-2026-06-22` header and may change without a
// deprecation period. It supersedes the Admin API endpoints at
// `/v1/organizations/tunnels`, which remain available during a migration window.
//
// Lists the certificates registered on a tunnel. Archived certificates are
// excluded unless include_archived is set.
func (r *BetaTunnelCertificateService) List(ctx context.Context, tunnelID string, params BetaTunnelCertificateListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaTunnelCertificate], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "mcp-tunnels-2026-06-22"), option.WithResponseInto(&raw)}, opts...)
	if tunnelID == "" {
		err = errors.New("missing required tunnel_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/tunnels/%s/certificates?beta=true", tunnelID)
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

// The Tunnels API is in research preview. It requires the
// `anthropic-beta: mcp-tunnels-2026-06-22` header and may change without a
// deprecation period. It supersedes the Admin API endpoints at
// `/v1/organizations/tunnels`, which remain available during a migration window.
//
// Lists the certificates registered on a tunnel. Archived certificates are
// excluded unless include_archived is set.
func (r *BetaTunnelCertificateService) ListAutoPaging(ctx context.Context, tunnelID string, params BetaTunnelCertificateListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaTunnelCertificate] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, tunnelID, params, opts...))
}

// The Tunnels API is in research preview. It requires the
// `anthropic-beta: mcp-tunnels-2026-06-22` header and may change without a
// deprecation period. It supersedes the Admin API endpoints at
// `/v1/organizations/tunnels`, which remain available during a migration window.
//
// Archives a tunnel certificate, removing it from the set Anthropic trusts for the
// tunnel. The certificate record is retained. Archiving the last non-archived
// certificate is permitted; the tunnel rejects MCP traffic until a new certificate
// is added.
func (r *BetaTunnelCertificateService) Archive(ctx context.Context, certificateID string, params BetaTunnelCertificateArchiveParams, opts ...option.RequestOption) (res *BetaTunnelCertificate, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "mcp-tunnels-2026-06-22")}, opts...)
	if params.TunnelID == "" {
		err = errors.New("missing required tunnel_id parameter")
		return nil, err
	}
	if certificateID == "" {
		err = errors.New("missing required certificate_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/tunnels/%s/certificates/%s/archive?beta=true", params.TunnelID, certificateID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// A CA certificate attached to a tunnel.
type BetaTunnelCertificate struct {
	// Unique identifier for the certificate, prefixed with `tcrt_`.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	ExpiresAt time.Time `json:"expires_at" api:"required" format:"date-time"`
	// Lowercase hex SHA-256 fingerprint of the certificate's DER encoding.
	Fingerprint string `json:"fingerprint" api:"required"`
	// ID of the tunnel the certificate is registered against.
	TunnelID string                     `json:"tunnel_id" api:"required"`
	Type     constant.TunnelCertificate `json:"type" default:"tunnel_certificate"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ArchivedAt  respjson.Field
		CreatedAt   respjson.Field
		ExpiresAt   respjson.Field
		Fingerprint respjson.Field
		TunnelID    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTunnelCertificate) RawJSON() string { return r.JSON.raw }
func (r *BetaTunnelCertificate) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaTunnelCertificateNewParams struct {
	// PEM-encoded X.509 CA certificate. Must contain exactly one certificate and no
	// private-key material. Maximum 8KB.
	CaCertificatePem string `json:"ca_certificate_pem" api:"required"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaTunnelCertificateNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaTunnelCertificateNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaTunnelCertificateNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaTunnelCertificateGetParams struct {
	TunnelID string `path:"tunnel_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaTunnelCertificateListParams struct {
	// Whether to include archived certificates in the results. Defaults to false.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum number of certificates to return per page. Defaults to 20, maximum 1000.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor from a previous `list_tunnel_certificates` response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaTunnelCertificateListParams]'s query parameters as
// `url.Values`.
func (r BetaTunnelCertificateListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaTunnelCertificateArchiveParams struct {
	TunnelID string `path:"tunnel_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
