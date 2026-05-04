// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go/config"
	"github.com/anthropics/anthropic-sdk-go/internal/auth"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// Client creates a struct with services and top level methods that help with
// interacting with the anthropic API. You should not instantiate this client
// directly, and instead use the [NewClient] method instead.
type Client struct {
	Options     []option.RequestOption
	Completions CompletionService
	Messages    MessageService
	Models      ModelService
	Beta        BetaService
}

// DefaultClientOptions walks the default credential chain per the
// cross-SDK credential precedence spec:
//
//  1. ANTHROPIC_API_KEY
//  2. ANTHROPIC_AUTH_TOKEN
//  3. Explicit profile via ANTHROPIC_PROFILE (surfaces the error if the
//     named profile is missing — the user explicitly selected it)
//  4. Env-var federation (ANTHROPIC_FEDERATION_RULE_ID +
//     ANTHROPIC_ORGANIZATION_ID + ANTHROPIC_IDENTITY_TOKEN_FILE / _TOKEN)
//  5. Fallback profile (active_config file or literal "default" — a
//     quiet miss when absent, so a WIF-configured machine with a
//     leftover default profile still uses WIF)
//
// When no source produces a credential, the first request fails with an
// [auth.NoCredentialsError]. If ANTHROPIC_PROFILE points at a missing or
// invalid profile, the first request instead fails with a wrapped
// profile-load error naming the profile. An explicit credential option
// passed to [NewClient] (e.g. [option.WithAPIKey] or [option.WithAuthToken])
// suppresses both paths. Also honors ANTHROPIC_BASE_URL.
func DefaultClientOptions() []option.RequestOption {
	defaults := []option.RequestOption{
		option.WithHTTPClient(defaultHTTPClient()),
		option.WithEnvironmentProduction(),
	}
	if o, ok := os.LookupEnv("ANTHROPIC_BASE_URL"); ok {
		defaults = append(defaults, option.WithBaseURL(o))
	}

	statuses := []auth.CredentialSourceStatus{}

	if v, ok := os.LookupEnv("ANTHROPIC_API_KEY"); ok && v != "" {
		defaults = append(defaults, option.WithAPIKey(v))
		return defaults
	}
	statuses = append(statuses, auth.CredentialSourceStatus{
		Name:  "ANTHROPIC_API_KEY env var",
		State: auth.CredentialSourceNotSet,
	})

	if v, ok := os.LookupEnv("ANTHROPIC_AUTH_TOKEN"); ok && v != "" {
		defaults = append(defaults, option.WithAuthToken(v))
		return defaults
	}
	statuses = append(statuses, auth.CredentialSourceStatus{
		Name:  "ANTHROPIC_AUTH_TOKEN env var",
		State: auth.CredentialSourceNotSet,
	})

	// Step 3: explicit profile via ANTHROPIC_PROFILE. The user named a
	// specific profile, so a load failure is surfaced immediately — do
	// not fall through to env federation or the fallback profile.
	if profile, ok := os.LookupEnv("ANTHROPIC_PROFILE"); ok && profile != "" {
		cfg, err := config.LoadProfile(config.DefaultDir(), profile)
		if err != nil {
			return append(defaults, explicitProfileErrorOption(profile, err))
		}
		return append(defaults, option.WithConfig(cfg))
	}

	// Step 4: env-var federation. Beats the fallback profile so a
	// WIF-configured machine with a leftover default profile file still
	// uses WIF.
	envResult, envDetail, envState := auth.EnvCredentials()
	if envResult != nil {
		defaults = append(defaults, auth.WithAuthMiddleware(envResult.Provider))
		return defaults
	}
	envFederationStatus := auth.CredentialSourceStatus{
		Name:   "env federation (ANTHROPIC_FEDERATION_RULE_ID + ANTHROPIC_ORGANIZATION_ID + ANTHROPIC_IDENTITY_TOKEN_FILE)",
		State:  envState,
		Detail: envDetail,
	}

	// Step 5: fallback profile (active_config or literal "default"). A
	// missing profile here is a quiet miss — fall through to the no-
	// credentials aggregate.
	fallbackStatus, fallbackOpt := tryLoadFallbackProfile()
	if fallbackOpt != nil {
		return append(defaults, fallbackOpt)
	}
	if o, ok := os.LookupEnv("ANTHROPIC_CUSTOM_HEADERS"); ok {
		for _, line := range strings.Split(o, "\n") {
			colon := strings.Index(line, ":")
			if colon >= 0 {
				defaults = append(defaults, option.WithHeader(strings.TrimSpace(line[:colon]), strings.TrimSpace(line[colon+1:])))
			}
		}
	}

	statuses = append(statuses, envFederationStatus, fallbackStatus)
	defaults = append(defaults, noCredentialsSentinel(statuses))
	return defaults
}

