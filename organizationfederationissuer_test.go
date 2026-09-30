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

func TestOrganizationFederationIssuerNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Organization.Federation.Issuers.New(context.TODO(), anthropic.OrganizationFederationIssuerNewParams{
		IssuerURL: "x",
		Name:      "x",
		CheckJTI:  anthropic.Bool(true),
		JWKS: anthropic.OrganizationFederationIssuerNewParamsJWKSUnion{
			OfDiscovery: &anthropic.JWKSDiscoveryParam{
				CACertPEM:     anthropic.String("ca_cert_pem"),
				DiscoveryBase: anthropic.String("discovery_base"),
			},
		},
		MaxJWTLifetimeSeconds: anthropic.Int(1),
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOrganizationFederationIssuerGet(t *testing.T) {
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
	_, err := client.Organization.Federation.Issuers.Get(context.TODO(), "federation_issuer_id")
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOrganizationFederationIssuerUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Organization.Federation.Issuers.Update(
		context.TODO(),
		"federation_issuer_id",
		anthropic.OrganizationFederationIssuerUpdateParams{
			CheckJTI:  anthropic.Bool(true),
			IssuerURL: anthropic.String("x"),
			JWKS: anthropic.OrganizationFederationIssuerUpdateParamsJWKSUnion{
				OfDiscovery: &anthropic.JWKSDiscoveryParam{
					CACertPEM:     anthropic.String("ca_cert_pem"),
					DiscoveryBase: anthropic.String("discovery_base"),
				},
			},
			JWKSPollingDisabled:   anthropic.Bool(true),
			MaxJWTLifetimeSeconds: anthropic.Int(1),
			Name:                  anthropic.String("x"),
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

func TestOrganizationFederationIssuerListWithOptionalParams(t *testing.T) {
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
	_, err := client.Organization.Federation.Issuers.List(context.TODO(), anthropic.OrganizationFederationIssuerListParams{
		IncludeArchived: anthropic.Bool(true),
		Limit:           anthropic.Int(1),
		Page:            anthropic.String("page"),
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOrganizationFederationIssuerArchive(t *testing.T) {
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
	_, err := client.Organization.Federation.Issuers.Archive(context.TODO(), "federation_issuer_id")
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
