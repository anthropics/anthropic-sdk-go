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
//
// NOTE: this only exercises the inline-document request chain with ASCII /
// base64 payloads. See marshalto_proto.go for the full list of untested cases
// (HTML-special chars, quoted fields, unions with 0/2 variants, extras,
// overrides, explicit-null, time fields) that a generalized fix must cover.
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
