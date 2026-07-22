package apijson_test

import (
	"encoding/json"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
)

// A type that looks apijson-native but opts out via UnmarshalAPIJSON,
// so its custom UnmarshalJSON must run even when nested.
type optedOut struct {
	Name string `json:"name"`
	ran  bool
	paramObject
}

func (o *optedOut) UnmarshalAPIJSON(data []byte) error { return o.UnmarshalJSON(data) }
func (o *optedOut) UnmarshalJSON(data []byte) error {
	o.ran = true
	return apijson.UnmarshalRoot(data, o)
}

type optedOutHolder struct {
	Inner optedOut `json:"inner"`
	paramObject
}

func (h *optedOutHolder) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, h) }

func TestCustomUnmarshalerRunsWhenNested(t *testing.T) {
	var h optedOutHolder
	if err := json.Unmarshal([]byte(`{"inner":{"name":"x"}}`), &h); err != nil {
		t.Fatal(err)
	}
	if !h.Inner.ran || h.Inner.Name != "x" {
		t.Fatalf("nested UnmarshalJSON was bypassed despite the UnmarshalAPIJSON opt-out: %#v", h.Inner)
	}
}
