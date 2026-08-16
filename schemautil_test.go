package anthropic

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/invopop/jsonschema"
	orderedmap "github.com/pb33f/ordered-map/v2"
)

func ptr[T any](v T) *T { return &v }

func props(pairs ...any) *orderedmap.OrderedMap[string, *jsonschema.Schema] {
	m := orderedmap.New[string, *jsonschema.Schema]()
	for i := 0; i < len(pairs); i += 2 {
		m.Set(pairs[i].(string), pairs[i+1].(*jsonschema.Schema))
	}
	return m
}

// normalizeJSON round-trips JSON through any to normalize key ordering,
// since json.Marshal sorts map keys deterministically.
func normalizeJSON(s string) string {
	var v any
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return s
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func TestTransformSchema(t *testing.T) {
	tests := []struct {
		name        string
		input       *jsonschema.Schema
		expected    string // expected JSON after transform
		wantErrPath string // non-empty if error is expected
	}{
		{
			name:     "nil schema is a no-op",
			input:    nil,
			expected: "",
		},
		{
			name: "integer with unsupported constraints moves them to description",
			input: &jsonschema.Schema{
				Type:        "integer",
				Minimum:     "1",
				Maximum:     "10",
				Description: "A number",
			},
			expected: `{"type":"integer","description":"A number\n\n{maximum: 10, minimum: 1}"}`,
		},
		{
			name: "object gets additionalProperties false and recurses into properties",
			input: &jsonschema.Schema{
				Type: "object",
				Properties: props(
					"name", &jsonschema.Schema{Type: "string"},
					"age", &jsonschema.Schema{Type: "integer"},
				),
				Required: []string{"name"},
			},
			expected: `{"type":"object","properties":{"name":{"type":"string"},"age":{"type":"integer"}},"additionalProperties":false,"required":["name"]}`,
		},
		{
			name: "object without properties gets empty properties",
			input: &jsonschema.Schema{
				Type: "object",
			},
			expected: `{"type":"object","properties":{},"additionalProperties":false}`,
		},
		{
			name: "supported string format is preserved",
			input: &jsonschema.Schema{
				Type:   "string",
				Format: "date-time",
			},
			expected: `{"type":"string","format":"date-time"}`,
		},
		{
			name: "unsupported string format moves to description",
			input: &jsonschema.Schema{
				Type:   "string",
				Format: "binary",
			},
			expected: `{"type":"string","description":"{format: binary}"}`,
		},
		{
			name: "array with minItems 1 is preserved",
			input: &jsonschema.Schema{
				Type:     "array",
				Items:    &jsonschema.Schema{Type: "string"},
				MinItems: ptr(uint64(1)),
			},
			expected: `{"type":"array","items":{"type":"string"},"minItems":1}`,
		},
		{
			name: "array with minItems 5 moves to description",
			input: &jsonschema.Schema{
				Type:     "array",
				Items:    &jsonschema.Schema{Type: "string"},
				MinItems: ptr(uint64(5)),
			},
			expected: `{"type":"array","items":{"type":"string"},"description":"{minItems: 5}"}`,
		},
		{
			name: "array recurses into items",
			input: &jsonschema.Schema{
				Type: "array",
				Items: &jsonschema.Schema{
					Type: "object",
					Properties: props(
						"id", &jsonschema.Schema{Type: "integer"},
					),
				},
			},
			expected: `{"type":"array","items":{"type":"object","properties":{"id":{"type":"integer"}},"additionalProperties":false}}`,
		},
		{
			name: "$ref strips all other properties",
			input: &jsonschema.Schema{
				Ref:  "#/definitions/Person",
				Type: "object",
			},
			expected: `{"$ref":"#/definitions/Person"}`,
		},
		{
			name: "oneOf converted to anyOf",
			input: &jsonschema.Schema{
				OneOf: []*jsonschema.Schema{
					{Type: "string"},
					{Type: "number"},
				},
			},
			expected: `{"anyOf":[{"type":"string"},{"type":"number"}]}`,
		},
		{
			name: "anyOf variants are recursively transformed",
			input: &jsonschema.Schema{
				AnyOf: []*jsonschema.Schema{
					{Type: "string"},
					{Type: "object", Properties: props("x", &jsonschema.Schema{Type: "integer"})},
				},
			},
			expected: `{"anyOf":[{"type":"string"},{"type":"object","properties":{"x":{"type":"integer"}},"additionalProperties":false}]}`,
		},
		{
			name: "anyOf genuine empty schema is preserved",
			input: &jsonschema.Schema{
				AnyOf: []*jsonschema.Schema{
					{},
					{Type: "integer"},
				},
			},
			expected: `{"anyOf":[true,{"type":"integer"}]}`,
		},
		{
			name: "$defs are recursively transformed",
			input: &jsonschema.Schema{
				Type: "object",
				Definitions: jsonschema.Definitions{
					"Person": &jsonschema.Schema{
						Type: "object",
						Properties: props(
							"name", &jsonschema.Schema{Type: "string"},
						),
					},
				},
				Properties: props(
					"user", &jsonschema.Schema{Ref: "#/$defs/Person"},
				),
			},
			expected: `{"type":"object","$defs":{"Person":{"type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false}},"properties":{"user":{"$ref":"#/$defs/Person"}},"additionalProperties":false}`,
		},
		{
			name: "description only schema without type returns error",
			input: &jsonschema.Schema{
				Description: "orphan",
			},
			wantErrPath: "$",
		},
		{
			name: "minLength only schema without type returns error",
			input: &jsonschema.Schema{
				MinLength: ptr(uint64(3)),
			},
			wantErrPath: "$",
		},
		{
			name: "$schema version is stripped",
			input: &jsonschema.Schema{
				Version: "https://json-schema.org/draft/2020-12/schema",
				Type:    "object",
				Properties: props(
					"name", &jsonschema.Schema{Type: "string"},
				),
			},
			expected: `{"type":"object","properties":{"name":{"type":"string"}},"additionalProperties":false,"description":"{$schema: https://json-schema.org/draft/2020-12/schema}"}`,
		},
		{
			name: "anyOf containing nil variant reports error",
			input: &jsonschema.Schema{
				AnyOf: []*jsonschema.Schema{
					nil,
					{Type: "string"},
				},
			},
			wantErrPath: "$.anyOf[0]",
		},
		{
			name: "invalid anyOf variant reports error",
			input: &jsonschema.Schema{
				AnyOf: []*jsonschema.Schema{
					{Type: "string"},
					{Description: "orphan: no type, no anyOf"},
				},
			},
			wantErrPath: "$.anyOf[1]",
		},
		{
			name: "invalid oneOf variant reports error during conversion to anyOf",
			input: &jsonschema.Schema{
				OneOf: []*jsonschema.Schema{
					{Type: "number"},
					{Description: "orphan"},
				},
			},
			wantErrPath: "$.anyOf[1]",
		},
		{
			name: "invalid allOf variant reports error",
			input: &jsonschema.Schema{
				AllOf: []*jsonschema.Schema{
					{Type: "string"},
					{MinLength: ptr(uint64(3))},
				},
			},
			wantErrPath: "$.allOf[1]",
		},
		{
			name: "invalid object property reports error",
			input: &jsonschema.Schema{
				Type: "object",
				Properties: props(
					"payload", &jsonschema.Schema{MinLength: ptr(uint64(3))},
				),
			},
			wantErrPath: "$.properties.payload",
		},
		{
			name: "invalid array items reports error",
			input: &jsonschema.Schema{
				Type:  "array",
				Items: &jsonschema.Schema{MinLength: ptr(uint64(3))},
			},
			wantErrPath: "$.items",
		},
		{
			name: "invalid $defs definition reports error",
			input: &jsonschema.Schema{
				Type: "object",
				Definitions: jsonschema.Definitions{
					"T": &jsonschema.Schema{MinLength: ptr(uint64(3))},
				},
				Properties: props(
					"ref", &jsonschema.Schema{Ref: "#/$defs/T"},
				),
			},
			wantErrPath: "$.$defs.T",
		},
		{
			name: "invalid dictionary additionalProperties reports error",
			input: &jsonschema.Schema{
				Type:                 "object",
				AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			},
			wantErrPath: "$.additionalProperties",
		},
		{
			name: "deeply nested invalid schema reports path",
			input: &jsonschema.Schema{
				Type: "object",
				AdditionalProperties: &jsonschema.Schema{
					Type: "object",
					Properties: props(
						"x", &jsonschema.Schema{Not: &jsonschema.Schema{}},
					),
				},
			},
			wantErrPath: "$.additionalProperties.properties.x",
		},
		{
			name: "dictionary additionalProperties empty schema is preserved",
			input: &jsonschema.Schema{
				Type:                 "object",
				AdditionalProperties: &jsonschema.Schema{},
			},
			expected: `{"type":"object","additionalProperties":true}`,
		},
		{
			name: "dictionary additionalProperties typed schema is preserved",
			input: &jsonschema.Schema{
				Type:                 "object",
				AdditionalProperties: &jsonschema.Schema{Type: "string"},
			},
			expected: `{"type":"object","additionalProperties":{"type":"string"}}`,
		},
		{
			name: "dictionary additionalProperties false schema is preserved",
			input: &jsonschema.Schema{
				Type:                 "object",
				AdditionalProperties: jsonschema.FalseSchema,
			},
			expected: `{"type":"object","additionalProperties":false}`,
		},
		{
			name: "unsupported not field renders as JSON in description",
			input: &jsonschema.Schema{
				Type: "string",
				Not:  &jsonschema.Schema{Type: "number"},
			},
			expected: `{"type":"string","description":"{not: {\"type\":\"number\"}}"}`,
		},
		{
			name: "enum is preserved",
			input: &jsonschema.Schema{
				Type: "string",
				Enum: []any{"red", "green", "blue"},
			},
			expected: `{"type":"string","enum":["red","green","blue"]}`,
		},
		{
			name: "enum without type is preserved",
			input: &jsonschema.Schema{
				Enum: []any{"a", "b"},
			},
			expected: `{"enum":["a","b"]}`,
		},
		{
			name: "const is preserved",
			input: &jsonschema.Schema{
				Type:  "string",
				Const: "fixed",
			},
			expected: `{"type":"string","const":"fixed"}`,
		},
		{
			name: "const without type is preserved",
			input: &jsonschema.Schema{
				Const: 42,
			},
			expected: `{"const":42}`,
		},
		{
			name: "string pattern is preserved",
			input: &jsonschema.Schema{
				Type:    "string",
				Pattern: "^[a-z]+$",
			},
			expected: `{"type":"string","pattern":"^[a-z]+$"}`,
		},
		{
			name: "allOf variants are preserved and recursively transformed",
			input: &jsonschema.Schema{
				AllOf: []*jsonschema.Schema{
					{Type: "object", Properties: props("a", &jsonschema.Schema{Type: "string"})},
					{Type: "object", Properties: props("b", &jsonschema.Schema{Type: "integer"}), Required: []string{"b"}},
				},
			},
			expected: `{"allOf":[{"type":"object","properties":{"a":{"type":"string"}},"additionalProperties":false},{"type":"object","properties":{"b":{"type":"integer"}},"additionalProperties":false,"required":["b"]}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input == nil {
				// Just verify no panic
				if err := transformSchema(nil); err != nil {
					t.Fatalf("unexpected error for nil schema: %v", err)
				}
				return
			}
			err := transformSchema(tt.input)
			if tt.wantErrPath != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErrPath)
				}
				if !strings.Contains(err.Error(), tt.wantErrPath) {
					t.Fatalf("expected error to contain %q, got %v", tt.wantErrPath, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			gotBytes, _ := json.Marshal(tt.input)
			got := normalizeJSON(string(gotBytes))
			want := normalizeJSON(tt.expected)
			if got != want {
				// Pretty print for readable diffs
				var gotAny, wantAny any
				json.Unmarshal([]byte(got), &gotAny)
				json.Unmarshal([]byte(tt.expected), &wantAny)
				gotPretty, _ := json.MarshalIndent(gotAny, "", "  ")
				wantPretty, _ := json.MarshalIndent(wantAny, "", "  ")
				t.Errorf("transformSchema() mismatch:\ngot:\n%s\nwant:\n%s", gotPretty, wantPretty)
			}
		})
	}
}

func TestTransformSchemaFromReflector(t *testing.T) {
	type OrderItem struct {
		Name     string  `json:"name"`
		Quantity int     `json:"quantity" jsonschema:"minimum=1"`
		Price    float64 `json:"price"`
	}
	type Order struct {
		Items    []OrderItem `json:"items" jsonschema:"description=List of items"`
		Total    float64     `json:"total"`
		Currency string      `json:"currency" jsonschema:"enum=USD,enum=EUR"`
	}

	reflector := jsonschema.Reflector{DoNotReference: true}
	schema := reflector.Reflect(&Order{})
	transformSchema(schema)

	result, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(result, &m); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Top-level object must have additionalProperties: false
	if ap, ok := m["additionalProperties"]; !ok || ap != false {
		t.Errorf("expected top-level additionalProperties=false, got %v", m["additionalProperties"])
	}

	// $schema should be stripped (moved to description)
	if _, ok := m["$schema"]; ok {
		t.Error("expected $schema to be stripped from output")
	}

	// Properties should exist
	props, ok := m["properties"].(map[string]any)
	if !ok {
		t.Fatal("expected properties map")
	}
	for _, key := range []string{"items", "total", "currency"} {
		if _, ok := props[key]; !ok {
			t.Errorf("expected property %q", key)
		}
	}

	// Nested items schema should also have additionalProperties: false
	itemsProp, ok := props["items"].(map[string]any)
	if !ok {
		t.Fatal("expected items property to be a map")
	}
	itemsSchema, ok := itemsProp["items"].(map[string]any)
	if !ok {
		t.Fatal("expected items.items (array element schema) to be a map")
	}
	if ap, ok := itemsSchema["additionalProperties"]; !ok || ap != false {
		t.Errorf("expected nested additionalProperties=false, got %v", itemsSchema["additionalProperties"])
	}

	// Unsupported constraint (minimum=1 on quantity) should be in description
	quantitySchema, ok := itemsSchema["properties"].(map[string]any)["quantity"].(map[string]any)
	if !ok {
		t.Fatal("expected quantity property")
	}
	desc, _ := quantitySchema["description"].(string)
	if desc == "" {
		t.Error("expected quantity to have a description with the minimum constraint")
	}
}

// Keywords like maxLength and maxItems are pointer-typed on jsonschema.Schema, so
// the description must render the pointed-to value rather than the address.
func TestTransformSchemaFromReflectorPointerKeywords(t *testing.T) {
	type Input struct {
		Code string   `json:"code" jsonschema:"minLength=2,maxLength=8"`
		Tags []string `json:"tags" jsonschema:"maxItems=3"`
	}

	reflector := jsonschema.Reflector{DoNotReference: true}
	schema := reflector.Reflect(&Input{})
	transformSchema(schema)

	tests := []struct {
		property string
		want     string
	}{
		{"code", "{maxLength: 8, minLength: 2}"},
		{"tags", "{maxItems: 3}"},
	}
	for _, tt := range tests {
		t.Run(tt.property+" keywords render as values in the description", func(t *testing.T) {
			prop, ok := schema.Properties.Get(tt.property)
			if !ok {
				t.Fatalf("expected property %q", tt.property)
			}
			if strings.Contains(prop.Description, "0x") {
				t.Errorf("description contains a pointer address: %q", prop.Description)
			}
			if prop.Description != tt.want {
				t.Errorf("description = %q, want %q", prop.Description, tt.want)
			}
		})
	}
}

func TestFormatExtraValue(t *testing.T) {
	tests := []struct {
		name string
		v    any
		want string
	}{
		{"pointer to integer renders the value", ptr(uint64(8)), "8"},
		{"nil pointer renders as null", (*uint64)(nil), "null"},
		{"untyped nil renders as null", nil, "null"},
		{"plain integer", uint64(3), "3"},
		{"string stays unquoted", "date", "date"},
		{"pointer to schema is JSON-marshaled", &jsonschema.Schema{Type: "string"}, `{"type":"string"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatExtraValue(tt.v); got != tt.want {
				t.Errorf("formatExtraValue(%#v) = %q, want %q", tt.v, got, tt.want)
			}
		})
	}
}

func TestTransformSchemaMap(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected map[string]any
	}{
		{
			name: "basic integer with unsupported properties",
			input: map[string]any{
				"type":        "integer",
				"minimum":     1,
				"maximum":     10,
				"description": "A number",
			},
			expected: map[string]any{
				"type":        "integer",
				"description": "A number\n\n{maximum: 10, minimum: 1}",
			},
		},
		{
			name: "object with properties",
			input: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{
						"type": "string",
					},
					"age": map[string]any{
						"type": "integer",
					},
				},
				"required":             []string{"name"},
				"additionalProperties": true,
			},
			expected: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"name": map[string]any{
						"type": "string",
					},
					"age": map[string]any{
						"type": "integer",
					},
				},
				"required":             []string{"name"},
				"additionalProperties": false,
			},
		},
		{
			name: "string with supported format",
			input: map[string]any{
				"type":   "string",
				"format": "date-time",
			},
			expected: map[string]any{
				"type":   "string",
				"format": "date-time",
			},
		},
		{
			name: "string with unsupported format",
			input: map[string]any{
				"type":   "string",
				"format": "binary",
			},
			expected: map[string]any{
				"type":        "string",
				"description": "{format: binary}",
			},
		},
		{
			name: "array with minItems 1",
			input: map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
				"minItems": 1,
			},
			expected: map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
				"minItems": 1,
			},
		},
		{
			name: "array with minItems 5",
			input: map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
				"minItems": 5,
			},
			expected: map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "string",
				},
				"description": "{minItems: 5}",
			},
		},
		{
			name: "schema with $ref",
			input: map[string]any{
				"$ref": "#/definitions/Person",
				"type": "object",
			},
			expected: map[string]any{
				"$ref": "#/definitions/Person",
			},
		},
		{
			name: "schema with anyOf",
			input: map[string]any{
				"anyOf": []any{
					map[string]any{
						"type": "string",
					},
					map[string]any{
						"type": "number",
					},
				},
			},
			expected: map[string]any{
				"anyOf": []any{
					map[string]any{
						"type": "string",
					},
					map[string]any{
						"type": "number",
					},
				},
			},
		},
		{
			name: "schema with oneOf converted to anyOf",
			input: map[string]any{
				"oneOf": []any{
					map[string]any{
						"type": "string",
					},
					map[string]any{
						"type": "number",
					},
				},
			},
			expected: map[string]any{
				"anyOf": []any{
					map[string]any{
						"type": "string",
					},
					map[string]any{
						"type": "number",
					},
				},
			},
		},
		{
			name: "schema with $defs",
			input: map[string]any{
				"type": "object",
				"$defs": map[string]any{
					"Person": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"name": map[string]any{
								"type": "string",
							},
						},
					},
				},
				"properties": map[string]any{
					"user": map[string]any{
						"$ref": "#/$defs/Person",
					},
				},
			},
			expected: map[string]any{
				"type": "object",
				"$defs": map[string]any{
					"Person": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"name": map[string]any{
								"type": "string",
							},
						},
						"additionalProperties": false,
					},
				},
				"properties": map[string]any{
					"user": map[string]any{
						"$ref": "#/$defs/Person",
					},
				},
				"additionalProperties": false,
			},
		},
		{
			name:     "nil schema returns nil",
			input:    nil,
			expected: nil,
		},
		{
			name: "schema without type preserved via fallback",
			input: map[string]any{
				"description": "A schema without type",
			},
			expected: map[string]any{
				"description": "A schema without type",
			},
		},
		{
			name: "root minLength without type preserved via fallback",
			input: map[string]any{
				"minLength": 3,
			},
			expected: map[string]any{
				"minLength": 3,
			},
		},
		{
			name: "invalid additionalProperties preserved via fallback",
			input: map[string]any{
				"type": "object",
				"additionalProperties": map[string]any{
					"not": map[string]any{},
				},
			},
			expected: map[string]any{
				"type": "object",
				"additionalProperties": map[string]any{
					"not": map[string]any{},
				},
			},
		},
		{
			name: "dictionary additionalProperties empty schema preserved",
			input: map[string]any{
				"type":                 "object",
				"additionalProperties": map[string]any{},
			},
			expected: map[string]any{
				"type":                 "object",
				"additionalProperties": true,
			},
		},
		{
			name: "dictionary additionalProperties typed schema preserved",
			input: map[string]any{
				"type": "object",
				"additionalProperties": map[string]any{
					"type": "string",
				},
			},
			expected: map[string]any{
				"type": "object",
				"additionalProperties": map[string]any{
					"type": "string",
				},
			},
		},
		{
			name: "array items false remains false",
			input: map[string]any{
				"type":  "array",
				"items": false,
			},
			expected: map[string]any{
				"type":  "array",
				"items": false,
			},
		},
		{
			name: "object property false remains false",
			input: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"x": false,
				},
			},
			expected: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"x": false,
				},
				"additionalProperties": false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformSchemaMap(tt.input)
			// Compare via JSON to avoid Go type mismatches (e.g. []string vs []any,
			// int vs float64) that arise from the JSON round-trip in transformSchemaMap.
			resultJSON, _ := json.Marshal(result)
			expectedJSON, _ := json.Marshal(tt.expected)
			if string(resultJSON) != string(expectedJSON) {
				resultPretty, _ := json.MarshalIndent(result, "", "  ")
				expectedPretty, _ := json.MarshalIndent(tt.expected, "", "  ")
				t.Errorf("transformSchemaMap() mismatch:\ngot:\n%s\nwant:\n%s", resultPretty, expectedPretty)
			}
		})
	}
}

