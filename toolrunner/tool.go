package toolrunner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/invopop/jsonschema"
)

// schemaValidator holds a parsed JSON Schema for runtime validation.
// It is compiled once at tool creation time from the tool's schema definition.
type schemaValidator struct {
	raw              map[string]any
	compiledPatterns map[uintptr]*regexp.Regexp // Pre-compiled regex patterns per schema node
	definitionErrors []error                    // Schema definition errors detected at construction time
}

func schemaKey(schema map[string]any) uintptr {
	return reflect.ValueOf(schema).Pointer()
}

func displayPath(path string) string {
	if path == "" {
		return "input"
	}
	return path
}

func joinPath(base, segment string) string {
	if base == "" {
		return segment
	}
	return base + "." + segment
}

func itemPath(base string, index int) string {
	return fmt.Sprintf("%s[%d]", displayPath(base), index)
}

func cloneRefStack(refStack map[string]bool) map[string]bool {
	if len(refStack) == 0 {
		return nil
	}
	cloned := make(map[string]bool, len(refStack))
	for ref, seen := range refStack {
		cloned[ref] = seen
	}
	return cloned
}

func isSupportedSchemaType(t string) bool {
	switch t {
	case "string", "number", "integer", "boolean", "array", "object", "null":
		return true
	default:
		return false
	}
}

func schemaTypes(typeValue any) ([]string, error) {
	switch t := typeValue.(type) {
	case nil:
		return nil, nil
	case string:
		if !isSupportedSchemaType(t) {
			return nil, fmt.Errorf("unsupported schema type %q", t)
		}
		return []string{t}, nil
	case []any:
		if len(t) == 0 {
			return nil, fmt.Errorf("type array must not be empty")
		}
		types := make([]string, 0, len(t))
		for _, entry := range t {
			typeName, ok := entry.(string)
			if !ok {
				return nil, fmt.Errorf("type array entries must be strings, got %T", entry)
			}
			if !isSupportedSchemaType(typeName) {
				return nil, fmt.Errorf("unsupported schema type %q", typeName)
			}
			types = append(types, typeName)
		}
		return types, nil
	default:
		return nil, fmt.Errorf("invalid type declaration %T", typeValue)
	}
}

func schemaLooksObject(schema map[string]any) bool {
	if schema == nil {
		return false
	}
	if typeValue, ok := schema["type"]; ok {
		types, err := schemaTypes(typeValue)
		if err != nil {
			return false
		}
		for _, t := range types {
			if t == "object" {
				return true
			}
		}
		return false
	}

	_, hasProps := schema["properties"]
	_, hasRequired := schema["required"]
	_, hasAdditional := schema["additionalProperties"]
	return hasProps || hasRequired || hasAdditional
}

func schemaLooksArray(schema map[string]any) bool {
	if schema == nil {
		return false
	}
	if typeValue, ok := schema["type"]; ok {
		types, err := schemaTypes(typeValue)
		if err != nil {
			return false
		}
		for _, t := range types {
			if t == "array" {
				return true
			}
		}
		return false
	}

	_, hasItems := schema["items"]
	_, hasMinItems := schema["minItems"]
	_, hasMaxItems := schema["maxItems"]
	return hasItems || hasMinItems || hasMaxItems
}

func rootSchemaShouldBeValidated(raw map[string]any) bool {
	if raw == nil {
		return false
	}
	if _, ok := raw["$ref"]; ok {
		return true
	}
	if _, ok := raw["anyOf"]; ok {
		return true
	}
	if _, ok := raw["oneOf"]; ok {
		return true
	}
	return schemaLooksObject(raw)
}

func (v *schemaValidator) addDefinitionError(err error) {
	if err != nil {
		v.definitionErrors = append(v.definitionErrors, err)
	}
}

func (v *schemaValidator) schemaDefinitionError() error {
	if v == nil || len(v.definitionErrors) == 0 {
		return nil
	}
	return errors.Join(v.definitionErrors...)
}

