package anthropic_test

import (
	"encoding/json"
	"testing"

	"github.com/tidwall/gjson"

	"github.com/anthropics/anthropic-sdk-go"
)

func unmarshalBetaContentBlockParam(t *testing.T, jsonData string) anthropic.BetaContentBlockParamUnion {
	t.Helper()
	var block anthropic.BetaContentBlockUnion
	if err := json.Unmarshal([]byte(jsonData), &block); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}
	return block.ToParam()
}

func TestBetaTextCitationToParamKeepsAllFields(t *testing.T) {
	t.Run("page_location keeps cited_text", func(t *testing.T) {
		result := unmarshalBetaContentBlockParam(t, `{"type":"text","text":"x","citations":[{"type":"page_location","cited_text":"quoted","document_index":2,"document_title":"Doc","start_page_number":3,"end_page_number":4}]}`)
		c := result.OfText.Citations[0].OfPageLocation
		if c == nil {
			t.Fatal("Expected OfPageLocation to be non-nil")
		}
		if c.CitedText != "quoted" {
			t.Errorf("Expected cited_text to survive ToParam, got %q", c.CitedText)
		}
	})

	t.Run("search_result_location keeps search_result_index", func(t *testing.T) {
		result := unmarshalBetaContentBlockParam(t, `{"type":"text","text":"x","citations":[{"type":"search_result_location","cited_text":"quoted","title":"T","source":"src-1","search_result_index":5,"start_block_index":1,"end_block_index":2}]}`)
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
		result := unmarshalBetaContentBlockParam(t, `{"type":"text","text":"x","citations":[{"type":"web_search_result_location","cited_text":"quoted","title":"T","url":"https://example.com","encrypted_index":"enc-1"}]}`)
		c := result.OfText.Citations[0].OfWebSearchResultLocation
		if c == nil {
			t.Fatal("Expected OfWebSearchResultLocation to be non-nil")
		}
		if c.URL != "https://example.com" || c.EncryptedIndex != "enc-1" {
			t.Errorf("Expected url/encrypted_index to survive ToParam, got url=%q encrypted_index=%q", c.URL, c.EncryptedIndex)
		}
	})
}

func TestBetaAccumulatePreservesWireJSON(t *testing.T) {
	toolResult := `{"type":"bash_code_execution_tool_result","tool_use_id":"srvtoolu_01","content":{"type":"bash_code_execution_result","stdout":"","stderr":"","return_code":0,"content":[{"type":"bash_code_execution_output","file_id":"file_011ABC"}]}}`
	events := []string{
		`{"type":"message_start","message":{"id":"msg_01","type":"message","role":"assistant","model":"claude-haiku-4-5","content":[],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}}`,
		`{"type":"content_block_start","index":0,"content_block":` + toolResult + `}`,
		`{"type":"content_block_stop","index":0}`,
		`{"type":"content_block_start","index":1,"content_block":{"type":"tool_use","id":"toolu_01","name":"create_pdf","input":{}}}`,
		`{"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"{\"filename\": \"bla"}}`,
		`{"type":"content_block_stop","index":1}`,
		`{"type":"message_delta","delta":{"stop_reason":"max_tokens","stop_sequence":null},"usage":{"output_tokens":4096}}`,
		`{"type":"message_stop"}`,
	}

	var message anthropic.BetaMessage
	for _, eventJSON := range events {
		var event anthropic.BetaRawMessageStreamEventUnion
		if err := json.Unmarshal([]byte(eventJSON), &event); err != nil {
			t.Fatalf("Failed to unmarshal event: %v", err)
		}
		if err := message.Accumulate(event); err != nil {
			t.Errorf("Accumulate(%s) returned an error: %v", event.Type, err)
		}
	}

	if got := message.Content[0].RawJSON(); got != toolResult {
		t.Errorf("Expected content block 0 to keep its wire JSON\n got: %s\nwant: %s", got, toolResult)
	}

	wantMessage := `{"id":"msg_01","type":"message","role":"assistant","model":"claude-haiku-4-5","content":[` + toolResult +
		`,{"type":"tool_use","id":"toolu_01","name":"create_pdf","input":{}}],"stop_reason":"max_tokens","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":4096}}`
	if got := message.RawJSON(); got != wantMessage {
		t.Errorf("Expected the accumulated message JSON to match the wire\n got: %s\nwant: %s", got, wantMessage)
	}

	block, ok := message.Content[0].AsAny().(anthropic.BetaBashCodeExecutionToolResultBlock)
	if !ok {
		t.Fatalf("Expected BetaBashCodeExecutionToolResultBlock, got %T", message.Content[0].AsAny())
	}
	if len(block.Content.Content) != 1 || block.Content.Content[0].FileID != "file_011ABC" {
		t.Errorf("Expected typed access to yield file_011ABC, got %+v", block.Content.Content)
	}
}

