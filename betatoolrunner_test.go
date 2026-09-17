package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
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
	weather := &stubBetaTool{name: "weather"}
	removedClient := newTestToolRunnerClient(messagesServer(t))
	removed := removedClient.Beta.Messages.NewToolRunner(
		[]BetaTool{weather},
		BetaToolRunnerParams{BetaMessageNewParams: BetaMessageNewParams{
			Model:     "m",
			MaxTokens: 512,
			Messages: []BetaMessageParam{
				systemToolChange(NewBetaToolRemovalBlock(weatherRef())),
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

type recordingTool struct {
	name  string
	calls int
}

func (r *recordingTool) Name() string                          { return r.name }
func (r *recordingTool) Description() string                   { return "records invocations" }
func (r *recordingTool) InputSchema() BetaToolInputSchemaParam { return BetaToolInputSchemaParam{} }
func (r *recordingTool) Execute(ctx context.Context, input json.RawMessage) ([]BetaToolResultBlockParamContentUnion, error) {
	r.calls++
	return []BetaToolResultBlockParamContentUnion{{OfText: &BetaTextBlockParam{Text: "ok"}}}, nil
}

// A cut-off turn may carry a tool call whose input was cut mid-stream, and a
// finished or unfinished turn may carry one nobody is waiting on; the runner
// executes tool calls only when the turn stopped for tool use.
func TestExecuteToolsSkipsMaxTokensTurn(t *testing.T) {
	var toolUse BetaContentBlockUnion
	if err := toolUse.UnmarshalJSON([]byte(`{"type":"tool_use","id":"toolu_1","name":"rec","input":{}}`)); err != nil {
		t.Fatalf("building tool_use block: %v", err)
	}

	for _, tt := range []struct {
		name       string
		stopReason BetaStopReason
		wantCalls  int
		wantResult bool
	}{
		{"max_tokens turn skips execution", BetaStopReasonMaxTokens, 0, false},
		{"context window turn skips execution", BetaStopReasonModelContextWindowExceeded, 0, false},
		{"refusal turn skips execution", BetaStopReasonRefusal, 0, false},
		{"end_turn skips execution", BetaStopReasonEndTurn, 0, false},
		{"pause_turn skips execution", BetaStopReasonPauseTurn, 0, false},
		{"compaction turn skips execution", BetaStopReasonCompaction, 0, false},
		{"tool_use turn executes", BetaStopReasonToolUse, 1, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tool := &recordingTool{name: "rec"}
			base := newBetaToolRunnerBase(nil, []BetaTool{tool}, BetaToolRunnerParams{}, nil)
			message := &BetaMessage{
				StopReason: tt.stopReason,
				Content:    []BetaContentBlockUnion{toolUse},
			}
			result, err := base.executeTools(context.Background(), message)
			if err != nil {
				t.Fatalf("executeTools: %v", err)
			}
			if tool.calls != tt.wantCalls {
				t.Errorf("tool executed %d times, want %d", tool.calls, tt.wantCalls)
			}
			if (result != nil) != tt.wantResult {
				t.Errorf("result presence = %v, want %v", result != nil, tt.wantResult)
			}
		})
	}
}

// Go cannot check that a switch covers every BetaStopReason constant, so this
// hand-listed table is the exhaustiveness check: keep it in step with the
// generated constants and decide each new one's bucket here.
func TestDetermineNextStepFromStopReason(t *testing.T) {
	want := []struct {
		reason   BetaStopReason
		nextStep toolRunnerStep
	}{
		{BetaStopReasonEndTurn, stepStop},
		{BetaStopReasonMaxTokens, stepStop},
		{BetaStopReasonStopSequence, stepStop},
		{BetaStopReasonToolUse, stepRunTools},
		{BetaStopReasonPauseTurn, stepResume},
		{BetaStopReasonCompaction, stepResume},
		{BetaStopReasonRefusal, stepStop},
		{BetaStopReasonModelContextWindowExceeded, stepStop},
		{"", stepStop},
		{"some_future_reason", stepStop},
	}
	for _, tt := range want {
		if got := determineNextStepFromStopReason(tt.reason); got != tt.nextStep {
			t.Errorf("determineNextStepFromStopReason(%q) = %d, want %d", tt.reason, got, tt.nextStep)
		}
	}
}

// containerServer scripts POST /v1/messages like messagesServer, but the
// first turn ran in a server-assigned container; it records each request body.
func containerServer(t *testing.T, bodies *[]map[string]any) *httptest.Server {
	t.Helper()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var got map[string]any
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		*bodies = append(*bodies, got)
		var body string
		if calls.Add(1) == 1 {
			body = `{"id":"msg_1","type":"message","role":"assistant","model":"m","content":[{"type":"tool_use","id":"toolu_1","name":"weather","input":{"city":"SF"}}],"container":{"id":"container_123","expires_at":"2025-01-01T00:00:00Z","skills":[]},"stop_reason":"tool_use","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
		} else {
			body = `{"id":"msg_2","type":"message","role":"assistant","model":"m","content":[{"type":"text","text":"done"}],"container":null,"stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	return server
}

// The follow-up request must name the container the previous turn ran in,
// unless the caller pinned a container themselves.
func TestBetaToolRunner_ForwardsContainer(t *testing.T) {
	for _, tt := range []struct {
		name   string
		pinned BetaMessageNewParamsContainerUnion
		want   any
	}{
		{"adopts server container", BetaMessageNewParamsContainerUnion{}, "container_123"},
		{"keeps pinned id", BetaMessageNewParamsContainerUnion{OfString: String("container_mine")}, "container_mine"},
		{"fills pinned params without id", BetaMessageNewParamsContainerUnion{OfContainers: &BetaContainerParams{}}, map[string]any{"id": "container_123"}},
		{"keeps pinned params id", BetaMessageNewParamsContainerUnion{OfContainers: &BetaContainerParams{ID: String("container_mine")}}, map[string]any{"id": "container_mine"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var bodies []map[string]any
			client := newTestToolRunnerClient(containerServer(t, &bodies))
			runner := client.Beta.Messages.NewToolRunner(
				[]BetaTool{&stubBetaTool{name: "weather"}},
				BetaToolRunnerParams{BetaMessageNewParams: BetaMessageNewParams{
					Model:     "m",
					MaxTokens: 512,
					Container: tt.pinned,
					Messages:  []BetaMessageParam{NewBetaUserMessage(NewBetaTextBlock("What's the weather in SF?"))},
				}, MaxIterations: 5},
			)
			if _, err := runner.RunToCompletion(context.Background()); err != nil {
				t.Fatalf("RunToCompletion: %v", err)
			}
			if len(bodies) != 2 {
				t.Fatalf("expected 2 requests, got %d", len(bodies))
			}
			gotJSON, _ := json.Marshal(bodies[1]["container"])
			wantJSON, _ := json.Marshal(tt.want)
			if string(gotJSON) != string(wantJSON) {
				t.Fatalf("follow-up container = %s, want %s", gotJSON, wantJSON)
			}
		})
	}
}

func TestBetaToolRunnerStreaming_ForwardsContainer(t *testing.T) {
	var bodies []map[string]any
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var got map[string]any
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		bodies = append(bodies, got)
		w.Header().Set("Content-Type", "text/event-stream")
		var events []string
		if calls.Add(1) == 1 {
			events = []string{
				`event: message_start` + "\n" + `data: {"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","model":"m","content":[],"container":null,"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}}`,
				`event: content_block_start` + "\n" + `data: {"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":"toolu_1","name":"weather","input":{}}}`,
				`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":"{\"city\":\"SF\"}"}}`,
				`event: content_block_stop` + "\n" + `data: {"type":"content_block_stop","index":0}`,
				`event: message_delta` + "\n" + `data: {"type":"message_delta","delta":{"stop_reason":"tool_use","stop_sequence":null,"container":{"id":"container_123","expires_at":"2025-01-01T00:00:00Z","skills":[]}},"usage":{"output_tokens":5}}`,
				`event: message_stop` + "\n" + `data: {"type":"message_stop"}`,
			}
		} else {
			events = []string{
				`event: message_start` + "\n" + `data: {"type":"message_start","message":{"id":"msg_2","type":"message","role":"assistant","model":"m","content":[],"container":null,"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}}`,
				`event: content_block_start` + "\n" + `data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`,
				`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"done"}}`,
				`event: content_block_stop` + "\n" + `data: {"type":"content_block_stop","index":0}`,
				`event: message_delta` + "\n" + `data: {"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null},"usage":{"output_tokens":5}}`,
				`event: message_stop` + "\n" + `data: {"type":"message_stop"}`,
			}
		}
		for _, e := range events {
			_, _ = w.Write([]byte(e + "\n\n"))
		}
	}))
	t.Cleanup(server.Close)

	client := newTestToolRunnerClient(server)
	runner := client.Beta.Messages.NewToolRunnerStreaming(
		[]BetaTool{&stubBetaTool{name: "weather"}},
		BetaToolRunnerParams{BetaMessageNewParams: BetaMessageNewParams{
			Model:     "m",
			MaxTokens: 512,
			Messages:  []BetaMessageParam{NewBetaUserMessage(NewBetaTextBlock("What's the weather in SF?"))},
		}, MaxIterations: 5},
	)
	for events := range runner.AllStreaming(context.Background()) {
		for _, err := range events {
			if err != nil {
				t.Fatalf("streaming: %v", err)
			}
		}
	}
	if len(bodies) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(bodies))
	}
	if got := bodies[1]["container"]; got != "container_123" {
		t.Fatalf("follow-up container = %v, want container_123", got)
	}
}

