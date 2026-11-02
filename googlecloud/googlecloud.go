// Package googlecloud provides client configuration for Claude Models on Google
// Cloud (Anthropic managed) — the first-party Anthropic API served through the
// Google Cloud gateway.
//
// This is distinct from the older [github.com/anthropics/anthropic-sdk-go/v2/vertex]
// package, which targets the :rawPredict publisher-model API (publisher
// model IDs, messages only). This client speaks the first-party Anthropic API:
// requests pass through the gateway unchanged — standard /v1/* paths and
// standard model names.
package googlecloud

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

// ClientConfig holds the configuration for creating a Claude Models on Google
// Cloud (Anthropic managed) client.
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

	// WorkspaceID is the Anthropic workspace ID. Resolved by precedence:
	//  1. ClientConfig.WorkspaceID
	//  2. The ANTHROPIC_GOOGLE_CLOUD_WORKSPACE_ID environment variable
	//
	// Required when the base URL is derived, except by organization-level
	// routes (paths under /v1/organizations, and /v1/user_profiles and the
	// paths under it), which are sent to the project's organization endpoint.
	// A request that needs a workspace ID and has none fails before it is sent.
	WorkspaceID string

	// BaseURL overrides the default gateway base URL for every route. Resolved
	// by precedence: ClientConfig.BaseURL > ANTHROPIC_GOOGLE_CLOUD_BASE_URL env >
	// derived from project, location, and workspace ID.
	BaseURL string

	// TokenSource overrides Application Default Credentials for authentication.
	// When nil, credentials are resolved via Google ADC with the cloud-platform scope.
	TokenSource oauth2.TokenSource

	// SkipAuth skips authentication, for when a gateway or proxy handles
	// authentication upstream. No token is attached when SkipAuth is set.
	// Mutually exclusive with TokenSource — setting both is a
	// construction-time error.
	SkipAuth bool
}

// Client provides access to the Anthropic API via the Google Cloud gateway. It
// mirrors the surface of [anthropic.Client].
type Client struct {
	Options      []option.RequestOption
	Messages     anthropic.MessageService
	Models       anthropic.ModelService
	Files        anthropic.FileService
	Skills       anthropic.SkillService
	UserProfiles anthropic.UserProfileService
	Organization anthropic.OrganizationService
	Beta         anthropic.BetaService
}

// NewClient creates a new Claude Models on Google Cloud (Anthropic managed)
// client with the given configuration.
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
		Options:      opts,
		Messages:     anthropic.NewMessageService(opts...),
		Models:       anthropic.NewModelService(opts...),
		Files:        anthropic.NewFileService(opts...),
		Skills:       anthropic.NewSkillService(opts...),
		UserProfiles: anthropic.NewUserProfileService(opts...),
		Organization: anthropic.NewOrganizationService(opts...),
		Beta:         anthropic.NewBetaService(opts...),
	}, nil
}
