package tracing

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

func attributeValues(attrs []attribute.KeyValue) map[attribute.Key]any {
	values := map[attribute.Key]any{}
	for _, kv := range attrs {
		values[kv.Key] = kv.Value.AsInterface()
	}
	return values
}

func TestBodyAttributes(t *testing.T) {
	rows := []bodyAttribute{
		{key: "holder", property: "holder", typ: stringType},
		{key: "limit", property: "limits", path: "daily", typ: integerType},
		{key: "rate", property: "rate", typ: numberType},
		{key: "virtual", property: "virtual", typ: booleanType},
		{key: "tags", property: "tags", typ: stringType, list: true},
		{key: "amount", property: "funding", path: "amount"},
	}
	cases := map[string]struct {
		body string
		want map[attribute.Key]any
	}{
		"a value of each type": {
			`{"holder":"ada","limits":{"daily":5},"rate":0.5,"virtual":true,"tags":["a","b"],"funding":{"amount":7}}`,
			map[attribute.Key]any{"holder": "ada", "limit": int64(5), "rate": 0.5, "virtual": true, "tags": []string{"a", "b"}, "amount": int64(7)},
		},
		"a property that the body lacks sets nothing": {
			`{"holder":"ada","limits":{}}`,
			map[attribute.Key]any{"holder": "ada"},
		},
		"a value of another type sets nothing": {
			`{"holder":1,"limits":{"daily":"5"},"rate":"0.5","virtual":"true","tags":["a",1],"funding":{"amount":null}}`,
			map[attribute.Key]any{},
		},
		"a property with more than one type has the type of its value": {
			`{"funding":{"amount":"seven"}}`,
			map[attribute.Key]any{"amount": "seven"},
		},
		"a number with a fraction is not an integer": {
			`{"funding":{"amount":7.5}}`,
			map[attribute.Key]any{"amount": 7.5},
		},
		"an empty string sets nothing": {
			`{"holder":"","limits":{"daily":0},"funding":{"amount":""}}`,
			map[attribute.Key]any{"limit": int64(0)},
		},
		"an empty array is an empty list": {
			`{"tags":[]}`,
			map[attribute.Key]any{"tags": []string{}},
		},
		"a body that is not an object sets nothing": {
			`["holder"]`,
			map[attribute.Key]any{},
		},
	}
	for name, test := range cases {
		if got := attributeValues(bodyAttributes(rows, gjson.Parse(test.body))); !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: got %v, want %v", name, got, test.want)
		}
	}
}

// The attributes of a span live until it is exported: a string among them that
// shared the memory of the body would keep the whole body for as long.
func TestBodyAttributesCopyTheirStrings(t *testing.T) {
	rows := []bodyAttribute{
		{key: "holder", property: "holder", typ: stringType},
		{key: "tags", property: "tags", typ: stringType, list: true},
		{key: "amount", property: "funding", path: "amount"},
	}
	body := strings.Repeat(" ", 64) + `{"holder":"ada","tags":["a","b"],"funding":{"amount":"seven"}}`
	within := func(s string) bool {
		start, at := uintptr(unsafe.Pointer(unsafe.StringData(body))), uintptr(unsafe.Pointer(unsafe.StringData(s)))
		return at >= start && at < start+uintptr(len(body))
	}
	attrs := bodyAttributes(rows, gjson.Parse(body))
	if len(attrs) != len(rows) {
		t.Fatalf("got %d attributes, want %d", len(attrs), len(rows))
	}
	for _, kv := range attrs {
		strs := kv.Value.AsStringSlice()
		if kv.Value.Type() == attribute.STRING {
			strs = []string{kv.Value.AsString()}
		}
		for _, s := range strs {
			if within(s) {
				t.Errorf("%s: the string %q shares the memory of the body", kv.Key, s)
			}
		}
	}
}

func TestHeaderAttributes(t *testing.T) {
	rows := []headerAttribute{
		{key: "request", header: "request-id"},
		{key: "betas", header: "acme-beta", list: true},
	}
	cases := map[string]struct {
		header http.Header
		want   map[attribute.Key]any
	}{
		"one string for the header, and one for each value in a list": {
			http.Header{"Request-Id": {"req_1"}, "Acme-Beta": {"a, b", " c ,"}},
			map[attribute.Key]any{"request": "req_1", "betas": []string{"a", "b", "c"}},
		},
		"a header that is absent or empty sets nothing": {
			http.Header{"Request-Id": {""}, "Acme-Beta": {" , "}},
			map[attribute.Key]any{},
		},
	}
	for name, test := range cases {
		if got := attributeValues(headerAttributes(rows, test.header)); !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: got %v, want %v", name, got, test.want)
		}
	}
}

func TestSpanRequestAttributes(t *testing.T) {
	// Which spans have a table depends on the API, so the test sets one.
	defer func(tables map[string]spanTable) { spanAttributes = tables }(spanAttributes)
	spanAttributes = map[string]spanTable{"acme.cards.create": {
		constants:      []attribute.KeyValue{attribute.String("operation", "create")},
		requestHeaders: []headerAttribute{{key: "betas", header: "acme-beta", list: true}},
		stream:         "stream",
	}}
	header := http.Header{"Acme-Beta": {"a"}}
	cases := map[string]struct {
		span    string
		streams bool
		want    map[attribute.Key]any
	}{
		"a request that streams": {
			"acme.cards.create", true,
			map[attribute.Key]any{"operation": "create", "betas": []string{"a"}, "stream": true},
		},
		"a request that does not stream": {
			"acme.cards.create", false,
			map[attribute.Key]any{"operation": "create", "betas": []string{"a"}},
		},
		"a span without a table": {
			"acme.cards.list", true,
			map[attribute.Key]any{},
		},
	}
	for name, test := range cases {
		if got := attributeValues(SpanRequestAttributes(test.span, header, test.streams)); !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s: got %v, want %v", name, got, test.want)
		}
	}
}

