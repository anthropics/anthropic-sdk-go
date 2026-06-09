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

// BetaDeploymentRunService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaDeploymentRunService] method instead.
type BetaDeploymentRunService struct {
	Options []option.RequestOption
}

// NewBetaDeploymentRunService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaDeploymentRunService(opts ...option.RequestOption) (r BetaDeploymentRunService) {
	r = BetaDeploymentRunService{}
	r.Options = opts
	return
}

// Get Deployment Run
func (r *BetaDeploymentRunService) Get(ctx context.Context, deploymentRunID string, query BetaDeploymentRunGetParams, opts ...option.RequestOption) (res *BetaManagedAgentsDeploymentRun, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if deploymentRunID == "" {
		err = errors.New("missing required deployment_run_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/deployment_runs/%s?beta=true", deploymentRunID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List Deployment Runs
func (r *BetaDeploymentRunService) List(ctx context.Context, params BetaDeploymentRunListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaManagedAgentsDeploymentRun], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/deployment_runs?beta=true"
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

// List Deployment Runs
func (r *BetaDeploymentRunService) ListAutoPaging(ctx context.Context, params BetaDeploymentRunListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsDeploymentRun] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// The deployment's agent was archived.
type BetaManagedAgentsAgentArchivedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "agent_archived_error".
	Type BetaManagedAgentsAgentArchivedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAgentArchivedRunError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAgentArchivedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsAgentArchivedRunErrorType string

const (
	BetaManagedAgentsAgentArchivedRunErrorTypeAgentArchivedError BetaManagedAgentsAgentArchivedRunErrorType = "agent_archived_error"
)

// A persistent, append-only record of a single deployment execution. Records
// session creation success or failure — no session lifecycle tracking.
type BetaManagedAgentsDeploymentRun struct {
	// Unique identifier for this run (`drun_...`).
	ID string `json:"id" api:"required"`
	// A resolved agent reference with a concrete version.
	Agent BetaManagedAgentsAgentReference `json:"agent" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// ID of the deployment that produced this run.
	DeploymentID string `json:"deployment_id" api:"required"`
	// Why the run failed to create a session. The type identifies the failure; message
	// is human-readable detail.
	Error BetaManagedAgentsDeploymentRunErrorUnion `json:"error" api:"required"`
	// Populated on success. Null on creation failure. Exactly one of session_id or
	// error is non-null.
	SessionID string `json:"session_id" api:"required"`
	// Describes what triggered a deployment run, with trigger-specific metadata.
	TriggerContext BetaManagedAgentsTriggerContextUnion `json:"trigger_context" api:"required"`
	// Any of "deployment_run".
	Type BetaManagedAgentsDeploymentRunType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		Agent          respjson.Field
		CreatedAt      respjson.Field
		DeploymentID   respjson.Field
		Error          respjson.Field
		SessionID      respjson.Field
		TriggerContext respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsDeploymentRun) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsDeploymentRun) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaManagedAgentsDeploymentRunErrorUnion contains all possible properties and
// values from [BetaManagedAgentsEnvironmentArchivedRunError],
// [BetaManagedAgentsAgentArchivedRunError],
// [BetaManagedAgentsEnvironmentNotFoundRunError],
// [BetaManagedAgentsVaultNotFoundRunError],
// [BetaManagedAgentsVaultArchivedRunError],
// [BetaManagedAgentsFileNotFoundRunError],
// [BetaManagedAgentsMemoryStoreArchivedRunError],
// [BetaManagedAgentsSkillNotFoundRunError],
// [BetaManagedAgentsSessionResourceNotFoundRunError],
// [BetaManagedAgentsWorkspaceArchivedRunError],
// [BetaManagedAgentsOrganizationDisabledRunError],
// [BetaManagedAgentsSessionRateLimitedRunError],
// [BetaManagedAgentsSessionCreationRejectedRunError],
// [BetaManagedAgentsUnknownRunError],
// [BetaManagedAgentsSelfHostedResourcesUnsupportedRunError],
// [BetaManagedAgentsMCPEgressBlockedRunError].
//
// Use the [BetaManagedAgentsDeploymentRunErrorUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsDeploymentRunErrorUnion struct {
	Message string `json:"message"`
	// Any of "environment_archived_error", "agent_archived_error",
	// "environment_not_found_error", "vault_not_found_error", "vault_archived_error",
	// "file_not_found_error", "memory_store_archived_error", "skill_not_found_error",
	// "session_resource_not_found_error", "workspace_archived_error",
	// "organization_disabled_error", "session_rate_limited_error",
	// "session_creation_rejected_error", "unknown_error",
	// "self_hosted_resources_unsupported_error", "mcp_egress_blocked_error".
	Type string `json:"type"`
	JSON struct {
		Message respjson.Field
		Type    respjson.Field
		raw     string
	} `json:"-"`
}

// anyBetaManagedAgentsDeploymentRunError is implemented by each variant of
// [BetaManagedAgentsDeploymentRunErrorUnion] to add type safety for the return
// type of [BetaManagedAgentsDeploymentRunErrorUnion.AsAny]
type anyBetaManagedAgentsDeploymentRunError interface {
	implBetaManagedAgentsDeploymentRunErrorUnion()
}

func (BetaManagedAgentsEnvironmentArchivedRunError) implBetaManagedAgentsDeploymentRunErrorUnion() {}
func (BetaManagedAgentsAgentArchivedRunError) implBetaManagedAgentsDeploymentRunErrorUnion()       {}
func (BetaManagedAgentsEnvironmentNotFoundRunError) implBetaManagedAgentsDeploymentRunErrorUnion() {}
func (BetaManagedAgentsVaultNotFoundRunError) implBetaManagedAgentsDeploymentRunErrorUnion()       {}
func (BetaManagedAgentsVaultArchivedRunError) implBetaManagedAgentsDeploymentRunErrorUnion()       {}
func (BetaManagedAgentsFileNotFoundRunError) implBetaManagedAgentsDeploymentRunErrorUnion()        {}
func (BetaManagedAgentsMemoryStoreArchivedRunError) implBetaManagedAgentsDeploymentRunErrorUnion() {}
func (BetaManagedAgentsSkillNotFoundRunError) implBetaManagedAgentsDeploymentRunErrorUnion()       {}
func (BetaManagedAgentsSessionResourceNotFoundRunError) implBetaManagedAgentsDeploymentRunErrorUnion() {
}
func (BetaManagedAgentsWorkspaceArchivedRunError) implBetaManagedAgentsDeploymentRunErrorUnion()    {}
func (BetaManagedAgentsOrganizationDisabledRunError) implBetaManagedAgentsDeploymentRunErrorUnion() {}
func (BetaManagedAgentsSessionRateLimitedRunError) implBetaManagedAgentsDeploymentRunErrorUnion()   {}
func (BetaManagedAgentsSessionCreationRejectedRunError) implBetaManagedAgentsDeploymentRunErrorUnion() {
}
func (BetaManagedAgentsUnknownRunError) implBetaManagedAgentsDeploymentRunErrorUnion() {}
func (BetaManagedAgentsSelfHostedResourcesUnsupportedRunError) implBetaManagedAgentsDeploymentRunErrorUnion() {
}
func (BetaManagedAgentsMCPEgressBlockedRunError) implBetaManagedAgentsDeploymentRunErrorUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsDeploymentRunErrorUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsEnvironmentArchivedRunError:
//	case anthropic.BetaManagedAgentsAgentArchivedRunError:
//	case anthropic.BetaManagedAgentsEnvironmentNotFoundRunError:
//	case anthropic.BetaManagedAgentsVaultNotFoundRunError:
//	case anthropic.BetaManagedAgentsVaultArchivedRunError:
//	case anthropic.BetaManagedAgentsFileNotFoundRunError:
//	case anthropic.BetaManagedAgentsMemoryStoreArchivedRunError:
//	case anthropic.BetaManagedAgentsSkillNotFoundRunError:
//	case anthropic.BetaManagedAgentsSessionResourceNotFoundRunError:
//	case anthropic.BetaManagedAgentsWorkspaceArchivedRunError:
//	case anthropic.BetaManagedAgentsOrganizationDisabledRunError:
//	case anthropic.BetaManagedAgentsSessionRateLimitedRunError:
//	case anthropic.BetaManagedAgentsSessionCreationRejectedRunError:
//	case anthropic.BetaManagedAgentsUnknownRunError:
//	case anthropic.BetaManagedAgentsSelfHostedResourcesUnsupportedRunError:
//	case anthropic.BetaManagedAgentsMCPEgressBlockedRunError:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsDeploymentRunErrorUnion) AsAny() anyBetaManagedAgentsDeploymentRunError {
	switch u.Type {
	case "environment_archived_error":
		return u.AsEnvironmentArchivedError()
	case "agent_archived_error":
		return u.AsAgentArchivedError()
	case "environment_not_found_error":
		return u.AsEnvironmentNotFoundError()
	case "vault_not_found_error":
		return u.AsVaultNotFoundError()
	case "vault_archived_error":
		return u.AsVaultArchivedError()
	case "file_not_found_error":
		return u.AsFileNotFoundError()
	case "memory_store_archived_error":
		return u.AsMemoryStoreArchivedError()
	case "skill_not_found_error":
		return u.AsSkillNotFoundError()
	case "session_resource_not_found_error":
		return u.AsSessionResourceNotFoundError()
	case "workspace_archived_error":
		return u.AsWorkspaceArchivedError()
	case "organization_disabled_error":
		return u.AsOrganizationDisabledError()
	case "session_rate_limited_error":
		return u.AsSessionRateLimitedError()
	case "session_creation_rejected_error":
		return u.AsSessionCreationRejectedError()
	case "unknown_error":
		return u.AsUnknownError()
	case "self_hosted_resources_unsupported_error":
		return u.AsSelfHostedResourcesUnsupportedError()
	case "mcp_egress_blocked_error":
		return u.AsMCPEgressBlockedError()
	}
	return nil
}

func (u BetaManagedAgentsDeploymentRunErrorUnion) AsEnvironmentArchivedError() (v BetaManagedAgentsEnvironmentArchivedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentRunErrorUnion) AsAgentArchivedError() (v BetaManagedAgentsAgentArchivedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentRunErrorUnion) AsEnvironmentNotFoundError() (v BetaManagedAgentsEnvironmentNotFoundRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentRunErrorUnion) AsVaultNotFoundError() (v BetaManagedAgentsVaultNotFoundRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentRunErrorUnion) AsVaultArchivedError() (v BetaManagedAgentsVaultArchivedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentRunErrorUnion) AsFileNotFoundError() (v BetaManagedAgentsFileNotFoundRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentRunErrorUnion) AsMemoryStoreArchivedError() (v BetaManagedAgentsMemoryStoreArchivedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentRunErrorUnion) AsSkillNotFoundError() (v BetaManagedAgentsSkillNotFoundRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentRunErrorUnion) AsSessionResourceNotFoundError() (v BetaManagedAgentsSessionResourceNotFoundRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentRunErrorUnion) AsWorkspaceArchivedError() (v BetaManagedAgentsWorkspaceArchivedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentRunErrorUnion) AsOrganizationDisabledError() (v BetaManagedAgentsOrganizationDisabledRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentRunErrorUnion) AsSessionRateLimitedError() (v BetaManagedAgentsSessionRateLimitedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentRunErrorUnion) AsSessionCreationRejectedError() (v BetaManagedAgentsSessionCreationRejectedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentRunErrorUnion) AsUnknownError() (v BetaManagedAgentsUnknownRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentRunErrorUnion) AsSelfHostedResourcesUnsupportedError() (v BetaManagedAgentsSelfHostedResourcesUnsupportedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsDeploymentRunErrorUnion) AsMCPEgressBlockedError() (v BetaManagedAgentsMCPEgressBlockedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsDeploymentRunErrorUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsDeploymentRunErrorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsDeploymentRunType string

const (
	BetaManagedAgentsDeploymentRunTypeDeploymentRun BetaManagedAgentsDeploymentRunType = "deployment_run"
)

// The deployment's environment was archived.
type BetaManagedAgentsEnvironmentArchivedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "environment_archived_error".
	Type BetaManagedAgentsEnvironmentArchivedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsEnvironmentArchivedRunError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsEnvironmentArchivedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsEnvironmentArchivedRunErrorType string

const (
	BetaManagedAgentsEnvironmentArchivedRunErrorTypeEnvironmentArchivedError BetaManagedAgentsEnvironmentArchivedRunErrorType = "environment_archived_error"
)

// The deployment's environment no longer exists.
type BetaManagedAgentsEnvironmentNotFoundRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "environment_not_found_error".
	Type BetaManagedAgentsEnvironmentNotFoundRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsEnvironmentNotFoundRunError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsEnvironmentNotFoundRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsEnvironmentNotFoundRunErrorType string

const (
	BetaManagedAgentsEnvironmentNotFoundRunErrorTypeEnvironmentNotFoundError BetaManagedAgentsEnvironmentNotFoundRunErrorType = "environment_not_found_error"
)

// A file resource referenced by the deployment no longer exists.
type BetaManagedAgentsFileNotFoundRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "file_not_found_error".
	Type BetaManagedAgentsFileNotFoundRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsFileNotFoundRunError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsFileNotFoundRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsFileNotFoundRunErrorType string

const (
	BetaManagedAgentsFileNotFoundRunErrorTypeFileNotFoundError BetaManagedAgentsFileNotFoundRunErrorType = "file_not_found_error"
)

// The run was started manually by creating a session directly against the
// deployment.
type BetaManagedAgentsManualTriggerContext struct {
	// Any of "manual".
	Type BetaManagedAgentsManualTriggerContextType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsManualTriggerContext) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsManualTriggerContext) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsManualTriggerContextType string

const (
	BetaManagedAgentsManualTriggerContextTypeManual BetaManagedAgentsManualTriggerContextType = "manual"
)

// An MCP server host used by the deployment's agent is blocked by the
// environment's network policy.
type BetaManagedAgentsMCPEgressBlockedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "mcp_egress_blocked_error".
	Type BetaManagedAgentsMCPEgressBlockedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMCPEgressBlockedRunError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsMCPEgressBlockedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMCPEgressBlockedRunErrorType string

const (
	BetaManagedAgentsMCPEgressBlockedRunErrorTypeMCPEgressBlockedError BetaManagedAgentsMCPEgressBlockedRunErrorType = "mcp_egress_blocked_error"
)

// A memory store referenced by the deployment is archived.
type BetaManagedAgentsMemoryStoreArchivedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "memory_store_archived_error".
	Type BetaManagedAgentsMemoryStoreArchivedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMemoryStoreArchivedRunError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsMemoryStoreArchivedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMemoryStoreArchivedRunErrorType string

const (
	BetaManagedAgentsMemoryStoreArchivedRunErrorTypeMemoryStoreArchivedError BetaManagedAgentsMemoryStoreArchivedRunErrorType = "memory_store_archived_error"
)

// The deployment's organization is disabled.
type BetaManagedAgentsOrganizationDisabledRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "organization_disabled_error".
	Type BetaManagedAgentsOrganizationDisabledRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsOrganizationDisabledRunError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsOrganizationDisabledRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsOrganizationDisabledRunErrorType string

const (
	BetaManagedAgentsOrganizationDisabledRunErrorTypeOrganizationDisabledError BetaManagedAgentsOrganizationDisabledRunErrorType = "organization_disabled_error"
)

// The run was fired by the deployment's cron schedule.
type BetaManagedAgentsScheduleTriggerContext struct {
	// A timestamp in RFC 3339 format
	ScheduledAt time.Time `json:"scheduled_at" api:"required" format:"date-time"`
	// Any of "schedule".
	Type BetaManagedAgentsScheduleTriggerContextType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ScheduledAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsScheduleTriggerContext) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsScheduleTriggerContext) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsScheduleTriggerContextType string

const (
	BetaManagedAgentsScheduleTriggerContextTypeSchedule BetaManagedAgentsScheduleTriggerContextType = "schedule"
)

// The deployment configures resources, but its environment is self-hosted and
// cannot mount them.
type BetaManagedAgentsSelfHostedResourcesUnsupportedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "self_hosted_resources_unsupported_error".
	Type BetaManagedAgentsSelfHostedResourcesUnsupportedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSelfHostedResourcesUnsupportedRunError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSelfHostedResourcesUnsupportedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSelfHostedResourcesUnsupportedRunErrorType string

const (
	BetaManagedAgentsSelfHostedResourcesUnsupportedRunErrorTypeSelfHostedResourcesUnsupportedError BetaManagedAgentsSelfHostedResourcesUnsupportedRunErrorType = "self_hosted_resources_unsupported_error"
)

// The session create request was rejected with a non-retryable validation error.
type BetaManagedAgentsSessionCreationRejectedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "session_creation_rejected_error".
	Type BetaManagedAgentsSessionCreationRejectedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionCreationRejectedRunError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionCreationRejectedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSessionCreationRejectedRunErrorType string

const (
	BetaManagedAgentsSessionCreationRejectedRunErrorTypeSessionCreationRejectedError BetaManagedAgentsSessionCreationRejectedRunErrorType = "session_creation_rejected_error"
)

// Session creation was rejected due to rate limiting. The schedule keeps firing;
// subsequent runs may succeed.
type BetaManagedAgentsSessionRateLimitedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "session_rate_limited_error".
	Type BetaManagedAgentsSessionRateLimitedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionRateLimitedRunError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionRateLimitedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSessionRateLimitedRunErrorType string

const (
	BetaManagedAgentsSessionRateLimitedRunErrorTypeSessionRateLimitedError BetaManagedAgentsSessionRateLimitedRunErrorType = "session_rate_limited_error"
)

// A referenced resource no longer exists and its kind was not reported.
type BetaManagedAgentsSessionResourceNotFoundRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "session_resource_not_found_error".
	Type BetaManagedAgentsSessionResourceNotFoundRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionResourceNotFoundRunError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionResourceNotFoundRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSessionResourceNotFoundRunErrorType string

const (
	BetaManagedAgentsSessionResourceNotFoundRunErrorTypeSessionResourceNotFoundError BetaManagedAgentsSessionResourceNotFoundRunErrorType = "session_resource_not_found_error"
)

// A skill referenced by the deployment's agent no longer exists.
type BetaManagedAgentsSkillNotFoundRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "skill_not_found_error".
	Type BetaManagedAgentsSkillNotFoundRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSkillNotFoundRunError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSkillNotFoundRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsSkillNotFoundRunErrorType string

const (
	BetaManagedAgentsSkillNotFoundRunErrorTypeSkillNotFoundError BetaManagedAgentsSkillNotFoundRunErrorType = "skill_not_found_error"
)

// BetaManagedAgentsTriggerContextUnion contains all possible properties and values
// from [BetaManagedAgentsScheduleTriggerContext],
// [BetaManagedAgentsManualTriggerContext].
//
// Use the [BetaManagedAgentsTriggerContextUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsTriggerContextUnion struct {
	// This field is from variant [BetaManagedAgentsScheduleTriggerContext].
	ScheduledAt time.Time `json:"scheduled_at"`
	// Any of "schedule", "manual".
	Type string `json:"type"`
	JSON struct {
		ScheduledAt respjson.Field
		Type        respjson.Field
		raw         string
	} `json:"-"`
}

// anyBetaManagedAgentsTriggerContext is implemented by each variant of
// [BetaManagedAgentsTriggerContextUnion] to add type safety for the return type of
// [BetaManagedAgentsTriggerContextUnion.AsAny]
type anyBetaManagedAgentsTriggerContext interface {
	implBetaManagedAgentsTriggerContextUnion()
}

func (BetaManagedAgentsScheduleTriggerContext) implBetaManagedAgentsTriggerContextUnion() {}
func (BetaManagedAgentsManualTriggerContext) implBetaManagedAgentsTriggerContextUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsTriggerContextUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsScheduleTriggerContext:
//	case anthropic.BetaManagedAgentsManualTriggerContext:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsTriggerContextUnion) AsAny() anyBetaManagedAgentsTriggerContext {
	switch u.Type {
	case "schedule":
		return u.AsSchedule()
	case "manual":
		return u.AsManual()
	}
	return nil
}

func (u BetaManagedAgentsTriggerContextUnion) AsSchedule() (v BetaManagedAgentsScheduleTriggerContext) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsTriggerContextUnion) AsManual() (v BetaManagedAgentsManualTriggerContext) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsTriggerContextUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsTriggerContextUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// What triggered a deployment run.
type BetaManagedAgentsTriggerType string

const (
	BetaManagedAgentsTriggerTypeSchedule BetaManagedAgentsTriggerType = "schedule"
	BetaManagedAgentsTriggerTypeManual   BetaManagedAgentsTriggerType = "manual"
)

// An unknown or unexpected error caused the run to fail. A fallback variant;
// clients that do not recognize a new error type can match on message alone.
type BetaManagedAgentsUnknownRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "unknown_error".
	Type BetaManagedAgentsUnknownRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsUnknownRunError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsUnknownRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsUnknownRunErrorType string

const (
	BetaManagedAgentsUnknownRunErrorTypeUnknownError BetaManagedAgentsUnknownRunErrorType = "unknown_error"
)

// A vault referenced by the deployment is archived.
type BetaManagedAgentsVaultArchivedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "vault_archived_error".
	Type BetaManagedAgentsVaultArchivedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsVaultArchivedRunError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsVaultArchivedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsVaultArchivedRunErrorType string

const (
	BetaManagedAgentsVaultArchivedRunErrorTypeVaultArchivedError BetaManagedAgentsVaultArchivedRunErrorType = "vault_archived_error"
)

// A vault referenced by the deployment no longer exists.
type BetaManagedAgentsVaultNotFoundRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "vault_not_found_error".
	Type BetaManagedAgentsVaultNotFoundRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsVaultNotFoundRunError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsVaultNotFoundRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsVaultNotFoundRunErrorType string

const (
	BetaManagedAgentsVaultNotFoundRunErrorTypeVaultNotFoundError BetaManagedAgentsVaultNotFoundRunErrorType = "vault_not_found_error"
)

// The deployment's workspace was archived.
type BetaManagedAgentsWorkspaceArchivedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "workspace_archived_error".
	Type BetaManagedAgentsWorkspaceArchivedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsWorkspaceArchivedRunError) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsWorkspaceArchivedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsWorkspaceArchivedRunErrorType string

const (
	BetaManagedAgentsWorkspaceArchivedRunErrorTypeWorkspaceArchivedError BetaManagedAgentsWorkspaceArchivedRunErrorType = "workspace_archived_error"
)

type BetaDeploymentRunGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaDeploymentRunListParams struct {
	// Return runs created strictly after this time (exclusive).
	CreatedAtGt param.Opt[time.Time] `query:"created_at[gt],omitzero" format:"date-time" json:"-"`
	// Return runs created at or after this time (inclusive).
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return runs created strictly before this time (exclusive).
	CreatedAtLt param.Opt[time.Time] `query:"created_at[lt],omitzero" format:"date-time" json:"-"`
	// Return runs created at or before this time (inclusive).
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// Filter to a specific deployment. Omit to list across all deployments in the
	// workspace. Filtering by a non-existent deployment_id returns 200 with empty
	// data.
	DeploymentID param.Opt[string] `query:"deployment_id,omitzero" json:"-"`
	// Filter: true for runs with non-null error, false for runs with non-null
	// session_id. Omit for all.
	HasError param.Opt[bool] `query:"has_error,omitzero" json:"-"`
	// Maximum results per page. Default 20, maximum 1000.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor. Pass next_page from the previous response. Invalid or
	// expired cursors return 400.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Filter runs by what triggered them. Omit to return all runs.
	//
	// Any of "schedule", "manual".
	TriggerType BetaManagedAgentsTriggerType `query:"trigger_type,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaDeploymentRunListParams]'s query parameters as
// `url.Values`.
func (r BetaDeploymentRunListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
