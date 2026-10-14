package anthropic_test

import (
	"cmp"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
)

const contentResponseJSON = `{"id":"msg_1","type":"message","role":"assistant","model":"claude-served",
	"content":[{"type":"thinking","thinking":"secret","signature":"s"},{"type":"text","text":"Sunny."},
		{"type":"tool_use","id":"toolu_9","name":"get_weather","input":{"city":"Paris"}}],
	"stop_reason":"tool_use","usage":{"input_tokens":1,"output_tokens":1}}`

func contentParams() anthropic.MessageNewParams {
	return anthropic.MessageNewParams{
		Model:     "claude-requested",
		MaxTokens: 100,
		System:    []anthropic.TextBlockParam{{Text: "Be brief."}},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Weather?"), anthropic.NewImageBlockBase64("image/png", "AAAA")),
			anthropic.NewAssistantMessage(anthropic.NewToolUseBlock("toolu_1", map[string]any{"city": "Paris"}, "get_weather")),
			anthropic.NewUserMessage(anthropic.NewToolResultBlock("toolu_1", "18C", false)),
		},
	}
}

const contentInputJSON = `[
	{"role":"user","parts":[{"type":"text","content":"Weather?"},{"type":"anthropic_image"}]},
	{"role":"assistant","parts":[{"type":"tool_call","id":"toolu_1","name":"get_weather","arguments":{"city":"Paris"}}]},
	{"role":"user","parts":[{"type":"tool_call_response","id":"toolu_1","response":[{"type":"text","text":"18C"}]}]}
]`

const contentOutputJSON = `[{"role":"assistant","parts":[{"type":"reasoning","content":"secret"},{"type":"text","content":"Sunny."},
	{"type":"tool_call","id":"toolu_9","name":"get_weather","arguments":{"city":"Paris"}}]}]`

func TestOpenTelemetryContent(t *testing.T) {
	for name, opts := range map[string][]option.RequestOption{
		"option": {option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)},
		"env":    nil,
	} {
		t.Run(name, func(t *testing.T) {
			if opts == nil {
				t.Setenv("ANTHROPIC_OPEN_TELEMETRY_TRACES_CONTENT_MODE", "CONTENT")
			}
			client, recorder := newTracedClient(t, respond(200, "application/json", contentResponseJSON), opts...)
			if _, err := client.Messages.New(context.Background(), contentParams()); err != nil {
				t.Fatal(err)
			}
			span := onlySpan(t, recorder)
			wantJSON(t, span, "gen_ai.system_instructions", `[{"type":"text","content":"Be brief."}]`)
			wantJSON(t, span, "gen_ai.input.messages", contentInputJSON)
			wantJSON(t, span, "gen_ai.output.messages", contentOutputJSON)
			wantAttrs(t, span, map[string]any{"anthropic.content.truncated": nil})
		})
	}
}

// thinkingResponseJSON has thinking, a hidden thinking block, whose thinking is
// empty, and redacted thinking.
const thinkingResponseJSON = `{"id":"msg_1","type":"message","role":"assistant","model":"claude-served",
	"content":[{"type":"thinking","thinking":"Paris is warm.","signature":"sig_1"},{"type":"thinking","thinking":"","signature":"sig_2"},
		{"type":"redacted_thinking","data":"encrypted"},{"type":"text","text":"Sunny."}],
	"stop_reason":"end_turn","usage":{"input_tokens":1,"output_tokens":1}}`

const thinkingOutputJSON = `[{"role":"assistant","parts":[{"type":"reasoning","content":"Paris is warm."},{"type":"reasoning","content":""},
	{"type":"anthropic_redacted_thinking"},{"type":"text","content":"Sunny."}]}]`

// thinkingStreamBody is the stream of the Message of thinkingResponseJSON.
const thinkingStreamBody = "event: message_start\n" +
	`data: {"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","model":"claude-served","content":[],"usage":{"input_tokens":1,"output_tokens":1}}}` + "\n\n" +
	"event: content_block_start\n" +
	`data: {"type":"content_block_start","index":0,"content_block":{"type":"thinking","thinking":"","signature":""}}` + "\n\n" +
	"event: content_block_delta\n" +
	`data: {"type":"content_block_delta","index":0,"delta":{"type":"thinking_delta","thinking":"Paris is "}}` + "\n\n" +
	"event: content_block_delta\n" +
	`data: {"type":"content_block_delta","index":0,"delta":{"type":"thinking_delta","thinking":"warm."}}` + "\n\n" +
	"event: content_block_delta\n" +
	`data: {"type":"content_block_delta","index":0,"delta":{"type":"signature_delta","signature":"sig_1"}}` + "\n\n" +
	"event: content_block_stop\n" + `data: {"type":"content_block_stop","index":0}` + "\n\n" +
	"event: content_block_start\n" +
	`data: {"type":"content_block_start","index":1,"content_block":{"type":"thinking","thinking":"","signature":""}}` + "\n\n" +
	"event: content_block_delta\n" +
	`data: {"type":"content_block_delta","index":1,"delta":{"type":"signature_delta","signature":"sig_2"}}` + "\n\n" +
	"event: content_block_stop\n" + `data: {"type":"content_block_stop","index":1}` + "\n\n" +
	"event: content_block_start\n" +
	`data: {"type":"content_block_start","index":2,"content_block":{"type":"redacted_thinking","data":"encrypted"}}` + "\n\n" +
	"event: content_block_stop\n" + `data: {"type":"content_block_stop","index":2}` + "\n\n" +
	"event: content_block_start\n" +
	`data: {"type":"content_block_start","index":3,"content_block":{"type":"text","text":""}}` + "\n\n" +
	"event: content_block_delta\n" +
	`data: {"type":"content_block_delta","index":3,"delta":{"type":"text_delta","text":"Sunny."}}` + "\n\n" +
	"event: content_block_stop\n" + `data: {"type":"content_block_stop","index":3}` + "\n\n" +
	"event: message_delta\n" +
	`data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":1}}` + "\n\n" +
	"event: message_stop\n" +
	`data: {"type":"message_stop"}` + "\n\n"

func TestOpenTelemetryThinkingContent(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	// Thinking is a reasoning part with the text and no signature, a hidden block's
	// content is empty, and redacted thinking holds its type only, with the prefix.
	for _, api := range []string{"messages", "beta"} {
		t.Run(api, func(t *testing.T) {
			streamed, unstreamed := streamedAndNot(t, api, thinkingStreamBody, thinkingResponseJSON, content)
			for _, span := range []sdktrace.ReadOnlySpan{streamed, unstreamed} {
				wantJSON(t, span, "gen_ai.output.messages", thinkingOutputJSON)
				for _, kv := range span.Attributes() {
					if value := kv.Value.Emit(); strings.Contains(value, "sig") || strings.Contains(value, "encrypted") {
						t.Errorf("attribute %s = %s, want no signature and no redacted data", kv.Key, value)
					}
				}
			}
		})
	}

	t.Run("an earlier turn in the request", func(t *testing.T) {
		client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON), content)
		params := contentParams()
		params.Messages = []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Weather?")),
			anthropic.NewAssistantMessage(anthropic.NewThinkingBlock("sig_1", "Paris is warm."), anthropic.NewRedactedThinkingBlock("encrypted"),
				anthropic.NewTextBlock("Sunny.")),
			anthropic.NewUserMessage(anthropic.NewTextBlock("Tomorrow?")),
		}
		if _, err := client.Messages.New(context.Background(), params); err != nil {
			t.Fatal(err)
		}
		span := onlySpan(t, recorder)
		wantJSON(t, span, "gen_ai.input.messages", `[
			{"role":"user","parts":[{"type":"text","content":"Weather?"}]},
			{"role":"assistant","parts":[{"type":"reasoning","content":"Paris is warm."},{"type":"anthropic_redacted_thinking"},{"type":"text","content":"Sunny."}]},
			{"role":"user","parts":[{"type":"text","content":"Tomorrow?"}]}]`)
		if got, _ := attr(span, "gen_ai.input.messages"); strings.Contains(got.AsString(), "sig_1") || strings.Contains(got.AsString(), "encrypted") {
			t.Errorf("input messages = %s, want no signature and no redacted data", got.AsString())
		}
	})

	t.Run("metadata only", func(t *testing.T) {
		streamed, unstreamed := streamedAndNot(t, "messages", thinkingStreamBody, thinkingResponseJSON)
		for _, span := range []sdktrace.ReadOnlySpan{streamed, unstreamed} {
			for _, kv := range span.Attributes() {
				if strings.Contains(kv.Value.Emit(), "Paris is warm.") {
					t.Errorf("attribute %s = %s, want no thinking", kv.Key, kv.Value.Emit())
				}
			}
		}
	})
}

func TestOpenTelemetryToolResultContentInARequest(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON), content)
	params := contentParams()
	params.Messages = []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewDocumentBlock(anthropic.PlainTextSourceParam{Data: "Forecasts"})),
		anthropic.NewUserMessage(
			anthropic.ContentBlockParamUnion{OfToolResult: &anthropic.ToolResultBlockParam{ToolUseID: "toolu_1", Content: []anthropic.ToolResultBlockParamContentUnion{
				{OfText: &anthropic.TextBlockParam{
					Text:         "sunny",
					Citations:    []anthropic.TextCitationParamUnion{{OfWebSearchResultLocation: &anthropic.CitationWebSearchResultLocationParam{CitedText: "sun", URL: "https://example.com", EncryptedIndex: "opaque"}}},
					CacheControl: anthropic.NewCacheControlEphemeralParam(),
				}},
				{OfSearchResult: &anthropic.SearchResultBlockParam{Source: "https://example.com", Title: "Forecast", Content: []anthropic.TextBlockParam{{Text: "rain"}}}},
				{OfImage: &anthropic.ImageBlockParam{Source: anthropic.ImageBlockParamSourceUnion{
					OfBase64: &anthropic.Base64ImageSourceParam{Data: "AAAA", MediaType: anthropic.Base64ImageSourceMediaTypeImagePNG},
				}}},
			}}},
			anthropic.ContentBlockParamUnion{OfToolResult: &anthropic.ToolResultBlockParam{ToolUseID: "toolu_2"}},
		),
	}
	if _, err := client.Messages.New(context.Background(), params); err != nil {
		t.Fatal(err)
	}
	// A block of a message that is no text, tool call, tool result or thinking has
	// the prefix. The response of a tool result is the blocks as they are sent: text
	// and search results with all their fields but the encrypted index of the
	// citation, any other block its type only, without the prefix, and none for a
	// result without content.
	span := onlySpan(t, recorder)
	if got, _ := attr(span, "gen_ai.input.messages"); strings.Contains(got.AsString(), "opaque") {
		t.Errorf("input messages = %s, want no encrypted index", got.AsString())
	}
	wantJSON(t, span, "gen_ai.input.messages", `[
		{"role":"user","parts":[{"type":"anthropic_document"}]},
		{"role":"user","parts":[
			{"type":"tool_call_response","id":"toolu_1","response":[
				{"type":"text","text":"sunny","citations":[{"type":"web_search_result_location","cited_text":"sun","url":"https://example.com"}],
					"cache_control":{"type":"ephemeral"}},
				{"type":"search_result","source":"https://example.com","title":"Forecast","content":[{"type":"text","text":"rain"}]},
				{"type":"image"}]},
			{"type":"tool_call_response","id":"toolu_2","response":[]}]}]`)
}

