package anthropic

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const threadsBeta = "message-threads-2026-08-12"

// threadResponse is one scripted answer of a threadServer.
type threadResponse struct {
	status int
	header map[string]string
	// body is a JSON body, or the events of a stream when stream is set.
	body   string
	stream bool
	// cutOff drops the connection once body has been written.
	cutOff bool
	// held keeps the connection open once body has been written, until the
	// caller hangs up.
	held bool
}

// threadRequest is one request a threadServer received.
type threadRequest struct {
	body   []byte
	header http.Header
}

// field returns the JSON of a body field, or "" when the body has none.
func (r threadRequest) field(path string) string {
	return gjson.GetBytes(r.body, path).Raw
}

// betas returns the tokens of the anthropic-beta header, sorted.
func (r threadRequest) betas() []string {
	betas := []string{}
	for _, value := range r.header.Values("anthropic-beta") {
		for _, beta := range strings.Split(value, ",") {
			betas = append(betas, strings.TrimSpace(beta))
		}
	}
	slices.Sort(betas)
	return betas
}

// threadServer answers POST /v1/messages with the scripted responses in
// order and records every request. A request past the end of the script fails
// the test.
type threadServer struct {
	*httptest.Server
	mu       sync.Mutex
	received []threadRequest
}

func newThreadServer(t *testing.T, script ...threadResponse) *threadServer {
	t.Helper()
	server := &threadServer{}
	server.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading request body: %v", err)
			return
		}
		server.mu.Lock()
		server.received = append(server.received, threadRequest{body: body, header: r.Header.Clone()})
		n := len(server.received)
		server.mu.Unlock()
		if n > len(script) {
			t.Errorf("unexpected request %d: only %d responses scripted: %s", n, len(script), body)
			http.Error(w, "unscripted", http.StatusNotImplemented)
			return
		}
		response := script[n-1]
		for name, value := range response.header {
			w.Header().Set(name, value)
		}
		if response.stream {
			w.Header().Set("Content-Type", "text/event-stream")
		} else if response.status == 0 && gjson.GetBytes(body, "stream").Bool() {
			// A scripted message is streamed when the request asks for a stream.
			w.Header().Set("Content-Type", "text/event-stream")
			writeMessageEvents(w, response.body)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
		}
		if response.status != 0 {
			w.WriteHeader(response.status)
		}
		_, _ = io.WriteString(w, response.body)
		if response.held {
			w.(http.Flusher).Flush()
			<-r.Context().Done()
		}
		if response.cutOff {
			w.(http.Flusher).Flush()
			panic(http.ErrAbortHandler)
		}
	}))
	t.Cleanup(server.Close)
	return server
}

func (s *threadServer) client(opts ...option.RequestOption) Client {
	return NewClient(slices.Concat([]option.RequestOption{
		option.WithBaseURL(s.URL),
		option.WithAPIKey("test-key"),
		option.WithMaxRetries(0),
	}, opts)...)
}

func (s *threadServer) requests() []threadRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.received)
}

// requireRequests requires that the server received n requests and returns them.
func (s *threadServer) requireRequests(t *testing.T, n int) []threadRequest {
	t.Helper()
	requests := s.requests()
	require.Len(t, requests, n)
	return requests
}

// newThread returns a thread on a client of a server that answers with script.
func newThread(t *testing.T, params BetaMessageThreadCreateParams, script ...threadResponse) (*BetaMessageThread, *threadServer) {
	t.Helper()
	server := newThreadServer(t, script...)
	client := server.client()
	return client.Beta.Messages.NewThread(params), server
}

func threadParams() BetaMessageThreadCreateParams {
	return BetaMessageThreadCreateParams{
		Model:     "claude-sonnet-4-5",
		MaxTokens: 1024,
	}
}

func continueUnthreadedParams() BetaMessageThreadCreateParams {
	params := threadParams()
	params.UnsupportedRequestBehavior = BetaMessageThreadUnsupportedRequestBehaviorContinueUnthreaded
	return params
}

func userMessage(text string) BetaMessageParam {
	return NewBetaUserMessage(NewBetaTextBlock(text))
}

// markedUserMessage is a user message with a cache breakpoint on its block.
func markedUserMessage(text string) BetaMessageParam {
	return NewBetaUserMessage(BetaContentBlockParamUnion{OfText: &BetaTextBlockParam{Text: text, CacheControl: NewBetaCacheControlEphemeralParam()}})
}

func assistantMessage(text string) BetaMessageParam {
	return BetaMessageParam{OfAssistant: &BetaAssistantMessageParam{Content: []BetaContentBlockParamUnion{NewBetaTextBlock(text)}}}
}

// turn returns the params of a turn that sends one user message.
func turn(text string) BetaMessageThreadSendParams {
	return BetaMessageThreadSendParams{Messages: []BetaMessageParam{userMessage(text)}}
}

// markedTurn returns the params of a turn that sends one marked user message.
func markedTurn(text string) BetaMessageThreadSendParams {
	return BetaMessageThreadSendParams{Messages: []BetaMessageParam{markedUserMessage(text)}}
}

// noMessages is a turn that sends what is staged, or resumes.
var noMessages = BetaMessageThreadSendParams{}

func weatherToolParam() BetaToolUnionParam {
	return BetaToolUnionParam{OfTool: &BetaToolParam{
		Name:        "get_weather",
		Description: String("Get the current weather in a given city."),
		InputSchema: BetaToolInputSchemaParam{
			Properties: map[string]any{"city": map[string]any{"type": "string"}},
			Required:   []string{"city"},
		},
	}}
}

func userJSON(text string) string {
	return fmt.Sprintf(`{"role":"user","content":[{"type":"text","text":%q}]}`, text)
}

func markedUserJSON(text string) string {
	return fmt.Sprintf(`{"role":"user","content":[{"type":"text","text":%q,"cache_control":{"type":"ephemeral"}}]}`, text)
}

func assistantJSON(text string) string {
	return fmt.Sprintf(`{"role":"assistant","content":[{"type":"text","text":%q}]}`, text)
}

func messagesOf(messages ...string) string {
	return "[" + strings.Join(messages, ",") + "]"
}

func continueJSON(previousMessageID string) string {
	return fmt.Sprintf(`{"type":"continue","previous_message_id":%q}`, previousMessageID)
}

const createJSON = `{"type":"create"}`

func replyJSON(id, text, stopReason string) string {
	return fmt.Sprintf(`{"id":%q,"type":"message","role":"assistant","model":"claude-sonnet-4-5","content":[{"type":"text","text":%q}],"stop_reason":%q,"stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`, id, text, stopReason)
}

// reply is a whole reply with one text block.
func reply(id, text string) threadResponse {
	return threadResponse{body: replyJSON(id, text, "end_turn")}
}

// pausedReply is a reply that stopped with pause_turn.
func pausedReply(id, text string) threadResponse {
	return threadResponse{body: replyJSON(id, text, "pause_turn")}
}

// refusedReply is a reply that stopped with refusal.
func refusedReply(id string) threadResponse {
	return threadResponse{body: replyJSON(id, "I can't help with that.", "refusal")}
}

const (
	fallbackModel     = "claude-haiku-4-5"
	fallbackBeta      = "server-side-fallback-2026-07-01"
	fallbackBlockJSON = `{"type":"fallback","from":{"model":"claude-sonnet-4-5"},"to":{"model":"claude-haiku-4-5"},"trigger":{"type":"refusal","category":"general_harms"}}`
)

// writtenByFallback returns response as fallbackModel writes it: the content
// starts with a fallback block, and the usage has a fallback_message entry.
func writtenByFallback(response threadResponse) threadResponse {
	iteration := func(kind, model string) string {
		return fmt.Sprintf(`{"type":%q,"model":%q,"input_tokens":1,"output_tokens":1,"cache_creation_input_tokens":0,"cache_read_input_tokens":0,"cache_creation":{"ephemeral_1h_input_tokens":0,"ephemeral_5m_input_tokens":0}}`, kind, model)
	}
	content := append([]string{fallbackBlockJSON}, rawArrayItems(response.body, "content")...)
	response.body, _ = sjson.SetRaw(response.body, "content", messagesOf(content...))
	response.body, _ = sjson.Set(response.body, "model", fallbackModel)
	response.body, _ = sjson.SetRaw(response.body, "usage.iterations", messagesOf(
		iteration("message", "claude-sonnet-4-5"), iteration("fallback_message", fallbackModel), iteration("advisor_message", "claude-opus-4-5"),
	))
	return response
}

func rawArrayItems(json, path string) []string {
	var items []string
	for _, item := range gjson.Get(json, path).Array() {
		items = append(items, item.Raw)
	}
	return items
}

// fallbackAssistantJSON is the reply of writtenByFallback in the history.
func fallbackAssistantJSON(text string) string {
	return fmt.Sprintf(`{"role":"assistant","content":[%s,{"type":"text","text":%q}]}`, fallbackBlockJSON, text)
}

// fallbackParams are params that name fallbackModel as a fallback.
func fallbackParams() BetaMessageThreadCreateParams {
	params := continueUnthreadedParams()
	params.Fallbacks = BetaFallbacksParamUnion{OfBetaFallbackArray: []BetaFallbackParam{{Model: fallbackModel}}}
	params.Betas = []AnthropicBeta{fallbackBeta}
	return params
}

// replyInContainer is a reply that names a container.
func replyInContainer(id, text, containerID string) threadResponse {
	body := strings.Replace(replyJSON(id, text, "end_turn"), `"stop_sequence":null`, fmt.Sprintf(`"stop_sequence":null,"container":{"id":%q,"expires_at":"2026-08-19T00:00:00Z"}`, containerID), 1)
	return threadResponse{body: body}
}

// threadRejection is an error response whose details carry errorCode.
func threadRejection(status int, errorCode string) threadResponse {
	return threadResponse{
		status: status,
		header: map[string]string{"x-should-retry": "false"},
		body:   fmt.Sprintf(`{"type":"error","error":{"type":"invalid_request_error","message":"rejected","details":{"error_code":%q}}}`, errorCode),
	}
}

// unsupportedRequest is a 400 thread_unsupported_request that names unsupported.
func unsupportedRequest(unsupported string) threadResponse {
	return threadResponse{
		status: http.StatusBadRequest,
		body:   fmt.Sprintf(`{"type":"error","error":{"type":"invalid_request_error","message":"rejected","details":{"error_code":"thread_unsupported_request","unsupported":%s}}}`, unsupported),
	}
}

// plainError is an error response without details.
func plainError(status int, message string) threadResponse {
	return threadResponse{
		status: status,
		body:   fmt.Sprintf(`{"type":"error","error":{"type":"invalid_request_error","message":%q}}`, message),
	}
}

func requireStatus(t *testing.T, err error, status int) {
	t.Helper()
	var apiErr *Error
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, status, apiErr.StatusCode)
}

// requireLocalError requires an error that did not come from a response.
func requireLocalError(t *testing.T, err error) {
	t.Helper()
	var apiErr *Error
	require.Error(t, err)
	require.False(t, errors.As(err, &apiErr), "expected an error raised without a response, got %v", err)
}

func requireMessages(t *testing.T, messages []BetaMessageParam, want string) {
	t.Helper()
	requireJSONEqual(t, messagesJSON(t, messages), want)
}

// requireCreate requires a request that creates a thread with messages.
func requireCreate(t *testing.T, request threadRequest, messages string) {
	t.Helper()
	requireJSONEqual(t, request.field("thread"), createJSON)
	requireJSONEqual(t, request.field("messages"), messages)
}

// requireContinue requires a request that continues from previousMessageID
// with messages, and carries neither system nor tools.
func requireContinue(t *testing.T, request threadRequest, previousMessageID string, messages string) {
	t.Helper()
	requireJSONEqual(t, request.field("thread"), continueJSON(previousMessageID))
	requireJSONEqual(t, request.field("messages"), messages)
	require.Empty(t, request.field("system"))
	require.Empty(t, request.field("tools"))
}

// requireUnthreaded requires a request with messages, no thread field and no
// threads beta.
func requireUnthreaded(t *testing.T, request threadRequest, messages string) {
	t.Helper()
	require.Empty(t, request.field("thread"))
	requireJSONEqual(t, request.field("messages"), messages)
	require.NotContains(t, request.betas(), threadsBeta)
}

