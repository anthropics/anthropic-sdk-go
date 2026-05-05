package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go/config"
	"github.com/anthropics/anthropic-sdk-go/internal"
)

// Credentials file schema — credentials/<name>.json. Reads only; all writes
// go through config.WriteCredentials for a single-source-of-truth on the
// atomic + fsync + Chmod guarantees.
type credentialsTokenData struct {
	Type         string `json:"type"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    *int64 `json:"expires_at,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func readCredentialsFile(path string) (*credentialsTokenData, error) {
	if err := checkCredentialsFileSafety(path); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cred credentialsTokenData
	if err := json.Unmarshal(data, &cred); err != nil {
		return nil, err
	}
	return &cred, nil
}

// checkCredentialsFileSafety refuses credentials files exposed to other
// UIDs on the host. Symlinks could redirect a refresh-token write; any
// mode&0o077 bit lets another UID read or inject a token. Mode bits are
// not meaningful on Windows, so that check is skipped there.
func checkCredentialsFileSafety(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode().Type()&fs.ModeSymlink != 0 {
		return fmt.Errorf("credentials file %q is a symlink; refusing to read (set ANTHROPIC_CREDENTIALS_PATH to the real file)", path)
	}
	if runtime.GOOS == "windows" {
		return nil
	}
	mode := info.Mode().Perm()
	if mode&0o077 != 0 {
		return fmt.Errorf("credentials file %q has unsafe permissions %#o; refusing to read (must not be accessible by group or other; chmod 600 %s)", path, mode, path)
	}
	return nil
}

// freshAccessToken returns a non-nil [AccessToken] only when the credentials
// contain a non-empty token that won't expire within [MandatoryRefreshThreshold].
// If ExpiresAt is nil the token is treated as non-expiring.
func (cred *credentialsTokenData) freshAccessToken() *AccessToken {
	if cred.AccessToken == "" {
		return nil
	}
	if cred.ExpiresAt != nil {
		exp := time.Unix(*cred.ExpiresAt, 0)
		if time.Until(exp) <= MandatoryRefreshThreshold {
			return nil
		}
		return &AccessToken{Token: cred.AccessToken, ExpiresAt: &exp}
	}
	return &AccessToken{Token: cred.AccessToken}
}

// ResolveCredentials builds a [CredentialsResult] from a [config.Config].
// The [AuthenticationInfo.CredentialsPath] must already be resolved to an
// absolute path (or left empty if credentials are not file-based).
func ResolveCredentials(cfg *config.Config) (*CredentialsResult, error) {
	if cfg.AuthenticationInfo == nil {
		return nil, &CredentialResolutionError{
			Message: "config is missing authentication",
		}
	}
	switch cfg.AuthenticationInfo.Type {
	case config.AuthenticationTypeOIDCFederation:
		return loadOIDCFederationProfile(cfg)
	case config.AuthenticationTypeUserOAuth:
		return loadUserOAuthProfile(cfg)
	default:
		return nil, &CredentialResolutionError{
			Message: fmt.Sprintf("unknown authentication.type %q", cfg.AuthenticationInfo.Type),
		}
	}
}

