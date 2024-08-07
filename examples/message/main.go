package main

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	client := anthropic.NewClient()

	content := "Write me a function to call the Anthropic message API in Node.js using the Anthropic Typescript SDK."

	println("[user]: " + content)

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: anthropic.Int(1024),
		System: anthropic.F([]anthropic.TextBlockParam{
			anthropic.NewTextBlock("Be very serious"),
		}),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(content)),
		}),
		Model:         anthropic.F(anthropic.ModelClaude_3_5_Sonnet_20240620),
		StopSequences: anthropic.F([]string{"```\n"}),
	})
	if err != nil {
		panic(err)
	}

	println("[assistant]: " + message.Content[0].Text + message.StopSequence)
}
