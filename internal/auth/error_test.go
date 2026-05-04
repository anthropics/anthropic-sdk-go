package auth

import (
	"errors"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/config"
)

func TestNoCredentialsError_IsSentinel(t *testing.T) {
	err := &NoCredentialsError{}
	if !errors.Is(err, ErrNoCredentials) {
		t.Fatal("errors.Is(NoCredentialsError, ErrNoCredentials) should be true")
	}
	wrapped := errors.Join(errors.New("leading"), err)
	if !errors.Is(wrapped, ErrNoCredentials) {
		t.Fatal("errors.Is should see through wrappers")
	}
	other := errors.New("unrelated")
	if errors.Is(err, other) {
		t.Fatal("errors.Is should not match unrelated errors")
	}
}

func TestOAuthTokenError_RedactsSensitiveJSONKeys(t *testing.T) {
	err := &OAuthTokenError{
		StatusCode: 400,
		Body:       `{"error":"invalid_grant","access_token":"sk-ant-oat01-secret","refresh_token":"rt-secret","assertion":"eyJhbGciOi..."}`,
	}
	msg := err.Error()
	for _, secret := range []string{"sk-ant-oat01-secret", "rt-secret", "eyJhbGciOi"} {
		if strings.Contains(msg, secret) {
			t.Errorf("error message contains unredacted secret %q: %s", secret, msg)
		}
	}
	if !strings.Contains(msg, "invalid_grant") {
		t.Errorf("expected non-sensitive fields preserved: %s", msg)
	}
}

// TestOAuthTokenError_RedactsJSONWithoutRFCCode verifies that JSON bodies
// without an RFC 6749 `error` field are redacted to an empty allowed-keys
// object. The previous denylist behavior would have leaked any non-listed
// fields; the allowlist drops them by default.
func TestOAuthTokenError_RedactsJSONWithoutRFCCode(t *testing.T) {
	err := &OAuthTokenError{
		StatusCode: 500,
		Body:       `{"message":"server error","access_token":"sk-ant-oat01-secret","refresh_token":"rt-secret"}`,
	}
	msg := err.Error()
	for _, secret := range []string{"sk-ant-oat01-secret", "rt-secret", "server error"} {
		if strings.Contains(msg, secret) {
			t.Errorf("error message contains unredacted secret %q: %s", secret, msg)
		}
	}
}

func TestOAuthTokenError_TruncatesLongBody(t *testing.T) {
	// Bodies WITH an RFC `error` key bypass redaction (formatted directly),
	// so to exercise the truncation path we use a body the RFC parser
	// rejects (no `error` key) but that the allowlist preserves a long
	// allowed value from. This is the realistic shape for a server that
	// responds with `{"error_uri":"..."}` alone.
	uri := strings.Repeat("x", config.OAuthErrorBodyMaxLen+200)
	body := `{"error_uri":"` + uri + `"}`
	err := &OAuthTokenError{StatusCode: 500, Body: body}
	msg := err.Error()
	if !strings.Contains(msg, "[truncated]") {
		t.Errorf("expected truncation marker: %s", msg[len(msg)-40:])
	}
	// Allow ~256 bytes of OAuthTokenError formatting overhead.
	if len(msg) > config.OAuthErrorBodyMaxLen+256 {
		t.Errorf("message too long after truncation: %d bytes", len(msg))
	}
}

func TestOAuthTokenError_NonJSONBodyRedacted(t *testing.T) {
	// A non-JSON body (HTML 5xx page, plain text from a misbehaving proxy)
	// must not pass through to the error message. Replaces the previous
	// denylist behavior of "if not JSON, just include the raw body."
	err := &OAuthTokenError{StatusCode: 503, Body: "upstream secret leaked here"}
	msg := err.Error()
	if strings.Contains(msg, "secret") {
		t.Errorf("non-JSON body must not pass through: %s", msg)
	}
	if !strings.Contains(msg, "[redacted; not a JSON error response]") {
		t.Errorf("expected non-JSON redaction marker: %s", msg)
	}
}

