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
	Model  BetaManagedAgentsModelConfig       `json:"model" api:"required"`
	Name   string                             `json:"name" api:"required"`
	Skills []BetaManagedAgentsAgentSkillUnion `json:"skills" api:"required"`
	System string                             `json:"system" api:"required"`
	Tools  []BetaManagedAgentsAgentToolUnion  `json:"tools" api:"required"`
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
	// This field is a union of [[]BetaManagedAgentsAgentToolConfig],
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

// Configuration for a specific agent tool.
type BetaManagedAgentsAgentToolConfig struct {
	Enabled bool `json:"enabled" api:"required"`
	// Built-in agent tool identifier.
	//
	// Any of "bash", "edit", "read", "write", "glob", "grep", "web_fetch",
	// "web_search".
	Name BetaManagedAgentsAgentToolConfigName `json:"name" api:"required"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsAgentToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
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
func (r BetaManagedAgentsAgentToolConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Built-in agent tool identifier.
type BetaManagedAgentsAgentToolConfigName string

const (
	BetaManagedAgentsAgentToolConfigNameBash      BetaManagedAgentsAgentToolConfigName = "bash"
	BetaManagedAgentsAgentToolConfigNameEdit      BetaManagedAgentsAgentToolConfigName = "edit"
	BetaManagedAgentsAgentToolConfigNameRead      BetaManagedAgentsAgentToolConfigName = "read"
	BetaManagedAgentsAgentToolConfigNameWrite     BetaManagedAgentsAgentToolConfigName = "write"
	BetaManagedAgentsAgentToolConfigNameGlob      BetaManagedAgentsAgentToolConfigName = "glob"
	BetaManagedAgentsAgentToolConfigNameGrep      BetaManagedAgentsAgentToolConfigName = "grep"
	BetaManagedAgentsAgentToolConfigNameWebFetch  BetaManagedAgentsAgentToolConfigName = "web_fetch"
	BetaManagedAgentsAgentToolConfigNameWebSearch BetaManagedAgentsAgentToolConfigName = "web_search"
)

// BetaManagedAgentsAgentToolConfigPermissionPolicyUnion contains all possible
// properties and values from [BetaManagedAgentsAlwaysAllowPolicy],
// [BetaManagedAgentsAlwaysAskPolicy].
//
// Use the [BetaManagedAgentsAgentToolConfigPermissionPolicyUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsAgentToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyBetaManagedAgentsAgentToolConfigPermissionPolicy is implemented by each
// variant of [BetaManagedAgentsAgentToolConfigPermissionPolicyUnion] to add type
// safety for the return type of
// [BetaManagedAgentsAgentToolConfigPermissionPolicyUnion.AsAny]
type anyBetaManagedAgentsAgentToolConfigPermissionPolicy interface {
	implBetaManagedAgentsAgentToolConfigPermissionPolicyUnion()
}

func (BetaManagedAgentsAlwaysAllowPolicy) implBetaManagedAgentsAgentToolConfigPermissionPolicyUnion() {
}
func (BetaManagedAgentsAlwaysAskPolicy) implBetaManagedAgentsAgentToolConfigPermissionPolicyUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsAgentToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAlwaysAllowPolicy:
//	case anthropic.BetaManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsAgentToolConfigPermissionPolicyUnion) AsAny() anyBetaManagedAgentsAgentToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u BetaManagedAgentsAgentToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v BetaManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsAgentToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v BetaManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsAgentToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsAgentToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for a specific tool within a toolset.
//
// The property Name is required.
type BetaManagedAgentsAgentToolConfigParams struct {
	// Built-in agent tool identifier.
	//
	// Any of "bash", "edit", "read", "write", "glob", "grep", "web_fetch",
	// "web_search".
	Name BetaManagedAgentsAgentToolConfigParamsName `json:"name,omitzero" api:"required"`
	// Whether this tool is enabled and available to Claude. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy BetaManagedAgentsAgentToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	paramObj
}

func (r BetaManagedAgentsAgentToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsAgentToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsAgentToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Built-in agent tool identifier.
type BetaManagedAgentsAgentToolConfigParamsName string

const (
	BetaManagedAgentsAgentToolConfigParamsNameBash      BetaManagedAgentsAgentToolConfigParamsName = "bash"
	BetaManagedAgentsAgentToolConfigParamsNameEdit      BetaManagedAgentsAgentToolConfigParamsName = "edit"
	BetaManagedAgentsAgentToolConfigParamsNameRead      BetaManagedAgentsAgentToolConfigParamsName = "read"
	BetaManagedAgentsAgentToolConfigParamsNameWrite     BetaManagedAgentsAgentToolConfigParamsName = "write"
	BetaManagedAgentsAgentToolConfigParamsNameGlob      BetaManagedAgentsAgentToolConfigParamsName = "glob"
	BetaManagedAgentsAgentToolConfigParamsNameGrep      BetaManagedAgentsAgentToolConfigParamsName = "grep"
	BetaManagedAgentsAgentToolConfigParamsNameWebFetch  BetaManagedAgentsAgentToolConfigParamsName = "web_fetch"
	BetaManagedAgentsAgentToolConfigParamsNameWebSearch BetaManagedAgentsAgentToolConfigParamsName = "web_search"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsAgentToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *BetaManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *BetaManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsAgentToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *BetaManagedAgentsAgentToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsAgentToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsAgentToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsAgentToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[BetaManagedAgentsAlwaysAskPolicyParam]("always_ask"),
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
	Configs []BetaManagedAgentsAgentToolConfig `json:"configs" api:"required"`
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

// Configuration for built-in agent tools. Use this to enable or disable groups of
// tools available to the agent.
//
// The property Type is required.
type BetaManagedAgentsAgentToolset20260401Params struct {
	// Any of "agent_toolset_20260401".
	Type BetaManagedAgentsAgentToolset20260401ParamsType `json:"type,omitzero" api:"required"`
	// Per-tool configuration overrides.
	Configs []BetaManagedAgentsAgentToolConfigParams `json:"configs,omitzero"`
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
	// JSON Schema properties defining the tool's input parameters.
	Properties map[string]any `json:"properties" api:"nullable"`
	// List of required property names.
	Required []string `json:"required"`
	// Must be 'object' for tool input schemas.
	//
	// Any of "object".
	Type BetaManagedAgentsCustomToolInputSchemaType `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Properties  respjson.Field
		Required    respjson.Field
		Type        respjson.Field
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

// Must be 'object' for tool input schemas.
type BetaManagedAgentsCustomToolInputSchemaType string

const (
	BetaManagedAgentsCustomToolInputSchemaTypeObject BetaManagedAgentsCustomToolInputSchemaType = "object"
)

// JSON Schema for custom tool input parameters.
type BetaManagedAgentsCustomToolInputSchemaParam struct {
	// JSON Schema properties defining the tool's input parameters.
	Properties map[string]any `json:"properties,omitzero"`
	// List of required property names.
	Required []string `json:"required,omitzero"`
	// Must be 'object' for tool input schemas.
	//
	// Any of "object".
	Type BetaManagedAgentsCustomToolInputSchemaType `json:"type,omitzero"`
	paramObj
}

func (r BetaManagedAgentsCustomToolInputSchemaParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsCustomToolInputSchemaParam
	return param.MarshalObject(r, (*shadow)(&r))
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
	// use the tool. 1-1024 characters.
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

// The model that will power your agent.\n\nSee
// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
// details and options.
type BetaManagedAgentsModel = string

const (
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
	// The model that will power your agent.\n\nSee
	// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	ID BetaManagedAgentsModel `json:"id" api:"required"`
	// Inference speed mode. `fast` provides significantly faster output token
	// generation at premium pricing. Not all models support `fast`; invalid
	// combinations are rejected at create time.
	//
	// Any of "standard", "fast".
	Speed BetaManagedAgentsModelConfigSpeed `json:"speed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Speed       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsModelConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsModelConfig) UnmarshalJSON(data []byte) error {
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
	// The model that will power your agent.\n\nSee
	// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	ID BetaManagedAgentsModel `json:"id,omitzero" api:"required"`
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

// Inference speed mode. `fast` provides significantly faster output token
// generation at premium pricing. Not all models support `fast`; invalid
// combinations are rejected at create time.
type BetaManagedAgentsModelConfigParamsSpeed string

const (
	BetaManagedAgentsModelConfigParamsSpeedStandard BetaManagedAgentsModelConfigParamsSpeed = "standard"
	BetaManagedAgentsModelConfigParamsSpeedFast     BetaManagedAgentsModelConfigParamsSpeed = "fast"
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

type BetaAgentNewParams struct {
	// Model identifier. Accepts the
	// [model string](https://platform.claude.com/docs/en/about-claude/models/overview#latest-models-comparison),
	// e.g. `claude-opus-4-6`, or a `model_config` object for additional configuration
	// control
	Model BetaManagedAgentsModelConfigParams `json:"model,omitzero" api:"required"`
	// Human-readable name for the agent. 1-256 characters.
	Name string `json:"name" api:"required"`
	// Description of what the agent does. Up to 2048 characters.
	Description param.Opt[string] `json:"description,omitzero"`
	// System prompt for the agent. Up to 100,000 characters.
	System param.Opt[string] `json:"system,omitzero"`
	// MCP servers this agent connects to. Maximum 20. Names must be unique within the
	// array.
	MCPServers []BetaManagedAgentsURLMCPServerParams `json:"mcp_servers,omitzero"`
	// Arbitrary key-value metadata. Maximum 16 pairs, keys up to 64 chars, values up
	// to 512 chars.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Skills available to the agent. Maximum 20.
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

// Can have the runtime types [_[]BetaManagedAgentsAgentToolConfigParams],
// [_[]BetaManagedAgentsMCPToolConfigParams]
type betaAgentNewParamsToolUnionConfigs struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]anthropic.BetaManagedAgentsAgentToolConfigParams:
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
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaAgentUpdateParams struct {
	// The agent's current version, used to prevent concurrent overwrites. Obtain this
	// value from a create or retrieve response. The request fails if this does not
	// match the server's current version.
	Version int64 `json:"version" api:"required"`
	// Description. Up to 2048 characters. Omit to preserve; send empty string or null
	// to clear.
	Description param.Opt[string] `json:"description,omitzero"`
	// System prompt. Up to 100,000 characters. Omit to preserve; send empty string or
	// null to clear.
	System param.Opt[string] `json:"system,omitzero"`
	// Human-readable name. 1-256 characters. Omit to preserve. Cannot be cleared.
	Name param.Opt[string] `json:"name,omitzero"`
	// MCP servers. Full replacement. Omit to preserve; send empty array or null to
	// clear. Names must be unique. Maximum 20.
	MCPServers []BetaManagedAgentsURLMCPServerParams `json:"mcp_servers,omitzero"`
	// Metadata patch. Set a key to a string to upsert it, or to null to delete it.
	// Omit the field to preserve. The stored bag is limited to 16 keys (up to 64 chars
	// each) with values up to 512 chars.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Skills. Full replacement. Omit to preserve; send empty array or null to clear.
	// Maximum 20.
	Skills []BetaManagedAgentsSkillParamsUnion `json:"skills,omitzero"`
	// Tool configurations available to the agent. Full replacement. Omit to preserve;
	// send empty array or null to clear. Maximum of 128 tools across all toolsets
	// allowed.
	Tools []BetaAgentUpdateParamsToolUnion `json:"tools,omitzero"`
	// Model identifier. Accepts the
	// [model string](https://platform.claude.com/docs/en/about-claude/models/overview#latest-models-comparison),
	// e.g. `claude-opus-4-6`, or a `model_config` object for additional configuration
	// control. Omit to preserve. Cannot be cleared.
	Model BetaManagedAgentsModelConfigParams `json:"model,omitzero"`
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

// Can have the runtime types [_[]BetaManagedAgentsAgentToolConfigParams],
// [_[]BetaManagedAgentsMCPToolConfigParams]
type betaAgentUpdateParamsToolUnionConfigs struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]anthropic.BetaManagedAgentsAgentToolConfigParams:
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
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaAgentArchiveParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