// newSchemaValidator creates a validator from a raw schema map.
// Returns nil if the schema is not an object-like tool input schema.
func newSchemaValidator(raw map[string]any) *schemaValidator {
	if raw == nil {
		return nil
	}
	if !rootSchemaShouldBeValidated(raw) {
		return nil
	}

	v := &schemaValidator{
		raw:              raw,
		compiledPatterns: make(map[uintptr]*regexp.Regexp),
	}
	v.prepareSchema("", raw, make(map[uintptr]bool))
	return v
}

func (v *schemaValidator) prepareSchema(path string, schema map[string]any, seen map[uintptr]bool) {
	if schema == nil {
		return
	}
	key := schemaKey(schema)
	if seen[key] {
		return
	}
	seen[key] = true

	if typeValue, ok := schema["type"]; ok {
		if _, err := schemaTypes(typeValue); err != nil {
			v.addDefinitionError(fmt.Errorf("invalid schema type for '%s': %w", displayPath(path), err))
		}
	}

	if pattern, ok := schema["pattern"].(string); ok {
		// JSON Schema's `pattern` keyword is enforced using Go's regexp engine (RE2).
		re, err := regexp.Compile(pattern)
		if err != nil {
			v.addDefinitionError(fmt.Errorf("invalid pattern '%s' for property '%s': %v", pattern, displayPath(path), err))
		} else {
			v.compiledPatterns[key] = re
		}
	}

	if ref, ok := schema["$ref"].(string); ok {
		resolved, err := v.resolveRefSchema(ref)
		if err != nil {
			v.addDefinitionError(err)
		} else {
			v.prepareSchema(path, resolved, seen)
		}
	}

	if defs, ok := schema["$defs"].(map[string]any); ok {
		names := make([]string, 0, len(defs))
		for name := range defs {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			defSchema, ok := defs[name].(map[string]any)
			if !ok {
				continue
			}
			v.prepareSchema(joinPath("$defs", name), defSchema, seen)
		}
	}

	if props, ok := schema["properties"].(map[string]any); ok {
		names := make([]string, 0, len(props))
		for name := range props {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			propSchema, ok := props[name].(map[string]any)
			if !ok {
				continue
			}
			v.prepareSchema(joinPath(path, name), propSchema, seen)
		}
	}

	if items, ok := schema["items"].(map[string]any); ok {
		v.prepareSchema(joinPath(path, "[]"), items, seen)
	}

	if additionalSchema, ok := schema["additionalProperties"].(map[string]any); ok {
		v.prepareSchema(joinPath(path, "*"), additionalSchema, seen)
	}

	for _, keyword := range []string{"anyOf", "oneOf"} {
		variants, ok := schema[keyword].([]any)
		if !ok {
			continue
		}
		for i, variant := range variants {
			variantSchema, ok := variant.(map[string]any)
			if !ok {
				continue
			}
			v.prepareSchema(fmt.Sprintf("%s.%s[%d]", displayPath(path), keyword, i), variantSchema, seen)
		}
	}
}

// validate checks an unmarshaled JSON value against the schema.
// It recursively enforces the supported JSON Schema subset: required fields,
// additionalProperties, type/union checks, enum, pattern, string and numeric
// bounds, array length and item validation, nested object properties, anyOf/oneOf,
// and local $ref/$defs resolution.
func (v *schemaValidator) validate(input any) error {
	if v == nil {
		return nil
	}
	if err := v.schemaDefinitionError(); err != nil {
		return err
	}
	obj, ok := input.(map[string]any)
	if !ok {
		return fmt.Errorf("expected object, got %T", input)
	}
	return v.validateValue("", obj, v.raw, nil)
}

