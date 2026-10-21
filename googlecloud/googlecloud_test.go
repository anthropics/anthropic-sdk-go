package googlecloud

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"golang.org/x/oauth2"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
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

// excludedServices are anthropic.Client fields intentionally NOT mirrored here:
// the gateway does not serve UserProfiles endpoints.
var excludedServices = map[string]bool{"UserProfiles": true}

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

func TestNewClientNeedsNoWorkspaceID(t *testing.T) {
	clearEnv(t)
	for name, cfg := range map[string]ClientConfig{
		"explicit base URL":              {BaseURL: "https://proxy.example.com"},
		"derived base URL":               {Project: "my-project"},
		"derived base URL with SkipAuth": {Project: "my-project", SkipAuth: true},
	} {
		t.Run(name, func(t *testing.T) {
			if !cfg.SkipAuth {
				cfg.TokenSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "tok"})
			}
			if _, err := NewClient(context.Background(), cfg); err != nil {
				t.Fatalf("expected no error without WorkspaceID, got: %v", err)
			}
		})
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

	if got := captured.Headers.Values("Authorization"); len(got) != 1 || got[0] != "Bearer tok" {
		t.Errorf("Authorization = %q, want exactly [%q] (TokenSource value, not ANTHROPIC_API_KEY)", got, "Bearer tok")
	}
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

// --- Organization-level routes ---

const (
	workspaceBaseURL    = "https://claude.googleapis.com/v1/projects/my-project/locations/us-central1/workspaces/wrkspc_123/invoke"
	organizationBaseURL = "https://claude.googleapis.com/v1/projects/my-project/locations/us-central1/claudeOrganization/invoke"
)

