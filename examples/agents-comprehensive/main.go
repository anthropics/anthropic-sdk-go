package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

const (
	mcpServerName = "github"
	mcpServerURL  = "https://api.githubcopilot.com/mcp/"

	prompt = "Hi! List every tool and skill you have access to, grouped by where they " +
		"came from (built-in toolset, custom tool, MCP server, skills)."
)

func main() {
	client := anthropic.NewClient()
	ctx := context.TODO()

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		panic("GITHUB_TOKEN is required (use a fine-grained PAT with public-repo read only)")
	}

	// Create an environment
	environment, err := client.Beta.Environments.New(ctx, anthropic.BetaEnvironmentNewParams{
		Name: "comprehensive-example-environment",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created environment:", environment.ID)

	// Create a vault and store the MCP server credential in it
	vault, err := client.Beta.Vaults.New(ctx, anthropic.BetaVaultNewParams{
		DisplayName: "comprehensive-example-vault",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created vault:", vault.ID)

	credential, err := client.Beta.Vaults.Credentials.New(ctx, vault.ID, anthropic.BetaVaultCredentialNewParams{
		DisplayName: param.NewOpt("github-mcp"),
		Auth: anthropic.BetaVaultCredentialNewParamsAuthUnion{
			OfStaticBearer: &anthropic.BetaManagedAgentsStaticBearerCreateParams{
				Type:         anthropic.BetaManagedAgentsStaticBearerCreateParamsTypeStaticBearer,
				MCPServerURL: mcpServerURL,
				Token:        githubToken,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created credential:", credential.ID)

	// Upload a custom skill
	skillFile, err := os.Open("agents-comprehensive/greeting-SKILL.md")
	if err != nil {
		panic(err)
	}
	defer skillFile.Close()

	skill, err := client.Beta.Skills.New(ctx, anthropic.BetaSkillNewParams{
		DisplayTitle: param.NewOpt(fmt.Sprintf("comprehensive-greeting-%d", time.Now().UnixMilli())),
		Files:        []io.Reader{namedReader{skillFile, "greeting/SKILL.md"}},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created skill:", skill.ID)

	// Create v1 of the agent with the built-in toolset, an MCP server, and a custom tool
	agentV1, err := client.Beta.Agents.New(ctx, anthropic.BetaAgentNewParams{
		Name: "comprehensive-example-agent",
		Model: anthropic.BetaManagedAgentsModelConfigParams{
			ID: anthropic.BetaManagedAgentsModelClaudeSonnet4_6,
		},
		System: param.NewOpt("You are a helpful assistant."),
		MCPServers: []anthropic.BetaManagedAgentsURLMCPServerParams{
			{
				Type: anthropic.BetaManagedAgentsURLMCPServerParamsTypeURL,
				Name: mcpServerName,
				URL:  mcpServerURL,
			},
		},
		Tools: []anthropic.BetaAgentNewParamsToolUnion{
			{
				OfAgentToolset20260401: &anthropic.BetaManagedAgentsAgentToolset20260401Params{
					Type: anthropic.BetaManagedAgentsAgentToolset20260401ParamsTypeAgentToolset20260401,
				},
			},
			{
				OfMCPToolset: &anthropic.BetaManagedAgentsMCPToolsetParams{
					Type:          anthropic.BetaManagedAgentsMCPToolsetParamsTypeMCPToolset,
					MCPServerName: mcpServerName,
				},
			},
			{
				OfCustom: &anthropic.BetaManagedAgentsCustomToolParams{
					Type:        anthropic.BetaManagedAgentsCustomToolParamsTypeCustom,
					Name:        "get_weather",
					Description: "Look up the current weather for a city.",
					InputSchema: anthropic.BetaManagedAgentsCustomToolInputSchemaParam{
						Type:       anthropic.BetaManagedAgentsCustomToolInputSchemaTypeObject,
						Properties: map[string]any{"city": map[string]any{"type": "string"}},
						Required:   []string{"city"},
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created agent v1:", agentV1.ID)

	// Patch the agent to v2 by adding skills; each update bumps the version
	agent, err := client.Beta.Agents.Update(ctx, agentV1.ID, anthropic.BetaAgentUpdateParams{
		Version: agentV1.Version,
		Skills: []anthropic.BetaManagedAgentsSkillParamsUnion{
			{
				OfCustom: &anthropic.BetaManagedAgentsCustomSkillParams{
					Type:    anthropic.BetaManagedAgentsCustomSkillParamsTypeCustom,
					SkillID: skill.ID,
				},
			},
			{
				OfAnthropic: &anthropic.BetaManagedAgentsAnthropicSkillParams{
					Type:    anthropic.BetaManagedAgentsAnthropicSkillParamsTypeAnthropic,
					SkillID: "xlsx",
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Patched agent to v2:", agent.ID)

	// List agent versions
	versions := client.Beta.Agents.Versions.ListAutoPaging(ctx, agent.ID, anthropic.BetaAgentVersionListParams{})
	for versions.Next() {
		v := versions.Current()
		fmt.Printf("  version %d (created %s)\n", v.Version, v.CreatedAt)
	}
	if versions.Err() != nil {
		panic(versions.Err())
	}

	// Create a session pinned to v2; the vault supplies the MCP credential
	session, err := client.Beta.Sessions.New(ctx, anthropic.BetaSessionNewParams{
		EnvironmentID: environment.ID,
		Agent: anthropic.BetaSessionNewParamsAgentUnion{
			OfBetaManagedAgentsAgents: &anthropic.BetaManagedAgentsAgentParams{
				ID:      agent.ID,
				Type:    anthropic.BetaManagedAgentsAgentParamsTypeAgent,
				Version: param.NewOpt(agent.Version),
			},
		},
		VaultIDs: []string{vault.ID},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created session:", session.ID)

	// Send a prompt and stream events, answering the custom tool if called
	fmt.Println("Streaming events:")
	_, err = client.Beta.Sessions.Events.Send(ctx, session.ID, anthropic.BetaSessionEventSendParams{
		Events: []anthropic.BetaManagedAgentsEventParamsUnion{
			{
				OfUserMessage: &anthropic.BetaManagedAgentsUserMessageEventParams{
					Type: anthropic.BetaManagedAgentsUserMessageEventParamsTypeUserMessage,
					Content: []anthropic.BetaManagedAgentsUserMessageEventParamsContentUnion{
						{
							OfText: &anthropic.BetaManagedAgentsTextBlockParam{
								Text: prompt,
								Type: anthropic.BetaManagedAgentsTextBlockTypeText,
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	stream := client.Beta.Sessions.Events.StreamEvents(ctx, session.ID, anthropic.BetaSessionEventStreamParams{})
	for stream.Next() {
		event := stream.Current()
		data, _ := json.MarshalIndent(event, "", "  ")
		fmt.Println(string(data))

		if event.Type == "agent.custom_tool_use" && event.Name == "get_weather" {
			_, err = client.Beta.Sessions.Events.Send(ctx, session.ID, anthropic.BetaSessionEventSendParams{
				Events: []anthropic.BetaManagedAgentsEventParamsUnion{
					{
						OfUserCustomToolResult: &anthropic.BetaManagedAgentsUserCustomToolResultEventParams{
							Type:            anthropic.BetaManagedAgentsUserCustomToolResultEventParamsTypeUserCustomToolResult,
							CustomToolUseID: event.ID,
							Content: []anthropic.BetaManagedAgentsUserCustomToolResultEventParamsContentUnion{
								{
									OfText: &anthropic.BetaManagedAgentsTextBlockParam{
										Text: `{"temperature_c": 14}`,
										Type: anthropic.BetaManagedAgentsTextBlockTypeText,
									},
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		}

		if event.Type == "session.status_idle" && event.StopReason.Type == "end_turn" {
			break
		}
	}
	if stream.Err() != nil {
		panic(stream.Err())
	}
}

// namedReader wraps an io.Reader with a custom filename for multipart uploads.
type namedReader struct {
	io.Reader
	name string
}

func (r namedReader) Filename() string { return r.name }
