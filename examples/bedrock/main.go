package main

import (
	"context"
	"log"
	"net/http"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/bedrock"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func main() {
	client := anthropic.NewClient(
		// Register middleware before the Bedrock option so it observes
		// Anthropic-shaped requests (POST /v1/messages, model in the body);
		// the Bedrock adaptation runs closest to the wire.
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			log.Printf("request: %s %s", req.Method, req.URL.Path)
			return next(req)
		}),
		bedrock.WithLoadDefaultConfig(context.Background()),
	)

	content := "Write me a function to call the Anthropic message API in Node.js using the Anthropic Typescript SDK."

	println("[user]: " + content)

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(content)),
		},
		Model:         "us.anthropic.claude-sonnet-4-5-20250929-v1:0",
		StopSequences: []string{"```\n"},
	})
	if err != nil {
		panic(err)
	}

	println("[assistant]: " + message.Content[0].Text + message.StopSequence)
}
