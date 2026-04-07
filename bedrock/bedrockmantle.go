package bedrock

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/awsauth"
	"github.com/anthropics/anthropic-sdk-go/option"
)

const mantleServiceName = "bedrock-mantle"

// MantleClientConfig holds the configuration for creating an Anthropic Bedrock Mantle client.
type MantleClientConfig struct {
	// APIKey is the Anthropic API key for x-api-key authentication.
	// Takes precedence over AWS credentials. When no AWS auth args are set, falls back
	// to the AWS_BEARER_TOKEN_BEDROCK environment variable (then ANTHROPIC_AWS_API_KEY)
	// before trying SigV4.
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

	// AWSRegion is the AWS region for the base URL and SigV4 signing.
	// Resolved by precedence: MantleClientConfig.AWSRegion > AWS_REGION env var.
	AWSRegion string

	// BaseURL overrides the default base URL.
	// Resolved by precedence: MantleClientConfig.BaseURL > ANTHROPIC_BEDROCK_MANTLE_BASE_URL env >
	// https://bedrock-mantle.{region}.api.aws/anthropic
	BaseURL string

	// SkipAuth skips Mantle-specific authentication (API key and SigV4).
	// This is useful when a gateway or proxy handles authentication on your behalf.
	// Note: when using [NewMantleClient], the base SDK may still send an X-Api-Key header
	// if the ANTHROPIC_API_KEY environment variable is set.
	SkipAuth bool
}

// MantleClient provides access to the Anthropic Bedrock Mantle API.
// Only the Messages API (/v1/messages) and its subpaths are supported.
type MantleClient struct {
	Options  []option.RequestOption
	Messages anthropic.MessageService
	Beta     MantleBetaService
}

// MantleBetaService exposes only the beta resources supported by Bedrock Mantle.
type MantleBetaService struct {
	Options  []option.RequestOption
	Messages anthropic.BetaMessageService
}

// NewMantleClient creates a new Bedrock Mantle client with the given configuration.
// Only the Messages API (/v1/messages) and its subpaths are supported on Bedrock Mantle.
//
// Any additional [option.RequestOption] values are applied after the client's
// internal options (base URL, auth, etc.), so they can be used to set custom
// headers, timeouts, middleware, and other request-level settings.
//
// Auth is resolved by precedence:
//  1. APIKey arg (x-api-key header)
//  2. AWSAccessKey + AWSSecretAccessKey args (SigV4)
//  3. AWSProfile arg (SigV4 via provider chain)
//  4. AWS_BEARER_TOKEN_BEDROCK env var, then ANTHROPIC_AWS_API_KEY (x-api-key header)
//  5. Default AWS credential chain (SigV4)
func NewMantleClient(ctx context.Context, cfg MantleClientConfig, opts ...option.RequestOption) (*MantleClient, error) {
	baseOpts, err := awsauth.CreateClientOptions(ctx, mantleToInternalConfig(cfg), mantleResolveParams())
	if err != nil {
		return nil, err
	}

	// We intentionally do not call anthropic.DefaultClientOptions() here.
	// The Mantle client resolves its own base URL, auth, and workspace ID — the
	// base SDK defaults (ANTHROPIC_API_KEY, ANTHROPIC_BASE_URL) do not apply.
	//
	// User-provided opts are appended last so they take highest precedence.
	opts = append(baseOpts, opts...)

	return &MantleClient{
		Options:  opts,
		Messages: anthropic.NewMessageService(opts...),
		Beta: MantleBetaService{
			Options:  opts,
			Messages: anthropic.NewBetaMessageService(opts...),
		},
	}, nil
}

func mantleResolveParams() awsauth.ResolveParams {
	return awsauth.ResolveParams{
		EnvAPIKey:         "AWS_BEARER_TOKEN_BEDROCK",
		EnvAPIKeyFallback: "ANTHROPIC_AWS_API_KEY",
		EnvBaseURL:        "ANTHROPIC_BEDROCK_MANTLE_BASE_URL",
		DeriveBaseURL:     func(region string) string { return fmt.Sprintf("https://bedrock-mantle.%s.api.aws/anthropic", region) },
		ServiceName:       mantleServiceName,
	}
}

func mantleToInternalConfig(cfg MantleClientConfig) awsauth.ClientConfig {
	return awsauth.ClientConfig{
		APIKey:             cfg.APIKey,
		AWSAccessKey:       cfg.AWSAccessKey,
		AWSSecretAccessKey: cfg.AWSSecretAccessKey,
		AWSSessionToken:    cfg.AWSSessionToken,
		AWSProfile:         cfg.AWSProfile,
		AWSRegion:          cfg.AWSRegion,
		BaseURL:            cfg.BaseURL,
		SkipAuth:           cfg.SkipAuth,
	}
}
