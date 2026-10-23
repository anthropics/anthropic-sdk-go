package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/config"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

const successResponse = `{"id":"msg_1","type":"message","role":"assistant","content":[{"type":"text","text":"hi"}],"model":"claude-sonnet-4-20250514","stop_reason":"end_turn","usage":{"input_tokens":1,"output_tokens":1}}`

// unsetEnv removes an env var for the duration of the test. t.Setenv
// registers the original value for restore at test end; Unsetenv then
// actually removes it for the duration of this test.
func unsetEnv(t *testing.T, key string) {
	t.Helper()
	t.Setenv(key, "")
	os.Unsetenv(key)
}

// isolateAuthEnv points ANTHROPIC_CONFIG_DIR at a fresh temp dir and
// unsets the federation env vars so a test is not affected by whatever
// host profile or env state the developer has configured.
func isolateAuthEnv(t *testing.T) {
	t.Helper()
	t.Setenv("ANTHROPIC_CONFIG_DIR", t.TempDir())
	unsetEnv(t, "ANTHROPIC_PROFILE")
	unsetEnv(t, "ANTHROPIC_FEDERATION_RULE_ID")
	unsetEnv(t, "ANTHROPIC_ORGANIZATION_ID")
	unsetEnv(t, "ANTHROPIC_IDENTITY_TOKEN")
	unsetEnv(t, "ANTHROPIC_IDENTITY_TOKEN_FILE")
	unsetEnv(t, "ANTHROPIC_WORKSPACE_ID")
}

var defaultParams = anthropic.MessageNewParams{
	Model:     "claude-sonnet-4-20250514",
	MaxTokens: 10,
	Messages: []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock("hi")),
	},
}

func TestIntegration_BearerTokenAndBetaHeader(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	var receivedAuth string
	var receivedBeta string

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		receivedBeta = r.Header.Get("anthropic-beta")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successResponse))
	}))
	defer apiServer.Close()

	dir := t.TempDir()
	credPath := filepath.Join(dir, "credentials.json")
	os.WriteFile(credPath, []byte(`{"type":"oauth_token","access_token":"my-access-tok"}`), 0600)

	client := anthropic.NewClient(
		option.WithBaseURL(apiServer.URL),
		option.WithConfig(&config.Config{
			AuthenticationInfo: &config.AuthenticationInfo{
				Type:            config.AuthenticationTypeUserOAuth,
				CredentialsPath: credPath,
				UserOAuth:       &config.UserOAuth{},
			},
		}),
	)

	_, err := client.Messages.New(context.Background(), defaultParams)
	if err != nil {
		t.Fatal(err)
	}

	if receivedAuth != "Bearer my-access-tok" {
		t.Fatalf("got Authorization %q, want %q", receivedAuth, "Bearer my-access-tok")
	}
	if !strings.Contains(receivedBeta, "oauth-2025-04-20") {
		t.Fatalf("anthropic-beta %q does not contain API oauth header", receivedBeta)
	}
	if strings.Contains(receivedBeta, "oidc-federation-2026-04-01") {
		t.Fatalf("anthropic-beta %q should not contain federation header on API requests", receivedBeta)
	}
}

func TestIntegration_WithFederationTokenProvider(t *testing.T) {
	var tokenExchangeBeta string
	var apiBeta string
	var apiAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/oauth/token" {
			tokenExchangeBeta = r.Header.Get("anthropic-beta")
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			if body["assertion"] != "custom-jwt" {
				t.Errorf("got assertion %v, want %q", body["assertion"], "custom-jwt")
			}
			if body["workspace_id"] != "wrkspc_x" {
				t.Errorf("got workspace_id %v, want %q", body["workspace_id"], "wrkspc_x")
			}
			json.NewEncoder(w).Encode(map[string]any{"access_token": "exchanged-tok"})
			return
		}
		apiBeta = r.Header.Get("anthropic-beta")
		apiAuth = r.Header.Get("Authorization")
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)
	t.Setenv("ANTHROPIC_BASE_URL", server.URL)

	var providerCalls int
	client := anthropic.NewClient(
		option.WithFederationTokenProvider(
			func(_ context.Context) (string, error) {
				providerCalls++
				return "custom-jwt", nil
			},
			option.FederationOptions{
				FederationRuleID: "rule-1",
				OrganizationID:   "org-1",
				WorkspaceID:      "wrkspc_x",
			},
		),
	)

	_, err := client.Messages.New(context.Background(), defaultParams)
	if err != nil {
		t.Fatal(err)
	}
	if providerCalls != 1 {
		t.Fatalf("expected provider called once, got %d", providerCalls)
	}
	if !strings.Contains(tokenExchangeBeta, "oidc-federation-2026-04-01") {
		t.Fatalf("token exchange beta %q missing federation header", tokenExchangeBeta)
	}
	if apiAuth != "Bearer exchanged-tok" {
		t.Fatalf("got api auth %q, want %q", apiAuth, "Bearer exchanged-tok")
	}
	if !strings.Contains(apiBeta, "oauth-2025-04-20") {
		t.Fatalf("api beta %q missing API oauth header", apiBeta)
	}
}

func TestIntegration_TokenCachedAcrossRequests(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	var tokenCalls atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/oauth/token" {
			tokenCalls.Add(1)
			expiresIn := 3600
			json.NewEncoder(w).Encode(map[string]any{
				"access_token": "cached-tok",
				"expires_in":   expiresIn,
			})
			return
		}
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	dir := t.TempDir()
	tokenPath := filepath.Join(dir, "token")
	os.WriteFile(tokenPath, []byte("my-jwt"), 0600)

	client := anthropic.NewClient(
		option.WithConfig(&config.Config{
			BaseURL:        server.URL,
			OrganizationID: "org-1",
			AuthenticationInfo: &config.AuthenticationInfo{
				Type:            config.AuthenticationTypeOIDCFederation,
				CredentialsPath: filepath.Join(dir, "credentials.json"),
				OIDCFederation: &config.OIDCFederation{
					FederationRuleID: "rule-1",
					IdentityToken: &config.IdentityTokenConfig{
						Source: config.IdentityTokenSourceFile,
						Path:   tokenPath,
					},
				},
			},
		}),
	)

	for range 3 {
		_, err := client.Messages.New(context.Background(), defaultParams)
		if err != nil {
			t.Fatal(err)
		}
	}

	if n := tokenCalls.Load(); n != 1 {
		t.Fatalf("expected 1 token exchange call, got %d", n)
	}
}

// TestIntegration_InMemoryFederationConfigSendsWorkspaceID verifies that an
// in-memory federation config passed via option.WithConfig forwards the
// top-level workspace_id into the jwt-bearer exchange body — and ONLY there.
// The anthropic-workspace-id header must be suppressed on API requests for
// federation profiles (the minted token is already workspace-scoped, so the
// header would be ignored). This mirrors the file-loaded path in
// loadOIDCFederationProfile, which omits WorkspaceID from CredentialsResult.
func TestIntegration_InMemoryFederationConfigSendsWorkspaceID(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	var exchangeBody map[string]any
	var apiWorkspaceHeader string
	var apiWorkspaceHeaderPresent bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/oauth/token" {
			_ = json.NewDecoder(r.Body).Decode(&exchangeBody)
			json.NewEncoder(w).Encode(map[string]any{"access_token": "exchanged-tok"})
			return
		}
		apiWorkspaceHeader = r.Header.Get("anthropic-workspace-id")
		_, apiWorkspaceHeaderPresent = r.Header[http.CanonicalHeaderKey("anthropic-workspace-id")]
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	dir := t.TempDir()
	tokenPath := filepath.Join(dir, "token")
	os.WriteFile(tokenPath, []byte("my-jwt"), 0600)

	client := anthropic.NewClient(
		option.WithConfig(&config.Config{
			BaseURL:        server.URL,
			OrganizationID: "org-1",
			WorkspaceID:    "wrkspc_x",
			AuthenticationInfo: &config.AuthenticationInfo{
				Type:            config.AuthenticationTypeOIDCFederation,
				CredentialsPath: filepath.Join(dir, "credentials.json"),
				OIDCFederation: &config.OIDCFederation{
					FederationRuleID: "rule-1",
					IdentityToken: &config.IdentityTokenConfig{
						Source: config.IdentityTokenSourceFile,
						Path:   tokenPath,
					},
				},
			},
		}),
	)
	if _, err := client.Messages.New(context.Background(), defaultParams); err != nil {
		t.Fatal(err)
	}
	if exchangeBody["workspace_id"] != "wrkspc_x" {
		t.Fatalf("got exchange-body workspace_id %v, want %q", exchangeBody["workspace_id"], "wrkspc_x")
	}
	if apiWorkspaceHeaderPresent {
		t.Fatalf("anthropic-workspace-id header should be suppressed for in-memory federation config, got %q", apiWorkspaceHeader)
	}
}

