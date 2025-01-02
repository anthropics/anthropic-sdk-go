// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/param"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/tidwall/gjson"
)

// BetaMessageBatchService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaMessageBatchService] method instead.
type BetaMessageBatchService struct {
	Options []option.RequestOption
}

// NewBetaMessageBatchService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewBetaMessageBatchService(opts ...option.RequestOption) (r *BetaMessageBatchService) {
	r = &BetaMessageBatchService{}
	r.Options = opts
	return
}

// Send a batch of Message creation requests.
//
// The Message Batches API can be used to process multiple Messages API requests at
// once. Once a Message Batch is created, it begins processing immediately. Batches
// can take up to 24 hours to complete.
func (r *BetaMessageBatchService) New(ctx context.Context, params BetaMessageBatchNewParams, opts ...option.RequestOption) (res *BetaMessageBatch, err error) {
	if params.Betas.Present {
		opts = append(opts, option.WithHeader("betas", fmt.Sprintf("%s", params.Betas)))
	}
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "message-batches-2024-09-24")}, opts...)
	path := "v1/messages/batches?beta=true"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return
}

// This endpoint is idempotent and can be used to poll for Message Batch
// completion. To access the results of a Message Batch, make a request to the
// `results_url` field in the response.
func (r *BetaMessageBatchService) Get(ctx context.Context, messageBatchID string, query BetaMessageBatchGetParams, opts ...option.RequestOption) (res *BetaMessageBatch, err error) {
	if query.Betas.Present {
		opts = append(opts, option.WithHeader("betas", fmt.Sprintf("%s", query.Betas)))
	}
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "message-batches-2024-09-24")}, opts...)
	if messageBatchID == "" {
		err = errors.New("missing required message_batch_id parameter")
		return
	}
	path := fmt.Sprintf("v1/messages/batches/%s?beta=true", messageBatchID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// List all Message Batches within a Workspace. Most recently created batches are
// returned first.
func (r *BetaMessageBatchService) List(ctx context.Context, params BetaMessageBatchListParams, opts ...option.RequestOption) (res *pagination.Page[BetaMessageBatch], err error) {
	var raw *http.Response
	if params.Betas.Present {
		opts = append(opts, option.WithHeader("betas", fmt.Sprintf("%s", params.Betas)))
	}
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "message-batches-2024-09-24"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/messages/batches?beta=true"
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

// List all Message Batches within a Workspace. Most recently created batches are
// returned first.
func (r *BetaMessageBatchService) ListAutoPaging(ctx context.Context, params BetaMessageBatchListParams, opts ...option.RequestOption) *pagination.PageAutoPager[BetaMessageBatch] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

// This endpoint is idempotent and can be used to poll for Message Batch
// completion. To access the results of a Message Batch, make a request to the
// `results_url` field in the response.
func (r *BetaMessageBatchService) Delete(ctx context.Context, messageBatchID string, body BetaMessageBatchDeleteParams, opts ...option.RequestOption) (res *BetaDeletedMessageBatch, err error) {
	if body.Betas.Present {
		opts = append(opts, option.WithHeader("betas", fmt.Sprintf("%s", body.Betas)))
	}
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "message-batches-2024-09-24")}, opts...)
	if messageBatchID == "" {
		err = errors.New("missing required message_batch_id parameter")
		return
	}
	path := fmt.Sprintf("v1/messages/batches/%s?beta=true", messageBatchID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Batches may be canceled any time before processing ends. Once cancellation is
// initiated, the batch enters a `canceling` state, at which time the system may
// complete any in-progress, non-interruptible requests before finalizing
// cancellation.
//
// The number of canceled requests is specified in `request_counts`. To determine
// which requests were canceled, check the individual results within the batch.
// Note that cancellation may not result in any canceled requests if they were
// non-interruptible.
func (r *BetaMessageBatchService) Cancel(ctx context.Context, messageBatchID string, body BetaMessageBatchCancelParams, opts ...option.RequestOption) (res *BetaMessageBatch, err error) {
	if body.Betas.Present {
		opts = append(opts, option.WithHeader("betas", fmt.Sprintf("%s", body.Betas)))
	}
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "message-batches-2024-09-24")}, opts...)
	if messageBatchID == "" {
		err = errors.New("missing required message_batch_id parameter")
		return
	}
	path := fmt.Sprintf("v1/messages/batches/%s/cancel?beta=true", messageBatchID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

// Streams the results of a Message Batch as a `.jsonl` file.
//
// Each line in the file is a JSON object containing the result of a single request
// in the Message Batch. Results are not guaranteed to be in the same order as
// requests. Use the `custom_id` field to match results to requests.
func (r *BetaMessageBatchService) Results(ctx context.Context, messageBatchID string, query BetaMessageBatchResultsParams, opts ...option.RequestOption) (res *http.Response, err error) {
	if query.Betas.Present {
		opts = append(opts, option.WithHeader("betas", fmt.Sprintf("%s", query.Betas)))
	}
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "message-batches-2024-09-24"), option.WithHeader("Accept", "application/binary")}, opts...)
	if messageBatchID == "" {
		err = errors.New("missing required message_batch_id parameter")
		return
	}
	path := fmt.Sprintf("v1/messages/batches/%s/results?beta=true", messageBatchID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

type BetaDeletedMessageBatch struct {
	// ID of the Message Batch.
	ID string `json:"id,required"`
	// Deleted object type.
	//
	// For Message Batches, this is always `"message_batch_deleted"`.
	Type BetaDeletedMessageBatchType `json:"type,required"`
	JSON betaDeletedMessageBatchJSON `json:"-"`
}

// betaDeletedMessageBatchJSON contains the JSON metadata for the struct
// [BetaDeletedMessageBatch]
type betaDeletedMessageBatchJSON struct {
	ID          apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaDeletedMessageBatch) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaDeletedMessageBatchJSON) RawJSON() string {
	return r.raw
}

