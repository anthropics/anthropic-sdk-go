package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go/config"
)

// ErrNoCredentials is the sentinel for the case where the default credential
// chain produced no usable credential. Every [NoCredentialsError] value
// matches this sentinel via [errors.Is], so callers can branch on the
// class of failure without a type assertion.
var ErrNoCredentials = errors.New("no Anthropic credentials found")

// OAuthTokenError is returned when an OAuth token request (workload identity
// exchange or authorized_user refresh) fails with a non-200 status.
type OAuthTokenError struct {
	StatusCode int
	Body       string
	// RequestID is the value of the server's Request-Id response header, if
	// any. Empty when the response had no such header. Included in the
	// formatted error so support can correlate user reports with server logs.
	RequestID string
}

// rfc6749ErrorBody is the subset of the RFC 6749 OAuth error response the SDK
// surfaces in user-facing error messages. Extra fields are ignored.
type rfc6749ErrorBody struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (e *OAuthTokenError) Error() string {
	parts := []string{fmt.Sprintf("oauth token request failed (status %d)", e.StatusCode)}
	if e.RequestID != "" {
		parts = append(parts, fmt.Sprintf("request id %s", e.RequestID))
	}

	var parsed rfc6749ErrorBody
	if err := json.Unmarshal([]byte(e.Body), &parsed); err == nil && parsed.Error != "" {
		if parsed.ErrorDescription != "" {
			parts = append(parts, fmt.Sprintf("%s: %s", parsed.Error, parsed.ErrorDescription))
		} else {
			parts = append(parts, parsed.Error)
		}
	} else {
		parts = append(parts, config.RedactOAuthErrorBody(e.Body))
	}

	msg := strings.Join(parts, "; ")
	if shouldSuggestRelogin(e.StatusCode, parsed.Error) {
		msg += " — re-run `anthropic auth login` to re-authenticate"
	}
	return msg
}

// shouldSuggestRelogin returns true when a token-endpoint failure is most
// likely caused by an expired or revoked refresh credential. 5xx is treated
// as transient; blaming the user's credentials on a server outage is wrong.
func shouldSuggestRelogin(status int, rfcCode string) bool {
	if rfcCode == "invalid_grant" {
		return true
	}
	return status == 400 || status == 401 || status == 403
}

// CredentialSourceState categorizes the outcome of one entry in the default
// credential chain when no source produced a usable credential.
type CredentialSourceState int

const (
	CredentialSourceNotSet CredentialSourceState = iota
	CredentialSourceNotFound
	CredentialSourceLoadFailed
	CredentialSourcePartial
)

// CredentialSourceStatus records what the default resolution chain observed
// for one credential source when all sources failed. Detail is the load
// error for [CredentialSourceLoadFailed], the remediation hint for
// [CredentialSourceNotFound], or the missing-var list for
// [CredentialSourcePartial]; it is empty for [CredentialSourceNotSet].
type CredentialSourceStatus struct {
	Name   string
	State  CredentialSourceState
	Detail string
}

// NoCredentialsError is returned by the default credential chain when no
// source produced a usable credential. Its Error method renders every
// source the SDK tried and a brief remediation hint.
type NoCredentialsError struct {
	Sources []CredentialSourceStatus
}

// Is reports whether target is [ErrNoCredentials]. Enables callers to
// detect the no-credentials class with [errors.Is] in addition to
// [errors.As] against the concrete type.
func (e *NoCredentialsError) Is(target error) bool {
	return target == ErrNoCredentials
}

func (e *NoCredentialsError) Error() string {
	var b strings.Builder
	b.WriteString("no Anthropic credentials found. The SDK tried these sources in order:")
	for i, s := range e.Sources {
		fmt.Fprintf(&b, "\n  %d. %s: %s", i+1, s.Name, formatSourceState(s))
	}
	b.WriteString("\nTo fix:")
	b.WriteString("\n  - run `anthropic auth login` to interactively authenticate, or")
	b.WriteString("\n  - set ANTHROPIC_API_KEY, or")
	b.WriteString("\n  - configure workload identity federation with ANTHROPIC_FEDERATION_RULE_ID, ANTHROPIC_ORGANIZATION_ID, and ANTHROPIC_IDENTITY_TOKEN_FILE")
	return b.String()
}

func formatSourceState(s CredentialSourceStatus) string {
	switch s.State {
	case CredentialSourceNotFound:
		if s.Detail != "" {
			return "not found (" + s.Detail + ")"
		}
		return "not found"
	case CredentialSourceLoadFailed:
		return "load failed: " + s.Detail
	case CredentialSourcePartial:
		return s.Detail
	default:
		return "not set"
	}
}

// CredentialResolutionError is returned when a credential provider cannot be
// constructed from the available configuration.
type CredentialResolutionError struct {
	Message string
	Err     error
}

func (e *CredentialResolutionError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("credential resolution error: %s: %v", e.Message, e.Err)
	}
	return fmt.Sprintf("credential resolution error: %s", e.Message)
}

func (e *CredentialResolutionError) Unwrap() error {
	return e.Err
}
