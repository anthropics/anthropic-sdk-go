package main

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	// Zero-config workload identity authentication via environment variables.
	// Set the following env vars before running:
	//
	//   ANTHROPIC_FEDERATION_RULE_ID — the federation rule ID
	//   ANTHROPIC_ORGANIZATION_ID   — the organization ID
	//   ANTHROPIC_IDENTITY_TOKEN    — a literal JWT identity token
	//     (or ANTHROPIC_IDENTITY_TOKEN_FILE — path to a file containing the JWT)
	//
	// Optional:
	//   ANTHROPIC_SERVICE_ACCOUNT_ID — service account ID
	//
	// When these are set, NewClient() automatically exchanges the identity token
	// for a short-lived Anthropic access token. If an API key is also set, the
	// API key takes precedence.
	client := anthropic.NewClient()

	content := "Write me a function to call the Anthropic message API in Node.js using the Anthropic Typescript SDK."

	println("[user]: " + content)

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(content)),
		},
		Model:         anthropic.ModelClaudeSonnet4_5_20250929,
		StopSequences: []string{"```\n"},
	})
	if err != nil {
		panic(err)
	}

	println("[assistant]: " + message.Content[0].Text + message.StopSequence)
}
