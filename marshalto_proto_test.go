package anthropic_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	anthropic "github.com/anthropics/anthropic-sdk-go"
)

// payload returns an n-byte string. 'A' is in the base64 alphabet and needs no
// JSON escaping, matching a real base64 document body; only the size drives cost.
func payload(n int) string { return strings.Repeat("A", n) }

func documentParams(n int) anthropic.MessageNewParams {
	return anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5,
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewDocumentBlock(anthropic.Base64PDFSourceParam{Data: payload(n)}),
			),
		},
	}
}

// TestBufferDirectMatchesStock is the correctness gate: buffer-direct output
// MUST be byte-for-byte identical to stock MarshalJSON on the document chain.
func TestBufferDirectMatchesStock(t *testing.T) {
	for _, n := range []int{0, 1, 100, 1 << 20} {
		p := documentParams(n)
		stock, err := p.MarshalJSON()
		if err != nil {
			t.Fatalf("stock MarshalJSON: %v", err)
		}
		direct, err := anthropic.MarshalBufferDirect(p)
		if err != nil {
			t.Fatalf("MarshalBufferDirect: %v", err)
		}
		if !bytes.Equal(stock, direct) {
			t.Fatalf("n=%d output mismatch:\n stock=%s\ndirect=%s", n, stock, direct)
		}
	}
}

// assertSameAsStock marshals p both ways and fails unless the output is
// byte-for-byte identical (and the error behavior matches).
func assertSameAsStock(t *testing.T, name string, p anthropic.MessageNewParams) {
	t.Helper()
	stock, stockErr := p.MarshalJSON()
	direct, directErr := anthropic.MarshalBufferDirect(p)
	if (stockErr == nil) != (directErr == nil) {
		t.Fatalf("%s: error mismatch: stock=%v direct=%v", name, stockErr, directErr)
	}
	if stockErr != nil {
		return // both errored; acceptable
	}
	if !bytes.Equal(stock, direct) {
		t.Fatalf("%s: output mismatch:\n stock=%s\ndirect=%s", name, stock, direct)
	}
}

// TestBufferDirectDifferential is the broad regression gate. It exercises the
// previously-untested cases — strings with HTML-special chars and control
// runes, multiple content blocks, fallback-to-stock types (TextBlockParam has
// no MarshalJSONTo), extra fields, Override, explicit-null, unions with 0/1
// present variants, and metadata — asserting buffer-direct stays byte-identical
// to stock MarshalJSON across all of them.
func TestBufferDirectDifferential(t *testing.T) {
	base := func() anthropic.MessageNewParams {
		return anthropic.MessageNewParams{Model: anthropic.ModelClaudeSonnet4_5, MaxTokens: 1024}
	}

	t.Run("text_with_html_special_chars", func(t *testing.T) {
		p := base()
		p.Messages = []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(`a < b && c > d, "quoted" </script>`)),
		}
		assertSameAsStock(t, "html_special", p)
	})

	t.Run("text_with_unicode_and_control", func(t *testing.T) {
		p := base()
		p.Messages = []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("emoji 😀     tab\tnewline\n nul\x00 end")),
		}
		assertSameAsStock(t, "unicode_control", p)
	})

	t.Run("mixed_content_blocks", func(t *testing.T) {
		p := base()
		p.Messages = []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock("before <b>"),
				anthropic.NewDocumentBlock(anthropic.Base64PDFSourceParam{Data: payload(256)}),
				anthropic.NewTextBlock("after & co"),
			),
		}
		assertSameAsStock(t, "mixed_blocks", p)
	})

	t.Run("system_prompt_with_specials", func(t *testing.T) {
		p := base()
		p.System = []anthropic.TextBlockParam{{Text: `system < & > "x"`}}
		p.Messages = []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hi"))}
		assertSameAsStock(t, "system_specials", p)
	})

	t.Run("metadata_user_id", func(t *testing.T) {
		p := base()
		p.Metadata = anthropic.MetadataParam{UserID: anthropic.String("user <id> & more")}
		p.Messages = []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hi"))}
		assertSameAsStock(t, "metadata", p)
	})

	t.Run("extra_fields_on_request", func(t *testing.T) {
		p := base()
		p.Messages = []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock("hi"))}
		p.SetExtraFields(map[string]any{"x_custom": "value & <stuff>", "x_num": 42})
		assertSameAsStock(t, "extra_fields", p)
	})

	t.Run("extra_fields_on_nested_block", func(t *testing.T) {
		p := base()
		tb := anthropic.TextBlockParam{Text: "hi"}
		tb.SetExtraFields(map[string]any{"x_nested": "<v>"})
		p.Messages = []anthropic.MessageParam{{Role: anthropic.MessageParamRoleUser, Content: []anthropic.ContentBlockParamUnion{{OfText: &tb}}}}
		assertSameAsStock(t, "extra_nested", p)
	})

	t.Run("empty_messages", func(t *testing.T) {
		p := base()
		p.Messages = []anthropic.MessageParam{}
		assertSameAsStock(t, "empty_messages", p)
	})

	t.Run("document_union_single_variant", func(t *testing.T) {
		p := base()
		p.Messages = []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewDocumentBlock(anthropic.Base64PDFSourceParam{Data: payload(128)})),
		}
		assertSameAsStock(t, "doc_union_one", p)
	})

	t.Run("content_union_empty_must_error_same_way", func(t *testing.T) {
		// A ContentBlockParamUnion with no variant set: both paths must agree
		// (stock errors "expected union to have only one present variant").
		p := base()
		p.Messages = []anthropic.MessageParam{
			{Role: anthropic.MessageParamRoleUser, Content: []anthropic.ContentBlockParamUnion{{}}},
		}
		assertSameAsStock(t, "content_union_empty", p)
	})

	t.Run("nested_text_citations_with_specials", func(t *testing.T) {
		p := base()
		p.Messages = []anthropic.MessageParam{
			anthropic.NewAssistantMessage(anthropic.NewTextBlock("answer with <tag> & \"quote\"")),
		}
		assertSameAsStock(t, "citations", p)
	})
}

