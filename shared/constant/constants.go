package constant

import (
	shimjson "github.com/anthropics/anthropic-sdk-go/v2/internal/encoding/json"
	"github.com/anthropics/anthropic-sdk-go/v2/internal/paramutil"
)

// ModelNonStreamingTokens defines the maximum tokens for models that should limit
// non-streaming requests.
var ModelNonStreamingTokens = map[string]int{
	"claude-opus-4-20250514":                  8192,
	"claude-4-opus-20250514":                  8192,
	"claude-opus-4-0":                         8192,
	"anthropic.claude-opus-4-20250514-v1:0":   8192,
	"claude-opus-4@20250514":                  8192,
	"anthropic.claude-opus-4-1-20250805-v1:0": 8192,
	"claude-opus-4-1@20250805":                8192,
}

type Constant[T any] interface {
	Default() T
}

// ValueOf gives the default value of a constant from its type. It's helpful when
// constructing constants as variants in a one-of. Note that empty structs are
// marshalled by default. Usage: constant.ValueOf[constant.Foo]()
func ValueOf[T Constant[T]]() T {
	var t T
	return t.Default()
}

type Adaptive string                                // Always "adaptive"
type Advisor string                                 // Always "advisor"
type Advisor20260301 string                         // Always "advisor_20260301"
type AdvisorMessage string                          // Always "advisor_message"
type AdvisorRedactedResult string                   // Always "advisor_redacted_result"
type AdvisorResult string                           // Always "advisor_result"
type AdvisorToolResult string                       // Always "advisor_tool_result"
type AdvisorToolResultError string                  // Always "advisor_tool_result_error"
type Agent string                                   // Always "agent"
type AgentArchivedError string                      // Always "agent_archived_error"
type AgentToolset20260401 string                    // Always "agent_toolset_20260401"
type AgentWithOverrides string                      // Always "agent_with_overrides"
type AgentArchived string                           // Always "agent.archived"
type AgentCreated string                            // Always "agent.created"
type AgentCustomToolUse string                      // Always "agent.custom_tool_use"
type AgentDeleted string                            // Always "agent.deleted"
type AgentMCPToolResult string                      // Always "agent.mcp_tool_result"
type AgentMCPToolUse string                         // Always "agent.mcp_tool_use"
type AgentMessage string                            // Always "agent.message"
type AgentThinking string                           // Always "agent.thinking"
type AgentThreadContextCompacted string             // Always "agent.thread_context_compacted"
type AgentThreadMessageReceived string              // Always "agent.thread_message_received"
type AgentThreadMessageSent string                  // Always "agent.thread_message_sent"
type AgentToolResult string                         // Always "agent.tool_result"
type AgentToolUse string                            // Always "agent.tool_use"
type AgentUpdated string                            // Always "agent.updated"
type All string                                     // Always "all"
type AllConnectors string                           // Always "all_connectors"
type Allow string                                   // Always "allow"
type AlwaysAllow string                             // Always "always_allow"
type AlwaysAsk string                               // Always "always_ask"
type Anthropic string                               // Always "anthropic"
type Any string                                     // Always "any"
type APIActor string                                // Always "api_actor"
type APIError string                                // Always "api_error"
type APIKey string                                  // Always "api_key"
type ApplicationPDF string                          // Always "application/pdf"
type Approximate string                             // Always "approximate"
type Ask string                                     // Always "ask"
type Assistant string                               // Always "assistant"
type Attached string                                // Always "attached"
type AuthenticationError string                     // Always "authentication_error"
type Auto string                                    // Always "auto"
type AWS string                                     // Always "aws"
type AWSRole string                                 // Always "aws_role"
type Azure string                                   // Always "azure"
type Ban string                                     // Always "ban"
type Base64 string                                  // Always "base64"
type Bash string                                    // Always "bash"
type Bash20241022 string                            // Always "bash_20241022"
type Bash20250124 string                            // Always "bash_20250124"
type BashCodeExecutionOutput string                 // Always "bash_code_execution_output"
type BashCodeExecutionResult string                 // Always "bash_code_execution_result"
type BashCodeExecutionToolResult string             // Always "bash_code_execution_tool_result"
type BashCodeExecutionToolResultError string        // Always "bash_code_execution_tool_result_error"
type Batch string                                   // Always "batch"
type BetweenTools string                            // Always "between_tools"
type BillingError string                            // Always "billing_error"
type Branch string                                  // Always "branch"
type Browser string                                 // Always "browser"
type BrowserState string                            // Always "browser_state"
type BrowserToolset20260801 string                  // Always "browser_toolset_20260801"
type BudgetReached string                           // Always "budget_reached"
type Caller string                                  // Always "caller"
type Canceled string                                // Always "canceled"
type CharLocation string                            // Always "char_location"
type CitationsDelta string                          // Always "citations_delta"
type ClearThinking20251015 string                   // Always "clear_thinking_20251015"
type ClearToolUses20250919 string                   // Always "clear_tool_uses_20250919"
type ClientSecretBasic string                       // Always "client_secret_basic"
type ClientSecretPost string                        // Always "client_secret_post"
type CloseTab string                                // Always "close_tab"
type Cloud string                                   // Always "cloud"
type CodeExecution string                           // Always "code_execution"
type CodeExecution20250522 string                   // Always "code_execution_20250522"
type CodeExecution20250825 string                   // Always "code_execution_20250825"
type CodeExecution20260120 string                   // Always "code_execution_20260120"
type CodeExecution20260521 string                   // Always "code_execution_20260521"
type CodeExecutionOutput string                     // Always "code_execution_output"
type CodeExecutionResult string                     // Always "code_execution_result"
type CodeExecutionToolResult string                 // Always "code_execution_tool_result"
type CodeExecutionToolResultError string            // Always "code_execution_tool_result_error"
type Commit string                                  // Always "commit"
type Compact20260112 string                         // Always "compact_20260112"
type Compaction string                              // Always "compaction"
type CompactionDelta string                         // Always "compaction_delta"
type Completed string                               // Always "completed"
type ComplianceSettings string                      // Always "compliance_settings"
type Computer string                                // Always "computer"
type Computer20241022 string                        // Always "computer_20241022"
type Computer20250124 string                        // Always "computer_20250124"
type Computer20251124 string                        // Always "computer_20251124"
type ComputerToolset20260801 string                 // Always "computer_toolset_20260801"
type ConflictError string                           // Always "conflict_error"
type Connector string                               // Always "connector"
type ConnectorScope string                          // Always "connector_scope"
type ConnectorTool string                           // Always "connector_tool"
type ContainerUpload string                         // Always "container_upload"
type Content string                                 // Always "content"
type ContentBlockDelta string                       // Always "content_block_delta"
type ContentBlockLocation string                    // Always "content_block_location"
type ContentBlockStart string                       // Always "content_block_start"
type ContentBlockStop string                        // Always "content_block_stop"
type ContentDelta string                            // Always "content_delta"
type ContentSha256 string                           // Always "content_sha256"
type Continue string                                // Always "continue"
type Coordinate string                              // Always "coordinate"
type Coordinator string                             // Always "coordinator"
type Create string                                  // Always "create"
type CreateNew string                               // Always "create_new"
type CredentialHostUnreachableError string          // Always "credential_host_unreachable_error"
type Cron string                                    // Always "cron"
type CursorPosition string                          // Always "cursor_position"
type Custom string                                  // Always "custom"
type Default string                                 // Always "default"
type Delete string                                  // Always "delete"
type Deny string                                    // Always "deny"
type Deployment string                              // Always "deployment"
type DeploymentRun string                           // Always "deployment_run"
type DeploymentRunFailed string                     // Always "deployment_run.failed"
type DeploymentRunStarted string                    // Always "deployment_run.started"
type DeploymentRunSucceeded string                  // Always "deployment_run.succeeded"
type DeploymentArchived string                      // Always "deployment.archived"
type DeploymentCreated string                       // Always "deployment.created"
type DeploymentDeleted string                       // Always "deployment.deleted"
type DeploymentPaused string                        // Always "deployment.paused"
type DeploymentUnpaused string                      // Always "deployment.unpaused"
type DeploymentUpdated string                       // Always "deployment.updated"
type Direct string                                  // Always "direct"
type Disabled string                                // Always "disabled"
type Discovery string                               // Always "discovery"
type Document string                                // Always "document"
type DoubleClick string                             // Always "double_click"
type DownloadCompleted string                       // Always "download_completed"
type DownloadFailed string                          // Always "download_failed"
type DownloadStarted string                         // Always "download_started"
type Dream string                                   // Always "dream"
type Edit string                                    // Always "edit"
type Enabled string                                 // Always "enabled"
type EncryptedCodeExecutionResult string            // Always "encrypted_code_execution_result"
type EndTurn string                                 // Always "end_turn"
type EnrollEndUser string                           // Always "enroll_end_user"
type EnrollmentURL string                           // Always "enrollment_url"
type Environment string                             // Always "environment"
type EnvironmentArchivedError string                // Always "environment_archived_error"
type EnvironmentDeletedResponseType string          // Always "environment_deleted"
type EnvironmentNotFoundError string                // Always "environment_not_found_error"
type EnvironmentVariable string                     // Always "environment_variable"
type EnvironmentArchived string                     // Always "environment.archived"
type EnvironmentCreated string                      // Always "environment.created"
type EnvironmentDeleted string                      // Always "environment.deleted"
type EnvironmentUpdated string                      // Always "environment.updated"
type Ephemeral string                               // Always "ephemeral"
type Error string                                   // Always "error"
type Errored string                                 // Always "errored"
type Event string                                   // Always "event"
type EventDelta string                              // Always "event_delta"
type EventStart string                              // Always "event_start"
type Except string                                  // Always "except"
type Exclude string                                 // Always "exclude"
type Exhausted string                               // Always "exhausted"
type Expired string                                 // Always "expired"
type ExplicitURL string                             // Always "explicit_url"
type ExternalKey string                             // Always "external_key"
type ExternalKeyDeleted string                      // Always "external_key_deleted"
type ExternalKeyValidation string                   // Always "external_key_validation"
type Failure string                                 // Always "failure"
type Fallback string                                // Always "fallback"
type FallbackMessage string                         // Always "fallback_message"
type FederatedActor string                          // Always "federated_actor"
type FederationIssuer string                        // Always "federation_issuer"
type FederationRule string                          // Always "federation_rule"
type FederationRuleWorkspace string                 // Always "federation_rule_workspace"
type FederationRuleWorkspaceDeleted string          // Always "federation_rule_workspace_deleted"
type File string                                    // Always "file"
type FileNotFoundError string                       // Always "file_not_found_error"
type FileUpload string                              // Always "file_upload"
type Files string                                   // Always "files"
type Find string                                    // Always "find"
type FormInput string                               // Always "form_input"
type GCP string                                     // Always "gcp"
type GetPageText string                             // Always "get_page_text"
type GitHubRepository string                        // Always "github_repository"
type Glob string                                    // Always "glob"
type Grep string                                    // Always "grep"
type High string                                    // Always "high"
type HoldKey string                                 // Always "hold_key"
type Hover string                                   // Always "hover"
type Image string                                   // Always "image"
type Include string                                 // Always "include"
type Inline string                                  // Always "inline"
type InputJSONDelta string                          // Always "input_json_delta"
type InputTokens string                             // Always "input_tokens"
type Insert string                                  // Always "insert"
type InvalidAddress string                          // Always "invalid_address"
type InvalidRequestError string                     // Always "invalid_request_error"
type Invite string                                  // Always "invite"
type InviteDeleted string                           // Always "invite_deleted"
type JavascriptExec string                          // Always "javascript_exec"
type JSONSchema string                              // Always "json_schema"
type Key string                                     // Always "key"
type LeftClick string                               // Always "left_click"
type LeftClickDrag string                           // Always "left_click_drag"
type LeftMouseDown string                           // Always "left_mouse_down"
type LeftMouseUp string                             // Always "left_mouse_up"
type Limit string                                   // Always "limit"
type Limited string                                 // Always "limited"
type ListTabs string                                // Always "list_tabs"
type Low string                                     // Always "low"
type Manual string                                  // Always "manual"
type Max string                                     // Always "max"
type MaxWorkflowRunsError string                    // Always "max_workflow_runs_error"
type MCPAuthenticationFailedError string            // Always "mcp_authentication_failed_error"
type MCPConnectionFailedError string                // Always "mcp_connection_failed_error"
type MCPEgressBlockedError string                   // Always "mcp_egress_blocked_error"
type MCPOAuth string                                // Always "mcp_oauth"
type MCPToolListing string                          // Always "mcp_tool_listing"
type MCPToolReference string                        // Always "mcp_tool_reference"
type MCPToolResult string                           // Always "mcp_tool_result"
type MCPToolUse string                              // Always "mcp_tool_use"
type MCPToolset string                              // Always "mcp_toolset"
type MCPToolsetReference string                     // Always "mcp_toolset_reference"
type Medium string                                  // Always "medium"
type Memory string                                  // Always "memory"
type Memory20250818 string                          // Always "memory_20250818"
type MemoryDeleted string                           // Always "memory_deleted"
type MemoryPathConflictError string                 // Always "memory_path_conflict_error"
type MemoryPreconditionFailedError string           // Always "memory_precondition_failed_error"
type MemoryPrefix string                            // Always "memory_prefix"
type MemoryStore string                             // Always "memory_store"
type MemoryStoreArchivedError string                // Always "memory_store_archived_error"
type MemoryStoreDeletedResponseType string          // Always "memory_store_deleted"
type MemoryStoreArchived string                     // Always "memory_store.archived"
type MemoryStoreCreated string                      // Always "memory_store.created"
type MemoryStoreDeleted string                      // Always "memory_store.deleted"
type MemoryVersion string                           // Always "memory_version"
type Message string                                 // Always "message"
type MessageBatch string                            // Always "message_batch"
type MessageBatchDeleted string                     // Always "message_batch_deleted"
type MessageDelta string                            // Always "message_delta"
type MessageStart string                            // Always "message_start"
type MessageStop string                             // Always "message_stop"
type MessagesChanged string                         // Always "messages_changed"
type MiddleClick string                             // Always "middle_click"
type Model string                                   // Always "model"
type ModelChanged string                            // Always "model_changed"
type ModelGroup string                              // Always "model_group"
type ModelOverloadedError string                    // Always "model_overloaded_error"
type ModelRateLimitedError string                   // Always "model_rate_limited_error"
type ModelRequestFailedError string                 // Always "model_request_failed_error"
type MouseMove string                               // Always "mouse_move"
type Multiagent20261001 string                      // Always "multiagent_20261001"
type Navigate string                                // Always "navigate"
type NewTab string                                  // Always "new_tab"
type None string                                    // Always "none"
type NotApplied string                              // Always "not_applied"
type NotFoundError string                           // Always "not_found_error"
type Object string                                  // Always "object"
type OIDC string                                    // Always "oidc"
type Only string                                    // Always "only"
type Organization string                            // Always "organization"
type OrganizationDisabledError string               // Always "organization_disabled_error"
type OrganizationMember string                      // Always "organization_member"
type OrganizationService string                     // Always "organization_service"
type OutcomeEvaluation string                       // Always "outcome_evaluation"
type OverloadedError string                         // Always "overloaded_error"
type PageLocation string                            // Always "page_location"
type PermissionError string                         // Always "permission_error"
type Plugin string                                  // Always "plugin"
type PluginDeleted string                           // Always "plugin_deleted"
type PluginInstallationSetting string               // Always "plugin_installation_setting"
type PluginInstallationSettingDeleted string        // Always "plugin_installation_setting_deleted"
type PluginMarketplace string                       // Always "plugin_marketplace"
type PluginMarketplaceValidationReport string       // Always "plugin_marketplace_validation_report"
type PluginShare string                             // Always "plugin_share"
type PluginVersion string                           // Always "plugin_version"
type PreviousMessageNotFound string                 // Always "previous_message_not_found"
type ProgramError string                            // Always "program_error"
type RateLimit string                               // Always "rate_limit"
type RateLimitError string                          // Always "rate_limit_error"
type RBACGroup string                               // Always "rbac_group"
type RBACGroupDeleted string                        // Always "rbac_group_deleted"
type RBACGroupMember string                         // Always "rbac_group_member"
type RBACGroupMemberDeleted string                  // Always "rbac_group_member_deleted"
type RBACRole string                                // Always "rbac_role"
type RBACRolePermission string                      // Always "rbac_role_permission"
type Read string                                    // Always "read"
type ReadConsole string                             // Always "read_console"
type ReadNetwork string                             // Always "read_network"
type ReadPage string                                // Always "read_page"
type Redacted string                                // Always "redacted"
type RedactedThinking string                        // Always "redacted_thinking"
type Redeemed string                                // Always "redeemed"
type Redirect string                                // Always "redirect"
type Ref string                                     // Always "ref"
type Refusal string                                 // Always "refusal"
type Rename string                                  // Always "rename"
type RepositoryAuthenticationError string           // Always "repository_authentication_error"
type RepositoryCheckoutError string                 // Always "repository_checkout_error"
type RepositoryCloneError string                    // Always "repository_clone_error"
type RepositoryForbiddenError string                // Always "repository_forbidden_error"
type RepositoryNotFoundError string                 // Always "repository_not_found_error"
type RequiresAction string                          // Always "requires_action"
type Restriction string                             // Always "restriction"
type RetriesExhausted string                        // Always "retries_exhausted"
type Retrying string                                // Always "retrying"
type RightClick string                              // Always "right_click"
type Schedule string                                // Always "schedule"
type ScopedAPIKeyActor string                       // Always "scoped_api_key_actor"
type Screenshot string                              // Always "screenshot"
type Scroll string                                  // Always "scroll"
type ScrollTo string                                // Always "scroll_to"
type SearchResult string                            // Always "search_result"
type SearchResultLocation string                    // Always "search_result_location"
type SeatTier string                                // Always "seat_tier"
type Self string                                    // Always "self"
type SelfHosted string                              // Always "self_hosted"
type SelfHostedResourcesUnsupportedError string     // Always "self_hosted_resources_unsupported_error"
type ServerToolUse string                           // Always "server_tool_use"
type ServiceAccount string                          // Always "service_account"
type ServiceAccountActor string                     // Always "service_account_actor"
type ServiceAccountWorkspaceMember string           // Always "service_account_workspace_member"
type ServiceAccountWorkspaceMemberDeleted string    // Always "service_account_workspace_member_deleted"
type Session string                                 // Always "session"
type SessionActor string                            // Always "session_actor"
type SessionCreationRejectedError string            // Always "session_creation_rejected_error"
type SessionDeletedResponseType string              // Always "session_deleted"
type SessionRateLimitedError string                 // Always "session_rate_limited_error"
type SessionResourceDeleted string                  // Always "session_resource_deleted"
type SessionResourceNotFoundError string            // Always "session_resource_not_found_error"
type SessionThread string                           // Always "session_thread"
type SessionArchived string                         // Always "session.archived"
type SessionBudgetReached string                    // Always "session.budget_reached"
type SessionCreated string                          // Always "session.created"
type SessionDeleted string                          // Always "session.deleted"
type SessionError string                            // Always "session.error"
type SessionIdled string                            // Always "session.idled"
type SessionOutcomeEvaluationEnded string           // Always "session.outcome_evaluation_ended"
type SessionPending string                          // Always "session.pending"
type SessionRequiresAction string                   // Always "session.requires_action"
type SessionRunning string                          // Always "session.running"
type SessionStatusIdle string                       // Always "session.status_idle"
type SessionStatusIdled string                      // Always "session.status_idled"
type SessionStatusRescheduled string                // Always "session.status_rescheduled"
type SessionStatusRunStarted string                 // Always "session.status_run_started"
type SessionStatusRunning string                    // Always "session.status_running"
type SessionStatusTerminated string                 // Always "session.status_terminated"
type SessionThreadCreated string                    // Always "session.thread_created"
type SessionThreadIdled string                      // Always "session.thread_idled"
type SessionThreadStatusIdle string                 // Always "session.thread_status_idle"
type SessionThreadStatusRescheduled string          // Always "session.thread_status_rescheduled"
type SessionThreadStatusRunning string              // Always "session.thread_status_running"
type SessionThreadStatusTerminated string           // Always "session.thread_status_terminated"
type SessionThreadTerminated string                 // Always "session.thread_terminated"
type SessionUpdated string                          // Always "session.updated"
type SessionUsage string                            // Always "session.usage"
type Sessions string                                // Always "sessions"
type SignatureDelta string                          // Always "signature_delta"
type Skill string                                   // Always "skill"
type SkillDeleted string                            // Always "skill_deleted"
type SkillNotFoundError string                      // Always "skill_not_found_error"
type SkillVersion string                            // Always "skill_version"
type SkillVersionDeleted string                     // Always "skill_version_deleted"
type Skills string                                  // Always "skills"
type SpanModelRequestEnd string                     // Always "span.model_request_end"
type SpanModelRequestStart string                   // Always "span.model_request_start"
type SpanOutcomeEvaluationEnd string                // Always "span.outcome_evaluation_end"
type SpanOutcomeEvaluationOngoing string            // Always "span.outcome_evaluation_ongoing"
type SpanOutcomeEvaluationStart string              // Always "span.outcome_evaluation_start"
type SpendLimit string                              // Always "spend_limit"
type SpendLimitDeleted string                       // Always "spend_limit_deleted"
type SpendLimitIncreaseRequest string               // Always "spend_limit_increase_request"
type StaticBearer string                            // Always "static_bearer"
type Stopped string                                 // Always "stopped"
type StrReplace string                              // Always "str_replace"
type StrReplaceBasedEditTool string                 // Always "str_replace_based_edit_tool"
type StrReplaceEditor string                        // Always "str_replace_editor"
type Succeeded string                               // Always "succeeded"
type Success string                                 // Always "success"
type Summarize string                               // Always "summarize"
type SustainedFailure string                        // Always "sustained_failure"
type SwitchTab string                               // Always "switch_tab"
type System string                                  // Always "system"
type SystemChanged string                           // Always "system_changed"
type SystemMessage string                           // Always "system.message"
type TabOpened string                               // Always "tab_opened"
type TelemetryDestination string                    // Always "telemetry_destination"
type TelemetryDestinationDeleted string             // Always "telemetry_destination_deleted"
type TelemetryDestinationTest string                // Always "telemetry_destination_test"
type TelemetryDestinationWorkspace string           // Always "telemetry_destination_workspace"
type TelemetryDestinationWorkspaceDeleted string    // Always "telemetry_destination_workspace_deleted"
type Terminal string                                // Always "terminal"
type Text string                                    // Always "text"
type TextDelta string                               // Always "text_delta"
type TextEditor20241022 string                      // Always "text_editor_20241022"
type TextEditor20250124 string                      // Always "text_editor_20250124"
type TextEditor20250429 string                      // Always "text_editor_20250429"
type TextEditor20250728 string                      // Always "text_editor_20250728"
type TextEditorCodeExecutionCreateResult string     // Always "text_editor_code_execution_create_result"
type TextEditorCodeExecutionStrReplaceResult string // Always "text_editor_code_execution_str_replace_result"
type TextEditorCodeExecutionToolResult string       // Always "text_editor_code_execution_tool_result"
type TextEditorCodeExecutionToolResultError string  // Always "text_editor_code_execution_tool_result_error"
type TextEditorCodeExecutionViewResult string       // Always "text_editor_code_execution_view_result"
type TextPlain string                               // Always "text/plain"
type Thinking string                                // Always "thinking"
type ThinkingDelta string                           // Always "thinking_delta"
type ThinkingDropped string                         // Always "thinking_dropped"
type ThinkingMismatchAllowed string                 // Always "thinking_mismatch_allowed"
type ThinkingTurns string                           // Always "thinking_turns"
type ThreadLimitError string                        // Always "thread_limit_error"
type TimeoutError string                            // Always "timeout_error"
type TokenCount string                              // Always "token_count"
type Tokens string                                  // Always "tokens"
type Tool string                                    // Always "tool"
type ToolAddition string                            // Always "tool_addition"
type ToolDefinition string                          // Always "tool_definition"
type ToolReference string                           // Always "tool_reference"
type ToolRemoval string                             // Always "tool_removal"
type ToolResult string                              // Always "tool_result"
type ToolSearchToolBm25 string                      // Always "tool_search_tool_bm25"
type ToolSearchToolRegex string                     // Always "tool_search_tool_regex"
type ToolSearchToolResult string                    // Always "tool_search_tool_result"
type ToolSearchToolResultError string               // Always "tool_search_tool_result_error"
type ToolSearchToolSearchResult string              // Always "tool_search_tool_search_result"
type ToolUse string                                 // Always "tool_use"
type ToolUses string                                // Always "tool_uses"
type ToolsChanged string                            // Always "tools_changed"
type TripleClick string                             // Always "triple_click"
type Tunnel string                                  // Always "tunnel"
type TunnelCertificate string                       // Always "tunnel_certificate"
type TunnelToken string                             // Always "tunnel_token"
type Type string                                    // Always "type"
type Unattached string                              // Always "unattached"
type Unavailable string                             // Always "unavailable"
type UnknownError string                            // Always "unknown_error"
type Unrestricted string                            // Always "unrestricted"
type UpdateExisting string                          // Always "update_existing"
type URL string                                     // Always "url"
type UsageTier string                               // Always "usage_tier"
type User string                                    // Always "user"
type UserActor string                               // Always "user_actor"
type UserDeleted string                             // Always "user_deleted"
type UserProfile string                             // Always "user_profile"
type UserProfileActionURL string                    // Always "user_profile_action_url"
type UserProfileDeleted string                      // Always "user_profile_deleted"
type UserProfileEnforcementActionOpened string      // Always "user_profile.enforcement_action_opened"
type UserProfileEnforcementActionResolved string    // Always "user_profile.enforcement_action_resolved"
type UserCustomToolResult string                    // Always "user.custom_tool_result"
type UserDefineOutcome string                       // Always "user.define_outcome"
type UserInterrupt string                           // Always "user.interrupt"
type UserMessage string                             // Always "user.message"
type UserToolConfirmation string                    // Always "user.tool_confirmation"
type UserToolResult string                          // Always "user.tool_result"
type Vault string                                   // Always "vault"
type VaultArchivedError string                      // Always "vault_archived_error"
type VaultCredential string                         // Always "vault_credential"
type VaultCredentialDeletedResponseType string      // Always "vault_credential_deleted"
type VaultCredentialValidation string               // Always "vault_credential_validation"
type VaultCredentialArchived string                 // Always "vault_credential.archived"
type VaultCredentialCreated string                  // Always "vault_credential.created"
type VaultCredentialDeleted string                  // Always "vault_credential.deleted"
type VaultCredentialRefreshFailed string            // Always "vault_credential.refresh_failed"
type VaultDeletedResponseType string                // Always "vault_deleted"
type VaultNotFoundError string                      // Always "vault_not_found_error"
type VaultArchived string                           // Always "vault.archived"
type VaultCreated string                            // Always "vault.created"
type VaultDeleted string                            // Always "vault.deleted"
type View string                                    // Always "view"
type Wait string                                    // Always "wait"
type WebFetch string                                // Always "web_fetch"
type WebFetch20250910 string                        // Always "web_fetch_20250910"
type WebFetch20260209 string                        // Always "web_fetch_20260209"
type WebFetch20260309 string                        // Always "web_fetch_20260309"
type WebFetch20260318 string                        // Always "web_fetch_20260318"
type WebFetchResult string                          // Always "web_fetch_result"
type WebFetchToolResult string                      // Always "web_fetch_tool_result"
type WebFetchToolResultError string                 // Always "web_fetch_tool_result_error"
type WebSearch string                               // Always "web_search"
type WebSearch20250305 string                       // Always "web_search_20250305"
type WebSearch20260209 string                       // Always "web_search_20260209"
type WebSearch20260318 string                       // Always "web_search_20260318"
type WebSearchResult string                         // Always "web_search_result"
type WebSearchResultLocation string                 // Always "web_search_result_location"
type WebSearchToolResult string                     // Always "web_search_tool_result"
type WebSearchToolResultError string                // Always "web_search_tool_result_error"
type WebhookEndpointDeleted string                  // Always "webhook_endpoint_deleted"
type Work string                                    // Always "work"
type WorkHeartbeat string                           // Always "work_heartbeat"
type WorkQueueStats string                          // Always "work_queue_stats"
type WorkflowRunCreated string                      // Always "workflow_run.created"
type WorkflowRunError string                        // Always "workflow_run.error"
type WorkflowRunPhaseEnded string                   // Always "workflow_run.phase_ended"
type WorkflowRunPhaseStarted string                 // Always "workflow_run.phase_started"
type WorkflowRunStatusEnded string                  // Always "workflow_run.status_ended"
type WorkflowRunStatusIdle string                   // Always "workflow_run.status_idle"
type WorkflowRunStatusRunning string                // Always "workflow_run.status_running"
type Workspace string                               // Always "workspace"
type WorkspaceArchivedError string                  // Always "workspace_archived_error"
type WorkspaceMember string                         // Always "workspace_member"
type WorkspaceMemberDeleted string                  // Always "workspace_member_deleted"
type WorkspaceRateLimit string                      // Always "workspace_rate_limit"
type Write string                                   // Always "write"
type Xhigh string                                   // Always "xhigh"
type Zoom string                                    // Always "zoom"

