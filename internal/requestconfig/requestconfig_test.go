package requestconfig

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/internal/apierror"
)

// mockBaseURL is a helper function to create a URL for testing
func mockBaseURL(server *httptest.Server) *url.URL {
	u, _ := url.Parse(server.URL)
	return u
}

// TestErrorWithRequestID tests that RequestID is properly extracted from response headers
// and included in the Error struct when API errors occur
func TestErrorWithRequestID(t *testing.T) {
	// Create a test server that simulates an API returning an error with request-id
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the request-id header - using the same header name as in the code
		w.Header().Set("request-id", "req_123456789")
		// Return a 400 error with error JSON
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"type":"invalid_request_error","message":"Invalid request"}}`))
	}))
	defer server.Close()

	// Create a request to the test server
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a RequestConfig with proper BaseURL
	cfg := &RequestConfig{
		Context:    context.Background(),
		Request:    req,
		BaseURL:    mockBaseURL(server),
		HTTPClient: http.DefaultClient,
	}

	// Execute the request, which should return an error
	err = cfg.Execute()
	if err == nil {
		t.Fatal("Expected an error, but got nil")
	}

	// The error should be of type *apierror.Error
	apiErr, ok := err.(*apierror.Error)
	if !ok {
		t.Fatalf("Expected error of type *apierror.Error, got %T", err)
	}

	// Verify that RequestID field was properly set from the header
	expectedRequestID := "req_123456789"
	if apiErr.RequestID != expectedRequestID {
		t.Errorf("Expected RequestID to be %s, got %s", expectedRequestID, apiErr.RequestID)
	}

	// Verify that the error message includes the RequestID
	errorMsg := apiErr.Error()
	if !strings.Contains(errorMsg, "Request-ID: req_123456789") {
		t.Errorf("Error message should contain request ID, got: %s", errorMsg)
	}
}
