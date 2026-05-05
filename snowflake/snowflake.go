package snowflake

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const DefaultVersion = "snowflake-2025-01-01"

// WithAccount returns a request option that configures the SDK to use Snowflake Cortex
// with the given account identifier and token. The account identifier is used to construct
// the base URL (https://<account>.snowflakecomputing.com/) and the token is used for
// Bearer authentication.
//
// If the token is empty, the SNOWFLAKE_AUTH_TOKEN environment variable is used.
//
// The account parameter should be in the format described in the Snowflake documentation
// for account identifiers (e.g. "myorg-myaccount").
//
// For more details, see https://docs.snowflake.com/en/user-guide/snowflake-cortex/cortex-rest-api
func WithAccount(account string, token string) option.RequestOption {
	if account == "" {
		panic("snowflake: account must be provided")
	}

	if token == "" {
		token = os.Getenv("SNOWFLAKE_AUTH_TOKEN")
	}

	if token == "" {
		panic("snowflake: token must be provided or SNOWFLAKE_AUTH_TOKEN environment variable must be set")
	}

	// Normalize the account: strip any trailing .snowflakecomputing.com if present
	account = strings.TrimSuffix(account, ".snowflakecomputing.com")

	baseURL := fmt.Sprintf("https://%s.snowflakecomputing.com/", account)
	middleware := cortexMiddleware(token)

	return requestconfig.RequestOptionFunc(func(rc *requestconfig.RequestConfig) error {
		return rc.Apply(
			option.WithBaseURL(baseURL),
			option.WithMiddleware(middleware),
		)
	})
}

// WithBaseURL returns a request option that configures the SDK to use Snowflake Cortex
// with a custom base URL and token. This is useful when you need to specify a full URL
// directly (e.g., for proxies or custom deployments).
//
// If the token is empty, the SNOWFLAKE_AUTH_TOKEN environment variable is used.
func WithBaseURL(baseURL string, token string) option.RequestOption {
	if baseURL == "" {
		panic("snowflake: baseURL must be provided")
	}

	if token == "" {
		token = os.Getenv("SNOWFLAKE_AUTH_TOKEN")
	}

	if token == "" {
		panic("snowflake: token must be provided or SNOWFLAKE_AUTH_TOKEN environment variable must be set")
	}

	middleware := cortexMiddleware(token)

	return requestconfig.RequestOptionFunc(func(rc *requestconfig.RequestConfig) error {
		return rc.Apply(
			option.WithBaseURL(baseURL),
			option.WithMiddleware(middleware),
		)
	})
}

// cortexMiddleware returns middleware that transforms Anthropic API requests into
// Snowflake Cortex REST API format.
//
// The middleware:
//   - Rewrites /v1/messages to /api/v2/cortex/inference:complete
//   - Rewrites /v1/messages/count_tokens to /api/v2/cortex/inference:complete (count tokens)
//   - Adds the anthropic_version field to the request body if not present
//   - Sets the Authorization header with the Bearer token
//   - Sets required Accept and Content-Type headers
func cortexMiddleware(token string) option.Middleware {
	return func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		if r.Body != nil {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return nil, err
			}
			r.Body.Close()

			if !gjson.GetBytes(body, "anthropic_version").Exists() {
				body, _ = sjson.SetBytes(body, "anthropic_version", DefaultVersion)
			}

			if r.URL.Path == "/v1/messages" && r.Method == http.MethodPost {
				r.URL.Path = "/api/v2/cortex/inference:complete"
			}

			if r.URL.Path == "/v1/messages/count_tokens" && r.Method == http.MethodPost {
				r.URL.Path = "/api/v2/cortex/inference:complete"
			}

			reader := bytes.NewReader(body)
			r.Body = io.NopCloser(reader)
			r.GetBody = func() (io.ReadCloser, error) {
				_, err := reader.Seek(0, 0)
				return io.NopCloser(reader), err
			}
			r.ContentLength = int64(len(body))
		}

		// Set authentication header
		r.Header.Set("Authorization", "Bearer "+token)

		// Set required headers for Cortex REST API
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Accept", "application/json, text/event-stream")

		return next(r)
	}
}
