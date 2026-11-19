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

func TestBetaOrganizationClaudeCodeManagedSettingConfigurationGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.ClaudeCode.ManagedSettings.Configuration.Get(context.TODO(), anthropic.BetaOrganizationClaudeCodeManagedSettingConfigurationGetParams{
		Betas: []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2024_10_22},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaOrganizationClaudeCodeManagedSettingConfigurationUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.ClaudeCode.ManagedSettings.Configuration.Update(context.TODO(), anthropic.BetaOrganizationClaudeCodeManagedSettingConfigurationUpdateParams{
		Mode: anthropic.BetaClaudeCodeManagedSettingsModeParamUnion{
			OfOrganization: &anthropic.BetaClaudeCodeManagedSettingsOrganizationModeParam{},
		},
		RBACGroupOrder:        []string{"string"},
		RBACGroupOrderVersion: anthropic.String("rbac_group_order_version"),
		Betas:                 []anthropic.AnthropicBeta{anthropic.AnthropicBetaComputerUse2024_10_22},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
