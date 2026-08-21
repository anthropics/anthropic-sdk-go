// Package compat tests the hand-written parts of the SDK against the generated
// types, so that a field or variant added by codegen fails a test here instead
// of going unhandled.
//
// This file synthesizes response JSON from those types: one fully populated
// value per type, plus one variation per field (empty list, zero, null, and
// each union variant), using the matching param struct to decide which fields a
// request can carry. It lives in this package because reflecting over the
// As<Variant>() accessors would defeat dead-code elimination in a package users
// link.
package compat

import (
	"encoding/json"
	"maps"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// A variant is one member of a response union.
type variant struct {
	// Accessor is the As* method that returns the variant.
	Accessor string
	// Type is what the accessor returns: a struct, or a slice for a list member.
	Type reflect.Type
}

// unionVariants returns the members of a response union in accessor order, or
// nothing for a type that is not a union.
func unionVariants(union reflect.Type) []variant {
	var out []variant
	for i := 0; i < union.NumMethod(); i++ {
		m := union.Method(i)
		if !strings.HasPrefix(m.Name, "As") || m.Type.NumIn() != 1 || m.Type.NumOut() != 1 {
			continue
		}
		// AsAny returns the interface the variants share, not a variant.
		if typ := m.Type.Out(0); typ.Kind() != reflect.Interface {
			out = append(out, variant{m.Name, typ})
		}
	}
	return out
}

// wireType returns the discriminator a variant carries on the wire.
func wireType(t *testing.T, v variant) string {
	t.Helper()
	field, ok := v.Type.FieldByName("Type")
	if !ok {
		t.Fatalf("%s has no Type field", v.Type)
	}
	return field.Tag.Get("default")
}

// A shape is one synthetic wire value for a response type, ready for
// json.Marshal.
type shape struct {
	// Name is "populated" or the path of the one difference from it, such as
	// "content=AsResponseWebFetchToolResultError/error_code=empty".
	Name  string
	Value any
}

// requestJSON marshals a param and normalizes it with comparableJSON.
func requestJSON(t *testing.T, param any) []byte {
	t.Helper()
	data, err := json.Marshal(param)
	if err != nil {
		t.Fatalf("Failed to marshal param: %v", err)
	}
	return comparableJSON(t, data)
}

// comparableJSON re-encodes JSON with sorted keys so that two documents compare
// with bytes.Equal. It drops nulls, which a request may omit, but keeps empty
// lists and zero values, since losing those is what the tests catch.
func comparableJSON(t *testing.T, data []byte) []byte {
	t.Helper()
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		t.Fatalf("Failed to parse JSON %s: %v", data, err)
	}
	out, err := json.Marshal(withoutNulls(v))
	if err != nil {
		t.Fatalf("Failed to re-encode JSON: %v", err)
	}
	return out
}

func withoutNulls(v any) any {
	switch v := v.(type) {
	case map[string]any:
		for key, member := range v {
			if member == nil {
				delete(v, key)
				continue
			}
			v[key] = withoutNulls(member)
		}
	case []any:
		for i := range v {
			v[i] = withoutNulls(v[i])
		}
	}
	return v
}

// maxDepth exceeds any real nesting, so reaching it means the walk found a
// cycle.
const maxDepth = 8

// A synthesizer generates response shapes for one API surface, consulting its
// param structs for what a request can carry.
type synthesizer struct {
	t *testing.T
	// params maps each wire "type" to the param structs that send it, which may
	// be several.
	params map[string][]reflect.Type
	// skipped records each response field with no request counterpart as
	// "Struct.key".
	skipped map[string]bool
}

// newSynthesizer indexes every param struct reachable from requestUnion.
func newSynthesizer(t *testing.T, requestUnion reflect.Type) *synthesizer {
	s := &synthesizer{t: t, params: map[string][]reflect.Type{}, skipped: map[string]bool{}}
	s.indexParams(requestUnion, map[reflect.Type]bool{})
	return s
}

