package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go/v2/internal/apierror"
	"github.com/anthropics/anthropic-sdk-go/v2/packages/param"
	"github.com/anthropics/anthropic-sdk-go/v2/shared"
)

// aliased to make [param.APIUnion] private when embedding
type paramUnion = param.APIUnion

// aliased to make [param.APIObject] private when embedding
type paramObj = param.APIObject

type Error = apierror.Error

// This is an alias to an internal type.
type ErrorCode = apierror.ErrorCode

// The email address already belongs to a user in the organization. The message
// names the user.
//
// Returned with HTTP status 400 and [shared.AlreadyMemberErrorDetails] details.
const ErrorCodeAlreadyMember = apierror.ErrorCodeAlreadyMember

// The tools available after a message exceed a limit on their number or their
// size.
//
// Returned with HTTP status 400 and
// [shared.AvailableToolsLimitExceededErrorDetails] details.
const ErrorCodeAvailableToolsLimitExceeded = apierror.ErrorCodeAvailableToolsLimitExceeded

// The email address's domain is not one of the verified domains of the
// organization's parent organization.
//
// Returned with HTTP status 400 and
// [shared.CannotInviteEmailsWithUnverifiedDomainsErrorDetails] details.
const ErrorCodeCannotInviteEmailsWithUnverifiedDomains = apierror.ErrorCodeCannotInviteEmailsWithUnverifiedDomains

// The user is the organization's only member. Add another member before removing
// this one.
//
// Returned with HTTP status 400 and [shared.CannotRemoveLastMemberErrorDetails]
// details.
const ErrorCodeCannotRemoveLastMember = apierror.ErrorCodeCannotRemoveLastMember

// The API key belongs to the Claude Code workspace, and the user who created it is
// no longer a member of the organization or of that workspace.
//
// Returned with HTTP status 400 and
// [shared.ClaudeCodeKeyCreatorNotMemberErrorDetails] details.
const ErrorCodeClaudeCodeKeyCreatorNotMember = apierror.ErrorCodeClaudeCodeKeyCreatorNotMember

// The request comes from a Claude Code version older than the minimum that the
// requested model or the organization requires. Update Claude Code, then retry.
//
// Returned with HTTP status 400 and [shared.ClaudeCodeVersionTooOldErrorDetails]
// details.
const ErrorCodeClaudeCodeVersionTooOld = apierror.ErrorCodeClaudeCodeVersionTooOld

// The workspace's customer-managed encryption key is disabled or revoked.
// Re-enable the key, then retry.
//
// Returned with HTTP status 400 and [shared.CMEKKeyDisabledErrorDetails] details.
const ErrorCodeCMEKKeyDisabled = apierror.ErrorCodeCMEKKeyDisabled

// The key management service's network access controls blocked the request to the
// workspace's customer-managed encryption key. The key itself may be healthy.
//
// Returned with HTTP status 400 and [shared.CMEKKeyNetworkBlockedErrorDetails]
// details.
const ErrorCodeCMEKKeyNetworkBlocked = apierror.ErrorCodeCMEKKeyNetworkBlocked

// `messages` holds more than one `compaction` block with content. Send only the
// block the conversation returned last.
//
// Returned with HTTP status 400 and [shared.CompactionBlockAmbiguousErrorDetails]
// details.
const ErrorCodeCompactionBlockAmbiguous = apierror.ErrorCodeCompactionBlockAmbiguous

// The `compaction` block is not the first block of the conversation, or is in a
// `system` message.
//
// Returned with HTTP status 400 and [shared.CompactionBlockMisplacedErrorDetails]
// details.
const ErrorCodeCompactionBlockMisplaced = apierror.ErrorCodeCompactionBlockMisplaced

// The `compaction` block's `content` differs from the text its `signature`
// recorded.
//
// Returned with HTTP status 400 and [shared.CompactionContentMismatchErrorDetails]
// details.
const ErrorCodeCompactionContentMismatch = apierror.ErrorCodeCompactionContentMismatch

// The compaction request ends on an `assistant` turn that is not finished.
//
// Returned with HTTP status 400 and [shared.CompactionIncompleteTurnErrorDetails]
// details.
const ErrorCodeCompactionIncompleteTurn = apierror.ErrorCodeCompactionIncompleteTurn

// `messages` holds no `user` or `assistant` content to summarize or adopt.
//
// Returned with HTTP status 400 and
// [shared.CompactionNothingToSummarizeErrorDetails] details.
const ErrorCodeCompactionNothingToSummarize = apierror.ErrorCodeCompactionNothingToSummarize

// The `compaction` block's `signature` does not verify.
//
// Returned with HTTP status 400 and
// [shared.CompactionSignatureInvalidErrorDetails] details.
const ErrorCodeCompactionSignatureInvalid = apierror.ErrorCodeCompactionSignatureInvalid

// The conversation references more distinct tools than a `compaction` block can
// record.
//
// Returned with HTTP status 400 and
// [shared.CompactionTooManyToolReferencesErrorDetails] details.
const ErrorCodeCompactionTooManyToolReferences = apierror.ErrorCodeCompactionTooManyToolReferences

// The `compaction` block's `tool_changes` differs from what was returned with the
// block.
//
// Returned with HTTP status 400 and
// [shared.CompactionToolChangesMismatchErrorDetails] details.
const ErrorCodeCompactionToolChangesMismatch = apierror.ErrorCodeCompactionToolChangesMismatch

// The organization restricts invites to an allowed list of email domains, and the
// list is empty or does not include the email address's domain.
//
// Returned with HTTP status 400 and [shared.EmailDomainNotAllowedErrorDetails]
// details.
const ErrorCodeEmailDomainNotAllowed = apierror.ErrorCodeEmailDomainNotAllowed

// The environment is archived, so it cannot be updated.
//
// Returned with HTTP status 400 and [shared.BetaEnvironmentArchivedErrorDetails]
// details.
const ErrorCodeBetaEnvironmentArchived = apierror.ErrorCodeBetaEnvironmentArchived

