package apijson

import (
	"errors"
	"reflect"
	"sync"

	"github.com/anthropics/anthropic-sdk-go/packages/param"

	"github.com/tidwall/gjson"
)

var apiUnionType = reflect.TypeOf(param.APIUnion{})
var apiObjectType = reflect.TypeOf(param.APIObject{})

var nativeTypeCache sync.Map // map[reflect.Type]bool

// CustomUnmarshaler indicates that a type is not "apijson native".
// apijson native types use the custom [json.Unmarshaler] only as an entry
// point, and therefore are skipped when nested.
type CustomUnmarshaler interface {
	UnmarshalAPIJSON([]byte) error
}

var customUnmarshalerType = reflect.TypeOf((*CustomUnmarshaler)(nil)).Elem()

func implementsCustomUnmarshaler(t reflect.Type) bool {
	return t.Implements(customUnmarshalerType) || reflect.PointerTo(t).Implements(customUnmarshalerType)
}

// isApijsonNative reports whether t is a generated struct whose
// UnmarshalJSON is the trivial delegation to [UnmarshalRoot], therefore
// the UnmarshalJSON function can be skipped for nested values.
func isApijsonNative(t reflect.Type) bool {
	if cached, ok := nativeTypeCache.Load(t); ok {
		return cached.(bool)
	}
	native := computeApijsonNative(t)
	nativeTypeCache.Store(t, native)
	return native
}

func computeApijsonNative(t reflect.Type) bool {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return false
	}
	if implementsCustomUnmarshaler(t) {
		return false
	}
	if isParamStruct(t) || isStructUnion(t) {
		return true
	}
	if f, ok := t.FieldByName("JSON"); ok && f.Type.Kind() == reflect.Struct {
		if _, ok := f.Type.FieldByName("raw"); ok {
			return true
		}
	}
	return false
}

var paramObjectIdxCache sync.Map // map[reflect.Type]int

func paramObjectIndex(t reflect.Type) int {
	return paramEmbedIndex(t, apiObjectType, &paramObjectIdxCache)
}

func paramEmbedIndex(t reflect.Type, embed reflect.Type, cache *sync.Map) int {
	if cached, ok := cache.Load(t); ok {
		return cached.(int)
	}
	idx := -1
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.Anonymous && f.Type == embed {
				idx = i
				break
			}
		}
	}
	cache.Store(t, idx)
	return idx
}

// isParamStruct reports whether t is a generated request param struct.
func isParamStruct(t reflect.Type) bool {
	return paramObjectIndex(t) >= 0
}

func isStructUnion(t reflect.Type) bool {
	if t.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Type == apiUnionType && t.Field(i).Anonymous {
			return true
		}
	}
	return false
}

func RegisterDiscriminatedUnion[T any](key string, mappings map[string]reflect.Type) {
	var t T
	entry := unionEntry{
		discriminatorKey: key,
		variants:         []UnionVariant{},
	}
	for k, typ := range mappings {
		entry.variants = append(entry.variants, UnionVariant{
			DiscriminatorValue: k,
			Type:               typ,
		})
	}
	unionRegistry[reflect.TypeOf(t)] = entry
}

// deriveDiscriminator picks the `default:`-tagged field whose values
// route the most variants uniquely, and returns that field's JSON
// name plus a value→variant-index map.
func deriveDiscriminator(variants []reflect.StructField) (key string, index map[string]int) {
	byKey := map[string]map[string]int{}
	for idx, v := range variants {
		t := v.Type
		for t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			continue
		}
		for i := 0; i < t.NumField(); i++ {
			ptag, ok := parseJSONStructTag(t.Field(i))
			if !ok {
				continue
			}
			def, ok := ptag.defaultValue.(string)
			if !ok {
				continue
			}
			m := byKey[ptag.name]
			if m == nil {
				m = map[string]int{}
				byKey[ptag.name] = m
			}
			if _, dup := m[def]; dup {
				m[def] = -1
			} else {
				m[def] = idx
			}
		}
	}
	for k, m := range byKey {
		usable := map[string]int{}
		for v, idx := range m {
			if idx >= 0 {
				usable[v] = idx
			}
		}
		if len(usable) > len(index) || (len(usable) == len(index) && k < key) {
			key, index = k, usable
		}
	}
	return key, index
}