func (v *schemaValidator) validateValue(path string, value any, schema map[string]any, refStack map[string]bool) error {
	if schema == nil {
		return nil
	}

	if ref, ok := schema["$ref"].(string); ok {
		if refStack == nil {
			refStack = make(map[string]bool)
		}
		if refStack[ref] {
			return fmt.Errorf("cyclic schema reference '%s'", ref)
		}
		resolved, err := v.resolveRefSchema(ref)
		if err != nil {
			return err
		}
		refStack[ref] = true
		err = v.validateValue(path, value, resolved, refStack)
		delete(refStack, ref)
		return err
	}

	if variants, ok := schema["anyOf"]; ok {
		if err := v.validateVariants(path, value, variants, false, refStack); err != nil {
			return err
		}
	}
	if variants, ok := schema["oneOf"]; ok {
		if err := v.validateVariants(path, value, variants, true, refStack); err != nil {
			return err
		}
	}

	if typeValue, ok := schema["type"]; ok {
		if err := validateType(displayPath(path), value, typeValue); err != nil {
			return err
		}
	}
	if _, hasType := schema["type"]; !hasType && schemaLooksObject(schema) {
		if _, ok := value.(map[string]any); !ok {
			return fmt.Errorf("property '%s' expected type object, got %s", displayPath(path), valueTypeName(value))
		}
	}
	if _, hasType := schema["type"]; !hasType && schemaLooksArray(schema) {
		if _, ok := value.([]any); !ok {
			return fmt.Errorf("property '%s' expected type array, got %s", displayPath(path), valueTypeName(value))
		}
	}

	if enumVals, ok := schema["enum"].([]any); ok {
		found := false
		for _, e := range enumVals {
			if reflect.DeepEqual(e, value) {
				found = true
				break
			}
		}
		if !found {
			allowed := make([]string, len(enumVals))
			for i, e := range enumVals {
				allowed[i] = fmt.Sprintf("%v", e)
			}
			return fmt.Errorf("property '%s' value '%v' is not one of the allowed values: [%s]", displayPath(path), value, strings.Join(allowed, ", "))
		}
	}

	if s, ok := value.(string); ok {
		if _, hasPattern := schema["pattern"].(string); hasPattern {
			if re, ok := v.compiledPatterns[schemaKey(schema)]; ok && !re.MatchString(s) {
				return fmt.Errorf("property '%s' value does not match pattern '%s'", displayPath(path), schema["pattern"])
			}
		}
		charCount := utf8.RuneCountInString(s)
		if minLen, ok := toFloat64(schema["minLength"]); ok {
			if float64(charCount) < minLen {
				return fmt.Errorf("property '%s' length %d is less than minimum %d", displayPath(path), charCount, int(minLen))
			}
		}
		if maxLen, ok := toFloat64(schema["maxLength"]); ok {
			if float64(charCount) > maxLen {
				return fmt.Errorf("property '%s' length %d exceeds maximum %d", displayPath(path), charCount, int(maxLen))
			}
		}
	}

	if f, ok := value.(float64); ok {
		if min, ok := toFloat64(schema["minimum"]); ok {
			if f < min {
				return fmt.Errorf("property '%s' value %v is less than minimum %v", displayPath(path), f, min)
			}
		}
		if max, ok := toFloat64(schema["maximum"]); ok {
			if f > max {
				return fmt.Errorf("property '%s' value %v exceeds maximum %v", displayPath(path), f, max)
			}
		}
		if eMin, ok := toFloat64(schema["exclusiveMinimum"]); ok {
			if f <= eMin {
				return fmt.Errorf("property '%s' value %v must be greater than %v", displayPath(path), f, eMin)
			}
		}
		if eMax, ok := toFloat64(schema["exclusiveMaximum"]); ok {
			if f >= eMax {
				return fmt.Errorf("property '%s' value %v must be less than %v", displayPath(path), f, eMax)
			}
		}
	}

	if arr, ok := value.([]any); ok {
		if err := v.validateArray(path, arr, schema, refStack); err != nil {
			return err
		}
	}
	if obj, ok := value.(map[string]any); ok {
		if err := v.validateObject(path, obj, schema, refStack); err != nil {
			return err
		}
	}

	return nil
}

