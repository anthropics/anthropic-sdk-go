package main

import (
	"context"
	"fmt"
	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	client := anthropic.NewClient(option.WithHeader("anthropic-beta", "mcp-client-2025-04-04"))

	content := "Use the eBird API to fetch the hotspot details of McGolrick park (L2987624)"

	mcpServers := []anthropic.BetaRequestMCPServerURLDefinitionParam{
		{
			URL:  "https://remote-ebird-mcp-server-authless.dev-66f.workers.dev/sse",
			Name: "ebird",
			ToolConfiguration: anthropic.BetaRequestMCPServerToolConfigurationParam{
				Enabled: anthropic.Bool(true),
			},
		},
	}

	stream := client.Beta.Messages.NewStreaming(context.TODO(), anthropic.BetaMessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(content)),
		},
		MCPServers:    mcpServers,
		Model:         anthropic.ModelClaude3_7Sonnet20250219,
		StopSequences: []string{"```\n"},
	})

	//message := anthropic.BetaMessage{}

	for stream.Next() {
		event := stream.Current()
		//err := message.Accumulate(event)

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
