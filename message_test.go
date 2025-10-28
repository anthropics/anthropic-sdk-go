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
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfText: &anthropic.TextBlockParam{Text: "What is a quaternion?", CacheControl: anthropic.CacheControlEphemeralParam{TTL: anthropic.CacheControlEphemeralTTLTTL5m}, Citations: []anthropic.TextCitationParamUnion{{
					OfCharLocation: &anthropic.CitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
				}}},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaude3_7SonnetLatest,
		Metadata: anthropic.MetadataParam{
			UserID: anthropic.String("13803d75-b4b5-4c3e-b2a2-6f21399b021b"),
		},
		ServiceTier:   anthropic.MessageNewParamsServiceTierAuto,
		StopSequences: []string{"string"},
		System: []anthropic.TextBlockParam{{Text: "x", CacheControl: anthropic.CacheControlEphemeralParam{TTL: anthropic.CacheControlEphemeralTTLTTL5m}, Citations: []anthropic.TextCitationParamUnion{{
			OfCharLocation: &anthropic.CitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
		}}}},
		Temperature: anthropic.Float(1),
		Thinking: anthropic.ThinkingConfigParamUnion{
			OfEnabled: &anthropic.ThinkingConfigEnabledParam{
				BudgetTokens: 1024,
			},
		},
		ToolChoice: anthropic.ToolChoiceUnionParam{
			OfAuto: &anthropic.ToolChoiceAutoParam{
				DisableParallelToolUse: anthropic.Bool(true),
			},
		},
		Tools: []anthropic.ToolUnionParam{{
			OfTool: &anthropic.ToolParam{
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"location": "bar",
						"unit":     "bar",
					},
					Required: []string{"location"},
				},
				Name: "name",
				CacheControl: anthropic.CacheControlEphemeralParam{
					TTL: anthropic.CacheControlEphemeralTTLTTL5m,
				},
				Description: anthropic.String("Get the current weather in a given location"),
				Type:        anthropic.ToolTypeCustom,
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
				OfText: &anthropic.TextBlockParam{Text: "What is a quaternion?", CacheControl: anthropic.CacheControlEphemeralParam{TTL: anthropic.CacheControlEphemeralTTLTTL5m}, Citations: []anthropic.TextCitationParamUnion{{
					OfCharLocation: &anthropic.CitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
				}}},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaude3_7SonnetLatest,
		System: anthropic.MessageCountTokensParamsSystemUnion{
			OfTextBlockArray: []anthropic.TextBlockParam{{
				Text: "Today's date is 2024-06-01.",
				CacheControl: anthropic.CacheControlEphemeralParam{
					TTL: anthropic.CacheControlEphemeralTTLTTL5m,
				},
				Citations: []anthropic.TextCitationParamUnion{{
					OfCharLocation: &anthropic.CitationCharLocationParam{
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
			OfEnabled: &anthropic.ThinkingConfigEnabledParam{
				BudgetTokens: 1024,
			},
		},
		ToolChoice: anthropic.ToolChoiceUnionParam{
			OfAuto: &anthropic.ToolChoiceAutoParam{
				DisableParallelToolUse: anthropic.Bool(true),
			},
		},
		Tools: []anthropic.MessageCountTokensToolUnionParam{{
			OfTool: &anthropic.ToolParam{
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"location": "bar",
						"unit":     "bar",
					},
					Required: []string{"location"},
				},
				Name: "name",
				CacheControl: anthropic.CacheControlEphemeralParam{
					TTL: anthropic.CacheControlEphemeralTTLTTL5m,
				},
				Description: anthropic.String("Get the current weather in a given location"),
				Type:        anthropic.ToolTypeCustom,
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
