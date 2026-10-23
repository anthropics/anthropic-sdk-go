package bedrock

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/tidwall/gjson"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/awsauth"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

func makeStaticAWSConfig(region string) awssdk.Config {
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

// writeMessagesResponse writes a canned successful Messages API response.
func writeMessagesResponse(w http.ResponseWriter) {
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

type mantleCapturedRequest struct {
	Headers http.Header
	URL     string
}

func mantleMessagesHandler(captured *mantleCapturedRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		captured.Headers = r.Header.Clone()
		captured.URL = r.URL.String()
		writeMessagesResponse(w)
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

	if got := captured.Headers.Get("Authorization"); got != "Bearer my-api-key" {
		t.Errorf("expected Authorization %q, got %q", "Bearer my-api-key", got)
	}
	if captured.Headers.Get("X-Api-Key") != "" {
		t.Error("expected no X-Api-Key header in bearer auth mode")
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
	cfg := makeStaticAWSConfig("us-east-1")
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

	if got := captured.Headers.Get("Authorization"); got != "Bearer aws-fallback-key" {
		t.Errorf("expected Authorization %q (AWS fallback), got %q", "Bearer aws-fallback-key", got)
	}
}

func TestMantleAPIKeyMantleEnvOverridesAWSEnv(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "mantle-key")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "aws-key")
	t.Setenv("AWS_REGION", "us-east-1")

	client, captured := newTestMantleClient(t, MantleClientConfig{})
	sendTestMantleRequest(t, client)

	if got := captured.Headers.Get("Authorization"); got != "Bearer mantle-key" {
		t.Errorf("expected Authorization %q (mantle-specific), got %q", "Bearer mantle-key", got)
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
	t.Setenv("ANTHROPIC_API_KEY", "should-not-appear")

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

	if got := captured.Headers.Get("Authorization"); got != "Bearer test-key" {
		t.Errorf("expected Authorization %q, got %q", "Bearer test-key", got)
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

	// User-provided opts should override internal opts (e.g. override the Authorization header)
	client, captured := newTestMantleClient(t, MantleClientConfig{
		APIKey:    "internal-key",
		AWSRegion: "us-east-1",
	}, option.WithHeader("Authorization", "Bearer override-key"))
	sendTestMantleRequest(t, client)

	if got := captured.Headers.Get("Authorization"); got != "Bearer override-key" {
		t.Errorf("expected Authorization %q (from RequestOption override), got %q", "Bearer override-key", got)
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

// TestMantleUserMiddlewareRunsBeforeSigV4Signing verifies that SigV4 signing
// runs after (inside) middleware registered through client options, so the
// signature covers the middleware's request mutations.
func TestMantleUserMiddlewareRunsBeforeSigV4Signing(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "")

	var observedAuth, observedAmzDate string
	mutator := func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		observedAuth = r.Header.Get("Authorization")
		observedAmzDate = r.Header.Get("X-Amz-Date")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		var parsed map[string]any
		if err := json.Unmarshal(body, &parsed); err != nil {
			return nil, err
		}
		parsed["metadata"] = map[string]any{"user_id": "injected-by-middleware"}
		mutated, err := json.Marshal(parsed)
		if err != nil {
			return nil, err
		}
		r.Body = io.NopCloser(bytes.NewReader(mutated))
		r.ContentLength = int64(len(mutated))
		return next(r)
	}

	var wireBody map[string]any
	var wireAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wireAuth = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&wireBody); err != nil {
			t.Errorf("failed to decode wire body: %v", err)
		}
		writeMessagesResponse(w)
	}))
	t.Cleanup(server.Close)

	client, err := NewMantleClient(context.Background(), MantleClientConfig{
		BaseURL:            server.URL,
		AWSRegion:          "us-east-1",
		AWSAccessKey:       "test-access-key",
		AWSSecretAccessKey: "test-secret-key",
	}, option.WithMiddleware(mutator))
	if err != nil {
		t.Fatalf("NewMantleClient failed: %v", err)
	}
	sendTestMantleRequest(t, client)

	// The user middleware runs before signing, so no signature exists yet.
	if observedAuth != "" {
		t.Errorf("expected user middleware to run before signing, but saw Authorization %q", observedAuth)
	}
	if observedAmzDate != "" {
		t.Errorf("expected user middleware to run before signing, but saw X-Amz-Date %q", observedAmzDate)
	}

	// The wire request is signed and carries the middleware's mutation.
	if !strings.HasPrefix(wireAuth, "AWS4-HMAC-SHA256") {
		t.Errorf("expected SigV4 Authorization on the wire, got %q", wireAuth)
	}
	metadata, _ := wireBody["metadata"].(map[string]any)
	if metadata["user_id"] != "injected-by-middleware" {
		t.Errorf("expected middleware mutation in the signed wire body, got metadata %v", wireBody["metadata"])
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
		t.Errorf("expected no X-Api-Key header when using SigV4, got %q", got)
	}
	// Authorization header should be SigV4, not Bearer
	if auth := captured.Headers.Get("Authorization"); !strings.HasPrefix(auth, "AWS4-HMAC-SHA256") {
		t.Errorf("expected SigV4 Authorization header, got %q", auth)
	}
}

