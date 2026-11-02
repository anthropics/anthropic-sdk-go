package anthropic_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

func TestBetaSessionThreadGetWithOptionalParams(t *testing.T) {
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
		option.WithQuery("beta", "true"),
	)
	_, err := client.Beta.Sessions.Threads.Get(
		context.TODO(),
		"sthr_011CZkZVWa6oJjw1rgXZpnBt",
		anthropic.BetaSessionThreadGetParams{
			SessionID:   "sesn_011CZkZAtmR3yMPDzynEDxu7",
			Betas:       []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
			WorkspaceID: anthropic.String("wrkspc_011CZkZaBF1tNoB5wlCeusgy"),
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

func TestBetaSessionThreadListWithOptionalParams(t *testing.T) {
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
		option.WithQuery("beta", "true"),
	)
	_, err := client.Beta.Sessions.Threads.List(
		context.TODO(),
		"sesn_011CZkZAtmR3yMPDzynEDxu7",
		anthropic.BetaSessionThreadListParams{
			Limit:       anthropic.Int(0),
			Page:        anthropic.String("page"),
			Statuses:    []anthropic.BetaManagedAgentsSessionThreadStatus{anthropic.BetaManagedAgentsSessionThreadStatusRunning},
			Betas:       []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
			WorkspaceID: anthropic.String("wrkspc_011CZkZaBF1tNoB5wlCeusgy"),
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

func TestBetaSessionThreadArchiveWithOptionalParams(t *testing.T) {
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
		option.WithQuery("beta", "true"),
	)
	_, err := client.Beta.Sessions.Threads.Archive(
		context.TODO(),
		"sthr_011CZkZVWa6oJjw1rgXZpnBt",
		anthropic.BetaSessionThreadArchiveParams{
			SessionID:   "sesn_011CZkZAtmR3yMPDzynEDxu7",
			Betas:       []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
			WorkspaceID: anthropic.String("wrkspc_011CZkZaBF1tNoB5wlCeusgy"),
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
