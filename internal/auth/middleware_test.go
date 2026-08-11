package auth

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
)

func configFor(t *testing.T, base string) *requestconfig.RequestConfig {
	t.Helper()
	u, err := url.Parse(base)
	if err != nil {
		t.Fatal(err)
	}
	return &requestconfig.RequestConfig{BaseURL: u}
}

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
	mw := NewMiddleware(cache, configFor(t, server.URL))

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
	mw := NewMiddleware(cache, configFor(t, server.URL))

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

func TestMiddleware_ExchangesAtConfiguredBaseURL(t *testing.T) {
	var providerBaseURLs []string
	provider := func(_ context.Context, baseURL string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		providerBaseURLs = append(providerBaseURLs, baseURL)
		return &AccessToken{Token: "tok-1"}, nil
	}
	next := func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: http.NoBody, Request: req}, nil
	}
	mw := NewMiddleware(NewTokenCache(provider, nil), configFor(t, "https://api.example.test/prefix/"))

	req, _ := http.NewRequest(http.MethodGet, "https://API.example.test/prefix/v1/models", nil)
	if _, err := mw(req, next); err != nil {
		t.Fatal(err)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer tok-1" {
		t.Fatalf("same-origin request got Authorization %q, want bearer", got)
	}
	if len(providerBaseURLs) != 1 || providerBaseURLs[0] != "https://api.example.test" {
		t.Fatalf("provider saw base URLs %q, want just the configured origin", providerBaseURLs)
	}

	for _, foreign := range []string{
		"https://attacker.example.test/v1/models",
		"https://api.example.test:8443/v1/models",
		"http://api.example.test/v1/models",
	} {
		req, _ := http.NewRequest(http.MethodGet, foreign, nil)
		if _, err := mw(req, next); err != nil {
			t.Fatalf("%s: %v", foreign, err)
		}
		if got := req.Header.Get("Authorization"); got != "" {
			t.Errorf("%s: Authorization %q attached to a foreign origin", foreign, got)
		}
		if got := req.Header.Get("anthropic-beta"); got != "" {
			t.Errorf("%s: anthropic-beta %q attached to a foreign origin", foreign, got)
		}
	}
	if len(providerBaseURLs) != 1 {
		t.Fatalf("provider was re-invoked for foreign origins: %q", providerBaseURLs)
	}
}

func TestMiddleware_UnresolvedBaseURLFails(t *testing.T) {
	provider := func(_ context.Context, _ string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		t.Fatal("provider must not be invoked without a configured base URL")
		return nil, nil
	}
	mw := NewMiddleware(NewTokenCache(provider, nil), &requestconfig.RequestConfig{})
	req, _ := http.NewRequest(http.MethodGet, "https://api.example.test/v1/models", nil)
	_, err := mw(req, func(*http.Request) (*http.Response, error) {
		t.Fatal("request must not be sent")
		return nil, nil
	})
	if err == nil || !strings.Contains(err.Error(), "base URL is not set") {
		t.Fatalf("got %v, want base URL error", err)
	}
}
