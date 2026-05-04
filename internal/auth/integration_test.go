package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/config"
	"github.com/anthropics/anthropic-sdk-go/internal/auth"
	"github.com/anthropics/anthropic-sdk-go/option"
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

func TestIntegration_WithProfile_APIKeyPreemptsLoadError(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if got := r.Header.Get("X-Api-Key"); got != "sk-test-key" {
			t.Errorf("got X-Api-Key %q, want %q", got, "sk-test-key")
		}
		w.Write([]byte(successResponse))
	}))
	defer server.Close()

	client := anthropic.NewClient(
		option.WithoutEnvironmentDefaults(),
		option.WithBaseURL(server.URL),
		option.WithProfile("does-not-exist"),
		option.WithAPIKey("sk-test-key"),
	)
	if _, err := client.Messages.New(context.Background(), defaultParams); err != nil {
		t.Fatalf("WithAPIKey should preempt a broken WithProfile: %v", err)
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
	if !strings.Contains(msg, "anthropic auth login") {
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
		sawHeader = r.Header.Get("X-Api-Key")
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
	if sawHeader != "sk-test-key" {
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

	if receivedHeaders.Get("Authorization") != "" {
		t.Fatal("expected no Authorization header when using API key")
	}
	if receivedHeaders.Get("X-Api-Key") != "sk-test-key" {
		t.Fatalf("got X-Api-Key %q, want %q", receivedHeaders.Get("X-Api-Key"), "sk-test-key")
	}
}

// TestIntegration_APIKeyWithProfileLogsWarning verifies that when a static
// ANTHROPIC_API_KEY is passed alongside a profile config, the API key still
// wins (per the spec's credential precedence) but the user gets a one-shot
// warning so they can figure out why their profile's credentials are being
// ignored.
func TestIntegration_APIKeyWithProfileLogsWarning(t *testing.T) {
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
	credPath := filepath.Join(dir, "credentials.json")
	os.WriteFile(credPath, []byte(`{"type":"oauth_token","access_token":"profile-tok"}`), 0600)

	var logMu sync.Mutex
	var logBuf bytes.Buffer
	log.SetOutput(writerFunc(func(p []byte) (int, error) {
		logMu.Lock()
		defer logMu.Unlock()
		return logBuf.Write(p)
	}))
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	readLog := func() string {
		logMu.Lock()
		defer logMu.Unlock()
		return logBuf.String()
	}

	client := anthropic.NewClient(
		option.WithBaseURL(apiServer.URL),
		option.WithAPIKey("sk-test-key"),
		option.WithConfig(&config.Config{
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

	if receivedHeaders.Get("X-Api-Key") != "sk-test-key" {
		t.Fatalf("got X-Api-Key %q, want %q", receivedHeaders.Get("X-Api-Key"), "sk-test-key")
	}
	if receivedHeaders.Get("Authorization") == "Bearer profile-tok" {
		t.Fatal("profile credential should not have won over static API key")
	}

	out := readLog()
	if !strings.Contains(out, "ANTHROPIC_API_KEY") {
		t.Fatalf("expected warning mentioning ANTHROPIC_API_KEY, got: %q", out)
	}
	if !strings.Contains(out, "profile") && !strings.Contains(out, "config") {
		t.Fatalf("expected warning to mention the profile/config being shadowed, got: %q", out)
	}
}

// TestIntegration_WithConfigQuiet_SuppressesShadowWarning verifies that
// option.WithConfigQuiet behaves identically to WithConfig for credential
// precedence (API key still wins) but does NOT emit the config-shadowed
// warning, so callers with their own diagnostic don't double-warn.
func TestIntegration_WithConfigQuiet_SuppressesShadowWarning(t *testing.T) {
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
	credPath := filepath.Join(dir, "credentials.json")
	os.WriteFile(credPath, []byte(`{"type":"oauth_token","access_token":"profile-tok"}`), 0600)

	var logMu sync.Mutex
	var logBuf bytes.Buffer
	log.SetOutput(writerFunc(func(p []byte) (int, error) {
		logMu.Lock()
		defer logMu.Unlock()
		return logBuf.Write(p)
	}))
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	client := anthropic.NewClient(
		option.WithBaseURL(apiServer.URL),
		option.WithAPIKey("sk-test-key"),
		option.WithConfigQuiet(&config.Config{
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

	if receivedHeaders.Get("X-Api-Key") != "sk-test-key" {
		t.Fatalf("got X-Api-Key %q, want %q (API key precedence must be unchanged)",
			receivedHeaders.Get("X-Api-Key"), "sk-test-key")
	}

	logMu.Lock()
	out := logBuf.String()
	logMu.Unlock()
	if out != "" {
		t.Fatalf("expected no log output from WithConfigQuiet, got: %q", out)
	}
}

// TestIntegration_WithoutEnvironmentDefaults_NoAutoloadShadowWarning proves
// that when the environment would otherwise autoload a profile (via
// ANTHROPIC_PROFILE), passing option.WithoutEnvironmentDefaults() causes
// NewClient to skip DefaultClientOptions entirely so no autoloaded
// WithConfig is added — and therefore an explicit WithAPIKey produces zero
// shadow warnings. This closes the gap WithConfigQuiet alone leaves: it
// quiets the caller's own WithConfig but not the SDK's autoloaded one.
func TestIntegration_WithoutEnvironmentDefaults_NoAutoloadShadowWarning(t *testing.T) {
	unsetEnv(t, "ANTHROPIC_API_KEY")
	unsetEnv(t, "ANTHROPIC_AUTH_TOKEN")
	isolateAuthEnv(t)
	auth.ResetWarnOnceForTest()

	// Stage a profile under ANTHROPIC_CONFIG_DIR and point ANTHROPIC_PROFILE
	// at it so DefaultClientOptions, if invoked, would inject a WithConfig.
	cfgDir := os.Getenv("ANTHROPIC_CONFIG_DIR")
	credPath := config.ProfileCredentialsPath(cfgDir, "work")
	if err := config.WriteCredentials(credPath, config.Credentials{AccessToken: "profile-tok"}); err != nil {
		t.Fatal(err)
	}
	if err := config.SaveProfile(cfgDir, "work", &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeUserOAuth,
			CredentialsPath: credPath,
			UserOAuth:       &config.UserOAuth{},
		},
	}); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ANTHROPIC_PROFILE", "work")

	var receivedHeaders http.Header
	apiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedHeaders = r.Header
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(successResponse))
	}))
	defer apiServer.Close()

	var logMu sync.Mutex
	var logBuf bytes.Buffer
	log.SetOutput(writerFunc(func(p []byte) (int, error) {
		logMu.Lock()
		defer logMu.Unlock()
		return logBuf.Write(p)
	}))
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	client := anthropic.NewClient(
		option.WithoutEnvironmentDefaults(),
		option.WithBaseURL(apiServer.URL),
		option.WithAPIKey("sk-test-key"),
	)

	if _, err := client.Messages.New(context.Background(), defaultParams); err != nil {
		t.Fatal(err)
	}

	if receivedHeaders.Get("X-Api-Key") != "sk-test-key" {
		t.Fatalf("got X-Api-Key %q, want sk-test-key", receivedHeaders.Get("X-Api-Key"))
	}
	if receivedHeaders.Get("Authorization") != "" {
		t.Fatalf("autoloaded profile credential leaked: Authorization=%q", receivedHeaders.Get("Authorization"))
	}

	logMu.Lock()
	out := logBuf.String()
	logMu.Unlock()
	if out != "" {
		t.Fatalf("expected zero log output with WithoutEnvironmentDefaults, got: %q", out)
	}
}

type writerFunc func([]byte) (int, error)

func (f writerFunc) Write(p []byte) (int, error) { return f(p) }

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

// TestDefaultClient_ExplicitAPIKeyOverridesBrokenProfile verifies that a
// caller-supplied option.WithAPIKey preempts the broken-profile error
// from explicitProfileErrorOption. The user is asking for an explicit
// credential override; the shell's ANTHROPIC_PROFILE should not strand them.
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
	if got := receivedHeaders.Get("X-Api-Key"); got != "sk-test-key" {
		t.Errorf("X-Api-Key: got %q, want %q", got, "sk-test-key")
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

// TestIntegration_APIKeyShadowWarning_InvertedOrder locks in that the
// shadow warning fires regardless of the order the caller passes
// option.WithConfig and option.WithAPIKey. The earlier implementation
// read r.APIKey at WithConfig's Apply time, which missed the case where
// WithAPIKey was passed AFTER WithConfig in the arg list.
func TestIntegration_APIKeyShadowWarning_InvertedOrder(t *testing.T) {
	auth.ResetWarnOnceForTest()
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
	credPath := filepath.Join(dir, "credentials.json")
	os.WriteFile(credPath, []byte(`{"type":"oauth_token","access_token":"profile-tok"}`), 0600)

	var logMu sync.Mutex
	var logBuf bytes.Buffer
	log.SetOutput(writerFunc(func(p []byte) (int, error) {
		logMu.Lock()
		defer logMu.Unlock()
		return logBuf.Write(p)
	}))
	t.Cleanup(func() { log.SetOutput(os.Stderr) })
	readLog := func() string {
		logMu.Lock()
		defer logMu.Unlock()
		return logBuf.String()
	}

	// Note: WithConfig is passed BEFORE WithAPIKey. The warning must
	// still fire because the final request-time state has both.
	client := anthropic.NewClient(
		option.WithBaseURL(apiServer.URL),
		option.WithConfig(&config.Config{
			AuthenticationInfo: &config.AuthenticationInfo{
				Type:            config.AuthenticationTypeUserOAuth,
				CredentialsPath: credPath,
				UserOAuth:       &config.UserOAuth{},
			},
		}),
		option.WithAPIKey("sk-test-key"),
	)

	if _, err := client.Messages.New(context.Background(), defaultParams); err != nil {
		t.Fatal(err)
	}

	if receivedHeaders.Get("X-Api-Key") != "sk-test-key" {
		t.Fatalf("got X-Api-Key %q, want sk-test-key (API key should still win)", receivedHeaders.Get("X-Api-Key"))
	}

	out := readLog()
	if !strings.Contains(out, "ANTHROPIC_API_KEY") {
		t.Fatalf("expected warning mentioning ANTHROPIC_API_KEY regardless of option order, got: %q", out)
	}
	if !strings.Contains(out, "profile") && !strings.Contains(out, "config") {
		t.Fatalf("expected warning to mention profile/config being shadowed, got: %q", out)
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

// TestWithConfig_ExplicitAPIKeyOverridesBrokenProfile verifies that a
// WithAPIKey passed alongside a WithConfig whose credentials cannot be
// resolved still authenticates the request — i.e. the resolution error is
// deferred to a request-time middleware with tier-1 escape hatches, not
// returned at option-Apply time where it would short-circuit later options.
func TestWithConfig_ExplicitAPIKeyOverridesBrokenProfile(t *testing.T) {
	auth.ResetWarnOnceForTest()
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
		option.WithAPIKey("sk-x"),
	)
	if _, err := client.Messages.New(context.Background(), defaultParams); err != nil {
		t.Fatalf("expected request to succeed via WithAPIKey despite broken profile, got: %v", err)
	}
	if got := receivedHeaders.Get("X-Api-Key"); got != "sk-x" {
		t.Fatalf("got X-Api-Key %q, want sk-x", got)
	}
	if receivedHeaders.Get("Authorization") != "" {
		t.Fatalf("Authorization should not be set when X-Api-Key is, got %q", receivedHeaders.Get("Authorization"))
	}
}

// TestWithConfig_ExplicitAuthTokenOverridesBrokenProfile is the WithAuthToken
// variant of the above.
func TestWithConfig_ExplicitAuthTokenOverridesBrokenProfile(t *testing.T) {
	auth.ResetWarnOnceForTest()
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
		option.WithAuthToken("tok-x"),
	)
	if _, err := client.Messages.New(context.Background(), defaultParams); err != nil {
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
	auth.ResetWarnOnceForTest()
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
		option.WithAPIKey("sk-x"),
	)
	if _, err := client.Messages.New(context.Background(), defaultParams); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hit.Load() {
		t.Fatal("request did not reach the profile's BaseURL")
	}
}

// TestWithConfig_BrokenProfileNoOverrideStillFails verifies the negative:
// when no tier-1 credential preempts the profile, the deferred resolution
// error surfaces on the first request.
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
