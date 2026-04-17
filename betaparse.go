package anthropic

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/invopop/jsonschema"
)

// ErrStructuredOutputParse is returned (wrapped) by [BetaMessageService.New]
// when the API request succeeds but the response body can't be unmarshaled
// into the struct pointer passed as Schema. The *BetaMessage is still
// returned in that case — check it alongside the error:
//
//	msg, err := client.Beta.Messages.New(ctx, params)
//	if errors.Is(err, anthropic.ErrStructuredOutputParse) {
//	    // msg is valid; the model's response didn't match the struct shape.
//	}
var ErrStructuredOutputParse = errors.New("anthropic: failed to parse structured output")

// schemaCache caches the final json.RawMessage (post-transform) keyed by
// reflect.Type, so repeated requests with the same struct type skip reflection,
// transformation, and marshaling.
//
// Entries are never evicted. This is fine for the common case where a program
// uses a fixed set of struct types (typically < 100, ~1KB each). Programs that
// synthesize types dynamically via reflect.StructOf with per-request shapes
// should precompute their schemas instead of relying on this cache.
var schemaCache sync.Map // reflect.Type → json.RawMessage

// schemaToRaw converts a value to a json.RawMessage JSON schema suitable for the wire.
// If v is already a json.RawMessage, it is returned as-is.
// If v is a map[string]any, it is marshaled directly to json.RawMessage.
// If v is a struct pointer, the SDK reflects the type into a JSON schema, transforms it,
// and returns the result as json.RawMessage.
func schemaToRaw(v any) (json.RawMessage, error) {
	if v == nil {
		return nil, nil
	}

	switch s := v.(type) {
	case json.RawMessage:
		return s, nil
	case map[string]any:
		b, err := json.Marshal(s)
		if err != nil {
			return nil, fmt.Errorf("anthropic: failed to marshal schema: %w", err)
		}
		return json.RawMessage(b), nil
	}

	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return nil, fmt.Errorf("anthropic: Schema must be a non-nil pointer to a struct, map[string]any, or json.RawMessage, got %T", v)
	}
	if val.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("anthropic: Schema must be a pointer to a struct, got pointer to %s", val.Elem().Kind())
	}

	// Cache the final json.RawMessage by reflect.Type so repeated requests
	// with the same struct type skip reflection, transformation, and marshaling.
	t := val.Elem().Type()
	if cached, ok := schemaCache.Load(t); ok {
		return cached.(json.RawMessage), nil
	}

	reflector := jsonschema.Reflector{DoNotReference: true}
	schema := reflector.Reflect(val.Interface())

	// Transform the schema in-place on the typed struct, then marshal once.
	transformSchema(schema)

	result, err := json.Marshal(schema)
	if err != nil {
		return nil, fmt.Errorf("anthropic: failed to marshal JSON schema: %w", err)
	}

	raw := json.RawMessage(result)
	schemaCache.Store(t, raw)
	return raw, nil
}

// outputFormatDest checks both OutputFormat.Schema and OutputConfig.Format.Schema
// for a struct pointer and returns it. Returns (nil, false) if neither has one.
func outputFormatDest(params BetaMessageNewParams) (any, bool) {
	for _, schema := range []any{params.OutputFormat.Schema, params.OutputConfig.Format.Schema} {
		if schema == nil {
			continue
		}
		val := reflect.ValueOf(schema)
		if val.Kind() == reflect.Ptr && !val.IsNil() && val.Elem().Kind() == reflect.Struct {
			return schema, true
		}
	}
	return nil, false
}

// parseOutputContent finds the first text content block in the message
// and unmarshals it into dest. Errors are wrapped with ErrStructuredOutputParse
// so callers can distinguish parse failures via errors.Is.
func parseOutputContent(msg *BetaMessage, dest any) error {
	for _, block := range msg.Content {
		if block.Type == "text" {
			if err := json.Unmarshal([]byte(block.Text), dest); err != nil {
				return fmt.Errorf("%w: %w", ErrStructuredOutputParse, err)
			}
			return nil
		}
	}
	return fmt.Errorf("%w: no text content block found in response", ErrStructuredOutputParse)
}