func (v *schemaValidator) validateVariants(path string, value any, variantsValue any, exactlyOne bool, refStack map[string]bool) error {
	variants, ok := variantsValue.([]any)
	if !ok {
		return nil
	}

	matches := 0
	for _, variant := range variants {
		variantSchema, ok := variant.(map[string]any)
		if !ok {
			continue
		}
		if err := v.validateValue(path, value, variantSchema, cloneRefStack(refStack)); err == nil {
			matches++
			if !exactlyOne {
				return nil
			}
		}
	}

	if exactlyOne {
		if matches == 1 {
			return nil
		}
		return fmt.Errorf("property '%s' must match exactly one allowed schema variant", displayPath(path))
	}
	return fmt.Errorf("property '%s' did not match any allowed schema variant", displayPath(path))
}

func (v *schemaValidator) validateObject(path string, obj map[string]any, schema map[string]any, refStack map[string]bool) error {
	if req, ok := schema["required"].([]any); ok {
		for _, r := range req {
			name, _ := r.(string)
			if name == "" {
				continue
			}
			if _, exists := obj[name]; !exists {
				return fmt.Errorf("missing required property '%s'", joinPath(path, name))
			}
		}
	}

	props, _ := schema["properties"].(map[string]any)
	additional := schema["additionalProperties"]
	for key, val := range obj {
		propPath := joinPath(path, key)
		if propRaw, defined := props[key]; defined {
			propSchema, ok := propRaw.(map[string]any)
			if !ok {
				continue
			}
			if err := v.validateValue(propPath, val, propSchema, refStack); err != nil {
				return err
			}
			continue
		}

		if additionalSchema, ok := additional.(map[string]any); ok {
			if err := v.validateValue(propPath, val, additionalSchema, refStack); err != nil {
				return err
			}
			continue
		}

		if val, isBool := additional.(bool); isBool && !val {
			return fmt.Errorf("additional property '%s' is not allowed", propPath)
		}
	}
	return nil
}

func (v *schemaValidator) validateArray(path string, arr []any, schema map[string]any, refStack map[string]bool) error {
	if minItems, ok := toFloat64(schema["minItems"]); ok {
		if float64(len(arr)) < minItems {
			return fmt.Errorf("property '%s' has %d items, minimum is %d", displayPath(path), len(arr), int(minItems))
		}
	}
	if maxItems, ok := toFloat64(schema["maxItems"]); ok {
		if float64(len(arr)) > maxItems {
			return fmt.Errorf("property '%s' has %d items, maximum is %d", displayPath(path), len(arr), int(maxItems))
		}
	}
	items, ok := schema["items"].(map[string]any)
	if !ok {
		return nil
	}
	for i, item := range arr {
		if err := v.validateValue(itemPath(path, i), item, items, refStack); err != nil {
			return err
		}
	}
	return nil
}

func (v *schemaValidator) resolveRefSchema(ref string) (map[string]any, error) {
	if ref == "#" {
		return v.raw, nil
	}
	if !strings.HasPrefix(ref, "#/") {
		return nil, fmt.Errorf("unsupported schema reference '%s': only local '#/...' refs are supported", ref)
	}

	current := any(v.raw)
	for _, token := range strings.Split(strings.TrimPrefix(ref, "#/"), "/") {
		decoded := strings.ReplaceAll(strings.ReplaceAll(token, "~1", "/"), "~0", "~")
		obj, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("schema reference '%s' does not resolve to an object schema", ref)
		}
		next, ok := obj[decoded]
		if !ok {
			return nil, fmt.Errorf("schema reference '%s' could not be resolved", ref)
		}
		current = next
	}

	resolved, ok := current.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("schema reference '%s' does not resolve to an object schema", ref)
	}
	return resolved, nil
}

func valueTypeName(value any) string {
	switch value.(type) {
	case nil:
		return "null"
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "boolean"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	default:
		return fmt.Sprintf("%T", value)
	}
}

