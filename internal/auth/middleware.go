package auth

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
)

// applyBearerAuth sets the Authorization and anthropic-beta headers for
// OAuth/bearer-based credentials on outgoing API requests.
func applyBearerAuth(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
	existing := req.Header.Get("anthropic-beta")
	switch {
	case existing == "":
		req.Header.Set("anthropic-beta", OAuthAPIBetaHeader)
	case !strings.Contains(existing, OAuthAPIBetaHeader):
		req.Header.Set("anthropic-beta", existing+","+OAuthAPIBetaHeader)
	}
}

type authMiddlewareFunc func(*http.Request, func(*http.Request) (*http.Response, error)) (*http.Response, error)

// Middleware is the bearer-auth HTTP middleware function shape. Exported so
// callers that construct the middleware lazily inside another middleware
// (e.g. option.WithConfig's deferred resolve) can hold the value.
type Middleware = authMiddlewareFunc

// NewProviderMiddleware constructs a bearer-auth middleware bound to provider
// and the given HTTP handler. Prefer [WithAuthMiddleware] when auth can be
// installed at option-Apply time; use this only when the [TokenProvider]
// isn't known until request time.
func NewProviderMiddleware(provider TokenProvider, handler func(*http.Request) (*http.Response, error)) Middleware {
	return authMiddleware(NewTokenCache(provider, handler))
}

func authMiddleware(cache *TokenCache) authMiddlewareFunc {
	return func(req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error) {
		// Skip if static auth headers are already set.
		if req.Header.Get("X-Api-Key") != "" || req.Header.Get("Authorization") != "" {
			return next(req)
		}

		baseURL := req.URL.Scheme + "://" + req.URL.Host
		token, err := cache.Token(req.Context(), baseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to get credentials token: %w", err)
		}
		applyBearerAuth(req, token)

		resp, err := next(req)
		if err != nil {
			return resp, err
		}

		if resp.StatusCode != http.StatusUnauthorized {
			return resp, nil
		}

		if !cache.Invalidate() {
			return resp, nil
		}

		if req.GetBody == nil && req.Body != nil {
			warnOnce(
				"401-retry-unreplayable-body",
				"401 response; request body cannot be replayed (no GetBody), skipping credential refresh retry. Use a buffered body or fix credentials.",
			)
			return resp, nil
		}

		token, tokenErr := cache.ForceToken(req.Context(), baseURL)
		if tokenErr != nil {
			warnOnce(
				"401-retry-refresh-failed",
				"401 response; credential refresh also failed, returning original 401: %v",
				tokenErr,
			)
			return resp, nil
		}
		applyBearerAuth(req, token)

		if req.GetBody != nil {
			newBody, bodyErr := req.GetBody()
			if bodyErr != nil {
				resp.Body.Close()
				return nil, bodyErr
			}
			req.Body = newBody
		}

		resp.Body.Close()
		return next(req)
	}
}

// clientKey identifies the HTTP transport a [TokenCache] should bind to.
// Keying the per-option cache map on this struct ensures two clients that
// share the same [WithAuthMiddleware] option value each get their own cache
// bound to their own transport, rather than the second client silently
// reusing the first client's transport for token exchange.
type clientKey struct {
	httpClient *http.Client
	customDoer requestconfig.HTTPDoer
}

// WithAuthMiddleware returns a [requestconfig.RequestOptionFunc] that appends
// HTTP middleware which authenticates requests using a cached bearer token.
//
// A separate [TokenCache] is lazily constructed per distinct HTTP transport
// (the combination of [requestconfig.RequestConfig.HTTPClient] and
// [requestconfig.RequestConfig.CustomHTTPDoer]) that applies this option, so
// the same option value can safely be shared across multiple clients built
// with different transports. Within one client, the cache is constructed
// exactly once and reused for every request so its in-memory token, refresh
// state, and 401 invalidation signal are preserved.
//
// The base URL for token-exchange / refresh HTTP calls is derived from the
// outgoing request URL at the moment the middleware runs, so this option can
// be listed in any order relative to [option.WithBaseURL] /
// [option.WithEnvironmentProduction].
//
// If the request already has an X-Api-Key or Authorization header set
// (e.g. via [option.WithAPIKey] or [option.WithAuthToken]), the middleware
// is a no-op pass-through.
func WithAuthMiddleware(provider TokenProvider) requestconfig.RequestOptionFunc {
	var (
		mu    sync.Mutex
		byKey = map[clientKey]authMiddlewareFunc{}
	)
	return func(r *requestconfig.RequestConfig) error {
		key := clientKey{httpClient: r.HTTPClient, customDoer: r.CustomHTTPDoer}

		mu.Lock()
		mw, ok := byKey[key]
		if !ok {
			handler := r.HTTPClient.Do
			if r.CustomHTTPDoer != nil {
				handler = r.CustomHTTPDoer.Do
			}
			cache := NewTokenCache(provider, handler)
			mw = authMiddleware(cache)
			byKey[key] = mw
		}
		mu.Unlock()

		r.Middlewares = append(r.Middlewares, mw)
		return nil
	}
}
