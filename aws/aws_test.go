package aws

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/awsauth"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

func makeStaticConfig(region string) awssdk.Config {
	return awssdk.Config{
		Region: region,
		Credentials: credentials.StaticCredentialsProvider{
			Value: awssdk.Credentials{
				AccessKeyID:     "test-access-key",
				SecretAccessKey: "test-secret-key",
			},
		},
	}
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
			"model":         "claude-sonnet-4-6-20250514",
			"stop_reason":   "end_turn",
			"stop_sequence": nil,
			"usage":         map[string]any{"input_tokens": 1, "output_tokens": 1},
		})
	}
}

func newTestClient(t *testing.T, cfg ClientConfig) (*Client, *capturedRequest) {
	t.Helper()
	var captured capturedRequest
	server := httptest.NewServer(messagesHandler(&captured))
	t.Cleanup(server.Close)

	cfg.BaseURL = server.URL
	client, err := NewClient(context.Background(), cfg)
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	return client, &captured
}

func sendTestRequest(t *testing.T, client *Client) {
	t.Helper()
	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     "claude-sonnet-4-6-20250514",
		MaxTokens: 1,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("hi")),
		},
	})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

// --- Service sync test ---

// excludedServices are anthropic.Client fields intentionally not mirrored here.
// Organization is the Admin API, which takes a first-party admin key.
var excludedServices = map[string]bool{"Organization": true}

func TestClientServicesMatchAnthropicClient(t *testing.T) {
	awsType := reflect.TypeOf(Client{})
	anthropicType := reflect.TypeOf(anthropic.Client{})

	for i := 0; i < anthropicType.NumField(); i++ {
		field := anthropicType.Field(i)
		awsField, ok := awsType.FieldByName(field.Name)
		if excludedServices[field.Name] {
			if ok {
				t.Errorf("aws.Client exposes %q, which is intentionally excluded", field.Name)
			}
			continue
		}
		if !ok {
			t.Errorf("aws.Client is missing field %q (type %s) from anthropic.Client", field.Name, field.Type)
			continue
		}
		if awsField.Type != field.Type {
			t.Errorf("aws.Client.%s has type %s, expected %s", field.Name, awsField.Type, field.Type)
		}
	}
}

// --- Validation tests ---

func TestNewClientRequiresWorkspaceID(t *testing.T) {
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	_, err := NewClient(context.Background(), ClientConfig{
		AWSRegion:          "us-east-1",
		AWSAccessKey:       "key",
		AWSSecretAccessKey: "secret",
	})
	if err == nil {
		t.Fatal("expected error when WorkspaceID is missing")
	}
}

func TestNewClientRequiresBaseURLOrRegion(t *testing.T) {
	t.Setenv("AWS_REGION", "")
	t.Setenv("ANTHROPIC_AWS_BASE_URL", "")
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "default")

	_, err := NewClient(context.Background(), ClientConfig{
		APIKey: "my-key",
	})
	if err == nil {
		t.Fatal("expected error when neither base URL nor region is available")
	}
}

func TestNewClientRegionRequiredForSigV4(t *testing.T) {
	t.Setenv("AWS_REGION", "")
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "default")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	_, err := NewClient(context.Background(), ClientConfig{
		AWSAccessKey:       "key",
		AWSSecretAccessKey: "secret",
	})
	if err == nil {
		t.Fatal("expected error when region is missing and SigV4 is used")
	}
}

// --- API key mode tests ---

func TestAPIKeyModeHeaders(t *testing.T) {
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	client, captured := newTestClient(t, ClientConfig{
		APIKey:      "my-api-key",
		AWSRegion:   "us-east-1",
		WorkspaceID: "ws-123",
	})
	sendTestRequest(t, client)

	if got := captured.Headers.Get("Authorization"); got != "Bearer my-api-key" {
		t.Errorf("expected Authorization %q, got %q", "Bearer my-api-key", got)
	}
	if got := captured.Headers.Get("Anthropic-Workspace-Id"); got != "ws-123" {
		t.Errorf("expected anthropic-workspace-id %q, got %q", "ws-123", got)
	}
	if got := captured.Headers.Get("X-Api-Key"); got != "" {
		t.Errorf("expected no x-api-key header, got %q", got)
	}
}

