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

func TestBetaPromptCachingMessageNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.PromptCaching.Messages.New(context.TODO(), anthropic.BetaPromptCachingMessageNewParams{
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.PromptCachingBetaMessageParam{{
			Content: anthropic.F([]anthropic.PromptCachingBetaMessageParamContentUnion{anthropic.PromptCachingBetaTextBlockParam{Text: anthropic.F("What is a quaternion?"), Type: anthropic.F(anthropic.PromptCachingBetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.PromptCachingBetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.PromptCachingBetaCacheControlEphemeralTypeEphemeral)})}}),
			Role:    anthropic.F(anthropic.PromptCachingBetaMessageParamRoleUser),
		}}),
		Model: anthropic.F(anthropic.ModelClaude_3_5_Sonnet_20240620),
		Metadata: anthropic.F(anthropic.BetaPromptCachingMessageNewParamsMetadata{
			UserID: anthropic.F("13803d75-b4b5-4c3e-b2a2-6f21399b021b"),
		}),
		StopSequences: anthropic.F([]string{"string", "string", "string"}),
		System: anthropic.F[anthropic.BetaPromptCachingMessageNewParamsSystemUnion](anthropic.BetaPromptCachingMessageNewParamsSystemArray([]anthropic.PromptCachingBetaTextBlockParam{{
			Text: anthropic.F("Today's date is 2024-06-01."),
			Type: anthropic.F(anthropic.PromptCachingBetaTextBlockParamTypeText),
			CacheControl: anthropic.F(anthropic.PromptCachingBetaCacheControlEphemeralParam{
				Type: anthropic.F(anthropic.PromptCachingBetaCacheControlEphemeralTypeEphemeral),
			}),
		}})),
		Temperature: anthropic.F(1.000000),
		ToolChoice: anthropic.F[anthropic.BetaPromptCachingMessageNewParamsToolChoiceUnion](anthropic.BetaPromptCachingMessageNewParamsToolChoiceToolChoiceAuto{
			Type:                   anthropic.F(anthropic.BetaPromptCachingMessageNewParamsToolChoiceToolChoiceAutoTypeAuto),
			DisableParallelToolUse: anthropic.F(true),
		}),
		Tools: anthropic.F([]anthropic.PromptCachingBetaToolParam{{
			InputSchema: anthropic.F(anthropic.PromptCachingBetaToolInputSchemaParam{
				Type: anthropic.F(anthropic.PromptCachingBetaToolInputSchemaTypeObject),
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
			Name: anthropic.F("x"),
			CacheControl: anthropic.F(anthropic.PromptCachingBetaCacheControlEphemeralParam{
				Type: anthropic.F(anthropic.PromptCachingBetaCacheControlEphemeralTypeEphemeral),
			}),
			Description: anthropic.F("Get the current weather in a given location"),
		}, {
			InputSchema: anthropic.F(anthropic.PromptCachingBetaToolInputSchemaParam{
				Type: anthropic.F(anthropic.PromptCachingBetaToolInputSchemaTypeObject),
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
			Name: anthropic.F("x"),
			CacheControl: anthropic.F(anthropic.PromptCachingBetaCacheControlEphemeralParam{
				Type: anthropic.F(anthropic.PromptCachingBetaCacheControlEphemeralTypeEphemeral),
			}),
			Description: anthropic.F("Get the current weather in a given location"),
		}, {
			InputSchema: anthropic.F(anthropic.PromptCachingBetaToolInputSchemaParam{
				Type: anthropic.F(anthropic.PromptCachingBetaToolInputSchemaTypeObject),
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
			Name: anthropic.F("x"),
			CacheControl: anthropic.F(anthropic.PromptCachingBetaCacheControlEphemeralParam{
				Type: anthropic.F(anthropic.PromptCachingBetaCacheControlEphemeralTypeEphemeral),
			}),
			Description: anthropic.F("Get the current weather in a given location"),
		}}),
		TopK: anthropic.F(int64(5)),
		TopP: anthropic.F(0.700000),
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
