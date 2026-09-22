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

// Start an asynchronous job that uses past sessions to produce a reorganized
// version of a memory store and get back the dream to poll for the result.
//
// By default the dream writes its result to a new memory store and doesn't change
// the input memory store. The response has `status` set to `pending` and an empty
// `outputs` array. Poll the dream until `status` is `completed`, `failed`, or
// `canceled`.
//
// See the
// [Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#create-a-dream)
// to learn more about creating dreams.
func (r *BetaDreamService) New(ctx context.Context, params BetaDreamNewParams, opts ...option.RequestOption) (res *BetaDream, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	if !param.IsOmitted(params.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", params.WorkspaceID.Value)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "dreaming-2026-04-21")}, opts...)
	path := "v1/dreams?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get a dream by ID to check its status, output memory store, and token usage.
//
// Archived dreams are returned too.
//
// See the
// [Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#track-progress)
// for how to poll a dream and what each status means.
func (r *BetaDreamService) Get(ctx context.Context, dreamID string, query BetaDreamGetParams, opts ...option.RequestOption) (res *BetaDream, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	if !param.IsOmitted(query.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", query.WorkspaceID.Value)))
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

// List the dreams in the workspace, newest first.
//
// Archived dreams are left out unless `include_archived` is `true`.
//
// See the
// [Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#list-dreams)
// for how to page through dreams.
func (r *BetaDreamService) List(ctx context.Context, params BetaDreamListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaDream], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	if !param.IsOmitted(params.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", params.WorkspaceID.Value)))
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

// List the dreams in the workspace, newest first.
//
// Archived dreams are left out unless `include_archived` is `true`.
//
// See the
// [Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#list-dreams)
// for how to page through dreams.
func (r *BetaDreamService) ListAutoPaging(ctx context.Context, params BetaDreamListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaDream] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Hide a `completed`, `failed`, or `canceled` dream from the default list of
// dreams.
//
// Archiving a `pending` or `running` dream returns a 400 error, so cancel it
// first. Archiving an archived dream returns it unchanged. An archived dream can
// still be fetched by ID. Archiving can't be undone.
//
// See the
// [Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#archive-a-dream)
// to learn more about archiving dreams.
func (r *BetaDreamService) Archive(ctx context.Context, dreamID string, body BetaDreamArchiveParams, opts ...option.RequestOption) (res *BetaDream, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	if !param.IsOmitted(body.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", body.WorkspaceID.Value)))
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

// Stop a `pending` or `running` dream.
//
// The response shows `status` as `canceled`, unless the dream reached `completed`
// or `failed` first. `usage` can keep changing after the response. Canceling a
// `canceled` dream returns it unchanged. Canceling a `completed` or `failed` dream
// returns a 400 error.
//
// See the
// [Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#cancel-a-dream)
// to learn more about canceling dreams.
func (r *BetaDreamService) Cancel(ctx context.Context, dreamID string, body BetaDreamCancelParams, opts ...option.RequestOption) (res *BetaDream, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("anthropic-beta", fmt.Sprintf("%v", v)))
	}
	if !param.IsOmitted(body.WorkspaceID) {
		opts = append(opts, option.WithHeader("anthropic-workspace-id", fmt.Sprintf("%v", body.WorkspaceID.Value)))
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

// An asynchronous job that reads a memory store and past sessions, then writes a
// reorganized version of that memory store.
//
// By default the dream writes its result to a new memory store and doesn't change
// the input memory store. With `output_behavior` set to `update_existing`, it
// writes its result into the input memory store instead. The Dreams API is in
// research preview, so this resource can still change.
//
// See the
// [Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#how-it-works)
// for what a dream reads and produces.
type BetaDream struct {
	// The unique ID of the dream (`drm_...`).
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	EndedAt time.Time `json:"ended_at" api:"required" format:"date-time"`
	// Failure detail for a Dream whose `status` is `failed`.
	Error BetaDreamError `json:"error" api:"required"`
	// The sources that the dream reads, from the request that created it.
	Inputs []BetaDreamInputUnion `json:"inputs" api:"required"`
	// The guidance given when the dream was created, or `null` if none was given.
	Instructions string `json:"instructions" api:"required"`
	// The model that runs a dream, from the request that created it.
	//
	// The dream uses this model for all of its work. The response always gives the
	// model as an object, even if the request gave only a model ID.
	Model BetaDreamModelConfig `json:"model" api:"required"`
	// Which memory store a dream writes its result to. Defaults to `create_new` when
	// left out of a create request.
	OutputBehavior BetaOutputBehaviorUnion `json:"output_behavior" api:"required"`
	// The memory store that holds the dream's result, as a one-item array, or an empty
	// array until the dream records that memory store.
	//
	// The array is empty while the dream is `pending` and for a short time after it
	// starts `running`. It can stay empty if the dream fails or is canceled before
	// then. The memory store holds the complete result only once `status` is
	// `completed`.
	//
	// See the
	// [Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#use-the-output)
	// for how to review and use the result.
	Outputs []BetaDreamOutput `json:"outputs" api:"required"`
	// The ID of the session that runs the dream (`sesn_...`), or `null` if that
	// session hasn't started.
	//
	// Stream that session's events to follow what the dream reads and writes.
	//
	// See the
	// [Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#watch-the-pipeline-run)
	// for how to watch a running dream.
	SessionID string `json:"session_id" api:"required"`
	// Where a dream is in its lifecycle.
	//
	// `completed`, `failed`, and `canceled` are final: once a dream has one of these
	// statuses, its status doesn't change again.
	//
	// See the
	// [Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#lifecycle)
	// for what each status means.
	//
	// Any of "pending", "running", "completed", "failed", "canceled".
	Status BetaDreamStatus `json:"status" api:"required"`
	// Any of "dream".
	Type BetaDreamType `json:"type" api:"required"`
	// The tokens that a dream has used so far.
	//
	// The counts are zero while the dream is `pending` and update while it is
	// `running`. They can keep changing after a cancel.
	//
	// See the
	// [Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#billing)
	// for how dreams are billed. See the
	// [prompt caching guide](https://platform.claude.com/docs/en/build-with-claude/prompt-caching#tracking-cache-performance)
	// for how the input token counts add up.
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
	// A human-readable explanation of why the dream failed.
	Message string `json:"message" api:"required"`
	// A code for why the dream failed, such as `timeout` or `internal_error`.
	//
	// The
	// [Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#errors)
	// lists common error codes and when they occur.
	Type string `json:"type" api:"required"`
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

// The memory store that a dream reads, given as an entry in `inputs`.
//
// With `output_behavior` set to `update_existing`, the dream writes its result
// into this memory store. Otherwise the dream doesn't change it.
type BetaDreamMemoryStoreInput struct {
	// The ID of the memory store for the dream to read (`memstore_...`).
	//
	// The memory store must be in the same workspace as the dream and must not be
	// archived.
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

// The memory store that a dream reads, given as an entry in `inputs`.
//
// With `output_behavior` set to `update_existing`, the dream writes its result
// into this memory store. Otherwise the dream doesn't change it.
//
// The properties MemoryStoreID, Type are required.
type BetaDreamMemoryStoreInputParam struct {
	// The ID of the memory store for the dream to read (`memstore_...`).
	//
	// The memory store must be in the same workspace as the dream and must not be
	// archived.
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

// The model that runs a dream, from the request that created it.
//
// The dream uses this model for all of its work. The response always gives the
// model as an object, even if the request gave only a model ID.
type BetaDreamModelConfig struct {
	// The ID of the model that runs the dream, as given in the request that created
	// it.
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

// The object form of `model` in a request to create a dream.
//
// The property ID is required.
type BetaDreamModelConfigParam struct {
	// The ID of the model to run the dream with.
	//
	// The ID can be 1 to 256 characters long.
	//
	// The
	// [limits table in the Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#limits)
	// lists the supported models.
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

// The memory store that holds a dream's result, as an entry in `outputs`.
type BetaDreamOutput struct {
	// The ID of the memory store that the dream writes its result to (`memstore_...`).
	//
	// With `output_behavior` set to `create_new`, this is a new memory store. With
	// `update_existing`, it is the input memory store.
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

// The sessions that a dream reads, given as an entry in `inputs`.
type BetaDreamSessionsInput struct {
	// The IDs of the sessions whose transcripts the dream reads (`sesn_...`).
	//
	// Give 1 to 100 IDs, with no duplicates. Each session must be in the same
	// workspace as the dream. Responses list the IDs in sorted order.
	//
	// The
	// [limits table in the Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#limits)
	// lists all the limits on a dream.
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

// The sessions that a dream reads, given as an entry in `inputs`.
//
// The properties SessionIDs, Type are required.
type BetaDreamSessionsInputParam struct {
	// The IDs of the sessions whose transcripts the dream reads (`sesn_...`).
	//
	// Give 1 to 100 IDs, with no duplicates. Each session must be in the same
	// workspace as the dream. Responses list the IDs in sorted order.
	//
	// The
	// [limits table in the Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#limits)
	// lists all the limits on a dream.
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

// Where a dream is in its lifecycle.
//
// `completed`, `failed`, and `canceled` are final: once a dream has one of these
// statuses, its status doesn't change again.
//
// See the
// [Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#lifecycle)
// for what each status means.
type BetaDreamStatus string

const (
	// The dream is waiting to start and hasn't read its inputs yet.
	//
	// `outputs` is empty and every `usage` count is zero.
	BetaDreamStatusPending BetaDreamStatus = "pending"
	// The dream is reading its inputs and writing its result.
	//
	// `usage` updates while the dream has this status.
	BetaDreamStatusRunning BetaDreamStatus = "running"
	// The dream finished and its output memory store holds the complete result.
	BetaDreamStatusCompleted BetaDreamStatus = "completed"
	// The dream stopped with an error, which `error` describes.
	//
	// If `outputs` references a memory store, that memory store keeps what the dream
	// wrote before it stopped.
	BetaDreamStatusFailed BetaDreamStatus = "failed"
	// A cancel request stopped the dream before it reached `completed` or `failed`.
	//
	// If `outputs` references a memory store, that memory store keeps what the dream
	// wrote. `usage` can keep changing after the cancel.
	BetaDreamStatusCanceled BetaDreamStatus = "canceled"
)

// The tokens that a dream has used so far.
//
// The counts are zero while the dream is `pending` and update while it is
// `running`. They can keep changing after a cancel.
//
// See the
// [Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#billing)
// for how dreams are billed. See the
// [prompt caching guide](https://platform.claude.com/docs/en/build-with-claude/prompt-caching#tracking-cache-performance)
// for how the input token counts add up.
type BetaDreamUsage struct {
	// The dream's input tokens that were written to the prompt cache, for both the
	// 5-minute and 1-hour cache durations.
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens" api:"required"`
	// The dream's input tokens that were read from the prompt cache.
	CacheReadInputTokens int64 `json:"cache_read_input_tokens" api:"required"`
	// The dream's input tokens that weren't read from or written to the prompt cache.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// The tokens that the model generated for the dream.
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

// Write the result to a new memory store that starts as a copy of the input memory
// store. This is the default.
//
// The new memory store is in the same workspace as the dream. The dream doesn't
// change the input memory store.
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

// Write the result to a new memory store that starts as a copy of the input memory
// store. This is the default.
//
// The new memory store is in the same workspace as the dream. The dream doesn't
// change the input memory store.
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

// Write the result into the input memory store instead of a new memory store.
//
// The credential must be allowed to write memory stores, or the request returns a
// 403 error. While another `update_existing` dream on the same memory store hasn't
// fully stopped, the request returns a 409 error.
type BetaOutputBehaviorUpdateExisting struct {
	// The ID of the memory store for the dream to write its result to
	// (`memstore_...`). It must be the memory store in the `memory_store` entry of
	// `inputs`.
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

// Write the result into the input memory store instead of a new memory store.
//
// The credential must be allowed to write memory stores, or the request returns a
// 403 error. While another `update_existing` dream on the same memory store hasn't
// fully stopped, the request returns a 409 error.
//
// The properties MemoryStoreID, Type are required.
type BetaOutputBehaviorUpdateExistingParam struct {
	// The ID of the memory store for the dream to write its result to
	// (`memstore_...`). It must be the memory store in the `memory_store` entry of
	// `inputs`.
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
	// The memory store and sessions for the dream to read, as exactly one
	// `memory_store` entry and exactly one `sessions` entry.
	Inputs []BetaDreamInputUnionParam `json:"inputs,omitzero" api:"required"`
	// The model that runs a dream, given as a model ID or as an object with `id` and
	// `speed`.
	//
	// In the object form, `speed` can only be `standard`.
	//
	// The
	// [limits table in the Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#limits)
	// lists the supported models.
	Model BetaDreamNewParamsModelUnion `json:"model,omitzero" api:"required"`
	// Guidance that steers how the dream reads the sessions and organizes the output
	// memory store, from 1 to 4,096 characters.
	//
	// See the
	// [Dreams guide](https://platform.claude.com/docs/en/managed-agents/dreams#steer-with-instructions)
	// for what kinds of instructions work well.
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// Which memory store a dream writes its result to. Defaults to `create_new` when
	// left out of a create request.
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
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaDreamListParams struct {
	// Return only dreams created after this time (exclusive), in RFC 3339.
	CreatedAtGt param.Opt[time.Time] `query:"created_at[gt],omitzero" format:"date-time" json:"-"`
	// Return only dreams created before this time (exclusive), in RFC 3339.
	CreatedAtLt param.Opt[time.Time] `query:"created_at[lt],omitzero" format:"date-time" json:"-"`
	// Whether to include archived dreams. Defaults to `false`.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// The maximum number of dreams to return, from 1 to 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// The cursor for the page to return, taken from `next_page` in a previous
	// response.
	//
	// Leave it out to get the first page.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// Return only dreams that have one of these statuses.
	//
	// Repeat the parameter to give more than one status. Leave it out to return dreams
	// of every status.
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
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

type BetaDreamCancelParams struct {
	// Optional header to select the Workspace for this request. The value is a
	// Workspace ID (for example, `wrkspc_011CZkZaBF1tNoB5wlCeusgy`).
	//
	// Only needed for credentials that can act on more than one Workspace. A
	// credential that belongs to a specific Workspace may omit it; if sent, it must
	// match that Workspace.
	WorkspaceID param.Opt[string] `header:"anthropic-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}