func (d *decoderBuilder) newStructUnionDecoder(t reflect.Type) decoderFunc {
	type variantDecoder struct {
		decoder decoderFunc
		field   reflect.StructField
	}
	var decoders []variantDecoder
	var fields []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous && field.Type == apiUnionType {
			continue
		}
		decoders = append(decoders, variantDecoder{d.typeDecoder(field.Type), field})
		fields = append(fields, field)
	}

	// To more performantly pick the correct variant, derive a
	// discriminator from the `default:` tags.
	discrimKey, discrimIndex := deriveDiscriminator(fields)

	type discriminatedDecoder struct {
		variantDecoder
		discriminator any
	}
	discriminatedDecoders := []discriminatedDecoder{}
	unionEntry, discriminated := unionRegistry[t]
	for _, variant := range unionEntry.variants {
		// For each union variant, find a matching decoder and save it
		for _, decoder := range decoders {
			if decoder.field.Type.Elem() == variant.Type {
				discriminatedDecoders = append(discriminatedDecoders, discriminatedDecoder{
					decoder,
					variant.DiscriminatorValue,
				})
				break
			}
		}
	}

	return func(n gjson.Result, v reflect.Value, state *decoderState) error {
		if discriminated && n.Type == gjson.JSON && len(unionEntry.discriminatorKey) != 0 {
			discriminator := n.Get(EscapeSJSONKey(unionEntry.discriminatorKey)).Value()
			for _, decoder := range discriminatedDecoders {
				if discriminator == decoder.discriminator {
					inner := v.FieldByIndex(decoder.field.Index)
					return decoder.decoder(n, inner, state)
				}
			}
			return errors.New("apijson: was not able to find discriminated union variant")
		}

		scanned := discrimKey != "" && n.Type == gjson.JSON
		if scanned {
			if i, ok := discrimIndex[n.Get(EscapeSJSONKey(discrimKey)).String()]; ok {
				sub := decoderState{strict: state.strict}
				inner := v.FieldByIndex(decoders[i].field.Index)
				if decoders[i].decoder(n, inner, &sub) == nil {
					state.exactness.absorb(sub.exactness)
					return nil
				}
			}
		}

		var best exactness
		bestVariant := -1
		for i, decoder := range decoders {
			// Pointers are used to discern JSON object variants from value variants
			if n.Type != gjson.JSON && decoder.field.Type.Kind() == reflect.Ptr {
				continue
			}

			sub := decoderState{strict: state.strict}
			inner := v.FieldByIndex(decoder.field.Index)
			if err := decoder.decoder(n, inner, &sub); err != nil {
				continue
			}
			if bestVariant < 0 || sub.exactness.betterThan(best) {
				best, bestVariant = sub.exactness, i
				// A perfect fit cannot be beaten once the discriminator
				// scan (or the fit itself) has ruled out a later hit.
				if best.perfect() && (scanned || best.constHit()) {
					break
				}
			}
		}

		if bestVariant < 0 {
			return errors.New("apijson: was not able to coerce type as union")
		}

		if state.strict && !best.perfect() {
			return errors.New("apijson: was not able to coerce type as union strictly")
		}
		state.exactness.absorb(best)

		for i := 0; i < len(decoders); i++ {
			if i == bestVariant {
				continue
			}
			v.FieldByIndex(decoders[i].field.Index).SetZero()
		}

		return nil
	}
}

// newUnionDecoder returns a decoderFunc that deserializes into a union using an
// algorithm roughly similar to Pydantic's [smart algorithm].
//
// Conceptually this is equivalent to choosing the best schema based on how 'exact'
// the deserialization is for each of the schemas.
//
// If there is a tie in the level of exactness, then the tie is broken
// left-to-right.
//
// [smart algorithm]: https://docs.pydantic.dev/latest/concepts/unions/#smart-mode
func (d *decoderBuilder) newUnionDecoder(t reflect.Type) decoderFunc {
	unionEntry, ok := unionRegistry[t]
	if !ok {
		panic("apijson: couldn't find union of type " + t.String() + " in union registry")
	}
	decoders := []decoderFunc{}
	for _, variant := range unionEntry.variants {
		decoder := d.typeDecoder(variant.Type)
		decoders = append(decoders, decoder)
	}
	return func(n gjson.Result, v reflect.Value, state *decoderState) error {
		// If there is a discriminator match, circumvent the exactness logic entirely
		for idx, variant := range unionEntry.variants {
			decoder := decoders[idx]
			if variant.TypeFilter != n.Type {
				continue
			}

			if len(unionEntry.discriminatorKey) != 0 {
				discriminatorValue := n.Get(EscapeSJSONKey(unionEntry.discriminatorKey)).Value()
				if discriminatorValue == variant.DiscriminatorValue {
					inner := reflect.New(variant.Type).Elem()
					err := decoder(n, inner, state)
					v.Set(inner)
					return err
				}
			}
		}

		var bestExactness exactness
		found := false
		for idx, variant := range unionEntry.variants {
			decoder := decoders[idx]
			if variant.TypeFilter != n.Type {
				continue
			}
			sub := decoderState{strict: state.strict}
			inner := reflect.New(variant.Type).Elem()
			err := decoder(n, inner, &sub)
			if err != nil {
				continue
			}
			if !found || sub.exactness.betterThan(bestExactness) {
				v.Set(inner)
				bestExactness = sub.exactness
				found = true
				if bestExactness.constHit() && bestExactness.perfect() {
					return nil
				}
			}
		}

		if !found {
			return errors.New("apijson: was not able to coerce type as union")
		}

		if state.strict && !bestExactness.perfect() {
			return errors.New("apijson: was not able to coerce type as union strictly")
		}
		state.exactness.absorb(bestExactness)

		return nil
	}
}
