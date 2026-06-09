package betafallback_test

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/lib/betafallback"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sseTransport returns canned responses in order — SSE by default — and
// records request bodies and beta headers.
type sseTransport struct {
	t            *testing.T
	responses    []string
	statuses     []int    // parallel to responses; missing entries are 200
	contentTypes []string // parallel to responses; missing entries are text/event-stream
	bodies       []map[string]any
	betas        [][]string
}

func (s *sseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	buf, err := io.ReadAll(req.Body)
	require.NoError(s.t, err)
	var body map[string]any
	require.NoError(s.t, json.Unmarshal(buf, &body))
	s.bodies = append(s.bodies, body)
	s.betas = append(s.betas, req.Header.Values("anthropic-beta"))
	require.NotEmpty(s.t, s.responses, "more requests than scripted responses")
	next := s.responses[0]
	s.responses = s.responses[1:]
	status := http.StatusOK
	if len(s.statuses) > 0 {
		status = s.statuses[0]
		s.statuses = s.statuses[1:]
	}
	contentType := "text/event-stream"
	if len(s.contentTypes) > 0 {
		contentType = s.contentTypes[0]
		s.contentTypes = s.contentTypes[1:]
	}
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{contentType}},
		Body:       io.NopCloser(strings.NewReader(next)),
		Request:    req,
	}, nil
}

func event(name, data string) string {
	return fmt.Sprintf("event: %s\ndata: %s\n\n", name, data)
}

// refusalStream is a stream from model that refuses mid-text with the given
// stop_details fields.
func refusalStream(model, stopDetails string) string {
	return event("message_start", fmt.Sprintf(`{"type":"message_start","message":{"type":"message","id":"msg_primary","role":"assistant","content":[],"model":%q,"stop_reason":null,"stop_sequence":null,"stop_details":null,"usage":{"input_tokens":10,"output_tokens":1,"cache_creation_input_tokens":0,"cache_read_input_tokens":0}}}`, model)) +
		event("content_block_start", `{"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`) +
		event("ping", `{"type": "ping"}`) +
		event("content_block_delta", `{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"partial "}}`) +
		event("content_block_stop", `{"type":"content_block_stop","index":0}`) +
		event("message_delta", fmt.Sprintf(`{"type":"message_delta","delta":{"stop_reason":"refusal","stop_sequence":null,"stop_details":{"type":"refusal","category":null,"explanation":null,%s}},"usage":{"input_tokens":10,"output_tokens":5,"cache_read_input_tokens":0,"cache_creation_input_tokens":0}}`, stopDetails)) +
		event("message_stop", `{"type":"message_stop"}`)
}

// cutRefusalStream refuses with content block 0 still open.
func cutRefusalStream(model, stopDetails string) string {
	full := refusalStream(model, stopDetails)
	return strings.Replace(full, event("content_block_stop", `{"type":"content_block_stop","index":0}`), "", 1)
}

func servedStream(model string) string {
	return event("message_start", fmt.Sprintf(`{"type":"message_start","message":{"type":"message","id":"msg_retry","role":"assistant","content":[],"model":%q,"stop_reason":null,"stop_sequence":null,"stop_details":null,"usage":{"input_tokens":12,"output_tokens":1,"cache_creation_input_tokens":0,"cache_read_input_tokens":0}}}`, model)) +
		event("content_block_start", `{"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`) +
		event("content_block_delta", `{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"served"}}`) +
		event("content_block_stop", `{"type":"content_block_stop","index":0}`) +
		event("message_delta", `{"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null,"stop_details":null},"usage":{"input_tokens":12,"output_tokens":3,"cache_read_input_tokens":0,"cache_creation_input_tokens":0}}`) +
		event("message_stop", `{"type":"message_stop"}`)
}

const tokenWithClaim = `"fallback_credit_token":"credit-token-a","fallback_has_prefill_claim":true`
const tokenNoClaim = `"fallback_credit_token":"credit-token-a","fallback_has_prefill_claim":false`
const noToken = `"fallback_credit_token":null,"fallback_has_prefill_claim":null`

func streamingFallbackClient(t *testing.T, transport *sseTransport, fallbacks []anthropic.BetaFallbackParam) anthropic.Client {
	transport.t = t
	return anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{Transport: transport}),
		option.WithMaxRetries(0),
		option.WithMiddleware(betafallback.BetaRefusalFallbackMiddleware(fallbacks)),
	)
}

