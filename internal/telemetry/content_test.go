package telemetry

import (
	"bytes"
	"encoding/json"
	"math"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/tidwall/gjson"
	"go.opentelemetry.io/otel/attribute"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/tracing"
)

// contentLimit is the size limit of a call whose options set none.
const contentLimit = tracing.DefaultMaxContentBytes

// message is one message in the content attributes, as encoding/json writes it.
type message struct {
	Role  string `json:"role"`
	Parts []part `json:"parts"`
}

func textMessages(texts ...string) []message {
	messages := make([]message, len(texts))
	for i, text := range texts {
		messages[i] = message{Role: "user", Parts: []part{textPart{Type: "text", Content: text}}}
	}
	return messages
}

// textInput returns the content of a request whose messages are one user
// message for each of texts, with that text.
func textInput(texts ...string) inputContent {
	var input inputContent
	for _, text := range texts {
		input.messages = append(input.messages, newMessageContent("user", textParts(text)))
	}
	return input
}

// outputMessages returns gen_ai.output.messages of one message, which has parts.
func outputMessages(parts ...part) []messageContent {
	return []messageContent{newMessageContent("assistant", parts)}
}

// outputContent returns the serialized gen_ai.output.messages of one message,
// which has parts, whole.
func outputContent(parts ...part) string {
	messages := outputMessages(parts...)
	return messagesJSON(messages, [][]string{messages[0].parts})
}

// contentAttributes returns the values of gen_ai.input.messages and
// gen_ai.output.messages in attrs, and anthropic.content.truncated when attrs
// has it.
func contentAttributes(t *testing.T, attrs []attribute.KeyValue) (input, output string, truncated *attribute.Value) {
	t.Helper()
	for _, kv := range attrs {
		switch kv.Key {
		case GenAIInputMessagesKey:
			input = kv.Value.AsString()
		case GenAIOutputMessagesKey:
			output = kv.Value.AsString()
		case AnthropicContentTruncatedKey:
			value := kv.Value
			truncated = &value
		}
	}
	return input, output, truncated
}

func textParts(texts ...string) []part {
	parts := make([]part, len(texts))
	for i, text := range texts {
		parts[i] = textPart{Type: "text", Content: text}
	}
	return parts
}

// wantTruncated fails the test unless truncated is true when want, and absent otherwise.
func wantTruncated(t *testing.T, truncated *attribute.Value, want bool) {
	t.Helper()
	switch {
	case !want && truncated != nil:
		t.Errorf("anthropic.content.truncated = %v, want it absent", truncated.Emit())
	case want && (truncated == nil || truncated.Type() != attribute.BOOL || !truncated.AsBool()):
		t.Errorf("anthropic.content.truncated = %v, want true", truncated)
	}
}

// wantSize fails the test unless size is the length of the values of the
// content attributes in attrs.
func wantSize(t *testing.T, attrs []attribute.KeyValue, size int) {
	t.Helper()
	total := 0
	for _, kv := range attrs {
		switch kv.Key {
		case GenAIRequestStopSequencesKey:
			for _, sequence := range kv.Value.AsStringSlice() {
				total += len(sequence)
			}
		case GenAISystemInstructionsKey, GenAIInputMessagesKey, GenAIOutputMessagesKey:
			total += len(kv.Value.AsString())
		}
	}
	if size != total {
		t.Errorf("the reported size is %d, want %d, the length of the values", size, total)
	}
}

// sizeBoundContent is the content attributes of a case of
// testdata/content_size_bound.json. A value that is null is an absent attribute.
type sizeBoundContent struct {
	StopSequences      []string        `json:"stop_sequences"`
	SystemInstructions json.RawMessage `json:"system_instructions"`
	InputMessages      json.RawMessage `json:"input_messages"`
	OutputMessages     json.RawMessage `json:"output_messages"`
}

// sizeBoundParts returns the JSON of each message part in raws, as the SDK
// serializes that part.
func sizeBoundParts(t *testing.T, raws []json.RawMessage) []string {
	t.Helper()
	parts := make([]part, len(raws))
	for i, raw := range raws {
		var fields struct {
			Type, Content, ID, Name string
			Arguments               json.RawMessage
		}
		if err := json.Unmarshal(raw, &fields); err != nil {
			t.Fatal(err)
		}
		switch fields.Type {
		case "text":
			parts[i] = textPart{Type: "text", Content: fields.Content}
		case "reasoning":
			parts[i] = reasoningPart{Type: "reasoning", Content: fields.Content}
		case "tool_call":
			parts[i] = toolCallRequestPart{Type: "tool_call", ID: fields.ID, Name: fields.Name, Arguments: fields.Arguments}
		default:
			parts[i] = genericPart{Type: fields.Type}
		}
		if got := marshal(parts[i]); got != compactJSON(t, raw) {
			t.Fatalf("the part %s is serialized as %s", raw, got)
		}
	}
	return marshalParts(parts)
}