// The External Key can't be deleted or reconfigured while the organization's
// Workspace Policy lists it as an allowed key.
//
// Returned with HTTP status 400 and
// [shared.ExternalKeyInWorkspacePolicyErrorDetails] details.
const ErrorCodeExternalKeyInWorkspacePolicy = apierror.ErrorCodeExternalKeyInWorkspacePolicy

// The file has `downloadable: false`. Only files generated by a tool, such as the
// code execution tool, can be downloaded.
//
// Returned with HTTP status 400 and [shared.FileNotDownloadableErrorDetails]
// details.
const ErrorCodeFileNotDownloadable = apierror.ErrorCodeFileNotDownloadable

// `config.type` differs from the environment's stored type, which cannot change
// after creation.
//
// Returned with HTTP status 400 and
// [shared.BetaInvalidConfigTypeChangeErrorDetails] details.
const ErrorCodeBetaInvalidConfigTypeChange = apierror.ErrorCodeBetaInvalidConfigTypeChange

// `metadata` has too many keys, a key or value that is too long, or a reserved
// key. The message names each problem.
//
// Returned with HTTP status 400 and [shared.BetaInvalidMetadataErrorDetails]
// details.
const ErrorCodeBetaInvalidMetadata = apierror.ErrorCodeBetaInvalidMetadata

// A pending invite already exists for the email address. The message names the
// invite.
//
// Returned with HTTP status 400 and [shared.InviteAlreadyExistsErrorDetails]
// details.
const ErrorCodeInviteAlreadyExists = apierror.ErrorCodeInviteAlreadyExists

// The External Key is bound to a user profile, and the Workspace does not have
// that user profile. A Workspace can use only a key bound to its own user profile
// or a key that isn't bound to one.
//
// Returned with HTTP status 400 and [shared.KeyProfileMismatchErrorDetails]
// details.
const ErrorCodeKeyProfileMismatch = apierror.ErrorCodeKeyProfileMismatch

// The organization has reached its limit of active API keys.
//
// Returned with HTTP status 400 and [shared.OrgAPIKeyLimitExceededErrorDetails]
// details.
const ErrorCodeOrgAPIKeyLimitExceeded = apierror.ErrorCodeOrgAPIKeyLimitExceeded

// The organization has reached the member limit that Anthropic set for it. Pending
// invites count toward the limit.
//
// Returned with HTTP status 400 and [shared.OrgMemberLimitExceededErrorDetails]
// details.
const ErrorCodeOrgMemberLimitExceeded = apierror.ErrorCodeOrgMemberLimitExceeded

// Anthropic has disabled the organization for trust and safety reasons. An
// organization disabled for any other reason, such as billing, does not get this
// code. `appeal_url`, when present, is where the organization can appeal.
//
// Returned with HTTP status 400 and [shared.OrganizationOnHoldErrorDetails]
// details.
const ErrorCodeOrganizationOnHold = apierror.ErrorCodeOrganizationOnHold

// The endpoint does not support the organization's type. Each endpoint's
// description says which organization types it supports.
//
// Returned with HTTP status 400 and
// [shared.OrganizationTypeNotSupportedErrorDetails] details.
const ErrorCodeOrganizationTypeNotSupported = apierror.ErrorCodeOrganizationTypeNotSupported

// A role attached to an RBAC Group in the request grants owner access, so this API
// does not add members to that group (directly or through an invite), remove its
// members, or delete it. Make the change in the organization's admin settings in
// Claude.
//
// Returned with HTTP status 400 and
// [shared.RBACGroupGrantsOwnerAccessErrorDetails] details.
const ErrorCodeRBACGroupGrantsOwnerAccess = apierror.ErrorCodeRBACGroupGrantsOwnerAccess

// The API could not finish checking whether an RBAC Group in the request grants
// owner access, so it made no change. Retrying the same request will not help.
// Make the change in the organization's admin settings in Claude.
//
// Returned with HTTP status 400 and
// [shared.RBACGroupOwnerAccessUnconfirmedErrorDetails] details.
const ErrorCodeRBACGroupOwnerAccessUnconfirmed = apierror.ErrorCodeRBACGroupOwnerAccessUnconfirmed

// The request body is not valid for its `Content-Encoding`.
//
// Returned with HTTP status 400 and
// [shared.RequestBodyEncodingInvalidErrorDetails] details.
const ErrorCodeRequestBodyEncodingInvalid = apierror.ErrorCodeRequestBodyEncodingInvalid

// The request's `Content-Encoding` is not supported.
//
// Returned with HTTP status 400 and
// [shared.RequestBodyEncodingUnsupportedErrorDetails] details.
const ErrorCodeRequestBodyEncodingUnsupported = apierror.ErrorCodeRequestBodyEncodingUnsupported

// The version in `plugin_version_id` did not pass its content scan, so it cannot
// be served.
//
// Returned with HTTP status 400 and [shared.BetaScanFailedErrorDetails] details.
const ErrorCodeBetaScanFailed = apierror.ErrorCodeBetaScanFailed

// The request set the channel up, and `settings` was then refused. The channel
// stays set up with every setting unset.
//
// Returned with HTTP status 400 and
// [shared.BetaSlackChannelSetUpSettingsRefusedErrorDetails] details.
const ErrorCodeBetaSlackChannelSetUpSettingsRefused = apierror.ErrorCodeBetaSlackChannelSetUpSettingsRefused

// The organization's file storage quota is exhausted. Delete files to free up
// space, then retry the upload.
//
// Returned with HTTP status 400 and [shared.StorageQuotaExceededErrorDetails]
// details.
const ErrorCodeStorageQuotaExceeded = apierror.ErrorCodeStorageQuotaExceeded

// The response that `previous_message_id` names has already been continued, and
// each response can be continued only once. Continue from the Thread's most recent
// response `id`.
//
// Returned with HTTP status 400 and [shared.ThreadAlreadyContinuedErrorDetails]
// details.
const ErrorCodeThreadAlreadyContinued = apierror.ErrorCodeThreadAlreadyContinued

