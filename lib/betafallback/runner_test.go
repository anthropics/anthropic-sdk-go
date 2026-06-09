package betafallback_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/lib/betafallback"
	"github.com/anthropics/anthropic-sdk-go/toolrunner"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Composition tests: the tool runner drives the conversation loop while the
// fallback middleware handles refusals underneath it. The runner must never
// see a refusal — only the served (possibly spliced) message — and the turns
// it echoes back must be valid on the wire.

type echoInput struct {
	Text string `json:"text" jsonschema:"required,description=Text to echo"`
}

func echoTool(t *testing.T) anthropic.BetaTool {
	tool, err := toolrunner.NewBetaToolFromJSONSchema("echo", "Echoes text",
		func(ctx context.Context, in echoInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: in.Text},
			}, nil
		})
	require.NoError(t, err)
	return tool
}

func recordingEchoTool(t *testing.T, calls *[]string) anthropic.BetaTool {
	tool, err := toolrunner.NewBetaToolFromJSONSchema("echo", "Echoes text",
		func(ctx context.Context, in echoInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			*calls = append(*calls, in.Text)
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: in.Text},
			}, nil
		})
	require.NoError(t, err)
	return tool
}

// splicedResponse is a message as the middleware's stream splice produces it:
// the refused model's blocks, a fallback boundary, then the serving model's.
func splicedResponse(blocks string, stopReason string) string {
	return fmt.Sprintf(`{
		"id": "msg_s", "type": "message", "role": "assistant", "model": "fallback-model",
		"content": [%s], "stop_reason": %q, "stop_sequence": null,
		"usage": {"input_tokens": 1, "output_tokens": 1}
	}`, blocks, stopReason)
}

func toolUseResponse(model string) string {
	return fmt.Sprintf(`{
		"id": "msg_t", "type": "message", "role": "assistant", "model": %q,
		"content": [{"type": "tool_use", "id": "toolu_1", "name": "echo", "input": {"text": "hi"}}],
		"stop_reason": "tool_use", "stop_sequence": null,
		"usage": {"input_tokens": 1, "output_tokens": 1}
	}`, model)
}

func fixtureSSE(t *testing.T, name string) string {
	t.Helper()
	buf, err := os.ReadFile(filepath.Join("testdata", "fallbackstream", name))
	require.NoError(t, err)
	return string(buf)
}

func TestToolRunnerWithFallbackMiddlewareMidLoopRefusal(t *testing.T) {
	// The runner drives the loop; the middleware handles a refusal on the
	// tool-result turn invisibly, and the pin carries on the shared state.
	client, transport := fallbackTestClient(t,
		[]string{
			toolUseResponse("primary-model"),
			refusalResponse("primary-model", "credit-token"),
			messageResponse("fallback-model"),
		},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	state := &betafallback.BetaFallbackState{}
	runner := client.Beta.Messages.NewToolRunner(
		[]anthropic.BetaTool{echoTool(t)},
		anthropic.BetaToolRunnerParams{BetaMessageNewParams: fallbackTestParams},
		betafallback.WithBetaFallbackState(state),
	)
	final, err := runner.RunToCompletion(context.Background())
	require.NoError(t, err)
	require.NotNil(t, final)
	assert.Equal(t, anthropic.Model("fallback-model"), final.Model)

	// Three wire requests: opening turn, tool-result turn (refused), retry.
	require.Len(t, transport.bodies, 3)
	refusedTurn, _ := json.Marshal(transport.bodies[1]["messages"])
	retryTurn, _ := json.Marshal(transport.bodies[2]["messages"])
	assert.JSONEq(t, string(refusedTurn), string(retryTurn), "the retry resends the refused turn's messages — tool results included")
	assert.Equal(t, "credit-token", transport.bodies[2]["fallback_credit_token"], "the retry redeems the refusal's token")
	assert.Equal(t, "fallback-model", transport.bodies[2]["model"])
	assert.Equal(t, 0, state.Index(), "the pin carries across runner turns")
}

func TestStreamingToolRunnerEchoesTheSplicedTurn(t *testing.T) {
	// A spliced turn (refused partial + boundary + fallback's tool_use) is
	// echoed by the runner's next request and survives the wire intact.
	transport := &sseTransport{responses: []string{
		fixtureSSE(t, "refusal.sse"),
		fixtureSSE(t, "fallback-tooluse.sse"),
		fixtureSSE(t, "fallback-end.sse"),
	}}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})

	runner := client.Beta.Messages.NewToolRunnerStreaming(
		[]anthropic.BetaTool{echoTool(t)},
		anthropic.BetaToolRunnerParams{BetaMessageNewParams: fallbackTestParams},
	)
	for events, err := range runner.AllStreaming(context.Background()) {
		require.NoError(t, err)
		for _, eventErr := range events {
			require.NoError(t, eventErr)
		}
	}
	require.NoError(t, runner.Err())

	// Three wire requests: refused opening turn, the middleware's hop, the
	// runner's tool-result turn echoing the spliced assistant message.
	require.Len(t, transport.bodies, 3)
	assert.Equal(t, "credit-token-a", transport.bodies[1]["fallback_credit_token"], "the hop redeems the token")
	echoed, err := json.Marshal(transport.bodies[2]["messages"])
	require.NoError(t, err)
	echo := string(echoed)
	// The echoed history is trimmed: the fallback block and the refused
	// model's thinking are removed; the tool call that got a result stays.
	assert.NotContains(t, echo, `"type":"fallback"`, "the fallback block is removed from the echo")
	assert.NotContains(t, echo, `"thinking":"Simple educational question. Benign."`, "the refused model's thinking is removed from the echo")
	assert.Contains(t, echo, `"type":"tool_use"`, "the fallback's tool call is echoed")
	assert.Contains(t, echo, `"type":"tool_result"`)
}

