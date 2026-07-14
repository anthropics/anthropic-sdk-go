package main

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/snowflake"
)

func main() {
	// Set the SNOWFLAKE_AUTH_TOKEN environment variable, or pass the token directly.
	client := anthropic.NewClient(
		snowflake.WithAccount("your-account", ""),
	)

	content := "Write me a function to call the Anthropic message API in Node.js using the Anthropic Typescript SDK."

	println("[user]: " + content)

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(content)),
		},
		Model:         "claude-3-5-sonnet",
		StopSequences: []string{"```\n"},
	})
	if err != nil {
		panic(err)
	}

	println("[assistant]: " + message.Content[0].Text + message.StopSequence)
}
