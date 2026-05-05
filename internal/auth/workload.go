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
	// WorkspaceID is an optional `wrkspc_*` tagged ID, or the literal
	// "default" to scope the token to the organization's default workspace.
	// When omitted the server picks the rule's sole enabled workspace, else
	// the org default if the rule covers it. Required when the rule enables
	// more than one non-default workspace, or to target a specific workspace
	// other than the one the server would pick. The minted token is
	// workspace-scoped: per-request workspace selection (the
	// anthropic-workspace-id header) is not supported for federation
	// tokens — switching workspaces requires a new token exchange with a
	// different WorkspaceID.
	WorkspaceID string
	BaseURL     string // optional override; if empty, uses baseURL from TokenCache
}

type tokenExchangeRequest struct {
	GrantType        string `json:"grant_type"`
	Assertion        string `json:"assertion"`
	FederationRuleID string `json:"federation_rule_id"`
	OrganizationID   string `json:"organization_id"`
	ServiceAccountID string `json:"service_account_id,omitempty"`
	WorkspaceID      string `json:"workspace_id,omitempty"`
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
			WorkspaceID:      cfg.WorkspaceID,
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
			// A 401 is the auth-failure case worth a hint: point the operator
			// at the federation rule and the authentication-events log in
			// Claude Console. When no WorkspaceID is configured, also surface
			// the most common ambiguous-401 cause — a federation rule spanning
			// multiple workspaces — so the operator doesn't have to dig
			// through docs. Other statuses (5xx, non-401 4xx) get no hint:
			// they don't indicate a config problem this guidance would fix.
			var hint string
			if resp.StatusCode == http.StatusUnauthorized {
				hintParts := []string{
					"Ensure your federation rule matches your identity token",
				}
				if cfg.WorkspaceID == "" {
					hintParts = append(hintParts,
						"If your federation rule is scoped to multiple workspaces, set the "+
							"ANTHROPIC_WORKSPACE_ID environment variable, the 'workspace_id' "+
							"config key, or the WorkspaceID field on option.FederationOptions")
				}
				hintParts = append(hintParts,
					"View your authentication events in the Workload identity page of "+
						"Claude Console for more details")
				hint = strings.Join(hintParts, ". ") + "."
			}
			return nil, &OAuthTokenError{
				StatusCode:       resp.StatusCode,
				Body:             string(respBody),
				RequestID:        resp.Header.Get("Request-Id"),
				Hint:             hint,
				WorkloadIdentity: true,
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