func loadOIDCFederationProfile(cfg *config.Config) (*CredentialsResult, error) {
	oidc := cfg.AuthenticationInfo.OIDCFederation
	if oidc == nil {
		return nil, &CredentialResolutionError{
			Message: "oidc_federation config missing oidc_federation sub-object",
		}
	}
	if oidc.FederationRuleID == "" || cfg.OrganizationID == "" {
		return nil, &CredentialResolutionError{
			Message: "oidc_federation config missing federation_rule_id or organization_id",
		}
	}

	var identityProvider IdentityTokenProvider
	switch {
	case oidc.IdentityToken != nil:
		if oidc.IdentityToken.Source != config.IdentityTokenSourceFile {
			return nil, &CredentialResolutionError{
				Message: fmt.Sprintf("oidc_federation identity_token.source %q is not supported (only %q)", oidc.IdentityToken.Source, config.IdentityTokenSourceFile),
			}
		}
		if oidc.IdentityToken.Path == "" {
			return nil, &CredentialResolutionError{
				Message: fmt.Sprintf("oidc_federation identity_token.source %q requires a non-empty path", config.IdentityTokenSourceFile),
			}
		}
		identityProvider = &IdentityTokenFile{Path: oidc.IdentityToken.Path}
	default:
		identityProvider = identityProviderFromEnv()
		if identityProvider == nil {
			return nil, &CredentialResolutionError{
				Message: fmt.Sprintf("oidc_federation config requires identity_token with source %q, or %s / %s environment variable", config.IdentityTokenSourceFile, EnvIdentityTokenFile, EnvIdentityToken),
			}
		}
	}

	exchangeProvider := NewOIDCFederationCredentials(OIDCFederationConfig{
		IdentityProvider: identityProvider,
		FederationRuleID: oidc.FederationRuleID,
		OrganizationID:   cfg.OrganizationID,
		ServiceAccountID: oidc.ServiceAccountID,
		WorkspaceID:      cfg.WorkspaceID,
		BaseURL:          cfg.BaseURL,
	})

	credPath := cfg.AuthenticationInfo.CredentialsPath

	// Wrap with cache-aware provider that checks the credentials file first.
	provider := func(ctx context.Context, baseURL string, handler func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		// Try cached credentials file, unless the caller signaled a force-
		// refresh (e.g. after a 401 invalidation in the auth middleware).
		if !isForceRefresh(ctx) {
			if cred, err := readCredentialsFile(credPath); err == nil {
				if token := cred.freshAccessToken(); token != nil {
					return token, nil
				}
			} else if !os.IsNotExist(err) {
				warnOnce("workload-cache-read:"+credPath,
					"failed to read workload-identity token cache %q: %v (continuing with fresh exchange)", credPath, err)
			}
		}

		// Exchange for a new token.
		token, err := exchangeProvider(ctx, baseURL, handler)
		if err != nil {
			return nil, err
		}

		// Write cache back (best-effort; log once on failure so a misconfigured
		// cache directory does not cause a silent exchange-per-request loop).
		if token.ExpiresAt != nil {
			cacheData := config.Credentials{
				AccessToken: token.Token,
				ExpiresAt:   token.ExpiresAt,
			}
			if writeErr := config.WriteCredentials(credPath, cacheData); writeErr != nil {
				warnOnce("workload-cache-write:"+credPath,
					"failed to write workload-identity token cache %q: %v", credPath, writeErr)
			}
		}

		return token, nil
	}

	// For federation profiles workspace_id is sent in the jwt-bearer
	// exchange body, not as a request header (the minted token is already
	// workspace-scoped, so the header would be ignored). WorkspaceID is
	// therefore intentionally omitted from the CredentialsResult here.
	return &CredentialsResult{
		Provider: provider,
		BaseURL:  cfg.BaseURL,
	}, nil
}

