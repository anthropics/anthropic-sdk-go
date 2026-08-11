package option_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestIdentityTokenFile_ReadsAndTrims(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token")
	if err := os.WriteFile(path, []byte("  the-jwt\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	tok, err := option.IdentityTokenFile(path)(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "the-jwt" {
		t.Errorf("got %q, want %q", tok, "the-jwt")
	}
}

func TestIdentityTokenFile_RereadsOnEachCall(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token")
	if err := os.WriteFile(path, []byte("v1"), 0o600); err != nil {
		t.Fatal(err)
	}
	fn := option.IdentityTokenFile(path)

	got1, err := fn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got1 != "v1" {
		t.Errorf("first call got %q, want %q", got1, "v1")
	}

	// Simulate token rotation.
	if err := os.WriteFile(path, []byte("v2"), 0o600); err != nil {
		t.Fatal(err)
	}
	got2, err := fn(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got2 != "v2" {
		t.Errorf("after rotation got %q, want %q (function did not re-read)", got2, "v2")
	}
}

func TestIdentityTokenFile_Errors(t *testing.T) {
	ctx := context.Background()

	t.Run("empty path", func(t *testing.T) {
		if _, err := option.IdentityTokenFile("")(ctx); err == nil || !strings.Contains(err.Error(), "path is empty") {
			t.Errorf("got %v, want error containing 'path is empty'", err)
		}
	})

	t.Run("missing file", func(t *testing.T) {
		missing := filepath.Join(t.TempDir(), "does-not-exist")
		if _, err := option.IdentityTokenFile(missing)(ctx); err == nil || !strings.Contains(err.Error(), "read") {
			t.Errorf("got %v, want read error", err)
		}
	})

	t.Run("whitespace only", func(t *testing.T) {
		empty := filepath.Join(t.TempDir(), "empty")
		if err := os.WriteFile(empty, []byte("   \n\t\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := option.IdentityTokenFile(empty)(ctx); err == nil || !strings.Contains(err.Error(), "is empty") {
			t.Errorf("got %v, want 'is empty' error", err)
		}
	})
}

// TestJoin_AppliesInOrder verifies that a joined option is equivalent to
// passing its members individually at the same position: members apply in
// order (a later WithMaxRetries wins) and the first error stops the join.
func TestJoin_AppliesInOrder(t *testing.T) {
	cfg := &requestconfig.RequestConfig{}
	joined := option.Join(
		option.WithMaxRetries(1),
		option.WithBaseURL("https://joined.example/"),
		option.WithMaxRetries(7),
	)
	if err := joined.Apply(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxRetries != 7 {
		t.Errorf("got MaxRetries %d, want the later member (7) to win", cfg.MaxRetries)
	}
	if got := cfg.BaseURL.String(); got != "https://joined.example/" {
		t.Errorf("got base URL %q, want %q", got, "https://joined.example/")
	}

	boom := errors.New("boom")
	applied := false
	failing := option.Join(
		requestconfig.RequestOptionFunc(func(*requestconfig.RequestConfig) error { return boom }),
		requestconfig.RequestOptionFunc(func(*requestconfig.RequestConfig) error { applied = true; return nil }),
	)
	if err := failing.Apply(&requestconfig.RequestConfig{}); !errors.Is(err, boom) {
		t.Fatalf("got error %v, want the member's error", err)
	}
	if applied {
		t.Error("members after a failing member must not be applied")
	}
}

// TestHasWithoutEnvironmentDefaults_SeesThroughJoin pins the contract that
// platform option constructors rely on: a marker returned inside a Join
// (bedrock/vertex bundle it with their own options) must still be detected
// by anthropic.NewClient, however deeply the joins nest — while a Join with
// no marker must not be mistaken for one.
func TestHasWithoutEnvironmentDefaults_SeesThroughJoin(t *testing.T) {
	noop := option.WithHeader("x-noop", "1")
	testCases := []struct {
		name string
		opts []option.RequestOption
		want bool
	}{
		{
			name: "no options",
			opts: nil,
			want: false,
		},
		{
			name: "unrelated options only",
			opts: []option.RequestOption{noop, option.WithAPIKey("k")},
			want: false,
		},
		{
			name: "marker at the top level",
			opts: []option.RequestOption{noop, option.WithoutEnvironmentDefaults()},
			want: true,
		},
		{
			name: "marker nested one join deep",
			opts: []option.RequestOption{noop, option.Join(noop, option.WithoutEnvironmentDefaults(), noop)},
			want: true,
		},
		{
			name: "marker nested two joins deep",
			opts: []option.RequestOption{option.Join(noop, option.Join(noop, option.WithoutEnvironmentDefaults()))},
			want: true,
		},
		{
			name: "joins without a marker",
			opts: []option.RequestOption{option.Join(noop, option.Join(noop, option.WithAPIKey("k"))), option.Join()},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := option.HasWithoutEnvironmentDefaults(tc.opts); got != tc.want {
				t.Errorf("HasWithoutEnvironmentDefaults() = %v, want %v", got, tc.want)
			}
		})
	}
}
