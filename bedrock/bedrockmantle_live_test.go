package bedrock_test

import (
	"context"
	"os"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/bedrock"
)

// Live integration tests for Bedrock Mantle. Skipped unless ANTHROPIC_LIVE=1.
//
// Required env vars vary by auth mode:
//
//	API key mode:  AWS_BEARER_TOKEN_BEDROCK (or ANTHROPIC_AWS_API_KEY),
//	               AWS_REGION
//
//	SigV4 mode:    AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION
//
// Run: ANTHROPIC_LIVE=1 go test ./bedrock/... -run TestLiveMantle -v

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

func sendMantleMessage(t *testing.T, client *bedrock.MantleClient) {
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

func TestLiveMantleAPIKey(t *testing.T) {
	skipUnlessLive(t)
	requireEnv(t, "AWS_REGION")

	// Need at least one of these for the API key
	apiKey := os.Getenv("AWS_BEARER_TOKEN_BEDROCK")
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_AWS_API_KEY")
	}
	if apiKey == "" {
		t.Fatal("required env var AWS_BEARER_TOKEN_BEDROCK or ANTHROPIC_AWS_API_KEY is not set")
	}

	client, err := bedrock.NewMantleClient(context.Background(), bedrock.MantleClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		t.Fatalf("NewMantleClient failed: %v", err)
	}

	sendMantleMessage(t, client)
}

func TestLiveMantleSigV4ExplicitCreds(t *testing.T) {
	skipUnlessLive(t)
	requireEnv(t, "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_REGION")

	client, err := bedrock.NewMantleClient(context.Background(), bedrock.MantleClientConfig{
		AWSAccessKey:       os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWSSessionToken:    os.Getenv("AWS_SESSION_TOKEN"),
	})
	if err != nil {
		t.Fatalf("NewMantleClient failed: %v", err)
	}

	sendMantleMessage(t, client)
}

func TestLiveMantleSigV4DefaultChain(t *testing.T) {
	skipUnlessLive(t)
	requireEnv(t, "AWS_REGION")

	// Clear all API key env vars so the default AWS credential chain is used
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "")

	client, err := bedrock.NewMantleClient(context.Background(), bedrock.MantleClientConfig{})
	if err != nil {
		t.Fatalf("NewMantleClient failed (default AWS credential chain): %v", err)
	}

	sendMantleMessage(t, client)
}

func TestLiveMantleSigV4ProfileFromCredentialsFile(t *testing.T) {
	skipUnlessLive(t)
	requireEnv(t, "AWS_REGION", "AWS_PROFILE")

	// Clear explicit creds and API keys so the SDK must resolve from ~/.aws/credentials
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("ANTHROPIC_AWS_API_KEY", "")
	t.Setenv("ANTHROPIC_API_KEY", "")
	t.Setenv("AWS_ACCESS_KEY_ID", "")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "")
	t.Setenv("AWS_SESSION_TOKEN", "")

	client, err := bedrock.NewMantleClient(context.Background(), bedrock.MantleClientConfig{
		AWSProfile: os.Getenv("AWS_PROFILE"),
	})
	if err != nil {
		t.Fatalf("NewMantleClient failed (profile from credentials file): %v", err)
	}

	sendMantleMessage(t, client)
}
