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

func TestBetaOrganizationClaudeTagBundleTargetNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.ClaudeTag.Bundles.Targets.New(
		context.TODO(),
		"capp_011CZkZjLP1cWxj5EtMmdgbq",
		anthropic.BetaOrganizationClaudeTagBundleTargetNewParams{
			Target: anthropic.BetaClaudeTagBundleTargetCreateUnionParam{
				OfSlackChannel: &anthropic.BetaClaudeTagSlackChannelBundleTargetCreateParam{
					SlackChannelID: "xxxxxxxx",
				},
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

func TestBetaOrganizationClaudeTagBundleTargetGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.ClaudeTag.Bundles.Targets.Get(
		context.TODO(),
		"btgt_01SvpyRQK3nN6eH9uhRA7n1qHq8TzW5mKc2YbX7vNd4RfJ",
		anthropic.BetaOrganizationClaudeTagBundleTargetGetParams{
			BundleID: "capp_011CZkZjLP1cWxj5EtMmdgbq",
			Betas:    []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
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

func TestBetaOrganizationClaudeTagBundleTargetListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.ClaudeTag.Bundles.Targets.List(
		context.TODO(),
		"capp_011CZkZjLP1cWxj5EtMmdgbq",
		anthropic.BetaOrganizationClaudeTagBundleTargetListParams{
			Limit: anthropic.Int(0),
			Page:  anthropic.String("page"),
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

func TestBetaOrganizationClaudeTagBundleTargetDeleteWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.ClaudeTag.Bundles.Targets.Delete(
		context.TODO(),
		"btgt_01SvpyRQK3nN6eH9uhRA7n1qHq8TzW5mKc2YbX7vNd4RfJ",
		anthropic.BetaOrganizationClaudeTagBundleTargetDeleteParams{
			BundleID: "capp_011CZkZjLP1cWxj5EtMmdgbq",
			Betas:    []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
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
