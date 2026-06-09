package betafallback

// Streaming half of [BetaRefusalFallbackMiddleware]: a refused streaming
// attempt is retried down the chain and the retry's events are spliced onto
// the open stream in server-parity framing — one message_start, a fallback
// boundary block at each model seam, monotonic block indices, and a terminal
// usage.iterations ledger covering every hop. The splicer is pull-based: hop
// requests are issued from Read on the consumer's goroutine.
//
// The streaming contract differs from the non-streaming loop where the wire
// does:
//   - A mid-stream refusal retries only when it minted a credit token; with
//     content already streamed and no token, the refusal surfaces untouched.
//     A pre-stream refusal retries either way — nothing streamed, so the
//     retry is invisible.
//   - When the refusal advertises stop_details.fallback_has_prefill_claim,
//     the retry appends one trailing assistant turn echoing the refused
//     partial verbatim; the continuation grows across hops.
//   - A hop whose request fails (HTTP error, transport error, or a non-SSE
//     200) was never reached: the token and continuation carry to the next
//     entry and no boundary is emitted for it.
//   - The non-streaming loop does not synthesize usage.iterations; only the
//     fallback boundary blocks are prepended.
//   - The non-streaming loop drops a refused hop's partial output; only the
//     serving hop's content reaches the caller.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"sort"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
)

// handleStreaming runs the streaming fallback path for an eligible request.
// body is the parsed original request body (stream: true).
func handleStreaming(
	req *http.Request,
	next func(*http.Request) (*http.Response, error),
	orig []byte,
	body map[string]json.RawMessage,
	fallbacks []anthropic.BetaFallbackParam,
	state *BetaFallbackState,
) (*http.Response, error) {
	send := func(payload []byte) (*http.Response, error) {
		r := requestWithBody(req, payload)
		ensureFallbackCreditBeta(r.Header)
		return next(r)
	}

	orig = trimHistory(body, orig)

	chainStart := 0
	attempt := orig
	if state != nil && state.Index() >= 0 {
		var err error
		if attempt, _, err = mergeFallback(body, fallbacks[state.Index()], ""); err != nil {
			return nil, err
		}
		chainStart = state.Index() + 1
	}

	res, err := send(attempt)
	if err != nil || !isEventStream(res) {
		// Errors and non-SSE responses surface verbatim; there is no refusal
		// to act on.
		return res, err
	}

	s := &streamSplicer{
		res:       res,
		dec:       ssestream.NewDecoder(res),
		body:      body,
		fallbacks: fallbacks,
		next:      chainStart,
		state:     state,
		send:      send,
		hopModel:  gjson.GetBytes(attempt, "model").String(),
		open:      map[int]bool{},
		blocks:    map[int]*rawBlock{},
	}
	s.effModel = s.hopModel

	out := new(http.Response)
	*out = *res
	out.Body = s
	out.ContentLength = -1
	out.Header = res.Header.Clone()
	out.Header.Del("Content-Length")
	return out, nil
}

// isEventStream reports whether res is a readable SSE stream. An encoded
// body (caller-set Accept-Encoding) can't be inspected.
func isEventStream(res *http.Response) bool {
	if res == nil || res.StatusCode != http.StatusOK || res.Header.Get("Content-Encoding") != "" {
		return false
	}
	mediaType, _, _ := mime.ParseMediaType(res.Header.Get("Content-Type"))
	return mediaType == "text/event-stream"
}

// rawBlock accumulates one content block of the current hop from its start
// event and deltas, so a prefill claim can echo it verbatim.
type rawBlock struct {
	start     []byte // raw content_block JSON from content_block_start
	blockType string
	text      strings.Builder
	thinking  strings.Builder
	signature strings.Builder
	partial   strings.Builder // input_json_delta fragments
	citations []json.RawMessage
}

// finalize folds the accumulated deltas back into the block's JSON.
func (b *rawBlock) finalize() json.RawMessage {
	out := b.start
	set := func(path string, v any) {
		if patched, err := sjson.SetBytes(out, path, v); err == nil {
			out = patched
		}
	}
	if b.text.Len() > 0 {
		set("text", gjson.GetBytes(b.start, "text").String()+b.text.String())
	}
	if b.thinking.Len() > 0 {
		set("thinking", gjson.GetBytes(b.start, "thinking").String()+b.thinking.String())
	}
	if b.signature.Len() > 0 {
		set("signature", gjson.GetBytes(b.start, "signature").String()+b.signature.String())
	}
	if raw := b.partial.String(); raw != "" && gjson.Valid(raw) {
		if patched, err := sjson.SetRawBytes(out, "input", []byte(raw)); err == nil {
			out = patched
		}
	}
	for _, c := range b.citations {
		if patched, err := sjson.SetRawBytes(out, "citations.-1", c); err == nil {
			out = patched
		}
	}
	return out
}

