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

func TestBetaSessionResourceGetWithOptionalParams(t *testing.T) {
	t.Skip("prism can't find endpoint with beta only tag")
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
	_, err := client.Beta.Sessions.Resources.Get(
		context.TODO(),
		"sesrsc_011CZkZBJq5dWxk9fVLNcPht",
		anthropic.BetaSessionResourceGetParams{
			SessionID: "sesn_011CZkZAtmR3yMPDzynEDxu7",
			Betas:     []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
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

func TestBetaSessionResourceUpdateWithOptionalParams(t *testing.T) {
	t.Skip("prism can't find endpoint with beta only tag")
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
	_, err := client.Beta.Sessions.Resources.Update(
		context.TODO(),
		"sesrsc_011CZkZBJq5dWxk9fVLNcPht",
		anthropic.BetaSessionResourceUpdateParams{
			SessionID:          "sesn_011CZkZAtmR3yMPDzynEDxu7",
			AuthorizationToken: "ghp_exampletoken",
			Betas:              []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
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

func TestBetaSessionResourceListWithOptionalParams(t *testing.T) {
	t.Skip("prism can't find endpoint with beta only tag")
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
	_, err := client.Beta.Sessions.Resources.List(
		context.TODO(),
		"sesn_011CZkZAtmR3yMPDzynEDxu7",
		anthropic.BetaSessionResourceListParams{
			Limit: anthropic.Int(0),
			Page:  anthropic.String("page"),
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

func TestBetaSessionResourceDeleteWithOptionalParams(t *testing.T) {
	t.Skip("prism can't find endpoint with beta only tag")
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
	_, err := client.Beta.Sessions.Resources.Delete(
		context.TODO(),
		"sesrsc_011CZkZBJq5dWxk9fVLNcPht",
		anthropic.BetaSessionResourceDeleteParams{
			SessionID: "sesn_011CZkZAtmR3yMPDzynEDxu7",
			Betas:     []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
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

func TestBetaSessionResourceAddWithOptionalParams(t *testing.T) {
	t.Skip("prism can't find endpoint with beta only tag")
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
	_, err := client.Beta.Sessions.Resources.Add(
		context.TODO(),
		"sesn_011CZkZAtmR3yMPDzynEDxu7",
		anthropic.BetaSessionResourceAddParams{
			BetaManagedAgentsFileResourceParams: anthropic.BetaManagedAgentsFileResourceParams{
				FileID:    "file_011CNha8iCJcU1wXNR6q4V8w",
				Type:      anthropic.BetaManagedAgentsFileResourceParamsTypeFile,
				MountPath: anthropic.String("/uploads/receipt.pdf"),
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
