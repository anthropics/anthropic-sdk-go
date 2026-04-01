package aws_test

import (
	"context"
	"os"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/aws"
)

// Live integration tests for the AWS gateway client. Skipped unless ANTHROPIC_LIVE=1.
//
// Required env vars vary by auth mode:
//
//	API key mode:  ANTHROPIC_AWS_API_KEY, ANTHROPIC_AWS_WORKSPACE_ID, AWS_REGION
//
//	SigV4 mode:    AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION,
//	               ANTHROPIC_AWS_WORKSPACE_ID
//
// Run: ANTHROPIC_LIVE=1 go test ./aws/... -run TestLiveAWS -v

func skipUnlessLive(t *testing.T) {
	t.Helper()
	if os.Getenv("ANTHROPIC_LIVE") != "1" {
		t.Skip("set ANTHROPIC_LIVE=1 to run live integration tests")
	}
}

func requireEnv(t *testing.T, names ...string) {
	t.Helper()
	for _, name := range names {
		if os.Getenv(name) == "" {
			t.Fatalf("required env var %s is not set", name)
		}
	}
}

func liveModel() string {
	if m := os.Getenv("ANTHROPIC_LIVE_MODEL"); m != "" {
		return m
	}
	return "claude-sonnet-4-6"
}

func sendAWSMessage(t *testing.T, client *aws.Client) {
	t.Helper()

	message, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     liveModel(),
		MaxTokens: 32,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Say exactly: hello")),
		},
	})
	if err != nil {
		t.Fatalf("Messages.New failed: %v", err)
	}
	if len(message.Content) == 0 {
		t.Fatal("expected non-empty content in response")
	}
	t.Logf("response: %s", message.Content[0].Text)
}

func TestLiveAWSAPIKey(t *testing.T) {
	skipUnlessLive(t)
	requireEnv(t, "ANTHROPIC_AWS_API_KEY", "ANTHROPIC_AWS_WORKSPACE_ID", "AWS_REGION")

	client, err := aws.NewClient(context.Background(), aws.ClientConfig{})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	sendAWSMessage(t, client)
}

func TestLiveAWSSigV4ExplicitCreds(t *testing.T) {
	skipUnlessLive(t)
	requireEnv(t, "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_REGION", "ANTHROPIC_AWS_WORKSPACE_ID")

	// Clear API key so SigV4 is used
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	client, err := aws.NewClient(context.Background(), aws.ClientConfig{
		AWSAccessKey:       os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSSessionToken:    os.Getenv("AWS_SESSION_TOKEN"),
	})
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	sendAWSMessage(t, client)
}

func TestLiveAWSSigV4DefaultChain(t *testing.T) {
	skipUnlessLive(t)
	requireEnv(t, "AWS_REGION", "ANTHROPIC_AWS_WORKSPACE_ID")

	// Clear all API key env vars so the default AWS credential chain is used
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")

	client, err := aws.NewClient(context.Background(), aws.ClientConfig{})
	if err != nil {
		t.Fatalf("NewClient failed (default AWS credential chain): %v", err)
	}

	sendAWSMessage(t, client)
}

func TestLiveAWSSigV4ProfileFromCredentialsFile(t *testing.T) {
	skipUnlessLive(t)
	requireEnv(t, "AWS_REGION", "ANTHROPIC_AWS_WORKSPACE_ID", "AWS_PROFILE")

	// Clear explicit creds and API keys so the SDK must resolve from ~/.aws/credentials
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "")
	t.Setenv("AWS_SESSION_TOKEN", "")

	client, err := aws.NewClient(context.Background(), aws.ClientConfig{
		AWSProfile: os.Getenv("AWS_PROFILE"),
	})
	if err != nil {
		t.Fatalf("NewClient failed (profile from credentials file): %v", err)
	}

	sendAWSMessage(t, client)
}