type seam struct{ from, to string }

// block returns the fallback boundary content block as raw JSON — the shape
// both the streaming splice and the non-streaming prepend emit, and that
// trimFallbackTurns keys off on history replay.
func (s seam) block() json.RawMessage {
	block, _ := json.Marshal(map[string]any{
		"type": "fallback",
		"from": map[string]any{"model": s.from},
		"to":   map[string]any{"model": s.to},
	})
	return block
}

type heldRefusal struct {
	startRaw []byte // raw message_start data; nil if none arrived
	deltaRaw []byte // raw terminal message_delta data
	stopRaw  []byte // raw message_stop data
}

// streamSplicer is an io.ReadCloser forwarding the current hop's SSE stream,
// retrying refusals down the chain, and splicing retries' events into one
// app-visible stream.
type streamSplicer struct {
	out  bytes.Buffer
	done bool
	err  error

	res *http.Response
	dec ssestream.Decoder

	body      map[string]json.RawMessage
	fallbacks []anthropic.BetaFallbackParam
	next      int // next chain index to try
	state     *BetaFallbackState
	send      func(body []byte) (*http.Response, error)

	// Output framing.
	wireOpen     bool
	nextIndex    int // next free output block index
	hopOffset    int // output index of the current hop's wire index 0
	pendingSeams []seam
	primaryID    string // first hop's message id; spliced starts adopt it

	// Current hop.
	hopModel    string // model string the hop was requested as
	effModel    string // last seam's to.model within the hop, else hopModel
	isRetryHop  bool
	startRaw    []byte // held message_start data (emitted when the wire opens)
	contentSeen bool
	open        map[int]bool
	blocks      map[int]*rawBlock
	blockOrder  []int
	sawSeam     bool

	// Chain state.
	token          string
	continuation   []json.RawMessage
	lastClaimCount int
	ledger         []json.RawMessage
	held           *heldRefusal
	terminalSent   bool // a terminal message_delta reached the consumer
	stopSent       bool // a message_stop reached the consumer

	// lastHopStatus is the HTTP status of the most recent failed hop
	// attempt, 0 for transport errors.
	lastHopStatus int
}

func (s *streamSplicer) Read(p []byte) (int, error) {
	for s.out.Len() == 0 && !s.done {
		s.advance()
	}
	if s.out.Len() > 0 {
		return s.out.Read(p)
	}
	if s.err != nil {
		return 0, s.err
	}
	return 0, io.EOF
}

func (s *streamSplicer) Close() error {
	s.done = true
	if s.res != nil {
		return s.res.Body.Close()
	}
	return nil
}

func (s *streamSplicer) finish(err error) {
	s.done = true
	s.err = err
	if s.res != nil {
		s.res.Body.Close()
	}
}

func (s *streamSplicer) emit(eventType string, data []byte) {
	fmt.Fprintf(&s.out, "event: %s\ndata: %s\n\n", eventType, bytes.TrimRight(data, "\n"))
}

// advance consumes one event from the current hop's stream.
func (s *streamSplicer) advance() {
	if !s.dec.Next() {
		err := s.dec.Err()
		if err == nil && s.terminalSent && !s.stopSent {
			// The hop's wire ended without dispatching message_stop (an
			// unterminated final frame); the spliced stream still completes.
			s.emit("message_stop", []byte(`{"type": "message_stop"}`))
			s.stopSent = true
		}
		s.finish(err)
		return
	}
	evt := s.dec.Event()
	if evt.Type == "message_stop" {
		s.stopSent = true
	}
	switch evt.Type {
	case "ping":
		// Dropped: a held message_start must stay the first event out.
	case "message_start":
		s.startRaw = append([]byte(nil), evt.Data...)
		if !s.isRetryHop {
			s.primaryID = gjson.GetBytes(evt.Data, "message.id").String()
		}
	case "content_block_start":
		s.ensureWireOpen()
		index := int(gjson.GetBytes(evt.Data, "index").Int())
		s.contentSeen = true
		s.open[index] = true
		block := gjson.GetBytes(evt.Data, "content_block")
		b := &rawBlock{start: []byte(block.Raw), blockType: block.Get("type").String()}
		if b.blockType == "fallback" {
			s.sawSeam = true
			s.effModel = block.Get("to.model").String()
		}
		s.blocks[index] = b
		s.blockOrder = append(s.blockOrder, index)
		s.emitReindexed(evt, index)
	case "content_block_delta":
		s.ensureWireOpen()
		index := int(gjson.GetBytes(evt.Data, "index").Int())
		s.contentSeen = true
		s.foldDelta(index, evt.Data)
		s.emitReindexed(evt, index)
	case "content_block_stop":
		s.ensureWireOpen()
		index := int(gjson.GetBytes(evt.Data, "index").Int())
		s.contentSeen = true
		delete(s.open, index)
		s.emitReindexed(evt, index)
	case "message_delta":
		s.handleTerminal(evt)
	default:
		// message_stop and anything unrecognised pass through.
		s.emit(evt.Type, evt.Data)
	}
}

