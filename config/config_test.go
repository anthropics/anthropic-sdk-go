package config_test

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/config"
)

// setupConfigDir creates a temp config dir with configs/<profile>.json written
// from cfgData. Returns the config dir path.
func setupConfigDir(t *testing.T, profile string, cfgData map[string]any) string {
	t.Helper()
	dir := t.TempDir()
	configsDir := filepath.Join(dir, "configs")
	os.MkdirAll(configsDir, 0755)
	b, _ := json.MarshalIndent(cfgData, "", "  ")
	os.WriteFile(filepath.Join(configsDir, profile+".json"), b, 0644)
	return dir
}

// unsetEnv removes an environment variable for the duration of a test. It
// uses t.Setenv first to register the original value for automatic
// restoration at test end, then calls os.Unsetenv to actually remove it.
func unsetEnv(t *testing.T, key string) {
	t.Helper()
	t.Setenv(key, "")
	os.Unsetenv(key)
}

// oidcFederationProfile is a minimal valid oidc_federation config body.
func oidcFederationProfile() map[string]any {
	return map[string]any{
		"authentication": map[string]any{
			"type":               "oidc_federation",
			"federation_rule_id": "fdrl_test",
			"identity_token": map[string]any{
				"source": "file",
				"path":   "/tmp/token",
			},
		},
	}
}

// userOAuthProfile is a minimal valid user_oauth config body.
func userOAuthProfile() map[string]any {
	return map[string]any{
		"authentication": map[string]any{
			"type": "user_oauth",
		},
	}
}

func TestLoadConfig_ProfileFromEnv(t *testing.T) {
	body := oidcFederationProfile()
	body["base_url"] = "https://work.example.com"
	dir := setupConfigDir(t, "work", body)
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "work")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AuthenticationInfo.Type != config.AuthenticationTypeOIDCFederation {
		t.Fatalf("got type %q, want %q", cfg.AuthenticationInfo.Type, config.AuthenticationTypeOIDCFederation)
	}
	if cfg.BaseURL != "https://work.example.com" {
		t.Fatalf("got base_url %q, want %q", cfg.BaseURL, "https://work.example.com")
	}
}

func TestLoadConfig_ProfileFromActiveConfigFile(t *testing.T) {
	dir := setupConfigDir(t, "staging", userOAuthProfile())
	os.WriteFile(filepath.Join(dir, "active_config"), []byte("staging"), 0644)
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	unsetEnv(t, "ANTHROPIC_PROFILE")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AuthenticationInfo.Type != config.AuthenticationTypeUserOAuth {
		t.Fatalf("got type %q, want %q", cfg.AuthenticationInfo.Type, config.AuthenticationTypeUserOAuth)
	}
}

func TestLoadConfig_DefaultProfile(t *testing.T) {
	dir := setupConfigDir(t, "default", userOAuthProfile())
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	unsetEnv(t, "ANTHROPIC_PROFILE")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AuthenticationInfo.Type != config.AuthenticationTypeUserOAuth {
		t.Fatalf("got type %q, want %q", cfg.AuthenticationInfo.Type, config.AuthenticationTypeUserOAuth)
	}
}

func TestLoadConfig_EnvProfileTakesPrecedenceOverActiveConfig(t *testing.T) {
	envBody := oidcFederationProfile()
	envBody["base_url"] = "https://from-env.example.com"
	dir := setupConfigDir(t, "from-env", envBody)

	fileBody := oidcFederationProfile()
	fileBody["base_url"] = "https://from-file.example.com"
	setupB, _ := json.MarshalIndent(fileBody, "", "  ")
	os.MkdirAll(filepath.Join(dir, "configs"), 0755)
	os.WriteFile(filepath.Join(dir, "configs", "from-file.json"), setupB, 0644)
	os.WriteFile(filepath.Join(dir, "active_config"), []byte("from-file"), 0644)

	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "from-env")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BaseURL != "https://from-env.example.com" {
		t.Fatalf("got base_url %q, want %q", cfg.BaseURL, "https://from-env.example.com")
	}
}