func (c Adaptive) Default() Adaptive                             { return "adaptive" }
func (c Advisor) Default() Advisor                               { return "advisor" }
func (c Advisor20260301) Default() Advisor20260301               { return "advisor_20260301" }
func (c AdvisorMessage) Default() AdvisorMessage                 { return "advisor_message" }
func (c AdvisorRedactedResult) Default() AdvisorRedactedResult   { return "advisor_redacted_result" }
func (c AdvisorResult) Default() AdvisorResult                   { return "advisor_result" }
func (c AdvisorToolResult) Default() AdvisorToolResult           { return "advisor_tool_result" }
func (c AdvisorToolResultError) Default() AdvisorToolResultError { return "advisor_tool_result_error" }
func (c Agent) Default() Agent                                   { return "agent" }
func (c AgentArchivedError) Default() AgentArchivedError         { return "agent_archived_error" }
func (c AgentToolset20260401) Default() AgentToolset20260401     { return "agent_toolset_20260401" }
func (c AgentWithOverrides) Default() AgentWithOverrides         { return "agent_with_overrides" }
func (c AgentArchived) Default() AgentArchived                   { return "agent.archived" }
func (c AgentCreated) Default() AgentCreated                     { return "agent.created" }
func (c AgentCustomToolUse) Default() AgentCustomToolUse         { return "agent.custom_tool_use" }
func (c AgentDeleted) Default() AgentDeleted                     { return "agent.deleted" }
func (c AgentMCPToolResult) Default() AgentMCPToolResult         { return "agent.mcp_tool_result" }
func (c AgentMCPToolUse) Default() AgentMCPToolUse               { return "agent.mcp_tool_use" }
func (c AgentMessage) Default() AgentMessage                     { return "agent.message" }
func (c AgentThinking) Default() AgentThinking                   { return "agent.thinking" }
func (c AgentThreadContextCompacted) Default() AgentThreadContextCompacted {
	return "agent.thread_context_compacted"
}
func (c AgentThreadMessageReceived) Default() AgentThreadMessageReceived {
	return "agent.thread_message_received"
}
func (c AgentThreadMessageSent) Default() AgentThreadMessageSent { return "agent.thread_message_sent" }
func (c AgentToolResult) Default() AgentToolResult               { return "agent.tool_result" }
func (c AgentToolUse) Default() AgentToolUse                     { return "agent.tool_use" }
func (c AgentUpdated) Default() AgentUpdated                     { return "agent.updated" }
func (c All) Default() All                                       { return "all" }
func (c AllConnectors) Default() AllConnectors                   { return "all_connectors" }
func (c Allow) Default() Allow                                   { return "allow" }
func (c AlwaysAllow) Default() AlwaysAllow                       { return "always_allow" }
func (c AlwaysAsk) Default() AlwaysAsk                           { return "always_ask" }
func (c Anthropic) Default() Anthropic                           { return "anthropic" }
func (c Any) Default() Any                                       { return "any" }
func (c APIActor) Default() APIActor                             { return "api_actor" }
func (c APIError) Default() APIError                             { return "api_error" }
func (c APIKey) Default() APIKey                                 { return "api_key" }
func (c ApplicationPDF) Default() ApplicationPDF                 { return "application/pdf" }
func (c Approximate) Default() Approximate                       { return "approximate" }
func (c Ask) Default() Ask                                       { return "ask" }
func (c Assistant) Default() Assistant                           { return "assistant" }
func (c Attached) Default() Attached                             { return "attached" }
func (c AuthenticationError) Default() AuthenticationError       { return "authentication_error" }
func (c Auto) Default() Auto                                     { return "auto" }
func (c AWS) Default() AWS                                       { return "aws" }
func (c AWSRole) Default() AWSRole                               { return "aws_role" }
func (c Azure) Default() Azure                                   { return "azure" }
func (c Ban) Default() Ban                                       { return "ban" }
func (c Base64) Default() Base64                                 { return "base64" }
func (c Bash) Default() Bash                                     { return "bash" }
func (c Bash20241022) Default() Bash20241022                     { return "bash_20241022" }
func (c Bash20250124) Default() Bash20250124                     { return "bash_20250124" }
func (c BashCodeExecutionOutput) Default() BashCodeExecutionOutput {
	return "bash_code_execution_output"
}
func (c BashCodeExecutionResult) Default() BashCodeExecutionResult {
	return "bash_code_execution_result"
}
func (c BashCodeExecutionToolResult) Default() BashCodeExecutionToolResult {
	return "bash_code_execution_tool_result"
}
func (c BashCodeExecutionToolResultError) Default() BashCodeExecutionToolResultError {
	return "bash_code_execution_tool_result_error"
}
func (c Batch) Default() Batch                                   { return "batch" }
func (c BetweenTools) Default() BetweenTools                     { return "between_tools" }
func (c BillingError) Default() BillingError                     { return "billing_error" }
func (c Branch) Default() Branch                                 { return "branch" }
func (c Browser) Default() Browser                               { return "browser" }
func (c BrowserState) Default() BrowserState                     { return "browser_state" }
func (c BrowserToolset20260801) Default() BrowserToolset20260801 { return "browser_toolset_20260801" }
func (c BudgetReached) Default() BudgetReached                   { return "budget_reached" }
func (c Caller) Default() Caller                                 { return "caller" }
func (c Canceled) Default() Canceled                             { return "canceled" }
func (c CharLocation) Default() CharLocation                     { return "char_location" }
func (c CitationsDelta) Default() CitationsDelta                 { return "citations_delta" }
func (c ClearThinking20251015) Default() ClearThinking20251015   { return "clear_thinking_20251015" }
func (c ClearToolUses20250919) Default() ClearToolUses20250919   { return "clear_tool_uses_20250919" }
func (c ClientSecretBasic) Default() ClientSecretBasic           { return "client_secret_basic" }
func (c ClientSecretPost) Default() ClientSecretPost             { return "client_secret_post" }
func (c CloseTab) Default() CloseTab                             { return "close_tab" }
func (c Cloud) Default() Cloud                                   { return "cloud" }
func (c CodeExecution) Default() CodeExecution                   { return "code_execution" }
func (c CodeExecution20250522) Default() CodeExecution20250522   { return "code_execution_20250522" }
func (c CodeExecution20250825) Default() CodeExecution20250825   { return "code_execution_20250825" }
func (c CodeExecution20260120) Default() CodeExecution20260120   { return "code_execution_20260120" }
func (c CodeExecution20260521) Default() CodeExecution20260521   { return "code_execution_20260521" }
func (c CodeExecutionOutput) Default() CodeExecutionOutput       { return "code_execution_output" }
func (c CodeExecutionResult) Default() CodeExecutionResult       { return "code_execution_result" }
func (c CodeExecutionToolResult) Default() CodeExecutionToolResult {
	return "code_execution_tool_result"
}
func (c CodeExecutionToolResultError) Default() CodeExecutionToolResultError {
	return "code_execution_tool_result_error"
}
func (c Commit) Default() Commit                         { return "commit" }
func (c Compact20260112) Default() Compact20260112       { return "compact_20260112" }
func (c Compaction) Default() Compaction                 { return "compaction" }
func (c CompactionDelta) Default() CompactionDelta       { return "compaction_delta" }
func (c Completed) Default() Completed                   { return "completed" }
func (c ComplianceSettings) Default() ComplianceSettings { return "compliance_settings" }
func (c Computer) Default() Computer                     { return "computer" }
func (c Computer20241022) Default() Computer20241022     { return "computer_20241022" }
func (c Computer20250124) Default() Computer20250124     { return "computer_20250124" }
func (c Computer20251124) Default() Computer20251124     { return "computer_20251124" }
func (c ComputerToolset20260801) Default() ComputerToolset20260801 {
	return "computer_toolset_20260801"
}
func (c ConflictError) Default() ConflictError               { return "conflict_error" }
func (c Connector) Default() Connector                       { return "connector" }
func (c ConnectorScope) Default() ConnectorScope             { return "connector_scope" }
func (c ConnectorTool) Default() ConnectorTool               { return "connector_tool" }
func (c ContainerUpload) Default() ContainerUpload           { return "container_upload" }
func (c Content) Default() Content                           { return "content" }
func (c ContentBlockDelta) Default() ContentBlockDelta       { return "content_block_delta" }
func (c ContentBlockLocation) Default() ContentBlockLocation { return "content_block_location" }
func (c ContentBlockStart) Default() ContentBlockStart       { return "content_block_start" }
func (c ContentBlockStop) Default() ContentBlockStop         { return "content_block_stop" }
func (c ContentDelta) Default() ContentDelta                 { return "content_delta" }
func (c ContentSha256) Default() ContentSha256               { return "content_sha256" }
func (c Continue) Default() Continue                         { return "continue" }
func (c Coordinate) Default() Coordinate                     { return "coordinate" }
func (c Coordinator) Default() Coordinator                   { return "coordinator" }
func (c Create) Default() Create                             { return "create" }
func (c CreateNew) Default() CreateNew                       { return "create_new" }
func (c CredentialHostUnreachableError) Default() CredentialHostUnreachableError {
	return "credential_host_unreachable_error"
}
func (c Cron) Default() Cron                                     { return "cron" }
func (c CursorPosition) Default() CursorPosition                 { return "cursor_position" }
func (c Custom) Default() Custom                                 { return "custom" }
func (c Default) Default() Default                               { return "default" }
func (c Delete) Default() Delete                                 { return "delete" }
func (c Deny) Default() Deny                                     { return "deny" }
func (c Deployment) Default() Deployment                         { return "deployment" }
func (c DeploymentRun) Default() DeploymentRun                   { return "deployment_run" }
func (c DeploymentRunFailed) Default() DeploymentRunFailed       { return "deployment_run.failed" }
func (c DeploymentRunStarted) Default() DeploymentRunStarted     { return "deployment_run.started" }
func (c DeploymentRunSucceeded) Default() DeploymentRunSucceeded { return "deployment_run.succeeded" }
func (c DeploymentArchived) Default() DeploymentArchived         { return "deployment.archived" }
func (c DeploymentCreated) Default() DeploymentCreated           { return "deployment.created" }
func (c DeploymentDeleted) Default() DeploymentDeleted           { return "deployment.deleted" }
func (c DeploymentPaused) Default() DeploymentPaused             { return "deployment.paused" }
func (c DeploymentUnpaused) Default() DeploymentUnpaused         { return "deployment.unpaused" }
func (c DeploymentUpdated) Default() DeploymentUpdated           { return "deployment.updated" }
func (c Direct) Default() Direct                                 { return "direct" }
func (c Disabled) Default() Disabled                             { return "disabled" }
func (c Discovery) Default() Discovery                           { return "discovery" }
func (c Document) Default() Document                             { return "document" }
func (c DoubleClick) Default() DoubleClick                       { return "double_click" }
func (c DownloadCompleted) Default() DownloadCompleted           { return "download_completed" }
func (c DownloadFailed) Default() DownloadFailed                 { return "download_failed" }
func (c DownloadStarted) Default() DownloadStarted               { return "download_started" }
func (c Dream) Default() Dream                                   { return "dream" }
func (c Edit) Default() Edit                                     { return "edit" }
func (c Enabled) Default() Enabled                               { return "enabled" }
func (c EncryptedCodeExecutionResult) Default() EncryptedCodeExecutionResult {
	return "encrypted_code_execution_result"
}
func (c EndTurn) Default() EndTurn             { return "end_turn" }
func (c EnrollEndUser) Default() EnrollEndUser { return "enroll_end_user" }
func (c EnrollmentURL) Default() EnrollmentURL { return "enrollment_url" }
func (c Environment) Default() Environment     { return "environment" }
func (c EnvironmentArchivedError) Default() EnvironmentArchivedError {
	return "environment_archived_error"
}
func (c EnvironmentDeletedResponseType) Default() EnvironmentDeletedResponseType {
	return "environment_deleted"
}
func (c EnvironmentNotFoundError) Default() EnvironmentNotFoundError {
	return "environment_not_found_error"
}
func (c EnvironmentVariable) Default() EnvironmentVariable     { return "environment_variable" }
func (c EnvironmentArchived) Default() EnvironmentArchived     { return "environment.archived" }
func (c EnvironmentCreated) Default() EnvironmentCreated       { return "environment.created" }
func (c EnvironmentDeleted) Default() EnvironmentDeleted       { return "environment.deleted" }
func (c EnvironmentUpdated) Default() EnvironmentUpdated       { return "environment.updated" }
func (c Ephemeral) Default() Ephemeral                         { return "ephemeral" }
func (c Error) Default() Error                                 { return "error" }
func (c Errored) Default() Errored                             { return "errored" }
func (c Event) Default() Event                                 { return "event" }
func (c EventDelta) Default() EventDelta                       { return "event_delta" }
func (c EventStart) Default() EventStart                       { return "event_start" }
func (c Except) Default() Except                               { return "except" }
func (c Exclude) Default() Exclude                             { return "exclude" }
func (c Exhausted) Default() Exhausted                         { return "exhausted" }
func (c Expired) Default() Expired                             { return "expired" }
func (c ExplicitURL) Default() ExplicitURL                     { return "explicit_url" }
func (c ExternalKey) Default() ExternalKey                     { return "external_key" }
func (c ExternalKeyDeleted) Default() ExternalKeyDeleted       { return "external_key_deleted" }
func (c ExternalKeyValidation) Default() ExternalKeyValidation { return "external_key_validation" }
func (c Failure) Default() Failure                             { return "failure" }
func (c Fallback) Default() Fallback                           { return "fallback" }
func (c FallbackMessage) Default() FallbackMessage             { return "fallback_message" }
func (c FederatedActor) Default() FederatedActor               { return "federated_actor" }
func (c FederationIssuer) Default() FederationIssuer           { return "federation_issuer" }
func (c FederationRule) Default() FederationRule               { return "federation_rule" }
func (c FederationRuleWorkspace) Default() FederationRuleWorkspace {
	return "federation_rule_workspace"
}
func (c FederationRuleWorkspaceDeleted) Default() FederationRuleWorkspaceDeleted {
	return "federation_rule_workspace_deleted"
}
func (c File) Default() File                                 { return "file" }
func (c FileNotFoundError) Default() FileNotFoundError       { return "file_not_found_error" }
func (c FileUpload) Default() FileUpload                     { return "file_upload" }
func (c Files) Default() Files                               { return "files" }
func (c Find) Default() Find                                 { return "find" }
func (c FormInput) Default() FormInput                       { return "form_input" }
func (c GCP) Default() GCP                                   { return "gcp" }
func (c GetPageText) Default() GetPageText                   { return "get_page_text" }
func (c GitHubRepository) Default() GitHubRepository         { return "github_repository" }
func (c Glob) Default() Glob                                 { return "glob" }
func (c Grep) Default() Grep                                 { return "grep" }
func (c High) Default() High                                 { return "high" }
func (c HoldKey) Default() HoldKey                           { return "hold_key" }
func (c Hover) Default() Hover                               { return "hover" }
func (c Image) Default() Image                               { return "image" }
func (c Include) Default() Include                           { return "include" }
func (c Inline) Default() Inline                             { return "inline" }
func (c InputJSONDelta) Default() InputJSONDelta             { return "input_json_delta" }
func (c InputTokens) Default() InputTokens                   { return "input_tokens" }
func (c Insert) Default() Insert                             { return "insert" }
func (c InvalidAddress) Default() InvalidAddress             { return "invalid_address" }
func (c InvalidRequestError) Default() InvalidRequestError   { return "invalid_request_error" }
func (c Invite) Default() Invite                             { return "invite" }
func (c InviteDeleted) Default() InviteDeleted               { return "invite_deleted" }
func (c JavascriptExec) Default() JavascriptExec             { return "javascript_exec" }
func (c JSONSchema) Default() JSONSchema                     { return "json_schema" }
func (c Key) Default() Key                                   { return "key" }
func (c LeftClick) Default() LeftClick                       { return "left_click" }
func (c LeftClickDrag) Default() LeftClickDrag               { return "left_click_drag" }
func (c LeftMouseDown) Default() LeftMouseDown               { return "left_mouse_down" }
func (c LeftMouseUp) Default() LeftMouseUp                   { return "left_mouse_up" }
func (c Limit) Default() Limit                               { return "limit" }
func (c Limited) Default() Limited                           { return "limited" }
func (c ListTabs) Default() ListTabs                         { return "list_tabs" }
func (c Low) Default() Low                                   { return "low" }
func (c Manual) Default() Manual                             { return "manual" }
func (c Max) Default() Max                                   { return "max" }
func (c MaxWorkflowRunsError) Default() MaxWorkflowRunsError { return "max_workflow_runs_error" }
func (c MCPAuthenticationFailedError) Default() MCPAuthenticationFailedError {
	return "mcp_authentication_failed_error"
}
func (c MCPConnectionFailedError) Default() MCPConnectionFailedError {
	return "mcp_connection_failed_error"
}
func (c MCPEgressBlockedError) Default() MCPEgressBlockedError { return "mcp_egress_blocked_error" }
func (c MCPOAuth) Default() MCPOAuth                           { return "mcp_oauth" }
func (c MCPToolListing) Default() MCPToolListing               { return "mcp_tool_listing" }
func (c MCPToolReference) Default() MCPToolReference           { return "mcp_tool_reference" }
func (c MCPToolResult) Default() MCPToolResult                 { return "mcp_tool_result" }
func (c MCPToolUse) Default() MCPToolUse                       { return "mcp_tool_use" }
func (c MCPToolset) Default() MCPToolset                       { return "mcp_toolset" }
func (c MCPToolsetReference) Default() MCPToolsetReference     { return "mcp_toolset_reference" }
func (c Medium) Default() Medium                               { return "medium" }
func (c Memory) Default() Memory                               { return "memory" }
func (c Memory20250818) Default() Memory20250818               { return "memory_20250818" }
func (c MemoryDeleted) Default() MemoryDeleted                 { return "memory_deleted" }
func (c MemoryPathConflictError) Default() MemoryPathConflictError {
	return "memory_path_conflict_error"
}
func (c MemoryPreconditionFailedError) Default() MemoryPreconditionFailedError {
	return "memory_precondition_failed_error"
}
func (c MemoryPrefix) Default() MemoryPrefix { return "memory_prefix" }
func (c MemoryStore) Default() MemoryStore   { return "memory_store" }
func (c MemoryStoreArchivedError) Default() MemoryStoreArchivedError {
	return "memory_store_archived_error"
}
func (c MemoryStoreDeletedResponseType) Default() MemoryStoreDeletedResponseType {
	return "memory_store_deleted"
}
func (c MemoryStoreArchived) Default() MemoryStoreArchived     { return "memory_store.archived" }
func (c MemoryStoreCreated) Default() MemoryStoreCreated       { return "memory_store.created" }
func (c MemoryStoreDeleted) Default() MemoryStoreDeleted       { return "memory_store.deleted" }
func (c MemoryVersion) Default() MemoryVersion                 { return "memory_version" }
func (c Message) Default() Message                             { return "message" }
func (c MessageBatch) Default() MessageBatch                   { return "message_batch" }
func (c MessageBatchDeleted) Default() MessageBatchDeleted     { return "message_batch_deleted" }
func (c MessageDelta) Default() MessageDelta                   { return "message_delta" }
func (c MessageStart) Default() MessageStart                   { return "message_start" }
func (c MessageStop) Default() MessageStop                     { return "message_stop" }
func (c MessagesChanged) Default() MessagesChanged             { return "messages_changed" }
func (c MiddleClick) Default() MiddleClick                     { return "middle_click" }
func (c Model) Default() Model                                 { return "model" }
func (c ModelChanged) Default() ModelChanged                   { return "model_changed" }
func (c ModelGroup) Default() ModelGroup                       { return "model_group" }
func (c ModelOverloadedError) Default() ModelOverloadedError   { return "model_overloaded_error" }
func (c ModelRateLimitedError) Default() ModelRateLimitedError { return "model_rate_limited_error" }
func (c ModelRequestFailedError) Default() ModelRequestFailedError {
	return "model_request_failed_error"
}
func (c MouseMove) Default() MouseMove                   { return "mouse_move" }
func (c Multiagent20261001) Default() Multiagent20261001 { return "multiagent_20261001" }
func (c Navigate) Default() Navigate                     { return "navigate" }
func (c NewTab) Default() NewTab                         { return "new_tab" }
func (c None) Default() None                             { return "none" }
func (c NotApplied) Default() NotApplied                 { return "not_applied" }
func (c NotFoundError) Default() NotFoundError           { return "not_found_error" }
func (c Object) Default() Object                         { return "object" }
func (c OIDC) Default() OIDC                             { return "oidc" }
func (c Only) Default() Only                             { return "only" }
func (c Organization) Default() Organization             { return "organization" }
func (c OrganizationDisabledError) Default() OrganizationDisabledError {
	return "organization_disabled_error"
}
func (c OrganizationMember) Default() OrganizationMember   { return "organization_member" }
func (c OrganizationService) Default() OrganizationService { return "organization_service" }
func (c OutcomeEvaluation) Default() OutcomeEvaluation     { return "outcome_evaluation" }
func (c OverloadedError) Default() OverloadedError         { return "overloaded_error" }
func (c PageLocation) Default() PageLocation               { return "page_location" }
func (c PermissionError) Default() PermissionError         { return "permission_error" }
func (c Plugin) Default() Plugin                           { return "plugin" }
func (c PluginDeleted) Default() PluginDeleted             { return "plugin_deleted" }
func (c PluginInstallationSetting) Default() PluginInstallationSetting {
	return "plugin_installation_setting"
}
func (c PluginInstallationSettingDeleted) Default() PluginInstallationSettingDeleted {
	return "plugin_installation_setting_deleted"
}
func (c PluginMarketplace) Default() PluginMarketplace { return "plugin_marketplace" }
func (c PluginMarketplaceValidationReport) Default() PluginMarketplaceValidationReport {
	return "plugin_marketplace_validation_report"
}
func (c PluginShare) Default() PluginShare     { return "plugin_share" }
func (c PluginVersion) Default() PluginVersion { return "plugin_version" }
func (c PreviousMessageNotFound) Default() PreviousMessageNotFound {
	return "previous_message_not_found"
}
func (c ProgramError) Default() ProgramError                     { return "program_error" }
func (c RateLimit) Default() RateLimit                           { return "rate_limit" }
func (c RateLimitError) Default() RateLimitError                 { return "rate_limit_error" }
func (c RBACGroup) Default() RBACGroup                           { return "rbac_group" }
func (c RBACGroupDeleted) Default() RBACGroupDeleted             { return "rbac_group_deleted" }
func (c RBACGroupMember) Default() RBACGroupMember               { return "rbac_group_member" }
func (c RBACGroupMemberDeleted) Default() RBACGroupMemberDeleted { return "rbac_group_member_deleted" }
func (c RBACRole) Default() RBACRole                             { return "rbac_role" }
func (c RBACRolePermission) Default() RBACRolePermission         { return "rbac_role_permission" }
func (c Read) Default() Read                                     { return "read" }
func (c ReadConsole) Default() ReadConsole                       { return "read_console" }
func (c ReadNetwork) Default() ReadNetwork                       { return "read_network" }
func (c ReadPage) Default() ReadPage                             { return "read_page" }
func (c Redacted) Default() Redacted                             { return "redacted" }
func (c RedactedThinking) Default() RedactedThinking             { return "redacted_thinking" }
func (c Redeemed) Default() Redeemed                             { return "redeemed" }
func (c Redirect) Default() Redirect                             { return "redirect" }
func (c Ref) Default() Ref                                       { return "ref" }
func (c Refusal) Default() Refusal                               { return "refusal" }
func (c Rename) Default() Rename                                 { return "rename" }
func (c RepositoryAuthenticationError) Default() RepositoryAuthenticationError {
	return "repository_authentication_error"
}
func (c RepositoryCheckoutError) Default() RepositoryCheckoutError {
	return "repository_checkout_error"
}
func (c RepositoryCloneError) Default() RepositoryCloneError { return "repository_clone_error" }
func (c RepositoryForbiddenError) Default() RepositoryForbiddenError {
	return "repository_forbidden_error"
}
func (c RepositoryNotFoundError) Default() RepositoryNotFoundError {
	return "repository_not_found_error"
}
func (c RequiresAction) Default() RequiresAction             { return "requires_action" }
func (c Restriction) Default() Restriction                   { return "restriction" }
func (c RetriesExhausted) Default() RetriesExhausted         { return "retries_exhausted" }
func (c Retrying) Default() Retrying                         { return "retrying" }
func (c RightClick) Default() RightClick                     { return "right_click" }
func (c Schedule) Default() Schedule                         { return "schedule" }
func (c ScopedAPIKeyActor) Default() ScopedAPIKeyActor       { return "scoped_api_key_actor" }
func (c Screenshot) Default() Screenshot                     { return "screenshot" }
func (c Scroll) Default() Scroll                             { return "scroll" }
func (c ScrollTo) Default() ScrollTo                         { return "scroll_to" }
func (c SearchResult) Default() SearchResult                 { return "search_result" }
func (c SearchResultLocation) Default() SearchResultLocation { return "search_result_location" }
func (c SeatTier) Default() SeatTier                         { return "seat_tier" }
func (c Self) Default() Self                                 { return "self" }
func (c SelfHosted) Default() SelfHosted                     { return "self_hosted" }
func (c SelfHostedResourcesUnsupportedError) Default() SelfHostedResourcesUnsupportedError {
	return "self_hosted_resources_unsupported_error"
}
func (c ServerToolUse) Default() ServerToolUse             { return "server_tool_use" }
func (c ServiceAccount) Default() ServiceAccount           { return "service_account" }
func (c ServiceAccountActor) Default() ServiceAccountActor { return "service_account_actor" }
func (c ServiceAccountWorkspaceMember) Default() ServiceAccountWorkspaceMember {
	return "service_account_workspace_member"
}
func (c ServiceAccountWorkspaceMemberDeleted) Default() ServiceAccountWorkspaceMemberDeleted {
	return "service_account_workspace_member_deleted"
}
func (c Session) Default() Session           { return "session" }
func (c SessionActor) Default() SessionActor { return "session_actor" }
func (c SessionCreationRejectedError) Default() SessionCreationRejectedError {
	return "session_creation_rejected_error"
}
func (c SessionDeletedResponseType) Default() SessionDeletedResponseType { return "session_deleted" }
func (c SessionRateLimitedError) Default() SessionRateLimitedError {
	return "session_rate_limited_error"
}
func (c SessionResourceDeleted) Default() SessionResourceDeleted { return "session_resource_deleted" }
func (c SessionResourceNotFoundError) Default() SessionResourceNotFoundError {
	return "session_resource_not_found_error"
}
func (c SessionThread) Default() SessionThread               { return "session_thread" }
func (c SessionArchived) Default() SessionArchived           { return "session.archived" }
func (c SessionBudgetReached) Default() SessionBudgetReached { return "session.budget_reached" }
func (c SessionCreated) Default() SessionCreated             { return "session.created" }
func (c SessionDeleted) Default() SessionDeleted             { return "session.deleted" }
func (c SessionError) Default() SessionError                 { return "session.error" }
func (c SessionIdled) Default() SessionIdled                 { return "session.idled" }
func (c SessionOutcomeEvaluationEnded) Default() SessionOutcomeEvaluationEnded {
	return "session.outcome_evaluation_ended"
}
func (c SessionPending) Default() SessionPending               { return "session.pending" }
func (c SessionRequiresAction) Default() SessionRequiresAction { return "session.requires_action" }
func (c SessionRunning) Default() SessionRunning               { return "session.running" }
func (c SessionStatusIdle) Default() SessionStatusIdle         { return "session.status_idle" }
func (c SessionStatusIdled) Default() SessionStatusIdled       { return "session.status_idled" }
func (c SessionStatusRescheduled) Default() SessionStatusRescheduled {
	return "session.status_rescheduled"
}
func (c SessionStatusRunStarted) Default() SessionStatusRunStarted {
	return "session.status_run_started"
}
func (c SessionStatusRunning) Default() SessionStatusRunning { return "session.status_running" }
func (c SessionStatusTerminated) Default() SessionStatusTerminated {
	return "session.status_terminated"
}
func (c SessionThreadCreated) Default() SessionThreadCreated { return "session.thread_created" }
func (c SessionThreadIdled) Default() SessionThreadIdled     { return "session.thread_idled" }
func (c SessionThreadStatusIdle) Default() SessionThreadStatusIdle {
	return "session.thread_status_idle"
}
func (c SessionThreadStatusRescheduled) Default() SessionThreadStatusRescheduled {
	return "session.thread_status_rescheduled"
}
func (c SessionThreadStatusRunning) Default() SessionThreadStatusRunning {
	return "session.thread_status_running"
}
func (c SessionThreadStatusTerminated) Default() SessionThreadStatusTerminated {
	return "session.thread_status_terminated"
}
func (c SessionThreadTerminated) Default() SessionThreadTerminated {
	return "session.thread_terminated"
}
func (c SessionUpdated) Default() SessionUpdated               { return "session.updated" }
func (c SessionUsage) Default() SessionUsage                   { return "session.usage" }
func (c Sessions) Default() Sessions                           { return "sessions" }
func (c SignatureDelta) Default() SignatureDelta               { return "signature_delta" }
func (c Skill) Default() Skill                                 { return "skill" }
func (c SkillDeleted) Default() SkillDeleted                   { return "skill_deleted" }
func (c SkillNotFoundError) Default() SkillNotFoundError       { return "skill_not_found_error" }
func (c SkillVersion) Default() SkillVersion                   { return "skill_version" }
func (c SkillVersionDeleted) Default() SkillVersionDeleted     { return "skill_version_deleted" }
func (c Skills) Default() Skills                               { return "skills" }
func (c SpanModelRequestEnd) Default() SpanModelRequestEnd     { return "span.model_request_end" }
func (c SpanModelRequestStart) Default() SpanModelRequestStart { return "span.model_request_start" }
func (c SpanOutcomeEvaluationEnd) Default() SpanOutcomeEvaluationEnd {
	return "span.outcome_evaluation_end"
}
func (c SpanOutcomeEvaluationOngoing) Default() SpanOutcomeEvaluationOngoing {
	return "span.outcome_evaluation_ongoing"
}
func (c SpanOutcomeEvaluationStart) Default() SpanOutcomeEvaluationStart {
	return "span.outcome_evaluation_start"
}
func (c SpendLimit) Default() SpendLimit               { return "spend_limit" }
func (c SpendLimitDeleted) Default() SpendLimitDeleted { return "spend_limit_deleted" }
func (c SpendLimitIncreaseRequest) Default() SpendLimitIncreaseRequest {
	return "spend_limit_increase_request"
}
func (c StaticBearer) Default() StaticBearer { return "static_bearer" }
func (c Stopped) Default() Stopped           { return "stopped" }
func (c StrReplace) Default() StrReplace     { return "str_replace" }
func (c StrReplaceBasedEditTool) Default() StrReplaceBasedEditTool {
	return "str_replace_based_edit_tool"
}
func (c StrReplaceEditor) Default() StrReplaceEditor         { return "str_replace_editor" }
func (c Succeeded) Default() Succeeded                       { return "succeeded" }
func (c Success) Default() Success                           { return "success" }
func (c Summarize) Default() Summarize                       { return "summarize" }
func (c SustainedFailure) Default() SustainedFailure         { return "sustained_failure" }
func (c SwitchTab) Default() SwitchTab                       { return "switch_tab" }
func (c System) Default() System                             { return "system" }
func (c SystemChanged) Default() SystemChanged               { return "system_changed" }
func (c SystemMessage) Default() SystemMessage               { return "system.message" }
func (c TabOpened) Default() TabOpened                       { return "tab_opened" }
func (c TelemetryDestination) Default() TelemetryDestination { return "telemetry_destination" }
func (c TelemetryDestinationDeleted) Default() TelemetryDestinationDeleted {
	return "telemetry_destination_deleted"
}
func (c TelemetryDestinationTest) Default() TelemetryDestinationTest {
	return "telemetry_destination_test"
}
func (c TelemetryDestinationWorkspace) Default() TelemetryDestinationWorkspace {
	return "telemetry_destination_workspace"
}
func (c TelemetryDestinationWorkspaceDeleted) Default() TelemetryDestinationWorkspaceDeleted {
	return "telemetry_destination_workspace_deleted"
}
func (c Terminal) Default() Terminal                     { return "terminal" }
func (c Text) Default() Text                             { return "text" }
func (c TextDelta) Default() TextDelta                   { return "text_delta" }
func (c TextEditor20241022) Default() TextEditor20241022 { return "text_editor_20241022" }
func (c TextEditor20250124) Default() TextEditor20250124 { return "text_editor_20250124" }
func (c TextEditor20250429) Default() TextEditor20250429 { return "text_editor_20250429" }
func (c TextEditor20250728) Default() TextEditor20250728 { return "text_editor_20250728" }
func (c TextEditorCodeExecutionCreateResult) Default() TextEditorCodeExecutionCreateResult {
	return "text_editor_code_execution_create_result"
}
func (c TextEditorCodeExecutionStrReplaceResult) Default() TextEditorCodeExecutionStrReplaceResult {
	return "text_editor_code_execution_str_replace_result"
}
func (c TextEditorCodeExecutionToolResult) Default() TextEditorCodeExecutionToolResult {
	return "text_editor_code_execution_tool_result"
}
func (c TextEditorCodeExecutionToolResultError) Default() TextEditorCodeExecutionToolResultError {
	return "text_editor_code_execution_tool_result_error"
}
func (c TextEditorCodeExecutionViewResult) Default() TextEditorCodeExecutionViewResult {
	return "text_editor_code_execution_view_result"
}
func (c TextPlain) Default() TextPlain             { return "text/plain" }
func (c Thinking) Default() Thinking               { return "thinking" }
func (c ThinkingDelta) Default() ThinkingDelta     { return "thinking_delta" }
func (c ThinkingDropped) Default() ThinkingDropped { return "thinking_dropped" }
func (c ThinkingMismatchAllowed) Default() ThinkingMismatchAllowed {
	return "thinking_mismatch_allowed"
}
func (c ThinkingTurns) Default() ThinkingTurns               { return "thinking_turns" }
func (c ThreadLimitError) Default() ThreadLimitError         { return "thread_limit_error" }
func (c TimeoutError) Default() TimeoutError                 { return "timeout_error" }
func (c TokenCount) Default() TokenCount                     { return "token_count" }
func (c Tokens) Default() Tokens                             { return "tokens" }
func (c Tool) Default() Tool                                 { return "tool" }
func (c ToolAddition) Default() ToolAddition                 { return "tool_addition" }
func (c ToolDefinition) Default() ToolDefinition             { return "tool_definition" }
func (c ToolReference) Default() ToolReference               { return "tool_reference" }
func (c ToolRemoval) Default() ToolRemoval                   { return "tool_removal" }
func (c ToolResult) Default() ToolResult                     { return "tool_result" }
func (c ToolSearchToolBm25) Default() ToolSearchToolBm25     { return "tool_search_tool_bm25" }
func (c ToolSearchToolRegex) Default() ToolSearchToolRegex   { return "tool_search_tool_regex" }
func (c ToolSearchToolResult) Default() ToolSearchToolResult { return "tool_search_tool_result" }
func (c ToolSearchToolResultError) Default() ToolSearchToolResultError {
	return "tool_search_tool_result_error"
}
func (c ToolSearchToolSearchResult) Default() ToolSearchToolSearchResult {
	return "tool_search_tool_search_result"
}
func (c ToolUse) Default() ToolUse                           { return "tool_use" }
func (c ToolUses) Default() ToolUses                         { return "tool_uses" }
func (c ToolsChanged) Default() ToolsChanged                 { return "tools_changed" }
func (c TripleClick) Default() TripleClick                   { return "triple_click" }
func (c Tunnel) Default() Tunnel                             { return "tunnel" }
func (c TunnelCertificate) Default() TunnelCertificate       { return "tunnel_certificate" }
func (c TunnelToken) Default() TunnelToken                   { return "tunnel_token" }
func (c Type) Default() Type                                 { return "type" }
func (c Unattached) Default() Unattached                     { return "unattached" }
func (c Unavailable) Default() Unavailable                   { return "unavailable" }
func (c UnknownError) Default() UnknownError                 { return "unknown_error" }
func (c Unrestricted) Default() Unrestricted                 { return "unrestricted" }
func (c UpdateExisting) Default() UpdateExisting             { return "update_existing" }
func (c URL) Default() URL                                   { return "url" }
func (c UsageTier) Default() UsageTier                       { return "usage_tier" }
func (c User) Default() User                                 { return "user" }
func (c UserActor) Default() UserActor                       { return "user_actor" }
func (c UserDeleted) Default() UserDeleted                   { return "user_deleted" }
func (c UserProfile) Default() UserProfile                   { return "user_profile" }
func (c UserProfileActionURL) Default() UserProfileActionURL { return "user_profile_action_url" }
func (c UserProfileDeleted) Default() UserProfileDeleted     { return "user_profile_deleted" }
func (c UserProfileEnforcementActionOpened) Default() UserProfileEnforcementActionOpened {
	return "user_profile.enforcement_action_opened"
}
func (c UserProfileEnforcementActionResolved) Default() UserProfileEnforcementActionResolved {
	return "user_profile.enforcement_action_resolved"
}
func (c UserCustomToolResult) Default() UserCustomToolResult { return "user.custom_tool_result" }
func (c UserDefineOutcome) Default() UserDefineOutcome       { return "user.define_outcome" }
func (c UserInterrupt) Default() UserInterrupt               { return "user.interrupt" }
func (c UserMessage) Default() UserMessage                   { return "user.message" }
func (c UserToolConfirmation) Default() UserToolConfirmation { return "user.tool_confirmation" }
func (c UserToolResult) Default() UserToolResult             { return "user.tool_result" }
func (c Vault) Default() Vault                               { return "vault" }
func (c VaultArchivedError) Default() VaultArchivedError     { return "vault_archived_error" }
func (c VaultCredential) Default() VaultCredential           { return "vault_credential" }
func (c VaultCredentialDeletedResponseType) Default() VaultCredentialDeletedResponseType {
	return "vault_credential_deleted"
}
func (c VaultCredentialValidation) Default() VaultCredentialValidation {
	return "vault_credential_validation"
}
func (c VaultCredentialArchived) Default() VaultCredentialArchived {
	return "vault_credential.archived"
}
func (c VaultCredentialCreated) Default() VaultCredentialCreated { return "vault_credential.created" }
func (c VaultCredentialDeleted) Default() VaultCredentialDeleted { return "vault_credential.deleted" }
func (c VaultCredentialRefreshFailed) Default() VaultCredentialRefreshFailed {
	return "vault_credential.refresh_failed"
}
func (c VaultDeletedResponseType) Default() VaultDeletedResponseType { return "vault_deleted" }
func (c VaultNotFoundError) Default() VaultNotFoundError             { return "vault_not_found_error" }
func (c VaultArchived) Default() VaultArchived                       { return "vault.archived" }
func (c VaultCreated) Default() VaultCreated                         { return "vault.created" }
func (c VaultDeleted) Default() VaultDeleted                         { return "vault.deleted" }
func (c View) Default() View                                         { return "view" }
func (c Wait) Default() Wait                                         { return "wait" }
func (c WebFetch) Default() WebFetch                                 { return "web_fetch" }
func (c WebFetch20250910) Default() WebFetch20250910                 { return "web_fetch_20250910" }
func (c WebFetch20260209) Default() WebFetch20260209                 { return "web_fetch_20260209" }
func (c WebFetch20260309) Default() WebFetch20260309                 { return "web_fetch_20260309" }
func (c WebFetch20260318) Default() WebFetch20260318                 { return "web_fetch_20260318" }
func (c WebFetchResult) Default() WebFetchResult                     { return "web_fetch_result" }
func (c WebFetchToolResult) Default() WebFetchToolResult             { return "web_fetch_tool_result" }
func (c WebFetchToolResultError) Default() WebFetchToolResultError {
	return "web_fetch_tool_result_error"
}
func (c WebSearch) Default() WebSearch                 { return "web_search" }
func (c WebSearch20250305) Default() WebSearch20250305 { return "web_search_20250305" }
func (c WebSearch20260209) Default() WebSearch20260209 { return "web_search_20260209" }
func (c WebSearch20260318) Default() WebSearch20260318 { return "web_search_20260318" }
func (c WebSearchResult) Default() WebSearchResult     { return "web_search_result" }
func (c WebSearchResultLocation) Default() WebSearchResultLocation {
	return "web_search_result_location"
}
func (c WebSearchToolResult) Default() WebSearchToolResult { return "web_search_tool_result" }
func (c WebSearchToolResultError) Default() WebSearchToolResultError {
	return "web_search_tool_result_error"
}
func (c WebhookEndpointDeleted) Default() WebhookEndpointDeleted { return "webhook_endpoint_deleted" }
func (c Work) Default() Work                                     { return "work" }
func (c WorkHeartbeat) Default() WorkHeartbeat                   { return "work_heartbeat" }
func (c WorkQueueStats) Default() WorkQueueStats                 { return "work_queue_stats" }
func (c WorkflowRunCreated) Default() WorkflowRunCreated         { return "workflow_run.created" }
func (c WorkflowRunError) Default() WorkflowRunError             { return "workflow_run.error" }
func (c WorkflowRunPhaseEnded) Default() WorkflowRunPhaseEnded   { return "workflow_run.phase_ended" }
func (c WorkflowRunPhaseStarted) Default() WorkflowRunPhaseStarted {
	return "workflow_run.phase_started"
}
func (c WorkflowRunStatusEnded) Default() WorkflowRunStatusEnded { return "workflow_run.status_ended" }
func (c WorkflowRunStatusIdle) Default() WorkflowRunStatusIdle   { return "workflow_run.status_idle" }
func (c WorkflowRunStatusRunning) Default() WorkflowRunStatusRunning {
	return "workflow_run.status_running"
}
func (c Workspace) Default() Workspace                           { return "workspace" }
func (c WorkspaceArchivedError) Default() WorkspaceArchivedError { return "workspace_archived_error" }
func (c WorkspaceMember) Default() WorkspaceMember               { return "workspace_member" }
func (c WorkspaceMemberDeleted) Default() WorkspaceMemberDeleted { return "workspace_member_deleted" }
func (c WorkspaceRateLimit) Default() WorkspaceRateLimit         { return "workspace_rate_limit" }
func (c Write) Default() Write                                   { return "write" }
func (c Xhigh) Default() Xhigh                                   { return "xhigh" }
func (c Zoom) Default() Zoom                                     { return "zoom" }

