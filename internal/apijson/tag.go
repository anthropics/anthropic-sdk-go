package apijson

import (
	"reflect"
	"strings"
)

const (
	jsonStructTag   = "json"
	formatStructTag = "format"
)

type parsedStructTag struct {
	name     string
	required bool
	extras   bool
	metadata bool
	inline   bool
}

func parseJSONStructTag(field reflect.StructField) (tag parsedStructTag, ok bool) {
	raw, ok := field.Tag.Lookup(jsonStructTag)
	if !ok {
		return tag, ok
	}
	parts := strings.Split(raw, ",")
	if len(parts) == 0 {
		return tag, false
	}
	tag.name = parts[0]
	for _, part := range parts[1:] {
		switch part {
		case "required":
			tag.required = true
		case "extras":
			tag.extras = true
		case "metadata":
			tag.metadata = true
		case "inline":
			tag.inline = true
		}
	}
	return tag, ok
}

func parseFormatStructTag(field reflect.StructField) (format string, ok bool) {
	format, ok = field.Tag.Lookup(formatStructTag)
	return format, ok
}
