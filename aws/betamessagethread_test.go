package aws

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
	t.Setenv("ANTHROPIC_AWS_WORKSPACE_ID", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "")
	client, captured := newTestClient(t, ClientConfig{AWSRegion: "us-east-1", AWSAccessKey: "test-access-key", AWSSecretAccessKey: "test-secret-key", WorkspaceID: "ws-test"})
	thread := client.Beta.Messages.NewThread(anthropic.BetaMessageThreadCreateParams{
		Model:     "claude-sonnet-4-5",
		MaxTokens: 1024,
	})

	_, err := thread.Send(context.Background(), hello())

	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	if got := captured.Headers.Get("anthropic-beta"); got != "message-threads-2026-08-12" {
		t.Errorf("Expected anthropic-beta %q, got %q", "message-threads-2026-08-12", got)
	}
	if got := thread.PreviousMessageID(); got != "msg_test" {
		t.Errorf("Expected the next turn to continue from %q, got %q", "msg_test", got)
	}
}
