// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

func TestMessageParamMarshalUnmarshal(t *testing.T) {
	original := anthropic.MessageParam{
		Role: anthropic.MessageParamRoleUser,
		Content: []anthropic.ContentBlockParamUnion{
			anthropic.ContentBlockParamOfRequestTextBlock("Hello, world!"),
			anthropic.ContentBlockParamOfRequestTextBlock("This is a second message"),
		},
	}

	// Marshal the original to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal MessageParam: %v", err)
	}

	// Unmarshal back to a new struct
	var unmarshaled anthropic.MessageParam
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal MessageParam: %v", err)
	}

	// Verify that the content was properly preserved
	if len(unmarshaled.Content) != len(original.Content) {
		t.Errorf("Content length mismatch. Original: %d, Unmarshaled: %d",
			len(original.Content), len(unmarshaled.Content))
	}

	// Check text content of each block
	for i, originalBlock := range original.Content {
		if i >= len(unmarshaled.Content) {
			t.Fatalf("Missing content block at index %d", i)
		}

		originalText := originalBlock.GetText()
		unmarshaledText := unmarshaled.Content[i].GetText()

		if originalText == nil || unmarshaledText == nil {
			t.Errorf("Text is nil at index %d. Original: %v, Unmarshaled: %v",
				i, originalText, unmarshaledText)
			continue
		}

		if *originalText != *unmarshaledText {
			t.Errorf("Content mismatch at index %d. Expected: %q, Got: %q",
				i, *originalText, *unmarshaledText)
		}
	}
}

func TestMessageNewWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := anthropic.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("my-anthropic-api-key"),
	)
	_, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfRequestTextBlock: &anthropic.TextBlockParam{Text: "What is a quaternion?", CacheControl: anthropic.CacheControlEphemeralParam{}, Citations: []anthropic.TextCitationParamUnion{{
					OfRequestCharLocationCitation: &anthropic.CitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
				}}},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaude3_7SonnetLatest,
		Metadata: anthropic.MetadataParam{
			UserID: anthropic.String("13803d75-b4b5-4c3e-b2a2-6f21399b021b"),
		},
		StopSequences: []string{"string"},
		System: []anthropic.TextBlockParam{{Text: "x", CacheControl: anthropic.CacheControlEphemeralParam{}, Citations: []anthropic.TextCitationParamUnion{{
			OfRequestCharLocationCitation: &anthropic.CitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
		}}}},
		Temperature: anthropic.Float(1),
		Thinking: anthropic.ThinkingConfigParamUnion{
			OfThinkingConfigEnabled: &anthropic.ThinkingConfigEnabledParam{
				BudgetTokens: 1024,
			},
		},
		ToolChoice: anthropic.ToolChoiceUnionParam{
			OfToolChoiceAuto: &anthropic.ToolChoiceAutoParam{
				DisableParallelToolUse: anthropic.Bool(true),
			},
		},
		Tools: []anthropic.ToolUnionParam{{
			OfTool: &anthropic.ToolParam{
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]interface{}{
						"location": map[string]interface{}{
							"description": "The city and state, e.g. San Francisco, CA",
							"type":        "string",
						},
						"unit": map[string]interface{}{
							"description": "Unit for the output - one of (celsius, fahrenheit)",
							"type":        "string",
						},
					},
				},
				Name:         "name",
				CacheControl: anthropic.CacheControlEphemeralParam{},
				Description:  anthropic.String("Get the current weather in a given location"),
			},
		}},
		TopK: anthropic.Int(5),
		TopP: anthropic.Float(0.7),
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestMessageCountTokensWithOptionalParams(t *testing.T) {
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := anthropic.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("my-anthropic-api-key"),
	)
	_, err := client.Messages.CountTokens(context.TODO(), anthropic.MessageCountTokensParams{
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfRequestTextBlock: &anthropic.TextBlockParam{Text: "What is a quaternion?", CacheControl: anthropic.CacheControlEphemeralParam{}, Citations: []anthropic.TextCitationParamUnion{{
					OfRequestCharLocationCitation: &anthropic.CitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
				}}},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaude3_7SonnetLatest,
		System: anthropic.MessageCountTokensParamsSystemUnion{
			OfMessageCountTokenssSystemArray: []anthropic.TextBlockParam{{
				Text:         "Today's date is 2024-06-01.",
				CacheControl: anthropic.CacheControlEphemeralParam{},
				Citations: []anthropic.TextCitationParamUnion{{
					OfRequestCharLocationCitation: &anthropic.CitationCharLocationParam{
						CitedText:      "cited_text",
						DocumentIndex:  0,
						DocumentTitle:  anthropic.String("x"),
						EndCharIndex:   0,
						StartCharIndex: 0,
					},
				}},
			}},
		},
		Thinking: anthropic.ThinkingConfigParamUnion{
			OfThinkingConfigEnabled: &anthropic.ThinkingConfigEnabledParam{
				BudgetTokens: 1024,
			},
		},
		ToolChoice: anthropic.ToolChoiceUnionParam{
			OfToolChoiceAuto: &anthropic.ToolChoiceAutoParam{
				DisableParallelToolUse: anthropic.Bool(true),
			},
		},
		Tools: []anthropic.MessageCountTokensToolUnionParam{{
			OfTool: &anthropic.ToolParam{
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]interface{}{
						"location": map[string]interface{}{
							"description": "The city and state, e.g. San Francisco, CA",
							"type":        "string",
						},
						"unit": map[string]interface{}{
							"description": "Unit for the output - one of (celsius, fahrenheit)",
							"type":        "string",
						},
					},
				},
				Name:         "name",
				CacheControl: anthropic.CacheControlEphemeralParam{},
				Description:  anthropic.String("Get the current weather in a given location"),
			},
		}},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestContentBlockParamUnionMarshalUnmarshal(t *testing.T) {
	testCases := []struct {
		name     string
		original anthropic.ContentBlockParamUnion
	}{
		{
			name: "TextBlock with CharLocationCitation",
			original: anthropic.ContentBlockParamUnion{
				OfRequestTextBlock: &anthropic.TextBlockParam{
					Text: "Hello, world!",
					Citations: []anthropic.TextCitationParamUnion{
						{
							OfRequestCharLocationCitation: &anthropic.CitationCharLocationParam{
								CitedText:      "Hello",
								DocumentIndex:  0,
								DocumentTitle:  anthropic.String("Document Title"),
								EndCharIndex:   5,
								StartCharIndex: 0,
								Type:           constant.CharLocation("char_location"),
							},
						},
					},
				},
			},
		},
		{
			name: "TextBlock with PageLocationCitation",
			original: anthropic.ContentBlockParamUnion{
				OfRequestTextBlock: &anthropic.TextBlockParam{
					Text: "Page citation example",
					Citations: []anthropic.TextCitationParamUnion{
						{
							OfRequestPageLocationCitation: &anthropic.CitationPageLocationParam{
								CitedText:       "Page citation",
								DocumentIndex:   1,
								DocumentTitle:   anthropic.String("Book Title"),
								EndPageNumber:   42,
								StartPageNumber: 40,
								Type:            constant.PageLocation("page_location"),
							},
						},
					},
				},
			},
		},
		{
			name: "TextBlock with ContentBlockLocationCitation",
			original: anthropic.ContentBlockParamUnion{
				OfRequestTextBlock: &anthropic.TextBlockParam{
					Text: "Content block citation example",
					Citations: []anthropic.TextCitationParamUnion{
						{
							OfRequestContentBlockLocationCitation: &anthropic.CitationContentBlockLocationParam{
								CitedText:       "Content block",
								DocumentIndex:   2,
								DocumentTitle:   anthropic.String("Document Title"),
								EndBlockIndex:   3,
								StartBlockIndex: 1,
								Type:            constant.ContentBlockLocation("content_block_location"),
							},
						},
					},
				},
			},
		},
		{
			name: "TextBlock with multiple citation types",
			original: anthropic.ContentBlockParamUnion{
				OfRequestTextBlock: &anthropic.TextBlockParam{
					Text: "Mixed citations example",
					Citations: []anthropic.TextCitationParamUnion{
						{
							OfRequestCharLocationCitation: &anthropic.CitationCharLocationParam{
								CitedText:      "Char location",
								DocumentIndex:  0,
								DocumentTitle:  anthropic.String("Document A"),
								EndCharIndex:   12,
								StartCharIndex: 0,
								Type:           constant.CharLocation("char_location"),
							},
						},
						{
							OfRequestPageLocationCitation: &anthropic.CitationPageLocationParam{
								CitedText:       "Page location",
								DocumentIndex:   1,
								DocumentTitle:   anthropic.String("Document B"),
								EndPageNumber:   30,
								StartPageNumber: 25,
								Type:            constant.PageLocation("page_location"),
							},
						},
					},
				},
			},
		},
		{
			name: "ImageBlock",
			original: anthropic.ContentBlockParamUnion{
				OfRequestImageBlock: &anthropic.ImageBlockParam{
					Type: constant.Image("image"),
					Source: anthropic.ImageBlockParamSourceUnion{
						OfBase64ImageSource: &anthropic.Base64ImageSourceParam{
							Data:      "base64encodeddata",
							MediaType: anthropic.Base64ImageSourceMediaTypeImageJPEG,
							Type:      constant.Base64("base64"),
						},
					},
				},
			},
		},
		{
			name: "ToolUseBlock",
			original: anthropic.ContentBlockParamOfRequestToolUseBlock(
				"tool-123",
				map[string]string{"param1": "value1", "param2": "value2"},
				"weather_tool",
			),
		},
		{
			name: "ToolResultBlock",
			original: anthropic.ContentBlockParamUnion{
				OfRequestToolResultBlock: &anthropic.ToolResultBlockParam{
					ToolUseID: "tool-123",
					IsError:   anthropic.Bool(true),
					Type:      constant.ToolResult("tool_result"),
					Content: []anthropic.ToolResultBlockParamContentUnion{
						{
							OfRequestTextBlock: &anthropic.TextBlockParam{
								Text: "Tool result content",
								Type: constant.Text("text"),
							},
						},
					},
				},
			},
		},
		{
			name: "DocumentBlock",
			original: anthropic.ContentBlockParamUnion{
				OfRequestDocumentBlock: &anthropic.DocumentBlockParam{
					Type: constant.Document("document"),
					Source: anthropic.DocumentBlockParamSourceUnion{
						OfPlainTextSource: &anthropic.PlainTextSourceParam{
							Data:      "Document content",
							MediaType: constant.TextPlain("text/plain"),
							Type:      constant.Text("text"),
						},
					},
					Context: anthropic.String("Additional context"),
					Title:   anthropic.String("Document Title"),
				},
			},
		},
		{
			name:     "ThinkingBlock",
			original: anthropic.ContentBlockParamOfRequestThinkingBlock("signature123", "Thinking process content"),
		},
		{
			name:     "RedactedThinkingBlock",
			original: anthropic.ContentBlockParamOfRequestRedactedThinkingBlock("redacted content"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal the original to JSON
			jsonData, err := json.Marshal(tc.original)
			if err != nil {
				t.Fatalf("Failed to marshal ContentBlockParamUnion: %v", err)
			}

			// Unmarshal back to a new struct
			var unmarshaled anthropic.ContentBlockParamUnion
			if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal ContentBlockParamUnion: %v", err)
			}

			// Marshal the unmarshaled struct again for comparison
			jsonDataRemashal, err := json.Marshal(unmarshaled)
			if err != nil {
				t.Fatalf("Failed to re-marshal ContentBlockParamUnion: %v", err)
			}

			// Instead of comparing the entire JSON structure, just verify the type field
			// and a few key properties specific to each block type
			var originalMap, remarshaled map[string]interface{}
			if err := json.Unmarshal(jsonData, &originalMap); err != nil {
				t.Fatalf("Failed to unmarshal for comparison: %v", err)
			}
			if err := json.Unmarshal(jsonDataRemashal, &remarshaled); err != nil {
				t.Fatalf("Failed to unmarshal re-marshaled data: %v", err)
			}

			// Check that the type field is preserved
			if originalType, ok := originalMap["type"]; ok {
				if remarshaledType, ok := remarshaled["type"]; ok {
					if originalType != remarshaledType {
						t.Errorf("Type mismatch: %v vs %v", originalType, remarshaledType)
					}
				} else {
					t.Errorf("Type field missing in remarshaled JSON")
				}
			}

			// Check block-specific fields
			switch tc.name {
			case "TextBlock with CharLocationCitation":
				// Citations validation is complex with nested types - already tested in TestTextCitationParamUnionMarshalUnmarshal
				if unmarshaled.GetText() == nil {
					t.Errorf("Text field is nil in unmarshaled data")
				} else if *unmarshaled.GetText() != "Hello, world!" {
					t.Errorf("Text mismatch: %s vs Hello, world!", *unmarshaled.GetText())
				}
			case "ToolUseBlock":
				if unmarshaled.GetID() == nil {
					t.Errorf("ID field is nil in unmarshaled data")
				} else if *unmarshaled.GetID() != "tool-123" {
					t.Errorf("ID mismatch: %s vs tool-123", *unmarshaled.GetID())
				}
				if unmarshaled.GetName() == nil {
					t.Errorf("Name field is nil in unmarshaled data")
				} else if *unmarshaled.GetName() != "weather_tool" {
					t.Errorf("Name mismatch: %s vs weather_tool", *unmarshaled.GetName())
				}
			case "ToolResultBlock":
				if unmarshaled.GetToolUseID() == nil {
					t.Errorf("ToolUseID field is nil in unmarshaled data")
				} else if *unmarshaled.GetToolUseID() != "tool-123" {
					t.Errorf("ToolUseID mismatch: %s vs tool-123", *unmarshaled.GetToolUseID())
				}
			case "DocumentBlock":
				if unmarshaled.GetTitle() == nil {
					t.Errorf("Title field is nil in unmarshaled data")
				} else if *unmarshaled.GetTitle() != "Document Title" {
					t.Errorf("Title mismatch: %s vs Document Title", *unmarshaled.GetTitle())
				}
			case "ThinkingBlock":
				if unmarshaled.GetThinking() == nil {
					t.Errorf("Thinking field is nil in unmarshaled data")
				} else if *unmarshaled.GetThinking() != "Thinking process content" {
					t.Errorf("Thinking mismatch: %s vs Thinking process content", *unmarshaled.GetThinking())
				}
			case "RedactedThinkingBlock":
				if unmarshaled.GetData() == nil {
					t.Errorf("Data field is nil in unmarshaled data")
				} else if *unmarshaled.GetData() != "redacted content" {
					t.Errorf("Data mismatch: %s vs redacted content", *unmarshaled.GetData())
				}
			}

			// Compare the original and re-marshaled values
			if !reflect.DeepEqual(originalMap, remarshaled) {
				t.Fatalf("Unmarshaled-then-marshaled value differs from original.\nOriginal: %+v\nUnmarshaled: %+v", originalMap, remarshaled)
			}
		})
	}
}

