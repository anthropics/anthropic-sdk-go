package betafallback_test

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/lib/betafallback"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// scriptedTransport returns canned JSON responses in order and records every
// request body it sees. Statuses and headers parallel responses; missing
// entries are 200.
type scriptedTransport struct {
	t         *testing.T
	responses []string
	statuses  []int
	headers   []http.Header
	bodies    []map[string]any
	betas     [][]string
}

func (s *scriptedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
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
	header := http.Header{"Content-Type": []string{"application/json"}}
	if len(s.headers) > 0 {
		maps.Copy(header, s.headers[0])
		s.headers = s.headers[1:]
	}
	return &http.Response{
		StatusCode: status,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(next)),
	}, nil
}

func (s *scriptedTransport) models() []string {
	models := make([]string, len(s.bodies))
	for i, body := range s.bodies {
		models[i], _ = body["model"].(string)
	}
	return models
}

func messageResponse(model string) string {
	return fmt.Sprintf(`{
		"id": "msg_1", "type": "message", "role": "assistant", "model": %q,
		"content": [], "stop_reason": "end_turn", "stop_sequence": null,
		"usage": {"input_tokens": 1, "output_tokens": 1}
	}`, model)
}

func refusalResponse(model string, creditToken any) string {
	token, _ := json.Marshal(creditToken)
	return fmt.Sprintf(`{
		"id": "msg_1", "type": "message", "role": "assistant", "model": %q,
		"content": [], "stop_reason": "refusal", "stop_sequence": null,
		"stop_details": {"type": "refusal", "category": null, "explanation": null,
			"fallback_credit_token": %s, "recommended_model": null},
		"usage": {"input_tokens": 1, "output_tokens": 1}
	}`, model, token)
}

func fallbackTestClient(t *testing.T, responses []string, middleware option.Middleware) (anthropic.Client, *scriptedTransport) {
	transport := &scriptedTransport{t: t, responses: responses}
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{Transport: transport}),
		option.WithMaxRetries(0),
		option.WithMiddleware(middleware),
	)
	return client, transport
}

var fallbackTestParams = anthropic.BetaMessageNewParams{
	Model:     "primary-model",
	MaxTokens: 1024,
	Messages:  []anthropic.BetaMessageParam{anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("hi"))},
}

