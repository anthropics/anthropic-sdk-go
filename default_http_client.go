// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"net/http"
	"time"
)

// defaultResponseHeaderTimeout bounds the time between a fully written request
// and the server's response headers. It does not apply to the response body,
// so long-running streams are unaffected. Without this, a server that accepts
// the connection but never responds would hang the request indefinitely.
const defaultResponseHeaderTimeout = 10 * time.Minute

// defaultHTTPClient returns an [*http.Client] used when the caller does not
// supply one via [option.WithHTTPClient]. When [http.DefaultTransport] is the
// stdlib [*http.Transport], it clones it and adds a
// [http.Transport.ResponseHeaderTimeout] so stuck connections fail fast
// instead of compounding across retries. When [http.DefaultTransport] has
// been replaced (for example, wrapped by otelhttp.NewTransport for
// distributed tracing), the wrapped transport is used as-is — preserving the
// caller's instrumentation, with the tradeoff that the default
// ResponseHeaderTimeout does not apply.
func defaultHTTPClient() *http.Client {
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		cloned := t.Clone()
		cloned.ResponseHeaderTimeout = defaultResponseHeaderTimeout
		return &http.Client{Transport: cloned}
	}
	// http.DefaultTransport has been replaced with a wrapped/custom
	// RoundTripper. Preserve the caller's transport rather than panicking.
	return &http.Client{Transport: http.DefaultTransport}
}
