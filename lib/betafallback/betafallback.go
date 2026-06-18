// Package betafallback retries refused Messages requests down a fallback
// chain, client-side. Beta surface; may change.
package betafallback

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// BetaFallbackState keeps the requests that share it on the model that
// accepted, so a conversation never re-asks a model that already refused.
// Use one zero-value state per conversation, passed with
// [WithBetaFallbackState] on each of its requests. Safe for concurrent use;
// when requests race, the last pin written wins.
type BetaFallbackState struct {
	// 1-based; zero targets the original request params.
	pinned atomic.Int64
}

// Index reports the pinned fallback's zero-based index, or -1 when no
// fallback is pinned. Use it with SetIndex to persist a pin across
// processes. The pin is positional: against a chain that may be reordered,
// persist the model instead.
func (s *BetaFallbackState) Index() int { return int(s.pinned.Load()) - 1 }

// SetIndex pins the fallback at zero-based index i. Negative resets to the
// original params; a pin past the configured chain errors when used.
func (s *BetaFallbackState) SetIndex(i int) {
	s.pinned.Store(int64(max(i, -1) + 1))
}

type betaFallbackStateKey struct{}

// WithBetaFallbackState returns a request option carrying state for [BetaRefusalFallbackMiddleware].
// Pass the same state on every request that should share the pin —
// typically the turns of one conversation.
func WithBetaFallbackState(state *BetaFallbackState) option.RequestOption {
	return requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
		ctx := context.WithValue(cfg.Request.Context(), betaFallbackStateKey{}, state)
		cfg.Request = cfg.Request.WithContext(ctx)
		return nil
	})
}

