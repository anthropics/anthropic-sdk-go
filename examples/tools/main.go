package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	client := anthropic.NewClient()

	content := "What is the weather in San Francisco, CA?"

	println(color("[user]: ") + content)

	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(content)),
	}

	tools := []anthropic.ToolParam{
		{
			Name:        anthropic.F("get_coordinates"),
			Description: anthropic.F("Accepts a place as an address, then returns the latitude and longitude coordinates."),
			InputSchema: anthropic.F[interface{}](map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "The location to look up.",
					},
				},
			}),
		},
		{
			Name: anthropic.F("get_temperature_unit"),
			InputSchema: anthropic.F[interface{}](map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"country": map[string]interface{}{
						"type":        "string",
						"description": "The country",
					},
				},
			}),
		},
		{
			Name:        anthropic.F("get_weather"),
			Description: anthropic.F("Get the weather at a specific location"),
			InputSchema: anthropic.F[interface{}](map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"lat": map[string]interface{}{
						"type":        "number",
						"description": "The lattitude of the location to check weather.",
					},
					"long": map[string]interface{}{
						"type":        "number",
						"description": "The longitude of the location to check weather.",
					},
					"unit": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"celsius", "fahrenheit"},
						"description": "Unit for the output",
					},
				},
			}),
		},
	}

	for {
		message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
			Model:     anthropic.F(anthropic.ModelClaude_3_5_Sonnet_20240620),
			MaxTokens: anthropic.Int(1024),
			Messages:  anthropic.F(messages),
			Tools:     anthropic.F(tools),
		})

		if err != nil {
			panic(err)
		}

		print(color("[assistant]: "))

		for _, block := range message.Content {
			switch block := block.AsUnion().(type) {
			case anthropic.TextBlock:
				println(block.Text)
				println()
			case anthropic.ToolUseBlock:
				println(block.Name + ": " + string(block.Input))
				println()
			}
		}
		println()

		messages = append(messages, message.ToParam())
		toolResults := []anthropic.MessageParamContentUnion{}

		for _, block := range message.Content {
			if block.Type == anthropic.ContentBlockTypeToolUse {
				print(color("[user (" + block.Name + ")]: "))

				var response interface{}
				switch block.Name {
				case "get_coordinates":
					var input struct {
						Location string `json:"location"`
					}
					err := json.Unmarshal(block.Input, &input)
					if err != nil {
						panic(err)
					}
					response = GetCoordinates(input.Location)
				case "get_temperature_unit":
					var input struct {
						Country string `json:"country"`
					}
					err := json.Unmarshal(block.Input, &input)
					if err != nil {
						panic(err)
					}
					response = GetTemperatureUnit(input.Country)
				case "get_weather":
					var input struct {
						Lat  float64 `json:"lat"`
						Long float64 `json:"long"`
						Unit string  `json:"unit"`
					}
					err := json.Unmarshal(block.Input, &input)
					if err != nil {
						panic(err)
					}
					response = GetWeather(input.Lat, input.Long, input.Unit)
				}

				b, err := json.Marshal(response)
				if err != nil {
					panic(err)
				}

				println(string(b))

				toolResults = append(toolResults, anthropic.NewToolResultBlock(block.ID, string(b), false))
			}
		}

		if len(toolResults) == 0 {
			break
		}

		messages = append(messages, anthropic.NewUserMessage(toolResults...))
	}
}

type CoordinateResponse struct {
	Long float64 `json:"long"`
	Lat  float64 `json:"lat"`
}

func GetCoordinates(location string) CoordinateResponse {
	return CoordinateResponse{
		Long: -122.4194,
		Lat:  37.7749,
	}
}

func GetTemperatureUnit(country string) string {
	return "farenheit"
}

type WeatherResponse struct {
	Unit        string  `json:"unit"`
	Temperature float64 `json:"temperature"`
}

func GetWeather(lat, long float64, unit string) WeatherResponse {
	return WeatherResponse{
		Unit:        "farenheit",
		Temperature: 122,
	}
}

func color(s string) string {
	return fmt.Sprintf("\033[1;%sm%s\033[0m", "33", s)
}
