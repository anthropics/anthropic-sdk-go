package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"strings"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/stainlessheader"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/ssestream"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// BetaMessageThreadUnsupportedRequestBehavior is what a [BetaMessageThread] does when
// the API answers a turn with the error code thread_unsupported_request.
type BetaMessageThreadUnsupportedRequestBehavior string

const (
	// BetaMessageThreadUnsupportedRequestBehaviorError returns the error to the caller. It
	// is the default.
	BetaMessageThreadUnsupportedRequestBehaviorError BetaMessageThreadUnsupportedRequestBehavior = "error"
	// BetaMessageThreadUnsupportedRequestBehaviorContinueUnthreaded sends the turn again
	// with the whole history and no thread, and sends every later turn the same way.
	BetaMessageThreadUnsupportedRequestBehaviorContinueUnthreaded BetaMessageThreadUnsupportedRequestBehavior = "continue_unthreaded"
)

// BetaMessageThread is a conversation that the API keeps between requests, so
// that each turn sends only the messages that are new. Create one with
// [BetaMessageService.NewThread].
//
// A BetaMessageThread runs one turn at a time and is NOT safe for concurrent use.
type BetaMessageThread struct {
	// Params holds the request fields sent with every turn, and the whole history
	// in Messages. Assign to it between turns to change all later turns. System and
	// Tools cannot change once a turn has created the thread.
	Params BetaMessageThreadCreateParams
	// Options are applied to every request, after the client's options and
	// before the options of the turn.
	Options []option.RequestOption

	service *BetaMessageService
	// stored is what has been sent, so the API may hold it. The next turn sends
	// the messages of Params.Messages that follow it.
	stored []betaThreadEntry
	link   betaThreadLink
	inTurn bool
	// The turn in progress was given the givenCount messages of Params.Messages
	// from givenStart on. givenCount is 0 once the turn is recorded.
	givenStart, givenCount int
}

// betaThreadLink is the thread field of the next request: betaThreadCreate,
// betaThreadContinue or betaThreadNone.
type betaThreadLink interface{ isBetaThreadLink() }

type betaThreadCreate struct {
	// afterNotFound is set when the request replays a turn that the API answered
	// with thread_not_found.
	afterNotFound bool
}

type betaThreadContinue struct {
	previousMessageID string
	createdWith       betaThreadCreateOnlyFields
	// afterNotFound is set on the first request that continues such a replay.
	afterNotFound bool
}

// betaThreadNone sends no thread field, for the life of the thread.
type betaThreadNone struct{}

func (betaThreadCreate) isBetaThreadLink()   {}
func (betaThreadContinue) isBetaThreadLink() {}
func (betaThreadNone) isBetaThreadLink()     {}

// betaThreadCreateOnlyFields holds the request JSON of the fields that only a
// request that creates a thread may carry.
type betaThreadCreateOnlyFields struct {
	system string
	tools  string
}

// betaThreadEntry is one message of the stored history.
type betaThreadEntry struct {
	param BetaMessageParam
	// reply is set when the message came from the API rather than the caller.
	reply *BetaMessage
	// partial is set when the stream of reply ended before message_stop.
	partial bool
}

type betaThreadRequest struct {
	params  BetaMessageNewParams
	options []option.RequestOption
	// createOnly holds the create-only fields of the turn. A request that
	// continues a thread leaves them out, and they equal what it was created with.
	createOnly betaThreadCreateOnlyFields
}

// NewThread creates a [BetaMessageThread] and sends no request. params holds what
// [BetaMessageService.New] takes, except Thread.
func (r *BetaMessageService) NewThread(params BetaMessageThreadCreateParams, opts ...option.RequestOption) *BetaMessageThread {
	thread := &BetaMessageThread{
		Params:  params,
		Options: slices.Clone(opts),
		service: r,
		link:    betaThreadCreate{},
	}
	thread.Params.Messages = nil
	thread.Append(params.Messages...)
	return thread
}

