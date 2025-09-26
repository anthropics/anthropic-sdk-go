// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package constant

import (
	shimjson "github.com/anthropics/anthropic-sdk-go/internal/encoding/json"
)

// ModelNonStreamingTokens defines the maximum tokens for models that should limit
// non-streaming requests.
var ModelNonStreamingTokens = map[string]int{
	"claude-opus-4-20250514":                  8192,
	"claude-4-opus-20250514":                  8192,
	"claude-opus-4-0":                         8192,
	"anthropic.claude-opus-4-20250514-v1:0":   8192,
	"claude-opus-4@20250514":                  8192,
	"claude-opus-4-1-20250805":                8192,
	"anthropic.claude-opus-4-1-20250805-v1:0": 8192,
	"claude-opus-4-1@20250805":                8192,
}

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

type (
	Any                                     string // Always "any"
	APIError                                string // Always "api_error"
	ApplicationPDF                          string // Always "application/pdf"
	Approximate                             string // Always "approximate"
	Assistant                               string // Always "assistant"
	AuthenticationError                     string // Always "authentication_error"
	Auto                                    string // Always "auto"
	Base64                                  string // Always "base64"
	Bash                                    string // Always "bash"
	Bash20241022                            string // Always "bash_20241022"
	Bash20250124                            string // Always "bash_20250124"
	BashCodeExecutionOutput                 string // Always "bash_code_execution_output"
	BashCodeExecutionResult                 string // Always "bash_code_execution_result"
	BashCodeExecutionToolResult             string // Always "bash_code_execution_tool_result"
	BashCodeExecutionToolResultError        string // Always "bash_code_execution_tool_result_error"
	BillingError                            string // Always "billing_error"
	Canceled                                string // Always "canceled"
	CharLocation                            string // Always "char_location"
	CitationsDelta                          string // Always "citations_delta"
	CodeExecution                           string // Always "code_execution"
	CodeExecution20250522                   string // Always "code_execution_20250522"
	CodeExecution20250825                   string // Always "code_execution_20250825"
	CodeExecutionOutput                     string // Always "code_execution_output"
	CodeExecutionResult                     string // Always "code_execution_result"
	CodeExecutionToolResult                 string // Always "code_execution_tool_result"
	CodeExecutionToolResultError            string // Always "code_execution_tool_result_error"
	Completion                              string // Always "completion"
	Computer                                string // Always "computer"
	Computer20241022                        string // Always "computer_20241022"
	Computer20250124                        string // Always "computer_20250124"
	ContainerUpload                         string // Always "container_upload"
	Content                                 string // Always "content"
	ContentBlockDelta                       string // Always "content_block_delta"
	ContentBlockLocation                    string // Always "content_block_location"
	ContentBlockStart                       string // Always "content_block_start"
	ContentBlockStop                        string // Always "content_block_stop"
	Disabled                                string // Always "disabled"
	Document                                string // Always "document"
	Enabled                                 string // Always "enabled"
	Ephemeral                               string // Always "ephemeral"
	Error                                   string // Always "error"
	Errored                                 string // Always "errored"
	Expired                                 string // Always "expired"
	File                                    string // Always "file"
	Image                                   string // Always "image"
	InputJSONDelta                          string // Always "input_json_delta"
	InvalidRequestError                     string // Always "invalid_request_error"
	MCPToolResult                           string // Always "mcp_tool_result"
	MCPToolUse                              string // Always "mcp_tool_use"
	Message                                 string // Always "message"
	MessageBatch                            string // Always "message_batch"
	MessageBatchDeleted                     string // Always "message_batch_deleted"
	MessageDelta                            string // Always "message_delta"
	MessageStart                            string // Always "message_start"
	MessageStop                             string // Always "message_stop"
	Model                                   string // Always "model"
	None                                    string // Always "none"
	NotFoundError                           string // Always "not_found_error"
	Object                                  string // Always "object"
	OverloadedError                         string // Always "overloaded_error"
	PageLocation                            string // Always "page_location"
	PermissionError                         string // Always "permission_error"
	RateLimitError                          string // Always "rate_limit_error"
	RedactedThinking                        string // Always "redacted_thinking"
	SearchResult                            string // Always "search_result"
	SearchResultLocation                    string // Always "search_result_location"
	ServerToolUse                           string // Always "server_tool_use"
	SignatureDelta                          string // Always "signature_delta"
	StrReplaceBasedEditTool                 string // Always "str_replace_based_edit_tool"
	StrReplaceEditor                        string // Always "str_replace_editor"
	Succeeded                               string // Always "succeeded"
	Text                                    string // Always "text"
	TextDelta                               string // Always "text_delta"
	TextEditor20241022                      string // Always "text_editor_20241022"
	TextEditor20250124                      string // Always "text_editor_20250124"
	TextEditor20250429                      string // Always "text_editor_20250429"
	TextEditor20250728                      string // Always "text_editor_20250728"
	TextEditorCodeExecutionCreateResult     string // Always "text_editor_code_execution_create_result"
	TextEditorCodeExecutionStrReplaceResult string // Always "text_editor_code_execution_str_replace_result"
	TextEditorCodeExecutionToolResult       string // Always "text_editor_code_execution_tool_result"
	TextEditorCodeExecutionToolResultError  string // Always "text_editor_code_execution_tool_result_error"
	TextEditorCodeExecutionViewResult       string // Always "text_editor_code_execution_view_result"
	TextPlain                               string // Always "text/plain"
	Thinking                                string // Always "thinking"
	ThinkingDelta                           string // Always "thinking_delta"
	TimeoutError                            string // Always "timeout_error"
	Tool                                    string // Always "tool"
	ToolResult                              string // Always "tool_result"
	ToolUse                                 string // Always "tool_use"
	URL                                     string // Always "url"
	WebFetch                                string // Always "web_fetch"
	WebFetch20250910                        string // Always "web_fetch_20250910"
	WebFetchResult                          string // Always "web_fetch_result"
	WebFetchToolResult                      string // Always "web_fetch_tool_result"
	WebFetchToolResultError                 string // Always "web_fetch_tool_result_error"
	WebSearch                               string // Always "web_search"
	WebSearch20250305                       string // Always "web_search_20250305"
	WebSearchResult                         string // Always "web_search_result"
	WebSearchResultLocation                 string // Always "web_search_result_location"
	WebSearchToolResult                     string // Always "web_search_tool_result"
	WebSearchToolResultError                string // Always "web_search_tool_result_error"
)

