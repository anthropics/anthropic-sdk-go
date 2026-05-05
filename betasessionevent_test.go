// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestBetaSessionEventListWithOptionalParams(t *testing.T) {
	t.Skip("buildURL drops path-level query params (SDK-4349)")
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
	_, err := client.Beta.Sessions.Events.List(
		context.TODO(),
		"sesn_011CZkZAtmR3yMPDzynEDxu7",
		anthropic.BetaSessionEventListParams{
			CreatedAtGt:  anthropic.Time(time.Now()),
			CreatedAtGte: anthropic.Time(time.Now()),
			CreatedAtLt:  anthropic.Time(time.Now()),
			CreatedAtLte: anthropic.Time(time.Now()),
			Limit:        anthropic.Int(0),
			Order:        anthropic.BetaSessionEventListParamsOrderAsc,
			Page:         anthropic.String("page"),
			Types:        []string{"string"},
			Betas:        []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
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

func TestBetaSessionEventSendWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Sessions.Events.Send(
		context.TODO(),
		"sesn_011CZkZAtmR3yMPDzynEDxu7",
		anthropic.BetaSessionEventSendParams{
			Events: []anthropic.BetaManagedAgentsEventParamsUnion{{
				OfUserMessage: &anthropic.BetaManagedAgentsUserMessageEventParams{
					Content: []anthropic.BetaManagedAgentsUserMessageEventParamsContentUnion{{
						OfText: &anthropic.BetaManagedAgentsTextBlockParam{
							Text: "Where is my order #1234?",
							Type: anthropic.BetaManagedAgentsTextBlockTypeText,
						},
					}},
					Type: anthropic.BetaManagedAgentsUserMessageEventParamsTypeUserMessage,
				},
			}},
			Betas: []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
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
