// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
)

// BetaSessionEventService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaSessionEventService] method instead.
type BetaSessionEventService struct {
	Options []option.RequestOption
}

// NewBetaSessionEventService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaSessionEventService(opts ...option.RequestOption) (r BetaSessionEventService) {
	r = BetaSessionEventService{}
	r.Options = opts
	return
}

// List Events
func (r *BetaSessionEventService) List(ctx context.Context, sessionID string, params BetaSessionEventListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaManagedAgentsSessionEventUnion], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01"), option.WithResponseInto(&raw)}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/sessions/%s/events?beta=true", sessionID)
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, params, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// List Events
func (r *BetaSessionEventService) ListAutoPaging(ctx context.Context, sessionID string, params BetaSessionEventListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsSessionEventUnion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, sessionID, params, opts...))
}

// Send Events
func (r *BetaSessionEventService) Send(ctx context.Context, sessionID string, params BetaSessionEventSendParams, opts ...option.RequestOption) (res *BetaManagedAgentsSendSessionEvents, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/sessions/%s/events?beta=true", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Stream Events
func (r *BetaSessionEventService) StreamEvents(ctx context.Context, sessionID string, query BetaSessionEventStreamParams, opts ...option.RequestOption) (stream *ssestream.Stream[BetaManagedAgentsStreamSessionEventsUnion]) {
	var (
		raw *http.Response
		err error
	)
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return ssestream.NewStream[BetaManagedAgentsStreamSessionEventsUnion](nil, err)
	}
	path := fmt.Sprintf("v1/sessions/%s/events/stream?beta=true", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &raw, opts...)
	return ssestream.NewStream[BetaManagedAgentsStreamSessionEventsUnion](ssestream.NewDecoder(raw), err)
}

// Event emitted when the agent calls a custom tool. The session goes idle until
// the client sends a `user.custom_tool_result` event with the result.
type BetaManagedAgentsAgentCustomToolUseEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Input parameters for the tool call.
	Input map[string]any `json:"input" api:"required"`
	// Name of the custom tool being called.
	Name string `json:"name" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "agent.custom_tool_use".
	Type BetaManagedAgentsAgentCustomToolUseEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Input       respjson.Field
		Name        respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentCustomToolUseEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentCustomToolUseEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAgentCustomToolUseEventType string

const (
	BetaManagedAgentsAgentCustomToolUseEventTypeAgentCustomToolUse BetaManagedAgentsAgentCustomToolUseEventType = "agent.custom_tool_use"
)

// Event representing the result of an MCP tool execution.
type BetaManagedAgentsAgentMCPToolResultEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// The id of the `agent.mcp_tool_use` event this result corresponds to.
	MCPToolUseID string `json:"mcp_tool_use_id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "agent.mcp_tool_result".
	Type BetaManagedAgentsAgentMCPToolResultEventType `json:"type" api:"required"`
	// The result content returned by the tool.
	Content []BetaManagedAgentsAgentMCPToolResultEventContentUnion `json:"content"`
	// Whether the tool execution resulted in an error.
	IsError bool `json:"is_error" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		MCPToolUseID respjson.Field
		ProcessedAt  respjson.Field
		Type         respjson.Field
		Content      respjson.Field
		IsError      respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentMCPToolResultEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentMCPToolResultEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAgentMCPToolResultEventType string

const (
	BetaManagedAgentsAgentMCPToolResultEventTypeAgentMCPToolResult BetaManagedAgentsAgentMCPToolResultEventType = "agent.mcp_tool_result"
)

// BetaManagedAgentsAgentMCPToolResultEventContentUnion contains all possible
// properties and values from [BetaManagedAgentsTextBlock],
// [BetaManagedAgentsImageBlock], [BetaManagedAgentsDocumentBlock].
//
// Use the [BetaManagedAgentsAgentMCPToolResultEventContentUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsAgentMCPToolResultEventContentUnion struct {
	// This field is from variant [BetaManagedAgentsTextBlock].
	Text string `json:"text"`
	// Any of "text", "image", "document".
	Type string `json:"type"`
	// This field is a union of [BetaManagedAgentsImageBlockSourceUnion],
	// [BetaManagedAgentsDocumentBlockSourceUnion]
	Source BetaManagedAgentsAgentMCPToolResultEventContentUnionSource `json:"source"`
	// This field is from variant [BetaManagedAgentsDocumentBlock].
	Context string `json:"context"`
	// This field is from variant [BetaManagedAgentsDocumentBlock].
	Title string `json:"title"`
	JSON  struct {
		Text    respjson.Field
		Type    respjson.Field
		Source  respjson.Field
		Context respjson.Field
		Title   respjson.Field
		raw     string
	} `json:"-"`
}

// anyBetaManagedAgentsAgentMCPToolResultEventContent is implemented by each
// variant of [BetaManagedAgentsAgentMCPToolResultEventContentUnion] to add type
// safety for the return type of
// [BetaManagedAgentsAgentMCPToolResultEventContentUnion.AsAny]
type anyBetaManagedAgentsAgentMCPToolResultEventContent interface {
	implBetaManagedAgentsAgentMCPToolResultEventContentUnion()
}

func (BetaManagedAgentsTextBlock) implBetaManagedAgentsAgentMCPToolResultEventContentUnion()     {}
func (BetaManagedAgentsImageBlock) implBetaManagedAgentsAgentMCPToolResultEventContentUnion()    {}
func (BetaManagedAgentsDocumentBlock) implBetaManagedAgentsAgentMCPToolResultEventContentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsAgentMCPToolResultEventContentUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsTextBlock:
//	case anthropic.BetaManagedAgentsImageBlock:
//	case anthropic.BetaManagedAgentsDocumentBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsAgentMCPToolResultEventContentUnion) AsAny() anyBetaManagedAgentsAgentMCPToolResultEventContent {
	switch u.Type {
	case "text":
		return u.AsText()
	case "image":
		return u.AsImage()
	case "document":
		return u.AsDocument()
	}
	return nil
}

func (u BetaManagedAgentsAgentMCPToolResultEventContentUnion) AsText() (v BetaManagedAgentsTextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsAgentMCPToolResultEventContentUnion) AsImage() (v BetaManagedAgentsImageBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsAgentMCPToolResultEventContentUnion) AsDocument() (v BetaManagedAgentsDocumentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsAgentMCPToolResultEventContentUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsAgentMCPToolResultEventContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsAgentMCPToolResultEventContentUnionSource is an implicit
// subunion of [BetaManagedAgentsAgentMCPToolResultEventContentUnion].
// BetaManagedAgentsAgentMCPToolResultEventContentUnionSource provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsAgentMCPToolResultEventContentUnion].
type BetaManagedAgentsAgentMCPToolResultEventContentUnionSource struct {
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	FileID    string `json:"file_id"`
	JSON      struct {
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

func (r *BetaManagedAgentsAgentMCPToolResultEventContentUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event emitted when the agent invokes a tool provided by an MCP server.
type BetaManagedAgentsAgentMCPToolUseEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Input parameters for the tool call.
	Input map[string]any `json:"input" api:"required"`
	// Name of the MCP server providing the tool.
	MCPServerName string `json:"mcp_server_name" api:"required"`
	// Name of the MCP tool being used.
	Name string `json:"name" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "agent.mcp_tool_use".
	Type BetaManagedAgentsAgentMCPToolUseEventType `json:"type" api:"required"`
	// AgentEvaluatedPermission enum
	//
	// Any of "allow", "ask", "deny".
	EvaluatedPermission BetaManagedAgentsAgentMCPToolUseEventEvaluatedPermission `json:"evaluated_permission"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                  respjson.Field
		Input               respjson.Field
		MCPServerName       respjson.Field
		Name                respjson.Field
		ProcessedAt         respjson.Field
		Type                respjson.Field
		EvaluatedPermission respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentMCPToolUseEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentMCPToolUseEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAgentMCPToolUseEventType string

const (
	BetaManagedAgentsAgentMCPToolUseEventTypeAgentMCPToolUse BetaManagedAgentsAgentMCPToolUseEventType = "agent.mcp_tool_use"
)

// AgentEvaluatedPermission enum
type BetaManagedAgentsAgentMCPToolUseEventEvaluatedPermission string

const (
	BetaManagedAgentsAgentMCPToolUseEventEvaluatedPermissionAllow BetaManagedAgentsAgentMCPToolUseEventEvaluatedPermission = "allow"
	BetaManagedAgentsAgentMCPToolUseEventEvaluatedPermissionAsk   BetaManagedAgentsAgentMCPToolUseEventEvaluatedPermission = "ask"
	BetaManagedAgentsAgentMCPToolUseEventEvaluatedPermissionDeny  BetaManagedAgentsAgentMCPToolUseEventEvaluatedPermission = "deny"
)

// An agent response event in the session conversation.
type BetaManagedAgentsAgentMessageEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Array of text blocks comprising the agent response.
	Content []BetaManagedAgentsTextBlock `json:"content" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "agent.message".
	Type BetaManagedAgentsAgentMessageEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Content     respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentMessageEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentMessageEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAgentMessageEventType string

const (
	BetaManagedAgentsAgentMessageEventTypeAgentMessage BetaManagedAgentsAgentMessageEventType = "agent.message"
)

// Indicates the agent is making forward progress via extended thinking. A progress
// signal, not a content carrier.
type BetaManagedAgentsAgentThinkingEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "agent.thinking".
	Type BetaManagedAgentsAgentThinkingEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentThinkingEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentThinkingEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAgentThinkingEventType string

const (
	BetaManagedAgentsAgentThinkingEventTypeAgentThinking BetaManagedAgentsAgentThinkingEventType = "agent.thinking"
)

// Indicates that context compaction (summarization) occurred during the session.
type BetaManagedAgentsAgentThreadContextCompactedEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "agent.thread_context_compacted".
	Type BetaManagedAgentsAgentThreadContextCompactedEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentThreadContextCompactedEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentThreadContextCompactedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAgentThreadContextCompactedEventType string

const (
	BetaManagedAgentsAgentThreadContextCompactedEventTypeAgentThreadContextCompacted BetaManagedAgentsAgentThreadContextCompactedEventType = "agent.thread_context_compacted"
)

// Event representing the result of an agent tool execution.
type BetaManagedAgentsAgentToolResultEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// The id of the `agent.tool_use` event this result corresponds to.
	ToolUseID string `json:"tool_use_id" api:"required"`
	// Any of "agent.tool_result".
	Type BetaManagedAgentsAgentToolResultEventType `json:"type" api:"required"`
	// The result content returned by the tool.
	Content []BetaManagedAgentsAgentToolResultEventContentUnion `json:"content"`
	// Whether the tool execution resulted in an error.
	IsError bool `json:"is_error" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		ToolUseID   respjson.Field
		Type        respjson.Field
		Content     respjson.Field
		IsError     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentToolResultEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentToolResultEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAgentToolResultEventType string

const (
	BetaManagedAgentsAgentToolResultEventTypeAgentToolResult BetaManagedAgentsAgentToolResultEventType = "agent.tool_result"
)

// BetaManagedAgentsAgentToolResultEventContentUnion contains all possible
// properties and values from [BetaManagedAgentsTextBlock],
// [BetaManagedAgentsImageBlock], [BetaManagedAgentsDocumentBlock].
//
// Use the [BetaManagedAgentsAgentToolResultEventContentUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsAgentToolResultEventContentUnion struct {
	// This field is from variant [BetaManagedAgentsTextBlock].
	Text string `json:"text"`
	// Any of "text", "image", "document".
	Type string `json:"type"`
	// This field is a union of [BetaManagedAgentsImageBlockSourceUnion],
	// [BetaManagedAgentsDocumentBlockSourceUnion]
	Source BetaManagedAgentsAgentToolResultEventContentUnionSource `json:"source"`
	// This field is from variant [BetaManagedAgentsDocumentBlock].
	Context string `json:"context"`
	// This field is from variant [BetaManagedAgentsDocumentBlock].
	Title string `json:"title"`
	JSON  struct {
		Text    respjson.Field
		Type    respjson.Field
		Source  respjson.Field
		Context respjson.Field
		Title   respjson.Field
		raw     string
	} `json:"-"`
}

// anyBetaManagedAgentsAgentToolResultEventContent is implemented by each variant
// of [BetaManagedAgentsAgentToolResultEventContentUnion] to add type safety for
// the return type of [BetaManagedAgentsAgentToolResultEventContentUnion.AsAny]
type anyBetaManagedAgentsAgentToolResultEventContent interface {
	implBetaManagedAgentsAgentToolResultEventContentUnion()
}

func (BetaManagedAgentsTextBlock) implBetaManagedAgentsAgentToolResultEventContentUnion()     {}
func (BetaManagedAgentsImageBlock) implBetaManagedAgentsAgentToolResultEventContentUnion()    {}
func (BetaManagedAgentsDocumentBlock) implBetaManagedAgentsAgentToolResultEventContentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsAgentToolResultEventContentUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsTextBlock:
//	case anthropic.BetaManagedAgentsImageBlock:
//	case anthropic.BetaManagedAgentsDocumentBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsAgentToolResultEventContentUnion) AsAny() anyBetaManagedAgentsAgentToolResultEventContent {
	switch u.Type {
	case "text":
		return u.AsText()
	case "image":
		return u.AsImage()
	case "document":
		return u.AsDocument()
	}
	return nil
}

func (u BetaManagedAgentsAgentToolResultEventContentUnion) AsText() (v BetaManagedAgentsTextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsAgentToolResultEventContentUnion) AsImage() (v BetaManagedAgentsImageBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsAgentToolResultEventContentUnion) AsDocument() (v BetaManagedAgentsDocumentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsAgentToolResultEventContentUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsAgentToolResultEventContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsAgentToolResultEventContentUnionSource is an implicit subunion
// of [BetaManagedAgentsAgentToolResultEventContentUnion].
// BetaManagedAgentsAgentToolResultEventContentUnionSource provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsAgentToolResultEventContentUnion].
type BetaManagedAgentsAgentToolResultEventContentUnionSource struct {
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	FileID    string `json:"file_id"`
	JSON      struct {
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

func (r *BetaManagedAgentsAgentToolResultEventContentUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event emitted when the agent invokes a built-in agent tool.
type BetaManagedAgentsAgentToolUseEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Input parameters for the tool call.
	Input map[string]any `json:"input" api:"required"`
	// Name of the agent tool being used.
	Name string `json:"name" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "agent.tool_use".
	Type BetaManagedAgentsAgentToolUseEventType `json:"type" api:"required"`
	// AgentEvaluatedPermission enum
	//
	// Any of "allow", "ask", "deny".
	EvaluatedPermission BetaManagedAgentsAgentToolUseEventEvaluatedPermission `json:"evaluated_permission"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                  respjson.Field
		Input               respjson.Field
		Name                respjson.Field
		ProcessedAt         respjson.Field
		Type                respjson.Field
		EvaluatedPermission respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentToolUseEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentToolUseEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAgentToolUseEventType string

const (
	BetaManagedAgentsAgentToolUseEventTypeAgentToolUse BetaManagedAgentsAgentToolUseEventType = "agent.tool_use"
)

// AgentEvaluatedPermission enum
type BetaManagedAgentsAgentToolUseEventEvaluatedPermission string

const (
	BetaManagedAgentsAgentToolUseEventEvaluatedPermissionAllow BetaManagedAgentsAgentToolUseEventEvaluatedPermission = "allow"
	BetaManagedAgentsAgentToolUseEventEvaluatedPermissionAsk   BetaManagedAgentsAgentToolUseEventEvaluatedPermission = "ask"
	BetaManagedAgentsAgentToolUseEventEvaluatedPermissionDeny  BetaManagedAgentsAgentToolUseEventEvaluatedPermission = "deny"
)

// Base64-encoded document data.
type BetaManagedAgentsBase64DocumentSource struct {
	// Base64-encoded document data.
	Data string `json:"data" api:"required"`
	// MIME type of the document (e.g., "application/pdf").
	MediaType string `json:"media_type" api:"required"`
	// Any of "base64".
	Type BetaManagedAgentsBase64DocumentSourceType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		MediaType   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsBase64DocumentSource) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsBase64DocumentSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsBase64DocumentSource to a
// BetaManagedAgentsBase64DocumentSourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsBase64DocumentSourceParam.Overrides()
func (r BetaManagedAgentsBase64DocumentSource) ToParam() BetaManagedAgentsBase64DocumentSourceParam {
	return param.Override[BetaManagedAgentsBase64DocumentSourceParam](json.RawMessage(r.RawJSON()))
}

type BetaManagedAgentsBase64DocumentSourceType string

const (
	BetaManagedAgentsBase64DocumentSourceTypeBase64 BetaManagedAgentsBase64DocumentSourceType = "base64"
)

// Base64-encoded document data.
//
// The properties Data, MediaType, Type are required.
type BetaManagedAgentsBase64DocumentSourceParam struct {
	// Base64-encoded document data.
	Data string `json:"data" api:"required"`
	// MIME type of the document (e.g., "application/pdf").
	MediaType string `json:"media_type" api:"required"`
	// Any of "base64".
	Type BetaManagedAgentsBase64DocumentSourceType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsBase64DocumentSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsBase64DocumentSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsBase64DocumentSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Base64-encoded image data.
type BetaManagedAgentsBase64ImageSource struct {
	// Base64-encoded image data.
	Data string `json:"data" api:"required"`
	// MIME type of the image (e.g., "image/png", "image/jpeg", "image/gif",
	// "image/webp").
	MediaType string `json:"media_type" api:"required"`
	// Any of "base64".
	Type BetaManagedAgentsBase64ImageSourceType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		MediaType   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsBase64ImageSource) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsBase64ImageSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsBase64ImageSource to a
// BetaManagedAgentsBase64ImageSourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsBase64ImageSourceParam.Overrides()
func (r BetaManagedAgentsBase64ImageSource) ToParam() BetaManagedAgentsBase64ImageSourceParam {
	return param.Override[BetaManagedAgentsBase64ImageSourceParam](json.RawMessage(r.RawJSON()))
}

type BetaManagedAgentsBase64ImageSourceType string

const (
	BetaManagedAgentsBase64ImageSourceTypeBase64 BetaManagedAgentsBase64ImageSourceType = "base64"
)

// Base64-encoded image data.
//
// The properties Data, MediaType, Type are required.
type BetaManagedAgentsBase64ImageSourceParam struct {
	// Base64-encoded image data.
	Data string `json:"data" api:"required"`
	// MIME type of the image (e.g., "image/png", "image/jpeg", "image/gif",
	// "image/webp").
	MediaType string `json:"media_type" api:"required"`
	// Any of "base64".
	Type BetaManagedAgentsBase64ImageSourceType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsBase64ImageSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsBase64ImageSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsBase64ImageSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The caller's organization or workspace cannot make model requests — out of
// credits or spend limit reached. Retrying with the same credentials will not
// succeed; the caller must resolve the billing state.
type BetaManagedAgentsBillingError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// What the client should do next in response to this error.
	RetryStatus BetaManagedAgentsBillingErrorRetryStatusUnion `json:"retry_status" api:"required"`
	// Any of "billing_error".
	Type BetaManagedAgentsBillingErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		RetryStatus respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsBillingError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsBillingError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsBillingErrorRetryStatusUnion contains all possible properties
// and values from [BetaManagedAgentsRetryStatusRetrying],
// [BetaManagedAgentsRetryStatusExhausted], [BetaManagedAgentsRetryStatusTerminal].
//
// Use the [BetaManagedAgentsBillingErrorRetryStatusUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsBillingErrorRetryStatusUnion struct {
	// Any of "retrying", "exhausted", "terminal".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsBillingErrorRetryStatus is implemented by each variant of
// [BetaManagedAgentsBillingErrorRetryStatusUnion] to add type safety for the
// return type of [BetaManagedAgentsBillingErrorRetryStatusUnion.AsAny]
type anyBetaManagedAgentsBillingErrorRetryStatus interface {
	implBetaManagedAgentsBillingErrorRetryStatusUnion()
}

func (BetaManagedAgentsRetryStatusRetrying) implBetaManagedAgentsBillingErrorRetryStatusUnion()  {}
func (BetaManagedAgentsRetryStatusExhausted) implBetaManagedAgentsBillingErrorRetryStatusUnion() {}
func (BetaManagedAgentsRetryStatusTerminal) implBetaManagedAgentsBillingErrorRetryStatusUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsBillingErrorRetryStatusUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsRetryStatusRetrying:
//	case anthropic.BetaManagedAgentsRetryStatusExhausted:
//	case anthropic.BetaManagedAgentsRetryStatusTerminal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsBillingErrorRetryStatusUnion) AsAny() anyBetaManagedAgentsBillingErrorRetryStatus {
	switch u.Type {
	case "retrying":
		return u.AsRetrying()
	case "exhausted":
		return u.AsExhausted()
	case "terminal":
		return u.AsTerminal()
	}
	return nil
}

func (u BetaManagedAgentsBillingErrorRetryStatusUnion) AsRetrying() (v BetaManagedAgentsRetryStatusRetrying) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsBillingErrorRetryStatusUnion) AsExhausted() (v BetaManagedAgentsRetryStatusExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsBillingErrorRetryStatusUnion) AsTerminal() (v BetaManagedAgentsRetryStatusTerminal) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsBillingErrorRetryStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsBillingErrorRetryStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsBillingErrorType string

const (
	BetaManagedAgentsBillingErrorTypeBillingError BetaManagedAgentsBillingErrorType = "billing_error"
)

// Document content, either specified directly as base64 data, as text, or as a
// reference via a URL.
type BetaManagedAgentsDocumentBlock struct {
	// Union type for document source variants.
	Source BetaManagedAgentsDocumentBlockSourceUnion `json:"source" api:"required"`
	// Any of "document".
	Type BetaManagedAgentsDocumentBlockType `json:"type" api:"required"`
	// Additional context about the document for the model.
	Context string `json:"context" api:"nullable"`
	// The title of the document.
	Title string `json:"title" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Source      respjson.Field
		Type        respjson.Field
		Context     respjson.Field
		Title       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsDocumentBlock) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsDocumentBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsDocumentBlock to a
// BetaManagedAgentsDocumentBlockParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsDocumentBlockParam.Overrides()
func (r BetaManagedAgentsDocumentBlock) ToParam() BetaManagedAgentsDocumentBlockParam {
	return param.Override[BetaManagedAgentsDocumentBlockParam](json.RawMessage(r.RawJSON()))
}

// BetaManagedAgentsDocumentBlockSourceUnion contains all possible properties and
// values from [BetaManagedAgentsBase64DocumentSource],
// [BetaManagedAgentsPlainTextDocumentSource],
// [BetaManagedAgentsURLDocumentSource], [BetaManagedAgentsFileDocumentSource].
//
// Use the [BetaManagedAgentsDocumentBlockSourceUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsDocumentBlockSourceUnion struct {
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	// Any of "base64", "text", "url", "file".
	Type string `json:"type"`
	// This field is from variant [BetaManagedAgentsURLDocumentSource].
	URL string `json:"url"`
	// This field is from variant [BetaManagedAgentsFileDocumentSource].
	FileID string `json:"file_id"`
	JSON   struct {
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

// anyBetaManagedAgentsDocumentBlockSource is implemented by each variant of
// [BetaManagedAgentsDocumentBlockSourceUnion] to add type safety for the return
// type of [BetaManagedAgentsDocumentBlockSourceUnion.AsAny]
type anyBetaManagedAgentsDocumentBlockSource interface {
	implBetaManagedAgentsDocumentBlockSourceUnion()
}

func (BetaManagedAgentsBase64DocumentSource) implBetaManagedAgentsDocumentBlockSourceUnion()    {}
func (BetaManagedAgentsPlainTextDocumentSource) implBetaManagedAgentsDocumentBlockSourceUnion() {}
func (BetaManagedAgentsURLDocumentSource) implBetaManagedAgentsDocumentBlockSourceUnion()       {}
func (BetaManagedAgentsFileDocumentSource) implBetaManagedAgentsDocumentBlockSourceUnion()      {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsDocumentBlockSourceUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsBase64DocumentSource:
//	case anthropic.BetaManagedAgentsPlainTextDocumentSource:
//	case anthropic.BetaManagedAgentsURLDocumentSource:
//	case anthropic.BetaManagedAgentsFileDocumentSource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsDocumentBlockSourceUnion) AsAny() anyBetaManagedAgentsDocumentBlockSource {
	switch u.Type {
	case "base64":
		return u.AsBase64()
	case "text":
		return u.AsText()
	case "url":
		return u.AsURL()
	case "file":
		return u.AsFile()
	}
	return nil
}

func (u BetaManagedAgentsDocumentBlockSourceUnion) AsBase64() (v BetaManagedAgentsBase64DocumentSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDocumentBlockSourceUnion) AsText() (v BetaManagedAgentsPlainTextDocumentSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDocumentBlockSourceUnion) AsURL() (v BetaManagedAgentsURLDocumentSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDocumentBlockSourceUnion) AsFile() (v BetaManagedAgentsFileDocumentSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsDocumentBlockSourceUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsDocumentBlockSourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsDocumentBlockType string

const (
	BetaManagedAgentsDocumentBlockTypeDocument BetaManagedAgentsDocumentBlockType = "document"
)

// Document content, either specified directly as base64 data, as text, or as a
// reference via a URL.
//
// The properties Source, Type are required.
type BetaManagedAgentsDocumentBlockParam struct {
	// Union type for document source variants.
	Source BetaManagedAgentsDocumentBlockSourceUnionParam `json:"source,omitzero" api:"required"`
	// Any of "document".
	Type BetaManagedAgentsDocumentBlockType `json:"type,omitzero" api:"required"`
	// Additional context about the document for the model.
	Context param.Opt[string] `json:"context,omitzero"`
	// The title of the document.
	Title param.Opt[string] `json:"title,omitzero"`
	paramObj
}

func (r BetaManagedAgentsDocumentBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsDocumentBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsDocumentBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsDocumentBlockSourceUnionParam struct {
	OfBase64 *BetaManagedAgentsBase64DocumentSourceParam    `json:",omitzero,inline"`
	OfText   *BetaManagedAgentsPlainTextDocumentSourceParam `json:",omitzero,inline"`
	OfURL    *BetaManagedAgentsURLDocumentSourceParam       `json:",omitzero,inline"`
	OfFile   *BetaManagedAgentsFileDocumentSourceParam      `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsDocumentBlockSourceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfBase64, u.OfText, u.OfURL, u.OfFile)
}
func (u *BetaManagedAgentsDocumentBlockSourceUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsDocumentBlockSourceUnionParam) asAny() any {
	if !param.IsOmitted(u.OfBase64) {
		return u.OfBase64
	} else if !param.IsOmitted(u.OfText) {
		return u.OfText
	} else if !param.IsOmitted(u.OfURL) {
		return u.OfURL
	} else if !param.IsOmitted(u.OfFile) {
		return u.OfFile
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsDocumentBlockSourceUnionParam) GetURL() *string {
	if vt := u.OfURL; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsDocumentBlockSourceUnionParam) GetFileID() *string {
	if vt := u.OfFile; vt != nil {
		return &vt.FileID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsDocumentBlockSourceUnionParam) GetData() *string {
	if vt := u.OfBase64; vt != nil {
		return (*string)(&vt.Data)
	} else if vt := u.OfText; vt != nil {
		return (*string)(&vt.Data)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsDocumentBlockSourceUnionParam) GetMediaType() *string {
	if vt := u.OfBase64; vt != nil {
		return (*string)(&vt.MediaType)
	} else if vt := u.OfText; vt != nil {
		return (*string)(&vt.MediaType)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsDocumentBlockSourceUnionParam) GetType() *string {
	if vt := u.OfBase64; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfURL; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFile; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsDocumentBlockSourceUnionParam](
		"type",
		apijson.Discriminator[BetaManagedAgentsBase64DocumentSourceParam]("base64"),
		apijson.Discriminator[BetaManagedAgentsPlainTextDocumentSourceParam]("text"),
		apijson.Discriminator[BetaManagedAgentsURLDocumentSourceParam]("url"),
		apijson.Discriminator[BetaManagedAgentsFileDocumentSourceParam]("file"),
	)
}

func BetaManagedAgentsEventParamsOfUserMessage(content []BetaManagedAgentsUserMessageEventParamsContentUnion) BetaManagedAgentsEventParamsUnion {
	var userMessage BetaManagedAgentsUserMessageEventParams
	userMessage.Content = content
	return BetaManagedAgentsEventParamsUnion{OfUserMessage: &userMessage}
}

func BetaManagedAgentsEventParamsOfUserInterrupt(type_ BetaManagedAgentsUserInterruptEventParamsType) BetaManagedAgentsEventParamsUnion {
	var userInterrupt BetaManagedAgentsUserInterruptEventParams
	userInterrupt.Type = type_
	return BetaManagedAgentsEventParamsUnion{OfUserInterrupt: &userInterrupt}
}

func BetaManagedAgentsEventParamsOfUserToolConfirmation(result BetaManagedAgentsUserToolConfirmationEventParamsResult, toolUseID string, type_ BetaManagedAgentsUserToolConfirmationEventParamsType) BetaManagedAgentsEventParamsUnion {
	var userToolConfirmation BetaManagedAgentsUserToolConfirmationEventParams
	userToolConfirmation.Result = result
	userToolConfirmation.ToolUseID = toolUseID
	userToolConfirmation.Type = type_
	return BetaManagedAgentsEventParamsUnion{OfUserToolConfirmation: &userToolConfirmation}
}

func BetaManagedAgentsEventParamsOfUserCustomToolResult(customToolUseID string) BetaManagedAgentsEventParamsUnion {
	var userCustomToolResult BetaManagedAgentsUserCustomToolResultEventParams
	userCustomToolResult.CustomToolUseID = customToolUseID
	return BetaManagedAgentsEventParamsUnion{OfUserCustomToolResult: &userCustomToolResult}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsEventParamsUnion struct {
	OfUserMessage          *BetaManagedAgentsUserMessageEventParams          `json:",omitzero,inline"`
	OfUserInterrupt        *BetaManagedAgentsUserInterruptEventParams        `json:",omitzero,inline"`
	OfUserToolConfirmation *BetaManagedAgentsUserToolConfirmationEventParams `json:",omitzero,inline"`
	OfUserCustomToolResult *BetaManagedAgentsUserCustomToolResultEventParams `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsEventParamsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfUserMessage, u.OfUserInterrupt, u.OfUserToolConfirmation, u.OfUserCustomToolResult)
}
func (u *BetaManagedAgentsEventParamsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsEventParamsUnion) asAny() any {
	if !param.IsOmitted(u.OfUserMessage) {
		return u.OfUserMessage
	} else if !param.IsOmitted(u.OfUserInterrupt) {
		return u.OfUserInterrupt
	} else if !param.IsOmitted(u.OfUserToolConfirmation) {
		return u.OfUserToolConfirmation
	} else if !param.IsOmitted(u.OfUserCustomToolResult) {
		return u.OfUserCustomToolResult
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsEventParamsUnion) GetResult() *string {
	if vt := u.OfUserToolConfirmation; vt != nil {
		return (*string)(&vt.Result)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsEventParamsUnion) GetToolUseID() *string {
	if vt := u.OfUserToolConfirmation; vt != nil {
		return &vt.ToolUseID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsEventParamsUnion) GetDenyMessage() *string {
	if vt := u.OfUserToolConfirmation; vt != nil && vt.DenyMessage.Valid() {
		return &vt.DenyMessage.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsEventParamsUnion) GetCustomToolUseID() *string {
	if vt := u.OfUserCustomToolResult; vt != nil {
		return &vt.CustomToolUseID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsEventParamsUnion) GetIsError() *bool {
	if vt := u.OfUserCustomToolResult; vt != nil && vt.IsError.Valid() {
		return &vt.IsError.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsEventParamsUnion) GetType() *string {
	if vt := u.OfUserMessage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfUserInterrupt; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfUserToolConfirmation; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfUserCustomToolResult; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u BetaManagedAgentsEventParamsUnion) GetContent() (res betaManagedAgentsEventParamsUnionContent) {
	if vt := u.OfUserMessage; vt != nil {
		res.any = &vt.Content
	} else if vt := u.OfUserCustomToolResult; vt != nil {
		res.any = &vt.Content
	}
	return
}

// Can have the runtime types
// [_[]BetaManagedAgentsUserMessageEventParamsContentUnion],
// [_[]BetaManagedAgentsUserCustomToolResultEventParamsContentUnion]
type betaManagedAgentsEventParamsUnionContent struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]anthropic.BetaManagedAgentsUserMessageEventParamsContentUnion:
//	case *[]anthropic.BetaManagedAgentsUserCustomToolResultEventParamsContentUnion:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u betaManagedAgentsEventParamsUnionContent) AsAny() any { return u.any }

func init() {
	apijson.RegisterUnion[BetaManagedAgentsEventParamsUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsUserMessageEventParams]("user.message"),
		apijson.Discriminator[BetaManagedAgentsUserInterruptEventParams]("user.interrupt"),
		apijson.Discriminator[BetaManagedAgentsUserToolConfirmationEventParams]("user.tool_confirmation"),
		apijson.Discriminator[BetaManagedAgentsUserCustomToolResultEventParams]("user.custom_tool_result"),
	)
}

// Document referenced by file ID.
type BetaManagedAgentsFileDocumentSource struct {
	// ID of a previously uploaded file.
	FileID string `json:"file_id" api:"required"`
	// Any of "file".
	Type BetaManagedAgentsFileDocumentSourceType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FileID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsFileDocumentSource) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsFileDocumentSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsFileDocumentSource to a
// BetaManagedAgentsFileDocumentSourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsFileDocumentSourceParam.Overrides()
func (r BetaManagedAgentsFileDocumentSource) ToParam() BetaManagedAgentsFileDocumentSourceParam {
	return param.Override[BetaManagedAgentsFileDocumentSourceParam](json.RawMessage(r.RawJSON()))
}

type BetaManagedAgentsFileDocumentSourceType string

const (
	BetaManagedAgentsFileDocumentSourceTypeFile BetaManagedAgentsFileDocumentSourceType = "file"
)

// Document referenced by file ID.
//
// The properties FileID, Type are required.
type BetaManagedAgentsFileDocumentSourceParam struct {
	// ID of a previously uploaded file.
	FileID string `json:"file_id" api:"required"`
	// Any of "file".
	Type BetaManagedAgentsFileDocumentSourceType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsFileDocumentSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsFileDocumentSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsFileDocumentSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Image referenced by file ID.
type BetaManagedAgentsFileImageSource struct {
	// ID of a previously uploaded file.
	FileID string `json:"file_id" api:"required"`
	// Any of "file".
	Type BetaManagedAgentsFileImageSourceType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FileID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsFileImageSource) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsFileImageSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsFileImageSource to a
// BetaManagedAgentsFileImageSourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsFileImageSourceParam.Overrides()
func (r BetaManagedAgentsFileImageSource) ToParam() BetaManagedAgentsFileImageSourceParam {
	return param.Override[BetaManagedAgentsFileImageSourceParam](json.RawMessage(r.RawJSON()))
}

type BetaManagedAgentsFileImageSourceType string

const (
	BetaManagedAgentsFileImageSourceTypeFile BetaManagedAgentsFileImageSourceType = "file"
)

// Image referenced by file ID.
//
// The properties FileID, Type are required.
type BetaManagedAgentsFileImageSourceParam struct {
	// ID of a previously uploaded file.
	FileID string `json:"file_id" api:"required"`
	// Any of "file".
	Type BetaManagedAgentsFileImageSourceType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsFileImageSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsFileImageSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsFileImageSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Image content specified directly as base64 data or as a reference via a URL.
type BetaManagedAgentsImageBlock struct {
	// Union type for image source variants.
	Source BetaManagedAgentsImageBlockSourceUnion `json:"source" api:"required"`
	// Any of "image".
	Type BetaManagedAgentsImageBlockType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Source      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsImageBlock) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsImageBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsImageBlock to a
// BetaManagedAgentsImageBlockParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsImageBlockParam.Overrides()
func (r BetaManagedAgentsImageBlock) ToParam() BetaManagedAgentsImageBlockParam {
	return param.Override[BetaManagedAgentsImageBlockParam](json.RawMessage(r.RawJSON()))
}

// BetaManagedAgentsImageBlockSourceUnion contains all possible properties and
// values from [BetaManagedAgentsBase64ImageSource],
// [BetaManagedAgentsURLImageSource], [BetaManagedAgentsFileImageSource].
//
// Use the [BetaManagedAgentsImageBlockSourceUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsImageBlockSourceUnion struct {
	// This field is from variant [BetaManagedAgentsBase64ImageSource].
	Data string `json:"data"`
	// This field is from variant [BetaManagedAgentsBase64ImageSource].
	MediaType string `json:"media_type"`
	// Any of "base64", "url", "file".
	Type string `json:"type"`
	// This field is from variant [BetaManagedAgentsURLImageSource].
	URL string `json:"url"`
	// This field is from variant [BetaManagedAgentsFileImageSource].
	FileID string `json:"file_id"`
	JSON   struct {
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

// anyBetaManagedAgentsImageBlockSource is implemented by each variant of
// [BetaManagedAgentsImageBlockSourceUnion] to add type safety for the return type
// of [BetaManagedAgentsImageBlockSourceUnion.AsAny]
type anyBetaManagedAgentsImageBlockSource interface {
	implBetaManagedAgentsImageBlockSourceUnion()
}

func (BetaManagedAgentsBase64ImageSource) implBetaManagedAgentsImageBlockSourceUnion() {}
func (BetaManagedAgentsURLImageSource) implBetaManagedAgentsImageBlockSourceUnion()    {}
func (BetaManagedAgentsFileImageSource) implBetaManagedAgentsImageBlockSourceUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsImageBlockSourceUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsBase64ImageSource:
//	case anthropic.BetaManagedAgentsURLImageSource:
//	case anthropic.BetaManagedAgentsFileImageSource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsImageBlockSourceUnion) AsAny() anyBetaManagedAgentsImageBlockSource {
	switch u.Type {
	case "base64":
		return u.AsBase64()
	case "url":
		return u.AsURL()
	case "file":
		return u.AsFile()
	}
	return nil
}

func (u BetaManagedAgentsImageBlockSourceUnion) AsBase64() (v BetaManagedAgentsBase64ImageSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsImageBlockSourceUnion) AsURL() (v BetaManagedAgentsURLImageSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsImageBlockSourceUnion) AsFile() (v BetaManagedAgentsFileImageSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsImageBlockSourceUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsImageBlockSourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsImageBlockType string

const (
	BetaManagedAgentsImageBlockTypeImage BetaManagedAgentsImageBlockType = "image"
)

// Image content specified directly as base64 data or as a reference via a URL.
//
// The properties Source, Type are required.
type BetaManagedAgentsImageBlockParam struct {
	// Union type for image source variants.
	Source BetaManagedAgentsImageBlockSourceUnionParam `json:"source,omitzero" api:"required"`
	// Any of "image".
	Type BetaManagedAgentsImageBlockType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsImageBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsImageBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsImageBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsImageBlockSourceUnionParam struct {
	OfBase64 *BetaManagedAgentsBase64ImageSourceParam `json:",omitzero,inline"`
	OfURL    *BetaManagedAgentsURLImageSourceParam    `json:",omitzero,inline"`
	OfFile   *BetaManagedAgentsFileImageSourceParam   `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsImageBlockSourceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfBase64, u.OfURL, u.OfFile)
}
func (u *BetaManagedAgentsImageBlockSourceUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsImageBlockSourceUnionParam) asAny() any {
	if !param.IsOmitted(u.OfBase64) {
		return u.OfBase64
	} else if !param.IsOmitted(u.OfURL) {
		return u.OfURL
	} else if !param.IsOmitted(u.OfFile) {
		return u.OfFile
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsImageBlockSourceUnionParam) GetData() *string {
	if vt := u.OfBase64; vt != nil {
		return &vt.Data
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsImageBlockSourceUnionParam) GetMediaType() *string {
	if vt := u.OfBase64; vt != nil {
		return &vt.MediaType
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsImageBlockSourceUnionParam) GetURL() *string {
	if vt := u.OfURL; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsImageBlockSourceUnionParam) GetFileID() *string {
	if vt := u.OfFile; vt != nil {
		return &vt.FileID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsImageBlockSourceUnionParam) GetType() *string {
	if vt := u.OfBase64; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfURL; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFile; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsImageBlockSourceUnionParam](
		"type",
		apijson.Discriminator[BetaManagedAgentsBase64ImageSourceParam]("base64"),
		apijson.Discriminator[BetaManagedAgentsURLImageSourceParam]("url"),
		apijson.Discriminator[BetaManagedAgentsFileImageSourceParam]("file"),
	)
}

// Authentication to an MCP server failed.
type BetaManagedAgentsMCPAuthenticationFailedError struct {
	// Name of the MCP server that failed authentication.
	MCPServerName string `json:"mcp_server_name" api:"required"`
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// What the client should do next in response to this error.
	RetryStatus BetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion `json:"retry_status" api:"required"`
	// Any of "mcp_authentication_failed_error".
	Type BetaManagedAgentsMCPAuthenticationFailedErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MCPServerName respjson.Field
		Message       respjson.Field
		RetryStatus   respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMCPAuthenticationFailedError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsMCPAuthenticationFailedError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion contains all
// possible properties and values from [BetaManagedAgentsRetryStatusRetrying],
// [BetaManagedAgentsRetryStatusExhausted], [BetaManagedAgentsRetryStatusTerminal].
//
// Use the [BetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion struct {
	// Any of "retrying", "exhausted", "terminal".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsMCPAuthenticationFailedErrorRetryStatus is implemented by
// each variant of [BetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion]
// to add type safety for the return type of
// [BetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion.AsAny]
type anyBetaManagedAgentsMCPAuthenticationFailedErrorRetryStatus interface {
	implBetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion()
}

func (BetaManagedAgentsRetryStatusRetrying) implBetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion() {
}
func (BetaManagedAgentsRetryStatusExhausted) implBetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion() {
}
func (BetaManagedAgentsRetryStatusTerminal) implBetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsRetryStatusRetrying:
//	case anthropic.BetaManagedAgentsRetryStatusExhausted:
//	case anthropic.BetaManagedAgentsRetryStatusTerminal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion) AsAny() anyBetaManagedAgentsMCPAuthenticationFailedErrorRetryStatus {
	switch u.Type {
	case "retrying":
		return u.AsRetrying()
	case "exhausted":
		return u.AsExhausted()
	case "terminal":
		return u.AsTerminal()
	}
	return nil
}

func (u BetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion) AsRetrying() (v BetaManagedAgentsRetryStatusRetrying) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion) AsExhausted() (v BetaManagedAgentsRetryStatusExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion) AsTerminal() (v BetaManagedAgentsRetryStatusTerminal) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *BetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMCPAuthenticationFailedErrorType string

const (
	BetaManagedAgentsMCPAuthenticationFailedErrorTypeMCPAuthenticationFailedError BetaManagedAgentsMCPAuthenticationFailedErrorType = "mcp_authentication_failed_error"
)

// Failed to connect to an MCP server.
type BetaManagedAgentsMCPConnectionFailedError struct {
	// Name of the MCP server that failed to connect.
	MCPServerName string `json:"mcp_server_name" api:"required"`
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// What the client should do next in response to this error.
	RetryStatus BetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion `json:"retry_status" api:"required"`
	// Any of "mcp_connection_failed_error".
	Type BetaManagedAgentsMCPConnectionFailedErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MCPServerName respjson.Field
		Message       respjson.Field
		RetryStatus   respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMCPConnectionFailedError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsMCPConnectionFailedError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion contains all possible
// properties and values from [BetaManagedAgentsRetryStatusRetrying],
// [BetaManagedAgentsRetryStatusExhausted], [BetaManagedAgentsRetryStatusTerminal].
//
// Use the [BetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion struct {
	// Any of "retrying", "exhausted", "terminal".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsMCPConnectionFailedErrorRetryStatus is implemented by each
// variant of [BetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion] to add
// type safety for the return type of
// [BetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion.AsAny]
type anyBetaManagedAgentsMCPConnectionFailedErrorRetryStatus interface {
	implBetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion()
}

func (BetaManagedAgentsRetryStatusRetrying) implBetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion() {
}
func (BetaManagedAgentsRetryStatusExhausted) implBetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion() {
}
func (BetaManagedAgentsRetryStatusTerminal) implBetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsRetryStatusRetrying:
//	case anthropic.BetaManagedAgentsRetryStatusExhausted:
//	case anthropic.BetaManagedAgentsRetryStatusTerminal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion) AsAny() anyBetaManagedAgentsMCPConnectionFailedErrorRetryStatus {
	switch u.Type {
	case "retrying":
		return u.AsRetrying()
	case "exhausted":
		return u.AsExhausted()
	case "terminal":
		return u.AsTerminal()
	}
	return nil
}

func (u BetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion) AsRetrying() (v BetaManagedAgentsRetryStatusRetrying) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion) AsExhausted() (v BetaManagedAgentsRetryStatusExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion) AsTerminal() (v BetaManagedAgentsRetryStatusTerminal) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *BetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMCPConnectionFailedErrorType string

const (
	BetaManagedAgentsMCPConnectionFailedErrorTypeMCPConnectionFailedError BetaManagedAgentsMCPConnectionFailedErrorType = "mcp_connection_failed_error"
)

// The model is currently overloaded. Emitted after automatic retries are
// exhausted.
type BetaManagedAgentsModelOverloadedError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// What the client should do next in response to this error.
	RetryStatus BetaManagedAgentsModelOverloadedErrorRetryStatusUnion `json:"retry_status" api:"required"`
	// Any of "model_overloaded_error".
	Type BetaManagedAgentsModelOverloadedErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		RetryStatus respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsModelOverloadedError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsModelOverloadedError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsModelOverloadedErrorRetryStatusUnion contains all possible
// properties and values from [BetaManagedAgentsRetryStatusRetrying],
// [BetaManagedAgentsRetryStatusExhausted], [BetaManagedAgentsRetryStatusTerminal].
//
// Use the [BetaManagedAgentsModelOverloadedErrorRetryStatusUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsModelOverloadedErrorRetryStatusUnion struct {
	// Any of "retrying", "exhausted", "terminal".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsModelOverloadedErrorRetryStatus is implemented by each
// variant of [BetaManagedAgentsModelOverloadedErrorRetryStatusUnion] to add type
// safety for the return type of
// [BetaManagedAgentsModelOverloadedErrorRetryStatusUnion.AsAny]
type anyBetaManagedAgentsModelOverloadedErrorRetryStatus interface {
	implBetaManagedAgentsModelOverloadedErrorRetryStatusUnion()
}

func (BetaManagedAgentsRetryStatusRetrying) implBetaManagedAgentsModelOverloadedErrorRetryStatusUnion() {
}
func (BetaManagedAgentsRetryStatusExhausted) implBetaManagedAgentsModelOverloadedErrorRetryStatusUnion() {
}
func (BetaManagedAgentsRetryStatusTerminal) implBetaManagedAgentsModelOverloadedErrorRetryStatusUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsModelOverloadedErrorRetryStatusUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsRetryStatusRetrying:
//	case anthropic.BetaManagedAgentsRetryStatusExhausted:
//	case anthropic.BetaManagedAgentsRetryStatusTerminal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsModelOverloadedErrorRetryStatusUnion) AsAny() anyBetaManagedAgentsModelOverloadedErrorRetryStatus {
	switch u.Type {
	case "retrying":
		return u.AsRetrying()
	case "exhausted":
		return u.AsExhausted()
	case "terminal":
		return u.AsTerminal()
	}
	return nil
}

func (u BetaManagedAgentsModelOverloadedErrorRetryStatusUnion) AsRetrying() (v BetaManagedAgentsRetryStatusRetrying) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsModelOverloadedErrorRetryStatusUnion) AsExhausted() (v BetaManagedAgentsRetryStatusExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsModelOverloadedErrorRetryStatusUnion) AsTerminal() (v BetaManagedAgentsRetryStatusTerminal) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsModelOverloadedErrorRetryStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsModelOverloadedErrorRetryStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsModelOverloadedErrorType string

const (
	BetaManagedAgentsModelOverloadedErrorTypeModelOverloadedError BetaManagedAgentsModelOverloadedErrorType = "model_overloaded_error"
)

// The model request was rate-limited.
type BetaManagedAgentsModelRateLimitedError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// What the client should do next in response to this error.
	RetryStatus BetaManagedAgentsModelRateLimitedErrorRetryStatusUnion `json:"retry_status" api:"required"`
	// Any of "model_rate_limited_error".
	Type BetaManagedAgentsModelRateLimitedErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		RetryStatus respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsModelRateLimitedError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsModelRateLimitedError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsModelRateLimitedErrorRetryStatusUnion contains all possible
// properties and values from [BetaManagedAgentsRetryStatusRetrying],
// [BetaManagedAgentsRetryStatusExhausted], [BetaManagedAgentsRetryStatusTerminal].
//
// Use the [BetaManagedAgentsModelRateLimitedErrorRetryStatusUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsModelRateLimitedErrorRetryStatusUnion struct {
	// Any of "retrying", "exhausted", "terminal".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsModelRateLimitedErrorRetryStatus is implemented by each
// variant of [BetaManagedAgentsModelRateLimitedErrorRetryStatusUnion] to add type
// safety for the return type of
// [BetaManagedAgentsModelRateLimitedErrorRetryStatusUnion.AsAny]
type anyBetaManagedAgentsModelRateLimitedErrorRetryStatus interface {
	implBetaManagedAgentsModelRateLimitedErrorRetryStatusUnion()
}

func (BetaManagedAgentsRetryStatusRetrying) implBetaManagedAgentsModelRateLimitedErrorRetryStatusUnion() {
}
func (BetaManagedAgentsRetryStatusExhausted) implBetaManagedAgentsModelRateLimitedErrorRetryStatusUnion() {
}
func (BetaManagedAgentsRetryStatusTerminal) implBetaManagedAgentsModelRateLimitedErrorRetryStatusUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsModelRateLimitedErrorRetryStatusUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsRetryStatusRetrying:
//	case anthropic.BetaManagedAgentsRetryStatusExhausted:
//	case anthropic.BetaManagedAgentsRetryStatusTerminal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsModelRateLimitedErrorRetryStatusUnion) AsAny() anyBetaManagedAgentsModelRateLimitedErrorRetryStatus {
	switch u.Type {
	case "retrying":
		return u.AsRetrying()
	case "exhausted":
		return u.AsExhausted()
	case "terminal":
		return u.AsTerminal()
	}
	return nil
}

func (u BetaManagedAgentsModelRateLimitedErrorRetryStatusUnion) AsRetrying() (v BetaManagedAgentsRetryStatusRetrying) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsModelRateLimitedErrorRetryStatusUnion) AsExhausted() (v BetaManagedAgentsRetryStatusExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsModelRateLimitedErrorRetryStatusUnion) AsTerminal() (v BetaManagedAgentsRetryStatusTerminal) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsModelRateLimitedErrorRetryStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsModelRateLimitedErrorRetryStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsModelRateLimitedErrorType string

const (
	BetaManagedAgentsModelRateLimitedErrorTypeModelRateLimitedError BetaManagedAgentsModelRateLimitedErrorType = "model_rate_limited_error"
)

// A model request failed for a reason other than overload or rate-limiting.
type BetaManagedAgentsModelRequestFailedError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// What the client should do next in response to this error.
	RetryStatus BetaManagedAgentsModelRequestFailedErrorRetryStatusUnion `json:"retry_status" api:"required"`
	// Any of "model_request_failed_error".
	Type BetaManagedAgentsModelRequestFailedErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		RetryStatus respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsModelRequestFailedError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsModelRequestFailedError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsModelRequestFailedErrorRetryStatusUnion contains all possible
// properties and values from [BetaManagedAgentsRetryStatusRetrying],
// [BetaManagedAgentsRetryStatusExhausted], [BetaManagedAgentsRetryStatusTerminal].
//
// Use the [BetaManagedAgentsModelRequestFailedErrorRetryStatusUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsModelRequestFailedErrorRetryStatusUnion struct {
	// Any of "retrying", "exhausted", "terminal".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsModelRequestFailedErrorRetryStatus is implemented by each
// variant of [BetaManagedAgentsModelRequestFailedErrorRetryStatusUnion] to add
// type safety for the return type of
// [BetaManagedAgentsModelRequestFailedErrorRetryStatusUnion.AsAny]
type anyBetaManagedAgentsModelRequestFailedErrorRetryStatus interface {
	implBetaManagedAgentsModelRequestFailedErrorRetryStatusUnion()
}

func (BetaManagedAgentsRetryStatusRetrying) implBetaManagedAgentsModelRequestFailedErrorRetryStatusUnion() {
}
func (BetaManagedAgentsRetryStatusExhausted) implBetaManagedAgentsModelRequestFailedErrorRetryStatusUnion() {
}
func (BetaManagedAgentsRetryStatusTerminal) implBetaManagedAgentsModelRequestFailedErrorRetryStatusUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsModelRequestFailedErrorRetryStatusUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsRetryStatusRetrying:
//	case anthropic.BetaManagedAgentsRetryStatusExhausted:
//	case anthropic.BetaManagedAgentsRetryStatusTerminal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsModelRequestFailedErrorRetryStatusUnion) AsAny() anyBetaManagedAgentsModelRequestFailedErrorRetryStatus {
	switch u.Type {
	case "retrying":
		return u.AsRetrying()
	case "exhausted":
		return u.AsExhausted()
	case "terminal":
		return u.AsTerminal()
	}
	return nil
}

func (u BetaManagedAgentsModelRequestFailedErrorRetryStatusUnion) AsRetrying() (v BetaManagedAgentsRetryStatusRetrying) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsModelRequestFailedErrorRetryStatusUnion) AsExhausted() (v BetaManagedAgentsRetryStatusExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsModelRequestFailedErrorRetryStatusUnion) AsTerminal() (v BetaManagedAgentsRetryStatusTerminal) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsModelRequestFailedErrorRetryStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsModelRequestFailedErrorRetryStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsModelRequestFailedErrorType string

const (
	BetaManagedAgentsModelRequestFailedErrorTypeModelRequestFailedError BetaManagedAgentsModelRequestFailedErrorType = "model_request_failed_error"
)

// Plain text document content.
type BetaManagedAgentsPlainTextDocumentSource struct {
	// The plain text content.
	Data string `json:"data" api:"required"`
	// MIME type of the text content. Must be "text/plain".
	//
	// Any of "text/plain".
	MediaType BetaManagedAgentsPlainTextDocumentSourceMediaType `json:"media_type" api:"required"`
	// Any of "text".
	Type BetaManagedAgentsPlainTextDocumentSourceType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		MediaType   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsPlainTextDocumentSource) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsPlainTextDocumentSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsPlainTextDocumentSource to a
// BetaManagedAgentsPlainTextDocumentSourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsPlainTextDocumentSourceParam.Overrides()
func (r BetaManagedAgentsPlainTextDocumentSource) ToParam() BetaManagedAgentsPlainTextDocumentSourceParam {
	return param.Override[BetaManagedAgentsPlainTextDocumentSourceParam](json.RawMessage(r.RawJSON()))
}

// MIME type of the text content. Must be "text/plain".
type BetaManagedAgentsPlainTextDocumentSourceMediaType string

const (
	BetaManagedAgentsPlainTextDocumentSourceMediaTypeTextPlain BetaManagedAgentsPlainTextDocumentSourceMediaType = "text/plain"
)

type BetaManagedAgentsPlainTextDocumentSourceType string

const (
	BetaManagedAgentsPlainTextDocumentSourceTypeText BetaManagedAgentsPlainTextDocumentSourceType = "text"
)

// Plain text document content.
//
// The properties Data, MediaType, Type are required.
type BetaManagedAgentsPlainTextDocumentSourceParam struct {
	// The plain text content.
	Data string `json:"data" api:"required"`
	// MIME type of the text content. Must be "text/plain".
	//
	// Any of "text/plain".
	MediaType BetaManagedAgentsPlainTextDocumentSourceMediaType `json:"media_type,omitzero" api:"required"`
	// Any of "text".
	Type BetaManagedAgentsPlainTextDocumentSourceType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsPlainTextDocumentSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsPlainTextDocumentSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsPlainTextDocumentSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// This turn is dead; queued inputs are flushed and the session returns to idle.
// Client may send a new prompt.
type BetaManagedAgentsRetryStatusExhausted struct {
	// Any of "exhausted".
	Type BetaManagedAgentsRetryStatusExhaustedType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsRetryStatusExhausted) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsRetryStatusExhausted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsRetryStatusExhaustedType string

const (
	BetaManagedAgentsRetryStatusExhaustedTypeExhausted BetaManagedAgentsRetryStatusExhaustedType = "exhausted"
)

// The server is retrying automatically. Client should wait; the same error type
// may fire again as retrying, then once as exhausted when the retry budget runs
// out.
type BetaManagedAgentsRetryStatusRetrying struct {
	// Any of "retrying".
	Type BetaManagedAgentsRetryStatusRetryingType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsRetryStatusRetrying) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsRetryStatusRetrying) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsRetryStatusRetryingType string

const (
	BetaManagedAgentsRetryStatusRetryingTypeRetrying BetaManagedAgentsRetryStatusRetryingType = "retrying"
)

// The session encountered a terminal error and will transition to `terminated`
// state.
type BetaManagedAgentsRetryStatusTerminal struct {
	// Any of "terminal".
	Type BetaManagedAgentsRetryStatusTerminalType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsRetryStatusTerminal) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsRetryStatusTerminal) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsRetryStatusTerminalType string

const (
	BetaManagedAgentsRetryStatusTerminalTypeTerminal BetaManagedAgentsRetryStatusTerminalType = "terminal"
)

// Events that were successfully sent to the session.
type BetaManagedAgentsSendSessionEvents struct {
	// Sent events
	Data []BetaManagedAgentsSendSessionEventsDataUnion `json:"data"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSendSessionEvents) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSendSessionEvents) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSendSessionEventsDataUnion contains all possible properties and
// values from [BetaManagedAgentsUserMessageEvent],
// [BetaManagedAgentsUserInterruptEvent],
// [BetaManagedAgentsUserToolConfirmationEvent],
// [BetaManagedAgentsUserCustomToolResultEvent].
//
// Use the [BetaManagedAgentsSendSessionEventsDataUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsSendSessionEventsDataUnion struct {
	ID string `json:"id"`
	// This field is a union of [[]BetaManagedAgentsUserMessageEventContentUnion],
	// [[]BetaManagedAgentsUserCustomToolResultEventContentUnion]
	Content BetaManagedAgentsSendSessionEventsDataUnionContent `json:"content"`
	// Any of "user.message", "user.interrupt", "user.tool_confirmation",
	// "user.custom_tool_result".
	Type        string    `json:"type"`
	ProcessedAt time.Time `json:"processed_at"`
	// This field is from variant [BetaManagedAgentsUserToolConfirmationEvent].
	Result BetaManagedAgentsUserToolConfirmationEventResult `json:"result"`
	// This field is from variant [BetaManagedAgentsUserToolConfirmationEvent].
	ToolUseID string `json:"tool_use_id"`
	// This field is from variant [BetaManagedAgentsUserToolConfirmationEvent].
	DenyMessage string `json:"deny_message"`
	// This field is from variant [BetaManagedAgentsUserCustomToolResultEvent].
	CustomToolUseID string `json:"custom_tool_use_id"`
	// This field is from variant [BetaManagedAgentsUserCustomToolResultEvent].
	IsError bool `json:"is_error"`
	JSON    struct {
		ID              respjson.Field
		Content         respjson.Field
		Type            respjson.Field
		ProcessedAt     respjson.Field
		Result          respjson.Field
		ToolUseID       respjson.Field
		DenyMessage     respjson.Field
		CustomToolUseID respjson.Field
		IsError         respjson.Field
		raw             string
	} `json:"-"`
}

// anyBetaManagedAgentsSendSessionEventsData is implemented by each variant of
// [BetaManagedAgentsSendSessionEventsDataUnion] to add type safety for the return
// type of [BetaManagedAgentsSendSessionEventsDataUnion.AsAny]
type anyBetaManagedAgentsSendSessionEventsData interface {
	implBetaManagedAgentsSendSessionEventsDataUnion()
}

func (BetaManagedAgentsUserMessageEvent) implBetaManagedAgentsSendSessionEventsDataUnion()          {}
func (BetaManagedAgentsUserInterruptEvent) implBetaManagedAgentsSendSessionEventsDataUnion()        {}
func (BetaManagedAgentsUserToolConfirmationEvent) implBetaManagedAgentsSendSessionEventsDataUnion() {}
func (BetaManagedAgentsUserCustomToolResultEvent) implBetaManagedAgentsSendSessionEventsDataUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsSendSessionEventsDataUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsUserMessageEvent:
//	case anthropic.BetaManagedAgentsUserInterruptEvent:
//	case anthropic.BetaManagedAgentsUserToolConfirmationEvent:
//	case anthropic.BetaManagedAgentsUserCustomToolResultEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsSendSessionEventsDataUnion) AsAny() anyBetaManagedAgentsSendSessionEventsData {
	switch u.Type {
	case "user.message":
		return u.AsUserMessage()
	case "user.interrupt":
		return u.AsUserInterrupt()
	case "user.tool_confirmation":
		return u.AsUserToolConfirmation()
	case "user.custom_tool_result":
		return u.AsUserCustomToolResult()
	}
	return nil
}

func (u BetaManagedAgentsSendSessionEventsDataUnion) AsUserMessage() (v BetaManagedAgentsUserMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSendSessionEventsDataUnion) AsUserInterrupt() (v BetaManagedAgentsUserInterruptEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSendSessionEventsDataUnion) AsUserToolConfirmation() (v BetaManagedAgentsUserToolConfirmationEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSendSessionEventsDataUnion) AsUserCustomToolResult() (v BetaManagedAgentsUserCustomToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsSendSessionEventsDataUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsSendSessionEventsDataUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSendSessionEventsDataUnionContent is an implicit subunion of
// [BetaManagedAgentsSendSessionEventsDataUnion].
// BetaManagedAgentsSendSessionEventsDataUnionContent provides convenient access to
// the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsSendSessionEventsDataUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfBetaManagedAgentsUserMessageEventContentArray
// OfBetaManagedAgentsUserCustomToolResultEventContentArray]
type BetaManagedAgentsSendSessionEventsDataUnionContent struct {
	// This field will be present if the value is a
	// [[]BetaManagedAgentsUserMessageEventContentUnion] instead of an object.
	OfBetaManagedAgentsUserMessageEventContentArray []BetaManagedAgentsUserMessageEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]BetaManagedAgentsUserCustomToolResultEventContentUnion] instead of an object.
	OfBetaManagedAgentsUserCustomToolResultEventContentArray []BetaManagedAgentsUserCustomToolResultEventContentUnion `json:",inline"`
	JSON                                                     struct {
		OfBetaManagedAgentsUserMessageEventContentArray          respjson.Field
		OfBetaManagedAgentsUserCustomToolResultEventContentArray respjson.Field
		raw                                                      string
	} `json:"-"`
}

func (r *BetaManagedAgentsSendSessionEventsDataUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a session has been deleted. Terminates any active event stream — no
// further events will be emitted for this session.
type BetaManagedAgentsSessionDeletedEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "session.deleted".
	Type BetaManagedAgentsSessionDeletedEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionDeletedEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionDeletedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSessionDeletedEventType string

const (
	BetaManagedAgentsSessionDeletedEventTypeSessionDeleted BetaManagedAgentsSessionDeletedEventType = "session.deleted"
)

// The agent completed its turn naturally and is ready for the next user message.
type BetaManagedAgentsSessionEndTurn struct {
	// Any of "end_turn".
	Type BetaManagedAgentsSessionEndTurnType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionEndTurn) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionEndTurn) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSessionEndTurnType string

const (
	BetaManagedAgentsSessionEndTurnTypeEndTurn BetaManagedAgentsSessionEndTurnType = "end_turn"
)

// An error event indicating a problem occurred during session execution.
type BetaManagedAgentsSessionErrorEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// An unknown or unexpected error occurred during session execution. A fallback
	// variant; clients that don't recognize a new error code can match on
	// `retry_status` and `message` alone.
	Error BetaManagedAgentsSessionErrorEventErrorUnion `json:"error" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "session.error".
	Type BetaManagedAgentsSessionErrorEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Error       respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionErrorEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionErrorEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSessionErrorEventErrorUnion contains all possible properties
// and values from [BetaManagedAgentsUnknownError],
// [BetaManagedAgentsModelOverloadedError],
// [BetaManagedAgentsModelRateLimitedError],
// [BetaManagedAgentsModelRequestFailedError],
// [BetaManagedAgentsMCPConnectionFailedError],
// [BetaManagedAgentsMCPAuthenticationFailedError],
// [BetaManagedAgentsBillingError].
//
// Use the [BetaManagedAgentsSessionErrorEventErrorUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsSessionErrorEventErrorUnion struct {
	Message string `json:"message"`
	// This field is a union of [BetaManagedAgentsUnknownErrorRetryStatusUnion],
	// [BetaManagedAgentsModelOverloadedErrorRetryStatusUnion],
	// [BetaManagedAgentsModelRateLimitedErrorRetryStatusUnion],
	// [BetaManagedAgentsModelRequestFailedErrorRetryStatusUnion],
	// [BetaManagedAgentsMCPConnectionFailedErrorRetryStatusUnion],
	// [BetaManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion],
	// [BetaManagedAgentsBillingErrorRetryStatusUnion]
	RetryStatus BetaManagedAgentsSessionErrorEventErrorUnionRetryStatus `json:"retry_status"`
	// Any of "unknown_error", "model_overloaded_error", "model_rate_limited_error",
	// "model_request_failed_error", "mcp_connection_failed_error",
	// "mcp_authentication_failed_error", "billing_error".
	Type          string `json:"type"`
	MCPServerName string `json:"mcp_server_name"`
	JSON          struct {
		Message       respjson.Field
		RetryStatus   respjson.Field
		Type          respjson.Field
		MCPServerName respjson.Field
		raw           string
	} `json:"-"`
}

// anyBetaManagedAgentsSessionErrorEventError is implemented by each variant of
// [BetaManagedAgentsSessionErrorEventErrorUnion] to add type safety for the return
// type of [BetaManagedAgentsSessionErrorEventErrorUnion.AsAny]
type anyBetaManagedAgentsSessionErrorEventError interface {
	implBetaManagedAgentsSessionErrorEventErrorUnion()
}

func (BetaManagedAgentsUnknownError) implBetaManagedAgentsSessionErrorEventErrorUnion()             {}
func (BetaManagedAgentsModelOverloadedError) implBetaManagedAgentsSessionErrorEventErrorUnion()     {}
func (BetaManagedAgentsModelRateLimitedError) implBetaManagedAgentsSessionErrorEventErrorUnion()    {}
func (BetaManagedAgentsModelRequestFailedError) implBetaManagedAgentsSessionErrorEventErrorUnion()  {}
func (BetaManagedAgentsMCPConnectionFailedError) implBetaManagedAgentsSessionErrorEventErrorUnion() {}
func (BetaManagedAgentsMCPAuthenticationFailedError) implBetaManagedAgentsSessionErrorEventErrorUnion() {
}
func (BetaManagedAgentsBillingError) implBetaManagedAgentsSessionErrorEventErrorUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsSessionErrorEventErrorUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsUnknownError:
//	case anthropic.BetaManagedAgentsModelOverloadedError:
//	case anthropic.BetaManagedAgentsModelRateLimitedError:
//	case anthropic.BetaManagedAgentsModelRequestFailedError:
//	case anthropic.BetaManagedAgentsMCPConnectionFailedError:
//	case anthropic.BetaManagedAgentsMCPAuthenticationFailedError:
//	case anthropic.BetaManagedAgentsBillingError:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsSessionErrorEventErrorUnion) AsAny() anyBetaManagedAgentsSessionErrorEventError {
	switch u.Type {
	case "unknown_error":
		return u.AsUnknownError()
	case "model_overloaded_error":
		return u.AsModelOverloadedError()
	case "model_rate_limited_error":
		return u.AsModelRateLimitedError()
	case "model_request_failed_error":
		return u.AsModelRequestFailedError()
	case "mcp_connection_failed_error":
		return u.AsMCPConnectionFailedError()
	case "mcp_authentication_failed_error":
		return u.AsMCPAuthenticationFailedError()
	case "billing_error":
		return u.AsBillingError()
	}
	return nil
}

func (u BetaManagedAgentsSessionErrorEventErrorUnion) AsUnknownError() (v BetaManagedAgentsUnknownError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionErrorEventErrorUnion) AsModelOverloadedError() (v BetaManagedAgentsModelOverloadedError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionErrorEventErrorUnion) AsModelRateLimitedError() (v BetaManagedAgentsModelRateLimitedError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionErrorEventErrorUnion) AsModelRequestFailedError() (v BetaManagedAgentsModelRequestFailedError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionErrorEventErrorUnion) AsMCPConnectionFailedError() (v BetaManagedAgentsMCPConnectionFailedError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionErrorEventErrorUnion) AsMCPAuthenticationFailedError() (v BetaManagedAgentsMCPAuthenticationFailedError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionErrorEventErrorUnion) AsBillingError() (v BetaManagedAgentsBillingError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsSessionErrorEventErrorUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsSessionErrorEventErrorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSessionErrorEventErrorUnionRetryStatus is an implicit subunion
// of [BetaManagedAgentsSessionErrorEventErrorUnion].
// BetaManagedAgentsSessionErrorEventErrorUnionRetryStatus provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsSessionErrorEventErrorUnion].
type BetaManagedAgentsSessionErrorEventErrorUnionRetryStatus struct {
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

func (r *BetaManagedAgentsSessionErrorEventErrorUnionRetryStatus) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSessionErrorEventType string

const (
	BetaManagedAgentsSessionErrorEventTypeSessionError BetaManagedAgentsSessionErrorEventType = "session.error"
)

// BetaManagedAgentsSessionEventUnion contains all possible properties and values
// from [BetaManagedAgentsUserMessageEvent], [BetaManagedAgentsUserInterruptEvent],
// [BetaManagedAgentsUserToolConfirmationEvent],
// [BetaManagedAgentsUserCustomToolResultEvent],
// [BetaManagedAgentsAgentCustomToolUseEvent],
// [BetaManagedAgentsAgentMessageEvent], [BetaManagedAgentsAgentThinkingEvent],
// [BetaManagedAgentsAgentMCPToolUseEvent],
// [BetaManagedAgentsAgentMCPToolResultEvent],
// [BetaManagedAgentsAgentToolUseEvent], [BetaManagedAgentsAgentToolResultEvent],
// [BetaManagedAgentsAgentThreadContextCompactedEvent],
// [BetaManagedAgentsSessionErrorEvent],
// [BetaManagedAgentsSessionStatusRescheduledEvent],
// [BetaManagedAgentsSessionStatusRunningEvent],
// [BetaManagedAgentsSessionStatusIdleEvent],
// [BetaManagedAgentsSessionStatusTerminatedEvent],
// [BetaManagedAgentsSpanModelRequestStartEvent],
// [BetaManagedAgentsSpanModelRequestEndEvent],
// [BetaManagedAgentsSessionDeletedEvent].
//
// Use the [BetaManagedAgentsSessionEventUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsSessionEventUnion struct {
	ID string `json:"id"`
	// This field is a union of [[]BetaManagedAgentsUserMessageEventContentUnion],
	// [[]BetaManagedAgentsUserCustomToolResultEventContentUnion],
	// [[]BetaManagedAgentsTextBlock],
	// [[]BetaManagedAgentsAgentMCPToolResultEventContentUnion],
	// [[]BetaManagedAgentsAgentToolResultEventContentUnion]
	Content BetaManagedAgentsSessionEventUnionContent `json:"content"`
	// Any of "user.message", "user.interrupt", "user.tool_confirmation",
	// "user.custom_tool_result", "agent.custom_tool_use", "agent.message",
	// "agent.thinking", "agent.mcp_tool_use", "agent.mcp_tool_result",
	// "agent.tool_use", "agent.tool_result", "agent.thread_context_compacted",
	// "session.error", "session.status_rescheduled", "session.status_running",
	// "session.status_idle", "session.status_terminated", "span.model_request_start",
	// "span.model_request_end", "session.deleted".
	Type        string    `json:"type"`
	ProcessedAt time.Time `json:"processed_at"`
	// This field is from variant [BetaManagedAgentsUserToolConfirmationEvent].
	Result    BetaManagedAgentsUserToolConfirmationEventResult `json:"result"`
	ToolUseID string                                           `json:"tool_use_id"`
	// This field is from variant [BetaManagedAgentsUserToolConfirmationEvent].
	DenyMessage string `json:"deny_message"`
	// This field is from variant [BetaManagedAgentsUserCustomToolResultEvent].
	CustomToolUseID string `json:"custom_tool_use_id"`
	IsError         bool   `json:"is_error"`
	Input           any    `json:"input"`
	Name            string `json:"name"`
	// This field is from variant [BetaManagedAgentsAgentMCPToolUseEvent].
	MCPServerName       string `json:"mcp_server_name"`
	EvaluatedPermission string `json:"evaluated_permission"`
	// This field is from variant [BetaManagedAgentsAgentMCPToolResultEvent].
	MCPToolUseID string `json:"mcp_tool_use_id"`
	// This field is from variant [BetaManagedAgentsSessionErrorEvent].
	Error BetaManagedAgentsSessionErrorEventErrorUnion `json:"error"`
	// This field is from variant [BetaManagedAgentsSessionStatusIdleEvent].
	StopReason BetaManagedAgentsSessionStatusIdleEventStopReasonUnion `json:"stop_reason"`
	// This field is from variant [BetaManagedAgentsSpanModelRequestEndEvent].
	ModelRequestStartID string `json:"model_request_start_id"`
	// This field is from variant [BetaManagedAgentsSpanModelRequestEndEvent].
	ModelUsage BetaManagedAgentsSpanModelUsage `json:"model_usage"`
	JSON       struct {
		ID                  respjson.Field
		Content             respjson.Field
		Type                respjson.Field
		ProcessedAt         respjson.Field
		Result              respjson.Field
		ToolUseID           respjson.Field
		DenyMessage         respjson.Field
		CustomToolUseID     respjson.Field
		IsError             respjson.Field
		Input               respjson.Field
		Name                respjson.Field
		MCPServerName       respjson.Field
		EvaluatedPermission respjson.Field
		MCPToolUseID        respjson.Field
		Error               respjson.Field
		StopReason          respjson.Field
		ModelRequestStartID respjson.Field
		ModelUsage          respjson.Field
		raw                 string
	} `json:"-"`
}

// anyBetaManagedAgentsSessionEvent is implemented by each variant of
// [BetaManagedAgentsSessionEventUnion] to add type safety for the return type of
// [BetaManagedAgentsSessionEventUnion.AsAny]
type anyBetaManagedAgentsSessionEvent interface {
	implBetaManagedAgentsSessionEventUnion()
}

func (BetaManagedAgentsUserMessageEvent) implBetaManagedAgentsSessionEventUnion()                 {}
func (BetaManagedAgentsUserInterruptEvent) implBetaManagedAgentsSessionEventUnion()               {}
func (BetaManagedAgentsUserToolConfirmationEvent) implBetaManagedAgentsSessionEventUnion()        {}
func (BetaManagedAgentsUserCustomToolResultEvent) implBetaManagedAgentsSessionEventUnion()        {}
func (BetaManagedAgentsAgentCustomToolUseEvent) implBetaManagedAgentsSessionEventUnion()          {}
func (BetaManagedAgentsAgentMessageEvent) implBetaManagedAgentsSessionEventUnion()                {}
func (BetaManagedAgentsAgentThinkingEvent) implBetaManagedAgentsSessionEventUnion()               {}
func (BetaManagedAgentsAgentMCPToolUseEvent) implBetaManagedAgentsSessionEventUnion()             {}
func (BetaManagedAgentsAgentMCPToolResultEvent) implBetaManagedAgentsSessionEventUnion()          {}
func (BetaManagedAgentsAgentToolUseEvent) implBetaManagedAgentsSessionEventUnion()                {}
func (BetaManagedAgentsAgentToolResultEvent) implBetaManagedAgentsSessionEventUnion()             {}
func (BetaManagedAgentsAgentThreadContextCompactedEvent) implBetaManagedAgentsSessionEventUnion() {}
func (BetaManagedAgentsSessionErrorEvent) implBetaManagedAgentsSessionEventUnion()                {}
func (BetaManagedAgentsSessionStatusRescheduledEvent) implBetaManagedAgentsSessionEventUnion()    {}
func (BetaManagedAgentsSessionStatusRunningEvent) implBetaManagedAgentsSessionEventUnion()        {}
func (BetaManagedAgentsSessionStatusIdleEvent) implBetaManagedAgentsSessionEventUnion()           {}
func (BetaManagedAgentsSessionStatusTerminatedEvent) implBetaManagedAgentsSessionEventUnion()     {}
func (BetaManagedAgentsSpanModelRequestStartEvent) implBetaManagedAgentsSessionEventUnion()       {}
func (BetaManagedAgentsSpanModelRequestEndEvent) implBetaManagedAgentsSessionEventUnion()         {}
func (BetaManagedAgentsSessionDeletedEvent) implBetaManagedAgentsSessionEventUnion()              {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsSessionEventUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsUserMessageEvent:
//	case anthropic.BetaManagedAgentsUserInterruptEvent:
//	case anthropic.BetaManagedAgentsUserToolConfirmationEvent:
//	case anthropic.BetaManagedAgentsUserCustomToolResultEvent:
//	case anthropic.BetaManagedAgentsAgentCustomToolUseEvent:
//	case anthropic.BetaManagedAgentsAgentMessageEvent:
//	case anthropic.BetaManagedAgentsAgentThinkingEvent:
//	case anthropic.BetaManagedAgentsAgentMCPToolUseEvent:
//	case anthropic.BetaManagedAgentsAgentMCPToolResultEvent:
//	case anthropic.BetaManagedAgentsAgentToolUseEvent:
//	case anthropic.BetaManagedAgentsAgentToolResultEvent:
//	case anthropic.BetaManagedAgentsAgentThreadContextCompactedEvent:
//	case anthropic.BetaManagedAgentsSessionErrorEvent:
//	case anthropic.BetaManagedAgentsSessionStatusRescheduledEvent:
//	case anthropic.BetaManagedAgentsSessionStatusRunningEvent:
//	case anthropic.BetaManagedAgentsSessionStatusIdleEvent:
//	case anthropic.BetaManagedAgentsSessionStatusTerminatedEvent:
//	case anthropic.BetaManagedAgentsSpanModelRequestStartEvent:
//	case anthropic.BetaManagedAgentsSpanModelRequestEndEvent:
//	case anthropic.BetaManagedAgentsSessionDeletedEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsSessionEventUnion) AsAny() anyBetaManagedAgentsSessionEvent {
	switch u.Type {
	case "user.message":
		return u.AsUserMessage()
	case "user.interrupt":
		return u.AsUserInterrupt()
	case "user.tool_confirmation":
		return u.AsUserToolConfirmation()
	case "user.custom_tool_result":
		return u.AsUserCustomToolResult()
	case "agent.custom_tool_use":
		return u.AsAgentCustomToolUse()
	case "agent.message":
		return u.AsAgentMessage()
	case "agent.thinking":
		return u.AsAgentThinking()
	case "agent.mcp_tool_use":
		return u.AsAgentMCPToolUse()
	case "agent.mcp_tool_result":
		return u.AsAgentMCPToolResult()
	case "agent.tool_use":
		return u.AsAgentToolUse()
	case "agent.tool_result":
		return u.AsAgentToolResult()
	case "agent.thread_context_compacted":
		return u.AsAgentThreadContextCompacted()
	case "session.error":
		return u.AsSessionError()
	case "session.status_rescheduled":
		return u.AsSessionStatusRescheduled()
	case "session.status_running":
		return u.AsSessionStatusRunning()
	case "session.status_idle":
		return u.AsSessionStatusIdle()
	case "session.status_terminated":
		return u.AsSessionStatusTerminated()
	case "span.model_request_start":
		return u.AsSpanModelRequestStart()
	case "span.model_request_end":
		return u.AsSpanModelRequestEnd()
	case "session.deleted":
		return u.AsSessionDeleted()
	}
	return nil
}

func (u BetaManagedAgentsSessionEventUnion) AsUserMessage() (v BetaManagedAgentsUserMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsUserInterrupt() (v BetaManagedAgentsUserInterruptEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsUserToolConfirmation() (v BetaManagedAgentsUserToolConfirmationEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsUserCustomToolResult() (v BetaManagedAgentsUserCustomToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsAgentCustomToolUse() (v BetaManagedAgentsAgentCustomToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsAgentMessage() (v BetaManagedAgentsAgentMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsAgentThinking() (v BetaManagedAgentsAgentThinkingEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsAgentMCPToolUse() (v BetaManagedAgentsAgentMCPToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsAgentMCPToolResult() (v BetaManagedAgentsAgentMCPToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsAgentToolUse() (v BetaManagedAgentsAgentToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsAgentToolResult() (v BetaManagedAgentsAgentToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsAgentThreadContextCompacted() (v BetaManagedAgentsAgentThreadContextCompactedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsSessionError() (v BetaManagedAgentsSessionErrorEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsSessionStatusRescheduled() (v BetaManagedAgentsSessionStatusRescheduledEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsSessionStatusRunning() (v BetaManagedAgentsSessionStatusRunningEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsSessionStatusIdle() (v BetaManagedAgentsSessionStatusIdleEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsSessionStatusTerminated() (v BetaManagedAgentsSessionStatusTerminatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsSpanModelRequestStart() (v BetaManagedAgentsSpanModelRequestStartEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsSpanModelRequestEnd() (v BetaManagedAgentsSpanModelRequestEndEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionEventUnion) AsSessionDeleted() (v BetaManagedAgentsSessionDeletedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsSessionEventUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsSessionEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSessionEventUnionContent is an implicit subunion of
// [BetaManagedAgentsSessionEventUnion]. BetaManagedAgentsSessionEventUnionContent
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsSessionEventUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfBetaManagedAgentsUserMessageEventContentArray
// OfBetaManagedAgentsUserCustomToolResultEventContentArray
// OfBetaManagedAgentsTextBlockArray
// OfBetaManagedAgentsAgentMCPToolResultEventContentArray
// OfBetaManagedAgentsAgentToolResultEventContentArray]
type BetaManagedAgentsSessionEventUnionContent struct {
	// This field will be present if the value is a
	// [[]BetaManagedAgentsUserMessageEventContentUnion] instead of an object.
	OfBetaManagedAgentsUserMessageEventContentArray []BetaManagedAgentsUserMessageEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]BetaManagedAgentsUserCustomToolResultEventContentUnion] instead of an object.
	OfBetaManagedAgentsUserCustomToolResultEventContentArray []BetaManagedAgentsUserCustomToolResultEventContentUnion `json:",inline"`
	// This field will be present if the value is a [[]BetaManagedAgentsTextBlock]
	// instead of an object.
	OfBetaManagedAgentsTextBlockArray []BetaManagedAgentsTextBlock `json:",inline"`
	// This field will be present if the value is a
	// [[]BetaManagedAgentsAgentMCPToolResultEventContentUnion] instead of an object.
	OfBetaManagedAgentsAgentMCPToolResultEventContentArray []BetaManagedAgentsAgentMCPToolResultEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]BetaManagedAgentsAgentToolResultEventContentUnion] instead of an object.
	OfBetaManagedAgentsAgentToolResultEventContentArray []BetaManagedAgentsAgentToolResultEventContentUnion `json:",inline"`
	JSON                                                struct {
		OfBetaManagedAgentsUserMessageEventContentArray          respjson.Field
		OfBetaManagedAgentsUserCustomToolResultEventContentArray respjson.Field
		OfBetaManagedAgentsTextBlockArray                        respjson.Field
		OfBetaManagedAgentsAgentMCPToolResultEventContentArray   respjson.Field
		OfBetaManagedAgentsAgentToolResultEventContentArray      respjson.Field
		raw                                                      string
	} `json:"-"`
}

func (r *BetaManagedAgentsSessionEventUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The agent is idle waiting on one or more blocking user-input events (tool
// confirmation, custom tool result, etc.). Resolving all of them transitions the
// session back to running.
type BetaManagedAgentsSessionRequiresAction struct {
	// The ids of events the agent is blocked on. Resolving fewer than all re-emits
	// `session.status_idle` with the remainder.
	EventIDs []string `json:"event_ids" api:"required"`
	// Any of "requires_action".
	Type BetaManagedAgentsSessionRequiresActionType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventIDs    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionRequiresAction) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionRequiresAction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSessionRequiresActionType string

const (
	BetaManagedAgentsSessionRequiresActionTypeRequiresAction BetaManagedAgentsSessionRequiresActionType = "requires_action"
)

// The turn ended because the retry budget was exhausted (`max_iterations` hit or
// an error escalated to `retry_status: 'exhausted'`).
type BetaManagedAgentsSessionRetriesExhausted struct {
	// Any of "retries_exhausted".
	Type BetaManagedAgentsSessionRetriesExhaustedType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionRetriesExhausted) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionRetriesExhausted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSessionRetriesExhaustedType string

const (
	BetaManagedAgentsSessionRetriesExhaustedTypeRetriesExhausted BetaManagedAgentsSessionRetriesExhaustedType = "retries_exhausted"
)

// Indicates the agent has paused and is awaiting user input.
type BetaManagedAgentsSessionStatusIdleEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// The agent completed its turn naturally and is ready for the next user message.
	StopReason BetaManagedAgentsSessionStatusIdleEventStopReasonUnion `json:"stop_reason" api:"required"`
	// Any of "session.status_idle".
	Type BetaManagedAgentsSessionStatusIdleEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		StopReason  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionStatusIdleEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionStatusIdleEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSessionStatusIdleEventStopReasonUnion contains all possible
// properties and values from [BetaManagedAgentsSessionEndTurn],
// [BetaManagedAgentsSessionRequiresAction],
// [BetaManagedAgentsSessionRetriesExhausted].
//
// Use the [BetaManagedAgentsSessionStatusIdleEventStopReasonUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsSessionStatusIdleEventStopReasonUnion struct {
	// Any of "end_turn", "requires_action", "retries_exhausted".
	Type string `json:"type"`
	// This field is from variant [BetaManagedAgentsSessionRequiresAction].
	EventIDs []string `json:"event_ids"`
	JSON     struct {
		Type     respjson.Field
		EventIDs respjson.Field
		raw      string
	} `json:"-"`
}

// anyBetaManagedAgentsSessionStatusIdleEventStopReason is implemented by each
// variant of [BetaManagedAgentsSessionStatusIdleEventStopReasonUnion] to add type
// safety for the return type of
// [BetaManagedAgentsSessionStatusIdleEventStopReasonUnion.AsAny]
type anyBetaManagedAgentsSessionStatusIdleEventStopReason interface {
	implBetaManagedAgentsSessionStatusIdleEventStopReasonUnion()
}

func (BetaManagedAgentsSessionEndTurn) implBetaManagedAgentsSessionStatusIdleEventStopReasonUnion() {}
func (BetaManagedAgentsSessionRequiresAction) implBetaManagedAgentsSessionStatusIdleEventStopReasonUnion() {
}
func (BetaManagedAgentsSessionRetriesExhausted) implBetaManagedAgentsSessionStatusIdleEventStopReasonUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsSessionStatusIdleEventStopReasonUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsSessionEndTurn:
//	case anthropic.BetaManagedAgentsSessionRequiresAction:
//	case anthropic.BetaManagedAgentsSessionRetriesExhausted:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsSessionStatusIdleEventStopReasonUnion) AsAny() anyBetaManagedAgentsSessionStatusIdleEventStopReason {
	switch u.Type {
	case "end_turn":
		return u.AsEndTurn()
	case "requires_action":
		return u.AsRequiresAction()
	case "retries_exhausted":
		return u.AsRetriesExhausted()
	}
	return nil
}

func (u BetaManagedAgentsSessionStatusIdleEventStopReasonUnion) AsEndTurn() (v BetaManagedAgentsSessionEndTurn) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionStatusIdleEventStopReasonUnion) AsRequiresAction() (v BetaManagedAgentsSessionRequiresAction) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionStatusIdleEventStopReasonUnion) AsRetriesExhausted() (v BetaManagedAgentsSessionRetriesExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsSessionStatusIdleEventStopReasonUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsSessionStatusIdleEventStopReasonUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSessionStatusIdleEventType string

const (
	BetaManagedAgentsSessionStatusIdleEventTypeSessionStatusIdle BetaManagedAgentsSessionStatusIdleEventType = "session.status_idle"
)

// Indicates the session is recovering from an error state and is rescheduled for
// execution.
type BetaManagedAgentsSessionStatusRescheduledEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "session.status_rescheduled".
	Type BetaManagedAgentsSessionStatusRescheduledEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionStatusRescheduledEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionStatusRescheduledEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSessionStatusRescheduledEventType string

const (
	BetaManagedAgentsSessionStatusRescheduledEventTypeSessionStatusRescheduled BetaManagedAgentsSessionStatusRescheduledEventType = "session.status_rescheduled"
)

// Indicates the session is actively running and the agent is working.
type BetaManagedAgentsSessionStatusRunningEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "session.status_running".
	Type BetaManagedAgentsSessionStatusRunningEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionStatusRunningEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionStatusRunningEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSessionStatusRunningEventType string

const (
	BetaManagedAgentsSessionStatusRunningEventTypeSessionStatusRunning BetaManagedAgentsSessionStatusRunningEventType = "session.status_running"
)

// Indicates the session has terminated, either due to an error or completion.
type BetaManagedAgentsSessionStatusTerminatedEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "session.status_terminated".
	Type BetaManagedAgentsSessionStatusTerminatedEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionStatusTerminatedEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionStatusTerminatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSessionStatusTerminatedEventType string

const (
	BetaManagedAgentsSessionStatusTerminatedEventTypeSessionStatusTerminated BetaManagedAgentsSessionStatusTerminatedEventType = "session.status_terminated"
)

// Emitted when a model request completes.
type BetaManagedAgentsSpanModelRequestEndEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Whether the model request resulted in an error.
	IsError bool `json:"is_error" api:"required"`
	// The id of the corresponding `span.model_request_start` event.
	ModelRequestStartID string `json:"model_request_start_id" api:"required"`
	// Token usage for a single model request.
	ModelUsage BetaManagedAgentsSpanModelUsage `json:"model_usage" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "span.model_request_end".
	Type BetaManagedAgentsSpanModelRequestEndEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                  respjson.Field
		IsError             respjson.Field
		ModelRequestStartID respjson.Field
		ModelUsage          respjson.Field
		ProcessedAt         respjson.Field
		Type                respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSpanModelRequestEndEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSpanModelRequestEndEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSpanModelRequestEndEventType string

const (
	BetaManagedAgentsSpanModelRequestEndEventTypeSpanModelRequestEnd BetaManagedAgentsSpanModelRequestEndEventType = "span.model_request_end"
)

// Emitted when a model request is initiated by the agent.
type BetaManagedAgentsSpanModelRequestStartEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "span.model_request_start".
	Type BetaManagedAgentsSpanModelRequestStartEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSpanModelRequestStartEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSpanModelRequestStartEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSpanModelRequestStartEventType string

const (
	BetaManagedAgentsSpanModelRequestStartEventTypeSpanModelRequestStart BetaManagedAgentsSpanModelRequestStartEventType = "span.model_request_start"
)

// Token usage for a single model request.
type BetaManagedAgentsSpanModelUsage struct {
	// Tokens used to create prompt cache in this request.
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens" api:"required"`
	// Tokens read from prompt cache in this request.
	CacheReadInputTokens int64 `json:"cache_read_input_tokens" api:"required"`
	// Input tokens consumed by this request.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// Output tokens generated by this request.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// Inference speed mode. `fast` provides significantly faster output token
	// generation at premium pricing. Not all models support `fast`; invalid
	// combinations are rejected at create time.
	//
	// Any of "standard", "fast".
	Speed BetaManagedAgentsSpanModelUsageSpeed `json:"speed" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CacheCreationInputTokens respjson.Field
		CacheReadInputTokens     respjson.Field
		InputTokens              respjson.Field
		OutputTokens             respjson.Field
		Speed                    respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSpanModelUsage) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSpanModelUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inference speed mode. `fast` provides significantly faster output token
// generation at premium pricing. Not all models support `fast`; invalid
// combinations are rejected at create time.
type BetaManagedAgentsSpanModelUsageSpeed string

const (
	BetaManagedAgentsSpanModelUsageSpeedStandard BetaManagedAgentsSpanModelUsageSpeed = "standard"
	BetaManagedAgentsSpanModelUsageSpeedFast     BetaManagedAgentsSpanModelUsageSpeed = "fast"
)

// BetaManagedAgentsStreamSessionEventsUnion contains all possible properties and
// values from [BetaManagedAgentsUserMessageEvent],
// [BetaManagedAgentsUserInterruptEvent],
// [BetaManagedAgentsUserToolConfirmationEvent],
// [BetaManagedAgentsUserCustomToolResultEvent],
// [BetaManagedAgentsAgentCustomToolUseEvent],
// [BetaManagedAgentsAgentMessageEvent], [BetaManagedAgentsAgentThinkingEvent],
// [BetaManagedAgentsAgentMCPToolUseEvent],
// [BetaManagedAgentsAgentMCPToolResultEvent],
// [BetaManagedAgentsAgentToolUseEvent], [BetaManagedAgentsAgentToolResultEvent],
// [BetaManagedAgentsAgentThreadContextCompactedEvent],
// [BetaManagedAgentsSessionErrorEvent],
// [BetaManagedAgentsSessionStatusRescheduledEvent],
// [BetaManagedAgentsSessionStatusRunningEvent],
// [BetaManagedAgentsSessionStatusIdleEvent],
// [BetaManagedAgentsSessionStatusTerminatedEvent],
// [BetaManagedAgentsSpanModelRequestStartEvent],
// [BetaManagedAgentsSpanModelRequestEndEvent],
// [BetaManagedAgentsSessionDeletedEvent].
//
// Use the [BetaManagedAgentsStreamSessionEventsUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsStreamSessionEventsUnion struct {
	ID string `json:"id"`
	// This field is a union of [[]BetaManagedAgentsUserMessageEventContentUnion],
	// [[]BetaManagedAgentsUserCustomToolResultEventContentUnion],
	// [[]BetaManagedAgentsTextBlock],
	// [[]BetaManagedAgentsAgentMCPToolResultEventContentUnion],
	// [[]BetaManagedAgentsAgentToolResultEventContentUnion]
	Content BetaManagedAgentsStreamSessionEventsUnionContent `json:"content"`
	// Any of "user.message", "user.interrupt", "user.tool_confirmation",
	// "user.custom_tool_result", "agent.custom_tool_use", "agent.message",
	// "agent.thinking", "agent.mcp_tool_use", "agent.mcp_tool_result",
	// "agent.tool_use", "agent.tool_result", "agent.thread_context_compacted",
	// "session.error", "session.status_rescheduled", "session.status_running",
	// "session.status_idle", "session.status_terminated", "span.model_request_start",
	// "span.model_request_end", "session.deleted".
	Type        string    `json:"type"`
	ProcessedAt time.Time `json:"processed_at"`
	// This field is from variant [BetaManagedAgentsUserToolConfirmationEvent].
	Result    BetaManagedAgentsUserToolConfirmationEventResult `json:"result"`
	ToolUseID string                                           `json:"tool_use_id"`
	// This field is from variant [BetaManagedAgentsUserToolConfirmationEvent].
	DenyMessage string `json:"deny_message"`
	// This field is from variant [BetaManagedAgentsUserCustomToolResultEvent].
	CustomToolUseID string `json:"custom_tool_use_id"`
	IsError         bool   `json:"is_error"`
	Input           any    `json:"input"`
	Name            string `json:"name"`
	// This field is from variant [BetaManagedAgentsAgentMCPToolUseEvent].
	MCPServerName       string `json:"mcp_server_name"`
	EvaluatedPermission string `json:"evaluated_permission"`
	// This field is from variant [BetaManagedAgentsAgentMCPToolResultEvent].
	MCPToolUseID string `json:"mcp_tool_use_id"`
	// This field is from variant [BetaManagedAgentsSessionErrorEvent].
	Error BetaManagedAgentsSessionErrorEventErrorUnion `json:"error"`
	// This field is from variant [BetaManagedAgentsSessionStatusIdleEvent].
	StopReason BetaManagedAgentsSessionStatusIdleEventStopReasonUnion `json:"stop_reason"`
	// This field is from variant [BetaManagedAgentsSpanModelRequestEndEvent].
	ModelRequestStartID string `json:"model_request_start_id"`
	// This field is from variant [BetaManagedAgentsSpanModelRequestEndEvent].
	ModelUsage BetaManagedAgentsSpanModelUsage `json:"model_usage"`
	JSON       struct {
		ID                  respjson.Field
		Content             respjson.Field
		Type                respjson.Field
		ProcessedAt         respjson.Field
		Result              respjson.Field
		ToolUseID           respjson.Field
		DenyMessage         respjson.Field
		CustomToolUseID     respjson.Field
		IsError             respjson.Field
		Input               respjson.Field
		Name                respjson.Field
		MCPServerName       respjson.Field
		EvaluatedPermission respjson.Field
		MCPToolUseID        respjson.Field
		Error               respjson.Field
		StopReason          respjson.Field
		ModelRequestStartID respjson.Field
		ModelUsage          respjson.Field
		raw                 string
	} `json:"-"`
}

// anyBetaManagedAgentsStreamSessionEvents is implemented by each variant of
// [BetaManagedAgentsStreamSessionEventsUnion] to add type safety for the return
// type of [BetaManagedAgentsStreamSessionEventsUnion.AsAny]
type anyBetaManagedAgentsStreamSessionEvents interface {
	implBetaManagedAgentsStreamSessionEventsUnion()
}

func (BetaManagedAgentsUserMessageEvent) implBetaManagedAgentsStreamSessionEventsUnion()          {}
func (BetaManagedAgentsUserInterruptEvent) implBetaManagedAgentsStreamSessionEventsUnion()        {}
func (BetaManagedAgentsUserToolConfirmationEvent) implBetaManagedAgentsStreamSessionEventsUnion() {}
func (BetaManagedAgentsUserCustomToolResultEvent) implBetaManagedAgentsStreamSessionEventsUnion() {}
func (BetaManagedAgentsAgentCustomToolUseEvent) implBetaManagedAgentsStreamSessionEventsUnion()   {}
func (BetaManagedAgentsAgentMessageEvent) implBetaManagedAgentsStreamSessionEventsUnion()         {}
func (BetaManagedAgentsAgentThinkingEvent) implBetaManagedAgentsStreamSessionEventsUnion()        {}
func (BetaManagedAgentsAgentMCPToolUseEvent) implBetaManagedAgentsStreamSessionEventsUnion()      {}
func (BetaManagedAgentsAgentMCPToolResultEvent) implBetaManagedAgentsStreamSessionEventsUnion()   {}
func (BetaManagedAgentsAgentToolUseEvent) implBetaManagedAgentsStreamSessionEventsUnion()         {}
func (BetaManagedAgentsAgentToolResultEvent) implBetaManagedAgentsStreamSessionEventsUnion()      {}
func (BetaManagedAgentsAgentThreadContextCompactedEvent) implBetaManagedAgentsStreamSessionEventsUnion() {
}
func (BetaManagedAgentsSessionErrorEvent) implBetaManagedAgentsStreamSessionEventsUnion() {}
func (BetaManagedAgentsSessionStatusRescheduledEvent) implBetaManagedAgentsStreamSessionEventsUnion() {
}
func (BetaManagedAgentsSessionStatusRunningEvent) implBetaManagedAgentsStreamSessionEventsUnion() {}
func (BetaManagedAgentsSessionStatusIdleEvent) implBetaManagedAgentsStreamSessionEventsUnion()    {}
func (BetaManagedAgentsSessionStatusTerminatedEvent) implBetaManagedAgentsStreamSessionEventsUnion() {
}
func (BetaManagedAgentsSpanModelRequestStartEvent) implBetaManagedAgentsStreamSessionEventsUnion() {}
func (BetaManagedAgentsSpanModelRequestEndEvent) implBetaManagedAgentsStreamSessionEventsUnion()   {}
func (BetaManagedAgentsSessionDeletedEvent) implBetaManagedAgentsStreamSessionEventsUnion()        {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsStreamSessionEventsUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsUserMessageEvent:
//	case anthropic.BetaManagedAgentsUserInterruptEvent:
//	case anthropic.BetaManagedAgentsUserToolConfirmationEvent:
//	case anthropic.BetaManagedAgentsUserCustomToolResultEvent:
//	case anthropic.BetaManagedAgentsAgentCustomToolUseEvent:
//	case anthropic.BetaManagedAgentsAgentMessageEvent:
//	case anthropic.BetaManagedAgentsAgentThinkingEvent:
//	case anthropic.BetaManagedAgentsAgentMCPToolUseEvent:
//	case anthropic.BetaManagedAgentsAgentMCPToolResultEvent:
//	case anthropic.BetaManagedAgentsAgentToolUseEvent:
//	case anthropic.BetaManagedAgentsAgentToolResultEvent:
//	case anthropic.BetaManagedAgentsAgentThreadContextCompactedEvent:
//	case anthropic.BetaManagedAgentsSessionErrorEvent:
//	case anthropic.BetaManagedAgentsSessionStatusRescheduledEvent:
//	case anthropic.BetaManagedAgentsSessionStatusRunningEvent:
//	case anthropic.BetaManagedAgentsSessionStatusIdleEvent:
//	case anthropic.BetaManagedAgentsSessionStatusTerminatedEvent:
//	case anthropic.BetaManagedAgentsSpanModelRequestStartEvent:
//	case anthropic.BetaManagedAgentsSpanModelRequestEndEvent:
//	case anthropic.BetaManagedAgentsSessionDeletedEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsStreamSessionEventsUnion) AsAny() anyBetaManagedAgentsStreamSessionEvents {
	switch u.Type {
	case "user.message":
		return u.AsUserMessage()
	case "user.interrupt":
		return u.AsUserInterrupt()
	case "user.tool_confirmation":
		return u.AsUserToolConfirmation()
	case "user.custom_tool_result":
		return u.AsUserCustomToolResult()
	case "agent.custom_tool_use":
		return u.AsAgentCustomToolUse()
	case "agent.message":
		return u.AsAgentMessage()
	case "agent.thinking":
		return u.AsAgentThinking()
	case "agent.mcp_tool_use":
		return u.AsAgentMCPToolUse()
	case "agent.mcp_tool_result":
		return u.AsAgentMCPToolResult()
	case "agent.tool_use":
		return u.AsAgentToolUse()
	case "agent.tool_result":
		return u.AsAgentToolResult()
	case "agent.thread_context_compacted":
		return u.AsAgentThreadContextCompacted()
	case "session.error":
		return u.AsSessionError()
	case "session.status_rescheduled":
		return u.AsSessionStatusRescheduled()
	case "session.status_running":
		return u.AsSessionStatusRunning()
	case "session.status_idle":
		return u.AsSessionStatusIdle()
	case "session.status_terminated":
		return u.AsSessionStatusTerminated()
	case "span.model_request_start":
		return u.AsSpanModelRequestStart()
	case "span.model_request_end":
		return u.AsSpanModelRequestEnd()
	case "session.deleted":
		return u.AsSessionDeleted()
	}
	return nil
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsUserMessage() (v BetaManagedAgentsUserMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsUserInterrupt() (v BetaManagedAgentsUserInterruptEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsUserToolConfirmation() (v BetaManagedAgentsUserToolConfirmationEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsUserCustomToolResult() (v BetaManagedAgentsUserCustomToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsAgentCustomToolUse() (v BetaManagedAgentsAgentCustomToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsAgentMessage() (v BetaManagedAgentsAgentMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsAgentThinking() (v BetaManagedAgentsAgentThinkingEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsAgentMCPToolUse() (v BetaManagedAgentsAgentMCPToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsAgentMCPToolResult() (v BetaManagedAgentsAgentMCPToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsAgentToolUse() (v BetaManagedAgentsAgentToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsAgentToolResult() (v BetaManagedAgentsAgentToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsAgentThreadContextCompacted() (v BetaManagedAgentsAgentThreadContextCompactedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsSessionError() (v BetaManagedAgentsSessionErrorEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsSessionStatusRescheduled() (v BetaManagedAgentsSessionStatusRescheduledEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsSessionStatusRunning() (v BetaManagedAgentsSessionStatusRunningEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsSessionStatusIdle() (v BetaManagedAgentsSessionStatusIdleEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsSessionStatusTerminated() (v BetaManagedAgentsSessionStatusTerminatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsSpanModelRequestStart() (v BetaManagedAgentsSpanModelRequestStartEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsSpanModelRequestEnd() (v BetaManagedAgentsSpanModelRequestEndEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionEventsUnion) AsSessionDeleted() (v BetaManagedAgentsSessionDeletedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsStreamSessionEventsUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsStreamSessionEventsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsStreamSessionEventsUnionContent is an implicit subunion of
// [BetaManagedAgentsStreamSessionEventsUnion].
// BetaManagedAgentsStreamSessionEventsUnionContent provides convenient access to
// the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsStreamSessionEventsUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfBetaManagedAgentsUserMessageEventContentArray
// OfBetaManagedAgentsUserCustomToolResultEventContentArray
// OfBetaManagedAgentsTextBlockArray
// OfBetaManagedAgentsAgentMCPToolResultEventContentArray
// OfBetaManagedAgentsAgentToolResultEventContentArray]
type BetaManagedAgentsStreamSessionEventsUnionContent struct {
	// This field will be present if the value is a
	// [[]BetaManagedAgentsUserMessageEventContentUnion] instead of an object.
	OfBetaManagedAgentsUserMessageEventContentArray []BetaManagedAgentsUserMessageEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]BetaManagedAgentsUserCustomToolResultEventContentUnion] instead of an object.
	OfBetaManagedAgentsUserCustomToolResultEventContentArray []BetaManagedAgentsUserCustomToolResultEventContentUnion `json:",inline"`
	// This field will be present if the value is a [[]BetaManagedAgentsTextBlock]
	// instead of an object.
	OfBetaManagedAgentsTextBlockArray []BetaManagedAgentsTextBlock `json:",inline"`
	// This field will be present if the value is a
	// [[]BetaManagedAgentsAgentMCPToolResultEventContentUnion] instead of an object.
	OfBetaManagedAgentsAgentMCPToolResultEventContentArray []BetaManagedAgentsAgentMCPToolResultEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]BetaManagedAgentsAgentToolResultEventContentUnion] instead of an object.
	OfBetaManagedAgentsAgentToolResultEventContentArray []BetaManagedAgentsAgentToolResultEventContentUnion `json:",inline"`
	JSON                                                struct {
		OfBetaManagedAgentsUserMessageEventContentArray          respjson.Field
		OfBetaManagedAgentsUserCustomToolResultEventContentArray respjson.Field
		OfBetaManagedAgentsTextBlockArray                        respjson.Field
		OfBetaManagedAgentsAgentMCPToolResultEventContentArray   respjson.Field
		OfBetaManagedAgentsAgentToolResultEventContentArray      respjson.Field
		raw                                                      string
	} `json:"-"`
}

func (r *BetaManagedAgentsStreamSessionEventsUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Regular text content.
type BetaManagedAgentsTextBlock struct {
	// The text content.
	Text string `json:"text" api:"required"`
	// Any of "text".
	Type BetaManagedAgentsTextBlockType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsTextBlock) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsTextBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsTextBlock to a
// BetaManagedAgentsTextBlockParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsTextBlockParam.Overrides()
func (r BetaManagedAgentsTextBlock) ToParam() BetaManagedAgentsTextBlockParam {
	return param.Override[BetaManagedAgentsTextBlockParam](json.RawMessage(r.RawJSON()))
}

type BetaManagedAgentsTextBlockType string

const (
	BetaManagedAgentsTextBlockTypeText BetaManagedAgentsTextBlockType = "text"
)

// Regular text content.
//
// The properties Text, Type are required.
type BetaManagedAgentsTextBlockParam struct {
	// The text content.
	Text string `json:"text" api:"required"`
	// Any of "text".
	Type BetaManagedAgentsTextBlockType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsTextBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsTextBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsTextBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An unknown or unexpected error occurred during session execution. A fallback
// variant; clients that don't recognize a new error code can match on
// `retry_status` and `message` alone.
type BetaManagedAgentsUnknownError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// What the client should do next in response to this error.
	RetryStatus BetaManagedAgentsUnknownErrorRetryStatusUnion `json:"retry_status" api:"required"`
	// Any of "unknown_error".
	Type BetaManagedAgentsUnknownErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		RetryStatus respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsUnknownError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsUnknownError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsUnknownErrorRetryStatusUnion contains all possible properties
// and values from [BetaManagedAgentsRetryStatusRetrying],
// [BetaManagedAgentsRetryStatusExhausted], [BetaManagedAgentsRetryStatusTerminal].
//
// Use the [BetaManagedAgentsUnknownErrorRetryStatusUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsUnknownErrorRetryStatusUnion struct {
	// Any of "retrying", "exhausted", "terminal".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsUnknownErrorRetryStatus is implemented by each variant of
// [BetaManagedAgentsUnknownErrorRetryStatusUnion] to add type safety for the
// return type of [BetaManagedAgentsUnknownErrorRetryStatusUnion.AsAny]
type anyBetaManagedAgentsUnknownErrorRetryStatus interface {
	implBetaManagedAgentsUnknownErrorRetryStatusUnion()
}

func (BetaManagedAgentsRetryStatusRetrying) implBetaManagedAgentsUnknownErrorRetryStatusUnion()  {}
func (BetaManagedAgentsRetryStatusExhausted) implBetaManagedAgentsUnknownErrorRetryStatusUnion() {}
func (BetaManagedAgentsRetryStatusTerminal) implBetaManagedAgentsUnknownErrorRetryStatusUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsUnknownErrorRetryStatusUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsRetryStatusRetrying:
//	case anthropic.BetaManagedAgentsRetryStatusExhausted:
//	case anthropic.BetaManagedAgentsRetryStatusTerminal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsUnknownErrorRetryStatusUnion) AsAny() anyBetaManagedAgentsUnknownErrorRetryStatus {
	switch u.Type {
	case "retrying":
		return u.AsRetrying()
	case "exhausted":
		return u.AsExhausted()
	case "terminal":
		return u.AsTerminal()
	}
	return nil
}

func (u BetaManagedAgentsUnknownErrorRetryStatusUnion) AsRetrying() (v BetaManagedAgentsRetryStatusRetrying) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsUnknownErrorRetryStatusUnion) AsExhausted() (v BetaManagedAgentsRetryStatusExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsUnknownErrorRetryStatusUnion) AsTerminal() (v BetaManagedAgentsRetryStatusTerminal) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsUnknownErrorRetryStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsUnknownErrorRetryStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsUnknownErrorType string

const (
	BetaManagedAgentsUnknownErrorTypeUnknownError BetaManagedAgentsUnknownErrorType = "unknown_error"
)

// Document referenced by URL.
type BetaManagedAgentsURLDocumentSource struct {
	// Any of "url".
	Type BetaManagedAgentsURLDocumentSourceType `json:"type" api:"required"`
	// URL of the document to fetch.
	URL string `json:"url" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsURLDocumentSource) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsURLDocumentSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsURLDocumentSource to a
// BetaManagedAgentsURLDocumentSourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsURLDocumentSourceParam.Overrides()
func (r BetaManagedAgentsURLDocumentSource) ToParam() BetaManagedAgentsURLDocumentSourceParam {
	return param.Override[BetaManagedAgentsURLDocumentSourceParam](json.RawMessage(r.RawJSON()))
}

type BetaManagedAgentsURLDocumentSourceType string

const (
	BetaManagedAgentsURLDocumentSourceTypeURL BetaManagedAgentsURLDocumentSourceType = "url"
)

// Document referenced by URL.
//
// The properties Type, URL are required.
type BetaManagedAgentsURLDocumentSourceParam struct {
	// Any of "url".
	Type BetaManagedAgentsURLDocumentSourceType `json:"type,omitzero" api:"required"`
	// URL of the document to fetch.
	URL string `json:"url" api:"required"`
	paramObj
}

func (r BetaManagedAgentsURLDocumentSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsURLDocumentSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsURLDocumentSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Image referenced by URL.
type BetaManagedAgentsURLImageSource struct {
	// Any of "url".
	Type BetaManagedAgentsURLImageSourceType `json:"type" api:"required"`
	// URL of the image to fetch.
	URL string `json:"url" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsURLImageSource) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsURLImageSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsURLImageSource to a
// BetaManagedAgentsURLImageSourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsURLImageSourceParam.Overrides()
func (r BetaManagedAgentsURLImageSource) ToParam() BetaManagedAgentsURLImageSourceParam {
	return param.Override[BetaManagedAgentsURLImageSourceParam](json.RawMessage(r.RawJSON()))
}

type BetaManagedAgentsURLImageSourceType string

const (
	BetaManagedAgentsURLImageSourceTypeURL BetaManagedAgentsURLImageSourceType = "url"
)

// Image referenced by URL.
//
// The properties Type, URL are required.
type BetaManagedAgentsURLImageSourceParam struct {
	// Any of "url".
	Type BetaManagedAgentsURLImageSourceType `json:"type,omitzero" api:"required"`
	// URL of the image to fetch.
	URL string `json:"url" api:"required"`
	paramObj
}

func (r BetaManagedAgentsURLImageSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsURLImageSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsURLImageSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event sent by the client providing the result of a custom tool execution.
type BetaManagedAgentsUserCustomToolResultEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// The id of the `agent.custom_tool_use` event this result corresponds to, which
	// can be found in the last `session.status_idle`
	// [event's](https://platform.claude.com/docs/en/api/beta/sessions/events/list#beta_managed_agents_session_requires_action.event_ids)
	// `stop_reason.event_ids` field.
	CustomToolUseID string `json:"custom_tool_use_id" api:"required"`
	// Any of "user.custom_tool_result".
	Type BetaManagedAgentsUserCustomToolResultEventType `json:"type" api:"required"`
	// The result content returned by the tool.
	Content []BetaManagedAgentsUserCustomToolResultEventContentUnion `json:"content"`
	// Whether the tool execution resulted in an error.
	IsError bool `json:"is_error" api:"nullable"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"nullable" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		CustomToolUseID respjson.Field
		Type            respjson.Field
		Content         respjson.Field
		IsError         respjson.Field
		ProcessedAt     respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsUserCustomToolResultEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsUserCustomToolResultEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsUserCustomToolResultEventType string

const (
	BetaManagedAgentsUserCustomToolResultEventTypeUserCustomToolResult BetaManagedAgentsUserCustomToolResultEventType = "user.custom_tool_result"
)

// BetaManagedAgentsUserCustomToolResultEventContentUnion contains all possible
// properties and values from [BetaManagedAgentsTextBlock],
// [BetaManagedAgentsImageBlock], [BetaManagedAgentsDocumentBlock].
//
// Use the [BetaManagedAgentsUserCustomToolResultEventContentUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsUserCustomToolResultEventContentUnion struct {
	// This field is from variant [BetaManagedAgentsTextBlock].
	Text string `json:"text"`
	// Any of "text", "image", "document".
	Type string `json:"type"`
	// This field is a union of [BetaManagedAgentsImageBlockSourceUnion],
	// [BetaManagedAgentsDocumentBlockSourceUnion]
	Source BetaManagedAgentsUserCustomToolResultEventContentUnionSource `json:"source"`
	// This field is from variant [BetaManagedAgentsDocumentBlock].
	Context string `json:"context"`
	// This field is from variant [BetaManagedAgentsDocumentBlock].
	Title string `json:"title"`
	JSON  struct {
		Text    respjson.Field
		Type    respjson.Field
		Source  respjson.Field
		Context respjson.Field
		Title   respjson.Field
		raw     string
	} `json:"-"`
}

// anyBetaManagedAgentsUserCustomToolResultEventContent is implemented by each
// variant of [BetaManagedAgentsUserCustomToolResultEventContentUnion] to add type
// safety for the return type of
// [BetaManagedAgentsUserCustomToolResultEventContentUnion.AsAny]
type anyBetaManagedAgentsUserCustomToolResultEventContent interface {
	implBetaManagedAgentsUserCustomToolResultEventContentUnion()
}

func (BetaManagedAgentsTextBlock) implBetaManagedAgentsUserCustomToolResultEventContentUnion()     {}
func (BetaManagedAgentsImageBlock) implBetaManagedAgentsUserCustomToolResultEventContentUnion()    {}
func (BetaManagedAgentsDocumentBlock) implBetaManagedAgentsUserCustomToolResultEventContentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsUserCustomToolResultEventContentUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsTextBlock:
//	case anthropic.BetaManagedAgentsImageBlock:
//	case anthropic.BetaManagedAgentsDocumentBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsUserCustomToolResultEventContentUnion) AsAny() anyBetaManagedAgentsUserCustomToolResultEventContent {
	switch u.Type {
	case "text":
		return u.AsText()
	case "image":
		return u.AsImage()
	case "document":
		return u.AsDocument()
	}
	return nil
}

func (u BetaManagedAgentsUserCustomToolResultEventContentUnion) AsText() (v BetaManagedAgentsTextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsUserCustomToolResultEventContentUnion) AsImage() (v BetaManagedAgentsImageBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsUserCustomToolResultEventContentUnion) AsDocument() (v BetaManagedAgentsDocumentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsUserCustomToolResultEventContentUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsUserCustomToolResultEventContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsUserCustomToolResultEventContentUnionSource is an implicit
// subunion of [BetaManagedAgentsUserCustomToolResultEventContentUnion].
// BetaManagedAgentsUserCustomToolResultEventContentUnionSource provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsUserCustomToolResultEventContentUnion].
type BetaManagedAgentsUserCustomToolResultEventContentUnionSource struct {
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	FileID    string `json:"file_id"`
	JSON      struct {
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

func (r *BetaManagedAgentsUserCustomToolResultEventContentUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Parameters for providing the result of a custom tool execution.
//
// The properties CustomToolUseID, Type are required.
type BetaManagedAgentsUserCustomToolResultEventParams struct {
	// The id of the `agent.custom_tool_use` event this result corresponds to, which
	// can be found in the last `session.status_idle`
	// [event's](https://platform.claude.com/docs/en/api/beta/sessions/events/list#beta_managed_agents_session_requires_action.event_ids)
	// `stop_reason.event_ids` field.
	CustomToolUseID string `json:"custom_tool_use_id" api:"required"`
	// Any of "user.custom_tool_result".
	Type BetaManagedAgentsUserCustomToolResultEventParamsType `json:"type,omitzero" api:"required"`
	// Whether the tool execution resulted in an error.
	IsError param.Opt[bool] `json:"is_error,omitzero"`
	// The result content returned by the tool.
	Content []BetaManagedAgentsUserCustomToolResultEventParamsContentUnion `json:"content,omitzero"`
	paramObj
}

func (r BetaManagedAgentsUserCustomToolResultEventParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsUserCustomToolResultEventParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsUserCustomToolResultEventParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsUserCustomToolResultEventParamsType string

const (
	BetaManagedAgentsUserCustomToolResultEventParamsTypeUserCustomToolResult BetaManagedAgentsUserCustomToolResultEventParamsType = "user.custom_tool_result"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsUserCustomToolResultEventParamsContentUnion struct {
	OfText     *BetaManagedAgentsTextBlockParam     `json:",omitzero,inline"`
	OfImage    *BetaManagedAgentsImageBlockParam    `json:",omitzero,inline"`
	OfDocument *BetaManagedAgentsDocumentBlockParam `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsUserCustomToolResultEventParamsContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfText, u.OfImage, u.OfDocument)
}
func (u *BetaManagedAgentsUserCustomToolResultEventParamsContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsUserCustomToolResultEventParamsContentUnion) asAny() any {
	if !param.IsOmitted(u.OfText) {
		return u.OfText
	} else if !param.IsOmitted(u.OfImage) {
		return u.OfImage
	} else if !param.IsOmitted(u.OfDocument) {
		return u.OfDocument
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsUserCustomToolResultEventParamsContentUnion) GetText() *string {
	if vt := u.OfText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsUserCustomToolResultEventParamsContentUnion) GetContext() *string {
	if vt := u.OfDocument; vt != nil && vt.Context.Valid() {
		return &vt.Context.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsUserCustomToolResultEventParamsContentUnion) GetTitle() *string {
	if vt := u.OfDocument; vt != nil && vt.Title.Valid() {
		return &vt.Title.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsUserCustomToolResultEventParamsContentUnion) GetType() *string {
	if vt := u.OfText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfImage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDocument; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u BetaManagedAgentsUserCustomToolResultEventParamsContentUnion) GetSource() (res betaManagedAgentsUserCustomToolResultEventParamsContentUnionSource) {
	if vt := u.OfImage; vt != nil {
		res.any = vt.Source.asAny()
	} else if vt := u.OfDocument; vt != nil {
		res.any = vt.Source.asAny()
	}
	return
}

// Can have the runtime types [*BetaManagedAgentsBase64ImageSourceParam],
// [*BetaManagedAgentsURLImageSourceParam],
// [*BetaManagedAgentsFileImageSourceParam],
// [*BetaManagedAgentsBase64DocumentSourceParam],
// [*BetaManagedAgentsPlainTextDocumentSourceParam],
// [*BetaManagedAgentsURLDocumentSourceParam],
// [*BetaManagedAgentsFileDocumentSourceParam]
type betaManagedAgentsUserCustomToolResultEventParamsContentUnionSource struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *anthropic.BetaManagedAgentsBase64ImageSourceParam:
//	case *anthropic.BetaManagedAgentsURLImageSourceParam:
//	case *anthropic.BetaManagedAgentsFileImageSourceParam:
//	case *anthropic.BetaManagedAgentsBase64DocumentSourceParam:
//	case *anthropic.BetaManagedAgentsPlainTextDocumentSourceParam:
//	case *anthropic.BetaManagedAgentsURLDocumentSourceParam:
//	case *anthropic.BetaManagedAgentsFileDocumentSourceParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u betaManagedAgentsUserCustomToolResultEventParamsContentUnionSource) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u betaManagedAgentsUserCustomToolResultEventParamsContentUnionSource) GetData() *string {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsImageBlockSourceUnionParam:
		return vt.GetData()
	case *BetaManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetData()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u betaManagedAgentsUserCustomToolResultEventParamsContentUnionSource) GetMediaType() *string {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsImageBlockSourceUnionParam:
		return vt.GetMediaType()
	case *BetaManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetMediaType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u betaManagedAgentsUserCustomToolResultEventParamsContentUnionSource) GetType() *string {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsImageBlockSourceUnionParam:
		return vt.GetType()
	case *BetaManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u betaManagedAgentsUserCustomToolResultEventParamsContentUnionSource) GetURL() *string {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsImageBlockSourceUnionParam:
		return vt.GetURL()
	case *BetaManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetURL()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u betaManagedAgentsUserCustomToolResultEventParamsContentUnionSource) GetFileID() *string {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsImageBlockSourceUnionParam:
		return vt.GetFileID()
	case *BetaManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetFileID()
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsUserCustomToolResultEventParamsContentUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsTextBlockParam]("text"),
		apijson.Discriminator[BetaManagedAgentsImageBlockParam]("image"),
		apijson.Discriminator[BetaManagedAgentsDocumentBlockParam]("document"),
	)
}

// An interrupt event that pauses agent execution and returns control to the user.
type BetaManagedAgentsUserInterruptEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Any of "user.interrupt".
	Type BetaManagedAgentsUserInterruptEventType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"nullable" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ProcessedAt respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsUserInterruptEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsUserInterruptEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsUserInterruptEventType string

const (
	BetaManagedAgentsUserInterruptEventTypeUserInterrupt BetaManagedAgentsUserInterruptEventType = "user.interrupt"
)

// Parameters for sending an interrupt to pause the agent.
//
// The property Type is required.
type BetaManagedAgentsUserInterruptEventParams struct {
	// Any of "user.interrupt".
	Type BetaManagedAgentsUserInterruptEventParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsUserInterruptEventParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsUserInterruptEventParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsUserInterruptEventParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsUserInterruptEventParamsType string

const (
	BetaManagedAgentsUserInterruptEventParamsTypeUserInterrupt BetaManagedAgentsUserInterruptEventParamsType = "user.interrupt"
)

// A user message event in the session conversation.
type BetaManagedAgentsUserMessageEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Array of content blocks comprising the user message.
	Content []BetaManagedAgentsUserMessageEventContentUnion `json:"content" api:"required"`
	// Any of "user.message".
	Type BetaManagedAgentsUserMessageEventType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"nullable" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Content     respjson.Field
		Type        respjson.Field
		ProcessedAt respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsUserMessageEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsUserMessageEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsUserMessageEventContentUnion contains all possible properties
// and values from [BetaManagedAgentsTextBlock], [BetaManagedAgentsImageBlock],
// [BetaManagedAgentsDocumentBlock].
//
// Use the [BetaManagedAgentsUserMessageEventContentUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsUserMessageEventContentUnion struct {
	// This field is from variant [BetaManagedAgentsTextBlock].
	Text string `json:"text"`
	// Any of "text", "image", "document".
	Type string `json:"type"`
	// This field is a union of [BetaManagedAgentsImageBlockSourceUnion],
	// [BetaManagedAgentsDocumentBlockSourceUnion]
	Source BetaManagedAgentsUserMessageEventContentUnionSource `json:"source"`
	// This field is from variant [BetaManagedAgentsDocumentBlock].
	Context string `json:"context"`
	// This field is from variant [BetaManagedAgentsDocumentBlock].
	Title string `json:"title"`
	JSON  struct {
		Text    respjson.Field
		Type    respjson.Field
		Source  respjson.Field
		Context respjson.Field
		Title   respjson.Field
		raw     string
	} `json:"-"`
}

// anyBetaManagedAgentsUserMessageEventContent is implemented by each variant of
// [BetaManagedAgentsUserMessageEventContentUnion] to add type safety for the
// return type of [BetaManagedAgentsUserMessageEventContentUnion.AsAny]
type anyBetaManagedAgentsUserMessageEventContent interface {
	implBetaManagedAgentsUserMessageEventContentUnion()
}

func (BetaManagedAgentsTextBlock) implBetaManagedAgentsUserMessageEventContentUnion()     {}
func (BetaManagedAgentsImageBlock) implBetaManagedAgentsUserMessageEventContentUnion()    {}
func (BetaManagedAgentsDocumentBlock) implBetaManagedAgentsUserMessageEventContentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsUserMessageEventContentUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsTextBlock:
//	case anthropic.BetaManagedAgentsImageBlock:
//	case anthropic.BetaManagedAgentsDocumentBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsUserMessageEventContentUnion) AsAny() anyBetaManagedAgentsUserMessageEventContent {
	switch u.Type {
	case "text":
		return u.AsText()
	case "image":
		return u.AsImage()
	case "document":
		return u.AsDocument()
	}
	return nil
}

func (u BetaManagedAgentsUserMessageEventContentUnion) AsText() (v BetaManagedAgentsTextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsUserMessageEventContentUnion) AsImage() (v BetaManagedAgentsImageBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsUserMessageEventContentUnion) AsDocument() (v BetaManagedAgentsDocumentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsUserMessageEventContentUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsUserMessageEventContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsUserMessageEventContentUnionSource is an implicit subunion of
// [BetaManagedAgentsUserMessageEventContentUnion].
// BetaManagedAgentsUserMessageEventContentUnionSource provides convenient access
// to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsUserMessageEventContentUnion].
type BetaManagedAgentsUserMessageEventContentUnionSource struct {
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	FileID    string `json:"file_id"`
	JSON      struct {
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

func (r *BetaManagedAgentsUserMessageEventContentUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsUserMessageEventType string

const (
	BetaManagedAgentsUserMessageEventTypeUserMessage BetaManagedAgentsUserMessageEventType = "user.message"
)

// Parameters for sending a user message to the session.
//
// The properties Content, Type are required.
type BetaManagedAgentsUserMessageEventParams struct {
	// Array of content blocks for the user message.
	Content []BetaManagedAgentsUserMessageEventParamsContentUnion `json:"content,omitzero" api:"required"`
	// Any of "user.message".
	Type BetaManagedAgentsUserMessageEventParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsUserMessageEventParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsUserMessageEventParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsUserMessageEventParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsUserMessageEventParamsContentUnion struct {
	OfText     *BetaManagedAgentsTextBlockParam     `json:",omitzero,inline"`
	OfImage    *BetaManagedAgentsImageBlockParam    `json:",omitzero,inline"`
	OfDocument *BetaManagedAgentsDocumentBlockParam `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsUserMessageEventParamsContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfText, u.OfImage, u.OfDocument)
}
func (u *BetaManagedAgentsUserMessageEventParamsContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsUserMessageEventParamsContentUnion) asAny() any {
	if !param.IsOmitted(u.OfText) {
		return u.OfText
	} else if !param.IsOmitted(u.OfImage) {
		return u.OfImage
	} else if !param.IsOmitted(u.OfDocument) {
		return u.OfDocument
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsUserMessageEventParamsContentUnion) GetText() *string {
	if vt := u.OfText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsUserMessageEventParamsContentUnion) GetContext() *string {
	if vt := u.OfDocument; vt != nil && vt.Context.Valid() {
		return &vt.Context.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsUserMessageEventParamsContentUnion) GetTitle() *string {
	if vt := u.OfDocument; vt != nil && vt.Title.Valid() {
		return &vt.Title.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsUserMessageEventParamsContentUnion) GetType() *string {
	if vt := u.OfText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfImage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDocument; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u BetaManagedAgentsUserMessageEventParamsContentUnion) GetSource() (res betaManagedAgentsUserMessageEventParamsContentUnionSource) {
	if vt := u.OfImage; vt != nil {
		res.any = vt.Source.asAny()
	} else if vt := u.OfDocument; vt != nil {
		res.any = vt.Source.asAny()
	}
	return
}

// Can have the runtime types [*BetaManagedAgentsBase64ImageSourceParam],
// [*BetaManagedAgentsURLImageSourceParam],
// [*BetaManagedAgentsFileImageSourceParam],
// [*BetaManagedAgentsBase64DocumentSourceParam],
// [*BetaManagedAgentsPlainTextDocumentSourceParam],
// [*BetaManagedAgentsURLDocumentSourceParam],
// [*BetaManagedAgentsFileDocumentSourceParam]
type betaManagedAgentsUserMessageEventParamsContentUnionSource struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *anthropic.BetaManagedAgentsBase64ImageSourceParam:
//	case *anthropic.BetaManagedAgentsURLImageSourceParam:
//	case *anthropic.BetaManagedAgentsFileImageSourceParam:
//	case *anthropic.BetaManagedAgentsBase64DocumentSourceParam:
//	case *anthropic.BetaManagedAgentsPlainTextDocumentSourceParam:
//	case *anthropic.BetaManagedAgentsURLDocumentSourceParam:
//	case *anthropic.BetaManagedAgentsFileDocumentSourceParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u betaManagedAgentsUserMessageEventParamsContentUnionSource) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u betaManagedAgentsUserMessageEventParamsContentUnionSource) GetData() *string {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsImageBlockSourceUnionParam:
		return vt.GetData()
	case *BetaManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetData()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u betaManagedAgentsUserMessageEventParamsContentUnionSource) GetMediaType() *string {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsImageBlockSourceUnionParam:
		return vt.GetMediaType()
	case *BetaManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetMediaType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u betaManagedAgentsUserMessageEventParamsContentUnionSource) GetType() *string {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsImageBlockSourceUnionParam:
		return vt.GetType()
	case *BetaManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u betaManagedAgentsUserMessageEventParamsContentUnionSource) GetURL() *string {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsImageBlockSourceUnionParam:
		return vt.GetURL()
	case *BetaManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetURL()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u betaManagedAgentsUserMessageEventParamsContentUnionSource) GetFileID() *string {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsImageBlockSourceUnionParam:
		return vt.GetFileID()
	case *BetaManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetFileID()
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsUserMessageEventParamsContentUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsTextBlockParam]("text"),
		apijson.Discriminator[BetaManagedAgentsImageBlockParam]("image"),
		apijson.Discriminator[BetaManagedAgentsDocumentBlockParam]("document"),
	)
}

type BetaManagedAgentsUserMessageEventParamsType string

const (
	BetaManagedAgentsUserMessageEventParamsTypeUserMessage BetaManagedAgentsUserMessageEventParamsType = "user.message"
)

// A tool confirmation event that approves or denies a pending tool execution.
type BetaManagedAgentsUserToolConfirmationEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// UserToolConfirmationResult enum
	//
	// Any of "allow", "deny".
	Result BetaManagedAgentsUserToolConfirmationEventResult `json:"result" api:"required"`
	// The id of the `agent.tool_use` or `agent.mcp_tool_use` event this result
	// corresponds to, which can be found in the last `session.status_idle`
	// [event's](https://platform.claude.com/docs/en/api/beta/sessions/events/list#beta_managed_agents_session_requires_action.event_ids)
	// `stop_reason.event_ids` field.
	ToolUseID string `json:"tool_use_id" api:"required"`
	// Any of "user.tool_confirmation".
	Type BetaManagedAgentsUserToolConfirmationEventType `json:"type" api:"required"`
	// Optional message providing context for a 'deny' decision. Only allowed when
	// result is 'deny'.
	DenyMessage string `json:"deny_message" api:"nullable"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"nullable" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Result      respjson.Field
		ToolUseID   respjson.Field
		Type        respjson.Field
		DenyMessage respjson.Field
		ProcessedAt respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsUserToolConfirmationEvent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsUserToolConfirmationEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// UserToolConfirmationResult enum
type BetaManagedAgentsUserToolConfirmationEventResult string

const (
	BetaManagedAgentsUserToolConfirmationEventResultAllow BetaManagedAgentsUserToolConfirmationEventResult = "allow"
	BetaManagedAgentsUserToolConfirmationEventResultDeny  BetaManagedAgentsUserToolConfirmationEventResult = "deny"
)

type BetaManagedAgentsUserToolConfirmationEventType string

const (
	BetaManagedAgentsUserToolConfirmationEventTypeUserToolConfirmation BetaManagedAgentsUserToolConfirmationEventType = "user.tool_confirmation"
)

// Parameters for confirming or denying a tool execution request.
//
// The properties Result, ToolUseID, Type are required.
type BetaManagedAgentsUserToolConfirmationEventParams struct {
	// UserToolConfirmationResult enum
	//
	// Any of "allow", "deny".
	Result BetaManagedAgentsUserToolConfirmationEventParamsResult `json:"result,omitzero" api:"required"`
	// The id of the `agent.tool_use` or `agent.mcp_tool_use` event this result
	// corresponds to, which can be found in the last `session.status_idle`
	// [event's](https://platform.claude.com/docs/en/api/beta/sessions/events/list#beta_managed_agents_session_requires_action.event_ids)
	// `stop_reason.event_ids` field.
	ToolUseID string `json:"tool_use_id" api:"required"`
	// Any of "user.tool_confirmation".
	Type BetaManagedAgentsUserToolConfirmationEventParamsType `json:"type,omitzero" api:"required"`
	// Optional message providing context for a 'deny' decision. Only allowed when
	// result is 'deny'.
	DenyMessage param.Opt[string] `json:"deny_message,omitzero"`
	paramObj
}

func (r BetaManagedAgentsUserToolConfirmationEventParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsUserToolConfirmationEventParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsUserToolConfirmationEventParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// UserToolConfirmationResult enum
type BetaManagedAgentsUserToolConfirmationEventParamsResult string

const (
	BetaManagedAgentsUserToolConfirmationEventParamsResultAllow BetaManagedAgentsUserToolConfirmationEventParamsResult = "allow"
	BetaManagedAgentsUserToolConfirmationEventParamsResultDeny  BetaManagedAgentsUserToolConfirmationEventParamsResult = "deny"
)

type BetaManagedAgentsUserToolConfirmationEventParamsType string

const (
	BetaManagedAgentsUserToolConfirmationEventParamsTypeUserToolConfirmation BetaManagedAgentsUserToolConfirmationEventParamsType = "user.tool_confirmation"
)

type BetaSessionEventListParams struct {
	// Query parameter for limit
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor from a previous response's next_page.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Sort direction for results, ordered by created_at. Defaults to asc
	// (chronological).
	//
	// Any of "asc", "desc".
	Order BetaSessionEventListParamsOrder `query:"order,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaSessionEventListParams]'s query parameters as
// `url.Values`.
func (r BetaSessionEventListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort direction for results, ordered by created_at. Defaults to asc
// (chronological).
type BetaSessionEventListParamsOrder string

const (
	BetaSessionEventListParamsOrderAsc  BetaSessionEventListParamsOrder = "asc"
	BetaSessionEventListParamsOrderDesc BetaSessionEventListParamsOrder = "desc"
)

type BetaSessionEventSendParams struct {
	// Events to send to the `session`.
	Events []BetaManagedAgentsEventParamsUnion `json:"events,omitzero" api:"required"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaSessionEventSendParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaSessionEventSendParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaSessionEventSendParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaSessionEventStreamParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
