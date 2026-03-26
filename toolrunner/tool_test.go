package toolrunner_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/toolrunner"
)

// TestSchemaValidation verifies that the tool runner validates inputs
// against the JSON Schema before executing the handler. This prevents missing
// required fields, enum violations, and type mismatches from reaching handlers.
func TestSchemaValidation(t *testing.T) {
	t.Parallel()

	type StrictInput struct {
		City  string `json:"city"`
		Units string `json:"units,omitempty"`
	}

	weatherSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"city":  map[string]any{"type": "string"},
			"units": map[string]any{"type": "string", "enum": []string{"celsius", "fahrenheit"}},
		},
		"required": []string{"city"},
	}

	handlerCalled := false
	tool, err := toolrunner.NewBetaToolFromBytes("weather", "Get weather", mustMarshal(t, weatherSchema),
		func(ctx context.Context, input StrictInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			handlerCalled = true
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: fmt.Sprintf("Weather in %s (%s)", input.City, input.Units)},
			}, nil
		})
	if err != nil {
		t.Fatalf("create tool: %v", err)
	}

	t.Run("valid input passes validation", func(t *testing.T) {
		handlerCalled = false
		input := json.RawMessage(`{"city": "London", "units": "celsius"}`)
		result, err := tool.Execute(context.Background(), input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !handlerCalled {
			t.Fatal("handler was not called for valid input")
		}
		if result.OfText == nil || result.OfText.Text != "Weather in London (celsius)" {
			t.Fatalf("unexpected result: %+v", result)
		}
	})

	t.Run("missing required field rejected", func(t *testing.T) {
		handlerCalled = false
		input := json.RawMessage(`{"units": "celsius"}`)
		_, err := tool.Execute(context.Background(), input)
		if err == nil {
			t.Fatal("expected error for missing required field 'city', got nil")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called when schema validation fails")
		}
		if !strings.Contains(err.Error(), "schema validation failed") {
			t.Fatalf("error should mention schema validation, got: %v", err)
		}
	})

	t.Run("enum violation rejected", func(t *testing.T) {
		handlerCalled = false
		input := json.RawMessage(`{"city": "London", "units": "kelvin"}`)
		_, err := tool.Execute(context.Background(), input)
		if err == nil {
			t.Fatal("expected error for enum violation on 'units', got nil")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called when schema validation fails")
		}
		if !strings.Contains(err.Error(), "schema validation failed") {
			t.Fatalf("error should mention schema validation, got: %v", err)
		}
	})

	t.Run("wrong type rejected", func(t *testing.T) {
		handlerCalled = false
		input := json.RawMessage(`{"city": 12345}`)
		_, err := tool.Execute(context.Background(), input)
		if err == nil {
			t.Fatal("expected error for wrong type on 'city', got nil")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called when schema validation fails")
		}
	})

	t.Run("empty object rejected when required fields exist", func(t *testing.T) {
		handlerCalled = false
		input := json.RawMessage(`{}`)
		_, err := tool.Execute(context.Background(), input)
		if err == nil {
			t.Fatal("expected error for empty object with required fields, got nil")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called when schema validation fails")
		}
	})

	t.Run("optional field can be omitted", func(t *testing.T) {
		handlerCalled = false
		input := json.RawMessage(`{"city": "Tokyo"}`)
		_, err := tool.Execute(context.Background(), input)
		if err != nil {
			t.Fatalf("unexpected error for valid input without optional field: %v", err)
		}
		if !handlerCalled {
			t.Fatal("handler was not called for valid input")
		}
	})
}

// TestAdditionalPropertiesRejected verifies that additionalProperties:false
// blocks unknown keys from reaching the handler.
func TestAdditionalPropertiesRejected(t *testing.T) {
	t.Parallel()

	type StrictInput struct {
		Name string `json:"name"`
	}

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{"type": "string"},
		},
		"required":             []string{"name"},
		"additionalProperties": false,
	}

	handlerCalled := false
	tool, err := toolrunner.NewBetaToolFromBytes("strict", "Strict tool", mustMarshal(t, schema),
		func(ctx context.Context, input StrictInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			handlerCalled = true
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: "ok"},
			}, nil
		})
	if err != nil {
		t.Fatalf("create tool: %v", err)
	}

	t.Run("valid input accepted", func(t *testing.T) {
		handlerCalled = false
		input := json.RawMessage(`{"name": "test"}`)
		_, err := tool.Execute(context.Background(), input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !handlerCalled {
			t.Fatal("handler was not called")
		}
	})

	t.Run("extra property rejected", func(t *testing.T) {
		handlerCalled = false
		input := json.RawMessage(`{"name": "test", "extra": "x"}`)
		_, err := tool.Execute(context.Background(), input)
		if err == nil {
			t.Fatal("expected error for additional property, got nil")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called when additionalProperties is violated")
		}
		if !strings.Contains(err.Error(), "additional property") {
			t.Fatalf("error should mention additional property, got: %v", err)
		}
	})
}