// BetaRefusalFallbackMiddleware retries a refused Messages request on each
// model of fallbacks in turn, so a refusal costs a retry instead of
// surfacing to the caller. Only stop_reason "refusal" triggers a retry —
// never transport or API errors. There are deliberately no hooks: what
// happened is readable from the response, and middleware registered after
// this one observes every attempt.
//
// Streaming requests are retried in place: the retry's events are spliced
// onto the open stream behind a `fallback` boundary block, with one
// message_start, monotonic block indices, and a terminal usage.iterations
// ledger covering every hop — the same framing the server's own fallback
// chain produces. A mid-stream refusal retries only when it minted a
// fallback_credit_token (a pre-stream refusal retries either way), and when
// the refusal advertises a prefill claim the retry continues from the
// refused partial instead of starting over.
//
// A non-streaming message served by a fallback retry opens with one
// `fallback` content block per model boundary, matching the streaming
// splice's block shape; an exhausted chain's final refusal is returned
// verbatim.
//
//	client := anthropic.NewClient(
//		option.WithMiddleware(betafallback.BetaRefusalFallbackMiddleware(
//			[]anthropic.BetaFallbackParam{{Model: anthropic.ModelClaudeOpus4_5}},
//		)),
//	)
//
//	conversation := betafallback.WithBetaFallbackState(&betafallback.BetaFallbackState{})
//	message, err := client.Beta.Messages.New(ctx, params, conversation)
//
// Each fallback merges over the original params, overriding only the fields
// it sets. An exhausted chain returns the last refusal as a normal response;
// an empty chain disables the middleware.
//
// Only the beta Messages surface is handled (client.Beta.Messages, which
// routes ?beta=true); plain client.Messages requests pass through untouched.
// With the Bedrock or Vertex options, register this middleware first — it
// needs the Anthropic-shaped request, and after the platform transform it
// stands down.
//
// Handled requests are opted into the fallback-credit beta so retries
// redeem the refusal's fallback_credit_token. A request using server-side
// fallbacks errors: only one of the chains can adjudicate refusals.
//
// The client retry layer treats the chain as one try; a retried request
// re-enters it from the original params or the state's pin.
func BetaRefusalFallbackMiddleware(fallbacks []anthropic.BetaFallbackParam) option.Middleware {
	// An empty chain can never act on a refusal; stand down rather than opt
	// requests into the credit beta.
	if len(fallbacks) == 0 {
		return func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			return next(req)
		}
	}
	// A zero-model fallback merges to a no-op that silently re-sends the
	// refused request; fail loudly instead.
	var initErr error
	for i := range fallbacks {
		if fallbacks[i].Model == "" {
			initErr = fmt.Errorf("betafallback: fallbacks[%d] has no model", i)
			break
		}
	}
	return func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		// Platform middlewares run below user middleware, so this path is
		// canonical; the suffix match tolerates gateway base-URL prefixes.
		if req.Method != http.MethodPost || !strings.HasSuffix(req.URL.Path, "/v1/messages") || req.Body == nil {
			return next(req)
		}
		// The middleware is a beta feature: only the beta surface
		// (client.Beta.Messages, which requests v1/messages?beta=true) is
		// handled. The plain surface passes through untouched — no credit
		// beta armed, no retry on refusal.
		if req.URL.Query().Get("beta") != "true" {
			return next(req)
		}
		mediaType, _, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))
		if !strings.Contains(mediaType, "application/json") && !strings.HasSuffix(mediaType, "+json") {
			return next(req)
		}

		orig, err := io.ReadAll(req.Body)
		req.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("betafallback: reading request body: %w", err)
		}

		// Unparseable bodies pass through for the server to reject.
		var body map[string]json.RawMessage
		if json.Unmarshal(orig, &body) != nil {
			return next(requestWithBody(req, orig))
		}
		// model+messages+max_tokens is the create shape (token counting lacks
		// max_tokens; the platform transforms strip model).
		for _, key := range []string{"model", "messages", "max_tokens"} {
			if _, ok := body[key]; !ok {
				return next(requestWithBody(req, orig))
			}
		}
		// The two chains cannot adjudicate refusals together; fail loudly
		// rather than letting them race.
		if _, ok := body["fallbacks"]; ok {
			return nil, fmt.Errorf("Sending the `fallbacks:` request param is not supported when using the `BetaRefusalFallbackMiddleware` middleware. You should either remove the middleware and send `fallbacks:` with the `server-side-fallback-2026-06-01` beta header to let the API handle refusal fallbacks, or omit the `fallbacks:` param if you'd like the `BetaRefusalFallbackMiddleware` middleware to handle fallbacks on the client side.")
		}
		if string(bytes.TrimSpace(body["stream"])) == "true" {
			// A bad chain fails only requests the middleware would act on.
			if initErr != nil {
				return nil, initErr
			}
			state, _ := req.Context().Value(betaFallbackStateKey{}).(*BetaFallbackState)
			if state != nil && state.Index() >= len(fallbacks) {
				return nil, fmt.Errorf("betafallback: pinned fallback %d is out of range for a chain of %d; was the state shared with a different middleware?", state.Index(), len(fallbacks))
			}
			return handleStreaming(req, next, orig, body, fallbacks, state)
		}
		// A bad chain fails only requests the middleware would act on.
		if initErr != nil {
			return nil, initErr
		}

		state, _ := req.Context().Value(betaFallbackStateKey{}).(*BetaFallbackState)

		// index targets fallbacks[index]; -1 targets the original params.
		index := -1
		if state != nil {
			index = state.Index()
			if index >= len(fallbacks) {
				return nil, fmt.Errorf("betafallback: pinned fallback %d is out of range for a chain of %d; was the state shared with a different middleware?", index, len(fallbacks))
			}
		}

		// History from an earlier fallback contains blocks the server
		// rejects; trim them before sending.
		orig = trimHistory(body, orig)

		attemptRequest := func(body []byte) *http.Request {
			r := requestWithBody(req, body)
			ensureFallbackCreditBeta(r.Header)
			return r
		}

		attempt := orig
		if index >= 0 {
			if attempt, _, err = mergeFallback(body, fallbacks[index], ""); err != nil {
				return nil, err
			}
		}
		res, err := next(attemptRequest(attempt))

		// Seams mirror the streaming splice's block shape: one fallback
		// boundary block per retried hop, prepended to the served message's
		// content. They are recorded optimistically before each retry is
		// sent; prependSeams' status/stop_reason guards ensure they attach
		// only to a served 200, so a hop that errors or the exhausted
		// chain's final refusal never carries them.
		var seams []seam
		fromModel := gjson.GetBytes(attempt, "model").String()

		for err == nil && index < len(fallbacks)-1 {
			refused, creditToken, category, perr := refusedMessage(res)
			if perr != nil {
				return nil, perr
			}
			if !refused {
				break
			}

			index++
			fallback := fallbacks[index]
			if state != nil {
				state.SetIndex(index)
			}
			seams = append(seams, seam{from: fromModel, to: string(fallback.Model), category: category})
			fromModel = string(fallback.Model)
			var tokenSent bool
			if attempt, tokenSent, err = mergeFallback(body, fallback, creditToken); err != nil {
				return nil, err
			}
			res, err = next(attemptRequest(attempt))

			// A rejected redemption recovers with one tokenless resend.
			if err == nil && tokenSent && res.StatusCode == http.StatusBadRequest {
				res.Body.Close()
				if attempt, _, err = mergeFallback(body, fallback, ""); err != nil {
					return nil, err
				}
				res, err = next(attemptRequest(attempt))
			}
		}
		if err == nil && len(seams) > 0 {
			if perr := prependSeams(res, seams); perr != nil {
				return nil, perr
			}
		}
		return res, err
	}
}

