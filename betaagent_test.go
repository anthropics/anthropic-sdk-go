package anthropic_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

func TestBetaAgentNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.New(context.TODO(), anthropic.BetaAgentNewParams{
		Model: anthropic.BetaManagedAgentsModelConfigParams{
			ID: anthropic.BetaManagedAgentsModelClaudeOpus5_5,
			Effort: anthropic.BetaManagedAgentsModelConfigParamsEffortUnion{
				OfBetaManagedAgentsModelConfigsEffortBetaManagedAgentsEffortLevel: anthropic.String("low"),
			},
			InferenceGeo: anthropic.String("inference_geo"),
			Speed:        anthropic.BetaManagedAgentsModelConfigParamsSpeedStandard,
		},
		Name:        "My First Agent",
		Description: anthropic.String("A general-purpose starter agent."),
		ExecutionIdentity: anthropic.BetaManagedAgentsExecutionIdentityParamsUnion{
			OfAWSRole: &anthropic.BetaManagedAgentsAWSRoleExecutionIdentityParams{
				RoleARN: "role_arn",
			},
		},
		MCPServers: []anthropic.BetaManagedAgentsURLMCPServerParams{{
			Name: "example-mcp",
			URL:  "https://example-server.modelcontextprotocol.io/sse",
		}},
		Metadata: map[string]string{
			"foo": "bar",
		},
		Multiagent: anthropic.BetaManagedAgentsMultiagentParamsUnion{
			OfCoordinator: &anthropic.BetaManagedAgentsMultiagentCoordinatorParams{
				Agents: []anthropic.BetaManagedAgentsMultiagentRosterEntryParamsUnion{{
					OfString: anthropic.String("agent_011CZkYqphY8vELVzwCUpqiQ"),
				}, {
					OfBetaManagedAgentsMultiagentSelfs: &anthropic.BetaManagedAgentsMultiagentSelfParams{},
				}},
			},
		},
		Skills: []anthropic.BetaManagedAgentsSkillParamsUnion{{
			OfAnthropic: &anthropic.BetaManagedAgentsAnthropicSkillParams{
				SkillID: "xlsx",
				Version: anthropic.String("1"),
			},
		}},
		System: anthropic.String("You are a general-purpose agent that can research, write code, run commands, and use connected tools to complete the user's task end to end."),
		Tools: []anthropic.BetaAgentNewParamsToolUnion{{
			OfAgentToolset20260401: &anthropic.BetaManagedAgentsAgentToolset20260401Params{
				Configs: []anthropic.BetaManagedAgentsAgentToolConfigParamsUnion{{
					OfBash: &anthropic.BetaManagedAgentsBashToolConfigParams{
						Enabled: anthropic.Bool(true),
						PermissionPolicy: anthropic.BetaManagedAgentsBashToolConfigParamsPermissionPolicyUnion{
							OfAlwaysAllow: &anthropic.BetaManagedAgentsAlwaysAllowPolicyParam{},
						},
						Type: anthropic.BetaManagedAgentsBashToolConfigParamsTypeBash,
					},
				}},
				DefaultConfig: anthropic.BetaManagedAgentsAgentToolsetDefaultConfigParams{
					Enabled: anthropic.Bool(true),
					PermissionPolicy: anthropic.BetaManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion{
						OfAlwaysAllow: &anthropic.BetaManagedAgentsAlwaysAllowPolicyParam{},
					},
				},
			},
		}},
		Betas:       []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
		WorkspaceID: anthropic.String("wrkspc_011CZkZaBF1tNoB5wlCeusgy"),
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.Get(
		context.TODO(),
		"agent_011CZkYpogX7uDKUyvBTophP",
		anthropic.BetaAgentGetParams{
			Version:     anthropic.Int(0),
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

func TestBetaAgentUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.Update(
		context.TODO(),
		"agent_011CZkYpogX7uDKUyvBTophP",
		anthropic.BetaAgentUpdateParams{
			Description: anthropic.String("updated"),
			ExecutionIdentity: anthropic.BetaManagedAgentsExecutionIdentityParamsUnion{
				OfAWSRole: &anthropic.BetaManagedAgentsAWSRoleExecutionIdentityParams{
					RoleARN: "role_arn",
				},
			},
			MCPServers: []anthropic.BetaManagedAgentsURLMCPServerParams{{
				Name: "example-mcp",
				URL:  "https://example-server.modelcontextprotocol.io/sse",
			}},
			Metadata: map[string]string{
				"foo": "string",
			},
			Model: anthropic.BetaManagedAgentsModelConfigParams{
				ID: anthropic.BetaManagedAgentsModelClaudeOpus5_5,
				Effort: anthropic.BetaManagedAgentsModelConfigParamsEffortUnion{
					OfBetaManagedAgentsModelConfigsEffortBetaManagedAgentsEffortLevel: anthropic.String("low"),
				},
				InferenceGeo: anthropic.String("inference_geo"),
				Speed:        anthropic.BetaManagedAgentsModelConfigParamsSpeedStandard,
			},
			Multiagent: anthropic.BetaManagedAgentsMultiagentParamsUnion{
				OfCoordinator: &anthropic.BetaManagedAgentsMultiagentCoordinatorParams{
					Agents: []anthropic.BetaManagedAgentsMultiagentRosterEntryParamsUnion{{
						OfString: anthropic.String("agent_011CZkYqphY8vELVzwCUpqiQ"),
					}, {
						OfBetaManagedAgentsMultiagentSelfs: &anthropic.BetaManagedAgentsMultiagentSelfParams{},
					}},
				},
			},
			Name: anthropic.String("name"),
			Skills: []anthropic.BetaManagedAgentsSkillParamsUnion{{
				OfAnthropic: &anthropic.BetaManagedAgentsAnthropicSkillParams{
					SkillID: "xlsx",
					Version: anthropic.String("1"),
				},
			}},
			System: anthropic.String("You are a general-purpose agent that can research, write code, run commands, and use connected tools to complete the user's task end to end."),
			Tools: []anthropic.BetaAgentUpdateParamsToolUnion{{
				OfAgentToolset20260401: &anthropic.BetaManagedAgentsAgentToolset20260401Params{
					Configs: []anthropic.BetaManagedAgentsAgentToolConfigParamsUnion{{
						OfBash: &anthropic.BetaManagedAgentsBashToolConfigParams{
							Enabled: anthropic.Bool(true),
							PermissionPolicy: anthropic.BetaManagedAgentsBashToolConfigParamsPermissionPolicyUnion{
								OfAlwaysAllow: &anthropic.BetaManagedAgentsAlwaysAllowPolicyParam{},
							},
							Type: anthropic.BetaManagedAgentsBashToolConfigParamsTypeBash,
						},
					}},
					DefaultConfig: anthropic.BetaManagedAgentsAgentToolsetDefaultConfigParams{
						Enabled: anthropic.Bool(true),
						PermissionPolicy: anthropic.BetaManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion{
							OfAlwaysAllow: &anthropic.BetaManagedAgentsAlwaysAllowPolicyParam{},
						},
					},
				},
			}},
			Version:     anthropic.Int(1),
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

func TestBetaAgentListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.List(context.TODO(), anthropic.BetaAgentListParams{
		CreatedAtGte:    anthropic.Time(time.Now()),
		CreatedAtLte:    anthropic.Time(time.Now()),
		IncludeArchived: anthropic.Bool(true),
		Limit:           anthropic.Int(0),
		Page:            anthropic.String("page"),
		Betas:           []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
		WorkspaceID:     anthropic.String("wrkspc_011CZkZaBF1tNoB5wlCeusgy"),
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAgentArchiveWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.Archive(
		context.TODO(),
		"agent_011CZkYpogX7uDKUyvBTophP",
		anthropic.BetaAgentArchiveParams{
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
