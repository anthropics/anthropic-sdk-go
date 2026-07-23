// EDIT(begin): add custom options for JSON encoding
package json

import (
	stdjson "encoding/json"
	"reflect"

	"github.com/anthropics/anthropic-sdk-go/internal/encoding/json/shims"
)

type Option func(*encOpts)

// Every time a sub-type of [json.Marshaler] is encountered,
// skip a redundant and costly compaction step, trust it to self-compact.
//
// This is a divergence from the standard library behavior, and is only guaranteed
// safe with SDK types.
func WithSkipCompaction(b bool) Option {
	return func(eos *encOpts) {
		eos.skipCompaction = true
	}
}

func (eos encOpts) apply(opts ...Option) encOpts {
	for _, opt := range opts {
		opt(&eos)
	}
	return eos
}

var rawMessageType = shims.TypeFor[stdjson.RawMessage]()

// rawMessageEncoder is marshalerEncoder minus the [WithSkipCompaction] trust:
// a json.RawMessage returns caller-supplied bytes verbatim, so they always get
// the scanning pass, which compacts, HTML-escapes, and rejects invalid JSON.
func rawMessageEncoder(e *encodeState, v reflect.Value, opts encOpts) {
	opts.skipCompaction = false
	marshalerEncoder(e, v, opts)
}

// EDIT(end)
