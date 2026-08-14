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

// BetaDreamService contains methods and other services that help with interacting
// with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaDreamService] method instead.
type BetaDreamService struct {
	Options []option.RequestOption
}

// NewBetaDreamService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewBetaDreamService(opts ...option.RequestOption) (r BetaDreamService) {
	r = BetaDreamService{}
	r.Options = opts
	return
}

// Create a Dream
func (r *BetaDreamService) New(ctx context.Context, params BetaDreamNewParams, opts ...option.RequestOption) (res *BetaDream, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "dreaming-2026-04-21")}, opts...)
	path := "v1/dreams?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get a Dream
func (r *BetaDreamService) Get(ctx context.Context, dreamID string, query BetaDreamGetParams, opts ...option.RequestOption) (res *BetaDream, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "dreaming-2026-04-21")}, opts...)
	if dreamID == "" {
		err = errors.New("missing required dream_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/dreams/%s?beta=true", dreamID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List Dreams
func (r *BetaDreamService) List(ctx context.Context, params BetaDreamListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaDream], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "dreaming-2026-04-21"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/dreams?beta=true"
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

// List Dreams
func (r *BetaDreamService) ListAutoPaging(ctx context.Context, params BetaDreamListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaDream] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Archive a Dream
func (r *BetaDreamService) Archive(ctx context.Context, dreamID string, body BetaDreamArchiveParams, opts ...option.RequestOption) (res *BetaDream, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "dreaming-2026-04-21")}, opts...)
	if dreamID == "" {
		err = errors.New("missing required dream_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/dreams/%s/archive?beta=true", dreamID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Cancel a Dream
func (r *BetaDreamService) Cancel(ctx context.Context, dreamID string, body BetaDreamCancelParams, opts ...option.RequestOption) (res *BetaDream, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "dreaming-2026-04-21")}, opts...)
	if dreamID == "" {
		err = errors.New("missing required dream_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/dreams/%s/cancel?beta=true", dreamID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// An asynchronous memory-consolidation job that reads a memory store plus a set of
// session transcripts and writes consolidated memories into an output memory store
// — a new store by default, or an existing store chosen via output_behavior. The
// Dreams API is in research preview: the request and response shapes are volatile
// and may change without the deprecation period that applies to
// generally-available endpoints.
type BetaDream struct {
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	EndedAt time.Time `json:"ended_at" api:"required" format:"date-time"`
	// Failure detail for a Dream whose `status` is `failed`.
	Error        BetaDreamError        `json:"error" api:"required"`
	Inputs       []BetaDreamInputUnion `json:"inputs" api:"required"`
	Instructions string                `json:"instructions" api:"required"`
	// Model identifier and configuration applied to every pipeline stage. Same wire
	// shape as the Agents API ModelConfig.
	Model BetaDreamModelConfig `json:"model" api:"required"`
	// The default destination: the job creates a new output memory store as a clone of
	// the memory_store input and writes the consolidated memories into it. The input
	// store is never mutated.
	OutputBehavior BetaOutputBehaviorUnion `json:"output_behavior" api:"required"`
	Outputs        []BetaDreamOutput       `json:"outputs" api:"required"`
	SessionID      string                  `json:"session_id" api:"required"`
	// Lifecycle status of a Dream.
	//
	// Any of "pending", "running", "completed", "failed", "canceled".
	Status BetaDreamStatus `json:"status" api:"required"`
	// Any of "dream".
	Type BetaDreamType `json:"type" api:"required"`
	// Cumulative token usage for the dream across every pipeline stage.
	Usage BetaDreamUsage `json:"usage" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		ArchivedAt     respjson.Field
		CreatedAt      respjson.Field
		EndedAt        respjson.Field
		Error          respjson.Field
		Inputs         respjson.Field
		Instructions   respjson.Field
		Model          respjson.Field
		OutputBehavior respjson.Field
		Outputs        respjson.Field
		SessionID      respjson.Field
		Status         respjson.Field
		Type           respjson.Field
		Usage          respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDream) RawJSON() string { return r.JSON.raw }
func (r *BetaDream) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaDreamType string

const (
	BetaDreamTypeDream BetaDreamType = "dream"
)

// Failure detail for a Dream whose `status` is `failed`.
type BetaDreamError struct {
	Message string `json:"message" api:"required"`
	Type    string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDreamError) RawJSON() string { return r.JSON.raw }
func (r *BetaDreamError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaDreamInputUnion contains all possible properties and values from
// [BetaDreamMemoryStoreInput], [BetaDreamSessionsInput].
//
// Use the [BetaDreamInputUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaDreamInputUnion struct {
	// This field is from variant [BetaDreamMemoryStoreInput].
	MemoryStoreID string `json:"memory_store_id"`
	// Any of "memory_store", "sessions".
	Type string `json:"type"`
	// This field is from variant [BetaDreamSessionsInput].
	SessionIDs []string `json:"session_ids"`
	JSON       struct {
		MemoryStoreID respjson.Field
		Type          respjson.Field
		SessionIDs    respjson.Field
		raw           string
	} `json:"-"`
}

// anyBetaDreamInput is implemented by each variant of [BetaDreamInputUnion] to add
// type safety for the return type of [BetaDreamInputUnion.AsAny]
type anyBetaDreamInput interface {
	implBetaDreamInputUnion()
}

func (BetaDreamMemoryStoreInput) implBetaDreamInputUnion() {}
func (BetaDreamSessionsInput) implBetaDreamInputUnion()    {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaDreamInputUnion.AsAny().(type) {
//	case anthropic.BetaDreamMemoryStoreInput:
//	case anthropic.BetaDreamSessionsInput:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaDreamInputUnion) AsAny() anyBetaDreamInput {
	switch u.Type {
	case "memory_store":
		return u.AsMemoryStore()
	case "sessions":
		return u.AsSessions()
	}
	return nil
}

func (u BetaDreamInputUnion) AsMemoryStore() (v BetaDreamMemoryStoreInput) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaDreamInputUnion) AsSessions() (v BetaDreamSessionsInput) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaDreamInputUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaDreamInputUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaDreamInputUnion to a BetaDreamInputUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaDreamInputUnionParam.Overrides()
func (r BetaDreamInputUnion) ToParam() BetaDreamInputUnionParam {
	return param.Override[BetaDreamInputUnionParam](json.RawMessage(r.RawJSON()))
}

func BetaDreamInputParamOfMemoryStore(memoryStoreID string) BetaDreamInputUnionParam {
	var memoryStore BetaDreamMemoryStoreInputParam
	memoryStore.MemoryStoreID = memoryStoreID
	return BetaDreamInputUnionParam{OfMemoryStore: &memoryStore}
}

func BetaDreamInputParamOfSessions(sessionIDs []string) BetaDreamInputUnionParam {
	var sessions BetaDreamSessionsInputParam
	sessions.SessionIDs = sessionIDs
	return BetaDreamInputUnionParam{OfSessions: &sessions}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaDreamInputUnionParam struct {
	OfMemoryStore *BetaDreamMemoryStoreInputParam `json:",omitzero,inline"`
	OfSessions    *BetaDreamSessionsInputParam    `json:",omitzero,inline"`
	paramUnion
}

func (u BetaDreamInputUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfMemoryStore, u.OfSessions)
}
func (u *BetaDreamInputUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaDreamInputUnionParam) asAny() any {
	if !param.IsOmitted(u.OfMemoryStore) {
		return u.OfMemoryStore
	} else if !param.IsOmitted(u.OfSessions) {
		return u.OfSessions
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDreamInputUnionParam) GetMemoryStoreID() *string {
	if vt := u.OfMemoryStore; vt != nil {
		return &vt.MemoryStoreID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDreamInputUnionParam) GetSessionIDs() []string {
	if vt := u.OfSessions; vt != nil {
		return vt.SessionIDs
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaDreamInputUnionParam) GetType() *string {
	if vt := u.OfMemoryStore; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessions; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaDreamInputUnionParam](
		"type",
		apijson.Discriminator[BetaDreamMemoryStoreInputParam]("memory_store"),
		apijson.Discriminator[BetaDreamSessionsInputParam]("sessions"),
	)
}

// An input memory store the dream reads from. The dream never mutates this store
// unless it is also the destination: with output_behavior {type:
// "update_existing"} the job consolidates this store in place.
type BetaDreamMemoryStoreInput struct {
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// Any of "memory_store".
	Type BetaDreamMemoryStoreInputType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MemoryStoreID respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDreamMemoryStoreInput) RawJSON() string { return r.JSON.raw }
func (r *BetaDreamMemoryStoreInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaDreamMemoryStoreInput to a
// BetaDreamMemoryStoreInputParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaDreamMemoryStoreInputParam.Overrides()
func (r BetaDreamMemoryStoreInput) ToParam() BetaDreamMemoryStoreInputParam {
	return param.Override[BetaDreamMemoryStoreInputParam](json.RawMessage(r.RawJSON()))
}

type BetaDreamMemoryStoreInputType string

const (
	BetaDreamMemoryStoreInputTypeMemoryStore BetaDreamMemoryStoreInputType = "memory_store"
)

// An input memory store the dream reads from. The dream never mutates this store
// unless it is also the destination: with output_behavior {type:
// "update_existing"} the job consolidates this store in place.
//
// The properties MemoryStoreID, Type are required.
type BetaDreamMemoryStoreInputParam struct {
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// Any of "memory_store".
	Type BetaDreamMemoryStoreInputType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaDreamMemoryStoreInputParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaDreamMemoryStoreInputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaDreamMemoryStoreInputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Model identifier and configuration applied to every pipeline stage. Same wire
// shape as the Agents API ModelConfig.
type BetaDreamModelConfig struct {
	// Model identifier, e.g. "claude-opus-5". 1-256 characters.
	ID string `json:"id" api:"required"`
	// Inference speed mode. `fast` provides significantly faster output token
	// generation at premium pricing. Not all models support `fast`; invalid
	// combinations are rejected at create time.
	//
	// Any of "standard", "fast".
	Speed BetaDreamModelConfigSpeed `json:"speed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Speed       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDreamModelConfig) RawJSON() string { return r.JSON.raw }
func (r *BetaDreamModelConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inference speed mode. `fast` provides significantly faster output token
// generation at premium pricing. Not all models support `fast`; invalid
// combinations are rejected at create time.
type BetaDreamModelConfigSpeed string

const (
	BetaDreamModelConfigSpeedStandard BetaDreamModelConfigSpeed = "standard"
	BetaDreamModelConfigSpeedFast     BetaDreamModelConfigSpeed = "fast"
)

// Model identifier and configuration applied to every pipeline stage.
//
// The property ID is required.
type BetaDreamModelConfigParam struct {
	// Model identifier, e.g. "claude-opus-5". 1-256 characters.
	ID string `json:"id" api:"required"`
	// Inference speed mode. `fast` provides significantly faster output token
	// generation at premium pricing. Not all models support `fast`; invalid
	// combinations are rejected at create time.
	//
	// Any of "standard", "fast".
	Speed BetaDreamModelConfigParamSpeed `json:"speed,omitzero"`
	paramObj
}

func (r BetaDreamModelConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaDreamModelConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaDreamModelConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inference speed mode. `fast` provides significantly faster output token
// generation at premium pricing. Not all models support `fast`; invalid
// combinations are rejected at create time.
type BetaDreamModelConfigParamSpeed string

const (
	BetaDreamModelConfigParamSpeedStandard BetaDreamModelConfigParamSpeed = "standard"
	BetaDreamModelConfigParamSpeedFast     BetaDreamModelConfigParamSpeed = "fast"
)

// An output memory store the dream writes consolidated memories into.
type BetaDreamOutput struct {
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// Any of "memory_store".
	Type BetaDreamOutputType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MemoryStoreID respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDreamOutput) RawJSON() string { return r.JSON.raw }
func (r *BetaDreamOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaDreamOutputType string

const (
	BetaDreamOutputTypeMemoryStore BetaDreamOutputType = "memory_store"
)

// Input session transcripts the dream reads.
type BetaDreamSessionsInput struct {
	SessionIDs []string `json:"session_ids" api:"required"`
	// Any of "sessions".
	Type BetaDreamSessionsInputType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SessionIDs  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDreamSessionsInput) RawJSON() string { return r.JSON.raw }
func (r *BetaDreamSessionsInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaDreamSessionsInput to a BetaDreamSessionsInputParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaDreamSessionsInputParam.Overrides()
func (r BetaDreamSessionsInput) ToParam() BetaDreamSessionsInputParam {
	return param.Override[BetaDreamSessionsInputParam](json.RawMessage(r.RawJSON()))
}

type BetaDreamSessionsInputType string

const (
	BetaDreamSessionsInputTypeSessions BetaDreamSessionsInputType = "sessions"
)

// Input session transcripts the dream reads.
//
// The properties SessionIDs, Type are required.
type BetaDreamSessionsInputParam struct {
	SessionIDs []string `json:"session_ids,omitzero" api:"required"`
	// Any of "sessions".
	Type BetaDreamSessionsInputType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaDreamSessionsInputParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaDreamSessionsInputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaDreamSessionsInputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Lifecycle status of a Dream.
type BetaDreamStatus string

const (
	BetaDreamStatusPending   BetaDreamStatus = "pending"
	BetaDreamStatusRunning   BetaDreamStatus = "running"
	BetaDreamStatusCompleted BetaDreamStatus = "completed"
	BetaDreamStatusFailed    BetaDreamStatus = "failed"
	BetaDreamStatusCanceled  BetaDreamStatus = "canceled"
)

// Cumulative token usage for the dream across every pipeline stage.
type BetaDreamUsage struct {
	// Total tokens used to create prompt-cache entries (sum of all TTL tiers).
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens" api:"required"`
	// Total tokens read from prompt cache.
	CacheReadInputTokens int64 `json:"cache_read_input_tokens" api:"required"`
	// Total uncached input tokens consumed across every pipeline stage.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// Total output tokens generated across every pipeline stage.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CacheCreationInputTokens respjson.Field
		CacheReadInputTokens     respjson.Field
		InputTokens              respjson.Field
		OutputTokens             respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaDreamUsage) RawJSON() string { return r.JSON.raw }
func (r *BetaDreamUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaOutputBehaviorUnion contains all possible properties and values from
// [BetaOutputBehaviorCreateNew], [BetaOutputBehaviorUpdateExisting].
//
// Use the [BetaOutputBehaviorUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaOutputBehaviorUnion struct {
	// Any of "create_new", "update_existing".
	Type string `json:"type"`
	// This field is from variant [BetaOutputBehaviorUpdateExisting].
	MemoryStoreID string `json:"memory_store_id"`
	JSON          struct {
		Type          respjson.Field
		MemoryStoreID respjson.Field
		raw           string
	} `json:"-"`
}

// anyBetaOutputBehavior is implemented by each variant of
// [BetaOutputBehaviorUnion] to add type safety for the return type of
// [BetaOutputBehaviorUnion.AsAny]
type anyBetaOutputBehavior interface {
	implBetaOutputBehaviorUnion()
}

func (BetaOutputBehaviorCreateNew) implBetaOutputBehaviorUnion()      {}
func (BetaOutputBehaviorUpdateExisting) implBetaOutputBehaviorUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaOutputBehaviorUnion.AsAny().(type) {
//	case anthropic.BetaOutputBehaviorCreateNew:
//	case anthropic.BetaOutputBehaviorUpdateExisting:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaOutputBehaviorUnion) AsAny() anyBetaOutputBehavior {
	switch u.Type {
	case "create_new":
		return u.AsCreateNew()
	case "update_existing":
		return u.AsUpdateExisting()
	}
	return nil
}

func (u BetaOutputBehaviorUnion) AsCreateNew() (v BetaOutputBehaviorCreateNew) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaOutputBehaviorUnion) AsUpdateExisting() (v BetaOutputBehaviorUpdateExisting) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaOutputBehaviorUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaOutputBehaviorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaOutputBehaviorUnion to a BetaOutputBehaviorUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaOutputBehaviorUnionParam.Overrides()
func (r BetaOutputBehaviorUnion) ToParam() BetaOutputBehaviorUnionParam {
	return param.Override[BetaOutputBehaviorUnionParam](json.RawMessage(r.RawJSON()))
}

func BetaOutputBehaviorParamOfCreateNew(type_ BetaOutputBehaviorCreateNewType) BetaOutputBehaviorUnionParam {
	var createNew BetaOutputBehaviorCreateNewParam
	createNew.Type = type_
	return BetaOutputBehaviorUnionParam{OfCreateNew: &createNew}
}

func BetaOutputBehaviorParamOfUpdateExisting(memoryStoreID string) BetaOutputBehaviorUnionParam {
	var updateExisting BetaOutputBehaviorUpdateExistingParam
	updateExisting.MemoryStoreID = memoryStoreID
	return BetaOutputBehaviorUnionParam{OfUpdateExisting: &updateExisting}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaOutputBehaviorUnionParam struct {
	OfCreateNew      *BetaOutputBehaviorCreateNewParam      `json:",omitzero,inline"`
	OfUpdateExisting *BetaOutputBehaviorUpdateExistingParam `json:",omitzero,inline"`
	paramUnion
}

func (u BetaOutputBehaviorUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfCreateNew, u.OfUpdateExisting)
}
func (u *BetaOutputBehaviorUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaOutputBehaviorUnionParam) asAny() any {
	if !param.IsOmitted(u.OfCreateNew) {
		return u.OfCreateNew
	} else if !param.IsOmitted(u.OfUpdateExisting) {
		return u.OfUpdateExisting
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOutputBehaviorUnionParam) GetMemoryStoreID() *string {
	if vt := u.OfUpdateExisting; vt != nil {
		return &vt.MemoryStoreID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u BetaOutputBehaviorUnionParam) GetType() *string {
	if vt := u.OfCreateNew; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfUpdateExisting; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[BetaOutputBehaviorUnionParam](
		"type",
		apijson.Discriminator[BetaOutputBehaviorCreateNewParam]("create_new"),
		apijson.Discriminator[BetaOutputBehaviorUpdateExistingParam]("update_existing"),
	)
}

// The default destination: the job creates a new output memory store as a clone of
// the memory_store input and writes the consolidated memories into it. The input
// store is never mutated.
type BetaOutputBehaviorCreateNew struct {
	// Any of "create_new".
	Type BetaOutputBehaviorCreateNewType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOutputBehaviorCreateNew) RawJSON() string { return r.JSON.raw }
func (r *BetaOutputBehaviorCreateNew) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaOutputBehaviorCreateNew to a
// BetaOutputBehaviorCreateNewParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaOutputBehaviorCreateNewParam.Overrides()
func (r BetaOutputBehaviorCreateNew) ToParam() BetaOutputBehaviorCreateNewParam {
	return param.Override[BetaOutputBehaviorCreateNewParam](json.RawMessage(r.RawJSON()))
}

type BetaOutputBehaviorCreateNewType string

const (
	BetaOutputBehaviorCreateNewTypeCreateNew BetaOutputBehaviorCreateNewType = "create_new"
)

// The default destination: the job creates a new output memory store as a clone of
// the memory_store input and writes the consolidated memories into it. The input
// store is never mutated.
//
// The property Type is required.
type BetaOutputBehaviorCreateNewParam struct {
	// Any of "create_new".
	Type BetaOutputBehaviorCreateNewType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaOutputBehaviorCreateNewParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaOutputBehaviorCreateNewParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOutputBehaviorCreateNewParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The job writes the consolidated memories into this existing memory store instead
// of creating one. In EAP the store must be the job's own memory_store input, so
// the job consolidates the store in place.
type BetaOutputBehaviorUpdateExisting struct {
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// Any of "update_existing".
	Type BetaOutputBehaviorUpdateExistingType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MemoryStoreID respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOutputBehaviorUpdateExisting) RawJSON() string { return r.JSON.raw }
func (r *BetaOutputBehaviorUpdateExisting) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this BetaOutputBehaviorUpdateExisting to a
// BetaOutputBehaviorUpdateExistingParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// BetaOutputBehaviorUpdateExistingParam.Overrides()
func (r BetaOutputBehaviorUpdateExisting) ToParam() BetaOutputBehaviorUpdateExistingParam {
	return param.Override[BetaOutputBehaviorUpdateExistingParam](json.RawMessage(r.RawJSON()))
}

type BetaOutputBehaviorUpdateExistingType string

const (
	BetaOutputBehaviorUpdateExistingTypeUpdateExisting BetaOutputBehaviorUpdateExistingType = "update_existing"
)

// The job writes the consolidated memories into this existing memory store instead
// of creating one. In EAP the store must be the job's own memory_store input, so
// the job consolidates the store in place.
//
// The properties MemoryStoreID, Type are required.
type BetaOutputBehaviorUpdateExistingParam struct {
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// Any of "update_existing".
	Type BetaOutputBehaviorUpdateExistingType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r BetaOutputBehaviorUpdateExistingParam) MarshalJSON() (data []byte, err error) {
	type shadow BetaOutputBehaviorUpdateExistingParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOutputBehaviorUpdateExistingParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaDreamNewParams struct {
	Inputs []BetaDreamInputUnionParam `json:"inputs,omitzero" api:"required"`
	// Model identifier and configuration applied to every pipeline stage.
	Model        BetaDreamNewParamsModelUnion `json:"model,omitzero" api:"required"`
	Instructions param.Opt[string]            `json:"instructions,omitzero"`
	// The default destination: the job creates a new output memory store as a clone of
	// the memory_store input and writes the consolidated memories into it. The input
	// store is never mutated.
	OutputBehavior BetaOutputBehaviorUnionParam `json:"output_behavior,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

func (r BetaDreamNewParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaDreamNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaDreamNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type BetaDreamNewParamsModelUnion struct {
	OfString               param.Opt[string]          `json:",omitzero,inline"`
	OfBetaDreamModelConfig *BetaDreamModelConfigParam `json:",omitzero,inline"`
	paramUnion
}

func (u BetaDreamNewParamsModelUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfBetaDreamModelConfig)
}
func (u *BetaDreamNewParamsModelUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *BetaDreamNewParamsModelUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfBetaDreamModelConfig) {
		return u.OfBetaDreamModelConfig
	}
	return nil
}

type BetaDreamGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaDreamListParams struct {
	// Return dreams with `created_at` strictly after this timestamp (exclusive lower
	// bound, RFC 3339). Unset applies no lower bound.
	CreatedAtGt param.Opt[time.Time] `query:"created_at[gt],omitzero" format:"date-time" json:"-"`
	// Return dreams with `created_at` strictly before this timestamp (exclusive upper
	// bound, RFC 3339). Unset applies no upper bound.
	CreatedAtLt param.Opt[time.Time] `query:"created_at[lt],omitzero" format:"date-time" json:"-"`
	// Query parameter for include_archived
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Query parameter for limit
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Query parameter for page
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Filter by lifecycle status. Repeat the parameter to match any of multiple
	// statuses. Empty applies no status filter.
	Statuses []BetaDreamStatus `query:"statuses,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaDreamListParams]'s query parameters as `url.Values`.
func (r BetaDreamListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaDreamArchiveParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaDreamCancelParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
