package auth

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// advisoryRefreshBackoff caps how often a failed background refresh will
// re-run in the advisory window. Without this, a token endpoint that is
// down during the advisory window would be re-hit at request rate until
// the mandatory threshold is crossed.
const advisoryRefreshBackoff = 5 * time.Second

// TokenCache wraps a [TokenProvider] with two-tier proactive refresh and
// thread-safe caching.
//
// Refresh tiers:
//   - No cached token: fetch synchronously
//   - ExpiresAt == nil: cache forever
//   - More than [AdvisoryRefreshThreshold] remaining: return cached
//   - Between mandatory and advisory thresholds: try background refresh;
//     on failure, return the stale token
//   - Less than [MandatoryRefreshThreshold] remaining or expired: refresh
//     synchronously; on failure, return the error
//
// A cached token is scoped to the baseURL it was minted for: if a caller
// passes a different baseURL to a subsequent [TokenCache.Token] call, the
// cached token is discarded so a token issued by one environment can never
// be served for a request going to another.
type TokenCache struct {
	provider    TokenProvider
	httpHandler func(*http.Request) (*http.Response, error)

	mu                sync.Mutex
	cached            *AccessToken
	cachedBaseURL     string
	refresh           *refreshState
	lastAdvisoryError time.Time
}

type refreshState struct {
	done  chan struct{}
	token *AccessToken
	err   error
}

// NewTokenCache returns a [TokenCache] that delegates to the given provider.
// The handler is passed through to the provider on every fetch. The baseURL
// is resolved per-call in [TokenCache.Token] so the cache picks up whatever
// base URL the HTTP request is actually being sent to.
func NewTokenCache(provider TokenProvider, handler func(*http.Request) (*http.Response, error)) *TokenCache {
	return &TokenCache{provider: provider, httpHandler: handler}
}

// Token returns a valid access token, fetching or refreshing via the provider
// if necessary. baseURL is forwarded to the provider for token-exchange /
// refresh HTTP calls; a fresh cached token bypasses the provider entirely.
func (cache *TokenCache) Token(ctx context.Context, baseURL string) (string, error) {
	if ok, token := cache.check(ctx, baseURL); ok {
		return token, nil
	}
	return cache.fetch(ctx, baseURL)
}

// ForceToken is like [TokenCache.Token] but signals token providers that
// maintain their own caches (e.g., file-backed credential loaders) to bypass
// them and re-exchange the underlying credential. Call this after a 401 has
// invalidated the in-memory cache, so a stale on-disk token can't re-surface
// for the retry.
func (cache *TokenCache) ForceToken(ctx context.Context, baseURL string) (string, error) {
	return cache.Token(withForceRefresh(ctx), baseURL)
}

// Invalidate drops the in-memory cache so the next call to [TokenCache.Token]
// re-enters the provider. Pair with [TokenCache.ForceToken] on the retry path
// to also bypass provider-side caches. Returns true if there was a cached
// token to drop.
func (cache *TokenCache) Invalidate() bool {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if cache.cached == nil {
		return false
	}
	cache.cached = nil
	cache.cachedBaseURL = ""
	return true
}

// check inspects the cached token under the lock and returns whether
// the caller can use it directly or must fetch a new one.
func (cache *TokenCache) check(ctx context.Context, baseURL string) (bool, string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if cache.cached == nil {
		return false, ""
	}

	// Scope cached tokens to the baseURL they were minted for.
	if cache.cachedBaseURL != baseURL {
		cache.cached = nil
		cache.cachedBaseURL = ""
		return false, ""
	}

	if cache.cached.ExpiresAt == nil {
		// No expiry. Cache forever.
		return true, cache.cached.Token
	}

	remaining := time.Until(*cache.cached.ExpiresAt)

	if remaining > AdvisoryRefreshThreshold {
		// Fresh enough. Use the cached token.
		return true, cache.cached.Token
	}

	if remaining > MandatoryRefreshThreshold {
		// Advisory window. Kick off background refresh (subject to backoff
		// on repeated failure), return stale.
		if time.Since(cache.lastAdvisoryError) >= advisoryRefreshBackoff {
			cache.startBackgroundRefresh(ctx, baseURL)
		}
		return true, cache.cached.Token
	}

	// Mandatory window or expired.
	return false, ""
}

// fetch performs a synchronous token fetch. If another goroutine is already
// fetching, it waits for that result instead of starting a duplicate. Callers
// in the mandatory window (the only callers that reach fetch) never receive
// a stale token: if the refresh fails, the error propagates.
func (cache *TokenCache) fetch(ctx context.Context, baseURL string) (string, error) {
	isFetcher, refresh := cache.getOrStartRefresh()
	if !isFetcher {
		return awaitRefresh(ctx, refresh)
	}

	token, err := cache.provider(ctx, baseURL, cache.httpHandler)
	cache.completeRefresh(refresh, baseURL, token, err, false)
	if err != nil {
		return "", err
	}
	return token.Token, nil
}

// getOrStartRefresh returns the in-flight refresh if one exists, or starts a
// new one. The bool indicates whether the caller is the fetcher.
func (cache *TokenCache) getOrStartRefresh() (bool, *refreshState) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if refresh := cache.refresh; refresh != nil {
		return false, refresh
	}

	refresh := &refreshState{done: make(chan struct{})}
	cache.refresh = refresh
	return true, refresh
}

// awaitRefresh blocks until the in-flight refresh completes and returns its
// result. Piggy-backing callers never receive a stale token on failure; they
// are by construction in the mandatory window, where stale is unacceptable.
func awaitRefresh(ctx context.Context, refresh *refreshState) (string, error) {
	select {
	case <-refresh.done:
	case <-ctx.Done():
		return "", ctx.Err()
	}
	if refresh.err != nil {
		return "", refresh.err
	}
	return refresh.token.Token, nil
}

// completeRefresh records the result of a fetch, updates the cache, and
// signals all waiters. advisory indicates whether this refresh ran in the
// advisory (background) window — on failure, advisory refreshes arm the
// backoff instead of propagating the error.
func (cache *TokenCache) completeRefresh(refresh *refreshState, baseURL string, token *AccessToken, err error, advisory bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	refresh.token = token
	refresh.err = err
	if err == nil {
		cache.cached = token
		cache.cachedBaseURL = baseURL
		cache.lastAdvisoryError = time.Time{}
	} else if advisory {
		cache.lastAdvisoryError = time.Now()
	}
	cache.refresh = nil
	close(refresh.done)
}

// startBackgroundRefresh kicks off an advisory background refresh if one isn't
// already in-flight. Must be called with mu held.
func (cache *TokenCache) startBackgroundRefresh(ctx context.Context, baseURL string) {
	if cache.refresh != nil {
		return
	}

	refresh := &refreshState{done: make(chan struct{})}
	cache.refresh = refresh

	go func() {
		token, err := cache.provider(context.WithoutCancel(ctx), baseURL, cache.httpHandler)
		cache.completeRefresh(refresh, baseURL, token, err, true)
	}()
}

type forceRefreshKey struct{}

func withForceRefresh(ctx context.Context) context.Context {
	return context.WithValue(ctx, forceRefreshKey{}, true)
}

// isForceRefresh reports whether ctx carries a force-refresh signal. Token
// providers that maintain their own caches should skip them when true.
func isForceRefresh(ctx context.Context) bool {
	v, _ := ctx.Value(forceRefreshKey{}).(bool)
	return v
}
