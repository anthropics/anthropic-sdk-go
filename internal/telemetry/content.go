package telemetry

import (
	"encoding/json"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel/attribute"
)

// part is one message part in the content attributes: a [textPart], a
// [reasoningPart], a [toolCallRequestPart], a [toolCallResponsePart], or a
// [genericPart]. They follow the message parts of the generative-AI conventions.
type part interface{ messagePart() }

// textPart is a text content block.
type textPart struct {
	Type string `json:"type"` // "text"
	// Content is the block's text.
	Content string `json:"content"`
}

// reasoningPart is a thinking content block. It has no signature.
type reasoningPart struct {
	Type string `json:"type"` // "reasoning"
	// Content is the block's thinking, which is empty for a hidden block.
	Content string `json:"content"`
}

// toolCallRequestPart is a tool-use content block: tool_use, server_tool_use, or
// mcp_tool_use.
type toolCallRequestPart struct {
	Type string `json:"type"` // "tool_call"
	ID   string `json:"id"`
	Name string `json:"name"`
	// Arguments is the block's input, a JSON object.
	Arguments json.RawMessage `json:"arguments"`
	// MCPServerName is the server_name of an mcp_tool_use block.
	MCPServerName string `json:"anthropic_mcp_server_name,omitempty"`
	// ToolsetName is the toolset_name of a tool_use block that has one.
	ToolsetName string `json:"anthropic_toolset_name,omitempty"`
}

// toolCallResponsePart is a result content block: tool_result, mcp_tool_result,
// or the result block of a server tool, such as web_search_tool_result.
type toolCallResponsePart struct {
	Type string `json:"type"` // "tool_call_response"
	// ID is the block's tool_use_id.
	ID string `json:"id"`
	// Response is the block's content, as [toolResultContent] returns it for a
	// tool_result block, and as [resultBlockContent] does for any other.
	Response any `json:"response"`
}

// genericPart is any other content block. It holds the block's type, with the
// prefix anthropic_, and nothing else. Image and document data is left out
// because it is large and binary.
type genericPart struct {
	Type string `json:"type"`
}

func (textPart) messagePart()             {}
func (reasoningPart) messagePart()        {}
func (toolCallRequestPart) messagePart()  {}
func (toolCallResponsePart) messagePart() {}
func (genericPart) messagePart()          {}

// blockType returns the type of a content block of the API, or "unknown" when
// the block has no type that is a string.
func blockType(block gjson.Result) string {
	typ := block.Get("type")
	if typ.Type != gjson.String {
		return "unknown"
	}
	return typ.String()
}

// blockPart returns the message part for a content block of the API.
func blockPart(block gjson.Result) part {
	typ := blockType(block)
	switch typ {
	case "text":
		return textPart{Type: "text", Content: block.Get("text").String()}
	case "thinking":
		return reasoningPart{Type: "reasoning", Content: block.Get("thinking").String()}
	case "tool_use", "server_tool_use", "mcp_tool_use":
		call := toolCallRequestPart{
			Type:      "tool_call",
			ID:        block.Get("id").String(),
			Name:      block.Get("name").String(),
			Arguments: toolCallArguments(block.Get("input")),
		}
		switch typ {
		case "mcp_tool_use":
			call.MCPServerName = block.Get("server_name").String()
		case "tool_use":
			// A null toolset_name is no toolset name.
			if toolsetName := block.Get("toolset_name"); toolsetName.Type == gjson.String {
				call.ToolsetName = toolsetName.String()
			}
		}
		return call
	case "tool_result":
		return toolCallResponsePart{
			Type:     "tool_call_response",
			ID:       block.Get("tool_use_id").String(),
			Response: toolResultContent(block.Get("content")),
		}
	}
	if strings.HasSuffix(typ, "_tool_result") {
		return toolCallResponsePart{
			Type:     "tool_call_response",
			ID:       block.Get("tool_use_id").String(),
			Response: resultBlockContent(block.Get("content")),
		}
	}
	return genericPart{Type: "anthropic_" + typ}
}

