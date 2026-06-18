// Package stainlessheader is the single source of truth for the
// x-stainless-helper telemetry header — the key, the closed set of tag
// values shared across SDKs, and the append-don't-clobber request option.
package stainlessheader

import (
	"net/http"
	"strings"

	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
)

// Header is the helper-telemetry header key. Always this lowercase form;
// http.Header is case-insensitive but a single canonical casing keeps every
// call site greppable.
const Header = "x-stainless-helper"

// Value is the closed set of helper-telemetry tags, shared verbatim across
// SDKs. A typo at a call site is caught by the unused-const check rather than
// silently mistagged. Existing values keep their original spellings — telemetry
// consumers match on them, so renames lose history. New tags are hyphenated
// lowercase; add them here (and to the matching set in every other SDK) before
// using them.
type Value string

const (
	BetaToolRunner            Value = "BetaToolRunner"
	Compaction                Value = "compaction"
	EnvironmentsWorkPoller    Value = "environments-work-poller"
	EnvironmentsWorker        Value = "environments-worker"
	FallbackRefusalMiddleware Value = "fallback-refusal-middleware"
	SessionToolRunner         Value = "session-tool-runner"
)

// With returns a request option (assignable to [option.RequestOption]) that
// appends value to the x-stainless-helper header rather than replacing it —
// the backend logs the header as one opaque string, so a second header line or
// a clobbered value loses data. Existing tags keep their position; the new tag
// appends at the end; duplicates are dropped.
func With(value Value) requestconfig.RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		AppendHeaderValue(r.Request.Header, value)
		return nil
	})
}

// AppendHeaderValue appends value to h's x-stainless-helper entry in place,
// collapsing any existing values (under any casing, across multiple lines)
// into one comma-joined deduplicated string. Exposed for code paths that hold
// a [*http.Request] rather than a request option chain (e.g. middleware).
func AppendHeaderValue(h http.Header, value Value) {
	var tokens []string
	seen := map[string]struct{}{}
	add := func(tok string) {
		if tok == "" {
			return
		}
		if _, ok := seen[tok]; ok {
			return
		}
		seen[tok] = struct{}{}
		tokens = append(tokens, tok)
	}
	for _, existing := range h.Values(Header) {
		for _, tok := range strings.Split(existing, ",") {
			add(strings.TrimSpace(tok))
		}
	}
	add(string(value))
	h.Set(Header, strings.Join(tokens, ", "))
}
