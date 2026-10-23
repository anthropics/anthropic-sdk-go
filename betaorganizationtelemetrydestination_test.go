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

func TestBetaOrganizationTelemetryDestinationNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.TelemetryDestinations.New(context.TODO(), anthropic.BetaOrganizationTelemetryDestinationNewParams{
		EndpointURL: "x",
		Workspaces: anthropic.BetaTelemetryDestinationWorkspacesParamsUnion{
			OfAll: &anthropic.BetaTelemetryDestinationWorkspacesAllParams{},
		},
		DisplayName: anthropic.String("display_name"),
		Headers: map[string]string{
			"foo": "string",
		},
		Traces: anthropic.BetaTelemetryDestinationTracesParams{
			ContentMode:     anthropic.BetaTelemetryDestinationContentModeMetadataOnly,
			Enabled:         anthropic.Bool(true),
			EndpointURL:     anthropic.String("endpoint_url"),
			MaxContentBytes: anthropic.Int(0),
		},
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

func TestBetaOrganizationTelemetryDestinationGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.TelemetryDestinations.Get(
		context.TODO(),
		"destination_id",
		anthropic.BetaOrganizationTelemetryDestinationGetParams{
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

func TestBetaOrganizationTelemetryDestinationUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.TelemetryDestinations.Update(
		context.TODO(),
		"destination_id",
		anthropic.BetaOrganizationTelemetryDestinationUpdateParams{
			DisplayName: anthropic.String("display_name"),
			EndpointURL: anthropic.String("endpoint_url"),
			Headers: map[string]string{
				"foo": "string",
			},
			State: anthropic.BetaTelemetryDestinationStateParamsUnion{
				OfEnabled: &anthropic.BetaTelemetryDestinationStateEnabledParams{},
			},
			Traces: anthropic.BetaTelemetryDestinationTracesParams{
				ContentMode:     anthropic.BetaTelemetryDestinationContentModeMetadataOnly,
				Enabled:         anthropic.Bool(true),
				EndpointURL:     anthropic.String("endpoint_url"),
				MaxContentBytes: anthropic.Int(0),
			},
			Workspaces: anthropic.BetaTelemetryDestinationWorkspacesParamsUnion{
				OfAll: &anthropic.BetaTelemetryDestinationWorkspacesAllParams{},
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

func TestBetaOrganizationTelemetryDestinationListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.TelemetryDestinations.List(context.TODO(), anthropic.BetaOrganizationTelemetryDestinationListParams{
		Limit:        anthropic.Int(0),
		Page:         anthropic.String("page"),
		WorkspaceIDs: []string{"string"},
		Betas:        []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaOrganizationTelemetryDestinationDeleteWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.TelemetryDestinations.Delete(
		context.TODO(),
		"destination_id",
		anthropic.BetaOrganizationTelemetryDestinationDeleteParams{
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

func TestBetaOrganizationTelemetryDestinationTestWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Organization.TelemetryDestinations.Test(
		context.TODO(),
		"destination_id",
		anthropic.BetaOrganizationTelemetryDestinationTestParams{
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
