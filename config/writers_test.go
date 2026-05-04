package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/config"
)

func TestDefaultDir_HonorsEnv(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	if got := config.DefaultDir(); got != dir {
		t.Errorf("DefaultDir() = %q, want %q", got, dir)
	}
}

func TestPathHelpers(t *testing.T) {
	dir := "/etc/anthropic"
	want := map[string]string{
		"profile":     filepath.Join(dir, "configs", "work.json"),
		"credentials": filepath.Join(dir, "credentials", "work.json"),
		"active":      filepath.Join(dir, "active_config"),
	}
	if got := config.ProfilePath(dir, "work"); got != want["profile"] {
		t.Errorf("ProfilePath = %q, want %q", got, want["profile"])
	}
	if got := config.ProfileCredentialsPath(dir, "work"); got != want["credentials"] {
		t.Errorf("ProfileCredentialsPath = %q, want %q", got, want["credentials"])
	}
	if got := config.ActiveConfigPath(dir); got != want["active"] {
		t.Errorf("ActiveConfigPath = %q, want %q", got, want["active"])
	}
}

func TestDirHelpers_MatchPerProfileParents(t *testing.T) {
	dir := "/etc/anthropic"
	// ProfilesDir / CredentialsDir must be exactly the parent of the
	// per-profile paths so callers that previously derived the directory via
	// filepath.Dir(ProfilePath(dir, "_")) get the same result.
	if got, want := config.ProfilesDir(dir), filepath.Dir(config.ProfilePath(dir, "x")); got != want {
		t.Errorf("ProfilesDir = %q, want %q", got, want)
	}
	if got, want := config.CredentialsDir(dir), filepath.Dir(config.ProfileCredentialsPath(dir, "x")); got != want {
		t.Errorf("CredentialsDir = %q, want %q", got, want)
	}
}

func TestLoadProfile(t *testing.T) {
	dir := t.TempDir()
	// Point env resolution at a different profile so we prove LoadProfile
	// bypasses ANTHROPIC_PROFILE / active_config.
	t.Setenv("ANTHROPIC_PROFILE", "not-the-one-we-load")

	cfg := &config.Config{
		BaseURL:            "https://api.example.com",
		OrganizationID:     "org_abc",
		AuthenticationInfo: config.NewUserOAuthAuthentication("client-xyz"),
	}
	if err := config.SaveProfile(dir, "work", cfg); err != nil {
		t.Fatal(err)
	}

	loaded, err := config.LoadProfile(dir, "work")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.BaseURL != "https://api.example.com" {
		t.Errorf("BaseURL = %q", loaded.BaseURL)
	}
	if loaded.OrganizationID != "org_abc" {
		t.Errorf("OrganizationID = %q", loaded.OrganizationID)
	}
	if loaded.AuthenticationInfo == nil || loaded.AuthenticationInfo.UserOAuth == nil ||
		loaded.AuthenticationInfo.UserOAuth.ClientID != "client-xyz" {
		t.Errorf("AuthenticationInfo round-trip: %+v", loaded.AuthenticationInfo)
	}
	wantCreds := config.ProfileCredentialsPath(dir, "work")
	if loaded.AuthenticationInfo.CredentialsPath != wantCreds {
		t.Errorf("CredentialsPath = %q, want defaulted %q",
			loaded.AuthenticationInfo.CredentialsPath, wantCreds)
	}

	if _, err := config.LoadProfile(dir, "missing"); err == nil {
		t.Error("expected error for missing profile")
	}
}

