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
	"github.com/anthropics/anthropic-sdk-go/packages/jsonl"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/shared"
	"github.com/tidwall/gjson"
)

// MessageBatchService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewMessageBatchService] method instead.
type MessageBatchService struct {
	Options []option.RequestOption
}

// NewMessageBatchService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewMessageBatchService(opts ...option.RequestOption) (r *MessageBatchService) {
	r = &MessageBatchService{}
	r.Options = opts
	return
}

// Send a batch of Message creation requests.
//
// The Message Batches API can be used to process multiple Messages API requests at
// once. Once a Message Batch is created, it begins processing immediately. Batches
// can take up to 24 hours to complete.
//
// Learn more about the Message Batches API in our
// [user guide](/en/docs/build-with-claude/batch-processing)
func (r *MessageBatchService) New(ctx context.Context, body MessageBatchNewParams, opts ...option.RequestOption) (res *MessageBatch, err error) {
	opts = append(r.Options[:], opts...)
	path := "v1/messages/batches"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// This endpoint is idempotent and can be used to poll for Message Batch
// completion. To access the results of a Message Batch, make a request to the
// `results_url` field in the response.
//
// Learn more about the Message Batches API in our
// [user guide](/en/docs/build-with-claude/batch-processing)
func (r *MessageBatchService) Get(ctx context.Context, messageBatchID string, opts ...option.RequestOption) (res *MessageBatch, err error) {
	opts = append(r.Options[:], opts...)
	if messageBatchID == "" {
		err = errors.New("missing required message_batch_id parameter")
		return
	}
	path := fmt.Sprintf("v1/messages/batches/%s", messageBatchID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// List all Message Batches within a Workspace. Most recently created batches are
// returned first.
//
// Learn more about the Message Batches API in our
// [user guide](/en/docs/build-with-claude/batch-processing)
func (r *MessageBatchService) List(ctx context.Context, query MessageBatchListParams, opts ...option.RequestOption) (res *pagination.Page[MessageBatch], err error) {
	var raw *http.Response
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/messages/batches"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
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
//
// Learn more about the Message Batches API in our
// [user guide](/en/docs/build-with-claude/batch-processing)
func (r *MessageBatchService) ListAutoPaging(ctx context.Context, query MessageBatchListParams, opts ...option.RequestOption) *pagination.PageAutoPager[MessageBatch] {
	return pagination.NewPageAutoPager(r.List(ctx, query, opts...))
}

// Delete a Message Batch.
//
// Message Batches can only be deleted once they've finished processing. If you'd
// like to delete an in-progress batch, you must first cancel it.
//
// Learn more about the Message Batches API in our
// [user guide](/en/docs/build-with-claude/batch-processing)
func (r *MessageBatchService) Delete(ctx context.Context, messageBatchID string, opts ...option.RequestOption) (res *DeletedMessageBatch, err error) {
	opts = append(r.Options[:], opts...)
	if messageBatchID == "" {
		err = errors.New("missing required message_batch_id parameter")
		return
	}
	path := fmt.Sprintf("v1/messages/batches/%s", messageBatchID)
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
//
// Learn more about the Message Batches API in our
// [user guide](/en/docs/build-with-claude/batch-processing)
func (r *MessageBatchService) Cancel(ctx context.Context, messageBatchID string, opts ...option.RequestOption) (res *MessageBatch, err error) {
	opts = append(r.Options[:], opts...)
	if messageBatchID == "" {
		err = errors.New("missing required message_batch_id parameter")
		return
	}
	path := fmt.Sprintf("v1/messages/batches/%s/cancel", messageBatchID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

// Streams the results of a Message Batch as a `.jsonl` file.
//
// Each line in the file is a JSON object containing the result of a single request
// in the Message Batch. Results are not guaranteed to be in the same order as
// requests. Use the `custom_id` field to match results to requests.
//
// Learn more about the Message Batches API in our
// [user guide](/en/docs/build-with-claude/batch-processing)
func (r *MessageBatchService) ResultsStreaming(ctx context.Context, messageBatchID string, opts ...option.RequestOption) (stream *jsonl.Stream[MessageBatchIndividualResponse]) {
	var (
		raw *http.Response
		err error
	)
	opts = append(r.Options[:], opts...)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "application/x-jsonl")}, opts...)
	if messageBatchID == "" {
		err = errors.New("missing required message_batch_id parameter")
		return
	}
	path := fmt.Sprintf("v1/messages/batches/%s/results", messageBatchID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &raw, opts...)
	return jsonl.NewStream[MessageBatchIndividualResponse](raw, err)
}

type DeletedMessageBatch struct {
	// ID of the Message Batch.
	ID string `json:"id,required"`
	// Deleted object type.
	//
	// For Message Batches, this is always `"message_batch_deleted"`.
	Type DeletedMessageBatchType `json:"type,required"`
	JSON deletedMessageBatchJSON `json:"-"`
}

// deletedMessageBatchJSON contains the JSON metadata for the struct
// [DeletedMessageBatch]
type deletedMessageBatchJSON struct {
	ID          apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *DeletedMessageBatch) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r deletedMessageBatchJSON) RawJSON() string {
	return r.raw
}

// Deleted object type.
//
// For Message Batches, this is always `"message_batch_deleted"`.
type DeletedMessageBatchType string

const (
	DeletedMessageBatchTypeMessageBatchDeleted DeletedMessageBatchType = "message_batch_deleted"
)

func (r DeletedMessageBatchType) IsKnown() bool {
	switch r {
	case DeletedMessageBatchTypeMessageBatchDeleted:
		return true
	}
	return false
}

type MessageBatch struct {
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
	ProcessingStatus MessageBatchProcessingStatus `json:"processing_status,required"`
	// Tallies requests within the Message Batch, categorized by their status.
	//
	// Requests start as `processing` and move to one of the other statuses only once
	// processing of the entire batch ends. The sum of all values always matches the
	// total number of requests in the batch.
	RequestCounts MessageBatchRequestCounts `json:"request_counts,required"`
	// URL to a `.jsonl` file containing the results of the Message Batch requests.
	// Specified only once processing ends.
	//
	// Results in the file are not guaranteed to be in the same order as requests. Use
	// the `custom_id` field to match results to requests.
	ResultsURL string `json:"results_url,required,nullable"`
	// Object type.
	//
	// For Message Batches, this is always `"message_batch"`.
	Type MessageBatchType `json:"type,required"`
	JSON messageBatchJSON `json:"-"`
}

// messageBatchJSON contains the JSON metadata for the struct [MessageBatch]
type messageBatchJSON struct {
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

func (r *MessageBatch) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageBatchJSON) RawJSON() string {
	return r.raw
}

// Processing status of the Message Batch.
type MessageBatchProcessingStatus string

const (
	MessageBatchProcessingStatusInProgress MessageBatchProcessingStatus = "in_progress"
	MessageBatchProcessingStatusCanceling  MessageBatchProcessingStatus = "canceling"
	MessageBatchProcessingStatusEnded      MessageBatchProcessingStatus = "ended"
)

func (r MessageBatchProcessingStatus) IsKnown() bool {
	switch r {
	case MessageBatchProcessingStatusInProgress, MessageBatchProcessingStatusCanceling, MessageBatchProcessingStatusEnded:
		return true
	}
	return false
}

// Object type.
//
// For Message Batches, this is always `"message_batch"`.
type MessageBatchType string

const (
	MessageBatchTypeMessageBatch MessageBatchType = "message_batch"
)

func (r MessageBatchType) IsKnown() bool {
	switch r {
	case MessageBatchTypeMessageBatch:
		return true
	}
	return false
}

type MessageBatchCanceledResult struct {
	Type MessageBatchCanceledResultType `json:"type,required"`
	JSON messageBatchCanceledResultJSON `json:"-"`
}

// messageBatchCanceledResultJSON contains the JSON metadata for the struct
// [MessageBatchCanceledResult]
type messageBatchCanceledResultJSON struct {
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *MessageBatchCanceledResult) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageBatchCanceledResultJSON) RawJSON() string {
	return r.raw
}

func (r MessageBatchCanceledResult) implementsMessageBatchResult() {}

type MessageBatchCanceledResultType string

const (
	MessageBatchCanceledResultTypeCanceled MessageBatchCanceledResultType = "canceled"
)

func (r MessageBatchCanceledResultType) IsKnown() bool {
	switch r {
	case MessageBatchCanceledResultTypeCanceled:
		return true
	}
	return false
}

type MessageBatchErroredResult struct {
	Error shared.ErrorResponse          `json:"error,required"`
	Type  MessageBatchErroredResultType `json:"type,required"`
	JSON  messageBatchErroredResultJSON `json:"-"`
}

// messageBatchErroredResultJSON contains the JSON metadata for the struct
// [MessageBatchErroredResult]
type messageBatchErroredResultJSON struct {
	Error       apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *MessageBatchErroredResult) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageBatchErroredResultJSON) RawJSON() string {
	return r.raw
}

func (r MessageBatchErroredResult) implementsMessageBatchResult() {}

type MessageBatchErroredResultType string

const (
	MessageBatchErroredResultTypeErrored MessageBatchErroredResultType = "errored"
)

func (r MessageBatchErroredResultType) IsKnown() bool {
	switch r {
	case MessageBatchErroredResultTypeErrored:
		return true
	}
	return false
}

type MessageBatchExpiredResult struct {
	Type MessageBatchExpiredResultType `json:"type,required"`
	JSON messageBatchExpiredResultJSON `json:"-"`
}

// messageBatchExpiredResultJSON contains the JSON metadata for the struct
// [MessageBatchExpiredResult]
type messageBatchExpiredResultJSON struct {
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *MessageBatchExpiredResult) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageBatchExpiredResultJSON) RawJSON() string {
	return r.raw
}

