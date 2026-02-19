package testutil

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
)

// NewVCRHTTPClient creates an *http.Client wired to a go-vcr recorder.
// Cassette files are stored under testdata/cassettes.
// If ANTHROPIC_LIVE=1, the recorder runs in recording mode; otherwise replay-only.
func NewVCRHTTPClient(t *testing.T, cassetteName string) (*http.Client, *recorder.Recorder) {
	t.Helper()

	mode := recorder.ModeReplaying
	if os.Getenv("ANTHROPIC_LIVE") == "1" {
		mode = recorder.ModeRecording
	}

	// Let go-vcr handle the .yaml extension to avoid accidental double suffixes
	cassettePath := filepath.Join("testdata", "cassettes", cassetteName)
	r, err := recorder.NewAsMode(cassettePath, mode, nil)
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}

	t.Cleanup(func() {
		_ = r.Stop()
	})

	httpClient := &http.Client{Transport: r}
	return httpClient, r
}
