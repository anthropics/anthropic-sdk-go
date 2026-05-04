package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal"
)

// MaxAssertionSize bounds the JWT sent to /v1/oauth/token. Honest OIDC
// tokens are well under 16 KiB; reject larger inputs before shipping them
// at request rate.
const MaxAssertionSize = 16384

// OIDCFederationConfig configures an OIDC-federation [TokenProvider] that
// exchanges a third-party JWT for an Anthropic access token.
type OIDCFederationConfig struct {
	IdentityProvider IdentityTokenProvider
	FederationRuleID string
	OrganizationID   string
	ServiceAccountID string // optional
	BaseURL          string // optional override; if empty, uses baseURL from TokenCache
}

type tokenExchangeRequest struct {
	GrantType        string `json:"grant_type"`
	Assertion        string `json:"assertion"`
	FederationRuleID string `json:"federation_rule_id"`
	OrganizationID   string `json:"organization_id"`
	ServiceAccountID string `json:"service_account_id,omitempty"`
}

type tokenExchangeResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type,omitempty"`
	ExpiresIn   *int   `json:"expires_in,omitempty"`
}

// NewOIDCFederationCredentials returns a [TokenProvider] that exchanges an
// OIDC identity token for a short-lived Anthropic access token via the OAuth
// token endpoint.
func NewOIDCFederationCredentials(cfg OIDCFederationConfig) TokenProvider {
	return func(ctx context.Context, baseURL string, handler func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		effectiveBase := strings.TrimRight(cfg.BaseURL, "/")
		if effectiveBase == "" {
			effectiveBase = strings.TrimRight(baseURL, "/")
		}
		if err := requireSecureTokenEndpoint(effectiveBase); err != nil {
			return nil, err
		}
		jwt, err := cfg.IdentityProvider.GetIdentityToken(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get identity token: %w", err)
		}
		// post-fetch: bounds wire size, not allocation
		if len(jwt) > MaxAssertionSize {
			return nil, fmt.Errorf("identity token exceeds %d-byte limit (got %d bytes)", MaxAssertionSize, len(jwt))
		}

		body := tokenExchangeRequest{
			GrantType:        GrantTypeJWTBearer,
			Assertion:        jwt,
			FederationRuleID: cfg.FederationRuleID,
			OrganizationID:   cfg.OrganizationID,
			ServiceAccountID: cfg.ServiceAccountID,
		}
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal token exchange request: %w", err)
		}

		endpoint := effectiveBase + TokenEndpoint
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(bodyJSON))
		if err != nil {
			return nil, fmt.Errorf("failed to create token exchange request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "anthropic-sdk-go/"+internal.PackageVersion+" (oidc-federation)")
		// oauth-2025-04-20 unlocks the oauth/token endpoint family and is
		// required alongside the federation beta for jwt-bearer grants.
		req.Header.Set("anthropic-beta", OAuthAPIBetaHeader+","+FederationBetaHeader)

		resp, err := handler(req)
		if err != nil {
			return nil, fmt.Errorf("token exchange request failed: %w", err)
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		if err != nil {
			return nil, fmt.Errorf("failed to read token exchange response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, &OAuthTokenError{
				StatusCode: resp.StatusCode,
				Body:       string(respBody),
				RequestID:  resp.Header.Get("Request-Id"),
			}
		}

		var result tokenExchangeResponse
		if err := json.Unmarshal(respBody, &result); err != nil {
			return nil, fmt.Errorf("failed to parse token exchange response: %w", err)
		}

		if result.AccessToken == "" {
			return nil, fmt.Errorf("token exchange response missing access_token")
		}
		if result.TokenType != "" && !strings.EqualFold(result.TokenType, "Bearer") {
			return nil, fmt.Errorf("token exchange response: unsupported token_type %q (want Bearer)", result.TokenType)
		}

		token := &AccessToken{Token: result.AccessToken}
		if result.ExpiresIn != nil {
			exp := time.Now().Add(time.Duration(*result.ExpiresIn) * time.Second)
			token.ExpiresAt = &exp
		}
		return token, nil
	}
}
