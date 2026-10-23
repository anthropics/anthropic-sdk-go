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

func TestBetaOrganizationClaudeTagSlackChannelNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.ClaudeTag.SlackChannels.New(context.TODO(), anthropic.BetaOrganizationClaudeTagSlackChannelNewParams{
		Slack: anthropic.BetaClaudeTagSlackChannelIdentifiersParam{
			ChannelID:   "xx",
			WorkspaceID: "xx",
		},
		Settings: anthropic.BetaClaudeTagSettingsCreateParam{
			Enabled:       anthropic.Bool(true),
			EnvironmentID: anthropic.String("x"),
			GuestPolicy: anthropic.BetaClaudeTagGuestPolicyUnionParam{
				OfAllow: &anthropic.BetaClaudeTagGuestPolicyAllowParam{},
			},
			Instructions: anthropic.String("instructions"),
			Model:        anthropic.String("x"),
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

func TestBetaOrganizationClaudeTagSlackChannelGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.ClaudeTag.SlackChannels.Get(
		context.TODO(),
		"slch_01SvpyRQK3nN6eH9uhRA7n1q",
		anthropic.BetaOrganizationClaudeTagSlackChannelGetParams{
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

func TestBetaOrganizationClaudeTagSlackChannelUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.ClaudeTag.SlackChannels.Update(
		context.TODO(),
		"slch_01SvpyRQK3nN6eH9uhRA7n1q",
		anthropic.BetaOrganizationClaudeTagSlackChannelUpdateParams{
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

func TestBetaOrganizationClaudeTagSlackChannelListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.ClaudeTag.SlackChannels.List(context.TODO(), anthropic.BetaOrganizationClaudeTagSlackChannelListParams{
		SlackWorkspaceID: "slack_workspace_id",
		IncludeInherited: anthropic.Bool(true),
		Limit:            anthropic.Int(0),
		Page:             anthropic.String("page"),
		SlackChannelID:   anthropic.String("slack_channel_id"),
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

func TestBetaOrganizationClaudeTagSlackChannelDeleteWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.ClaudeTag.SlackChannels.Delete(
		context.TODO(),
		"slch_01SvpyRQK3nN6eH9uhRA7n1q",
		anthropic.BetaOrganizationClaudeTagSlackChannelDeleteParams{
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
