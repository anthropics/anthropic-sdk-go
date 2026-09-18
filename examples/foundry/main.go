package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/foundry"
)

// Set ANTHROPIC_FOUNDRY_API_KEY and ANTHROPIC_FOUNDRY_RESOURCE before running.
func main() {
	client, err := foundry.NewClient(foundry.ClientConfig{})
	if err != nil {
		panic(err)
	}

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Write me a haiku about Go.")),
		},
		// Model is your Foundry deployment's name, which defaults to the model ID.
		Model: anthropic.ModelClaudeSonnet5,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(message.Content[0].Text)
}
