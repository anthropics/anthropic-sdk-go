package googlecloud

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/oauth2"

	"github.com/anthropics/anthropic-sdk-go"
)

func clearEnv(t *testing.T) {
	t.Helper()
	t.Setenv(envProject, "")
	t.Setenv(envLocation, "")
	t.Setenv(envWorkspaceID, "")
	t.Setenv(envBaseURL, "")
	t.Setenv("GOOGLE_CLOUD_PROJECT", "")
}

type capturedRequest struct {
	Headers http.Header
	URL     string
}

func messagesHandler(captured *capturedRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		captured.Headers = r.Header.Clone()
		captured.URL = r.URL.String()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":            "msg_test",
			"type":          "message",
			"role":          "assistant",
			"content":       []map[string]any{{"type": "text", "text": "hi"}},
			"model":         "claude-sonnet-4-5",
			"stop_reason":   "end_turn",
			"stop_sequence": nil,
			"usage":         map[string]any{"input_tokens": 1, "output_tokens": 1},
		})
	}
}

// newTestClient stands up an httptest server and a client pointed at it, using a
// static token source so no real Google call is made.
func newTestClient(t *testing.T, cfg ClientConfig) (*Client, *capturedRequest) {
	t.Helper()
	var captured capturedRequest
	server := httptest.NewServer(messagesHandler(&captured))
	t.Cleanup(server.Close)

	cfg.BaseURL = server.URL
	if cfg.TokenSource == nil && !cfg.SkipAuth {
		cfg.TokenSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test-token"})
	}
	client, err := NewClient(context.Background(), cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	return client, &captured
}

func sendTestRequest(t *testing.T, client *Client) {
	t.Helper()
	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     "claude-sonnet-4-5",
		MaxTokens: 1,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("hi")),
		},
	})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

// --- Service surface parity ---

// excludedServices are anthropic.Client fields intentionally NOT mirrored here.
// Completions is the deprecated text-completions API.
var excludedServices = map[string]bool{"Completions": true}

func TestClientServicesMatchAnthropicClient(t *testing.T) {
	gcType := reflect.TypeFor[Client]()
	anthropicType := reflect.TypeFor[anthropic.Client]()

	// Every non-excluded base service must be present here with the same type —
	// catches the base client growing a resource this mirror would silently miss.
	for i := 0; i < anthropicType.NumField(); i++ {
		field := anthropicType.Field(i)
		if excludedServices[field.Name] {
			if _, ok := gcType.FieldByName(field.Name); ok {
				t.Errorf("Client exposes %q, which is intentionally excluded", field.Name)
			}
			continue
		}
		gcField, ok := gcType.FieldByName(field.Name)
		if !ok {
			t.Errorf("Client is missing field %q (type %s) from anthropic.Client", field.Name, field.Type)
			continue
		}
		if gcField.Type != field.Type {
			t.Errorf("Client.%s has type %s, expected %s", field.Name, gcField.Type, field.Type)
		}
	}

	// Conversely, every field here must exist on the base client (no stray fields).
	for i := 0; i < gcType.NumField(); i++ {
		field := gcType.Field(i)
		if _, ok := anthropicType.FieldByName(field.Name); !ok {
			t.Errorf("Client has field %q not present on anthropic.Client", field.Name)
		}
	}
}

// --- Auth header wiring ---

func TestBearerModeHeaders(t *testing.T) {
	clearEnv(t)
	t.Setenv("ANTHROPIC_API_KEY", "")

	client, captured := newTestClient(t, ClientConfig{
		WorkspaceID: "wrkspc_abc",
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "tok-123"}),
	})
	sendTestRequest(t, client)

	if got := captured.Headers.Get("Authorization"); got != "Bearer tok-123" {
		t.Errorf("Authorization = %q, want %q", got, "Bearer tok-123")
	}
	// The workspace id travels in the URL path only — the gateway mints the
	// header itself and rejects client-sent canonical-cased names.
	if got := captured.Headers.Get("Anthropic-Workspace-Id"); got != "" {
		t.Errorf("expected no anthropic-workspace-id header, got %q", got)
	}
	if got := captured.Headers.Get("X-Api-Key"); got != "" {
		t.Errorf("expected no x-api-key, got %q", got)
	}
}