// collectStream drains a streaming call, returning the accumulated message,
// the (type, index) sequence, and each message_delta's raw JSON.
func collectStream(t *testing.T, client anthropic.Client, ctx context.Context, params anthropic.BetaMessageNewParams, opts ...option.RequestOption) (anthropic.BetaMessage, []string, []string) {
	stream := client.Beta.Messages.NewStreaming(ctx, params, opts...)
	defer stream.Close()
	var msg anthropic.BetaMessage
	var sequence []string
	var deltas []string
	for stream.Next() {
		event := stream.Current()
		require.NoError(t, msg.Accumulate(event))
		switch event.Type {
		case "content_block_start", "content_block_delta", "content_block_stop":
			sequence = append(sequence, fmt.Sprintf("%s:%d", event.Type, event.Index))
		default:
			sequence = append(sequence, string(event.Type))
		}
		if event.Type == "message_delta" {
			deltas = append(deltas, event.RawJSON())
		}
	}
	require.NoError(t, stream.Err())
	return msg, sequence, deltas
}

func deltaIterations(t *testing.T, rawDelta string) []struct {
	Type  string `json:"type"`
	Model string `json:"model"`
} {
	var delta struct {
		Usage struct {
			Iterations []struct {
				Type  string `json:"type"`
				Model string `json:"model"`
			} `json:"iterations"`
		} `json:"usage"`
	}
	require.NoError(t, json.Unmarshal([]byte(rawDelta), &delta))
	return delta.Usage.Iterations
}

func TestStreamingRefusalSplicesIntoOneMessage(t *testing.T) {
	transport := &sseTransport{responses: []string{
		refusalStream("primary-model", tokenWithClaim), servedStream("fallback-model"),
	}}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})

	msg, sequence, deltas := collectStream(t, client, context.Background(), fallbackTestParams)

	// One continuous message: A's block, the boundary, B's block reindexed.
	assert.Equal(t, []string{
		"message_start",
		"content_block_start:0", "content_block_delta:0", "content_block_stop:0",
		"content_block_start:1", "content_block_stop:1", // fallback boundary
		"content_block_start:2", "content_block_delta:2", "content_block_stop:2",
		"message_delta", "message_stop",
	}, sequence)

	assert.Equal(t, anthropic.BetaStopReasonEndTurn, msg.StopReason)
	require.Len(t, msg.Content, 3)
	boundary := msg.Content[1].AsFallback()
	assert.Equal(t, anthropic.Model("primary-model"), boundary.From.Model)
	assert.Equal(t, anthropic.Model("fallback-model"), boundary.To.Model)
	assert.Equal(t, anthropic.Model("fallback-model"), msg.Model, "the seam relabels the accumulated model")

	// Terminal usage.iterations carries the whole chain.
	require.Len(t, deltas, 1)
	iterations := deltaIterations(t, deltas[0])
	require.Len(t, iterations, 2)
	assert.Equal(t, "message", iterations[0].Type)
	assert.Equal(t, "primary-model", iterations[0].Model)
	assert.Equal(t, "fallback_message", iterations[1].Type)
	assert.Equal(t, "fallback-model", iterations[1].Model)

	// The hop request: model swapped, token redeemed, the advertised claim
	// echoed verbatim as one trailing assistant turn.
	require.Len(t, transport.bodies, 2)
	retry := transport.bodies[1]
	assert.Equal(t, "fallback-model", retry["model"])
	assert.Equal(t, "credit-token-a", retry["fallback_credit_token"])
	assert.Equal(t, true, retry["stream"])
	messages := retry["messages"].([]any)
	require.Len(t, messages, 2, "the claim rides as one appended turn")
	claim := messages[1].(map[string]any)
	assert.Equal(t, "assistant", claim["role"])
	assert.Equal(t, []any{map[string]any{"type": "text", "text": "partial "}}, claim["content"], "echoed verbatim, whitespace included")
	for i, betas := range transport.betas {
		assert.Contains(t, betas, string(anthropic.AnthropicBetaFallbackCredit2026_06_01), "request %d", i)
	}
}

