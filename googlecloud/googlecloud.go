// Package googlecloud provides client configuration for Claude Platform on
// Google Cloud — the first-party Anthropic API served through the Google Cloud
// gateway.
//
// This is distinct from the older [github.com/anthropics/anthropic-sdk-go/vertex]
// package, which targets the :rawPredict publisher-model API (publisher
// model IDs, messages only). This client speaks the full first-party Anthropic
// API: requests pass through the gateway unchanged — standard /v1/* paths,
// standard model names, the complete API surface.
//
// The deprecated text Completions API is intentionally not exposed.
package googlecloud

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// ClientConfig holds the configuration for creating a Claude Platform on Google Cloud client.
type ClientConfig struct {
	// Project is the GCP consumer project ID. Resolved by precedence:
	//  1. ClientConfig.Project
	//  2. The ANTHROPIC_GOOGLE_CLOUD_PROJECT environment variable
	//  3. The GOOGLE_CLOUD_PROJECT environment variable
	//  4. The project reported by Application Default Credentials
	//
	// The ADC fallback (4) applies only when TokenSource is nil and SkipAuth is
	// false; with an explicit TokenSource or SkipAuth, set Project (or BaseURL)
	// explicitly.
	Project string

	// Location is the GCP location. Optional — defaults to "global", which is
	// the region the gateway should normally be addressed through. Resolved by
	// precedence:
	//  1. ClientConfig.Location
	//  2. The ANTHROPIC_GOOGLE_CLOUD_LOCATION environment variable
	//  3. "global"
	Location string

	// WorkspaceID is the Anthropic workspace ID. Required: resolved by
	// precedence
	//  1. ClientConfig.WorkspaceID
	//  2. The ANTHROPIC_GOOGLE_CLOUD_WORKSPACE_ID environment variable
	//
	// [NewClient] returns an error when neither is set, unless SkipAuth is true
	// and BaseURL is set explicitly.
	WorkspaceID string

	// BaseURL overrides the default gateway base URL. Resolved by precedence:
	// ClientConfig.BaseURL > ANTHROPIC_GOOGLE_CLOUD_BASE_URL env > derived from
	// project, location, and workspace ID.
	BaseURL string

	// TokenSource overrides Application Default Credentials for authentication.
	// When nil, credentials are resolved via Google ADC with the cloud-platform scope.
	TokenSource oauth2.TokenSource

	// SkipAuth skips authentication, for when a gateway or proxy handles
	// authentication upstream. No token is attached when SkipAuth is set.
	// A workspace ID is still needed to derive the base URL —
	// set BaseURL explicitly to construct without one. Mutually exclusive with
	// TokenSource — setting both is a construction-time error.
	SkipAuth bool
}

// Client provides access to the Anthropic API via the Google Cloud gateway. It
// mirrors the surface of [anthropic.Client]; the gateway proxies the entire
// first-party API. The deprecated Completions service is intentionally omitted.
type Client struct {
	Options  []option.RequestOption
	Messages anthropic.MessageService
	Models   anthropic.ModelService
	Beta     anthropic.BetaService
}

// NewClient creates a new Claude Platform on Google Cloud client with the given
// configuration.
//
// Authentication uses Google credentials by precedence:
//  1. TokenSource arg
//  2. Application Default Credentials (cloud-platform scope)
//
// The ctx is used for credential discovery during construction (e.g. probing
// the metadata server) and is not retained afterward; it does not need to
// outlive the client. Use a per-request context for individual API calls.
func NewClient(ctx context.Context, cfg ClientConfig) (*Client, error) {
	opts, err := createClientOptions(ctx, cfg)
	if err != nil {
		return nil, err
	}

	// We intentionally do not call anthropic.DefaultClientOptions() here.
	// This client resolves its own base URL, auth, and workspace ID — the base
	// SDK defaults (ANTHROPIC_API_KEY, ANTHROPIC_BASE_URL) do not apply.

	return &Client{
		Options:  opts,
		Messages: anthropic.NewMessageService(opts...),
		Models:   anthropic.NewModelService(opts...),
		Beta:     anthropic.NewBetaService(opts...),
	}, nil
}
