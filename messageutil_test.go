package anthropic_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/tidwall/gjson"

	"github.com/anthropics/anthropic-sdk-go"
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