func TestStreamingClaimFalseRetriesTheExactBody(t *testing.T) {
	transport := &sseTransport{responses: []string{
		refusalStream("primary-model", tokenNoClaim), servedStream("fallback-model"),
	}}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})

	_, _, _ = collectStream(t, client, context.Background(), fallbackTestParams)
	require.Len(t, transport.bodies, 2)
	retry := transport.bodies[1]
	assert.Equal(t, "credit-token-a", retry["fallback_credit_token"])
	assert.Len(t, retry["messages"].([]any), 1, "no claim advertised, no turn appended")
}

func TestStreamingMidStreamTokenlessRefusalSurfacesUntouched(t *testing.T) {
	// Content already streamed and no token minted: no retry, the refusal
	// passes through as-is.
	transport := &sseTransport{responses: []string{refusalStream("primary-model", noToken)}}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})

	msg, sequence, _ := collectStream(t, client, context.Background(), fallbackTestParams)
	assert.Equal(t, anthropic.BetaStopReasonRefusal, msg.StopReason)
	require.Len(t, transport.bodies, 1)
	assert.Contains(t, sequence, "message_stop")
	for _, block := range msg.Content {
		assert.NotEqual(t, "fallback", string(block.Type))
	}
}

func TestStreamingPreStreamTokenlessRefusalRetries(t *testing.T) {
	// Nothing streamed yet, so the retry is free and invisible — token or no
	// token.
	preStream := event("message_start", `{"type":"message_start","message":{"type":"message","id":"msg_primary","role":"assistant","content":[],"model":"primary-model","stop_reason":null,"stop_sequence":null,"stop_details":null,"usage":{"input_tokens":10,"output_tokens":1,"cache_creation_input_tokens":0,"cache_read_input_tokens":0}}}`) +
		event("message_delta", `{"type":"message_delta","delta":{"stop_reason":"refusal","stop_sequence":null,"stop_details":{"type":"refusal","category":"x","explanation":null,"fallback_credit_token":null}},"usage":{"input_tokens":10,"output_tokens":1,"cache_read_input_tokens":0,"cache_creation_input_tokens":0}}`) +
		event("message_stop", `{"type":"message_stop"}`)
	transport := &sseTransport{responses: []string{preStream, servedStream("fallback-model")}}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})

	msg, sequence, _ := collectStream(t, client, context.Background(), fallbackTestParams)
	assert.Equal(t, anthropic.BetaStopReasonEndTurn, msg.StopReason)
	require.Len(t, transport.bodies, 2)
	_, hasToken := transport.bodies[1]["fallback_credit_token"]
	assert.False(t, hasToken, "no token minted, none sent")

	// The refused attempt is invisible: one message_start, the boundary
	// queued ahead of the serving hop's content, the primary's message id.
	assert.Equal(t, "message_start", sequence[0])
	require.Len(t, msg.Content, 2)
	assert.Equal(t, anthropic.Model("primary-model"), msg.Content[0].AsFallback().From.Model)
	assert.Equal(t, "msg_primary", msg.ID, "the envelope keeps the primary's message id")
}

func TestStreamingFailedHopSkipsToTheNextEntry(t *testing.T) {
	// An HTTP-failed hop was never reached: same token and continuation move
	// to the next entry, and no boundary is left behind for it.
	transport := &sseTransport{
		responses: []string{refusalStream("primary-model", tokenWithClaim), `{}`, servedStream("fallback-2")},
		statuses:  []int{200, 500, 200},
	}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-1"}, {Model: "fallback-2"}})

	msg, _, _ := collectStream(t, client, context.Background(), fallbackTestParams)
	assert.Equal(t, anthropic.BetaStopReasonEndTurn, msg.StopReason)
	require.Len(t, transport.bodies, 3)
	assert.Equal(t, "credit-token-a", transport.bodies[2]["fallback_credit_token"])
	assert.Len(t, transport.bodies[2]["messages"].([]any), 2, "the claim carries to the next entry")

	var boundaries []string
	for _, block := range msg.Content {
		if fb := block.AsFallback(); fb.To.Model != "" {
			boundaries = append(boundaries, string(fb.From.Model)+"->"+string(fb.To.Model))
		}
	}
	assert.Equal(t, []string{"primary-model->fallback-2"}, boundaries, "no boundary for the hop that never engaged")
}

