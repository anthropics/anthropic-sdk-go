package main

import (
	"context"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/vertex"
)

func main() {
	region := getenv("CLOUD_ML_REGION", "us-east5")
	projectID := getenv("ANTHROPIC_VERTEX_PROJECT_ID", "id-xxx")

	client := anthropic.NewClient(
		vertex.WithGoogleAuth(context.Background(), region, projectID),
	)

	content := "Write me a function to call the Anthropic message API in Node.js using the Anthropic Typescript SDK."

	println("[user]: " + content)

	stream := client.Messages.NewStreaming(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(content)),
		},
		Model:         "claude-sonnet-5",
		StopSequences: []string{"```\n"},
	})

	print("[assistant]: ")

	for stream.Next() {
		event := stream.Current()

		switch variant := event.AsAny().(type) {
		case anthropic.ContentBlockDeltaEvent:
			if variant.Delta.Text != "" {
				print(variant.Delta.Text)
			}
		case anthropic.MessageDeltaEvent:
			if variant.Delta.StopSequence != "" {
				print(variant.Delta.StopSequence)
			}
		}

	}

	println()

	if stream.Err() != nil {
		panic(stream.Err())
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
