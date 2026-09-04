package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/vertex"
)

func main() {
	region := getenv("CLOUD_ML_REGION", "us-east5")
	projectID := getenv("ANTHROPIC_VERTEX_PROJECT_ID", "id-xxx")

	client := anthropic.NewClient(
		// Register middleware before the Vertex option so it observes
		// Anthropic-shaped requests (POST /v1/messages, model in the body);
		// the Vertex adaptation runs closest to the wire.
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			log.Printf("request: %s %s", req.Method, req.URL.Path)
			return next(req)
		}),
		vertex.WithGoogleAuth(context.Background(), region, projectID),
	)
	content := "Write me a function to call the Anthropic message API in Node.js using the Anthropic Typescript SDK."

	println("[user]: " + content)

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(content)),
		},
		Model:         "claude-sonnet-5",
		StopSequences: []string{"```\n"},
	})
	if err != nil {
		panic(err)
	}

	println("[assistant]: " + message.Content[0].Text + message.StopSequence)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