// Another request is still continuing from the same `previous_message_id`. Wait
// for it to complete and continue from its response `id`, or replay the full
// conversation with `thread: {"type": "create"}`.
//
// Returned with HTTP status 400 and [shared.ThreadContinueInFlightErrorDetails]
// details.
const ErrorCodeThreadContinueInFlight = apierror.ErrorCodeThreadContinueInFlight

// A field that must stay the same across a Thread's turns differs from the value
// stored with the Thread; the message names each such field. Send the stored
// values, or replay the full conversation with `thread: {"type": "create"}` to
// change them.
//
// Returned with HTTP status 400 and [shared.ThreadFingerprintMismatchErrorDetails]
// details.
const ErrorCodeThreadFingerprintMismatch = apierror.ErrorCodeThreadFingerprintMismatch

// The request combines `thread` with a request field or an organization setting
// that Threads do not support; `unsupported`, when present, names it. Send the
// full conversation without `thread`, or change that field.
//
// Returned with HTTP status 400 and [shared.ThreadUnsupportedRequestErrorDetails]
// details.
const ErrorCodeThreadUnsupportedRequest = apierror.ErrorCodeThreadUnsupportedRequest

// A tool made available in `messages` takes a name that an available tool of
// another kind, or of another MCP server, already holds.
//
// Returned with HTTP status 400 and [shared.ToolNameConflictErrorDetails] details.
const ErrorCodeToolNameConflict = apierror.ErrorCodeToolNameConflict

// A tool change or MCP tool listing references a tool or MCP server that is not
// defined at that position in the conversation.
//
// Returned with HTTP status 400 and [shared.ToolReferenceUnresolvedErrorDetails]
// details.
const ErrorCodeToolReferenceUnresolved = apierror.ErrorCodeToolReferenceUnresolved

// The request's `Content-Type` header names a charset other than UTF-8, ISO-8859-1
// or US-ASCII. Send the request with no charset, or with one of these.
//
// Returned with HTTP status 400 and [shared.BetaUnsupportedCharsetErrorDetails]
// details.
const ErrorCodeBetaUnsupportedCharset = apierror.ErrorCodeBetaUnsupportedCharset

// The request depends on a user profile that does not exist in the organization or
// is no longer active. The request either names this user profile by ID, or names
// a resource, such as an External Key, that is bound to it.
//
// Returned with HTTP status 400 and [shared.UserProfileNotFoundErrorDetails]
// details.
const ErrorCodeUserProfileNotFound = apierror.ErrorCodeUserProfileNotFound

// The workspace has reached its limit of active API keys.
//
// Returned with HTTP status 400 and
// [shared.WorkspaceAPIKeyLimitExceededErrorDetails] details.
const ErrorCodeWorkspaceAPIKeyLimitExceeded = apierror.ErrorCodeWorkspaceAPIKeyLimitExceeded

// The workspace has reached the member limit that Anthropic set for it.
//
// Returned with HTTP status 400 and
// [shared.WorkspaceMemberLimitExceededErrorDetails] details.
const ErrorCodeWorkspaceMemberLimitExceeded = apierror.ErrorCodeWorkspaceMemberLimitExceeded

// The request breaks the organization's Workspace Policy. The message says which
// rule it broke; an organization admin can see and change the policy in the
// Console.
//
// Returned with HTTP status 400 and [shared.WorkspacePolicyViolationErrorDetails]
// details.
const ErrorCodeWorkspacePolicyViolation = apierror.ErrorCodeWorkspacePolicyViolation

// The API key lacks the scope the request needs: `read:claude_tag` for a read,
// `write:claude_tag` for a change.
//
// Returned with HTTP status 403 and [shared.BetaAPIKeyScopeMissingErrorDetails]
// details.
const ErrorCodeBetaAPIKeyScopeMissing = apierror.ErrorCodeBetaAPIKeyScopeMissing

// Anthropic has temporarily disabled billing changes for the user account the
// credential belongs to. Contact support.
//
// Returned with HTTP status 403 and [shared.BetaBillingFrozenErrorDetails]
// details.
const ErrorCodeBetaBillingFrozen = apierror.ErrorCodeBetaBillingFrozen

// Changes through the Claude Tag API are not enabled for the organization.
//
// Returned with HTTP status 403 and
// [shared.BetaClaudeTagWritesNotEnabledErrorDetails] details.
const ErrorCodeBetaClaudeTagWritesNotEnabled = apierror.ErrorCodeBetaClaudeTagWritesNotEnabled

// The organization's compliance configuration does not let Claude Tag store
// anything new.
//
// Returned with HTTP status 403 and
// [shared.BetaComplianceConfigurationRestrictedErrorDetails] details.
const ErrorCodeBetaComplianceConfigurationRestricted = apierror.ErrorCodeBetaComplianceConfigurationRestricted

// Anthropic has turned off creating environments for this organization. Existing
// environments still work.
//
// Returned with HTTP status 403 and
// [shared.BetaEnvironmentCreateDisabledErrorDetails] details.
const ErrorCodeBetaEnvironmentCreateDisabled = apierror.ErrorCodeBetaEnvironmentCreateDisabled

// Anthropic has turned off updating environments for this organization. Existing
// environments still work as they are.
//
// Returned with HTTP status 403 and
// [shared.BetaEnvironmentUpdateDisabledErrorDetails] details.
const ErrorCodeBetaEnvironmentUpdateDisabled = apierror.ErrorCodeBetaEnvironmentUpdateDisabled

// The credential is not a scoped API key that is bound at the organization level.
//
// Returned with HTTP status 403 and
// [shared.BetaOrganizationLevelScopedAPIKeyRequiredErrorDetails] details.
const ErrorCodeBetaOrganizationLevelScopedAPIKeyRequired = apierror.ErrorCodeBetaOrganizationLevelScopedAPIKeyRequired

// The organization cannot retrieve, set or delete spend limits with an Admin API
// credential.
//
// Returned with HTTP status 403 and
// [shared.BetaSpendLimitsAdminAPIWritesNotEnabledErrorDetails] details.
const ErrorCodeBetaSpendLimitsAdminAPIWritesNotEnabled = apierror.ErrorCodeBetaSpendLimitsAdminAPIWritesNotEnabled