func (c Adaptive) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c Advisor) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c Advisor20260301) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c AdvisorMessage) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c AdvisorRedactedResult) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c AdvisorResult) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c AdvisorToolResult) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c AdvisorToolResultError) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c Agent) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c AgentArchivedError) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c AgentToolset20260401) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c AgentWithOverrides) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c AgentArchived) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c AgentCreated) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c AgentCustomToolUse) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c AgentDeleted) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c AgentMCPToolResult) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c AgentMCPToolUse) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c AgentMessage) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c AgentThinking) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c AgentThreadContextCompacted) MarshalJSON() ([]byte, error)          { return marshalString(c) }
func (c AgentThreadMessageReceived) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c AgentThreadMessageSent) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c AgentToolResult) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c AgentToolUse) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c AgentUpdated) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c All) MarshalJSON() ([]byte, error)                                  { return marshalString(c) }
func (c AllConnectors) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c Allow) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c AlwaysAllow) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c AlwaysAsk) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c Anthropic) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c Any) MarshalJSON() ([]byte, error)                                  { return marshalString(c) }
func (c APIActor) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c APIError) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c APIKey) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c ApplicationPDF) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c Approximate) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c Ask) MarshalJSON() ([]byte, error)                                  { return marshalString(c) }
func (c Assistant) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c Attached) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c AuthenticationError) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c Auto) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c AWS) MarshalJSON() ([]byte, error)                                  { return marshalString(c) }
func (c AWSRole) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c Azure) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c Ban) MarshalJSON() ([]byte, error)                                  { return marshalString(c) }
func (c Base64) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c Bash) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c Bash20241022) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c Bash20250124) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c BashCodeExecutionOutput) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c BashCodeExecutionResult) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c BashCodeExecutionToolResult) MarshalJSON() ([]byte, error)          { return marshalString(c) }
func (c BashCodeExecutionToolResultError) MarshalJSON() ([]byte, error)     { return marshalString(c) }
func (c Batch) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c BetweenTools) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c BillingError) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c Branch) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c Browser) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c BrowserState) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c BrowserToolset20260801) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c BudgetReached) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c Caller) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c Canceled) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c CharLocation) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c CitationsDelta) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c ClearThinking20251015) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c ClearToolUses20250919) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c ClientSecretBasic) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c ClientSecretPost) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c CloseTab) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c Cloud) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c CodeExecution) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c CodeExecution20250522) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c CodeExecution20250825) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c CodeExecution20260120) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c CodeExecution20260521) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c CodeExecutionOutput) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c CodeExecutionResult) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c CodeExecutionToolResult) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c CodeExecutionToolResultError) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c Commit) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c Compact20260112) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c Compaction) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c CompactionDelta) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c Completed) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c ComplianceSettings) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c Computer) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c Computer20241022) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c Computer20250124) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c Computer20251124) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c ComputerToolset20260801) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c ConflictError) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c Connector) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c ConnectorScope) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c ConnectorTool) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c ContainerUpload) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c Content) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c ContentBlockDelta) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c ContentBlockLocation) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c ContentBlockStart) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c ContentBlockStop) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c ContentDelta) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c ContentSha256) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c Continue) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c Coordinate) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c Coordinator) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c Create) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c CreateNew) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c CredentialHostUnreachableError) MarshalJSON() ([]byte, error)       { return marshalString(c) }
func (c Cron) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c CursorPosition) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c Custom) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c Default) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c Delete) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c Deny) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c Deployment) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c DeploymentRun) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c DeploymentRunFailed) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c DeploymentRunStarted) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c DeploymentRunSucceeded) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c DeploymentArchived) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c DeploymentCreated) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c DeploymentDeleted) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c DeploymentPaused) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c DeploymentUnpaused) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c DeploymentUpdated) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c Direct) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c Disabled) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c Discovery) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c Document) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c DoubleClick) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c DownloadCompleted) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c DownloadFailed) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c DownloadStarted) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c Dream) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c Edit) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c Enabled) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c EncryptedCodeExecutionResult) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c EndTurn) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c EnrollEndUser) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c EnrollmentURL) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c Environment) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c EnvironmentArchivedError) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c EnvironmentDeletedResponseType) MarshalJSON() ([]byte, error)       { return marshalString(c) }
func (c EnvironmentNotFoundError) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c EnvironmentVariable) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c EnvironmentArchived) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c EnvironmentCreated) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c EnvironmentDeleted) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c EnvironmentUpdated) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c Ephemeral) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c Error) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c Errored) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c Event) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c EventDelta) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c EventStart) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c Except) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c Exclude) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c Exhausted) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c Expired) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c ExplicitURL) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c ExternalKey) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c ExternalKeyDeleted) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c ExternalKeyValidation) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c Failure) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c Fallback) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c FallbackMessage) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c FederatedActor) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c FederationIssuer) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c FederationRule) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c FederationRuleWorkspace) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c FederationRuleWorkspaceDeleted) MarshalJSON() ([]byte, error)       { return marshalString(c) }
func (c File) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c FileNotFoundError) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c FileUpload) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c Files) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c Find) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c FormInput) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c GCP) MarshalJSON() ([]byte, error)                                  { return marshalString(c) }
func (c GetPageText) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c GitHubRepository) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c Glob) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c Grep) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c High) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c HoldKey) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c Hover) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c Image) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c Include) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c Inline) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c InputJSONDelta) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c InputTokens) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c Insert) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c InvalidAddress) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c InvalidRequestError) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c Invite) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c InviteDeleted) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c JavascriptExec) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c JSONSchema) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c Key) MarshalJSON() ([]byte, error)                                  { return marshalString(c) }
func (c LeftClick) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c LeftClickDrag) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c LeftMouseDown) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c LeftMouseUp) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c Limit) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c Limited) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c ListTabs) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c Low) MarshalJSON() ([]byte, error)                                  { return marshalString(c) }
func (c Manual) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c Max) MarshalJSON() ([]byte, error)                                  { return marshalString(c) }
func (c MaxWorkflowRunsError) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c MCPAuthenticationFailedError) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c MCPConnectionFailedError) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c MCPEgressBlockedError) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c MCPOAuth) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c MCPToolListing) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c MCPToolReference) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c MCPToolResult) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c MCPToolUse) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c MCPToolset) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c MCPToolsetReference) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c Medium) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c Memory) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c Memory20250818) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c MemoryDeleted) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c MemoryPathConflictError) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c MemoryPreconditionFailedError) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c MemoryPrefix) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c MemoryStore) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c MemoryStoreArchivedError) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c MemoryStoreDeletedResponseType) MarshalJSON() ([]byte, error)       { return marshalString(c) }
func (c MemoryStoreArchived) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c MemoryStoreCreated) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c MemoryStoreDeleted) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c MemoryVersion) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c Message) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c MessageBatch) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c MessageBatchDeleted) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c MessageDelta) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c MessageStart) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c MessageStop) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c MessagesChanged) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c MiddleClick) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c Model) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c ModelChanged) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c ModelGroup) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c ModelOverloadedError) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c ModelRateLimitedError) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c ModelRequestFailedError) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c MouseMove) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c Multiagent20261001) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c Navigate) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c NewTab) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c None) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c NotApplied) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c NotFoundError) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c Object) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c OIDC) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c Only) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c Organization) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c OrganizationDisabledError) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c OrganizationMember) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c OrganizationService) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c OutcomeEvaluation) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c OverloadedError) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c PageLocation) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c PermissionError) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c Plugin) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c PluginDeleted) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c PluginInstallationSetting) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c PluginInstallationSettingDeleted) MarshalJSON() ([]byte, error)     { return marshalString(c) }
func (c PluginMarketplace) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c PluginMarketplaceValidationReport) MarshalJSON() ([]byte, error)    { return marshalString(c) }
func (c PluginShare) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c PluginVersion) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c PreviousMessageNotFound) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c ProgramError) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c RateLimit) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c RateLimitError) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c RBACGroup) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c RBACGroupDeleted) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c RBACGroupMember) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c RBACGroupMemberDeleted) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c RBACRole) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c RBACRolePermission) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c Read) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c ReadConsole) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c ReadNetwork) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c ReadPage) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c Redacted) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c RedactedThinking) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c Redeemed) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c Redirect) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c Ref) MarshalJSON() ([]byte, error)                                  { return marshalString(c) }
func (c Refusal) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c Rename) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c RepositoryAuthenticationError) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c RepositoryCheckoutError) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c RepositoryCloneError) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c RepositoryForbiddenError) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c RepositoryNotFoundError) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c RequiresAction) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c Restriction) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c RetriesExhausted) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c Retrying) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c RightClick) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c Schedule) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c ScopedAPIKeyActor) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c Screenshot) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c Scroll) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c ScrollTo) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c SearchResult) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c SearchResultLocation) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c SeatTier) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c Self) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c SelfHosted) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c SelfHostedResourcesUnsupportedError) MarshalJSON() ([]byte, error)  { return marshalString(c) }
func (c ServerToolUse) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c ServiceAccount) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c ServiceAccountActor) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c ServiceAccountWorkspaceMember) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c ServiceAccountWorkspaceMemberDeleted) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c Session) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c SessionActor) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c SessionCreationRejectedError) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c SessionDeletedResponseType) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c SessionRateLimitedError) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c SessionResourceDeleted) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c SessionResourceNotFoundError) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c SessionThread) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c SessionArchived) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c SessionBudgetReached) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c SessionCreated) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c SessionDeleted) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c SessionError) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c SessionIdled) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c SessionOutcomeEvaluationEnded) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c SessionPending) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c SessionRequiresAction) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c SessionRunning) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c SessionStatusIdle) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c SessionStatusIdled) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c SessionStatusRescheduled) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c SessionStatusRunStarted) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c SessionStatusRunning) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c SessionStatusTerminated) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c SessionThreadCreated) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c SessionThreadIdled) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c SessionThreadStatusIdle) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c SessionThreadStatusRescheduled) MarshalJSON() ([]byte, error)       { return marshalString(c) }
func (c SessionThreadStatusRunning) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c SessionThreadStatusTerminated) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c SessionThreadTerminated) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c SessionUpdated) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c SessionUsage) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c Sessions) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c SignatureDelta) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c Skill) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c SkillDeleted) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c SkillNotFoundError) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c SkillVersion) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c SkillVersionDeleted) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c Skills) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c SpanModelRequestEnd) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c SpanModelRequestStart) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c SpanOutcomeEvaluationEnd) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c SpanOutcomeEvaluationOngoing) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c SpanOutcomeEvaluationStart) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c SpendLimit) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c SpendLimitDeleted) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c SpendLimitIncreaseRequest) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c StaticBearer) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c Stopped) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c StrReplace) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c StrReplaceBasedEditTool) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c StrReplaceEditor) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c Succeeded) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c Success) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c Summarize) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c SustainedFailure) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c SwitchTab) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c System) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c SystemChanged) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c SystemMessage) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c TabOpened) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c TelemetryDestination) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c TelemetryDestinationDeleted) MarshalJSON() ([]byte, error)          { return marshalString(c) }
func (c TelemetryDestinationTest) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c TelemetryDestinationWorkspace) MarshalJSON() ([]byte, error)        { return marshalString(c) }
func (c TelemetryDestinationWorkspaceDeleted) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c Terminal) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c Text) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c TextDelta) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c TextEditor20241022) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c TextEditor20250124) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c TextEditor20250429) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c TextEditor20250728) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c TextEditorCodeExecutionCreateResult) MarshalJSON() ([]byte, error)  { return marshalString(c) }
func (c TextEditorCodeExecutionStrReplaceResult) MarshalJSON() ([]byte, error) {
	return marshalString(c)
}
func (c TextEditorCodeExecutionToolResult) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c TextEditorCodeExecutionToolResultError) MarshalJSON() ([]byte, error) {
	return marshalString(c)
}
func (c TextEditorCodeExecutionViewResult) MarshalJSON() ([]byte, error)    { return marshalString(c) }
func (c TextPlain) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c Thinking) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c ThinkingDelta) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c ThinkingDropped) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c ThinkingMismatchAllowed) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c ThinkingTurns) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c ThreadLimitError) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c TimeoutError) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c TokenCount) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c Tokens) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c Tool) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c ToolAddition) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c ToolDefinition) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c ToolReference) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c ToolRemoval) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c ToolResult) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c ToolSearchToolBm25) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c ToolSearchToolRegex) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c ToolSearchToolResult) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c ToolSearchToolResultError) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c ToolSearchToolSearchResult) MarshalJSON() ([]byte, error)           { return marshalString(c) }
func (c ToolUse) MarshalJSON() ([]byte, error)                              { return marshalString(c) }
func (c ToolUses) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c ToolsChanged) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c TripleClick) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c Tunnel) MarshalJSON() ([]byte, error)                               { return marshalString(c) }
func (c TunnelCertificate) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c TunnelToken) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c Type) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c Unattached) MarshalJSON() ([]byte, error)                           { return marshalString(c) }
func (c Unavailable) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c UnknownError) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c Unrestricted) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c UpdateExisting) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c URL) MarshalJSON() ([]byte, error)                                  { return marshalString(c) }
func (c UsageTier) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c User) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c UserActor) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c UserDeleted) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c UserProfile) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c UserProfileActionURL) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c UserProfileDeleted) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c UserProfileEnforcementActionOpened) MarshalJSON() ([]byte, error)   { return marshalString(c) }
func (c UserProfileEnforcementActionResolved) MarshalJSON() ([]byte, error) { return marshalString(c) }
func (c UserCustomToolResult) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c UserDefineOutcome) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c UserInterrupt) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c UserMessage) MarshalJSON() ([]byte, error)                          { return marshalString(c) }
func (c UserToolConfirmation) MarshalJSON() ([]byte, error)                 { return marshalString(c) }
func (c UserToolResult) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c Vault) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c VaultArchivedError) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c VaultCredential) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c VaultCredentialDeletedResponseType) MarshalJSON() ([]byte, error)   { return marshalString(c) }
func (c VaultCredentialValidation) MarshalJSON() ([]byte, error)            { return marshalString(c) }
func (c VaultCredentialArchived) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c VaultCredentialCreated) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c VaultCredentialDeleted) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c VaultCredentialRefreshFailed) MarshalJSON() ([]byte, error)         { return marshalString(c) }
func (c VaultDeletedResponseType) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c VaultNotFoundError) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c VaultArchived) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c VaultCreated) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c VaultDeleted) MarshalJSON() ([]byte, error)                         { return marshalString(c) }
func (c View) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c Wait) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c WebFetch) MarshalJSON() ([]byte, error)                             { return marshalString(c) }
func (c WebFetch20250910) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c WebFetch20260209) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c WebFetch20260309) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c WebFetch20260318) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c WebFetchResult) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c WebFetchToolResult) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c WebFetchToolResultError) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c WebSearch) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c WebSearch20250305) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c WebSearch20260209) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c WebSearch20260318) MarshalJSON() ([]byte, error)                    { return marshalString(c) }
func (c WebSearchResult) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c WebSearchResultLocation) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c WebSearchToolResult) MarshalJSON() ([]byte, error)                  { return marshalString(c) }
func (c WebSearchToolResultError) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c WebhookEndpointDeleted) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c Work) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }
func (c WorkHeartbeat) MarshalJSON() ([]byte, error)                        { return marshalString(c) }
func (c WorkQueueStats) MarshalJSON() ([]byte, error)                       { return marshalString(c) }
func (c WorkflowRunCreated) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c WorkflowRunError) MarshalJSON() ([]byte, error)                     { return marshalString(c) }
func (c WorkflowRunPhaseEnded) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c WorkflowRunPhaseStarted) MarshalJSON() ([]byte, error)              { return marshalString(c) }
func (c WorkflowRunStatusEnded) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c WorkflowRunStatusIdle) MarshalJSON() ([]byte, error)                { return marshalString(c) }
func (c WorkflowRunStatusRunning) MarshalJSON() ([]byte, error)             { return marshalString(c) }
func (c Workspace) MarshalJSON() ([]byte, error)                            { return marshalString(c) }
func (c WorkspaceArchivedError) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c WorkspaceMember) MarshalJSON() ([]byte, error)                      { return marshalString(c) }
func (c WorkspaceMemberDeleted) MarshalJSON() ([]byte, error)               { return marshalString(c) }
func (c WorkspaceRateLimit) MarshalJSON() ([]byte, error)                   { return marshalString(c) }
func (c Write) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c Xhigh) MarshalJSON() ([]byte, error)                                { return marshalString(c) }
func (c Zoom) MarshalJSON() ([]byte, error)                                 { return marshalString(c) }

func marshalString[T interface {
	~string
	Constant[T]
}](v T) ([]byte, error) {
	return shimjson.Marshal(paramutil.AddrOrDefault(v))
}