// toolCallBlocks are the content blocks of a response in which the model
// called a server tool, a tool of an MCP server, and client tools, one of them a
// member of a toolset. Each has the input that its tool-use block streams, or
// none. Only a tool_use block's toolset_name is part of the API: the other two
// are there to show that they are left out.
var toolCallBlocks = []struct{ block, input string }{
	{`{"type":"server_tool_use","id":"srvtoolu_1","name":"web_search","toolset_name":"computer","input":{}}`, `{"query":"weather"}`},
	{`{"type":"web_search_tool_result","tool_use_id":"srvtoolu_1","content":[
		{"type":"web_search_result","url":"https://example.com","title":"Forecast","encrypted_content":"opaque","page_age":null}]}`, ""},
	{`{"type":"mcp_tool_use","id":"mcptoolu_1","name":"get_weather","server_name":"weather","toolset_name":"computer","input":{}}`, `{"city":"Paris"}`},
	{`{"type":"mcp_tool_result","tool_use_id":"mcptoolu_1","is_error":false,"content":"sunny"}`, ""},
	{`{"type":"code_execution_tool_result","tool_use_id":"srvtoolu_2","content":{"type":"code_execution_result","stdout":"1","stderr":"","return_code":0,"content":[]}}`, ""},
	{`{"type":"tool_use","id":"toolu_1","name":"screenshot","toolset_name":"computer","input":{}}`, ""},
	{`{"type":"tool_use","id":"toolu_2","name":"get_time","toolset_name":null,"input":{}}`, `{"zone":"UTC"}`},
	{`{"type":"tool_use","id":"toolu_3","name":"get_time","input":{}}`, `{"zone":"CET"}`},
}

const toolCallOutputJSON = `[{"role":"assistant","parts":[
	{"type":"tool_call","id":"srvtoolu_1","name":"web_search","arguments":{"query":"weather"}},
	{"type":"tool_call_response","id":"srvtoolu_1","response":[{"type":"web_search_result","url":"https://example.com","title":"Forecast","page_age":null}]},
	{"type":"tool_call","id":"mcptoolu_1","name":"get_weather","arguments":{"city":"Paris"},"anthropic_mcp_server_name":"weather"},
	{"type":"tool_call_response","id":"mcptoolu_1","response":[{"type":"text","text":"sunny"}]},
	{"type":"tool_call_response","id":"srvtoolu_2","response":{"type":"code_execution_result","stdout":"1","stderr":"","return_code":0,"content":[]}},
	{"type":"tool_call","id":"toolu_1","name":"screenshot","arguments":{},"anthropic_toolset_name":"computer"},
	{"type":"tool_call","id":"toolu_2","name":"get_time","arguments":{"zone":"UTC"}},
	{"type":"tool_call","id":"toolu_3","name":"get_time","arguments":{"zone":"CET"}}]}]`

func TestOpenTelemetryServerAndMCPToolContent(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	var blocks []string
	var stream strings.Builder
	stream.WriteString("event: message_start\n" +
		`data: {"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","model":"claude-served","content":[],"usage":{"input_tokens":1,"output_tokens":1}}}` + "\n\n")
	for i, b := range toolCallBlocks {
		index := strconv.Itoa(i)
		compact := strings.NewReplacer("\n", "", "\t", "").Replace(b.block)
		stream.WriteString("event: content_block_start\n" + `data: {"type":"content_block_start","index":` + index + `,"content_block":` + compact + `}` + "\n\n")
		if b.input != "" {
			stream.WriteString("event: content_block_delta\n" + `data: {"type":"content_block_delta","index":` + index +
				`,"delta":{"type":"input_json_delta","partial_json":` + strconv.Quote(b.input) + `}}` + "\n\n")
			compact = strings.Replace(compact, `"input":{}`, `"input":`+b.input, 1)
		}
		stream.WriteString("event: content_block_stop\n" + `data: {"type":"content_block_stop","index":` + index + `}` + "\n\n")
		blocks = append(blocks, compact)
	}
	stream.WriteString("event: message_delta\n" + `data: {"type":"message_delta","delta":{"stop_reason":"tool_use"},"usage":{"output_tokens":2}}` + "\n\n" +
		"event: message_stop\n" + `data: {"type":"message_stop"}` + "\n\n")
	message := `{"id":"msg_1","type":"message","role":"assistant","model":"claude-served","content":[` + strings.Join(blocks, ",") +
		`],"stop_reason":"tool_use","usage":{"input_tokens":1,"output_tokens":2}}`

	// A server tool's call and the call of an MCP server's tool are tool calls, and
	// their result blocks are the responses, without the encrypted content. Of the
	// tool calls, only a toolset member's tool_use has the toolset's name.
	for _, api := range []string{"messages", "beta"} {
		t.Run(api, func(t *testing.T) {
			streamed, unstreamed := streamedAndNot(t, api, stream.String(), message, content)
			for _, span := range []sdktrace.ReadOnlySpan{streamed, unstreamed} {
				wantJSON(t, span, "gen_ai.output.messages", toolCallOutputJSON)
				if got, _ := attr(span, "gen_ai.output.messages"); strings.Contains(got.AsString(), "opaque") {
					t.Errorf("output messages = %s, want no encrypted content", got.AsString())
				}
			}
		})
	}

	t.Run("an earlier turn in the request", func(t *testing.T) {
		client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON), content)
		params := betaContentParams()
		params.Messages = append(params.Messages, anthropic.BetaMessageParam{Role: anthropic.BetaMessageParamRoleAssistant, Content: []anthropic.BetaContentBlockParamUnion{
			{OfServerToolUse: &anthropic.BetaServerToolUseBlockParam{ID: "srvtoolu_1", Name: anthropic.BetaServerToolUseBlockParamNameWebSearch, Input: map[string]any{"query": "weather"}}},
			{OfWebSearchToolResult: &anthropic.BetaWebSearchToolResultBlockParam{ToolUseID: "srvtoolu_1", Content: anthropic.BetaWebSearchToolResultBlockParamContentUnion{
				OfResultBlock: []anthropic.BetaWebSearchResultBlockParam{{URL: "https://example.com", Title: "Forecast", EncryptedContent: "opaque"}},
			}}},
			{OfMCPToolUse: &anthropic.BetaMCPToolUseBlockParam{ID: "mcptoolu_1", Name: "get_weather", ServerName: "weather", Input: map[string]any{"city": "Paris"}}},
			{OfMCPToolResult: &anthropic.BetaRequestMCPToolResultBlockParam{ToolUseID: "mcptoolu_1", Content: anthropic.BetaRequestMCPToolResultBlockParamContentUnion{
				OfBetaMCPToolResultBlockContent: []anthropic.BetaTextBlockParam{{Text: "sunny"}},
			}}},
			{OfMCPToolResult: &anthropic.BetaRequestMCPToolResultBlockParam{ToolUseID: "mcptoolu_2"}},
			{OfToolUse: &anthropic.BetaToolUseBlockParam{ID: "toolu_1", Name: "screenshot", ToolsetName: anthropic.String("computer"), Input: map[string]any{}}},
			{OfToolUse: &anthropic.BetaToolUseBlockParam{ID: "toolu_2", Name: "get_time", ToolsetName: param.Null[string](), Input: map[string]any{"zone": "UTC"}}},
			{OfToolUse: &anthropic.BetaToolUseBlockParam{ID: "toolu_3", Name: "get_time", Input: map[string]any{"zone": "CET"}}},
		}})
		// Neither block has a toolset_name in the API.
		assistant := params.Messages[len(params.Messages)-1].Content
		assistant[0].OfServerToolUse.SetExtraFields(map[string]any{"toolset_name": "computer"})
		assistant[2].OfMCPToolUse.SetExtraFields(map[string]any{"toolset_name": "computer"})
		if _, err := client.Beta.Messages.New(context.Background(), params); err != nil {
			t.Fatal(err)
		}
		span := onlySpan(t, recorder)
		wantJSON(t, span, "gen_ai.input.messages", `[
			{"role":"user","parts":[{"type":"text","content":"Weather?"}]},
			{"role":"assistant","parts":[
				{"type":"tool_call","id":"srvtoolu_1","name":"web_search","arguments":{"query":"weather"}},
				{"type":"tool_call_response","id":"srvtoolu_1","response":[{"type":"web_search_result","url":"https://example.com","title":"Forecast"}]},
				{"type":"tool_call","id":"mcptoolu_1","name":"get_weather","arguments":{"city":"Paris"},"anthropic_mcp_server_name":"weather"},
				{"type":"tool_call_response","id":"mcptoolu_1","response":[{"type":"text","text":"sunny"}]},
				{"type":"tool_call_response","id":"mcptoolu_2","response":[]},
				{"type":"tool_call","id":"toolu_1","name":"screenshot","arguments":{},"anthropic_toolset_name":"computer"},
				{"type":"tool_call","id":"toolu_2","name":"get_time","arguments":{"zone":"UTC"}},
				{"type":"tool_call","id":"toolu_3","name":"get_time","arguments":{"zone":"CET"}}]}]`)
		if got, _ := attr(span, "gen_ai.input.messages"); strings.Contains(got.AsString(), "opaque") {
			t.Errorf("input messages = %s, want no encrypted content", got.AsString())
		}
	})
}