func (c Any) Default() Any                                 { return "any" }
func (c APIError) Default() APIError                       { return "api_error" }
func (c ApplicationPDF) Default() ApplicationPDF           { return "application/pdf" }
func (c Approximate) Default() Approximate                 { return "approximate" }
func (c Assistant) Default() Assistant                     { return "assistant" }
func (c AuthenticationError) Default() AuthenticationError { return "authentication_error" }
func (c Auto) Default() Auto                               { return "auto" }
func (c Base64) Default() Base64                           { return "base64" }
func (c Bash) Default() Bash                               { return "bash" }
func (c Bash20241022) Default() Bash20241022               { return "bash_20241022" }
func (c Bash20250124) Default() Bash20250124               { return "bash_20250124" }
func (c BashCodeExecutionOutput) Default() BashCodeExecutionOutput {
	return "bash_code_execution_output"
}

func (c BashCodeExecutionResult) Default() BashCodeExecutionResult {
	return "bash_code_execution_result"
}

func (c BashCodeExecutionToolResult) Default() BashCodeExecutionToolResult {
	return "bash_code_execution_tool_result"
}

func (c BashCodeExecutionToolResultError) Default() BashCodeExecutionToolResultError {
	return "bash_code_execution_tool_result_error"
}
func (c BillingError) Default() BillingError                   { return "billing_error" }
func (c Canceled) Default() Canceled                           { return "canceled" }
func (c CharLocation) Default() CharLocation                   { return "char_location" }
func (c CitationsDelta) Default() CitationsDelta               { return "citations_delta" }
func (c CodeExecution) Default() CodeExecution                 { return "code_execution" }
func (c CodeExecution20250522) Default() CodeExecution20250522 { return "code_execution_20250522" }
func (c CodeExecution20250825) Default() CodeExecution20250825 { return "code_execution_20250825" }
func (c CodeExecutionOutput) Default() CodeExecutionOutput     { return "code_execution_output" }
func (c CodeExecutionResult) Default() CodeExecutionResult     { return "code_execution_result" }
func (c CodeExecutionToolResult) Default() CodeExecutionToolResult {
	return "code_execution_tool_result"
}

