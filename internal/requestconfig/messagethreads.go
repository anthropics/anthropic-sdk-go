package requestconfig

import "slices"

// withoutMessageThreads does nothing when applied: [SupportsMessageThreads] looks
// for it in the list.
type withoutMessageThreads struct{}

func (withoutMessageThreads) Apply(*RequestConfig) error { return nil }

// WithoutMessageThreads returns a RequestOption that marks the client as one
// whose API has no message threads.
func WithoutMessageThreads() RequestOption {
	return withoutMessageThreads{}
}

// SupportsMessageThreads reports whether a client built from opts can use
// message threads. It looks inside a [JoinedOption] at any depth.
func SupportsMessageThreads(opts ...RequestOption) bool {
	return !slices.Contains(flatten(nil, opts), RequestOption(withoutMessageThreads{}))
}
