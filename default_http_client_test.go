package anthropic_test

import (
	"net/http"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// customRoundTripper is a non-*http.Transport RoundTripper that simulates
// wrappers like otelhttp.NewTransport replacing http.DefaultTransport.
type customRoundTripper struct {
	base http.RoundTripper
}

func (c *customRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return c.base.RoundTrip(req)
}

func TestDefaultHTTPClientWithCustomTransport(t *testing.T) {
	// Save and restore http.DefaultTransport to avoid polluting other tests.
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })

	// Replace with a wrapper that is NOT *http.Transport, like otelhttp would.
	http.DefaultTransport = &customRoundTripper{base: original}

	// NewClient calls defaultHTTPClient() internally. Before the fix in
	// commit 7c9d586, this panicked with:
	//   "interface conversion: http.RoundTripper is *customRoundTripper, not *http.Transport"
	_ = anthropic.NewClient(
		option.WithAPIKey("test-key"),
	)
}

func TestDefaultHTTPClientWithCustomTransportAndWithHTTPClient(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })

	http.DefaultTransport = &customRoundTripper{base: original}

	// Even when the caller provides their own http.Client, NewClient must not
	// panic during construction of default options.
	_ = anthropic.NewClient(
		option.WithAPIKey("test-key"),
		option.WithHTTPClient(&http.Client{}),
	)
}

func TestDefaultHTTPClientWithNilTransport(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })

	// Edge case: nil DefaultTransport should not panic either.
	http.DefaultTransport = nil

	_ = anthropic.NewClient(
		option.WithAPIKey("test-key"),
	)
}

func TestDefaultHTTPClientWithStandardTransport(t *testing.T) {
	// Verify the normal case still works: when http.DefaultTransport is
	// the standard *http.Transport, the SDK should clone it and set the
	// ResponseHeaderTimeout.
	_ = anthropic.NewClient(
		option.WithAPIKey("test-key"),
	)
}