// Deleted object type.
//
// For Message Batches, this is always `"message_batch_deleted"`.
type BetaDeletedMessageBatchType string

const (
	BetaDeletedMessageBatchTypeMessageBatchDeleted BetaDeletedMessageBatchType = "message_batch_deleted"
)

func (r BetaDeletedMessageBatchType) IsKnown() bool {
	switch r {
	case BetaDeletedMessageBatchTypeMessageBatchDeleted:
		return true
	}
	return false
}

type BetaMessageBatch struct {
	// Unique object identifier.
	//
	// The format and length of IDs may change over time.
	ID string `json:"id,required"`
	// RFC 3339 datetime string representing the time at which the Message Batch was
	// archived and its results became unavailable.
	ArchivedAt time.Time `json:"archived_at,required,nullable" format:"date-time"`
	// RFC 3339 datetime string representing the time at which cancellation was
	// initiated for the Message Batch. Specified only if cancellation was initiated.
	CancelInitiatedAt time.Time `json:"cancel_initiated_at,required,nullable" format:"date-time"`
	// RFC 3339 datetime string representing the time at which the Message Batch was
	// created.
	CreatedAt time.Time `json:"created_at,required" format:"date-time"`
	// RFC 3339 datetime string representing the time at which processing for the
	// Message Batch ended. Specified only once processing ends.
	//
	// Processing ends when every request in a Message Batch has either succeeded,
	// errored, canceled, or expired.
	EndedAt time.Time `json:"ended_at,required,nullable" format:"date-time"`
	// RFC 3339 datetime string representing the time at which the Message Batch will
	// expire and end processing, which is 24 hours after creation.
	ExpiresAt time.Time `json:"expires_at,required" format:"date-time"`
	// Processing status of the Message Batch.
	ProcessingStatus BetaMessageBatchProcessingStatus `json:"processing_status,required"`
	// Tallies requests within the Message Batch, categorized by their status.
	//
	// Requests start as `processing` and move to one of the other statuses only once
	// processing of the entire batch ends. The sum of all values always matches the
	// total number of requests in the batch.
	RequestCounts BetaMessageBatchRequestCounts `json:"request_counts,required"`
	// URL to a `.jsonl` file containing the results of the Message Batch requests.
	// Specified only once processing ends.
	//
	// Results in the file are not guaranteed to be in the same order as requests. Use
	// the `custom_id` field to match results to requests.
	ResultsURL string `json:"results_url,required,nullable"`
	// Object type.
	//
	// For Message Batches, this is always `"message_batch"`.
	Type BetaMessageBatchType `json:"type,required"`
	JSON betaMessageBatchJSON `json:"-"`
}

