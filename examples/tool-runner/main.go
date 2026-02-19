// Example demonstrating the Tool Runner framework
package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/toolrunner"
)

type CalculatorInput struct {
	Operation string  `json:"operation" jsonschema:"required,description=The arithmetic operation to perform,enum=add,subtract,multiply,divide"`
	A         float64 `json:"a" jsonschema:"required,description=The first number"`
	B         float64 `json:"b" jsonschema:"required,description=The second number"`
}

func calculate(ctx context.Context, calc CalculatorInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
	var floatResult float64
	fmt.Printf("🔧 Calculator tool called with: %+v\n", calc)
	switch calc.Operation {
	case "add":
		floatResult = calc.A + calc.B
	case "subtract":
		floatResult = calc.A - calc.B
	case "multiply":
		floatResult = calc.A * calc.B
	case "divide":
		if calc.B == 0 {
			return anthropic.BetaToolResultBlockParamContentUnion{}, fmt.Errorf("division by zero")
		}
		floatResult = calc.A / calc.B
	default:
		return anthropic.BetaToolResultBlockParamContentUnion{}, fmt.Errorf("unknown operation: %s", calc.Operation)
	}

	return anthropic.BetaToolResultBlockParamContentUnion{
		OfText: &anthropic.BetaTextBlockParam{Text: strconv.FormatFloat(floatResult, 'g', -1, 64)},
	}, nil
}

func main() {
	client := anthropic.NewClient()
	ctx := context.Background()

	calculatorTool, err := toolrunner.NewBetaToolFromJSONSchema("calculator", "Perform basic arithmetic operations", calculate)
	if err != nil {
		fmt.Printf("Error creating calculator tool: %v\n", err)
		return
	}

	fmt.Printf("Starting tool runner with calculator tool: %+v\n", calculatorTool)

	tools := []anthropic.BetaTool{calculatorTool}

	runner := client.Beta.Messages.NewToolRunner(tools, anthropic.BetaToolRunnerParams{
		BetaMessageNewParams: anthropic.BetaMessageNewParams{
			Model:     anthropic.ModelClaudeSonnet4_20250514,
			MaxTokens: 1000,
			Messages: []anthropic.BetaMessageParam{
				anthropic.NewBetaUserMessage(
					anthropic.NewBetaTextBlock("Calculate 15 * 23, then add 10 to the result"),
				),
			},
		},
		MaxIterations: 5,
	})

	finalMessage, err := runner.RunToCompletion(ctx)
	if err != nil {
		fmt.Printf("Error running tools: %v\n", err)
		return
	}

	fmt.Printf("Final message content:\n")
	fmt.Printf("%+v\n", finalMessage)
}