func TestRefusalFallbackMiddlewareRetriesWithFallbackParamsAndCreditToken(t *testing.T) {
	client, transport := fallbackTestClient(t,
		[]string{refusalResponse("primary-model", "credit-token"), messageResponse("fallback-model")},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	message, err := client.Beta.Messages.New(context.Background(), fallbackTestParams)
	require.NoError(t, err)
	assert.Equal(t, anthropic.Model("fallback-model"), message.Model)
	assert.Equal(t, anthropic.BetaStopReasonEndTurn, message.StopReason)
	assert.Equal(t, []string{"primary-model", "fallback-model"}, transport.models())
	assert.Equal(t, "credit-token", transport.bodies[1]["fallback_credit_token"])
	_, hasToken := transport.bodies[0]["fallback_credit_token"]
	assert.False(t, hasToken)
}

func TestRefusalFallbackMiddlewarePinsTheConversationToTheAcceptedFallback(t *testing.T) {
	client, transport := fallbackTestClient(t,
		[]string{refusalResponse("primary-model", nil), messageResponse("fallback-model"), messageResponse("fallback-model")},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	conversation := betafallback.WithBetaFallbackState(&betafallback.BetaFallbackState{})
	_, err := client.Beta.Messages.New(context.Background(), fallbackTestParams, conversation)
	require.NoError(t, err)

	// The follow-up goes straight to the pinned fallback in a single request.
	_, err = client.Beta.Messages.New(context.Background(), fallbackTestParams, conversation)
	require.NoError(t, err)
	assert.Equal(t, []string{"primary-model", "fallback-model", "fallback-model"}, transport.models())
}

func TestRefusalFallbackMiddlewareKeepsSeparateConversationsIndependent(t *testing.T) {
	client, transport := fallbackTestClient(t,
		[]string{refusalResponse("primary-model", nil), messageResponse("fallback-model"), messageResponse("primary-model")},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	_, err := client.Beta.Messages.New(context.Background(), fallbackTestParams,
		betafallback.WithBetaFallbackState(&betafallback.BetaFallbackState{}))
	require.NoError(t, err)

	_, err = client.Beta.Messages.New(context.Background(), fallbackTestParams,
		betafallback.WithBetaFallbackState(&betafallback.BetaFallbackState{}))
	require.NoError(t, err)
	assert.Equal(t, []string{"primary-model", "fallback-model", "primary-model"}, transport.models())
}

func TestRefusalFallbackMiddlewareLeavesAcceptedRequestsUntouched(t *testing.T) {
	client, transport := fallbackTestClient(t,
		[]string{messageResponse("primary-model")},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	message, err := client.Beta.Messages.New(context.Background(), fallbackTestParams,
		betafallback.WithBetaFallbackState(&betafallback.BetaFallbackState{}))
	require.NoError(t, err)
	assert.Equal(t, anthropic.Model("primary-model"), message.Model)
	require.Len(t, transport.bodies, 1)
	_, hasToken := transport.bodies[0]["fallback_credit_token"]
	assert.False(t, hasToken)
}

func TestRefusalFallbackMiddlewareReportsEachHopThroughTheChain(t *testing.T) {
	client, transport := fallbackTestClient(t,
		[]string{
			refusalResponse("primary-model", "token-1"),
			refusalResponse("fallback-1", "token-2"),
			messageResponse("fallback-2"),
		},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-1"}, {Model: "fallback-2"}},
		),
	)

	message, err := client.Beta.Messages.New(context.Background(), fallbackTestParams)
	require.NoError(t, err)
	assert.Equal(t, anthropic.Model("fallback-2"), message.Model)
	assert.Equal(t, []string{"primary-model", "fallback-1", "fallback-2"}, transport.models())
	assert.Equal(t, "token-1", transport.bodies[1]["fallback_credit_token"])
	assert.Equal(t, "token-2", transport.bodies[2]["fallback_credit_token"])
}

func TestRefusalFallbackMiddlewareReturnsTheLastRefusalWhenTheChainIsExhausted(t *testing.T) {
	client, transport := fallbackTestClient(t,
		[]string{refusalResponse("primary-model", nil), refusalResponse("fallback-model", nil)},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	message, err := client.Beta.Messages.New(context.Background(), fallbackTestParams)
	require.NoError(t, err)
	assert.Equal(t, anthropic.BetaStopReasonRefusal, message.StopReason)
	assert.Equal(t, anthropic.Model("fallback-model"), message.Model)
	assert.Len(t, transport.bodies, 2)
}

func TestRefusalFallbackMiddlewareAppliesFallbackOverridesAndPreservesOtherFields(t *testing.T) {
	client, transport := fallbackTestClient(t,
		[]string{refusalResponse("primary-model", nil), messageResponse("fallback-model")},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{
				Model:     "fallback-model",
				MaxTokens: anthropic.Int(2048),
				Thinking: anthropic.BetaFallbackParamThinkingUnion{
					OfDisabled: &anthropic.BetaThinkingConfigDisabledParam{},
				},
			}},
		),
	)

	params := fallbackTestParams
	params.Temperature = anthropic.Float(0.5)
	_, err := client.Beta.Messages.New(context.Background(), params)
	require.NoError(t, err)

	require.Len(t, transport.bodies, 2)
	retry := transport.bodies[1]
	assert.Equal(t, "fallback-model", retry["model"])
	assert.Equal(t, float64(2048), retry["max_tokens"])
	assert.Equal(t, map[string]any{"type": "disabled"}, retry["thinking"])
	assert.Equal(t, 0.5, retry["temperature"])
	assert.Equal(t, transport.bodies[0]["messages"], retry["messages"])
}

func TestBetaFallbackStateIsConcurrencySafe(t *testing.T) {
	state := &betafallback.BetaFallbackState{}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			for n := 0; n < 100; n++ {
				state.SetIndex(i % 2)
				_ = state.Index()
			}
		}()
	}
	wg.Wait()
	assert.Contains(t, []int{0, 1}, state.Index())
}

