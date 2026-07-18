package param

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	shimjson "github.com/anthropics/anthropic-sdk-go/internal/encoding/json"

	"github.com/tidwall/sjson"
)

// EncodedAsDate is not be stable and shouldn't be relied upon
type EncodedAsDate Opt[time.Time]

// If we want to set a literal key value into JSON using sjson, we need to make sure it doesn't have
// special characters that sjson interprets as a path.
var EscapeSJSONKey = strings.NewReplacer("\\", "\\\\", "|", "\\|", "#", "\\#", "@", "\\@", "*", "\\*", ".", "\\.", ":", "\\:", "?", "\\?").Replace

type forceOmit int

func (m EncodedAsDate) MarshalJSON() ([]byte, error) {
	underlying := Opt[time.Time](m)
	bytes := underlying.MarshalJSONWithTimeLayout("2006-01-02")
	if len(bytes) > 0 {
		return bytes, nil
	}
	return underlying.MarshalJSON()
}

// MarshalObject uses a shimmed 'encoding/json' from Go 1.24, to support the 'omitzero' tag
//
// Stability for the API of MarshalObject is not guaranteed.
func MarshalObject[T ParamStruct](f T, underlying any) ([]byte, error) {
	return MarshalWithExtras(f, underlying, f.extraFields())
}

// MarshalWithExtras is used to marshal a struct with additional properties.
//
// Stability for the API of MarshalWithExtras is not guaranteed.
func MarshalWithExtras[T ParamStruct, R any](f T, underlying any, extras map[string]R) ([]byte, error) {
	if f.null() {
		return []byte("null"), nil
	} else if len(extras) > 0 {
		bytes, err := shimjson.Marshal(underlying)
		if err != nil {
			return nil, err
		}
		for k, v := range extras {
			var a any = v
			if a == Omit {
				// Errors when handling ForceOmitted are ignored.
				if b, e := sjson.DeleteBytes(bytes, k); e == nil {
					bytes = b
				}
				continue
			}
			bytes, err = sjson.SetBytes(bytes, EscapeSJSONKey(k), v)
			if err != nil {
				return nil, err
			}
		}
		return bytes, nil
	} else if ovr, ok := f.Overrides(); ok {
		return shimjson.Marshal(ovr)
	} else {
		return shimjson.Marshal(underlying, shimjson.WithSkipCompaction(true))
	}
}

// PROTOTYPE(begin): buffer-direct variants of MarshalObject / MarshalUnion.
//
// These mirror MarshalWithExtras / MarshalUnion but, on the common fast path
// (no extras, no overrides, not explicit-null), encode the underlying value
// straight into the encoder's shared buffer via enc.Encode — so a nested
// param's payload is written once instead of being re-allocated as a fresh
// []byte at every nesting level. Slow paths fall back to the existing
// []byte-returning marshalers and WriteRaw the result.

// MarshalObjectTo is the buffer-direct counterpart of MarshalObject.
func MarshalObjectTo[T ParamStruct](enc *shimjson.DirectEncoder, f T, underlying any) {
	if f.null() {
		enc.WriteRaw([]byte("null"))
		return
	}
	if extras := f.extraFields(); len(extras) > 0 {
		b, err := MarshalWithExtras(f, underlying, extras)
		if err != nil {
			enc.Error(err)
			return
		}
		enc.WriteRaw(b)
		return
	}
	if ovr, ok := f.Overrides(); ok {
		b, err := shimjson.Marshal(ovr)
		if err != nil {
			enc.Error(err)
			return
		}
		enc.WriteRaw(b)
		return
	}
	enc.Encode(underlying)
}

// MarshalUnionTo is the buffer-direct counterpart of MarshalUnion. It mirrors
// MarshalUnion exactly — same variant-counting, the >1-present error, and the
// null/override handling when no variant is present — but encodes the present
// variant into the shared buffer instead of returning a fresh []byte.
func MarshalUnionTo[T ParamStruct](enc *shimjson.DirectEncoder, metadata T, variants ...any) {
	nPresent := 0
	presentIdx := -1
	for i, variant := range variants {
		if !IsOmitted(variant) {
			nPresent++
			presentIdx = i
		}
	}
	if nPresent == 0 || presentIdx == -1 {
		if metadata.null() {
			enc.WriteRaw([]byte("null"))
			return
		}
		if ovr, ok := metadata.Overrides(); ok {
			b, err := shimjson.Marshal(ovr)
			if err != nil {
				enc.Error(err)
				return
			}
			enc.WriteRaw(b)
			return
		}
		enc.WriteRaw([]byte("null"))
		return
	} else if nPresent > 1 {
		enc.Error(&json.MarshalerError{
			Type: typeFor[T](),
			Err:  fmt.Errorf("expected union to have only one present variant, got %d", nPresent),
		})
		return
	}
	enc.Encode(variants[presentIdx])
}

// PROTOTYPE(end)

// MarshalUnion uses a shimmed 'encoding/json' from Go 1.24, to support the 'omitzero' tag
//
// Stability for the API of MarshalUnion is not guaranteed.
func MarshalUnion[T ParamStruct](metadata T, variants ...any) ([]byte, error) {
	nPresent := 0
	presentIdx := -1
	for i, variant := range variants {
		if !IsOmitted(variant) {
			nPresent++
			presentIdx = i
		}
	}
	if nPresent == 0 || presentIdx == -1 {
		if metadata.null() {
			return []byte("null"), nil
		}
		if ovr, ok := metadata.Overrides(); ok {
			return shimjson.Marshal(ovr)
		}
		return []byte(`null`), nil
	} else if nPresent > 1 {
		return nil, &json.MarshalerError{
			Type: typeFor[T](),
			Err:  fmt.Errorf("expected union to have only one present variant, got %d", nPresent),
		}
	}
	return shimjson.Marshal(variants[presentIdx], shimjson.WithSkipCompaction(true))
}

// typeFor is shimmed from Go 1.23 "reflect" package
func typeFor[T any]() reflect.Type {
	var v T
	if t := reflect.TypeOf(v); t != nil {
		return t // optimize for T being a non-interface kind
	}
	return reflect.TypeOf((*T)(nil)).Elem() // only for an interface kind
}
