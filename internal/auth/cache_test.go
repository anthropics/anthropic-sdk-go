package auth

import (
	"context"
	"errors"
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

func TestTokenCache_ConcurrentFetchForOtherBaseURL(t *testing.T) {
	const baseA, baseB = "https://a.example.com", "https://b.example.com"
	started := make(chan struct{})
	release := make(chan struct{})
	var startOnce sync.Once
	cache := newTestCache(func(ctx context.Context, baseURL string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		if baseURL == baseA {
			startOnce.Do(func() { close(started) })
			<-release
		}
		return &AccessToken{Token: "tok " + baseURL}, nil
	})

	gotA := make(chan string, 1)
	go func() {
		token, err := cache.Token(context.Background(), baseA)
		if err != nil {
			t.Error(err)
		}
		gotA <- token
	}()
	<-started

	gotB := make(chan string, 1)
	go func() {
		token, err := cache.Token(context.Background(), baseB)
		if err != nil {
			t.Error(err)
		}
		gotB <- token
	}()

	select {
	case token := <-gotB:
		if token != "tok "+baseB {
			t.Fatalf("baseB got %q, want %q", token, "tok "+baseB)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("baseB waited for baseA's fetch")
	}
	close(release)

	if token := <-gotA; token != "tok "+baseA {
		t.Fatalf("baseA got %q, want %q", token, "tok "+baseA)
	}
}

func TestTokenCache_BackgroundRefreshForOtherBaseURL(t *testing.T) {
	const baseA, baseB = "https://a.example.com", "https://b.example.com"
	started := make(chan struct{})
	release := make(chan struct{})
	callsA := atomic.Int32{}
	cache := newTestCache(func(ctx context.Context, baseURL string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		if baseURL == baseA {
			switch callsA.Add(1) {
			case 1:
				exp := time.Now().Add(60 * time.Second) // within advisory window
				return &AccessToken{Token: "stale " + baseA, ExpiresAt: &exp}, nil
			case 2:
				close(started)
				<-release
			}
		}
		exp := time.Now().Add(10 * time.Minute)
		return &AccessToken{Token: "tok " + baseURL, ExpiresAt: &exp}, nil
	})

	// Prime the cache, then start a background refresh for baseA.
	for range 2 {
		if _, err := cache.Token(context.Background(), baseA); err != nil {
			t.Fatal(err)
		}
	}
	<-started

	gotB := make(chan string, 1)
	go func() {
		token, err := cache.Token(context.Background(), baseB)
		if err != nil {
			t.Error(err)
		}
		gotB <- token
	}()
	// Let the baseB call reach the cache while baseA's refresh is in flight.
	time.Sleep(50 * time.Millisecond)
	close(release)

	if token := <-gotB; token != "tok "+baseB {
		t.Fatalf("baseB got %q, want %q", token, "tok "+baseB)
	}
}

func TestTokenCache_FailedFetchForOtherBaseURLNotShared(t *testing.T) {
	const baseA, baseB = "https://a.example.com", "https://b.example.com"
	started := make(chan struct{})
	release := make(chan struct{})
	cache := newTestCache(func(ctx context.Context, baseURL string, _ func(*http.Request) (*http.Response, error)) (*AccessToken, error) {
		if baseURL == baseA {
			close(started)
			<-release
			return nil, errors.New("exchange failed")
		}
		return &AccessToken{Token: "tok " + baseURL}, nil
	})

	errA := make(chan error, 1)
	go func() {
		_, err := cache.Token(context.Background(), baseA)
		errA <- err
	}()
	<-started

	type result struct {
		token string
		err   error
	}
	gotB := make(chan result, 1)
	go func() {
		token, err := cache.Token(context.Background(), baseB)
		gotB <- result{token, err}
	}()
	// Let the baseB call reach the cache while baseA's fetch is in flight.
	time.Sleep(50 * time.Millisecond)
	close(release)

	if err := <-errA; err == nil {
		t.Fatal("baseA: expected an error")
	}
	if r := <-gotB; r.err != nil || r.token != "tok "+baseB {
		t.Fatalf("baseB got (%q, %v), want (%q, nil)", r.token, r.err, "tok "+baseB)
	}
}
