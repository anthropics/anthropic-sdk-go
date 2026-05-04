package config

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal"
)

// OAuth 2.0 wire-contract constants. These are the authoritative
// definitions; internal/auth re-exports them as local aliases so that
// in-tree auth code can reference them without importing config and
// forming a cycle.
const (
	// TokenEndpoint is the path of the Anthropic OAuth 2.0 token
	// endpoint — the destination for jwt-bearer exchanges,
	// refresh_token grants, and (future) authorization_code grants.
	TokenEndpoint = "/v1/oauth/token"

	// GrantTypeJWTBearer is the RFC 7523 grant type string used for
	// OIDC federation exchanges.
	GrantTypeJWTBearer = "urn:ietf:params:oauth:grant-type:jwt-bearer"

	// GrantTypeRefreshToken is the RFC 6749 §6 grant type string used
	// for rotating user_oauth access tokens.
	GrantTypeRefreshToken = "refresh_token"

	// OAuthAPIBetaHeader is the anthropic-beta value required on
	// authenticated API requests using an OAuth bearer token, and on
	// refresh_token grants against the token endpoint.
	OAuthAPIBetaHeader = "oauth-2025-04-20"

	// FederationBetaHeader is the anthropic-beta value required on
	// jwt-bearer exchanges against the token endpoint. It routes the
	// request to the Go userauth service; it must NOT be sent on
	// refresh_token grants, which are gateway-routed to the Python
	// oauth-server.
	FederationBetaHeader = "oidc-federation-2026-04-01"
)

const (
	defaultAPIBaseURL = "https://api.anthropic.com"
	// federationExchangeBetaValue is the combined anthropic-beta value
	// for jwt-bearer exchanges: oauth-2025-04-20 unlocks the oauth/token
	// endpoint family, and oidc-federation-2026-04-01 routes jwt-bearer
	// to the federation service.
	federationExchangeBetaValue = OAuthAPIBetaHeader + "," + FederationBetaHeader
)

// FederationExchangeParams captures the inputs needed to exchange a signed
// third-party assertion (GitHub OIDC, Kubernetes service account token,
// etc.) for a short-lived Anthropic access token via the jwt-bearer grant.
type FederationExchangeParams struct {
	// Assertion is the signed JWT presented to the token endpoint. Required.
	Assertion string

	// FederationRuleID is the tagged ID ("fdrl_...") of the OidcFederationRule
	// that governs the exchange. Required.
	FederationRuleID string

	// OrganizationID is the tagged ID of the Anthropic organization whose
	// credentials the exchange should mint. Required.
	OrganizationID string

	// ServiceAccountID is an optional "svac_..." target check for
	// federation rules with target_type=SERVICE_ACCOUNT. Leave empty for
	// user-targeted rules.
	ServiceAccountID string

	// BaseURL overrides the Anthropic API base URL. Defaults to
	// https://api.anthropic.com. A trailing slash is tolerated.
	BaseURL string

	// HTTPClient overrides the default HTTP client used for the exchange.
	// When nil, a client with a 30s timeout is used.
	HTTPClient *http.Client

	// UserAgent overrides the outgoing User-Agent header. When empty, the
	// helper sends "anthropic-sdk-go/<version> ExchangeFederationAssertion"
	// so the token endpoint's access logs identify the caller for
	// incident triage. Callers with their own tooling (e.g. `ant-cli/1.2.3`)
	// should set this so support tickets can point at the real binary.
	UserAgent string
}

// federationExchangeRequest mirrors the JSON body the server's token
// endpoint accepts. The REST gateway's AliasTransformation strips any
// field not in this exact set, so each JSON tag is load-bearing — sending
// "federation_rule" or "organization" (without "_id") causes the gateway
// to drop the value and the Connect handler to fail with an empty
// federation_rule_id error. Keep in sync with
// internal/auth/workload.go:tokenExchangeRequest.
type federationExchangeRequest struct {
	GrantType        string `json:"grant_type"`
	Assertion        string `json:"assertion"`
	FederationRuleID string `json:"federation_rule_id"`
	OrganizationID   string `json:"organization_id"`
	ServiceAccountID string `json:"service_account_id,omitempty"`
}

