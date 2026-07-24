package anthropic_test

import (
	"encoding/json"
	"testing"

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
