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
)

// BetaSessionThreadService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaSessionThreadService] method instead.
type BetaSessionThreadService struct {
	Options []option.RequestOption
	Events  BetaSessionThreadEventService
}

// NewBetaSessionThreadService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaSessionThreadService(opts ...option.RequestOption) (r BetaSessionThreadService) {
	r = BetaSessionThreadService{}
	r.Options = opts
	r.Events = NewBetaSessionThreadEventService(opts...)
	return
}

// Get Session Thread
func (r *BetaSessionThreadService) Get(ctx context.Context, threadID string, params BetaSessionThreadGetParams, opts ...option.RequestOption) (res *BetaManagedAgentsSessionThread, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if params.SessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/sessions/%s/threads/%s?beta=true", params.SessionID, threadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List Session Threads
func (r *BetaSessionThreadService) List(ctx context.Context, sessionID string, params BetaSessionThreadListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaManagedAgentsSessionThread], err error) {
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
	path := fmt.Sprintf("v1/sessions/%s/threads?beta=true", sessionID)
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

// List Session Threads
func (r *BetaSessionThreadService) ListAutoPaging(ctx context.Context, sessionID string, params BetaSessionThreadListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsSessionThread] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, sessionID, params, opts...))
}

// Archive Session Thread
func (r *BetaSessionThreadService) Archive(ctx context.Context, threadID string, params BetaSessionThreadArchiveParams, opts ...option.RequestOption) (res *BetaManagedAgentsSessionThread, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if params.SessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/sessions/%s/threads/%s/archive?beta=true", params.SessionID, threadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// An execution thread within a `session`. Each session has one primary thread plus
// zero or more child threads spawned by the coordinator.
type BetaManagedAgentsSessionThread struct {
	// Unique identifier for this thread.
	ID string `json:"id" api:"required"`
	// Resolved `agent` definition for a single `session_thread`. Snapshot of the agent
	// at thread creation time. The multiagent roster is not repeated here; read it
	// from `Session.agent`.
	Agent BetaManagedAgentsSessionThreadAgent `json:"agent" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Parent thread that spawned this thread. Null for the primary thread.
	ParentThreadID string `json:"parent_thread_id" api:"required"`
	// The session this thread belongs to.
	SessionID string `json:"session_id" api:"required"`
	// Timing statistics for a session thread.
	Stats BetaManagedAgentsSessionThreadStats `json:"stats" api:"required"`
	// SessionThreadStatus enum
	//
	// Any of "running", "idle", "rescheduling", "terminated".
	Status BetaManagedAgentsSessionThreadStatus `json:"status" api:"required"`
	// Any of "session_thread".
	Type BetaManagedAgentsSessionThreadType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// Cumulative token usage for a session thread across all turns.
	Usage BetaManagedAgentsSessionThreadUsage `json:"usage" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		Agent          respjson.Field
		ArchivedAt     respjson.Field
		CreatedAt      respjson.Field
		ParentThreadID respjson.Field
		SessionID      respjson.Field
		Stats          respjson.Field
		Status         respjson.Field
		Type           respjson.Field
		UpdatedAt      respjson.Field
		Usage          respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionThread) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionThread) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSessionThreadType string

const (
	BetaManagedAgentsSessionThreadTypeSessionThread BetaManagedAgentsSessionThreadType = "session_thread"
)