// textBlock is a text content block of the API.
type textBlock struct {
	Type string `json:"type"` // "text"
	Text string `json:"text"`
}

// typeOnlyBlock is a content block of the API that holds its type and nothing
// else.
type typeOnlyBlock struct {
	Type string `json:"type"`
}

// toolResultContent returns the content blocks of content, the content of a tool
// result: a string, a list of content blocks, or nothing. They are blocks of the
// API, not message parts. A string is one text block. A text block and a
// search_result block keep all their fields except those whose name starts with
// encrypted_. Any other block holds its type only, so the mapping does not
// recurse. Image and document data is left out because it is large and binary. A
// result without content, or whose content is null, has no blocks.
func toolResultContent(content gjson.Result) []any {
	if content.Type == gjson.String {
		return []any{textBlock{Type: "text", Text: content.String()}}
	}
	blocks := []any{}
	// ForEach passes a value that is no list to the function as it is.
	if !content.IsArray() {
		return blocks
	}
	content.ForEach(func(_, block gjson.Result) bool {
		switch typ := blockType(block); typ {
		case "text", "search_result":
			blocks = append(blocks, exportedValue(block))
		default:
			blocks = append(blocks, typeOnlyBlock{Type: typ})
		}
		return true
	})
	return blocks
}

// resultBlockContent returns the response of a result block other than
// tool_result, whose content is content: a string, a list, one object, or
// nothing. A string is one text block. In a list, a text block, a search_result
// block and a web_search_result are kept, and any other element holds its type
// only. One object, such as a code_execution_result, is kept. What is kept is as
// the API has it, without what [withoutUnexportedFields] leaves out. A block
// without content, or whose content is null, has an empty list.
func resultBlockContent(content gjson.Result) any {
	if content.Type == gjson.String {
		return []any{textBlock{Type: "text", Text: content.String()}}
	}
	if content.IsObject() {
		return exportedValue(content)
	}
	blocks := []any{}
	if !content.IsArray() {
		return blocks
	}
	content.ForEach(func(_, element gjson.Result) bool {
		switch typ := blockType(element); typ {
		case "text", "search_result", "web_search_result":
			blocks = append(blocks, exportedValue(element))
		default:
			blocks = append(blocks, typeOnlyBlock{Type: typ})
		}
		return true
	})
	return blocks
}

// exportedValue returns the JSON object value, a block or an object inside the
// content of a result, without what [withoutUnexportedFields] leaves out. An
// object that does not decode holds its type only.
func exportedValue(value gjson.Result) any {
	decoder := json.NewDecoder(strings.NewReader(value.Raw))
	// A number stays as the API wrote it.
	decoder.UseNumber()
	var decoded any
	if decoder.Decode(&decoded) != nil {
		return typeOnlyBlock{Type: blockType(value)}
	}
	return withoutUnexportedFields(decoded)
}

// withoutUnexportedFields returns value, a decoded JSON value inside the content
// of a result, without two kinds of field, at any depth. A field whose name
// starts with encrypted_ is left out because the API returns it as an opaque
// string. Image and document data is left out because it is large and binary: an
// image or document block holds its type only, and a
// text_editor_code_execution_view_result of an image or a PDF has no content.
func withoutUnexportedFields(value any) any {
	switch value := value.(type) {
	case []any:
		for i, element := range value {
			value[i] = withoutUnexportedFields(element)
		}
	case map[string]any:
		typ, _ := value["type"].(string)
		if typ == "image" || typ == "document" {
			return typeOnlyBlock{Type: typ}
		}
		if fileType, _ := value["file_type"].(string); typ == "text_editor_code_execution_view_result" && (fileType == "image" || fileType == "pdf") {
			delete(value, "content")
		}
		for name, field := range value {
			if strings.HasPrefix(name, "encrypted_") {
				delete(value, name)
				continue
			}
			value[name] = withoutUnexportedFields(field)
		}
	}
	return value
}