// sizeBoundMessages returns the messages whose JSON list is raw.
func sizeBoundMessages(t *testing.T, raw json.RawMessage) []messageContent {
	t.Helper()
	var decoded []struct {
		Role  string
		Parts []json.RawMessage
	}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	messages := make([]messageContent, len(decoded))
	for i, m := range decoded {
		messages[i] = newMessageContent(m.Role, nil)
		messages[i].parts = sizeBoundParts(t, m.Parts)
	}
	return messages
}

// present reports whether raw is a JSON value other than null.
func present(raw json.RawMessage) bool { return len(raw) > 0 && string(raw) != "null" }

// compactJSON returns raw without insignificant space, and "" for an absent value.
func compactJSON(t *testing.T, raw json.RawMessage) string {
	t.Helper()
	if !present(raw) {
		return ""
	}
	var compact bytes.Buffer
	if err := json.Compact(&compact, raw); err != nil {
		t.Fatal(err)
	}
	return compact.String()
}

// attributeValues returns the value of each content attribute in attrs that is
// a string, by its key.
func attributeValues(attrs []attribute.KeyValue) map[attribute.Key]string {
	values := map[attribute.Key]string{}
	for _, kv := range attrs {
		if kv.Value.Type() == attribute.STRING {
			values[kv.Key] = kv.Value.AsString()
		}
	}
	return values
}

// The cases of testdata/content_size_bound.json are the same for every SDK.
// Their sizes are the bytes of JSON without insignificant space, in which no
// character other than a quote, a backslash or a control character is escaped.
func TestContentSizeBoundCases(t *testing.T) {
	data, err := os.ReadFile("testdata/content_size_bound.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Name       string           `json:"name"`
		Bound      int              `json:"bound"`
		SizeBefore int              `json:"size_before"`
		Input      sizeBoundContent `json:"input"`
		Expected   sizeBoundContent `json:"expected"`
		Truncated  bool             `json:"expected_truncated"`
		SizeAfter  int              `json:"size_after"`
	}
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) == 0 {
		t.Fatal("no cases")
	}
	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			input := inputContent{stopSequences: test.Input.StopSequences}
			if present(test.Input.SystemInstructions) {
				var raws []json.RawMessage
				if err := json.Unmarshal(test.Input.SystemInstructions, &raws); err != nil {
					t.Fatal(err)
				}
				input.system = sizeBoundParts(t, raws)
				input.hasSystem = true
			}
			if present(test.Input.InputMessages) {
				input.messages = sizeBoundMessages(t, test.Input.InputMessages)
			}
			var output []messageContent
			hasOutput := present(test.Input.OutputMessages)
			if hasOutput {
				output = sizeBoundMessages(t, test.Input.OutputMessages)
			}

			// A request always has gen_ai.input.messages, so a case that has none
			// stands for an empty list there, whose two bytes its sizes leave out.
			empty := 0
			if !present(test.Input.InputMessages) {
				empty = len("[]")
			}

			whole, size := input.attributes(output, hasOutput, math.MaxInt)
			if size != test.SizeBefore+empty {
				t.Errorf("the whole content has %d bytes, want %d", size, test.SizeBefore+empty)
			}
			wantSize(t, whole, size)

			attrs, size := input.attributes(output, hasOutput, test.Bound+empty)
			if size != test.SizeAfter+empty {
				t.Errorf("the content has %d bytes, want %d", size, test.SizeAfter+empty)
			}
			wantSize(t, attrs, size)
			values := attributeValues(attrs)
			for key, want := range map[attribute.Key]json.RawMessage{
				GenAISystemInstructionsKey: test.Expected.SystemInstructions,
				GenAIInputMessagesKey:      test.Expected.InputMessages,
				GenAIOutputMessagesKey:     test.Expected.OutputMessages,
			} {
				if key == GenAIInputMessagesKey && !present(want) {
					want = json.RawMessage("[]")
				}
				got, ok := values[key]
				if ok != present(want) || got != compactJSON(t, want) {
					t.Errorf("%s = %q (present %t), want %q", key, got, ok, compactJSON(t, want))
				}
			}
			_, _, truncated := contentAttributes(t, attrs)
			wantTruncated(t, truncated, test.Truncated)
			wantStopSequences(t, whole, test.Input.StopSequences)
			wantStopSequences(t, attrs, test.Expected.StopSequences)
		})
	}
}

