package toolrunner_test

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal/testutil"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/toolrunner"
)

// schemaToBytes converts a map schema to JSON bytes for use with NewBetaToolFromBytes.
func schemaToBytes(t *testing.T, schema map[string]any) []byte {
	t.Helper()
	bytes, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("marshal schema: %v", err)
	}
	return bytes
}

// assistantText extracts concatenated assistant text blocks from a BetaMessage.
func assistantText(msg *anthropic.BetaMessage) string {
	var b strings.Builder
	for _, c := range msg.Content {
		if tb, ok := c.AsAny().(anthropic.BetaTextBlock); ok {
			b.WriteString(tb.Text)
		}
	}
	return b.String()
}

// Shared weather tool used by tests
type weatherRequest struct {
	City  string `json:"city"`
	Units string `json:"units,omitempty"`
}

var weatherSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"city":  map[string]any{"type": "string"},
		"units": map[string]any{"type": "string", "enum": []string{"celsius", "fahrenheit"}},
	},
	"required": []string{"city"},
}

func weatherTool(t *testing.T) anthropic.BetaTool {
	t.Helper()
	tool, err := toolrunner.NewBetaToolFromBytes("get_weather", "Get weather", schemaToBytes(t, weatherSchema),
		func(ctx context.Context, req weatherRequest) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			temp := 20
			if req.Units == "fahrenheit" {
				temp = 68
			}
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: fmt.Sprintf("The weather in %s is %d degrees %s.", req.City, temp, req.Units)},
			}, nil
		})
	if err != nil {
		t.Fatalf("create weather tool: %v", err)
	}
	return tool
}

func newClientWithVCR(t *testing.T, cassette string) anthropic.Client {
	t.Helper()
	httpClient, _ := testutil.NewVCRHTTPClient(t, cassette)
	return anthropic.NewClient(option.WithHTTPClient(httpClient))
}

// Test All() end-to-end

func TestToolRunner_All_Basic(t *testing.T) {
	t.Parallel()
	client := newClientWithVCR(t, "tool_runner_basic")
	tool := weatherTool(t)

	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{tool}, anthropic.BetaToolRunnerParams{
		BetaMessageNewParams: anthropic.BetaMessageNewParams{
			Model:     anthropic.ModelClaudeSonnet4_5,
			MaxTokens: 512,
			Messages: []anthropic.BetaMessageParam{
				anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("What's the weather in San Francisco? Use fahrenheit.")),
			},
		},
		MaxIterations: 5,
	})

	ctx := context.Background()
	var last *anthropic.BetaMessage
	for msg, err := range runner.All(ctx) {
		if err != nil {
			t.Fatalf("runner error: %v", err)
		}
		last = msg
	}
	if last == nil {
		t.Fatalf("no final message produced")
	}

	// Extract assistant text content concisely
	got := []byte(assistantText(last) + "\n")
	testutil.CompareGolden(t, "testdata/snapshots/tool_runner_basic.golden", got)
}

func TestToolRunner_RunToCompletion(t *testing.T) {
	t.Parallel()
	client := newClientWithVCR(t, "tool_runner_run_to_completion")
	tool := weatherTool(t)

	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{tool}, anthropic.BetaToolRunnerParams{
		BetaMessageNewParams: anthropic.BetaMessageNewParams{
			Model:     anthropic.ModelClaudeSonnet4_5,
			MaxTokens: 512,
			Messages: []anthropic.BetaMessageParam{
				anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("What's the weather in San Francisco? Use fahrenheit.")),
			},
		},
		MaxIterations: 5,
	})

	ctx := context.Background()
	last, err := runner.RunToCompletion(ctx)
	if err != nil {
		t.Fatalf("RunToCompletion: %v", err)
	}

	// Extract assistant text content concisely
	got := []byte(assistantText(last) + "\n")
	testutil.CompareGolden(t, "testdata/snapshots/tool_runner_run_to_completion.golden", got)
}

// Test NextMessage step-wise, ensuring an intermediate tool_result is appended, then final answer