// --- User profile ID tests ---

type mantleWireRequest struct {
	Path    string
	Query   string
	Headers http.Header
	Body    []byte
}

// newMantleWireServer answers Messages requests and records each request as it
// arrived, body included.
func newMantleWireServer(t *testing.T) (*httptest.Server, *mantleWireRequest) {
	t.Helper()
	var wire mantleWireRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read wire body: %v", err)
		}
		wire = mantleWireRequest{Path: r.URL.Path, Query: r.URL.RawQuery, Headers: r.Header.Clone(), Body: body}

		switch {
		case strings.HasSuffix(r.URL.Path, "/count_tokens"):
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"input_tokens":1}`)
		case strings.HasSuffix(r.URL.Path, "/batches"):
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"id":"msgbatch_test","type":"message_batch","processing_status":"in_progress"}`)
		case gjson.GetBytes(body, "stream").Bool():
			w.Header().Set("Content-Type", "text/event-stream")
			io.WriteString(w, "event: message_stop\ndata: {\"type\":\"message_stop\"}\n\n")
		default:
			writeMessagesResponse(w)
		}
	}))
	t.Cleanup(server.Close)
	return server, &wire
}

func TestMantleUserProfileID(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	messages := []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hi"))}
	betaMessages := []anthropic.BetaMessageParam{anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("hi"))}
	const model = "claude-sonnet-4-6-20250514"

	tests := []struct {
		name       string
		send       func(ctx context.Context, client *MantleClient) error
		wantPath   string
		wantBody   string // the body's user_profile_id; "" when the field must be absent
		wantHeader string
	}{
		{
			name: "messages",
			send: func(ctx context.Context, client *MantleClient) error {
				_, err := client.Messages.New(ctx, anthropic.MessageNewParams{
					Model: model, MaxTokens: 1, Messages: messages, UserProfileID: anthropic.String("uprof_1"),
				})
				return err
			},
			wantPath: "/v1/messages",
			wantBody: "uprof_1",
		},
		{
			name: "beta messages",
			send: func(ctx context.Context, client *MantleClient) error {
				_, err := client.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{
					Model: model, MaxTokens: 1, Messages: betaMessages, UserProfileID: anthropic.String("uprof_1"),
				})
				return err
			},
			wantPath: "/v1/messages",
			wantBody: "uprof_1",
		},
		{
			name: "streaming",
			send: func(ctx context.Context, client *MantleClient) error {
				stream := client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
					Model: model, MaxTokens: 1, Messages: messages, UserProfileID: anthropic.String("uprof_1"),
				})
				for stream.Next() {
				}
				return stream.Err()
			},
			wantPath: "/v1/messages",
			wantBody: "uprof_1",
		},
		{
			name: "beta streaming",
			send: func(ctx context.Context, client *MantleClient) error {
				stream := client.Beta.Messages.NewStreaming(ctx, anthropic.BetaMessageNewParams{
					Model: model, MaxTokens: 1, Messages: betaMessages, UserProfileID: anthropic.String("uprof_1"),
				})
				for stream.Next() {
				}
				return stream.Err()
			},
			wantPath: "/v1/messages",
			wantBody: "uprof_1",
		},
		{
			name: "header set by the caller",
			send: func(ctx context.Context, client *MantleClient) error {
				_, err := client.Messages.New(ctx, anthropic.MessageNewParams{
					Model: model, MaxTokens: 1, Messages: messages,
				}, option.WithHeader("anthropic-user-profile-id", "uprof_1"))
				return err
			},
			wantPath: "/v1/messages",
			wantBody: "uprof_1",
		},
		{
			name: "body value set by the caller wins",
			send: func(ctx context.Context, client *MantleClient) error {
				_, err := client.Messages.New(ctx, anthropic.MessageNewParams{
					Model: model, MaxTokens: 1, Messages: messages, UserProfileID: anthropic.String("uprof_1"),
				}, option.WithJSONSet("user_profile_id", "uprof_caller"))
				return err
			},
			wantPath: "/v1/messages",
			wantBody: "uprof_caller",
		},
		{
			name: "no profile ID",
			send: func(ctx context.Context, client *MantleClient) error {
				_, err := client.Messages.New(ctx, anthropic.MessageNewParams{
					Model: model, MaxTokens: 1, Messages: messages,
				})
				return err
			},
			wantPath: "/v1/messages",
		},
		{
			name: "body that is not a JSON object",
			send: func(ctx context.Context, client *MantleClient) error {
				_, err := client.Messages.New(ctx, anthropic.MessageNewParams{
					Model: model, MaxTokens: 1, Messages: messages, UserProfileID: anthropic.String("uprof_1"),
				}, option.WithRequestBody("application/json", []byte(`["not an object"]`)))
				return err
			},
			wantPath:   "/v1/messages",
			wantHeader: "uprof_1",
		},
		{
			name: "count tokens",
			send: func(ctx context.Context, client *MantleClient) error {
				_, err := client.Messages.CountTokens(ctx, anthropic.MessageCountTokensParams{
					Model: model, Messages: messages, UserProfileID: anthropic.String("uprof_1"),
				})
				return err
			},
			wantPath: "/v1/messages/count_tokens",
			wantBody: "uprof_1",
		},
		{
			name: "beta count tokens",
			send: func(ctx context.Context, client *MantleClient) error {
				_, err := client.Beta.Messages.CountTokens(ctx, anthropic.BetaMessageCountTokensParams{
					Model: model, Messages: betaMessages, UserProfileID: anthropic.String("uprof_1"),
				})
				return err
			},
			wantPath: "/v1/messages/count_tokens",
			wantBody: "uprof_1",
		},
		{
			name: "batches",
			send: func(ctx context.Context, client *MantleClient) error {
				_, err := client.Messages.Batches.New(ctx, anthropic.MessageBatchNewParams{
					Requests: []anthropic.MessageBatchNewParamsRequest{{
						CustomID: "request-1",
						Params:   anthropic.MessageBatchNewParamsRequestParams{Model: model, MaxTokens: 1, Messages: messages},
					}},
					UserProfileID: anthropic.String("uprof_1"),
				})
				return err
			},
			wantPath:   "/v1/messages/batches",
			wantHeader: "uprof_1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server, wire := newMantleWireServer(t)
			client, err := NewMantleClient(context.Background(), MantleClientConfig{
				APIKey:  "test-key",
				BaseURL: server.URL,
			})
			if err != nil {
				t.Fatalf("NewMantleClient failed: %v", err)
			}

			if err := test.send(context.Background(), client); err != nil {
				t.Fatalf("request failed: %v", err)
			}

			if wire.Path != test.wantPath || wire.Query != "" {
				t.Errorf("request went to %s?%s, want %s", wire.Path, wire.Query, test.wantPath)
			}
			if got := wire.Headers.Get("anthropic-user-profile-id"); got != test.wantHeader {
				t.Errorf("anthropic-user-profile-id header is %q, want %q", got, test.wantHeader)
			}
			field := gjson.GetBytes(wire.Body, "user_profile_id")
			if field.Exists() != (test.wantBody != "") || field.String() != test.wantBody {
				t.Errorf("body user_profile_id is %s, want %q; body: %s", field.Raw, test.wantBody, wire.Body)
			}
			if got := wire.Headers.Get("Content-Length"); got != strconv.Itoa(len(wire.Body)) {
				t.Errorf("Content-Length is %s for a body of %d bytes", got, len(wire.Body))
			}
		})
	}
}

