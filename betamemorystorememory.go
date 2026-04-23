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

// CreateMemory
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

// GetMemory
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

// UpdateMemory
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

// ListMemories
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

// ListMemories
func (r *BetaMemoryStoreMemoryService) ListAutoPaging(ctx context.Context, memoryStoreID string, params BetaMemoryStoreMemoryListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsMemoryListItemUnion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, memoryStoreID, params, opts...))
}

// DeleteMemory
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

type BetaManagedAgentsDeletedMemory struct {
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

type BetaManagedAgentsMemory struct {
	ID               string `json:"id" api:"required"`
	ContentSha256    string `json:"content_sha256" api:"required"`
	ContentSizeBytes int64  `json:"content_size_bytes" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt       time.Time `json:"created_at" api:"required" format:"date-time"`
	MemoryStoreID   string    `json:"memory_store_id" api:"required"`
	MemoryVersionID string    `json:"memory_version_id" api:"required"`
	Path            string    `json:"path" api:"required"`
	// Any of "memory".
	Type BetaManagedAgentsMemoryType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	Content   string    `json:"content" api:"nullable"`
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

type BetaManagedAgentsMemoryPrefix struct {
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

// MemoryView enum
type BetaManagedAgentsMemoryView string

const (
	BetaManagedAgentsMemoryViewBasic BetaManagedAgentsMemoryView = "basic"
	BetaManagedAgentsMemoryViewFull  BetaManagedAgentsMemoryView = "full"
)

// The property Type is required.
type BetaManagedAgentsPreconditionParam struct {
	// Any of "content_sha256".
	Type          BetaManagedAgentsPreconditionType `json:"type,omitzero" api:"required"`
	ContentSha256 param.Opt[string]                 `json:"content_sha256,omitzero"`
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
	Content param.Opt[string] `json:"content,omitzero" api:"required"`
	Path    string            `json:"path" api:"required"`
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
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
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
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaMemoryStoreMemoryUpdateParams struct {
	MemoryStoreID string            `path:"memory_store_id" api:"required" json:"-"`
	Content       param.Opt[string] `json:"content,omitzero"`
	Path          param.Opt[string] `json:"path,omitzero"`
	// Query parameter for view
	//
	// Any of "basic", "full".
	View         BetaManagedAgentsMemoryView        `query:"view,omitzero" json:"-"`
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
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
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
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
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
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
