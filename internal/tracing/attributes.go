package tracing

import (
	"net/http"
	"slices"
	"strings"

	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel/attribute"
)

// valueType is the type of the value of an attribute that a span reads from JSON:
// the type that the API description gives the property.
type valueType int

const (
	// anyType is the type of a property that the API description gives more than
	// one type. The attribute has the type of the value that the JSON has.
	anyType valueType = iota
	stringType
	integerType
	numberType
	booleanType
)

// bodyAttribute is an attribute that a span reads from the JSON body of a
// response.
type bodyAttribute struct {
	key attribute.Key
	// property is the property of the body that is the value, or that has the
	// value at path, a gjson path, when path is not empty.
	property, path string
	typ            valueType
	// list is whether the value is an array of typ.
	list bool
}

// headerAttribute is an attribute that a span reads from a header of a request
// or of a response.
type headerAttribute struct {
	key    attribute.Key
	header string
	// list is whether the attribute has one string for each value in the header.
	list bool
}

// spanTable is what the spans of an SDK method have of the method's request and
// response.
type spanTable struct {
	// constants are the attributes that have a value of their own.
	constants      []attribute.KeyValue
	requestHeaders []headerAttribute
	// stream is the attribute that is true when the method streams the response.
	stream   attribute.Key
	response []bodyAttribute
}

// spanAttributes has the table of each span name that has one.
var spanAttributes = map[string]spanTable{
	"anthropic.messages.create": {
		constants: []attribute.KeyValue{
			attribute.String("gen_ai.operation.name", "chat"),
		},
		requestHeaders: []headerAttribute{
			{key: "anthropic.beta", header: "anthropic-beta", list: true},
		},
		stream: "gen_ai.request.stream",
		response: []bodyAttribute{
			{key: "gen_ai.response.id", property: "id", typ: stringType},
			{key: "gen_ai.response.model", property: "model", typ: stringType},
			{key: "gen_ai.usage.cache_write.input_tokens", property: "usage", path: "cache_creation_input_tokens", typ: integerType},
			{key: "gen_ai.usage.cache_read.input_tokens", property: "usage", path: "cache_read_input_tokens", typ: integerType},
			{key: "anthropic.usage.cache_creation.ephemeral_5m_input_tokens", property: "usage", path: "cache_creation.ephemeral_5m_input_tokens", typ: integerType},
			{key: "anthropic.usage.cache_creation.ephemeral_1h_input_tokens", property: "usage", path: "cache_creation.ephemeral_1h_input_tokens", typ: integerType},
			{key: "anthropic.response.service_tier", property: "usage", path: "service_tier", typ: stringType},
			{key: "gen_ai.usage.output_tokens", property: "usage", path: "output_tokens", typ: integerType},
			{key: "gen_ai.usage.reasoning.output_tokens", property: "usage", path: "output_tokens_details.thinking_tokens", typ: integerType},
			{key: "anthropic.message.stop_reason", property: "stop_reason", typ: stringType},
		},
	},
	"anthropic.messages.count_tokens": {
		response: []bodyAttribute{
			{key: "anthropic.message.count_tokens.input_tokens", property: "input_tokens", typ: integerType},
		},
	},
	"anthropic.sessions.create": {
		response: []bodyAttribute{
			{key: "gen_ai.agent.name", property: "agent", path: "name", typ: stringType},
			{key: "gen_ai.request.model", property: "agent", path: "model.id", typ: stringType},
		},
	},
}

// responseHeaderAttributes are the attributes that every span has from the
// headers of the response that ends its call.
var responseHeaderAttributes = []headerAttribute{
	{key: "anthropic.request.id", header: "request-id"},
	{key: "anthropic.organization.id", header: "anthropic-organization-id"},
	{key: "anthropic.workspace.id", header: "anthropic-workspace-id"},
}

// ClientAttributes returns the attributes that every span starts with.
func ClientAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{attribute.String("gen_ai.provider.name", "anthropic")}
}

// SpanRequestAttributes returns the attributes that the span named name starts
// with of its request, beside the ones that the SDK method reads from its
// parameters: the ones with a value of their own, the ones that the request has
// in header, its header, and the one for a response that the method streams,
// when streams is true.
func SpanRequestAttributes(name string, header http.Header, streams bool) []attribute.KeyValue {
	table := spanAttributes[name]
	attrs := append(slices.Clone(table.constants), headerAttributes(table.requestHeaders, header)...)
	if streams && table.stream != "" {
		attrs = append(attrs, table.stream.Bool(true))
	}
	return attrs
}

// SetSpanResponseAttributes sets on the span of call the attributes that a span
// of its name has from body, the JSON of the call's response. requestconfig calls
// it with the body that it decoded into the value of the SDK method. For a
// response that the method streams, the SDK's own code calls it with what it has
// read of the stream. Bytes are copied only for a span that reads them.
func SetSpanResponseAttributes[JSON string | []byte](call *APICall, body JSON) {
	if rows := spanAttributes[call.Name].response; len(rows) > 0 && call.Span.IsRecording() {
		call.Span.SetAttributes(bodyAttributes(rows, gjson.Parse(string(body)))...)
	}
}

