package anthropic

import (
	"encoding/json"
	"time"

	"github.com/anthropics/anthropic-sdk-go/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/shared/constant"
)

// BetaOrganizationAnalyticsService contains methods and other services that help
// with interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewBetaOrganizationAnalyticsService] method instead.
type BetaOrganizationAnalyticsService struct {
	Options         []option.RequestOption
	Summaries       BetaOrganizationAnalyticsSummaryService
	Users           BetaOrganizationAnalyticsUserService
	Apps            BetaOrganizationAnalyticsAppService
	Connectors      BetaOrganizationAnalyticsConnectorService
	Plugins         BetaOrganizationAnalyticsPluginService
	Skills          BetaOrganizationAnalyticsSkillService
	Artifacts       BetaOrganizationAnalyticsArtifactService
	UsageReport     BetaOrganizationAnalyticsUsageReportService
	UserUsageReport BetaOrganizationAnalyticsUserUsageReportService
	CostReport      BetaOrganizationAnalyticsCostReportService
	UserCostReport  BetaOrganizationAnalyticsUserCostReportService
}

// NewBetaOrganizationAnalyticsService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewBetaOrganizationAnalyticsService(opts ...option.RequestOption) (r BetaOrganizationAnalyticsService) {
	r = BetaOrganizationAnalyticsService{}
	r.Options = opts
	r.Summaries = NewBetaOrganizationAnalyticsSummaryService(opts...)
	r.Users = NewBetaOrganizationAnalyticsUserService(opts...)
	r.Apps = NewBetaOrganizationAnalyticsAppService(opts...)
	r.Connectors = NewBetaOrganizationAnalyticsConnectorService(opts...)
	r.Plugins = NewBetaOrganizationAnalyticsPluginService(opts...)
	r.Skills = NewBetaOrganizationAnalyticsSkillService(opts...)
	r.Artifacts = NewBetaOrganizationAnalyticsArtifactService(opts...)
	r.UsageReport = NewBetaOrganizationAnalyticsUsageReportService(opts...)
	r.UserUsageReport = NewBetaOrganizationAnalyticsUserUsageReportService(opts...)
	r.CostReport = NewBetaOrganizationAnalyticsCostReportService(opts...)
	r.UserCostReport = NewBetaOrganizationAnalyticsUserCostReportService(opts...)
	return
}

// Artifact-creation activity for one (`artifact_type`, `is_shared`) bucket on a
// given day.
//
// Artifacts form a small finite cube — the canonical MIME type (8 values incl.
// `other`) crossed with shared-vs-private — so the response is the full set of
// non-empty buckets, not a ranked/paginated list. Claude Code and Cowork artifacts
// report under `text/html` and are counted from 2026-08-17 onward; earlier days
// contain claude.ai chat artifacts only. With `group_by[]=product` / `user_id` /
// `rbac_group_id` each row is further split by the flat group keys and counts are
// scoped to that cut.
type BetaAnalyticsArtifactActivity struct {
	// Canonical artifact MIME type (e.g. `text/markdown`, `application/vnd.ant.react`,
	// `image/svg+xml`), or `other`. Claude Code and Cowork artifacts report as
	// `text/html`.
	ArtifactType string `json:"artifact_type" api:"required"`
	// Number of artifacts created in this bucket on the requested day
	ArtifactsCreatedCount int64 `json:"artifacts_created_count" api:"required"`
	// Number of distinct users who created artifacts in this bucket on the requested
	// day
	DistinctUserCount int64 `json:"distinct_user_count" api:"required"`
	// Whether the artifacts in this bucket have ever been shared (a Claude Code /
	// Cowork artifact is shared once anyone beyond its creator may open it: named
	// members, the whole organization, or anyone with the link).
	IsShared bool `json:"is_shared" api:"required"`
	// Number of those artifacts that have been published (for Claude Code / Cowork
	// artifacts: open to anyone with the link); never exceeds
	// `artifacts_created_count`
	PublishedArtifactsCreatedCount int64 `json:"published_artifacts_created_count" api:"required"`
	// Product that produced this row's activity: one of `chat`, `claude_code`,
	// `cowork`, `office_agent`, or `chat_cowork_unified` (Chat and Cowork unified).
	// These are the canonical Cost & Usage product names; an `office_agent` row's
	// per-surface breakdown is in its `office_metrics`. On `/plugins` only `cowork`,
	// `claude_code` and `chat_cowork_unified` occur (the only surfaces with plugin
	// attribution); on `/artifacts` only `chat`, `claude_code`, `cowork` and
	// `chat_cowork_unified` occur (the surfaces that create artifacts);
	// `/apps/chat/projects` does not support the product dimension (a `product` entry
	// in `group_by[]` or `filter[]` there is rejected). Present only when the request
	// grouped by `product`.
	Product string `json:"product" api:"nullable"`
	// Tagged RBAC group identifier (`rbac_group_...`), matching the spend-limits API
	// spelling. Present only when the request grouped by `rbac_group_id`.
	RBACGroupID string `json:"rbac_group_id" api:"nullable"`
	// Resolved RBAC group display name, alongside `rbac_group_id` when name resolution
	// is available. Null if the group has been deleted or its name could not be
	// resolved; `rbac_group_id` remains the stable key.
	RBACGroupName string `json:"rbac_group_name" api:"nullable"`
	// Tagged user identifier (e.g. `user_...`). Present only when the request grouped
	// by `user_id`.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ArtifactType                   respjson.Field
		ArtifactsCreatedCount          respjson.Field
		DistinctUserCount              respjson.Field
		IsShared                       respjson.Field
		PublishedArtifactsCreatedCount respjson.Field
		Product                        respjson.Field
		RBACGroupID                    respjson.Field
		RBACGroupName                  respjson.Field
		UserID                         respjson.Field
		ExtraFields                    map[string]respjson.Field
		raw                            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsArtifactActivity) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsArtifactActivity) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Chat activity recorded while members had Chat and Cowork unified turned on.