func TestToolRunner_NextMessage_Step(t *testing.T) {
	t.Parallel()
	client := newClientWithVCR(t, "tool_runner_next_message")
	tool := weatherTool(t)

	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{tool}, anthropic.BetaToolRunnerParams{
		BetaMessageNewParams: anthropic.BetaMessageNewParams{
			Model:     anthropic.ModelClaudeSonnet4_5,
			MaxTokens: 512,
			Messages: []anthropic.BetaMessageParam{
				anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("What's the weather in SF? Use celsius.")),
			},
		},
		MaxIterations: 5,
	})

	ctx := context.Background()
	// Turn 1: expect tool_use
	msg, err := runner.NextMessage(ctx)
	if err != nil {
		t.Fatalf("NextMessage 1: %v", err)
	}
	if msg == nil {
		t.Fatalf("expected first message")
	}

	got := []byte(assistantText(msg) + "\n")
	testutil.CompareGolden(t, "testdata/snapshots/tool_runner_next_message_step_1.golden", got)
	// Turn 2: tool results sent, expect final assistant
	msg, err = runner.NextMessage(ctx)
	if err != nil {
		t.Fatalf("NextMessage 2: %v", err)
	}
	if msg == nil {
		t.Fatalf("expected final message")
	}

	got = []byte(assistantText(msg) + "\n")
	testutil.CompareGolden(t, "testdata/snapshots/tool_runner_next_message_step_2.golden", got)
}

// Test AllStreaming end-to-end collects final text and compares

func TestToolRunner_AllStreaming(t *testing.T) {
	t.Parallel()
	client := newClientWithVCR(t, "tool_runner_streaming_all")
	tool := weatherTool(t)

	runner := client.Beta.Messages.NewToolRunnerStreaming([]anthropic.BetaTool{tool}, anthropic.BetaToolRunnerParams{
		BetaMessageNewParams: anthropic.BetaMessageNewParams{
			Model:     anthropic.ModelClaudeSonnet4_5,
			MaxTokens: 512,
			Messages: []anthropic.BetaMessageParam{
				anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("Weather in SF in fahrenheit?")),
			},
		},
		MaxIterations: 5,
	})

	ctx := context.Background()
	assistantMessages := []string{}
	for iterator, err := range runner.AllStreaming(ctx) {
		if err != nil {
			t.Fatalf("stream err: %v", err)
		}
		stringBuilder := strings.Builder{}
		for ev := range iterator {
			switch evVariant := ev.AsAny().(type) {
			case anthropic.BetaRawContentBlockDeltaEvent:
				switch deltaVariant := evVariant.Delta.AsAny().(type) {
				case anthropic.BetaTextDelta:
					stringBuilder.WriteString(deltaVariant.Text)
				}
			}
		}
		assistantMessages = append(assistantMessages, stringBuilder.String())
	}

	got := []byte(strings.Join(assistantMessages, "\n"))
	testutil.CompareGolden(t, "testdata/snapshots/tool_runner_streaming_all.golden", got)
}

// Test NextStreaming for a single turn; verify event types set is stable

func TestToolRunner_NextStreaming_EventTypes(t *testing.T) {
	t.Parallel()
	client := newClientWithVCR(t, "tool_runner_next_streaming")
	tool := weatherTool(t)

	runner := client.Beta.Messages.NewToolRunnerStreaming([]anthropic.BetaTool{tool}, anthropic.BetaToolRunnerParams{
		BetaMessageNewParams: anthropic.BetaMessageNewParams{
			Model:     anthropic.ModelClaudeSonnet4_5,
			MaxTokens: 512,
			Messages: []anthropic.BetaMessageParam{
				anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("Weather in SF?")),
			},
		},
		MaxIterations: 1,
	})

	ctx := context.Background()
	events := runner.NextStreaming(ctx)

	eventsTypes := []string{}
	for ev := range events {
		eventsTypes = append(eventsTypes, ev.Type)
	}

	sort.Strings(eventsTypes)

	got := []byte(strings.Join(eventsTypes, "\n") + "\n")
	testutil.CompareGolden(t, "testdata/snapshots/tool_runner_next_streaming_types.golden", got)
}

// Test that tool error is surfaced as a tool_result with is_error and the flow completes