func TestStreamingExhaustedChainDegradesToTheHeldRefusal(t *testing.T) {
	transport := &sseTransport{
		responses: []string{refusalStream("primary-model", tokenNoClaim), `{}`, `{}`},
		statuses:  []int{200, 500, 503},
	}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-1"}, {Model: "fallback-2"}})

	msg, sequence, deltas := collectStream(t, client, context.Background(), fallbackTestParams)
	assert.Equal(t, anthropic.BetaStopReasonRefusal, msg.StopReason)
	assert.Contains(t, sequence, "message_stop")
	for _, block := range msg.Content {
		assert.NotEqual(t, "fallback", string(block.Type), "failed hops leave no boundary")
	}
	require.Len(t, deltas, 1)
	iterations := deltaIterations(t, deltas[0])
	require.Len(t, iterations, 1, "the held refusal's ledger survives exhaustion")
	assert.Equal(t, "primary-model", iterations[0].Model)
	assert.Equal(t, "message", iterations[0].Type)
	assert.False(t, msg.StopDetails.JSON.RecommendedModel.Valid() && msg.StopDetails.RecommendedModel != "", "non-429 failures recommend nothing")
}

func TestStreamingRateLimitedHopRecommendsTheUnreachableModel(t *testing.T) {
	transport := &sseTransport{
		responses: []string{refusalStream("primary-model", tokenNoClaim), `{}`},
		statuses:  []int{200, 429},
	}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})

	msg, _, _ := collectStream(t, client, context.Background(), fallbackTestParams)
	assert.Equal(t, anthropic.BetaStopReasonRefusal, msg.StopReason)
	assert.Equal(t, anthropic.Model("fallback-model"), msg.StopDetails.RecommendedModel)
}

func TestStreaming400OnAClaimedAttemptRetriesTheLastEntryWithoutTheClaim(t *testing.T) {
	transport := &sseTransport{
		responses: []string{refusalStream("primary-model", tokenWithClaim), `{}`, servedStream("fallback-model")},
		statuses:  []int{200, 400, 200},
	}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})

	msg, _, _ := collectStream(t, client, context.Background(), fallbackTestParams)
	assert.Equal(t, anthropic.BetaStopReasonEndTurn, msg.StopReason)
	require.Len(t, transport.bodies, 3)
	assert.Len(t, transport.bodies[1]["messages"].([]any), 2, "first attempt carries the claim")
	assert.Len(t, transport.bodies[2]["messages"].([]any), 1, "the degrade drops the claim turn")
	assert.Equal(t, "fallback-model", transport.bodies[2]["model"], "same entry retried")
	assert.Equal(t, "credit-token-a", transport.bodies[2]["fallback_credit_token"], "the token is kept")
}

func TestStreamingHopSuppliedIterationsRideThroughLabeled(t *testing.T) {
	// A server-tool run reports several iterations entries; unlabeled ones
	// can only be the hop's own attempts and ride through with its model.
	supplied := `,"iterations":[{"input_tokens":10,"output_tokens":2,"cache_read_input_tokens":0,"cache_creation_input_tokens":0,"type":"message"},{"input_tokens":20,"output_tokens":3,"cache_read_input_tokens":0,"cache_creation_input_tokens":0,"type":"message"}]`
	refusal := strings.Replace(refusalStream("primary-model", tokenNoClaim),
		`"cache_creation_input_tokens":0}}`, `"cache_creation_input_tokens":0`+supplied+`}}`, 1)
	transport := &sseTransport{responses: []string{refusal, servedStream("fallback-model")}}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})

	_, _, deltas := collectStream(t, client, context.Background(), fallbackTestParams)
	require.NotEmpty(t, deltas)
	iterations := deltaIterations(t, deltas[len(deltas)-1])
	require.Len(t, iterations, 3)
	assert.Equal(t, "primary-model", iterations[0].Model)
	assert.Equal(t, "primary-model", iterations[1].Model)
	assert.Equal(t, "fallback-model", iterations[2].Model)
	assert.Equal(t, "fallback_message", iterations[2].Type)
}