func TestRefusalFallbackMiddlewareErrorsOnAPinPastTheConfiguredChain(t *testing.T) {
	// A pin past the chain means the state was shared with a different
	// middleware; erroring beats silently picking another model.
	client, _ := fallbackTestClient(t,
		[]string{messageResponse("fallback-model")},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	stale := &betafallback.BetaFallbackState{}
	stale.SetIndex(1)
	_, err := client.Beta.Messages.New(context.Background(), fallbackTestParams,
		betafallback.WithBetaFallbackState(stale))
	require.ErrorContains(t, err, "out of range")
}

func TestRefusalFallbackMiddlewareSkipsRequestsItDoesNotApplyTo(t *testing.T) {
	middleware := betafallback.BetaRefusalFallbackMiddleware(
		[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
	)

	for _, test := range []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{"streaming request", http.MethodPost, "/v1/messages", `{"model": "primary-model", "stream": true}`},
		{"count_tokens shape", http.MethodPost, "/v1/messages/count_tokens", `{"model": "primary-model", "messages": []}`},
		{"non-POST", http.MethodGet, "/v1/messages", `{"model": "primary-model"}`},
		{"non-JSON body", http.MethodPost, "/v1/messages", `not json`},
	} {
		t.Run(test.name, func(t *testing.T) {
			calls := 0
			req, err := http.NewRequest(test.method, "https://api.anthropic.com"+test.path, strings.NewReader(test.body))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			res, err := middleware(req, func(req *http.Request) (*http.Response, error) {
				calls++
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				assert.Equal(t, test.body, string(body))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(refusalResponse("primary-model", nil))),
				}, nil
			})
			require.NoError(t, err)
			assert.Equal(t, 1, calls)
			assert.NotNil(t, res)
		})
	}
}