// toolUseServer answers every POST /v1/messages with a fresh tool_use turn
// (msg_1, msg_2, ...), so a runner only stops at MaxIterations.
func toolUseServer(t *testing.T, streaming bool) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		if !streaming {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"id":"msg_%d","type":"message","role":"assistant","model":"m","content":[{"type":"tool_use","id":"toolu_%d","name":"weather","input":{"city":"SF"}}],"stop_reason":"tool_use","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`, n, n)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		for _, e := range []string{
			fmt.Sprintf(`event: message_start`+"\n"+`data: {"type":"message_start","message":{"id":"msg_%d","type":"message","role":"assistant","model":"m","content":[],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}}`, n),
			fmt.Sprintf(`event: content_block_start`+"\n"+`data: {"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":"toolu_%d","name":"weather","input":{}}}`, n),
			`event: content_block_delta` + "\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":"{\"city\":\"SF\"}"}}`,
			`event: content_block_stop` + "\n" + `data: {"type":"content_block_stop","index":0}`,
			`event: message_delta` + "\n" + `data: {"type":"message_delta","delta":{"stop_reason":"tool_use","stop_sequence":null},"usage":{"output_tokens":5}}`,
			`event: message_stop` + "\n" + `data: {"type":"message_stop"}`,
		} {
			_, _ = w.Write([]byte(e + "\n\n"))
		}
	}))
	t.Cleanup(server.Close)
	return server, &calls
}

func newCappedRunnerParams(maxIterations int) BetaToolRunnerParams {
	return BetaToolRunnerParams{BetaMessageNewParams: BetaMessageNewParams{
		Model:     "m",
		MaxTokens: 512,
		Messages:  []BetaMessageParam{NewBetaUserMessage(NewBetaTextBlock("What's the weather in SF?"))},
	}, MaxIterations: maxIterations}
}

func requireYieldedOnce(t *testing.T, ids []string, want ...string) {
	t.Helper()
	if fmt.Sprint(ids) != fmt.Sprint(want) {
		t.Fatalf("yielded messages = %v, want %v", ids, want)
	}
}

// All() yields one message per API call: at MaxIterations the K-th message is
// yielded once and no further request is made.
func TestBetaToolRunner_All_MaxIterationsYieldsEachMessageOnce(t *testing.T) {
	const maxIterations = 3
	server, calls := toolUseServer(t, false)
	client := newTestToolRunnerClient(server)
	runner := client.Beta.Messages.NewToolRunner([]BetaTool{&stubBetaTool{name: "weather"}}, newCappedRunnerParams(maxIterations))

	var ids []string
	for msg, err := range runner.All(context.Background()) {
		if err != nil {
			t.Fatalf("All: %v", err)
		}
		ids = append(ids, msg.ID)
	}
	requireYieldedOnce(t, ids, "msg_1", "msg_2", "msg_3")
	if got := calls.Load(); got != maxIterations {
		t.Fatalf("expected %d requests, got %d", maxIterations, got)
	}
	if got := runner.IterationCount(); got != maxIterations {
		t.Fatalf("expected IterationCount %d, got %d", maxIterations, got)
	}
	if msg, err := runner.NextMessage(context.Background()); msg != nil || err != nil {
		t.Fatalf("NextMessage after completion = (%v, %v), want (nil, nil)", msg, err)
	}
}

// All() yields the final answer once when the model stops using tools.
func TestBetaToolRunner_All_YieldsFinalMessageOnce(t *testing.T) {
	client := newTestToolRunnerClient(messagesServer(t))
	runner := client.Beta.Messages.NewToolRunner([]BetaTool{&stubBetaTool{name: "weather"}}, newCappedRunnerParams(0))

	var ids []string
	for msg, err := range runner.All(context.Background()) {
		if err != nil {
			t.Fatalf("All: %v", err)
		}
		ids = append(ids, msg.ID)
	}
	requireYieldedOnce(t, ids, "msg_1", "msg_2")
	if got := runner.IterationCount(); got != 2 {
		t.Fatalf("expected IterationCount 2, got %d", got)
	}
	if last := runner.LastMessage(); last == nil || last.ID != "msg_2" {
		t.Fatalf("LastMessage = %v, want msg_2", last)
	}
}

// AllStreaming streams one turn per API call: at MaxIterations no message is
// streamed twice and no further request is made.
func TestBetaToolRunnerStreaming_AllStreaming_MaxIterationsStreamsEachMessageOnce(t *testing.T) {
	const maxIterations = 3
	server, calls := toolUseServer(t, true)
	client := newTestToolRunnerClient(server)
	runner := client.Beta.Messages.NewToolRunnerStreaming([]BetaTool{&stubBetaTool{name: "weather"}}, newCappedRunnerParams(maxIterations))

	var ids []string
	for events, err := range runner.AllStreaming(context.Background()) {
		if err != nil {
			t.Fatalf("AllStreaming: %v", err)
		}
		for event, err := range events {
			if err != nil {
				t.Fatalf("streaming: %v", err)
			}
			if start, ok := event.AsAny().(BetaRawMessageStartEvent); ok {
				ids = append(ids, start.Message.ID)
			}
		}
	}
	requireYieldedOnce(t, ids, "msg_1", "msg_2", "msg_3")
	if got := calls.Load(); got != maxIterations {
		t.Fatalf("expected %d requests, got %d", maxIterations, got)
	}
	if last := runner.LastMessage(); last == nil || last.ID != "msg_3" {
		t.Fatalf("LastMessage = %v, want msg_3", last)
	}
}

const (
	pausedTurnJSON = `{"id":"msg_paused","type":"message","role":"assistant","model":"m","content":[{"type":"text","text":"Let me look that up."},{"type":"server_tool_use","id":"srvtoolu_1","name":"web_search","input":{"query":"weather in SF"}}],"stop_reason":"pause_turn","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
	endTurnJSON    = `{"id":"msg_end","type":"message","role":"assistant","model":"m","content":[{"type":"text","text":"done"}],"stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
	// compactedTurnJSON is what pause_after_compaction returns: the summary
	// block, handed back before the model has answered.
	compactedTurnJSON = `{"id":"msg_compacted","type":"message","role":"assistant","model":"m","content":[{"type":"compaction","content":"Summary: the user asked about the weather in SF."}],"stop_reason":"compaction","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
	// cutOffToolUseJSON is a max_tokens turn that still carries a client tool call.
	cutOffToolUseJSON = `{"id":"msg_cut","type":"message","role":"assistant","model":"m","content":[{"type":"tool_use","id":"toolu_1","name":"rec","input":{}}],"stop_reason":"max_tokens","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
)

// scriptedMessagesServer answers POST /v1/messages with the scripted message
// bodies in order (as SSE when the request sets "stream") and records every
// request body. A scripted error body is answered with status 529. A request
// past the end of the script fails the test.
func scriptedMessagesServer(t *testing.T, script ...string) (*httptest.Server, func() [][]byte) {
	t.Helper()
	var (
		mu     sync.Mutex
		bodies [][]byte
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading request body: %v", err)
			http.Error(w, "bad body", http.StatusBadRequest)
			return
		}
		mu.Lock()
		bodies = append(bodies, body)
		n := len(bodies)
		mu.Unlock()
		if n > len(script) {
			t.Errorf("unexpected request %d: only %d responses scripted", n, len(script))
			http.Error(w, "unscripted", http.StatusInternalServerError)
			return
		}
		if gjson.Get(script[n-1], "type").String() == "error" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(529)
			_, _ = w.Write([]byte(script[n-1]))
			return
		}
		if gjson.GetBytes(body, "stream").Bool() {
			w.Header().Set("Content-Type", "text/event-stream")
			writeMessageEvents(w, script[n-1])
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(script[n-1]))
	}))
	t.Cleanup(server.Close)
	return server, func() [][]byte {
		mu.Lock()
		defer mu.Unlock()
		return append([][]byte(nil), bodies...)
	}
}

// writeMessageEvents replays a complete message body as a stream, one
// content_block_start per block carrying the whole block.
func writeMessageEvents(w io.Writer, message string) {
	start, _ := sjson.SetRaw(message, "content", "[]")
	start, _ = sjson.SetRaw(start, "stop_reason", "null")
	fmt.Fprintf(w, "event: message_start\ndata: {\"type\":\"message_start\",\"message\":%s}\n\n", start)
	for i, block := range gjson.Get(message, "content").Array() {
		fmt.Fprintf(w, "event: content_block_start\ndata: {\"type\":\"content_block_start\",\"index\":%d,\"content_block\":%s}\n\n", i, block.Raw)
		fmt.Fprintf(w, "event: content_block_stop\ndata: {\"type\":\"content_block_stop\",\"index\":%d}\n\n", i)
	}
	fmt.Fprintf(w, "event: message_delta\ndata: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":%s,\"stop_sequence\":null},\"usage\":%s}\n\n",
		gjson.Get(message, "stop_reason").Raw, gjson.Get(message, "usage").Raw)
	fmt.Fprint(w, "event: message_stop\ndata: {\"type\":\"message_stop\"}\n\n")
}

func requireJSONEqual(t *testing.T, got, want string) {
	t.Helper()
	var g, w any
	if err := json.Unmarshal([]byte(got), &g); err != nil {
		t.Fatalf("unmarshal got: %v\n%s", err, got)
	}
	if err := json.Unmarshal([]byte(want), &w); err != nil {
		t.Fatalf("unmarshal want: %v\n%s", err, want)
	}
	if !reflect.DeepEqual(g, w) {
		t.Fatalf("JSON mismatch\n got: %s\nwant: %s", got, want)
	}
}

func pauseTurnParams(maxIterations int) BetaToolRunnerParams {
	return BetaToolRunnerParams{BetaMessageNewParams: BetaMessageNewParams{
		Model:     "m",
		MaxTokens: 512,
		Messages:  []BetaMessageParam{NewBetaUserMessage(NewBetaTextBlock("What's the weather in SF?"))},
	}, MaxIterations: maxIterations}
}

// A pause_turn response has no client tool calls to answer but is not
// finished: the runner must send the paused turn back so the server resumes it.
func TestBetaToolRunner_PauseTurnResumes(t *testing.T) {
	requireTurnResumed(t, pausedTurnJSON)
}

// A compaction-stopped response is handed back before the model answers: the
// runner must send the compacted turn back unchanged, with no tool_result
// turn, so the model continues from the summary.
func TestBetaToolRunner_CompactionResumes(t *testing.T) {
	requireTurnResumed(t, compactedTurnJSON)
}

// requireTurnResumed scripts unfinishedTurn followed by a normal end_turn and
// checks, for both runners, that exactly two requests are made and the second
// ends with the unfinished assistant turn echoed verbatim.
func requireTurnResumed(t *testing.T, unfinishedTurn string) {
	t.Helper()
	for _, stream := range []bool{false, true} {
		t.Run(fmt.Sprintf("stream=%v", stream), func(t *testing.T) {
			server, requests := scriptedMessagesServer(t, unfinishedTurn, endTurnJSON)
			client := newTestToolRunnerClient(server)
			ctx := context.Background()

			var final *BetaMessage
			var history []BetaMessageParam
			if stream {
				runner := client.Beta.Messages.NewToolRunnerStreaming(nil, pauseTurnParams(5))
				for events, err := range runner.AllStreaming(ctx) {
					if err != nil {
						t.Fatalf("AllStreaming: %v", err)
					}
					for _, err := range events {
						if err != nil {
							t.Fatalf("streaming turn: %v", err)
						}
					}
				}
				final, history = runner.LastMessage(), runner.Messages()
			} else {
				runner := client.Beta.Messages.NewToolRunner(nil, pauseTurnParams(5))
				var err error
				if final, err = runner.RunToCompletion(ctx); err != nil {
					t.Fatalf("RunToCompletion: %v", err)
				}
				history = runner.Messages()
			}

			if final == nil {
				t.Fatal("expected a final message, got nil")
			}
			if final.ID != "msg_end" || final.StopReason != BetaStopReasonEndTurn {
				t.Fatalf("final message = %s (%s), want msg_end (end_turn)", final.ID, final.StopReason)
			}
			bodies := requests()
			if len(bodies) != 2 {
				t.Fatalf("expected 2 requests, got %d", len(bodies))
			}
			sent := gjson.GetBytes(bodies[1], "messages").Array()
			if len(sent) != 2 {
				t.Fatalf("second request should carry the user turn plus the unfinished turn, got %d messages", len(sent))
			}
			if role := sent[1].Get("role").String(); role != "assistant" {
				t.Fatalf("last message role = %q, want assistant", role)
			}
			requireJSONEqual(t, sent[1].Get("content").Raw, gjson.Get(unfinishedTurn, "content").Raw)
			if len(history) != 3 {
				t.Fatalf("expected user + unfinished + resumed turns in history, got %d messages", len(history))
			}
		})
	}
}

// A max_tokens turn is final even when it carries a client tool call: the
// runner makes no follow-up request and never invokes the tool.
func TestBetaToolRunner_MaxTokensToolUseStops(t *testing.T) {
	server, requests := scriptedMessagesServer(t, cutOffToolUseJSON)
	client := newTestToolRunnerClient(server)
	tool := &recordingTool{name: "rec"}
	runner := client.Beta.Messages.NewToolRunner([]BetaTool{tool}, pauseTurnParams(5))

	final, err := runner.RunToCompletion(context.Background())
	if err != nil {
		t.Fatalf("RunToCompletion: %v", err)
	}
	if final == nil {
		t.Fatal("expected a final message, got nil")
	}
	if final.StopReason != BetaStopReasonMaxTokens {
		t.Fatalf("final stop_reason = %s, want max_tokens", final.StopReason)
	}
	if got := len(requests()); got != 1 {
		t.Fatalf("expected 1 request, got %d", got)
	}
	if tool.calls != 0 {
		t.Fatalf("cut-off tool call must not execute, ran %d times", tool.calls)
	}
	if got := len(runner.Messages()); got != 2 {
		t.Fatalf("expected user + cut-off turns in history with no tool_result, got %d messages", got)
	}
}

// A server that keeps pausing is still bounded by MaxIterations.
func TestBetaToolRunner_PauseTurnStopsAtMaxIterations(t *testing.T) {
	server, requests := scriptedMessagesServer(t, pausedTurnJSON, pausedTurnJSON, pausedTurnJSON)
	client := newTestToolRunnerClient(server)
	runner := client.Beta.Messages.NewToolRunner(nil, pauseTurnParams(3))

	final, err := runner.RunToCompletion(context.Background())
	if err != nil {
		t.Fatalf("RunToCompletion: %v", err)
	}
	if final == nil {
		t.Fatal("expected a final message, got nil")
	}
	if final.StopReason != BetaStopReasonPauseTurn {
		t.Fatalf("final stop_reason = %s, want pause_turn", final.StopReason)
	}
	if got := len(requests()); got != 3 {
		t.Fatalf("expected 3 requests, got %d", got)
	}
}

const (
	toolUseTurnJSON = `{"id":"msg_tool","type":"message","role":"assistant","model":"m","content":[{"type":"tool_use","id":"toolu_1","name":"weather","input":{"city":"SF"}}],"stop_reason":"tool_use","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
	// compactionResponseJSON is what a request carrying `compaction` returns:
	// the signed summary block and nothing sampled after it.
	compactionResponseJSON = `{"id":"msg_compaction","type":"message","role":"assistant","model":"m","content":[{"type":"compaction","content":"Summary so far.","signature":"c2ln"}],"stop_reason":"compaction","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
	// compactionWithListingJSON carries a block type this SDK does not model
	// after the compaction block.
	compactionWithListingJSON = `{"id":"msg_compaction","type":"message","role":"assistant","model":"m","content":[{"type":"compaction","content":"Summary so far.","signature":"c2ln"},{"type":"some_future_listing","server":"docs","tools":[{"name":"search"}]}],"stop_reason":"compaction","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
	// A compaction request that produced no summary answers with a block whose
	// content is null, or with no content and the summarization call's own stop
	// reason.
	nullSummaryResponseJSON = `{"id":"msg_compaction","type":"message","role":"assistant","model":"m","content":[{"type":"compaction","content":null,"signature":null}],"stop_reason":"compaction","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
	noContentResponseJSON   = `{"id":"msg_compaction","type":"message","role":"assistant","model":"m","content":[],"stop_reason":"max_tokens","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
	// fallbackEndTurnJSON is a finished turn whose first attempt made a tool call
	// and then refused; the serving model wrote what follows the fallback block.
	fallbackEndTurnJSON = `{"id":"msg_fallback","type":"message","role":"assistant","model":"m","content":[{"type":"tool_use","id":"toolu_1","name":"weather","input":{}},{"type":"fallback","from":{"model":"m"},"to":{"model":"n"},"trigger":{"type":"refusal","category":"general_harms"}},{"type":"text","text":"done"}],"stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
	overloadedErrorJSON = `{"type":"error","error":{"type":"overloaded_error","message":"Overloaded"}}`
	refusedToolUseJSON  = `{"id":"msg_refused","type":"message","role":"assistant","model":"m","content":[{"type":"tool_use","id":"toolu_1","name":"weather","input":{}}],"stop_reason":"refusal","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
)

// turnRunner drives either runner one request at a time, so each compaction
// case runs against both.
type turnRunner struct {
	*betaToolRunnerBase
	next func(ctx context.Context) (*BetaMessage, error)
}

func newTurnRunner(client Client, stream bool, tools []BetaTool, params BetaToolRunnerParams) *turnRunner {
	if !stream {
		runner := client.Beta.Messages.NewToolRunner(tools, params)
		return &turnRunner{&runner.betaToolRunnerBase, runner.NextMessage}
	}
	runner := client.Beta.Messages.NewToolRunnerStreaming(tools, params)
	return &turnRunner{&runner.betaToolRunnerBase, func(ctx context.Context) (*BetaMessage, error) {
		var message BetaMessage
		events := 0
		for event, err := range runner.NextStreaming(ctx) {
			if err != nil {
				return nil, err
			}
			if err := message.Accumulate(event); err != nil {
				return nil, err
			}
			events++
		}
		if events == 0 {
			return nil, nil
		}
		return &message, nil
	}}
}

// run calls onMessage with each message until the run is over and returns the
// messages' ids.
func (r *turnRunner) run(t *testing.T, onMessage func(*BetaMessage)) []string {
	t.Helper()
	var ids []string
	for {
		message, err := r.next(context.Background())
		if err != nil {
			t.Fatalf("next turn: %v", err)
		}
		if message == nil {
			return ids
		}
		ids = append(ids, message.ID)
		if onMessage != nil {
			onMessage(message)
		}
	}
}

func compactOnToolUse(runner *turnRunner) func(*BetaMessage) {
	return func(message *BetaMessage) {
		if message.StopReason == BetaStopReasonToolUse {
			runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
		}
	}
}

func forEachRunner(t *testing.T, test func(t *testing.T, stream bool)) {
	t.Helper()
	for _, stream := range []bool{false, true} {
		t.Run(fmt.Sprintf("stream=%v", stream), func(t *testing.T) { test(t, stream) })
	}
}

// captureStderr redirects os.Stderr, where the SDK prints its warnings, and
// returns a function that restores it and reports what was written.
func captureStderr(t *testing.T) func() string {
	t.Helper()
	stderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stderr = w
	t.Cleanup(func() {
		os.Stderr = stderr
		_ = w.Close()
		_ = r.Close()
	})
	return func() string {
		os.Stderr = stderr
		_ = w.Close()
		written, _ := io.ReadAll(r)
		return string(written)
	}
}

func messagesJSON(t *testing.T, messages []BetaMessageParam) string {
	t.Helper()
	data, err := json.Marshal(messages)
	if err != nil {
		t.Fatalf("marshal messages: %v", err)
	}
	return string(data)
}

func requireRequestCount(t *testing.T, bodies [][]byte, want int) {
	t.Helper()
	if len(bodies) != want {
		t.Fatalf("expected %d requests, got %d", want, len(bodies))
	}
}

// requireCompactionAlone checks a request's messages are the compaction
// response as it came and nothing else.
func requireCompactionAlone(t *testing.T, body []byte, response string) {
	t.Helper()
	requireJSONEqual(t, gjson.GetBytes(body, "messages").Raw, `[{"role":"assistant","content":`+gjson.Get(response, "content").Raw+`}]`)
}

// A compaction asked for on a tool-use turn goes out as its own request once
// that turn's tool results are in, without context_management and without
// counting towards MaxIterations; the next request carries the response alone,
// as it came, including a block this SDK does not model.
func TestBetaToolRunner_CompactBeforeNextTurn_SentAfterToolResults(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		server, requests := scriptedMessagesServer(t, toolUseTurnJSON, compactionWithListingJSON, endTurnJSON)
		var betas []string
		client := NewClient(
			option.WithBaseURL(server.URL),
			option.WithAPIKey("test-key"),
			option.WithMaxRetries(0),
			option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
				betas = append(betas, req.Header.Get("anthropic-beta"))
				return next(req)
			}),
		)
		params := pauseTurnParams(2)
		params.Betas = []AnthropicBeta{AnthropicBetaCompact2026_09_04}
		params.ContextManagement = BetaContextManagementConfigParam{Edits: []BetaContextManagementConfigEditUnionParam{
			{OfClearToolUses20250919: &BetaClearToolUses20250919EditParam{}},
		}}
		runner := newTurnRunner(client, stream, []BetaTool{&stubBetaTool{name: "weather"}}, params)

		var stopReasons []BetaStopReason
		ids := runner.run(t, func(message *BetaMessage) {
			stopReasons = append(stopReasons, message.StopReason)
			compactOnToolUse(runner)(message)
		})
		requireYieldedOnce(t, ids, "msg_tool", "msg_compaction", "msg_end")
		if fmt.Sprint(stopReasons) != "[tool_use compaction end_turn]" {
			t.Fatalf("stop reasons = %v", stopReasons)
		}
		if got := runner.IterationCount(); got != 2 {
			t.Fatalf("IterationCount = %d, want 2: the compaction request is not an iteration", got)
		}

		bodies := requests()
		requireRequestCount(t, bodies, 3)
		first, compaction, after := bodies[0], bodies[1], bodies[2]
		requireJSONEqual(t, gjson.GetBytes(compaction, "compaction").Raw, `{"type":"summarize"}`)
		if gjson.GetBytes(compaction, "context_management").Exists() {
			t.Fatalf("compaction request must not carry context_management: %s", compaction)
		}
		requireJSONEqual(t, gjson.GetBytes(compaction, "messages").Raw, `[
			{"role":"user","content":[{"type":"text","text":"What's the weather in SF?"}]},
			{"role":"assistant","content":[{"type":"tool_use","id":"toolu_1","name":"weather","input":{"city":"SF"}}]},
			{"role":"user","content":[{"type":"tool_result","tool_use_id":"toolu_1","content":[{"type":"text","text":"ok from weather"}]}]}
		]`)

		requireCompactionAlone(t, after, compactionWithListingJSON)
		if got := gjson.GetBytes(after, "messages.0.content.1.type").String(); got != "some_future_listing" {
			t.Fatalf("second block type = %q, want the unmodeled block kept in place", got)
		}
		if gjson.GetBytes(after, "compaction").Exists() {
			t.Fatalf("request after the compaction must not carry compaction: %s", after)
		}
		requireJSONEqual(t, gjson.GetBytes(after, "context_management").Raw, gjson.GetBytes(first, "context_management").Raw)
		if fmt.Sprint(betas) != "[compact-2026-09-04 compact-2026-09-04 compact-2026-09-04]" {
			t.Fatalf("anthropic-beta headers = %v, want the caller's beta and nothing added", betas)
		}
	})
}

func TestBetaToolRunner_CompactBeforeNextTurn_BeforeFirstIterationIsFirstRequest(t *testing.T) {
	server, requests := scriptedMessagesServer(t, compactionResponseJSON, endTurnJSON)
	runner := newTurnRunner(newTestToolRunnerClient(server), false, nil, pauseTurnParams(1))
	runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})

	requireYieldedOnce(t, runner.run(t, nil), "msg_compaction", "msg_end")
	bodies := requests()
	requireRequestCount(t, bodies, 2)
	requireJSONEqual(t, gjson.GetBytes(bodies[0], "compaction").Raw, `{"type":"summarize"}`)
	if got := gjson.GetBytes(bodies[0], "messages.#").Int(); got != 1 {
		t.Fatalf("compaction request carried %d messages, want the saved history", got)
	}
	requireCompactionAlone(t, bodies[1], compactionResponseJSON)
}