func send(t *testing.T, thread *BetaMessageThread, params BetaMessageThreadSendParams, opts ...option.RequestOption) *BetaMessage {
	t.Helper()
	message, err := thread.Send(context.Background(), params, opts...)
	require.NoError(t, err)
	return message
}

// sendOrStream sends a turn, as a stream that is read to its end when stream
// is set, and returns how the turn ended.
func sendOrStream(thread *BetaMessageThread, stream bool, params BetaMessageThreadSendParams, opts ...option.RequestOption) error {
	return sendOrStreamContext(context.Background(), thread, stream, params, opts...)
}

func sendOrStreamContext(ctx context.Context, thread *BetaMessageThread, stream bool, params BetaMessageThreadSendParams, opts ...option.RequestOption) error {
	if !stream {
		_, err := thread.Send(ctx, params, opts...)
		return err
	}
	turn := thread.SendStreaming(ctx, params, opts...)
	defer turn.Close()
	for turn.Next() {
	}
	return turn.Err()
}

func TestBetaMessageThread_NewThread_SendsNothingAndHoldsTheSeed(t *testing.T) {
	params := threadParams()
	params.Messages = []BetaMessageParam{userMessage("Hello")}

	thread, server := newThread(t, params)

	require.Empty(t, server.requests())
	require.Equal(t, "", thread.PreviousMessageID())
	require.Nil(t, thread.LastMessage())
	requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello")))
}

func TestBetaMessageThread_Send_FirstTurnCreatesTheThreadWithSystemAndTools(t *testing.T) {
	params := threadParams()
	params.System = []BetaTextBlockParam{{Text: "Be brief."}}
	params.Tools = []BetaToolUnionParam{weatherToolParam()}
	thread, server := newThread(t, params, reply("msg_a", "Reply A"))

	message := send(t, thread, turn("Hello"))

	request := server.requests()[0]
	require.Equal(t, "msg_a", message.ID)
	requireCreate(t, request, messagesOf(userJSON("Hello")))
	requireJSONEqual(t, request.field("system"), `[{"type":"text","text":"Be brief."}]`)
	requireJSONEqual(t, request.field("tools.0.name"), `"get_weather"`)
	require.Equal(t, []string{threadsBeta}, request.betas())
	require.Empty(t, request.field("stream"))
}

func TestBetaMessageThread_Send_RecordsTheReplyInTheHistory(t *testing.T) {
	thread, _ := newThread(t, threadParams(), reply("msg_a", "Reply A"))

	message := send(t, thread, turn("Hello"))

	require.Equal(t, "msg_a", thread.PreviousMessageID())
	require.Same(t, message, thread.LastMessage())
	requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A")))
}

func TestBetaMessageThread_Send_SendsTheSeedFollowedByItsMessage(t *testing.T) {
	params := threadParams()
	params.Messages = []BetaMessageParam{userMessage("Hello"), assistantMessage("Hi there.")}
	thread, server := newThread(t, params, reply("msg_a", "Reply A"))

	send(t, thread, turn("Again"))

	requireCreate(t, server.requests()[0], messagesOf(userJSON("Hello"), assistantJSON("Hi there."), userJSON("Again")))
}

func TestBetaMessageThread_Send_SendsStagedMessagesInOrderAndUnmerged(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"))
	thread.Append(userMessage("part one"))
	thread.Append(userMessage("part two"))

	send(t, thread, noMessages)

	requireCreate(t, server.requests()[0], messagesOf(userJSON("part one"), userJSON("part two")))
}

func TestBetaMessageThread_Append_TheReplyJustRecordedIsANoOp(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))
	message := send(t, thread, turn("Hello"))

	thread.Append(message.ToParam())
	send(t, thread, turn("Again"))

	requireContinue(t, server.requests()[1], "msg_a", messagesOf(userJSON("Again")))
	requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again"), assistantJSON("Reply B")))
}

func TestBetaMessageThread_Send_OnANewThreadWithNothingStagedSendsACreateWithNoMessages(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, threadParams(), plainError(http.StatusBadRequest, "messages: at least one message is required"))

		err := sendOrStream(thread, stream, noMessages)

		requireStatus(t, err, http.StatusBadRequest)
		require.Len(t, server.requests(), 1)
		requireCreate(t, server.requests()[0], `[]`)
	})
}

func TestBetaMessageThread_Send_AfterAFinishedReplyWithNothingStagedSendsAContinueWithNoMessages(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), plainError(http.StatusBadRequest, "nothing to resume"))
		send(t, thread, turn("Hello"))

		err := sendOrStream(thread, stream, noMessages)

		requireStatus(t, err, http.StatusBadRequest)
		require.Len(t, server.requests(), 2)
		requireContinue(t, server.requests()[1], "msg_a", `[]`)
		require.Equal(t, "msg_a", thread.PreviousMessageID())
	})
}

func TestBetaMessageThread_Params_GivenToNewThreadStartANewThreadWithTheSameFieldsAndHistory(t *testing.T) {
	server := newThreadServer(t, reply("msg_a", "Reply A"), reply("msg_b", "Reply B"), reply("msg_c", "Reply C"))
	client := server.client()
	thread := client.Beta.Messages.NewThread(pinnedThreadParams())
	send(t, thread, turn("Hello"))
	thread.Append(userMessage("Again"))

	rebuilt := client.Beta.Messages.NewThread(thread.Params)
	send(t, rebuilt, turn("Once more"))
	send(t, thread, noMessages)

	requests := server.requireRequests(t, 3)
	requireCreate(t, requests[1], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again"), userJSON("Once more")))
	for _, field := range []string{"model", "max_tokens", "system", "tools"} {
		requireJSONEqual(t, requests[1].field(field), requests[0].field(field))
	}
	requireContinue(t, requests[2], "msg_a", messagesOf(userJSON("Again")))
	requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again"), assistantJSON("Reply C")))
}

func TestBetaMessageThread_Params_ACopyAndTheThreadDoNotReachEachOther(t *testing.T) {
	params := threadParams()
	params.Messages = []BetaMessageParam{userMessage("Seed")}
	thread, _ := newThread(t, params, reply("msg_a", "Reply A"))
	send(t, thread, turn("Hello"))
	history := []string{userJSON("Seed"), userJSON("Hello"), assistantJSON("Reply A")}

	copied := thread.Params
	thread.Append(userMessage("Theirs"))
	copied.Messages = append(copied.Messages, userMessage("Mine"))

	requireMessages(t, thread.Messages(), messagesOf(append(slices.Clone(history), userJSON("Theirs"))...))
	requireMessages(t, copied.Messages, messagesOf(append(slices.Clone(history), userJSON("Mine"))...))
}

func TestBetaMessageThread_Params_MessagesHoldsTheHistory(t *testing.T) {
	params := threadParams()
	params.Messages = []BetaMessageParam{userMessage("Seed")}
	thread, _ := newThread(t, params, reply("msg_a", "Reply A"), textCutOff("msg_b", "par"), plainError(http.StatusBadRequest, "rejected"))
	requireMessages(t, thread.Params.Messages, messagesOf(userJSON("Seed")))

	thread.Append(userMessage("Staged"))
	requireMessages(t, thread.Params.Messages, messagesOf(userJSON("Seed"), userJSON("Staged")))

	send(t, thread, turn("Hello"))
	whole := []string{userJSON("Seed"), userJSON("Staged"), userJSON("Hello"), assistantJSON("Reply A")}
	requireMessages(t, thread.Params.Messages, messagesOf(whole...))

	readStream(thread, turn("Again"))
	requireMessages(t, thread.Params.Messages, messagesOf(append(whole, userJSON("Again"), assistantJSON("par"))...))

	_, err := thread.Send(context.Background(), noMessages)
	requireStatus(t, err, http.StatusBadRequest)
	requireMessages(t, thread.Params.Messages, messagesOf(append(whole, userJSON("Again"))...))
	requireMessages(t, thread.Messages(), messagesOf(append(whole, userJSON("Again"))...))
}

func TestBetaMessageThread_Params_AssignedMessagesThatExtendTheHistoryContinueTheThread(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))
	send(t, thread, turn("Hello"))

	thread.Params.Messages = append(thread.Params.Messages, userMessage("Assigned"))
	require.Equal(t, "msg_a", thread.PreviousMessageID())
	send(t, thread, turn("Again"))

	requireContinue(t, server.requests()[1], "msg_a", messagesOf(userJSON("Assigned"), userJSON("Again")))
}

func TestBetaMessageThread_Params_AssignedMessagesThatRewriteTheHistoryStartANewThread(t *testing.T) {
	histories := map[string][]BetaMessageParam{
		"other":   {userMessage("Other"), assistantMessage("Reply A")},
		"shorter": {userMessage("Hello")},
		"none":    nil,
	}
	for name, history := range histories {
		t.Run(name, func(t *testing.T) {
			thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"), reply("msg_c", "Reply C"))
			send(t, thread, turn("Hello"))

			thread.Params.Messages = history
			require.Nil(t, thread.LastMessage())
			require.Equal(t, "", thread.PreviousMessageID())
			send(t, thread, turn("Again"))
			send(t, thread, turn("Once more"))

			requireCreate(t, server.requests()[1], messagesJSON(t, append(slices.Clone(history), userMessage("Again"))))
			requireContinue(t, server.requests()[2], "msg_b", messagesOf(userJSON("Once more")))
		})
	}
}

func TestBetaMessageThread_Params_ASentMessageChangedInPlaceStartsANewThread(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))
		send(t, thread, turn("Hello"))

		thread.Params.Messages[0].OfUser.Content = userMessage("Other").OfUser.Content
		require.Equal(t, "", thread.PreviousMessageID())
		require.NoError(t, sendOrStream(thread, stream, turn("Again")))

		requireCreate(t, server.requireRequests(t, 2)[1], messagesOf(userJSON("Other"), assistantJSON("Reply A"), userJSON("Again")))
	})
}

func TestBetaMessageThread_AMessageTheCallerChangesAfterGivingItStaysAsItWasGiven(t *testing.T) {
	ways := map[string]func(t *testing.T, own []BetaMessageParam, script ...threadResponse) (*BetaMessageThread, *threadServer){
		"seeded": func(t *testing.T, own []BetaMessageParam, script ...threadResponse) (*BetaMessageThread, *threadServer) {
			params := threadParams()
			params.Messages = own
			thread, server := newThread(t, params, script...)
			send(t, thread, noMessages)
			return thread, server
		},
		"appended": func(t *testing.T, own []BetaMessageParam, script ...threadResponse) (*BetaMessageThread, *threadServer) {
			thread, server := newThread(t, threadParams(), script...)
			thread.Append(own...)
			send(t, thread, noMessages)
			return thread, server
		},
		"sent": func(t *testing.T, own []BetaMessageParam, script ...threadResponse) (*BetaMessageThread, *threadServer) {
			thread, server := newThread(t, threadParams(), script...)
			send(t, thread, BetaMessageThreadSendParams{Messages: own})
			return thread, server
		},
		"streamed": func(t *testing.T, own []BetaMessageParam, script ...threadResponse) (*BetaMessageThread, *threadServer) {
			thread, server := newThread(t, threadParams(), script...)
			require.NoError(t, sendOrStream(thread, true, BetaMessageThreadSendParams{Messages: own}))
			return thread, server
		},
	}
	for name, give := range ways {
		t.Run(name, func(t *testing.T) {
			own := []BetaMessageParam{userMessage("Hello")}
			thread, server := give(t, own, reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))

			own[0].OfUser.Content = userMessage("Changed").OfUser.Content
			send(t, thread, turn("Again"))

			requireContinue(t, server.requireRequests(t, 2)[1], "msg_a", messagesOf(userJSON("Again")))
			requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again"), assistantJSON("Reply B")))
		})
	}
}

