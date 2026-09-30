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

// BetaOrganizationAnalyticsAppChatProjectService contains methods and other
// services that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationAnalyticsAppChatProjectService] method instead.
type BetaOrganizationAnalyticsAppChatProjectService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationAnalyticsAppChatProjectService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationAnalyticsAppChatProjectService(opts ...option.RequestOption) (r BetaOrganizationAnalyticsAppChatProjectService) {
	r = BetaOrganizationAnalyticsAppChatProjectService{}
	r.Options = opts
	return
}

// Get per-project activity for a given day, with cursor-based pagination.
//
// Returns activity metrics for each project in the organization, sorted by project
// ID. Use `group_by[]` to break projects out per member or per RBAC group, and
// `filter[]` to scope results; the parameter descriptions list the supported
// dimensions. Available to organizations on a Claude Enterprise plan. Requires an
// API key with the `read:analytics` scope.
func (r *BetaOrganizationAnalyticsAppChatProjectService) List(ctx context.Context, query BetaOrganizationAnalyticsAppChatProjectListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaAnalyticsProjectActivity], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/analytics/apps/chat/projects?beta=true"
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

// Get per-project activity for a given day, with cursor-based pagination.
//
// Returns activity metrics for each project in the organization, sorted by project
// ID. Use `group_by[]` to break projects out per member or per RBAC group, and
// `filter[]` to scope results; the parameter descriptions list the supported
// dimensions. Available to organizations on a Claude Enterprise plan. Requires an
// API key with the `read:analytics` scope.
func (r *BetaOrganizationAnalyticsAppChatProjectService) ListAutoPaging(ctx context.Context, query BetaOrganizationAnalyticsAppChatProjectListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaAnalyticsProjectActivity] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

type BetaOrganizationAnalyticsAppChatProjectListParams struct {
	// UTC date in YYYY-MM-DD format. The day to get project activity for. Data is
	// typically available with a 1-day lag (varies by query; the error for a
	// too-recent date names the latest available day) and may be revised by a few
	// percent over the following days. No earlier than 2026-01-01.
	Date param.Opt[time.Time] `query:"date,omitzero" format:"date" json:"-"`
	// UTC date in YYYY-MM-DD format. End of the date range (exclusive); only valid
	// with `starting_date`. Data is typically available with a 1-day lag (varies by
	// query; the error for a too-recent date names the latest available day), so this
	// can be at most today — which is also the default when omitted, resolved once
	// when the first page is served and reused for the rest of the pagination
	// sequence. At most 366 days after `starting_date`.
	EndingDate param.Opt[time.Time] `query:"ending_date,omitzero" format:"date" json:"-"`
	// Number of results per page (1-1000, default 100).
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Sort field. Restricted to the endpoint's sort column plus its rankable metrics
	// (metrics default to descending; a few metrics rank in date-range mode only, per
	// the endpoint's documented orderable set).
	OrderBy param.Opt[string] `query:"order_by,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page` field.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// UTC date in YYYY-MM-DD format. Start of a date range (inclusive). Enables rollup
	// mode: one row per entity aggregated over the whole range — addable counters are
	// summed across days, and a distinct count is never summed where summing could
	// double-count (a field's range value is recomputed exactly over the window,
	// approximate via HLL with typical error under 2%, null, or — for the
	// creation-event counts, whose per-day values cannot overlap — a per-day sum that
	// is itself exact; each field's own description says which). Use either `date` or
	// `starting_date`, not both. Data is typically available with a 1-day lag (varies
	// by query; the error for a too-recent date names the latest available day) and
	// may be revised by a few percent over the following days. No earlier than
	// 2026-01-01.
	StartingDate param.Opt[time.Time] `query:"starting_date,omitzero" format:"date" json:"-"`
	// Filters as `dimension:value`, e.g. `filter[]=rbac_group_id:{id}`. Repeat the
	// param for OR within a dimension and across dimensions for AND. Supported
	// dimensions on this endpoint: `project_id`, `rbac_group_id`, `user_id`. Value
	// forms: `project_id` takes a tagged project id (`claude_proj_...`);
	// `rbac_group_id` takes the tagged id (`rbac_group_...`, as emitted in responses
	// and by the spend-limits API) or a bare group UUID, and matches users who held
	// the group at any point during each covered UTC day (time-of-usage attribution);
	// `user_id` takes a tagged user id (`user_...`), as emitted in responses. An
	// unsupported dimension returns 400. At most 100 entries.
	Filter []string `query:"filter,omitzero" json:"-"`
	// Dimensions to break results out by (e.g. `group_by[]=user_id`). Supported on
	// this endpoint: `rbac_group_id`, `user_id`. Grouped rows carry the requested
	// dimension values as additional fields and paginate like ungrouped responses via
	// `next_page`; an unsupported dimension returns 400. `rbac_group_id` attributes a
	// user to every group they held at any point during each covered UTC day, so
	// grouped rows are not an exclusive partition and can sum above org-level totals.
	// At most 100 entries.
	//
	// Any of "rbac_group_id", "user_id".
	GroupBy []string `query:"group_by,omitzero" json:"-"`
	// Sort direction: `asc` or `desc`. Defaults to `asc` for the endpoint's sort
	// column and to `desc` when `order_by` names a metric (a top-N ranking). Applies
	// to `order_by`, or to the endpoint's default sort field when `order_by` is
	// omitted.
	//
	// Any of "asc", "desc".
	Order BetaOrganizationAnalyticsAppChatProjectListParamsOrder `query:"order,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationAnalyticsAppChatProjectListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationAnalyticsAppChatProjectListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort direction: `asc` or `desc`. Defaults to `asc` for the endpoint's sort
// column and to `desc` when `order_by` names a metric (a top-N ranking). Applies
// to `order_by`, or to the endpoint's default sort field when `order_by` is
// omitted.
type BetaOrganizationAnalyticsAppChatProjectListParamsOrder string

const (
	BetaOrganizationAnalyticsAppChatProjectListParamsOrderAsc  BetaOrganizationAnalyticsAppChatProjectListParamsOrder = "asc"
	BetaOrganizationAnalyticsAppChatProjectListParamsOrderDesc BetaOrganizationAnalyticsAppChatProjectListParamsOrder = "desc"
)
