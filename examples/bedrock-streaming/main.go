package main

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/bedrock"
)

func main() {
	client := anthropic.NewClient(
		bedrock.WithLoadDefaultConfig(context.Background()),
	)

	content := "Write me a function to call the Anthropic message API in Node.js using the Anthropic Typescript SDK."

	println("[user]: " + content)

	stream := client.Messages.NewStreaming(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: anthropic.Int(1024),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(content)),
		}),
		Model:         anthropic.F("anthropic.claude-3-sonnet-20240229-v1:0"),
		StopSequences: anthropic.F([]string{"```\n"}),
	})

	print("[assistant]: ")

	for stream.Next() {
		event := stream.Current()

		switch delta := event.Delta.(type) {
		case anthropic.ContentBlockDeltaEventDelta:
			if delta.Text != "" {
				print(delta.Text)
			}
		case anthropic.MessageDeltaEventDelta:
			if delta.StopSequence != "" {
				print(delta.StopSequence)
			}
		}
	}

	println()

	if stream.Err() != nil {
		panic(stream.Err())
	}
}