func TestStreaming400OnATokenedAttemptRetriesTheLastEntryTokenless(t *testing.T) {
	// Blanket retry chain: the 400 never surfaces; the token is dropped and
	// the same entry retried bare.
	transport := &sseTransport{
		responses: []string{refusalStream("primary-model", tokenNoClaim), `{}`, servedStream("fallback-model")},
		statuses:  []int{200, 400, 200},
	}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})

	msg, _, _ := collectStream(t, client, context.Background(), fallbackTestParams)
	assert.Equal(t, anthropic.BetaStopReasonEndTurn, msg.StopReason)
	require.Len(t, transport.bodies, 3)
	assert.Equal(t, "credit-token-a", transport.bodies[1]["fallback_credit_token"], "first attempt redeems")
	assert.Equal(t, "fallback-model", transport.bodies[2]["model"], "same entry retried")
	_, hasToken := transport.bodies[2]["fallback_credit_token"]
	assert.False(t, hasToken, "the tokenless retry drops the token")
}

func TestStreamingTrimsReplayedFallbackTurns(t *testing.T) {
	// The refused model's thinking and unfinished tool call are removed
	// along with the fallback block; the serving model's text stays.
	transport := &sseTransport{responses: []string{servedStream("fallback-model")}}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})

	params := fallbackTestParams
	msg, _, _ := collectStream(t, client, context.Background(), params,
		option.WithJSONSet("messages", []map[string]any{
			{"role": "user", "content": "hi"},
			{"role": "assistant", "content": []map[string]any{
				{"type": "thinking", "thinking": "hmm", "signature": "sig"},
				{"type": "text", "text": "partial"},
				{"type": "tool_use", "id": "tu_cut", "name": "search", "input": map[string]any{}},
				{"type": "fallback", "from": map[string]any{"model": "primary-model"}, "to": map[string]any{"model": "fallback-model"}},
				{"type": "text", "text": "served"},
			}},
			{"role": "user", "content": "more"},
		}))
	assert.Equal(t, anthropic.BetaStopReasonEndTurn, msg.StopReason)

	require.Len(t, transport.bodies, 1)
	sent := transport.bodies[0]["messages"].([]any)[1].(map[string]any)["content"].([]any)
	var types []string
	for _, block := range sent {
		types = append(types, block.(map[string]any)["type"].(string))
	}
	assert.Equal(t, []string{"text", "text"}, types)
}

func TestStreamingMidBlockRefusalClosesTheOpenBlock(t *testing.T) {
	transport := &sseTransport{responses: []string{
		cutRefusalStream("primary-model", tokenWithClaim), servedStream("fallback-model"),
	}}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})

	msg, sequence, _ := collectStream(t, client, context.Background(), fallbackTestParams)
	assert.Equal(t, anthropic.BetaStopReasonEndTurn, msg.StopReason)
	assert.Contains(t, sequence, "content_block_stop:0", "the cut-open block is closed before the boundary")

	// The synthetic close does not damage the claim: the partial still rides.
	require.Len(t, transport.bodies, 2)
	messages := transport.bodies[1]["messages"].([]any)
	require.Len(t, messages, 2)
}

func TestStreamingPinnedConversationStartsAtTheFallback(t *testing.T) {
	transport := &sseTransport{responses: []string{servedStream("fallback-model")}}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})

	state := &betafallback.BetaFallbackState{}
	state.SetIndex(0)
	msg, _, _ := collectStream(t, client, context.Background(), fallbackTestParams,
		betafallback.WithBetaFallbackState(state))
	assert.Equal(t, anthropic.BetaStopReasonEndTurn, msg.StopReason)
	require.Len(t, transport.bodies, 1)
	assert.Equal(t, "fallback-model", transport.bodies[0]["model"])
	_, hasToken := transport.bodies[0]["fallback_credit_token"]
	assert.False(t, hasToken, "a pinned fresh turn redeems no token")
}

func TestStreamingPinPastTheChainErrors(t *testing.T) {
	transport := &sseTransport{responses: []string{servedStream("fallback-model")}}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})

	stale := &betafallback.BetaFallbackState{}
	stale.SetIndex(1)
	stream := client.Beta.Messages.NewStreaming(context.Background(), fallbackTestParams,
		betafallback.WithBetaFallbackState(stale))
	for stream.Next() {
	}
	require.ErrorContains(t, stream.Err(), "out of range")
}