// TestIntegration_InMemoryNonFederationConfigSendsWorkspaceIDHeader is the
// regression counterpart to the federation test above: a non-federation
// in-memory config (user_oauth) with a top-level workspace_id must still
// emit the anthropic-workspace-id header on API requests. The federation
// suppression is gated on AuthenticationInfo.Type, not on whether the
// config is in-memory.
func TestIntegration_InMemoryNonFederationConfigSendsWorkspaceIDHeader(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	var apiWorkspaceHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiWorkspaceHeader = r.Header.Get("anthropic-workspace-id")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	dir := t.TempDir()
	credPath := filepath.Join(dir, "credentials.json")
	os.WriteFile(credPath, []byte(`{"type":"oauth_token","access_token":"my-access-tok"}`), 0600)

	client := anthropic.NewClient(
		option.WithConfig(&config.Config{
			BaseURL:     server.URL,
			WorkspaceID: "wrkspc_y",
			AuthenticationInfo: &config.AuthenticationInfo{
				Type:            config.AuthenticationTypeUserOAuth,
				CredentialsPath: credPath,
				UserOAuth:       &config.UserOAuth{},
			},
		}),
	)
	if _, err := client.Messages.New(context.Background(), defaultParams); err != nil {
		t.Fatal(err)
	}
	if apiWorkspaceHeader != "wrkspc_y" {
		t.Fatalf("got anthropic-workspace-id header %q, want %q", apiWorkspaceHeader, "wrkspc_y")
	}
}

// TestIntegration_EnvWorkspaceIDFillsUserOAuthHeader verifies that
// ANTHROPIC_WORKSPACE_ID fills workspace_id uniformly across profile types —
// not just federation. This pins the precedence model: explicit config >
// env var > profile, regardless of authentication.type. For user_oauth the
// filled value surfaces as the anthropic-workspace-id request header
// (federation routes it into the exchange body instead).
func TestIntegration_EnvWorkspaceIDFillsUserOAuthHeader(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	var apiWorkspaceHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiWorkspaceHeader = r.Header.Get("anthropic-workspace-id")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	dir := t.TempDir()
	configsDir := filepath.Join(dir, "configs")
	credsDir := filepath.Join(dir, "credentials")
	os.MkdirAll(configsDir, 0755)
	os.MkdirAll(credsDir, 0700)
	// user_oauth profile with NO workspace_id; ANTHROPIC_WORKSPACE_ID
	// fills the missing top-level config key.
	os.WriteFile(filepath.Join(configsDir, "default.json"), []byte(`{
		"base_url": "`+server.URL+`",
		"authentication": {"type": "user_oauth"}
	}`), 0644)
	os.WriteFile(filepath.Join(credsDir, "default.json"),
		[]byte(`{"type":"oauth_token","access_token":"tok"}`), 0600)
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_WORKSPACE_ID", "wrkspc_env")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.WorkspaceID != "wrkspc_env" {
		t.Fatalf("got cfg.WorkspaceID %q, want env fill-in", cfg.WorkspaceID)
	}

	client := anthropic.NewClient(option.WithConfig(cfg))
	if _, err := client.Messages.New(context.Background(), defaultParams); err != nil {
		t.Fatal(err)
	}
	if apiWorkspaceHeader != "wrkspc_env" {
		t.Fatalf("got anthropic-workspace-id header %q, want %q", apiWorkspaceHeader, "wrkspc_env")
	}
}

func TestIntegration_PerRequestWorkspaceIDOverridesConfigHeader(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	var apiWorkspaceHeaders []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiWorkspaceHeaders = r.Header.Values("anthropic-workspace-id")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	dir := t.TempDir()
	credPath := filepath.Join(dir, "credentials.json")
	os.WriteFile(credPath, []byte(`{"type":"oauth_token","access_token":"my-access-tok"}`), 0600)

	client := anthropic.NewClient(
		option.WithConfig(&config.Config{
			BaseURL:     server.URL,
			WorkspaceID: "wrkspc_config",
			AuthenticationInfo: &config.AuthenticationInfo{
				Type:            config.AuthenticationTypeUserOAuth,
				CredentialsPath: credPath,
				UserOAuth:       &config.UserOAuth{},
			},
		}),
	)
	params := defaultParams
	params.WorkspaceID = anthropic.String("wrkspc_request")
	if _, err := client.Messages.New(context.Background(), params); err != nil {
		t.Fatal(err)
	}
	if len(apiWorkspaceHeaders) != 1 || apiWorkspaceHeaders[0] != "wrkspc_request" {
		t.Fatalf("got anthropic-workspace-id header %q, want single %q", apiWorkspaceHeaders, "wrkspc_request")
	}
}

