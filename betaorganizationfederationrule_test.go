package anthropic_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

func TestBetaOrganizationFederationRuleNewWithOptionalParams(t *testing.T) {
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
		option.WithQuery("beta", "true"),
	)
	_, err := client.Beta.Organization.Federation.Rules.New(context.TODO(), anthropic.BetaOrganizationFederationRuleNewParams{
		IssuerID: "fdis_01SDCCSbTxrXDpWc1phhtcfK",
		Match: anthropic.BetaFederationRuleMatchParam{
			Audience: anthropic.String("audience"),
			Claims: map[string]string{
				"event_name":          "push",
				"ref":                 "refs/heads/main",
				"repository_id":       "123456789",
				"repository_owner_id": "987654321",
			},
			Condition:     anthropic.String("condition"),
			SubjectPrefix: anthropic.String("subject_prefix"),
		},
		Name:       "prod-deploy-pipeline",
		OAuthScope: "workspace:developer",
		Target: anthropic.BetaServiceAccountTargetParam{
			ServiceAccountID:   "svac_01SDCCSbTxrXDpWc1phhtcfK",
			ServiceAccountName: anthropic.String("service_account_name"),
		},
		AppliesToAllWorkspaces: anthropic.Bool(true),
		Attributes: map[string]string{
			"foo": "string",
		},
		Description:          anthropic.String("description"),
		TokenLifetimeSeconds: anthropic.Int(60),
		WorkspaceID:          anthropic.String("wrkspc_01JwQvzr7rXLA5AGx3HKfFUJ"),
		Betas:                []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2024_10_22},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaOrganizationFederationRuleGetWithOptionalParams(t *testing.T) {
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
		option.WithQuery("beta", "true"),
	)
	_, err := client.Beta.Organization.Federation.Rules.Get(
		context.TODO(),
		"federation_rule_id",
		anthropic.BetaOrganizationFederationRuleGetParams{
			Betas: []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2024_10_22},
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

func TestBetaOrganizationFederationRuleUpdateWithOptionalParams(t *testing.T) {
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
		option.WithQuery("beta", "true"),
	)
	_, err := client.Beta.Organization.Federation.Rules.Update(
		context.TODO(),
		"federation_rule_id",
		anthropic.BetaOrganizationFederationRuleUpdateParams{
			AppliesToAllWorkspaces: anthropic.Bool(true),
			Attributes: map[string]string{
				"foo": "string",
			},
			Description: anthropic.String("description"),
			Match: anthropic.BetaFederationRuleMatchParam{
				Audience: anthropic.String("audience"),
				Claims: map[string]string{
					"foo": "string",
				},
				Condition:     anthropic.String("condition"),
				SubjectPrefix: anthropic.String("subject_prefix"),
			},
			Name:       anthropic.String("x"),
			OAuthScope: anthropic.String("x"),
			Target: anthropic.BetaServiceAccountTargetParam{
				ServiceAccountID:   "svac_01SDCCSbTxrXDpWc1phhtcfK",
				ServiceAccountName: anthropic.String("service_account_name"),
			},
			TokenLifetimeSeconds: anthropic.Int(60),
			WorkspaceID:          anthropic.String("workspace_id"),
			Betas:                []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2024_10_22},
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

func TestBetaOrganizationFederationRuleListWithOptionalParams(t *testing.T) {
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
		option.WithQuery("beta", "true"),
	)
	_, err := client.Beta.Organization.Federation.Rules.List(context.TODO(), anthropic.BetaOrganizationFederationRuleListParams{
		IncludeArchived: anthropic.Bool(true),
		IssuerID:        anthropic.String("issuer_id"),
		Limit:           anthropic.Int(1),
		Page:            anthropic.String("page"),
		Betas:           []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2024_10_22},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaOrganizationFederationRuleArchiveWithOptionalParams(t *testing.T) {
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
		option.WithQuery("beta", "true"),
	)
	_, err := client.Beta.Organization.Federation.Rules.Archive(
		context.TODO(),
		"federation_rule_id",
		anthropic.BetaOrganizationFederationRuleArchiveParams{
			Betas: []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2024_10_22},
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
