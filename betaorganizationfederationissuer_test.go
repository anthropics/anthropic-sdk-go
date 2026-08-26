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

func TestBetaOrganizationFederationIssuerNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.Federation.Issuers.New(context.TODO(), anthropic.BetaOrganizationFederationIssuerNewParams{
		IssuerURL: "x",
		Name:      "x",
		CheckJTI:  anthropic.Bool(true),
		JWKS: anthropic.BetaOrganizationFederationIssuerNewParamsJWKSUnion{
			OfDiscovery: &anthropic.BetaJWKSDiscoveryParam{
				CACertPEM:     anthropic.String("ca_cert_pem"),
				DiscoveryBase: anthropic.String("discovery_base"),
			},
		},
		MaxJWTLifetimeSeconds: anthropic.Int(1),
		Betas:                 []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaOrganizationFederationIssuerGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.Federation.Issuers.Get(
		context.TODO(),
		"federation_issuer_id",
		anthropic.BetaOrganizationFederationIssuerGetParams{
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

func TestBetaOrganizationFederationIssuerUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.Federation.Issuers.Update(
		context.TODO(),
		"federation_issuer_id",
		anthropic.BetaOrganizationFederationIssuerUpdateParams{
			CheckJTI:  anthropic.Bool(true),
			IssuerURL: anthropic.String("x"),
			JWKS: anthropic.BetaOrganizationFederationIssuerUpdateParamsJWKSUnion{
				OfDiscovery: &anthropic.BetaJWKSDiscoveryParam{
					CACertPEM:     anthropic.String("ca_cert_pem"),
					DiscoveryBase: anthropic.String("discovery_base"),
				},
			},
			JWKSPollingDisabled:   anthropic.Bool(true),
			MaxJWTLifetimeSeconds: anthropic.Int(1),
			Name:                  anthropic.String("x"),
			Betas:                 []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
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

func TestBetaOrganizationFederationIssuerListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.Federation.Issuers.List(context.TODO(), anthropic.BetaOrganizationFederationIssuerListParams{
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

func TestBetaOrganizationFederationIssuerArchiveWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.Federation.Issuers.Archive(
		context.TODO(),
		"federation_issuer_id",
		anthropic.BetaOrganizationFederationIssuerArchiveParams{
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