// prependSeams rewrites a served message's content to open with the fallback
// boundary blocks, mirroring the streaming splice's block shape. An
// exhausted chain's refusal — and anything that isn't an inspectable served
// message — is left as written.
func prependSeams(res *http.Response, seams []seam) error {
	if res.StatusCode != http.StatusOK || res.Header.Get("Content-Encoding") != "" {
		return nil
	}
	buf, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return fmt.Errorf("betafallback: reading response body: %w", err)
	}
	res.Body = io.NopCloser(bytes.NewReader(buf))
	if gjson.GetBytes(buf, "stop_reason").String() == string(anthropic.BetaStopReasonRefusal) {
		return nil
	}
	content := gjson.GetBytes(buf, "content")
	if !content.IsArray() {
		return nil
	}
	blocks := make([]json.RawMessage, 0, len(seams)+len(content.Array()))
	for _, sm := range seams {
		blocks = append(blocks, sm.block())
	}
	for _, b := range content.Array() {
		blocks = append(blocks, json.RawMessage(b.Raw))
	}
	if patched, err := sjson.SetRawBytes(buf, "content", mustMarshal(blocks)); err == nil {
		res.Body = io.NopCloser(bytes.NewReader(patched))
		res.ContentLength = int64(len(patched))
		if res.Header == nil {
			res.Header = http.Header{}
		}
		res.Header.Set("Content-Length", strconv.Itoa(len(patched)))
	}
	return nil
}

// trimHistory rewrites body["messages"] through trimFallbackTurns and
// returns the request bytes to send. On no change — or if the body cannot
// be rebuilt — the original bytes are returned untouched.
func trimHistory(body map[string]json.RawMessage, orig []byte) []byte {
	trimmed, changed := trimFallbackTurns(body["messages"])
	if !changed {
		return orig
	}
	body["messages"] = trimmed
	rebuilt, err := json.Marshal(body)
	if err != nil {
		return orig
	}
	return rebuilt
}

// trimFallbackTurns removes content the server would reject from assistant
// turns that contain a fallback block: the fallback block itself (only
// accepted under a beta this middleware does not send), and everything
// before it that belongs to the model that refused — thinking, connector
// text, and tool calls that never got a result. Blocks after the fallback
// block are what the serving model produced and stay as written. Turns
// without a fallback block are never touched.
func trimFallbackTurns(messages json.RawMessage) (json.RawMessage, bool) {
	parsed := gjson.ParseBytes(messages)
	if !parsed.IsArray() {
		return messages, false
	}
	// Tool calls whose result appears anywhere in the history are kept.
	resolved := map[string]bool{}
	parsed.ForEach(func(_, msg gjson.Result) bool {
		content := msg.Get("content")
		if content.IsArray() {
			content.ForEach(func(_, block gjson.Result) bool {
				typ := block.Get("type").String()
				if typ == "tool_result" || strings.HasSuffix(typ, "_tool_result") {
					resolved[block.Get("tool_use_id").String()] = true
				}
				return true
			})
		}
		return true
	})

	changed := false
	var msgs []json.RawMessage
	parsed.ForEach(func(_, msg gjson.Result) bool {
		content := msg.Get("content")
		hasFallback := false
		if msg.Get("role").String() == "assistant" && content.IsArray() {
			content.ForEach(func(_, block gjson.Result) bool {
				if block.Get("type").String() == "fallback" {
					hasFallback = true
				}
				return !hasFallback
			})
		}
		if !hasFallback {
			msgs = append(msgs, json.RawMessage(msg.Raw))
			return true
		}
		all := content.Array()
		lastFallback := -1
		for i, block := range all {
			if block.Get("type").String() == "fallback" {
				lastFallback = i
			}
		}
		var blocks []json.RawMessage
		for i, block := range all {
			typ := block.Get("type").String()
			switch {
			case typ == "fallback":
			case i < lastFallback && (typ == "thinking" || typ == "redacted_thinking" || typ == "connector_text"):
			case i < lastFallback && typ == "tool_use":
			case i < lastFallback && typ == "server_tool_use" && !resolved[block.Get("id").String()]:
			default:
				blocks = append(blocks, json.RawMessage(block.Raw))
			}
		}
		changed = true
		rebuilt, err := sjson.SetRawBytes([]byte(msg.Raw), "content", mustMarshal(blocks))
		if err != nil {
			msgs = append(msgs, json.RawMessage(msg.Raw))
			return true
		}
		msgs = append(msgs, rebuilt)
		return true
	})
	if !changed {
		return messages, false
	}
	return mustMarshal(msgs), true
}