func TestLoadConfig_OIDCFederationAllFields(t *testing.T) {
	dir := setupConfigDir(t, "default", map[string]any{
		"base_url":        "https://api.example.com",
		"organization_id": "org-123",
		"workspace_id":    "wrkspc_456",
		"authentication": map[string]any{
			"type":               "oidc_federation",
			"credentials_path":   "/custom/creds.json",
			"federation_rule_id": "fdrl_789",
			"service_account_id": "svac_abc",
			"identity_token": map[string]any{
				"source": "file",
				"path":   "/tmp/token",
			},
		},
	})
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "default")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AuthenticationInfo.Type != config.AuthenticationTypeOIDCFederation {
		t.Fatalf("got type %q, want %q", cfg.AuthenticationInfo.Type, config.AuthenticationTypeOIDCFederation)
	}
	if cfg.BaseURL != "https://api.example.com" {
		t.Fatalf("got base_url %q", cfg.BaseURL)
	}
	if cfg.OrganizationID != "org-123" {
		t.Fatalf("got organization_id %q", cfg.OrganizationID)
	}
	if cfg.WorkspaceID != "wrkspc_456" {
		t.Fatalf("got workspace_id %q", cfg.WorkspaceID)
	}
	if cfg.AuthenticationInfo.CredentialsPath != "/custom/creds.json" {
		t.Fatalf("got credentials_path %q", cfg.AuthenticationInfo.CredentialsPath)
	}
	oidc := cfg.AuthenticationInfo.OIDCFederation
	if oidc == nil {
		t.Fatal("expected oidc_federation sub-object")
	}
	if oidc.FederationRuleID != "fdrl_789" {
		t.Fatalf("got federation_rule_id %q", oidc.FederationRuleID)
	}
	if oidc.ServiceAccountID != "svac_abc" {
		t.Fatalf("got service_account_id %q", oidc.ServiceAccountID)
	}
	if oidc.IdentityToken == nil || oidc.IdentityToken.Source != config.IdentityTokenSourceFile {
		t.Fatalf("got identity_token %+v", oidc.IdentityToken)
	}
	if oidc.IdentityToken.Path != "/tmp/token" {
		t.Fatalf("got identity_token.path %q", oidc.IdentityToken.Path)
	}
}

func TestLoadConfig_UserOAuthWithClientID(t *testing.T) {
	body := userOAuthProfile()
	body["authentication"].(map[string]any)["client_id"] = "client-xyz"
	dir := setupConfigDir(t, "default", body)
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "default")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AuthenticationInfo.UserOAuth == nil {
		t.Fatal("expected user_oauth sub-object")
	}
	if cfg.AuthenticationInfo.UserOAuth.ClientID != "client-xyz" {
		t.Fatalf("got client_id %q", cfg.AuthenticationInfo.UserOAuth.ClientID)
	}
}

func TestLoadConfig_DefaultCredentialsPath(t *testing.T) {
	dir := setupConfigDir(t, "myprofile", userOAuthProfile())
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "myprofile")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, "credentials", "myprofile.json")
	if cfg.AuthenticationInfo.CredentialsPath != want {
		t.Fatalf("got credentials_path %q, want %q", cfg.AuthenticationInfo.CredentialsPath, want)
	}
}

func TestLoadConfig_ExplicitCredentialsPathNotOverridden(t *testing.T) {
	body := userOAuthProfile()
	body["authentication"].(map[string]any)["credentials_path"] = "/explicit/path.json"
	dir := setupConfigDir(t, "default", body)
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "default")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AuthenticationInfo.CredentialsPath != "/explicit/path.json" {
		t.Fatalf("got credentials_path %q, want %q", cfg.AuthenticationInfo.CredentialsPath, "/explicit/path.json")
	}
}

func TestLoadConfig_MissingConfigFile(t *testing.T) {
	t.Setenv("ANTHROPIC_CONFIG_DIR", "/nonexistent/dir")
	t.Setenv("ANTHROPIC_PROFILE", "default")

	_, err := config.LoadConfig()
	if err == nil {
		t.Fatal("expected error for missing config file")
	}
}

func TestLoadConfig_EmptyActiveConfigFallsBackToDefault(t *testing.T) {
	dir := setupConfigDir(t, "default", userOAuthProfile())
	os.WriteFile(filepath.Join(dir, "active_config"), []byte("  \n  "), 0644)
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	unsetEnv(t, "ANTHROPIC_PROFILE")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AuthenticationInfo.Type != config.AuthenticationTypeUserOAuth {
		t.Fatalf("got type %q", cfg.AuthenticationInfo.Type)
	}
}

