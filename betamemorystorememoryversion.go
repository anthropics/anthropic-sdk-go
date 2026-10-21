package anthropic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// BetaMemoryStoreMemoryVersionService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaMemoryStoreMemoryVersionService] method instead.
type BetaMemoryStoreMemoryVersionService struct {
	Options []option.RequestOption
}

// NewBetaMemoryStoreMemoryVersionService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaMemoryStoreMemoryVersionService(opts ...option.RequestOption) (r BetaMemoryStoreMemoryVersionService) {
	r = BetaMemoryStoreMemoryVersionService{}
	r.Options = opts
	return
}

// Retrieve a Memory Version
func (r *BetaMemoryStoreMemoryVersionService) Get(ctx context.Context, memoryVersionID string, params BetaMemoryStoreMemoryVersionGetParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemoryVersion, err error) {
	if len(params.Betas) > 0 {
		headerValues := make([]string, len(params.Betas))
		for i, v := range params.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	if !param.IsOmitted(params.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", params.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "agent-memory-2026-07-22")}, opts...)
	if params.MemoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	if memoryVersionID == "" {
		err = errors.New("missing required memory_version_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s/memory_versions/%s?beta=true", url.PathEscape(params.MemoryStoreID), url.PathEscape(memoryVersionID))
	opts = append(opts, requestconfig.WithSpan("anthropic.memory_stores.memory_versions.retrieve", "/v1/memory_stores/{memory_store_id}/memory_versions/{memory_version_id}", "view"), requestconfig.WithSpanAttribute("gen_ai.memory.store.id", params.MemoryStoreID), requestconfig.WithSpanAttribute("anthropic.memory_store.memory_version.id", memoryVersionID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// List Memory Versions
func (r *BetaMemoryStoreMemoryVersionService) List(ctx context.Context, memoryStoreID string, params BetaMemoryStoreMemoryVersionListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaManagedAgentsMemoryVersion], err error) {
	var raw *http.Response
	if len(params.Betas) > 0 {
		headerValues := make([]string, len(params.Betas))
		for i, v := range params.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	if !param.IsOmitted(params.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", params.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "agent-memory-2026-07-22"), option.WithResponseInto(&raw)}, opts...)
	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s/memory_versions?beta=true", url.PathEscape(memoryStoreID))
	opts = append(opts, requestconfig.WithSpan("anthropic.memory_stores.memory_versions.list", "/v1/memory_stores/{memory_store_id}/memory_versions", "api_key_id", "created_at[gte]", "created_at[lte]", "limit", "memory_id", "operation", "page", "service_account_id", "session_id", "view"), requestconfig.WithSpanAttribute("gen_ai.memory.store.id", memoryStoreID))
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

// List Memory Versions
func (r *BetaMemoryStoreMemoryVersionService) ListAutoPaging(ctx context.Context, memoryStoreID string, params BetaMemoryStoreMemoryVersionListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsMemoryVersion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, memoryStoreID, params, opts...))
}

// Redact a Memory Version
func (r *BetaMemoryStoreMemoryVersionService) Redact(ctx context.Context, memoryVersionID string, params BetaMemoryStoreMemoryVersionRedactParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemoryVersion, err error) {
	if len(params.Betas) > 0 {
		headerValues := make([]string, len(params.Betas))
		for i, v := range params.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	if !param.IsOmitted(params.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", params.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "agent-memory-2026-07-22")}, opts...)
	if params.MemoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	if memoryVersionID == "" {
		err = errors.New("missing required memory_version_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s/memory_versions/%s/redact?beta=true", url.PathEscape(params.MemoryStoreID), url.PathEscape(memoryVersionID))
	opts = append(opts, requestconfig.WithSpan("anthropic.memory_stores.memory_versions.redact", "/v1/memory_stores/{memory_store_id}/memory_versions/{memory_version_id}/redact"), requestconfig.WithSpanAttribute("gen_ai.memory.store.id", params.MemoryStoreID), requestconfig.WithSpanAttribute("anthropic.memory_store.memory_version.id", memoryVersionID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// BetaManagedAgentsActorUnion contains all possible properties and values from
// [BetaManagedAgentsSessionActor], [BetaManagedAgentsAPIActor],
// [BetaManagedAgentsUserActor], [BetaManagedAgentsServiceAccountActor].
//
// Use the [BetaManagedAgentsActorUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsActorUnion struct {
	// This field is from variant [BetaManagedAgentsSessionActor].
	SessionID string `json:"session_id"`
	// Any of "session_actor", "api_actor", "user_actor", "service_account_actor".
	Type string `json:"type"`
	// This field is from variant [BetaManagedAgentsAPIActor].
	APIKeyID string `json:"api_key_id"`
	// This field is from variant [BetaManagedAgentsUserActor].
	UserID string `json:"user_id"`
	// This field is from variant [BetaManagedAgentsServiceAccountActor].
	ServiceAccountID string `json:"service_account_id"`
	JSON             struct {
		SessionID        respjson.Field
		Type             respjson.Field
		APIKeyID         respjson.Field
		UserID           respjson.Field
		ServiceAccountID respjson.Field
		raw              string
	} `json:"-"`
}

// anyBetaManagedAgentsActor is implemented by each variant of
// [BetaManagedAgentsActorUnion] to add type safety for the return type of
// [BetaManagedAgentsActorUnion.AsAny]
type anyBetaManagedAgentsActor interface {
	implBetaManagedAgentsActorUnion()
}

func (BetaManagedAgentsSessionActor) implBetaManagedAgentsActorUnion()        {}
func (BetaManagedAgentsAPIActor) implBetaManagedAgentsActorUnion()            {}
func (BetaManagedAgentsUserActor) implBetaManagedAgentsActorUnion()           {}
func (BetaManagedAgentsServiceAccountActor) implBetaManagedAgentsActorUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsActorUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsSessionActor:
//	case anthropic.BetaManagedAgentsAPIActor:
//	case anthropic.BetaManagedAgentsUserActor:
//	case anthropic.BetaManagedAgentsServiceAccountActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsActorUnion) AsAny() anyBetaManagedAgentsActor {
	switch u.Type {
	case "session_actor":
		return u.AsSessionActor()
	case "api_actor":
		return u.AsAPIActor()
	case "user_actor":
		return u.AsUserActor()
	case "service_account_actor":
		return u.AsServiceAccountActor()
	}
	return nil
}

func (u BetaManagedAgentsActorUnion) AsSessionActor() (v BetaManagedAgentsSessionActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsActorUnion) AsAPIActor() (v BetaManagedAgentsAPIActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsActorUnion) AsUserActor() (v BetaManagedAgentsUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsActorUnion) AsServiceAccountActor() (v BetaManagedAgentsServiceAccountActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsActorUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsActorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A direct caller of the public API, identified by the API key that authenticated
// the request.
type BetaManagedAgentsAPIActor struct {
	// ID of the API key (an `apikey_...` value). This identifies the key, not the
	// secret.
	APIKeyID string            `json:"api_key_id" api:"required"`
	Type     constant.APIActor `json:"type" default:"api_actor"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		APIKeyID    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsAPIActor) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsAPIActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// One entry in a memory's history, recording a single change and who made it.
//
// Every change to a memory adds one version. Only redaction can change a version
// after that. Versions belong to the memory store and stay after the memory is
// deleted. Each version is retained for at least the version retention period
// after it was written, unless the store itself is deleted.
//
// Retrieving a redacted version succeeds, with `content`, `path`,
// `content_size_bytes`, and `content_sha256` set to `null`. Check `redacted_at` to
// find out whether a version was redacted.
//
// See the
// [memory guide](https://platform.claude.com/docs/en/managed-agents/memory#audit-memory-changes)
// to learn more about memory versions.
type BetaManagedAgentsMemoryVersion struct {
	// Unique identifier for this version (a `memver_...` value).
	ID string `json:"id" api:"required"`
	// When this version was written, in RFC 3339 format.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// The ID of the memory that this version belongs to (`mem_...`).
	//
	// The ID stays valid after the memory is deleted. Pass it as `memory_id` to
	// [List Memory Versions](https://platform.claude.com/docs/en/api/beta/memory_stores/memory_versions/list)
	// to get the memory's retained versions, including its `deleted` version while the
	// lineage is retained.
	MemoryID string `json:"memory_id" api:"required"`
	// ID of the memory store this version belongs to (a `memstore_...` value).
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// The kind of mutation this version records: `created`, `modified`, or `deleted`.
	//
	// Any of "created", "modified", "deleted".
	Operation BetaManagedAgentsMemoryVersionOperation `json:"operation" api:"required"`
	Type      constant.MemoryVersion                  `json:"type" default:"memory_version"`
	// The memory's UTF-8 text content as of this version. `null` when `view=basic`,
	// when `operation` is `deleted`, or when `redacted_at` is set.
	Content string `json:"content" api:"nullable"`
	// Lowercase hex SHA-256 digest of `content` as of this version (64 characters).
	// `null` when `redacted_at` is set or `operation` is `deleted`. Populated
	// regardless of `view` otherwise.
	ContentSha256 string `json:"content_sha256" api:"nullable"`
	// Size of `content` in bytes as of this version. `null` when `redacted_at` is set
	// or `operation` is `deleted`. Populated regardless of `view` otherwise.
	ContentSizeBytes int64 `json:"content_size_bytes" api:"nullable"`
	// Who made this change to the memory, or `null` if the author wasn't recorded.
	//
	// Redacting the version doesn't change this field. A `session_actor` here is an
	// agent that wrote through the store's mounted filesystem at `/mnt/memory/`. The
	// version doesn't record the API key that created that session. So the author in
	// `created_by` is who made the change, not who is ultimately responsible for it.
	//
	// You can
	// [retrieve the session](https://platform.claude.com/docs/en/api/beta/sessions/retrieve)
	// to learn more about where the change came from.
	CreatedBy BetaManagedAgentsActorUnion `json:"created_by"`
	// The memory's path at the time of this write. `null` if and only if `redacted_at`
	// is set.
	Path string `json:"path" api:"nullable"`
	// When the version was redacted, in RFC 3339, or `null` if it hasn't been.
	//
	// When set, `content`, `path`, `content_size_bytes`, and `content_sha256` are
	// `null`.
	//
	// See the
	// [memory guide](https://platform.claude.com/docs/en/managed-agents/memory#redact-a-version)
	// to learn more about redaction.
	RedactedAt time.Time `json:"redacted_at" api:"nullable" format:"date-time"`
	// Who redacted this version, or `null` if it has not been redacted. In practice
	// always an `api_actor`, `user_actor`, or `service_account_actor` (agents do not
	// have a redact capability).
	RedactedBy BetaManagedAgentsActorUnion `json:"redacted_by"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		CreatedAt        respjson.Field
		MemoryID         respjson.Field
		MemoryStoreID    respjson.Field
		Operation        respjson.Field
		Type             respjson.Field
		Content          respjson.Field
		ContentSha256    respjson.Field
		ContentSizeBytes respjson.Field
		CreatedBy        respjson.Field
		Path             respjson.Field
		RedactedAt       respjson.Field
		RedactedBy       respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMemoryVersion) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsMemoryVersion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The kind of change that a memory version records.
//
// Every change to a memory adds one version. If you update a memory and nothing
// changes, the memory doesn't get a new version.
type BetaManagedAgentsMemoryVersionOperation string

const (
	// The memory was created. The first version in any memory's lineage.
	BetaManagedAgentsMemoryVersionOperationCreated BetaManagedAgentsMemoryVersionOperation = "created"
	// The memory's `content`, `path`, or both were changed via update. Writes the
	// agent makes through the filesystem mount also appear as `modified`.
	BetaManagedAgentsMemoryVersionOperationModified BetaManagedAgentsMemoryVersionOperation = "modified"
	// The memory was deleted. The `content`, `content_size_bytes`, and
	// `content_sha256` fields are `null` on this version. The preceding version, while
	// it is retained, records the deleted content's size and hash.
	BetaManagedAgentsMemoryVersionOperationDeleted BetaManagedAgentsMemoryVersionOperation = "deleted"
)

// A workload authenticated as a service account, for example via Workload Identity
// Federation.
type BetaManagedAgentsServiceAccountActor struct {
	// ID of the service account (a `svac_...` value).
	ServiceAccountID string                       `json:"service_account_id" api:"required"`
	Type             constant.ServiceAccountActor `json:"type" default:"service_account_actor"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ServiceAccountID respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsServiceAccountActor) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsServiceAccountActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An agent acting during a session, for example through the session's mounted
// filesystem. It names the session itself, not the user or API key that started
// the session.
type BetaManagedAgentsSessionActor struct {
	// ID of the session (a `sesn_...` value). Look up the session via
	// [Retrieve a session](/en/api/beta/sessions/retrieve) for further provenance.
	SessionID string                `json:"session_id" api:"required"`
	Type      constant.SessionActor `json:"type" default:"session_actor"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SessionID   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsSessionActor) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsSessionActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A human user, for example acting through the Anthropic Console.
type BetaManagedAgentsUserActor struct {
	Type constant.UserActor `json:"type" default:"user_actor"`
	// ID of the user (a `user_...` value).
	UserID string `json:"user_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		UserID      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsUserActor) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsUserActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaMemoryStoreMemoryVersionGetParams struct {
	// The ID of the memory store that holds the version (`memstore_...`).
	MemoryStoreID string `path:"memory_store_id" api:"required" json:"-"`
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// Whether a `memory` or `memory_version` in the response includes its `content`.
	// Defaults to `full` on retrieve operations and to `basic` on list, create, and
	// update operations.
	//
	// With `full`, list operations return at most 20 items per page.
	//
	// Any of "basic", "full".
	View BetaManagedAgentsMemoryView `query:"view,omitzero" json:"-"`
	// This endpoint is in beta: requests must send `agent-memory-2026-07-22` or
	// `managed-agents-2026-04-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaMemoryStoreMemoryVersionGetParams]'s query parameters
// as `url.Values`.
func (r BetaMemoryStoreMemoryVersionGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaMemoryStoreMemoryVersionListParams struct {
	// Return only versions written with the API key that has this ID.
	APIKeyID param.Opt[string] `query:"api_key_id,omitzero" json:"-"`
	// Return versions created at or after this time (inclusive).
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return versions created at or before this time (inclusive).
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// The maximum number of versions to return per page, from 1 to 100. Defaults
	// to 20.
	//
	// With `view=full`, a page holds at most 20 versions.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Return only versions of the memory with this ID (`mem_...`).
	//
	// The filter still works after the memory is deleted. While that memory's versions
	// are retained, the results include the version whose `operation` is `deleted`.
	MemoryID param.Opt[string] `query:"memory_id,omitzero" json:"-"`
	// The `next_page` value from a previous response, to get the next page. Omit it to
	// get the first page.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Return only versions written by the service account with this ID (`svac_...`).
	ServiceAccountID param.Opt[string] `query:"service_account_id,omitzero" json:"-"`
	// Return only versions written by the session with this ID.
	SessionID param.Opt[string] `query:"session_id,omitzero" json:"-"`
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// Return only versions that record this kind of change.
	//
	// Any of "created", "modified", "deleted".
	Operation BetaManagedAgentsMemoryVersionOperation `query:"operation,omitzero" json:"-"`
	// Whether a `memory` or `memory_version` in the response includes its `content`.
	// Defaults to `full` on retrieve operations and to `basic` on list, create, and
	// update operations.
	//
	// With `full`, list operations return at most 20 items per page.
	//
	// Any of "basic", "full".
	View BetaManagedAgentsMemoryView `query:"view,omitzero" json:"-"`
	// This endpoint is in beta: requests must send `agent-memory-2026-07-22` or
	// `managed-agents-2026-04-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaMemoryStoreMemoryVersionListParams]'s query parameters
// as `url.Values`.
func (r BetaMemoryStoreMemoryVersionListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaMemoryStoreMemoryVersionRedactParams struct {
	// The ID of the memory store that holds the version (`memstore_...`).
	MemoryStoreID string `path:"memory_store_id" api:"required" json:"-"`
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// This endpoint is in beta: requests must send `agent-memory-2026-07-22` or
	// `managed-agents-2026-04-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