func TestListProfiles(t *testing.T) {
	t.Run("nonexistent dir", func(t *testing.T) {
		got, err := config.ListProfiles(filepath.Join(t.TempDir(), "never-created"))
		if err != nil {
			t.Fatalf("err = %v, want nil", err)
		}
		if len(got) != 0 {
			t.Errorf("got %v, want empty", got)
		}
	})

	t.Run("empty configs dir", func(t *testing.T) {
		dir := t.TempDir()
		os.MkdirAll(filepath.Join(dir, "configs"), 0755)
		got, err := config.ListProfiles(dir)
		if err != nil {
			t.Fatalf("err = %v, want nil", err)
		}
		if len(got) != 0 {
			t.Errorf("got %v, want empty", got)
		}
	})

	t.Run("populated and sorted", func(t *testing.T) {
		dir := t.TempDir()
		cfg := &config.Config{AuthenticationInfo: config.NewUserOAuthAuthentication("c")}
		// Save out of order to prove the result is sorted.
		for _, name := range []string{"work", "default", "staging"} {
			if err := config.SaveProfile(dir, name, cfg); err != nil {
				t.Fatal(err)
			}
		}
		// Noise that must be ignored: a .tmp leftover and a stray subdir.
		os.WriteFile(filepath.Join(dir, "configs", "scratch.json.tmp"), []byte("{}"), 0644)
		os.Mkdir(filepath.Join(dir, "configs", "subdir"), 0755)

		got, err := config.ListProfiles(dir)
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"default", "staging", "work"}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestSaveProfile_WritesFileWithMode0644(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		BaseURL:            "https://api.example.com",
		OrganizationID:     "org_123",
		AuthenticationInfo: config.NewUserOAuthAuthentication("client-xyz"),
	}
	if err := config.SaveProfile(dir, "work", cfg); err != nil {
		t.Fatal(err)
	}

	target := config.ProfilePath(dir, "work")
	info, err := os.Stat(target)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" {
		if mode := info.Mode().Perm(); mode != 0644 {
			t.Errorf("config file mode = %o, want 0644 (configs are non-secret)", mode)
		}
	}
	// Parent directory should also be public (0755).
	if runtime.GOOS != "windows" {
		parentInfo, err := os.Stat(filepath.Dir(target))
		if err != nil {
			t.Fatal(err)
		}
		if mode := parentInfo.Mode().Perm(); mode != 0755 {
			t.Errorf("configs/ dir mode = %o, want 0755", mode)
		}
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	var decoded config.Config
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("re-decode: %v", err)
	}
	if decoded.Version != config.ConfigFileVersion {
		t.Errorf("Version = %q, want %q (SaveProfile must stamp it on write)",
			decoded.Version, config.ConfigFileVersion)
	}
	if decoded.BaseURL != "https://api.example.com" {
		t.Errorf("BaseURL round-trip: %q", decoded.BaseURL)
	}
	if decoded.AuthenticationInfo == nil || decoded.AuthenticationInfo.UserOAuth == nil ||
		decoded.AuthenticationInfo.UserOAuth.ClientID != "client-xyz" {
		t.Errorf("UserOAuth round-trip: %+v", decoded.AuthenticationInfo)
	}
	if _, err := os.Stat(target + ".tmp"); !os.IsNotExist(err) {
		t.Errorf("expected no .tmp leftover, stat err = %v", err)
	}
}

func TestSaveProfile_LoadConfigRoundTrip(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "rt")

	original := &config.Config{
		BaseURL:        "https://rt.example.com",
		OrganizationID: "org_rt",
		WorkspaceID:    "wrkspc_rt",
		AuthenticationInfo: config.NewOIDCFederationAuthentication(config.OIDCFederation{
			FederationRuleID: "fdrl_rt",
			ServiceAccountID: "svac_rt",
			IdentityToken: &config.IdentityTokenConfig{
				Source: config.IdentityTokenSourceFile,
				Path:   "/tmp/token",
			},
			Scope: "user:inference",
		}),
	}

	if err := config.SaveProfile(dir, "rt", original); err != nil {
		t.Fatal(err)
	}
	loaded, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.BaseURL != original.BaseURL ||
		loaded.OrganizationID != original.OrganizationID ||
		loaded.WorkspaceID != original.WorkspaceID {
		t.Errorf("top-level fields diverged: got %+v", loaded)
	}
	if loaded.AuthenticationInfo == nil || loaded.AuthenticationInfo.OIDCFederation == nil {
		t.Fatal("OIDCFederation missing after load")
	}
	if loaded.AuthenticationInfo.OIDCFederation.FederationRuleID != "fdrl_rt" {
		t.Errorf("FederationRuleID: %q", loaded.AuthenticationInfo.OIDCFederation.FederationRuleID)
	}
	if loaded.AuthenticationInfo.OIDCFederation.Scope != "user:inference" {
		t.Errorf("Scope round-trip: %q", loaded.AuthenticationInfo.OIDCFederation.Scope)
	}
}

// TestSaveProfile_LoadSaveKeepsCredentialsPathBlank verifies that a
// load → save cycle doesn't rewrite a blank credentials_path into an
// absolute path pinned to the current $HOME. LoadConfig defaults the
// field at read time for the SDK's own consumers; SaveProfile must
// strip that default back out before writing so the on-disk file stays
// relocatable.
func TestSaveProfile_LoadSaveKeepsCredentialsPathBlank(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "portable")

	// Initial save: caller never sets CredentialsPath.
	original := &config.Config{
		OrganizationID:     "org_portable",
		AuthenticationInfo: config.NewUserOAuthAuthentication("client-xyz"),
	}
	if err := config.SaveProfile(dir, "portable", original); err != nil {
		t.Fatal(err)
	}

	// Load populates the field with the default absolute path.
	loaded, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.AuthenticationInfo.CredentialsPath == "" {
		t.Fatal("LoadConfig should populate the default CredentialsPath for runtime use")
	}

	// Save again — the default path must not be persisted back.
	if err := config.SaveProfile(dir, "portable", loaded); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(config.ProfilePath(dir, "portable"))
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	auth, _ := raw["authentication"].(map[string]any)
	if _, present := auth["credentials_path"]; present {
		t.Errorf("credentials_path should be omitted on write, got: %s", string(data))
	}
}

