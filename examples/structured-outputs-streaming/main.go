package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

// WeatherQuery demonstrates structured output with streaming.
type WeatherQuery struct {
	Location    string   `json:"location" jsonschema:"title=Location,description=The city and state e.g. San Francisco CA"`
	Units       string   `json:"units,omitempty" jsonschema:"enum=celsius,enum=fahrenheit,default=fahrenheit"`
	Days        int      `json:"days" jsonschema:"minimum=1,maximum=7,description=Number of days to forecast"`
	IncludeWind bool     `json:"include_wind,omitempty"`
	Details     []string `json:"details,omitempty" jsonschema:"minItems=1,maxItems=5"`
}

func main() {
	client := anthropic.NewClient()

	// Pass a struct pointer as Schema — the JSON schema is auto-generated
	// on the wire. After accumulating the stream, use ParseOutput to parse.
	var weather WeatherQuery
	stream := client.Beta.Messages.NewStreaming(context.TODO(), anthropic.BetaMessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5,
		MaxTokens: 1024,
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("What's the weather like in San Francisco for the next 3 days? Include wind information.")),
		},
		OutputFormat: anthropic.BetaJSONOutputFormatParam{
			Schema: &weather,
		},
		Betas: []anthropic.AnthropicBeta{"structured-outputs-2025-11-13"},
	})

	var msg anthropic.BetaMessage
	for stream.Next() {
		evt := stream.Current()
		fmt.Printf("Event: %s\n", evt.Type)
		msg.Accumulate(evt)
	}

	if err := stream.Err(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if err := msg.ParseOutput(&weather); err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	fmt.Printf("\nStop reason: %s\n", msg.StopReason)
	fmt.Println("\nParsed response:")
	printJSON(weather)
}

func printJSON(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(string(b))
}