func TestTextCitationParamUnionMarshalUnmarshal(t *testing.T) {
	testCases := []struct {
		name     string
		original anthropic.TextCitationParamUnion
	}{
		{
			name: "CharLocationCitation",
			original: anthropic.TextCitationParamUnion{
				OfRequestCharLocationCitation: &anthropic.CitationCharLocationParam{
					CitedText:      "Cited text for char location",
					DocumentIndex:  0,
					DocumentTitle:  anthropic.String("Document A"),
					EndCharIndex:   25,
					StartCharIndex: 0,
					Type:           constant.CharLocation("char_location"),
				},
			},
		},
		{
			name: "PageLocationCitation",
			original: anthropic.TextCitationParamUnion{
				OfRequestPageLocationCitation: &anthropic.CitationPageLocationParam{
					CitedText:       "Cited text for page location",
					DocumentIndex:   1,
					DocumentTitle:   anthropic.String("Document B"),
					EndPageNumber:   45,
					StartPageNumber: 42,
					Type:            constant.PageLocation("page_location"),
				},
			},
		},
		{
			name: "ContentBlockLocationCitation",
			original: anthropic.TextCitationParamUnion{
				OfRequestContentBlockLocationCitation: &anthropic.CitationContentBlockLocationParam{
					CitedText:       "Cited text for content block location",
					DocumentIndex:   2,
					DocumentTitle:   anthropic.String("Document C"),
					EndBlockIndex:   5,
					StartBlockIndex: 3,
					Type:            constant.ContentBlockLocation("content_block_location"),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Marshal the original to JSON
			jsonData, err := json.Marshal(tc.original)
			if err != nil {
				t.Fatalf("Failed to marshal TextCitationParamUnion: %v", err)
			}

			// Unmarshal back to a new struct
			var unmarshaled anthropic.TextCitationParamUnion
			if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal TextCitationParamUnion: %v", err)
			}

			// Marshal the unmarshaled struct again for comparison
			jsonDataRemashal, err := json.Marshal(unmarshaled)
			if err != nil {
				t.Fatalf("Failed to re-marshal TextCitationParamUnion: %v", err)
			}

			// Compare the JSON representations
			var original, remashaled map[string]interface{}
			if err := json.Unmarshal(jsonData, &original); err != nil {
				t.Fatalf("Failed to unmarshal for comparison: %v", err)
			}
			if err := json.Unmarshal(jsonDataRemashal, &remashaled); err != nil {
				t.Fatalf("Failed to unmarshal re-marshaled data: %v", err)
			}

			// Compare the original and re-marshaled values
			if !reflect.DeepEqual(original, remashaled) {
				t.Errorf("Unmarshaled-then-marshaled value differs from original.\nOriginal: %+v\nUnmarshaled: %+v", original, remashaled)
			}

			// Verify specific fields based on the citation type
			switch tc.name {
			case "CharLocationCitation":
				if tc.original.GetType() == nil || unmarshaled.GetType() == nil {
					t.Errorf("Citation type is nil")
				} else if *tc.original.GetType() != *unmarshaled.GetType() {
					t.Errorf("Citation type mismatch: %s vs %s", *tc.original.GetType(), *unmarshaled.GetType())
				}
				if tc.original.GetCitedText() == nil || unmarshaled.GetCitedText() == nil {
					t.Errorf("CitedText is nil")
				} else if *tc.original.GetCitedText() != *unmarshaled.GetCitedText() {
					t.Errorf("CitedText mismatch: %s vs %s", *tc.original.GetCitedText(), *unmarshaled.GetCitedText())
				}
				if tc.original.GetStartCharIndex() == nil || unmarshaled.GetStartCharIndex() == nil {
					t.Errorf("StartCharIndex is nil")
				} else if *tc.original.GetStartCharIndex() != *unmarshaled.GetStartCharIndex() {
					t.Errorf("StartCharIndex mismatch: %d vs %d", *tc.original.GetStartCharIndex(), *unmarshaled.GetStartCharIndex())
				}
				if tc.original.GetEndCharIndex() == nil || unmarshaled.GetEndCharIndex() == nil {
					t.Errorf("EndCharIndex is nil")
				} else if *tc.original.GetEndCharIndex() != *unmarshaled.GetEndCharIndex() {
					t.Errorf("EndCharIndex mismatch: %d vs %d", *tc.original.GetEndCharIndex(), *unmarshaled.GetEndCharIndex())
				}
			case "PageLocationCitation":
				if tc.original.GetStartPageNumber() == nil || unmarshaled.GetStartPageNumber() == nil {
					t.Errorf("StartPageNumber is nil")
				} else if *tc.original.GetStartPageNumber() != *unmarshaled.GetStartPageNumber() {
					t.Errorf("StartPageNumber mismatch: %d vs %d", *tc.original.GetStartPageNumber(), *unmarshaled.GetStartPageNumber())
				}
				if tc.original.GetEndPageNumber() == nil || unmarshaled.GetEndPageNumber() == nil {
					t.Errorf("EndPageNumber is nil")
				} else if *tc.original.GetEndPageNumber() != *unmarshaled.GetEndPageNumber() {
					t.Errorf("EndPageNumber mismatch: %d vs %d", *tc.original.GetEndPageNumber(), *unmarshaled.GetEndPageNumber())
				}
			case "ContentBlockLocationCitation":
				if tc.original.GetStartBlockIndex() == nil || unmarshaled.GetStartBlockIndex() == nil {
					t.Errorf("StartBlockIndex is nil")
				} else if *tc.original.GetStartBlockIndex() != *unmarshaled.GetStartBlockIndex() {
					t.Errorf("StartBlockIndex mismatch: %d vs %d", *tc.original.GetStartBlockIndex(), *unmarshaled.GetStartBlockIndex())
				}
				if tc.original.GetEndBlockIndex() == nil || unmarshaled.GetEndBlockIndex() == nil {
					t.Errorf("EndBlockIndex is nil")
				} else if *tc.original.GetEndBlockIndex() != *unmarshaled.GetEndBlockIndex() {
					t.Errorf("EndBlockIndex mismatch: %d vs %d", *tc.original.GetEndBlockIndex(), *unmarshaled.GetEndBlockIndex())
				}
			}
		})
	}
}
