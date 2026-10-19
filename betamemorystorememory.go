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

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
)

// BetaMemoryStoreMemoryService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaMemoryStoreMemoryService] method instead.
type BetaMemoryStoreMemoryService struct {
	Options []option.RequestOption
}

// NewBetaMemoryStoreMemoryService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaMemoryStoreMemoryService(opts ...option.RequestOption) (r BetaMemoryStoreMemoryService) {
	r = BetaMemoryStoreMemoryService{}
	r.Options = opts
	return
}

// Create a Memory
func (r *BetaMemoryStoreMemoryService) New(ctx context.Context, memoryStoreID string, params BetaMemoryStoreMemoryNewParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemory, err error) {
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
	path := fmt.Sprintf("v1/memory_stores/%s/memories?beta=true", url.PathEscape(memoryStoreID))
	opts = append(opts, requestconfig.WithSpan("anthropic.memory_stores.memories.create", "/v1/memory_stores/{memory_store_id}/memories", "view"), requestconfig.WithSpanAttribute("gen_ai.memory.store.id", memoryStoreID), requestconfig.WithSpanCreatedID("gen_ai.memory.record.id", &res, func(created *BetaManagedAgentsMemory) string { return created.ID }))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Retrieve a Memory
func (r *BetaMemoryStoreMemoryService) Get(ctx context.Context, memoryID string, params BetaMemoryStoreMemoryGetParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemory, err error) {
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
	if memoryID == "" {
		err = errors.New("missing required memory_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s/memories/%s?beta=true", url.PathEscape(params.MemoryStoreID), url.PathEscape(memoryID))
	opts = append(opts, requestconfig.WithSpan("anthropic.memory_stores.memories.retrieve", "/v1/memory_stores/{memory_store_id}/memories/{memory_id}", "view"), requestconfig.WithSpanAttribute("gen_ai.memory.store.id", params.MemoryStoreID), requestconfig.WithSpanAttribute("gen_ai.memory.record.id", memoryID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Update a Memory
func (r *BetaMemoryStoreMemoryService) Update(ctx context.Context, memoryID string, params BetaMemoryStoreMemoryUpdateParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemory, err error) {
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
	if memoryID == "" {
		err = errors.New("missing required memory_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s/memories/%s?beta=true", url.PathEscape(params.MemoryStoreID), url.PathEscape(memoryID))
	opts = append(opts, requestconfig.WithSpan("anthropic.memory_stores.memories.update", "/v1/memory_stores/{memory_store_id}/memories/{memory_id}", "view"), requestconfig.WithSpanAttribute("gen_ai.memory.store.id", params.MemoryStoreID), requestconfig.WithSpanAttribute("gen_ai.memory.record.id", memoryID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List Memories
func (r *BetaMemoryStoreMemoryService) List(ctx context.Context, memoryStoreID string, params BetaMemoryStoreMemoryListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaManagedAgentsMemoryListItemUnion], err error) {
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
	path := fmt.Sprintf("v1/memory_stores/%s/memories?beta=true", url.PathEscape(memoryStoreID))
	opts = append(opts, requestconfig.WithSpan("anthropic.memory_stores.memories.list", "/v1/memory_stores/{memory_store_id}/memories", "depth", "limit", "page", "path_prefix", "view"), requestconfig.WithSpanAttribute("gen_ai.memory.store.id", memoryStoreID))
	opts = append(opts, requestconfig.WithSpanContentQuery("page"))
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

// List Memories
func (r *BetaMemoryStoreMemoryService) ListAutoPaging(ctx context.Context, memoryStoreID string, params BetaMemoryStoreMemoryListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsMemoryListItemUnion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, memoryStoreID, params, opts...))
}

// Delete a Memory
func (r *BetaMemoryStoreMemoryService) Delete(ctx context.Context, memoryID string, params BetaMemoryStoreMemoryDeleteParams, opts ...option.RequestOption) (res *BetaManagedAgentsDeletedMemory, err error) {
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
	if memoryID == "" {
		err = errors.New("missing required memory_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s/memories/%s?beta=true", url.PathEscape(params.MemoryStoreID), url.PathEscape(memoryID))
	opts = append(opts, requestconfig.WithSpan("anthropic.memory_stores.memories.delete", "/v1/memory_stores/{memory_store_id}/memories/{memory_id}", "expected_content_sha256"), requestconfig.WithSpanAttribute("gen_ai.memory.store.id", params.MemoryStoreID), requestconfig.WithSpanAttribute("gen_ai.memory.record.id", memoryID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, params, &res, opts...)
	return res, err
}

// The confirmation that
// [Delete a Memory](https://platform.claude.com/docs/en/api/beta/memory_stores/memories/delete)
// returns.
//
// Deleting a memory doesn't delete its versions.
// [List Memory Versions](https://platform.claude.com/docs/en/api/beta/memory_stores/memory_versions/list)
// returns them while they are retained: each version is kept for at least the
// version retention period after it was written, unless the store itself is
// deleted.
//
// See the
// [memory guide](https://platform.claude.com/docs/en/managed-agents/memory#audit-memory-changes)
// to learn more about how long versions are kept.
type BetaManagedAgentsDeletedMemory struct {
	// ID of the deleted memory (a `mem_...` value).
	ID string `json:"id" api:"required"`
	// Any of "memory_deleted".
	Type BetaManagedAgentsDeletedMemoryType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsDeletedMemory) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsDeletedMemory) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsDeletedMemoryType string

const (
	BetaManagedAgentsDeletedMemoryTypeMemoryDeleted BetaManagedAgentsDeletedMemoryType = "memory_deleted"
)

// A text document at a path inside a memory store.
//
// `content` is `null` when the response uses `view=basic`. `content_size_bytes`
// and `content_sha256` are always set, so you can detect changes without fetching
// the content. To retrieve, update, or delete a memory, pass its ID, which stays
// the same when its path changes.
//
// See the
// [memory guide](https://platform.claude.com/docs/en/managed-agents/memory#view-and-edit-memories)
// to learn more about working with memories.
type BetaManagedAgentsMemory struct {
	// The unique ID of the memory (`mem_...`). It stays the same when the memory's
	// path changes.
	//
	// To retrieve, update, or delete the memory, pass this ID.
	ID string `json:"id" api:"required"`
	// Lowercase hex SHA-256 digest of the UTF-8 `content` bytes (64 characters). The
	// server applies no normalization, so clients can compute the same hash locally
	// for staleness checks and as the value for a `content_sha256` precondition on
	// update. Always populated, regardless of `view`.
	ContentSha256 string `json:"content_sha256" api:"required"`
	// The size of the memory's `content` in UTF-8 bytes. It is set regardless of
	// `view`.
	ContentSizeBytes int64 `json:"content_size_bytes" api:"required"`
	// When this memory was created, in RFC 3339 format.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// ID of the memory store this memory belongs to (a `memstore_...` value).
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// The ID of the memory version that holds the memory's current content
	// (`memver_...`).
	//
	// Memory versions don't have a flag that marks the latest one. Compare their IDs
	// with this field instead.
	//
	// See the
	// [memory guide](https://platform.claude.com/docs/en/managed-agents/memory#audit-memory-changes)
	// to learn more about memory versions.
	MemoryVersionID string `json:"memory_version_id" api:"required"`
	// Hierarchical path of the memory within the store, e.g. `/projects/foo/notes.md`.
	// Always starts with `/`. Paths are case-sensitive and unique within a store.
	// Maximum 1,024 bytes.
	Path string `json:"path" api:"required"`
	// Any of "memory".
	Type BetaManagedAgentsMemoryType `json:"type" api:"required"`
	// When the memory was last changed, in RFC 3339.
	//
	// To find out who made the change, look at `created_by` on the memory version that
	// `memory_version_id` references.
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// The memory's UTF-8 text content, or `null` when the response uses `view=basic`.
	//
	// The maximum size is 100 kB (102,400 bytes).
	Content string `json:"content" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		ContentSha256    respjson.Field
		ContentSizeBytes respjson.Field
		CreatedAt        respjson.Field
		MemoryStoreID    respjson.Field
		MemoryVersionID  respjson.Field
		Path             respjson.Field
		Type             respjson.Field
		UpdatedAt        respjson.Field
		Content          respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMemory) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsMemory) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMemoryType string

const (
	BetaManagedAgentsMemoryTypeMemory BetaManagedAgentsMemoryType = "memory"
)

// BetaManagedAgentsMemoryListItemUnion contains all possible properties and values
// from [BetaManagedAgentsMemory], [BetaManagedAgentsMemoryPrefix].
//
// Use the [BetaManagedAgentsMemoryListItemUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsMemoryListItemUnion struct {
	// This field is from variant [BetaManagedAgentsMemory].
	ID string `json:"id"`
	// This field is from variant [BetaManagedAgentsMemory].
	ContentSha256 string `json:"content_sha256"`
	// This field is from variant [BetaManagedAgentsMemory].
	ContentSizeBytes int64 `json:"content_size_bytes"`
	// This field is from variant [BetaManagedAgentsMemory].
	CreatedAt time.Time `json:"created_at"`
	// This field is from variant [BetaManagedAgentsMemory].
	MemoryStoreID string `json:"memory_store_id"`
	// This field is from variant [BetaManagedAgentsMemory].
	MemoryVersionID string `json:"memory_version_id"`
	Path            string `json:"path"`
	// Any of "memory", "memory_prefix".
	Type string `json:"type"`
	// This field is from variant [BetaManagedAgentsMemory].
	UpdatedAt time.Time `json:"updated_at"`
	// This field is from variant [BetaManagedAgentsMemory].
	Content string `json:"content"`
	JSON    struct {
		ID               respjson.Field
		ContentSha256    respjson.Field
		ContentSizeBytes respjson.Field
		CreatedAt        respjson.Field
		MemoryStoreID    respjson.Field
		MemoryVersionID  respjson.Field
		Path             respjson.Field
		Type             respjson.Field
		UpdatedAt        respjson.Field
		Content          respjson.Field
		raw              string
	} `json:"-"`
}

// anyBetaManagedAgentsMemoryListItem is implemented by each variant of
// [BetaManagedAgentsMemoryListItemUnion] to add type safety for the return type of
// [BetaManagedAgentsMemoryListItemUnion.AsAny]
type anyBetaManagedAgentsMemoryListItem interface {
	implBetaManagedAgentsMemoryListItemUnion()
}

func (BetaManagedAgentsMemory) implBetaManagedAgentsMemoryListItemUnion()       {}
func (BetaManagedAgentsMemoryPrefix) implBetaManagedAgentsMemoryListItemUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsMemoryListItemUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsMemory:
//	case anthropic.BetaManagedAgentsMemoryPrefix:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaManagedAgentsMemoryListItemUnion) AsAny() anyBetaManagedAgentsMemoryListItem {
	switch u.Type {
	case "memory":
		return u.AsMemory()
	case "memory_prefix":
		return u.AsMemoryPrefix()
	}
	return nil
}

func (u BetaManagedAgentsMemoryListItemUnion) AsMemory() (v BetaManagedAgentsMemory) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaManagedAgentsMemoryListItemUnion) AsMemoryPrefix() (v BetaManagedAgentsMemoryPrefix) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsMemoryListItemUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsMemoryListItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A directory marker that
// [List Memories](https://platform.claude.com/docs/en/api/beta/memory_stores/memories/list)
// returns instead of the memories below the requested `depth`.
//
// It stands for one or more memories whose paths start with `path`. It isn't a
// stored object and doesn't have an ID. Each marker counts toward the page
// `limit`.
type BetaManagedAgentsMemoryPrefix struct {
	// The rolled-up path prefix, including a trailing `/` (e.g. `/projects/foo/`).
	// Pass this value as `path_prefix` on a subsequent list call to drill into the
	// directory.
	Path string `json:"path" api:"required"`
	// Any of "memory_prefix".
	Type BetaManagedAgentsMemoryPrefixType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Path        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMemoryPrefix) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsMemoryPrefix) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMemoryPrefixType string

const (
	BetaManagedAgentsMemoryPrefixTypeMemoryPrefix BetaManagedAgentsMemoryPrefixType = "memory_prefix"
)

// Whether a `memory` or `memory_version` in the response includes its `content`.
// Defaults to `full` on retrieve operations and to `basic` on list, create, and
// update operations.
//
// With `full`, list operations return at most 20 items per page.
type BetaManagedAgentsMemoryView string

const (
	// Return the object with `content` set to `null`. The `content_size_bytes` and
	// `content_sha256` fields remain populated, so sync clients can diff without
	// fetching content.
	BetaManagedAgentsMemoryViewBasic BetaManagedAgentsMemoryView = "basic"
	// Return the object with `content` populated. On list endpoints, `view=full` caps
	// `limit` at 20.
	BetaManagedAgentsMemoryViewFull BetaManagedAgentsMemoryView = "full"
)

// Optional condition that must hold for an update to apply. When omitted, the
// update is unconditional. Asserts the current state of the memory being updated.
// When an update changes `path`, the precondition still refers to the memory's
// current content, not the destination path. Currently the only supported variant
// is `content_sha256`.
//
// The property Type is required.
type BetaManagedAgentsPreconditionParam struct {
	// Any of "content_sha256".
	Type BetaManagedAgentsPreconditionType `json:"type,omitzero" api:"required"`
	// Expected `content_sha256` of the stored memory (64 lowercase hexadecimal
	// characters). Typically the `content_sha256` returned by a prior read or list
	// call. Because the server applies no content normalization, clients can also
	// compute this locally as the SHA-256 of the UTF-8 content bytes.
	ContentSha256 param.Opt[string] `json:"content_sha256,omitzero"`
	paramObj
}

func (r BetaManagedAgentsPreconditionParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaManagedAgentsPreconditionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaManagedAgentsPreconditionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsPreconditionType string

const (
	BetaManagedAgentsPreconditionTypeContentSha256 BetaManagedAgentsPreconditionType = "content_sha256"
)

type BetaMemoryStoreMemoryNewParams struct {
	// The UTF-8 text content of the new memory. Pass `""` to create an empty memory.
	//
	// The maximum size is 100 kB (102,400 bytes).
	Content param.Opt[string] `json:"content,omitzero" api:"required"`
	// The path of the new memory, such as `/projects/foo/notes.md`. Paths are
	// case-sensitive.
	//
	// A path starts with `/` and has 1 to 20 segments. A segment can't be empty and
	// can't be `.` or `..`. A path can't contain control characters, format
	// characters, or the line and paragraph separators U+2028 and U+2029. It must be
	// in Unicode Normalization Form C and at most 1,024 bytes long.
	Path string `json:"path" api:"required"`
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

func (r BetaMemoryStoreMemoryNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaMemoryStoreMemoryNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaMemoryStoreMemoryNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// URLQuery serializes [BetaMemoryStoreMemoryNewParams]'s query parameters as
// `url.Values`.
func (r BetaMemoryStoreMemoryNewParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaMemoryStoreMemoryGetParams struct {
	// The ID of the memory store that holds the memory (`memstore_...`).
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

// URLQuery serializes [BetaMemoryStoreMemoryGetParams]'s query parameters as
// `url.Values`.
func (r BetaMemoryStoreMemoryGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaMemoryStoreMemoryUpdateParams struct {
	// The ID of the memory store that holds the memory (`memstore_...`).
	MemoryStoreID string `path:"memory_store_id" api:"required" json:"-"`
	// New UTF-8 text content for the memory. Maximum 100 kB (102,400 bytes). Omit to
	// leave the content unchanged (e.g., for a rename-only update).
	Content param.Opt[string] `json:"content,omitzero"`
	// The new path of the memory, which renames it. Omit it to keep the current path.
	//
	// The memory keeps its `id`. A path starts with `/` and has 1 to 20 segments. A
	// segment can't be empty and can't be `.` or `..`. A path can't contain control
	// characters, format characters, or the line and paragraph separators U+2028 and
	// U+2029. It must be in Unicode Normalization Form C and at most 1,024 bytes long.
	// Paths are case-sensitive.
	Path param.Opt[string] `json:"path,omitzero"`
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
	// The condition that the memory must meet for the update to apply. Omit it to
	// update unconditionally.
	//
	// If the condition fails, the request fails with HTTP status 409.
	Precondition BetaManagedAgentsPreconditionParam `json:"precondition,omitzero"`
	// This endpoint is in beta: requests must send `agent-memory-2026-07-22` or
	// `managed-agents-2026-04-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaMemoryStoreMemoryUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaMemoryStoreMemoryUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaMemoryStoreMemoryUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// URLQuery serializes [BetaMemoryStoreMemoryUpdateParams]'s query parameters as
// `url.Values`.
func (r BetaMemoryStoreMemoryUpdateParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaMemoryStoreMemoryListParams struct {
	// How far below `path_prefix` to list. Defaults to `0`, which lists every memory
	// below the prefix.
	//
	// Pass `1` to list only the direct children. Deeper memories are then grouped into
	// `memory_prefix` items.
	Depth param.Opt[int64] `query:"depth,omitzero" json:"-"`
	// The maximum number of items to return per page, from 1 to 100. Defaults to 20.
	//
	// With `view=full`, a page holds at most 20 items. `memory` and `memory_prefix`
	// items both count toward the limit.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The `next_page` value from a previous response, to get the next page. Omit it to
	// get the first page.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Return only memories whose path starts with this prefix, such as `/notes/`. This
	// value appears in request URLs. Do not include secrets or personally identifiable
	// information.
	//
	// Omit it to list the whole store. The prefix must end with `/` and matches whole
	// segments, so `/notes/` matches `/notes/todo.md` but not
	// `/notes-archive/todo.md`.
	PathPrefix param.Opt[string] `query:"path_prefix,omitzero" json:"-"`
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// Whether each memory in the response includes `content`. Defaults to `basic`,
	// which sets it to `null`.
	//
	// With `full`, a page holds at most 20 items.
	//
	// Any of "basic", "full".
	View BetaManagedAgentsMemoryView `query:"view,omitzero" json:"-"`
	// This endpoint is in beta: requests must send `agent-memory-2026-07-22` or
	// `managed-agents-2026-04-01` in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaMemoryStoreMemoryListParams]'s query parameters as
// `url.Values`.
func (r BetaMemoryStoreMemoryListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaMemoryStoreMemoryDeleteParams struct {
	// The ID of the memory store that holds the memory (`memstore_...`).
	MemoryStoreID string `path:"memory_store_id" api:"required" json:"-"`
	// Delete the memory only if its current `content_sha256` equals this value, given
	// as 64 lowercase hexadecimal characters. Omit it to delete unconditionally.
	//
	// If the hashes differ, the request fails with HTTP status 409 and nothing is
	// deleted.
	ExpectedContentSha256 param.Opt[string] `query:"expected_content_sha256,omitzero" json:"-"`
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

// URLQuery serializes [BetaMemoryStoreMemoryDeleteParams]'s query parameters as
// `url.Values`.
func (r BetaMemoryStoreMemoryDeleteParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
