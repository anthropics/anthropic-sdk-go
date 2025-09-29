package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

func main() {
	client := anthropic.NewClient(option.WithHeader("anthropic-beta", anthropic.AnthropicBetaMCPClient2025_04_04))

	mcpServers := []anthropic.BetaRequestMCPServerURLDefinitionParam{
		{
			URL:                "http://example-server.modelcontextprotocol.io/sse",
			Name:               "example",
			AuthorizationToken: param.NewOpt("YOUR_TOKEN"),
			ToolConfiguration: anthropic.BetaRequestMCPServerToolConfigurationParam{
				Enabled:      anthropic.Bool(true),
				AllowedTools: []string{"echo", "add"},
			},
		},
	}

	stream := client.Beta.Messages.NewStreaming(context.TODO(), anthropic.BetaMessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("what is 1+1?")),
		},
		MCPServers:    mcpServers,
		Model:         anthropic.ModelClaudeSonnet4_5_20250929,
		StopSequences: []string{"```\n"},
	})

	message := anthropic.BetaMessage{}
	for stream.Next() {
		event := stream.Current()
		err := message.Accumulate(event)
		if err != nil {
			fmt.Printf("error accumulating event: %v\n", err)
			continue
		}

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