// What the last text part becomes when not even its empty content fits is the
// one case that the rule leaves open: here it stays, with no content.
func TestContentSizeBoundKeepsTheLastTextPartWithEmptyContentWhenNoneOfItFits(t *testing.T) {
	input := textInput("question")
	// The stop sequences fit, so they stay, and leave no room for the message.
	input.stopSequences = []string{strings.Repeat("x", 95)}
	attrs, size := input.attributes(nil, false, 100)
	got, _, truncated := contentAttributes(t, attrs)
	if want := marshal(textMessages("")); got != want {
		t.Errorf("gen_ai.input.messages = %s, want %s: the part stays, with empty content", got, want)
	}
	if want := 95 + len(marshal(textMessages(""))); size != want {
		t.Errorf("the values have %d bytes, want %d, over the limit", size, want)
	}
	wantTruncated(t, truncated, true)
	wantStopSequences(t, attrs, input.stopSequences)

	// The same beside output: the text is the last part of the list, and it stays.
	attrs, _ = input.attributes(outputMessages(genericPart{Type: "anthropic_image"}, textPart{Type: "text", Content: "answer"}), true, 100)
	got, output, truncated := contentAttributes(t, attrs)
	if want := `[{"role":"user","parts":[]}]`; got != want {
		t.Errorf("gen_ai.input.messages = %s, want %s", got, want)
	}
	if want := outputContent(textParts("")...); output != want {
		t.Errorf("gen_ai.output.messages = %s, want %s: the part stays, with empty content", output, want)
	}
	wantTruncated(t, truncated, true)

	// A text part that has no content to cut is not a reason to say that anything was.
	attrs, _ = textInput("").attributes(nil, false, 10)
	got, _, truncated = contentAttributes(t, attrs)
	if want := marshal(textMessages("")); got != want {
		t.Errorf("gen_ai.input.messages = %s, want %s", got, want)
	}
	wantTruncated(t, truncated, false)
}

func TestInputContentAttributesUnderTheLimitAreTheMarshaledValues(t *testing.T) {
	params := gjson.Parse(`{"system":[{"type":"text","text":"Be <brief> & kind."},{"type":"text","text":"Really."}],"messages":[
		{"role":"user","content":[{"type":"text","text":"Weather?"},{"type":"image","source":{"type":"base64","data":"AAAA"}}]},
		{"role":"assistant","content":[]},
		{"role":"user","content":"Tomorrow?"}]}`)
	output := []part{reasoningPart{Type: "reasoning", Content: "hm"}, textPart{Type: "text", Content: "Sunny."}}
	attrs, size := newInputContent(params.Get("system"), params.Get("messages"), nil).attributes(outputMessages(output...), true, contentLimit)
	values := attributeValues(attrs)
	for key, want := range map[attribute.Key]any{
		GenAISystemInstructionsKey: textParts("Be <brief> & kind.", "Really."),
		GenAIInputMessagesKey: []message{
			{Role: "user", Parts: []part{textPart{Type: "text", Content: "Weather?"}, genericPart{Type: "anthropic_image"}}},
			{Role: "assistant", Parts: []part{}},
			{Role: "user", Parts: textParts("Tomorrow?")},
		},
		GenAIOutputMessagesKey: []message{{Role: "assistant", Parts: output}},
	} {
		if values[key] != marshal(want) {
			t.Errorf("%s = %s, want %s", key, values[key], marshal(want))
		}
	}
	wantSize(t, attrs, size)
	_, _, truncated := contentAttributes(t, attrs)
	wantTruncated(t, truncated, false)
}

func TestInputContentAttributesRemovesTheFirstPartsThatDoNotFit(t *testing.T) {
	big := strings.Repeat("a", 225_000)
	params := gjson.Parse(marshal(map[string]any{
		"system":   strings.Repeat("s", 150_000),
		"messages": []map[string]any{{"role": "user", "content": "0" + big}, {"role": "user", "content": "1" + big}, {"role": "user", "content": "2" + big}},
	}))
	// Three messages of 225 kB and 50 kB of output fit on their own, but not
	// with 150 kB of system instructions, which go first.
	output := outputMessages(textParts(strings.Repeat("o", 50_000))...)
	attrs, size := newInputContent(params.Get("system"), params.Get("messages"), nil).attributes(output, true, contentLimit)
	input, got, truncated := contentAttributes(t, attrs)
	if system := systemInstructions(attrs); system != "[]" {
		t.Errorf("gen_ai.system_instructions have %d bytes, want an empty list", len(system))
	}
	if input != marshal(textMessages("0"+big, "1"+big, "2"+big)) || got != messagesJSON(output, [][]string{output[0].parts}) {
		t.Errorf("kept %d bytes of input and %d of output, want both whole", len(input), len(got))
	}
	wantSize(t, attrs, size)
	wantTruncated(t, truncated, true)

	// With 250 kB of output, the only part of the first message goes too, and the
	// message stays.
	output = outputMessages(textParts(strings.Repeat("o", 250_000))...)
	attrs, size = newInputContent(params.Get("system"), params.Get("messages"), nil).attributes(output, true, contentLimit)
	input, got, truncated = contentAttributes(t, attrs)
	want := []message{{Role: "user", Parts: []part{}}, {Role: "user", Parts: textParts("1" + big)}, {Role: "user", Parts: textParts("2" + big)}}
	if systemInstructions(attrs) != "[]" || input != marshal(want) || got != messagesJSON(output, [][]string{output[0].parts}) {
		t.Errorf("kept %d bytes of input and %d of output, want the first message without its part and the rest whole", len(input), len(got))
	}
	if size > contentLimit {
		t.Errorf("the values have %d bytes, more than the limit", size)
	}
	wantTruncated(t, truncated, true)

	// Without the system instructions, the messages and the first output fit.
	attrs, _ = newInputContent(gjson.Result{}, params.Get("messages"), nil).attributes(outputMessages(textParts(strings.Repeat("o", 50_000))...), true, contentLimit)
	input, _, truncated = contentAttributes(t, attrs)
	if input != marshal(textMessages("0"+big, "1"+big, "2"+big)) {
		t.Errorf("kept %d bytes of input without the system instructions, want all three messages", len(input))
	}
	wantTruncated(t, truncated, false)
}

