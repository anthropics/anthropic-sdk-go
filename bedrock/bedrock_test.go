package bedrock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

func TestBedrockURLEncoding(t *testing.T) {
	testCases := []struct {
		name           string
		model          string
		stream         bool
		expectedPath   string
		expectedRawPath string
	}{
		{
			name:           "regular model name",
			model:          "claude-3-sonnet",
			stream:         false,
			expectedPath:   "/model/claude-3-sonnet/invoke",
			expectedRawPath: "/model/claude-3-sonnet/invoke",
		},
		{
			name:           "regular model name with streaming",
			model:          "claude-3-sonnet",
			stream:         true,
			expectedPath:   "/model/claude-3-sonnet/invoke-with-response-stream",
			expectedRawPath: "/model/claude-3-sonnet/invoke-with-response-stream",
		},
		{
			name:           "inference profile ARN with slashes",
			model:          "arn:aws:bedrock:us-east-1:947123456126:application-inference-profile/xv9example4b",
			stream:         false,
			expectedPath:   "/model/arn:aws:bedrock:us-east-1:947123456126:application-inference-profile/xv9example4b/invoke",
			expectedRawPath: "/model/arn%3Aaws%3Abedrock%3Aus-east-1%3A947123456126%3Aapplication-inference-profile%2Fxv9example4b/invoke",
		},
		{
			name:           "inference profile ARN with streaming",
			model:          "arn:aws:bedrock:us-east-1:947123456126:application-inference-profile/xv9example4b",
			stream:         true,
			expectedPath:   "/model/arn:aws:bedrock:us-east-1:947123456126:application-inference-profile/xv9example4b/invoke-with-response-stream",
			expectedRawPath: "/model/arn%3Aaws%3Abedrock%3Aus-east-1%3A947123456126%3Aapplication-inference-profile%2Fxv9example4b/invoke-with-response-stream",
		},
		{
			name:           "foundation model ARN with colons",
			model:          "arn:aws:bedrock:us-east-1:123456789012:foundation-model/anthropic.claude-3-sonnet-20240229-v1:0",
			stream:         false,
			expectedPath:   "/model/arn:aws:bedrock:us-east-1:123456789012:foundation-model/anthropic.claude-3-sonnet-20240229-v1:0/invoke",
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
			requestBody := map[string]interface{}{
				"model":   tc.model,
				"stream":  tc.stream,
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

func TestURLEncodingCorrectness(t *testing.T) {
	// Test that our approach correctly handles URL encoding edge cases
	testModel := "arn:aws:bedrock:us-east-1:947123456126:application-inference-profile/xv9example4b"

	// Create URL with both Path and RawPath (like the PR does)
	u := &url.URL{
		Scheme:  "https",
		Host:    "bedrock-runtime.us-east-1.amazonaws.com",
		Path:    fmt.Sprintf("/model/%s/invoke", testModel),
		RawPath: fmt.Sprintf("/model/%s/invoke", url.QueryEscape(testModel)),
	}

	// Verify that parsing the URL string gives us back the original model name
	urlString := u.String()
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		t.Fatalf("Failed to parse URL: %v", err)
	}

	// The parsed URL should have the original unencoded path
	expectedPath := fmt.Sprintf("/model/%s/invoke", testModel)
	if parsedURL.Path != expectedPath {
		t.Errorf("Expected parsed Path %q, got %q", expectedPath, parsedURL.Path)
	}

	// The RawPath should be encoded
	expectedRawPath := fmt.Sprintf("/model/%s/invoke", url.QueryEscape(testModel))
	if parsedURL.RawPath != expectedRawPath {
		t.Errorf("Expected parsed RawPath %q, got %q", expectedRawPath, parsedURL.RawPath)
	}
}

func TestBedrockMiddlewarePreservesOtherRequests(t *testing.T) {
	// Test that the middleware doesn't interfere with non-bedrock requests
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

	// Create a request to a different endpoint
	req, err := http.NewRequest("GET", "https://bedrock-runtime.us-east-1.amazonaws.com/some-other-endpoint", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	originalPath := req.URL.Path

	// Apply middleware
	_, err = middleware(req, func(r *http.Request) (*http.Response, error) {
		// The path should remain unchanged for non-default endpoints
		if r.URL.Path != originalPath {
			t.Errorf("Expected path to remain %q, got %q", originalPath, r.URL.Path)
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