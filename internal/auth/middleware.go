package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

// NewMiddleware returns bearer-auth middleware serving tokens from cache.
// Tokens are exchanged at, and only attached to requests for, cfg.BaseURL —
// a request path may be an absolute URL for some other host.
func NewMiddleware(cache *TokenCache, cfg *requestconfig.RequestConfig) authMiddlewareFunc {
	return func(req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error) {
		// Skip if static auth headers are already set.
		if req.Header.Get("X-Api-Key") != "" || req.Header.Get("Authorization") != "" {
			return next(req)
		}

		if cfg.BaseURL == nil {
			return nil, errors.New("cannot resolve credentials: client base URL is not set")
		}
		if !sameOrigin(req.URL, cfg.BaseURL) {
			return next(req)
		}
		baseURL := cfg.BaseURL.Scheme + "://" + cfg.BaseURL.Host

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

func sameOrigin(u, base *url.URL) bool {
	return strings.EqualFold(u.Scheme, base.Scheme) && strings.EqualFold(u.Host, base.Host)
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
// The base URL is read at request time, so this option can be listed in any
// order relative to [option.WithBaseURL] / [option.WithEnvironmentProduction].
//
// If the request already has an X-Api-Key or Authorization header set
// (e.g. via [option.WithAPIKey] or [option.WithAuthToken]), the middleware
// is a no-op pass-through.
func WithAuthMiddleware(provider TokenProvider) requestconfig.RequestOptionFunc {
	var (
		mu    sync.Mutex
		byKey = map[clientKey]*TokenCache{}
	)
	return func(r *requestconfig.RequestConfig) error {
		key := clientKey{httpClient: r.HTTPClient, customDoer: r.CustomHTTPDoer}

		mu.Lock()
		cache, ok := byKey[key]
		if !ok {
			handler := r.HTTPClient.Do
			if r.CustomHTTPDoer != nil {
				handler = r.CustomHTTPDoer.Do
			}
			cache = NewTokenCache(provider, handler)
			byKey[key] = cache
		}
		mu.Unlock()

		r.Middlewares = append(r.Middlewares, NewMiddleware(cache, r))
		return nil
	}
}