func (r MessageBatchExpiredResult) implementsMessageBatchResult() {}

type MessageBatchExpiredResultType string

const (
	MessageBatchExpiredResultTypeExpired MessageBatchExpiredResultType = "expired"
)

func (r MessageBatchExpiredResultType) IsKnown() bool {
	switch r {
	case MessageBatchExpiredResultTypeExpired:
		return true
	}
	return false
}

// This is a single line in the response `.jsonl` file and does not represent the
// response as a whole.
type MessageBatchIndividualResponse struct {
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
	Result MessageBatchResult                 `json:"result,required"`
	JSON   messageBatchIndividualResponseJSON `json:"-"`
}

// messageBatchIndividualResponseJSON contains the JSON metadata for the struct
// [MessageBatchIndividualResponse]
type messageBatchIndividualResponseJSON struct {
	CustomID    apijson.Field
	Result      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *MessageBatchIndividualResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageBatchIndividualResponseJSON) RawJSON() string {
	return r.raw
}

type MessageBatchRequestCounts struct {
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
	Succeeded int64                         `json:"succeeded,required"`
	JSON      messageBatchRequestCountsJSON `json:"-"`
}

// messageBatchRequestCountsJSON contains the JSON metadata for the struct
// [MessageBatchRequestCounts]
type messageBatchRequestCountsJSON struct {
	Canceled    apijson.Field
	Errored     apijson.Field
	Expired     apijson.Field
	Processing  apijson.Field
	Succeeded   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *MessageBatchRequestCounts) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageBatchRequestCountsJSON) RawJSON() string {
	return r.raw
}