func TestRefusalFallbackMiddlewareErrorsWhenAFallbackHasNoModel(t *testing.T) {
	// A zero-model fallback merges to a no-op that silently re-sends the
	// refused request.
	middleware := betafallback.BetaRefusalFallbackMiddleware(
		[]anthropic.BetaFallbackParam{{MaxTokens: anthropic.Int(2048)}},
	)

	req, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages?beta=true",
		strings.NewReader(`{"model": "primary-model", "messages": [], "max_tokens": 16}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	_, err = middleware(req, func(*http.Request) (*http.Response, error) {
		t.Fatal("request should not be sent")
		return nil, nil
	})
	require.ErrorContains(t, err, "fallbacks[0] has no model")

	// Requests the middleware would not handle are unaffected by the bad config.
	other, err := http.NewRequest(http.MethodGet, "https://api.anthropic.com/v1/models", nil)
	require.NoError(t, err)
	res, err := middleware(other, func(*http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("{}"))}, nil
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	// Streaming requests are handled too; the bad chain fails them just as
	// loudly.
	streaming, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages?beta=true",
		strings.NewReader(`{"model": "primary-model", "messages": [], "max_tokens": 16, "stream": true}`))
	require.NoError(t, err)
	streaming.Header.Set("Content-Type", "application/json")
	_, err = middleware(streaming, func(*http.Request) (*http.Response, error) {
		t.Fatal("request should not be sent")
		return nil, nil
	})
	require.ErrorContains(t, err, "fallbacks[0] has no model")
}

func TestRefusalFallbackMiddlewareWithAnEmptyChainDoesNotOptIntoTheBeta(t *testing.T) {
	// The opt-in changes what the server returns; an empty chain gets no
	// benefit from it.
	client, transport := fallbackTestClient(t,
		[]string{messageResponse("primary-model")},
		betafallback.BetaRefusalFallbackMiddleware(nil),
	)

	_, err := client.Beta.Messages.New(context.Background(), fallbackTestParams)
	require.NoError(t, err)
	require.Len(t, transport.betas, 1)
	assert.NotContains(t, transport.betas[0], string(anthropic.AnthropicBetaFallbackCredit2026_06_01))
}

func TestRefusalFallbackMiddlewareDropsAStaleCreditToken(t *testing.T) {
	// A token on the original params was minted for an earlier refusal; it
	// must not survive onto a retry whose refusal supplied none.
	client, transport := fallbackTestClient(t,
		[]string{refusalResponse("primary-model", nil), messageResponse("fallback-model")},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	params := fallbackTestParams
	params.FallbackCreditToken = anthropic.String("stale-token")
	_, err := client.Beta.Messages.New(context.Background(), params)
	require.NoError(t, err)

	require.Len(t, transport.bodies, 2)
	assert.Equal(t, "stale-token", transport.bodies[0]["fallback_credit_token"])
	_, hasToken := transport.bodies[1]["fallback_credit_token"]
	assert.False(t, hasToken)
}

func TestRefusalFallbackMiddlewareRehydratesAPersistedPin(t *testing.T) {
	// A pin persisted by a previous process is restored with SetIndex and
	// readable with Index after the turn.
	client, transport := fallbackTestClient(t,
		[]string{refusalResponse("primary-model", nil), messageResponse("fallback-model"), messageResponse("fallback-model")},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	state := &betafallback.BetaFallbackState{}
	assert.Equal(t, -1, state.Index())
	_, err := client.Beta.Messages.New(context.Background(), fallbackTestParams,
		betafallback.WithBetaFallbackState(state))
	require.NoError(t, err)
	assert.Equal(t, 0, state.Index())

	rehydrated := &betafallback.BetaFallbackState{}
	rehydrated.SetIndex(state.Index())
	_, err = client.Beta.Messages.New(context.Background(), fallbackTestParams,
		betafallback.WithBetaFallbackState(rehydrated))
	require.NoError(t, err)
	assert.Equal(t, []string{"primary-model", "fallback-model", "fallback-model"}, transport.models())

	reset := &betafallback.BetaFallbackState{}
	reset.SetIndex(-5)
	assert.Equal(t, -1, reset.Index())
}

func TestRefusalFallbackMiddlewareSendsTheTokenWithAllOverrides(t *testing.T) {
	// The server adjudicates what a redemption permits; the client always
	// attaches the token.
	client, transport := fallbackTestClient(t,
		[]string{
			refusalResponse("primary-model", "token-1"),
			refusalResponse("fallback-1", "token-2"),
			messageResponse("fallback-2"),
		},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{
				{Model: "fallback-1", Thinking: anthropic.BetaFallbackParamThinkingUnion{
					OfDisabled: &anthropic.BetaThinkingConfigDisabledParam{},
				}},
				{Model: "fallback-2", MaxTokens: anthropic.Int(2048)},
			},
		),
	)

	_, err := client.Beta.Messages.New(context.Background(), fallbackTestParams)
	require.NoError(t, err)
	require.Len(t, transport.bodies, 3)
	assert.Equal(t, "token-1", transport.bodies[1]["fallback_credit_token"], "the token rides with every override")
	assert.Equal(t, "token-2", transport.bodies[2]["fallback_credit_token"])
}

func TestRefusalFallbackMiddlewareErrorsWhenServerSideFallbacksAreRequested(t *testing.T) {
	// Only one chain can adjudicate refusals; erroring beats silently
	// picking one.
	client, transport := fallbackTestClient(t,
		[]string{refusalResponse("primary-model", "credit-token")},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	params := fallbackTestParams
	params.Fallbacks = []anthropic.BetaFallbackParam{{Model: "server-side-fallback"}}
	_, err := client.Beta.Messages.New(context.Background(), params)
	require.ErrorContains(t, err, "Sending the `fallbacks:` request param is not supported when using the `BetaRefusalFallbackMiddleware` middleware. You should either remove the middleware and send `fallbacks:` with the `server-side-fallback-2026-06-01` beta header to let the API handle refusal fallbacks, or omit the `fallbacks:` param if you'd like the `BetaRefusalFallbackMiddleware` middleware to handle fallbacks on the client side.")
	assert.Empty(t, transport.bodies, "no request reaches the server")
}
func TestRefusalFallbackMiddlewareRetriesARejectedTokenOnceWithoutIt(t *testing.T) {
	// The credit contract's documented recovery for a redemption-failure 400
	// is to retry without the token.
	rejection := `{"type": "error", "error": {"type": "invalid_request_error", "message": "fallback_credit_token: does not match"}}`
	transport := &scriptedTransport{t: t}
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{Transport: transport}),
		option.WithMaxRetries(0),
		option.WithMiddleware(betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		)),
	)
	transport.statuses = []int{http.StatusOK, http.StatusBadRequest, http.StatusOK}
	transport.responses = []string{
		refusalResponse("primary-model", "credit-token"),
		rejection,
		messageResponse("fallback-model"),
	}

	message, err := client.Beta.Messages.New(context.Background(), fallbackTestParams)
	require.NoError(t, err)
	assert.Equal(t, anthropic.Model("fallback-model"), message.Model)
	require.Len(t, transport.bodies, 3)
	assert.Equal(t, "credit-token", transport.bodies[1]["fallback_credit_token"])
	_, hasToken := transport.bodies[2]["fallback_credit_token"]
	assert.False(t, hasToken, "recovery retry must omit the token")
}

func TestRefusalFallbackMiddlewareSurfacesAnUnrelatedRejectionAfterOneResend(t *testing.T) {
	// The middleware cannot tell a token rejection from any other 400, so an
	// unrelated one costs one duplicate tokenless request.
	rejection := `{"type": "error", "error": {"type": "invalid_request_error", "message": "max_tokens: too large"}}`
	transport := &scriptedTransport{
		t:         t,
		responses: []string{refusalResponse("primary-model", "credit-token"), rejection, rejection},
		statuses:  []int{http.StatusOK, http.StatusBadRequest, http.StatusBadRequest},
	}
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{Transport: transport}),
		option.WithMaxRetries(0),
		option.WithMiddleware(betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		)),
	)

	_, err := client.Beta.Messages.New(context.Background(), fallbackTestParams)
	var apierr *anthropic.Error
	require.ErrorAs(t, err, &apierr)
	assert.Equal(t, http.StatusBadRequest, apierr.StatusCode)
	require.Len(t, transport.bodies, 3)
	_, hasToken := transport.bodies[2]["fallback_credit_token"]
	assert.False(t, hasToken, "the resend must omit the token")
}

func TestRefusalFallbackMiddlewareReentersTheChainWhenTheClientRetries(t *testing.T) {
	overloaded := `{"type": "error", "error": {"type": "overloaded_error", "message": "Overloaded"}}`
	transport := &scriptedTransport{
		t: t,
		responses: []string{
			refusalResponse("primary-model", "token-1"),
			overloaded,
			refusalResponse("primary-model", "token-2"),
			messageResponse("fallback-model"),
		},
		statuses: []int{http.StatusOK, 529, http.StatusOK, http.StatusOK},
		headers:  []http.Header{nil, {"Retry-After": []string{"0"}}},
	}
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{Transport: transport}),
		option.WithMaxRetries(1),
		option.WithMiddleware(betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		)),
	)

	message, err := client.Beta.Messages.New(context.Background(), fallbackTestParams)
	require.NoError(t, err)
	assert.Equal(t, anthropic.Model("fallback-model"), message.Model)
	assert.Equal(t, []string{"primary-model", "fallback-model", "primary-model", "fallback-model"}, transport.models())
	assert.Equal(t, "token-1", transport.bodies[1]["fallback_credit_token"])
	assert.Equal(t, "token-2", transport.bodies[3]["fallback_credit_token"], "the re-entered chain redeems the fresh token")
}

func TestRefusalFallbackMiddlewareReentersAtThePinWhenTheClientRetries(t *testing.T) {
	overloaded := `{"type": "error", "error": {"type": "overloaded_error", "message": "Overloaded"}}`
	transport := &scriptedTransport{
		t: t,
		responses: []string{
			refusalResponse("primary-model", "credit-token"),
			overloaded,
			messageResponse("fallback-model"),
		},
		statuses: []int{http.StatusOK, 529, http.StatusOK},
		headers:  []http.Header{nil, {"Retry-After": []string{"0"}}},
	}
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"),
		option.WithHTTPClient(&http.Client{Transport: transport}),
		option.WithMaxRetries(1),
		option.WithMiddleware(betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		)),
	)

	message, err := client.Beta.Messages.New(context.Background(), fallbackTestParams,
		betafallback.WithBetaFallbackState(&betafallback.BetaFallbackState{}))
	require.NoError(t, err)
	assert.Equal(t, anthropic.Model("fallback-model"), message.Model)
	assert.Equal(t, []string{"primary-model", "fallback-model", "fallback-model"}, transport.models())
	_, hasToken := transport.bodies[2]["fallback_credit_token"]
	assert.False(t, hasToken, "the pinned re-entry has no refusal to redeem")
}

func TestRefusalFallbackMiddlewareSkipsEncodedResponseBodies(t *testing.T) {
	// A still-encoded body (caller-set Accept-Encoding) can't be inspected;
	// pass it through.
	middleware := betafallback.BetaRefusalFallbackMiddleware(
		[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
	)

	calls := 0
	req, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages?beta=true", strings.NewReader(`{"model": "primary-model"}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	res, err := middleware(req, func(*http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Encoding": []string{"gzip"}},
			Body:       io.NopCloser(strings.NewReader("\x1f\x8b not really gzip")),
		}, nil
	})
	require.NoError(t, err)
	assert.Equal(t, 1, calls)
	assert.NotNil(t, res)
}

