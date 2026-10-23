# Message Threads

A message thread (beta) keeps the conversation on the server, so each request sends only the new messages instead of the whole history. `client.Beta.Messages.NewThread(...)` takes what `New` takes and sends nothing until the first turn.

```go
thread := client.Beta.Messages.NewThread(anthropic.BetaMessageThreadCreateParams{
	Model:     anthropic.ModelClaudeSonnet5,
	MaxTokens: 1024,
	System:    []anthropic.BetaTextBlockParam{{Text: "You are a concise assistant."}},
})

first, err := thread.Send(ctx, anthropic.BetaMessageThreadSendParams{
	Messages: []anthropic.BetaMessageParam{anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("What is a goroutine?"))},
})
second, err := thread.Send(ctx, anthropic.BetaMessageThreadSendParams{
	Messages: []anthropic.BetaMessageParam{anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("And in one sentence?"))},
})
```

`thread.NewToolRunner(tools)` runs the [tool loop](./tools.md) over this thread.

See [examples/message-threads](./examples/message-threads).