// toolCallArguments returns the arguments of a tool call whose tool_use block
// has input. The request parameters leave a nil input out, and the arguments
// are an empty object then.
func toolCallArguments(input gjson.Result) json.RawMessage {
	if !input.Exists() {
		return json.RawMessage("{}")
	}
	return json.RawMessage(withoutHTMLEscapes(input.Raw))
}

// withoutHTMLEscapes returns raw, a JSON value, with <, > and & in place of the
// escapes that the default JSON encoder of Go writes for them. The SDK encodes a
// request with that encoder, so the JSON of a tool's input has them.
func withoutHTMLEscapes(raw string) string {
	if !strings.Contains(raw, `\u00`) {
		return raw
	}
	var typed strings.Builder
	for i := 0; i < len(raw); i++ {
		if raw[i] != '\\' || i+1 == len(raw) {
			typed.WriteByte(raw[i])
			continue
		}
		if raw[i+1] == 'u' && i+6 <= len(raw) {
			if char, ok := htmlEscapes[strings.ToLower(raw[i+2:i+6])]; ok {
				typed.WriteByte(char)
				i += 5
				continue
			}
		}
		// Any other escape is two bytes or more, and an escaped backslash starts none.
		typed.WriteString(raw[i : i+2])
		i++
	}
	return typed.String()
}

var htmlEscapes = map[string]byte{"003c": '<', "003e": '>', "0026": '&'}

// contentParts returns the message parts of content, which is a string or a
// list of content blocks.
func contentParts(content gjson.Result) []part {
	if content.Type == gjson.String {
		return []part{textPart{Type: "text", Content: content.String()}}
	}
	parts := []part{}
	content.ForEach(func(_, block gjson.Result) bool {
		parts = append(parts, blockPart(block))
		return true
	})
	return parts
}

// marshal returns the JSON of v, a value of the content attributes, with <, >
// and & as they are: a span is not HTML. It is "" when v cannot be encoded,
// which a message part whose arguments are not JSON would cause.
func marshal(v any) string {
	var data strings.Builder
	encoder := json.NewEncoder(&data)
	encoder.SetEscapeHTML(false)
	if encoder.Encode(v) != nil {
		return ""
	}
	return strings.TrimSuffix(data.String(), "\n")
}

// messageContent is one message of gen_ai.input.messages or
// gen_ai.output.messages: the JSON of the message up to its parts, and the JSON
// of each part.
type messageContent struct {
	prefix string
	parts  []string
}

// The JSON of a message after its parts.
const messageSuffix = `]}`

// newMessageContent serializes the message whose role is role, and each of
// parts once. A part that cannot be encoded is left out, so that the message
// stays JSON.
func newMessageContent(role string, parts []part) messageContent {
	return messageContent{prefix: `{"role":` + marshal(role) + `,"parts":[`, parts: marshalParts(parts)}
}

// inputContent is the serialized content of the parameters of a call to the
// Messages create method: the JSON of each part of gen_ai.system_instructions,
// when the request has a system prompt, each input message with the JSON of
// each of its parts, and the request's stop sequences. The span keeps it from
// when the call starts, whole, so that fitting the request's content again
// beside the output starts from every part, and does not parse and serialize
// the request a second time.
type inputContent struct {
	system        []string
	hasSystem     bool
	messages      []messageContent
	stopSequences []string
}

// newInputContent serializes the parameters system and messages, which are
// those JSON values: each part of the system instructions, and each part of
// each message, once. It keeps a copy of stopSequences, the parameter of that
// name, when all of them are valid UTF-8.
func newInputContent(system, messages gjson.Result, stopSequences []string) inputContent {
	var input inputContent
	// A string that is not valid UTF-8 can fail the whole export, so stop
	// sequences that have one are left out, and only they are.
	if !slices.ContainsFunc(stopSequences, func(s string) bool { return !utf8.ValidString(s) }) {
		input.stopSequences = slices.Clone(stopSequences)
	}
	if system.Exists() {
		input.system = marshalParts(contentParts(system))
		input.hasSystem = true
	}
	messages.ForEach(func(_, m gjson.Result) bool {
		parts := contentParts(m.Get("content"))
		// A message that has a part that cannot be encoded is left out.
		if message := newMessageContent(m.Get("role").String(), parts); len(message.parts) == len(parts) {
			input.messages = append(input.messages, message)
		}
		return true
	})
	return input
}

