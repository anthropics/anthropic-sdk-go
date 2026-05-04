// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package option

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/anthropics/anthropic-sdk-go/config"
	"github.com/anthropics/anthropic-sdk-go/internal/auth"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/tidwall/sjson"
)

// IdentityTokenFunc returns a fresh JWT identity token (e.g. from SPIFFE/SPIRE,
// the GCP metadata server, or any other OIDC-compatible provider). It is
// invoked once per federation token exchange; implementations should be safe
// for concurrent use.
type IdentityTokenFunc func(ctx context.Context) (string, error)

// FederationOptions configures a federation-based token exchange.
type FederationOptions struct {
	// FederationRuleID identifies the OidcFederationRule governing this
	// exchange. Required; must be a tagged ID with the "fdrl_" prefix.
	FederationRuleID string
	// OrganizationID is the UUID of the Anthropic organization the
	// federation rule belongs to. Required.
	OrganizationID string
	// ServiceAccountID is an optional expected-target check for federation
	// rules with target_type=SERVICE_ACCOUNT. Must be a tagged ID with the
	// "svac_" prefix. Omit for target_type=USER rules, where the principal
	// is derived from the JWT claims.
	ServiceAccountID string
}

// identityTokenFuncAdapter adapts an [IdentityTokenFunc] to the internal
// [auth.IdentityTokenProvider] interface without leaking the internal type.
type identityTokenFuncAdapter struct {
	fn IdentityTokenFunc
}

func (a identityTokenFuncAdapter) GetIdentityToken(ctx context.Context) (string, error) {
	if a.fn == nil {
		return "", fmt.Errorf("option: IdentityTokenFunc is nil")
	}
	return a.fn(ctx)
}

// errOption returns a [RequestOption] whose Apply unconditionally returns
// err. Used to defer validation errors from option constructors so they
// surface on the first request rather than at option-construction time.
func errOption(err error) RequestOption {
	return requestconfig.RequestOptionFunc(func(*requestconfig.RequestConfig) error {
		return err
	})
}

// IdentityTokenFile returns an [IdentityTokenFunc] that reads a JWT from path
// on each invocation. Use it for Kubernetes projected service account tokens
// and similar rotating on-disk credentials — the file is re-read on every
// federation exchange so rotations are picked up automatically. Surrounding
// whitespace is trimmed; an empty file is treated as an error.
//
// Passing this as the provider to [WithFederationTokenProvider] takes
// precedence over the ANTHROPIC_IDENTITY_TOKEN_FILE environment variable:
// explicit options always beat env-var auto-discovery. The env var only
// applies when the client is constructed without a [WithFederationTokenProvider]
// option (or equivalent).
//
// The returned function is safe for concurrent use.
//
// Typical usage:
//
//	client := anthropic.NewClient(
//	    option.WithFederationTokenProvider(
//	        option.IdentityTokenFile("/var/run/secrets/anthropic.com/token"),
//	        option.FederationOptions{
//	            FederationRuleID: "fdrl_...",
//	            OrganizationID:   os.Getenv("ANTHROPIC_ORGANIZATION_ID"),
//	        },
//	    ),
//	)
func IdentityTokenFile(path string) IdentityTokenFunc {
	return (&auth.IdentityTokenFile{Path: path}).GetIdentityToken
}

// WithFederationTokenProvider returns a [RequestOption] that authenticates
// requests using workload identity federation, exchanging a caller-supplied
// identity token for a short-lived Anthropic access token. Use this to
// integrate custom OIDC token sources (SPIFFE/SPIRE, cloud provider SDKs,
// etc.) without staging the token through a file.
//
// provider is called on each token exchange to fetch a fresh JWT.
// opts.FederationRuleID and opts.OrganizationID are required.
//
// The auth middleware is constructed once and reused across requests, so the
// [auth.TokenCache] built by [auth.WithAuthMiddleware] is shared — a fresh
// access token is cached in memory and only re-exchanged when it enters the
// refresh window.
func WithFederationTokenProvider(provider IdentityTokenFunc, opts FederationOptions) RequestOption {
	switch {
	case provider == nil:
		return errOption(fmt.Errorf("option: WithFederationTokenProvider: provider is nil"))
	case opts.FederationRuleID == "":
		return errOption(fmt.Errorf("option: WithFederationTokenProvider: FederationRuleID is required"))
	case opts.OrganizationID == "":
		return errOption(fmt.Errorf("option: WithFederationTokenProvider: OrganizationID is required"))
	}
	tokenProvider := auth.NewOIDCFederationCredentials(auth.OIDCFederationConfig{
		IdentityProvider: identityTokenFuncAdapter{fn: provider},
		FederationRuleID: opts.FederationRuleID,
		OrganizationID:   opts.OrganizationID,
		ServiceAccountID: opts.ServiceAccountID,
	})
	return auth.WithAuthMiddleware(tokenProvider)
}

