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

func TestBetaOrganizationClaudeTagSlackWorkspaceGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.ClaudeTag.SlackWorkspaces.Get(
		context.TODO(),
		"slws_01Hq8TzW5mKc2YbX7vNd4RfJ",
		anthropic.BetaOrganizationClaudeTagSlackWorkspaceGetParams{
			IncludeInherited: anthropic.Bool(true),
			Betas:            []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
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

func TestBetaOrganizationClaudeTagSlackWorkspaceUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.ClaudeTag.SlackWorkspaces.Update(
		context.TODO(),
		"slws_01Hq8TzW5mKc2YbX7vNd4RfJ",
		anthropic.BetaOrganizationClaudeTagSlackWorkspaceUpdateParams{
			Settings: anthropic.BetaClaudeTagSettingsUpdateParam{
				Enabled:       anthropic.Bool(true),
				EnvironmentID: anthropic.String("x"),
				GuestPolicy: anthropic.BetaClaudeTagGuestPolicyUnionParam{
					OfAllow: &anthropic.BetaClaudeTagGuestPolicyAllowParam{},
				},
				Instructions: anthropic.String("instructions"),
				Model:        anthropic.String("x"),
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

func TestBetaOrganizationClaudeTagSlackWorkspaceListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.ClaudeTag.SlackWorkspaces.List(context.TODO(), anthropic.BetaOrganizationClaudeTagSlackWorkspaceListParams{
		IncludeInherited: anthropic.Bool(true),
		Limit:            anthropic.Int(0),
		Page:             anthropic.String("page"),
		SlackWorkspaceID: anthropic.String("slack_workspace_id"),
		Betas:            []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
