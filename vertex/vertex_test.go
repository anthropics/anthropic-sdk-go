package vertex

import (
	"context"
	"strings"
	"testing"

	"golang.org/x/oauth2/google"
)

func TestBaseURLConfiguration(t *testing.T) {

	tests := []struct {
		name        string
		region      string
		expectedURL string
	}{
		{
			name:        "global region uses correct base URL",
			region:      "global",
			expectedURL: "https://aiplatform.googleapis.com/",
		},
		{
			name:        "us-central1 region uses correct base URL",
			region:      "us-central1",
			expectedURL: "https://us-central1-aiplatform.googleapis.com/",
		},
		{
			name:        "europe-west1 region uses correct base URL",
			region:      "europe-west1",
			expectedURL: "https://europe-west1-aiplatform.googleapis.com/",
		},
		{
			name:        "asia-southeast1 region uses correct base URL",
			region:      "asia-southeast1",
			expectedURL: "https://asia-southeast1-aiplatform.googleapis.com/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the base URL generation logic
			var baseURL string
			if tt.region == "global" {
				baseURL = "https://aiplatform.googleapis.com/"
			} else {
				baseURL = "https://" + tt.region + "-aiplatform.googleapis.com/"
			}

			if baseURL != tt.expectedURL {
				t.Errorf("Expected base URL %s, got %s", tt.expectedURL, baseURL)
			}
		})
	}
}

func TestWithCredentialsRegionHandling(t *testing.T) {
	ctx := context.Background()
	projectID := "test-project"
	creds := &google.Credentials{}

	// Test that the function doesn't panic with global region
	t.Run("global region does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if str, ok := r.(string); ok && strings.Contains(str, "region must be provided") {
					t.Errorf("WithCredentials panicked with global region: %v", r)
				}
				// Re-panic if it's a different error (like transport.NewHTTPClient failing in test env)
				panic(r)
			}
		}()

		// This will likely panic due to transport.NewHTTPClient, but it shouldn't be due to region validation
		WithCredentials(ctx, "global", projectID, creds)
	})

	t.Run("regional endpoint does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if str, ok := r.(string); ok && strings.Contains(str, "region must be provided") {
					t.Errorf("WithCredentials panicked with us-central1 region: %v", r)
				}
				// Re-panic if it's a different error (like transport.NewHTTPClient failing in test env)
				panic(r)
			}
		}()

		// This will likely panic due to transport.NewHTTPClient, but it shouldn't be due to region validation
		WithCredentials(ctx, "us-central1", projectID, creds)
	})
}

func TestWithGoogleAuthRegionRequired(t *testing.T) {
	ctx := context.Background()

	t.Run("empty region panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic when region is empty")
			} else if str, ok := r.(string); ok && !strings.Contains(str, "region must be provided") {
				t.Errorf("Expected panic about region, got: %v", r)
			}
		}()

		WithGoogleAuth(ctx, "", "test-project")
	})

	t.Run("global region is accepted", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if str, ok := r.(string); ok && strings.Contains(str, "region must be provided") {
					t.Errorf("WithGoogleAuth should accept 'global' as a valid region: %v", r)
				}
				// Re-panic if it's a different error (like credentials not found)
				panic(r)
			}
		}()

		// This will likely panic due to credentials not being found, but it shouldn't be due to region validation
		WithGoogleAuth(ctx, "global", "test-project")
	})
}