// RequestOption is an option for the requests made by the anthropic API Client
// which can be supplied to clients, services, and methods. You can read more about this functional
// options pattern in our [README].
//
// [README]: https://pkg.go.dev/github.com/anthropics/anthropic-sdk-go#readme-requestoptions
type RequestOption = requestconfig.RequestOption

// WithBaseURL returns a RequestOption that sets the BaseURL for the client.
//
// For security reasons, ensure that the base URL is trusted.
func WithBaseURL(base string) RequestOption {
	u, err := url.Parse(base)
	if err == nil && u.Path != "" && !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}

	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		if err != nil {
			return fmt.Errorf("requestoption: WithBaseURL failed to parse url %s", err)
		}

		r.BaseURL = u
		return nil
	})
}

// HTTPClient is primarily used to describe an [*http.Client], but also
// supports custom implementations.
//
// For bespoke implementations, prefer using an [*http.Client] with a
// custom transport. See [http.RoundTripper] for further information.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// WithHTTPClient returns a RequestOption that changes the underlying http client used to make this
// request, which by default is [http.DefaultClient].
//
// For custom uses cases, it is recommended to provide an [*http.Client] with a custom
// [http.RoundTripper] as its transport, rather than directly implementing [HTTPClient].
func WithHTTPClient(client HTTPClient) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		if client == nil {
			return fmt.Errorf("requestoption: custom http client cannot be nil")
		}

		if c, ok := client.(*http.Client); ok {
			// Prefer the native client if possible.
			r.HTTPClient = c
			r.CustomHTTPDoer = nil
		} else {
			r.CustomHTTPDoer = client
		}

		return nil
	})
}

// MiddlewareNext is a function which is called by a middleware to pass an HTTP request
// to the next stage in the middleware chain.
type MiddlewareNext = func(*http.Request) (*http.Response, error)

// Middleware is a function which intercepts HTTP requests, processing or modifying
// them, and then passing the request to the next middleware or handler
// in the chain by calling the provided MiddlewareNext function.
type Middleware = func(*http.Request, MiddlewareNext) (*http.Response, error)

// WithMiddleware returns a RequestOption that applies the given middleware
// to the requests made. Each middleware will execute in the order they were given.
func WithMiddleware(middlewares ...Middleware) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.Middlewares = append(r.Middlewares, middlewares...)
		return nil
	})
}

// WithMaxRetries returns a RequestOption that sets the maximum number of retries that the client
// attempts to make. When given 0, the client only makes one request. By
// default, the client retries two times.
//
// WithMaxRetries panics when retries is negative.
func WithMaxRetries(retries int) RequestOption {
	if retries < 0 {
		panic("option: cannot have fewer than 0 retries")
	}
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.MaxRetries = retries
		return nil
	})
}

// WithHeader returns a RequestOption that sets the header value to the associated key. It overwrites
// any value if there was one already present.
func WithHeader(key, value string) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.Request.Header.Set(key, value)
		return nil
	})
}

// WithHeaderAdd returns a RequestOption that adds the header value to the associated key. It appends
// onto any existing values.
func WithHeaderAdd(key, value string) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.Request.Header.Add(key, value)
		return nil
	})
}

// WithHeaderDel returns a RequestOption that deletes the header value(s) associated with the given key.
func WithHeaderDel(key string) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.Request.Header.Del(key)
		return nil
	})
}