// `settings.environment_id` names no Environment within the organization's reach.
//
// Returned with HTTP status 404 and [shared.BetaEnvironmentNotFoundErrorDetails]
// details.
const ErrorCodeBetaEnvironmentNotFound = apierror.ErrorCodeBetaEnvironmentNotFound

// The file's `expires_at` has passed, so its content is no longer available.
//
// Returned with HTTP status 404 and [shared.FileExpiredErrorDetails] details.
const ErrorCodeFileExpired = apierror.ErrorCodeFileExpired

// `settings.environment_id` names no Environment within the organization's reach,
// and the request had set the channel up by then. The channel stays set up with
// every setting unset.
//
// Returned with HTTP status 404 and
// [shared.BetaSlackChannelSetUpEnvironmentNotFoundErrorDetails] details.
const ErrorCodeBetaSlackChannelSetUpEnvironmentNotFound = apierror.ErrorCodeBetaSlackChannelSetUpEnvironmentNotFound

// No usable Thread state was found for `previous_message_id`. Replay the full
// conversation with `thread: {"type": "create"}` to start a new Thread.
//
// Returned with HTTP status 404 and [shared.ThreadNotFoundErrorDetails] details.
const ErrorCodeThreadNotFound = apierror.ErrorCodeThreadNotFound

// Another change to the organization's members is in progress. Retry the request
// once it completes.
//
// Returned with HTTP status 409 and
// [shared.ConcurrentMembershipChangeErrorDetails] details.
const ErrorCodeConcurrentMembershipChange = apierror.ErrorCodeConcurrentMembershipChange

// The external key changed while the request was validating it. Send the request
// again.
//
// Returned with HTTP status 409 and [shared.ExternalKeyChangedErrorDetails]
// details.
const ErrorCodeExternalKeyChanged = apierror.ErrorCodeExternalKeyChanged

// The marketplace already holds this name. `plugin_id` is the Plugin that holds
// it, and can be absent.
//
// Returned with HTTP status 409 and [shared.BetaPluginNameTakenErrorDetails]
// details.
const ErrorCodeBetaPluginNameTaken = apierror.ErrorCodeBetaPluginNameTaken

// The content scan of the version in `plugin_version_id` is still running. Retry
// once the scan finishes.
//
// Returned with HTTP status 409 and [shared.BetaScanPendingErrorDetails] details.
const ErrorCodeBetaScanPending = apierror.ErrorCodeBetaScanPending

// One of the Plugin's skills has the name of an organization skill. `skill_name`
// is that name.
//
// Returned with HTTP status 409 and [shared.BetaSkillNameTakenErrorDetails]
// details.
const ErrorCodeBetaSkillNameTaken = apierror.ErrorCodeBetaSkillNameTaken

// A custom skill the request names is stored under an encryption key state the
// workspace cannot currently decrypt.
//
// Returned with HTTP status 409 and
// [shared.SkillReadBlockedUndecryptableErrorDetails] details.
const ErrorCodeSkillReadBlockedUndecryptable = apierror.ErrorCodeSkillReadBlockedUndecryptable

// The request set the channel up, and something else then changed the channel or
// removed it. Nothing of `settings` was stored.
//
// Returned with HTTP status 409 and
// [shared.BetaSlackChannelSetUpSettingsConflictErrorDetails] details.
const ErrorCodeBetaSlackChannelSetUpSettingsConflict = apierror.ErrorCodeBetaSlackChannelSetUpSettingsConflict

// `expected_last_heartbeat` does not match the stored heartbeat, so another worker
// may hold the lease.
//
// Returned with HTTP status 412 and
// [shared.BetaHeartbeatPreconditionFailedErrorDetails] details.
const ErrorCodeBetaHeartbeatPreconditionFailed = apierror.ErrorCodeBetaHeartbeatPreconditionFailed

// The compressed request body exceeds the request size limit once decoded.
//
// Returned with HTTP status 413 and
// [shared.RequestBodyTooLargeAfterDecodingErrorDetails] details.
const ErrorCodeRequestBodyTooLargeAfterDecoding = apierror.ErrorCodeRequestBodyTooLargeAfterDecoding

// The organization has crossed the monthly API usage threshold of its API tier.
//
// Returned with HTTP status 429 and [shared.EnforcedSpendLimitReachedErrorDetails]
// details.
const ErrorCodeEnforcedSpendLimitReached = apierror.ErrorCodeEnforcedSpendLimitReached

// The workspace's encryption settings could not be read, so the request was
// refused. Retry in a few minutes.
//
// Returned with HTTP status 503 and [shared.CMEKContextUnavailableErrorDetails]
// details.
const ErrorCodeCMEKContextUnavailable = apierror.ErrorCodeCMEKContextUnavailable

// The workspace's customer-managed encryption key could not be reached. Retry the
// request; if this continues, check that the key is available in your key
// management service, or contact Anthropic support.
//
// Returned with HTTP status 503 and [shared.CMEKKeyUnreachableErrorDetails]
// details.
const ErrorCodeCMEKKeyUnreachable = apierror.ErrorCodeCMEKKeyUnreachable

// The Plugin version was stored but is not yet usable in claude.ai. `plugin_id`
// and `plugin_version_id` identify what was stored.
//
// Returned with HTTP status 503 and [shared.BetaRegistrationPendingErrorDetails]
// details.
const ErrorCodeBetaRegistrationPending = apierror.ErrorCodeBetaRegistrationPending

// The compaction could not be completed for a transient server reason. Retry the
// request.
//
// Returned with HTTP status 529 and [shared.CompactionUnavailableErrorDetails]
// details.
const ErrorCodeCompactionUnavailable = apierror.ErrorCodeCompactionUnavailable

// This is an alias to an internal type.
type APIErrorObject = shared.APIErrorObject

// This is an alias to an internal type.
type AuthenticationError = shared.AuthenticationError

// This is an alias to an internal type.
type BillingError = shared.BillingError

