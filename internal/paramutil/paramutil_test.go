package paramutil_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/paramutil"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
)

type color string

func (color) Default() color { return "blue" }

type union struct {
	OfString *string
	OfInt    *int
	param.APIUnion
}

type plainStruct struct {
	Name string
	param.APIObject
}

func TestAddrOrDefault(t *testing.T) {
	if got := *paramutil.AddrOrDefault(color("")); got != "blue" {
		t.Fatalf("empty value = %q, want the default %q", got, "blue")
	}
	if got := *paramutil.AddrOrDefault(color("red")); got != "red" {
		t.Fatalf("set value = %q, want %q", got, "red")
	}
}

func TestAddrIfPresent(t *testing.T) {
	if got := paramutil.AddrIfPresent(param.NewOpt("x")); got == nil || *got != "x" {
		t.Fatalf("included Opt = %v, want a pointer to %q", got, "x")
	}
	if got := paramutil.AddrIfPresent(param.Opt[string]{}); got != nil {
		t.Fatalf("omitted Opt = %v, want nil", got)
	}
	if got := paramutil.AddrIfPresent(param.Null[string]()); got != nil {
		t.Fatalf("null Opt = %v, want nil", got)
	}
	// A present zero value is still present.
	if got := paramutil.AddrIfPresent(param.NewOpt(0)); got == nil || *got != 0 {
		t.Fatalf("included zero Opt = %v, want a pointer to 0", got)
	}
}

func TestToOpt(t *testing.T) {
	valid := paramutil.ToOpt("x", respjson.NewField(`"x"`))
	if !valid.Valid() || valid.Value != "x" {
		t.Fatalf("valid field = %+v, want an included Opt holding %q", valid, "x")
	}

	null := paramutil.ToOpt("", respjson.NewField("null"))
	if null.Valid() || !param.IsNull(null) {
		t.Fatalf("null field = %+v, want a null Opt", null)
	}

	for name, field := range map[string]respjson.Field{
		"omitted": {},
		"invalid": respjson.NewInvalidField(`{"unexpected":true}`),
	} {
		got := paramutil.ToOpt("ignored", field)
		if got.Valid() || param.IsNull(got) || !param.IsOmitted(got) {
			t.Fatalf("%s field = %+v, want an omitted Opt", name, got)
		}
	}
}

func TestVariantFromUnion(t *testing.T) {
	s := "hello"
	n := 7

	got, err := paramutil.VariantFromUnion(reflect.ValueOf(union{OfString: &s}))
	if err != nil || got != &s {
		t.Fatalf("string variant = (%v, %v), want (%p, nil)", got, err, &s)
	}

	// A pointer to the union is accepted too.
	got, err = paramutil.VariantFromUnion(reflect.ValueOf(&union{OfInt: &n}))
	if err != nil || got != &n {
		t.Fatalf("pointer to union = (%v, %v), want (%p, nil)", got, err, &n)
	}
}

func TestVariantFromUnionErrors(t *testing.T) {
	s := "hello"
	n := 7

	cases := []struct {
		name  string
		value reflect.Value
		want  string
	}{
		{"no variants", reflect.ValueOf(union{}), "no variants"},
		{"multiple variants", reflect.ValueOf(union{OfString: &s, OfInt: &n}), "multiple variants"},
		{"not a struct", reflect.ValueOf(3), "non-struct"},
		{"struct without an embedded union", reflect.ValueOf(plainStruct{Name: "x"}), "non-union"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := paramutil.VariantFromUnion(tc.value)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want one containing %q", err, tc.want)
			}
			if got != nil {
				t.Fatalf("variant = %v, want nil alongside the error", got)
			}
		})
	}
}