func TestIntegration_401RetryWithInvalidation(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	var apiCalls atomic.Int32
	var tokenCalls atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/oauth/token" {
			n := tokenCalls.Add(1)
			json.NewEncoder(w).Encode(map[string]any{
				"access_token": "tok-" + string(rune('0'+n)),
			})
			return
		}
		n := apiCalls.Add(1)
		if n == 1 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{}`))
			return
		}
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	dir := t.TempDir()
	tokenPath := filepath.Join(dir, "token")
	os.WriteFile(tokenPath, []byte("my-jwt"), 0600)

	client := anthropic.NewClient(
		option.WithConfig(&config.Config{
			BaseURL:        server.URL,
			OrganizationID: "org-1",
			AuthenticationInfo: &config.AuthenticationInfo{
				Type:            config.AuthenticationTypeOIDCFederation,
				CredentialsPath: filepath.Join(dir, "credentials.json"),
				OIDCFederation: &config.OIDCFederation{
					FederationRuleID: "rule-1",
					IdentityToken: &config.IdentityTokenConfig{
						Source: config.IdentityTokenSourceFile,
						Path:   tokenPath,
					},
				},
			},
		}),
		option.WithMaxRetries(0), // disable normal retries to isolate 401 retry
	)

	_, err := client.Messages.New(context.Background(), defaultParams)
	if err != nil {
		t.Fatal(err)
	}

	if n := apiCalls.Load(); n != 2 {
		t.Fatalf("expected 2 API calls (initial + retry), got %d", n)
	}
	if n := tokenCalls.Load(); n != 2 {
		t.Fatalf("expected 2 token calls (initial + after invalidation), got %d", n)
	}
}

func TestIntegration_ZeroConfigWorkloadIdentity(t *testing.T) {
	// Single server handles both OAuth token exchange and API requests.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/oauth/token" {
			json.NewEncoder(w).Encode(map[string]any{
				"access_token": "zero-config-tok",
			})
			return
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer zero-config-tok" {
			t.Errorf("got Authorization %q, want %q", authHeader, "Bearer zero-config-tok")
		}
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)
	unsetEnv(t, "ANTHROPIC_CREDENTIALS_FILE")
	unsetEnv(t, "ANTHROPIC_PROFILE")

	// Point the config dir at an empty temp dir so the default credential
	// chain falls through the profile step to the env-federation step
	// instead of finding whatever profile the test host already has.
	t.Setenv("ANTHROPIC_CONFIG_DIR", t.TempDir())

	// Set base URL via env so DefaultClientOptions picks it up for both
	// the API and the OAuth token endpoint.
	t.Setenv("ANTHROPIC_BASE_URL", server.URL)
	t.Setenv("ANTHROPIC_IDENTITY_TOKEN", "my-jwt")
	t.Setenv("ANTHROPIC_FEDERATION_RULE_ID", "rule-1")
	t.Setenv("ANTHROPIC_ORGANIZATION_ID", "org-1")

	client := anthropic.NewClient()

	_, err := client.Messages.New(context.Background(), defaultParams)
	if err != nil {
		t.Fatal(err)
	}
}

// TestIntegration_ZeroConfigProfile verifies that anthropic.NewClient() with
// no options picks up a profile from the config directory, per the spec's
// credential precedence chain.
func TestIntegration_ZeroConfigProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if got := r.Header.Get("Authorization"); got != "Bearer profile-tok" {
			t.Errorf("got Authorization %q, want %q", got, "Bearer profile-tok")
		}
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)
	unsetEnv(t, "ANTHROPIC_FEDERATION_RULE_ID")
	unsetEnv(t, "ANTHROPIC_ORGANIZATION_ID")
	unsetEnv(t, "ANTHROPIC_IDENTITY_TOKEN_FILE")
	unsetEnv(t, "ANTHROPIC_IDENTITY_TOKEN")
	unsetEnv(t, "ANTHROPIC_PROFILE")

	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "configs"), 0755)
	os.MkdirAll(filepath.Join(dir, "credentials"), 0700)
	os.WriteFile(
		filepath.Join(dir, "configs", "default.json"),
		[]byte(`{"authentication":{"type":"user_oauth"}}`),
		0644,
	)
	os.WriteFile(
		filepath.Join(dir, "credentials", "default.json"),
		[]byte(`{"type":"oauth_token","access_token":"profile-tok"}`),
		0600,
	)

	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_BASE_URL", server.URL)

	client := anthropic.NewClient()
	if _, err := client.Messages.New(context.Background(), defaultParams); err != nil {
		t.Fatalf("zero-config NewClient() should find the default profile: %v", err)
	}
}

func TestIntegration_WithProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if got := r.Header.Get("Authorization"); got != "Bearer staging-tok" {
			t.Errorf("got Authorization %q, want %q", got, "Bearer staging-tok")
		}
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "configs"), 0755)
	os.MkdirAll(filepath.Join(dir, "credentials"), 0700)
	os.WriteFile(
		filepath.Join(dir, "configs", "staging.json"),
		[]byte(`{"authentication":{"type":"user_oauth"}}`),
		0644,
	)
	os.WriteFile(
		filepath.Join(dir, "credentials", "staging.json"),
		[]byte(`{"type":"oauth_token","access_token":"staging-tok"}`),
		0600,
	)
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)

	client := anthropic.NewClient(
		option.WithoutEnvironmentDefaults(),
		option.WithBaseURL(server.URL),
		option.WithProfile("staging"),
	)
	if _, err := client.Messages.New(context.Background(), defaultParams); err != nil {
		t.Fatalf("WithProfile should authenticate from the named profile: %v", err)
	}
}

func TestIntegration_WithProfile_MissingProfileSurfacesError(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	client := anthropic.NewClient(
		option.WithoutEnvironmentDefaults(),
		option.WithBaseURL("http://unused.invalid"),
		option.WithMaxRetries(0),
		option.WithProfile("does-not-exist"),
	)
	_, err := client.Messages.New(context.Background(), defaultParams)
	if err == nil {
		t.Fatal("expected load error from WithProfile with missing profile")
	}
	if !strings.Contains(err.Error(), `WithProfile("does-not-exist")`) {
		t.Fatalf("error should name the profile and option, got: %v", err)
	}
}

func TestIntegration_WithProfile_RequestAPIKeyReplacesLoadError(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if got := r.Header.Get("Authorization"); got != "Bearer sk-test-key" {
			t.Errorf("got Authorization %q, want %q", got, "Bearer sk-test-key")
		}
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	client := anthropic.NewClient(
		option.WithoutEnvironmentDefaults(),
		option.WithBaseURL(server.URL),
		option.WithProfile("does-not-exist"),
	)
	if _, err := client.Messages.New(context.Background(), defaultParams, option.WithAPIKey("sk-test-key")); err != nil {
		t.Fatalf("WithAPIKey on the request should replace a broken WithProfile: %v", err)
	}
}

func TestIntegration_WithProfile_EmptyName(t *testing.T) {
	isolateAuthEnv(t)
	client := anthropic.NewClient(
		option.WithoutEnvironmentDefaults(),
		option.WithBaseURL("http://unused.invalid"),
		option.WithMaxRetries(0),
		option.WithProfile(""),
	)
	_, err := client.Messages.New(context.Background(), defaultParams)
	if err == nil || !strings.Contains(err.Error(), "profile name is empty") {
		t.Fatalf("expected empty-name error, got: %v", err)
	}
}

// TestIntegration_NoCredentialsAggregatedError verifies that when no source
// provides credentials, the first API request returns an aggregated error
// listing every source the SDK tried and what each failed with.
func TestIntegration_NoCredentialsAggregatedError(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)
	unsetEnv(t, "ANTHROPIC_FEDERATION_RULE_ID")
	unsetEnv(t, "ANTHROPIC_ORGANIZATION_ID")
	unsetEnv(t, "ANTHROPIC_IDENTITY_TOKEN_FILE")
	unsetEnv(t, "ANTHROPIC_IDENTITY_TOKEN")
	unsetEnv(t, "ANTHROPIC_PROFILE")

	dir := t.TempDir() // empty config dir - no profile file
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_BASE_URL", "http://127.0.0.1:1") // never reached

	client := anthropic.NewClient()
	_, err := client.Messages.New(context.Background(), defaultParams)
	if err == nil {
		t.Fatal("expected error when no credentials source is configured")
	}
	msg := err.Error()
	if !strings.Contains(msg, "no Anthropic credentials") {
		t.Errorf("expected aggregated error message, got: %v", err)
	}
	if !strings.Contains(msg, "ANTHROPIC_API_KEY") {
		t.Errorf("expected source list to mention ANTHROPIC_API_KEY, got: %v", err)
	}
	if !strings.Contains(msg, "profile") {
		t.Errorf("expected source list to mention profile, got: %v", err)
	}
	if !strings.Contains(msg, "federation") {
		t.Errorf("expected source list to mention federation, got: %v", err)
	}
	if !strings.Contains(msg, "ant auth login") {
		t.Errorf("expected remediation hint in message, got: %v", err)
	}
}

// TestIntegration_NoCredentialsPartialFederation verifies that when the user
// set some-but-not-all federation env vars, the aggregated error tells them
// which specific vars are missing instead of silently falling through to the
// "nothing is set" case.
func TestIntegration_NoCredentialsPartialFederation(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)
	unsetEnv(t, "ANTHROPIC_ORGANIZATION_ID")
	unsetEnv(t, "ANTHROPIC_IDENTITY_TOKEN_FILE")
	unsetEnv(t, "ANTHROPIC_IDENTITY_TOKEN")
	unsetEnv(t, "ANTHROPIC_PROFILE")

	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_BASE_URL", "http://127.0.0.1:1")
	t.Setenv("ANTHROPIC_FEDERATION_RULE_ID", "rule-1") // only one of three set

	client := anthropic.NewClient()
	_, err := client.Messages.New(context.Background(), defaultParams)
	if err == nil {
		t.Fatal("expected error when federation env vars are only partially set")
	}
	msg := err.Error()
	if !strings.Contains(msg, "partial configuration") {
		t.Errorf("expected aggregated error to mark the env federation source as partial, got: %v", err)
	}
	if !strings.Contains(msg, "missing: ANTHROPIC_ORGANIZATION_ID") {
		t.Errorf("expected partial-config detail to name missing ANTHROPIC_ORGANIZATION_ID first, got: %v", err)
	}
}

// TestIntegration_AggregatedErrorNotShownWhenKeySet verifies the aggregated
// error does NOT fire when a credential source actually works.
func TestIntegration_AggregatedErrorNotShownWhenKeySet(t *testing.T) {
	var sawHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sawHeader = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)
	unsetEnv(t, "ANTHROPIC_PROFILE")
	t.Setenv("ANTHROPIC_API_KEY", "sk-test-key")
	t.Setenv("ANTHROPIC_CONFIG_DIR", t.TempDir())
	t.Setenv("ANTHROPIC_BASE_URL", server.URL)

	client := anthropic.NewClient()
	if _, err := client.Messages.New(context.Background(), defaultParams); err != nil {
		t.Fatalf("expected success when API key is set: %v", err)
	}
	if sawHeader != "Bearer sk-test-key" {
		t.Fatalf("API key not forwarded: %q", sawHeader)
	}
}

func TestIntegration_APIKeyPrecedencePreserved(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	var receivedHeaders http.Header

	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successResponse))
	}))
	defer apiServer.Close()

	client := anthropic.NewClient(
		option.WithBaseURL(apiServer.URL),
		option.WithAPIKey("sk-test-key"),
	)

	_, err := client.Messages.New(context.Background(), defaultParams)
	if err != nil {
		t.Fatal(err)
	}

	if got := receivedHeaders.Get("Authorization"); got != "Bearer sk-test-key" {
		t.Fatalf("got Authorization %q, want %q", got, "Bearer sk-test-key")
	}
	if got := receivedHeaders.Get("X-Api-Key"); got != "" {
		t.Fatalf("expected no X-Api-Key header, got %q", got)
	}
}

// TestIntegration_SharedFederationOptionPerClientTransport verifies that a
// single WithFederationTokenProvider option value shared across two clients
// with different HTTP transports does not cause the second client's token
// exchange to run on the first client's transport.
func TestIntegration_SharedFederationOptionPerClientTransport(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	var serverA, serverB *httptest.Server
	serverA = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/oauth/token" {
			json.NewEncoder(w).Encode(map[string]any{"access_token": "tok-A"})
			return
		}
		if got := r.Header.Get("Authorization"); got != "Bearer tok-A" {
			t.Errorf("serverA got Authorization %q, want %q", got, "Bearer tok-A")
		}
		w.Write([]byte(successResponse))
	}))
	defer serverA.Close()
	serverB = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/oauth/token" {
			json.NewEncoder(w).Encode(map[string]any{"access_token": "tok-B"})
			return
		}
		if got := r.Header.Get("Authorization"); got != "Bearer tok-B" {
			t.Errorf("serverB got Authorization %q, want %q", got, "Bearer tok-B")
		}
		w.Write([]byte(successResponse))
	}))
	defer serverB.Close()

	// Count which underlying transport each client actually uses for its
	// token-exchange traffic. If the fix is wrong, client B's token exchange
	// runs on client A's transport and the counter for A sees both.
	var aTransportCalls, bTransportCalls atomic.Int32
	transportA := &countingTransport{n: &aTransportCalls, next: http.DefaultTransport}
	transportB := &countingTransport{n: &bTransportCalls, next: http.DefaultTransport}

	// Shared option value, deliberately reused across both clients.
	sharedAuth := option.WithFederationTokenProvider(
		func(_ context.Context) (string, error) { return "jwt", nil },
		option.FederationOptions{FederationRuleID: "rule-1", OrganizationID: "org-1"},
	)

	clientA := anthropic.NewClient(
		option.WithBaseURL(serverA.URL),
		option.WithHTTPClient(&http.Client{Transport: transportA}),
		sharedAuth,
	)
	clientB := anthropic.NewClient(
		option.WithBaseURL(serverB.URL),
		option.WithHTTPClient(&http.Client{Transport: transportB}),
		sharedAuth,
	)

	if _, err := clientA.Messages.New(context.Background(), defaultParams); err != nil {
		t.Fatal(err)
	}
	if _, err := clientB.Messages.New(context.Background(), defaultParams); err != nil {
		t.Fatal(err)
	}

	// Each client must have driven its own transport for both the token
	// exchange and the API request (2 calls each).
	if n := aTransportCalls.Load(); n != 2 {
		t.Fatalf("expected 2 calls on transport A, got %d", n)
	}
	if n := bTransportCalls.Load(); n != 2 {
		t.Fatalf("expected 2 calls on transport B, got %d", n)
	}
}

type countingTransport struct {
	n    *atomic.Int32
	next http.RoundTripper
}

func (t *countingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.n.Add(1)
	return t.next.RoundTrip(req)
}

// TestIntegration_WithConfigBaseURLOrderIndependent asserts that an explicit
// WithBaseURL wins over config.BaseURL regardless of the order in which
// WithBaseURL and WithConfig are passed.
func TestIntegration_WithConfigBaseURLOrderIndependent(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	var explicitHits atomic.Int32
	explicit := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		explicitHits.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successResponse))
	}))
	defer explicit.Close()

	// config.BaseURL points at an address that would fail if actually used.
	cfg := &config.Config{
		BaseURL: "http://127.0.0.1:1",
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeUserOAuth,
			CredentialsPath: writeStaticCred(t),
			UserOAuth:       &config.UserOAuth{},
		},
	}

	for _, name := range []string{"base_first", "config_first"} {
		t.Run(name, func(t *testing.T) {
			explicitHits.Store(0)
			var client anthropic.Client
			if name == "base_first" {
				client = anthropic.NewClient(
					option.WithBaseURL(explicit.URL),
					option.WithConfig(cfg),
				)
			} else {
				client = anthropic.NewClient(
					option.WithConfig(cfg),
					option.WithBaseURL(explicit.URL),
				)
			}
			if _, err := client.Messages.New(context.Background(), defaultParams); err != nil {
				t.Fatal(err)
			}
			if explicitHits.Load() == 0 {
				t.Fatal("explicit base URL did not receive the request")
			}
		})
	}
}

// TestIntegration_WithConfigMiddlewareCountStable makes several requests with
// the same WithConfig option value and checks that each request sees exactly
// one auth middleware invocation per request (not N on the Nth request).
func TestIntegration_WithConfigMiddlewareCountStable(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	var apiCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCalls.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	client := anthropic.NewClient(
		option.WithBaseURL(server.URL),
		option.WithConfig(&config.Config{
			AuthenticationInfo: &config.AuthenticationInfo{
				Type:            config.AuthenticationTypeUserOAuth,
				CredentialsPath: writeStaticCred(t),
				UserOAuth:       &config.UserOAuth{},
			},
		}),
	)

	for i := 0; i < 5; i++ {
		if _, err := client.Messages.New(context.Background(), defaultParams); err != nil {
			t.Fatalf("request %d: %v", i, err)
		}
	}
	if n := apiCalls.Load(); n != 5 {
		t.Fatalf("expected 5 API hits, got %d", n)
	}
}

func writeStaticCred(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	credPath := filepath.Join(dir, "creds.json")
	os.WriteFile(credPath, []byte(`{"type":"oauth_token","access_token":"tok"}`), 0600)
	return credPath
}

// TestDefaultClient_ExplicitProfileBeatsEnvFederation verifies the spec rule
// that ANTHROPIC_PROFILE takes precedence over the env-var federation path:
// when both are set, the profile's federation_rule_id wins on the wire.
func TestDefaultClient_ExplicitProfileBeatsEnvFederation(t *testing.T) {
	var sawBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/oauth/token" {
			_ = json.NewDecoder(r.Body).Decode(&sawBody)
			json.NewEncoder(w).Encode(map[string]any{"access_token": "tok"})
			return
		}
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	dir := t.TempDir()
	tokenPath := filepath.Join(dir, "token")
	os.WriteFile(tokenPath, []byte("profile-jwt"), 0600)
	os.MkdirAll(filepath.Join(dir, "configs"), 0755)
	os.WriteFile(filepath.Join(dir, "configs", "dev.json"), []byte(`{
  "authentication": {
    "type": "oidc_federation",
    "federation_rule_id": "fdrl_from_profile",
    "identity_token": {"source": "file", "path": "`+tokenPath+`"}
  },
  "organization_id": "org_from_profile"
}`), 0644)

	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "dev")
	t.Setenv("ANTHROPIC_BASE_URL", server.URL)
	t.Setenv("ANTHROPIC_FEDERATION_RULE_ID", "fdrl_from_env")
	t.Setenv("ANTHROPIC_ORGANIZATION_ID", "org_from_env")
	t.Setenv("ANTHROPIC_IDENTITY_TOKEN", "env-jwt")

	client := anthropic.NewClient()
	if _, err := client.Messages.New(context.Background(), defaultParams); err != nil {
		t.Fatal(err)
	}
	if got := sawBody["federation_rule_id"]; got != "fdrl_from_profile" {
		t.Errorf("exchange federation_rule_id: got %v, want fdrl_from_profile (profile should beat env)", got)
	}
	if got := sawBody["assertion"]; got != "profile-jwt" {
		t.Errorf("exchange assertion: got %v, want profile-jwt (profile's identity_token should win)", got)
	}
}

// TestDefaultClient_EnvFederationBeatsFallbackProfile verifies the spec rule
// that a fallback profile (active_config / "default" — no explicit
// ANTHROPIC_PROFILE) loses to the direct env-var federation path. A leftover
// default.json on a WIF-configured machine must not silently replace WIF.
func TestDefaultClient_EnvFederationBeatsFallbackProfile(t *testing.T) {
	var sawBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/oauth/token" {
			_ = json.NewDecoder(r.Body).Decode(&sawBody)
			json.NewEncoder(w).Encode(map[string]any{"access_token": "tok"})
			return
		}
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)
	unsetEnv(t, "ANTHROPIC_PROFILE")

	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "configs"), 0755)
	os.WriteFile(filepath.Join(dir, "configs", "default.json"), []byte(`{
  "authentication": {
    "type": "oidc_federation",
    "federation_rule_id": "fdrl_from_profile",
    "identity_token": {"source": "file", "path": "/tmp/unused"}
  },
  "organization_id": "org_from_profile"
}`), 0644)

	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_BASE_URL", server.URL)
	t.Setenv("ANTHROPIC_FEDERATION_RULE_ID", "fdrl_from_env")
	t.Setenv("ANTHROPIC_ORGANIZATION_ID", "org_from_env")
	t.Setenv("ANTHROPIC_IDENTITY_TOKEN", "env-jwt")

	client := anthropic.NewClient()
	if _, err := client.Messages.New(context.Background(), defaultParams); err != nil {
		t.Fatal(err)
	}
	if got := sawBody["federation_rule_id"]; got != "fdrl_from_env" {
		t.Errorf("exchange federation_rule_id: got %v, want fdrl_from_env (env should beat fallback profile)", got)
	}
	if got := sawBody["assertion"]; got != "env-jwt" {
		t.Errorf("exchange assertion: got %v, want env-jwt", got)
	}
}

// TestDefaultClient_ExplicitProfileMissingFails verifies that when
// ANTHROPIC_PROFILE names a profile that does not exist, the SDK surfaces
// the error instead of silently falling through to env federation or the
// "no credentials" aggregate. The user explicitly asked for a profile.
func TestDefaultClient_ExplicitProfileMissingFails(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "nonexistent")
	t.Setenv("ANTHROPIC_BASE_URL", "http://127.0.0.1:1")

	client := anthropic.NewClient()
	_, err := client.Messages.New(context.Background(), defaultParams)
	if err == nil {
		t.Fatal("expected an error when ANTHROPIC_PROFILE names a missing profile")
	}
	msg := err.Error()
	if !strings.Contains(msg, "nonexistent") {
		t.Errorf("expected error to name the missing profile, got: %v", err)
	}
}

// TestDefaultClient_ExplicitAPIKeyOverridesBrokenProfile verifies that
// ANTHROPIC_PROFILE is not read when the caller supplies option.WithAPIKey.
// The user is asking for an explicit credential; the shell's
// ANTHROPIC_PROFILE should not strand them.
func TestDefaultClient_ExplicitAPIKeyOverridesBrokenProfile(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	var receivedHeaders http.Header
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successResponse))
	}))
	defer apiServer.Close()

	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "nonexistent")

	client := anthropic.NewClient(
		option.WithBaseURL(apiServer.URL),
		option.WithAPIKey("sk-test-key"),
	)
	_, err := client.Messages.New(context.Background(), defaultParams)
	if err != nil {
		t.Fatalf("expected explicit API key to override broken profile, got: %v", err)
	}
	if got := receivedHeaders.Get("Authorization"); got != "Bearer sk-test-key" {
		t.Errorf("Authorization: got %q, want %q", got, "Bearer sk-test-key")
	}
}

// TestDefaultClient_ExplicitAuthTokenOverridesBrokenProfile is the
// companion for option.WithAuthToken.
func TestDefaultClient_ExplicitAuthTokenOverridesBrokenProfile(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	var receivedHeaders http.Header
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successResponse))
	}))
	defer apiServer.Close()

	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "nonexistent")

	client := anthropic.NewClient(
		option.WithBaseURL(apiServer.URL),
		option.WithAuthToken("sk-test-token"),
	)
	_, err := client.Messages.New(context.Background(), defaultParams)
	if err != nil {
		t.Fatalf("expected explicit auth token to override broken profile, got: %v", err)
	}
	if got := receivedHeaders.Get("Authorization"); got != "Bearer sk-test-token" {
		t.Errorf("Authorization: got %q, want %q", got, "Bearer sk-test-token")
	}
}

// brokenUserOAuthConfig returns a user_oauth Config whose CredentialsPath
// points at a missing file under t.TempDir, so ResolveCredentials errors.
// BaseURL is set so callers can assert it applies independently.
func brokenUserOAuthConfig(t *testing.T, baseURL string) *config.Config {
	t.Helper()
	return &config.Config{
		BaseURL: baseURL,
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeUserOAuth,
			CredentialsPath: filepath.Join(t.TempDir(), "missing.json"),
			UserOAuth:       &config.UserOAuth{ClientID: "x"},
		},
	}
}

// TestWithConfig_RequestAPIKeyReplacesBrokenProfile verifies that a
// WithAPIKey passed on a request authenticates it when the client's
// WithConfig has credentials that cannot be resolved — i.e. the resolution
// error is deferred to a request-time middleware, not returned at
// option-Apply time where it would short-circuit later options.
func TestWithConfig_RequestAPIKeyReplacesBrokenProfile(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	var receivedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	client := anthropic.NewClient(
		option.WithConfig(brokenUserOAuthConfig(t, server.URL)),
	)
	if _, err := client.Messages.New(context.Background(), defaultParams, option.WithAPIKey("sk-x")); err != nil {
		t.Fatalf("expected request to succeed via WithAPIKey despite broken profile, got: %v", err)
	}
	if got := receivedHeaders.Get("Authorization"); got != "Bearer sk-x" {
		t.Fatalf("got Authorization %q, want Bearer sk-x", got)
	}
}

// TestWithConfig_RequestAuthTokenReplacesBrokenProfile is the WithAuthToken
// variant of the above.
func TestWithConfig_RequestAuthTokenReplacesBrokenProfile(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	var receivedHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	client := anthropic.NewClient(
		option.WithConfig(brokenUserOAuthConfig(t, server.URL)),
	)
	if _, err := client.Messages.New(context.Background(), defaultParams, option.WithAuthToken("tok-x")); err != nil {
		t.Fatalf("expected request to succeed via WithAuthToken despite broken profile, got: %v", err)
	}
	if got := receivedHeaders.Get("Authorization"); got != "Bearer tok-x" {
		t.Fatalf("got Authorization %q, want %q", got, "Bearer tok-x")
	}
}

// TestWithConfig_BaseURLAppliesEvenWithBrokenCreds verifies that the
// profile's BaseURL is honored regardless of credential-resolution outcome —
// it's applied at option-Apply time, not gated behind ResolveCredentials.
func TestWithConfig_BaseURLAppliesEvenWithBrokenCreds(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	var hit atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hit.Store(true)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	// No WithBaseURL — only the profile's BaseURL.
	client := anthropic.NewClient(
		option.WithConfig(brokenUserOAuthConfig(t, server.URL)),
	)
	if _, err := client.Messages.New(context.Background(), defaultParams, option.WithAPIKey("sk-x")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hit.Load() {
		t.Fatal("request did not reach the profile's BaseURL")
	}
}

// TestWithConfig_BrokenProfileNoOverrideStillFails verifies the negative:
// when no credential option on the request replaces the profile, the
// deferred resolution error surfaces on the first request.
func TestWithConfig_BrokenProfileNoOverrideStillFails(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	client := anthropic.NewClient(
		option.WithConfig(brokenUserOAuthConfig(t, "http://127.0.0.1:1")),
		option.WithMaxRetries(0),
	)
	_, err := client.Messages.New(context.Background(), defaultParams)
	if err == nil {
		t.Fatal("expected credential-resolution error to surface, got nil")
	}
	if !strings.Contains(err.Error(), "missing.json") && !strings.Contains(err.Error(), "credential") {
		t.Fatalf("expected the deferred resolve error, got: %v", err)
	}
}

// TestDefaultClient_NoExplicitProfileFallsThrough verifies the full
// no-credentials chain: no env API key, no ANTHROPIC_PROFILE, no profile
// file on disk, no env federation — the SDK returns the aggregated
// NoCredentialsError.
func TestDefaultClient_NoExplicitProfileFallsThrough(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)
	unsetEnv(t, "ANTHROPIC_PROFILE")

	dir := t.TempDir() // no configs/ subdir
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_BASE_URL", "http://127.0.0.1:1")

	client := anthropic.NewClient()
	_, err := client.Messages.New(context.Background(), defaultParams)
	if err == nil {
		t.Fatal("expected NoCredentialsError when nothing is set")
	}
	if !strings.Contains(err.Error(), "no Anthropic credentials") {
		t.Errorf("expected aggregated error, got: %v", err)
	}
}

type recordedRequest struct {
	method string
	url    string
	auth   string
}

type recordingHandler struct {
	mu   sync.Mutex
	reqs []recordedRequest
}

func (h *recordingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.reqs = append(h.reqs, recordedRequest{method: r.Method, url: r.URL.String(), auth: r.Header.Get("Authorization")})
	h.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{}`))
}

