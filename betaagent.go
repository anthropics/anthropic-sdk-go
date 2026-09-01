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
	"github.com/anthropics/anthropic-sdk-go/internal/paramutil"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// BetaAgentService contains methods and other services that help with interacting
// with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaAgentService] method instead.
type BetaAgentService struct {
	Options  []option.RequestOption
	Versions BetaAgentVersionService
}

// NewBetaAgentService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaAgentService(opts ...option.RequestOption) (r BetaAgentService) {
	r = BetaAgentService{}
	r.Options = opts
	r.Versions = NewBetaAgentVersionService(opts...)
	return
}

// Create Agent
func (r *BetaAgentService) New(ctx context.Context, params BetaAgentNewParams, opts ...option.RequestOption) (res *BetaManagedAgentsAgent, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	path := "v1/agents?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get Agent
func (r *BetaAgentService) Get(ctx context.Context, agentID string, params BetaAgentGetParams, opts ...option.RequestOption) (res *BetaManagedAgentsAgent, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if agentID == "" {
		err = errors.New("missing required agent_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/agents/%s?beta=true", agentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Update Agent
func (r *BetaAgentService) Update(ctx context.Context, agentID string, params BetaAgentUpdateParams, opts ...option.RequestOption) (res *BetaManagedAgentsAgent, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if agentID == "" {
		err = errors.New("missing required agent_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/agents/%s?beta=true", agentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List Agents
func (r *BetaAgentService) List(ctx context.Context, params BetaAgentListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaManagedAgentsAgent], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/agents?beta=true"
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

// List Agents
func (r *BetaAgentService) ListAutoPaging(ctx context.Context, params BetaAgentListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsAgent] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Archive Agent
func (r *BetaAgentService) Archive(ctx context.Context, agentID string, body BetaAgentArchiveParams, opts ...option.RequestOption) (res *BetaManagedAgentsAgent, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if agentID == "" {
		err = errors.New("missing required agent_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/agents/%s/archive?beta=true", agentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Platform advisor roster entry: a model the session's primary thread may consult
// mid-turn.
type BetaManagedAgentsAdvisor struct {
	// The advisor model id.
	Model string `json:"model" api:"required"`
	// Any of "advisor".
	Type BetaManagedAgentsAdvisorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Model       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAdvisor) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAdvisor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAdvisorType string

const (
	BetaManagedAgentsAdvisorTypeAdvisor BetaManagedAgentsAdvisorType = "advisor"
)

// A Managed Agents `agent`.
type BetaManagedAgentsAgent struct {
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	CreatedAt   time.Time                                 `json:"created_at" api:"required" format:"date-time"`
	Description string                                    `json:"description" api:"required"`
	MCPServers  []BetaManagedAgentsMCPServerURLDefinition `json:"mcp_servers" api:"required"`
	Metadata    map[string]string                         `json:"metadata" api:"required"`
	// Model identifier and configuration.
	Model BetaManagedAgentsModelConfig `json:"model" api:"required"`
	// Resolved coordinator topology with a concrete agent roster.
	Multiagent BetaManagedAgentsMultiagent        `json:"multiagent" api:"required"`
	Name       string                             `json:"name" api:"required"`
	Skills     []BetaManagedAgentsAgentSkillUnion `json:"skills" api:"required"`
	System     string                             `json:"system" api:"required"`
	Tools      []BetaManagedAgentsAgentToolUnion  `json:"tools" api:"required"`
	// Any of "agent".
	Type BetaManagedAgentsAgentType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// The agent's current version. Starts at 1 and increments when the agent is
	// modified.
	Version int64 `json:"version" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ArchivedAt  respjson.Field
		CreatedAt   respjson.Field
		Description respjson.Field
		MCPServers  respjson.Field
		Metadata    respjson.Field
		Model       respjson.Field
		Multiagent  respjson.Field
		Name        respjson.Field
		Skills      respjson.Field
		System      respjson.Field
		Tools       respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsAgentSkillUnion contains all possible properties and values
// from [BetaManagedAgentsAnthropicSkill], [BetaManagedAgentsCustomSkill].
//
// Use the [BetaManagedAgentsAgentSkillUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsAgentSkillUnion struct {
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

// anyBetaManagedAgentsAgentSkill is implemented by each variant of
// [BetaManagedAgentsAgentSkillUnion] to add type safety for the return type of
// [BetaManagedAgentsAgentSkillUnion.AsAny]
type anyBetaManagedAgentsAgentSkill interface {
	implBetaManagedAgentsAgentSkillUnion()
}

func (BetaManagedAgentsAnthropicSkill) implBetaManagedAgentsAgentSkillUnion() {}
func (BetaManagedAgentsCustomSkill) implBetaManagedAgentsAgentSkillUnion()    {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsAgentSkillUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAnthropicSkill:
//	case anthropic.BetaManagedAgentsCustomSkill:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsAgentSkillUnion) AsAny() anyBetaManagedAgentsAgentSkill {
	switch u.Type {
	case "anthropic":
		return u.AsAnthropic()
	case "custom":
		return u.AsCustom()
	}
	return nil
}

func (u BetaManagedAgentsAgentSkillUnion) AsAnthropic() (v BetaManagedAgentsAnthropicSkill) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsAgentSkillUnion) AsCustom() (v BetaManagedAgentsCustomSkill) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsAgentSkillUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsAgentSkillUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsAgentToolUnion contains all possible properties and values from
// [BetaManagedAgentsAgentToolset20260401], [BetaManagedAgentsMCPToolset],
// [BetaManagedAgentsCustomTool].
//
// Use the [BetaManagedAgentsAgentToolUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsAgentToolUnion struct {
	// This field is a union of [[]BetaManagedAgentsAgentToolConfigUnion],
	// [[]BetaManagedAgentsMCPToolConfig]
	Configs BetaManagedAgentsAgentToolUnionConfigs `json:"configs"`
	// This field is a union of [BetaManagedAgentsAgentToolsetDefaultConfig],
	// [BetaManagedAgentsMCPToolsetDefaultConfig]
	DefaultConfig BetaManagedAgentsAgentToolUnionDefaultConfig `json:"default_config"`
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

// anyBetaManagedAgentsAgentTool is implemented by each variant of
// [BetaManagedAgentsAgentToolUnion] to add type safety for the return type of
// [BetaManagedAgentsAgentToolUnion.AsAny]
type anyBetaManagedAgentsAgentTool interface {
	implBetaManagedAgentsAgentToolUnion()
}

func (BetaManagedAgentsAgentToolset20260401) implBetaManagedAgentsAgentToolUnion() {}
func (BetaManagedAgentsMCPToolset) implBetaManagedAgentsAgentToolUnion()           {}
func (BetaManagedAgentsCustomTool) implBetaManagedAgentsAgentToolUnion()           {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsAgentToolUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAgentToolset20260401:
//	case anthropic.BetaManagedAgentsMCPToolset:
//	case anthropic.BetaManagedAgentsCustomTool:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsAgentToolUnion) AsAny() anyBetaManagedAgentsAgentTool {
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

func (u BetaManagedAgentsAgentToolUnion) AsAgentToolset20260401() (v BetaManagedAgentsAgentToolset20260401) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsAgentToolUnion) AsMCPToolset() (v BetaManagedAgentsMCPToolset) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsAgentToolUnion) AsCustom() (v BetaManagedAgentsCustomTool) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsAgentToolUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsAgentToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsAgentToolUnionConfigs is an implicit subunion of
// [BetaManagedAgentsAgentToolUnion]. BetaManagedAgentsAgentToolUnionConfigs
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsAgentToolUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfBetaManagedAgentsAgentToolConfigArray
// OfBetaManagedAgentsMCPToolConfigArray]
type BetaManagedAgentsAgentToolUnionConfigs struct {
	// This field will be present if the value is a
	// [[]BetaManagedAgentsAgentToolConfigUnion] instead of an object.
	OfBetaManagedAgentsAgentToolConfigArray []BetaManagedAgentsAgentToolConfigUnion `json:",inline"`
	// This field will be present if the value is a [[]BetaManagedAgentsMCPToolConfig]
	// instead of an object.
	OfBetaManagedAgentsMCPToolConfigArray []BetaManagedAgentsMCPToolConfig `json:",inline"`
	JSON                                  struct {
		OfBetaManagedAgentsAgentToolConfigArray respjson.Field
		OfBetaManagedAgentsMCPToolConfigArray   respjson.Field
		raw                                     string
	} `json:"-"`
}

func (r *BetaManagedAgentsAgentToolUnionConfigs) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsAgentToolUnionDefaultConfig is an implicit subunion of
// [BetaManagedAgentsAgentToolUnion]. BetaManagedAgentsAgentToolUnionDefaultConfig
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsAgentToolUnion].
type BetaManagedAgentsAgentToolUnionDefaultConfig struct {
	Enabled bool `json:"enabled"`
	// This field is a union of
	// [BetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion],
	// [BetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion]
	PermissionPolicy BetaManagedAgentsAgentToolUnionDefaultConfigPermissionPolicy `json:"permission_policy"`
	JSON             struct {
		Enabled          respjson.Field
		PermissionPolicy respjson.Field
		raw              string
	} `json:"-"`
}

func (r *BetaManagedAgentsAgentToolUnionDefaultConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsAgentToolUnionDefaultConfigPermissionPolicy is an implicit
// subunion of [BetaManagedAgentsAgentToolUnion].
// BetaManagedAgentsAgentToolUnionDefaultConfigPermissionPolicy provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsAgentToolUnion].
type BetaManagedAgentsAgentToolUnionDefaultConfigPermissionPolicy struct {
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

func (r *BetaManagedAgentsAgentToolUnionDefaultConfigPermissionPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAgentType string

const (
	BetaManagedAgentsAgentTypeAgent BetaManagedAgentsAgentType = "agent"
)

// A resolved agent reference with a concrete version.
type BetaManagedAgentsAgentReference struct {
	ID string `json:"id" api:"required"`
	// Any of "agent".
	Type    BetaManagedAgentsAgentReferenceType `json:"type" api:"required"`
	Version int64                               `json:"version" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentReference) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentReference) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAgentReferenceType string

const (
	BetaManagedAgentsAgentReferenceTypeAgent BetaManagedAgentsAgentReferenceType = "agent"
)

// BetaManagedAgentsAgentToolConfigUnion contains all possible properties and
// values from [BetaManagedAgentsBashToolConfig],
// [BetaManagedAgentsEditToolConfig], [BetaManagedAgentsReadToolConfig],
// [BetaManagedAgentsWriteToolConfig], [BetaManagedAgentsGlobToolConfig],
// [BetaManagedAgentsGrepToolConfig], [BetaManagedAgentsWebFetchToolConfig],
// [BetaManagedAgentsWebSearchToolConfig].
//
// Use the [BetaManagedAgentsAgentToolConfigUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsAgentToolConfigUnion struct {
	Enabled bool   `json:"enabled"`
	Name    string `json:"name"`
	// This field is a union of [BetaManagedAgentsBashToolConfigPermissionPolicyUnion],
	// [BetaManagedAgentsEditToolConfigPermissionPolicyUnion],
	// [BetaManagedAgentsReadToolConfigPermissionPolicyUnion],
	// [BetaManagedAgentsWriteToolConfigPermissionPolicyUnion],
	// [BetaManagedAgentsGlobToolConfigPermissionPolicyUnion],
	// [BetaManagedAgentsGrepToolConfigPermissionPolicyUnion],
	// [BetaManagedAgentsWebFetchToolConfigPermissionPolicyUnion],
	// [BetaManagedAgentsWebSearchToolConfigPermissionPolicyUnion]
	PermissionPolicy BetaManagedAgentsAgentToolConfigUnionPermissionPolicy `json:"permission_policy"`
	// Any of "bash", "edit", "read", "write", "glob", "grep", "web_fetch",
	// "web_search".
	Type           string   `json:"type"`
	AllowedDomains []string `json:"allowed_domains"`
	BlockedDomains []string `json:"blocked_domains"`
	// This field is from variant [BetaManagedAgentsWebFetchToolConfig].
	MaxContentTokens int64 `json:"max_content_tokens"`
	// This field is from variant [BetaManagedAgentsWebSearchToolConfig].
	UserLocation BetaManagedAgentsUserLocation `json:"user_location"`
	JSON         struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		AllowedDomains   respjson.Field
		BlockedDomains   respjson.Field
		MaxContentTokens respjson.Field
		UserLocation     respjson.Field
		raw              string
	} `json:"-"`
}

// anyBetaManagedAgentsAgentToolConfig is implemented by each variant of
// [BetaManagedAgentsAgentToolConfigUnion] to add type safety for the return type
// of [BetaManagedAgentsAgentToolConfigUnion.AsAny]
type anyBetaManagedAgentsAgentToolConfig interface {
	implBetaManagedAgentsAgentToolConfigUnion()
}

func (BetaManagedAgentsBashToolConfig) implBetaManagedAgentsAgentToolConfigUnion()      {}
func (BetaManagedAgentsEditToolConfig) implBetaManagedAgentsAgentToolConfigUnion()      {}
func (BetaManagedAgentsReadToolConfig) implBetaManagedAgentsAgentToolConfigUnion()      {}
func (BetaManagedAgentsWriteToolConfig) implBetaManagedAgentsAgentToolConfigUnion()     {}
func (BetaManagedAgentsGlobToolConfig) implBetaManagedAgentsAgentToolConfigUnion()      {}
func (BetaManagedAgentsGrepToolConfig) implBetaManagedAgentsAgentToolConfigUnion()      {}
func (BetaManagedAgentsWebFetchToolConfig) implBetaManagedAgentsAgentToolConfigUnion()  {}
func (BetaManagedAgentsWebSearchToolConfig) implBetaManagedAgentsAgentToolConfigUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsAgentToolConfigUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsBashToolConfig:
//	case anthropic.BetaManagedAgentsEditToolConfig:
//	case anthropic.BetaManagedAgentsReadToolConfig:
//	case anthropic.BetaManagedAgentsWriteToolConfig:
//	case anthropic.BetaManagedAgentsGlobToolConfig:
//	case anthropic.BetaManagedAgentsGrepToolConfig:
//	case anthropic.BetaManagedAgentsWebFetchToolConfig:
//	case anthropic.BetaManagedAgentsWebSearchToolConfig:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsAgentToolConfigUnion) AsAny() anyBetaManagedAgentsAgentToolConfig {
	switch u.Type {
	case "bash":
		return u.AsBash()
	case "edit":
		return u.AsEdit()
	case "read":
		return u.AsRead()
	case "write":
		return u.AsWrite()
	case "glob":
		return u.AsGlob()
	case "grep":
		return u.AsGrep()
	case "web_fetch":
		return u.AsWebFetch()
	case "web_search":
		return u.AsWebSearch()
	}
	return nil
}

func (u BetaManagedAgentsAgentToolConfigUnion) AsBash() (v BetaManagedAgentsBashToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsAgentToolConfigUnion) AsEdit() (v BetaManagedAgentsEditToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsAgentToolConfigUnion) AsRead() (v BetaManagedAgentsReadToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsAgentToolConfigUnion) AsWrite() (v BetaManagedAgentsWriteToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsAgentToolConfigUnion) AsGlob() (v BetaManagedAgentsGlobToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsAgentToolConfigUnion) AsGrep() (v BetaManagedAgentsGrepToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsAgentToolConfigUnion) AsWebFetch() (v BetaManagedAgentsWebFetchToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsAgentToolConfigUnion) AsWebSearch() (v BetaManagedAgentsWebSearchToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsAgentToolConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsAgentToolConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsAgentToolConfigUnionPermissionPolicy is an implicit subunion of
// [BetaManagedAgentsAgentToolConfigUnion].
// BetaManagedAgentsAgentToolConfigUnionPermissionPolicy provides convenient access
// to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsAgentToolConfigUnion].
type BetaManagedAgentsAgentToolConfigUnionPermissionPolicy struct {
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

func (r *BetaManagedAgentsAgentToolConfigUnionPermissionPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsAgentToolConfigParamsUnion struct {
	OfBash      *BetaManagedAgentsBashToolConfigParams      `json:",omitzero,inline"`
	OfEdit      *BetaManagedAgentsEditToolConfigParams      `json:",omitzero,inline"`
	OfRead      *BetaManagedAgentsReadToolConfigParams      `json:",omitzero,inline"`
	OfWrite     *BetaManagedAgentsWriteToolConfigParams     `json:",omitzero,inline"`
	OfGlob      *BetaManagedAgentsGlobToolConfigParams      `json:",omitzero,inline"`
	OfGrep      *BetaManagedAgentsGrepToolConfigParams      `json:",omitzero,inline"`
	OfWebFetch  *BetaManagedAgentsWebFetchToolConfigParams  `json:",omitzero,inline"`
	OfWebSearch *BetaManagedAgentsWebSearchToolConfigParams `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsAgentToolConfigParamsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfBash,
		u.OfEdit,
		u.OfRead,
		u.OfWrite,
		u.OfGlob,
		u.OfGrep,
		u.OfWebFetch,
		u.OfWebSearch)
}
func (u *BetaManagedAgentsAgentToolConfigParamsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsAgentToolConfigParamsUnion) asAny() any {
	if !param.IsOmitted(u.OfBash) {
		return u.OfBash
	} else if !param.IsOmitted(u.OfEdit) {
		return u.OfEdit
	} else if !param.IsOmitted(u.OfRead) {
		return u.OfRead
	} else if !param.IsOmitted(u.OfWrite) {
		return u.OfWrite
	} else if !param.IsOmitted(u.OfGlob) {
		return u.OfGlob
	} else if !param.IsOmitted(u.OfGrep) {
		return u.OfGrep
	} else if !param.IsOmitted(u.OfWebFetch) {
		return u.OfWebFetch
	} else if !param.IsOmitted(u.OfWebSearch) {
		return u.OfWebSearch
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsAgentToolConfigParamsUnion) GetMaxContentTokens() *int64 {
	if vt := u.OfWebFetch; vt != nil && vt.MaxContentTokens.Valid() {
		return &vt.MaxContentTokens.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsAgentToolConfigParamsUnion) GetUserLocation() *BetaManagedAgentsUserLocationParam {
	if vt := u.OfWebSearch; vt != nil {
		return &vt.UserLocation
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsAgentToolConfigParamsUnion) GetName() *string {
	if vt := u.OfBash; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfEdit; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfRead; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfWrite; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfGlob; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfGrep; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfWebFetch; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfWebSearch; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsAgentToolConfigParamsUnion) GetEnabled() *bool {
	if vt := u.OfBash; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	} else if vt := u.OfEdit; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	} else if vt := u.OfRead; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	} else if vt := u.OfWrite; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	} else if vt := u.OfGlob; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	} else if vt := u.OfGrep; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	} else if vt := u.OfWebFetch; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	} else if vt := u.OfWebSearch; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsAgentToolConfigParamsUnion) GetType() *string {
	if vt := u.OfBash; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfEdit; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRead; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWrite; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfGlob; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfGrep; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWebFetch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWebSearch; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u BetaManagedAgentsAgentToolConfigParamsUnion) GetPermissionPolicy() (res betaManagedAgentsAgentToolConfigParamsUnionPermissionPolicy) {
	if vt := u.OfBash; vt != nil {
		res.any = vt.PermissionPolicy.asAny()
	} else if vt := u.OfEdit; vt != nil {
		res.any = vt.PermissionPolicy.asAny()
	} else if vt := u.OfRead; vt != nil {
		res.any = vt.PermissionPolicy.asAny()
	} else if vt := u.OfWrite; vt != nil {
		res.any = vt.PermissionPolicy.asAny()
	} else if vt := u.OfGlob; vt != nil {
		res.any = vt.PermissionPolicy.asAny()
	} else if vt := u.OfGrep; vt != nil {
		res.any = vt.PermissionPolicy.asAny()
	} else if vt := u.OfWebFetch; vt != nil {
		res.any = vt.PermissionPolicy.asAny()
	} else if vt := u.OfWebSearch; vt != nil {
		res.any = vt.PermissionPolicy.asAny()
	}
	return
}

// Can have the runtime types [*BetaManagedAgentsAlwaysAllowPolicyParam],
// [*BetaManagedAgentsAlwaysAskPolicyParam]
type betaManagedAgentsAgentToolConfigParamsUnionPermissionPolicy struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *anthropic.BetaManagedAgentsAlwaysAllowPolicyParam:
//	case *anthropic.BetaManagedAgentsAlwaysAskPolicyParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u betaManagedAgentsAgentToolConfigParamsUnionPermissionPolicy) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u betaManagedAgentsAgentToolConfigParamsUnionPermissionPolicy) GetType() *string {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsBashToolConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *BetaManagedAgentsEditToolConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *BetaManagedAgentsReadToolConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *BetaManagedAgentsWriteToolConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *BetaManagedAgentsGlobToolConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *BetaManagedAgentsGrepToolConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *BetaManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *BetaManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	}
	return nil
}

// Returns a pointer to the underlying variant's AllowedDomains property, if
// present.
func (u BetaManagedAgentsAgentToolConfigParamsUnion) GetAllowedDomains() []string {
	if vt := u.OfWebFetch; vt != nil {
		return vt.AllowedDomains
	} else if vt := u.OfWebSearch; vt != nil {
		return vt.AllowedDomains
	}
	return nil
}

// Returns a pointer to the underlying variant's BlockedDomains property, if
// present.
func (u BetaManagedAgentsAgentToolConfigParamsUnion) GetBlockedDomains() []string {
	if vt := u.OfWebFetch; vt != nil {
		return vt.BlockedDomains
	} else if vt := u.OfWebSearch; vt != nil {
		return vt.BlockedDomains
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsAgentToolConfigParamsUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsBashToolConfigParams]("bash"),
		apijson.Discriminator[BetaManagedAgentsEditToolConfigParams]("edit"),
		apijson.Discriminator[BetaManagedAgentsReadToolConfigParams]("read"),
		apijson.Discriminator[BetaManagedAgentsWriteToolConfigParams]("write"),
		apijson.Discriminator[BetaManagedAgentsGlobToolConfigParams]("glob"),
		apijson.Discriminator[BetaManagedAgentsGrepToolConfigParams]("grep"),
		apijson.Discriminator[BetaManagedAgentsWebFetchToolConfigParams]("web_fetch"),
		apijson.Discriminator[BetaManagedAgentsWebSearchToolConfigParams]("web_search"),
	)
}

// Resolved default configuration for agent tools.
type BetaManagedAgentsAgentToolsetDefaultConfig struct {
	Enabled bool `json:"enabled" api:"required"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		PermissionPolicy respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentToolsetDefaultConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentToolsetDefaultConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion contains all
// possible properties and values from [BetaManagedAgentsAlwaysAllowPolicy],
// [BetaManagedAgentsAlwaysAskPolicy].
//
// Use the [BetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicy is implemented by
// each variant of
// [BetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion] to add type
// safety for the return type of
// [BetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion.AsAny]
type anyBetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicy interface {
	implBetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion()
}

func (BetaManagedAgentsAlwaysAllowPolicy) implBetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion() {
}
func (BetaManagedAgentsAlwaysAskPolicy) implBetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAlwaysAllowPolicy:
//	case anthropic.BetaManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion) AsAny() anyBetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u BetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion) AsAlwaysAllow() (v BetaManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion) AsAlwaysAsk() (v BetaManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *BetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Default configuration for all tools in a toolset.
type BetaManagedAgentsAgentToolsetDefaultConfigParams struct {
	// Whether tools are enabled and available to Claude by default. Defaults to true
	// if not specified.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	paramObj
}

func (r BetaManagedAgentsAgentToolsetDefaultConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsAgentToolsetDefaultConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsAgentToolsetDefaultConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *BetaManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *BetaManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *BetaManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[BetaManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type BetaManagedAgentsAgentToolset20260401 struct {
	Configs []BetaManagedAgentsAgentToolConfigUnion `json:"configs" api:"required"`
	// Resolved default configuration for agent tools.
	DefaultConfig BetaManagedAgentsAgentToolsetDefaultConfig `json:"default_config" api:"required"`
	// Any of "agent_toolset_20260401".
	Type BetaManagedAgentsAgentToolset20260401Type `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Configs       respjson.Field
		DefaultConfig respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentToolset20260401) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentToolset20260401) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAgentToolset20260401Type string

const (
	BetaManagedAgentsAgentToolset20260401TypeAgentToolset20260401 BetaManagedAgentsAgentToolset20260401Type = "agent_toolset_20260401"
)

// Input payload for the `bash` tool of the `agent_toolset_20260401` toolset. All
// fields are optional; a normal invocation supplies `command`, while
// `restart=true` (with no `command`) reboots the runner-side bash session.
type BetaManagedAgentsAgentToolset20260401BashInput struct {
	// Shell command to execute. Omit only when `restart` is true.
	Command string `json:"command"`
	// When true, restart the persistent bash session instead of running a command.
	// Subsequent calls without `restart` will run against the fresh session.
	Restart bool `json:"restart"`
	// Per-call timeout in milliseconds. Defaults to the runner-wide tool timeout when
	// omitted or zero.
	TimeoutMs int64 `json:"timeout_ms"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Command     respjson.Field
		Restart     respjson.Field
		TimeoutMs   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentToolset20260401BashInput) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentToolset20260401BashInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Input payload for the `edit` tool. Performs a string replacement in the named
// file; by default `old_string` must occur exactly once.
type BetaManagedAgentsAgentToolset20260401EditInput struct {
	// Path of the file to edit.
	FilePath string `json:"file_path" api:"required"`
	// Replacement text.
	NewString string `json:"new_string" api:"required"`
	// Substring to find and replace.
	OldString string `json:"old_string" api:"required"`
	// When true, replace every occurrence of `old_string` instead of requiring a
	// unique match.
	ReplaceAll bool `json:"replace_all"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FilePath    respjson.Field
		NewString   respjson.Field
		OldString   respjson.Field
		ReplaceAll  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentToolset20260401EditInput) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentToolset20260401EditInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Input payload for the `glob` tool. Returns paths matching a doublestar glob
// pattern, newest first.
type BetaManagedAgentsAgentToolset20260401GlobInput struct {
	// Doublestar glob pattern (e.g. `**/*.go`). Absolute patterns are only permitted
	// when the runner is configured to allow them.
	Pattern string `json:"pattern" api:"required"`
	// Optional directory root to search under. Defaults to the runner's working
	// directory.
	Path string `json:"path"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Pattern     respjson.Field
		Path        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentToolset20260401GlobInput) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentToolset20260401GlobInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Input payload for the `grep` tool. Searches file contents for a regular
// expression, returning matching lines.
type BetaManagedAgentsAgentToolset20260401GrepInput struct {
	// Regular expression to search for.
	Pattern string `json:"pattern" api:"required"`
	// Optional directory root to search under. Defaults to the runner's working
	// directory.
	Path string `json:"path"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Pattern     respjson.Field
		Path        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentToolset20260401GrepInput) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentToolset20260401GrepInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for built-in agent tools. Use this to enable or disable groups of
// tools available to the agent.
//
// The property Type is required.
type BetaManagedAgentsAgentToolset20260401Params struct {
	// Any of "agent_toolset_20260401".
	Type BetaManagedAgentsAgentToolset20260401ParamsType `json:"type,omitzero" api:"required"`
	// Per-tool configuration overrides.
	Configs []BetaManagedAgentsAgentToolConfigParamsUnion `json:"configs,omitzero"`
	// Default configuration for all tools in a toolset.
	DefaultConfig BetaManagedAgentsAgentToolsetDefaultConfigParams `json:"default_config,omitzero"`
	paramObj
}

func (r BetaManagedAgentsAgentToolset20260401Params) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsAgentToolset20260401Params
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsAgentToolset20260401Params) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAgentToolset20260401ParamsType string

const (
	BetaManagedAgentsAgentToolset20260401ParamsTypeAgentToolset20260401 BetaManagedAgentsAgentToolset20260401ParamsType = "agent_toolset_20260401"
)

// Input payload for the `read` tool. Reads file contents relative to the runner's
// working directory (or absolute when the runner permits).
type BetaManagedAgentsAgentToolset20260401ReadInput struct {
	// Path of the file to read.
	FilePath string `json:"file_path" api:"required"`
	// Optional `[start_line, end_line]` 1-indexed inclusive range. When omitted the
	// entire file is returned. `end_line` of 0 or negative means "to end of file".
	ViewRange []int64 `json:"view_range"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FilePath    respjson.Field
		ViewRange   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentToolset20260401ReadInput) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentToolset20260401ReadInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Input payload for the `write` tool. Writes (overwriting) the entire file
// contents.
type BetaManagedAgentsAgentToolset20260401WriteInput struct {
	// Full file contents to write.
	Content string `json:"content" api:"required"`
	// Path of the file to write.
	FilePath string `json:"file_path" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content     respjson.Field
		FilePath    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentToolset20260401WriteInput) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentToolset20260401WriteInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Tool calls are automatically approved without user confirmation.
type BetaManagedAgentsAlwaysAllowPolicy struct {
	// Any of "always_allow".
	Type BetaManagedAgentsAlwaysAllowPolicyType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAlwaysAllowPolicy) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAlwaysAllowPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsAlwaysAllowPolicy to a
// BetaManagedAgentsAlwaysAllowPolicyParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsAlwaysAllowPolicyParam.Overrides()
func (r BetaManagedAgentsAlwaysAllowPolicy) ToParam() BetaManagedAgentsAlwaysAllowPolicyParam {
	return param.Override[BetaManagedAgentsAlwaysAllowPolicyParam](json.RawMessage(r.RawJSON()))
}

type BetaManagedAgentsAlwaysAllowPolicyType string

const (
	BetaManagedAgentsAlwaysAllowPolicyTypeAlwaysAllow BetaManagedAgentsAlwaysAllowPolicyType = "always_allow"
)

// Tool calls are automatically approved without user confirmation.
//
// The property Type is required.
type BetaManagedAgentsAlwaysAllowPolicyParam struct {
	// Any of "always_allow".
	Type BetaManagedAgentsAlwaysAllowPolicyType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsAlwaysAllowPolicyParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsAlwaysAllowPolicyParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsAlwaysAllowPolicyParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Tool calls require user confirmation before execution.
type BetaManagedAgentsAlwaysAskPolicy struct {
	// Any of "always_ask".
	Type BetaManagedAgentsAlwaysAskPolicyType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAlwaysAskPolicy) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAlwaysAskPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsAlwaysAskPolicy to a
// BetaManagedAgentsAlwaysAskPolicyParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsAlwaysAskPolicyParam.Overrides()
func (r BetaManagedAgentsAlwaysAskPolicy) ToParam() BetaManagedAgentsAlwaysAskPolicyParam {
	return param.Override[BetaManagedAgentsAlwaysAskPolicyParam](json.RawMessage(r.RawJSON()))
}

type BetaManagedAgentsAlwaysAskPolicyType string

const (
	BetaManagedAgentsAlwaysAskPolicyTypeAlwaysAsk BetaManagedAgentsAlwaysAskPolicyType = "always_ask"
)

// Tool calls require user confirmation before execution.
//
// The property Type is required.
type BetaManagedAgentsAlwaysAskPolicyParam struct {
	// Any of "always_ask".
	Type BetaManagedAgentsAlwaysAskPolicyType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsAlwaysAskPolicyParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsAlwaysAskPolicyParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsAlwaysAskPolicyParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A resolved Anthropic-managed skill.
type BetaManagedAgentsAnthropicSkill struct {
	SkillID string `json:"skill_id" api:"required"`
	// Any of "anthropic".
	Type    BetaManagedAgentsAnthropicSkillType `json:"type" api:"required"`
	Version string                              `json:"version" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SkillID     respjson.Field
		Type        respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAnthropicSkill) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAnthropicSkill) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAnthropicSkillType string

const (
	BetaManagedAgentsAnthropicSkillTypeAnthropic BetaManagedAgentsAnthropicSkillType = "anthropic"
)

// An Anthropic-managed skill.
//
// The properties SkillID, Type are required.
type BetaManagedAgentsAnthropicSkillParams struct {
	// Identifier of the Anthropic skill (e.g., "xlsx").
	SkillID string `json:"skill_id" api:"required"`
	// Any of "anthropic".
	Type BetaManagedAgentsAnthropicSkillParamsType `json:"type,omitzero" api:"required"`
	// Version to pin. Defaults to latest if omitted.
	Version param.Opt[string] `json:"version,omitzero"`
	paramObj
}

func (r BetaManagedAgentsAnthropicSkillParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsAnthropicSkillParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsAnthropicSkillParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAnthropicSkillParamsType string

const (
	BetaManagedAgentsAnthropicSkillParamsTypeAnthropic BetaManagedAgentsAnthropicSkillParamsType = "anthropic"
)

// Configuration for the bash tool.
type BetaManagedAgentsBashToolConfig struct {
	Enabled bool          `json:"enabled" api:"required"`
	Name    constant.Bash `json:"name" default:"bash"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsBashToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	Type             constant.Bash                                        `json:"type" default:"bash"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsBashToolConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsBashToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsBashToolConfigPermissionPolicyUnion contains all possible
// properties and values from [BetaManagedAgentsAlwaysAllowPolicy],
// [BetaManagedAgentsAlwaysAskPolicy].
//
// Use the [BetaManagedAgentsBashToolConfigPermissionPolicyUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsBashToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsBashToolConfigPermissionPolicy is implemented by each
// variant of [BetaManagedAgentsBashToolConfigPermissionPolicyUnion] to add type
// safety for the return type of
// [BetaManagedAgentsBashToolConfigPermissionPolicyUnion.AsAny]
type anyBetaManagedAgentsBashToolConfigPermissionPolicy interface {
	implBetaManagedAgentsBashToolConfigPermissionPolicyUnion()
}

func (BetaManagedAgentsAlwaysAllowPolicy) implBetaManagedAgentsBashToolConfigPermissionPolicyUnion() {
}
func (BetaManagedAgentsAlwaysAskPolicy) implBetaManagedAgentsBashToolConfigPermissionPolicyUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsBashToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAlwaysAllowPolicy:
//	case anthropic.BetaManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsBashToolConfigPermissionPolicyUnion) AsAny() anyBetaManagedAgentsBashToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u BetaManagedAgentsBashToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v BetaManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsBashToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v BetaManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsBashToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsBashToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for the bash tool.
//
// The property Name is required.
type BetaManagedAgentsBashToolConfigParams struct {
	// Whether this tool is enabled and available to Claude. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsBashToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	// Any of "bash".
	Type BetaManagedAgentsBashToolConfigParamsType `json:"type,omitzero"`
	// Must be "bash".
	//
	// This field can be elided, and will marshal its zero value as "bash".
	Name constant.Bash `json:"name" default:"bash"`
	paramObj
}

func (r BetaManagedAgentsBashToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsBashToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsBashToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsBashToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *BetaManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *BetaManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsBashToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *BetaManagedAgentsBashToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsBashToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsBashToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsBashToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[BetaManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type BetaManagedAgentsBashToolConfigParamsType string

const (
	BetaManagedAgentsBashToolConfigParamsTypeBash BetaManagedAgentsBashToolConfigParamsType = "bash"
)

// A resolved user-created custom skill.
type BetaManagedAgentsCustomSkill struct {
	SkillID string `json:"skill_id" api:"required"`
	// Any of "custom".
	Type    BetaManagedAgentsCustomSkillType `json:"type" api:"required"`
	Version string                           `json:"version" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SkillID     respjson.Field
		Type        respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsCustomSkill) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsCustomSkill) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsCustomSkillType string

const (
	BetaManagedAgentsCustomSkillTypeCustom BetaManagedAgentsCustomSkillType = "custom"
)

// A user-created custom skill.
//
// The properties SkillID, Type are required.
type BetaManagedAgentsCustomSkillParams struct {
	// Tagged ID of the custom skill (e.g., "skill_01XJ5...").
	SkillID string `json:"skill_id" api:"required"`
	// Any of "custom".
	Type BetaManagedAgentsCustomSkillParamsType `json:"type,omitzero" api:"required"`
	// Version to pin. Defaults to latest if omitted.
	Version param.Opt[string] `json:"version,omitzero"`
	paramObj
}

func (r BetaManagedAgentsCustomSkillParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsCustomSkillParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsCustomSkillParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsCustomSkillParamsType string

const (
	BetaManagedAgentsCustomSkillParamsTypeCustom BetaManagedAgentsCustomSkillParamsType = "custom"
)

// A custom tool as returned in API responses.
type BetaManagedAgentsCustomTool struct {
	Description string `json:"description" api:"required"`
	// JSON Schema for custom tool input parameters.
	InputSchema BetaManagedAgentsCustomToolInputSchema `json:"input_schema" api:"required"`
	Name        string                                 `json:"name" api:"required"`
	// Any of "custom".
	Type BetaManagedAgentsCustomToolType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		InputSchema respjson.Field
		Name        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsCustomTool) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsCustomTool) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsCustomToolType string

const (
	BetaManagedAgentsCustomToolTypeCustom BetaManagedAgentsCustomToolType = "custom"
)

// JSON Schema for custom tool input parameters.
type BetaManagedAgentsCustomToolInputSchema struct {
	Type        constant.Object `json:"type" default:"object"`
	Properties  map[string]any  `json:"properties" api:"nullable"`
	Required    []string        `json:"required" api:"nullable"`
	ExtraFields map[string]any  `json:"" api:"extrafields"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		Properties  respjson.Field
		Required    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsCustomToolInputSchema) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsCustomToolInputSchema) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsCustomToolInputSchema to a
// BetaManagedAgentsCustomToolInputSchemaParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsCustomToolInputSchemaParam.Overrides()
func (r BetaManagedAgentsCustomToolInputSchema) ToParam() BetaManagedAgentsCustomToolInputSchemaParam {
	return param.Override[BetaManagedAgentsCustomToolInputSchemaParam](json.RawMessage(r.RawJSON()))
}

// JSON Schema for custom tool input parameters.
//
// The property Type is required.
type BetaManagedAgentsCustomToolInputSchemaParam struct {
	Properties map[string]any `json:"properties,omitzero"`
	Required   []string       `json:"required,omitzero"`
	// This field can be elided, and will marshal its zero value as "object".
	Type        constant.Object `json:"type" default:"object"`
	ExtraFields map[string]any  `json:"-"`
	paramObj
}

func (r BetaManagedAgentsCustomToolInputSchemaParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsCustomToolInputSchemaParam
	return param.MarshalWithExtras(r, (*shadow)(&r), r.ExtraFields)
}
func (r *BetaManagedAgentsCustomToolInputSchemaParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A custom tool that is executed by the API client rather than the agent. When the
// agent calls this tool, an `agent.custom_tool_use` event is emitted and the
// session goes idle, waiting for the client to provide the result via a
// `user.custom_tool_result` event.
//
// The properties Description, InputSchema, Name, Type are required.
type BetaManagedAgentsCustomToolParams struct {
	// Description of what the tool does, shown to the agent to help it decide when to
	// use the tool.
	Description string `json:"description" api:"required"`
	// JSON Schema for custom tool input parameters.
	InputSchema BetaManagedAgentsCustomToolInputSchemaParam `json:"input_schema,omitzero" api:"required"`
	// Unique name for the tool. 1-128 characters; letters, digits, underscores, and
	// hyphens.
	Name string `json:"name" api:"required"`
	// Any of "custom".
	Type BetaManagedAgentsCustomToolParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsCustomToolParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsCustomToolParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsCustomToolParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsCustomToolParamsType string

const (
	BetaManagedAgentsCustomToolParamsTypeCustom BetaManagedAgentsCustomToolParamsType = "custom"
)

// Configuration for the edit tool.
type BetaManagedAgentsEditToolConfig struct {
	Enabled bool          `json:"enabled" api:"required"`
	Name    constant.Edit `json:"name" default:"edit"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsEditToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	Type             constant.Edit                                        `json:"type" default:"edit"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsEditToolConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsEditToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsEditToolConfigPermissionPolicyUnion contains all possible
// properties and values from [BetaManagedAgentsAlwaysAllowPolicy],
// [BetaManagedAgentsAlwaysAskPolicy].
//
// Use the [BetaManagedAgentsEditToolConfigPermissionPolicyUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsEditToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsEditToolConfigPermissionPolicy is implemented by each
// variant of [BetaManagedAgentsEditToolConfigPermissionPolicyUnion] to add type
// safety for the return type of
// [BetaManagedAgentsEditToolConfigPermissionPolicyUnion.AsAny]
type anyBetaManagedAgentsEditToolConfigPermissionPolicy interface {
	implBetaManagedAgentsEditToolConfigPermissionPolicyUnion()
}

func (BetaManagedAgentsAlwaysAllowPolicy) implBetaManagedAgentsEditToolConfigPermissionPolicyUnion() {
}
func (BetaManagedAgentsAlwaysAskPolicy) implBetaManagedAgentsEditToolConfigPermissionPolicyUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsEditToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAlwaysAllowPolicy:
//	case anthropic.BetaManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsEditToolConfigPermissionPolicyUnion) AsAny() anyBetaManagedAgentsEditToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u BetaManagedAgentsEditToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v BetaManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsEditToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v BetaManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsEditToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsEditToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for the edit tool.
//
// The property Name is required.
type BetaManagedAgentsEditToolConfigParams struct {
	// Whether this tool is enabled and available to Claude. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsEditToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	// Any of "edit".
	Type BetaManagedAgentsEditToolConfigParamsType `json:"type,omitzero"`
	// Must be "edit".
	//
	// This field can be elided, and will marshal its zero value as "edit".
	Name constant.Edit `json:"name" default:"edit"`
	paramObj
}

func (r BetaManagedAgentsEditToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsEditToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsEditToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsEditToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *BetaManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *BetaManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsEditToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *BetaManagedAgentsEditToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsEditToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsEditToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsEditToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[BetaManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type BetaManagedAgentsEditToolConfigParamsType string

const (
	BetaManagedAgentsEditToolConfigParamsTypeEdit BetaManagedAgentsEditToolConfigParamsType = "edit"
)

// High effort. Favors reasoning depth.
type BetaManagedAgentsEffortHigh struct {
	// Any of "high".
	Type BetaManagedAgentsEffortHighType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsEffortHigh) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsEffortHigh) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsEffortHigh to a
// BetaManagedAgentsEffortHighParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsEffortHighParam.Overrides()
func (r BetaManagedAgentsEffortHigh) ToParam() BetaManagedAgentsEffortHighParam {
	return param.Override[BetaManagedAgentsEffortHighParam](json.RawMessage(r.RawJSON()))
}

type BetaManagedAgentsEffortHighType string

const (
	BetaManagedAgentsEffortHighTypeHigh BetaManagedAgentsEffortHighType = "high"
)

// High effort. Favors reasoning depth.
//
// The property Type is required.
type BetaManagedAgentsEffortHighParam struct {
	// Any of "high".
	Type BetaManagedAgentsEffortHighType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsEffortHighParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsEffortHighParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsEffortHighParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Low effort. Favors latency over reasoning depth.
type BetaManagedAgentsEffortLow struct {
	// Any of "low".
	Type BetaManagedAgentsEffortLowType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsEffortLow) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsEffortLow) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsEffortLow to a
// BetaManagedAgentsEffortLowParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsEffortLowParam.Overrides()
func (r BetaManagedAgentsEffortLow) ToParam() BetaManagedAgentsEffortLowParam {
	return param.Override[BetaManagedAgentsEffortLowParam](json.RawMessage(r.RawJSON()))
}

type BetaManagedAgentsEffortLowType string

const (
	BetaManagedAgentsEffortLowTypeLow BetaManagedAgentsEffortLowType = "low"
)

// Low effort. Favors latency over reasoning depth.
//
// The property Type is required.
type BetaManagedAgentsEffortLowParam struct {
	// Any of "low".
	Type BetaManagedAgentsEffortLowType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsEffortLowParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsEffortLowParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsEffortLowParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Maximum effort. Favors reasoning depth over latency.
type BetaManagedAgentsEffortMax struct {
	// Any of "max".
	Type BetaManagedAgentsEffortMaxType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsEffortMax) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsEffortMax) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsEffortMax to a
// BetaManagedAgentsEffortMaxParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsEffortMaxParam.Overrides()
func (r BetaManagedAgentsEffortMax) ToParam() BetaManagedAgentsEffortMaxParam {
	return param.Override[BetaManagedAgentsEffortMaxParam](json.RawMessage(r.RawJSON()))
}

type BetaManagedAgentsEffortMaxType string

const (
	BetaManagedAgentsEffortMaxTypeMax BetaManagedAgentsEffortMaxType = "max"
)

// Maximum effort. Favors reasoning depth over latency.
//
// The property Type is required.
type BetaManagedAgentsEffortMaxParam struct {
	// Any of "max".
	Type BetaManagedAgentsEffortMaxType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsEffortMaxParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsEffortMaxParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsEffortMaxParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Medium effort. Balances latency and reasoning depth.
type BetaManagedAgentsEffortMedium struct {
	// Any of "medium".
	Type BetaManagedAgentsEffortMediumType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsEffortMedium) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsEffortMedium) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsEffortMedium to a
// BetaManagedAgentsEffortMediumParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsEffortMediumParam.Overrides()
func (r BetaManagedAgentsEffortMedium) ToParam() BetaManagedAgentsEffortMediumParam {
	return param.Override[BetaManagedAgentsEffortMediumParam](json.RawMessage(r.RawJSON()))
}

type BetaManagedAgentsEffortMediumType string

const (
	BetaManagedAgentsEffortMediumTypeMedium BetaManagedAgentsEffortMediumType = "medium"
)

// Medium effort. Balances latency and reasoning depth.
//
// The property Type is required.
type BetaManagedAgentsEffortMediumParam struct {
	// Any of "medium".
	Type BetaManagedAgentsEffortMediumType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsEffortMediumParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsEffortMediumParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsEffortMediumParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Extra-high effort. Not all models accept this level.
type BetaManagedAgentsEffortXhigh struct {
	// Any of "xhigh".
	Type BetaManagedAgentsEffortXhighType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsEffortXhigh) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsEffortXhigh) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsEffortXhigh to a
// BetaManagedAgentsEffortXhighParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsEffortXhighParam.Overrides()
func (r BetaManagedAgentsEffortXhigh) ToParam() BetaManagedAgentsEffortXhighParam {
	return param.Override[BetaManagedAgentsEffortXhighParam](json.RawMessage(r.RawJSON()))
}

type BetaManagedAgentsEffortXhighType string

const (
	BetaManagedAgentsEffortXhighTypeXhigh BetaManagedAgentsEffortXhighType = "xhigh"
)

// Extra-high effort. Not all models accept this level.
//
// The property Type is required.
type BetaManagedAgentsEffortXhighParam struct {
	// Any of "xhigh".
	Type BetaManagedAgentsEffortXhighType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsEffortXhighParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsEffortXhighParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsEffortXhighParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for the glob tool.
type BetaManagedAgentsGlobToolConfig struct {
	Enabled bool          `json:"enabled" api:"required"`
	Name    constant.Glob `json:"name" default:"glob"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsGlobToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	Type             constant.Glob                                        `json:"type" default:"glob"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsGlobToolConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsGlobToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsGlobToolConfigPermissionPolicyUnion contains all possible
// properties and values from [BetaManagedAgentsAlwaysAllowPolicy],
// [BetaManagedAgentsAlwaysAskPolicy].
//
// Use the [BetaManagedAgentsGlobToolConfigPermissionPolicyUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsGlobToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsGlobToolConfigPermissionPolicy is implemented by each
// variant of [BetaManagedAgentsGlobToolConfigPermissionPolicyUnion] to add type
// safety for the return type of
// [BetaManagedAgentsGlobToolConfigPermissionPolicyUnion.AsAny]
type anyBetaManagedAgentsGlobToolConfigPermissionPolicy interface {
	implBetaManagedAgentsGlobToolConfigPermissionPolicyUnion()
}

func (BetaManagedAgentsAlwaysAllowPolicy) implBetaManagedAgentsGlobToolConfigPermissionPolicyUnion() {
}
func (BetaManagedAgentsAlwaysAskPolicy) implBetaManagedAgentsGlobToolConfigPermissionPolicyUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsGlobToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAlwaysAllowPolicy:
//	case anthropic.BetaManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsGlobToolConfigPermissionPolicyUnion) AsAny() anyBetaManagedAgentsGlobToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u BetaManagedAgentsGlobToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v BetaManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsGlobToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v BetaManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsGlobToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsGlobToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for the glob tool.
//
// The property Name is required.
type BetaManagedAgentsGlobToolConfigParams struct {
	// Whether this tool is enabled and available to Claude. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsGlobToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	// Any of "glob".
	Type BetaManagedAgentsGlobToolConfigParamsType `json:"type,omitzero"`
	// Must be "glob".
	//
	// This field can be elided, and will marshal its zero value as "glob".
	Name constant.Glob `json:"name" default:"glob"`
	paramObj
}

func (r BetaManagedAgentsGlobToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsGlobToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsGlobToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsGlobToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *BetaManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *BetaManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsGlobToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *BetaManagedAgentsGlobToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsGlobToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsGlobToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsGlobToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[BetaManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type BetaManagedAgentsGlobToolConfigParamsType string

const (
	BetaManagedAgentsGlobToolConfigParamsTypeGlob BetaManagedAgentsGlobToolConfigParamsType = "glob"
)

// Configuration for the grep tool.
type BetaManagedAgentsGrepToolConfig struct {
	Enabled bool          `json:"enabled" api:"required"`
	Name    constant.Grep `json:"name" default:"grep"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsGrepToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	Type             constant.Grep                                        `json:"type" default:"grep"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsGrepToolConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsGrepToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsGrepToolConfigPermissionPolicyUnion contains all possible
// properties and values from [BetaManagedAgentsAlwaysAllowPolicy],
// [BetaManagedAgentsAlwaysAskPolicy].
//
// Use the [BetaManagedAgentsGrepToolConfigPermissionPolicyUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsGrepToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsGrepToolConfigPermissionPolicy is implemented by each
// variant of [BetaManagedAgentsGrepToolConfigPermissionPolicyUnion] to add type
// safety for the return type of
// [BetaManagedAgentsGrepToolConfigPermissionPolicyUnion.AsAny]
type anyBetaManagedAgentsGrepToolConfigPermissionPolicy interface {
	implBetaManagedAgentsGrepToolConfigPermissionPolicyUnion()
}

func (BetaManagedAgentsAlwaysAllowPolicy) implBetaManagedAgentsGrepToolConfigPermissionPolicyUnion() {
}
func (BetaManagedAgentsAlwaysAskPolicy) implBetaManagedAgentsGrepToolConfigPermissionPolicyUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsGrepToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAlwaysAllowPolicy:
//	case anthropic.BetaManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsGrepToolConfigPermissionPolicyUnion) AsAny() anyBetaManagedAgentsGrepToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u BetaManagedAgentsGrepToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v BetaManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsGrepToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v BetaManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsGrepToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsGrepToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for the grep tool.
//
// The property Name is required.
type BetaManagedAgentsGrepToolConfigParams struct {
	// Whether this tool is enabled and available to Claude. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsGrepToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	// Any of "grep".
	Type BetaManagedAgentsGrepToolConfigParamsType `json:"type,omitzero"`
	// Must be "grep".
	//
	// This field can be elided, and will marshal its zero value as "grep".
	Name constant.Grep `json:"name" default:"grep"`
	paramObj
}

func (r BetaManagedAgentsGrepToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsGrepToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsGrepToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsGrepToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *BetaManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *BetaManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsGrepToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *BetaManagedAgentsGrepToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsGrepToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsGrepToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsGrepToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[BetaManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type BetaManagedAgentsGrepToolConfigParamsType string

const (
	BetaManagedAgentsGrepToolConfigParamsTypeGrep BetaManagedAgentsGrepToolConfigParamsType = "grep"
)

// URL-based MCP server connection as returned in API responses.
type BetaManagedAgentsMCPServerURLDefinition struct {
	Name string `json:"name" api:"required"`
	// Any of "url".
	Type BetaManagedAgentsMCPServerURLDefinitionType `json:"type" api:"required"`
	URL  string                                      `json:"url" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Type        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMCPServerURLDefinition) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsMCPServerURLDefinition) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMCPServerURLDefinitionType string

const (
	BetaManagedAgentsMCPServerURLDefinitionTypeURL BetaManagedAgentsMCPServerURLDefinitionType = "url"
)

// Resolved configuration for a specific MCP tool.
type BetaManagedAgentsMCPToolConfig struct {
	Enabled bool   `json:"enabled" api:"required"`
	Name    string `json:"name" api:"required"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsMCPToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMCPToolConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsMCPToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsMCPToolConfigPermissionPolicyUnion contains all possible
// properties and values from [BetaManagedAgentsAlwaysAllowPolicy],
// [BetaManagedAgentsAlwaysAskPolicy].
//
// Use the [BetaManagedAgentsMCPToolConfigPermissionPolicyUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsMCPToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsMCPToolConfigPermissionPolicy is implemented by each variant
// of [BetaManagedAgentsMCPToolConfigPermissionPolicyUnion] to add type safety for
// the return type of [BetaManagedAgentsMCPToolConfigPermissionPolicyUnion.AsAny]
type anyBetaManagedAgentsMCPToolConfigPermissionPolicy interface {
	implBetaManagedAgentsMCPToolConfigPermissionPolicyUnion()
}

func (BetaManagedAgentsAlwaysAllowPolicy) implBetaManagedAgentsMCPToolConfigPermissionPolicyUnion() {}
func (BetaManagedAgentsAlwaysAskPolicy) implBetaManagedAgentsMCPToolConfigPermissionPolicyUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsMCPToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAlwaysAllowPolicy:
//	case anthropic.BetaManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsMCPToolConfigPermissionPolicyUnion) AsAny() anyBetaManagedAgentsMCPToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u BetaManagedAgentsMCPToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v BetaManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsMCPToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v BetaManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsMCPToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsMCPToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for a specific MCP tool.
//
// The property Name is required.
type BetaManagedAgentsMCPToolConfigParams struct {
	// Name of the MCP tool to configure. 1-128 characters.
	Name string `json:"name" api:"required"`
	// Whether this tool is enabled. Overrides the `default_config` setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsMCPToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	paramObj
}

func (r BetaManagedAgentsMCPToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsMCPToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsMCPToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsMCPToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *BetaManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *BetaManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsMCPToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *BetaManagedAgentsMCPToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsMCPToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsMCPToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsMCPToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[BetaManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type BetaManagedAgentsMCPToolset struct {
	Configs []BetaManagedAgentsMCPToolConfig `json:"configs" api:"required"`
	// Resolved default configuration for all tools from an MCP server.
	DefaultConfig BetaManagedAgentsMCPToolsetDefaultConfig `json:"default_config" api:"required"`
	MCPServerName string                                   `json:"mcp_server_name" api:"required"`
	// Any of "mcp_toolset".
	Type BetaManagedAgentsMCPToolsetType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Configs       respjson.Field
		DefaultConfig respjson.Field
		MCPServerName respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMCPToolset) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsMCPToolset) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMCPToolsetType string

const (
	BetaManagedAgentsMCPToolsetTypeMCPToolset BetaManagedAgentsMCPToolsetType = "mcp_toolset"
)

// Resolved default configuration for all tools from an MCP server.
type BetaManagedAgentsMCPToolsetDefaultConfig struct {
	Enabled bool `json:"enabled" api:"required"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		PermissionPolicy respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMCPToolsetDefaultConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsMCPToolsetDefaultConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion contains all
// possible properties and values from [BetaManagedAgentsAlwaysAllowPolicy],
// [BetaManagedAgentsAlwaysAskPolicy].
//
// Use the [BetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicy is implemented by
// each variant of [BetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion]
// to add type safety for the return type of
// [BetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion.AsAny]
type anyBetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicy interface {
	implBetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion()
}

func (BetaManagedAgentsAlwaysAllowPolicy) implBetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion() {
}
func (BetaManagedAgentsAlwaysAskPolicy) implBetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAlwaysAllowPolicy:
//	case anthropic.BetaManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion) AsAny() anyBetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u BetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion) AsAlwaysAllow() (v BetaManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion) AsAlwaysAsk() (v BetaManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *BetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Default configuration for all tools from an MCP server.
type BetaManagedAgentsMCPToolsetDefaultConfigParams struct {
	// Whether tools are enabled by default. Defaults to true if not specified.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	paramObj
}

func (r BetaManagedAgentsMCPToolsetDefaultConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsMCPToolsetDefaultConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsMCPToolsetDefaultConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *BetaManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *BetaManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *BetaManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[BetaManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

// Configuration for tools from an MCP server defined in `mcp_servers`.
//
// The properties MCPServerName, Type are required.
type BetaManagedAgentsMCPToolsetParams struct {
	// Name of the MCP server. Must match a server name from the mcp_servers array.
	// 1-255 characters.
	MCPServerName string `json:"mcp_server_name" api:"required"`
	// Any of "mcp_toolset".
	Type BetaManagedAgentsMCPToolsetParamsType `json:"type,omitzero" api:"required"`
	// Per-tool configuration overrides.
	Configs []BetaManagedAgentsMCPToolConfigParams `json:"configs,omitzero"`
	// Default configuration for all tools from an MCP server.
	DefaultConfig BetaManagedAgentsMCPToolsetDefaultConfigParams `json:"default_config,omitzero"`
	paramObj
}

func (r BetaManagedAgentsMCPToolsetParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsMCPToolsetParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsMCPToolsetParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMCPToolsetParamsType string

const (
	BetaManagedAgentsMCPToolsetParamsTypeMCPToolset BetaManagedAgentsMCPToolsetParamsType = "mcp_toolset"
)

// The model that will power your agent.
//
// See [models](https://docs.anthropic.com/en/docs/models-overview) for additional
// details and options.
type BetaManagedAgentsModel = string

const (
	BetaManagedAgentsModelClaudeFable5_1           BetaManagedAgentsModel = "claude-fable-5-1"
	BetaManagedAgentsModelClaudeSonnet5            BetaManagedAgentsModel = "claude-sonnet-5"
	BetaManagedAgentsModelClaudeFable5             BetaManagedAgentsModel = "claude-fable-5"
	BetaManagedAgentsModelClaudeOpus5              BetaManagedAgentsModel = "claude-opus-5"
	BetaManagedAgentsModelClaudeOpus4_8            BetaManagedAgentsModel = "claude-opus-4-8"
	BetaManagedAgentsModelClaudeOpus4_7            BetaManagedAgentsModel = "claude-opus-4-7"
	BetaManagedAgentsModelClaudeOpus4_6            BetaManagedAgentsModel = "claude-opus-4-6"
	BetaManagedAgentsModelClaudeSonnet4_6          BetaManagedAgentsModel = "claude-sonnet-4-6"
	BetaManagedAgentsModelClaudeHaiku4_5           BetaManagedAgentsModel = "claude-haiku-4-5"
	BetaManagedAgentsModelClaudeHaiku4_5_20251001  BetaManagedAgentsModel = "claude-haiku-4-5-20251001"
	BetaManagedAgentsModelClaudeOpus4_5            BetaManagedAgentsModel = "claude-opus-4-5"
	BetaManagedAgentsModelClaudeOpus4_5_20251101   BetaManagedAgentsModel = "claude-opus-4-5-20251101"
	BetaManagedAgentsModelClaudeSonnet4_5          BetaManagedAgentsModel = "claude-sonnet-4-5"
	BetaManagedAgentsModelClaudeSonnet4_5_20250929 BetaManagedAgentsModel = "claude-sonnet-4-5-20250929"
)

// Model identifier and configuration.
type BetaManagedAgentsModelConfig struct {
	// The model that will power your agent.
	//
	// See [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	ID BetaManagedAgentsModel `json:"id" api:"required"`
	// How hard Claude works on each turn. Sets `output_config.effort` on every
	// Messages call the session makes.
	Effort BetaManagedAgentsModelConfigEffortUnion `json:"effort"`
	// Geographic region for model inference. When unset, requests fall through to the
	// workspace's default_inference_geo.
	InferenceGeo string `json:"inference_geo"`
	// Inference speed mode. `fast` provides significantly faster output token
	// generation at premium pricing. Not all models support `fast`; invalid
	// combinations are rejected at create time.
	//
	// Any of "standard", "fast".
	Speed BetaManagedAgentsModelConfigSpeed `json:"speed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		Effort       respjson.Field
		InferenceGeo respjson.Field
		Speed        respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsModelConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsModelConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsModelConfigEffortUnion contains all possible properties and
// values from [BetaManagedAgentsEffortLow], [BetaManagedAgentsEffortMedium],
// [BetaManagedAgentsEffortHigh], [BetaManagedAgentsEffortXhigh],
// [BetaManagedAgentsEffortMax].
//
// Use the [BetaManagedAgentsModelConfigEffortUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsModelConfigEffortUnion struct {
	// Any of "low", "medium", "high", "xhigh", "max".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsModelConfigEffort is implemented by each variant of
// [BetaManagedAgentsModelConfigEffortUnion] to add type safety for the return type
// of [BetaManagedAgentsModelConfigEffortUnion.AsAny]
type anyBetaManagedAgentsModelConfigEffort interface {
	implBetaManagedAgentsModelConfigEffortUnion()
}

func (BetaManagedAgentsEffortLow) implBetaManagedAgentsModelConfigEffortUnion()    {}
func (BetaManagedAgentsEffortMedium) implBetaManagedAgentsModelConfigEffortUnion() {}
func (BetaManagedAgentsEffortHigh) implBetaManagedAgentsModelConfigEffortUnion()   {}
func (BetaManagedAgentsEffortXhigh) implBetaManagedAgentsModelConfigEffortUnion()  {}
func (BetaManagedAgentsEffortMax) implBetaManagedAgentsModelConfigEffortUnion()    {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsModelConfigEffortUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsEffortLow:
//	case anthropic.BetaManagedAgentsEffortMedium:
//	case anthropic.BetaManagedAgentsEffortHigh:
//	case anthropic.BetaManagedAgentsEffortXhigh:
//	case anthropic.BetaManagedAgentsEffortMax:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsModelConfigEffortUnion) AsAny() anyBetaManagedAgentsModelConfigEffort {
	switch u.Type {
	case "low":
		return u.AsLow()
	case "medium":
		return u.AsMedium()
	case "high":
		return u.AsHigh()
	case "xhigh":
		return u.AsXhigh()
	case "max":
		return u.AsMax()
	}
	return nil
}

func (u BetaManagedAgentsModelConfigEffortUnion) AsLow() (v BetaManagedAgentsEffortLow) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsModelConfigEffortUnion) AsMedium() (v BetaManagedAgentsEffortMedium) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsModelConfigEffortUnion) AsHigh() (v BetaManagedAgentsEffortHigh) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsModelConfigEffortUnion) AsXhigh() (v BetaManagedAgentsEffortXhigh) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsModelConfigEffortUnion) AsMax() (v BetaManagedAgentsEffortMax) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsModelConfigEffortUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsModelConfigEffortUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inference speed mode. `fast` provides significantly faster output token
// generation at premium pricing. Not all models support `fast`; invalid
// combinations are rejected at create time.
type BetaManagedAgentsModelConfigSpeed string

const (
	BetaManagedAgentsModelConfigSpeedStandard BetaManagedAgentsModelConfigSpeed = "standard"
	BetaManagedAgentsModelConfigSpeedFast     BetaManagedAgentsModelConfigSpeed = "fast"
)

// An object that defines additional configuration control over model use
//
// The property ID is required.
type BetaManagedAgentsModelConfigParams struct {
	// The model that will power your agent.
	//
	// See [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	ID BetaManagedAgentsModel `json:"id,omitzero" api:"required"`
	// Geographic region for model inference. When unset, requests fall through to the
	// workspace's default_inference_geo. On update, `model` is whole-object
	// replacement — omitting inference_geo clears it.
	InferenceGeo param.Opt[string] `json:"inference_geo,omitzero"`
	// How hard Claude works on each inference call. Accepts a bare level string
	// (`"high"`) or `{"type": "high"}`. On create, omitting it resolves the per-model
	// default; on update, omitting it leaves the stored value unchanged.
	Effort BetaManagedAgentsModelConfigParamsEffortUnion `json:"effort,omitzero"`
	// Inference speed mode. `fast` provides significantly faster output token
	// generation at premium pricing. Not all models support `fast`; invalid
	// combinations are rejected at create time.
	//
	// Any of "standard", "fast".
	Speed BetaManagedAgentsModelConfigParamsSpeed `json:"speed,omitzero"`
	paramObj
}

func (r BetaManagedAgentsModelConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsModelConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsModelConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsModelConfigParamsEffortUnion struct {
	// Check if union is this variant with
	// !param.IsOmitted(union.OfBetaManagedAgentsModelConfigsEffortBetaManagedAgentsEffortLevel)
	OfBetaManagedAgentsModelConfigsEffortBetaManagedAgentsEffortLevel param.Opt[string]                   `json:",omitzero,inline"`
	OfBetaManagedAgentsEffortLow                                      *BetaManagedAgentsEffortLowParam    `json:",omitzero,inline"`
	OfBetaManagedAgentsEffortMedium                                   *BetaManagedAgentsEffortMediumParam `json:",omitzero,inline"`
	OfBetaManagedAgentsEffortHigh                                     *BetaManagedAgentsEffortHighParam   `json:",omitzero,inline"`
	OfBetaManagedAgentsEffortXhigh                                    *BetaManagedAgentsEffortXhighParam  `json:",omitzero,inline"`
	OfBetaManagedAgentsEffortMax                                      *BetaManagedAgentsEffortMaxParam    `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsModelConfigParamsEffortUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfBetaManagedAgentsModelConfigsEffortBetaManagedAgentsEffortLevel,
		u.OfBetaManagedAgentsEffortLow,
		u.OfBetaManagedAgentsEffortMedium,
		u.OfBetaManagedAgentsEffortHigh,
		u.OfBetaManagedAgentsEffortXhigh,
		u.OfBetaManagedAgentsEffortMax)
}
func (u *BetaManagedAgentsModelConfigParamsEffortUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsModelConfigParamsEffortUnion) asAny() any {
	if !param.IsOmitted(u.OfBetaManagedAgentsModelConfigsEffortBetaManagedAgentsEffortLevel) {
		return &u.OfBetaManagedAgentsModelConfigsEffortBetaManagedAgentsEffortLevel
	} else if !param.IsOmitted(u.OfBetaManagedAgentsEffortLow) {
		return u.OfBetaManagedAgentsEffortLow
	} else if !param.IsOmitted(u.OfBetaManagedAgentsEffortMedium) {
		return u.OfBetaManagedAgentsEffortMedium
	} else if !param.IsOmitted(u.OfBetaManagedAgentsEffortHigh) {
		return u.OfBetaManagedAgentsEffortHigh
	} else if !param.IsOmitted(u.OfBetaManagedAgentsEffortXhigh) {
		return u.OfBetaManagedAgentsEffortXhigh
	} else if !param.IsOmitted(u.OfBetaManagedAgentsEffortMax) {
		return u.OfBetaManagedAgentsEffortMax
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsModelConfigParamsEffortUnion) GetType() *string {
	if vt := u.OfBetaManagedAgentsEffortLow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfBetaManagedAgentsEffortMedium; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfBetaManagedAgentsEffortHigh; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfBetaManagedAgentsEffortXhigh; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfBetaManagedAgentsEffortMax; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// How hard Claude works on each turn. Higher levels favor reasoning depth over
// latency. Not all models accept every level; invalid combinations are rejected at
// create time.
type BetaManagedAgentsModelConfigParamsEffortBetaManagedAgentsEffortLevel string

const (
	BetaManagedAgentsModelConfigParamsEffortBetaManagedAgentsEffortLevelLow    BetaManagedAgentsModelConfigParamsEffortBetaManagedAgentsEffortLevel = "low"
	BetaManagedAgentsModelConfigParamsEffortBetaManagedAgentsEffortLevelMedium BetaManagedAgentsModelConfigParamsEffortBetaManagedAgentsEffortLevel = "medium"
	BetaManagedAgentsModelConfigParamsEffortBetaManagedAgentsEffortLevelHigh   BetaManagedAgentsModelConfigParamsEffortBetaManagedAgentsEffortLevel = "high"
	BetaManagedAgentsModelConfigParamsEffortBetaManagedAgentsEffortLevelXhigh  BetaManagedAgentsModelConfigParamsEffortBetaManagedAgentsEffortLevel = "xhigh"
	BetaManagedAgentsModelConfigParamsEffortBetaManagedAgentsEffortLevelMax    BetaManagedAgentsModelConfigParamsEffortBetaManagedAgentsEffortLevel = "max"
)

// Inference speed mode. `fast` provides significantly faster output token
// generation at premium pricing. Not all models support `fast`; invalid
// combinations are rejected at create time.
type BetaManagedAgentsModelConfigParamsSpeed string

const (
	BetaManagedAgentsModelConfigParamsSpeedStandard BetaManagedAgentsModelConfigParamsSpeed = "standard"
	BetaManagedAgentsModelConfigParamsSpeedFast     BetaManagedAgentsModelConfigParamsSpeed = "fast"
)

// Sentinel roster entry meaning "the agent that owns this configuration". Resolved
// server-side to a concrete agent reference.
//
// The property Type is required.
type BetaManagedAgentsMultiagentSelfParams struct {
	// Any of "self".
	Type BetaManagedAgentsMultiagentSelfParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsMultiagentSelfParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsMultiagentSelfParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsMultiagentSelfParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMultiagentSelfParamsType string

const (
	BetaManagedAgentsMultiagentSelfParamsTypeSelf BetaManagedAgentsMultiagentSelfParamsType = "self"
)

// Configuration for the read tool.
type BetaManagedAgentsReadToolConfig struct {
	Enabled bool          `json:"enabled" api:"required"`
	Name    constant.Read `json:"name" default:"read"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsReadToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	Type             constant.Read                                        `json:"type" default:"read"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsReadToolConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsReadToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsReadToolConfigPermissionPolicyUnion contains all possible
// properties and values from [BetaManagedAgentsAlwaysAllowPolicy],
// [BetaManagedAgentsAlwaysAskPolicy].
//
// Use the [BetaManagedAgentsReadToolConfigPermissionPolicyUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsReadToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsReadToolConfigPermissionPolicy is implemented by each
// variant of [BetaManagedAgentsReadToolConfigPermissionPolicyUnion] to add type
// safety for the return type of
// [BetaManagedAgentsReadToolConfigPermissionPolicyUnion.AsAny]
type anyBetaManagedAgentsReadToolConfigPermissionPolicy interface {
	implBetaManagedAgentsReadToolConfigPermissionPolicyUnion()
}

func (BetaManagedAgentsAlwaysAllowPolicy) implBetaManagedAgentsReadToolConfigPermissionPolicyUnion() {
}
func (BetaManagedAgentsAlwaysAskPolicy) implBetaManagedAgentsReadToolConfigPermissionPolicyUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsReadToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAlwaysAllowPolicy:
//	case anthropic.BetaManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsReadToolConfigPermissionPolicyUnion) AsAny() anyBetaManagedAgentsReadToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u BetaManagedAgentsReadToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v BetaManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsReadToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v BetaManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsReadToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsReadToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for the read tool.
//
// The property Name is required.
type BetaManagedAgentsReadToolConfigParams struct {
	// Whether this tool is enabled and available to Claude. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsReadToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	// Any of "read".
	Type BetaManagedAgentsReadToolConfigParamsType `json:"type,omitzero"`
	// Must be "read".
	//
	// This field can be elided, and will marshal its zero value as "read".
	Name constant.Read `json:"name" default:"read"`
	paramObj
}

func (r BetaManagedAgentsReadToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsReadToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsReadToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsReadToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *BetaManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *BetaManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsReadToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *BetaManagedAgentsReadToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsReadToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsReadToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsReadToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[BetaManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type BetaManagedAgentsReadToolConfigParamsType string

const (
	BetaManagedAgentsReadToolConfigParamsTypeRead BetaManagedAgentsReadToolConfigParamsType = "read"
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
	// This field is a union of [[]BetaManagedAgentsAgentToolConfigUnion],
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
	// [[]BetaManagedAgentsAgentToolConfigUnion] instead of an object.
	OfBetaManagedAgentsAgentToolConfigArray []BetaManagedAgentsAgentToolConfigUnion `json:",inline"`
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

func BetaManagedAgentsSkillParamsOfAnthropic(skillID string) BetaManagedAgentsSkillParamsUnion {
	var anthropic BetaManagedAgentsAnthropicSkillParams
	anthropic.SkillID = skillID
	return BetaManagedAgentsSkillParamsUnion{OfAnthropic: &anthropic}
}

func BetaManagedAgentsSkillParamsOfCustom(skillID string) BetaManagedAgentsSkillParamsUnion {
	var custom BetaManagedAgentsCustomSkillParams
	custom.SkillID = skillID
	return BetaManagedAgentsSkillParamsUnion{OfCustom: &custom}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsSkillParamsUnion struct {
	OfAnthropic *BetaManagedAgentsAnthropicSkillParams `json:",omitzero,inline"`
	OfCustom    *BetaManagedAgentsCustomSkillParams    `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsSkillParamsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAnthropic, u.OfCustom)
}
func (u *BetaManagedAgentsSkillParamsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsSkillParamsUnion) asAny() any {
	if !param.IsOmitted(u.OfAnthropic) {
		return u.OfAnthropic
	} else if !param.IsOmitted(u.OfCustom) {
		return u.OfCustom
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsSkillParamsUnion) GetSkillID() *string {
	if vt := u.OfAnthropic; vt != nil {
		return (*string)(&vt.SkillID)
	} else if vt := u.OfCustom; vt != nil {
		return (*string)(&vt.SkillID)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsSkillParamsUnion) GetType() *string {
	if vt := u.OfAnthropic; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCustom; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsSkillParamsUnion) GetVersion() *string {
	if vt := u.OfAnthropic; vt != nil && vt.Version.Valid() {
		return &vt.Version.Value
	} else if vt := u.OfCustom; vt != nil && vt.Version.Valid() {
		return &vt.Version.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsSkillParamsUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsAnthropicSkillParams]("anthropic"),
		apijson.Discriminator[BetaManagedAgentsCustomSkillParams]("custom"),
	)
}

// URL-based MCP server connection.
//
// The properties Name, Type, URL are required.
type BetaManagedAgentsURLMCPServerParams struct {
	// Unique name for this server, referenced by mcp_toolset configurations. 1-255
	// characters.
	Name string `json:"name" api:"required"`
	// Any of "url".
	Type BetaManagedAgentsURLMCPServerParamsType `json:"type,omitzero" api:"required"`
	// Endpoint URL for the MCP server.
	URL string `json:"url" api:"required"`
	paramObj
}

func (r BetaManagedAgentsURLMCPServerParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsURLMCPServerParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsURLMCPServerParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsURLMCPServerParamsType string

const (
	BetaManagedAgentsURLMCPServerParamsTypeURL BetaManagedAgentsURLMCPServerParamsType = "url"
)

// Approximate user location for search result localization.
type BetaManagedAgentsUserLocation struct {
	// Location precision. Only "approximate" is supported.
	Type constant.Approximate `json:"type" default:"approximate"`
	// City name.
	City string `json:"city" api:"nullable"`
	// Two-letter ISO 3166-1 country code, uppercase.
	Country string `json:"country" api:"nullable"`
	// Region or state name.
	Region string `json:"region" api:"nullable"`
	// IANA timezone identifier, e.g. "America/Los_Angeles".
	Timezone string `json:"timezone" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		City        respjson.Field
		Country     respjson.Field
		Region      respjson.Field
		Timezone    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsUserLocation) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsUserLocation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsUserLocation to a
// BetaManagedAgentsUserLocationParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsUserLocationParam.Overrides()
func (r BetaManagedAgentsUserLocation) ToParam() BetaManagedAgentsUserLocationParam {
	return param.Override[BetaManagedAgentsUserLocationParam](json.RawMessage(r.RawJSON()))
}

// Approximate user location for search result localization.
//
// The property Type is required.
type BetaManagedAgentsUserLocationParam struct {
	// City name.
	City param.Opt[string] `json:"city,omitzero"`
	// Two-letter ISO 3166-1 country code, uppercase.
	Country param.Opt[string] `json:"country,omitzero"`
	// Region or state name.
	Region param.Opt[string] `json:"region,omitzero"`
	// IANA timezone identifier, e.g. "America/Los_Angeles".
	Timezone param.Opt[string] `json:"timezone,omitzero"`
	// Location precision. Only "approximate" is supported.
	//
	// This field can be elided, and will marshal its zero value as "approximate".
	Type constant.Approximate `json:"type" default:"approximate"`
	paramObj
}

func (r BetaManagedAgentsUserLocationParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsUserLocationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsUserLocationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for the web_fetch tool.
type BetaManagedAgentsWebFetchToolConfig struct {
	Enabled bool              `json:"enabled" api:"required"`
	Name    constant.WebFetch `json:"name" default:"web_fetch"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsWebFetchToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	Type             constant.WebFetch                                        `json:"type" default:"web_fetch"`
	AllowedDomains   []string                                                 `json:"allowed_domains"`
	BlockedDomains   []string                                                 `json:"blocked_domains"`
	MaxContentTokens int64                                                    `json:"max_content_tokens" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		AllowedDomains   respjson.Field
		BlockedDomains   respjson.Field
		MaxContentTokens respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsWebFetchToolConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsWebFetchToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsWebFetchToolConfigPermissionPolicyUnion contains all possible
// properties and values from [BetaManagedAgentsAlwaysAllowPolicy],
// [BetaManagedAgentsAlwaysAskPolicy].
//
// Use the [BetaManagedAgentsWebFetchToolConfigPermissionPolicyUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsWebFetchToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsWebFetchToolConfigPermissionPolicy is implemented by each
// variant of [BetaManagedAgentsWebFetchToolConfigPermissionPolicyUnion] to add
// type safety for the return type of
// [BetaManagedAgentsWebFetchToolConfigPermissionPolicyUnion.AsAny]
type anyBetaManagedAgentsWebFetchToolConfigPermissionPolicy interface {
	implBetaManagedAgentsWebFetchToolConfigPermissionPolicyUnion()
}

func (BetaManagedAgentsAlwaysAllowPolicy) implBetaManagedAgentsWebFetchToolConfigPermissionPolicyUnion() {
}
func (BetaManagedAgentsAlwaysAskPolicy) implBetaManagedAgentsWebFetchToolConfigPermissionPolicyUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsWebFetchToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAlwaysAllowPolicy:
//	case anthropic.BetaManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsWebFetchToolConfigPermissionPolicyUnion) AsAny() anyBetaManagedAgentsWebFetchToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u BetaManagedAgentsWebFetchToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v BetaManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsWebFetchToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v BetaManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsWebFetchToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsWebFetchToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for the web_fetch tool.
//
// The property Name is required.
type BetaManagedAgentsWebFetchToolConfigParams struct {
	// Whether this tool is enabled and available to Claude. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Maximum number of tokens of fetched text content to include in context per call.
	// Does not apply to binary content such as PDFs.
	MaxContentTokens param.Opt[int64] `json:"max_content_tokens,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	// Only fetch URLs whose host is one of these domains or a subdomain of one. Each
	// entry is a plain hostname like "docs.example.com" (no scheme, port, or path). At
	// most 64 entries; an empty list is rejected (omit the field instead). Cannot be
	// combined with blocked_domains.
	AllowedDomains []string `json:"allowed_domains,omitzero"`
	// Never fetch URLs whose host is one of these domains or a subdomain of one. Each
	// entry is a plain hostname like "ads.example.com" (no scheme, port, or path). At
	// most 64 entries; an empty list is rejected (omit the field instead). Cannot be
	// combined with allowed_domains.
	BlockedDomains []string `json:"blocked_domains,omitzero"`
	// Any of "web_fetch".
	Type BetaManagedAgentsWebFetchToolConfigParamsType `json:"type,omitzero"`
	// Must be "web_fetch".
	//
	// This field can be elided, and will marshal its zero value as "web_fetch".
	Name constant.WebFetch `json:"name" default:"web_fetch"`
	paramObj
}

func (r BetaManagedAgentsWebFetchToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsWebFetchToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsWebFetchToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *BetaManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *BetaManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *BetaManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[BetaManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type BetaManagedAgentsWebFetchToolConfigParamsType string

const (
	BetaManagedAgentsWebFetchToolConfigParamsTypeWebFetch BetaManagedAgentsWebFetchToolConfigParamsType = "web_fetch"
)

// Configuration for the web_search tool.
type BetaManagedAgentsWebSearchToolConfig struct {
	Enabled bool               `json:"enabled" api:"required"`
	Name    constant.WebSearch `json:"name" default:"web_search"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsWebSearchToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	Type             constant.WebSearch                                        `json:"type" default:"web_search"`
	AllowedDomains   []string                                                  `json:"allowed_domains"`
	BlockedDomains   []string                                                  `json:"blocked_domains"`
	// Approximate user location for search result localization.
	UserLocation BetaManagedAgentsUserLocation `json:"user_location" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		AllowedDomains   respjson.Field
		BlockedDomains   respjson.Field
		UserLocation     respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsWebSearchToolConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsWebSearchToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsWebSearchToolConfigPermissionPolicyUnion contains all possible
// properties and values from [BetaManagedAgentsAlwaysAllowPolicy],
// [BetaManagedAgentsAlwaysAskPolicy].
//
// Use the [BetaManagedAgentsWebSearchToolConfigPermissionPolicyUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsWebSearchToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsWebSearchToolConfigPermissionPolicy is implemented by each
// variant of [BetaManagedAgentsWebSearchToolConfigPermissionPolicyUnion] to add
// type safety for the return type of
// [BetaManagedAgentsWebSearchToolConfigPermissionPolicyUnion.AsAny]
type anyBetaManagedAgentsWebSearchToolConfigPermissionPolicy interface {
	implBetaManagedAgentsWebSearchToolConfigPermissionPolicyUnion()
}

func (BetaManagedAgentsAlwaysAllowPolicy) implBetaManagedAgentsWebSearchToolConfigPermissionPolicyUnion() {
}
func (BetaManagedAgentsAlwaysAskPolicy) implBetaManagedAgentsWebSearchToolConfigPermissionPolicyUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsWebSearchToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAlwaysAllowPolicy:
//	case anthropic.BetaManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsWebSearchToolConfigPermissionPolicyUnion) AsAny() anyBetaManagedAgentsWebSearchToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u BetaManagedAgentsWebSearchToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v BetaManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsWebSearchToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v BetaManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsWebSearchToolConfigPermissionPolicyUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *BetaManagedAgentsWebSearchToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for the web_search tool.
//
// The property Name is required.
type BetaManagedAgentsWebSearchToolConfigParams struct {
	// Whether this tool is enabled and available to Claude. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	// Only return search results whose host is one of these domains or a subdomain of
	// one. Each entry is a plain hostname like "docs.example.com" (no scheme or port;
	// an optional path suffix is accepted). At most 64 entries; an empty list is
	// rejected (omit the field instead). Cannot be combined with blocked_domains.
	AllowedDomains []string `json:"allowed_domains,omitzero"`
	// Never return search results whose host is one of these domains or a subdomain of
	// one. Each entry is a plain hostname like "ads.example.com" (no scheme or port;
	// an optional path suffix is accepted). At most 64 entries; an empty list is
	// rejected (omit the field instead). Cannot be combined with allowed_domains.
	BlockedDomains []string `json:"blocked_domains,omitzero"`
	// Any of "web_search".
	Type BetaManagedAgentsWebSearchToolConfigParamsType `json:"type,omitzero"`
	// Approximate user location for search result localization.
	UserLocation BetaManagedAgentsUserLocationParam `json:"user_location,omitzero"`
	// Must be "web_search".
	//
	// This field can be elided, and will marshal its zero value as "web_search".
	Name constant.WebSearch `json:"name" default:"web_search"`
	paramObj
}

func (r BetaManagedAgentsWebSearchToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsWebSearchToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsWebSearchToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *BetaManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *BetaManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *BetaManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[BetaManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type BetaManagedAgentsWebSearchToolConfigParamsType string

const (
	BetaManagedAgentsWebSearchToolConfigParamsTypeWebSearch BetaManagedAgentsWebSearchToolConfigParamsType = "web_search"
)

// Configuration for the write tool.
type BetaManagedAgentsWriteToolConfig struct {
	Enabled bool           `json:"enabled" api:"required"`
	Name    constant.Write `json:"name" default:"write"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsWriteToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	Type             constant.Write                                        `json:"type" default:"write"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsWriteToolConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsWriteToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsWriteToolConfigPermissionPolicyUnion contains all possible
// properties and values from [BetaManagedAgentsAlwaysAllowPolicy],
// [BetaManagedAgentsAlwaysAskPolicy].
//
// Use the [BetaManagedAgentsWriteToolConfigPermissionPolicyUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsWriteToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsWriteToolConfigPermissionPolicy is implemented by each
// variant of [BetaManagedAgentsWriteToolConfigPermissionPolicyUnion] to add type
// safety for the return type of
// [BetaManagedAgentsWriteToolConfigPermissionPolicyUnion.AsAny]
type anyBetaManagedAgentsWriteToolConfigPermissionPolicy interface {
	implBetaManagedAgentsWriteToolConfigPermissionPolicyUnion()
}

func (BetaManagedAgentsAlwaysAllowPolicy) implBetaManagedAgentsWriteToolConfigPermissionPolicyUnion() {
}
func (BetaManagedAgentsAlwaysAskPolicy) implBetaManagedAgentsWriteToolConfigPermissionPolicyUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsWriteToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAlwaysAllowPolicy:
//	case anthropic.BetaManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsWriteToolConfigPermissionPolicyUnion) AsAny() anyBetaManagedAgentsWriteToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u BetaManagedAgentsWriteToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v BetaManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsWriteToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v BetaManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsWriteToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsWriteToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for the write tool.
//
// The property Name is required.
type BetaManagedAgentsWriteToolConfigParams struct {
	// Whether this tool is enabled and available to Claude. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsWriteToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	// Any of "write".
	Type BetaManagedAgentsWriteToolConfigParamsType `json:"type,omitzero"`
	// Must be "write".
	//
	// This field can be elided, and will marshal its zero value as "write".
	Name constant.Write `json:"name" default:"write"`
	paramObj
}

func (r BetaManagedAgentsWriteToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsWriteToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsWriteToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsWriteToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *BetaManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *BetaManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsWriteToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *BetaManagedAgentsWriteToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsWriteToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsWriteToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsWriteToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[BetaManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type BetaManagedAgentsWriteToolConfigParamsType string

const (
	BetaManagedAgentsWriteToolConfigParamsTypeWrite BetaManagedAgentsWriteToolConfigParamsType = "write"
)

type BetaAgentNewParams struct {
	// Model identifier. Accepts the
	// [model string](https://platform.claude.com/docs/en/about-claude/models/overview#latest-models-comparison),
	// e.g. `claude-opus-5`, or a `model_config` object for additional configuration
	// control
	Model BetaManagedAgentsModelConfigParams `json:"model,omitzero" api:"required"`
	// Human-readable name for the agent.
	Name string `json:"name" api:"required"`
	// Description of what the agent does.
	Description param.Opt[string] `json:"description,omitzero"`
	// System prompt for the agent.
	System param.Opt[string] `json:"system,omitzero"`
	// MCP servers this agent connects to. Maximum 20. Names must be unique within the
	// array. Every server must be referenced by an `mcp_toolset` in `tools`;
	// unreferenced servers are rejected. See the
	// [MCP connector guide](https://platform.claude.com/docs/en/managed-agents/mcp-connector).
	MCPServers []BetaManagedAgentsURLMCPServerParams `json:"mcp_servers,omitzero"`
	// Arbitrary key-value metadata. Maximum 16 pairs, keys up to 64 chars, values up
	// to 512 chars.
	Metadata map[string]string `json:"metadata,omitzero"`
	// A coordinator topology: the session's primary thread orchestrates work by
	// spawning session threads, each running an agent drawn from the `agents` roster.
	Multiagent BetaManagedAgentsMultiagentParams `json:"multiagent,omitzero"`
	// Skills available to the agent.
	Skills []BetaManagedAgentsSkillParamsUnion `json:"skills,omitzero"`
	// Tool configurations available to the agent. Maximum of 128 tools across all
	// toolsets allowed.
	Tools []BetaAgentNewParamsToolUnion `json:"tools,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaAgentNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaAgentNewParamsToolUnion struct {
	OfAgentToolset20260401 *BetaManagedAgentsAgentToolset20260401Params `json:",omitzero,inline"`
	OfMCPToolset           *BetaManagedAgentsMCPToolsetParams           `json:",omitzero,inline"`
	OfCustom               *BetaManagedAgentsCustomToolParams           `json:",omitzero,inline"`
	paramUnion
}

func (u BetaAgentNewParamsToolUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAgentToolset20260401, u.OfMCPToolset, u.OfCustom)
}
func (u *BetaAgentNewParamsToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaAgentNewParamsToolUnion) asAny() any {
	if !param.IsOmitted(u.OfAgentToolset20260401) {
		return u.OfAgentToolset20260401
	} else if !param.IsOmitted(u.OfMCPToolset) {
		return u.OfMCPToolset
	} else if !param.IsOmitted(u.OfCustom) {
		return u.OfCustom
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaAgentNewParamsToolUnion) GetMCPServerName() *string {
	if vt := u.OfMCPToolset; vt != nil {
		return &vt.MCPServerName
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaAgentNewParamsToolUnion) GetDescription() *string {
	if vt := u.OfCustom; vt != nil {
		return &vt.Description
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaAgentNewParamsToolUnion) GetInputSchema() *BetaManagedAgentsCustomToolInputSchemaParam {
	if vt := u.OfCustom; vt != nil {
		return &vt.InputSchema
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaAgentNewParamsToolUnion) GetName() *string {
	if vt := u.OfCustom; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaAgentNewParamsToolUnion) GetType() *string {
	if vt := u.OfAgentToolset20260401; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMCPToolset; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCustom; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u BetaAgentNewParamsToolUnion) GetConfigs() (res betaAgentNewParamsToolUnionConfigs) {
	if vt := u.OfAgentToolset20260401; vt != nil {
		res.any = &vt.Configs
	} else if vt := u.OfMCPToolset; vt != nil {
		res.any = &vt.Configs
	}
	return
}

// Can have the runtime types [_[]BetaManagedAgentsAgentToolConfigParamsUnion],
// [_[]BetaManagedAgentsMCPToolConfigParams]
type betaAgentNewParamsToolUnionConfigs struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]anthropic.BetaManagedAgentsAgentToolConfigParamsUnion:
//	case *[]anthropic.BetaManagedAgentsMCPToolConfigParams:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u betaAgentNewParamsToolUnionConfigs) AsAny() any { return u.any }

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u BetaAgentNewParamsToolUnion) GetDefaultConfig() (res betaAgentNewParamsToolUnionDefaultConfig) {
	if vt := u.OfAgentToolset20260401; vt != nil {
		res.any = &vt.DefaultConfig
	} else if vt := u.OfMCPToolset; vt != nil {
		res.any = &vt.DefaultConfig
	}
	return
}

// Can have the runtime types [*BetaManagedAgentsAgentToolsetDefaultConfigParams],
// [*BetaManagedAgentsMCPToolsetDefaultConfigParams]
type betaAgentNewParamsToolUnionDefaultConfig struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *anthropic.BetaManagedAgentsAgentToolsetDefaultConfigParams:
//	case *anthropic.BetaManagedAgentsMCPToolsetDefaultConfigParams:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u betaAgentNewParamsToolUnionDefaultConfig) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u betaAgentNewParamsToolUnionDefaultConfig) GetEnabled() *bool {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsAgentToolsetDefaultConfigParams:
		return paramutil.AddrIfPresent(vt.Enabled)
	case *BetaManagedAgentsMCPToolsetDefaultConfigParams:
		return paramutil.AddrIfPresent(vt.Enabled)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u betaAgentNewParamsToolUnionDefaultConfig) GetPermissionPolicy() (res betaAgentNewParamsToolUnionDefaultConfigPermissionPolicy) {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsAgentToolsetDefaultConfigParams:
		res.any = vt.PermissionPolicy
	case *BetaManagedAgentsMCPToolsetDefaultConfigParams:
		res.any = vt.PermissionPolicy
	}
	return res
}

// Can have the runtime types [*BetaManagedAgentsAlwaysAllowPolicyParam],
// [*BetaManagedAgentsAlwaysAskPolicyParam]
type betaAgentNewParamsToolUnionDefaultConfigPermissionPolicy struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *anthropic.BetaManagedAgentsAlwaysAllowPolicyParam:
//	case *anthropic.BetaManagedAgentsAlwaysAskPolicyParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u betaAgentNewParamsToolUnionDefaultConfigPermissionPolicy) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u betaAgentNewParamsToolUnionDefaultConfigPermissionPolicy) GetType() *string {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *BetaManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaAgentNewParamsToolUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsAgentToolset20260401Params]("agent_toolset_20260401"),
		apijson.Discriminator[BetaManagedAgentsMCPToolsetParams]("mcp_toolset"),
		apijson.Discriminator[BetaManagedAgentsCustomToolParams]("custom"),
	)
}

type BetaAgentGetParams struct {
	// Agent version. Omit for the most recent version. Must be at least 1 if
	// specified.
	Version param.Opt[int64] `query:"version,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentGetParams]'s query parameters as `url.Values`.
func (r BetaAgentGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaAgentUpdateParams struct {
	// Description. Omit to preserve; send empty string or null to clear.
	Description param.Opt[string] `json:"description,omitzero"`
	// System prompt. Omit to preserve; send empty string or null to clear.
	System param.Opt[string] `json:"system,omitzero"`
	// Human-readable name. Must be non-empty. Omit to preserve. Cannot be cleared.
	Name param.Opt[string] `json:"name,omitzero"`
	// The agent's current version, used to prevent concurrent overwrites. Obtain this
	// value from a create or retrieve response. Must be at least 1 if specified. When
	// supplied, the request fails if it does not match the server's current version;
	// omit to apply the update unconditionally.
	Version param.Opt[int64] `json:"version,omitzero"`
	// MCP servers. Full replacement. Omit to preserve; send empty array or `null` to
	// clear. Names must be unique. Maximum 20. Every server must be referenced by an
	// `mcp_toolset` in the agent's resulting `tools`; unreferenced servers are
	// rejected. See the
	// [MCP connector guide](https://platform.claude.com/docs/en/managed-agents/mcp-connector).
	MCPServers []BetaManagedAgentsURLMCPServerParams `json:"mcp_servers,omitzero"`
	// Metadata patch. Set a key to a string to upsert it, or to null to delete it.
	// Omit the field to preserve. The stored bag is limited to 16 keys (up to 64 chars
	// each) with values up to 512 chars.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Skills. Full replacement. Omit to preserve; send empty array or null to clear.
	Skills []BetaManagedAgentsSkillParamsUnion `json:"skills,omitzero"`
	// Tool configurations available to the agent. Full replacement. Omit to preserve;
	// send empty array or null to clear. Maximum of 128 tools across all toolsets
	// allowed.
	Tools []BetaAgentUpdateParamsToolUnion `json:"tools,omitzero"`
	// Model identifier. Accepts the
	// [model string](https://platform.claude.com/docs/en/about-claude/models/overview#latest-models-comparison),
	// e.g. `claude-opus-5`, or a `model_config` object for additional configuration
	// control. Omit to preserve. Cannot be cleared.
	Model BetaManagedAgentsModelConfigParams `json:"model,omitzero"`
	// A coordinator topology: the session's primary thread orchestrates work by
	// spawning session threads, each running an agent drawn from the `agents` roster.
	Multiagent BetaManagedAgentsMultiagentParams `json:"multiagent,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaAgentUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaAgentUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaAgentUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaAgentUpdateParamsToolUnion struct {
	OfAgentToolset20260401 *BetaManagedAgentsAgentToolset20260401Params `json:",omitzero,inline"`
	OfMCPToolset           *BetaManagedAgentsMCPToolsetParams           `json:",omitzero,inline"`
	OfCustom               *BetaManagedAgentsCustomToolParams           `json:",omitzero,inline"`
	paramUnion
}

func (u BetaAgentUpdateParamsToolUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAgentToolset20260401, u.OfMCPToolset, u.OfCustom)
}
func (u *BetaAgentUpdateParamsToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaAgentUpdateParamsToolUnion) asAny() any {
	if !param.IsOmitted(u.OfAgentToolset20260401) {
		return u.OfAgentToolset20260401
	} else if !param.IsOmitted(u.OfMCPToolset) {
		return u.OfMCPToolset
	} else if !param.IsOmitted(u.OfCustom) {
		return u.OfCustom
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaAgentUpdateParamsToolUnion) GetMCPServerName() *string {
	if vt := u.OfMCPToolset; vt != nil {
		return &vt.MCPServerName
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaAgentUpdateParamsToolUnion) GetDescription() *string {
	if vt := u.OfCustom; vt != nil {
		return &vt.Description
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaAgentUpdateParamsToolUnion) GetInputSchema() *BetaManagedAgentsCustomToolInputSchemaParam {
	if vt := u.OfCustom; vt != nil {
		return &vt.InputSchema
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaAgentUpdateParamsToolUnion) GetName() *string {
	if vt := u.OfCustom; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaAgentUpdateParamsToolUnion) GetType() *string {
	if vt := u.OfAgentToolset20260401; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMCPToolset; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCustom; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u BetaAgentUpdateParamsToolUnion) GetConfigs() (res betaAgentUpdateParamsToolUnionConfigs) {
	if vt := u.OfAgentToolset20260401; vt != nil {
		res.any = &vt.Configs
	} else if vt := u.OfMCPToolset; vt != nil {
		res.any = &vt.Configs
	}
	return
}

// Can have the runtime types [_[]BetaManagedAgentsAgentToolConfigParamsUnion],
// [_[]BetaManagedAgentsMCPToolConfigParams]
type betaAgentUpdateParamsToolUnionConfigs struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]anthropic.BetaManagedAgentsAgentToolConfigParamsUnion:
//	case *[]anthropic.BetaManagedAgentsMCPToolConfigParams:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u betaAgentUpdateParamsToolUnionConfigs) AsAny() any { return u.any }

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u BetaAgentUpdateParamsToolUnion) GetDefaultConfig() (res betaAgentUpdateParamsToolUnionDefaultConfig) {
	if vt := u.OfAgentToolset20260401; vt != nil {
		res.any = &vt.DefaultConfig
	} else if vt := u.OfMCPToolset; vt != nil {
		res.any = &vt.DefaultConfig
	}
	return
}

// Can have the runtime types [*BetaManagedAgentsAgentToolsetDefaultConfigParams],
// [*BetaManagedAgentsMCPToolsetDefaultConfigParams]
type betaAgentUpdateParamsToolUnionDefaultConfig struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *anthropic.BetaManagedAgentsAgentToolsetDefaultConfigParams:
//	case *anthropic.BetaManagedAgentsMCPToolsetDefaultConfigParams:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u betaAgentUpdateParamsToolUnionDefaultConfig) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u betaAgentUpdateParamsToolUnionDefaultConfig) GetEnabled() *bool {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsAgentToolsetDefaultConfigParams:
		return paramutil.AddrIfPresent(vt.Enabled)
	case *BetaManagedAgentsMCPToolsetDefaultConfigParams:
		return paramutil.AddrIfPresent(vt.Enabled)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u betaAgentUpdateParamsToolUnionDefaultConfig) GetPermissionPolicy() (res betaAgentUpdateParamsToolUnionDefaultConfigPermissionPolicy) {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsAgentToolsetDefaultConfigParams:
		res.any = vt.PermissionPolicy
	case *BetaManagedAgentsMCPToolsetDefaultConfigParams:
		res.any = vt.PermissionPolicy
	}
	return res
}

// Can have the runtime types [*BetaManagedAgentsAlwaysAllowPolicyParam],
// [*BetaManagedAgentsAlwaysAskPolicyParam]
type betaAgentUpdateParamsToolUnionDefaultConfigPermissionPolicy struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *anthropic.BetaManagedAgentsAlwaysAllowPolicyParam:
//	case *anthropic.BetaManagedAgentsAlwaysAskPolicyParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u betaAgentUpdateParamsToolUnionDefaultConfigPermissionPolicy) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u betaAgentUpdateParamsToolUnionDefaultConfigPermissionPolicy) GetType() *string {
	switch vt := u.any.(type) {
	case *BetaManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *BetaManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaAgentUpdateParamsToolUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsAgentToolset20260401Params]("agent_toolset_20260401"),
		apijson.Discriminator[BetaManagedAgentsMCPToolsetParams]("mcp_toolset"),
		apijson.Discriminator[BetaManagedAgentsCustomToolParams]("custom"),
	)
}

type BetaAgentListParams struct {
	// Return agents created at or after this time (inclusive).
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return agents created at or before this time (inclusive).
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// Include archived agents in results. Defaults to false.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum results per page. Default 20, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor from a previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaAgentListParams]'s query parameters as `url.Values`.
func (r BetaAgentListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaAgentArchiveParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
