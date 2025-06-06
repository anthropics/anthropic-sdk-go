package apierror

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRequestIDFromHeaders simulates a real HTTP server that returns errors with request-id headers
func TestRequestIDFromHeaders(t *testing.T) {
	// Create a test server that simulates the Anthropic API returning an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the request-id header
		w.Header().Set("request-id", "req_test_server_123")
		// Return a 400 error with error JSON
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"type":"invalid_request_error","message":"Invalid request"}}`))
	}))
	defer server.Close()

	// Make a request to the test server
	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, _ := http.DefaultClient.Do(req)

	// Read the response body
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	// Create our API error manually
	apiErr := &Error{
		Request:    req,
		Response:   resp,
		StatusCode: resp.StatusCode,
		RequestID:  resp.Header.Get("request-id"),
	}
	apiErr.UnmarshalJSON(body)

	// Verify the RequestID field
	if apiErr.RequestID != "req_test_server_123" {
		t.Errorf("Expected RequestID to be %s, got %s", "req_test_server_123", apiErr.RequestID)
	}

	// Verify that the error message includes the RequestID
	if !strings.Contains(apiErr.Error(), "Request-ID: req_test_server_123") {
		t.Errorf("Error message should contain request ID: %s", apiErr.Error())
	}
}
