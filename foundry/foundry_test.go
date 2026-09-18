package foundry

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func clearEnv(t *testing.T) {
	t.Helper()
	t.Setenv(envAPIKey, "")
	t.Setenv(envResource, "")
	t.Setenv(envBaseURL, "")
	// The developer's own first-party credentials must not reach the tests.
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("ANTHROPIC_AUTH_TOKEN", "")
	t.Setenv("ANTHROPIC_BASE_URL", "")
	t.Setenv("ANTHROPIC_PROFILE", "")
}

type recordedRequest struct {
	Method  string
	Path    string
	Headers http.Header
}

// newServer records every request and answers each with a minimal message.
func newServer(t *testing.T) (*httptest.Server, *[]recordedRequest) {
	t.Helper()
	var requests []recordedRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, recordedRequest{Method: r.Method, Path: r.URL.Path, Headers: r.Header.Clone()})
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":            "msg_test",
			"type":          "message",
			"role":          "assistant",
			"content":       []map[string]any{{"type": "text", "text": "hi"}},
			"model":         "claude-sonnet-4-5",
			"stop_reason":   "end_turn",
			"stop_sequence": nil,
			"usage":         map[string]any{"input_tokens": 1, "output_tokens": 1},
		})
	}))
	t.Cleanup(server.Close)
	return server, &requests
}

func newTestClient(t *testing.T, server *httptest.Server, cfg ClientConfig, opts ...option.RequestOption) *Client {
	t.Helper()
	// Mounted under a path prefix, like the real service, to catch a client that drops it.
	if cfg.BaseURL == "" && cfg.Resource == "" {
		cfg.BaseURL = server.URL + "/anthropic"
	}
	client, err := NewClient(cfg, opts...)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}

func sendMessage(t *testing.T, client *Client) {
	t.Helper()
	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     "claude-sonnet-4-5",
		MaxTokens: 1,
		Messages:  []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hi"))},
	})
	if err != nil {
		t.Fatalf("Messages.New: %v", err)
	}
}

func onlyRequest(t *testing.T, requests *[]recordedRequest) recordedRequest {
	t.Helper()
	if len(*requests) != 1 {
		t.Fatalf("expected exactly one request, got %d", len(*requests))
	}
	return (*requests)[0]
}

func staticToken(token string) func(context.Context) (string, error) {
	return func(context.Context) (string, error) { return token, nil }
}

// --- Wire shape ---

func TestAPIKeyRequestIsFirstPartyShaped(t *testing.T) {
	clearEnv(t)
	server, requests := newServer(t)
	client := newTestClient(t, server, ClientConfig{APIKey: "fake-foundry-api-key"})
	sendMessage(t, client)

	req := onlyRequest(t, requests)
	if req.Method != http.MethodPost || req.Path != "/anthropic/v1/messages" {
		t.Errorf("request = %s %s, want POST /anthropic/v1/messages", req.Method, req.Path)
	}
	if got := req.Headers.Values("X-Api-Key"); len(got) != 1 || got[0] != "fake-foundry-api-key" {
		t.Errorf("x-api-key = %q, want exactly [fake-foundry-api-key]", got)
	}
	if got := req.Headers.Get("Anthropic-Version"); got == "" {
		t.Error("expected the anthropic-version header")
	}
	if got := req.Headers.Get("Authorization"); got != "" {
		t.Errorf("expected no Authorization header in API-key mode, got %q", got)
	}
	if got := req.Headers.Get("Api-Key"); got != "" {
		t.Errorf("expected no api-key header, got %q", got)
	}
}

func TestTokenProviderSendsBearerOnly(t *testing.T) {
	clearEnv(t)
	server, requests := newServer(t)
	client := newTestClient(t, server, ClientConfig{AzureADTokenProvider: staticToken("entra-token")})
	sendMessage(t, client)

	req := onlyRequest(t, requests)
	if got := req.Headers.Values("Authorization"); len(got) != 1 || got[0] != "Bearer entra-token" {
		t.Errorf("Authorization = %q, want exactly [Bearer entra-token]", got)
	}
	if got := req.Headers.Get("X-Api-Key"); got != "" {
		t.Errorf("expected no x-api-key header in token mode, got %q", got)
	}
}

