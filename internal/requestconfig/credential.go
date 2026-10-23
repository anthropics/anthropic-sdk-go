package requestconfig

import "fmt"

// A client has one credential, and this file is where that is enforced.
//
// The SDK keeps no separate list of client options. The environment defaults, the
// arguments of anthropic.NewClient and the arguments of a method reach
// [NewRequestConfig] as one flat list, once for every request:
//
//	defaults | NewClient's arguments | the method's arguments
//
// Two cases that look the same in that list have different outcomes:
//
//   - Two different credential options in one call are a mistake, so the request
//     fails.
//   - A credential option in a later call replaces the one from an earlier call.
//     That is how a method call overrides the client's credential.
//
// To tell them apart, NewClient puts a [CallBoundary] at each | above.
//
// [resolveCredential] settles both cases on the list, before any option is
// applied. A replaced credential is never applied, so nothing has to be undone:
// the request is what it would be had that option not been passed. The option
// that is left keeps its position, so a WithHeader or a WithMiddleware before or
// after it behaves as it does next to any other option.
//
// The last boundary ends the client's options. When the credential came before it,
// the base URL in effect there is recorded as [RequestConfig.ClientBaseURL]: what a
// method's WithBaseURL overrides.

// CredentialOption is a [RequestOption] that supplies the credential requests are
// sent with. Of the CredentialOptions given to [NewRequestConfig], only the last
// supplies it. Each has to be in the list itself or in a [JoinedOption]: one that
// another option's function applies is not seen.
type CredentialOption struct {
	// Name is the option's name in the option package, for error messages.
	Name string
	// Settings applies what a profile holds besides its credential: a base URL and
	// a workspace. They stay in effect when an option without Settings replaces
	// the profile's credential, so that a request can swap the token and still
	// reach the profile's deployment. It is nil when the option has none.
	Settings   RequestOption
	Credential RequestOption
}

func (o CredentialOption) Apply(r *RequestConfig) error {
	if o.Settings == nil {
		return o.Credential.Apply(r)
	}
	return r.Apply(o.Settings, o.Credential)
}

// CallBoundary is a [RequestOption] that separates the options passed in one call
// from the options passed in another. It does nothing when applied:
// [resolveCredential] looks for it in the list.
type CallBoundary struct{}

func (CallBoundary) Apply(*RequestConfig) error { return nil }

// resolveCredential returns opts without the [CredentialOption]s that a later one
// replaces. It fails when one call passed two different ones.
func resolveCredential(opts []RequestOption) ([]RequestOption, error) {
	opts = flatten(nil, opts)

	// fromClient is set at the first boundary, which anthropic.NewClient puts
	// before its arguments. Options given straight to a service constructor, such
	// as anthropic.NewMessageService, have no boundaries. There the service's
	// options and a method's cannot be told apart, so nothing is an error and the
	// last CredentialOption wins.
	fromClient := false
	// inCall names the CredentialOption seen since the last boundary. The same
	// one twice is allowed, as it is for any other option.
	inCall := ""
	// credential and settings index the last CredentialOption and the last one
	// that has Settings. boundary indexes the last CallBoundary.
	credential, settings, boundary := -1, -1, -1

	for i, o := range opts {
		switch o := o.(type) {
		case CallBoundary:
			fromClient, inCall, boundary = true, "", i
		case CredentialOption:
			if fromClient && inCall != "" && inCall != o.Name {
				return nil, fmt.Errorf(
					"option: pass at most one of WithAPIKey, WithAuthToken, WithConfig, WithConfigQuiet, WithProfile and WithFederationTokenProvider; got %s and %s",
					inCall, o.Name,
				)
			}
			inCall, credential = o.Name, i
			if o.Settings != nil {
				settings = i
			}
		}
	}

	// flatten returned a new slice, so filtering it in place leaves the caller's alone.
	resolved := opts[:0]
	for i, o := range opts {
		c, ok := o.(CredentialOption)
		switch {
		case i == boundary && credential < boundary:
			resolved = append(resolved, RequestOptionFunc(recordClientBaseURL))
		case !ok || i == credential:
			resolved = append(resolved, o)
		case i == settings:
			resolved = append(resolved, c.Settings)
		}
	}
	return resolved, nil
}

func recordClientBaseURL(r *RequestConfig) error {
	r.ClientBaseURL = r.BaseURL
	if r.ClientBaseURL == nil {
		r.ClientBaseURL = r.DefaultBaseURL
	}
	return nil
}

// JoinedOption is a [RequestOption] that applies its members in order. It is a
// distinct type, not a closure, so that the members can be looked at.
type JoinedOption []RequestOption

func (opts JoinedOption) Apply(r *RequestConfig) error {
	return r.Apply(opts...)
}

// flatten appends opts to dst, with every [JoinedOption] replaced by its members.
func flatten(dst, opts []RequestOption) []RequestOption {
	for _, o := range opts {
		if joined, ok := o.(JoinedOption); ok {
			dst = flatten(dst, joined)
		} else {
			dst = append(dst, o)
		}
	}
	return dst
}

// HasCredentialOption reports whether opts contains a [CredentialOption], looking
// inside a [JoinedOption] at any depth.
func HasCredentialOption(opts []RequestOption) bool {
	for _, o := range flatten(nil, opts) {
		if _, ok := o.(CredentialOption); ok {
			return true
		}
	}
	return false
}
