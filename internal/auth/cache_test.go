package auth

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func newTestCache(provider TokenProvider) *TokenCache {
	return NewTokenCache(provider, nil)
}

func TestTokenCache_FirstCallFetches(t *testing.T) {
	calls := 0
	cache := newTestCache(func(ctx context.Context, _ string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		calls++
		return &AccessToken{Token: "tok-1"}, nil
	})

	token, err := cache.Token(context.Background(), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if token != "tok-1" {
		t.Fatalf("got %q, want %q", token, "tok-1")
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestTokenCache_NoExpiryNeverRefreshes(t *testing.T) {
	calls := 0
	cache := newTestCache(func(ctx context.Context, _ string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		calls++
		return &AccessToken{Token: "tok-forever", ExpiresAt: nil}, nil
	})

	for range 5 {
		token, err := cache.Token(context.Background(), "https://example.com")
		if err != nil {
			t.Fatal(err)
		}
		if token != "tok-forever" {
			t.Fatalf("got %q, want %q", token, "tok-forever")
		}
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestTokenCache_FreshTokenNotRefetched(t *testing.T) {
	calls := 0
	cache := newTestCache(func(ctx context.Context, _ string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		calls++
		exp := time.Now().Add(10 * time.Minute)
		return &AccessToken{Token: "fresh", ExpiresAt: &exp}, nil
	})

	for range 3 {
		token, err := cache.Token(context.Background(), "https://example.com")
		if err != nil {
			t.Fatal(err)
		}
		if token != "fresh" {
			t.Fatalf("got %q, want %q", token, "fresh")
		}
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestTokenCache_AdvisoryRefreshSuccess(t *testing.T) {
	callCount := atomic.Int32{}

	cache := newTestCache(func(ctx context.Context, _ string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		n := callCount.Add(1)
		if n == 1 {
			exp := time.Now().Add(60 * time.Second) // within advisory window
			return &AccessToken{Token: "stale", ExpiresAt: &exp}, nil
		}
		exp := time.Now().Add(10 * time.Minute)
		return &AccessToken{Token: "refreshed", ExpiresAt: &exp}, nil
	})

	// Prime the cache.
	token, err := cache.Token(context.Background(), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if token != "stale" {
		t.Fatalf("got %q, want %q", token, "stale")
	}

	// Second call should return stale and trigger background refresh.
	token, err = cache.Token(context.Background(), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if token != "stale" {
		t.Fatalf("got %q, want %q", token, "stale")
	}

	// Poll until the background refresh completes rather than sleeping a
	// fixed duration, which is inherently flaky.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		token, err = cache.Token(context.Background(), "https://example.com")
		if err != nil {
			t.Fatal(err)
		}
		if token == "refreshed" {
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
	t.Fatal("background refresh did not complete in time")
}

func TestTokenCache_AdvisoryRefreshFailureServesStale(t *testing.T) {
	callCount := atomic.Int32{}
	refreshDone := make(chan struct{})

	cache := newTestCache(func(ctx context.Context, _ string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		n := callCount.Add(1)
		if n == 1 {
			exp := time.Now().Add(60 * time.Second)
			return &AccessToken{Token: "stale", ExpiresAt: &exp}, nil
		}
		defer close(refreshDone)
		return nil, fmt.Errorf("refresh failed")
	})

	// Prime.
	_, _ = cache.Token(context.Background(), "https://example.com")

	// Should return stale despite background refresh failure.
	token, err := cache.Token(context.Background(), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if token != "stale" {
		t.Fatalf("got %q, want %q", token, "stale")
	}

	<-refreshDone
}

func TestTokenCache_MandatoryRefreshFailureReturnsError(t *testing.T) {
	callCount := atomic.Int32{}
	cache := newTestCache(func(ctx context.Context, _ string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		n := callCount.Add(1)
		if n == 1 {
			exp := time.Now().Add(10 * time.Second) // within mandatory window
			return &AccessToken{Token: "expiring", ExpiresAt: &exp}, nil
		}
		return nil, fmt.Errorf("refresh failed")
	})

	// Prime.
	_, _ = cache.Token(context.Background(), "https://example.com")

	// Mandatory refresh should propagate the error.
	_, err := cache.Token(context.Background(), "https://example.com")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestTokenCache_ExpiredIsMandatory(t *testing.T) {
	callCount := atomic.Int32{}
	cache := newTestCache(func(ctx context.Context, _ string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		n := callCount.Add(1)
		if n == 1 {
			exp := time.Now().Add(-1 * time.Second) // already expired
			return &AccessToken{Token: "expired", ExpiresAt: &exp}, nil
		}
		exp := time.Now().Add(10 * time.Minute)
		return &AccessToken{Token: "new", ExpiresAt: &exp}, nil
	})

	// Prime with expired token.
	_, _ = cache.Token(context.Background(), "https://example.com")

	token, err := cache.Token(context.Background(), "https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	if token != "new" {
		t.Fatalf("got %q, want %q", token, "new")
	}
}

func TestTokenCache_Invalidate(t *testing.T) {
	callCount := atomic.Int32{}
	cache := newTestCache(func(ctx context.Context, _ string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		n := callCount.Add(1)
		return &AccessToken{Token: fmt.Sprintf("tok-%d", n)}, nil
	})

	token, _ := cache.Token(context.Background(), "https://example.com")
	if token != "tok-1" {
		t.Fatalf("got %q, want %q", token, "tok-1")
	}

	cache.Invalidate()

	token, _ = cache.Token(context.Background(), "https://example.com")
	if token != "tok-2" {
		t.Fatalf("got %q, want %q", token, "tok-2")
	}
}

func TestTokenCache_ConcurrentFetchDeduplication(t *testing.T) {
	calls := atomic.Int32{}
	cache := newTestCache(func(ctx context.Context, _ string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		calls.Add(1)
		time.Sleep(50 * time.Millisecond)
		return &AccessToken{Token: "concurrent"}, nil
	})

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token, err := cache.Token(context.Background(), "https://example.com")
			if err != nil {
				t.Error(err)
				return
			}
			if token != "concurrent" {
				t.Errorf("got %q, want %q", token, "concurrent")
			}
		}()
	}
	wg.Wait()

	if n := calls.Load(); n != 1 {
		t.Fatalf("expected 1 provider call, got %d", n)
	}
}
