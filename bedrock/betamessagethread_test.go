package bedrock

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
)

func threadParams() anthropic.BetaMessageThreadCreateParams {
	return anthropic.BetaMessageThreadCreateParams{
		Model:     "claude-sonnet-4-5",
		MaxTokens: 1024,
	}
}

func hello() anthropic.BetaMessageThreadSendParams {
	return anthropic.BetaMessageThreadSendParams{
		Messages: []anthropic.BetaMessageParam{anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("Hello"))},
	}
}

// requireLocalError requires an error that did not come from a response.
func requireLocalError(t *testing.T, err error) {
	t.Helper()
	var apiErr *anthropic.Error
	if err == nil || errors.As(err, &apiErr) {
		t.Fatalf("Expected an error raised without a response, got: %v", err)
	}
}

func requireNoRequests(t *testing.T, requests *recordedRequest) {
	t.Helper()
	if got := requests.count.Load(); got != 0 {
		t.Fatalf("Expected no request, got %d", got)
	}
}

func setStaticAWSEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	t.Setenv("AWS_ACCESS_KEY_ID", "test-access-key")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test-secret-key")
	t.Setenv("AWS_REGION", "us-east-1")
	t.Setenv("AWS_CONFIG_FILE", t.TempDir()+"/config")
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", t.TempDir()+"/credentials")
	t.Setenv("AWS_EC2_METADATA_DISABLED", "true")
}

func TestWithConfig_MessageThreadSendFailsLocally(t *testing.T) {
	t.Setenv("AWS_BEARER_TOKEN_BEDROCK", "")
	server, requests := recordingServer(t)
	client := anthropic.NewClient(WithConfig(makeStaticAWSConfig("us-east-1")), option.WithBaseURL(server.URL))
	thread := client.Beta.Messages.NewThread(threadParams())

	_, err := thread.Send(context.Background(), hello())

	requireLocalError(t, err)
	requireNoRequests(t, requests)
}

func TestWithLoadDefaultConfig_MessageThreadSendFailsLocally(t *testing.T) {
	setStaticAWSEnvironment(t)
	server, requests := recordingServer(t)
	client := anthropic.NewClient(WithLoadDefaultConfig(context.Background()), option.WithBaseURL(server.URL))
	thread := client.Beta.Messages.NewThread(threadParams())

	_, err := thread.Send(context.Background(), hello())

	requireLocalError(t, err)
	requireNoRequests(t, requests)
}

func TestNewMantleClient_MessageThreadSendFailsLocally(t *testing.T) {
	server, requests := recordingServer(t)
	client, err := NewMantleClient(context.Background(), MantleClientConfig{APIKey: "test-key", AWSRegion: "us-east-1", BaseURL: server.URL})
	if err != nil {
		t.Fatalf("NewMantleClient failed: %v", err)
	}
	thread := client.Beta.Messages.NewThread(threadParams())

	_, err = thread.Send(context.Background(), hello())

	requireLocalError(t, err)
	requireNoRequests(t, requests)
}

// recordingServer answers every request with a message and counts the requests.
func recordingServer(t *testing.T) (*httptest.Server, *recordedRequest) {
	t.Helper()
	recorded := &recordedRequest{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorded.count.Add(1)
		writeMessagesResponse(w)
	}))
	t.Cleanup(server.Close)
	return server, recorded
}

type recordedRequest struct {
	count atomic.Int32
}
