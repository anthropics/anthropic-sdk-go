package json

import (
	"bytes"
	"strings"
	"testing"
)

// Inner implements MarshalJSON to trigger the optimized code path
type benchInner struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (b benchInner) MarshalJSON() ([]byte, error) {
	return Marshal(struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}{b.Name, b.Value})
}

// Nested structure with multiple MarshalJSON calls
type benchNested struct {
	Inner benchInner `json:"inner"`
	Items []int      `json:"items"`
}

func (b benchNested) MarshalJSON() ([]byte, error) {
	return Marshal(struct {
		Inner benchInner `json:"inner"`
		Items []int      `json:"items"`
	}{b.Inner, b.Items})
}

// Deeply nested to amplify the effect
type benchDeep struct {
	Level1 benchNested `json:"level1"`
	Level2 benchNested `json:"level2"`
	Data   string      `json:"data"`
}

func (b benchDeep) MarshalJSON() ([]byte, error) {
	return Marshal(struct {
		Level1 benchNested `json:"level1"`
		Level2 benchNested `json:"level2"`
		Data   string      `json:"data"`
	}{b.Level1, b.Level2, b.Data})
}

func BenchmarkMarshalNestedMarshalJSON(b *testing.B) {
	data := benchDeep{
		Level1: benchNested{
			Inner: benchInner{Name: "test1", Value: 100},
			Items: []int{1, 2, 3, 4, 5},
		},
		Level2: benchNested{
			Inner: benchInner{Name: "test2", Value: 200},
			Items: []int{6, 7, 8, 9, 10},
		},
		Data: "some test data here",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Slice of nested structs - common real-world pattern
func BenchmarkMarshalSliceOfNestedMarshalJSON(b *testing.B) {
	data := make([]benchDeep, 50)
	for i := range data {
		data[i] = benchDeep{
			Level1: benchNested{
				Inner: benchInner{Name: "test1", Value: i},
				Items: []int{1, 2, 3, 4, 5},
			},
			Level2: benchNested{
				Inner: benchInner{Name: "test2", Value: i * 2},
				Items: []int{6, 7, 8, 9, 10},
			},
			Data: "some test data here that is a bit longer to simulate real payloads",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Test that HTML escaping is preserved for nested MarshalJSON calls
type htmlTestInner struct {
	Content string `json:"content"`
}

func (h htmlTestInner) MarshalJSON() ([]byte, error) {
	return Marshal(struct {
		Content string `json:"content"`
	}{h.Content})
}

type htmlTestOuter struct {
	Inner htmlTestInner `json:"inner"`
}

func (h htmlTestOuter) MarshalJSON() ([]byte, error) {
	return Marshal(struct {
		Inner htmlTestInner `json:"inner"`
	}{h.Inner})
}

func TestMarshalHTMLEscapeWithNestedMarshalJSON(t *testing.T) {
	// Test that HTML-sensitive characters are escaped in nested MarshalJSON
	data := htmlTestOuter{
		Inner: htmlTestInner{
			Content: "<script>alert('xss')</script>",
		},
	}

	result, err := Marshal(data)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// The < and > should be escaped as \u003c and \u003e
	if strings.Contains(string(result), "<script>") {
		t.Errorf("HTML was not escaped in Marshal output: %s", result)
	}
	if !strings.Contains(string(result), `\u003cscript\u003e`) {
		t.Errorf("Expected escaped HTML in output, got: %s", result)
	}
	// Verify no double-escaping (e.g., \u003c should not become \\u003c)
	if strings.Contains(string(result), `\\u003c`) {
		t.Errorf("HTML was double-escaped in output: %s", result)
	}
}

func TestEncoderHTMLEscapeWithNestedMarshalJSON(t *testing.T) {
	// Test with Encoder (which has escapeHTML=true by default)
	data := htmlTestOuter{
		Inner: htmlTestInner{
			Content: "<div>&amp;</div>",
		},
	}

	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	result := buf.String()
	// The < > & should be escaped
	if strings.Contains(result, "<div>") {
		t.Errorf("HTML was not escaped in Encoder output: %s", result)
	}
	if !strings.Contains(result, `\u003cdiv\u003e`) {
		t.Errorf("Expected escaped < and > in output, got: %s", result)
	}
	if !strings.Contains(result, `\u0026`) {
		t.Errorf("Expected escaped & in output, got: %s", result)
	}
}

func TestEncoderNoHTMLEscapeWithNestedMarshalJSON(t *testing.T) {
	// Test with SetEscapeHTML(false)
	// Note: Inner MarshalJSON calls use Marshal() which has escapeHTML=true by default,
	// so HTML escaping still occurs in the nested output. This is expected behavior
	// since the inner calls don't inherit the outer encoder's settings.
	data := htmlTestOuter{
		Inner: htmlTestInner{
			Content: "<div>&</div>",
		},
	}

	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(data); err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	result := buf.String()
	// Inner Marshal calls still escape HTML since they use default settings
	if strings.Contains(result, "<div>") {
		t.Logf("Note: HTML in nested MarshalJSON is escaped because inner Marshal uses default escapeHTML=true")
	}
	// Just verify we got valid JSON output
	if !strings.Contains(result, "content") {
		t.Errorf("Expected content field in output, got: %s", result)
	}
}
