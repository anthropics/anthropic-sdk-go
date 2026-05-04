package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type staticIdentity struct{ token string }

func (s *staticIdentity) GetIdentityToken(_ context.Context) (string, error) {
	return s.token, nil
}

type identityFunc func(ctx context.Context) (string, error)

func (f identityFunc) GetIdentityToken(ctx context.Context) (string, error) {
	return f(ctx)
}

func TestWorkloadIdentityCredentials_Exchange(t *testing.T) {
	var receivedBody tokenExchangeRequest
	var receivedBetaHeader string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedBetaHeader = r.Header.Get("anthropic-beta")
		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			t.Fatal(err)
		}
		expiresIn := 3600
		json.NewEncoder(w).Encode(tokenExchangeResponse{
			AccessToken: "access-tok-123",
			ExpiresIn:   &expiresIn,
		})
	}))
	defer server.Close()

	provider := NewOIDCFederationCredentials(OIDCFederationConfig{
		IdentityProvider: &staticIdentity{token: "my-jwt"},
		FederationRuleID: "rule-1",
		OrganizationID:   "org-1",
		ServiceAccountID: "sa-1",
		BaseURL:          server.URL,
	})

	token, err := provider(context.Background(), "", http.DefaultClient.Do)
	if err != nil {
		t.Fatal(err)
	}
	if token.Token != "access-tok-123" {
		t.Fatalf("got %q, want %q", token.Token, "access-tok-123")
	}
	if token.ExpiresAt == nil {
		t.Fatal("expected ExpiresAt to be set")
	}
	if !strings.Contains(receivedBetaHeader, FederationBetaHeader) {
		t.Fatalf("got beta header %q, want to contain %q", receivedBetaHeader, FederationBetaHeader)
	}
	if !strings.Contains(receivedBetaHeader, OAuthAPIBetaHeader) {
		t.Fatalf("got beta header %q, want to contain %q", receivedBetaHeader, OAuthAPIBetaHeader)
	}
	if receivedBody.GrantType != GrantTypeJWTBearer {
		t.Fatalf("got grant_type %q, want %q", receivedBody.GrantType, GrantTypeJWTBearer)
	}
	if receivedBody.Assertion != "my-jwt" {
		t.Fatalf("got assertion %q, want %q", receivedBody.Assertion, "my-jwt")
	}
	if receivedBody.FederationRuleID != "rule-1" {
		t.Fatalf("got federation_rule_id %q, want %q", receivedBody.FederationRuleID, "rule-1")
	}
	if receivedBody.OrganizationID != "org-1" {
		t.Fatalf("got organization_id %q, want %q", receivedBody.OrganizationID, "org-1")
	}
	if receivedBody.ServiceAccountID != "sa-1" {
		t.Fatalf("got service_account_id %q, want %q", receivedBody.ServiceAccountID, "sa-1")
	}
}

func TestWorkloadIdentityCredentials_OptionalFields(t *testing.T) {
	var receivedBody tokenExchangeRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&receivedBody); err != nil {
			t.Fatal(err)
		}
		json.NewEncoder(w).Encode(tokenExchangeResponse{AccessToken: "tok"})
	}))
	defer server.Close()

	provider := NewOIDCFederationCredentials(OIDCFederationConfig{
		IdentityProvider: &staticIdentity{token: "jwt"},
		FederationRuleID: "rule-1",
		OrganizationID:   "org-1",
		BaseURL:          server.URL,
	})

	token, err := provider(context.Background(), "", http.DefaultClient.Do)
	if err != nil {
		t.Fatal(err)
	}
	if token.ExpiresAt != nil {
		t.Fatal("expected ExpiresAt to be nil when expires_in is absent")
	}
	if receivedBody.ServiceAccountID != "" {
		t.Fatalf("expected empty service_account_id, got %q", receivedBody.ServiceAccountID)
	}
}

func TestWorkloadIdentityCredentials_ReInvokesIdentityProvider(t *testing.T) {
	calls := 0
	identity := identityFunc(func(_ context.Context) (string, error) {
		calls++
		return "jwt-" + string(rune('0'+calls)), nil
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(tokenExchangeResponse{AccessToken: "tok"})
	}))
	defer server.Close()

	provider := NewOIDCFederationCredentials(OIDCFederationConfig{
		IdentityProvider: identity,
		FederationRuleID: "rule-1",
		OrganizationID:   "org-1",
		BaseURL:          server.URL,
	})

	_, _ = provider(context.Background(), "", http.DefaultClient.Do)
	_, _ = provider(context.Background(), "", http.DefaultClient.Do)

	if calls != 2 {
		t.Fatalf("expected 2 identity provider calls, got %d", calls)
	}
}

