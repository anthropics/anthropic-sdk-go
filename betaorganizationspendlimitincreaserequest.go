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
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// BetaOrganizationSpendLimitIncreaseRequestService contains methods and other
// services that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationSpendLimitIncreaseRequestService] method instead.
type BetaOrganizationSpendLimitIncreaseRequestService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationSpendLimitIncreaseRequestService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationSpendLimitIncreaseRequestService(opts ...option.RequestOption) (r BetaOrganizationSpendLimitIncreaseRequestService) {
	r = BetaOrganizationSpendLimitIncreaseRequestService{}
	r.Options = opts
	return
}

// Retrieve a spend limit increase request.
//
// While `pending`, the response includes a live `spend_summary` for the requester
// at the request's period.
func (r *BetaOrganizationSpendLimitIncreaseRequestService) Get(ctx context.Context, spendLimitIncreaseRequestID string, opts ...option.RequestOption) (res *BetaSpendLimitIncreaseRequest, err error) {
	opts = slices.Concat(r.Options, opts)
	if spendLimitIncreaseRequestID == "" {
		err = errors.New("missing required spend_limit_increase_request_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/spend_limit_increase_requests/%s?beta=true", url.PathEscape(spendLimitIncreaseRequestID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List spend limit increase requests, most recent first.
//
// Pending requests include a live `spend_summary` for the requester. Requests
// whose requester is no longer a member are excluded.
func (r *BetaOrganizationSpendLimitIncreaseRequestService) List(ctx context.Context, query BetaOrganizationSpendLimitIncreaseRequestListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaSpendLimitIncreaseRequest], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/spend_limit_increase_requests?beta=true"
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

// List spend limit increase requests, most recent first.
//
// Pending requests include a live `spend_summary` for the requester. Requests
// whose requester is no longer a member are excluded.
func (r *BetaOrganizationSpendLimitIncreaseRequestService) ListAutoPaging(ctx context.Context, query BetaOrganizationSpendLimitIncreaseRequestListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaSpendLimitIncreaseRequest] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

// Approve a pending spend limit increase request.
//
// Writes a per-user spend limit at `amount` for the requester and transitions the
// request to `approved`. `period` defaults to the period the member was blocked
// on. Anthropic emails the requester unless `suppress_notification` is set.
func (r *BetaOrganizationSpendLimitIncreaseRequestService) Approve(ctx context.Context, spendLimitIncreaseRequestID string, body BetaOrganizationSpendLimitIncreaseRequestApproveParams, opts ...option.RequestOption) (res *BetaOrganizationSpendLimitIncreaseRequestApproveResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if spendLimitIncreaseRequestID == "" {
		err = errors.New("missing required spend_limit_increase_request_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/spend_limit_increase_requests/%s/approve?beta=true", url.PathEscape(spendLimitIncreaseRequestID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Deny a pending spend limit increase request.
//
// Idempotent on `denied`; denying an already-`approved` request returns 400.
// Anthropic emails the requester unless `suppress_notification` is set.
func (r *BetaOrganizationSpendLimitIncreaseRequestService) Deny(ctx context.Context, spendLimitIncreaseRequestID string, body BetaOrganizationSpendLimitIncreaseRequestDenyParams, opts ...option.RequestOption) (res *BetaSpendLimitIncreaseRequest, err error) {
	opts = slices.Concat(r.Options, opts)
	if spendLimitIncreaseRequestID == "" {
		err = errors.New("missing required spend_limit_increase_request_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/spend_limit_increase_requests/%s/deny?beta=true", url.PathEscape(spendLimitIncreaseRequestID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type BetaSpendLimitIncreaseRequest struct {
	ID        string                                  `json:"id" api:"required"`
	Actor     BetaSpendLimitIncreaseRequestActorUnion `json:"actor" api:"required"`
	CreatedAt time.Time                               `json:"created_at" api:"required" format:"date-time"`
	// Any of "daily", "monthly", "weekly".
	Period     BetaSpendLimitPeriod                         `json:"period" api:"required"`
	ResolvedAt time.Time                                    `json:"resolved_at" api:"required" format:"date-time"`
	ResolvedBy BetaSpendLimitIncreaseRequestResolvedByUnion `json:"resolved_by" api:"required"`
	// Per-member effective-limit report row (`GET /spend_limits/effective`).
	SpendSummary BetaSpendSummary `json:"spend_summary" api:"required"`
	// Any of "approved", "denied", "pending".
	Status BetaSpendLimitIncreaseRequestStatus `json:"status" api:"required"`
	Type   constant.SpendLimitIncreaseRequest  `json:"type" default:"spend_limit_increase_request"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		Actor        respjson.Field
		CreatedAt    respjson.Field
		Period       respjson.Field
		ResolvedAt   respjson.Field
		ResolvedBy   respjson.Field
		SpendSummary respjson.Field
		Status       respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaSpendLimitIncreaseRequest) RawJSON() string { return r.JSON.raw }
func (r *BetaSpendLimitIncreaseRequest) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaSpendLimitIncreaseRequestActorUnion contains all possible properties and
// values from [BetaSpendLimitUserActor], [BetaSpendLimitScopedAPIKeyActor].
//
// Use the [BetaSpendLimitIncreaseRequestActorUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaSpendLimitIncreaseRequestActorUnion struct {
	// This field is from variant [BetaSpendLimitUserActor].
	Deleted bool `json:"deleted"`
	// This field is from variant [BetaSpendLimitUserActor].
	EmailAddress string `json:"email_address"`
	// This field is from variant [BetaSpendLimitUserActor].
	Name string `json:"name"`
	// Any of "user_actor", "scoped_api_key_actor".
	Type string `json:"type"`
	// This field is from variant [BetaSpendLimitUserActor].
	UserID string `json:"user_id"`
	// This field is from variant [BetaSpendLimitScopedAPIKeyActor].
	ScopedAPIKeyID string `json:"scoped_api_key_id"`
	JSON           struct {
		Deleted        respjson.Field
		EmailAddress   respjson.Field
		Name           respjson.Field
		Type           respjson.Field
		UserID         respjson.Field
		ScopedAPIKeyID respjson.Field
		raw            string
	} `json:"-"`
}

// anyBetaSpendLimitIncreaseRequestActor is implemented by each variant of
// [BetaSpendLimitIncreaseRequestActorUnion] to add type safety for the return type
// of [BetaSpendLimitIncreaseRequestActorUnion.AsAny]
type anyBetaSpendLimitIncreaseRequestActor interface {
	implBetaSpendLimitIncreaseRequestActorUnion()
}

func (BetaSpendLimitUserActor) implBetaSpendLimitIncreaseRequestActorUnion()         {}
func (BetaSpendLimitScopedAPIKeyActor) implBetaSpendLimitIncreaseRequestActorUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaSpendLimitIncreaseRequestActorUnion.AsAny().(type) {
//	case anthropic.BetaSpendLimitUserActor:
//	case anthropic.BetaSpendLimitScopedAPIKeyActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaSpendLimitIncreaseRequestActorUnion) AsAny() anyBetaSpendLimitIncreaseRequestActor {
	switch u.Type {
	case "user_actor":
		return u.AsUserActor()
	case "scoped_api_key_actor":
		return u.AsScopedAPIKeyActor()
	}
	return nil
}

func (u BetaSpendLimitIncreaseRequestActorUnion) AsUserActor() (v BetaSpendLimitUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendLimitIncreaseRequestActorUnion) AsScopedAPIKeyActor() (v BetaSpendLimitScopedAPIKeyActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaSpendLimitIncreaseRequestActorUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaSpendLimitIncreaseRequestActorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaSpendLimitIncreaseRequestResolvedByUnion contains all possible properties
// and values from [BetaSpendLimitUserActor], [BetaSpendLimitScopedAPIKeyActor].
//
// Use the [BetaSpendLimitIncreaseRequestResolvedByUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaSpendLimitIncreaseRequestResolvedByUnion struct {
	// This field is from variant [BetaSpendLimitUserActor].
	Deleted bool `json:"deleted"`
	// This field is from variant [BetaSpendLimitUserActor].
	EmailAddress string `json:"email_address"`
	// This field is from variant [BetaSpendLimitUserActor].
	Name string `json:"name"`
	// Any of "user_actor", "scoped_api_key_actor".
	Type string `json:"type"`
	// This field is from variant [BetaSpendLimitUserActor].
	UserID string `json:"user_id"`
	// This field is from variant [BetaSpendLimitScopedAPIKeyActor].
	ScopedAPIKeyID string `json:"scoped_api_key_id"`
	JSON           struct {
		Deleted        respjson.Field
		EmailAddress   respjson.Field
		Name           respjson.Field
		Type           respjson.Field
		UserID         respjson.Field
		ScopedAPIKeyID respjson.Field
		raw            string
	} `json:"-"`
}

// anyBetaSpendLimitIncreaseRequestResolvedBy is implemented by each variant of
// [BetaSpendLimitIncreaseRequestResolvedByUnion] to add type safety for the return
// type of [BetaSpendLimitIncreaseRequestResolvedByUnion.AsAny]
type anyBetaSpendLimitIncreaseRequestResolvedBy interface {
	implBetaSpendLimitIncreaseRequestResolvedByUnion()
}

func (BetaSpendLimitUserActor) implBetaSpendLimitIncreaseRequestResolvedByUnion()         {}
func (BetaSpendLimitScopedAPIKeyActor) implBetaSpendLimitIncreaseRequestResolvedByUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaSpendLimitIncreaseRequestResolvedByUnion.AsAny().(type) {
//	case anthropic.BetaSpendLimitUserActor:
//	case anthropic.BetaSpendLimitScopedAPIKeyActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaSpendLimitIncreaseRequestResolvedByUnion) AsAny() anyBetaSpendLimitIncreaseRequestResolvedBy {
	switch u.Type {
	case "user_actor":
		return u.AsUserActor()
	case "scoped_api_key_actor":
		return u.AsScopedAPIKeyActor()
	}
	return nil
}

func (u BetaSpendLimitIncreaseRequestResolvedByUnion) AsUserActor() (v BetaSpendLimitUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaSpendLimitIncreaseRequestResolvedByUnion) AsScopedAPIKeyActor() (v BetaSpendLimitScopedAPIKeyActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaSpendLimitIncreaseRequestResolvedByUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaSpendLimitIncreaseRequestResolvedByUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaSpendLimitIncreaseRequestStatus string

const (
	BetaSpendLimitIncreaseRequestStatusApproved BetaSpendLimitIncreaseRequestStatus = "approved"
	BetaSpendLimitIncreaseRequestStatusDenied   BetaSpendLimitIncreaseRequestStatus = "denied"
	BetaSpendLimitIncreaseRequestStatusPending  BetaSpendLimitIncreaseRequestStatus = "pending"
)

type BetaOrganizationSpendLimitIncreaseRequestApproveResponse struct {
	ID        string                                                             `json:"id" api:"required"`
	Actor     BetaOrganizationSpendLimitIncreaseRequestApproveResponseActorUnion `json:"actor" api:"required"`
	CreatedAt time.Time                                                          `json:"created_at" api:"required" format:"date-time"`
	// Any of "daily", "monthly", "weekly".
	Period     BetaSpendLimitPeriod                                                    `json:"period" api:"required"`
	ResolvedAt time.Time                                                               `json:"resolved_at" api:"required" format:"date-time"`
	ResolvedBy BetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedByUnion `json:"resolved_by" api:"required"`
	// A configured spend limit: a cap on metered spend for one scope and period.
	SpendLimit BetaSpendLimit `json:"spend_limit" api:"required"`
	// Per-member effective-limit report row (`GET /spend_limits/effective`).
	SpendSummary BetaSpendSummary `json:"spend_summary" api:"required"`
	// Any of "approved", "denied", "pending".
	Status BetaSpendLimitIncreaseRequestStatus `json:"status" api:"required"`
	Type   constant.SpendLimitIncreaseRequest  `json:"type" default:"spend_limit_increase_request"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		Actor        respjson.Field
		CreatedAt    respjson.Field
		Period       respjson.Field
		ResolvedAt   respjson.Field
		ResolvedBy   respjson.Field
		SpendLimit   respjson.Field
		SpendSummary respjson.Field
		Status       respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaOrganizationSpendLimitIncreaseRequestApproveResponse) RawJSON() string { return r.JSON.raw }
func (r *BetaOrganizationSpendLimitIncreaseRequestApproveResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaOrganizationSpendLimitIncreaseRequestApproveResponseActorUnion contains all
// possible properties and values from [BetaSpendLimitUserActor],
// [BetaSpendLimitScopedAPIKeyActor].
//
// Use the
// [BetaOrganizationSpendLimitIncreaseRequestApproveResponseActorUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaOrganizationSpendLimitIncreaseRequestApproveResponseActorUnion struct {
	// This field is from variant [BetaSpendLimitUserActor].
	Deleted bool `json:"deleted"`
	// This field is from variant [BetaSpendLimitUserActor].
	EmailAddress string `json:"email_address"`
	// This field is from variant [BetaSpendLimitUserActor].
	Name string `json:"name"`
	// Any of "user_actor", "scoped_api_key_actor".
	Type string `json:"type"`
	// This field is from variant [BetaSpendLimitUserActor].
	UserID string `json:"user_id"`
	// This field is from variant [BetaSpendLimitScopedAPIKeyActor].
	ScopedAPIKeyID string `json:"scoped_api_key_id"`
	JSON           struct {
		Deleted        respjson.Field
		EmailAddress   respjson.Field
		Name           respjson.Field
		Type           respjson.Field
		UserID         respjson.Field
		ScopedAPIKeyID respjson.Field
		raw            string
	} `json:"-"`
}

// anyBetaOrganizationSpendLimitIncreaseRequestApproveResponseActor is implemented
// by each variant of
// [BetaOrganizationSpendLimitIncreaseRequestApproveResponseActorUnion] to add type
// safety for the return type of
// [BetaOrganizationSpendLimitIncreaseRequestApproveResponseActorUnion.AsAny]
type anyBetaOrganizationSpendLimitIncreaseRequestApproveResponseActor interface {
	implBetaOrganizationSpendLimitIncreaseRequestApproveResponseActorUnion()
}

func (BetaSpendLimitUserActor) implBetaOrganizationSpendLimitIncreaseRequestApproveResponseActorUnion() {
}
func (BetaSpendLimitScopedAPIKeyActor) implBetaOrganizationSpendLimitIncreaseRequestApproveResponseActorUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaOrganizationSpendLimitIncreaseRequestApproveResponseActorUnion.AsAny().(type) {
//	case anthropic.BetaSpendLimitUserActor:
//	case anthropic.BetaSpendLimitScopedAPIKeyActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaOrganizationSpendLimitIncreaseRequestApproveResponseActorUnion) AsAny() anyBetaOrganizationSpendLimitIncreaseRequestApproveResponseActor {
	switch u.Type {
	case "user_actor":
		return u.AsUserActor()
	case "scoped_api_key_actor":
		return u.AsScopedAPIKeyActor()
	}
	return nil
}

func (u BetaOrganizationSpendLimitIncreaseRequestApproveResponseActorUnion) AsUserActor() (v BetaSpendLimitUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaOrganizationSpendLimitIncreaseRequestApproveResponseActorUnion) AsScopedAPIKeyActor() (v BetaSpendLimitScopedAPIKeyActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaOrganizationSpendLimitIncreaseRequestApproveResponseActorUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *BetaOrganizationSpendLimitIncreaseRequestApproveResponseActorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedByUnion contains
// all possible properties and values from [BetaSpendLimitUserActor],
// [BetaSpendLimitScopedAPIKeyActor].
//
// Use the
// [BetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedByUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedByUnion struct {
	// This field is from variant [BetaSpendLimitUserActor].
	Deleted bool `json:"deleted"`
	// This field is from variant [BetaSpendLimitUserActor].
	EmailAddress string `json:"email_address"`
	// This field is from variant [BetaSpendLimitUserActor].
	Name string `json:"name"`
	// Any of "user_actor", "scoped_api_key_actor".
	Type string `json:"type"`
	// This field is from variant [BetaSpendLimitUserActor].
	UserID string `json:"user_id"`
	// This field is from variant [BetaSpendLimitScopedAPIKeyActor].
	ScopedAPIKeyID string `json:"scoped_api_key_id"`
	JSON           struct {
		Deleted        respjson.Field
		EmailAddress   respjson.Field
		Name           respjson.Field
		Type           respjson.Field
		UserID         respjson.Field
		ScopedAPIKeyID respjson.Field
		raw            string
	} `json:"-"`
}

// anyBetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedBy is
// implemented by each variant of
// [BetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedByUnion] to add
// type safety for the return type of
// [BetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedByUnion.AsAny]
type anyBetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedBy interface {
	implBetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedByUnion()
}

func (BetaSpendLimitUserActor) implBetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedByUnion() {
}
func (BetaSpendLimitScopedAPIKeyActor) implBetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedByUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedByUnion.AsAny().(type) {
//	case anthropic.BetaSpendLimitUserActor:
//	case anthropic.BetaSpendLimitScopedAPIKeyActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedByUnion) AsAny() anyBetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedBy {
	switch u.Type {
	case "user_actor":
		return u.AsUserActor()
	case "scoped_api_key_actor":
		return u.AsScopedAPIKeyActor()
	}
	return nil
}

func (u BetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedByUnion) AsUserActor() (v BetaSpendLimitUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u BetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedByUnion) AsScopedAPIKeyActor() (v BetaSpendLimitScopedAPIKeyActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedByUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *BetaOrganizationSpendLimitIncreaseRequestApproveResponseResolvedByUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationSpendLimitIncreaseRequestListParams struct {
	// Opaque cursor from a previous response's `next_page`.
	Page  param.Opt[string] `query:"page,omitzero" json:"-"`
	Limit param.Opt[int64]  `query:"limit,omitzero" json:"-"`
	// Filter by requester, as `user_...` tagged IDs.
	ActorIDs []string `query:"actor_ids,omitzero" json:"-"`
	// Filter by status. Omit to return all.
	Status []BetaSpendLimitIncreaseRequestStatus `query:"status,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationSpendLimitIncreaseRequestListParams]'s
// query parameters as `url.Values`.
func (r BetaOrganizationSpendLimitIncreaseRequestListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type BetaOrganizationSpendLimitIncreaseRequestApproveParams struct {
	// New per-user spend limit as a non-negative integer decimal string (minor units).
	Amount               string          `json:"amount" api:"required"`
	SuppressNotification param.Opt[bool] `json:"suppress_notification,omitzero"`
	// Any of "daily", "monthly", "weekly".
	Period BetaSpendLimitPeriod `json:"period,omitzero"`
	paramObj
}

func (r BetaOrganizationSpendLimitIncreaseRequestApproveParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationSpendLimitIncreaseRequestApproveParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationSpendLimitIncreaseRequestApproveParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaOrganizationSpendLimitIncreaseRequestDenyParams struct {
	SuppressNotification param.Opt[bool] `json:"suppress_notification,omitzero"`
	paramObj
}

func (r BetaOrganizationSpendLimitIncreaseRequestDenyParams) MarshalJSON() (data []byte, err error) {
	type shadow BetaOrganizationSpendLimitIncreaseRequestDenyParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BetaOrganizationSpendLimitIncreaseRequestDenyParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
