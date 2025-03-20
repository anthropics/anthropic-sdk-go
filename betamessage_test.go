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

func TestBetaMessageNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Messages.New(context.TODO(), anthropic.BetaMessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.BetaMessageParam{{
			Content: []anthropic.BetaContentBlockParamUnion{{
				OfRequestTextBlock: &anthropic.BetaTextBlockParam{Text: "What is a quaternion?", CacheControl: anthropic.BetaCacheControlEphemeralParam{}, Citations: []anthropic.BetaTextCitationParamUnion{{
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
		System: []anthropic.BetaTextBlockParam{{Text: "x", CacheControl: anthropic.BetaCacheControlEphemeralParam{}, Citations: []anthropic.BetaTextCitationParamUnion{{
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
				CacheControl: anthropic.BetaCacheControlEphemeralParam{},
				Description:  anthropic.String("Get the current weather in a given location"),
				Type:         anthropic.BetaToolTypeCustom,
			},
		}},
		TopK:  anthropic.Int(5),
		TopP:  anthropic.Float(0.7),
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

func TestBetaMessageCountTokensWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Messages.CountTokens(context.TODO(), anthropic.BetaMessageCountTokensParams{
		Messages: []anthropic.BetaMessageParam{{
			Content: []anthropic.BetaContentBlockParamUnion{{
				OfRequestTextBlock: &anthropic.BetaTextBlockParam{Text: "What is a quaternion?", CacheControl: anthropic.BetaCacheControlEphemeralParam{}, Citations: []anthropic.BetaTextCitationParamUnion{{
					OfRequestCharLocationCitation: &anthropic.BetaCitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
				}}},
			}},
			Role: anthropic.BetaMessageParamRoleUser,
		}},
		Model: anthropic.ModelClaude3_7SonnetLatest,
		System: anthropic.BetaMessageCountTokensParamsSystemUnion{
			OfBetaMessageCountTokenssSystemArray: []anthropic.BetaTextBlockParam{{
				Text:         "Today's date is 2024-06-01.",
				CacheControl: anthropic.BetaCacheControlEphemeralParam{},
				Citations: []anthropic.BetaTextCitationParamUnion{{
					OfRequestCharLocationCitation: &anthropic.BetaCitationCharLocationParam{
						CitedText:      "cited_text",
						DocumentIndex:  0,
						DocumentTitle:  anthropic.String("x"),
						EndCharIndex:   0,
						StartCharIndex: 0,
					},
				}},
			}},
		},
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
		Tools: []anthropic.BetaMessageCountTokensParamsToolUnion{{
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
				CacheControl: anthropic.BetaCacheControlEphemeralParam{},
				Description:  anthropic.String("Get the current weather in a given location"),
				Type:         anthropic.BetaToolTypeCustom,
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