func TestLoadConfig_MissingAuthenticationInfo(t *testing.T) {
	dir := setupConfigDir(t, "default", map[string]any{
		"base_url": "https://api.example.com",
	})
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "default")

	_, err := config.LoadConfig()
	if err == nil {
		t.Fatal("expected error for missing authentication")
	}
}

// TestLoadConfig_CrossVariantFieldTolerated asserts that a field belonging
// to the other variant (e.g. client_id on an oidc_federation body) is
// silently ignored per the credentials-file-format spec's tolerance rule.
// The declared variant still decodes cleanly.
func TestLoadConfig_CrossVariantFieldTolerated(t *testing.T) {
	dir := setupConfigDir(t, "default", map[string]any{
		"authentication": map[string]any{
			"type":               "oidc_federation",
			"federation_rule_id": "fdrl_test",
			"client_id":          "ignored-cross-variant",
		},
	})
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "default")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AuthenticationInfo.Type != config.AuthenticationTypeOIDCFederation {
		t.Errorf("type: %q", cfg.AuthenticationInfo.Type)
	}
	if cfg.AuthenticationInfo.OIDCFederation.FederationRuleID != "fdrl_test" {
		t.Errorf("FederationRuleID: %q", cfg.AuthenticationInfo.OIDCFederation.FederationRuleID)
	}
	if cfg.AuthenticationInfo.UserOAuth != nil {
		t.Errorf("UserOAuth should be nil, got %+v", cfg.AuthenticationInfo.UserOAuth)
	}
}

// TestLoadConfig_BothVariantFieldsTolerated asserts that a payload carrying
// fields from both variants at once is accepted: the declared variant's
// fields populate its sub-struct and the other variant's fields are
// ignored.
func TestLoadConfig_BothVariantFieldsTolerated(t *testing.T) {
	dir := setupConfigDir(t, "default", map[string]any{
		"authentication": map[string]any{
			"type":               "user_oauth",
			"federation_rule_id": "ignored-fdrl",
			"client_id":          "client-xyz",
		},
	})
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "default")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AuthenticationInfo.UserOAuth == nil || cfg.AuthenticationInfo.UserOAuth.ClientID != "client-xyz" {
		t.Errorf("UserOAuth.ClientID: %+v", cfg.AuthenticationInfo.UserOAuth)
	}
	if cfg.AuthenticationInfo.OIDCFederation != nil {
		t.Errorf("OIDCFederation should be nil, got %+v", cfg.AuthenticationInfo.OIDCFederation)
	}
}

func TestLoadConfig_UnknownAuthenticationType(t *testing.T) {
	dir := setupConfigDir(t, "default", map[string]any{
		"authentication": map[string]any{
			"type": "totally_made_up",
		},
	})
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "default")

	_, err := config.LoadConfig()
	if err == nil {
		t.Fatal("expected error for unknown authentication type")
	}
}

func TestLoadConfig_RejectsProfileWithPathTraversal(t *testing.T) {
	dir := setupConfigDir(t, "default", userOAuthProfile())
	// Write a config file at the location the traversal path would resolve
	// to, so the only thing stopping the load is the validator.
	os.MkdirAll(filepath.Join(dir, "configs"), 0755)
	evilBody := userOAuthProfile()
	evilB, _ := json.Marshal(evilBody)
	os.WriteFile(filepath.Join(dir, "configs", "evil.json"), evilB, 0644)

	for _, bad := range []string{
		"../etc/passwd",
		"..",
		"../../evil",
		`nested\name`,
		"with/slash",
		"",
	} {
		t.Run(bad, func(t *testing.T) {
			t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
			t.Setenv("ANTHROPIC_PROFILE", bad)
			if _, err := config.LoadConfig(); err == nil {
				t.Fatalf("expected error for profile %q", bad)
			}
		})
	}
}