// ensureWireOpen emits the held message_start and any queued seam blocks the
// first time output framing is needed.
func (s *streamSplicer) ensureWireOpen() {
	if s.wireOpen {
		return
	}
	s.wireOpen = true
	start := s.startRaw
	if start == nil {
		return
	}
	if s.isRetryHop && s.primaryID != "" {
		// The envelope keeps the primary's message id no matter which hop
		// opens the wire.
		if patched, err := sjson.SetBytes(start, "message.id", s.primaryID); err == nil {
			start = patched
		}
	}
	s.emit("message_start", start)
	for _, sm := range s.pendingSeams {
		s.emitSeam(sm)
	}
	s.pendingSeams = nil
	s.hopOffset = s.nextIndex
}

func (s *streamSplicer) emitSeam(sm seam) {
	index := s.nextIndex
	s.nextIndex++
	start, _ := json.Marshal(map[string]any{
		"type":          "content_block_start",
		"index":         index,
		"content_block": sm.block(),
	})
	s.emit("content_block_start", start)
	stop, _ := json.Marshal(map[string]any{"type": "content_block_stop", "index": index})
	s.emit("content_block_stop", stop)
}

// emitReindexed forwards a block event, remapping its index into output
// space. Other fields keep their wire bytes.
func (s *streamSplicer) emitReindexed(evt ssestream.Event, wireIndex int) {
	out := s.hopOffset + wireIndex
	if out+1 > s.nextIndex {
		s.nextIndex = out + 1
	}
	data := evt.Data
	if out != wireIndex {
		if patched, err := sjson.SetBytes(data, "index", out); err == nil {
			data = patched
		}
	}
	s.emit(evt.Type, data)
}

func (s *streamSplicer) foldDelta(index int, data []byte) {
	b := s.blocks[index]
	if b == nil {
		return
	}
	delta := gjson.GetBytes(data, "delta")
	switch delta.Get("type").String() {
	case "text_delta":
		b.text.WriteString(delta.Get("text").String())
	case "thinking_delta":
		b.thinking.WriteString(delta.Get("thinking").String())
	case "signature_delta":
		b.signature.WriteString(delta.Get("signature").String())
	case "input_json_delta":
		b.partial.WriteString(delta.Get("partial_json").String())
	case "citations_delta":
		if c := delta.Get("citation"); c.Exists() {
			b.citations = append(b.citations, json.RawMessage(c.Raw))
		}
	}
}

// closeOpenBlocks synthesizes content_block_stop for blocks a refusal cut
// mid-stream, in ascending index order, so the boundary lands on complete
// framing.
func (s *streamSplicer) closeOpenBlocks() {
	if len(s.open) == 0 {
		return
	}
	indices := make([]int, 0, len(s.open))
	for index := range s.open {
		indices = append(indices, index)
	}
	sort.Ints(indices)
	for _, index := range indices {
		data, _ := json.Marshal(map[string]any{"type": "content_block_stop", "index": s.hopOffset + index})
		s.emit("content_block_stop", data)
	}
	s.open = map[int]bool{}
}