func TestBetaMessageThread_Messages_AChangeToTheCopyDoesNotReachTheThread(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))
	send(t, thread, turn("Hello"))

	thread.Messages()[0].OfUser.Content = userMessage("Other").OfUser.Content
	send(t, thread, turn("Again"))

	requireContinue(t, server.requireRequests(t, 2)[1], "msg_a", messagesOf(userJSON("Again")))
	requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again"), assistantJSON("Reply B")))
}

func TestBetaMessageThread_SendStreaming_AMessageAddedWhileTheStreamIsOpenFollowsTheReply(t *testing.T) {
	ways := map[string]func(thread *BetaMessageThread, message BetaMessageParam){
		"appended": func(thread *BetaMessageThread, message BetaMessageParam) { thread.Append(message) },
		"assigned": func(thread *BetaMessageThread, message BetaMessageParam) {
			thread.Params.Messages = append(thread.Params.Messages, message)
		},
	}
	for name, add := range ways {
		t.Run(name, func(t *testing.T) {
			thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))
			stream := thread.SendStreaming(context.Background(), turn("Hello"))
			require.True(t, stream.Next())

			add(thread, userMessage("Added"))
			read(stream)
			requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Added")))
			send(t, thread, noMessages)

			requireContinue(t, server.requireRequests(t, 2)[1], "msg_a", messagesOf(userJSON("Added")))
		})
	}
}

func TestBetaMessageThread_Params_MessagesAssignedBeforeTheFirstTurnAreSent(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"))
	thread.Params.Messages = []BetaMessageParam{userMessage("Assigned")}

	send(t, thread, turn("Hello"))

	requireCreate(t, server.requests()[0], messagesOf(userJSON("Assigned"), userJSON("Hello")))
}

// Betas and headers.

func TestBetaMessageThread_Send_SendsTheThreadsBetaOnceWhenTheCallerNamesIt(t *testing.T) {
	params := threadParams()
	params.Betas = []AnthropicBeta{AnthropicBetaContextManagement2025_06_27, AnthropicBetaMessageThreads2026_08_12}
	thread, server := newThread(t, params, reply("msg_a", "Reply A"))

	send(t, thread, turn("Hello"))

	require.Equal(t, []string{"context-management-2025-06-27", threadsBeta}, server.requests()[0].betas())
}

func TestBetaMessageThread_Send_AddsTheThreadsBetaToABetaHeaderGivenAtCreation(t *testing.T) {
	server := newThreadServer(t, reply("msg_a", "Reply A"))
	client := server.client()
	thread := client.Beta.Messages.NewThread(threadParams(), option.WithHeader("anthropic-beta", "context-management-2025-06-27"))

	send(t, thread, turn("Hello"))

	require.Equal(t, []string{"context-management-2025-06-27", threadsBeta}, server.requests()[0].betas())
}

func TestBetaMessageThread_Send_UnthreadedRequestCarriesTheCallersBetaHeaderAsGiven(t *testing.T) {
	server := newThreadServer(t, unsupportedRequest(`["inference_hooks"]`), reply("msg_a", "Reply A"))
	client := server.client()
	thread := client.Beta.Messages.NewThread(continueUnthreadedParams(), option.WithHeader("anthropic-beta", "context-management-2025-06-27"))

	send(t, thread, turn("Hello"))

	require.Equal(t, []string{"context-management-2025-06-27"}, server.requests()[1].betas())
}

func TestBetaMessageThread_Send_EveryRequestCarriesTheHelperHeader(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, continueUnthreadedParams(),
			reply("msg_a", "Reply A"),
			threadRejection(http.StatusNotFound, "thread_not_found"), reply("msg_b", "Reply B"),
			unsupportedRequest(`["inference_hooks"]`), reply("msg_c", "Reply C"),
		)

		for _, text := range []string{"Hello", "Again", "Once more"} {
			require.NoError(t, sendOrStream(thread, stream, turn(text)))
		}

		requests := server.requireRequests(t, 5)
		var links []string
		for _, request := range requests {
			links = append(links, gjson.GetBytes(request.body, "thread.type").String())
			require.Equal(t, []string{"message-thread"}, request.header.Values("x-stainless-helper"))
		}
		require.Equal(t, []string{"create", "continue", "create", "continue", ""}, links)
	})
}

func TestBetaMessageThread_Send_TheCallersOptionsKeepTheHelperHeader(t *testing.T) {
	server := newThreadServer(t, reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))
	client := server.client()
	thread := client.Beta.Messages.NewThread(threadParams(), option.WithHeader("x-stainless-helper", "mine"))

	send(t, thread, turn("Hello"), option.WithHeader("x-turn", "1"))
	thread.Options = nil
	send(t, thread, turn("Again"))

	requests := server.requireRequests(t, 2)
	require.Equal(t, []string{"mine, message-thread"}, requests[0].header.Values("x-stainless-helper"))
	require.Equal(t, []string{"message-thread"}, requests[1].header.Values("x-stainless-helper"))
}

func TestBetaMessageThread_Send_AnOptionOfTheTurnReplacesTheThreadsForThatRequest(t *testing.T) {
	overloaded := threadResponse{status: http.StatusInternalServerError, header: map[string]string{"retry-after-ms": "1"}, body: `{}`}
	server := newThreadServer(t, overloaded, overloaded, reply("msg_a", "Reply A"))
	client := server.client()
	thread := client.Beta.Messages.NewThread(threadParams(),
		option.WithHeader("x-replaced", "thread"),
		option.WithHeader("x-removed", "thread"),
		option.WithHeader("x-kept", "thread"),
		option.WithRequestTimeout(30*time.Second),
		option.WithMaxRetries(1),
	)

	_, err := thread.Send(context.Background(), turn("Hello"),
		option.WithHeader("x-replaced", "turn"),
		option.WithHeaderDel("x-removed"),
		option.WithRequestTimeout(10*time.Second),
		option.WithMaxRetries(0),
	)
	requireStatus(t, err, http.StatusInternalServerError)
	first := server.requireRequests(t, 1)[0].header
	send(t, thread, noMessages)
	last := server.requireRequests(t, 3)[2].header

	for name, want := range map[string][2]string{
		"x-replaced":          {"turn", "thread"},
		"x-removed":           {"", "thread"},
		"x-kept":              {"thread", "thread"},
		"x-stainless-timeout": {"10", "30"},
	} {
		require.Equal(t, want, [2]string{first.Get(name), last.Get(name)}, name)
	}
}

// Fields.

// setNonZero sets value to something other than its zero value. Of a struct it
// sets the first field only, because a union takes one variant.
func setNonZero(t *testing.T, value reflect.Value) {
	t.Helper()
	switch value.Kind() {
	case reflect.String:
		value.SetString("x")
	case reflect.Int64:
		value.SetInt(1)
	case reflect.Float64:
		value.SetFloat(1)
	case reflect.Bool:
		value.SetBool(true)
	case reflect.Interface:
		value.Set(reflect.ValueOf(map[string]any{"type": "object"}))
	case reflect.Slice:
		value.Set(reflect.MakeSlice(value.Type(), 1, 1))
		setNonZero(t, value.Index(0))
	case reflect.Pointer:
		value.Set(reflect.New(value.Type().Elem()))
		setNonZero(t, value.Elem())
	case reflect.Struct:
		setNonZero(t, value.Field(0))
	}
	require.False(t, value.IsZero(), "setNonZero cannot set a %s", value.Type())
}

// setEveryField sets every exported field of the struct that params points to,
// and the field of the same name on want if there is one.
func setEveryField(t *testing.T, params any, want *BetaMessageNewParams) {
	t.Helper()
	fields := reflect.ValueOf(params).Elem()
	for i := range fields.NumField() {
		field := fields.Type().Field(i)
		if !field.IsExported() {
			continue
		}
		setNonZero(t, fields.Field(i))
		if wanted := reflect.ValueOf(want).Elem().FieldByName(field.Name); wanted.IsValid() {
			setNonZero(t, wanted)
		}
	}
}

func TestBetaMessageThread_Send_SendsEveryFieldAsNewDoes(t *testing.T) {
	routes := map[string]func(t *testing.T, client Client, want *BetaMessageNewParams){
		"given to NewThread": func(t *testing.T, client Client, want *BetaMessageNewParams) {
			var params BetaMessageThreadCreateParams
			setEveryField(t, &params, want)
			send(t, client.Beta.Messages.NewThread(params), noMessages)
		},
		"given to the turn": func(t *testing.T, client Client, want *BetaMessageNewParams) {
			var params BetaMessageThreadSendParams
			setEveryField(t, &params, want)
			send(t, client.Beta.Messages.NewThread(BetaMessageThreadCreateParams{}), params)
		},
		"given to a tool runner": func(t *testing.T, client Client, want *BetaMessageNewParams) {
			setEveryField(t, &BetaMessageThreadCreateParams{}, want)
			// The runner takes its tools apart from the params, and refuses Compaction.
			want.Tools, want.Compaction = []BetaToolUnionParam{}, BetaCompactionConfigUnionParam{}
			_, err := client.Beta.Messages.NewToolRunner(nil, BetaToolRunnerParams{BetaMessageNewParams: *want}).NextMessage(context.Background())
			require.NoError(t, err)
		},
	}
	for route, sendEveryField := range routes {
		t.Run(route, func(t *testing.T) {
			server := newThreadServer(t, reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))
			client := server.client()
			var want BetaMessageNewParams

			sendEveryField(t, client, &want)
			_, err := client.Beta.Messages.New(context.Background(), want)
			require.NoError(t, err)

			threaded, plain := server.requests()[0], server.requests()[1]
			body, err := sjson.DeleteBytes(threaded.body, "thread")
			require.NoError(t, err)
			requireJSONEqual(t, string(body), string(plain.body))
			require.ElementsMatch(t, append(plain.betas(), threadsBeta), threaded.betas())
			for name := range plain.header {
				if strings.HasPrefix(name, "Anthropic-") && name != "Anthropic-Beta" {
					require.Equal(t, plain.header.Values(name), threaded.header.Values(name), name)
				}
			}
		})
	}
}

func TestBetaMessageThread_Send_FieldOfOneTurnAppliesToThatRequestOnly(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))
	params := turn("Hello")
	params.MaxTokens = 12
	params.StopSequences = []string{"END"}

	send(t, thread, params)
	send(t, thread, turn("Again"))

	first, second := server.requests()[0], server.requests()[1]
	requireJSONEqual(t, first.field("max_tokens"), `12`)
	requireJSONEqual(t, first.field("stop_sequences"), `["END"]`)
	requireJSONEqual(t, second.field("max_tokens"), `1024`)
	require.Empty(t, second.field("stop_sequences"))
}

func TestBetaMessageThread_Params_AssignedFieldsGoOutOnEveryLaterRequest(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"), reply("msg_c", "Reply C"))
	send(t, thread, turn("Hello"))

	thread.Params.MaxTokens = 77
	thread.Params.StopSequences = []string{"END"}
	send(t, thread, turn("Again"))
	send(t, thread, turn("Once more"))

	second, third := server.requests()[1], server.requests()[2]
	requireContinue(t, second, "msg_a", messagesOf(userJSON("Again")))
	requireJSONEqual(t, second.field("max_tokens"), `77`)
	requireJSONEqual(t, second.field("stop_sequences"), `["END"]`)
	requireContinue(t, third, "msg_b", messagesOf(userJSON("Once more")))
	requireJSONEqual(t, third.field("max_tokens"), `77`)
}

func TestBetaMessageThread_Send_FingerprintMismatchLeavesTheThreadAsItWas(t *testing.T) {
	thread, _ := newThread(t, threadParams(), reply("msg_a", "Reply A"), threadRejection(http.StatusBadRequest, "thread_fingerprint_mismatch"))
	send(t, thread, turn("Hello"))

	_, err := thread.Send(context.Background(), turn("Again"))

	requireStatus(t, err, http.StatusBadRequest)
	require.Equal(t, "msg_a", thread.PreviousMessageID())
	requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A")))
}

// pinnedThreadParams are params with a System and Tools.
func pinnedThreadParams() BetaMessageThreadCreateParams {
	params := continueUnthreadedParams()
	params.System = []BetaTextBlockParam{{Text: "Be brief."}}
	params.Tools = []BetaToolUnionParam{weatherToolParam()}
	return params
}