// Processing result for this request.
//
// Contains a Message output if processing was successful, an error response if
// processing failed, or the reason why processing was not attempted, such as
// cancellation or expiration.
type MessageBatchResult struct {
	Type    MessageBatchResultType `json:"type,required"`
	Error   shared.ErrorResponse   `json:"error"`
	Message Message                `json:"message"`
	JSON    messageBatchResultJSON `json:"-"`
	union   MessageBatchResultUnion
}

// messageBatchResultJSON contains the JSON metadata for the struct
// [MessageBatchResult]
type messageBatchResultJSON struct {
	Type        apijson.Field
	Error       apijson.Field
	Message     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r messageBatchResultJSON) RawJSON() string {
	return r.raw
}

func (r *MessageBatchResult) UnmarshalJSON(data []byte) (err error) {
	*r = MessageBatchResult{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [MessageBatchResultUnion] interface which you can cast to the
// specific types for more type safety.
//
// Possible runtime types of the union are [MessageBatchSucceededResult],
// [MessageBatchErroredResult], [MessageBatchCanceledResult],
// [MessageBatchExpiredResult].
func (r MessageBatchResult) AsUnion() MessageBatchResultUnion {
	return r.union
}

// Processing result for this request.
//
// Contains a Message output if processing was successful, an error response if
// processing failed, or the reason why processing was not attempted, such as
// cancellation or expiration.
//
// Union satisfied by [MessageBatchSucceededResult], [MessageBatchErroredResult],
// [MessageBatchCanceledResult] or [MessageBatchExpiredResult].
type MessageBatchResultUnion interface {
	implementsMessageBatchResult()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*MessageBatchResultUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(MessageBatchSucceededResult{}),
			DiscriminatorValue: "succeeded",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(MessageBatchErroredResult{}),
			DiscriminatorValue: "errored",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(MessageBatchCanceledResult{}),
			DiscriminatorValue: "canceled",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(MessageBatchExpiredResult{}),
			DiscriminatorValue: "expired",
		},
	)
}

