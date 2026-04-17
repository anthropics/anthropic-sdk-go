package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	client := anthropic.NewClient()

	// Use json.RawMessage when you already have a JSON schema as raw bytes.
	// Unlike map[string]any, json.RawMessage serializes deterministically,
	// so prompt caching works correctly.
	schema := json.RawMessage(`{
		"type": "object",
		"properties": {
			"location": {
				"type": "string",
				"description": "The city and state, e.g. San Francisco, CA"
			},
			"units": {
				"type": "string",
				"description": "Temperature units",
				"enum": ["celsius", "fahrenheit"]
			},
			"days": {
				"type": "integer",
				"description": "Number of days to forecast"
			}
		},
		"required": ["location", "units", "days"],
		"additionalProperties": false
	}`)

	msg, err := client.Beta.Messages.New(context.TODO(), anthropic.BetaMessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5,
		MaxTokens: 1024,
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("What's the weather like in San Francisco for the next 3 days?")),
		},
		OutputFormat: anthropic.BetaJSONOutputFormatParam{
			Schema: schema,
		},
		Betas: []anthropic.AnthropicBeta{"structured-outputs-2025-11-13"},
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Stop reason: %s\n", msg.StopReason)
	for _, block := range msg.Content {
		if block.Type == "text" {
			// Parse manually since we used a raw schema (no struct pointer for auto-parse)
			var result map[string]any
			if err := json.Unmarshal([]byte(block.Text), &result); err != nil {
				fmt.Printf("Parse error: %v\n", err)
				return
			}
			printJSON(result)
		}
	}
}

func printJSON(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(string(b))
}