// handleTerminal decides what a hop's message_delta means: serve, retry, or
// surface.
func (s *streamSplicer) handleTerminal(evt ssestream.Event) {
	parsed := gjson.ParseBytes(evt.Data)
	refused := parsed.Get("delta.stop_reason").String() == string(anthropic.BetaStopReasonRefusal)

	if !refused {
		s.ensureWireOpen()
		s.terminalSent = true
		if !s.isRetryHop {
			// The original attempt served: the middleware is the identity.
			s.emit(evt.Type, evt.Data)
			return
		}
		s.emit(evt.Type, s.rewriteTerminal(evt.Data, true, false))
		return
	}

	token := parsed.Get("delta.stop_details.fallback_credit_token").String()
	// Token required once content has streamed; a pre-stream refusal retries
	// free — nothing reached the consumer.
	canRetry := token != "" || !s.contentSeen

	if s.next < len(s.fallbacks) && canRetry {
		s.ledger = append(s.ledger, s.hopIterations(evt.Data, false)...)
		s.lastClaimCount = 0
		if parsed.Get("delta.stop_details.fallback_has_prefill_claim").Type == gjson.True {
			claim := s.claimBlocks()
			s.continuation = append(s.continuation, claim...)
			s.lastClaimCount = len(claim)
		}
		s.token = token
		s.held = &heldRefusal{startRaw: s.startRaw, deltaRaw: append([]byte(nil), evt.Data...), stopRaw: s.readStop()}
		s.res.Body.Close()
		s.openNextHop()
		return
	}

	if !s.isRetryHop {
		// Untouched surface: no retry ever engaged, so the original wire
		// bytes pass through, terminal included.
		s.ensureWireOpen()
		s.terminalSent = true
		s.emit(evt.Type, evt.Data)
		return
	}

	// A retry hop's refusal surfaces as the chain's terminal: merged
	// iterations with this hop as the fallback_message completer, and
	// recommended_model stamped null when the wire left it out.
	s.ensureWireOpen()
	s.terminalSent = true
	s.emit(evt.Type, s.rewriteTerminal(evt.Data, true, true))
}

// readStop reads the message_stop that follows a held refusal delta.
func (s *streamSplicer) readStop() []byte {
	for s.dec.Next() {
		evt := s.dec.Event()
		if evt.Type == "ping" {
			continue
		}
		if evt.Type == "message_stop" {
			return append([]byte(nil), evt.Data...)
		}
		break
	}
	data, _ := json.Marshal(map[string]any{"type": "message_stop"})
	return data
}

// rewriteTerminal rebuilds a terminal message_delta: usage.iterations becomes
// the whole-chain ledger (plus this hop, last entry typed fallback_message),
// and optionally recommended_model is stamped null when absent.
func (s *streamSplicer) rewriteTerminal(data []byte, includeSelf bool, stampRecommended bool) []byte {
	iterations := append([]json.RawMessage{}, s.ledger...)
	if includeSelf {
		iterations = append(iterations, s.hopIterations(data, true)...)
	}
	out := data
	if patched, err := sjson.SetRawBytes(out, "usage.iterations", mustMarshal(iterations)); err == nil {
		out = patched
	}
	if stampRecommended && !gjson.GetBytes(out, "delta.stop_details.recommended_model").Exists() {
		if patched, err := sjson.SetRawBytes(out, "delta.stop_details.recommended_model", []byte("null")); err == nil {
			out = patched
		}
	}
	return out
}

// hopIterations builds the current hop's usage.iterations contribution from
// its terminal delta. A hop reporting its own iterations (a server-stitched
// envelope, or a server-tool loop) contributes them verbatim; otherwise one
// entry is synthesized from its delta usage. Every non-final entry is
// (re)typed message — those attempts did not serve; the final hop's last
// entry is the fallback_message completer. The hop's model is stamped only
// when the contribution is unambiguously the hop itself: a single entry. A
// multi-entry array attributes per-iteration usage the wire did not break
// down by model.
func (s *streamSplicer) hopIterations(deltaData []byte, final bool) []json.RawMessage {
	var entries []json.RawMessage
	if reported := gjson.GetBytes(deltaData, "usage.iterations"); reported.IsArray() && len(reported.Array()) > 0 {
		for _, item := range reported.Array() {
			entries = append(entries, json.RawMessage(item.Raw))
		}
	} else {
		entries = []json.RawMessage{synthesizedEntry(deltaData)}
	}
	for i, entry := range entries {
		last := i == len(entries)-1
		typ := gjson.GetBytes(entry, "type").String()
		switch {
		case final && last:
			typ = "fallback_message"
		case typ == "fallback_message":
			// The envelope's serving hop declined the overall request; it is
			// an ordinary attempt in the merged ledger.
			typ = "message"
		case typ == "":
			typ = "message"
		}
		patched, err := sjson.SetBytes(entry, "type", typ)
		if err != nil {
			continue
		}
		// Entries the hop supplied unlabeled can only be the hop's own
		// attempts; a server-stitched envelope labels its merged ledger.
		if !gjson.GetBytes(patched, "model").Exists() {
			if p, err := sjson.SetBytes(patched, "model", s.hopModel); err == nil {
				patched = p
			}
		}
		entries[i] = patched
	}
	return entries
}