func TestOpenTelemetryMetadataOnly(t *testing.T) {
	for name, test := range map[string]struct {
		opts []option.RequestOption
		env  map[string]string
	}{
		"default": {},
		"option over the environment variable": {
			opts: []option.RequestOption{option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeMetadataOnly)},
			env:  map[string]string{"ANTHROPIC_OPEN_TELEMETRY_TRACES_CONTENT_MODE": "content"},
		},
		"a mode that is not one of the two": {
			opts: []option.RequestOption{option.WithOpenTelemetryTracesContentMode("everything")},
			env:  map[string]string{"ANTHROPIC_OPEN_TELEMETRY_TRACES_CONTENT_MODE": "content"},
		},
		"an environment variable that is not one of the two": {
			env: map[string]string{"ANTHROPIC_OPEN_TELEMETRY_TRACES_CONTENT_MODE": "everything"},
		},
		// The variable had this name before the content mode became a setting of the traces.
		"the former name of the environment variable": {
			env: map[string]string{"ANTHROPIC_OPEN_TELEMETRY_CONTENT_MODE": "content"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			for key, value := range test.env {
				t.Setenv(key, value)
			}
			client, recorder := newTracedClient(t, respond(200, "application/json", contentResponseJSON), test.opts...)
			if _, err := client.Messages.New(context.Background(), contentParams()); err != nil {
				t.Fatal(err)
			}
			wantAttrs(t, onlySpan(t, recorder), map[string]any{
				"gen_ai.system_instructions":  nil,
				"gen_ai.input.messages":       nil,
				"gen_ai.output.messages":      nil,
				"anthropic.content.truncated": nil,
				"gen_ai.response.id":          "msg_1",
			})
		})
	}
}

const contentStreamBody = "event: message_start\n" +
	`data: {"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","model":"claude-served","content":[],"usage":{"input_tokens":1,"output_tokens":1}}}` + "\n\n" +
	"event: content_block_start\n" +
	`data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}` + "\n\n" +
	"event: content_block_delta\n" +
	`data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Sun"}}` + "\n\n" +
	"event: content_block_delta\n" +
	`data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"ny."}}` + "\n\n" +
	"event: content_block_start\n" +
	`data: {"type":"content_block_start","index":1,"content_block":{"type":"tool_use","id":"toolu_9","name":"get_weather","input":{}}}` + "\n\n" +
	"event: content_block_delta\n" +
	`data: {"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"{\"city\": "}}` + "\n\n" +
	"event: content_block_delta\n" +
	`data: {"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"\"Paris\"}"}}` + "\n\n" +
	"event: message_delta\n" +
	`data: {"type":"message_delta","delta":{"stop_reason":"tool_use"},"usage":{"output_tokens":5}}` + "\n\n" +
	"event: message_stop\n" +
	`data: {"type":"message_stop"}` + "\n\n"

// contentStreamMessageJSON is the Message of contentStreamBody, not streamed.
const contentStreamMessageJSON = `{"id":"msg_1","type":"message","role":"assistant","model":"claude-served",
	"content":[{"type":"text","text":"Sunny."},{"type":"tool_use","id":"toolu_9","name":"get_weather","input":{"city":"Paris"}}],
	"stop_reason":"tool_use","usage":{"input_tokens":1,"output_tokens":5}}`

func betaContentParams() anthropic.BetaMessageNewParams {
	return anthropic.BetaMessageNewParams{
		Model:     "claude-requested",
		MaxTokens: 100,
		Messages:  []anthropic.BetaMessageParam{anthropic.NewBetaUserMessage(anthropic.NewBetaTextBlock("Weather?"))},
	}
}

// streamedAndNot makes the Messages call of api ("messages" or "beta") twice,
// streamed from stream and not streamed from message, reads both to the end, and
// returns the span of each.
func streamedAndNot(t *testing.T, api, stream, message string, opts ...option.RequestOption) (streamed, unstreamed sdktrace.ReadOnlySpan) {
	t.Helper()
	ctx := context.Background()
	client, recorder := newTracedClient(t, respond(200, "text/event-stream", stream), opts...)
	var err error
	if api == "beta" {
		s := client.Beta.Messages.NewStreaming(ctx, betaContentParams())
		for s.Next() {
		}
		err = s.Err()
	} else {
		s := client.Messages.NewStreaming(ctx, contentParams())
		for s.Next() {
		}
		err = s.Err()
	}
	if err != nil {
		t.Fatal(err)
	}
	streamed = onlySpan(t, recorder)

	client, recorder = newTracedClient(t, respond(200, "application/json", message), opts...)
	if api == "beta" {
		_, err = client.Beta.Messages.New(ctx, betaContentParams())
	} else {
		_, err = client.Messages.New(ctx, contentParams())
	}
	if err != nil {
		t.Fatal(err)
	}
	return streamed, onlySpan(t, recorder)
}

func TestOpenTelemetryStreamContent(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	for _, api := range []string{"messages", "beta"} {
		t.Run(api, func(t *testing.T) {
			streamed, unstreamed := streamedAndNot(t, api, contentStreamBody, contentStreamMessageJSON, content)

			wantJSON(t, streamed, "gen_ai.output.messages", `[{"role":"assistant","parts":[{"type":"text","content":"Sunny."},
				{"type":"tool_call","id":"toolu_9","name":"get_weather","arguments":{"city":"Paris"}}]}]`)
			got, _ := attr(streamed, "gen_ai.output.messages")
			want, _ := attr(unstreamed, "gen_ai.output.messages")
			if got.AsString() != want.AsString() {
				t.Errorf("streamed output messages = %s, want those of the same Message not streamed, %s", got.AsString(), want.AsString())
			}
			if api == "messages" {
				wantJSON(t, streamed, "gen_ai.input.messages", contentInputJSON)
			}
			wantAttrs(t, streamed, map[string]any{
				"anthropic.message.stop_reason": "tool_use",
				"gen_ai.usage.output_tokens":    int64(5),
			})
		})
	}
}

func TestOpenTelemetryStreamContentIsTheSameWhenTheApplicationAccumulatesToo(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	ctx := context.Background()
	want := `[{"role":"assistant","parts":[{"type":"text","content":"Sunny."},
		{"type":"tool_call","id":"toolu_9","name":"get_weather","arguments":{"city":"Paris"}}]}]`

	// The span builds the streamed Message on a value of its own, so a Message that
	// the application builds from the same events changes nothing, and is whole too.
	t.Run("messages", func(t *testing.T) {
		client, recorder := newTracedClient(t, respond(200, "text/event-stream", contentStreamBody), content)
		stream := client.Messages.NewStreaming(ctx, contentParams())
		var message anthropic.Message
		for stream.Next() {
			if err := message.Accumulate(stream.Current()); err != nil {
				t.Fatal(err)
			}
		}
		if err := stream.Err(); err != nil {
			t.Fatal(err)
		}
		if len(message.Content) != 2 || message.Content[0].Text != "Sunny." {
			t.Errorf("the application accumulated %d content blocks, want the text and the tool call", len(message.Content))
		}
		wantJSON(t, onlySpan(t, recorder), "gen_ai.output.messages", want)
	})
	t.Run("beta", func(t *testing.T) {
		client, recorder := newTracedClient(t, respond(200, "text/event-stream", contentStreamBody), content)
		stream := client.Beta.Messages.NewStreaming(ctx, betaContentParams())
		var message anthropic.BetaMessage
		for stream.Next() {
			if err := message.Accumulate(stream.Current()); err != nil {
				t.Fatal(err)
			}
		}
		if err := stream.Err(); err != nil {
			t.Fatal(err)
		}
		if len(message.Content) != 2 || message.Content[0].Text != "Sunny." {
			t.Errorf("the application accumulated %d content blocks, want the text and the tool call", len(message.Content))
		}
		wantJSON(t, onlySpan(t, recorder), "gen_ai.output.messages", want)
	})
}

func TestOpenTelemetryStreamContentDoesNotDecodeTheEventsAgain(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	start, stop, _ := strings.Cut(streamBody, "event: content_block_stop\n")
	stop = "event: content_block_stop\n" + stop
	delta := "event: content_block_delta\n" + `data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Sunny. "}}` + "\n\n"

	// allocations returns how many allocations an event of the stream costs when
	// read reads it: those of reading a long stream less those of reading a short
	// one, which cost the same but for their events.
	allocations := func(read func(client anthropic.Client), opts ...option.RequestOption) float64 {
		const short, long = 100, 300
		perStream := func(deltas int) float64 {
			client, recorder := newTracedClient(t, respond(200, "text/event-stream", start+strings.Repeat(delta, deltas)+stop), opts...)
			allocs := testing.AllocsPerRun(5, func() { read(client) })
			for _, span := range recorder.Spans() {
				if _, ok := attr(span, "gen_ai.output.messages"); !ok {
					t.Fatal("a span has no output messages")
				}
			}
			return allocs
		}
		return (perStream(long) - perStream(short)) / (long - short)
	}

	for api, read := range map[string]func(accumulate bool) func(anthropic.Client){
		"messages": func(accumulate bool) func(anthropic.Client) {
			return func(client anthropic.Client) {
				stream := client.Messages.NewStreaming(context.Background(), contentParams())
				var message anthropic.Message
				for stream.Next() {
					if accumulate && message.Accumulate(stream.Current()) != nil {
						t.Fatal("the event was not accumulated")
					}
				}
				if err := stream.Err(); err != nil {
					t.Fatal(err)
				}
			}
		},
		"beta": func(accumulate bool) func(anthropic.Client) {
			return func(client anthropic.Client) {
				stream := client.Beta.Messages.NewStreaming(context.Background(), betaContentParams())
				var message anthropic.BetaMessage
				for stream.Next() {
					if accumulate && message.Accumulate(stream.Current()) != nil {
						t.Fatal("the event was not accumulated")
					}
				}
				if err := stream.Err(); err != nil {
					t.Fatal(err)
				}
			}
		},
	} {
		t.Run(api, func(t *testing.T) {
			// The span builds its Message from the values that the stream decoded, so an
			// event costs what it costs an application that accumulates without a span.
			// Decoding the event again would cost over twenty allocations more.
			traced := allocations(read(false), content)
			untraced := allocations(read(true), option.WithoutOpenTelemetry())
			if traced > untraced+2 {
				t.Errorf("an event costs %.1f allocations with content on the span, and %.1f when the application accumulates without a span", traced, untraced)
			}
		})
	}
}

func TestOpenTelemetryStreamMetadataOnly(t *testing.T) {
	for _, api := range []string{"messages", "beta"} {
		streamed, unstreamed := streamedAndNot(t, api, contentStreamBody, contentStreamMessageJSON)
		wantAttrs(t, streamed, map[string]any{
			"gen_ai.output.messages":        nil,
			"gen_ai.input.messages":         nil,
			"gen_ai.response.id":            "msg_1",
			"anthropic.message.stop_reason": "tool_use",
			"gen_ai.usage.output_tokens":    int64(5),
		})
		// The streamed Message is the one that the SDK's accumulator builds, so
		// its span has the response attributes of the same Message not streamed.
		for _, kv := range unstreamed.Attributes() {
			if key := string(kv.Key); strings.HasPrefix(key, "gen_ai.response.") || strings.HasPrefix(key, "gen_ai.usage.") {
				wantAttrs(t, streamed, map[string]any{key: kv.Value.AsInterface()})
			}
		}
	}
}

