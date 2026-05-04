package auth

import (
	"context"
	"net/http"
	"time"
)

// AccessToken represents an access token with an optional expiry time.
type AccessToken struct {
	Token     string
	ExpiresAt *time.Time // nil means no expiry (cache forever)
}

// TokenProvider fetches a fresh [AccessToken], potentially performing network calls.
// The handler parameter is the resolved HTTP transport from [requestconfig.RequestConfig]
// (CustomHTTPDoer ?? HTTPClient).
//
// Implementations must be safe for concurrent use because [TokenCache] may
// invoke the provider from a background goroutine while another call is
// in-flight on the main request path.
type TokenProvider func(ctx context.Context, baseURL string, handler func(*http.Request) (*http.Response, error)) (*AccessToken, error)

// IdentityTokenProvider provides a JWT identity token for federation.
type IdentityTokenProvider interface {
	GetIdentityToken(ctx context.Context) (string, error)
}

// CredentialsResult bundles a [TokenProvider] with config-level metadata that
// needs to be propagated to the client (e.g. base_url, workspace_id).
//
// OrganizationID is intentionally absent: the server exposes the caller's
// organization only as the anthropic-organization-id *response* header, not
// as a request header, so there is no way for the SDK to act on a
// config-level organization ID today.
type CredentialsResult struct {
	Provider    TokenProvider
	BaseURL     string // from config file; empty if not set
	WorkspaceID string
}
