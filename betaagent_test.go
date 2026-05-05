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
			ID:    anthropic.BetaManagedAgentsModelClaudeOpus4_6,
			Speed: anthropic.BetaManagedAgentsModelConfigParamsSpeedStandard,
		},
		Name:        "My First Agent",
		Description: anthropic.String("A general-purpose starter agent."),
		MCPServers: []anthropic.BetaManagedAgentsURLMCPServerParams{{
			Name: "example-mcp",
			Type: anthropic.BetaManagedAgentsURLMCPServerParamsTypeURL,
			URL:  "https://example-server.modelcontextprotocol.io/sse",
		}},
		Metadata: map[string]string{
			"foo": "bar",
		},
		Multiagent: anthropic.BetaManagedAgentsMultiagentParams{
			Agents: []anthropic.BetaManagedAgentsMultiagentRosterEntryParamsUnion{{
				OfString: anthropic.String("agent_011CZkYqphY8vELVzwCUpqiQ"),
			}, {
				OfBetaManagedAgentsMultiagentSelfs: &anthropic.BetaManagedAgentsMultiagentSelfParams{
					Type: anthropic.BetaManagedAgentsMultiagentSelfParamsTypeSelf,
				},
			}},
			Type: anthropic.BetaManagedAgentsMultiagentParamsTypeCoordinator,
		},
		Skills: []anthropic.BetaManagedAgentsSkillParamsUnion{{
			OfAnthropic: &anthropic.BetaManagedAgentsAnthropicSkillParams{
				SkillID: "xlsx",
				Type:    anthropic.BetaManagedAgentsAnthropicSkillParamsTypeAnthropic,
				Version: anthropic.String("1"),
			},
		}},
		System: anthropic.String("You are a general-purpose agent that can research, write code, run commands, and use connected tools to complete the user's task end to end."),
		Tools: []anthropic.BetaAgentNewParamsToolUnion{{
			OfAgentToolset20260401: &anthropic.BetaManagedAgentsAgentToolset20260401Params{
				Type: anthropic.BetaManagedAgentsAgentToolset20260401ParamsTypeAgentToolset20260401,
				Configs: []anthropic.BetaManagedAgentsAgentToolConfigParams{{
					Name:    anthropic.BetaManagedAgentsAgentToolConfigParamsNameBash,
					Enabled: anthropic.Bool(true),
					PermissionPolicy: anthropic.BetaManagedAgentsAgentToolConfigParamsPermissionPolicyUnion{
						OfAlwaysAllow: &anthropic.BetaManagedAgentsAlwaysAllowPolicyParam{
							Type: anthropic.BetaManagedAgentsAlwaysAllowPolicyTypeAlwaysAllow,
						},
					},
				}},
				DefaultConfig: anthropic.BetaManagedAgentsAgentToolsetDefaultConfigParams{
					Enabled: anthropic.Bool(true),
					PermissionPolicy: anthropic.BetaManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion{
						OfAlwaysAllow: &anthropic.BetaManagedAgentsAlwaysAllowPolicyParam{
							Type: anthropic.BetaManagedAgentsAlwaysAllowPolicyTypeAlwaysAllow,
						},
					},
				},
			},
		}},
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

func TestBetaAgentGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.Get(
		context.TODO(),
		"agent_011CZkYpogX7uDKUyvBTophP",
		anthropic.BetaAgentGetParams{
			Version: anthropic.Int(0),
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
			Version:     1,
			Description: anthropic.String("description"),
			MCPServers: []anthropic.BetaManagedAgentsURLMCPServerParams{{
				Name: "example-mcp",
				Type: anthropic.BetaManagedAgentsURLMCPServerParamsTypeURL,
				URL:  "https://example-server.modelcontextprotocol.io/sse",
			}},
			Metadata: map[string]string{
				"foo": "string",
			},
			Model: anthropic.BetaManagedAgentsModelConfigParams{
				ID:    anthropic.BetaManagedAgentsModelClaudeOpus4_6,
				Speed: anthropic.BetaManagedAgentsModelConfigParamsSpeedStandard,
			},
			Multiagent: anthropic.BetaManagedAgentsMultiagentParams{
				Agents: []anthropic.BetaManagedAgentsMultiagentRosterEntryParamsUnion{{
					OfString: anthropic.String("agent_011CZkYqphY8vELVzwCUpqiQ"),
				}, {
					OfBetaManagedAgentsMultiagentSelfs: &anthropic.BetaManagedAgentsMultiagentSelfParams{
						Type: anthropic.BetaManagedAgentsMultiagentSelfParamsTypeSelf,
					},
				}},
				Type: anthropic.BetaManagedAgentsMultiagentParamsTypeCoordinator,
			},
			Name: anthropic.String("name"),
			Skills: []anthropic.BetaManagedAgentsSkillParamsUnion{{
				OfAnthropic: &anthropic.BetaManagedAgentsAnthropicSkillParams{
					SkillID: "xlsx",
					Type:    anthropic.BetaManagedAgentsAnthropicSkillParamsTypeAnthropic,
					Version: anthropic.String("1"),
				},
			}},
			System: anthropic.String("You are a general-purpose agent that can research, write code, run commands, and use connected tools to complete the user's task end to end."),
			Tools: []anthropic.BetaAgentUpdateParamsToolUnion{{
				OfAgentToolset20260401: &anthropic.BetaManagedAgentsAgentToolset20260401Params{
					Type: anthropic.BetaManagedAgentsAgentToolset20260401ParamsTypeAgentToolset20260401,
					Configs: []anthropic.BetaManagedAgentsAgentToolConfigParams{{
						Name:    anthropic.BetaManagedAgentsAgentToolConfigParamsNameBash,
						Enabled: anthropic.Bool(true),
						PermissionPolicy: anthropic.BetaManagedAgentsAgentToolConfigParamsPermissionPolicyUnion{
							OfAlwaysAllow: &anthropic.BetaManagedAgentsAlwaysAllowPolicyParam{
								Type: anthropic.BetaManagedAgentsAlwaysAllowPolicyTypeAlwaysAllow,
							},
						},
					}},
					DefaultConfig: anthropic.BetaManagedAgentsAgentToolsetDefaultConfigParams{
						Enabled: anthropic.Bool(true),
						PermissionPolicy: anthropic.BetaManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion{
							OfAlwaysAllow: &anthropic.BetaManagedAgentsAlwaysAllowPolicyParam{
								Type: anthropic.BetaManagedAgentsAlwaysAllowPolicyTypeAlwaysAllow,
							},
						},
					},
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

func TestBetaAgentListWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Agents.List(context.TODO(), anthropic.BetaAgentListParams{
		CreatedAtGte:    anthropic.Time(time.Now()),
		CreatedAtLte:    anthropic.Time(time.Now()),
		IncludeArchived: anthropic.Bool(true),
		Limit:           anthropic.Int(0),
		Page:            anthropic.String("page"),
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
