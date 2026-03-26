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
	compiledPatterns map[string]*regexp.Regexp // Pre-compiled regex patterns per property
	patternErrors    map[string]error          // Compile errors for invalid patterns
}

func (v *schemaValidator) schemaDefinitionError() error {
	if v == nil || len(v.patternErrors) == 0 {
		return nil
	}

	names := make([]string, 0, len(v.patternErrors))
	for name := range v.patternErrors {
		names = append(names, name)
	}
	sort.Strings(names)

	errs := make([]error, 0, len(names))
	for _, name := range names {
		errs = append(errs, v.patternErrors[name])
	}
	return errors.Join(errs...)
}

// newSchemaValidator creates a validator from a raw schema map.
// Returns nil if the schema is not an object-type schema (validation will be skipped).
// Recognizes object schemas by: type "object", type array containing "object",
// or missing type with properties, required, or additionalProperties present.
func newSchemaValidator(raw map[string]any) *schemaValidator {
	if raw == nil {
		return nil
	}

	isObject := false
	if typeVal, ok := raw["type"]; ok {
		switch t := typeVal.(type) {
		case string:
			isObject = t == "object"
		case []any:
			for _, v := range t {
				if s, ok := v.(string); ok && s == "object" {
					isObject = true
					break
				}
			}
		}
	} else {
		// No type field: infer object from object-specific keywords
		_, hasProps := raw["properties"]
		_, hasReq := raw["required"]
		_, hasAdditional := raw["additionalProperties"]
		isObject = hasProps || hasReq || hasAdditional
	}

	if !isObject {
		return nil
	}

	v := &schemaValidator{
		raw:              raw,
		compiledPatterns: make(map[string]*regexp.Regexp),
		patternErrors:    make(map[string]error),
	}

	// Pre-compile regex patterns from property schemas.
	// JSON Schema's `pattern` keyword is enforced using Go's regexp engine (RE2).
	if props, ok := raw["properties"].(map[string]any); ok {
		for name, propRaw := range props {
			propSchema, ok := propRaw.(map[string]any)
			if !ok {
				continue
			}
			if pattern, ok := propSchema["pattern"].(string); ok {
				re, err := regexp.Compile(pattern)
				if err != nil {
					v.patternErrors[name] = fmt.Errorf("invalid pattern '%s' for property '%s': %v", pattern, name, err)
				} else {
					v.compiledPatterns[name] = re
				}
			}
		}
	}

	return v
}

// validate checks an unmarshaled JSON value against the schema.
// It enforces: required fields, additionalProperties, property types, enum
// constraints, pattern, string length bounds, and numeric bounds.
func (v *schemaValidator) validate(input any) error {
	if v == nil {
		return nil
	}
	obj, ok := input.(map[string]any)
	if !ok {
		return fmt.Errorf("expected object, got %T", input)
	}

	// Check required fields
	if req, ok := v.raw["required"].([]any); ok {
		for _, r := range req {
			name, _ := r.(string)
			if name == "" {
				continue
			}
			if _, exists := obj[name]; !exists {
				return fmt.Errorf("missing required property '%s'", name)
			}
		}
	}

	props, _ := v.raw["properties"].(map[string]any)

	// Check additionalProperties
	if additional, ok := v.raw["additionalProperties"]; ok {
		if val, isBool := additional.(bool); isBool && !val {
			for key := range obj {
				if props == nil {
					return fmt.Errorf("additional property '%s' is not allowed", key)
				}
				if _, defined := props[key]; !defined {
					return fmt.Errorf("additional property '%s' is not allowed", key)
				}
			}
		}
	}

	// Check property constraints
	if props == nil {
		return nil
	}
	for key, val := range obj {
		propSchema, ok := props[key].(map[string]any)
		if !ok {
			continue
		}
		if err := v.validateProperty(key, val, propSchema); err != nil {
			return err
		}
	}
	return nil
}