func TestToolRunner_ToolCallError_ThenSuccess(t *testing.T) {
	t.Parallel()
	client := newClientWithVCR(t, "tool_runner_tool_call_error")
	called := false
	tool, err := toolrunner.NewBetaToolFromBytes("get_weather", "Get weather", schemaToBytes(t, weatherSchema),
		func(ctx context.Context, req weatherRequest) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			if !called {
				called = true
				return anthropic.BetaToolResultBlockParamContentUnion{}, fmt.Errorf("Unexpected error, try again")
			}
			return anthropic.BetaToolResultBlockParamContentUnion{OfText: &anthropic.BetaTextBlockParam{Text: "Sunny 68°F"}}, nil
		})
	if err != nil {
		t.Fatalf("create tool: %v", err)
	}

	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{tool}, anthropic.BetaToolRunnerParams{
		BetaMessageNewParams: anthropic.BetaMessageNewParams{
			Model:     anthropic.ModelClaudeSonnet4_5,
			MaxTokens: 512,
			Messages: []anthropic.BetaMessageParam{
				anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("Weather in San Francisco?")),
			},
		},
	})

	ctx := context.Background()
	msg, err := runner.RunToCompletion(ctx)
	if err != nil {
		t.Fatalf("RunToCompletion: %v", err)
	}

	messages := runner.Messages()
	if len(messages) == 0 {
		t.Fatalf("expected messages in runner")
	}

	// look through all the messages to find the tool results
	// 1. should be an error
	// 2. should be a success
	toolResultBlocks := []*anthropic.BetaToolResultBlockParam{}
	for _, msg := range messages {
		for _, c := range msg.Content {
			if tr := c.OfToolResult; tr != nil {
				toolResultBlocks = append(toolResultBlocks, tr)
			}
		}
	}

	if len(toolResultBlocks) != 2 {
		t.Fatalf("expected 2 tool result blocks, got %d", len(toolResultBlocks))
	}

	errorToolResultBlock := toolResultBlocks[0]
	if !errorToolResultBlock.IsError.Value {
		t.Fatalf("expected first tool result to be an error")
	}
	errorText := errorToolResultBlock.Content[0].OfText.Text
	if !strings.Contains(errorText, "Unexpected error") {
		t.Fatalf("expected error message in tool result, got: %s", errorText)
	}

	successToolResultBlock := toolResultBlocks[1]
	if successToolResultBlock.IsError.Value {
		t.Fatalf("expected second tool result to be a success")
	}
	successText := successToolResultBlock.Content[0].OfText.Text
	if successText != "Sunny 68°F" {
		t.Fatalf("expected success message in tool result, got: %s", successText)
	}

	// Final assistant golden snapshot and iteration count
	testutil.CompareGolden(t, "testdata/snapshots/tool_runner_tool_call_error_assistant.golden", []byte(assistantText(msg)+"\n"))
	if runner.IterationCount() != 3 {
		t.Fatalf("expected 3 iterations, got %d", runner.IterationCount())
	}
}

// Test custom handling: intercept tool_use, push our own tool_result, and disable tools for next turn

func TestToolRunner_CustomHandlingWithPushMessages(t *testing.T) {
	t.Parallel()
	client := newClientWithVCR(t, "tool_runner_custom_handling")
	tool := weatherTool(t)

	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{tool}, anthropic.BetaToolRunnerParams{
		BetaMessageNewParams: anthropic.BetaMessageNewParams{
			Model:     anthropic.ModelClaudeSonnet4_5,
			MaxTokens: 512,
			Messages: []anthropic.BetaMessageParam{
				anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("What's the weather in SF in celsius?")),
			},
		},
	})

	ctx := context.Background()
	// First assistant message with tool_use
	msg, err := runner.NextMessage(ctx)
	if err != nil || msg == nil {
		t.Fatalf("turn 1: %v %v", msg, err)
	}
	// Find first tool_use id
	var toolUseID string
	for _, c := range msg.Content {
		if tu, ok := c.AsAny().(anthropic.BetaToolUseBlock); ok {
			toolUseID = tu.ID
			break
		}
	}
	if toolUseID == "" {
		t.Fatalf("expected a tool_use block")
	}
	// Build a new runner with our custom tool_result appended to messages to avoid
	// automatic execution for the prior assistant tool_use turn.
	msgs := runner.Messages()
	msgs = append(msgs, anthropic.NewBetaUserMessage(
		anthropic.BetaContentBlockParamUnion{OfToolResult: &anthropic.BetaToolResultBlockParam{ToolUseID: toolUseID, Content: []anthropic.BetaToolResultBlockParamContentUnion{{OfText: &anthropic.BetaTextBlockParam{Text: "Celsius 20°C"}}}}},
	))

	// No tools so the next turn is just the assistant producing final text
	runner2 := client.Beta.Messages.NewToolRunner(nil, anthropic.BetaToolRunnerParams{
		BetaMessageNewParams: anthropic.BetaMessageNewParams{
			Model:     anthropic.ModelClaudeSonnet4_5,
			MaxTokens: 512,
			Messages:  msgs,
		},
	})

	// Next turn should finalize with assistant text
	msg, err = runner2.NextMessage(ctx)
	if err != nil || msg == nil {
		t.Fatalf("turn 2: %v %v", msg, err)
	}
}