// betaMessageBatchJSON contains the JSON metadata for the struct
// [BetaMessageBatch]
type betaMessageBatchJSON struct {
	ID                apijson.Field
	ArchivedAt        apijson.Field
	CancelInitiatedAt apijson.Field
	CreatedAt         apijson.Field
	EndedAt           apijson.Field
	ExpiresAt         apijson.Field
	ProcessingStatus  apijson.Field
	RequestCounts     apijson.Field
	ResultsURL        apijson.Field
	Type              apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *BetaMessageBatch) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaMessageBatchJSON) RawJSON() string {
	return r.raw
}

// Processing status of the Message Batch.
type BetaMessageBatchProcessingStatus string

const (
	BetaMessageBatchProcessingStatusInProgress BetaMessageBatchProcessingStatus = "in_progress"
	BetaMessageBatchProcessingStatusCanceling  BetaMessageBatchProcessingStatus = "canceling"
	BetaMessageBatchProcessingStatusEnded      BetaMessageBatchProcessingStatus = "ended"
)

func (r BetaMessageBatchProcessingStatus) IsKnown() bool {
	switch r {
	case BetaMessageBatchProcessingStatusInProgress, BetaMessageBatchProcessingStatusCanceling, BetaMessageBatchProcessingStatusEnded:
		return true
	}
	return false
}

// Object type.
//
// For Message Batches, this is always `"message_batch"`.
type BetaMessageBatchType string

const (
	BetaMessageBatchTypeMessageBatch BetaMessageBatchType = "message_batch"
)

func (r BetaMessageBatchType) IsKnown() bool {
	switch r {
	case BetaMessageBatchTypeMessageBatch:
		return true
	}
	return false
}

type BetaMessageBatchCanceledResult struct {
	Type BetaMessageBatchCanceledResultType `json:"type,required"`
	JSON betaMessageBatchCanceledResultJSON `json:"-"`
}

// betaMessageBatchCanceledResultJSON contains the JSON metadata for the struct
// [BetaMessageBatchCanceledResult]
type betaMessageBatchCanceledResultJSON struct {
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaMessageBatchCanceledResult) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaMessageBatchCanceledResultJSON) RawJSON() string {
	return r.raw
}

func (r BetaMessageBatchCanceledResult) implementsBetaMessageBatchResult() {}

type BetaMessageBatchCanceledResultType string

const (
	BetaMessageBatchCanceledResultTypeCanceled BetaMessageBatchCanceledResultType = "canceled"
)

func (r BetaMessageBatchCanceledResultType) IsKnown() bool {
	switch r {
	case BetaMessageBatchCanceledResultTypeCanceled:
		return true
	}
	return false
}

type BetaMessageBatchErroredResult struct {
	Error BetaErrorResponse                 `json:"error,required"`
	Type  BetaMessageBatchErroredResultType `json:"type,required"`
	JSON  betaMessageBatchErroredResultJSON `json:"-"`
}

// betaMessageBatchErroredResultJSON contains the JSON metadata for the struct
// [BetaMessageBatchErroredResult]
type betaMessageBatchErroredResultJSON struct {
	Error       apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaMessageBatchErroredResult) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaMessageBatchErroredResultJSON) RawJSON() string {
	return r.raw
}

func (r BetaMessageBatchErroredResult) implementsBetaMessageBatchResult() {}

type BetaMessageBatchErroredResultType string

const (
	BetaMessageBatchErroredResultTypeErrored BetaMessageBatchErroredResultType = "errored"
)

func (r BetaMessageBatchErroredResultType) IsKnown() bool {
	switch r {
	case BetaMessageBatchErroredResultTypeErrored:
		return true
	}
	return false
}

type BetaMessageBatchExpiredResult struct {
	Type BetaMessageBatchExpiredResultType `json:"type,required"`
	JSON betaMessageBatchExpiredResultJSON `json:"-"`
}

// betaMessageBatchExpiredResultJSON contains the JSON metadata for the struct
// [BetaMessageBatchExpiredResult]
type betaMessageBatchExpiredResultJSON struct {
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaMessageBatchExpiredResult) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaMessageBatchExpiredResultJSON) RawJSON() string {
	return r.raw
}

func (r BetaMessageBatchExpiredResult) implementsBetaMessageBatchResult() {}

type BetaMessageBatchExpiredResultType string

const (
	BetaMessageBatchExpiredResultTypeExpired BetaMessageBatchExpiredResultType = "expired"
)

