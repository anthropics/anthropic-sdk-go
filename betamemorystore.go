package anthropic

import (
	"context"
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

// BetaMemoryStoreService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaMemoryStoreService] method instead.
type BetaMemoryStoreService struct {
	Options        []option.RequestOption
	Memories       BetaMemoryStoreMemoryService
	MemoryVersions BetaMemoryStoreMemoryVersionService
}

// NewBetaMemoryStoreService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaMemoryStoreService(opts ...option.RequestOption) (r BetaMemoryStoreService) {
	r = BetaMemoryStoreService{}
	r.Options = opts
	r.Memories = NewBetaMemoryStoreMemoryService(opts...)
	r.MemoryVersions = NewBetaMemoryStoreMemoryVersionService(opts...)
	return
}

// Create a Memory Store
func (r *BetaMemoryStoreService) New(ctx context.Context, params BetaMemoryStoreNewParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemoryStore, err error) {
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
	path := "v1/memory_stores?beta=true"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.memory_stores.create", "/v1/memory_stores"), requestconfig.WithSpanCreatedID("gen_ai.memory.store.id", &res, func(created *BetaManagedAgentsMemoryStore) string { return created.ID }))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Retrieve a Memory Store
func (r *BetaMemoryStoreService) Get(ctx context.Context, memoryStoreID string, query BetaMemoryStoreGetParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemoryStore, err error) {
	if len(query.Betas) > 0 {
		headerValues := make([]string, len(query.Betas))
		for i, v := range query.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	if !param.IsOmitted(query.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", query.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "agent-memory-2026-07-22")}, opts...)
	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s?beta=true", url.PathEscape(memoryStoreID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.memory_stores.retrieve", "/v1/memory_stores/{memory_store_id}"), requestconfig.WithSpanAttribute("gen_ai.memory.store.id", memoryStoreID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a Memory Store
func (r *BetaMemoryStoreService) Update(ctx context.Context, memoryStoreID string, params BetaMemoryStoreUpdateParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemoryStore, err error) {
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
	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s?beta=true", url.PathEscape(memoryStoreID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.memory_stores.update", "/v1/memory_stores/{memory_store_id}"), requestconfig.WithSpanAttribute("gen_ai.memory.store.id", memoryStoreID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List Memory Stores
func (r *BetaMemoryStoreService) List(ctx context.Context, params BetaMemoryStoreListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaManagedAgentsMemoryStore], err error) {
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
	path := "v1/memory_stores?beta=true"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.memory_stores.list", "/v1/memory_stores", "created_at[gte]", "created_at[lte]", "include_archived", "limit", "page"))
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

// List Memory Stores
func (r *BetaMemoryStoreService) ListAutoPaging(ctx context.Context, params BetaMemoryStoreListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsMemoryStore] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Delete a Memory Store
func (r *BetaMemoryStoreService) Delete(ctx context.Context, memoryStoreID string, body BetaMemoryStoreDeleteParams, opts ...option.RequestOption) (res *BetaManagedAgentsDeletedMemoryStore, err error) {
	if len(body.Betas) > 0 {
		headerValues := make([]string, len(body.Betas))
		for i, v := range body.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	if !param.IsOmitted(body.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", body.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "agent-memory-2026-07-22")}, opts...)
	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s?beta=true", url.PathEscape(memoryStoreID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.memory_stores.delete", "/v1/memory_stores/{memory_store_id}"), requestconfig.WithSpanAttribute("gen_ai.memory.store.id", memoryStoreID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Archive a Memory Store
func (r *BetaMemoryStoreService) Archive(ctx context.Context, memoryStoreID string, body BetaMemoryStoreArchiveParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemoryStore, err error) {
	if len(body.Betas) > 0 {
		headerValues := make([]string, len(body.Betas))
		for i, v := range body.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	if !param.IsOmitted(body.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", body.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "agent-memory-2026-07-22")}, opts...)
	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s/archive?beta=true", url.PathEscape(memoryStoreID))
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.memory_stores.archive", "/v1/memory_stores/{memory_store_id}/archive"), requestconfig.WithSpanAttribute("gen_ai.memory.store.id", memoryStoreID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Confirmation that a `memory_store` was deleted.
type BetaManagedAgentsDeletedMemoryStore struct {
	// ID of the deleted memory store (a `memstore_...` identifier). The store and all
	// its memories and versions are no longer retrievable.
	ID   string                                  `json:"id" api:"required"`
	Type constant.MemoryStoreDeletedResponseType `json:"type" default:"memory_store_deleted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsDeletedMemoryStore) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsDeletedMemoryStore) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A `memory_store`: a named container for agent memories, scoped to a workspace.
// Attach a store to a session via `resources[]` to mount it as a directory the
// agent can read and write.
type BetaManagedAgentsMemoryStore struct {
	// The unique ID of the memory store (`memstore_...`). Use it to attach the store
	// to a session.
	ID string `json:"id" api:"required"`
	// When the memory store was archived, in RFC 3339, or `null` if it hasn't been.
	//
	// Archiving can't be undone. An archived store can't be attached to new sessions.
	// You can't create or update its memories or update the store. You can still read
	// and delete its memories, redact their versions, and delete the store.
	//
	// See the
	// [memory guide](https://platform.claude.com/docs/en/managed-agents/memory#archive-a-store)
	// to learn more about archiving.
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// When the memory store was created, in RFC 3339.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Free-text description of what the store contains, up to 1,024 characters.
	// Included in the agent's system prompt when the store is attached, so word it to
	// be useful to the agent. Empty string when unset.
	Description string `json:"description" api:"required"`
	// Key-value pairs for your own use, such as the end user that the store belongs
	// to.
	//
	// A store has at most 16 pairs. Keys are 1 to 64 characters long and values at
	// most 512 characters. You can't filter lists by metadata.
	Metadata map[string]string `json:"metadata" api:"required"`
	// Human-readable name for the store. 1–255 characters. The store's mount-path slug
	// under `/mnt/memory/` is derived from this name.
	Name string               `json:"name" api:"required"`
	Type constant.MemoryStore `json:"type" default:"memory_store"`
	// When the memory store was last changed, in RFC 3339.
	//
	// Changing its `name`, `description`, or `metadata` updates it. Archiving the
	// store does too. Changes to the store's memories don't.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ArchivedAt  respjson.Field
		CreatedAt   respjson.Field
		Description respjson.Field
		Metadata    respjson.Field
		Name        respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMemoryStore) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsMemoryStore) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaMemoryStoreNewParams struct {
	// The name of the memory store, 1 to 255 characters long, without control
	// characters. Names don't have to be unique in a workspace.
	//
	// In a session, the store's directory under `/mnt/memory/` is based on this name.
	// The name is lowercased and every run of characters other than ASCII letters and
	// digits becomes one hyphen. The session's memory store resource gives the exact
	// path in `mount_path`.
	//
	// See the
	// [memory guide](https://platform.claude.com/docs/en/managed-agents/memory#how-the-agent-accesses-memory)
	// to learn more about how agents access memory stores.
	Name string `json:"name" api:"required"`
	// Free-text description of what the store contains, up to 1,024 characters.
	// Included in the agent's system prompt when the store is attached, so word it to
	// be useful to the agent.
	Description param.Opt[string] `json:"description,omitzero"`
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// Key-value pairs for your own use, such as the end user that the store belongs
	// to. The agent doesn't see them.
	//
	// A store has at most 16 pairs. Keys are 1 to 64 characters long and values at
	// most 512 characters.
	Metadata map[string]string `json:"metadata,omitzero"`
	// This endpoint is in beta: requests must send `agent-memory-2026-07-22` or
	// `managed-agents-2026-04-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaMemoryStoreNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaMemoryStoreNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaMemoryStoreNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaMemoryStoreGetParams struct {
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

type BetaMemoryStoreUpdateParams struct {
	// New description for the store, up to 1,024 characters. Pass an empty string to
	// clear it.
	Description param.Opt[string] `json:"description,omitzero"`
	// The new name of the memory store, 1 to 255 characters long, without control
	// characters.
	//
	// Sessions created after the change use the new name for the store's directory
	// under `/mnt/memory/`.
	Name param.Opt[string] `json:"name,omitzero"`
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// Changes to apply to the memory store's metadata. Keys you leave out keep their
	// values.
	//
	// Set a key to a string to add or replace it, or to `null` to remove it. After the
	// change, the metadata can hold at most 16 key-value pairs. Keys are 1 to 64
	// characters long and values at most 512 characters.
	Metadata map[string]string `json:"metadata,omitzero"`
	// This endpoint is in beta: requests must send `agent-memory-2026-07-22` or
	// `managed-agents-2026-04-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaMemoryStoreUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaMemoryStoreUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaMemoryStoreUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaMemoryStoreListParams struct {
	// Return only memory stores created at or after this time (inclusive).
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return only memory stores created at or before this time (inclusive).
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// When `true`, archived stores are included in the results. Defaults to `false`
	// (archived stores are excluded).
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum number of stores to return per page. Must be between 1 and 100. Defaults
	// to 20 when omitted.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The `next_page` value from a previous response, to get the next page. Omit it to
	// get the first page.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
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

// URLQuery serializes [BetaMemoryStoreListParams]'s query parameters as
// `url.Values`.
func (r BetaMemoryStoreListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaMemoryStoreDeleteParams struct {
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

type BetaMemoryStoreArchiveParams struct {
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