func TestInputContentAttributesFillTheLimitExactly(t *testing.T) {
	input := textInput("old", strings.Repeat("b", 1000), strings.Repeat("c", 1000))
	// No room is set aside for a marker: only the first part goes when the rest fit to the byte.
	want := []message{{Role: "user", Parts: []part{}}, {Role: "user", Parts: textParts(strings.Repeat("b", 1000))}, {Role: "user", Parts: textParts(strings.Repeat("c", 1000))}}
	limit := len(marshal(want))
	attrs, size := input.attributes(nil, false, limit)
	got, _, truncated := contentAttributes(t, attrs)
	if got != marshal(want) || size != limit {
		t.Errorf("kept %d bytes, want the newest two messages whole in %d", size, limit)
	}
	wantTruncated(t, truncated, true)

	want[1].Parts = []part{}
	attrs, size = input.attributes(nil, false, limit-1)
	got, _, truncated = contentAttributes(t, attrs)
	if got != marshal(want) {
		t.Errorf("one byte less kept %d bytes, want only the newest message whole", size)
	}
	wantSize(t, attrs, size)
	wantTruncated(t, truncated, true)

	// The content that was fitted is whole the next time.
	attrs, _ = input.attributes(nil, false, contentLimit)
	got, _, truncated = contentAttributes(t, attrs)
	if got != marshal(textMessages("old", strings.Repeat("b", 1000), strings.Repeat("c", 1000))) {
		t.Errorf("gen_ai.input.messages = %.80s, want every message whole", got)
	}
	wantTruncated(t, truncated, false)
}

func TestInputContentAttributesCountsTheStopSequencesAndRemovesThemWholeWhenTheyAloneAreOverTheLimit(t *testing.T) {
	big := strings.Repeat("a", 225_000)
	messages := gjson.Parse(marshal([]map[string]any{{"role": "user", "content": "0" + big}, {"role": "user", "content": "1" + big}, {"role": "user", "content": "2" + big}}))
	// Three messages of 225 kB fit on their own, but not with 150 kB of stop sequences.
	stop := []string{strings.Repeat("x", 90_000), strings.Repeat("y", 60_000)}
	attrs, size := newInputContent(gjson.Result{}, messages, stop).attributes(nil, false, contentLimit)
	input, _, truncated := contentAttributes(t, attrs)
	if want := marshal([]message{{Role: "user", Parts: []part{}}, {Role: "user", Parts: textParts("1" + big)}, {Role: "user", Parts: textParts("2" + big)}}); input != want {
		t.Errorf("gen_ai.input.messages have %d bytes, want the first message without its part and the newest two whole", len(input))
	}
	if want := 150_000 + len(input); size != want {
		t.Errorf("the values have %d bytes, want %d: the stop sequences and the input messages", size, want)
	}
	wantTruncated(t, truncated, true)
	wantStopSequences(t, attrs, stop)

	// Stop sequences that fill the limit to the byte stay whole, and every part
	// goes but the last, a text part, which stays without content.
	stop = []string{strings.Repeat("x", contentLimit-len("END")), "END"}
	output := outputMessages(textParts("one", "two")...)
	attrs, size = newInputContent(gjson.Result{}, messages, stop).attributes(output, true, contentLimit)
	input, got, truncated := contentAttributes(t, attrs)
	if input != marshal([]message{{Role: "user", Parts: []part{}}, {Role: "user", Parts: []part{}}, {Role: "user", Parts: []part{}}}) || got != outputContent(textParts("")...) {
		t.Errorf("gen_ai.input.messages = %.120s and gen_ai.output.messages = %.80s, want three messages with no parts and one with an empty text part", input, got)
	}
	if size <= contentLimit {
		t.Errorf("the values have %d bytes, want them over the limit", size)
	}
	wantSize(t, attrs, size)
	wantTruncated(t, truncated, true)
	wantStopSequences(t, attrs, stop)

	// One byte more and they are over the limit on their own: the attribute is
	// left out, and every part has room.
	stop[0] += "x"
	attrs, size = newInputContent(gjson.Result{}, messages, stop).attributes(output, true, contentLimit)
	input, got, truncated = contentAttributes(t, attrs)
	if input != marshal(textMessages("0"+big, "1"+big, "2"+big)) || got != outputContent(textParts("one", "two")...) {
		t.Errorf("kept %d bytes of input and %d of output, want both whole", len(input), len(got))
	}
	if want := len(input) + len(got); size != want {
		t.Errorf("the values have %d bytes, want %d: the stop sequences count as none", size, want)
	}
	wantTruncated(t, truncated, true)
	wantStopSequences(t, attrs, nil)

	// The same when the request's content is set before the output is known.
	attrs, _ = newInputContent(gjson.Result{}, messages, stop).attributes(nil, false, contentLimit)
	_, _, truncated = contentAttributes(t, attrs)
	wantTruncated(t, truncated, true)
	wantStopSequences(t, attrs, nil)

	// A request without stop sequences has no such attribute, and neither has one
	// whose stop sequences are not all valid UTF-8, however long: those are not
	// left out for their size, so nothing says that content was removed. The other
	// attributes stay.
	for _, none := range [][]string{nil, {}, {"END", "\xff\xfe"}, {strings.Repeat("x", contentLimit), "\xff\xfe"}} {
		attrs, _ = newInputContent(gjson.Result{}, messages, none).attributes(nil, false, contentLimit)
		input, _, truncated = contentAttributes(t, attrs)
		if input != marshal(textMessages("0"+big, "1"+big, "2"+big)) {
			t.Errorf("kept %d bytes of input, want every message whole", len(input))
		}
		wantTruncated(t, truncated, false)
		wantStopSequences(t, attrs, nil)
	}
}

