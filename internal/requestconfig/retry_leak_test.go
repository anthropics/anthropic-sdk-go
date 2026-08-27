package requestconfig

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

type capturingDoer struct {
	mu    sync.Mutex
	seen  []context.Context
	inner http.RoundTripper
}

func (d *capturingDoer) Do(req *http.Request) (*http.Response, error) {
	d.mu.Lock()
	d.seen = append(d.seen, req.Context())
	d.mu.Unlock()
	return d.inner.RoundTrip(req)
}

// TestRetryDoesNotLeakPerAttemptTimeouts verifies that retry timeout contexts are cancelled after each attempt.
func TestRetryDoesNotLeakPerAttemptTimeouts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-should-retry", "true")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	doer := &capturingDoer{inner: http.DefaultTransport}
	var raw *http.Response
	cfg := &RequestConfig{
		Context:        context.Background(),
		Request:        req,
		BaseURL:        mockBaseURL(server),
		HTTPClient:     http.DefaultClient,
		CustomHTTPDoer: doer,
		MaxRetries:     3,
		RequestTimeout: time.Hour,
		ResponseInto:   &raw,
	}
	if err := cfg.Execute(); err == nil {
		t.Fatal("Execute returned nil, want an error from repeated 429 responses")
	}

	if raw == nil || raw.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("final response was not the intended 429: %+v", raw)
	}
	if want := cfg.MaxRetries + 1; len(doer.seen) != want {
		t.Fatalf("expected %d attempts (initial + %d retries), got %d", want, cfg.MaxRetries, len(doer.seen))
	}

	for i, ctx := range doer.seen[:len(doer.seen)-1] {
		select {
		case <-ctx.Done():
		default:
			t.Fatalf("attempt %d: per-request timeout context was not cancelled after Execute returned (timer leaked)", i)
		}
	}
}
