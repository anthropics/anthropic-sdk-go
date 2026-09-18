# Tool Helpers

The SDK provides helper functions for defining tools and automatically running the conversation loop between Claude and your tools until Claude produces a final response.

## Defining Tools

Use the `toolrunner` package to create a `BetaTool` that combines the tool definition with its execution handler. There are three ways to create tools:

- `NewBetaToolFromJSONSchema` - Automatically generates schema from a struct with `jsonschema` tags (recommended)
- `NewBetaToolFromBytes` - Creates a tool from JSON schema bytes
- `NewBetaTool` - Creates a tool from an explicit `BetaToolInputSchemaParam`

The generic type parameter is automatically inferred from your handler function's signature, so you don't need to specify it explicitly.

### Automatic Schema Generation from Structs (Recommended)

The easiest approach is to use `NewBetaToolFromJSONSchema`, which automatically generates the schema from your struct using `jsonschema` tags:

```go
type GetWeatherInput struct {
	City  string `json:"city" jsonschema:"required,description=The city name"`
	Units string `json:"units,omitempty" jsonschema:"enum=celsius,enum=fahrenheit,description=Temperature units"`
}

weatherTool, err := toolrunner.NewBetaToolFromJSONSchema(
	"get_weather",
	"Get current weather for a city",
	func(ctx context.Context, input GetWeatherInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
		return anthropic.BetaToolResultBlockParamContentUnion{
			OfText: &anthropic.BetaTextBlockParam{
				Text: fmt.Sprintf("Weather in %s: 72°F, sunny", input.City),
			},
		}, nil
	},
)
```

### Using JSON Bytes

You can provide the schema as JSON bytes using `NewBetaToolFromBytes`:

```go
type GetWeatherInput struct {
	City string `json:"city"`
}

weatherTool, err := toolrunner.NewBetaToolFromBytes(
	"get_weather",
	"Get current weather for a city",
	[]byte(`{
		"type": "object",
		"properties": {
			"city": {"type": "string", "description": "The city name"}
		},
		"required": ["city"]
	}`),
	func(ctx context.Context, input GetWeatherInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
		// Your handler here
	},
)
```

### Using an Explicit Schema

For full control, use `NewBetaTool` with a `BetaToolInputSchemaParam` directly:

```go
weatherTool := toolrunner.NewBetaTool(
	"get_weather",
	"Get current weather for a city",
	anthropic.BetaToolInputSchemaParam{
		Properties: map[string]any{
			"city": map[string]any{
				"type":        "string",
				"description": "The city name",
			},
		},
	},
	handler,
)
```

### Raw JSON Input

If you prefer to handle JSON parsing yourself, use `json.RawMessage` or `[]byte` as the input type:

```go
rawTool, err := toolrunner.NewBetaToolFromBytes(
	"process_data",
	"Process raw JSON data",
	schemaBytes,
	func(ctx context.Context, input json.RawMessage) (anthropic.BetaToolResultBlockParamContentUnion, error) {
		// Parse the JSON yourself
		var data map[string]any
		json.Unmarshal(input, &data)
		// ...
	},
)

## Tool Runner

The `BetaToolRunner` automatically handles the conversation loop between Claude and your tools. On each iteration, it:

1. Sends the current messages to Claude
2. If the message's `StopReason` is `tool_use`, executes the requested tools in parallel and adds the results to the conversation
3. If it is `pause_turn` or `compaction`, sends the turn back unchanged so the server can resume it
4. Repeats until Claude stops for any other reason (e.g. `end_turn`, `max_tokens`), returning that final message without running tools

### Basic Usage

```go
tools := []anthropic.BetaTool{weatherTool}

runner := client.Beta.Messages.NewToolRunner(tools, anthropic.BetaToolRunnerParams{
	BetaMessageNewParams: anthropic.BetaMessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		MaxTokens: 1024,
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("What's the weather in Tokyo?")),
		},
	},
})

// Run the entire conversation to completion
message, err := runner.RunToCompletion(context.Background())
```

### Iterating Over Messages

Use `All()` to iterate over each message in the conversation:

```go
for message, err := range runner.All(ctx) {
	if err != nil {
		log.Fatal(err)
	}
	for _, block := range message.Content {
		switch b := block.AsAny().(type) {
		case anthropic.BetaTextBlock:
			fmt.Println("[assistant]:", b.Text)
		case anthropic.BetaToolUseBlock:
			fmt.Printf("[tool call]: %s(%v)\n", b.Name, b.Input)
		}
	}
}
```

### Step-by-Step Iteration

For more control, use `NextMessage()` to advance one turn at a time:

```go
for {
	message, err := runner.NextMessage(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if message == nil {
		break // Conversation complete
	}
	// Process the message...
}
```

### Streaming

Use `BetaToolRunnerStreaming` via `NewToolRunnerStreaming()` for streaming responses:

```go
runner := client.Beta.Messages.NewToolRunnerStreaming(tools, anthropic.BetaToolRunnerParams{
	BetaMessageNewParams: anthropic.BetaMessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		MaxTokens: 1024,
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("What's the weather in Tokyo?")),
		},
	},
})