func TestWorkspaceIDFromEnv(t *testing.T) {
	clearEnv(t)
	t.Setenv(envWorkspaceID, "env-workspace")

	client, captured := newTestClient(t, ClientConfig{})
	sendTestRequest(t, client)

	// An env-resolved workspace id feeds URL derivation only — never a header.
	if got := captured.Headers.Get("Anthropic-Workspace-Id"); got != "" {
		t.Errorf("expected no anthropic-workspace-id header, got %q", got)
	}
}

// --- Validation ---

func TestNewClientRequiresWorkspaceID(t *testing.T) {
	clearEnv(t)
	_, err := NewClient(context.Background(), ClientConfig{
		BaseURL:     "https://proxy.example.com",
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "tok"}),
	})
	if err == nil {
		t.Fatal("expected error when WorkspaceID is missing")
	}
}

func TestExplicitBaseURLWithoutLocation(t *testing.T) {
	clearEnv(t)
	_, err := NewClient(context.Background(), ClientConfig{
		BaseURL:     "https://proxy.example.com",
		WorkspaceID: "wrkspc_abc",
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "tok"}),
	})
	if err != nil {
		t.Fatalf("expected no error with explicit BaseURL and no Location, got: %v", err)
	}
}

// --- SkipAuth ---

func TestSkipAuthNoWorkspaceRequired(t *testing.T) {
	clearEnv(t)
	_, err := NewClient(context.Background(), ClientConfig{
		BaseURL:  "https://proxy.example.com",
		SkipAuth: true,
	})
	if err != nil {
		t.Fatalf("expected no error with SkipAuth, got: %v", err)
	}
}

func TestSkipAuthNoAuthHeaders(t *testing.T) {
	clearEnv(t)
	t.Setenv("ANTHROPIC_API_KEY", "")

	client, captured := newTestClient(t, ClientConfig{SkipAuth: true})
	sendTestRequest(t, client)

	if captured.Headers.Get("Authorization") != "" {
		t.Error("expected no Authorization header with SkipAuth")
	}
	if captured.Headers.Get("Anthropic-Workspace-Id") != "" {
		t.Error("expected no anthropic-workspace-id header with SkipAuth")
	}
}

func TestSkipAuthRejectsProvidedCreds(t *testing.T) {
	clearEnv(t)
	_, err := NewClient(context.Background(), ClientConfig{
		BaseURL:     "https://proxy.example.com",
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "tok"}),
		SkipAuth:    true,
	})
	if err == nil {
		t.Fatal("expected error when SkipAuth and TokenSource are both set")
	}
	if !strings.Contains(err.Error(), "mutually exclusive") {
		t.Errorf("error = %q, want mention of mutual exclusivity", err)
	}
}

// --- base-SDK env isolation ---
//
// The googlecloud client constructs services directly (no DefaultClientOptions),
// so the base SDK's env-var credential chain — ANTHROPIC_API_KEY,
// ANTHROPIC_AUTH_TOKEN, ANTHROPIC_BASE_URL, ANTHROPIC_PROFILE/CONFIG_DIR,
// ANTHROPIC_CUSTOM_HEADERS — must contribute nothing to the request.

func TestDoesNotLeakAnthropicAPIKey(t *testing.T) {
	clearEnv(t)
	t.Setenv("ANTHROPIC_API_KEY", "should-not-appear")

	client, captured := newTestClient(t, ClientConfig{
		WorkspaceID: "wrkspc_abc",
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "tok"}),
	})
	sendTestRequest(t, client)

	if got := captured.Headers.Get("X-Api-Key"); got != "" {
		t.Errorf("expected no x-api-key header, got %q", got)
	}
}

func TestDoesNotLeakAnthropicAuthToken(t *testing.T) {
	clearEnv(t)
	t.Setenv("ANTHROPIC_AUTH_TOKEN", "env-auth-token")

	client, captured := newTestClient(t, ClientConfig{
		WorkspaceID: "wrkspc_abc",
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "ts-token"}),
	})
	sendTestRequest(t, client)

	if got := captured.Headers.Values("Authorization"); len(got) != 1 || got[0] != "Bearer ts-token" {
		t.Errorf("Authorization = %q, want exactly [%q] (TokenSource value, not ANTHROPIC_AUTH_TOKEN)", got, "Bearer ts-token")
	}
	if got := captured.Headers.Get("X-Api-Key"); got != "" {
		t.Errorf("expected no x-api-key header, got %q", got)
	}
}

