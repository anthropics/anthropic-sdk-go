package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	client := anthropic.NewClient()
	ctx := context.TODO()

	environment, err := client.Beta.Environments.New(ctx, anthropic.BetaEnvironmentNewParams{
		Name: "streaming-deltas-example",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created environment:", environment.ID)

	agent, err := client.Beta.Agents.New(ctx, anthropic.BetaAgentNewParams{
		Name: "streaming-deltas-example",
		Model: anthropic.BetaManagedAgentsModelConfigParams{
			ID: anthropic.BetaManagedAgentsModelClaudeSonnet4_6,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created agent:", agent.ID)

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

	stream := client.Beta.Sessions.Events.StreamEvents(ctx, session.ID, anthropic.BetaSessionEventStreamParams{
		EventDeltas: []anthropic.BetaManagedAgentsDeltaType{anthropic.BetaManagedAgentsDeltaTypeAgentMessage},
	})

	_, err = client.Beta.Sessions.Events.Send(ctx, session.ID, anthropic.BetaSessionEventSendParams{
		Events: []anthropic.BetaManagedAgentsEventParamsUnion{
			{
				OfUserMessage: &anthropic.BetaManagedAgentsUserMessageEventParams{
					Type: anthropic.BetaManagedAgentsUserMessageEventParamsTypeUserMessage,
					Content: []anthropic.BetaManagedAgentsUserMessageEventParamsContentUnion{
						{
							OfText: &anthropic.BetaManagedAgentsTextBlockParam{
								Text: "Write a short haiku about the ocean.",
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

	var previews anthropic.BetaManagedAgentsEventAccumulator

	fmt.Println("\nStreaming:")
	for stream.Next() {
		event := stream.Current()
		previews.Accumulate(event)

		switch event.Type {
		case "event_delta":
			fmt.Printf("\r%s", previews.AgentMessageText(event.EventID))

		case "agent.message":
			fmt.Println()
			fmt.Println("[final]", previews.AgentMessageText(event.ID))

		case "session.status_idle":
			if event.StopReason.Type == "end_turn" {
				return
			}

		case "session.error":
			fmt.Println("[error]", event.Error.Type, event.Error.Message)
		}
	}
	if stream.Err() != nil {
		panic(stream.Err())
	}
}
