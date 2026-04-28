// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"context"
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

// Create a memory store
func (r *BetaMemoryStoreService) New(ctx context.Context, params BetaMemoryStoreNewParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemoryStore, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	path := "v1/memory_stores?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Retrieve a memory store
func (r *BetaMemoryStoreService) Get(ctx context.Context, memoryStoreID string, query BetaMemoryStoreGetParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemoryStore, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s?beta=true", memoryStoreID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a memory store
func (r *BetaMemoryStoreService) Update(ctx context.Context, memoryStoreID string, params BetaMemoryStoreUpdateParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemoryStore, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s?beta=true", memoryStoreID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List memory stores
func (r *BetaMemoryStoreService) List(ctx context.Context, params BetaMemoryStoreListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaManagedAgentsMemoryStore], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/memory_stores?beta=true"
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

// List memory stores
func (r *BetaMemoryStoreService) ListAutoPaging(ctx context.Context, params BetaMemoryStoreListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsMemoryStore] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Delete a memory store
func (r *BetaMemoryStoreService) Delete(ctx context.Context, memoryStoreID string, body BetaMemoryStoreDeleteParams, opts ...option.RequestOption) (res *BetaManagedAgentsDeletedMemoryStore, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s?beta=true", memoryStoreID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Archive a memory store
func (r *BetaMemoryStoreService) Archive(ctx context.Context, memoryStoreID string, body BetaMemoryStoreArchiveParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemoryStore, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s/archive?beta=true", memoryStoreID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Confirmation that a `memory_store` was deleted.
type BetaManagedAgentsDeletedMemoryStore struct {
	// ID of the deleted memory store (a `memstore_...` identifier). The store and all
	// its memories and versions are no longer retrievable.
	ID string `json:"id" api:"required"`
	// Any of "memory_store_deleted".
	Type BetaManagedAgentsDeletedMemoryStoreType `json:"type" api:"required"`
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

type BetaManagedAgentsDeletedMemoryStoreType string

const (
	BetaManagedAgentsDeletedMemoryStoreTypeMemoryStoreDeleted BetaManagedAgentsDeletedMemoryStoreType = "memory_store_deleted"
)

// A `memory_store`: a named container for agent memories, scoped to a workspace.
// Attach a store to a session via `resources[]` to mount it as a directory the
// agent can read and write.
type BetaManagedAgentsMemoryStore struct {
	// Unique identifier for the memory store (a `memstore_...` tagged ID). Use this
	// when attaching the store to a session, or in the `{memory_store_id}` path
	// parameter of subsequent calls.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Human-readable name for the store. 1–255 characters. The store's mount-path slug
	// under `/mnt/memory/` is derived from this name.
	Name string `json:"name" api:"required"`
	// Any of "memory_store".
	Type BetaManagedAgentsMemoryStoreType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"nullable" format:"date-time"`
	// Free-text description of what the store contains, up to 1024 characters.
	// Included in the agent's system prompt when the store is attached, so word it to
	// be useful to the agent. Empty string when unset.
	Description string `json:"description"`
	// Arbitrary key-value tags for your own bookkeeping (such as the end user a store
	// belongs to). Up to 16 pairs; keys 1–64 characters; values up to 512 characters.
	// Returned on retrieve/list but not filterable.
	Metadata map[string]string `json:"metadata"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		ArchivedAt  respjson.Field
		Description respjson.Field
		Metadata    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaManagedAgentsMemoryStore) RawJSON() string { return r.JSON.raw }
func (r *BetaManagedAgentsMemoryStore) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaManagedAgentsMemoryStoreType string

const (
	BetaManagedAgentsMemoryStoreTypeMemoryStore BetaManagedAgentsMemoryStoreType = "memory_store"
)

type BetaMemoryStoreNewParams struct {
	// Human-readable name for the store. Required; 1–255 characters; no control
	// characters. The mount-path slug under `/mnt/memory/` is derived from this name
	// (lowercased, non-alphanumeric runs collapsed to a hyphen). Names need not be
	// unique within a workspace.
	Name string `json:"name" api:"required"`
	// Free-text description of what the store contains, up to 1024 characters.
	// Included in the agent's system prompt when the store is attached, so word it to
	// be useful to the agent.
	Description param.Opt[string] `json:"description,omitzero"`
	// Arbitrary key-value tags for your own bookkeeping (such as the end user a store
	// belongs to). Up to 16 pairs; keys 1–64 characters; values up to 512 characters.
	// Not visible to the agent.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
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
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaMemoryStoreUpdateParams struct {
	// New description for the store, up to 1024 characters. Pass an empty string to
	// clear it.
	Description param.Opt[string] `json:"description,omitzero"`
	// New human-readable name for the store. 1–255 characters; no control characters.
	// Renaming changes the slug used for the store's `mount_path` in sessions created
	// after the update.
	Name param.Opt[string] `json:"name,omitzero"`
	// Metadata patch. Set a key to a string to upsert it, or to null to delete it.
	// Omit the field to preserve. The stored bag is limited to 16 keys (up to 64 chars
	// each) with values up to 512 chars.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
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
	// Return only stores whose `created_at` is at or after this time (inclusive). Sent
	// on the wire as `created_at[gte]`.
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return only stores whose `created_at` is at or before this time (inclusive).
	// Sent on the wire as `created_at[lte]`.
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// When `true`, archived stores are included in the results. Defaults to `false`
	// (archived stores are excluded).
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum number of stores to return per page. Must be between 1 and 100. Defaults
	// to 20 when omitted.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor (a `page_...` value). Pass the `next_page` value from a
	// previous response to fetch the next page; omit for the first page.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaMemoryStoreListParams]'s query parameters as
// `url.Values`.
func (r BetaMemoryStoreListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaMemoryStoreDeleteParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaMemoryStoreArchiveParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