// The API key belongs to the Claude Code workspace, and the user who created it is
// no longer a member of the organization or of that workspace.
//
// This is an alias to an internal type.
type ClaudeCodeKeyCreatorNotMemberErrorDetails = shared.ClaudeCodeKeyCreatorNotMemberErrorDetails

// The request comes from a Claude Code version older than the minimum that the
// requested model or the organization requires. Update Claude Code, then retry.
//
// This is an alias to an internal type.
type ClaudeCodeVersionTooOldErrorDetails = shared.ClaudeCodeVersionTooOldErrorDetails

// The workspace's encryption settings could not be read, so the request was
// refused. Retry in a few minutes.
//
// This is an alias to an internal type.
type CMEKContextUnavailableErrorDetails = shared.CMEKContextUnavailableErrorDetails

// The workspace's customer-managed encryption key is disabled or revoked.
// Re-enable the key, then retry.
//
// This is an alias to an internal type.
type CMEKKeyDisabledErrorDetails = shared.CMEKKeyDisabledErrorDetails

// The key management service's network access controls blocked the request to the
// workspace's customer-managed encryption key. The key itself may be healthy.
//
// This is an alias to an internal type.
type CMEKKeyNetworkBlockedErrorDetails = shared.CMEKKeyNetworkBlockedErrorDetails

// The workspace's customer-managed encryption key could not be reached. Retry the
// request; if this continues, check that the key is available in your key
// management service, or contact Anthropic support.
//
// This is an alias to an internal type.
type CMEKKeyUnreachableErrorDetails = shared.CMEKKeyUnreachableErrorDetails

// The organization has crossed the monthly API usage threshold of its API tier.
//
// This is an alias to an internal type.
type EnforcedSpendLimitReachedErrorDetails = shared.EnforcedSpendLimitReachedErrorDetails

// This is an alias to an internal type.
type ErrorDetailsUnion = shared.ErrorDetailsUnion

// This is an alias to an internal type.
type ErrorObjectUnion = shared.ErrorObjectUnion

// This is an alias to an internal type.
type ErrorResponse = shared.ErrorResponse

// This is an alias to an internal type.
type ErrorType = shared.ErrorType

// Equals "invalid_request_error"
const ErrorTypeInvalidRequestError = shared.ErrorTypeInvalidRequestError

// Equals "authentication_error"
const ErrorTypeAuthenticationError = shared.ErrorTypeAuthenticationError

// Equals "billing_error"
const ErrorTypeBillingError = shared.ErrorTypeBillingError

// Equals "permission_error"
const ErrorTypePermissionError = shared.ErrorTypePermissionError

// Equals "not_found_error"
const ErrorTypeNotFoundError = shared.ErrorTypeNotFoundError

// Equals "rate_limit_error"
const ErrorTypeRateLimitError = shared.ErrorTypeRateLimitError

// Equals "timeout_error"
const ErrorTypeTimeoutError = shared.ErrorTypeTimeoutError

// Equals "api_error"
const ErrorTypeAPIError = shared.ErrorTypeAPIError

// Equals "overloaded_error"
const ErrorTypeOverloadedError = shared.ErrorTypeOverloadedError

// Equals "memory_precondition_failed_error"
const ErrorTypeMemoryPreconditionFailedError = shared.ErrorTypeMemoryPreconditionFailedError

// Equals "memory_path_conflict_error"
const ErrorTypeMemoryPathConflictError = shared.ErrorTypeMemoryPathConflictError

// Equals "conflict_error"
const ErrorTypeConflictError = shared.ErrorTypeConflictError

// This is an alias to an internal type.
type GatewayTimeoutError = shared.GatewayTimeoutError

// This is an alias to an internal type.
type InvalidRequestError = shared.InvalidRequestError

// This is an alias to an internal type.
type NotFoundError = shared.NotFoundError

// Anthropic has disabled the organization for trust and safety reasons. An
// organization disabled for any other reason, such as billing, does not get this
// code. `appeal_url`, when present, is where the organization can appeal.
//
// This is an alias to an internal type.
type OrganizationOnHoldErrorDetails = shared.OrganizationOnHoldErrorDetails

// The endpoint does not support the organization's type. Each endpoint's
// description says which organization types it supports.
//
// This is an alias to an internal type.
type OrganizationTypeNotSupportedErrorDetails = shared.OrganizationTypeNotSupportedErrorDetails

// This is an alias to an internal type.
type OverloadedError = shared.OverloadedError

// This is an alias to an internal type.
type PermissionError = shared.PermissionError

// This is an alias to an internal type.
type RateLimitError = shared.RateLimitError

// The request body is not valid for its `Content-Encoding`.
//
// This is an alias to an internal type.
type RequestBodyEncodingInvalidErrorDetails = shared.RequestBodyEncodingInvalidErrorDetails

// The request's `Content-Encoding` is not supported.
//
// This is an alias to an internal type.
type RequestBodyEncodingUnsupportedErrorDetails = shared.RequestBodyEncodingUnsupportedErrorDetails

// The compressed request body exceeds the request size limit once decoded.
//
// This is an alias to an internal type.
type RequestBodyTooLargeAfterDecodingErrorDetails = shared.RequestBodyTooLargeAfterDecodingErrorDetails

// The tools available after a message exceed a limit on their number or their
// size.
//
// This is an alias to an internal type.
type AvailableToolsLimitExceededErrorDetails = shared.AvailableToolsLimitExceededErrorDetails

// `messages` holds more than one `compaction` block with content. Send only the
// block the conversation returned last.
//
// This is an alias to an internal type.
type CompactionBlockAmbiguousErrorDetails = shared.CompactionBlockAmbiguousErrorDetails

// The `compaction` block is not the first block of the conversation, or is in a
// `system` message.
//
// This is an alias to an internal type.
type CompactionBlockMisplacedErrorDetails = shared.CompactionBlockMisplacedErrorDetails

// The `compaction` block's `content` differs from the text its `signature`
// recorded.
//
// This is an alias to an internal type.
type CompactionContentMismatchErrorDetails = shared.CompactionContentMismatchErrorDetails

