package config

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// RequireSecureTokenEndpoint rejects base URLs that would cause a JWT
// assertion, authorization code or refresh token to be sent over cleartext
// HTTP. Local development hosts (localhost, 127.0.0.1, ::1) are allowed.
func RequireSecureTokenEndpoint(base string) error {
	if base == "" {
		return nil
	}
	u, err := url.Parse(base)
	if err != nil {
		return fmt.Errorf("invalid token endpoint base URL %q: %w", base, err)
	}
	if u.Scheme == "https" {
		return nil
	}
	if u.Scheme == "http" && isLoopbackHost(u.Hostname()) {
		return nil
	}
	return fmt.Errorf("refusing to send credential over non-https token endpoint %q", base)
}

// NewTokenHTTPClient returns the client to send token requests with when the
// caller has not supplied one.
func NewTokenHTTPClient() *http.Client {
	return WithoutRedirects(&http.Client{Timeout: 30 * time.Second})
}

// WithoutRedirects returns a copy of client that does not follow redirects,
// for sending token requests: a 307 or 308 would replay the body, which
// carries the credential, to a host [RequireSecureTokenEndpoint] never saw.
// The transport, timeout and cookie jar carry over, and client is not modified.
func WithoutRedirects(client *http.Client) *http.Client {
	c := *client
	c.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &c
}

func isLoopbackHost(host string) bool {
	host = strings.ToLower(host)
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}
