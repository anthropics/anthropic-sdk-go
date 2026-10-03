package anthropic_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/tidwall/gjson"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

func unmarshalContentBlockParam(t *testing.T, jsonData string) anthropic.ContentBlockParamUnion {
	t.Helper()
	var block anthropic.ContentBlockUnion
	err := json.Unmarshal([]byte(jsonData), &block)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	result := block.ToParam()
	return result
}

func TestContentBlockUnionToParam(t *testing.T) {
	t.Run("TextBlock with text only", func(t *testing.T) {
		result := unmarshalContentBlockParam(t, `{"type":"text","text":"Hello, world!"}`)
		if result.OfText == nil {
			t.Fatal("Expected OfText to be non-nil")
		}
		if result.OfText.Text != "Hello, world!" {
			t.Errorf("Expected text 'Hello, world!', got '%s'", result.OfText.Text)
		}
		if result.OfText.Type != constant.Text("text") {
			t.Errorf("Expected type 'text', got '%s'", result.OfText.Type)
		}
	})

	t.Run("WebSearchToolResultBlock with search results", func(t *testing.T) {
		result := unmarshalContentBlockParam(t, `{"type":"web_search_tool_result","tool_use_id":"test123","content":[{"type":"web_search_result","title":"Test Web Title","url":"https://test.com","encrypted_content":"abc123","page_age":"1 day ago"}]}`)
		var block anthropic.ContentBlockUnion
		if err := json.Unmarshal([]byte(`{"type":"web_search_tool_result","tool_use_id":"test123","content":[{"type":"web_search_result","title":"Test Web Title","url":"https://test.com","encrypted_content":"abc123","page_age":"1 day ago"}]}`), &block); err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}

		if len(block.Content.OfWebSearchResultBlockArray) != 1 {
			t.Errorf("Expected Content.OfWebSearchResultBlockArray to have 1 result, got %d", len(block.Content.OfWebSearchResultBlockArray))
		}
		if len(block.Content.OfWebSearchResultBlockArray) > 0 && block.Content.OfWebSearchResultBlockArray[0].Title != "Test Web Title" {
			t.Errorf("Expected title '', got '%s'", block.Content.OfWebSearchResultBlockArray[0].Title)
		}

		if result.OfWebSearchToolResult == nil {
			t.Fatal("Expected OfWebSearchToolResult to be non-nil")
		}
		if len(result.OfWebSearchToolResult.Content.OfWebSearchToolResultBlockItem) != 1 {
			t.Errorf("Expected 1 search result in param, got %d", len(result.OfWebSearchToolResult.Content.OfWebSearchToolResultBlockItem))
		}
	})

	t.Run("ToolUseBlock keeps unknown caller", func(t *testing.T) {
		result := unmarshalContentBlockParam(t, `{"type":"tool_use","id":"toolu_1","name":"f","input":{},"caller":{"type":"future_caller","x":1}}`)
		b, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("Failed to marshal: %v", err)
		}
		if got := gjson.GetBytes(b, "caller").Raw; got != `{"type":"future_caller","x":1}` {
			t.Errorf("Expected unknown caller to survive ToParam, got %s", got)
		}
	})
}

