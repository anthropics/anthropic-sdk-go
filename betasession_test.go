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

func TestBetaSessionNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Sessions.New(context.TODO(), anthropic.BetaSessionNewParams{
		Agent: anthropic.BetaSessionNewParamsAgentUnion{
			OfString: anthropic.String("agent_011CZkYpogX7uDKUyvBTophP"),
		},
		EnvironmentID: "env_011CZkZ9X2dpNyB7HsEFoRfW",
		Metadata: map[string]string{
			"foo": "string",
		},
		Resources: []anthropic.BetaSessionNewParamsResourceUnion{{
			OfFile: &anthropic.BetaManagedAgentsFileResourceParams{
				FileID:    "file_011CNha8iCJcU1wXNR6q4V8w",
				Type:      anthropic.BetaManagedAgentsFileResourceParamsTypeFile,
				MountPath: anthropic.String("/uploads/receipt.pdf"),
			},
		}},
		Title:    anthropic.String("Order #1234 inquiry"),
		VaultIDs: []string{"string"},
		Betas:    []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaSessionGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Sessions.Get(
		context.TODO(),
		"sesn_011CZkZAtmR3yMPDzynEDxu7",
		anthropic.BetaSessionGetParams{
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

func TestBetaSessionUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Sessions.Update(
		context.TODO(),
		"sesn_011CZkZAtmR3yMPDzynEDxu7",
		anthropic.BetaSessionUpdateParams{
			Metadata: map[string]string{
				"foo": "string",
			},
			Title:    anthropic.String("Order #1234 inquiry"),
			VaultIDs: []string{"string"},
			Betas:    []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
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

func TestBetaSessionListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Sessions.List(context.TODO(), anthropic.BetaSessionListParams{
		AgentID:         anthropic.String("agent_id"),
		AgentVersion:    anthropic.Int(0),
		CreatedAtGt:     anthropic.Time(time.Now()),
		CreatedAtGte:    anthropic.Time(time.Now()),
		CreatedAtLt:     anthropic.Time(time.Now()),
		CreatedAtLte:    anthropic.Time(time.Now()),
		IncludeArchived: anthropic.Bool(true),
		Limit:           anthropic.Int(0),
		MemoryStoreID:   anthropic.String("memory_store_id"),
		Order:           anthropic.BetaSessionListParamsOrderAsc,
		Page:            anthropic.String("page"),
		Statuses:        []string{"rescheduling"},
		Betas:           []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaSessionDeleteWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Sessions.Delete(
		context.TODO(),
		"sesn_011CZkZAtmR3yMPDzynEDxu7",
		anthropic.BetaSessionDeleteParams{
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

func TestBetaSessionArchiveWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Sessions.Archive(
		context.TODO(),
		"sesn_011CZkZAtmR3yMPDzynEDxu7",
		anthropic.BetaSessionArchiveParams{
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