func TestBetaMessageThread_Send_AChangedSystemOrToolsIsRefusedOnATurnThatContinuesTheThread(t *testing.T) {
	changes := map[string]func(params *BetaMessageThreadCreateParams){
		"System": func(params *BetaMessageThreadCreateParams) {
			params.System = []BetaTextBlockParam{{Text: "Be verbose."}}
		},
		"Tools": func(params *BetaMessageThreadCreateParams) { params.Tools = nil },
		"System and Tools": func(params *BetaMessageThreadCreateParams) {
			params.System, params.Tools = nil, nil
		},
	}
	for field, change := range changes {
		t.Run(field, func(t *testing.T) {
			forEachRunner(t, func(t *testing.T, stream bool) {
				thread, server := newThread(t, pinnedThreadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))
				send(t, thread, turn("Hello"))
				change(&thread.Params)

				err := sendOrStream(thread, stream, turn("Again"))

				requireLocalError(t, err)
				require.ErrorContains(t, err, field+" cannot be changed on a message thread")
				require.Equal(t, strings.Contains(field, "Tools"), strings.HasSuffix(err.Error(), "send tool_addition or tool_removal blocks in a system message"))
				require.Len(t, server.requests(), 1)
				require.Equal(t, "msg_a", thread.PreviousMessageID())

				thread.Params.System, thread.Params.Tools = pinnedThreadParams().System, pinnedThreadParams().Tools
				require.NoError(t, sendOrStream(thread, stream, turn("Again")))

				requireContinue(t, server.requests()[1], "msg_a", messagesOf(userJSON("Again")))
			})
		})
	}
}

func TestBetaMessageThread_Send_AToolRemovalInASystemMessageGoesOutOnTheContinue(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, pinnedThreadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))
		send(t, thread, turn("Hello"))
		thread.Append(systemToolChange(NewBetaToolRemovalBlock(BetaToolChangeToolReferenceParam{Name: "get_weather"})))

		require.NoError(t, sendOrStream(thread, stream, turn("Again")))

		requireContinue(t, server.requireRequests(t, 2)[1], "msg_a", messagesOf(toolChangeMessageJSON(toolRemovalJSON("get_weather")), userJSON("Again")))
	})
}

// Extra fields are marshaled in a random order, so a single turn could pass by
// chance.
func TestBetaMessageThread_Send_ToolsWithExtraFieldsContinueTheThread(t *testing.T) {
	params := pinnedThreadParams()
	params.Tools[0].OfTool.InputSchema.ExtraFields = map[string]any{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5, "f": 6}
	script := make([]threadResponse, 10)
	for i := range script {
		script[i] = reply(fmt.Sprintf("msg_%d", i), "Reply")
	}
	thread, server := newThread(t, params, script...)

	for range script {
		send(t, thread, turn("Hello"))
	}

	for _, request := range server.requireRequests(t, len(script))[1:] {
		requireJSONEqual(t, request.field("thread.type"), `"continue"`)
	}
}

func TestBetaMessageThread_Send_AChangedSystemIsSentByATurnThatCreatesTheThread(t *testing.T) {
	thread, server := newThread(t, pinnedThreadParams(), reply("msg_a", "Reply A"))
	thread.Params.System = []BetaTextBlockParam{{Text: "Be verbose."}}

	send(t, thread, turn("Hello"))

	requireJSONEqual(t, server.requests()[0].field("system.0.text"), `"Be verbose."`)
}

func TestBetaMessageThread_Send_TheSystemOfTheTurnThatCreatesTheThreadStaysTheThreads(t *testing.T) {
	thread, server := newThread(t, pinnedThreadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))
	first := turn("Hello")
	first.System = []BetaTextBlockParam{{Text: "Be verbose."}}
	send(t, thread, first)

	send(t, thread, turn("Again"))

	requests := server.requireRequests(t, 2)
	requireJSONEqual(t, requests[0].field("system.0.text"), `"Be verbose."`)
	requireContinue(t, requests[1], "msg_a", messagesOf(userJSON("Again")))
	require.Equal(t, first.System, thread.Params.System)
}

func TestBetaMessageThread_Fork_WithAChangedSystemLeavesTheOriginalContinuing(t *testing.T) {
	thread, server := newThread(t, pinnedThreadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"), reply("msg_c", "Reply C"))
	send(t, thread, turn("Hello"))
	fork := thread.Fork()
	fork.Params.System = []BetaTextBlockParam{{Text: "Be verbose."}}
	send(t, fork, turn("Again"))

	send(t, thread, turn("Once more"))

	requireContinue(t, server.requests()[2], "msg_a", messagesOf(userJSON("Once more")))
}

// Recovery.

// requireReplayedUnderCreate requires that a continue rejected with status and
// errorCode is sent again as a create inside the same Send, and that the next
// turn continues from the reply of that create.
func requireReplayedUnderCreate(t *testing.T, status int, errorCode string) {
	t.Helper()
	params := threadParams()
	params.System = []BetaTextBlockParam{{Text: "Be brief."}}
	params.Tools = []BetaToolUnionParam{weatherToolParam()}
	thread, server := newThread(t, params, reply("msg_a", "Reply A"), threadRejection(status, errorCode), reply("msg_b", "Reply B"), reply("msg_c", "Reply C"))
	send(t, thread, turn("Hello"))

	message := send(t, thread, turn("Again"))
	send(t, thread, turn("Once more"))

	requests := server.requests()
	require.Equal(t, "msg_b", message.ID)
	requireContinue(t, requests[1], "msg_a", messagesOf(userJSON("Again")))
	requireCreate(t, requests[2], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again")))
	requireJSONEqual(t, requests[2].field("system"), `[{"type":"text","text":"Be brief."}]`)
	requireJSONEqual(t, requests[2].field("tools.0.name"), `"get_weather"`)
	requireContinue(t, requests[3], "msg_b", messagesOf(userJSON("Once more")))
}

func TestBetaMessageThread_Send_ThreadNotFoundIsReplayedUnderCreate(t *testing.T) {
	requireReplayedUnderCreate(t, http.StatusNotFound, "thread_not_found")
}

func TestBetaMessageThread_Send_ThreadAlreadyContinuedIsReplayedUnderCreate(t *testing.T) {
	requireReplayedUnderCreate(t, http.StatusBadRequest, "thread_already_continued")
}

func TestBetaMessageThread_Send_ThreadContinueInFlightIsReplayedUnderCreate(t *testing.T) {
	requireReplayedUnderCreate(t, http.StatusBadRequest, "thread_continue_in_flight")
}

func TestBetaMessageThread_Send_ANotFoundWithoutAThreadErrorCodeIsReturned(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), plainError(http.StatusNotFound, "model not found"))
	send(t, thread, turn("Hello"))

	_, err := thread.Send(context.Background(), turn("Again"))

	requireStatus(t, err, http.StatusNotFound)
	require.Len(t, server.requests(), 2)
	require.Equal(t, "msg_a", thread.PreviousMessageID())
}

func TestBetaMessageThread_Send_ThreadErrorOnACreateIsReturnedAfterOneRequest(t *testing.T) {
	thread, server := newThread(t, threadParams(), threadRejection(http.StatusNotFound, "thread_not_found"))

	_, err := thread.Send(context.Background(), turn("Hello"))

	requireStatus(t, err, http.StatusNotFound)
	require.Len(t, server.requests(), 1)
}

func TestBetaMessageThread_Send_TurnAfterAFailedReplayIsACreate(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), threadRejection(http.StatusNotFound, "thread_not_found"), plainError(http.StatusBadRequest, "bad request"), reply("msg_b", "Reply B"))
	send(t, thread, turn("Hello"))
	_, _ = thread.Send(context.Background(), turn("Again"))

	send(t, thread, turn("Again"))

	requireCreate(t, server.requests()[3], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again")))
}

func TestBetaMessageThread_Send_UnknownUnsupportedValueIsReturnedByDefault(t *testing.T) {
	thread, server := newThread(t, threadParams(), unsupportedRequest(`["something_new"]`))

	_, err := thread.Send(context.Background(), turn("Hello"))

	requireStatus(t, err, http.StatusBadRequest)
	require.Len(t, server.requests(), 1)
}

func TestBetaMessageThread_Send_ContinueUnthreadedSendsTheTurnAgainWithoutThread(t *testing.T) {
	thread, server := newThread(t, continueUnthreadedParams(), reply("msg_a", "Reply A"), unsupportedRequest(`["inference_hooks"]`), reply("msg_b", "Reply B"))
	send(t, thread, turn("Hello"))

	message := send(t, thread, turn("Again"))

	require.Equal(t, "msg_b", message.ID)
	requireContinue(t, server.requests()[1], "msg_a", messagesOf(userJSON("Again")))
	requireUnthreaded(t, server.requests()[2], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again")))
	require.Equal(t, "", thread.PreviousMessageID())
}

func notFound() threadResponse {
	return threadRejection(http.StatusNotFound, "thread_not_found")
}

func TestBetaMessageThread_Send_ContinueUnthreadedStopsThreadingWhenTheReplayedThreadIsNotFoundEither(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, continueUnthreadedParams(), reply("msg_a", "Reply A"), notFound(), reply("msg_b", "Reply B"), notFound(), reply("msg_c", "Reply C"), reply("msg_d", "Reply D"))
		send(t, thread, turn("Hello"))
		require.NoError(t, sendOrStream(thread, stream, turn("Again")))

		require.NoError(t, sendOrStream(thread, stream, turn("Once more")))
		require.NoError(t, sendOrStream(thread, stream, turn("Last")))

		requests := server.requireRequests(t, 6)
		history := []string{userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again"), assistantJSON("Reply B"), userJSON("Once more")}
		requireCreate(t, requests[2], messagesOf(history[:3]...))
		requireContinue(t, requests[3], "msg_b", messagesOf(userJSON("Once more")))
		requireUnthreaded(t, requests[4], messagesOf(history...))
		requireUnthreaded(t, requests[5], messagesOf(append(history, assistantJSON("Reply C"), userJSON("Last"))...))
	})
}

func TestBetaMessageThread_Send_ByDefaultAReplayedThreadThatIsNotFoundIsReplayedAgain(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), notFound(), reply("msg_b", "Reply B"), notFound(), reply("msg_c", "Reply C"), reply("msg_d", "Reply D"))
		send(t, thread, turn("Hello"))
		require.NoError(t, sendOrStream(thread, stream, turn("Again")))

		require.NoError(t, sendOrStream(thread, stream, turn("Once more")))
		require.NoError(t, sendOrStream(thread, stream, turn("Last")))

		requests := server.requireRequests(t, 6)
		requireContinue(t, requests[3], "msg_b", messagesOf(userJSON("Once more")))
		requireCreate(t, requests[4], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again"), assistantJSON("Reply B"), userJSON("Once more")))
		requireContinue(t, requests[5], "msg_c", messagesOf(userJSON("Last")))
	})
}

func TestBetaMessageThread_Send_ContinueUnthreadedKeepsThreadingWhenAContinueSucceededSinceTheReplay(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, continueUnthreadedParams(), reply("msg_a", "Reply A"), notFound(), reply("msg_b", "Reply B"), reply("msg_c", "Reply C"), notFound(), reply("msg_d", "Reply D"))
		send(t, thread, turn("Hello"))
		require.NoError(t, sendOrStream(thread, stream, turn("Again")))
		require.NoError(t, sendOrStream(thread, stream, turn("Once more")))

		require.NoError(t, sendOrStream(thread, stream, turn("Last")))

		requests := server.requireRequests(t, 6)
		requireContinue(t, requests[3], "msg_b", messagesOf(userJSON("Once more")))
		requireContinue(t, requests[4], "msg_c", messagesOf(userJSON("Last")))
		requireJSONEqual(t, requests[5].field("thread"), createJSON)
		require.Equal(t, "msg_d", thread.PreviousMessageID())
	})
}