// Calling again before the compaction is sent replaces the earlier call: one
// compaction request goes out, with the last config, as given.
func TestBetaToolRunner_CompactBeforeNextTurn_LastConfigWinsAndIsSentAsGiven(t *testing.T) {
	server, requests := scriptedMessagesServer(t, toolUseTurnJSON, compactionResponseJSON, endTurnJSON)
	runner := newTurnRunner(newTestToolRunnerClient(server), false, []BetaTool{&stubBetaTool{name: "weather"}}, pauseTurnParams(0))

	runner.run(t, func(message *BetaMessage) {
		if message.StopReason != BetaStopReasonToolUse {
			return
		}
		runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
		runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{
			OfSummarize: &BetaSummarizeCompactionParam{Instructions: String("  ")},
		})
	})

	bodies := requests()
	requireRequestCount(t, bodies, 3)
	requireJSONEqual(t, gjson.GetBytes(bodies[1], "compaction").Raw, `{"type":"summarize","instructions":"  "}`)
	if gjson.GetBytes(bodies[2], "compaction").Exists() {
		t.Fatalf("only one compaction request should go out: %s", bodies[2])
	}
}

// A paused turn is resumed and finished before the compaction goes out.
func TestBetaToolRunner_CompactBeforeNextTurn_WaitsOutPausedTurn(t *testing.T) {
	server, requests := scriptedMessagesServer(t, pausedTurnJSON, toolUseTurnJSON, compactionResponseJSON, endTurnJSON)
	runner := newTurnRunner(newTestToolRunnerClient(server), false, []BetaTool{&stubBetaTool{name: "weather"}}, pauseTurnParams(0))

	ids := runner.run(t, func(message *BetaMessage) {
		if message.StopReason == BetaStopReasonPauseTurn {
			runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
		}
	})
	requireYieldedOnce(t, ids, "msg_paused", "msg_tool", "msg_compaction", "msg_end")
	bodies := requests()
	requireRequestCount(t, bodies, 4)
	if gjson.GetBytes(bodies[1], "compaction").Exists() {
		t.Fatalf("the paused turn must be resumed before compacting: %s", bodies[1])
	}
	requireJSONEqual(t, gjson.GetBytes(bodies[2], "compaction").Raw, `{"type":"summarize"}`)
	if got := gjson.GetBytes(bodies[2], "messages.#").Int(); got != 4 {
		t.Fatalf("compaction request carried %d messages, want user, paused turn, resumed turn, tool results", got)
	}
}

