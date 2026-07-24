package anthropic

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/option"
)

// messagesServer scripts POST /v1/messages: the first call returns a tool_use
// for "weather", the second a final text answer.
func messagesServer(t *testing.T) *httptest.Server {
	t.Helper()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/messages" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			http.Error(w, "unexpected", http.StatusNotImplemented)
			return
		}
		var body string
		if calls.Add(1) == 1 {
			body = `{"id":"msg_1","type":"message","role":"assistant","model":"m","content":[{"type":"tool_use","id":"toolu_1","name":"weather","input":{"city":"SF"}}],"stop_reason":"tool_use","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
		} else {
			body = `{"id":"msg_2","type":"message","role":"assistant","model":"m","content":[{"type":"text","text":"done"}],"stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	return server
}

func newTestToolRunnerClient(server *httptest.Server) Client {
	return NewClient(
		option.WithBaseURL(server.URL),
		option.WithAPIKey("test-key"),
		option.WithMaxRetries(0),
	)
}

// runToToolResults drives the runner to completion and returns every
// tool_result block it appended to the conversation.
func runToToolResults(t *testing.T, runner *BetaToolRunner) []*BetaToolResultBlockParam {
	t.Helper()
	if _, err := runner.RunToCompletion(context.Background()); err != nil {
		t.Fatalf("RunToCompletion: %v", err)
	}
	var results []*BetaToolResultBlockParam
	for _, msg := range runner.Messages() {
		for _, c := range msg.Content {
			if c.OfToolResult != nil {
				results = append(results, c.OfToolResult)
			}
		}
	}
	return results
}

func toolResultJSON(t *testing.T, block *BetaToolResultBlockParam) string {
	t.Helper()
	data, err := json.Marshal(block)
	if err != nil {
		t.Fatalf("marshal tool_result: %v", err)
	}
	return string(data)
}

func systemToolChange(block BetaContentBlockParamUnion) BetaMessageParam {
	return BetaMessageParam{Role: BetaMessageParamRoleSystem, Content: []BetaContentBlockParamUnion{block}}
}

func weatherRef() BetaToolChangeToolReferenceParam {
	return BetaToolChangeToolReferenceParam{Name: "weather"}
}

// A tool_use for a tool dropped by tool_removal must be answered exactly like
// a tool that was never registered, without invoking the local tool.
func TestBetaToolRunner_ToolRemoval_MatchesUnknownTool(t *testing.T) {
	removals := map[string]BetaContentBlockParamUnion{
		"typed":    NewBetaToolRemovalBlock(weatherRef()),
		"mid_conv": NewBetaMidConvSystemBlock([]BetaMidConversationSystemBlockParamContentUnion{{OfToolRemoval: NewBetaToolRemovalBlock(weatherRef()).OfToolRemoval}}),
	}
	for name, removal := range removals {
		t.Run(name, func(t *testing.T) {
			weather := &stubBetaTool{name: "weather"}
			removedClient := newTestToolRunnerClient(messagesServer(t))
			removed := removedClient.Beta.Messages.NewToolRunner(
				[]BetaTool{weather},
				BetaToolRunnerParams{BetaMessageNewParams: BetaMessageNewParams{
					Model:     "m",
					MaxTokens: 512,
					Messages: []BetaMessageParam{
						systemToolChange(removal),
						NewBetaUserMessage(NewBetaTextBlock("What's the weather in SF?")),
					},
				}, MaxIterations: 5},
			)
			removedResults := runToToolResults(t, removed)
			if weather.runs.Load() != 0 {
				t.Fatalf("removed tool must not execute, ran %d times", weather.runs.Load())
			}

			// Reference: the same call against a runner that never had the tool.
			neverClient := newTestToolRunnerClient(messagesServer(t))
			never := neverClient.Beta.Messages.NewToolRunner(
				nil,
				BetaToolRunnerParams{BetaMessageNewParams: BetaMessageNewParams{
					Model:     "m",
					MaxTokens: 512,
					Messages: []BetaMessageParam{
						NewBetaUserMessage(NewBetaTextBlock("What's the weather in SF?")),
					},
				}, MaxIterations: 5},
			)
			neverResults := runToToolResults(t, never)

			if len(removedResults) != 1 || len(neverResults) != 1 {
				t.Fatalf("expected one tool_result each, got %d and %d", len(removedResults), len(neverResults))
			}
			if !removedResults[0].IsError.Value {
				t.Fatalf("expected removed-tool result to be an error")
			}
			if got, want := toolResultJSON(t, removedResults[0]), toolResultJSON(t, neverResults[0]); got != want {
				t.Fatalf("removed-tool result differs from never-defined tool result\n got: %s\nwant: %s", got, want)
			}
		})
	}
}

// A tool_addition after an earlier tool_removal re-enables the tool.
func TestBetaToolRunner_ToolAddition_ReenablesTool(t *testing.T) {
	weather := &stubBetaTool{name: "weather"}
	client := newTestToolRunnerClient(messagesServer(t))
	runner := client.Beta.Messages.NewToolRunner(
		[]BetaTool{weather},
		BetaToolRunnerParams{BetaMessageNewParams: BetaMessageNewParams{
			Model:     "m",
			MaxTokens: 512,
			Messages: []BetaMessageParam{
				systemToolChange(NewBetaToolRemovalBlock(weatherRef())),
				systemToolChange(NewBetaToolAdditionBlock(weatherRef())),
				NewBetaUserMessage(NewBetaTextBlock("What's the weather in SF?")),
			},
		}, MaxIterations: 5},
	)
	results := runToToolResults(t, runner)
	if weather.runs.Load() != 1 {
		t.Fatalf("re-added tool should run once, ran %d times", weather.runs.Load())
	}
	if len(results) != 1 {
		t.Fatalf("expected one tool_result, got %d", len(results))
	}
	if results[0].IsError.Value {
		t.Fatalf("expected successful tool_result, got error")
	}
	if len(results[0].Content) != 1 || results[0].Content[0].OfText == nil || results[0].Content[0].OfText.Text != "ok from weather" {
		t.Fatalf("unexpected tool_result content: %+v", results[0].Content)
	}
}

// newWeatherRunner builds a runner whose scripted server first asks for the
// "weather" tool and then answers. prefix messages (e.g. system tool changes)
// are placed before the single user turn.
func newWeatherRunner(t *testing.T, weather *stubBetaTool, prefix ...BetaMessageParam) *BetaToolRunner {
	t.Helper()
	client := newTestToolRunnerClient(messagesServer(t))
	messages := append([]BetaMessageParam{}, prefix...)
	messages = append(messages, NewBetaUserMessage(NewBetaTextBlock("What's the weather in SF?")))
	return client.Beta.Messages.NewToolRunner(
		[]BetaTool{weather},
		BetaToolRunnerParams{BetaMessageNewParams: BetaMessageNewParams{
			Model:     "m",
			MaxTokens: 512,
			Messages:  messages,
		}, MaxIterations: 5},
	)
}

func requireToolNotFound(t *testing.T, results []*BetaToolResultBlockParam) {
	t.Helper()
	if len(results) != 1 {
		t.Fatalf("expected one tool_result, got %d", len(results))
	}
	if !results[0].IsError.Value {
		t.Fatalf("expected removed-tool result to be an error")
	}
	content := results[0].Content
	if len(content) != 1 || content[0].OfText == nil || content[0].OfText.Text != "Error: Tool 'weather' not found" {
		t.Fatalf("expected not-found tool_result, got %+v", content)
	}
}

// A tool_removal supplied through AppendMessages before the model is asked
// (not in the initial params) must be honored: the tool_use it answers is not
// executed and resolves to the not-found result.
func TestBetaToolRunner_ToolRemoval_AppendMessagesBeforeCall(t *testing.T) {
	weather := &stubBetaTool{name: "weather"}
	runner := newWeatherRunner(t, weather)
	runner.AppendMessages(systemToolChange(NewBetaToolRemovalBlock(weatherRef())))

	requireToolNotFound(t, runToToolResults(t, runner))
	if got := weather.runs.Load(); got != 0 {
		t.Fatalf("removed tool must not execute, ran %d times", got)
	}
}

// A tool_removal appended in the dispatch window — after NextMessage has
// returned the assistant's tool_use but before the following NextMessage
// executes it — must also stop the call. Uses the exported Params.Messages
// mutation path.
func TestBetaToolRunner_ToolRemoval_DispatchWindowParamsMutation(t *testing.T) {
	weather := &stubBetaTool{name: "weather"}
	runner := newWeatherRunner(t, weather)

	msg, err := runner.NextMessage(context.Background())
	if err != nil {
		t.Fatalf("NextMessage: %v", err)
	}
	if msg == nil || msg.StopReason != BetaStopReasonToolUse {
		t.Fatalf("expected assistant tool_use turn, got %+v", msg)
	}
	// Tool execution is deferred to the next NextMessage call, so this is the
	// window in which a removal can still take effect.
	if got := weather.runs.Load(); got != 0 {
		t.Fatalf("tool must not run before the next NextMessage call, ran %d times", got)
	}
	runner.Params.Messages = append(runner.Params.Messages, systemToolChange(NewBetaToolRemovalBlock(weatherRef())))

	requireToolNotFound(t, runToToolResults(t, runner))
	if got := weather.runs.Load(); got != 0 {
		t.Fatalf("removed tool must not execute, ran %d times", got)
	}
}

// A tool_addition appended in the dispatch window re-enables a tool that
// the initial params had removed.
func TestBetaToolRunner_ToolAddition_AppendMessagesInDispatchWindow(t *testing.T) {
	weather := &stubBetaTool{name: "weather"}
	runner := newWeatherRunner(t, weather, systemToolChange(NewBetaToolRemovalBlock(weatherRef())))

	if _, err := runner.NextMessage(context.Background()); err != nil {
		t.Fatalf("NextMessage: %v", err)
	}
	runner.AppendMessages(systemToolChange(NewBetaToolAdditionBlock(weatherRef())))

	results := runToToolResults(t, runner)
	if got := weather.runs.Load(); got != 1 {
		t.Fatalf("re-added tool should run once, ran %d times", got)
	}
	if len(results) != 1 || results[0].IsError.Value {
		t.Fatalf("expected one successful tool_result, got %+v", results)
	}
}