// stopSequencesSpan makes the Messages call of api ("messages" or "beta"),
// streamed or not, with the stop sequences stop, reads it to the end, and returns
// its span.
func stopSequencesSpan(t *testing.T, api string, streamed bool, stop []string, opts ...option.RequestOption) sdktrace.ReadOnlySpan {
	t.Helper()
	ctx := context.Background()
	contentType, body := "application/json", contentStreamMessageJSON
	if streamed {
		contentType, body = "text/event-stream", contentStreamBody
	}
	client, recorder := newTracedClient(t, respond(200, contentType, body), opts...)
	params, betaParams := contentParams(), betaContentParams()
	params.StopSequences, betaParams.StopSequences = stop, stop
	var err error
	switch {
	case api == "beta" && streamed:
		s := client.Beta.Messages.NewStreaming(ctx, betaParams)
		for s.Next() {
		}
		err = s.Err()
	case api == "beta":
		_, err = client.Beta.Messages.New(ctx, betaParams)
	case streamed:
		s := client.Messages.NewStreaming(ctx, params)
		for s.Next() {
		}
		err = s.Err()
	default:
		_, err = client.Messages.New(ctx, params)
	}
	if err != nil {
		t.Fatal(err)
	}
	return onlySpan(t, recorder)
}

func TestOpenTelemetryStopSequences(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	stop := []string{"END", "\n\nHuman:"}
	for _, api := range []string{"messages", "beta"} {
		for _, streamed := range []bool{false, true} {
			t.Run(api+map[bool]string{true: " streamed", false: ""}[streamed], func(t *testing.T) {
				span := stopSequencesSpan(t, api, streamed, stop, content)
				// A list of strings, as anthropic.beta is, and not JSON.
				if got, _ := attr(span, "gen_ai.request.stop_sequences"); got.Type() != attribute.STRINGSLICE {
					t.Errorf("gen_ai.request.stop_sequences is a %v, want a list of strings", got.Type())
				}
				wantAttrs(t, span, map[string]any{
					"gen_ai.request.stop_sequences": stop,
					"anthropic.content.truncated":   nil,
				})

				// The stop sequences are content: a span without content has none.
				wantAttrs(t, stopSequencesSpan(t, api, streamed, stop), map[string]any{"gen_ai.request.stop_sequences": nil})
				// A request that sets none, or an empty list, has no such attribute.
				wantAttrs(t, stopSequencesSpan(t, api, streamed, nil, content), map[string]any{"gen_ai.request.stop_sequences": nil})
				wantAttrs(t, stopSequencesSpan(t, api, streamed, []string{}, content), map[string]any{"gen_ai.request.stop_sequences": nil})
				// A stop sequence that is not valid UTF-8 can fail the export: the span
				// leaves the attribute out, and keeps the rest of the content.
				invalid := stopSequencesSpan(t, api, streamed, []string{"END", "\xff\xfe"}, content)
				wantAttrs(t, invalid, map[string]any{"gen_ai.request.stop_sequences": nil, "anthropic.content.truncated": nil})
				for _, key := range []string{"gen_ai.input.messages", "gen_ai.output.messages"} {
					if _, ok := attr(invalid, key); !ok {
						t.Errorf("attribute %s is absent, want only the stop sequences left out", key)
					}
				}

				// Stop sequences that fill the size limit stay. One byte less, and they
				// are over it on their own: the span leaves the attribute out and says so.
				size := int64(len(stop[0]) + len(stop[1]))
				filled := stopSequencesSpan(t, api, streamed, stop, content, option.WithOpenTelemetryTracesMaxContentBytes(size))
				wantAttrs(t, filled, map[string]any{"gen_ai.request.stop_sequences": stop, "anthropic.content.truncated": true})
				over := stopSequencesSpan(t, api, streamed, stop, content, option.WithOpenTelemetryTracesMaxContentBytes(size-1))
				wantAttrs(t, over, map[string]any{"gen_ai.request.stop_sequences": nil, "anthropic.content.truncated": true})
				for _, key := range []string{"gen_ai.input.messages", "gen_ai.output.messages"} {
					if _, ok := attr(over, key); !ok {
						t.Errorf("attribute %s is absent, want only the stop sequences left out", key)
					}
				}
			})
		}
	}
}

func TestOpenTelemetryStopSequencesCountTowardTheSizeBoundAndGoWholeWhenTheyAloneAreOverIt(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	big := strings.Repeat("a", 140_000)
	params := messageParams()
	params.Messages = nil
	for i := 0; i < 5; i++ {
		params.Messages = append(params.Messages, anthropic.NewUserMessage(anthropic.NewTextBlock(string(rune('0'+i))+big)))
	}

	// Five messages of 140 KB fit on their own. The 450 KB of stop sequences leave
	// room for the parts of the newest two, and the others stay without parts.
	params.StopSequences = []string{strings.Repeat("x", 300_000), strings.Repeat("y", 150_000)}
	client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON), content)
	if _, err := client.Messages.New(context.Background(), params); err != nil {
		t.Fatal(err)
	}
	span := onlySpan(t, recorder)
	wantPartCounts(t, span, "gen_ai.input.messages", 0, 0, 0, 1, 1)
	wantTextParts(t, span, "gen_ai.input.messages", len(big)+1, "3a", "4a")
	wantJSON(t, span, "gen_ai.output.messages", `[{"role":"assistant","parts":[{"type":"text","content":"hi"}]}]`)
	wantAttrs(t, span, map[string]any{
		"gen_ai.request.stop_sequences": params.StopSequences,
		"anthropic.content.truncated":   true,
	})

	// Stop sequences that are over the limit on their own are left out whole, so
	// every part has room.
	params.StopSequences = []string{strings.Repeat("x", 1_100_000), "END"}
	client, recorder = newTracedClient(t, respond(200, "application/json", messageJSON), content)
	if _, err := client.Messages.New(context.Background(), params); err != nil {
		t.Fatal(err)
	}
	span = onlySpan(t, recorder)
	wantPartCounts(t, span, "gen_ai.input.messages", 1, 1, 1, 1, 1)
	wantJSON(t, span, "gen_ai.output.messages", `[{"role":"assistant","parts":[{"type":"text","content":"hi"}]}]`)
	wantAttrs(t, span, map[string]any{
		"gen_ai.request.stop_sequences": nil,
		"anthropic.content.truncated":   true,
	})
}

// fallbackStreamBody is the stream of a beta Message that a fallback model took
// over: message_start names the requested model, and the fallback block the one
// that served the response.
const fallbackStreamBody = "event: message_start\n" +
	`data: {"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","model":"claude-requested","content":[],"usage":{"input_tokens":1,"output_tokens":1}}}` + "\n\n" +
	"event: content_block_start\n" +
	`data: {"type":"content_block_start","index":0,"content_block":{"type":"fallback","from":{"model":"claude-requested"},"to":{"model":"claude-fallback"},"trigger":{"type":"refusal"}}}` + "\n\n" +
	"event: content_block_stop\n" + `data: {"type":"content_block_stop","index":0}` + "\n\n" +
	"event: content_block_start\n" +
	`data: {"type":"content_block_start","index":1,"content_block":{"type":"text","text":""}}` + "\n\n" +
	"event: content_block_delta\n" +
	`data: {"type":"content_block_delta","index":1,"delta":{"type":"text_delta","text":"Sunny."}}` + "\n\n" +
	"event: content_block_stop\n" + `data: {"type":"content_block_stop","index":1}` + "\n\n" +
	"event: message_delta\n" +
	`data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":5}}` + "\n\n" +
	"event: message_stop\n" + `data: {"type":"message_stop"}` + "\n\n"

// fallbackMessageJSON is the same Message, not streamed.
const fallbackMessageJSON = `{"id":"msg_1","type":"message","role":"assistant","model":"claude-fallback",
	"content":[{"type":"fallback","from":{"model":"claude-requested"},"to":{"model":"claude-fallback"},"trigger":{"type":"refusal"}},{"type":"text","text":"Sunny."}],
	"stop_reason":"end_turn","usage":{"input_tokens":1,"output_tokens":5}}`

func TestOpenTelemetryStreamFallbackModel(t *testing.T) {
	for name, opts := range map[string][]option.RequestOption{
		"metadata only": nil,
		"content":       {option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)},
	} {
		t.Run(name, func(t *testing.T) {
			streamed, unstreamed := streamedAndNot(t, "beta", fallbackStreamBody, fallbackMessageJSON, opts...)
			// The model that served the response, as on the same Message not streamed.
			wantAttrs(t, unstreamed, map[string]any{"gen_ai.response.model": "claude-fallback"})
			wantAttrs(t, streamed, map[string]any{"gen_ai.response.model": "claude-fallback", "gen_ai.output.messages": func() any {
				if got, ok := attr(unstreamed, "gen_ai.output.messages"); ok {
					return got.AsString()
				}
				return nil
			}()})
		})
	}
}

func TestOpenTelemetryContentToolUseWithoutInput(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON), content)

	// A nil input is left out of the request.
	params := messageParams()
	params.Messages = []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock("weather?")),
		anthropic.NewAssistantMessage(anthropic.NewToolUseBlock("toolu_1", nil, "get_weather")),
		anthropic.NewUserMessage(anthropic.NewToolResultBlock("toolu_1", "sunny", false)),
	}
	if _, err := client.Messages.New(context.Background(), params); err != nil {
		t.Fatal(err)
	}
	wantJSON(t, onlySpan(t, recorder), "gen_ai.input.messages", `[
		{"role":"user","parts":[{"type":"text","content":"weather?"}]},
		{"role":"assistant","parts":[{"type":"tool_call","id":"toolu_1","name":"get_weather","arguments":{}}]},
		{"role":"user","parts":[{"type":"tool_call_response","id":"toolu_1","response":[{"type":"text","text":"sunny"}]}]}]`)
}

// textStreamBody is the stream of a Message whose text comes in the deltas.
func textStreamBody(deltas ...string) string {
	body := "event: message_start\n" +
		`data: {"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","model":"claude-served","content":[],"usage":{"input_tokens":1,"output_tokens":1}}}` + "\n\n" +
		"event: content_block_start\n" + `data: {"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}` + "\n\n"
	var events strings.Builder
	for _, delta := range deltas {
		events.WriteString("event: content_block_delta\n" +
			`data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"` + delta + `"}}` + "\n\n")
	}
	return body + events.String() + "event: content_block_stop\n" + `data: {"type":"content_block_stop","index":0}` + "\n\n" +
		"event: message_delta\n" + `data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":2}}` + "\n\n" +
		"event: message_stop\n" + `data: {"type":"message_stop"}` + "\n\n"
}

