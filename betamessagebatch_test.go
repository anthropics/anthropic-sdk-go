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
		Requests: anthropic.F([]anthropic.BetaMessageBatchNewParamsRequest{{
			CustomID: anthropic.F("my-custom-id-1"),
			Params: anthropic.F(anthropic.BetaMessageBatchNewParamsRequestsParams{
				MaxTokens: anthropic.F(int64(1024)),
				Messages: anthropic.F([]anthropic.BetaMessageParam{{
					Content: anthropic.F([]anthropic.BetaContentBlockParamUnion{anthropic.BetaTextBlockParam{Text: anthropic.F("What is a quaternion?"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)}), Citations: anthropic.F([]anthropic.BetaTextCitationParamUnion{anthropic.BetaCitationCharLocationParam{CitedText: anthropic.F("cited_text"), DocumentIndex: anthropic.F(int64(0)), DocumentTitle: anthropic.F("x"), EndCharIndex: anthropic.F(int64(0)), StartCharIndex: anthropic.F(int64(0)), Type: anthropic.F(anthropic.BetaCitationCharLocationParamTypeCharLocation)}})}}),
					Role:    anthropic.F(anthropic.BetaMessageParamRoleUser),
				}}),
				Model: anthropic.F(anthropic.ModelClaude3_7SonnetLatest),
				Metadata: anthropic.F(anthropic.BetaMetadataParam{
					UserID: anthropic.F("13803d75-b4b5-4c3e-b2a2-6f21399b021b"),
				}),
				StopSequences: anthropic.F([]string{"string"}),
				Stream:        anthropic.F(true),
				System:        anthropic.F([]anthropic.BetaTextBlockParam{{Text: anthropic.F("x"), Type: anthropic.F(anthropic.BetaTextBlockParamTypeText), CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral)}), Citations: anthropic.F([]anthropic.BetaTextCitationParamUnion{anthropic.BetaCitationCharLocationParam{CitedText: anthropic.F("cited_text"), DocumentIndex: anthropic.F(int64(0)), DocumentTitle: anthropic.F("x"), EndCharIndex: anthropic.F(int64(0)), StartCharIndex: anthropic.F(int64(0)), Type: anthropic.F(anthropic.BetaCitationCharLocationParamTypeCharLocation)}})}}),
				Temperature:   anthropic.F(1.000000),
				Thinking: anthropic.F[anthropic.BetaThinkingConfigParamUnion](anthropic.BetaThinkingConfigEnabledParam{
					BudgetTokens: anthropic.F(int64(1024)),
					Type:         anthropic.F(anthropic.BetaThinkingConfigEnabledTypeEnabled),
				}),
				ToolChoice: anthropic.F[anthropic.BetaToolChoiceUnionParam](anthropic.BetaToolChoiceAutoParam{
					Type:                   anthropic.F(anthropic.BetaToolChoiceAutoTypeAuto),
					DisableParallelToolUse: anthropic.F(true),
				}),
				Tools: anthropic.F([]anthropic.BetaToolUnionUnionParam{anthropic.BetaToolComputerUse20241022Param{
					DisplayHeightPx: anthropic.F(int64(1)),
					DisplayWidthPx:  anthropic.F(int64(1)),
					Name:            anthropic.F(anthropic.BetaToolComputerUse20241022NameComputer),
					Type:            anthropic.F(anthropic.BetaToolComputerUse20241022TypeComputer20241022),
					CacheControl: anthropic.F(anthropic.BetaCacheControlEphemeralParam{
						Type: anthropic.F(anthropic.BetaCacheControlEphemeralTypeEphemeral),
					}),
					DisplayNumber: anthropic.F(int64(0)),
				}}),
				TopK: anthropic.F(int64(5)),
				TopP: anthropic.F(0.700000),
			}),
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
			Betas: anthropic.F([]anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24}),
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
		Betas:    anthropic.F([]anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24}),
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
			Betas: anthropic.F([]anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24}),
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
			Betas: anthropic.F([]anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24}),
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
