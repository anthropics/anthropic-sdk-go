package config_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/config"
)

func TestExchangeFederationAssertion_Success(t *testing.T) {
	var gotReq *http.Request
	var gotBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotReq = r.Clone(context.Background())
		gotBody, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Request-Id", "req_success")
		expiresIn := 3600
		json.NewEncoder(w).Encode(map[string]any{
			"access_token": "fed-access",
			"token_type":   "Bearer",
			"expires_in":   expiresIn,
		})
	}))
	defer server.Close()

	creds, err := config.ExchangeFederationAssertion(context.Background(), config.FederationExchangeParams{
		Assertion:        "my-jwt",
		FederationRuleID: "fdrl_abc",
		OrganizationID:   "org_123",
		ServiceAccountID: "svac_def",
		WorkspaceID:      "wrkspc_x",
		BaseURL:          server.URL,
	})
	if err != nil {
		t.Fatal(err)
	}
	if creds.AccessToken != "fed-access" {
		t.Errorf("AccessToken: %q", creds.AccessToken)
	}
	if creds.RefreshToken != "" {
		t.Errorf("federation should not return refresh_token, got %q", creds.RefreshToken)
	}
	if creds.ExpiresAt == nil {
		t.Error("ExpiresAt should be populated")
	}

	if gotReq.URL.Path != "/v1/oauth/token" {
		t.Errorf("path: %q", gotReq.URL.Path)
	}
	if ct := gotReq.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type: %q", ct)
	}
	beta := gotReq.Header.Get("anthropic-beta")
	if !strings.Contains(beta, "oauth-2025-04-20") || !strings.Contains(beta, "oidc-federation-2026-04-01") {
		t.Errorf("anthropic-beta header missing required values: %q", beta)
	}

	// The JSON body's field names must match the REST gateway's alias
	// transformation exactly — grant_type, assertion, federation_rule_id,
	// organization_id, service_account_id, workspace_id. Anything else is
	// silently dropped by the gateway and produces an empty-field error
	// downstream.
	var decoded map[string]any
	if err := json.Unmarshal(gotBody, &decoded); err != nil {
		t.Fatalf("body must be JSON: %v, got=%q", err, gotBody)
	}
	want := map[string]any{
		"grant_type":         "urn:ietf:params:oauth:grant-type:jwt-bearer",
		"assertion":          "my-jwt",
		"federation_rule_id": "fdrl_abc",
		"organization_id":    "org_123",
		"service_account_id": "svac_def",
		"workspace_id":       "wrkspc_x",
	}
	for k, v := range want {
		if decoded[k] != v {
			t.Errorf("body[%q] = %v, want %v", k, decoded[k], v)
		}
	}
	for _, forbidden := range []string{"federation_rule", "organization"} {
		if _, present := decoded[forbidden]; present {
			t.Errorf("body must not use alias %q, got %q", forbidden, gotBody)
		}
	}
}

func TestExchangeFederationAssertion_SendsDefaultUserAgent(t *testing.T) {
	var gotUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		json.NewEncoder(w).Encode(map[string]any{"access_token": "t"})
	}))
	defer server.Close()

	if _, err := config.ExchangeFederationAssertion(context.Background(), config.FederationExchangeParams{
		Assertion:        "jwt",
		FederationRuleID: "fdrl_1",
		OrganizationID:   "org_1",
		BaseURL:          server.URL,
	}); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(gotUA, "anthropic-sdk-go/") {
		t.Errorf("default UA should start with anthropic-sdk-go/, got %q", gotUA)
	}
	if !strings.Contains(gotUA, "ExchangeFederationAssertion") {
		t.Errorf("default UA should name the helper, got %q", gotUA)
	}
}

func TestExchangeFederationAssertion_CustomUserAgent(t *testing.T) {
	var gotUA string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		json.NewEncoder(w).Encode(map[string]any{"access_token": "t"})
	}))
	defer server.Close()

	if _, err := config.ExchangeFederationAssertion(context.Background(), config.FederationExchangeParams{
		Assertion:        "jwt",
		FederationRuleID: "fdrl_1",
		OrganizationID:   "org_1",
		BaseURL:          server.URL,
		UserAgent:        "ant-cli/1.2.3",
	}); err != nil {
		t.Fatal(err)
	}
	if gotUA != "ant-cli/1.2.3" {
		t.Errorf("custom UA not honoured, got %q", gotUA)
	}
}

