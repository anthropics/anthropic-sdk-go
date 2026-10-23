# Claude SDK for Go

<!-- x-release-please-start-version -->

<a href="https://pkg.go.dev/github.com/anthropics/anthropic-sdk-go/v2"><img src="https://pkg.go.dev/badge/github.com/anthropics/anthropic-sdk-go/v2.svg" alt="Go Reference"></a>

<!-- x-release-please-end -->

The Claude SDK for Go provides access to the [Claude API](https://docs.anthropic.com/en/api/) from Go applications.

## Documentation

Full documentation is available at **[platform.claude.com/docs/en/api/sdks/go](https://platform.claude.com/docs/en/api/sdks/go)**.

## Installation

<!-- x-release-please-start-version -->

```go
import (
	"github.com/anthropics/anthropic-sdk-go/v2" // imported as anthropic
)
```

<!-- x-release-please-end -->

Or explicitly add the dependency:

<!-- x-release-please-start-version -->

```sh
go get -u 'github.com/anthropics/anthropic-sdk-go/v2@v2.3.0'
```

<!-- x-release-please-end -->

## Getting started

```go
package main

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go/v2"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
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
		Model: anthropic.ModelClaudeOpus5_5,
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", message.Content)
}
```

## OpenTelemetry

The SDK reports each API call as an OpenTelemetry span when your application has registered a tracer provider.
No client option is needed for that.
The spans carry no content of requests or responses unless you set `option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)` or `ANTHROPIC_OPEN_TELEMETRY_TRACES_CONTENT_MODE=content`.
A span then holds at most 750,000 bytes of content, and has `anthropic.content.truncated` when content was cut or left out to fit; `option.WithOpenTelemetryTracesMaxContentBytes(n)` or `ANTHROPIC_OPEN_TELEMETRY_TRACES_MAX_CONTENT_BYTES=n` sets another maximum.
Each request carries the trace context of its span in the headers of your registered propagator; `option.WithOpenTelemetryPropagation(false)` or `ANTHROPIC_OPEN_TELEMETRY_PROPAGATION=false` leaves them out, for an application whose HTTP client instrumentation already sends them.

```go
client := anthropic.NewClient(option.WithTracerProvider(provider)) // a provider other than the registered one
quiet := anthropic.NewClient(option.WithoutOpenTelemetry())        // no spans; ANTHROPIC_OPEN_TELEMETRY=false does the same

client.Messages.New(ctx, params, option.WithOpenTelemetryConversationID(conversationID)) // sets gen_ai.conversation.id
```

`WithTracerProvider`, `WithOpenTelemetryPropagation`, `WithoutOpenTelemetry`, `WithOpenTelemetryTraces`, `WithOpenTelemetryTracesContentMode`, and `WithOpenTelemetryTracesMaxContentBytes` are client options.
Go has one request option type, so a method call accepts them too; passed that way, they apply to that call only.

## Requirements

Go 1.25+

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md).

## License

This project is licensed under the Apache License, Version 2.0. See the [LICENSE](LICENSE) file for details. Versions released before the change from the MIT License, and third-party contributions merged before it, remain under the MIT License. See [LICENSE-MIT](LICENSE-MIT).
