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

func TestBetaMemoryStoreMemoryNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.MemoryStores.Memories.New(
		context.TODO(),
		"memory_store_id",
		anthropic.BetaMemoryStoreMemoryNewParams{
			Content: anthropic.String("content"),
			Path:    "xx",
			View:    anthropic.BetaManagedAgentsMemoryViewBasic,
			Betas:   []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
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

func TestBetaMemoryStoreMemoryGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.MemoryStores.Memories.Get(
		context.TODO(),
		"memory_id",
		anthropic.BetaMemoryStoreMemoryGetParams{
			MemoryStoreID: "memory_store_id",
			View:          anthropic.BetaManagedAgentsMemoryViewBasic,
			Betas:         []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
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

func TestBetaMemoryStoreMemoryUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.MemoryStores.Memories.Update(
		context.TODO(),
		"memory_id",
		anthropic.BetaMemoryStoreMemoryUpdateParams{
			MemoryStoreID: "memory_store_id",
			View:          anthropic.BetaManagedAgentsMemoryViewBasic,
			Content:       anthropic.String("content"),
			Path:          anthropic.String("xx"),
			Precondition: anthropic.BetaManagedAgentsPreconditionParam{
				Type:          anthropic.BetaManagedAgentsPreconditionTypeContentSha256,
				ContentSha256: anthropic.String("content_sha256"),
			},
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

func TestBetaMemoryStoreMemoryListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.MemoryStores.Memories.List(
		context.TODO(),
		"memory_store_id",
		anthropic.BetaMemoryStoreMemoryListParams{
			Depth:      anthropic.Int(0),
			Limit:      anthropic.Int(0),
			Order:      anthropic.BetaMemoryStoreMemoryListParamsOrderAsc,
			OrderBy:    anthropic.String("order_by"),
			Page:       anthropic.String("page"),
			PathPrefix: anthropic.String("path_prefix"),
			View:       anthropic.BetaManagedAgentsMemoryViewBasic,
			Betas:      []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
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

func TestBetaMemoryStoreMemoryDeleteWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.MemoryStores.Memories.Delete(
		context.TODO(),
		"memory_id",
		anthropic.BetaMemoryStoreMemoryDeleteParams{
			MemoryStoreID:         "memory_store_id",
			ExpectedContentSha256: anthropic.String("expected_content_sha256"),
			Betas:                 []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
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