// Resolved `agent` definition for a single `session_thread`. Snapshot of the agent
// at thread creation time. The multiagent roster is not repeated here; read it
// from `Session.agent`.
type BetaManagedAgentsSessionThreadAgent struct {
	ID          string                                    `json:"id" api:"required"`
	Description string                                    `json:"description" api:"required"`
	MCPServers  []BetaManagedAgentsMCPServerURLDefinition `json:"mcp_servers" api:"required"`
	// Model identifier and configuration.
	Model  BetaManagedAgentsModelConfig                    `json:"model" api:"required"`
	Name   string                                          `json:"name" api:"required"`
	Skills []BetaManagedAgentsSessionThreadAgentSkillUnion `json:"skills" api:"required"`
	System string                                          `json:"system" api:"required"`
	Tools  []BetaManagedAgentsSessionThreadAgentToolUnion  `json:"tools" api:"required"`
	// Any of "agent".
	Type    BetaManagedAgentsSessionThreadAgentType `json:"type" api:"required"`
	Version int64                                   `json:"version" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Description respjson.Field
		MCPServers  respjson.Field
		Model       respjson.Field
		Name        respjson.Field
		Skills      respjson.Field
		System      respjson.Field
		Tools       respjson.Field
		Type        respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionThreadAgent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionThreadAgent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSessionThreadAgentSkillUnion contains all possible properties
// and values from [BetaManagedAgentsAnthropicSkill],
// [BetaManagedAgentsCustomSkill].
//
// Use the [BetaManagedAgentsSessionThreadAgentSkillUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsSessionThreadAgentSkillUnion struct {
	SkillID string `json:"skill_id"`
	// Any of "anthropic", "custom".
	Type    string `json:"type"`
	Version string `json:"version"`
	JSON    struct {
		SkillID respjson.Field
		Type    respjson.Field
		Version respjson.Field
		raw     string
	} `json:"-"`
}

// anyBetaManagedAgentsSessionThreadAgentSkill is implemented by each variant of
// [BetaManagedAgentsSessionThreadAgentSkillUnion] to add type safety for the
// return type of [BetaManagedAgentsSessionThreadAgentSkillUnion.AsAny]
type anyBetaManagedAgentsSessionThreadAgentSkill interface {
	implBetaManagedAgentsSessionThreadAgentSkillUnion()
}

func (BetaManagedAgentsAnthropicSkill) implBetaManagedAgentsSessionThreadAgentSkillUnion() {}
func (BetaManagedAgentsCustomSkill) implBetaManagedAgentsSessionThreadAgentSkillUnion()    {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsSessionThreadAgentSkillUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAnthropicSkill:
//	case anthropic.BetaManagedAgentsCustomSkill:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsSessionThreadAgentSkillUnion) AsAny() anyBetaManagedAgentsSessionThreadAgentSkill {
	switch u.Type {
	case "anthropic":
		return u.AsAnthropic()
	case "custom":
		return u.AsCustom()
	}
	return nil
}

func (u BetaManagedAgentsSessionThreadAgentSkillUnion) AsAnthropic() (v BetaManagedAgentsAnthropicSkill) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionThreadAgentSkillUnion) AsCustom() (v BetaManagedAgentsCustomSkill) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsSessionThreadAgentSkillUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsSessionThreadAgentSkillUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSessionThreadAgentToolUnion contains all possible properties
// and values from [BetaManagedAgentsAgentToolset20260401],
// [BetaManagedAgentsMCPToolset], [BetaManagedAgentsCustomTool].
//
// Use the [BetaManagedAgentsSessionThreadAgentToolUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsSessionThreadAgentToolUnion struct {
	// This field is a union of [[]BetaManagedAgentsAgentToolConfig],
	// [[]BetaManagedAgentsMCPToolConfig]
	Configs BetaManagedAgentsSessionThreadAgentToolUnionConfigs `json:"configs"`
	// This field is a union of [BetaManagedAgentsAgentToolsetDefaultConfig],
	// [BetaManagedAgentsMCPToolsetDefaultConfig]
	DefaultConfig BetaManagedAgentsSessionThreadAgentToolUnionDefaultConfig `json:"default_config"`
	// Any of "agent_toolset_20260401", "mcp_toolset", "custom".
	Type string `json:"type"`
	// This field is from variant [BetaManagedAgentsMCPToolset].
	MCPServerName string `json:"mcp_server_name"`
	// This field is from variant [BetaManagedAgentsCustomTool].
	Description string `json:"description"`
	// This field is from variant [BetaManagedAgentsCustomTool].
	InputSchema BetaManagedAgentsCustomToolInputSchema `json:"input_schema"`
	// This field is from variant [BetaManagedAgentsCustomTool].
	Name string `json:"name"`
	JSON struct {
		Configs       respjson.Field
		DefaultConfig respjson.Field
		Type          respjson.Field
		MCPServerName respjson.Field
		Description   respjson.Field
		InputSchema   respjson.Field
		Name          respjson.Field
		raw           string
	} `json:"-"`
}

// anyBetaManagedAgentsSessionThreadAgentTool is implemented by each variant of
// [BetaManagedAgentsSessionThreadAgentToolUnion] to add type safety for the return
// type of [BetaManagedAgentsSessionThreadAgentToolUnion.AsAny]
type anyBetaManagedAgentsSessionThreadAgentTool interface {
	implBetaManagedAgentsSessionThreadAgentToolUnion()
}

func (BetaManagedAgentsAgentToolset20260401) implBetaManagedAgentsSessionThreadAgentToolUnion() {}
func (BetaManagedAgentsMCPToolset) implBetaManagedAgentsSessionThreadAgentToolUnion()           {}
func (BetaManagedAgentsCustomTool) implBetaManagedAgentsSessionThreadAgentToolUnion()           {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsSessionThreadAgentToolUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAgentToolset20260401:
//	case anthropic.BetaManagedAgentsMCPToolset:
//	case anthropic.BetaManagedAgentsCustomTool:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsSessionThreadAgentToolUnion) AsAny() anyBetaManagedAgentsSessionThreadAgentTool {
	switch u.Type {
	case "agent_toolset_20260401":
		return u.AsAgentToolset20260401()
	case "mcp_toolset":
		return u.AsMCPToolset()
	case "custom":
		return u.AsCustom()
	}
	return nil
}

func (u BetaManagedAgentsSessionThreadAgentToolUnion) AsAgentToolset20260401() (v BetaManagedAgentsAgentToolset20260401) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionThreadAgentToolUnion) AsMCPToolset() (v BetaManagedAgentsMCPToolset) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionThreadAgentToolUnion) AsCustom() (v BetaManagedAgentsCustomTool) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsSessionThreadAgentToolUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsSessionThreadAgentToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSessionThreadAgentToolUnionConfigs is an implicit subunion of
// [BetaManagedAgentsSessionThreadAgentToolUnion].
// BetaManagedAgentsSessionThreadAgentToolUnionConfigs provides convenient access
// to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsSessionThreadAgentToolUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfBetaManagedAgentsAgentToolConfigArray
// OfBetaManagedAgentsMCPToolConfigArray]
type BetaManagedAgentsSessionThreadAgentToolUnionConfigs struct {
	// This field will be present if the value is a
	// [[]BetaManagedAgentsAgentToolConfig] instead of an object.
	OfBetaManagedAgentsAgentToolConfigArray []BetaManagedAgentsAgentToolConfig `json:",inline"`
	// This field will be present if the value is a [[]BetaManagedAgentsMCPToolConfig]
	// instead of an object.
	OfBetaManagedAgentsMCPToolConfigArray []BetaManagedAgentsMCPToolConfig `json:",inline"`
	JSON                                  struct {
		OfBetaManagedAgentsAgentToolConfigArray respjson.Field
		OfBetaManagedAgentsMCPToolConfigArray   respjson.Field
		raw                                     string
	} `json:"-"`
}

func (r *BetaManagedAgentsSessionThreadAgentToolUnionConfigs) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSessionThreadAgentToolUnionDefaultConfig is an implicit
// subunion of [BetaManagedAgentsSessionThreadAgentToolUnion].
// BetaManagedAgentsSessionThreadAgentToolUnionDefaultConfig provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsSessionThreadAgentToolUnion].
type BetaManagedAgentsSessionThreadAgentToolUnionDefaultConfig struct {
	Enabled bool `json:"enabled"`
	// This field is a union of
	// [BetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion],
	// [BetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion]
	PermissionPolicy BetaManagedAgentsSessionThreadAgentToolUnionDefaultConfigPermissionPolicy `json:"permission_policy"`
	JSON             struct {
		Enabled          respjson.Field
		PermissionPolicy respjson.Field
		raw              string
	} `json:"-"`
}

func (r *BetaManagedAgentsSessionThreadAgentToolUnionDefaultConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSessionThreadAgentToolUnionDefaultConfigPermissionPolicy is an
// implicit subunion of [BetaManagedAgentsSessionThreadAgentToolUnion].
// BetaManagedAgentsSessionThreadAgentToolUnionDefaultConfigPermissionPolicy
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsSessionThreadAgentToolUnion].
type BetaManagedAgentsSessionThreadAgentToolUnionDefaultConfigPermissionPolicy struct {
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

func (r *BetaManagedAgentsSessionThreadAgentToolUnionDefaultConfigPermissionPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSessionThreadAgentType string

const (
	BetaManagedAgentsSessionThreadAgentTypeAgent BetaManagedAgentsSessionThreadAgentType = "agent"
)

// Timing statistics for a session thread.
type BetaManagedAgentsSessionThreadStats struct {
	// Cumulative time in seconds the thread spent actively running. Excludes idle
	// time.
	ActiveSeconds float64 `json:"active_seconds"`
	// Elapsed time since thread creation in seconds. For archived threads, frozen at
	// the final update.
	DurationSeconds float64 `json:"duration_seconds"`
	// Time in seconds for the thread to begin running. Zero for child threads, which
	// start immediately.
	StartupSeconds float64 `json:"startup_seconds"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ActiveSeconds   respjson.Field
		DurationSeconds respjson.Field
		StartupSeconds  respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionThreadStats) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionThreadStats) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SessionThreadStatus enum