// wantStopSequences fails the test unless attrs has gen_ai.request.stop_sequences
// as the list of strings want, or lacks it when want is nil.
func wantStopSequences(t *testing.T, attrs []attribute.KeyValue, want []string) {
	t.Helper()
	for _, kv := range attrs {
		if kv.Key != GenAIRequestStopSequencesKey {
			continue
		}
		if want == nil {
			t.Errorf("gen_ai.request.stop_sequences has %d elements, want it absent", len(kv.Value.AsStringSlice()))
		} else if kv.Value.Type() != attribute.STRINGSLICE || !slices.Equal(kv.Value.AsStringSlice(), want) {
			t.Errorf("gen_ai.request.stop_sequences is a %v of %d elements, want the %d stop sequences as a list of strings", kv.Value.Type(), len(kv.Value.AsStringSlice()), len(want))
		}
		return
	}
	if want != nil {
		t.Error("gen_ai.request.stop_sequences is absent")
	}
}

func TestInputContentAttributesWithoutOutputHasNoOutputMessages(t *testing.T) {
	attrs, _ := newInputContent(gjson.Result{}, gjson.Parse(`[{"role":"user","content":"hi"}]`), nil).attributes(nil, false, contentLimit)
	for _, kv := range attrs {
		if kv.Key == GenAIOutputMessagesKey {
			t.Errorf("gen_ai.output.messages = %s, want it absent", kv.Value.AsString())
		}
	}
}

// systemInstructions returns the value of gen_ai.system_instructions in attrs.
func systemInstructions(attrs []attribute.KeyValue) string {
	for _, kv := range attrs {
		if kv.Key == GenAISystemInstructionsKey {
			return kv.Value.AsString()
		}
	}
	return ""
}

func TestInputContentAttributesFitsEveryPartOfTheRequestAgainBesideTheOutput(t *testing.T) {
	text := strings.Repeat("h", 200)
	input := textInput(text)
	limit := len(marshal(textMessages(text))) - 50

	// When the call starts, the text is the last part, so its end goes.
	attrs, size := input.attributes(nil, false, limit)
	got, _, truncated := contentAttributes(t, attrs)
	if want := marshal(textMessages(text[:150])); got != want || size != limit {
		t.Errorf("gen_ai.input.messages = %s in %d bytes, want the first 150 bytes of the text in %d", got, size, limit)
	}
	wantTruncated(t, truncated, true)

	// Beside output it is no longer the last part. It goes whole, as it would if
	// its end had not gone before, and the output has the room that it leaves.
	output := outputMessages(textParts(strings.Repeat("o", 100))...)
	attrs, size = input.attributes(output, true, limit)
	got, out, truncated := contentAttributes(t, attrs)
	if got != `[{"role":"user","parts":[]}]` || out != messagesJSON(output, [][]string{output[0].parts}) {
		t.Errorf("gen_ai.input.messages = %s and gen_ai.output.messages = %s, want a message with no parts and the output whole", got, out)
	}
	wantSize(t, attrs, size)
	wantTruncated(t, truncated, true)
}

