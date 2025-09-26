package bedrock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/credentials"
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
