package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/config"
)

func writeCredentials(t *testing.T, path string, data map[string]any) {
	t.Helper()
	os.MkdirAll(filepath.Dir(path), 0700)
	b, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile(path, b, 0600)
}

func TestResolveCredentials_OIDCFederation(t *testing.T) {
	var receivedBody tokenExchangeRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			t.Fatal(err)
		}
		expiresIn := 3600
		json.NewEncoder(w).Encode(tokenExchangeResponse{
			AccessToken: "exchanged-tok",
			ExpiresIn:   &expiresIn,
		})
	}))
	defer server.Close()

	dir := t.TempDir()
	tokenPath := filepath.Join(dir, "identity-token")
	os.WriteFile(tokenPath, []byte("my-jwt"), 0600)

	result, err := ResolveCredentials(&config.Config{
		OrganizationID: "org-1",
		WorkspaceID:    "wrkspc_x",
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeOIDCFederation,
			CredentialsPath: filepath.Join(dir, "credentials", "creds.json"),
			OIDCFederation: &config.OIDCFederation{
				FederationRuleID: "fdrl_1",
				IdentityToken: &config.IdentityTokenConfig{
					Source: config.IdentityTokenSourceFile,
					Path:   tokenPath,
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	// For federation profiles workspace_id is sent in the exchange body, not
	// as a request header — the result must not surface a header value.
	if result.WorkspaceID != "" {
		t.Fatalf("federation result should not carry a header workspace_id, got %q", result.WorkspaceID)
	}
	tok, err := result.Provider(context.Background(), server.URL, http.DefaultClient.Do)
	if err != nil {
		t.Fatal(err)
	}
	if tok.Token != "exchanged-tok" {
		t.Fatalf("got %q, want %q", tok.Token, "exchanged-tok")
	}
	if receivedBody.WorkspaceID != "wrkspc_x" {
		t.Fatalf("got exchange-body workspace_id %q, want %q", receivedBody.WorkspaceID, "wrkspc_x")
	}
}

func TestResolveCredentials_OIDCFederationWithCache(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not call token endpoint when cache is fresh")
	}))
	defer server.Close()

	credPath := filepath.Join(t.TempDir(), "creds.json")
	freshAt := time.Now().Add(10 * time.Minute).Unix()
	writeCredentials(t, credPath, map[string]any{
		"type":         "oauth_token",
		"access_token": "cached-tok",
		"expires_at":   freshAt,
	})

	result, err := ResolveCredentials(&config.Config{
		OrganizationID: "org-1",
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeOIDCFederation,
			CredentialsPath: credPath,
			OIDCFederation: &config.OIDCFederation{
				FederationRuleID: "fdrl_1",
				IdentityToken: &config.IdentityTokenConfig{
					Source: config.IdentityTokenSourceFile,
					Path:   "/dev/null",
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	tok, err := result.Provider(context.Background(), server.URL, http.DefaultClient.Do)
	if err != nil {
		t.Fatal(err)
	}
	if tok.Token != "cached-tok" {
		t.Fatalf("got %q, want %q", tok.Token, "cached-tok")
	}
}

func TestResolveCredentials_OIDCFederationMissingOrganization(t *testing.T) {
	dir := t.TempDir()
	tokenPath := filepath.Join(dir, "token")
	os.WriteFile(tokenPath, []byte("jwt"), 0600)

	_, err := ResolveCredentials(&config.Config{
		// missing organization_id
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeOIDCFederation,
			CredentialsPath: filepath.Join(dir, "creds.json"),
			OIDCFederation: &config.OIDCFederation{
				FederationRuleID: "fdrl_1",
				IdentityToken: &config.IdentityTokenConfig{
					Source: config.IdentityTokenSourceFile,
					Path:   tokenPath,
				},
			},
		},
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestResolveCredentials_OIDCFederationMissingFederationRuleID(t *testing.T) {
	dir := t.TempDir()
	tokenPath := filepath.Join(dir, "token")
	os.WriteFile(tokenPath, []byte("jwt"), 0600)

	_, err := ResolveCredentials(&config.Config{
		OrganizationID: "org-1",
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeOIDCFederation,
			CredentialsPath: filepath.Join(dir, "creds.json"),
			OIDCFederation: &config.OIDCFederation{
				// missing federation_rule_id
				IdentityToken: &config.IdentityTokenConfig{
					Source: config.IdentityTokenSourceFile,
					Path:   tokenPath,
				},
			},
		},
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestResolveCredentials_UserOAuthRefresh(t *testing.T) {
	refreshCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshCalls++
		beta := r.Header.Get("anthropic-beta")
		if strings.Contains(beta, "oidc-federation-2026-04-01") {
			t.Errorf("refresh request must not carry oidc-federation beta header, got %q", beta)
		}
		if !strings.Contains(beta, "oauth-2025-04-20") {
			t.Errorf("refresh request missing required oauth beta header, got %q", beta)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if _, ok := body["client_secret"]; ok {
			t.Fatal("client_secret should not be sent in refresh request")
		}
		expiresIn := 3600
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "refreshed-tok",
			"refresh_token": "new-refresh-tok",
			"expires_in":    expiresIn,
		})
	}))
	defer server.Close()

	credPath := filepath.Join(t.TempDir(), "creds.json")
	expiredAt := time.Now().Add(-1 * time.Minute).Unix()
	writeCredentials(t, credPath, map[string]any{
		"type":          "oauth_token",
		"access_token":  "old-tok",
		"expires_at":    expiredAt,
		"refresh_token": "old-refresh",
	})

	result, err := ResolveCredentials(&config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeUserOAuth,
			CredentialsPath: credPath,
			UserOAuth: &config.UserOAuth{
				ClientID: "my-client",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	tok, err := result.Provider(context.Background(), server.URL, http.DefaultClient.Do)
	if err != nil {
		t.Fatal(err)
	}
	if tok.Token != "refreshed-tok" {
		t.Fatalf("got %q, want %q", tok.Token, "refreshed-tok")
	}
	if refreshCalls != 1 {
		t.Fatalf("expected 1 refresh call, got %d", refreshCalls)
	}

	// Verify the credentials file was updated.
	data, _ := os.ReadFile(credPath)
	var updated credentialsTokenData
	json.Unmarshal(data, &updated)
	if updated.RefreshToken != "new-refresh-tok" {
		t.Fatalf("credentials refresh_token not updated, got %q", updated.RefreshToken)
	}
	if updated.Type != "oauth_token" {
		t.Fatalf("credentials type not set, got %q", updated.Type)
	}
}

func TestResolveCredentials_UserOAuthFreshNoRefresh(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not call token endpoint")
	}))
	defer server.Close()

	credPath := filepath.Join(t.TempDir(), "creds.json")
	freshAt := time.Now().Add(10 * time.Minute).Unix()
	writeCredentials(t, credPath, map[string]any{
		"type":          "oauth_token",
		"access_token":  "fresh-tok",
		"expires_at":    freshAt,
		"refresh_token": "refresh",
	})

	result, err := ResolveCredentials(&config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeUserOAuth,
			CredentialsPath: credPath,
			UserOAuth: &config.UserOAuth{
				ClientID: "my-client",
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	tok, err := result.Provider(context.Background(), server.URL, http.DefaultClient.Do)
	if err != nil {
		t.Fatal(err)
	}
	if tok.Token != "fresh-tok" {
		t.Fatalf("got %q, want %q", tok.Token, "fresh-tok")
	}
}

func TestResolveCredentials_UserOAuthNoClientIDStaticToken(t *testing.T) {
	credPath := filepath.Join(t.TempDir(), "creds.json")
	writeCredentials(t, credPath, map[string]any{
		"type":         "oauth_token",
		"access_token": "static-tok",
	})

	result, err := ResolveCredentials(&config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeUserOAuth,
			CredentialsPath: credPath,
			UserOAuth:       &config.UserOAuth{},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	tok, err := result.Provider(context.Background(), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if tok.Token != "static-tok" {
		t.Fatalf("got %q, want %q", tok.Token, "static-tok")
	}
}

// TestResolveCredentials_UserOAuthNoClientIDReReads verifies that the static
// (no-client-id) user_oauth path re-reads the credentials file on every
// provider call so externally-rotated tokens (e.g. from a credential daemon)
// are picked up. This is the same behaviour the deleted "external" loader
// provided.
func TestResolveCredentials_UserOAuthNoClientIDReReads(t *testing.T) {
	credPath := filepath.Join(t.TempDir(), "creds.json")
	writeCredentials(t, credPath, map[string]any{
		"type":         "oauth_token",
		"access_token": "tok-1",
	})

	result, err := ResolveCredentials(&config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeUserOAuth,
			CredentialsPath: credPath,
			UserOAuth:       &config.UserOAuth{},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	tok, err := result.Provider(context.Background(), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if tok.Token != "tok-1" {
		t.Fatalf("got %q, want %q", tok.Token, "tok-1")
	}

	// Rotate the file underneath us and verify the next call sees it.
	writeCredentials(t, credPath, map[string]any{
		"type":         "oauth_token",
		"access_token": "tok-2",
	})

	tok, err = result.Provider(context.Background(), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if tok.Token != "tok-2" {
		t.Fatalf("got %q, want %q", tok.Token, "tok-2")
	}
}

// TestResolveCredentials_UserOAuthNoClientIDExpiredFailsFast verifies that
// when a user_oauth profile has no client_id (no refresh possible) and the
// on-disk token is already expired, ResolveCredentials returns an error
// upfront instead of handing an expired token back to the cache, which
// would otherwise 401-loop until maxRetries.
func TestResolveCredentials_UserOAuthNoClientIDExpiredFailsFast(t *testing.T) {
	credPath := filepath.Join(t.TempDir(), "creds.json")
	expiredAt := time.Now().Add(-time.Hour).Unix()
	writeCredentials(t, credPath, map[string]any{
		"type":         "oauth_token",
		"access_token": "stale-tok",
		"expires_at":   expiredAt,
	})

	_, err := ResolveCredentials(&config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeUserOAuth,
			CredentialsPath: credPath,
			UserOAuth:       &config.UserOAuth{},
		},
	})
	if err == nil {
		t.Fatal("expected error for expired static user_oauth token, got nil")
	}
	if !strings.Contains(err.Error(), "expired") {
		t.Errorf("error should mention expiry, got: %v", err)
	}
}

// TestResolveCredentials_OIDCFederationEnvIdentityToken verifies that a
// profile which omits identity_token falls back to ANTHROPIC_IDENTITY_TOKEN
// (the literal token env var) in addition to the _FILE variant.
func TestResolveCredentials_OIDCFederationEnvIdentityToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expiresIn := 3600
		json.NewEncoder(w).Encode(tokenExchangeResponse{
			AccessToken: "exchanged-from-env-literal",
			ExpiresIn:   &expiresIn,
		})
	}))
	defer server.Close()

	t.Setenv(EnvIdentityTokenFile, "")
	os.Unsetenv(EnvIdentityTokenFile)
	t.Setenv(EnvIdentityToken, "literal-jwt-xyz")

	dir := t.TempDir()
	result, err := ResolveCredentials(&config.Config{
		OrganizationID: "org-1",
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeOIDCFederation,
			CredentialsPath: filepath.Join(dir, "credentials", "creds.json"),
			OIDCFederation: &config.OIDCFederation{
				FederationRuleID: "fdrl_1",
			},
		},
	})
	if err != nil {
		t.Fatalf("expected env-literal fallback to succeed, got: %v", err)
	}
	tok, err := result.Provider(context.Background(), server.URL, http.DefaultClient.Do)
	if err != nil {
		t.Fatal(err)
	}
	if tok.Token != "exchanged-from-env-literal" {
		t.Fatalf("unexpected token: %q", tok.Token)
	}
}

func TestResolveCredentials_BaseURLFromConfig(t *testing.T) {
	credPath := filepath.Join(t.TempDir(), "creds.json")
	writeCredentials(t, credPath, map[string]any{
		"type":         "oauth_token",
		"access_token": "tok",
	})

	result, err := ResolveCredentials(&config.Config{
		BaseURL: "https://staging.anthropic.com",
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeUserOAuth,
			CredentialsPath: credPath,
			UserOAuth:       &config.UserOAuth{},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.BaseURL != "https://staging.anthropic.com" {
		t.Fatalf("got base_url %q, want %q", result.BaseURL, "https://staging.anthropic.com")
	}
}

func TestResolveCredentials_WorkspaceID(t *testing.T) {
	credPath := filepath.Join(t.TempDir(), "creds.json")
	writeCredentials(t, credPath, map[string]any{
		"type":         "oauth_token",
		"access_token": "tok",
	})

	result, err := ResolveCredentials(&config.Config{
		WorkspaceID: "wrkspc_01test",
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeUserOAuth,
			CredentialsPath: credPath,
			UserOAuth:       &config.UserOAuth{},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.WorkspaceID != "wrkspc_01test" {
		t.Fatalf("got workspace_id %q, want %q", result.WorkspaceID, "wrkspc_01test")
	}
}

func TestResolveCredentials_MissingAuthenticationInfo(t *testing.T) {
	_, err := ResolveCredentials(&config.Config{})
	if err == nil {
		t.Fatal("expected error for missing authentication")
	}
}

func TestResolveCredentials_UnknownType(t *testing.T) {
	_, err := ResolveCredentials(&config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{
			Type: "unknown_type",
		},
	})
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

// TestCheckCredentialsFileSafety_RejectsWorldReadable verifies that a
// credentials file with world-readable mode is refused. A bearer token
// must not be readable by another UID on the host.
func TestCheckCredentialsFileSafety_RejectsWorldReadable(t *testing.T) {
	credPath := filepath.Join(t.TempDir(), "creds.json")
	writeCredentials(t, credPath, map[string]any{"type": "oauth_token", "access_token": "tok"})
	if err := os.Chmod(credPath, 0644); err != nil {
		t.Fatal(err)
	}
	_, err := readCredentialsFile(credPath)
	if err == nil {
		t.Fatal("expected error for world-readable credentials file")
	}
	if !strings.Contains(err.Error(), "unsafe permissions") {
		t.Errorf("expected unsafe-permissions error, got: %v", err)
	}
}

// TestCheckCredentialsFileSafety_RejectsGroupReadable verifies that a
// credentials file with group-readable mode is also refused (matches TS
// behavior; Python warns but Go is stricter to match TS).
func TestCheckCredentialsFileSafety_RejectsGroupReadable(t *testing.T) {
	credPath := filepath.Join(t.TempDir(), "creds.json")
	writeCredentials(t, credPath, map[string]any{"type": "oauth_token", "access_token": "tok"})
	if err := os.Chmod(credPath, 0640); err != nil {
		t.Fatal(err)
	}
	_, err := readCredentialsFile(credPath)
	if err == nil {
		t.Fatal("expected error for group-readable credentials file")
	}
}

// TestCheckCredentialsFileSafety_RejectsWorldWritable verifies that a
// credentials file with world-writable mode is refused. A writable file
// lets an attacker inject a token the SDK then presents as the caller's
// identity — at least as dangerous as a readable one.
func TestCheckCredentialsFileSafety_RejectsWorldWritable(t *testing.T) {
	credPath := filepath.Join(t.TempDir(), "creds.json")
	writeCredentials(t, credPath, map[string]any{"type": "oauth_token", "access_token": "tok"})
	if err := os.Chmod(credPath, 0602); err != nil {
		t.Fatal(err)
	}
	_, err := readCredentialsFile(credPath)
	if err == nil {
		t.Fatal("expected error for world-writable credentials file")
	}
	if !strings.Contains(err.Error(), "unsafe permissions") {
		t.Errorf("expected unsafe-permissions error, got: %v", err)
	}
}

// TestCheckCredentialsFileSafety_AcceptsOwnerOnly verifies that 0600
// credentials read normally — the safety check is opt-in to risky modes.
func TestCheckCredentialsFileSafety_AcceptsOwnerOnly(t *testing.T) {
	credPath := filepath.Join(t.TempDir(), "creds.json")
	writeCredentials(t, credPath, map[string]any{"type": "oauth_token", "access_token": "tok"})
	if err := os.Chmod(credPath, 0600); err != nil {
		t.Fatal(err)
	}
	cred, err := readCredentialsFile(credPath)
	if err != nil {
		t.Fatalf("expected 0600 file to read cleanly, got: %v", err)
	}
	if cred.AccessToken != "tok" {
		t.Errorf("got access_token %q, want %q", cred.AccessToken, "tok")
	}
}

// TestCheckCredentialsFileSafety_RejectsSymlink verifies that a symlink
// at the credentials path is refused. An attacker-controlled symlink in
// the credentials dir could redirect a refresh-write to an arbitrary
// location.
func TestCheckCredentialsFileSafety_RejectsSymlink(t *testing.T) {
	dir := t.TempDir()
	realPath := filepath.Join(dir, "real.json")
	writeCredentials(t, realPath, map[string]any{"type": "oauth_token", "access_token": "tok"})
	if err := os.Chmod(realPath, 0600); err != nil {
		t.Fatal(err)
	}
	linkPath := filepath.Join(dir, "creds.json")
	if err := os.Symlink(realPath, linkPath); err != nil {
		t.Skipf("symlink unsupported on this platform: %v", err)
	}
	_, err := readCredentialsFile(linkPath)
	if err == nil {
		t.Fatal("expected error for symlink at credentials path")
	}
	if !strings.Contains(err.Error(), "symlink") {
		t.Errorf("expected symlink error, got: %v", err)
	}
}

// TestResolveCredentials_UserOAuthRefresh_SetsUserAgent verifies that the
// refresh-token POST carries an SDK User-Agent header. Mirrors the public
// ExchangeFederationAssertion behavior; previously absent on the internal
// refresh path.
func TestResolveCredentials_UserOAuthRefresh_SetsUserAgent(t *testing.T) {
	var receivedUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUA = r.Header.Get("User-Agent")
		expiresIn := 3600
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "new-tok",
			"refresh_token": "new-rt",
			"expires_in":    expiresIn,
		})
	}))
	defer server.Close()

	credPath := filepath.Join(t.TempDir(), "creds.json")
	expired := time.Now().Add(-time.Hour).Unix()
	writeCredentials(t, credPath, map[string]any{
		"type":          "oauth_token",
		"access_token":  "stale",
		"refresh_token": "rt",
		"expires_at":    expired,
	})
	if err := os.Chmod(credPath, 0600); err != nil {
		t.Fatal(err)
	}

	result, err := ResolveCredentials(&config.Config{
		BaseURL: server.URL,
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeUserOAuth,
			CredentialsPath: credPath,
			UserOAuth:       &config.UserOAuth{ClientID: "cid"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := result.Provider(context.Background(), "", http.DefaultClient.Do); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(receivedUA, "anthropic-sdk-go/") {
		t.Errorf("got User-Agent %q, want prefix anthropic-sdk-go/", receivedUA)
	}
	if !strings.Contains(receivedUA, "user-oauth-refresh") {
		t.Errorf("got User-Agent %q, want context user-oauth-refresh", receivedUA)
	}
}