func validateType(name string, value any, typeValue any) error {
	expectedTypes, err := schemaTypes(typeValue)
	if err != nil {
		return fmt.Errorf("invalid schema type for '%s': %w", name, err)
	}
	if len(expectedTypes) == 0 {
		return nil
	}

	for _, expected := range expectedTypes {
		switch expected {
		case "string":
			if _, ok := value.(string); ok {
				return nil
			}
		case "number":
			if _, ok := value.(float64); ok {
				return nil
			}
		case "integer":
			if f, ok := value.(float64); ok && f == float64(int64(f)) {
				return nil
			}
		case "boolean":
			if _, ok := value.(bool); ok {
				return nil
			}
		case "array":
			if _, ok := value.([]any); ok {
				return nil
			}
		case "object":
			if _, ok := value.(map[string]any); ok {
				return nil
			}
		case "null":
			if value == nil {
				return nil
			}
		}
	}

	if len(expectedTypes) == 1 && expectedTypes[0] == "integer" {
		if f, ok := value.(float64); ok {
			return fmt.Errorf("property '%s' expected integer, got float %v", name, f)
		}
	}
	if len(expectedTypes) == 1 {
		return fmt.Errorf("property '%s' expected type %s, got %s", name, expectedTypes[0], valueTypeName(value))
	}
	return fmt.Errorf("property '%s' expected one of types [%s], got %s", name, strings.Join(expectedTypes, ", "), valueTypeName(value))
}

// toFloat64 extracts a float64 from a JSON-decoded number (which is always float64).
func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	default:
		return 0, false
	}
}

// betaTool is the internal generic implementation of anthropic.BetaTool.
// Users never see this type directly - they work with the BetaTool interface.
// The generic type parameter T is used internally for type-safe JSON unmarshaling.
type betaTool[T any] struct {
	name        string
	description string
	schema      anthropic.BetaToolInputSchemaParam
	rawSchema   map[string]any // Original schema for validation (avoids marshal roundtrip losses)
	validator   *schemaValidator
	handler     func(context.Context, T) (anthropic.BetaToolResultBlockParamContentUnion, error)
}

func (t *betaTool[T]) Name() string                                    { return t.name }
func (t *betaTool[T]) Description() string                             { return t.description }
func (t *betaTool[T]) InputSchema() anthropic.BetaToolInputSchemaParam { return t.schema }

func (t *betaTool[T]) Execute(ctx context.Context, input json.RawMessage) (anthropic.BetaToolResultBlockParamContentUnion, error) {
	parsed, err := t.parse(input)
	if err != nil {
		return anthropic.BetaToolResultBlockParamContentUnion{}, fmt.Errorf("failed to parse tool input: %w", err)
	}
	return t.handler(ctx, parsed)
}

// parse validates the input against the tool's JSON Schema and then unmarshals
// it into the target type T. Validation enforces required fields, additionalProperties,
// type correctness, enum constraints, pattern, string length bounds, numeric bounds,
// and array item counts before the handler runs.
// This method handles special cases for json.RawMessage and []byte type parameters.
func (t *betaTool[T]) parse(input json.RawMessage) (T, error) {
	var parsed T

	switch any(parsed).(type) {
	case json.RawMessage:
		if result, ok := any(input).(T); ok {
			return result, nil
		}
		return parsed, fmt.Errorf("type assertion failed for json.RawMessage")
	case []byte:
		if result, ok := any([]byte(input)).(T); ok {
			return result, nil
		}
		return parsed, fmt.Errorf("type assertion failed for []byte")
	default:
		// Validate against JSON Schema before unmarshaling into the typed struct
		if t.validator != nil {
			var inputData any
			if err := json.Unmarshal(input, &inputData); err != nil {
				return parsed, fmt.Errorf("invalid JSON: %w", err)
			}
			if err := t.validator.validate(inputData); err != nil {
				return parsed, fmt.Errorf("schema validation failed: %w", err)
			}
		}

		if err := json.Unmarshal(input, &parsed); err != nil {
			return parsed, err
		}
		return parsed, nil
	}
}

