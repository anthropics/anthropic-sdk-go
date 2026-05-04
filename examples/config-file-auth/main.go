package main

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/config"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func main() {
	// LoadConfig reads from ~/.config/anthropic/configs/<profile>.json.
	// The profile is resolved from ANTHROPIC_PROFILE, then the active_config
	// file, then "default". The config directory can be overridden with
	// ANTHROPIC_CONFIG_DIR.
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	client := anthropic.NewClient(
		option.WithConfig(cfg),
	)

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