func TestInputContentAttributesRemovesPartsOfSeveralOutputMessagesFromTheFirst(t *testing.T) {
	first, second := textParts("first "+strings.Repeat("f", 100), "second"), textParts("third", "fourth")
	output := []messageContent{newMessageContent("assistant", first), newMessageContent("assistant", second)}
	whole := marshal([]message{{Role: "assistant", Parts: first}, {Role: "assistant", Parts: second}})
	input := textInput("question")
	attrs, _ := input.attributes(output, true, contentLimit)
	if _, got, _ := contentAttributes(t, attrs); got != whole {
		t.Errorf("gen_ai.output.messages = %s, want %s", got, whole)
	}

	// The question, and the first part of the first output message.
	attrs, size := input.attributes(output, true, len(`[{"role":"user","parts":[]}]`)+len(whole)-100)
	got, out, truncated := contentAttributes(t, attrs)
	want := marshal([]message{{Role: "assistant", Parts: textParts("second")}, {Role: "assistant", Parts: second}})
	if got != `[{"role":"user","parts":[]}]` || out != want {
		t.Errorf("gen_ai.input.messages = %s and gen_ai.output.messages = %s, want %s", got, out, want)
	}
	wantSize(t, attrs, size)
	wantTruncated(t, truncated, true)

	// Every part but the last, whose end goes. The first message stays, with no parts.
	shells := len(`[{"role":"user","parts":[]}]`) + len(marshal([]message{{Role: "assistant", Parts: []part{}}, {Role: "assistant", Parts: textParts("")}}))
	attrs, size = input.attributes(output, true, shells+3)
	_, out, truncated = contentAttributes(t, attrs)
	want = marshal([]message{{Role: "assistant", Parts: []part{}}, {Role: "assistant", Parts: textParts("fou")}})
	if out != want || size != shells+3 {
		t.Errorf("gen_ai.output.messages = %s in %d bytes, want %s in %d", out, size, want, shells+3)
	}
	wantTruncated(t, truncated, true)
}

func TestCutTextCutsBetweenCharacters(t *testing.T) {
	// Characters of three and four bytes, and then some that encoding/json
	// escapes: in two bytes, in six, and a byte that is not UTF-8, in six.
	text := strings.Repeat("\u4e16\U0001f600\n\x01<\u2028\xff", 3)
	part := func(content string) string { return marshal(textPart{Type: "text", Content: content}) }
	for limit := len(part("")) - 5; limit <= len(part(text)); limit++ {
		// The most whole characters whose JSON fits, or none.
		want := part("")
		for end := range text {
			if value := part(text[:end]); len(value) <= limit {
				want = value
			}
		}
		if value := part(text); len(value) <= limit {
			want = value
		}
		if got := cutText(text, limit); got != want {
			t.Errorf("limit %d: kept %d bytes, want %d: %s", limit, len(got), len(want), want)
		}
	}
}

func TestInputContentAttributesRemovesAPartThatIsNotTextWhole(t *testing.T) {
	call := toolCallRequestPart{Type: "tool_call", ID: "toolu_1", Name: "f", Arguments: json.RawMessage(`{"content":"yyyyy"}`)}
	result := toolCallResponsePart{Type: "tool_call_response", ID: "toolu_1", Response: []any{textBlock{Type: "text", Text: "one"}, textBlock{Type: "text", Text: "two"}}}
	reasoning := reasoningPart{Type: "reasoning", Content: "zzzzz"}
	image := genericPart{Type: "anthropic_image"}
	text := textPart{Type: "text", Content: "xxxxx"}
	size := func(parts ...part) int { return len("[]") + len(outputContent(parts...)) }
	for name, test := range map[string]struct {
		parts []part
		limit int
		want  []part
	}{
		"a part with only a type goes whole":             {[]part{image, text}, size(image, text) - 1, []part{text}},
		"reasoning is not cut as text is":                {[]part{text, reasoning}, size(text, reasoning) - 1, []part{reasoning}},
		"reasoning that is the last part goes whole":     {[]part{text, reasoning}, size(reasoning) - 1, []part{}},
		"a tool call goes whole":                         {[]part{call, text}, size(call, text) - 1, []part{text}},
		"a tool call that is the last part goes whole":   {[]part{text, call}, size(call) - 1, []part{}},
		"a tool result goes whole, not block by block":   {[]part{result, text}, size(result, text) - 1, []part{text}},
		"a tool result that is the last part goes whole": {[]part{result}, size(result) - 1, []part{}},
		"text that is not the last part is not cut":      {[]part{text, image}, size(text, image) - 1, []part{image}},
		"text that is the last part is cut from its end": {[]part{call, text}, size(text) - 1, []part{textPart{Type: "text", Content: "xxxx"}}},
		"nothing of the last text is left, and it stays": {[]part{call, text}, size() - 100, []part{textPart{Type: "text"}}},
		"the last part that is not text leaves no part":  {[]part{text, image}, size() - 100, []part{}},
	} {
		attrs, reported := inputContent{}.attributes(outputMessages(test.parts...), true, test.limit)
		_, got, truncated := contentAttributes(t, attrs)
		if got != outputContent(test.want...) {
			t.Errorf("%s: got %s, want %s", name, got, outputContent(test.want...))
		}
		wantSize(t, attrs, reported)
		wantTruncated(t, truncated, true)
	}
}

