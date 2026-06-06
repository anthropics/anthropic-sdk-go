package vertex

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

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
