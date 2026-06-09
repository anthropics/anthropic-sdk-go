// Falls back to a second model when the primary refuses, two ways:
// server-side via the fallbacks param (preferred), and client-side via
// betafallback.BetaRefusalFallbackMiddleware for providers without
// server-side support.
//
// Requires ANTHROPIC_API_KEY.
package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/lib/betafallback"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func main() {
	ctx := context.Background()

	// 1. Server-side fallbacks (preferred): the API retries a refusal itself —
	// one request, a plain client, no client-side logic. Use this when talking
	// to the API directly.
	client := anthropic.NewClient()
	served, err := client.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{
		MaxTokens: 1024,
		Model:     anthropic.ModelClaudeFable5,
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("Some prompt that triggers a refusal")),
		},
		Fallbacks: []anthropic.BetaFallbackParam{{Model: anthropic.ModelClaudeOpus4_8}},
		Betas:     []anthropic.AnthropicBeta{anthropic.AnthropicBetaServerSideFallback2026_06_01},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("server-side, served by:", served.Model)

	// If your provider doesn't support server-side fallbacks, register the
	// client-side middleware instead:
	fallbackClient := anthropic.NewClient(
		option.WithMiddleware(betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: anthropic.ModelClaudeOpus4_8}},
		)),
	)
	state := betafallback.WithBetaFallbackState(&betafallback.BetaFallbackState{}) // pins follow-ups to the model that accepted

	// 2. Streaming: on a refusal the middleware retries and splices the
	// fallback's events onto the open stream — one continuous message, with a
	// `fallback` content block marking the model boundary.
	stream := fallbackClient.Beta.Messages.NewStreaming(ctx, anthropic.BetaMessageNewParams{
		MaxTokens: 1024,
		Model:     anthropic.ModelClaudeFable5,
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("Some prompt that triggers a refusal")),
		},
	}, state)
	defer stream.Close()

	var streamed anthropic.BetaMessage
	for stream.Next() {
		event := stream.Current()
		if err := streamed.Accumulate(event); err != nil {
			panic(err)
		}
		switch event := event.AsAny().(type) {
		case anthropic.BetaRawContentBlockStartEvent:
			// the fallback block marks the splice point
			if fallback, ok := event.ContentBlock.AsAny().(anthropic.BetaFallbackBlock); ok {
				fmt.Printf("\n--- fell back: %s -> %s ---\n", fallback.From.Model, fallback.To.Model)
			}
		case anthropic.BetaRawContentBlockDeltaEvent:
			if delta, ok := event.Delta.AsAny().(anthropic.BetaTextDelta); ok {
				fmt.Print(delta.Text)
			}
		}
	}
	if err := stream.Err(); err != nil {
		panic(err)
	}
	fmt.Println("\nstreaming, served by:", streamed.Model)

	// 3. Non-streaming: same middleware, the retry just happens before you
	// get the message back.
	message, err := fallbackClient.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{
		MaxTokens: 1024,
		Model:     anthropic.ModelClaudeFable5,
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("Some prompt that triggers a refusal")),
		},
	}, state) // reusing the state keeps the conversation pinned
	if err != nil {
		panic(err)
	}
	fmt.Println("non-streaming, served by:", message.Model)
}
