// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestBetaVaultCredentialNewWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := anthropic.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("my-anthropic-api-key"),
	)
	_, err := client.Beta.Vaults.Credentials.New(
		context.TODO(),
		"vlt_011CZkZDLs7fYzm1hXNPeRjv",
		anthropic.BetaVaultCredentialNewParams{
			Auth: anthropic.BetaVaultCredentialNewParamsAuthUnion{
				OfStaticBearer: &anthropic.BetaManagedAgentsStaticBearerCreateParams{
					Token:        "bearer_exampletoken",
					MCPServerURL: "https://example-server.modelcontextprotocol.io/sse",
					Type:         anthropic.BetaManagedAgentsStaticBearerCreateParamsTypeStaticBearer,
				},
			},
			DisplayName: anthropic.String("Example credential"),
			Metadata: map[string]string{
				"environment": "production",
			},
			Betas: []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
		},
	)
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaVaultCredentialGetWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := anthropic.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("my-anthropic-api-key"),
	)
	_, err := client.Beta.Vaults.Credentials.Get(
		context.TODO(),
		"vcrd_011CZkZEMt8gZan2iYOQfSkw",
		anthropic.BetaVaultCredentialGetParams{
			VaultID: "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			Betas:   []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
		},
	)
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaVaultCredentialUpdateWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := anthropic.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("my-anthropic-api-key"),
	)
	_, err := client.Beta.Vaults.Credentials.Update(
		context.TODO(),
		"vcrd_011CZkZEMt8gZan2iYOQfSkw",
		anthropic.BetaVaultCredentialUpdateParams{
			VaultID: "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			Auth: anthropic.BetaVaultCredentialUpdateParamsAuthUnion{
				OfMCPOAuth: &anthropic.BetaManagedAgentsMCPOAuthUpdateParams{
					Type:        anthropic.BetaManagedAgentsMCPOAuthUpdateParamsTypeMCPOAuth,
					AccessToken: anthropic.String("x"),
					ExpiresAt:   anthropic.Time(time.Now()),
					Refresh: anthropic.BetaManagedAgentsMCPOAuthRefreshUpdateParams{
						RefreshToken: anthropic.String("x"),
						Scope:        anthropic.String("scope"),
						TokenEndpointAuth: anthropic.BetaManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion{
							OfClientSecretBasic: &anthropic.BetaManagedAgentsTokenEndpointAuthBasicUpdateParam{
								Type:         anthropic.BetaManagedAgentsTokenEndpointAuthBasicUpdateParamTypeClientSecretBasic,
								ClientSecret: anthropic.String("x"),
							},
						},
					},
				},
			},
			DisplayName: anthropic.String("Example credential"),
			Metadata: map[string]string{
				"environment": "production",
			},
			Betas: []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
		},
	)
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaVaultCredentialListWithOptionalParams(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := anthropic.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("my-anthropic-api-key"),
	)
	_, err := client.Beta.Vaults.Credentials.List(
		context.TODO(),
		"vlt_011CZkZDLs7fYzm1hXNPeRjv",
		anthropic.BetaVaultCredentialListParams{
			IncludeArchived: anthropic.Bool(true),
			Limit:           anthropic.Int(0),
			Page:            anthropic.String("page"),
			Betas:           []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
		},
	)
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaVaultCredentialDeleteWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := anthropic.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("my-anthropic-api-key"),
	)
	_, err := client.Beta.Vaults.Credentials.Delete(
		context.TODO(),
		"vcrd_011CZkZEMt8gZan2iYOQfSkw",
		anthropic.BetaVaultCredentialDeleteParams{
			VaultID: "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			Betas:   []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
		},
	)
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaVaultCredentialArchiveWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := anthropic.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("my-anthropic-api-key"),
	)
	_, err := client.Beta.Vaults.Credentials.Archive(
		context.TODO(),
		"vcrd_011CZkZEMt8gZan2iYOQfSkw",
		anthropic.BetaVaultCredentialArchiveParams{
			VaultID: "vlt_011CZkZDLs7fYzm1hXNPeRjv",
			Betas:   []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
		},
	)
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
