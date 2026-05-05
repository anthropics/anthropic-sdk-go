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

// Retrieve a memory version
func (r *BetaMemoryStoreMemoryVersionService) Get(ctx context.Context, memoryVersionID string, params BetaMemoryStoreMemoryVersionGetParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemoryVersion, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if params.MemoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	if memoryVersionID == "" {
		err = errors.New("missing required memory_version_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s/memory_versions/%s?beta=true", params.MemoryStoreID, memoryVersionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// List memory versions
func (r *BetaMemoryStoreMemoryVersionService) List(ctx context.Context, memoryStoreID string, params BetaMemoryStoreMemoryVersionListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaManagedAgentsMemoryVersion], err error) {
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
	path := fmt.Sprintf("v1/memory_stores/%s/memory_versions?beta=true", memoryStoreID)
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

// List memory versions
func (r *BetaMemoryStoreMemoryVersionService) ListAutoPaging(ctx context.Context, memoryStoreID string, params BetaMemoryStoreMemoryVersionListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaManagedAgentsMemoryVersion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, memoryStoreID, params, opts...))
}

// Redact a memory version
func (r *BetaMemoryStoreMemoryVersionService) Redact(ctx context.Context, memoryVersionID string, params BetaMemoryStoreMemoryVersionRedactParams, opts ...option.RequestOption) (res *BetaManagedAgentsMemoryVersion, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "managed-agents-2026-04-01")}, opts...)
	if params.MemoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	if memoryVersionID == "" {
		err = errors.New("missing required memory_version_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/memory_stores/%s/memory_versions/%s/redact?beta=true", params.MemoryStoreID, memoryVersionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// BetaManagedAgentsActorUnion contains all possible properties and values from
// [BetaManagedAgentsSessionActor], [BetaManagedAgentsAPIActor],
// [BetaManagedAgentsUserActor].
//
// Use the [BetaManagedAgentsActorUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaManagedAgentsActorUnion struct {
	// This field is from variant [BetaManagedAgentsSessionActor].
	SessionID string `json:"session_id"`
	// Any of "session_actor", "api_actor", "user_actor".
	Type string `json:"type"`
	// This field is from variant [BetaManagedAgentsAPIActor].
	APIKeyID string `json:"api_key_id"`
	// This field is from variant [BetaManagedAgentsUserActor].
	UserID string `json:"user_id"`
	JSON   struct {
		SessionID respjson.Field
		Type      respjson.Field
		APIKeyID  respjson.Field
		UserID    respjson.Field
		raw       string
	} `json:"-"`
}

// anyBetaManagedAgentsActor is implemented by each variant of
// [BetaManagedAgentsActorUnion] to add type safety for the return type of
// [BetaManagedAgentsActorUnion.AsAny]
type anyBetaManagedAgentsActor interface {
	implBetaManagedAgentsActorUnion()
}

func (BetaManagedAgentsSessionActor) implBetaManagedAgentsActorUnion() {}
func (BetaManagedAgentsAPIActor) implBetaManagedAgentsActorUnion()     {}
func (BetaManagedAgentsUserActor) implBetaManagedAgentsActorUnion()    {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaManagedAgentsActorUnion.AsAny().(type) {
//	case anthropic.BetaManagedAgentsSessionActor:
//	case anthropic.BetaManagedAgentsAPIActor:
//	case anthropic.BetaManagedAgentsUserActor:
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

// Returns the unmodified JSON received from the API
func (u BetaManagedAgentsActorUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaManagedAgentsActorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Attribution for a write made directly via the public API (outside of any
// session).
type BetaManagedAgentsAPIActor struct {
	// ID of the API key that performed the write. This identifies the key, not the
	// secret.
	APIKeyID string `json:"api_key_id" api:"required"`
	// Any of "api_actor".
	Type BetaManagedAgentsAPIActorType `json:"type" api:"required"`
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

type BetaManagedAgentsAPIActorType string

const (
	BetaManagedAgentsAPIActorTypeAPIActor BetaManagedAgentsAPIActorType = "api_actor"
)

// A `memory_version` object: one immutable, attributed row in a memory's
// append-only history. Every non-no-op mutation to a memory produces a new
// version. Versions belong to the store (not the individual memory) and persist
// after the memory is deleted. Retrieving a redacted version returns 200 with
// `content`, `path`, `content_size_bytes`, and `content_sha256` set to `null`;
// branch on `redacted_at`, not HTTP status.
type BetaManagedAgentsMemoryVersion struct {
	// Unique identifier for this version (a `memver_...` value).
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// ID of the memory this version snapshots (a `mem_...` value). Remains valid after
	// the memory is deleted; pass it as `memory_id` to
	// [List memory versions](/en/api/beta/memory_stores/memory_versions/list) to
	// retrieve the full lineage including the `deleted` row.
	MemoryID string `json:"memory_id" api:"required"`
	// ID of the memory store this version belongs to (a `memstore_...` value).
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// The kind of mutation a `memory_version` records. Every non-no-op mutation to a
	// memory appends exactly one version row with one of these values.
	//
	// Any of "created", "modified", "deleted".
	Operation BetaManagedAgentsMemoryVersionOperation `json:"operation" api:"required"`
	// Any of "memory_version".
	Type BetaManagedAgentsMemoryVersionType `json:"type" api:"required"`
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
	// Identifies who performed a write or redact operation. Captured at write time on
	// the `memory_version` row. The API key that created a session is not recorded on
	// agent writes; attribution answers who made the write, not who is ultimately
	// responsible. Look up session provenance separately via the
	// [Sessions API](/en/api/sessions-retrieve).
	CreatedBy BetaManagedAgentsActorUnion `json:"created_by"`
	// The memory's path at the time of this write. `null` if and only if `redacted_at`
	// is set.
	Path string `json:"path" api:"nullable"`
	// A timestamp in RFC 3339 format
	RedactedAt time.Time `json:"redacted_at" api:"nullable" format:"date-time"`
	// Identifies who performed a write or redact operation. Captured at write time on
	// the `memory_version` row. The API key that created a session is not recorded on
	// agent writes; attribution answers who made the write, not who is ultimately
	// responsible. Look up session provenance separately via the
	// [Sessions API](/en/api/sessions-retrieve).
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

type BetaManagedAgentsMemoryVersionType string

const (
	BetaManagedAgentsMemoryVersionTypeMemoryVersion BetaManagedAgentsMemoryVersionType = "memory_version"
)

// The kind of mutation a `memory_version` records. Every non-no-op mutation to a
// memory appends exactly one version row with one of these values.
type BetaManagedAgentsMemoryVersionOperation string

const (
	BetaManagedAgentsMemoryVersionOperationCreated  BetaManagedAgentsMemoryVersionOperation = "created"
	BetaManagedAgentsMemoryVersionOperationModified BetaManagedAgentsMemoryVersionOperation = "modified"
	BetaManagedAgentsMemoryVersionOperationDeleted  BetaManagedAgentsMemoryVersionOperation = "deleted"
)

// Attribution for a write made by an agent during a session, through the mounted
// filesystem at `/mnt/memory/`.
type BetaManagedAgentsSessionActor struct {
	// ID of the session that performed the write (a `sesn_...` value). Look up the
	// session via [Retrieve a session](/en/api/sessions-retrieve) for further
	// provenance.
	SessionID string `json:"session_id" api:"required"`
	// Any of "session_actor".
	Type BetaManagedAgentsSessionActorType `json:"type" api:"required"`
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

type BetaManagedAgentsSessionActorType string

const (
	BetaManagedAgentsSessionActorTypeSessionActor BetaManagedAgentsSessionActorType = "session_actor"
)

// Attribution for a write made by a human user through the Anthropic Console.
type BetaManagedAgentsUserActor struct {
	// Any of "user_actor".
	Type BetaManagedAgentsUserActorType `json:"type" api:"required"`
	// ID of the user who performed the write (a `user_...` value).
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

type BetaManagedAgentsUserActorType string

const (
	BetaManagedAgentsUserActorTypeUserActor BetaManagedAgentsUserActorType = "user_actor"
)

type BetaMemoryStoreMemoryVersionGetParams struct {
	MemoryStoreID string `path:"memory_store_id" api:"required" json:"-"`
	// Query parameter for view
	//
	// Any of "basic", "full".
	View BetaManagedAgentsMemoryView `query:"view,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
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
	// Query parameter for api_key_id
	APIKeyID param.Opt[string] `query:"api_key_id,omitzero" json:"-"`
	// Return versions created at or after this time (inclusive).
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return versions created at or before this time (inclusive).
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// Query parameter for limit
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Query parameter for memory_id
	MemoryID param.Opt[string] `query:"memory_id,omitzero" json:"-"`
	// Query parameter for page
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Query parameter for session_id
	SessionID param.Opt[string] `query:"session_id,omitzero" json:"-"`
	// Query parameter for operation
	//
	// Any of "created", "modified", "deleted".
	Operation BetaManagedAgentsMemoryVersionOperation `query:"operation,omitzero" json:"-"`
	// Query parameter for view
	//
	// Any of "basic", "full".
	View BetaManagedAgentsMemoryView `query:"view,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
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
	MemoryStoreID string `path:"memory_store_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