// A compaction asked for on the final turn still goes out, even when that
// turn was also the last iteration allowed; the run then stops.
func TestBetaToolRunner_CompactBeforeNextTurn_FinalTurnCompactsThenStops(t *testing.T) {
	for _, finalTurn := range []struct{ name, json string }{
		{"end_turn", endTurnJSON},
		{"a tool call only before the fallback block", fallbackEndTurnJSON},
	} {
		t.Run(finalTurn.name, func(t *testing.T) {
			forEachRunner(t, func(t *testing.T, stream bool) {
				server, requests := scriptedMessagesServer(t, finalTurn.json, compactionResponseJSON)
				runner := newTurnRunner(newTestToolRunnerClient(server), stream, nil, pauseTurnParams(1))

				ids := runner.run(t, func(message *BetaMessage) {
					runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
				})
				requireYieldedOnce(t, ids, gjson.Get(finalTurn.json, "id").String(), "msg_compaction")
				if !runner.IsCompleted() {
					t.Fatal("expected the run to be completed")
				}
				bodies := requests()
				requireRequestCount(t, bodies, 2)
				requireJSONEqual(t, gjson.GetBytes(bodies[1], "compaction").Raw, `{"type":"summarize"}`)
				if got := gjson.GetBytes(bodies[1], "messages.#").Int(); got != 2 {
					t.Fatalf("compaction request carried %d messages, want the user turn and the final answer", got)
				}
				requireJSONEqual(t, messagesJSON(t, runner.Messages()), `[{"role":"assistant","content":`+gjson.Get(compactionResponseJSON, "content").Raw+`}]`)
			})
		})
	}
}