func parseSchemaMap(s map[string]any) (anthropic.BetaToolInputSchemaParam, error) {
	bytes, err := json.Marshal(s)
	if err != nil {
		return anthropic.BetaToolInputSchemaParam{}, fmt.Errorf("failed to marshal schema: %w", err)
	}

	var schema anthropic.BetaToolInputSchemaParam
	if err := json.Unmarshal(bytes, &schema); err != nil {
		return anthropic.BetaToolInputSchemaParam{}, fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	return schema, nil
}

// NewBetaToolFromBytes creates a BetaTool from JSON schema bytes.
func NewBetaToolFromBytes[T any](
	name, description string,
	schemaJSON []byte,
	handler func(context.Context, T) (anthropic.BetaToolResultBlockParamContentUnion, error),
) (anthropic.BetaTool, error) {
	var schema anthropic.BetaToolInputSchemaParam
	if err := schema.UnmarshalJSON(schemaJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}
	// Preserve raw schema for validation (avoids BetaToolInputSchemaParam marshal losses)
	var rawSchema map[string]any
	if err := json.Unmarshal(schemaJSON, &rawSchema); err != nil {
		return nil, fmt.Errorf("failed to parse raw schema: %w", err)
	}
	validator := newSchemaValidator(rawSchema)
	if err := validator.schemaDefinitionError(); err != nil {
		return nil, fmt.Errorf("invalid tool schema: %w", err)
	}
	return newBetaToolWithValidator(name, description, schema, rawSchema, validator, handler), nil
}

// NewBetaToolFromJSONSchema creates a BetaTool by inferring the schema from struct type T using reflection.
// The struct should use jsonschema tags to define the schema (e.g., `jsonschema:"required,description=..."`).
func NewBetaToolFromJSONSchema[T any](
	name, description string,
	handler func(context.Context, T) (anthropic.BetaToolResultBlockParamContentUnion, error),
) (anthropic.BetaTool, error) {
	var zeroValue T
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties:  false,
		RequiredFromJSONSchemaTags: true,
		DoNotReference:             true,
	}

	schema := reflector.Reflect(zeroValue)

	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	var schemaMap map[string]any
	if err := json.Unmarshal(schemaBytes, &schemaMap); err != nil {
		return nil, err
	}

	inputSchema, err := parseSchemaMap(schemaMap)
	if err != nil {
		return nil, err
	}

	validator := newSchemaValidator(schemaMap)
	if err := validator.schemaDefinitionError(); err != nil {
		return nil, fmt.Errorf("invalid tool schema: %w", err)
	}
	return newBetaToolWithValidator(name, description, inputSchema, schemaMap, validator, handler), nil
}

// newBetaTool is the internal constructor that accepts both the typed schema param
// and the raw schema map (to avoid losing fields like additionalProperties during marshal).
func newBetaTool[T any](
	name, description string,
	schema anthropic.BetaToolInputSchemaParam,
	rawSchema map[string]any,
	handler func(context.Context, T) (anthropic.BetaToolResultBlockParamContentUnion, error),
) anthropic.BetaTool {
	return newBetaToolWithValidator(name, description, schema, rawSchema, newSchemaValidator(rawSchema), handler)
}

func newBetaToolWithValidator[T any](
	name, description string,
	schema anthropic.BetaToolInputSchemaParam,
	rawSchema map[string]any,
	validator *schemaValidator,
	handler func(context.Context, T) (anthropic.BetaToolResultBlockParamContentUnion, error),
) anthropic.BetaTool {
	return &betaTool[T]{
		name:        name,
		description: description,
		schema:      schema,
		rawSchema:   rawSchema,
		validator:   validator,
		handler:     handler,
	}
}

// NewBetaTool creates a BetaTool with a BetaToolInputSchemaParam directly.
// The schema is parsed at creation time for efficient runtime validation.
func NewBetaTool[T any](
	name, description string,
	schema anthropic.BetaToolInputSchemaParam,
	handler func(context.Context, T) (anthropic.BetaToolResultBlockParamContentUnion, error),
) anthropic.BetaTool {
	// Extract raw schema via marshal roundtrip (best-effort; some fields may be lost)
	var rawSchema map[string]any
	if b, err := json.Marshal(schema); err == nil {
		json.Unmarshal(b, &rawSchema)
	}
	return newBetaTool(name, description, schema, rawSchema, handler)
}
