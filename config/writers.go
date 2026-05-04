package config

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const credentialsTypeOAuthToken = "oauth_token"

// CredentialsFileVersion is the version written to credentials/<profile>.json.
// Absent on read implies "1.0". Major.minor format: bump major on
// backwards-incompatible shape changes, minor on additive changes.
const CredentialsFileVersion = "1.0"

// Credentials is the in-memory form of a credentials/<profile>.json file.
// Not every field is populated for every authentication variant: federation
// grants, for example, do not return a refresh token.
type Credentials struct {
	AccessToken  string
	RefreshToken string
	// ExpiresAt is the absolute expiry time of AccessToken. Nil means the
	// token's lifetime is unknown to the SDK (treated as non-expiring by
	// the in-memory cache).
	ExpiresAt *time.Time

	// Scope, OrganizationUUID, OrganizationName, and AccountEmail record
	// what the current AccessToken was actually granted/minted for. They are
	// written on every login and reflect the token's view of the world at
	// mint time. This is distinct from [Config.OrganizationID] on the profile
	// config side, which is the user's intended target org; the two may
	// diverge if the user is reassigned or the token was minted before the
	// profile was edited.
	Scope            string
	OrganizationUUID string
	OrganizationName string
	// AccountEmail is the email of the account that minted the token, taken
	// from the /v1/oauth/token response's account.email_address.
	AccountEmail string
	// WorkspaceID and WorkspaceName record the workspace the token was
	// bound to at mint time (when the authorization carried one). Stored
	// as the tagged `wrkspc_...` form — the same format the CLI flag,
	// profile config, and anthropic-workspace-id header accept. Sourced
	// from the /v1/oauth/token response's workspace.id. Empty for tokens
	// that aren't workspace-scoped.
	WorkspaceID   string
	WorkspaceName string
}

