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

func TestBetaDeploymentNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Deployments.New(context.TODO(), anthropic.BetaDeploymentNewParams{
		Agent: anthropic.BetaDeploymentNewParamsAgentUnion{
			OfString: anthropic.String("string"),
		},
		EnvironmentID: "x",
		InitialEvents: []anthropic.BetaManagedAgentsDeploymentInitialEventParamsUnion{{
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
		Name:        "x",
		Description: anthropic.String("description"),
		Metadata: map[string]string{
			"foo": "string",
		},
		Resources: []anthropic.BetaDeploymentNewParamsResourceUnion{{
			OfFile: &anthropic.BetaManagedAgentsFileResourceParams{
				FileID:    "file_011CNha8iCJcU1wXNR6q4V8w",
				Type:      anthropic.BetaManagedAgentsFileResourceParamsTypeFile,
				MountPath: anthropic.String("/uploads/receipt.pdf"),
			},
		}},
		Schedule: anthropic.BetaManagedAgentsScheduleParams{
			Expression: "x",
			Timezone:   "x",
			Type:       anthropic.BetaManagedAgentsScheduleParamsTypeCron,
		},
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

func TestBetaDeploymentGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Deployments.Get(
		context.TODO(),
		"deployment_id",
		anthropic.BetaDeploymentGetParams{
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

func TestBetaDeploymentUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Deployments.Update(
		context.TODO(),
		"deployment_id",
		anthropic.BetaDeploymentUpdateParams{
			Agent: anthropic.BetaDeploymentUpdateParamsAgentUnion{
				OfString: anthropic.String("string"),
			},
			Description:   anthropic.String("description"),
			EnvironmentID: anthropic.String("environment_id"),
			InitialEvents: []anthropic.BetaManagedAgentsDeploymentInitialEventParamsUnion{{
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
			Metadata: map[string]string{
				"foo": "string",
			},
			Name: anthropic.String("name"),
			Resources: []anthropic.BetaDeploymentUpdateParamsResourceUnion{{
				OfFile: &anthropic.BetaManagedAgentsFileResourceParams{
					FileID:    "file_011CNha8iCJcU1wXNR6q4V8w",
					Type:      anthropic.BetaManagedAgentsFileResourceParamsTypeFile,
					MountPath: anthropic.String("/uploads/receipt.pdf"),
				},
			}},
			Schedule: anthropic.BetaManagedAgentsScheduleParams{
				Expression: "x",
				Timezone:   "x",
				Type:       anthropic.BetaManagedAgentsScheduleParamsTypeCron,
			},
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

func TestBetaDeploymentListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Deployments.List(context.TODO(), anthropic.BetaDeploymentListParams{
		AgentID:         anthropic.String("agent_id"),
		CreatedAtGte:    anthropic.Time(time.Now()),
		CreatedAtLte:    anthropic.Time(time.Now()),
		IncludeArchived: anthropic.Bool(true),
		Limit:           anthropic.Int(0),
		Page:            anthropic.String("page"),
		Status:          anthropic.BetaManagedAgentsDeploymentStatusActive,
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

func TestBetaDeploymentArchiveWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Deployments.Archive(
		context.TODO(),
		"deployment_id",
		anthropic.BetaDeploymentArchiveParams{
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

func TestBetaDeploymentPauseWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Deployments.Pause(
		context.TODO(),
		"deployment_id",
		anthropic.BetaDeploymentPauseParams{
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

func TestBetaDeploymentRunWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Deployments.Run(
		context.TODO(),
		"deployment_id",
		anthropic.BetaDeploymentRunParams{
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

func TestBetaDeploymentUnpauseWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Deployments.Unpause(
		context.TODO(),
		"deployment_id",
		anthropic.BetaDeploymentUnpauseParams{
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
