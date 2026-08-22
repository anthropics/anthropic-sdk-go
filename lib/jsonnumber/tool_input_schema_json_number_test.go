package jsonnumber_test

import (
	"bytes"
	"encoding/json"
	"testing"

	anthropic "github.com/anthropics/anthropic-sdk-go"
)

func TestToolInputSchema_StdlibJSONNumberRemainsJSONNumber(t *testing.T) {
	tool := anthropic.ToolUnionParam{
		OfTool: &anthropic.ToolParam{
			Name: "bounded_search",
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: map[string]any{
					"page_size": map[string]any{
						"type":    "integer",
						"minimum": json.Number("1"),
						"maximum": json.Number("9007199254740993"),
					},
				},
			},
		},
	}

	got, err := json.Marshal(tool)
	if err != nil {
		t.Fatal(err)
	}

	dec := json.NewDecoder(bytes.NewReader(got))
	dec.UseNumber()
	var decoded struct {
		InputSchema struct {
			Properties map[string]map[string]any `json:"properties"`
		} `json:"input_schema"`
	}
	if err := dec.Decode(&decoded); err != nil {
		t.Fatalf("decode %s: %v", got, err)
	}

	pageSize := decoded.InputSchema.Properties["page_size"]
	if pageSize["minimum"] != json.Number("1") {
		t.Errorf("minimum: got %#v, want json.Number(1) in %s", pageSize["minimum"], got)
	}
	if pageSize["maximum"] != json.Number("9007199254740993") {
		t.Errorf("maximum: got %#v, want exact integer in %s", pageSize["maximum"], got)
	}
}