// validateProperty checks a single property value against its schema definition.
// Enforces: type, enum, pattern (pre-compiled), minLength, maxLength, minimum,
// maximum, minItems, maxItems.
func (v *schemaValidator) validateProperty(name string, value any, propSchema map[string]any) error {
	// Type check
	if expectedType, ok := propSchema["type"].(string); ok {
		if err := checkType(name, value, expectedType); err != nil {
			return err
		}
	}

	// Enum check (type-strict: string "1" does not match number 1)
	if enumVals, ok := propSchema["enum"].([]any); ok {
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
			return fmt.Errorf("property '%s' value '%v' is not one of the allowed values: [%s]", name, value, strings.Join(allowed, ", "))
		}
	}

	// String constraints
	if s, ok := value.(string); ok {
		if _, hasPattern := propSchema["pattern"].(string); hasPattern {
			// Check for pre-compiled pattern error (invalid regex in schema)
			if err, hasFail := v.patternErrors[name]; hasFail {
				return err
			}
			// Use pre-compiled regex
			if re, ok := v.compiledPatterns[name]; ok && !re.MatchString(s) {
				return fmt.Errorf("property '%s' value does not match pattern '%s'", name, propSchema["pattern"])
			}
		}
		charCount := utf8.RuneCountInString(s)
		if minLen, ok := toFloat64(propSchema["minLength"]); ok {
			if float64(charCount) < minLen {
				return fmt.Errorf("property '%s' length %d is less than minimum %d", name, charCount, int(minLen))
			}
		}
		if maxLen, ok := toFloat64(propSchema["maxLength"]); ok {
			if float64(charCount) > maxLen {
				return fmt.Errorf("property '%s' length %d exceeds maximum %d", name, charCount, int(maxLen))
			}
		}
	}

	// Numeric constraints
	if f, ok := value.(float64); ok {
		if min, ok := toFloat64(propSchema["minimum"]); ok {
			if f < min {
				return fmt.Errorf("property '%s' value %v is less than minimum %v", name, f, min)
			}
		}
		if max, ok := toFloat64(propSchema["maximum"]); ok {
			if f > max {
				return fmt.Errorf("property '%s' value %v exceeds maximum %v", name, f, max)
			}
		}
		if eMin, ok := toFloat64(propSchema["exclusiveMinimum"]); ok {
			if f <= eMin {
				return fmt.Errorf("property '%s' value %v must be greater than %v", name, f, eMin)
			}
		}
		if eMax, ok := toFloat64(propSchema["exclusiveMaximum"]); ok {
			if f >= eMax {
				return fmt.Errorf("property '%s' value %v must be less than %v", name, f, eMax)
			}
		}
	}

	// Array constraints
	if arr, ok := value.([]any); ok {
		if minItems, ok := toFloat64(propSchema["minItems"]); ok {
			if float64(len(arr)) < minItems {
				return fmt.Errorf("property '%s' has %d items, minimum is %d", name, len(arr), int(minItems))
			}
		}
		if maxItems, ok := toFloat64(propSchema["maxItems"]); ok {
			if float64(len(arr)) > maxItems {
				return fmt.Errorf("property '%s' has %d items, maximum is %d", name, len(arr), int(maxItems))
			}
		}
	}

	return nil
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

// checkType validates that a JSON value matches the expected JSON Schema type.
func checkType(name string, value any, expected string) error {
	switch expected {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("property '%s' expected type string, got %T", name, value)
		}
	case "number":
		if _, ok := value.(float64); !ok {
			return fmt.Errorf("property '%s' expected type number, got %T", name, value)
		}
	case "integer":
		f, ok := value.(float64)
		if !ok {
			return fmt.Errorf("property '%s' expected type integer, got %T", name, value)
		}
		if f != float64(int64(f)) {
			return fmt.Errorf("property '%s' expected integer, got float %v", name, f)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("property '%s' expected type boolean, got %T", name, value)
		}
	case "array":
		if _, ok := value.([]any); !ok {
			return fmt.Errorf("property '%s' expected type array, got %T", name, value)
		}
	case "object":
		if _, ok := value.(map[string]any); !ok {
			return fmt.Errorf("property '%s' expected type object, got %T", name, value)
		}
	}
	return nil
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
