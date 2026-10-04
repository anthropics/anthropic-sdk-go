package anthropic

import (
	"encoding/json"
	"reflect"
	"testing"

	validator "github.com/google/jsonschema-go/jsonschema"
	"github.com/invopop/jsonschema"
)

func conjunctionSchema(t *testing.T) map[string]any {
	t.Helper()
	var schema map[string]any
	err := json.Unmarshal([]byte(`{
        "type": "object",
        "properties": {
            "value": {
                "type": "string",
                "anyOf": [{"enum": ["a", "b"]}],
                "oneOf": [{"const": "b"}, {"const": "c"}],
                "allOf": [{"type": "string", "pattern": "^[a-z]$"}]
            }
        },
        "required": ["value"]
    }`), &schema)
	if err != nil {
		t.Fatal(err)
	}
	return schema
}

func assertConjunctionValues(t *testing.T, schema any) {
	t.Helper()
	encoded, err := json.Marshal(schema)
	if err != nil {
		t.Fatal(err)
	}
	var parsed validator.Schema
	if err := json.Unmarshal(encoded, &parsed); err != nil {
		t.Fatal(err)
	}
	resolved, err := parsed.Resolve(nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range []string{"a", "b", "c", "d"} {
		valid := resolved.Validate(map[string]any{"value": value}) == nil
		if valid != (value == "b") {
			t.Errorf("value %q: accepted=%v, want %v; schema=%s", value, valid, value == "b", encoded)
		}
	}
}

func TestSchemaSiblingUnionsPreserveConjunction(t *testing.T) {
	input := conjunctionSchema(t)
	before, _ := json.Marshal(input)
	transformed := transformSchemaMap(input)
	assertConjunctionValues(t, transformed)
	after, _ := json.Marshal(input)
	if string(before) != string(after) {
		t.Fatal("transform mutated the caller's schema")
	}
	if !reflect.DeepEqual(transformed, transformSchemaMap(transformed)) {
		t.Fatal("transform is not idempotent")
	}
}

func TestSchemaSiblingUnionsPublicHelpers(t *testing.T) {
	for _, helper := range []string{"output", "tool"} {
		t.Run(helper, func(t *testing.T) {
			input := conjunctionSchema(t)
			if helper == "output" {
				assertConjunctionValues(t, BetaJSONSchemaOutputFormat(input).Schema)
			} else {
				assertConjunctionValues(t, BetaToolInputSchema(input))
			}
		})
	}
}

func TestSchemaSiblingUnionsKeepExistingAllOfAndRecurse(t *testing.T) {
	for _, existing := range []bool{false, true} {
		t.Run(map[bool]string{false: "new conjunction", true: "existing conjunction"}[existing], func(t *testing.T) {
			variant := &jsonschema.Schema{Type: "object", Properties: props("count", &jsonschema.Schema{Type: "integer", Minimum: "2"})}
			schema := &jsonschema.Schema{
				AnyOf: []*jsonschema.Schema{{Type: "object"}},
				OneOf: []*jsonschema.Schema{variant},
			}
			prior := &jsonschema.Schema{Type: "object", Description: "existing constraint"}
			if existing {
				schema.AllOf = []*jsonschema.Schema{prior}
			}
			transformSchema(schema)
			want := 1
			if existing {
				want = 2
			}
			if len(schema.AllOf) != want {
				t.Fatalf("allOf has %d entries, want %d", len(schema.AllOf), want)
			}
			if existing && schema.AllOf[0] != prior {
				t.Fatal("existing allOf entry was replaced")
			}
			converted := schema.AllOf[want-1]
			if len(converted.AnyOf) != 1 || converted.AnyOf[0] != variant {
				t.Fatal("oneOf alternatives were not retained")
			}
			count, _ := variant.Properties.Get("count")
			if count.Minimum != "" || count.Description != "{minimum: 2}" {
				t.Errorf("nested constraint was not transformed: %+v", count)
			}
			if variant.AdditionalProperties != jsonschema.FalseSchema {
				t.Fatal("nested object was not transformed")
			}
			if len(schema.OneOf) != 0 || len(schema.AnyOf) != 1 {
				t.Fatal("sibling unions changed unexpectedly")
			}
		})
	}
}

func TestSchemaSiblingUnionsAtNestedLocations(t *testing.T) {
	for _, location := range []string{"property", "items", "definition"} {
		t.Run(location, func(t *testing.T) {
			leaf := &jsonschema.Schema{AnyOf: []*jsonschema.Schema{{Type: "string"}}, OneOf: []*jsonschema.Schema{{Const: "b"}}}
			schema := &jsonschema.Schema{Type: "object"}
			switch location {
			case "property":
				schema.Properties = props("value", leaf)
			case "items":
				schema.Properties = props("values", &jsonschema.Schema{Type: "array", Items: leaf})
			case "definition":
				schema.Definitions = jsonschema.Definitions{"value": leaf}
			}
			transformSchema(schema)
			if len(leaf.AllOf) != 1 || len(leaf.AllOf[0].AnyOf) != 1 || leaf.AllOf[0].AnyOf[0].Const != "b" {
				t.Errorf("nested oneOf was discarded: %+v", leaf)
			}
		})
	}
}