func TestOpenTelemetryStreamContentSizeBound(t *testing.T) {
	piece := strings.Repeat("x", 120_000)
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	client, recorder := newTracedClient(t, respond(200, "text/event-stream", textStreamBody(piece, piece, piece)), content)

	// Five messages of 140 KB each fit on their own. The 360 KB of output leave
	// room for the parts of the newest two.
	big := strings.Repeat("a", 140_000)
	params := messageParams()
	params.Messages = nil
	for i := 0; i < 5; i++ {
		params.Messages = append(params.Messages, anthropic.NewUserMessage(anthropic.NewTextBlock(string(rune('0'+i))+big)))
	}
	stream := client.Messages.NewStreaming(context.Background(), params)
	for stream.Next() {
	}
	if err := stream.Err(); err != nil {
		t.Fatal(err)
	}

	span := onlySpan(t, recorder)
	var messages []struct {
		Parts []map[string]any `json:"parts"`
	}
	value, _ := attr(span, "gen_ai.output.messages")
	if err := json.Unmarshal([]byte(value.AsString()), &messages); err != nil {
		t.Fatal(err)
	}
	if len(messages) != 1 || len(messages[0].Parts) != 1 || messages[0].Parts[0]["content"] != piece+piece+piece {
		t.Errorf("gen_ai.output.messages has %d bytes, want one message with all of the streamed text", len(value.AsString()))
	}
	value, _ = attr(span, "gen_ai.input.messages")
	messages = nil
	if err := json.Unmarshal([]byte(value.AsString()), &messages); err != nil {
		t.Fatal(err)
	}
	if len(messages) != 5 {
		t.Fatalf("kept %d messages, want all 5", len(messages))
	}
	for i, m := range messages[:3] {
		if len(m.Parts) != 0 {
			t.Errorf("message %d has %d parts, want none", i, len(m.Parts))
		}
	}
	if len(messages[3].Parts) != 1 {
		t.Fatalf("the fourth message has %d parts, want its text alone and no marker", len(messages[3].Parts))
	}
	if text, _ := messages[3].Parts[0]["content"].(string); !strings.HasPrefix(text, "3a") || len(text) != len(big)+1 {
		t.Error("the fourth message does not have its text, whole")
	}
	wantAttrs(t, span, map[string]any{"anthropic.content.truncated": true})
}

// eventsBefore returns the events of body, an event stream, that come before the
// event whose data has marker.
func eventsBefore(body, marker string) string {
	before, _, _ := strings.Cut(body, marker)
	return before[:strings.LastIndex(before, "event: ")]
}

// endingBody is a response body that serves text and then fails its reads with
// err, or, when err is nil, with the error of ctx once ctx ends.
type endingBody struct {
	text io.Reader
	ctx  context.Context
	err  error
}

func (b *endingBody) Read(p []byte) (int, error) {
	if n, err := b.text.Read(p); err != io.EOF {
		return n, err
	}
	if b.err != nil {
		return 0, b.err
	}
	<-b.ctx.Done()
	return 0, b.ctx.Err()
}

func (*endingBody) Close() error { return nil }

// respondAndEnd returns a transport that answers every request with an event
// stream that serves text and then ends as an [endingBody] with err does.
func respondAndEnd(text string, err error) *closureTransport {
	return &closureTransport{fn: func(req *http.Request) (*http.Response, error) {
		body := &endingBody{text: strings.NewReader(text), ctx: req.Context(), err: err}
		return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"text/event-stream"}}, Body: body, Request: req}, nil
	}}
}

// readContentStream reads the Message stream of api ("messages" or "beta") until
// it ends or, when stopAt is not empty, until its first event of that type.
// Then it calls stopped, when it is not nil, reads one more event, closes the
// stream, and returns the error of the stream.
func readContentStream(ctx context.Context, client anthropic.Client, api, stopAt string, stopped func()) error {
	read := func(next func() bool, current func() string, streamErr, closeStream func() error) error {
		for next() {
			if stopAt != "" && current() == stopAt {
				break
			}
		}
		if stopped != nil {
			stopped()
			next()
		}
		err := streamErr()
		closeStream()
		return err
	}
	if api == "beta" {
		stream := client.Beta.Messages.NewStreaming(ctx, betaContentParams())
		return read(stream.Next, func() string { return stream.Current().Type }, stream.Err, stream.Close)
	}
	stream := client.Messages.NewStreaming(ctx, contentParams())
	return read(stream.Next, func() string { return stream.Current().Type }, stream.Err, stream.Close)
}

func TestOpenTelemetryStreamContentEndedBeforeStop(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	afterFirstDelta := eventsBefore(contentStreamBody, `"text":"ny."`)
	inToolCall := eventsBefore(contentStreamBody, `\"Paris\"`)
	errorEvent := "event: error\n" + `data: {"type":"error","error":{"type":"overloaded_error","message":"Overloaded"}}` + "\n\n"
	const firstDelta = `[{"role":"assistant","parts":[{"type":"text","content":"Sun"}]}]`

	for name, test := range map[string]struct {
		transport *closureTransport
		// stopAt is the type of the event at which the application stops reading, and cancel is whether it cancels the call then.
		stopAt string
		cancel bool
		// failed is whether the stream fails.
		failed    bool
		status    codes.Code
		errorType any
		output    string
	}{
		"the application closes the stream after a text delta": {
			transport: respond(200, "text/event-stream", contentStreamBody), stopAt: "content_block_delta",
			status: codes.Unset, output: firstDelta,
		},
		"the application cancels the call after a text delta": {
			transport: respondAndEnd(afterFirstDelta, nil), stopAt: "content_block_delta", cancel: true,
			failed: true, status: codes.Unset, output: firstDelta,
		},
		"the server ends the stream after a text delta": {
			transport: respond(200, "text/event-stream", afterFirstDelta),
			status:    codes.Error, errorType: "api_error", output: firstDelta,
		},
		// The input of a tool call that the stream did not finish is empty.
		"the server ends the stream inside a tool call": {
			transport: respond(200, "text/event-stream", inToolCall),
			status:    codes.Error, errorType: "api_error",
			output: `[{"role":"assistant","parts":[{"type":"text","content":"Sunny."},
				{"type":"tool_call","id":"toolu_9","name":"get_weather","arguments":{}}]}]`,
		},
		"the stream fails with an error event": {
			transport: respond(200, "text/event-stream", afterFirstDelta+errorEvent),
			failed:    true, status: codes.Error, errorType: "overloaded_error", output: firstDelta,
		},
		"a read of the stream fails": {
			transport: respondAndEnd(afterFirstDelta, io.ErrUnexpectedEOF),
			failed:    true, status: codes.Error, errorType: "connection_error", output: firstDelta,
		},
		"the stream ends before a content block": {
			transport: respond(200, "text/event-stream", eventsBefore(contentStreamBody, `"content_block":{"type":"text"`)),
			status:    codes.Error, errorType: "api_error", output: `[{"role":"assistant","parts":[]}]`,
		},
	} {
		for _, api := range []string{"messages", "beta"} {
			t.Run(name+"/"+api, func(t *testing.T) {
				client, recorder := newTracedClient(t, test.transport, content)
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				var stopped func()
				if test.cancel {
					stopped = cancel
				}
				if err := readContentStream(ctx, client, api, test.stopAt, stopped); (err != nil) != test.failed {
					t.Fatalf("the stream ended with the error %v", err)
				}

				span := onlySpan(t, recorder)
				if span.Status().Code != test.status {
					t.Errorf("status %v, want %v", span.Status().Code, test.status)
				}
				wantJSON(t, span, "gen_ai.output.messages", test.output)
				if api == "messages" {
					wantJSON(t, span, "gen_ai.input.messages", contentInputJSON)
				}
				// What only a whole Message has stays off the span.
				wantAttrs(t, span, map[string]any{
					"error.type":                     test.errorType,
					"gen_ai.response.id":             "msg_1",
					"anthropic.message.stop_reason":  nil,
					"gen_ai.response.finish_reasons": nil,
					"gen_ai.usage.output_tokens":     nil,
				})
			})
		}
	}
}

func TestOpenTelemetryStreamContentEndedBeforeStopWithoutAMessage(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	afterFirstDelta := eventsBefore(contentStreamBody, `"text":"ny."`)
	_, withoutStart, _ := strings.Cut(afterFirstDelta, "\n\n")
	rejected := afterFirstDelta + "event: content_block_delta\n" +
		`data: {"type":"content_block_delta","index":7,"delta":{"type":"text_delta","text":"no block"}}` + "\n\n"

	for name, test := range map[string]struct {
		body string
		opts []option.RequestOption
	}{
		"no message_start arrived":           {withoutStart, []option.RequestOption{content}},
		"the accumulator rejected an event":  {rejected, []option.RequestOption{content}},
		"the content is not captured":        {afterFirstDelta, nil},
		"the stream ends before every event": {"", []option.RequestOption{content}},
	} {
		for _, api := range []string{"messages", "beta"} {
			t.Run(name+"/"+api, func(t *testing.T) {
				client, recorder := newTracedClient(t, respond(200, "text/event-stream", test.body), test.opts...)
				if err := readContentStream(context.Background(), client, api, "", nil); err != nil {
					t.Fatal(err)
				}
				span := onlySpan(t, recorder)
				wantAttrs(t, span, map[string]any{"gen_ai.output.messages": nil, "error.type": "api_error"})
			})
		}
	}
}

func TestOpenTelemetryStreamContentEndedBeforeStopThatTheSDKDoesNotAccumulate(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	afterFirstDelta := eventsBefore(contentStreamBody, `"text":"ny."`)

	// The application reads the body itself: no Message is built from it.
	t.Run("body", func(t *testing.T) {
		client, recorder := newTracedClient(t, respond(200, "text/event-stream", contentStreamBody), content)
		var raw *http.Response
		_ = client.Messages.NewStreaming(context.Background(), contentParams(), option.WithResponseBodyInto(&raw))
		if _, err := io.ReadFull(raw.Body, make([]byte, len(afterFirstDelta))); err != nil {
			t.Fatal(err)
		}
		raw.Body.Close()
		span := onlySpan(t, recorder)
		wantJSON(t, span, "gen_ai.input.messages", contentInputJSON)
		wantAttrs(t, span, map[string]any{"gen_ai.output.messages": nil, "gen_ai.response.id": nil})
	})
	// The application decodes the events into a type of its own, which the SDK's accumulator does not take.
	t.Run("events", func(t *testing.T) {
		client, recorder := newTracedClient(t, respond(200, "text/event-stream", afterFirstDelta), content)
		var raw *http.Response
		_ = client.Messages.NewStreaming(context.Background(), contentParams(), option.WithResponseBodyInto(&raw))
		stream := ssestream.NewStream[json.RawMessage](ssestream.NewDecoder(raw), nil)
		events := 0
		for stream.Next() {
			events++
		}
		if err := stream.Err(); err != nil || events != 3 {
			t.Fatalf("the application read %d events and the error %v, want 3 events", events, err)
		}
		span := onlySpan(t, recorder)
		wantJSON(t, span, "gen_ai.input.messages", contentInputJSON)
		wantAttrs(t, span, map[string]any{"gen_ai.output.messages": nil, "error.type": "api_error"})
	})
}

