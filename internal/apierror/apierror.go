// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

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
	StatusCode int
	Request    *http.Request
	Response   *http.Response
	RequestID  string

	errorType shared.ErrorType
}

// Type returns the error type from the API response body, e.g.
// "rate_limit_error" or "overloaded_error". Returns "" if the
// response body did not contain a recognized error type.
func (r *Error) Type() shared.ErrorType { return r.errorType }

// Returns the unmodified JSON received from the API
func (r Error) RawJSON() string { return r.JSON.raw }
func (r *Error) UnmarshalJSON(data []byte) error {
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

func (r *Error) Error() string {
	// Attempt to re-populate the response body
	statusInfo := fmt.Sprintf("%s %q: %d %s", r.Request.Method, r.Request.URL, r.Response.StatusCode, http.StatusText(r.Response.StatusCode))

	if r.RequestID != "" {
		statusInfo += fmt.Sprintf(" (Request-ID: %s)", r.RequestID)
	}

	return fmt.Sprintf("%s %s", statusInfo, r.JSON.raw)
}

func (r *Error) DumpRequest(body bool) []byte {
	if r.Request.GetBody != nil {
		r.Request.Body, _ = r.Request.GetBody()
	}
	out, _ := httputil.DumpRequestOut(r.Request, body)
	return out
}

func (r *Error) DumpResponse(body bool) []byte {
	out, _ := httputil.DumpResponse(r.Response, body)
	return out
}