func TestTokenProviderIsCalledPerAttemptWithRequestContext(t *testing.T) {
	clearEnv(t)
	var calls atomic.Int32
	var attempts atomic.Int32
	var seen []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Header.Get("Authorization"))
		if attempts.Add(1) == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"msg_test","type":"message","role":"assistant","content":[],"model":"m","usage":{"input_tokens":1,"output_tokens":1}}`))
	}))
	t.Cleanup(server.Close)

	type ctxKey struct{}
	client := newTestClient(t, server, ClientConfig{
		AzureADTokenProvider: func(ctx context.Context) (string, error) {
			if ctx.Value(ctxKey{}) != "marker" {
				t.Error("provider did not receive the request context")
			}
			return "token-" + string(rune('a'+calls.Add(1)-1)), nil
		},
	})
	ctx := context.WithValue(context.Background(), ctxKey{}, "marker")
	_, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     "m",
		MaxTokens: 1,
		Messages:  []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hi"))},
	}, option.WithMaxRetries(1))
	if err != nil {
		t.Fatalf("Messages.New: %v", err)
	}
	// A retry must ask the provider again, since that is where refresh happens.
	if want := []string{"Bearer token-a", "Bearer token-b"}; !reflect.DeepEqual(seen, want) {
		t.Errorf("Authorization per attempt = %q, want %q", seen, want)
	}
}

func TestExplicitAuthorizationHeaderWinsOverProvider(t *testing.T) {
	clearEnv(t)
	server, requests := newServer(t)
	provider := func(context.Context) (string, error) {
		t.Error("provider must not be called when Authorization is already set")
		return "", nil
	}
	client := newTestClient(t, server, ClientConfig{AzureADTokenProvider: provider},
		option.WithHeader("Authorization", "Bearer caller-token"))
	sendMessage(t, client)

	if got := onlyRequest(t, requests).Headers.Values("Authorization"); len(got) != 1 || got[0] != "Bearer caller-token" {
		t.Errorf("Authorization = %q, want exactly [Bearer caller-token]", got)
	}
}

func TestTokenProviderErrors(t *testing.T) {
	clearEnv(t)
	providerErr := errors.New("credential unavailable")
	tests := []struct {
		name     string
		provider func(context.Context) (string, error)
		want     string
		wantIs   error
	}{
		{
			name:     "provider error is wrapped and preserved for errors.Is",
			provider: func(context.Context) (string, error) { return "", providerErr },
			want:     "AzureADTokenProvider failed",
			wantIs:   providerErr,
		},
		{
			name:     "empty token is rejected instead of sending 'Bearer '",
			provider: staticToken(""),
			want:     "empty token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, requests := newServer(t)
			// Without this the SDK would retry the failing provider with backoff.
			client := newTestClient(t, server, ClientConfig{AzureADTokenProvider: tt.provider}, option.WithMaxRetries(0))
			_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
				Model:     "m",
				MaxTokens: 1,
				Messages:  []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hi"))},
			})
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want one containing %q", err, tt.want)
			}
			if tt.wantIs != nil && !errors.Is(err, tt.wantIs) {
				t.Errorf("errors.Is(err, providerErr) = false; err = %v", err)
			}
			if len(*requests) != 0 {
				t.Errorf("expected no request to reach the server, got %d", len(*requests))
			}
		})
	}
}

// --- Config resolution ---

func TestResolveConfig(t *testing.T) {
	provider := staticToken("tok")
	tests := []struct {
		name    string
		env     map[string]string
		cfg     ClientConfig
		want    ClientConfig // compared field by field, since funcs are not comparable
		wantErr string
	}{
		{
			name: "resource derives the services.ai.azure.com base URL",
			cfg:  ClientConfig{APIKey: "k", Resource: "my-resource"},
			want: ClientConfig{APIKey: "k", Resource: "my-resource", BaseURL: "https://my-resource.services.ai.azure.com/anthropic/"},
		},
		{
			name: "everything resolves from the environment",
			env:  map[string]string{envAPIKey: "env-key", envResource: "env-resource"},
			want: ClientConfig{APIKey: "env-key", Resource: "env-resource", BaseURL: "https://env-resource.services.ai.azure.com/anthropic/"},
		},
		{
			name: "explicit fields beat the environment",
			env:  map[string]string{envAPIKey: "env-key", envResource: "env-resource"},
			cfg:  ClientConfig{APIKey: "k", Resource: "r"},
			want: ClientConfig{APIKey: "k", Resource: "r", BaseURL: "https://r.services.ai.azure.com/anthropic/"},
		},
		{
			name: "base URL from the environment is used verbatim",
			env:  map[string]string{envAPIKey: "env-key", envBaseURL: "http://127.0.0.1:1/mock"},
			want: ClientConfig{APIKey: "env-key", BaseURL: "http://127.0.0.1:1/mock"},
		},
		{
			name:    "no credential",
			cfg:     ClientConfig{Resource: "r"},
			wantErr: envAPIKey,
		},
		{
			name:    "explicit key and token provider are mutually exclusive",
			cfg:     ClientConfig{APIKey: "k", AzureADTokenProvider: provider, Resource: "r"},
			wantErr: "mutually exclusive",
		},
		{
			name: "token provider takes precedence over an environment key",
			env:  map[string]string{envAPIKey: "env-key"},
			cfg:  ClientConfig{AzureADTokenProvider: provider, Resource: "r"},
			want: ClientConfig{Resource: "r", BaseURL: "https://r.services.ai.azure.com/anthropic/"},
		},
		{
			name:    "no endpoint",
			cfg:     ClientConfig{APIKey: "k"},
			wantErr: envResource,
		},
		{
			name:    "resource and base URL are mutually exclusive",
			cfg:     ClientConfig{APIKey: "k", Resource: "r", BaseURL: "https://example.test/anthropic"},
			wantErr: "mutually exclusive",
		},
		{
			name:    "environment resource alongside an explicit base URL is still an error",
			env:     map[string]string{envResource: "env-resource"},
			cfg:     ClientConfig{APIKey: "k", BaseURL: "https://example.test/anthropic"},
			wantErr: "mutually exclusive",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv(t)
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			got, err := resolveConfig(tt.cfg)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %v, want one containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.APIKey != tt.want.APIKey || got.Resource != tt.want.Resource || got.BaseURL != tt.want.BaseURL {
				t.Errorf("resolved = {APIKey:%q Resource:%q BaseURL:%q}, want {APIKey:%q Resource:%q BaseURL:%q}",
					got.APIKey, got.Resource, got.BaseURL, tt.want.APIKey, tt.want.Resource, tt.want.BaseURL)
			}
			if (got.AzureADTokenProvider != nil) != (tt.cfg.AzureADTokenProvider != nil) {
				t.Errorf("token provider presence changed during resolution")
			}
		})
	}
}

func TestDerivedBaseURLMountsV1UnderAnthropicPrefix(t *testing.T) {
	clearEnv(t)
	var got string
	transport := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		got = r.URL.String()
		return nil, errors.New("stop")
	})
	client, err := NewClient(ClientConfig{APIKey: "k", Resource: "my-resource"},
		option.WithHTTPClient(&http.Client{Transport: transport}), option.WithMaxRetries(0))
	if err != nil {
		t.Fatal(err)
	}
	_, _ = client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model: "m", MaxTokens: 1,
		Messages: []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hi"))},
	})
	if want := "https://my-resource.services.ai.azure.com/anthropic/v1/messages"; got != want {
		t.Errorf("request URL = %q, want %q", got, want)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// --- First-party isolation ---

func TestFirstPartyEnvironmentDoesNotLeak(t *testing.T) {
	clearEnv(t)
	t.Setenv("ANTHROPIC_API_KEY", "leaked-1p-key")
	t.Setenv("ANTHROPIC_AUTH_TOKEN", "leaked-1p-token")
	t.Setenv("ANTHROPIC_BASE_URL", "http://127.0.0.1:1/should-not-be-used")

	server, requests := newServer(t)
	client := newTestClient(t, server, ClientConfig{APIKey: "fake-foundry-api-key"})
	sendMessage(t, client)

	req := onlyRequest(t, requests)
	if got := req.Headers.Values("X-Api-Key"); len(got) != 1 || got[0] != "fake-foundry-api-key" {
		t.Errorf("x-api-key = %q, want exactly [fake-foundry-api-key]", got)
	}
	if got := req.Headers.Get("Authorization"); got != "" {
		t.Errorf("ANTHROPIC_AUTH_TOKEN leaked as Authorization %q", got)
	}
}

// --- Service surface ---

// Services Foundry does not serve. A service added to anthropic.Client fails
// the test below until someone decides whether Foundry serves it.
var excludedServices = map[string]bool{"Completions": true, "Models": true}

func TestClientServicesMatchAnthropicClient(t *testing.T) {
	foundryType := reflect.TypeFor[Client]()
	anthropicType := reflect.TypeFor[anthropic.Client]()

	for i := 0; i < anthropicType.NumField(); i++ {
		field := anthropicType.Field(i)
		foundryField, ok := foundryType.FieldByName(field.Name)
		if excludedServices[field.Name] {
			if ok {
				t.Errorf("Client exposes %q, which Foundry does not serve", field.Name)
			}
			continue
		}
		if !ok {
			t.Errorf("Client is missing %q (%s) from anthropic.Client", field.Name, field.Type)
			continue
		}
		if foundryField.Type != field.Type {
			t.Errorf("Client.%s has type %s, want %s", field.Name, foundryField.Type, field.Type)
		}
	}
	for i := 0; i < foundryType.NumField(); i++ {
		name := foundryType.Field(i).Name
		if _, ok := anthropicType.FieldByName(name); !ok {
			t.Errorf("Client has field %q not present on anthropic.Client", name)
		}
	}
}

// --- Options composition ---

func TestOptionsBuildAFirstPartyClient(t *testing.T) {
	clearEnv(t)
	// Only the auth token is set: the default chain stops at an API key, which
	// the Foundry key would overwrite, hiding the leak.
	t.Setenv("ANTHROPIC_AUTH_TOKEN", "leaked-1p-token")
	server, requests := newServer(t)
	fc := newTestClient(t, server, ClientConfig{APIKey: "fake-foundry-api-key"})

	client := anthropic.NewClient(fc.Options...)
	if _, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model: "m", MaxTokens: 1,
		Messages: []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hi"))},
	}); err != nil {
		t.Fatal(err)
	}
	req := onlyRequest(t, requests)
	if req.Path != "/anthropic/v1/messages" {
		t.Errorf("path = %q", req.Path)
	}
	if got := req.Headers.Values("X-Api-Key"); len(got) != 1 || got[0] != "fake-foundry-api-key" {
		t.Errorf("x-api-key = %q, want exactly [fake-foundry-api-key]", got)
	}
	if got := req.Headers.Get("Authorization"); got != "" {
		t.Errorf("ANTHROPIC_AUTH_TOKEN leaked as Authorization %q", got)
	}
}

func TestCallerOptionsCanOverrideHeaders(t *testing.T) {
	clearEnv(t)
	server, requests := newServer(t)
	client := newTestClient(t, server, ClientConfig{APIKey: "k"}, option.WithHeader("x-custom", "1"))
	sendMessage(t, client)
	if got := onlyRequest(t, requests).Headers.Get("X-Custom"); got != "1" {
		t.Errorf("x-custom = %q, want 1", got)
	}
}