func (r BetaMessageBatchExpiredResultType) IsKnown() bool {
	switch r {
	case BetaMessageBatchExpiredResultTypeExpired:
		return true
	}
	return false
}

type BetaMessageBatchIndividualResponse struct {
	// Developer-provided ID created for each request in a Message Batch. Useful for
	// matching results to requests, as results may be given out of request order.
	//
	// Must be unique for each request within the Message Batch.
	CustomID string `json:"custom_id,required"`
	// Processing result for this request.
	//
	// Contains a Message output if processing was successful, an error response if
	// processing failed, or the reason why processing was not attempted, such as
	// cancellation or expiration.
	Result BetaMessageBatchResult                 `json:"result,required"`
	JSON   betaMessageBatchIndividualResponseJSON `json:"-"`
}

// betaMessageBatchIndividualResponseJSON contains the JSON metadata for the struct
// [BetaMessageBatchIndividualResponse]
type betaMessageBatchIndividualResponseJSON struct {
	CustomID    apijson.Field
	Result      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaMessageBatchIndividualResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaMessageBatchIndividualResponseJSON) RawJSON() string {
	return r.raw
}

type BetaMessageBatchRequestCounts struct {
	// Number of requests in the Message Batch that have been canceled.
	//
	// This is zero until processing of the entire Message Batch has ended.
	Canceled int64 `json:"canceled,required"`
	// Number of requests in the Message Batch that encountered an error.
	//
	// This is zero until processing of the entire Message Batch has ended.
	Errored int64 `json:"errored,required"`
	// Number of requests in the Message Batch that have expired.
	//
	// This is zero until processing of the entire Message Batch has ended.
	Expired int64 `json:"expired,required"`
	// Number of requests in the Message Batch that are processing.
	Processing int64 `json:"processing,required"`
	// Number of requests in the Message Batch that have completed successfully.
	//
	// This is zero until processing of the entire Message Batch has ended.
	Succeeded int64                             `json:"succeeded,required"`
	JSON      betaMessageBatchRequestCountsJSON `json:"-"`
}

// betaMessageBatchRequestCountsJSON contains the JSON metadata for the struct
// [BetaMessageBatchRequestCounts]
type betaMessageBatchRequestCountsJSON struct {
	Canceled    apijson.Field
	Errored     apijson.Field
	Expired     apijson.Field
	Processing  apijson.Field
	Succeeded   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaMessageBatchRequestCounts) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaMessageBatchRequestCountsJSON) RawJSON() string {
	return r.raw
}

// Processing result for this request.
//
// Contains a Message output if processing was successful, an error response if
// processing failed, or the reason why processing was not attempted, such as
// cancellation or expiration.
type BetaMessageBatchResult struct {
	Type    BetaMessageBatchResultType `json:"type,required"`
	Error   BetaErrorResponse          `json:"error"`
	Message BetaMessage                `json:"message"`
	JSON    betaMessageBatchResultJSON `json:"-"`
	union   BetaMessageBatchResultUnion
}

// betaMessageBatchResultJSON contains the JSON metadata for the struct
// [BetaMessageBatchResult]
type betaMessageBatchResultJSON struct {
	Type        apijson.Field
	Error       apijson.Field
	Message     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r betaMessageBatchResultJSON) RawJSON() string {
	return r.raw
}

func (r *BetaMessageBatchResult) UnmarshalJSON(data []byte) (err error) {
	*r = BetaMessageBatchResult{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [BetaMessageBatchResultUnion] interface which you can cast to
// the specific types for more type safety.
//
// Possible runtime types of the union are [BetaMessageBatchSucceededResult],
// [BetaMessageBatchErroredResult], [BetaMessageBatchCanceledResult],
// [BetaMessageBatchExpiredResult].
func (r BetaMessageBatchResult) AsUnion() BetaMessageBatchResultUnion {
	return r.union
}

// Processing result for this request.
//
// Contains a Message output if processing was successful, an error response if
// processing failed, or the reason why processing was not attempted, such as
// cancellation or expiration.
//
// Union satisfied by [BetaMessageBatchSucceededResult],
// [BetaMessageBatchErroredResult], [BetaMessageBatchCanceledResult] or
// [BetaMessageBatchExpiredResult].
type BetaMessageBatchResultUnion interface {
	implementsBetaMessageBatchResult()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*BetaMessageBatchResultUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaMessageBatchSucceededResult{}),
			DiscriminatorValue: "succeeded",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaMessageBatchErroredResult{}),
			DiscriminatorValue: "errored",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaMessageBatchCanceledResult{}),
			DiscriminatorValue: "canceled",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(BetaMessageBatchExpiredResult{}),
			DiscriminatorValue: "expired",
		},
	)
}

