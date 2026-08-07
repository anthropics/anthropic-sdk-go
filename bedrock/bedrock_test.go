package bedrock

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream"
	"github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream/eventstreamapi"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestBedrockURLEncoding(t *testing.T) {
	testCases := []struct {
		name            string
		model           string
		stream          bool
		expectedPath    string
		expectedRawPath string
	}{
		{
			name:            "regular model name",
			model:           "claude-3-sonnet",
			stream:          false,
			expectedPath:    "/model/claude-3-sonnet/invoke",
			expectedRawPath: "/model/claude-3-sonnet/invoke",
		},
		{
			name:            "regular model name with streaming",
			model:           "claude-3-sonnet",
			stream:          true,
			expectedPath:    "/model/claude-3-sonnet/invoke-with-response-stream",
			expectedRawPath: "/model/claude-3-sonnet/invoke-with-response-stream",
		},
		{
			name:            "inference profile ARN with slashes",
			model:           "arn:aws:bedrock:us-east-1:947123456126:application-inference-profile/xv9example4b",
			stream:          false,
			expectedPath:    "/model/arn:aws:bedrock:us-east-1:947123456126:application-inference-profile/xv9example4b/invoke",
			expectedRawPath: "/model/arn%3Aaws%3Abedrock%3Aus-east-1%3A947123456126%3Aapplication-inference-profile%2Fxv9example4b/invoke",
		},
		{
			name:            "inference profile ARN with streaming",
			model:           "arn:aws:bedrock:us-east-1:947123456126:application-inference-profile/xv9example4b",
			stream:          true,
			expectedPath:    "/model/arn:aws:bedrock:us-east-1:947123456126:application-inference-profile/xv9example4b/invoke-with-response-stream",
			expectedRawPath: "/model/arn%3Aaws%3Abedrock%3Aus-east-1%3A947123456126%3Aapplication-inference-profile%2Fxv9example4b/invoke-with-response-stream",
		},
		{
			name:            "foundation model ARN with colons",
			model:           "arn:aws:bedrock:us-east-1:123456789012:foundation-model/anthropic.claude-3-sonnet-20240229-v1:0",
			stream:          false,
			expectedPath:    "/model/arn:aws:bedrock:us-east-1:123456789012:foundation-model/anthropic.claude-3-sonnet-20240229-v1:0/invoke",
			expectedRawPath: "/model/arn%3Aaws%3Abedrock%3Aus-east-1%3A123456789012%3Afoundation-model%2Fanthropic.claude-3-sonnet-20240229-v1%3A0/invoke",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock AWS config
			cfg := aws.Config{
				Region: "us-east-1",
				Credentials: credentials.StaticCredentialsProvider{
					Value: aws.Credentials{
						AccessKeyID:     "test-access-key",
						SecretAccessKey: "test-secret-key",
					},
				},
			}

			signer := v4.NewSigner()
			middleware := bedrockMiddleware(signer, cfg)

			// Create request body
			requestBody := map[string]any{
				"model":  tc.model,
				"stream": tc.stream,
				"messages": []map[string]string{
					{"role": "user", "content": "Hello"},
				},
			}

			bodyBytes, err := json.Marshal(requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			// Create HTTP request
			req, err := http.NewRequest("POST", "https://bedrock-runtime.us-east-1.amazonaws.com/v1/messages", bytes.NewReader(bodyBytes))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Apply middleware
			_, err = middleware(req, func(r *http.Request) (*http.Response, error) {
				// Verify the URL paths are set correctly
				if r.URL.Path != tc.expectedPath {
					t.Errorf("Expected Path %q, got %q", tc.expectedPath, r.URL.Path)
				}

				if r.URL.RawPath != tc.expectedRawPath {
					t.Errorf("Expected RawPath %q, got %q", tc.expectedRawPath, r.URL.RawPath)
				}

				// Verify that the URL string contains the properly encoded path
				urlString := r.URL.String()
				expectedURL := fmt.Sprintf("https://bedrock-runtime.us-east-1.amazonaws.com%s", tc.expectedRawPath)
				if urlString != expectedURL {
					t.Errorf("Expected URL %q, got %q", expectedURL, urlString)
				}

				// Return a dummy response
				return &http.Response{
					StatusCode: 200,
					Body:       http.NoBody,
				}, nil
			})

			if err != nil {
				t.Fatalf("Middleware failed: %v", err)
			}
		})
	}
}