// WithQuery returns a RequestOption that sets the query value to the associated key. It overwrites
// any value if there was one already present.
func WithQuery(key, value string) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		query := r.Request.URL.Query()
		query.Set(key, value)
		r.Request.URL.RawQuery = query.Encode()
		return nil
	})
}

// WithQueryAdd returns a RequestOption that adds the query value to the associated key. It appends
// onto any existing values.
func WithQueryAdd(key, value string) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		query := r.Request.URL.Query()
		query.Add(key, value)
		r.Request.URL.RawQuery = query.Encode()
		return nil
	})
}

// WithQueryDel returns a RequestOption that deletes the query value(s) associated with the key.
func WithQueryDel(key string) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		query := r.Request.URL.Query()
		query.Del(key)
		r.Request.URL.RawQuery = query.Encode()
		return nil
	})
}

// WithJSONSet returns a RequestOption that sets the body's JSON value associated with the key.
// The key accepts a string as defined by the [sjson format].
//
// [sjson format]: https://github.com/tidwall/sjson
func WithJSONSet(key string, value any) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) (err error) {
		var b []byte

		if r.Body == nil {
			b, err = sjson.SetBytes(nil, key, value)
			if err != nil {
				return err
			}
		} else if buffer, ok := r.Body.(*bytes.Buffer); ok {
			b = buffer.Bytes()
			b, err = sjson.SetBytes(b, key, value)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("cannot use WithJSONSet on a body that is not serialized as *bytes.Buffer")
		}

		r.Body = bytes.NewBuffer(b)
		return nil
	})
}

// WithJSONDel returns a RequestOption that deletes the body's JSON value associated with the key.
// The key accepts a string as defined by the [sjson format].
//
// [sjson format]: https://github.com/tidwall/sjson
func WithJSONDel(key string) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) (err error) {
		if buffer, ok := r.Body.(*bytes.Buffer); ok {
			b := buffer.Bytes()
			b, err = sjson.DeleteBytes(b, key)
			if err != nil {
				return err
			}
			r.Body = bytes.NewBuffer(b)
			return nil
		}

		return fmt.Errorf("cannot use WithJSONDel on a body that is not serialized as *bytes.Buffer")
	})
}

// WithResponseBodyInto returns a RequestOption that overwrites the deserialization target with
// the given destination. If provided, we don't deserialize into the default struct.
func WithResponseBodyInto(dst any) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.ResponseBodyInto = dst
		return nil
	})
}

// WithResponseInto returns a RequestOption that copies the [*http.Response] into the given address.
func WithResponseInto(dst **http.Response) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.ResponseInto = dst
		return nil
	})
}

// WithRequestBody returns a RequestOption that provides a custom serialized body with the given
// content type.
//
// body accepts an io.Reader or raw []bytes.
func WithRequestBody(contentType string, body any) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		if reader, ok := body.(io.Reader); ok {
			r.Body = reader
			return r.Apply(WithHeader("Content-Type", contentType))
		}

		if b, ok := body.([]byte); ok {
			r.Body = bytes.NewBuffer(b)
			return r.Apply(WithHeader("Content-Type", contentType))
		}

		return fmt.Errorf("body must be a byte slice or implement io.Reader")
	})
}

// WithRequestTimeout returns a RequestOption that sets the timeout for
// each request attempt. This should be smaller than the timeout defined in
// the context, which spans all retries.
func WithRequestTimeout(dur time.Duration) RequestOption {
	// we need this to be a PreRequestOptionFunc so that it can be applied at the endpoint level
	// see: CalculateNonStreamingTimeout
	return requestconfig.PreRequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.RequestTimeout = dur
		return nil
	})
}

// withConfigClientKey identifies the HTTP transport a per-WithConfig
// auth.Middleware should bind to. Mirrors auth.clientKey so the lazily-
// constructed bearer middleware is cached per transport.
type withConfigClientKey struct {
	httpClient *http.Client
	customDoer requestconfig.HTTPDoer
}

