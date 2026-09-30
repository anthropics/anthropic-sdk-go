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

// BetaOrganizationAnalyticsUserCostReportService contains methods and other
// services that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationAnalyticsUserCostReportService] method instead.
type BetaOrganizationAnalyticsUserCostReportService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationAnalyticsUserCostReportService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationAnalyticsUserCostReportService(opts ...option.RequestOption) (r BetaOrganizationAnalyticsUserCostReportService) {
	r = BetaOrganizationAnalyticsUserCostReportService{}
	r.Options = opts
	return
}

// Get per-user cost in USD across a date range.
//
// Returns one row per user, ranked by spend. Use this to see which users account
// for the most cost. Only cost attributable to a seat user is included; for
// organization-wide totals including direct API-key and automation traffic, use
// the bucketed `/v1/organizations/analytics/cost_report` endpoint. Available to
// organizations on a Claude Enterprise plan. Requires an API key with the
// `read:analytics` scope.
func (r *BetaOrganizationAnalyticsUserCostReportService) List(ctx context.Context, query BetaOrganizationAnalyticsUserCostReportListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaAnalyticsCostUsersItem], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/analytics/user_cost_report?beta=true"
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

// Get per-user cost in USD across a date range.
//
// Returns one row per user, ranked by spend. Use this to see which users account
// for the most cost. Only cost attributable to a seat user is included; for
// organization-wide totals including direct API-key and automation traffic, use
// the bucketed `/v1/organizations/analytics/cost_report` endpoint. Available to
// organizations on a Claude Enterprise plan. Requires an API key with the
// `read:analytics` scope.
func (r *BetaOrganizationAnalyticsUserCostReportService) ListAutoPaging(ctx context.Context, query BetaOrganizationAnalyticsUserCostReportListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaAnalyticsCostUsersItem] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

