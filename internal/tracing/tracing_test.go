package tracing

import (
	"errors"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func TestHandleRecoveredReportsToTheErrorHandler(t *testing.T) {
	var handled []error
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) { handled = append(handled, err) }))

	cause := errors.New("boom")
	HandleRecovered(nil)
	HandleRecovered(cause)
	HandleRecovered("not an error")
	if len(handled) != 2 || !errors.Is(handled[0], cause) || !strings.Contains(handled[1].Error(), "not an error") {
		t.Errorf("the error handler got %v, want the recovered error and the recovered value", handled)
	}
}

func TestDefinedQuery(t *testing.T) {
	defined := []string{"created_at[gte]", "ids", "limit", "page"}
	cases := map[string]struct{ rawQuery, want string }{
		"no query string":                         {"", ""},
		"nothing that the endpoint defines":       {"beta=true&trace=on", ""},
		"the order and the bytes as sent":         {"page=a%2Fb+c&beta=true&limit=5&ids%5B%5D=b&ids[]=a", "page=a%2Fb+c&limit=5&ids%5B%5D=b&ids[]=a"},
		"a name with brackets of its own":         {"created_at%5Bgte%5D=1&created_at%5Blte%5D=2&created_at=3", "created_at%5Bgte%5D=1"},
		"brackets the endpoint does not list":     {"ids%5B0%5D%5Bx%5D=a&idsx%5B%5D=b&%5Bids%5D=c", ""},
		"a key under a name that takes one value": {"limit%5Bx%5D=1&limit=5&ids%5B%5D=a", "limit=5&ids%5B%5D=a"},
		"a name that is encoded":                  {"%6Cimit=5&li%mit=6", "%6Cimit=5"},
		"a name without a value, and empty pairs": {"&limit&&page=", "limit&page="},
		"a name is not a prefix of another":       {"limits=1&pages=2", ""},
	}
	for name, test := range cases {
		if got := definedQuery(test.rawQuery, defined); got != test.want {
			t.Errorf("%s: definedQuery(%q) = %q, want %q", name, test.rawQuery, got, test.want)
		}
	}
	if got := definedQuery("limit=1", nil); got != "" {
		t.Errorf("an endpoint that defines no parameter keeps %q", got)
	}
}

// A parameter that is content wherever it is defined makes the query content
// for an endpoint that names no content parameter of its own.
func TestURLAttributesWithoutEndpointContentQueryParams(t *testing.T) {
	for _, name := range contentQueryParams {
		u, err := url.Parse("https://api.example.com/v1/things?limit=5&" + name + "=a")
		if err != nil {
			t.Fatal(err)
		}
		for _, kv := range URLAttributes(u, []string{name, "limit"}, nil, false) {
			if kv.Key == URLQueryKey || kv.Key == URLFullKey {
				t.Errorf("%s: the span has %s=%q, want none", name, kv.Key, kv.Value.AsString())
			}
		}
		if got := len(URLAttributes(u, []string{name, "limit"}, nil, true)); got != 4 {
			t.Errorf("%s: %d attributes with content on, want 4", name, got)
		}
	}
}

func TestURLAttributes(t *testing.T) {
	// Which parameters have content for a value depends on the API, so the test sets them.
	defer func(params []string) { contentQueryParams = params }(contentQueryParams)
	contentQueryParams = []string{"email", "path_prefix"}
	cases := map[string]struct {
		// defined are the parameters of the endpoint, email and limit unless the case
		// has its own, and endpointContent the ones that are content on it only.
		defined, endpointContent []string
		content                  bool
		want                     map[attribute.Key]string
	}{
		"https://user:secret@api.example.com/v1/files?beta=true&limit=5#top": {
			want: map[attribute.Key]string{URLSchemeKey: "https", URLPathKey: "/v1/files", URLQueryKey: "limit=5", URLFullKey: "https://api.example.com/v1/files?limit=5"},
		},
		"http://[::1]:8080/model/us.claude%3A0/invoke?beta=true": {
			want: map[attribute.Key]string{URLSchemeKey: "http", URLPathKey: "/model/us.claude%3A0/invoke", URLFullKey: "http://[::1]:8080/model/us.claude%3A0/invoke"},
		},
		// A parameter whose value is content takes the query and the full URL with it,
		// though the endpoint names none of its own.
		"https://api.example.com/v1/users?limit=5&em%61il=a%40b.c": {
			want: map[attribute.Key]string{URLSchemeKey: "https", URLPathKey: "/v1/users"},
		},
		"https://api.example.com/v1/users?email=a%40b.c&limit=5": {
			content: true,
			want:    map[attribute.Key]string{URLSchemeKey: "https", URLPathKey: "/v1/users", URLQueryKey: "email=a%40b.c&limit=5", URLFullKey: "https://api.example.com/v1/users?email=a%40b.c&limit=5"},
		},
		// One that the endpoint does not define is not in the query at all.
		"https://api.example.com/v1/files?limit=5&path_prefix=%2Fa": {
			want: map[attribute.Key]string{URLSchemeKey: "https", URLPathKey: "/v1/files", URLQueryKey: "limit=5", URLFullKey: "https://api.example.com/v1/files?limit=5"},
		},
		// A parameter that is content on some endpoints only is so for the endpoints
		// that say so, and for no other.
		"https://api.example.com/v1/skills?filter%5B%5D=name%3AQuarterly&limit=5": {
			defined: []string{"filter", "limit"}, endpointContent: []string{"filter"},
			want: map[attribute.Key]string{URLSchemeKey: "https", URLPathKey: "/v1/skills"},
		},
		"https://api.example.com/v1/skills?email=a%40b.c&limit=5": {
			defined: []string{"email", "filter", "limit"}, endpointContent: []string{"filter"},
			want: map[attribute.Key]string{URLSchemeKey: "https", URLPathKey: "/v1/skills"},
		},
		"https://api.example.com/v1/users?filter%5B%5D=user_id%3Au_1&limit=5": {
			defined: []string{"filter", "limit"},
			want:    map[attribute.Key]string{URLSchemeKey: "https", URLPathKey: "/v1/users", URLQueryKey: "filter%5B%5D=user_id%3Au_1&limit=5", URLFullKey: "https://api.example.com/v1/users?filter%5B%5D=user_id%3Au_1&limit=5"},
		},
	}
	for raw, test := range cases {
		u, err := url.Parse(raw)
		if err != nil {
			t.Fatal(err)
		}
		if test.defined == nil {
			test.defined = []string{"email", "limit"}
		}
		got := map[attribute.Key]string{}
		for _, kv := range URLAttributes(u, test.defined, test.endpointContent, test.content) {
			got[kv.Key] = kv.Value.AsString()
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("URLAttributes(%s) = %v, want %v", raw, got, test.want)
		}
	}
}
