package bedrock

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestBedrockError(t *testing.T) {
	expectedErr := errors.New("simulated network error")

	client := anthropic.NewClient(
		option.WithMiddleware(func(r *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			return nil, expectedErr
		}),
	)

	// Attempt to make a request
	stream := client.Messages.NewStreaming(context.Background(), anthropic.MessageNewParams{
		MaxTokens: anthropic.Int(1024),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("test message")),
		}),
		Model: anthropic.F("anthropic.claude-3-sonnet-20240229-v1:0"),
	})

	for stream.Next() {
		stream.Current()
	}

	if stream.Err() != expectedErr {
		t.Fatal()
	}
}
