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
	t.Skip("prism validates based on the non-beta endpoint")
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
				OfText: &anthropic.BetaTextBlockParam{Text: "What is a quaternion?", CacheControl: anthropic.BetaCacheControlEphemeralParam{TTL: anthropic.BetaCacheControlEphemeralTTLTTL5m}, Citations: []anthropic.BetaTextCitationParamUnion{{
					OfCharLocation: &anthropic.BetaCitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
				}}},
			}},
			Role: anthropic.BetaMessageParamRoleUser,
		}},
		Model: anthropic.ModelClaude3_7SonnetLatest,
		Container: anthropic.BetaMessageNewParamsContainerUnion{
			OfContainers: &anthropic.BetaContainerParams{
				ID: anthropic.String("id"),
				Skills: []anthropic.BetaSkillParams{{
					SkillID: "x",
					Type:    anthropic.BetaSkillParamsTypeAnthropic,
					Version: anthropic.String("x"),
				}},
			},
		},
		ContextManagement: anthropic.BetaContextManagementConfigParam{
			Edits: []anthropic.BetaContextManagementConfigEditUnionParam{{
				OfClearToolUses20250919: &anthropic.BetaClearToolUses20250919EditParam{
					ClearAtLeast: anthropic.BetaInputTokensClearAtLeastParam{
						Value: 0,
					},
					ClearToolInputs: anthropic.BetaClearToolUses20250919EditClearToolInputsUnionParam{
						OfBool: anthropic.Bool(true),
					},
					ExcludeTools: []string{"string"},
					Keep: anthropic.BetaToolUsesKeepParam{
						Value: 0,
					},
					Trigger: anthropic.BetaClearToolUses20250919EditTriggerUnionParam{
						OfInputTokens: &anthropic.BetaInputTokensTriggerParam{
							Value: 1,
						},
					},
				},
			}},
		},
		MCPServers: []anthropic.BetaRequestMCPServerURLDefinitionParam{{
			Name:               "name",
			URL:                "url",
			AuthorizationToken: anthropic.String("authorization_token"),
			ToolConfiguration: anthropic.BetaRequestMCPServerToolConfigurationParam{
				AllowedTools: []string{"string"},
				Enabled:      anthropic.Bool(true),
			},
		}},
		Metadata: anthropic.BetaMetadataParam{
			UserID: anthropic.String("13803d75-b4b5-4c3e-b2a2-6f21399b021b"),
		},
		OutputFormat: anthropic.BetaJSONOutputFormatParam{
			Schema: map[string]any{
				"foo": "bar",
			},
		},
		ServiceTier:   anthropic.BetaMessageNewParamsServiceTierAuto,
		StopSequences: []string{"string"},
		System: []anthropic.BetaTextBlockParam{{Text: "x", CacheControl: anthropic.BetaCacheControlEphemeralParam{TTL: anthropic.BetaCacheControlEphemeralTTLTTL5m}, Citations: []anthropic.BetaTextCitationParamUnion{{
			OfCharLocation: &anthropic.BetaCitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
		}}}},
		Temperature: anthropic.Float(1),
		Thinking: anthropic.BetaThinkingConfigParamUnion{
			OfEnabled: &anthropic.BetaThinkingConfigEnabledParam{
				BudgetTokens: 1024,
			},
		},
		ToolChoice: anthropic.BetaToolChoiceUnionParam{
			OfAuto: &anthropic.BetaToolChoiceAutoParam{
				DisableParallelToolUse: anthropic.Bool(true),
			},
		},
		Tools: []anthropic.BetaToolUnionParam{{
			OfTool: &anthropic.BetaToolParam{
				InputSchema: anthropic.BetaToolInputSchemaParam{
					Properties: map[string]any{
						"location": "bar",
						"unit":     "bar",
					},
					Required: []string{"location"},
				},
				Name: "name",
				CacheControl: anthropic.BetaCacheControlEphemeralParam{
					TTL: anthropic.BetaCacheControlEphemeralTTLTTL5m,
				},
				Description: anthropic.String("Get the current weather in a given location"),
				Strict:      anthropic.Bool(true),
				Type:        anthropic.BetaToolTypeCustom,
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
	t.Skip("prism validates based on the non-beta endpoint")
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
				OfText: &anthropic.BetaTextBlockParam{Text: "What is a quaternion?", CacheControl: anthropic.BetaCacheControlEphemeralParam{TTL: anthropic.BetaCacheControlEphemeralTTLTTL5m}, Citations: []anthropic.BetaTextCitationParamUnion{{
					OfCharLocation: &anthropic.BetaCitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
				}}},
			}},
			Role: anthropic.BetaMessageParamRoleUser,
		}},
		Model: anthropic.ModelClaude3_7SonnetLatest,
		ContextManagement: anthropic.BetaContextManagementConfigParam{
			Edits: []anthropic.BetaContextManagementConfigEditUnionParam{{
				OfClearToolUses20250919: &anthropic.BetaClearToolUses20250919EditParam{
					ClearAtLeast: anthropic.BetaInputTokensClearAtLeastParam{
						Value: 0,
					},
					ClearToolInputs: anthropic.BetaClearToolUses20250919EditClearToolInputsUnionParam{
						OfBool: anthropic.Bool(true),
					},
					ExcludeTools: []string{"string"},
					Keep: anthropic.BetaToolUsesKeepParam{
						Value: 0,
					},
					Trigger: anthropic.BetaClearToolUses20250919EditTriggerUnionParam{
						OfInputTokens: &anthropic.BetaInputTokensTriggerParam{
							Value: 1,
						},
					},
				},
			}},
		},
		MCPServers: []anthropic.BetaRequestMCPServerURLDefinitionParam{{
			Name:               "name",
			URL:                "url",
			AuthorizationToken: anthropic.String("authorization_token"),
			ToolConfiguration: anthropic.BetaRequestMCPServerToolConfigurationParam{
				AllowedTools: []string{"string"},
				Enabled:      anthropic.Bool(true),
			},
		}},
		OutputFormat: anthropic.BetaJSONOutputFormatParam{
			Schema: map[string]any{
				"foo": "bar",
			},
		},
		System: anthropic.BetaMessageCountTokensParamsSystemUnion{
			OfBetaTextBlockArray: []anthropic.BetaTextBlockParam{{
				Text: "Today's date is 2024-06-01.",
				CacheControl: anthropic.BetaCacheControlEphemeralParam{
					TTL: anthropic.BetaCacheControlEphemeralTTLTTL5m,
				},
				Citations: []anthropic.BetaTextCitationParamUnion{{
					OfCharLocation: &anthropic.BetaCitationCharLocationParam{
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
			OfEnabled: &anthropic.BetaThinkingConfigEnabledParam{
				BudgetTokens: 1024,
			},
		},
		ToolChoice: anthropic.BetaToolChoiceUnionParam{
			OfAuto: &anthropic.BetaToolChoiceAutoParam{
				DisableParallelToolUse: anthropic.Bool(true),
			},
		},
		Tools: []anthropic.BetaMessageCountTokensParamsToolUnion{{
			OfTool: &anthropic.BetaToolParam{
				InputSchema: anthropic.BetaToolInputSchemaParam{
					Properties: map[string]any{
						"location": "bar",
						"unit":     "bar",
					},
					Required: []string{"location"},
				},
				Name: "name",
				CacheControl: anthropic.BetaCacheControlEphemeralParam{
					TTL: anthropic.BetaCacheControlEphemeralTTLTTL5m,
				},
				Description: anthropic.String("Get the current weather in a given location"),
				Strict:      anthropic.Bool(true),
				Type:        anthropic.BetaToolTypeCustom,
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