func TestAPIKeyFromEnvHeaders(t *testing.T) {
	t.Setenv("ANTHROPIC_AWS_API_KEY", "env-api-key")
	t.Setenv("AWS_REGION", "us-east-1")
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "")

	client, captured := newTestClient(t, ClientConfig{
		WorkspaceID: "ws-456",
	})
	sendTestRequest(t, client)

	if got := captured.Headers.Get("Authorization"); got != "Bearer env-api-key" {
		t.Errorf("expected Authorization %q, got %q", "Bearer env-api-key", got)
	}
}

func TestExplicitAPIKeyOverridesEnv(t *testing.T) {
	t.Setenv("ANTHROPIC_AWS_API_KEY", "env-api-key")
	t.Setenv("AWS_REGION", "us-east-1")

	client, captured := newTestClient(t, ClientConfig{
		APIKey:      "explicit-key",
		WorkspaceID: "ws-789",
	})
	sendTestRequest(t, client)

	if got := captured.Headers.Get("Authorization"); got != "Bearer explicit-key" {
		t.Errorf("expected Authorization %q, got %q", "Bearer explicit-key", got)
	}
}

// --- SigV4 mode tests ---

func TestSigV4ModeHeaders(t *testing.T) {
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "")

	client, captured := newTestClient(t, ClientConfig{
		AWSRegion:          "us-east-1",
		AWSAccessKey:       "test-access-key",
		AWSSecretAccessKey: "test-secret-key",
		WorkspaceID:        "ws-sigv4",
	})
	sendTestRequest(t, client)

	auth := captured.Headers.Get("Authorization")
	if auth == "" {
		t.Fatal("expected Authorization header in SigV4 mode")
	}
	if !strings.HasPrefix(auth, "AWS4-HMAC-SHA256") {
		t.Errorf("expected AWS4 signature, got: %s", auth)
	}
	if !strings.Contains(auth, defaultServiceName) {
		t.Errorf("expected service name %q in Authorization, got: %s", defaultServiceName, auth)
	}
	if captured.Headers.Get("X-Amz-Date") == "" {
		t.Error("expected X-Amz-Date header in SigV4 mode")
	}
	if got := captured.Headers.Get("Anthropic-Workspace-Id"); got != "ws-sigv4" {
		t.Errorf("expected anthropic-workspace-id %q, got %q", "ws-sigv4", got)
	}
}

func TestSigV4WithSessionToken(t *testing.T) {
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "")

	client, captured := newTestClient(t, ClientConfig{
		AWSRegion:          "us-east-1",
		AWSAccessKey:       "test-access-key",
		AWSSecretAccessKey: "test-secret-key",
		AWSSessionToken:    "test-session-token",
		WorkspaceID:        "ws-session",
	})
	sendTestRequest(t, client)

	if captured.Headers.Get("X-Amz-Security-Token") != "test-session-token" {
		t.Errorf("expected X-Amz-Security-Token %q, got %q", "test-session-token", captured.Headers.Get("X-Amz-Security-Token"))
	}
}

// --- Workspace ID tests ---

func TestWorkspaceIDFromEnv(t *testing.T) {
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "env-workspace")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "test-key")
	t.Setenv("AWS_REGION", "us-east-1")

	client, captured := newTestClient(t, ClientConfig{})
	sendTestRequest(t, client)

	if got := captured.Headers.Get("Anthropic-Workspace-Id"); got != "env-workspace" {
		t.Errorf("expected anthropic-workspace-id %q, got %q", "env-workspace", got)
	}
}

func TestPerRequestWorkspaceIDOverridesClientConfig(t *testing.T) {
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	client, captured := newTestClient(t, ClientConfig{
		APIKey:      "my-api-key",
		AWSRegion:   "us-east-1",
		WorkspaceID: "ws-client",
	})
	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:       "claude-sonnet-4-6-20250514",
		MaxTokens:   1,
		Messages:    []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hi"))},
		WorkspaceID: anthropic.String("ws-request"),
	})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if got := captured.Headers.Values("Anthropic-Workspace-Id"); len(got) != 1 || got[0] != "ws-request" {
		t.Errorf("expected single anthropic-workspace-id %q, got %q", "ws-request", got)
	}
}

// --- Base URL tests ---