func TestTextCitationToParamKeepsAllFields(t *testing.T) {
	t.Run("page_location keeps cited_text", func(t *testing.T) {
		result := unmarshalContentBlockParam(t, `{"type":"text","text":"x","citations":[{"type":"page_location","cited_text":"quoted","document_index":2,"document_title":"Doc","start_page_number":3,"end_page_number":4}]}`)
		c := result.OfText.Citations[0].OfPageLocation
		if c == nil {
			t.Fatal("Expected OfPageLocation to be non-nil")
		}
		if c.CitedText != "quoted" {
			t.Errorf("Expected cited_text to survive ToParam, got %q", c.CitedText)
		}
	})

	t.Run("search_result_location keeps index fields and source", func(t *testing.T) {
		result := unmarshalContentBlockParam(t, `{"type":"text","text":"x","citations":[{"type":"search_result_location","cited_text":"quoted","title":"T","source":"src-1","search_result_index":5,"start_block_index":1,"end_block_index":2}]}`)
		c := result.OfText.Citations[0].OfSearchResultLocation
		if c == nil {
			t.Fatal("Expected OfSearchResultLocation to be non-nil")
		}
		if c.Source != "src-1" || c.SearchResultIndex != 5 || c.StartBlockIndex != 1 || c.EndBlockIndex != 2 {
			t.Errorf("Expected source/index fields to survive ToParam, got source=%q search_result_index=%d start=%d end=%d",
				c.Source, c.SearchResultIndex, c.StartBlockIndex, c.EndBlockIndex)
		}
	})

	t.Run("web_search_result_location keeps url and encrypted_index", func(t *testing.T) {
		result := unmarshalContentBlockParam(t, `{"type":"text","text":"x","citations":[{"type":"web_search_result_location","cited_text":"quoted","title":"T","url":"https://example.com","encrypted_index":"enc-1"}]}`)
		c := result.OfText.Citations[0].OfWebSearchResultLocation
		if c == nil {
			t.Fatal("Expected OfWebSearchResultLocation to be non-nil")
		}
		if c.URL != "https://example.com" || c.EncryptedIndex != "enc-1" {
			t.Errorf("Expected url/encrypted_index to survive ToParam, got url=%q encrypted_index=%q", c.URL, c.EncryptedIndex)
		}
	})
}

func TestToolUseToParamKeepsToolsetName(t *testing.T) {
	t.Run("toolset member keeps toolset_name", func(t *testing.T) {
		result := unmarshalContentBlockParam(t, `{"type":"tool_use","id":"toolu_1","name":"screenshot","input":{},"toolset_name":"computer","caller":{"type":"direct"}}`)
		if result.OfToolUse == nil {
			t.Fatal("Expected OfToolUse to be non-nil")
		}
		b, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("Failed to marshal param: %v", err)
		}
		if got := gjson.GetBytes(b, "toolset_name").String(); got != "computer" {
			t.Errorf("Expected toolset_name to survive ToParam, got %q in %s", got, b)
		}
	})

	t.Run("non-member omits toolset_name", func(t *testing.T) {
		result := unmarshalContentBlockParam(t, `{"type":"tool_use","id":"toolu_1","name":"get_weather","input":{}}`)
		b, err := json.Marshal(result)
		if err != nil {
			t.Fatalf("Failed to marshal param: %v", err)
		}
		if gjson.GetBytes(b, "toolset_name").Exists() {
			t.Errorf("Expected toolset_name to be omitted, got %s", b)
		}
	})
}

// TestTextCitationToParamExhaustive guards against converter drift: every
// citation variant round-trips fully-populated JSON, and every exported
// field of the resulting param must be set — a spec-added field that the
// converter forgets to copy fails here without a hand-written assertion.
func TestTextCitationToParamExhaustive(t *testing.T) {
	cases := map[string]string{
		"char_location":              `{"type":"char_location","cited_text":"q","document_index":1,"document_title":"D","start_char_index":2,"end_char_index":3}`,
		"page_location":              `{"type":"page_location","cited_text":"q","document_index":1,"document_title":"D","start_page_number":2,"end_page_number":3}`,
		"content_block_location":     `{"type":"content_block_location","cited_text":"q","document_index":1,"document_title":"D","start_block_index":2,"end_block_index":3}`,
		"search_result_location":     `{"type":"search_result_location","cited_text":"q","title":"T","source":"s","search_result_index":1,"start_block_index":2,"end_block_index":3}`,
		"web_search_result_location": `{"type":"web_search_result_location","cited_text":"q","title":"T","url":"https://e.com","encrypted_index":"e"}`,
	}
	for name, citationJSON := range cases {
		t.Run(name, func(t *testing.T) {
			result := unmarshalContentBlockParam(t, `{"type":"text","text":"x","citations":[`+citationJSON+`]}`)
			assertNoZeroExportedFields(t, result.OfText.Citations[0])
		})
	}
}

