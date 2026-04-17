package anthropic

import (
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/invopop/jsonschema"
	orderedmap "github.com/wk8/go-ordered-map/v2"
)

// supportedSchemaKeySet is a set version of supportedSchemaKeys for O(1) lookups.
var supportedSchemaKeySet = func() map[string]bool {
	m := make(map[string]bool, len(supportedSchemaKeys))
	for _, k := range supportedSchemaKeys {
		m[k] = true
	}
	return m
}()

// transformSchema transforms a [jsonschema.Schema] in-place to ensure it
// conforms to the Anthropic API's expectations.
//
// The transformation process:
//   - Preserves $ref references
//   - Transforms $defs recursively
//   - Handles anyOf/oneOf by converting oneOf to anyOf
//   - Ensures objects have additionalProperties: false
//   - Filters string formats to only supported ones
//   - Limits array minItems to 0 or 1
//   - Appends unsupported properties to the description
func transformSchema(s *jsonschema.Schema) {
	if s == nil {
		return
	}

	// $ref is not supported alongside other properties
	if s.Ref != "" {
		*s = jsonschema.Schema{Ref: s.Ref}
		return
	}

	// Collect and clear fields not in the supported set
	extras := extractUnsupportedFields(s)

	// Transform $defs recursively
	for _, def := range s.Definitions {
		transformSchema(def)
	}

	// Convert oneOf to anyOf before recursing so variants are transformed once.
	if len(s.OneOf) > 0 && len(s.AnyOf) == 0 {
		s.AnyOf = s.OneOf
	}
	s.OneOf = nil

	// Recurse into anyOf variants, dropping any that transformSchema zeroed out
	// as invalid — a zero jsonschema.Schema marshals as the literal JSON `true`,
	// which would otherwise leak into the variant list as a match-everything.
	if len(s.AnyOf) > 0 {
		kept := s.AnyOf[:0]
		for _, variant := range s.AnyOf {
			transformSchema(variant)
			if variant != nil && !reflect.ValueOf(*variant).IsZero() {
				kept = append(kept, variant)
			}
		}
		s.AnyOf = kept
	}

	// Bail if no type and no anyOf — schema is invalid or a boolean schema.
	// Boolean schemas (JSON true/false) carry meaning in an unexported field
	// that zeroing would clear, flipping false→true. Detect them by checking
	// whether any exported field is non-zero: if not, the schema is either
	// boolean (preserve) or truly empty (zeroing is a no-op anyway).
	if s.Type == "" && len(s.AnyOf) == 0 {
		if !hasExportedContent(s) {
			return
		}
		*s = jsonschema.Schema{}
		return
	}

	switch s.Type {
	case "object":
		if s.Properties != nil {
			for pair := s.Properties.Oldest(); pair != nil; pair = pair.Next() {
				transformSchema(pair.Value)
			}
			s.AdditionalProperties = jsonschema.FalseSchema
		} else if s.AdditionalProperties != nil {
			// Dictionary-style schema (e.g. Go map types): no fixed properties,
			// additionalProperties describes the value type. Preserve and recurse.
			transformSchema(s.AdditionalProperties)
		} else {
			s.Properties = orderedmap.New[string, *jsonschema.Schema]()
			s.AdditionalProperties = jsonschema.FalseSchema
		}

	case "string":
		if s.Format != "" && !slices.Contains(supportedStringFormats, s.Format) {
			extras["format"] = s.Format
			s.Format = ""
		}

	case "array":
		if s.Items != nil {
			transformSchema(s.Items)
		}
		if s.MinItems != nil && *s.MinItems != 0 && *s.MinItems != 1 {
			extras["minItems"] = *s.MinItems
			s.MinItems = nil
		}
	}

	// Append unsupported properties to description
	if len(extras) > 0 {
		keys := make([]string, 0, len(extras))
		for k := range extras {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		parts := make([]string, 0, len(keys))
		for _, k := range keys {
			parts = append(parts, fmt.Sprintf("%s: %s", k, formatExtraValue(extras[k])))
		}

		extraStr := "{" + strings.Join(parts, ", ") + "}"
		if s.Description != "" {
			s.Description += "\n\n" + extraStr
		} else {
			s.Description = extraStr
		}
	}
}

// transformSchemaMap transforms a JSON schema map to conform to Anthropic's API
// requirements. It delegates to [transformSchema] via a JSON round-trip.
func transformSchemaMap(jsonSchema map[string]any) map[string]any {
	if jsonSchema == nil {
		return nil
	}

	b, err := json.Marshal(jsonSchema)
	if err != nil {
		return nil
	}

	var s jsonschema.Schema
	if err := json.Unmarshal(b, &s); err != nil {
		return nil
	}

	transformSchema(&s)

	result, err := json.Marshal(&s)
	if err != nil {
		return nil
	}

	var out map[string]any
	if err := json.Unmarshal(result, &out); err != nil {
		return nil
	}
	return out
}

// formatExtraValue renders a value extracted from the schema for inclusion in
// the description. Composite values (*Schema, []*Schema, maps) are
// JSON-marshaled so they render readably instead of as Go pointer dumps.
// Scalars use %v so strings stay unquoted.
func formatExtraValue(v any) string {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return "null"
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
		if b, err := json.Marshal(v); err == nil {
			return string(b)
		}
	}
	return fmt.Sprintf("%v", v)
}

