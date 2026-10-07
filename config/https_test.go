package config_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/config"
)

func TestRequireSecureTokenEndpoint(t *testing.T) {
	cases := map[string]struct {
		base    string
		wantErr string
	}{
		"https":                    {"https://api.anthropic.com", ""},
		"https on another host":    {"https://gateway.example.com/anthropic", ""},
		"empty uses the default":   {"", ""},
		"http on localhost":        {"http://localhost:8080", ""},
		"http on LOCALHOST":        {"http://LOCALHOST:8080", ""},
		"http on 127.0.0.1":        {"http://127.0.0.1:8080", ""},
		"http on ::1":              {"http://[::1]:8080", ""},
		"http on another host":     {"http://api.anthropic.com", "non-https token endpoint"},
		"http on a lookalike host": {"http://localhost.example.com", "non-https token endpoint"},
		"http on 127.0.0.2":        {"http://127.0.0.2", "non-https token endpoint"},
		"no scheme":                {"api.anthropic.com", "non-https token endpoint"},
		"another scheme":           {"ftp://api.anthropic.com", "non-https token endpoint"},
		"unparseable":              {"http://[::1", "invalid token endpoint base URL"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := config.RequireSecureTokenEndpoint(tc.base)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("error = %v, want one containing %q", err, tc.wantErr)
			}
		})
	}
}

type recordingTransport struct{ calls int }

func (r *recordingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	r.calls++
	return nil, http.ErrHandlerTimeout
}

func TestExchangeFederationAssertion_RefusesCleartextBaseURL(t *testing.T) {
	transport := &recordingTransport{}
	_, err := config.ExchangeFederationAssertion(context.Background(), config.FederationExchangeParams{
		Assertion:        "jwt",
		FederationRuleID: "fdrl_1",
		OrganizationID:   "org",
		BaseURL:          "http://gateway.example.com",
		HTTPClient:       &http.Client{Transport: transport},
	})
	if err == nil || !strings.Contains(err.Error(), "non-https token endpoint") {
		t.Fatalf("error = %v, want a refusal", err)
	}
	if transport.calls != 0 {
		t.Fatalf("the assertion was sent %d time(s)", transport.calls)
	}
}

func TestExchangeFederationAssertion_DoesNotFollowRedirects(t *testing.T) {
	supplied := &http.Client{}
	for name, client := range map[string]*http.Client{"default client": nil, "supplied client": supplied} {
		t.Run(name, func(t *testing.T) {
			replayed := 0
			elsewhere := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { replayed++ }))
			defer elsewhere.Close()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, elsewhere.URL+r.URL.Path, http.StatusTemporaryRedirect)
			}))
			defer server.Close()

			_, err := config.ExchangeFederationAssertion(context.Background(), config.FederationExchangeParams{
				Assertion:        "jwt",
				FederationRuleID: "fdrl_1",
				OrganizationID:   "org",
				BaseURL:          server.URL,
				HTTPClient:       client,
			})
			var exchangeErr *config.FederationExchangeError
			if !errors.As(err, &exchangeErr) || exchangeErr.StatusCode != http.StatusTemporaryRedirect {
				t.Fatalf("error = %v, want the 307 surfaced", err)
			}
			if replayed != 0 {
				t.Fatalf("the assertion was replayed to the redirect target %d time(s)", replayed)
			}
		})
	}
	if supplied.CheckRedirect != nil {
		t.Fatal("the caller's client was modified")
	}
}
