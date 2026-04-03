package vertex

import (
	"context"
	"testing"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
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
			name:        "europe region",
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