type BetaManagedAgentsSessionThreadStatus string

const (
	BetaManagedAgentsSessionThreadStatusRunning      BetaManagedAgentsSessionThreadStatus = "running"
	BetaManagedAgentsSessionThreadStatusIdle         BetaManagedAgentsSessionThreadStatus = "idle"
	BetaManagedAgentsSessionThreadStatusRescheduling BetaManagedAgentsSessionThreadStatus = "rescheduling"
	BetaManagedAgentsSessionThreadStatusTerminated   BetaManagedAgentsSessionThreadStatus = "terminated"
)

// Cumulative token usage for a session thread across all turns.
type BetaManagedAgentsSessionThreadUsage struct {
	// Prompt-cache creation token usage broken down by cache lifetime.
	CacheCreation BetaManagedAgentsCacheCreationUsage `json:"cache_creation"`
	// Total tokens read from prompt cache.
	CacheReadInputTokens int64 `json:"cache_read_input_tokens"`
	// Total input tokens consumed across all turns.
	InputTokens int64 `json:"input_tokens"`
	// Total output tokens generated across all turns.
	OutputTokens int64 `json:"output_tokens"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CacheCreation        respjson.Field
		CacheReadInputTokens respjson.Field
		InputTokens          respjson.Field
		OutputTokens         respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionThreadUsage) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionThreadUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsStreamSessionThreadEventsUnion contains all possible properties
// and values from [BetaManagedAgentsUserMessageEvent],
// [BetaManagedAgentsUserInterruptEvent],
// [BetaManagedAgentsUserToolConfirmationEvent],
// [BetaManagedAgentsUserCustomToolResultEvent],
// [BetaManagedAgentsAgentCustomToolUseEvent],
// [BetaManagedAgentsAgentMessageEvent], [BetaManagedAgentsAgentThinkingEvent],
// [BetaManagedAgentsAgentMCPToolUseEvent],
// [BetaManagedAgentsAgentMCPToolResultEvent],
// [BetaManagedAgentsAgentToolUseEvent], [BetaManagedAgentsAgentToolResultEvent],
// [BetaManagedAgentsAgentThreadMessageReceivedEvent],
// [BetaManagedAgentsAgentThreadMessageSentEvent],
// [BetaManagedAgentsAgentThreadContextCompactedEvent],
// [BetaManagedAgentsSessionErrorEvent],
// [BetaManagedAgentsSessionStatusRescheduledEvent],
// [BetaManagedAgentsSessionStatusRunningEvent],
// [BetaManagedAgentsSessionStatusIdleEvent],
// [BetaManagedAgentsSessionStatusTerminatedEvent],
// [BetaManagedAgentsSessionThreadCreatedEvent],
// [BetaManagedAgentsSpanOutcomeEvaluationStartEvent],
// [BetaManagedAgentsSpanOutcomeEvaluationEndEvent],
// [BetaManagedAgentsSpanModelRequestStartEvent],
// [BetaManagedAgentsSpanModelRequestEndEvent],
// [BetaManagedAgentsSpanOutcomeEvaluationOngoingEvent],
// [BetaManagedAgentsUserDefineOutcomeEvent],
// [BetaManagedAgentsSessionDeletedEvent],
// [BetaManagedAgentsSessionThreadStatusRunningEvent],
// [BetaManagedAgentsSessionThreadStatusIdleEvent],
// [BetaManagedAgentsSessionThreadStatusTerminatedEvent],
// [BetaManagedAgentsSessionThreadStatusRescheduledEvent].
//
// Use the [BetaManagedAgentsStreamSessionThreadEventsUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsStreamSessionThreadEventsUnion struct {
	ID string `json:"id"`
	// This field is a union of [[]BetaManagedAgentsUserMessageEventContentUnion],
	// [[]BetaManagedAgentsUserCustomToolResultEventContentUnion],
	// [[]BetaManagedAgentsTextBlock],
	// [[]BetaManagedAgentsAgentMCPToolResultEventContentUnion],
	// [[]BetaManagedAgentsAgentToolResultEventContentUnion],
	// [[]BetaManagedAgentsAgentThreadMessageReceivedEventContentUnion],
	// [[]BetaManagedAgentsAgentThreadMessageSentEventContentUnion]
	Content BetaManagedAgentsStreamSessionThreadEventsUnionContent `json:"content"`
	// Any of "user.message", "user.interrupt", "user.tool_confirmation",
	// "user.custom_tool_result", "agent.custom_tool_use", "agent.message",
	// "agent.thinking", "agent.mcp_tool_use", "agent.mcp_tool_result",
	// "agent.tool_use", "agent.tool_result", "agent.thread_message_received",
	// "agent.thread_message_sent", "agent.thread_context_compacted", "session.error",
	// "session.status_rescheduled", "session.status_running", "session.status_idle",
	// "session.status_terminated", "session.thread_created",
	// "span.outcome_evaluation_start", "span.outcome_evaluation_end",
	// "span.model_request_start", "span.model_request_end",
	// "span.outcome_evaluation_ongoing", "user.define_outcome", "session.deleted",
	// "session.thread_status_running", "session.thread_status_idle",
	// "session.thread_status_terminated", "session.thread_status_rescheduled".
	Type            string    `json:"type"`
	ProcessedAt     time.Time `json:"processed_at"`
	SessionThreadID string    `json:"session_thread_id"`
	Result          string    `json:"result"`
	ToolUseID       string    `json:"tool_use_id"`
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
	// This field is from variant [BetaManagedAgentsAgentThreadMessageReceivedEvent].
	FromSessionThreadID string `json:"from_session_thread_id"`
	// This field is from variant [BetaManagedAgentsAgentThreadMessageReceivedEvent].
	FromAgentName string `json:"from_agent_name"`
	// This field is from variant [BetaManagedAgentsAgentThreadMessageSentEvent].
	ToSessionThreadID string `json:"to_session_thread_id"`
	// This field is from variant [BetaManagedAgentsAgentThreadMessageSentEvent].
	ToAgentName string `json:"to_agent_name"`
	// This field is from variant [BetaManagedAgentsSessionErrorEvent].
	Error BetaManagedAgentsSessionErrorEventErrorUnion `json:"error"`
	// This field is a union of
	// [BetaManagedAgentsSessionStatusIdleEventStopReasonUnion],
	// [BetaManagedAgentsSessionThreadStatusIdleEventStopReasonUnion]
	StopReason BetaManagedAgentsStreamSessionThreadEventsUnionStopReason `json:"stop_reason"`
	AgentName  string                                                    `json:"agent_name"`
	Iteration  int64                                                     `json:"iteration"`
	OutcomeID  string                                                    `json:"outcome_id"`
	// This field is from variant [BetaManagedAgentsSpanOutcomeEvaluationEndEvent].
	Explanation string `json:"explanation"`
	// This field is from variant [BetaManagedAgentsSpanOutcomeEvaluationEndEvent].
	OutcomeEvaluationStartID string `json:"outcome_evaluation_start_id"`
	// This field is from variant [BetaManagedAgentsSpanOutcomeEvaluationEndEvent].
	Usage BetaManagedAgentsSpanModelUsage `json:"usage"`
	// This field is from variant [BetaManagedAgentsSpanModelRequestEndEvent].
	ModelRequestStartID string `json:"model_request_start_id"`
	// This field is from variant [BetaManagedAgentsSpanModelRequestEndEvent].
	ModelUsage BetaManagedAgentsSpanModelUsage `json:"model_usage"`
	// This field is from variant [BetaManagedAgentsUserDefineOutcomeEvent].
	Description string `json:"description"`
	// This field is from variant [BetaManagedAgentsUserDefineOutcomeEvent].
	MaxIterations int64 `json:"max_iterations"`
	// This field is from variant [BetaManagedAgentsUserDefineOutcomeEvent].
	Rubric BetaManagedAgentsUserDefineOutcomeEventRubricUnion `json:"rubric"`
	JSON   struct {
		ID                       respjson.Field
		Content                  respjson.Field
		Type                     respjson.Field
		ProcessedAt              respjson.Field
		SessionThreadID          respjson.Field
		Result                   respjson.Field
		ToolUseID                respjson.Field
		DenyMessage              respjson.Field
		CustomToolUseID          respjson.Field
		IsError                  respjson.Field
		Input                    respjson.Field
		Name                     respjson.Field
		MCPServerName            respjson.Field
		EvaluatedPermission      respjson.Field
		MCPToolUseID             respjson.Field
		FromSessionThreadID      respjson.Field
		FromAgentName            respjson.Field
		ToSessionThreadID        respjson.Field
		ToAgentName              respjson.Field
		Error                    respjson.Field
		StopReason               respjson.Field
		AgentName                respjson.Field
		Iteration                respjson.Field
		OutcomeID                respjson.Field
		Explanation              respjson.Field
		OutcomeEvaluationStartID respjson.Field
		Usage                    respjson.Field
		ModelRequestStartID      respjson.Field
		ModelUsage               respjson.Field
		Description              respjson.Field
		MaxIterations            respjson.Field
		Rubric                   respjson.Field
		raw                      string
	} `json:"-"`
}

// anyBetaManagedAgentsStreamSessionThreadEvents is implemented by each variant of
// [BetaManagedAgentsStreamSessionThreadEventsUnion] to add type safety for the
// return type of [BetaManagedAgentsStreamSessionThreadEventsUnion.AsAny]
type anyBetaManagedAgentsStreamSessionThreadEvents interface {
	implBetaManagedAgentsStreamSessionThreadEventsUnion()
}

func (BetaManagedAgentsUserMessageEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion()   {}
func (BetaManagedAgentsUserInterruptEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {}
func (BetaManagedAgentsUserToolConfirmationEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsUserCustomToolResultEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsAgentCustomToolUseEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsAgentMessageEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion()    {}
func (BetaManagedAgentsAgentThinkingEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion()   {}
func (BetaManagedAgentsAgentMCPToolUseEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {}
func (BetaManagedAgentsAgentMCPToolResultEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsAgentToolUseEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion()    {}
func (BetaManagedAgentsAgentToolResultEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {}
func (BetaManagedAgentsAgentThreadMessageReceivedEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsAgentThreadMessageSentEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsAgentThreadContextCompactedEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsSessionErrorEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {}
func (BetaManagedAgentsSessionStatusRescheduledEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsSessionStatusRunningEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsSessionStatusIdleEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsSessionStatusTerminatedEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsSessionThreadCreatedEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsSpanOutcomeEvaluationStartEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsSpanOutcomeEvaluationEndEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsSpanModelRequestStartEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsSpanModelRequestEndEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsSpanOutcomeEvaluationOngoingEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsUserDefineOutcomeEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsSessionDeletedEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {}
func (BetaManagedAgentsSessionThreadStatusRunningEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsSessionThreadStatusIdleEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsSessionThreadStatusTerminatedEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}
func (BetaManagedAgentsSessionThreadStatusRescheduledEvent) implBetaManagedAgentsStreamSessionThreadEventsUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsStreamSessionThreadEventsUnion.AsAny().(type) {
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
//	case anthropic.BetaManagedAgentsAgentThreadMessageReceivedEvent:
//	case anthropic.BetaManagedAgentsAgentThreadMessageSentEvent:
//	case anthropic.BetaManagedAgentsAgentThreadContextCompactedEvent:
//	case anthropic.BetaManagedAgentsSessionErrorEvent:
//	case anthropic.BetaManagedAgentsSessionStatusRescheduledEvent:
//	case anthropic.BetaManagedAgentsSessionStatusRunningEvent:
//	case anthropic.BetaManagedAgentsSessionStatusIdleEvent:
//	case anthropic.BetaManagedAgentsSessionStatusTerminatedEvent:
//	case anthropic.BetaManagedAgentsSessionThreadCreatedEvent:
//	case anthropic.BetaManagedAgentsSpanOutcomeEvaluationStartEvent:
//	case anthropic.BetaManagedAgentsSpanOutcomeEvaluationEndEvent:
//	case anthropic.BetaManagedAgentsSpanModelRequestStartEvent:
//	case anthropic.BetaManagedAgentsSpanModelRequestEndEvent:
//	case anthropic.BetaManagedAgentsSpanOutcomeEvaluationOngoingEvent:
//	case anthropic.BetaManagedAgentsUserDefineOutcomeEvent:
//	case anthropic.BetaManagedAgentsSessionDeletedEvent:
//	case anthropic.BetaManagedAgentsSessionThreadStatusRunningEvent:
//	case anthropic.BetaManagedAgentsSessionThreadStatusIdleEvent:
//	case anthropic.BetaManagedAgentsSessionThreadStatusTerminatedEvent:
//	case anthropic.BetaManagedAgentsSessionThreadStatusRescheduledEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsAny() anyBetaManagedAgentsStreamSessionThreadEvents {
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
	case "agent.thread_message_received":
		return u.AsAgentThreadMessageReceived()
	case "agent.thread_message_sent":
		return u.AsAgentThreadMessageSent()
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
	case "session.thread_created":
		return u.AsSessionThreadCreated()
	case "span.outcome_evaluation_start":
		return u.AsSpanOutcomeEvaluationStart()
	case "span.outcome_evaluation_end":
		return u.AsSpanOutcomeEvaluationEnd()
	case "span.model_request_start":
		return u.AsSpanModelRequestStart()
	case "span.model_request_end":
		return u.AsSpanModelRequestEnd()
	case "span.outcome_evaluation_ongoing":
		return u.AsSpanOutcomeEvaluationOngoing()
	case "user.define_outcome":
		return u.AsUserDefineOutcome()
	case "session.deleted":
		return u.AsSessionDeleted()
	case "session.thread_status_running":
		return u.AsSessionThreadStatusRunning()
	case "session.thread_status_idle":
		return u.AsSessionThreadStatusIdle()
	case "session.thread_status_terminated":
		return u.AsSessionThreadStatusTerminated()
	case "session.thread_status_rescheduled":
		return u.AsSessionThreadStatusRescheduled()
	}
	return nil
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsUserMessage() (v BetaManagedAgentsUserMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsUserInterrupt() (v BetaManagedAgentsUserInterruptEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsUserToolConfirmation() (v BetaManagedAgentsUserToolConfirmationEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsUserCustomToolResult() (v BetaManagedAgentsUserCustomToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsAgentCustomToolUse() (v BetaManagedAgentsAgentCustomToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsAgentMessage() (v BetaManagedAgentsAgentMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsAgentThinking() (v BetaManagedAgentsAgentThinkingEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsAgentMCPToolUse() (v BetaManagedAgentsAgentMCPToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsAgentMCPToolResult() (v BetaManagedAgentsAgentMCPToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsAgentToolUse() (v BetaManagedAgentsAgentToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsAgentToolResult() (v BetaManagedAgentsAgentToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsAgentThreadMessageReceived() (v BetaManagedAgentsAgentThreadMessageReceivedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsAgentThreadMessageSent() (v BetaManagedAgentsAgentThreadMessageSentEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsAgentThreadContextCompacted() (v BetaManagedAgentsAgentThreadContextCompactedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsSessionError() (v BetaManagedAgentsSessionErrorEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsSessionStatusRescheduled() (v BetaManagedAgentsSessionStatusRescheduledEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsSessionStatusRunning() (v BetaManagedAgentsSessionStatusRunningEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsSessionStatusIdle() (v BetaManagedAgentsSessionStatusIdleEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsSessionStatusTerminated() (v BetaManagedAgentsSessionStatusTerminatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsSessionThreadCreated() (v BetaManagedAgentsSessionThreadCreatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsSpanOutcomeEvaluationStart() (v BetaManagedAgentsSpanOutcomeEvaluationStartEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsSpanOutcomeEvaluationEnd() (v BetaManagedAgentsSpanOutcomeEvaluationEndEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsSpanModelRequestStart() (v BetaManagedAgentsSpanModelRequestStartEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsSpanModelRequestEnd() (v BetaManagedAgentsSpanModelRequestEndEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsSpanOutcomeEvaluationOngoing() (v BetaManagedAgentsSpanOutcomeEvaluationOngoingEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsUserDefineOutcome() (v BetaManagedAgentsUserDefineOutcomeEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsSessionDeleted() (v BetaManagedAgentsSessionDeletedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsSessionThreadStatusRunning() (v BetaManagedAgentsSessionThreadStatusRunningEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsSessionThreadStatusIdle() (v BetaManagedAgentsSessionThreadStatusIdleEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsSessionThreadStatusTerminated() (v BetaManagedAgentsSessionThreadStatusTerminatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsStreamSessionThreadEventsUnion) AsSessionThreadStatusRescheduled() (v BetaManagedAgentsSessionThreadStatusRescheduledEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsStreamSessionThreadEventsUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsStreamSessionThreadEventsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsStreamSessionThreadEventsUnionContent is an implicit subunion
// of [BetaManagedAgentsStreamSessionThreadEventsUnion].
// BetaManagedAgentsStreamSessionThreadEventsUnionContent provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsStreamSessionThreadEventsUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfBetaManagedAgentsUserMessageEventContentArray
// OfBetaManagedAgentsUserCustomToolResultEventContentArray
// OfBetaManagedAgentsTextBlockArray
// OfBetaManagedAgentsAgentMCPToolResultEventContentArray
// OfBetaManagedAgentsAgentToolResultEventContentArray
// OfBetaManagedAgentsAgentThreadMessageReceivedEventContentArray
// OfBetaManagedAgentsAgentThreadMessageSentEventContentArray]
type BetaManagedAgentsStreamSessionThreadEventsUnionContent struct {
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
	// This field will be present if the value is a
	// [[]BetaManagedAgentsAgentThreadMessageReceivedEventContentUnion] instead of an
	// object.
	OfBetaManagedAgentsAgentThreadMessageReceivedEventContentArray []BetaManagedAgentsAgentThreadMessageReceivedEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]BetaManagedAgentsAgentThreadMessageSentEventContentUnion] instead of an
	// object.
	OfBetaManagedAgentsAgentThreadMessageSentEventContentArray []BetaManagedAgentsAgentThreadMessageSentEventContentUnion `json:",inline"`
	JSON                                                       struct {
		OfBetaManagedAgentsUserMessageEventContentArray                respjson.Field
		OfBetaManagedAgentsUserCustomToolResultEventContentArray       respjson.Field
		OfBetaManagedAgentsTextBlockArray                              respjson.Field
		OfBetaManagedAgentsAgentMCPToolResultEventContentArray         respjson.Field
		OfBetaManagedAgentsAgentToolResultEventContentArray            respjson.Field
		OfBetaManagedAgentsAgentThreadMessageReceivedEventContentArray respjson.Field
		OfBetaManagedAgentsAgentThreadMessageSentEventContentArray     respjson.Field
		raw                                                            string
	} `json:"-"`
}

func (r *BetaManagedAgentsStreamSessionThreadEventsUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsStreamSessionThreadEventsUnionStopReason is an implicit
// subunion of [BetaManagedAgentsStreamSessionThreadEventsUnion].
// BetaManagedAgentsStreamSessionThreadEventsUnionStopReason provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsStreamSessionThreadEventsUnion].
type BetaManagedAgentsStreamSessionThreadEventsUnionStopReason struct {
	Type string `json:"type"`
	// This field is from variant
	// [BetaManagedAgentsSessionStatusIdleEventStopReasonUnion],
	// [BetaManagedAgentsSessionThreadStatusIdleEventStopReasonUnion].
	EventIDs []string `json:"event_ids"`
	JSON     struct {
		Type     respjson.Field
		EventIDs respjson.Field
		raw      string
	} `json:"-"`
}

func (r *BetaManagedAgentsStreamSessionThreadEventsUnionStopReason) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaSessionThreadGetParams struct {
	SessionID string `path:"session_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaSessionThreadListParams struct {
	// Maximum results per page. Defaults to 1000.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor from a previous response's next_page. Forward-only.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaSessionThreadListParams]'s query parameters as
// `url.Values`.
func (r BetaSessionThreadListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaSessionThreadArchiveParams struct {
	SessionID string `path:"session_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
