package bedrock

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/awsauth"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func makeMantleStaticConfig(region string) awssdk.Config {
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

type mantleCapturedRequest struct {
	Headers http.Header
	URL     string
}

func mantleMessagesHandler(captured *mantleCapturedRequest) http.HandlerFunc {
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

// newTestMantleClient creates a MantleClient pointed at a test server and returns
// the client and a struct that captures the request headers/URL.
func newTestMantleClient(t *testing.T, cfg MantleClientConfig, opts ...option.RequestOption) (*MantleClient, *mantleCapturedRequest) {
	t.Helper()
	var captured mantleCapturedRequest
	server := httptest.NewServer(mantleMessagesHandler(&captured))
	t.Cleanup(server.Close)

	cfg.BaseURL = server.URL
	client, err := NewMantleClient(context.Background(), cfg, opts...)
	if err != nil {
		t.Fatalf("NewMantleClient failed: %v", err)
	}
	return client, &captured
}

func sendTestMantleRequest(t *testing.T, client *MantleClient) {
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

// --- Validation tests ---

func TestMantleRequiresBaseURLOrRegion(t *testing.T) {
	t.Setenv("AWS_REGION", "")
	t.Setenv("ANTHROPIC_BEDROCK_MANTLE_BASE_URL", "")

	_, err := NewMantleClient(context.Background(), MantleClientConfig{
		APIKey: "my-key",
	})
	if err == nil {
		t.Fatal("expected error when neither base URL nor region is available")
	}
}

func TestMantleRegionRequiredForSigV4(t *testing.T) {
	t.Setenv("AWS_REGION", "")
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	_, err := NewMantleClient(context.Background(), MantleClientConfig{
		AWSAccessKey:       "key",
		AWSSecretAccessKey: "secret",
	})
	if err == nil {
		t.Fatal("expected error when region is missing and SigV4 is used")
	}
}

// --- API key mode tests ---

func TestMantleAPIKeyModeHeaders(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	client, captured := newTestMantleClient(t, MantleClientConfig{
		APIKey:    "my-api-key",
		AWSRegion: "us-east-1",
	})
	sendTestMantleRequest(t, client)

	if got := captured.Headers.Get("X-Api-Key"); got != "my-api-key" {
		t.Errorf("expected x-api-key %q, got %q", "my-api-key", got)
	}
	if captured.Headers.Get("Authorization") != "" {
		t.Error("expected no Authorization header in API key mode")
	}
}

// --- SigV4 mode tests ---

func TestMantleSigV4ModeHeaders(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "")

	client, captured := newTestMantleClient(t, MantleClientConfig{
		AWSRegion:          "us-east-1",
		AWSAccessKey:       "test-access-key",
		AWSSecretAccessKey: "test-secret-key",
	})
	sendTestMantleRequest(t, client)

	auth := captured.Headers.Get("Authorization")
	if auth == "" {
		t.Fatal("expected Authorization header in SigV4 mode")
	}
	if !strings.HasPrefix(auth, "AWS4-HMAC-SHA256") {
		t.Errorf("expected AWS4 signature, got: %s", auth)
	}
	if !strings.Contains(auth, mantleServiceName) {
		t.Errorf("expected service name %q in Authorization, got: %s", mantleServiceName, auth)
	}
	if captured.Headers.Get("X-Amz-Date") == "" {
		t.Error("expected X-Amz-Date header in SigV4 mode")
	}
}

// --- SigV4 middleware service name test ---

func TestMantleSigV4ServiceName(t *testing.T) {
	cfg := makeMantleStaticConfig("us-east-1")
	signer := v4.NewSigner()
	middleware := awsauth.SigV4Middleware(signer, cfg, mantleServiceName)

	req, err := http.NewRequest("POST", "https://bedrock-mantle.us-east-1.api.aws/anthropic/v1/messages", bytes.NewReader([]byte(`{}`)))
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
	expectedServiceFragment := "/" + mantleServiceName + "/"
	if !bytes.Contains([]byte(gotAuth), []byte(expectedServiceFragment)) {
		t.Errorf("expected Authorization to contain service %q, got: %s", mantleServiceName, gotAuth)
	}
}

// --- Env var fallback tests ---

func TestMantleAPIKeyFallbackToAWSEnv(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "aws-fallback-key")
	t.Setenv("AWS_REGION", "us-east-1")

	client, captured := newTestMantleClient(t, MantleClientConfig{})
	sendTestMantleRequest(t, client)

	if got := captured.Headers.Get("X-Api-Key"); got != "aws-fallback-key" {
		t.Errorf("expected x-api-key %q (AWS fallback), got %q", "aws-fallback-key", got)
	}
}