// hasExportedContent reports whether any exported field on the schema is non-zero.
// Boolean schemas (JSON true/false) store their value in an unexported field, so
// they return false here despite being meaningful.
func hasExportedContent(s *jsonschema.Schema) bool {
	v := reflect.ValueOf(s).Elem()
	t := v.Type()
	for i := range t.NumField() {
		if t.Field(i).IsExported() && !v.Field(i).IsZero() {
			return true
		}
	}
	return false
}

// extractUnsupportedFields uses reflection to find non-zero exported fields on
// the schema whose JSON key is not in the supported set. It collects their
// values into a map and zeros the fields, so they won't appear in the marshaled
// output. This avoids enumerating every field explicitly and automatically
// handles new fields added to the jsonschema library.
func extractUnsupportedFields(s *jsonschema.Schema) map[string]any {
	extras := make(map[string]any)
	v := reflect.ValueOf(s).Elem()
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		jsonTag := field.Tag.Get("json")
		name, _, _ := strings.Cut(jsonTag, ",")
		if name == "" || name == "-" {
			continue
		}

		if supportedSchemaKeySet[name] {
			continue
		}

		fv := v.Field(i)
		if fv.IsZero() {
			continue
		}

		extras[name] = fv.Interface()
		fv.SetZero()
	}

	// Handle the Extras map (tagged json:"-", inlined by Schema.MarshalJSON)
	for k, val := range s.Extras {
		if !supportedSchemaKeySet[k] {
			extras[k] = val
		}
	}
	s.Extras = nil

	return extras
}

// BetaJSONSchemaOutputFormat creates a BetaJSONOutputFormatParam from a JSON schema map.
// It transforms the schema to ensure compatibility with Anthropic's JSON schema requirements.
//
// Example:
//
//	schema := map[string]any{
//	    "type": "object",
//	    "properties": map[string]any{
//	        "name": map[string]any{"type": "string"},
//	        "age": map[string]any{"type": "integer", "minimum": 0},
//	    },
//	    "required": []string{"name"},
//	}
//	outputFormat := BetaJSONSchemaOutputFormat(schema)
//
//	msg, _ := client.Beta.Messages.New(ctx, anthropic.BetaMessageNewParams{
//	    Model: anthropic.Model("claude-sonnet-4-5"),
//	    Messages: anthropic.F([]anthropic.BetaMessageParam{...}),
//	    MaxTokens: 1024,
//	    OutputFormat: outputFormat,
//	})
func BetaJSONSchemaOutputFormat(jsonSchema map[string]any) BetaJSONOutputFormatParam {
	return BetaJSONOutputFormatParam{Schema: transformSchemaMap(jsonSchema)}
}

// BetaToolInputSchema creates a BetaToolInputSchemaParam from a JSON schema map.
// It transforms the schema to ensure compatibility with Anthropic's tool calling requirements.
func BetaToolInputSchema(jsonSchema map[string]any) BetaToolInputSchemaParam {
	return BetaToolInputSchemaParam{ExtraFields: transformSchemaMap(jsonSchema)}
}

var supportedStringFormats = []string{
	"date-time",
	"time",
	"date",
	"duration",
	"email",
	"hostname",
	"uri",
	"ipv4",
	"ipv6",
	"uuid",
}
var supportedSchemaKeys = []string{
	// Top-level schema keys
	"$ref",
	"$defs",
	"type",
	"anyOf",
	"oneOf",
	"description",
	"title",

	// Object-specific keys
	"properties",
	"additionalProperties",
	"required",

	// Array-specific keys
	"items",
	"minItems",

	// String-specific keys
	"format",
}