// Test max iterations stops further calls

func TestToolRunner_MaxIterations(t *testing.T) {
	t.Parallel()
	client := newClientWithVCR(t, "tool_runner_max_iterations")
	tool := weatherTool(t)

	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{tool}, anthropic.BetaToolRunnerParams{
		BetaMessageNewParams: anthropic.BetaMessageNewParams{
			Model:     anthropic.ModelClaudeSonnet4_5,
			MaxTokens: 512,
			Messages: []anthropic.BetaMessageParam{
				anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("Check weather in SF and NY, step by step")),
			},
		},
		MaxIterations: 2,
	})

	ctx := context.Background()
	for {
		msg, err := runner.NextMessage(ctx)
		if msg == nil {
			if err != nil {
				t.Fatalf("runner error: %v", err)
			}
			break
		}
	}
	if got := runner.IterationCount(); got != 2 {
		t.Fatalf("expected 2 iterations, got %d", got)
	}
}

// Test concurrent tool execution (multiple tools in one message)

func TestToolRunner_ConcurrentToolExecution(t *testing.T) {
	t.Parallel()
	client := newClientWithVCR(t, "tool_runner_concurrent")

	// Track execution with timing to verify concurrency
	var callCount atomic.Int32
	var executionTimes sync.Map
	startTime := time.Now()

	weatherTool, err := toolrunner.NewBetaToolFromBytes("get_weather", "Get weather for a city", schemaToBytes(t, weatherSchema),
		func(ctx context.Context, req weatherRequest) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			executionTimes.Store(req.City, time.Since(startTime))
			callCount.Add(1)
			// Small delay - if sequential this would take 3x longer
			time.Sleep(50 * time.Millisecond)
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: fmt.Sprintf("Weather in %s: Sunny 72°F", req.City)},
			}, nil
		})
	if err != nil {
		t.Fatalf("create weather tool: %v", err)
	}

	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{weatherTool}, anthropic.BetaToolRunnerParams{
		BetaMessageNewParams: anthropic.BetaMessageNewParams{
			Model:     anthropic.ModelClaudeSonnet4_5,
			MaxTokens: 512,
			Messages: []anthropic.BetaMessageParam{
				anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock(
					"What's the weather in San Francisco, New York, and London? Check all three cities at once.",
				)),
			},
		},
		MaxIterations: 5,
	})

	ctx := context.Background()
	_, err = runner.RunToCompletion(ctx)
	if err != nil {
		t.Fatalf("RunToCompletion: %v", err)
	}

	// Verify multiple tools were called
	count := callCount.Load()
	if count < 2 {
		t.Fatalf("expected at least 2 concurrent tool calls, got %d", count)
	}

	// Verify tool results are in the messages
	messages := runner.Messages()
	toolResultCount := 0
	for _, msg := range messages {
		for _, c := range msg.Content {
			if c.OfToolResult != nil {
				toolResultCount++
			}
		}
	}
	if toolResultCount < 2 {
		t.Fatalf("expected at least 2 tool results, got %d", toolResultCount)
	}
}

// Test context cancellation during tool execution