func TestOpenTelemetryContentOnError(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	client, recorder := newTracedClient(t, respond(400, "application/json",
		`{"type":"error","error":{"type":"invalid_request_error","message":"bad"}}`), content)

	if _, err := client.Messages.New(context.Background(), contentParams()); err == nil {
		t.Fatal("the call did not fail")
	}
	span := onlySpan(t, recorder)
	wantJSON(t, span, "gen_ai.input.messages", contentInputJSON)
	wantAttrs(t, span, map[string]any{"gen_ai.output.messages": nil})
}

func TestOpenTelemetryContentSizeBound(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON), content)

	// Five messages of 225 KB each: the parts of the oldest two must go.
	big := strings.Repeat("a", 225_000)
	params := messageParams()
	params.Messages = nil
	for i := 0; i < 5; i++ {
		params.Messages = append(params.Messages, anthropic.NewUserMessage(anthropic.NewTextBlock(string(rune('0'+i))+big)))
	}
	if _, err := client.Messages.New(context.Background(), params); err != nil {
		t.Fatal(err)
	}

	span := onlySpan(t, recorder)
	value, _ := attr(span, "gen_ai.input.messages")
	if n := len(value.AsString()); n > 750_000 {
		t.Errorf("gen_ai.input.messages has %d bytes, more than 0.75 MB", n)
	}
	var messages []struct {
		Parts []map[string]any `json:"parts"`
	}
	if err := json.Unmarshal([]byte(value.AsString()), &messages); err != nil {
		t.Fatal(err)
	}
	if len(messages) != 5 {
		t.Fatalf("kept %d messages, want all 5", len(messages))
	}
	if len(messages[0].Parts) != 0 || len(messages[1].Parts) != 0 {
		t.Errorf("the oldest two messages have %d and %d parts, want none", len(messages[0].Parts), len(messages[1].Parts))
	}
	if len(messages[2].Parts) != 1 {
		t.Fatalf("the third message has %d parts, want its text alone and no marker", len(messages[2].Parts))
	}
	if text, _ := messages[2].Parts[0]["content"].(string); !strings.HasPrefix(text, "2a") || len(text) != len(big)+1 {
		t.Error("the third message does not have its text, whole")
	}
	wantAttrs(t, span, map[string]any{"anthropic.content.truncated": true})
}

func TestOpenTelemetryTracesMaxContentBytes(t *testing.T) {
	maxBytes := func(n int64) []option.RequestOption {
		return []option.RequestOption{option.WithOpenTelemetryTracesMaxContentBytes(n)}
	}
	reply := strings.Repeat("h", 500)

	for name, test := range map[string]struct {
		// clientOpts are options of the client, and callOpts those of the method call.
		clientOpts, callOpts []option.RequestOption
		// env is the value of the environment variable envKey, or of
		// ANTHROPIC_OPEN_TELEMETRY_TRACES_MAX_CONTENT_BYTES when envKey is empty. No
		// variable is set when env is empty.
		env, envKey string
		// want is the maximum that the content is cut to, or 0 when all of it fits
		// the maximum, as it fits the default.
		want int
	}{
		"default":                              {},
		"option":                               {clientOpts: maxBytes(200), want: 200},
		"option of 1":                          {clientOpts: maxBytes(1), want: 1},
		"option of the method call":            {clientOpts: maxBytes(300), callOpts: maxBytes(200), want: 200},
		"environment variable":                 {env: "200", want: 200},
		"option over the environment variable": {clientOpts: maxBytes(200), env: "300", want: 200},
		"an option that all the content fits":  {clientOpts: maxBytes(1 << 20), env: "200"},
		// An option below 1 selects the default: the environment variable is for a client without the option.
		"option of 0":                                     {clientOpts: maxBytes(0), env: "200"},
		"option below 0":                                  {clientOpts: maxBytes(-200), env: "200"},
		"environment variable of 0":                       {env: "0"},
		"environment variable below 0":                    {env: "-200"},
		"environment variable with a fraction":            {env: "200.5"},
		"environment variable that is no number":          {env: "200 bytes"},
		"environment variable that is too large":          {env: "99999999999999999999"},
		"environment variable without traces in its name": {env: "200", envKey: "ANTHROPIC_OPEN_TELEMETRY_MAX_CONTENT_BYTES"},
	} {
		t.Run(name, func(t *testing.T) {
			if test.env != "" {
				t.Setenv(cmp.Or(test.envKey, "ANTHROPIC_OPEN_TELEMETRY_TRACES_MAX_CONTENT_BYTES"), test.env)
			}
			opts := append([]option.RequestOption{option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)}, test.clientOpts...)
			client, recorder := newTracedClient(t, respond(200, "application/json", textBlocksJSON(reply)), opts...)
			if _, err := client.Messages.New(context.Background(), messageParams(), test.callOpts...); err != nil {
				t.Fatal(err)
			}

			span := onlySpan(t, recorder)
			if test.want == 0 {
				wantJSON(t, span, "gen_ai.input.messages", `[{"role":"user","parts":[{"type":"text","content":"hello"}]}]`)
				wantJSON(t, span, "gen_ai.output.messages", `[{"role":"assistant","parts":[{"type":"text","content":"`+reply+`"}]}]`)
				wantAttrs(t, span, map[string]any{"anthropic.content.truncated": nil})
				return
			}
			// The input's part goes, and then the end of the text of the output's part,
			// the last one, which stays: with all of it cut when the maximum has no
			// room for the messages without their parts.
			emptied := len(`[{"role":"user","parts":[]}]`) + len(`[{"role":"assistant","parts":[{"type":"text","content":""}]}]`)
			kept := max(test.want-emptied, 0)
			wantJSON(t, span, "gen_ai.input.messages", `[{"role":"user","parts":[]}]`)
			wantJSON(t, span, "gen_ai.output.messages", `[{"role":"assistant","parts":[{"type":"text","content":"`+strings.Repeat("h", kept)+`"}]}]`)
			wantAttrs(t, span, map[string]any{"anthropic.content.truncated": true})
			if total := contentSize(span); total != emptied+kept {
				t.Errorf("the content attributes have %d bytes, want %d", total, emptied+kept)
			}
		})
	}
}

func TestOpenTelemetryTracesMaxContentBytesBoundsAStreamedMessageToo(t *testing.T) {
	const bound = 160
	opts := []option.RequestOption{
		option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent),
		option.WithOpenTelemetryTracesMaxContentBytes(bound),
	}
	for _, api := range []string{"messages", "beta"} {
		t.Run(api, func(t *testing.T) {
			whole, _ := streamedAndNot(t, api, contentStreamBody, contentStreamMessageJSON, opts[0])
			if total := contentSize(whole); total <= bound {
				t.Fatalf("the content attributes have %d bytes without the maximum, want more than %d", total, bound)
			}
			streamed, unstreamed := streamedAndNot(t, api, contentStreamBody, contentStreamMessageJSON, opts...)
			for _, span := range []sdktrace.ReadOnlySpan{streamed, unstreamed} {
				if total := contentSize(span); total > bound {
					t.Errorf("the content attributes have %d bytes, want at most %d", total, bound)
				}
				wantAttrs(t, span, map[string]any{"anthropic.content.truncated": true})
			}
		})
	}
}

// contentSize returns the bytes of the content attributes of span that the size
// bound can reduce.
func contentSize(span sdktrace.ReadOnlySpan) int {
	total := 0
	for _, key := range []string{"gen_ai.system_instructions", "gen_ai.input.messages", "gen_ai.output.messages"} {
		value, _ := attr(span, key)
		total += len(value.AsString())
	}
	return total
}

func TestOpenTelemetryContentSystemInstructionsOverTheLimitGoFirst(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON), content)

	params := messageParams()
	params.System = []anthropic.TextBlockParam{{Text: strings.Repeat("s", 1_100_000)}, {Text: "Be brief."}}
	params.Messages = []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock("first")),
		anthropic.NewAssistantMessage(anthropic.NewTextBlock("second")),
		anthropic.NewUserMessage(anthropic.NewTextBlock("third")),
	}
	if _, err := client.Messages.New(context.Background(), params); err != nil {
		t.Fatal(err)
	}

	// The first part of the system instructions goes, whole, and everything after it fits.
	span := onlySpan(t, recorder)
	wantJSON(t, span, "gen_ai.system_instructions", `[{"type":"text","content":"Be brief."}]`)
	wantJSON(t, span, "gen_ai.input.messages", `[{"role":"user","parts":[{"type":"text","content":"first"}]},
		{"role":"assistant","parts":[{"type":"text","content":"second"}]},{"role":"user","parts":[{"type":"text","content":"third"}]}]`)
	wantJSON(t, span, "gen_ai.output.messages", `[{"role":"assistant","parts":[{"type":"text","content":"hi"}]}]`)
	wantAttrs(t, span, map[string]any{"anthropic.content.truncated": true})
}

func TestOpenTelemetryContentTextCutWhenTheCallStartsGoesWholeOnceThereIsOutput(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	params := messageParams()
	params.Messages = []anthropic.MessageParam{anthropic.NewUserMessage(anthropic.NewTextBlock(strings.Repeat("h", 1_100_000)))}

	// When the call starts, the text is the last part of the content, so its end
	// goes, and a call that fails has no more than that.
	client, recorder := newTracedClient(t, respond(400, "application/json",
		`{"type":"error","error":{"type":"invalid_request_error","message":"bad"}}`), content)
	if _, err := client.Messages.New(context.Background(), params); err == nil {
		t.Fatal("the call did not fail")
	}
	span := onlySpan(t, recorder)
	kept := 750_000 - len(`[{"role":"user","parts":[{"type":"text","content":""}]}]`)
	wantJSON(t, span, "gen_ai.input.messages", `[{"role":"user","parts":[{"type":"text","content":"`+strings.Repeat("h", kept)+`"}]}]`)
	wantAttrs(t, span, map[string]any{"gen_ai.output.messages": nil, "anthropic.content.truncated": true})

	// Beside output the text is no longer the last part: it goes whole, its message
	// stays, and the output fits.
	client, recorder = newTracedClient(t, respond(200, "application/json", messageJSON), content)
	if _, err := client.Messages.New(context.Background(), params); err != nil {
		t.Fatal(err)
	}
	span = onlySpan(t, recorder)
	wantJSON(t, span, "gen_ai.input.messages", `[{"role":"user","parts":[]}]`)
	wantJSON(t, span, "gen_ai.output.messages", `[{"role":"assistant","parts":[{"type":"text","content":"hi"}]}]`)
	wantAttrs(t, span, map[string]any{"anthropic.content.truncated": true})
}

