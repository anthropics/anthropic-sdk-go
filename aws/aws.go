// Package aws provides client configuration for the Anthropic AWS gateway.
//
// This package name may shadow the AWS SDK's "aws" package. If you use both in the
// same file, alias one of them:
//
//	import (
//	    anthropicaws "github.com/anthropics/anthropic-sdk-go/v2/aws"
//	    awssdk "github.com/aws/aws-sdk-go-v2/aws"
//	)
package aws

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/awsauth"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

const defaultServiceName = "aws-external-anthropic"

// ClientConfig holds the configuration for creating an Anthropic AWS gateway client.
type ClientConfig struct {
	// APIKey is the Anthropic API key, sent as an Authorization: Bearer header.
	// Takes precedence over AWS credentials. When no AWS auth args are set, falls back
	// to the ANTHROPIC_AWS_API_KEY environment variable before trying SigV4.
	APIKey string

	// AWSAccessKey is the AWS access key ID for SigV4 authentication.
	// Must be paired with AWSSecretAccessKey. When unset, credentials are resolved
	// via the default AWS credential chain (env vars, shared credentials file, IAM roles, etc.).
	AWSAccessKey string

	// AWSSecretAccessKey is the AWS secret access key for SigV4 authentication.
	// When unset, credentials are resolved via the default AWS credential chain
	// (env vars, shared credentials file, IAM roles, etc.).
	AWSSecretAccessKey string

	// AWSSessionToken is the optional AWS session token for temporary credentials.
	// When unset, resolved via the default AWS credential chain if applicable.
	AWSSessionToken string

	// AWSProfile is the AWS named profile for credential resolution via the provider chain.
	AWSProfile string

	// AWSRegion is the AWS region for the gateway URL and SigV4 signing.
	// Resolved by precedence: ClientConfig.AWSRegion > AWS_REGION env var.
	AWSRegion string

	// WorkspaceID is sent as the anthropic-workspace-id header on every request.
	// Resolved by precedence: ClientConfig.WorkspaceID > ANTHROPIC_AWS_WORKSPACE_ID env var. Required.
	WorkspaceID string

	// BaseURL overrides the default gateway base URL.
	// Resolved by precedence: ClientConfig.BaseURL > ANTHROPIC_AWS_BASE_URL env > https://aws-external-anthropic.{region}.api.aws
	BaseURL string

	// SkipAuth skips all authentication (API key and SigV4) and the workspace ID requirement.
	// This is useful when a gateway or proxy handles authentication on your behalf.
	SkipAuth bool
}

// Client provides access to the Anthropic API via the AWS gateway.
//
// The gateway serves user profiles but not their action URLs, so
// UserProfiles.ActionURLs and Beta.UserProfiles.ActionURLs return an error
// instead of sending the request.
type Client struct {
	Options      []option.RequestOption
	Messages     anthropic.MessageService
	Models       anthropic.ModelService
	Files        anthropic.FileService
	Skills       anthropic.SkillService
	UserProfiles anthropic.UserProfileService
	Beta         anthropic.BetaService
}

// NewClient creates a new AWS gateway client with the given configuration.
//
// Any additional [option.RequestOption] values are applied after the client's
// internal options (base URL, auth, etc.), so they can be used to set custom
// headers, timeouts, middleware, and other request-level settings. When SigV4
// authentication is in use, the signing middleware runs after any middleware
// registered through these options, so the signature covers their request
// mutations. Per-request middleware passed at a method call site instead runs
// after signing and must not mutate the request.
//
// Auth is resolved by precedence:
//  1. APIKey arg (Authorization: Bearer header)
//  2. AWSAccessKey + AWSSecretAccessKey args (SigV4)
//  3. AWSProfile arg (SigV4 via provider chain)
//  4. ANTHROPIC_AWS_API_KEY env var (Authorization: Bearer header)
//  5. Default AWS credential chain (SigV4)
func NewClient(ctx context.Context, cfg ClientConfig, opts ...option.RequestOption) (*Client, error) {
	opts = slices.Concat(opts, []option.RequestOption{option.WithMiddleware(rejectUserProfileActionURLs)})

	opts, err := awsauth.CreateClientOptions(ctx, toInternalConfig(cfg), awsResolveParams(), opts...)
	if err != nil {
		return nil, err
	}
	opts = append(opts, requestconfig.WithOpenTelemetryPlatform("aws"))

	// We intentionally do not call anthropic.DefaultClientOptions() here.
	// The AWS client resolves its own base URL, auth, and workspace ID — the
	// base SDK defaults (ANTHROPIC_API_KEY, ANTHROPIC_BASE_URL) do not apply.

	return &Client{
		Options:      opts,
		Messages:     anthropic.NewMessageService(opts...),
		Models:       anthropic.NewModelService(opts...),
		Files:        anthropic.NewFileService(opts...),
		Skills:       anthropic.NewSkillService(opts...),
		UserProfiles: anthropic.NewUserProfileService(opts...),
		Beta:         anthropic.NewBetaService(opts...),
	}, nil
}

func rejectUserProfileActionURLs(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
	if strings.Contains(r.URL.Path, "/v1/user_profiles/") && strings.HasSuffix(r.URL.Path, "/action_urls") {
		return nil, errors.New("aws: the AWS gateway does not serve user profile action URLs; create them with a first-party client")
	}
	return next(r)
}

func awsResolveParams() awsauth.ResolveParams {
	return awsauth.ResolveParams{
		EnvAPIKey:      "ANTHROPIC_AWS_API_KEY",
		EnvWorkspaceID: "ANTHROPIC_AWS_WORKSPACE_ID",
		EnvBaseURL:     "ANTHROPIC_AWS_BASE_URL",
		DeriveBaseURL:  func(region string) string { return fmt.Sprintf("https://aws-external-anthropic.%s.api.aws", region) },
		ServiceName:    defaultServiceName,
	}
}

func toInternalConfig(cfg ClientConfig) awsauth.ClientConfig {
	return awsauth.ClientConfig{
		APIKey:             cfg.APIKey,
		AWSAccessKey:       cfg.AWSAccessKey,
		AWSSecretAccessKey: cfg.AWSSecretAccessKey,
		AWSSessionToken:    cfg.AWSSessionToken,
		AWSProfile:         cfg.AWSProfile,
		AWSRegion:          cfg.AWSRegion,
		WorkspaceID:        cfg.WorkspaceID,
		BaseURL:            cfg.BaseURL,
		SkipAuth:           cfg.SkipAuth,
	}
}