func TestRefusalFallbackMiddlewareOptsEveryAttemptIntoTheCreditBeta(t *testing.T) {
	// Minting and redeeming both require the fallback-credit beta header;
	// installing the middleware is the opt-in.
	client, transport := fallbackTestClient(t,
		[]string{refusalResponse("primary-model", "credit-token"), messageResponse("fallback-model")},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	_, err := client.Beta.Messages.New(context.Background(), fallbackTestParams)
	require.NoError(t, err)
	require.Len(t, transport.betas, 2)
	for i, betas := range transport.betas {
		assert.Contains(t, betas, string(anthropic.AnthropicBetaFallbackCredit2026_06_01), "request %d", i)
	}
}

func TestRefusalFallbackMiddlewareKeepsACallerSuppliedCreditBeta(t *testing.T) {
	client, transport := fallbackTestClient(t,
		[]string{messageResponse("primary-model")},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	_, err := client.Beta.Messages.New(context.Background(), fallbackTestParams,
		option.WithHeaderAdd("anthropic-beta", string(anthropic.AnthropicBetaFallbackCredit2026_06_01)))
	require.NoError(t, err)
	require.Len(t, transport.betas, 1)
	assert.Equal(t, []string{string(anthropic.AnthropicBetaFallbackCredit2026_06_01)}, transport.betas[0])
}

func TestRefusalFallbackMiddlewareDoesNotAddTheBetaToPassthroughRequests(t *testing.T) {
	middleware := betafallback.BetaRefusalFallbackMiddleware(
		[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
	)

	req, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages?beta=true",
		strings.NewReader(`{"model": "primary-model", "stream": true}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	_, err = middleware(req, func(seen *http.Request) (*http.Response, error) {
		assert.Empty(t, seen.Header.Values("anthropic-beta"))
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("{}"))}, nil
	})
	require.NoError(t, err)
}

func TestRefusalFallbackMiddlewarePassesErrorResponsesThrough(t *testing.T) {
	middleware := betafallback.BetaRefusalFallbackMiddleware(
		[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
	)

	calls := 0
	req, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages?beta=true", strings.NewReader(`{"model": "primary-model", "messages": [], "max_tokens": 16}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	errorBody := `{"type": "error", "error": {"type": "invalid_request_error", "message": "nope"}}`
	res, err := middleware(req, func(*http.Request) (*http.Response, error) {
		calls++
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader(errorBody)),
		}, nil
	})
	require.NoError(t, err)
	assert.Equal(t, 1, calls)
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Equal(t, errorBody, string(body))
}

func TestRefusalFallbackMiddlewarePassesNonJSONContentTypesThroughUnread(t *testing.T) {
	middleware := betafallback.BetaRefusalFallbackMiddleware(
		[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
	)

	req, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages?beta=true",
		strings.NewReader("--boundary--"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")
	calls := 0
	_, err = middleware(req, func(seen *http.Request) (*http.Response, error) {
		calls++
		body, err := io.ReadAll(seen.Body)
		require.NoError(t, err)
		assert.Equal(t, "--boundary--", string(body))
		assert.Empty(t, seen.Header.Values("anthropic-beta"))
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader("{}"))}, nil
	})
	require.NoError(t, err)
	assert.Equal(t, 1, calls)
}

// opaqueReader hides the concrete reader type so http.NewRequest cannot set
// GetBody, modeling requests built from plain io.Readers.
type opaqueReader struct{ io.Reader }

func TestRefusalFallbackMiddlewareWorksWithUnbufferableRequestBodies(t *testing.T) {
	middleware := betafallback.BetaRefusalFallbackMiddleware(
		[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
	)

	req, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages?beta=true",
		opaqueReader{strings.NewReader(`{"model": "primary-model", "max_tokens": 16, "messages": []}`)})
	require.NoError(t, err)
	require.Nil(t, req.GetBody, "test requires an unbufferable body")
	req.Header.Set("Content-Type", "application/json")

	var bodies []map[string]any
	responses := []string{refusalResponse("primary-model", "credit-token"), messageResponse("fallback-model")}
	res, err := middleware(req, func(seen *http.Request) (*http.Response, error) {
		buf, err := io.ReadAll(seen.Body)
		require.NoError(t, err)
		var body map[string]any
		require.NoError(t, json.Unmarshal(buf, &body))
		bodies = append(bodies, body)
		next := responses[0]
		responses = responses[1:]
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(next)),
		}, nil
	})
	require.NoError(t, err)
	require.Len(t, bodies, 2)
	assert.Equal(t, "fallback-model", bodies[1]["model"])
	assert.Equal(t, "credit-token", bodies[1]["fallback_credit_token"])
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), "end_turn")
}