// supportsMessageThreads reports whether the API that a request with opts
// reaches has message threads.
func (r *BetaMessageService) supportsMessageThreads(opts []option.RequestOption) bool {
	return requestconfig.SupportsMessageThreads(slices.Concat(r.Options, opts)...)
}

// Messages returns a copy of the history: every message that was sent, every
// reply, and the messages staged for the next turn.
func (t *BetaMessageThread) Messages() []BetaMessageParam {
	return cloneBetaMessageParams(t.Params.Messages)
}

// PreviousMessageID returns the ID of the reply that the next turn continues
// from, or "" when the next turn starts a new thread.
func (t *BetaMessageThread) PreviousMessageID() string {
	t.adoptHistory()
	if link, ok := t.link.(betaThreadContinue); ok {
		return link.previousMessageID
	}
	return ""
}

// LastMessage returns the most recent reply that the thread recorded, or nil when
// there is none.
func (t *BetaMessageThread) LastMessage() *BetaMessage {
	t.adoptHistory()
	for _, entry := range slices.Backward(t.stored) {
		if entry.reply != nil {
			return entry.reply
		}
	}
	return nil
}

// Append stages messages for the next turn without sending them.
func (t *BetaMessageThread) Append(messages ...BetaMessageParam) {
	t.Params.Messages = slices.Clip(t.Params.Messages)
	for _, message := range cloneBetaMessageParams(messages) {
		if !t.repeatsLastReply(message) {
			t.Params.Messages = append(t.Params.Messages, message)
		}
	}
}

// repeatsLastReply reports whether the history ends with a reply that message
// is a copy of.
func (t *BetaMessageThread) repeatsLastReply(message BetaMessageParam) bool {
	n := len(t.stored)
	if n == 0 || n != len(t.Params.Messages) {
		return false
	}
	last := t.stored[n-1]
	return last.reply != nil && betaThreadJSON(message) == betaThreadJSON(last.param)
}

// adoptHistory makes the next turn start a new thread when Params.Messages no
// longer begins with the stored history.
func (t *BetaMessageThread) adoptHistory() {
	kept := len(t.stored)
	if kept > len(t.Params.Messages) || !slices.EqualFunc(t.stored, t.Params.Messages[:kept], func(entry betaThreadEntry, message BetaMessageParam) bool {
		// Only a message that was built again needs its JSON compared.
		return reflect.DeepEqual(entry.param, message) || betaThreadJSON(entry.param) == betaThreadJSON(message)
	}) {
		t.stored = nil
		t.startNewThread()
	}
}

// startNewThread makes the next turn send the whole history under a new thread.
func (t *BetaMessageThread) startNewThread() {
	if _, unthreaded := t.link.(betaThreadNone); !unthreaded {
		t.link = betaThreadCreate{}
	}
}

// Fork returns a new, independent thread with a copy of this thread's history,
// request fields and options. It sends no request.
func (t *BetaMessageThread) Fork() *BetaMessageThread {
	fork := &BetaMessageThread{
		Params:  t.Params,
		Options: slices.Clone(t.Options),
		service: t.service,
		stored:  slices.Clone(t.stored),
		link:    t.link,
	}
	fork.Params.Messages = t.Messages()
	fork.startNewThread()
	return fork
}

// Send sends the staged messages and params.Messages as one turn and returns the
// reply. params and opts apply to this turn only. It returns an error while
// another turn is in progress.
func (t *BetaMessageThread) Send(ctx context.Context, params BetaMessageThreadSendParams, opts ...option.RequestOption) (*BetaMessage, error) {
	if err := t.startTurn(params, opts); err != nil {
		return nil, err
	}
	defer t.endTurn()
	for {
		request, err := t.planRequest(params, opts)
		if err != nil {
			return nil, err
		}
		message, err := t.service.New(ctx, request.params, request.options...)
		// A reply whose structured output does not parse has still arrived.
		if message != nil {
			t.recordReply(request, message)
			return message, err
		}
		if !t.recoversFrom(err) {
			return nil, err
		}
	}
}

