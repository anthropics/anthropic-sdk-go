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

func TestBetaMessageBatchNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Messages.Batches.New(context.TODO(), anthropic.BetaMessageBatchNewParams{
		Requests: []anthropic.BetaMessageBatchNewParamsRequest{{
			CustomID: "my-custom-id-1",
			Params: anthropic.BetaMessageBatchNewParamsRequestParams{
				MaxTokens: 1024,
				Messages: []anthropic.BetaMessageParam{{
					Content: []anthropic.BetaContentBlockParamUnion{{
						OfRequestTextBlock: &anthropic.BetaTextBlockParam{Text: "What is a quaternion?", CacheControl: anthropic.NewBetaCacheControlEphemeralParam(), Citations: []anthropic.BetaTextCitationParamUnion{{
							OfRequestCharLocationCitation: &anthropic.BetaCitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
						}}},
					}},
					Role: anthropic.BetaMessageParamRoleUser,
				}},
				Model: anthropic.ModelClaude3_7SonnetLatest,
				Metadata: anthropic.BetaMetadataParam{
					UserID: anthropic.String("13803d75-b4b5-4c3e-b2a2-6f21399b021b"),
				},
				StopSequences: []string{"string"},
				Stream:        anthropic.Bool(true),
				System: []anthropic.BetaTextBlockParam{{Text: "x", CacheControl: anthropic.NewBetaCacheControlEphemeralParam(), Citations: []anthropic.BetaTextCitationParamUnion{{
					OfRequestCharLocationCitation: &anthropic.BetaCitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
				}}}},
				Temperature: anthropic.Float(1),
				Thinking: anthropic.BetaThinkingConfigParamUnion{
					OfThinkingConfigEnabled: &anthropic.BetaThinkingConfigEnabledParam{
						BudgetTokens: 1024,
					},
				},
				ToolChoice: anthropic.BetaToolChoiceUnionParam{
					OfToolChoiceAuto: &anthropic.BetaToolChoiceAutoParam{
						DisableParallelToolUse: anthropic.Bool(true),
					},
				},
				Tools: []anthropic.BetaToolUnionParam{{
					OfTool: &anthropic.BetaToolParam{
						InputSchema: anthropic.BetaToolInputSchemaParam{
							Properties: map[string]interface{}{
								"location": map[string]interface{}{
									"description": "The city and state, e.g. San Francisco, CA",
									"type":        "string",
								},
								"unit": map[string]interface{}{
									"description": "Unit for the output - one of (celsius, fahrenheit)",
									"type":        "string",
								},
							},
						},
						Name:         "name",
						CacheControl: anthropic.NewBetaCacheControlEphemeralParam(),
						Description:  anthropic.String("Get the current weather in a given location"),
						Type:         anthropic.BetaToolTypeCustom,
					},
				}},
				TopK: anthropic.Int(5),
				TopP: anthropic.Float(0.7),
			},
		}},
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

func TestBetaMessageBatchGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Messages.Batches.Get(
		context.TODO(),
		"message_batch_id",
		anthropic.BetaMessageBatchGetParams{
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

func TestBetaMessageBatchListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Messages.Batches.List(context.TODO(), anthropic.BetaMessageBatchListParams{
		AfterID:  anthropic.String("after_id"),
		BeforeID: anthropic.String("before_id"),
		Limit:    anthropic.Int(1),
		Betas:    []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaMessageBatchDeleteWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Messages.Batches.Delete(
		context.TODO(),
		"message_batch_id",
		anthropic.BetaMessageBatchDeleteParams{
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

func TestBetaMessageBatchCancelWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Messages.Batches.Cancel(
		context.TODO(),
		"message_batch_id",
		anthropic.BetaMessageBatchCancelParams{
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