func TestMantleAPIKeyMantleEnvOverridesAWSEnv(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "mantle-key")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "aws-key")
	t.Setenv("AWS_REGION", "us-east-1")

	client, captured := newTestMantleClient(t, MantleClientConfig{})
	sendTestMantleRequest(t, client)

	if got := captured.Headers.Get("X-Api-Key"); got != "mantle-key" {
		t.Errorf("expected x-api-key %q (mantle-specific), got %q", "mantle-key", got)
	}
}

// --- Base URL tests ---

func TestMantleBaseURLDerivedFromRegion(t *testing.T) {
	t.Setenv("ANTHROPIC_BEDROCK_MANTLE_BASE_URL", "")
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	resolved, err := awsauth.ResolveConfig(mantleToInternalConfig(MantleClientConfig{
		AWSRegion: "us-west-2",
		APIKey:    "my-key",
	}), mantleResolveParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "https://bedrock-mantle.us-west-2.api.aws/anthropic"
	if resolved.BaseURL != expected {
		t.Errorf("expected base URL %q, got %q", expected, resolved.BaseURL)
	}
}

func TestMantleBaseURLFromEnv(t *testing.T) {
	t.Setenv("ANTHROPIC_BEDROCK_MANTLE_BASE_URL", "https://custom.mantle.example.com")
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "my-key")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	resolved, err := awsauth.ResolveConfig(mantleToInternalConfig(MantleClientConfig{}), mantleResolveParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved.BaseURL != "https://custom.mantle.example.com" {
		t.Errorf("expected base URL from env, got %q", resolved.BaseURL)
	}
}

func TestMantleBaseURLExplicitOverridesRegion(t *testing.T) {
	t.Setenv("ANTHROPIC_BEDROCK_MANTLE_BASE_URL", "")
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	resolved, err := awsauth.ResolveConfig(mantleToInternalConfig(MantleClientConfig{
		BaseURL:   "https://explicit.example.com",
		AWSRegion: "us-east-1",
		APIKey:    "my-key",
	}), mantleResolveParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved.BaseURL != "https://explicit.example.com" {
		t.Errorf("expected explicit base URL, got %q", resolved.BaseURL)
	}
}

// --- skipAuth tests ---

func TestMantleSkipAuthNoAuthRequired(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	_, err := NewMantleClient(context.Background(), MantleClientConfig{
		BaseURL:  "https://proxy.example.com",
		SkipAuth: true,
	})
	if err != nil {
		t.Fatalf("expected no error with skipAuth, got: %v", err)
	}
}

func TestMantleSkipAuthNoAuthHeaders(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	client, captured := newTestMantleClient(t, MantleClientConfig{
		SkipAuth: true,
	})
	sendTestMantleRequest(t, client)

	if captured.Headers.Get("Authorization") != "" {
		t.Error("expected no Authorization header with skipAuth")
	}
	if captured.Headers.Get("X-Amz-Date") != "" {
		t.Error("expected no X-Amz-Date header with skipAuth")
	}
}

// --- NewMantleClient tests ---

func TestNewMantleClientMessages(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	client, captured := newTestMantleClient(t, MantleClientConfig{
		APIKey:    "test-key",
		AWSRegion: "us-east-1",
	})
	sendTestMantleRequest(t, client)

	if got := captured.Headers.Get("X-Api-Key"); got != "test-key" {
		t.Errorf("expected x-api-key %q, got %q", "test-key", got)
	}
}

func TestNewMantleClientBetaMessages(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	var captured mantleCapturedRequest
	server := httptest.NewServer(mantleMessagesHandler(&captured))
	t.Cleanup(server.Close)

	client, err := NewMantleClient(context.Background(), MantleClientConfig{
		APIKey:    "test-key",
		AWSRegion: "us-east-1",
		BaseURL:   server.URL,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.Beta.Messages.New(context.Background(), anthropic.BetaMessageNewParams{
		Model:     "claude-sonnet-4-6-20250514",
		MaxTokens: 1,
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("hi")),
		},
	})
	if err != nil {
		t.Fatalf("beta messages request failed: %v", err)
	}
}

// --- RequestOption tests ---

func TestMantleClientRequestOptionsCustomHeader(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	client, captured := newTestMantleClient(t, MantleClientConfig{
		APIKey:    "test-key",
		AWSRegion: "us-east-1",
	}, option.WithHeader("X-Custom-Header", "custom-value"))
	sendTestMantleRequest(t, client)

	if got := captured.Headers.Get("X-Custom-Header"); got != "custom-value" {
		t.Errorf("expected X-Custom-Header %q, got %q", "custom-value", got)
	}
}

func TestMantleClientRequestOptionsOverrideInternal(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	// User-provided opts should override internal opts (e.g. override the API key header)
	client, captured := newTestMantleClient(t, MantleClientConfig{
		APIKey:    "internal-key",
		AWSRegion: "us-east-1",
	}, option.WithAPIKey("override-key"))
	sendTestMantleRequest(t, client)

	if got := captured.Headers.Get("X-Api-Key"); got != "override-key" {
		t.Errorf("expected X-Api-Key %q (from RequestOption override), got %q", "override-key", got)
	}
}

func TestMantleClientPerRequestOptionsOverrideClientOptions(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	client, captured := newTestMantleClient(t, MantleClientConfig{
		APIKey:    "test-key",
		AWSRegion: "us-east-1",
	}, option.WithHeader("X-Custom-Header", "client-level"))

	// Send request with per-request option that overrides the client-level header
	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     "claude-sonnet-4-6-20250514",
		MaxTokens: 1,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("hi")),
		},
	}, option.WithHeader("X-Custom-Header", "request-level"))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if got := captured.Headers.Get("X-Custom-Header"); got != "request-level" {
		t.Errorf("expected X-Custom-Header %q (per-request override), got %q", "request-level", got)
	}
}

func TestMantleClientRequestOptionsMiddleware(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	middlewareCalled := false
	mw := func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		middlewareCalled = true
		r.Header.Set("X-From-Middleware", "true")
		return next(r)
	}

	client, captured := newTestMantleClient(t, MantleClientConfig{
		APIKey:    "test-key",
		AWSRegion: "us-east-1",
	}, option.WithMiddleware(mw))
	sendTestMantleRequest(t, client)

	if !middlewareCalled {
		t.Error("expected middleware to be called")
	}
	if got := captured.Headers.Get("X-From-Middleware"); got != "true" {
		t.Errorf("expected X-From-Middleware %q, got %q", "true", got)
	}
}

// --- ANTHROPIC_API_KEY isolation test ---

func TestMantleDoesNotLeakAnthropicAPIKey(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "should-not-appear")
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	client, captured := newTestMantleClient(t, MantleClientConfig{
		AWSRegion:          "us-east-1",
		AWSAccessKey:       "test-access-key",
		AWSSecretAccessKey: "test-secret-key",
	})
	sendTestMantleRequest(t, client)

	if got := captured.Headers.Get("X-Api-Key"); got != "" {
		t.Errorf("expected no x-api-key header when using SigV4, got %q", got)
	}
}