func (h *recordingHandler) requests() []recordedRequest {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]recordedRequest(nil), h.reqs...)
}

func assertNoCredentials(t *testing.T, reqs []recordedRequest) {
	t.Helper()
	for _, r := range reqs {
		if strings.Contains(r.url, "/v1/oauth/token") {
			t.Errorf("token exchange was sent to a foreign origin: %s %s", r.method, r.url)
		}
		if r.auth != "" {
			t.Errorf("credential was sent to a foreign origin: %s %s Authorization=%q", r.method, r.url, r.auth)
		}
	}
}

func federationServer(t *testing.T) (*httptest.Server, *recordingHandler, *atomic.Int32) {
	t.Helper()
	var exchanges atomic.Int32
	api := &recordingHandler{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/oauth/token" {
			exchanges.Add(1)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"access_token": "exchanged-tok"})
			return
		}
		api.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)
	return server, api, &exchanges
}

func TestIntegration_AbsoluteURLForeignHost_FederationProvider(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	base, baseAPI, exchanges := federationServer(t)
	attacker := &recordingHandler{}
	attackerServer := httptest.NewServer(attacker)
	defer attackerServer.Close()

	var identityCalls atomic.Int32
	client := anthropic.NewClient(
		option.WithMaxRetries(0),
		option.WithBaseURL(base.URL),
		option.WithFederationTokenProvider(
			func(context.Context) (string, error) {
				identityCalls.Add(1)
				return "secret-jwt", nil
			},
			option.FederationOptions{FederationRuleID: "rule-1", OrganizationID: "org-1"},
		),
	)

	err := client.Get(context.Background(), attackerServer.URL+"/v1/collect", nil, nil)
	assertNoCredentials(t, attacker.requests())
	if err != nil {
		t.Fatalf("foreign-host request: %v", err)
	}
	if got := identityCalls.Load(); got != 0 {
		t.Errorf("identity token was fetched %d times for a foreign-host request, want 0", got)
	}
	if got := exchanges.Load(); got != 0 {
		t.Errorf("token exchange ran %d times for a foreign-host request, want 0", got)
	}

	if err := client.Get(context.Background(), "/v1/models", nil, nil); err != nil {
		t.Fatalf("same-origin request: %v", err)
	}
	if got := exchanges.Load(); got != 1 {
		t.Fatalf("token exchange ran %d times at the configured base URL, want 1", got)
	}
	reqs := baseAPI.requests()
	if len(reqs) != 1 || reqs[0].auth != "Bearer exchanged-tok" {
		t.Fatalf("configured base URL got %+v, want one request with the exchanged bearer", reqs)
	}
	assertNoCredentials(t, attacker.requests())
}