// TestPatternValidation verifies that pattern constraints on string
// properties are enforced at runtime.
func TestPatternValidation(t *testing.T) {
	t.Parallel()

	type URLInput struct {
		URL string `json:"url"`
	}

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"url": map[string]any{
				"type":    "string",
				"pattern": `^https://allowed\.example/`,
			},
		},
		"required": []string{"url"},
	}

	handlerCalled := false
	tool, err := toolrunner.NewBetaToolFromBytes("url_tool", "URL tool", mustMarshal(t, schema),
		func(ctx context.Context, input URLInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			handlerCalled = true
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: "ok"},
			}, nil
		})
	if err != nil {
		t.Fatalf("create tool: %v", err)
	}

	t.Run("matching pattern accepted", func(t *testing.T) {
		handlerCalled = false
		input := json.RawMessage(`{"url": "https://allowed.example/page"}`)
		_, err := tool.Execute(context.Background(), input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !handlerCalled {
			t.Fatal("handler was not called")
		}
	})

	t.Run("non-matching pattern rejected", func(t *testing.T) {
		handlerCalled = false
		input := json.RawMessage(`{"url": "https://evil.example/attack"}`)
		_, err := tool.Execute(context.Background(), input)
		if err == nil {
			t.Fatal("expected error for pattern violation, got nil")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called when pattern is violated")
		}
		if !strings.Contains(err.Error(), "pattern") {
			t.Fatalf("error should mention pattern, got: %v", err)
		}
	})
}

// TestStringLengthValidation verifies minLength and maxLength enforcement.
func TestStringLengthValidation(t *testing.T) {
	t.Parallel()

	type NameInput struct {
		Name string `json:"name"`
	}

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type":      "string",
				"minLength": 2,
				"maxLength": 10,
			},
		},
		"required": []string{"name"},
	}

	handlerCalled := false
	tool, err := toolrunner.NewBetaToolFromBytes("name_tool", "Name tool", mustMarshal(t, schema),
		func(ctx context.Context, input NameInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			handlerCalled = true
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: "ok"},
			}, nil
		})
	if err != nil {
		t.Fatalf("create tool: %v", err)
	}

	t.Run("valid length accepted", func(t *testing.T) {
		handlerCalled = false
		_, err := tool.Execute(context.Background(), json.RawMessage(`{"name": "Alice"}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !handlerCalled {
			t.Fatal("handler was not called")
		}
	})

	t.Run("too short rejected", func(t *testing.T) {
		handlerCalled = false
		_, err := tool.Execute(context.Background(), json.RawMessage(`{"name": "A"}`))
		if err == nil {
			t.Fatal("expected error for minLength violation")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called")
		}
	})

	t.Run("too long rejected", func(t *testing.T) {
		handlerCalled = false
		_, err := tool.Execute(context.Background(), json.RawMessage(`{"name": "VeryLongNameHere"}`))
		if err == nil {
			t.Fatal("expected error for maxLength violation")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called")
		}
	})
}

// TestNumericBoundsValidation verifies minimum and maximum enforcement.
func TestNumericBoundsValidation(t *testing.T) {
	t.Parallel()

	type AgeInput struct {
		Age int `json:"age"`
	}

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"age": map[string]any{
				"type":    "integer",
				"minimum": 0,
				"maximum": 150,
			},
		},
		"required": []string{"age"},
	}

	handlerCalled := false
	tool, err := toolrunner.NewBetaToolFromBytes("age_tool", "Age tool", mustMarshal(t, schema),
		func(ctx context.Context, input AgeInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			handlerCalled = true
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: "ok"},
			}, nil
		})
	if err != nil {
		t.Fatalf("create tool: %v", err)
	}

	t.Run("valid value accepted", func(t *testing.T) {
		handlerCalled = false
		_, err := tool.Execute(context.Background(), json.RawMessage(`{"age": 25}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !handlerCalled {
			t.Fatal("handler was not called")
		}
	})

	t.Run("below minimum rejected", func(t *testing.T) {
		handlerCalled = false
		_, err := tool.Execute(context.Background(), json.RawMessage(`{"age": -1}`))
		if err == nil {
			t.Fatal("expected error for minimum violation")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called")
		}
	})

	t.Run("above maximum rejected", func(t *testing.T) {
		handlerCalled = false
		_, err := tool.Execute(context.Background(), json.RawMessage(`{"age": 200}`))
		if err == nil {
			t.Fatal("expected error for maximum violation")
		}
		if handlerCalled {
			t.Fatal("handler should NOT be called")
		}
	})
}

// TestMissingTypeInference verifies that schemas without an explicit "type" field
// are still treated as object schemas when they contain object-specific keywords.
func TestMissingTypeInference(t *testing.T) {
	t.Parallel()

	type Input struct {
		Name string `json:"name"`
	}
	handler := func(ctx context.Context, input Input) (anthropic.BetaToolResultBlockParamContentUnion, error) {
		return anthropic.BetaToolResultBlockParamContentUnion{
			OfText: &anthropic.BetaTextBlockParam{Text: "ok"},
		}, nil
	}

	t.Run("no type with required only", func(t *testing.T) {
		schema := map[string]any{
			"required": []string{"name"},
		}
		tool, err := toolrunner.NewBetaToolFromBytes("t", "t", mustMarshal(t, schema), handler)
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		_, err = tool.Execute(context.Background(), json.RawMessage(`{}`))
		if err == nil {
			t.Fatal("expected error for missing required field in schema without type")
		}
	})

	t.Run("no type with properties only", func(t *testing.T) {
		schema := map[string]any{
			"properties": map[string]any{
				"name": map[string]any{"type": "string"},
			},
		}
		tool, err := toolrunner.NewBetaToolFromBytes("t", "t", mustMarshal(t, schema), handler)
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		// Should validate type even without top-level "type"
		_, err = tool.Execute(context.Background(), json.RawMessage(`{"name": 123}`))
		if err == nil {
			t.Fatal("expected type error in schema without type field")
		}
	})

	t.Run("no type with additionalProperties false only", func(t *testing.T) {
		schema := map[string]any{
			"additionalProperties": false,
		}
		tool, err := toolrunner.NewBetaToolFromBytes("t", "t", mustMarshal(t, schema), handler)
		if err != nil {
			t.Fatalf("create: %v", err)
		}
		_, err = tool.Execute(context.Background(), json.RawMessage(`{"x": 1}`))
		if err == nil {
			t.Fatal("expected error for additional property in schema with only additionalProperties:false")
		}
		if !strings.Contains(err.Error(), "additional property") {
			t.Fatalf("error should mention additional property, got: %v", err)
		}
	})
}

// TestAdditionalPropertiesNoPropsField verifies that additionalProperties:false
// rejects all keys when the properties field is absent entirely (not just empty).
func TestAdditionalPropertiesNoPropsField(t *testing.T) {
	t.Parallel()

	type Input struct{}
	schema := map[string]any{
		"type":                 "object",
		"additionalProperties": false,
	}
	tool, err := toolrunner.NewBetaToolFromBytes("t", "t", mustMarshal(t, schema),
		func(ctx context.Context, input Input) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: "ok"},
			}, nil
		})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	t.Run("empty object accepted", func(t *testing.T) {
		_, err := tool.Execute(context.Background(), json.RawMessage(`{}`))
		if err != nil {
			t.Fatalf("unexpected error for empty object: %v", err)
		}
	})

	t.Run("any key rejected", func(t *testing.T) {
		_, err := tool.Execute(context.Background(), json.RawMessage(`{"x": 1}`))
		if err == nil {
			t.Fatal("expected error for additional property with no properties defined")
		}
		if !strings.Contains(err.Error(), "additional property") {
			t.Fatalf("error should mention additional property, got: %v", err)
		}
	})
}

