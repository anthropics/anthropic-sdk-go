package anthropic_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestPathParamsRejectDotDotEscape(t *testing.T) {
	var gotURI string
	var hits int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		gotURI = r.URL.RequestURI()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{}`)
	}))
	defer server.Close()

	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithBaseURL(server.URL),
		option.WithMaxRetries(0),
	)
	_, err := client.Beta.Vaults.Credentials.Get(
		context.Background(),
		"../../vault_B/credentials/cred_x",
		anthropic.BetaVaultCredentialGetParams{VaultID: "vault_A"},
	)
	if err == nil {
		if !strings.Contains(gotURI, "/v1/vaults/vault_A/credentials/") {
			t.Fatalf("request escaped the intended prefix: %s", gotURI)
		}
		if strings.Contains(gotURI, "/vault_B/") && !strings.Contains(gotURI, "/vault_A/") {
			t.Fatalf("request walked to vault_B: %s", gotURI)
		}
		return
	}
	if hits != 0 {
		t.Fatalf("dot-dot id issued a request to %s (hits=%d)", gotURI, hits)
	}
}

func TestPathParamsQuestionMarkDoesNotInjectQuery(t *testing.T) {
	var got *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Clone(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{}`)
	}))
	defer server.Close()

	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithBaseURL(server.URL),
		option.WithMaxRetries(0),
	)
	_, _ = client.Beta.Vaults.Credentials.Get(
		context.Background(),
		"cred?injected=1",
		anthropic.BetaVaultCredentialGetParams{VaultID: "vault_A"},
	)
	if got == nil {
		t.Fatal("expected a request")
	}
	if _, ok := got.URL.Query()["injected"]; ok {
		t.Fatalf("id injected a query parameter: %s", got.URL.RequestURI())
	}
	if got.URL.Query().Get("beta") != "true" {
		t.Fatalf("expected beta=true query, got %s", got.URL.RawQuery)
	}
	if !strings.Contains(got.URL.EscapedPath(), "cred%3Finjected=1") {
		t.Fatalf("expected escaped id in path, got %s", got.URL.RequestURI())
	}
}

func TestPathParamsHashDoesNotDropSuffix(t *testing.T) {
	var got *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Clone(r.Context())
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "ok")
	}))
	defer server.Close()

	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithBaseURL(server.URL),
		option.WithMaxRetries(0),
	)
	resp, err := client.Files.Download(context.Background(), "foo#bar")
	if err != nil {
		t.Fatalf("Download: %v", err)
	}
	defer resp.Body.Close()
	if got == nil {
		t.Fatal("expected a request")
	}
	if got.URL.Fragment != "" {
		t.Fatalf("id truncated the path into a fragment: path=%q fragment=%q uri=%s", got.URL.Path, got.URL.Fragment, got.URL.RequestURI())
	}
	if !strings.HasSuffix(got.URL.Path, "/content") {
		t.Fatalf("id dropped the /content suffix: %s", got.URL.RequestURI())
	}
	if !strings.Contains(got.URL.EscapedPath(), "foo%23bar") {
		t.Fatalf("expected escaped id in path, got %s", got.URL.RequestURI())
	}
}

func TestPathParamsEmptyIDStillErrors(t *testing.T) {
	var hits int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{}`)
	}))
	defer server.Close()

	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithBaseURL(server.URL),
		option.WithMaxRetries(0),
	)
	_, err := client.Beta.Vaults.Credentials.Get(
		context.Background(),
		"",
		anthropic.BetaVaultCredentialGetParams{VaultID: "vault_A"},
	)
	if err == nil {
		t.Fatal("expected empty credential id to error")
	}
	if !strings.Contains(err.Error(), "missing required credential_id") {
		t.Fatalf("expected missing credential_id error, got %v", err)
	}
	if hits != 0 {
		t.Fatalf("empty id issued a request (hits=%d)", hits)
	}
}