func TestDoesNotLeakAnthropicBaseURL(t *testing.T) {
	clearEnv(t)
	t.Setenv("ANTHROPIC_BASE_URL", "https://should-not-be-used.example.com")

	client, captured := newTestClient(t, ClientConfig{WorkspaceID: "wrkspc_abc"})
	sendTestRequest(t, client)

	// The request reached the gateway test server (captured is populated), not
	// the env URL — newTestClient's BaseURL is the only URL the client sees.
	if captured.URL == "" {
		t.Fatal("request never reached the gateway test server")
	}
	if got := captured.Headers.Get("Authorization"); got != "Bearer test-token" {
		t.Errorf("Authorization = %q, want gateway bearer", got)
	}
}

func TestDoesNotLeakAnthropicProfile(t *testing.T) {
	clearEnv(t)

	// A real-shaped profile under ANTHROPIC_CONFIG_DIR/configs/<profile>.json with a
	// base_url and a credentials file — neither must surface on the wire.
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "configs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "credentials"), 0o700); err != nil {
		t.Fatal(err)
	}
	profile := `{"version":"1.0","base_url":"https://profile.example.com","authentication":{"type":"user_oauth","client_id":"cli_x"}}`
	if err := os.WriteFile(filepath.Join(dir, "configs", "leaky.json"), []byte(profile), 0o644); err != nil {
		t.Fatal(err)
	}
	creds := `{"access_token":"profile-access-token","token_type":"Bearer"}`
	if err := os.WriteFile(filepath.Join(dir, "credentials", "leaky.json"), []byte(creds), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("ANTHROPIC_CONFIG_DIR", dir)
	t.Setenv("ANTHROPIC_PROFILE", "leaky")

	client, captured := newTestClient(t, ClientConfig{
		WorkspaceID: "wrkspc_abc",
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "ts-token"}),
	})
	sendTestRequest(t, client)

	if captured.URL == "" {
		t.Fatal("request never reached the gateway test server")
	}
	if got := captured.Headers.Values("Authorization"); len(got) != 1 || got[0] != "Bearer ts-token" {
		t.Errorf("Authorization = %q, want exactly [%q] (TokenSource value, not profile token)", got, "Bearer ts-token")
	}
	if got := captured.Headers.Get("X-Api-Key"); got != "" {
		t.Errorf("expected no x-api-key header, got %q", got)
	}
}

func TestDoesNotLeakAnthropicCustomHeaders(t *testing.T) {
	clearEnv(t)
	// One benign header and one Authorization override; neither is honored
	// because the googlecloud client never walks DefaultClientOptions. The
	// bearer middleware's only-if-absent rule is therefore not exercised here —
	// the env header simply never reaches the request.
	t.Setenv("ANTHROPIC_CUSTOM_HEADERS", "X-Benign: from-env\nAuthorization: Bearer env-token")

	client, captured := newTestClient(t, ClientConfig{
		WorkspaceID: "wrkspc_abc",
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "ts-token"}),
	})
	sendTestRequest(t, client)

	if got := captured.Headers.Get("X-Benign"); got != "" {
		t.Errorf("X-Benign = %q, want absent (ANTHROPIC_CUSTOM_HEADERS not honored)", got)
	}
	if got := captured.Headers.Values("Authorization"); len(got) != 1 || got[0] != "Bearer ts-token" {
		t.Errorf("Authorization = %q, want exactly [%q]", got, "Bearer ts-token")
	}
}

// --- ADC consultation ---

func TestExplicitTokenSourceSkipsADC(t *testing.T) {
	clearEnv(t)
	// google.FindDefaultCredentials reads GOOGLE_APPLICATION_CREDENTIALS first; a
	// path that exists but is not valid JSON makes it fail. With an explicit
	// TokenSource the client must never reach that call — the request succeeds.
	bad := filepath.Join(t.TempDir(), "broken.json")
	if err := os.WriteFile(bad, []byte("not-json"), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", bad)

	client, captured := newTestClient(t, ClientConfig{
		WorkspaceID: "wrkspc_abc",
		TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "ts-token"}),
	})
	sendTestRequest(t, client)

	if got := captured.Headers.Get("Authorization"); got != "Bearer ts-token" {
		t.Errorf("Authorization = %q, want %q", got, "Bearer ts-token")
	}
}