// WithConfig returns a [RequestOption] that authenticates requests using
// the given [config.Config].
//
// Credential resolution is deferred to a request-time middleware so a
// later [WithAPIKey] / [WithAuthToken] (tier-1 explicit-credential
// precedence) can preempt the profile — even when the profile's
// credentials are missing or broken. The same escape-hatch shape as
// [explicitProfileErrorOption] applies: if the final RequestConfig has a
// static APIKey/AuthToken, or the request already carries an
// Authorization / X-Api-Key header, the profile is treated as shadowed
// and resolution is skipped entirely.
//
// The config's non-credential fields (BaseURL, WorkspaceID) apply at
// option-Apply time and do not depend on credential resolution. BaseURL
// only takes effect when no explicit [WithBaseURL] is present.
//
// When a static credential shadows the profile, a one-shot warning is
// emitted so users migrating from API keys can tell why their profile is
// being ignored. Use [WithConfigQuiet] to suppress that warning when the
// caller has its own, richer diagnostic.
func WithConfig(cfg *config.Config) RequestOption {
	return withConfig(cfg, false)
}

// WithConfigQuiet is identical to [WithConfig] but suppresses the one-shot
// "static credential shadows profile" warning. Intended for callers (such as
// the CLI) that detect the same condition themselves and emit a more
// specific multi-auth diagnostic, so users don't see two warnings for the
// same thing.
func WithConfigQuiet(cfg *config.Config) RequestOption {
	return withConfig(cfg, true)
}

// WithProfile returns a [RequestOption] that loads the named profile from
// the default config directory (see [config.DefaultDir]) and authenticates
// requests using it — equivalent to setting ANTHROPIC_PROFILE=name and
// constructing a zero-config client. Shorthand for:
//
//	cfg, err := config.LoadProfile(config.DefaultDir(), name)
//	... option.WithConfig(cfg)
//
// If the profile cannot be loaded, the error is deferred to the first
// request and is preempted by a static credential ([WithAPIKey] /
// [WithAuthToken]) — the same escape hatch as the ANTHROPIC_PROFILE env
// path. As with [WithConfig], a static credential shadows the profile at
// request time and a one-shot warning is emitted.
func WithProfile(name string) RequestOption {
	if name == "" {
		return errOption(fmt.Errorf("option: WithProfile: profile name is empty"))
	}
	cfg, err := config.LoadProfile(config.DefaultDir(), name)
	if err != nil {
		loadErr := fmt.Errorf("option: WithProfile(%q): %w", name, err)
		return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
			rc := r
			check := func(req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error) {
				if rc.APIKey != "" || rc.AuthToken != "" {
					return next(req)
				}
				if req.Header.Get("Authorization") != "" || req.Header.Get("X-Api-Key") != "" {
					return next(req)
				}
				return nil, loadErr
			}
			r.Middlewares = append(r.Middlewares, check)
			return nil
		})
	}
	return withConfig(cfg, false)
}

// withoutEnvironmentDefaultsOption is the marker type returned by
// [WithoutEnvironmentDefaults]. Its Apply is a no-op; it is detected by
// [HasWithoutEnvironmentDefaults] before options are applied.
type withoutEnvironmentDefaultsOption struct{}

func (withoutEnvironmentDefaultsOption) Apply(*requestconfig.RequestConfig) error { return nil }

// WithoutEnvironmentDefaults returns a marker [RequestOption] that, when
// passed to anthropic.NewClient, causes it to skip the environment-based
// credential autoload performed by anthropic.DefaultClientOptions
// (ANTHROPIC_API_KEY, ANTHROPIC_AUTH_TOKEN, ANTHROPIC_PROFILE, env
// federation, fallback profile, ANTHROPIC_BASE_URL). The hardcoded
// production base-URL default is still applied so callers that supply only
// credentials get a working client.
//
// Intended for callers that perform their own credential resolution (such
// as the CLI) and want full control over which auth source is used, without
// the SDK's autoloader contributing a second [WithConfig] that would emit a
// duplicate shadow-warning.
func WithoutEnvironmentDefaults() RequestOption {
	return withoutEnvironmentDefaultsOption{}
}

// HasWithoutEnvironmentDefaults reports whether opts contains a
// [WithoutEnvironmentDefaults] marker. Used by anthropic.NewClient to decide
// whether to prepend anthropic.DefaultClientOptions.
func HasWithoutEnvironmentDefaults(opts []RequestOption) bool {
	for _, o := range opts {
		if _, ok := o.(withoutEnvironmentDefaultsOption); ok {
			return true
		}
	}
	return false
}

