package apierror

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/anthropics/anthropic-sdk-go/v2/internal/apijson"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/respjson"
	"github.com/anthropics/anthropic-sdk-go/v2/shared"
	"github.com/tidwall/gjson"
)

// Error represents an error that originates from the API, i.e. when a request is
// made and the API returns a response with a HTTP status code. Other errors are
// not wrapped by this SDK.
type Error struct {
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
	StatusCode  int
	Request     *http.Request
	Response    *http.Response
	RequestID   string
	WorkspaceID string
	// A human-readable explanation of the error. Its wording can change, so don't
	// match on it.
	//
	// `error.message` from the response body, if present.
	Message string
	// Machine-readable detail about the cause of the error. `error_code` names the
	// cause; branch on it rather than on `message`. Absent when the error has no code.
	// Treat an unrecognized `error_code` as the bare `type`.
	//
	// `error.details` from the response body, if present.
	Details shared.ErrorDetailsUnion

	errorType shared.ErrorType
}

// Type returns the error type from the API response body, e.g.
// "rate_limit_error" or "overloaded_error". Returns "" if the
// response body did not contain a recognized error type.
func (r *Error) Type() shared.ErrorType { return r.errorType }

// Returns the unmodified JSON received from the API
func (r Error) RawJSON() string { return r.JSON.raw }

func (r *Error) UnmarshalJSON(data []byte) error {
	if !json.Valid(data) {
		r.JSON.raw = string(data)
		return nil
	}
	body := gjson.ParseBytes(data)
	for path, field := range map[string]any{
		"error.message": &r.Message,
		"error.type":    &r.errorType,
		"error.details": &r.Details,
	} {
		apijson.UnmarshalRoot([]byte(body.Get(path).Raw), field)
	}
	return apijson.UnmarshalRoot(data, r)
}

// UnmarshalAPIJSON marks Error as not apijson-native. See
// [apijson.CustomUnmarshaler].
func (r *Error) UnmarshalAPIJSON(data []byte) error { return r.UnmarshalJSON(data) }

func (r *Error) Error() string {
	msg := fmt.Sprintf("%s %q: %d %s", r.Request.Method, r.Request.URL, r.Response.StatusCode, http.StatusText(r.Response.StatusCode))
	if r.RequestID != "" {
		msg += fmt.Sprintf(" (Request-ID: %s)", r.RequestID)
	}
	if r.WorkspaceID != "" {
		msg += fmt.Sprintf(" (Workspace-ID: %s)", r.WorkspaceID)
	}
	if body := r.JSON.raw; body != "" {
		msg += " " + body
	}
	return msg
}

// ErrorCode returns `error.details.error_code` from the response body, or "" if it
// is absent or not a string.
//
// A stable code that names the cause.
func (r *Error) ErrorCode() ErrorCode {
	if code := gjson.Get(r.RawJSON(), "error.details.error_code"); code.Type == gjson.String {
		return ErrorCode(code.Str)
	}
	return ""
}

// DumpRequest returns the request in HTTP/1.x wire format, with sensitive header
// values redacted.
func (r *Error) DumpRequest(body bool) []byte {
	if r.Request.GetBody != nil {
		r.Request.Body, _ = r.Request.GetBody()
	}
	req := *r.Request
	req.Header = RedactHeaders(r.Request.Header)
	out, _ := httputil.DumpRequestOut(&req, body)
	// The dump may swap in an unread copy of the body; keep it for later reads.
	r.Request.Body = req.Body
	return out
}

// DumpResponse returns the response in HTTP/1.x wire format, with sensitive header
// values redacted.
func (r *Error) DumpResponse(body bool) []byte {
	resp := *r.Response
	resp.Header = RedactHeaders(r.Response.Header)
	out, _ := httputil.DumpResponse(&resp, body)
	// The dump may swap in an unread copy of the body; keep it for later reads.
	r.Response.Body = resp.Body
	return out
}

var sensitiveHeaders = []string{"authorization", "proxy-authorization", "api-key", "x-api-key", "x-amz-security-token", "cookie", "set-cookie"}