func TestBaseURLDerivedFromRegion(t *testing.T) {
	t.Setenv("ANTHROPIC_AWS_BASE_URL", "")
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	resolved, err := awsauth.ResolveConfig(toInternalConfig(ClientConfig{
		AWSRegion:   "us-west-2",
		WorkspaceID: "default",
		APIKey:      "my-key",
	}), awsResolveParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "https://aws-external-anthropic.us-west-2.api.aws"
	if resolved.BaseURL != expected {
		t.Errorf("expected base URL %q, got %q", expected, resolved.BaseURL)
	}
}

func TestBaseURLFromEnv(t *testing.T) {
	t.Setenv("ANTHROPIC_AWS_BASE_URL", "https://custom.gateway.example.com")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "my-key")
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "")

	resolved, err := awsauth.ResolveConfig(toInternalConfig(ClientConfig{
		WorkspaceID: "default",
	}), awsResolveParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved.BaseURL != "https://custom.gateway.example.com" {
		t.Errorf("expected base URL from env, got %q", resolved.BaseURL)
	}
}

func TestBaseURLExplicitOverridesRegion(t *testing.T) {
	t.Setenv("ANTHROPIC_AWS_BASE_URL", "")
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	resolved, err := awsauth.ResolveConfig(toInternalConfig(ClientConfig{
		BaseURL:     "https://explicit.example.com",
		AWSRegion:   "us-east-1",
		WorkspaceID: "default",
		APIKey:      "my-key",
	}), awsResolveParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved.BaseURL != "https://explicit.example.com" {
		t.Errorf("expected explicit base URL, got %q", resolved.BaseURL)
	}
}

// --- skipAuth tests ---

func TestSkipAuthNoWorkspaceRequired(t *testing.T) {
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	_, err := NewClient(context.Background(), ClientConfig{
		BaseURL:  "https://proxy.example.com",
		SkipAuth: true,
	})
	if err != nil {
		t.Fatalf("expected no error with skipAuth, got: %v", err)
	}
}

func TestSkipAuthNoAuthHeaders(t *testing.T) {
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "should-not-appear")

	client, captured := newTestClient(t, ClientConfig{
		SkipAuth: true,
	})
	sendTestRequest(t, client)

	if captured.Headers.Get("Authorization") != "" {
		t.Error("expected no Authorization header with skipAuth")
	}
	if captured.Headers.Get("X-Amz-Date") != "" {
		t.Error("expected no X-Amz-Date header with skipAuth")
	}
	if captured.Headers.Get("Anthropic-Workspace-Id") != "" {
		t.Error("expected no anthropic-workspace-id header with skipAuth")
	}
}

func TestSkipAuthIgnoresProvidedCreds(t *testing.T) {
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "")

	client, captured := newTestClient(t, ClientConfig{
		AWSAccessKey:       "key",
		AWSSecretAccessKey: "secret",
		AWSRegion:          "us-east-1",
		WorkspaceID:        "ws-123",
		SkipAuth:           true,
	})
	sendTestRequest(t, client)

	if captured.Headers.Get("Authorization") != "" {
		t.Error("expected no Authorization header with skipAuth even when credentials are provided")
	}
	if captured.Headers.Get("Anthropic-Workspace-Id") != "" {
		t.Error("expected no anthropic-workspace-id header with skipAuth even when workspace ID is provided")
	}
}

// --- ANTHROPIC_API_KEY isolation test ---

func TestDoesNotLeakAnthropicAPIKey(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "should-not-appear")
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	client, captured := newTestClient(t, ClientConfig{
		AWSRegion:          "us-east-1",
		AWSAccessKey:       "test-access-key",
		AWSSecretAccessKey: "test-secret-key",
		WorkspaceID:        "ws-test",
	})
	sendTestRequest(t, client)

	if got := captured.Headers.Get("X-Api-Key"); got != "" {
		t.Errorf("expected no x-api-key header when using SigV4, got %q", got)
	}
}

// --- Middleware unit tests ---

func TestAWSMiddlewareSigV4Signing(t *testing.T) {
	cfg := makeStaticConfig("us-east-1")
	signer := v4.NewSigner()
	middleware := awsauth.SigV4Middleware(signer, cfg, defaultServiceName)

	body := []byte(`{"messages":[{"role":"user","content":"hello"}]}`)
	req, err := http.NewRequest("POST", "https://aws-external-anthropic.us-east-1.api.aws/v1/messages", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = middleware(req, func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header to be set by SigV4 signing")
		}
		if r.Header.Get("X-Amz-Date") == "" {
			t.Error("expected X-Amz-Date header to be set by SigV4 signing")
		}
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		if !bytes.Equal(bodyBytes, body) {
			t.Errorf("expected body %q, got %q", body, bodyBytes)
		}
		return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
	})
	if err != nil {
		t.Fatalf("middleware failed: %v", err)
	}
}

