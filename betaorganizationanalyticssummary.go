package anthropic

import (
	"context"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal/apiquery"
	"github.com/anthropics/anthropic-sdk-go/internal/requestconfig"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/pagination"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
)

// BetaOrganizationAnalyticsSummaryService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationAnalyticsSummaryService] method instead.
type BetaOrganizationAnalyticsSummaryService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationAnalyticsSummaryService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationAnalyticsSummaryService(opts ...option.RequestOption) (r BetaOrganizationAnalyticsSummaryService) {
	r = BetaOrganizationAnalyticsSummaryService{}
	r.Options = opts
	return
}

// Get organization-wide activity summaries for a date range.
//
// Returns one entry per day from `starting_date` (inclusive) to `ending_date`
// (exclusive) in `data`, the same `data` / `next_page` envelope as the other
// analytics list endpoints; the series is currently returned in full, so
// `next_page` is always null. Data is typically available with a 1-day lag and may
// be revised by a few percent over the following days: when `ending_date` is
// omitted it defaults to the most recent available day + 1, so the last entry
// covers the most recent available day. The series can be scoped to an RBAC group
// via `filter[]=rbac_group_id:{id}`. Available to organizations on a Claude
// Enterprise plan. Requires an API key with the `read:analytics` scope.
func (r *BetaOrganizationAnalyticsSummaryService) List(ctx context.Context, query BetaOrganizationAnalyticsSummaryListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaAnalyticsSingleDayActivitySummary], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/analytics/summaries?beta=true"
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

// Get organization-wide activity summaries for a date range.
//
// Returns one entry per day from `starting_date` (inclusive) to `ending_date`
// (exclusive) in `data`, the same `data` / `next_page` envelope as the other
// analytics list endpoints; the series is currently returned in full, so
// `next_page` is always null. Data is typically available with a 1-day lag and may
// be revised by a few percent over the following days: when `ending_date` is
// omitted it defaults to the most recent available day + 1, so the last entry
// covers the most recent available day. The series can be scoped to an RBAC group
// via `filter[]=rbac_group_id:{id}`. Available to organizations on a Claude
// Enterprise plan. Requires an API key with the `read:analytics` scope.
func (r *BetaOrganizationAnalyticsSummaryService) ListAutoPaging(ctx context.Context, query BetaOrganizationAnalyticsSummaryListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaAnalyticsSingleDayActivitySummary] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

type BetaOrganizationAnalyticsSummaryListParams struct {
	// UTC date in YYYY-MM-DD format. Start of the date range (inclusive). Data is
	// typically available with a 1-day lag (varies by query; the error for a
	// too-recent date names the latest available day) and may be revised by a few
	// percent over the following days. No earlier than 2026-01-01.
	StartingDate time.Time `query:"starting_date" api:"required" format:"date" json:"-"`
	// UTC date in YYYY-MM-DD format. End of the date range (exclusive). Data is
	// typically available with a 1-day lag, so this can be at most today — which is
	// also the default when omitted, making the last entry cover the most recent
	// available day. Data may be revised by a few percent over the following days. The
	// range may span at most 366 days.
	EndingDate param.Opt[time.Time] `query:"ending_date,omitzero" format:"date" json:"-"`
	// Number of results per page (1-1000, default 100). The day series (at most 366
	// entries) is currently returned in full in a single page, so `limit` does not yet
	// shorten it.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page` field. `next_page` is
	// currently always null, so there is never a cursor to send.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Filters as `dimension:value`. Only `rbac_group_id` is supported (e.g.
	// `filter[]=rbac_group_id:{id}`); repeat the param to OR across groups. Scopes the
	// whole day series to members of the matching group(s), re-aggregated from
	// member-level activity — org-wide seat/invite fields and the adoption rates
	// derived from them are null on scoped rows. `rbac_group_id` accepts the tagged id
	// (`rbac_group_...`, as emitted in responses and by the spend-limits API) or a
	// bare group UUID, and matches users who held the group at any point during each
	// UTC day (time-of-usage attribution). At most 100 entries.
	Filter []string `query:"filter,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationAnalyticsSummaryListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationAnalyticsSummaryListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
