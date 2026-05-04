package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

// ConfigFileVersion is the version written to configs/<profile>.json. Absent
// on read implies "1.0". Two-part major.minor: bump major on backwards-
// incompatible shape changes (readers can reject), minor on additive
// changes that older readers can tolerate.
const ConfigFileVersion = "1.0"

// profileNamePattern restricts config profile names to characters that can
// never escape the configs/ directory (no path separators, no "..") and
// keeps the name shell-safe.
var profileNamePattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

func validateProfileName(name string) error {
	if name == "" {
		return fmt.Errorf("profile name is empty")
	}
	if name == "." || name == ".." {
		return fmt.Errorf("profile name %q is not allowed", name)
	}
	if strings.ContainsAny(name, `/\`) {
		return fmt.Errorf("profile name %q must not contain path separators", name)
	}
	if !profileNamePattern.MatchString(name) {
		return fmt.Errorf("profile name %q contains disallowed characters (allowed: letters, digits, '_', '.', '-')", name)
	}
	return nil
}

// Config holds the raw configuration for authenticating with the Anthropic
// API. It mirrors the data stored in config files (configs/<profile>.json)
// and can be constructed manually or loaded from disk with [LoadConfig].
//
// Authentication-mode-specific fields live inside [AuthenticationInfo], a
// tagged union discriminated on [AuthenticationInfo.Type]. Top-level fields
// apply to every profile regardless of authentication mode.
type Config struct {
	// Version is the file-format version. Set to [ConfigFileVersion] by
	// [SaveProfile] on every write; absent on disk implies "1.0".
	Version string `json:"version,omitempty"`

	// AuthenticationInfo describes how this profile authenticates. Required.
	AuthenticationInfo *AuthenticationInfo `json:"authentication"`

	// BaseURL overrides the default API base URL. Optional.
	BaseURL string `json:"base_url,omitempty"`

	// OrganizationID is the Anthropic organization the profile targets.
	OrganizationID string `json:"organization_id,omitempty"`

	// WorkspaceID scopes requests to a specific workspace.
	WorkspaceID string `json:"workspace_id,omitempty"`
}

// AuthenticationType is the discriminator for [AuthenticationInfo].
type AuthenticationType string

const (
	// AuthenticationTypeOIDCFederation exchanges a third-party OIDC JWT
	// ("identity token") for a short-lived Anthropic access token via the
	// jwt-bearer grant on /v1/oauth/token.
	AuthenticationTypeOIDCFederation AuthenticationType = "oidc_federation"

	// AuthenticationTypeUserOAuth authenticates with an access token minted
	// through a user-interactive OAuth flow (e.g. `anthropic login`), with
	// optional refresh-token rotation when a ClientID is configured.
	AuthenticationTypeUserOAuth AuthenticationType = "user_oauth"
)

// AuthenticationInfo is a tagged union discriminated on [AuthenticationInfo.Type].
// On the wire it is flat: the top-level JSON object holds `type`, the shared
// `credentials_path`, and the variant-specific fields all at the same level.
// In Go, the variant-specific fields are grouped into strongly-typed sub-structs
// ([OIDCFederation], [UserOAuth]) so callers get type safety and a clean
// non-nil check per variant. Exactly one sub-struct pointer is populated, and
// it must match [AuthenticationInfo.Type].
//
// Because the wire shape is flat, all variant field names share a single
// namespace at the JSON layer — new variants must pick field names that do
// not collide with shared fields or with each other.
//
// [AuthenticationInfo.UnmarshalJSON] silently ignores unknown fields at the
// JSON layer (per the credentials-file-format spec: "SDKs MUST silently
// ignore unrecognized top-level keys in both config and credentials files"
// — the same tolerance rule applies to the nested authentication object).
// Unknown authentication types still fail loud because the SDK has no way
// to meaningfully resolve credentials for an unknown variant.
type AuthenticationInfo struct {
	Type AuthenticationType `json:"-"`

	// CredentialsPath is the path to the credentials JSON file that stores
	// access / refresh tokens on disk. Leave empty to use the default path
	// (credentials/<profile>.json under the config directory). Shared across
	// all authentication types.
	CredentialsPath string `json:"-"`

	// OIDCFederation holds the fields for Type == AuthenticationTypeOIDCFederation.
	// Populated by UnmarshalJSON and inlined by MarshalJSON; never appears as
	// a nested JSON object.
	OIDCFederation *OIDCFederation `json:"-"`

	// UserOAuth holds the fields for Type == AuthenticationTypeUserOAuth.
	// Populated by UnmarshalJSON and inlined by MarshalJSON; never appears as
	// a nested JSON object.
	UserOAuth *UserOAuth `json:"-"`
}

// OIDCFederation configures a profile that authenticates by exchanging a
// third-party OIDC identity token for an Anthropic access token. Its fields
// are inlined into the parent [AuthenticationInfo] on the wire.
type OIDCFederation struct {
	// FederationRuleID is the tagged ID ("fdrl_...") of the OidcFederationRule
	// that governs the exchange. Required.
	FederationRuleID string `json:"federation_rule_id"`

	// ServiceAccountID is an optional expected-target check for federation
	// rules with target_type=SERVICE_ACCOUNT. Must be a "svac_..." tagged ID.
	// Omit for target_type=USER rules, where the principal is derived from
	// the JWT claims.
	ServiceAccountID string `json:"service_account_id,omitempty"`

	// IdentityToken describes how to obtain the OIDC assertion this profile
	// presents at token-exchange time.
	IdentityToken *IdentityTokenConfig `json:"identity_token,omitempty"`

	// Scope is the OAuth scope string (RFC 6749 §3.3 space-delimited form)
	// the profile expects to be granted. It is stored on the profile for
	// display and configuration purposes only — the SDK does NOT send it
	// on the jwt-bearer exchange. The granted scope is determined by the
	// federation rule on the server; IssueOAuthTokenRequest has no scope
	// field, and the REST gateway's alias transformation strips unknown
	// keys, so any attempt to wire this through to the exchange body
	// would be silently dropped. The granted scope appears on the token
	// response's `scope` field, but the SDK does not currently surface
	// it to callers.
	Scope string `json:"scope,omitempty"`
}

// UserOAuth configures a profile authenticated via a user-interactive OAuth
// flow. The access and (optional) refresh tokens live on disk under the
// profile's [AuthenticationInfo.CredentialsPath]. Its fields are inlined
// into the parent [AuthenticationInfo] on the wire.
type UserOAuth struct {
	// ClientID is the OAuth client ID used for refresh-token exchange. When
	// empty, the access token is treated as static (no refresh) and the
	// profile fails once it expires.
	ClientID string `json:"client_id,omitempty"`

	// Scope is the OAuth scope string (RFC 6749 §3.3 space-delimited form)
	// the profile was granted, captured at login time. The SDK does not
	// consult this on refresh — the oauth-server preserves the original
	// scope set when the refresh request omits scope — but `ant auth status`
	// and similar tools display it.
	Scope string `json:"scope,omitempty"`

	// ConsoleURL is the base URL of the OAuth /authorize page for
	// interactive login. The SDK does not consult it; it's CLI-only state
	// (the `ant auth login` target) that lives on the shared schema so
	// SaveProfile/LoadProfile round-trip it without triggering the
	// unknown-field warning.
	ConsoleURL string `json:"console_url,omitempty"`
}

// IdentityTokenSource is the source kind for an OIDC identity token.
type IdentityTokenSource string

const (
	// IdentityTokenSourceFile reads the token from a file on every exchange,
	// which supports rotated tokens (e.g. Kubernetes projected service
	// account tokens).
	IdentityTokenSourceFile IdentityTokenSource = "file"
)

// IdentityTokenConfig specifies how to obtain an OIDC identity token for
// federation exchange.
type IdentityTokenConfig struct {
	Source IdentityTokenSource `json:"source"`
	Path   string              `json:"path,omitempty"`
}

// NewUserOAuthAuthentication returns a populated [AuthenticationInfo] for
// the user_oauth variant. Pass an empty clientID to describe a profile whose
// access token is static (no refresh).
func NewUserOAuthAuthentication(clientID string) *AuthenticationInfo {
	return &AuthenticationInfo{
		Type:      AuthenticationTypeUserOAuth,
		UserOAuth: &UserOAuth{ClientID: clientID},
	}
}

// NewOIDCFederationAuthentication builds an [AuthenticationInfo] for the
// oidc_federation variant.
func NewOIDCFederationAuthentication(oidc OIDCFederation) *AuthenticationInfo {
	return &AuthenticationInfo{
		Type:           AuthenticationTypeOIDCFederation,
		OIDCFederation: &oidc,
	}
}

// sharedAuthFields are the fields present on every authentication variant.
// The discriminator read uses a tolerant decoder (unknown fields allowed)
// because the variant decode in the second pass is what rejects typos.
type sharedAuthFields struct {
	Type            AuthenticationType `json:"type"`
	CredentialsPath string             `json:"credentials_path,omitempty"`
}

// oidcFederationWire is the flat wire shape for oidc_federation: shared
// fields + inlined OIDCFederation fields. Unknown fields are tolerated
// per the credentials-file-format spec so forward-compatible additions
// and cross-SDK-written files don't fail loud.
type oidcFederationWire struct {
	sharedAuthFields
	FederationRuleID string               `json:"federation_rule_id"`
	ServiceAccountID string               `json:"service_account_id,omitempty"`
	IdentityToken    *IdentityTokenConfig `json:"identity_token,omitempty"`
	Scope            string               `json:"scope,omitempty"`
}

// userOAuthWire is the flat wire shape for user_oauth: shared fields +
// inlined UserOAuth fields.
type userOAuthWire struct {
	sharedAuthFields
	ClientID   string `json:"client_id,omitempty"`
	Scope      string `json:"scope,omitempty"`
	ConsoleURL string `json:"console_url,omitempty"`
}

// UnmarshalJSON decodes the flat tagged-union wire shape into the nested
// Go representation. Unknown fields are silently tolerated per the
// credentials-file-format spec but are also logged (warn-once) so a
// user who typo'd a field name at least sees that it was ignored.
// Unknown authentication types still fail loud.
func (a *AuthenticationInfo) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	var authType AuthenticationType
	if b, ok := raw["type"]; ok {
		if err := json.Unmarshal(b, &authType); err != nil {
			return err
		}
	}

	switch authType {
	case "":
		return fmt.Errorf("authentication.type is required")
	case AuthenticationTypeOIDCFederation:
		var w oidcFederationWire
		if err := json.Unmarshal(data, &w); err != nil {
			return fmt.Errorf("authentication.type %q: %w", authType, err)
		}
		unknown := unknownRawKeys(raw, knownOIDCFederationFields)
		*a = AuthenticationInfo{
			Type:            w.Type,
			CredentialsPath: w.CredentialsPath,
			OIDCFederation: &OIDCFederation{
				FederationRuleID: w.FederationRuleID,
				ServiceAccountID: w.ServiceAccountID,
				IdentityToken:    w.IdentityToken,
				Scope:            w.Scope,
			},
		}
		warnUnknownAuthFields(string(authType), unknown)
		return nil
	case AuthenticationTypeUserOAuth:
		var w userOAuthWire
		if err := json.Unmarshal(data, &w); err != nil {
			return fmt.Errorf("authentication.type %q: %w", authType, err)
		}
		unknown := unknownRawKeys(raw, knownUserOAuthFields)
		*a = AuthenticationInfo{
			Type:            w.Type,
			CredentialsPath: w.CredentialsPath,
			UserOAuth: &UserOAuth{
				ClientID:   w.ClientID,
				Scope:      w.Scope,
				ConsoleURL: w.ConsoleURL,
			},
		}
		warnUnknownAuthFields(string(authType), unknown)
		return nil
	default:
		return fmt.Errorf("authentication.type %q is not a known authentication type", authType)
	}
}

func unknownRawKeys(raw map[string]json.RawMessage, known map[string]struct{}) []string {
	var unknown []string
	for k := range raw {
		if _, ok := known[k]; ok {
			continue
		}
		unknown = append(unknown, k)
	}
	return unknown
}

// Keep in sync with [oidcFederationWire] and [userOAuthWire].
var (
	knownOIDCFederationFields = map[string]struct{}{
		"type":               {},
		"credentials_path":   {},
		"federation_rule_id": {},
		"service_account_id": {},
		"identity_token":     {},
		"scope":              {},
	}
	knownUserOAuthFields = map[string]struct{}{
		"type":             {},
		"credentials_path": {},
		"client_id":        {},
		"scope":            {},
		"console_url":      {},
	}
)

var (
	configWarnOnceMu  sync.Mutex
	configWarnOnceSet = map[string]bool{}
)

func configWarnOnce(key, format string, args ...any) {
	configWarnOnceMu.Lock()
	already := configWarnOnceSet[key]
	configWarnOnceSet[key] = true
	configWarnOnceMu.Unlock()
	if already {
		return
	}
	log.Printf("anthropic-sdk-go/config: "+format, args...)
}

func warnUnknownAuthFields(authType string, fields []string) {
	for _, f := range fields {
		configWarnOnce(
			"unknown-auth-field:"+authType+":"+f,
			"profile authentication.%s: unknown field %q (ignored)",
			authType, f,
		)
	}
}

// ResetConfigWarnOnceForTest clears the warn-once dedupe state so tests
// can observe a warning a prior test may already have triggered.
func ResetConfigWarnOnceForTest() {
	configWarnOnceMu.Lock()
	configWarnOnceSet = map[string]bool{}
	configWarnOnceMu.Unlock()
}

// MarshalJSON emits the flat tagged-union wire shape: shared fields, then
// the matching variant's fields inlined at the same level. Returns an error
// when the in-memory state is inconsistent (Type mismatches the populated
// sub-struct, or neither/both are set).
func (a AuthenticationInfo) MarshalJSON() ([]byte, error) {
	switch a.Type {
	case "":
		return nil, fmt.Errorf("AuthenticationInfo.Type is required")
	case AuthenticationTypeOIDCFederation:
		if a.OIDCFederation == nil {
			return nil, fmt.Errorf("AuthenticationInfo.Type=%q requires OIDCFederation", a.Type)
		}
		if a.UserOAuth != nil {
			return nil, fmt.Errorf("AuthenticationInfo.Type=%q must not set UserOAuth", a.Type)
		}
		return json.Marshal(oidcFederationWire{
			sharedAuthFields: sharedAuthFields{Type: a.Type, CredentialsPath: a.CredentialsPath},
			FederationRuleID: a.OIDCFederation.FederationRuleID,
			ServiceAccountID: a.OIDCFederation.ServiceAccountID,
			IdentityToken:    a.OIDCFederation.IdentityToken,
			Scope:            a.OIDCFederation.Scope,
		})
	case AuthenticationTypeUserOAuth:
		if a.UserOAuth == nil {
			return nil, fmt.Errorf("AuthenticationInfo.Type=%q requires UserOAuth", a.Type)
		}
		if a.OIDCFederation != nil {
			return nil, fmt.Errorf("AuthenticationInfo.Type=%q must not set OIDCFederation", a.Type)
		}
		return json.Marshal(userOAuthWire{
			sharedAuthFields: sharedAuthFields{Type: a.Type, CredentialsPath: a.CredentialsPath},
			ClientID:         a.UserOAuth.ClientID,
			Scope:            a.UserOAuth.Scope,
			ConsoleURL:       a.UserOAuth.ConsoleURL,
		})
	default:
		return nil, fmt.Errorf("AuthenticationInfo.Type %q is not a known authentication type", a.Type)
	}
}

// defaultConfigDir returns the platform-specific base directory for Anthropic
// config and credentials. It checks ANTHROPIC_CONFIG_DIR first, then falls
// back to the platform default. On non-Windows platforms XDG_CONFIG_HOME is
// honored per the XDG Base Directory spec before the $HOME/.config fallback.
// Returns an empty string when the platform home directory cannot be
// resolved (e.g., $HOME unset in a minimal container), so callers don't
// silently read/write to a CWD-relative path.
func defaultConfigDir() string {
	if dir, ok := os.LookupEnv("ANTHROPIC_CONFIG_DIR"); ok {
		return dir
	}
	if runtime.GOOS == "windows" {
		if appdata := os.Getenv("APPDATA"); appdata != "" {
			return filepath.Join(appdata, "Anthropic")
		}
		return ""
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "anthropic")
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, ".config", "anthropic")
}

// resolveProfile determines the active profile name.
// Resolution order: ANTHROPIC_PROFILE env -> active_config file -> "default".
func resolveProfile(configDir string) string {
	if profile, ok := os.LookupEnv("ANTHROPIC_PROFILE"); ok {
		return profile
	}
	data, err := os.ReadFile(filepath.Join(configDir, "active_config"))
	if err == nil {
		if name := strings.TrimSpace(string(data)); name != "" {
			return name
		}
	}
	return "default"
}

// loadProfile reads the config file for a specific profile and returns the
// parsed [Config]. The [Config.CredentialsPath] is defaulted to
// credentials/<profile>.json under configDir if not already set.
func loadProfile(configDir, profile string) (*Config, error) {
	if err := validateProfileName(profile); err != nil {
		return nil, fmt.Errorf("invalid profile %q: %w", profile, err)
	}
	configPath := filepath.Join(configDir, "configs", profile+".json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %q: %w", configPath, err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %q: %w", configPath, err)
	}

	if cfg.AuthenticationInfo == nil {
		return nil, fmt.Errorf("config file %q is missing authentication", configPath)
	}

	if cfg.AuthenticationInfo.CredentialsPath == "" {
		cfg.AuthenticationInfo.CredentialsPath = filepath.Join(configDir, "credentials", profile+".json")
	}

	fillMissingFromEnv(&cfg)

	return &cfg, nil
}

const (
	envBaseURL           = "ANTHROPIC_BASE_URL"
	envOrganizationID    = "ANTHROPIC_ORGANIZATION_ID"
	envFederationRuleID  = "ANTHROPIC_FEDERATION_RULE_ID"
	envServiceAccountID  = "ANTHROPIC_SERVICE_ACCOUNT_ID"
	envScope             = "ANTHROPIC_SCOPE"
	envIdentityTokenFile = "ANTHROPIC_IDENTITY_TOKEN_FILE"
)

// fillMissingFromEnv populates cfg fields from ANTHROPIC_* environment
// variables, but only when the corresponding field is empty in the loaded
// profile. Contract lives in the cross-SDK credential precedence spec
// (sdk_credential_precedence.md, section 1 interaction rules): "env vars
// fill in fields omitted from a profile file, but never override present
// ones." This keeps an explicit profile authoritative on a machine that
// also happens to have WIF env vars exported (e.g., a prod profile is
// not silently replaced by a leftover dev federation_rule_id env var).
//
// Empty env values are treated as unset rather than "clear the field" so
// an accidentally-blanked env var can't silently strip a required field.
func fillMissingFromEnv(cfg *Config) {
	if cfg.BaseURL == "" {
		if v, ok := lookupNonEmpty(envBaseURL); ok {
			cfg.BaseURL = v
		}
	}
	if cfg.OrganizationID == "" {
		if v, ok := lookupNonEmpty(envOrganizationID); ok {
			cfg.OrganizationID = v
		}
	}

	switch cfg.AuthenticationInfo.Type {
	case AuthenticationTypeOIDCFederation:
		oidc := cfg.AuthenticationInfo.OIDCFederation
		if oidc == nil {
			return
		}
		if oidc.FederationRuleID == "" {
			if v, ok := lookupNonEmpty(envFederationRuleID); ok {
				oidc.FederationRuleID = v
			}
		}
		if oidc.ServiceAccountID == "" {
			if v, ok := lookupNonEmpty(envServiceAccountID); ok {
				oidc.ServiceAccountID = v
			}
		}
		if oidc.Scope == "" {
			if v, ok := lookupNonEmpty(envScope); ok {
				oidc.Scope = v
			}
		}
		if oidc.IdentityToken == nil {
			if v, ok := lookupNonEmpty(envIdentityTokenFile); ok {
				oidc.IdentityToken = &IdentityTokenConfig{
					Source: IdentityTokenSourceFile,
					Path:   v,
				}
			}
		}
	case AuthenticationTypeUserOAuth:
		if cfg.AuthenticationInfo.UserOAuth == nil {
			return
		}
		if cfg.AuthenticationInfo.UserOAuth.Scope == "" {
			if v, ok := lookupNonEmpty(envScope); ok {
				cfg.AuthenticationInfo.UserOAuth.Scope = v
			}
		}
	}
}

// lookupNonEmpty returns the env var's value and true only when it is set
// and non-empty. An empty value is reported as unset.
func lookupNonEmpty(key string) (string, bool) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return "", false
	}
	return v, true
}

// LoadProfile loads the config for the named profile from the given config
// directory, bypassing ANTHROPIC_PROFILE / active_config resolution. Use
// [DefaultDir] for the standard location. This is the building block CLIs use
// to inspect or operate on a profile other than the currently-active one.
func LoadProfile(dir, profile string) (*Config, error) {
	return loadProfile(dir, profile)
}

// LoadConfig reads the raw configuration from the Anthropic config file system
// (configs/<profile>.json) and returns it without resolving credentials.
// Credential resolution is deferred until the config is passed to
// [option.WithConfig].
//
// The config directory and profile are resolved using the standard resolution
// order (ANTHROPIC_CONFIG_DIR, ANTHROPIC_PROFILE, active_config file, defaults).
func LoadConfig() (*Config, error) {
	configDir := defaultConfigDir()
	profile := resolveProfile(configDir)
	if err := validateProfileName(profile); err != nil {
		return nil, fmt.Errorf("invalid profile name from resolution: %w", err)
	}
	return loadProfile(configDir, profile)
}
