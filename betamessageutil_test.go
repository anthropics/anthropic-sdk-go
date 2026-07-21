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
