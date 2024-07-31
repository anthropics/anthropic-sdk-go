package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	client := anthropic.NewClient()

	content := "How many dogs are in this picture?"

	println("[user]: " + content)

	file, err := os.Open("./multimodal/nine_dogs.png")
	if err != nil {
		panic(fmt.Errorf("failed to open file: you should run this example from the root of the anthropic-go/examples directory: %w", err))
	}
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	fileEncoded := base64.StdEncoding.EncodeToString(fileBytes)

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: anthropic.Int(1024),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock(content),
				anthropic.NewImageBlockBase64("image/png", fileEncoded),
			),
		}),
		Model:         anthropic.F(anthropic.ModelClaude_3_5_Sonnet_20240620),
		StopSequences: anthropic.F([]string{"```\n"}),
	})
	if err != nil {
		panic(err)
	}

	println("[assistant]: " + message.Content[0].Text + message.StopSequence)
}
