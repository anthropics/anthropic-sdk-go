package anthropic_test

// Union variant selection for request params: the variant whose
// constants match wins; otherwise the closest structural fit.

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
)

func TestParamUnmarshalNullAnyFieldDecodesToNil(t *testing.T) {
	raw := `{"model": "m", "max_tokens": 1,
		"messages": [{"role": "assistant", "content": [{"type": "tool_use", "id": "t1", "name": "f", "input": null}]}]}`
	var p anthropic.MessageNewParams
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tu := p.Messages[0].Content[0].OfToolUse
	if tu == nil {
		t.Fatalf("tool_use variant not selected, got %#v", p.Messages[0].Content[0])
	}
	if tu.Input != nil {
		t.Errorf("input must decode to nil as encoding/json does, got %#v", tu.Input)
	}
}

// Unknown keys are dropped and never touch SetExtraFields.
func TestParamUnmarshalUnknownFieldsIgnored(t *testing.T) {
	raw := `{"model": "m", "max_tokens": 1, "some_future_field": {"a": 1},
		"messages": [{"role": "user", "content": [{"type": "text", "text": "hi", "x_meta": 1}]}]}`
	var p anthropic.MessageNewParams
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Messages[0].Content[0].OfText == nil || p.Messages[0].Content[0].OfText.Text != "hi" {
		t.Errorf("text block not decoded alongside unknown field, got %#v", p.Messages[0].Content[0])
	}
	if len(p.ExtraFields()) != 0 || len(p.Messages[0].Content[0].OfText.ExtraFields()) != 0 {
		t.Errorf("unknown keys must not populate the SetExtraFields escape hatch")
	}
	out, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	for _, dropped := range []string{"x_meta", "some_future_field"} {
		if strings.Contains(string(out), dropped) {
			t.Errorf("unknown key %s must be dropped on re-marshal, got %s", dropped, out)
		}
	}
}

func TestParamUnmarshalToolUnionVariantIdentity(t *testing.T) {
	cases := map[string]struct {
		raw    string
		picked func(u anthropic.ToolUnionParam) bool
	}{
		"web_search_20250305": {
			`{"type": "web_search_20250305", "name": "web_search", "max_uses": 3}`,
			func(u anthropic.ToolUnionParam) bool {
				return u.OfWebSearchTool20250305 != nil && u.OfWebSearchTool20250305.MaxUses.Value == 3
			},
		},
		"text_editor_20250429": {
			`{"type": "text_editor_20250429", "name": "str_replace_based_edit_tool"}`,
			func(u anthropic.ToolUnionParam) bool { return u.OfTextEditor20250429 != nil },
		},
		"code_execution_20260120": {
			`{"type": "code_execution_20260120", "name": "code_execution"}`,
			func(u anthropic.ToolUnionParam) bool { return u.OfCodeExecutionTool20260120 != nil },
		},
		"web_fetch_20250910": {
			`{"type": "web_fetch_20250910", "name": "web_fetch"}`,
			func(u anthropic.ToolUnionParam) bool { return u.OfWebFetchTool20250910 != nil },
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var u anthropic.ToolUnionParam
			if err := json.Unmarshal([]byte(tc.raw), &u); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.picked(u) {
				t.Fatalf("wrong variant selected for %s: %#v", tc.raw, u)
			}
		})
	}
}

// An unknown newer version lands in the closest sibling, so its
// config survives the round trip.
func TestParamUnmarshalUnknownToolVersionPreservesConfig(t *testing.T) {
	raw := `{"type":"web_search_20990101","name":"web_search","max_uses":5,"allowed_domains":["example.com"]}`
	var u anthropic.ToolUnionParam
	if err := json.Unmarshal([]byte(raw), &u); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.OfTool != nil {
		t.Fatalf("must not fall through to the custom-tool variant (drops config): %#v", u)
	}
	out, err := json.Marshal(u)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	assertJSONEqual(t, raw, string(out))
}

func TestParamUnmarshalCustomToolPicksOfTool(t *testing.T) {
	for _, raw := range []string{
		`{"name":"my_tool","input_schema":{"type":"object"}}`,
		`{"type":"custom","name":"my_tool","input_schema":{"type":"object"}}`,
	} {
		var u anthropic.ToolUnionParam
		if err := json.Unmarshal([]byte(raw), &u); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if u.OfTool == nil {
			t.Fatalf("expected the custom-tool variant for %s, got %#v", raw, u)
		}
	}
}

func assertJSONEqual(t *testing.T, want, got string) {
	t.Helper()
	var a, b any
	if err := json.Unmarshal([]byte(want), &a); err != nil {
		t.Fatalf("bad want json: %v", err)
	}
	if err := json.Unmarshal([]byte(got), &b); err != nil {
		t.Fatalf("bad got json: %v", err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("round trip lost data\n want %s\n  got %s", want, got)
	}
}

// Bare-string content promotes to a one-element text-block slice.
func TestParamUnmarshalStringPromotion(t *testing.T) {
	var p anthropic.MessageNewParams
	if err := json.Unmarshal([]byte(`{"model":"m","max_tokens":1,"system":"be brief"}`), &p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(p.System) != 1 || p.System[0].Text != "be brief" {
		t.Errorf("system string not promoted, got %#v", p.System)
	}
	var bm anthropic.BetaMessageParam
	if err := json.Unmarshal([]byte(`{"role":"user","content":"hello"}`), &bm); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bm.Content) != 1 || bm.Content[0].OfText == nil || bm.Content[0].OfText.Text != "hello" {
		t.Errorf("beta content string not promoted, got %#v", bm.Content)
	}
}