type BetaAnalyticsChatCoworkUnifiedChatMetrics struct {
	// Same measure as `chat_metrics.connectors_used_count`, for activity recorded
	// while members had Chat and Cowork unified turned on.
	ConnectorsUsedCount int64 `json:"connectors_used_count" api:"required"`
	// Same measure as `chat_metrics.distinct_artifacts_created_count`, for activity
	// recorded while members had Chat and Cowork unified turned on. Exact in
	// date-range mode: a creation belongs to exactly one day, so the per-day counts
	// never overlap and their sum over the window is the exact count of distinct
	// creations in it.
	DistinctArtifactsCreatedCount int64 `json:"distinct_artifacts_created_count" api:"required"`
	// Same measure as `chat_metrics.distinct_connectors_used_count`, for activity
	// recorded while members had Chat and Cowork unified turned on. Approximate (HLL,
	// typical error <2%) in date-range mode. Null on aggregated rows where a distinct
	// count cannot be computed.
	DistinctConnectorsUsedCount int64 `json:"distinct_connectors_used_count" api:"required"`
	// Same measure as `chat_metrics.distinct_conversation_count`, for activity
	// recorded while members had Chat and Cowork unified turned on. Approximate (HLL,
	// typical error <2%) in date-range mode. Null on aggregated rows where a distinct
	// count cannot be computed.
	DistinctConversationCount int64 `json:"distinct_conversation_count" api:"required"`
	// Same measure as `chat_metrics.distinct_files_uploaded_count`, for activity
	// recorded while members had Chat and Cowork unified turned on. It counts uploaded
	// files as well as files Claude created and images returned by Claude's tools,
	// such as screenshots. Approximate (HLL, typical error <2%) in date-range mode.
	// Null on aggregated rows where a distinct count cannot be computed.
	DistinctFilesUploadedCount int64 `json:"distinct_files_uploaded_count" api:"required"`
	// Same measure as `chat_metrics.distinct_projects_created_count`, for activity
	// recorded while members had Chat and Cowork unified turned on. Exact in
	// date-range mode: a creation belongs to exactly one day, so the per-day counts
	// never overlap and their sum over the window is the exact count of distinct
	// creations in it.
	DistinctProjectsCreatedCount int64 `json:"distinct_projects_created_count" api:"required"`
	// Same measure as `chat_metrics.distinct_projects_used_count`, for activity
	// recorded while members had Chat and Cowork unified turned on. Approximate (HLL,
	// typical error <2%) in date-range mode. Null on aggregated rows where a distinct
	// count cannot be computed.
	DistinctProjectsUsedCount int64 `json:"distinct_projects_used_count" api:"required"`
	// Always null: shared-artifact views are not currently measured.
	DistinctSharedArtifactsViewedCount int64 `json:"distinct_shared_artifacts_viewed_count" api:"required"`
	// Same measure as `chat_metrics.distinct_skills_used_count`, for activity recorded
	// while members had Chat and Cowork unified turned on. Approximate (HLL, typical
	// error <2%) in date-range mode. Null on aggregated rows where a distinct count
	// cannot be computed.
	DistinctSkillsUsedCount int64 `json:"distinct_skills_used_count" api:"required"`
	// Same measure as `chat_metrics.message_count`, for activity recorded while
	// members had Chat and Cowork unified turned on.
	MessageCount int64 `json:"message_count" api:"required"`
	// Same measure as `chat_metrics.shared_conversations_viewed_count`, for activity
	// recorded while members had Chat and Cowork unified turned on.
	SharedConversationsViewedCount int64 `json:"shared_conversations_viewed_count" api:"required"`
	// Same measure as `chat_metrics.thinking_message_count`, for activity recorded
	// while members had Chat and Cowork unified turned on.
	ThinkingMessageCount int64 `json:"thinking_message_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ConnectorsUsedCount                respjson.Field
		DistinctArtifactsCreatedCount      respjson.Field
		DistinctConnectorsUsedCount        respjson.Field
		DistinctConversationCount          respjson.Field
		DistinctFilesUploadedCount         respjson.Field
		DistinctProjectsCreatedCount       respjson.Field
		DistinctProjectsUsedCount          respjson.Field
		DistinctSharedArtifactsViewedCount respjson.Field
		DistinctSkillsUsedCount            respjson.Field
		MessageCount                       respjson.Field
		SharedConversationsViewedCount     respjson.Field
		ThinkingMessageCount               respjson.Field
		ExtraFields                        map[string]respjson.Field
		raw                                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsChatCoworkUnifiedChatMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsChatCoworkUnifiedChatMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Cowork session activity recorded while members had Chat and Cowork unified
// turned on.
type BetaAnalyticsChatCoworkUnifiedSessionsMetrics struct {
	// Same measure as `cowork_metrics.action_count`, for activity recorded while
	// members had Chat and Cowork unified turned on.
	ActionCount int64 `json:"action_count" api:"required"`
	// Same measure as `cowork_metrics.artifacts_created_count`, for activity recorded
	// while members had Chat and Cowork unified turned on. Exact in date-range mode: a
	// creation belongs to exactly one day, so the per-day counts never overlap and
	// their sum over the window is the exact count of distinct creations in it.
	ArtifactsCreatedCount int64 `json:"artifacts_created_count" api:"required"`
	// Same measure as `cowork_metrics.connectors_used_count`, for activity recorded
	// while members had Chat and Cowork unified turned on.
	ConnectorsUsedCount int64 `json:"connectors_used_count" api:"required"`
	// Same measure as `cowork_metrics.dispatch_turn_count`, for activity recorded
	// while members had Chat and Cowork unified turned on.
	DispatchTurnCount int64 `json:"dispatch_turn_count" api:"required"`
	// Same measure as `cowork_metrics.distinct_connectors_used_count`, for activity
	// recorded while members had Chat and Cowork unified turned on. Approximate (HLL,
	// typical error <2%) in date-range mode. Null on aggregated rows where a distinct
	// count cannot be computed.
	DistinctConnectorsUsedCount int64 `json:"distinct_connectors_used_count" api:"required"`
	// Same measure as `cowork_metrics.distinct_plugins_used_count`, for activity
	// recorded while members had Chat and Cowork unified turned on. Approximate (HLL,
	// typical error <2%) in date-range mode. Null on aggregated rows where a distinct
	// count cannot be computed.
	DistinctPluginsUsedCount int64 `json:"distinct_plugins_used_count" api:"required"`
	// Same measure as `cowork_metrics.distinct_session_count`, for activity recorded
	// while members had Chat and Cowork unified turned on. Approximate (HLL, typical
	// error <2%) in date-range mode. Null on aggregated rows where a distinct count
	// cannot be computed.
	DistinctSessionCount int64 `json:"distinct_session_count" api:"required"`
	// Same measure as `cowork_metrics.distinct_skills_used_count`, for activity
	// recorded while members had Chat and Cowork unified turned on. Approximate (HLL,
	// typical error <2%) in date-range mode. Null on aggregated rows where a distinct
	// count cannot be computed.
	DistinctSkillsUsedCount int64 `json:"distinct_skills_used_count" api:"required"`
	// Same measure as `cowork_metrics.edit_tool_count`, for activity recorded while
	// members had Chat and Cowork unified turned on.
	EditToolCount int64 `json:"edit_tool_count" api:"required"`
	// Same measure as `cowork_metrics.file_edit_count`, for activity recorded while
	// members had Chat and Cowork unified turned on.
	FileEditCount int64 `json:"file_edit_count" api:"required"`
	// Same measure as `cowork_metrics.message_count`, for activity recorded while
	// members had Chat and Cowork unified turned on.
	MessageCount int64 `json:"message_count" api:"required"`
	// Same measure as `cowork_metrics.multi_edit_tool_count`, for activity recorded
	// while members had Chat and Cowork unified turned on. Claude no longer has a
	// multi-edit tool, so expect 0 when not null; each edit is now a separate Edit
	// tool call, counted in `edit_tool_count` and `file_edit_count`.
	MultiEditToolCount int64 `json:"multi_edit_tool_count" api:"required"`
	// Same measure as `cowork_metrics.notebook_edit_tool_count`, for activity recorded
	// while members had Chat and Cowork unified turned on.
	NotebookEditToolCount int64 `json:"notebook_edit_tool_count" api:"required"`
	// Same measure as `cowork_metrics.plugins_used_count`, for activity recorded while
	// members had Chat and Cowork unified turned on.
	PluginsUsedCount int64 `json:"plugins_used_count" api:"required"`
	// Same measure as `cowork_metrics.sessions_with_file_edits_count`, for activity
	// recorded while members had Chat and Cowork unified turned on. Approximate (HLL,
	// typical error <2%) in date-range mode. Null on aggregated rows where a distinct
	// count cannot be computed.
	SessionsWithFileEditsCount int64 `json:"sessions_with_file_edits_count" api:"required"`
	// Same measure as `cowork_metrics.skills_used_count`, for activity recorded while
	// members had Chat and Cowork unified turned on.
	SkillsUsedCount int64 `json:"skills_used_count" api:"required"`
	// Same measure as `cowork_metrics.write_tool_count`, for activity recorded while
	// members had Chat and Cowork unified turned on.
	WriteToolCount int64 `json:"write_tool_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ActionCount                 respjson.Field
		ArtifactsCreatedCount       respjson.Field
		ConnectorsUsedCount         respjson.Field
		DispatchTurnCount           respjson.Field
		DistinctConnectorsUsedCount respjson.Field
		DistinctPluginsUsedCount    respjson.Field
		DistinctSessionCount        respjson.Field
		DistinctSkillsUsedCount     respjson.Field
		EditToolCount               respjson.Field
		FileEditCount               respjson.Field
		MessageCount                respjson.Field
		MultiEditToolCount          respjson.Field
		NotebookEditToolCount       respjson.Field
		PluginsUsedCount            respjson.Field
		SessionsWithFileEditsCount  respjson.Field
		SkillsUsedCount             respjson.Field
		WriteToolCount              respjson.Field
		ExtraFields                 map[string]respjson.Field
		raw                         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsChatCoworkUnifiedSessionsMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsChatCoworkUnifiedSessionsMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude.ai activity metrics for a single user on a given day.
type BetaAnalyticsChatMetrics struct {
	// Number of MCP connector invocations.
	ConnectorsUsedCount int64 `json:"connectors_used_count" api:"required"`
	// Number of distinct artifacts created. Exact in date-range mode: a creation
	// belongs to exactly one day, so the per-day counts never overlap and their sum
	// over the window is the exact count of distinct creations in it.
	DistinctArtifactsCreatedCount int64 `json:"distinct_artifacts_created_count" api:"required"`
	// Distinct claude.ai connectors this user used. Excludes calls whose connector
	// could not be identified and all calls from organizations with zero data
	// retention. Approximate (HLL, typical error <2%) in date-range mode. Null on
	// aggregated rows where a distinct count cannot be computed.
	DistinctConnectorsUsedCount int64 `json:"distinct_connectors_used_count" api:"required"`
	// Number of distinct conversations the user participated in. Approximate (HLL,
	// typical error <2%) in date-range mode. Null on aggregated rows where a distinct
	// count cannot be computed.
	DistinctConversationCount int64 `json:"distinct_conversation_count" api:"required"`
	// Number of distinct files uploaded. Approximate (HLL, typical error <2%) in
	// date-range mode. Null on aggregated rows where a distinct count cannot be
	// computed.
	DistinctFilesUploadedCount int64 `json:"distinct_files_uploaded_count" api:"required"`
	// Number of distinct projects created. Exact in date-range mode: a creation
	// belongs to exactly one day, so the per-day counts never overlap and their sum
	// over the window is the exact count of distinct creations in it.
	DistinctProjectsCreatedCount int64 `json:"distinct_projects_created_count" api:"required"`
	// Number of distinct projects used. Approximate (HLL, typical error <2%) in
	// date-range mode. Null on aggregated rows where a distinct count cannot be
	// computed.
	DistinctProjectsUsedCount int64 `json:"distinct_projects_used_count" api:"required"`
	// Number of distinct shared artifacts the user viewed. Approximate (HLL, typical
	// error <2%) in date-range mode. Null on aggregated rows where a distinct count
	// cannot be computed.
	DistinctSharedArtifactsViewedCount int64 `json:"distinct_shared_artifacts_viewed_count" api:"required"`
	// Number of distinct skills used. Approximate (HLL, typical error <2%) in
	// date-range mode. Null on aggregated rows where a distinct count cannot be
	// computed.
	DistinctSkillsUsedCount int64 `json:"distinct_skills_used_count" api:"required"`
	// Number of messages sent
	MessageCount int64 `json:"message_count" api:"required"`
	// Number of times the user opened a shared conversation in a project
	SharedConversationsViewedCount int64 `json:"shared_conversations_viewed_count" api:"required"`
	// Number of messages that used extended thinking
	ThinkingMessageCount int64 `json:"thinking_message_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ConnectorsUsedCount                respjson.Field
		DistinctArtifactsCreatedCount      respjson.Field
		DistinctConnectorsUsedCount        respjson.Field
		DistinctConversationCount          respjson.Field
		DistinctFilesUploadedCount         respjson.Field
		DistinctProjectsCreatedCount       respjson.Field
		DistinctProjectsUsedCount          respjson.Field
		DistinctSharedArtifactsViewedCount respjson.Field
		DistinctSkillsUsedCount            respjson.Field
		MessageCount                       respjson.Field
		SharedConversationsViewedCount     respjson.Field
		ThinkingMessageCount               respjson.Field
		ExtraFields                        map[string]respjson.Field
		raw                                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsChatMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsChatMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude Code activity metrics for a single user on a given day.
type BetaAnalyticsClaudeCodeMetrics struct {
	// Core Claude Code activity metrics for a single user on a given day.
	CoreMetrics BetaAnalyticsCoreCodeMetrics `json:"core_metrics" api:"required"`
	// Per-tool accepted/rejected counts for Claude Code file modification tools.
	ToolActions BetaAnalyticsToolActions `json:"tool_actions" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CoreMetrics respjson.Field
		ToolActions respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsClaudeCodeMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsClaudeCodeMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAnalyticsClaudeTagCategory string

const (
	BetaAnalyticsClaudeTagCategoryDm         BetaAnalyticsClaudeTagCategory = "dm"
	BetaAnalyticsClaudeTagCategoryEngaged    BetaAnalyticsClaudeTagCategory = "engaged"
	BetaAnalyticsClaudeTagCategoryMonitoring BetaAnalyticsClaudeTagCategory = "monitoring"
	BetaAnalyticsClaudeTagCategoryProactive  BetaAnalyticsClaudeTagCategory = "proactive"
	BetaAnalyticsClaudeTagCategoryScheduled  BetaAnalyticsClaudeTagCategory = "scheduled"
)

// Per-connector activity data for a given day.
type BetaAnalyticsConnectorActivity struct {
	// Claude.ai activity metrics for a single connector on a given day.
	ChatMetrics BetaAnalyticsConnectorChatMetrics `json:"chat_metrics" api:"required"`
	// Claude Code activity metrics for a single connector on a given day.
	ClaudeCodeMetrics BetaAnalyticsConnectorClaudeCodeMetrics `json:"claude_code_metrics" api:"required"`
	// Name of the connector. Some rows carry an opaque connector id here instead of a
	// readable name; `connector_display_name` holds the resolved name for those rows.
	ConnectorName string `json:"connector_name" api:"required"`
	// Cowork activity metrics for a single connector on a given day.
	CoworkMetrics BetaAnalyticsConnectorCoworkMetrics `json:"cowork_metrics" api:"required"`
	// Number of distinct users who used the connector on the requested day, or, in
	// date-range mode, over the requested window — recomputed as an exact distinct
	// count over the window's per-member daily rows, never a sum of per-day values.
	DistinctUserCount int64 `json:"distinct_user_count" api:"required"`
	// Office Agent activity metrics for a single connector on a given day, broken out
	// by Office product.
	OfficeMetrics BetaAnalyticsConnectorOfficeMetrics `json:"office_metrics" api:"required"`
	// Connector use recorded while members had Chat and Cowork unified (Cowork's
	// features inside claude.ai chat) turned on, split into chat conversations and
	// Cowork sessions. A count is null in date-range mode where it cannot be computed.
	// Omitted from the response on deployments that do not offer Chat and Cowork
	// unified.
	ChatCoworkUnifiedMetrics BetaAnalyticsConnectorActivityChatCoworkUnifiedMetrics `json:"chat_cowork_unified_metrics" api:"nullable"`
	// Human-readable display name for rows whose `connector_name` is an opaque
	// connector id rather than a readable name, resolved at request time from the
	// organization's connectors (including connectors that have since been removed).
	// `connector_name` remains the row's stable key for sorting and pagination, and
	// `filter[]=connector_name:{value}` also matches these rows by display name.
	// Display names are not unique, and the same connector's claude.ai usage can
	// appear under a separate row with a readable `connector_name`. Null when
	// `connector_name` is already a readable name, when the id cannot be resolved to
	// one of the organization's connectors, or when display-name resolution is not
	// enabled for this organization.
	ConnectorDisplayName string `json:"connector_display_name" api:"nullable"`
	// Number of distinct users whose use of this connector on the requested day ran on
	// their own individual credential, connected through their own consent flow.
	// Companion bucket to `managed_auth_distinct_user_count`, which carries the
	// measurement, attribution, and null rules. Users whose requests used no stored
	// credential count in neither bucket.
	IndividualAuthDistinctUserCount int64 `json:"individual_auth_distinct_user_count" api:"nullable"`
	// Number of distinct users whose use of this connector on the requested day ran on
	// Enterprise Managed Auth (an organization-managed credential provisioned through
	// the organization's identity provider), read from the token record each request
	// used. Null, never 0, when managed-auth reporting is not enabled for the
	// organization, the value cannot be attributed to the row, no credentialed
	// requests and no managed-token mint events (a managed credential being
	// provisioned for a user's use of the connector) were observed that day, or the
	// day predates 2026-07-01, the first day the backing data exists (forward-only
	// data, no backfill). When credentialed requests or mint events were observed and
	// attributed, both managed-auth fields populate, reporting 0 for a bucket with no
	// users; the two counts are independent, not a partition — a user whose requests
	// that day used both kinds of credential counts in both. Mint events carry user
	// but not surface attribution, so they count as observed auth activity on
	// `user_id` and `rbac_group_id` cuts — attributed to the user the credential was
	// provisioned for — but never on a cut that references `product` (group or
	// filter). Date-range rollup mode (`starting_date`/`ending_date`) computes both
	// fields exactly over the window — distinct users with at least one qualifying day
	// — when the whole window starts on or after 2026-07-01, with the null-versus-0
	// and mint-event rules applying with the window in place of the day; a range
	// starting earlier reports every managed-auth field as null, never a
	// partial-window value.
	ManagedAuthDistinctUserCount int64 `json:"managed_auth_distinct_user_count" api:"nullable"`
	// Product that produced this row's activity: one of `chat`, `claude_code`,
	// `cowork`, `office_agent`, or `chat_cowork_unified` (Chat and Cowork unified).
	// These are the canonical Cost & Usage product names; an `office_agent` row's
	// per-surface breakdown is in its `office_metrics`. On `/plugins` only `cowork`,
	// `claude_code` and `chat_cowork_unified` occur (the only surfaces with plugin
	// attribution); on `/artifacts` only `chat`, `claude_code`, `cowork` and
	// `chat_cowork_unified` occur (the surfaces that create artifacts);
	// `/apps/chat/projects` does not support the product dimension (a `product` entry
	// in `group_by[]` or `filter[]` there is rejected). Present only when the request
	// grouped by `product`.
	Product string `json:"product" api:"nullable"`
	// Tagged RBAC group identifier (`rbac_group_...`), matching the spend-limits API
	// spelling. Present only when the request grouped by `rbac_group_id`.
	RBACGroupID string `json:"rbac_group_id" api:"nullable"`
	// Resolved RBAC group display name, alongside `rbac_group_id` when name resolution
	// is available. Null if the group has been deleted or its name could not be
	// resolved; `rbac_group_id` remains the stable key.
	RBACGroupName string `json:"rbac_group_name" api:"nullable"`
	// Number of connector tool calls on the requested day whose trusted read-only
	// annotation marked them read-only. Call count, not distinct users. Every call
	// recorded on a classified surface lands in exactly one of `read_call_count`,
	// `write_call_count`, or `unclassified_call_count`, so the three sum to the day's
	// classified calls. Classification is forward-only per surface: claude.ai from
	// 2026-06-01, Claude Code from 2026-05-30, Claude in Office from 2026-05-29,
	// Cowork from 2026-06-02 (Cowork clients predating annotation forwarding land in
	// `unclassified_call_count`). Null, never 0, when the value cannot be stated: the
	// read/write split is not enabled for this organization, or the day predates
	// 2026-05-29. For a date-range total, sum the per-day values, but treat a window
	// that extends before 2026-05-29 as null rather than summing only its covered days
	// — date-range rollup mode (`starting_date`/`ending_date`) applies both rules
	// server-side.
	ReadCallCount int64 `json:"read_call_count" api:"nullable"`
	// Number of connector tool calls on the requested day with no trusted read-only
	// annotation — the annotation is optional in the MCP spec and is discarded when
	// connector access controls are active, so unclassified calls are common. This
	// field shows how much of the day's classified activity the read/write split
	// actually covers. Call count, not distinct users. One of the three
	// call-classification buckets; see `read_call_count` for the per-surface
	// data-start dates, null conditions, and date-range guidance.
	UnclassifiedCallCount int64 `json:"unclassified_call_count" api:"nullable"`
	// Tagged user identifier (e.g. `user_...`). Present only when the request grouped
	// by `user_id`.
	UserID string `json:"user_id" api:"nullable"`
	// Number of connector tool calls on the requested day whose trusted read-only
	// annotation marked them not read-only. Call count, not distinct users. One of the
	// three call-classification buckets; see `read_call_count` for the per-surface
	// data-start dates, null conditions, and date-range guidance.
	WriteCallCount int64 `json:"write_call_count" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ChatMetrics                     respjson.Field
		ClaudeCodeMetrics               respjson.Field
		ConnectorName                   respjson.Field
		CoworkMetrics                   respjson.Field
		DistinctUserCount               respjson.Field
		OfficeMetrics                   respjson.Field
		ChatCoworkUnifiedMetrics        respjson.Field
		ConnectorDisplayName            respjson.Field
		IndividualAuthDistinctUserCount respjson.Field
		ManagedAuthDistinctUserCount    respjson.Field
		Product                         respjson.Field
		RBACGroupID                     respjson.Field
		RBACGroupName                   respjson.Field
		ReadCallCount                   respjson.Field
		UnclassifiedCallCount           respjson.Field
		UserID                          respjson.Field
		WriteCallCount                  respjson.Field
		ExtraFields                     map[string]respjson.Field
		raw                             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsConnectorActivity) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsConnectorActivity) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Connector use recorded while members had Chat and Cowork unified (Cowork's
// features inside claude.ai chat) turned on, split into chat conversations and
// Cowork sessions. A count is null in date-range mode where it cannot be computed.
// Omitted from the response on deployments that do not offer Chat and Cowork
// unified.
type BetaAnalyticsConnectorActivityChatCoworkUnifiedMetrics struct {
	// A connector's use in chat conversations recorded while members had Chat and
	// Cowork unified turned on.
	Chat BetaAnalyticsConnectorChatCoworkUnifiedChatMetrics `json:"chat" api:"required"`
	// A connector's use in Cowork sessions recorded while members had Chat and Cowork
	// unified turned on.
	Sessions BetaAnalyticsConnectorChatCoworkUnifiedSessionsMetrics `json:"sessions" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Chat        respjson.Field
		Sessions    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsConnectorActivityChatCoworkUnifiedMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsConnectorActivityChatCoworkUnifiedMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A connector's use in chat conversations recorded while members had Chat and
// Cowork unified turned on.
type BetaAnalyticsConnectorChatCoworkUnifiedChatMetrics struct {
	// Same measure as `chat_metrics.distinct_conversation_connector_used_count`, for
	// activity recorded while members had Chat and Cowork unified turned on.
	// Approximate (HLL, typical error <2%) in date-range mode. Null on aggregated rows
	// where a distinct count cannot be computed.
	DistinctConversationConnectorUsedCount int64 `json:"distinct_conversation_connector_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctConversationConnectorUsedCount respjson.Field
		ExtraFields                            map[string]respjson.Field
		raw                                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsConnectorChatCoworkUnifiedChatMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsConnectorChatCoworkUnifiedChatMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A connector's use in Cowork sessions recorded while members had Chat and Cowork
// unified turned on.
type BetaAnalyticsConnectorChatCoworkUnifiedSessionsMetrics struct {
	// Same measure as `cowork_metrics.distinct_session_connector_used_count`, for
	// activity recorded while members had Chat and Cowork unified turned on.
	// Approximate (HLL, typical error <2%) in date-range mode. Null on aggregated rows
	// where a distinct count cannot be computed.
	DistinctSessionConnectorUsedCount int64 `json:"distinct_session_connector_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctSessionConnectorUsedCount respjson.Field
		ExtraFields                       map[string]respjson.Field
		raw                               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsConnectorChatCoworkUnifiedSessionsMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsConnectorChatCoworkUnifiedSessionsMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude.ai activity metrics for a single connector on a given day.
type BetaAnalyticsConnectorChatMetrics struct {
	// Number of distinct conversations in which the connector was used. Approximate
	// (HLL, typical error <2%) in date-range mode. Null on aggregated rows where a
	// distinct count cannot be computed.
	DistinctConversationConnectorUsedCount int64 `json:"distinct_conversation_connector_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctConversationConnectorUsedCount respjson.Field
		ExtraFields                            map[string]respjson.Field
		raw                                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsConnectorChatMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsConnectorChatMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude Code activity metrics for a single connector on a given day.
type BetaAnalyticsConnectorClaudeCodeMetrics struct {
	// Number of distinct Claude Code sessions in which the connector was used.
	// Approximate (HLL, typical error <2%) in date-range mode. Null on aggregated rows
	// where a distinct count cannot be computed.
	DistinctSessionConnectorUsedCount int64 `json:"distinct_session_connector_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctSessionConnectorUsedCount respjson.Field
		ExtraFields                       map[string]respjson.Field
		raw                               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsConnectorClaudeCodeMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsConnectorClaudeCodeMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Cowork activity metrics for a single connector on a given day.
type BetaAnalyticsConnectorCoworkMetrics struct {
	// Number of distinct Cowork sessions in which the connector was used. Approximate
	// (HLL, typical error <2%) in date-range mode. Null on aggregated rows where a
	// distinct count cannot be computed.
	DistinctSessionConnectorUsedCount int64 `json:"distinct_session_connector_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctSessionConnectorUsedCount respjson.Field
		ExtraFields                       map[string]respjson.Field
		raw                               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsConnectorCoworkMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsConnectorCoworkMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Office Agent activity metrics for a single connector on a given day, broken out
// by Office product.
type BetaAnalyticsConnectorOfficeMetrics struct {
	// Office Agent activity metrics for a single connector on a given day within one
	// Office product.
	Excel BetaAnalyticsConnectorOfficeProductMetrics `json:"excel" api:"required"`
	// Office Agent activity metrics for a single connector on a given day within one
	// Office product.
	Outlook BetaAnalyticsConnectorOfficeProductMetrics `json:"outlook" api:"required"`
	// Office Agent activity metrics for a single connector on a given day within one
	// Office product.
	Powerpoint BetaAnalyticsConnectorOfficeProductMetrics `json:"powerpoint" api:"required"`
	// Office Agent activity metrics for a single connector on a given day within one
	// Office product.
	Word BetaAnalyticsConnectorOfficeProductMetrics `json:"word" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Excel       respjson.Field
		Outlook     respjson.Field
		Powerpoint  respjson.Field
		Word        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsConnectorOfficeMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsConnectorOfficeMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Office Agent activity metrics for a single connector on a given day within one
// Office product.
type BetaAnalyticsConnectorOfficeProductMetrics struct {
	// Number of distinct Office Agent sessions in which the connector was used.
	// Approximate (HLL, typical error <2%) in date-range mode. Null on aggregated rows
	// where a distinct count cannot be computed.
	DistinctSessionConnectorUsedCount int64 `json:"distinct_session_connector_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctSessionConnectorUsedCount respjson.Field
		ExtraFields                       map[string]respjson.Field
		raw                               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsConnectorOfficeProductMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsConnectorOfficeProductMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAnalyticsContextWindow string

const (
	BetaAnalyticsContextWindowFrom0To200k  BetaAnalyticsContextWindow = "0-200k"
	BetaAnalyticsContextWindowFrom200kTo1M BetaAnalyticsContextWindow = "200k-1M"
)

// Core Claude Code activity metrics for a single user on a given day.
type BetaAnalyticsCoreCodeMetrics struct {
	// Number of artifacts created in Claude Code sessions: an artifact counts once, on
	// the day a session first saves it. Counted from 2026-08-17; 0 on earlier days.
	// Exact in date-range mode: a creation belongs to exactly one day, so the per-day
	// counts never overlap and their sum over the window is the exact count of
	// distinct creations in it.
	ArtifactsCreatedCount int64 `json:"artifacts_created_count" api:"required"`
	// Number of commits made via Claude Code
	CommitCount int64 `json:"commit_count" api:"required"`
	// Number of distinct Claude Code sessions. On aggregated rows and in date-range
	// mode: summed per-day distinct counts. A session essentially never spans a UTC
	// day, so the sum is in practice the true distinct count.
	DistinctSessionCount int64 `json:"distinct_session_count" api:"required"`
	// Lines of code added and removed via Claude Code.
	LinesOfCode BetaAnalyticsLinesOfCode `json:"lines_of_code" api:"required"`
	// Number of pull requests created via Claude Code
	PullRequestCount int64 `json:"pull_request_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ArtifactsCreatedCount respjson.Field
		CommitCount           respjson.Field
		DistinctSessionCount  respjson.Field
		LinesOfCode           respjson.Field
		PullRequestCount      respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsCoreCodeMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsCoreCodeMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAnalyticsCostBucketedResult struct {
	// Amount (post-discount, pre-credit) in fractional cents.
	Amount string `json:"amount" api:"required"`
	// Claude Tag (Claude in Slack) spend category: `engaged` (a person addressed
	// Claude in a channel or thread), `proactive` (Claude responded without being
	// addressed), `scheduled` (a scheduled routine ran), `monitoring` (Claude watching
	// a channel it was asked to monitor), or `dm` (direct messages with Claude).
	// Populated only when `claude_tag_category` is in `group_by[]`; null for usage
	// that is not Claude Tag. Direct-message usage is billed to the individual user
	// and is reported under that user's product, not under `claude-tag`. New
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
	// Cost component when `group_by[]=cost_type`; null otherwise (amount is the
	// combined total).
	//
	// Any of "code_execution", "tokens", "web_search".
	CostType BetaAnalyticsCostType `json:"cost_type" api:"required"`
	// Currency code for the cost amount. Currently always `"USD"`.
	Currency string `json:"currency" api:"required"`
	// Inference region of the usage or cost. Null unless `inference_geo` is in
	// `group_by[]`; it can also be null on grouped rows where the region is not set
	// (the rows that `inference_geos[]=not_available` matches).
	//
	// Any of "global", "us".
	InferenceGeo BetaAnalyticsCostBucketedResultInferenceGeo `json:"inference_geo" api:"required"`
	// List-price amount (pre-discount) in fractional cents.
	ListAmount string `json:"list_amount" api:"required"`
	// Model that produced the usage or cost, as a model name in the form the
	// `models[]` filter accepts (for example, `claude-opus-5`). Null unless `model` is
	// in `group_by[]`; it can also be null on grouped rows whose usage or cost is not
	// attributed to a specific model, such as code execution.
	Model string `json:"model" api:"required"`
	// Product surface that produced the usage or cost. Null unless product is in
	// `group_by[]`; it can also be null on grouped rows whose usage cannot be
	// attributed to a known surface. Values include `chat`, `claude_code`, `cowork`,
	// `office_agent`, `claude_in_chrome`, `claude_design`, `claude-tag`, and
	// `chat_cowork_unified`. `claude-tag` is Claude Tag, the Claude product in Slack.
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
	// Number of API requests in this row's scope. Null when `group_by` includes
	// `cost_type` or `token_type` (the count has no per-component attribution; read it
	// from the ungrouped response). For sandbox / code-execution events, this counts
	// execution spans rather than HTTP requests (these rows surface with
	// `product: null`).
	Requests int64 `json:"requests" api:"required"`
	// Slack channel the usage originated from. Populated only when `slack_channel_id`
	// is in `group_by[]`; null for usage outside Slack (and for rows recorded before
	// channel attribution was enabled).
	SlackChannelID string `json:"slack_channel_id" api:"required"`
	// Inference speed mode of the usage or cost: `fast` or `standard`. Null unless
	// `speed` is in `group_by[]`.
	//
	// Any of "fast", "standard".
	Speed BetaAnalyticsCostBucketedResultSpeed `json:"speed" api:"required"`
	// Token type when `group_by[]=token_type` and `cost_type=tokens`; null otherwise.
	//
	// Any of "cache_creation.ephemeral_1h_input_tokens",
	// "cache_creation.ephemeral_5m_input_tokens", "cache_read_input_tokens",
	// "output_tokens", "uncached_input_tokens".
	TokenType BetaAnalyticsTokenType `json:"token_type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Amount            respjson.Field
		ClaudeTagCategory respjson.Field
		ClaudeTagUserID   respjson.Field
		ContextWindow     respjson.Field
		CostType          respjson.Field
		Currency          respjson.Field
		InferenceGeo      respjson.Field
		ListAmount        respjson.Field
		Model             respjson.Field
		Product           respjson.Field
		RBACGroupID       respjson.Field
		Requests          respjson.Field
		SlackChannelID    respjson.Field
		Speed             respjson.Field
		TokenType         respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsCostBucketedResult) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsCostBucketedResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inference region of the usage or cost. Null unless `inference_geo` is in
// `group_by[]`; it can also be null on grouped rows where the region is not set
// (the rows that `inference_geos[]=not_available` matches).
type BetaAnalyticsCostBucketedResultInferenceGeo string

const (
	BetaAnalyticsCostBucketedResultInferenceGeoGlobal BetaAnalyticsCostBucketedResultInferenceGeo = "global"
	BetaAnalyticsCostBucketedResultInferenceGeoUs     BetaAnalyticsCostBucketedResultInferenceGeo = "us"
)

// Inference speed mode of the usage or cost: `fast` or `standard`. Null unless
// `speed` is in `group_by[]`.
type BetaAnalyticsCostBucketedResultSpeed string

const (
	BetaAnalyticsCostBucketedResultSpeedFast     BetaAnalyticsCostBucketedResultSpeed = "fast"
	BetaAnalyticsCostBucketedResultSpeedStandard BetaAnalyticsCostBucketedResultSpeed = "standard"
)

type BetaAnalyticsCostReportTimeBucket struct {
	// End of the time bucket (exclusive) in RFC 3339 format.
	EndingAt time.Time `json:"ending_at" api:"required" format:"date-time"`
	// Rows for this time bucket. Empty when the bucket has no data; otherwise a single
	// combined row when `group_by[]` is omitted, or one row per group (subject to the
	// per-bucket group cap described on the `group_by[]` parameter).
	Results []BetaAnalyticsCostBucketedResult `json:"results" api:"required"`
	// Start of the time bucket (inclusive) in RFC 3339 format.
	StartingAt time.Time `json:"starting_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EndingAt    respjson.Field
		Results     respjson.Field
		StartingAt  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsCostReportTimeBucket) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsCostReportTimeBucket) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAnalyticsCostType string

const (
	BetaAnalyticsCostTypeCodeExecution BetaAnalyticsCostType = "code_execution"
	BetaAnalyticsCostTypeTokens        BetaAnalyticsCostType = "tokens"
	BetaAnalyticsCostTypeWebSearch     BetaAnalyticsCostType = "web_search"
)

type BetaAnalyticsCostUsersItem struct {
	// The user this row's usage or cost is attributed to. Always a `user_actor`.
	Actor BetaAnalyticsCostUsersItemActorUnion `json:"actor" api:"required"`
	// Amount (post-discount, pre-credit) in fractional cents (minor units).
	Amount string `json:"amount" api:"required"`
	// Claude Tag (Claude in Slack) spend category: `engaged` (a person addressed
	// Claude in a channel or thread), `proactive` (Claude responded without being
	// addressed), `scheduled` (a scheduled routine ran), `monitoring` (Claude watching
	// a channel it was asked to monitor), or `dm` (direct messages with Claude).
	// Populated only when `claude_tag_category` is in `group_by[]`; null for usage
	// that is not Claude Tag. Direct-message usage is billed to the individual user
	// and is reported under that user's product, not under `claude-tag`. New
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
	// Cost component breakdown; null when returning the combined total.
	//
	// Any of "code_execution", "tokens", "web_search".
	CostType BetaAnalyticsCostType `json:"cost_type" api:"required"`
	// Currency code for the cost amount. Currently always `"USD"`.
	Currency string `json:"currency" api:"required"`
	// End of the row's UTC time bucket (exclusive), as an RFC 3339 timestamp; equal to
	// `starting_at` plus one `bucket_width`. Null unless `bucket_width` is set.
	EndingAt time.Time `json:"ending_at" api:"required" format:"date-time"`
	// Inference region of the usage or cost. Null unless `inference_geo` is in
	// `group_by[]`; it can also be null on grouped rows where the region is not set
	// (the rows that `inference_geos[]=not_available` matches).
	//
	// Any of "global", "us".
	InferenceGeo BetaAnalyticsCostUsersItemInferenceGeo `json:"inference_geo" api:"required"`
	// List-price amount (pre-discount) in fractional cents.
	ListAmount string `json:"list_amount" api:"required"`
	// Model that produced the usage or cost, as a model name in the form the
	// `models[]` filter accepts (for example, `claude-opus-5`). Null unless `model` is
	// in `group_by[]`; it can also be null on grouped rows whose usage or cost is not
	// attributed to a specific model, such as code execution.
	Model string `json:"model" api:"required"`
	// Product surface that produced the usage or cost. Null unless product is in
	// `group_by[]`; it can also be null on grouped rows whose usage cannot be
	// attributed to a known surface. Values include `chat`, `claude_code`, `cowork`,
	// `office_agent`, `claude_in_chrome`, `claude_design`, `claude-tag`, and
	// `chat_cowork_unified`. `claude-tag` is Claude Tag, the Claude product in Slack.
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
	// Number of API requests in this row's scope. Null when `group_by` includes
	// `cost_type` or `token_type` (the count has no per-component attribution; read it
	// from the ungrouped response). For sandbox / code-execution events, this counts
	// execution spans rather than HTTP requests (these rows surface with
	// `product: null`).
	Requests int64 `json:"requests" api:"required"`
	// Slack channel the usage originated from. Populated only when `slack_channel_id`
	// is in `group_by[]`; null for usage outside Slack (and for rows recorded before
	// channel attribution was enabled).
	SlackChannelID string `json:"slack_channel_id" api:"required"`
	// Inference speed mode of the usage or cost: `fast` or `standard`. Null unless
	// `speed` is in `group_by[]`.
	//
	// Any of "fast", "standard".
	Speed BetaAnalyticsCostUsersItemSpeed `json:"speed" api:"required"`
	// Start of the row's UTC time bucket (inclusive), as an RFC 3339 timestamp. Null
	// unless `bucket_width` is set; without `bucket_width`, each row aggregates the
	// full requested range.
	StartingAt time.Time `json:"starting_at" api:"required" format:"date-time"`
	// Token type when `cost_type` is `tokens`; null otherwise.
	//
	// Any of "cache_creation.ephemeral_1h_input_tokens",
	// "cache_creation.ephemeral_5m_input_tokens", "cache_read_input_tokens",
	// "output_tokens", "uncached_input_tokens".
	TokenType BetaAnalyticsTokenType `json:"token_type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Actor             respjson.Field
		Amount            respjson.Field
		ClaudeTagCategory respjson.Field
		ClaudeTagUserID   respjson.Field
		ContextWindow     respjson.Field
		CostType          respjson.Field
		Currency          respjson.Field
		EndingAt          respjson.Field
		InferenceGeo      respjson.Field
		ListAmount        respjson.Field
		Model             respjson.Field
		Product           respjson.Field
		RBACGroupID       respjson.Field
		Requests          respjson.Field
		SlackChannelID    respjson.Field
		Speed             respjson.Field
		StartingAt        respjson.Field
		TokenType         respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsCostUsersItem) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsCostUsersItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaAnalyticsCostUsersItemActorUnion contains all possible properties and values
// from [BetaAnalyticsUserActor].
//
// Use the [BetaAnalyticsCostUsersItemActorUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaAnalyticsCostUsersItemActorUnion struct {
	// This field is from variant [BetaAnalyticsUserActor].
	Deleted bool `json:"deleted"`
	// This field is from variant [BetaAnalyticsUserActor].
	EmailAddress string `json:"email_address"`
	// This field is from variant [BetaAnalyticsUserActor].
	Name string `json:"name"`
	// Any of "user_actor".
	Type string `json:"type"`
	// This field is from variant [BetaAnalyticsUserActor].
	UserID string `json:"user_id"`
	JSON   struct {
		Deleted      respjson.Field
		EmailAddress respjson.Field
		Name         respjson.Field
		Type         respjson.Field
		UserID       respjson.Field
		raw          string
	} `json:"-"`
}

// anyBetaAnalyticsCostUsersItemActor is implemented by each variant of
// [BetaAnalyticsCostUsersItemActorUnion] to add type safety for the return type of
// [BetaAnalyticsCostUsersItemActorUnion.AsAny]
type anyBetaAnalyticsCostUsersItemActor interface {
	implBetaAnalyticsCostUsersItemActorUnion()
}

func (BetaAnalyticsUserActor) implBetaAnalyticsCostUsersItemActorUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaAnalyticsCostUsersItemActorUnion.AsAny().(type) {
//	case anthropic.BetaAnalyticsUserActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaAnalyticsCostUsersItemActorUnion) AsAny() anyBetaAnalyticsCostUsersItemActor {
	switch u.Type {
	case "user_actor":
		return u.AsUserActor()
	}
	return nil
}

func (u BetaAnalyticsCostUsersItemActorUnion) AsUserActor() (v BetaAnalyticsUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaAnalyticsCostUsersItemActorUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaAnalyticsCostUsersItemActorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inference region of the usage or cost. Null unless `inference_geo` is in
// `group_by[]`; it can also be null on grouped rows where the region is not set
// (the rows that `inference_geos[]=not_available` matches).
type BetaAnalyticsCostUsersItemInferenceGeo string

const (
	BetaAnalyticsCostUsersItemInferenceGeoGlobal BetaAnalyticsCostUsersItemInferenceGeo = "global"
	BetaAnalyticsCostUsersItemInferenceGeoUs     BetaAnalyticsCostUsersItemInferenceGeo = "us"
)

// Inference speed mode of the usage or cost: `fast` or `standard`. Null unless
// `speed` is in `group_by[]`.
type BetaAnalyticsCostUsersItemSpeed string

const (
	BetaAnalyticsCostUsersItemSpeedFast     BetaAnalyticsCostUsersItemSpeed = "fast"
	BetaAnalyticsCostUsersItemSpeedStandard BetaAnalyticsCostUsersItemSpeed = "standard"
)

// Cowork activity metrics for a single user on a given day.
type BetaAnalyticsCoworkMetrics struct {
	// Number of tool actions completed in Cowork sessions
	ActionCount int64 `json:"action_count" api:"required"`
	// Number of artifacts created in Cowork sessions: an artifact counts once, on the
	// day a session first saves it. Counted from 2026-08-17; 0 on earlier days. Exact
	// in date-range mode: a creation belongs to exactly one day, so the per-day counts
	// never overlap and their sum over the window is the exact count of distinct
	// creations in it.
	ArtifactsCreatedCount int64 `json:"artifacts_created_count" api:"required"`
	// Total number of connector invocations in Cowork sessions
	ConnectorsUsedCount int64 `json:"connectors_used_count" api:"required"`
	// Number of Dispatch (background agent) turns completed
	DispatchTurnCount int64 `json:"dispatch_turn_count" api:"required"`
	// Number of distinct connectors used in Cowork sessions. Approximate (HLL, typical
	// error <2%) in date-range mode. Null on aggregated rows where a distinct count
	// cannot be computed.
	DistinctConnectorsUsedCount int64 `json:"distinct_connectors_used_count" api:"required"`
	// Number of distinct Cowork sessions. Approximate (HLL, typical error <2%) in
	// date-range mode. Null on aggregated rows where a distinct count cannot be
	// computed.
	DistinctSessionCount int64 `json:"distinct_session_count" api:"required"`
	// Number of distinct skills used in Cowork sessions. Approximate (HLL, typical
	// error <2%) in date-range mode. Null on aggregated rows where a distinct count
	// cannot be computed.
	DistinctSkillsUsedCount int64 `json:"distinct_skills_used_count" api:"required"`
	// Number of messages sent in Cowork sessions
	MessageCount int64 `json:"message_count" api:"required"`
	// Total number of skill invocations in Cowork sessions
	SkillsUsedCount int64 `json:"skills_used_count" api:"required"`
	// Number of distinct plugins used in Cowork sessions. Null while Cowork plugin-use
	// metrics are not enabled for this organization. Approximate (HLL, typical error
	// <2%) in date-range mode. Null on aggregated rows where a distinct count cannot
	// be computed.
	DistinctPluginsUsedCount int64 `json:"distinct_plugins_used_count" api:"nullable"`
	// Number of successful Edit tool calls in Cowork sessions. Null while the
	// file-edit metrics are not enabled for this organization.
	EditToolCount int64 `json:"edit_tool_count" api:"nullable"`
	// Number of successful file-edit tool calls (Edit, MultiEdit, Write, NotebookEdit)
	// in Cowork sessions. Null, never 0, while the file-edit metrics are not enabled
	// for this organization.
	FileEditCount int64 `json:"file_edit_count" api:"nullable"`
	// Number of successful MultiEdit tool calls in Cowork sessions. Null while the
	// file-edit metrics are not enabled for this organization.
	MultiEditToolCount int64 `json:"multi_edit_tool_count" api:"nullable"`
	// Number of successful NotebookEdit tool calls in Cowork sessions. Null while the
	// file-edit metrics are not enabled for this organization.
	NotebookEditToolCount int64 `json:"notebook_edit_tool_count" api:"nullable"`
	// Total number of plugin invocations in Cowork sessions. Null while Cowork
	// plugin-use metrics are not enabled for this organization.
	PluginsUsedCount int64 `json:"plugins_used_count" api:"nullable"`
	// Number of distinct Cowork sessions with at least one successful file-edit tool
	// call. Null while the file-edit metrics are not enabled for this organization.
	// Approximate (HLL, typical error <2%) in date-range mode. Null on aggregated rows
	// where a distinct count cannot be computed.
	SessionsWithFileEditsCount int64 `json:"sessions_with_file_edits_count" api:"nullable"`
	// Number of successful Write tool calls in Cowork sessions. Null while the
	// file-edit metrics are not enabled for this organization.
	WriteToolCount int64 `json:"write_tool_count" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ActionCount                 respjson.Field
		ArtifactsCreatedCount       respjson.Field
		ConnectorsUsedCount         respjson.Field
		DispatchTurnCount           respjson.Field
		DistinctConnectorsUsedCount respjson.Field
		DistinctSessionCount        respjson.Field
		DistinctSkillsUsedCount     respjson.Field
		MessageCount                respjson.Field
		SkillsUsedCount             respjson.Field
		DistinctPluginsUsedCount    respjson.Field
		EditToolCount               respjson.Field
		FileEditCount               respjson.Field
		MultiEditToolCount          respjson.Field
		NotebookEditToolCount       respjson.Field
		PluginsUsedCount            respjson.Field
		SessionsWithFileEditsCount  respjson.Field
		WriteToolCount              respjson.Field
		ExtraFields                 map[string]respjson.Field
		raw                         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsCoworkMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsCoworkMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude Design activity metrics for a single user on a given day.
type BetaAnalyticsDesignMetrics struct {
	// Number of distinct Claude Design projects created. Exact in date-range mode: a
	// creation belongs to exactly one day, so the per-day counts never overlap and
	// their sum over the window is the exact count of distinct creations in it.
	DistinctProjectsCreatedCount int64 `json:"distinct_projects_created_count" api:"required"`
	// Number of distinct Claude Design projects the user worked in. Approximate (HLL,
	// typical error <2%) in date-range mode. Null on aggregated rows where a distinct
	// count cannot be computed.
	DistinctProjectsUsedCount int64 `json:"distinct_projects_used_count" api:"required"`
	// Number of distinct Claude Design sessions. Approximate (HLL, typical error <2%)
	// in date-range mode. Null on aggregated rows where a distinct count cannot be
	// computed.
	DistinctSessionCount int64 `json:"distinct_session_count" api:"required"`
	// Number of messages sent in Claude Design sessions
	MessageCount int64 `json:"message_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctProjectsCreatedCount respjson.Field
		DistinctProjectsUsedCount    respjson.Field
		DistinctSessionCount         respjson.Field
		MessageCount                 respjson.Field
		ExtraFields                  map[string]respjson.Field
		raw                          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsDesignMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsDesignMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAnalyticsInferenceGeoFilter string

const (
	BetaAnalyticsInferenceGeoFilterGlobal       BetaAnalyticsInferenceGeoFilter = "global"
	BetaAnalyticsInferenceGeoFilterNotAvailable BetaAnalyticsInferenceGeoFilter = "not_available"
	BetaAnalyticsInferenceGeoFilterUs           BetaAnalyticsInferenceGeoFilter = "us"
)

// Lines of code added and removed via Claude Code.
type BetaAnalyticsLinesOfCode struct {
	// Lines of code added
	AddedCount int64 `json:"added_count" api:"required"`
	// Lines of code removed
	RemovedCount int64 `json:"removed_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AddedCount   respjson.Field
		RemovedCount respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsLinesOfCode) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsLinesOfCode) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Office Agent activity metrics for a single user on a given day, broken out by
// Office product.
type BetaAnalyticsOfficeMetrics struct {
	// Office Agent activity metrics for a single user on a given day within one Office
	// product.
	Excel BetaAnalyticsOfficeProductMetrics `json:"excel" api:"required"`
	// Office Agent activity metrics for a single user on a given day within one Office
	// product.
	Outlook BetaAnalyticsOfficeProductMetrics `json:"outlook" api:"required"`
	// Office Agent activity metrics for a single user on a given day within one Office
	// product.
	Powerpoint BetaAnalyticsOfficeProductMetrics `json:"powerpoint" api:"required"`
	// Office Agent activity metrics for a single user on a given day within one Office
	// product.
	Word BetaAnalyticsOfficeProductMetrics `json:"word" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Excel       respjson.Field
		Outlook     respjson.Field
		Powerpoint  respjson.Field
		Word        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsOfficeMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsOfficeMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Office Agent activity metrics for a single user on a given day within one Office
// product.
type BetaAnalyticsOfficeProductMetrics struct {
	// Number of MCP connector invocations
	ConnectorsUsedCount int64 `json:"connectors_used_count" api:"required"`
	// Number of distinct MCP connectors used. Approximate (HLL, typical error <2%) in
	// date-range mode. Null on aggregated rows where a distinct count cannot be
	// computed.
	DistinctConnectorsUsedCount int64 `json:"distinct_connectors_used_count" api:"required"`
	// Number of distinct Office Agent sessions. Approximate (HLL, typical error <2%)
	// in date-range mode. Null on aggregated rows where a distinct count cannot be
	// computed.
	DistinctSessionCount int64 `json:"distinct_session_count" api:"required"`
	// Number of distinct skills used. Approximate (HLL, typical error <2%) in
	// date-range mode. Null on aggregated rows where a distinct count cannot be
	// computed.
	DistinctSkillsUsedCount int64 `json:"distinct_skills_used_count" api:"required"`
	// Number of messages sent
	MessageCount int64 `json:"message_count" api:"required"`
	// Number of skill invocations
	SkillsUsedCount int64 `json:"skills_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ConnectorsUsedCount         respjson.Field
		DistinctConnectorsUsedCount respjson.Field
		DistinctSessionCount        respjson.Field
		DistinctSkillsUsedCount     respjson.Field
		MessageCount                respjson.Field
		SkillsUsedCount             respjson.Field
		ExtraFields                 map[string]respjson.Field
		raw                         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsOfficeProductMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsOfficeProductMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Per-plugin install + invocation activity for a given day.
//
// With `group_by[]=user_id` / `rbac_group_id` / `product` (`cowork`, `claude_code`
// and `chat_cowork_unified` only on this endpoint) each row is one (plugin, user),
// (plugin, group), or (plugin, product) cut: the flat `user_id` / `rbac_group_id`
// / `product` keys carry the cut and the counts are scoped to it.
type BetaAnalyticsPluginActivity struct {
	// Claude Code activity metrics for a single plugin on a given day.
	ClaudeCodeMetrics BetaAnalyticsPluginClaudeCodeMetrics `json:"claude_code_metrics" api:"required"`
	// Cowork activity metrics for a single plugin on a given day.
	CoworkMetrics BetaAnalyticsPluginCoworkMetrics `json:"cowork_metrics" api:"required"`
	// Number of distinct users with recorded install or invocation activity for the
	// plugin on the requested day (install-only users count), or, in date-range mode,
	// over the requested window — recomputed as an exact distinct count over the
	// window's per-member daily rows, never a sum of per-day values.
	DistinctUserCount int64 `json:"distinct_user_count" api:"required"`
	// Number of distinct users who installed the plugin on the requested day, or, in
	// date-range mode, over the requested window — recomputed as an exact distinct
	// count over the window's per-member daily rows, never a sum of per-day values.
	InstallCount int64 `json:"install_count" api:"required"`
	// Number of plugin invocations on the requested day
	InvocationCount int64 `json:"invocation_count" api:"required"`
	// Name of the plugin
	PluginName string `json:"plugin_name" api:"required"`
	// Plugin use recorded while members had Chat and Cowork unified (Cowork's features
	// inside claude.ai chat) turned on. A count is null in date-range mode where it
	// cannot be computed. Omitted from the response on deployments that do not offer
	// Chat and Cowork unified.
	ChatCoworkUnifiedMetrics BetaAnalyticsPluginActivityChatCoworkUnifiedMetrics `json:"chat_cowork_unified_metrics" api:"nullable"`
	// Stable plugin identifier when available (e.g. `serena@claude-plugins-official`).
	// Null for third-party Claude Code plugins (redacted at the source) and Cowork
	// slash commands that carry only a hashed id.
	PluginID string `json:"plugin_id" api:"nullable"`
	// Product that produced this row's activity: one of `chat`, `claude_code`,
	// `cowork`, `office_agent`, or `chat_cowork_unified` (Chat and Cowork unified).
	// These are the canonical Cost & Usage product names; an `office_agent` row's
	// per-surface breakdown is in its `office_metrics`. On `/plugins` only `cowork`,
	// `claude_code` and `chat_cowork_unified` occur (the only surfaces with plugin
	// attribution); on `/artifacts` only `chat`, `claude_code`, `cowork` and
	// `chat_cowork_unified` occur (the surfaces that create artifacts);
	// `/apps/chat/projects` does not support the product dimension (a `product` entry
	// in `group_by[]` or `filter[]` there is rejected). Present only when the request
	// grouped by `product`.
	Product string `json:"product" api:"nullable"`
	// Tagged RBAC group identifier (`rbac_group_...`), matching the spend-limits API
	// spelling. Present only when the request grouped by `rbac_group_id`.
	RBACGroupID string `json:"rbac_group_id" api:"nullable"`
	// Resolved RBAC group display name, alongside `rbac_group_id` when name resolution
	// is available. Null if the group has been deleted or its name could not be
	// resolved; `rbac_group_id` remains the stable key.
	RBACGroupName string `json:"rbac_group_name" api:"nullable"`
	// Tagged user identifier (e.g. `user_...`). Present only when the request grouped
	// by `user_id`.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ClaudeCodeMetrics        respjson.Field
		CoworkMetrics            respjson.Field
		DistinctUserCount        respjson.Field
		InstallCount             respjson.Field
		InvocationCount          respjson.Field
		PluginName               respjson.Field
		ChatCoworkUnifiedMetrics respjson.Field
		PluginID                 respjson.Field
		Product                  respjson.Field
		RBACGroupID              respjson.Field
		RBACGroupName            respjson.Field
		UserID                   respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsPluginActivity) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsPluginActivity) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Plugin use recorded while members had Chat and Cowork unified (Cowork's features
// inside claude.ai chat) turned on. A count is null in date-range mode where it
// cannot be computed. Omitted from the response on deployments that do not offer
// Chat and Cowork unified.
type BetaAnalyticsPluginActivityChatCoworkUnifiedMetrics struct {
	// Same measure as `cowork_metrics.distinct_session_plugin_used_count`, for
	// activity recorded while members had Chat and Cowork unified turned on. Null on
	// aggregated rows where a distinct count cannot be computed.
	DistinctSessionPluginUsedCount int64 `json:"distinct_session_plugin_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctSessionPluginUsedCount respjson.Field
		ExtraFields                    map[string]respjson.Field
		raw                            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsPluginActivityChatCoworkUnifiedMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsPluginActivityChatCoworkUnifiedMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude Code activity metrics for a single plugin on a given day.
type BetaAnalyticsPluginClaudeCodeMetrics struct {
	// Number of distinct Claude Code sessions in which the plugin was invoked. Null on
	// aggregated rows where a distinct count cannot be computed.
	DistinctSessionPluginUsedCount int64 `json:"distinct_session_plugin_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctSessionPluginUsedCount respjson.Field
		ExtraFields                    map[string]respjson.Field
		raw                            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsPluginClaudeCodeMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsPluginClaudeCodeMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Cowork activity metrics for a single plugin on a given day.
type BetaAnalyticsPluginCoworkMetrics struct {
	// Number of distinct Cowork sessions in which the plugin was invoked. Null on
	// aggregated rows where a distinct count cannot be computed.
	DistinctSessionPluginUsedCount int64 `json:"distinct_session_plugin_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctSessionPluginUsedCount respjson.Field
		ExtraFields                    map[string]respjson.Field
		raw                            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsPluginCoworkMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsPluginCoworkMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Publicly documented product surfaces. `claude-tag` is Claude Tag, the Claude
// product in Slack. `chat_cowork_unified` is Chat and Cowork unified, Cowork's
// features inside claude.ai chat: chat and Cowork usage by a member who has it
// turned on is reported under this value instead of `chat` or `cowork`. It is
// accepted as a filter only on deployments that offer Chat and Cowork unified.
type BetaAnalyticsProductFilter string

const (
	BetaAnalyticsProductFilterChat              BetaAnalyticsProductFilter = "chat"
	BetaAnalyticsProductFilterChatCoworkUnified BetaAnalyticsProductFilter = "chat_cowork_unified"
	BetaAnalyticsProductFilterClaudeTag         BetaAnalyticsProductFilter = "claude-tag"
	BetaAnalyticsProductFilterClaudeCode        BetaAnalyticsProductFilter = "claude_code"
	BetaAnalyticsProductFilterClaudeDesign      BetaAnalyticsProductFilter = "claude_design"
	BetaAnalyticsProductFilterClaudeInChrome    BetaAnalyticsProductFilter = "claude_in_chrome"
	BetaAnalyticsProductFilterCowork            BetaAnalyticsProductFilter = "cowork"
	BetaAnalyticsProductFilterOfficeAgent       BetaAnalyticsProductFilter = "office_agent"
)

// Per-project activity data for a given day.
type BetaAnalyticsProjectActivity struct {
	// Number of distinct users who used the project on the requested day, or, in
	// date-range mode, over the requested window — recomputed as an exact distinct
	// count over the window's per-member daily rows, never a sum of per-day values.
	DistinctUserCount int64 `json:"distinct_user_count" api:"required"`
	// Number of messages sent in the project on the requested day
	MessageCount int64 `json:"message_count" api:"required"`
	// Tagged project identifier (e.g. `claude_proj_...`)
	ProjectID string `json:"project_id" api:"required"`
	// Name of the project
	ProjectName string `json:"project_name" api:"required"`
	// Project creation timestamp in RFC 3339 format. Null if the project was deleted
	// before attribution was recorded.
	CreatedAt time.Time `json:"created_at" api:"nullable" format:"date-time"`
	// User who created the project. Null if the project was deleted before attribution
	// was recorded, or if the creator's account no longer exists.
	CreatedBy BetaAnalyticsUser `json:"created_by" api:"nullable"`
	// Number of distinct conversations in the project. Null on aggregated rows where a
	// distinct count cannot be computed.
	DistinctConversationCount int64 `json:"distinct_conversation_count" api:"nullable"`
	// Product that produced this row's activity: one of `chat`, `claude_code`,
	// `cowork`, `office_agent`, or `chat_cowork_unified` (Chat and Cowork unified).
	// These are the canonical Cost & Usage product names; an `office_agent` row's
	// per-surface breakdown is in its `office_metrics`. On `/plugins` only `cowork`,
	// `claude_code` and `chat_cowork_unified` occur (the only surfaces with plugin
	// attribution); on `/artifacts` only `chat`, `claude_code`, `cowork` and
	// `chat_cowork_unified` occur (the surfaces that create artifacts);
	// `/apps/chat/projects` does not support the product dimension (a `product` entry
	// in `group_by[]` or `filter[]` there is rejected). Present only when the request
	// grouped by `product`.
	Product string `json:"product" api:"nullable"`
	// Tagged RBAC group identifier (`rbac_group_...`), matching the spend-limits API
	// spelling. Present only when the request grouped by `rbac_group_id`.
	RBACGroupID string `json:"rbac_group_id" api:"nullable"`
	// Resolved RBAC group display name, alongside `rbac_group_id` when name resolution
	// is available. Null if the group has been deleted or its name could not be
	// resolved; `rbac_group_id` remains the stable key.
	RBACGroupName string `json:"rbac_group_name" api:"nullable"`
	// Tagged user identifier (e.g. `user_...`). Present only when the request grouped
	// by `user_id`.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctUserCount         respjson.Field
		MessageCount              respjson.Field
		ProjectID                 respjson.Field
		ProjectName               respjson.Field
		CreatedAt                 respjson.Field
		CreatedBy                 respjson.Field
		DistinctConversationCount respjson.Field
		Product                   respjson.Field
		RBACGroupID               respjson.Field
		RBACGroupName             respjson.Field
		UserID                    respjson.Field
		ExtraFields               map[string]respjson.Field
		raw                       string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsProjectActivity) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsProjectActivity) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude Science activity metrics for a single user on a given day.
type BetaAnalyticsScienceMetrics struct {
	// Number of delegations (handoffs to a specialized agent) in Claude Science
	// sessions
	DelegationCount int64 `json:"delegation_count" api:"required"`
	// Number of distinct Claude Science sessions. Approximate (HLL, typical error <2%)
	// in date-range mode. Null on aggregated rows where a distinct count cannot be
	// computed.
	DistinctSessionCount int64 `json:"distinct_session_count" api:"required"`
	// Number of messages sent in Claude Science sessions
	MessageCount int64 `json:"message_count" api:"required"`
	// Number of remote compute jobs launched from Claude Science sessions
	RemoteComputeJobCount int64 `json:"remote_compute_job_count" api:"required"`
	// Total number of skill invocations in Claude Science sessions
	SkillsUsedCount int64 `json:"skills_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DelegationCount       respjson.Field
		DistinctSessionCount  respjson.Field
		MessageCount          respjson.Field
		RemoteComputeJobCount respjson.Field
		SkillsUsedCount       respjson.Field
		ExtraFields           map[string]respjson.Field
		raw                   string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsScienceMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsScienceMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAnalyticsServerToolUse struct {
	// The number of web search requests made.
	WebSearchRequests int64 `json:"web_search_requests" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		WebSearchRequests respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsServerToolUse) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsServerToolUse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Per-day entry in the /summaries response.
type BetaAnalyticsSingleDayActivitySummary struct {
	// Number of seats currently assigned to members. Null when the response is scoped
	// to an RBAC group — seat assignment is org-wide and has no per-group analogue.
	AssignedSeatCount int64 `json:"assigned_seat_count" api:"required"`
	// Number of users with Cowork activity on the requested day
	CoworkDailyActiveUserCount int64 `json:"cowork_daily_active_user_count" api:"required"`
	// Number of users with Cowork activity in the 30-day rolling window
	CoworkMonthlyActiveUserCount int64 `json:"cowork_monthly_active_user_count" api:"required"`
	// Number of users with Cowork activity in the 7-day rolling window
	CoworkWeeklyActiveUserCount int64 `json:"cowork_weekly_active_user_count" api:"required"`
	// Number of users with token consumption on the requested day
	DailyActiveUserCount int64 `json:"daily_active_user_count" api:"required"`
	// Percentage of assigned seats with activity on the requested day
	// (`DAU / assigned_seat_count * 100`). Null when the response is scoped to an RBAC
	// group.
	DailyAdoptionRate float64 `json:"daily_adoption_rate" api:"required"`
	// End of the aggregation period (exclusive), UTC midnight in RFC 3339 format (e.g.
	// `2026-01-16T00:00:00Z`).
	EndingAt time.Time `json:"ending_at" api:"required" format:"date-time"`
	// Number of users with token consumption in the 30-day rolling window
	MonthlyActiveUserCount int64 `json:"monthly_active_user_count" api:"required"`
	// Percentage of assigned seats with activity in the 30-day rolling window
	// (`MAU / assigned_seat_count * 100`). Null when the response is scoped to an RBAC
	// group.
	MonthlyAdoptionRate float64 `json:"monthly_adoption_rate" api:"required"`
	// Number of pending invitations to join the organization. Null when the response
	// is scoped to an RBAC group.
	PendingInviteCount int64 `json:"pending_invite_count" api:"required"`
	// Start of the aggregation period (inclusive), UTC midnight in RFC 3339 format
	// (e.g. `2026-01-15T00:00:00Z`).
	StartingAt time.Time `json:"starting_at" api:"required" format:"date-time"`
	// Number of users with token consumption in the 7-day rolling window
	WeeklyActiveUserCount int64 `json:"weekly_active_user_count" api:"required"`
	// Percentage of assigned seats with activity in the 7-day rolling window
	// (`WAU / assigned_seat_count * 100`). Null when the response is scoped to an RBAC
	// group.
	WeeklyAdoptionRate float64 `json:"weekly_adoption_rate" api:"required"`
	// Number of users with activity in Chat and Cowork unified on the requested day.
	// Omitted from the response on deployments that do not offer Chat and Cowork
	// unified.
	ChatCoworkUnifiedDailyActiveUserCount int64 `json:"chat_cowork_unified_daily_active_user_count" api:"nullable"`
	// Number of users with activity in Chat and Cowork unified in the 28-day rolling
	// window (30 days when the request filters by `rbac_group_id`). Omitted from the
	// response on deployments that do not offer Chat and Cowork unified.
	ChatCoworkUnifiedMonthlyActiveUserCount int64 `json:"chat_cowork_unified_monthly_active_user_count" api:"nullable"`
	// Number of users with activity in Chat and Cowork unified in the 7-day rolling
	// window. Omitted from the response on deployments that do not offer Chat and
	// Cowork unified.
	ChatCoworkUnifiedWeeklyActiveUserCount int64 `json:"chat_cowork_unified_weekly_active_user_count" api:"nullable"`
	// Number of users with claude.ai (chat) activity on the requested day. Omitted
	// from the response while the per-product breakdown is not enabled for this
	// organization.
	ChatDailyActiveUserCount int64 `json:"chat_daily_active_user_count" api:"nullable"`
	// Number of users with claude.ai (chat) activity in the 30-day rolling window.
	// Omitted from the response while the per-product breakdown is not enabled for
	// this organization.
	ChatMonthlyActiveUserCount int64 `json:"chat_monthly_active_user_count" api:"nullable"`
	// Number of users with claude.ai (chat) activity in the 7-day rolling window.
	// Omitted from the response while the per-product breakdown is not enabled for
	// this organization.
	ChatWeeklyActiveUserCount int64 `json:"chat_weekly_active_user_count" api:"nullable"`
	// Number of users with Claude Code activity on the requested day. Omitted from the
	// response while the per-product breakdown is not enabled for this organization.
	ClaudeCodeDailyActiveUserCount int64 `json:"claude_code_daily_active_user_count" api:"nullable"`
	// Number of users with Claude Code activity in the 30-day rolling window. Omitted
	// from the response while the per-product breakdown is not enabled for this
	// organization.
	ClaudeCodeMonthlyActiveUserCount int64 `json:"claude_code_monthly_active_user_count" api:"nullable"`
	// Number of users with Claude Code activity in the 7-day rolling window. Omitted
	// from the response while the per-product breakdown is not enabled for this
	// organization.
	ClaudeCodeWeeklyActiveUserCount int64 `json:"claude_code_weekly_active_user_count" api:"nullable"`
	// Number of users with Claude Design activity on the requested day. Omitted from
	// the response while the per-product breakdown is not enabled for this
	// organization.
	ClaudeDesignDailyActiveUserCount int64 `json:"claude_design_daily_active_user_count" api:"nullable"`
	// Number of users with Claude Design activity in the 30-day rolling window.
	// Omitted from the response while the per-product breakdown is not enabled for
	// this organization.
	ClaudeDesignMonthlyActiveUserCount int64 `json:"claude_design_monthly_active_user_count" api:"nullable"`
	// Number of users with Claude Design activity in the 7-day rolling window. Omitted
	// from the response while the per-product breakdown is not enabled for this
	// organization.
	ClaudeDesignWeeklyActiveUserCount int64 `json:"claude_design_weekly_active_user_count" api:"nullable"`
	// Number of users with Claude in Office activity on the requested day. Omitted
	// from the response while the per-product breakdown is not enabled for this
	// organization.
	OfficeAgentDailyActiveUserCount int64 `json:"office_agent_daily_active_user_count" api:"nullable"`
	// Number of users with Claude in Office activity in the 30-day rolling window.
	// Omitted from the response while the per-product breakdown is not enabled for
	// this organization.
	OfficeAgentMonthlyActiveUserCount int64 `json:"office_agent_monthly_active_user_count" api:"nullable"`
	// Number of users with Claude in Office activity in the 7-day rolling window.
	// Omitted from the response while the per-product breakdown is not enabled for
	// this organization.
	OfficeAgentWeeklyActiveUserCount int64 `json:"office_agent_weekly_active_user_count" api:"nullable"`
	// Number of users with Claude Science activity on the requested day. Omitted from
	// the response while the per-product breakdown is not enabled for this
	// organization.
	ScienceDailyActiveUserCount int64 `json:"science_daily_active_user_count" api:"nullable"`
	// Number of users with a Claude Science seat entitlement (per-seat RBAC) at the
	// time of the daily snapshot. The funnel top; independent of the org-level Claude
	// Science toggle. Null when the response is scoped to an RBAC group — entitlement
	// is org-wide and has no per-group analogue. Omitted from the response while the
	// per-product breakdown is not enabled for this organization.
	ScienceEntitledUserCount int64 `json:"science_entitled_user_count" api:"nullable"`
	// Number of users with Claude Science activity in the 30-day rolling window.
	// Omitted from the response while the per-product breakdown is not enabled for
	// this organization.
	ScienceMonthlyActiveUserCount int64 `json:"science_monthly_active_user_count" api:"nullable"`
	// Number of users with Claude Science activity in the 7-day rolling window.
	// Omitted from the response while the per-product breakdown is not enabled for
	// this organization.
	ScienceWeeklyActiveUserCount int64 `json:"science_weekly_active_user_count" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AssignedSeatCount                       respjson.Field
		CoworkDailyActiveUserCount              respjson.Field
		CoworkMonthlyActiveUserCount            respjson.Field
		CoworkWeeklyActiveUserCount             respjson.Field
		DailyActiveUserCount                    respjson.Field
		DailyAdoptionRate                       respjson.Field
		EndingAt                                respjson.Field
		MonthlyActiveUserCount                  respjson.Field
		MonthlyAdoptionRate                     respjson.Field
		PendingInviteCount                      respjson.Field
		StartingAt                              respjson.Field
		WeeklyActiveUserCount                   respjson.Field
		WeeklyAdoptionRate                      respjson.Field
		ChatCoworkUnifiedDailyActiveUserCount   respjson.Field
		ChatCoworkUnifiedMonthlyActiveUserCount respjson.Field
		ChatCoworkUnifiedWeeklyActiveUserCount  respjson.Field
		ChatDailyActiveUserCount                respjson.Field
		ChatMonthlyActiveUserCount              respjson.Field
		ChatWeeklyActiveUserCount               respjson.Field
		ClaudeCodeDailyActiveUserCount          respjson.Field
		ClaudeCodeMonthlyActiveUserCount        respjson.Field
		ClaudeCodeWeeklyActiveUserCount         respjson.Field
		ClaudeDesignDailyActiveUserCount        respjson.Field
		ClaudeDesignMonthlyActiveUserCount      respjson.Field
		ClaudeDesignWeeklyActiveUserCount       respjson.Field
		OfficeAgentDailyActiveUserCount         respjson.Field
		OfficeAgentMonthlyActiveUserCount       respjson.Field
		OfficeAgentWeeklyActiveUserCount        respjson.Field
		ScienceDailyActiveUserCount             respjson.Field
		ScienceEntitledUserCount                respjson.Field
		ScienceMonthlyActiveUserCount           respjson.Field
		ScienceWeeklyActiveUserCount            respjson.Field
		ExtraFields                             map[string]respjson.Field
		raw                                     string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsSingleDayActivitySummary) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsSingleDayActivitySummary) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Per-skill activity data for a given day.
type BetaAnalyticsSkillActivity struct {
	// Claude.ai activity metrics for a single skill on a given day.
	ChatMetrics BetaAnalyticsSkillChatMetrics `json:"chat_metrics" api:"required"`
	// Claude Code activity metrics for a single skill on a given day.
	ClaudeCodeMetrics BetaAnalyticsSkillClaudeCodeMetrics `json:"claude_code_metrics" api:"required"`
	// Cowork activity metrics for a single skill on a given day.
	CoworkMetrics BetaAnalyticsSkillCoworkMetrics `json:"cowork_metrics" api:"required"`
	// Number of distinct users who used the skill on the requested day, or, in
	// date-range mode, over the requested window — recomputed as an exact distinct
	// count over the window's per-member daily rows, never a sum of per-day values. A
	// skill counts as used only when it is explicitly activated — the model (or the
	// user, via the skill's slash command) invokes it, reading its instructions into
	// context as part of that activation. Skills that are merely installed or listed
	// as available, or whose content reaches the context without an activation
	// (preloaded, hook-injected, or read as a plain file), are not counted.
	DistinctUserCount int64 `json:"distinct_user_count" api:"required"`
	// Office Agent activity metrics for a single skill on a given day, broken out by
	// Office product.
	OfficeMetrics BetaAnalyticsSkillOfficeMetrics `json:"office_metrics" api:"required"`
	// Name of the skill
	SkillName string `json:"skill_name" api:"required"`
	// List-price (rate-card) value of the member requests attributed to this skill, as
	// a decimal string in the minor unit of `currency` (cents for USD), from Claude
	// Code, Cowork, and Office Agent request-level attribution — the value of requests
	// that involved the skill, not the skill's incremental cost. Unlike
	// `estimated_overage_spend` this reflects usage value regardless of how it was
	// funded — seat-covered usage counts — but it is undiscounted and does not tie to
	// billed spend or the organization's spend reporting. claude.ai chat usage carries
	// no request-level attribution and contributes nothing: the field is null on
	// `chat` product rows and on `office_agent` product cuts dated before 2026-06-18
	// (the Office Agent attribution data-start), and on ungrouped rows it covers the
	// Claude Code + Cowork + Office Agent share only (null when no attributable usage
	// exists). Also null under the same conditions as `estimated_overage_spend` (spend
	// reporting not enabled for this organization, `office_agent` product cuts before
	// the 2026-06-18 data-start). "0" means attributable usage existed but none was
	// attributed to this skill. Addable across days: date-range rollup mode returns
	// the window's sum. On `group_by[]` and `filter[]` shapes both amounts can total
	// below the ungrouped value for the same skill over the same date or range: spend
	// attributed to a member–skill pair with no counted usage on that day is excluded
	// from those cuts.
	AttributedListPrice string `json:"attributed_list_price" api:"nullable"`
	// Skill use recorded while members had Chat and Cowork unified (Cowork's features
	// inside claude.ai chat) turned on, split into chat conversations and Cowork
	// sessions. A count is null in date-range mode where it cannot be computed.
	// Omitted from the response on deployments that do not offer Chat and Cowork
	// unified.
	ChatCoworkUnifiedMetrics BetaAnalyticsSkillActivityChatCoworkUnifiedMetrics `json:"chat_cowork_unified_metrics" api:"nullable"`
	// Currency for this row's monetary fields (`estimated_overage_spend` and
	// `attributed_list_price`), as an uppercase ISO-4217 code. Always "USD" when
	// either amount is populated; null whenever both amounts are null.
	Currency string `json:"currency" api:"nullable"`
	// Distinct accounts that enabled this skill on the requested day (claude.ai only —
	// the skill analog of plugin `install_count`). The count is org-wide: null when
	// enable reporting is not enabled for this organization, or when the request
	// scopes to `user_id` / `rbac_group_id` / `product` via `group_by[]` or `filter[]`
	// (an org-wide count would be misleading on per-cut rows). A distinct count, not
	// an event count: summing across days double-counts members who enable the skill
	// on more than one day, so it is also null in date-range rollup mode
	// (`starting_date`/`ending_date`).
	EnableCount int64 `json:"enable_count" api:"nullable"`
	// Estimated overage spend attributed to this skill, as a decimal string in the
	// minor unit of `currency` (cents for USD; "1250" is $12.50, fractional cents
	// possible) — an allocation of each member's daily post-discount, pre-credit
	// metered overage spend (the same cost basis as the organization's spend reporting
	// and the Cost & Usage API, so per-skill figures are directly comparable; spend
	// with no skill attribution — including any member-day without skill invocations —
	// is not represented, so skill rows sum to at most those totals) across the skills
	// the member used. Overage only: usage covered by included seat allowances bills
	// nothing and allocates $0 here — see `attributed_list_price` for the
	// funding-independent usage-value companion. Claude Code, Cowork, and Office Agent
	// spend use request-level skill attribution; claude.ai chat spend is approximated
	// proportionally to skill-invoking messages. An estimate, not a billing number —
	// and the cost of the requests/messages that involved the skill, not the skill's
	// incremental cost (the same request would still have cost something without the
	// skill active). "0" means no overage spend was attributed; null when spend
	// reporting is not enabled for this organization, on `office_agent` product cuts
	// dated before 2026-06-18 (the Office Agent attribution data-start). Addable
	// across days: date-range rollup mode (`starting_date`/`ending_date`) returns the
	// window's sum. With `group_by[]=user_id` each row carries the user's own
	// attributed spend. On `group_by[]` and `filter[]` shapes both amounts can total
	// below the ungrouped value for the same skill over the same date or range: spend
	// attributed to a member–skill pair with no counted usage on that day is excluded
	// from those cuts.
	EstimatedOverageSpend string `json:"estimated_overage_spend" api:"nullable"`
	// Total number of times this skill was invoked on the requested day (the skill
	// analog of plugin `invocation_count`). Unlike `distinct_user_count` — which
	// answers '# of users' — this is the true '# of uses'. A skill counts as used only
	// when it is explicitly activated — the model (or the user, via the skill's slash
	// command) invokes it, reading its instructions into context as part of that
	// activation. Skills that are merely installed or listed as available, or whose
	// content reaches the context without an activation (preloaded, hook-injected, or
	// read as a plain file), are not counted. Null when invocation reporting is not
	// enabled for this organization. Sum across a date range for total uses in the
	// window — date-range rollup mode (`starting_date`/`ending_date`) returns this sum
	// directly.
	InvocationCount int64 `json:"invocation_count" api:"nullable"`
	// Product that produced this row's activity: one of `chat`, `claude_code`,
	// `cowork`, `office_agent`, or `chat_cowork_unified` (Chat and Cowork unified).
	// These are the canonical Cost & Usage product names; an `office_agent` row's
	// per-surface breakdown is in its `office_metrics`. On `/plugins` only `cowork`,
	// `claude_code` and `chat_cowork_unified` occur (the only surfaces with plugin
	// attribution); on `/artifacts` only `chat`, `claude_code`, `cowork` and
	// `chat_cowork_unified` occur (the surfaces that create artifacts);
	// `/apps/chat/projects` does not support the product dimension (a `product` entry
	// in `group_by[]` or `filter[]` there is rejected). Present only when the request
	// grouped by `product`.
	Product string `json:"product" api:"nullable"`
	// Tagged RBAC group identifier (`rbac_group_...`), matching the spend-limits API
	// spelling. Present only when the request grouped by `rbac_group_id`.
	RBACGroupID string `json:"rbac_group_id" api:"nullable"`
	// Resolved RBAC group display name, alongside `rbac_group_id` when name resolution
	// is available. Null if the group has been deleted or its name could not be
	// resolved; `rbac_group_id` remains the stable key.
	RBACGroupName string `json:"rbac_group_name" api:"nullable"`
	// Skill share status (claude.ai only): one of `private`, `organization`, or
	// `public`. Null for skills used only in Claude Code or Office (no per-skill
	// share-status concept) and when share-status reporting is not yet available for
	// the organization. Filterable via `filter[]=share_status:{value}`.
	//
	// Any of "organization", "private", "public".
	ShareStatus BetaAnalyticsSkillActivityShareStatus `json:"share_status" api:"nullable"`
	// Human-readable display name for rows whose `skill_name` is an opaque skill id
	// (user/organization skill types and plugin-delivered skills, whose user-defined
	// names usage reports generally withhold). Organization-shared skills and skills
	// delivered by the organization's own plugins (its plugin marketplaces and its
	// library) resolve; plugin skill names are shown without their 'plugin:' prefix.
	// The literal 'unknown' bucket row gets a fixed 'Unknown skill' label. For a
	// member's own skill (private or personal-plugin) it is null, except when the
	// skill's owner used it from Claude Code or Cowork in the requested period: then
	// it shows the name that client reported at the time. Apart from that, the names
	// of members' own skills are not disclosed to analytics-key holders. Also null for
	// Anthropic-provided plugin skills (not resolved), for an organization skill or
	// plugin whose name can no longer be found (for example, one since deleted), when
	// `skill_name` is already a display name, or when display-name resolution is not
	// enabled for this organization.
	SkillDisplayName string `json:"skill_display_name" api:"nullable"`
	// Tagged user identifier (e.g. `user_...`). Present only when the request grouped
	// by `user_id`.
	UserID string `json:"user_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ChatMetrics              respjson.Field
		ClaudeCodeMetrics        respjson.Field
		CoworkMetrics            respjson.Field
		DistinctUserCount        respjson.Field
		OfficeMetrics            respjson.Field
		SkillName                respjson.Field
		AttributedListPrice      respjson.Field
		ChatCoworkUnifiedMetrics respjson.Field
		Currency                 respjson.Field
		EnableCount              respjson.Field
		EstimatedOverageSpend    respjson.Field
		InvocationCount          respjson.Field
		Product                  respjson.Field
		RBACGroupID              respjson.Field
		RBACGroupName            respjson.Field
		ShareStatus              respjson.Field
		SkillDisplayName         respjson.Field
		UserID                   respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsSkillActivity) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsSkillActivity) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Skill use recorded while members had Chat and Cowork unified (Cowork's features
// inside claude.ai chat) turned on, split into chat conversations and Cowork
// sessions. A count is null in date-range mode where it cannot be computed.
// Omitted from the response on deployments that do not offer Chat and Cowork
// unified.
type BetaAnalyticsSkillActivityChatCoworkUnifiedMetrics struct {
	// A skill's use in chat conversations recorded while members had Chat and Cowork
	// unified turned on.
	Chat BetaAnalyticsSkillChatCoworkUnifiedChatMetrics `json:"chat" api:"required"`
	// A skill's use in Cowork sessions recorded while members had Chat and Cowork
	// unified turned on.
	Sessions BetaAnalyticsSkillChatCoworkUnifiedSessionsMetrics `json:"sessions" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Chat        respjson.Field
		Sessions    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsSkillActivityChatCoworkUnifiedMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsSkillActivityChatCoworkUnifiedMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Skill share status (claude.ai only): one of `private`, `organization`, or
// `public`. Null for skills used only in Claude Code or Office (no per-skill
// share-status concept) and when share-status reporting is not yet available for
// the organization. Filterable via `filter[]=share_status:{value}`.
type BetaAnalyticsSkillActivityShareStatus string

const (
	BetaAnalyticsSkillActivityShareStatusOrganization BetaAnalyticsSkillActivityShareStatus = "organization"
	BetaAnalyticsSkillActivityShareStatusPrivate      BetaAnalyticsSkillActivityShareStatus = "private"
	BetaAnalyticsSkillActivityShareStatusPublic       BetaAnalyticsSkillActivityShareStatus = "public"
)

// A skill's use in chat conversations recorded while members had Chat and Cowork
// unified turned on.
type BetaAnalyticsSkillChatCoworkUnifiedChatMetrics struct {
	// Same measure as `chat_metrics.distinct_conversation_skill_used_count`, for
	// activity recorded while members had Chat and Cowork unified turned on.
	// Approximate (HLL, typical error <2%) in date-range mode. Null on aggregated rows
	// where a distinct count cannot be computed.
	DistinctConversationSkillUsedCount int64 `json:"distinct_conversation_skill_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctConversationSkillUsedCount respjson.Field
		ExtraFields                        map[string]respjson.Field
		raw                                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsSkillChatCoworkUnifiedChatMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsSkillChatCoworkUnifiedChatMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A skill's use in Cowork sessions recorded while members had Chat and Cowork
// unified turned on.
type BetaAnalyticsSkillChatCoworkUnifiedSessionsMetrics struct {
	// Same measure as `cowork_metrics.distinct_session_skill_used_count`, for activity
	// recorded while members had Chat and Cowork unified turned on. Approximate (HLL,
	// typical error <2%) in date-range mode. Null on aggregated rows where a distinct
	// count cannot be computed.
	DistinctSessionSkillUsedCount int64 `json:"distinct_session_skill_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctSessionSkillUsedCount respjson.Field
		ExtraFields                   map[string]respjson.Field
		raw                           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsSkillChatCoworkUnifiedSessionsMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsSkillChatCoworkUnifiedSessionsMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude.ai activity metrics for a single skill on a given day.
type BetaAnalyticsSkillChatMetrics struct {
	// Number of distinct conversations in which the skill was used. A skill counts as
	// used only when it is explicitly activated — the model (or the user, via the
	// skill's slash command) invokes it, reading its instructions into context as part
	// of that activation. Skills that are merely installed or listed as available, or
	// whose content reaches the context without an activation (preloaded,
	// hook-injected, or read as a plain file), are not counted. Approximate (HLL,
	// typical error <2%) in date-range mode. Null on aggregated rows where a distinct
	// count cannot be computed.
	DistinctConversationSkillUsedCount int64 `json:"distinct_conversation_skill_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctConversationSkillUsedCount respjson.Field
		ExtraFields                        map[string]respjson.Field
		raw                                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsSkillChatMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsSkillChatMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude Code activity metrics for a single skill on a given day.
type BetaAnalyticsSkillClaudeCodeMetrics struct {
	// Number of distinct Claude Code sessions in which the skill was used. A skill
	// counts as used only when it is explicitly activated — the model (or the user,
	// via the skill's slash command) invokes it, reading its instructions into context
	// as part of that activation. Skills that are merely installed or listed as
	// available, or whose content reaches the context without an activation
	// (preloaded, hook-injected, or read as a plain file), are not counted.
	// Approximate (HLL, typical error <2%) in date-range mode. Null on aggregated rows
	// where a distinct count cannot be computed.
	DistinctSessionSkillUsedCount int64 `json:"distinct_session_skill_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctSessionSkillUsedCount respjson.Field
		ExtraFields                   map[string]respjson.Field
		raw                           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsSkillClaudeCodeMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsSkillClaudeCodeMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Cowork activity metrics for a single skill on a given day.
type BetaAnalyticsSkillCoworkMetrics struct {
	// Number of distinct Cowork sessions in which the skill was used. A skill counts
	// as used only when it is explicitly activated — the model (or the user, via the
	// skill's slash command) invokes it, reading its instructions into context as part
	// of that activation. Skills that are merely installed or listed as available, or
	// whose content reaches the context without an activation (preloaded,
	// hook-injected, or read as a plain file), are not counted. Approximate (HLL,
	// typical error <2%) in date-range mode. Null on aggregated rows where a distinct
	// count cannot be computed.
	DistinctSessionSkillUsedCount int64 `json:"distinct_session_skill_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctSessionSkillUsedCount respjson.Field
		ExtraFields                   map[string]respjson.Field
		raw                           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsSkillCoworkMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsSkillCoworkMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Office Agent activity metrics for a single skill on a given day, broken out by
// Office product.
type BetaAnalyticsSkillOfficeMetrics struct {
	// Office Agent activity metrics for a single skill on a given day within one
	// Office product.
	Excel BetaAnalyticsSkillOfficeProductMetrics `json:"excel" api:"required"`
	// Office Agent activity metrics for a single skill on a given day within one
	// Office product.
	Outlook BetaAnalyticsSkillOfficeProductMetrics `json:"outlook" api:"required"`
	// Office Agent activity metrics for a single skill on a given day within one
	// Office product.
	Powerpoint BetaAnalyticsSkillOfficeProductMetrics `json:"powerpoint" api:"required"`
	// Office Agent activity metrics for a single skill on a given day within one
	// Office product.
	Word BetaAnalyticsSkillOfficeProductMetrics `json:"word" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Excel       respjson.Field
		Outlook     respjson.Field
		Powerpoint  respjson.Field
		Word        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsSkillOfficeMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsSkillOfficeMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Office Agent activity metrics for a single skill on a given day within one
// Office product.
type BetaAnalyticsSkillOfficeProductMetrics struct {
	// Number of distinct Office Agent sessions in which the skill was used. A skill
	// counts as used only when it is explicitly activated — the model (or the user,
	// via the skill's slash command) invokes it, reading its instructions into context
	// as part of that activation. Skills that are merely installed or listed as
	// available, or whose content reaches the context without an activation
	// (preloaded, hook-injected, or read as a plain file), are not counted.
	// Approximate (HLL, typical error <2%) in date-range mode. Null on aggregated rows
	// where a distinct count cannot be computed.
	DistinctSessionSkillUsedCount int64 `json:"distinct_session_skill_used_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		DistinctSessionSkillUsedCount respjson.Field
		ExtraFields                   map[string]respjson.Field
		raw                           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsSkillOfficeProductMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsSkillOfficeProductMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAnalyticsTokenType string

const (
	BetaAnalyticsTokenTypeCacheCreationEphemeral1hInputTokens BetaAnalyticsTokenType = "cache_creation.ephemeral_1h_input_tokens"
	BetaAnalyticsTokenTypeCacheCreationEphemeral5mInputTokens BetaAnalyticsTokenType = "cache_creation.ephemeral_5m_input_tokens"
	BetaAnalyticsTokenTypeCacheReadInputTokens                BetaAnalyticsTokenType = "cache_read_input_tokens"
	BetaAnalyticsTokenTypeOutputTokens                        BetaAnalyticsTokenType = "output_tokens"
	BetaAnalyticsTokenTypeUncachedInputTokens                 BetaAnalyticsTokenType = "uncached_input_tokens"
)

// Accepted/rejected counts for a single Claude Code tool type.
type BetaAnalyticsToolActionCounts struct {
	// Number of tool proposals accepted
	AcceptedCount int64 `json:"accepted_count" api:"required"`
	// Number of tool proposals rejected
	RejectedCount int64 `json:"rejected_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AcceptedCount respjson.Field
		RejectedCount respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsToolActionCounts) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsToolActionCounts) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Per-tool accepted/rejected counts for Claude Code file modification tools.
type BetaAnalyticsToolActions struct {
	// Accepted/rejected counts for a single Claude Code tool type.
	EditTool BetaAnalyticsToolActionCounts `json:"edit_tool" api:"required"`
	// Accepted/rejected counts for a single Claude Code tool type.
	MultiEditTool BetaAnalyticsToolActionCounts `json:"multi_edit_tool" api:"required"`
	// Accepted/rejected counts for a single Claude Code tool type.
	NotebookEditTool BetaAnalyticsToolActionCounts `json:"notebook_edit_tool" api:"required"`
	// Accepted/rejected counts for a single Claude Code tool type.
	WriteTool BetaAnalyticsToolActionCounts `json:"write_tool" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EditTool         respjson.Field
		MultiEditTool    respjson.Field
		NotebookEditTool respjson.Field
		WriteTool        respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsToolActions) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsToolActions) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAnalyticsUsageBucketedResult struct {
	// The number of input tokens for cache creation.
	CacheCreation BetaCacheCreation `json:"cache_creation" api:"required"`
	// The number of input tokens read from the cache.
	CacheReadInputTokens int64 `json:"cache_read_input_tokens" api:"required"`
	// Claude Tag (Claude in Slack) spend category: `engaged` (a person addressed
	// Claude in a channel or thread), `proactive` (Claude responded without being
	// addressed), `scheduled` (a scheduled routine ran), `monitoring` (Claude watching
	// a channel it was asked to monitor), or `dm` (direct messages with Claude).
	// Populated only when `claude_tag_category` is in `group_by[]`; null for usage
	// that is not Claude Tag. Direct-message usage is billed to the individual user
	// and is reported under that user's product, not under `claude-tag`. New
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
	// Inference region of the usage or cost. Null unless `inference_geo` is in
	// `group_by[]`; it can also be null on grouped rows where the region is not set
	// (the rows that `inference_geos[]=not_available` matches).
	//
	// Any of "global", "us".
	InferenceGeo BetaAnalyticsUsageBucketedResultInferenceGeo `json:"inference_geo" api:"required"`
	// Model that produced the usage or cost, as a model name in the form the
	// `models[]` filter accepts (for example, `claude-opus-5`). Null unless `model` is
	// in `group_by[]`; it can also be null on grouped rows whose usage or cost is not
	// attributed to a specific model, such as code execution.
	Model string `json:"model" api:"required"`
	// The number of output tokens generated.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// Product surface that produced the usage or cost. Null unless product is in
	// `group_by[]`; it can also be null on grouped rows whose usage cannot be
	// attributed to a known surface. Values include `chat`, `claude_code`, `cowork`,
	// `office_agent`, `claude_in_chrome`, `claude_design`, `claude-tag`, and
	// `chat_cowork_unified`. `claude-tag` is Claude Tag, the Claude product in Slack.
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
	// Number of API requests in this row's scope. For sandbox / code-execution events,
	// this counts execution spans rather than HTTP requests (these rows surface with
	// `product: null`).
	Requests int64 `json:"requests" api:"required"`
	// Server-side tool usage metrics.
	ServerToolUse BetaAnalyticsServerToolUse `json:"server_tool_use" api:"required"`
	// Slack channel the usage originated from. Populated only when `slack_channel_id`
	// is in `group_by[]`; null for usage outside Slack (and for rows recorded before
	// channel attribution was enabled).
	SlackChannelID string `json:"slack_channel_id" api:"required"`
	// Inference speed mode of the usage or cost: `fast` or `standard`. Null unless
	// `speed` is in `group_by[]`.
	//
	// Any of "fast", "standard".
	Speed BetaAnalyticsUsageBucketedResultSpeed `json:"speed" api:"required"`
	// The number of uncached input tokens processed.
	UncachedInputTokens int64 `json:"uncached_input_tokens" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CacheCreation        respjson.Field
		CacheReadInputTokens respjson.Field
		ClaudeTagCategory    respjson.Field
		ClaudeTagUserID      respjson.Field
		ContextWindow        respjson.Field
		InferenceGeo         respjson.Field
		Model                respjson.Field
		OutputTokens         respjson.Field
		Product              respjson.Field
		RBACGroupID          respjson.Field
		Requests             respjson.Field
		ServerToolUse        respjson.Field
		SlackChannelID       respjson.Field
		Speed                respjson.Field
		UncachedInputTokens  respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsUsageBucketedResult) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsUsageBucketedResult) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inference region of the usage or cost. Null unless `inference_geo` is in
// `group_by[]`; it can also be null on grouped rows where the region is not set
// (the rows that `inference_geos[]=not_available` matches).
type BetaAnalyticsUsageBucketedResultInferenceGeo string

const (
	BetaAnalyticsUsageBucketedResultInferenceGeoGlobal BetaAnalyticsUsageBucketedResultInferenceGeo = "global"
	BetaAnalyticsUsageBucketedResultInferenceGeoUs     BetaAnalyticsUsageBucketedResultInferenceGeo = "us"
)

// Inference speed mode of the usage or cost: `fast` or `standard`. Null unless
// `speed` is in `group_by[]`.
type BetaAnalyticsUsageBucketedResultSpeed string

const (
	BetaAnalyticsUsageBucketedResultSpeedFast     BetaAnalyticsUsageBucketedResultSpeed = "fast"
	BetaAnalyticsUsageBucketedResultSpeedStandard BetaAnalyticsUsageBucketedResultSpeed = "standard"
)

type BetaAnalyticsUsageReportTimeBucket struct {
	// End of the time bucket (exclusive) in RFC 3339 format.
	EndingAt time.Time `json:"ending_at" api:"required" format:"date-time"`
	// Rows for this time bucket. Empty when the bucket has no data; otherwise a single
	// combined row when `group_by[]` is omitted, or one row per group (subject to the
	// per-bucket group cap described on the `group_by[]` parameter).
	Results []BetaAnalyticsUsageBucketedResult `json:"results" api:"required"`
	// Start of the time bucket (inclusive) in RFC 3339 format.
	StartingAt time.Time `json:"starting_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EndingAt    respjson.Field
		Results     respjson.Field
		StartingAt  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsUsageReportTimeBucket) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsUsageReportTimeBucket) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAnalyticsUsageUsersItem struct {
	// The user this row's usage or cost is attributed to. Always a `user_actor`.
	Actor BetaAnalyticsUsageUsersItemActorUnion `json:"actor" api:"required"`
	// The number of input tokens for cache creation.
	CacheCreation BetaCacheCreation `json:"cache_creation" api:"required"`
	// The number of input tokens read from the cache.
	CacheReadInputTokens int64 `json:"cache_read_input_tokens" api:"required"`
	// Claude Tag (Claude in Slack) spend category: `engaged` (a person addressed
	// Claude in a channel or thread), `proactive` (Claude responded without being
	// addressed), `scheduled` (a scheduled routine ran), `monitoring` (Claude watching
	// a channel it was asked to monitor), or `dm` (direct messages with Claude).
	// Populated only when `claude_tag_category` is in `group_by[]`; null for usage
	// that is not Claude Tag. Direct-message usage is billed to the individual user
	// and is reported under that user's product, not under `claude-tag`. New
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
	// End of the row's UTC time bucket (exclusive), as an RFC 3339 timestamp; equal to
	// `starting_at` plus one `bucket_width`. Null unless `bucket_width` is set.
	EndingAt time.Time `json:"ending_at" api:"required" format:"date-time"`
	// Inference region of the usage or cost. Null unless `inference_geo` is in
	// `group_by[]`; it can also be null on grouped rows where the region is not set
	// (the rows that `inference_geos[]=not_available` matches).
	//
	// Any of "global", "us".
	InferenceGeo BetaAnalyticsUsageUsersItemInferenceGeo `json:"inference_geo" api:"required"`
	// Model that produced the usage or cost, as a model name in the form the
	// `models[]` filter accepts (for example, `claude-opus-5`). Null unless `model` is
	// in `group_by[]`; it can also be null on grouped rows whose usage or cost is not
	// attributed to a specific model, such as code execution.
	Model string `json:"model" api:"required"`
	// The number of output tokens generated.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// Product surface that produced the usage or cost. Null unless product is in
	// `group_by[]`; it can also be null on grouped rows whose usage cannot be
	// attributed to a known surface. Values include `chat`, `claude_code`, `cowork`,
	// `office_agent`, `claude_in_chrome`, `claude_design`, `claude-tag`, and
	// `chat_cowork_unified`. `claude-tag` is Claude Tag, the Claude product in Slack.
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
	// Number of API requests in this row's scope. For sandbox / code-execution events,
	// this counts execution spans rather than HTTP requests (these rows surface with
	// `product: null`).
	Requests int64 `json:"requests" api:"required"`
	// Server-side tool usage metrics.
	ServerToolUse BetaAnalyticsServerToolUse `json:"server_tool_use" api:"required"`
	// Slack channel the usage originated from. Populated only when `slack_channel_id`
	// is in `group_by[]`; null for usage outside Slack (and for rows recorded before
	// channel attribution was enabled).
	SlackChannelID string `json:"slack_channel_id" api:"required"`
	// Inference speed mode of the usage or cost: `fast` or `standard`. Null unless
	// `speed` is in `group_by[]`.
	//
	// Any of "fast", "standard".
	Speed BetaAnalyticsUsageUsersItemSpeed `json:"speed" api:"required"`
	// Start of the row's UTC time bucket (inclusive), as an RFC 3339 timestamp. Null
	// unless `bucket_width` is set; without `bucket_width`, each row aggregates the
	// full requested range.
	StartingAt time.Time `json:"starting_at" api:"required" format:"date-time"`
	// Total token count across all token types. This is the value the default
	// `order_by` (`total_tokens`) sorts on.
	TotalTokens int64 `json:"total_tokens" api:"required"`
	// The number of uncached input tokens processed.
	UncachedInputTokens int64 `json:"uncached_input_tokens" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Actor                respjson.Field
		CacheCreation        respjson.Field
		CacheReadInputTokens respjson.Field
		ClaudeTagCategory    respjson.Field
		ClaudeTagUserID      respjson.Field
		ContextWindow        respjson.Field
		EndingAt             respjson.Field
		InferenceGeo         respjson.Field
		Model                respjson.Field
		OutputTokens         respjson.Field
		Product              respjson.Field
		RBACGroupID          respjson.Field
		Requests             respjson.Field
		ServerToolUse        respjson.Field
		SlackChannelID       respjson.Field
		Speed                respjson.Field
		StartingAt           respjson.Field
		TotalTokens          respjson.Field
		UncachedInputTokens  respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsUsageUsersItem) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsUsageUsersItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// BetaAnalyticsUsageUsersItemActorUnion contains all possible properties and
// values from [BetaAnalyticsUserActor].
//
// Use the [BetaAnalyticsUsageUsersItemActorUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type BetaAnalyticsUsageUsersItemActorUnion struct {
	// This field is from variant [BetaAnalyticsUserActor].
	Deleted bool `json:"deleted"`
	// This field is from variant [BetaAnalyticsUserActor].
	EmailAddress string `json:"email_address"`
	// This field is from variant [BetaAnalyticsUserActor].
	Name string `json:"name"`
	// Any of "user_actor".
	Type string `json:"type"`
	// This field is from variant [BetaAnalyticsUserActor].
	UserID string `json:"user_id"`
	JSON   struct {
		Deleted      respjson.Field
		EmailAddress respjson.Field
		Name         respjson.Field
		Type         respjson.Field
		UserID       respjson.Field
		raw          string
	} `json:"-"`
}

// anyBetaAnalyticsUsageUsersItemActor is implemented by each variant of
// [BetaAnalyticsUsageUsersItemActorUnion] to add type safety for the return type
// of [BetaAnalyticsUsageUsersItemActorUnion.AsAny]
type anyBetaAnalyticsUsageUsersItemActor interface {
	implBetaAnalyticsUsageUsersItemActorUnion()
}

func (BetaAnalyticsUserActor) implBetaAnalyticsUsageUsersItemActorUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := BetaAnalyticsUsageUsersItemActorUnion.AsAny().(type) {
//	case anthropic.BetaAnalyticsUserActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u BetaAnalyticsUsageUsersItemActorUnion) AsAny() anyBetaAnalyticsUsageUsersItemActor {
	switch u.Type {
	case "user_actor":
		return u.AsUserActor()
	}
	return nil
}

func (u BetaAnalyticsUsageUsersItemActorUnion) AsUserActor() (v BetaAnalyticsUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u BetaAnalyticsUsageUsersItemActorUnion) RawJSON() string { return u.JSON.raw }

func (r *BetaAnalyticsUsageUsersItemActorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inference region of the usage or cost. Null unless `inference_geo` is in
// `group_by[]`; it can also be null on grouped rows where the region is not set
// (the rows that `inference_geos[]=not_available` matches).
type BetaAnalyticsUsageUsersItemInferenceGeo string

const (
	BetaAnalyticsUsageUsersItemInferenceGeoGlobal BetaAnalyticsUsageUsersItemInferenceGeo = "global"
	BetaAnalyticsUsageUsersItemInferenceGeoUs     BetaAnalyticsUsageUsersItemInferenceGeo = "us"
)

// Inference speed mode of the usage or cost: `fast` or `standard`. Null unless
// `speed` is in `group_by[]`.
type BetaAnalyticsUsageUsersItemSpeed string

const (
	BetaAnalyticsUsageUsersItemSpeedFast     BetaAnalyticsUsageUsersItemSpeed = "fast"
	BetaAnalyticsUsageUsersItemSpeedStandard BetaAnalyticsUsageUsersItemSpeed = "standard"
)

// A user in the organization, identified by tagged id and email address.
type BetaAnalyticsUser struct {
	// Tagged user identifier (e.g. `user_...`)
	ID string `json:"id" api:"required"`
	// Email address of the user
	EmailAddress string `json:"email_address" api:"required"`
	// Object type. Always `user`.
	Type constant.User `json:"type" default:"user"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		EmailAddress respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsUser) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsUser) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Per-user activity data for a given day.
type BetaAnalyticsUserActivity struct {
	// Claude.ai activity metrics for a single user on a given day.
	ChatMetrics BetaAnalyticsChatMetrics `json:"chat_metrics" api:"required"`
	// Claude Code activity metrics for a single user on a given day.
	ClaudeCodeMetrics BetaAnalyticsClaudeCodeMetrics `json:"claude_code_metrics" api:"required"`
	// Cowork activity metrics for a single user on a given day.
	CoworkMetrics BetaAnalyticsCoworkMetrics `json:"cowork_metrics" api:"required"`
	// Claude Design activity metrics for a single user on a given day.
	DesignMetrics BetaAnalyticsDesignMetrics `json:"design_metrics" api:"required"`
	// Office Agent activity metrics for a single user on a given day, broken out by
	// Office product.
	OfficeMetrics BetaAnalyticsOfficeMetrics `json:"office_metrics" api:"required"`
	// Claude Science activity metrics for a single user on a given day.
	ScienceMetrics BetaAnalyticsScienceMetrics `json:"science_metrics" api:"required"`
	// Number of web searches performed
	WebSearchCount int64 `json:"web_search_count" api:"required"`
	// Activity recorded while the member had Chat and Cowork unified (Cowork's
	// features inside claude.ai chat) turned on, split into `chat` (chat activity) and
	// `sessions` (Cowork activity). Omitted from the response on deployments that do
	// not offer Chat and Cowork unified.
	ChatCoworkUnifiedMetrics BetaAnalyticsUserActivityChatCoworkUnifiedMetrics `json:"chat_cowork_unified_metrics" api:"nullable"`
	// Number of distinct active users represented by this row. Only set for grouped
	// rollups (`group_by[]`); null for per-user rows. In date-range mode, recomputed
	// as an exact distinct count of the group's active members over the requested
	// window, never a sum of per-day values.
	DistinctUserCount int64 `json:"distinct_user_count" api:"nullable"`
	// Most recent UTC day (YYYY-MM-DD) on which the user had any counted activity,
	// within the requested window: equal to the requested `date` in single-day mode,
	// and to the latest active day from `starting_date` (inclusive) to `ending_date`
	// (exclusive) in date-range rollup mode — never a day earlier than the window
	// start. On filtered requests (`filter[]`) only days matching the filter count:
	// with `filter[]=rbac_group_id:{id}` it is the last day the user was active while
	// a member of that group, consistent with the row's other metrics. On grouped
	// (`group_by[]`) rows it is the latest day any member of the group was active (the
	// requested `date` in single-day mode). Omitted from the response while
	// last-activity reporting is not enabled for this organization.
	LastActivityDate time.Time `json:"last_activity_date" api:"nullable" format:"date"`
	// Tagged RBAC group identifier (`rbac_group_...`), matching the spend-limits API
	// spelling. Present only when the request grouped by `rbac_group_id`.
	RBACGroupID string `json:"rbac_group_id" api:"nullable"`
	// Resolved RBAC group display name, alongside `rbac_group_id` when name resolution
	// is available. Null if the group has been deleted or its name could not be
	// resolved; `rbac_group_id` remains the stable key.
	RBACGroupName string `json:"rbac_group_name" api:"nullable"`
	// The user this row describes. Null on rows aggregated across users.
	User BetaAnalyticsUser `json:"user" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ChatMetrics              respjson.Field
		ClaudeCodeMetrics        respjson.Field
		CoworkMetrics            respjson.Field
		DesignMetrics            respjson.Field
		OfficeMetrics            respjson.Field
		ScienceMetrics           respjson.Field
		WebSearchCount           respjson.Field
		ChatCoworkUnifiedMetrics respjson.Field
		DistinctUserCount        respjson.Field
		LastActivityDate         respjson.Field
		RBACGroupID              respjson.Field
		RBACGroupName            respjson.Field
		User                     respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsUserActivity) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsUserActivity) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Activity recorded while the member had Chat and Cowork unified (Cowork's
// features inside claude.ai chat) turned on, split into `chat` (chat activity) and
// `sessions` (Cowork activity). Omitted from the response on deployments that do
// not offer Chat and Cowork unified.
type BetaAnalyticsUserActivityChatCoworkUnifiedMetrics struct {
	// Chat activity recorded while members had Chat and Cowork unified turned on.
	Chat BetaAnalyticsChatCoworkUnifiedChatMetrics `json:"chat" api:"required"`
	// Cowork session activity recorded while members had Chat and Cowork unified
	// turned on.
	Sessions BetaAnalyticsChatCoworkUnifiedSessionsMetrics `json:"sessions" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Chat        respjson.Field
		Sessions    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsUserActivityChatCoworkUnifiedMetrics) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsUserActivityChatCoworkUnifiedMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BetaAnalyticsUserActor struct {
	// True when the account has been deleted, or when the user is no longer a member
	// of the organization or its associated organizations (for example, their
	// membership was removed or they were deprovisioned via your identity provider).
	// `email_address` stays populated for removed users and is null when the account
	// has been deleted. `name` follows the rules described on that field. The
	// `user_id` is still populated for reconciliation.
	Deleted bool `json:"deleted" api:"required"`
	// The user's email address, including for users who are no longer members of the
	// organization or its associated organizations. Null when the account has been
	// deleted (check `deleted`) and for system-minted service accounts, which have no
	// person's mailbox behind them (check `name`).
	EmailAddress string `json:"email_address" api:"required"`
	// The user's full name. Null when the user has not set a name. Returns
	// `"Deleted User"` when the account itself has been deleted, or when the user is
	// no longer a member of the organization or its associated organizations and the
	// organization has chosen to hide the names of removed users. Otherwise, the name
	// stays populated for removed users. Rows for system-minted service accounts
	// render the service name (for example, `"Claude Security"` for usage by
	// Anthropic's security-patching service) or null.
	Name string `json:"name" api:"required"`
	// Actor type. Always `"user_actor"`.
	Type constant.UserActor `json:"type" default:"user_actor"`
	// Tagged user ID.
	UserID string `json:"user_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Deleted      respjson.Field
		EmailAddress respjson.Field
		Name         respjson.Field
		Type         respjson.Field
		UserID       respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r BetaAnalyticsUserActor) RawJSON() string { return r.JSON.raw }
func (r *BetaAnalyticsUserActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