// TestSaveProfile_ExplicitCredentialsPathPreserved asserts that a caller
// who deliberately sets a non-default CredentialsPath still sees it
// round-trip through the save path.
func TestSaveProfile_ExplicitCredentialsPathPreserved(t *testing.T) {
	dir := t.TempDir()
	cfg := &config.Config{
		AuthenticationInfo: &config.AuthenticationInfo{
			Type:            config.AuthenticationTypeUserOAuth,
			CredentialsPath: "/custom/explicit/path.json",
			UserOAuth:       &config.UserOAuth{ClientID: "c"},
		},
	}
	if err := config.SaveProfile(dir, "explicit", cfg); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(config.ProfilePath(dir, "explicit"))
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	auth, _ := raw["authentication"].(map[string]any)
	if auth["credentials_path"] != "/custom/explicit/path.json" {
		t.Errorf("explicit credentials_path lost: %v", auth["credentials_path"])
	}
	// And the caller's original struct must not be mutated by the save.
	if cfg.AuthenticationInfo.CredentialsPath != "/custom/explicit/path.json" {
		t.Errorf("SaveProfile mutated caller's cfg: %q", cfg.AuthenticationInfo.CredentialsPath)
	}
}

func TestSaveProfile_Validation(t *testing.T) {
	dir := t.TempDir()
	good := &config.Config{AuthenticationInfo: config.NewUserOAuthAuthentication("c")}

	cases := []struct {
		name    string
		dir     string
		profile string
		cfg     *config.Config
	}{
		{"empty dir", "", "work", good},
		{"traversal profile", dir, "../evil", good},
		{"nil cfg", dir, "work", nil},
		{"missing auth", dir, "work", &config.Config{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := config.SaveProfile(tc.dir, tc.profile, tc.cfg); err == nil {
				t.Error("expected error")
			}
		})
	}
}

