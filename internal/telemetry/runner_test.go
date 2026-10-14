package telemetry

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

// The cases of testdata/tool_result_size_bound.json are the same for every SDK.
// Their sizes are the bytes of JSON without insignificant space.
func TestToolResultContentSizeBoundCases(t *testing.T) {
	data, err := os.ReadFile("testdata/tool_result_size_bound.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Name       string `json:"name"`
		Bound      int    `json:"bound"`
		SizeBefore int    `json:"size_before"`
		Input      struct {
			Arguments json.RawMessage `json:"tool_call_arguments"`
			Content   json.RawMessage `json:"tool_result_content"`
		} `json:"input"`
		// Arguments that are null are an absent gen_ai.tool.call.arguments.
		Expected struct {
			Arguments json.RawMessage `json:"tool_call_arguments"`
			Content   json.RawMessage `json:"tool_result_content"`
		} `json:"expected"`
		Truncated bool `json:"expected_truncated"`
		SizeAfter int  `json:"size_after"`
	}
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("no cases")
	}
	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			arguments := json.RawMessage(compactJSON(t, test.Input.Arguments))
			result := json.RawMessage(test.Input.Content)
			wholeAttrs := toolUseContentAttributes(contentLimit, arguments, result, false)
			whole := attributeValues(wholeAttrs)
			if size := len(whole[GenAIToolCallArgumentsKey]) + len(whole[AnthropicToolResultContentKey]); size != test.SizeBefore {
				t.Errorf("the whole content has %d bytes, want %d", size, test.SizeBefore)
			}
			_, _, wholeTruncated := contentAttributes(t, wholeAttrs)
			wantTruncated(t, wholeTruncated, false)
			if whole[GenAIToolCallArgumentsKey] != string(arguments) {
				t.Errorf("gen_ai.tool.call.arguments = %s under the default limit, want %s", whole[GenAIToolCallArgumentsKey], arguments)
			}

			attrs := toolUseContentAttributes(test.Bound, arguments, result, false)
			values := attributeValues(attrs)
			got := values[AnthropicToolResultContentKey]
			if size := len(values[GenAIToolCallArgumentsKey]) + len(got); size != test.SizeAfter {
				t.Errorf("the content has %d bytes, want %d", size, test.SizeAfter)
			}
			gotArguments, ok := values[GenAIToolCallArgumentsKey]
			if want := compactJSON(t, test.Expected.Arguments); ok != present(test.Expected.Arguments) || gotArguments != want {
				t.Errorf("gen_ai.tool.call.arguments = %q (present %t), want %q", gotArguments, ok, want)
			}
			_, _, truncated := contentAttributes(t, attrs)
			wantTruncated(t, truncated, test.Truncated)
			// The fields of a block are in the order of their names, not that of the case.
			var gotBlocks, wantBlocks any
			if err := json.Unmarshal([]byte(got), &gotBlocks); err != nil {
				t.Fatalf("anthropic.tool_result.content is not JSON: %v", err)
			}
			if err := json.Unmarshal(test.Expected.Content, &wantBlocks); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(gotBlocks, wantBlocks) {
				t.Errorf("anthropic.tool_result.content = %s, want %s", got, test.Expected.Content)
			}
		})
	}
}

func TestToolUseContentAttributesLeaveOutArgumentsThatAloneAreOverTheLimit(t *testing.T) {
	arguments := json.RawMessage(`{"city":"Paris, the capital of France"}`)
	const ok = `[{"type":"text","text":"ok"}]`
	for name, test := range map[string]struct {
		limit    int
		canceled bool
		// arguments and result are the values of the attributes, or "" for an absent one.
		arguments, result string
		truncated         bool
	}{
		"arguments that fill the limit stay, and leave the result no room": {limit: len(arguments), arguments: string(arguments), result: `[]`, truncated: true},
		"one byte over the limit, they go, and the result has all of it":   {limit: len(arguments) - 1, result: ok, truncated: true},
		"both fit":                             {limit: len(arguments) + len(ok), arguments: string(arguments), result: ok},
		"a result that is over the limit too":  {limit: len(ok) - 1, result: `[]`, truncated: true},
		"a canceled run has no result to drop": {limit: len(arguments) - 1, canceled: true, truncated: true},
		"a canceled run whose arguments fit":   {limit: len(arguments), canceled: true, arguments: string(arguments)},
	} {
		t.Run(name, func(t *testing.T) {
			attrs := toolUseContentAttributes(test.limit, arguments, "ok", test.canceled)
			values := attributeValues(attrs)
			for key, want := range map[attribute.Key]string{
				GenAIToolCallArgumentsKey:     test.arguments,
				AnthropicToolResultContentKey: test.result,
			} {
				if got, ok := values[key]; got != want || ok != (want != "") {
					t.Errorf("%s = %q (present %t), want %q", key, got, ok, want)
				}
			}
			_, _, truncated := contentAttributes(t, attrs)
			wantTruncated(t, truncated, test.truncated)
		})
	}
}
