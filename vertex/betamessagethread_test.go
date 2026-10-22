package vertex

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/anthropics/anthropic-sdk-go/v2"
	sdkoption "github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/toolrunner"
)

func staticCredentials() *google.Credentials {
	return &google.Credentials{TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "test-token"})}
}

// setApplicationDefaultCredentials points Application Default Credentials at a
// file that is read without a request.
func setApplicationDefaultCredentials(t *testing.T) {
	t.Helper()
	credFile := filepath.Join(t.TempDir(), "creds.json")
	credJSON := `{"type":"authorized_user","client_id":"id","client_secret":"secret","refresh_token":"token"}`
	if err := os.WriteFile(credFile, []byte(credJSON), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credFile)
}

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

func TestWithCredentials_MessageThreadSendFailsLocally(t *testing.T) {
	server, requests := recordingServer(t)
	client := anthropic.NewClient(WithCredentials(context.Background(), "us-east5", "proj", staticCredentials()), sdkoption.WithBaseURL(server.URL))
	thread := client.Beta.Messages.NewThread(threadParams())

	_, err := thread.Send(context.Background(), hello())

	requireLocalError(t, err)
	requireNoRequests(t, requests)
}

func TestWithGoogleAuth_MessageThreadSendFailsLocally(t *testing.T) {
	setApplicationDefaultCredentials(t)
	server, requests := recordingServer(t)
	client := anthropic.NewClient(WithGoogleAuth(context.Background(), "us-east5", "proj"), sdkoption.WithBaseURL(server.URL))
	thread := client.Beta.Messages.NewThread(threadParams())

	_, err := thread.Send(context.Background(), hello())

	requireLocalError(t, err)
	requireNoRequests(t, requests)
}

// recordingServer answers the first requests with first and every later one
// with a message, and records the last request.
func recordingServer(t *testing.T, first ...string) (*httptest.Server, *recordedRequest) {
	t.Helper()
	recorded := &recordedRequest{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		count := int(recorded.count.Add(1))
		recorded.body.Store(string(body))
		recorded.betaHeader.Store(strings.Join(r.Header.Values("anthropic-beta"), ","))
		w.Header().Set("Content-Type", "application/json")
		if count <= len(first) {
			_, _ = io.WriteString(w, first[count-1])
			return
		}
		_, _ = w.Write([]byte(`{"id":"msg_a","type":"message","role":"assistant","content":[{"type":"text","text":"hi"}],"model":"claude-sonnet-4-5","stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`))
	}))
	t.Cleanup(server.Close)
	return server, recorded
}

type recordedRequest struct {
	count      atomic.Int32
	body       atomic.Value
	betaHeader atomic.Value
}

func TestWithCredentials_ToolRunnerSendsTheWholeHistoryWithoutAThread(t *testing.T) {
	server, recorded := recordingServer(t, `{"id":"msg_call","type":"message","role":"assistant","content":[{"type":"tool_use","id":"toolu_1","name":"weather","input":{}}],"model":"claude-sonnet-4-5","stop_reason":"tool_use","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`)
	client := anthropic.NewClient(WithCredentials(context.Background(), "us-east5", "proj", staticCredentials()), sdkoption.WithBaseURL(server.URL))
	weather := toolrunner.NewBetaTool("weather", "weather", anthropic.BetaToolInputSchemaParam{}, func(context.Context, struct{}) (anthropic.BetaToolResultBlockParamContentUnion, error) {
		return anthropic.BetaToolResultBlockParamContentUnion{OfText: &anthropic.BetaTextBlockParam{Text: "sunny"}}, nil
	})
	params := anthropic.BetaToolRunnerParams{BetaMessageNewParams: anthropic.BetaMessageNewParams{
		Model:     threadParams().Model,
		MaxTokens: threadParams().MaxTokens,
		Messages:  hello().Messages,
	}}
	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{weather}, params)

	_, err := runner.RunToCompletion(context.Background())

	if err != nil {
		t.Fatalf("RunToCompletion failed: %v", err)
	}
	if got := recorded.count.Load(); got != 2 {
		t.Fatalf("Expected 2 requests, got %d", got)
	}
	body := recorded.body.Load().(string)
	for _, want := range []string{`"name":"weather"`, `"text":"Hello"`, `"text":"sunny"`} {
		if !strings.Contains(body, want) {
			t.Errorf("Expected %s in the second request, got body %s", want, body)
		}
	}
	if strings.Contains(body, `"thread"`) {
		t.Errorf("Expected a request without a thread, got body %s", body)
	}
	if got := recorded.betaHeader.Load().(string); got != "" {
		t.Errorf("Expected no anthropic-beta header, got %q", got)
	}
}