// assertNoZeroExportedFields fails for any exported zero-valued field in the
// single set variant of a param union, walking one pointer level.
func assertNoZeroExportedFields(t *testing.T, union any) {
	t.Helper()
	uv := reflect.ValueOf(union)
	var variant reflect.Value
	for i := 0; i < uv.NumField(); i++ {
		f := uv.Field(i)
		if f.Kind() == reflect.Pointer && !f.IsNil() {
			variant = f.Elem()
			break
		}
	}
	if !variant.IsValid() {
		t.Fatal("Expected exactly one non-nil variant in the citation param union")
	}
	vt := variant.Type()
	for i := 0; i < vt.NumField(); i++ {
		field := vt.Field(i)
		if !field.IsExported() || field.Anonymous {
			continue
		}
		if variant.Field(i).IsZero() {
			t.Errorf("Expected %s.%s to be populated by toParamUnion, got the zero value", vt.Name(), field.Name)
		}
	}
}

// messageStartWithUsage seeds the fields message_delta never re-sends, so a
// test can tell "the delta overwrote it" from "the delta left it alone".
const messageStartWithUsage = `{"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","content":[],"stop_reason":null,"stop_sequence":null,"stop_details":null,"usage":{"input_tokens":10,"output_tokens":1,"cache_creation_input_tokens":2,"cache_read_input_tokens":3,"cache_creation":{"ephemeral_5m_input_tokens":2,"ephemeral_1h_input_tokens":0},"service_tier":"standard"}}}`

// accumulateMessage folds raw stream events into a Message the way a caller of
// Messages.NewStreaming does.
func accumulateMessage(t *testing.T, events ...string) anthropic.Message {
	t.Helper()
	message := anthropic.Message{}
	for _, raw := range events {
		event := anthropic.MessageStreamEventUnion{}
		if err := event.UnmarshalJSON([]byte(raw)); err != nil {
			t.Fatalf("Failed to unmarshal event %s: %v", raw, err)
		}
		if err := message.Accumulate(event); err != nil {
			t.Fatalf("Accumulate(%s): %v", raw, err)
		}
	}
	return message
}

// TestAccumulateMessageDeltaAppliesEveryField pins the accumulated Message to
// what the non-streaming Message would hold: every value the final
// message_delta carries reaches it.
func TestAccumulateMessageDeltaAppliesEveryField(t *testing.T) {
	message := accumulateMessage(t,
		messageStartWithUsage,
		`{"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null,"stop_details":null,"container":{"id":"container_1","expires_at":"2025-01-01T00:00:00Z"}},"usage":{"input_tokens":12,"output_tokens":99,"cache_creation_input_tokens":4,"cache_read_input_tokens":5,"server_tool_use":{"web_search_requests":2,"web_fetch_requests":1},"output_tokens_details":{"thinking_tokens":7}}}`,
		`{"type":"message_stop"}`,
	)

	if message.Container.ID != "container_1" {
		t.Errorf("Expected container from the delta, got %q", message.Container.ID)
	}
	if got := gjson.Get(message.RawJSON(), "container.id").String(); got != "container_1" {
		t.Errorf("Expected container in the accumulated JSON, got %q", got)
	}
	if message.StopReason != anthropic.StopReasonEndTurn {
		t.Errorf("Expected stop_reason end_turn, got %q", message.StopReason)
	}
	if message.Usage.OutputTokensDetails.ThinkingTokens != 7 {
		t.Errorf("Expected output_tokens_details from the delta, got %d", message.Usage.OutputTokensDetails.ThinkingTokens)
	}
	if message.Usage.ServerToolUse.WebSearchRequests != 2 || message.Usage.ServerToolUse.WebFetchRequests != 1 {
		t.Errorf("Expected server_tool_use from the delta, got %+v", message.Usage.ServerToolUse)
	}
	// Cumulative totals: the delta's value replaces the message_start one.
	if message.Usage.InputTokens != 12 || message.Usage.OutputTokens != 99 {
		t.Errorf("Expected input/output tokens 12/99, got %d/%d", message.Usage.InputTokens, message.Usage.OutputTokens)
	}
	if message.Usage.CacheCreationInputTokens != 4 || message.Usage.CacheReadInputTokens != 5 {
		t.Errorf("Expected cache tokens 4/5, got %d/%d", message.Usage.CacheCreationInputTokens, message.Usage.CacheReadInputTokens)
	}
}