// TestLoadConfig_ToleratesUnknownTopLevelFields covers the spec's tolerance
// rule on the top-level Config object.
func TestLoadConfig_ToleratesUnknownTopLevelFields(t *testing.T) {
	body := oidcFederationProfile()
	body["future_field"] = "ignored-by-this-sdk-version"
	body["_comment"] = "human-readable annotation"
	dir := setupConfigDir(t, "default", body)
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "default")
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AuthenticationInfo.Type != config.AuthenticationTypeOIDCFederation {
		t.Errorf("type: %q", cfg.AuthenticationInfo.Type)
	}
}

// TestLoadConfig_ToleratesUnknownVariantFields covers the tolerance rule
// inside AuthenticationInfo.UnmarshalJSON. A typo'd field name (e.g.
// "federation_rule" instead of "federation_rule_id") is silently ignored
// — the typo surfaces as a missing-required-field error at credential
// resolution time, not a parse error.
func TestLoadConfig_ToleratesUnknownVariantFields(t *testing.T) {
	dir := setupConfigDir(t, "default", map[string]any{
		"authentication": map[string]any{
			"type":               "oidc_federation",
			"federation_rule_id": "fdrl_ok",
			"federation_rule":    "typo-ignored",
			"_comment":           "per-spec tolerance",
		},
	})
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "default")
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AuthenticationInfo.OIDCFederation.FederationRuleID != "fdrl_ok" {
		t.Errorf("FederationRuleID: %q", cfg.AuthenticationInfo.OIDCFederation.FederationRuleID)
	}
}

// TestLoadConfig_UnknownVariantFieldLogsWarning verifies that the tolerance
// rule also logs a warn-once message naming the unknown field, so a user
// whose typo silently vanished can at least see that it happened.
func TestLoadConfig_UnknownVariantFieldLogsWarning(t *testing.T) {
	dir := setupConfigDir(t, "default", map[string]any{
		"authentication": map[string]any{
			"type":               "oidc_federation",
			"federation_rule_id": "fdrl_ok",
			"federaton_rule_id":  "typo-ignored",
		},
	})
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "default")
	config.ResetConfigWarnOnceForTest()

	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "federaton_rule_id") {
		t.Errorf("expected warning mentioning unknown field, got: %q", out)
	}
	// Regression check: the correctly-spelled field must still load through
	// even when an unknown field is alongside it.
	if cfg.AuthenticationInfo == nil || cfg.AuthenticationInfo.OIDCFederation == nil ||
		cfg.AuthenticationInfo.OIDCFederation.FederationRuleID != "fdrl_ok" {
		t.Errorf("correct field dropped: got %+v", cfg.AuthenticationInfo)
	}
}