type BetaOrganizationAnalyticsUserCostReportListParams struct {
	// Start of range, inclusive. RFC 3339 tz-aware. Must be within the last 365 days
	// and no earlier than 2026-01-01T00:00:00Z.
	StartingAt time.Time `query:"starting_at" api:"required" format:"date-time" json:"-"`
	// End of range, exclusive. When omitted, defaults to the earlier of now and
	// `starting_at` + 31 days. The range may span at most 31 days.
	EndingAt param.Opt[time.Time] `query:"ending_at,omitzero" format:"date-time" json:"-"`
	// Opaque cursor from a previous response's `next_page` field.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// If true, omit rows for users who are deleted (`deleted: true`). A page may
	// contain fewer than `limit` rows; use `has_more` and `next_page` to paginate as
	// usual.
	ExcludeDeletedUsers param.Opt[bool] `query:"exclude_deleted_users,omitzero" json:"-"`
	// Number of rows per page (1-1000, default 20). One row per actor unless
	// `group_by[]` or `bucket_width` splits an actor across rows;
	// `cost_type`/`token_type` fan-out rows (cost endpoint only) are the exception —
	// they do not count toward this limit, so `data` can exceed it.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Time-bucket granularity. When set, each row's `starting_at` and `ending_at` are
	// populated and one actor may span several rows (one per time bucket with usage).
	// The time bucket counts toward `limit`, so one page can return multiple rows for
	// the same actor. `ending_at` is required when `bucket_width` is set, and with
	// `bucket_width="1m"` the range may span at most 24 hours. When omitted, each row
	// aggregates the full `[starting_at, ending_at)` range.
	//
	// Any of "1d", "1h", "1m".
	BucketWidth BetaOrganizationAnalyticsUserCostReportListParamsBucketWidth `query:"bucket_width,omitzero" json:"-"`
	// Filter to Claude Tag (Claude in Slack) usage in specific spend categories. Usage
	// with no category never matches. `dm` usage is reported under the user's product
	// rather than `claude-tag`, so combining this filter with `products[]=claude-tag`
	// excludes it. Use `group_by[]=claude_tag_category` to break out per-category
	// values.
	ClaudeTagCategories []BetaAnalyticsClaudeTagCategory `query:"claude_tag_categories,omitzero" json:"-"`
	// Filter to Claude Tag (Claude in Slack) usage attributed to specific Slack users,
	// by Slack user ID (for example `U0123ABCDEF`), not claude.ai user ID. Usage that
	// is not Claude Tag, and Claude Tag usage not attributed to a single user, never
	// matches. Use `group_by[]=claude_tag_user_id` to break out per-user values.
	ClaudeTagUserIDs []string `query:"claude_tag_user_ids,omitzero" json:"-"`
	// Filter to specific context-window pricing tiers. Use `group_by[]=context_window`
	// to break out per-tier values.
	ContextWindows []BetaAnalyticsContextWindow `query:"context_windows,omitzero" json:"-"`
	// Break each actor's row out by the given dimensions. Accepts the same values as
	// the bucketed `/cost_report` endpoint. The `product`, `model`, `context_window`,
	// `inference_geo`, and `speed` dimensions — and the time bucket, when
	// `bucket_width` is set — count toward `limit`. `cost_type` and `token_type` do
	// not: `cost_type` returns one row per cost component (tokens, web search, code
	// execution); `token_type` returns one row per token type, each with
	// `cost_type: "tokens"`; combining both returns the per-token-type rows plus the
	// web-search and code-execution rows. A page can therefore contain more rows than
	// `limit` when `cost_type` or `token_type` is requested.
	//
	// Any of "claude_tag_category", "claude_tag_user_id", "context_window",
	// "cost_type", "inference_geo", "model", "product", "rbac_group_id",
	// "slack_channel_id", "speed", "token_type".
	GroupBy []string `query:"group_by,omitzero" json:"-"`
	// Filter to specific inference regions. `not_available` matches rows where the
	// region is unset. Use `group_by[]=inference_geo` to break out per-region values.
	InferenceGeos []BetaAnalyticsInferenceGeoFilter `query:"inference_geos,omitzero" json:"-"`
	// Models to include. Defaults to all models. Use `group_by[]=model` to break out
	// per-model values.
	Models []string `query:"models,omitzero" json:"-"`
	// Product surfaces to include. Defaults to all products.
	Products []BetaAnalyticsProductFilter `query:"products,omitzero" json:"-"`
	// Filter to usage attributed to specific RBAC groups. Accepts tagged RBAC group
	// IDs (`rbac_group_...`) or bare group UUIDs. A row matches when the user belonged
	// to any of the listed groups on the (UTC) day the usage occurred; usage with no
	// group attribution never matches.
	RBACGroupIDs []string `query:"rbac_group_ids,omitzero" json:"-"`
	// Filter to usage originating from specific Slack channels. Use
	// `group_by[]=slack_channel_id` to break out per-channel values.
	SlackChannelIDs []string `query:"slack_channel_ids,omitzero" json:"-"`
	// Filter to fast or standard inference mode. Use `group_by[]=speed` to break out
	// per-mode values.
	//
	// Any of "fast", "standard".
	Speeds []string `query:"speeds,omitzero" json:"-"`
	// Filter to specific users by tagged user ID.
	UserIDs []string `query:"user_ids,omitzero" json:"-"`
	// Sort direction. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaOrganizationAnalyticsUserCostReportListParamsOrder `query:"order,omitzero" json:"-"`
	// Metric to rank actors by. Defaults to `amount`.
	//
	// Any of "amount", "list_amount".
	OrderBy BetaOrganizationAnalyticsUserCostReportListParamsOrderBy `query:"order_by,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationAnalyticsUserCostReportListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationAnalyticsUserCostReportListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Time-bucket granularity. When set, each row's `starting_at` and `ending_at` are
// populated and one actor may span several rows (one per time bucket with usage).
// The time bucket counts toward `limit`, so one page can return multiple rows for
// the same actor. `ending_at` is required when `bucket_width` is set, and with
// `bucket_width="1m"` the range may span at most 24 hours. When omitted, each row
// aggregates the full `[starting_at, ending_at)` range.
type BetaOrganizationAnalyticsUserCostReportListParamsBucketWidth string

const (
	BetaOrganizationAnalyticsUserCostReportListParamsBucketWidthDay    BetaOrganizationAnalyticsUserCostReportListParamsBucketWidth = "1d"
	BetaOrganizationAnalyticsUserCostReportListParamsBucketWidthHour   BetaOrganizationAnalyticsUserCostReportListParamsBucketWidth = "1h"
	BetaOrganizationAnalyticsUserCostReportListParamsBucketWidthMinute BetaOrganizationAnalyticsUserCostReportListParamsBucketWidth = "1m"
)

// Sort direction. Defaults to `desc`.
type BetaOrganizationAnalyticsUserCostReportListParamsOrder string

const (
	BetaOrganizationAnalyticsUserCostReportListParamsOrderAsc  BetaOrganizationAnalyticsUserCostReportListParamsOrder = "asc"
	BetaOrganizationAnalyticsUserCostReportListParamsOrderDesc BetaOrganizationAnalyticsUserCostReportListParamsOrder = "desc"
)

// Metric to rank actors by. Defaults to `amount`.
type BetaOrganizationAnalyticsUserCostReportListParamsOrderBy string

const (
	BetaOrganizationAnalyticsUserCostReportListParamsOrderByAmount     BetaOrganizationAnalyticsUserCostReportListParamsOrderBy = "amount"
	BetaOrganizationAnalyticsUserCostReportListParamsOrderByListAmount BetaOrganizationAnalyticsUserCostReportListParamsOrderBy = "list_amount"
)
