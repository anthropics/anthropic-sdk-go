package main

import (
	"context"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/bedrock"
)

func main() {
	token := os.Getenv("AWS_BEARER_TOKEN_BEDROCK")
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	client := anthropic.NewClient(
		bedrock.WithBearerToken(token, region),
	)

	content := "Write a haiku about Go programming."

	println("[user]: " + content)

	stream := client.Messages.NewStreaming(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(content)),
		},
		Model: "us.anthropic.claude-sonnet-4-20250514-v1:0",
	})

	print("[assistant]: ")

	for stream.Next() {
		event := stream.Current()

		switch eventVariant := event.AsAny().(type) {
		case anthropic.MessageDeltaEvent:
			print(eventVariant.Delta.StopSequence)
		case anthropic.ContentBlockDeltaEvent:
			switch deltaVariant := eventVariant.Delta.AsAny().(type) {
			case anthropic.TextDelta:
				print(deltaVariant.Text)
			}
		}
	}

	println()

	if stream.Err() != nil {
		println("Stream error:", stream.Err().Error())
	}
}