// TestAuthenticationInfo_RoundTrip marshals a populated AuthenticationInfo
// for each variant and asserts that unmarshalling the result reproduces the
// original.
func TestAuthenticationInfo_RoundTrip(t *testing.T) {
	cases := []struct {
		name  string
		value config.AuthenticationInfo
	}{
		{
			name: "oidc_federation_all_fields",
			value: config.AuthenticationInfo{
				Type:            config.AuthenticationTypeOIDCFederation,
				CredentialsPath: "/tmp/creds.json",
				OIDCFederation: &config.OIDCFederation{
					FederationRuleID: "fdrl_abc",
					ServiceAccountID: "svac_def",
					IdentityToken: &config.IdentityTokenConfig{
						Source: config.IdentityTokenSourceFile,
						Path:   "/tmp/token",
					},
					Scope: "user:inference",
				},
			},
		},
		{
			name: "oidc_federation_minimal",
			value: config.AuthenticationInfo{
				Type: config.AuthenticationTypeOIDCFederation,
				OIDCFederation: &config.OIDCFederation{
					FederationRuleID: "fdrl_min",
				},
			},
		},
		{
			name: "user_oauth_with_client_id_and_scope",
			value: config.AuthenticationInfo{
				Type:            config.AuthenticationTypeUserOAuth,
				CredentialsPath: "/tmp/creds.json",
				UserOAuth: &config.UserOAuth{
					ClientID: "client-xyz",
					Scope:    "user:profile user:inference",
				},
			},
		},
		{
			name: "user_oauth_with_console_url",
			value: config.AuthenticationInfo{
				Type: config.AuthenticationTypeUserOAuth,
				UserOAuth: &config.UserOAuth{
					ClientID:   "client-xyz",
					ConsoleURL: "https://staging.example/console",
				},
			},
		},
		{
			name: "user_oauth_static",
			value: config.AuthenticationInfo{
				Type:      config.AuthenticationTypeUserOAuth,
				UserOAuth: &config.UserOAuth{},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.value)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var got config.AuthenticationInfo
			if err := json.Unmarshal(b, &got); err != nil {
				t.Fatalf("unmarshal %s: %v", string(b), err)
			}
			// Compare field by field so we catch pointer-equality traps.
			if got.Type != tc.value.Type {
				t.Errorf("Type: got %q, want %q", got.Type, tc.value.Type)
			}
			if got.CredentialsPath != tc.value.CredentialsPath {
				t.Errorf("CredentialsPath: got %q, want %q", got.CredentialsPath, tc.value.CredentialsPath)
			}
			switch tc.value.Type {
			case config.AuthenticationTypeOIDCFederation:
				if got.OIDCFederation == nil {
					t.Fatal("OIDCFederation nil after round-trip")
				}
				if got.OIDCFederation.FederationRuleID != tc.value.OIDCFederation.FederationRuleID {
					t.Errorf("FederationRuleID: got %q, want %q",
						got.OIDCFederation.FederationRuleID, tc.value.OIDCFederation.FederationRuleID)
				}
				if got.OIDCFederation.ServiceAccountID != tc.value.OIDCFederation.ServiceAccountID {
					t.Errorf("ServiceAccountID: got %q, want %q",
						got.OIDCFederation.ServiceAccountID, tc.value.OIDCFederation.ServiceAccountID)
				}
				if got.OIDCFederation.Scope != tc.value.OIDCFederation.Scope {
					t.Errorf("Scope: got %q, want %q",
						got.OIDCFederation.Scope, tc.value.OIDCFederation.Scope)
				}
				wantIT := tc.value.OIDCFederation.IdentityToken
				gotIT := got.OIDCFederation.IdentityToken
				if (wantIT == nil) != (gotIT == nil) {
					t.Fatalf("IdentityToken presence: got %v, want %v", gotIT, wantIT)
				}
				if wantIT != nil && (gotIT.Source != wantIT.Source || gotIT.Path != wantIT.Path) {
					t.Errorf("IdentityToken: got %+v, want %+v", gotIT, wantIT)
				}
				if got.UserOAuth != nil {
					t.Errorf("UserOAuth should be nil, got %+v", got.UserOAuth)
				}
			case config.AuthenticationTypeUserOAuth:
				if got.UserOAuth == nil {
					t.Fatal("UserOAuth nil after round-trip")
				}
				if got.UserOAuth.ClientID != tc.value.UserOAuth.ClientID {
					t.Errorf("ClientID: got %q, want %q",
						got.UserOAuth.ClientID, tc.value.UserOAuth.ClientID)
				}
				if got.UserOAuth.Scope != tc.value.UserOAuth.Scope {
					t.Errorf("Scope: got %q, want %q",
						got.UserOAuth.Scope, tc.value.UserOAuth.Scope)
				}
				if got.OIDCFederation != nil {
					t.Errorf("OIDCFederation should be nil, got %+v", got.OIDCFederation)
				}
			}
		})
	}
}

// TestAuthenticationInfo_FlatWireShape pins the exact JSON emitted by
// MarshalJSON so a future refactor can't silently re-introduce nesting.
func TestAuthenticationInfo_FlatWireShape(t *testing.T) {
	value := config.AuthenticationInfo{
		Type:            config.AuthenticationTypeOIDCFederation,
		CredentialsPath: "/tmp/creds.json",
		OIDCFederation: &config.OIDCFederation{
			FederationRuleID: "fdrl_abc",
			Scope:            "user:inference",
		},
	}
	b, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var asMap map[string]any
	if err := json.Unmarshal(b, &asMap); err != nil {
		t.Fatal(err)
	}
	// Variant-specific field must sit at the top level, not nested.
	if _, nested := asMap["oidc_federation"]; nested {
		t.Errorf("output nests oidc_federation key: %s", string(b))
	}
	if asMap["federation_rule_id"] != "fdrl_abc" {
		t.Errorf("expected flat federation_rule_id, got: %s", string(b))
	}
	if asMap["type"] != "oidc_federation" {
		t.Errorf("missing or wrong type: %s", string(b))
	}
	if asMap["credentials_path"] != "/tmp/creds.json" {
		t.Errorf("missing credentials_path: %s", string(b))
	}
	if asMap["scope"] != "user:inference" {
		t.Errorf("expected flat scope at auth level, got: %s", string(b))
	}
}