// recordRequests answers every request in place of the network and records
// what was asked for, so a derived gateway URL can be asserted on.
func recordRequests(captured *[]capturedRequest) option.RequestOption {
	return option.WithHTTPClient(&http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			*captured = append(*captured, capturedRequest{Headers: r.Header.Clone(), URL: r.URL.String()})
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"data":[]}`)),
			}, nil
		}),
	})
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func newDerivedTestClient(t *testing.T, cfg ClientConfig) *Client {
	t.Helper()
	clearEnv(t)
	cfg.Project = "my-project"
	cfg.Location = "us-central1"
	if cfg.TokenSource == nil && !cfg.SkipAuth {
		cfg.TokenSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test-token"})
	}
	client, err := NewClient(context.Background(), cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	return client
}

func newMessage(client *Client, opts ...option.RequestOption) error {
	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     "claude-sonnet-4-5",
		MaxTokens: 1,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("hi")),
		},
	}, opts...)
	return err
}

func listWorkspaces(client *Client, opts ...option.RequestOption) error {
	_, err := client.Organization.Workspaces.List(context.Background(), anthropic.OrganizationWorkspaceListParams{}, opts...)
	return err
}

func listBetaWorkspaces(client *Client, opts ...option.RequestOption) error {
	_, err := client.Beta.Organization.Workspaces.List(context.Background(), anthropic.BetaOrganizationWorkspaceListParams{}, opts...)
	return err
}

func requestedURLs(captured []capturedRequest) []string {
	urls := make([]string, len(captured))
	for i, c := range captured {
		urls[i] = c.URL
	}
	return urls
}

func TestDerivedBaseURLByRoute(t *testing.T) {
	tests := []struct {
		name string
		call func(*Client, ...option.RequestOption) error
		want string
	}{
		{"workspace route", newMessage, workspaceBaseURL + "/v1/messages"},
		{"organization route", listWorkspaces, organizationBaseURL + "/v1/organizations/workspaces"},
		{"beta organization route", listBetaWorkspaces, organizationBaseURL + "/v1/organizations/workspaces?beta=true"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newDerivedTestClient(t, ClientConfig{WorkspaceID: "wrkspc_123"})

			var captured []capturedRequest
			// A request option on the call leaves the derived base URL in place.
			if err := tt.call(client, recordRequests(&captured), option.WithMaxRetries(0)); err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if got := requestedURLs(captured); !slices.Equal(got, []string{tt.want}) {
				t.Errorf("requested %q, want [%q]", got, tt.want)
			}
			if got := captured[0].Headers.Get("Authorization"); got != "Bearer test-token" {
				t.Errorf("Authorization = %q, want %q", got, "Bearer test-token")
			}
		})
	}
}

func TestDerivedBaseURLDefaultsToGlobalLocation(t *testing.T) {
	clearEnv(t)
	client, err := NewClient(context.Background(), ClientConfig{
		Project:     "my-project",
		WorkspaceID: "wrkspc_123",
		SkipAuth:    true,
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	var captured []capturedRequest
	if err := newMessage(client, recordRequests(&captured)); err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if err := listWorkspaces(client, recordRequests(&captured)); err != nil {
		t.Fatalf("request failed: %v", err)
	}

	want := []string{
		"https://claude.googleapis.com/v1/projects/my-project/locations/global/workspaces/wrkspc_123/invoke/v1/messages",
		"https://claude.googleapis.com/v1/projects/my-project/locations/global/claudeOrganization/invoke/v1/organizations/workspaces",
	}
	if got := requestedURLs(captured); !slices.Equal(got, want) {
		t.Errorf("requested %q, want %q", got, want)
	}
	if got := captured[0].Headers.Get("Authorization"); got != "" {
		t.Errorf("Authorization = %q, want none with SkipAuth", got)
	}
}

func TestDerivedBaseURLLocationFromEnv(t *testing.T) {
	clearEnv(t)
	t.Setenv(envLocation, "europe-west4")
	client, err := NewClient(context.Background(), ClientConfig{
		Project:     "my-project",
		WorkspaceID: "wrkspc_123",
		SkipAuth:    true,
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	var captured []capturedRequest
	if err := newMessage(client, recordRequests(&captured)); err != nil {
		t.Fatalf("request failed: %v", err)
	}

	want := "https://claude.googleapis.com/v1/projects/my-project/locations/europe-west4/workspaces/wrkspc_123/invoke/v1/messages"
	if got := requestedURLs(captured); !slices.Equal(got, []string{want}) {
		t.Errorf("requested %q, want [%q]", got, want)
	}
}

func TestOrganizationRouteFollowsPagination(t *testing.T) {
	client := newDerivedTestClient(t, ClientConfig{})

	var urls []string
	pages := option.WithHTTPClient(&http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			urls = append(urls, r.URL.String())
			body := `{"data":[{"id":"wrkspc_1"}],"has_more":true,"first_id":"wrkspc_1","last_id":"wrkspc_1"}`
			if len(urls) > 1 {
				body = `{"data":[{"id":"wrkspc_2"}],"has_more":false,"first_id":"wrkspc_2","last_id":"wrkspc_2"}`
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(body)),
			}, nil
		}),
	})

	iter := client.Organization.Workspaces.ListAutoPaging(context.Background(), anthropic.OrganizationWorkspaceListParams{}, pages)
	for iter.Next() {
	}
	if err := iter.Err(); err != nil {
		t.Fatalf("pagination failed: %v", err)
	}

	want := []string{
		organizationBaseURL + "/v1/organizations/workspaces",
		organizationBaseURL + "/v1/organizations/workspaces?after_id=wrkspc_1",
	}
	if !slices.Equal(urls, want) {
		t.Errorf("requested %q, want %q", urls, want)
	}
}

func TestOrganizationRouteNeedsNoWorkspaceID(t *testing.T) {
	for name, cfg := range map[string]ClientConfig{
		"bearer":   {},
		"SkipAuth": {SkipAuth: true},
	} {
		t.Run(name, func(t *testing.T) {
			client := newDerivedTestClient(t, cfg)

			var captured []capturedRequest
			if err := listBetaWorkspaces(client, recordRequests(&captured)); err != nil {
				t.Fatalf("request failed: %v", err)
			}

			want := organizationBaseURL + "/v1/organizations/workspaces?beta=true"
			if got := requestedURLs(captured); !slices.Equal(got, []string{want}) {
				t.Errorf("requested %q, want [%q]", got, want)
			}
		})
	}
}

func TestWorkspaceRouteWithoutWorkspaceIDFailsBeforeSending(t *testing.T) {
	for name, cfg := range map[string]ClientConfig{
		"bearer":   {},
		"SkipAuth": {SkipAuth: true},
	} {
		t.Run(name, func(t *testing.T) {
			client := newDerivedTestClient(t, cfg)

			var captured []capturedRequest
			err := newMessage(client, recordRequests(&captured), option.WithMaxRetries(0))
			if err == nil {
				t.Fatal("expected error when a workspace route has no workspace ID")
			}
			for _, want := range []string{"WorkspaceID", envWorkspaceID} {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error = %q, want mention of %q", err, want)
				}
			}
			if len(captured) != 0 {
				t.Errorf("requested %q, want no request", requestedURLs(captured))
			}
		})
	}
}

func TestWorkspaceRouteWithoutWorkspaceIDFetchesNoToken(t *testing.T) {
	var fetched bool
	client := newDerivedTestClient(t, ClientConfig{
		TokenSource: tokenSourceFunc(func() (*oauth2.Token, error) {
			fetched = true
			return &oauth2.Token{AccessToken: "tok"}, nil
		}),
	})

	var captured []capturedRequest
	if err := newMessage(client, recordRequests(&captured), option.WithMaxRetries(0)); err == nil {
		t.Fatal("expected error when a workspace route has no workspace ID")
	}
	if fetched {
		t.Error("token fetched for a request that is never sent")
	}
}

type tokenSourceFunc func() (*oauth2.Token, error)

func (f tokenSourceFunc) Token() (*oauth2.Token, error) { return f() }

func TestExplicitBaseURLTakesEveryRoute(t *testing.T) {
	tests := []struct {
		name string
		cfg  ClientConfig
		env  string
		opts []option.RequestOption
	}{
		{name: "ClientConfig.BaseURL", cfg: ClientConfig{BaseURL: "https://proxy.example.com", WorkspaceID: "wrkspc_123"}},
		{name: "ClientConfig.BaseURL without workspace ID", cfg: ClientConfig{BaseURL: "https://proxy.example.com"}},
		{name: "environment variable", cfg: ClientConfig{Project: "my-project", WorkspaceID: "wrkspc_123"}, env: "https://proxy.example.com"},
		{name: "environment variable without workspace ID", cfg: ClientConfig{Project: "my-project"}, env: "https://proxy.example.com"},
		{
			name: "request option",
			cfg:  ClientConfig{Project: "my-project", WorkspaceID: "wrkspc_123"},
			opts: []option.RequestOption{option.WithBaseURL("https://proxy.example.com")},
		},
		{
			name: "request option without workspace ID",
			cfg:  ClientConfig{Project: "my-project"},
			opts: []option.RequestOption{option.WithBaseURL("https://proxy.example.com")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t)
			t.Setenv(envBaseURL, tt.env)
			tt.cfg.TokenSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test-token"})
			client, err := NewClient(context.Background(), tt.cfg)
			if err != nil {
				t.Fatalf("NewClient failed: %v", err)
			}

			var captured []capturedRequest
			opts := append(tt.opts, recordRequests(&captured), option.WithMaxRetries(0))
			for _, call := range []func(*Client, ...option.RequestOption) error{newMessage, listWorkspaces, listBetaWorkspaces} {
				if err := call(client, opts...); err != nil {
					t.Fatalf("request failed: %v", err)
				}
			}

			want := []string{
				"https://proxy.example.com/v1/messages",
				"https://proxy.example.com/v1/organizations/workspaces",
				"https://proxy.example.com/v1/organizations/workspaces?beta=true",
			}
			if got := requestedURLs(captured); !slices.Equal(got, want) {
				t.Errorf("requested %q, want %q", got, want)
			}
			for _, c := range captured {
				if got := c.Headers.Get("Authorization"); got != "Bearer test-token" {
					t.Errorf("Authorization = %q for %s, want %q", got, c.URL, "Bearer test-token")
				}
			}
		})
	}
}

func TestADCProjectReachesOrganizationRoute(t *testing.T) {
	clearEnv(t)
	path := filepath.Join(t.TempDir(), "adc.json")
	if err := os.WriteFile(path, []byte(adcFixture), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", path)

	client, err := NewClient(context.Background(), ClientConfig{Location: "us-central1"})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	var captured []capturedRequest
	// Authorization is preset so the bearer middleware never asks the fixture
	// for a token, which would call the Google token endpoint.
	err = listBetaWorkspaces(client, recordRequests(&captured), option.WithHeader("Authorization", "Bearer preset"))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	want := "https://claude.googleapis.com/v1/projects/adc-project/locations/us-central1/claudeOrganization/invoke/v1/organizations/workspaces?beta=true"
	if got := requestedURLs(captured); !slices.Equal(got, []string{want}) {
		t.Errorf("requested %q, want [%q]", got, want)
	}
}