// SpanResponsePaths returns the gjson path of each value that the span named
// name reads from the JSON of its response. The SDK's own code checks against it
// what it lists of these paths.
func SpanResponsePaths(name string) []string {
	var paths []string
	for _, row := range spanAttributes[name].response {
		paths = append(paths, strings.TrimSuffix(row.property+"."+row.path, "."))
	}
	return paths
}

// setResponseHeaderAttributes sets [responseHeaderAttributes] on the span of
// call from its response, when one arrived.
func setResponseHeaderAttributes(call *APICall) {
	if call.Response != nil {
		call.Span.SetAttributes(headerAttributes(responseHeaderAttributes, call.Response.Header)...)
	}
}

// bodyAttributes returns the attributes that body, the JSON of a response's
// body, has a value for among rows. A property that body lacks, a value of
// another type than its row's, and an empty string set no attribute.
func bodyAttributes(rows []bodyAttribute, body gjson.Result) []attribute.KeyValue {
	var attrs []attribute.KeyValue
	// One pass over body, whatever the number of rows: the body of a response can be large.
	body.ForEach(func(property, value gjson.Result) bool {
		for _, row := range rows {
			if row.property != property.Str {
				continue
			}
			found := value
			if row.path != "" {
				found = value.Get(row.path)
			}
			if kv, ok := row.keyValue(found); ok {
				attrs = append(attrs, kv)
			}
		}
		return true
	})
	return attrs
}

// keyValue returns the attribute of row with value, and false when value is not
// of the row's type or is an empty string.
func (row bodyAttribute) keyValue(value gjson.Result) (kv attribute.KeyValue, ok bool) {
	if row.list {
		return row.listKeyValue(value)
	}
	typ := row.typ
	if typ == anyType {
		typ = typeOf(value)
	}
	if !typ.has(value) || typ == stringType && value.Str == "" {
		return kv, false
	}
	switch typ {
	case stringType:
		// A string of gjson can share the memory of the body: the attribute, which
		// lives until the span is exported, must not keep the whole body with it.
		return row.key.String(strings.Clone(value.Str)), true
	case integerType:
		return row.key.Int64(value.Int()), true
	case numberType:
		return row.key.Float64(value.Num), true
	default:
		return row.key.Bool(value.Bool()), true
	}
}

// listKeyValue is [bodyAttribute.keyValue] for a row whose value is an array.
func (row bodyAttribute) listKeyValue(value gjson.Result) (kv attribute.KeyValue, ok bool) {
	if !value.IsArray() {
		return kv, false
	}
	items := value.Array()
	switch row.typ {
	case stringType:
		values, ok := each(items, row.typ, func(item gjson.Result) string { return strings.Clone(item.Str) })
		return row.key.StringSlice(values), ok
	case integerType:
		values, ok := each(items, row.typ, gjson.Result.Int)
		return row.key.Int64Slice(values), ok
	case numberType:
		values, ok := each(items, row.typ, gjson.Result.Float)
		return row.key.Float64Slice(values), ok
	case booleanType:
		values, ok := each(items, row.typ, gjson.Result.Bool)
		return row.key.BoolSlice(values), ok
	default:
		return kv, false
	}
}

// each returns what as gives each of items, and false when one of them is not of typ.
func each[T any](items []gjson.Result, typ valueType, as func(gjson.Result) T) ([]T, bool) {
	values := make([]T, len(items))
	for i, item := range items {
		if !typ.has(item) {
			return nil, false
		}
		values[i] = as(item)
	}
	return values, true
}

// has reports whether value, a JSON value, is of typ.
func (typ valueType) has(value gjson.Result) bool {
	switch typ {
	case stringType:
		return value.Type == gjson.String
	case integerType, numberType:
		return value.Type == gjson.Number
	case booleanType:
		return value.Type == gjson.True || value.Type == gjson.False
	default:
		return false
	}
}

// typeOf returns the type of value, a JSON value, and [anyType] for a value
// that no attribute can have: null, an array or an object.
func typeOf(value gjson.Result) valueType {
	switch value.Type {
	case gjson.String:
		return stringType
	case gjson.Number:
		if strings.ContainsAny(value.Raw, ".eE") {
			return numberType
		}
		return integerType
	case gjson.True, gjson.False:
		return booleanType
	default:
		return anyType
	}
}

// headerAttributes returns the attributes that header has a value for among rows.
func headerAttributes(rows []headerAttribute, header http.Header) []attribute.KeyValue {
	var attrs []attribute.KeyValue
	for _, row := range rows {
		if !row.list {
			if value := header.Get(row.header); value != "" {
				attrs = append(attrs, row.key.String(value))
			}
			continue
		}
		var values []string
		for _, line := range header.Values(row.header) {
			for _, value := range strings.Split(line, ",") {
				if value = strings.TrimSpace(value); value != "" {
					values = append(values, value)
				}
			}
		}
		if len(values) > 0 {
			attrs = append(attrs, row.key.StringSlice(values))
		}
	}
	return attrs
}
