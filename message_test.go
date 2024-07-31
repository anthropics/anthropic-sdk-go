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

func TestMessageNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.MessageParam{{
			Role:    anthropic.F(anthropic.MessageParamRoleUser),
			Content: anthropic.F([]anthropic.MessageParamContentUnion{anthropic.TextBlockParam{Type: anthropic.F(anthropic.TextBlockParamTypeText), Text: anthropic.F("What is a quaternion?")}}),
		}}),
		Model: anthropic.F(anthropic.ModelClaude_3_5_Sonnet_20240620),
		Metadata: anthropic.F(anthropic.MessageNewParamsMetadata{
			UserID: anthropic.F("13803d75-b4b5-4c3e-b2a2-6f21399b021b"),
		}),
		StopSequences: anthropic.F([]string{"string", "string", "string"}),
		System: anthropic.F[anthropic.MessageNewParamsSystemUnion](anthropic.MessageNewParamsSystemArray([]anthropic.TextBlockParam{{
			Type: anthropic.F(anthropic.TextBlockParamTypeText),
			Text: anthropic.F("Today's date is 2024-06-01."),
		}})),
		Temperature: anthropic.F(1.000000),
		ToolChoice: anthropic.F[anthropic.MessageNewParamsToolChoiceUnion](anthropic.MessageNewParamsToolChoiceToolChoiceAuto{
			Type: anthropic.F(anthropic.MessageNewParamsToolChoiceToolChoiceAutoTypeAuto),
		}),
		Tools: anthropic.F([]anthropic.ToolParam{{
			Description: anthropic.F("Get the current weather in a given location"),
			Name:        anthropic.F("x"),
			InputSchema: anthropic.F(anthropic.ToolInputSchemaParam{
				Type: anthropic.F(anthropic.ToolInputSchemaTypeObject),
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
		}, {
			Description: anthropic.F("Get the current weather in a given location"),
			Name:        anthropic.F("x"),
			InputSchema: anthropic.F(anthropic.ToolInputSchemaParam{
				Type: anthropic.F(anthropic.ToolInputSchemaTypeObject),
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
		}, {
			Description: anthropic.F("Get the current weather in a given location"),
			Name:        anthropic.F("x"),
			InputSchema: anthropic.F(anthropic.ToolInputSchemaParam{
				Type: anthropic.F(anthropic.ToolInputSchemaTypeObject),
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