type MessageBatchResultType string

const (
	MessageBatchResultTypeSucceeded MessageBatchResultType = "succeeded"
	MessageBatchResultTypeErrored   MessageBatchResultType = "errored"
	MessageBatchResultTypeCanceled  MessageBatchResultType = "canceled"
	MessageBatchResultTypeExpired   MessageBatchResultType = "expired"
)

func (r MessageBatchResultType) IsKnown() bool {
	switch r {
	case MessageBatchResultTypeSucceeded, MessageBatchResultTypeErrored, MessageBatchResultTypeCanceled, MessageBatchResultTypeExpired:
		return true
	}
	return false
}

type MessageBatchSucceededResult struct {
	Message Message                         `json:"message,required"`
	Type    MessageBatchSucceededResultType `json:"type,required"`
	JSON    messageBatchSucceededResultJSON `json:"-"`
}

// messageBatchSucceededResultJSON contains the JSON metadata for the struct
// [MessageBatchSucceededResult]
type messageBatchSucceededResultJSON struct {
	Message     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *MessageBatchSucceededResult) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r messageBatchSucceededResultJSON) RawJSON() string {
	return r.raw
}

func (r MessageBatchSucceededResult) implementsMessageBatchResult() {}

type MessageBatchSucceededResultType string

const (
	MessageBatchSucceededResultTypeSucceeded MessageBatchSucceededResultType = "succeeded"
)

func (r MessageBatchSucceededResultType) IsKnown() bool {
	switch r {
	case MessageBatchSucceededResultTypeSucceeded:
		return true
	}
	return false
}

type MessageBatchNewParams struct {
	// List of requests for prompt completion. Each is an individual request to create
	// a Message.
	Requests param.Field[[]MessageBatchNewParamsRequest] `json:"requests,required"`
}

func (r MessageBatchNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type MessageBatchNewParamsRequest struct {
	// Developer-provided ID created for each request in a Message Batch. Useful for
	// matching results to requests, as results may be given out of request order.
	//
	// Must be unique for each request within the Message Batch.
	CustomID param.Field[string] `json:"custom_id,required"`
	// Messages API creation parameters for the individual request.
	//
	// See the [Messages API reference](/en/api/messages) for full documentation on
	// available parameters.
	Params param.Field[MessageBatchNewParamsRequestsParams] `json:"params,required"`
}

func (r MessageBatchNewParamsRequest) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Messages API creation parameters for the individual request.
//
// See the [Messages API reference](/en/api/messages) for full documentation on
// available parameters.
type MessageBatchNewParamsRequestsParams struct {
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
	Messages param.Field[[]MessageParam] `json:"messages,required"`
	// The model that will complete your prompt.\n\nSee
	// [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	Model param.Field[Model] `json:"model,required"`
	// An object describing metadata about the request.
	Metadata param.Field[MetadataParam] `json:"metadata"`
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
	System param.Field[[]TextBlockParam] `json:"system"`
	// Amount of randomness injected into the response.
	//
	// Defaults to `1.0`. Ranges from `0.0` to `1.0`. Use `temperature` closer to `0.0`
	// for analytical / multiple choice, and closer to `1.0` for creative and
	// generative tasks.
	//
	// Note that even with `temperature` of `0.0`, the results will not be fully
	// deterministic.
	Temperature param.Field[float64] `json:"temperature"`
	// Configuration for enabling Claude's extended thinking.
	//
	// When enabled, responses include `thinking` content blocks showing Claude's
	// thinking process before the final answer. Requires a minimum budget of 1,024
	// tokens and counts towards your `max_tokens` limit.
	//
	// See
	// [extended thinking](https://docs.anthropic.com/en/docs/build-with-claude/extended-thinking)
	// for details.
	Thinking param.Field[ThinkingConfigParamUnion] `json:"thinking"`
	// How the model should use the provided tools. The model can use a specific tool,
	// any available tool, or decide by itself.
	ToolChoice param.Field[ToolChoiceUnionParam] `json:"tool_choice"`
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
	//   - `input_schema`: [JSON schema](https://json-schema.org/draft/2020-12) for the
	//     tool `input` shape that the model will produce in `tool_use` output content
	//     blocks.
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
	Tools param.Field[[]ToolUnionUnionParam] `json:"tools"`
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

func (r MessageBatchNewParamsRequestsParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type MessageBatchListParams struct {
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
}

// URLQuery serializes [MessageBatchListParams]'s query parameters as `url.Values`.
func (r MessageBatchListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
