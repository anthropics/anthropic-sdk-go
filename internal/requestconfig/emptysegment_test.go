package requestconfig

import "testing"

func TestIsEmptySegment(t *testing.T) {
	for segment, want := range map[string]bool{
		"":     true,
		";":    true,
		";x":   true,
		";a=b": true,
		";..":  true,
		"a":    false,
		"a;":   false,
		"a;b":  false,
		"..":   false,
		"..;x": false,
		":x;y": false,
		"%3Bx": false,
	} {
		if got := IsEmptySegment(segment); got != want {
			t.Errorf("IsEmptySegment(%q) = %t, want %t", segment, got, want)
		}
	}
}
