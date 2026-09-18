package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"os"

	"github.com/anthropics/anthropic-sdk-go/internal/stainlessheader"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"golang.org/x/sync/errgroup"
)

// BetaTool represents a tool that can be executed by the BetaToolRunner.
type BetaTool interface {
	// Name returns the tool's name
	Name() string
	// Description returns the tool's description
	Description() string
	// InputSchema returns the JSON schema for the tool's input
	InputSchema() BetaToolInputSchemaParam
	// Execute runs the tool with raw JSON input and returns one or more result blocks.
	Execute(ctx context.Context, input json.RawMessage) ([]BetaToolResultBlockParamContentUnion, error)
}

// BetaToolRunnerParams contains parameters for creating a BetaToolRunner or BetaToolRunnerStreaming.
type BetaToolRunnerParams struct {
	BetaMessageNewParams
	// MaxIterations limits the number of API calls. When set to 0 (the default),
	// there is no limit and the runner continues until the model stops using tools.
	// Compaction requests sent for CompactBeforeNextTurn are not counted.
	MaxIterations int
}

// betaToolRunnerBase holds state and logic shared by BetaToolRunner and BetaToolRunnerStreaming.
type betaToolRunnerBase struct {
	messageService *BetaMessageService
	// Params contains the configuration for the tool runner.
	// This field is exported so users can modify parameters directly.
	Params         BetaToolRunnerParams
	toolMap        map[string]BetaTool
	iterationCount int
	lastMessage    *BetaMessage
	completed      bool
	opts           []option.RequestOption
	err            error

	pendingCompaction *BetaCompactionConfigUnionParam
	// compacting is set from sending a compaction request until the next request
	// is built, or until that request fails. CompactBeforeNextTurn does nothing
	// meanwhile: there is nothing new to summarize.
	compacting bool
}

func newBetaToolRunnerBase(messageService *BetaMessageService, tools []BetaTool, params BetaToolRunnerParams, opts []option.RequestOption) betaToolRunnerBase {
	toolMap := make(map[string]BetaTool)
	apiTools := make([]BetaToolUnionParam, len(tools))

	for i, tool := range tools {
		toolMap[tool.Name()] = tool
		apiTools[i] = BetaToolUnionParam{
			OfTool: &BetaToolParam{
				Name:        tool.Name(),
				Description: String(tool.Description()),
				InputSchema: tool.InputSchema(),
			},
		}
	}

	// Add tools to the API params
	params.BetaMessageNewParams.Tools = apiTools
	params.Messages = append([]BetaMessageParam{}, params.Messages...)

	opts = append([]option.RequestOption{stainlessheader.With(stainlessheader.BetaToolRunner)}, opts...)

	return betaToolRunnerBase{
		messageService: messageService,
		Params:         params,
		toolMap:        toolMap,
		opts:           opts,
	}
}

// LastMessage returns the most recent assistant message, or nil if no messages have been received yet.
func (b *betaToolRunnerBase) LastMessage() *BetaMessage {
	return b.lastMessage
}

// AppendMessages adds messages to the conversation history.
// This is a convenience method equivalent to:
//
//	runner.Params.Messages = append(runner.Params.Messages, messages...)
func (b *betaToolRunnerBase) AppendMessages(messages ...BetaMessageParam) {
	b.Params.Messages = append(b.Params.Messages, messages...)
}

// Messages returns a copy of the current conversation history.
// The returned slice can be safely modified without affecting the runner's state.
func (b *betaToolRunnerBase) Messages() []BetaMessageParam {
	result := make([]BetaMessageParam, len(b.Params.Messages))
	copy(result, b.Params.Messages)
	return result
}

// CompactBeforeNextTurn compacts the conversation before the model's next turn.
// Once the current turn has finished, including any tool calls, the runner
// requests a summary, replaces the message history with the compaction response
// and returns that response like any other message. Calling this again before
// the compaction runs replaces the earlier call. compaction is the same config
// as [BetaMessageNewParams.Compaction]; its zero value means {"type": "summarize"}.
func (b *betaToolRunnerBase) CompactBeforeNextTurn(compaction BetaCompactionConfigUnionParam) {
	if b.compacting {
		return
	}
	if param.IsOmitted(compaction) {
		compaction.OfSummarize = &BetaSummarizeCompactionParam{}
	}
	b.pendingCompaction = &compaction
}