// FederationExchangeError is returned by [ExchangeFederationAssertion] when
// the token endpoint responds with a non-2xx status. The server body is
// kept verbatim so the caller can surface the exact upstream message; the
// Request-Id header is captured separately so support tickets can include
// a correlation identifier.
type FederationExchangeError struct {
	StatusCode int
	Body       string
	RequestID  string
}

func (e *FederationExchangeError) Error() string {
	redacted := RedactOAuthErrorBody(e.Body)
	if e.RequestID != "" {
		return fmt.Sprintf("federation exchange failed (status %d, request-id %s): %s",
			e.StatusCode, e.RequestID, redacted)
	}
	return fmt.Sprintf("federation exchange failed (status %d): %s", e.StatusCode, redacted)
}

// ExchangeFederationAssertion performs an OAuth 2.0 jwt-bearer exchange
// against the Anthropic token endpoint and returns the minted credentials.
//
// Federation grants do not return a refresh token — callers re-exchange
// their assertion on expiry. The returned [*Credentials] therefore has an
// empty RefreshToken but a populated ExpiresAt whenever the server returns
// an expires_in field.
//
// On non-2xx responses the returned error is a [*FederationExchangeError]
// carrying the server body and the Request-Id response header.
func ExchangeFederationAssertion(ctx context.Context, params FederationExchangeParams) (*Credentials, error) {
	if params.Assertion == "" {
		return nil, errors.New("ExchangeFederationAssertion: Assertion is required")
	}
	if params.FederationRuleID == "" {
		return nil, errors.New("ExchangeFederationAssertion: FederationRuleID is required")
	}
	if params.OrganizationID == "" {
		return nil, errors.New("ExchangeFederationAssertion: OrganizationID is required")
	}

	base := strings.TrimRight(params.BaseURL, "/")
	if base == "" {
		base = defaultAPIBaseURL
	}

	bodyJSON, err := json.Marshal(federationExchangeRequest{
		GrantType:        GrantTypeJWTBearer,
		Assertion:        params.Assertion,
		FederationRuleID: params.FederationRuleID,
		OrganizationID:   params.OrganizationID,
		ServiceAccountID: params.ServiceAccountID,
	})
	if err != nil {
		return nil, fmt.Errorf("ExchangeFederationAssertion: marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base+TokenEndpoint, bytes.NewReader(bodyJSON))
	if err != nil {
		return nil, fmt.Errorf("ExchangeFederationAssertion: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-beta", federationExchangeBetaValue)
	ua := params.UserAgent
	if ua == "" {
		ua = "anthropic-sdk-go/" + internal.PackageVersion + " ExchangeFederationAssertion"
	}
	req.Header.Set("User-Agent", ua)

	client := params.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ExchangeFederationAssertion: post: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("ExchangeFederationAssertion: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &FederationExchangeError{
			StatusCode: resp.StatusCode,
			Body:       string(body),
			RequestID:  resp.Header.Get("Request-Id"),
		}
	}

	var parsed struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type,omitempty"`
		ExpiresIn   *int   `json:"expires_in,omitempty"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, fmt.Errorf("ExchangeFederationAssertion: parse response: %w", err)
	}
	if parsed.AccessToken == "" {
		return nil, errors.New("ExchangeFederationAssertion: response missing access_token")
	}
	if parsed.TokenType != "" && !strings.EqualFold(parsed.TokenType, "Bearer") {
		return nil, fmt.Errorf("ExchangeFederationAssertion: unsupported token_type %q (want Bearer)", parsed.TokenType)
	}

	creds := &Credentials{AccessToken: parsed.AccessToken}
	if parsed.ExpiresIn != nil {
		exp := time.Now().Add(time.Duration(*parsed.ExpiresIn) * time.Second)
		creds.ExpiresAt = &exp
	}
	return creds, nil
}