func TestStreamingServerSideFallbacksError(t *testing.T) {
	// Only one chain can adjudicate refusals; erroring beats silently
	// picking one.
	transport := &sseTransport{responses: []string{servedStream("server-fallback")}}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})

	params := fallbackTestParams
	params.Fallbacks = []anthropic.BetaFallbackParam{{Model: "server-fallback"}}
	stream := client.Beta.Messages.NewStreaming(context.Background(), params)
	for stream.Next() {
	}
	require.ErrorContains(t, stream.Err(), "Sending the `fallbacks:` request param is not supported when using the `BetaRefusalFallbackMiddleware` middleware. You should either remove the middleware and send `fallbacks:` with the `server-side-fallback-2026-06-01` beta header to let the API handle refusal fallbacks, or omit the `fallbacks:` param if you'd like the `BetaRefusalFallbackMiddleware` middleware to handle fallbacks on the client side.")
	assert.Empty(t, transport.bodies, "no request reaches the server")
}

func TestStreamingEmptyChainIgnoresAPin(t *testing.T) {
	transport := &sseTransport{responses: []string{servedStream("primary-model")}}
	client := streamingFallbackClient(t, transport, nil)

	state := &betafallback.BetaFallbackState{}
	state.SetIndex(3)
	msg, _, _ := collectStream(t, client, context.Background(), fallbackTestParams,
		betafallback.WithBetaFallbackState(state))
	assert.Equal(t, anthropic.BetaStopReasonEndTurn, msg.StopReason)
	require.Len(t, transport.bodies, 1)
	assert.Equal(t, "primary-model", transport.bodies[0]["model"], "no pin to apply against an empty chain")
	assert.Empty(t, transport.betas[0])
}

func TestStreamingNonSSEHopIsAFailureNotTruncation(t *testing.T) {
	transport := &sseTransport{
		responses:    []string{refusalStream("primary-model", tokenNoClaim), `{"type": "message"}`},
		contentTypes: []string{"text/event-stream", "application/json"},
	}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})

	msg, sequence, _ := collectStream(t, client, context.Background(), fallbackTestParams)
	assert.Equal(t, anthropic.BetaStopReasonRefusal, msg.StopReason, "must terminate, not truncate")
	assert.Contains(t, sequence, "message_stop")
}

func TestStreamingEncodedResponsePassesThroughUnread(t *testing.T) {
	// A still-encoded body (caller-set Accept-Encoding) can't be spliced.
	middleware := betafallback.BetaRefusalFallbackMiddleware(
		[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
	)
	req, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages?beta=true",
		strings.NewReader(`{"model": "primary-model", "messages": [], "max_tokens": 16, "stream": true}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	calls := 0
	res, err := middleware(req, func(*http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"text/event-stream"}, "Content-Encoding": []string{"gzip"}},
			Body:       io.NopCloser(strings.NewReader("\x1f\x8b not really gzip")),
		}, nil
	})
	require.NoError(t, err)
	assert.Equal(t, 1, calls)
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, "\x1f\x8b not really gzip", string(body))
}

func TestStreamingFallbackRedirectsThroughAPlatformTransform(t *testing.T) {
	// The fallback middleware runs above the platform transform: each hop's
	// model swap must redirect through it, and the signature must cover the
	// body the middleware produced.
	var paths []string
	transport := &sseTransport{responses: []string{
		refusalStream("primary-model", tokenNoClaim), servedStream("fallback-model"),
	}}
	transport.t = t
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			paths = append(paths, req.URL.Path)
			sig := req.Header.Get("X-Fake-Signature")
			buf, err := io.ReadAll(req.Body)
			require.NoError(t, err)
			assert.Equal(t, fmt.Sprintf("%x", sha256.Sum256(buf)), sig, "signature must cover the body the middleware produced")
			req.Body = io.NopCloser(strings.NewReader(string(buf)))
			return transport.RoundTrip(req)
		})}),
		option.WithMaxRetries(0),
		option.WithMiddleware(
			betafallback.BetaRefusalFallbackMiddleware(
				[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
			),
			fakePlatformTransform,
		),
	)

	msg, _, _ := collectStream(t, client, context.Background(), fallbackTestParams)
	assert.Equal(t, anthropic.BetaStopReasonEndTurn, msg.StopReason)
	assert.Equal(t, []string{"/model/primary-model/invoke", "/model/fallback-model/invoke"}, paths)
	require.Len(t, transport.bodies, 2)
	_, hasModel := transport.bodies[1]["model"]
	assert.False(t, hasModel, "the transform moved the model out of the body")
}
