package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	client := anthropic.NewClient()
	ctx := context.TODO()

	// Create an environment
	environment, err := client.Beta.Environments.New(ctx, anthropic.BetaEnvironmentNewParams{
		Name: "simple-example-environment",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created environment:", environment.ID)

	// Create an agent
	agent, err := client.Beta.Agents.New(ctx, anthropic.BetaAgentNewParams{
		Name: "simple-example-agent",
		Model: anthropic.BetaManagedAgentsModelConfigParams{
			ID: anthropic.BetaManagedAgentsModelClaudeSonnet4_6,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created agent:", agent.ID)

	// Create a session
	session, err := client.Beta.Sessions.New(ctx, anthropic.BetaSessionNewParams{
		EnvironmentID: environment.ID,
		Agent: anthropic.BetaSessionNewParamsAgentUnion{
			OfBetaManagedAgentsAgents: &anthropic.BetaManagedAgentsAgentParams{
				ID:   agent.ID,
				Type: anthropic.BetaManagedAgentsAgentParamsTypeAgent,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created session:", session.ID)

	// Send a user message
	_, err = client.Beta.Sessions.Events.Send(ctx, session.ID, anthropic.BetaSessionEventSendParams{
		Events: []anthropic.BetaManagedAgentsEventParamsUnion{
			{
				OfUserMessage: &anthropic.BetaManagedAgentsUserMessageEventParams{
					Type: anthropic.BetaManagedAgentsUserMessageEventParamsTypeUserMessage,
					Content: []anthropic.BetaManagedAgentsUserMessageEventParamsContentUnion{
						{
							OfText: &anthropic.BetaManagedAgentsTextBlockParam{
								Text: "Hello Claude!",
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
	fmt.Println("Streaming events:")
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