func TestBetaMessageThread_Send_ContinueUnthreadedKeepsThreadingWhenOneOfTheTwoRejectionsIsNotThreadNotFound(t *testing.T) {
	for _, errorCode := range []string{"thread_already_continued", "thread_continue_in_flight"} {
		other := threadRejection(http.StatusBadRequest, errorCode)
		for name, rejections := range map[string][2]threadResponse{
			errorCode + " then thread_not_found": {other, notFound()},
			"thread_not_found then " + errorCode: {notFound(), other},
		} {
			t.Run(name, func(t *testing.T) {
				thread, server := newThread(t, continueUnthreadedParams(), reply("msg_a", "Reply A"), rejections[0], reply("msg_b", "Reply B"), rejections[1], reply("msg_c", "Reply C"))
				send(t, thread, turn("Hello"))
				send(t, thread, turn("Again"))

				send(t, thread, turn("Once more"))

				requireJSONEqual(t, server.requireRequests(t, 5)[4].field("thread"), createJSON)
				require.Equal(t, "msg_c", thread.PreviousMessageID())
			})
		}
	}
}

func TestBetaMessageThread_Params_GivenToNewThreadKeepTheUnsupportedRequestBehavior(t *testing.T) {
	server := newThreadServer(t, unsupportedRequest(`["inference_hooks"]`), reply("msg_a", "Reply A"))
	client := server.client()
	thread := client.Beta.Messages.NewThread(continueUnthreadedParams())

	rebuilt := client.Beta.Messages.NewThread(thread.Params)
	send(t, rebuilt, turn("Hello"))

	requireCreate(t, server.requests()[0], messagesOf(userJSON("Hello")))
	requireUnthreaded(t, server.requests()[1], messagesOf(userJSON("Hello")))
}

func TestBetaMessageThread_Params_AnAssignedUnsupportedRequestBehaviorAppliesToLaterTurns(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), unsupportedRequest(`["inference_hooks"]`), reply("msg_b", "Reply B"))
	send(t, thread, turn("Hello"))

	thread.Params.UnsupportedRequestBehavior = BetaMessageThreadUnsupportedRequestBehaviorContinueUnthreaded
	send(t, thread, turn("Again"))

	requireContinue(t, server.requests()[1], "msg_a", messagesOf(userJSON("Again")))
	requireUnthreaded(t, server.requests()[2], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again")))
}

// Refusals.

func TestBetaMessageThread_Send_ARefusedTurnIsLeftOutOfTheNewThreadThatTheNextTurnStarts(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, pinnedThreadParams(), reply("msg_a", "Reply A"), refusedReply("msg_b"), reply("msg_c", "Reply C"), reply("msg_d", "Reply D"))
		send(t, thread, turn("Hello"))

		require.NoError(t, sendOrStream(thread, stream, turn("Refused")))
		require.NoError(t, sendOrStream(thread, stream, turn("Again")))
		require.NoError(t, sendOrStream(thread, stream, turn("Once more")))

		requests := server.requireRequests(t, 4)
		requireContinue(t, requests[1], "msg_a", messagesOf(userJSON("Refused")))
		requireCreate(t, requests[2], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again")))
		requireJSONEqual(t, requests[2].field("system.0.text"), `"Be brief."`)
		requireJSONEqual(t, requests[2].field("tools.0.name"), `"get_weather"`)
		requireContinue(t, requests[3], "msg_c", messagesOf(userJSON("Once more")))
	})
}

func TestBetaMessageThread_Send_ARefusedTurnIsNotInTheHistory(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, _ := newThread(t, threadParams(), reply("msg_a", "Reply A"), refusedReply("msg_b"))
		send(t, thread, turn("Hello"))

		require.NoError(t, sendOrStream(thread, stream, turn("Refused")))

		want := messagesOf(userJSON("Hello"), assistantJSON("Reply A"))
		requireMessages(t, thread.Messages(), want)
		requireMessages(t, thread.Params.Messages, want)
		require.Equal(t, "msg_a", thread.LastMessage().ID)
		require.Equal(t, "", thread.PreviousMessageID())
	})
}

func TestBetaMessageThread_Send_ReturnsARefusedReplyAsItIs(t *testing.T) {
	thread, _ := newThread(t, threadParams(), refusedReply("msg_a"))

	message := send(t, thread, turn("Refused"))

	requireJSONEqual(t, message.RawJSON(), refusedReply("msg_a").body)
}

func TestBetaMessageThread_SendStreaming_ForwardsEveryEventOfARefusedReply(t *testing.T) {
	thread, _ := newThread(t, threadParams(), refusedReply("msg_a"))
	stream := thread.SendStreaming(context.Background(), turn("Refused"))
	defer stream.Close()

	var message BetaMessage
	for stream.Next() {
		require.NoError(t, message.Accumulate(stream.Current()))
	}

	require.NoError(t, stream.Err())
	require.Equal(t, "msg_a", message.ID)
	require.Equal(t, BetaStopReasonRefusal, message.StopReason)
	require.Equal(t, "I can't help with that.", message.Content[0].Text)
}

func TestBetaMessageThread_SendStreaming_ClosedAfterTheRefusalAndBeforeMessageStopIsNotRecorded(t *testing.T) {
	held := heldStream(messageStartEvent("msg_b"), textEvents(0, "I can't help with that."), messageDeltaEvent("refusal"))
	thread, server := newThread(t, pinnedThreadParams(), reply("msg_a", "Reply A"), held, reply("msg_c", "Reply C"))
	send(t, thread, turn("Hello"))
	thread.Append(userMessage("Staged"))
	stream := thread.SendStreaming(context.Background(), turn("Refused"))
	for stream.Next() && stream.Current().Type != "message_delta" {
	}

	require.NoError(t, stream.Close())

	requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Staged")))
	require.Equal(t, "msg_a", thread.LastMessage().ID)
	send(t, thread, turn("Again"))
	request := server.requireRequests(t, 3)[2]
	requireCreate(t, request, messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Staged"), userJSON("Again")))
	require.NotEmpty(t, request.field("system"))
	require.NotEmpty(t, request.field("tools"))
}

func TestBetaMessageThread_Send_ARefusedFirstTurnLeavesTheSeedStagedAndKeepsNoneOfItsFields(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		params := pinnedThreadParams()
		params.Messages = []BetaMessageParam{userMessage("Hello"), assistantMessage("Hi there.")}
		thread, server := newThread(t, params, refusedReply("msg_a"), reply("msg_b", "Reply B"))
		refused := turn("Refused")
		refused.System = []BetaTextBlockParam{{Text: "Be verbose."}}
		refused.Tools = []BetaToolUnionParam{}

		require.NoError(t, sendOrStream(thread, stream, refused))
		require.Nil(t, thread.LastMessage())
		require.NoError(t, sendOrStream(thread, stream, turn("Again")))

		request := server.requireRequests(t, 2)[1]
		requireCreate(t, request, messagesOf(userJSON("Hello"), assistantJSON("Hi there."), userJSON("Again")))
		requireJSONEqual(t, request.field("system.0.text"), `"Be brief."`)
		requireJSONEqual(t, request.field("tools.0.name"), `"get_weather"`)
	})
}

func TestBetaMessageThread_Send_ARefusedTurnLeavesAppendedMessagesStaged(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), refusedReply("msg_b"), reply("msg_c", "Reply C"))
		send(t, thread, turn("Hello"))
		thread.Append(userMessage("Staged"))

		require.NoError(t, sendOrStream(thread, stream, turn("Refused")))
		require.NoError(t, sendOrStream(thread, stream, turn("Again")))

		requests := server.requireRequests(t, 3)
		requireContinue(t, requests[1], "msg_a", messagesOf(userJSON("Staged"), userJSON("Refused")))
		requireCreate(t, requests[2], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Staged"), userJSON("Again")))
	})
}

func TestBetaMessageThread_Send_ARefusedTurnWithoutAThreadIsLeftOutOfTheNextRequest(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, continueUnthreadedParams(), unsupportedRequest(`["inference_hooks"]`), reply("msg_a", "Reply A"), refusedReply("msg_b"), reply("msg_c", "Reply C"))
		send(t, thread, turn("Hello"))

		require.NoError(t, sendOrStream(thread, stream, turn("Refused")))
		require.NoError(t, sendOrStream(thread, stream, turn("Again")))

		requireUnthreaded(t, server.requireRequests(t, 4)[3], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again")))
	})
}

// Fallback models.

func TestBetaMessageThread_Send_AReplyWrittenByAFallbackModelStartsANewThreadAndChangesNoField(t *testing.T) {
	for name, fallbacks := range map[string]BetaFallbacksParamUnion{
		"named fallbacks":   fallbackParams().Fallbacks,
		"default fallbacks": BetaFallbacksParamOfDefault(),
	} {
		t.Run(name, func(t *testing.T) {
			forEachRunner(t, func(t *testing.T, stream bool) {
				params := fallbackParams()
				params.Fallbacks = fallbacks
				params.System, params.Tools = pinnedThreadParams().System, pinnedThreadParams().Tools
				thread, server := newThread(t, params,
					writtenByFallback(reply("msg_a", "Reply A")), writtenByFallback(reply("msg_b", "Reply B")), reply("msg_c", "Reply C"), reply("msg_d", "Reply D"))

				require.NoError(t, sendOrStream(thread, stream, turn("Hello")))
				require.Equal(t, "", thread.PreviousMessageID())
				require.Equal(t, "msg_a", thread.LastMessage().ID)
				for _, text := range []string{"Again", "Once more", "Last"} {
					require.NoError(t, sendOrStream(thread, stream, turn(text)))
				}

				requests := server.requireRequests(t, 4)
				history := []string{userJSON("Hello"), fallbackAssistantJSON("Reply A"), userJSON("Again"), fallbackAssistantJSON("Reply B"), userJSON("Once more")}
				requireCreate(t, requests[1], messagesOf(history[:3]...))
				requireCreate(t, requests[2], messagesOf(history...))
				requireContinue(t, requests[3], "msg_c", messagesOf(userJSON("Last")))
				require.NotEmpty(t, requests[0].field("fallbacks"))
				for _, request := range requests {
					requireJSONEqual(t, request.field("model"), `"claude-sonnet-4-5"`)
					requireJSONEqual(t, request.field("fallbacks"), requests[0].field("fallbacks"))
					require.Equal(t, []string{threadsBeta, fallbackBeta}, request.betas())
				}
				otherFields := func(request threadRequest) string {
					body, _ := sjson.DeleteBytes(request.body, "messages")
					return string(body)
				}
				requireJSONEqual(t, otherFields(requests[1]), otherFields(requests[0]))
				requireJSONEqual(t, otherFields(requests[2]), otherFields(requests[0]))
				params.Messages = thread.Params.Messages
				require.Equal(t, params, thread.Params)
			})
		})
	}
}

func TestBetaMessageThread_SendStreaming_ReadsTheFallbackEntryFromMessageDelta(t *testing.T) {
	usage := gjson.Get(writtenByFallback(reply("msg_a", "Reply A")).body, "usage").Raw
	events := messageStartEvent("msg_a") + textEvents(0, "Reply A") +
		threadEvent("message_delta", `{"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null},"usage":`+usage+`}`) +
		messageStopEvent()
	thread, server := newThread(t, fallbackParams(), threadResponse{stream: true, body: events}, reply("msg_b", "Reply B"))
	readStream(thread, turn("Hello"))

	send(t, thread, turn("Again"))

	request := server.requireRequests(t, 2)[1]
	requireCreate(t, request, messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again")))
	requireJSONEqual(t, request.field("model"), `"claude-sonnet-4-5"`)
	requireJSONEqual(t, request.field("fallbacks"), `[{"model":"`+fallbackModel+`"}]`)
}

func TestBetaMessageThread_Send_AThreadStartedAfterAFallbackReplyIsNotTakenForTheReplayOfALostOne(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, fallbackParams(), writtenByFallback(reply("msg_a", "Reply A")), reply("msg_b", "Reply B"), notFound(), reply("msg_c", "Reply C"))
		require.NoError(t, sendOrStream(thread, stream, turn("Hello")))
		require.NoError(t, sendOrStream(thread, stream, turn("Again")))

		require.NoError(t, sendOrStream(thread, stream, turn("Once more")))

		requests := server.requireRequests(t, 4)
		requireContinue(t, requests[2], "msg_b", messagesOf(userJSON("Once more")))
		requireJSONEqual(t, requests[3].field("thread"), createJSON)
	})
}

