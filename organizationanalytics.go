package anthropic

import (
	"time"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/option"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared/constant"
)

// OrganizationAnalyticsService contains methods and other services that help with
// interacting with the anthropic API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOrganizationAnalyticsService] method instead.
type OrganizationAnalyticsService struct {
	Options    []option.RequestOption
	Summaries  OrganizationAnalyticsSummaryService
	Users      OrganizationAnalyticsUserService
	Apps       OrganizationAnalyticsAppService
	Connectors OrganizationAnalyticsConnectorService
	Plugins    OrganizationAnalyticsPluginService
	Skills     OrganizationAnalyticsSkillService
	Artifacts  OrganizationAnalyticsArtifactService
}

// NewOrganizationAnalyticsService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewOrganizationAnalyticsService(opts ...option.RequestOption) (r OrganizationAnalyticsService) {
	r = OrganizationAnalyticsService{}
	r.Options = opts
	r.Summaries = NewOrganizationAnalyticsSummaryService(opts...)
	r.Users = NewOrganizationAnalyticsUserService(opts...)
	r.Apps = NewOrganizationAnalyticsAppService(opts...)
	r.Connectors = NewOrganizationAnalyticsConnectorService(opts...)
	r.Plugins = NewOrganizationAnalyticsPluginService(opts...)
	r.Skills = NewOrganizationAnalyticsSkillService(opts...)
	r.Artifacts = NewOrganizationAnalyticsArtifactService(opts...)
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
type AnalyticsArtifactActivity struct {
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
func (r AnalyticsArtifactActivity) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsArtifactActivity) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A dimension to group artifact activity by.
type AnalyticsArtifactGroupBy string

const (
	AnalyticsArtifactGroupByProduct     AnalyticsArtifactGroupBy = "product"
	AnalyticsArtifactGroupByRBACGroupID AnalyticsArtifactGroupBy = "rbac_group_id"
	AnalyticsArtifactGroupByUserID      AnalyticsArtifactGroupBy = "user_id"
)

