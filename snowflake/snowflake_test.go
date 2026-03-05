package snowflake

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestCortexURLRewriting(t *testing.T) {
	testCases := []struct {
		name         string
		path         string
		method       string
		expectedPath string
	}{
		{
			name:         "messages endpoint",
			path:         "/v1/messages",
			method:       http.MethodPost,
			expectedPath: "/api/v2/cortex/inference:complete",
		},
		{
			name:         "count tokens endpoint",
			path:         "/v1/messages/count_tokens",
			method:       http.MethodPost,
			expectedPath: "/api/v2/cortex/inference:complete",
		},
		{
			name:         "non-messages endpoint unchanged",
			path:         "/v1/other",
			method:       http.MethodPost,
			expectedPath: "/v1/other",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			middleware := cortexMiddleware("test-token")

			requestBody := map[string]any{
				"model":  "claude-3-5-sonnet",
				"stream": false,
				"messages": []map[string]string{
					{"role": "user", "content": "Hello"},
				},
			}

			bodyBytes, err := json.Marshal(requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req, err := http.NewRequest(tc.method, "https://myorg-myaccount.snowflakecomputing.com"+tc.path, bytes.NewReader(bodyBytes))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			_, err = middleware(req, func(r *http.Request) (*http.Response, error) {
				if r.URL.Path != tc.expectedPath {
					t.Errorf("Expected Path %q, got %q", tc.expectedPath, r.URL.Path)
				}

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

func TestCortexAuthorizationHeader(t *testing.T) {
	token := "test-jwt-token"
	middleware := cortexMiddleware(token)

	requestBody := map[string]any{
		"model": "claude-3-5-sonnet",
		"messages": []map[string]string{
			{"role": "user", "content": "Hello"},
		},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://myorg-myaccount.snowflakecomputing.com/v1/messages", bytes.NewReader(bodyBytes))
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

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type %q, got %q", "application/json", contentType)
		}

		accept := r.Header.Get("Accept")
		if accept != "application/json, text/event-stream" {
			t.Errorf("Expected Accept %q, got %q", "application/json, text/event-stream", accept)
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

func TestCortexAnthropicVersionInjected(t *testing.T) {
	middleware := cortexMiddleware("test-token")

	requestBody := map[string]any{
		"model": "claude-3-5-sonnet",
		"messages": []map[string]string{
			{"role": "user", "content": "Hello"},
		},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://myorg-myaccount.snowflakecomputing.com/v1/messages", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = middleware(req, func(r *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read body: %v", err)
		}

		var parsed map[string]any
		if err := json.Unmarshal(body, &parsed); err != nil {
			t.Fatalf("Failed to unmarshal body: %v", err)
		}

		version, ok := parsed["anthropic_version"]
		if !ok {
			t.Fatal("Expected anthropic_version to be set in body")
		}
		if version != DefaultVersion {
			t.Errorf("Expected anthropic_version %q, got %q", DefaultVersion, version)
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

func TestCortexAnthropicVersionNotOverwritten(t *testing.T) {
	middleware := cortexMiddleware("test-token")

	customVersion := "custom-version-2024-01-01"
	requestBody := map[string]any{
		"model":             "claude-3-5-sonnet",
		"anthropic_version": customVersion,
		"messages": []map[string]string{
			{"role": "user", "content": "Hello"},
		},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://myorg-myaccount.snowflakecomputing.com/v1/messages", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = middleware(req, func(r *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read body: %v", err)
		}

		var parsed map[string]any
		if err := json.Unmarshal(body, &parsed); err != nil {
			t.Fatalf("Failed to unmarshal body: %v", err)
		}

		version := parsed["anthropic_version"]
		if version != customVersion {
			t.Errorf("Expected anthropic_version %q to not be overwritten, got %q", customVersion, version)
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

func TestCortexStreamingRequest(t *testing.T) {
	middleware := cortexMiddleware("test-token")

	requestBody := map[string]any{
		"model":  "claude-3-5-sonnet",
		"stream": true,
		"messages": []map[string]string{
			{"role": "user", "content": "Hello"},
		},
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://myorg-myaccount.snowflakecomputing.com/v1/messages", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = middleware(req, func(r *http.Request) (*http.Response, error) {
		// Verify the URL is correctly rewritten
		if r.URL.Path != "/api/v2/cortex/inference:complete" {
			t.Errorf("Expected path /api/v2/cortex/inference:complete, got %s", r.URL.Path)
		}

		// Verify the stream field is preserved (Cortex handles streaming via SSE natively)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read body: %v", err)
		}

		var parsed map[string]any
		if err := json.Unmarshal(body, &parsed); err != nil {
			t.Fatalf("Failed to unmarshal body: %v", err)
		}

		stream, ok := parsed["stream"]
		if !ok {
			t.Fatal("Expected stream field to be preserved in body")
		}
		if stream != true {
			t.Errorf("Expected stream to be true, got %v", stream)
		}

		// Verify model is preserved (unlike Bedrock, Cortex keeps model in body)
		model, ok := parsed["model"]
		if !ok {
			t.Fatal("Expected model field to be preserved in body")
		}
		if model != "claude-3-5-sonnet" {
			t.Errorf("Expected model to be claude-3-5-sonnet, got %v", model)
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

func TestWithAccountPanicsOnEmptyAccount(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for empty account")
		}
	}()
	WithAccount("", "token")
}

func TestWithAccountPanicsOnEmptyToken(t *testing.T) {
	// Ensure the env var is not set
	os.Unsetenv("SNOWFLAKE_AUTH_TOKEN")

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for empty token without env var")
		}
	}()
	WithAccount("myaccount", "")
}

func TestWithAccountUsesEnvVar(t *testing.T) {
	os.Setenv("SNOWFLAKE_AUTH_TOKEN", "env-token")
	defer os.Unsetenv("SNOWFLAKE_AUTH_TOKEN")

	// Should not panic when env var is set
	opt := WithAccount("myaccount", "")
	if opt == nil {
		t.Error("Expected a non-nil option")
	}
}

func TestWithAccountStripsSuffix(t *testing.T) {
	// Verify it doesn't panic and works with the suffix
	opt := WithAccount("myaccount.snowflakecomputing.com", "test-token")
	if opt == nil {
		t.Error("Expected a non-nil option")
	}
}

func TestWithBaseURLPanicsOnEmptyURL(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for empty base URL")
		}
	}()
	WithBaseURL("", "token")
}

func TestWithBaseURLPanicsOnEmptyToken(t *testing.T) {
	os.Unsetenv("SNOWFLAKE_AUTH_TOKEN")

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for empty token without env var")
		}
	}()
	WithBaseURL("https://example.com", "")
}