func TestBetaMessageThread_Send_ARefusedReplyOfAFallbackModelLeavesTheModelAsItWas(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, fallbackParams(), writtenByFallback(refusedReply("msg_a")), reply("msg_b", "Reply B"))

		require.NoError(t, sendOrStream(thread, stream, turn("Refused")))
		require.NoError(t, sendOrStream(thread, stream, turn("Again")))

		request := server.requireRequests(t, 2)[1]
		requireCreate(t, request, messagesOf(userJSON("Again")))
		requireJSONEqual(t, request.field("model"), `"claude-sonnet-4-5"`)
		requireJSONEqual(t, request.field("fallbacks"), `[{"model":"`+fallbackModel+`"}]`)
	})
}

func TestBetaMessageThread_Send_AReplyOfTheRequestedModelKeepsTheThreadAndItsFallbacks(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, fallbackParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))

		require.NoError(t, sendOrStream(thread, stream, turn("Hello")))
		require.NoError(t, sendOrStream(thread, stream, turn("Again")))

		request := server.requireRequests(t, 2)[1]
		requireContinue(t, request, "msg_a", messagesOf(userJSON("Again")))
		requireJSONEqual(t, request.field("model"), `"claude-sonnet-4-5"`)
		requireJSONEqual(t, request.field("fallbacks"), `[{"model":"`+fallbackModel+`"}]`)
	})
}

func TestBetaMessageThread_Send_WithoutAThreadAReplyWrittenByAFallbackModelLeavesTheModelAsItWas(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, fallbackParams(), unsupportedRequest(`["inference_hooks"]`), writtenByFallback(reply("msg_a", "Reply A")), reply("msg_b", "Reply B"))

		require.NoError(t, sendOrStream(thread, stream, turn("Hello")))
		require.NoError(t, sendOrStream(thread, stream, turn("Again")))

		request := server.requireRequests(t, 3)[2]
		requireUnthreaded(t, request, messagesOf(userJSON("Hello"), fallbackAssistantJSON("Reply A"), userJSON("Again")))
		requireJSONEqual(t, request.field("model"), `"claude-sonnet-4-5"`)
		requireJSONEqual(t, request.field("fallbacks"), `[{"model":"`+fallbackModel+`"}]`)
		require.Equal(t, Model("claude-sonnet-4-5"), thread.Params.Model)
	})
}

// Cache breakpoints.

func TestBetaMessageThread_Send_AReplayDropsTheCacheControlOfTheRecordedMessages(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, threadParams(),
			reply("msg_1", "Reply 1"), reply("msg_2", "Reply 2"), reply("msg_3", "Reply 3"), reply("msg_4", "Reply 4"), reply("msg_5", "Reply 5"),
			threadRejection(http.StatusNotFound, "thread_not_found"), reply("msg_6", "Reply 6"))
		given := make([]BetaMessageParam, 6)
		var sent, history []string
		for i := range given {
			text := fmt.Sprintf("Turn %d", i+1)
			given[i] = markedUserMessage(text)
			sent = append(sent, userJSON(text), assistantJSON(fmt.Sprintf("Reply %d", i+1)))
			history = append(history, markedUserJSON(text), assistantJSON(fmt.Sprintf("Reply %d", i+1)))
		}

		for _, message := range given {
			require.NoError(t, sendOrStream(thread, stream, BetaMessageThreadSendParams{Messages: []BetaMessageParam{message}}))
		}

		requests := server.requireRequests(t, 7)
		requireContinue(t, requests[4], "msg_4", messagesOf(markedUserJSON("Turn 5")))
		requireContinue(t, requests[5], "msg_5", messagesOf(markedUserJSON("Turn 6")))
		requireCreate(t, requests[6], messagesOf(append(sent[:10:10], markedUserJSON("Turn 6"))...))
		requireMessages(t, thread.Messages(), messagesOf(history...))
		requireMessages(t, thread.Params.Messages, messagesOf(history...))
		for i, message := range given {
			requireMessages(t, []BetaMessageParam{message}, messagesOf(history[2*i]))
		}
	})
}

func TestBetaMessageThread_Send_TheFirstCreateCarriesTheCacheControlOfTheSeed(t *testing.T) {
	params := threadParams()
	params.Messages = []BetaMessageParam{markedUserMessage("Hello"), assistantMessage("Hi there.")}
	thread, server := newThread(t, params, reply("msg_a", "Reply A"))

	send(t, thread, markedTurn("Again"))

	requireCreate(t, server.requests()[0], messagesOf(markedUserJSON("Hello"), assistantJSON("Hi there."), markedUserJSON("Again")))
}

func TestBetaMessageThread_Send_AReplayCarriesTheCacheControlOfAMessageThatWasNeverAnswered(t *testing.T) {
	const unanswered = `{"role":"user","content":[{"type":"text","text":"Again","cache_control":{"type":"ephemeral"}},{"type":"text","text":"and again","cache_control":{"type":"ephemeral"}}]}`
	for name, test := range map[string]struct {
		next BetaMessageThreadSendParams
		sent []string
	}{
		"with no messages": {next: noMessages},
		"with a message":   {next: markedTurn("Once more"), sent: []string{markedUserJSON("Once more")}},
	} {
		t.Run(name, func(t *testing.T) {
			thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), cutOffStream(messageStartEvent("msg_b")), reply("msg_c", "Reply C"))
			send(t, thread, markedTurn("Hello"))
			readStream(thread, BetaMessageThreadSendParams{Messages: []BetaMessageParam{NewBetaUserMessage(
				markedUserMessage("Again").OfUser.Content[0],
				markedUserMessage("and again").OfUser.Content[0],
			)}})

			send(t, thread, test.next)

			requireCreate(t, server.requireRequests(t, 3)[2], messagesOf(append([]string{userJSON("Hello"), assistantJSON("Reply A"), unanswered}, test.sent...)...))
		})
	}
}

func TestBetaMessageThread_Send_AReplayDropsTheCacheControlOfNestedBlocks(t *testing.T) {
	marker := NewBetaCacheControlEphemeralParam()
	markedText := BetaTextBlockParam{Text: "ok", CacheControl: marker}
	markedAddition := BetaRequestToolAdditionBlockParam{
		Tool: BetaRequestToolAdditionBlockToolUnionParam{OfToolDefinition: &BetaToolChangeToolDefinitionParam{Definition: BetaToolUnionParam{OfTool: &BetaToolParam{
			Name:          "lookup",
			InputSchema:   BetaToolInputSchemaParam{Properties: map[string]any{"cache_control": map[string]any{"type": "string"}}},
			InputExamples: []map[string]any{{"cache_control": "no-store"}},
			CacheControl:  marker,
		}}}},
		CacheControl: marker,
	}
	const definitionFieldsJSON = `"name":"lookup","input_schema":{"type":"object","properties":{"cache_control":{"type":"string"}}},"input_examples":[{"cache_control":"no-store"}]`
	for name, test := range map[string]struct {
		message BetaMessageParam
		want    string
	}{
		"in a tool_result": {
			message: NewBetaUserMessage(BetaContentBlockParamUnion{OfToolResult: &BetaToolResultBlockParam{
				ToolUseID:    "a",
				Content:      []BetaToolResultBlockParamContentUnion{{OfText: &markedText}},
				CacheControl: marker,
			}}),
			want: toolResultMessageJSON(toolResultBlockJSON("a", "ok")),
		},
		"in a search_result": {
			message: NewBetaUserMessage(BetaContentBlockParamUnion{OfSearchResult: &BetaSearchResultBlockParam{
				Content: []BetaTextBlockParam{markedText},
				Source:  "https://example.com",
				Title:   "Example",
			}}),
			want: `{"role":"user","content":[{"type":"search_result","content":[{"type":"text","text":"ok"}],"source":"https://example.com","title":"Example"}]}`,
		},
		"in the source of a document": {
			message: NewBetaUserMessage(BetaContentBlockParamUnion{OfDocument: &BetaRequestDocumentBlockParam{
				Source: BetaRequestDocumentBlockSourceUnionParam{OfContent: &BetaContentBlockSourceParam{
					Content: param.Override[BetaContentBlockSourceContentUnionParam]([]BetaTextBlockParam{markedText}),
				}},
				CacheControl: marker,
			}}),
			want: `{"role":"user","content":[{"type":"document","source":{"type":"content","content":[{"type":"text","text":"ok"}]}}]}`,
		},
		"in the tool_references of a tool_search_tool_result": {
			message: BetaMessageParam{OfAssistant: &BetaAssistantMessageParam{Content: []BetaContentBlockParamUnion{{OfToolSearchToolResult: &BetaToolSearchToolResultBlockParam{
				ToolUseID: "a",
				Content: BetaToolSearchToolResultBlockParamContentUnion{OfRequestToolSearchToolSearchResultBlock: &BetaToolSearchToolSearchResultBlockParam{
					ToolReferences: []BetaToolReferenceBlockParam{{ToolName: "lookup", CacheControl: marker}},
				}},
			}}}}},
			want: `{"role":"assistant","content":[{"type":"tool_search_tool_result","tool_use_id":"a","content":{"type":"tool_search_tool_search_result","tool_references":[{"type":"tool_reference","tool_name":"lookup"}]}}]}`,
		},
		"in the definition of a tool_addition, beside its input_schema and input_examples": {
			message: systemToolChange(BetaContentBlockParamUnion{OfToolAddition: &markedAddition}),
			want:    toolChangeMessageJSON(toolAdditionJSON(`{` + definitionFieldsJSON + `}`)),
		},
		"beside the tool_changes of a compaction": {
			message: BetaMessageParam{OfAssistant: &BetaAssistantMessageParam{Content: []BetaContentBlockParamUnion{{OfCompaction: &BetaCompactionBlockParam{
				Content:      String("Summary"),
				ToolChanges:  []BetaCompactionBlockParamToolChangeUnion{{OfToolAddition: &markedAddition}},
				CacheControl: marker,
			}}}}},
			want: `{"role":"assistant","content":[{"type":"compaction","content":"Summary","tool_changes":[{"type":"tool_addition","cache_control":{"type":"ephemeral"},"tool":{"type":"tool_definition","definition":{` + definitionFieldsJSON + `,"cache_control":{"type":"ephemeral"}}}}]}]}`,
		},
		"beside the input of a tool_use": {
			message: BetaMessageParam{OfAssistant: &BetaAssistantMessageParam{Content: []BetaContentBlockParamUnion{{OfToolUse: &BetaToolUseBlockParam{
				ID:           "a",
				Name:         "lookup",
				Input:        map[string]any{"cache_control": "no-store"},
				CacheControl: marker,
			}}}}},
			want: `{"role":"assistant","content":[{"type":"tool_use","id":"a","name":"lookup","input":{"cache_control":"no-store"}}]}`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), threadRejection(http.StatusNotFound, "thread_not_found"), reply("msg_b", "Reply B"))
			send(t, thread, BetaMessageThreadSendParams{Messages: []BetaMessageParam{test.message}})

			send(t, thread, turn("Again"))

			requireCreate(t, server.requireRequests(t, 3)[2], messagesOf(test.want, assistantJSON("Reply A"), userJSON("Again")))
		})
	}
}

func TestBetaMessageThread_Send_AReplayKeepsTheCacheControlOutsideTheMessages(t *testing.T) {
	params := pinnedThreadParams()
	params.System[0].CacheControl = NewBetaCacheControlEphemeralParam()
	params.Tools[0].OfTool.CacheControl = NewBetaCacheControlEphemeralParam()
	params.CacheControl = NewBetaCacheControlEphemeralParam()
	thread, server := newThread(t, params, reply("msg_a", "Reply A"), threadRejection(http.StatusNotFound, "thread_not_found"), reply("msg_b", "Reply B"))
	send(t, thread, markedTurn("Hello"))

	send(t, thread, turn("Again"))

	request := server.requireRequests(t, 3)[2]
	requireCreate(t, request, messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again")))
	for _, field := range []string{"system.0.cache_control", "tools.0.cache_control", "cache_control"} {
		requireJSONEqual(t, request.field(field), `{"type":"ephemeral"}`)
	}
}

func TestBetaMessageThread_Send_WithoutAThreadDropsTheCacheControlOfTheRecordedMessages(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, continueUnthreadedParams(), unsupportedRequest(`["inference_hooks"]`), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))

		require.NoError(t, sendOrStream(thread, stream, markedTurn("Hello")))
		require.NoError(t, sendOrStream(thread, stream, markedTurn("Again")))

		requests := server.requireRequests(t, 3)
		requireUnthreaded(t, requests[1], messagesOf(markedUserJSON("Hello")))
		requireUnthreaded(t, requests[2], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), markedUserJSON("Again")))
	})
}