func (s *synthesizer) indexParams(typ reflect.Type, seen map[reflect.Type]bool) {
	for typ.Kind() == reflect.Pointer || typ.Kind() == reflect.Slice || typ.Kind() == reflect.Map {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct || seen[typ] {
		return
	}
	seen[typ] = true
	if wire := discriminator(typ); wire != "" {
		s.params[wire] = append(s.params[wire], typ)
	}
	for i := 0; i < typ.NumField(); i++ {
		s.indexParams(typ.Field(i).Type, seen)
	}
}

// discriminator returns a struct's constant wire "type", or "" if it has none.
func discriminator(typ reflect.Type) string {
	if field, ok := typ.FieldByName("Type"); ok {
		return field.Tag.Get("default")
	}
	return ""
}

// responseShapes returns the populated shape of typ followed by its variations.
func (s *synthesizer) responseShapes(typ reflect.Type) []shape {
	s.t.Helper()
	return s.shapes(typ, nil, 0)
}

// shapes generates typ, with param naming the request-side field that carries
// it when known.
func (s *synthesizer) shapes(typ reflect.Type, param *reflect.StructField, depth int) []shape {
	if depth > maxDepth {
		s.t.Fatalf("%s nests deeper than %d levels", typ, maxDepth)
	}
	if typ == reflect.TypeFor[json.RawMessage]() {
		return []shape{{"populated", map[string]any{"key": "value"}}, {"empty object", map[string]any{}}}
	}
	// constant.* discriminators report their only value through Default.
	if def := reflect.Zero(typ).MethodByName("Default"); def.IsValid() {
		return []shape{{"populated", def.Call(nil)[0].Interface()}}
	}
	if variants := unionVariants(typ); len(variants) > 0 {
		return s.unionShapes(variants, depth)
	}
	switch typ.Kind() {
	case reflect.Pointer:
		return s.shapes(typ.Elem(), param, depth)
	case reflect.Struct:
		return s.structShapes(typ, param, depth)
	case reflect.Slice:
		return s.sliceShapes(typ, depth)
	case reflect.String:
		if typ.PkgPath() != "" {
			// Named string types are enums, which are never empty.
			return []shape{{"populated", "value"}}
		}
		return []shape{{"populated", "value"}, {"empty", ""}}
	case reflect.Bool:
		return []shape{{"populated", true}, {"false", false}}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []shape{{"populated", 7}, {"zero", 0}}
	case reflect.Float32, reflect.Float64:
		return []shape{{"populated", 7.5}, {"zero", 0}}
	case reflect.Interface, reflect.Map:
		return []shape{{"populated", map[string]any{"key": "value"}}}
	}
	s.t.Fatalf("Don't know how to generate a %s", typ)
	return nil
}

func (s *synthesizer) unionShapes(variants []variant, depth int) []shape {
	var out []shape
	for _, v := range variants {
		for _, sh := range s.shapes(v.Type, nil, depth+1) {
			out = append(out, shape{v.Accessor + "/" + sh.Name, sh.Value})
		}
	}
	return out
}

func (s *synthesizer) sliceShapes(typ reflect.Type, depth int) []shape {
	items := s.shapes(typ.Elem(), nil, depth+1)
	out := []shape{{"populated", []any{items[0].Value}}, {"empty", []any{}}}
	for _, item := range items[1:] {
		out = append(out, shape{"item " + item.Name, []any{item.Value}})
	}
	return out
}

// structShapes populates every field the request can carry, then varies one at
// a time while the rest stay populated.
func (s *synthesizer) structShapes(typ reflect.Type, param *reflect.StructField, depth int) []shape {
	type variation struct {
		key string
		shape
	}
	paramStruct := s.paramStruct(typ, param)
	populated := map[string]any{}
	var variations []variation
	for i := 0; i < typ.NumField(); i++ {
		key := wireKey(typ.Field(i))
		if key == "" {
			continue
		}
		paramField, ok := requestField(paramStruct, key)
		if paramStruct != nil && !ok {
			s.skipped[typ.Name()+"."+key] = true
			continue
		}
		fieldShapes := s.fieldShapes(typ.Field(i).Type, paramField, depth)
		populated[key] = fieldShapes[0].Value
		for _, sh := range fieldShapes[1:] {
			variations = append(variations, variation{key, sh})
		}
	}
	out := []shape{{"populated", populated}}
	for _, v := range variations {
		obj := maps.Clone(populated)
		obj[v.key] = v.Value
		out = append(out, shape{v.key + "=" + v.Name, obj})
	}
	return out
}

// paramStruct returns param's type when it is a struct, and otherwise the param
// struct with typ's wire type that shares the most fields.
func (s *synthesizer) paramStruct(typ reflect.Type, param *reflect.StructField) reflect.Type {
	if param != nil && param.Type.Kind() == reflect.Struct {
		return param.Type
	}
	var best reflect.Type
	bestShared := -1
	for _, candidate := range s.params[discriminator(typ)] {
		if shared := sharedFields(typ, candidate); shared > bestShared {
			best, bestShared = candidate, shared
		}
	}
	return best
}

func sharedFields(response, param reflect.Type) int {
	shared := 0
	for i := 0; i < response.NumField(); i++ {
		if key := wireKey(response.Field(i)); key != "" {
			if _, ok := requestField(param, key); ok {
				shared++
			}
		}
	}
	return shared
}

// fieldShapes generates one field, adding a null shape for a param.Opt
// counterpart and dropping the zero shape for a scalar omitzero one.
func (s *synthesizer) fieldShapes(typ reflect.Type, param *reflect.StructField, depth int) []shape {
	out := s.shapes(typ, param, depth+1)
	switch {
	case param == nil:
	case isOpt(param.Type):
		out = append(out, shape{"null", nil})
	case strings.Contains(param.Tag.Get("json"), ",omitzero") && isScalar(param.Type):
		out = slices.DeleteFunc(out, func(sh shape) bool { return reflect.ValueOf(sh.Value).IsZero() })
	}
	return out
}

// requestField looks up the field of a param struct by wire key.
func requestField(paramStruct reflect.Type, key string) (*reflect.StructField, bool) {
	if paramStruct == nil {
		return nil, false
	}
	for i := 0; i < paramStruct.NumField(); i++ {
		field := paramStruct.Field(i)
		if wireKey(field) == key {
			return &field, true
		}
	}
	return nil, false
}

func isOpt(typ reflect.Type) bool {
	return strings.HasSuffix(typ.PkgPath(), "/packages/param") && strings.HasPrefix(typ.Name(), "Opt[")
}

func isScalar(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

// wireKey returns a field's JSON key, or "" for a field not on the wire such as
// the JSON metadata struct.
func wireKey(field reflect.StructField) string {
	if !field.IsExported() || field.Anonymous {
		return ""
	}
	key, _, _ := strings.Cut(field.Tag.Get("json"), ",")
	if key == "-" {
		return ""
	}
	return key
}
