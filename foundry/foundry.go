// Package foundry configures the SDK for Claude on Microsoft Foundry, which
// serves the Anthropic API under https://{resource}.services.ai.azure.com/anthropic/.
//
// Configuration falls back to the ANTHROPIC_FOUNDRY_API_KEY,
// ANTHROPIC_FOUNDRY_RESOURCE and ANTHROPIC_FOUNDRY_BASE_URL environment
// variables. For Microsoft Entra ID, adapt an azidentity credential:
//
//	AzureADTokenProvider: func(ctx context.Context) (string, error) {
//		tok, err := cred.GetToken(ctx, policy.TokenRequestOptions{Scopes: []string{foundry.EntraIDScope}})
//		return tok.Token, err
//	},
//
// On Foundry, the model in a request is the name of your deployment, which
// defaults to the model ID. Foundry does not serve the Models API, which
// [Client] leaves out, or Message Batches. Files and Skills need a deployment
// hosted on Anthropic; one hosted on Azure rejects those requests with a 400.
package foundry

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

const (
	envAPIKey   = "ANTHROPIC_FOUNDRY_API_KEY"
	envResource = "ANTHROPIC_FOUNDRY_RESOURCE"
	envBaseURL  = "ANTHROPIC_FOUNDRY_BASE_URL"
)

// EntraIDScope is the scope to request for [ClientConfig.AzureADTokenProvider] tokens.
const EntraIDScope = "https://ai.azure.com/.default"

const baseURLTemplate = "https://%s.services.ai.azure.com/anthropic/"

// ClientConfig configures [NewClient]. Empty fields fall back to the
// ANTHROPIC_FOUNDRY_* environment variables.
type ClientConfig struct {
	// APIKey is sent as the x-api-key header. Mutually exclusive with
	// AzureADTokenProvider.
	APIKey string

	// AzureADTokenProvider returns an Entra ID access token for [EntraIDScope],
	// sent as a bearer token. It is called on every attempt, including retries,
	// so it should cache tokens; azidentity credentials do.
	AzureADTokenProvider func(ctx context.Context) (string, error)

	// Resource is the name in https://{resource}.services.ai.azure.com/anthropic/.
	// Mutually exclusive with BaseURL.
	Resource string

	// BaseURL replaces the resource-derived URL, for example for a private link
	// or proxy. It must include the /anthropic path.
	BaseURL string
}

// Client is [anthropic.Client] without the services Foundry does not serve.
type Client struct {
	// Options includes [option.WithoutEnvironmentDefaults], so passing it to
	// anthropic.NewClient does not pick up first-party credentials.
	Options  []option.RequestOption
	Messages anthropic.MessageService
	Files    anthropic.FileService
	Skills   anthropic.SkillService
	Beta     anthropic.BetaService
}

// NewClient returns a client for cfg. It never reads the first-party
// ANTHROPIC_API_KEY, ANTHROPIC_BASE_URL or config profiles, so a first-party
// credential cannot reach Foundry. opts apply after the endpoint and credential.
func NewClient(cfg ClientConfig, opts ...option.RequestOption) (*Client, error) {
	options, err := createClientOptions(cfg, opts)
	if err != nil {
		return nil, err
	}
	return &Client{
		Options:  options,
		Messages: anthropic.NewMessageService(options...),
		Files:    anthropic.NewFileService(options...),
		Skills:   anthropic.NewSkillService(options...),
		Beta:     anthropic.NewBetaService(options...),
	}, nil
}

// resolveConfig applies the environment fallbacks and requires exactly one
// credential and one endpoint.
func resolveConfig(cfg ClientConfig) (ClientConfig, error) {
	if cfg.APIKey != "" && cfg.AzureADTokenProvider != nil {
		return cfg, errors.New("foundry: ClientConfig.APIKey and ClientConfig.AzureADTokenProvider are mutually exclusive; set one or the other")
	}
	if cfg.AzureADTokenProvider == nil {
		if cfg.APIKey == "" {
			cfg.APIKey = os.Getenv(envAPIKey)
		}
		if cfg.APIKey == "" {
			return cfg, fmt.Errorf("foundry: no credentials found; set ClientConfig.APIKey, ClientConfig.AzureADTokenProvider, or the %s environment variable", envAPIKey)
		}
	}

	if cfg.Resource == "" {
		cfg.Resource = os.Getenv(envResource)
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = os.Getenv(envBaseURL)
	}
	switch {
	case cfg.BaseURL != "" && cfg.Resource != "":
		return cfg, fmt.Errorf("foundry: ClientConfig.BaseURL (or %s) and ClientConfig.Resource (or %s) are mutually exclusive; set one or the other", envBaseURL, envResource)
	case cfg.BaseURL == "" && cfg.Resource == "":
		return cfg, fmt.Errorf("foundry: no endpoint found; set ClientConfig.Resource, ClientConfig.BaseURL, or the %s or %s environment variable", envResource, envBaseURL)
	case cfg.BaseURL == "":
		cfg.BaseURL = fmt.Sprintf(baseURLTemplate, cfg.Resource)
	}
	return cfg, nil
}

// createClientOptions registers the bearer middleware last, so the credential
// stays closest to the wire and caller middleware can set its own Authorization.
func createClientOptions(cfg ClientConfig, userOpts []option.RequestOption) ([]option.RequestOption, error) {
	resolved, err := resolveConfig(cfg)
	if err != nil {
		return nil, err
	}

	// Stops anthropic.NewClient(client.Options...) from adding first-party credentials.
	opts := make([]option.RequestOption, 0, len(userOpts)+3)
	opts = append(opts, option.WithoutEnvironmentDefaults(), option.WithBaseURL(resolved.BaseURL))
	if resolved.AzureADTokenProvider == nil {
		opts = append(opts, option.WithAPIKey(resolved.APIKey))
	}
	opts = append(opts, userOpts...)
	if resolved.AzureADTokenProvider != nil {
		opts = append(opts, option.WithMiddleware(bearerMiddleware(resolved.AzureADTokenProvider)))
	}
	return opts, nil
}

// bearerMiddleware calls the provider on every attempt, so refresh stays with
// the provider. A provider error is returned as a transport error and retried.
func bearerMiddleware(provider func(context.Context) (string, error)) option.Middleware {
	return func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		if r.Header.Get("Authorization") == "" {
			token, err := provider(r.Context())
			if err != nil {
				return nil, fmt.Errorf("foundry: AzureADTokenProvider failed: %w", err)
			}
			if token == "" {
				return nil, errors.New("foundry: AzureADTokenProvider returned an empty token")
			}
			r.Header.Set("Authorization", "Bearer "+token)
		}
		return next(r)
	}
}