// A final turn that was cut off with a tool call nobody answered cannot be
// compacted: the pending compaction is dropped with a warning.
func TestBetaToolRunner_CompactBeforeNextTurn_SkippedAfterUnansweredToolCall(t *testing.T) {
	for _, finalTurn := range []struct{ name, json string }{
		{"max_tokens", cutOffToolUseJSON},
		{"refusal", refusedToolUseJSON},
	} {
		t.Run(finalTurn.name, func(t *testing.T) {
			forEachRunner(t, func(t *testing.T, stream bool) {
				warned := captureStderr(t)
				server, requests := scriptedMessagesServer(t, finalTurn.json)
				runner := newTurnRunner(newTestToolRunnerClient(server), stream, []BetaTool{&recordingTool{name: "rec"}}, pauseTurnParams(0))

				ids := runner.run(t, func(message *BetaMessage) {
					runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
				})
				requireYieldedOnce(t, ids, gjson.Get(finalTurn.json, "id").String())
				requireRequestCount(t, requests(), 1)
				if got := warned(); !strings.Contains(got, "skipped the pending compaction") || !strings.Contains(got, finalTurn.name) {
					t.Fatalf("expected a warning naming the stop reason, got %q", got)
				}
			})
		})
	}
}

// A run cut short by MaxIterations after a tool turn drops the pending
// compaction without a word, as it drops the tool calls.
func TestBetaToolRunner_CompactBeforeNextTurn_DroppedAtMaxIterations(t *testing.T) {
	warned := captureStderr(t)
	server, requests := scriptedMessagesServer(t, toolUseTurnJSON)
	runner := newTurnRunner(newTestToolRunnerClient(server), false, []BetaTool{&stubBetaTool{name: "weather"}}, pauseTurnParams(1))

	requireYieldedOnce(t, runner.run(t, compactOnToolUse(runner)), "msg_tool")
	requireRequestCount(t, requests(), 1)
	if got := warned(); got != "" {
		t.Fatalf("expected no warning, got %q", got)
	}
}