type BetaMessageBatchResultType string

const (
	BetaMessageBatchResultTypeSucceeded BetaMessageBatchResultType = "succeeded"
	BetaMessageBatchResultTypeErrored   BetaMessageBatchResultType = "errored"
	BetaMessageBatchResultTypeCanceled  BetaMessageBatchResultType = "canceled"
	BetaMessageBatchResultTypeExpired   BetaMessageBatchResultType = "expired"
)

func (r BetaMessageBatchResultType) IsKnown() bool {
	switch r {
	case BetaMessageBatchResultTypeSucceeded, BetaMessageBatchResultTypeErrored, BetaMessageBatchResultTypeCanceled, BetaMessageBatchResultTypeExpired:
		return true
	}
	return false
}

type BetaMessageBatchSucceededResult struct {
	Message BetaMessage                         `json:"message,required"`
	Type    BetaMessageBatchSucceededResultType `json:"type,required"`
	JSON    betaMessageBatchSucceededResultJSON `json:"-"`
}

// betaMessageBatchSucceededResultJSON contains the JSON metadata for the struct
// [BetaMessageBatchSucceededResult]
type betaMessageBatchSucceededResultJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BetaMessageBatchSucceededResult) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r betaMessageBatchSucceededResultJSON) RawJSON() string {
	return r.raw
}

func (r BetaMessageBatchSucceededResult) implementsBetaMessageBatchResult() {}

type BetaMessageBatchSucceededResultType string

const (
	BetaMessageBatchSucceededResultTypeSucceeded BetaMessageBatchSucceededResultType = "succeeded"
)

func (r BetaMessageBatchSucceededResultType) IsKnown() bool {
	switch r {
	case BetaMessageBatchSucceededResultTypeSucceeded:
		return true
	}
	return false
}

type BetaMessageBatchNewParams struct {
	// List of requests for prompt completion. Each is an individual request to create
	// a Message.
	Requests param.Field[[]BetaMessageBatchNewParamsRequest] `json:"requests,required"`
	// Optional header to specify the beta version(s) you want to use.
	Betas param.Field[[]AnthropicBeta] `header:"anthropic-beta"`
}

