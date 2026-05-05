package auth

import (
	"context"
	"encoding/json"
	"errors"
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
	var receivedRaw map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&receivedRaw); err != nil {
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
	if _, present := receivedRaw["service_account_id"]; present {
		t.Fatalf("service_account_id must be omitted when empty, body=%v", receivedRaw)
	}
	if _, present := receivedRaw["workspace_id"]; present {
		t.Fatalf("workspace_id must be omitted when empty, body=%v", receivedRaw)
	}
}

// TestWorkloadIdentityCredentials_WorkspaceIDIncluded verifies that an
// explicitly configured WorkspaceID is sent as workspace_id in the
// jwt-bearer exchange body so the server can mint a workspace-scoped token.
func TestWorkloadIdentityCredentials_WorkspaceIDIncluded(t *testing.T) {
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
		WorkspaceID:      "wrkspc_01abc",
		BaseURL:          server.URL,
	})

	if _, err := provider(context.Background(), "", http.DefaultClient.Do); err != nil {
		t.Fatal(err)
	}
	if receivedBody.WorkspaceID != "wrkspc_01abc" {
		t.Fatalf("got workspace_id %q, want %q", receivedBody.WorkspaceID, "wrkspc_01abc")
	}
}

// TestWorkloadIdentityCredentials_WorkspaceIDDefaultSentinel verifies the
// literal "default" sentinel is passed through unchanged so callers can
// pin the exchange to the organization's default workspace.
func TestWorkloadIdentityCredentials_WorkspaceIDDefaultSentinel(t *testing.T) {
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
		WorkspaceID:      "default",
		BaseURL:          server.URL,
	})

	if _, err := provider(context.Background(), "", http.DefaultClient.Do); err != nil {
		t.Fatal(err)
	}
	if receivedBody.WorkspaceID != "default" {
		t.Fatalf("got workspace_id %q, want %q", receivedBody.WorkspaceID, "default")
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

// errStatusServer starts a one-shot httptest server that replies with the
// given status and an RFC 6749-shaped error body. The caller owns Close.
func errStatusServer(t *testing.T, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(`{"error": "unauthorized"}`))
	}))
}

// TestWorkloadIdentityCredentials_401NoWorkspaceIDIncludesHint verifies that
// a 401 with no WorkspaceID surfaces the full hint: the always-on federation
// guidance and auth-event pointer, plus the workspace-scope clause — such a
// 401 is most often a federation rule spanning multiple workspaces, and the
// server can't pick one for us.
func TestWorkloadIdentityCredentials_401NoWorkspaceIDIncludesHint(t *testing.T) {
	server := errStatusServer(t, http.StatusUnauthorized)
	defer server.Close()

	provider := NewOIDCFederationCredentials(OIDCFederationConfig{
		IdentityProvider: &staticIdentity{token: "jwt"},
		FederationRuleID: "rule-1",
		OrganizationID:   "org-1",
		BaseURL:          server.URL,
	})

	_, err := provider(context.Background(), "", http.DefaultClient.Do)
	if err == nil {
		t.Fatal("expected error for 401")
	}
	msg := err.Error()
	if !strings.Contains(msg, "Ensure your federation rule matches your identity token") {
		t.Errorf("expected hint to include federation-rule guidance, got: %s", msg)
	}
	if !strings.Contains(msg, "ANTHROPIC_WORKSPACE_ID") {
		t.Errorf("expected hint to mention ANTHROPIC_WORKSPACE_ID, got: %s", msg)
	}
	if !strings.Contains(msg, "scoped to multiple workspaces") {
		t.Errorf("expected hint to mention 'scoped to multiple workspaces', got: %s", msg)
	}
	if !strings.Contains(msg, "workspace_id") {
		t.Errorf("expected hint to mention the 'workspace_id' config key, got: %s", msg)
	}
	// The hint must point at a type consumers can actually import — the
	// public option.FederationOptions, not the internal OIDCFederationConfig.
	if !strings.Contains(msg, "option.FederationOptions") {
		t.Errorf("expected hint to reference option.FederationOptions, got: %s", msg)
	}
	if strings.Contains(msg, "OIDCFederationConfig") {
		t.Errorf("hint must not reference internal OIDCFederationConfig, got: %s", msg)
	}
	if !strings.Contains(msg, "View your authentication events") {
		t.Errorf("expected hint to point at the authentication-events page, got: %s", msg)
	}
	// The relogin suffix is user_oauth-only; it must never run into this hint.
	if strings.Contains(msg, "anthropic auth login") {
		t.Errorf("workload-identity 401 must not suggest interactive relogin, got: %s", msg)
	}
}