// With no summary the history is kept, a warning is printed and the run goes
// on; the response that came back does not become the last message.
func TestBetaToolRunner_CompactBeforeNextTurn_NoSummaryKeepsHistory(t *testing.T) {
	for _, tt := range []struct{ name, response string }{
		{"a compaction block with null content", nullSummaryResponseJSON},
		{"no content at all", noContentResponseJSON},
	} {
		t.Run(tt.name, func(t *testing.T) {
			forEachRunner(t, func(t *testing.T, stream bool) {
				warned := captureStderr(t)
				server, requests := scriptedMessagesServer(t, toolUseTurnJSON, tt.response, endTurnJSON)
				runner := newTurnRunner(newTestToolRunnerClient(server), stream, []BetaTool{&stubBetaTool{name: "weather"}}, pauseTurnParams(0))

				ids := runner.run(t, func(message *BetaMessage) {
					compactOnToolUse(runner)(message)
					if message.ID == "msg_compaction" && runner.LastMessage().ID != "msg_tool" {
						t.Errorf("LastMessage = %s, want the turn before the failed compaction", runner.LastMessage().ID)
					}
				})
				requireYieldedOnce(t, ids, "msg_tool", "msg_compaction", "msg_end")
				bodies := requests()
				requireRequestCount(t, bodies, 3)
				requireJSONEqual(t, gjson.GetBytes(bodies[2], "messages").Raw, gjson.GetBytes(bodies[1], "messages").Raw)
				if got := warned(); !strings.Contains(got, "no summary") {
					t.Fatalf("expected a warning about the missing summary, got %q", got)
				}
			})
		})
	}
}