// TestAccumulateMessageDeltaKeepsOmittedUsage covers the other half of the
// contract: message_delta omits the counters that do not apply, and it never
// re-sends service_tier or the cache_creation breakdown at all.
func TestAccumulateMessageDeltaKeepsOmittedUsage(t *testing.T) {
	message := accumulateMessage(t,
		messageStartWithUsage,
		`{"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null,"stop_details":null},"usage":{"output_tokens":25}}`,
		`{"type":"message_stop"}`,
	)

	if message.Usage.OutputTokens != 25 {
		t.Errorf("Expected output_tokens to be overwritten with 25, got %d", message.Usage.OutputTokens)
	}
	if message.Usage.InputTokens != 10 {
		t.Errorf("Expected input_tokens to keep the message_start value 10, got %d", message.Usage.InputTokens)
	}
	if message.Usage.CacheCreationInputTokens != 2 || message.Usage.CacheReadInputTokens != 3 {
		t.Errorf("Expected cache tokens to keep the message_start values 2/3, got %d/%d", message.Usage.CacheCreationInputTokens, message.Usage.CacheReadInputTokens)
	}
	if message.Usage.CacheCreation.Ephemeral5mInputTokens != 2 {
		t.Errorf("Expected cache_creation to survive, got %+v", message.Usage.CacheCreation)
	}
	if message.Usage.ServiceTier != anthropic.UsageServiceTierStandard {
		t.Errorf("Expected service_tier to survive, got %q", message.Usage.ServiceTier)
	}
	if got := gjson.Get(message.RawJSON(), "usage.service_tier").String(); got != "standard" {
		t.Errorf("Expected service_tier in the accumulated JSON, got %q", got)
	}
	if message.ID != "msg_1" || message.Role != "assistant" {
		t.Errorf("Expected id/role to survive, got %q/%q", message.ID, message.Role)
	}
}

// TestAccumulateMessageDeltaStopDetails checks the refusal detail reaches the
// message, and that the last delta wins even when its stop_details is null —
// the key is always sent and null is a meaningful final value, as it is for
// stop_reason and stop_sequence.
func TestAccumulateMessageDeltaStopDetails(t *testing.T) {
	message := accumulateMessage(t,
		messageStartWithUsage,
		`{"type":"message_delta","delta":{"stop_reason":"refusal","stop_sequence":null,"stop_details":{"type":"refusal","category":"bio","explanation":"nope"}},"usage":{"output_tokens":5}}`,
	)
	if message.StopDetails.Category != anthropic.RefusalStopDetailsCategoryBio {
		t.Errorf("Expected stop_details from the delta, got %+v", message.StopDetails)
	}

	message = accumulateMessage(t,
		messageStartWithUsage,
		`{"type":"message_delta","delta":{"stop_reason":"refusal","stop_sequence":null,"stop_details":{"type":"refusal","category":"bio","explanation":"nope"}},"usage":{"output_tokens":5}}`,
		`{"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null,"stop_details":null},"usage":{"output_tokens":9}}`,
	)
	if message.StopDetails.Category != "" {
		t.Errorf("Expected the last delta's null stop_details to win, got %+v", message.StopDetails)
	}
}

func TestContentBlockToParamPreservesUnknownVariants(t *testing.T) {
	cases := []string{
		`{"type":"future_result","payload":{"number":9007199254740993,"empty":null,"flag":false,"items":["a",2]},"content":"opaque"}`,
		`{"type":"future_text","text":"verbatim <tag> & unicode ☃","extra":{}}`,
		`{"type":"future_tool","id":"call_1","name":"unregistered","input":{"x":1}}`,
	}
	for _, raw := range cases {
		t.Run(gjson.Get(raw, "type").String(), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("ToParam panicked on unknown content: %v", r)
				}
			}()
			var block anthropic.ContentBlockUnion
			if err := json.Unmarshal([]byte(raw), &block); err != nil {
				t.Fatal(err)
			}
			encoded, err := json.Marshal(block.ToParam())
			if err != nil {
				t.Fatal(err)
			}
			// encoding/json may escape HTML even for a raw payload; preserve its
			// normal encoding behavior while keeping unknown fields and numbers.
			expected, err := json.Marshal(json.RawMessage(raw))
			if err != nil {
				t.Fatal(err)
			}
			if string(encoded) != string(expected) {
				t.Errorf("unknown content changed: got %s, want %s", encoded, expected)
			}
			if block.RawJSON() != raw {
				t.Fatal("ToParam modified the response block")
			}
			var beta anthropic.BetaContentBlockUnion
			if err := json.Unmarshal([]byte(raw), &beta); err != nil {
				t.Fatal(err)
			}
			betaEncoded, err := json.Marshal(beta.ToParam())
			if err != nil {
				t.Fatal(err)
			}
			if string(betaEncoded) != string(encoded) {
				t.Errorf("GA and beta differed: %s vs %s", encoded, betaEncoded)
			}
		})
	}
}