func TestToolRunner_ContextCancellation(t *testing.T) {
	t.Parallel()
	client := newClientWithVCR(t, "tool_runner_context_cancel")

	toolStarted := make(chan struct{})
	toolCompleted := make(chan struct{})

	slowSchema := map[string]any{
		"type":       "object",
		"properties": map[string]any{"input": map[string]any{"type": "string"}},
	}
	slowTool, err := toolrunner.NewBetaToolFromBytes("slow_tool", "A slow tool", schemaToBytes(t, slowSchema),
		func(ctx context.Context, req struct{ Input string }) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			close(toolStarted)
			// Wait for context cancellation or timeout
			select {
			case <-ctx.Done():
				return anthropic.BetaToolResultBlockParamContentUnion{}, ctx.Err()
			case <-time.After(5 * time.Second):
				close(toolCompleted)
				return anthropic.BetaToolResultBlockParamContentUnion{
					OfText: &anthropic.BetaTextBlockParam{Text: "completed"},
				}, nil
			}
		})
	if err != nil {
		t.Fatalf("create slow tool: %v", err)
	}

	runner := client.Beta.Messages.NewToolRunner([]anthropic.BetaTool{slowTool}, anthropic.BetaToolRunnerParams{
		BetaMessageNewParams: anthropic.BetaMessageNewParams{
			Model:     anthropic.ModelClaudeSonnet4_5,
			MaxTokens: 512,
			Messages: []anthropic.BetaMessageParam{
				anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("Call the slow_tool with input 'test'")),
			},
		},
		MaxIterations: 5,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the runner in a goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		runner.RunToCompletion(ctx)
	}()

	// Wait for tool to start, then cancel
	select {
	case <-toolStarted:
		cancel()
	case <-time.After(10 * time.Second):
		t.Fatal("tool never started")
	}

	// Verify runner completes quickly after cancellation (not waiting 5 seconds)
	select {
	case <-done:
		// Good - runner completed
	case <-toolCompleted:
		t.Fatal("tool completed without cancellation being respected")
	case <-time.After(2 * time.Second):
		t.Fatal("runner did not complete promptly after cancellation")
	}
}

// Test malformed JSON input error handling through Execute

func TestToolRunner_MalformedJSONInput(t *testing.T) {
	t.Parallel()

	type StrictInput struct {
		RequiredField string `json:"required_field"`
		NumberField   int    `json:"number_field"`
	}

	strictSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"required_field": map[string]any{"type": "string"},
			"number_field":   map[string]any{"type": "integer"},
		},
		"required": []string{"required_field"},
	}
	tool, err := toolrunner.NewBetaToolFromBytes("strict_tool", "A tool with strict input", schemaToBytes(t, strictSchema),
		func(ctx context.Context, input StrictInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: "success"},
			}, nil
		})
	if err != nil {
		t.Fatalf("create tool: %v", err)
	}

	// Test Execute with valid JSON
	validJSON := json.RawMessage(`{"required_field": "test", "number_field": 42}`)
	result, err := tool.Execute(context.Background(), validJSON)
	if err != nil {
		t.Fatalf("Execute valid JSON failed: %v", err)
	}
	if result.OfText == nil || result.OfText.Text != "success" {
		t.Fatalf("Execute returned unexpected result: %+v", result)
	}

	// Test Execute with malformed JSON (invalid syntax)
	malformedJSON := json.RawMessage(`{"required_field": "test", "number_field": }`)
	_, err = tool.Execute(context.Background(), malformedJSON)
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}

	// Test Execute with type mismatch (string instead of int)
	typeMismatchJSON := json.RawMessage(`{"required_field": "test", "number_field": "not a number"}`)
	_, err = tool.Execute(context.Background(), typeMismatchJSON)
	if err == nil {
		t.Fatal("expected error for type mismatch, got nil")
	}

	// Test Execute with invalid JSON propagates error
	invalidJSON := json.RawMessage(`{invalid json}`)
	_, err = tool.Execute(context.Background(), invalidJSON)
	if err == nil {
		t.Fatal("expected error for invalid JSON in Execute")
	}
}