func TestNextTurnEchoStripsARefusalCutToolUse(t *testing.T) {
	// A refusal that fires after a completed tool_use leaves it orphaned in
	// the spliced turn — stop_reason is not tool_use, so no result follows.
	// ToParam keeps the orphan, but echoing the seam-bearing turn back must
	// strip it before the wire (conformance: multiturn-echo-orphan-tooluse).
	cut := event("message_start", `{"type":"message_start","message":{"type":"message","id":"msg_primary","role":"assistant","content":[],"model":"primary-model","stop_reason":null,"stop_sequence":null,"stop_details":null,"usage":{"input_tokens":9,"output_tokens":1,"cache_creation_input_tokens":0,"cache_read_input_tokens":0}}}`) +
		event("content_block_start", `{"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":"toolu_orphan","name":"echo","input":{"text":"hi"}}}`) +
		event("content_block_stop", `{"type":"content_block_stop","index":0}`) +
		event("message_delta", fmt.Sprintf(`{"type":"message_delta","delta":{"stop_reason":"refusal","stop_sequence":null,"stop_details":{"type":"refusal","category":null,"explanation":null,%s}},"usage":{"input_tokens":9,"output_tokens":12,"cache_read_input_tokens":0,"cache_creation_input_tokens":0}}`, tokenNoClaim)) +
		event("message_stop", `{"type":"message_stop"}`)
	transport := &sseTransport{responses: []string{
		cut, servedStream("fallback-model"), servedStream("fallback-model"),
	}}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})

	msg, _, _ := collectStream(t, client, context.Background(), fallbackTestParams)
	assert.Equal(t, anthropic.BetaStopReasonEndTurn, msg.StopReason)
	spliced, err := json.Marshal(msg.ToParam())
	require.NoError(t, err)
	assert.Contains(t, string(spliced), `"id":"toolu_orphan"`, "the spliced turn itself still carries the cut tool_use")

	// Echo the turn into the next request: the orphan must not reach the wire.
	params := fallbackTestParams
	params.Messages = append(append([]anthropic.BetaMessageParam{}, params.Messages...),
		msg.ToParam(), anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("continue")))
	_, _, _ = collectStream(t, client, context.Background(), params)

	require.Len(t, transport.bodies, 3)
	next, err := json.Marshal(transport.bodies[2]["messages"])
	require.NoError(t, err)
	assert.NotContains(t, string(next), `"type":"fallback"`, "the fallback block is removed from the echo")
	assert.NotContains(t, string(next), "toolu_orphan", "the orphaned tool_use is stripped from the echoed turn")
}