// keptSpan is a span that keeps the attributes that it is given.
type keptSpan struct {
	noop.Span
	records bool
	attrs   []attribute.KeyValue
}

func (span *keptSpan) IsRecording() bool { return span.records }

func (span *keptSpan) SetAttributes(attrs ...attribute.KeyValue) {
	span.attrs = append(span.attrs, attrs...)
}

func TestSetSpanResponseAttributes(t *testing.T) {
	defer func(tables map[string]spanTable) { spanAttributes = tables }(spanAttributes)
	spanAttributes = map[string]spanTable{"acme.cards.create": {
		response: []bodyAttribute{
			{key: "id", property: "id", typ: stringType},
			{key: "available", property: "balance", path: "available", typ: numberType},
		},
	}}
	if got, want := SpanResponsePaths("acme.cards.create"), []string{"id", "balance.available"}; !reflect.DeepEqual(got, want) {
		t.Errorf("the paths of the response: got %v, want %v", got, want)
	}
	body := `{"id":"card_1"}`
	cases := map[string]struct {
		span    string
		records bool
		want    map[attribute.Key]any
	}{
		"a span that records":         {"acme.cards.create", true, map[attribute.Key]any{"id": "card_1"}},
		"a span that does not record": {"acme.cards.create", false, map[attribute.Key]any{}},
		"a span without a table":      {"acme.cards.list", true, map[attribute.Key]any{}},
	}
	for name, test := range cases {
		ofString, ofBytes := &keptSpan{records: test.records}, &keptSpan{records: test.records}
		SetSpanResponseAttributes(&APICall{Name: test.span, Span: ofString}, body)
		SetSpanResponseAttributes(&APICall{Name: test.span, Span: ofBytes}, []byte(body))
		for kind, span := range map[string]*keptSpan{"string": ofString, "bytes": ofBytes} {
			if got := attributeValues(span.attrs); !reflect.DeepEqual(got, test.want) {
				t.Errorf("%s, from %s: got %v, want %v", name, kind, got, test.want)
			}
		}
	}

	// The body of a response can be large: a span that reads none of it does not copy it.
	call, large := &APICall{Name: "acme.cards.list", Span: &keptSpan{records: true}}, make([]byte, 1<<20)
	if allocs := testing.AllocsPerRun(10, func() { SetSpanResponseAttributes(call, large) }); allocs != 0 {
		t.Errorf("a span without a table: got %v allocations, want 0", allocs)
	}
}

// pickySpan is a span that panics when it is given the attribute refused.
type pickySpan struct {
	keptSpan
	refused attribute.Key
	status  codes.Code
	ended   bool
}

func (span *pickySpan) SetAttributes(attrs ...attribute.KeyValue) {
	for _, kv := range attrs {
		if kv.Key == span.refused {
			panic("refused")
		}
	}
	span.keptSpan.SetAttributes(attrs...)
}

func (span *pickySpan) SetStatus(code codes.Code, _ string) { span.status = code }

func (span *pickySpan) End(...trace.SpanEndOption) { span.ended = true }

func TestResponseAttributesThatPanic(t *testing.T) {
	defer func(rows []headerAttribute) { responseHeaderAttributes = rows }(responseHeaderAttributes)
	responseHeaderAttributes = []headerAttribute{{key: "request", header: "request-id"}}
	defer otel.SetErrorHandler(otel.GetErrorHandler())
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(error) {}))

	cases := map[string]struct {
		refused attribute.Key
		want    map[attribute.Key]any
		status  codes.Code
	}{
		"an attribute of a header": {
			"request",
			map[attribute.Key]any{HTTPResponseStatusCodeKey: int64(200), HTTPRequestResendCountKey: int64(1)},
			codes.Ok,
		},
		"the status code": {HTTPResponseStatusCodeKey, map[attribute.Key]any{"request": "req_1"}, codes.Unset},
	}
	for name, test := range cases {
		resent := func(body io.ReadCloser) (*pickySpan, *APICall) {
			span := &pickySpan{keptSpan: keptSpan{records: true}, refused: test.refused}
			call := &APICall{Span: span, Response: &http.Response{
				StatusCode: 200,
				Header:     http.Header{"Request-Id": {"req_1"}},
				Body:       body,
			}}
			call.requests.Store(2)
			return span, call
		}

		span, call := resent(http.NoBody)
		EndAPICallSpan(call, nil)
		if got := attributeValues(span.attrs); !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s, a call that has its result: got %v, want %v", name, got, test.want)
		}
		if span.status != test.status || !span.ended {
			t.Errorf("%s, a call that has its result: got status %v and ended %v, want %v and true", name, span.status, span.ended, test.status)
		}

		span, call = resent(io.NopCloser(strings.NewReader("{}")))
		EndAPICallSpanWhenBodyEnds(call, StreamObserver{})
		if got := attributeValues(span.attrs); !reflect.DeepEqual(got, test.want) {
			t.Errorf("%s, a body that is yet to be read: got %v, want %v", name, got, test.want)
		}
		if _, err := io.Copy(io.Discard, call.Response.Body); err != nil || span.status != codes.Ok || !span.ended {
			t.Errorf("%s, a body that was read: got error %v, status %v and ended %v, want nil, Ok and true", name, err, span.status, span.ended)
		}
	}
}
