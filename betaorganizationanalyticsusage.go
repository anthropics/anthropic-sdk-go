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
)

// BetaOrganizationAnalyticsUsageService contains methods and other services that
// help with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationAnalyticsUsageService] method instead.
type BetaOrganizationAnalyticsUsageService struct {
	Options []option.RequestOption
}

// NewBetaOrganizationAnalyticsUsageService generates a new service that applies
// the given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationAnalyticsUsageService(opts ...option.RequestOption) (r BetaOrganizationAnalyticsUsageService) {
	r = BetaOrganizationAnalyticsUsageService{}
	r.Options = opts
	return
}

// Get token usage across a date range, broken out by dimension.
//
// Returns one row per combination of the requested `group_by[]` dimensions (and
// time bucket, when `time_bucket_width` is set) that has usage in the range,
// ranked by the `order_by` metric or ordered by time bucket; every combination is
// returned across pages. With no `group_by[]` and no `time_bucket_width`, a single
// row aggregates the whole range. Includes all usage, whether or not it is
// attributable to a seat user. Requires an API key with the `read:analytics`
// scope, or an Admin API key.
func (r *BetaOrganizationAnalyticsUsageService) List(ctx context.Context, params BetaOrganizationAnalyticsUsageListParams, opts ...option.RequestOption) (res *pagination.ReportPageCursor[BetaAnalyticsUsageItem], err error) {
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
	opts = append([]option.RequestOption{option.WithHeader("anthropic-beta", "analytics-usage-cost-2026-08-24"), option.WithResponseInto(&raw)}, opts...)
	path := "v1/organizations/analytics/usage?beta=true"
	opts = append(opts, requestconfig.WithQualifiers(map[string]string{"beta": "true"}))
	opts = append(opts, requestconfig.WithSpan("anthropic.organization.analytics.usage.list", "/v1/organizations/analytics/usage", "actor_ids", "actor_types", "agent_ids", "api_key_ids", "claude_tag_categories", "claude_tag_user_ids", "context_windows", "deployment_ids", "ending_at", "environment_ids", "exclude_deleted_users", "group_by", "inference_geos", "limit", "models", "order", "order_by", "page", "products", "rbac_group_ids", "service_tiers", "session_ids", "slack_channel_ids", "speeds", "starting_at", "time_bucket_width", "user_profile_ids", "workspace_ids"))
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

// Get token usage across a date range, broken out by dimension.
//
// Returns one row per combination of the requested `group_by[]` dimensions (and
// time bucket, when `time_bucket_width` is set) that has usage in the range,
// ranked by the `order_by` metric or ordered by time bucket; every combination is
// returned across pages. With no `group_by[]` and no `time_bucket_width`, a single
// row aggregates the whole range. Includes all usage, whether or not it is
// attributable to a seat user. Requires an API key with the `read:analytics`
// scope, or an Admin API key.
func (r *BetaOrganizationAnalyticsUsageService) ListAutoPaging(ctx context.Context, params BetaOrganizationAnalyticsUsageListParams, opts ...option.RequestOption) *pagination.ReportPageCursorAutoPager[BetaAnalyticsUsageItem] {
	return pagination.NewReportPageCursorAutoPager(r.List(ctx, params, opts...))
}

// The dimension values one usage row is broken out by; each is null unless its
// dimension is in `group_by[]`.
type BetaAnalyticsUsageGroup struct {
	// The identity the usage is attributed to: the user if one is attributed, else the
	// service account, else the federated (cloud-platform) identity, else the API key.
	// Populated only when `actor` is in `group_by[]`; null for usage with no
	// attributable identity, and for users' and federated identities' usage when the
	// reader is a member who can see only some workspaces. Distinct identities the
	// server cannot render also read null, each as its own row.
	Actor BetaAnalyticsActorUnion `json:"actor" api:"required"`
	// Kind of identity the usage is attributed to; matches `actor.type` when `actor`
	// is set. Populated only when `actor_type` is in `group_by[]`.
	//
	// Any of "api_actor", "federated_actor", "service_account_actor", "user_actor".
	ActorType BetaAnalyticsActorType `json:"actor_type" api:"required"`
	// Claude Managed Agents agent whose session incurred the usage. Populated only
	// when `agent_id` is in `group_by[]`; null for usage outside a managed session.
	AgentID string `json:"agent_id" api:"required"`
	// API key that authenticated the requests, as a tagged API key ID, whichever
	// identity `actor` resolves to. Populated only when `api_key_id` is in
	// `group_by[]`; null for usage not authenticated by an API key.
	APIKeyID string `json:"api_key_id" api:"required"`
	// Claude Tag (Claude in Slack) spend category: `engaged` (a person addressed
	// Claude in a channel or thread), `proactive` (Claude responded without being
	// addressed), `scheduled` (a scheduled routine ran), `monitoring` (Claude watching
	// a channel it was asked to monitor), or `dm` (direct messages with Claude).
	// Populated only when `claude_tag_category` is in `group_by[]`; null for usage
	// that is not Claude Tag. Direct-message usage is billed to the individual user
	// and is reported under that user's product, not under `claude_tag`. New
	// categories may be added over time.
	//
	// Any of "dm", "engaged", "monitoring", "proactive", "scheduled".
	ClaudeTagCategory BetaAnalyticsClaudeTagCategory `json:"claude_tag_category" api:"required"`
	// Slack user ID (for example `U0123ABCDEF`) of the member the Claude Tag (Claude
	// in Slack) usage is attributed to, not a claude.ai user ID. Populated only when
	// `claude_tag_user_id` is in `group_by[]`; null for usage that is not Claude Tag
	// and for Claude Tag usage that is not attributed to a single user (for example
	// `monitoring`, and `proactive` usage Claude initiated), so per-user rows can sum
	// to less than the Claude Tag total. Cannot be combined with
	// `group_by[]=rbac_group_id` or the `rbac_group_ids[]` filter.
	ClaudeTagUserID string `json:"claude_tag_user_id" api:"required"`
	// Context-window pricing tier of the usage or cost. Null unless `context_window`
	// is in `group_by[]`; it can also be null on grouped rows with no context-window
	// tier, such as code execution.
	//
	// Any of "0-200k", "200k-1M".
	ContextWindow BetaAnalyticsContextWindow `json:"context_window" api:"required"`
	// Deployment that started the Claude Managed Agents session. Populated only when
	// `deployment_id` is in `group_by[]`; null for usage outside a managed session and
	// for sessions no deployment started.
	DeploymentID string `json:"deployment_id" api:"required"`
	// Claude Managed Agents environment the session ran in. Populated only when
	// `environment_id` is in `group_by[]`; null for usage outside a managed session.
	EnvironmentID string `json:"environment_id" api:"required"`
	// Inference region of the usage or cost. Null unless `inference_geo` is in
	// `group_by[]`; it can also be null on grouped rows where the region is not set
	// (the rows that `inference_geos[]=not_available` matches).
	//
	// Any of "global", "us".
	InferenceGeo BetaAnalyticsUsageCostInferenceGeo `json:"inference_geo" api:"required"`
	// Model that produced the usage or cost, as a model name in the form the
	// `models[]` filter accepts (for example, `claude-opus-5-5`). Null unless `model`
	// is in `group_by[]`; it can also be null on grouped rows whose usage or cost is
	// not attributed to a specific model, such as code execution.
	Model string `json:"model" api:"required"`
	// Product surface that produced the usage or cost. Null unless product is in
	// `group_by[]`; it can also be null on grouped rows whose usage cannot be
	// attributed to a known surface. Values include `chat`, `claude_code`, `cowork`,
	// `office_agent`, `claude_in_chrome`, `claude_design`, `claude_tag`, and
	// `chat_cowork_unified`. `claude_tag` is Claude Tag, the Claude product in Slack.
	// `chat_cowork_unified` is Chat and Cowork unified, Cowork's features inside
	// claude.ai chat: chat and Cowork usage by a member who has it turned on is
	// reported under this value instead of `chat` or `cowork`. It is accepted as a
	// filter only on deployments that offer Chat and Cowork unified. Some unattributed
	// usage is reported as "other".
	Product string `json:"product" api:"required"`
	// RBAC group (team) the usage is attributed to, in the public tagged
	// `rbac_group_...` spelling — the same spelling the activity resources use for
	// this key, so the same team has one id across resources and it round-trips as an
	// `rbac_group_ids[]` filter value. Populated only when `rbac_group_id` is in
	// `group_by[]`. Any-membership semantics: a user in several groups contributes
	// their full usage to each of those groups' rows, so the named-group rows overlap
	// and their sum can exceed the org total. A null value is the single unassigned
	// row: users in no group on that (UTC) day. For the true org total, run the same
	// query without `group_by[]`.
	RBACGroupID string `json:"rbac_group_id" api:"required"`
	// Service tier the usage was served at. Populated only when `service_tier` is in
	// `group_by[]`; null for usage with no service tier (for example, code execution)
	// and for Claude Enterprise usage.
	//
	// Any of "batch", "flex", "flex_discount", "priority", "priority_on_demand",
	// "standard".
	ServiceTier BetaAnalyticsServiceTier `json:"service_tier" api:"required"`
	// Claude Managed Agents session that incurred the usage. Populated only when
	// `session_id` is in `group_by[]`; null for usage outside a managed session.
	SessionID string `json:"session_id" api:"required"`
	// Slack channel the usage originated from. Populated only when `slack_channel_id`
	// is in `group_by[]`; null for usage outside Slack (and for rows recorded before
	// channel attribution was enabled).
	SlackChannelID string `json:"slack_channel_id" api:"required"`
	// Inference speed mode of the usage or cost: `fast` or `standard`. Null unless
	// `speed` is in `group_by[]`.
	//
	// Any of "fast", "standard".
	Speed BetaAnalyticsSpeed `json:"speed" api:"required"`
	// User profile the usage was attributed to, as a tagged `uprof_` ID. Populated
	// only when `user_profile_id` is in `group_by[]`; null for usage that named no
	// profile and for Claude Enterprise usage.
	UserProfileID string `json:"user_profile_id" api:"required"`
	// Workspace the usage occurred in, as a tagged workspace ID, including the
	// organization's default workspace. Populated only when `workspace_id` is in
	// `group_by[]`; null for Claude Enterprise usage.
	WorkspaceID string `json:"workspace_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Actor             respjson.Field
		ActorType         respjson.Field
		AgentID           respjson.Field
		APIKeyID          respjson.Field
		ClaudeTagCategory respjson.Field
		ClaudeTagUserID   respjson.Field
		ContextWindow     respjson.Field
		DeploymentID      respjson.Field
		EnvironmentID     respjson.Field
		InferenceGeo      respjson.Field
		Model             respjson.Field
		Product           respjson.Field
		RBACGroupID       respjson.Field
		ServiceTier       respjson.Field
		SessionID         respjson.Field
		SlackChannelID    respjson.Field
		Speed             respjson.Field
		UserProfileID     respjson.Field
		WorkspaceID       respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsUsageGroup) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsUsageGroup) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAnalyticsUsageGroupBy string

const (
	BetaAnalyticsUsageGroupByActor             BetaAnalyticsUsageGroupBy = "actor"
	BetaAnalyticsUsageGroupByActorType         BetaAnalyticsUsageGroupBy = "actor_type"
	BetaAnalyticsUsageGroupByAgentID           BetaAnalyticsUsageGroupBy = "agent_id"
	BetaAnalyticsUsageGroupByAPIKeyID          BetaAnalyticsUsageGroupBy = "api_key_id"
	BetaAnalyticsUsageGroupByClaudeTagCategory BetaAnalyticsUsageGroupBy = "claude_tag_category"
	BetaAnalyticsUsageGroupByClaudeTagUserID   BetaAnalyticsUsageGroupBy = "claude_tag_user_id"
	BetaAnalyticsUsageGroupByContextWindow     BetaAnalyticsUsageGroupBy = "context_window"
	BetaAnalyticsUsageGroupByDeploymentID      BetaAnalyticsUsageGroupBy = "deployment_id"
	BetaAnalyticsUsageGroupByEnvironmentID     BetaAnalyticsUsageGroupBy = "environment_id"
	BetaAnalyticsUsageGroupByInferenceGeo      BetaAnalyticsUsageGroupBy = "inference_geo"
	BetaAnalyticsUsageGroupByModel             BetaAnalyticsUsageGroupBy = "model"
	BetaAnalyticsUsageGroupByProduct           BetaAnalyticsUsageGroupBy = "product"
	BetaAnalyticsUsageGroupByRBACGroupID       BetaAnalyticsUsageGroupBy = "rbac_group_id"
	BetaAnalyticsUsageGroupByServiceTier       BetaAnalyticsUsageGroupBy = "service_tier"
	BetaAnalyticsUsageGroupBySessionID         BetaAnalyticsUsageGroupBy = "session_id"
	BetaAnalyticsUsageGroupBySlackChannelID    BetaAnalyticsUsageGroupBy = "slack_channel_id"
	BetaAnalyticsUsageGroupBySpeed             BetaAnalyticsUsageGroupBy = "speed"
	BetaAnalyticsUsageGroupByUserProfileID     BetaAnalyticsUsageGroupBy = "user_profile_id"
	BetaAnalyticsUsageGroupByWorkspaceID       BetaAnalyticsUsageGroupBy = "workspace_id"
)

type BetaAnalyticsUsageItem struct {
	// End of the period this row covers (exclusive), as an RFC 3339 UTC timestamp: the
	// end of the row's time bucket when a time-bucket width is set; otherwise the
	// requested `ending_at` (or, when omitted, the earlier of now and `starting_at`
	// plus 31 days) rounded down to the hour, because without a time bucket the report
	// covers whole hours.
	EndingAt time.Time `json:"ending_at" api:"required" format:"date-time"`
	// The dimension values this row's usage is broken out by.
	Group BetaAnalyticsUsageGroup `json:"group" api:"required"`
	// Start of the period this row covers (inclusive), as an RFC 3339 UTC timestamp:
	// the start of the row's time bucket when a time-bucket width is set, otherwise
	// the requested `starting_at` rounded down to the hour.
	StartingAt time.Time `json:"starting_at" api:"required" format:"date-time"`
	// Token usage for this row.
	Usage BetaAnalyticsUsageMetrics `json:"usage" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EndingAt    respjson.Field
		Group       respjson.Field
		StartingAt  respjson.Field
		Usage       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsUsageItem) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsUsageItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The token usage of one row.
type BetaAnalyticsUsageMetrics struct {
	// The number of input tokens for cache creation.
	CacheCreation BetaCacheCreation `json:"cache_creation" api:"required"`
	// The number of input tokens read from the cache.
	CacheReadInputTokens int64 `json:"cache_read_input_tokens" api:"required"`
	// The number of output tokens generated.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// Number of API requests in this row's scope. For sandbox / code-execution events,
	// this counts execution spans rather than HTTP requests (these rows surface with
	// `product: null`).
	Requests int64 `json:"requests" api:"required"`
	// Server-side tool usage metrics.
	ServerToolUse BetaAnalyticsServerToolUse `json:"server_tool_use" api:"required"`
	// Total token count across all token types. This is the value the default
	// `order_by` (`total_tokens`) sorts on.
	TotalTokens int64 `json:"total_tokens" api:"required"`
	// The number of uncached input tokens processed.
	UncachedInputTokens int64 `json:"uncached_input_tokens" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CacheCreation        respjson.Field
		CacheReadInputTokens respjson.Field
		OutputTokens         respjson.Field
		Requests             respjson.Field
		ServerToolUse        respjson.Field
		TotalTokens          respjson.Field
		UncachedInputTokens  respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsUsageMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsUsageMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAnalyticsUsageOrderBy string

const (
	BetaAnalyticsUsageOrderByOutputTokens        BetaAnalyticsUsageOrderBy = "output_tokens"
	BetaAnalyticsUsageOrderByRequests            BetaAnalyticsUsageOrderBy = "requests"
	BetaAnalyticsUsageOrderByTime                BetaAnalyticsUsageOrderBy = "time"
	BetaAnalyticsUsageOrderByTotalTokens         BetaAnalyticsUsageOrderBy = "total_tokens"
	BetaAnalyticsUsageOrderByUncachedInputTokens BetaAnalyticsUsageOrderBy = "uncached_input_tokens"
)

type BetaOrganizationAnalyticsUsageListParams struct {
	// Start of range, inclusive. RFC 3339 tz-aware. Must be within the last 365 days
	// and no earlier than 2026-01-01T00:00:00Z.
	StartingAt time.Time `query:"starting_at" api:"required" format:"date-time" json:"-"`
	// End of range, exclusive. Required when `time_bucket_width` is set. Without
	// `time_bucket_width`, defaults to the earlier of now and `starting_at` + 31 days.
	// The range may span at most 31 days.
	EndingAt param.Opt[time.Time] `query:"ending_at,omitzero" format:"date-time" json:"-"`
	// Opaque cursor from a previous response's `next_page` field.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// If true, omit rows whose `actor` is a `user_actor` with `deleted: true`.
	// Requires `group_by[]=actor`. Rows for other actor types, and rows whose `actor`
	// is null, are always kept. A page may contain fewer than `limit` rows; use
	// `has_more` and `next_page` to paginate as usual.
	ExcludeDeletedUsers param.Opt[bool] `query:"exclude_deleted_users,omitzero" json:"-"`
	// Number of rows per page (1-1000, default 20). `cost_type`/`token_type` fan-out
	// rows (cost endpoint only) do not count toward this limit, so `data` can exceed
	// it.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Filter to usage attributed to specific identities, by tagged ID: a `user_...`,
	// `svac_...` (service account), or `apikey_...` ID. An `apikey_` entry matches
	// only usage attributed to the key itself, not usage on that key attributed to a
	// user or service account (use `api_key_ids[]` for every request a key
	// authenticated). Federated identities have no tagged ID; select them with
	// `actor_types[]`. Use `group_by[]=actor` to break out per-identity values.
	ActorIDs []string `query:"actor_ids,omitzero" json:"-"`
	// Filter to usage attributed to specific kinds of identity. Use
	// `group_by[]=actor_type` to break out per-kind values.
	ActorTypes []BetaAnalyticsActorType `query:"actor_types,omitzero" json:"-"`
	// Filter to usage incurred by sessions of specific Claude Managed Agents agents.
	// Use `group_by[]=agent_id` to break out per-agent values. Claude Developer
	// Platform organizations only.
	AgentIDs []string `query:"agent_ids,omitzero" json:"-"`
	// Filter to usage authenticated by specific API keys, by tagged API key ID
	// (`apikey_...`), whichever identity the usage is attributed to. Include an empty
	// entry (`api_key_ids[]=`) to also match usage that no API key authenticated. Use
	// `group_by[]=api_key_id` to break out per-key values; usage without a key has a
	// null `api_key_id`. Claude Developer Platform organizations only.
	APIKeyIDs []string `query:"api_key_ids,omitzero" json:"-"`
	// Filter to Claude Tag (Claude in Slack) usage in specific spend categories. Usage
	// with no category never matches. `dm` usage is reported under the user's product
	// rather than `claude_tag`, so combining this filter with `products[]=claude_tag`
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
	// Filter to usage incurred by Claude Managed Agents sessions that specific
	// deployments started. Use `group_by[]=deployment_id` to break out per-deployment
	// values. Claude Developer Platform organizations only.
	DeploymentIDs []string `query:"deployment_ids,omitzero" json:"-"`
	// Filter to usage incurred by Claude Managed Agents sessions that ran in specific
	// environments. Use `group_by[]=environment_id` to break out per-environment
	// values. Claude Developer Platform organizations only.
	EnvironmentIDs []string `query:"environment_ids,omitzero" json:"-"`
	// Dimensions to break usage out by (at most 5): one row per combination of their
	// values with usage in the range (and per time bucket, when `time_bucket_width` is
	// set), every combination returned across pages rather than a capped top set. When
	// omitted, a single row aggregates the whole range.
	GroupBy []BetaAnalyticsUsageGroupBy `query:"group_by,omitzero" json:"-"`
	// Filter to specific inference regions. `not_available` matches rows where the
	// region is unset. Use `group_by[]=inference_geo` to break out per-region values.
	InferenceGeos []BetaAnalyticsInferenceGeoFilter `query:"inference_geos,omitzero" json:"-"`
	// Models to include. Defaults to all models. Use `group_by[]=model` to break out
	// per-model values.
	Models []string `query:"models,omitzero" json:"-"`
	// Product surfaces to include. Defaults to all products.
	Products []BetaAnalyticsUsageCostProductFilter `query:"products,omitzero" json:"-"`
	// Filter to usage attributed to specific RBAC groups. Accepts tagged RBAC group
	// IDs (`rbac_group_...`) or bare group UUIDs. A row matches when the user belonged
	// to any of the listed groups on the (UTC) day the usage occurred; usage with no
	// group attribution never matches.
	RBACGroupIDs []string `query:"rbac_group_ids,omitzero" json:"-"`
	// Filter to specific service tiers. Use `group_by[]=service_tier` to break out
	// per-tier values. Claude Developer Platform organizations only; matches nothing
	// for a Claude Enterprise organization.
	ServiceTiers []BetaAnalyticsServiceTier `query:"service_tiers,omitzero" json:"-"`
	// Filter to usage incurred by specific Claude Managed Agents sessions, by tagged
	// session ID (`sesn_...`). Use `group_by[]=session_id` to break out per-session
	// values. Claude Developer Platform organizations only.
	SessionIDs []string `query:"session_ids,omitzero" json:"-"`
	// Filter to usage originating from specific Slack channels. Use
	// `group_by[]=slack_channel_id` to break out per-channel values.
	SlackChannelIDs []string `query:"slack_channel_ids,omitzero" json:"-"`
	// Filter to fast or standard inference mode. Use `group_by[]=speed` to break out
	// per-mode values.
	Speeds []BetaAnalyticsSpeed `query:"speeds,omitzero" json:"-"`
	// Filter to usage attributed to specific user profiles, by tagged user profile ID
	// (`uprof_...`). Use `group_by[]=user_profile_id` to break out per-profile values.
	// Claude Developer Platform organizations only.
	UserProfileIDs []string `query:"user_profile_ids,omitzero" json:"-"`
	// Filter to usage in specific workspaces, by tagged workspace ID (`wrkspc_...`).
	// Use `group_by[]=workspace_id` to break out per-workspace values. Claude
	// Developer Platform organizations only; matches nothing for a Claude Enterprise
	// organization.
	WorkspaceIDs []string `query:"workspace_ids,omitzero" json:"-"`
	// Sort direction of the `order_by` metric, or of the time buckets when `order_by`
	// is `time` and `time_bucket_width` is set. Defaults to `desc`.
	//
	// Any of "asc", "desc".
	Order BetaAnalyticsUsageCostSortOrder `query:"order,omitzero" json:"-"`
	// Row order. A metric ranks rows by that metric in the `order` direction. `time`
	// orders rows by time bucket in the `order` direction (newest first by default)
	// and, within a bucket, by `total_tokens`, highest first. `time` takes effect only
	// with `time_bucket_width`; without it, `time` is the same as `total_tokens`.
	// Defaults to `total_tokens`.
	//
	// Any of "output_tokens", "requests", "time", "total_tokens",
	// "uncached_input_tokens".
	OrderBy BetaAnalyticsUsageOrderBy `query:"order_by,omitzero" json:"-"`
	// Time-bucket granularity. When set, rows are returned per time bucket (one per
	// bucket and dimension combination with usage) and each row's `starting_at` and
	// `ending_at` are that bucket's bounds; the time bucket counts toward `limit`.
	// `ending_at` is required when `time_bucket_width` is set, and with
	// `time_bucket_width="1m"` the range may span at most 24 hours. When omitted, each
	// row aggregates the requested range with `starting_at` and `ending_at` each
	// rounded down to the hour.
	//
	// Any of "1d", "1h", "1m".
	TimeBucketWidth BetaAnalyticsTimeBucketWidth `query:"time_bucket_width,omitzero" json:"-"`
	// This endpoint is in beta: requests must send `analytics-usage-cost-2026-08-24`
	// in this header.
	Betas []AnthropicBeta `header:"anthropic-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [BetaOrganizationAnalyticsUsageListParams]'s query
// parameters as `url.Values`.
func (r BetaOrganizationAnalyticsUsageListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