// The compaction request ends on an `assistant` turn that is not finished.
//
// This is an alias to an internal type.
type CompactionIncompleteTurnErrorDetails = shared.CompactionIncompleteTurnErrorDetails

// `messages` holds no `user` or `assistant` content to summarize or adopt.
//
// This is an alias to an internal type.
type CompactionNothingToSummarizeErrorDetails = shared.CompactionNothingToSummarizeErrorDetails

// The `compaction` block's `signature` does not verify.
//
// This is an alias to an internal type.
type CompactionSignatureInvalidErrorDetails = shared.CompactionSignatureInvalidErrorDetails

// The conversation references more distinct tools than a `compaction` block can
// record.
//
// This is an alias to an internal type.
type CompactionTooManyToolReferencesErrorDetails = shared.CompactionTooManyToolReferencesErrorDetails

// The `compaction` block's `tool_changes` differs from what was returned with the
// block.
//
// This is an alias to an internal type.
type CompactionToolChangesMismatchErrorDetails = shared.CompactionToolChangesMismatchErrorDetails

// The response that `previous_message_id` names has already been continued, and
// each response can be continued only once. Continue from the Thread's most recent
// response `id`.
//
// This is an alias to an internal type.
type ThreadAlreadyContinuedErrorDetails = shared.ThreadAlreadyContinuedErrorDetails

// Another request is still continuing from the same `previous_message_id`. Wait
// for it to complete and continue from its response `id`, or replay the full
// conversation with `thread: {"type": "create"}`.
//
// This is an alias to an internal type.
type ThreadContinueInFlightErrorDetails = shared.ThreadContinueInFlightErrorDetails

// A field that must stay the same across a Thread's turns differs from the value
// stored with the Thread; the message names each such field. Send the stored
// values, or replay the full conversation with `thread: {"type": "create"}` to
// change them.
//
// This is an alias to an internal type.
type ThreadFingerprintMismatchErrorDetails = shared.ThreadFingerprintMismatchErrorDetails

// The request combines `thread` with a request field or an organization setting
// that Threads do not support; `unsupported`, when present, names it. Send the
// full conversation without `thread`, or change that field.
//
// This is an alias to an internal type.
type ThreadUnsupportedRequestErrorDetails = shared.ThreadUnsupportedRequestErrorDetails

// This is an alias to an internal type.
type ThreadUnsupportedFeatureName = shared.ThreadUnsupportedFeatureName

// Equals "cache_control"
const ThreadUnsupportedFeatureNameCacheControl = shared.ThreadUnsupportedFeatureNameCacheControl

// Equals "context_management"
const ThreadUnsupportedFeatureNameContextManagement = shared.ThreadUnsupportedFeatureNameContextManagement

// Equals "fallbacks"
const ThreadUnsupportedFeatureNameFallbacks = shared.ThreadUnsupportedFeatureNameFallbacks

// Equals "hipaa"
const ThreadUnsupportedFeatureNameHipaa = shared.ThreadUnsupportedFeatureNameHipaa

// Equals "inference_hooks"
const ThreadUnsupportedFeatureNameInferenceHooks = shared.ThreadUnsupportedFeatureNameInferenceHooks

// Equals "thinking"
const ThreadUnsupportedFeatureNameThinking = shared.ThreadUnsupportedFeatureNameThinking

// Equals "zdr"
const ThreadUnsupportedFeatureNameZdr = shared.ThreadUnsupportedFeatureNameZdr

// A tool made available in `messages` takes a name that an available tool of
// another kind, or of another MCP server, already holds.
//
// This is an alias to an internal type.
type ToolNameConflictErrorDetails = shared.ToolNameConflictErrorDetails

// A tool change or MCP tool listing references a tool or MCP server that is not
// defined at that position in the conversation.
//
// This is an alias to an internal type.
type ToolReferenceUnresolvedErrorDetails = shared.ToolReferenceUnresolvedErrorDetails

// No usable Thread state was found for `previous_message_id`. Replay the full
// conversation with `thread: {"type": "create"}` to start a new Thread.
//
// This is an alias to an internal type.
type ThreadNotFoundErrorDetails = shared.ThreadNotFoundErrorDetails

// A custom skill the request names is stored under an encryption key state the
// workspace cannot currently decrypt.
//
// This is an alias to an internal type.
type SkillReadBlockedUndecryptableErrorDetails = shared.SkillReadBlockedUndecryptableErrorDetails

// The compaction could not be completed for a transient server reason. Retry the
// request.
//
// This is an alias to an internal type.
type CompactionUnavailableErrorDetails = shared.CompactionUnavailableErrorDetails

// The organization's file storage quota is exhausted. Delete files to free up
// space, then retry the upload.
//
// This is an alias to an internal type.
type StorageQuotaExceededErrorDetails = shared.StorageQuotaExceededErrorDetails

// The file has `downloadable: false`. Only files generated by a tool, such as the
// code execution tool, can be downloaded.
//
// This is an alias to an internal type.
type FileNotDownloadableErrorDetails = shared.FileNotDownloadableErrorDetails

// The file's `expires_at` has passed, so its content is no longer available.
//
// This is an alias to an internal type.
type FileExpiredErrorDetails = shared.FileExpiredErrorDetails

// The organization has reached its limit of active API keys.
//
// This is an alias to an internal type.
type OrgAPIKeyLimitExceededErrorDetails = shared.OrgAPIKeyLimitExceededErrorDetails

// The workspace has reached its limit of active API keys.
//
// This is an alias to an internal type.
type WorkspaceAPIKeyLimitExceededErrorDetails = shared.WorkspaceAPIKeyLimitExceededErrorDetails

// The request depends on a user profile that does not exist in the organization or
// is no longer active. The request either names this user profile by ID, or names
// a resource, such as an External Key, that is bound to it.
//
// This is an alias to an internal type.
type UserProfileNotFoundErrorDetails = shared.UserProfileNotFoundErrorDetails

// The External Key can't be deleted or reconfigured while the organization's
// Workspace Policy lists it as an allowed key.
//
// This is an alias to an internal type.
type ExternalKeyInWorkspacePolicyErrorDetails = shared.ExternalKeyInWorkspacePolicyErrorDetails

