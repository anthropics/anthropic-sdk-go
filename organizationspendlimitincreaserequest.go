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

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// OrganizationSpendLimitIncreaseRequestService contains methods and other services
// that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationSpendLimitIncreaseRequestService] method instead.
type OrganizationSpendLimitIncreaseRequestService struct {
	Options []option.RequestOption
}

// NewOrganizationSpendLimitIncreaseRequestService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewOrganizationSpendLimitIncreaseRequestService(opts ...option.RequestOption) (r OrganizationSpendLimitIncreaseRequestService) {
	r = OrganizationSpendLimitIncreaseRequestService{}
	r.Options = opts
	return
}

// Retrieve a spend limit increase request.
//
// While `pending`, the response includes a live `spend_summary` for the requester
// at the request's period.
func (r *OrganizationSpendLimitIncreaseRequestService) Get(ctx context.Context, spendLimitIncreaseRequestID string, opts ...option.RequestOption) (res *SpendLimitIncreaseRequest, err error) {
	opts = slices.Concat(r.Options, opts)
	if spendLimitIncreaseRequestID == "" {
		err = errors.New("missing required spend_limit_increase_request_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/spend_limit_increase_requests/%s", requestconfig.EscapePathSegment(spendLimitIncreaseRequestID))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.spend_limits.increase_requests.retrieve", "/v1/organizations/spend_limit_increase_requests/{spend_limit_increase_request_id}"), requestconfig.WithSpanAttribute("anthropic.organization.spend_limit_increase_request.id", spendLimitIncreaseRequestID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List spend limit increase requests, most recent first.
//
// Pending requests include a live `spend_summary` for the requester. Requests
// whose requester is no longer a member are excluded.
func (r *OrganizationSpendLimitIncreaseRequestService) List(ctx context.Context, query OrganizationSpendLimitIncreaseRequestListParams, opts ...option.RequestOption) (res *pagination.PageCursor[SpendLimitIncreaseRequest], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/spend_limit_increase_requests"
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.spend_limits.increase_requests.list", "/v1/organizations/spend_limit_increase_requests", "actor_ids", "limit", "page", "status"))
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
func (r *OrganizationSpendLimitIncreaseRequestService) ListAutoPaging(ctx context.Context, query OrganizationSpendLimitIncreaseRequestListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[SpendLimitIncreaseRequest] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

// Approve a pending spend limit increase request.
//
// Writes a per-user spend limit at `amount` for the requester and transitions the
// request to `approved`. `period` defaults to the period the member was blocked
// on. Anthropic emails the requester unless `suppress_notification` is set.
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeSandboxEgressRestricted] (HTTP 403)
func (r *OrganizationSpendLimitIncreaseRequestService) Approve(ctx context.Context, spendLimitIncreaseRequestID string, body OrganizationSpendLimitIncreaseRequestApproveParams, opts ...option.RequestOption) (res *OrganizationSpendLimitIncreaseRequestApproveResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if spendLimitIncreaseRequestID == "" {
		err = errors.New("missing required spend_limit_increase_request_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/spend_limit_increase_requests/%s/approve", requestconfig.EscapePathSegment(spendLimitIncreaseRequestID))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.spend_limits.increase_requests.approve", "/v1/organizations/spend_limit_increase_requests/{spend_limit_increase_request_id}/approve"), requestconfig.WithSpanAttribute("anthropic.organization.spend_limit_increase_request.id", spendLimitIncreaseRequestID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Deny a pending spend limit increase request.
//
// Idempotent on `denied`; denying an already-`approved` request returns 400.
// Anthropic emails the requester unless `suppress_notification` is set.
//
// When the returned error is an [*apierror.Error], compare the
// [apierror.Error.ErrorCode] against these values.
//
//   - [ErrorCodeSandboxEgressRestricted] (HTTP 403)
func (r *OrganizationSpendLimitIncreaseRequestService) Deny(ctx context.Context, spendLimitIncreaseRequestID string, body OrganizationSpendLimitIncreaseRequestDenyParams, opts ...option.RequestOption) (res *SpendLimitIncreaseRequest, err error) {
	opts = slices.Concat(r.Options, opts)
	if spendLimitIncreaseRequestID == "" {
		err = errors.New("missing required spend_limit_increase_request_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/organizations/spend_limit_increase_requests/%s/deny", requestconfig.EscapePathSegment(spendLimitIncreaseRequestID))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.spend_limits.increase_requests.deny", "/v1/organizations/spend_limit_increase_requests/{spend_limit_increase_request_id}/deny"), requestconfig.WithSpanAttribute("anthropic.organization.spend_limit_increase_request.id", spendLimitIncreaseRequestID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type SpendLimitIncreaseRequest struct {
	ID string `json:"id" api:"required"`
	// Who submitted the increase request.
	Actor     SpendLimitIncreaseRequestActorUnion `json:"actor" api:"required"`
	CreatedAt time.Time                           `json:"created_at" api:"required" format:"date-time"`
	// Any of "daily", "monthly", "weekly".
	Period     SpendLimitPeriod `json:"period" api:"required"`
	ResolvedAt time.Time        `json:"resolved_at" api:"required" format:"date-time"`
	// Who approved or denied the increase request.
	ResolvedBy SpendLimitIncreaseRequestResolvedByUnion `json:"resolved_by" api:"required"`
	// Per-member effective-limit report row (`GET /spend_limits/effective`).
	SpendSummary SpendSummary `json:"spend_summary" api:"required"`
	// Any of "approved", "denied", "pending".
	Status SpendLimitIncreaseRequestStatus    `json:"status" api:"required"`
	Type   constant.SpendLimitIncreaseRequest `json:"type" default:"spend_limit_increase_request"`
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
func (r SpendLimitIncreaseRequest) RawJSON() string { return r.JSON.raw }
func (r *SpendLimitIncreaseRequest) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SpendLimitIncreaseRequestActorUnion contains all possible properties and values
// from [SpendLimitUserActor], [SpendLimitScopedAPIKeyActor].
//
// Use the [SpendLimitIncreaseRequestActorUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type SpendLimitIncreaseRequestActorUnion struct {
	// This field is from variant [SpendLimitUserActor].
	Deleted bool `json:"deleted"`
	// This field is from variant [SpendLimitUserActor].
	EmailAddress string `json:"email_address"`
	// This field is from variant [SpendLimitUserActor].
	Name string `json:"name"`
	// Any of "user_actor", "scoped_api_key_actor".
	Type string `json:"type"`
	// This field is from variant [SpendLimitUserActor].
	UserID string `json:"user_id"`
	// This field is from variant [SpendLimitScopedAPIKeyActor].
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

// anySpendLimitIncreaseRequestActor is implemented by each variant of
// [SpendLimitIncreaseRequestActorUnion] to add type safety for the return type of
// [SpendLimitIncreaseRequestActorUnion.AsAny]
type anySpendLimitIncreaseRequestActor interface {
	implSpendLimitIncreaseRequestActorUnion()
}

func (SpendLimitUserActor) implSpendLimitIncreaseRequestActorUnion()         {}
func (SpendLimitScopedAPIKeyActor) implSpendLimitIncreaseRequestActorUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := SpendLimitIncreaseRequestActorUnion.AsAny().(type) {
//	case anthropic.SpendLimitUserActor:
//	case anthropic.SpendLimitScopedAPIKeyActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u SpendLimitIncreaseRequestActorUnion) AsAny() anySpendLimitIncreaseRequestActor {
	switch u.Type {
	case "user_actor":
		return u.AsUserActor()
	case "scoped_api_key_actor":
		return u.AsScopedAPIKeyActor()
	}
	return nil
}

func (u SpendLimitIncreaseRequestActorUnion) AsUserActor() (v SpendLimitUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendLimitIncreaseRequestActorUnion) AsScopedAPIKeyActor() (v SpendLimitScopedAPIKeyActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u SpendLimitIncreaseRequestActorUnion) RawJSON() string { return u.JSON.raw }

func (r *SpendLimitIncreaseRequestActorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SpendLimitIncreaseRequestResolvedByUnion contains all possible properties and
// values from [SpendLimitUserActor], [SpendLimitScopedAPIKeyActor].
//
// Use the [SpendLimitIncreaseRequestResolvedByUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type SpendLimitIncreaseRequestResolvedByUnion struct {
	// This field is from variant [SpendLimitUserActor].
	Deleted bool `json:"deleted"`
	// This field is from variant [SpendLimitUserActor].
	EmailAddress string `json:"email_address"`
	// This field is from variant [SpendLimitUserActor].
	Name string `json:"name"`
	// Any of "user_actor", "scoped_api_key_actor".
	Type string `json:"type"`
	// This field is from variant [SpendLimitUserActor].
	UserID string `json:"user_id"`
	// This field is from variant [SpendLimitScopedAPIKeyActor].
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

// anySpendLimitIncreaseRequestResolvedBy is implemented by each variant of
// [SpendLimitIncreaseRequestResolvedByUnion] to add type safety for the return
// type of [SpendLimitIncreaseRequestResolvedByUnion.AsAny]
type anySpendLimitIncreaseRequestResolvedBy interface {
	implSpendLimitIncreaseRequestResolvedByUnion()
}

func (SpendLimitUserActor) implSpendLimitIncreaseRequestResolvedByUnion()         {}
func (SpendLimitScopedAPIKeyActor) implSpendLimitIncreaseRequestResolvedByUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := SpendLimitIncreaseRequestResolvedByUnion.AsAny().(type) {
//	case anthropic.SpendLimitUserActor:
//	case anthropic.SpendLimitScopedAPIKeyActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u SpendLimitIncreaseRequestResolvedByUnion) AsAny() anySpendLimitIncreaseRequestResolvedBy {
	switch u.Type {
	case "user_actor":
		return u.AsUserActor()
	case "scoped_api_key_actor":
		return u.AsScopedAPIKeyActor()
	}
	return nil
}

func (u SpendLimitIncreaseRequestResolvedByUnion) AsUserActor() (v SpendLimitUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SpendLimitIncreaseRequestResolvedByUnion) AsScopedAPIKeyActor() (v SpendLimitScopedAPIKeyActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u SpendLimitIncreaseRequestResolvedByUnion) RawJSON() string { return u.JSON.raw }

func (r *SpendLimitIncreaseRequestResolvedByUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SpendLimitIncreaseRequestStatus string

const (
	SpendLimitIncreaseRequestStatusApproved SpendLimitIncreaseRequestStatus = "approved"
	SpendLimitIncreaseRequestStatusDenied   SpendLimitIncreaseRequestStatus = "denied"
	SpendLimitIncreaseRequestStatusPending  SpendLimitIncreaseRequestStatus = "pending"
)

type OrganizationSpendLimitIncreaseRequestApproveResponse struct {
	ID string `json:"id" api:"required"`
	// Who submitted the increase request.
	Actor     SpendLimitIncreaseRequestActorUnion `json:"actor" api:"required"`
	CreatedAt time.Time                           `json:"created_at" api:"required" format:"date-time"`
	// Any of "daily", "monthly", "weekly".
	Period     SpendLimitPeriod `json:"period" api:"required"`
	ResolvedAt time.Time        `json:"resolved_at" api:"required" format:"date-time"`
	// Who approved or denied the increase request.
	ResolvedBy SpendLimitIncreaseRequestResolvedByUnion `json:"resolved_by" api:"required"`
	// A configured spend limit: a cap on metered spend for one scope and period.
	SpendLimit SpendLimit `json:"spend_limit" api:"required"`
	// Per-member effective-limit report row (`GET /spend_limits/effective`).
	SpendSummary SpendSummary `json:"spend_summary" api:"required"`
	// Any of "approved", "denied", "pending".
	Status SpendLimitIncreaseRequestStatus    `json:"status" api:"required"`
	Type   constant.SpendLimitIncreaseRequest `json:"type" default:"spend_limit_increase_request"`
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
func (r OrganizationSpendLimitIncreaseRequestApproveResponse) RawJSON() string { return r.JSON.raw }
func (r *OrganizationSpendLimitIncreaseRequestApproveResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationSpendLimitIncreaseRequestListParams struct {
	// Opaque cursor from a previous response's `next_page`.
	Page  param.Opt[string] `query:"page,omitzero" json:"-"`
	Limit param.Opt[int64]  `query:"limit,omitzero" json:"-"`
	// Filter by requester, as `user_...` tagged IDs.
	ActorIDs []string `query:"actor_ids,omitzero" json:"-"`
	// Filter by status. Omit to return all.
	Status []SpendLimitIncreaseRequestStatus `query:"status,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [OrganizationSpendLimitIncreaseRequestListParams]'s query
// parameters as `url.Values`.
func (r OrganizationSpendLimitIncreaseRequestListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type OrganizationSpendLimitIncreaseRequestApproveParams struct {
	// New per-user spend limit as a non-negative integer decimal string (minor units).
	Amount               string          `json:"amount" api:"required"`
	SuppressNotification param.Opt[bool] `json:"suppress_notification,omitzero"`
	// Any of "daily", "monthly", "weekly".
	Period SpendLimitPeriod `json:"period,omitzero"`
	paramObj
}

func (r OrganizationSpendLimitIncreaseRequestApproveParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationSpendLimitIncreaseRequestApproveParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationSpendLimitIncreaseRequestApproveParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type OrganizationSpendLimitIncreaseRequestDenyParams struct {
	SuppressNotification param.Opt[bool] `json:"suppress_notification,omitzero"`
	paramObj
}

func (r OrganizationSpendLimitIncreaseRequestDenyParams) MarshalJSON() (data []byte, err error) {
	type shadow OrganizationSpendLimitIncreaseRequestDenyParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OrganizationSpendLimitIncreaseRequestDenyParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