func TestWriteCredentials_WritesFileWithMode0600(t *testing.T) {
	dir := t.TempDir()
	path := config.ProfileCredentialsPath(dir, "myprof")

	expiry := time.Now().Add(time.Hour).Truncate(time.Second)
	creds := config.Credentials{
		AccessToken:  "access-123",
		RefreshToken: "refresh-456",
		ExpiresAt:    &expiry,
	}
	if err := config.WriteCredentials(path, creds); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if runtime.GOOS != "windows" {
		if mode := info.Mode().Perm(); mode != 0600 {
			t.Errorf("credentials file mode = %o, want 0600 (credentials are secret)", mode)
		}
		parentInfo, err := os.Stat(filepath.Dir(path))
		if err != nil {
			t.Fatal(err)
		}
		if mode := parentInfo.Mode().Perm(); mode != 0700 {
			t.Errorf("credentials/ dir mode = %o, want 0700", mode)
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	if raw["version"] != config.CredentialsFileVersion ||
		raw["type"] != "oauth_token" ||
		raw["access_token"] != "access-123" ||
		raw["refresh_token"] != "refresh-456" {
		t.Errorf("unexpected wire contents: %v", raw)
	}
	if raw["expires_at"] == nil {
		t.Error("expires_at missing")
	}
}

func TestWriteCredentials_EmptyPath(t *testing.T) {
	if err := config.WriteCredentials("", config.Credentials{AccessToken: "x"}); err == nil {
		t.Error("expected error for empty path")
	}
}

func TestCredentials_RoundTrip(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cases := []config.Credentials{
		{AccessToken: "a", RefreshToken: "r", ExpiresAt: &now},
		{AccessToken: "a"},
		{AccessToken: "a", RefreshToken: "r"},
		{
			AccessToken:      "a",
			RefreshToken:     "r",
			ExpiresAt:        &now,
			Scope:            "user:inference user:profile",
			OrganizationUUID: "11111111-2222-3333-4444-555555555555",
			OrganizationName: "Acme Inc",
			AccountEmail:     "user@example.com",
		},
	}
	for _, in := range cases {
		b, err := json.Marshal(in)
		if err != nil {
			t.Fatal(err)
		}
		var out config.Credentials
		if err := json.Unmarshal(b, &out); err != nil {
			t.Fatal(err)
		}
		if out.AccessToken != in.AccessToken || out.RefreshToken != in.RefreshToken {
			t.Errorf("token round-trip failed: in=%+v out=%+v", in, out)
		}
		if out.Scope != in.Scope || out.OrganizationUUID != in.OrganizationUUID ||
			out.OrganizationName != in.OrganizationName || out.AccountEmail != in.AccountEmail {
			t.Errorf("token-time field round-trip failed: in=%+v out=%+v", in, out)
		}
		if (in.ExpiresAt == nil) != (out.ExpiresAt == nil) {
			t.Errorf("expiry presence mismatch: in=%+v out=%+v", in.ExpiresAt, out.ExpiresAt)
		}
		if in.ExpiresAt != nil && !out.ExpiresAt.Equal(*in.ExpiresAt) {
			t.Errorf("expiry value mismatch: in=%v out=%v", in.ExpiresAt, out.ExpiresAt)
		}
	}
}

func TestWriteCredentials_RoundTripsTokenTimeFields(t *testing.T) {
	dir := t.TempDir()
	path := config.ProfileCredentialsPath(dir, "p")
	expiry := time.Unix(1_700_000_000, 0)
	in := config.Credentials{
		AccessToken:      "tok",
		RefreshToken:     "ref",
		ExpiresAt:        &expiry,
		Scope:            "user:inference",
		OrganizationUUID: "org-uuid",
		OrganizationName: "Org Name",
		AccountEmail:     "user@example.com",
	}
	if err := config.WriteCredentials(path, in); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var out config.Credentials
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	// Compare via reflect.DeepEqual after normalizing ExpiresAt zone (the
	// wire form is unix seconds, so location is dropped).
	if out.ExpiresAt == nil || !out.ExpiresAt.Equal(*in.ExpiresAt) {
		t.Fatalf("ExpiresAt round-trip: in=%v out=%v", in.ExpiresAt, out.ExpiresAt)
	}
	in.ExpiresAt, out.ExpiresAt = nil, nil
	if !reflect.DeepEqual(in, out) {
		t.Errorf("WriteCredentials round-trip mismatch:\n in=%+v\nout=%+v", in, out)
	}
}

// TestLoadProfile_MissingVersionDecodesAsZero pins the read-side compatibility
// contract: a config file written by an older SDK (before the version field
// existed) loads cleanly with cfg.Version == "". Future SDKs that bump
// ConfigFileVersion can rely on this to distinguish "old file" from "current"
// without an error path.
func TestLoadProfile_MissingVersionDecodesAsZero(t *testing.T) {
	dir := t.TempDir()
	path := config.ProfilePath(dir, "old")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	// Hand-rolled config JSON without the version key, mimicking a file from
	// an older SDK release.
	body := `{"authentication":{"type":"user_oauth","client_id":"cid"}}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "old")
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig on no-version file: %v", err)
	}
	if cfg.Version != "" {
		t.Errorf("cfg.Version = %q, want empty (missing-key default)", cfg.Version)
	}
}

// TestCredentials_UnmarshalMissingVersion pins the equivalent contract for
// the credentials file: a no-version file unmarshals without error. The
// public Credentials struct does not surface Version (it lives only on the
// wire shape), so we just assert the unmarshal succeeds and other fields
// populate.
func TestCredentials_UnmarshalMissingVersion(t *testing.T) {
	body := `{"type":"oauth_token","access_token":"tok","refresh_token":"ref"}`
	var creds config.Credentials
	if err := json.Unmarshal([]byte(body), &creds); err != nil {
		t.Fatalf("Unmarshal no-version credentials: %v", err)
	}
	if creds.AccessToken != "tok" || creds.RefreshToken != "ref" {
		t.Errorf("fields not populated: %+v", creds)
	}
}

func TestSetActiveProfile_WritesPointer(t *testing.T) {
	dir := t.TempDir()
	if err := config.SetActiveProfile(dir, "work"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(config.ActiveConfigPath(dir))
	if err != nil {
		t.Fatal(err)
	}
	if got := string(data); got != "work\n" {
		t.Errorf("active_config contents = %q", got)
	}
}

func TestSetActiveProfile_SwitchesLoadedProfile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "")
	os.Unsetenv("ANTHROPIC_PROFILE")

	for _, name := range []string{"alpha", "beta"} {
		cfg := &config.Config{
			AuthenticationInfo: config.NewUserOAuthAuthentication("cid-" + name),
		}
		if err := config.SaveProfile(dir, name, cfg); err != nil {
			t.Fatal(err)
		}
	}

	if err := config.SetActiveProfile(dir, "beta"); err != nil {
		t.Fatal(err)
	}
	loaded, err := config.LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	if loaded.AuthenticationInfo.UserOAuth.ClientID != "cid-beta" {
		t.Errorf("active profile not honoured: %+v", loaded.AuthenticationInfo)
	}
}

func TestSetActiveProfile_Validation(t *testing.T) {
	dir := t.TempDir()
	if err := config.SetActiveProfile("", "work"); err == nil {
		t.Error("expected error for empty dir")
	}
	if err := config.SetActiveProfile(dir, "../evil"); err == nil {
		t.Error("expected error for traversal profile")
	}
}

func TestDeleteProfile_RemovesBothFilesAndClearsPointer(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "configs"), 0755)
	os.MkdirAll(filepath.Join(dir, "credentials"), 0755)
	os.WriteFile(config.ProfilePath(dir, "work"), []byte("{}"), 0600)
	os.WriteFile(config.ProfileCredentialsPath(dir, "work"), []byte("{}"), 0600)
	os.WriteFile(config.ActiveConfigPath(dir), []byte("work\n"), 0644)

	if err := config.DeleteProfile(dir, "work"); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{
		config.ProfilePath(dir, "work"),
		config.ProfileCredentialsPath(dir, "work"),
		config.ActiveConfigPath(dir),
	} {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Errorf("%s still exists (err %v)", p, err)
		}
	}
}

func TestDeleteProfile_LeavesOtherPointer(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "configs"), 0755)
	os.WriteFile(config.ProfilePath(dir, "work"), []byte("{}"), 0600)
	os.WriteFile(config.ActiveConfigPath(dir), []byte("other\n"), 0644)

	if err := config.DeleteProfile(dir, "work"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(config.ActiveConfigPath(dir))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "other\n" {
		t.Errorf("active_config was clobbered: %q", data)
	}
}

func TestDeleteProfile_Idempotent(t *testing.T) {
	dir := t.TempDir()
	if err := config.DeleteProfile(dir, "nonexistent"); err != nil {
		t.Errorf("first delete returned %v", err)
	}
	if err := config.DeleteProfile(dir, "nonexistent"); err != nil {
		t.Errorf("second delete returned %v", err)
	}
}

// TestWriteCredentials_UniqueTmpFilename verifies that two concurrent
// writes to the same target use different sibling tmp filenames so they
// don't trample each other. The previous fixed "<path>.tmp" was a race
// — a second writer would overwrite the first writer's tmp before its
// rename completed.
func TestWriteCredentials_UniqueTmpFilename(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "creds.json")
	exp := time.Now().Add(time.Hour)
	if err := config.WriteCredentials(target, config.Credentials{AccessToken: "a", ExpiresAt: &exp}); err != nil {
		t.Fatal(err)
	}
	// No fixed-name .tmp leftover from the previous behavior.
	if _, err := os.Stat(target + ".tmp"); !os.IsNotExist(err) {
		t.Errorf("expected no fixed-name .tmp leftover, stat err = %v", err)
	}
	// No leftover .tmp at all (success path cleans them up).
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".tmp" {
			t.Errorf("found .tmp leftover after successful write: %s", e.Name())
		}
	}
}

// TestWriteCredentials_RejectsSymlinkTarget verifies that an attacker who
// pre-plants a symlink at the credentials path cannot redirect the write.
// The read-side check (checkCredentialsFileSafety) is insufficient on its
// own — the write path is the injection vector for an attacker who can
// place a symlink before first-write.
func TestWriteCredentials_RejectsSymlinkTarget(t *testing.T) {
	dir := t.TempDir()
	victim := filepath.Join(dir, "victim.json")
	if err := os.WriteFile(victim, []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "creds.json")
	if err := os.Symlink(victim, target); err != nil {
		t.Skipf("symlink unsupported on this platform: %v", err)
	}
	exp := time.Now().Add(time.Hour)
	err := config.WriteCredentials(target, config.Credentials{AccessToken: "tok", ExpiresAt: &exp})
	if err == nil {
		t.Fatal("expected error when target is a symlink")
	}
	if !strings.Contains(err.Error(), "symlink") {
		t.Errorf("expected symlink error, got: %v", err)
	}
	// Victim must be untouched.
	got, err := os.ReadFile(victim)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "{}" {
		t.Errorf("victim file was modified through symlink: %q", got)
	}
	// No .tmp debris left behind on the failure path.
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".tmp" {
			t.Errorf("found .tmp debris after symlink rejection: %s", e.Name())
		}
	}
}
