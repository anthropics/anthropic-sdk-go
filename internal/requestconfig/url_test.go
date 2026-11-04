package requestconfig

import (
	"testing"
)

func TestValidateRequestPath(t *testing.T) {
	cases := map[string]struct {
		u       string
		refused bool
	}{
		"no dot segment":                 {"v1/files/abc", false},
		"..":                             {"v1/files/..", true},
		"%2e in the middle of the path":  {"v1/files/%2e/content", true},
		"..%3B then more":                {"v1/files/..%3Bx", true},
		".%3B then more":                 {"v1/files/.%3Bx", true},
		"..%3B at the end":               {"v1/files/..%3B", true},
		"..; with a raw ;":               {"v1/files/..;x", true},
		"encoded dots and lowercase %3b": {"v1/files/%2e%2E%3bx", true},
		"%3B first":                      {"v1/files/%3Bx", true},
		"%3B alone":                      {"v1/files/%3B", true},
		"a raw ; first":                  {"v1/files/;x", true},
		"..;x as the host of //host/x":   {"//..;x/y", true},
		"..;x as the user of //u@host/x": {"//..;x@host/y", true},
		"a plain user in //u@host/x":     {"//a@host/y", false},
		"%3B after other characters":     {"v1/files/a%3Bb", false},
		"%3B last":                       {"v1/files/a%3B", false},
		"three dots then %3B":            {"v1/files/...%3Bx", false},
		"dots after other characters":    {"v1/files/a..%3Bx", false},
		"dots whose % is itself escaped": {"v1/files/%252e%252e%3Bx", false},
	}
	for name, test := range cases {
		err := ValidateRequestPath(test.u)
		if refused := err != nil; refused != test.refused {
			t.Errorf("%s: ValidateRequestPath(%q) = %v, want refused %t", name, test.u, err, test.refused)
		}
	}
}