// TestToolRunner_SchemaValidation verifies that the tool runner validates inputs
// against the JSON Schema before executing the handler. This prevents missing
// required fields, enum violations, and type mismatches from reaching handlers.
func TestToolRunner_SchemaValidation(t *testing.T) {
	t.Parallel()

	type StrictInput struct {
		City  string `json:"city"`
		Units string `json:"units,omitempty"`
	}

	handlerCalled := false
	tool, err := toolrunner.NewBetaToolFromBytes("weather", "Get weather", schemaToBytes(t, weatherSchema),
		func(ctx context.Context, input StrictInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			handlerCalled = true
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: fmt.Sprintf("Weather in %s (%s)", input.City, input.Units)},
			}, nil
		})
	if err != nil {
		t.Fatalf("create tool: %v", err)
	}

	t.Run("valid input passes validation", func(t *testing.T) {
		handlerCalled = false
		input := json.RawMessage(`{"city": "London", "units": "celsius"}`)
		result, err := tool.Execute(context.Background(), input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !handlerCalled {
			t.Fatal("handler was not called for valid input")
		}
		if result.OfText == nil || result.OfText.Text != "Weather in London (celsius)" {
			t.Fatalf("unexpected result: %+v", result)
		}
	})

	t.Run("missing required field rejected", func(t *testing.T) {
		handlerCalled = false
		// "city" is required but missing
		input := json.RawMessage(`{"units": "celsius"}`)
		_, err := tool.Execute(context.Background(), input)
		if err == nil {
			t.Fatal("expected error for missing required field 'city', got nil")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called when schema validation fails")
		}
		if !strings.Contains(err.Error(), "schema validation failed") {
			t.Fatalf("error should mention schema validation, got: %v", err)
		}
	})

	t.Run("enum violation rejected", func(t *testing.T) {
		handlerCalled = false
		// "units" must be "celsius" or "fahrenheit"
		input := json.RawMessage(`{"city": "London", "units": "kelvin"}`)
		_, err := tool.Execute(context.Background(), input)
		if err == nil {
			t.Fatal("expected error for enum violation on 'units', got nil")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called when schema validation fails")
		}
		if !strings.Contains(err.Error(), "schema validation failed") {
			t.Fatalf("error should mention schema validation, got: %v", err)
		}
	})

	t.Run("wrong type rejected", func(t *testing.T) {
		handlerCalled = false
		// "city" should be string, not number
		input := json.RawMessage(`{"city": 12345}`)
		_, err := tool.Execute(context.Background(), input)
		if err == nil {
			t.Fatal("expected error for wrong type on 'city', got nil")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called when schema validation fails")
		}
	})

	t.Run("empty object rejected when required fields exist", func(t *testing.T) {
		handlerCalled = false
		input := json.RawMessage(`{}`)
		_, err := tool.Execute(context.Background(), input)
		if err == nil {
			t.Fatal("expected error for empty object with required fields, got nil")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called when schema validation fails")
		}
	})

	t.Run("optional field can be omitted", func(t *testing.T) {
		handlerCalled = false
		// "units" is optional, only "city" is required
		input := json.RawMessage(`{"city": "Tokyo"}`)
		_, err := tool.Execute(context.Background(), input)
		if err != nil {
			t.Fatalf("unexpected error for valid input without optional field: %v", err)
		}
		if !handlerCalled {
			t.Fatal("handler was not called for valid input")
		}
	})
}

// TestToolRunner_AdditionalPropertiesRejected verifies that additionalProperties:false
// blocks unknown keys from reaching the handler.
func TestToolRunner_AdditionalPropertiesRejected(t *testing.T) {
	t.Parallel()

	type StrictInput struct {
		Name string `json:"name"`
	}

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{"type": "string"},
		},
		"required":             []string{"name"},
		"additionalProperties": false,
	}

	handlerCalled := false
	tool, err := toolrunner.NewBetaToolFromBytes("strict", "Strict tool", schemaToBytes(t, schema),
		func(ctx context.Context, input StrictInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			handlerCalled = true
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: "ok"},
			}, nil
		})
	if err != nil {
		t.Fatalf("create tool: %v", err)
	}

	t.Run("valid input accepted", func(t *testing.T) {
		handlerCalled = false
		input := json.RawMessage(`{"name": "test"}`)
		_, err := tool.Execute(context.Background(), input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !handlerCalled {
			t.Fatal("handler was not called")
		}
	})

	t.Run("extra property rejected", func(t *testing.T) {
		handlerCalled = false
		input := json.RawMessage(`{"name": "test", "extra": "x"}`)
		_, err := tool.Execute(context.Background(), input)
		if err == nil {
			t.Fatal("expected error for additional property, got nil")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called when additionalProperties is violated")
		}
		if !strings.Contains(err.Error(), "additional property") {
			t.Fatalf("error should mention additional property, got: %v", err)
		}
	})
}