// Same as above for a config-file profile with no base_url of its own.
func TestIntegration_AbsoluteURLForeignHost_ProfileFederation(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	base, baseAPI, exchanges := federationServer(t)
	attacker := &recordingHandler{}
	attackerServer := httptest.NewServer(attacker)
	defer attackerServer.Close()

	dir := t.TempDir()
	tokenPath := filepath.Join(dir, "identity-token")
	os.WriteFile(tokenPath, []byte("secret-jwt"), 0600)
	client := anthropic.NewClient(
		option.WithMaxRetries(0),
		option.WithBaseURL(base.URL),
		option.WithConfig(&config.Config{
			OrganizationID: "org-1",
			AuthenticationInfo: &config.AuthenticationInfo{
				Type:            config.AuthenticationTypeOIDCFederation,
				CredentialsPath: filepath.Join(dir, "credentials", "creds.json"),
				OIDCFederation: &config.OIDCFederation{
					FederationRuleID: "rule-1",
					IdentityToken:    &config.IdentityTokenConfig{Source: config.IdentityTokenSourceFile, Path: tokenPath},
				},
			},
		}),
	)

	err := client.Post(context.Background(), attackerServer.URL+"/v1/collect", nil, nil)
	assertNoCredentials(t, attacker.requests())
	if err != nil {
		t.Fatalf("foreign-host request: %v", err)
	}
	if got := exchanges.Load(); got != 0 {
		t.Errorf("token exchange ran %d times for a foreign-host request, want 0", got)
	}

	if err := client.Get(context.Background(), "/v1/models", nil, nil); err != nil {
		t.Fatalf("same-origin request: %v", err)
	}
	if got := exchanges.Load(); got != 1 {
		t.Fatalf("token exchange ran %d times at the configured base URL, want 1", got)
	}
	if reqs := baseAPI.requests(); len(reqs) != 1 || reqs[0].auth != "Bearer exchanged-tok" {
		t.Fatalf("configured base URL got %+v, want one request with the exchanged bearer", reqs)
	}
	assertNoCredentials(t, attacker.requests())
}