func withConfig(cfg *config.Config, quiet bool) RequestOption {
	var (
		shadowOnce  sync.Once
		resolveOnce sync.Once
		provider    auth.TokenProvider
		resolveErr  error

		mwMu sync.Mutex
		mwBy = map[withConfigClientKey]auth.Middleware{}
	)
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		// Non-credential config — applied unconditionally so a profile's
		// base_url / workspace_id are honored even when its credentials
		// are shadowed by a later WithAPIKey or fail to resolve.
		if cfg.BaseURL != "" && r.BaseURL == nil {
			if err := WithBaseURL(cfg.BaseURL).Apply(r); err != nil {
				return err
			}
		}
		if cfg.WorkspaceID != "" {
			if err := WithHeader("anthropic-workspace-id", cfg.WorkspaceID).Apply(r); err != nil {
				return err
			}
		}

		// Both checks below are request-time middlewares so they observe
		// the RequestConfig after ALL options have applied — order of
		// WithConfig vs WithAPIKey/WithAuthToken in the caller's option
		// list doesn't matter.
		rc := r
		if !quiet {
			shadowCheck := func(req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error) {
				shadowOnce.Do(func() {
					switch {
					case rc.APIKey != "":
						auth.WarnConfigShadowed("ANTHROPIC_API_KEY", detectShadowSource("ANTHROPIC_API_KEY"))
					case rc.AuthToken != "":
						auth.WarnConfigShadowed("ANTHROPIC_AUTH_TOKEN", detectShadowSource("ANTHROPIC_AUTH_TOKEN"))
					}
				})
				return next(req)
			}
			r.Middlewares = append(r.Middlewares, shadowCheck)
		}

		credCheck := func(req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error) {
			if rc.APIKey != "" || rc.AuthToken != "" {
				return next(req)
			}
			if req.Header.Get("Authorization") != "" || req.Header.Get("X-Api-Key") != "" {
				return next(req)
			}
			resolveOnce.Do(func() {
				result, err := auth.ResolveCredentials(cfg)
				if err != nil {
					resolveErr = err
					return
				}
				provider = result.Provider
			})
			if resolveErr != nil {
				return nil, resolveErr
			}

			handler := rc.HTTPClient.Do
			if rc.CustomHTTPDoer != nil {
				handler = rc.CustomHTTPDoer.Do
			}
			key := withConfigClientKey{httpClient: rc.HTTPClient, customDoer: rc.CustomHTTPDoer}
			mwMu.Lock()
			inner, ok := mwBy[key]
			if !ok {
				inner = auth.NewProviderMiddleware(provider, handler)
				mwBy[key] = inner
			}
			mwMu.Unlock()
			return inner(req, next)
		}
		r.Middlewares = append(r.Middlewares, credCheck)
		return nil
	})
}

// detectShadowSource classifies where the shadowing static credential
// came from by inspecting whether the named env var is set. When the env
// var has a non-empty value the source is [auth.ConfigShadowFromEnv];
// otherwise the credential must have come from an explicit option call.
func detectShadowSource(envVar string) auth.ConfigShadowSource {
	if v, ok := os.LookupEnv(envVar); ok && v != "" {
		return auth.ConfigShadowFromEnv
	}
	return auth.ConfigShadowFromOption
}

// WithEnvironmentProduction returns a RequestOption that sets the current
// environment to be the "production" environment. An environment specifies which base URL
// to use by default.
func WithEnvironmentProduction() RequestOption {
	return requestconfig.WithDefaultBaseURL("https://api.anthropic.com/")
}

// WithAPIKey returns a RequestOption that sets the client setting "api_key".
func WithAPIKey(value string) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.APIKey = value
		return r.Apply(WithHeader("X-Api-Key", r.APIKey))
	})
}

// WithAuthToken returns a RequestOption that sets the client setting "auth_token".
func WithAuthToken(value string) RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.AuthToken = value
		return r.Apply(WithHeader("authorization", fmt.Sprintf("Bearer %s", r.AuthToken)))
	})
}
