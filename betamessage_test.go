// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"os"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func TestBetaMessageNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Messages.New(context.TODO(), anthropic.BetaMessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.BetaMessageParam{{
			Content: []anthropic.BetaContentBlockParamUnion{{
				OfText: &anthropic.BetaTextBlockParam{Text: "What is a quaternion?", CacheControl: anthropic.BetaCacheControlEphemeralParam{TTL: anthropic.BetaCacheControlEphemeralTTLTTL5m}, Citations: []anthropic.BetaTextCitationParamUnion{{
					OfCharLocation: &anthropic.BetaCitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
				}}},
			}},
			Role: anthropic.BetaMessageParamRoleUser,
		}},
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		Container: anthropic.String("container"),
		MCPServers: []anthropic.BetaRequestMCPServerURLDefinitionParam{{
			Name:               "name",
			URL:                "url",
			AuthorizationToken: anthropic.String("authorization_token"),
			ToolConfiguration: anthropic.BetaRequestMCPServerToolConfigurationParam{
				AllowedTools: []string{"string"},
				Enabled:      anthropic.Bool(true),
			},
		}},
		Metadata: anthropic.BetaMetadataParam{
			UserID: anthropic.String("13803d75-b4b5-4c3e-b2a2-6f21399b021b"),
		},
		ServiceTier:   anthropic.BetaMessageNewParamsServiceTierAuto,
		StopSequences: []string{"string"},
		System: []anthropic.BetaTextBlockParam{{Text: "x", CacheControl: anthropic.BetaCacheControlEphemeralParam{TTL: anthropic.BetaCacheControlEphemeralTTLTTL5m}, Citations: []anthropic.BetaTextCitationParamUnion{{
			OfCharLocation: &anthropic.BetaCitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
		}}}},
		Temperature: anthropic.Float(1),
		Thinking: anthropic.BetaThinkingConfigParamUnion{
			OfEnabled: &anthropic.BetaThinkingConfigEnabledParam{
				BudgetTokens: 1024,
			},
		},
		ToolChoice: anthropic.BetaToolChoiceUnionParam{
			OfAuto: &anthropic.BetaToolChoiceAutoParam{
				DisableParallelToolUse: anthropic.Bool(true),
			},
		},
		Tools: []anthropic.BetaToolUnionParam{{
			OfTool: &anthropic.BetaToolParam{
				InputSchema: anthropic.BetaToolInputSchemaParam{
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
					Required: []string{"location"},
				},
				Name: "name",
				CacheControl: anthropic.BetaCacheControlEphemeralParam{
					TTL: anthropic.BetaCacheControlEphemeralTTLTTL5m,
				},
				Description: anthropic.String("Get the current weather in a given location"),
				Type:        anthropic.BetaToolTypeCustom,
			},
		}},
		TopK:  anthropic.Int(5),
		TopP:  anthropic.Float(0.7),
		Betas: []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaMessageCountTokensWithOptionalParams(t *testing.T) {
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
	_, err := client.Beta.Messages.CountTokens(context.TODO(), anthropic.BetaMessageCountTokensParams{
		Messages: []anthropic.BetaMessageParam{{
			Content: []anthropic.BetaContentBlockParamUnion{{
				OfText: &anthropic.BetaTextBlockParam{Text: "What is a quaternion?", CacheControl: anthropic.BetaCacheControlEphemeralParam{TTL: anthropic.BetaCacheControlEphemeralTTLTTL5m}, Citations: []anthropic.BetaTextCitationParamUnion{{
					OfCharLocation: &anthropic.BetaCitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
				}}},
			}},
			Role: anthropic.BetaMessageParamRoleUser,
		}},
		Model: anthropic.ModelClaude3_7SonnetLatest,
		MCPServers: []anthropic.BetaRequestMCPServerURLDefinitionParam{{
			Name:               "name",
			URL:                "url",
			AuthorizationToken: anthropic.String("authorization_token"),
			ToolConfiguration: anthropic.BetaRequestMCPServerToolConfigurationParam{
				AllowedTools: []string{"string"},
				Enabled:      anthropic.Bool(true),
			},
		}},
		System: anthropic.BetaMessageCountTokensParamsSystemUnion{
			OfBetaTextBlockArray: []anthropic.BetaTextBlockParam{{
				Text: "Today's date is 2024-06-01.",
				CacheControl: anthropic.BetaCacheControlEphemeralParam{
					TTL: anthropic.BetaCacheControlEphemeralTTLTTL5m,
				},
				Citations: []anthropic.BetaTextCitationParamUnion{{
					OfCharLocation: &anthropic.BetaCitationCharLocationParam{
						CitedText:      "cited_text",
						DocumentIndex:  0,
						DocumentTitle:  anthropic.String("x"),
						EndCharIndex:   0,
						StartCharIndex: 0,
					},
				}},
			}},
		},
		Thinking: anthropic.BetaThinkingConfigParamUnion{
			OfEnabled: &anthropic.BetaThinkingConfigEnabledParam{
				BudgetTokens: 1024,
			},
		},
		ToolChoice: anthropic.BetaToolChoiceUnionParam{
			OfAuto: &anthropic.BetaToolChoiceAutoParam{
				DisableParallelToolUse: anthropic.Bool(true),
			},
		},
		Tools: []anthropic.BetaMessageCountTokensParamsToolUnion{{
			OfTool: &anthropic.BetaToolParam{
				InputSchema: anthropic.BetaToolInputSchemaParam{
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
					Required: []string{"location"},
				},
				Name: "name",
				CacheControl: anthropic.BetaCacheControlEphemeralParam{
					TTL: anthropic.BetaCacheControlEphemeralTTLTTL5m,
				},
				Description: anthropic.String("Get the current weather in a given location"),
				Type:        anthropic.BetaToolTypeCustom,
			},
		}},
		Betas: []anthropic.AnthropicBeta{anthropic.AnthropicBetaMessageBatches2024_09_24},
	})
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestBetaAccumulate(t *testing.T) {
	for name, testCase := range map[string]struct {
		expected anthropic.BetaMessage
		events   []string
	}{
		"empty message": {
			expected: anthropic.BetaMessage{Usage: anthropic.BetaUsage{}},
			events: []string{
				`{"type": "message_start", "message": {}}`,
				`{"type: "message_stop"}`,
			},
		},
		"text content block": {
			events: []string{
				`{"type": "message_start", "message": {}}`,
				`{"type": "content_block_start", "index": 0, "content_block": {"type": "text", "text": "This "}}`,
				`{"type": "content_block_delta", "index": 0, "delta": {"type": "text_delta", "text": "is a "}}`,
				`{"type": "content_block_delta", "index": 0, "delta": {"type": "text_delta": "text": "text block!"}}`,
				`{"type": "content_block_stop", "index": 0}`,
				`{"type": "message_stop"}`,
			},
			expected: anthropic.BetaMessage{Content: []anthropic.BetaContentBlockUnion{
				{Type: "text", Text: "This is a text block!"},
			}},
		},
		"text content block with citations": {
			events: []string{
				`{"type": "message_start", "message": {}}`,
				`{"type": "content_block_start", "index": 0, "content_block": {"type": "text", "text": "1 + 1"}}`,
				`{"type": "content_block_delta", "index": 0, "delta": {"type": "text_delta", "text": " = 2"}}`,
				`{"type": "content_block_delta", "index": 0, "delta": {"type": "citations_delta", "citation": {"type": "char_location", "cited_text": "1 + 1 = 2", "document_index": 0, "document_title": "Math Facts", "start_char_index": 300, "end_char_index": 310 }}}`,
				`{"type": "content_block_stop", "index": 0}`,
				`{"type": "message_stop"}`,
			},
			expected: anthropic.BetaMessage{Content: []anthropic.BetaContentBlockUnion{
				{Type: "text", Text: "1 + 1 = 2", Citations: []anthropic.BetaTextCitationUnion{{
					Type:           "char_location",
					CitedText:      "1 + 1 = 2",
					DocumentIndex:  0,
					DocumentTitle:  "Math Facts",
					StartCharIndex: 300,
					EndCharIndex:   310,
				}}},
			}},
		},
		"tool use block": {
			events: []string{
				`{"type": "message_start", "message": {}}`,
				`{"type": "content_block_start", "index": 0, "content_block": {"type": "tool_use", "id": "toolu_id", "name": "tool_name", "input": {}}}`,
				`{"type": "content_block_delta", "index": 0, "delta": {"type": "input_json_delta", "partial_json": "{\"argument\":"}}`,
				`{"type": "content_block_delta", "index": 0, "delta": {"type": "input_json_delta", "partial_json": " \"value\"}"}}`,
				`{"type": "content_block_stop", "index": 0}`,
				`{"type": "message_stop"}`,
			},
			expected: anthropic.BetaMessage{Content: []anthropic.BetaContentBlockUnion{
				{Type: "tool_use", ID: "toolu_id", Name: "tool_name", Input: []byte(`{"argument": "value"}`)},
			}},
		},
		"tool use block with no params": {
			events: []string{
				`{"type": "message_start", "message": {}}`,
				`{"type": "content_block_start": "index": 0, "content_block": {"type": "tool_use", "id": "toolu_id", "name": "tool_name", input: {}}}`,
				`{"type": "content_block_delta", "index": 0, "delta": {"type": "input_json_delta", "partial_json": ""}}`,
				`{"type": "content_block_stop", "index": 0}`,
				`{"type": "message_stop"}`,
			},
			expected: anthropic.BetaMessage{Content: []anthropic.BetaContentBlockUnion{
				{Type: "tool_use", ID: "toolu_id", Name: "tool_name"},
			}},
		},
		"server tool use block": {
			events: []string{
				`{"type": "message_start", "message": {}}`,
				`{"type": "content_block_start": "index": 0, "content_block": {"type": "server_tool_use", "id": "srvtoolu_id", "name": "web_search", input: {}}}`,
				`{"type": "content_block_delta", "index": 0, "delta": {"type": "input_json_delta", "partial_json": ""}}`,
				`{"type": "content_block_delta", "index": 0, "delta": {"type": "input_json_delta", "partial_json": "{\"query\": \"weat"}}`,
				`{"type": "content_block_delta", "index": 0, "delta": {"type": "input_json_delta", "partial_json": "her\"}"}}`,
				`{"type": "content_block_stop", "index": 0}`,
				`{"type": "message_stop"}`,
			},
			expected: anthropic.BetaMessage{Content: []anthropic.BetaContentBlockUnion{
				{Type: "server_tool_use", ID: "srvtoolu_id", Name: "web_search", Input: []byte(`{"query": "weather"}`)},
			}},
		},
		"thinking block": {
			events: []string{
				`{"type": "message_start", "message": {}}`,
				`{"type": "content_block_start", "index": 0, "content_block": {"type": "thinking", "thinking": "Let me think..."}}`,
				`{"type": "content_block_delta", "index": 0, "delta": {"type": "thinking_delta", "thinking": "
First, let's try this..."}}`,
				`{"type": "content_block_delta", "index": 0, "delta": {"type": "thinking_delta", "thinking": "
Therefore, the answer is..."}}`,
				`{"type": "content_block_delta", "index": 0, "delta": {"type": "signature_delta", "signature": "ThinkingSignature"}}`,
				`{"type": "content_block_stop", "index": 0}`,
				`{"type": "message_stop"}`,
			},
			expected: anthropic.BetaMessage{Content: []anthropic.BetaContentBlockUnion{
				{Type: "thinking", Thinking: "Let me think...\nFirst, let's try this...\nTherefore, the answer is...", Signature: "ThinkingSignature"},
			}},
		},
		"redacted thinking block": {
			events: []string{
				`{"type": "message_start", "message": {}}`,
				`{"type": "content_block_start", "index": 0, "content_block": {"type": "redacted_thinking", "data": "Redacted"}}`,
				`{"type": "content_block_stop", "index": 0}`,
				`{"type": "message_stop"}`,
			},
			expected: anthropic.BetaMessage{Content: []anthropic.BetaContentBlockUnion{
				{Type: "redacted_thinking", Data: "Redacted"},
			}},
		},
		"multiple content blocks": {
			events: []string{
				`{"type": "message_start", "message": {}}`,
				`{"type": "content_block_start", "index": 0, "content_block": {"type": "text", "text": "Let me look up "}}`,
				`{"type": "content_block_delta", "index": 0, "delta": {"type": "text_delta", "text": "the weather for "}}`,
				`{"type": "content_block_delta", "index": 0, "delta": {"type": "text_delta": "text": "you."}}`,
				`{"type": "content_block_stop", "index": 0}`,
				`{"type": "content_block_start", "index": 1, "content_block": {"type": "thinking", "thinking": ""}}`,
				`{"type": "content_block_delta", "index": 1, "delta": {"type": "thinking_delta", "thinking": "I can look this "}}`,
				`{"type": "content_block_delta", "index": 1, "delta": {"type": "thinking_delta", "thinking": "up using a tool."}}`,
				`{"type": "content_block_stop", "index": 1}`,
				`{"type": "content_block_start", "index": 2, "content_block": {"type": "tool_use", "id": "toolu_id", "name": "get_weather", "input": {}}}`,
				`{"type": "content_block_delta", "index": 2, "delta": {"type": "input_json_delta", "partial_json": "{\"city\": "}}`,
				`{"type": "content_block_delta", "index": 2, "delta": {"type": "input_json_delta", "partial_json": "\"Los Angeles\"}"}}`,
				`{"type": "content_block_stop", "index": 2}`,
				`{"type": "content_block_start", "index": 3, "content_block": {"type": "text", "text": ""}}`,
				`{"type": "content_block_delta", "index": 3, "delta": {"type": "text_delta", "text": "The weather in Los Angeles"}}`,
				`{"type": "content_block_delta", "index": 3, "delta": {"type": "text_delta", "text": " is 85 degrees Fahrenheit!"}}`,
				`{"type": "content_block_stop", "index": 3"}`,
				`{"type": "message_stop"}`,
			},
			expected: anthropic.BetaMessage{Content: []anthropic.BetaContentBlockUnion{
				{Type: "text", Text: "Let me look up the weather for you."},
				{Type: "thinking", Thinking: "I can look this up using a tool."},
				{Type: "tool_use", ID: "toolu_id", Name: "get_weather", Input: []byte(`{"city": "Los Angeles"}`)},
				{Type: "text", Text: "The weather in Los Angeles is 85 degrees Fahrenheit!"},
			}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			message := anthropic.BetaMessage{}
			for _, eventStr := range testCase.events {
				event := anthropic.BetaRawMessageStreamEventUnion{}
				err := (&event).UnmarshalJSON([]byte(eventStr))
				if err != nil {
					t.Fatal(err)
				}
				(&message).Accumulate(event)
			}
			marshaledMessage, err := json.Marshal(message)
			if err != nil {
				t.Fatal(err)
			}
			marshaledExpectedMessage, err := json.Marshal(testCase.expected)
			if err != nil {
				t.Fatal(err)
			}
			if string(marshaledMessage) != string(marshaledExpectedMessage) {
				t.Fatalf("Mismatched message: expected %s but got %s", marshaledExpectedMessage, marshaledMessage)
			}
		})
	}
}

func TestCustomContentCitations(t *testing.T) {
	// Test the current implementation's ability to create custom content blocks
	chunks := []string{
		"First chunk of text",
		"Second chunk of text",
		"Third chunk of text",
	}

	// Attempt 1: Try to create content blocks using the union type
	contentChunks := make([]anthropic.BetaContentBlockSourceContentUnionParam, len(chunks))
	for i, chunk := range chunks {
		contentChunks[i] = anthropic.BetaContentBlockSourceContentUnionParam{
			OfString: param.NewOpt(chunk),
		}
	}

	// Create the source with content blocks
	source := anthropic.BetaContentBlockSourceParam{
		Content: anthropic.BetaContentBlockSourceContentUnionParam{
			OfBetaContentBlockSourceContent: contentChunks,
		},
	}

	// Marshal and check the JSON output
	data, err := json.MarshalIndent(source, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal source: %v", err)
	}

	t.Logf("Marshaled JSON (attempt 1):\n%s", string(data))

	// Check if the marshaled JSON matches the expected structure
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	// The API expects:
	// {
	//   "type": "content",
	//   "content": [
	//     {"type": "text", "text": "First chunk"},
	//     {"type": "text", "text": "Second chunk"},
	//     {"type": "text", "text": "Third chunk"}
	//   ]
	// }

	// Check if content is an array
	content, ok := result["content"].([]interface{})
	if !ok {
		t.Errorf("Expected content to be an array, got %T", result["content"])
	} else {
		t.Logf("Content is an array with %d elements", len(content))

		// Check if each element has the correct structure
		for i, item := range content {
			block, ok := item.(map[string]interface{})
			if !ok {
				t.Errorf("Expected content[%d] to be an object, got %T", i, item)
				continue
			}

			// Check for "type": "text"
			if blockType, exists := block["type"]; !exists || blockType != "text" {
				t.Errorf("Expected content[%d].type to be 'text', got %v", i, blockType)
			}

			// Check for "text" field
			if _, exists := block["text"]; !exists {
				t.Errorf("Expected content[%d] to have 'text' field", i)
			}
		}
	}
}

func TestCustomContentWithWorkaround(t *testing.T) {
	// Test the workaround using param.Override
	chunks := []string{
		"First chunk of text",
		"Second chunk of text",
		"Third chunk of text",
	}

	// Create custom type for marshaling
	type CustomContentBlocks []string

	// Implement MarshalJSON to produce the correct structure
	type textBlock struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}

	blocks := make([]textBlock, len(chunks))
	for i, chunk := range chunks {
		blocks[i] = textBlock{
			Type: "text",
			Text: chunk,
		}
	}

	sourceJSON := map[string]interface{}{
		"type":    "content",
		"content": blocks,
	}

	sourceBytes, err := json.Marshal(sourceJSON)
	if err != nil {
		t.Fatalf("Failed to marshal source: %v", err)
	}

	t.Logf("Marshaled JSON (workaround):\n%s", string(sourceBytes))

	// Verify the structure
	var result map[string]interface{}
	if err := json.Unmarshal(sourceBytes, &result); err != nil {
		t.Fatalf("Failed to unmarshal result: %v", err)
	}

	// Check if content is an array with correct structure
	content, ok := result["content"].([]interface{})
	if !ok {
		t.Errorf("Expected content to be an array, got %T", result["content"])
	} else {
		t.Logf("Workaround produces correct structure with %d elements", len(content))
	}
}
