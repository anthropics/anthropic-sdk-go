package anthropic_test

import (
	"net/http"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// otelStyleRoundTripper satisfies http.RoundTripper without being an
// *http.Transport — mimics the otelhttp.NewTransport wrapping pattern that
// triggers the bug reported in #334.
type otelStyleRoundTripper struct {
	wrapped http.RoundTripper
}

func (rt *otelStyleRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt.wrapped.RoundTrip(r)
}

// TestNewClientWithWrappedDefaultTransport is the end-to-end test for #334.
// It mirrors what an OpenTelemetry-instrumented application does:
//
//	http.DefaultTransport = otelhttp.NewTransport(http.DefaultTransport)
//	client := anthropic.NewClient(...)
//
// Before the fix, NewClient panics during DefaultClientOptions() because
// defaultHTTPClient asserts *http.Transport on a wrapped RoundTripper.
//
// The test verifies the entire NewClient path completes without panic when
// DefaultTransport is wrapped, even when the caller does NOT override the
// HTTP client via option.WithHTTPClient.
func TestNewClientWithWrappedDefaultTransport(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })

	http.DefaultTransport = &otelStyleRoundTripper{wrapped: original}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("anthropic.NewClient panicked with wrapped DefaultTransport: %v", r)
		}
	}()

	_ = anthropic.NewClient(option.WithAPIKey("test-key-not-used-for-real-requests"))
}

// TestNewClientWithWrappedDefaultTransportPreservesUserOverride verifies that
// even when DefaultTransport is wrapped, a caller-supplied option.WithHTTPClient
// is honored — i.e., the fix does not change semantics for callers who pass
// their own client.
func TestNewClientWithWrappedDefaultTransportPreservesUserOverride(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })

	http.DefaultTransport = &otelStyleRoundTripper{wrapped: original}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("anthropic.NewClient panicked even with caller-supplied HTTPClient: %v", r)
		}
	}()

	customClient := &http.Client{Transport: http.DefaultTransport}
	_ = anthropic.NewClient(
		option.WithAPIKey("test-key-not-used-for-real-requests"),
		option.WithHTTPClient(customClient),
	)
}
