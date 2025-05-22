package main

import (
	"context"
	"fmt"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
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
	client := anthropic.NewClient(option.WithHeader("anthropic-beta", anthropic.AnthropicBetaMCPClient2025_04_04),
		option.WithMiddleware(VerboseLoggingMiddleware()))

	mcpServers := []anthropic.BetaRequestMCPServerURLDefinitionParam{
		{
			URL:                "http://example-server.modelcontextprotocol.io/sse",
			Name:               "example",
			AuthorizationToken: param.NewOpt("YOUR_TOKEN"),
			ToolConfiguration: anthropic.BetaRequestMCPServerToolConfigurationParam{
				Enabled:      anthropic.Bool(true),
				AllowedTools: []string{"echo", "add"},
			},
		},
	}

	stream := client.Beta.Messages.NewStreaming(context.TODO(), anthropic.BetaMessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("what is 1+1?")),
		},
		MCPServers:    mcpServers,
		Model:         anthropic.ModelClaude3_7Sonnet20250219,
		StopSequences: []string{"```\n"},
	})

	message := anthropic.BetaMessage{}
	for stream.Next() {
		event := stream.Current()
		err := message.Accumulate(event)
		if err != nil {
			fmt.Printf("error accumulating event: %v\n", err)
			continue
		}

		switch eventVariant := event.AsAny().(type) {
		case anthropic.BetaRawMessageDeltaEvent:
			print(eventVariant.Delta.StopSequence)
		case anthropic.BetaRawContentBlockDeltaEvent:
			switch deltaVariant := eventVariant.Delta.AsAny().(type) {
			case anthropic.BetaTextDelta:
				print(deltaVariant.Text)
			}
		default:
			fmt.Printf("%+v\n", eventVariant)
		}
	}

	println()

	if stream.Err() != nil {
		panic(stream.Err())
	}

}
