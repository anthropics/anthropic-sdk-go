package anthropic

// PROTOTYPE: buffer-direct marshaling for fix #1 of the request-marshal memory
// amplification report. These MarshalJSONTo methods let the encoder write each
// nested param's payload once into a single shared buffer (instead of every
// level allocating a fresh []byte of its whole subtree). Only the types on the
// inline-document request chain are implemented here; the real fix would have
// the code generator emit MarshalJSONTo for all param types.
//
// Enabled per-call via MarshalBufferDirect; the stock MarshalJSON path is
// unchanged, so before/after can be benchmarked in one binary.

import (
	shimjson "github.com/anthropics/anthropic-sdk-go/internal/encoding/json"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

// MarshalBufferDirect marshals v with buffer-direct encoding enabled.
//
// WithSkipCompaction(true) matches the option the stock param marshalers pass,
// so any subtree that falls back to the []byte MarshalJSON path (a type without
// MarshalJSONTo) is incorporated exactly as stock would — keeping output
// byte-identical regardless of how deep the implemented chain reaches.
func MarshalBufferDirect(v any) ([]byte, error) {
	return shimjson.Marshal(v, shimjson.WithBufferDirect(true), shimjson.WithSkipCompaction(true))
}

func (r MessageNewParams) MarshalJSONTo(enc *shimjson.DirectEncoder) {
	type shadow MessageNewParams
	param.MarshalObjectTo(enc, r, (*shadow)(&r))
}

func (r MessageParam) MarshalJSONTo(enc *shimjson.DirectEncoder) {
	type shadow MessageParam
	param.MarshalObjectTo(enc, r, (*shadow)(&r))
}

func (r DocumentBlockParam) MarshalJSONTo(enc *shimjson.DirectEncoder) {
	type shadow DocumentBlockParam
	param.MarshalObjectTo(enc, r, (*shadow)(&r))
}

func (r Base64PDFSourceParam) MarshalJSONTo(enc *shimjson.DirectEncoder) {
	type shadow Base64PDFSourceParam
	param.MarshalObjectTo(enc, r, (*shadow)(&r))
}

// MarshalJSONTo mirrors ContentBlockParamUnion.MarshalJSON's variant list so the
// >1-present error and null/override handling match stock exactly.
func (u ContentBlockParamUnion) MarshalJSONTo(enc *shimjson.DirectEncoder) {
	param.MarshalUnionTo(enc, u, u.OfText,
		u.OfImage,
		u.OfDocument,
		u.OfSearchResult,
		u.OfThinking,
		u.OfRedactedThinking,
		u.OfToolUse,
		u.OfToolResult,
		u.OfServerToolUse,
		u.OfWebSearchToolResult,
		u.OfWebFetchToolResult,
		u.OfCodeExecutionToolResult,
		u.OfBashCodeExecutionToolResult,
		u.OfTextEditorCodeExecutionToolResult,
		u.OfToolSearchToolResult,
		u.OfContainerUpload,
		u.OfMidConvSystem)
}

func (u DocumentBlockParamSourceUnion) MarshalJSONTo(enc *shimjson.DirectEncoder) {
	param.MarshalUnionTo(enc, u, u.OfBase64, u.OfText, u.OfContent, u.OfURL)
}