func TestAWSMiddlewareEmptyBody(t *testing.T) {
	cfg := makeStaticConfig("us-east-1")
	signer := v4.NewSigner()
	middleware := awsauth.SigV4Middleware(signer, cfg, defaultServiceName)

	req, err := http.NewRequest("GET", "https://aws-external-anthropic.us-east-1.api.aws/v1/models", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	_, err = middleware(req, func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Authorization") == "" {
			t.Error("expected Authorization header from SigV4 signing even with empty body")
		}
		return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
	})
	if err != nil {
		t.Fatalf("middleware failed: %v", err)
	}
}

func TestAWSMiddlewareServiceName(t *testing.T) {
	cfg := makeStaticConfig("eu-west-1")
	signer := v4.NewSigner()
	middleware := awsauth.SigV4Middleware(signer, cfg, defaultServiceName)

	req, err := http.NewRequest("POST", "https://aws-external-anthropic.eu-west-1.api.aws/v1/messages", bytes.NewReader([]byte(`{}`)))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	var gotAuth string
	_, err = middleware(req, func(r *http.Request) (*http.Response, error) {
		gotAuth = r.Header.Get("Authorization")
		return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
	})
	if err != nil {
		t.Fatalf("middleware failed: %v", err)
	}

	if gotAuth == "" {
		t.Fatal("expected non-empty Authorization header")
	}
	expectedServiceFragment := "/" + defaultServiceName + "/"
	if !bytes.Contains([]byte(gotAuth), []byte(expectedServiceFragment)) {
		t.Errorf("expected Authorization to contain service %q, got: %s", defaultServiceName, gotAuth)
	}
}

