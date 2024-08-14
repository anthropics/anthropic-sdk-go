package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
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
			InputSchema: anthropic.F(GetCoordinatesInputSchema),
		},
		{
			Name:        anthropic.F("get_temperature_unit"),
			InputSchema: anthropic.F(GetTemperatureUnitInputSchema),
		},
		{
			Name:        anthropic.F("get_weather"),
			Description: anthropic.F("Get the weather at a specific location"),
			InputSchema: anthropic.F(GetWeatherInputSchema),
		},
	}

	for {
		stream := client.Messages.NewStreaming(context.TODO(), anthropic.MessageNewParams{
			Model:     anthropic.F(anthropic.ModelClaude_3_5_Sonnet_20240620),
			MaxTokens: anthropic.Int(1024),
			Messages:  anthropic.F(messages),
			Tools:     anthropic.F(tools),
		})

		print(color("[assistant]: "))

		message := anthropic.Message{}
		for stream.Next() {
			event := stream.Current()
			err := message.Accumulate(event)
			if err != nil {
				panic(err)
			}

			switch event := event.AsUnion().(type) {
			case anthropic.ContentBlockStartEvent:
				if event.ContentBlock.Name != "" {
					print(event.ContentBlock.Name + ": ")
				}
			case anthropic.ContentBlockDeltaEvent:
				print(event.Delta.Text)
				print(event.Delta.PartialJSON)
			case anthropic.ContentBlockStopEvent:
				println()
				println()
			case anthropic.MessageStopEvent:
				println()
			}
		}

		if stream.Err() != nil {
			panic(stream.Err())
		}

		messages = append(messages, message.ToParam())
		toolResults := []anthropic.MessageParamContentUnion{}

		for _, block := range message.Content {
			if block.Type == anthropic.ContentBlockTypeToolUse {
				print(color("[user (" + block.Name + ")]: "))

				var response interface{}
				switch block.Name {
				case "get_coordinates":
					input := GetCoordinatesInput{}
					err := json.Unmarshal(block.Input, &input)
					if err != nil {
						panic(err)
					}
					response = GetCoordinates(input.Location)
				case "get_temperature_unit":
					input := GetTemperatureUnitInput{}
					err := json.Unmarshal(block.Input, &input)
					if err != nil {
						panic(err)
					}
					response = GetTemperatureUnit(input.Country)
				case "get_weather":
					input := GetWeatherInput{}
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

		messages = append(messages, anthropic.MessageParam{
			Role:    anthropic.F(anthropic.MessageParamRoleUser),
			Content: anthropic.F(toolResults),
		})
	}
}

// Get Coordinates

type GetCoordinatesInput struct {
	Location string `json:"location" jsonschema_description:"The location to look up."`
}

var GetCoordinatesInputSchema = GenerateSchema[GetCoordinatesInput]()

type GetCoordinateResponse struct {
	Long float64 `json:"long"`
	Lat  float64 `json:"lat"`
}

func GetCoordinates(location string) GetCoordinateResponse {
	return GetCoordinateResponse{
		Long: -122.4194,
		Lat:  37.7749,
	}
}

// Get Temperature Unit

type GetTemperatureUnitInput struct {
	Country string `json:"country" jsonschema_description:"The country"`
}

var GetTemperatureUnitInputSchema = GenerateSchema[GetTemperatureUnitInput]()

func GetTemperatureUnit(country string) string {
	return "farenheit"
}

// Get Weather

type GetWeatherInput struct {
	Lat  float64 `json:"lat" jsonschema_description:"The latitude of the location to check weather."`
	Long float64 `json:"long" jsonschema_description:"The longitude of the location to check weather."`
	Unit string  `json:"unit" jsonschema_description:"Unit for the output"`
}

var GetWeatherInputSchema = GenerateSchema[GetWeatherInput]()

type GetWeatherResponse struct {
	Unit        string  `json:"unit"`
	Temperature float64 `json:"temperature"`
}

func GetWeather(lat, long float64, unit string) GetWeatherResponse {
	return GetWeatherResponse{
		Unit:        "farenheit",
		Temperature: 122,
	}
}

func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	return reflector.Reflect(v)
}

func color(s string) string {
	return fmt.Sprintf("\033[1;%sm%s\033[0m", "33", s)
}
