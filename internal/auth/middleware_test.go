package auth

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestMiddleware_401WithUnreplayableBodyWarns verifies that the 401-retry path
// emits a one-shot warning when it has to give up because the request body
// cannot be replayed (no GetBody). Today this returns the 401 silently.
func TestMiddleware_401WithUnreplayableBodyWarns(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	provider := func(_ context.Context, _ string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		return &AccessToken{Token: "tok-1"}, nil
	}
	cache := NewTokenCache(provider, http.DefaultClient.Do)
	mw := authMiddleware(cache)

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, server.URL, io.NopCloser(strings.NewReader("body")))
	req.GetBody = nil // unreplayable

	resetWarnOnceForTest(t)
	var buf syncBuffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	resp, err := mw(req, http.DefaultClient.Do)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	if !strings.Contains(buf.String(), "cannot be replayed") {
		t.Fatalf("expected warning about unreplayable body, got: %q", buf.String())
	}
}

// TestMiddleware_401ForceTokenErrorWarns verifies that when ForceToken itself
// returns an error during the 401 retry, the middleware logs the underlying
// cause instead of silently returning the original 401.
func TestMiddleware_401ForceTokenErrorWarns(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	var callCount int
	forceErr := errors.New("token endpoint unreachable at http://elsewhere")
	provider := func(_ context.Context, _ string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		callCount++
		if callCount == 1 {
			return &AccessToken{Token: "tok-1"}, nil
		}
		return nil, forceErr
	}
	cache := NewTokenCache(provider, http.DefaultClient.Do)
	mw := authMiddleware(cache)

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, server.URL, bytes.NewReader([]byte("body")))
	req.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader([]byte("body"))), nil }

	resetWarnOnceForTest(t)
	var buf syncBuffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)

	resp, err := mw(req, http.DefaultClient.Do)
	if err != nil {
		t.Fatalf("expected nil error from middleware so caller gets the original 401, got %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
	resp.Body.Close()

	out := buf.String()
	if !strings.Contains(out, "token endpoint unreachable") {
		t.Fatalf("expected warning to include underlying error, got: %q", out)
	}
}