// attributes returns the content attributes of the call, and the length of
// their values: gen_ai.request.stop_sequences, gen_ai.system_instructions,
// gen_ai.input.messages, and, when hasOutput, gen_ai.output.messages, which
// holds the messages output. The stop sequences are a list of strings, not
// JSON, and their length is that of the strings. The attribute is absent when
// the request has none.
//
// Together the values stay within limit bytes, apart from the structure that
// is kept. Stop sequences that are over the limit on their own are left out
// whole, and the attribute with them. Then the parts of the system
// instructions, then those of each input message, then those of each output
// message, are one list, and whole parts go from its start while the values are
// over the limit and more than one part is left. When the one part left is a
// text part, the end of its content goes, as [cutText] cuts it, and the part
// stays; any other part goes whole. Nothing else is removed: a message stays
// with no parts, and the system instructions stay as an empty list, so those and
// the last text part can keep the values over a small limit. When anything was
// removed, anthropic.content.truncated is among the attributes; it is true or
// absent, never false.
func (input inputContent) attributes(output []messageContent, hasOutput bool, limit int) (attrs []attribute.KeyValue, size int) {
	truncated := false
	if len(input.stopSequences) > 0 {
		stopSequencesSize := 0
		for _, sequence := range input.stopSequences {
			stopSequencesSize += len(sequence)
		}
		if stopSequencesSize > limit {
			truncated = true
		} else {
			size += stopSequencesSize
			attrs = append(attrs, GenAIRequestStopSequencesKey.StringSlice(input.stopSequences))
		}
	}
	size += messagesSize(input.messages)
	messages := input.messages
	if hasOutput {
		size += messagesSize(output)
		messages = slices.Concat(input.messages, output)
	}
	// kept is the JSON of the parts that stay, of the system instructions and of
	// each message, in the order in which the parts go, and left is how many they are.
	kept := make([][]string, 0, 1+len(messages))
	if input.hasSystem {
		kept = append(kept, input.system)
		size += listSize(input.system, "[", "]")
	}
	for _, m := range messages {
		kept = append(kept, m.parts)
	}
	left := 0
	for _, parts := range kept {
		left += len(parts)
	}

fit:
	for i, parts := range kept {
		for _, j := range parts {
			if size <= limit {
				break fit
			}
			if left == 1 && gjson.Get(j, "type").String() == "text" {
				cut := cutText(gjson.Get(j, "content").String(), limit-(size-len(j)))
				kept[i] = []string{cut}
				size += len(cut) - len(j)
				truncated = truncated || cut != j
				break fit
			}
			kept[i] = kept[i][1:]
			size -= len(j)
			// The comma after it, unless it was the last.
			if len(kept[i]) > 0 {
				size--
			}
			left--
			truncated = true
		}
	}

	if input.hasSystem {
		attrs = append(attrs, GenAISystemInstructionsKey.String("["+strings.Join(kept[0], ",")+"]"))
		kept = kept[1:]
	}
	attrs = append(attrs, GenAIInputMessagesKey.String(messagesJSON(input.messages, kept)))
	if hasOutput {
		attrs = append(attrs, GenAIOutputMessagesKey.String(messagesJSON(output, kept[len(input.messages):])))
	}
	if truncated {
		attrs = append(attrs, AnthropicContentTruncatedKey.Bool(true))
	}
	return attrs, size
}