func TestExchangeFederationAssertion_OmitsServiceAccountWhenEmpty(t *testing.T) {
	var gotBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotBody, _ = io.ReadAll(r.Body)
		json.NewEncoder(w).Encode(map[string]any{"access_token": "t"})
	}))
	defer server.Close()

	if _, err := config.ExchangeFederationAssertion(context.Background(), config.FederationExchangeParams{
		Assertion:        "jwt",
		FederationRuleID: "fdrl_1",
		OrganizationID:   "org_1",
		BaseURL:          server.URL,
	}); err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(gotBody, &decoded); err != nil {
		t.Fatal(err)
	}
	if _, present := decoded["service_account_id"]; present {
		t.Errorf("service_account_id must be omitted when empty, body=%q", gotBody)
	}
	if _, present := decoded["workspace_id"]; present {
		t.Errorf("workspace_id must be omitted when empty, body=%q", gotBody)
	}
}

func TestExchangeFederationAssertion_ErrorSurfacesBodyAndRequestID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Request-Id", "req_oops")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"invalid_grant","error_description":"assertion expired"}`))
	}))
	defer server.Close()

	_, err := config.ExchangeFederationAssertion(context.Background(), config.FederationExchangeParams{
		Assertion:        "jwt",
		FederationRuleID: "fdrl_1",
		OrganizationID:   "org_1",
		BaseURL:          server.URL,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	var fe *config.FederationExchangeError
	if !errors.As(err, &fe) {
		t.Fatalf("expected *FederationExchangeError, got %T: %v", err, err)
	}
	if fe.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode: %d", fe.StatusCode)
	}
	if fe.RequestID != "req_oops" {
		t.Errorf("RequestID: %q", fe.RequestID)
	}
	if !strings.Contains(fe.Body, "assertion expired") {
		t.Errorf("Body missing server text: %q", fe.Body)
	}
	if !strings.Contains(fe.Error(), "req_oops") {
		t.Errorf("Error() should include request id: %q", fe.Error())
	}
}

func TestFederationExchangeError_RedactsSensitiveBody(t *testing.T) {
	err := &config.FederationExchangeError{
		StatusCode: 400,
		RequestID:  "req_test",
		Body:       `{"error":"invalid_request","error_description":"bad","assertion":"eyJhbGciOi...SECRET","refresh_token":"rt-secret","weird_new_field":"also-secret"}`,
	}
	msg := err.Error()
	// Allowlist redaction (RFC 6749 §5.2): drop everything except error,
	// error_description, error_uri. Future fields the server adds are
	// dropped by default rather than relying on a denylist update.
	for _, secret := range []string{"eyJhbGciOi", "SECRET", "rt-secret", "weird_new_field", "also-secret", "assertion", "refresh_token"} {
		if strings.Contains(msg, secret) {
			t.Errorf("federation error message leaks %q: %s", secret, msg)
		}
	}
	if !strings.Contains(msg, `"error":"invalid_request"`) {
		t.Errorf("expected RFC 6749 error key kept: %s", msg)
	}
	if !strings.Contains(msg, `"error_description":"bad"`) {
		t.Errorf("expected RFC 6749 error_description kept: %s", msg)
	}
	if !strings.Contains(msg, "req_test") {
		t.Errorf("expected request id in message: %s", msg)
	}
}

func TestFederationExchangeError_TruncatesLongBody(t *testing.T) {
	// Long JSON body containing an error key; allowlist preserves only that
	// key and its (large) error_description value, which then truncates.
	desc := strings.Repeat("x", config.OAuthErrorBodyMaxLen+200)
	body := `{"error":"server_error","error_description":"` + desc + `"}`
	err := &config.FederationExchangeError{StatusCode: 500, Body: body}
	msg := err.Error()
	if !strings.Contains(msg, "[truncated]") {
		t.Errorf("expected truncation marker: %s", msg)
	}
}

func TestFederationExchangeError_NonJSONBodyRedacted(t *testing.T) {
	// Non-JSON bodies (e.g., HTML 5xx pages from a misconfigured proxy)
	// must not be passed through to the user-facing error message.
	err := &config.FederationExchangeError{StatusCode: 502, Body: "<html>upstream secret leaked here</html>"}
	msg := err.Error()
	if strings.Contains(msg, "secret") {
		t.Errorf("non-JSON body must not pass through: %s", msg)
	}
	if !strings.Contains(msg, "[redacted; not a JSON error response]") {
		t.Errorf("expected non-JSON redaction marker: %s", msg)
	}
}

func TestExchangeFederationAssertion_ValidatesInputs(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name   string
		params config.FederationExchangeParams
	}{
		{"missing assertion", config.FederationExchangeParams{FederationRuleID: "r", OrganizationID: "o"}},
		{"missing rule", config.FederationExchangeParams{Assertion: "a", OrganizationID: "o"}},
		{"missing org", config.FederationExchangeParams{Assertion: "a", FederationRuleID: "r"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := config.ExchangeFederationAssertion(ctx, tc.params); err == nil {
				t.Error("expected error")
			}
		})
	}
}
