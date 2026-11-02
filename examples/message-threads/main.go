// Example of a conversation on a message thread, where each turn sends only
// the messages that are new.
package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go/v2"
)

func main() {
	client := anthropic.NewClient()
	ctx := context.Background()

	thread := client.Beta.Messages.NewThread(anthropic.BetaMessageThreadCreateParams{
		Model:     anthropic.ModelClaudeSonnet5_5,
		MaxTokens: 1024,
		System:    []anthropic.BetaTextBlockParam{{Text: "You are a concise assistant."}},
	})

	questions := []string{
		"What is a goroutine?",
		"How is it different from a thread of the operating system?",
	}
	for _, question := range questions {
		fmt.Printf("> %s\n", question)
		reply, err := thread.Send(ctx, anthropic.BetaMessageThreadSendParams{
			Messages: []anthropic.BetaMessageParam{
				anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(question)),
			},
		})
		if err != nil {
			fmt.Printf("Error sending the turn: %v\n", err)
			return
		}
		for _, block := range reply.Content {
			if block.Type == "text" {
				fmt.Println(block.Text)
			}
		}
	}

	fmt.Println("> And in one sentence?")
	stream := thread.SendStreaming(ctx, anthropic.BetaMessageThreadSendParams{
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("And in one sentence?")),
		},
	})
	defer stream.Close()
	for stream.Next() {
		if delta := stream.Current().Delta; delta.Type == "text_delta" {
			fmt.Print(delta.Text)
		}
	}
	fmt.Println()
	if err := stream.Err(); err != nil {
		fmt.Printf("Error reading the stream: %v\n", err)
		return
	}

	fmt.Printf("The thread holds %d messages.\n", len(thread.Messages()))
}
