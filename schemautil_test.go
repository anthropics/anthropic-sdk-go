package anthropic

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestTransformSchema(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformSchema(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				resultJSON, _ := json.MarshalIndent(result, "", "  ")
				expectedJSON, _ := json.MarshalIndent(tt.expected, "", "  ")
				t.Errorf("TransformSchema() mismatch:\ngot:\n%s\nwant:\n%s", resultJSON, expectedJSON)
			}
		})
	}
}
