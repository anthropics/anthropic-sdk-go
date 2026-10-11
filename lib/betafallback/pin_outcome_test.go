package betafallback_test

import (
	"context"
	"encoding/json"
	"errors"
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

func TestFallbackStateOnlyRemembersServingModels(t *testing.T) {
	for _, streaming := range []bool{false, true} {
		for _, initial := range []int{-1, 0} {
			for _, outcome := range []string{"served", "refusal", "http", "transport", "incomplete"} {
				t.Run(fmt.Sprintf("stream=%v/pin=%d/%s", streaming, initial, outcome), func(t *testing.T) {
					state := &betafallback.BetaFallbackState{}
					state.SetIndex(initial)
					chain := []anthropic.BetaFallbackParam{{Model: "fallback-0"}}
					if initial == 0 {
						chain = append(chain, anthropic.BetaFallbackParam{Model: "fallback-1"})
					}
					target := len(chain) - 1
					primary := "primary-model"
					if initial == 0 {
						primary = "fallback-0"
					}
					candidate := string(chain[target].Model)
					var bodies []map[string]any
					transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
						raw, err := io.ReadAll(req.Body)
						require.NoError(t, err)
						var body map[string]any
						require.NoError(t, json.Unmarshal(raw, &body))
						bodies = append(bodies, body)
						count := len(bodies)
						response := messageResponse(primary)
						contentType := "application/json"
						status := 200
						if streaming {
							response = servedStream(primary)
							contentType = "text/event-stream"
						}
						if count == 1 {
							response = refusalResponse(primary, "token")
							if streaming {
								response = refusalStream(primary, tokenNoClaim)
							}
						} else if count == 2 {
							assert.Equal(t, initial, state.Index(), "must not pin before the request succeeds")
							switch outcome {
							case "served":
								response = messageResponse(candidate)
								if streaming {
									response = servedStream(candidate)
								}
							case "refusal":
								response = refusalResponse(candidate, nil)
								if streaming {
									response = refusalStream(candidate, noToken)
								}
							case "http":
								status = 503
								response = `{"type":"error","error":{"type":"overloaded_error","message":"unavailable"}}`
								contentType = "application/json"
							case "transport":
								return nil, errors.New("connection failed")
							case "incomplete":
								response = `{"type":"message","content":{},"stop_reason":"end_turn"}`
								if streaming {
									response = event("message_start", fmt.Sprintf(`{"type":"message_start","message":{"type":"message","id":"msg_incomplete","role":"assistant","model":%q,"content":[],"usage":{"input_tokens":1,"output_tokens":0}}}`, candidate))
								}
							}
						}
						return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{contentType}}, Body: io.NopCloser(strings.NewReader(response)), Request: req}, nil
					})
					client := anthropic.NewClient(option.WithAPIKey("test-key"), option.WithMaxRetries(0), option.WithHTTPClient(&http.Client{Transport: transport}), option.WithMiddleware(betafallback.BetaRefusalFallbackMiddleware(chain)))
					send := func() error {
						if streaming {
							stream := client.Beta.Messages.NewStreaming(context.Background(), fallbackTestParams, betafallback.WithBetaFallbackState(state))
							defer stream.Close()
							for stream.Next() {
							}
							return stream.Err()
						}
						_, err := client.Beta.Messages.New(context.Background(), fallbackTestParams, betafallback.WithBetaFallbackState(state))
						return err
					}
					err := send()
					if outcome == "served" || outcome == "refusal" {
						require.NoError(t, err)
					}
					require.Len(t, bodies, 2)
					expected := initial
					model := primary
					if outcome == "served" {
						expected = target
						model = candidate
					}
					assert.Equal(t, expected, state.Index())
					require.NoError(t, send())
					require.Len(t, bodies, 3)
					assert.Equal(t, model, bodies[2]["model"], "next request must use the last successful route")
				})
			}
		}
	}
}

func TestFallbackStateDoesNotPinAbandonedStreams(t *testing.T) {
	state := &betafallback.BetaFallbackState{}
	transport := &sseTransport{responses: []string{refusalStream("primary-model", tokenNoClaim), servedStream("fallback-model")}}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})
	stream := client.Beta.Messages.NewStreaming(context.Background(), fallbackTestParams, betafallback.WithBetaFallbackState(state))
	found := false
	for stream.Next() {
		current := stream.Current()
		if current.Type == "content_block_start" && current.ContentBlock.Type == "fallback" {
			found = true
			break
		}
	}
	require.True(t, found)
	require.NoError(t, stream.Close())
	assert.Equal(t, -1, state.Index())
}

func TestFallbackStateRetainsImplicitStopCompletion(t *testing.T) {
	state := &betafallback.BetaFallbackState{}
	terminal := strings.TrimSuffix(servedStream("fallback-model"), event("message_stop", `{"type":"message_stop"}`))
	transport := &sseTransport{responses: []string{refusalStream("primary-model", tokenNoClaim), terminal}}
	client := streamingFallbackClient(t, transport, []anthropic.BetaFallbackParam{{Model: "fallback-model"}})
	message, sequence, _ := collectStream(t, client, context.Background(), fallbackTestParams, betafallback.WithBetaFallbackState(state))
	assert.Equal(t, anthropic.BetaStopReasonEndTurn, message.StopReason)
	assert.Equal(t, "message_stop", sequence[len(sequence)-1])
	assert.Equal(t, 0, state.Index())
}