// TestToolRunner_PatternValidation verifies that pattern constraints on string
// properties are enforced at runtime.
func TestToolRunner_PatternValidation(t *testing.T) {
	t.Parallel()

	type URLInput struct {
		URL string `json:"url"`
	}

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"url": map[string]any{
				"type":    "string",
				"pattern": `^https://allowed\.example/`,
			},
		},
		"required": []string{"url"},
	}

	handlerCalled := false
	tool, err := toolrunner.NewBetaToolFromBytes("url_tool", "URL tool", schemaToBytes(t, schema),
		func(ctx context.Context, input URLInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			handlerCalled = true
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: "ok"},
			}, nil
		})
	if err != nil {
		t.Fatalf("create tool: %v", err)
	}

	t.Run("matching pattern accepted", func(t *testing.T) {
		handlerCalled = false
		input := json.RawMessage(`{"url": "https://allowed.example/page"}`)
		_, err := tool.Execute(context.Background(), input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !handlerCalled {
			t.Fatal("handler was not called")
		}
	})

	t.Run("non-matching pattern rejected", func(t *testing.T) {
		handlerCalled = false
		input := json.RawMessage(`{"url": "https://evil.example/attack"}`)
		_, err := tool.Execute(context.Background(), input)
		if err == nil {
			t.Fatal("expected error for pattern violation, got nil")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called when pattern is violated")
		}
		if !strings.Contains(err.Error(), "pattern") {
			t.Fatalf("error should mention pattern, got: %v", err)
		}
	})
}

// TestToolRunner_StringLengthValidation verifies minLength and maxLength enforcement.
func TestToolRunner_StringLengthValidation(t *testing.T) {
	t.Parallel()

	type NameInput struct {
		Name string `json:"name"`
	}

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type":      "string",
				"minLength": 2,
				"maxLength": 10,
			},
		},
		"required": []string{"name"},
	}

	handlerCalled := false
	tool, err := toolrunner.NewBetaToolFromBytes("name_tool", "Name tool", schemaToBytes(t, schema),
		func(ctx context.Context, input NameInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			handlerCalled = true
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: "ok"},
			}, nil
		})
	if err != nil {
		t.Fatalf("create tool: %v", err)
	}

	t.Run("valid length accepted", func(t *testing.T) {
		handlerCalled = false
		_, err := tool.Execute(context.Background(), json.RawMessage(`{"name": "Alice"}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !handlerCalled {
			t.Fatal("handler was not called")
		}
	})

	t.Run("too short rejected", func(t *testing.T) {
		handlerCalled = false
		_, err := tool.Execute(context.Background(), json.RawMessage(`{"name": "A"}`))
		if err == nil {
			t.Fatal("expected error for minLength violation")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called")
		}
	})

	t.Run("too long rejected", func(t *testing.T) {
		handlerCalled = false
		_, err := tool.Execute(context.Background(), json.RawMessage(`{"name": "VeryLongNameHere"}`))
		if err == nil {
			t.Fatal("expected error for maxLength violation")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called")
		}
	})
}

// TestToolRunner_NumericBoundsValidation verifies minimum and maximum enforcement.
func TestToolRunner_NumericBoundsValidation(t *testing.T) {
	t.Parallel()

	type AgeInput struct {
		Age int `json:"age"`
	}

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"age": map[string]any{
				"type":    "integer",
				"minimum": 0,
				"maximum": 150,
			},
		},
		"required": []string{"age"},
	}

	handlerCalled := false
	tool, err := toolrunner.NewBetaToolFromBytes("age_tool", "Age tool", schemaToBytes(t, schema),
		func(ctx context.Context, input AgeInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			handlerCalled = true
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: "ok"},
			}, nil
		})
	if err != nil {
		t.Fatalf("create tool: %v", err)
	}

	t.Run("valid value accepted", func(t *testing.T) {
		handlerCalled = false
		_, err := tool.Execute(context.Background(), json.RawMessage(`{"age": 25}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !handlerCalled {
			t.Fatal("handler was not called")
		}
	})

	t.Run("below minimum rejected", func(t *testing.T) {
		handlerCalled = false
		_, err := tool.Execute(context.Background(), json.RawMessage(`{"age": -1}`))
		if err == nil {
			t.Fatal("expected error for minimum violation")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called")
		}
	})

	t.Run("above maximum rejected", func(t *testing.T) {
		handlerCalled = false
		_, err := tool.Execute(context.Background(), json.RawMessage(`{"age": 200}`))
		if err == nil {
			t.Fatal("expected error for maximum violation")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called")
		}
	})
}
