//go:build go1.25

package json

import (
	"testing"
	"testing/synctest"
)

func TestMarshalRecursiveTypeAcrossSynctestBubble(t *testing.T) {
	type node struct {
		Next *node `json:"next,omitempty"`
	}

	// Build the recursive encoder in a bubble, leaving its indirect encoder
	// cached for later calls outside the bubble.
	synctest.Test(t, func(t *testing.T) {
		got, err := Marshal(node{})
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != `{}` {
			t.Fatalf("Marshal(node{}) = %s, want {}", got)
		}
	})

	got, err := Marshal(node{Next: &node{}})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != `{"next":{}}` {
		t.Fatalf("Marshal(nested node) = %s, want {\"next\":{}}", got)
	}
}