// SendStreaming is like [BetaMessageThread.Send] and returns the stream of the
// reply. The turn is in progress until the stream ends or is closed, so always
// call Close on the stream.
func (t *BetaMessageThread) SendStreaming(ctx context.Context, params BetaMessageThreadSendParams, opts ...option.RequestOption) *BetaMessageThreadStream {
	if err := t.startTurn(params, opts); err != nil {
		return &BetaMessageThreadStream{stream: ssestream.NewStream[BetaRawMessageStreamEventUnion](nil, err)}
	}
	for {
		request, err := t.planRequest(params, opts)
		if err != nil {
			t.endTurn()
			return &BetaMessageThreadStream{stream: ssestream.NewStream[BetaRawMessageStreamEventUnion](nil, err)}
		}
		stream := t.service.NewStreaming(ctx, request.params, request.options...)
		if stream.Err() == nil {
			return &BetaMessageThreadStream{stream: stream, turn: &betaThreadStreamTurn{thread: t, request: request}}
		}
		if !t.recoversFrom(stream.Err()) {
			t.endTurn()
			return &BetaMessageThreadStream{stream: stream}
		}
	}
}

// startTurn stages the messages of a turn and marks the turn as in progress.
func (t *BetaMessageThread) startTurn(params BetaMessageThreadSendParams, opts []option.RequestOption) error {
	if !t.service.supportsMessageThreads(slices.Concat(t.Options, opts)) {
		return errors.New("anthropic: message threads are not available on this client; use Beta.Messages.New with the whole history instead")
	}
	if t.inTurn {
		return errors.New("anthropic: another turn is in progress on this message thread; wait for it to end, or close its stream")
	}
	t.adoptHistory()
	t.givenStart = len(t.Params.Messages)
	t.Append(params.Messages...)
	t.givenCount = len(t.Params.Messages) - t.givenStart
	// With nothing new to send, the turn whose stream ended early is sent again,
	// so its partial reply goes.
	if n := len(t.stored); len(t.Params.Messages) == n && n > 0 && t.stored[n-1].partial {
		t.stored = t.stored[:n-1]
		t.Params.Messages = t.Params.Messages[:n-1]
	}
	t.inTurn = true
	return nil
}

// endTurn marks the turn in progress as ended. A turn that was not recorded
// takes back the messages it was given.
func (t *BetaMessageThread) endTurn() {
	t.inTurn = false
	start, end := t.givenStart, t.givenStart+t.givenCount
	t.givenCount = 0
	if start < end && end <= len(t.Params.Messages) {
		t.Params.Messages = slices.Concat(t.Params.Messages[:start], t.Params.Messages[end:])
	}
}

// planRequest returns the next request of a turn, or an error when the turn
// cannot be sent. It changes nothing on the thread.
func (t *BetaMessageThread) planRequest(params BetaMessageThreadSendParams, opts []option.RequestOption) (betaThreadRequest, error) {
	request := betaThreadRequest{
		params:  overlayBetaThreadParams(t.Params.toBetaMessageNewParams(), params.toBetaMessageNewParams()),
		options: slices.Concat(t.Options, opts, []option.RequestOption{stainlessheader.With(stainlessheader.MessageThread)}),
	}
	request.createOnly = betaThreadCreateOnlyFields{
		system: betaThreadJSON(request.params.System),
		tools:  betaThreadJSON(request.params.Tools),
	}

	messages := t.Params.Messages
	switch link := t.link.(type) {
	case betaThreadCreate:
		request.params.Thread.OfCreate = &BetaThreadCreateParams{}
	case betaThreadContinue:
		// The API takes the create-only fields from the thread, and would not see
		// that one of them changed.
		if err := request.createOnly.errChangedSince(link.createdWith); err != nil {
			return betaThreadRequest{}, err
		}
		request.params.Thread = BetaMessageThreadParamOfContinue(link.previousMessageID)
		request.params.System = nil
		request.params.Tools = nil
		messages = messages[len(t.stored):]
	}
	// Never nil: a turn with no new message sends an empty list. A copy: what is
	// stored has to stay as it was sent.
	request.params.Messages = cloneBetaMessageParams(messages)
	// Each turn may have set a cache breakpoint, and a request takes a limited number.
	if request.params.Thread.OfContinue == nil {
		answered := false
		for i, entry := range slices.Backward(t.stored) {
			if answered = answered || entry.reply != nil; answered {
				request.params.Messages[i] = withoutBetaCacheControl(messages[i])
			}
		}
	}

	if _, unthreaded := t.link.(betaThreadNone); !unthreaded && !slices.Contains(request.params.Betas, AnthropicBetaMessageThreads2026_08_12) {
		request.params.Betas = append(slices.Clip(request.params.Betas), AnthropicBetaMessageThreads2026_08_12)
	}
	return request, nil
}

