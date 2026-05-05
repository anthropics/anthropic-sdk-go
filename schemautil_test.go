package anthropic

import (
	"encoding/json"
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
		name     string
		input    *jsonschema.Schema
		expected string // expected JSON after transform
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
			name: "no type and no anyOf clears the schema",
			input: &jsonschema.Schema{
				Description: "orphan",
			},
			// Annotation-only schemas are zeroed; a zero jsonschema.Schema
			// marshals as `true` (the "match everything" schema).
			expected: `true`,
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
			name: "invalid anyOf variants are dropped, not serialized as true",
			input: &jsonschema.Schema{
				AnyOf: []*jsonschema.Schema{
					{Type: "string"},
					{Description: "orphan: no type, no anyOf"},
				},
			},
			expected: `{"anyOf":[{"type":"string"}]}`,
		},
		{
			name: "invalid oneOf variants are dropped during conversion to anyOf",
			input: &jsonschema.Schema{
				OneOf: []*jsonschema.Schema{
					{Type: "number"},
					{Description: "orphan"},
				},
			},
			expected: `{"anyOf":[{"type":"number"}]}`,
		},
		{
			name: "unsupported allOf field renders as JSON in description",
			input: &jsonschema.Schema{
				Type: "object",
				AllOf: []*jsonschema.Schema{
					{Type: "string"},
				},
				Properties: props("x", &jsonschema.Schema{Type: "integer"}),
			},
			expected: `{"type":"object","properties":{"x":{"type":"integer"}},"additionalProperties":false,"description":"{allOf: [{\"type\":\"string\"}]}"}`,
		},
		{
			name: "unsupported not field renders as JSON in description",
			input: &jsonschema.Schema{
				Type: "string",
				Not:  &jsonschema.Schema{Type: "number"},
			},
			expected: `{"type":"string","description":"{not: {\"type\":\"number\"}}"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input == nil {
				// Just verify no panic
				transformSchema(nil)
				return
			}
			transformSchema(tt.input)
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
			name: "schema without type returns nil",
			input: map[string]any{
				"description": "A schema without type",
			},
			expected: nil,
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