// synthesizedEntry derives one ledger entry from a hop that reported no
// iterations: the base usage fields from its terminal delta, zero where the
// delta left them out and null cache_creation — message_start's snapshot is
// a running total, not this attempt's bill.
func synthesizedEntry(deltaData []byte) json.RawMessage {
	delta := gjson.GetBytes(deltaData, "usage")
	entry := []byte(`{}`)
	for _, key := range []string{"input_tokens", "output_tokens", "cache_read_input_tokens", "cache_creation_input_tokens"} {
		raw := "0"
		if value := delta.Get(key); value.Exists() {
			raw = value.Raw
		}
		if patched, err := sjson.SetRawBytes(entry, key, []byte(raw)); err == nil {
			entry = patched
		}
	}
	raw := "null"
	if value := delta.Get("cache_creation"); value.Exists() {
		raw = value.Raw
	}
	if patched, err := sjson.SetRawBytes(entry, "cache_creation", []byte(raw)); err == nil {
		entry = patched
	}
	entry, _ = sjson.SetBytes(entry, "type", "message")
	return entry
}

// claimBlocks reconstructs the refused hop's streamed content for a prefill
// claim. Fallback seam blocks never ride a claim; when the hop was itself a
// stitched envelope only the content after its last seam is the partial, and
// its trailing whitespace is the server's, not the model's.
func (s *streamSplicer) claimBlocks() []json.RawMessage {
	var finalized []json.RawMessage
	for _, index := range s.blockOrder {
		b := s.blocks[index]
		if b.blockType == "fallback" {
			finalized = nil // restart after the seam
			continue
		}
		finalized = append(finalized, b.finalize())
	}
	if s.sawSeam && len(finalized) > 0 {
		last := finalized[len(finalized)-1]
		if gjson.GetBytes(last, "type").String() == "text" {
			text := gjson.GetBytes(last, "text").String()
			if trimmed := strings.TrimRight(text, " \t\r\n"); trimmed != text {
				if patched, err := sjson.SetBytes(last, "text", trimmed); err == nil {
					finalized[len(finalized)-1] = patched
				}
			}
		}
	}
	return finalized
}

// hopBody builds one retry attempt: the original body with the entry merged
// over it, the credit token attached, and the continuation appended as one
// trailing assistant turn.
func (s *streamSplicer) hopBody(entry anthropic.BetaFallbackParam, token string, continuation []json.RawMessage) ([]byte, error) {
	merged, _, err := mergeFallback(s.body, entry, token)
	if err != nil {
		return nil, err
	}
	if len(continuation) == 0 {
		return merged, nil
	}
	turn, err := json.Marshal(map[string]any{"role": "assistant", "content": continuation})
	if err != nil {
		return nil, err
	}
	var messages []json.RawMessage
	gjson.GetBytes(merged, "messages").ForEach(func(_, msg gjson.Result) bool {
		messages = append(messages, json.RawMessage(msg.Raw))
		return true
	})
	messages = append(messages, turn)
	return sjson.SetRawBytes(merged, "messages", mustMarshal(messages))
}

