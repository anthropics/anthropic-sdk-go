// EDIT(begin): add custom options for JSON encoding
package json

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

// PROTOTYPE(begin): opt-in buffer-direct encoding of nested SDK marshalers.
// When set, the encoder encodes any value implementing [MarshalerTo] straight
// into the shared output buffer instead of calling MarshalJSON() (which
// allocates a fresh []byte of the whole subtree at every nesting level).
func WithBufferDirect(b bool) Option {
	return func(eos *encOpts) {
		eos.bufferDirect = b
	}
}

// PROTOTYPE(end)

func (eos encOpts) apply(opts ...Option) encOpts {
	for _, opt := range opts {
		opt(&eos)
	}
	return eos
}

// EDIT(end)
