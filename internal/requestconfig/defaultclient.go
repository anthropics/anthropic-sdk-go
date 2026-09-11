package requestconfig

import (
	"net/http"
	"time"
)

// defaultResponseHeaderTimeout bounds the time between a fully written request
// and the server's response headers. It does not apply to the response body,
// so long-running streams are unaffected. Without this, a server that accepts
// the connection but never responds would hang the request indefinitely.
const defaultResponseHeaderTimeout = 10 * time.Minute

// DefaultHTTPClient returns a new [http.Client] for clients that were not given
// one with option.WithHTTPClient: a clone of [http.DefaultTransport] with a
// response-header timeout, or [http.DefaultTransport] itself if it has been
// replaced by a wrapper such as otelhttp.
func DefaultHTTPClient() *http.Client {
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		t = t.Clone()
		t.ResponseHeaderTimeout = defaultResponseHeaderTimeout
		return &http.Client{Transport: t}
	}
	return &http.Client{Transport: http.DefaultTransport}
}