func loadUserOAuthProfile(cfg *config.Config) (*CredentialsResult, error) {
	userOAuth := cfg.AuthenticationInfo.UserOAuth
	if userOAuth == nil {
		return nil, &CredentialResolutionError{
			Message: "user_oauth config missing user_oauth sub-object",
		}
	}

	credPath := cfg.AuthenticationInfo.CredentialsPath
	cred, err := readCredentialsFile(credPath)
	if err != nil {
		return nil, &CredentialResolutionError{
			Message: fmt.Sprintf("failed to read credentials file %q", credPath),
			Err:     err,
		}
	}

	if cred.AccessToken == "" {
		return nil, &CredentialResolutionError{
			Message: fmt.Sprintf("user_oauth credentials in %q missing access_token", credPath),
		}
	}

	if userOAuth.ClientID == "" {
		// Without a client_id we cannot refresh. Fail fast on an already-
		// expired token so callers get a clean error instead of a 401-retry
		// loop until maxRetries on a silently stale bearer.
		if cred.freshAccessToken() == nil {
			return nil, &CredentialResolutionError{
				Message: fmt.Sprintf("user_oauth credentials in %q: access_token is expired or missing and no refresh is available (client_id is empty)", credPath),
			}
		}
		// Re-read on every call so externally-rotated tokens (e.g. from a
		// credential daemon) are picked up.
		staticProvider := func(_ context.Context, _ string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
			current, err := readCredentialsFile(credPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read credentials file %q: %w", credPath, err)
			}
			token := current.freshAccessToken()
			if token == nil {
				return nil, fmt.Errorf("user_oauth credentials in %q: access_token is expired or missing and no refresh is available (client_id is empty)", credPath)
			}
			return token, nil
		}
		return &CredentialsResult{
			Provider:    staticProvider,
			BaseURL:     cfg.BaseURL,
			WorkspaceID: cfg.WorkspaceID,
		}, nil
	}

	if cred.RefreshToken == "" {
		return nil, &CredentialResolutionError{
			Message: fmt.Sprintf("user_oauth credentials in %q has client_id in config but missing refresh_token", credPath),
		}
	}

	provider := func(ctx context.Context, baseURL string, handler func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		base := strings.TrimRight(cfg.BaseURL, "/")
		if base == "" {
			base = strings.TrimRight(baseURL, "/")
		}
		if err := requireSecureTokenEndpoint(base); err != nil {
			return nil, err
		}
		// Re-read credentials file to check if token is still fresh. Skip
		// the freshness shortcut when force-refresh is signaled (e.g. after
		// a 401 in the auth middleware).
		current, err := readCredentialsFile(credPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read credentials file %q: %w", credPath, err)
		}
		if !isForceRefresh(ctx) {
			if token := current.freshAccessToken(); token != nil {
				return token, nil
			}
		}

		// Refresh the token. `scope` is intentionally omitted: the Python
		// oauth-server preserves the original token's full scope set when
		// scope is absent on the request (see the early-return in
		// api/api/oauth_server/oauth/grants/refresh_token_grant.py).
		refreshReq := map[string]string{
			"grant_type":    GrantTypeRefreshToken,
			"refresh_token": current.RefreshToken,
			"client_id":     userOAuth.ClientID,
		}
		refreshBody, _ := json.Marshal(refreshReq)

		endpoint := base + TokenEndpoint
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(refreshBody))
		if err != nil {
			return nil, fmt.Errorf("failed to create refresh request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "anthropic-sdk-go/"+internal.PackageVersion+" (user-oauth-refresh)")
		// The oauth-2025-04-20 beta is required by the Python oauth-server
		// to accept refresh-token grants for user_oauth credentials.
		//
		// Do NOT add FederationBetaHeader here. POST /v1/oauth/token is
		// gateway-routed: with oidc-federation-2026-04-01 present it goes
		// to the Go userauth service, which only handles the jwt-bearer
		// grant. Refresh-token grants must fall through to the Python
		// oauth-server, where the federation beta must be absent.
		req.Header.Set("anthropic-beta", OAuthAPIBetaHeader)

		resp, err := handler(req)
		if err != nil {
			return nil, fmt.Errorf("refresh request failed: %w", err)
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		if err != nil {
			return nil, fmt.Errorf("failed to read refresh response: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, &OAuthTokenError{
				StatusCode: resp.StatusCode,
				Body:       string(respBody),
				RequestID:  resp.Header.Get("Request-Id"),
			}
		}

		var result struct {
			AccessToken  string `json:"access_token"`
			TokenType    string `json:"token_type,omitempty"`
			RefreshToken string `json:"refresh_token"`
			ExpiresIn    *int   `json:"expires_in,omitempty"`
		}
		if err := json.Unmarshal(respBody, &result); err != nil {
			return nil, fmt.Errorf("failed to parse refresh response: %w", err)
		}
		if result.TokenType != "" && !strings.EqualFold(result.TokenType, "Bearer") {
			return nil, fmt.Errorf("refresh response: unsupported token_type %q (want Bearer)", result.TokenType)
		}

		updated := config.Credentials{
			AccessToken:  result.AccessToken,
			RefreshToken: current.RefreshToken,
		}
		if result.RefreshToken != "" {
			updated.RefreshToken = result.RefreshToken
		}
		token := &AccessToken{Token: result.AccessToken}
		if result.ExpiresIn != nil {
			exp := time.Now().Add(time.Duration(*result.ExpiresIn) * time.Second)
			updated.ExpiresAt = &exp
			token.ExpiresAt = &exp
		}

		if err := config.WriteCredentials(credPath, updated); err != nil {
			return nil, fmt.Errorf("failed to write updated credentials: %w", err)
		}
		return token, nil
	}

	return &CredentialsResult{
		Provider:    provider,
		BaseURL:     cfg.BaseURL,
		WorkspaceID: cfg.WorkspaceID,
	}, nil
}