// RedactHeaders returns headers with every value of each sensitive header
// replaced by "***", without modifying headers. The result may be headers
// itself, so callers must not modify it.
func RedactHeaders(headers http.Header) http.Header {
	var redacted http.Header
	for _, name := range sensitiveHeaders {
		values := headers.Values(name)
		if len(values) == 0 {
			continue
		}
		if redacted == nil {
			redacted = headers.Clone()
		}
		redacted.Del(name)
		for range values {
			redacted.Add(name, "***")
		}
	}
	if redacted == nil {
		return headers
	}
	return redacted
}

// ErrorCode identifies the specific cause of an [Error]. The API may return codes
// that this version of the library does not know.
type ErrorCode string

const (
	// The email address already belongs to a user in the organization. The message
	// names the user.
	//
	// Returned with HTTP status 400 and [shared.AlreadyMemberErrorDetails] details.
	ErrorCodeAlreadyMember ErrorCode = "already_member"

	// The tools available after a message exceed a limit on their number or their
	// size.
	//
	// Returned with HTTP status 400 and
	// [shared.AvailableToolsLimitExceededErrorDetails] details.
	ErrorCodeAvailableToolsLimitExceeded ErrorCode = "available_tools_limit_exceeded"

	// The email address's domain is not one of the verified domains of the
	// organization's parent organization.
	//
	// Returned with HTTP status 400 and
	// [shared.CannotInviteEmailsWithUnverifiedDomainsErrorDetails] details.
	ErrorCodeCannotInviteEmailsWithUnverifiedDomains ErrorCode = "cannot_invite_emails_with_unverified_domains"

	// The user is the organization's only member. Add another member before removing
	// this one.
	//
	// Returned with HTTP status 400 and [shared.CannotRemoveLastMemberErrorDetails]
	// details.
	ErrorCodeCannotRemoveLastMember ErrorCode = "cannot_remove_last_member"

	// The API key belongs to the Claude Code workspace, and the user who created it is
	// no longer a member of the organization or of that workspace.
	//
	// Returned with HTTP status 400 and
	// [shared.ClaudeCodeKeyCreatorNotMemberErrorDetails] details.
	ErrorCodeClaudeCodeKeyCreatorNotMember ErrorCode = "claude_code_key_creator_not_member"

	// The request comes from a Claude Code version older than the minimum that the
	// requested model or the organization requires. Update Claude Code, then retry.
	//
	// Returned with HTTP status 400 and [shared.ClaudeCodeVersionTooOldErrorDetails]
	// details.
	ErrorCodeClaudeCodeVersionTooOld ErrorCode = "claude_code_version_too_old"

	// The workspace's customer-managed encryption key is disabled or revoked.
	// Re-enable the key, then retry.
	//
	// Returned with HTTP status 400 and [shared.CMEKKeyDisabledErrorDetails] details.
	ErrorCodeCMEKKeyDisabled ErrorCode = "cmek_key_disabled"

	// The key management service's network access controls blocked the request to the
	// workspace's customer-managed encryption key. The key itself may be healthy.
	//
	// Returned with HTTP status 400 and [shared.CMEKKeyNetworkBlockedErrorDetails]
	// details.
	ErrorCodeCMEKKeyNetworkBlocked ErrorCode = "cmek_key_network_blocked"

	// `messages` holds more than one `compaction` block with content. Send only the
	// block the conversation returned last.
	//
	// Returned with HTTP status 400 and [shared.CompactionBlockAmbiguousErrorDetails]
	// details.
	ErrorCodeCompactionBlockAmbiguous ErrorCode = "compaction_block_ambiguous"

	// The `compaction` block is not the first block of the conversation, or is in a
	// `system` message.
	//
	// Returned with HTTP status 400 and [shared.CompactionBlockMisplacedErrorDetails]
	// details.
	ErrorCodeCompactionBlockMisplaced ErrorCode = "compaction_block_misplaced"

	// The `compaction` block's `content` differs from the text its `signature`
	// recorded.
	//
	// Returned with HTTP status 400 and [shared.CompactionContentMismatchErrorDetails]
	// details.
	ErrorCodeCompactionContentMismatch ErrorCode = "compaction_content_mismatch"

	// The compaction request ends on an `assistant` turn that is not finished.
	//
	// Returned with HTTP status 400 and [shared.CompactionIncompleteTurnErrorDetails]
	// details.
	ErrorCodeCompactionIncompleteTurn ErrorCode = "compaction_incomplete_turn"

	// `messages` holds no `user` or `assistant` content to summarize or adopt.
	//
	// Returned with HTTP status 400 and
	// [shared.CompactionNothingToSummarizeErrorDetails] details.
	ErrorCodeCompactionNothingToSummarize ErrorCode = "compaction_nothing_to_summarize"

	// The `compaction` block's `signature` does not verify.
	//
	// Returned with HTTP status 400 and
	// [shared.CompactionSignatureInvalidErrorDetails] details.
	ErrorCodeCompactionSignatureInvalid ErrorCode = "compaction_signature_invalid"

	// The conversation references more distinct tools than a `compaction` block can
	// record.
	//
	// Returned with HTTP status 400 and
	// [shared.CompactionTooManyToolReferencesErrorDetails] details.
	ErrorCodeCompactionTooManyToolReferences ErrorCode = "compaction_too_many_tool_references"

	// The `compaction` block's `tool_changes` differs from what was returned with the
	// block.
	//
	// Returned with HTTP status 400 and
	// [shared.CompactionToolChangesMismatchErrorDetails] details.
	ErrorCodeCompactionToolChangesMismatch ErrorCode = "compaction_tool_changes_mismatch"

	// The organization restricts invites to an allowed list of email domains, and the
	// list is empty or does not include the email address's domain.
	//
	// Returned with HTTP status 400 and [shared.EmailDomainNotAllowedErrorDetails]
	// details.
	ErrorCodeEmailDomainNotAllowed ErrorCode = "email_domain_not_allowed"

	// The environment is archived, so it cannot be updated.
	//
	// Returned with HTTP status 400 and [shared.BetaEnvironmentArchivedErrorDetails]
	// details.
	ErrorCodeBetaEnvironmentArchived ErrorCode = "environment_archived"

	// The External Key can't be deleted or reconfigured while the organization's
	// Workspace Policy lists it as an allowed key.
	//
	// Returned with HTTP status 400 and
	// [shared.ExternalKeyInWorkspacePolicyErrorDetails] details.
	ErrorCodeExternalKeyInWorkspacePolicy ErrorCode = "external_key_in_workspace_policy"

	// The file has `downloadable: false`. Only files generated by a tool, such as the
	// code execution tool, can be downloaded.
	//
	// Returned with HTTP status 400 and [shared.FileNotDownloadableErrorDetails]
	// details.
	ErrorCodeFileNotDownloadable ErrorCode = "file_not_downloadable"

	// `config.type` differs from the environment's stored type, which cannot change
	// after creation.
	//
	// Returned with HTTP status 400 and
	// [shared.BetaInvalidConfigTypeChangeErrorDetails] details.
	ErrorCodeBetaInvalidConfigTypeChange ErrorCode = "invalid_config_type_change"

	// `metadata` has too many keys, a key or value that is too long, or a reserved
	// key. The message names each problem.
	//
	// Returned with HTTP status 400 and [shared.BetaInvalidMetadataErrorDetails]
	// details.
	ErrorCodeBetaInvalidMetadata ErrorCode = "invalid_metadata"

	// A pending invite already exists for the email address. The message names the
	// invite.
	//
	// Returned with HTTP status 400 and [shared.InviteAlreadyExistsErrorDetails]
	// details.
	ErrorCodeInviteAlreadyExists ErrorCode = "invite_already_exists"

	// The External Key is bound to a user profile, and the Workspace does not have
	// that user profile. A Workspace can use only a key bound to its own user profile
	// or a key that isn't bound to one.
	//
	// Returned with HTTP status 400 and [shared.KeyProfileMismatchErrorDetails]
	// details.
	ErrorCodeKeyProfileMismatch ErrorCode = "key_profile_mismatch"

	// The organization has reached its limit of active API keys.
	//
	// Returned with HTTP status 400 and [shared.OrgAPIKeyLimitExceededErrorDetails]
	// details.
	ErrorCodeOrgAPIKeyLimitExceeded ErrorCode = "org_api_key_limit_exceeded"

	// The organization has reached the member limit that Anthropic set for it. Pending
	// invites count toward the limit.
	//
	// Returned with HTTP status 400 and [shared.OrgMemberLimitExceededErrorDetails]
	// details.
	ErrorCodeOrgMemberLimitExceeded ErrorCode = "org_member_limit_exceeded"

	// Anthropic has disabled the organization for trust and safety reasons. An
	// organization disabled for any other reason, such as billing, does not get this
	// code. `appeal_url`, when present, is where the organization can appeal.
	//
	// Returned with HTTP status 400 and [shared.OrganizationOnHoldErrorDetails]
	// details.
	ErrorCodeOrganizationOnHold ErrorCode = "organization_on_hold"

	// The endpoint does not support the organization's type. Each endpoint's
	// description says which organization types it supports.
	//
	// Returned with HTTP status 400 and
	// [shared.OrganizationTypeNotSupportedErrorDetails] details.
	ErrorCodeOrganizationTypeNotSupported ErrorCode = "organization_type_not_supported"

	// A role attached to an RBAC Group in the request grants owner access, so this API
	// does not add members to that group (directly or through an invite), remove its
	// members, or delete it. Make the change in the organization's admin settings in
	// Claude.
	//
	// Returned with HTTP status 400 and
	// [shared.RBACGroupGrantsOwnerAccessErrorDetails] details.
	ErrorCodeRBACGroupGrantsOwnerAccess ErrorCode = "rbac_group_grants_owner_access"

	// The API could not finish checking whether an RBAC Group in the request grants
	// owner access, so it made no change. Retrying the same request will not help.
	// Make the change in the organization's admin settings in Claude.
	//
	// Returned with HTTP status 400 and
	// [shared.RBACGroupOwnerAccessUnconfirmedErrorDetails] details.
	ErrorCodeRBACGroupOwnerAccessUnconfirmed ErrorCode = "rbac_group_owner_access_unconfirmed"

	// The request body is not valid for its `Content-Encoding`.
	//
	// Returned with HTTP status 400 and
	// [shared.RequestBodyEncodingInvalidErrorDetails] details.
	ErrorCodeRequestBodyEncodingInvalid ErrorCode = "request_body_encoding_invalid"

	// The request's `Content-Encoding` is not supported.
	//
	// Returned with HTTP status 400 and
	// [shared.RequestBodyEncodingUnsupportedErrorDetails] details.
	ErrorCodeRequestBodyEncodingUnsupported ErrorCode = "request_body_encoding_unsupported"

	// The version in `plugin_version_id` did not pass its content scan, so it cannot
	// be served.
	//
	// Returned with HTTP status 400 and [shared.BetaScanFailedErrorDetails] details.
	ErrorCodeBetaScanFailed ErrorCode = "scan_failed"

	// The request set the channel up, and `settings` was then refused. The channel
	// stays set up with every setting unset.
	//
	// Returned with HTTP status 400 and
	// [shared.BetaSlackChannelSetUpSettingsRefusedErrorDetails] details.
	ErrorCodeBetaSlackChannelSetUpSettingsRefused ErrorCode = "slack_channel_set_up_settings_refused"

	// The organization's file storage quota is exhausted. Delete files to free up
	// space, then retry the upload.
	//
	// Returned with HTTP status 400 and [shared.StorageQuotaExceededErrorDetails]
	// details.
	ErrorCodeStorageQuotaExceeded ErrorCode = "storage_quota_exceeded"

	// The response that `previous_message_id` names has already been continued, and
	// each response can be continued only once. Continue from the Thread's most recent
	// response `id`.
	//
	// Returned with HTTP status 400 and [shared.ThreadAlreadyContinuedErrorDetails]
	// details.
	ErrorCodeThreadAlreadyContinued ErrorCode = "thread_already_continued"

	// Another request is still continuing from the same `previous_message_id`. Wait
	// for it to complete and continue from its response `id`, or replay the full
	// conversation with `thread: {"type": "create"}`.
	//
	// Returned with HTTP status 400 and [shared.ThreadContinueInFlightErrorDetails]
	// details.
	ErrorCodeThreadContinueInFlight ErrorCode = "thread_continue_in_flight"

	// A field that must stay the same across a Thread's turns differs from the value
	// stored with the Thread; the message names each such field. Send the stored
	// values, or replay the full conversation with `thread: {"type": "create"}` to
	// change them.
	//
	// Returned with HTTP status 400 and [shared.ThreadFingerprintMismatchErrorDetails]
	// details.
	ErrorCodeThreadFingerprintMismatch ErrorCode = "thread_fingerprint_mismatch"

	// The request combines `thread` with a request field or an organization setting
	// that Threads do not support; `unsupported`, when present, names it. Send the
	// full conversation without `thread`, or change that field.
	//
	// Returned with HTTP status 400 and [shared.ThreadUnsupportedRequestErrorDetails]
	// details.
	ErrorCodeThreadUnsupportedRequest ErrorCode = "thread_unsupported_request"

	// A tool made available in `messages` takes a name that an available tool of
	// another kind, or of another MCP server, already holds.
	//
	// Returned with HTTP status 400 and [shared.ToolNameConflictErrorDetails] details.
	ErrorCodeToolNameConflict ErrorCode = "tool_name_conflict"

	// A tool change or MCP tool listing references a tool or MCP server that is not
	// defined at that position in the conversation.
	//
	// Returned with HTTP status 400 and [shared.ToolReferenceUnresolvedErrorDetails]
	// details.
	ErrorCodeToolReferenceUnresolved ErrorCode = "tool_reference_unresolved"

	// The request's `Content-Type` header names a charset other than UTF-8, ISO-8859-1
	// or US-ASCII. Send the request with no charset, or with one of these.
	//
	// Returned with HTTP status 400 and [shared.BetaUnsupportedCharsetErrorDetails]
	// details.
	ErrorCodeBetaUnsupportedCharset ErrorCode = "unsupported_charset"

	// The request depends on a user profile that does not exist in the organization or
	// is no longer active. The request either names this user profile by ID, or names
	// a resource, such as an External Key, that is bound to it.
	//
	// Returned with HTTP status 400 and [shared.UserProfileNotFoundErrorDetails]
	// details.
	ErrorCodeUserProfileNotFound ErrorCode = "user_profile_not_found"

	// The workspace has reached its limit of active API keys.
	//
	// Returned with HTTP status 400 and
	// [shared.WorkspaceAPIKeyLimitExceededErrorDetails] details.
	ErrorCodeWorkspaceAPIKeyLimitExceeded ErrorCode = "workspace_api_key_limit_exceeded"

	// The workspace has reached the member limit that Anthropic set for it.
	//
	// Returned with HTTP status 400 and
	// [shared.WorkspaceMemberLimitExceededErrorDetails] details.
	ErrorCodeWorkspaceMemberLimitExceeded ErrorCode = "workspace_member_limit_exceeded"

	// The request breaks the organization's Workspace Policy. The message says which
	// rule it broke; an organization admin can see and change the policy in the
	// Console.
	//
	// Returned with HTTP status 400 and [shared.WorkspacePolicyViolationErrorDetails]
	// details.
	ErrorCodeWorkspacePolicyViolation ErrorCode = "workspace_policy_violation"

	// The API key lacks the scope the request needs: `read:claude_tag` for a read,
	// `write:claude_tag` for a change.
	//
	// Returned with HTTP status 403 and [shared.BetaAPIKeyScopeMissingErrorDetails]
	// details.
	ErrorCodeBetaAPIKeyScopeMissing ErrorCode = "api_key_scope_missing"

	// Anthropic has temporarily disabled billing changes for the user account the
	// credential belongs to. Contact support.
	//
	// Returned with HTTP status 403 and [shared.BetaBillingFrozenErrorDetails]
	// details.
	ErrorCodeBetaBillingFrozen ErrorCode = "billing_frozen"

	// Changes through the Claude Tag API are not enabled for the organization.
	//
	// Returned with HTTP status 403 and
	// [shared.BetaClaudeTagWritesNotEnabledErrorDetails] details.
	ErrorCodeBetaClaudeTagWritesNotEnabled ErrorCode = "claude_tag_writes_not_enabled"

	// The organization's compliance configuration does not let Claude Tag store
	// anything new.
	//
	// Returned with HTTP status 403 and
	// [shared.BetaComplianceConfigurationRestrictedErrorDetails] details.
	ErrorCodeBetaComplianceConfigurationRestricted ErrorCode = "compliance_configuration_restricted"

	// Anthropic has turned off creating environments for this organization. Existing
	// environments still work.
	//
	// Returned with HTTP status 403 and
	// [shared.BetaEnvironmentCreateDisabledErrorDetails] details.
	ErrorCodeBetaEnvironmentCreateDisabled ErrorCode = "environment_create_disabled"

	// Anthropic has turned off updating environments for this organization. Existing
	// environments still work as they are.
	//
	// Returned with HTTP status 403 and
	// [shared.BetaEnvironmentUpdateDisabledErrorDetails] details.
	ErrorCodeBetaEnvironmentUpdateDisabled ErrorCode = "environment_update_disabled"

	// The credential is not a scoped API key that is bound at the organization level.
	//
	// Returned with HTTP status 403 and
	// [shared.BetaOrganizationLevelScopedAPIKeyRequiredErrorDetails] details.
	ErrorCodeBetaOrganizationLevelScopedAPIKeyRequired ErrorCode = "organization_level_scoped_api_key_required"

	// The organization cannot retrieve, set or delete spend limits with an Admin API
	// credential.
	//
	// Returned with HTTP status 403 and
	// [shared.BetaSpendLimitsAdminAPIWritesNotEnabledErrorDetails] details.
	ErrorCodeBetaSpendLimitsAdminAPIWritesNotEnabled ErrorCode = "spend_limits_admin_api_writes_not_enabled"

	// `settings.environment_id` names no Environment within the organization's reach.
	//
	// Returned with HTTP status 404 and [shared.BetaEnvironmentNotFoundErrorDetails]
	// details.
	ErrorCodeBetaEnvironmentNotFound ErrorCode = "environment_not_found"

	// The file's `expires_at` has passed, so its content is no longer available.
	//
	// Returned with HTTP status 404 and [shared.FileExpiredErrorDetails] details.
	ErrorCodeFileExpired ErrorCode = "file_expired"

	// `settings.environment_id` names no Environment within the organization's reach,
	// and the request had set the channel up by then. The channel stays set up with
	// every setting unset.
	//
	// Returned with HTTP status 404 and
	// [shared.BetaSlackChannelSetUpEnvironmentNotFoundErrorDetails] details.
	ErrorCodeBetaSlackChannelSetUpEnvironmentNotFound ErrorCode = "slack_channel_set_up_environment_not_found"

	// No usable Thread state was found for `previous_message_id`. Replay the full
	// conversation with `thread: {"type": "create"}` to start a new Thread.
	//
	// Returned with HTTP status 404 and [shared.ThreadNotFoundErrorDetails] details.
	ErrorCodeThreadNotFound ErrorCode = "thread_not_found"

	// Another change to the organization's members is in progress. Retry the request
	// once it completes.
	//
	// Returned with HTTP status 409 and
	// [shared.ConcurrentMembershipChangeErrorDetails] details.
	ErrorCodeConcurrentMembershipChange ErrorCode = "concurrent_membership_change"

	// The external key changed while the request was validating it. Send the request
	// again.
	//
	// Returned with HTTP status 409 and [shared.ExternalKeyChangedErrorDetails]
	// details.
	ErrorCodeExternalKeyChanged ErrorCode = "external_key_changed"

	// The marketplace already holds this name. `plugin_id` is the Plugin that holds
	// it, and can be absent.
	//
	// Returned with HTTP status 409 and [shared.BetaPluginNameTakenErrorDetails]
	// details.
	ErrorCodeBetaPluginNameTaken ErrorCode = "plugin_name_taken"

	// The content scan of the version in `plugin_version_id` is still running. Retry
	// once the scan finishes.
	//
	// Returned with HTTP status 409 and [shared.BetaScanPendingErrorDetails] details.
	ErrorCodeBetaScanPending ErrorCode = "scan_pending"

	// One of the Plugin's skills has the name of an organization skill. `skill_name`
	// is that name.
	//
	// Returned with HTTP status 409 and [shared.BetaSkillNameTakenErrorDetails]
	// details.
	ErrorCodeBetaSkillNameTaken ErrorCode = "skill_name_taken"

	// A custom skill the request names is stored under an encryption key state the
	// workspace cannot currently decrypt.
	//
	// Returned with HTTP status 409 and
	// [shared.SkillReadBlockedUndecryptableErrorDetails] details.
	ErrorCodeSkillReadBlockedUndecryptable ErrorCode = "skill_read_blocked_undecryptable"

	// The request set the channel up, and something else then changed the channel or
	// removed it. Nothing of `settings` was stored.
	//
	// Returned with HTTP status 409 and
	// [shared.BetaSlackChannelSetUpSettingsConflictErrorDetails] details.
	ErrorCodeBetaSlackChannelSetUpSettingsConflict ErrorCode = "slack_channel_set_up_settings_conflict"

	// `expected_last_heartbeat` does not match the stored heartbeat, so another worker
	// may hold the lease.
	//
	// Returned with HTTP status 412 and
	// [shared.BetaHeartbeatPreconditionFailedErrorDetails] details.
	ErrorCodeBetaHeartbeatPreconditionFailed ErrorCode = "heartbeat_precondition_failed"

	// The compressed request body exceeds the request size limit once decoded.
	//
	// Returned with HTTP status 413 and
	// [shared.RequestBodyTooLargeAfterDecodingErrorDetails] details.
	ErrorCodeRequestBodyTooLargeAfterDecoding ErrorCode = "request_body_too_large_after_decoding"

	// The organization has crossed the monthly API usage threshold of its API tier.
	//
	// Returned with HTTP status 429 and [shared.EnforcedSpendLimitReachedErrorDetails]
	// details.
	ErrorCodeEnforcedSpendLimitReached ErrorCode = "enforced_spend_limit_reached"

	// The request would exceed a ramp limit, which caps how fast usage can grow. Retry
	// after the time in the `retry-after` header.
	//
	// Returned with HTTP status 429 and [shared.RampLimitExceededErrorDetails]
	// details.
	ErrorCodeRampLimitExceeded ErrorCode = "ramp_limit_exceeded"

	// The workspace's encryption settings could not be read, so the request was
	// refused. Retry in a few minutes.
	//
	// Returned with HTTP status 503 and [shared.CMEKContextUnavailableErrorDetails]
	// details.
	ErrorCodeCMEKContextUnavailable ErrorCode = "cmek_context_unavailable"

	// The workspace's customer-managed encryption key could not be reached. Retry the
	// request; if this continues, check that the key is available in your key
	// management service, or contact Anthropic support.
	//
	// Returned with HTTP status 503 and [shared.CMEKKeyUnreachableErrorDetails]
	// details.
	ErrorCodeCMEKKeyUnreachable ErrorCode = "cmek_key_unreachable"

	// The Plugin version was stored but is not yet usable in claude.ai. `plugin_id`
	// and `plugin_version_id` identify what was stored.
	//
	// Returned with HTTP status 503 and [shared.BetaRegistrationPendingErrorDetails]
	// details.
	ErrorCodeBetaRegistrationPending ErrorCode = "registration_pending"

	// The compaction could not be completed for a transient server reason. Retry the
	// request.
	//
	// Returned with HTTP status 529 and [shared.CompactionUnavailableErrorDetails]
	// details.
	ErrorCodeCompactionUnavailable ErrorCode = "compaction_unavailable"
)
