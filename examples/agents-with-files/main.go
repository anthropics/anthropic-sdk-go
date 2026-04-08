package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

func main() {
	client := anthropic.NewClient()
	ctx := context.TODO()

	// Create an environment
	environment, err := client.Beta.Environments.New(ctx, anthropic.BetaEnvironmentNewParams{
		Name: "files-example-environment",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created environment:", environment.ID)

	// Create an agent with the built-in toolset and an always-allow permission policy
	agent, err := client.Beta.Agents.New(ctx, anthropic.BetaAgentNewParams{
		Name: "files-example-agent",
		Model: anthropic.BetaManagedAgentsModelConfigParams{
			ID: anthropic.BetaManagedAgentsModelClaudeSonnet4_6,
		},
		Tools: []anthropic.BetaAgentNewParamsToolUnion{
			{
				OfAgentToolset20260401: &anthropic.BetaManagedAgentsAgentToolset20260401Params{
					Type: anthropic.BetaManagedAgentsAgentToolset20260401ParamsTypeAgentToolset20260401,
					DefaultConfig: anthropic.BetaManagedAgentsAgentToolsetDefaultConfigParams{
						Enabled: param.NewOpt(true),
						PermissionPolicy: anthropic.BetaManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion{
							OfAlwaysAllow: &anthropic.BetaManagedAgentsAlwaysAllowPolicyParam{
								Type: anthropic.BetaManagedAgentsAlwaysAllowPolicyTypeAlwaysAllow,
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
	fmt.Println("Created agent:", agent.ID)

	// Upload a file
	csvFile, err := os.Open("agents-with-files/data.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	file, err := client.Beta.Files.Upload(ctx, anthropic.BetaFileUploadParams{
		File: csvFile,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Uploaded file:", file.ID)

	// Create a session with the file mounted as a resource
	session, err := client.Beta.Sessions.New(ctx, anthropic.BetaSessionNewParams{
		EnvironmentID: environment.ID,
		Agent: anthropic.BetaSessionNewParamsAgentUnion{
			OfBetaManagedAgentsAgents: &anthropic.BetaManagedAgentsAgentParams{
				ID:      agent.ID,
				Type:    anthropic.BetaManagedAgentsAgentParamsTypeAgent,
				Version: param.NewOpt(agent.Version),
			},
		},
		Resources: []anthropic.BetaSessionNewParamsResourceUnion{
			{
				OfFile: &anthropic.BetaManagedAgentsFileResourceParams{
					Type:      anthropic.BetaManagedAgentsFileResourceParamsTypeFile,
					FileID:    file.ID,
					MountPath: param.NewOpt("data.csv"),
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created session:", session.ID)

	// Send a prompt asking the agent to read the mounted file
	fmt.Println("Streaming events:")
	_, err = client.Beta.Sessions.Events.Send(ctx, session.ID, anthropic.BetaSessionEventSendParams{
		Events: []anthropic.BetaManagedAgentsEventParamsUnion{
			{
				OfUserMessage: &anthropic.BetaManagedAgentsUserMessageEventParams{
					Type: anthropic.BetaManagedAgentsUserMessageEventParamsTypeUserMessage,
					Content: []anthropic.BetaManagedAgentsUserMessageEventParamsContentUnion{
						{
							OfText: &anthropic.BetaManagedAgentsTextBlockParam{
								Text: "Read /uploads/data.csv and tell me the column names.",
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

	// Stream events until the session goes idle
	stream := client.Beta.Sessions.Events.StreamEvents(ctx, session.ID, anthropic.BetaSessionEventStreamParams{})
	for stream.Next() {
		event := stream.Current()
		data, _ := json.MarshalIndent(event, "", "  ")
		fmt.Println(string(data))
		if event.Type == "session.status_idle" {
			break
		}
	}
	if stream.Err() != nil {
		panic(stream.Err())
	}
}
