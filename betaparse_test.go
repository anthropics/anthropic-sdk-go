package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/anthropics/anthropic-sdk-go/option"
)

type testOrder struct {
	Items    []testOrderItem `json:"items" jsonschema:"description=List of items in the order"`
	Total    float64         `json:"total"`
	Currency string          `json:"currency" jsonschema:"enum=USD,enum=EUR"`
}

type testOrderItem struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func TestSchemaToRaw(t *testing.T) {
	t.Run("struct pointer generates schema with additionalProperties false", func(t *testing.T) {
		order := &testOrder{}
		raw, err := schemaToRaw(order)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if raw == nil {
			t.Fatal("expected non-nil raw schema")
		}

		var schemaMap map[string]any
		if err := json.Unmarshal(raw, &schemaMap); err != nil {
			t.Fatalf("unexpected unmarshal error: %v", err)
		}

		schemaType, ok := schemaMap["type"].(string)
		if !ok || schemaType != "object" {
			t.Fatalf("expected schema type 'object', got %v", schemaMap["type"])
		}

		props, ok := schemaMap["properties"].(map[string]any)
		if !ok {
			t.Fatal("expected properties to be a map")
		}
		for _, key := range []string{"items", "total", "currency"} {
			if _, ok := props[key]; !ok {
				t.Fatalf("expected '%s' property in schema", key)
			}
		}

		if ap, ok := schemaMap["additionalProperties"]; !ok || ap != false {
			t.Fatalf("expected additionalProperties=false, got %v", ap)
		}
	})

	t.Run("struct pointer preserves map field additionalProperties schema", func(t *testing.T) {
		type payloadWithLabels struct {
			Labels map[string]string `json:"labels"`
		}

		raw, err := schemaToRaw(&payloadWithLabels{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if raw == nil {
			t.Fatal("expected non-nil raw schema")
		}

		var schemaMap map[string]any
		if err := json.Unmarshal(raw, &schemaMap); err != nil {
			t.Fatalf("unexpected unmarshal error: %v", err)
		}
		got, _ := json.Marshal(schemaMap["properties"].(map[string]any)["labels"])
		want := `{"type":"object","additionalProperties":{"type":"string"}}`
		if normalizeJSON(string(got)) != normalizeJSON(want) {
			t.Fatalf("expected labels schema %s, got %s", want, got)
		}
	})

	t.Run("map[string]any marshals to json.RawMessage", func(t *testing.T) {
		m := map[string]any{"type": "object"}
		raw, err := schemaToRaw(m)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if raw == nil {
			t.Fatal("expected non-nil raw for map[string]any input")
		}
		var result map[string]any
		if err := json.Unmarshal(raw, &result); err != nil {
			t.Fatalf("unexpected unmarshal error: %v", err)
		}
		if result["type"] != "object" {
			t.Fatalf("expected type 'object', got %v", result["type"])
		}
	})

	t.Run("json.RawMessage returned as-is", func(t *testing.T) {
		input := json.RawMessage(`{"type":"object"}`)
		raw, err := schemaToRaw(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(raw) != string(input) {
			t.Fatalf("expected raw to be returned as-is, got %s", raw)
		}
	})

	t.Run("non-pointer returns error", func(t *testing.T) {
		_, err := schemaToRaw(testOrder{})
		if err == nil {
			t.Fatal("expected error for non-pointer")
		}
	})

	t.Run("nil pointer returns error", func(t *testing.T) {
		_, err := schemaToRaw((*testOrder)(nil))
		if err == nil {
			t.Fatal("expected error for nil pointer")
		}
	})

	t.Run("pointer to non-struct returns error", func(t *testing.T) {
		s := "hello"
		_, err := schemaToRaw(&s)
		if err == nil {
			t.Fatal("expected error for pointer to string")
		}
	})
}

func clearSchemaCache() {
	schemaCache.Range(func(key, _ any) bool {
		schemaCache.Delete(key)
		return true
	})
}

func TestSchemaCaching(t *testing.T) {
	t.Run("same type returns cached result", func(t *testing.T) {
		clearSchemaCache()

		raw1, err := schemaToRaw(&testOrder{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		raw2, err := schemaToRaw(&testOrder{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(raw1) != string(raw2) {
			t.Fatalf("expected identical results, got %s and %s", raw1, raw2)
		}
	})

	t.Run("different types produce different schemas", func(t *testing.T) {
		clearSchemaCache()

		type Alpha struct {
			Name string `json:"name"`
		}
		type Beta struct {
			Count int `json:"count"`
		}

		raw1, err := schemaToRaw(&Alpha{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		raw2, err := schemaToRaw(&Beta{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(raw1) == string(raw2) {
			t.Fatal("expected different schemas for different types")
		}
	})

	t.Run("struct field values do not affect cached schema", func(t *testing.T) {
		clearSchemaCache()

		raw1, err := schemaToRaw(&testOrder{Total: 100, Currency: "USD"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		raw2, err := schemaToRaw(&testOrder{Total: 200, Currency: "EUR"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(raw1) != string(raw2) {
			t.Fatalf("expected identical schemas regardless of field values, got %s and %s", raw1, raw2)
		}
	})

	t.Run("non-struct types are not cached", func(t *testing.T) {
		clearSchemaCache()

		m1 := map[string]any{"type": "object"}
		raw1, err := schemaToRaw(m1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		m2 := map[string]any{"type": "string"}
		raw2, err := schemaToRaw(m2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(raw1) == string(raw2) {
			t.Fatal("expected different results for different maps")
		}
	})
}

func TestParseOutputContent(t *testing.T) {
	t.Run("parses text block into struct", func(t *testing.T) {
		msg := &BetaMessage{
			Content: []BetaContentBlockUnion{
				{
					Type: "text",
					Text: `{"items":[{"name":"Widget","quantity":2,"price":9.99}],"total":19.98,"currency":"USD"}`,
				},
			},
		}

		var order testOrder
		err := parseOutputContent(msg, &order)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(order.Items) != 1 || order.Items[0].Name != "Widget" {
			t.Errorf("unexpected items: %+v", order.Items)
		}
		if order.Total != 19.98 {
			t.Errorf("expected total 19.98, got %f", order.Total)
		}
		if order.Currency != "USD" {
			t.Errorf("expected currency 'USD', got '%s'", order.Currency)
		}
	})

	t.Run("skips non-text blocks", func(t *testing.T) {
		msg := &BetaMessage{
			Content: []BetaContentBlockUnion{
				{Type: "thinking", Thinking: "Let me think..."},
				{Type: "text", Text: `{"items":[],"total":0,"currency":"EUR"}`},
			},
		}

		var order testOrder
		err := parseOutputContent(msg, &order)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if order.Currency != "EUR" {
			t.Errorf("expected currency 'EUR', got '%s'", order.Currency)
		}
	})

	t.Run("returns error when no text block found", func(t *testing.T) {
		msg := &BetaMessage{
			Content: []BetaContentBlockUnion{
				{Type: "thinking", Thinking: "hmm"},
			},
		}

		var order testOrder
		err := parseOutputContent(msg, &order)
		if err == nil {
			t.Fatal("expected error for missing text block")
		}
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		msg := &BetaMessage{
			Content: []BetaContentBlockUnion{
				{Type: "text", Text: "not valid json"},
			},
		}

		var order testOrder
		err := parseOutputContent(msg, &order)
		if err == nil {
			t.Fatal("expected error for invalid JSON")
		}
		if !errors.Is(err, ErrStructuredOutputParse) {
			t.Errorf("expected error to wrap ErrStructuredOutputParse, got %v", err)
		}
	})

	t.Run("missing text block error wraps sentinel", func(t *testing.T) {
		msg := &BetaMessage{
			Content: []BetaContentBlockUnion{
				{Type: "thinking", Thinking: "hmm"},
			},
		}

		var order testOrder
		err := parseOutputContent(msg, &order)
		if !errors.Is(err, ErrStructuredOutputParse) {
			t.Errorf("expected error to wrap ErrStructuredOutputParse, got %v", err)
		}
	})
}

func TestNewAutoParseWithMockServer(t *testing.T) {
	responseJSON := `{
		"id": "msg_123",
		"type": "message",
		"role": "assistant",
		"model": "claude-sonnet-4-5-20250514",
		"stop_reason": "end_turn",
		"stop_sequence": null,
		"usage": {"input_tokens": 100, "output_tokens": 50},
		"content": [{
			"type": "text",
			"text": "{\"items\":[{\"name\":\"Laptop\",\"quantity\":1,\"price\":999.99}],\"total\":999.99,\"currency\":\"USD\"}"
		}]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request body includes the generated schema
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}

		outputFormat, ok := body["output_format"].(map[string]any)
		if !ok {
			t.Error("expected output_format in request body")
		} else {
			if outputFormat["type"] != "json_schema" {
				t.Errorf("expected type 'json_schema', got %v", outputFormat["type"])
			}
			schema, ok := outputFormat["schema"].(map[string]any)
			if !ok || schema == nil {
				t.Error("expected schema in output_format")
			} else {
				if schema["additionalProperties"] != false {
					t.Error("expected additionalProperties=false in schema")
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, responseJSON)
	}))
	defer server.Close()

	client := NewClient(
		option.WithBaseURL(server.URL),
		option.WithAPIKey("test-key"),
	)

	var order testOrder
	msg, err := client.Beta.Messages.New(context.Background(), BetaMessageNewParams{
		Model:     ModelClaudeSonnet4_5,
		MaxTokens: 1024,
		Messages: []BetaMessageParam{
			NewBetaUserMessage(NewBetaTextBlock("Order a laptop")),
		},
		OutputFormat: BetaJSONOutputFormatParam{
			Schema: &order,
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if msg.ID != "msg_123" {
		t.Errorf("expected message ID 'msg_123', got '%s'", msg.ID)
	}

	if len(order.Items) != 1 || order.Items[0].Name != "Laptop" {
		t.Fatalf("expected 1 item named 'Laptop', got %+v", order.Items)
	}
	if order.Total != 999.99 {
		t.Errorf("expected total 999.99, got %f", order.Total)
	}
}

func TestStreamingWithParseOutput(t *testing.T) {
	events := []string{
		`event: message_start`,
		`data: {"type":"message_start","message":{"id":"msg_456","type":"message","role":"assistant","model":"claude-sonnet-4-5-20250514","content":[],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":100,"output_tokens":0}}}`,
		``,
		`event: content_block_start`,
		`data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`,
		``,
		`event: content_block_delta`,
		`data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"{\"items\":[{\"name\":\"Phone\""}}`,
		``,
		`event: content_block_delta`,
		`data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":",\"quantity\":2,\"price\":499.99}],\"total\":999.98,\"currency\":\"EUR\"}"}}`,
		``,
		`event: content_block_stop`,
		`data: {"type":"content_block_stop","index":0}`,
		``,
		`event: message_delta`,
		`data: {"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null},"usage":{"output_tokens":30}}`,
		``,
		`event: message_stop`,
		`data: {"type":"message_stop"}`,
		``,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			t.Fatal("expected http.Flusher")
		}

		for _, event := range events {
			fmt.Fprintf(w, "%s\n", event)
		}
		flusher.Flush()
	}))
	defer server.Close()

	client := NewClient(
		option.WithBaseURL(server.URL),
		option.WithAPIKey("test-key"),
	)

	var order testOrder
	stream := client.Beta.Messages.NewStreaming(context.Background(), BetaMessageNewParams{
		Model:     ModelClaudeSonnet4_5,
		MaxTokens: 1024,
		Messages: []BetaMessageParam{
			NewBetaUserMessage(NewBetaTextBlock("Order a phone")),
		},
		OutputFormat: BetaJSONOutputFormatParam{
			Schema: &order,
		},
	})

	var msg BetaMessage
	for stream.Next() {
		msg.Accumulate(stream.Current())
	}

	if err := stream.Err(); err != nil {
		t.Fatalf("unexpected stream error: %v", err)
	}

	if err := msg.ParseOutput(&order); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if msg.ID != "msg_456" {
		t.Errorf("expected message ID 'msg_456', got '%s'", msg.ID)
	}

	if len(order.Items) != 1 || order.Items[0].Name != "Phone" {
		t.Fatalf("expected 1 item named 'Phone', got %+v", order.Items)
	}
	if order.Items[0].Quantity != 2 {
		t.Errorf("expected quantity 2, got %d", order.Items[0].Quantity)
	}
	if order.Total != 999.98 {
		t.Errorf("expected total 999.98, got %f", order.Total)
	}
	if order.Currency != "EUR" {
		t.Errorf("expected currency 'EUR', got '%s'", order.Currency)
	}
}

func TestMarshalJSONSerializesSchemaNotValues(t *testing.T) {
	type Order struct {
		Total float64 `json:"total"`
	}

	data, _ := json.Marshal(BetaJSONOutputFormatParam{Schema: &Order{Total: 42}})

	var got map[string]any
	json.Unmarshal(data, &got)
	schema := got["schema"].(map[string]any)

	if schema["type"] != "object" {
		t.Fatalf("expected JSON Schema, got struct values: %s", data)
	}
	if schema["additionalProperties"] != false {
		t.Fatalf("expected additionalProperties=false: %s", data)
	}
	props := schema["properties"].(map[string]any)
	total := props["total"].(map[string]any)
	if total["type"] != "number" {
		t.Fatalf("expected properties.total.type='number': %s", data)
	}
}

func TestMarshalJSONSchemaStructPointer(t *testing.T) {
	type simple struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	param := BetaJSONOutputFormatParam{
		Schema: &simple{},
	}

	data, err := json.Marshal(param)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if result["type"] != "json_schema" {
		t.Errorf("expected type 'json_schema', got %v", result["type"])
	}

	schema, ok := result["schema"].(map[string]any)
	if !ok {
		t.Fatal("expected schema to be a map")
	}

	if schema["type"] != "object" {
		t.Errorf("expected schema type 'object', got %v", schema["type"])
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("expected properties to be a map")
	}
	if _, ok := props["name"]; !ok {
		t.Error("expected 'name' property in schema")
	}
	if _, ok := props["age"]; !ok {
		t.Error("expected 'age' property in schema")
	}

	if schema["additionalProperties"] != false {
		t.Error("expected additionalProperties=false")
	}
}

func TestMarshalJSONSchemaMapPassthrough(t *testing.T) {
	// Backward compat: map[string]any as Schema still works
	schemaMap := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{"type": "string"},
		},
	}

	param := BetaJSONOutputFormatParam{
		Schema: schemaMap,
	}

	data, err := json.Marshal(param)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	schema, ok := result["schema"].(map[string]any)
	if !ok {
		t.Fatal("expected schema to be a map")
	}
	if schema["type"] != "object" {
		t.Errorf("expected schema type 'object', got %v", schema["type"])
	}
}