func TestBedrockBetaHeadersReRoutedThroughBody(t *testing.T) {
	// Create a mock AWS config
	cfg := aws.Config{
		Region: "us-east-1",
		Credentials: credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "test-access-key",
				SecretAccessKey: "test-secret-key",
			},
		},
	}

	signer := v4.NewSigner()
	middleware := bedrockMiddleware(signer, cfg)

	// Create HTTP request with beta headers
	type fakeRequest struct {
		Model         string              `json:"model"`
		AnthropicBeta []string            `json:"anthropic_beta,omitempty"`
		Messages      []map[string]string `json:"messages"`
	}
	reqBody := fakeRequest{
		Model: "fake-model",
		Messages: []map[string]string{
			{"role": "user", "content": "Hello"},
		},
	}
	requestBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://bedrock-runtime.us-east-1.amazonaws.com/v1/messages", bytes.NewReader(requestBodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("anthropic-beta", "beta-feature-1")
	req.Header.Add("anthropic-beta", "beta-feature-2")

	// Apply middleware
	_, err = middleware(req, func(r *http.Request) (*http.Response, error) {
		// Read the modified body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		var modifiedBody fakeRequest
		err = json.Unmarshal(bodyBytes, &modifiedBody)
		if err != nil {
			t.Fatalf("Failed to unmarshal modified body: %v", err)
		}

		// Verify that the anthropic_beta field is present in the body
		expectedBetas := []string{"beta-feature-1", "beta-feature-2"}
		if len(modifiedBody.AnthropicBeta) != len(expectedBetas) {
			t.Fatalf("Expected %d beta features, got %d", len(expectedBetas), len(modifiedBody.AnthropicBeta))
		}
		for i, beta := range expectedBetas {
			if modifiedBody.AnthropicBeta[i] != beta {
				t.Errorf("Expected beta feature %q, got %q", beta, modifiedBody.AnthropicBeta[i])
			}
		}

		// Return a dummy response
		return &http.Response{
			StatusCode: 200,
			Body:       http.NoBody,
		}, nil
	})

	if err != nil {
		t.Fatalf("Middleware failed: %v", err)
	}
}

func TestBedrockBearerToken(t *testing.T) {
	token := "test-bearer-token"
	region := "us-west-2"

	cfg := aws.Config{
		Region:                  region,
		BearerAuthTokenProvider: NewStaticBearerTokenProvider(token),
	}
	middleware := bedrockMiddleware(nil, cfg)

	requestBody := map[string]any{
		"model": "claude-3-sonnet",
		"messages": []map[string]string{
			{"role": "user", "content": "Hello"},
		},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://bedrock-runtime.us-west-2.amazonaws.com/v1/messages", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = middleware(req, func(r *http.Request) (*http.Response, error) {
		authHeader := r.Header.Get("Authorization")
		expectedAuth := "Bearer " + token
		if authHeader != expectedAuth {
			t.Errorf("Expected Authorization header %q, got %q", expectedAuth, authHeader)
		}

		if r.Header.Get("X-Amz-Date") != "" {
			t.Error("Expected no AWS SigV4 headers when using bearer token")
		}

		return &http.Response{
			StatusCode: 200,
			Body:       http.NoBody,
		}, nil
	})

	if err != nil {
		t.Fatalf("Middleware failed: %v", err)
	}
}

// TestBedrockWithConfigRequiresCredentials verifies that a config with no
// bearer token and no AWS credentials fails with a clear setup error rather
// than a nil-pointer panic on the first request.
func TestBedrockWithConfigRequiresCredentials(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")

	client := anthropic.NewClient(
		option.WithoutEnvironmentDefaults(),
		WithConfig(aws.Config{Region: "us-east-1"}),
	)

	_, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     "claude-3-sonnet",
		MaxTokens: 1,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("hi")),
		},
	})

	if err == nil || !strings.Contains(err.Error(), "expected AWS credentials to be set") {
		t.Fatalf("Expected credentials error, got: %v", err)
	}
}

