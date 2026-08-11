package vertex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	sdkoption "github.com/anthropics/anthropic-sdk-go/option"
)

func TestBaseURLForRegion(t *testing.T) {
	testCases := []struct {
		name        string
		region      string
		expectedURL string
	}{
		{
			name:        "global region",
			region:      "global",
			expectedURL: "https://aiplatform.googleapis.com/",
		},
		{
			name:        "us region",
			region:      "us",
			expectedURL: "https://aiplatform.us.rep.googleapis.com/",
		},
		{
			name:        "specific region",
			region:      "us-central1",
			expectedURL: "https://us-central1-aiplatform.googleapis.com/",
		},
		{
			name:        "eu region",
			region:      "eu",
			expectedURL: "https://aiplatform.eu.rep.googleapis.com/",
		},
		{
			name:        "specific european region",
			region:      "europe-west1",
			expectedURL: "https://europe-west1-aiplatform.googleapis.com/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			creds := &google.Credentials{
				TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "fake"}),
			}
			opt := WithCredentials(context.Background(), tc.region, "test-project", creds)

			cfg := &requestconfig.RequestConfig{}
			if err := opt.Apply(cfg); err != nil {
				t.Fatalf("Failed to apply option: %v", err)
			}

			if cfg.BaseURL.String() != tc.expectedURL {
				t.Errorf("Expected base URL %q, got %q", tc.expectedURL, cfg.BaseURL.String())
			}
		})
	}
}

// TestVertexUserMiddlewareObservesAnthropicShape verifies the documented
// ordering: middleware registered before the Vertex option observes the
// Anthropic-shaped request, while the wire receives the rewritten Vertex
// request with OAuth authorization.
func TestVertexUserMiddlewareObservesAnthropicShape(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "")

	var wirePath, wireAuth string
	var wireBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wirePath = r.URL.Path
		wireAuth = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&wireBody); err != nil {
			t.Errorf("Failed to decode wire body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id": "msg_test", "type": "message", "role": "assistant",
			"content": []map[string]any{{"type": "text", "text": "hi"}},
			"model":   "claude-3-sonnet", "stop_reason": "end_turn",
			"usage": map[string]any{"input_tokens": 1, "output_tokens": 1},
		})
	}))
	t.Cleanup(server.Close)

	var observedPath, observedAuth string
	var observedBody map[string]any
	spy := func(r *http.Request, next sdkoption.MiddlewareNext) (*http.Response, error) {
		observedPath = r.URL.Path
		observedAuth = r.Header.Get("Authorization")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(body, &observedBody); err != nil {
			return nil, err
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
		return next(r)
	}

	creds := &google.Credentials{
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "fake"}),
	}
	client := anthropic.NewClient(
		sdkoption.WithoutEnvironmentDefaults(),
		sdkoption.WithMiddleware(spy),
		WithCredentials(context.Background(), "us-central1", "test-project", creds),
		sdkoption.WithBaseURL(server.URL),
	)

	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     "claude-3-sonnet",
		MaxTokens: 1,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("hi")),
		},
	})
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	// The spy (outside the Vertex adaptation) sees the Anthropic shape,
	// before OAuth authorization is attached at the transport.
	if observedPath != "/v1/messages" {
		t.Errorf("Expected middleware to observe path %q, got %q", "/v1/messages", observedPath)
	}
	if observedBody["model"] != "claude-3-sonnet" {
		t.Errorf("Expected middleware to observe model in body, got %v", observedBody["model"])
	}
	if observedAuth != "" {
		t.Errorf("Expected middleware to observe no Authorization header, got %q", observedAuth)
	}

	// The wire sees the rewritten, authorized Vertex shape.
	expectedWirePath := "/v1/projects/test-project/locations/us-central1/publishers/anthropic/models/claude-3-sonnet:rawPredict"
	if wirePath != expectedWirePath {
		t.Errorf("Expected wire path %q, got %q", expectedWirePath, wirePath)
	}
	if _, ok := wireBody["model"]; ok {
		t.Error("Expected model to be removed from the wire body")
	}
	if wireBody["anthropic_version"] != DefaultVersion {
		t.Errorf("Expected anthropic_version %q on the wire, got %v", DefaultVersion, wireBody["anthropic_version"])
	}
	if wireAuth != "Bearer fake" {
		t.Errorf("Expected OAuth Authorization on the wire, got %q", wireAuth)
	}
}

