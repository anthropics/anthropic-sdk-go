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

// BetaSessionService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaSessionService] method instead.
type BetaSessionService struct {
	Options   []option.RequestOption
	Events    BetaSessionEventService
	Resources BetaSessionResourceService
}

// NewBetaSessionService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaSessionService(opts ...option.RequestOption) (r BetaSessionService) {
	r = BetaSessionService{}
	r.Options = opts
	r.Events = NewBetaSessionEventService(opts...)
	r.Resources = NewBetaSessionResourceService(opts...)
	return
}

// Create Session
func (r *BetaSessionService) New(ctx context.Context, params BetaSessionNewParams, opts ...option.RequestOption) (res *BetaManagedAgentsSession, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	path := "v1/sessions?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get Session
func (r *BetaSessionService) Get(ctx context.Context, sessionID string, query BetaSessionGetParams, opts ...option.RequestOption) (res *BetaManagedAgentsSession, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/sessions/%s?beta=true", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Session
func (r *BetaSessionService) Update(ctx context.Context, sessionID string, params BetaSessionUpdateParams, opts ...option.RequestOption) (res *BetaManagedAgentsSession, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/sessions/%s?beta=true", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List Sessions
func (r *BetaSessionService) List(ctx context.Context, params BetaSessionListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaManagedAgentsSession], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/sessions?beta=true"
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

// List Sessions
func (r *BetaSessionService) ListAutoPaging(ctx context.Context, params BetaSessionListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsSession] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Delete Session
func (r *BetaSessionService) Delete(ctx context.Context, sessionID string, body BetaSessionDeleteParams, opts ...option.RequestOption) (res *BetaManagedAgentsDeletedSession, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/sessions/%s?beta=true", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Archive Session
func (r *BetaSessionService) Archive(ctx context.Context, sessionID string, body BetaSessionArchiveParams, opts ...option.RequestOption) (res *BetaManagedAgentsSession, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/sessions/%s/archive?beta=true", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Specification for an Agent. Provide a specific `version` or use the short-form
// `agent="agent_id"` for the most recent version
//
// The properties ID, Type are required.
type BetaManagedAgentsAgentParams struct {
	// The `agent` ID.
	ID string `json:"id" api:"required"`
	// Any of "agent".
	Type BetaManagedAgentsAgentParamsType `json:"type,omitzero" api:"required"`
	// The specific `agent` version to use. Omit to use the latest version. Must be at
	// least 1 if specified.
	Version param.Opt[int64] `json:"version,omitzero"`
	paramObj
}

func (r BetaManagedAgentsAgentParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsAgentParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsAgentParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAgentParamsType string

const (
	BetaManagedAgentsAgentParamsTypeAgent BetaManagedAgentsAgentParamsType = "agent"
)

type BetaManagedAgentsBranchCheckout struct {
	// Branch name to check out.
	Name string `json:"name" api:"required"`
	// Any of "branch".
	Type BetaManagedAgentsBranchCheckoutType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsBranchCheckout) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsBranchCheckout) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsBranchCheckout to a
// BetaManagedAgentsBranchCheckoutParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsBranchCheckoutParam.Overrides()
func (r BetaManagedAgentsBranchCheckout) ToParam() BetaManagedAgentsBranchCheckoutParam {
	return param.Override[BetaManagedAgentsBranchCheckoutParam](json.RawMessage(r.RawJSON()))
}

type BetaManagedAgentsBranchCheckoutType string

const (
	BetaManagedAgentsBranchCheckoutTypeBranch BetaManagedAgentsBranchCheckoutType = "branch"
)

// The properties Name, Type are required.
type BetaManagedAgentsBranchCheckoutParam struct {
	// Branch name to check out.
	Name string `json:"name" api:"required"`
	// Any of "branch".
	Type BetaManagedAgentsBranchCheckoutType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsBranchCheckoutParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsBranchCheckoutParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsBranchCheckoutParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Prompt-cache creation token usage broken down by cache lifetime.
type BetaManagedAgentsCacheCreationUsage struct {
	// Tokens used to create 1-hour ephemeral cache entries.
	Ephemeral1hInputTokens int64 `json:"ephemeral_1h_input_tokens"`
	// Tokens used to create 5-minute ephemeral cache entries.
	Ephemeral5mInputTokens int64 `json:"ephemeral_5m_input_tokens"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Ephemeral1hInputTokens respjson.Field
		Ephemeral5mInputTokens respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsCacheCreationUsage) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsCacheCreationUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsCommitCheckout struct {
	// Full commit SHA to check out.
	Sha string `json:"sha" api:"required"`
	// Any of "commit".
	Type BetaManagedAgentsCommitCheckoutType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Sha         respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsCommitCheckout) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsCommitCheckout) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaManagedAgentsCommitCheckout to a
// BetaManagedAgentsCommitCheckoutParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaManagedAgentsCommitCheckoutParam.Overrides()
func (r BetaManagedAgentsCommitCheckout) ToParam() BetaManagedAgentsCommitCheckoutParam {
	return param.Override[BetaManagedAgentsCommitCheckoutParam](json.RawMessage(r.RawJSON()))
}

type BetaManagedAgentsCommitCheckoutType string

const (
	BetaManagedAgentsCommitCheckoutTypeCommit BetaManagedAgentsCommitCheckoutType = "commit"
)

// The properties Sha, Type are required.
type BetaManagedAgentsCommitCheckoutParam struct {
	// Full commit SHA to check out.
	Sha string `json:"sha" api:"required"`
	// Any of "commit".
	Type BetaManagedAgentsCommitCheckoutType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaManagedAgentsCommitCheckoutParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsCommitCheckoutParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsCommitCheckoutParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Confirmation that a `session` has been permanently deleted.
type BetaManagedAgentsDeletedSession struct {
	ID string `json:"id" api:"required"`
	// Any of "session_deleted".
	Type BetaManagedAgentsDeletedSessionType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsDeletedSession) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsDeletedSession) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsDeletedSessionType string

const (
	BetaManagedAgentsDeletedSessionTypeSessionDeleted BetaManagedAgentsDeletedSessionType = "session_deleted"
)

// Mount a file uploaded via the Files API into the session.
//
// The properties FileID, Type are required.
type BetaManagedAgentsFileResourceParams struct {
	// ID of a previously uploaded file.
	FileID string `json:"file_id" api:"required"`
	// Any of "file".
	Type BetaManagedAgentsFileResourceParamsType `json:"type,omitzero" api:"required"`
	// Mount path in the container. Defaults to `/mnt/session/uploads/<file_id>`.
	MountPath param.Opt[string] `json:"mount_path,omitzero"`
	paramObj
}

func (r BetaManagedAgentsFileResourceParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsFileResourceParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsFileResourceParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsFileResourceParamsType string

const (
	BetaManagedAgentsFileResourceParamsTypeFile BetaManagedAgentsFileResourceParamsType = "file"
)

// Mount a GitHub repository into the session's container.
//
// The properties AuthorizationToken, Type, URL are required.
type BetaManagedAgentsGitHubRepositoryResourceParams struct {
	// GitHub authorization token used to clone the repository.
	AuthorizationToken string `json:"authorization_token" api:"required"`
	// Any of "github_repository".
	Type BetaManagedAgentsGitHubRepositoryResourceParamsType `json:"type,omitzero" api:"required"`
	// Github URL of the repository
	URL string `json:"url" api:"required"`
	// Mount path in the container. Defaults to `/workspace/<repo-name>`.
	MountPath param.Opt[string] `json:"mount_path,omitzero"`
	// Branch or commit to check out. Defaults to the repository's default branch.
	Checkout BetaManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion `json:"checkout,omitzero"`
	paramObj
}

func (r BetaManagedAgentsGitHubRepositoryResourceParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsGitHubRepositoryResourceParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsGitHubRepositoryResourceParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsGitHubRepositoryResourceParamsType string

const (
	BetaManagedAgentsGitHubRepositoryResourceParamsTypeGitHubRepository BetaManagedAgentsGitHubRepositoryResourceParamsType = "github_repository"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion struct {
	OfBranch *BetaManagedAgentsBranchCheckoutParam `json:",omitzero,inline"`
	OfCommit *BetaManagedAgentsCommitCheckoutParam `json:",omitzero,inline"`
	paramUnion
}

func (u BetaManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfBranch, u.OfCommit)
}
func (u *BetaManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion) asAny() any {
	if !param.IsOmitted(u.OfBranch) {
		return u.OfBranch
	} else if !param.IsOmitted(u.OfCommit) {
		return u.OfCommit
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion) GetName() *string {
	if vt := u.OfBranch; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion) GetSha() *string {
	if vt := u.OfCommit; vt != nil {
		return &vt.Sha
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion) GetType() *string {
	if vt := u.OfBranch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCommit; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsBranchCheckoutParam]("branch"),
		apijson.Discriminator[BetaManagedAgentsCommitCheckoutParam]("commit"),
	)
}

// A Managed Agents `session`.
type BetaManagedAgentsSession struct {
	ID string `json:"id" api:"required"`
	// Resolved `agent` definition for a `session`. Snapshot of the `agent` at
	// `session` creation time.
	Agent BetaManagedAgentsSessionAgent `json:"agent" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	CreatedAt     time.Time                               `json:"created_at" api:"required" format:"date-time"`
	EnvironmentID string                                  `json:"environment_id" api:"required"`
	Metadata      map[string]string                       `json:"metadata" api:"required"`
	Resources     []BetaManagedAgentsSessionResourceUnion `json:"resources" api:"required"`
	// Timing statistics for a session.
	Stats BetaManagedAgentsSessionStats `json:"stats" api:"required"`
	// SessionStatus enum
	//
	// Any of "rescheduling", "running", "idle", "terminated".
	Status BetaManagedAgentsSessionStatus `json:"status" api:"required"`
	Title  string                         `json:"title" api:"required"`
	// Any of "session".
	Type BetaManagedAgentsSessionType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// Cumulative token usage for a session across all turns.
	Usage BetaManagedAgentsSessionUsage `json:"usage" api:"required"`
	// Vault IDs attached to the session at creation. Empty when no vaults were
	// supplied.
	VaultIDs []string `json:"vault_ids" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		Agent         respjson.Field
		ArchivedAt    respjson.Field
		CreatedAt     respjson.Field
		EnvironmentID respjson.Field
		Metadata      respjson.Field
		Resources     respjson.Field
		Stats         respjson.Field
		Status        respjson.Field
		Title         respjson.Field
		Type          respjson.Field
		UpdatedAt     respjson.Field
		Usage         respjson.Field
		VaultIDs      respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSession) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSession) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SessionStatus enum
type BetaManagedAgentsSessionStatus string

const (
	BetaManagedAgentsSessionStatusRescheduling BetaManagedAgentsSessionStatus = "rescheduling"
	BetaManagedAgentsSessionStatusRunning      BetaManagedAgentsSessionStatus = "running"
	BetaManagedAgentsSessionStatusIdle         BetaManagedAgentsSessionStatus = "idle"
	BetaManagedAgentsSessionStatusTerminated   BetaManagedAgentsSessionStatus = "terminated"
)

type BetaManagedAgentsSessionType string

const (
	BetaManagedAgentsSessionTypeSession BetaManagedAgentsSessionType = "session"
)

// Resolved `agent` definition for a `session`. Snapshot of the `agent` at
// `session` creation time.
type BetaManagedAgentsSessionAgent struct {
	ID          string                                    `json:"id" api:"required"`
	Description string                                    `json:"description" api:"required"`
	MCPServers  []BetaManagedAgentsMCPServerURLDefinition `json:"mcp_servers" api:"required"`
	// Model identifier and configuration.
	Model  BetaManagedAgentsModelConfig              `json:"model" api:"required"`
	Name   string                                    `json:"name" api:"required"`
	Skills []BetaManagedAgentsSessionAgentSkillUnion `json:"skills" api:"required"`
	System string                                    `json:"system" api:"required"`
	Tools  []BetaManagedAgentsSessionAgentToolUnion  `json:"tools" api:"required"`
	// Any of "agent".
	Type    BetaManagedAgentsSessionAgentType `json:"type" api:"required"`
	Version int64                             `json:"version" api:"required"`
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
func (r BetaManagedAgentsSessionAgent) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionAgent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSessionAgentSkillUnion contains all possible properties and
// values from [BetaManagedAgentsAnthropicSkill], [BetaManagedAgentsCustomSkill].
//
// Use the [BetaManagedAgentsSessionAgentSkillUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsSessionAgentSkillUnion struct {
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

// anyBetaManagedAgentsSessionAgentSkill is implemented by each variant of
// [BetaManagedAgentsSessionAgentSkillUnion] to add type safety for the return type
// of [BetaManagedAgentsSessionAgentSkillUnion.AsAny]
type anyBetaManagedAgentsSessionAgentSkill interface {
	implBetaManagedAgentsSessionAgentSkillUnion()
}

func (BetaManagedAgentsAnthropicSkill) implBetaManagedAgentsSessionAgentSkillUnion() {}
func (BetaManagedAgentsCustomSkill) implBetaManagedAgentsSessionAgentSkillUnion()    {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsSessionAgentSkillUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAnthropicSkill:
//	case anthropic.BetaManagedAgentsCustomSkill:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsSessionAgentSkillUnion) AsAny() anyBetaManagedAgentsSessionAgentSkill {
	switch u.Type {
	case "anthropic":
		return u.AsAnthropic()
	case "custom":
		return u.AsCustom()
	}
	return nil
}

func (u BetaManagedAgentsSessionAgentSkillUnion) AsAnthropic() (v BetaManagedAgentsAnthropicSkill) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionAgentSkillUnion) AsCustom() (v BetaManagedAgentsCustomSkill) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsSessionAgentSkillUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsSessionAgentSkillUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSessionAgentToolUnion contains all possible properties and
// values from [BetaManagedAgentsAgentToolset20260401],
// [BetaManagedAgentsMCPToolset], [BetaManagedAgentsCustomTool].
//
// Use the [BetaManagedAgentsSessionAgentToolUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsSessionAgentToolUnion struct {
	// This field is a union of [[]BetaManagedAgentsAgentToolConfig],
	// [[]BetaManagedAgentsMCPToolConfig]
	Configs BetaManagedAgentsSessionAgentToolUnionConfigs `json:"configs"`
	// This field is a union of [BetaManagedAgentsAgentToolsetDefaultConfig],
	// [BetaManagedAgentsMCPToolsetDefaultConfig]
	DefaultConfig BetaManagedAgentsSessionAgentToolUnionDefaultConfig `json:"default_config"`
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

// anyBetaManagedAgentsSessionAgentTool is implemented by each variant of
// [BetaManagedAgentsSessionAgentToolUnion] to add type safety for the return type
// of [BetaManagedAgentsSessionAgentToolUnion.AsAny]
type anyBetaManagedAgentsSessionAgentTool interface {
	implBetaManagedAgentsSessionAgentToolUnion()
}

func (BetaManagedAgentsAgentToolset20260401) implBetaManagedAgentsSessionAgentToolUnion() {}
func (BetaManagedAgentsMCPToolset) implBetaManagedAgentsSessionAgentToolUnion()           {}
func (BetaManagedAgentsCustomTool) implBetaManagedAgentsSessionAgentToolUnion()           {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsSessionAgentToolUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsAgentToolset20260401:
//	case anthropic.BetaManagedAgentsMCPToolset:
//	case anthropic.BetaManagedAgentsCustomTool:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsSessionAgentToolUnion) AsAny() anyBetaManagedAgentsSessionAgentTool {
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

func (u BetaManagedAgentsSessionAgentToolUnion) AsAgentToolset20260401() (v BetaManagedAgentsAgentToolset20260401) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionAgentToolUnion) AsMCPToolset() (v BetaManagedAgentsMCPToolset) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsSessionAgentToolUnion) AsCustom() (v BetaManagedAgentsCustomTool) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsSessionAgentToolUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsSessionAgentToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSessionAgentToolUnionConfigs is an implicit subunion of
// [BetaManagedAgentsSessionAgentToolUnion].
// BetaManagedAgentsSessionAgentToolUnionConfigs provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsSessionAgentToolUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfBetaManagedAgentsAgentToolConfigArray
// OfBetaManagedAgentsMCPToolConfigArray]
type BetaManagedAgentsSessionAgentToolUnionConfigs struct {
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

func (r *BetaManagedAgentsSessionAgentToolUnionConfigs) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSessionAgentToolUnionDefaultConfig is an implicit subunion of
// [BetaManagedAgentsSessionAgentToolUnion].
// BetaManagedAgentsSessionAgentToolUnionDefaultConfig provides convenient access
// to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsSessionAgentToolUnion].
type BetaManagedAgentsSessionAgentToolUnionDefaultConfig struct {
	Enabled bool `json:"enabled"`
	// This field is a union of
	// [BetaManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion],
	// [BetaManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion]
	PermissionPolicy BetaManagedAgentsSessionAgentToolUnionDefaultConfigPermissionPolicy `json:"permission_policy"`
	JSON             struct {
		Enabled          respjson.Field
		PermissionPolicy respjson.Field
		raw              string
	} `json:"-"`
}

func (r *BetaManagedAgentsSessionAgentToolUnionDefaultConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsSessionAgentToolUnionDefaultConfigPermissionPolicy is an
// implicit subunion of [BetaManagedAgentsSessionAgentToolUnion].
// BetaManagedAgentsSessionAgentToolUnionDefaultConfigPermissionPolicy provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [BetaManagedAgentsSessionAgentToolUnion].
type BetaManagedAgentsSessionAgentToolUnionDefaultConfigPermissionPolicy struct {
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

func (r *BetaManagedAgentsSessionAgentToolUnionDefaultConfigPermissionPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSessionAgentType string

const (
	BetaManagedAgentsSessionAgentTypeAgent BetaManagedAgentsSessionAgentType = "agent"
)

// Timing statistics for a session.
type BetaManagedAgentsSessionStats struct {
	// Cumulative time in seconds the session spent in running status. Excludes idle
	// time.
	ActiveSeconds float64 `json:"active_seconds"`
	// Elapsed time since session creation in seconds. For terminated sessions, frozen
	// at the final update.
	DurationSeconds float64 `json:"duration_seconds"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ActiveSeconds   respjson.Field
		DurationSeconds respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionStats) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionStats) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Cumulative token usage for a session across all turns.
type BetaManagedAgentsSessionUsage struct {
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
func (r BetaManagedAgentsSessionUsage) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaSessionNewParams struct {
	// Agent identifier. Accepts the `agent` ID string, which pins the latest version
	// for the session, or an `agent` object with both id and version specified.
	Agent BetaSessionNewParamsAgentUnion `json:"agent,omitzero" api:"required"`
	// ID of the `environment` defining the container configuration for this session.
	EnvironmentID string `json:"environment_id" api:"required"`
	// Human-readable session title.
	Title param.Opt[string] `json:"title,omitzero"`
	// Arbitrary key-value metadata attached to the session. Maximum 16 pairs, keys up
	// to 64 chars, values up to 512 chars.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Resources (e.g. repositories, files) to mount into the session's container.
	Resources []BetaSessionNewParamsResourceUnion `json:"resources,omitzero"`
	// Vault IDs for stored credentials the agent can use during the session.
	VaultIDs []string `json:"vault_ids,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaSessionNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaSessionNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaSessionNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaSessionNewParamsAgentUnion struct {
	OfString                  param.Opt[string]             `json:",omitzero,inline"`
	OfBetaManagedAgentsAgents *BetaManagedAgentsAgentParams `json:",omitzero,inline"`
	paramUnion
}

func (u BetaSessionNewParamsAgentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfBetaManagedAgentsAgents)
}
func (u *BetaSessionNewParamsAgentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaSessionNewParamsAgentUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfBetaManagedAgentsAgents) {
		return u.OfBetaManagedAgentsAgents
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaSessionNewParamsResourceUnion struct {
	OfGitHubRepository *BetaManagedAgentsGitHubRepositoryResourceParams `json:",omitzero,inline"`
	OfFile             *BetaManagedAgentsFileResourceParams             `json:",omitzero,inline"`
	paramUnion
}

func (u BetaSessionNewParamsResourceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfGitHubRepository, u.OfFile)
}
func (u *BetaSessionNewParamsResourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaSessionNewParamsResourceUnion) asAny() any {
	if !param.IsOmitted(u.OfGitHubRepository) {
		return u.OfGitHubRepository
	} else if !param.IsOmitted(u.OfFile) {
		return u.OfFile
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaSessionNewParamsResourceUnion) GetAuthorizationToken() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.AuthorizationToken
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaSessionNewParamsResourceUnion) GetURL() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaSessionNewParamsResourceUnion) GetCheckout() *BetaManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.Checkout
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaSessionNewParamsResourceUnion) GetFileID() *string {
	if vt := u.OfFile; vt != nil {
		return &vt.FileID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaSessionNewParamsResourceUnion) GetType() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFile; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaSessionNewParamsResourceUnion) GetMountPath() *string {
	if vt := u.OfGitHubRepository; vt != nil && vt.MountPath.Valid() {
		return &vt.MountPath.Value
	} else if vt := u.OfFile; vt != nil && vt.MountPath.Valid() {
		return &vt.MountPath.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaSessionNewParamsResourceUnion](
		"type",
		apijson.Discriminator[BetaManagedAgentsGitHubRepositoryResourceParams]("github_repository"),
		apijson.Discriminator[BetaManagedAgentsFileResourceParams]("file"),
	)
}

type BetaSessionGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaSessionUpdateParams struct {
	// Human-readable session title.
	Title param.Opt[string] `json:"title,omitzero"`
	// Metadata patch. Set a key to a string to upsert it, or to null to delete it.
	// Omit the field to preserve.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Vault IDs (`vlt_*`) to attach to the session. Not yet supported; requests
	// setting this field are rejected. Reserved for future use.
	VaultIDs []string `json:"vault_ids,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaSessionUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaSessionUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaSessionUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaSessionListParams struct {
	// Filter sessions created with this agent ID.
	AgentID param.Opt[string] `query:"agent_id,omitzero" json:"-"`
	// Filter by agent version. Only applies when agent_id is also set.
	AgentVersion param.Opt[int64] `query:"agent_version,omitzero" json:"-"`
	// Return sessions created after this time (exclusive).
	CreatedAtGt param.Opt[time.Time] `query:"created_at[gt],omitzero" format:"date-time" json:"-"`
	// Return sessions created at or after this time (inclusive).
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return sessions created before this time (exclusive).
	CreatedAtLt param.Opt[time.Time] `query:"created_at[lt],omitzero" format:"date-time" json:"-"`
	// Return sessions created at or before this time (inclusive).
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// When true, includes archived sessions. Default: false (exclude archived).
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum number of results to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor from a previous response's next_page.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Sort direction for results, ordered by created_at. Defaults to desc (newest
	// first).
	//
	// Any of "asc", "desc".
	Order BetaSessionListParamsOrder `query:"order,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaSessionListParams]'s query parameters as `url.Values`.
func (r BetaSessionListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort direction for results, ordered by created_at. Defaults to desc (newest
// first).
type BetaSessionListParamsOrder string

const (
	BetaSessionListParamsOrderAsc  BetaSessionListParamsOrder = "asc"
	BetaSessionListParamsOrderDesc BetaSessionListParamsOrder = "desc"
)

type BetaSessionDeleteParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaSessionArchiveParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
