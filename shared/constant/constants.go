// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package constant

import (
	"encoding/json"
)

type Constant[T any] interface {
	Default() T
}

// ValueOf gives the default value of a constant from its type. It's helpful when
// constructing constants as variants in a one-of. Note that empty structs are
// marshalled by default. Usage: constant.ValueOf[constant.Foo]()
func ValueOf[T Constant[T]]() T {
	var t T
	return t.Default()
}

type Any string                  // Always "any"
type APIError string             // Always "api_error"
type ApplicationPDF string       // Always "application/pdf"
type Assistant string            // Always "assistant"
type AuthenticationError string  // Always "authentication_error"
type Auto string                 // Always "auto"
type Base64 string               // Always "base64"
type Bash string                 // Always "bash"
type Bash20241022 string         // Always "bash_20241022"
type Bash20250124 string         // Always "bash_20250124"
type BillingError string         // Always "billing_error"
type Canceled string             // Always "canceled"
type CharLocation string         // Always "char_location"
type CitationsDelta string       // Always "citations_delta"
type Completion string           // Always "completion"
type Computer string             // Always "computer"
type Computer20241022 string     // Always "computer_20241022"
type Computer20250124 string     // Always "computer_20250124"
type Content string              // Always "content"
type ContentBlockDelta string    // Always "content_block_delta"
type ContentBlockLocation string // Always "content_block_location"
type ContentBlockStart string    // Always "content_block_start"
type ContentBlockStop string     // Always "content_block_stop"
type Disabled string             // Always "disabled"
type Document string             // Always "document"
type Enabled string              // Always "enabled"
type Ephemeral string            // Always "ephemeral"
type Error string                // Always "error"
type Errored string              // Always "errored"
type Expired string              // Always "expired"
type Image string                // Always "image"
type InputJSONDelta string       // Always "input_json_delta"
type InvalidRequestError string  // Always "invalid_request_error"
type Message string              // Always "message"
type MessageBatch string         // Always "message_batch"
type MessageBatchDeleted string  // Always "message_batch_deleted"
type MessageDelta string         // Always "message_delta"
type MessageStart string         // Always "message_start"
type MessageStop string          // Always "message_stop"
type Model string                // Always "model"
type None string                 // Always "none"
type NotFoundError string        // Always "not_found_error"
type Object string               // Always "object"
type OverloadedError string      // Always "overloaded_error"
type PageLocation string         // Always "page_location"
type PermissionError string      // Always "permission_error"
type RateLimitError string       // Always "rate_limit_error"
type RedactedThinking string     // Always "redacted_thinking"
type SignatureDelta string       // Always "signature_delta"
type StrReplaceEditor string     // Always "str_replace_editor"
type Succeeded string            // Always "succeeded"
type Text string                 // Always "text"
type TextDelta string            // Always "text_delta"
type TextEditor20241022 string   // Always "text_editor_20241022"
type TextEditor20250124 string   // Always "text_editor_20250124"
type TextPlain string            // Always "text/plain"
type Thinking string             // Always "thinking"
type ThinkingDelta string        // Always "thinking_delta"
type TimeoutError string         // Always "timeout_error"
type Tool string                 // Always "tool"
type ToolResult string           // Always "tool_result"
type ToolUse string              // Always "tool_use"
type URL string                  // Always "url"