// TestWorkloadIdentityCredentials_RejectsOversizedAssertion checks that an
// identity-token payload over the 16 KiB cap is rejected before any request
// goes out, so an over-large or attacker-crafted assertion can't be
// shipped to the token endpoint at request rate.
func TestWorkloadIdentityCredentials_RejectsOversizedAssertion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("server should not be hit when assertion is over-sized")
	}))
	defer server.Close()

	provider := NewOIDCFederationCredentials(OIDCFederationConfig{
		IdentityProvider: &staticIdentity{token: strings.Repeat("a", MaxAssertionSize+1)},
		FederationRuleID: "rule-1",
		OrganizationID:   "org-1",
		BaseURL:          server.URL,
	})

	_, err := provider(context.Background(), "", http.DefaultClient.Do)
	if err == nil {
		t.Fatal("expected error for over-sized assertion")
	}
	if !strings.Contains(err.Error(), "16384") {
		t.Errorf("expected size limit in error, got: %v", err)
	}
}

// TestWorkloadIdentityCredentials_RejectsOversizedResponseBody checks that
// a token-endpoint response over 1 MiB is rejected rather than read into
// memory in full. Mirrors the Python SDK's bound.
func TestWorkloadIdentityCredentials_RejectsOversizedResponseBody(t *testing.T) {
	huge := strings.Repeat("x", 2<<20) // 2 MiB > 1 MiB cap
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Wrap in valid JSON so the parse doesn't fail before the read does.
		io.WriteString(w, `{"access_token":"`+huge+`"}`)
	}))
	defer server.Close()

	provider := NewOIDCFederationCredentials(OIDCFederationConfig{
		IdentityProvider: &staticIdentity{token: "jwt"},
		FederationRuleID: "rule-1",
		OrganizationID:   "org-1",
		BaseURL:          server.URL,
	})

	_, err := provider(context.Background(), "", http.DefaultClient.Do)
	if err == nil {
		t.Fatal("expected parse error from truncated 2 MiB response")
	}
}

// TestWorkloadIdentityCredentials_RejectsNonBearerTokenType checks that the
// federation grant rejects a token_type other than Bearer (case-insensitive).
// Python validates this; Go user_oauth refresh did too, but federation did not
// until parity sweep.
func TestWorkloadIdentityCredentials_RejectsNonBearerTokenType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(tokenExchangeResponse{
			AccessToken: "tok",
			TokenType:   "Mac", // not Bearer
		})
	}))
	defer server.Close()

	provider := NewOIDCFederationCredentials(OIDCFederationConfig{
		IdentityProvider: &staticIdentity{token: "jwt"},
		FederationRuleID: "rule-1",
		OrganizationID:   "org-1",
		BaseURL:          server.URL,
	})

	_, err := provider(context.Background(), "", http.DefaultClient.Do)
	if err == nil {
		t.Fatal("expected error for non-Bearer token_type")
	}
	if !strings.Contains(err.Error(), "Mac") || !strings.Contains(err.Error(), "Bearer") {
		t.Errorf("expected token_type and Bearer in error, got: %v", err)
	}
}

// TestWorkloadIdentityCredentials_SetsUserAgent verifies the federation
// token POST carries an SDK User-Agent so server-side logs can attribute
// the call to this SDK version.
func TestWorkloadIdentityCredentials_SetsUserAgent(t *testing.T) {
	var receivedUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedUA = r.Header.Get("User-Agent")
		json.NewEncoder(w).Encode(tokenExchangeResponse{AccessToken: "tok"})
	}))
	defer server.Close()

	provider := NewOIDCFederationCredentials(OIDCFederationConfig{
		IdentityProvider: &staticIdentity{token: "jwt"},
		FederationRuleID: "rule-1",
		OrganizationID:   "org-1",
		BaseURL:          server.URL,
	})

	_, err := provider(context.Background(), "", http.DefaultClient.Do)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(receivedUA, "anthropic-sdk-go/") {
		t.Errorf("got User-Agent %q, want prefix anthropic-sdk-go/", receivedUA)
	}
	if !strings.Contains(receivedUA, "oidc-federation") {
		t.Errorf("got User-Agent %q, want context oidc-federation", receivedUA)
	}
}

func TestWorkloadIdentityCredentials_ErrorStatus(t *testing.T) {
	for _, status := range []int{403, 503} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			w.Write([]byte(`{"error": "bad"}`))
		}))

		provider := NewOIDCFederationCredentials(OIDCFederationConfig{
			IdentityProvider: &staticIdentity{token: "jwt"},
			FederationRuleID: "rule-1",
			OrganizationID:   "org-1",
			BaseURL:          server.URL,
		})

		_, err := provider(context.Background(), "", http.DefaultClient.Do)
		server.Close()

		if err == nil {
			t.Fatalf("expected error for status %d", status)
		}
		wierr, ok := err.(*OAuthTokenError)
		if !ok {
			t.Fatalf("expected *OAuthTokenError, got %T", err)
		}
		if wierr.StatusCode != status {
			t.Fatalf("got status %d, want %d", wierr.StatusCode, status)
		}
	}
}
