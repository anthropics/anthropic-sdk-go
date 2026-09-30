package anthropic

import (
	"encoding/json"
	"testing"

	"github.com/invopop/jsonschema"
)

// A schema whose root is a $ref still has to keep its $defs: the definitions are
// what the reference resolves against. transformSchema replaced the whole node
// with {Ref: ...}, so the submitted schema carried a dangling reference.
// anthropic-sdk-python and anthropic-sdk-typescript already handle this ordering;
// this is the same fix for the Go SDK.

func personDef() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:       "object",
		Properties: props("name", &jsonschema.Schema{Type: "string"}),
	}
}

func TestTransformSchemaKeepsDefsAlongsideRef(t *testing.T) {
	s := &jsonschema.Schema{
		Ref:         "#/$defs/Person",
		Definitions: jsonschema.Definitions{"Person": personDef()},
	}

	transformSchema(s)

	if s.Ref != "#/$defs/Person" {
		t.Fatalf("Ref = %q, want it preserved", s.Ref)
	}
	def, ok := s.Definitions["Person"]
	if !ok {
		b, _ := json.Marshal(s)
		t.Fatalf("$defs was dropped, leaving a dangling $ref: %s", b)
	}
	// The definition must have been transformed, not merely carried over.
	b, err := json.Marshal(def)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if v, ok := got["additionalProperties"]; !ok || v != false {
		t.Fatalf("$def was not transformed: %s", b)
	}
}

func TestTransformSchemaMapKeepsDefsAlongsideRef(t *testing.T) {
	in := map[string]any{
		"$ref": "#/$defs/Person",
		"$defs": map[string]any{
			"Person": map[string]any{
				"type":       "object",
				"properties": map[string]any{"name": map[string]any{"type": "string"}},
			},
		},
	}

	out := transformSchemaMap(in)

	defs, ok := out["$defs"].(map[string]any)
	if !ok {
		b, _ := json.Marshal(out)
		t.Fatalf("$defs dropped from transformed map: %s", b)
	}
	if _, ok := defs["Person"]; !ok {
		t.Fatalf("$defs lost its Person entry: %#v", defs)
	}
	if out["$ref"] != "#/$defs/Person" {
		t.Fatalf("$ref = %#v, want #/$defs/Person", out["$ref"])
	}
}

// The public constructor must not emit a schema that references something it
// removed.
func TestBetaToolInputSchemaKeepsDefsAlongsideRef(t *testing.T) {
	param := BetaToolInputSchema(map[string]any{
		"$ref": "#/$defs/Person",
		"$defs": map[string]any{
			"Person": map[string]any{
				"type":       "object",
				"properties": map[string]any{"name": map[string]any{"type": "string"}},
			},
		},
	})

	b, err := json.Marshal(param)
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if got["$ref"] != "#/$defs/Person" {
		t.Fatalf("marshalled param lost $ref: %s", b)
	}
	if _, ok := got["$defs"]; !ok {
		t.Fatalf("marshalled param has a dangling $ref with no $defs: %s", b)
	}
}

// Guard the existing contract so the fix cannot over-correct: a $ref node still
// drops unrelated siblings such as type.
func TestTransformSchemaRefStillDropsOtherSiblings(t *testing.T) {
	s := &jsonschema.Schema{
		Ref:  "#/definitions/Person",
		Type: "object",
	}
	transformSchema(s)

	b, _ := json.Marshal(s)
	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	if _, ok := got["type"]; ok {
		t.Fatalf("$ref node must not carry sibling type: %s", b)
	}
	if got["$ref"] != "#/definitions/Person" {
		t.Fatalf("$ref = %#v, want #/definitions/Person", got["$ref"])
	}
}
