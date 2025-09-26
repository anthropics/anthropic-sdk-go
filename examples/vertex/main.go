package main

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go/vertex"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	client := anthropic.NewClient(
		vertex.WithGoogleAuth(context.Background(), "us-central1", "id-xxx"),
	)
	content := "Write me a function to call the Anthropic message API in Node.js using the Anthropic Typescript SDK."

	println("[user]: " + content)

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(content)),
		},
		Model:         "claude-sonnet-4-v1@20250514",
		StopSequences: []string{"```\n"},
	})
	if err != nil {
		panic(err)
	}

	println("[assistant]: " + message.Content[0].Text + message.StopSequence)
}