func TestBetaMessageThread_Fork_DropsTheCacheControlOfTheRecordedMessagesAndKeepsThatOfTheStaged(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))
	send(t, thread, markedTurn("Hello"))
	thread.Append(markedUserMessage("Staged"))
	fork := thread.Fork()

	send(t, fork, markedTurn("Fork"))

	requireCreate(t, server.requireRequests(t, 2)[1], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), markedUserJSON("Staged"), markedUserJSON("Fork")))
	requireMessages(t, fork.Messages(), messagesOf(markedUserJSON("Hello"), assistantJSON("Reply A"), markedUserJSON("Staged"), markedUserJSON("Fork"), assistantJSON("Reply B")))
}

// Unfinished turns.

func TestBetaMessageThread_Send_WithNoMessagesResumesAPausedTurn(t *testing.T) {
	thread, server := newThread(t, threadParams(), pausedReply("msg_a", "Searching"), reply("msg_b", "Reply B"))
	send(t, thread, turn("Hello"))

	send(t, thread, noMessages)

	requireContinue(t, server.requests()[1], "msg_a", `[]`)
}

// Container.

func TestBetaMessageThread_Send_SendsTheContainerOfTheReplyOnLaterTurns(t *testing.T) {
	thread, server := newThread(t, threadParams(), replyInContainer("msg_a", "Reply A", "container_a"), reply("msg_b", "Reply B"))

	send(t, thread, turn("Hello"))
	send(t, thread, turn("Again"))

	require.Empty(t, server.requests()[0].field("container"))
	requireJSONEqual(t, server.requests()[1].field("container"), `"container_a"`)
}

func TestBetaMessageThread_Send_KeepsTheContainerTheCallerGave(t *testing.T) {
	params := threadParams()
	params.Container = BetaMessageNewParamsContainerUnion{OfString: String("container_mine")}
	thread, server := newThread(t, params, replyInContainer("msg_a", "Reply A", "container_a"), reply("msg_b", "Reply B"))

	send(t, thread, turn("Hello"))
	send(t, thread, turn("Again"))

	requireJSONEqual(t, server.requests()[0].field("container"), `"container_mine"`)
	requireJSONEqual(t, server.requests()[1].field("container"), `"container_mine"`)
}

// Fork.

func TestBetaMessageThread_Fork_LeavesTheOriginalContinuingFromItsOwnReply(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"), reply("msg_c", "Reply C"))
	send(t, thread, turn("Hello"))
	fork := thread.Fork()
	send(t, fork, turn("Fork"))

	send(t, thread, turn("Original"))

	requireContinue(t, server.requests()[2], "msg_a", messagesOf(userJSON("Original")))
	requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Original"), assistantJSON("Reply C")))
	requireMessages(t, fork.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Fork"), assistantJSON("Reply B")))
}

func TestBetaMessageThread_Fork_AMessageChangedInPlaceOnTheForkLeavesTheOriginalContinuing(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))
	send(t, thread, turn("Hello"))

	thread.Fork().Params.Messages[0].OfUser.Content = userMessage("Other").OfUser.Content
	send(t, thread, turn("Original"))

	requireContinue(t, server.requireRequests(t, 2)[1], "msg_a", messagesOf(userJSON("Original")))
	requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Original"), assistantJSON("Reply B")))
}

func TestBetaMessageThread_Fork_KeepsTheUnsupportedRequestBehavior(t *testing.T) {
	thread, server := newThread(t, continueUnthreadedParams(), unsupportedRequest(`["inference_hooks"]`), reply("msg_a", "Reply A"))
	fork := thread.Fork()

	send(t, fork, turn("Hello"))

	requireCreate(t, server.requests()[0], messagesOf(userJSON("Hello")))
	requireUnthreaded(t, server.requests()[1], messagesOf(userJSON("Hello")))
}

func TestBetaMessageThread_Fork_OfAnUnthreadedThreadStaysUnthreaded(t *testing.T) {
	thread, server := newThread(t, continueUnthreadedParams(), unsupportedRequest(`["inference_hooks"]`), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))
	send(t, thread, turn("Hello"))

	fork := thread.Fork()
	send(t, fork, turn("Fork"))

	requireUnthreaded(t, server.requests()[2], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Fork")))
}

func TestBetaMessageThread_Fork_OfAFinishedThreadSendsTheWholeHistory(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))
		send(t, thread, turn("Hello"))
		fork := thread.Fork()

		require.NoError(t, sendOrStream(fork, stream, noMessages))

		requireCreate(t, server.requests()[1], messagesOf(userJSON("Hello"), assistantJSON("Reply A")))
	})
}

func TestBetaMessageThread_Fork_SendsWhatTheOriginalHadStaged(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), reply("msg_b", "Reply B"))
	send(t, thread, turn("Hello"))
	thread.Append(userMessage("Again"))
	fork := thread.Fork()

	send(t, fork, noMessages)

	requireCreate(t, server.requests()[1], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again")))
	requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again")))
}

// Structured output.

type threadForecast struct {
	Forecast string `json:"forecast"`
}

func forecastReply(id, forecast string) threadResponse {
	return reply(id, fmt.Sprintf(`{"forecast": %q}`, forecast))
}

func TestBetaMessageThread_Send_SendsTheOutputFormatsAsTheyWereGiven(t *testing.T) {
	var forecast threadForecast
	params := threadParams()
	params.OutputFormat = BetaJSONOutputFormatParam{Schema: &threadForecast{}}
	params.OutputConfig.Format = BetaJSONOutputFormatParam{Schema: &forecast}
	thread, server := newThread(t, params, forecastReply("msg_a", "Sunny"))

	send(t, thread, turn("Today?"))

	request := server.requests()[0]
	requireJSONEqual(t, request.field("output_format.type"), `"json_schema"`)
	requireJSONEqual(t, request.field("output_config.format.type"), `"json_schema"`)
	require.Equal(t, threadForecast{Forecast: "Sunny"}, forecast)
}

func TestBetaMessageThread_Send_RecordsAReplyWhoseStructuredOutputDoesNotParse(t *testing.T) {
	params := threadParams()
	params.OutputConfig.Format = BetaJSONOutputFormatParam{Schema: &threadForecast{}}
	thread, _ := newThread(t, params, reply("msg_a", "not JSON"))

	message, err := thread.Send(context.Background(), turn("Today?"))

	require.ErrorIs(t, err, ErrStructuredOutputParse)
	require.Equal(t, "msg_a", message.ID)
	require.Equal(t, "msg_a", thread.PreviousMessageID())
	requireMessages(t, thread.Messages(), messagesOf(userJSON("Today?"), assistantJSON("not JSON")))
}

// Provider clients.

func TestBetaMessageThread_Send_OnAClientWithoutMessageThreadsFailsLocally(t *testing.T) {
	server := newThreadServer(t)
	client := server.client(requestconfig.WithoutMessageThreads())
	thread := client.Beta.Messages.NewThread(threadParams())

	_, err := thread.Send(context.Background(), turn("Hello"))

	requireLocalError(t, err)
	require.Empty(t, server.requests())
}

// Streams.

func threadEvent(event, data string) string {
	return "event: " + event + "\ndata: " + data + "\n\n"
}

func messageStartEvent(id string) string {
	return threadEvent("message_start", fmt.Sprintf(`{"type":"message_start","message":{"id":%q,"type":"message","role":"assistant","model":"claude-sonnet-4-5","content":[],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}}`, id))
}

func textStartEvent(index int) string {
	return threadEvent("content_block_start", fmt.Sprintf(`{"type":"content_block_start","index":%d,"content_block":{"type":"text","text":""}}`, index))
}

func textDeltaEvent(index int, text string) string {
	return threadEvent("content_block_delta", fmt.Sprintf(`{"type":"content_block_delta","index":%d,"delta":{"type":"text_delta","text":%q}}`, index, text))
}

func toolUseStartEvent(index int, id, name string) string {
	return threadEvent("content_block_start", fmt.Sprintf(`{"type":"content_block_start","index":%d,"content_block":{"type":"tool_use","id":%q,"name":%q,"input":{}}}`, index, id, name))
}

func inputDeltaEvent(index int, partialJSON string) string {
	return threadEvent("content_block_delta", fmt.Sprintf(`{"type":"content_block_delta","index":%d,"delta":{"type":"input_json_delta","partial_json":%q}}`, index, partialJSON))
}

func blockStopEvent(index int) string {
	return threadEvent("content_block_stop", fmt.Sprintf(`{"type":"content_block_stop","index":%d}`, index))
}

func messageDeltaEvent(stopReason string) string {
	return threadEvent("message_delta", fmt.Sprintf(`{"type":"message_delta","delta":{"stop_reason":%q,"stop_sequence":null},"usage":{"output_tokens":1}}`, stopReason))
}

func messageStopEvent() string {
	return threadEvent("message_stop", `{"type":"message_stop"}`)
}

// textEvents returns the events of a finished text block.
func textEvents(index int, text string) string {
	return textStartEvent(index) + textDeltaEvent(index, text) + blockStopEvent(index)
}

// streamedReply is a whole streamed reply with one text block.
func streamedReply(id, text string) threadResponse {
	return streamedReplyThatStopped(id, text, "end_turn")
}

func streamedReplyThatStopped(id, text, stopReason string) threadResponse {
	return threadResponse{stream: true, body: messageStartEvent(id) + textEvents(0, text) + messageDeltaEvent(stopReason) + messageStopEvent()}
}

// cutOffStream is a stream whose connection drops after events.
func cutOffStream(events ...string) threadResponse {
	return threadResponse{stream: true, cutOff: true, body: strings.Join(events, "")}
}

// heldStream is a stream that stays open after events until the caller hangs
// up.
func heldStream(events ...string) threadResponse {
	return threadResponse{stream: true, held: true, body: strings.Join(events, "")}
}

// textCutOff is a stream that drops inside the text block of reply id, after
// text.
func textCutOff(id, text string) threadResponse {
	return cutOffStream(messageStartEvent(id), textStartEvent(0), textDeltaEvent(0, text))
}

// read reads stream to its end and closes it. It ignores how the stream ended.
func read(stream *BetaMessageThreadStream) {
	for stream.Next() {
	}
	_ = stream.Close()
}

// readStream sends a turn as a stream and reads it to its end.
func readStream(thread *BetaMessageThread, params BetaMessageThreadSendParams, opts ...option.RequestOption) {
	read(thread.SendStreaming(context.Background(), params, opts...))
}

// readUntilText reads stream until its first text delta.
func readUntilText(t *testing.T, stream *BetaMessageThreadStream) {
	t.Helper()
	for stream.Next() {
		if stream.Current().Delta.Text != "" {
			return
		}
	}
	t.Fatalf("the stream ended before any text: %v", stream.Err())
}

func TestBetaMessageThread_SendStreaming_FirstTurnCreatesTheThreadAndAsksForAStream(t *testing.T) {
	thread, server := newThread(t, threadParams(), streamedReply("msg_a", "Reply A"))

	readStream(thread, turn("Hello"))

	request := server.requests()[0]
	requireCreate(t, request, messagesOf(userJSON("Hello")))
	requireJSONEqual(t, request.field("stream"), `true`)
	require.Equal(t, []string{threadsBeta}, request.betas())
}

