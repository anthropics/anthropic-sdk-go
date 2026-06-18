package stainlessheader

import (
	"net/http"
	"testing"
)

func TestAppendHeaderValue(t *testing.T) {
	cases := []struct {
		name string
		h    http.Header
		add  Value
		want string
	}{
		{
			name: "empty",
			h:    http.Header{},
			add:  BetaToolRunner,
			want: "BetaToolRunner",
		},
		{
			name: "appends to existing",
			h:    http.Header{"X-Stainless-Helper": []string{"mcp_tool"}},
			add:  BetaToolRunner,
			want: "mcp_tool, BetaToolRunner",
		},
		{
			name: "dedups",
			h:    http.Header{"X-Stainless-Helper": []string{"compaction"}},
			add:  Compaction,
			want: "compaction",
		},
		{
			name: "collapses multiple lines",
			h:    http.Header{"X-Stainless-Helper": []string{"a", "b, c"}},
			add:  Compaction,
			want: "a, b, c, compaction",
		},
		{
			name: "trims whitespace",
			h:    http.Header{"X-Stainless-Helper": []string{" a ,  b "}},
			add:  Compaction,
			want: "a, b, compaction",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			AppendHeaderValue(tc.h, tc.add)
			if got := tc.h.Get(Header); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
			if vals := tc.h.Values(Header); len(vals) != 1 {
				t.Errorf("expected exactly one header line, got %d: %v", len(vals), vals)
			}
		})
	}
}
