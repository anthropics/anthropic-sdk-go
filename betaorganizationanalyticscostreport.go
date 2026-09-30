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

// BetaOrganizationAnalyticsCostReportService contains methods and other services
// that help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationAnalyticsCostReportService] method instead.
type BetaOrganizationAnalyticsCostReportService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationAnalyticsCostReportService generates a new service that
// applies the given options to each request. These options are applied after the
// parent client's options (if there is one), and before any request-specific
// options.
func NewBetaOrganizationAnalyticsCostReportService(opts ...option.RequestOption) (r BetaOrganizationAnalyticsCostReportService) {
	r = BetaOrganizationAnalyticsCostReportService{}
	r.Options = opts
	return
}

// Get cost in USD over time across a date range.
//
// Returns cost bucketed by minute, hour, or day, optionally broken down by
// product, model, context window, inference region, speed, cost type, or token
// type. Available to organizations on a Claude Enterprise plan. Requires an API
// key with the `read:analytics` scope.
func (r *BetaOrganizationAnalyticsCostReportService) List(ctx context.Context, query BetaOrganizationAnalyticsCostReportListParams, opts ...option.RequestOption) (res *pagination.PageCursor[BetaAnalyticsCostReportTimeBucket], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/analytics/cost_report?beta=true"
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

// Get cost in USD over time across a date range.
//
// Returns cost bucketed by minute, hour, or day, optionally broken down by
// product, model, context window, inference region, speed, cost type, or token
// type. Available to organizations on a Claude Enterprise plan. Requires an API
// key with the `read:analytics` scope.
func (r *BetaOrganizationAnalyticsCostReportService) ListAutoPaging(ctx context.Context, query BetaOrganizationAnalyticsCostReportListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[BetaAnalyticsCostReportTimeBucket] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, query, opts...))
}

type BetaOrganizationAnalyticsCostReportListParams struct {
	// Start of range, inclusive. RFC 3339 tz-aware. Must be within the last 365 days
	// and no earlier than 2026-01-01T00:00:00Z.
	StartingAt time.Time `query:"starting_at" api:"required" format:"date-time" json:"-"`
	// End of range, exclusive. When omitted, defaults to the earlier of now and
	// `starting_at` + 31 days. The range may span at most 31 days.
	EndingAt param.Opt[time.Time] `query:"ending_at,omitzero" format:"date-time" json:"-"`
	// Maximum number of time buckets per page. Defaults and caps vary by
	// `bucket_width` (`1d`: default 7, max 31; `1h`: default 24, max 168; `1m`:
	// default 60, max 256).
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page` field.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
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
	// Dimensions to break each time bucket out by. Defaults to no grouping (one total
	// per bucket). Each bucket reports at most its top 100 groups; a group beyond that
	// cap has no row in that bucket (there is no remainder row), so grouped buckets
	// are not exhaustive when a dimension has more than 100 distinct values.
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
	// Product surfaces to include. Defaults to all products. Use `group_by[]=product`
	// to break out per-product values.
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
	// Time bucket granularity.
	//
	// Any of "1d", "1h", "1m".
	BucketWidth BetaOrganizationAnalyticsCostReportListParamsBucketWidth `query:"bucket_width,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationAnalyticsCostReportListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationAnalyticsCostReportListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Time bucket granularity.
type BetaOrganizationAnalyticsCostReportListParamsBucketWidth string

const (
	BetaOrganizationAnalyticsCostReportListParamsBucketWidthDay    BetaOrganizationAnalyticsCostReportListParamsBucketWidth = "1d"
	BetaOrganizationAnalyticsCostReportListParamsBucketWidthHour   BetaOrganizationAnalyticsCostReportListParamsBucketWidth = "1h"
	BetaOrganizationAnalyticsCostReportListParamsBucketWidthMinute BetaOrganizationAnalyticsCostReportListParamsBucketWidth = "1m"
)