// errChangedSince returns the error of a turn whose create-only fields differ
// from created, or nil when they are the same.
func (f betaThreadCreateOnlyFields) errChangedSince(created betaThreadCreateOnlyFields) error {
	var changed []string
	if f.system != created.system {
		changed = append(changed, "System")
	}
	if f.tools != created.tools {
		changed = append(changed, "Tools")
	}
	if len(changed) == 0 {
		return nil
	}
	message := "anthropic: " + strings.Join(changed, " and ") + " cannot be changed on a message thread because the API keeps the System and Tools the thread was created with; to change them, start a new thread with Fork"
	if f.tools != created.tools {
		message += "; to add or remove tools on this thread, send tool_addition or tool_removal blocks in a system message"
	}
	return errors.New(message)
}

// overlayBetaThreadParams returns base with every non-zero field of override,
// except Messages. It walks the fields by reflection, so that a field added to
// the params is overlaid without a change here.
func overlayBetaThreadParams(base, override BetaMessageNewParams) BetaMessageNewParams {
	from := reflect.ValueOf(override)
	to := reflect.ValueOf(&base).Elem()
	for i := range from.NumField() {
		field := from.Type().Field(i)
		if !field.IsExported() || field.Name == "Messages" {
			continue
		}
		if value := from.Field(i); !value.IsZero() {
			to.Field(i).Set(value)
		}
	}
	return base
}

// recordReply makes the next turn continue from message, the reply to request. A
// refused turn is not recorded: the API may hold it, and a request that continues
// the thread cannot remove it, so the next turn starts a new thread.
func (t *BetaMessageThread) recordReply(request betaThreadRequest, message *BetaMessage) {
	if message.StopReason == BetaStopReasonRefusal {
		t.startNewThread()
		return
	}
	t.storeTurn(request, message, false)
	// What the thread was created with is the thread's from then on, also when it
	// was given to that turn only.
	create, created := t.link.(betaThreadCreate)
	if created {
		t.Params.System, t.Params.Tools = request.params.System, request.params.Tools
	}
	if _, unthreaded := t.link.(betaThreadNone); !unthreaded {
		t.link = betaThreadContinue{previousMessageID: message.ID, createdWith: request.createOnly, afterNotFound: create.afterNotFound}
		// A thread keeps one model, so the API does not store a turn that a fallback
		// model wrote.
		for _, iteration := range message.Usage.Iterations {
			if iteration.Type == "fallback_message" {
				t.startNewThread()
			}
		}
	}
	adoptBetaContainer(&t.Params.Container, message)
}

// recordPartialReply records a turn whose stream ended after message_start and
// before message_stop. The API holds either the whole reply or none of it, so the
// next turn starts a new thread.
func (t *BetaMessageThread) recordPartialReply(request betaThreadRequest, message *BetaMessage) {
	t.storeTurn(request, message, true)
	t.startNewThread()
}