// The email address already belongs to a user in the organization. The message
// names the user.
//
// This is an alias to an internal type.
type AlreadyMemberErrorDetails = shared.AlreadyMemberErrorDetails

// The email address's domain is not one of the verified domains of the
// organization's parent organization.
//
// This is an alias to an internal type.
type CannotInviteEmailsWithUnverifiedDomainsErrorDetails = shared.CannotInviteEmailsWithUnverifiedDomainsErrorDetails

// The organization restricts invites to an allowed list of email domains, and the
// list is empty or does not include the email address's domain.
//
// This is an alias to an internal type.
type EmailDomainNotAllowedErrorDetails = shared.EmailDomainNotAllowedErrorDetails

// A pending invite already exists for the email address. The message names the
// invite.
//
// This is an alias to an internal type.
type InviteAlreadyExistsErrorDetails = shared.InviteAlreadyExistsErrorDetails

// The organization has reached the member limit that Anthropic set for it. Pending
// invites count toward the limit.
//
// This is an alias to an internal type.
type OrgMemberLimitExceededErrorDetails = shared.OrgMemberLimitExceededErrorDetails

// A role attached to an RBAC Group in the request grants owner access, so this API
// does not add members to that group (directly or through an invite), remove its
// members, or delete it. Make the change in the organization's admin settings in
// Claude.
//
// This is an alias to an internal type.
type RBACGroupGrantsOwnerAccessErrorDetails = shared.RBACGroupGrantsOwnerAccessErrorDetails

// The API could not finish checking whether an RBAC Group in the request grants
// owner access, so it made no change. Retrying the same request will not help.
// Make the change in the organization's admin settings in Claude.
//
// This is an alias to an internal type.
type RBACGroupOwnerAccessUnconfirmedErrorDetails = shared.RBACGroupOwnerAccessUnconfirmedErrorDetails

// The user is the organization's only member. Add another member before removing
// this one.
//
// This is an alias to an internal type.
type CannotRemoveLastMemberErrorDetails = shared.CannotRemoveLastMemberErrorDetails

// Another change to the organization's members is in progress. Retry the request
// once it completes.
//
// This is an alias to an internal type.
type ConcurrentMembershipChangeErrorDetails = shared.ConcurrentMembershipChangeErrorDetails

// The External Key is bound to a user profile, and the Workspace does not have
// that user profile. A Workspace can use only a key bound to its own user profile
// or a key that isn't bound to one.
//
// This is an alias to an internal type.
type KeyProfileMismatchErrorDetails = shared.KeyProfileMismatchErrorDetails

// The request breaks the organization's Workspace Policy. The message says which
// rule it broke; an organization admin can see and change the policy in the
// Console.
//
// This is an alias to an internal type.
type WorkspacePolicyViolationErrorDetails = shared.WorkspacePolicyViolationErrorDetails

// This is an alias to an internal type.
type WorkspacePolicyField = shared.WorkspacePolicyField

// Equals "external_key_id"
const WorkspacePolicyFieldExternalKeyID = shared.WorkspacePolicyFieldExternalKeyID

// The external key changed while the request was validating it. Send the request
// again.
//
// This is an alias to an internal type.
type ExternalKeyChangedErrorDetails = shared.ExternalKeyChangedErrorDetails

// The workspace has reached the member limit that Anthropic set for it.
//
// This is an alias to an internal type.
type WorkspaceMemberLimitExceededErrorDetails = shared.WorkspaceMemberLimitExceededErrorDetails

// Anthropic has turned off creating environments for this organization. Existing
// environments still work.
//
// This is an alias to an internal type.
type BetaEnvironmentCreateDisabledErrorDetails = shared.BetaEnvironmentCreateDisabledErrorDetails

// The environment is archived, so it cannot be updated.
//
// This is an alias to an internal type.
type BetaEnvironmentArchivedErrorDetails = shared.BetaEnvironmentArchivedErrorDetails

// `config.type` differs from the environment's stored type, which cannot change
// after creation.
//
// This is an alias to an internal type.
type BetaInvalidConfigTypeChangeErrorDetails = shared.BetaInvalidConfigTypeChangeErrorDetails

// Anthropic has turned off updating environments for this organization. Existing
// environments still work as they are.
//
// This is an alias to an internal type.
type BetaEnvironmentUpdateDisabledErrorDetails = shared.BetaEnvironmentUpdateDisabledErrorDetails

// `expected_last_heartbeat` does not match the stored heartbeat, so another worker
// may hold the lease.
//
// This is an alias to an internal type.
type BetaHeartbeatPreconditionFailedErrorDetails = shared.BetaHeartbeatPreconditionFailedErrorDetails

// A work item's lease as the server holds it.
//
// This is an alias to an internal type.
type BetaWorkLeaseState = shared.BetaWorkLeaseState

// Current state of the work item
//
// This is an alias to an internal type.
type BetaWorkLeaseStateState = shared.BetaWorkLeaseStateState

// Equals "queued"
const BetaWorkLeaseStateStateQueued = shared.BetaWorkLeaseStateStateQueued

// Equals "starting"
const BetaWorkLeaseStateStateStarting = shared.BetaWorkLeaseStateStateStarting

// Equals "active"
const BetaWorkLeaseStateStateActive = shared.BetaWorkLeaseStateStateActive

// Equals "stopping"
const BetaWorkLeaseStateStateStopping = shared.BetaWorkLeaseStateStateStopping

// Equals "stopped"
const BetaWorkLeaseStateStateStopped = shared.BetaWorkLeaseStateStateStopped

// `metadata` has too many keys, a key or value that is too long, or a reserved
// key. The message names each problem.
//
// This is an alias to an internal type.
type BetaInvalidMetadataErrorDetails = shared.BetaInvalidMetadataErrorDetails

// Anthropic has temporarily disabled billing changes for the user account the
// credential belongs to. Contact support.
//
// This is an alias to an internal type.
type BetaBillingFrozenErrorDetails = shared.BetaBillingFrozenErrorDetails

