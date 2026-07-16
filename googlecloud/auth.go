// Config resolution and bearer-token auth for the Anthropic Google Cloud client.

package googlecloud

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/anthropics/anthropic-sdk-go/option"
)

// Environment variables consulted as fallbacks for the corresponding
// [ClientConfig] fields.
const (
	envProject     = "ANTHROPIC_GOOGLE_CLOUD_PROJECT"
	envLocation    = "ANTHROPIC_GOOGLE_CLOUD_LOCATION"
	envWorkspaceID = "ANTHROPIC_GOOGLE_CLOUD_WORKSPACE_ID"
	envBaseURL     = "ANTHROPIC_GOOGLE_CLOUD_BASE_URL"
)

// cloudPlatformScope is the OAuth2 scope required to call the gateway.
const cloudPlatformScope = "https://www.googleapis.com/auth/cloud-platform"

// baseURLTemplate is the gateway base URL; override via BaseURL if needed.
// %[1]s is the project, %[2]s the location, %[3]s the workspace ID.
const baseURLTemplate = "https://claude.googleapis.com/v1alpha/projects/%[1]s/locations/%[2]s/workspaces/%[3]s/invoke"

// defaultLocation is the location used when none is configured; the gateway
// should always be addressed via the global region.
const defaultLocation = "global"

// resolveConfig applies environment-variable fallbacks and derives the base URL
// when a project is known. It does not touch the network; ADC-based project
// backfill happens in createClientOptions. Location defaults to "global" when
// not configured. The workspace ID is embedded in the derived base URL, so it
// is required unless SkipAuth is set together with an explicit BaseURL.
func resolveConfig(cfg ClientConfig) (ClientConfig, error) {
	cfg.Project = firstNonEmpty(cfg.Project, os.Getenv(envProject), os.Getenv("GOOGLE_CLOUD_PROJECT"))
	cfg.Location = firstNonEmpty(cfg.Location, os.Getenv(envLocation), defaultLocation)
	cfg.BaseURL = firstNonEmpty(cfg.BaseURL, os.Getenv(envBaseURL))
	cfg.WorkspaceID = firstNonEmpty(cfg.WorkspaceID, os.Getenv(envWorkspaceID))

	if cfg.SkipAuth && cfg.TokenSource != nil {
		return cfg, fmt.Errorf("googlecloud: SkipAuth and TokenSource are mutually exclusive; set one or the other")
	}

	// The workspace ID is required unless SkipAuth is set together with an
	// explicit BaseURL (no URL to derive).
	if cfg.WorkspaceID == "" && !(cfg.SkipAuth && cfg.BaseURL != "") {
		return cfg, fmt.Errorf("googlecloud: no workspace ID found; set WorkspaceID or the %s environment variable", envWorkspaceID)
	}

	if cfg.BaseURL == "" && cfg.Project != "" {
		cfg.BaseURL = deriveBaseURL(cfg.Project, cfg.Location, cfg.WorkspaceID)
	}

	return cfg, nil
}

// createClientOptions resolves configuration and returns the request options that
// configure an Anthropic client for the Google Cloud gateway: the base URL
// (which embeds the workspace ID) and a Google bearer-token middleware.
func createClientOptions(ctx context.Context, cfg ClientConfig) ([]option.RequestOption, error) {
	resolved, err := resolveConfig(cfg)
	if err != nil {
		return nil, err
	}

	// Resolve credentials first: ADC also yields a project we can use to backfill
	// both an unset Project and the derived base URL.
	var ts oauth2.TokenSource
	if !resolved.SkipAuth {
		ts = resolved.TokenSource
		if ts == nil {
			creds, err := google.FindDefaultCredentials(ctx, cloudPlatformScope)
			if err != nil {
				return nil, fmt.Errorf("googlecloud: failed to load Google application default credentials: %w", err)
			}
			ts = creds.TokenSource
			if resolved.Project == "" {
				resolved.Project = creds.ProjectID
			}
		}
		// Guarantee the per-request caching the middleware relies on, regardless of
		// what the caller passed (ADC's source already caches; a naive caller source
		// would otherwise be hit on every request and retry).
		ts = oauth2.ReuseTokenSource(nil, ts)
	}

	if resolved.BaseURL == "" {
		if resolved.Project == "" {
			return nil, fmt.Errorf("googlecloud: no project found; set Project, set the %s environment variable, or configure application default credentials with a project", envProject)
		}
		resolved.BaseURL = deriveBaseURL(resolved.Project, resolved.Location, resolved.WorkspaceID)
	}

	opts := []option.RequestOption{option.WithBaseURL(resolved.BaseURL)}

	if resolved.SkipAuth {
		return opts, nil
	}

	opts = append(opts, option.WithMiddleware(bearerMiddleware(ts)))
	return opts, nil
}

// bearerMiddleware attaches a Google OAuth2 access token to each request as
// Authorization: Bearer, only if the header is not already set. The token source
// caches and refreshes internally, so this is cheap per request. Unlike SigV4, a
// bearer token does not depend on the request body, so the body is never read or
// buffered here — streaming and multipart uploads pass through untouched.
func bearerMiddleware(ts oauth2.TokenSource) option.Middleware {
	return func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		if r.Header.Get("Authorization") == "" {
			tok, err := ts.Token()
			if err != nil {
				return nil, fmt.Errorf("googlecloud: failed to fetch Google access token: %w", err)
			}
			tok.SetAuthHeader(r)
		}
		return next(r)
	}
}

func deriveBaseURL(project, location, workspaceID string) string {
	return fmt.Sprintf(baseURLTemplate, project, location, workspaceID)
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