// storeTurn adds what request sent to the stored history, followed by reply if
// there is one. In Params.Messages the reply goes after what was sent, ahead of
// what was added during the turn.
func (t *BetaMessageThread) storeTurn(request betaThreadRequest, reply *BetaMessage, partial bool) {
	t.givenCount = 0
	sent := request.params.Messages
	if request.params.Thread.OfContinue == nil {
		sent = sent[len(t.stored):]
	}
	for _, message := range sent {
		t.stored = append(t.stored, betaThreadEntry{param: message})
	}
	if reply != nil {
		t.Params.Messages = slices.Insert(slices.Clip(t.Params.Messages), min(len(t.stored), len(t.Params.Messages)), reply.ToParam())
		t.stored = append(t.stored, betaThreadEntry{param: reply.ToParam(), reply: reply, partial: partial})
	}
}

// adoptBetaContainer makes later requests run in the container that message
// names, unless the caller named one: container-bound server tools reject a
// request that leaves out the container the turn before it ran in.
func adoptBetaContainer(container *BetaMessageNewParamsContainerUnion, message *BetaMessage) {
	id := message.Container.ID
	if id == "" {
		return
	}
	switch {
	case container.OfContainers != nil:
		if !container.OfContainers.ID.Valid() {
			pinned := *container.OfContainers
			pinned.ID = String(id)
			container.OfContainers = &pinned
		}
	case !container.OfString.Valid():
		container.OfString = String(id)
	}
}

// recoversFrom reports whether the turn should be sent again because err is a
// rejection of the thread field that the thread has now adjusted for. The link
// only moves from continue to create to none: at most three requests a turn.
func (t *BetaMessageThread) recoversFrom(err error) bool {
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		return false
	}
	continueUnthreaded := t.Params.UnsupportedRequestBehavior == BetaMessageThreadUnsupportedRequestBehaviorContinueUnthreaded
	switch code := gjson.Get(apiErr.RawJSON(), "error.details.error_code").String(); code {
	case "thread_not_found", "thread_already_continued", "thread_continue_in_flight":
		if link, ok := t.link.(betaThreadContinue); ok {
			notFound := code == "thread_not_found"
			// The API does not hold the thread that it has just created, so it may be
			// unable to store this conversation.
			if notFound && link.afterNotFound && continueUnthreaded {
				t.link = betaThreadNone{}
			} else {
				t.link = betaThreadCreate{afterNotFound: notFound}
			}
			return true
		}
	case "thread_unsupported_request":
		if _, unthreaded := t.link.(betaThreadNone); !unthreaded && continueUnthreaded {
			t.link = betaThreadNone{}
			return true
		}
	}
	return false
}

// betaThreadJSON returns the request JSON of value, to tell whether two values
// are sent alike. The keys of every object are sorted, because extra fields are
// marshaled in no fixed order.
func betaThreadJSON(value any) string {
	data, _ := json.Marshal(value)
	var decoded any
	if json.Unmarshal(data, &decoded) == nil {
		data, _ = json.Marshal(decoded)
	}
	return string(data)
}

// withoutBetaCacheControl returns a copy of message that is sent without the
// cache_control of its content blocks.
func withoutBetaCacheControl(message BetaMessageParam) BetaMessageParam {
	data, err := json.Marshal(message)
	if err != nil {
		return message
	}
	stripped := withoutBlockCacheControl(gjson.ParseBytes(data))
	if stripped == string(data) {
		return message
	}
	return param.Override[BetaMessageParam](json.RawMessage(stripped))
}

// withoutBlockCacheControl returns the JSON of block without its cache_control and
// that of the blocks and tool definitions under it. It enters only the keys that
// hold those: a cache_control key in a free-form field is not a breakpoint.
func withoutBlockCacheControl(block gjson.Result) string {
	if !block.IsObject() {
		return block.Raw
	}
	var members []string
	block.ForEach(func(key, value gjson.Result) bool {
		raw := value.Raw
		switch key.Str {
		case "cache_control":
			return true
		case "content", "source", "tool", "definition", "tool_references":
			if value.IsArray() {
				var blocks []string
				for _, nested := range value.Array() {
					blocks = append(blocks, withoutBlockCacheControl(nested))
				}
				raw = "[" + strings.Join(blocks, ",") + "]"
			} else {
				raw = withoutBlockCacheControl(value)
			}
		}
		members = append(members, key.Raw+":"+raw)
		return true
	})
	return "{" + strings.Join(members, ",") + "}"
}