// TestMantleUserProfileIDIsSigned verifies that the SigV4 signature covers the
// request as sent: the body that carries the profile ID, and no profile header.
func TestMantleUserProfileIDIsSigned(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "")

	server, wire := newMantleWireServer(t)
	client, err := NewMantleClient(context.Background(), MantleClientConfig{
		BaseURL:            server.URL,
		AWSRegion:          "us-east-1",
		AWSAccessKey:       "test-access-key",
		AWSSecretAccessKey: "test-secret-key",
	})
	if err != nil {
		t.Fatalf("NewMantleClient failed: %v", err)
	}

	_, err = client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:         "claude-sonnet-4-6-20250514",
		MaxTokens:     1,
		Messages:      []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hi"))},
		UserProfileID: anthropic.String("uprof_1"),
	})
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if got := gjson.GetBytes(wire.Body, "user_profile_id").String(); got != "uprof_1" {
		t.Errorf("body user_profile_id is %q, want %q; body: %s", got, "uprof_1", wire.Body)
	}
	if got := wire.Headers.Get("anthropic-user-profile-id"); got != "" {
		t.Errorf("anthropic-user-profile-id header is %q, want none", got)
	}

	// Sign the request as it arrived, with the same key and time, and expect
	// the Authorization header it came with.
	auth := wire.Headers.Get("Authorization")
	_, signedHeaders, found := strings.Cut(auth, "SignedHeaders=")
	if !found {
		t.Fatalf("expected a SigV4 Authorization header, got %q", auth)
	}
	signedHeaders, _, _ = strings.Cut(signedHeaders, ",")

	resigned, err := http.NewRequest(http.MethodPost, server.URL+wire.Path, bytes.NewReader(wire.Body))
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	for _, name := range strings.Split(signedHeaders, ";") {
		if name != "host" && name != "content-length" {
			resigned.Header[http.CanonicalHeaderKey(name)] = wire.Headers.Values(name)
		}
	}
	signingTime, err := time.Parse("20060102T150405Z", wire.Headers.Get("X-Amz-Date"))
	if err != nil {
		t.Fatalf("failed to parse X-Amz-Date: %v", err)
	}
	creds, err := makeStaticAWSConfig("us-east-1").Credentials.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("failed to retrieve credentials: %v", err)
	}
	hash := sha256.Sum256(wire.Body)
	err = awsauth.SignHTTP(context.Background(), v4.NewSigner(), creds, resigned, hex.EncodeToString(hash[:]), mantleServiceName, "us-east-1", signingTime)
	if err != nil {
		t.Fatalf("failed to sign: %v", err)
	}
	if got := resigned.Header.Get("Authorization"); got != auth {
		t.Errorf("signature does not cover the request as sent:\n sent: %s\n want: %s", auth, got)
	}
}

// Mantle serves only the Messages API, so a user profiles service must not appear on either surface.
func TestMantleClientExposesOnlyMessages(t *testing.T) {
	for _, service := range []any{MantleClient{}, MantleBetaService{}} {
		serviceType := reflect.TypeOf(service)
		for i := range serviceType.NumField() {
			switch name := serviceType.Field(i).Name; name {
			case "Options", "Messages", "Beta":
			default:
				t.Errorf("%s has unexpected service %s", serviceType.Name(), name)
			}
		}
	}
}