// openNextHop walks the remaining chain after a held refusal. Failed hops
// are skipped with the token and continuation intact. A 400 must not
// surface a refusal the app never saw, so the last entry descends a ladder,
// shedding one ingredient per rung: first the echoed partial goes (the
// appended assistant turn may be all the server rejected), then the token
// (it is not consumed by a rejected redemption; only the reprice is lost).
// A chain that never engages degrades to the held refusal.
func (s *streamSplicer) openNextHop() {
	lastStatus := 0
	var lastEntry anthropic.BetaFallbackParam
	for s.next < len(s.fallbacks) {
		entry := s.fallbacks[s.next]
		index := s.next
		s.next++
		if s.state != nil {
			s.state.SetIndex(index)
		}
		res, err := s.tryHop(entry, s.token, s.continuation)
		if err != nil {
			s.finish(err)
			return
		}
		if res != nil {
			s.engage(res, entry)
			return
		}
		lastStatus = s.lastHopStatus
		lastEntry = entry
	}

	// The server adjudicates a claim; its 400 may reject only the appended
	// turn, so the last entry gets one attempt without this hop's partial.
	if lastStatus == http.StatusBadRequest && s.lastClaimCount > 0 {
		trimmed := s.continuation[:len(s.continuation)-s.lastClaimCount]
		res, err := s.tryHop(lastEntry, s.token, trimmed)
		if err != nil {
			s.finish(err)
			return
		}
		if res != nil {
			s.continuation = trimmed
			s.lastClaimCount = 0
			s.engage(res, lastEntry)
			return
		}
		lastStatus = s.lastHopStatus
	}

	// Final rung: same entry, no echoed partial, no token. A rejected
	// redemption does not consume the token, so dropping it costs only the
	// reprice — never the conversation.
	if lastStatus == http.StatusBadRequest && s.token != "" {
		trimmed := s.continuation[:len(s.continuation)-s.lastClaimCount]
		res, err := s.tryHop(lastEntry, "", trimmed)
		if err != nil {
			s.finish(err)
			return
		}
		if res != nil {
			s.token = ""
			s.continuation = trimmed
			s.lastClaimCount = 0
			s.engage(res, lastEntry)
			return
		}
		lastStatus = s.lastHopStatus
	}

	s.degrade(lastStatus, lastEntry)
}

// tryHop issues one hop attempt. It returns the response when the hop
// engaged (a readable SSE stream), nil when the hop failed and the chain
// should move on, or an error for unrecoverable problems building the
// attempt.
func (s *streamSplicer) tryHop(entry anthropic.BetaFallbackParam, token string, continuation []json.RawMessage) (*http.Response, error) {
	payload, err := s.hopBody(entry, token, continuation)
	if err != nil {
		return nil, err
	}
	res, err := s.send(payload)
	if err != nil {
		s.lastHopStatus = 0
		return nil, nil
	}
	if !isEventStream(res) {
		s.lastHopStatus = res.StatusCode
		res.Body.Close()
		return nil, nil
	}
	return res, nil
}

// engage switches the splicer onto a hop that returned a live stream,
// emitting (or queueing) the fallback boundary between the previous
// content-bearing model and this hop.
func (s *streamSplicer) engage(res *http.Response, entry anthropic.BetaFallbackParam) {
	sm := seam{from: s.effModel, to: string(entry.Model)}
	if s.wireOpen {
		s.closeOpenBlocks()
		s.emitSeam(sm)
		s.hopOffset = s.nextIndex
	} else {
		s.pendingSeams = append(s.pendingSeams, sm)
	}

	s.hopModel = string(entry.Model)
	s.effModel = s.hopModel
	s.isRetryHop = true
	s.contentSeen = false
	s.sawSeam = false
	s.open = map[int]bool{}
	s.blocks = map[int]*rawBlock{}
	s.blockOrder = nil
	s.startRaw = nil

	s.res = res
	s.dec = ssestream.NewDecoder(res)
}

// degrade surfaces the held refusal after every remaining hop failed: its
// terminal events replay with the chain's ledger and a recommended_model
// verdict — the unreachable entry on a rate limit, null otherwise.
func (s *streamSplicer) degrade(lastStatus int, lastEntry anthropic.BetaFallbackParam) {
	held := s.held
	if held == nil {
		s.finish(nil)
		return
	}
	if !s.wireOpen {
		s.wireOpen = true
		start := held.startRaw
		if start != nil {
			if s.primaryID != "" {
				if patched, err := sjson.SetBytes(start, "message.id", s.primaryID); err == nil {
					start = patched
				}
			}
			s.emit("message_start", start)
		}
		for _, sm := range s.pendingSeams {
			s.emitSeam(sm)
		}
		s.pendingSeams = nil
	}

	recommended := []byte("null")
	if lastStatus == http.StatusTooManyRequests {
		recommended, _ = json.Marshal(string(lastEntry.Model))
	}
	data := held.deltaRaw
	existing := gjson.GetBytes(data, "delta.stop_details.recommended_model")
	if !existing.Exists() || existing.Type == gjson.Null {
		if patched, err := sjson.SetRawBytes(data, "delta.stop_details.recommended_model", recommended); err == nil {
			data = patched
		}
	}
	if len(s.ledger) > 0 {
		if patched, err := sjson.SetRawBytes(data, "usage.iterations", mustMarshal(s.ledger)); err == nil {
			data = patched
		}
	}
	s.emit("message_delta", data)
	s.emit("message_stop", held.stopRaw)
	s.finish(nil)
}

func mustMarshal(v any) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		return []byte("null")
	}
	return data
}