func TestAWSMiddlewareQueryParamsIncludedInSignature(t *testing.T) {
	cfg := makeStaticConfig("us-east-1")
	signer := v4.NewSigner()
	middleware := awsauth.SigV4Middleware(signer, cfg, defaultServiceName)

	body := []byte(`{}`)
	urlWithQuery := "https://aws-external-anthropic.us-east-1.api.aws/v1/messages?beta=true&stream=false"
	req, err := http.NewRequest("POST", urlWithQuery, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	var signedAuth string
	_, err = middleware(req, func(r *http.Request) (*http.Response, error) {
		signedAuth = r.Header.Get("Authorization")
		if r.URL.RawQuery == "" {
			t.Error("expected query params to be preserved after signing")
		}
		if r.URL.Query().Get("beta") != "true" || r.URL.Query().Get("stream") != "false" {
			t.Errorf("query params were modified: got %s", r.URL.RawQuery)
		}
		return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
	})
	if err != nil {
		t.Fatalf("middleware failed: %v", err)
	}

	middlewareNoQuery := awsauth.SigV4Middleware(signer, cfg, defaultServiceName)
	reqNoQuery, err := http.NewRequest("POST", "https://aws-external-anthropic.us-east-1.api.aws/v1/messages", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	var signedAuthNoQuery string
	_, err = middlewareNoQuery(reqNoQuery, func(r *http.Request) (*http.Response, error) {
		signedAuthNoQuery = r.Header.Get("Authorization")
		return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
	})
	if err != nil {
		t.Fatalf("middleware failed: %v", err)
	}

	if signedAuth == signedAuthNoQuery {
		t.Error("expected different signatures for requests with and without query params")
	}
}

func TestAWSMiddlewareBodyReplayed(t *testing.T) {
	cfg := makeStaticConfig("us-east-1")
	signer := v4.NewSigner()
	middleware := awsauth.SigV4Middleware(signer, cfg, defaultServiceName)

	originalBody := []byte(`{"test":"body content"}`)
	req, err := http.NewRequest("POST", "https://aws-external-anthropic.us-east-1.api.aws/v1/messages", bytes.NewReader(originalBody))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	_, err = middleware(req, func(r *http.Request) (*http.Response, error) {
		if r.GetBody == nil {
			t.Error("expected GetBody to be set for retry support")
		}
		rc, err := r.GetBody()
		if err != nil {
			t.Fatalf("GetBody failed: %v", err)
		}
		replayed, err := io.ReadAll(rc)
		if err != nil {
			t.Fatalf("failed to read replayed body: %v", err)
		}
		if !bytes.Equal(replayed, originalBody) {
			t.Errorf("replayed body %q does not match original %q", replayed, originalBody)
		}
		return &http.Response{StatusCode: 200, Body: http.NoBody}, nil
	})
	if err != nil {
		t.Fatalf("middleware failed: %v", err)
	}
}

// --- User profiles ---

func TestUserProfilesCRUD(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name       string
		send       func(client *Client) error
		wantMethod string
		wantPath   string
	}{
		{
			name: "create",
			send: func(client *Client) error {
				_, err := client.UserProfiles.New(ctx, anthropic.UserProfileNewParams{Name: anthropic.String("Example")})
				return err
			},
			wantMethod: http.MethodPost,
			wantPath:   "/v1/user_profiles",
		},
		{
			name: "get",
			send: func(client *Client) error {
				_, err := client.UserProfiles.Get(ctx, "uprof_1", anthropic.UserProfileGetParams{})
				return err
			},
			wantMethod: http.MethodGet,
			wantPath:   "/v1/user_profiles/uprof_1",
		},
		{
			name: "update",
			send: func(client *Client) error {
				_, err := client.UserProfiles.Update(ctx, "uprof_1", anthropic.UserProfileUpdateParams{Name: anthropic.String("Example")})
				return err
			},
			wantMethod: http.MethodPost,
			wantPath:   "/v1/user_profiles/uprof_1",
		},
		{
			name: "list",
			send: func(client *Client) error {
				_, err := client.UserProfiles.List(ctx, anthropic.UserProfileListParams{})
				return err
			},
			wantMethod: http.MethodGet,
			wantPath:   "/v1/user_profiles",
		},
		{
			name: "delete",
			send: func(client *Client) error {
				_, err := client.UserProfiles.Delete(ctx, "uprof_1", anthropic.UserProfileDeleteParams{})
				return err
			},
			wantMethod: http.MethodDelete,
			wantPath:   "/v1/user_profiles/uprof_1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotMethod, gotPath string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotMethod, gotPath = r.Method, r.URL.Path
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{}`))
			}))
			t.Cleanup(server.Close)

			client, err := NewClient(ctx, ClientConfig{APIKey: "test-key", WorkspaceID: "ws-123", BaseURL: server.URL})
			if err != nil {
				t.Fatalf("NewClient failed: %v", err)
			}
			if err := tt.send(client); err != nil {
				t.Fatalf("request failed: %v", err)
			}
			if gotMethod != tt.wantMethod || gotPath != tt.wantPath {
				t.Errorf("got %s %s, want %s %s", gotMethod, gotPath, tt.wantMethod, tt.wantPath)
			}
		})
	}
}

func TestUserProfileIDHeader(t *testing.T) {
	client, captured := newTestClient(t, ClientConfig{APIKey: "test-key", WorkspaceID: "ws-123"})

	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:         "claude-sonnet-4-6-20250514",
		MaxTokens:     1,
		Messages:      []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hi"))},
		UserProfileID: anthropic.String("uprof_1"),
	})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if got := captured.Headers.Get("anthropic-user-profile-id"); got != "uprof_1" {
		t.Errorf("anthropic-user-profile-id = %q, want %q", got, "uprof_1")
	}
}

func TestUserProfileActionURLsAreRejected(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name string
		send func(client *Client) error
	}{
		{
			name: "stable",
			send: func(client *Client) error {
				_, err := client.UserProfiles.ActionURLs.New(ctx, "uprof_1", anthropic.UserProfileActionURLNewParams{
					Action: anthropic.UserProfileActionUnionParam{OfEnrollEndUser: &anthropic.UserProfileEnrollEndUserActionParam{}},
				})
				return err
			},
		},
		{
			name: "beta",
			send: func(client *Client) error {
				_, err := client.Beta.UserProfiles.ActionURLs.New(ctx, "uprof_1", anthropic.BetaUserProfileActionURLNewParams{
					Action: anthropic.BetaUserProfileActionUnionParam{OfEnrollEndUser: &anthropic.BetaUserProfileEnrollEndUserActionParam{}},
				})
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requests := 0
			server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { requests++ }))
			t.Cleanup(server.Close)

			client, err := NewClient(ctx, ClientConfig{APIKey: "test-key", WorkspaceID: "ws-123", BaseURL: server.URL}, option.WithMaxRetries(0))
			if err != nil {
				t.Fatalf("NewClient failed: %v", err)
			}

			err = tt.send(client)
			if err == nil || !strings.Contains(err.Error(), "does not serve user profile action URLs") {
				t.Errorf("expected the action URLs error, got %v", err)
			}
			if requests != 0 {
				t.Errorf("expected no request to be sent, got %d", requests)
			}
		})
	}
}
