# Anthropic Go API Library

<a href="https://pkg.go.dev/github.com/anthropics/anthropic-sdk-go"><img src="https://pkg.go.dev/badge/github.com/anthropics/anthropic-sdk-go.svg" alt="Go Reference"></a>

The Anthropic Go library provides access to the [Claude API](https://docs.anthropic.com/en/api/) from Go applications.

## Documentation

Full documentation is available at **[docs.anthropic.com/en/api/sdks/go](https://docs.anthropic.com/en/api/sdks/go)**.

## Installation

```go
import (
	"github.com/anthropics/anthropic-sdk-go" // imported as anthropic
)
```

Or explicitly add the dependency:

```sh
go get -u github.com/anthropics/anthropic-sdk-go
```

## Getting started

```go
package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func main() {
	client := anthropic.NewClient(
		option.WithAPIKey("my-anthropic-api-key"), // defaults to os.LookupEnv("ANTHROPIC_API_KEY")
	)
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What is a quaternion?")),
		},
		Model: anthropic.ModelClaudeOpus4_6,
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", message.Content)
}
```

## Requirements

Go 1.22+

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md).

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