func TestTransformSchemaMapFallbackImmutability(t *testing.T) {
	orig := map[string]any{
		"type": "object",
		"additionalProperties": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"x": map[string]any{
					"not": map[string]any{},
				},
			},
		},
	}
	expectedOrig := map[string]any{
		"type": "object",
		"additionalProperties": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"x": map[string]any{
					"not": map[string]any{},
				},
			},
		},
	}

	fallback := transformSchemaMap(orig)
	if fallback == nil {
		t.Fatal("expected non-nil fallback map")
	}

	// Mutate the returned fallback
	fallback["type"] = "mutated_type"
	ap := fallback["additionalProperties"].(map[string]any)
	ap["type"] = "mutated_ap_type"
	props := ap["properties"].(map[string]any)
	props["x"] = "mutated_x"

	// Verify caller input is untouched
	if !reflect.DeepEqual(orig, expectedOrig) {
		t.Errorf("original map was mutated:\ngot:  %+v\nwant: %+v", orig, expectedOrig)
	}
}

func TestDefsDeterministicTraversal(t *testing.T) {
	// Create schema with multiple invalid $defs
	s := &jsonschema.Schema{
		Type: "object",
		Definitions: jsonschema.Definitions{
			"Z": &jsonschema.Schema{MinLength: ptr(uint64(1))},
			"A": &jsonschema.Schema{MinLength: ptr(uint64(1))},
			"M": &jsonschema.Schema{MinLength: ptr(uint64(1))},
			"B": &jsonschema.Schema{MinLength: ptr(uint64(1))},
		},
		Properties: props("ref", &jsonschema.Schema{Ref: "#/$defs/A"}),
	}

	// Over multiple runs, the error must always deterministically point to A
	for i := 0; i < 50; i++ {
		// Make a shallow copy of s and its Definitions map to test traversal order
		copyDefs := make(jsonschema.Definitions, len(s.Definitions))
		for k, v := range s.Definitions {
			copyDef := *v
			copyDefs[k] = &copyDef
		}
		testSchema := *s
		testSchema.Definitions = copyDefs

		err := transformSchema(&testSchema)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "$.$defs.A") {
			t.Fatalf("iteration %d: expected error path $.$defs.A, got: %v", i, err)
		}
	}
}