// textBlocksJSON is a Message whose content is one text block for each of texts.
func textBlocksJSON(texts ...string) string {
	blocks := make([]string, len(texts))
	for i, text := range texts {
		blocks[i] = `{"type":"text","text":"` + text + `"}`
	}
	return `{"id":"msg_1","type":"message","role":"assistant","model":"claude-served","stop_reason":"end_turn",
		"content":[` + strings.Join(blocks, ",") + `],"usage":{"input_tokens":1,"output_tokens":1}}`
}

// textBlocksStreamBody is the stream of a Message whose content is one text block for each of texts.
func textBlocksStreamBody(texts ...string) string {
	var events strings.Builder
	events.WriteString("event: message_start\n" +
		`data: {"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","model":"claude-served","content":[],"usage":{"input_tokens":1,"output_tokens":1}}}` + "\n\n")
	for i, text := range texts {
		index := strconv.Itoa(i)
		events.WriteString("event: content_block_start\n" + `data: {"type":"content_block_start","index":` + index + `,"content_block":{"type":"text","text":""}}` + "\n\n" +
			"event: content_block_delta\n" + `data: {"type":"content_block_delta","index":` + index + `,"delta":{"type":"text_delta","text":"` + text + `"}}` + "\n\n" +
			"event: content_block_stop\n" + `data: {"type":"content_block_stop","index":` + index + `}` + "\n\n")
	}
	events.WriteString("event: message_delta\n" + `data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":2}}` + "\n\n" +
		"event: message_stop\n" + `data: {"type":"message_stop"}` + "\n\n")
	return events.String()
}

// wantPartCounts fails the test unless the attribute key of span holds one
// message for each of want, with that many parts.
func wantPartCounts(t *testing.T, span sdktrace.ReadOnlySpan, key string, want ...int) {
	t.Helper()
	value, _ := attr(span, key)
	var messages []struct {
		Parts []map[string]any `json:"parts"`
	}
	if err := json.Unmarshal([]byte(value.AsString()), &messages); err != nil {
		t.Fatalf("%s is not JSON: %v", key, err)
	}
	got := make([]int, len(messages))
	for i, m := range messages {
		if m.Parts == nil {
			t.Errorf("%s: message %d has no list of parts", key, i)
		}
		got[i] = len(m.Parts)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s: the messages have %v parts, want %v", key, got, want)
	}
}

// wantTextParts fails the test unless the parts of the messages that the
// attribute key of span holds, or for gen_ai.output.messages of its one message,
// are text parts whose texts begin with want, in order, each of size bytes.
func wantTextParts(t *testing.T, span sdktrace.ReadOnlySpan, key string, size int, want ...string) {
	t.Helper()
	value, _ := attr(span, key)
	var messages []struct {
		Parts []map[string]any `json:"parts"`
	}
	if err := json.Unmarshal([]byte(value.AsString()), &messages); err != nil {
		t.Fatalf("%s is not JSON: %v", key, err)
	}
	var parts []map[string]any
	for _, m := range messages {
		parts = append(parts, m.Parts...)
	}
	if key == "gen_ai.output.messages" && len(messages) != 1 {
		t.Fatalf("%s has %d messages, want 1", key, len(messages))
	}
	if len(parts) != len(want) {
		t.Fatalf("%s has %d parts, want %d", key, len(parts), len(want))
	}
	for i, p := range parts {
		if text, _ := p["content"].(string); p["type"] != "text" || !strings.HasPrefix(text, want[i]) || len(text) != size {
			t.Errorf("%s: part %d is not the text that begins with %q, whole", key, i, want[i])
		}
	}
}

func TestOpenTelemetryContentSizeBoundRemovesOutputPartsFromTheStart(t *testing.T) {
	// Two input messages of 300 KB and three output parts of 350 KB. The system
	// instructions go first, then the parts of both input messages, because the
	// output alone is over the limit, and then the first part of the output.
	big, piece := strings.Repeat("a", 300_000), strings.Repeat("o", 350_000)
	params := messageParams()
	params.System = []anthropic.TextBlockParam{{Text: "Be brief."}}
	params.Messages = []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock("0" + big)),
		anthropic.NewUserMessage(anthropic.NewTextBlock("1" + big)),
	}
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	check := func(t *testing.T, span sdktrace.ReadOnlySpan) {
		t.Helper()
		wantJSON(t, span, "gen_ai.system_instructions", `[]`)
		wantJSON(t, span, "gen_ai.input.messages", `[{"role":"user","parts":[]},{"role":"user","parts":[]}]`)
		wantTextParts(t, span, "gen_ai.output.messages", len(piece)+1, "1o", "2o")
		wantAttrs(t, span, map[string]any{"anthropic.content.truncated": true})
	}

	t.Run("message", func(t *testing.T) {
		client, recorder := newTracedClient(t, respond(200, "application/json", textBlocksJSON("0"+piece, "1"+piece, "2"+piece)), content)
		if _, err := client.Messages.New(context.Background(), params); err != nil {
			t.Fatal(err)
		}
		check(t, onlySpan(t, recorder))
	})
	t.Run("stream", func(t *testing.T) {
		client, recorder := newTracedClient(t, respond(200, "text/event-stream", textBlocksStreamBody("0"+piece, "1"+piece, "2"+piece)), content)
		stream := client.Messages.NewStreaming(context.Background(), params)
		for stream.Next() {
		}
		if err := stream.Err(); err != nil {
			t.Fatal(err)
		}
		check(t, onlySpan(t, recorder))
	})
}

func TestOpenTelemetryContentRemovesAToolResultOverTheLimitWhole(t *testing.T) {
	// The blocks of a tool result count toward the size of their part, so a search
	// result of more than 0.75 MB takes its whole part with it, and only that: the
	// text after it stays in the message, and so does the next message.
	params := messageParams()
	params.Messages = []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.ContentBlockParamUnion{OfToolResult: &anthropic.ToolResultBlockParam{ToolUseID: "toolu_1", Content: []anthropic.ToolResultBlockParamContentUnion{
			{OfText: &anthropic.TextBlockParam{Text: "One result."}},
			{OfSearchResult: &anthropic.SearchResultBlockParam{Source: "https://example.com", Title: "Forecast",
				Content: []anthropic.TextBlockParam{{Text: strings.Repeat("a", 750_000)}}}},
		}}}, anthropic.NewTextBlock("And today?")),
		anthropic.NewUserMessage(anthropic.NewTextBlock("Tomorrow?")),
	}
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON), content)
	if _, err := client.Messages.New(context.Background(), params); err != nil {
		t.Fatal(err)
	}
	span := onlySpan(t, recorder)
	wantJSON(t, span, "gen_ai.input.messages", `[{"role":"user","parts":[{"type":"text","content":"And today?"}]},{"role":"user","parts":[{"type":"text","content":"Tomorrow?"}]}]`)
	wantJSON(t, span, "gen_ai.output.messages", `[{"role":"assistant","parts":[{"type":"text","content":"hi"}]}]`)
	wantAttrs(t, span, map[string]any{"anthropic.content.truncated": true})
}

func TestOpenTelemetryContentRemovesAPartWhoseResultObjectIsOverTheLimit(t *testing.T) {
	// A result that is one object counts whole toward the size of its part, so one
	// of more than 0.75 MB goes with its whole part, and nothing of it is kept.
	result := `{"type":"code_execution_tool_result","tool_use_id":"srvtoolu_1","content":{"type":"code_execution_result","stdout":"` + strings.Repeat("a", 750_000) + `","stderr":"","return_code":0,"content":[]}}`
	response := func(blocks string) string {
		return `{"id":"msg_1","type":"message","role":"assistant","model":"claude-served","content":[` + blocks + `],
			"stop_reason":"end_turn","usage":{"input_tokens":1,"output_tokens":1}}`
	}
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	for name, test := range map[string]struct{ blocks, want string }{
		// The input's part goes before any part of the output does.
		"before text, it goes and the text stays": {result + `,{"type":"text","text":"Done."}`, `[{"role":"assistant","parts":[{"type":"text","content":"Done."}]}]`},
		// The text before it goes first. The result is then the last part, which is
		// not text, so it goes too.
		"as the last part, it goes as every part before it": {`{"type":"text","text":"Running."},` + result, `[{"role":"assistant","parts":[]}]`},
	} {
		t.Run(name, func(t *testing.T) {
			client, recorder := newTracedClient(t, respond(200, "application/json", response(test.blocks)), content)
			if _, err := client.Messages.New(context.Background(), messageParams()); err != nil {
				t.Fatal(err)
			}
			span := onlySpan(t, recorder)
			wantJSON(t, span, "gen_ai.input.messages", `[{"role":"user","parts":[]}]`)
			wantJSON(t, span, "gen_ai.output.messages", test.want)
			wantAttrs(t, span, map[string]any{"anthropic.content.truncated": true})
		})
	}
}

