package option_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

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