func TestUnknownContentSurvivesHTTPMessageReplay(t *testing.T) {
	for _, streaming := range []bool{false, true} {
		t.Run(fmt.Sprintf("streaming=%t", streaming), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("message replay panicked: %v", r)
				}
			}()
			const unknown = `{"type":"future_result","payload":{"number":9007199254740993,"empty":null},"token":"opaque"}`
			const final = `{"id":"msg_2","type":"message","role":"assistant","model":"m","content":[{"type":"text","text":"continued"}],"stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/v1/messages" {
					t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
					http.Error(w, "read error", 500)
					return
				}
				if requests.Add(1) == 1 {
					if streaming {
						w.Header().Set("Content-Type", "text/event-stream")
						events := []string{
							messageStartWithUsage,
							`{"type":"content_block_start","index":0,"content_block":{"type":"text","text":"before"}}`,
							`{"type":"content_block_stop","index":0}`,
							`{"type":"content_block_start","index":1,"content_block":` + unknown + `}`,
							`{"type":"content_block_stop","index":1}`,
							`{"type":"content_block_start","index":2,"content_block":{"type":"text","text":"after"}}`,
							`{"type":"content_block_stop","index":2}`,
							`{"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null},"usage":{"output_tokens":3}}`,
							`{"type":"message_stop"}`,
						}
						for _, event := range events {
							fmt.Fprintf(w, "event: %s\ndata: %s\n\n", gjson.Get(event, "type").String(), event)
						}
					} else {
						w.Header().Set("Content-Type", "application/json")
						fmt.Fprint(w, `{"id":"msg_1","type":"message","role":"assistant","model":"m","content":[{"type":"text","text":"before"},`+unknown+`,{"type":"text","text":"after"}],"stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":3}}`)
					}
					return
				}
				replay := gjson.GetBytes(body, "messages.0.content")
				if got := replay.Get("1").Raw; got != unknown {
					t.Errorf("unknown block not preserved in HTTP request: %s", got)
				}
				if len(replay.Array()) != 3 || replay.Get("0.text").String() != "before" || replay.Get("2.text").String() != "after" {
					t.Errorf("content order changed: %s", replay.Raw)
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, final)
			}))
			defer server.Close()
			client := anthropic.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test-key"), option.WithMaxRetries(0))
			params := anthropic.MessageNewParams{Model: "m", MaxTokens: 16, Messages: []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hello"))}}
			var message anthropic.Message
			if streaming {
				stream := client.Messages.NewStreaming(context.Background(), params)
				defer stream.Close()
				for stream.Next() {
					if err := message.Accumulate(stream.Current()); err != nil {
						t.Fatal(err)
					}
				}
				if err := stream.Err(); err != nil {
					t.Fatal(err)
				}
			} else {
				response, err := client.Messages.New(context.Background(), params)
				if err != nil {
					t.Fatal(err)
				}
				message = *response
			}
			params.Messages = []anthropic.MessageParam{message.ToParam(), anthropic.NewUserMessage(anthropic.NewTextBlock("continue"))}
			response, err := client.Messages.New(context.Background(), params)
			if err != nil {
				t.Fatal(err)
			}
			if response.Content[0].Text != "continued" {
				t.Fatalf("unexpected response: %+v", response)
			}
			if requests.Load() != 2 {
				t.Fatalf("expected two HTTP requests, got %d", requests.Load())
			}
		})
	}
}