var sizesMiB = []int{1, 4, 8}

// BenchmarkStockMarshalDocumentRequest measures the current per-request
// allocation: ~7x the serialized body size for an inline base64 document.
func BenchmarkStockMarshalDocumentRequest(b *testing.B) {
	for _, mib := range sizesMiB {
		b.Run(fmt.Sprintf("%dMiB", mib), func(b *testing.B) {
			p := documentParams(mib << 20)
			b.ReportAllocs()
			b.ResetTimer()
			var out []byte
			for i := 0; i < b.N; i++ {
				var err error
				out, err = p.MarshalJSON()
				if err != nil {
					b.Fatal(err)
				}
			}
			b.StopTimer()
			b.ReportMetric(float64(len(out)), "outputBytes")
		})
	}
}

// BenchmarkBufferDirectDocumentRequest measures the buffer-direct path: ~1.4x
// the body size (the per-nesting-level amplification is removed; the residual is
// the final whole-buffer copy + buffer doubling — see follow-ups in the report).
func BenchmarkBufferDirectDocumentRequest(b *testing.B) {
	for _, mib := range sizesMiB {
		b.Run(fmt.Sprintf("%dMiB", mib), func(b *testing.B) {
			p := documentParams(mib << 20)
			b.ReportAllocs()
			b.ResetTimer()
			var out []byte
			for i := 0; i < b.N; i++ {
				var err error
				out, err = anthropic.MarshalBufferDirect(p)
				if err != nil {
					b.Fatal(err)
				}
			}
			b.StopTimer()
			b.ReportMetric(float64(len(out)), "outputBytes")
		})
	}
}

// BenchmarkStdlibFlatBaseline marshals the same-sized payload as a single string
// field with the standard library — the ~1x floor the SDK should approach.
func BenchmarkStdlibFlatBaseline(b *testing.B) {
	for _, mib := range sizesMiB {
		b.Run(fmt.Sprintf("%dMiB", mib), func(b *testing.B) {
			v := map[string]string{"data": payload(mib << 20)}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := json.Marshal(v); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// FuzzBufferDirectMatchesStock fuzzes the string-escaping path (the likeliest
// source of an escapeHTML/quoted divergence) by feeding arbitrary text into a
// content block and asserting buffer-direct stays byte-identical to stock.
func FuzzBufferDirectMatchesStock(f *testing.F) {
	for _, s := range []string{"", "a < b & c > d", "\"q\"", "\\n\t\x00", "😀/script", "{\"k\":1}"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, text string) {
		p := anthropic.MessageNewParams{
			Model:     anthropic.ModelClaudeSonnet4_5,
			MaxTokens: 1024,
			System:    []anthropic.TextBlockParam{{Text: text}},
			Messages:  []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock(text))},
		}
		stock, errA := p.MarshalJSON()
		direct, errB := anthropic.MarshalBufferDirect(p)
		if (errA == nil) != (errB == nil) {
			t.Fatalf("error mismatch for %q: stock=%v direct=%v", text, errA, errB)
		}
		if errA == nil && !bytes.Equal(stock, direct) {
			t.Fatalf("mismatch for %q:\n stock=%s\ndirect=%s", text, stock, direct)
		}
	})
}