// After a final-turn compaction that produced no summary, the run's final
// message is still the model's answer.
func TestBetaToolRunner_CompactBeforeNextTurn_NoSummaryOnFinalTurnKeepsFinalMessage(t *testing.T) {
	captureStderr(t)
	server, requests := scriptedMessagesServer(t, endTurnJSON, noContentResponseJSON)
	client := newTestToolRunnerClient(server)
	runner := client.Beta.Messages.NewToolRunner(nil, pauseTurnParams(0))
	ctx := context.Background()

	if _, err := runner.NextMessage(ctx); err != nil {
		t.Fatalf("final turn: %v", err)
	}
	runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
	final, err := runner.RunToCompletion(ctx)
	if err != nil {
		t.Fatalf("RunToCompletion: %v", err)
	}
	if final == nil || final.ID != "msg_end" {
		t.Fatalf("final message = %v, want msg_end", final)
	}
	requireRequestCount(t, requests(), 2)
}

// A call made while the compaction response is being handled is ignored, so a
// token threshold that also matches that response does not compact twice.
func TestBetaToolRunner_CompactBeforeNextTurn_IgnoredOnCompactionResponse(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		server, requests := scriptedMessagesServer(t, toolUseTurnJSON, compactionResponseJSON, endTurnJSON)
		runner := newTurnRunner(newTestToolRunnerClient(server), stream, []BetaTool{&stubBetaTool{name: "weather"}}, pauseTurnParams(0))

		ids := runner.run(t, func(message *BetaMessage) {
			if message.StopReason != BetaStopReasonEndTurn {
				runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
			}
		})
		requireYieldedOnce(t, ids, "msg_tool", "msg_compaction", "msg_end")
		bodies := requests()
		requireRequestCount(t, bodies, 3)
		if gjson.GetBytes(bodies[2], "compaction").Exists() {
			t.Fatalf("the call made on the compaction response must be ignored: %s", bodies[2])
		}
	})
}