func TestBedrockWithConfigUsesAWSHTTPClient(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")

	var customClientHits atomic.Int32
	customClient := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			customClientHits.Add(1)
			if r.URL.Path != "/model/claude-3-sonnet/invoke" {
				t.Errorf("Expected Bedrock path %q, got %q", "/model/claude-3-sonnet/invoke", r.URL.Path)
			}
			return &http.Response{
				StatusCode: 200,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(`{"id":"msg_test","type":"message","role":"assistant","content":[{"type":"text","text":"hi"}],"model":"claude-3-sonnet","stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`)),
				Request:    r,
			}, nil
		}),
	}

	var fallbackHits atomic.Int32
	fallbackServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fallbackHits.Add(1)
		writeMessagesResponse(w)
	}))
	t.Cleanup(fallbackServer.Close)

	cfg := makeStaticAWSConfig("us-east-1")
	cfg.HTTPClient = customClient
	client := anthropic.NewClient(
		option.WithoutEnvironmentDefaults(),
		WithConfig(cfg),
		option.WithBaseURL(fallbackServer.URL),
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

	if got := customClientHits.Load(); got != 1 {
		t.Fatalf("Expected AWS config HTTP client to handle 1 request, got %d", got)
	}
	if got := fallbackHits.Load(); got != 0 {
		t.Fatalf("Expected fallback server to receive no requests, got %d", got)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// --- EventStream → SSE response normalization tests ---

// encodeChunkFrame writes an EventStream "chunk" event frame whose payload
// carries the given Anthropic event JSON, the way Bedrock streams responses.
func encodeChunkFrame(t *testing.T, w io.Writer, eventJSON string) {
	t.Helper()
	payload, err := json.Marshal(eventstreamChunk{Bytes: base64.StdEncoding.EncodeToString([]byte(eventJSON))})
	if err != nil {
		t.Fatalf("Failed to marshal chunk payload: %v", err)
	}
	msg := eventstream.Message{Payload: payload}
	msg.Headers.Set(eventstreamapi.MessageTypeHeader, eventstream.StringValue(eventstreamapi.EventMessageType))
	msg.Headers.Set(eventstreamapi.EventTypeHeader, eventstream.StringValue("chunk"))
	if err := eventstream.NewEncoder().Encode(w, msg); err != nil {
		t.Fatalf("Failed to encode event frame: %v", err)
	}
}

// encodeExceptionFrame writes an EventStream exception frame, the way Bedrock
// reports mid-stream errors such as throttling.
func encodeExceptionFrame(t *testing.T, w io.Writer, exceptionType, message string) {
	t.Helper()
	payload, err := json.Marshal(map[string]string{"message": message})
	if err != nil {
		t.Fatalf("Failed to marshal exception payload: %v", err)
	}
	msg := eventstream.Message{Payload: payload}
	msg.Headers.Set(eventstreamapi.MessageTypeHeader, eventstream.StringValue(eventstreamapi.ExceptionMessageType))
	msg.Headers.Set(eventstreamapi.ExceptionTypeHeader, eventstream.StringValue(exceptionType))
	if err := eventstream.NewEncoder().Encode(w, msg); err != nil {
		t.Fatalf("Failed to encode exception frame: %v", err)
	}
}

// applyStreamingMiddleware runs bedrockMiddleware over a fake streaming
// request, with the wire responding with the given EventStream body, and
// returns the response the middleware produced.
func applyStreamingMiddleware(t *testing.T, frames *bytes.Buffer) *http.Response {
	t.Helper()
	middleware := bedrockMiddleware(v4.NewSigner(), makeStaticAWSConfig("us-east-1"))

	body := `{"model": "claude-3-sonnet", "stream": true, "messages": [{"role": "user", "content": "Hello"}]}`
	req, err := http.NewRequest("POST", "https://bedrock-runtime.us-east-1.amazonaws.com/v1/messages", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	res, err := middleware(req, func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Header:     http.Header{"Content-Type": []string{"application/vnd.amazon.eventstream"}},
			Body:       io.NopCloser(frames),
		}, nil
	})
	if err != nil {
		t.Fatalf("Middleware failed: %v", err)
	}
	return res
}

func TestBedrockStreamingResponseNormalizedToSSE(t *testing.T) {
	messageStartJSON := `{"type":"message_start","message":{"id":"msg_test"}}`
	deltaJSON := `{"type":"content_block_delta","delta":{"type":"text_delta","text":"Hi"}}`
	frames := &bytes.Buffer{}
	encodeChunkFrame(t, frames, messageStartJSON)
	encodeChunkFrame(t, frames, deltaJSON)

	res := applyStreamingMiddleware(t, frames)

	if got := res.Header.Get("Content-Type"); got != "text/event-stream" {
		t.Errorf("Expected Content-Type %q, got %q", "text/event-stream", got)
	}
	sse, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read normalized body: %v", err)
	}
	expected := "event: message_start\ndata: " + messageStartJSON + "\n\n" +
		"event: content_block_delta\ndata: " + deltaJSON + "\n\n"
	if string(sse) != expected {
		t.Errorf("Expected SSE body %q, got %q", expected, string(sse))
	}
}

func TestBedrockStreamingExceptionSurfacesAsBodyError(t *testing.T) {
	messageStartJSON := `{"type":"message_start","message":{"id":"msg_test"}}`
	frames := &bytes.Buffer{}
	encodeChunkFrame(t, frames, messageStartJSON)
	encodeExceptionFrame(t, frames, "ThrottlingException", "Too many requests")

	res := applyStreamingMiddleware(t, frames)

	sse, err := io.ReadAll(res.Body)
	if err == nil {
		t.Fatal("Expected an error reading a stream containing an exception frame")
	}
	expectedErr := "received exception ThrottlingException: Too many requests"
	if err.Error() != expectedErr {
		t.Errorf("Expected error %q, got %q", expectedErr, err.Error())
	}
	// Events decoded before the exception must still be delivered.
	expectedSSE := "event: message_start\ndata: " + messageStartJSON + "\n\n"
	if string(sse) != expectedSSE {
		t.Errorf("Expected SSE body %q before the error, got %q", expectedSSE, string(sse))
	}
}