func TestToolRunnerSkipsTheRefusedAttemptsPreSeamToolUse(t *testing.T) {
	// A refusal cut after a complete tool_use leaves it before the seam in
	// the spliced turn. The runner must not execute it: the middleware
	// strips pre-seam tool calls from the echoed history, so answering one
	// would put an orphaned tool_result on the wire.
	client, transport := fallbackTestClient(t,
		[]string{
			splicedResponse(`
				{"type": "tool_use", "id": "toolu_orphan", "name": "echo", "input": {"text": "pre"}},
				{"type": "fallback", "from": {"model": "primary-model"}, "to": {"model": "fallback-model"}},
				{"type": "text", "text": "done"}`, "end_turn"),
			messageResponse("fallback-model"),
		},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	var calls []string
	runner := client.Beta.Messages.NewToolRunner(
		[]anthropic.BetaTool{recordingEchoTool(t, &calls)},
		anthropic.BetaToolRunnerParams{BetaMessageNewParams: fallbackTestParams},
	)
	final, err := runner.RunToCompletion(context.Background())
	require.NoError(t, err)
	require.NotNil(t, final)

	assert.Empty(t, calls, "the pre-seam tool call is not executed")
	assert.Len(t, transport.bodies, 1, "no tool-result turn follows a turn whose only tool calls are pre-seam")
}

func TestToolRunnerAnswersThePostSeamToolUse(t *testing.T) {
	// Control: the serving model's own tool call — after the seam — is
	// executed and answered as usual; the pre-seam orphan never reaches
	// the wire in any form.
	client, transport := fallbackTestClient(t,
		[]string{
			splicedResponse(`
				{"type": "tool_use", "id": "toolu_orphan", "name": "echo", "input": {"text": "pre"}},
				{"type": "fallback", "from": {"model": "primary-model"}, "to": {"model": "fallback-model"}},
				{"type": "tool_use", "id": "toolu_live", "name": "echo", "input": {"text": "post"}}`, "tool_use"),
			messageResponse("fallback-model"),
		},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	var calls []string
	runner := client.Beta.Messages.NewToolRunner(
		[]anthropic.BetaTool{recordingEchoTool(t, &calls)},
		anthropic.BetaToolRunnerParams{BetaMessageNewParams: fallbackTestParams},
	)
	_, err := runner.RunToCompletion(context.Background())
	require.NoError(t, err)

	assert.Equal(t, []string{"post"}, calls, "only the post-seam tool call executes")
	require.Len(t, transport.bodies, 2)
	echoed, err := json.Marshal(transport.bodies[1]["messages"])
	require.NoError(t, err)
	echo := string(echoed)
	assert.Contains(t, echo, `"tool_use_id":"toolu_live"`, "the post-seam tool call is answered")
	assert.NotContains(t, echo, "toolu_orphan", "nothing in the next turn references the pre-seam tool call")
}

// refusalToolUseResponse is a refusal that fired after the model emitted a
// complete tool_use — e.g. the last hop of an exhausted fallback chain.
func refusalToolUseResponse(model string) string {
	return fmt.Sprintf(`{
		"id": "msg_r", "type": "message", "role": "assistant", "model": %q,
		"content": [{"type": "tool_use", "id": "toolu_dead", "name": "echo", "input": {"text": "dead"}}],
		"stop_reason": "refusal", "stop_sequence": null,
		"stop_details": {"type": "refusal", "category": null, "explanation": null,
			"fallback_credit_token": "credit-token", "recommended_model": null},
		"usage": {"input_tokens": 1, "output_tokens": 1}
	}`, model)
}

func TestToolRunnerEndsTheLoopOnARefusalTerminatedTurn(t *testing.T) {
	// An exhausted fallback chain surfaces the last refusal as the final
	// message. Its tool calls belong to a dead conversation: the runner must
	// not execute them, and the loop ends with the refusal as the final
	// message instead of sending another turn.
	client, transport := fallbackTestClient(t,
		[]string{
			refusalToolUseResponse("primary-model"),
			refusalToolUseResponse("fallback-model"),
		},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	var calls []string
	runner := client.Beta.Messages.NewToolRunner(
		[]anthropic.BetaTool{recordingEchoTool(t, &calls)},
		anthropic.BetaToolRunnerParams{BetaMessageNewParams: fallbackTestParams},
	)
	final, err := runner.RunToCompletion(context.Background())
	require.NoError(t, err)
	require.NotNil(t, final)

	assert.Equal(t, anthropic.BetaStopReasonRefusal, final.StopReason, "the refusal is the final message")
	assert.Empty(t, calls, "tool calls on a refusal-terminated turn are not executed")
	assert.Len(t, transport.bodies, 2, "no wire request follows the surfaced refusal (both bodies are the middleware's own hops)")
}
