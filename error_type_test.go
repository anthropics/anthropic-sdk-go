package anthropic

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/internal/apierror"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
	"github.com/anthropics/anthropic-sdk-go/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStreamingErrorType(t *testing.T) {
	sseBody := "event: error\ndata: {\"type\":\"error\",\"error\":{\"type\":\"overloaded_error\",\"message\":\"Overloaded\"}}\n\n"
	httpResp := &http.Response{
		StatusCode: 200,
		Header:     http.Header{"Request-Id": []string{"req_stream123"}},
		Body:       io.NopCloser(strings.NewReader(sseBody)),
		Request:    mustNewRequest("POST", "https://api.anthropic.com/v1/messages"),
	}

	stream := ssestream.NewStream[json.RawMessage](ssestream.NewDecoder(httpResp), nil)

	stream.Next()
	require.Error(t, stream.Err())

	var apierr *apierror.Error
	require.True(t, errors.As(stream.Err(), &apierr))
	assert.Equal(t, shared.ErrorTypeOverloadedError, apierr.Type())
	assert.Equal(t, 200, apierr.StatusCode)
	assert.Equal(t, "req_stream123", apierr.RequestID)
	assert.Contains(t, apierr.Error(), "200 OK")
	assert.Contains(t, apierr.Error(), "overloaded_error")
}

func TestStreamingErrorMalformedBody(t *testing.T) {
	sseBody := "event: error\ndata: not valid json\n\n"
	httpResp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(sseBody)),
		Request:    mustNewRequest("POST", "https://api.anthropic.com/v1/messages"),
	}

	stream := ssestream.NewStream[json.RawMessage](ssestream.NewDecoder(httpResp), nil)

	stream.Next()
	require.Error(t, stream.Err())

	// Even with malformed JSON, UnmarshalJSON succeeds (gjson is lenient),
	// so we still get an *apierror.Error — but Type() is empty.
	var apierr *apierror.Error
	require.True(t, errors.As(stream.Err(), &apierr))
	assert.Equal(t, shared.ErrorType(""), apierr.Type())
}

// When the HTTP response has no Request (e.g. a synthetic response in tests),
// the decoder falls back to a plain-text error instead of constructing an
// *apierror.Error, since apierror.Error.Error() requires a non-nil Request.
func TestStreamingErrorNoRequest(t *testing.T) {
	sseBody := "event: error\ndata: {\"type\":\"error\",\"error\":{\"type\":\"overloaded_error\",\"message\":\"Overloaded\"}}\n\n"
	httpResp := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(sseBody)),
	}

	stream := ssestream.NewStream[json.RawMessage](ssestream.NewDecoder(httpResp), nil)

	stream.Next()
	require.Error(t, stream.Err())

	var apierr *apierror.Error
	assert.False(t, errors.As(stream.Err(), &apierr), "expected plain error, not *apierror.Error, when Request is nil")
	assert.Contains(t, stream.Err().Error(), "received error while streaming:")
}

func mustNewRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}
