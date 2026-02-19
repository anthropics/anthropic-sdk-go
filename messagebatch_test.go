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

func TestMessageBatchNew(t *testing.T) {
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
	_, err := client.Messages.Batches.New(context.TODO(), anthropic.MessageBatchNewParams{
		Requests: []anthropic.MessageBatchNewParamsRequest{{
			CustomID: "my-custom-id-1",
			Params: anthropic.MessageBatchNewParamsRequestParams{
				MaxTokens: 1024,
				Messages: []anthropic.MessageParam{{
					Content: []anthropic.ContentBlockParamUnion{{
						OfText: &anthropic.TextBlockParam{
							Text: "x",
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
						},
					}},
					Role: anthropic.MessageParamRoleUser,
				}},
				Model: anthropic.ModelClaudeOpus4_6,
				CacheControl: anthropic.CacheControlEphemeralParam{
					TTL: anthropic.CacheControlEphemeralTTLTTL5m,
				},
				Container:    anthropic.String("container"),
				InferenceGeo: anthropic.String("inference_geo"),
				Metadata: anthropic.MetadataParam{
					UserID: anthropic.String("13803d75-b4b5-4c3e-b2a2-6f21399b021b"),
				},
				OutputConfig: anthropic.OutputConfigParam{
					Effort: anthropic.OutputConfigEffortLow,
					Format: anthropic.JSONOutputFormatParam{
						Schema: map[string]any{
							"foo": "bar",
						},
					},
				},
				ServiceTier:   "auto",
				StopSequences: []string{"string"},
				Stream:        anthropic.Bool(true),
				System: []anthropic.TextBlockParam{{
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
						Name:           "name",
						AllowedCallers: []string{"direct"},
						CacheControl: anthropic.CacheControlEphemeralParam{
							TTL: anthropic.CacheControlEphemeralTTLTTL5m,
						},
						DeferLoading:        anthropic.Bool(true),
						Description:         anthropic.String("Get the current weather in a given location"),
						EagerInputStreaming: anthropic.Bool(true),
						InputExamples: []map[string]any{{
							"foo": "bar",
						}},
						Strict: anthropic.Bool(true),
						Type:   anthropic.ToolTypeCustom,
					},
				}},
				TopK: anthropic.Int(5),
				TopP: anthropic.Float(0.7),
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

func TestMessageBatchGet(t *testing.T) {
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
	_, err := client.Messages.Batches.Get(context.TODO(), "message_batch_id")
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestMessageBatchListWithOptionalParams(t *testing.T) {
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
	_, err := client.Messages.Batches.List(context.TODO(), anthropic.MessageBatchListParams{
		AfterID:  anthropic.String("after_id"),
		BeforeID: anthropic.String("before_id"),
		Limit:    anthropic.Int(1),
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestMessageBatchDelete(t *testing.T) {
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
	_, err := client.Messages.Batches.Delete(context.TODO(), "message_batch_id")
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestMessageBatchCancel(t *testing.T) {
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
	_, err := client.Messages.Batches.Cancel(context.TODO(), "message_batch_id")
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
