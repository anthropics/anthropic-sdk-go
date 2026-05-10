package anthropic

import (
	"net/http"
	"testing"
)

// roundTripperFunc satisfies http.RoundTripper without being an *http.Transport.
// Used to simulate the otelhttp wrapping case from #334.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// TestDefaultHTTPClientWithWrappedTransport verifies that defaultHTTPClient
// does not panic when http.DefaultTransport has been replaced with a
// RoundTripper that is not an *http.Transport. This is the bug reported in
// #334 (otelhttp.NewTransport wrapping).
func TestDefaultHTTPClientWithWrappedTransport(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })

	http.DefaultTransport = roundTripperFunc(func(*http.Request) (*http.Response, error) {
		return nil, nil
	})

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("defaultHTTPClient panicked when DefaultTransport is not *http.Transport: %v", r)
		}
	}()

	client := defaultHTTPClient()
	if client == nil {
		t.Fatal("defaultHTTPClient returned nil")
	}
	if client.Transport == nil {
		t.Fatal("defaultHTTPClient returned client with nil Transport")
	}
}

// TestDefaultHTTPClientWithStdlibTransport verifies the common path — when
// http.DefaultTransport is the stdlib *http.Transport, defaultHTTPClient
// returns a clone with ResponseHeaderTimeout configured.
func TestDefaultHTTPClientWithStdlibTransport(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })

	http.DefaultTransport = &http.Transport{}

	client := defaultHTTPClient()
	if client == nil {
		t.Fatal("defaultHTTPClient returned nil")
	}

	transport, ok := client.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("expected client.Transport to be *http.Transport, got %T", client.Transport)
	}
	if transport.ResponseHeaderTimeout != defaultResponseHeaderTimeout {
		t.Fatalf("ResponseHeaderTimeout: got %v, want %v",
			transport.ResponseHeaderTimeout, defaultResponseHeaderTimeout)
	}
}
