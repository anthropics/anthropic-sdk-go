// Example demonstrating the Tool Runner framework with streaming
package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/toolrunner"
)

type WeatherRequest struct {
	City    string `json:"city" jsonschema:"required,description=The city to get weather for"`
	Units   string `json:"units,omitempty" jsonschema:"description=Temperature units,enum=celsius,fahrenheit"`
	Include string `json:"include,omitempty" jsonschema:"description=Additional data to include,enum=forecast,historical"`
}

func colorWith(code string, s string) string { return fmt.Sprintf("\033[1;%sm%s\033[0m", code, s) }
func colorUser(s string) string              { return colorWith("36", s) } // cyan
func colorAssistant(s string) string         { return colorWith("32", s) } // green
func colorTool(s string) string              { return colorWith("33", s) } // yellow
func colorThinking(s string) string          { return colorWith("90", s) } // grey

func getWeather(ctx context.Context, req WeatherRequest) (anthropic.BetaToolResultBlockParamContentUnion, error) {
	fmt.Printf("%s%s %+v\n", colorTool("[tool get_weather]: "), "called with", req)

	temp := 22
	if req.Units == "fahrenheit" {
		temp = 72
	}

	return anthropic.BetaToolResultBlockParamContentUnion{
		OfText: &anthropic.BetaTextBlockParam{
			Text: fmt.Sprintf("The current weather in %s is %d degrees %s. Tomorrow's weather will be cloudy and colder.", req.City, temp, req.Units),
		},
	}, nil
}

func main() {
	client := anthropic.NewClient()
	ctx := context.Background()

	weatherTool, err := toolrunner.NewBetaToolFromJSONSchema("get_weather", "Get current weather information for a city", getWeather)
	if err != nil {
		fmt.Printf("Error creating weather tool: %v\n", err)
		return
	}

	userQuestion := "What's the weather like in San Francisco? Please use Fahrenheit and include the forecast."
	fmt.Println(colorUser("[user]: ") + userQuestion)

	tools := []anthropic.BetaTool{weatherTool}

	runner := client.Beta.Messages.NewToolRunnerStreaming(tools, anthropic.BetaToolRunnerParams{
		BetaMessageNewParams: anthropic.BetaMessageNewParams{
			Model:     anthropic.ModelClaudeSonnet4_20250514,
			MaxTokens: 1000,
			Messages: []anthropic.BetaMessageParam{
				anthropic.NewBetaUserMessage(
					anthropic.NewBetaTextBlock(userQuestion),
				),
			},
		},
		MaxIterations: 5,
	})

	for eventsIterator := range runner.AllStreaming(ctx) {
		_ = runner.IterationCount()
		for event, err := range eventsIterator {
			if err != nil {
				fmt.Printf("%s %+v\n", colorWith("31", "[error]: "+err.Error()), event)
				return
			}
			switch eventVariant := event.AsAny().(type) {
			case anthropic.BetaRawMessageStartEvent:
				fmt.Print(colorAssistant("[assistant]: "))
			case anthropic.BetaRawContentBlockStartEvent:
				switch cb := eventVariant.ContentBlock.AsAny().(type) {
				case anthropic.BetaToolUseBlock:
					// Assistant is initiating a tool call; stream its JSON input deltas next
					label := fmt.Sprintf("[tool call %s]: ", cb.Name)
					fmt.Print(colorTool(label))
				case anthropic.BetaTextBlock:
					// nothing, normal assistant text will follow via deltas
				case anthropic.BetaThinkingBlock:
					fmt.Print(colorThinking("[assistant thinking]: "))
				}
			case anthropic.BetaRawContentBlockDeltaEvent:
				switch deltaVariant := eventVariant.Delta.AsAny().(type) {
				case anthropic.BetaTextDelta:
					fmt.Print(colorAssistant(deltaVariant.Text))
				case anthropic.BetaInputJSONDelta:
					if deltaVariant.PartialJSON != "" {
						fmt.Print(colorTool(deltaVariant.PartialJSON))
					}
				case anthropic.BetaThinkingDelta:
					fmt.Print(colorThinking(deltaVariant.Thinking))
				}
			case anthropic.BetaRawContentBlockStopEvent:
				fmt.Println()
			case anthropic.BetaRawMessageDeltaEvent:
				// No visible text here; keep for completeness
			case anthropic.BetaRawMessageStopEvent:
				fmt.Println()
			}
		}
	}
}
