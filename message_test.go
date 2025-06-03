// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

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
				OfText: &anthropic.TextBlockParam{Text: "What is a quaternion?", CacheControl: anthropic.NewCacheControlEphemeralParam(), Citations: []anthropic.TextCitationParamUnion{{
					OfCharLocation: &anthropic.CitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
				}}},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaude3_7SonnetLatest,
		Metadata: anthropic.MetadataParam{
			UserID: anthropic.String("13803d75-b4b5-4c3e-b2a2-6f21399b021b"),
		},
		ServiceTier:   anthropic.MessageNewParamsServiceTierAuto,
		StopSequences: []string{"string"},
		System: []anthropic.TextBlockParam{{Text: "x", CacheControl: anthropic.NewCacheControlEphemeralParam(), Citations: []anthropic.TextCitationParamUnion{{
			OfCharLocation: &anthropic.CitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
		}}}},
		Temperature: anthropic.Float(1),
		Thinking: anthropic.ThinkingConfigParamUnion{
			OfEnabled: &anthropic.ThinkingConfigEnabledParam{
				BudgetTokens: 1024,
			},
		},
		ToolChoice: anthropic.ToolChoiceUnionParam{
			OfAuto: &anthropic.ToolChoiceAutoParam{
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
					Required: []string{"location"},
				},
				Name:         "name",
				CacheControl: anthropic.NewCacheControlEphemeralParam(),
				Description:  anthropic.String("Get the current weather in a given location"),
				Type:         anthropic.ToolTypeCustom,
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
				OfText: &anthropic.TextBlockParam{Text: "What is a quaternion?", CacheControl: anthropic.NewCacheControlEphemeralParam(), Citations: []anthropic.TextCitationParamUnion{{
					OfCharLocation: &anthropic.CitationCharLocationParam{CitedText: "cited_text", DocumentIndex: 0, DocumentTitle: anthropic.String("x"), EndCharIndex: 0, StartCharIndex: 0},
				}}},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: anthropic.ModelClaude3_7SonnetLatest,
		System: anthropic.MessageCountTokensParamsSystemUnion{
			OfTextBlockArray: []anthropic.TextBlockParam{{
				Text:         "Today's date is 2024-06-01.",
				CacheControl: anthropic.NewCacheControlEphemeralParam(),
				Citations: []anthropic.TextCitationParamUnion{{
					OfCharLocation: &anthropic.CitationCharLocationParam{
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
			OfEnabled: &anthropic.ThinkingConfigEnabledParam{
				BudgetTokens: 1024,
			},
		},
		ToolChoice: anthropic.ToolChoiceUnionParam{
			OfAuto: &anthropic.ToolChoiceAutoParam{
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
					Required: []string{"location"},
				},
				Name:         "name",
				CacheControl: anthropic.NewCacheControlEphemeralParam(),
				Description:  anthropic.String("Get the current weather in a given location"),
				Type:         anthropic.ToolTypeCustom,
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

func TestAccumulate(t *testing.T) {
	for name, testCase := range map[string]struct {
		expected anthropic.Message
		events   []string
	}{
		"empty message": {
			expected: anthropic.Message{Usage: anthropic.Usage{}},
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
			expected: anthropic.Message{Content: []anthropic.ContentBlockUnion{
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
			expected: anthropic.Message{Content: []anthropic.ContentBlockUnion{
				{Type: "text", Text: "1 + 1 = 2", Citations: []anthropic.TextCitationUnion{{
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
			expected: anthropic.Message{Content: []anthropic.ContentBlockUnion{
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
			expected: anthropic.Message{Content: []anthropic.ContentBlockUnion{
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
			expected: anthropic.Message{Content: []anthropic.ContentBlockUnion{
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
			expected: anthropic.Message{Content: []anthropic.ContentBlockUnion{
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
			expected: anthropic.Message{Content: []anthropic.ContentBlockUnion{
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
			expected: anthropic.Message{Content: []anthropic.ContentBlockUnion{
				{Type: "text", Text: "Let me look up the weather for you."},
				{Type: "thinking", Thinking: "I can look this up using a tool."},
				{Type: "tool_use", ID: "toolu_id", Name: "get_weather", Input: []byte(`{"city": "Los Angeles"}`)},
				{Type: "text", Text: "The weather in Los Angeles is 85 degrees Fahrenheit!"},
			}},
		},
	} {
		t.Run(name, func(t *testing.T) {
			message := anthropic.Message{}
			for _, eventStr := range testCase.events {
				event := anthropic.MessageStreamEventUnion{}
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

func TestMessageNewWithNonStreamingTimeoutLimits(t *testing.T) {
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

	// Set a model with known token limits
	model := anthropic.ModelClaudeOpus4_20250514
	testModelLimit := constant.ModelNonStreamingTokens[string(model)]

	// This test verifies that we can still create a message with tokens below the limit
	safeParams := anthropic.MessageNewParams{
		MaxTokens: int64(testModelLimit - 1000), // Well below the limit
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfText: &anthropic.TextBlockParam{Text: "What is a quaternion?"},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: model,
	}

	_, err := client.Messages.New(context.TODO(), safeParams)
	if err != nil {
		var apierr *anthropic.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("Expected no error for tokens below limit, got: %v", err)
	}

	// This test verifies that we get an error when exceeding the limit
	unsafeParams := anthropic.MessageNewParams{
		MaxTokens: int64(testModelLimit + 1000), // Exceed the limit
		Messages: []anthropic.MessageParam{{
			Content: []anthropic.ContentBlockParamUnion{{
				OfText: &anthropic.TextBlockParam{Text: "What is a quaternion?"},
			}},
			Role: anthropic.MessageParamRoleUser,
		}},
		Model: model,
	}

	_, err = client.Messages.New(context.TODO(), unsafeParams)
	if err == nil {
		t.Fatal("Expected error for tokens above limit, got nil")
	}
}

func TestCalculateNonStreamingTimeout(t *testing.T) {
	// Store original model token limits to restore after test
	originalModelTokenLimits := make(map[string]int)
	for k, v := range constant.ModelNonStreamingTokens {
		originalModelTokenLimits[k] = v
	}
	defer func() {
		// Restore original model token limits
		constant.ModelNonStreamingTokens = originalModelTokenLimits
	}()

	// Set up a test model for consistent testing
	constant.ModelNonStreamingTokens = map[string]int{
		"test-model": 8192,
	}

	tests := []struct {
		name          string
		maxTokens     int
		model         string
		expectTimeout time.Duration
		expectError   bool
		opts          []option.RequestOption
	}{
		{
			name:          "small token count returns expected timeout",
			maxTokens:     1000,
			model:         "any-model",
			expectTimeout: time.Duration(float64(60*60) * float64(1000) / 128000.0 * float64(time.Second)),
			expectError:   false,
		},
		{
			name:          "large token count above default time limit throws error",
			maxTokens:     100000,
			model:         "any-model",
			expectTimeout: 0,
			expectError:   true,
		},
		{
			name:          "token count above model specific limit throws error",
			maxTokens:     9000,
			model:         "test-model",
			expectTimeout: 0,
			expectError:   true,
		},
		{
			name:          "token count below model specific limit is ok",
			maxTokens:     8000,
			model:         "test-model",
			expectTimeout: time.Duration(float64(60*60) * float64(8000) / 128000.0 * float64(time.Second)),
			expectError:   false,
		},
		{
			name:          "user supplies custom timeout so we'll use that instead",
			maxTokens:     1000,
			model:         "test-model",
			expectTimeout: 30 * time.Second, // Custom timeout
			expectError:   false,
			opts: []option.RequestOption{
				option.WithRequestTimeout(30 * time.Second),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			timeout, err := anthropic.CalculateNonStreamingTimeout(tc.maxTokens, anthropic.Model(tc.model), tc.opts)

			if tc.expectError && err == nil {
				t.Error("Expected error but got nil")
			}

			if !tc.expectError && err != nil {
				t.Errorf("Did not expect error but got: %v", err)
			}

			if timeout != tc.expectTimeout {
				t.Errorf("Expected timeout %v but got %v", tc.expectTimeout, timeout)
			}
		})
	}
}

// Test specific model limits
func TestModelLimits(t *testing.T) {
	// Verify the model limits are defined for opus-4 models
	if _, exists := constant.ModelNonStreamingTokens["claude-opus-4-20250514"]; !exists {
		t.Error("Expected model limit for claude-opus-4-20250514 but not found")
	}

	if _, exists := constant.ModelNonStreamingTokens["anthropic.claude-opus-4-20250514-v1:0"]; !exists {
		t.Error("Expected model limit for anthropic.claude-opus-4-20250514-v1:0 but not found")
	}

	if _, exists := constant.ModelNonStreamingTokens["claude-opus-4@20250514"]; !exists {
		t.Error("Expected model limit for claude-opus-4@20250514 but not found")
	}
}
