package main

import (
	"context"
	"fmt"
	"github.com/anthropics/anthropic-sdk-go/option"
	"io"
	"net/http"

	"github.com/anthropics/anthropic-sdk-go"
)

func VerboseLoggingMiddleware() option.Middleware {
	return func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		fmt.Printf("Request: %s %s\n", req.Method, req.URL.String())

		if req.Header != nil {
			for key, values := range req.Header {
				for _, value := range values {
					fmt.Printf("Header: %s: %s\n", key, value)
				}
			}
		}
		if req.Body != nil {
			body, err := req.GetBody()
			if err != nil {
				fmt.Printf("Error getting request body: %v\n", err)
				return nil, err
			}
			bodyBytes, err := io.ReadAll(body)
			if err != nil {
				fmt.Printf("Error reading request body: %v\n", err)
				return nil, err
			}
			fmt.Printf("Request body: %s\n", string(bodyBytes))

		}
		resp, err := next(req)
		if resp.Header != nil {
			for key, values := range resp.Header {
				for _, value := range values {
					fmt.Printf("Response Header: %s: %s\n", key, value)
				}
			}
		}
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return nil, err
		}
		fmt.Printf("Response: %d\n", resp.StatusCode)
		return resp, nil
	}
}

func main() {
	client := anthropic.NewClient(option.WithMiddleware(VerboseLoggingMiddleware()))

	content := "Write me a function to call the Anthropic message API in Node.js using the Anthropic Typescript SDK."

	println("[user]: " + content)

	stream := client.Messages.NewStreaming(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(content)),
		},
		Model:         anthropic.ModelClaude_3_5_Sonnet_20240620,
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