// Chat activity recorded while members had Chat and Cowork unified turned on.
type AnalyticsChatCoworkUnifiedChatMetrics struct {
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
func (r AnalyticsChatCoworkUnifiedChatMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsChatCoworkUnifiedChatMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Activity in Chat and Cowork unified for one user on one day: `chat` holds chat
// activity and `sessions` holds Cowork session activity.
type AnalyticsChatCoworkUnifiedMetrics struct {
	// Chat activity recorded while members had Chat and Cowork unified turned on.
	Chat AnalyticsChatCoworkUnifiedChatMetrics `json:"chat" api:"required"`
	// Cowork session activity recorded while members had Chat and Cowork unified
	// turned on.
	Sessions AnalyticsChatCoworkUnifiedSessionsMetrics `json:"sessions" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Chat        respjson.Field
		Sessions    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AnalyticsChatCoworkUnifiedMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsChatCoworkUnifiedMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Cowork session activity recorded while members had Chat and Cowork unified
// turned on.
type AnalyticsChatCoworkUnifiedSessionsMetrics struct {
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
func (r AnalyticsChatCoworkUnifiedSessionsMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsChatCoworkUnifiedSessionsMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude.ai activity metrics for a single user on a given day.
type AnalyticsChatMetrics struct {
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
func (r AnalyticsChatMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsChatMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude Code activity metrics for a single user on a given day.
type AnalyticsClaudeCodeMetrics struct {
	// Core Claude Code activity metrics for a single user on a given day.
	CoreMetrics AnalyticsCoreCodeMetrics `json:"core_metrics" api:"required"`
	// Per-tool accepted/rejected counts for Claude Code file modification tools.
	ToolActions AnalyticsToolActions `json:"tool_actions" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CoreMetrics respjson.Field
		ToolActions respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AnalyticsClaudeCodeMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsClaudeCodeMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude Tag activity metrics for a single linked user on a given day.
type AnalyticsClaudeTagMetrics struct {
	// 1 on a day the user was active in Claude Tag in Slack, otherwise 0. Summed over
	// the window in date-range mode, so it is the number of active days; on grouped
	// rows, the total of each member's active days. Null when Claude Tag data cannot
	// be read for the request. Otherwise, 0 on a row filtered by project, and null on
	// any other row when no covered day was measured.
	ActiveDayCount int64 `json:"active_day_count" api:"required"`
	// Number of the user's Slack messages that made them active: @-mentions, direct
	// messages, group direct message asks, and thread posts Claude answered. Other
	// messages are not counted. Null when Claude Tag data cannot be read for the
	// request. Otherwise, 0 on a row filtered by project, and null on any other row
	// when no covered day was measured.
	ActiveMessageCount int64 `json:"active_message_count" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ActiveDayCount     respjson.Field
		ActiveMessageCount respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AnalyticsClaudeTagMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsClaudeTagMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Per-connector activity data for a given day.
type AnalyticsConnectorActivity struct {
	// Claude.ai activity metrics for a single connector on a given day.
	ChatMetrics AnalyticsConnectorChatMetrics `json:"chat_metrics" api:"required"`
	// Claude Code activity metrics for a single connector on a given day.
	ClaudeCodeMetrics AnalyticsConnectorClaudeCodeMetrics `json:"claude_code_metrics" api:"required"`
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
	ConnectorDisplayName string `json:"connector_display_name" api:"required"`
	// Name of the connector. Some rows carry an opaque connector id here instead of a
	// readable name; `connector_display_name` holds the resolved name for those rows.
	ConnectorName string `json:"connector_name" api:"required"`
	// Cowork activity metrics for a single connector on a given day.
	CoworkMetrics AnalyticsConnectorCoworkMetrics `json:"cowork_metrics" api:"required"`
	// Number of distinct users who used the connector on the requested day, or, in
	// date-range mode, over the requested window — recomputed as an exact distinct
	// count over the window's per-member daily rows, never a sum of per-day values.
	DistinctUserCount int64 `json:"distinct_user_count" api:"required"`
	// Number of distinct users whose use of this connector on the requested day ran on
	// their own individual credential, connected through their own consent flow.
	// Companion bucket to `managed_auth_distinct_user_count`, which carries the
	// measurement, attribution, and null rules. Users whose requests used no stored
	// credential count in neither bucket.
	IndividualAuthDistinctUserCount int64 `json:"individual_auth_distinct_user_count" api:"required"`
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
	ManagedAuthDistinctUserCount int64 `json:"managed_auth_distinct_user_count" api:"required"`
	// Office Agent activity metrics for a single connector on a given day, broken out
	// by Office product.
	OfficeMetrics AnalyticsConnectorOfficeMetrics `json:"office_metrics" api:"required"`
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
	ReadCallCount int64 `json:"read_call_count" api:"required"`
	// Number of connector tool calls on the requested day with no trusted read-only
	// annotation — the annotation is optional in the MCP spec and is discarded when
	// connector access controls are active, so unclassified calls are common. This
	// field shows how much of the day's classified activity the read/write split
	// actually covers. Call count, not distinct users. One of the three
	// call-classification buckets; see `read_call_count` for the per-surface
	// data-start dates, null conditions, and date-range guidance.
	UnclassifiedCallCount int64 `json:"unclassified_call_count" api:"required"`
	// Number of connector tool calls on the requested day whose trusted read-only
	// annotation marked them not read-only. Call count, not distinct users. One of the
	// three call-classification buckets; see `read_call_count` for the per-surface
	// data-start dates, null conditions, and date-range guidance.
	WriteCallCount int64 `json:"write_call_count" api:"required"`
	// Connector use recorded while members had Chat and Cowork unified (Cowork's
	// features inside claude.ai chat) turned on, split into chat conversations and
	// Cowork sessions. A count is null in date-range mode where it cannot be computed.
	// Omitted from the response where Chat and Cowork unified is not available.
	ChatCoworkUnifiedMetrics AnalyticsConnectorChatCoworkUnifiedMetrics `json:"chat_cowork_unified_metrics" api:"nullable"`
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
		ChatMetrics                     respjson.Field
		ClaudeCodeMetrics               respjson.Field
		ConnectorDisplayName            respjson.Field
		ConnectorName                   respjson.Field
		CoworkMetrics                   respjson.Field
		DistinctUserCount               respjson.Field
		IndividualAuthDistinctUserCount respjson.Field
		ManagedAuthDistinctUserCount    respjson.Field
		OfficeMetrics                   respjson.Field
		ReadCallCount                   respjson.Field
		UnclassifiedCallCount           respjson.Field
		WriteCallCount                  respjson.Field
		ChatCoworkUnifiedMetrics        respjson.Field
		Product                         respjson.Field
		RBACGroupID                     respjson.Field
		RBACGroupName                   respjson.Field
		UserID                          respjson.Field
		ExtraFields                     map[string]respjson.Field
		raw                             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AnalyticsConnectorActivity) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsConnectorActivity) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A connector's use in chat conversations recorded while members had Chat and
// Cowork unified turned on.
type AnalyticsConnectorChatCoworkUnifiedChatMetrics struct {
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
func (r AnalyticsConnectorChatCoworkUnifiedChatMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsConnectorChatCoworkUnifiedChatMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A connector's use in Chat and Cowork unified on one day: `chat` holds its use in
// chat conversations and `sessions` its use in Cowork sessions.
type AnalyticsConnectorChatCoworkUnifiedMetrics struct {
	// A connector's use in chat conversations recorded while members had Chat and
	// Cowork unified turned on.
	Chat AnalyticsConnectorChatCoworkUnifiedChatMetrics `json:"chat" api:"required"`
	// A connector's use in Cowork sessions recorded while members had Chat and Cowork
	// unified turned on.
	Sessions AnalyticsConnectorChatCoworkUnifiedSessionsMetrics `json:"sessions" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Chat        respjson.Field
		Sessions    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AnalyticsConnectorChatCoworkUnifiedMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsConnectorChatCoworkUnifiedMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A connector's use in Cowork sessions recorded while members had Chat and Cowork
// unified turned on.
type AnalyticsConnectorChatCoworkUnifiedSessionsMetrics struct {
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
func (r AnalyticsConnectorChatCoworkUnifiedSessionsMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsConnectorChatCoworkUnifiedSessionsMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude.ai activity metrics for a single connector on a given day.
type AnalyticsConnectorChatMetrics struct {
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
func (r AnalyticsConnectorChatMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsConnectorChatMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude Code activity metrics for a single connector on a given day.
type AnalyticsConnectorClaudeCodeMetrics struct {
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
func (r AnalyticsConnectorClaudeCodeMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsConnectorClaudeCodeMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Cowork activity metrics for a single connector on a given day.
type AnalyticsConnectorCoworkMetrics struct {
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
func (r AnalyticsConnectorCoworkMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsConnectorCoworkMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A dimension to group connector activity by.
type AnalyticsConnectorGroupBy string

const (
	AnalyticsConnectorGroupByProduct     AnalyticsConnectorGroupBy = "product"
	AnalyticsConnectorGroupByRBACGroupID AnalyticsConnectorGroupBy = "rbac_group_id"
	AnalyticsConnectorGroupByUserID      AnalyticsConnectorGroupBy = "user_id"
)

// Office Agent activity metrics for a single connector on a given day, broken out
// by Office product.
type AnalyticsConnectorOfficeMetrics struct {
	// Office Agent activity metrics for a single connector on a given day within one
	// Office product.
	Excel AnalyticsConnectorOfficeProductMetrics `json:"excel" api:"required"`
	// Office Agent activity metrics for a single connector on a given day within one
	// Office product.
	Outlook AnalyticsConnectorOfficeProductMetrics `json:"outlook" api:"required"`
	// Office Agent activity metrics for a single connector on a given day within one
	// Office product.
	Powerpoint AnalyticsConnectorOfficeProductMetrics `json:"powerpoint" api:"required"`
	// Office Agent activity metrics for a single connector on a given day within one
	// Office product.
	Word AnalyticsConnectorOfficeProductMetrics `json:"word" api:"required"`
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
func (r AnalyticsConnectorOfficeMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsConnectorOfficeMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Office Agent activity metrics for a single connector on a given day within one
// Office product.
type AnalyticsConnectorOfficeProductMetrics struct {
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
func (r AnalyticsConnectorOfficeProductMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsConnectorOfficeProductMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Core Claude Code activity metrics for a single user on a given day.
type AnalyticsCoreCodeMetrics struct {
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
	LinesOfCode AnalyticsLinesOfCode `json:"lines_of_code" api:"required"`
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
func (r AnalyticsCoreCodeMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsCoreCodeMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Cowork activity metrics for a single user on a given day.
type AnalyticsCoworkMetrics struct {
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
	// Number of distinct plugins used in Cowork sessions. Null while Cowork plugin-use
	// metrics are not enabled for this organization. Approximate (HLL, typical error
	// <2%) in date-range mode. Null on aggregated rows where a distinct count cannot
	// be computed.
	DistinctPluginsUsedCount int64 `json:"distinct_plugins_used_count" api:"required"`
	// Number of distinct Cowork sessions. Approximate (HLL, typical error <2%) in
	// date-range mode. Null on aggregated rows where a distinct count cannot be
	// computed.
	DistinctSessionCount int64 `json:"distinct_session_count" api:"required"`
	// Number of distinct skills used in Cowork sessions. Approximate (HLL, typical
	// error <2%) in date-range mode. Null on aggregated rows where a distinct count
	// cannot be computed.
	DistinctSkillsUsedCount int64 `json:"distinct_skills_used_count" api:"required"`
	// Number of successful Edit tool calls in Cowork sessions. Null while the
	// file-edit metrics are not enabled for this organization.
	EditToolCount int64 `json:"edit_tool_count" api:"required"`
	// Number of successful file-edit tool calls (Edit, MultiEdit, Write, NotebookEdit)
	// in Cowork sessions. Null, never 0, while the file-edit metrics are not enabled
	// for this organization.
	FileEditCount int64 `json:"file_edit_count" api:"required"`
	// Number of messages sent in Cowork sessions
	MessageCount int64 `json:"message_count" api:"required"`
	// Number of successful MultiEdit tool calls in Cowork sessions. Null while the
	// file-edit metrics are not enabled for this organization.
	MultiEditToolCount int64 `json:"multi_edit_tool_count" api:"required"`
	// Number of successful NotebookEdit tool calls in Cowork sessions. Null while the
	// file-edit metrics are not enabled for this organization.
	NotebookEditToolCount int64 `json:"notebook_edit_tool_count" api:"required"`
	// Total number of plugin invocations in Cowork sessions. Null while Cowork
	// plugin-use metrics are not enabled for this organization.
	PluginsUsedCount int64 `json:"plugins_used_count" api:"required"`
	// Number of distinct Cowork sessions with at least one successful file-edit tool
	// call. Null while the file-edit metrics are not enabled for this organization.
	// Approximate (HLL, typical error <2%) in date-range mode. Null on aggregated rows
	// where a distinct count cannot be computed.
	SessionsWithFileEditsCount int64 `json:"sessions_with_file_edits_count" api:"required"`
	// Total number of skill invocations in Cowork sessions
	SkillsUsedCount int64 `json:"skills_used_count" api:"required"`
	// Number of successful Write tool calls in Cowork sessions. Null while the
	// file-edit metrics are not enabled for this organization.
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
func (r AnalyticsCoworkMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsCoworkMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude Design activity metrics for a single user on a given day.
type AnalyticsDesignMetrics struct {
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
func (r AnalyticsDesignMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsDesignMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Lines of code added and removed via Claude Code.
type AnalyticsLinesOfCode struct {
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
func (r AnalyticsLinesOfCode) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsLinesOfCode) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Office Agent activity metrics for a single user on a given day, broken out by
// Office product.
type AnalyticsOfficeMetrics struct {
	// Office Agent activity metrics for a single user on a given day within one Office
	// product.
	Excel AnalyticsOfficeProductMetrics `json:"excel" api:"required"`
	// Office Agent activity metrics for a single user on a given day within one Office
	// product.
	Outlook AnalyticsOfficeProductMetrics `json:"outlook" api:"required"`
	// Office Agent activity metrics for a single user on a given day within one Office
	// product.
	Powerpoint AnalyticsOfficeProductMetrics `json:"powerpoint" api:"required"`
	// Office Agent activity metrics for a single user on a given day within one Office
	// product.
	Word AnalyticsOfficeProductMetrics `json:"word" api:"required"`
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
func (r AnalyticsOfficeMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsOfficeMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Office Agent activity metrics for a single user on a given day within one Office
// product.
type AnalyticsOfficeProductMetrics struct {
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
func (r AnalyticsOfficeProductMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsOfficeProductMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Sort direction: `asc` or `desc`. Defaults to `asc` for the endpoint's sort
// column and to `desc` when `order_by` names a metric (a top-N ranking). Applies
// to `order_by`, or to the endpoint's default sort field when `order_by` is
// omitted.
type AnalyticsOrder string

const (
	AnalyticsOrderAsc  AnalyticsOrder = "asc"
	AnalyticsOrderDesc AnalyticsOrder = "desc"
)

// Per-plugin install + invocation activity for a given day.
//
// With `group_by[]=user_id` / `rbac_group_id` / `product` (`cowork`, `claude_code`
// and `chat_cowork_unified` only on this endpoint) each row is one (plugin, user),
// (plugin, group), or (plugin, product) cut: the flat `user_id` / `rbac_group_id`
// / `product` keys carry the cut and the counts are scoped to it.
type AnalyticsPluginActivity struct {
	// Claude Code activity metrics for a single plugin on a given day.
	ClaudeCodeMetrics AnalyticsPluginClaudeCodeMetrics `json:"claude_code_metrics" api:"required"`
	// Cowork activity metrics for a single plugin on a given day.
	CoworkMetrics AnalyticsPluginCoworkMetrics `json:"cowork_metrics" api:"required"`
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
	// Stable plugin identifier when available (e.g. `serena@claude-plugins-official`).
	// Null for third-party Claude Code plugins (redacted at the source) and Cowork
	// slash commands that carry only a hashed id.
	PluginID string `json:"plugin_id" api:"required"`
	// Name of the plugin
	PluginName string `json:"plugin_name" api:"required"`
	// Plugin use recorded while members had Chat and Cowork unified (Cowork's features
	// inside claude.ai chat) turned on. A count is null in date-range mode where it
	// cannot be computed. Omitted from the response where Chat and Cowork unified is
	// not available.
	ChatCoworkUnifiedMetrics AnalyticsPluginChatCoworkUnifiedMetrics `json:"chat_cowork_unified_metrics" api:"nullable"`
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
		PluginID                 respjson.Field
		PluginName               respjson.Field
		ChatCoworkUnifiedMetrics respjson.Field
		Product                  respjson.Field
		RBACGroupID              respjson.Field
		RBACGroupName            respjson.Field
		UserID                   respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AnalyticsPluginActivity) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsPluginActivity) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A plugin's use in Cowork sessions recorded while members had Chat and Cowork
// unified turned on.
type AnalyticsPluginChatCoworkUnifiedMetrics struct {
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
func (r AnalyticsPluginChatCoworkUnifiedMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsPluginChatCoworkUnifiedMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude Code activity metrics for a single plugin on a given day.
type AnalyticsPluginClaudeCodeMetrics struct {
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
func (r AnalyticsPluginClaudeCodeMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsPluginClaudeCodeMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Cowork activity metrics for a single plugin on a given day.
type AnalyticsPluginCoworkMetrics struct {
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
func (r AnalyticsPluginCoworkMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsPluginCoworkMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A dimension to group plugin activity by.
type AnalyticsPluginGroupBy string

const (
	AnalyticsPluginGroupByProduct     AnalyticsPluginGroupBy = "product"
	AnalyticsPluginGroupByRBACGroupID AnalyticsPluginGroupBy = "rbac_group_id"
	AnalyticsPluginGroupByUserID      AnalyticsPluginGroupBy = "user_id"
)

// Per-project activity data for a given day.
type AnalyticsProjectActivity struct {
	// Project creation timestamp in RFC 3339 format. Null if the project was deleted
	// before attribution was recorded.
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// User who created the project. Null if the project was deleted before attribution
	// was recorded, or if the creator's account no longer exists.
	CreatedBy AnalyticsUser `json:"created_by" api:"required"`
	// Number of distinct conversations in the project. Null on aggregated rows where a
	// distinct count cannot be computed.
	DistinctConversationCount int64 `json:"distinct_conversation_count" api:"required"`
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
		CreatedAt                 respjson.Field
		CreatedBy                 respjson.Field
		DistinctConversationCount respjson.Field
		DistinctUserCount         respjson.Field
		MessageCount              respjson.Field
		ProjectID                 respjson.Field
		ProjectName               respjson.Field
		RBACGroupID               respjson.Field
		RBACGroupName             respjson.Field
		UserID                    respjson.Field
		ExtraFields               map[string]respjson.Field
		raw                       string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AnalyticsProjectActivity) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsProjectActivity) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A dimension to group chat project activity by.
type AnalyticsProjectGroupBy string

const (
	AnalyticsProjectGroupByRBACGroupID AnalyticsProjectGroupBy = "rbac_group_id"
	AnalyticsProjectGroupByUserID      AnalyticsProjectGroupBy = "user_id"
)

// Claude Science activity metrics for a single user on a given day.
type AnalyticsScienceMetrics struct {
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
func (r AnalyticsScienceMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsScienceMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Per-day entry in the /summaries response.
type AnalyticsSingleDayActivitySummary struct {
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
	// Omitted from the response where Chat and Cowork unified is not available.
	ChatCoworkUnifiedDailyActiveUserCount int64 `json:"chat_cowork_unified_daily_active_user_count" api:"nullable"`
	// Number of users with activity in Chat and Cowork unified in the 28-day rolling
	// window (30 days when the request filters by `rbac_group_id`). Omitted from the
	// response where Chat and Cowork unified is not available.
	ChatCoworkUnifiedMonthlyActiveUserCount int64 `json:"chat_cowork_unified_monthly_active_user_count" api:"nullable"`
	// Number of users with activity in Chat and Cowork unified in the 7-day rolling
	// window. Omitted from the response where Chat and Cowork unified is not
	// available.
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
	// Number of Slack users active in Claude Tag on the requested day. Counts distinct
	// Slack users, including users who have not linked a Claude account, so it can
	// exceed `daily_active_user_count` and `assigned_seat_count`. Null on a day with
	// no measured value, when Claude Tag data cannot be read for the request, and when
	// the response is scoped to an RBAC group, because an unlinked user belongs to no
	// group. Omitted from the response on deployments that do not offer Claude Tag
	// reporting.
	ClaudeTagDailyActiveUserCount int64 `json:"claude_tag_daily_active_user_count" api:"nullable"`
	// Number of Slack users active in Claude Tag in the 28-day rolling window. Counts
	// distinct Slack users, including users who have not linked a Claude account, so
	// it can exceed `daily_active_user_count` and `assigned_seat_count`. Null on a day
	// with no measured value, when Claude Tag data cannot be read for the request, and
	// when the response is scoped to an RBAC group, because an unlinked user belongs
	// to no group. Omitted from the response on deployments that do not offer Claude
	// Tag reporting.
	ClaudeTagMonthlyActiveUserCount int64 `json:"claude_tag_monthly_active_user_count" api:"nullable"`
	// Number of Slack users active in Claude Tag in the 7-day rolling window. Counts
	// distinct Slack users, including users who have not linked a Claude account, so
	// it can exceed `daily_active_user_count` and `assigned_seat_count`. Null on a day
	// with no measured value, when Claude Tag data cannot be read for the request, and
	// when the response is scoped to an RBAC group, because an unlinked user belongs
	// to no group. Omitted from the response on deployments that do not offer Claude
	// Tag reporting.
	ClaudeTagWeeklyActiveUserCount int64 `json:"claude_tag_weekly_active_user_count" api:"nullable"`
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
		ClaudeTagDailyActiveUserCount           respjson.Field
		ClaudeTagMonthlyActiveUserCount         respjson.Field
		ClaudeTagWeeklyActiveUserCount          respjson.Field
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
func (r AnalyticsSingleDayActivitySummary) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsSingleDayActivitySummary) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Per-skill activity data for a given day.
type AnalyticsSkillActivity struct {
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
	AttributedListPrice string `json:"attributed_list_price" api:"required"`
	// Claude.ai activity metrics for a single skill on a given day.
	ChatMetrics AnalyticsSkillChatMetrics `json:"chat_metrics" api:"required"`
	// Claude Code activity metrics for a single skill on a given day.
	ClaudeCodeMetrics AnalyticsSkillClaudeCodeMetrics `json:"claude_code_metrics" api:"required"`
	// Cowork activity metrics for a single skill on a given day.
	CoworkMetrics AnalyticsSkillCoworkMetrics `json:"cowork_metrics" api:"required"`
	// Currency for this row's monetary fields (`estimated_overage_spend` and
	// `attributed_list_price`), as an uppercase ISO-4217 code. Always "USD" when
	// either amount is populated; null whenever both amounts are null.
	Currency string `json:"currency" api:"required"`
	// Number of distinct users who used the skill on the requested day, or, in
	// date-range mode, over the requested window — recomputed as an exact distinct
	// count over the window's per-member daily rows, never a sum of per-day values. A
	// skill counts as used only when it is explicitly activated — the model (or the
	// user, via the skill's slash command) invokes it, reading its instructions into
	// context as part of that activation. Skills that are merely installed or listed
	// as available, or whose content reaches the context without an activation
	// (preloaded, hook-injected, or read as a plain file), are not counted.
	DistinctUserCount int64 `json:"distinct_user_count" api:"required"`
	// Distinct accounts that enabled this skill on the requested day (claude.ai only —
	// the skill analog of plugin `install_count`). The count is org-wide: null when
	// enable reporting is not enabled for this organization, or when the request
	// scopes to `user_id` / `rbac_group_id` / `product` via `group_by[]` or `filter[]`
	// (an org-wide count would be misleading on per-cut rows). A distinct count, not
	// an event count: summing across days double-counts members who enable the skill
	// on more than one day, so it is also null in date-range rollup mode
	// (`starting_date`/`ending_date`).
	EnableCount int64 `json:"enable_count" api:"required"`
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
	EstimatedOverageSpend string `json:"estimated_overage_spend" api:"required"`
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
	InvocationCount int64 `json:"invocation_count" api:"required"`
	// Office Agent activity metrics for a single skill on a given day, broken out by
	// Office product.
	OfficeMetrics AnalyticsSkillOfficeMetrics `json:"office_metrics" api:"required"`
	// Skill share status (claude.ai only): one of `private`, `organization`, or
	// `public`. Null for skills used only in Claude Code or Office (no per-skill
	// share-status concept) and when share-status reporting is not yet available for
	// the organization. Filterable via `filter[]=share_status:{value}`.
	//
	// Any of "organization", "private", "public".
	ShareStatus AnalyticsSkillActivityShareStatus `json:"share_status" api:"required"`
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
	SkillDisplayName string `json:"skill_display_name" api:"required"`
	// Name of the skill
	SkillName string `json:"skill_name" api:"required"`
	// Skill use recorded while members had Chat and Cowork unified (Cowork's features
	// inside claude.ai chat) turned on, split into chat conversations and Cowork
	// sessions. A count is null in date-range mode where it cannot be computed.
	// Omitted from the response where Chat and Cowork unified is not available.
	ChatCoworkUnifiedMetrics AnalyticsSkillChatCoworkUnifiedMetrics `json:"chat_cowork_unified_metrics" api:"nullable"`
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
		AttributedListPrice      respjson.Field
		ChatMetrics              respjson.Field
		ClaudeCodeMetrics        respjson.Field
		CoworkMetrics            respjson.Field
		Currency                 respjson.Field
		DistinctUserCount        respjson.Field
		EnableCount              respjson.Field
		EstimatedOverageSpend    respjson.Field
		InvocationCount          respjson.Field
		OfficeMetrics            respjson.Field
		ShareStatus              respjson.Field
		SkillDisplayName         respjson.Field
		SkillName                respjson.Field
		ChatCoworkUnifiedMetrics respjson.Field
		Product                  respjson.Field
		RBACGroupID              respjson.Field
		RBACGroupName            respjson.Field
		UserID                   respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AnalyticsSkillActivity) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsSkillActivity) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Skill share status (claude.ai only): one of `private`, `organization`, or
// `public`. Null for skills used only in Claude Code or Office (no per-skill
// share-status concept) and when share-status reporting is not yet available for
// the organization. Filterable via `filter[]=share_status:{value}`.
type AnalyticsSkillActivityShareStatus string

const (
	AnalyticsSkillActivityShareStatusOrganization AnalyticsSkillActivityShareStatus = "organization"
	AnalyticsSkillActivityShareStatusPrivate      AnalyticsSkillActivityShareStatus = "private"
	AnalyticsSkillActivityShareStatusPublic       AnalyticsSkillActivityShareStatus = "public"
)

// A skill's use in chat conversations recorded while members had Chat and Cowork
// unified turned on.
type AnalyticsSkillChatCoworkUnifiedChatMetrics struct {
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
func (r AnalyticsSkillChatCoworkUnifiedChatMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsSkillChatCoworkUnifiedChatMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A skill's use in Chat and Cowork unified on one day: `chat` holds its use in
// chat conversations and `sessions` its use in Cowork sessions.
type AnalyticsSkillChatCoworkUnifiedMetrics struct {
	// A skill's use in chat conversations recorded while members had Chat and Cowork
	// unified turned on.
	Chat AnalyticsSkillChatCoworkUnifiedChatMetrics `json:"chat" api:"required"`
	// A skill's use in Cowork sessions recorded while members had Chat and Cowork
	// unified turned on.
	Sessions AnalyticsSkillChatCoworkUnifiedSessionsMetrics `json:"sessions" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Chat        respjson.Field
		Sessions    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AnalyticsSkillChatCoworkUnifiedMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsSkillChatCoworkUnifiedMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A skill's use in Cowork sessions recorded while members had Chat and Cowork
// unified turned on.
type AnalyticsSkillChatCoworkUnifiedSessionsMetrics struct {
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
func (r AnalyticsSkillChatCoworkUnifiedSessionsMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsSkillChatCoworkUnifiedSessionsMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude.ai activity metrics for a single skill on a given day.
type AnalyticsSkillChatMetrics struct {
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
func (r AnalyticsSkillChatMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsSkillChatMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Claude Code activity metrics for a single skill on a given day.
type AnalyticsSkillClaudeCodeMetrics struct {
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
func (r AnalyticsSkillClaudeCodeMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsSkillClaudeCodeMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Cowork activity metrics for a single skill on a given day.
type AnalyticsSkillCoworkMetrics struct {
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
func (r AnalyticsSkillCoworkMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsSkillCoworkMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A dimension to group skill activity by.
type AnalyticsSkillGroupBy string

const (
	AnalyticsSkillGroupByProduct     AnalyticsSkillGroupBy = "product"
	AnalyticsSkillGroupByRBACGroupID AnalyticsSkillGroupBy = "rbac_group_id"
	AnalyticsSkillGroupByUserID      AnalyticsSkillGroupBy = "user_id"
)

// Office Agent activity metrics for a single skill on a given day, broken out by
// Office product.
type AnalyticsSkillOfficeMetrics struct {
	// Office Agent activity metrics for a single skill on a given day within one
	// Office product.
	Excel AnalyticsSkillOfficeProductMetrics `json:"excel" api:"required"`
	// Office Agent activity metrics for a single skill on a given day within one
	// Office product.
	Outlook AnalyticsSkillOfficeProductMetrics `json:"outlook" api:"required"`
	// Office Agent activity metrics for a single skill on a given day within one
	// Office product.
	Powerpoint AnalyticsSkillOfficeProductMetrics `json:"powerpoint" api:"required"`
	// Office Agent activity metrics for a single skill on a given day within one
	// Office product.
	Word AnalyticsSkillOfficeProductMetrics `json:"word" api:"required"`
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
func (r AnalyticsSkillOfficeMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsSkillOfficeMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Office Agent activity metrics for a single skill on a given day within one
// Office product.
type AnalyticsSkillOfficeProductMetrics struct {
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
func (r AnalyticsSkillOfficeProductMetrics) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsSkillOfficeProductMetrics) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Accepted/rejected counts for a single Claude Code tool type.
type AnalyticsToolActionCounts struct {
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
func (r AnalyticsToolActionCounts) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsToolActionCounts) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Per-tool accepted/rejected counts for Claude Code file modification tools.
type AnalyticsToolActions struct {
	// Accepted/rejected counts for a single Claude Code tool type.
	EditTool AnalyticsToolActionCounts `json:"edit_tool" api:"required"`
	// Accepted/rejected counts for a single Claude Code tool type.
	MultiEditTool AnalyticsToolActionCounts `json:"multi_edit_tool" api:"required"`
	// Accepted/rejected counts for a single Claude Code tool type.
	NotebookEditTool AnalyticsToolActionCounts `json:"notebook_edit_tool" api:"required"`
	// Accepted/rejected counts for a single Claude Code tool type.
	WriteTool AnalyticsToolActionCounts `json:"write_tool" api:"required"`
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
func (r AnalyticsToolActions) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsToolActions) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A user in the organization, identified by tagged id and email address.
type AnalyticsUser struct {
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
func (r AnalyticsUser) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsUser) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Per-user activity data for a given day.
type AnalyticsUserActivity struct {
	// Claude.ai activity metrics for a single user on a given day.
	ChatMetrics AnalyticsChatMetrics `json:"chat_metrics" api:"required"`
	// Claude Code activity metrics for a single user on a given day.
	ClaudeCodeMetrics AnalyticsClaudeCodeMetrics `json:"claude_code_metrics" api:"required"`
	// Cowork activity metrics for a single user on a given day.
	CoworkMetrics AnalyticsCoworkMetrics `json:"cowork_metrics" api:"required"`
	// Claude Design activity metrics for a single user on a given day.
	DesignMetrics AnalyticsDesignMetrics `json:"design_metrics" api:"required"`
	// Number of distinct active users represented by this row. Only set for grouped
	// rollups (`group_by[]`); null for per-user rows. In date-range mode, recomputed
	// as an exact distinct count of the group's active members over the requested
	// window, never a sum of per-day values.
	DistinctUserCount int64 `json:"distinct_user_count" api:"required"`
	// Office Agent activity metrics for a single user on a given day, broken out by
	// Office product.
	OfficeMetrics AnalyticsOfficeMetrics `json:"office_metrics" api:"required"`
	// Claude Science activity metrics for a single user on a given day.
	ScienceMetrics AnalyticsScienceMetrics `json:"science_metrics" api:"required"`
	// The user this row describes. Null on rows aggregated across users.
	User AnalyticsUser `json:"user" api:"required"`
	// Number of web searches performed
	WebSearchCount int64 `json:"web_search_count" api:"required"`
	// Activity recorded while the member had Chat and Cowork unified (Cowork's
	// features inside claude.ai chat) turned on, split into `chat` (chat activity) and
	// `sessions` (Cowork activity). Omitted from the response where Chat and Cowork
	// unified is not available.
	ChatCoworkUnifiedMetrics AnalyticsChatCoworkUnifiedMetrics `json:"chat_cowork_unified_metrics" api:"nullable"`
	// Claude Tag activity in Slack by a user who has linked a Claude account. Activity
	// by Slack users who have not linked an account appears only in the `/summaries`
	// Claude Tag counts, so these rows do not add up to those counts. Omitted from the
	// response on deployments that do not offer Claude Tag reporting.
	ClaudeTagMetrics AnalyticsClaudeTagMetrics `json:"claude_tag_metrics" api:"nullable"`
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
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ChatMetrics              respjson.Field
		ClaudeCodeMetrics        respjson.Field
		CoworkMetrics            respjson.Field
		DesignMetrics            respjson.Field
		DistinctUserCount        respjson.Field
		OfficeMetrics            respjson.Field
		ScienceMetrics           respjson.Field
		User                     respjson.Field
		WebSearchCount           respjson.Field
		ChatCoworkUnifiedMetrics respjson.Field
		ClaudeTagMetrics         respjson.Field
		LastActivityDate         respjson.Field
		RBACGroupID              respjson.Field
		RBACGroupName            respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r AnalyticsUserActivity) RawJSON() string { return r.JSON.raw }
func (r *AnalyticsUserActivity) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A dimension to group user activity by.
type AnalyticsUserGroupBy string

const (
	AnalyticsUserGroupByRBACGroupID AnalyticsUserGroupBy = "rbac_group_id"
)