func (c CodeExecutionToolResultError) Default() CodeExecutionToolResultError {
	return "code_execution_tool_result_error"
}
func (c Completion) Default() Completion                     { return "completion" }
func (c Computer) Default() Computer                         { return "computer" }
func (c Computer20241022) Default() Computer20241022         { return "computer_20241022" }
func (c Computer20250124) Default() Computer20250124         { return "computer_20250124" }
func (c ContainerUpload) Default() ContainerUpload           { return "container_upload" }
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
func (c File) Default() File                                 { return "file" }
func (c Image) Default() Image                               { return "image" }
func (c InputJSONDelta) Default() InputJSONDelta             { return "input_json_delta" }
func (c InvalidRequestError) Default() InvalidRequestError   { return "invalid_request_error" }
func (c MCPToolResult) Default() MCPToolResult               { return "mcp_tool_result" }
func (c MCPToolUse) Default() MCPToolUse                     { return "mcp_tool_use" }
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
func (c SearchResult) Default() SearchResult                 { return "search_result" }
func (c SearchResultLocation) Default() SearchResultLocation { return "search_result_location" }
func (c ServerToolUse) Default() ServerToolUse               { return "server_tool_use" }
func (c SignatureDelta) Default() SignatureDelta             { return "signature_delta" }
func (c StrReplaceBasedEditTool) Default() StrReplaceBasedEditTool {
	return "str_replace_based_edit_tool"
}
func (c StrReplaceEditor) Default() StrReplaceEditor     { return "str_replace_editor" }
func (c Succeeded) Default() Succeeded                   { return "succeeded" }
func (c Text) Default() Text                             { return "text" }
func (c TextDelta) Default() TextDelta                   { return "text_delta" }
func (c TextEditor20241022) Default() TextEditor20241022 { return "text_editor_20241022" }
func (c TextEditor20250124) Default() TextEditor20250124 { return "text_editor_20250124" }
func (c TextEditor20250429) Default() TextEditor20250429 { return "text_editor_20250429" }
func (c TextEditor20250728) Default() TextEditor20250728 { return "text_editor_20250728" }
func (c TextEditorCodeExecutionCreateResult) Default() TextEditorCodeExecutionCreateResult {
	return "text_editor_code_execution_create_result"
}

func (c TextEditorCodeExecutionStrReplaceResult) Default() TextEditorCodeExecutionStrReplaceResult {
	return "text_editor_code_execution_str_replace_result"
}

func (c TextEditorCodeExecutionToolResult) Default() TextEditorCodeExecutionToolResult {
	return "text_editor_code_execution_tool_result"
}

func (c TextEditorCodeExecutionToolResultError) Default() TextEditorCodeExecutionToolResultError {
	return "text_editor_code_execution_tool_result_error"
}

