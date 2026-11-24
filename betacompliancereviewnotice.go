package anthropic

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
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

// BetaComplianceReviewNoticeService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaComplianceReviewNoticeService] method instead.
type BetaComplianceReviewNoticeService struct {
	Options []option.RequestOption
}

// NewBetaComplianceReviewNoticeService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaComplianceReviewNoticeService(opts ...option.RequestOption) (r BetaComplianceReviewNoticeService) {
	r = BetaComplianceReviewNoticeService{}
	r.Options = opts
	return
}

// List Anthropic's notices to the organizations that the API key can reach, newest
// first by default. A notice is returned only if the API key can read every case
// in its `case_ids`.
//
// To poll, store the largest `created_at` you have read, and send it as the
// at-or-after filter on `created_at`. Read every page, match results by `id`, and
// store the new largest value after the last page.
//
// An Admin API key with the `read:compliance_review_cases` scope can call this
// endpoint.
func (r *BetaComplianceReviewNoticeService) List(ctx context.Context, params BetaComplianceReviewNoticeListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaComplianceReviewNotice], err error) {
	var raw *http.Response
	if len(params.Betas) > 0 {
		headerValues := make([]string, len(params.Betas))
		for i, v := range params.Betas {
			headerValues[i] = fmt.Sprintf("%v", v)
		}
		opts = append(opts, requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
			cfg.Request.Header.Set("anthropic-beta", strings.Join(append(headerValues, cfg.Request.Header.Values("anthropic-beta")...), ","))
			return nil
		}))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/compliance/review/notices"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.compliance.review.notices.list", "/v1/compliance/review/notices", "case_id", "created_at[gt]", "created_at[gte]", "created_at[lt]", "created_at[lte]", "limit", "order", "page"))
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

// List Anthropic's notices to the organizations that the API key can reach, newest
// first by default. A notice is returned only if the API key can read every case
// in its `case_ids`.
//
// To poll, store the largest `created_at` you have read, and send it as the
// at-or-after filter on `created_at`. Read every page, match results by `id`, and
// store the new largest value after the last page.
//
// An Admin API key with the `read:compliance_review_cases` scope can call this
// endpoint.
func (r *BetaComplianceReviewNoticeService) ListAutoPaging(ctx context.Context, params BetaComplianceReviewNoticeListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaComplianceReviewNotice] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// A message from Anthropic to an organization about its review cases or about the
// organization as a whole, such as a request to respond by a date, a warning, or a
// question. A notice doesn't change after Anthropic sends it.
type BetaComplianceReviewNotice struct {
	// The unique ID of the notice (`rvn_...`).
	ID string `json:"id" api:"required"`
	// The IDs of the cases that the notice is about (`rvc_...`). Empty when the notice
	// is about the whole organization.
	CaseIDs []string `json:"case_ids" api:"required"`
	// When Anthropic sent the notice.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// What kind of notice this is.
	Kind BetaComplianceReviewNoticeKind `json:"kind" api:"required"`
	// The text of the notice.
	Message string `json:"message" api:"required"`
	// The UUID of the organization that the notice is for.
	OrganizationID string `json:"organization_id" api:"required"`
	// The time by which Anthropic asks for an answer. Null when there is no deadline.
	// Missing it changes nothing on a case, and the time to respond to a case can end
	// before it.
	RespondBy time.Time `json:"respond_by" api:"required" format:"date-time"`
	// The object type, always `compliance_review_notice`.
	Type constant.ComplianceReviewNotice `json:"type" default:"compliance_review_notice"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		CaseIDs        respjson.Field
		CreatedAt      respjson.Field
		Kind           respjson.Field
		Message        respjson.Field
		OrganizationID respjson.Field
		RespondBy      respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewNotice) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewNotice) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The kind of a review notice.
type BetaComplianceReviewNoticeKind struct {
	// `response_requested` asks you to respond to the cases in `case_ids`. For
	// `warning` and `question`, `message` explains how to answer if an answer is
	// needed. New values can appear, so accept a value that you don't know.
	//
	// Any of "response_requested", "warning", "question".
	Type BetaComplianceReviewNoticeKindType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaComplianceReviewNoticeKind) RawJSON() string { return r.JSON.raw }
func (r *BetaComplianceReviewNoticeKind) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// `response_requested` asks you to respond to the cases in `case_ids`. For
// `warning` and `question`, `message` explains how to answer if an answer is
// needed. New values can appear, so accept a value that you don't know.
type BetaComplianceReviewNoticeKindType string

const (
	BetaComplianceReviewNoticeKindTypeResponseRequested BetaComplianceReviewNoticeKindType = "response_requested"
	BetaComplianceReviewNoticeKindTypeWarning           BetaComplianceReviewNoticeKindType = "warning"
	BetaComplianceReviewNoticeKindTypeQuestion          BetaComplianceReviewNoticeKindType = "question"
)

type BetaComplianceReviewNoticeListParams struct {
	// Return only notices whose `case_ids` include this case (`rvc_...`). A case that
	// the API key can't read returns an empty list.
	CaseID param.Opt[string] `query:"case_id,omitzero" json:"-"`
	// Return only notices whose `created_at` is after this time.
	CreatedAtGt param.Opt[time.Time] `query:"created_at[gt],omitzero" format:"date-time" json:"-"`
	// Return notices whose `created_at` is at or after this time. Notices slightly
	// older than this time are also returned, so that a poll misses no notice.
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return only notices whose `created_at` is before this time.
	CreatedAtLt param.Opt[time.Time] `query:"created_at[lt],omitzero" format:"date-time" json:"-"`
	// Return only notices whose `created_at` is at or before this time.
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// Number of notices to return per page, between 1 and 100. Defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page`. Leave it out to get the
	// first page. A cursor works only with the filters and sort order of the request
	// that returned it.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Sort direction, by `created_at` and then by `id`. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaComplianceReviewOrder `query:"order,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaComplianceReviewNoticeListParams]'s query parameters as
// `url.Values`.
func (r BetaComplianceReviewNoticeListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