func (c Any) Default() Any                                   { return "any" }
func (c APIError) Default() APIError                         { return "api_error" }
func (c ApplicationPDF) Default() ApplicationPDF             { return "application/pdf" }
func (c Assistant) Default() Assistant                       { return "assistant" }
func (c AuthenticationError) Default() AuthenticationError   { return "authentication_error" }
func (c Auto) Default() Auto                                 { return "auto" }
func (c Base64) Default() Base64                             { return "base64" }
func (c Bash) Default() Bash                                 { return "bash" }
func (c Bash20241022) Default() Bash20241022                 { return "bash_20241022" }
func (c Bash20250124) Default() Bash20250124                 { return "bash_20250124" }
func (c BillingError) Default() BillingError                 { return "billing_error" }
func (c Canceled) Default() Canceled                         { return "canceled" }
func (c CharLocation) Default() CharLocation                 { return "char_location" }
func (c CitationsDelta) Default() CitationsDelta             { return "citations_delta" }
func (c Completion) Default() Completion                     { return "completion" }
func (c Computer) Default() Computer                         { return "computer" }
func (c Computer20241022) Default() Computer20241022         { return "computer_20241022" }
func (c Computer20250124) Default() Computer20250124         { return "computer_20250124" }
func (c Content) Default() Content                           { return "content" }
func (c ContentBlockDelta) Default() ContentBlockDelta       { return "content_block_delta" }
func (c ContentBlockLocation) Default() ContentBlockLocation { return "content_block_location" }
func (c ContentBlockStart) Default() ContentBlockStart       { return "content_block_start" }
func (c ContentBlockStop) Default() ContentBlockStop         { return "content_block_stop" }
func (c Disabled) Default() Disabled                         { return "disabled" }
func (c Document) Default() Document                         { return "document" }
func (c Enabled) Default() Enabled                           { return "enabled" }
func (c Ephemeral) Default() Ephemeral                       { return "ephemeral" }
func (c Error) Default() Error                               { return "error" }
func (c Errored) Default() Errored                           { return "errored" }
func (c Expired) Default() Expired                           { return "expired" }
func (c Image) Default() Image                               { return "image" }
func (c InputJSONDelta) Default() InputJSONDelta             { return "input_json_delta" }
func (c InvalidRequestError) Default() InvalidRequestError   { return "invalid_request_error" }
func (c Message) Default() Message                           { return "message" }
func (c MessageBatch) Default() MessageBatch                 { return "message_batch" }
func (c MessageBatchDeleted) Default() MessageBatchDeleted   { return "message_batch_deleted" }
func (c MessageDelta) Default() MessageDelta                 { return "message_delta" }
func (c MessageStart) Default() MessageStart                 { return "message_start" }
func (c MessageStop) Default() MessageStop                   { return "message_stop" }
func (c Model) Default() Model                               { return "model" }
func (c None) Default() None                                 { return "none" }
func (c NotFoundError) Default() NotFoundError               { return "not_found_error" }
func (c Object) Default() Object                             { return "object" }
func (c OverloadedError) Default() OverloadedError           { return "overloaded_error" }
func (c PageLocation) Default() PageLocation                 { return "page_location" }
func (c PermissionError) Default() PermissionError           { return "permission_error" }
func (c RateLimitError) Default() RateLimitError             { return "rate_limit_error" }
func (c RedactedThinking) Default() RedactedThinking         { return "redacted_thinking" }
func (c SignatureDelta) Default() SignatureDelta             { return "signature_delta" }
func (c StrReplaceEditor) Default() StrReplaceEditor         { return "str_replace_editor" }
func (c Succeeded) Default() Succeeded                       { return "succeeded" }
func (c Text) Default() Text                                 { return "text" }
func (c TextDelta) Default() TextDelta                       { return "text_delta" }
func (c TextEditor20241022) Default() TextEditor20241022     { return "text_editor_20241022" }
func (c TextEditor20250124) Default() TextEditor20250124     { return "text_editor_20250124" }
func (c TextPlain) Default() TextPlain                       { return "text/plain" }
func (c Thinking) Default() Thinking                         { return "thinking" }
func (c ThinkingDelta) Default() ThinkingDelta               { return "thinking_delta" }
func (c TimeoutError) Default() TimeoutError                 { return "timeout_error" }
func (c Tool) Default() Tool                                 { return "tool" }
func (c ToolResult) Default() ToolResult                     { return "tool_result" }
func (c ToolUse) Default() ToolUse                           { return "tool_use" }
func (c URL) Default() URL                                   { return "url" }

func (c Any) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c APIError) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c ApplicationPDF) MarshalJSON() ([]byte, error)       { return marshalString(c) }
func (c Assistant) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c AuthenticationError) MarshalJSON() ([]byte, error)  { return marshalString(c) }
func (c Auto) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c Base64) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c Bash) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c Bash20241022) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c Bash20250124) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c BillingError) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c Canceled) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c CharLocation) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c CitationsDelta) MarshalJSON() ([]byte, error)       { return marshalString(c) }
func (c Completion) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c Computer) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c Computer20241022) MarshalJSON() ([]byte, error)     { return marshalString(c) }
func (c Computer20250124) MarshalJSON() ([]byte, error)     { return marshalString(c) }
func (c Content) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c ContentBlockDelta) MarshalJSON() ([]byte, error)    { return marshalString(c) }
func (c ContentBlockLocation) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c ContentBlockStart) MarshalJSON() ([]byte, error)    { return marshalString(c) }
func (c ContentBlockStop) MarshalJSON() ([]byte, error)     { return marshalString(c) }
func (c Disabled) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c Document) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c Enabled) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c Ephemeral) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c Error) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c Errored) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c Expired) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c Image) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c InputJSONDelta) MarshalJSON() ([]byte, error)       { return marshalString(c) }
func (c InvalidRequestError) MarshalJSON() ([]byte, error)  { return marshalString(c) }
func (c Message) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c MessageBatch) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c MessageBatchDeleted) MarshalJSON() ([]byte, error)  { return marshalString(c) }
func (c MessageDelta) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c MessageStart) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c MessageStop) MarshalJSON() ([]byte, error)          { return marshalString(c) }
func (c Model) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c None) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c NotFoundError) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c Object) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c OverloadedError) MarshalJSON() ([]byte, error)      { return marshalString(c) }
func (c PageLocation) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c PermissionError) MarshalJSON() ([]byte, error)      { return marshalString(c) }
func (c RateLimitError) MarshalJSON() ([]byte, error)       { return marshalString(c) }
func (c RedactedThinking) MarshalJSON() ([]byte, error)     { return marshalString(c) }
func (c SignatureDelta) MarshalJSON() ([]byte, error)       { return marshalString(c) }
func (c StrReplaceEditor) MarshalJSON() ([]byte, error)     { return marshalString(c) }
func (c Succeeded) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c Text) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c TextDelta) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c TextEditor20241022) MarshalJSON() ([]byte, error)   { return marshalString(c) }
func (c TextEditor20250124) MarshalJSON() ([]byte, error)   { return marshalString(c) }
func (c TextPlain) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c Thinking) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c ThinkingDelta) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c TimeoutError) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c Tool) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c ToolResult) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c ToolUse) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c URL) MarshalJSON() ([]byte, error)                  { return marshalString(c) }

type constant[T any] interface {
	Constant[T]
	*T
}

func marshalString[T ~string, PT constant[T]](v T) ([]byte, error) {
	var zero T
	if v == zero {
		v = PT(&v).Default()
	}
	return json.Marshal(string(v))
}