func TestContentThatCannotBeEncodedIsLeftOutOfItsList(t *testing.T) {
	broken := toolCallRequestPart{Type: "tool_call", ID: "toolu_1", Name: "get_weather", Arguments: json.RawMessage("{")}
	// The input of the second message's tool call is not JSON.
	messages := gjson.Parse(`[{"role":"user","content":"first"},
		{"role":"assistant","content":[{"type":"tool_use","id":"toolu_1","name":"get_weather","input":{"city":}}]},
		{"role":"user","content":"third"}]`)
	attrs, _ := newInputContent(gjson.Result{}, messages, nil).attributes(nil, false, contentLimit)
	if got, _, _ := contentAttributes(t, attrs); got != marshal(textMessages("first", "third")) {
		t.Errorf("input messages = %s, want the two that encode, as a JSON list", got)
	}
	parts := []part{textPart{Type: "text", Content: "one"}, broken, textPart{Type: "text", Content: "three"}}
	if got := outputContent(parts...); got != outputContent(textParts("one", "three")...) {
		t.Errorf("output messages = %s, want the two parts that encode, as JSON", got)
	}
	if got := toolCallArguments(gjson.Parse(`{"type":"tool_use"}`).Get("input")); string(got) != "{}" {
		t.Errorf("the arguments of a tool call without input = %s, want an empty object", got)
	}
}

func TestToolResultContentIsMappedWithoutRecursion(t *testing.T) {
	// The API allows no tool_result inside a tool_result. One that is there all
	// the same is a block with its type only, like an image.
	block := gjson.Parse(`{"type":"tool_result","tool_use_id":"toolu_1","content":[
		{"type":"text","text":"sunny"},
		{"type":"image","source":{"type":"base64","data":"AAAA"}},
		{"type":"tool_result","tool_use_id":"toolu_2","content":[{"type":"text","text":"nested"}]},
		{"type":"tool_use","id":"toolu_3","name":"get_weather","input":{"city":"Paris"}}]}`)
	want := `{"type":"tool_call_response","id":"toolu_1","response":[{"text":"sunny","type":"text"},{"type":"image"},{"type":"tool_result"},{"type":"tool_use"}]}`
	if got := marshal(blockPart(block)); got != want {
		t.Errorf("part = %s, want %s", got, want)
	}

	// However deep the nesting, the part tree is two levels deep.
	const depth = 100_000
	nested := strings.Repeat(`{"type":"tool_result","tool_use_id":"t","content":[`, depth) + strings.Repeat(`]}`, depth)
	if got, want := marshal(blockPart(gjson.Parse(nested))), `{"type":"tool_call_response","id":"t","response":[{"type":"tool_result"}]}`; got != want {
		t.Errorf("part of %d nested tool results has %d bytes, want %s", depth, len(got), want)
	}
}

func TestInputContentAttributesOfNoMessagesAreAnEmptyList(t *testing.T) {
	for _, limit := range []int{contentLimit, 0, -100} {
		attrs, size := inputContent{}.attributes(nil, false, limit)
		if got, _, truncated := contentAttributes(t, attrs); got != "[]" || truncated != nil || size != len("[]") {
			t.Errorf("limit %d: got %s in %d bytes (truncated %v), want [] and nothing removed", limit, got, size, truncated)
		}
	}
}