func (r BetaMessageBatchNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaMessageBatchNewParamsRequest struct {
	// Developer-provided ID created for each request in a Message Batch. Useful for
	// matching results to requests, as results may be given out of request order.
	//
	// Must be unique for each request within the Message Batch.
	CustomID param.Field[string] `json:"custom_id,required"`
	// Messages API creation parameters for the individual request.
	//
	// See the [Messages API reference](/en/api/messages) for full documentation on
	// available parameters.
	Params param.Field[BetaMessageBatchNewParamsRequestsParams] `json:"params,required"`
}

func (r BetaMessageBatchNewParamsRequest) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Messages API creation parameters for the individual request.
//
// See the [Messages API reference](/en/api/messages) for full documentation on
// available parameters.
type BetaMessageBatchNewParamsRequestsParams struct {
	// The maximum number of tokens to generate before stopping.
	//
	// Note that our models may stop _before_ reaching this maximum. This parameter
	// only specifies the absolute maximum number of tokens to generate.
	//
	// Different models have different maximum values for this parameter. See
	// [models](https://docs.anthropic.com/en/docs/models-overview) for details.
	MaxTokens param.Field[int64] `json:"max_tokens,required"`
	// Input messages.
	//
	// Our models are trained to operate on alternating `user` and `assistant`
	// conversational turns. When creating a new `Message`, you specify the prior
	// conversational turns with the `messages` parameter, and the model then generates
	// the next `Message` in the conversation. Consecutive `user` or `assistant` turns
	// in your request will be combined into a single turn.
	//
	// Each input message must be an object with a `role` and `content`. You can
	// specify a single `user`-role message, or you can include multiple `user` and
	// `assistant` messages.
	//
	// If the final message uses the `assistant` role, the response content will
	// continue immediately from the content in that message. This can be used to
	// constrain part of the model's response.
	//
	// Example with a single `user` message:
	//
	// ```json
	// [{ "role": "user", "content": "Hello, Claude" }]
	// ```
	//
	// Example with multiple conversational turns:
	//
	// ```json
	// [
	//
	//	{ "role": "user", "content": "Hello there." },
	//	{ "role": "assistant", "content": "Hi, I'm Claude. How can I help you?" },
	//	{ "role": "user", "content": "Can you explain LLMs in plain English?" }
	//
	// ]
	// ```
	//
	// Example with a partially-filled response from Claude:
	//
	// ```json
	// [
	//
	//	{
	//	  "role": "user",
	//	  "content": "What's the Greek name for Sun? (A) Sol (B) Helios (C) Sun"
	//	},
	//	{ "role": "assistant", "content": "The best answer is (" }
	//
	// ]
	// ```
	//
	// Each input message `content` may be either a single `string` or an array of
	// content blocks, where each block has a specific `type`. Using a `string` for
	// `content` is shorthand for an array of one content block of type `"text"`. The
	// following input messages are equivalent:
	//
	// ```json
	// { "role": "user", "content": "Hello, Claude" }
	// ```
	//
	// ```json
	// { "role": "user", "content": [{ "type": "text", "text": "Hello, Claude" }] }
	// ```
	//
	// Starting with Claude 3 models, you can also send image content blocks:
	//
	// ```json
	//
	//	{
	//	  "role": "user",
	//	  "content": [
	//	    {
	//	      "type": "image",
	//	      "source": {
	//	        "type": "base64",
	//	        "media_type": "image/jpeg",
	//	        "data": "/9j/4AAQSkZJRg..."
	//	      }
	//	    },
	//	    { "type": "text", "text": "What is in this image?" }
	//	  ]
	//	}
	//
	// ```
	//
	// We currently support the `base64` source type for images, and the `image/jpeg`,
	// `image/png`, `image/gif`, and `image/webp` media types.
	//
	// See [examples](https://docs.anthropic.com/en/api/messages-examples#vision) for
	// more input examples.
	//
	// Note that if you want to include a
	// [system prompt](https://docs.anthropic.com/en/docs/system-prompts), you can use
	// the top-level `system` parameter — there is no `"system"` role for input
	// messages in the Messages API.
	Messages param.Field[[]BetaMessageParam] `json:"messages,required"`
	// The model that will complete your prompt.\n\nSee
	// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	Model param.Field[Model] `json:"model,required"`
	// An object describing metadata about the request.
	Metadata param.Field[BetaMetadataParam] `json:"metadata"`
	// Custom text sequences that will cause the model to stop generating.
	//
	// Our models will normally stop when they have naturally completed their turn,
	// which will result in a response `stop_reason` of `"end_turn"`.
	//
	// If you want the model to stop generating when it encounters custom strings of
	// text, you can use the `stop_sequences` parameter. If the model encounters one of
	// the custom sequences, the response `stop_reason` value will be `"stop_sequence"`
	// and the response `stop_sequence` value will contain the matched stop sequence.
	StopSequences param.Field[[]string] `json:"stop_sequences"`
	// Whether to incrementally stream the response using server-sent events.
	//
	// See [streaming](https://docs.anthropic.com/en/api/messages-streaming) for
	// details.
	Stream param.Field[bool] `json:"stream"`
	// System prompt.
	//
	// A system prompt is a way of providing context and instructions to Claude, such
	// as specifying a particular goal or role. See our
	// [guide to system prompts](https://docs.anthropic.com/en/docs/system-prompts).
	System param.Field[[]BetaTextBlockParam] `json:"system"`
	// Amount of randomness injected into the response.
	//
	// Defaults to `1.0`. Ranges from `0.0` to `1.0`. Use `temperature` closer to `0.0`
	// for analytical / multiple choice, and closer to `1.0` for creative and
	// generative tasks.
	//
	// Note that even with `temperature` of `0.0`, the results will not be fully
	// deterministic.
	Temperature param.Field[float64] `json:"temperature"`
	// How the model should use the provided tools. The model can use a specific tool,
	// any available tool, or decide by itself.
	ToolChoice param.Field[BetaToolChoiceUnionParam] `json:"tool_choice"`
	// Definitions of tools that the model may use.
	//
	// If you include `tools` in your API request, the model may return `tool_use`
	// content blocks that represent the model's use of those tools. You can then run
	// those tools using the tool input generated by the model and then optionally
	// return results back to the model using `tool_result` content blocks.
	//
	// Each tool definition includes:
	//
	//   - `name`: Name of the tool.
	//   - `description`: Optional, but strongly-recommended description of the tool.
	//   - `input_schema`: [JSON schema](https://json-schema.org/) for the tool `input`
	//     shape that the model will produce in `tool_use` output content blocks.
	//
	// For example, if you defined `tools` as:
	//
	// ```json
	// [
	//
	//	{
	//	  "name": "get_stock_price",
	//	  "description": "Get the current stock price for a given ticker symbol.",
	//	  "input_schema": {
	//	    "type": "object",
	//	    "properties": {
	//	      "ticker": {
	//	        "type": "string",
	//	        "description": "The stock ticker symbol, e.g. AAPL for Apple Inc."
	//	      }
	//	    },
	//	    "required": ["ticker"]
	//	  }
	//	}
	//
	// ]
	// ```
	//
	// And then asked the model "What's the S&P 500 at today?", the model might produce
	// `tool_use` content blocks in the response like this:
	//
	// ```json
	// [
	//
	//	{
	//	  "type": "tool_use",
	//	  "id": "toolu_01D7FLrfh4GYq7yT1ULFeyMV",
	//	  "name": "get_stock_price",
	//	  "input": { "ticker": "^GSPC" }
	//	}
	//
	// ]
	// ```
	//
	// You might then run your `get_stock_price` tool with `{"ticker": "^GSPC"}` as an
	// input, and return the following back to the model in a subsequent `user`
	// message:
	//
	// ```json
	// [
	//
	//	{
	//	  "type": "tool_result",
	//	  "tool_use_id": "toolu_01D7FLrfh4GYq7yT1ULFeyMV",
	//	  "content": "259.75 USD"
	//	}
	//
	// ]
	// ```
	//
	// Tools can be used for workflows that include running client-side tools and
	// functions, or more generally whenever you want the model to produce a particular
	// JSON structure of output.
	//
	// See our [guide](https://docs.anthropic.com/en/docs/tool-use) for more details.
	Tools param.Field[[]BetaToolUnionUnionParam] `json:"tools"`
	// Only sample from the top K options for each subsequent token.
	//
	// Used to remove "long tail" low probability responses.
	// [Learn more technical details here](https://towardsdatascience.com/how-to-sample-from-language-models-682bceb97277).
	//
	// Recommended for advanced use cases only. You usually only need to use
	// `temperature`.
	TopK param.Field[int64] `json:"top_k"`
	// Use nucleus sampling.
	//
	// In nucleus sampling, we compute the cumulative distribution over all the options
	// for each subsequent token in decreasing probability order and cut it off once it
	// reaches a particular probability specified by `top_p`. You should either alter
	// `temperature` or `top_p`, but not both.
	//
	// Recommended for advanced use cases only. You usually only need to use
	// `temperature`.
	TopP param.Field[float64] `json:"top_p"`
}

func (r BetaMessageBatchNewParamsRequestsParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type BetaMessageBatchGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas param.Field[[]AnthropicBeta] `header:"anthropic-beta"`
}

type BetaMessageBatchListParams struct {
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately after this object.
	AfterID param.Field[string] `query:"after_id"`
	// ID of the object to use as a cursor for pagination. When provided, returns the
	// page of results immediately before this object.
	BeforeID param.Field[string] `query:"before_id"`
	// Number of items to return per page.
	//
	// Defaults to `20`. Ranges from `1` to `1000`.
	Limit param.Field[int64] `query:"limit"`
	// Optional header to specify the beta version(s) you want to use.
	Betas param.Field[[]AnthropicBeta] `header:"anthropic-beta"`
}

// URLQuery serializes [BetaMessageBatchListParams]'s query parameters as
// `url.Values`.
func (r BetaMessageBatchListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaMessageBatchDeleteParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas param.Field[[]AnthropicBeta] `header:"anthropic-beta"`
}

type BetaMessageBatchCancelParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas param.Field[[]AnthropicBeta] `header:"anthropic-beta"`
}

type BetaMessageBatchResultsParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas param.Field[[]AnthropicBeta] `header:"anthropic-beta"`
}
