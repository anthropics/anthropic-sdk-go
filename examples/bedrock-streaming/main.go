package main

import (
	"context"
	"log"
	"net/http"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/bedrock"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// loggingMiddleware logs each request and response. Because it is registered
// before bedrock.WithLoadDefaultConfig, it runs outside the Bedrock
// adaptation: it observes the Anthropic-shaped request (POST /v1/messages,
// model in the body, no AWS signature) and the SSE-formatted streaming
// response — exactly what it would observe against the first-party API.
func loggingMiddleware(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
	log.Printf("request: %s %s", req.Method, req.URL.Path)
	res, err := next(req)
	if err == nil {
		log.Printf("response: %s (%s)", res.Status, res.Header.Get("Content-Type"))
	}
	return res, err
}

func main() {
	client := anthropic.NewClient(
		// Register middleware before the Bedrock option so it sees
		// Anthropic-shaped traffic; the Bedrock adaptation (URL/body rewrite,
		// SigV4 signing, SSE normalization) runs closest to the wire.
		option.WithMiddleware(loggingMiddleware),
		bedrock.WithLoadDefaultConfig(context.Background()),
	)

	content := "Write me a function to call the Anthropic message API in Node.js using the Anthropic Typescript SDK."

	println("[user]: " + content)

	stream := client.Messages.NewStreaming(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(content)),
		},
		Model:         "us.anthropic.claude-sonnet-4-5-20250929-v1:0",
		StopSequences: []string{"```\n"},
	})

	print("[assistant]: ")

	for stream.Next() {
		event := stream.Current()

		switch eventVariant := event.AsAny().(type) {
		case anthropic.MessageDeltaEvent:
			print(eventVariant.Delta.StopSequence)
		case anthropic.ContentBlockDeltaEvent:
			switch deltaVariant := eventVariant.Delta.AsAny().(type) {
			case anthropic.TextDelta:
				print(deltaVariant.Text)
			}
		}
	}

	println()

	if stream.Err() != nil {
		panic(stream.Err())
	}
}