// TestWithGoogleAuthDefaultsCloudPlatformScope drives the workload-identity
// (external_account) flow against a local STS endpoint: without an explicit
// scope from the caller, the token exchange must request cloud-platform —
// external-account credentials fail to mint tokens with no scope at all.
func TestWithGoogleAuthDefaultsCloudPlatformScope(t *testing.T) {
	scopes := make(chan string, 1)
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The STS exchange is the only POST; the probe GET must not report.
		if r.Method != http.MethodPost {
			return
		}
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse token request form: %v", err)
		}
		scopes <- r.Form.Get("scope")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"tok","issued_token_type":"urn:ietf:params:oauth:token-type:access_token","token_type":"Bearer","expires_in":3600}`))
	}))
	t.Cleanup(tokenServer.Close)

	dir := t.TempDir()
	subjectTokenFile := filepath.Join(dir, "subject-token")
	if err := os.WriteFile(subjectTokenFile, []byte("subject-token"), 0600); err != nil {
		t.Fatal(err)
	}
	credFile := filepath.Join(dir, "creds.json")
	credJSON := fmt.Sprintf(`{
		"type": "external_account",
		"audience": "//iam.googleapis.com/projects/1/locations/global/workloadIdentityPools/pool/providers/provider",
		"subject_token_type": "urn:ietf:params:oauth:token-type:jwt",
		"token_url": %q,
		"credential_source": {"file": %q}
	}`, tokenServer.URL, subjectTokenFile)
	if err := os.WriteFile(credFile, []byte(credJSON), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credFile)

	var rc requestconfig.RequestConfig
	if err := rc.Apply(WithGoogleAuth(context.Background(), "us-east5", "proj")); err != nil {
		t.Fatalf("Expected option to apply, got: %v", err)
	}

	// Any request through the option's HTTP client first performs the token
	// exchange; the probe target just needs to answer.
	req, _ := http.NewRequest(http.MethodGet, tokenServer.URL+"/probe", nil)
	resp, err := rc.HTTPClient.Do(req)
	if err != nil {
		t.Fatalf("Expected the probe request to complete, got: %v", err)
	}
	resp.Body.Close()

	select {
	case got := <-scopes:
		if got != cloudPlatformScope {
			t.Errorf("Expected default scope %q in the token exchange, got %q", cloudPlatformScope, got)
		}
	default:
		t.Fatal("Expected a token exchange to hit the local STS endpoint")
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// TestWithCredentialsPreservesUserHTTPClient is the regression test for a
// user-supplied [sdkoption.WithHTTPClient] being replaced by the Vertex
// option: the custom transport must carry the request, with OAuth
// authorization layered on top and the Vertex URL rewrite still applied.
func TestWithCredentialsPreservesUserHTTPClient(t *testing.T) {
	const cannedMessage = `{"id":"msg_test","type":"message","role":"assistant",` +
		`"content":[{"type":"text","text":"hi"}],"model":"claude-sonnet-4-5","stop_reason":"end_turn",` +
		`"usage":{"input_tokens":1,"output_tokens":1}}`

	// Only reachable if the user-supplied transport is bypassed.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("Expected the request to go through the user-supplied transport, but it reached the wire directly")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(cannedMessage))
	}))
	t.Cleanup(server.Close)

	var called bool
	var gotAuth, gotURL string
	recorder := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		called = true
		gotAuth = r.Header.Get("Authorization")
		gotURL = r.URL.String()
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(cannedMessage)),
			Request:    r,
		}, nil
	})
	userClient := &http.Client{Transport: recorder, Timeout: 7 * time.Second}

	creds := &google.Credentials{
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test-token"}),
	}
	client := anthropic.NewClient(
		sdkoption.WithoutEnvironmentDefaults(),
		sdkoption.WithHTTPClient(userClient),
		WithCredentials(context.Background(), "us-east5", "proj", creds),
		sdkoption.WithBaseURL(server.URL),
	)

	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     "claude-sonnet-4-5",
		MaxTokens: 1,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("hi")),
		},
	})
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if !called {
		t.Fatal("Expected the user-supplied transport to carry the request")
	}
	if gotAuth != "Bearer test-token" {
		t.Errorf("Expected OAuth Authorization %q on the wire, got %q", "Bearer test-token", gotAuth)
	}
	expectedURL := server.URL + "/v1/projects/proj/locations/us-east5/publishers/anthropic/models/claude-sonnet-4-5:rawPredict"
	if gotURL != expectedURL {
		t.Errorf("Expected wire URL %q, got %q", expectedURL, gotURL)
	}

	// The applied client keeps the user's settings and the user's client is
	// left untouched.
	rc := requestconfig.RequestConfig{HTTPClient: userClient}
	if err := rc.Apply(WithCredentials(context.Background(), "us-east5", "proj", creds)); err != nil {
		t.Fatalf("Failed to apply option: %v", err)
	}
	if rc.HTTPClient == userClient {
		t.Error("Expected a wrapped copy of the user client, got the same pointer")
	}
	if rc.HTTPClient.Timeout != 7*time.Second {
		t.Errorf("Expected Timeout %v to be preserved, got %v", 7*time.Second, rc.HTTPClient.Timeout)
	}
	if _, ok := userClient.Transport.(roundTripperFunc); !ok {
		t.Errorf("Expected the user client's transport to be left unmodified, got %T", userClient.Transport)
	}
}

// TestWithCredentialsInstallsGoogleClientByDefault verifies that without a
// user-supplied client the option still installs its own authorized client.
func TestWithCredentialsInstallsGoogleClientByDefault(t *testing.T) {
	creds := &google.Credentials{
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test-token"}),
	}
	for _, tc := range []struct {
		name   string
		client *http.Client
	}{
		{name: "nil client", client: nil},
		{name: "http.DefaultClient", client: http.DefaultClient},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rc := requestconfig.RequestConfig{HTTPClient: tc.client}
			if err := rc.Apply(WithCredentials(context.Background(), "us-east5", "proj", creds)); err != nil {
				t.Fatalf("Failed to apply option: %v", err)
			}
			if rc.HTTPClient == nil || rc.HTTPClient == http.DefaultClient {
				t.Errorf("Expected an authorized Google client to be installed, got %v", rc.HTTPClient)
			}
		})
	}
}

// TestVertexClientIgnoresConfigStore is a regression test for the
// first-party credential chain leaking into Vertex clients: a client
// constructed with static Google credentials must never consult the shared
// config store (ANTHROPIC_CONFIG_DIR), so no store-derived value — the
// resolvable profile's workspace id or its access token — may reach the
// wire. Authorization must stay the static OAuth token.
func TestVertexClientIgnoresConfigStore(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("ANTHROPIC_AUTH_TOKEN", "")
	unsetEnv(t, "ANTHROPIC_PROFILE")
	unsetEnv(t, "ANTHROPIC_FEDERATION_RULE_ID")
	unsetEnv(t, "ANTHROPIC_ORGANIZATION_ID")
	unsetEnv(t, "ANTHROPIC_IDENTITY_TOKEN")
	unsetEnv(t, "ANTHROPIC_IDENTITY_TOKEN_FILE")
	unsetEnv(t, "ANTHROPIC_WORKSPACE_ID")
	writeFakeConfigStore(t)

	var wireHeaders http.Header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wireHeaders = r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id": "msg_test", "type": "message", "role": "assistant",
			"content": []map[string]any{{"type": "text", "text": "hi"}},
			"model":   "claude-3-sonnet", "stop_reason": "end_turn",
			"usage": map[string]any{"input_tokens": 1, "output_tokens": 1},
		})
	}))
	t.Cleanup(server.Close)

	creds := &google.Credentials{
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "fake-vertex-access-token"}),
	}
	client := anthropic.NewClient(
		WithCredentials(context.Background(), "us-central1", "test-project", creds),
		sdkoption.WithBaseURL(server.URL),
	)
	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     "claude-3-sonnet",
		MaxTokens: 1,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("hi")),
		},
	})
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if got, ok := wireHeaders[http.CanonicalHeaderKey("anthropic-workspace-id")]; ok {
		t.Errorf("Expected no anthropic-workspace-id header, got %q", got)
	}
	if got := wireHeaders.Get("X-Api-Key"); got != "" {
		t.Errorf("Expected no X-Api-Key header, got %q", got)
	}
	if auth := wireHeaders.Get("Authorization"); auth != "Bearer fake-vertex-access-token" {
		t.Errorf("Expected the static OAuth Authorization on the wire, got %q", auth)
	}
}

// writeFakeConfigStore points ANTHROPIC_CONFIG_DIR at a temp config store
// holding a resolvable default profile with a workspace id, so a client
// that (incorrectly) walks the first-party credential chain would pick it
// up. Values are deliberately fake.
func writeFakeConfigStore(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "configs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "credentials"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "configs", "default.json"), []byte(`{
		"authentication": {"type": "user_oauth"},
		"workspace_id": "wrkspc_fake_store_value"
	}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "credentials", "default.json"),
		[]byte(`{"type":"oauth_token","access_token":"fake-store-access-token"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
}

// unsetEnv unsets an env var for the duration of the test, restoring the
// original value afterwards (t.Setenv registers the restore).
func unsetEnv(t *testing.T, key string) {
	t.Helper()
	t.Setenv(key, "")
	os.Unsetenv(key)
}