func TestOpenTelemetryContentInputIsFittedAgainBesideTheOutput(t *testing.T) {
	// 700 KB of input fits when the call starts. The 60 KB of output do not fit beside
	// it, so the request's content is fitted again when the response arrives: the
	// system instructions go, then the part of the oldest message, and the rest and
	// the output are as they were.
	piece, reply := strings.Repeat("a", 70_000), strings.Repeat("o", 60_000)
	params := messageParams()
	params.System = []anthropic.TextBlockParam{{Text: "Be <brief> & kind."}}
	params.Messages = nil
	for i := 0; i < 10; i++ {
		params.Messages = append(params.Messages, anthropic.NewUserMessage(anthropic.NewTextBlock(string(rune('0'+i))+piece)))
	}
	params.Messages = append(params.Messages,
		anthropic.NewAssistantMessage(anthropic.NewToolUseBlock("toolu_1", map[string]any{"q": "<a&b>"}, "get_weather")),
		anthropic.NewUserMessage(anthropic.NewToolResultBlock("toolu_1", "sunny", false)))
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	client, recorder := newTracedClient(t, respond(200, "application/json", textBlocksJSON(reply)), content)
	if _, err := client.Messages.New(context.Background(), params); err != nil {
		t.Fatal(err)
	}

	span := onlySpan(t, recorder)
	input, _ := attr(span, "gen_ai.input.messages")
	var want strings.Builder
	want.WriteString(`[{"role":"user","parts":[]},`)
	for i := 1; i < 10; i++ {
		want.WriteString(`{"role":"user","parts":[{"type":"text","content":"` + string(rune('0'+i)) + piece + `"}]},`)
	}
	// The values are as typed, without the escapes of the default JSON encoder.
	want.WriteString(`{"role":"assistant","parts":[{"type":"tool_call","id":"toolu_1","name":"get_weather","arguments":{"q":"<a&b>"}}]},` +
		`{"role":"user","parts":[{"type":"tool_call_response","id":"toolu_1","response":[{"text":"sunny","type":"text"}]}]}]`)
	if input.AsString() != want.String() {
		t.Errorf("gen_ai.input.messages has %d bytes, want the oldest message without its part and the newest eleven byte for byte in %d", len(input.AsString()), want.Len())
	}
	wantAttrs(t, span, map[string]any{
		"gen_ai.system_instructions":  `[]`,
		"gen_ai.output.messages":      `[{"role":"assistant","parts":[{"type":"text","content":"` + reply + `"}]}]`,
		"anthropic.content.truncated": true,
	})
}

func TestOpenTelemetryContentTheLastTextPartOverTheLimitIsCutFromItsEnd(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	client, recorder := newTracedClient(t, respond(200, "application/json", textBlocksJSON(strings.Repeat("h", 1_100_000))), content)
	if _, err := client.Messages.New(context.Background(), messageParams()); err != nil {
		t.Fatal(err)
	}

	// The input's part goes first, though that does not help, and then the end
	// of the text of the output's part, the last one, which stays.
	span := onlySpan(t, recorder)
	wantJSON(t, span, "gen_ai.input.messages", `[{"role":"user","parts":[]}]`)
	kept := 750_000 - len(`[{"role":"user","parts":[]}]`) - len(`[{"role":"assistant","parts":[{"type":"text","content":""}]}]`)
	wantJSON(t, span, "gen_ai.output.messages", `[{"role":"assistant","parts":[{"type":"text","content":"`+strings.Repeat("h", kept)+`"}]}]`)
	wantAttrs(t, span, map[string]any{"anthropic.content.truncated": true})
	if total := contentSize(span); total != 750_000 {
		t.Errorf("the content attributes have %d bytes, want 0.75 MB", total)
	}
}

func TestOpenTelemetryContentSizeBoundRemovesAToolCallWholeAndCutsTextBetweenCharacters(t *testing.T) {
	// The tool call goes whole, and then the end of the text after it, the last
	// part, whose characters are four bytes each.
	text := strings.Repeat("\U0001f600", 300_000)
	response := `{"id":"msg_1","type":"message","role":"assistant","model":"claude-served","content":[
		{"type":"tool_use","id":"toolu_1","name":"get_weather","input":{"city":"Paris"}},{"type":"text","text":"` + text + `"}],
		"stop_reason":"end_turn","usage":{"input_tokens":1,"output_tokens":1}}`
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	client, recorder := newTracedClient(t, respond(200, "application/json", response), content)
	if _, err := client.Messages.New(context.Background(), messageParams()); err != nil {
		t.Fatal(err)
	}
	span := onlySpan(t, recorder)
	wantJSON(t, span, "gen_ai.input.messages", `[{"role":"user","parts":[]}]`)
	room := 750_000 - len(`[{"role":"user","parts":[]}]`) - len(`[{"role":"assistant","parts":[{"type":"text","content":""}]}]`)
	wantJSON(t, span, "gen_ai.output.messages", `[{"role":"assistant","parts":[{"type":"text","content":"`+text[:room-room%4]+`"}]}]`)
	if total := contentSize(span); total > 750_000 || total <= 750_000-4 {
		t.Errorf("the content attributes have %d bytes, want the most whole characters that fit in 0.75 MB", total)
	}
	wantAttrs(t, span, map[string]any{"anthropic.content.truncated": true})
}

func TestOpenTelemetryContentNoInputMessagesAreAnEmptyList(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	client, recorder := newTracedClient(t, respond(200, "application/json", messageJSON), content)

	params := messageParams()
	params.System = []anthropic.TextBlockParam{{Text: strings.Repeat("s", 1_100_000)}}
	params.Messages = []anthropic.MessageParam{}
	if _, err := client.Messages.New(context.Background(), params); err != nil {
		t.Fatal(err)
	}

	span := onlySpan(t, recorder)
	if value, _ := attr(span, "gen_ai.input.messages"); value.AsString() != "[]" {
		t.Errorf("gen_ai.input.messages = %s, want []", value.AsString())
	}
	// The system instructions are over the limit, so their part went, and the
	// emptied list stays. The output fits then.
	wantJSON(t, span, "gen_ai.system_instructions", `[]`)
	wantJSON(t, span, "gen_ai.output.messages", `[{"role":"assistant","parts":[{"type":"text","content":"hi"}]}]`)
	wantAttrs(t, span, map[string]any{"anthropic.content.truncated": true})
}

func TestOpenTelemetryStreamContentMalformedIndex(t *testing.T) {
	body := "event: message_start\n" +
		`data: {"type":"message_start","message":{"id":"msg_1","type":"message","role":"assistant","model":"claude-served","content":[],"usage":{"input_tokens":1,"output_tokens":1}}}` + "\n\n"
	events := []string{
		`{"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`,
		`{"type":"content_block_start","index":-1,"content_block":{"type":"text","text":"negative"}}`,
		`{"type":"content_block_start","index":4000000000,"content_block":{"type":"text","text":"huge"}}`,
		`{"type":"content_block_delta","index":7,"delta":{"type":"text_delta","text":"no block"}}`,
		`{"type":"content_block_delta","index":-1,"delta":{"type":"text_delta","text":"negative"}}`,
		`{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Sunny."}}`,
	}
	for _, event := range events {
		kind := "content_block_delta"
		if strings.Contains(event, "content_block_start") {
			kind = "content_block_start"
		}
		body += "event: " + kind + "\ndata: " + event + "\n\n"
	}
	body += "event: message_delta\n" + `data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":2}}` + "\n\n" +
		"event: message_stop\n" + `data: {"type":"message_stop"}` + "\n\n"
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	client, recorder := newTracedClient(t, respond(200, "text/event-stream", body), content)

	stream := client.Messages.NewStreaming(context.Background(), contentParams())
	var indexes []int64
	for stream.Next() {
		if event := stream.Current(); strings.HasPrefix(event.Type, "content_block_") {
			indexes = append(indexes, event.Index)
		}
	}
	if err := stream.Err(); err != nil {
		t.Fatal(err)
	}
	// The application gets every event, to take or leave. The span records no
	// content from a stream that the SDK's accumulator rejects.
	if want := []int64{0, -1, 4000000000, 7, -1, 0}; !reflect.DeepEqual(indexes, want) {
		t.Errorf("the application got content block events with indexes %v, want %v", indexes, want)
	}
	span := onlySpan(t, recorder)
	if span.Status().Code != codes.Ok {
		t.Errorf("status %v, want OK", span.Status().Code)
	}
	wantJSON(t, span, "gen_ai.input.messages", contentInputJSON)
	wantAttrs(t, span, map[string]any{
		"gen_ai.output.messages":        nil,
		"gen_ai.response.id":            "msg_1",
		"anthropic.message.stop_reason": "end_turn",
		"gen_ai.usage.output_tokens":    int64(2),
	})
}

const markupText = `<b>R&D</b>`

// markupResponseJSON is a Message whose text and tool input have the characters
// that the default JSON encoder of Go writes as escapes.
const markupResponseJSON = `{"id":"msg_1","type":"message","role":"assistant","model":"claude-served",
	"content":[{"type":"text","text":"<b>R&D</b>"},{"type":"tool_use","id":"toolu_9","name":"get_weather","input":{"city":"<b>R&D</b>"}}],
	"stop_reason":"tool_use","usage":{"input_tokens":1,"output_tokens":1}}`

// wantMarkupAsTyped checks that the content attribute key of span has markupText
// want times, as it is typed, and no escape in its place.
func wantMarkupAsTyped(t *testing.T, span sdktrace.ReadOnlySpan, key string, want int) {
	t.Helper()
	got, _ := attr(span, key)
	if count := strings.Count(got.AsString(), markupText); count != want {
		t.Errorf("%s: %s has %s %d times, want %d: %s", span.Name(), key, markupText, count, want, got.AsString())
	}
	for _, escape := range []string{"\\u003c", "\\u003e", "\\u0026"} {
		if strings.Contains(strings.ToLower(got.AsString()), escape) {
			t.Errorf("%s: %s has the escape %s: %s", span.Name(), key, escape, got.AsString())
		}
	}
}

func TestOpenTelemetryContentIsNotHTMLEscaped(t *testing.T) {
	content := option.WithOpenTelemetryTracesContentMode(option.OpenTelemetryContentModeContent)
	params := anthropic.MessageNewParams{
		Model:     "claude-requested",
		MaxTokens: 100,
		System:    []anthropic.TextBlockParam{{Text: markupText}},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(markupText)),
			// The note is a backslash and five characters, not an escape.
			anthropic.NewAssistantMessage(anthropic.NewToolUseBlock("toolu_1", map[string]any{"city": markupText, "note": "\\u003c"}, "get_weather")),
			anthropic.NewUserMessage(anthropic.NewToolResultBlock("toolu_1", markupText, false)),
		},
	}
	client, recorder := newTracedClient(t, respond(200, "application/json", markupResponseJSON), content)
	if _, err := client.Messages.New(context.Background(), params); err != nil {
		t.Fatal(err)
	}
	span := onlySpan(t, recorder)
	wantMarkupAsTyped(t, span, "gen_ai.system_instructions", 1)
	wantJSON(t, span, "gen_ai.input.messages", `[
		{"role":"user","parts":[{"type":"text","content":"<b>R&D</b>"}]},
		{"role":"assistant","parts":[{"type":"tool_call","id":"toolu_1","name":"get_weather","arguments":{"city":"<b>R&D</b>","note":"\\u003c"}}]},
		{"role":"user","parts":[{"type":"tool_call_response","id":"toolu_1","response":[{"type":"text","text":"<b>R&D</b>"}]}]}]`)
	wantMarkupAsTyped(t, span, "gen_ai.output.messages", 2)
}
