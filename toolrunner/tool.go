package toolrunner

import (
	"context"
	"encoding/json"
	"fmt"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/invopop/jsonschema"
)

// betaTool is the internal generic implementation of anthropic.BetaTool.
// Users never see this type directly - they work with the BetaTool interface.
// The generic type parameter T is used internally for type-safe JSON unmarshaling.
type betaTool[T any] struct {
	name        string
	description string
	schema      anthropic.BetaToolInputSchemaParam
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

// parse validates and parses the input according to the tool's schema.
// This method handles special cases for json.RawMessage and []byte type parameters.
func (t *betaTool[T]) parse(input json.RawMessage) (T, error) {
	var parsed T

	switch any(parsed).(type) {
	case json.RawMessage:
		// If T is json.RawMessage, return the input as is
		if result, ok := any(input).(T); ok {
			return result, nil
		}
		return parsed, fmt.Errorf("type assertion failed for json.RawMessage")
	case []byte:
		// If T is []byte, return the raw JSON input as bytes
		if result, ok := any([]byte(input)).(T); ok {
			return result, nil
		}
		return parsed, fmt.Errorf("type assertion failed for []byte")
	default:
		// For all other types (structs), unmarshal the input
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
	return NewBetaTool(name, description, schema, handler), nil
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

	return NewBetaTool(name, description, inputSchema, handler), nil
}

// NewBetaTool creates a BetaTool with a BetaToolInputSchemaParam directly.
func NewBetaTool[T any](
	name, description string,
	schema anthropic.BetaToolInputSchemaParam,
	handler func(context.Context, T) (anthropic.BetaToolResultBlockParamContentUnion, error),
) anthropic.BetaTool {
	return &betaTool[T]{
		name:        name,
		description: description,
		schema:      schema,
		handler:     handler,
	}
}