// TestEnumCrossTypeMismatch verifies that enum matching is type-strict:
// string "1" must not match numeric enum value 1.
func TestEnumCrossTypeMismatch(t *testing.T) {
	t.Parallel()

	type Input struct {
		Code any `json:"code"`
	}
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"code": map[string]any{"enum": []any{1, 2, 3}},
		},
		"required": []string{"code"},
	}
	tool, err := toolrunner.NewBetaToolFromBytes("t", "t", mustMarshal(t, schema),
		func(ctx context.Context, input Input) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: "ok"},
			}, nil
		})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	t.Run("numeric value matches numeric enum", func(t *testing.T) {
		_, err := tool.Execute(context.Background(), json.RawMessage(`{"code": 1}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("string value does not match numeric enum", func(t *testing.T) {
		_, err := tool.Execute(context.Background(), json.RawMessage(`{"code": "1"}`))
		if err == nil {
			t.Fatal("expected error: string '1' should not match numeric enum value 1")
		}
	})
}

// TestInvalidRegexPattern verifies that an invalid regex pattern in the schema
// causes validation to fail with an error instead of silently passing.
func TestInvalidRegexPattern(t *testing.T) {
	t.Parallel()

	type Input struct {
		Value string `json:"value"`
	}
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"value": map[string]any{
				"type":    "string",
				"pattern": "[invalid(regex",
			},
		},
		"required": []string{"value"},
	}
	tool, err := toolrunner.NewBetaToolFromBytes("t", "t", mustMarshal(t, schema),
		func(ctx context.Context, input Input) (anthropic.BetaToolResultBlockParamContentUnion, error) {
			return anthropic.BetaToolResultBlockParamContentUnion{
				OfText: &anthropic.BetaTextBlockParam{Text: "ok"},
			}, nil
		})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	_, err = tool.Execute(context.Background(), json.RawMessage(`{"value": "anything"}`))
	if err == nil {
		t.Fatal("expected error for invalid regex pattern, got nil")
	}
	if !strings.Contains(err.Error(), "invalid pattern") {
		t.Fatalf("error should mention invalid pattern, got: %v", err)
	}
}

// mustMarshal is a test helper that marshals a value to JSON bytes or fails the test.
func mustMarshal(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}
