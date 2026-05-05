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

// Create a memory
func (r *BetaMemoryStoreMemoryService) New(ctx context.Context, memoryStoreID string, params BetaMemoryStoreMemoryNewParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemory, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s/memories?beta=true", memoryStoreID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Retrieve a memory
func (r *BetaMemoryStoreMemoryService) Get(ctx context.Context, memoryID string, params BetaMemoryStoreMemoryGetParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemory, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if params.MemoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	if memoryID == "" {
		err = errors.New("missing required memory_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s/memories/%s?beta=true", params.MemoryStoreID, memoryID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Update a memory
func (r *BetaMemoryStoreMemoryService) Update(ctx context.Context, memoryID string, params BetaMemoryStoreMemoryUpdateParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemory, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if params.MemoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	if memoryID == "" {
		err = errors.New("missing required memory_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s/memories/%s?beta=true", params.MemoryStoreID, memoryID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List memories
func (r *BetaMemoryStoreMemoryService) List(ctx context.Context, memoryStoreID string, params BetaMemoryStoreMemoryListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaManagedAgentsMemoryListItemUnion], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01"), option.WithResponseInto(&raw)}, opts...)
	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s/memories?beta=true", memoryStoreID)
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

// List memories
func (r *BetaMemoryStoreMemoryService) ListAutoPaging(ctx context.Context, memoryStoreID string, params BetaMemoryStoreMemoryListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsMemoryListItemUnion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, memoryStoreID, params, opts...))
}

// Delete a memory
func (r *BetaMemoryStoreMemoryService) Delete(ctx context.Context, memoryID string, params BetaMemoryStoreMemoryDeleteParams, opts ...option.RequestOption) (res *BetaManagedAgentsDeletedMemory, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if params.MemoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	if memoryID == "" {
		err = errors.New("missing required memory_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s/memories/%s?beta=true", params.MemoryStoreID, memoryID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, params, &res, opts...)
	return res, err
}

// Tombstone returned by
// [Delete a memory](/en/api/beta/memory_stores/memories/delete). The memory's
// version history persists and remains listable via
// [List memory versions](/en/api/beta/memory_stores/memory_versions/list) until
// the store itself is deleted.
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

// A `memory` object: a single text document at a hierarchical path inside a memory
// store. The `content` field is populated when `view=full` and `null` when
// `view=basic`; the `content_size_bytes` and `content_sha256` fields are always
// populated so sync clients can diff without fetching content. Memories are
// addressed by their `mem_...` ID; the path is the create key and can be changed
// via update.
type BetaManagedAgentsMemory struct {
	// Unique identifier for this memory (a `mem_...` value). Stable across renames;
	// use this ID, not the path, to read, update, or delete the memory.
	ID string `json:"id" api:"required"`
	// Lowercase hex SHA-256 digest of the UTF-8 `content` bytes (64 characters). The
	// server applies no normalization, so clients can compute the same hash locally
	// for staleness checks and as the value for a `content_sha256` precondition on
	// update. Always populated, regardless of `view`.
	ContentSha256 string `json:"content_sha256" api:"required"`
	// Size of `content` in bytes (the UTF-8 plaintext length). Always populated,
	// regardless of `view`.
	ContentSizeBytes int64 `json:"content_size_bytes" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// ID of the memory store this memory belongs to (a `memstore_...` value).
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// ID of the `memory_version` representing this memory's current content (a
	// `memver_...` value). This is the authoritative head pointer; `memory_version`
	// objects do not carry an `is_latest` flag, so compare against this field instead.
	// Enumerate the full history via
	// [List memory versions](/en/api/beta/memory_stores/memory_versions/list).
	MemoryVersionID string `json:"memory_version_id" api:"required"`
	// Hierarchical path of the memory within the store, e.g. `/projects/foo/notes.md`.
	// Always starts with `/`. Paths are case-sensitive and unique within a store.
	// Maximum 1,024 bytes.
	Path string `json:"path" api:"required"`
	// Any of "memory".
	Type BetaManagedAgentsMemoryType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// The memory's UTF-8 text content. Populated when `view=full`; `null` when
	// `view=basic`. Maximum 100 kB (102,400 bytes).
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

// A rolled-up directory marker returned by
// [List memories](/en/api/beta/memory_stores/memories/list) when `depth` is set.
// Indicates that one or more memories exist deeper than the requested depth under
// this prefix. This is a list-time rollup, not a stored resource; it has no ID and
// no lifecycle. Each prefix counts toward the page `limit` and interleaves with
// `memory` items in path order.
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

// Selects which projection of a `memory` or `memory_version` the server returns.
// `basic` returns the object with `content` set to `null`; `full` populates
// `content`. When omitted, the default is endpoint-specific: retrieve operations
// default to `full`; list, create, and update operations default to `basic`.
// Listing with `view=full` caps `limit` at 20.
type BetaManagedAgentsMemoryView string

const (
	BetaManagedAgentsMemoryViewBasic BetaManagedAgentsMemoryView = "basic"
	BetaManagedAgentsMemoryViewFull  BetaManagedAgentsMemoryView = "full"
)

// Optimistic-concurrency precondition: the update applies only if the memory's
// stored `content_sha256` equals the supplied value. On mismatch, the request
// returns `memory_precondition_failed_error` (HTTP 409); re-read the memory and
// retry against the fresh state. If the precondition fails but the stored state
// already exactly matches the requested `content` and `path`, the server returns
// 200 instead of 409.
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
	// UTF-8 text content for the new memory. Maximum 100 kB (102,400 bytes). Required;
	// pass `""` explicitly to create an empty memory.
	Content param.Opt[string] `json:"content,omitzero" api:"required"`
	// Hierarchical path for the new memory, e.g. `/projects/foo/notes.md`. Must start
	// with `/`, contain at least one non-empty segment, and be at most 1,024 bytes.
	// Must not contain empty segments, `.` or `..` segments, control or format
	// characters, and must be NFC-normalized. Paths are case-sensitive.
	Path string `json:"path" api:"required"`
	// Query parameter for view
	//
	// Any of "basic", "full".
	View BetaManagedAgentsMemoryView `query:"view,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
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
	MemoryStoreID string `path:"memory_store_id" api:"required" json:"-"`
	// Query parameter for view
	//
	// Any of "basic", "full".
	View BetaManagedAgentsMemoryView `query:"view,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
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
	MemoryStoreID string `path:"memory_store_id" api:"required" json:"-"`
	// New UTF-8 text content for the memory. Maximum 100 kB (102,400 bytes). Omit to
	// leave the content unchanged (e.g., for a rename-only update).
	Content param.Opt[string] `json:"content,omitzero"`
	// New path for the memory (a rename). Must start with `/`, contain at least one
	// non-empty segment, and be at most 1,024 bytes. Must not contain empty segments,
	// `.` or `..` segments, control or format characters, and must be NFC-normalized.
	// Paths are case-sensitive. The memory's `id` is preserved across renames. Omit to
	// leave the path unchanged.
	Path param.Opt[string] `json:"path,omitzero"`
	// Query parameter for view
	//
	// Any of "basic", "full".
	View BetaManagedAgentsMemoryView `query:"view,omitzero" json:"-"`
	// Optimistic-concurrency precondition: the update applies only if the memory's
	// stored `content_sha256` equals the supplied value. On mismatch, the request
	// returns `memory_precondition_failed_error` (HTTP 409); re-read the memory and
	// retry against the fresh state. If the precondition fails but the stored state
	// already exactly matches the requested `content` and `path`, the server returns
	// 200 instead of 409.
	Precondition BetaManagedAgentsPreconditionParam `json:"precondition,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
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
	// Query parameter for depth
	Depth param.Opt[int64] `query:"depth,omitzero" json:"-"`
	// Query parameter for limit
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Query parameter for order_by
	OrderBy param.Opt[string] `query:"order_by,omitzero" json:"-"`
	// Query parameter for page
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional path prefix filter (raw string-prefix match; include a trailing slash
	// for directory-scoped lists). This value appears in request URLs. Do not include
	// secrets or personally identifiable information.
	PathPrefix param.Opt[string] `query:"path_prefix,omitzero" json:"-"`
	// Query parameter for order
	//
	// Any of "asc", "desc".
	Order BetaMemoryStoreMemoryListParamsOrder `query:"order,omitzero" json:"-"`
	// Query parameter for view
	//
	// Any of "basic", "full".
	View BetaManagedAgentsMemoryView `query:"view,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
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

// Query parameter for order
type BetaMemoryStoreMemoryListParamsOrder string

const (
	BetaMemoryStoreMemoryListParamsOrderAsc  BetaMemoryStoreMemoryListParamsOrder = "asc"
	BetaMemoryStoreMemoryListParamsOrderDesc BetaMemoryStoreMemoryListParamsOrder = "desc"
)

type BetaMemoryStoreMemoryDeleteParams struct {
	MemoryStoreID string `path:"memory_store_id" api:"required" json:"-"`
	// Query parameter for expected_content_sha256
	ExpectedContentSha256 param.Opt[string] `query:"expected_content_sha256,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
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
