package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	client := anthropic.NewClient()
	ctx := context.TODO()

	// NOTE: Outcome events (user.define_outcome) are not yet supported in the Go SDK.

	// Create an environment
	environment, err := client.Beta.Environments.New(ctx, anthropic.BetaEnvironmentNewParams{
		Name: "example-environment",
	})
	if err != nil {
		panic(err)
	}

	// Create an agent
	agent, err := client.Beta.Agents.New(ctx, anthropic.BetaAgentNewParams{
		Name: "example-agent",
		Model: anthropic.BetaManagedAgentsModelConfigParams{
			ID: anthropic.BetaManagedAgentsModelClaudeSonnet4_6,
		},
	})
	if err != nil {
		panic(err)
	}

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
}