// tryLoadFallbackProfile attempts the step-5 fallback profile lookup:
// active_config file, otherwise literal "default". A missing profile is
// reported as a silent-miss status (the caller will fall through to the
// no-credentials aggregate); any other load error is reported as a
// load-failure status so the user sees the specific OS error.
func tryLoadFallbackProfile() (auth.CredentialSourceStatus, option.RequestOption) {
	cfg, err := config.LoadConfig()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return auth.CredentialSourceStatus{
				Name:   "profile config file",
				State:  auth.CredentialSourceNotFound,
				Detail: "run `anthropic auth login` to create one",
			}, nil
		}
		return auth.CredentialSourceStatus{
			Name:   "profile config file",
			State:  auth.CredentialSourceLoadFailed,
			Detail: err.Error(),
		}, nil
	}
	return auth.CredentialSourceStatus{Name: "profile config file"}, option.WithConfig(cfg)
}

// explicitProfileErrorOption installs a middleware that fails the request
// with the underlying load error, unless a caller-supplied credential
// option preempts the profile. Used when ANTHROPIC_PROFILE names a profile
// whose config file cannot be loaded.
func explicitProfileErrorOption(profile string, loadErr error) option.RequestOption {
	profileErr := fmt.Errorf("ANTHROPIC_PROFILE=%q: %w", profile, loadErr)
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		cfg := r
		check := func(req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error) {
			if cfg.APIKey != "" || cfg.AuthToken != "" {
				return next(req)
			}
			if req.Header.Get("Authorization") != "" || req.Header.Get("X-Api-Key") != "" {
				return next(req)
			}
			return nil, profileErr
		}
		r.Middlewares = append(r.Middlewares, check)
		return nil
	})
}

func noCredentialsSentinel(statuses []auth.CredentialSourceStatus) option.RequestOption {
	preBuiltErr := &auth.NoCredentialsError{Sources: statuses}
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		cfg := r
		check := func(req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error) {
			if cfg.APIKey != "" || cfg.AuthToken != "" {
				return next(req)
			}
			if len(cfg.Middlewares) > 1 {
				return next(req)
			}
			if req.Header.Get("Authorization") != "" || req.Header.Get("X-Api-Key") != "" {
				return next(req)
			}
			return nil, preBuiltErr
		}
		r.Middlewares = append(r.Middlewares, check)
		return nil
	})
}

// NewClient generates a new client with the default option read from the
// environment (ANTHROPIC_API_KEY, ANTHROPIC_AUTH_TOKEN, ANTHROPIC_BASE_URL). The
// option passed in as arguments are applied after these default arguments, and all
// option will be passed down to the services and requests that this client makes.
//
// Pass [option.WithoutEnvironmentDefaults] to skip the environment-based
// credential autoload entirely (only the hardcoded production base-URL
// default is kept). Use this when the caller does its own credential
// resolution and wants the SDK to contribute nothing from the environment.
func NewClient(opts ...option.RequestOption) (r Client) {
	var defaults []option.RequestOption
	if option.HasWithoutEnvironmentDefaults(opts) {
		defaults = []option.RequestOption{option.WithEnvironmentProduction()}
	} else {
		defaults = DefaultClientOptions()
	}
	opts = append(defaults, opts...)

	r = Client{Options: opts}

	r.Completions = NewCompletionService(opts...)
	r.Messages = NewMessageService(opts...)
	r.Models = NewModelService(opts...)
	r.Beta = NewBetaService(opts...)

	return
}

