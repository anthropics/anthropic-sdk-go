package main

import (
	"context"
	"fmt"
	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	client := anthropic.NewClient(option.WithHeader("anthropic-beta", "mcp-client-2025-04-04"))

	mcpServers := []anthropic.BetaRequestMCPServerURLDefinitionParam{
		{
			URL:  "http://example-server.modelcontextprotocol.io/sse",
			Name: "example",
			ToolConfiguration: anthropic.BetaRequestMCPServerToolConfigurationParam{
				Enabled:      anthropic.Bool(true),
				AllowedTools: []string{"echo", "add"},
			},
		},
	}

	stream := client.Beta.Messages.NewStreaming(context.TODO(), anthropic.BetaMessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("Calculate 1+2")),
		},
		MCPServers:    mcpServers,
		Model:         anthropic.ModelClaude3_7Sonnet20250219,
		StopSequences: []string{"```\n"},
	})

	for stream.Next() {
		event := stream.Current()

		switch eventVariant := event.AsAny().(type) {
		case anthropic.BetaRawMessageDeltaEvent:
			print(eventVariant.Delta.StopSequence)
		case anthropic.BetaRawContentBlockDeltaEvent:
			switch deltaVariant := eventVariant.Delta.AsAny().(type) {
			case anthropic.BetaTextDelta:
				print(deltaVariant.Text)
			}
		default:
			fmt.Printf("%+v\n", eventVariant)
		}
	}

	println()

	if stream.Err() != nil {
		panic(stream.Err())
	}

}