// --- Middleware ordering tests ---

// TestBedrockUserMiddlewareObservesAnthropicShape verifies the documented
// ordering: middleware registered before the Bedrock option observes the
// Anthropic-shaped, unsigned request, while the wire receives the rewritten,
// signed Bedrock request.
func TestBedrockUserMiddlewareObservesAnthropicShape(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")

	var wirePath, wireAuth string
	var wireBody map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wirePath = r.URL.Path
		wireAuth = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&wireBody); err != nil {
			t.Errorf("Failed to decode wire body: %v", err)
		}
		writeMessagesResponse(w)
	}))
	t.Cleanup(server.Close)

	var observedPath, observedAuth string
	var observedBody map[string]any
	spy := func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
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

	client := anthropic.NewClient(
		option.WithoutEnvironmentDefaults(),
		option.WithMiddleware(spy),
		WithConfig(makeStaticAWSConfig("us-east-1")),
		option.WithBaseURL(server.URL),
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

	// The spy (outside the Bedrock adaptation) sees the Anthropic shape.
	if observedPath != "/v1/messages" {
		t.Errorf("Expected middleware to observe path %q, got %q", "/v1/messages", observedPath)
	}
	if observedBody["model"] != "claude-3-sonnet" {
		t.Errorf("Expected middleware to observe model in body, got %v", observedBody["model"])
	}
	if observedAuth != "" {
		t.Errorf("Expected middleware to observe no Authorization header, got %q", observedAuth)
	}

	// The wire sees the rewritten, signed Bedrock shape.
	if wirePath != "/model/claude-3-sonnet/invoke" {
		t.Errorf("Expected wire path %q, got %q", "/model/claude-3-sonnet/invoke", wirePath)
	}
	if _, ok := wireBody["model"]; ok {
		t.Error("Expected model to be removed from the wire body")
	}
	if wireBody["anthropic_version"] != DefaultVersion {
		t.Errorf("Expected anthropic_version %q on the wire, got %v", DefaultVersion, wireBody["anthropic_version"])
	}
	if !strings.HasPrefix(wireAuth, "AWS4-HMAC-SHA256") {
		t.Errorf("Expected SigV4 Authorization on the wire, got %q", wireAuth)
	}
}

// TestBedrockStreamingEndToEnd verifies that an EventStream wire response
// decodes into the same stream events a first-party SSE response would.
func TestBedrockStreamingEndToEnd(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")

	eventJSONs := []string{
		`{"type":"message_start","message":{"id":"msg_test","type":"message","role":"assistant","content":[],"model":"claude-3-sonnet","usage":{"input_tokens":1,"output_tokens":1}}}`,
		`{"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`,
		`{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hi"}}`,
		`{"type":"content_block_stop","index":0}`,
		`{"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":2}}`,
		`{"type":"message_stop"}`,
	}
	frames := &bytes.Buffer{}
	for _, eventJSON := range eventJSONs {
		encodeChunkFrame(t, frames, eventJSON)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.amazon.eventstream")
		w.Write(frames.Bytes())
	}))
	t.Cleanup(server.Close)

	client := anthropic.NewClient(
		option.WithoutEnvironmentDefaults(),
		WithConfig(makeStaticAWSConfig("us-east-1")),
		option.WithBaseURL(server.URL),
	)

	stream := client.Messages.NewStreaming(context.Background(), anthropic.MessageNewParams{
		Model:     "claude-3-sonnet",
		MaxTokens: 1,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("hi")),
		},
	})

	var gotTypes []string
	for stream.Next() {
		gotTypes = append(gotTypes, string(stream.Current().Type))
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("Expected no stream error, got: %v", err)
	}

	expectedTypes := []string{
		"message_start", "content_block_start", "content_block_delta",
		"content_block_stop", "message_delta", "message_stop",
	}
	if len(gotTypes) != len(expectedTypes) {
		t.Fatalf("Expected %d events %v, got %d: %v", len(expectedTypes), expectedTypes, len(gotTypes), gotTypes)
	}
	for i, expected := range expectedTypes {
		if gotTypes[i] != expected {
			t.Errorf("Expected event %d to be %q, got %q", i, expected, gotTypes[i])
		}
	}
}
