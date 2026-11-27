package anthropic_test

import (
	"encoding/json"
	"testing"

	anthropic "github.com/anthropics/anthropic-sdk-go"
)

func TestTextEditorCodeExecutionToolResultBlockRoundtrip(t *testing.T) {
	// A non-error text_editor code-execution result block, as the API returns it.
	wire := `{
		"type": "text_editor_code_execution_tool_result",
		"tool_use_id": "srvtoolu_1",
		"content": {
			"type": "text_editor_code_execution_view_result",
			"content": "line1\nline2\n",
			"file_type": "text",
			"num_lines": 2,
			"start_line": 1,
			"total_lines": 2
		}
	}`

	var block anthropic.ContentBlockUnion
	if err := json.Unmarshal([]byte(wire), &block); err != nil {
		t.Fatal(err)
	}

	param := block.ToParam()
	out, err := json.Marshal(param)
	if err != nil {
		t.Fatal(err)
	}

	// Verify content is a JSON object, not a quoted string
	var result map[string]json.RawMessage
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatal(err)
	}

	contentRaw := result["content"]
	if len(contentRaw) == 0 {
		t.Fatal("content field missing")
	}

	// The content should start with '{' (object), not '"' (string)
	if contentRaw[0] == '"' {
		t.Fatalf("content is double-encoded as JSON string: %s", string(contentRaw))
	}

	// Verify it's a valid JSON object
	var contentObj map[string]interface{}
	if err := json.Unmarshal(contentRaw, &contentObj); err != nil {
		t.Fatalf("content is not a valid JSON object: %v", err)
	}

	if contentObj["type"] != "text_editor_code_execution_view_result" {
		t.Fatalf("unexpected content type: %v", contentObj["type"])
	}
}