for eventsIterator := range runner.AllStreaming(ctx) {
	for event, err := range eventsIterator {
		if err != nil {
			log.Fatal(err)
		}
		switch e := event.AsAny().(type) {
		case anthropic.BetaRawContentBlockDeltaEvent:
			switch delta := e.Delta.AsAny().(type) {
			case anthropic.BetaTextDelta:
				fmt.Print(delta.Text)
			}
		}
	}
}
```

Or use `NextStreaming()` for step-by-step streaming:

```go
for !runner.IsCompleted() {
	for event, err := range runner.NextStreaming(ctx) {
		// Handle streaming events...
	}
}
```

## Configuration

### Max Iterations

Limit the number of API calls to prevent runaway loops. When set to 0 (the default), there is no limit and the runner continues until the model's stop reason ends the loop:

```go
runner := client.Beta.Messages.NewToolRunner(tools, anthropic.BetaToolRunnerParams{
	// ...
	MaxIterations: 10, // Stop after 10 API calls (0 = no limit)
})
```

### Modifying Parameters Mid-Conversation

The `Params` field is exported, so you can modify parameters directly:

```go
// Update maximum tokens
runner.Params.MaxTokens = 2048

// Update maximum iterations
runner.Params.MaxIterations = 10

// Update system prompt
runner.Params.System = []anthropic.BetaTextBlockParam{
	{Text: "You are a helpful assistant."},
}

// Add messages to the conversation (direct field access)
runner.Params.Messages = append(runner.Params.Messages, anthropic.NewBetaUserMessage(
	anthropic.NewBetaTextBlock("Now check the weather in London too"),
))

// Or use the convenience method
runner.AppendMessages(anthropic.NewBetaUserMessage(
	anthropic.NewBetaTextBlock("Now check the weather in London too"),
))
```

### Compacting the Conversation

With the `compact-2026-09-04` beta you decide when a conversation is compacted: a request with the `Compaction` param returns a single `compaction` block, which then replaces the messages it summarizes. In a tool runner, call `CompactBeforeNextTurn()` and the runner does this for you:

```go
runner := client.Beta.Messages.NewToolRunner(tools, anthropic.BetaToolRunnerParams{
	BetaMessageNewParams: anthropic.BetaMessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		MaxTokens: 1024,
		Betas:     []anthropic.AnthropicBeta{anthropic.AnthropicBetaCompact2026_09_04},
		Messages: []anthropic.BetaMessageParam{
			anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("Find every page that mentions rate limits.")),
		},
	},
})

for message, err := range runner.All(ctx) {
	if err != nil {
		log.Fatal(err)
	}
	if message.Usage.InputTokens > 100_000 {
		runner.CompactBeforeNextTurn(anthropic.BetaCompactionConfigUnionParam{})
	}
}
```

The call only schedules the compaction. Once the current turn has finished, including any tool calls, the runner requests a summary, replaces its message history with the compaction response the API returns, and carries on. A turn that was paused (`pause_turn`) is resumed and finished first. If the current turn is the last one, the runner compacts and then stops, so the last message of the run is the compaction response. If you call it before the first turn, the compaction is the first request.

The compaction response is returned like any other message and doesn't count towards `MaxIterations`. Its `StopReason` is `compaction`, the summary is in its first content block, and its `Usage.InputTokens` is the size of the history that was just summarized. Calling `CompactBeforeNextTurn()` while handling that message does nothing, so a threshold like the one above doesn't compact twice. With the streaming runner, read the message from `runner.LastMessage()` once the turn's events have been consumed.

`CompactBeforeNextTurn()` takes the same config as the `Compaction` field of `BetaMessageNewParams`, and the zero value means `{"type": "summarize"}`. For example, to give your own summarization instructions:

```go
runner.CompactBeforeNextTurn(anthropic.BetaCompactionConfigUnionParam{
	OfSummarize: &anthropic.BetaSummarizeCompactionParam{
		Instructions: anthropic.String("Keep the page URLs found so far."),
	},
})
```

A few things to know:

- Calling it again before the compaction runs replaces the pending one.
- The runner doesn't add the beta for you, so pass it in `Betas`.
- `ContextManagement` is left out of the compaction request, because the API doesn't accept the two together, and is sent again afterwards. While `ContextManagement` has a compaction edit (`compact_20260112`) the compaction is not sent: the next `NextMessage()` or `NextStreaming()` call returns an error and the pending compaction is dropped.
- With the streaming runner, don't replace or append to `Params.Messages` while the compaction response is streaming, because that response is about to replace them. If you do, the stream ends with an error and your messages are kept. Other params can still be changed.
- If the API returns no summary, the runner prints a warning to stderr and keeps the history as it is.
- If the run ends on a turn that was cut short with tool calls that never ran (`max_tokens`, for example), the pending compaction is skipped with a warning printed to stderr. It is also skipped if the run stops at `MaxIterations` or you stop iterating.
- The `Compaction` param itself can't be set in the runner's `Params`, because every request in the loop would compact again: the next `NextMessage()` or `NextStreaming()` call clears it and returns an error.
- If the compaction request itself fails, the error is returned and the compaction is not retried. Call `CompactBeforeNextTurn()` again if you carry on.

### Inspecting State

```go
// Get most recent assistant message
lastMsg := runner.LastMessage()