// Execute makes a request with the given context, method, URL, request params,
// response, and request options. This is useful for hitting undocumented endpoints
// while retaining the base URL, auth, retries, and other options from the client.
//
// If a byte slice or an [io.Reader] is supplied to params, it will be used as-is
// for the request body.
//
// The params is by default serialized into the body using [encoding/json]. If your
// type implements a MarshalJSON function, it will be used instead to serialize the
// request. If a URLQuery method is implemented, the returned [url.Values] will be
// used as query strings to the url.
//
// If your params struct uses [param.Field], you must provide either [MarshalJSON],
// [URLQuery], and/or [MarshalForm] functions. It is undefined behavior to use a
// struct uses [param.Field] without specifying how it is serialized.
//
// Any "…Params" object defined in this library can be used as the request
// argument. Note that 'path' arguments will not be forwarded into the url.
//
// The response body will be deserialized into the res variable, depending on its
// type:
//
//   - A pointer to a [*http.Response] is populated by the raw response.
//   - A pointer to a byte array will be populated with the contents of the request
//     body.
//   - A pointer to any other type uses this library's default JSON decoding, which
//     respects UnmarshalJSON if it is defined on the type.
//   - A nil value will not read the response body.
//
// For even greater flexibility, see [option.WithResponseInto] and
// [option.WithResponseBodyInto].
func (r *Client) Execute(ctx context.Context, method string, path string, params any, res any, opts ...option.RequestOption) error {
	opts = slices.Concat(r.Options, opts)
	return requestconfig.ExecuteNewRequest(ctx, method, path, params, res, opts...)
}

// Get makes a GET request with the given URL, params, and optionally deserializes
// to a response. See [Execute] documentation on the params and response.
func (r *Client) Get(ctx context.Context, path string, params any, res any, opts ...option.RequestOption) error {
	return r.Execute(ctx, http.MethodGet, path, params, res, opts...)
}

// Post makes a POST request with the given URL, params, and optionally
// deserializes to a response. See [Execute] documentation on the params and
// response.
func (r *Client) Post(ctx context.Context, path string, params any, res any, opts ...option.RequestOption) error {
	return r.Execute(ctx, http.MethodPost, path, params, res, opts...)
}

// Put makes a PUT request with the given URL, params, and optionally deserializes
// to a response. See [Execute] documentation on the params and response.
func (r *Client) Put(ctx context.Context, path string, params any, res any, opts ...option.RequestOption) error {
	return r.Execute(ctx, http.MethodPut, path, params, res, opts...)
}

// Patch makes a PATCH request with the given URL, params, and optionally
// deserializes to a response. See [Execute] documentation on the params and
// response.
func (r *Client) Patch(ctx context.Context, path string, params any, res any, opts ...option.RequestOption) error {
	return r.Execute(ctx, http.MethodPatch, path, params, res, opts...)
}

// Delete makes a DELETE request with the given URL, params, and optionally
// deserializes to a response. See [Execute] documentation on the params and
// response.
func (r *Client) Delete(ctx context.Context, path string, params any, res any, opts ...option.RequestOption) error {
	return r.Execute(ctx, http.MethodDelete, path, params, res, opts...)
}

// CalculateNonStreamingTimeout calculates the appropriate timeout for a non-streaming request
// based on the maximum number of tokens and the model's non-streaming token limit
func CalculateNonStreamingTimeout(maxTokens int, model Model, opts []option.RequestOption) (time.Duration, error) {
	preCfg, err := requestconfig.PreRequestOptions(opts...)
	if err != nil {
		return 0, fmt.Errorf("error applying request options: %w", err)
	}
	// if the user has set a specific request timeout, use that
	if preCfg.RequestTimeout != 0 {
		return preCfg.RequestTimeout, nil
	}

	maximumTime := time.Hour // 1 hour
	defaultTime := 10 * time.Minute

	expectedTime := time.Duration(float64(maximumTime) * float64(maxTokens) / 128000.0)

	// If the model has a non-streaming token limit and max_tokens exceeds it,
	// or if the expected time exceeds default time, require streaming
	maxNonStreamingTokens, hasLimit := constant.ModelNonStreamingTokens[string(model)]
	if expectedTime > defaultTime || (hasLimit && maxTokens > maxNonStreamingTokens) {
		return 0, fmt.Errorf("streaming is required for operations that may take longer than 10 minutes")
	}

	return defaultTime, nil
}