type credentialsWireShape struct {
	Version          string `json:"version,omitempty"`
	Type             string `json:"type"`
	AccessToken      string `json:"access_token"`
	ExpiresAt        *int64 `json:"expires_at,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	Scope            string `json:"scope,omitempty"`
	OrganizationUUID string `json:"organization_uuid,omitempty"`
	OrganizationName string `json:"organization_name,omitempty"`
	AccountEmail     string `json:"account_email,omitempty"`
	WorkspaceID      string `json:"workspace_id,omitempty"`
	WorkspaceName    string `json:"workspace_name,omitempty"`
}

func (c Credentials) MarshalJSON() ([]byte, error) {
	w := credentialsWireShape{
		Version:          CredentialsFileVersion,
		Type:             credentialsTypeOAuthToken,
		AccessToken:      c.AccessToken,
		RefreshToken:     c.RefreshToken,
		Scope:            c.Scope,
		OrganizationUUID: c.OrganizationUUID,
		OrganizationName: c.OrganizationName,
		AccountEmail:     c.AccountEmail,
		WorkspaceID:      c.WorkspaceID,
		WorkspaceName:    c.WorkspaceName,
	}
	if c.ExpiresAt != nil {
		exp := c.ExpiresAt.Unix()
		w.ExpiresAt = &exp
	}
	return json.Marshal(w)
}

// UnmarshalJSON decodes a credentials/<profile>.json file. Unknown fields
// are silently tolerated per the credentials-file-format spec. A missing
// "type" is treated as equivalent to "oauth_token" — this is a concession
// to interop with external tooling (credential daemons, sidecars) that
// may write plain bearer-token blobs without the discriminator. Every
// SDK writer emits "type" on write, so missing-type files can only come
// from outside the SDK; rejecting them would break valid integrations.
// A type that is set but not "oauth_token" still fails loud.
func (c *Credentials) UnmarshalJSON(data []byte) error {
	var w credentialsWireShape
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}
	if w.Type != "" && w.Type != credentialsTypeOAuthToken {
		return fmt.Errorf("credentials: unknown type %q", w.Type)
	}
	c.AccessToken = w.AccessToken
	c.RefreshToken = w.RefreshToken
	c.Scope = w.Scope
	c.OrganizationUUID = w.OrganizationUUID
	c.OrganizationName = w.OrganizationName
	c.AccountEmail = w.AccountEmail
	c.WorkspaceID = w.WorkspaceID
	c.WorkspaceName = w.WorkspaceName
	c.ExpiresAt = nil
	if w.ExpiresAt != nil {
		exp := time.Unix(*w.ExpiresAt, 0)
		c.ExpiresAt = &exp
	}
	return nil
}

// DefaultDir returns the SDK's default configuration directory (the same
// one [LoadConfig] reads from when ANTHROPIC_CONFIG_DIR is unset). Returns
// an empty string if the platform home directory cannot be resolved — the
// writers in this package surface that as an explicit error.
func DefaultDir() string { return defaultConfigDir() }

// configsSubdir and credentialsSubdir are the directories under the config
// root that hold one JSON file per profile. Defined once here so the layout
// has a single source of truth across path helpers and [ListProfiles].
const (
	configsSubdir     = "configs"
	credentialsSubdir = "credentials"
)

// ProfilesDir returns the directory containing profile config JSON files.
func ProfilesDir(dir string) string { return filepath.Join(dir, configsSubdir) }

// CredentialsDir returns the directory containing per-profile credentials files.
func CredentialsDir(dir string) string { return filepath.Join(dir, credentialsSubdir) }

// ProfilePath returns the path to configs/<profile>.json under dir.
func ProfilePath(dir, profile string) string {
	return filepath.Join(ProfilesDir(dir), profile+".json")
}

// ProfileCredentialsPath returns the path to credentials/<profile>.json under dir.
func ProfileCredentialsPath(dir, profile string) string {
	return filepath.Join(CredentialsDir(dir), profile+".json")
}

// ActiveConfigPath returns the path to the active_config pointer file under dir.
func ActiveConfigPath(dir string) string {
	return filepath.Join(dir, "active_config")
}

// ListProfiles returns the names of all profiles stored under dir's configs
// subdirectory, sorted for stable output. A missing configs directory is
// treated as "no profiles" and returns a nil slice with no error, so callers
// can enumerate a fresh config root without special-casing first run.
func ListProfiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(ProfilesDir(dir))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name, ok := strings.CutSuffix(e.Name(), ".json")
		if !ok || name == "" {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

// Profile filesystem modes. Configs are non-secret and use the "public"
// mode so another UID on the same host (e.g. a sidecar, or `ant profile
// list` running as a different user inside a pod) can read them;
// credentials are secret and use the tighter mode.
//
// These track the credentials-file-format spec's split rationale:
// configs/ (0755 dirs, 0644 files — checkin-safe) vs credentials/
// (0700 dirs, 0600 files — secrets).
const (
	publicDirMode  os.FileMode = 0755
	publicFileMode os.FileMode = 0644
	secretDirMode  os.FileMode = 0700
	secretFileMode os.FileMode = 0600
)

// writeFileAtomic writes to a unique sibling tmp file then renames onto
// the target (last-writer-wins). Tmp name is unique per write so
// concurrent writers don't collide on the temp; mode is Chmod'd
// explicitly so umask can't loosen secrets; tmp and parent dir are
// fsync'd for crash durability; tmp is removed on any failure.
func writeFileAtomic(path string, data []byte, dirMode, fileMode os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, dirMode); err != nil {
		return err
	}

	// Base prefix is for debuggability (which target a leftover tmp belongs to).
	f, err := os.CreateTemp(dir, filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	tmpPath := f.Name()
	// From here on, any failure must remove the tmp file.
	cleanup := func(err error) error {
		os.Remove(tmpPath)
		return err
	}
	if _, err := f.Write(data); err != nil {
		f.Close()
		return cleanup(err)
	}
	// CreateTemp defaults to 0600; Chmod enforces fileMode for both paths.
	if err := f.Chmod(fileMode); err != nil {
		f.Close()
		return cleanup(err)
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return cleanup(err)
	}
	if err := f.Close(); err != nil {
		return cleanup(err)
	}
	// Refuse a pre-planted symlink at the target — would let an attacker
	// redirect the bearer-token write to a path of their choosing.
	if fi, err := os.Lstat(path); err == nil && fi.Mode()&fs.ModeSymlink != 0 {
		return cleanup(fmt.Errorf("refusing to write %s: target is a symlink", path))
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return cleanup(err)
	}
	// fsync the parent directory so the rename is durable across a crash.
	// On Windows this is a no-op / unsupported — the Sync() call may
	// return an error, which we swallow because the rename itself is
	// already visible to other processes.
	if d, err := os.Open(dir); err == nil {
		_ = d.Sync()
		d.Close()
	}
	return nil
}

// validateDirAndProfile rejects empty dirs and unsafe profile names so
// callers that pass zero values or untrusted input get a clean error
// instead of a path-traversal.
func validateDirAndProfile(dir, profile string) error {
	if dir == "" {
		return fmt.Errorf("config dir is empty (DefaultDir returned empty; set ANTHROPIC_CONFIG_DIR)")
	}
	if err := validateProfileName(profile); err != nil {
		return err
	}
	return nil
}

// SaveProfile persists cfg to configs/<profile>.json under dir. The write
// is atomic (.tmp sibling + rename); the target file is mode 0644 and the
// configs/ parent is 0755, matching the spec's "non-secret, checkin-safe"
// positioning for config files (other UIDs on the host — a sidecar, or
// the CLI running under a different user inside a pod — must be able to
// read them).
//
// If cfg.AuthenticationInfo.CredentialsPath exactly matches the default
// resolved path for this profile (i.e. it was populated by LoadConfig
// defaulting the value, not set explicitly by the caller), it is cleared
// on write. Otherwise a load → save round-trip would pin the profile
// config file to the current $HOME via an absolute path, breaking the
// "checkin-safe, relocatable" design goal.
func SaveProfile(dir, profile string, cfg *Config) error {
	if err := validateDirAndProfile(dir, profile); err != nil {
		return fmt.Errorf("SaveProfile: %w", err)
	}
	if cfg == nil {
		return fmt.Errorf("SaveProfile: cfg is nil")
	}
	if cfg.AuthenticationInfo == nil {
		return fmt.Errorf("SaveProfile: cfg.AuthenticationInfo is required")
	}

	toWrite := *cfg
	toWrite.Version = ConfigFileVersion
	authCopy := *cfg.AuthenticationInfo
	toWrite.AuthenticationInfo = &authCopy
	if authCopy.CredentialsPath == ProfileCredentialsPath(dir, profile) {
		authCopy.CredentialsPath = ""
	}

	body, err := json.MarshalIndent(&toWrite, "", "  ")
	if err != nil {
		return fmt.Errorf("SaveProfile: marshal: %w", err)
	}
	target := ProfilePath(dir, profile)
	if err := writeFileAtomic(target, body, publicDirMode, publicFileMode); err != nil {
		return fmt.Errorf("SaveProfile: write %q: %w", target, err)
	}
	return nil
}

// WriteCredentials persists creds to path atomically. The target file is
// mode 0600 and the parent directory 0700. Callers typically build path
// with [ProfileCredentialsPath]; the signature takes a raw path so tests
// and non-profile layouts can target arbitrary locations.
func WriteCredentials(path string, creds Credentials) error {
	if path == "" {
		return fmt.Errorf("WriteCredentials: path is empty")
	}
	body, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return fmt.Errorf("WriteCredentials: marshal: %w", err)
	}
	if err := writeFileAtomic(path, body, secretDirMode, secretFileMode); err != nil {
		return fmt.Errorf("WriteCredentials: write %q: %w", path, err)
	}
	return nil
}