// A compaction request that fails is not retried: the error is returned, the
// next turn is an ordinary request, and a call made after the failure is sent
// as usual.
func TestBetaToolRunner_CompactBeforeNextTurn_FailedRequestIsNotRetried(t *testing.T) {
	for _, tt := range []struct {
		name      string
		askAgain  bool
		afterward []string
	}{
		{"the run goes on without it", false, []string{endTurnJSON}},
		{"a call made after the failure is sent", true, []string{compactionResponseJSON, endTurnJSON}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			forEachRunner(t, func(t *testing.T, stream bool) {
				script := append([]string{toolUseTurnJSON, overloadedErrorJSON}, tt.afterward...)
				server, requests := scriptedMessagesServer(t, script...)
				runner := newTurnRunner(newTestToolRunnerClient(server), stream, []BetaTool{&stubBetaTool{name: "weather"}}, pauseTurnParams(0))
				ctx := context.Background()

				if _, err := runner.next(ctx); err != nil {
					t.Fatalf("first turn: %v", err)
				}
				runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})
				if _, err := runner.next(ctx); err == nil || runner.Err() == nil {
					t.Fatalf("expected the failed compaction request to be reported, got %v (Err() = %v)", err, runner.Err())
				}
				if tt.askAgain {
					runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{
						OfSummarize: &BetaSummarizeCompactionParam{Instructions: String("Second try.")},
					})
				}

				ids := runner.run(t, nil)
				bodies := requests()
				requireRequestCount(t, bodies, len(script))
				if !tt.askAgain {
					requireYieldedOnce(t, ids, "msg_end")
					if gjson.GetBytes(bodies[2], "compaction").Exists() {
						t.Fatalf("the failed compaction must not be retried: %s", bodies[2])
					}
					return
				}
				requireYieldedOnce(t, ids, "msg_compaction", "msg_end")
				requireJSONEqual(t, gjson.GetBytes(bodies[2], "compaction").Raw, `{"type":"summarize","instructions":"Second try."}`)
				requireCompactionAlone(t, bodies[3], compactionResponseJSON)
			})
		})
	}
}

// Params is an exported field, so replacing the messages while the compaction
// response streams cannot be stopped where it happens; it is reported once the
// response has been read, and the compaction response does not replace them.
func TestBetaToolRunnerStreaming_CompactBeforeNextTurn_MessagesReplacedWhileCompacting(t *testing.T) {
	extra := NewBetaUserMessage(NewBetaTextBlock("One more thing."))
	for _, tt := range []struct {
		name    string
		change  func(runner *BetaToolRunnerStreaming)
		refused bool
	}{
		{"AppendMessages", func(runner *BetaToolRunnerStreaming) { runner.AppendMessages(extra) }, true},
		{"a new slice of the same length", func(runner *BetaToolRunnerStreaming) {
			runner.Params.Messages = append([]BetaMessageParam{}, runner.Params.Messages...)
		}, true},
		{"another param", func(runner *BetaToolRunnerStreaming) { runner.Params.MaxTokens = 1024 }, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := scriptedMessagesServer(t, compactionResponseJSON)
			client := newTestToolRunnerClient(server)
			runner := client.Beta.Messages.NewToolRunnerStreaming(nil, pauseTurnParams(0))
			runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})

			var got error
			changed := false
			for _, err := range runner.NextStreaming(context.Background()) {
				if err != nil {
					got = err
					break
				}
				if !changed {
					tt.change(runner)
					changed = true
				}
			}
			if !tt.refused {
				if got != nil {
					t.Fatalf("changing another param while compacting must be allowed: %v", got)
				}
				requireJSONEqual(t, messagesJSON(t, runner.Messages()), `[{"role":"assistant","content":`+gjson.Get(compactionResponseJSON, "content").Raw+`}]`)
				return
			}
			if got == nil || !strings.Contains(got.Error(), "while the conversation was being compacted") {
				t.Fatalf("expected the change to be refused, got %v", got)
			}
			if runner.Err() != got {
				t.Fatalf("Err() = %v, want %v", runner.Err(), got)
			}
			if role := runner.Messages()[0].Role; role != BetaMessageParamRoleUser {
				t.Fatalf("the caller's messages should be kept, first role = %q", role)
			}
		})
	}
}

// Set in the runner's own params, compaction would compact on every request:
// it is refused once and cleared, so the run can go on without it.
func TestBetaToolRunner_CompactionParamRefused(t *testing.T) {
	server, requests := scriptedMessagesServer(t, endTurnJSON)
	params := pauseTurnParams(0)
	params.Compaction = BetaCompactionConfigUnionParam{OfSummarize: &BetaSummarizeCompactionParam{}}
	runner := newTurnRunner(newTestToolRunnerClient(server), false, nil, params)

	_, err := runner.next(context.Background())
	if err == nil || !strings.Contains(err.Error(), "CompactBeforeNextTurn") {
		t.Fatalf("expected an error pointing at CompactBeforeNextTurn, got %v", err)
	}
	if runner.Err() != err {
		t.Fatalf("Err() = %v, want %v", runner.Err(), err)
	}
	requireRequestCount(t, requests(), 0)

	requireYieldedOnce(t, runner.run(t, nil), "msg_end")
	if gjson.GetBytes(requests()[0], "compaction").Exists() {
		t.Fatal("the refused param must not be sent")
	}
}

// The compaction request goes out without context_management, so a compaction
// edit there has to be refused by the runner: the API would only notice on the
// request after the paid compaction. The refused call is dropped.
func TestBetaToolRunner_CompactBeforeNextTurn_RefusedBesideCompactionEdit(t *testing.T) {
	server, requests := scriptedMessagesServer(t, endTurnJSON)
	params := pauseTurnParams(0)
	params.ContextManagement = BetaContextManagementConfigParam{Edits: []BetaContextManagementConfigEditUnionParam{
		{OfClearToolUses20250919: &BetaClearToolUses20250919EditParam{}},
		{OfCompact20260112: &BetaCompact20260112EditParam{}},
	}}
	runner := newTurnRunner(newTestToolRunnerClient(server), false, nil, params)
	runner.CompactBeforeNextTurn(BetaCompactionConfigUnionParam{})

	_, err := runner.next(context.Background())
	if err == nil || !strings.Contains(err.Error(), "compaction edit") {
		t.Fatalf("expected the call to be refused, got %v", err)
	}
	requireRequestCount(t, requests(), 0)

	requireYieldedOnce(t, runner.run(t, nil), "msg_end")
	if gjson.GetBytes(requests()[0], "compaction").Exists() {
		t.Fatal("the refused compaction must not be sent later")
	}
}