func TestToolResultContentIsAListOfContentBlocks(t *testing.T) {
	for name, test := range map[string]struct{ result, want string }{
		"a string is one text block": {`{"content":"sunny"}`, `[{"type":"text","text":"sunny"}]`},
		"a text block and a search_result block keep their fields but the encrypted ones, and any other block holds its type only": {
			`{"content":[{"type":"text","text":"sunny","citations":[{"type":"web_search_result_location","url":"https://example.com","encrypted_index":"opaque"}],"cache_control":{"type":"ephemeral"}},
				{"type":"search_result","source":"https://example.com","title":"Forecast","content":[{"type":"text","text":"rain","citations":[{"type":"web_search_result_location","encrypted_index":"opaque"}]}]},
				{"type":"image","source":{"type":"base64","data":"AAAA"}}]}`,
			`[{"cache_control":{"type":"ephemeral"},"citations":[{"type":"web_search_result_location","url":"https://example.com"}],"text":"sunny","type":"text"},` +
				`{"content":[{"citations":[{"type":"web_search_result_location"}],"text":"rain","type":"text"}],"source":"https://example.com","title":"Forecast","type":"search_result"},{"type":"image"}]`,
		},
		"an empty list":          {`{"content":[]}`, `[]`},
		"no content is no block": {`{}`, `[]`},
		"null is no block":       {`{"content":null}`, `[]`},
	} {
		t.Run(name, func(t *testing.T) {
			if got := marshal(toolResultContent(gjson.Parse(test.result).Get("content"))); got != test.want {
				t.Errorf("content = %s, want %s", got, test.want)
			}
		})
	}
}

func TestResultBlockContent(t *testing.T) {
	for name, test := range map[string]struct{ block, want string }{
		"a string is one text block": {`{"content":"sunny"}`, `[{"type":"text","text":"sunny"}]`},
		"in a list, text, search results and web search results are kept without their encrypted fields, and any other element holds its type only": {
			`{"content":[{"type":"text","text":"rain","citations":null},
				{"type":"web_search_result","url":"https://example.com","title":"Forecast","encrypted_content":"opaque","page_age":null},
				{"type":"search_result","source":"https://example.com","title":"Forecast","content":[{"type":"text","text":"rain"}]},
				{"type":"image","source":{"type":"base64","data":"AAAA"}},
				{"data":"no type"}]}`,
			`[{"citations":null,"text":"rain","type":"text"},{"page_age":null,"title":"Forecast","type":"web_search_result","url":"https://example.com"},` +
				`{"content":[{"text":"rain","type":"text"}],"source":"https://example.com","title":"Forecast","type":"search_result"},{"type":"image"},{"type":"unknown"}]`,
		},
		"one object is kept": {
			`{"content":{"type":"code_execution_result","stdout":"1\n","stderr":"","return_code":0,"content":[{"type":"code_execution_output","file_id":"file_1"}]}}`,
			`{"content":[{"file_id":"file_1","type":"code_execution_output"}],"return_code":0,"stderr":"","stdout":"1\n","type":"code_execution_result"}`,
		},
		"an error object is kept": {
			`{"content":{"type":"web_search_tool_result_error","error_code":"max_uses_exceeded"}}`,
			`{"error_code":"max_uses_exceeded","type":"web_search_tool_result_error"}`,
		},
		"an encrypted field is left out at any depth": {
			`{"content":{"type":"encrypted_code_execution_result","encrypted_stdout":"opaque","return_code":0,"content":[{"type":"code_execution_output","encrypted_name":"opaque","file_id":"file_1"}]}}`,
			`{"content":[{"file_id":"file_1","type":"code_execution_output"}],"return_code":0,"type":"encrypted_code_execution_result"}`,
		},
		"the content of a viewed image is left out": {
			`{"content":{"type":"text_editor_code_execution_view_result","file_type":"image","content":"AAAA","num_lines":null}}`,
			`{"file_type":"image","num_lines":null,"type":"text_editor_code_execution_view_result"}`,
		},
		"the content of a viewed text file is kept": {
			`{"content":{"type":"text_editor_code_execution_view_result","file_type":"text","content":"hello","num_lines":1}}`,
			`{"content":"hello","file_type":"text","num_lines":1,"type":"text_editor_code_execution_view_result"}`,
		},
		"a document inside an object holds its type only": {
			`{"content":{"type":"web_fetch_result","url":"https://example.com","retrieved_at":null,"content":{"type":"document","title":"Forecast","source":{"type":"text","data":"rain"}}}}`,
			`{"content":{"type":"document"},"retrieved_at":null,"type":"web_fetch_result","url":"https://example.com"}`,
		},
		"no content is an empty list": {`{}`, `[]`},
		"null is an empty list":       {`{"content":null}`, `[]`},
	} {
		t.Run(name, func(t *testing.T) {
			if got := marshal(resultBlockContent(gjson.Parse(test.block).Get("content"))); got != test.want {
				t.Errorf("response = %s, want %s", got, test.want)
			}
		})
	}
}

func TestBlockWithoutATypeIsUnknown(t *testing.T) {
	if got, want := marshal(blockPart(gjson.Parse(`{"text":"sunny"}`))), `{"type":"anthropic_unknown"}`; got != want {
		t.Errorf("part = %s, want %s", got, want)
	}
	if got, want := marshal(toolResultContent(gjson.Parse(`[{"type":7}]`))), `[{"type":"unknown"}]`; got != want {
		t.Errorf("content = %s, want %s", got, want)
	}
}