// TestBetaTextCitationToParamExhaustive mirrors the non-beta drift guard for
// the beta citation converters.
func TestBetaTextCitationToParamExhaustive(t *testing.T) {
	cases := map[string]string{
		"char_location":              `{"type":"char_location","cited_text":"q","document_index":1,"document_title":"D","start_char_index":2,"end_char_index":3}`,
		"page_location":              `{"type":"page_location","cited_text":"q","document_index":1,"document_title":"D","start_page_number":2,"end_page_number":3}`,
		"content_block_location":     `{"type":"content_block_location","cited_text":"q","document_index":1,"document_title":"D","start_block_index":2,"end_block_index":3}`,
		"search_result_location":     `{"type":"search_result_location","cited_text":"q","title":"T","source":"s","search_result_index":1,"start_block_index":2,"end_block_index":3}`,
		"web_search_result_location": `{"type":"web_search_result_location","cited_text":"q","title":"T","url":"https://e.com","encrypted_index":"e"}`,
	}
	for name, citationJSON := range cases {
		t.Run(name, func(t *testing.T) {
			result := unmarshalBetaContentBlockParam(t, `{"type":"text","text":"x","citations":[`+citationJSON+`]}`)
			assertNoZeroExportedFields(t, result.OfText.Citations[0])
		})
	}
}

// betaMessageStartWithUsage seeds the fields message_delta never re-sends, so a
// test can tell "the delta overwrote it" from "the delta left it alone".
const betaMessageStartWithUsage = `{"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","content":[],"stop_reason":null,"stop_sequence":null,"stop_details":null,"usage":{"input_tokens":10,"output_tokens":1,"cache_creation_input_tokens":2,"cache_read_input_tokens":3,"cache_creation":{"ephemeral_5m_input_tokens":2,"ephemeral_1h_input_tokens":0},"service_tier":"standard","inference_geo":"us"}}}`

// accumulateBetaMessage folds raw stream events into a BetaMessage the way a
// caller of Beta.Messages.NewStreaming does.
func accumulateBetaMessage(t *testing.T, events ...string) anthropic.BetaMessage {
	t.Helper()
	message := anthropic.BetaMessage{}
	for _, raw := range events {
		event := anthropic.BetaRawMessageStreamEventUnion{}
		if err := event.UnmarshalJSON([]byte(raw)); err != nil {
			t.Fatalf("Failed to unmarshal event %s: %v", raw, err)
		}
		if err := message.Accumulate(event); err != nil {
			t.Fatalf("Accumulate(%s): %v", raw, err)
		}
	}
	return message
}