// requestWithBody clones req around a replayable body.
func requestWithBody(req *http.Request, body []byte) *http.Request {
	r := req.Clone(req.Context())
	r.Body = io.NopCloser(bytes.NewReader(body))
	r.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(body)), nil }
	r.ContentLength = int64(len(body))
	return r
}

// ensureFallbackCreditBeta opts a request into the fallback-credit beta.
func ensureFallbackCreditBeta(h http.Header) {
	if hasFallbackCreditBeta(h) {
		return
	}
	h.Add("anthropic-beta", string(anthropic.AnthropicBetaFallbackCredit2026_06_01))
}

// hasFallbackCreditBeta reports whether the headers already carry the
// fallback-credit beta.
func hasFallbackCreditBeta(h http.Header) bool {
	for _, value := range h.Values("anthropic-beta") {
		for _, beta := range strings.Split(value, ",") {
			if strings.TrimSpace(beta) == string(anthropic.AnthropicBetaFallbackCredit2026_06_01) {
				return true
			}
		}
	}
	return false
}

// mergeFallback overlays fallback's fields and the credit token onto the
// original body, reporting whether the token was included. The body stays a
// raw map: a typed round-trip would drop unknown fields and rewrite
// encodings that prompt caching and the token's body match depend on.
func mergeFallback(body map[string]json.RawMessage, fallback anthropic.BetaFallbackParam, creditToken string) ([]byte, bool, error) {
	overrides, err := json.Marshal(fallback)
	if err != nil {
		return nil, false, fmt.Errorf("betafallback: marshaling fallback params: %w", err)
	}
	var overlay map[string]json.RawMessage
	if err := json.Unmarshal(overrides, &overlay); err != nil {
		return nil, false, fmt.Errorf("betafallback: parsing fallback params: %w", err)
	}

	merged := make(map[string]json.RawMessage, len(body)+len(overlay)+1)
	for k, v := range body {
		merged[k] = v
	}
	for k, v := range overlay {
		merged[k] = v
	}

	// A token cannot combine with server-side fallback params, and a stale
	// token from the original body cannot match this retry.
	delete(merged, "fallback")
	delete(merged, "fallbacks")
	if creditToken != "" {
		token, _ := json.Marshal(creditToken)
		merged["fallback_credit_token"] = token
	} else {
		delete(merged, "fallback_credit_token")
	}
	out, err := json.Marshal(merged)
	return out, creditToken != "", err
}

// refusedMessage reports whether res is a message that stopped with a refusal,
// its credit token, and its policy category (nil when not surfaced), restoring
// res.Body for the caller.
func refusedMessage(res *http.Response) (refused bool, creditToken string, category any, err error) {
	if res.StatusCode != http.StatusOK {
		return false, "", nil, nil
	}
	// A still-encoded body (caller-set Accept-Encoding) can't be inspected.
	if res.Header.Get("Content-Encoding") != "" {
		return false, "", nil, nil
	}
	buf, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return false, "", nil, fmt.Errorf("betafallback: reading response body: %w", err)
	}
	res.Body = io.NopCloser(bytes.NewReader(buf))
	if gjson.GetBytes(buf, "stop_reason").String() != string(anthropic.BetaStopReasonRefusal) {
		return false, "", nil, nil
	}
	if cat := gjson.GetBytes(buf, "stop_details.category"); cat.Type == gjson.String {
		category = cat.String()
	}
	return true, gjson.GetBytes(buf, "stop_details.fallback_credit_token").String(), category, nil
}
