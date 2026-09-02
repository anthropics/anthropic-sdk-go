package anthropic_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestPathParamsRejectDotDotEscape(t *testing.T) {
	var gotURI string
	var gotEscaped string
	var hits int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		gotURI = r.URL.RequestURI()
		gotEscaped = r.URL.EscapedPath()
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
		"../../vault_B/credentials/cred_x",
		anthropic.BetaVaultCredentialGetParams{VaultID: "vault_A"},
	)
	if hits != 1 {
		t.Fatalf("expected a request, got hits=%d uri=%s", hits, gotURI)
	}
	if !strings.HasPrefix(gotEscaped, "/v1/vaults/vault_A/credentials/") {
		t.Fatalf("request escaped the intended prefix: %s", gotURI)
	}
	if strings.Contains(gotEscaped, "/vault_B/") {
		t.Fatalf("request walked to vault_B: %s", gotURI)
	}
	if !strings.Contains(gotEscaped, "..%2F..%2Fvault_B%2Fcredentials%2Fcred_x") {
		t.Fatalf("expected encoded traversal to stay in the id segment, got %s", gotURI)
	}
}

func TestPathParamsSlashStaysInSegment(t *testing.T) {
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
		"foo/bar",
		anthropic.BetaVaultCredentialGetParams{VaultID: "vault_A"},
	)
	if got == nil {
		t.Fatal("expected a request")
	}
	escaped := got.URL.EscapedPath()
	if !strings.Contains(escaped, "foo%2Fbar") {
		t.Fatalf("expected foo/bar to be a single escaped segment, got %s", got.URL.RequestURI())
	}
	if strings.Contains(escaped, "/foo/bar") {
		t.Fatalf("slash in id became an extra path segment: %s", got.URL.RequestURI())
	}
	if got.URL.Query().Get("beta") != "true" {
		t.Fatalf("expected beta=true query, got %s", got.URL.RawQuery)
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

func TestGeneratedResourcePathsUseFormatPath(t *testing.T) {
	matches, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range matches {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Contains(src, []byte("File generated from our OpenAPI spec by Stainless")) {
			continue
		}
		for i, line := range strings.Split(string(src), "\n") {
			if strings.Contains(line, "path := fmt.Sprintf(") {
				t.Errorf("%s:%d: generated path still uses fmt.Sprintf; use requestconfig.FormatPath so IDs are escaped before interpolation", path, i+1)
			}
		}
	}
}