func TestStructuredSchemaHelpersPreserveFalseSubschemas(t *testing.T) {
	tests := []struct {
		name     string
		got      any
		expected string
	}{
		{
			name:     "BetaJSONSchemaOutputFormat preserves false child schemas",
			got:      BetaJSONSchemaOutputFormat(map[string]any{"type": "array", "items": false}).Schema,
			expected: `{"type":"array","items":false}`,
		},
		{
			name:     "BetaToolInputSchema preserves false child schemas",
			got:      BetaToolInputSchema(map[string]any{"type": "object", "properties": map[string]any{"x": false}}).ExtraFields,
			expected: `{"type":"object","properties":{"x":false},"additionalProperties":false}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := json.Marshal(tt.got)
			if normalizeJSON(string(got)) != normalizeJSON(tt.expected) {
				t.Errorf("schema helper mismatch:\ngot:\n%s\nwant:\n%s", got, tt.expected)
			}
		})
	}
}

func TestStructuredSchemaHelpersWithInvalidSchema(t *testing.T) {
	invalidSchema := map[string]any{
		"type": "object",
		"additionalProperties": map[string]any{
			"not": map[string]any{},
		},
	}

	tests := []struct {
		name string
		got  any
	}{
		{
			name: "BetaJSONSchemaOutputFormat preserves invalid schema",
			got:  BetaJSONSchemaOutputFormat(invalidSchema).Schema,
		},
		{
			name: "BetaToolInputSchema preserves invalid schema",
			got:  BetaToolInputSchema(invalidSchema).ExtraFields,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBytes, _ := json.Marshal(tt.got)
			wantBytes, _ := json.Marshal(invalidSchema)
			if string(gotBytes) != string(wantBytes) {
				t.Errorf("schema helper modified invalid schema:\ngot:  %s\nwant: %s", gotBytes, wantBytes)
			}
		})
	}
}