// TestWorkloadIdentityCredentials_401WithWorkspaceIDOmitsWorkspaceHint
// verifies the workspace-scope clause is suppressed when WorkspaceID is
// already set: a 401 then has some other cause (revoked rule, expired
// assertion, ...) and that clause is noise. The always-on federation
// guidance and auth-event pointer still apply.
func TestWorkloadIdentityCredentials_401WithWorkspaceIDOmitsWorkspaceHint(t *testing.T) {
	server := errStatusServer(t, http.StatusUnauthorized)
	defer server.Close()

	provider := NewOIDCFederationCredentials(OIDCFederationConfig{
		IdentityProvider: &staticIdentity{token: "jwt"},
		FederationRuleID: "rule-1",
		OrganizationID:   "org-1",
		WorkspaceID:      "wrkspc_x",
		BaseURL:          server.URL,
	})

	_, err := provider(context.Background(), "", http.DefaultClient.Do)
	if err == nil {
		t.Fatal("expected error for 401")
	}
	msg := err.Error()
	if !strings.Contains(msg, "Ensure your federation rule") {
		t.Errorf("expected hint to include federation-rule guidance, got: %s", msg)
	}
	if !strings.Contains(msg, "View your authentication events") {
		t.Errorf("expected hint to point at the authentication-events page, got: %s", msg)
	}
	if strings.Contains(msg, "ANTHROPIC_WORKSPACE_ID") {
		t.Errorf("workspace-scope clause must be omitted when WorkspaceID is set, got: %s", msg)
	}
	if strings.Contains(msg, "scoped to multiple workspaces") {
		t.Errorf("workspace-scope clause must be omitted when WorkspaceID is set, got: %s", msg)
	}
	if strings.Contains(msg, "option.FederationOptions") {
		t.Errorf("workspace-scope clause must be omitted when WorkspaceID is set, got: %s", msg)
	}
}

// TestWorkloadIdentityCredentials_NoReloginSuggestion verifies that
// workload-identity token-exchange failures never suggest re-running
// `anthropic auth login`. Machine credentials have no interactive browser
// login to re-run; that remediation is only meaningful for the user_oauth
// flow. The cases below are exactly the statuses (and the invalid_grant
// RFC code) that trigger the relogin suffix on the user_oauth path — see
// shouldSuggestRelogin and TestOAuthTokenError_*SuggestsRelogin.
//
// Only the 401 row also carries the workload-identity hint; with the relogin
// suffix suppressed it must join directly onto the parsed RFC 6749 body as
// `<body>. Ensure your federation rule...` and read as prose.
func TestWorkloadIdentityCredentials_NoReloginSuggestion(t *testing.T) {
	cases := []struct {
		name     string
		status   int
		body     string
		wantHint bool
	}{
		{"401_unauthorized", http.StatusUnauthorized, `{"error":"unauthorized"}`, true},
		{"403_forbidden", http.StatusForbidden, `{"error":"access_denied"}`, false},
		{"400_bad_request", http.StatusBadRequest, `{"error":"invalid_request"}`, false},
		{"400_invalid_grant", http.StatusBadRequest, `{"error":"invalid_grant","error_description":"assertion expired"}`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.status)
				w.Write([]byte(tc.body))
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
				t.Fatalf("expected error for status %d", tc.status)
			}
			var oauthErr *OAuthTokenError
			if !errors.As(err, &oauthErr) {
				t.Fatalf("expected *OAuthTokenError, got %T: %v", err, err)
			}
			if !oauthErr.WorkloadIdentity {
				t.Error("expected WorkloadIdentity to be set on workload-identity errors")
			}
			msg := err.Error()
			if strings.Contains(msg, "anthropic auth login") {
				t.Errorf("workload-identity error must not suggest interactive relogin, got: %s", msg)
			}
			if tc.wantHint {
				if !strings.Contains(msg, ". Ensure your federation rule matches your identity token") {
					t.Errorf("expected hint to follow the body with a period-space join, got: %s", msg)
				}
				if !strings.Contains(msg, "View your authentication events") {
					t.Errorf("expected hint to point at the authentication-events page, got: %s", msg)
				}
			} else {
				if strings.Contains(msg, "Ensure your federation rule") {
					t.Errorf("hint must be omitted on non-401, got: %s", msg)
				}
				if strings.Contains(msg, "View your authentication events") {
					t.Errorf("hint must be omitted on non-401, got: %s", msg)
				}
			}
		})
	}
}

// TestWorkloadIdentityCredentials_Non401NoWorkspaceIDOmitsHint verifies the
// hint is 401-specific: a 5xx or a non-401 4xx shouldn't suggest a config
// change.
func TestWorkloadIdentityCredentials_Non401NoWorkspaceIDOmitsHint(t *testing.T) {
	server := errStatusServer(t, http.StatusInternalServerError)
	defer server.Close()

	provider := NewOIDCFederationCredentials(OIDCFederationConfig{
		IdentityProvider: &staticIdentity{token: "jwt"},
		FederationRuleID: "rule-1",
		OrganizationID:   "org-1",
		BaseURL:          server.URL,
	})

	_, err := provider(context.Background(), "", http.DefaultClient.Do)
	if err == nil {
		t.Fatal("expected error for 500")
	}
	msg := err.Error()
	if strings.Contains(msg, "Ensure your federation rule") {
		t.Errorf("hint must be omitted on non-401, got: %s", msg)
	}
	if strings.Contains(msg, "View your authentication events") {
		t.Errorf("hint must be omitted on non-401, got: %s", msg)
	}
	if strings.Contains(msg, "ANTHROPIC_WORKSPACE_ID") {
		t.Errorf("hint must be omitted on non-401, got: %s", msg)
	}
}