// TestLoadConfig_EnvDoesNotOverrideProfileFields locks in the spec's
// fill-missing-only semantics: if a profile declares a field, an env var of
// the same name must NOT silently replace it. This keeps an explicit
// profile authoritative on a machine that also exports WIF env vars.
func TestLoadConfig_EnvDoesNotOverrideProfileFields(t *testing.T) {
	body := oidcFederationProfile()
	body["base_url"] = "https://file.example.com"
	body["organization_id"] = "org_from_file"
	body["workspace_id"] = "wrkspc_from_file"
	auth := body["authentication"].(map[string]any)
	auth["federation_rule_id"] = "fdrl_from_file"
	auth["service_account_id"] = "svac_from_file"
	auth["scope"] = "scope_from_file"
	auth["identity_token"] = map[string]any{
		"source": "file",
		"path":   "/tmp/file-token",
	}

	dir := setupConfigDir(t, "default", body)
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	unsetEnv(t, "ANTHROPIC_PROFILE")
	t.Setenv("ANTHROPIC_BASE_URL", "https://env.example.com")
	t.Setenv("ANTHROPIC_ORGANIZATION_ID", "org_from_env")
	t.Setenv("ANTHROPIC_WORKSPACE_ID", "wrkspc_from_env")
	t.Setenv("ANTHROPIC_FEDERATION_RULE_ID", "fdrl_from_env")
	t.Setenv("ANTHROPIC_SERVICE_ACCOUNT_ID", "svac_from_env")
	t.Setenv("ANTHROPIC_SCOPE", "scope_from_env")
	t.Setenv("ANTHROPIC_IDENTITY_TOKEN_FILE", "/tmp/env-token")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BaseURL != "https://file.example.com" {
		t.Errorf("base_url: got %q, want file value (env must not override)", cfg.BaseURL)
	}
	if cfg.OrganizationID != "org_from_file" {
		t.Errorf("organization_id: got %q, want file value", cfg.OrganizationID)
	}
	if cfg.WorkspaceID != "wrkspc_from_file" {
		t.Errorf("workspace_id: got %q, want file value", cfg.WorkspaceID)
	}
	oidc := cfg.AuthenticationInfo.OIDCFederation
	if oidc.FederationRuleID != "fdrl_from_file" {
		t.Errorf("federation_rule_id: got %q, want file value", oidc.FederationRuleID)
	}
	if oidc.ServiceAccountID != "svac_from_file" {
		t.Errorf("service_account_id: got %q, want file value", oidc.ServiceAccountID)
	}
	if oidc.Scope != "scope_from_file" {
		t.Errorf("scope: got %q, want file value", oidc.Scope)
	}
	if oidc.IdentityToken == nil || oidc.IdentityToken.Path != "/tmp/file-token" {
		t.Errorf("identity_token: got %+v, want file value", oidc.IdentityToken)
	}
}

