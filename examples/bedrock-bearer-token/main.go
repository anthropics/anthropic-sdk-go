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

	content := "Write me a function to call the Anthropic message API in Node.js using the Anthropic Typescript SDK."

	println("[user]: " + content)

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(content)),
		},
		Model:         "us.anthropic.claude-sonnet-4-20250514-v1:0",
		StopSequences: []string{"```\n"},
	})
	if err != nil {
		panic(err)
	}

	println("[assistant]: " + message.Content[0].Text + message.StopSequence)
}