// fakePlatformTransform stands in for the Bedrock/Vertex adapters below user
// middleware: it moves the model from the body into the URL and "signs" the
// final body.
func fakePlatformTransform(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
	buf, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	var body map[string]json.RawMessage
	if err := json.Unmarshal(buf, &body); err != nil {
		return nil, err
	}
	var model string
	_ = json.Unmarshal(body["model"], &model)
	delete(body, "model")
	rewritten, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	r := req.Clone(req.Context())
	r.URL.Path = "/model/" + model + "/invoke"
	r.Body = io.NopCloser(strings.NewReader(string(rewritten)))
	r.ContentLength = int64(len(rewritten))
	r.Header.Set("X-Fake-Signature", fmt.Sprintf("%x", sha256.Sum256(rewritten)))
	return next(r)
}

func TestRefusalFallbackMiddlewareRedirectsThroughAPlatformTransform(t *testing.T) {
	// The fallback middleware runs above the platform transform: its model
	// swap must redirect the attempt through the transform and the signature
	// must cover the rewritten body.
	var paths []string
	transport := &scriptedTransport{t: t, responses: []string{
		refusalResponse("primary-model", "credit-token"), messageResponse("fallback-model"),
	}}
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

	message, err := client.Beta.Messages.New(context.Background(), fallbackTestParams)
	require.NoError(t, err)
	assert.Equal(t, anthropic.Model("fallback-model"), message.Model)
	assert.Equal(t, []string{"/model/primary-model/invoke", "/model/fallback-model/invoke"}, paths)
	require.Len(t, transport.bodies, 2)
	_, hasModel := transport.bodies[1]["model"]
	assert.False(t, hasModel, "the transform moved the model out of the body")
	assert.Equal(t, "credit-token", transport.bodies[1]["fallback_credit_token"])
}