// messagesSize returns the length of the JSON of the list of messages, whole.
func messagesSize(messages []messageContent) int {
	size := 2 + max(len(messages)-1, 0)
	for _, m := range messages {
		size += listSize(m.parts, m.prefix, messageSuffix)
	}
	return size
}

// messagesJSON returns the JSON of the list of messages, each with those of
// its parts whose JSON is in the element of parts at its index.
func messagesJSON(messages []messageContent, parts [][]string) string {
	size := 2 + max(len(messages)-1, 0)
	for i, m := range messages {
		size += listSize(parts[i], m.prefix, messageSuffix)
	}
	var value strings.Builder
	value.Grow(size)
	value.WriteString("[")
	for i, m := range messages {
		if i > 0 {
			value.WriteString(",")
		}
		value.WriteString(m.prefix)
		for k, j := range parts[i] {
			if k > 0 {
				value.WriteString(",")
			}
			value.WriteString(j)
		}
		value.WriteString(messageSuffix)
	}
	value.WriteString("]")
	return value.String()
}

// marshalParts returns the JSON of each of parts that can be encoded. A part
// that cannot be encoded is left out, so that the list stays JSON.
func marshalParts(parts []part) []string {
	jsons := make([]string, 0, len(parts))
	for _, p := range parts {
		if j := marshal(p); j != "" {
			jsons = append(jsons, j)
		}
	}
	return jsons
}

// listSize returns the length of the items jsons between prefix and suffix,
// with a comma between each two.
func listSize(jsons []string, prefix, suffix string) int {
	size := len(prefix) + len(suffix) + max(len(jsons)-1, 0)
	for _, j := range jsons {
		size += len(j)
	}
	return size
}

// cutText returns the JSON of a text part within limit bytes, whose content is
// the longest start of content that fits: the end is cut off, never inside the
// UTF-8 encoding of a character. The content is empty when none of it fits, and
// the JSON is then longer than a limit that has no room for the part.
func cutText(content string, limit int) string {
	cut := textPart{Type: "text"}
	room := limit - len(marshal(cut))
	// A byte of the content takes at least one byte of its JSON, so the bytes
	// that fit are among the first room. In their JSON each character is
	// itself, or an escape sequence, which starts with a backslash.
	head := content[:min(len(content), max(room, 0))]
	for len(head) > 0 && len(head) < len(content) && !utf8.RuneStart(content[len(head)]) {
		head = head[:len(head)-1]
	}
	escaped := marshal(head)
	length, escapedLength := 0, 0
	for length < len(head) {
		_, width := utf8.DecodeRuneInString(head[length:])
		step := width
		if escaped[1+escapedLength] == '\\' {
			step = 2
			if escaped[2+escapedLength] == 'u' {
				step = 6
			}
		}
		if escapedLength+step > room {
			break
		}
		escapedLength += step
		length += width
	}
	cut.Content = content[:length]
	return marshal(cut)
}

// truncateToolResultContent returns anthropic.tool_result.content within limit
// bytes, from blocks, the content blocks of a tool's result, dropping blocks from
// the end, the last one first, and reports whether it dropped any. The list is
// empty when not even the first block fits, as when the tool's arguments leave
// no room for it. No block is shortened.
func truncateToolResultContent(blocks []any, limit int) (value string, truncated bool) {
	// A block that cannot be encoded is left out, so that the list stays JSON.
	jsons := make([]string, 0, len(blocks))
	for _, block := range blocks {
		if j := marshal(block); j != "" {
			jsons = append(jsons, j)
		}
	}
	// The JSON of a list is its brackets, its items, and a comma between each two.
	size := 2 + max(len(jsons)-1, 0)
	for _, j := range jsons {
		size += len(j)
	}
	end := len(jsons)
	for end > 0 && size > limit {
		end--
		size -= len(jsons[end])
		// The comma before it, unless it was the first.
		if end > 0 {
			size--
		}
	}
	return "[" + strings.Join(jsons[:end], ",") + "]", end < len(jsons)
}
