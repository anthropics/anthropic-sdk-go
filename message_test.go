// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestMessageParamMarshalUnmarshal(t *testing.T) {
	original := anthropic.MessageParam{
		Role: anthropic.MessageParamRoleUser,
		Content: []anthropic.ContentBlockParamUnion{
			anthropic.ContentBlockParamOfRequestTextBlock("Hello, world!"),
			anthropic.ContentBlockParamOfRequestTextBlock("This is a second message"),
		},
	}

	// Marshal the original to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal MessageParam: %v", err)
	}

	// Unmarshal back to a new struct
	var unmarshaled anthropic.MessageParam
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal MessageParam: %v", err)
	}

	// Verify that the content was properly preserved
	if len(unmarshaled.Content) != len(original.Content) {
		t.Errorf("Content length mismatch. Original: %d, Unmarshaled: %d",
			len(original.Content), len(unmarshaled.Content))
	}

	// Check text content of each block
	for i, originalBlock := range original.Content {
		if i >= len(unmarshaled.Content) {
			t.Fatalf("Missing content block at index %d", i)
		}

		originalText := originalBlock.GetText()
		unmarshaledText := unmarshaled.Content[i].GetText()

		if originalText == nil || unmarshaledText == nil {
			t.Errorf("Text is nil at index %d. Original: %v, Unmarshaled: %v",
				i, originalText, unmarshaledText)
			continue
		}

		if *originalText != *unmarshaledText {
			t.Errorf("Content mismatch at index %d. Expected: %q, Got: %q",
				i, *originalText, *unmarshaledText)
		}
	}
}

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
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfRequestTextBlock: &anthropic.TextBlockParam{Text: "What is a quaternion?", CacheControl: anthropic.CacheControlEphemeralParam{}, Citations: []anthropic.TextCitationParamUnion{{
					OfRequestCharLocationCitation: &anthropic.CitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
				}}},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaude3_7SonnetLatest,
		Metadata: anthropic.MetadataParam{
			UserID: anthropic.String("13803d75-b4b5-4c3e-b2a2-6f21399b021b"),
		},
		StopSequences: []string{"string"},
		System: []anthropic.TextBlockParam{{Text: "x", CacheControl: anthropic.CacheControlEphemeralParam{}, Citations: []anthropic.TextCitationParamUnion{{
			OfRequestCharLocationCitation: &anthropic.CitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
		}}}},
		Temperature: anthropic.Float(1),
		Thinking: anthropic.ThinkingConfigParamUnion{
			OfThinkingConfigEnabled: &anthropic.ThinkingConfigEnabledParam{
				BudgetTokens: 1024,
			},
		},
		ToolChoice: anthropic.ToolChoiceUnionParam{
			OfToolChoiceAuto: &anthropic.ToolChoiceAutoParam{
				DisableParallelToolUse: anthropic.Bool(true),
			},
		},
		Tools: []anthropic.ToolUnionParam{{
			OfTool: &anthropic.ToolParam{
				InputSchema: anthropic.ToolInputSchemaParam{
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
				CacheControl: anthropic.CacheControlEphemeralParam{},
				Description:  anthropic.String("Get the current weather in a given location"),
			},
		}},
		TopK: anthropic.Int(5),
		TopP: anthropic.Float(0.7),
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestMessageCountTokensWithOptionalParams(t *testing.T) {
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
	_, err := client.Messages.CountTokens(context.TODO(), anthropic.MessageCountTokensParams{
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfRequestTextBlock: &anthropic.TextBlockParam{Text: "What is a quaternion?", CacheControl: anthropic.CacheControlEphemeralParam{}, Citations: []anthropic.TextCitationParamUnion{{
					OfRequestCharLocationCitation: &anthropic.CitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
				}}},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaude3_7SonnetLatest,
		System: anthropic.MessageCountTokensParamsSystemUnion{
			OfMessageCountTokenssSystemArray: []anthropic.TextBlockParam{{
				Text:         "Today's date is 2024-06-01.",
				CacheControl: anthropic.CacheControlEphemeralParam{},
				Citations: []anthropic.TextCitationParamUnion{{
					OfRequestCharLocationCitation: &anthropic.CitationCharLocationParam{
						CitedText:      "cited_text",
						DocumentIndex:  0,
						DocumentTitle:  anthropic.String("x"),
						EndCharIndex:   0,
						StartCharIndex: 0,
					},
				}},
			}},
		},
		Thinking: anthropic.ThinkingConfigParamUnion{
			OfThinkingConfigEnabled: &anthropic.ThinkingConfigEnabledParam{
				BudgetTokens: 1024,
			},
		},
		ToolChoice: anthropic.ToolChoiceUnionParam{
			OfToolChoiceAuto: &anthropic.ToolChoiceAutoParam{
				DisableParallelToolUse: anthropic.Bool(true),
			},
		},
		Tools: []anthropic.MessageCountTokensToolUnionParam{{
			OfTool: &anthropic.ToolParam{
				InputSchema: anthropic.ToolInputSchemaParam{
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
				CacheControl: anthropic.CacheControlEphemeralParam{},
				Description:  anthropic.String("Get the current weather in a given location"),
			},
		}},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
