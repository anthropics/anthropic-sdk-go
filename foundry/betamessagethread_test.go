package foundry

import (
	"context"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/v2"
)

func hello() anthropic.BetaMessageThreadSendParams {
	return anthropic.BetaMessageThreadSendParams{
		Messages: []anthropic.BetaMessageParam{anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("Hello"))},
	}
}

func TestNewClient_MessageThreadSendCreatesAThread(t *testing.T) {
	clearEnv(t)
	server, requests := newServer(t)
	client := newTestClient(t, server, ClientConfig{APIKey: "test-key"})
	thread := client.Beta.Messages.NewThread(anthropic.BetaMessageThreadCreateParams{
		Model:     "claude-sonnet-4-5",
		MaxTokens: 1024,
	})

	_, err := thread.Send(context.Background(), hello())

	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if got := onlyRequest(t, requests).Headers.Get("anthropic-beta"); got != "message-threads-2026-08-12" {
		t.Errorf("Expected anthropic-beta %q, got %q", "message-threads-2026-08-12", got)
	}
	if got := thread.PreviousMessageID(); got != "msg_test" {
		t.Errorf("Expected the next turn to continue from %q, got %q", "msg_test", got)
	}
}