func TestBetaMessageThread_SendStreaming_ForwardsEveryEvent(t *testing.T) {
	thread, _ := newThread(t, threadParams(), streamedReply("msg_a", "Reply A"))
	stream := thread.SendStreaming(context.Background(), turn("Hello"))
	defer stream.Close()
	types := []string{}

	for stream.Next() {
		types = append(types, stream.Current().Type)
	}

	require.NoError(t, stream.Err())
	require.Equal(t, []string{"message_start", "content_block_start", "content_block_delta", "content_block_stop", "message_delta", "message_stop"}, types)
}

func TestBetaMessageThread_SendStreaming_RecordsTheReplyAtMessageStop(t *testing.T) {
	thread, _ := newThread(t, threadParams(), streamedReply("msg_a", "Reply A"))

	readStream(thread, turn("Hello"))

	require.Equal(t, "msg_a", thread.PreviousMessageID())
	requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A")))
}

func TestBetaMessageThread_SendStreaming_ClosingAfterMessageStopKeepsTheWholeReply(t *testing.T) {
	thread, _ := newThread(t, threadParams(), streamedReply("msg_a", "Reply A"))
	stream := thread.SendStreaming(context.Background(), turn("Hello"))
	for stream.Next() && stream.Current().Type != "message_stop" {
	}

	err := stream.Close()

	require.NoError(t, err)
	require.Equal(t, "msg_a", thread.PreviousMessageID())
}

func TestBetaMessageThread_SendStreaming_RejectedContinueIsReplayedUnderCreateInTheSameStream(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), threadRejection(http.StatusNotFound, "thread_not_found"), streamedReply("msg_b", "Reply B"))
	send(t, thread, turn("Hello"))

	readStream(thread, turn("Again"))

	replay := server.requests()[2]
	requireCreate(t, replay, messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Again")))
	requireJSONEqual(t, replay.field("stream"), `true`)
	require.Equal(t, "msg_b", thread.PreviousMessageID())
}

func TestBetaMessageThread_Send_AfterAFailedStreamThatWasNeverReadIsNotRefused(t *testing.T) {
	thread, _ := newThread(t, threadParams(), plainError(http.StatusBadRequest, "bad request"), reply("msg_a", "Reply A"))
	_ = thread.SendStreaming(context.Background(), turn("Hello"))

	_, err := thread.Send(context.Background(), noMessages)

	require.NoError(t, err)
}

func TestBetaMessageThread_LastMessage_AfterACutOffStreamIsThePartialReply(t *testing.T) {
	thread, _ := newThread(t, threadParams(), reply("msg_a", "Reply A"), textCutOff("msg_b", "one"))
	send(t, thread, turn("Hello"))

	readStream(thread, turn("Count"))

	require.Equal(t, "msg_b", thread.LastMessage().ID)
	require.Equal(t, "one", thread.LastMessage().Content[0].Text)
}

func TestBetaMessageThread_SendStreaming_ClosedBeforeMessageStopRecordsWhatWasRead(t *testing.T) {
	held := heldStream(messageStartEvent("msg_b"), textStartEvent(0), textDeltaEvent(0, "one"))
	thread, _ := newThread(t, threadParams(), reply("msg_a", "Reply A"), held)
	send(t, thread, turn("Hello"))
	stream := thread.SendStreaming(context.Background(), turn("Count"))
	readUntilText(t, stream)

	err := stream.Close()

	require.NoError(t, err)
	require.Empty(t, thread.PreviousMessageID())
	requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Count"), assistantJSON("one")))
}

func TestBetaMessageThread_SendStreaming_PartialReply(t *testing.T) {
	tests := map[string]struct {
		// events follow message_start, and the connection drops after them.
		events   string
		recorded string
	}{
		"keeps a block that finished": {
			textEvents(0, "Let me check.") + textStartEvent(1) + textDeltaEvent(1, "It is"),
			partialJSON(textBlockJSON("Let me check."), textBlockJSON("It is")),
		},
		"drops an unfinished tool call": {
			textEvents(0, "Let me check.") + toolUseStartEvent(1, "toolu_a", "get_weather") + inputDeltaEvent(1, `{"city": "Pa`),
			partialJSON(textBlockJSON("Let me check.")),
		},
		"drops a text block without text": {
			textEvents(0, "Let me check.") + textStartEvent(1),
			partialJSON(textBlockJSON("Let me check.")),
		},
		"keeps a tool call that finished": {
			finishedToolCall(0, "toolu_a") + messageDeltaEvent("tool_use"),
			partialJSON(toolCallJSON("toolu_a")),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			thread, _ := newThread(t, threadParams(), reply("msg_a", "Reply A"), cutOffStream(messageStartEvent("msg_b"), test.events))
			send(t, thread, turn("Hello"))

			readStream(thread, turn("Weather?"))

			requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Weather?"), test.recorded))
		})
	}
}

// finishedToolCall returns the events of a tool call that finished.
func finishedToolCall(index int, id string) string {
	return toolUseStartEvent(index, id, "get_weather") + inputDeltaEvent(index, `{"city": "Paris"}`) + blockStopEvent(index)
}

func toolCallJSON(id string) string {
	return fmt.Sprintf(`{"type":"tool_use","id":%q,"name":"get_weather","input":{"city":"Paris"}}`, id)
}

func partialJSON(blocks ...string) string {
	return `{"role":"assistant","content":` + messagesOf(blocks...) + `}`
}

func textBlockJSON(text string) string {
	return fmt.Sprintf(`{"type":"text","text":%q}`, text)
}

func TestBetaMessageThread_Send_WithNoMessagesAfterACutOffStreamSendsTheTurnAgainWithoutThePartialReply(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), textCutOff("msg_b", "one"), reply("msg_c", "Reply C"))
	send(t, thread, turn("Hello"))
	readStream(thread, turn("Count"))

	send(t, thread, noMessages)

	requireCreate(t, server.requests()[2], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Count")))
	requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Count"), assistantJSON("Reply C")))
}

func TestBetaMessageThread_Send_AMessageAfterACutOffStreamStartsANewThreadWithThePartialReply(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), textCutOff("msg_b", "one"), reply("msg_c", "Reply C"), reply("msg_d", "Reply D"))
		send(t, thread, turn("Hello"))
		readStream(thread, turn("Count"))

		require.NoError(t, sendOrStream(thread, stream, turn("Never mind")))
		send(t, thread, turn("Thanks"))

		requests := server.requireRequests(t, 4)
		requireCreate(t, requests[2], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Count"), assistantJSON("one"), userJSON("Never mind")))
		requireContinue(t, requests[3], "msg_c", messagesOf(userJSON("Thanks")))
	})
}

func TestBetaMessageThread_SendStreaming_CutOffBeforeAnyContentAddsNoReply(t *testing.T) {
	thread, _ := newThread(t, threadParams(), reply("msg_a", "Reply A"), cutOffStream(messageStartEvent("msg_b")))
	send(t, thread, turn("Hello"))

	readStream(thread, turn("Count"))

	require.Empty(t, thread.PreviousMessageID())
	requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Count")))
}

func TestBetaMessageThread_ACallThatFailsBeforeAReplyStartsTakesBackItsMessages(t *testing.T) {
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	failures := map[string]struct {
		// answers is how the server answers the call. A call with none sends nothing.
		answers []threadResponse
		ctx     context.Context
		system  []BetaTextBlockParam
		opts    []option.RequestOption
	}{
		"rejected":           {answers: []threadResponse{plainError(http.StatusBadRequest, "bad request")}},
		"overloaded":         {answers: []threadResponse{plainError(529, "overloaded")}},
		"connection dropped": {answers: []threadResponse{cutOffStream()}},
		"no connection": {opts: []option.RequestOption{option.WithMiddleware(func(*http.Request, option.MiddlewareNext) (*http.Response, error) {
			return nil, errors.New("connection refused")
		})}},
		"cancelled context": {ctx: cancelled},
		"changed System":    {system: []BetaTextBlockParam{{Text: "Be verbose."}}},
	}
	for name, failure := range failures {
		t.Run(name, func(t *testing.T) {
			forEachRunner(t, func(t *testing.T, stream bool) {
				script := slices.Concat(
					[]threadResponse{reply("msg_a", "Reply A")},
					failure.answers, []threadResponse{reply("msg_b", "Reply B")},
					failure.answers, []threadResponse{reply("msg_c", "Reply C")},
				)
				thread, server := newThread(t, threadParams(), script...)
				fail := func(text string) {
					t.Helper()
					params, ctx := turn(text), context.Background()
					params.System = failure.system
					if failure.ctx != nil {
						ctx = failure.ctx
					}
					require.Error(t, sendOrStreamContext(ctx, thread, stream, params, failure.opts...))
				}
				send(t, thread, turn("Hello"))

				fail("Question")
				require.Equal(t, "msg_a", thread.PreviousMessageID())
				requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello"), assistantJSON("Reply A")))
				send(t, thread, turn("Question"))
				thread.Append(userMessage("Staged"))
				fail("Never mind")
				send(t, thread, noMessages)

				requests := server.requireRequests(t, len(script))
				requireContinue(t, requests[1+len(failure.answers)], "msg_a", messagesOf(userJSON("Question")))
				requireContinue(t, requests[len(script)-1], "msg_b", messagesOf(userJSON("Staged")))
				requireMessages(t, thread.Messages(), messagesOf(
					userJSON("Hello"), assistantJSON("Reply A"),
					userJSON("Question"), assistantJSON("Reply B"),
					userJSON("Staged"), assistantJSON("Reply C"),
				))
			})
		})
	}
}

func TestBetaMessageThread_SendStreaming_ClosedBeforeMessageStartTakesBackItsMessage(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), heldStream(), reply("msg_b", "Reply B"))
	send(t, thread, turn("Hello"))
	require.NoError(t, thread.SendStreaming(context.Background(), turn("Never mind")).Close())

	send(t, thread, turn("Again"))

	requireContinue(t, server.requests()[2], "msg_a", messagesOf(userJSON("Again")))
}

func TestBetaMessageThread_SendStreaming_AMessageAppendedToAStreamThatFailsBeforeMessageStartStaysStaged(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), heldStream(), reply("msg_b", "Reply B"))
	send(t, thread, turn("Hello"))
	stream := thread.SendStreaming(context.Background(), turn("Never mind"))
	thread.Append(userMessage("Staged"))
	require.NoError(t, stream.Close())

	send(t, thread, noMessages)

	requireContinue(t, server.requests()[2], "msg_a", messagesOf(userJSON("Staged")))
}

func TestBetaMessageThread_Send_AfterAStreamCutOffAtMessageStartSendsItsMessageAgain(t *testing.T) {
	thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), cutOffStream(messageStartEvent("msg_b")), reply("msg_c", "Reply C"))
	send(t, thread, turn("Hello"))
	readStream(thread, turn("Count"))

	send(t, thread, noMessages)

	requireCreate(t, server.requests()[2], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Count")))
}

func TestBetaMessageThread_Send_ARejectedMessageAfterACutOffStreamLeavesThePartialReplyAsItWas(t *testing.T) {
	forEachRunner(t, func(t *testing.T, stream bool) {
		thread, server := newThread(t, threadParams(), reply("msg_a", "Reply A"), textCutOff("msg_b", "one"), plainError(http.StatusBadRequest, "bad request"), reply("msg_c", "Reply C"))
		send(t, thread, turn("Hello"))
		readStream(thread, turn("Count"))
		requireStatus(t, sendOrStream(thread, stream, turn("Never mind")), http.StatusBadRequest)

		send(t, thread, noMessages)

		requests := server.requireRequests(t, 4)
		requireCreate(t, requests[2], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Count"), assistantJSON("one"), userJSON("Never mind")))
		requireCreate(t, requests[3], messagesOf(userJSON("Hello"), assistantJSON("Reply A"), userJSON("Count")))
	})
}

func TestBetaMessageThread_Send_WhileAStreamIsOpenFailsLocally(t *testing.T) {
	thread, server := newThread(t, threadParams(), streamedReply("msg_a", "Reply A"))
	stream := thread.SendStreaming(context.Background(), turn("Hello"))
	defer stream.Close()

	_, err := thread.Send(context.Background(), turn("Interrupting"))

	requireLocalError(t, err)
	require.Len(t, server.requests(), 1)
	requireMessages(t, thread.Messages(), messagesOf(userJSON("Hello")))
}
