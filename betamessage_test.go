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
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.BetaMessageParam{{
			Content: anthropic.F([]anthropic.BetaContentBlockParamUnion{anthropic.BetaTextBlockParam{Text: anthropic.F("What is a quaternion?"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)})}}),
			Role:    anthropic.F(anthropic.BetaMessageParamRoleUser),
		}}),
		Model: anthropic.F(anthropic.ModelClaude3_5HaikuLatest),
		Metadata: anthropic.F(anthropic.BetaMetadataParam{
			UserID: anthropic.F("13803d75-b4b5-4c3e-b2a2-6f21399b021b"),
		}),
		StopSequences: anthropic.F([]string{"string"}),
		System:        anthropic.F([]anthropic.BetaTextBlockParam{{Text: anthropic.F("x"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)})}}),
		Temperature:   anthropic.F(1.000000),
		ToolChoice: anthropic.F[anthropic.BetaToolChoiceUnionParam](anthropic.BetaToolChoiceAutoParam{
			Type:                   anthropic.F(anthropic.BetaToolChoiceAutoTypeAuto),
			DisableParallelToolUse: anthropic.F(true),
		}),
		Tools: anthropic.F([]anthropic.BetaToolUnionUnionParam{anthropic.BetaToolParam{
			InputSchema: anthropic.F(anthropic.BetaToolInputSchemaParam{
				Type: anthropic.F(anthropic.BetaToolInputSchemaTypeObject),
				Properties: anthropic.F[any](map[string]interface{}{
					"location": map[string]interface{}{
						"description": "The city and state, e.g. San Francisco, CA",
						"type":        "string",
					},
					"unit": map[string]interface{}{
						"description": "Unit for the output - one of (celsius, fahrenheit)",
						"type":        "string",
					},
				}),
			}),
			Name: anthropic.F("name"),
			CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{
				Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral),
			}),
			Description: anthropic.F("Get the current weather in a given location"),
			Type:        anthropic.F(anthropic.BetaToolTypeCustom),
		}}),
		TopK:  anthropic.F(int64(5)),
		TopP:  anthropic.F(0.700000),
		Betas: anthropic.F([]anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24}),
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
		Messages: anthropic.F([]anthropic.BetaMessageParam{{
			Content: anthropic.F([]anthropic.BetaContentBlockParamUnion{anthropic.BetaTextBlockParam{Text: anthropic.F("What is a quaternion?"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)})}}),
			Role:    anthropic.F(anthropic.BetaMessageParamRoleUser),
		}}),
		Model: anthropic.F(anthropic.ModelClaude3_5HaikuLatest),
		System: anthropic.F[anthropic.BetaMessageCountTokensParamsSystemUnion](anthropic.BetaMessageCountTokensParamsSystemArray([]anthropic.BetaTextBlockParam{{
			Text: anthropic.F("Today's date is 2024-06-01."),
			Type: anthropic.F(anthropic.BetaTextBlockParamTypeText),
			CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{
				Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral),
			}),
		}})),
		ToolChoice: anthropic.F[anthropic.BetaToolChoiceUnionParam](anthropic.BetaToolChoiceAutoParam{
			Type:                   anthropic.F(anthropic.BetaToolChoiceAutoTypeAuto),
			DisableParallelToolUse: anthropic.F(true),
		}),
		Tools: anthropic.F([]anthropic.BetaMessageCountTokensParamsToolUnion{anthropic.BetaToolParam{
			InputSchema: anthropic.F(anthropic.BetaToolInputSchemaParam{
				Type: anthropic.F(anthropic.BetaToolInputSchemaTypeObject),
				Properties: anthropic.F[any](map[string]interface{}{
					"location": map[string]interface{}{
						"description": "The city and state, e.g. San Francisco, CA",
						"type":        "string",
					},
					"unit": map[string]interface{}{
						"description": "Unit for the output - one of (celsius, fahrenheit)",
						"type":        "string",
					},
				}),
			}),
			Name: anthropic.F("name"),
			CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{
				Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral),
			}),
			Description: anthropic.F("Get the current weather in a given location"),
			Type:        anthropic.F(anthropic.BetaToolTypeCustom),
		}}),
		Betas: anthropic.F([]anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24}),
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