func TestOAuthTokenError_JSONNullBodyRedactsToEmpty(t *testing.T) {
	// A `null` body is valid JSON: redaction must not classify it as
	// non-JSON. The filter loop on a nil map serializes to "{}".
	err := &OAuthTokenError{StatusCode: 400, Body: "null"}
	msg := err.Error()
	if strings.Contains(msg, "[redacted; not a JSON error response]") {
		t.Errorf("`null` is valid JSON; must not get non-JSON marker: %s", msg)
	}
	if !strings.Contains(msg, "{}") {
		t.Errorf("expected empty-object redaction for null body: %s", msg)
	}
}

// TestOAuthTokenError_InvalidGrantSuggestsRelogin checks that a refresh-token
// rejection (400 with error=invalid_grant) surfaces the RFC 6749 error_description
// and a remediation hint telling the user to re-authenticate.
func TestOAuthTokenError_InvalidGrantSuggestsRelogin(t *testing.T) {
	err := &OAuthTokenError{
		StatusCode: 400,
		Body:       `{"error":"invalid_grant","error_description":"refresh token expired"}`,
	}
	msg := err.Error()
	if !strings.Contains(msg, "refresh token expired") {
		t.Errorf("expected error_description in message, got: %s", msg)
	}
	if !strings.Contains(msg, "anthropic auth login") {
		t.Errorf("expected relogin hint in message, got: %s", msg)
	}
	if !strings.Contains(msg, "invalid_grant") {
		t.Errorf("expected error code in message, got: %s", msg)
	}
}

// TestOAuthTokenError_401UnauthorizedSuggestsRelogin checks that a 401 with any
// body still triggers the re-login hint because a token-endpoint 401 almost
// always means the refresh credential has been revoked.
func TestOAuthTokenError_401UnauthorizedSuggestsRelogin(t *testing.T) {
	err := &OAuthTokenError{StatusCode: 401, Body: `{"error":"unauthorized_client"}`}
	msg := err.Error()
	if !strings.Contains(msg, "anthropic auth login") {
		t.Errorf("expected relogin hint on 401, got: %s", msg)
	}
}

// TestOAuthTokenError_500DoesNotSuggestRelogin verifies that transient server
// errors do NOT blame the user's credentials.
func TestOAuthTokenError_500DoesNotSuggestRelogin(t *testing.T) {
	err := &OAuthTokenError{StatusCode: 500, Body: `{"error":"server_error"}`}
	msg := err.Error()
	if strings.Contains(msg, "anthropic auth login") {
		t.Errorf("5xx should not suggest re-login: %s", msg)
	}
}

// TestOAuthTokenError_RequestIDInMessage verifies the HTTP Request-Id response
// header is captured into the error and printed so support can correlate.
func TestOAuthTokenError_RequestIDInMessage(t *testing.T) {
	err := &OAuthTokenError{
		StatusCode: 400,
		Body:       `{"error":"invalid_grant"}`,
		RequestID:  "req_abc123",
	}
	msg := err.Error()
	if !strings.Contains(msg, "req_abc123") {
		t.Errorf("expected Request-Id in message, got: %s", msg)
	}
}

// TestOAuthTokenError_InvalidGrantStillRedactsSecrets makes sure the new
// RFC 6749 parsing path does not accidentally bypass body redaction.
func TestOAuthTokenError_InvalidGrantStillRedactsSecrets(t *testing.T) {
	err := &OAuthTokenError{
		StatusCode: 400,
		Body:       `{"error":"invalid_grant","refresh_token":"rt-secret","assertion":"eyJhbGciOi.."}`,
	}
	msg := err.Error()
	if strings.Contains(msg, "rt-secret") {
		t.Errorf("refresh_token leaked: %s", msg)
	}
	if strings.Contains(msg, "eyJhbGciOi..") {
		t.Errorf("assertion leaked: %s", msg)
	}
}