// recordingTransport records requests for origins no test server listens on.
type recordingTransport struct {
	recordingHandler
}

func (rt *recordingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, r)
	resp := rec.Result()
	resp.Request = r
	return resp, nil
}

func TestIntegration_AbsoluteURLSchemeDowngrade_NoCredentials(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	credPath := filepath.Join(t.TempDir(), "credentials.json")
	os.WriteFile(credPath, []byte(`{"type":"oauth_token","access_token":"static-tok"}`), 0600)

	transport := &recordingTransport{}
	client := anthropic.NewClient(
		option.WithHTTPClient(&http.Client{Transport: transport}),
		option.WithBaseURL("https://api.example.test"),
		option.WithConfig(&config.Config{
			AuthenticationInfo: &config.AuthenticationInfo{
				Type:            config.AuthenticationTypeUserOAuth,
				CredentialsPath: credPath,
				UserOAuth:       &config.UserOAuth{},
			},
		}),
	)

	if err := client.Get(context.Background(), "https://api.example.test/v1/models", nil, nil); err != nil {
		t.Fatalf("same-origin absolute URL: %v", err)
	}
	if reqs := transport.requests(); len(reqs) != 1 || reqs[0].auth != "Bearer static-tok" {
		t.Fatalf("same-origin absolute URL got %+v, want the bearer attached", reqs)
	}

	if err := client.Get(context.Background(), "http://api.example.test/v1/models", nil, nil); err != nil {
		t.Fatalf("downgraded absolute URL: %v", err)
	}
	reqs := transport.requests()
	if len(reqs) != 2 {
		t.Fatalf("got %d requests, want 2", len(reqs))
	}
	if reqs[1].auth != "" {
		t.Fatalf("bearer was attached to cleartext request %s: Authorization=%q", reqs[1].url, reqs[1].auth)
	}
}

// exchangeServer answers token exchanges with token and records every other
// request. If hold is non-nil, every exchange waits for hold to be closed
// before answering, and the first one closes held when it arrives.
func exchangeServer(t *testing.T, token string, held chan<- struct{}, hold <-chan struct{}) (*httptest.Server, *recordingHandler) {
	t.Helper()
	api := &recordingHandler{}
	var heldOnce sync.Once
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/oauth/token" {
			if hold != nil {
				heldOnce.Do(func() { close(held) })
				<-hold
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"access_token": token})
			return
		}
		api.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)
	return server, api
}

func assertBearer(t *testing.T, server string, reqs []recordedRequest, want string) {
	t.Helper()
	if len(reqs) != 1 {
		t.Errorf("%s got %d API requests, want 1", server, len(reqs))
	}
	for _, r := range reqs {
		if r.auth != want {
			t.Errorf("%s got Authorization %q, want %q", server, r.auth, want)
		}
	}
}

// A request that overrides the base URL with option.WithBaseURL, sent while
// the client's own token exchange is in flight, must get a token for its own
// base URL.
func TestIntegration_PerRequestBaseURLDuringTokenExchange(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	held, hold := make(chan struct{}), make(chan struct{})
	serverA, apiA := exchangeServer(t, "tok-A", held, hold)
	serverB, apiB := exchangeServer(t, "tok-B", nil, nil)
	release := sync.OnceFunc(func() { close(hold) })
	t.Cleanup(release)

	client := anthropic.NewClient(
		option.WithMaxRetries(0),
		option.WithBaseURL(serverA.URL),
		option.WithFederationTokenProvider(
			func(context.Context) (string, error) { return "jwt", nil },
			option.FederationOptions{FederationRuleID: "rule-1", OrganizationID: "org-1"},
		),
	)

	errA := make(chan error, 1)
	go func() { errA <- client.Get(context.Background(), "/v1/models", nil, nil) }()
	<-held
	errB := make(chan error, 1)
	go func() {
		errB <- client.Get(context.Background(), "/v1/models", nil, nil, option.WithBaseURL(serverB.URL))
	}()
	// Let the second request reach the token cache before the first exchange finishes.
	time.Sleep(100 * time.Millisecond)
	release()

	if err := <-errA; err != nil {
		t.Fatal(err)
	}
	if err := <-errB; err != nil {
		t.Fatal(err)
	}
	assertBearer(t, "serverA", apiA.requests(), "Bearer tok-A")
	assertBearer(t, "serverB", apiB.requests(), "Bearer tok-B")
}