func (c TextEditorCodeExecutionViewResult) Default() TextEditorCodeExecutionViewResult {
	return "text_editor_code_execution_view_result"
}
func (c TextPlain) Default() TextPlain                   { return "text/plain" }
func (c Thinking) Default() Thinking                     { return "thinking" }
func (c ThinkingDelta) Default() ThinkingDelta           { return "thinking_delta" }
func (c TimeoutError) Default() TimeoutError             { return "timeout_error" }
func (c Tool) Default() Tool                             { return "tool" }
func (c ToolResult) Default() ToolResult                 { return "tool_result" }
func (c ToolUse) Default() ToolUse                       { return "tool_use" }
func (c URL) Default() URL                               { return "url" }
func (c WebFetch) Default() WebFetch                     { return "web_fetch" }
func (c WebFetch20250910) Default() WebFetch20250910     { return "web_fetch_20250910" }
func (c WebFetchResult) Default() WebFetchResult         { return "web_fetch_result" }
func (c WebFetchToolResult) Default() WebFetchToolResult { return "web_fetch_tool_result" }
func (c WebFetchToolResultError) Default() WebFetchToolResultError {
	return "web_fetch_tool_result_error"
}
func (c WebSearch) Default() WebSearch                 { return "web_search" }
func (c WebSearch20250305) Default() WebSearch20250305 { return "web_search_20250305" }
func (c WebSearchResult) Default() WebSearchResult     { return "web_search_result" }
func (c WebSearchResultLocation) Default() WebSearchResultLocation {
	return "web_search_result_location"
}
func (c WebSearchToolResult) Default() WebSearchToolResult { return "web_search_tool_result" }
func (c WebSearchToolResultError) Default() WebSearchToolResultError {
	return "web_search_tool_result_error"
}

func (c Any) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c APIError) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c ApplicationPDF) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c Approximate) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c Assistant) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c AuthenticationError) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c Auto) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c Base64) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c Bash) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c Bash20241022) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c Bash20250124) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c BashCodeExecutionOutput) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c BashCodeExecutionResult) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c BashCodeExecutionToolResult) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c BashCodeExecutionToolResultError) MarshalJSON() ([]byte, error)    { return marshalString(c) }
func (c BillingError) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c Canceled) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c CharLocation) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c CitationsDelta) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c CodeExecution) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c CodeExecution20250522) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c CodeExecution20250825) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c CodeExecutionOutput) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c CodeExecutionResult) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c CodeExecutionToolResult) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c CodeExecutionToolResultError) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c Completion) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c Computer) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c Computer20241022) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c Computer20250124) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c ContainerUpload) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c Content) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c ContentBlockDelta) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c ContentBlockLocation) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c ContentBlockStart) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c ContentBlockStop) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c Disabled) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c Document) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c Enabled) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c Ephemeral) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c Error) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c Errored) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c Expired) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c File) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c Image) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c InputJSONDelta) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c InvalidRequestError) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c MCPToolResult) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c MCPToolUse) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c Message) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c MessageBatch) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c MessageBatchDeleted) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c MessageDelta) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c MessageStart) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c MessageStop) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c Model) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c None) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c NotFoundError) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c Object) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c OverloadedError) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c PageLocation) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c PermissionError) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c RateLimitError) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c RedactedThinking) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c SearchResult) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c SearchResultLocation) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c ServerToolUse) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c SignatureDelta) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c StrReplaceBasedEditTool) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c StrReplaceEditor) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c Succeeded) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c Text) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c TextDelta) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c TextEditor20241022) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c TextEditor20250124) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c TextEditor20250429) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c TextEditor20250728) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c TextEditorCodeExecutionCreateResult) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c TextEditorCodeExecutionStrReplaceResult) MarshalJSON() ([]byte, error) {
	return marshalString(c)
}
func (c TextEditorCodeExecutionToolResult) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c TextEditorCodeExecutionToolResultError) MarshalJSON() ([]byte, error) {
	return marshalString(c)
}
func (c TextEditorCodeExecutionViewResult) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c TextPlain) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c Thinking) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c ThinkingDelta) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c TimeoutError) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c Tool) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c ToolResult) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c ToolUse) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c URL) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c WebFetch) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c WebFetch20250910) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c WebFetchResult) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c WebFetchToolResult) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c WebFetchToolResultError) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c WebSearch) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c WebSearch20250305) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c WebSearchResult) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c WebSearchResultLocation) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c WebSearchToolResult) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c WebSearchToolResultError) MarshalJSON() ([]byte, error)          { return marshalString(c) }

type constant[T any] interface {
	Constant[T]
	*T
}

func marshalString[T ~string, PT constant[T]](v T) ([]byte, error) {
	var zero T
	if v == zero {
		v = PT(&v).Default()
	}
	return shimjson.Marshal(string(v))
}
