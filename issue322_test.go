package anthropic_test

import (
	"encoding/json"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
)

// TestIssue322_CodeExecutionToolResult_RoundTrip verifies that a full
// Unmarshal → Marshal → Unmarshal → Marshal round-trip on a
// code_execution_tool_result block preserves all nested fields.
func TestIssue322_CodeExecutionToolResult_RoundTrip(t *testing.T) {
	inputJSON := `{
		"type": "code_execution_tool_result",
		"tool_use_id": "toolu_01A02B3C",
		"content": {
			"type": "code_execution_result",
			"stdout": "hello world",
			"stderr": "warning: something",
			"return_code": 0
		}
	}`

	// --- First pass: Unmarshal ---
	var block1 anthropic.ContentBlockParamUnion
	if err := json.Unmarshal([]byte(inputJSON), &block1); err != nil {
		t.Fatalf("first Unmarshal failed: %v", err)
	}

	if block1.OfCodeExecutionToolResult == nil {
		t.Fatal("OfCodeExecutionToolResult is nil after first unmarshal")
	}
	res1 := block1.OfCodeExecutionToolResult.Content.OfRequestCodeExecutionResultBlock
	if res1 == nil {
		t.Fatal("OfRequestCodeExecutionResultBlock is nil — fields were lost during first unmarshal")
	}
	if res1.Stdout != "hello world" {
		t.Errorf("pass1 stdout: want 'hello world', got '%s'", res1.Stdout)
	}
	if res1.Stderr != "warning: something" {
		t.Errorf("pass1 stderr: want 'warning: something', got '%s'", res1.Stderr)
	}
	if res1.ReturnCode != 0 {
		t.Errorf("pass1 return_code: want 0, got %d", res1.ReturnCode)
	}

	// --- First Marshal ---
	marshaled1, err := json.Marshal(block1)
	if err != nil {
		t.Fatalf("first Marshal failed: %v", err)
	}
	assertJSONContentFields(t, "after first marshal", marshaled1, "hello world", "warning: something", 0)

	// --- Second Unmarshal (round-trip) ---
	var block2 anthropic.ContentBlockParamUnion
	if err := json.Unmarshal(marshaled1, &block2); err != nil {
		t.Fatalf("second Unmarshal failed: %v", err)
	}

	if block2.OfCodeExecutionToolResult == nil {
		t.Fatal("OfCodeExecutionToolResult is nil after second unmarshal")
	}
	res2 := block2.OfCodeExecutionToolResult.Content.OfRequestCodeExecutionResultBlock
	if res2 == nil {
		t.Fatal("OfRequestCodeExecutionResultBlock is nil after second unmarshal — round-trip broken")
	}
	if res2.Stdout != "hello world" {
		t.Errorf("pass2 stdout: want 'hello world', got '%s'", res2.Stdout)
	}
	if res2.Stderr != "warning: something" {
		t.Errorf("pass2 stderr: want 'warning: something', got '%s'", res2.Stderr)
	}
	if res2.ReturnCode != 0 {
		t.Errorf("pass2 return_code: want 0, got %d", res2.ReturnCode)
	}

	// --- Second Marshal (final verification) ---
	marshaled2, err := json.Marshal(block2)
	if err != nil {
		t.Fatalf("second Marshal failed: %v", err)
	}
	assertJSONContentFields(t, "after second marshal", marshaled2, "hello world", "warning: something", 0)
}

// TestIssue322_CodeExecutionToolResult_ZeroValues ensures fields with empty
// string / zero int are NOT dropped after a round-trip.
func TestIssue322_CodeExecutionToolResult_ZeroValues(t *testing.T) {
	inputJSON := `{
		"type": "code_execution_tool_result",
		"tool_use_id": "toolu_zero",
		"content": {
			"type": "code_execution_result",
			"stdout": "",
			"stderr": "",
			"return_code": 0
		}
	}`

	var block anthropic.ContentBlockParamUnion
	if err := json.Unmarshal([]byte(inputJSON), &block); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if block.OfCodeExecutionToolResult == nil {
		t.Fatal("OfCodeExecutionToolResult is nil")
	}
	res := block.OfCodeExecutionToolResult.Content.OfRequestCodeExecutionResultBlock
	if res == nil {
		t.Fatal("OfRequestCodeExecutionResultBlock is nil — zero-value content lost during unmarshal")
	}

	marshaled, err := json.Marshal(block)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var outer map[string]any
	if err := json.Unmarshal(marshaled, &outer); err != nil {
		t.Fatal(err)
	}
	content, ok := outer["content"].(map[string]any)
	if !ok {
		t.Fatalf("content field missing or not a map: %s", string(marshaled))
	}
	if content["type"] != "code_execution_result" {
		t.Errorf("content.type: want 'code_execution_result', got '%v'", content["type"])
	}
	// Verify zero-value fields are present in the JSON output
	if _, exists := content["stdout"]; !exists {
		t.Error("content.stdout was dropped from JSON output")
	}
	if _, exists := content["stderr"]; !exists {
		t.Error("content.stderr was dropped from JSON output")
	}
	if _, exists := content["return_code"]; !exists {
		t.Error("content.return_code was dropped from JSON output")
	}
}

// assertJSONContentFields unmarshals the outer JSON and checks the nested
// content object for the expected stdout, stderr, and return_code values.
func assertJSONContentFields(t *testing.T, label string, data []byte, wantStdout, wantStderr string, wantRC float64) {
	t.Helper()
	var outer map[string]any
	if err := json.Unmarshal(data, &outer); err != nil {
		t.Fatalf("%s: json.Unmarshal failed: %v", label, err)
	}
	content, ok := outer["content"].(map[string]any)
	if !ok {
		t.Fatalf("%s: content field missing or not a map: %s", label, string(data))
	}
	if content["type"] != "code_execution_result" {
		t.Errorf("%s: content.type: want 'code_execution_result', got '%v'", label, content["type"])
	}
	if content["stdout"] != wantStdout {
		t.Errorf("%s: content.stdout: want '%s', got '%v'", label, wantStdout, content["stdout"])
	}
	if content["stderr"] != wantStderr {
		t.Errorf("%s: content.stderr: want '%s', got '%v'", label, wantStderr, content["stderr"])
	}
	if content["return_code"] != wantRC {
		t.Errorf("%s: content.return_code: want %v, got %v", label, wantRC, content["return_code"])
	}
}