// tokenServer answers token requests with an expiring token, which the SDK saves
// to the credentials file, and passes every other request to api.
func tokenServer(t *testing.T, token string, api http.Handler) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var tokenRequests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/oauth/token" {
			api.ServeHTTP(w, r)
			return
		}
		tokenRequests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"access_token": token, "refresh_token": "rotated-refresh", "expires_in": 3600})
	}))
	t.Cleanup(server.Close)
	return server, &tokenRequests
}

func isolateProfileEnv(t *testing.T) {
	t.Helper()
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	unsetEnv(t, "ANTHROPIC_BASE_URL")
	isolateAuthEnv(t)
}

func federationConfig(t *testing.T, credPath string) *config.Config {
	t.Helper()
	identityTokenPath := filepath.Join(t.TempDir(), "identity-token")
	if err := os.WriteFile(identityTokenPath, []byte("jwt"), 0600); err != nil {
		t.Fatal(err)
	}
	return &config.Config{
		OrganizationID: "org-1",
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeOIDCFederation,
			CredentialsPath: credPath,
			OIDCFederation: &config.OIDCFederation{
				FederationRuleID: "rule-1",
				IdentityToken:    &config.IdentityTokenConfig{Source: config.IdentityTokenSourceFile, Path: identityTokenPath},
			},
		},
	}
}

func userOAuthConfig(credPath string) *config.Config {
	return &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeUserOAuth,
			CredentialsPath: credPath,
			UserOAuth:       &config.UserOAuth{ClientID: "client-1"},
		},
	}
}

func saveCredentials(t *testing.T, accessToken string, expiresIn time.Duration) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "creds.json")
	expiresAt := time.Now().Add(expiresIn)
	credentials := config.Credentials{AccessToken: accessToken, RefreshToken: "saved-refresh", ExpiresAt: &expiresAt}
	if err := config.WriteCredentials(path, credentials); err != nil {
		t.Fatal(err)
	}
	return path
}

func fileContents(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// clientBaseURLSources are the ways a client gets the base URL that a method's
// WithBaseURL overrides.
var clientBaseURLSources = []struct {
	name    string
	options func(cfg *config.Config, baseURL string) []option.RequestOption
}{
	{
		name: "WithBaseURL",
		options: func(cfg *config.Config, baseURL string) []option.RequestOption {
			return []option.RequestOption{option.WithMaxRetries(0), option.WithConfig(cfg), option.WithBaseURL(baseURL)}
		},
	},
	{
		name: "default base URL",
		options: func(cfg *config.Config, baseURL string) []option.RequestOption {
			return []option.RequestOption{option.WithMaxRetries(0), option.WithConfig(cfg), requestconfig.WithDefaultBaseURL(baseURL)}
		},
	},
	{
		name: "profile base_url",
		options: func(cfg *config.Config, baseURL string) []option.RequestOption {
			cfg.BaseURL = baseURL
			return []option.RequestOption{option.WithMaxRetries(0), option.WithConfig(cfg)}
		},
	},
}

// Server B stands for a proxy that forwards API paths and has no token endpoint,
// so its one request has to carry the file's token.
func TestIntegration_PerRequestBaseURL_FederationUsesValidFileToken(t *testing.T) {
	for _, source := range clientBaseURLSources {
		t.Run(source.name, func(t *testing.T) {
			isolateProfileEnv(t)
			apiB := &recordingHandler{}
			serverA, exchangesA := tokenServer(t, "tok-A", &recordingHandler{})
			serverB := httptest.NewServer(apiB)
			t.Cleanup(serverB.Close)
			credPath := saveCredentials(t, "saved-tok", time.Hour)
			before := fileContents(t, credPath)
			client := anthropic.NewClient(source.options(federationConfig(t, credPath), serverA.URL)...)

			err := client.Get(context.Background(), "/v1/models", nil, nil, option.WithBaseURL(serverB.URL))

			if err != nil {
				t.Fatal(err)
			}
			assertBearer(t, "serverB", apiB.requests(), "Bearer saved-tok")
			if got := exchangesA.Load(); got != 0 {
				t.Errorf("serverA got %d exchanges, want 0", got)
			}
			if got := fileContents(t, credPath); got != before {
				t.Errorf("credentials file changed to %s, want %s", got, before)
			}
		})
	}
}

// Server B rejects its first token, so the forced exchange of the 401 retry has
// to leave the file alone too. A later request to server A must not get the token
// obtained for server B.
func TestIntegration_PerRequestBaseURL_FederationDoesNotSaveOverrideToken(t *testing.T) {
	saved := []struct {
		name          string
		expiresIn     time.Duration
		wantExchanges int32
		wantTokenA    string
	}{
		{name: "fresh token", expiresIn: time.Hour, wantExchanges: 1, wantTokenA: "saved-tok"},
		{name: "expired token", expiresIn: -time.Hour, wantExchanges: 2, wantTokenA: "tok-A"},
	}
	for _, source := range clientBaseURLSources {
		for _, file := range saved {
			t.Run(source.name+"/"+file.name, func(t *testing.T) {
				isolateProfileEnv(t)
				var apiRequests atomic.Int32
				var lastBearerB atomic.Value
				apiA := &recordingHandler{}
				apiB := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					lastBearerB.Store(r.Header.Get("Authorization"))
					if apiRequests.Add(1) == 1 {
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
					w.Write([]byte(`{}`))
				})
				serverA, exchangesA := tokenServer(t, "tok-A", apiA)
				serverB, exchangesB := tokenServer(t, "tok-B", apiB)
				credPath := saveCredentials(t, "saved-tok", file.expiresIn)
				before := fileContents(t, credPath)
				cfg := federationConfig(t, credPath)
				client := anthropic.NewClient(source.options(cfg, serverA.URL)...)

				err := client.Get(context.Background(), "/v1/models", nil, nil, option.WithBaseURL(serverB.URL))

				if err != nil {
					t.Fatal(err)
				}
				wantBearer, wantExchangesA, wantExchangesB := "Bearer tok-B", int32(0), file.wantExchanges
				if cfg.BaseURL != "" {
					wantBearer, wantExchangesA, wantExchangesB = "Bearer tok-A", file.wantExchanges, 0
				}
				if got := lastBearerB.Load(); got != wantBearer {
					t.Errorf("serverB got Authorization %q, want %q", got, wantBearer)
				}
				if got := exchangesA.Load(); got != wantExchangesA {
					t.Errorf("serverA got %d exchanges, want %d", got, wantExchangesA)
				}
				if got := exchangesB.Load(); got != wantExchangesB {
					t.Errorf("serverB got %d exchanges, want %d", got, wantExchangesB)
				}
				if got := fileContents(t, credPath); got != before {
					t.Errorf("credentials file changed to %s, want %s", got, before)
				}

				if err := client.Get(context.Background(), "/v1/models", nil, nil); err != nil {
					t.Fatal(err)
				}
				assertBearer(t, "serverA", apiA.requests(), "Bearer "+file.wantTokenA)
			})
		}
	}
}

func TestIntegration_PerRequestBaseURL_ClientBaseURLUsesCredentialsFile(t *testing.T) {
	saved := []struct {
		name          string
		expiresIn     time.Duration
		wantToken     string
		wantExchanges int32
	}{
		{name: "fresh token is read", expiresIn: time.Hour, wantToken: "saved-tok", wantExchanges: 0},
		{name: "expired token is replaced", expiresIn: -time.Hour, wantToken: "tok-A", wantExchanges: 1},
	}
	for _, source := range clientBaseURLSources {
		for _, file := range saved {
			t.Run(source.name+"/"+file.name, func(t *testing.T) {
				isolateProfileEnv(t)
				api := &recordingHandler{}
				server, exchanges := tokenServer(t, "tok-A", api)
				credPath := saveCredentials(t, "saved-tok", file.expiresIn)
				client := anthropic.NewClient(source.options(federationConfig(t, credPath), server.URL)...)

				err := client.Get(context.Background(), "/v1/models", nil, nil)

				if err != nil {
					t.Fatal(err)
				}
				assertBearer(t, "server", api.requests(), "Bearer "+file.wantToken)
				if got := exchanges.Load(); got != file.wantExchanges {
					t.Errorf("server got %d exchanges, want %d", got, file.wantExchanges)
				}
				if got := fileContents(t, credPath); !strings.Contains(got, file.wantToken) {
					t.Errorf("credentials file is %s, want it to hold %s", got, file.wantToken)
				}
			})
		}
	}
}