// IterationCount returns the number of API calls made so far, not counting
// compaction requests sent for CompactBeforeNextTurn.
// This is incremented each time a turn makes an API call.
func (b *betaToolRunnerBase) IterationCount() int {
	return b.iterationCount
}

// IsCompleted returns true if the conversation has finished, either because
// the model stopped using tools or the maximum iteration limit was reached.
func (b *betaToolRunnerBase) IsCompleted() bool {
	return b.completed
}

// Err returns the last error that occurred during iteration, if any.
// This is useful when using All() or AllStreaming() to check for errors
// after the iteration completes.
func (b *betaToolRunnerBase) Err() error {
	return b.err
}

// adoptContainer carries the container the last turn ran in onto the next
// request: container-bound server tools reject a follow-up that omits it, so
// its id is forwarded unless the caller pinned one themselves.
func (b *betaToolRunnerBase) adoptContainer(message *BetaMessage) {
	id := message.Container.ID
	if id == "" {
		return
	}
	container := &b.Params.Container
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

// toolRunnerStep is what the runner does with a turn, decided by its stop reason.
type toolRunnerStep int

const (
	// stepRunTools answers the turn's client tool calls and continues, or stops
	// when there are none.
	stepRunTools toolRunnerStep = iota
	// stepResume sends the unfinished turn back unchanged, without executing
	// any tool calls it carries, so the server continues it.
	stepResume
	// stepStop ends the conversation with the turn as the final message,
	// without executing its tool calls.
	stepStop
)

// determineNextStepFromStopReason classifies every generated stop reason
// explicitly so a new one has to be placed in a bucket here; values unknown at
// runtime stop like end_turn rather than erroring.
func determineNextStepFromStopReason(reason BetaStopReason) toolRunnerStep {
	switch reason {
	case BetaStopReasonToolUse:
		return stepRunTools
	case BetaStopReasonPauseTurn:
		// A long-running server tool paused the turn; sending it back unchanged
		// resumes it.
		return stepResume
	case BetaStopReasonCompaction:
		// pause_after_compaction hands the turn back before the model answers;
		// sending it back unchanged continues it.
		return stepResume
	case BetaStopReasonRefusal:
		// A refusal-terminated turn is terminal: its tool calls belong to a dead
		// conversation — executing them fires side effects the caller never
		// confirmed and produces tool_results that cannot be coherently replayed.
		return stepStop
	case BetaStopReasonMaxTokens, BetaStopReasonModelContextWindowExceeded:
		// A cut-off turn left its last call's arguments incomplete.
		return stepStop
	case BetaStopReasonEndTurn, BetaStopReasonStopSequence:
		return stepStop
	default:
		// An unrecognized stop reason stops like end_turn.
		return stepStop
	}
}

// executeTools processes any tool use blocks in the given message and returns a tool result message.
// Returns:
//   - (result, nil) if tools executed successfully
//   - (nil, nil) if the turn did not stop for tool use or has no client tool calls
//   - (nil, ctx.Err()) if context was cancelled
func (b *betaToolRunnerBase) executeTools(ctx context.Context, message *BetaMessage) (*BetaMessageParam, error) {
	toolUseBlocks := clientToolCalls(message)
	if len(toolUseBlocks) == 0 {
		return nil, nil
	}

	// Execute all tools in parallel using errgroup for proper cancellation handling
	results := make([]BetaContentBlockParamUnion, len(toolUseBlocks))
	available := b.availableToolNames()

	g, gctx := errgroup.WithContext(ctx)
	for i, toolUse := range toolUseBlocks {
		g.Go(func() error {
			// Check for cancellation before executing tool
			select {
			case <-gctx.Done():
				return gctx.Err()
			default:
			}
			result := b.executeToolUse(gctx, toolUse, available)
			results[i] = BetaContentBlockParamUnion{OfToolResult: &result}
			return nil // tool errors become result content, not Go errors
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Create user message with tool results
	userMessage := NewBetaUserMessage(results...)
	return &userMessage, nil
}

// clientToolCalls returns the tool calls the runner answers for the given
// message: none unless the turn stopped for tool use.
func clientToolCalls(message *BetaMessage) []BetaToolUseBlock {
	if determineNextStepFromStopReason(message.StopReason) != stepRunTools {
		return nil
	}
	return unansweredToolCalls(message)
}

// unansweredToolCalls returns the message's client tool calls that a
// tool_result can still answer.
func unansweredToolCalls(message *BetaMessage) []BetaToolUseBlock {
	// Tool calls before the last fallback block belong to the attempt that
	// refused; the fallback middleware strips them from replayed history, so
	// answering them would orphan the tool_result.
	seam := -1
	for i, block := range message.Content {
		if block.Type == "fallback" {
			seam = i
		}
	}

	var toolUseBlocks []BetaToolUseBlock

	// Find all tool use blocks in the message
	for i, block := range message.Content {
		if i > seam && block.Type == "tool_use" {
			toolUseBlocks = append(toolUseBlocks, block.AsToolUse())
		}
	}
	return toolUseBlocks
}

func newBetaToolResultErrorBlockParam(toolUseID string, errorText string) BetaToolResultBlockParam {
	return NewBetaToolResultTextBlockParam(toolUseID, errorText, true)
}

// availableToolNames returns the tool names currently offered to the model:
// every registered tool, minus names dropped by tool_removal blocks in
// role "system" messages, plus names re-enabled by later tool_addition
// blocks. Removal is only a hint — the model can still call a removed tool —
// so a removed tool must resolve to the same not-found result as one that
// was never registered.
func (b *betaToolRunnerBase) availableToolNames() map[string]struct{} {
	available := make(map[string]struct{}, len(b.toolMap))
	for name := range b.toolMap {
		available[name] = struct{}{}
	}
	for _, message := range b.Params.Messages {
		if message.Role != BetaMessageParamRoleSystem {
			continue
		}
		for _, block := range message.Content {
			applyToolChange(block, available)
		}
	}
	return available
}

// applyToolChange folds one system-message content block into the available
// tool-name set. The populated Of* variant is the discriminator; every other
// block type, including ones added after this SDK was generated, leaves the
// set untouched.
func applyToolChange(block BetaContentBlockParamUnion, available map[string]struct{}) {
	switch {
	case block.OfToolRemoval != nil:
		if name, ok := removalRefName(block.OfToolRemoval.Tool); ok {
			delete(available, name)
		}
	case block.OfToolAddition != nil:
		if name, ok := additionRefName(block.OfToolAddition.Tool); ok {
			available[name] = struct{}{}
		}
	}
}

// additionRefName and removalRefName resolve the referenced tool's name across
// the generated addition/removal tool unions. Only a tool_reference names a
// locally runnable tool; MCP references run server-side and unknown variants
// are ignored.
func additionRefName(u BetaRequestToolAdditionBlockToolUnionParam) (string, bool) {
	switch {
	case u.OfToolReference != nil:
		return u.OfToolReference.Name, u.OfToolReference.Name != ""
	default:
		return "", false
	}
}

func removalRefName(u BetaRequestToolRemovalBlockToolUnionParam) (string, bool) {
	switch {
	case u.OfToolReference != nil:
		return u.OfToolReference.Name, u.OfToolReference.Name != ""
	default:
		return "", false
	}
}

// executeToolUse executes a single tool use block and returns the result.
// available is the current set of offered tool names (see availableToolNames).
func (b *betaToolRunnerBase) executeToolUse(ctx context.Context, toolUse BetaToolUseBlock, available map[string]struct{}) BetaToolResultBlockParam {
	_, allowed := available[toolUse.Name]
	tool, exists := b.toolMap[toolUse.Name]
	if !exists || !allowed {
		return newBetaToolResultErrorBlockParam(
			toolUse.ID,
			fmt.Sprintf("Error: Tool '%s' not found", toolUse.Name),
		)
	}

	// Parse and execute the tool
	inputBytes, err := json.Marshal(toolUse.Input)
	if err != nil {
		return newBetaToolResultErrorBlockParam(
			toolUse.ID,
			fmt.Sprintf("Error: Failed to marshal tool input: %v", err),
		)
	}

	content, err := tool.Execute(ctx, inputBytes)
	if err != nil {
		return newBetaToolResultErrorBlockParam(
			toolUse.ID,
			fmt.Sprintf("Error: %v", err),
		)
	}

	return BetaToolResultBlockParam{
		ToolUseID: toolUse.ID,
		Content:   content,
	}
}

// toolRunnerRequest is the next request the runner sends.
type toolRunnerRequest struct {
	params BetaMessageNewParams
}

// nextRequest answers the last turn's tool calls and works out what to send
// next. It returns nil once the run is over.
func (b *betaToolRunnerBase) nextRequest(ctx context.Context) (*toolRunnerRequest, error) {
	if b.completed {
		return nil, nil
	}
	if !param.IsOmitted(b.Params.Compaction) {
		// Unset it so the refusal is reported once; left set, it would block
		// every later turn.
		b.Params.Compaction = BetaCompactionConfigUnionParam{}
		return nil, errors.New("anthropic: Params.Compaction cannot be set on a tool runner because every request in the loop would compact again, so it was cleared; call CompactBeforeNextTurn when the conversation should be compacted instead")
	}
	if b.pendingCompaction != nil {
		if err := b.checkCanCompact(); err != nil {
			b.pendingCompaction = nil
			return nil, err
		}
	}

	// The response to a compaction request has nothing to answer or resume,
	// whatever it stopped for.
	turn := b.lastMessage
	if b.compacting {
		turn = nil
	}
	b.compacting = false

	paused := false
	if turn != nil {
		paused = determineNextStepFromStopReason(turn.StopReason) == stepResume
		if !paused && len(clientToolCalls(turn)) == 0 {
			return b.finalTurnRequest(turn), nil
		}
	}
	if b.Params.MaxIterations > 0 && b.iterationCount >= b.Params.MaxIterations {
		b.completed = true
		return nil, nil
	}
	if turn != nil && !paused {
		toolMessage, err := b.executeTools(ctx, turn)
		if err != nil {
			return nil, err
		}
		if toolMessage != nil {
			b.Params.Messages = append(b.Params.Messages, *toolMessage)
		}
	}

	// The API cannot compact a conversation that stops mid-turn, so a paused
	// turn is resumed first.
	if b.pendingCompaction != nil && !paused {
		return b.takeCompactionRequest(), nil
	}
	b.iterationCount++
	return &toolRunnerRequest{params: b.Params.BetaMessageNewParams}, nil
}

// finalTurnRequest ends the run. It returns one last request when a compaction
// is pending and the API can accept it after turn.
func (b *betaToolRunnerBase) finalTurnRequest(turn *BetaMessage) *toolRunnerRequest {
	b.completed = true
	if b.pendingCompaction == nil {
		return nil
	}
	// A turn that was cut short can end with tool calls that are never run, and
	// the API cannot compact a conversation whose last turn has an unanswered
	// tool call.
	if len(unansweredToolCalls(turn)) > 0 {
		fmt.Fprintf(os.Stderr, "Warning: The tool runner skipped the pending compaction because the last turn (stop_reason %q) ended with tool calls that were not run. Call CompactBeforeNextTurn again if you continue the conversation.\n", turn.StopReason)
		b.pendingCompaction = nil
		return nil
	}
	return b.takeCompactionRequest()
}

func (b *betaToolRunnerBase) checkCanCompact() error {
	// The compaction request goes out without context_management, so the API
	// cannot reject this pairing there: it would run and bill the compaction,
	// then reject the next request.
	for _, edit := range b.Params.ContextManagement.Edits {
		if edit.OfCompact20260112 != nil {
			return errors.New("anthropic: CompactBeforeNextTurn cannot be used while Params.ContextManagement has a compaction edit because the API does not accept a compaction block together with one; remove the edit and call it again")
		}
	}
	return nil
}

func (b *betaToolRunnerBase) takeCompactionRequest() *toolRunnerRequest {
	params := b.Params.BetaMessageNewParams
	params.Compaction = *b.pendingCompaction
	// The API refuses compaction alongside context_management; later requests
	// keep it.
	params.ContextManagement = BetaContextManagementConfigParam{}
	b.pendingCompaction = nil
	b.compacting = true
	return &toolRunnerRequest{params: params}
}

// handleResponse records the response to request and carries it into the
// conversation history.
func (b *betaToolRunnerBase) handleResponse(request *toolRunnerRequest, message *BetaMessage) error {
	if b.compacting {
		return b.finishCompaction(request, message)
	}
	b.lastMessage = message
	b.Params.Messages = append(b.Params.Messages, message.ToParam())
	b.adoptContainer(message)
	return nil
}

func (b *betaToolRunnerBase) finishCompaction(request *toolRunnerRequest, message *BetaMessage) error {
	// Params is an exported field, so a change made while the response was
	// streaming can only be noticed here: appending changes the length and
	// assigning a new slice changes the first element's address.
	sent, current := request.params.Messages, b.Params.Messages
	if len(current) != len(sent) || (len(sent) > 0 && &current[0] != &sent[0]) {
		return errors.New("anthropic: the tool runner's messages were changed while the conversation was being compacted, so the compaction response did not replace them; change Params.Messages once the compaction response has been read")
	}
	// A failed compaction comes back as a block with null content, or as no
	// block at all.
	for _, block := range message.Content {
		if block.Type == "compaction" && block.Content.OfString != "" {
			// The response has to be sent back as it came, first, replacing the
			// messages it summarizes.
			b.lastMessage = message
			b.Params.Messages = []BetaMessageParam{message.ToParam()}
			return nil
		}
	}
	fmt.Fprint(os.Stderr, "Warning: Compaction produced no summary, so the tool runner kept the conversation as it is.\n")
	return nil
}

// BetaToolRunner manages the automatic conversation loop between the assistant and tools
// using non-streaming API calls. It implements an iterator pattern for processing
// conversation turns.
//
// A BetaToolRunner is NOT safe for concurrent use. All methods must be called
// from a single goroutine. However, tool handlers ARE called concurrently
// when multiple tools are invoked in a single turn - ensure your handlers
// are thread-safe.
type BetaToolRunner struct {
	betaToolRunnerBase
}

// NewToolRunner creates a BetaToolRunner that automatically handles the loop between
// the model generating tool calls, executing those tool calls, and sending the
// results back to the model until a final answer is produced or the maximum
// number of iterations is reached.
func (r *BetaMessageService) NewToolRunner(tools []BetaTool, params BetaToolRunnerParams, opts ...option.RequestOption) *BetaToolRunner {
	return &BetaToolRunner{
		betaToolRunnerBase: newBetaToolRunnerBase(r, tools, params, opts),
	}
}

// NextMessage advances the conversation by one turn. It executes any pending tool calls
// from the previous message, then makes an API call to get the model's next response.
//
// Returns:
//   - (message, nil) on success with the assistant's response
//   - (nil, nil) when the conversation is complete (no more tool calls or max iterations reached)
//   - (nil, error) if an error occurred during tool execution or API call
func (r *BetaToolRunner) NextMessage(ctx context.Context) (*BetaMessage, error) {
	// Execute any pending tool calls from the last message
	request, err := r.nextRequest(ctx)
	if err != nil {
		r.err = err
		return nil, err
	}
	if request == nil {
		return nil, nil
	}

	// Make API call
	message, err := r.messageService.New(ctx, request.params, r.opts...)
	if err != nil {
		r.compacting = false
		r.err = err
		return nil, fmt.Errorf("failed to get next message: %w", err)
	}

	if err := r.handleResponse(request, message); err != nil {
		r.err = err
		return nil, err
	}

	return message, nil
}

// RunToCompletion repeatedly calls NextMessage until the conversation is complete,
// either because the model stopped using tools or the maximum iteration limit was reached.
//
// Returns the final assistant message and any error that occurred.
func (r *BetaToolRunner) RunToCompletion(ctx context.Context) (*BetaMessage, error) {
	for {
		message, err := r.NextMessage(ctx)
		if err != nil {
			return nil, err
		}
		if message == nil {
			return r.lastMessage, nil
		}
	}
}

// All returns an iterator that yields all messages until the conversation completes.
// This is a convenience method for iterating over the entire conversation.
//
// Example usage:
//
//	for message, err := range runner.All(ctx) {
//	    if err != nil {
//	        return err
//	    }
//	    // process message
//	}
func (r *BetaToolRunner) All(ctx context.Context) iter.Seq2[*BetaMessage, error] {
	return func(yield func(*BetaMessage, error) bool) {
		for {
			message, err := r.NextMessage(ctx)
			r.err = err
			if message == nil {
				if err != nil {
					yield(nil, err)
				}
				return
			}
			if !yield(message, err) {
				return
			}
		}
	}
}

// BetaToolRunnerStreaming manages the automatic conversation loop between the assistant
// and tools using streaming API calls. It implements an iterator pattern for processing
// streaming events across conversation turns.
//
// A BetaToolRunnerStreaming is NOT safe for concurrent use. All methods must be called
// from a single goroutine. However, tool handlers ARE called concurrently
// when multiple tools are invoked in a single turn - ensure your handlers
// are thread-safe.
type BetaToolRunnerStreaming struct {
	betaToolRunnerBase
}

// NewToolRunnerStreaming creates a BetaToolRunnerStreaming that automatically handles
// the loop between the model generating tool calls, executing those tool calls, and
// sending the results back to the model using streaming API calls until a final answer
// is produced or the maximum number of iterations is reached.
func (r *BetaMessageService) NewToolRunnerStreaming(tools []BetaTool, params BetaToolRunnerParams, opts ...option.RequestOption) *BetaToolRunnerStreaming {
	return &BetaToolRunnerStreaming{
		betaToolRunnerBase: newBetaToolRunnerBase(r, tools, params, opts),
	}
}

// NextStreaming advances the conversation by one turn with streaming. It executes any
// pending tool calls from the previous message, then makes a streaming API call.
//
// Returns an iterator that yields streaming events as they arrive. The iterator should
// be fully consumed to ensure the message is properly accumulated for subsequent turns.
//
// If an error occurs, it will be yielded as the second value in the iterator pair.
// Check IsCompleted() after consuming the iterator to determine if the conversation
// has finished.
func (r *BetaToolRunnerStreaming) NextStreaming(ctx context.Context) iter.Seq2[BetaRawMessageStreamEventUnion, error] {
	return func(yield func(BetaRawMessageStreamEventUnion, error) bool) {
		// Execute any pending tool calls from the last message
		request, err := r.nextRequest(ctx)
		if err != nil {
			r.err = err
			yield(BetaRawMessageStreamEventUnion{}, err)
			return
		}
		if request == nil {
			return
		}

		// Make streaming API call
		stream := r.messageService.NewStreaming(ctx, request.params, r.opts...)
		defer stream.Close()
		responded := false
		defer func() {
			if !responded {
				r.compacting = false
			}
		}()

		// We need to collect the final message from the stream for the next iteration
		finalMessage := &BetaMessage{}
		for stream.Next() {
			event := stream.Current()
			err := finalMessage.Accumulate(event)
			if err != nil {
				r.err = fmt.Errorf("failed to accumulate streaming event: %w", err)
				yield(BetaRawMessageStreamEventUnion{}, r.err)
				return
			}

			if !yield(event, nil) {
				return
			}
		}

		// Check for stream errors after the loop exits
		if stream.Err() != nil {
			r.err = stream.Err()
			yield(BetaRawMessageStreamEventUnion{}, r.err)
			return
		}

		responded = true
		if err := r.handleResponse(request, finalMessage); err != nil {
			r.err = err
			yield(BetaRawMessageStreamEventUnion{}, err)
		}
	}
}

// AllStreaming returns an iterator of iterators, where each inner iterator yields
// streaming events for a single turn of the conversation. The outer iterator continues
// until the conversation completes.
//
// Example usage:
//
//	for events, err := range runner.AllStreaming(ctx) {
//	    if err != nil {
//	        return err
//	    }
//	    for event, err := range events {
//	        if err != nil {
//	            return err
//	        }
//	        // process streaming event
//	    }
//	}
func (r *BetaToolRunnerStreaming) AllStreaming(ctx context.Context) iter.Seq2[iter.Seq2[BetaRawMessageStreamEventUnion, error], error] {
	return func(yield func(iter.Seq2[BetaRawMessageStreamEventUnion, error], error) bool) {
		for !r.completed {
			eventSeq := r.NextStreaming(ctx)
			if !yield(eventSeq, nil) {
				return
			}
		}
	}
}
