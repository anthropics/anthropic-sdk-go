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
func MarshalBufferDirect(v any) ([]byte, error) {
	return shimjson.Marshal(v, shimjson.WithBufferDirect(true))
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

func (u ContentBlockParamUnion) MarshalJSONTo(enc *shimjson.DirectEncoder) {
	param.MarshalUnionValueTo(enc, u.asAny())
}

func (u DocumentBlockParamSourceUnion) MarshalJSONTo(enc *shimjson.DirectEncoder) {
	param.MarshalUnionValueTo(enc, u.asAny())
}
