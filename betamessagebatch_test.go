// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
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
		Requests: anthropic.F([]anthropic.BetaMessageBatchNewParamsRequest{{
			CustomID: anthropic.F("my-custom-id-1"),
			Params: anthropic.F(anthropic.BetaMessageBatchNewParamsRequestsParams{
				MaxTokens: anthropic.F(int64(1024)),
				Messages: anthropic.F([]anthropic.BetaMessageParam{{
					Content: anthropic.F([]anthropic.BetaContentBlockParamUnion{anthropic.BetaTextBlockParam{Text: anthropic.F("What is a quaternion?"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)})}}),
					Role:    anthropic.F(anthropic.BetaMessageParamRoleUser),
				}}),
				Model: anthropic.F(anthropic.ModelClaude3_5HaikuLatest),
				Metadata: anthropic.F(anthropic.BetaMetadataParam{
					UserID: anthropic.F("13803d75-b4b5-4c3e-b2a2-6f21399b021b"),
				}),
				StopSequences: anthropic.F([]string{"string", "string", "string"}),
				Stream:        anthropic.F(true),
				System:        anthropic.F([]anthropic.BetaTextBlockParam{{Text: anthropic.F("x"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)})}, {Text: anthropic.F("x"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)})}, {Text: anthropic.F("x"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)})}}),
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
					Name: anthropic.F("x"),
					CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{
						Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral),
					}),
					Description: anthropic.F("Get the current weather in a given location"),
					Type:        anthropic.F(anthropic.BetaToolTypeCustom),
				}, anthropic.BetaToolParam{
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
					Name: anthropic.F("x"),
					CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{
						Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral),
					}),
					Description: anthropic.F("Get the current weather in a given location"),
					Type:        anthropic.F(anthropic.BetaToolTypeCustom),
				}, anthropic.BetaToolParam{
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
					Name: anthropic.F("x"),
					CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{
						Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral),
					}),
					Description: anthropic.F("Get the current weather in a given location"),
					Type:        anthropic.F(anthropic.BetaToolTypeCustom),
				}}),
				TopK: anthropic.F(int64(5)),
				TopP: anthropic.F(0.700000),
			}),
		}, {
			CustomID: anthropic.F("my-custom-id-1"),
			Params: anthropic.F(anthropic.BetaMessageBatchNewParamsRequestsParams{
				MaxTokens: anthropic.F(int64(1024)),
				Messages: anthropic.F([]anthropic.BetaMessageParam{{
					Content: anthropic.F([]anthropic.BetaContentBlockParamUnion{anthropic.BetaTextBlockParam{Text: anthropic.F("What is a quaternion?"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)})}}),
					Role:    anthropic.F(anthropic.BetaMessageParamRoleUser),
				}}),
				Model: anthropic.F(anthropic.ModelClaude3_5HaikuLatest),
				Metadata: anthropic.F(anthropic.BetaMetadataParam{
					UserID: anthropic.F("13803d75-b4b5-4c3e-b2a2-6f21399b021b"),
				}),
				StopSequences: anthropic.F([]string{"string", "string", "string"}),
				Stream:        anthropic.F(true),
				System:        anthropic.F([]anthropic.BetaTextBlockParam{{Text: anthropic.F("x"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)})}, {Text: anthropic.F("x"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)})}, {Text: anthropic.F("x"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)})}}),
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
					Name: anthropic.F("x"),
					CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{
						Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral),
					}),
					Description: anthropic.F("Get the current weather in a given location"),
					Type:        anthropic.F(anthropic.BetaToolTypeCustom),
				}, anthropic.BetaToolParam{
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
					Name: anthropic.F("x"),
					CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{
						Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral),
					}),
					Description: anthropic.F("Get the current weather in a given location"),
					Type:        anthropic.F(anthropic.BetaToolTypeCustom),
				}, anthropic.BetaToolParam{
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
					Name: anthropic.F("x"),
					CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{
						Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral),
					}),
					Description: anthropic.F("Get the current weather in a given location"),
					Type:        anthropic.F(anthropic.BetaToolTypeCustom),
				}}),
				TopK: anthropic.F(int64(5)),
				TopP: anthropic.F(0.700000),
			}),
		}, {
			CustomID: anthropic.F("my-custom-id-1"),
			Params: anthropic.F(anthropic.BetaMessageBatchNewParamsRequestsParams{
				MaxTokens: anthropic.F(int64(1024)),
				Messages: anthropic.F([]anthropic.BetaMessageParam{{
					Content: anthropic.F([]anthropic.BetaContentBlockParamUnion{anthropic.BetaTextBlockParam{Text: anthropic.F("What is a quaternion?"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)})}}),
					Role:    anthropic.F(anthropic.BetaMessageParamRoleUser),
				}}),
				Model: anthropic.F(anthropic.ModelClaude3_5HaikuLatest),
				Metadata: anthropic.F(anthropic.BetaMetadataParam{
					UserID: anthropic.F("13803d75-b4b5-4c3e-b2a2-6f21399b021b"),
				}),
				StopSequences: anthropic.F([]string{"string", "string", "string"}),
				Stream:        anthropic.F(true),
				System:        anthropic.F([]anthropic.BetaTextBlockParam{{Text: anthropic.F("x"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)})}, {Text: anthropic.F("x"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)})}, {Text: anthropic.F("x"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)})}}),
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
					Name: anthropic.F("x"),
					CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{
						Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral),
					}),
					Description: anthropic.F("Get the current weather in a given location"),
					Type:        anthropic.F(anthropic.BetaToolTypeCustom),
				}, anthropic.BetaToolParam{
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
					Name: anthropic.F("x"),
					CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{
						Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral),
					}),
					Description: anthropic.F("Get the current weather in a given location"),
					Type:        anthropic.F(anthropic.BetaToolTypeCustom),
				}, anthropic.BetaToolParam{
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
					Name: anthropic.F("x"),
					CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{
						Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral),
					}),
					Description: anthropic.F("Get the current weather in a given location"),
					Type:        anthropic.F(anthropic.BetaToolTypeCustom),
				}}),
				TopK: anthropic.F(int64(5)),
				TopP: anthropic.F(0.700000),
			}),
		}}),
		Betas: anthropic.F([]anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24, anthropic.AnthropicBetaMessageBatches2024_09_24, anthropic.AnthropicBetaMessageBatches2024_09_24}),
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
			Betas: anthropic.F([]anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24, anthropic.AnthropicBetaMessageBatches2024_09_24, anthropic.AnthropicBetaMessageBatches2024_09_24}),
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
		AfterID:  anthropic.F("after_id"),
		BeforeID: anthropic.F("before_id"),
		Limit:    anthropic.F(int64(1)),
		Betas:    anthropic.F([]anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24, anthropic.AnthropicBetaMessageBatches2024_09_24, anthropic.AnthropicBetaMessageBatches2024_09_24}),
	})
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
			Betas: anthropic.F([]anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24, anthropic.AnthropicBetaMessageBatches2024_09_24, anthropic.AnthropicBetaMessageBatches2024_09_24}),
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

func TestBetaMessageBatchResultsWithOptionalParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("abc"))
	}))
	defer server.Close()
	baseURL := server.URL
	client := anthropic.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("my-anthropic-api-key"),
	)
	resp, err := client.Beta.Messages.Batches.Results(
		context.TODO(),
		"message_batch_id",
		anthropic.BetaMessageBatchResultsParams{
			Betas: anthropic.F([]anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24, anthropic.AnthropicBetaMessageBatches2024_09_24, anthropic.AnthropicBetaMessageBatches2024_09_24}),
		},
	)
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
	if !bytes.Equal(b, []byte("abc")) {
		t.Fatalf("return value not %s: %s", "abc", b)
	}
}