// The organization cannot retrieve, set or delete spend limits with an Admin API
// credential.
//
// This is an alias to an internal type.
type BetaSpendLimitsAdminAPIWritesNotEnabledErrorDetails = shared.BetaSpendLimitsAdminAPIWritesNotEnabledErrorDetails

// The request's `Content-Type` header names a charset other than UTF-8, ISO-8859-1
// or US-ASCII. Send the request with no charset, or with one of these.
//
// This is an alias to an internal type.
type BetaUnsupportedCharsetErrorDetails = shared.BetaUnsupportedCharsetErrorDetails

// The marketplace already holds this name. `plugin_id` is the Plugin that holds
// it, and can be absent.
//
// This is an alias to an internal type.
type BetaPluginNameTakenErrorDetails = shared.BetaPluginNameTakenErrorDetails

// One of the Plugin's skills has the name of an organization skill. `skill_name`
// is that name.
//
// This is an alias to an internal type.
type BetaSkillNameTakenErrorDetails = shared.BetaSkillNameTakenErrorDetails

// The Plugin version was stored but is not yet usable in claude.ai. `plugin_id`
// and `plugin_version_id` identify what was stored.
//
// This is an alias to an internal type.
type BetaRegistrationPendingErrorDetails = shared.BetaRegistrationPendingErrorDetails

// The version in `plugin_version_id` did not pass its content scan, so it cannot
// be served.
//
// This is an alias to an internal type.
type BetaScanFailedErrorDetails = shared.BetaScanFailedErrorDetails

// The content scan of the version in `plugin_version_id` is still running. Retry
// once the scan finishes.
//
// This is an alias to an internal type.
type BetaScanPendingErrorDetails = shared.BetaScanPendingErrorDetails

// The request set the channel up, and `settings` was then refused. The channel
// stays set up with every setting unset.
//
// This is an alias to an internal type.
type BetaSlackChannelSetUpSettingsRefusedErrorDetails = shared.BetaSlackChannelSetUpSettingsRefusedErrorDetails

// The API key lacks the scope the request needs: `read:claude_tag` for a read,
// `write:claude_tag` for a change.
//
// This is an alias to an internal type.
type BetaAPIKeyScopeMissingErrorDetails = shared.BetaAPIKeyScopeMissingErrorDetails

// Changes through the Claude Tag API are not enabled for the organization.
//
// This is an alias to an internal type.
type BetaClaudeTagWritesNotEnabledErrorDetails = shared.BetaClaudeTagWritesNotEnabledErrorDetails

// The organization's compliance configuration does not let Claude Tag store
// anything new.
//
// This is an alias to an internal type.
type BetaComplianceConfigurationRestrictedErrorDetails = shared.BetaComplianceConfigurationRestrictedErrorDetails

// The credential is not a scoped API key that is bound at the organization level.
//
// This is an alias to an internal type.
type BetaOrganizationLevelScopedAPIKeyRequiredErrorDetails = shared.BetaOrganizationLevelScopedAPIKeyRequiredErrorDetails

// `settings.environment_id` names no Environment within the organization's reach.
//
// This is an alias to an internal type.
type BetaEnvironmentNotFoundErrorDetails = shared.BetaEnvironmentNotFoundErrorDetails

// `settings.environment_id` names no Environment within the organization's reach,
// and the request had set the channel up by then. The channel stays set up with
// every setting unset.
//
// This is an alias to an internal type.
type BetaSlackChannelSetUpEnvironmentNotFoundErrorDetails = shared.BetaSlackChannelSetUpEnvironmentNotFoundErrorDetails

// The request set the channel up, and something else then changed the channel or
// removed it. Nothing of `settings` was stored.
//
// This is an alias to an internal type.
type BetaSlackChannelSetUpSettingsConflictErrorDetails = shared.BetaSlackChannelSetUpSettingsConflictErrorDetails

// Backward-compatible type aliases for renamed types.
// These ensure existing code continues to compile after
// UserLocation and ErrorCode types were consolidated.

type WebSearchToolRequestErrorErrorCode = WebSearchToolResultErrorCode
type WebSearchToolResultErrorErrorCode = WebSearchToolResultErrorCode
type WebSearchTool20250305UserLocationParam = UserLocationParam
type WebSearchTool20260209UserLocationParam = UserLocationParam

const WebSearchToolRequestErrorErrorCodeInvalidToolInput = WebSearchToolResultErrorCodeInvalidToolInput
const WebSearchToolRequestErrorErrorCodeUnavailable = WebSearchToolResultErrorCodeUnavailable
const WebSearchToolRequestErrorErrorCodeMaxUsesExceeded = WebSearchToolResultErrorCodeMaxUsesExceeded
const WebSearchToolRequestErrorErrorCodeTooManyRequests = WebSearchToolResultErrorCodeTooManyRequests
const WebSearchToolRequestErrorErrorCodeQueryTooLong = WebSearchToolResultErrorCodeQueryTooLong
const WebSearchToolRequestErrorErrorCodeRequestTooLarge = WebSearchToolResultErrorCodeRequestTooLarge
const WebSearchToolResultErrorErrorCodeInvalidToolInput = WebSearchToolResultErrorCodeInvalidToolInput
const WebSearchToolResultErrorErrorCodeUnavailable = WebSearchToolResultErrorCodeUnavailable
const WebSearchToolResultErrorErrorCodeMaxUsesExceeded = WebSearchToolResultErrorCodeMaxUsesExceeded
const WebSearchToolResultErrorErrorCodeTooManyRequests = WebSearchToolResultErrorCodeTooManyRequests
const WebSearchToolResultErrorErrorCodeQueryTooLong = WebSearchToolResultErrorCodeQueryTooLong
const WebSearchToolResultErrorErrorCodeRequestTooLarge = WebSearchToolResultErrorCodeRequestTooLarge

// Deprecated: use [BetaCurrencyUSD]. Kept so code written before the USD
// initialism casing continues to compile.
const BetaCurrencyUsd = BetaCurrencyUSD
