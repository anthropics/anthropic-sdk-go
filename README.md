# Claude SDK for Go

<!-- x-release-please-start-version -->

<a href="https://pkg.go.dev/github.com/anthropics/anthropic-sdk-go"><img src="https://pkg.go.dev/badge/github.com/anthropics/anthropic-sdk-go.svg" alt="Go Reference"></a>

<!-- x-release-please-end -->

The Claude SDK for Go provides access to the [Claude API](https://docs.anthropic.com/en/api/) from Go applications.

## Documentation

Full documentation is available at **[platform.claude.com/docs/en/api/sdks/go](https://platform.claude.com/docs/en/api/sdks/go)**.

## Installation

<!-- x-release-please-start-version -->

```go
import (
	"github.com/anthropics/anthropic-sdk-go" // imported as anthropic
)
```

<!-- x-release-please-end -->

Or explicitly add the dependency:

<!-- x-release-please-start-version -->

```sh
go get -u 'github.com/anthropics/anthropic-sdk-go@v1.38.0'
```

<!-- x-release-please-end -->

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
