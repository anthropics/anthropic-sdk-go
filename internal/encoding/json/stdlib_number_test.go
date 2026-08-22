package json_test

import (
	"bytes"
	stdjson "encoding/json"
	"testing"

	shimjson "github.com/anthropics/anthropic-sdk-go/internal/encoding/json"
)

func TestMarshalStdlibJSONNumber(t *testing.T) {
	payload := map[string]any{
		"minimum": stdjson.Number("1"),
		"maximum": stdjson.Number("9007199254740993"),
	}

	got, err := shimjson.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	dec := stdjson.NewDecoder(bytes.NewReader(got))
	dec.UseNumber()
	var decoded map[string]any
	if err := dec.Decode(&decoded); err != nil {
		t.Fatalf("decode %s: %v", got, err)
	}

	if decoded["minimum"] != stdjson.Number("1") {
		t.Errorf("minimum: got %#v, want json.Number(1)", decoded["minimum"])
	}
	if decoded["maximum"] != stdjson.Number("9007199254740993") {
		t.Errorf("maximum: got %#v, want json.Number(9007199254740993) in %s", decoded["maximum"], got)
	}
}