// TestLoadConfig_EnvFillsMissingFields covers the other half of the
// fill-missing rule: env vars populate fields the profile OMITTED. A
// profile can declare federation_rule_id and organization_id once on disk
// and let each pod inject its own identity_token via env.
func TestLoadConfig_EnvFillsMissingFields(t *testing.T) {
	// Profile omits every field the env can fill.
	body := map[string]any{
		"authentication": map[string]any{
			"type": "oidc_federation",
		},
	}
	dir := setupConfigDir(t, "default", body)
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	unsetEnv(t, "ANTHROPIC_PROFILE")
	t.Setenv("ANTHROPIC_ORGANIZATION_ID", "org_from_env")
	t.Setenv("ANTHROPIC_WORKSPACE_ID", "wrkspc_from_env")
	t.Setenv("ANTHROPIC_FEDERATION_RULE_ID", "fdrl_from_env")
	t.Setenv("ANTHROPIC_SERVICE_ACCOUNT_ID", "svac_from_env")
	t.Setenv("ANTHROPIC_SCOPE", "scope_from_env")
	t.Setenv("ANTHROPIC_IDENTITY_TOKEN_FILE", "/tmp/env-token")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.OrganizationID != "org_from_env" {
		t.Errorf("organization_id: got %q, want env fill-in", cfg.OrganizationID)
	}
	if cfg.WorkspaceID != "wrkspc_from_env" {
		t.Errorf("workspace_id: got %q, want env fill-in", cfg.WorkspaceID)
	}
	oidc := cfg.AuthenticationInfo.OIDCFederation
	if oidc.FederationRuleID != "fdrl_from_env" {
		t.Errorf("federation_rule_id: got %q, want env fill-in", oidc.FederationRuleID)
	}
	if oidc.ServiceAccountID != "svac_from_env" {
		t.Errorf("service_account_id: got %q, want env fill-in", oidc.ServiceAccountID)
	}
	if oidc.Scope != "scope_from_env" {
		t.Errorf("scope: got %q, want env fill-in", oidc.Scope)
	}
	if oidc.IdentityToken == nil || oidc.IdentityToken.Path != "/tmp/env-token" {
		t.Errorf("identity_token: got %+v, want env fill-in", oidc.IdentityToken)
	}
}

func TestLoadConfig_EnvFillsMissingUserOAuthScope(t *testing.T) {
	// Profile omits scope entirely.
	body := userOAuthProfile()
	dir := setupConfigDir(t, "default", body)
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	unsetEnv(t, "ANTHROPIC_PROFILE")
	t.Setenv("ANTHROPIC_SCOPE", "scope_from_env")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AuthenticationInfo.UserOAuth == nil {
		t.Fatal("expected user_oauth variant")
	}
	if got := cfg.AuthenticationInfo.UserOAuth.Scope; got != "scope_from_env" {
		t.Errorf("user_oauth scope: got %q, want env fill-in", got)
	}
}

func TestLoadConfig_EnvDoesNotOverrideUserOAuthScope(t *testing.T) {
	body := userOAuthProfile()
	auth := body["authentication"].(map[string]any)
	auth["scope"] = "scope_from_file"

	dir := setupConfigDir(t, "default", body)
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	unsetEnv(t, "ANTHROPIC_PROFILE")
	t.Setenv("ANTHROPIC_SCOPE", "scope_from_env")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AuthenticationInfo.UserOAuth == nil {
		t.Fatal("expected user_oauth variant")
	}
	if got := cfg.AuthenticationInfo.UserOAuth.Scope; got != "scope_from_file" {
		t.Errorf("user_oauth scope: got %q, want file value (env must not override)", got)
	}
}

func TestLoadConfig_EnvOverrideIgnoresEmptyValue(t *testing.T) {
	body := oidcFederationProfile()
	body["base_url"] = "https://file.example.com"
	dir := setupConfigDir(t, "default", body)
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	unsetEnv(t, "ANTHROPIC_PROFILE")
	t.Setenv("ANTHROPIC_BASE_URL", "")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BaseURL != "https://file.example.com" {
		t.Errorf("empty env var should not clear file field, got %q", cfg.BaseURL)
	}
}

func TestLoadConfig_XDGConfigHome(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG not applicable on Windows")
	}
	xdg := t.TempDir()
	configsDir := filepath.Join(xdg, "anthropic", "configs")
	os.MkdirAll(configsDir, 0755)
	body, _ := json.MarshalIndent(oidcFederationProfile(), "", "  ")
	os.WriteFile(filepath.Join(configsDir, "default.json"), body, 0644)

	unsetEnv(t, "ANTHROPIC_CONFIG_DIR")
	unsetEnv(t, "ANTHROPIC_PROFILE")
	t.Setenv("XDG_CONFIG_HOME", xdg)
	t.Setenv("HOME", t.TempDir()) // ensure fallback path is not used

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig with XDG_CONFIG_HOME set: %v", err)
	}
	if cfg.AuthenticationInfo.Type != config.AuthenticationTypeOIDCFederation {
		t.Errorf("expected to load profile from XDG_CONFIG_HOME, got type %q", cfg.AuthenticationInfo.Type)
	}
}
