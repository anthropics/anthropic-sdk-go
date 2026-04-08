// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestBetaEnvironmentNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Environments.New(context.TODO(), anthropic.BetaEnvironmentNewParams{
		Name: "python-data-analysis",
		Config: anthropic.BetaCloudConfigParams{
			Networking: anthropic.BetaCloudConfigParamsNetworkingUnion{
				OfLimited: &anthropic.BetaLimitedNetworkParams{
					AllowMCPServers:      anthropic.Bool(true),
					AllowPackageManagers: anthropic.Bool(true),
					AllowedHosts:         []string{"api.example.com"},
				},
			},
			Packages: anthropic.BetaPackagesParams{
				Apt:   []string{"string"},
				Cargo: []string{"string"},
				Gem:   []string{"string"},
				Go:    []string{"string"},
				Npm:   []string{"string"},
				Pip:   []string{"pandas", "numpy"},
				Type:  anthropic.BetaPackagesParamsTypePackages,
			},
		},
		Description: anthropic.String("Python environment with data-analysis packages."),
		Metadata: map[string]string{
			"foo": "string",
		},
		Betas: []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaEnvironmentGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Environments.Get(
		context.TODO(),
		"env_011CZkZ9X2dpNyB7HsEFoRfW",
		anthropic.BetaEnvironmentGetParams{
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

func TestBetaEnvironmentUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Environments.Update(
		context.TODO(),
		"env_011CZkZ9X2dpNyB7HsEFoRfW",
		anthropic.BetaEnvironmentUpdateParams{
			Config: anthropic.BetaCloudConfigParams{
				Networking: anthropic.BetaCloudConfigParamsNetworkingUnion{
					OfLimited: &anthropic.BetaLimitedNetworkParams{
						AllowMCPServers:      anthropic.Bool(true),
						AllowPackageManagers: anthropic.Bool(true),
						AllowedHosts:         []string{"api.example.com"},
					},
				},
				Packages: anthropic.BetaPackagesParams{
					Apt:   []string{"string"},
					Cargo: []string{"string"},
					Gem:   []string{"string"},
					Go:    []string{"string"},
					Npm:   []string{"string"},
					Pip:   []string{"pandas", "numpy"},
					Type:  anthropic.BetaPackagesParamsTypePackages,
				},
			},
			Description: anthropic.String("Python environment with data-analysis packages."),
			Metadata: map[string]string{
				"foo": "string",
			},
			Name:  anthropic.String("x"),
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

func TestBetaEnvironmentListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Environments.List(context.TODO(), anthropic.BetaEnvironmentListParams{
		IncludeArchived: anthropic.Bool(true),
		Limit:           anthropic.Int(1),
		Page:            anthropic.String("page"),
		Betas:           []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaEnvironmentDeleteWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Environments.Delete(
		context.TODO(),
		"env_011CZkZ9X2dpNyB7HsEFoRfW",
		anthropic.BetaEnvironmentDeleteParams{
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

func TestBetaEnvironmentArchiveWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Environments.Archive(
		context.TODO(),
		"env_011CZkZ9X2dpNyB7HsEFoRfW",
		anthropic.BetaEnvironmentArchiveParams{
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
