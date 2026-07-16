package googlecloud

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/oauth2"
)

// --- resolveConfig ---

func TestResolveDerivesBaseURLFromLocationAndProject(t *testing.T) {
	clearEnv(t)
	rc, err := resolveConfig(ClientConfig{
		Project:     "my-project",
		Location:    "us-central1",
		WorkspaceID: "wrkspc_123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "https://claude.googleapis.com/v1alpha/projects/my-project/locations/us-central1/workspaces/wrkspc_123/invoke"
	if rc.BaseURL != want {
		t.Errorf("derived base URL = %q, want %q", rc.BaseURL, want)
	}
}

func TestResolveLocationDefaultsToGlobal(t *testing.T) {
	clearEnv(t)
	rc, err := resolveConfig(ClientConfig{
		Project:     "my-project",
		WorkspaceID: "wrkspc_123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rc.Location != "global" {
		t.Errorf("location = %q, want default %q", rc.Location, "global")
	}
	want := "https://claude.googleapis.com/v1alpha/projects/my-project/locations/global/workspaces/wrkspc_123/invoke"
	if rc.BaseURL != want {
		t.Errorf("derived base URL = %q, want %q", rc.BaseURL, want)
	}
}

func TestResolveExplicitBaseURLNeedsNoLocation(t *testing.T) {
	clearEnv(t)
	rc, err := resolveConfig(ClientConfig{
		BaseURL:     "https://proxy.example.com",
		WorkspaceID: "wrkspc_123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rc.BaseURL != "https://proxy.example.com" {
		t.Errorf("base URL = %q, want explicit value", rc.BaseURL)
	}
}

func TestResolveBaseURLFromEnvOverridesDerivation(t *testing.T) {
	clearEnv(t)
	t.Setenv(envBaseURL, "https://env.example.com")
	rc, err := resolveConfig(ClientConfig{
		Project:     "my-project",
		Location:    "us-central1",
		WorkspaceID: "wrkspc_123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rc.BaseURL != "https://env.example.com" {
		t.Errorf("base URL = %q, want env value", rc.BaseURL)
	}
}

func TestResolveLocationFromEnv(t *testing.T) {
	clearEnv(t)
	t.Setenv(envLocation, "europe-west4")
	rc, err := resolveConfig(ClientConfig{
		Project:     "my-project",
		WorkspaceID: "wrkspc_123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rc.Location != "europe-west4" {
		t.Errorf("location = %q, want env value", rc.Location)
	}
	want := "https://claude.googleapis.com/v1alpha/projects/my-project/locations/europe-west4/workspaces/wrkspc_123/invoke"
	if rc.BaseURL != want {
		t.Errorf("derived base URL = %q, want %q", rc.BaseURL, want)
	}
}

func TestResolveProjectFromGoogleCloudProjectEnv(t *testing.T) {
	clearEnv(t)
	t.Setenv("GOOGLE_CLOUD_PROJECT", "gcp-env-project")
	rc, err := resolveConfig(ClientConfig{
		Location:    "us-central1",
		WorkspaceID: "wrkspc_123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rc.Project != "gcp-env-project" {
		t.Errorf("project = %q, want GOOGLE_CLOUD_PROJECT value", rc.Project)
	}
}

func TestResolveAnthropicProjectEnvBeatsGoogleCloudProject(t *testing.T) {
	clearEnv(t)
	t.Setenv(envProject, "anthropic-env-project")
	t.Setenv("GOOGLE_CLOUD_PROJECT", "gcp-env-project")
	rc, err := resolveConfig(ClientConfig{
		Location:    "us-central1",
		WorkspaceID: "wrkspc_123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rc.Project != "anthropic-env-project" {
		t.Errorf("project = %q, want %s value (precedence over GOOGLE_CLOUD_PROJECT)", rc.Project, envProject)
	}
}

func TestResolveExplicitConfigOverridesEnv(t *testing.T) {
	t.Setenv(envProject, "env-project")
	t.Setenv(envLocation, "env-location")
	t.Setenv(envWorkspaceID, "env-workspace")
	t.Setenv(envBaseURL, "")
	t.Setenv("GOOGLE_CLOUD_PROJECT", "gcp-env-project")
	rc, err := resolveConfig(ClientConfig{
		Project:     "explicit-project",
		Location:    "us-central1",
		WorkspaceID: "explicit-workspace",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rc.Project != "explicit-project" {
		t.Errorf("project = %q, want explicit", rc.Project)
	}
	if rc.Location != "us-central1" {
		t.Errorf("location = %q, want explicit", rc.Location)
	}
	if rc.WorkspaceID != "explicit-workspace" {
		t.Errorf("workspace = %q, want explicit", rc.WorkspaceID)
	}
}

func TestResolveWorkspaceFromEnv(t *testing.T) {
	clearEnv(t)
	t.Setenv(envWorkspaceID, "env-workspace")
	rc, err := resolveConfig(ClientConfig{
		BaseURL: "https://proxy.example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rc.WorkspaceID != "env-workspace" {
		t.Errorf("workspace = %q, want env value", rc.WorkspaceID)
	}
}

func TestResolveRequiresWorkspace(t *testing.T) {
	clearEnv(t)
	_, err := resolveConfig(ClientConfig{
		BaseURL: "https://proxy.example.com",
	})
	if err == nil {
		t.Fatal("expected error when WorkspaceID is missing")
	}
}

func TestResolveSkipAuthWithExplicitBaseURLNeedsNoWorkspace(t *testing.T) {
	clearEnv(t)
	_, err := resolveConfig(ClientConfig{
		BaseURL:  "https://proxy.example.com",
		SkipAuth: true,
	})
	if err != nil {
		t.Fatalf("unexpected error with SkipAuth: %v", err)
	}
}

func TestResolveSkipAuthDerivationStillRequiresWorkspace(t *testing.T) {
	clearEnv(t)
	// The workspace ID is part of the derived base URL, so SkipAuth alone no
	// longer waives it — only an explicit BaseURL does.
	_, err := resolveConfig(ClientConfig{
		Project:  "my-project",
		SkipAuth: true,
	})
	if err == nil {
		t.Fatal("expected error when SkipAuth derivation has no WorkspaceID")
	}
	for _, want := range []string{"WorkspaceID", envWorkspaceID} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error = %q, want mention of %q", err, want)
		}
	}
}

func TestResolveSkipAuthDerivesURLWithWorkspace(t *testing.T) {
	clearEnv(t)
	rc, err := resolveConfig(ClientConfig{
		Project:     "my-project",
		WorkspaceID: "wrkspc_123",
		SkipAuth:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "https://claude.googleapis.com/v1alpha/projects/my-project/locations/global/workspaces/wrkspc_123/invoke"
	if rc.BaseURL != want {
		t.Errorf("derived base URL = %q, want %q", rc.BaseURL, want)
	}
}

func TestResolveSkipAuthAndTokenSourceMutuallyExclusive(t *testing.T) {
	clearEnv(t)
	_, err := resolveConfig(ClientConfig{
		BaseURL:     "https://proxy.example.com",
		SkipAuth:    true,
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "tok"}),
	})
	if err == nil {
		t.Fatal("expected error when SkipAuth and TokenSource are both set")
	}
}

func TestResolveEmptyEnvTreatedAsUnset(t *testing.T) {
	// An exported-but-empty env var is the unset case, not "use empty string".
	t.Setenv(envProject, "")
	t.Setenv(envLocation, "")
	t.Setenv(envWorkspaceID, "")
	t.Setenv(envBaseURL, "")
	t.Setenv("GOOGLE_CLOUD_PROJECT", "")

	t.Run("workspace still required", func(t *testing.T) {
		_, err := resolveConfig(ClientConfig{BaseURL: "https://proxy.example.com"})
		if err == nil {
			t.Fatal("expected workspace-required error when env workspace is empty")
		}
	})
	t.Run("project not satisfied by empty env", func(t *testing.T) {
		rc, err := resolveConfig(ClientConfig{WorkspaceID: "wrkspc_123", BaseURL: "https://proxy.example.com"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rc.Project != "" {
			t.Errorf("Project = %q, want empty (env was empty)", rc.Project)
		}
	})
	t.Run("base URL not satisfied by empty env", func(t *testing.T) {
		rc, err := resolveConfig(ClientConfig{Project: "p", Location: "us-central1", WorkspaceID: "wrkspc_123"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rc.BaseURL == "" {
			t.Error("BaseURL = empty, want derived (empty env must not short-circuit derivation)")
		}
	})
}

// --- ADC fallback ---

// adcFixture is a minimal service-account credentials file. The private key is
// never parsed: google.FindDefaultCredentials only reads the JSON and stores
// the bytes on the jwt.Config; key parsing happens at Token() time, which this
// test never reaches. Asserting only construction-time behavior keeps the test
// hermetic — exercising Token() would call the Google token endpoint.
const adcFixture = `{
  "type": "service_account",
  "project_id": "adc-project",
  "private_key_id": "k",
  "private_key": "-----BEGIN PRIVATE KEY-----\nZmFrZQ==\n-----END PRIVATE KEY-----\n",
  "client_email": "test@adc-project.iam.gserviceaccount.com",
  "client_id": "1",
  "token_uri": "https://oauth2.googleapis.com/token"
}`

func TestADCFallbackBackfillsProject(t *testing.T) {
	clearEnv(t)
	path := filepath.Join(t.TempDir(), "adc.json")
	if err := os.WriteFile(path, []byte(adcFixture), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", path)

	// No Project, no BaseURL, no TokenSource: createClientOptions must fall back
	// to ADC, read project_id from the fixture, and derive the base URL from it.
	// A failure to consult ADC would surface as the "no project found" error.
	_, err := createClientOptions(context.Background(), ClientConfig{
		Location:    "us-central1",
		WorkspaceID: "wrkspc_123",
	})
	if err != nil {
		t.Fatalf("createClientOptions with ADC fixture failed: %v", err)
	}
}

func TestADCFallbackInvalidCredentialsSurfacesError(t *testing.T) {
	clearEnv(t)
	path := filepath.Join(t.TempDir(), "adc.json")
	if err := os.WriteFile(path, []byte("not-json"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", path)

	_, err := createClientOptions(context.Background(), ClientConfig{
		Location:    "us-central1",
		WorkspaceID: "wrkspc_123",
	})
	if err == nil {
		t.Fatal("expected ADC load error when GOOGLE_APPLICATION_CREDENTIALS is invalid")
	}
}

// --- bearerMiddleware ---

func TestBearerMiddlewareSetsAuthHeader(t *testing.T) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test-token"})
	mw := bearerMiddleware(ts)

	req, _ := http.NewRequest("POST", "https://example.com/v1/messages", nil)
	var got string
	_, err := mw(req, func(r *http.Request) (*http.Response, error) {
		got = r.Header.Get("Authorization")
		return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
	})
	if err != nil {
		t.Fatalf("middleware error: %v", err)
	}
	if got != "Bearer test-token" {
		t.Errorf("Authorization = %q, want %q", got, "Bearer test-token")
	}
}

func TestBearerMiddlewareDoesNotOverrideExistingAuth(t *testing.T) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test-token"})
	mw := bearerMiddleware(ts)

	req, _ := http.NewRequest("POST", "https://example.com/v1/messages", nil)
	req.Header.Set("Authorization", "Bearer preexisting")
	var got string
	_, err := mw(req, func(r *http.Request) (*http.Response, error) {
		got = r.Header.Get("Authorization")
		return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
	})
	if err != nil {
		t.Fatalf("middleware error: %v", err)
	}
	if got != "Bearer preexisting" {
		t.Errorf("Authorization = %q, want preexisting value untouched", got)
	}
}

func TestBearerMiddlewareDoesNotBufferBody(t *testing.T) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test-token"})
	mw := bearerMiddleware(ts)

	// A request with a nil GetBody must reach next untouched: the bearer middleware
	// has no reason to read or reset the body (unlike SigV4).
	req, _ := http.NewRequest("POST", "https://example.com/v1/messages", http.NoBody)
	req.GetBody = nil
	_, err := mw(req, func(r *http.Request) (*http.Response, error) {
		if r.GetBody != nil {
			t.Error("bearer middleware should not set GetBody")
		}
		return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
	})
	if err != nil {
		t.Fatalf("middleware error: %v", err)
	}
}
