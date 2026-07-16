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

// BetaTunnelService contains methods and other services that help with interacting
// with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaTunnelService] method instead.
type BetaTunnelService struct {
	Options      []option.RequestOption
	Certificates BetaTunnelCertificateService
}

// NewBetaTunnelService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaTunnelService(opts ...option.RequestOption) (r BetaTunnelService) {
	r = BetaTunnelService{}
	r.Options = opts
	r.Certificates = NewBetaTunnelCertificateService(opts...)
	return
}

// The Tunnels API is in research preview. It requires the
// `anthropic-beta: mcp-tunnels-2026-06-22` header and may change without a
// deprecation period. It supersedes the Admin API endpoints at
// `/v1/organizations/tunnels`, which remain available during a migration window.
//
// Creates a tunnel. Creation allocates a fresh hostname and provisions the tunnel;
// it is not idempotent. The new tunnel rejects MCP traffic until at least one CA
// certificate is added.
func (r *BetaTunnelService) New(ctx context.Context, params BetaTunnelNewParams, opts ...option.RequestOption) (res *BetaTunnel, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "mcp-tunnels-2026-06-22")}, opts...)
	path := "v1/tunnels?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// The Tunnels API is in research preview. It requires the
// `anthropic-beta: mcp-tunnels-2026-06-22` header and may change without a
// deprecation period. It supersedes the Admin API endpoints at
// `/v1/organizations/tunnels`, which remain available during a migration window.
//
// Fetches a tunnel by ID.
func (r *BetaTunnelService) Get(ctx context.Context, tunnelID string, query BetaTunnelGetParams, opts ...option.RequestOption) (res *BetaTunnel, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "mcp-tunnels-2026-06-22")}, opts...)
	if tunnelID == "" {
		err = errors.New("missing required tunnel_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/tunnels/%s?beta=true", tunnelID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// The Tunnels API is in research preview. It requires the
// `anthropic-beta: mcp-tunnels-2026-06-22` header and may change without a
// deprecation period. It supersedes the Admin API endpoints at
// `/v1/organizations/tunnels`, which remain available during a migration window.
//
// Lists tunnels. Results are ordered by creation time, newest first; archived
// tunnels are excluded unless include_archived is set.
func (r *BetaTunnelService) List(ctx context.Context, params BetaTunnelListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaTunnel], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "mcp-tunnels-2026-06-22"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/tunnels?beta=true"
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
// Lists tunnels. Results are ordered by creation time, newest first; archived
// tunnels are excluded unless include_archived is set.
func (r *BetaTunnelService) ListAutoPaging(ctx context.Context, params BetaTunnelListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaTunnel] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// The Tunnels API is in research preview. It requires the
// `anthropic-beta: mcp-tunnels-2026-06-22` header and may change without a
// deprecation period. It supersedes the Admin API endpoints at
// `/v1/organizations/tunnels`, which remain available during a migration window.
//
// Archives a tunnel. Archival is irreversible: every non-archived certificate on
// the tunnel is archived in the same operation, the hostname is retired and never
// re-allocated, and the tunnel token is invalidated. Retrying against an
// already-archived tunnel returns the existing record unchanged.
func (r *BetaTunnelService) Archive(ctx context.Context, tunnelID string, body BetaTunnelArchiveParams, opts ...option.RequestOption) (res *BetaTunnel, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "mcp-tunnels-2026-06-22")}, opts...)
	if tunnelID == "" {
		err = errors.New("missing required tunnel_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/tunnels/%s/archive?beta=true", tunnelID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// The Tunnels API is in research preview. It requires the
// `anthropic-beta: mcp-tunnels-2026-06-22` header and may change without a
// deprecation period. It supersedes the Admin API endpoints at
// `/v1/organizations/tunnels`, which remain available during a migration window.
//
// Reveals a tunnel's connector token. The value is fetched live on each call;
// Anthropic does not store it. Repeated calls return the same value until the
// token is rotated. Exposed as POST so the token does not appear in intermediary
// access logs.
func (r *BetaTunnelService) RevealToken(ctx context.Context, tunnelID string, body BetaTunnelRevealTokenParams, opts ...option.RequestOption) (res *BetaTunnelToken, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "mcp-tunnels-2026-06-22")}, opts...)
	if tunnelID == "" {
		err = errors.New("missing required tunnel_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/tunnels/%s/reveal_token?beta=true", tunnelID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// The Tunnels API is in research preview. It requires the
// `anthropic-beta: mcp-tunnels-2026-06-22` header and may change without a
// deprecation period. It supersedes the Admin API endpoints at
// `/v1/organizations/tunnels`, which remain available during a migration window.
//
// Rotates a tunnel's connector token. Rotation invalidates the current token for
// new connections and returns a fresh value; established connections are not
// severed. A connector restarted after rotation must use the new value.
func (r *BetaTunnelService) RotateToken(ctx context.Context, tunnelID string, params BetaTunnelRotateTokenParams, opts ...option.RequestOption) (res *BetaTunnelToken, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "mcp-tunnels-2026-06-22")}, opts...)
	if tunnelID == "" {
		err = errors.New("missing required tunnel_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/tunnels/%s/rotate_token?beta=true", tunnelID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// An MCP tunnel.
type BetaTunnel struct {
	// Unique identifier for the tunnel, prefixed with `tnl_`.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Human-readable name for the tunnel (1-255 characters). Null if unset.
	DisplayName string `json:"display_name" api:"required"`
	// Anthropic-assigned hostname for the tunnel. MCP server URLs whose host is a
	// subdomain of this value are routed through the tunnel. Globally unique and never
	// reused, even after the tunnel is archived.
	Domain string          `json:"domain" api:"required"`
	Type   constant.Tunnel `json:"type" default:"tunnel"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ArchivedAt  respjson.Field
		CreatedAt   respjson.Field
		DisplayName respjson.Field
		Domain      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTunnel) RawJSON() string { return r.JSON.raw }
func (r *BetaTunnel) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A tunnel's connector token.
type BetaTunnelToken struct {
	// Stable identifier for the current token value. Changes when the token is
	// rotated.
	ID string `json:"id" api:"required"`
	// The connector token used to run the tunnel. Treat as a credential.
	TunnelToken string               `json:"tunnel_token" api:"required"`
	Type        constant.TunnelToken `json:"type" default:"tunnel_token"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		TunnelToken respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaTunnelToken) RawJSON() string { return r.JSON.raw }
func (r *BetaTunnelToken) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaTunnelNewParams struct {
	// Optional human-readable name for the tunnel (1-255 characters).
	DisplayName param.Opt[string] `json:"display_name,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaTunnelNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaTunnelNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaTunnelNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaTunnelGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaTunnelListParams struct {
	// Whether to include archived tunnels in the results. Defaults to false.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum number of tunnels to return per page. Defaults to 20, maximum 1000.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor from a previous `list_tunnels` response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaTunnelListParams]'s query parameters as `url.Values`.
func (r BetaTunnelListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaTunnelArchiveParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaTunnelRevealTokenParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaTunnelRotateTokenParams struct {
	// Optional free-text reason for the rotation, recorded for audit.
	Reason param.Opt[string] `json:"reason,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaTunnelRotateTokenParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaTunnelRotateTokenParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaTunnelRotateTokenParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