// BetaMessageThreadStream is the stream of one turn of a [BetaMessageThread].
// It is read like an [ssestream.Stream]:
//
//	stream := thread.SendStreaming(ctx, params)
//	defer stream.Close()
//	for stream.Next() {
//		event := stream.Current()
//	}
//	if err := stream.Err(); err != nil {
//		return err
//	}
type BetaMessageThreadStream struct {
	// stream holds the error of a turn that could not start.
	stream *ssestream.Stream[BetaRawMessageStreamEventUnion]
	// turn is nil once the thread has recorded how the turn ended.
	turn *betaThreadStreamTurn
}

type betaThreadStreamTurn struct {
	thread  *BetaMessageThread
	request betaThreadRequest
	// reply is the reply so far. It has an ID once message_start has arrived.
	reply BetaMessage
	// finishedBlocks counts the content blocks that received content_block_stop.
	finishedBlocks int
}

// Next reads the next event and reports whether there is one. It returns false
// at the end of the stream and when the stream fails; Err tells the two apart.
func (s *BetaMessageThreadStream) Next() bool {
	if !s.stream.Next() {
		s.endTurn()
		return false
	}
	turn := s.turn
	if turn == nil {
		return true
	}
	event := s.stream.Current()
	if err := turn.reply.Accumulate(event); err != nil {
		s.Close()
		s.stream = ssestream.NewStream[BetaRawMessageStreamEventUnion](nil, err)
		return false
	}
	switch event.Type {
	case "content_block_stop":
		turn.finishedBlocks++
	case "message_stop":
		s.turn = nil
		turn.thread.recordReply(turn.request, &turn.reply)
		turn.thread.endTurn()
	}
	return true
}

// Current returns the event that the last call to Next read.
func (s *BetaMessageThreadStream) Current() BetaRawMessageStreamEventUnion {
	return s.stream.Current()
}

// Err returns the error that ended the stream, or the error that kept the turn
// from starting.
func (s *BetaMessageThreadStream) Err() error {
	return s.stream.Err()
}

// Close ends the stream and the turn.
func (s *BetaMessageThreadStream) Close() error {
	s.endTurn()
	return s.stream.Close()
}

// endTurn records a turn whose stream ended before message_stop.
func (s *BetaMessageThreadStream) endTurn() {
	turn := s.turn
	if turn == nil {
		return
	}
	s.turn = nil
	switch {
	case turn.reply.StopReason == BetaStopReasonRefusal:
		turn.thread.recordReply(turn.request, &turn.reply)
	case turn.reply.ID != "":
		turn.thread.recordPartialReply(turn.request, partialBetaThreadReply(&turn.reply, turn.finishedBlocks))
	}
	turn.thread.endTurn()
}

// partialBetaThreadReply returns the part of reply that can be sent back, or nil:
// its finished blocks and the text so far of an unfinished text block. Any other
// unfinished block has an incomplete input or signature.
func partialBetaThreadReply(reply *BetaMessage, finishedBlocks int) *BetaMessage {
	finishedBlocks = min(finishedBlocks, len(reply.Content))
	content := slices.Clone(reply.Content[:finishedBlocks])
	if finishedBlocks < len(reply.Content) {
		if unfinished := reply.Content[finishedBlocks]; unfinished.Type == "text" && unfinished.Text != "" {
			refreshBetaContentBlockRaw(&unfinished)
			content = append(content, unfinished)
		}
	}
	if len(content) == 0 {
		return nil
	}
	raw, err := sjson.SetRaw(reply.RawJSON(), "content", rawArray(content))
	if err != nil {
		return nil
	}
	partial := &BetaMessage{}
	if err := partial.UnmarshalJSON([]byte(raw)); err != nil {
		return nil
	}
	return partial
}