func TestIntegration_PerRequestBaseURL_UserOAuthRefreshesWhereClientDoes(t *testing.T) {
	for _, source := range clientBaseURLSources {
		t.Run(source.name, func(t *testing.T) {
			isolateProfileEnv(t)
			apiB := &recordingHandler{}
			serverA, refreshesA := tokenServer(t, "tok-A", &recordingHandler{})
			serverB, refreshesB := tokenServer(t, "tok-B", apiB)
			credPath := saveCredentials(t, "saved-tok", -time.Hour)
			client := anthropic.NewClient(source.options(userOAuthConfig(credPath), serverA.URL)...)

			err := client.Get(context.Background(), "/v1/models", nil, nil, option.WithBaseURL(serverB.URL))

			if err != nil {
				t.Fatal(err)
			}
			assertBearer(t, "serverB", apiB.requests(), "Bearer tok-A")
			if got := refreshesA.Load(); got != 1 {
				t.Errorf("serverA got %d refreshes, want 1", got)
			}
			if got := refreshesB.Load(); got != 0 {
				t.Errorf("serverB got %d refreshes, want 0", got)
			}
			if got := fileContents(t, credPath); !strings.Contains(got, "rotated-refresh") {
				t.Errorf("credentials file is %s, want the rotated refresh token", got)
			}
		})
	}

	t.Run("profile base_url differs from client base URL", func(t *testing.T) {
		isolateProfileEnv(t)
		apiB := &recordingHandler{}
		serverClient, refreshesClient := tokenServer(t, "tok-client", &recordingHandler{})
		serverPin, refreshesPin := tokenServer(t, "tok-pin", &recordingHandler{})
		serverB, refreshesB := tokenServer(t, "tok-B", apiB)
		cfg := userOAuthConfig(saveCredentials(t, "saved-tok", -time.Hour))
		cfg.BaseURL = serverPin.URL
		client := anthropic.NewClient(option.WithMaxRetries(0), option.WithConfig(cfg), option.WithBaseURL(serverClient.URL))

		err := client.Get(context.Background(), "/v1/models", nil, nil, option.WithBaseURL(serverB.URL))

		if err != nil {
			t.Fatal(err)
		}
		assertBearer(t, "serverB", apiB.requests(), "Bearer tok-pin")
		if got := refreshesPin.Load(); got != 1 {
			t.Errorf("serverPin got %d refreshes, want 1", got)
		}
		if got := refreshesClient.Load() + refreshesB.Load(); got != 0 {
			t.Errorf("serverClient and serverB got %d refreshes, want 0", got)
		}
	})
}

// The token goes to the request's base URL, so a plain-HTTP, non-loopback override
// fails before anything is sent. A refresh goes to the client's base URL, so such a
// client fails only when the request needs a refresh.
func TestIntegration_PerRequestBaseURL_UserOAuthPlainHTTP(t *testing.T) {
	const secure, plain = "https://api.example.test", "http://plain.example.test"
	tests := []struct {
		name      string
		client    string
		override  string
		expiresIn time.Duration
		wantSent  []recordedRequest
	}{
		{name: "override/fresh token", client: secure, override: plain, expiresIn: time.Hour},
		{name: "override/expired token", client: secure, override: plain, expiresIn: -time.Hour},
		{name: "client/expired token", client: plain, override: secure, expiresIn: -time.Hour},
		{
			name: "client/fresh token", client: plain, override: secure, expiresIn: time.Hour,
			wantSent: []recordedRequest{{method: http.MethodGet, url: secure + "/v1/models", auth: "Bearer saved-tok"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolateProfileEnv(t)
			transport := &recordingTransport{}
			credPath := saveCredentials(t, "saved-tok", tt.expiresIn)
			before := fileContents(t, credPath)
			client := anthropic.NewClient(
				option.WithMaxRetries(0),
				option.WithHTTPClient(&http.Client{Transport: transport}),
				option.WithConfig(userOAuthConfig(credPath)),
				option.WithBaseURL(tt.client),
			)

			err := client.Get(context.Background(), "/v1/models", nil, nil, option.WithBaseURL(tt.override))

			if tt.wantSent != nil {
				if err != nil {
					t.Errorf("got error %v, want none", err)
				}
			} else if err == nil || !strings.Contains(err.Error(), "refusing to send credential over non-https token endpoint") {
				t.Errorf("got error %v, want the non-https refusal", err)
			}
			if got := transport.requests(); !slices.Equal(got, tt.wantSent) {
				t.Errorf("sent %+v, want %+v", got, tt.wantSent)
			}
			if got := fileContents(t, credPath); got != before {
				t.Errorf("credentials file changed to %s, want %s", got, before)
			}
		})
	}
}

// A profile passed to the method is the method's own, even when the client has the
// same one, so it refreshes where it would without a client.
func TestIntegration_PerRequestBaseURL_MethodProfileIsNotAnOverride(t *testing.T) {
	tests := []struct {
		name          string
		pinnedToOther bool
	}{
		{name: "no base_url"},
		{name: "base_url", pinnedToOther: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isolateProfileEnv(t)
			serverA, refreshesA := tokenServer(t, "tok-A", &recordingHandler{})
			serverB, refreshesB := tokenServer(t, "tok-B", &recordingHandler{})
			serverPin, refreshesPin := tokenServer(t, "tok-pin", &recordingHandler{})
			cfg := userOAuthConfig(saveCredentials(t, "saved-tok", -time.Hour))
			wantRefreshesB, wantRefreshesPin := int32(1), int32(0)
			if tt.pinnedToOther {
				cfg.BaseURL = serverPin.URL
				wantRefreshesB, wantRefreshesPin = 0, 1
			}
			client := anthropic.NewClient(option.WithMaxRetries(0), option.WithBaseURL(serverA.URL), option.WithConfig(cfg))

			err := client.Get(context.Background(), "/v1/models", nil, nil, option.WithConfig(cfg), option.WithBaseURL(serverB.URL))

			if err != nil {
				t.Fatal(err)
			}
			if got := refreshesA.Load(); got != 0 {
				t.Errorf("serverA got %d refreshes, want 0", got)
			}
			if got := refreshesB.Load(); got != wantRefreshesB {
				t.Errorf("serverB got %d refreshes, want %d", got, wantRefreshesB)
			}
			if got := refreshesPin.Load(); got != wantRefreshesPin {
				t.Errorf("serverPin got %d refreshes, want %d", got, wantRefreshesPin)
			}
		})
	}
}

// A service built without a client has no boundary to record its base URL at, so
// nothing counts as an override and the token from server B is saved.
func TestIntegration_PerRequestBaseURL_ServiceWithoutClientSavesToFile(t *testing.T) {
	isolateProfileEnv(t)
	apiB := &recordingHandler{}
	serverA, exchangesA := tokenServer(t, "tok-A", &recordingHandler{})
	serverB, exchangesB := tokenServer(t, "tok-B", apiB)
	credPath := saveCredentials(t, "saved-tok", -time.Hour)
	service := anthropic.NewModelService(option.WithMaxRetries(0), option.WithBaseURL(serverA.URL), option.WithConfig(federationConfig(t, credPath)))

	_, err := service.List(context.Background(), anthropic.ModelListParams{}, option.WithBaseURL(serverB.URL))

	if err != nil {
		t.Fatal(err)
	}
	assertBearer(t, "serverB", apiB.requests(), "Bearer tok-B")
	if got := exchangesA.Load(); got != 0 {
		t.Errorf("serverA got %d exchanges, want 0", got)
	}
	if got := exchangesB.Load(); got != 1 {
		t.Errorf("serverB got %d exchanges, want 1", got)
	}
	if got := fileContents(t, credPath); !strings.Contains(got, "tok-B") {
		t.Errorf("credentials file is %s, want it to hold tok-B", got)
	}
}