// TestBetaAccumulateMessageDeltaAppliesEveryField pins the accumulated
// BetaMessage to what the non-streaming BetaMessage would hold, including the
// event-level context_management report, which message_start never carries.
func TestBetaAccumulateMessageDeltaAppliesEveryField(t *testing.T) {
	message := accumulateBetaMessage(t,
		betaMessageStartWithUsage,
		`{"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null,"stop_details":null,"container":{"id":"container_1","expires_at":"2025-01-01T00:00:00Z"}},"context_management":{"applied_edits":[{"type":"clear_tool_uses_20250919","cleared_input_tokens":120,"cleared_tool_uses":3}]},"usage":{"input_tokens":12,"output_tokens":99,"cache_creation_input_tokens":4,"cache_read_input_tokens":5,"server_tool_use":{"web_search_requests":2,"web_fetch_requests":1},"output_tokens_details":{"thinking_tokens":7},"iterations":[{"type":"message","model":"model","input_tokens":12,"output_tokens":99,"cache_creation_input_tokens":4,"cache_read_input_tokens":5}]}}`,
		`{"type":"message_stop"}`,
	)

	if len(message.ContextManagement.AppliedEdits) != 1 {
		t.Fatalf("Expected context_management from the delta event, got %+v", message.ContextManagement)
	}
	if edit := message.ContextManagement.AppliedEdits[0]; edit.ClearedToolUses != 3 {
		t.Errorf("Expected the applied edit to survive, got %+v", edit)
	}
	if got := gjson.Get(message.RawJSON(), "context_management.applied_edits.0.cleared_tool_uses").Int(); got != 3 {
		t.Errorf("Expected context_management in the accumulated JSON, got %d", got)
	}
	if message.Container.ID != "container_1" {
		t.Errorf("Expected container from the delta, got %q", message.Container.ID)
	}
	if message.Usage.OutputTokensDetails.ThinkingTokens != 7 {
		t.Errorf("Expected output_tokens_details from the delta, got %d", message.Usage.OutputTokensDetails.ThinkingTokens)
	}
	if message.Usage.ServerToolUse.WebSearchRequests != 2 || message.Usage.ServerToolUse.WebFetchRequests != 1 {
		t.Errorf("Expected server_tool_use from the delta, got %+v", message.Usage.ServerToolUse)
	}
	if len(message.Usage.Iterations) != 1 {
		t.Errorf("Expected iterations from the delta, got %+v", message.Usage.Iterations)
	}
	if message.Usage.InputTokens != 12 || message.Usage.OutputTokens != 99 {
		t.Errorf("Expected input/output tokens 12/99, got %d/%d", message.Usage.InputTokens, message.Usage.OutputTokens)
	}
	if message.Usage.CacheCreationInputTokens != 4 || message.Usage.CacheReadInputTokens != 5 {
		t.Errorf("Expected cache tokens 4/5, got %d/%d", message.Usage.CacheCreationInputTokens, message.Usage.CacheReadInputTokens)
	}
}

// TestBetaAccumulateMessageDeltaKeepsOmittedUsage covers the other half of the
// contract: message_delta omits the counters that do not apply, and it never
// re-sends service_tier, inference_geo or the cache_creation breakdown at all.
func TestBetaAccumulateMessageDeltaKeepsOmittedUsage(t *testing.T) {
	message := accumulateBetaMessage(t,
		betaMessageStartWithUsage,
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
	if message.Usage.ServiceTier != anthropic.BetaUsageServiceTierStandard {
		t.Errorf("Expected service_tier to survive, got %q", message.Usage.ServiceTier)
	}
	if message.Usage.InferenceGeo != "us" {
		t.Errorf("Expected inference_geo to survive, got %q", message.Usage.InferenceGeo)
	}
	if message.ID != "msg_1" || message.Role != "assistant" {
		t.Errorf("Expected id/role to survive, got %q/%q", message.ID, message.Role)
	}
}

// TestBetaAccumulateMessageDeltaStopDetails checks the refusal detail reaches
// the message, and that the last delta wins even when its stop_details is null
// — the key is always sent and null is a meaningful final value, as it is for
// stop_reason and stop_sequence.
func TestBetaAccumulateMessageDeltaStopDetails(t *testing.T) {
	message := accumulateBetaMessage(t,
		betaMessageStartWithUsage,
		`{"type":"message_delta","delta":{"stop_reason":"refusal","stop_sequence":null,"stop_details":{"type":"refusal","category":"bio","explanation":"nope"}},"usage":{"output_tokens":5}}`,
	)
	if message.StopDetails.Category != anthropic.BetaRefusalStopDetailsCategoryBio {
		t.Errorf("Expected stop_details from the delta, got %+v", message.StopDetails)
	}

	message = accumulateBetaMessage(t,
		betaMessageStartWithUsage,
		`{"type":"message_delta","delta":{"stop_reason":"refusal","stop_sequence":null,"stop_details":{"type":"refusal","category":"bio","explanation":"nope"}},"usage":{"output_tokens":5}}`,
		`{"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null,"stop_details":null},"usage":{"output_tokens":9}}`,
	)
	if message.StopDetails.Category != "" {
		t.Errorf("Expected the last delta's null stop_details to win, got %+v", message.StopDetails)
	}
}