// roundTripFunc adapts a function to http.RoundTripper.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

func TestRefusalFallbackMiddlewareLeavesTokenCountingAloneEvenWhenPinned(t *testing.T) {
	// Token-counting bodies share model+messages but lack max_tokens; a
	// pinned state must not rewrite them.
	client, transport := fallbackTestClient(t,
		[]string{`{"input_tokens": 10}`},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model", MaxTokens: anthropic.Int(2048)}},
		),
	)

	state := &betafallback.BetaFallbackState{}
	state.SetIndex(0)
	_, err := client.Beta.Messages.CountTokens(context.Background(), anthropic.BetaMessageCountTokensParams{
		Model:    "primary-model",
		Messages: []anthropic.BetaMessageParam{anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("hi"))},
	}, betafallback.WithBetaFallbackState(state))
	require.NoError(t, err)
	require.Len(t, transport.bodies, 1)
	assert.Equal(t, "primary-model", transport.bodies[0]["model"], "pin must not rewrite token counting")
	_, hasMax := transport.bodies[0]["max_tokens"]
	assert.False(t, hasMax)
	assert.Empty(t, transport.betas[0])
}

func TestRefusalFallbackMiddlewareArmsTheCreditBetaOnEveryLeg(t *testing.T) {
	client, transport := fallbackTestClient(t,
		[]string{refusalResponse("primary-model", "credit-token"), messageResponse("fallback-model")},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	_, err := client.Beta.Messages.New(context.Background(), fallbackTestParams)
	require.NoError(t, err)
	for i := range transport.betas {
		assert.Contains(t, transport.betas[i], "fallback-credit-2026-06-01", "request %d", i)
	}
}

func TestRefusalFallbackMiddlewareTrimsReplayedFallbackTurns(t *testing.T) {
	// A fallback block in replayed history is rejected by the server; the
	// request goes out with the turn trimmed.
	client, transport := fallbackTestClient(t,
		[]string{messageResponse("primary-model")},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)

	history := []map[string]any{
		{"role": "user", "content": "hi"},
		{"role": "assistant", "content": []map[string]any{
			{"type": "fallback", "from": map[string]any{"model": "primary-model"}, "to": map[string]any{"model": "fallback-model"}},
			{"type": "text", "text": "answer"},
		}},
		{"role": "user", "content": "continue"},
	}
	_, err := client.Beta.Messages.New(context.Background(), fallbackTestParams,
		option.WithJSONSet("messages", history))
	require.NoError(t, err)
	assert.Equal(t, []string{"primary-model"}, transport.models(), "a fallback block does not pin")
	sent := transport.bodies[0]["messages"].([]any)
	for _, block := range sent[1].(map[string]any)["content"].([]any) {
		assert.NotEqual(t, "fallback", block.(map[string]any)["type"], "the fallback block is removed")
	}
}

func TestRefusalFallbackMiddlewareStandsDownOffTheBetaSurface(t *testing.T) {
	// The plain (non-beta) Messages surface gets stock behavior: one wire
	// request, no credit beta armed, the refusal surfaces verbatim.
	middleware := betafallback.BetaRefusalFallbackMiddleware(
		[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
	)
	req, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages",
		strings.NewReader(`{"model": "primary-model", "messages": [], "max_tokens": 16}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	calls := 0
	res, err := middleware(req, func(seen *http.Request) (*http.Response, error) {
		calls++
		assert.Empty(t, seen.Header.Values("anthropic-beta"))
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(refusalResponse("primary-model", "credit-token"))),
		}, nil
	})
	require.NoError(t, err)
	assert.Equal(t, 1, calls)
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Contains(t, string(body), `"refusal"`)
}

func TestRefusalFallbackMiddlewarePrependsFallbackBlockWhenCallerArmedTheBeta(t *testing.T) {
	// A caller already opted into the credit beta sees the same stitched
	// envelope the server's own chain produces: a fallback boundary block
	// prepended to the serving hop's content.
	middleware := betafallback.BetaRefusalFallbackMiddleware(
		[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
	)
	req, err := http.NewRequest(http.MethodPost, "https://api.anthropic.com/v1/messages?beta=true",
		strings.NewReader(`{"model": "primary-model", "messages": [], "max_tokens": 16}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("anthropic-beta", string(anthropic.AnthropicBetaFallbackCredit2026_06_01))

	served := `{
		"id": "msg_2", "type": "message", "role": "assistant", "model": "fallback-model",
		"content": [{"type": "text", "text": "served"}], "stop_reason": "end_turn",
		"stop_sequence": null, "usage": {"input_tokens": 1, "output_tokens": 1}
	}`
	responses := []string{refusalResponse("primary-model", "credit-token"), served}
	res, err := middleware(req, func(*http.Request) (*http.Response, error) {
		next := responses[0]
		responses = responses[1:]
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(next)),
		}, nil
	})
	require.NoError(t, err)
	buf, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	var message struct {
		Content []map[string]any `json:"content"`
	}
	require.NoError(t, json.Unmarshal(buf, &message))
	require.Len(t, message.Content, 2)
	assert.Equal(t, map[string]any{
		"type":    "fallback",
		"from":    map[string]any{"model": "primary-model"},
		"to":      map[string]any{"model": "fallback-model"},
		"trigger": map[string]any{"type": "refusal", "category": nil},
	}, message.Content[0])
	assert.Equal(t, "served", message.Content[1]["text"])
	assert.Equal(t, int64(len(buf)), res.ContentLength)
}

func TestRefusalFallbackMiddlewarePrependsFallbackBlockWhenMiddlewareArmedTheBeta(t *testing.T) {
	// The seam is unconditional: an auto-armed retry stitches the same
	// envelope a caller-armed one does.
	client, transport := fallbackTestClient(t,
		[]string{refusalResponse("primary-model", "credit-token"), `{
			"id": "msg_2", "type": "message", "role": "assistant", "model": "fallback-model",
			"content": [{"type": "text", "text": "served"}], "stop_reason": "end_turn",
			"stop_sequence": null, "usage": {"input_tokens": 1, "output_tokens": 1}
		}`},
		betafallback.BetaRefusalFallbackMiddleware(
			[]anthropic.BetaFallbackParam{{Model: "fallback-model"}},
		),
	)
	message, err := client.Beta.Messages.New(context.Background(), fallbackTestParams)
	require.NoError(t, err)
	require.Len(t, transport.bodies, 2)
	require.Len(t, message.Content, 2)
	assert.Equal(t, "fallback", string(message.Content[0].Type))
	assert.Equal(t, "served", message.Content[1].Text)
}
