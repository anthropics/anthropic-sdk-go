// Changes the effort level partway through a conversation by appending a
// directive-only role:"system" message that carries a per-message
// output_config. The directive applies from that point in the transcript on.
//
// Requires ANTHROPIC_API_KEY and the mid-conversation-output-config beta.
package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	ctx := context.Background()
	client := anthropic.NewClient()

	messages := []anthropic.BetaMessageParam{
		anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("Give me a one-line summary of the Go memory model.")),
	}

	first, err := client.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{
		MaxTokens: 1024,
		Model:     anthropic.ModelClaudeOpus4_8,
		Messages:  messages,
		Betas:     []anthropic.AnthropicBeta{anthropic.AnthropicBetaMidConversationOutputConfig2026_07_01},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("[assistant]:", first.Content[0].Text)

	// Raise the effort for the rest of the conversation. NewBetaSystemMessage
	// with no blocks serializes as {"role":"system","content":[],"output_config":{...}};
	// the API requires the content key even when it is empty.
	messages = append(messages,
		first.ToParam(),
		anthropic.NewBetaSystemMessage(anthropic.BetaSystemMessageOutputConfigParam{
			Effort: anthropic.BetaSystemMessageOutputConfigEffortHigh,
		}),
		anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("Now explain happens-before in detail, with an example.")),
	)

	second, err := client.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{
		MaxTokens: 2048,
		Model:     anthropic.ModelClaudeOpus4_8,
		Messages:  messages,
		Betas:     []anthropic.AnthropicBeta{anthropic.AnthropicBetaMidConversationOutputConfig2026_07_01},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("[assistant, effort=high]:", second.Content[0].Text)
}
