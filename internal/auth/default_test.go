package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/config"
)

// unsetEnv removes an environment variable for the duration of a test.
// t.Setenv registers the original value for automatic restoration at test
// end; os.Unsetenv then actually removes the var for this test.
func unsetEnv(t *testing.T, key string) {
	t.Helper()
	t.Setenv(key, "")
	os.Unsetenv(key)
}

// userOAuthProfileBody returns a minimal user_oauth config map suitable for
// the profile-system tests.
func userOAuthProfileBody() map[string]any {
	return map[string]any{
		"authentication": map[string]any{
			"type": "user_oauth",
		},
	}
}

// setupProfile creates a config dir with configs/<name>.json and optionally
// credentials/<name>.json, returning the config dir path.
func setupProfile(t *testing.T, name string, cfgData map[string]any, creds map[string]any) string {
	t.Helper()
	dir := t.TempDir()
	configsDir := filepath.Join(dir, "configs")
	credsDir := filepath.Join(dir, "credentials")
	os.MkdirAll(configsDir, 0755)
	os.MkdirAll(credsDir, 0700)

	b, _ := json.MarshalIndent(cfgData, "", "  ")
	os.WriteFile(filepath.Join(configsDir, name+".json"), b, 0644)

	if creds != nil {
		b, _ = json.MarshalIndent(creds, "", "  ")
		os.WriteFile(filepath.Join(credsDir, name+".json"), b, 0600)
	}
	return dir
}

func TestLoadConfig_ProfileSystem(t *testing.T) {
	dir := setupProfile(t, "default", userOAuthProfileBody(), map[string]any{
		"type":         "oauth_token",
		"access_token": "profile-tok",
	})
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	result, err := ResolveCredentials(cfg)
	if err != nil {
		t.Fatal(err)
	}
	tok, err := result.Provider(context.Background(), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if tok.Token != "profile-tok" {
		t.Fatalf("got %q, want %q", tok.Token, "profile-tok")
	}
}

func TestLoadConfig_ProfileEnv(t *testing.T) {
	dir := setupProfile(t, "work", userOAuthProfileBody(), map[string]any{
		"type":         "oauth_token",
		"access_token": "work-tok",
	})
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "work")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	result, err := ResolveCredentials(cfg)
	if err != nil {
		t.Fatal(err)
	}
	tok, err := result.Provider(context.Background(), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if tok.Token != "work-tok" {
		t.Fatalf("got %q, want %q", tok.Token, "work-tok")
	}
}

func TestLoadConfig_ActiveConfigFile(t *testing.T) {
	dir := setupProfile(t, "staging", userOAuthProfileBody(), map[string]any{
		"type":         "oauth_token",
		"access_token": "staging-tok",
	})
	os.WriteFile(filepath.Join(dir, "active_config"), []byte("staging"), 0644)
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	result, err := ResolveCredentials(cfg)
	if err != nil {
		t.Fatal(err)
	}
	tok, err := result.Provider(context.Background(), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	if tok.Token != "staging-tok" {
		t.Fatalf("got %q, want %q", tok.Token, "staging-tok")
	}
}

func TestLoadConfig_NothingSetReturnsError(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)

	_, err := config.LoadConfig()
	if err == nil {
		t.Fatal("expected error when no config file exists")
	}
}

func TestEnvCredentials_WorkloadIdentity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(tokenExchangeResponse{AccessToken: "wi-tok"})
	}))
	defer server.Close()

	dir := t.TempDir()
	tokenPath := filepath.Join(dir, "token")
	os.WriteFile(tokenPath, []byte("my-jwt"), 0600)

	// Set config dir to empty dir so no profile config exists.
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv(EnvIdentityTokenFile, tokenPath)
	t.Setenv(EnvFederationRuleID, "rule-1")
	t.Setenv(EnvOrganizationID, "org-1")

	result, _, _ := EnvCredentials()
	if result == nil {
		t.Fatal("expected non-nil")
	}
	tok, err := result.Provider(context.Background(), server.URL, http.DefaultClient.Do)
	if err != nil {
		t.Fatal(err)
	}
	if tok.Token != "wi-tok" {
		t.Fatalf("got %q, want %q", tok.Token, "wi-tok")
	}
}

func TestEnvCredentials_LiteralToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(tokenExchangeResponse{AccessToken: "wi-tok"})
	}))
	defer server.Close()

	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv(EnvIdentityToken, "literal-jwt")
	t.Setenv(EnvFederationRuleID, "rule-1")
	t.Setenv(EnvOrganizationID, "org-1")

	result, _, _ := EnvCredentials()
	if result == nil {
		t.Fatal("expected non-nil")
	}
	tok, err := result.Provider(context.Background(), server.URL, http.DefaultClient.Do)
	if err != nil {
		t.Fatal(err)
	}
	if tok.Token != "wi-tok" {
		t.Fatalf("got %q, want %q", tok.Token, "wi-tok")
	}
}

func TestEnvCredentials_RequiresAllThree(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	// Only federation rule, missing org ID and token.
	t.Setenv(EnvFederationRuleID, "rule-1")

	result, _, _ := EnvCredentials()
	if result != nil {
		t.Fatal("expected nil when org ID and identity token missing")
	}
}

func TestEnvCredentials_NothingSetReturnsNil(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	unsetEnv(t, EnvIdentityTokenFile)
	unsetEnv(t, EnvIdentityToken)
	unsetEnv(t, EnvFederationRuleID)
	unsetEnv(t, EnvOrganizationID)

	result, _, _ := EnvCredentials()
	if result != nil {
		t.Fatal("expected nil when nothing is set")
	}
}