// Get full conversation history (returns a copy)
messages := runner.Messages()

// Check iteration count
count := runner.IterationCount()

// Check if completed
if runner.IsCompleted() {
	// ...
}
```

## Error Handling

Tool execution errors are automatically converted to error results and sent back to Claude, allowing it to recover or try a different approach:

```go
func handler(ctx context.Context, input MyInput) (anthropic.BetaToolResultBlockParamContentUnion, error) {
	if input.City == "" {
		return anthropic.BetaToolResultBlockParamContentUnion{}, errors.New("city is required")
	}
	// ...
}
```

The error message will be sent to Claude as a tool result with `is_error: true`.

## Parallel Tool Execution

When Claude requests multiple tool calls in a single message, they are executed in parallel using an `errgroup`. This provides:

- Concurrent execution for better performance
- Proper context cancellation handling
- Results returned in the correct order

## Managed-agents sessions

The same `anthropic.BetaTool` shape works for managed-agents sessions. Two helpers cover the self-hosted side:

- `client.Beta.Sessions.Events.NewToolRunner(ctx, sessionID, anthropic.SessionToolRunnerOptions{...})` — the sessions-side counterpart to `client.Beta.Messages.NewToolRunner`. The session id is a positional argument (matching `list`/`send`/`stream` on the events resource); the options struct carries the tool registry and tuning knobs. It attaches to a session's event stream, dispatches the registered tools on both `agent.tool_use` (builtin tools, answered with `user.tool_result`) and `agent.custom_tool_use` (user-defined function tools, answered with `user.custom_tool_result`), and stops after the session is idle past `MaxIdle`. It does *only* that — no work claiming, lease heartbeating, or skill download.
- `environments.NewEnvironmentWorker(client, environments.EnvironmentWorkerOptions{...})` (in `github.com/anthropics/anthropic-sdk-go/lib/environments`) — the full self-hosted runner: it composes `environments.WorkPoller` (claim work) with a per-session `SessionToolRunner`, sets up the workdir + downloads the session agent's skills, heartbeats the work-item lease in parallel, force-stops the work on exit, and loops. A single `EnvironmentKey` authorizes everything — both the work-poll calls and the per-session calls. `worker.Run(ctx)` drives the poll loop (requires `EnvironmentID` + `EnvironmentKey`); `worker.HandleItem(ctx, environments.HandleItemOptions{...})` runs that same per-item flow (skills + run + heartbeat + force-stop) once for a work item you have already claimed yourself. Each `HandleItemOptions` field — `WorkID` / `EnvironmentID` / `SessionID` / `EnvironmentKey` — falls back to `ANTHROPIC_WORK_ID` / `ANTHROPIC_ENVIRONMENT_ID` / `ANTHROPIC_SESSION_ID` / `ANTHROPIC_ENVIRONMENT_KEY` when left empty (and `EnvironmentKey` also falls back to the worker's own `EnvironmentKey` option), so inside an `ant worker poll --on-work` hook (which exports all of them) it is just `worker.HandleItem(ctx, environments.HandleItemOptions{})`. If you are iterating `environments.WorkPoller` yourself, pass the claimed item through: `worker.HandleItem(ctx, environments.HandleItemOptions{WorkID: work.ID, EnvironmentID: work.EnvironmentID, SessionID: work.Data.ID, EnvironmentKey: environmentKey})`. When another process already acknowledged the item and sent its first heartbeat, pass that heartbeat response's `last_heartbeat` as `ExpectedLastHeartbeat` (or export `ANTHROPIC_WORK_LAST_HEARTBEAT`) and `HandleItem` continues that lease instead of starting one.

The standard `agent_toolset_20260401` tools (`bash`, `read`, `write`, `edit`, `glob`, `grep`), the workdir/skills `AgentToolContext`, and the skill-download helper live in `github.com/anthropics/anthropic-sdk-go/tools/agenttoolset`; `agenttoolset.BetaAgentToolset20260401(env)` returns them as a plain `[]anthropic.BetaTool` you can filter or extend. The file tools confine to the workdir (symlink-aware) and are safe without a sandbox; `bash` is unrestricted and should run inside one.

## Examples

See the [examples](./examples) directory for complete working examples:

- [examples/tool-runner](./examples/tool-runner) - Basic tool runner usage
- [examples/tool-runner-streaming](./examples/tool-runner-streaming) - Streaming with tool runner
- [examples/managed-agents-self-hosted-sandbox-worker](./examples/managed-agents-self-hosted-sandbox-worker) - Self-hosted environment worker