// SetActiveProfile writes the active_config pointer under dir so that
// subsequent [LoadConfig] calls (without ANTHROPIC_PROFILE set) resolve
// to profile. The pointer file sits next to configs/ and is written with
// the "public" config modes.
func SetActiveProfile(dir, profile string) error {
	if err := validateDirAndProfile(dir, profile); err != nil {
		return fmt.Errorf("SetActiveProfile: %w", err)
	}
	target := ActiveConfigPath(dir)
	if err := writeFileAtomic(target, []byte(profile+"\n"), publicDirMode, publicFileMode); err != nil {
		return fmt.Errorf("SetActiveProfile: write %q: %w", target, err)
	}
	return nil
}

// DeleteProfile removes configs/<profile>.json and credentials/<profile>.json
// under dir. If active_config currently points at profile, the pointer file
// is also cleared so the next [LoadConfig] call falls back to "default".
// Missing files are not an error — DeleteProfile is idempotent.
func DeleteProfile(dir, profile string) error {
	if err := validateDirAndProfile(dir, profile); err != nil {
		return fmt.Errorf("DeleteProfile: %w", err)
	}

	for _, p := range []string{
		ProfilePath(dir, profile),
		ProfileCredentialsPath(dir, profile),
	} {
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("DeleteProfile: remove %q: %w", p, err)
		}
	}

	if readActiveConfigPointer(dir) == profile {
		if err := os.Remove(ActiveConfigPath(dir)); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("DeleteProfile: clear active_config: %w", err)
		}
	}
	return nil
}

func readActiveConfigPointer(dir string) string {
	data, err := os.ReadFile(ActiveConfigPath(dir))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
