package apierror

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared"
)

// Error represents an error that originates from the API, i.e. when a request is
// made and the API returns a response with a HTTP status code. Other errors are
// not wrapped by this SDK.
type Error struct {
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
	StatusCode  int
	Request     *http.Request
	Response    *http.Response
	RequestID   string
	WorkspaceID string

	errorType shared.ErrorType
}

// Type returns the error type from the API response body, e.g.
// "rate_limit_error" or "overloaded_error". Returns "" if the
// response body did not contain a recognized error type.
func (r *Error) Type() shared.ErrorType { return r.errorType }

// Returns the unmodified JSON received from the API
func (r Error) RawJSON() string { return r.JSON.raw }

func (r *Error) UnmarshalJSON(data []byte) error {
	if !json.Valid(data) {
		r.JSON.raw = string(data)
		return nil
	}
	if err := apijson.UnmarshalRoot(data, r); err != nil {
		return err
	}
	// Extract error type from the standard {"error":{"type":"..."}} envelope.
	var envelope struct {
		Error struct {
			Type shared.ErrorType `json:"type"`
		} `json:"error"`
	}
	if json.Unmarshal(data, &envelope) == nil {
		r.errorType = envelope.Error.Type
	}
	return nil
}

// UnmarshalAPIJSON marks Error as not apijson-native. See
// [apijson.CustomUnmarshaler].
func (r *Error) UnmarshalAPIJSON(data []byte) error { return r.UnmarshalJSON(data) }

func (r *Error) Error() string {
	msg := fmt.Sprintf("%s %q: %d %s", r.Request.Method, r.Request.URL, r.Response.StatusCode, http.StatusText(r.Response.StatusCode))
	if r.RequestID != "" {
		msg += fmt.Sprintf(" (Request-ID: %s)", r.RequestID)
	}
	if r.WorkspaceID != "" {
		msg += fmt.Sprintf(" (Workspace-ID: %s)", r.WorkspaceID)
	}
	if body := r.JSON.raw; body != "" {
		msg += " " + body
	}
	return msg
}

// DumpRequest returns the request in HTTP/1.x wire format, with sensitive header
// values redacted.
func (r *Error) DumpRequest(body bool) []byte {
	if r.Request.GetBody != nil {
		r.Request.Body, _ = r.Request.GetBody()
	}
	req := *r.Request
	req.Header = RedactHeaders(r.Request.Header)
	out, _ := httputil.DumpRequestOut(&req, body)
	// The dump may swap in an unread copy of the body; keep it for later reads.
	r.Request.Body = req.Body
	return out
}

// DumpResponse returns the response in HTTP/1.x wire format, with sensitive header
// values redacted.
func (r *Error) DumpResponse(body bool) []byte {
	resp := *r.Response
	resp.Header = RedactHeaders(r.Response.Header)
	out, _ := httputil.DumpResponse(&resp, body)
	// The dump may swap in an unread copy of the body; keep it for later reads.
	r.Response.Body = resp.Body
	return out
}

var sensitiveHeaders = []string{"authorization", "api-key", "x-api-key", "cookie", "set-cookie"}

// RedactHeaders returns headers with every value of each sensitive header
// replaced by "***", without modifying headers. The result may be headers
// itself, so callers must not modify it.
func RedactHeaders(headers http.Header) http.Header {
	var redacted http.Header
	for _, name := range sensitiveHeaders {
		values := headers.Values(name)
		if len(values) == 0 {
			continue
		}
		if redacted == nil {
			redacted = headers.Clone()
		}
		redacted.Del(name)
		for range values {
			redacted.Add(name, "***")
		}
	}
	if redacted == nil {
		return headers
	}
	return redacted
}
